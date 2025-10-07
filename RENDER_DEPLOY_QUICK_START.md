# Deploy to Render - 10 Minute Setup

Deploy your **complete StaffFind application** (frontend + backend + database) to Render for **FREE**!

## 🎯 What You'll Get

- ✅ Frontend: `https://stafind-frontend.onrender.com`
- ✅ Backend API: `https://stafind-backend.onrender.com`
- ✅ PostgreSQL database (free for 90 days)
- ✅ Automatic deployments on git push
- ✅ Free SSL certificates
- ✅ **Total cost: $0/month**

---

## ⚡ 10-Minute Setup

### Step 1: Prepare Your Repository (1 minute)

The `render.yaml` file is already created for you! Just verify it exists:

```bash
ls render.yaml  # Should exist
```

If you want to use **Supabase instead of Render PostgreSQL**, update `render.yaml`:

```yaml
# In render.yaml, remove the databases section and update backend envVars:
envVars:
  - key: DB_PROVIDER
    value: supabase
  - key: DATABASE_URL
    value: your-supabase-connection-string  # Add your Supabase URL
```

### Step 2: Push to Git (1 minute)

```bash
git add render.yaml
git commit -m "Add Render deployment config"
git push
```

### Step 3: Connect to Render (2 minutes)

1. Go to [dashboard.render.com](https://dashboard.render.com)
2. Sign up/Login (free, no credit card required)
3. Click **"New +"** → **"Blueprint"**
4. Click **"Connect a repository"**
5. Select your StaffFind repository
6. Click **"Apply"**

Render will automatically:
- ✅ Create PostgreSQL database
- ✅ Create backend service
- ✅ Create frontend service
- ✅ Link them together
- ✅ Start deploying

### Step 4: Set Environment Variables (3 minutes)

**After the initial deployment**, you need to cross-reference the URLs:

#### 4a. Update Backend CORS

1. In Render dashboard, go to **stafind-backend** service
2. Go to **Environment** tab
3. Add environment variable:
   - **Key:** `CORS_ALLOWED_ORIGINS`
   - **Value:** `https://stafind-frontend.onrender.com` (your frontend URL)
4. Click **"Save Changes"** (backend will redeploy)

#### 4b. Update Frontend API URL

1. Go to **stafind-frontend** service
2. Go to **Environment** tab
3. Add/update environment variable:
   - **Key:** `NEXT_PUBLIC_API_URL`
   - **Value:** `https://stafind-backend.onrender.com` (your backend URL)
4. Click **"Save Changes"** (frontend will redeploy)

### Step 5: Run Migrations (3 minutes)

**Option A: Using db-clean-go (Easiest)**

```bash
# Get database URL from Render dashboard (stafind-db → Connection → External Database URL)
export DATABASE_URL="postgresql://stafind:password@hostname.oregon-postgres.render.com/stafind"

# Create .env with Render database
cd backend
cat > .env << EOF
DB_PROVIDER=postgres
DATABASE_URL=$DATABASE_URL
EOF

# Run migrations
make db-clean-go
```

**Option B: Via Render Shell**

1. Go to backend service in Render dashboard
2. Click **"Shell"** tab
3. Run:
   ```bash
   cd backend
   go run cmd/flyway-cli/main.go migrate
   ```

### Step 6: Test Your App (1 minute)

1. Visit your frontend URL: `https://stafind-frontend.onrender.com`
2. Try logging in
3. Create an employee
4. Test the search functionality

**🎉 You're live!**

---

## 🔧 Troubleshooting

### "Backend build failed"

**Check:**
1. Render logs (Dashboard → Backend service → Logs)
2. Ensure `go.mod` and `go.sum` are in the repo
3. Verify build command path is correct

**Common fix:**
```bash
# Make sure go.mod is correct
cd backend
go mod tidy
git add go.mod go.sum
git commit -m "Update Go modules"
git push
```

### "Frontend build failed"

**Check:**
1. Frontend logs
2. Ensure `package.json` is correct
3. Node version compatibility

**Common fix:**
```bash
cd frontend
npm install
git add package-lock.json
git commit -m "Update dependencies"
git push
```

### "CORS errors in browser"

**Fix:**
1. Verify `CORS_ALLOWED_ORIGINS` in backend includes your frontend URL
2. Make sure it's `https://` not `http://`
3. Redeploy backend after changing

### "Can't connect to database"

**Fix:**
1. Check `DATABASE_URL` is set in backend environment
2. Verify database is running (Render dashboard → stafind-db)
3. Run migrations

### "Free database expired after 90 days"

**Options:**
1. **Upgrade to Render paid database** (~$7/month)
2. **Switch to Supabase** (forever free, see SUPABASE_SETUP.md):
   - Update `DATABASE_URL` in backend service
   - Set `DB_PROVIDER=supabase`
   - Redeploy

---

## 📊 Render Free Tier Limits

| Resource | Free Tier | Notes |
|----------|-----------|-------|
| **Web Services** | 750 hours/month/service | ~31 days |
| **Database** | 90 days free | Then $7/month or switch to Supabase |
| **Bandwidth** | Unlimited | No bandwidth charges |
| **Build Minutes** | Unlimited | No build time limits |
| **Services** | Unlimited | Deploy as many as you want |

**Free tier is generous!** Perfect for development and small apps.

---

## 🔄 Auto-Deployment

After initial setup, Render automatically deploys when you push to Git:

```bash
# Make changes
git add .
git commit -m "Add new feature"
git push

# Render automatically:
# 1. Detects the push
# 2. Builds frontend and backend
# 3. Deploys both services
# 4. No manual intervention needed!
```

---

## 🌐 Custom Domain

Want to use your own domain?

1. Go to your service in Render dashboard
2. Click **"Settings"** → **"Custom Domain"**
3. Add your domain (e.g., `app.yourdomain.com`)
4. Update DNS records as shown
5. SSL certificate automatically provisioned!

---

## 📈 Monitoring

Render provides:
- ✅ Real-time logs
- ✅ Metrics (CPU, Memory, Requests)
- ✅ Deploy history
- ✅ Health checks

Access from service dashboard → **Metrics** or **Logs** tab.

---

## 🚀 Migration from Development to Render

If you're currently running locally:

```bash
# 1. Commit your code
git add .
git commit -m "Prepare for deployment"
git push

# 2. Go to Render dashboard
# 3. New + → Blueprint
# 4. Connect repository
# 5. Done! Render reads render.yaml and deploys everything
```

---

## 💡 Pro Tips

### 1. Use Supabase for Database (Recommended)

Instead of Render PostgreSQL (90-day limit), use Supabase (forever free):

**Update backend environment in Render:**
```env
DB_PROVIDER=supabase
DATABASE_URL=postgresql://postgres.xxx:password@aws-0-us-east-1.pooler.supabase.com:6543/postgres
```

**Remove the database section from `render.yaml`** before deploying.

### 2. Set Up Staging Environment

Create a separate blueprint for staging:
- Use `develop` branch
- Separate database
- Test before production

### 3. Enable Deploy Previews

Render can create preview deployments for pull requests:
1. Settings → Deploy → Pull Request Previews
2. Enable for both services
3. Every PR gets a preview URL!

### 4. Monitor Your Free Hours

Check usage at dashboard.render.com → Account → Usage

**Free tier:** 750 hours/month per service
- Backend: ~31 days
- Frontend: ~31 days
- Both can run 24/7 on free tier!

---

## 🎉 Summary

**Deploy everything in 10 minutes:**

1. ✅ Push `render.yaml` to Git (already created for you!)
2. ✅ Connect repository in Render dashboard
3. ✅ Wait for deployment (~5 minutes)
4. ✅ Update environment variables with URLs
5. ✅ Run migrations
6. ✅ Test your app!

**Free tier includes:**
- Frontend hosting
- Backend hosting  
- PostgreSQL database (90 days)
- Automatic SSL
- Auto-deployments

**Recommended:** Use Supabase for the database (forever free) instead of Render PostgreSQL (90-day trial).

---

## 📚 Next Steps

After deploying:

1. ✅ Switch to Supabase for database (see [SUPABASE_SETUP.md](SUPABASE_SETUP.md))
2. ✅ Set up custom domain
3. ✅ Configure error tracking
4. ✅ Set up monitoring

**Need help?** Check [FULLSTACK_DEPLOYMENT.md](FULLSTACK_DEPLOYMENT.md) for detailed instructions!

**Happy deploying! 🚀**

