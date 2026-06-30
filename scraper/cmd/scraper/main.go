package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"

	"github.com/wesellis/terminal-news/scraper/internal/classifier"
	"github.com/wesellis/terminal-news/scraper/internal/deduplicator"
	"github.com/wesellis/terminal-news/scraper/internal/newsapi"
	"github.com/wesellis/terminal-news/scraper/internal/parser"
	"github.com/wesellis/terminal-news/scraper/internal/storage"
	"github.com/wesellis/terminal-news/scraper/pkg/types"
)

type Aggregator struct {
	db            *sqlx.DB
	storage       *storage.Storage
	feedParser    *parser.FeedParser
	newsAPIClient *newsapi.Client
	deduplicator  *deduplicator.Deduplicator
	classifier    *classifier.Classifier
	cron          *cron.Cron
}

func main() {
	log.Println("Starting Terminal News Aggregator...")

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize aggregator
	agg, err := NewAggregator()
	if err != nil {
		log.Fatalf("Failed to initialize aggregator: %v", err)
	}
	defer agg.Close()

	// Run immediate fetch on startup
	log.Println("Running initial article fetch...")
	agg.FetchAll()

	// Start cron scheduler
	log.Println("Starting cron scheduler...")
	agg.StartScheduler()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down aggregator...")
}

func NewAggregator() (*Aggregator, error) {
	// Get database URL from environment
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	// Connect to database
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Initialize components
	stor := storage.NewStorage(db)
	feedParser := parser.NewFeedParser()
	dedup := deduplicator.NewDeduplicator(stor)
	class := classifier.NewClassifier()

	// Initialize NewsAPI client (optional)
	var newsAPIClient *newsapi.Client
	newsAPIKey := os.Getenv("NEWSAPI_KEY")
	if newsAPIKey != "" {
		newsAPIClient = newsapi.NewClient(newsAPIKey)
		if err := newsAPIClient.ValidateAPIKey(); err != nil {
			log.Printf("WARNING: NewsAPI key invalid or rate limited: %v", err)
			newsAPIClient = nil
		} else {
			log.Println("NewsAPI client initialized successfully")
		}
	} else {
		log.Println("No NewsAPI key provided, skipping NewsAPI integration")
	}

	// Initialize cron scheduler
	c := cron.New()

	return &Aggregator{
		db:            db,
		storage:       stor,
		feedParser:    feedParser,
		newsAPIClient: newsAPIClient,
		deduplicator:  dedup,
		classifier:    class,
		cron:          c,
	}, nil
}

func (a *Aggregator) Close() {
	if a.cron != nil {
		a.cron.Stop()
	}
	if a.db != nil {
		a.db.Close()
	}
}

// StartScheduler starts the cron jobs
func (a *Aggregator) StartScheduler() {
	// Fetch RSS feeds every 15 minutes
	a.cron.AddFunc("*/15 * * * *", func() {
		log.Println("Cron: Fetching RSS feeds...")
		a.FetchRSSFeeds()
	})

	// Fetch NewsAPI every 30 minutes (if available)
	if a.newsAPIClient != nil {
		a.cron.AddFunc("*/30 * * * *", func() {
			log.Println("Cron: Fetching NewsAPI articles...")
			a.FetchNewsAPI()
		})
	}

	// Clean old articles daily at 2 AM
	a.cron.AddFunc("0 2 * * *", func() {
		log.Println("Cron: Cleaning old articles...")
		count, err := a.storage.CleanOldArticles(90)
		if err != nil {
			log.Printf("Failed to clean old articles: %v", err)
		} else {
			log.Printf("Cleaned %d old articles", count)
		}
	})

	// Clear deduplication cache every hour
	a.cron.AddFunc("0 * * * *", func() {
		log.Println("Cron: Clearing deduplication cache...")
		a.deduplicator.ClearCache()
	})

	a.cron.Start()
	log.Println("Scheduler started")
}

// FetchAll fetches from all sources
func (a *Aggregator) FetchAll() {
	start := time.Now()
	totalArticles := 0

	// Fetch RSS feeds
	rssCount := a.FetchRSSFeeds()
	totalArticles += rssCount

	// Fetch NewsAPI
	if a.newsAPIClient != nil {
		newsAPICount := a.FetchNewsAPI()
		totalArticles += newsAPICount
	}

	elapsed := time.Since(start)
	log.Printf("Fetch complete: %d articles in %.2f seconds", totalArticles, elapsed.Seconds())

	// Show storage stats
	a.PrintStats()
}

// FetchRSSFeeds fetches articles from all RSS feeds
func (a *Aggregator) FetchRSSFeeds() int {
	sources := parser.GetFeedSources()
	totalArticles := make([]types.ParsedArticle, 0)

	for _, source := range sources {
		if !source.Enabled {
			continue
		}

		log.Printf("Fetching RSS feed: %s", source.Name)
		articles, err := a.feedParser.ParseFeed(source.URL, source.Name, source.Category)
		if err != nil {
			log.Printf("ERROR: Failed to fetch %s: %v", source.Name, err)
			continue
		}

		totalArticles = append(totalArticles, articles...)

		// Small delay to be respectful
		time.Sleep(200 * time.Millisecond)
	}

	log.Printf("Fetched %d articles from RSS feeds", len(totalArticles))

	// Process articles
	return a.ProcessArticles(totalArticles)
}

// FetchNewsAPI fetches articles from NewsAPI
func (a *Aggregator) FetchNewsAPI() int {
	if a.newsAPIClient == nil {
		return 0
	}

	totalArticles := make([]types.ParsedArticle, 0)

	// Fetch top headlines by category
	articles, err := a.newsAPIClient.FetchMultipleCategories("us", 100)
	if err != nil {
		log.Printf("ERROR: Failed to fetch NewsAPI articles: %v", err)
		return 0
	}

	totalArticles = append(totalArticles, articles...)
	log.Printf("Fetched %d articles from NewsAPI", len(totalArticles))

	// Process articles
	return a.ProcessArticles(totalArticles)
}

// ProcessArticles deduplicates, classifies, and stores articles
func (a *Aggregator) ProcessArticles(articles []types.ParsedArticle) int {
	if len(articles) == 0 {
		return 0
	}

	log.Printf("Processing %d articles...", len(articles))

	// Step 1: Deduplicate
	articles = a.deduplicator.Deduplicate(articles)
	log.Printf("After deduplication: %d unique articles", len(articles))

	// Step 2: Classify
	for i := range articles {
		a.classifier.Classify(&articles[i])
	}

	// Step 3: Store
	stored, err := a.storage.StoreArticles(articles)
	if err != nil {
		log.Printf("ERROR: Failed to store articles: %v", err)
		return 0
	}

	log.Printf("Successfully stored %d articles", stored)
	return stored
}

// PrintStats prints current storage statistics
func (a *Aggregator) PrintStats() {
	totalCount, err := a.storage.GetArticleCount()
	if err != nil {
		log.Printf("Failed to get article count: %v", err)
		return
	}

	sourceCount, err := a.storage.GetArticleCountBySource()
	if err != nil {
		log.Printf("Failed to get source counts: %v", err)
		return
	}

	log.Println("=== Storage Statistics ===")
	log.Printf("Total articles: %d", totalCount)
	log.Println("Articles by source:")
	for source, count := range sourceCount {
		log.Printf("  %s: %d", source, count)
	}
	log.Println("==========================")
}
