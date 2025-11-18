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
- ⏳ Comments system
- ⏳ Classifieds CRUD
- ⏳ Weather API integration
- ⏳ Stripe payments
- ⏳ WebSocket real-time updates
- ⏳ Background jobs/scheduler

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

### Comments (Not Yet Implemented)
- `GET /api/v1/articles/{id}/comments` - Get comments
- `POST /api/v1/articles/{id}/comments` - Post comment
- `PUT /api/v1/comments/{id}` - Update comment
- `DELETE /api/v1/comments/{id}` - Delete comment

### Classifieds (Not Yet Implemented)
- `GET /api/v1/classifieds` - List classifieds
- `GET /api/v1/classifieds/{id}` - Get classified
- `POST /api/v1/classifieds` - Create classified
- `PUT /api/v1/classifieds/{id}` - Update classified
- `DELETE /api/v1/classifieds/{id}` - Delete classified

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
│   │   │   ├── comments.go     # Stub
│   │   │   ├── classifieds.go  # Stub
│   │   │   ├── weather.go      # Stub
│   │   │   ├── payments.go     # Stub
│   │   │   └── websocket.go    # Stub
│   │   └── middleware/
│   │       └── auth.go          # ✅ Implemented
│   ├── database/
│   │   └── db.go                # ✅ Implemented
│   └── services/
│       ├── auth.go              # ✅ Implemented
│       ├── articles.go          # ✅ Implemented
│       ├── votes.go             # ✅ Implemented
│       ├── comments.go          # Stub
│       ├── classifieds.go       # Stub
│       ├── payments.go          # Stub
│       └── scheduler.go         # Stub
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

### Week 1-2: Core Features
1. ✅ Implement articles API endpoints
2. ✅ Implement voting system
3. ✅ Add Redis caching for rankings
4. Complete user profile endpoints

### Week 3-4: Extended Features
5. Implement comments system
6. Implement classifieds CRUD
7. Add full-text search
8. Implement ranking algorithms

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
- Auth endpoints are READY to use ✅
- Articles endpoints are READY to use ✅
- Voting endpoints are READY to use ✅
- Test credentials can be created via register endpoint
- Access token expires in 15 minutes
- Refresh token expires in 7 days
- Include `Authorization: Bearer <token>` header for protected endpoints
- All article feeds include Redis caching (5min for hot/controversial, 3min for rising)
- Voting automatically invalidates article caches for real-time ranking updates

**For Dev 3 (Scraper):**
- Database schema is ready
- Use the same database connection setup
- Articles table structure is in `database/migrations/001_initial_schema.sql`

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

**Status:** Core features implemented (Auth + Articles + Voting)
**Next:** Implement comments system and classifieds CRUD
