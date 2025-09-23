-- Job Request-related SQL queries

-- Get all job requests
-- Query name: get_all_job_requests
SELECT id, title, description, department, required_skills, preferred_skills,
       experience_level, location, priority, status, created_by, created_at, updated_at
FROM job_requests
ORDER BY created_at DESC;

-- Get job request by ID
-- Query name: get_job_request_by_id
SELECT id, title, description, department, required_skills, preferred_skills,
       experience_level, location, priority, status, created_by, created_at, updated_at
FROM job_requests
WHERE id = $1;

-- Create new job request
-- Query name: create_job_request
INSERT INTO job_requests (title, description, department, required_skills, preferred_skills,
                         experience_level, location, priority, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, created_at, updated_at;

-- Update job request
-- Query name: update_job_request
UPDATE job_requests 
SET title = $1, description = $2, department = $3, required_skills = $4, preferred_skills = $5,
    experience_level = $6, location = $7, priority = $8, updated_at = CURRENT_TIMESTAMP
WHERE id = $9;

-- Delete job request
-- Query name: delete_job_request
DELETE FROM job_requests WHERE id = $1;
