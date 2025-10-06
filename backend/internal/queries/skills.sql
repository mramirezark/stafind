-- Skill-related SQL queries

-- Get all skills with categories
-- Query name: get_all_skills
SELECT s.id, s.name, c.id as category_id, c.name as category_name
FROM skills s
LEFT JOIN skills_categories sc ON s.id = sc.skill_id
LEFT JOIN categories c ON sc.category_id = c.id
ORDER BY s.name;

-- Get skill by ID
-- Query name: get_skill_by_id
SELECT id, name FROM skills WHERE id = $1;

-- Get skill by name
-- Query name: get_skill_by_name
SELECT id, name FROM skills WHERE name = $1;

-- Get skills by category ID
-- Query name: get_skills_by_category_id
SELECT s.id, s.name FROM skills s
INNER JOIN skills_categories sc ON s.id = sc.skill_id
WHERE sc.category_id = $1
ORDER BY s.name;

-- Get skills with their categories
-- Query name: get_skills_with_categories
SELECT s.id, s.name, c.id as category_id, c.name as category_name
FROM skills s
LEFT JOIN skills_categories sc ON s.id = sc.skill_id
LEFT JOIN categories c ON sc.category_id = c.id
ORDER BY s.name, c.name;

-- Search skills by name (case-insensitive)
-- Query name: search_skills
SELECT id, name FROM skills WHERE LOWER(name) LIKE LOWER($1) ORDER BY name;

-- Get all categories
-- Query name: get_all_categories
SELECT id, name FROM categories ORDER BY name;

-- Get category by ID
-- Query name: get_category_by_id
SELECT id, name FROM categories WHERE id = $1;

-- Get category by name
-- Query name: get_category_by_name
SELECT id, name FROM categories WHERE name = $1;

-- Get categories with skill count
-- Query name: get_categories_with_skill_count
SELECT c.id, c.name, COUNT(sc.skill_id) as skill_count
FROM categories c
LEFT JOIN skills_categories sc ON c.id = sc.category_id
GROUP BY c.id, c.name
ORDER BY c.name;

-- Create category
-- Query name: create_category
INSERT INTO categories (name) VALUES ($1) RETURNING id, name;

-- Update category
-- Query name: update_category
UPDATE categories SET name = $1 WHERE id = $2 RETURNING id, name;

-- Delete category
-- Query name: delete_category
DELETE FROM categories WHERE id = $1;

-- Delete categories batch
-- Query name: delete_categories_batch
DELETE FROM categories WHERE id = ANY($1);

-- Get category stats
-- Query name: get_category_stats
SELECT 
    COUNT(DISTINCT s.id) as total_skills,
    COUNT(DISTINCT c.id) as total_categories,
    s.id as most_popular_skill_id,
    s.name as most_popular_skill_name
FROM categories c
LEFT JOIN skills_categories sc ON c.id = sc.category_id
LEFT JOIN skills s ON sc.skill_id = s.id
LEFT JOIN employee_skills es ON s.id = es.skill_id
GROUP BY s.id, s.name
ORDER BY COUNT(es.employee_id) DESC
LIMIT 1;
SELECT id, name FROM skills 
WHERE LOWER(name) LIKE LOWER('%' || $1 || '%') 
ORDER BY name;

-- Get popular skills (most used by employees)
-- Query name: get_popular_skills
SELECT s.id, s.name, s.category, COUNT(es.skill_id) as employee_count
FROM skills s
LEFT JOIN employee_skills es ON s.id = es.skill_id
GROUP BY s.id, s.name, s.category
ORDER BY employee_count DESC, s.name
LIMIT $1;

-- Get skills by IDs
-- Query name: get_skills_by_ids
SELECT id, name FROM skills WHERE id = ANY($1) ORDER BY name;

-- Get skills with employee count
-- Query name: get_skills_with_employee_count
SELECT s.id, s.name, s.category, COUNT(es.skill_id) as employee_count
FROM skills s
LEFT JOIN employee_skills es ON s.id = es.skill_id
GROUP BY s.id, s.name, s.category
ORDER BY s.name;

-- Create new skill
-- Query name: create_skill
INSERT INTO skills (name)
VALUES ($1)
RETURNING id, name;

-- Create multiple skills
-- Query name: create_skills_batch
INSERT INTO skills (name) VALUES 
-- This will be dynamically generated based on the number of skills

-- Update skill
-- Query name: update_skill
UPDATE skills 
SET name = $1
WHERE id = $2
RETURNING id, name;

-- Delete skill
-- Query name: delete_skill
DELETE FROM skills WHERE id = $1;

-- Delete multiple skills
-- Query name: delete_skills_batch
DELETE FROM skills WHERE id = ANY($1);

-- Add skill to category
-- Query name: add_skill_to_category
INSERT INTO skills_categories (skill_id, category_id)
VALUES ($1, $2)
ON CONFLICT (skill_id, category_id) DO NOTHING;

-- Remove skill from category
-- Query name: remove_skill_from_category
DELETE FROM skills_categories 
WHERE skill_id = $1 AND category_id = $2;

-- Remove all skill categories
-- Query name: remove_all_skill_categories
DELETE FROM skills_categories 
WHERE skill_id = $1;

-- Get skill categories
-- Query name: get_skill_categories
SELECT c.id, c.name FROM categories c
INNER JOIN skills_categories sc ON c.id = sc.category_id
WHERE sc.skill_id = $1
ORDER BY c.name;

-- Get skill statistics
-- Query name: get_skill_stats
WITH most_popular AS (
    SELECT s.id, s.name
    FROM skills s
    LEFT JOIN employee_skills es ON s.id = es.skill_id
    GROUP BY s.id, s.name
    ORDER BY COUNT(es.skill_id) DESC
    LIMIT 1
)
SELECT 
    (SELECT COUNT(*) FROM skills) as total_skills,
    (SELECT COUNT(*) FROM categories) as total_categories,
    (SELECT id FROM most_popular) as most_popular_id,
    (SELECT name FROM most_popular) as most_popular_name;
