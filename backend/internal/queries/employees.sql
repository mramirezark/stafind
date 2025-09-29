-- Employee-related SQL queries

-- Get all employees with their basic information
-- Query name: get_all_employees
SELECT e.id, e.name, e.email, e.department, e.level, e.location, e.bio, e.current_project, e.created_at, e.updated_at
FROM employees e
ORDER BY e.name;

-- Get employee by ID
-- Query name: get_employee_by_id
SELECT e.id, e.name, e.email, e.department, e.level, e.location, e.bio, e.current_project, e.created_at, e.updated_at
FROM employees e
WHERE e.id = $1;

-- Create new employee
-- Query name: create_employee
INSERT INTO employees (name, email, department, level, location, bio, current_project)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, created_at, updated_at;

-- Update employee
-- Query name: update_employee
UPDATE employees 
SET name = $1, email = $2, department = $3, level = $4, location = $5, bio = $6, current_project = $7, updated_at = CURRENT_TIMESTAMP
WHERE id = $8;

-- Delete employee
-- Query name: delete_employee
DELETE FROM employees WHERE id = $1;

-- Get employee skills
-- Query name: get_employee_skills
SELECT s.id, s.name, s.category, es.proficiency_level, es.years_experience
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
