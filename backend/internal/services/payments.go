package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/customer"
	"github.com/stripe/stripe-go/v76/paymentintent"
	"github.com/stripe/stripe-go/v76/price"
	"github.com/stripe/stripe-go/v76/product"
	"github.com/stripe/stripe-go/v76/subscription"
)

var (
	ErrPaymentFailed    = errors.New("payment failed")
	ErrInvalidAmount    = errors.New("invalid payment amount")
	ErrCustomerNotFound = errors.New("customer not found")
)

type PaymentService struct {
	db *sqlx.DB
}

func NewPaymentService(db *sqlx.DB) *PaymentService {
	// Initialize Stripe with API key
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	return &PaymentService{db: db}
}

// Payment represents a payment record
type Payment struct {
	ID                int64     `json:"id" db:"id"`
	UserID            int64     `json:"user_id" db:"user_id"`
	StripeCustomerID  string    `json:"stripe_customer_id" db:"stripe_customer_id"`
	StripePaymentID   string    `json:"stripe_payment_id" db:"stripe_payment_id"`
	Amount            int64     `json:"amount" db:"amount"` // in cents
	Currency          string    `json:"currency" db:"currency"`
	Status            string    `json:"status" db:"status"`
	PaymentType       string    `json:"payment_type" db:"payment_type"` // classified_boost, sponsor_subscription
	RelatedID         *int64    `json:"related_id,omitempty" db:"related_id"` // classified_id or subscription_id
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

// CreateCustomer creates a Stripe customer for a user
func (s *PaymentService) CreateCustomer(ctx context.Context, userID int64, email, name string) (string, error) {
	// Check if customer already exists
	var existingCustomerID string
	err := s.db.GetContext(ctx, &existingCustomerID,
		`SELECT stripe_customer_id FROM users WHERE id = $1 AND stripe_customer_id IS NOT NULL`,
		userID)

	if err == nil && existingCustomerID != "" {
		return existingCustomerID, nil
	}

	// Create new Stripe customer
	params := &stripe.CustomerParams{
		Email: stripe.String(email),
		Name:  stripe.String(name),
		Metadata: map[string]string{
			"user_id": fmt.Sprintf("%d", userID),
		},
	}

	cust, err := customer.New(params)
	if err != nil {
		return "", fmt.Errorf("failed to create Stripe customer: %w", err)
	}

	// Update user with Stripe customer ID
	_, err = s.db.ExecContext(ctx,
		`UPDATE users SET stripe_customer_id = $1 WHERE id = $2`,
		cust.ID, userID)

	if err != nil {
		return "", fmt.Errorf("failed to save Stripe customer ID: %w", err)
	}

	return cust.ID, nil
}

// CreateClassifiedBoostPayment creates a payment for boosting a classified ad
func (s *PaymentService) CreateClassifiedBoostPayment(ctx context.Context, userID, classifiedID int64, durationDays int) (*Payment, error) {
	// Verify user owns the classified
	var ownerID int64
	err := s.db.GetContext(ctx, &ownerID, `SELECT user_id FROM classifieds WHERE id = $1`, classifiedID)
	if err != nil {
		return nil, fmt.Errorf("classified not found")
	}
	if ownerID != userID {
		return nil, fmt.Errorf("unauthorized: you can only boost your own classifieds")
	}

	// Calculate amount based on duration
	// Base price: $5 for 7 days, $10 for 30 days
	var amount int64
	switch durationDays {
	case 7:
		amount = 500 // $5.00
	case 30:
		amount = 1000 // $10.00
	default:
		return nil, ErrInvalidAmount
	}

	// Get user info
	var email, username string
	err = s.db.GetContext(ctx, &email, `SELECT email FROM users WHERE id = $1`, userID)
	if err != nil {
		return nil, err
	}
	err = s.db.GetContext(ctx, &username, `SELECT username FROM users WHERE id = $1`, userID)
	if err != nil {
		username = "User"
	}

	// Get or create Stripe customer
	customerID, err := s.CreateCustomer(ctx, userID, email, username)
	if err != nil {
		return nil, err
	}

	// Create payment intent
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		Customer: stripe.String(customerID),
		Metadata: map[string]string{
			"user_id":       fmt.Sprintf("%d", userID),
			"classified_id": fmt.Sprintf("%d", classifiedID),
			"type":          "classified_boost",
			"duration_days": fmt.Sprintf("%d", durationDays),
		},
		Description: stripe.String(fmt.Sprintf("Boost classified ad #%d for %d days", classifiedID, durationDays)),
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	// Save payment record
	payment := &Payment{
		UserID:           userID,
		StripeCustomerID: customerID,
		StripePaymentID:  pi.ID,
		Amount:           amount,
		Currency:         "usd",
		Status:           string(pi.Status),
		PaymentType:      "classified_boost",
		RelatedID:        &classifiedID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err = s.db.GetContext(ctx, &payment.ID,
		`INSERT INTO payments (user_id, stripe_customer_id, stripe_payment_id, amount, currency, status, payment_type, related_id, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		 RETURNING id`,
		payment.UserID, payment.StripeCustomerID, payment.StripePaymentID, payment.Amount,
		payment.Currency, payment.Status, payment.PaymentType, payment.RelatedID,
		payment.CreatedAt, payment.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to save payment: %w", err)
	}

	return payment, nil
}

// CreateSponsorSubscription creates a recurring subscription for sponsors
func (s *PaymentService) CreateSponsorSubscription(ctx context.Context, userID int64, tier string) (*Payment, error) {
	// Get price ID based on tier
	priceID, amount := s.getSponsorPriceID(tier)
	if priceID == "" {
		return nil, ErrInvalidAmount
	}

	// Get user info
	var email, username string
	err := s.db.GetContext(ctx, &email, `SELECT email FROM users WHERE id = $1`, userID)
	if err != nil {
		return nil, err
	}
	err = s.db.GetContext(ctx, &username, `SELECT username FROM users WHERE id = $1`, userID)
	if err != nil {
		username = "User"
	}

	// Get or create Stripe customer
	customerID, err := s.CreateCustomer(ctx, userID, email, username)
	if err != nil {
		return nil, err
	}

	// Create subscription
	params := &stripe.SubscriptionParams{
		Customer: stripe.String(customerID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(priceID),
			},
		},
		Metadata: map[string]string{
			"user_id": fmt.Sprintf("%d", userID),
			"type":    "sponsor_subscription",
			"tier":    tier,
		},
	}

	sub, err := subscription.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	// Save payment record
	payment := &Payment{
		UserID:           userID,
		StripeCustomerID: customerID,
		StripePaymentID:  sub.ID,
		Amount:           amount,
		Currency:         "usd",
		Status:           string(sub.Status),
		PaymentType:      "sponsor_subscription",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err = s.db.GetContext(ctx, &payment.ID,
		`INSERT INTO payments (user_id, stripe_customer_id, stripe_payment_id, amount, currency, status, payment_type, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		 RETURNING id`,
		payment.UserID, payment.StripeCustomerID, payment.StripePaymentID, payment.Amount,
		payment.Currency, payment.Status, payment.PaymentType, payment.CreatedAt, payment.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to save payment: %w", err)
	}

	// Update user to sponsor status
	_, err = s.db.ExecContext(ctx,
		`UPDATE users SET account_tier = 'sponsor', sponsor_tier = $1, sponsor_expires_at = NULL WHERE id = $2`,
		tier, userID)

	if err != nil {
		return nil, fmt.Errorf("failed to update user to sponsor: %w", err)
	}

	return payment, nil
}

// HandlePaymentSuccess handles successful payment (called from webhook)
func (s *PaymentService) HandlePaymentSuccess(ctx context.Context, paymentIntentID string) error {
	// Update payment status
	_, err := s.db.ExecContext(ctx,
		`UPDATE payments SET status = 'succeeded', updated_at = $1 WHERE stripe_payment_id = $2`,
		time.Now(), paymentIntentID)

	if err != nil {
		return err
	}

	// Get payment details
	var payment Payment
	err = s.db.GetContext(ctx, &payment,
		`SELECT * FROM payments WHERE stripe_payment_id = $1`, paymentIntentID)

	if err != nil {
		return err
	}

	// Handle based on payment type
	switch payment.PaymentType {
	case "classified_boost":
		if payment.RelatedID != nil {
			// Activate classified boost
			_, err = s.db.ExecContext(ctx,
				`UPDATE classifieds SET is_premium = true, premium_expires_at = NOW() + INTERVAL '7 days' WHERE id = $1`,
				*payment.RelatedID)
		}
	case "sponsor_subscription":
		// Already handled in CreateSponsorSubscription
	}

	return err
}

// HandleSubscriptionCanceled handles canceled subscriptions (called from webhook)
func (s *PaymentService) HandleSubscriptionCanceled(ctx context.Context, subscriptionID string) error {
	// Get subscription payment
	var payment Payment
	err := s.db.GetContext(ctx, &payment,
		`SELECT * FROM payments WHERE stripe_payment_id = $1 AND payment_type = 'sponsor_subscription'`,
		subscriptionID)

	if err != nil {
		return err
	}

	// Downgrade user from sponsor
	_, err = s.db.ExecContext(ctx,
		`UPDATE users SET account_tier = 'free', sponsor_tier = NULL, sponsor_expires_at = NOW() + INTERVAL '30 days' WHERE id = $1`,
		payment.UserID)

	return err
}

// GetPaymentHistory retrieves payment history for a user
func (s *PaymentService) GetPaymentHistory(ctx context.Context, userID int64) ([]Payment, error) {
	var payments []Payment
	err := s.db.SelectContext(ctx, &payments,
		`SELECT * FROM payments WHERE user_id = $1 ORDER BY created_at DESC LIMIT 50`,
		userID)

	return payments, err
}

// getSponsorPriceID returns the Stripe price ID for sponsor tiers
func (s *PaymentService) getSponsorPriceID(tier string) (string, int64) {
	// These will be created in Stripe dashboard or via API
	switch tier {
	case "bronze":
		return "price_bronze_sponsor", 999  // $9.99/month
	case "silver":
		return "price_silver_sponsor", 2999 // $29.99/month
	case "gold":
		return "price_gold_sponsor", 9999   // $99.99/month
	default:
		return "", 0
	}
}

// SetupStripeProducts creates products and prices in Stripe (one-time setup)
func (s *PaymentService) SetupStripeProducts() error {
	// Create Terminal News product
	productParams := &stripe.ProductParams{
		Name:        stripe.String("Terminal News Premium Features"),
		Description: stripe.String("Premium features for Terminal News"),
	}
	prod, err := product.New(productParams)
	if err != nil {
		return err
	}

	// Create price for classified boost (one-time $5)
	_, err = price.New(&stripe.PriceParams{
		Product:    stripe.String(prod.ID),
		Currency:   stripe.String(string(stripe.CurrencyUSD)),
		UnitAmount: stripe.Int64(500),
		Nickname:   stripe.String("7-day classified boost"),
	})
	if err != nil {
		return err
	}

	// Create sponsor subscription prices
	tiers := []struct {
		name   string
		amount int64
	}{
		{"Bronze Sponsor", 999},
		{"Silver Sponsor", 2999},
		{"Gold Sponsor", 9999},
	}

	for _, tier := range tiers {
		_, err = price.New(&stripe.PriceParams{
			Product:    stripe.String(prod.ID),
			Currency:   stripe.String(string(stripe.CurrencyUSD)),
			UnitAmount: stripe.Int64(tier.amount),
			Recurring: &stripe.PriceRecurringParams{
				Interval: stripe.String(string(stripe.PriceRecurringIntervalMonth)),
			},
			Nickname: stripe.String(tier.name),
		})
		if err != nil {
			return err
		}
	}

	return nil
}
