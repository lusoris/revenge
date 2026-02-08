package apikeys

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lusoris/revenge/internal/infra/database/db"
	"go.uber.org/zap"
)

const (
	// KeyPrefix is the prefix for all API keys
	KeyPrefix = "rv_"
	// KeyLength is the length of the random part (32 bytes = 64 hex chars)
	KeyLength = 32
	// DefaultMaxKeysPerUser is the maximum number of active keys per user
	DefaultMaxKeysPerUser = 10
)

var (
	ErrKeyNotFound      = errors.New("API key not found")
	ErrKeyInactive      = errors.New("API key is inactive")
	ErrKeyExpired       = errors.New("API key has expired")
	ErrMaxKeysExceeded  = errors.New("maximum number of API keys exceeded")
	ErrInvalidKeyFormat = errors.New("invalid API key format")
	ErrInvalidScope     = errors.New("invalid scope")
)

// Service implements API keys business logic
type Service struct {
	repo           Repository
	logger         *zap.Logger
	maxKeysPerUser int
	defaultExpiry  time.Duration // 0 = never expire
}

// NewService creates a new API keys service
func NewService(repo Repository, logger *zap.Logger, maxKeysPerUser int, defaultExpiry time.Duration) *Service {
	if maxKeysPerUser <= 0 {
		maxKeysPerUser = DefaultMaxKeysPerUser
	}
	return &Service{
		repo:           repo,
		logger:         logger,
		maxKeysPerUser: maxKeysPerUser,
		defaultExpiry:  defaultExpiry,
	}
}

