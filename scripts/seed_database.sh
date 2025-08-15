#!/bin/bash

# Database seeding script for Study Platform
# This script populates the database with development data

set -e  # Exit on any error

# Database configuration
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-2345}
DB_NAME=${DB_NAME:-studyplatform}
DB_USER=${DB_USER:-admin}
DB_PASSWORD=${DB_PASSWORD:-password}

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Study Platform Database Seeding Script${NC}"
echo -e "${BLUE}=====================================${NC}"

# Check if PostgreSQL client is available
if ! command -v psql &> /dev/null; then
    echo -e "${RED}Error: psql command not found. Please install PostgreSQL client.${NC}"
    exit 1
fi

# Test database connection
echo -e "${YELLOW}Testing database connection...${NC}"
export PGPASSWORD=$DB_PASSWORD
if ! psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "SELECT 1;" > /dev/null 2>&1; then
    echo -e "${RED}Error: Cannot connect to database.${NC}"
    echo -e "${RED}Please ensure the database is running and connection parameters are correct.${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Database connection successful${NC}"

# Check if seeds directory exists
SEEDS_DIR="$(dirname "$0")/../seeds"
if [ ! -d "$SEEDS_DIR" ]; then
    echo -e "${RED}Error: Seeds directory not found at $SEEDS_DIR${NC}"
    exit 1
fi

# Change to seeds directory
cd "$SEEDS_DIR"

echo -e "${YELLOW}Starting database seeding...${NC}"
echo -e "${YELLOW}Note: This will clear existing data and repopulate with seed data.${NC}"

# Ask for confirmation unless --force flag is provided
if [ "$1" != "--force" ]; then
    read -p "Are you sure you want to continue? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${YELLOW}Seeding cancelled.${NC}"
        exit 0
    fi
fi

# Run the master seed script
echo -e "${YELLOW}Running seed scripts...${NC}"
if psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f run_seeds_fixed.sql; then
    echo -e "${GREEN}✓ Database seeding completed successfully!${NC}"
    echo
    echo -e "${BLUE}Seed data includes:${NC}"
    echo -e "  • 10 test users (admin, instructors, students)"
    echo -e "  • 7 courses with lectures (free and paid)"
    echo -e "  • Sample enrollments and progress data"
    echo -e "  • Payment methods and transaction history"
    echo -e "  • Forum topics, posts, and votes"
    echo -e "  • Chat sessions and conversation history"
    echo
    echo -e "${BLUE}Test Credentials:${NC}"
    echo -e "  • Admin: admin@studyplatform.com / password123"
    echo -e "  • Instructor: john@studyplatform.com / password123"
    echo -e "  • Student: alice@example.com / password123"
    echo -e "  • Test User: test@example.com / password123"
    echo
else
    echo -e "${RED}Error: Database seeding failed.${NC}"
    exit 1
fi

echo -e "${GREEN}Database seeding completed successfully!${NC}"