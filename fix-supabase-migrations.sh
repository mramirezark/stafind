#!/bin/bash

# Fix Supabase Migration Issues
# This script cleans your Supabase database and runs migrations fresh

echo "🔧 Fixing Supabase Migration Issues"
echo "===================================="
echo ""

# Your Supabase connection string (with URL-encoded password)
DATABASE_URL="postgresql://postgres.sijmfiggqrornhubfsyf:73GW8mZ%26QBaDuvu@aws-0-us-east-1.pooler.supabase.com:6543/postgres"

echo "⚠️  WARNING: This will DELETE ALL DATA from your Supabase database!"
echo "Are you sure you want to continue? (type 'yes' to confirm): "
read confirmation

if [ "$confirmation" != "yes" ]; then
    echo "❌ Cancelled."
    exit 0
fi

# Create temporary .env
cd backend
cat > .env << EOF
DB_PROVIDER=supabase
DATABASE_URL=$DATABASE_URL
EOF

echo ""
echo "🗑️  Cleaning Supabase database..."
go run cmd/db-clean/main.go << EOF
yes
EOF

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Supabase database cleaned and migrated successfully!"
    echo ""
    echo "Next steps:"
    echo "1. Push your code to Git (if not already done)"
    echo "2. Redeploy on Render"
    echo "3. Your backend should start successfully!"
else
    echo ""
    echo "❌ Failed to clean database. Check the error above."
    exit 1
fi

# Clean up
rm .env

echo ""
echo "🎉 All done!"

