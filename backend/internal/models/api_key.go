package models

import (
	"time"
)

// APIKey represents an API key in the system
type APIKey struct {
	ID          int        `json:"id" db:"id"`
	KeyHash     string     `json:"-" db:"key_hash"` // Never return the actual key
	ServiceName string     `json:"service_name" db:"service_name"`
	Description string     `json:"description" db:"description"`
	IsActive    bool       `json:"is_active" db:"is_active"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	LastUsedAt  *time.Time `json:"last_used_at,omitempty" db:"last_used_at"`
	Permissions []string   `json:"permissions,omitempty"` // What this key can do
}

// CreateAPIKeyRequest represents the request to create a new API key
type CreateAPIKeyRequest struct {
	ServiceName string     `json:"service_name" binding:"required"`
	Description string     `json:"description"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	Permissions []string   `json:"permissions,omitempty"`
}

// APIKeyResponse represents the response when creating an API key
type APIKeyResponse struct {
	ID          int        `json:"id"`
	APIKey      string     `json:"api_key"` // Only returned once when created
	ServiceName string     `json:"service_name"`
	Description string     `json:"description"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	Permissions []string   `json:"permissions,omitempty"`
}
