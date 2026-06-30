#!/bin/bash

# Terminal News Scraper Setup Script

set -e

echo "🚀 Setting up Terminal News Scraper..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21+ first."
    exit 1
fi

echo "✅ Go version: $(go version)"

# Check if .env exists
if [ ! -f .env ]; then
    echo "📝 Creating .env from template..."
    cp .env.example .env
    echo "⚠️  Please edit .env with your configuration before running!"
else
    echo "✅ .env file already exists"
fi

# Install dependencies
echo "📦 Installing dependencies..."
go mod download
go mod tidy

# Create bin directory
mkdir -p bin

# Build the scraper
echo "🔨 Building scraper..."
go build -o bin/scraper cmd/scraper/main.go

echo ""
echo "✅ Setup complete!"
echo ""
echo "Next steps:"
echo "  1. Edit .env with your DATABASE_URL and API keys"
echo "  2. Ensure PostgreSQL is running"
echo "  3. Run: ./bin/scraper"
echo ""
echo "Or use Docker:"
echo "  docker build -t terminal-news-scraper ."
echo "  docker run --env-file .env terminal-news-scraper"
echo ""
