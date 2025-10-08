-- Query name: create_api_key
INSERT INTO api_keys (key_hash, service_name, description, is_active, created_at, expires_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, created_at

-- Query name: get_api_key_by_id
SELECT id, key_hash, service_name, description, is_active, created_at, expires_at, last_used_at
FROM api_keys
WHERE id = $1

-- Query name: get_api_key_by_hash
SELECT id, key_hash, service_name, description, is_active, created_at, expires_at, last_used_at
FROM api_keys
WHERE key_hash = $1

-- Query name: get_all_api_keys
SELECT id, key_hash, service_name, description, is_active, created_at, expires_at, last_used_at
FROM api_keys
ORDER BY created_at DESC
LIMIT $1 OFFSET $2

-- Query name: update_api_key
UPDATE api_keys
SET service_name = $2, description = $3, is_active = $4, expires_at = $5
WHERE id = $1

-- Query name: deactivate_api_key
UPDATE api_keys SET is_active = false WHERE id = $1

-- Query name: update_api_key_last_used
UPDATE api_keys SET last_used_at = $1 WHERE key_hash = $2

-- Query name: delete_api_key
DELETE FROM api_keys WHERE id = $1
