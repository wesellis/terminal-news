# Technical Architecture

## System Overview

Terminal News is a distributed system with a native terminal client connecting to a centralized backend API for shared data and community features.

```
┌─────────────────────────────────────────────────────────┐
│                   Terminal Client                       │
│  (Native CLI - Rust/Go)                                │
│  - Terminal UI rendering                                │
│  - User input handling                                  │
│  - Local caching                                        │
│  - Offline mode                                         │
└────────────────┬────────────────────────────────────────┘
                 │
                 │ HTTPS/WebSocket
                 ▼
┌─────────────────────────────────────────────────────────┐
│                   API Server                            │
│  (Go/Node.js)                                          │
│  - REST API endpoints                                   │
│  - WebSocket server (real-time updates)                │
│  - Authentication & session management                  │
│  - Rate limiting                                        │
└────────────┬────────────────────────┬───────────────────┘
             │                        │
             ▼                        ▼
┌────────────────────┐    ┌──────────────────────┐
│   PostgreSQL       │    │      Redis           │
│                    │    │                      │
│ - Users            │    │ - Ranking cache      │
│ - Articles         │    │ - Session store      │
│ - Votes            │    │ - Rate limiting      │
│ - Comments         │    │ - Real-time scores   │
│ - Classifieds      │    └──────────────────────┘
└────────────────────┘
             │
             ▼
┌─────────────────────────────────────────────────────────┐
│               External Services                         │
│                                                         │
│  - News APIs (NewsAPI, Guardian, NYT)                  │
│  - RSS Feed Aggregator                                  │
│  - NOAA Weather API                                     │
└─────────────────────────────────────────────────────────┘
```

## Client Architecture

### Technology Choice: Go + Bubbletea

**Why Go:**
- Single binary compilation (easy distribution)
- Excellent cross-platform support
- Strong standard library
- Fast execution
- Low memory footprint

**UI Framework: Bubbletea**
- Modern, composable TUI framework
- Elm-inspired architecture (model-update-view)
- Active development and community
- Great for complex interactive UIs

### Client Components

```
terminal-news/
├── cmd/
│   └── terminal-news/
│       └── main.go              # Entry point
├── internal/
│   ├── ui/
│   │   ├── app.go              # Main application model
│   │   ├── tabs/
│   │   │   ├── hot.go          # Hot news view
│   │   │   ├── controversial.go
│   │   │   ├── rising.go
│   │   │   ├── profile.go
│   │   │   ├── classifieds.go
│   │   │   └── weather.go
│   │   └── components/
│   │       ├── article_list.go
│   │       ├── comment_view.go
│   │       └── input.go
│   ├── api/
│   │   ├── client.go           # API client
│   │   └── websocket.go        # Real-time connection
│   ├── cache/
│   │   └── store.go            # Local cache (SQLite)
│   └── config/
│       └── config.go           # User settings
└── pkg/
    └── models/
        └── types.go            # Shared data models
```

### Client Features

**Offline Support**
- SQLite cache for read articles
- Queue actions (votes, comments) for sync
- Graceful degradation

**Keyboard Navigation**
```
Tab/Shift+Tab  - Switch between sections
↑/↓            - Navigate items
Enter          - Open article/thread
L              - Like
D              - Dislike
C              - Comment
R              - Refresh
Q              - Quit
```

## Backend Architecture

### Technology Choice: Go

**API Server Stack:**
- Go (net/http + gorilla/mux or chi)
- PostgreSQL (primary database)
- Redis (caching + real-time)
- Docker (deployment)

### Database Schema

```sql
-- Users
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    location VARCHAR(100),  -- For local news/weather
    created_at TIMESTAMP DEFAULT NOW()
);

-- Articles
CREATE TABLE articles (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    source VARCHAR(100),
    content TEXT,
    published_at TIMESTAMP,
    fetched_at TIMESTAMP DEFAULT NOW(),
    category VARCHAR(50)
);

-- Votes
CREATE TABLE votes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    article_id INTEGER REFERENCES articles(id),
    vote_type VARCHAR(10) CHECK (vote_type IN ('open', 'like', 'dislike')),
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, article_id, vote_type)
);

-- Comments
CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    article_id INTEGER REFERENCES articles(id),
    parent_id INTEGER REFERENCES comments(id),
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Classifieds
CREATE TABLE classifieds (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    title VARCHAR(200) NOT NULL,
    description TEXT,
    category VARCHAR(50),
    location VARCHAR(100),
    price DECIMAL(10,2),
    is_premium BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP
);
```

