# 🚀 Developer Quick Start

**Get coding in 10 minutes.**

---

## For the Impatient

```bash
# 1. Clone and setup
git clone <your-repo-url> terminal-news
cd terminal-news
./scripts/dev-setup.sh

# 2. Add API keys to .env
nano .env  # Add NEWS_API_KEY, GUARDIAN_API_KEY

# 3. Start everything
docker-compose -f docker-compose.dev.yml up

# 4. In another terminal, run CLI
cd cli && go run cmd/main.go
```

**Done!** You're now running Terminal News locally.

---

## Project Layout (Where to Code)

```
terminal-news/
│
├── backend/              ← Backend API work here
│   └── cmd/api/main.go   ← API entry point
│
├── cli/                  ← CLI client work here
│   └── cmd/main.go       ← CLI entry point
│
├── scraper/              ← News scraper work here
│   └── cmd/main.go       ← Scraper entry point
│
├── shared/models/        ← Shared data models
│   └── models.go         ← Add new models here
│
└── database/migrations/  ← Database changes
    └── 00X_name.sql      ← Add new migrations here
```

---

## What Each Service Does

### Backend (`backend/`)
**API Server on port 8080**
- Handles HTTP requests
- Authentication & authorization
- CRUD operations
- Serves data to CLI

**You'll work here if:** Building API endpoints, auth, business logic

### CLI (`cli/`)
**Terminal UI Client**
- Bubbletea terminal interface
- User interactions
- Calls backend API
- Local caching

**You'll work here if:** Building UI, user experience, terminal interactions

### Scraper (`scraper/`)
**News Aggregation Service**
- Fetches from RSS feeds & APIs
- Parses articles
- Stores in database
- Runs every 15 minutes

**You'll work here if:** Adding news sources, parsing logic

---

## First Tasks (Phase 1 MVP)

Pick a task based on your interest:

### Backend Tasks
- [ ] Implement user registration endpoint
- [ ] Implement login endpoint with JWT
- [ ] Create article listing endpoint (hot feed)
- [ ] Create vote endpoint
- [ ] Add health check endpoint

### CLI Tasks
- [ ] Create basic TUI layout with Bubbletea
- [ ] Implement login screen
- [ ] Implement article list view
- [ ] Add keyboard navigation
- [ ] Implement article voting

### Scraper Tasks
- [ ] Create RSS feed fetcher
- [ ] Add NewsAPI integration
- [ ] Implement deduplication logic
- [ ] Create article parser
- [ ] Add scheduling logic

---

## Development Commands

```bash
# Start database only
docker-compose -f docker-compose.dev.yml up -d postgres redis

# Run API server (with hot reload)
cd backend && go run cmd/api/main.go

# Run CLI client
cd cli && go run cmd/main.go

# Run scraper
cd scraper && go run cmd/main.go

# Run tests
make test

# Format code
make fmt

# Lint code
make lint
```

---

## Making Your First Change

### Example: Add a New API Endpoint

**1. Create handler function** (`backend/internal/api/handlers.go`):

```go
func (s *Server) HandleGetWeather(w http.ResponseWriter, r *http.Request) {
    // TODO: Implement weather endpoint
    s.respond(w, r, map[string]string{"temp": "72F"}, http.StatusOK)
}
```

**2. Register route** (`backend/internal/api/routes.go`):

```go
r.Get("/api/v1/weather", s.HandleGetWeather)
```

**3. Test it:**

```bash
curl http://localhost:8080/api/v1/weather
```

**4. Commit:**

```bash
git add .
git commit -m "feat: add weather endpoint"
git push
```

---

## Database Changes

### Adding a new table:

**1. Create migration** (`database/migrations/003_add_table.sql`):

```sql
BEGIN;

CREATE TABLE my_new_table (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

COMMIT;
```

**2. Run migration:**

```bash
./scripts/run-migrations.sh up
```

**3. Add model** (`shared/models/models.go`):

```go
type MyNewTable struct {
    ID        int64     `json:"id" db:"id"`
    Name      string    `json:"name" db:"name"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}
```

---

## Testing

### Unit Tests

```bash
# Run all tests
cd backend && go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test ./internal/api -run TestCreateUser -v
```

### Manual Testing

```bash
# Start services
docker-compose -f docker-compose.dev.yml up

