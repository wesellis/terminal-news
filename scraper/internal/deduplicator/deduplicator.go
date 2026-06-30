package deduplicator

import (
	"crypto/md5"
	"fmt"
	"log"
	"strings"

	"github.com/wesellis/terminal-news/scraper/internal/storage"
	"github.com/wesellis/terminal-news/scraper/pkg/types"
)

type Deduplicator struct {
	storage *storage.Storage
	cache   map[string]bool // In-memory cache for this session
}

func NewDeduplicator(storage *storage.Storage) *Deduplicator {
	return &Deduplicator{
		storage: storage,
		cache:   make(map[string]bool),
	}
}

// Deduplicate removes duplicate articles from the slice
func (d *Deduplicator) Deduplicate(articles []types.ParsedArticle) []types.ParsedArticle {
	unique := make([]types.ParsedArticle, 0, len(articles))
	seen := make(map[string]bool)

	for _, article := range articles {
		// Generate multiple hashes for different duplicate detection methods
		urlHash := d.hashURL(article.URL)
		titleHash := d.hashTitle(article.Title)
		contentHash := d.hashContent(article.Title, article.Summary)

		// Check if we've seen this URL
		if seen[urlHash] {
			log.Printf("Duplicate URL detected (in batch): %s", article.Title)
			continue
		}

		// Check if we've seen similar title
		if seen[titleHash] {
			log.Printf("Duplicate title detected (in batch): %s", article.Title)
			continue
		}

		// Check if we've seen similar content
		if seen[contentHash] {
			log.Printf("Duplicate content detected (in batch): %s", article.Title)
			continue
		}

		// Check against database if storage available
		if d.storage != nil {
			exists, err := d.storage.GetArticleByExternalID(article.ExternalID, article.FetchSource)
			if err == nil && exists {
				log.Printf("Duplicate article in database: %s", article.Title)
				continue
			}
		}

		// Check in-memory cache (for current session)
		cacheKey := fmt.Sprintf("%s:%s", article.FetchSource, article.ExternalID)
		if d.cache[cacheKey] {
			log.Printf("Duplicate article in cache: %s", article.Title)
			continue
		}

		// Mark as seen
		seen[urlHash] = true
		seen[titleHash] = true
		seen[contentHash] = true
		d.cache[cacheKey] = true

		unique = append(unique, article)
	}

	removed := len(articles) - len(unique)
	if removed > 0 {
		log.Printf("Removed %d duplicates from batch of %d articles", removed, len(articles))
	}

	return unique
}

// IsDuplicate checks if a single article is a duplicate
func (d *Deduplicator) IsDuplicate(article types.ParsedArticle) bool {
	// Check in-memory cache first (fastest)
	cacheKey := fmt.Sprintf("%s:%s", article.FetchSource, article.ExternalID)
	if d.cache[cacheKey] {
		return true
	}

	// Check database
	if d.storage != nil {
		exists, err := d.storage.GetArticleByExternalID(article.ExternalID, article.FetchSource)
		if err == nil && exists {
			return true
		}
	}

	return false
}

// FindSimilarArticles finds articles with similar titles
func (d *Deduplicator) FindSimilarArticles(article types.ParsedArticle) ([]types.ParsedArticle, error) {
	if d.storage == nil {
		return nil, nil
	}

	// Use fuzzy matching with 70% similarity threshold
	return d.storage.FindDuplicatesByTitle(article.Title, 0.7)
}

// CalculateSimilarity calculates text similarity between two strings (0.0 to 1.0)
func (d *Deduplicator) CalculateSimilarity(text1, text2 string) float64 {
	// Normalize texts
	t1 := d.normalizeText(text1)
	t2 := d.normalizeText(text2)

	// Use Jaccard similarity (intersection over union of words)
	words1 := d.extractWords(t1)
	words2 := d.extractWords(t2)

	if len(words1) == 0 || len(words2) == 0 {
		return 0.0
	}

	// Calculate intersection
	intersection := 0
	seen := make(map[string]bool)
	for _, word := range words1 {
		seen[word] = true
	}
	for _, word := range words2 {
		if seen[word] {
			intersection++
		}
	}

	// Calculate union
	union := len(words1) + len(words2) - intersection

	if union == 0 {
		return 0.0
	}

	return float64(intersection) / float64(union)
}

