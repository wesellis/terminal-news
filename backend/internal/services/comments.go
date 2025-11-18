package services

import "github.com/jmoiron/sqlx"

type CommentService struct {
	db *sqlx.DB
}

func NewCommentService(db *sqlx.DB) *CommentService {
	return &CommentService{db: db}
}
