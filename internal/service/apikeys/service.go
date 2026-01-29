// Package apikeys provides API key management services.
package apikeys

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/lusoris/revenge/internal/infra/database/db"
)

var (
	// ErrKeyNotFound indicates the API key was not found.
	ErrKeyNotFound = errors.New("api key not found")
	// ErrKeyExpired indicates the API key has expired.
	ErrKeyExpired = errors.New("api key expired")
	// ErrInvalidKey indicates the API key is invalid.
	ErrInvalidKey = errors.New("invalid api key")
)

// Service provides API key management operations.
type Service struct {
	queries *db.Queries
	logger  *slog.Logger
}

// NewService creates a new API keys service.
func NewService(queries *db.Queries, logger *slog.Logger) *Service {
	return &Service{
		queries: queries,
		logger:  logger.With(slog.String("service", "apikeys")),
	}
}

// CreateParams contains parameters for creating an API key.
type CreateParams struct {
	UserID    uuid.UUID
	Name      string
	Scopes    []string
	ExpiresAt *time.Time
}

// CreateResult contains the result of creating an API key.
type CreateResult struct {
	Key    *db.ApiKey // Database record
	RawKey string     // The actual key value (only returned on creation)
}

// Create creates a new API key.
// The raw key is only returned once and cannot be retrieved later.
func (s *Service) Create(ctx context.Context, params CreateParams) (*CreateResult, error) {
	// Generate a random 32-byte key
	keyBytes := make([]byte, 32)
	if _, err := rand.Read(keyBytes); err != nil {
		return nil, err
	}

	// Encode as base64 for the raw key
	rawKey := base64.URLEncoding.EncodeToString(keyBytes)

	// Hash the key for storage
	hash := sha256.Sum256([]byte(rawKey))
	keyHash := hex.EncodeToString(hash[:])

	// Get prefix (first 8 chars of raw key)
	keyPrefix := rawKey[:8]

	// Default scopes to empty if nil
	scopes := params.Scopes
	if scopes == nil {
		scopes = []string{}
	}

	// Convert expiresAt to pgtype.Timestamptz
	var expiresAt pgtype.Timestamptz
	if params.ExpiresAt != nil {
		expiresAt = pgtype.Timestamptz{Time: *params.ExpiresAt, Valid: true}
	}

	apiKey, err := s.queries.CreateAPIKey(ctx, db.CreateAPIKeyParams{
		UserID:    params.UserID,
		Name:      params.Name,
		KeyHash:   keyHash,
		KeyPrefix: keyPrefix,
		Scopes:    scopes,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return nil, err
	}

	s.logger.Info("API key created",
		slog.String("key_id", apiKey.ID.String()),
		slog.String("user_id", params.UserID.String()),
		slog.String("name", params.Name),
	)

	return &CreateResult{
		Key:    &apiKey,
		RawKey: rawKey,
	}, nil
}

// Validate validates an API key and returns the associated record.
func (s *Service) Validate(ctx context.Context, rawKey string) (*db.ApiKey, error) {
	// Hash the provided key
	hash := sha256.Sum256([]byte(rawKey))
	keyHash := hex.EncodeToString(hash[:])

	// Look up by hash
	apiKey, err := s.queries.GetAPIKeyByHash(ctx, keyHash)
	if err != nil {
		return nil, ErrKeyNotFound
	}

	// Check expiration (already handled in query, but double-check)
	if apiKey.ExpiresAt.Valid && apiKey.ExpiresAt.Time.Before(time.Now()) {
		return nil, ErrKeyExpired
	}

	// Update usage statistics
	_ = s.queries.UpdateAPIKeyUsage(ctx, apiKey.ID)

	return &apiKey, nil
}

// GetByID retrieves an API key by ID.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*db.ApiKey, error) {
	apiKey, err := s.queries.GetAPIKeyByID(ctx, id)
	if err != nil {
		return nil, ErrKeyNotFound
	}
	return &apiKey, nil
}

// ListByUser returns all API keys for a user.
func (s *Service) ListByUser(ctx context.Context, userID uuid.UUID) ([]db.ApiKey, error) {
	return s.queries.ListAPIKeysByUser(ctx, userID)
}

// Delete deletes an API key.
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.queries.DeleteAPIKey(ctx, id); err != nil {
		return err
	}

	s.logger.Info("API key deleted", slog.String("key_id", id.String()))
	return nil
}

// DeleteExpired removes all expired API keys.
func (s *Service) DeleteExpired(ctx context.Context) error {
	if err := s.queries.DeleteExpiredAPIKeys(ctx); err != nil {
		return err
	}

	s.logger.Info("Expired API keys deleted")
	return nil
}

// HasScope checks if an API key has the specified scope.
func HasScope(apiKey *db.ApiKey, scope string) bool {
	for _, s := range apiKey.Scopes {
		if s == scope || s == "*" {
			return true
		}
	}
	return false
}
