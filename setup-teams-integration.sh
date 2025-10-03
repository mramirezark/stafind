#!/bin/bash

# Setup script for Teams Google Drive Resource Search Integration
# This script helps you set up the Teams integration for searching resources

echo "🚀 Setting up Teams Google Drive Resource Search Integration"
echo "============================================================"

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
mkdir -p n8n-workflows/teams-integration/credentials
mkdir -p n8n-workflows/teams-integration/data
mkdir -p n8n-workflows/teams-integration/logs

# Check if workflow file exists
if [ ! -f "n8n-workflows/Teams-GoogleDrive-Resource-Search.json" ]; then
    echo "❌ Teams workflow file not found: n8n-workflows/Teams-GoogleDrive-Resource-Search.json"
    exit 1
fi

echo "✅ Teams workflow file found"

# Create environment file template
echo "📝 Creating environment file template..."
cat > n8n-workflows/teams-integration/.env << EOF
# Teams Google Drive Resource Search Environment Variables

# Microsoft Teams Configuration
TEAMS_CHANNEL_ID=your_teams_channel_id_here
TEAMS_BOT_TOKEN=your_teams_bot_token_here
TEAMS_APP_ID=your_teams_app_id_here
TEAMS_APP_PASSWORD=your_teams_app_password_here

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

# Workflow Configuration
SEARCH_BATCH_SIZE=5
MIN_MATCH_SCORE=60
MAX_FILES_PER_SEARCH=50
SEARCH_TIMEOUT_SECONDS=300

# Language Configuration
DEFAULT_LANGUAGE=spanish
SUPPORTED_LANGUAGES=spanish,english
EOF

echo "✅ Environment file created: n8n-workflows/teams-integration/.env"

# Create credentials directory structure
echo "📁 Setting up credentials directory..."
mkdir -p n8n-workflows/teams-integration/credentials/microsoft-teams
mkdir -p n8n-workflows/teams-integration/credentials/google-drive
mkdir -p n8n-workflows/teams-integration/credentials/openai
mkdir -p n8n-workflows/teams-integration/credentials/google-sheets

# Create setup instructions
echo "📋 Creating setup instructions..."
cat > n8n-workflows/teams-integration/SETUP_INSTRUCTIONS.md << 'EOF'
# Teams Integration Setup Instructions

## Prerequisites
1. N8N instance (local or cloud)
2. Microsoft Teams workspace with admin permissions
3. Google Cloud Platform account
4. OpenAI API key
5. Google Drive folder with resumes
6. (Optional) Google Sheets for logging

## Step 1: Microsoft Teams App Setup

