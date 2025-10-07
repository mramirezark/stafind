-- Category-related SQL queries

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

-- Get skills by category ID
-- Query name: get_skills_by_category_id
SELECT s.id, s.name FROM skills s
INNER JOIN skills_categories sc ON s.id = sc.skill_id
WHERE sc.category_id = $1
ORDER BY s.name;

-- Create new category
-- Query name: create_category
INSERT INTO categories (name)
VALUES ($1)
RETURNING id, name;

-- Update category
-- Query name: update_category
UPDATE categories 
SET name = $1
WHERE id = $2
RETURNING id, name;

-- Delete category
-- Query name: delete_category
DELETE FROM categories WHERE id = $1;

-- Delete multiple categories
-- Query name: delete_categories_batch
DELETE FROM categories WHERE id = ANY($1);

-- Get category statistics
-- Query name: get_category_stats
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