func (d *Deduplicator) hashURL(url string) string {
	// Normalize URL (remove query params, trailing slash, etc.)
	url = strings.ToLower(url)
	url = strings.Split(url, "?")[0] // Remove query string
	url = strings.TrimSuffix(url, "/")

	hash := md5.Sum([]byte(url))
	return fmt.Sprintf("url:%x", hash)
}

func (d *Deduplicator) hashTitle(title string) string {
	// Normalize title (lowercase, remove punctuation, extra spaces)
	normalized := d.normalizeText(title)

	// Take first N significant words for comparison
	words := d.extractWords(normalized)
	if len(words) > 8 {
		words = words[:8]
	}

	key := strings.Join(words, " ")
	hash := md5.Sum([]byte(key))
	return fmt.Sprintf("title:%x", hash)
}

func (d *Deduplicator) hashContent(title, summary string) string {
	// Combine title and summary for content hash
	content := title + " " + summary
	normalized := d.normalizeText(content)

	// Take first 100 characters for quick comparison
	if len(normalized) > 100 {
		normalized = normalized[:100]
	}

	hash := md5.Sum([]byte(normalized))
	return fmt.Sprintf("content:%x", hash)
}

func (d *Deduplicator) normalizeText(text string) string {
	// Convert to lowercase
	text = strings.ToLower(text)

	// Remove common punctuation
	replacements := []string{
		".", "", ",", "", "!", "", "?", "",
		":", "", ";", "", "'", "", "\"", "",
		"(", "", ")", "", "[", "", "]", "",
		"{", "", "}", "", "-", " ", "_", " ",
	}
	replacer := strings.NewReplacer(replacements...)
	text = replacer.Replace(text)

	// Normalize whitespace
	text = strings.Join(strings.Fields(text), " ")

	return strings.TrimSpace(text)
}

func (d *Deduplicator) extractWords(text string) []string {
	// Split into words and filter out stop words
	words := strings.Fields(text)
	stopWords := d.getStopWords()

	filtered := make([]string, 0, len(words))
	for _, word := range words {
		if len(word) < 3 { // Skip very short words
			continue
		}
		if stopWords[word] { // Skip stop words
			continue
		}
		filtered = append(filtered, word)
	}

	return filtered
}

func (d *Deduplicator) getStopWords() map[string]bool {
	// Common English stop words
	return map[string]bool{
		"the": true, "be": true, "to": true, "of": true, "and": true,
		"a": true, "in": true, "that": true, "have": true, "i": true,
		"it": true, "for": true, "not": true, "on": true, "with": true,
		"he": true, "as": true, "you": true, "do": true, "at": true,
		"this": true, "but": true, "his": true, "by": true, "from": true,
		"they": true, "we": true, "say": true, "her": true, "she": true,
		"or": true, "an": true, "will": true, "my": true, "one": true,
		"all": true, "would": true, "there": true, "their": true, "what": true,
		"so": true, "up": true, "out": true, "if": true, "about": true,
		"who": true, "get": true, "which": true, "go": true, "me": true,
		"when": true, "make": true, "can": true, "like": true, "time": true,
		"no": true, "just": true, "him": true, "know": true, "take": true,
		"people": true, "into": true, "year": true, "your": true, "good": true,
		"some": true, "could": true, "them": true, "see": true, "other": true,
		"than": true, "then": true, "now": true, "look": true, "only": true,
		"come": true, "its": true, "over": true, "think": true, "also": true,
		"back": true, "after": true, "use": true, "two": true, "how": true,
		"our": true, "work": true, "first": true, "well": true, "way": true,
		"even": true, "new": true, "want": true, "because": true, "any": true,
		"these": true, "give": true, "day": true, "most": true, "us": true,
	}
}

// ClearCache clears the in-memory cache
func (d *Deduplicator) ClearCache() {
	d.cache = make(map[string]bool)
	log.Println("Deduplication cache cleared")
}

// GetCacheSize returns the size of the current cache
func (d *Deduplicator) GetCacheSize() int {
	return len(d.cache)
}
