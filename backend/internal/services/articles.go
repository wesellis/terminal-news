package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

var (
	ErrArticleNotFound = errors.New("article not found")
)

type ArticleService struct {
	db  *sqlx.DB
	rdb *redis.Client
}

func NewArticleService(db *sqlx.DB, rdb *redis.Client) *ArticleService {
	return &ArticleService{
		db:  db,
		rdb: rdb,
	}
}

// Article represents an article with its ranking data
type Article struct {
	ID          int64     `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	URL         string    `json:"url" db:"url"`
	Content     *string   `json:"content,omitempty" db:"content"`
	ImageURL    *string   `json:"image_url,omitempty" db:"image_url"`
	Source      *string   `json:"source,omitempty" db:"source"`
	Author      *string   `json:"author,omitempty" db:"author"`
	PublishedAt time.Time `json:"published_at" db:"published_at"`
	Category    *string   `json:"category,omitempty" db:"category"`
	Tags        []string  `json:"tags,omitempty" db:"tags"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`

	// Ranking data (from join with article_rankings)
	OpenCount        int     `json:"open_count" db:"open_count"`
	LikeCount        int     `json:"like_count" db:"like_count"`
	DislikeCount     int     `json:"dislike_count" db:"dislike_count"`
	TotalScore       int     `json:"total_score" db:"total_score"`
	ControversyScore float64 `json:"controversy_score" db:"controversy_score"`
	TotalEngagement  int     `json:"total_engagement" db:"total_engagement"`
	HotRank          float64 `json:"hot_rank" db:"hot_rank"`
}

// ArticleListResponse is the response for article list endpoints
type ArticleListResponse struct {
	Articles   []Article `json:"articles"`
	TotalCount int       `json:"total_count"`
	Limit      int       `json:"limit"`
	Offset     int       `json:"offset"`
}

// GetHotArticles retrieves articles ordered by hot rank
func (s *ArticleService) GetHotArticles(ctx context.Context, limit, offset int) (*ArticleListResponse, error) {
	cacheKey := fmt.Sprintf("articles:hot:%d:%d", limit, offset)

	// Try to get from cache first
	cached, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var response ArticleListResponse
		if err := json.Unmarshal([]byte(cached), &response); err == nil {
			return &response, nil
		}
	}

	// Query from database
	query := `
		SELECT
			a.id, a.title, a.url, a.content, a.image_url, a.source, a.author,
			a.published_at, a.category, a.tags, a.created_at,
			COALESCE(r.open_count, 0) AS open_count,
			COALESCE(r.like_count, 0) AS like_count,
			COALESCE(r.dislike_count, 0) AS dislike_count,
			COALESCE(r.total_score, 0) AS total_score,
			COALESCE(r.controversy_score, 0) AS controversy_score,
			COALESCE(r.total_engagement, 0) AS total_engagement,
			COALESCE(r.hot_rank, 0) AS hot_rank
		FROM articles a
		LEFT JOIN article_rankings r ON a.id = r.article_id
		WHERE a.published_at > NOW() - INTERVAL '7 days'
		ORDER BY r.hot_rank DESC NULLS LAST, a.published_at DESC
		LIMIT $1 OFFSET $2
	`

	var articles []Article
	if err := s.db.SelectContext(ctx, &articles, query, limit, offset); err != nil {
		return nil, fmt.Errorf("failed to get hot articles: %w", err)
	}

	// Get total count
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM articles WHERE published_at > NOW() - INTERVAL '7 days'`
	if err := s.db.GetContext(ctx, &totalCount, countQuery); err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	response := &ArticleListResponse{
		Articles:   articles,
		TotalCount: totalCount,
		Limit:      limit,
		Offset:     offset,
	}

	// Cache the response for 5 minutes
	if data, err := json.Marshal(response); err == nil {
		s.rdb.Set(ctx, cacheKey, data, 5*time.Minute)
	}

	return response, nil
}

// GetControversialArticles retrieves articles ordered by controversy score
func (s *ArticleService) GetControversialArticles(ctx context.Context, limit, offset int) (*ArticleListResponse, error) {
	cacheKey := fmt.Sprintf("articles:controversial:%d:%d", limit, offset)

	// Try to get from cache first
	cached, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var response ArticleListResponse
		if err := json.Unmarshal([]byte(cached), &response); err == nil {
			return &response, nil
		}
	}

	// Query from database
	query := `
		SELECT
			a.id, a.title, a.url, a.content, a.image_url, a.source, a.author,
			a.published_at, a.category, a.tags, a.created_at,
			COALESCE(r.open_count, 0) AS open_count,
			COALESCE(r.like_count, 0) AS like_count,
			COALESCE(r.dislike_count, 0) AS dislike_count,
			COALESCE(r.total_score, 0) AS total_score,
			COALESCE(r.controversy_score, 0) AS controversy_score,
			COALESCE(r.total_engagement, 0) AS total_engagement,
			COALESCE(r.hot_rank, 0) AS hot_rank
		FROM articles a
		LEFT JOIN article_rankings r ON a.id = r.article_id
		WHERE a.published_at > NOW() - INTERVAL '7 days'
			AND r.controversy_score > 0
		ORDER BY r.controversy_score DESC, r.total_engagement DESC
		LIMIT $1 OFFSET $2
	`

	var articles []Article
	if err := s.db.SelectContext(ctx, &articles, query, limit, offset); err != nil {
		return nil, fmt.Errorf("failed to get controversial articles: %w", err)
	}

	// Get total count
	var totalCount int
	countQuery := `
		SELECT COUNT(*)
		FROM articles a
		LEFT JOIN article_rankings r ON a.id = r.article_id
		WHERE a.published_at > NOW() - INTERVAL '7 days'
			AND r.controversy_score > 0
	`
	if err := s.db.GetContext(ctx, &totalCount, countQuery); err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	response := &ArticleListResponse{
		Articles:   articles,
		TotalCount: totalCount,
		Limit:      limit,
		Offset:     offset,
	}

	// Cache the response for 5 minutes
	if data, err := json.Marshal(response); err == nil {
		s.rdb.Set(ctx, cacheKey, data, 5*time.Minute)
	}

	return response, nil
}

