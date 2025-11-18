# Developer 1: Backend API & Infrastructure Guide

## Your Mission
Build the entire backend API server, database layer, payment systems, and infrastructure that powers Terminal News.

---

## 🚀 IMMEDIATE SETUP (Day 1)

### 1. Clone and Navigate
```bash
cd C:\Users\wesle\Dropbox\GITHUB\02 - Applications\terminal-news
mkdir -p backend
cd backend
```

### 2. Initialize Go Module
```bash
go mod init github.com/yourusername/terminal-news/backend
```

### 3. Install Dependencies
```bash
# Core dependencies
go get -u github.com/go-chi/chi/v5
go get -u github.com/go-chi/chi/v5/middleware
go get -u github.com/go-chi/cors
go get -u github.com/golang-jwt/jwt/v5
go get -u github.com/lib/pq
go get -u github.com/jmoiron/sqlx
go get -u github.com/go-redis/redis/v8
go get -u github.com/gorilla/websocket
go get -u github.com/stripe/stripe-go/v74
go get -u golang.org/x/crypto/bcrypt
go get -u github.com/joho/godotenv
go get -u github.com/golang-migrate/migrate/v4
go get -u github.com/golang-migrate/migrate/v4/database/postgres
go get -u github.com/golang-migrate/migrate/v4/source/file
```

### 4. Create Initial Structure
```bash
# From terminal-news/backend directory
mkdir -p cmd/server
mkdir -p internal/api/handlers
mkdir -p internal/api/middleware
mkdir -p internal/database
mkdir -p internal/models
mkdir -p internal/services
mkdir -p internal/utils
mkdir -p pkg/websocket
mkdir -p migrations
mkdir -p scripts
```

### 5. Database Setup
```bash
# Install PostgreSQL locally or use Docker
docker run --name terminal-news-postgres \
  -e POSTGRES_DB=terminalnews \
  -e POSTGRES_USER=tnuser \
  -e POSTGRES_PASSWORD=tnpass123 \
  -p 5432:5432 \
  -d postgres:15

# Redis
docker run --name terminal-news-redis \
  -p 6379:6379 \
  -d redis:7-alpine
```

### 6. Environment File
Create `C:\Users\wesle\Dropbox\GITHUB\02 - Applications\terminal-news\backend\.env`:
```env
DATABASE_URL=postgres://tnuser:tnpass123@localhost:5432/terminalnews?sslmode=disable
REDIS_URL=redis://localhost:6379
JWT_SECRET=your-super-secret-key-change-this
STRIPE_SECRET_KEY=sk_test_xxxxx
STRIPE_WEBHOOK_SECRET=whsec_xxxxx
PORT=8080
ENVIRONMENT=development
```

---

## 📁 FILE REFERENCES

### Database Schema
**USE THIS**: `C:\Users\wesle\Dropbox\GITHUB\02 - Applications\terminal-news\docs\DATABASE_SCHEMA.md`
- Contains all tables, indexes, functions
- Copy SQL from here for migrations

### API Architecture
**READ THIS**: `C:\Users\wesle\Dropbox\GITHUB\02 - Applications\terminal-news\docs\ARCHITECTURE.md`
- System design details
- API endpoint specifications

### Business Logic
**REFERENCE**: `C:\Users\wesle\Dropbox\GITHUB\02 - Applications\terminal-news\docs\BUSINESS_MODEL.md`
- Pricing tiers
- Payment flow

---

## 🔨 WEEK 1-2: Core API Foundation

