package moderation

import (
	"log"
	"regexp"
	"strings"
	"unicode"

	"github.com/jmoiron/sqlx"
	"github.com/wesellis/terminal-news/scraper/pkg/types"
)

type Moderator struct {
	db                *sqlx.DB
	spamPatterns      []*regexp.Regexp
	spamKeywords      []string
	trustedSources    map[string]bool
	suspiciousDomains []string
}

func NewModerator(db *sqlx.DB) *Moderator {
	return &Moderator{
		db:                db,
		spamPatterns:      initSpamPatterns(),
		spamKeywords:      initSpamKeywords(),
		trustedSources:    initTrustedSources(),
		suspiciousDomains: initSuspiciousDomains(),
	}
}

func initSpamPatterns() []*regexp.Regexp {
	patterns := []string{
		`(?i)\b(viagra|cialis|levitra)\b`,
		`(?i)\b(casino|gambling|poker|slots)\b`,
		`(?i)\b(weight.?loss|diet.?pills|lose.?\d+.?pounds)\b`,
		`(?i)\b(click.?here|buy.?now|limited.?time|act.?now)\b`,
		`(?i)\b(work.?from.?home|make.?\$\d+|earn.?money.?online)\b`,
		`(?i)\b(hot.?singles|adult.?dating|xxx)\b`,
		`\b[A-Z]{15,}\b`,          // All caps spam
		`(.)\1{10,}`,              // Repeated characters
		`https?://bit\.ly|tinyurl\.com|goo\.gl`, // URL shorteners
	}

	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, pattern := range patterns {
		if re, err := regexp.Compile(pattern); err == nil {
			compiled = append(compiled, re)
		}
	}

	return compiled
}

func initSpamKeywords() []string {
	return []string{
		"free money", "get rich quick", "mlm", "pyramid scheme",
		"miracle cure", "doctors hate", "one weird trick",
		"nigerian prince", "wire transfer", "bitcoin doubler",
		"adult content", "nsfw", "18+",
	}
}

func initTrustedSources() map[string]bool {
	return map[string]bool{
		"BBC":             true,
		"CNN":             true,
		"Reuters":         true,
		"AP News":         true,
		"The Guardian":    true,
		"NY Times":        true,
		"Washington Post": true,
		"WSJ":             true,
		"Bloomberg":       true,
		"TechCrunch":      true,
		"The Verge":       true,
		"Ars Technica":    true,
	}
}

func initSuspiciousDomains() []string {
	return []string{
		"bit.ly", "tinyurl.com", "goo.gl", "ow.ly",
		"blogspot.com", "wordpress.com", "tumblr.com",
		".tk", ".ml", ".ga", ".cf", // Free TLDs often used for spam
	}
}

func (m *Moderator) IsSpam(article types.ParsedArticle) bool {
	score := 0.0

	// Check if from trusted source (-50 points, likely not spam)
	if m.trustedSources[article.Source] {
		score -= 50
	}

	// Check spam patterns (+20 points each)
	text := strings.ToLower(article.Title + " " + article.Summary)
	for _, pattern := range m.spamPatterns {
		if pattern.MatchString(text) {
			score += 20
			log.Printf("Spam pattern matched: %v in %s", pattern, article.Title)
		}
	}

	// Check spam keywords (+10 points each)
	for _, keyword := range m.spamKeywords {
		if strings.Contains(text, keyword) {
			score += 10
		}
	}

	// Check suspicious domains (+15 points)
	for _, domain := range m.suspiciousDomains {
		if strings.Contains(article.URL, domain) {
			score += 15
		}
	}

	// Check for excessive caps (+10 points)
	if m.hasExcessiveCaps(article.Title) {
		score += 10
	}

	// Check for excessive punctuation (+10 points)
	if m.hasExcessivePunctuation(article.Title) {
		score += 10
	}

	// Check content quality
	if len(article.Summary) < 50 {
		score += 5 // Too short, might be spam
	}

	if len(article.Title) < 10 {
		score += 5 // Title too short
	}

	// If score > 30, likely spam
	isSpam := score > 30

	if isSpam {
		log.Printf("Article flagged as spam: %s (score: %.2f)", article.Title, score)
	}

	return isSpam
}

func (m *Moderator) hasExcessiveCaps(text string) bool {
	caps := 0
	total := 0

	for _, r := range text {
		if unicode.IsLetter(r) {
			total++
			if unicode.IsUpper(r) {
				caps++
			}
		}
	}

	if total == 0 {
		return false
	}

	return float64(caps)/float64(total) > 0.5
}

func (m *Moderator) hasExcessivePunctuation(text string) bool {
	punctuation := 0

	for _, r := range text {
		if r == '!' || r == '?' || r == '$' {
			punctuation++
		}
	}

	return punctuation > 3
}

func (m *Moderator) CalculateTrustScore(article types.ParsedArticle) float64 {
	score := 0.5 // Start at neutral

	// Trusted source boost
	if m.trustedSources[article.Source] {
		score += 0.3
	}

	// Has author
	if article.Author != "" {
		score += 0.1
	}

	// Reasonable length
	if len(article.Summary) >= 100 && len(article.Summary) <= 1000 {
		score += 0.1
	}

	// No spam indicators
	if !m.IsSpam(article) {
		score += 0.1
	}

	// Cap at 1.0
	if score > 1.0 {
		score = 1.0
	}

	return score
}

func (m *Moderator) QueueForReview(articles []types.ParsedArticle) error {
	query := `
		INSERT INTO moderation_queue (
			item_type, item_id, user_id,
			flag_reason, auto_flagged, ai_confidence
		) VALUES ($1, $2, $3, $4, $5, $6)
	`

	for _, article := range articles {
		// Store article first
		var articleID int64
		err := m.db.Get(&articleID,
			"INSERT INTO articles (title, url, source) VALUES ($1, $2, $3) RETURNING id",
			article.Title, article.URL, article.Source)

		if err != nil {
			continue
		}

		// Add to moderation queue
		_, err = m.db.Exec(query,
			"article",
			articleID,
			1, // System user
			"spam",
			true,
			m.CalculateTrustScore(article),
		)

		if err != nil {
			log.Printf("Failed to queue for review: %v", err)
		}
	}

	return nil
}
