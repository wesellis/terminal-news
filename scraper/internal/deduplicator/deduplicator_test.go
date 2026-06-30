package deduplicator_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wesellis/terminal-news/scraper/internal/deduplicator"
	"github.com/wesellis/terminal-news/scraper/pkg/types"
)

func TestDeduplicate(t *testing.T) {
	dedup := deduplicator.NewDeduplicator(nil)

	articles := []types.ParsedArticle{
		{
			Title:       "Breaking: Major Tech News",
			URL:         "https://example.com/news1",
			ExternalID:  "id1",
			FetchSource: "rss",
		},
		{
			Title:       "Breaking: Major Tech News", // Duplicate title
			URL:         "https://example.com/news2",
			ExternalID:  "id2",
			FetchSource: "rss",
		},
		{
			Title:       "Different News Story",
			URL:         "https://example.com/news3",
			ExternalID:  "id3",
			FetchSource: "rss",
		},
		{
			Title:       "Another Story",
			URL:         "https://example.com/news1", // Duplicate URL
			ExternalID:  "id4",
			FetchSource: "rss",
		},
	}

	unique := dedup.Deduplicate(articles)
	assert.Len(t, unique, 2, "Should have 2 unique articles")
}

func TestCalculateSimilarity(t *testing.T) {
	dedup := deduplicator.NewDeduplicator(nil)

	tests := []struct {
		text1    string
		text2    string
		expected float64
		minSim   float64
	}{
		{
			text1:    "Apple announces new iPhone",
			text2:    "Apple announces new iPhone",
			expected: 1.0,
			minSim:   0.99,
		},
		{
			text1:  "Apple announces new iPhone with AI features",
			text2:  "Apple announces new iPhone",
			minSim: 0.5,
		},
		{
			text1:  "Completely different article about politics",
			text2:  "Tech company releases new product",
			minSim: 0.0,
		},
	}

	for _, test := range tests {
		sim := dedup.CalculateSimilarity(test.text1, test.text2)
		if test.expected > 0 {
			assert.InDelta(t, test.expected, sim, 0.1, "Similarity mismatch")
		} else {
			assert.GreaterOrEqual(t, sim, test.minSim, "Minimum similarity not met")
		}
	}
}

func TestIsDuplicate(t *testing.T) {
	dedup := deduplicator.NewDeduplicator(nil)

	article1 := types.ParsedArticle{
		Title:       "Original Article",
		URL:         "https://example.com/original",
		ExternalID:  "original1",
		FetchSource: "rss",
	}

	// First article should not be duplicate
	assert.False(t, dedup.IsDuplicate(article1))

	// Add to deduplicator by calling Deduplicate
	dedup.Deduplicate([]types.ParsedArticle{article1})

	// Same article should now be duplicate
	assert.True(t, dedup.IsDuplicate(article1))
}

func TestClearCache(t *testing.T) {
	dedup := deduplicator.NewDeduplicator(nil)

	articles := []types.ParsedArticle{
		{
			Title:       "Article 1",
			URL:         "https://example.com/1",
			ExternalID:  "1",
			FetchSource: "rss",
		},
	}

	dedup.Deduplicate(articles)
	assert.Greater(t, dedup.GetCacheSize(), 0, "Cache should not be empty")

	dedup.ClearCache()
	assert.Equal(t, 0, dedup.GetCacheSize(), "Cache should be empty after clear")
}

func BenchmarkDeduplicate(b *testing.B) {
	dedup := deduplicator.NewDeduplicator(nil)

	// Generate 1000 test articles
	articles := make([]types.ParsedArticle, 1000)
	for i := 0; i < 1000; i++ {
		articles[i] = types.ParsedArticle{
			Title:       "Article " + string(rune(i)),
			URL:         "https://example.com/" + string(rune(i)),
			ExternalID:  string(rune(i)),
			FetchSource: "rss",
			PublishedAt: time.Now(),
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dedup.Deduplicate(articles)
	}
}