### Create Main Server
`C:\Users\wesle\Dropbox\GITHUB\02 - Applications\terminal-news\backend\cmd\server\main.go`:
```go
package main

import (
    "log"
    "net/http"
    "os"
    
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/go-chi/cors"
    "github.com/joho/godotenv"
    _ "github.com/lib/pq"
)

func main() {
    // Load .env
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }
    
    // Initialize database
    db := initDB()
    defer db.Close()
    
    // Initialize Redis
    rdb := initRedis()
    defer rdb.Close()
    
    // Setup router
    r := chi.NewRouter()
    
    // Middleware
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(middleware.RequestID)
    r.Use(cors.Handler(cors.Options{
        AllowedOrigins:   []string{"*"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
        ExposedHeaders:   []string{"Link"},
        AllowCredentials: true,
        MaxAge:           300,
    }))
    
    // Routes
    r.Route("/api", func(r chi.Router) {
        // Auth
        r.Post("/auth/register", handleRegister)
        r.Post("/auth/login", handleLogin)
        r.Post("/auth/logout", handleLogout)
        
        // Articles
        r.Get("/articles", handleGetArticles)
        r.Get("/articles/{id}", handleGetArticle)
        
        // Votes
        r.Post("/articles/{id}/vote", handleVote)
        
        // Comments
        r.Get("/articles/{id}/comments", handleGetComments)
        r.Post("/articles/{id}/comments", handlePostComment)
        
        // Classifieds
        r.Get("/classifieds", handleGetClassifieds)
        r.Post("/classifieds", handlePostClassified)
        r.Get("/classifieds/{id}", handleGetClassified)
        r.Put("/classifieds/{id}", handleUpdateClassified)
        r.Delete("/classifieds/{id}", handleDeleteClassified)
        
        // Weather
        r.Get("/weather", handleGetWeather)
        
        // User
        r.Get("/user/profile", handleGetProfile)
        r.Get("/user/activity", handleGetActivity)
    })
    
    // WebSocket
    r.Get("/ws", handleWebSocket)
    
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    log.Printf("Server starting on port %s", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}
```

### Database Connection
`C:\Users\wesle\Dropbox\GITHUB\02 - Applications\terminal-news\backend\internal\database\db.go`:
```go
package database

import (
    "database/sql"
    "log"
    "os"
    
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

var DB *sqlx.DB

func InitDB() *sqlx.DB {
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        log.Fatal("DATABASE_URL not set")
    }
    
    db, err := sqlx.Connect("postgres", dsn)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    
    // Set connection pool settings
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    
    DB = db
    return db
}
```

### User Model
`C:\Users\wesle\Dropbox\GITHUB\02 - Applications\terminal-news\backend\internal\models\user.go`:
```go
package models

import (
    "time"
    "database/sql"
)

type User struct {
    ID           int64          `db:"id" json:"id"`
    Username     string         `db:"username" json:"username"`
    Email        string         `db:"email" json:"email"`
    PasswordHash string         `db:"password_hash" json:"-"`
    DisplayName  sql.NullString `db:"display_name" json:"display_name,omitempty"`
    Bio          sql.NullString `db:"bio" json:"bio,omitempty"`
    Location     sql.NullString `db:"location" json:"location,omitempty"`
    Website      sql.NullString `db:"website" json:"website,omitempty"`
    Karma        int            `db:"karma" json:"karma"`
    TrustScore   float64        `db:"trust_score" json:"trust_score"`
    IsBanned     bool           `db:"is_banned" json:"-"`
    IsModerator  bool           `db:"is_moderator" json:"is_moderator"`
    IsAdmin      bool           `db:"is_admin" json:"is_admin"`
    CreatedAt    time.Time      `db:"created_at" json:"created_at"`
    UpdatedAt    time.Time      `db:"updated_at" json:"updated_at"`
}
```

### Authentication Handler
`C:\Users\wesle\Dropbox\GITHUB\02 - Applications\terminal-news\backend\internal\api\handlers\auth.go`:
```go
package handlers

import (
    "encoding/json"
    "net/http"
    "time"
    
    "golang.org/x/crypto/bcrypt"
    "github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type RegisterRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
    var req RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Hash password
    hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Failed to hash password", http.StatusInternalServerError)
        return
    }
    
    // Insert user into database
    query := `
        INSERT INTO users (username, email, password_hash)
        VALUES ($1, $2, $3)
        RETURNING id, username, email, created_at
    `
    
    // Execute query and return user
    // ... implementation
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
    var req LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Fetch user from database
    // Verify password
    // Generate JWT token
    // Return token
}
```

### Create Migrations
`C:\Users\wesle\Dropbox\GITHUB\02 - Applications\terminal-news\backend\migrations\001_initial_schema.up.sql`:
```sql
-- Copy the entire schema from DATABASE_SCHEMA.md
-- Start with users table
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    -- ... rest of schema
);

-- Continue with all other tables from DATABASE_SCHEMA.md
```

