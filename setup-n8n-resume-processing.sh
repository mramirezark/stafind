#!/bin/bash

# Setup script for N8N Resume Processing Workflow
# This script helps you set up the Google Drive resume processing workflow

echo "🚀 Setting up N8N Resume Processing Workflow"
echo "============================================="

# Check if n8n is available
if ! command -v n8n &> /dev/null; then
    echo "❌ N8N is not installed. Please install N8N first:"
    echo "   npm install -g n8n"
    echo "   or use Docker: docker run -it --rm --name n8n -p 5678:5678 n8nio/n8n"
    exit 1
fi

echo "✅ N8N is available"

# Create necessary directories
echo "📁 Creating directories..."
mkdir -p n8n-workflows/credentials
mkdir -p n8n-workflows/data

# Check if workflow file exists
if [ ! -f "n8n-workflows/Google-Drive-Resume-Processing.json" ]; then
    echo "❌ Workflow file not found: n8n-workflows/Google-Drive-Resume-Processing.json"
    exit 1
fi

echo "✅ Workflow file found"

# Create environment file template
echo "📝 Creating environment file template..."
cat > n8n-workflows/.env << EOF
# N8N Resume Processing Environment Variables

# Google Drive Configuration
GOOGLE_DRIVE_FOLDER_ID=your_google_drive_folder_id_here
GOOGLE_DRIVE_CREDENTIALS_PATH=./credentials/google-drive-credentials.json

# Google Sheets Configuration (Optional)
GOOGLE_SHEETS_ID=your_google_sheets_id_here

# OpenAI Configuration
OPENAI_API_KEY=your_openai_api_key_here

# Backend API Configuration
BACKEND_URL=http://localhost:8080
BACKEND_API_KEY=your_backend_api_key_here

# Slack Configuration (Optional)
SLACK_BOT_TOKEN=your_slack_bot_token_here
SLACK_CHANNEL=#resume-processing

# Workflow Configuration
SCHEDULE_CRON=0 */6 * * *
PROCESS_EXISTING_FILES=false
MAX_FILE_SIZE_MB=10
EOF

echo "✅ Environment file created: n8n-workflows/.env"

# Create credentials directory structure
echo "📁 Setting up credentials directory..."
mkdir -p n8n-workflows/credentials/google-drive
mkdir -p n8n-workflows/credentials/openai
mkdir -p n8n-workflows/credentials/slack

# Create setup instructions
echo "📋 Creating setup instructions..."
cat > n8n-workflows/SETUP_INSTRUCTIONS.md << 'EOF'
# N8N Resume Processing Setup Instructions

## Prerequisites
1. N8N instance (local or cloud)
2. Google Cloud Platform account
3. OpenAI API key
4. Google Drive folder with resumes
5. (Optional) Google Sheets for data storage
6. (Optional) Slack workspace for notifications

## Step 1: Google Drive API Setup

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select existing one
3. Enable Google Drive API:
   - Go to "APIs & Services" → "Library"
   - Search for "Google Drive API"
   - Click "Enable"
4. Create OAuth 2.0 credentials:
   - Go to "APIs & Services" → "Credentials"
   - Click "Create Credentials" → "OAuth 2.0 Client IDs"
   - Application type: "Web application"
   - Add authorized redirect URIs:
     - For local n8n: `http://localhost:5678/rest/oauth2-credential/callback`
     - For cloud n8n: `https://your-n8n-instance.com/rest/oauth2-credential/callback`
5. Download credentials JSON file
6. Save it as `credentials/google-drive-credentials.json`

## Step 2: OpenAI API Setup

