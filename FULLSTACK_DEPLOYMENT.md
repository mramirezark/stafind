# Full-Stack Deployment Guide - Deploy Everything Together

Deploy your complete StaffFind application (frontend + backend + database) on a single platform for free!

## 🎯 Best Free Full-Stack Platforms

| Platform | Free Tier | Best For | Difficulty |
|----------|-----------|----------|------------|
| **Render** | ✅ 750 hrs/month | Easiest setup | ⭐ Easy |
| **Railway** | ✅ $5 credit/month | Best DX | ⭐⭐ Easy |
| **Fly.io** | ✅ 3 VMs free | Most powerful | ⭐⭐ Medium |
| **Oracle Cloud** | ✅ Forever free VPS | Full control | ⭐⭐⭐ Advanced |

---

## Option 1: Render (Recommended) ⭐

**Why Render:**
- ✅ True free tier (750 hours/month per service)
- ✅ Deploy frontend, backend, and database
- ✅ No credit card required
- ✅ Auto-deploy from Git
- ✅ Free PostgreSQL database (90 days, then expires)

### Deploy to Render

#### 1. Create `render.yaml` in project root

```yaml
services:
  # Backend Service
  - type: web
    name: stafind-backend
    env: go
    region: oregon
    buildCommand: cd backend && go build -o server cmd/server/main.go
    startCommand: cd backend && ./server
    envVars:
      - key: PORT
        value: 8080
      - key: GIN_MODE
        value: release
      - key: DB_PROVIDER
        value: postgres
      - key: DATABASE_URL
        fromDatabase:
          name: stafind-db
          property: connectionString
      - key: CORS_ALLOWED_ORIGINS
        value: https://stafind-frontend.onrender.com
    healthCheckPath: /health

  # Frontend Service
  - type: web
    name: stafind-frontend
    env: node
    region: oregon
    buildCommand: cd frontend && npm install && npm run build
    startCommand: cd frontend && npm start
    envVars:
      - key: NODE_ENV
        value: production
      - key: NEXT_PUBLIC_API_URL
        value: https://stafind-backend.onrender.com

databases:
  # PostgreSQL Database
  - name: stafind-db
    databaseName: stafind
    user: stafind
    region: oregon
    plan: free
    ipAllowList: []
```

#### 2. Deploy via Render Dashboard

