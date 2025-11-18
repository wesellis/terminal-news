# Terminal News - Development Setup Guide

Quick start guide for developers to get up and running.

---

## Prerequisites

Before you begin, install:

- **Go 1.21+** - [Download](https://golang.org/dl/)
- **Docker & Docker Compose** - [Download](https://www.docker.com/get-started)
- **Git** - [Download](https://git-scm.com/downloads)
- **PostgreSQL client** (psql) - For running migrations manually
- **Make** (optional) - For using Makefile commands

---

## Quick Start (Automated)

### 1. Clone the repository

```bash
git clone https://github.com/yourusername/terminal-news.git
cd terminal-news
```

### 2. Run setup script

```bash
./scripts/dev-setup.sh
```

This script will:
- Check prerequisites
- Create .env file
- Start Docker services (PostgreSQL, Redis)
- Run database migrations
- Install Go dependencies

### 3. Add API keys

Edit `.env` and add your API keys:

```bash
# News APIs
NEWS_API_KEY=your_key_here          # Get from newsapi.org
GUARDIAN_API_KEY=your_key_here      # Get from open-platform.theguardian.com

# Payment (for testing)
STRIPE_SECRET_KEY=sk_test_...       # Get from stripe.com
```

### 4. Start services

```bash
# Option 1: Use Docker Compose (all services)
docker-compose -f docker-compose.dev.yml up

# Option 2: Run individually
make dev-api      # API server on :8080
make dev-scraper  # News scraper
make dev-cli      # CLI client
```

### 5. Test it works

```bash
# Check API health
curl http://localhost:8080/health

# Run CLI
cd cli
go run cmd/main.go
```

---

## Manual Setup

### 1. Start Database Services

```bash
docker-compose -f docker-compose.dev.yml up -d postgres redis
```

### 2. Run Migrations

```bash
./scripts/run-migrations.sh up postgres://postgres:postgres@localhost:5432/terminalnews_dev?sslmode=disable
```

Or manually:

```bash
psql postgres://postgres:postgres@localhost:5432/terminalnews_dev?sslmode=disable \
  -f database/migrations/001_initial_schema.sql

psql postgres://postgres:postgres@localhost:5432/terminalnews_dev?sslmode=disable \
  -f database/migrations/002_triggers_and_functions.sql
```

### 3. Install Dependencies

```bash
# Backend
cd backend && go mod download && cd ..

# Scraper
cd scraper && go mod download && cd ..

# CLI
cd cli && go mod download && cd ..
```

### 4. Configure Environment

```bash
cp .env.example .env
# Edit .env with your values
```

### 5. Run Services

```bash
# Terminal 1: API Server
cd backend
go run cmd/api/main.go

# Terminal 2: Scraper
cd scraper
go run cmd/main.go

# Terminal 3: CLI
cd cli
go run cmd/main.go
```

---

## Project Structure

```
terminal-news/
├── backend/              # API Server
│   ├── cmd/api/         # Main entry point
│   ├── internal/        # Internal packages
│   │   ├── api/        # HTTP handlers
│   │   ├── auth/       # Authentication
│   │   ├── database/   # DB connections
│   │   ├── middleware/ # HTTP middleware
│   │   ├── models/     # Data models (aliases shared/)
│   │   └── services/   # Business logic
│   └── pkg/            # Public packages
│
├── cli/                 # Terminal Client
│   ├── cmd/            # Main entry point
│   ├── internal/       # Internal packages
│   │   ├── ui/        # Bubbletea UI components
│   │   ├── api/       # API client
│   │   ├── config/    # Configuration
│   │   └── cache/     # Local caching
│   └── pkg/           # Public packages
│
├── scraper/            # News Aggregation Service
│   ├── cmd/           # Main entry point
│   ├── internal/      # Internal packages
│   │   ├── fetchers/ # RSS/API fetchers
│   │   ├── parser/   # Content parsing
│   │   └── storage/  # Database storage
│   └── pkg/          # Public packages
│
├── shared/             # Shared Code
│   ├── models/        # Data models
│   ├── types/         # Type definitions
│   ├── utils/         # Utility functions
│   └── config/        # Shared configuration
│
├── database/           # Database
│   ├── migrations/    # SQL migrations
│   └── seeds/         # Seed data
│
├── docker/             # Docker Configurations
│   ├── Dockerfile.api.dev
│   ├── Dockerfile.scraper.dev
│   └── Dockerfile.cli.dev
│
├── .github/workflows/  # CI/CD
│   └── ci.yml
│
└── scripts/            # Utility Scripts
    ├── dev-setup.sh
    └── run-migrations.sh
```

---

## Development Workflow

### Daily Development

```bash
# Start services
docker-compose -f docker-compose.dev.yml up

# Make changes to code (hot reload enabled in dev mode)

# Run tests
make test

# Lint code
make lint

# Format code
make fmt
```

### Working on Backend

```bash
cd backend

# Run API server (with hot reload)
go run cmd/api/main.go

# Run tests
go test ./...

# Run specific test
go test ./internal/api -v -run TestCreateUser
```

### Working on CLI

```bash
cd cli

# Run CLI
go run cmd/main.go

# Build binary
go build -o terminal-news cmd/main.go

# Run binary
./terminal-news
```

### Working on Scraper

```bash
cd scraper

# Run scraper
go run cmd/main.go

# Test fetcher
go test ./internal/fetchers -v
```

---

## Database Management

### Access Database

```bash
# Via Docker
docker-compose -f docker-compose.dev.yml exec postgres psql -U postgres terminalnews_dev

# Via local psql
psql postgres://postgres:postgres@localhost:5432/terminalnews_dev
```

### Run Migrations

```bash
./scripts/run-migrations.sh up
```

### Reset Database

```bash
# Stop services
docker-compose -f docker-compose.dev.yml down

# Remove volumes (deletes data!)
docker volume rm terminalnews_postgres_dev_data

# Start services and run migrations
docker-compose -f docker-compose.dev.yml up -d postgres
./scripts/run-migrations.sh up
```

### Seed Test Data

```bash
psql postgres://postgres:postgres@localhost:5432/terminalnews_dev \
  -f database/seeds/test_data.sql
```

---

## Common Commands

### Using Makefile

```bash
make help           # Show all commands
make dev            # Start all services
make test           # Run all tests
make lint           # Run linters
make fmt            # Format code
make build          # Build all binaries
make clean          # Clean build artifacts
```

### Using Docker Compose

```bash
# Start all services
docker-compose -f docker-compose.dev.yml up

# Start specific service
docker-compose -f docker-compose.dev.yml up api

# View logs
docker-compose -f docker-compose.dev.yml logs -f api

# Stop services
docker-compose -f docker-compose.dev.yml down

# Rebuild services
docker-compose -f docker-compose.dev.yml up --build
```

---

## Testing

### Run All Tests

```bash
make test
```

### Run Specific Tests

```bash
# Backend tests
cd backend && go test ./...

# With coverage
cd backend && go test -cover ./...

# Specific package
cd backend && go test ./internal/api

# Specific test
cd backend && go test ./internal/api -run TestCreateUser
```

### Integration Tests

```bash
# Start test database
docker-compose -f docker-compose.test.yml up -d

# Run integration tests
go test -tags=integration ./...

# Stop test database
docker-compose -f docker-compose.test.yml down
```

---

## Troubleshooting

### Port Already in Use

```bash
# Find process using port 8080
lsof -i :8080

# Kill process
kill -9 <PID>
```

### Database Connection Failed

```bash
# Check if PostgreSQL is running
docker-compose -f docker-compose.dev.yml ps postgres

# Check logs
docker-compose -f docker-compose.dev.yml logs postgres

# Restart PostgreSQL
docker-compose -f docker-compose.dev.yml restart postgres
```

### Migrations Failed

```bash
# Check database exists
psql postgres://postgres:postgres@localhost:5432/postgres -c "\l"

# Drop and recreate database
psql postgres://postgres:postgres@localhost:5432/postgres -c "DROP DATABASE terminalnews_dev;"
psql postgres://postgres:postgres@localhost:5432/postgres -c "CREATE DATABASE terminalnews_dev;"

# Re-run migrations
./scripts/run-migrations.sh up
```

### Go Module Issues

```bash
# Clear module cache
go clean -modcache

# Re-download dependencies
cd backend && go mod download
cd ../scraper && go mod download
cd ../cli && go mod download
```

---

## API Documentation

Once the API is running, view documentation at:

**Swagger UI:** http://localhost:8080/docs (when implemented)

**Health Check:** http://localhost:8080/health

**Example Requests:**

```bash
# Get hot articles
curl http://localhost:8080/api/v1/articles/hot

# Register user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'
```

---

## Environment Variables

Key environment variables in `.env`:

```bash
# Database
DATABASE_URL=postgres://postgres:postgres@localhost:5432/terminalnews_dev?sslmode=disable

# Redis
REDIS_URL=redis://localhost:6379/0

# API
API_PORT=8080
JWT_SECRET=your_secret_here

# News APIs
NEWS_API_KEY=your_key
GUARDIAN_API_KEY=your_key

# Stripe (optional for local dev)
STRIPE_SECRET_KEY=sk_test_...

# Logging
LOG_LEVEL=debug
```

---

## Next Steps

After setup:

1. **Read the docs:** Check `/docs` folder for architecture details
2. **Pick a task:** See `ROADMAP.md` for Phase 1 tasks
3. **Write code:** Follow `CONTRIBUTING.md` guidelines
4. **Submit PR:** Push changes and open pull request

---

## Getting Help

- **Documentation:** `/docs` folder
- **Issues:** GitHub Issues
- **Questions:** GitHub Discussions (when available)
- **Slack/Discord:** (Link when available)

---

**Ready to code!** 🚀
