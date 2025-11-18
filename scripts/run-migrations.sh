#!/bin/bash
# Run database migrations
# Usage: ./scripts/run-migrations.sh [up|down] [DATABASE_URL]

set -e

ACTION=${1:-up}
DATABASE_URL=${2:-$DATABASE_URL}

if [ -z "$DATABASE_URL" ]; then
    echo "❌ DATABASE_URL not set"
    echo "Usage: ./scripts/run-migrations.sh [up|down] [DATABASE_URL]"
    echo "Or set DATABASE_URL environment variable"
    exit 1
fi

echo "🔄 Running migrations ($ACTION)..."

if [ "$ACTION" = "up" ]; then
    # Run all migrations
    for migration in database/migrations/*.sql; do
        echo "Running $migration..."
        psql "$DATABASE_URL" -f "$migration"
    done
    echo "✅ Migrations complete"
elif [ "$ACTION" = "down" ]; then
    # For down migrations, you'd need separate down.sql files
    echo "⚠️  Down migrations not implemented yet"
    echo "To rollback, restore from backup or recreate database"
else
    echo "❌ Invalid action: $ACTION"
    echo "Use 'up' or 'down'"
    exit 1
fi
