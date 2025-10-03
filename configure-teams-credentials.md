# Quick Teams Credentials Configuration for n8n

## Immediate Steps to Fix Your Teams Message Trigger

### 1. First, you need these values from Azure Portal:

1. **Go to Azure Portal** → **App registrations** → **Your app**
2. **Copy these values:**
   - Application (client) ID
   - Directory (tenant) ID  
   - Client secret (create new if needed)

### 2. Create Credential in n8n:

1. **In n8n, go to "Credentials"**
2. **Click "Create New"**
3. **Select "Microsoft Teams OAuth2 API"**
4. **Fill in:**
   ```
   Client ID: [Your Application ID from Azure]
   Client Secret: [Your Client Secret from Azure]
   Tenant ID: [Your Directory ID from Azure]
   ```
5. **Click "Save"**

### 3. Configure the Teams Message Trigger Node:

1. **In your workflow, click on the Teams Message Trigger node**
2. **Set these values:**
   ```
   Credential to connect with: [Select the credential you just created]
   Resource: Message
   Operation: Create  
   Chat Name or ID: [This will populate once credential is set]
   Message Type: Text
   Message: [Leave empty - will receive incoming messages]
   ```

### 4. Quick Test:

1. **Save and activate your workflow**
2. **In Teams, mention your bot: `@Resource Search Bot hola`**
3. **Check if n8n receives the message**

## If you don't have a Teams app yet:

### Quick Teams App Creation:

1. **Go to [Teams App Studio](https://appstudio.teams.microsoft.com/)**
2. **Click "New app"**
3. **Fill basic info:**
   - Name: `Resource Search Bot`
   - Description: `AI resource search assistant`
4. **Go to "Bot" tab → "Create bot"**
5. **Use your Azure App ID from Step 1**
6. **Download and install the app to Teams**

## Alternative: Use Webhook Instead

If Teams app setup is complex, you can use a webhook trigger:

1. **Replace Teams Message Trigger with "Webhook" node**
2. **Set webhook URL: `/webhook/teams-resource-search`**
3. **Use this URL in Teams webhook connector**
4. **Configure Teams connector to send messages to this webhook**

## Common Issues & Solutions:

### ❌ "Error fetching options from Microsoft Teams"
**Solution:** Credential not configured or app not installed in Teams

### ❌ "Select Credential" dropdown empty  
**Solution:** Create Microsoft Teams OAuth2 API credential first

### ❌ Bot not responding in Teams
**Solution:** Check if workflow is active and webhook is accessible

### ❌ Permission denied
**Solution:** Grant admin consent in Azure Portal → API permissions

## Need Help?

1. **Check the detailed guide:** `teams-app-setup-guide.md`
2. **Run setup script:** `./setup-teams-integration.sh`
3. **Test integration:** Use the sample messages provided

## Quick Commands:

```bash
# Run the complete setup
./setup-teams-integration.sh

# Test the integration  
cd n8n-workflows/teams-integration
./test-integration.sh
```
