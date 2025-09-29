package repositories

import (
	"database/sql"
	"fmt"
	"stafind-backend/internal/models"
	"time"
)

type apiKeyRepository struct {
	*BaseRepository
}

// NewAPIKeyRepository creates a new API key repository
func NewAPIKeyRepository(db *sql.DB) (APIKeyRepository, error) {
	baseRepo, err := NewBaseRepository(db)
	if err != nil {
		return nil, err
	}

	return &apiKeyRepository{BaseRepository: baseRepo}, nil
}

// Create creates a new API key
func (r *apiKeyRepository) Create(key *models.APIKey) (*models.APIKey, error) {
	query := r.MustGetQuery("create_api_key")

	var id int
	var createdAt time.Time
	err := r.db.QueryRow(
		query,
		key.KeyHash,
		key.ServiceName,
		key.Description,
		key.IsActive,
		key.CreatedAt,
		key.ExpiresAt,
	).Scan(&id, &createdAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create API key: %v", err)
	}

	key.ID = id
	key.CreatedAt = createdAt
	return key, nil
}

// GetByID retrieves an API key by ID
func (r *apiKeyRepository) GetByID(id int) (*models.APIKey, error) {
	query := r.MustGetQuery("get_api_key_by_id")

	key := &models.APIKey{}
	var expiresAt, lastUsedAt sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&key.ID,
		&key.KeyHash,
		&key.ServiceName,
		&key.Description,
		&key.IsActive,
		&key.CreatedAt,
		&expiresAt,
		&lastUsedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("API key not found")
		}
		return nil, fmt.Errorf("failed to get API key: %v", err)
	}

	if expiresAt.Valid {
		key.ExpiresAt = &expiresAt.Time
	}
	if lastUsedAt.Valid {
		key.LastUsedAt = &lastUsedAt.Time
	}

	return key, nil
}

// GetByHash retrieves an API key by its hash
func (r *apiKeyRepository) GetByHash(hash string) (*models.APIKey, error) {
	query := r.MustGetQuery("get_api_key_by_hash")

	key := &models.APIKey{}
	var expiresAt, lastUsedAt sql.NullTime

	err := r.db.QueryRow(query, hash).Scan(
		&key.ID,
		&key.KeyHash,
		&key.ServiceName,
		&key.Description,
		&key.IsActive,
		&key.CreatedAt,
		&expiresAt,
		&lastUsedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("API key not found")
		}
		return nil, fmt.Errorf("failed to get API key: %v", err)
	}

	if expiresAt.Valid {
		key.ExpiresAt = &expiresAt.Time
	}
	if lastUsedAt.Valid {
		key.LastUsedAt = &lastUsedAt.Time
	}

	return key, nil
}

// GetAll retrieves all API keys with pagination
func (r *apiKeyRepository) GetAll(limit, offset int) ([]models.APIKey, error) {
	query := r.MustGetQuery("get_all_api_keys")

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get API keys: %v", err)
	}
	defer rows.Close()

	var keys []models.APIKey
	for rows.Next() {
		key := models.APIKey{}
		var expiresAt, lastUsedAt sql.NullTime

		err := rows.Scan(
			&key.ID,
			&key.KeyHash,
			&key.ServiceName,
			&key.Description,
			&key.IsActive,
			&key.CreatedAt,
			&expiresAt,
			&lastUsedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan API key: %v", err)
		}

		if expiresAt.Valid {
			key.ExpiresAt = &expiresAt.Time
		}
		if lastUsedAt.Valid {
			key.LastUsedAt = &lastUsedAt.Time
		}

		keys = append(keys, key)
	}

	return keys, nil
}

// Update updates an API key
func (r *apiKeyRepository) Update(id int, key *models.APIKey) error {
	query := r.MustGetQuery("update_api_key")

	_, err := r.db.Exec(query, id, key.ServiceName, key.Description, key.IsActive, key.ExpiresAt)
	if err != nil {
		return fmt.Errorf("failed to update API key: %v", err)
	}

	return nil
}

// Deactivate deactivates an API key
func (r *apiKeyRepository) Deactivate(id int) error {
	query := r.MustGetQuery("deactivate_api_key")

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to deactivate API key: %v", err)
	}

	return nil
}

// UpdateLastUsed updates the last used timestamp for an API key
func (r *apiKeyRepository) UpdateLastUsed(hash string) error {
	query := r.MustGetQuery("update_api_key_last_used")

	_, err := r.db.Exec(query, time.Now(), hash)
	if err != nil {
		return fmt.Errorf("failed to update last used timestamp: %v", err)
	}

	return nil
}

// Delete deletes an API key
func (r *apiKeyRepository) Delete(id int) error {
	query := r.MustGetQuery("delete_api_key")

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete API key: %v", err)
	}

	return nil
}
