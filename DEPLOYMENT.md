# EthosView Deployment Guide

This guide covers deploying EthosView to production using Vercel for the frontend and various options for the backend.

## 🚀 Quick Start

### Frontend (Vercel)
1. **Deploy to Vercel**: Push to GitHub and connect to Vercel
2. **Set Environment Variables**: Configure in Vercel dashboard
3. **Deploy**: Automatic deployment on push

### Backend Options
- **Railway**: Simple Docker deployment
- **Render**: Managed Docker containers
- **DigitalOcean**: Droplet with Docker Compose
- **AWS/GCP**: Container services

## 📋 Prerequisites

- GitHub repository with your code
- Supabase project with database setup
- Domain name (optional but recommended)

## 🎯 Frontend Deployment (Vercel)

### 1. Prepare Repository
```bash
# Ensure your code is pushed to GitHub
git add .
git commit -m "Prepare for Vercel deployment"
git push origin main
```

### 2. Deploy to Vercel

#### Option A: Vercel CLI
```bash
# Install Vercel CLI
npm i -g vercel

# Navigate to frontend directory
cd ethosview-frontend

# Deploy
vercel

# Follow the prompts:
# - Set up and deploy? Y
# - Which scope? (your account)
# - Link to existing project? N
# - Project name: ethosview-frontend
# - Directory: ./
# - Override settings? N
```

