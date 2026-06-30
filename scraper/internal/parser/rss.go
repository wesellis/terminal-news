package parser

import (
	"crypto/md5"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/wesellis/terminal-news/scraper/pkg/types"
)

type FeedParser struct {
	parser *gofeed.Parser
}

func NewFeedParser() *FeedParser {
	return &FeedParser{
		parser: gofeed.NewParser(),
	}
}

// ParseFeed fetches and parses an RSS feed
func (fp *FeedParser) ParseFeed(url, sourceName, defaultCategory string) ([]types.ParsedArticle, error) {
	feed, err := fp.parser.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse feed %s: %w", url, err)
	}

	articles := make([]types.ParsedArticle, 0, len(feed.Items))

	for _, item := range feed.Items {
		article := fp.convertToArticle(item, sourceName, defaultCategory)
		articles = append(articles, article)
	}

	log.Printf("Parsed %d articles from %s", len(articles), sourceName)
	return articles, nil
}

func (fp *FeedParser) convertToArticle(item *gofeed.Item, source, category string) types.ParsedArticle {
	// Extract content (try description first, then content)
	content := ""
	if item.Description != "" {
		content = fp.cleanHTML(item.Description)
	} else if item.Content != "" {
		content = fp.cleanHTML(item.Content)
	}

	// Get author
	author := ""
	if item.Author != nil {
		author = item.Author.Name
	}

	// Get published time
	publishedAt := time.Now()
	if item.PublishedParsed != nil {
		publishedAt = *item.PublishedParsed
	} else if item.UpdatedParsed != nil {
		publishedAt = *item.UpdatedParsed
	}

	// Extract image URL
	imageURL := ""
	if item.Image != nil {
		imageURL = item.Image.URL
	} else if item.Enclosures != nil && len(item.Enclosures) > 0 {
		// Check for image enclosures
		for _, enc := range item.Enclosures {
			if strings.HasPrefix(enc.Type, "image/") {
				imageURL = enc.URL
				break
			}
		}
	}

	// Generate external ID from URL or GUID
	externalID := item.GUID
	if externalID == "" {
		externalID = fp.generateHash(item.Link)
	}

	return types.ParsedArticle{
		Title:       fp.cleanTitle(item.Title),
		URL:         item.Link,
		Summary:     fp.truncate(content, 500),
		Content:     content,
		Author:      author,
		Source:      source,
		Category:    category,
		Tags:        item.Categories,
		ImageURL:    imageURL,
		PublishedAt: publishedAt,
		FetchedAt:   time.Now(),
		ExternalID:  externalID,
		FetchSource: "rss",
	}
}

func (fp *FeedParser) cleanHTML(html string) string {
	// Simple HTML tag removal (more sophisticated cleaning can be added)
	text := strings.ReplaceAll(html, "<br>", "\n")
	text = strings.ReplaceAll(text, "<br/>", "\n")
	text = strings.ReplaceAll(text, "<br />", "\n")
	text = strings.ReplaceAll(text, "</p>", "\n\n")

	// Remove all HTML tags
	result := ""
	inTag := false
	for _, r := range text {
		if r == '<' {
			inTag = true
			continue
		}
		if r == '>' {
			inTag = false
			continue
		}
		if !inTag {
			result += string(r)
		}
	}

	// Clean up extra whitespace
	result = strings.TrimSpace(result)
	lines := strings.Split(result, "\n")
	cleaned := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			cleaned = append(cleaned, line)
		}
	}

	return strings.Join(cleaned, "\n")
}

func (fp *FeedParser) cleanTitle(title string) string {
	// Remove excessive whitespace
	title = strings.TrimSpace(title)
	title = strings.Join(strings.Fields(title), " ")
	return title
}

func (fp *FeedParser) truncate(text string, length int) string {
	if len(text) <= length {
		return text
	}

	// Try to truncate at word boundary
	truncated := text[:length]
	lastSpace := strings.LastIndex(truncated, " ")
	if lastSpace > length-50 { // If space is reasonably close to target
		truncated = truncated[:lastSpace]
	}

	return truncated + "..."
}

func (fp *FeedParser) generateHash(text string) string {
	hash := md5.Sum([]byte(text))
	return fmt.Sprintf("%x", hash)
}

// GetFeedSources returns a list of RSS feed sources to scrape
func GetFeedSources() []types.FeedSource {
	return []types.FeedSource{
		// Tech News
		{Name: "TechCrunch", URL: "https://techcrunch.com/feed/", Category: "tech", Enabled: true},
		{Name: "The Verge", URL: "https://www.theverge.com/rss/index.xml", Category: "tech", Enabled: true},
		{Name: "Ars Technica", URL: "https://feeds.arstechnica.com/arstechnica/index", Category: "tech", Enabled: true},
		{Name: "Hacker News", URL: "https://news.ycombinator.com/rss", Category: "tech", Enabled: true},
		{Name: "Wired", URL: "https://www.wired.com/feed/rss", Category: "tech", Enabled: true},
		{Name: "CNET", URL: "https://www.cnet.com/rss/news/", Category: "tech", Enabled: true},

		// World News
		{Name: "BBC News", URL: "http://feeds.bbci.co.uk/news/rss.xml", Category: "world", Enabled: true},
		{Name: "Reuters", URL: "https://www.reutersagency.com/feed/?taxonomy=best-topics&post_type=best", Category: "world", Enabled: true},
		{Name: "Al Jazeera", URL: "https://www.aljazeera.com/xml/rss/all.xml", Category: "world", Enabled: true},
		{Name: "NPR", URL: "https://feeds.npr.org/1001/rss.xml", Category: "world", Enabled: true},

		// Business
		{Name: "Bloomberg", URL: "https://www.bloomberg.com/feed/podcast/the-bloomberg-advantage.xml", Category: "business", Enabled: true},
		{Name: "Financial Times", URL: "https://www.ft.com/?format=rss", Category: "business", Enabled: true},
		{Name: "WSJ", URL: "https://feeds.a.dj.com/rss/RSSWorldNews.xml", Category: "business", Enabled: true},
		{Name: "CNBC", URL: "https://www.cnbc.com/id/100003114/device/rss/rss.html", Category: "business", Enabled: true},

		// Science
		{Name: "Science Daily", URL: "https://www.sciencedaily.com/rss/all.xml", Category: "science", Enabled: true},
		{Name: "Nature", URL: "https://www.nature.com/nature.rss", Category: "science", Enabled: true},
		{Name: "Phys.org", URL: "https://phys.org/rss-feed/", Category: "science", Enabled: true},

		// Entertainment
		{Name: "Variety", URL: "https://variety.com/feed/", Category: "entertainment", Enabled: true},
		{Name: "Hollywood Reporter", URL: "https://www.hollywoodreporter.com/feed/", Category: "entertainment", Enabled: true},

		// General
		{Name: "NY Times", URL: "https://rss.nytimes.com/services/xml/rss/nyt/HomePage.xml", Category: "general", Enabled: true},
		{Name: "Washington Post", URL: "https://feeds.washingtonpost.com/rss/world", Category: "general", Enabled: true},
		{Name: "The Guardian", URL: "https://www.theguardian.com/world/rss", Category: "general", Enabled: true},
		{Name: "CNN", URL: "http://rss.cnn.com/rss/cnn_topstories.rss", Category: "general", Enabled: true},
	}
}
