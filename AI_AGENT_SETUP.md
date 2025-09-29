# AI Agent Setup Guide

This guide explains how to set up the AI Agent feature that integrates with Microsoft Teams and uses n8n for workflow automation.

## Overview

The AI Agent feature provides:
- Microsoft Teams message processing
- Text extraction from attachments
- OpenAI-powered skill extraction
- Employee matching and ranking
- Automated response generation
- Error handling and admin notifications

## Prerequisites

1. **OpenAI API Key**: Required for skill extraction
2. **Microsoft Teams Webhook**: For sending responses back to Teams
3. **n8n Instance**: For workflow automation
4. **SMTP Configuration**: For admin error notifications

## Environment Variables

Add these variables to your `.env` file:

```bash
# OpenAI Configuration
OPENAI_API_KEY=your_openai_api_key_here

# Microsoft Teams Configuration
TEAMS_WEBHOOK_URL=your_teams_webhook_url_here

# Email Configuration for Admin Notifications
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your_email@gmail.com
SMTP_PASS=your_app_password_here
ADMIN_EMAIL=admin@yourcompany.com

# Stafind API Configuration (for n8n)
STAFIND_API_URL=http://localhost:8080
STAFIND_API_TOKEN=your_api_token_here
```

## Database Setup

The AI Agent feature requires additional database tables. Run the migration:

```bash
# The migration V7__Create_ai_agent_tables.sql will be automatically applied
# when you start the server with Flyway migrations enabled
```

## n8n Workflow Setup

1. **Import the Workflow**:
   - Copy the workflow from `n8n-workflows/ai-agent-teams-workflow.json`
   - Import it into your n8n instance

2. **Configure Environment Variables in n8n**:
   - `STAFIND_API_URL`: Your Stafind API URL
   - `STAFIND_API_TOKEN`: API token for authentication
   - `TEAMS_WEBHOOK_URL`: Microsoft Teams webhook URL

3. **Activate the Workflow**:
   - Enable the workflow in n8n
   - Note the webhook URL for Microsoft Teams integration

## Microsoft Teams Setup

1. **Create a Webhook**:
   - Go to your Teams channel
   - Click on the three dots (...) next to the channel name
   - Select "Connectors"
   - Find "Incoming Webhook" and configure it
   - Copy the webhook URL

2. **Configure the Bot**:
   - Set up a bot that can receive messages
   - Configure it to send messages to the n8n webhook URL

## API Endpoints

The AI Agent provides these endpoints:

### Public Endpoints
- `POST /api/v1/webhooks/teams` - Webhook for n8n integration

### Protected Endpoints (require authentication)
- `GET /api/v1/ai-agent/requests` - List AI agent requests
- `GET /api/v1/ai-agent/requests/:id` - Get specific request
- `POST /api/v1/ai-agent/requests` - Create new request
- `POST /api/v1/ai-agent/requests/:id/process` - Process request
- `POST /api/v1/ai-agent/extract-skills` - Extract skills from text

## Usage Flow

1. **User sends message in Teams** with job requirements or skills
2. **n8n receives the message** via webhook
3. **n8n processes the message**:
   - Checks for attachments
   - Extracts text if needed
   - Sends data to Stafind API
4. **Stafind API processes the request**:
   - Extracts skills using OpenAI
   - Finds matching employees
   - Generates explanations
5. **Response is sent back to Teams** via n8n

## Error Handling

- **API Errors**: Logged to database and admin notified via email
- **OpenAI Errors**: Handled gracefully with fallback responses
- **Teams Errors**: Logged and retried if possible

## Monitoring

- Check AI agent requests in the admin panel
- Monitor error logs in the database
- Review Teams message delivery status

## Troubleshooting

### Common Issues

1. **OpenAI API Key Invalid**:
   - Verify the API key is correct
   - Check API usage limits

2. **Teams Webhook Not Working**:
   - Verify webhook URL is correct
   - Check Teams connector configuration

3. **Database Migration Fails**:
   - Ensure database is accessible
   - Check Flyway configuration

4. **n8n Workflow Not Triggering**:
   - Verify webhook URL in n8n
   - Check n8n execution logs

### Debug Mode

Enable debug logging by setting:
```bash
LOG_LEVEL=debug
```

## Security Considerations

1. **API Authentication**: All protected endpoints require valid JWT tokens
2. **Webhook Security**: Consider adding webhook signature verification
3. **Data Privacy**: Ensure OpenAI API usage complies with your data policies
4. **Rate Limiting**: Implement rate limiting for webhook endpoints

## Performance Optimization

1. **Async Processing**: AI requests are processed asynchronously
2. **Caching**: Consider caching skill extraction results
3. **Database Indexing**: Ensure proper indexes on AI agent tables
4. **Resource Limits**: Monitor OpenAI API usage and costs

## Future Enhancements

- Support for more file types (Word, Excel, etc.)
- Advanced skill matching algorithms
- Integration with other messaging platforms
- Real-time processing status updates
- Advanced analytics and reporting
