package services

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
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

// TODO: Implement article methods
// - GetArticles(feed string, limit, offset int)
// - GetArticle(id int64)
// - GetHotArticles()
// - GetControversialArticles()
// - GetRisingArticles()
