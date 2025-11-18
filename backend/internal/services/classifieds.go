package services

import "github.com/jmoiron/sqlx"

type ClassifiedService struct {
	db *sqlx.DB
}

func NewClassifiedService(db *sqlx.DB) *ClassifiedService {
	return &ClassifiedService{db: db}
}
