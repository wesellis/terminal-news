#!/bin/bash
# Database Migration Script for Terminal News

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get database URL from environment or default
DATABASE_URL=${DATABASE_URL:-"postgres://postgres:postgres@localhost:5432/terminalnews_dev?sslmode=disable"}

echo -e "${GREEN}Terminal News - Database Migration Tool${NC}"
echo ""

# Check if psql is available
if ! command -v psql &> /dev/null; then
    echo -e "${RED}Error: psql is not installed${NC}"
    echo "Install PostgreSQL client tools"
    exit 1
fi

# Check if database is reachable
echo -e "${YELLOW}Checking database connection...${NC}"
if ! psql "$DATABASE_URL" -c "SELECT 1" > /dev/null 2>&1; then
    echo -e "${RED}Error: Cannot connect to database${NC}"
    echo "Database URL: $DATABASE_URL"
    echo ""
    echo "Make sure PostgreSQL is running:"
    echo "  docker-compose -f docker-compose.dev.yml up -d postgres"
    exit 1
fi

echo -e "${GREEN}✓ Database connection successful${NC}"
echo ""

# Function to apply migration
apply_migration() {
    local file=$1
    local name=$(basename "$file")

    echo -e "${YELLOW}Applying migration: $name${NC}"

    if psql "$DATABASE_URL" -f "$file"; then
        echo -e "${GREEN}✓ Migration applied: $name${NC}"
        return 0
    else
        echo -e "${RED}✗ Migration failed: $name${NC}"
        return 1
    fi
}

# Main migration logic
MIGRATION_DIR="$(dirname "$0")/../database/migrations"

case "${1:-up}" in
    up)
        echo "Applying all migrations..."
        echo ""

        for migration in "$MIGRATION_DIR"/*.sql; do
            if [[ $migration == *"down"* ]]; then
                continue
            fi
            apply_migration "$migration"
            echo ""
        done

        echo -e "${GREEN}All migrations applied successfully!${NC}"
        ;;

    down)
        echo -e "${RED}Rolling back migrations not implemented yet${NC}"
        echo "To manually rollback, drop and recreate the database"
        exit 1
        ;;

    status)
        echo "Checking migration status..."
        psql "$DATABASE_URL" -c "
            SELECT
                schemaname,
                tablename
            FROM pg_tables
            WHERE schemaname = 'public'
            ORDER BY tablename;
        "
        ;;

    reset)
        echo -e "${YELLOW}WARNING: This will drop all tables!${NC}"
        read -p "Are you sure? (yes/no): " confirm

        if [ "$confirm" != "yes" ]; then
            echo "Aborted"
            exit 0
        fi

        echo "Dropping all tables..."
        psql "$DATABASE_URL" -c "
            DROP SCHEMA public CASCADE;
            CREATE SCHEMA public;
            GRANT ALL ON SCHEMA public TO postgres;
            GRANT ALL ON SCHEMA public TO public;
        "

        echo -e "${GREEN}Database reset complete${NC}"
        echo "Run: ./scripts/migrate.sh up"
        ;;

    *)
        echo "Usage: $0 {up|down|status|reset}"
        echo ""
        echo "Commands:"
        echo "  up     - Apply all pending migrations"
        echo "  down   - Rollback last migration (not implemented)"
        echo "  status - Show current database tables"
        echo "  reset  - Drop all tables and reset database"
        exit 1
        ;;
esac