---

## 🔨 WEEK 3-4: Advanced Features

### Ranking System Implementation
`C:\Users\wesle\Dropbox\GITHUB\02 - Applications\terminal-news\backend\internal\services\ranking.go`:
```go
package services

import (
    "context"
    "log"
    "time"
)

// Implement the ranking algorithm from DATABASE_SCHEMA.md
func RefreshRankings(ctx context.Context) error {
    query := `REFRESH MATERIALIZED VIEW CONCURRENTLY article_rankings`
    // Execute every 5 minutes
    return nil
}

func CalculateHotScore(likes, dislikes, opens int, publishedAt time.Time) float64 {
    // Implement hot ranking algorithm
    score := float64(opens*1 + likes*2 + dislikes*-1)
    hoursAgo := time.Since(publishedAt).Hours()
    timeDecay := 1.0 / math.Pow(hoursAgo+2, 1.5)
    return score * timeDecay
}
```

### Redis Caching
`C:\Users\wesle\Dropbox\GITHUB\02 - Applications\terminal-news\backend\internal\services\cache.go`:
```go
package services

import (
    "context"
    "encoding/json"
    "time"
    
    "github.com/go-redis/redis/v8"
)

type CacheService struct {
    client *redis.Client
}

func NewCacheService(client *redis.Client) *CacheService {
    return &CacheService{client: client}
}

func (c *CacheService) GetArticleRankings(ctx context.Context, feed string) ([]Article, error) {
    key := fmt.Sprintf("rankings:%s", feed)
    // Get from Redis
    // If miss, get from DB and cache
}

func (c *CacheService) InvalidateArticleRankings(ctx context.Context) {
    // Clear all ranking caches
}
```

### WebSocket Real-time Updates
`C:\Users\wesle\Dropbox\GITHUB\02 - Applications\terminal-news\backend\pkg\websocket\hub.go`:
```go
package websocket

import (
    "log"
    "github.com/gorilla/websocket"
)

type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
}

func NewHub() *Hub {
    return &Hub{
        broadcast:  make(chan []byte),
        register:   make(chan *Client),
        unregister: make(chan *Client),
        clients:    make(map[*Client]bool),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.clients[client] = true
            
        case client := <-h.unregister:
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
            }
            
        case message := <-h.broadcast:
            for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(h.clients, client)
                }
            }
        }
    }
}
```

---

## 🔨 WEEK 5-6: Payment Integration

### Stripe Setup
`C:\Users\wesle\Dropbox\GITHUB\02 - Applications\terminal-news\backend\internal\services\payments.go`:
```go
package services

import (
    "github.com/stripe/stripe-go/v74"
    "github.com/stripe/stripe-go/v74/customer"
    "github.com/stripe/stripe-go/v74/paymentintent"
    "github.com/stripe/stripe-go/v74/subscription"
)

func InitStripe() {
    stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
}

func CreateCustomer(email, name string) (*stripe.Customer, error) {
    params := &stripe.CustomerParams{
        Email: stripe.String(email),
        Name:  stripe.String(name),
    }
    return customer.New(params)
}

func CreateClassifiedPayment(amount int64, classifiedID int64) (*stripe.PaymentIntent, error) {
    params := &stripe.PaymentIntentParams{
        Amount:   stripe.Int64(amount),
        Currency: stripe.String(string(stripe.CurrencyUSD)),
        Metadata: map[string]string{
            "classified_id": fmt.Sprintf("%d", classifiedID),
            "type":         "classified_premium",
        },
    }
    return paymentintent.New(params)
}

func CreateSponsorSubscription(customerID string, priceID string) (*stripe.Subscription, error) {
    // Create recurring subscription for sponsors
}
```

