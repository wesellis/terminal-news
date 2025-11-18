package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

var (
	ErrClassifiedNotFound = errors.New("classified not found")
	ErrClassifiedInactive = errors.New("classified is inactive")
)

type ClassifiedService struct {
	db *sqlx.DB
}

func NewClassifiedService(db *sqlx.DB) *ClassifiedService {
	return &ClassifiedService{db: db}
}

// Classified represents a classified ad listing
type Classified struct {
	ID              int64          `json:"id" db:"id"`
	UserID          int64          `json:"user_id" db:"user_id"`
	Title           string         `json:"title" db:"title"`
	Description     string         `json:"description" db:"description"`
	Price           *float64       `json:"price,omitempty" db:"price"`
	Category        string         `json:"category" db:"category"`
	Subcategory     *string        `json:"subcategory,omitempty" db:"subcategory"`
	City            string         `json:"city" db:"city"`
	State           *string        `json:"state,omitempty" db:"state"`
	Country         string         `json:"country" db:"country"`
	Lat             *float64       `json:"lat,omitempty" db:"lat"`
	Lng             *float64       `json:"lng,omitempty" db:"lng"`
	ContactEmail    *string        `json:"contact_email,omitempty" db:"contact_email"`
	ContactPhone    *string        `json:"contact_phone,omitempty" db:"contact_phone"`
	ContactMethod   string         `json:"contact_method" db:"contact_method"`
	IsActive        bool           `json:"is_active" db:"is_active"`
	IsPremium       bool           `json:"is_premium" db:"is_premium"`
	IsFlagged       bool           `json:"is_flagged" db:"is_flagged"`
	FlagCount       int            `json:"flag_count" db:"flag_count"`
	PremiumUntil    sql.NullTime   `json:"premium_until,omitempty" db:"premium_until"`
	BoostCount      int            `json:"boost_count" db:"boost_count"`
	LastBoostedAt   sql.NullTime   `json:"last_boosted_at,omitempty" db:"last_boosted_at"`
	ViewCount       int            `json:"view_count" db:"view_count"`
	ContactCount    int            `json:"contact_count" db:"contact_count"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at" db:"updated_at"`
	ExpiresAt       sql.NullTime   `json:"expires_at,omitempty" db:"expires_at"`

	// Joined data
	Username        string         `json:"username,omitempty" db:"username"`
}

// CreateClassifiedRequest is the request body for creating a classified
type CreateClassifiedRequest struct {
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Price         *float64 `json:"price,omitempty"`
	Category      string   `json:"category"`
	Subcategory   *string  `json:"subcategory,omitempty"`
	City          string   `json:"city"`
	State         *string  `json:"state,omitempty"`
	Country       string   `json:"country"`
	Lat           *float64 `json:"lat,omitempty"`
	Lng           *float64 `json:"lng,omitempty"`
	ContactEmail  *string  `json:"contact_email,omitempty"`
	ContactPhone  *string  `json:"contact_phone,omitempty"`
	ContactMethod string   `json:"contact_method"`
	ExpiresInDays int      `json:"expires_in_days"` // Default 30 days
}

// UpdateClassifiedRequest is the request body for updating a classified
type UpdateClassifiedRequest struct {
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Price         *float64 `json:"price,omitempty"`
	ContactEmail  *string  `json:"contact_email,omitempty"`
	ContactPhone  *string  `json:"contact_phone,omitempty"`
	ContactMethod string   `json:"contact_method"`
}

// ClassifiedListResponse is the response for classified list endpoints
type ClassifiedListResponse struct {
	Classifieds []Classified `json:"classifieds"`
	TotalCount  int          `json:"total_count"`
	Limit       int          `json:"limit"`
	Offset      int          `json:"offset"`
}

// CreateClassified creates a new classified ad
func (s *ClassifiedService) CreateClassified(ctx context.Context, userID int64, req *CreateClassifiedRequest) (*Classified, error) {
	// Validate input
	if len(req.Title) < 10 {
		return nil, errors.New("title must be at least 10 characters")
	}
	if len(req.Description) < 20 {
		return nil, errors.New("description must be at least 20 characters")
	}

	// Set default expiration
	expiresInDays := req.ExpiresInDays
	if expiresInDays <= 0 {
		expiresInDays = 30 // Default 30 days
	}
	expiresAt := time.Now().AddDate(0, 0, expiresInDays)

	// Insert classified
	query := `
		INSERT INTO classifieds (
			user_id, title, description, price, category, subcategory,
			city, state, country, lat, lng,
			contact_email, contact_phone, contact_method,
			expires_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id, user_id, title, description, price, category, subcategory,
		          city, state, country, lat, lng, contact_email, contact_phone, contact_method,
		          is_active, is_premium, is_flagged, flag_count, premium_until,
		          boost_count, last_boosted_at, view_count, contact_count,
		          created_at, updated_at, expires_at
	`

	var classified Classified
	err := s.db.GetContext(ctx, &classified, query,
		userID, req.Title, req.Description, req.Price, req.Category, req.Subcategory,
		req.City, req.State, req.Country, req.Lat, req.Lng,
		req.ContactEmail, req.ContactPhone, req.ContactMethod,
		expiresAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create classified: %w", err)
	}

	return &classified, nil
}

