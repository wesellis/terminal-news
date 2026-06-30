package classifier

import (
	"strings"

	"github.com/wesellis/terminal-news/scraper/pkg/types"
)

type Classifier struct {
	categories map[string][]string
}

func NewClassifier() *Classifier {
	return &Classifier{
		categories: initCategoryKeywords(),
	}
}

func initCategoryKeywords() map[string][]string {
	return map[string][]string{
		"tech": {
			"software", "hardware", "computer", "internet", "app", "startup",
			"ai", "artificial intelligence", "machine learning", "blockchain",
			"crypto", "bitcoin", "programming", "code", "developer", "tech",
			"silicon valley", "google", "apple", "microsoft", "facebook",
			"amazon", "tesla", "spacex", "data", "cloud", "cybersecurity",
			"algorithm", "api", "database", "framework", "git", "javascript",
			"python", "rust", "golang", "react", "kubernetes", "docker",
		},
		"politics": {
			"election", "president", "congress", "senate", "democrat",
			"republican", "government", "policy", "law", "legislation",
			"vote", "campaign", "politician", "parliament", "minister",
			"governor", "mayor", "political", "democracy", "constitution",
		},
		"business": {
			"market", "stock", "economy", "finance", "bank", "investment",
			"trading", "earnings", "revenue", "profit", "ipo", "merger",
			"acquisition", "startup", "entrepreneur", "ceo", "company",
			"corporation", "business", "industry", "commerce", "retail",
			"gdp", "inflation", "recession", "growth", "quarter",
		},
		"science": {
			"research", "study", "scientist", "experiment", "discovery",
			"physics", "chemistry", "biology", "medicine", "health",
			"disease", "vaccine", "drug", "treatment", "climate", "space",
			"nasa", "astronomy", "planet", "universe", "quantum", "dna",
			"gene", "evolution", "ecology", "environment",
		},
		"sports": {
			"game", "match", "tournament", "championship", "league",
			"team", "player", "coach", "score", "win", "loss", "football",
			"basketball", "baseball", "soccer", "tennis", "golf", "olympics",
			"athlete", "sport", "nfl", "nba", "mlb", "fifa", "uefa",
		},
		"entertainment": {
			"movie", "film", "actor", "actress", "director", "hollywood",
			"netflix", "streaming", "tv", "television", "show", "series",
			"music", "album", "song", "artist", "concert", "tour", "award",
			"oscar", "grammy", "celebrity", "entertainment",
		},
	}
}

func (c *Classifier) Classify(article *types.ParsedArticle) {
	// If category already set and not "general", keep it
	if article.Category != "" && article.Category != "general" {
		return
	}

	// Combine title and content for analysis
	text := strings.ToLower(article.Title + " " + article.Summary)

	// Count keyword matches for each category
	scores := make(map[string]int)

	for category, keywords := range c.categories {
		for _, keyword := range keywords {
			if strings.Contains(text, keyword) {
				scores[category]++
			}
		}
	}

	// Find category with highest score
	maxScore := 0
	bestCategory := "general"

	for category, score := range scores {
		if score > maxScore {
			maxScore = score
			bestCategory = category
		}
	}

	article.Category = bestCategory
}

func (c *Classifier) ExtractKeywords(text string, count int) []string {
	// Simple keyword extraction - count word frequency
	words := strings.Fields(strings.ToLower(text))
	freq := make(map[string]int)

	// Count word frequency
	for _, word := range words {
		// Skip short words and common words
		if len(word) < 4 || c.isStopWord(word) {
			continue
		}
		freq[word]++
	}

	// Get top N words
	keywords := make([]string, 0, count)
	for i := 0; i < count && len(freq) > 0; i++ {
		maxWord := ""
		maxCount := 0
		for word, cnt := range freq {
			if cnt > maxCount {
				maxCount = cnt
				maxWord = word
			}
		}
		if maxWord != "" {
			keywords = append(keywords, maxWord)
			delete(freq, maxWord)
		}
	}

	return keywords
}

func (c *Classifier) isStopWord(word string) bool {
	stopWords := map[string]bool{
		"the": true, "and": true, "for": true, "that": true, "this": true,
		"with": true, "from": true, "have": true, "they": true, "will": true,
		"what": true, "when": true, "where": true, "which": true, "said": true,
	}
	return stopWords[word]
}
