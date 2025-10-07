-- CV Extract SQL queries

-- Create CV extraction record
-- Query name: create_cv_extract
INSERT INTO cv_extract (
    extract_request_id, status, num_files, files_processed, files_failed,
    total_processing_time_ms, average_processing_time_ms,
    started_at, completed_at, error_message, metadata
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING id, created_at, updated_at

-- Get CV extraction by ID
-- Query name: get_cv_extract_by_id
SELECT id, extract_request_id, status, num_files, files_processed, files_failed,
       total_processing_time_ms, average_processing_time_ms,
       started_at, completed_at, error_message, metadata,
       created_at, updated_at
FROM cv_extract
WHERE id = $1

-- Get CV extraction by extract_request_id
-- Query name: get_cv_extract_by_request_id
SELECT id, extract_request_id, status, num_files, files_processed, files_failed,
       total_processing_time_ms, average_processing_time_ms,
       started_at, completed_at, error_message, metadata,
       created_at, updated_at
FROM cv_extract
WHERE extract_request_id = $1

-- Upsert CV extraction (insert or update)
-- Query name: upsert_cv_extract
INSERT INTO cv_extract (
    extract_request_id, status, num_files, files_processed, files_failed,
    total_processing_time_ms, average_processing_time_ms,
    started_at, completed_at, error_message, metadata
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
ON CONFLICT (extract_request_id) 
DO UPDATE SET
    status = EXCLUDED.status,
    files_processed = EXCLUDED.files_processed,
    files_failed = EXCLUDED.files_failed,
    total_processing_time_ms = EXCLUDED.total_processing_time_ms,
    average_processing_time_ms = EXCLUDED.average_processing_time_ms,
    completed_at = EXCLUDED.completed_at,
    error_message = EXCLUDED.error_message,
    metadata = EXCLUDED.metadata,
    updated_at = CURRENT_TIMESTAMP
RETURNING id, extract_request_id, status, num_files, files_processed, files_failed,
          total_processing_time_ms, average_processing_time_ms,
          started_at, completed_at, error_message, metadata,
          created_at, updated_at

-- Update CV extraction
-- Query name: update_cv_extract
UPDATE cv_extract 
SET status = $2, files_processed = $3, files_failed = $4,
    total_processing_time_ms = $5, average_processing_time_ms = $6,
    completed_at = $7, error_message = $8, metadata = $9,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, status, num_files, files_processed, files_failed,
          total_processing_time_ms, average_processing_time_ms,
          started_at, completed_at, error_message, metadata,
          created_at, updated_at

-- List CV extraction records with pagination
-- Query name: list_cv_extract
SELECT id, extract_request_id, status, num_files, files_processed, files_failed,
       total_processing_time_ms, average_processing_time_ms,
       started_at, completed_at, error_message, metadata,
       created_at, updated_at
FROM cv_extract
WHERE ($1::text IS NULL OR status = $1)
  AND ($2::timestamp IS NULL OR started_at >= $2)
  AND ($3::timestamp IS NULL OR started_at <= $3)
ORDER BY created_at DESC
LIMIT $4 OFFSET $5

-- Count CV extraction records
-- Query name: count_cv_extract
SELECT COUNT(*)
FROM cv_extract
WHERE ($1::text IS NULL OR status = $1)
  AND ($2::timestamp IS NULL OR started_at >= $2)
  AND ($3::timestamp IS NULL OR started_at <= $3)

-- Get CV extraction statistics
-- Query name: get_cv_extract_stats
SELECT 
    COUNT(*) as total_extractions,
    COUNT(CASE WHEN status = 'completed' THEN 1 END) as successful_extractions,
    COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_extractions,
    AVG(total_processing_time_ms) as average_processing_time_ms,
    SUM(files_processed) as total_files_processed,
    SUM(files_failed) as total_files_failed,
    CASE 
        WHEN COUNT(*) > 0 THEN 
            (COUNT(CASE WHEN status = 'completed' THEN 1 END)::float / COUNT(*)::float) * 100
        ELSE 0 
    END as success_rate
FROM cv_extract
WHERE ($1::timestamp IS NULL OR started_at >= $1)
  AND ($2::timestamp IS NULL OR started_at <= $2)

-- Get recent CV extraction records
-- Query name: get_recent_cv_extract
SELECT id, extract_request_id, status, num_files, files_processed, files_failed,
       total_processing_time_ms, average_processing_time_ms,
       started_at, completed_at, error_message, metadata,
       created_at, updated_at
FROM cv_extract
ORDER BY created_at DESC
LIMIT $1

-- Delete CV extraction record
-- Query name: delete_cv_extract
DELETE FROM cv_extract WHERE id = $1

-- Get CV extraction by status
-- Query name: get_cv_extract_by_status
SELECT id, extract_request_id, status, num_files, files_processed, files_failed,
       total_processing_time_ms, average_processing_time_ms,
       started_at, completed_at, error_message, metadata,
       created_at, updated_at
FROM cv_extract
WHERE status = $1
ORDER BY created_at DESC