#### Option B: Vercel Dashboard
1. Go to [vercel.com](https://vercel.com)
2. Click "New Project"
3. Import your GitHub repository
4. **Important**: Set Root Directory to `ethosview-frontend`
5. Configure environment variables (see below)
6. Deploy

### 3. Environment Variables (Vercel)

In Vercel dashboard, go to Settings → Environment Variables:

```
NEXT_PUBLIC_API_BASE_URL=https://your-backend-url.com
NEXT_PUBLIC_SUPABASE_URL=https://your-project.supabase.co
NEXT_PUBLIC_SUPABASE_ANON_KEY=your_supabase_anon_key
```

### 4. Custom Domain (Optional)
1. In Vercel dashboard → Settings → Domains
2. Add your domain
3. Configure DNS records as instructed

## 🔧 Backend Deployment Options

### Option 1: Railway (Recommended)

Railway is the easiest option for Docker deployments:

1. **Sign up**: [railway.app](https://railway.app)
2. **Create Project**: "Deploy from GitHub repo"
3. **Select Repository**: Your EthosView repository
4. **Configure**:
   - Root Directory: Leave empty (uses root)
   - Dockerfile: Uses existing `Dockerfile`
5. **Environment Variables**:
   ```
   DB_HOST=aws-1-us-east-2.pooler.supabase.com
   DB_PORT=5432
   DB_NAME=postgres
   DB_USER=postgres.your-project-id
   DB_PASSWORD=your_supabase_password
   DB_SSL_MODE=require
   REDIS_HOST=redis:6379
   PORT=8080
   ENVIRONMENT=production
   SUPABASE_URL=https://your-project.supabase.co
   SUPABASE_ANON_KEY=your_supabase_anon_key
   SUPABASE_SERVICE_ROLE_KEY=your_supabase_service_role_key
   JWT_SECRET=your_secure_jwt_secret_32_chars_minimum
   ```
6. **Deploy**: Railway automatically builds and deploys

### Option 2: Render

1. **Sign up**: [render.com](https://render.com)
2. **Create Web Service**: "Build and deploy from a Git repository"
3. **Configure**:
   - Build Command: `docker build -t ethosview-backend .`
   - Start Command: `docker run -p 10000:8080 ethosview-backend`
   - Environment: Docker
4. **Environment Variables**: Same as Railway
5. **Deploy**

### Option 3: DigitalOcean Droplet

1. **Create Droplet**: Ubuntu 22.04 LTS, 2GB RAM minimum
2. **Install Docker**:
   ```bash
   sudo apt update
   sudo apt install docker.io docker-compose
   sudo systemctl start docker
   sudo systemctl enable docker
   ```
3. **Clone Repository**:
   ```bash
   git clone https://github.com/your-username/ethosview.git
   cd ethosview
   ```
4. **Configure Environment**:
   ```bash
   cp .env.production.example .env.production
   nano .env.production  # Edit with your values
   ```
5. **Deploy**:
   ```bash
   docker-compose -f docker-compose.prod.yml up -d
   ```

### Option 4: AWS/GCP Container Services

#### AWS ECS/Fargate
1. **Build Image**: Push to ECR
2. **Create Task Definition**: Use your Dockerfile
3. **Create Service**: Configure networking and load balancing
4. **Environment Variables**: Set in task definition

#### Google Cloud Run
1. **Build Image**: `gcloud builds submit --tag gcr.io/PROJECT-ID/ethosview-backend`
2. **Deploy**: `gcloud run deploy --image gcr.io/PROJECT-ID/ethosview-backend`
3. **Environment Variables**: Set in Cloud Run service

## 🔗 Connecting Frontend to Backend

After deploying both:

1. **Update Frontend Environment Variables**:
   ```
   NEXT_PUBLIC_API_BASE_URL=https://your-backend-url.com
   ```

2. **Redeploy Frontend**: Vercel will automatically redeploy when you push changes

## 🗄️ Database Setup

### Supabase Configuration
1. **Create Project**: [supabase.com](https://supabase.com)
2. **Run Migration**:
   ```bash
   # Connect to your Supabase database
   psql "postgresql://postgres.your-project-id:password@aws-1-us-east-2.pooler.supabase.com:5432/postgres"
   
   # Run the migration
   \i supabase_migration.sql
   ```
3. **Get Connection Details**: Project Settings → Database

## 🔒 Security Checklist

### Frontend (Vercel)
- ✅ Environment variables configured
- ✅ HTTPS enabled (automatic with Vercel)
- ✅ CORS properly configured

### Backend
- ✅ Environment variables secured
- ✅ JWT secret is strong (32+ characters)
- ✅ Database credentials secured
- ✅ Rate limiting enabled
- ✅ HTTPS/SSL enabled
- ✅ CORS configured for your domain

## 📊 Monitoring & Maintenance

### Health Checks
- **Frontend**: Vercel provides automatic monitoring
- **Backend**: Use `/health/live` endpoint

### Logs
- **Frontend**: Vercel dashboard → Functions tab
- **Backend**: 
  - Railway: Dashboard logs
  - Render: Logs tab
  - DigitalOcean: `docker-compose logs -f`

### Updates
1. **Push Changes**: `git push origin main`
2. **Frontend**: Auto-deploys with Vercel
3. **Backend**: Restart service or auto-deploy (depending on platform)

## 🆘 Troubleshooting

### Common Issues

#### Frontend Build Failures
```bash
# Check build locally
cd ethosview-frontend
pnpm build

# Common fixes:
# - Check environment variables
# - Ensure all dependencies are in package.json
# - Check TypeScript errors
```

#### Backend Connection Issues
```bash
# Test database connection
curl https://your-backend-url.com/health/live

# Check logs
# Railway: railway logs
# Render: Check logs tab
# DigitalOcean: docker-compose logs backend
```

#### CORS Issues
- Ensure `CORS_ORIGIN` includes your Vercel domain
- Check that frontend URL is correct in backend environment

### Performance Optimization

#### Frontend
- Enable Vercel Analytics
- Use Vercel's Edge Functions for API routes
- Optimize images with Next.js Image component

#### Backend
- Use Redis caching
- Enable compression middleware
- Monitor database query performance
- Use connection pooling

## 📈 Scaling

### Frontend (Vercel)
- Automatic scaling with Vercel
- Consider Vercel Pro for higher limits
- Use Edge Functions for global performance

### Backend
- **Railway**: Automatic scaling
- **Render**: Configure scaling in dashboard
- **DigitalOcean**: Use Load Balancer + multiple droplets
- **AWS/GCP**: Use managed services with auto-scaling

## 💰 Cost Estimation

### Free Tiers
- **Vercel**: Free tier supports most use cases
- **Railway**: $5/month for hobby projects
- **Render**: Free tier available
- **Supabase**: Free tier with generous limits

### Production Costs
- **Vercel Pro**: $20/month
- **Railway**: $5-20/month
- **Render**: $7-25/month
- **DigitalOcean**: $12-24/month

## 🔄 CI/CD Pipeline

### GitHub Actions (Optional)
Create `.github/workflows/deploy.yml`:

```yaml
name: Deploy
on:
  push:
    branches: [main]
jobs:
  deploy-frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Deploy to Vercel
        uses: amondnet/vercel-action@v20
        with:
          vercel-token: ${{ secrets.VERCEL_TOKEN }}
          vercel-org-id: ${{ secrets.ORG_ID }}
          vercel-project-id: ${{ secrets.PROJECT_ID }}
          working-directory: ./ethosview-frontend
```

## 📞 Support

- **Vercel**: [vercel.com/docs](https://vercel.com/docs)
- **Railway**: [docs.railway.app](https://docs.railway.app)
- **Render**: [render.com/docs](https://render.com/docs)
- **Supabase**: [supabase.com/docs](https://supabase.com/docs)

---

## 🎉 You're Ready to Deploy!

Follow the steps above to get EthosView running in production. Start with the Vercel deployment for the frontend, then choose your preferred backend hosting option.

**Recommended Path**: Vercel (Frontend) + Railway (Backend) for the fastest setup.