// GetClassifieds retrieves a list of active classifieds with optional filtering
func (s *ClassifiedService) GetClassifieds(ctx context.Context, category, city, state string, limit, offset int) (*ClassifiedListResponse, error) {
	// Build dynamic query
	baseQuery := `
		SELECT
			c.id, c.user_id, c.title, c.description, c.price, c.category, c.subcategory,
			c.city, c.state, c.country, c.lat, c.lng,
			c.contact_email, c.contact_phone, c.contact_method,
			c.is_active, c.is_premium, c.is_flagged, c.flag_count, c.premium_until,
			c.boost_count, c.last_boosted_at, c.view_count, c.contact_count,
			c.created_at, c.updated_at, c.expires_at,
			u.username
		FROM classifieds c
		JOIN users u ON c.user_id = u.id
		WHERE c.is_active = TRUE
	`

	countQuery := `SELECT COUNT(*) FROM classifieds c WHERE c.is_active = TRUE`

	args := []interface{}{}
	argIndex := 1

	// Add filters
	if category != "" {
		baseQuery += fmt.Sprintf(" AND c.category = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND category = $%d", argIndex)
		args = append(args, category)
		argIndex++
	}

	if city != "" {
		baseQuery += fmt.Sprintf(" AND LOWER(c.city) = LOWER($%d)", argIndex)
		countQuery += fmt.Sprintf(" AND LOWER(city) = LOWER($%d)", argIndex)
		args = append(args, city)
		argIndex++
	}

	if state != "" {
		baseQuery += fmt.Sprintf(" AND LOWER(c.state) = LOWER($%d)", argIndex)
		countQuery += fmt.Sprintf(" AND LOWER(state) = LOWER($%d)", argIndex)
		args = append(args, state)
		argIndex++
	}

	// Order by: premium first, then most recent, then most boosted
	baseQuery += ` ORDER BY c.is_premium DESC, c.created_at DESC, c.boost_count DESC`

	// Add pagination
	baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	// Execute query
	var classifieds []Classified
	if err := s.db.SelectContext(ctx, &classifieds, baseQuery, args...); err != nil {
		return nil, fmt.Errorf("failed to get classifieds: %w", err)
	}

	// Get total count
	var totalCount int
	countArgs := args[:len(args)-2] // Remove limit and offset
	if err := s.db.GetContext(ctx, &totalCount, countQuery, countArgs...); err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	return &ClassifiedListResponse{
		Classifieds: classifieds,
		TotalCount:  totalCount,
		Limit:       limit,
		Offset:      offset,
	}, nil
}

// GetClassified retrieves a single classified by ID and increments view count
func (s *ClassifiedService) GetClassified(ctx context.Context, id int64) (*Classified, error) {
	// Get classified
	query := `
		SELECT
			c.id, c.user_id, c.title, c.description, c.price, c.category, c.subcategory,
			c.city, c.state, c.country, c.lat, c.lng,
			c.contact_email, c.contact_phone, c.contact_method,
			c.is_active, c.is_premium, c.is_flagged, c.flag_count, c.premium_until,
			c.boost_count, c.last_boosted_at, c.view_count, c.contact_count,
			c.created_at, c.updated_at, c.expires_at,
			u.username
		FROM classifieds c
		JOIN users u ON c.user_id = u.id
		WHERE c.id = $1
	`

	var classified Classified
	if err := s.db.GetContext(ctx, &classified, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrClassifiedNotFound
		}
		return nil, fmt.Errorf("failed to get classified: %w", err)
	}

	// Increment view count
	_, _ = s.db.ExecContext(ctx, `UPDATE classifieds SET view_count = view_count + 1 WHERE id = $1`, id)

	return &classified, nil
}