1. Go to [dashboard.render.com](https://dashboard.render.com)
2. Click **"New"** → **"Blueprint"**
3. Connect your Git repository
4. Render will detect `render.yaml` and create all services
5. Click **"Apply"**

**Done!** Your app will be live at:
- Frontend: `https://stafind-frontend.onrender.com`
- Backend: `https://stafind-backend.onrender.com`

#### 3. Run Migrations

After deployment, run migrations:

```bash
# Get database URL from Render dashboard
export DATABASE_URL="postgresql://user:password@hostname/database"

# Run migrations locally against Render database
cd backend
go run cmd/flyway-cli/main.go migrate
```

---

## Option 2: Railway (Best Developer Experience) ⭐

**Why Railway:**
- ✅ $5 free credit per month (~500 hours)
- ✅ Excellent CLI
- ✅ Built-in PostgreSQL
- ✅ Automatic HTTPS
- ✅ Environment variable management

### Deploy to Railway

#### 1. Install Railway CLI

```bash
npm i -g @railway/cli
railway login
```

#### 2. Create `railway.json` in project root

```json
{
  "$schema": "https://railway.app/railway.schema.json",
  "build": {
    "builder": "NIXPACKS"
  },
  "deploy": {
    "numReplicas": 1,
    "sleepApplication": false,
    "restartPolicyType": "ON_FAILURE",
    "restartPolicyMaxRetries": 10
  }
}
```

#### 3. Create Nixpacks config files

**`backend/nixpacks.toml`:**

```toml
[phases.setup]
nixPkgs = ["go_1_21"]

[phases.build]
cmds = ["go build -o server cmd/server/main.go"]

[phases.deploy]
cmd = "./server"
```

**`frontend/nixpacks.toml`:**

```toml
[phases.setup]
nixPkgs = ["nodejs-18_x"]

[phases.build]
cmds = ["npm install", "npm run build"]

[phases.deploy]
cmd = "npm start"
```

#### 4. Deploy

```bash
# Initialize Railway project
railway init

# Add PostgreSQL
railway add

# Create backend service
railway service create backend
cd backend
railway up

# Create frontend service
railway service create frontend
cd ../frontend
railway up

# Link services and set environment variables in Railway dashboard
```

#### 5. Set Environment Variables

In Railway dashboard:

**Backend:**
```env
PORT=8080
GIN_MODE=release
DB_PROVIDER=postgres
DATABASE_URL=${{Postgres.DATABASE_URL}}
CORS_ALLOWED_ORIGINS=${{Frontend.RAILWAY_PUBLIC_DOMAIN}}
```

**Frontend:**
```env
NEXT_PUBLIC_API_URL=${{Backend.RAILWAY_PUBLIC_DOMAIN}}
```

---

## Option 3: Fly.io (Most Powerful) ⭐⭐

**Why Fly.io:**
- ✅ 3 small VMs free (256MB RAM each)
- ✅ Full Docker support
- ✅ Global deployment
- ✅ Persistent volumes

### Deploy to Fly.io with Docker Compose

#### 1. Install flyctl

```bash
curl -L https://fly.io/install.sh | sh
flyctl auth login
```

#### 2. Create `fly.toml` in project root

```toml
app = "stafind"
primary_region = "iad"

[build]
  dockerfile = "Dockerfile.production"

[env]
  PORT = "8080"
  GIN_MODE = "release"

[[services]]
  internal_port = 8080
  protocol = "tcp"

  [[services.ports]]
    handlers = ["http"]
    port = 80
    force_https = true

  [[services.ports]]
    handlers = ["tls", "http"]
    port = 443

[http_service]
  internal_port = 8080
  force_https = true
  auto_stop_machines = true
  auto_start_machines = true
  min_machines_running = 0

[[vm]]
  cpu_kind = "shared"
  cpus = 1
  memory_mb = 256
```

#### 3. Create Production Dockerfile

**`Dockerfile.production`:**

```dockerfile
# Build backend
FROM golang:1.21-alpine AS backend-builder
WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN go build -o server cmd/server/main.go

# Build frontend
FROM node:18-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# Production image
FROM alpine:latest
WORKDIR /app

# Install Node.js for Next.js
RUN apk add --no-cache nodejs npm

# Copy backend
COPY --from=backend-builder /app/backend/server /app/backend/server
COPY --from=backend-builder /app/backend/flyway_migrations /app/backend/flyway_migrations

# Copy frontend
COPY --from=frontend-builder /app/frontend/.next /app/frontend/.next
COPY --from=frontend-builder /app/frontend/node_modules /app/frontend/node_modules
COPY --from=frontend-builder /app/frontend/package.json /app/frontend/package.json
COPY --from=frontend-builder /app/frontend/public /app/frontend/public

# Install nginx to serve both
RUN apk add --no-cache nginx

# Nginx config
COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 8080

# Start script
COPY start.sh /app/start.sh
RUN chmod +x /app/start.sh

CMD ["/app/start.sh"]
```

#### 4. Create start script

**`start.sh`:**

```bash
#!/bin/sh

# Start backend
cd /app/backend
./server &

# Start frontend
cd /app/frontend
npm start &

# Start nginx
nginx -g 'daemon off;'
```

#### 5. Create nginx config

**`nginx.conf`:**

```nginx
events {
    worker_connections 1024;
}

http {
    server {
        listen 8080;

        # Frontend
        location / {
            proxy_pass http://localhost:3000;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection 'upgrade';
            proxy_set_header Host $host;
            proxy_cache_bypass $http_upgrade;
        }

        # Backend API
        location /api/ {
            proxy_pass http://localhost:8080;
            proxy_http_version 1.1;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
        }
    }
}
```

#### 6. Deploy

```bash
# Create Fly app
flyctl apps create stafind

# Create Postgres database
flyctl postgres create --name stafind-db

# Attach database to app
flyctl postgres attach stafind-db

# Deploy
flyctl deploy

# Run migrations
flyctl ssh console
cd /app/backend
./server migrate
```

---

## Option 4: Render with Docker (Simplest Full-Stack) ⭐

Use your existing `docker-compose.yml`!

#### 1. Create `Dockerfile.render` in project root

```dockerfile
# Multi-stage build for both frontend and backend

# Backend build
FROM golang:1.21-alpine AS backend
WORKDIR /app/backend
COPY backend/ ./
RUN go build -o server cmd/server/main.go

# Frontend build
FROM node:18-alpine AS frontend
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
ENV NEXT_PUBLIC_API_URL=http://localhost:8080
RUN npm run build

# Final image
FROM node:18-alpine
WORKDIR /app

# Install Go runtime
RUN apk add --no-cache go

# Copy built files
COPY --from=backend /app/backend/server /app/backend/
COPY --from=backend /app/backend/flyway_migrations /app/backend/flyway_migrations
COPY --from=frontend /app/frontend/.next /app/frontend/.next
COPY --from=frontend /app/frontend/node_modules /app/frontend/node_modules
COPY --from=frontend /app/frontend/package.json /app/frontend/
COPY --from=frontend /app/frontend/public /app/frontend/public

# Startup script
RUN echo '#!/bin/sh\n\
cd /app/backend && ./server &\n\
cd /app/frontend && npm start\n\
' > /app/start.sh && chmod +x /app/start.sh

EXPOSE 8080 3000

CMD ["/app/start.sh"]
```

#### 2. Deploy to Render

1. Go to [dashboard.render.com](https://dashboard.render.com)
2. New → Web Service
3. Connect Git repository
4. Select **Docker** runtime
5. Set:
   - **Dockerfile path:** `Dockerfile.render`
   - **Port:** 3000
6. Add PostgreSQL database
7. Deploy!

---

## Option 5: Oracle Cloud (Forever Free VPS) 💎

**Why Oracle Cloud:**
- ✅ **Forever free** (not a trial)
- ✅ 2 VMs with 1GB RAM each
- ✅ 200GB storage
- ✅ Full control

### Deploy to Oracle Cloud

#### 1. Create Oracle Cloud Account

1. Sign up at [cloud.oracle.com](https://cloud.oracle.com)
2. Create Compute Instance (Always Free tier)
3. Choose Ubuntu 22.04
4. Save SSH key

#### 2. Connect and Setup

```bash
# SSH into your VM
ssh ubuntu@your-vm-ip

# Install Docker
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker ubuntu

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Clone your repository
git clone https://github.com/yourusername/stafind.git
cd stafind
```

#### 3. Update docker-compose for production

**`docker-compose.production.yml`:**

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: stafind
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped

  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: postgres
      DB_PASSWORD: ${DB_PASSWORD}
      DB_NAME: stafind
      PORT: 8080
      GIN_MODE: release
    depends_on:
      - postgres
    restart: unless-stopped

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    environment:
      NEXT_PUBLIC_API_URL: http://your-vm-ip:8080
    ports:
      - "80:3000"
      - "8080:8080"
    depends_on:
      - backend
    restart: unless-stopped

volumes:
  postgres_data:
```

#### 4. Deploy

```bash
# Create .env file
echo "DB_PASSWORD=your_secure_password" > .env

# Start services
docker-compose -f docker-compose.production.yml up -d

# Run migrations
docker-compose exec backend go run cmd/flyway-cli/main.go migrate

# View logs
docker-compose logs -f
```

#### 5. Setup Nginx reverse proxy (optional but recommended)

```bash
sudo apt install nginx

# Configure nginx
sudo nano /etc/nginx/sites-available/stafind
```

**Nginx config:**

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }

    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

```bash
sudo ln -s /etc/nginx/sites-available/stafind /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

---

## Comparison Matrix

| Feature | Render | Railway | Fly.io | Oracle Cloud |
|---------|--------|---------|--------|--------------|
| **Free Tier** | 750 hrs/service | $5 credit/mo | 3 VMs | Forever |
| **Setup Time** | 10 min | 15 min | 20 min | 30 min |
| **Database** | ✅ Free 90 days | ✅ Included | ⚠️ Extra cost | ✅ Self-hosted |
| **Auto-deploy** | ✅ Yes | ✅ Yes | ✅ Yes | ❌ Manual |
| **SSL/HTTPS** | ✅ Auto | ✅ Auto | ✅ Auto | ⚠️ Setup required |
| **Docker Support** | ✅ Yes | ✅ Yes | ✅ Yes | ✅ Yes |
| **Difficulty** | ⭐ Easy | ⭐ Easy | ⭐⭐ Medium | ⭐⭐⭐ Advanced |

---

## 🏆 Recommendation

**For Quick Start:** → **Render** (easiest, everything in one place)

**For Best DX:** → **Railway** (excellent developer experience)

**For Long-term Free:** → **Oracle Cloud** (forever free, full control)

---

## 📋 Quick Deploy Commands

### Render (Easiest)

```bash
# Just create render.yaml and push to Git
git add render.yaml
git commit -m "Add Render config"
git push

# Then connect in Render dashboard
```

### Railway

```bash
npm i -g @railway/cli
railway login
railway init
railway up
```

### Fly.io

```bash
curl -L https://fly.io/install.sh | sh
flyctl launch
flyctl deploy
```

### Oracle Cloud

```bash
ssh ubuntu@your-vm-ip
git clone your-repo
docker-compose up -d
```

---

## 💰 Cost Comparison (Monthly)

| Platform | Free Tier | After Free Tier |
|----------|-----------|-----------------|
| **Render** | $0 (750 hrs) | ~$7/service |
| **Railway** | $0 ($5 credit) | Pay as you go |
| **Fly.io** | $0 (3 VMs) | ~$2/VM extra |
| **Oracle Cloud** | **$0 forever** | $0 (always free) |

---

## 🎯 Next Steps

1. Choose your platform
2. Follow the specific guide above
3. Deploy!
4. Celebrate! 🎉

All platforms support **deploying frontend + backend + database together** for **free**!

Need help with a specific platform? Let me know! 🚀

