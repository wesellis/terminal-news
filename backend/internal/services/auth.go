package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db *sqlx.DB
}

type User struct {
	ID             int64          `db:"id" json:"id"`
	Username       string         `db:"username" json:"username"`
	Email          string         `db:"email" json:"email"`
	PasswordHash   string         `db:"password_hash" json:"-"`
	DisplayName    sql.NullString `db:"display_name" json:"display_name,omitempty"`
	Bio            sql.NullString `db:"bio" json:"bio,omitempty"`
	Location       sql.NullString `db:"location" json:"location,omitempty"`
	Karma          int            `db:"karma" json:"karma"`
	TrustScore     float64        `db:"trust_score" json:"trust_score"`
	IsBanned       bool           `db:"is_banned" json:"-"`
	IsModerator    bool           `db:"is_moderator" json:"is_moderator"`
	IsAdmin        bool           `db:"is_admin" json:"is_admin"`
	EmailVerified  bool           `db:"email_verified" json:"email_verified"`
	LastActiveAt   sql.NullTime   `db:"last_active_at" json:"last_active_at,omitempty"`
	CreatedAt      time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at" json:"updated_at"`
}

type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrUsernameTaken     = errors.New("username already taken")
	ErrEmailTaken        = errors.New("email already taken")
	ErrInvalidToken      = errors.New("invalid token")
	ErrUserBanned        = errors.New("user is banned")
)

func NewAuthService(db *sqlx.DB) *AuthService {
	return &AuthService{db: db}
}

// Register creates a new user account
func (s *AuthService) Register(ctx context.Context, username, email, password string) (*User, error) {
	// Validate input
	if len(username) < 3 || len(username) > 50 {
		return nil, errors.New("username must be between 3 and 50 characters")
	}
	if len(password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Insert user
	query := `
		INSERT INTO users (username, email, password_hash, trust_score, karma)
		VALUES ($1, $2, $3, 0.5, 0)
		RETURNING id, username, email, karma, trust_score, is_moderator, is_admin,
		          email_verified, created_at, updated_at
	`

	user := &User{}
	err = s.db.QueryRowxContext(ctx, query, username, email, string(hashedPassword)).StructScan(user)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_username_key\"" {
			return nil, ErrUsernameTaken
		}
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			return nil, ErrEmailTaken
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(ctx context.Context, username, password string) (*User, *TokenPair, error) {
	// Get user by username or email
	query := `
		SELECT id, username, email, password_hash, karma, trust_score,
		       is_banned, is_moderator, is_admin, email_verified,
		       created_at, updated_at
		FROM users
		WHERE username = $1 OR email = $1
	`

	user := &User{}
	err := s.db.GetContext(ctx, user, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, ErrUserNotFound
		}
		return nil, nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user is banned
	if user.IsBanned {
		return nil, nil, ErrUserBanned
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, nil, ErrInvalidPassword
	}

	// Generate tokens
	tokens, err := s.GenerateTokens(user)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Update last active
	_, err = s.db.ExecContext(ctx, "UPDATE users SET last_active_at = NOW() WHERE id = $1", user.ID)
	if err != nil {
		// Log but don't fail
		fmt.Printf("Warning: failed to update last_active_at: %v\n", err)
	}

	return user, tokens, nil
}

// GenerateTokens creates access and refresh tokens
func (s *AuthService) GenerateTokens(user *User) (*TokenPair, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET not set")
	}

	// Access token (15 minutes)
	accessExpiry := time.Now().Add(15 * time.Minute)
	accessClaims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Refresh token (7 days)
	refreshExpiry := time.Now().Add(7 * 24 * time.Hour)
	refreshClaims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresAt:    accessExpiry.Unix(),
	}, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET not set")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(ctx context.Context, userID int64) (*User, error) {
	query := `
		SELECT id, username, email, karma, trust_score, display_name, bio, location,
		       is_banned, is_moderator, is_admin, email_verified,
		       last_active_at, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &User{}
	err := s.db.GetContext(ctx, user, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user.IsBanned {
		return nil, ErrUserBanned
	}

	return user, nil
}

// RefreshAccessToken generates a new access token from a refresh token
func (s *AuthService) RefreshAccessToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
	// Validate refresh token
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// Get user to ensure they still exist and aren't banned
	user, err := s.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	// Generate new token pair
	return s.GenerateTokens(user)
}