### Ranking Algorithm

**Score Calculation:**
```
score = (opens * 1) + (likes * 2) + (dislikes * -1)
time_decay = 1 / (hours_since_published + 2)^1.5
final_rank = score * time_decay
```

**Controversy Score:**
```
controversy = min(likes, dislikes) / max(likes, dislikes)
engagement = likes + dislikes + opens
controversial_rank = controversy * engagement
```

**Rising Score:**
```
recent_engagement = votes_last_hour / total_votes
velocity = (score_now - score_1hr_ago) / 1hr
rising_rank = velocity * recent_engagement
```

### API Endpoints

```
Authentication:
POST   /api/auth/register
POST   /api/auth/login
POST   /api/auth/logout

Articles:
GET    /api/articles?feed=hot|controversial|rising&offset=0&limit=50
GET    /api/articles/:id
POST   /api/articles/:id/vote     { type: "open|like|dislike" }

Comments:
GET    /api/articles/:id/comments
POST   /api/articles/:id/comments  { content: "...", parent_id: null }

Classifieds:
GET    /api/classifieds?location=...&category=...
POST   /api/classifieds
PUT    /api/classifieds/:id
DELETE /api/classifieds/:id

User:
GET    /api/user/profile
GET    /api/user/activity
GET    /api/user/classifieds

Weather:
GET    /api/weather?location=...

WebSocket:
WS     /api/ws  (real-time updates)
```

### News Aggregation Service

**Background Worker:**
- Runs every 5-15 minutes
- Fetches from RSS feeds + News APIs
- Deduplicates articles (title similarity)
- Categorizes content
- Stores in database

**Sources:**
- NewsAPI.org (70+ sources)
- RSS feeds (CNN, BBC, Reuters, etc.)
- Guardian API
- NYT API (if budget allows)
- Reddit (r/news, r/worldnews via API)

## Deployment Architecture

### Digital Ocean Setup

```
Droplet Configuration:
- Ubuntu 22.04 LTS
- 2GB RAM / 1 vCPU (start)
- 50GB SSD

Services:
- Docker + Docker Compose
- Nginx (reverse proxy + SSL)
- PostgreSQL (containerized)
- Redis (containerized)
- Go API server (containerized)
```

### Docker Compose Structure

```yaml
version: '3.8'
services:
  api:
    build: ./api
    ports:
      - "8080:8080"
    depends_on:
      - db
      - redis
    environment:
      - DATABASE_URL=postgres://...
      - REDIS_URL=redis://redis:6379

  db:
    image: postgres:15
    volumes:
      - postgres_data:/var/lib/postgresql/data
    environment:
      - POSTGRES_DB=terminalnews
      - POSTGRES_PASSWORD=...

  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - /etc/letsencrypt:/etc/letsencrypt

volumes:
  postgres_data:
  redis_data:
```

## Security Considerations

**Authentication:**
- JWT tokens (short-lived access + refresh tokens)
- Bcrypt password hashing
- Rate limiting on auth endpoints

**API Security:**
- HTTPS only
- CORS configuration
- Input validation & sanitization
- SQL injection prevention (prepared statements)
- XSS prevention (even in terminal context)

**User Privacy:**
- Minimal data collection
- Optional analytics
- No tracking cookies
- Location data stored locally only

## Scalability Plan

**Phase 1 (0-1k users):**
- Single droplet
- Minimal caching

**Phase 2 (1k-10k users):**
- Upgrade droplet (4GB RAM)
- Redis caching layer
- CDN for API responses (Cloudflare)

**Phase 3 (10k+ users):**
- Separate database server
- Read replicas
- Load balancer
- Multi-region deployment

## Performance Targets

- Article list load: <200ms
- Vote registration: <100ms
- Comment post: <300ms
- WebSocket latency: <50ms
- Client startup: <500ms

## Monitoring & Observability

**Metrics:**
- Prometheus + Grafana
- API response times
- Database query performance
- Active users
- Vote/comment rates

**Logging:**
- Structured logging (JSON)
- Log aggregation (Loki or ELK)
- Error tracking (Sentry)

**Alerts:**
- API downtime
- Database connection issues
- High error rates
- Disk space warnings
