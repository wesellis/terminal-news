# Terminal News - Project Makefile
# Coordinates all three components: Backend, CLI, Scraper

.PHONY: help install build test clean docker-up docker-down dev fmt

# Default target
help:
	@echo "Terminal News - Available Commands"
	@echo ""
	@echo "Quick Start:"
	@echo "  make setup       - Install deps + start Docker + migrations"
	@echo "  make dev         - Start all services in dev mode"
	@echo ""
	@echo "Build:"
	@echo "  make install     - Download all Go dependencies"
	@echo "  make build       - Build all binaries (backend, cli, scraper)"
	@echo "  make build-backend  - Build backend only"
	@echo "  make build-cli      - Build CLI only"
	@echo "  make build-scraper  - Build scraper only"
	@echo ""
	@echo "Run:"
	@echo "  make run-backend - Run API server"
	@echo "  make run-cli     - Run terminal client"
	@echo "  make run-scraper - Run news scraper"
	@echo ""
	@echo "Testing:"
	@echo "  make test        - Run all tests"
	@echo "  make test-scraper - Test scraper components"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-up   - Start all Docker services"
	@echo "  make docker-down - Stop Docker services"
	@echo "  make docker-logs - View logs"
	@echo ""
	@echo "Database:"
	@echo "  make migrate-up  - Apply migrations"
	@echo "  make db-reset    - Reset database"
	@echo ""
	@echo "Utilities:"
	@echo "  make clean       - Clean all build artifacts"
	@echo "  make fmt         - Format all code"

# Install all dependencies
install:
	@echo "Installing Backend dependencies..."
	cd backend && go mod download
	@echo "Installing CLI dependencies..."
	cd cli && go mod download
	@echo "Installing Scraper dependencies..."
	cd scraper && go mod download
	@echo "✅ All dependencies installed"

# Build all components
build: build-backend build-cli build-scraper

build-backend:
	@echo "Building Backend..."
	cd backend && mkdir -p bin && go build -o bin/backend.exe cmd/server/main.go
	@echo "✅ Backend built: backend/bin/backend.exe"

build-cli:
	@echo "Building CLI..."
	cd cli && mkdir -p bin && go build -o bin/terminal-news.exe cmd/terminal-news/main.go
	@echo "✅ CLI built: cli/bin/terminal-news.exe"

build-scraper:
	@echo "Building Scraper..."
	cd scraper && mkdir -p bin && go build -o bin/scraper.exe cmd/scraper/main.go
	@echo "✅ Scraper built: scraper/bin/scraper.exe"

# Run components
run-backend: build-backend
	@echo "Starting Backend API..."
	cd backend && ./bin/backend.exe

run-cli: build-cli
	@echo "Starting CLI..."
	cd cli && ./bin/terminal-news.exe

run-scraper: build-scraper
	@echo "Starting Scraper..."
	cd scraper && ./bin/scraper.exe

# Testing
test:
	@echo "Running Backend tests..."
	cd backend && go test ./...
	@echo "Running CLI tests..."
	cd cli && go test ./...
	@echo "Running Scraper tests..."
	cd scraper && go test ./...

test-scraper:
	@echo "Testing Scraper components..."
	cd scraper && go run cmd/test/main.go
	@echo ""
	cd scraper && go run cmd/test-dedup/main.go
	@echo ""
	cd scraper && go run cmd/test-classifier/main.go

# Docker commands
docker-up:
	@echo "Starting Docker services..."
	docker-compose -f docker-compose.dev.yml up -d
	@echo "✅ Docker services started"

docker-down:
	@echo "Stopping Docker services..."
	docker-compose -f docker-compose.dev.yml down

docker-logs:
	docker-compose -f docker-compose.dev.yml logs -f

docker-restart:
	@echo "Restarting Docker services..."
	docker-compose -f docker-compose.dev.yml restart

# Database
migrate-up:
	@echo "Applying database migrations..."
	@echo "⚠️  Migrations need to be applied manually for now"
	@echo "Run: psql $$DATABASE_URL < database/migrations/001_initial_schema.sql"
	@echo "Run: psql $$DATABASE_URL < database/migrations/002_triggers_and_functions.sql"

db-reset:
	@echo "Resetting database..."
	docker-compose -f docker-compose.dev.yml down -v postgres
	docker-compose -f docker-compose.dev.yml up -d postgres
	@echo "Waiting for postgres..."
	@timeout 10
	@echo "Apply migrations with: make migrate-up"

# Development mode
dev: docker-up
	@echo "Starting development environment..."
	@echo "Backend will be available at http://localhost:8080"
	@echo "Run in separate terminals:"
	@echo "  make run-backend"
	@echo "  make run-scraper"
	@echo "  make run-cli"

# Code quality
fmt:
	@echo "Formatting code..."
	cd backend && go fmt ./...
	cd cli && go fmt ./...
	cd scraper && go fmt ./...
	@echo "✅ Code formatted"

# Clean artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf backend/bin
	rm -rf cli/bin
	rm -rf scraper/bin
	@echo "✅ Clean complete"

# Quick setup for new developers
setup: install docker-up
	@echo ""
	@echo "✅ Setup complete!"
	@echo ""
	@echo "Next steps:"
	@echo "  1. Apply migrations: make migrate-up"
	@echo "  2. Start backend: make run-backend"
	@echo "  3. Start scraper: make run-scraper"
	@echo "  4. Start CLI: make run-cli"
