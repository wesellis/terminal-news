package storage

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/wesellis/terminal-news/scraper/pkg/types"
)

type Storage struct {
	db *sqlx.DB
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{db: db}
}

// StoreArticles stores multiple articles, skipping duplicates
func (s *Storage) StoreArticles(articles []types.ParsedArticle) (int, error) {
	stored := 0

	for _, article := range articles {
		err := s.StoreArticle(article)
		if err != nil {
			// Check if it's a duplicate error
			if pqErr, ok := err.(*pq.Error); ok {
				if pqErr.Code == "23505" { // unique_violation
					log.Printf("Duplicate article skipped: %s", article.Title)
					continue
				}
			}
			log.Printf("Failed to store article %s: %v", article.Title, err)
			continue
		}
		stored++
	}

	log.Printf("Stored %d/%d articles", stored, len(articles))
	return stored, nil
}

// StoreArticle stores a single article
func (s *Storage) StoreArticle(article types.ParsedArticle) error {
	query := `
		INSERT INTO articles (
			title, url, content, image_url, source, author,
			published_at, category, tags, external_id, fetch_source,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW()
		)
		ON CONFLICT (external_id, fetch_source) DO NOTHING
		RETURNING id
	`

	var id int64
	err := s.db.Get(&id, query,
		article.Title,
		article.URL,
		article.Summary, // Store summary as content for now
		article.ImageURL,
		article.Source,
		article.Author,
		article.PublishedAt,
		article.Category,
		pq.Array(article.Tags),
		article.ExternalID,
		article.FetchSource,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			// Article was a duplicate (ON CONFLICT DO NOTHING)
			return nil
		}
		return fmt.Errorf("failed to insert article: %w", err)
	}

	return nil
}

// GetArticleByExternalID checks if an article already exists
func (s *Storage) GetArticleByExternalID(externalID, fetchSource string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM articles
			WHERE external_id = $1 AND fetch_source = $2
		)
	`

	var exists bool
	err := s.db.Get(&exists, query, externalID, fetchSource)
	return exists, err
}

// GetRecentArticlesBySource gets recent articles from a specific source
func (s *Storage) GetRecentArticlesBySource(source string, limit int) ([]types.ParsedArticle, error) {
	query := `
		SELECT
			id, title, url, content, image_url, source, author,
			published_at, category, tags, external_id, fetch_source,
			created_at
		FROM articles
		WHERE source = $1
		ORDER BY published_at DESC
		LIMIT $2
	`

	var articles []struct {
		ID          int64          `db:"id"`
		Title       string         `db:"title"`
		URL         string         `db:"url"`
		Content     string         `db:"content"`
		ImageURL    sql.NullString `db:"image_url"`
		Source      string         `db:"source"`
		Author      sql.NullString `db:"author"`
		PublishedAt sql.NullTime   `db:"published_at"`
		Category    string         `db:"category"`
		Tags        pq.StringArray `db:"tags"`
		ExternalID  string         `db:"external_id"`
		FetchSource string         `db:"fetch_source"`
		CreatedAt   sql.NullTime   `db:"created_at"`
	}

	err := s.db.Select(&articles, query, source, limit)
	if err != nil {
		return nil, err
	}

	result := make([]types.ParsedArticle, len(articles))
	for i, a := range articles {
		result[i] = types.ParsedArticle{
			Title:       a.Title,
			URL:         a.URL,
			Summary:     a.Content,
			ImageURL:    a.ImageURL.String,
			Source:      a.Source,
			Author:      a.Author.String,
			Category:    a.Category,
			Tags:        a.Tags,
			ExternalID:  a.ExternalID,
			FetchSource: a.FetchSource,
		}
		if a.PublishedAt.Valid {
			result[i].PublishedAt = a.PublishedAt.Time
		}
		if a.CreatedAt.Valid {
			result[i].FetchedAt = a.CreatedAt.Time
		}
	}

	return result, nil
}

// GetArticleCount returns total article count
func (s *Storage) GetArticleCount() (int, error) {
	var count int
	err := s.db.Get(&count, "SELECT COUNT(*) FROM articles")
	return count, err
}

// GetArticleCountBySource returns article count per source
func (s *Storage) GetArticleCountBySource() (map[string]int, error) {
	query := `
		SELECT source, COUNT(*) as count
		FROM articles
		GROUP BY source
		ORDER BY count DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var source string
		var count int
		if err := rows.Scan(&source, &count); err != nil {
			continue
		}
		counts[source] = count
	}

	return counts, nil
}

// CleanOldArticles removes articles older than specified days
func (s *Storage) CleanOldArticles(days int) (int, error) {
	query := `
		DELETE FROM articles
		WHERE published_at < NOW() - INTERVAL '$1 days'
	`

	result, err := s.db.Exec(query, days)
	if err != nil {
		return 0, err
	}

	rows, _ := result.RowsAffected()
	log.Printf("Cleaned %d old articles (older than %d days)", rows, days)
	return int(rows), nil
}

// FindDuplicatesByTitle finds articles with similar titles
func (s *Storage) FindDuplicatesByTitle(title string, threshold float64) ([]types.ParsedArticle, error) {
	// Use PostgreSQL similarity function for fuzzy matching
	query := `
		SELECT id, title, url, source, similarity(title, $1) as sim
		FROM articles
		WHERE similarity(title, $1) > $2
		ORDER BY sim DESC
		LIMIT 10
	`

	var results []struct {
		ID         int64   `db:"id"`
		Title      string  `db:"title"`
		URL        string  `db:"url"`
		Source     string  `db:"source"`
		Similarity float64 `db:"sim"`
	}

	err := s.db.Select(&results, query, title, threshold)
	if err != nil {
		// Fallback to simple matching if similarity extension not available
		return s.findDuplicatesByTitleSimple(title)
	}

	articles := make([]types.ParsedArticle, len(results))
	for i, r := range results {
		articles[i] = types.ParsedArticle{
			Title:  r.Title,
			URL:    r.URL,
			Source: r.Source,
		}
	}

	return articles, nil
}

func (s *Storage) findDuplicatesByTitleSimple(title string) ([]types.ParsedArticle, error) {
	// Simple LIKE-based matching as fallback
	query := `
		SELECT id, title, url, source
		FROM articles
		WHERE LOWER(title) LIKE LOWER($1)
		LIMIT 10
	`

	pattern := "%" + strings.Join(strings.Fields(title)[:min(len(strings.Fields(title)), 5)], "%") + "%"

	var results []struct {
		ID     int64  `db:"id"`
		Title  string `db:"title"`
		URL    string `db:"url"`
		Source string `db:"source"`
	}

	err := s.db.Select(&results, query, pattern)
	if err != nil {
		return nil, err
	}

	articles := make([]types.ParsedArticle, len(results))
	for i, r := range results {
		articles[i] = types.ParsedArticle{
			Title:  r.Title,
			URL:    r.URL,
			Source: r.Source,
		}
	}

	return articles, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
