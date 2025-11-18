# Terminal News - Backend API

Go-based REST API server for Terminal News.

## Status: Foundation Complete ✅

### Implemented:
- ✅ Server setup with Chi router
- ✅ Database connection (PostgreSQL + Redis)
- ✅ User authentication (register/login)
- ✅ JWT token generation and validation
- ✅ Auth middleware
- ✅ All API route stubs

### To Implement:
- ✅ Articles API endpoints
- ✅ Voting system
- ✅ Comments system
- ✅ WebSocket real-time updates
- ✅ Background jobs/scheduler
- ✅ Classifieds CRUD
- ✅ Weather API integration
- ✅ Rate limiting
- ✅ Security middleware
- ⏳ Stripe payments

## Quick Start

### Prerequisites
- Go 1.21+
- PostgreSQL 15+
- Redis 7+

### Setup

```bash
# 1. Install Go dependencies
go mod download

# 2. Copy environment file
cp .env.example .env
# Edit .env with your database credentials

# 3. Start PostgreSQL and Redis
# Using Docker:
docker-compose -f ../docker-compose.dev.yml up -d postgres redis

# 4. Run migrations
psql $DATABASE_URL -f ../database/migrations/001_initial_schema.sql
psql $DATABASE_URL -f ../database/migrations/002_triggers_and_functions.sql

# 5. Run the server
go run cmd/server/main.go
```

Server will start on http://localhost:8080

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register new user ✅
- `POST /api/v1/auth/login` - Login user ✅
- `POST /api/v1/auth/refresh` - Refresh access token ✅

### Articles
- `GET /api/v1/articles` - List articles (with optional ?feed=hot/controversial/rising) ✅
- `GET /api/v1/articles/hot` - Hot feed ✅
- `GET /api/v1/articles/controversial` - Controversial feed ✅
- `GET /api/v1/articles/rising` - Rising feed ✅
- `GET /api/v1/articles/{id}` - Get article ✅

### User (Partially Implemented)
- `GET /api/v1/user/profile` - Get current user profile ✅
- `PUT /api/v1/user/profile` - Update profile (stub)
- `GET /api/v1/user/activity` - Get user activity (stub)

### Voting
- `POST /api/v1/articles/{id}/vote` - Vote on article ✅
- `DELETE /api/v1/articles/{id}/vote` - Remove vote ✅

### Comments
- `GET /api/v1/articles/{id}/comments` - Get comments in tree structure ✅
- `POST /api/v1/articles/{id}/comments` - Post comment (supports nested replies) ✅
- `PUT /api/v1/comments/{id}` - Update comment ✅
- `DELETE /api/v1/comments/{id}` - Delete comment (soft delete) ✅

### Classifieds
- `GET /api/v1/classifieds` - List classifieds (filter by category/city/state) ✅
- `GET /api/v1/classifieds?lat=...&lng=...&radius=...` - Geographic search ✅
- `GET /api/v1/classifieds/{id}` - Get classified ✅
- `POST /api/v1/classifieds` - Create classified ✅
- `PUT /api/v1/classifieds/{id}` - Update classified ✅
- `DELETE /api/v1/classifieds/{id}` - Delete classified ✅

## Testing Auth Endpoints

### Register
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

### Get Profile (requires token)
```bash
curl http://localhost:8080/api/v1/user/profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Get Hot Articles
```bash
curl http://localhost:8080/api/v1/articles/hot?limit=10&offset=0
```

### Get Article by ID
```bash
curl http://localhost:8080/api/v1/articles/1
```

### Vote on Article (requires token)
```bash
# Vote with "like"
curl -X POST http://localhost:8080/api/v1/articles/1/vote \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"vote_type": "like"}'

# Vote with "open" (tracking article opens)
curl -X POST http://localhost:8080/api/v1/articles/1/vote \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"vote_type": "open"}'

# Vote with "dislike"
curl -X POST http://localhost:8080/api/v1/articles/1/vote \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"vote_type": "dislike"}'
```

### Remove Vote (requires token)
```bash
curl -X DELETE "http://localhost:8080/api/v1/articles/1/vote?vote_type=like" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Get Comments for Article
```bash
curl http://localhost:8080/api/v1/articles/1/comments
```

