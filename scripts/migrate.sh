#!/bin/bash

# Database migration script for EthosView
# Week 2: Database Schema & Basic Structure

set -e

echo "🚀 Starting EthosView database migration..."

# Check if PostgreSQL container is running
if ! docker ps | grep -q ethosview-postgres; then
    echo "❌ PostgreSQL container is not running. Please start the services first:"
    echo "   docker-compose up -d"
    exit 1
fi

echo "📊 Running database migrations..."

# Run the initial schema migration
echo "Creating database schema..."
docker exec -i ethosview-postgres psql -U postgres -d ethosview < scripts/migrations/001_initial_schema.sql

echo "✅ Database schema created successfully!"

# Run the seed data
echo "🌱 Seeding sample data..."
docker exec -i ethosview-postgres psql -U postgres -d ethosview < scripts/seeds/sample_data.sql

echo "✅ Sample data seeded successfully!"

echo "🎉 Database migration completed!"
echo ""
echo "📋 Summary:"
echo "   - Database schema created"
echo "   - Sample companies added (10 companies)"
echo "   - Sample ESG scores added (16 scores)"
echo "   - Indexes created for performance"
echo "   - Triggers created for updated_at timestamps"
echo ""
echo "🔗 You can now test the API endpoints:"
echo "   - GET /api/v1/companies"
echo "   - GET /api/v1/esg/scores"
echo "   - GET /api/v1/companies/sectors"
