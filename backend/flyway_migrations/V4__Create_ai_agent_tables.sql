-- Create AI agent request table
CREATE TABLE ai_agent_requests (
    id SERIAL PRIMARY KEY,
    teams_message_id VARCHAR(255) NOT NULL UNIQUE,
    channel_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    user_name VARCHAR(255) NOT NULL,
    message_text TEXT,
    attachment_url TEXT,
    extracted_text TEXT,
    extracted_skills TEXT[], -- Array of skill names
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    error TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMP WITH TIME ZONE
);

-- Create index for faster lookups
CREATE INDEX idx_ai_agent_requests_status ON ai_agent_requests(status);
CREATE INDEX idx_ai_agent_requests_teams_message_id ON ai_agent_requests(teams_message_id);
CREATE INDEX idx_ai_agent_requests_created_at ON ai_agent_requests(created_at);

-- Create AI agent response table for logging responses
CREATE TABLE ai_agent_responses (
    id SERIAL PRIMARY KEY,
    request_id INTEGER NOT NULL REFERENCES ai_agent_requests(id) ON DELETE CASCADE,
    matches JSONB NOT NULL, -- Store the matches as JSON
    summary TEXT,
    processing_time_ms BIGINT,
    status VARCHAR(50) NOT NULL,
    error TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index for faster lookups
CREATE INDEX idx_ai_agent_responses_request_id ON ai_agent_responses(request_id);
CREATE INDEX idx_ai_agent_responses_created_at ON ai_agent_responses(created_at);
