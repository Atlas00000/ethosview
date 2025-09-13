#!/bin/bash

# EthosView Deployment Script
# This script helps deploy EthosView to production

set -e

echo "🚀 EthosView Deployment Script"
echo "================================"

# Check if we're in the right directory
if [ ! -f "docker-compose.yml" ]; then
    echo "❌ Error: Please run this script from the EthosView root directory"
    exit 1
fi

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check prerequisites
echo "📋 Checking prerequisites..."

if ! command_exists git; then
    echo "❌ Git is not installed. Please install Git first."
    exit 1
fi

if ! command_exists docker; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

echo "✅ Prerequisites check passed"

# Get deployment type
echo ""
echo "🎯 Choose deployment type:"
echo "1) Frontend only (Vercel)"
echo "2) Backend only (Docker)"
echo "3) Full stack (Frontend + Backend)"
echo "4) Exit"

read -p "Enter your choice (1-4): " choice

case $choice in
    1)
        echo ""
        echo "🎨 Deploying Frontend to Vercel..."
        
        # Check if Vercel CLI is installed
        if ! command_exists vercel; then
            echo "📦 Installing Vercel CLI..."
            npm install -g vercel
        fi
        
        cd ethosview-frontend
        
        echo "🔧 Building frontend..."
        pnpm build
        
        echo "🚀 Deploying to Vercel..."
        vercel --prod
        
        echo "✅ Frontend deployment complete!"
        echo "📝 Don't forget to set environment variables in Vercel dashboard:"
        echo "   - NEXT_PUBLIC_API_BASE_URL"
        echo "   - NEXT_PUBLIC_SUPABASE_URL"
        echo "   - NEXT_PUBLIC_SUPABASE_ANON_KEY"
        ;;
        
    2)
        echo ""
        echo "🔧 Deploying Backend with Docker..."
        
        # Check if .env.production exists
        if [ ! -f ".env.production" ]; then
            echo "⚠️  .env.production not found. Creating from example..."
            cp .env.production.example .env.production
            echo "📝 Please edit .env.production with your actual values"
            echo "   Then run this script again."
            exit 1
        fi
        
        echo "🐳 Building and starting backend..."
        docker-compose -f docker-compose.prod.yml up -d --build
        
        echo "✅ Backend deployment complete!"
        echo "🔍 Check status with: docker-compose -f docker-compose.prod.yml ps"
        echo "📊 View logs with: docker-compose -f docker-compose.prod.yml logs -f"
        ;;
        
    3)
        echo ""
        echo "🌟 Deploying Full Stack..."
        
        # Deploy backend first
        echo "🔧 Deploying Backend..."
        
        if [ ! -f ".env.production" ]; then
            echo "⚠️  .env.production not found. Creating from example..."
            cp .env.production.example .env.production
            echo "📝 Please edit .env.production with your actual values"
            echo "   Then run this script again."
            exit 1
        fi
        
        docker-compose -f docker-compose.prod.yml up -d --build
        
        # Wait for backend to be ready
        echo "⏳ Waiting for backend to be ready..."
        sleep 10
        
        # Get backend URL (you'll need to update this with your actual URL)
        echo "📝 Please note your backend URL for frontend deployment"
        
        # Deploy frontend
        echo "🎨 Deploying Frontend..."
        
        if ! command_exists vercel; then
            echo "📦 Installing Vercel CLI..."
            npm install -g vercel
        fi
        
        cd ethosview-frontend
        pnpm build
        vercel --prod
        
        echo "✅ Full stack deployment complete!"
        echo "📝 Don't forget to update frontend environment variables with your backend URL"
        ;;
        
    4)
        echo "👋 Goodbye!"
        exit 0
        ;;
        
    *)
        echo "❌ Invalid choice. Please run the script again."
        exit 1
        ;;
esac

echo ""
echo "🎉 Deployment completed successfully!"
echo ""
echo "📚 Next steps:"
echo "   - Check the DEPLOYMENT.md file for detailed instructions"
echo "   - Monitor your deployment with the provided health check endpoints"
echo "   - Set up monitoring and alerts for production"
echo ""
echo "🔗 Useful links:"
echo "   - Frontend: https://vercel.com/dashboard"
echo "   - Backend logs: docker-compose -f docker-compose.prod.yml logs -f"
echo "   - Health check: curl https://your-backend-url.com/health/live"
