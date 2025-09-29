# Backend Debugging Guide

## Setup

1. **Start the database:**
   ```bash
   make dev
   ```

2. **Run migrations:**
   ```bash
   make migrate
   ```

3. **Start debugging in VS Code:**
   - Open VS Code
   - Go to Run and Debug (Ctrl+Shift+D)
   - Select "Debug Backend Server"
   - Press F5 or click the play button

## Debugging Steps

### 1. Set Breakpoints

Set breakpoints in these key locations:

- `backend/internal/handlers/ai_agent_handlers.go:119` - WebhookHandler start
- `backend/internal/handlers/ai_agent_handlers.go:155` - Before CreateAIAgentRequest
- `backend/internal/services/ai_agent_service.go:44` - CreateAIAgentRequest start
- `backend/internal/repositories/ai_agent_repository.go:29` - Repository Create start

### 2. Test the Webhook

Run the debug test script:
```bash
node debug-webhook.js
```

Or use VS Code task:
- Press Ctrl+Shift+P
- Type "Tasks: Run Task"
- Select "Test Webhook"

### 3. Monitor Debug Output

The debug logs will show:
- Webhook data received
- AI agent request creation
- Database query execution
- Any errors that occur

### 4. Common Issues to Check

1. **Database Connection:**
   - Check if database is running
   - Verify connection parameters

2. **Array Field Issues:**
   - Check if `extracted_skills` field is being set incorrectly
   - Verify PostgreSQL array handling

3. **JSON Parsing:**
   - Check if webhook data is parsed correctly
   - Verify null value handling

## Debug Configuration

The debug configuration is in `.vscode/launch.json`:
- Environment variables are set for database connection
- Program points to the main server file
- Working directory is set to backend folder

## Troubleshooting

### If debugging doesn't start:
1. Make sure Go extension is installed
2. Check if database is running
3. Verify environment variables

### If webhook still fails:
1. Check the debug output for specific error messages
2. Look at the database logs
3. Verify the database schema

## Next Steps

Once you identify where the error occurs:
1. Add more specific breakpoints around that area
2. Check the exact values being passed to the database
3. Verify the database schema matches the code expectations
