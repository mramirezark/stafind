-- Create CV extract table to monitor extraction processes
CREATE TABLE cv_extract (
    id SERIAL PRIMARY KEY,
    extract_request_id VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, processing, completed, failed
    num_files INTEGER NOT NULL DEFAULT 0,
    files_processed INTEGER NOT NULL DEFAULT 0,
    files_failed INTEGER NOT NULL DEFAULT 0,
    total_processing_time_ms BIGINT, -- Total processing time in milliseconds
    average_processing_time_ms BIGINT, -- Average processing time per file in milliseconds
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP NULL,
    error_message TEXT,
    metadata JSONB, -- Additional tracking metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX idx_cv_extract_extract_request_id ON cv_extract(extract_request_id);
CREATE INDEX idx_cv_extract_status ON cv_extract(status);
CREATE INDEX idx_cv_extract_started_at ON cv_extract(started_at);
CREATE INDEX idx_cv_extract_completed_at ON cv_extract(completed_at);
CREATE INDEX idx_cv_extract_created_at ON cv_extract(created_at);