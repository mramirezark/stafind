-- Match-related SQL queries

-- Get matches by job request ID
-- Query name: get_matches_by_job_request
SELECT m.id, m.job_request_id, m.employee_id, m.match_score, m.matching_skills, m.notes, m.created_at,
       e.id as employee_id, e.name as employee_name, e.email as employee_email, e.department as employee_department, 
       e.level as employee_level, e.location as employee_location, e.bio as employee_bio, 
       e.created_at as employee_created_at, e.updated_at as employee_updated_at
FROM matches m
JOIN employees e ON m.employee_id = e.id
WHERE m.job_request_id = $1
ORDER BY m.match_score DESC;

-- Create match
-- Query name: create_match
INSERT INTO matches (job_request_id, employee_id, match_score, matching_skills, notes)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at;

-- Delete matches by job request ID
-- Query name: delete_matches_by_job_request
DELETE FROM matches WHERE job_request_id = $1;

-- Get employee skills for match display
-- Query name: get_employee_skills_for_match
SELECT s.id, s.name, s.category
FROM skills s
JOIN employee_skills es ON s.id = es.skill_id
WHERE es.employee_id = $1
ORDER BY s.name;

-- Get all matches with employee information
-- Query name: get_all_matches
SELECT m.id, m.employee_id, m.match_score, m.matching_skills, m.notes, m.created_at,
       e.id, e.name, e.email, e.department, e.level, e.location, e.bio, e.current_project, e.created_at, e.updated_at
FROM matches m
LEFT JOIN employees e ON m.employee_id = e.id
ORDER BY m.created_at DESC;

-- Delete a match by ID
-- Query name: delete_match
DELETE FROM matches WHERE id = $1;