### Webhook Handler
`C:\Users\wesle\Dropbox\GITHUB\02 - Applications\terminal-news\backend\internal\api\handlers\webhook.go`:
```go
package handlers

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
    
    "github.com/stripe/stripe-go/v74/webhook"
)

func HandleStripeWebhook(w http.ResponseWriter, r *http.Request) {
    const MaxBodyBytes = int64(65536)
    r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
    
    payload, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusServiceUnavailable)
        return
    }
    
    endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
    event, err := webhook.ConstructEvent(payload, 
        r.Header.Get("Stripe-Signature"), endpointSecret)
    
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    switch event.Type {
    case "payment_intent.succeeded":
        // Handle successful payment
        // Update classified to premium
        
    case "customer.subscription.created":
        // Handle new subscription
        // Activate sponsor
        
    case "customer.subscription.deleted":
        // Handle cancelled subscription
        // Deactivate sponsor
    }
    
    w.WriteHeader(http.StatusOK)
}
```

---

## 🔨 WEEK 7-8: Background Jobs & Automation

### Background Worker
`C:\Users\wesle\Dropbox\GITHUB\02 - Applications\terminal-news\backend\internal\workers\scheduler.go`:
```go
package workers

import (
    "context"
    "log"
    "time"
)

type Scheduler struct {
    db    *sqlx.DB
    redis *redis.Client
}

func (s *Scheduler) Start(ctx context.Context) {
    // Refresh rankings every 5 minutes
    go s.runEvery(ctx, 5*time.Minute, s.RefreshRankings)
    
    // Expire old classifieds every hour
    go s.runEvery(ctx, 1*time.Hour, s.ExpireClassifieds)
    
    // Clean up old audit logs daily
    go s.runEvery(ctx, 24*time.Hour, s.CleanupAuditLogs)
    
    // Send billing reminders daily at 9am
    go s.runDaily(ctx, 9, 0, s.SendBillingReminders)
}

func (s *Scheduler) runEvery(ctx context.Context, d time.Duration, f func() error) {
    ticker := time.NewTicker(d)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            if err := f(); err != nil {
                log.Printf("Scheduler error: %v", err)
            }
        }
    }
}

func (s *Scheduler) RefreshRankings() error {
    _, err := s.db.Exec(`REFRESH MATERIALIZED VIEW CONCURRENTLY article_rankings`)
    return err
}

func (s *Scheduler) ExpireClassifieds() error {
    _, err := s.db.Exec(`
        UPDATE classifieds 
        SET is_active = FALSE 
        WHERE expires_at < NOW() AND is_active = TRUE
    `)
    return err
}
```

### API Rate Limiting
`C:\Users\wesle\Dropbox\GITHUB\02 - Applications\terminal-news\backend\internal\api\middleware\ratelimit.go`:
```go
package middleware

import (
    "net/http"
    "time"
    
    "github.com/go-redis/redis/v8"
)

func RateLimit(rdb *redis.Client, requests int, window time.Duration) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Get client IP
            ip := r.RemoteAddr
            
            // Check rate limit in Redis
            key := fmt.Sprintf("ratelimit:%s", ip)
            
            // Increment counter
            pipe := rdb.Pipeline()
            incr := pipe.Incr(r.Context(), key)
            pipe.Expire(r.Context(), key, window)
            _, err := pipe.Exec(r.Context())
            
            if err != nil {
                http.Error(w, "Internal error", http.StatusInternalServerError)
                return
            }
            
            if incr.Val() > int64(requests) {
                http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}
```

---

## 🔨 WEEK 9: Testing & Deployment

### Integration Tests
`C:\Users\wesle\Dropbox\GITHUB\02 - Applications\terminal-news\backend\internal\api\handlers\auth_test.go`:
```go
package handlers_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestRegisterUser(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer db.Close()
    
    // Create request
    user := RegisterRequest{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "secure123",
    }
    
    body, _ := json.Marshal(user)
    req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
    w := httptest.NewRecorder()
    
    // Call handler
    HandleRegister(w, req)
    
    // Check response
    if w.Code != http.StatusCreated {
        t.Errorf("Expected status 201, got %d", w.Code)
    }
}
```

### Docker Configuration
`C:\Users\wesle\Dropbox\GITHUB\02 - Applications\terminal-news\backend\Dockerfile`:
```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/server .
COPY --from=builder /app/.env .

EXPOSE 8080
CMD ["./server"]
```