### Post Comment (requires token)
```bash
# Top-level comment
curl -X POST http://localhost:8080/api/v1/articles/1/comments \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content": "Great article!"}'

# Reply to comment
curl -X POST http://localhost:8080/api/v1/articles/1/comments \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content": "I agree!", "parent_id": 5}'
```

### Update Comment (requires token)
```bash
curl -X PUT http://localhost:8080/api/v1/comments/5 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content": "Updated comment text"}'
```

### Delete Comment (requires token)
```bash
curl -X DELETE http://localhost:8080/api/v1/comments/5 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Connect to WebSocket
```bash
# Using wscat (npm install -g wscat)
wscat -c ws://localhost:8080/ws

# With authentication token
wscat -c "ws://localhost:8080/ws?token=YOUR_ACCESS_TOKEN"
```

### Create Classified (requires token)
```bash
curl -X POST http://localhost:8080/api/v1/classifieds \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Vintage bicycle for sale",
    "description": "Classic 1970s Schwinn in excellent condition",
    "price": 250.00,
    "category": "for-sale",
    "subcategory": "bicycles",
    "city": "Portland",
    "state": "OR",
    "country": "US",
    "lat": 45.5152,
    "lng": -122.6784,
    "contact_email": "seller@example.com",
    "contact_method": "email",
    "expires_in_days": 30
  }'
```

### Search Classifieds by Location
```bash
# Find classifieds within 25 miles of coordinates
curl "http://localhost:8080/api/v1/classifieds?lat=45.5152&lng=-122.6784&radius=25"

# Filter by category
curl "http://localhost:8080/api/v1/classifieds?category=for-sale"

# Filter by city
curl "http://localhost:8080/api/v1/classifieds?city=Portland&state=OR"
```

### Get Weather
```bash
# Get weather for Portland, OR
curl "http://localhost:8080/api/v1/weather?lat=45.5152&lng=-122.6784"
```

## Project Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # Entry point
├── internal/
│   ├── api/
│   │   ├── handlers/            # HTTP handlers
│   │   │   ├── auth.go         # ✅ Implemented
│   │   │   ├── articles.go     # ✅ Implemented
│   │   │   ├── votes.go        # ✅ Implemented
│   │   │   ├── comments.go     # ✅ Implemented
│   │   │   ├── classifieds.go  # ✅ Implemented
│   │   │   ├── weather.go      # ✅ Implemented
│   │   │   ├── payments.go     # Stub
│   │   │   └── websocket.go    # ✅ Implemented
│   │   └── middleware/
│   │       └── auth.go          # ✅ Implemented
│   ├── database/
│   │   └── db.go                # ✅ Implemented
│   ├── services/
│   │   ├── auth.go              # ✅ Implemented
│   │   ├── articles.go          # ✅ Implemented
│   │   ├── votes.go             # ✅ Implemented
│   │   ├── comments.go          # ✅ Implemented
│   │   ├── classifieds.go       # ✅ Implemented
│   │   └── payments.go          # Stub
│   ├── external/
│   │   └── weather.go           # ✅ Implemented (NOAA API)
│   └── workers/
│       └── scheduler.go         # ✅ Implemented
├── pkg/
│   └── websocket/
│       ├── hub.go               # ✅ Implemented
│       └── client.go            # ✅ Implemented
├── go.mod
├── go.sum
├── .env.example
└── README.md
```

## Development

### Run with hot reload
```bash
# Install air for hot reload
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
```

### Run tests
```bash
go test ./...
```

### Format code
```bash
go fmt ./...
```

### Lint
```bash
golangci-lint run
```

## Next Steps (Dev 1 Tasks)

### Week 1-2: Core Features ✅ COMPLETE
1. ✅ Implement articles API endpoints
2. ✅ Implement voting system
3. ✅ Add Redis caching for rankings
4. ✅ Implement comments system
5. ✅ Implement WebSocket real-time updates
6. ✅ Implement background scheduler