// CreateKey generates a new API key for a user
func (s *Service) CreateKey(ctx context.Context, userID uuid.UUID, req CreateKeyRequest) (*CreateKeyResponse, error) {
	// Check max keys per user
	count, err := s.repo.CountUserAPIKeys(ctx, userID)
	if err != nil {
		s.logger.Error("failed to count user API keys",
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to count API keys: %w", err)
	}

	if count >= int64(s.maxKeysPerUser) {
		return nil, ErrMaxKeysExceeded
	}

	// Validate scopes
	if err := s.validateScopes(req.Scopes); err != nil {
		return nil, err
	}

	// Generate random key
	rawKey, keyHash, keyPrefix, err := s.generateKey()
	if err != nil {
		s.logger.Error("failed to generate API key",
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}

	// Use default expiry if not specified
	expiresAt := req.ExpiresAt
	if expiresAt == nil && s.defaultExpiry > 0 {
		expires := time.Now().Add(s.defaultExpiry)
		expiresAt = &expires
	}

	// Convert to pgtype.Timestamptz
	var expiresAtPg pgtype.Timestamptz
	if expiresAt != nil {
		expiresAtPg = pgtype.Timestamptz{
			Time:  *expiresAt,
			Valid: true,
		}
	}

	// Create key in database
	dbKey, err := s.repo.CreateAPIKey(ctx, db.CreateAPIKeyParams{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		KeyHash:     keyHash,
		KeyPrefix:   keyPrefix,
		Scopes:      req.Scopes,
		ExpiresAt:   expiresAtPg,
	})
	if err != nil {
		s.logger.Error("failed to create API key",
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to create key: %w", err)
	}

	return &CreateKeyResponse{
		Key:    s.dbKeyToAPIKey(dbKey),
		RawKey: rawKey,
	}, nil
}

// GetKey retrieves an API key by ID (without raw key)
func (s *Service) GetKey(ctx context.Context, keyID uuid.UUID) (*APIKey, error) {
	dbKey, err := s.repo.GetAPIKey(ctx, keyID)
	if err != nil {
		return nil, ErrKeyNotFound
	}

	key := s.dbKeyToAPIKey(dbKey)
	return &key, nil
}

// ListUserKeys lists active API keys for a user
func (s *Service) ListUserKeys(ctx context.Context, userID uuid.UUID) ([]APIKey, error) {
	dbKeys, err := s.repo.ListActiveUserAPIKeys(ctx, userID)
	if err != nil {
		s.logger.Error("failed to list user API keys",
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to list keys: %w", err)
	}

	keys := make([]APIKey, len(dbKeys))
	for i, dbKey := range dbKeys {
		keys[i] = s.dbKeyToAPIKey(dbKey)
	}

	return keys, nil
}

// ValidateKey validates a raw API key and returns the associated key data
func (s *Service) ValidateKey(ctx context.Context, rawKey string) (*APIKey, error) {
	// Check key format
	if !s.isValidKeyFormat(rawKey) {
		return nil, ErrInvalidKeyFormat
	}

	// Hash the key
	keyHash := s.hashKey(rawKey)

	// Get key from database
	dbKey, err := s.repo.GetAPIKeyByHash(ctx, keyHash)
	if err != nil {
		return nil, ErrKeyNotFound
	}

	// Check if active
	if !dbKey.IsActive {
		return nil, ErrKeyInactive
	}

	// Check expiry
	if dbKey.ExpiresAt.Valid && time.Now().After(dbKey.ExpiresAt.Time) {
		return nil, ErrKeyExpired
	}

	// Update last used timestamp (async, don't wait)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.repo.UpdateAPIKeyLastUsed(ctx, dbKey.ID); err != nil {
			s.logger.Warn("failed to update API key last used",
				zap.String("key_id", dbKey.ID.String()),
				zap.Error(err),
			)
		}
	}()

	key := s.dbKeyToAPIKey(dbKey)
	return &key, nil
}

// RevokeKey revokes an API key
func (s *Service) RevokeKey(ctx context.Context, keyID uuid.UUID) error {
	return s.repo.RevokeAPIKey(ctx, keyID)
}

// CheckScope checks if an API key has a required scope
func (s *Service) CheckScope(ctx context.Context, keyID uuid.UUID, requiredScope string) (bool, error) {
	key, err := s.GetKey(ctx, keyID)
	if err != nil {
		return false, err
	}

	for _, scope := range key.Scopes {
		if scope == requiredScope || scope == "admin" {
			return true, nil
		}
	}

	return false, nil
}

// UpdateScopes updates the scopes of an API key
func (s *Service) UpdateScopes(ctx context.Context, keyID uuid.UUID, scopes []string) error {
	if err := s.validateScopes(scopes); err != nil {
		return err
	}

	return s.repo.UpdateAPIKeyScopes(ctx, keyID, scopes)
}

// CleanupExpiredKeys deletes expired and inactive keys
func (s *Service) CleanupExpiredKeys(ctx context.Context) error {
	return s.repo.DeleteExpiredAPIKeys(ctx)
}

// ============================================================================
// Private helpers
// ============================================================================

// generateKey generates a random API key with format: rv_<64 hex chars>
func (s *Service) generateKey() (rawKey, keyHash, keyPrefix string, err error) {
	// Generate random bytes
	randomBytes := make([]byte, KeyLength)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", "", "", err
	}

	// Encode to hex
	hexRandom := hex.EncodeToString(randomBytes)

	// Format: rv_<64 hex chars>
	rawKey = KeyPrefix + hexRandom

	// Hash the key (SHA-256)
	keyHash = s.hashKey(rawKey)

	// Store first 8 chars as prefix for identification
	keyPrefix = rawKey[:8] // "rv_xxxxx"

	return rawKey, keyHash, keyPrefix, nil
}

// hashKey creates a SHA-256 hash of the key
func (s *Service) hashKey(rawKey string) string {
	hash := sha256.Sum256([]byte(rawKey))
	return hex.EncodeToString(hash[:])
}

// isValidKeyFormat checks if a key has the correct format
func (s *Service) isValidKeyFormat(key string) bool {
	// Must have prefix
	if len(key) < len(KeyPrefix) {
		return false
	}

	// Must start with prefix
	if key[:len(KeyPrefix)] != KeyPrefix {
		return false
	}

	// Must have correct length (prefix + 64 hex chars)
	expectedLen := len(KeyPrefix) + (KeyLength * 2)
	return len(key) == expectedLen
}

// validateScopes validates API key scopes
func (s *Service) validateScopes(scopes []string) error {
	validScopes := map[string]bool{
		"read":  true,
		"write": true,
		"admin": true,
	}

	for _, scope := range scopes {
		if !validScopes[scope] {
			return fmt.Errorf("%w: %s", ErrInvalidScope, scope)
		}
	}

	return nil
}

// dbKeyToAPIKey converts a database key to an API key
func (s *Service) dbKeyToAPIKey(dbKey db.SharedApiKey) APIKey {
	var expiresAt *time.Time
	if dbKey.ExpiresAt.Valid {
		expiresAt = &dbKey.ExpiresAt.Time
	}

	var lastUsedAt *time.Time
	if dbKey.LastUsedAt.Valid {
		lastUsedAt = &dbKey.LastUsedAt.Time
	}

	return APIKey{
		ID:          dbKey.ID,
		UserID:      dbKey.UserID,
		Name:        dbKey.Name,
		Description: dbKey.Description,
		KeyPrefix:   dbKey.KeyPrefix,
		Scopes:      dbKey.Scopes,
		IsActive:    dbKey.IsActive,
		ExpiresAt:   expiresAt,
		LastUsedAt:  lastUsedAt,
		CreatedAt:   dbKey.CreatedAt,
		UpdatedAt:   dbKey.UpdatedAt,
	}
}
