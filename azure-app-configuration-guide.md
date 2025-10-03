# Azure App Configuration Guide - Step by Step

## 🎯 Goal: Fix AADSTS500113 Error by Configuring Redirect URI

### Step 1: Access Azure Portal

1. **Go to [Azure Portal](https://portal.azure.com/)**
2. **Sign in** with your Microsoft 365 admin account
3. **Search for "App registrations"** in the top search bar
4. **Click on "App registrations"**

### Step 2: Find or Create Your App

#### Option A: If you already have an app
1. **Look for your existing app** (might be named "Resource Search Bot" or similar)
2. **Click on the app name** to open it

#### Option B: If you need to create a new app
1. **Click "New registration"**
2. **Fill in the details:**
   - **Name**: `Resource Search Bot`
   - **Supported account types**: `Accounts in this organizational directory only`
   - **Redirect URI**: Leave empty for now (we'll add it later)
3. **Click "Register"**

### Step 3: Configure Authentication (FIX THE ERROR)

1. **In your app, click "Authentication" in the left menu**
2. **Scroll down to "Redirect URIs" section**
3. **Click "Add a platform"**
4. **Select "Web"**
5. **Add the redirect URI based on your n8n setup:**

#### For Local n8n (localhost):
```
http://localhost:5678/rest/oauth2-credential/callback
```

#### For n8n Cloud:
```
https://your-account.app.n8n.cloud/rest/oauth2-credential/callback
```

#### For Self-hosted n8n:
```
https://your-domain.com/rest/oauth2-credential/callback
```

6. **Click "Save"**

### Step 4: Configure API Permissions

1. **Click "API permissions" in the left menu**
2. **Click "Add a permission"**
3. **Select "Microsoft Graph"**
4. **Choose "Application permissions"** (not Delegated)
5. **Add these permissions:**
   - `Channel.ReadBasic.All`
   - `ChannelMessage.Read.All`
   - `Chat.ReadWrite`
   - `Team.ReadBasic.All`
   - `User.Read.All`
6. **Click "Add permissions"**
7. **Click "Grant admin consent for [Your Organization]"**
8. **Confirm the permissions**

### Step 5: Create Client Secret

1. **Click "Certificates & secrets" in the left menu**
2. **Click "New client secret"**
3. **Add description**: `n8n integration secret`
4. **Set expiration**: `24 months` (recommended)
5. **Click "Add"**
6. **⚠️ IMPORTANT: Copy the secret value immediately** (it won't be shown again)

### Step 6: Get Your App Credentials

1. **Click "Overview" in the left menu**
2. **Copy these values:**
   - **Application (client) ID** - Copy this
   - **Directory (tenant) ID** - Copy this
   - **Client secret** - Use the value from Step 5

### Step 7: Configure n8n Credential

1. **Go to your n8n instance**
2. **Click "Credentials" in the left menu**
3. **Click "Create New"**
4. **Search for "Microsoft Teams"**
5. **Select "Microsoft Teams OAuth2 API"**
6. **Fill in the credentials:**
   ```
   Client ID: [Your Application ID from Step 6]
   Client Secret: [Your Client Secret from Step 6]
   Tenant ID: [Your Directory ID from Step 6]
   ```
7. **Click "Save"**

### Step 8: Test the Configuration

1. **Go back to your workflow**
2. **Click on the Teams Message Trigger node**
3. **Select the credential you just created**
4. **The "Error fetching options" should disappear**
5. **You should now see channel options in the dropdown**

## 🔧 Troubleshooting Common Issues

### Issue: "Still getting AADSTS500113 error"
**Solution**: 
- Double-check the redirect URI is exactly correct
- Make sure there are no extra spaces or characters
- Try using `http://localhost:5678` instead of `https://localhost:5678`

### Issue: "No channels showing in dropdown"
**Solution**:
- Make sure the app has the correct permissions
- Check that admin consent was granted
- Verify the app is installed in Teams

### Issue: "Permission denied errors"
**Solution**:
- Go back to API permissions
- Make sure all required permissions are added
- Grant admin consent again

### Issue: "App not found in Teams"
**Solution**:
- The app needs to be installed in Teams
- Go to Teams → Apps → Manage your apps → Upload custom app
- Upload the Teams app package

## 📋 Quick Checklist

- [ ] Azure app created/configured
- [ ] Redirect URI added correctly
- [ ] API permissions configured
- [ ] Admin consent granted
- [ ] Client secret created and copied
- [ ] n8n credential created
- [ ] Teams app installed (if using Teams integration)
- [ ] Workflow tested

## 🚀 Alternative: Use Webhook Instead

If you're still having issues with OAuth2, consider using the webhook approach:

1. **Replace Teams Message Trigger with Webhook node**
2. **Set webhook path**: `teams-resource-search`
3. **In Teams**: Channel → Connectors → Incoming Webhook
4. **Set webhook URL**: `http://localhost:5678/webhook/teams-resource-search`

This is actually simpler and more reliable for this use case!

## 📞 Need Help?

If you're still stuck:
1. **Check the exact error message** in n8n
2. **Verify the redirect URI** matches your n8n setup exactly
3. **Try the webhook alternative** - it's much simpler
4. **Check n8n logs** for more detailed error information
