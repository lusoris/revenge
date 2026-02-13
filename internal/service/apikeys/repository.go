package apikeys

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/infra/database/db"
)

// Repository defines data access for API keys
type Repository interface {
	// Key management
	CreateAPIKey(ctx context.Context, params db.CreateAPIKeyParams) (db.SharedApiKey, error)
	GetAPIKey(ctx context.Context, id uuid.UUID) (db.SharedApiKey, error)
	GetAPIKeyByHash(ctx context.Context, keyHash string) (db.SharedApiKey, error)
	GetAPIKeyByPrefix(ctx context.Context, keyPrefix string) (db.SharedApiKey, error)
	ListUserAPIKeys(ctx context.Context, userID uuid.UUID) ([]db.SharedApiKey, error)
	ListActiveUserAPIKeys(ctx context.Context, userID uuid.UUID) ([]db.SharedApiKey, error)
	CountUserAPIKeys(ctx context.Context, userID uuid.UUID) (int64, error)
	RevokeAPIKey(ctx context.Context, id uuid.UUID) error
	UpdateAPIKeyLastUsed(ctx context.Context, id uuid.UUID) error
	UpdateAPIKeyScopes(ctx context.Context, id uuid.UUID, scopes []string) error
	DeleteAPIKey(ctx context.Context, id uuid.UUID) error
	DeleteExpiredAPIKeys(ctx context.Context) error
}

// CreateKeyRequest contains data for creating an API key
type CreateKeyRequest struct {
	Name        string
	Description *string
	Scopes      []string
	ExpiresAt   *time.Time
}

// APIKey represents an API key (without the raw key)
type APIKey struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        string
	Description *string
	KeyPrefix   string // First 8 chars for identification
	Scopes      []string
	IsActive    bool
	ExpiresAt   *time.Time
	LastUsedAt  *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// CreateKeyResponse contains the created key with the raw token
type CreateKeyResponse struct {
	Key    APIKey
	RawKey string // Only returned on creation!
}

// Service defines the API keys service interface.
// All consumers should depend on this interface, not the concrete implementation.
type Service interface {
	// CreateKey generates a new API key for a user.
	CreateKey(ctx context.Context, userID uuid.UUID, req CreateKeyRequest) (*CreateKeyResponse, error)
	// GetKey retrieves an API key by ID (without raw key).
	GetKey(ctx context.Context, keyID uuid.UUID) (*APIKey, error)
	// ListUserKeys lists active API keys for a user.
	ListUserKeys(ctx context.Context, userID uuid.UUID) ([]APIKey, error)
	// ValidateKey validates a raw API key and returns the associated key data.
	ValidateKey(ctx context.Context, rawKey string) (*APIKey, error)
	// RevokeKey revokes an API key.
	RevokeKey(ctx context.Context, keyID uuid.UUID) error
	// CheckScope checks if an API key has a required scope.
	CheckScope(ctx context.Context, keyID uuid.UUID, requiredScope string) (bool, error)
	// UpdateScopes updates the scopes of an API key.
	UpdateScopes(ctx context.Context, keyID uuid.UUID, scopes []string) error
	// CleanupExpiredKeys deletes expired and inactive keys.
	CleanupExpiredKeys(ctx context.Context) error
}
