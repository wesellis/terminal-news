#!/bin/bash
# Development Environment Setup Script
# Run this to set up your local development environment

set -e

echo "🚀 Terminal News - Development Setup"
echo "======================================"

# Check prerequisites
echo ""
echo "Checking prerequisites..."

# Check Docker
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi
echo "✅ Docker found"

# Check Docker Compose
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi
echo "✅ Docker Compose found"

# Check Go
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21+ first."
    exit 1
fi
echo "✅ Go found: $(go version)"

# Check Git
if ! command -v git &> /dev/null; then
    echo "❌ Git is not installed. Please install Git first."
    exit 1
fi
echo "✅ Git found"

echo ""
echo "Setting up environment..."

# Create .env if it doesn't exist
if [ ! -f .env ]; then
    echo "Creating .env file from example..."
    cp .env.example .env
    echo "✅ .env file created"
    echo "⚠️  Please edit .env and add your API keys"
else
    echo "✅ .env file already exists"
fi

# Start Docker services
echo ""
echo "Starting Docker services..."
docker-compose -f docker-compose.dev.yml up -d postgres redis

# Wait for PostgreSQL
echo "Waiting for PostgreSQL..."
until docker-compose -f docker-compose.dev.yml exec -T postgres pg_isready -U postgres; do
    sleep 1
done
echo "✅ PostgreSQL is ready"

# Run migrations
echo ""
echo "Running database migrations..."
for migration in database/migrations/*.sql; do
    echo "Running $migration..."
    docker-compose -f docker-compose.dev.yml exec -T postgres \
        psql -U postgres -d terminalnews_dev -f - < "$migration"
done
echo "✅ Migrations complete"

# Install Go dependencies
echo ""
echo "Installing Go dependencies..."

cd backend && go mod download && cd ..
cd scraper && go mod download && cd ..
cd cli && go mod download && cd ..

echo "✅ Go dependencies installed"

echo ""
echo "======================================"
echo "✅ Setup complete!"
echo ""
echo "Next steps:"
echo "1. Edit .env and add your API keys"
echo "2. Run: make dev-api     (start API server)"
echo "3. Run: make dev-scraper (start scraper)"
echo "4. Run: make dev-cli     (start CLI client)"
echo ""
echo "Or run all services: docker-compose -f docker-compose.dev.yml up"
echo ""
