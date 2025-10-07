-- Create API keys table for external service authentication
CREATE TABLE api_keys (
    id SERIAL PRIMARY KEY,
    key_hash VARCHAR(255) NOT NULL UNIQUE,
    service_name VARCHAR(100) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP,
    last_used_at TIMESTAMP
);

-- Create index for faster lookups
CREATE INDEX idx_api_keys_hash ON api_keys(key_hash);
CREATE INDEX idx_api_keys_service ON api_keys(service_name);
CREATE INDEX idx_api_keys_active ON api_keys(is_active);

-- Insert a default API key for development
-- Hash of 'dev-api-key-12345' using SHA-256
INSERT INTO api_keys (key_hash, service_name, description, is_active, expires_at)
VALUES (
    '8264dc9f07e749d9c2ffead0b25de8cb22bed7af774e189ef224ae015908776b',
    'development',
    'Default development API key',
    true,
    NOW() + INTERVAL '1 year'
);

-- Insert the user's specific API key
-- Hash of 'a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456' using SHA-256
INSERT INTO api_keys (key_hash, service_name, description, is_active, expires_at)
VALUES (
    'c5012edc5a4d3afcbaa410c25cc4bbeada36309a1d9561d169daef902b0327a2',
    'user-service',
    'User provided API key for testing',
    true,
    NOW() + INTERVAL '1 year'
);
