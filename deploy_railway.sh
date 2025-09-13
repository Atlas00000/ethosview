#!/bin/bash

# Railway Deployment Script for EthosView
# This script helps you deploy to Railway with proper database setup

echo "🚀 EthosView Railway Deployment Script"
echo "======================================"

# Check if Railway CLI is installed
if ! command -v railway &> /dev/null; then
    echo "❌ Railway CLI not found. Please install it first:"
    echo "   npm install -g @railway/cli"
    echo "   or visit: https://docs.railway.app/develop/cli"
    exit 1
fi

echo "✅ Railway CLI found"

# Check if user is logged in
if ! railway whoami &> /dev/null; then
    echo "🔐 Please log in to Railway:"
    railway login
fi

echo "✅ Logged in to Railway"

# Deploy the application
echo "📦 Deploying application to Railway..."
railway up

echo "✅ Deployment complete!"
echo ""
echo "🔧 Next steps:"
echo "1. Add PostgreSQL service in Railway dashboard"
echo "2. Run the database migration:"
echo "   railway run psql < railway_migration.sql"
echo "3. Test your deployment:"
echo "   railway open"
echo ""
echo "📊 Your Railway project should now be running!"
echo "🌐 Check the logs with: railway logs"
