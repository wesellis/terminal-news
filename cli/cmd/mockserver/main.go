package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Mock data structures matching shared models
type Article struct {
	ID              int64     `json:"id"`
	Title           string    `json:"title"`
	URL             string    `json:"url"`
	Content         string    `json:"content"`
	ImageURL        string    `json:"image_url"`
	Source          string    `json:"source"`
	Author          string    `json:"author"`
	PublishedAt     time.Time `json:"published_at"`
	Category        string    `json:"category"`
	Tags            []string  `json:"tags"`
	// Ranking fields (embedded ArticleRanking)
	OpenCount           int     `json:"open_count"`
	LikeCount           int     `json:"like_count"`
	DislikeCount        int     `json:"dislike_count"`
	TotalScore          int     `json:"total_score"`
	ControversyScore    float64 `json:"controversy_score"`
	TotalEngagement     int     `json:"total_engagement"`
	HoursSincePublished float64 `json:"hours_since_published"`
	HotRank             float64 `json:"hot_rank"`
}

type Comment struct {
	ID         int64      `json:"id"`
	ArticleID  int64      `json:"article_id"`
	ParentID   *int64     `json:"parent_id"`
	UserID     int64      `json:"user_id"`
	Content    string     `json:"content"`
	Upvotes    int        `json:"upvotes"`
	Downvotes  int        `json:"downvotes"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	EditedAt   *time.Time `json:"edited_at"`
	IsDeleted  bool       `json:"is_deleted"`
	IsFlagged  bool       `json:"is_flagged"`
	FlagCount  int        `json:"flag_count"`
	// CommentWithUser fields
	Username string `json:"username"`
	Karma    int    `json:"karma"`
}

type Classified struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Category     string    `json:"category"`
	Price        float64   `json:"price"`
	City         string    `json:"city"`
	State        string    `json:"state"`
	ContactEmail string    `json:"contact_email"`
	ContactPhone string    `json:"contact_phone"`
	IsPremium    bool      `json:"is_premium"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type Weather struct {
	Location  string          `json:"location"`
	Current   CurrentWeather  `json:"current"`
	Forecast  []ForecastDay   `json:"forecast"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type CurrentWeather struct {
	Temperature int    `json:"temperature"`
	Condition   string `json:"condition"`
	Humidity    int    `json:"humidity"`
	WindSpeed   int    `json:"wind_speed"`
}

type ForecastDay struct {
	Date      string `json:"date"`
	High      int    `json:"high"`
	Low       int    `json:"low"`
	Condition string `json:"condition"`
}

func main() {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	// Health check
	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, 200, map[string]string{"status": "ok"})
	})

	// Auth endpoints
	r.Post("/api/auth/login", handleLogin)
	r.Post("/api/auth/register", handleRegister)
	r.Post("/api/auth/logout", handleLogout)

	// Article endpoints
	r.Get("/api/articles", handleGetArticles)
	r.Get("/api/articles/{id}", handleGetArticle)
	r.Post("/api/articles/{id}/vote", handleVoteArticle)

	// Comment endpoints
	r.Get("/api/articles/{id}/comments", handleGetComments)
	r.Post("/api/articles/{id}/comments", handlePostComment)

	// Classifieds endpoints
	r.Get("/api/classifieds", handleGetClassifieds)
	r.Post("/api/classifieds", handlePostClassified)

	// Weather endpoint
	r.Get("/api/weather", handleGetWeather)

	// Profile endpoints
	r.Get("/api/profile", handleGetProfile)
	r.Get("/api/users/{id}/activity", handleGetActivity)

	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println("  Terminal News - Mock API Server")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println("")
	log.Println("  Running on: http://localhost:8080")
	log.Println("  Health:     http://localhost:8080/api/health")
	log.Println("")
	log.Println("  Press Ctrl+C to stop")
	log.Println("")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 200, map[string]interface{}{
		"token":    "mock-jwt-token-" + fmt.Sprint(time.Now().Unix()),
		"username": "testuser",
		"user_id":  1,
	})
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 201, map[string]interface{}{
		"token":    "mock-jwt-token-" + fmt.Sprint(time.Now().Unix()),
		"username": "newuser",
		"user_id":  2,
	})
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 200, map[string]string{"message": "Logged out successfully"})
}

func handleGetArticles(w http.ResponseWriter, r *http.Request) {
	feed := r.URL.Query().Get("feed")
	if feed == "" {
		feed = "hot"
	}

	articles := generateMockArticles(50, feed)
	respondJSON(w, 200, map[string]interface{}{
		"articles": articles,
		"total":    50,
	})
}

func handleGetArticle(w http.ResponseWriter, r *http.Request) {
	_ = chi.URLParam(r, "id") // articleID not used in mock
	article := generateMockArticles(1, "hot")[0]
	article.ID = 1
	respondJSON(w, 200, map[string]interface{}{
		"article": article,
	})
}

func handleVoteArticle(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 200, map[string]string{"message": "Vote recorded"})
}

func handleGetComments(w http.ResponseWriter, r *http.Request) {
	comments := generateMockComments(20)
	respondJSON(w, 200, map[string]interface{}{
		"comments": comments,
		"total":    20,
	})
}

func handlePostComment(w http.ResponseWriter, r *http.Request) {
	comment := Comment{
		ID:        rand.Int63n(10000),
		ArticleID: 1,
		Username:  "testuser",
		Content:   "This is a test comment",
		CreatedAt: time.Now(),
	}
	respondJSON(w, 201, map[string]interface{}{
		"comment": comment,
	})
}

func handleGetClassifieds(w http.ResponseWriter, r *http.Request) {
	classifieds := generateMockClassifieds(30)
	respondJSON(w, 200, map[string]interface{}{
		"classifieds": classifieds,
		"total":       30,
	})
}

func handlePostClassified(w http.ResponseWriter, r *http.Request) {
	classified := Classified{
		ID:        rand.Int63n(10000),
		Title:     "New Classified",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}
	respondJSON(w, 201, map[string]interface{}{
		"classified": classified,
	})
}

func handleGetWeather(w http.ResponseWriter, r *http.Request) {
	location := r.URL.Query().Get("location")
	if location == "" {
		location = "San Francisco, CA"
	}

	weather := generateMockWeather(location)
	respondJSON(w, 200, map[string]interface{}{
		"weather": weather,
	})
}

func handleGetProfile(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 200, map[string]interface{}{
		"user": map[string]interface{}{
			"id":            1,
			"username":      "testuser",
			"email":         "test@example.com",
			"karma":         1234,
			"article_count": 42,
			"comment_count": 156,
			"vote_count":    789,
			"created_at":    time.Now().Add(-365 * 24 * time.Hour),
			"last_active":   time.Now(),
		},
	})
}

func handleGetActivity(w http.ResponseWriter, r *http.Request) {
	activity := []map[string]interface{}{
		{
			"type":        "comment",
			"description": "Commented on 'Breaking News Story'",
			"created_at":  time.Now().Add(-1 * time.Hour),
		},
		{
			"type":        "vote",
			"description": "Upvoted 'Tech Article'",
			"created_at":  time.Now().Add(-2 * time.Hour),
		},
	}
	respondJSON(w, 200, map[string]interface{}{
		"activity": activity,
		"total":    2,
	})
}

func generateMockArticles(count int, feed string) []Article {
	articles := make([]Article, count)

	sources := []string{"TechCrunch", "Reuters", "BBC", "HackerNews", "The Verge", "Ars Technica"}
	categories := []string{"tech", "business", "science", "politics", "world", "sports"}
	titles := []string{
		"OpenAI Announces Major Breakthrough in AI Research",
		"Global Markets Rally on Economic Data",
		"New Discovery Could Revolutionize Medicine",
		"Political Leaders Meet for Climate Summit",
		"Tech Giant Unveils Latest Innovation",
		"Breaking: Major Development in Space Exploration",
		"Scientists Make Surprising Finding",
		"Economic Forecast Shows Positive Trends",
		"Technology Startup Raises Record Funding",
		"International Agreement Reached on Key Issue",
	}

	for i := 0; i < count; i++ {
		likeCount := rand.Intn(1000) + 10
		dislikeCount := rand.Intn(100)
		openCount := rand.Intn(5000) + 100
		totalEngagement := likeCount + dislikeCount + openCount/10
		hoursSince := float64(i)

		articles[i] = Article{
			ID:          int64(i + 1),
			Title:       titles[i%len(titles)] + fmt.Sprintf(" (#%d)", i+1),
			URL:         fmt.Sprintf("https://example.com/article/%d", i+1),
			Content:     "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris.",
			ImageURL:    fmt.Sprintf("https://example.com/images/%d.jpg", i+1),
			Source:      sources[i%len(sources)],
			Author:      "John Doe",
			PublishedAt: time.Now().Add(-time.Duration(i) * time.Hour),
			Category:    categories[i%len(categories)],
			Tags:        []string{"tech", "news"},
			// Ranking fields
			OpenCount:           openCount,
			LikeCount:           likeCount,
			DislikeCount:        dislikeCount,
			TotalScore:          likeCount - dislikeCount,
			ControversyScore:    float64(dislikeCount) / float64(likeCount+1),
			TotalEngagement:     totalEngagement,
			HoursSincePublished: hoursSince,
			HotRank:             float64(likeCount) / (hoursSince + 2),
		}
	}

	return articles
}

func generateMockComments(count int) []Comment {
	comments := make([]Comment, count)

	usernames := []string{"alice", "bob", "charlie", "diana", "eve", "frank"}
	contents := []string{
		"Great article! Thanks for sharing.",
		"I disagree with this perspective.",
		"Very interesting point.",
		"Can you provide more sources?",
		"This is exactly what I was looking for.",
		"Excellent analysis.",
	}

	for i := 0; i < count; i++ {
		var parentID *int64
		if i > 0 && rand.Float32() > 0.5 {
			pid := int64(rand.Intn(i))
			parentID = &pid
		}

		comments[i] = Comment{
			ID:        int64(i + 1),
			ArticleID: 1,
			ParentID:  parentID,
			UserID:    int64(rand.Intn(len(usernames)) + 1),
			Content:   contents[i%len(contents)] + fmt.Sprintf(" (Comment #%d)", i+1),
			Upvotes:   rand.Intn(50),
			Downvotes: rand.Intn(10),
			CreatedAt: time.Now().Add(-time.Duration(i) * time.Minute * 10),
			UpdatedAt: time.Now().Add(-time.Duration(i) * time.Minute * 10),
			EditedAt:  nil,
			IsDeleted: false,
			IsFlagged: false,
			FlagCount: 0,
			Username:  usernames[i%len(usernames)],
			Karma:     rand.Intn(1000) + 10,
		}
	}

	return comments
}

func generateMockClassifieds(count int) []Classified {
	classifieds := make([]Classified, count)

	categories := []string{"For Sale", "Jobs", "Housing", "Services", "Events", "Gigs"}
	titles := []string{
		"MacBook Pro M3 - Excellent Condition",
		"Senior Software Engineer Position",
		"2BR Apartment Downtown",
		"Professional Photography Services",
		"Local Startup Networking Event",
		"Freelance Web Development",
	}
	cities := []string{"San Francisco", "New York", "Austin", "Seattle", "Boston", "Denver"}

	for i := 0; i < count; i++ {
		classifieds[i] = Classified{
			ID:           int64(i + 1),
			Title:        titles[i%len(titles)] + fmt.Sprintf(" (#%d)", i+1),
			Description:  "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Contact for more details.",
			Category:     categories[i%len(categories)],
			Price:        float64(rand.Intn(5000) + 100),
			City:         cities[i%len(cities)],
			State:        "CA",
			ContactEmail: "contact@example.com",
			ContactPhone: "(555) 123-4567",
			IsPremium:    rand.Float32() > 0.8,
			CreatedAt:    time.Now().Add(-time.Duration(i) * time.Hour * 24),
			ExpiresAt:    time.Now().Add(time.Duration(30-i) * time.Hour * 24),
		}
	}

	return classifieds
}

func generateMockWeather(location string) Weather {
	conditions := []string{"Clear", "Partly Cloudy", "Cloudy", "Rainy", "Sunny"}
	forecast := make([]ForecastDay, 5)

	for i := 0; i < 5; i++ {
		forecast[i] = ForecastDay{
			Date:      time.Now().Add(time.Duration(i) * 24 * time.Hour).Format("2006-01-02"),
			High:      rand.Intn(30) + 50,
			Low:       rand.Intn(20) + 40,
			Condition: conditions[rand.Intn(len(conditions))],
		}
	}

	return Weather{
		Location: location,
		Current: CurrentWeather{
			Temperature: rand.Intn(30) + 50,
			Condition:   conditions[rand.Intn(len(conditions))],
			Humidity:    rand.Intn(40) + 40,
			WindSpeed:   rand.Intn(20) + 5,
		},
		Forecast:  forecast,
		UpdatedAt: time.Now(),
	}
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
