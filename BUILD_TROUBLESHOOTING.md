# Build Troubleshooting Guide

Common build issues and their solutions for deploying StaffFind.

## Issue: Node Version Too New (24.x)

### Error Message
```
Using Node.js version 24.9.0
npm error Options:
npm error [--include <prod|dev|optional|peer>...]
```

### Root Cause

Render may use Node 24.x by default, which is too new and has compatibility issues with many packages.

### ✅ Solution

Lock to a stable LTS version (Node 20) in `render.yaml`:

```yaml
services:
  - type: web
    name: frontend
    envVars:
      - key: NODE_VERSION
        value: 20.18.0  # Lock to Node 20 LTS
```

**Already fixed in your `render.yaml`!** ✅

You can also create `.node-version` or `.nvmrc` files:

```bash
# Create .node-version
echo "20.18.0" > frontend/.node-version

# Or .nvmrc
echo "20.18.0" > frontend/.nvmrc
```

**Already created for you!** ✅

---

## Issue: TypeScript Not Found During Build

### Error Message
```
Please install typescript and @types/node by running:
	npm install --save-dev typescript @types/node
```

### Root Cause

devDependencies not being installed during build.

### ✅ Solution

Use `npm install` instead of `npm ci` with production flags:

**In `render.yaml`:**
```yaml
buildCommand: npm install && npm run build
```

This automatically installs devDependencies without special flags.

**Already fixed in your `render.yaml`!** ✅

---

## Issue: Module Not Found - Can't Resolve '@/components'

### Error Message
```
Module not found: Can't resolve '@/components'
./app/login/page.tsx
./app/page.tsx
```

### Root Cause

TypeScript path aliases (`@/components`) not configured properly for production builds.

### ✅ Solution

**1. Update `tsconfig.json` with explicit paths:**

```json
{
  "compilerOptions": {
    "baseUrl": ".",
    "paths": {
      "@/*": ["./*"],
      "@/components": ["./components"],
      "@/components/*": ["./components/*"],
      "@/lib/*": ["./lib/*"],
      "@/hooks/*": ["./hooks/*"],
      "@/services/*": ["./services/*"],
      "@/types/*": ["./types/*"],
      "@/utils/*": ["./utils/*"]
    }
  }
}
```

**2. Create `jsconfig.json` with same paths:**

```json
{
  "compilerOptions": {
    "baseUrl": ".",
    "paths": {
      "@/*": ["./*"],
      "@/components": ["./components"],
      "@/components/*": ["./components/*"]
    }
  }
}
```

**Already fixed in your project!** ✅

---

## Issue: Build Works Locally But Fails on Render

### Common Causes

1. **Different Node versions**
   - Solution: Specify Node version in `render.yaml`

2. **Missing files in Git**
   - Solution: Check `.gitignore` isn't ignoring required files

3. **Environment variables not set**
   - Solution: Set in Render dashboard

4. **Case-sensitive imports** (Linux vs macOS/Windows)
   - Solution: Ensure import paths match file names exactly

### ✅ Solutions

**1. Lock Node version in `package.json`:**

```json
{
  "engines": {
    "node": ">=18.0.0",
    "npm": ">=9.0.0"
  }
}
```

**Already configured!** ✅

**2. Verify all files are committed:**

```bash
git status
git add .
git commit -m "Add all files"
```

**3. Test locally with production build:**

```bash
cd frontend
NODE_ENV=production npm ci
npm run build
```

---

## Issue: Backend Build Fails

### Error: "Package Not Found"

**Solution:**

```bash
cd backend
go mod tidy
go mod download
git add go.mod go.sum
git commit -m "Update Go modules"
git push
```

### Error: "Cannot Find flyway_migrations"

**Solution:**

Ensure migrations are copied in build:

```yaml
# In render.yaml, make sure rootDir is set correctly
rootDir: backend
buildCommand: go build -o server cmd/server/main.go
```

The migrations folder will be included automatically with `rootDir`.

---

