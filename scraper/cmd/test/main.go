package main

import (
	"fmt"
	"log"

	"github.com/wesellis/terminal-news/scraper/internal/parser"
)

func main() {
	log.Println("=== RSS Parser Test (No Database Required) ===")

	p := parser.NewFeedParser()

	// Test with a few known working RSS feeds
	testFeeds := []struct {
		url      string
		source   string
		category string
	}{
		{"https://feeds.bbci.co.uk/news/technology/rss.xml", "BBC Tech", "tech"},
		{"https://hnrss.org/frontpage", "HackerNews", "tech"},
		{"https://techcrunch.com/feed/", "TechCrunch", "tech"},
	}

	totalArticles := 0

	for _, feed := range testFeeds {
		log.Printf("\n--- Testing %s ---", feed.source)
		log.Printf("URL: %s", feed.url)

		articles, err := p.ParseFeed(feed.url, feed.source, feed.category)

		if err != nil {
			log.Printf("❌ ERROR: %v", err)
			continue
		}

		log.Printf("✅ SUCCESS: Got %d articles from %s", len(articles), feed.source)
		totalArticles += len(articles)

		if len(articles) > 0 {
			fmt.Printf("\nFirst article from %s:\n", feed.source)
			fmt.Printf("  Title: %s\n", articles[0].Title)
			fmt.Printf("  URL: %s\n", articles[0].URL)
			fmt.Printf("  Published: %s\n", articles[0].PublishedAt)
			fmt.Printf("  Category: %s\n", articles[0].Category)
		}
	}

	log.Printf("\n=== Test Summary ===")
	log.Printf("Total articles fetched: %d", totalArticles)
	if totalArticles > 0 {
		log.Println("✅ RSS Parser is working!")
	} else {
		log.Println("❌ No articles fetched - check network/URLs")
	}
}
