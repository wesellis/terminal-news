package services

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type VoteService struct {
	db  *sqlx.DB
	rdb *redis.Client
}

func NewVoteService(db *sqlx.DB, rdb *redis.Client) *VoteService {
	return &VoteService{db: db, rdb: rdb}
}
