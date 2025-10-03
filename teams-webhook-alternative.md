# Teams Webhook Alternative Setup

## Quick Solution: Use Webhook Instead of OAuth2

If you're getting OAuth2 errors, here's a simpler webhook-based approach:

### Step 1: Replace Teams Message Trigger with Webhook

1. **Delete the Teams Message Trigger node**
2. **Add a Webhook node instead**
3. **Configure the webhook:**
   ```
   HTTP Method: POST
   Path: teams-resource-search
   Response Mode: On Received
   ```

### Step 2: Set Up Teams Webhook Connector

1. **Go to Teams** → **Your Channel** → **...** → **Connectors**
2. **Find "Incoming Webhook"** → **Configure**
3. **Set webhook URL:** `https://your-n8n-instance.com/webhook/teams-resource-search`
4. **Name:** "Resource Search Bot"
5. **Upload icon** (optional)
6. **Save**

### Step 3: Test the Webhook

1. **In Teams, click the webhook connector**
2. **Send a test message:**
   ```json
   {
     "text": "Dame recursos con Java",
     "user": "Test User",
     "channel": "General"
   }
   ```

### Step 4: Update Your Workflow

1. **The webhook will receive the message in n8n**
2. **Add a "Set" node to structure the data:**
   ```javascript
   return {
     json: {
       text: $json.text,
       user: $json.user || "Unknown",
       channel: $json.channel || "Unknown",
       timestamp: new Date().toISOString()
     }
   };
   ```

## Benefits of Webhook Approach:

✅ **No OAuth2 complexity**  
✅ **No redirect URI issues**  
✅ **Faster setup**  
✅ **More reliable**  
✅ **Easier to debug**  

## Webhook Message Format:

The webhook will receive messages like:
```json
{
  "text": "Dame los recursos con más experiencia en Java",
  "user": "John Doe", 
  "channel": "General",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## Update Your Workflow:

1. **Replace Teams Message Trigger** with **Webhook node**
2. **Keep the rest of your workflow the same**
3. **The AI analysis and Google Drive search will work identically**

This approach is actually simpler and more reliable than OAuth2 for this use case!
