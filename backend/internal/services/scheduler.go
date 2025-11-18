package services

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Scheduler struct {
	db  *sqlx.DB
	rdb *redis.Client
}

func NewScheduler(db *sqlx.DB, rdb *redis.Client) *Scheduler {
	return &Scheduler{db: db, rdb: rdb}
}

func (s *Scheduler) Start(ctx context.Context) {
	// TODO: Implement background jobs
	// - Refresh rankings every 5 minutes
	// - Expire classifieds hourly
	// - Clean up audit logs daily
}
