#!/bin/bash

# Skills Enhancement Script
# This script will enhance your skills table with comprehensive modern technologies

echo "🚀 Starting Skills Table Enhancement..."

# Check if we have a database connection
if [ -z "$DATABASE_URL" ] && [ -z "$DB_HOST" ]; then
    echo "❌ Error: Database connection not configured"
    echo "Please set DATABASE_URL or DB_HOST environment variables"
    echo "Example: export DATABASE_URL='postgres://user:pass@localhost:5432/dbname'"
    exit 1
fi

# Run the enhancement SQL script
echo "📊 Running skills enhancement SQL..."
if [ ! -z "$DATABASE_URL" ]; then
    psql "$DATABASE_URL" -f run_skills_enhancement.sql
else
    psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -f run_skills_enhancement.sql
fi

if [ $? -eq 0 ]; then
    echo "✅ Skills table enhanced successfully!"
    echo ""
    echo "📈 What was added:"
    echo "   • 300+ new modern technologies"
    echo "   • Proper categorization for all skills"
    echo "   • Popularity scores based on 2025 trends"
    echo "   • Descriptions for major technologies"
    echo "   • Performance indexes"
    echo ""
    echo "🎯 Categories included:"
    echo "   • Programming Languages (40+)"
    echo "   • Frontend Frameworks (15+)"
    echo "   • Backend Frameworks (20+)"
    echo "   • Databases (25+)"
    echo "   • Cloud Platforms (15+)"
    echo "   • DevOps Tools (25+)"
    echo "   • API Technologies (10+)"
    echo "   • Architecture Patterns (15+)"
    echo "   • Testing Frameworks (15+)"
    echo "   • Mobile Development (10+)"
    echo "   • Data Science & AI (15+)"
    echo "   • Security (15+)"
    echo "   • Soft Skills (15+)"
    echo ""
    echo "🔧 Next steps:"
    echo "   1. Test your NER service with the enhanced database"
    echo "   2. Update your skill extraction to use the new categories"
    echo "   3. Consider implementing the DatabaseSkillExtractor service"
else
    echo "❌ Error: Failed to enhance skills table"
    exit 1
fi

echo "🎉 Enhancement complete!"
