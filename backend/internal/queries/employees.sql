-- Employee-related SQL queries

-- Get all employees with their basic information
-- Query name: get_all_employees
SELECT e.id, e.name, e.email, e.department, e.level, e.location, e.bio, e.current_project, e.resume_url, e.created_at, e.updated_at
FROM employees e
ORDER BY e.name;

-- Get all employees with skills in a single query (no N+1)
-- Query name: get_all_employees_with_skills
SELECT 
    e.id, e.name, e.email, e.department, e.level, e.location, e.bio, e.current_project, e.resume_url, e.created_at, e.updated_at,
    s.id as skill_id, s.name as skill_name, es.proficiency_level, es.years_experience
FROM employees e
LEFT JOIN employee_skills es ON e.id = es.employee_id
LEFT JOIN skills s ON es.skill_id = s.id
ORDER BY e.id, s.name;

-- Get employee by ID
-- Query name: get_employee_by_id
SELECT e.id, e.name, e.email, e.department, e.level, e.location, e.bio, e.current_project, e.resume_url, e.created_at, e.updated_at
FROM employees e
WHERE e.id = $1;

-- Get employee by email with extraction data
-- Query name: get_employee_by_email
SELECT e.id, e.name, e.email, e.department, e.level, e.location, e.bio, e.current_project, e.resume_url,
       e.original_text, e.extracted_data, e.extraction_timestamp, e.extraction_source, e.extraction_status,
       e.created_at, e.updated_at
FROM employees e
WHERE e.email = $1;

-- Create new employee
-- Query name: create_employee
INSERT INTO employees (name, email, department, level, location, bio, current_project, resume_url)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, created_at, updated_at;

-- Create new employee with extraction data
-- Query name: create_employee_with_extraction
INSERT INTO employees (name, email, department, level, location, bio, current_project, resume_url,
                      original_text, extracted_data, extraction_timestamp, extraction_source, extraction_status)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING id, created_at, updated_at;

-- Update employee
-- Query name: update_employee
UPDATE employees 
SET name = $1, email = $2, department = $3, level = $4, location = $5, bio = $6, current_project = $7, updated_at = CURRENT_TIMESTAMP
WHERE id = $8;

-- Update employee with extraction data
-- Query name: update_employee_extraction
UPDATE employees 
SET name = $1, email = $2, department = $3, level = $4, location = $5, bio = $6, current_project = $7, resume_url = $8,
    original_text = $9, extracted_data = $10, extraction_timestamp = $11, extraction_source = $12, extraction_status = $13,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $14;

-- Delete employee
-- Query name: delete_employee
DELETE FROM employees WHERE id = $1;

-- Get employee skills
-- Query name: get_employee_skills
SELECT s.id, s.name, es.proficiency_level, es.years_experience
FROM skills s
JOIN employee_skills es ON s.id = es.skill_id
WHERE es.employee_id = $1
ORDER BY s.name;

-- Add employee skill relationship
-- Query name: add_employee_skill
INSERT INTO employee_skills (employee_id, skill_id, proficiency_level, years_experience)
VALUES ($1, $2, $3, $4);

-- Remove all employee skills
-- Query name: remove_employee_skills
DELETE FROM employee_skills WHERE employee_id = $1;

-- Get employees with specific skills (optimized for matching)
-- Query name: get_employees_with_skills
SELECT DISTINCT e.id, e.name, e.email, e.department, e.level, e.location, e.bio, e.current_project, e.resume_url, e.created_at, e.updated_at
FROM employees e
JOIN employee_skills es ON e.id = es.employee_id
JOIN skills s ON es.skill_id = s.id
WHERE LOWER(s.name) = ANY($1)
ORDER BY e.name;

-- Get employees with skills in a single query (no N+1)
-- Query name: get_employees_with_skills_optimized
SELECT 
    e.id, e.name, e.email, e.department, e.level, e.location, e.bio, e.current_project, e.resume_url, e.created_at, e.updated_at,
    s.id as skill_id, s.name as skill_name, es.proficiency_level, es.years_experience
FROM employees e
JOIN employee_skills es ON e.id = es.employee_id
JOIN skills s ON es.skill_id = s.id
WHERE e.id IN (
    SELECT DISTINCT e2.id
    FROM employees e2
    JOIN employee_skills es2 ON e2.id = es2.employee_id
    JOIN skills s2 ON es2.skill_id = s2.id
    WHERE LOWER(s2.name) = ANY($1)
)
ORDER BY e.id, s.name;

-- Get employees with skills and their matching skills (for scoring)
-- Query name: get_employees_with_matching_skills
SELECT e.id, e.name, e.email, e.department, e.level, e.location, e.bio, e.current_project, e.created_at, e.updated_at,
       s.name as skill_name, es.proficiency_level, es.years_experience
FROM employees e
JOIN employee_skills es ON e.id = es.employee_id
JOIN skills s ON es.skill_id = s.id
WHERE LOWER(s.name) = ANY($1)
ORDER BY e.id, s.name;