// UpdateClassified updates a classified ad
func (s *ClassifiedService) UpdateClassified(ctx context.Context, id, userID int64, req *UpdateClassifiedRequest) (*Classified, error) {
	// Check ownership
	var ownerID int64
	err := s.db.GetContext(ctx, &ownerID, `SELECT user_id FROM classifieds WHERE id = $1`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrClassifiedNotFound
		}
		return nil, fmt.Errorf("failed to check ownership: %w", err)
	}

	if ownerID != userID {
		return nil, ErrUnauthorized
	}

	// Update classified
	query := `
		UPDATE classifieds
		SET title = $1, description = $2, price = $3,
		    contact_email = $4, contact_phone = $5, contact_method = $6,
		    updated_at = NOW()
		WHERE id = $7 AND is_active = TRUE
		RETURNING id, user_id, title, description, price, category, subcategory,
		          city, state, country, lat, lng, contact_email, contact_phone, contact_method,
		          is_active, is_premium, is_flagged, flag_count, premium_until,
		          boost_count, last_boosted_at, view_count, contact_count,
		          created_at, updated_at, expires_at
	`

	var classified Classified
	err = s.db.GetContext(ctx, &classified, query,
		req.Title, req.Description, req.Price,
		req.ContactEmail, req.ContactPhone, req.ContactMethod,
		id,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrClassifiedInactive
		}
		return nil, fmt.Errorf("failed to update classified: %w", err)
	}

	return &classified, nil
}

// DeleteClassified soft-deletes a classified (sets is_active = false)
func (s *ClassifiedService) DeleteClassified(ctx context.Context, id, userID int64) error {
	// Check ownership
	var ownerID int64
	err := s.db.GetContext(ctx, &ownerID, `SELECT user_id FROM classifieds WHERE id = $1`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrClassifiedNotFound
		}
		return fmt.Errorf("failed to check ownership: %w", err)
	}

	if ownerID != userID {
		return ErrUnauthorized
	}

	// Soft delete
	query := `UPDATE classifieds SET is_active = FALSE, updated_at = NOW() WHERE id = $1`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete classified: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrClassifiedNotFound
	}

	return nil
}

// SearchClassifiedsByLocation finds classifieds near a geographic point
func (s *ClassifiedService) SearchClassifiedsByLocation(ctx context.Context, lat, lng, radiusMiles float64, limit, offset int) (*ClassifiedListResponse, error) {
	// Use Haversine formula to calculate distance
	// radiusMiles converted to degrees (approximate: 1 degree ≈ 69 miles)
	radiusDegrees := radiusMiles / 69.0

	query := `
		SELECT
			c.id, c.user_id, c.title, c.description, c.price, c.category, c.subcategory,
			c.city, c.state, c.country, c.lat, c.lng,
			c.contact_email, c.contact_phone, c.contact_method,
			c.is_active, c.is_premium, c.is_flagged, c.flag_count, c.premium_until,
			c.boost_count, c.last_boosted_at, c.view_count, c.contact_count,
			c.created_at, c.updated_at, c.expires_at,
			u.username,
			(
				3959 * acos(
					cos(radians($1)) * cos(radians(c.lat)) *
					cos(radians(c.lng) - radians($2)) +
					sin(radians($1)) * sin(radians(c.lat))
				)
			) AS distance
		FROM classifieds c
		JOIN users u ON c.user_id = u.id
		WHERE c.is_active = TRUE
		  AND c.lat IS NOT NULL
		  AND c.lng IS NOT NULL
		  AND c.lat BETWEEN $1 - $3 AND $1 + $3
		  AND c.lng BETWEEN $2 - $3 AND $2 + $3
		ORDER BY distance ASC, c.is_premium DESC
		LIMIT $4 OFFSET $5
	`

	var classifieds []Classified
	if err := s.db.SelectContext(ctx, &classifieds, query, lat, lng, radiusDegrees, limit, offset); err != nil {
		return nil, fmt.Errorf("failed to search classifieds by location: %w", err)
	}

	// Get total count (approximate - within bounding box)
	countQuery := `
		SELECT COUNT(*)
		FROM classifieds
		WHERE is_active = TRUE
		  AND lat IS NOT NULL
		  AND lng IS NOT NULL
		  AND lat BETWEEN $1 - $3 AND $1 + $3
		  AND lng BETWEEN $2 - $3 AND $2 + $3
	`

	var totalCount int
	if err := s.db.GetContext(ctx, &totalCount, countQuery, lat, lng, radiusDegrees); err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	return &ClassifiedListResponse{
		Classifieds: classifieds,
		TotalCount:  totalCount,
		Limit:       limit,
		Offset:      offset,
	}, nil
}

// IncrementContactCount increments the contact count for a classified
func (s *ClassifiedService) IncrementContactCount(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, `UPDATE classifieds SET contact_count = contact_count + 1 WHERE id = $1`, id)
	return err
}
