# Terminal News - Makefile

.PHONY: help build run test clean docker-up docker-down install dev fmt lint

# Default target
help:
	@echo "Terminal News - Available Commands"
	@echo ""
	@echo "Development:"
	@echo "  make install    - Install dependencies"
	@echo "  make dev        - Run in development mode"
	@echo "  make build      - Build binaries"
	@echo "  make run        - Run the client"
	@echo "  make api        - Run the API server"
	@echo ""
	@echo "Testing:"
	@echo "  make test       - Run tests"
	@echo "  make test-v     - Run tests (verbose)"
	@echo "  make coverage   - Generate test coverage"
	@echo ""
	@echo "Code Quality:"
	@echo "  make fmt        - Format code"
	@echo "  make lint       - Run linters"
	@echo "  make vet        - Run go vet"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-up  - Start Docker services"
	@echo "  make docker-down - Stop Docker services"
	@echo "  make docker-logs - View Docker logs"
	@echo ""
	@echo "Database:"
	@echo "  make migrate-up   - Run database migrations"
	@echo "  make migrate-down - Rollback migrations"
	@echo "  make db-reset     - Reset database"
	@echo ""
	@echo "Utilities:"
	@echo "  make clean      - Clean build artifacts"
	@echo "  make deps       - Update dependencies"

# Install dependencies
install:
	@echo "Installing dependencies..."
	go mod download
	go mod verify
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Development mode (with auto-reload)
dev:
	@echo "Starting development mode..."
	go run cmd/terminal-news/main.go

# Build binaries
build:
	@echo "Building binaries..."
	mkdir -p bin
	go build -o bin/terminal-news cmd/terminal-news/main.go
	go build -o bin/api cmd/api/main.go
	go build -o bin/worker cmd/worker/main.go
	@echo "Binaries built in ./bin/"

# Build for all platforms
build-all:
	@echo "Building for all platforms..."
	mkdir -p bin
	GOOS=darwin GOARCH=amd64 go build -o bin/terminal-news-darwin-amd64 cmd/terminal-news/main.go
	GOOS=darwin GOARCH=arm64 go build -o bin/terminal-news-darwin-arm64 cmd/terminal-news/main.go
	GOOS=linux GOARCH=amd64 go build -o bin/terminal-news-linux-amd64 cmd/terminal-news/main.go
	GOOS=linux GOARCH=arm64 go build -o bin/terminal-news-linux-arm64 cmd/terminal-news/main.go
	GOOS=windows GOARCH=amd64 go build -o bin/terminal-news-windows-amd64.exe cmd/terminal-news/main.go
	@echo "Cross-platform binaries built!"

# Run the terminal client
run: build
	@echo "Running Terminal News..."
	./bin/terminal-news

# Run API server
api:
	@echo "Starting API server..."
	go run cmd/api/main.go

# Run worker
worker:
	@echo "Starting news aggregation worker..."
	go run cmd/worker/main.go

# Tests
test:
	@echo "Running tests..."
	go test ./...

test-v:
	@echo "Running tests (verbose)..."
	go test -v ./...

coverage:
	@echo "Generating test coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Code quality
fmt:
	@echo "Formatting code..."
	go fmt ./...

lint:
	@echo "Running linters..."
	golangci-lint run

vet:
	@echo "Running go vet..."
	go vet ./...

# Docker commands
docker-up:
	@echo "Starting Docker services..."
	docker-compose up -d

docker-down:
	@echo "Stopping Docker services..."
	docker-compose down

docker-logs:
	docker-compose logs -f

docker-build:
	@echo "Building Docker images..."
	docker-compose build

docker-reset:
	@echo "Resetting Docker environment..."
	docker-compose down -v
	docker-compose up -d

# Database migrations
migrate-up:
	@echo "Running database migrations..."
	go run cmd/migrate/main.go up

migrate-down:
	@echo "Rolling back database migrations..."
	go run cmd/migrate/main.go down

migrate-create:
	@echo "Creating new migration..."
	@read -p "Migration name: " name; \
	go run cmd/migrate/main.go create $$name

db-reset:
	@echo "Resetting database..."
	docker-compose down -v db
	docker-compose up -d db
	sleep 5
	make migrate-up
	@echo "Database reset complete!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf dist/
	rm -f coverage.out coverage.html
	go clean

# Update dependencies
deps:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

# Create release
release:
	@echo "Creating release..."
	@read -p "Version (e.g., v1.0.0): " version; \
	git tag -a $$version -m "Release $$version"; \
	git push origin $$version
	@echo "Release tagged. CI will build and publish."

# Quick start for new developers
setup: install docker-up migrate-up
	@echo ""
	@echo "Setup complete! You can now run:"
	@echo "  make dev    - to start the client"
	@echo "  make api    - to start the API server"