# Test API
curl http://localhost:8080/api/v1/articles/hot

# Test CLI
cd cli && go run cmd/main.go
```

---

## Common Issues

### "Port 8080 already in use"

```bash
# Find and kill process
lsof -i :8080
kill -9 <PID>
```

### "Database connection failed"

```bash
# Restart PostgreSQL
docker-compose -f docker-compose.dev.yml restart postgres

# Check logs
docker-compose -f docker-compose.dev.yml logs postgres
```

### "Module not found"

```bash
# Download dependencies
cd backend && go mod download
```

---

## Code Style

### Go Conventions

```go
// Good
func GetArticles(limit int) ([]Article, error) {
    // Implementation
}

// Bad
func get_articles(l int) {
    // Implementation
}
```

### Commit Messages

```bash
# Good
git commit -m "feat: add weather widget to CLI"
git commit -m "fix: resolve database connection timeout"
git commit -m "docs: update API documentation"

# Bad
git commit -m "stuff"
git commit -m "fixed bug"
```

### PR Process

1. Create feature branch: `git checkout -b feature/my-feature`
2. Make changes
3. Write tests
4. Run `make test && make lint`
5. Commit changes
6. Push: `git push origin feature/my-feature`
7. Open PR on GitHub
8. Wait for review

---

## Getting Help

**Read first:**
- `SETUP.md` - Detailed setup guide
- `docs/ARCHITECTURE.md` - System design
- `docs/ROADMAP.md` - What to build

**Ask questions:**
- GitHub Issues - Bug reports
- GitHub Discussions - Questions
- Slack/Discord - Chat (link TBD)

**Debug tips:**
- Check logs: `docker-compose logs -f api`
- Read error messages carefully
- Search GitHub issues first
- Ask in Discord if stuck

---

## What to Build First

### MVP Must-Haves (Phase 1):

**Backend:**
1. User auth (register/login)
2. Article CRUD
3. Voting system
4. Basic API endpoints

**CLI:**
1. Login screen
2. Article list (hot feed)
3. Voting (L/D keys)
4. Basic navigation

**Scraper:**
1. RSS feed parser
2. NewsAPI integration
3. Deduplication
4. Auto-scheduling

**Target:** Working news reader in 4 weeks

---

## Daily Workflow

### Morning Routine

```bash
# 1. Pull latest changes
git pull origin main

# 2. Start services
docker-compose -f docker-compose.dev.yml up -d

# 3. Check what's working
curl http://localhost:8080/health

# 4. Pick a task from roadmap
# 5. Create feature branch
git checkout -b feature/my-task

# 6. Start coding!
```

### Before Committing

```bash
# 1. Run tests
make test

# 2. Lint code
make lint

# 3. Format code
make fmt

# 4. Commit
git add .
git commit -m "feat: describe what you did"

# 5. Push
git push origin feature/my-task

# 6. Open PR
```

---

## Resources

**Go Learning:**
- [Go Tour](https://tour.golang.org/)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go by Example](https://gobyexample.com/)

**Bubbletea (CLI):**
- [Bubbletea Tutorial](https://github.com/charmbracelet/bubbletea/tree/master/tutorials)
- [Examples](https://github.com/charmbracelet/bubbletea/tree/master/examples)

**PostgreSQL:**
- [PostgreSQL Tutorial](https://www.postgresqltutorial.com/)

**Docker:**
- [Docker Get Started](https://docs.docker.com/get-started/)

---

## Your First Day Checklist

- [ ] Clone repository
- [ ] Run `./scripts/dev-setup.sh`
- [ ] Add API keys to `.env`
- [ ] Start services with Docker Compose
- [ ] Verify API health: `curl http://localhost:8080/health`
- [ ] Run tests: `make test`
- [ ] Read `docs/ROADMAP.md`
- [ ] Pick a task
- [ ] Create feature branch
- [ ] Start coding!

---

## Tips for Success

1. **Start small** - Pick one small task, complete it, commit
2. **Test often** - Run tests after every change
3. **Ask questions** - No question is dumb
4. **Read code** - Best way to learn the codebase
5. **Write tests** - They'll save you time later
6. **Commit frequently** - Small commits are easier to review
7. **Document as you go** - Future you will thank you
8. **Have fun!** - This is a cool project

---

**Now go build something awesome!** 🚀
