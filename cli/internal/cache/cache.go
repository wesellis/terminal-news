package cache

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/wesellis/terminal-news/cli/internal/models"
)

// Cache manages local SQLite database for offline support
type Cache struct {
	db *sqlx.DB
}

// New creates a new cache instance
func New(dbPath string) (*Cache, error) {
	db, err := sqlx.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	cache := &Cache{db: db}

	// Create tables
	if err := cache.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return cache, nil
}

// createTables creates all necessary cache tables
func (c *Cache) createTables() error {
	schema := `
	CREATE TABLE IF NOT EXISTS articles (
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL,
		url TEXT NOT NULL,
		source TEXT,
		summary TEXT,
		published_at DATETIME,
		upvotes INTEGER DEFAULT 0,
		downvotes INTEGER DEFAULT 0,
		views INTEGER DEFAULT 0,
		comment_count INTEGER DEFAULT 0,
		is_hot BOOLEAN DEFAULT 0,
		is_rising BOOLEAN DEFAULT 0,
		cached_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_articles_cached ON articles(cached_at DESC);
	CREATE INDEX IF NOT EXISTS idx_articles_hot ON articles(is_hot) WHERE is_hot = 1;
	CREATE INDEX IF NOT EXISTS idx_articles_rising ON articles(is_rising) WHERE is_rising = 1;

	CREATE TABLE IF NOT EXISTS read_articles (
		article_id INTEGER PRIMARY KEY,
		read_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY,
		article_id INTEGER NOT NULL,
		user_id INTEGER,
		username TEXT,
		content TEXT NOT NULL,
		parent_id INTEGER,
		upvotes INTEGER DEFAULT 0,
		downvotes INTEGER DEFAULT 0,
		created_at DATETIME,
		cached_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (article_id) REFERENCES articles(id)
	);

	CREATE INDEX IF NOT EXISTS idx_comments_article ON comments(article_id);

	CREATE TABLE IF NOT EXISTS classifieds (
		id INTEGER PRIMARY KEY,
		user_id INTEGER,
		username TEXT,
		title TEXT NOT NULL,
		description TEXT,
		price REAL,
		category TEXT,
		city TEXT,
		state TEXT,
		is_premium BOOLEAN DEFAULT 0,
		is_active BOOLEAN DEFAULT 1,
		created_at DATETIME,
		cached_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_classifieds_category ON classifieds(category);
	CREATE INDEX IF NOT EXISTS idx_classifieds_location ON classifieds(city, state);

	CREATE TABLE IF NOT EXISTS offline_queue (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		action TEXT NOT NULL,
		payload TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS user_settings (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS weather_cache (
		location TEXT PRIMARY KEY,
		data TEXT NOT NULL,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := c.db.Exec(schema)
	return err
}

// Article operations

func (c *Cache) SaveArticles(articles []models.Article) error {
	tx, err := c.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT OR REPLACE INTO articles
		(id, title, url, source, summary, published_at, upvotes, downvotes, views, comment_count, is_hot, is_rising)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	for _, article := range articles {
		// Cache only core Article fields (not ranking fields)
		_, err := tx.Exec(query,
			article.ID,
			article.Title,
			article.URL,
			article.Source,
			article.Content, // Changed from Summary
			article.PublishedAt,
			0, // Upvotes - not in shared model
			0, // Downvotes - not in shared model
			0, // Views - not in shared model
			0, // CommentCount - not in shared model
			false, // IsHot - not in shared model
			false, // IsRising - not in shared model
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (c *Cache) GetArticles(feed string, limit int) ([]models.Article, error) {
	var articles []models.Article

	query := `
		SELECT id, title, url, source, summary, published_at,
		       upvotes, downvotes, views, comment_count, is_hot, is_rising, cached_at
		FROM articles
	`

	switch feed {
	case "hot":
		query += " WHERE is_hot = 1 ORDER BY cached_at DESC"
	case "rising":
		query += " WHERE is_rising = 1 ORDER BY cached_at DESC"
	default:
		query += " ORDER BY cached_at DESC"
	}

	query += " LIMIT ?"

	err := c.db.Select(&articles, query, limit)
	return articles, err
}

func (c *Cache) GetArticle(id int64) (*models.Article, error) {
	var article models.Article

	query := `
		SELECT id, title, url, source, summary, published_at,
		       upvotes, downvotes, views, comment_count, is_hot, is_rising, cached_at
		FROM articles
		WHERE id = ?
	`

	err := c.db.Get(&article, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &article, err
}

func (c *Cache) MarkArticleRead(articleID int64) error {
	_, err := c.db.Exec(
		"INSERT OR REPLACE INTO read_articles (article_id) VALUES (?)",
		articleID,
	)
	return err
}

// Comment operations

func (c *Cache) SaveComments(comments []models.Comment) error {
	tx, err := c.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT OR REPLACE INTO comments
		(id, article_id, user_id, username, content, parent_id, upvotes, downvotes, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	for _, comment := range comments {
		_, err := tx.Exec(query,
			comment.ID,
			comment.ArticleID,
			comment.UserID,
			"", // username - not in base Comment model
			comment.Content,
			comment.ParentID,
			comment.Upvotes,
			comment.Downvotes,
			comment.CreatedAt,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (c *Cache) GetComments(articleID int64) ([]models.Comment, error) {
	var comments []models.Comment

	query := `
		SELECT id, article_id, user_id, username, content, parent_id,
		       upvotes, downvotes, created_at
		FROM comments
		WHERE article_id = ?
		ORDER BY created_at ASC
	`

	err := c.db.Select(&comments, query, articleID)
	return comments, err
}

// Classifieds operations

func (c *Cache) SaveClassifieds(classifieds []models.Classified) error {
	tx, err := c.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT OR REPLACE INTO classifieds
		(id, user_id, username, title, description, price, category, city, state, is_premium, is_active, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	for _, classified := range classifieds {
		_, err := tx.Exec(query,
			classified.ID,
			classified.UserID,
			"", // username - not in Classified model
			classified.Title,
			classified.Description,
			classified.Price,
			classified.Category,
			classified.City,
			classified.State,
			classified.IsPremium,
			classified.IsActive,
			classified.CreatedAt,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (c *Cache) GetClassifieds(category, location string, limit int) ([]models.Classified, error) {
	var classifieds []models.Classified

	query := `
		SELECT id, user_id, username, title, description, price, category,
		       city, state, is_premium, is_active, created_at
		FROM classifieds
		WHERE is_active = 1
	`

	args := []interface{}{}

	if category != "" {
		query += " AND category = ?"
		args = append(args, category)
	}

	if location != "" {
		query += " AND (city LIKE ? OR state LIKE ?)"
		args = append(args, "%"+location+"%", "%"+location+"%")
	}

	query += " ORDER BY is_premium DESC, created_at DESC LIMIT ?"
	args = append(args, limit)

	err := c.db.Select(&classifieds, query, args...)
	return classifieds, err
}

// Offline queue operations

func (c *Cache) QueueAction(action string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = c.db.Exec(
		"INSERT INTO offline_queue (action, payload) VALUES (?, ?)",
		action, string(data),
	)

	return err
}

type QueuedAction struct {
	ID      int64  `db:"id"`
	Action  string `db:"action"`
	Payload string `db:"payload"`
}

func (c *Cache) GetQueuedActions() ([]QueuedAction, error) {
	var actions []QueuedAction

	err := c.db.Select(&actions, "SELECT id, action, payload FROM offline_queue ORDER BY created_at")
	return actions, err
}

func (c *Cache) DeleteQueuedAction(id int64) error {
	_, err := c.db.Exec("DELETE FROM offline_queue WHERE id = ?", id)
	return err
}

func (c *Cache) ClearQueue() error {
	_, err := c.db.Exec("DELETE FROM offline_queue")
	return err
}

// Settings operations

func (c *Cache) GetSetting(key string) (string, error) {
	var value string
	err := c.db.Get(&value, "SELECT value FROM user_settings WHERE key = ?", key)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return value, err
}

func (c *Cache) SetSetting(key, value string) error {
	_, err := c.db.Exec(
		"INSERT OR REPLACE INTO user_settings (key, value) VALUES (?, ?)",
		key, value,
	)
	return err
}

// Weather cache

func (c *Cache) SaveWeather(location string, weather *models.Weather) error {
	data, err := json.Marshal(weather)
	if err != nil {
		return err
	}

	_, err = c.db.Exec(
		"INSERT OR REPLACE INTO weather_cache (location, data) VALUES (?, ?)",
		location, string(data),
	)

	return err
}

func (c *Cache) GetWeather(location string, maxAge time.Duration) (*models.Weather, error) {
	var data string
	var updatedAt time.Time

	err := c.db.QueryRow(
		"SELECT data, updated_at FROM weather_cache WHERE location = ?",
		location,
	).Scan(&data, &updatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Check if cache is too old
	if time.Since(updatedAt) > maxAge {
		return nil, nil
	}

	var weather models.Weather
	if err := json.Unmarshal([]byte(data), &weather); err != nil {
		return nil, err
	}

	return &weather, nil
}

// Cleanup operations

func (c *Cache) CleanupOldArticles(olderThan time.Time) error {
	_, err := c.db.Exec("DELETE FROM articles WHERE cached_at < ?", olderThan)
	return err
}

func (c *Cache) CleanupOldComments(olderThan time.Time) error {
	_, err := c.db.Exec("DELETE FROM comments WHERE cached_at < ?", olderThan)
	return err
}

func (c *Cache) CleanupOldClassifieds(olderThan time.Time) error {
	_, err := c.db.Exec("DELETE FROM classifieds WHERE cached_at < ?", olderThan)
	return err
}

// Close closes the database connection
func (c *Cache) Close() error {
	return c.db.Close()
}
