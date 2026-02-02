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