// GetRisingArticles retrieves articles with rapid recent engagement
func (s *ArticleService) GetRisingArticles(ctx context.Context, limit, offset int) (*ArticleListResponse, error) {
	cacheKey := fmt.Sprintf("articles:rising:%d:%d", limit, offset)

	// Try to get from cache first
	cached, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var response ArticleListResponse
		if err := json.Unmarshal([]byte(cached), &response); err == nil {
			return &response, nil
		}
	}

	// Rising: Articles published in last 48 hours with increasing engagement
	query := `
		SELECT
			a.id, a.title, a.url, a.content, a.image_url, a.source, a.author,
			a.published_at, a.category, a.tags, a.created_at,
			COALESCE(r.open_count, 0) AS open_count,
			COALESCE(r.like_count, 0) AS like_count,
			COALESCE(r.dislike_count, 0) AS dislike_count,
			COALESCE(r.total_score, 0) AS total_score,
			COALESCE(r.controversy_score, 0) AS controversy_score,
			COALESCE(r.total_engagement, 0) AS total_engagement,
			COALESCE(r.hot_rank, 0) AS hot_rank,
			-- Rising score: engagement / hours since published
			CASE
				WHEN EXTRACT(EPOCH FROM (NOW() - a.published_at)) / 3600 > 0
				THEN COALESCE(r.total_engagement, 0)::FLOAT / (EXTRACT(EPOCH FROM (NOW() - a.published_at)) / 3600)
				ELSE 0
			END AS rising_score
		FROM articles a
		LEFT JOIN article_rankings r ON a.id = r.article_id
		WHERE a.published_at > NOW() - INTERVAL '48 hours'
			AND r.total_engagement >= 5
		ORDER BY rising_score DESC, r.total_engagement DESC
		LIMIT $1 OFFSET $2
	`

	var articles []Article
	if err := s.db.SelectContext(ctx, &articles, query, limit, offset); err != nil {
		return nil, fmt.Errorf("failed to get rising articles: %w", err)
	}

	// Get total count
	var totalCount int
	countQuery := `
		SELECT COUNT(*)
		FROM articles a
		LEFT JOIN article_rankings r ON a.id = r.article_id
		WHERE a.published_at > NOW() - INTERVAL '48 hours'
			AND r.total_engagement >= 5
	`
	if err := s.db.GetContext(ctx, &totalCount, countQuery); err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	response := &ArticleListResponse{
		Articles:   articles,
		TotalCount: totalCount,
		Limit:      limit,
		Offset:     offset,
	}

	// Cache the response for 3 minutes (shorter for rising)
	if data, err := json.Marshal(response); err == nil {
		s.rdb.Set(ctx, cacheKey, data, 3*time.Minute)
	}

	return response, nil
}

// GetArticles retrieves articles with optional feed type
func (s *ArticleService) GetArticles(ctx context.Context, feed string, limit, offset int) (*ArticleListResponse, error) {
	switch feed {
	case "hot":
		return s.GetHotArticles(ctx, limit, offset)
	case "controversial":
		return s.GetControversialArticles(ctx, limit, offset)
	case "rising":
		return s.GetRisingArticles(ctx, limit, offset)
	default:
		// Default to hot feed
		return s.GetHotArticles(ctx, limit, offset)
	}
}

// GetArticle retrieves a single article by ID
func (s *ArticleService) GetArticle(ctx context.Context, id int64) (*Article, error) {
	cacheKey := fmt.Sprintf("article:%d", id)

	// Try to get from cache first
	cached, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var article Article
		if err := json.Unmarshal([]byte(cached), &article); err == nil {
			return &article, nil
		}
	}

	// Query from database
	query := `
		SELECT
			a.id, a.title, a.url, a.content, a.image_url, a.source, a.author,
			a.published_at, a.category, a.tags, a.created_at,
			COALESCE(r.open_count, 0) AS open_count,
			COALESCE(r.like_count, 0) AS like_count,
			COALESCE(r.dislike_count, 0) AS dislike_count,
			COALESCE(r.total_score, 0) AS total_score,
			COALESCE(r.controversy_score, 0) AS controversy_score,
			COALESCE(r.total_engagement, 0) AS total_engagement,
			COALESCE(r.hot_rank, 0) AS hot_rank
		FROM articles a
		LEFT JOIN article_rankings r ON a.id = r.article_id
		WHERE a.id = $1
	`

	var article Article
	if err := s.db.GetContext(ctx, &article, query, id); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, ErrArticleNotFound
		}
		return nil, fmt.Errorf("failed to get article: %w", err)
	}

	// Cache the article for 10 minutes
	if data, err := json.Marshal(article); err == nil {
		s.rdb.Set(ctx, cacheKey, data, 10*time.Minute)
	}

	return &article, nil
}