1. Go to [OpenAI Platform](https://platform.openai.com/)
2. Create API key:
   - Go to "API Keys" section
   - Click "Create new secret key"
   - Copy the key
3. Add funds to your account (pay-per-use)
4. Add the API key to your environment file

## Step 3: Google Sheets Setup (Optional)

1. Create a new Google Sheet
2. Copy the sheet ID from the URL
3. Create headers in the first row:
   - Candidate Name, Email, Phone, Location
   - Seniority Level, Years Experience, Current Role
   - Technical Skills, Soft Skills, Programming Languages
   - Tools, Frameworks, Experience Summary
   - Projects Summary, Education, Certifications
   - Professional Summary, Google Drive ID, Filename
   - File Size, Created Time, Modified Time, Web View Link
   - Processing Timestamp

## Step 4: Slack Setup (Optional)

1. Go to [Slack API](https://api.slack.com/)
2. Create a new app
3. Add OAuth scopes: `chat:write`, `channels:read`
4. Install app to your workspace
5. Copy the Bot User OAuth Token

## Step 5: Configure the Workflow

1. Update `n8n-workflows/.env` with your actual values
2. Import the workflow:
   - Open n8N
   - Go to "Workflows" → "Import from File"
   - Upload `Google-Drive-Resume-Processing.json`
3. Configure credentials in n8N:
   - Google Drive OAuth2 API
   - OpenAI API
   - Google Sheets OAuth2 API (optional)
   - Slack OAuth2 API (optional)

## Step 6: Test the Workflow

1. Save the workflow
2. Click "Execute Workflow" to test manually
3. Check the execution log for any errors
4. Verify data appears in your storage destinations

## Troubleshooting

### Common Issues:
1. **Authentication Errors**: Check OAuth2 credentials and redirect URIs
2. **File Processing Failures**: Verify file permissions and formats
3. **AI Extraction Issues**: Check OpenAI API key and credits
4. **Data Storage Problems**: Verify permissions and endpoints

### Debug Mode:
Enable debug mode in n8N to see detailed execution logs.

## Security Notes:
- Store API keys securely
- Use least-privilege access
- Regularly rotate credentials
- Monitor API usage and costs
EOF

echo "✅ Setup instructions created: n8n-workflows/SETUP_INSTRUCTIONS.md"

# Create a simple test script
echo "🧪 Creating test script..."
cat > n8n-workflows/test-workflow.sh << 'EOF'
#!/bin/bash

echo "🧪 Testing N8N Resume Processing Workflow"
echo "=========================================="

# Check if n8n is running
if ! curl -s http://localhost:5678 > /dev/null; then
    echo "❌ N8N is not running. Please start N8N first:"
    echo "   n8n start"
    echo "   or with Docker: docker run -it --rm --name n8n -p 5678:5678 n8nio/n8n"
    exit 1
fi

echo "✅ N8N is running"

# Check if workflow file exists
if [ ! -f "Google-Drive-Resume-Processing.json" ]; then
    echo "❌ Workflow file not found"
    exit 1
fi

echo "✅ Workflow file found"

# Test workflow import
echo "📥 Testing workflow import..."
if n8n import:workflow --input=Google-Drive-Resume-Processing.json; then
    echo "✅ Workflow imported successfully"
else
    echo "❌ Failed to import workflow"
    exit 1
fi

echo "🎉 Test completed successfully!"
echo "Next steps:"
echo "1. Configure your credentials in N8N"
echo "2. Update the workflow variables"
echo "3. Test the workflow manually"
echo "4. Set up the schedule trigger"
EOF

chmod +x n8n-workflows/test-workflow.sh

echo "✅ Test script created: n8n-workflows/test-workflow.sh"

echo ""
echo "🎉 Setup completed successfully!"
echo ""
echo "Next steps:"
echo "1. 📖 Read the setup instructions: n8n-workflows/SETUP_INSTRUCTIONS.md"
echo "2. ⚙️  Configure your credentials in n8n-workflows/.env"
echo "3. 🚀 Start N8N: n8n start"
echo "4. 📥 Import the workflow in N8N"
echo "5. 🧪 Test the workflow: cd n8n-workflows && ./test-workflow.sh"
echo ""
echo "📚 Documentation: n8n-workflows/README.md"
echo "🔧 Configuration: n8n-workflows/.env"
echo "📋 Instructions: n8n-workflows/SETUP_INSTRUCTIONS.md"