### GitHub Actions CI/CD
`C:\Users\wesle\Dropbox\GITHUB\02 - Applications\terminal-news\.github\workflows\backend.yml`:
```yaml
name: Backend CI/CD

on:
  push:
    branches: [main]
    paths:
      - 'backend/**'
  pull_request:
    branches: [main]
    paths:
      - 'backend/**'

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: postgres
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432
          
      redis:
        image: redis:7-alpine
        ports:
          - 6379:6379
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    
    - name: Install dependencies
      working-directory: ./backend
      run: go mod download
    
    - name: Run migrations
      working-directory: ./backend
      run: |
        go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
        migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/test?sslmode=disable" up
    
    - name: Run tests
      working-directory: ./backend
      run: go test -v ./...
    
    - name: Build
      working-directory: ./backend
      run: go build -v ./cmd/server
```

---

## 📋 DAILY CHECKLIST

### Week 1-2 Goals
- [ ] Database connected and migrations running
- [ ] User registration/login working
- [ ] Basic CRUD for articles, votes, comments
- [ ] JWT authentication implemented
- [ ] Redis caching layer setup
- [ ] WebSocket server running

### Week 3-4 Goals
- [ ] Ranking algorithms implemented
- [ ] Materialized views refreshing
- [ ] Full classifieds CRUD
- [ ] Geographic search working
- [ ] Rate limiting active
- [ ] Security middleware complete

### Week 5-6 Goals
- [ ] Stripe integration complete
- [ ] Payment webhooks handling
- [ ] Subscription management
- [ ] Premium features logic
- [ ] Automated billing system
- [ ] Background jobs running

### Week 7-8 Goals
- [ ] All cron jobs scheduled
- [ ] Performance optimized
- [ ] Monitoring setup
- [ ] API documentation complete
- [ ] Load testing passed
- [ ] Docker deployment ready

### Week 9 Goals
- [ ] All tests passing
- [ ] CI/CD pipeline working
- [ ] Production deployment script
- [ ] Performance benchmarks met
- [ ] Security audit complete
- [ ] Ready for launch

---

## 🔗 RESOURCES

### GitHub References
- [Chi Router Examples](https://github.com/go-chi/chi/tree/master/_examples)
- [SQLX Tutorial](https://github.com/jmoiron/sqlx)
- [Stripe Go Examples](https://github.com/stripe/stripe-go/tree/master/_examples)
- [JWT Go Examples](https://github.com/golang-jwt/jwt/tree/main/cmd/jwt)
- [WebSocket Chat Example](https://github.com/gorilla/websocket/tree/master/examples/chat)

### Documentation
- [PostgreSQL Docs](https://www.postgresql.org/docs/15/)
- [Redis Commands](https://redis.io/commands/)
- [Stripe API Reference](https://stripe.com/docs/api)
- [Go Best Practices](https://go.dev/doc/effective_go)

### Testing Tools
- [Postman](https://www.postman.com/) - API testing
- [pgAdmin](https://www.pgadmin.org/) - Database management
- [Redis Insight](https://redis.com/redis-enterprise/redis-insight/) - Redis GUI
- [k6](https://k6.io/) - Load testing

---

## 🚨 CRITICAL PATHS

Your work blocks the other developers after Week 2. Priority order:

1. **Auth endpoints** (Dev 2 needs this for login screen)
2. **Articles GET endpoints** (Dev 2 needs for displaying news)
3. **WebSocket setup** (Dev 2 needs for real-time)
4. **Votes/Comments POST** (Core functionality)
5. **Classifieds CRUD** (Major feature)
6. **Payment system** (Revenue generation)

---

## 💬 COMMUNICATION

### Daily Standup Questions
1. What endpoints did I complete yesterday?
2. What endpoints am I working on today?
3. What do the other devs need from me?

### Share with Team
- API documentation (Swagger/OpenAPI)
- Database connection details
- Test user credentials
- WebSocket event types
- Webhook URLs for testing

---

## 🎯 SUCCESS METRICS

By end of Week 9, you should have:
- ✅ 40+ API endpoints implemented
- ✅ 90%+ test coverage
- ✅ <100ms average response time
- ✅ Support for 1000+ concurrent users
- ✅ Zero security vulnerabilities
- ✅ Complete API documentation
- ✅ Automated deployment pipeline
- ✅ Production-ready infrastructure

Ready to build the backbone of Terminal News! 🚀
