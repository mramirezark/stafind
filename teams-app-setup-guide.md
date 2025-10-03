# Microsoft Teams App Setup Guide

## Step 1: Create Teams App in Azure Portal

### 1.1 Access Azure Portal
1. Go to [Azure Portal](https://portal.azure.com/)
2. Sign in with your Microsoft 365 admin account
3. Search for "App registrations" in the search bar

### 1.2 Create New App Registration
1. Click "New registration"
2. Fill in the details:
   - **Name**: `Resource Search Bot`
   - **Supported account types**: `Accounts in this organizational directory only`
   - **Redirect URI**: `Web` → `https://your-n8n-instance.com/rest/oauth2-credential/callback`
3. Click "Register"

### 1.3 Get Application Credentials
1. Copy these values (you'll need them for n8n):
   - **Application (client) ID**: Found in "Overview" section
   - **Directory (tenant) ID**: Found in "Overview" section

### 1.4 Create Client Secret
1. Go to "Certificates & secrets"
2. Click "New client secret"
3. Add description: `n8n integration secret`
4. Set expiration: `24 months` (recommended)
5. Click "Add"
6. **IMPORTANT**: Copy the secret value immediately (it won't be shown again)

## Step 2: Configure API Permissions

### 2.1 Add Microsoft Graph Permissions
1. Go to "API permissions"
2. Click "Add a permission"
3. Select "Microsoft Graph"
4. Choose "Application permissions" (not Delegated)
5. Add these permissions:
   - `Channel.ReadBasic.All`
   - `ChannelMessage.Read.All`
   - `ChannelMessage.Send`
   - `Chat.ReadWrite`
   - `Team.ReadBasic.All`
6. Click "Add permissions"

### 2.2 Grant Admin Consent
1. Click "Grant admin consent for [Your Organization]"
2. Confirm the permissions

## Step 3: Create Teams App Manifest

### 3.1 Download App Studio
1. Go to [Teams App Studio](https://appstudio.teams.microsoft.com/)
2. Sign in with the same Microsoft 365 account

### 3.2 Create New App
1. Click "New app"
2. Fill in basic info:
   - **App name**: `Resource Search Bot`
   - **Short description**: `AI-powered resource search assistant`
   - **Long description**: `Search for team members with specific skills and experience using natural language queries`

### 3.3 Configure App Details
1. Go to "App details" tab
2. Upload app icons (192x192 and 32x32 pixels)
3. Set accent color: `#0078D4` (Teams blue)

### 3.4 Configure Bot
1. Go to "Bot" tab
2. Click "Create bot"
3. Fill in details:
   - **Bot name**: `Resource Search Assistant`
   - **Bot handle**: `resourcesearch`
   - **Description**: `AI-powered resource search assistant`
   - **App ID**: Use the Application ID from Azure (Step 1.3)
4. Upload bot icon (same as app icons)

### 3.5 Configure Permissions
1. Go to "Permissions" tab
2. Add scopes:
   - `team`
   - `groupchat`
3. Add permissions:
   - `Microsoft Graph` → `ChannelMessage.Read.All`
   - `Microsoft Graph` → `ChannelMessage.Send`
   - `Microsoft Graph` → `Chat.ReadWrite`

## Step 4: Configure Messaging Endpoint

### 4.1 Set Up Webhook (For n8n)
1. Go to "Messaging" tab
2. Set **Messaging endpoint**: `https://your-n8n-instance.com/webhook/teams-messages`
3. Set **Bot framework app ID**: Use the Application ID from Azure

## Step 5: Install App to Teams

### 5.1 Download App Package
1. Go to "Test and distribute" tab
2. Click "Download"
3. Save the `.zip` file

### 5.2 Install to Teams
1. Open Microsoft Teams
2. Go to "Apps" → "Manage your apps"
3. Click "Upload an app"
4. Select "Upload a custom app"
5. Upload the downloaded `.zip` file
6. Click "Add"

### 5.3 Add to Channel
1. Go to the channel where you want to test
2. Click the "..." menu → "Manage team"
3. Go to "Apps" tab
4. Find your app and click "Add"

## Step 6: Get Channel Information

### 6.1 Get Channel ID
1. In Teams, go to your channel
2. Click the channel name at the top
3. Copy the channel ID from the URL or use the Teams web interface
4. The URL format is: `https://teams.microsoft.com/l/channel/[CHANNEL_ID]/General`

### 6.2 Test Bot Installation
1. In the channel, type `@Resource Search Bot`
2. The bot should appear in the mention list
3. Try sending a message to the bot

## Step 7: Configure n8n

### 7.1 Create Teams Credential in n8n
1. In n8n, go to "Credentials"
2. Click "Create New"
3. Select "Microsoft Teams OAuth2 API"
4. Fill in the details:
   - **Client ID**: Application ID from Azure (Step 1.3)
   - **Client Secret**: Secret value from Azure (Step 1.4)
   - **Tenant ID**: Directory ID from Azure (Step 1.3)
5. Click "Save"

### 7.2 Configure Teams Message Trigger Node
1. Open your workflow
2. Click on the "Teams Message Trigger" node
3. Configure:
   - **Credential**: Select the credential you just created
   - **Resource**: `Message`
   - **Operation**: `Create`
   - **Chat Name or ID**: Select your channel (should now work)
   - **Message Type**: `Text`
   - **Message**: Leave empty (will be populated by incoming messages)

## Step 8: Test the Integration

### 8.1 Test Message Flow
1. Save your n8n workflow
2. Activate the workflow
3. In Teams, send a message to the bot
4. Check n8n execution logs

### 8.2 Sample Test Messages
Try these messages in Teams:
- `@Resource Search Bot Dame recursos con Java`
- `@Resource Search Bot Busca empleados con Python`
- `@Resource Search Bot ¿Quién tiene experiencia en React?`

## Troubleshooting

### Common Issues:

1. **"Error fetching options from Microsoft Teams"**
   - Check if credentials are correctly configured
   - Verify the app has proper permissions
   - Ensure the app is installed in Teams

2. **Bot not responding**
   - Check if the workflow is active in n8n
   - Verify the webhook endpoint is accessible
   - Check n8n execution logs for errors

3. **Permission denied errors**
   - Ensure admin consent was granted in Azure
   - Check if the app has the required Graph permissions
   - Verify the bot is added to the channel

4. **Messages not triggering workflow**
   - Check if the channel ID is correct
   - Verify the message format
   - Ensure the workflow trigger is properly configured

### Debug Steps:
1. Check Azure Portal for any permission issues
2. Verify Teams app installation
3. Test webhook connectivity
4. Review n8n execution logs
5. Check network connectivity between Teams and n8n

## Security Notes:
- Keep client secrets secure
- Use HTTPS for webhook endpoints
- Regularly rotate secrets
- Monitor app permissions
- Review access logs regularly