### Create Teams App
1. Go to [Microsoft Teams Developer Portal](https://dev.teams.microsoft.com/)
2. Sign in with your Microsoft 365 account
3. Create a new app:
   - Click "New app"
   - Choose "Build an app for your org"
   - Enter app name: "Resource Search Bot"
   - Upload app icons (192x192 and 32x32 pixels)

### Configure Bot
1. In your app, go to "Bot" section
2. Click "Create bot"
3. Choose "Create a new bot"
4. Configure bot settings:
   - Bot name: "Resource Search Assistant"
   - Bot handle: "resourcesearch"
   - Description: "AI-powered resource search assistant"
   - Upload bot icon

### Set Permissions
1. Go to "Permissions" section
2. Add required permissions:
   - `Microsoft Graph` → `Chat.ReadWrite`
   - `Microsoft Graph` → `ChannelMessage.Read.All`
   - `Microsoft Graph` → `ChannelMessage.Send`

### Get Credentials
1. Go to "Authentication" section
2. Copy the following values:
   - Application (client) ID
   - Directory (tenant) ID
   - Client secret (create new if needed)

### Install App
1. Go to "Publishing" section
2. Click "Publish your app"
3. Install app to your organization
4. Add app to desired Teams channels

## Step 2: Google Drive API Setup

### Enable APIs
1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create new project or select existing
3. Enable APIs:
   - Google Drive API
   - Google Sheets API (optional)

### Create Service Account
1. Go to "IAM & Admin" → "Service Accounts"
2. Create service account:
   - Name: "n8n-drive-access"
   - Description: "Service account for n8n Google Drive access"
3. Create and download JSON key file
4. Save as `credentials/google-drive-credentials.json`

### Configure Folder Access
1. Open your Google Drive
2. Navigate to resume folder
3. Right-click → "Share"
4. Add service account email (from JSON file)
5. Set permission to "Editor"
6. Copy folder ID from URL

## Step 3: OpenAI API Setup

1. Go to [OpenAI Platform](https://platform.openai.com/)
2. Create API key:
   - Go to "API Keys" section
   - Click "Create new secret key"
   - Copy the key
3. Add funds to your account
4. Set usage limits if needed

## Step 4: Google Sheets Setup (Optional)

1. Create new Google Sheet
2. Name it "Resource Search Log"
3. Create headers in first row:
   - Search ID, Original Request, User Name, Channel Name
   - Primary Skill, Experience Level, Years Experience Min
   - Total Files Processed, Matching Candidates
   - Candidates Found, Processed At
4. Copy sheet ID from URL
5. Share with service account (Editor permission)

## Step 5: Configure N8N

### Import Workflow
1. Open N8N
2. Go to "Workflows" → "Import from File"
3. Upload `Teams-GoogleDrive-Resource-Search.json`

### Configure Credentials

#### Microsoft Teams OAuth2 API
1. Create new credential: "Microsoft Teams OAuth2 API"
2. Enter:
   - Client ID: Your Teams app client ID
   - Client Secret: Your Teams app client secret
   - Tenant ID: Your directory tenant ID
   - Scopes: `https://graph.microsoft.com/.default`

#### Google Drive OAuth2 API
1. Create new credential: "Google Drive OAuth2 API"
2. Upload service account JSON file
3. Or use OAuth2 flow with your Google account

#### OpenAI API
1. Create new credential: "OpenAI API"
2. Enter your API key

#### Google Sheets OAuth2 API (Optional)
1. Create new credential: "Google Sheets OAuth2 API"
2. Use same credentials as Google Drive

### Update Workflow Variables
1. Update these variables in workflow nodes:
   - Google Drive folder ID
   - Teams channel ID
   - Backend API URL
   - Google Sheets ID (if using)

## Step 6: Test the Integration

### Test Teams Connection
1. Send test message in Teams channel
2. Check N8N execution logs
3. Verify bot responds

### Test Google Drive Access
1. Check if workflow can list folder contents
2. Verify file download works
3. Test text extraction

### Test AI Processing
1. Send Spanish query: "Dame recursos con Java"
2. Verify AI extracts criteria correctly
3. Check response generation

### End-to-End Test
1. Send complete query in Teams
2. Monitor workflow execution
3. Verify response in Teams
4. Check data in Google Sheets/backend

## Troubleshooting

### Common Issues

1. **Teams Bot Not Responding**:
   - Check bot is installed in channel
   - Verify OAuth2 credentials
   - Check webhook URLs

2. **Google Drive Access Denied**:
   - Verify service account permissions
   - Check folder sharing settings
   - Confirm API quotas

3. **AI Processing Errors**:
   - Check OpenAI API key and credits
   - Verify prompt formatting
   - Review response parsing logic

4. **No Files Found**:
   - Check Google Drive folder ID
   - Verify file formats are supported
   - Confirm folder contains resume files

### Debug Steps
1. Enable debug mode in N8N
2. Check execution logs for each node
3. Test individual nodes manually
4. Verify all credentials are working

## Security Best Practices

1. **API Keys**:
   - Store securely in N8N credentials
   - Rotate regularly
   - Monitor usage

2. **Access Control**:
   - Use least-privilege permissions
   - Limit folder access to necessary files
   - Monitor access logs

3. **Data Privacy**:
   - Ensure resume data is handled securely
   - Implement data retention policies
   - Follow GDPR requirements

## Performance Optimization

1. **Batch Processing**:
   - Adjust batch size based on performance
   - Process files in parallel when possible
   - Implement caching for processed files

2. **Rate Limiting**:
   - Respect API rate limits
   - Implement backoff strategies
   - Queue requests when needed

3. **Monitoring**:
   - Set up alerts for failures
   - Monitor execution times
   - Track API usage and costs
EOF

echo "✅ Setup instructions created: n8n-workflows/teams-integration/SETUP_INSTRUCTIONS.md"

# Create test script
echo "🧪 Creating test script..."
cat > n8n-workflows/teams-integration/test-integration.sh << 'EOF'
#!/bin/bash

echo "🧪 Testing Teams Google Drive Resource Search Integration"
echo "========================================================="

# Check if n8n is running
if ! curl -s http://localhost:5678 > /dev/null; then
    echo "❌ N8N is not running. Please start N8N first:"
    echo "   n8n start"
    echo "   or with Docker: docker run -it --rm --name n8n -p 5678:5678 n8nio/n8n"
    exit 1
fi

echo "✅ N8N is running"

# Check if workflow file exists
if [ ! -f "../Teams-GoogleDrive-Resource-Search.json" ]; then
    echo "❌ Workflow file not found"
    exit 1
fi

echo "✅ Workflow file found"

# Test workflow import
echo "📥 Testing workflow import..."
if n8n import:workflow --input=../Teams-GoogleDrive-Resource-Search.json; then
    echo "✅ Workflow imported successfully"
else
    echo "❌ Failed to import workflow"
    exit 1
fi

# Check environment file
if [ ! -f ".env" ]; then
    echo "❌ Environment file not found. Please configure .env file"
    exit 1
fi

echo "✅ Environment file found"

# Test API connections
echo "🔗 Testing API connections..."

# Test OpenAI API (if key is provided)
if [ -n "$OPENAI_API_KEY" ] && [ "$OPENAI_API_KEY" != "your_openai_api_key_here" ]; then
    echo "🧠 Testing OpenAI API..."
    if curl -s -H "Authorization: Bearer $OPENAI_API_KEY" \
            -H "Content-Type: application/json" \
            -d '{"model":"gpt-3.5-turbo","messages":[{"role":"user","content":"test"}],"max_tokens":5}' \
            https://api.openai.com/v1/chat/completions > /dev/null; then
        echo "✅ OpenAI API connection successful"
    else
        echo "⚠️  OpenAI API connection failed - check your API key"
    fi
else
    echo "⚠️  OpenAI API key not configured"
fi

# Test Google Drive API (if credentials exist)
if [ -f "credentials/google-drive-credentials.json" ]; then
    echo "📁 Google Drive credentials found"
    echo "✅ Google Drive credentials file exists"
else
    echo "⚠️  Google Drive credentials not found"
fi

# Test Teams configuration
if [ -n "$TEAMS_CHANNEL_ID" ] && [ "$TEAMS_CHANNEL_ID" != "your_teams_channel_id_here" ]; then
    echo "💬 Teams channel ID configured"
else
    echo "⚠️  Teams channel ID not configured"
fi

echo ""
echo "🎉 Integration test completed!"
echo ""
echo "Next steps:"
echo "1. 📖 Read the setup instructions: SETUP_INSTRUCTIONS.md"
echo "2. ⚙️  Configure your credentials in .env"
echo "3. 🔧 Set up Microsoft Teams app"
echo "4. 🚀 Test the complete workflow"
echo "5. 💬 Send test message in Teams: 'Dame recursos con Java'"
EOF

chmod +x n8n-workflows/teams-integration/test-integration.sh

echo "✅ Test script created: n8n-workflows/teams-integration/test-integration.sh"

# Create sample Teams messages
echo "💬 Creating sample Teams messages..."
cat > n8n-workflows/teams-integration/sample-messages.txt << 'EOF'
# Sample Teams Messages for Testing

## Basic Skill Searches
- "Dame los recursos con más experiencia que tengamos en Java"
- "Busca empleados con habilidades en Python"
- "¿Quién tiene experiencia en React?"
- "Encuentra desarrolladores con JavaScript"

## Experience Level Searches
- "Busca empleados senior con más de 5 años de experiencia"
- "Dame recursos junior con Python"
- "¿Hay desarrolladores lead con Node.js?"
- "Encuentra arquitectos con experiencia en microservicios"

## Multiple Skills Searches
- "Busca quien tenga React, Node.js y AWS"
- "¿Quién tiene Python, Django y PostgreSQL?"
- "Encuentra desarrolladores con Java, Spring y Docker"
- "Dame recursos con Angular, TypeScript y MongoDB"

## Specific Technology Searches
- "¿Quién tiene experiencia en Kubernetes?"
- "Busca empleados con Docker"
- "Encuentra desarrolladores con AWS Lambda"
- "¿Hay recursos con experiencia en GraphQL?"

## Experience Duration Searches
- "Dame recursos con más de 8 años de experiencia"
- "Busca empleados con menos de 3 años"
- "¿Quién tiene entre 5 y 10 años de experiencia?"
- "Encuentra desarrolladores experimentados"

## Location-based Searches
- "¿Hay recursos en Madrid con Java?"
- "Busca empleados remotos con Python"
- "Dame desarrolladores en Barcelona"
- "¿Quién está disponible para viajar?"

## Project-based Searches
- "¿Quién ha trabajado en proyectos de e-commerce?"
- "Busca empleados con experiencia en fintech"
- "¿Hay desarrolladores con proyectos de IoT?"
- "Dame recursos con experiencia en startups"

## Framework-specific Searches
- "¿Quién tiene experiencia en Spring Boot?"
- "Busca empleados con Angular"
- "Encuentra desarrolladores con Laravel"
- "¿Hay recursos con Django?"

## Database Searches
- "¿Quién tiene experiencia con PostgreSQL?"
- "Busca empleados con MongoDB"
- "Encuentra desarrolladores con Redis"
- "¿Hay recursos con Oracle?"

## Cloud Platform Searches
- "¿Quién tiene experiencia en AWS?"
- "Busca empleados con Azure"
- "Encuentra desarrolladores con Google Cloud"
- "¿Hay recursos con multi-cloud?"
EOF

echo "✅ Sample messages created: n8n-workflows/teams-integration/sample-messages.txt"

echo ""
echo "🎉 Teams Integration Setup Completed Successfully!"
echo ""
echo "📁 Files created:"
echo "   - Environment template: n8n-workflows/teams-integration/.env"
echo "   - Setup instructions: n8n-workflows/teams-integration/SETUP_INSTRUCTIONS.md"
echo "   - Test script: n8n-workflows/teams-integration/test-integration.sh"
echo "   - Sample messages: n8n-workflows/teams-integration/sample-messages.txt"
echo ""
echo "🚀 Next Steps:"
echo "1. 📖 Read: n8n-workflows/teams-integration/SETUP_INSTRUCTIONS.md"
echo "2. ⚙️  Configure: n8n-workflows/teams-integration/.env"
echo "3. 🔧 Set up Microsoft Teams app"
echo "4. 🚀 Start N8N: n8n start"
echo "5. 📥 Import workflow in N8N"
echo "6. 🧪 Test: cd n8n-workflows/teams-integration && ./test-integration.sh"
echo "7. 💬 Send test message in Teams"
echo ""
echo "📚 Documentation: n8n-workflows/TEAMS_INTEGRATION_GUIDE.md"
echo "💬 Sample queries: n8n-workflows/teams-integration/sample-messages.txt"
