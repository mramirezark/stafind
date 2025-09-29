# Microsoft Teams Integration Setup Guide

## Option 1: Microsoft Graph API (Real-time Monitoring)

### Prerequisites
- Azure subscription
- Teams admin permissions
- Node.js installed

### Step 1: Azure App Registration

1. **Go to Azure Portal** → **Azure Active Directory** → **App registrations**
2. **Click "New registration"**
3. **Fill in details:**
   - Name: "Stafind Teams Integration"
   - Supported account types: "Accounts in this organizational directory only"
   - Redirect URI: "Web" → "http://localhost:3000/callback"
4. **Click "Register"**
5. **Copy the "Application (client) ID" and "Directory (tenant) ID"**

### Step 2: Configure API Permissions

1. **Go to "API permissions"**
2. **Click "Add a permission"**
3. **Select "Microsoft Graph"**
4. **Add these permissions:**
   - `Channel.ReadBasic.All`
   - `ChannelMessage.Read.All`
   - `Chat.Read.All`
   - `User.Read.All`
5. **Click "Grant admin consent"**

### Step 3: Create Client Secret

1. **Go to "Certificates & secrets"**
2. **Click "New client secret"**
3. **Add description: "Stafind Teams Secret"**
4. **Set expiration (recommend 24 months)**
5. **Copy the secret value (you won't see it again)**

### Step 4: Install Dependencies

```bash
npm install @azure/msal-node @microsoft/microsoft-graph-client
```

### Step 5: Create Graph API Client

```javascript
const { Client } = require('@microsoft/microsoft-graph-client');
const { ConfidentialClientApplication } = require('@azure/msal-node');

const msalConfig = {
  auth: {
    clientId: 'YOUR_CLIENT_ID',
    clientSecret: 'YOUR_CLIENT_SECRET',
    authority: 'https://login.microsoftonline.com/YOUR_TENANT_ID'
  }
};

const cca = new ConfidentialClientApplication(msalConfig);

async function getAccessToken() {
  const clientCredentialRequest = {
    scopes: ['https://graph.microsoft.com/.default'],
  };

  try {
    const response = await cca.acquireTokenSilent(clientCredentialRequest);
    return response.accessToken;
  } catch (error) {
    const response = await cca.acquireTokenByClientCredential(clientCredentialRequest);
    return response.accessToken;
  }
}

const graphClient = Client.initWithMiddleware({
  authProvider: {
    getAccessToken: getAccessToken
  }
});
```

### Step 6: Monitor Channel Messages

```javascript
async function monitorChannelMessages(teamId, channelId) {
  try {
    // Get recent messages
    const messages = await graphClient
      .teams(teamId)
      .channels(channelId)
      .messages
      .get();

    for (const message of messages.value) {
      if (message.body && message.body.content) {
        // Send to n8n webhook
        await sendToN8N({
          message_id: message.id,
          channel_id: channelId,
          user_id: message.from.user.id,
          user_name: message.from.user.displayName,
          message_text: message.body.content,
          timestamp: message.createdDateTime
        });
      }
    }
  } catch (error) {
    console.error('Error monitoring messages:', error);
  }
}

async function sendToN8N(messageData) {
  const n8nWebhookUrl = 'https://primary-production-674a.up.railway.app/webhook-test/teams-channel-monitor';
  
  try {
    const response = await fetch(n8nWebhookUrl, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(messageData)
    });
    
    if (response.ok) {
      console.log('✅ Message sent to n8n successfully');
    } else {
      console.log('❌ Failed to send to n8n:', response.status);
    }
  } catch (error) {
    console.error('❌ Error sending to n8n:', error);
  }
}
```

## Option 2: Teams Bot Framework (Simpler)

### Step 1: Create Bot in Azure

1. **Go to Azure Portal** → **Create a resource**
2. **Search for "Azure Bot"** → **Create**
3. **Fill in details:**
   - Bot name: "stafind-ai-bot"
   - Subscription: Your subscription
   - Resource group: Create new or use existing
   - Pricing tier: F0 (Free)
4. **Click "Review + create"** → **Create**

### Step 2: Configure Bot

1. **Go to your bot resource**
2. **Copy the "Microsoft App ID"**
3. **Go to "Configuration"** → **Add Microsoft App ID and Password**
4. **Create a new password** → **Copy the password**

### Step 3: Add Bot to Teams

1. **Go to Teams** → **Apps** → **Manage your apps**
2. **Click "Upload a custom app"**
3. **Upload the bot manifest** (see below)

### Bot Manifest Example

```json
{
  "$schema": "https://developer.microsoft.com/en-us/json-schemas/teams/v1.16/MicrosoftTeams.schema.json",
  "manifestVersion": "1.16",
  "version": "1.0.0",
  "id": "YOUR_BOT_APP_ID",
  "packageName": "com.stafind.aibot",
  "developer": {
    "name": "Stafind",
    "websiteUrl": "https://stafind.com",
    "privacyUrl": "https://stafind.com/privacy",
    "termsOfUseUrl": "https://stafind.com/terms"
  },
  "icons": {
    "color": "color.png",
    "outline": "outline.png"
  },
  "name": {
    "short": "Stafind AI",
    "full": "Stafind AI Agent"
  },
  "description": {
    "short": "AI-powered employee matching",
    "full": "Stafind AI Agent helps match employees with job requirements"
  },
  "accentColor": "#FFFFFF",
  "bots": [
    {
      "botId": "YOUR_BOT_APP_ID",
      "scopes": ["personal", "team", "groupchat"],
      "supportsFiles": false,
      "isNotificationOnly": false
    }
  ],
  "permissions": ["identity", "messageTeamMembers"],
  "validDomains": []
}
```

## Option 3: Teams Webhook (Easiest for Testing)

### Step 1: Create Webhook in Teams Channel

1. **Go to your Teams channel**
2. **Click "..." menu** → **Connectors**
3. **Find "Incoming Webhook"** → **Configure**
4. **Name: "Stafind AI Agent"**
5. **Copy the webhook URL**

### Step 2: Test Webhook

Use the provided `teams-webhook-test.js` script to send test messages.

## Testing Your Integration

### Test with n8n Workflow

1. **Import one of the Teams workflows** into n8n
2. **Activate the workflow**
3. **Send a test message** to your Teams channel
4. **Check n8n execution logs** to see if the message was processed

### Sample Test Message

```json
{
  "message_type": "message",
  "message_id": "test-123",
  "channel_id": "general",
  "user_id": "user-123",
  "user_name": "John Doe",
  "message_text": "Looking for React developers with 3+ years experience in TypeScript and Node.js"
}
```

## Troubleshooting

### Common Issues

1. **Permission denied**: Check API permissions in Azure
2. **Webhook not receiving**: Verify webhook URL and n8n workflow activation
3. **Message format**: Ensure message structure matches expected format

### Debug Steps

1. **Check Azure logs** for authentication issues
2. **Verify n8n webhook URL** is correct
3. **Test with simple message** first
4. **Check Teams channel permissions**

## Next Steps

1. **Choose your integration method** (Graph API recommended)
2. **Set up authentication** (Azure App Registration)
3. **Test with sample messages**
4. **Monitor n8n execution logs**
5. **Scale to production** with proper error handling
