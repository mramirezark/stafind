-- Role-related SQL queries

-- Create role
-- Query name: create_role
INSERT INTO roles (name, description)
VALUES ($1, $2)
RETURNING id, created_at, updated_at

-- Get role by ID
-- Query name: get_role_by_id
SELECT id, name, description, created_at, updated_at
FROM roles
WHERE id = $1

-- Get role by name
-- Query name: get_role_by_name
SELECT id, name, description, created_at, updated_at
FROM roles
WHERE name = $1

-- Update role
-- Query name: update_role
UPDATE roles 
SET name = $1, description = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $3

-- Delete role
-- Query name: delete_role
DELETE FROM roles WHERE id = $1

-- List all roles
-- Query name: list_roles
SELECT id, name, description, created_at, updated_at
FROM roles
ORDER BY name
