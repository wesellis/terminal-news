package main

import (
	"fmt"
	"log"

	"github.com/wesellis/terminal-news/scraper/internal/deduplicator"
	"github.com/wesellis/terminal-news/scraper/pkg/types"
)

func main() {
	log.Println("=== Deduplication Test (No Database) ===")

	// Create deduplicator without database (nil storage)
	dedup := deduplicator.NewDeduplicator(nil)

	// Create test articles with duplicates
	articles := []types.ParsedArticle{
		{
			Title:   "Breaking: Major Tech Announcement",
			URL:     "https://example.com/article1",
			Summary: "A major technology company announced a breakthrough today",
			Source:  "TechNews",
		},
		{
			Title:   "Breaking: Major Tech Announcement", // Exact duplicate title
			URL:     "https://example.com/article1",      // Exact duplicate URL
			Summary: "A major technology company announced a breakthrough today",
			Source:  "TechDaily",
		},
		{
			Title:   "Major Tech Announcement Breaks Records", // Similar title
			URL:     "https://example.com/article2",
			Summary: "A major technology company announced a breakthrough today",
			Source:  "NewsHub",
		},
		{
			Title:   "Completely Different Article",
			URL:     "https://example.com/article3",
			Summary: "This is about something totally different",
			Source:  "OtherNews",
		},
	}

	log.Printf("Testing with %d articles (including duplicates)", len(articles))

	// Test deduplication
	unique := dedup.Deduplicate(articles)

	log.Printf("\n✅ Deduplication complete!")
	log.Printf("Input: %d articles", len(articles))
	log.Printf("Output: %d unique articles", len(unique))
	log.Printf("Removed: %d duplicates", len(articles)-len(unique))

	fmt.Println("\nUnique articles:")
	for i, article := range unique {
		fmt.Printf("%d. %s (from %s)\n", i+1, article.Title, article.Source)
	}

	// Test similarity calculation
	fmt.Println("\n=== Similarity Testing ===")
	text1 := "Breaking: Major Tech Announcement"
	text2 := "Major Tech Announcement Breaks Records"
	text3 := "Completely Different Article"

	sim1 := dedup.CalculateSimilarity(text1, text2)
	sim2 := dedup.CalculateSimilarity(text1, text3)

	fmt.Printf("Similarity between similar titles: %.2f (should be high)\n", sim1)
	fmt.Printf("Similarity between different titles: %.2f (should be low)\n", sim2)

	if sim1 > 0.5 && sim2 < 0.3 {
		log.Println("\n✅ Similarity algorithm working correctly!")
	} else {
		log.Println("\n⚠️ Similarity thresholds may need adjustment")
	}

	// Test cache
	fmt.Println("\n=== Cache Testing ===")
	log.Printf("Cache size: %d", dedup.GetCacheSize())

	// Add more articles
	moreArticles := []types.ParsedArticle{
		{Title: "New Article 1", URL: "https://example.com/new1", ExternalID: "new1", FetchSource: "test"},
		{Title: "New Article 2", URL: "https://example.com/new2", ExternalID: "new2", FetchSource: "test"},
	}

	dedup.Deduplicate(moreArticles)
	log.Printf("Cache size after adding articles: %d", dedup.GetCacheSize())

	// Clear cache
	dedup.ClearCache()
	log.Printf("Cache size after clearing: %d", dedup.GetCacheSize())

	log.Println("\n✅ Deduplication test complete!")
}
