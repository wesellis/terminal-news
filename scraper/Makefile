.PHONY: help build run test clean docker-build docker-run install deps fmt lint

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps: ## Download dependencies
	go mod download
	go mod tidy

build: ## Build the scraper binary
	go build -o bin/scraper cmd/scraper/main.go

run: ## Run the scraper locally
	go run cmd/scraper/main.go

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

fmt: ## Format code
	go fmt ./...
	goimports -w .

lint: ## Run linters
	golangci-lint run ./...

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out

docker-build: ## Build Docker image
	docker build -t terminal-news-scraper:latest .

docker-run: ## Run Docker container
	docker run --rm \
		--env-file .env \
		terminal-news-scraper:latest

docker-push: ## Push Docker image
	docker tag terminal-news-scraper:latest yourregistry/terminal-news-scraper:latest
	docker push yourregistry/terminal-news-scraper:latest

install: deps build ## Install dependencies and build

dev: ## Run in development mode with auto-reload
	air # Requires air: go install github.com/cosmtrek/air@latest

stats: ## Show code statistics
	@echo "Lines of code:"
	@find . -name '*.go' | xargs wc -l | tail -1
	@echo "\nFiles:"
	@find . -name '*.go' | wc -l

validate: fmt lint test ## Run all validation checks

.DEFAULT_GOAL := help
