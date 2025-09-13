# 🚀 Vercel Deployment Guide for EthosView

## 📋 **Pre-Deployment Checklist**

✅ **Backend Ready:**
- Railway backend deployed at `https://ethosview-production.up.railway.app`
- PostgreSQL database with all tables (companies, esg_scores, users, stock_prices, financial_indicators, market_data)
- Redis cache configured
- All API endpoints working

✅ **Frontend Ready:**
- Next.js 15.5.2 with React 19
- All components working locally
- Charts and graphs loading with real data
- API integration tested

## 🔧 **Vercel Deployment Steps**

### **Step 1: Connect to Vercel**
1. Go to [vercel.com](https://vercel.com) and sign in
2. Click "New Project"
3. Import your GitHub repository: `Atlas00000/ethosview`
4. Select the `ethosview-frontend` folder as the root directory

### **Step 2: Configure Build Settings**
Vercel will auto-detect Next.js, but verify these settings:
- **Framework Preset:** Next.js
- **Root Directory:** `ethosview-frontend`
- **Build Command:** `pnpm build`
- **Output Directory:** `.next`
- **Install Command:** `corepack enable && corepack prepare pnpm@9.0.0 --activate && pnpm install --frozen-lockfile`

### **Step 3: Add Environment Variables**
In Vercel dashboard → Settings → Environment Variables, add:

```
NEXT_PUBLIC_API_URL=https://ethosview-production.up.railway.app
NEXT_PUBLIC_SUPABASE_URL=https://wrxyobquvqbwuinlikur.supabase.co
NEXT_PUBLIC_SUPABASE_ANON_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6IndyeHlvYnF1dnFid3Vpbmxpa3VyIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NTc2OTc1MDksImV4cCI6MjA3MzI3MzUwOX0.irO5qSxvXyViVvXyp_1Nn3vLIqNb9oWk4LXJ53d86n4
```

### **Step 4: Deploy**
1. Click "Deploy"
2. Wait for build to complete (should take 2-3 minutes)
3. Your app will be live at `https://your-project-name.vercel.app`

## 🎯 **Post-Deployment Verification**

### **Test These Features:**
1. **Homepage loads** with hero section
2. **Market data** displays in real-time
3. **ESG scores** show with company data
4. **Charts and graphs** render correctly
5. **Sector heatmap** displays with colors
6. **Financial data** shows P/E ratios and stock prices
7. **API calls** work (check Network tab in DevTools)

### **Expected Performance:**
- **First Load:** < 3 seconds
- **Subsequent Loads:** < 1 second (with caching)
- **API Response:** < 500ms (Railway backend)

## 🔍 **Troubleshooting**

### **Build Failures:**
- Check that all dependencies are in `package.json`
- Verify TypeScript compilation passes locally
- Ensure all imports are correct

### **Runtime Errors:**
- Check environment variables are set correctly
- Verify API URLs are accessible
- Check browser console for CORS issues

### **API Issues:**
- Test Railway backend directly: `https://ethosview-production.up.railway.app/health`
- Verify database connection in Railway logs
- Check Redis connectivity

## 📊 **Performance Optimization**

### **Already Configured:**
- ✅ **Next.js 15** with App Router
- ✅ **Turbopack** for faster builds
- ✅ **Static generation** where possible
- ✅ **Image optimization** with Next.js Image
- ✅ **Caching headers** in vercel.json

### **Future Enhancements:**
- Add CDN for static assets
- Implement service worker for offline support
- Add performance monitoring (Sentry, Vercel Analytics)

## 🌐 **Custom Domain (Optional)**

1. In Vercel dashboard → Settings → Domains
2. Add your custom domain
3. Update DNS records as instructed
4. SSL certificate will be automatically provisioned

## 🔄 **Continuous Deployment**

Your setup includes:
- ✅ **Automatic deployments** on git push to main
- ✅ **Preview deployments** for pull requests
- ✅ **Build optimization** with Vercel's edge network

## 📈 **Monitoring**

Monitor your deployment:
- **Vercel Dashboard:** Real-time logs and metrics
- **Railway Dashboard:** Backend performance and database
- **Browser DevTools:** Frontend performance and errors

---

## 🎉 **Success!**

Once deployed, your EthosView platform will be:
- **Frontend:** Hosted on Vercel's global CDN
- **Backend:** Running on Railway with PostgreSQL + Redis
- **Performance:** Optimized for speed and reliability
- **Scalability:** Auto-scaling based on traffic

**Your EthosView platform is now production-ready!** 🚀
