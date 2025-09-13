# 🚀 EthosView Quick Deploy Guide

Get EthosView running in production in under 10 minutes!

## ⚡ Super Quick Start

### 1. Frontend (Vercel) - 2 minutes
```bash
# Install Vercel CLI
npm i -g vercel

# Deploy frontend
cd ethosview-frontend
vercel --prod
```

### 2. Backend (Railway) - 5 minutes
1. Go to [railway.app](https://railway.app)
2. "Deploy from GitHub repo"
3. Select your repository
4. Add environment variables (see below)
5. Deploy!

### 3. Environment Variables

#### Frontend (Vercel Dashboard)
```
NEXT_PUBLIC_API_BASE_URL=https://your-railway-backend-url.com
NEXT_PUBLIC_SUPABASE_URL=https://your-project.supabase.co
NEXT_PUBLIC_SUPABASE_ANON_KEY=your_supabase_anon_key
```

#### Backend (Railway Dashboard)
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

## 🎯 Alternative: Use Our Script

```bash
# Make sure you're in the EthosView root directory
./deploy.sh

# Follow the interactive prompts
```

## 📋 Prerequisites

- GitHub repository with your code
- Supabase project (free tier works)
- 10 minutes of your time

## 🔗 After Deployment

1. **Update Frontend**: Add your backend URL to Vercel environment variables
2. **Test**: Visit your Vercel URL
3. **Monitor**: Check Railway dashboard for backend logs

## 🆘 Need Help?

- **Detailed Guide**: See `DEPLOYMENT.md`
- **Issues**: Check the troubleshooting section
- **Support**: Open a GitHub issue

## 💰 Cost

- **Vercel**: Free tier (perfect for most use cases)
- **Railway**: $5/month (hobby plan)
- **Supabase**: Free tier (generous limits)

**Total**: ~$5/month for full production deployment!

---

## 🎉 That's It!

Your ESG analytics platform is now live in production! 🚀
