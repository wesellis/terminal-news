package models

// This file now imports from shared models and adds CLI-specific extensions
// The canonical models live in shared/models/models.go

import (
	"time"

	shared "github.com/wesellis/terminal-news/shared/models"
)

// Re-export shared models for compatibility
type Article = shared.Article
type Comment = shared.Comment
type Classified = shared.Classified
type User = shared.User
type Vote = shared.Vote

// ArticleWithRanking combines article data with ranking info (from shared)
type ArticleWithRanking = shared.ArticleWithRanking

// CommentWithUser includes user information with comment (from shared)
type CommentWithUser = shared.CommentWithUser

// CLI-specific extensions

// Activity represents a user activity item (CLI-specific)
type Activity struct {
	ID        int64     `json:"id"`
	Type      string    `json:"type"` // vote, comment, classified
	ItemID    int64     `json:"item_id"`
	ItemTitle string    `json:"item_title"`
	Action    string    `json:"action"`
	CreatedAt time.Time `json:"created_at"`
}

// Weather represents current weather data (CLI-specific)
type Weather struct {
	Location  string         `json:"location"`
	Current   CurrentWeather `json:"current"`
	Forecast  []DayForecast  `json:"forecast"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type CurrentWeather struct {
	Temperature int     `json:"temperature"`
	FeelsLike   int     `json:"feels_like"`
	Condition   string  `json:"condition"`
	Humidity    int     `json:"humidity"`
	WindSpeed   int     `json:"wind_speed"`
	WindDir     string  `json:"wind_dir"`
	Pressure    float64 `json:"pressure"`
}

type DayForecast struct {
	Day       string `json:"day"`
	High      int    `json:"high"`
	Low       int    `json:"low"`
	Condition string `json:"condition"`
}

// API Response types (CLI-specific)
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type ArticlesResponse struct {
	Articles []ArticleWithRanking `json:"articles"` // Use ArticleWithRanking from shared
	Total    int                  `json:"total"`
	Page     int                  `json:"page"`
}

type CommentsResponse struct {
	Comments []CommentWithUser `json:"comments"` // Use CommentWithUser from shared
	Total    int               `json:"total"`
}

type ClassifiedsResponse struct {
	Classifieds []Classified `json:"classifieds"`
	Total       int          `json:"total"`
}

// WebSocket message types (CLI-specific)
type WSMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type VoteUpdate struct {
	ArticleID int64 `json:"article_id"`
	Upvotes   int   `json:"upvotes"`
	Downvotes int   `json:"downvotes"`
}

// CommentTree is a helper for building threaded comment displays (CLI-specific)
type CommentTree struct {
	Comment
	Username string        `json:"username"`
	Karma    int           `json:"karma"`
	Depth    int           `json:"depth"`
	Children []CommentTree `json:"children"`
}
