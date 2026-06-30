#!/bin/bash
# scripts/test_integration.sh
# Terminal News CLI Integration Test Script

set -e

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  Terminal News - Integration Test"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Check if backend is running
echo "⏳ Checking if backend is running..."
if ! curl -s http://localhost:8080/api/health > /dev/null 2>&1; then
    echo "❌ ERROR: Backend not running on port 8080"
    echo ""
    echo "Please start the backend first:"
    echo "  cd ../backend && go run cmd/server/main.go"
    echo ""
    exit 1
fi
echo "✓ Backend is running"
echo ""

# Check Go installation
echo "⏳ Checking Go installation..."
if ! command -v go &> /dev/null; then
    echo "❌ ERROR: Go is not installed"
    echo ""
    echo "Please install Go 1.21+ from https://golang.org/dl/"
    echo ""
    exit 1
fi
echo "✓ Go is installed ($(go version))"
echo ""

# Download dependencies
echo "⏳ Downloading dependencies..."
go mod download
echo "✓ Dependencies downloaded"
echo ""

# Build CLI
echo "⏳ Building CLI..."
go build -o bin/terminal-news cmd/terminal-news/main.go

if [ $? -ne 0 ]; then
    echo "❌ ERROR: Build failed"
    exit 1
fi
echo "✓ CLI built successfully"
echo ""

# Create test config
echo "⏳ Creating test configuration..."
cat > test_config.yaml << EOF
api:
  base_url: "http://localhost:8080/api"
  websocket_url: "ws://localhost:8080/ws"
  timeout: 30s

cache:
  database_path: "./test_cache.db"
  ttl: 3600
  max_articles: 1000
  max_comments: 5000

ui:
  theme: "default"
  compact_mode: false
  show_emojis: false
  articles_per_page: 50
  refresh_interval: 300

user:
  location: "San Francisco, CA"
  default_tab: "hot"

offline: false
EOF
echo "✓ Test configuration created"
echo ""

# Run CLI
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  Starting Terminal News CLI"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Press Ctrl+C to exit"
echo ""

./bin/terminal-news --config test_config.yaml

# Cleanup
echo ""
echo "⏳ Cleaning up test files..."
rm -f test_config.yaml test_cache.db
echo "✓ Cleanup complete"
