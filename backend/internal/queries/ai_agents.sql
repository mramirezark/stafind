-- AI Agent-related SQL queries

-- Create AI agent request
-- Query name: create_ai_agent_request
INSERT INTO ai_agent_requests (teams_message_id, channel_id, user_id, user_name, message_text, attachment_url, status, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, created_at

-- Get AI agent request by ID
-- Query name: get_ai_agent_request_by_id
SELECT id, teams_message_id, channel_id, user_id, user_name, message_text, attachment_url, 
       extracted_text, extracted_skills, status, error, created_at, processed_at
FROM ai_agent_requests 
WHERE id = $1

-- Get AI agent request by Teams message ID
-- Query name: get_ai_agent_request_by_teams_message_id
SELECT id, teams_message_id, channel_id, user_id, user_name, message_text, attachment_url, 
       extracted_text, extracted_skills, status, error, created_at, processed_at
FROM ai_agent_requests 
WHERE teams_message_id = $1

-- Update AI agent request
-- Query name: update_ai_agent_request
UPDATE ai_agent_requests 
SET extracted_text = $1, extracted_skills = $2, status = $3, error = $4, processed_at = $5
WHERE id = $6

-- Update AI agent request status only
-- Query name: update_ai_agent_status
UPDATE ai_agent_requests 
SET status = $1
WHERE id = $2

-- Get all AI agent requests with pagination
-- Query name: get_all_ai_agent_requests
SELECT id, teams_message_id, channel_id, user_id, user_name, message_text, attachment_url, 
       extracted_text, extracted_skills, status, error, created_at, processed_at
FROM ai_agent_requests 
ORDER BY created_at DESC
LIMIT $1 OFFSET $2

-- Save AI agent response
-- Query name: save_ai_agent_response
INSERT INTO ai_agent_responses (request_id, matches, summary, processing_time_ms, status, error, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)

-- Get AI agent response by request ID
-- Query name: get_ai_agent_response_by_request_id
SELECT request_id, matches, summary, processing_time_ms, status, error, created_at
FROM ai_agent_responses 
WHERE request_id = $1
