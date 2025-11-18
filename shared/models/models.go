package models

import "time"

// User represents a user in the system
type User struct {
	ID             int64     `json:"id" db:"id"`
	Username       string    `json:"username" db:"username"`
	Email          string    `json:"email" db:"email"`
	PasswordHash   string    `json:"-" db:"password_hash"`
	DisplayName    string    `json:"display_name" db:"display_name"`
	Bio            string    `json:"bio" db:"bio"`
	Location       string    `json:"location" db:"location"`
	Karma          int       `json:"karma" db:"karma"`
	TrustScore     float64   `json:"trust_score" db:"trust_score"`
	EmailVerified  bool      `json:"email_verified" db:"email_verified"`
	IsBanned       bool      `json:"is_banned" db:"is_banned"`
	IsModerator    bool      `json:"is_moderator" db:"is_moderator"`
	IsAdmin        bool      `json:"is_admin" db:"is_admin"`
	LastActiveAt   time.Time `json:"last_active_at" db:"last_active_at"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// Article represents a news article
type Article struct {
	ID          int64     `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	URL         string    `json:"url" db:"url"`
	Content     string    `json:"content" db:"content"`
	ImageURL    string    `json:"image_url" db:"image_url"`
	Source      string    `json:"source" db:"source"`
	Author      string    `json:"author" db:"author"`
	PublishedAt time.Time `json:"published_at" db:"published_at"`
	Category    string    `json:"category" db:"category"`
	Tags        []string  `json:"tags" db:"tags"`
	ExternalID  string    `json:"external_id" db:"external_id"`
	FetchSource string    `json:"fetch_source" db:"fetch_source"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Vote represents a user vote on an article
type Vote struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	ArticleID int64     `json:"article_id" db:"article_id"`
	VoteType  string    `json:"vote_type" db:"vote_type"` // open, like, dislike
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Comment represents a comment on an article
type Comment struct {
	ID         int64     `json:"id" db:"id"`
	UserID     int64     `json:"user_id" db:"user_id"`
	ArticleID  int64     `json:"article_id" db:"article_id"`
	ParentID   *int64    `json:"parent_id" db:"parent_id"`
	Content    string    `json:"content" db:"content"`
	IsDeleted  bool      `json:"is_deleted" db:"is_deleted"`
	IsFlagged  bool      `json:"is_flagged" db:"is_flagged"`
	FlagCount  int       `json:"flag_count" db:"flag_count"`
	Upvotes    int       `json:"upvotes" db:"upvotes"`
	Downvotes  int       `json:"downvotes" db:"downvotes"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	EditedAt   *time.Time `json:"edited_at" db:"edited_at"`
}

// Classified represents a classified ad
type Classified struct {
	ID             int64      `json:"id" db:"id"`
	UserID         int64      `json:"user_id" db:"user_id"`
	Title          string     `json:"title" db:"title"`
	Description    string     `json:"description" db:"description"`
	Price          *float64   `json:"price" db:"price"`
	Category       string     `json:"category" db:"category"`
	Subcategory    string     `json:"subcategory" db:"subcategory"`
	City           string     `json:"city" db:"city"`
	State          string     `json:"state" db:"state"`
	Country        string     `json:"country" db:"country"`
	ContactEmail   string     `json:"contact_email" db:"contact_email"`
	ContactPhone   string     `json:"contact_phone" db:"contact_phone"`
	ContactMethod  string     `json:"contact_method" db:"contact_method"`
	IsActive       bool       `json:"is_active" db:"is_active"`
	IsPremium      bool       `json:"is_premium" db:"is_premium"`
	IsFlagged      bool       `json:"is_flagged" db:"is_flagged"`
	PremiumUntil   *time.Time `json:"premium_until" db:"premium_until"`
	LastBoostedAt  *time.Time `json:"last_boosted_at" db:"last_boosted_at"`
	ViewCount      int        `json:"view_count" db:"view_count"`
	ContactCount   int        `json:"contact_count" db:"contact_count"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
	ExpiresAt      *time.Time `json:"expires_at" db:"expires_at"`
}

// ArticleRanking represents cached article rankings
type ArticleRanking struct {
	ArticleID          int64     `json:"article_id" db:"article_id"`
	OpenCount          int       `json:"open_count" db:"open_count"`
	LikeCount          int       `json:"like_count" db:"like_count"`
	DislikeCount       int       `json:"dislike_count" db:"dislike_count"`
	TotalScore         int       `json:"total_score" db:"total_score"`
	ControversyScore   float64   `json:"controversy_score" db:"controversy_score"`
	TotalEngagement    int       `json:"total_engagement" db:"total_engagement"`
	HoursSincePublished float64  `json:"hours_since_published" db:"hours_since_published"`
	HotRank            float64   `json:"hot_rank" db:"hot_rank"`
	LastUpdated        time.Time `json:"last_updated" db:"last_updated"`
}

// APIResponse is a generic API response wrapper
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// PaginatedResponse wraps paginated data
type PaginatedResponse struct {
	Items      interface{} `json:"items"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	TotalPages int         `json:"total_pages"`
}

// ArticleWithRanking combines article data with ranking info
type ArticleWithRanking struct {
	Article
	ArticleRanking
	UserVote *string `json:"user_vote,omitempty"` // user's vote if logged in
}

// CommentWithUser includes user information with comment
type CommentWithUser struct {
	Comment
	Username string `json:"username"`
	Karma    int    `json:"karma"`
}
