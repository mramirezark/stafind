#!/bin/bash

# Test Build Script - Verify everything builds correctly before deploying

set -e

echo "🔨 Testing StaffFind Build Process"
echo "=================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test Backend Build
echo "📦 Testing Backend Build..."
cd backend
if go build -o /tmp/stafind-server cmd/server/main.go; then
    echo -e "${GREEN}✅ Backend build successful${NC}"
    rm -f /tmp/stafind-server
else
    echo -e "${RED}❌ Backend build failed${NC}"
    exit 1
fi
cd ..

echo ""

# Test Frontend Build
echo "📦 Testing Frontend Build..."
cd frontend

# Check if node_modules exists
if [ ! -d "node_modules" ]; then
    echo -e "${YELLOW}⚠️  node_modules not found. Running npm install...${NC}"
    npm install
fi

# Test the build
if npm run build; then
    echo -e "${GREEN}✅ Frontend build successful${NC}"
else
    echo -e "${RED}❌ Frontend build failed${NC}"
    echo ""
    echo "Common fixes:"
    echo "1. Make sure all imports use correct paths"
    echo "2. Check tsconfig.json has correct path aliases"
    echo "3. Run 'npm install' to ensure all dependencies are installed"
    exit 1
fi

cd ..

echo ""
echo -e "${GREEN}🎉 All builds successful!${NC}"
echo ""
echo "Next steps:"
echo "1. Commit your changes: git add . && git commit -m 'Fix build'"
echo "2. Push to Git: git push"
echo "3. Deploy to Render!"

