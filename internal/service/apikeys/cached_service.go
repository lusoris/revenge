package apikeys

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"github.com/lusoris/revenge/internal/infra/cache"
)

// Ensure CachedService implements Service.
var _ Service = (*CachedService)(nil)

// CachedService wraps the API keys Service with L1+L2 caching for ValidateKey.
// ValidateKey is on the authentication hot path — every API-key-authenticated
// request calls it. Caching avoids a DB round-trip per request.
//
// Cache invalidation strategy:
//   - RevokeKey: deletes the cached entry by key hash (requires DB lookup first).
//   - UpdateScopes: same — invalidate after scope change.
//   - Short TTL (30s): limits stale-key exposure window after revocation
//     from a code path that doesn't have the hash (defense in depth).
type CachedService struct {
	Service // embed the interface (delegates all uncached methods)
	cache   *cache.Cache
	logger  *slog.Logger
	repo    Repository // needed to look up key hash on revoke/update
}

// NewCachedService creates a cached API keys service.
// If c is nil, all calls pass through to the underlying service unchanged.
func NewCachedService(svc Service, repo Repository, c *cache.Cache, logger *slog.Logger) *CachedService {
	return &CachedService{
		Service: svc,
		cache:   c,
		logger:  logger.With("component", "apikeys-cache"),
		repo:    repo,
	}
}

// ValidateKey validates a raw API key with cache-aside.
// Cache key is the SHA-256 hash of the raw key (same hash the DB uses),
// so we never store the raw key in cache.
func (s *CachedService) ValidateKey(ctx context.Context, rawKey string) (*APIKey, error) {
	// Format validation first — cheap, avoids cache/DB for garbage input.
	if len(rawKey) < len(KeyPrefix) || rawKey[:len(KeyPrefix)] != KeyPrefix {
		return nil, ErrInvalidKeyFormat
	}
	expectedLen := len(KeyPrefix) + (KeyLength * 2)
	if len(rawKey) != expectedLen {
		return nil, ErrInvalidKeyFormat
	}

	keyHash := hashRawKey(rawKey)
	cacheKey := cache.APIKeyByHashKey(keyHash)

	return cache.Get(ctx, s.cache, cacheKey, cache.APIKeyTTL, func(ctx context.Context) (*APIKey, error) {
		return s.Service.ValidateKey(ctx, rawKey)
	})
}

// RevokeKey revokes an API key and invalidates the cache entry.
func (s *CachedService) RevokeKey(ctx context.Context, keyID uuid.UUID) error {
	// Look up the key hash BEFORE revoking so we can invalidate the cache.
	dbKey, lookupErr := s.repo.GetAPIKey(ctx, keyID)

	if err := s.Service.RevokeKey(ctx, keyID); err != nil {
		return err
	}

	// Invalidate cache if we got the hash.
	if lookupErr == nil {
		s.invalidateKeyCache(ctx, dbKey.KeyHash)
	}

	return nil
}

// UpdateScopes updates API key scopes and invalidates the cache.
func (s *CachedService) UpdateScopes(ctx context.Context, keyID uuid.UUID, scopes []string) error {
	// Look up the key hash BEFORE updating so we can invalidate the cache.
	dbKey, lookupErr := s.repo.GetAPIKey(ctx, keyID)

	if err := s.Service.UpdateScopes(ctx, keyID, scopes); err != nil {
		return err
	}

	// Invalidate cache if we got the hash.
	if lookupErr == nil {
		s.invalidateKeyCache(ctx, dbKey.KeyHash)
	}

	return nil
}

// invalidateKeyCache removes a validated key from the cache by its hash.
func (s *CachedService) invalidateKeyCache(ctx context.Context, keyHash string) {
	cacheKey := cache.APIKeyByHashKey(keyHash)
	if err := s.cache.Delete(ctx, cacheKey); err != nil {
		s.logger.Warn("failed to invalidate API key cache",
			slog.String("cache_key", cacheKey),
			slog.Any("error", err),
		)
	}
}

// hashRawKey delegates to the shared hashAPIKey function.
func hashRawKey(rawKey string) string {
	return hashAPIKey(rawKey)
}
