# Netlify Deployment - Quick Start

Quick guide to deploy StaffFind frontend to Netlify in 5 minutes.

## ⚡ Quick Steps

### 1. Deploy Backend First (Required)

Netlify only hosts your frontend. Deploy your backend to Railway:

```bash
# Install Railway CLI
npm i -g @railway/cli

# Login and deploy
railway login
cd backend
railway init
railway up
```

**Copy your Railway backend URL** (e.g., `https://stafind-backend.railway.app`)

---

### 2. Update Frontend Configuration

**Edit `frontend/netlify.toml`:**

Replace `your-backend-url.railway.app` with your actual backend URL:

```toml
[[redirects]]
  from = "/api/*"
  to = "https://stafind-backend.railway.app/api/:splat"  # ← Change this!
  status = 200
  force = true
```

**Create `frontend/.env.production`:**

```env
NEXT_PUBLIC_API_URL=https://stafind-backend.railway.app
```

---

### 3. Deploy to Netlify

#### Option A: Via Netlify UI (Easiest)

1. Go to [app.netlify.com](https://app.netlify.com)
2. Click **"Add new site"** → **"Import an existing project"**
3. Connect your Git repository (GitHub/GitLab/Bitbucket)
4. Configure build settings:
   - **Base directory:** `frontend`
   - **Build command:** `npm run build`
   - **Publish directory:** `frontend/.next`
5. Add environment variable:
   - **Key:** `NEXT_PUBLIC_API_URL`
   - **Value:** `https://your-backend.railway.app`
6. Click **"Deploy site"**

#### Option B: Via Netlify CLI

```bash
# Install Netlify CLI
npm install -g netlify-cli

# Login to Netlify
netlify login

# Deploy from project root
netlify deploy --dir=frontend --prod
```

---

### 4. Configure Database (Supabase)

1. Create Supabase project at [supabase.com](https://supabase.com)
2. Get your connection string from Settings → Database
3. Add to Railway backend environment:
   ```env
   DB_PROVIDER=supabase
   DATABASE_URL=postgresql://postgres.xxx:password@aws-0-us-east-1.pooler.supabase.com:6543/postgres
   ```
4. Run migrations:
   ```bash
   cd backend
   export DATABASE_URL="your-supabase-url"
   go run cmd/flyway-cli/main.go migrate
   ```

---

## ✅ Verification

After deployment, test:

1. **Frontend loads:** Visit your Netlify URL
2. **API works:** Try logging in
3. **Database connected:** Create an employee

---

## 🔧 Common Issues

### Issue: "API calls failing with 404"

**Solution:** Update the backend URL in `netlify.toml` and redeploy:

```bash
cd frontend
netlify deploy --prod
```

### Issue: "Build failed on Netlify"

**Solution:** Check build logs. Common fixes:

1. **Missing dependencies:**
   ```bash
   cd frontend
   npm install
   git add package-lock.json
   git commit -m "Update dependencies"
   git push
   ```

2. **Node version mismatch:**
   - Add to `netlify.toml`:
     ```toml
     [build.environment]
       NODE_VERSION = "18"
     ```

### Issue: "CORS errors in browser"

**Solution:** Update backend CORS settings in Railway:

```env
CORS_ALLOWED_ORIGINS=https://your-app.netlify.app
```

---

## 📝 Environment Variables

Set these in Netlify dashboard (Site settings → Environment variables):

| Variable | Value | Required |
|----------|-------|----------|
| `NEXT_PUBLIC_API_URL` | Your Railway backend URL | ✅ Yes |
| `NODE_VERSION` | `18` | Recommended |

---

## 🚀 Deploy Updates

After making changes:

```bash
# Commit your changes
git add .
git commit -m "Update feature"
git push

# Netlify automatically deploys from your main branch!
# Or deploy manually:
cd frontend
netlify deploy --prod
```

---

## 💰 Cost

**Netlify Free Tier:**
- ✅ 100GB bandwidth/month
- ✅ 300 build minutes/month
- ✅ Automatic deployments
- ✅ Custom domain
- ✅ HTTPS included

**Perfect for development and small apps!**

---

## 🎯 Complete Setup Summary

```
┌─────────────────────────────────────────────┐
│                                             │
│  1. Backend → Railway                       │
│     └─ Database → Supabase                 │
│                                             │
│  2. Frontend → Netlify                      │
│     └─ Connects to Railway backend         │
│                                             │
│  Total Time: ~10 minutes                    │
│  Total Cost: $0 (all free tiers)           │
│                                             │
└─────────────────────────────────────────────┘
```

---

## 📚 Next Steps

After deploying:

1. ✅ Set up custom domain in Netlify
2. ✅ Enable branch deploys for staging
3. ✅ Configure error tracking (Sentry)
4. ✅ Set up monitoring (UptimeRobot)
5. ✅ Add Google Analytics (optional)

---

## 🆘 Need Help?

- **Netlify Docs:** [docs.netlify.com](https://docs.netlify.com)
- **Railway Docs:** [docs.railway.app](https://docs.railway.app)
- **Supabase Docs:** [supabase.com/docs](https://supabase.com/docs)
- **Full Guide:** See [DEPLOYMENT_GUIDE.md](DEPLOYMENT_GUIDE.md)

**Happy deploying! 🎉**

