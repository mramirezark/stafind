package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"stafind-backend/internal/constants"
	"stafind-backend/internal/models"
	"stafind-backend/internal/repositories"
	"time"
)

type apiKeyService struct {
	apiKeyRepo repositories.APIKeyRepository
}

// NewAPIKeyService creates a new API key service
func NewAPIKeyService(apiKeyRepo repositories.APIKeyRepository) APIKeyService {
	return &apiKeyService{
		apiKeyRepo: apiKeyRepo,
	}
}

// generateAPIKey generates a secure random API key
func (s *apiKeyService) generateAPIKey() (string, error) {
	// Generate 32 random bytes
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Convert to hex string and add prefix
	key := "sk_" + hex.EncodeToString(bytes)
	return key, nil
}

// hashAPIKey hashes an API key for secure storage
func (s *apiKeyService) hashAPIKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}

// CreateAPIKey creates a new API key
func (s *apiKeyService) CreateAPIKey(req *models.CreateAPIKeyRequest) (*models.APIKeyResponse, error) {
	// Generate new API key
	apiKey, err := s.generateAPIKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate API key: %v", err)
	}

	// Hash the key for storage
	keyHash := s.hashAPIKey(apiKey)

	// Create API key record
	apiKeyRecord := &models.APIKey{
		KeyHash:     keyHash,
		ServiceName: req.ServiceName,
		Description: req.Description,
		IsActive:    true,
		CreatedAt:   time.Now(),
		ExpiresAt:   req.ExpiresAt,
	}

	// Save to database
	createdKey, err := s.apiKeyRepo.Create(apiKeyRecord)
	if err != nil {
		return nil, fmt.Errorf("failed to save API key: %v", err)
	}

	// Return response with the actual key (only time it's returned)
	return &models.APIKeyResponse{
		ID:          createdKey.ID,
		APIKey:      apiKey, // Only returned once!
		ServiceName: createdKey.ServiceName,
		Description: createdKey.Description,
		IsActive:    createdKey.IsActive,
		CreatedAt:   createdKey.CreatedAt,
		ExpiresAt:   createdKey.ExpiresAt,
		Permissions: req.Permissions,
	}, nil
}

// ValidateAPIKey validates an API key
func (s *apiKeyService) ValidateAPIKey(key string) (*models.APIKey, error) {
	// First try to look up the key as-is (in case it's already hashed)
	apiKey, err := s.apiKeyRepo.GetByHash(key)
	if err != nil {
		// If not found, try hashing the key and looking up again
		keyHash := s.hashAPIKey(key)
		apiKey, err = s.apiKeyRepo.GetByHash(keyHash)
		if err != nil {
			return nil, fmt.Errorf("invalid API key")
		}
	}

	// Check if key is active
	if !apiKey.IsActive {
		return nil, fmt.Errorf("API key is deactivated")
	}

	// Check if key has expired
	if apiKey.ExpiresAt != nil && time.Now().After(*apiKey.ExpiresAt) {
		return nil, fmt.Errorf("API key has expired")
	}

	// Update last used timestamp
	go s.UpdateLastUsed(key)

	return apiKey, nil
}

// GetAPIKeys retrieves all API keys (without the actual keys)
func (s *apiKeyService) GetAPIKeys(limit, offset int) ([]models.APIKey, error) {
	return s.apiKeyRepo.GetAll(limit, offset)
}

// DeactivateAPIKey deactivates an API key
func (s *apiKeyService) DeactivateAPIKey(id int) error {
	return s.apiKeyRepo.Deactivate(id)
}

// UpdateLastUsed updates the last used timestamp for an API key
func (s *apiKeyService) UpdateLastUsed(key string) error {
	keyHash := s.hashAPIKey(key)
	return s.apiKeyRepo.UpdateLastUsed(keyHash)
}

// RotateAPIKey creates a new API key and deactivates the old one
func (s *apiKeyService) RotateAPIKey(oldKeyID int) (*models.APIKeyResponse, error) {
	// Get the old key to preserve some metadata
	oldKey, err := s.apiKeyRepo.GetByID(oldKeyID)
	if err != nil {
		return nil, fmt.Errorf("old API key not found: %v", err)
	}

	// Deactivate the old key
	err = s.apiKeyRepo.Deactivate(oldKeyID)
	if err != nil {
		return nil, fmt.Errorf("failed to deactivate old API key: %v", err)
	}

	// Generate new API key
	newAPIKey, err := s.generateAPIKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate new API key: %v", err)
	}

	// Hash the new key
	keyHash := s.hashAPIKey(newAPIKey)

	// Create new API key record with same metadata as old one
	newKeyRecord := &models.APIKey{
		KeyHash:     keyHash,
		ServiceName: oldKey.ServiceName,
		Description: oldKey.Description + " (rotated)",
		IsActive:    true,
		CreatedAt:   time.Now(),
		ExpiresAt:   oldKey.ExpiresAt, // Keep same expiration
	}

	// Save new key
	createdKey, err := s.apiKeyRepo.Create(newKeyRecord)
	if err != nil {
		return nil, fmt.Errorf("failed to create new API key: %v", err)
	}

	// Return the new key (only time it's returned)
	return &models.APIKeyResponse{
		ID:          createdKey.ID,
		APIKey:      newAPIKey, // New key shown only once!
		ServiceName: createdKey.ServiceName,
		Description: createdKey.Description,
		IsActive:    createdKey.IsActive,
		CreatedAt:   createdKey.CreatedAt,
		ExpiresAt:   createdKey.ExpiresAt,
		Permissions: []string{}, // You might want to preserve permissions
	}, nil
}

// GetEnvironmentAPIKey gets API key from environment variables (fallback)
func GetEnvironmentAPIKey() string {
	return os.Getenv(constants.EnvExternalAPIKey)
}

// ValidateEnvironmentAPIKey validates against environment variable
func ValidateEnvironmentAPIKey(key string) bool {
	envKey := GetEnvironmentAPIKey()
	if envKey == "" {
		// Default for development
		envKey = "dev-api-key-12345"
	}
	return key == envKey
}
