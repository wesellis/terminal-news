.PHONY: build run test clean install lint format help dev release

# Variables
BINARY_NAME=terminal-news
BUILD_DIR=bin
DIST_DIR=dist
MAIN_PATH=cmd/terminal-news/main.go

# Colors for output
GREEN=\033[0;32m
YELLOW=\033[1;33m
NC=\033[0m # No Color

## help: Display this help message
help:
	@echo "$(GREEN)Terminal News CLI - Available Commands:$(NC)"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

## build: Build the application
build:
	@echo "$(GREEN)Building $(BINARY_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "$(GREEN)✓ Build complete! Binary at $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

## run: Run the application
run:
	@echo "$(GREEN)Running $(BINARY_NAME)...$(NC)"
	@go run $(MAIN_PATH)

## dev: Run with auto-reload (requires nodemon)
dev:
	@echo "$(YELLOW)Starting development mode with auto-reload...$(NC)"
	@nodemon --exec go run $(MAIN_PATH) --signal SIGTERM

## offline: Run in offline mode
offline:
	@echo "$(GREEN)Running in offline mode...$(NC)"
	@go run $(MAIN_PATH) --offline

##@ Testing

## test: Run all tests
test:
	@echo "$(GREEN)Running tests...$(NC)"
	@go test -v ./...

## test-coverage: Run tests with coverage
test-coverage:
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Coverage report generated: coverage.html$(NC)"

## test-race: Run tests with race detector
test-race:
	@echo "$(GREEN)Running tests with race detector...$(NC)"
	@go test -race -v ./...

## bench: Run benchmarks
bench:
	@echo "$(GREEN)Running benchmarks...$(NC)"
	@go test -bench=. -benchmem ./...

##@ Code Quality

## lint: Run linter (requires golangci-lint)
lint:
	@echo "$(GREEN)Running linter...$(NC)"
	@golangci-lint run ./...

## format: Format code
format:
	@echo "$(GREEN)Formatting code...$(NC)"
	@go fmt ./...
	@gofmt -s -w .
	@goimports -w .
	@echo "$(GREEN)✓ Code formatted$(NC)"

## vet: Run go vet
vet:
	@echo "$(GREEN)Running go vet...$(NC)"
	@go vet ./...

## tidy: Tidy go modules
tidy:
	@echo "$(GREEN)Tidying go modules...$(NC)"
	@go mod tidy
	@echo "$(GREEN)✓ Modules tidied$(NC)"

##@ Installation

## install: Install the binary to $GOPATH/bin
install: build
	@echo "$(GREEN)Installing $(BINARY_NAME) to $$GOPATH/bin...$(NC)"
	@cp $(BUILD_DIR)/$(BINARY_NAME) $$GOPATH/bin/
	@echo "$(GREEN)✓ Installed! Run with: $(BINARY_NAME)$(NC)"

## uninstall: Remove the binary from $GOPATH/bin
uninstall:
	@echo "$(YELLOW)Uninstalling $(BINARY_NAME)...$(NC)"
	@rm -f $$GOPATH/bin/$(BINARY_NAME)
	@echo "$(GREEN)✓ Uninstalled$(NC)"

##@ Release

## release: Build for all platforms
release: clean
	@echo "$(GREEN)Building for all platforms...$(NC)"
	@mkdir -p $(DIST_DIR)

	# Linux AMD64
	@echo "$(YELLOW)Building for Linux AMD64...$(NC)"
	@GOOS=linux GOARCH=amd64 go build -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)

	# macOS AMD64 (Intel)
	@echo "$(YELLOW)Building for macOS AMD64...$(NC)"
	@GOOS=darwin GOARCH=amd64 go build -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)

	# macOS ARM64 (Apple Silicon)
	@echo "$(YELLOW)Building for macOS ARM64...$(NC)"
	@GOOS=darwin GOARCH=arm64 go build -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)

	# Windows AMD64
	@echo "$(YELLOW)Building for Windows AMD64...$(NC)"
	@GOOS=windows GOARCH=amd64 go build -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)

	# Create archives
	@echo "$(GREEN)Creating release archives...$(NC)"
	@cd $(DIST_DIR) && tar -czf $(BINARY_NAME)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64
	@cd $(DIST_DIR) && tar -czf $(BINARY_NAME)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64
	@cd $(DIST_DIR) && tar -czf $(BINARY_NAME)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64
	@cd $(DIST_DIR) && zip -q $(BINARY_NAME)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe

	@echo "$(GREEN)✓ Release builds complete! Files in $(DIST_DIR)/$(NC)"

## release-checksums: Generate SHA256 checksums for release files
release-checksums:
	@echo "$(GREEN)Generating checksums...$(NC)"
	@cd $(DIST_DIR) && sha256sum *.tar.gz *.zip > SHA256SUMS
	@echo "$(GREEN)✓ Checksums generated$(NC)"

##@ Cleanup

## clean: Remove build artifacts
clean:
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	@rm -rf $(BUILD_DIR) $(DIST_DIR) coverage.out coverage.html
	@echo "$(GREEN)✓ Clean complete$(NC)"

## clean-cache: Remove local cache database
clean-cache:
	@echo "$(YELLOW)Removing cache database...$(NC)"
	@rm -f ~/.terminal-news/cache.db
	@echo "$(GREEN)✓ Cache cleaned$(NC)"

## clean-all: Remove all build artifacts and cache
clean-all: clean clean-cache
	@echo "$(GREEN)✓ Full cleanup complete$(NC)"

##@ Dependencies

## deps: Download dependencies
deps:
	@echo "$(GREEN)Downloading dependencies...$(NC)"
	@go mod download
	@echo "$(GREEN)✓ Dependencies downloaded$(NC)"

## deps-upgrade: Upgrade all dependencies
deps-upgrade:
	@echo "$(GREEN)Upgrading dependencies...$(NC)"
	@go get -u ./...
	@go mod tidy
	@echo "$(GREEN)✓ Dependencies upgraded$(NC)"

## deps-verify: Verify dependencies
deps-verify:
	@echo "$(GREEN)Verifying dependencies...$(NC)"
	@go mod verify
	@echo "$(GREEN)✓ Dependencies verified$(NC)"

##@ Tools

## install-tools: Install development tools
install-tools:
	@echo "$(GREEN)Installing development tools...$(NC)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@echo "$(GREEN)✓ Tools installed$(NC)"

##@ Documentation

## docs: Generate documentation
docs:
	@echo "$(GREEN)Generating documentation...$(NC)"
	@godoc -http=:6060
	@echo "$(GREEN)Documentation server running at http://localhost:6060$(NC)"

##@ Info

## version: Show version information
version:
	@echo "$(GREEN)Terminal News CLI v0.1.0$(NC)"
	@echo "Go version: $$(go version)"

## size: Show binary sizes
size:
	@echo "$(GREEN)Binary sizes:$(NC)"
	@ls -lh $(BUILD_DIR)/$(BINARY_NAME) 2>/dev/null || echo "Build the binary first with 'make build'"
	@echo ""
	@echo "$(GREEN)Release sizes:$(NC)"
	@ls -lh $(DIST_DIR)/ 2>/dev/null || echo "Build releases first with 'make release'"

# Default target
.DEFAULT_GOAL := help
