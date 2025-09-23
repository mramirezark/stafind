-- Skill-related SQL queries

-- Get all skills
-- Query name: get_all_skills
SELECT id, name, category FROM skills ORDER BY name;

-- Get skill by ID
-- Query name: get_skill_by_id
SELECT id, name, category FROM skills WHERE id = $1;

-- Get skill by name
-- Query name: get_skill_by_name
SELECT id, name, category FROM skills WHERE name = $1;

-- Create new skill
-- Query name: create_skill
INSERT INTO skills (name, category)
VALUES ($1, $2)
RETURNING id, created_at;

-- Update skill
-- Query name: update_skill
UPDATE skills 
SET name = $1, category = $2
WHERE id = $3
RETURNING id, name, category;

-- Delete skill
-- Query name: delete_skill
DELETE FROM skills WHERE id = $1;