## Issue: Build Cache Warnings

### Warning Message
```
⚠ No build cache found. Please configure build caching for faster rebuilds.
```

### Impact
Builds will be slower but will still work.

### Solution (Optional)

**For Render, add to environment:**

```yaml
envVars:
  - key: NEXT_PRIVATE_BUILD_CACHE
    value: true
```

**Not critical - builds work without it.**

---

## Issue: Out of Memory During Build

### Error Message
```
FATAL ERROR: Reached heap limit Allocation failed - JavaScript heap out of memory
```

### Solution

**Increase Node memory in `render.yaml`:**

```yaml
envVars:
  - key: NODE_OPTIONS
    value: "--max-old-space-size=4096"
```

---

## Testing Builds Locally

### Test Backend Build

```bash
cd backend
go build -o /tmp/test-server cmd/server/main.go
echo "✅ Backend build successful"
```

### Test Frontend Build

```bash
cd frontend

# Test with production settings
NODE_ENV=production npm ci
npm run build

# Or use the test script
cd ..
./test-build.sh
```

---

## Build Command Reference

### Frontend

**Development:**
```bash
npm run dev
```

**Production build:**
```bash
npm ci
npm run build
npm start
```

**With explicit devDependencies:**
```bash
NPM_CONFIG_PRODUCTION=false npm ci
npm run build
```

### Backend

**Development:**
```bash
go run cmd/server/main.go
```

**Production build:**
```bash
go build -o server cmd/server/main.go
./server
```

**With migrations:**
```bash
go build -o server cmd/server/main.go
# Ensure flyway_migrations/ folder is in same directory
./server
```

---

## Render-Specific Fixes

### Frontend Build Command (Current - Correct)

```yaml
buildCommand: npm ci && npm run build
envVars:
  - key: NPM_CONFIG_PRODUCTION
    value: false  # ← This is the key fix!
```

### Backend Build Command (Current - Correct)

```yaml
rootDir: backend
buildCommand: go build -o server cmd/server/main.go
startCommand: ./server
```

---

## Common Render Environment Variables

### Frontend

```yaml
envVars:
  - key: NEXT_PUBLIC_API_URL
    value: https://your-backend.onrender.com
  - key: NPM_CONFIG_PRODUCTION
    value: false  # Install devDependencies
  - key: NODE_OPTIONS
    value: "--max-old-space-size=4096"  # If OOM errors
```

### Backend

```yaml
envVars:
  - key: PORT
    value: 8080
  - key: GIN_MODE
    value: release
  - key: DB_PROVIDER
    value: supabase
  - key: DATABASE_URL
    value: your-database-url
  - key: CORS_ALLOWED_ORIGINS
    value: https://your-frontend.onrender.com
```

---

## Verification Checklist

Before deploying, verify:

- [ ] `typescript` in `devDependencies` ✅
- [ ] `@types/node` in `devDependencies` ✅
- [ ] `NPM_CONFIG_PRODUCTION=false` in render.yaml ✅
- [ ] Path aliases in `tsconfig.json` ✅
- [ ] `jsconfig.json` created ✅
- [ ] `rootDir` set correctly in render.yaml ✅
- [ ] All files committed to Git
- [ ] Test build locally: `./test-build.sh`

---

## Quick Fix Commands

### If build still fails, try:

```bash
cd frontend

# Clean install
rm -rf node_modules package-lock.json
npm install

# Test build
npm run build

# If successful, commit package-lock.json
git add package-lock.json
git commit -m "Update package-lock.json"
git push
```

---

## Summary

**Your fixes are already in place:**

✅ `NPM_CONFIG_PRODUCTION=false` in `render.yaml`
✅ TypeScript and @types/node in `package.json`
✅ Path aliases configured in `tsconfig.json`
✅ `jsconfig.json` created

**Test locally:**
```bash
./test-build.sh
```

**Then deploy:**
```bash
git add .
git commit -m "Fix build configuration"
git push
```

Your build should now work on Render! 🎉