### Week 3-4: Extended Features ✅ COMPLETE
1. ✅ Implement classifieds CRUD with geographic search
2. ✅ Weather API integration (NOAA)
3. Complete user profile update endpoints
4. Add full-text search
5. Add rate limiting middleware

### Week 5-6: Payments
9. Stripe integration
10. Payment webhooks
11. Subscription management

### Week 7-8: Real-time & Background
12. WebSocket implementation
13. Background job scheduler
14. Auto-refresh rankings
15. Auto-expire classifieds

## Notes for Other Devs

**For Dev 2 (CLI):**
- ✅ Auth endpoints (register, login, refresh token)
- ✅ Articles endpoints (hot/controversial/rising feeds with Redis caching)
- ✅ Voting endpoints (track opens, likes, dislikes)
- ✅ Comments endpoints (create, read, update, delete with tree structure)
- ✅ Classifieds endpoints (full CRUD with geographic search)
- ✅ Weather endpoint (NOAA integration for local weather)
- ✅ WebSocket endpoint for real-time updates
- ✅ User activity endpoint (view comment/vote history)
- Access token expires in 15 minutes, refresh token expires in 7 days
- Include `Authorization: Bearer <token>` header for protected endpoints
- Redis caching: 5min for hot/controversial, 3min for rising, 10min for individual articles
- Background scheduler refreshes rankings every 5 minutes automatically
- Classifieds auto-expire after 30 days (configurable)
- Geographic search uses Haversine formula for radius-based queries

**For Dev 3 (Scraper):**
- Database schema is ready
- Use the same database connection setup
- Articles table structure is in `database/migrations/001_initial_schema.sql`

## Rate Limiting

The API implements Redis-based rate limiting to prevent abuse:

- **Global Rate Limit**: 100 requests per minute per IP address
- **Headers Returned**:
  - `X-RateLimit-Limit`: Maximum requests allowed
  - `X-RateLimit-Remaining`: Requests remaining in current window
  - `X-RateLimit-Reset`: Unix timestamp when limit resets
  - `Retry-After`: Seconds to wait before retrying (when limited)

When rate limit is exceeded, the API returns:
- Status: `429 Too Many Requests`
- Body: `{"error": "Rate limit exceeded. Please try again later."}`

Rate limiting uses Redis for distributed tracking, so it works across multiple server instances.

## Security Features

The API includes multiple security middleware layers:

- **Security Headers**:
  - `X-Content-Type-Options: nosniff` - Prevent MIME sniffing
  - `X-Frame-Options: DENY` - Prevent clickjacking
  - `X-XSS-Protection: 1; mode=block` - XSS protection
  - `Referrer-Policy: strict-origin-when-cross-origin`
  - `Content-Security-Policy: default-src 'self'`
  - `Permissions-Policy` - Restrict browser features

- **Panic Recovery**: Catches panics and returns 500 errors
- **Request Timeouts**: 60-second timeout on all requests
- **CORS**: Configurable cross-origin resource sharing
- **JWT Authentication**: Secure token-based auth with 15-min expiry

## Environment Variables

See `.env.example` for all required environment variables.

Critical ones:
- `DATABASE_URL` - PostgreSQL connection string
- `REDIS_URL` - Redis connection string
- `JWT_SECRET` - Secret key for JWT signing (CHANGE IN PRODUCTION!)

## Troubleshooting

### "connect: connection refused"
- Make sure PostgreSQL and Redis are running
- Check DATABASE_URL and REDIS_URL in .env

### "database does not exist"
- Create database: `createdb terminalnews_dev`
- Run migrations

### "JWT_SECRET not set"
- Copy .env.example to .env
- Set JWT_SECRET to a random string

---

**Status:** Week 1-4 complete! Full featured backend ready for CLI integration
**Implemented:** Auth, Articles, Voting, Comments, Classifieds (with geo-search), Weather (NOAA), WebSocket, Background Jobs
**Next:** Stripe payments, rate limiting, full-text search
**Note:** Weather, classifieds, and local news are all location-aware using lat/lng coordinates
