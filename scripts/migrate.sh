#!/bin/bash

# Database migration script for EthosView
# This script applies all database migrations in order

set -e

# Database configuration
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_NAME=${DB_NAME:-ethosview}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-password}

echo "Starting database migrations..."

# Apply migrations in order
echo "Applying initial schema migration..."
psql "host=$DB_HOST port=$DB_PORT dbname=$DB_NAME user=$DB_USER password=$DB_PASSWORD" -f scripts/migrations/001_initial_schema.sql

echo "Applying financial data migration..."
psql "host=$DB_HOST port=$DB_PORT dbname=$DB_NAME user=$DB_USER password=$DB_PASSWORD" -f scripts/migrations/002_financial_data.sql

echo "Applying performance optimization migration..."
psql "host=$DB_HOST port=$DB_PORT dbname=$DB_NAME user=$DB_USER password=$DB_PASSWORD" -f scripts/migrations/003_performance_optimization.sql

echo "Database migrations completed successfully!"

# Optional: Run seed data
if [ "$1" = "--seed" ]; then
    echo "Applying seed data..."
    psql "host=$DB_HOST port=$DB_PORT dbname=$DB_NAME user=$DB_USER password=$DB_PASSWORD" -f scripts/seeds/sample_data.sql
    psql "host=$DB_HOST port=$DB_PORT dbname=$DB_NAME user=$DB_USER password=$DB_PASSWORD" -f scripts/seeds/financial_data.sql
    echo "Seed data applied successfully!"
fi
