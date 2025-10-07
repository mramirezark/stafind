# Deployment Guide - StaffFind

Complete guide for deploying the StaffFind application to production.

## Architecture Overview

StaffFind is a full-stack application with:
- **Frontend:** Next.js 14 (React)
- **Backend:** Go (Fiber framework)
- **Database:** PostgreSQL or Supabase

## Deployment Strategy

### Recommended Approach

```
┌─────────────────────────────────────────────────┐
│                                                 │
│  Frontend (Next.js)    →    Netlify/Vercel    │
│  Backend (Go)          →    Railway/Render     │
│  Database (PostgreSQL) →    Supabase           │
│                                                 │
└─────────────────────────────────────────────────┘
```

---

## Option 1: Netlify (Frontend Only) ⭐

**What Netlify is good for:**
- ✅ Static sites and Next.js frontend
- ✅ Serverless functions (Node.js/Go functions, not full servers)
- ✅ Free tier with good performance
- ✅ Automatic deployments from Git

**Limitation:**
- ❌ Cannot host your Go backend server (it's not a serverless function)

### Deploy Frontend to Netlify

#### 1. Prepare Frontend

**Create `frontend/netlify.toml`:**

```toml
[build]
  command = "npm run build"
  publish = ".next"

[[redirects]]
  from = "/api/*"
  to = "https://your-backend-url.railway.app/api/:splat"
  status = 200
  force = true

[[redirects]]
  from = "/*"
  to = "/index.html"
  status = 200
```

**Update `frontend/.env.production`:**

```env
NEXT_PUBLIC_API_URL=https://your-backend-url.railway.app
```

#### 2. Deploy to Netlify

**Via Netlify UI:**

1. Go to [app.netlify.com](https://app.netlify.com)
2. Click "Add new site" → "Import an existing project"
3. Connect your Git repository
4. Configure build settings:
   - **Base directory:** `frontend`
   - **Build command:** `npm run build`
   - **Publish directory:** `frontend/.next`
5. Add environment variables:
   - `NEXT_PUBLIC_API_URL`: Your backend URL
6. Click "Deploy site"

**Via Netlify CLI:**

```bash
# Install Netlify CLI
npm install -g netlify-cli

# Login
netlify login

# Deploy from frontend directory
cd frontend
netlify init
netlify deploy --prod
```

---

## Option 2: Vercel (Frontend Alternative)

Vercel is made by the creators of Next.js - excellent choice for frontend!

### Deploy Frontend to Vercel

```bash
# Install Vercel CLI
npm i -g vercel

# Deploy from frontend directory
cd frontend
vercel

# Production deployment
vercel --prod
```

**Configuration:**
- **Framework Preset:** Next.js
- **Root Directory:** `frontend`
- **Environment Variables:**
  - `NEXT_PUBLIC_API_URL`: Your backend URL

---

## Backend Deployment Options

### Option 1: Railway (Recommended) ⭐

**Why Railway:**
- ✅ Easy Go deployment
- ✅ Built-in PostgreSQL
- ✅ Automatic SSL
- ✅ Git-based deployments
- ✅ Free tier available

#### Deploy to Railway

**1. Create `Procfile` in `backend/`:**

```
web: ./server
```

**2. Create `railway.toml` in project root:**

```toml
[build]
  builder = "nixpacks"
  buildCommand = "cd backend && go build -o server cmd/server/main.go"

[deploy]
  startCommand = "cd backend && ./server"
  healthcheckPath = "/health"
  healthcheckTimeout = 100
```

**3. Deploy:**

```bash
# Install Railway CLI
npm i -g @railway/cli

# Login
railway login

# Link project
railway link

# Add PostgreSQL (optional, or use Supabase)
railway add postgresql

# Deploy
railway up
```

**4. Set Environment Variables:**

Go to Railway dashboard → Variables:

```env
DB_PROVIDER=postgres
DB_HOST=${PGHOST}
DB_PORT=${PGPORT}
DB_USER=${PGUSER}
DB_PASSWORD=${PGPASSWORD}
DB_NAME=${PGDATABASE}
DB_SSLMODE=require
PORT=8080
GIN_MODE=release
```

---

### Option 2: Render

**Deploy to Render:**

1. Go to [dashboard.render.com](https://dashboard.render.com)
2. New → Web Service
3. Connect Git repository
4. Configure:
   - **Name:** stafind-backend
   - **Environment:** Go
   - **Build Command:** `cd backend && go build -o server cmd/server/main.go`
   - **Start Command:** `cd backend && ./server`
   - **Environment Variables:** (same as Railway)

---

### Option 3: Fly.io

**Deploy to Fly.io:**

```bash
# Install flyctl
curl -L https://fly.io/install.sh | sh

# Login
flyctl auth login

# Initialize from backend directory
cd backend
flyctl launch

# Deploy
flyctl deploy
```

**Create `backend/fly.toml`:**

```toml
app = "stafind-backend"
primary_region = "iad"

[build]
  builder = "paketobuildpacks/builder:base"

[env]
  PORT = "8080"
  GIN_MODE = "release"

[[services]]
  http_checks = []
  internal_port = 8080
  processes = ["app"]
  protocol = "tcp"
  script_checks = []

  [services.concurrency]
    hard_limit = 25
    soft_limit = 20
    type = "connections"

  [[services.ports]]
    force_https = true
    handlers = ["http"]
    port = 80

  [[services.ports]]
    handlers = ["tls", "http"]
    port = 443

  [[services.tcp_checks]]
    grace_period = "1s"
    interval = "15s"
    restart_limit = 0
    timeout = "2s"
```

---

## Database: Supabase Setup

### 1. Create Supabase Project

1. Go to [supabase.com](https://supabase.com)
2. Create new project
3. Note your database password
4. Go to Settings → Database
5. Copy the Connection String (URI) under "Connection pooling"

### 2. Configure Backend

Set these environment variables in your backend deployment (Railway/Render/Fly.io):

```env
DB_PROVIDER=supabase
DATABASE_URL=postgresql://postgres.xxx:your-password@aws-0-us-east-1.pooler.supabase.com:6543/postgres
```

### 3. Run Migrations

```bash
# Option 1: Run locally against Supabase
cd backend
export DATABASE_URL="your-supabase-connection-string"
go run cmd/flyway-cli/main.go migrate

# Option 2: Use db-clean-go for fresh setup
make db-clean-go
```

---

## Complete Deployment Checklist

### Phase 1: Database Setup

- [ ] Create Supabase project
- [ ] Save database password
- [ ] Get connection string (pooled, port 6543)
- [ ] Run migrations
- [ ] Verify tables created

### Phase 2: Backend Deployment

- [ ] Choose platform (Railway/Render/Fly.io)
- [ ] Connect Git repository
- [ ] Set environment variables
- [ ] Configure build commands
- [ ] Deploy backend
- [ ] Test health endpoint: `https://your-backend.railway.app/health`
- [ ] Note backend URL for frontend

### Phase 3: Frontend Deployment

- [ ] Choose platform (Netlify/Vercel)
- [ ] Update `.env.production` with backend URL
- [ ] Configure build settings
- [ ] Deploy frontend
- [ ] Test frontend loads
- [ ] Verify API calls work

### Phase 4: Verification

- [ ] Test login functionality
- [ ] Test employee creation
- [ ] Test skill matching
- [ ] Test AI agent features
- [ ] Check all API endpoints work

---

## Environment Variables Reference

### Frontend (Netlify/Vercel)

```env
NEXT_PUBLIC_API_URL=https://your-backend.railway.app
```

### Backend (Railway/Render/Fly.io)

```env
# Database (Supabase)
DB_PROVIDER=supabase
DATABASE_URL=postgresql://postgres.xxx:password@aws-0-us-east-1.pooler.supabase.com:6543/postgres

# Server
PORT=8080
GIN_MODE=release
LOG_LEVEL=info

# CORS (allow your frontend domain)
CORS_ALLOWED_ORIGINS=https://your-app.netlify.app,https://your-app.vercel.app

# Optional
HUGGINGFACE_API_KEY=your-key-here
```

---

## Netlify Configuration Files

### Frontend: netlify.toml

```toml
[build]
  base = "frontend"
  command = "npm run build"
  publish = ".next"

[build.environment]
  NODE_VERSION = "18"

[[redirects]]
  from = "/api/*"
  to = "https://your-backend.railway.app/api/:splat"
  status = 200
  force = true
  headers = {X-From = "Netlify"}

[[redirects]]
  from = "/*"
  to = "/index.html"
  status = 200

[[headers]]
  for = "/*"
  [headers.values]
    X-Frame-Options = "DENY"
    X-XSS-Protection = "1; mode=block"
    X-Content-Type-Options = "nosniff"
    Referrer-Policy = "strict-origin-when-cross-origin"
```

### Frontend: _headers (in frontend/public/)

```
/*
  X-Frame-Options: DENY
  X-XSS-Protection: 1; mode=block
  X-Content-Type-Options: nosniff
  Referrer-Policy: strict-origin-when-cross-origin
```

---

## Troubleshooting

### Frontend can't connect to backend

**Problem:** API calls failing with CORS errors

**Solution:**
1. Check `NEXT_PUBLIC_API_URL` is set correctly
2. Verify backend CORS settings allow your frontend domain:
   ```env
   CORS_ALLOWED_ORIGINS=https://your-app.netlify.app
   ```

### Backend deployment fails

**Problem:** Build errors on Railway/Render

**Solution:**
1. Ensure Go version matches (check `go.mod`)
2. Verify build command is correct
3. Check logs for specific errors

### Database connection fails

**Problem:** Backend can't connect to Supabase

**Solution:**
1. Verify `DATABASE_URL` is correct
2. Check SSL mode is `require` for Supabase
3. Ensure IP isn't blocked (Supabase allows all by default)
4. Test connection locally first:
   ```bash
   psql "your-database-url"
   ```

### Migrations not running

**Problem:** Tables don't exist after deployment

**Solution:**
1. Run migrations manually:
   ```bash
   export DATABASE_URL="your-supabase-url"
   go run cmd/flyway-cli/main.go migrate
   ```
2. Or use the startup command to run migrations
3. Check migration files exist in deployment

---

## Production Best Practices

### Security

1. **Use environment variables** for all secrets
2. **Enable HTTPS** (automatic on Netlify/Railway/Vercel)
3. **Set proper CORS** origins
4. **Use Supabase Row Level Security** if needed
5. **Rotate database passwords** regularly

### Performance

1. **Enable caching** on Netlify/Vercel
2. **Use Supabase connection pooling** (port 6543)
3. **Set appropriate `DB_MAX_OPEN_CONNS`** based on your plan
4. **Monitor backend performance**

### Monitoring

1. **Set up uptime monitoring** (UptimeRobot, etc.)
2. **Check application logs** regularly
3. **Monitor database usage** in Supabase dashboard
4. **Set up error tracking** (Sentry, etc.)

### Backups

1. **Supabase automatic backups** (daily on free tier)
2. **Manual backups before major changes:**
   ```bash
   pg_dump "your-database-url" > backup.sql
   ```
3. **Test restore process**

---

## Cost Estimation

### Free Tier Setup

| Service | Free Tier | Limits |
|---------|-----------|--------|
| **Netlify** | ✅ | 100GB bandwidth, 300 build minutes |
| **Vercel** | ✅ | 100GB bandwidth, unlimited builds |
| **Railway** | ✅ $5 credit/month | ~500 hours |
| **Render** | ✅ | 750 hours/month (free tier) |
| **Supabase** | ✅ | 500MB database, 2GB bandwidth |

**Total:** $0/month for moderate usage!

### Paid Plans (when you scale)

| Service | Paid Plan | Cost |
|---------|-----------|------|
| **Netlify** | Pro | $19/month |
| **Vercel** | Pro | $20/month |
| **Railway** | Usage-based | ~$5-20/month |
| **Supabase** | Pro | $25/month |

---

## Quick Deploy Commands

### Deploy Everything

```bash
# 1. Deploy Backend to Railway
cd backend
railway up

# 2. Deploy Frontend to Netlify
cd ../frontend
netlify deploy --prod

# 3. Run migrations on Supabase
export DATABASE_URL="your-supabase-url"
cd ../backend
go run cmd/flyway-cli/main.go migrate
```

### Update Deployment

```bash
# Backend (Railway auto-deploys on git push)
git push

# Frontend (Netlify auto-deploys on git push)
git push

# Or manual:
cd frontend
netlify deploy --prod
```

---

## Alternative: Full Docker Deployment

If you want to deploy everything together, consider:
- **AWS ECS/Fargate**
- **Google Cloud Run**
- **DigitalOcean App Platform**
- **Your own VPS** (DigitalOcean, Linode, Hetzner)

See `docker-compose.yml` in the project for containerized setup.

---

## Summary

**Recommended Production Setup:**

1. **Frontend:** Netlify or Vercel (free tier)
2. **Backend:** Railway (easiest) or Render
3. **Database:** Supabase (free tier, 500MB)

**Deployment Time:** ~30 minutes

**Monthly Cost:** $0 (all free tiers) for moderate usage

**Next Steps:**
1. Set up Supabase database
2. Deploy backend to Railway
3. Deploy frontend to Netlify
4. Test everything works!

Need help with a specific platform? Let me know! 🚀

