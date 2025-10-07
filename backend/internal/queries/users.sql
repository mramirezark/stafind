-- User-related SQL queries

-- Create user
-- Query name: create_user
INSERT INTO users (username, email, password_hash, first_name, last_name, role_id, is_active)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, created_at, updated_at

-- Get user by ID with role information
-- Query name: get_user_by_id
SELECT u.id, u.username, u.email, u.password_hash, u.first_name, u.last_name, 
       u.role_id, u.is_active, u.last_login, u.created_at, u.updated_at,
       r.id, r.name, r.description, r.created_at, r.updated_at
FROM users u
LEFT JOIN roles r ON u.role_id = r.id
WHERE u.id = $1

-- Get user by email with role information
-- Query name: get_user_by_email
SELECT u.id, u.username, u.email, u.password_hash, u.first_name, u.last_name, 
       u.role_id, u.is_active, u.last_login, u.created_at, u.updated_at,
       r.id, r.name, r.description, r.created_at, r.updated_at
FROM users u
LEFT JOIN roles r ON u.role_id = r.id
WHERE u.email = $1

-- Get user by username with role information
-- Query name: get_user_by_username
SELECT u.id, u.username, u.email, u.password_hash, u.first_name, u.last_name, 
       u.role_id, u.is_active, u.last_login, u.created_at, u.updated_at,
       r.id, r.name, r.description, r.created_at, r.updated_at
FROM users u
LEFT JOIN roles r ON u.role_id = r.id
WHERE u.username = $1

-- Update user information
-- Query name: update_user
UPDATE users 
SET username = $1, email = $2, first_name = $3, last_name = $4, role_id = $5, 
    is_active = $6, updated_at = CURRENT_TIMESTAMP
WHERE id = $7

-- Delete user
-- Query name: delete_user
DELETE FROM users WHERE id = $1

-- List users with pagination and role information
-- Query name: list_users
SELECT u.id, u.email, u.first_name, u.last_name, u.role_id, u.is_active, 
       u.last_login, u.created_at, u.updated_at,
       r.id, r.name, r.description, r.created_at, r.updated_at
FROM users u
LEFT JOIN roles r ON u.role_id = r.id
ORDER BY u.created_at DESC
LIMIT $1 OFFSET $2

-- Get total number of users
-- Query name: get_user_count
SELECT COUNT(*) FROM users

-- Create user session
-- Query name: create_user_session
INSERT INTO user_sessions (user_id, token_hash, expires_at)
VALUES ($1, $2, $3)
RETURNING id, created_at

-- Get user session by token hash
-- Query name: get_user_session
SELECT id, user_id, token_hash, expires_at, created_at, is_revoked
FROM user_sessions
WHERE token_hash = $1 AND is_revoked = FALSE

-- Revoke user session
-- Query name: revoke_user_session
UPDATE user_sessions SET is_revoked = TRUE WHERE token_hash = $1

-- Clean up expired user sessions
-- Query name: cleanup_expired_sessions
DELETE FROM user_sessions WHERE expires_at < CURRENT_TIMESTAMP OR is_revoked = TRUE

-- Get all roles for a user
-- Query name: get_user_roles
SELECT r.id, r.name, r.description, r.created_at, r.updated_at
FROM roles r
INNER JOIN user_roles ur ON r.id = ur.role_id
WHERE ur.user_id = $1

-- Assign role to user
-- Query name: assign_role_to_user
INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2) ON CONFLICT (user_id, role_id) DO NOTHING

-- Remove role from user
-- Query name: remove_role_from_user
DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2
