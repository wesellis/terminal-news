# Getting Started with Terminal News

This guide will help you get Terminal News up and running on your local machine for development and testing.

## Prerequisites

Before you begin, ensure you have the following installed:

- **Go** 1.21 or higher ([Download](https://golang.org/dl/))
- **PostgreSQL** 15+ ([Download](https://www.postgresql.org/download/))
- **Redis** 7+ ([Download](https://redis.io/download))
- **Git** ([Download](https://git-scm.com/downloads))
- **Docker** (optional, recommended) ([Download](https://www.docker.com/get-started))

### Verify Installation

```bash
# Check Go version
go version

# Check PostgreSQL (if installed directly)
psql --version

# Check Redis (if installed directly)
redis-cli --version

# Check Docker (if using Docker)
docker --version
docker-compose --version
```

## Quick Start (Docker)

The fastest way to get started is using Docker:

```bash
# 1. Clone the repository
git clone https://github.com/YOUR_USERNAME/terminal-news.git
cd terminal-news

# 2. Copy environment file
cp .env.example .env

# 3. Edit .env and add your API keys (optional for basic testing)
# You'll need keys from:
# - NewsAPI.org (free tier available)
# - Guardian API (free)
# - NOAA (usually no key needed)

# 4. Start services with Docker
make docker-up
# Or: docker-compose up -d

# 5. Run database migrations
make migrate-up

# 6. Build and run the client
make build
./bin/terminal-news
```

You should now see the Terminal News interface!

## Manual Setup (Without Docker)

If you prefer to run services directly:

### 1. Install Database Services

**PostgreSQL:**
```bash
# macOS (Homebrew)
brew install postgresql@15
brew services start postgresql@15

# Ubuntu/Debian
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql

# Create database
createdb terminalnews
```

**Redis:**
```bash
# macOS (Homebrew)
brew install redis
brew services start redis

# Ubuntu/Debian
sudo apt install redis-server
sudo systemctl start redis
```

### 2. Clone and Configure

```bash
# Clone repository
git clone https://github.com/YOUR_USERNAME/terminal-news.git
cd terminal-news

# Install Go dependencies
go mod download

# Copy and configure environment
cp .env.example .env

# Edit .env with your database credentials
# Update DATABASE_PASSWORD, etc.
```

### 3. Run Database Migrations

```bash
# Create the database schema
make migrate-up

# Or manually:
go run cmd/migrate/main.go up
```

### 4. Start the API Server

In one terminal:
```bash
make api
# Or: go run cmd/api/main.go
```

The API should now be running at `http://localhost:8080`

### 5. Start the News Worker (Optional)

In another terminal:
```bash
make worker
# Or: go run cmd/worker/main.go
```

This fetches news from various sources and populates the database.

### 6. Run the Terminal Client

In a third terminal:
```bash
make run
# Or: go run cmd/terminal-news/main.go
```

## Configuration

### Environment Variables

Key configuration options in `.env`:

```bash
# API Connection
API_BASE_URL=http://localhost:8080

# Database
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_NAME=terminalnews
DATABASE_USER=postgres
DATABASE_PASSWORD=your_password

# News Sources
NEWS_API_KEY=your_newsapi_key        # Get from newsapi.org
GUARDIAN_API_KEY=your_guardian_key   # Get from open-platform.theguardian.com

# Features
ENABLE_CLASSIFIEDS=true
ENABLE_COMMENTS=true
```

### Getting API Keys

**NewsAPI (Required for news):**
1. Go to [newsapi.org](https://newsapi.org/)
2. Click "Get API Key"
3. Sign up for free tier (500 requests/day)
4. Copy key to `.env`

**Guardian API (Optional):**
1. Go to [open-platform.theguardian.com](https://open-platform.theguardian.com/)
2. Register for a developer key
3. Copy key to `.env`

**NOAA Weather (No key needed):**
- Weather works out of the box using public NOAA API

## First Run

### Create an Account

When you first run Terminal News:

1. You'll see the login screen
2. Press `R` to register
3. Enter username, email, password
4. Press `Enter` to create account

### Navigate the Interface

**Keyboard Shortcuts:**
- `Tab` - Switch between sections (Hot, Controversial, Rising, etc.)
- `↑/↓` - Navigate articles
- `O` - Open article in browser
- `L` - Like article
- `D` - Dislike article
- `C` - View comments
- `?` - Show help
- `Q` - Quit

### Explore Features

1. **Hot News** - Browse trending articles
2. **Vote** - Like/dislike to influence rankings
3. **Comments** - Press `C` on any article to read/post comments
4. **Weather** - Tab to Weather to see local forecast
5. **Classifieds** - Post or browse local listings

## Development Workflow

### Making Changes

```bash
# 1. Create a feature branch
git checkout -b feature/my-feature

# 2. Make your changes

# 3. Format code
make fmt

# 4. Run tests
make test

# 5. Build
make build

# 6. Test your changes
./bin/terminal-news

# 7. Commit and push
git add .
git commit -m "feat: add my feature"
git push origin feature/my-feature
```

### Running Tests

```bash
# Run all tests
make test

# Verbose output
make test-v

# With coverage
make coverage
# Open coverage.html in browser
```

### Code Style

```bash
# Format code
make fmt

# Run linter
make lint

# Run go vet
make vet

# All at once
make fmt && make lint && make vet && make test
```

## Common Issues

### "Connection refused" errors

**Problem:** Can't connect to API or database

**Solution:**
```bash
# Check services are running
docker-compose ps

# Or manually:
pg_isready -h localhost
redis-cli ping

# Restart services
make docker-down
make docker-up
```

### "Database does not exist"

**Problem:** Database not created

**Solution:**
```bash
# With Docker
make db-reset

# Manually
createdb terminalnews
make migrate-up
```

### "No articles showing"

**Problem:** Database is empty

**Solution:**
```bash
# Start the worker to fetch news
make worker

# Wait 1-2 minutes, then restart client
# Or seed with test data (when available):
# make seed
```

### Port already in use

**Problem:** Port 8080 or 5432 already taken

**Solution:**
```bash
# Find process using port
lsof -i :8080

# Kill process
kill -9 <PID>

# Or change port in .env
API_PORT=8081
```

### API key errors

**Problem:** "API key invalid" or rate limit errors

**Solution:**
- Check your API keys in `.env`
- Free tier NewsAPI is limited to 100 requests/day
- Use multiple RSS feeds as backup (no key needed)

## Next Steps

Now that you have Terminal News running:

1. **Explore the codebase** - See `docs/ARCHITECTURE.md`
2. **Read contribution guide** - See `CONTRIBUTING.md`
3. **Pick an issue** - Look for `good first issue` labels
4. **Join discussions** - GitHub Discussions (when available)

## Useful Commands

```bash
# Development
make dev            # Run client in dev mode
make api            # Run API server
make worker         # Run news worker

# Testing
make test           # Run tests
make coverage       # Test coverage

# Docker
make docker-up      # Start all services
make docker-down    # Stop all services
make docker-logs    # View logs

# Database
make migrate-up     # Run migrations
make migrate-down   # Rollback
make db-reset       # Reset database

# Build
make build          # Build binaries
make build-all      # Build for all platforms

# Utilities
make clean          # Clean artifacts
make deps           # Update dependencies
```

## Resources

- **Documentation:** `/docs` folder
- **API Reference:** `docs/API.md` (when available)
- **UI Mockups:** `design/UI_MOCKUPS.md`
- **Architecture:** `docs/ARCHITECTURE.md`

## Getting Help

If you run into issues:

1. Check this guide
2. Search existing GitHub issues
3. Create a new issue with details
4. Tag with `question` or `help wanted`

**Happy coding!** 🚀
