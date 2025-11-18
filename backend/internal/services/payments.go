package services

import "github.com/jmoiron/sqlx"

type PaymentService struct {
	db *sqlx.DB
}

func NewPaymentService(db *sqlx.DB) *PaymentService {
	return &PaymentService{db: db}
}
