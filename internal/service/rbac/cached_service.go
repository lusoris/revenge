package rbac

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/infra/cache"
	"go.uber.org/zap"
)

// CachedService wraps the RBAC Service with caching support.
type CachedService struct {
	*Service
	cache  *cache.Cache
	logger *zap.Logger
}

// NewCachedService creates a new cached RBAC service.
// If cache is nil, it falls back to the underlying service without caching.
func NewCachedService(svc *Service, cache *cache.Cache, logger *zap.Logger) *CachedService {
	return &CachedService{
		Service: svc,
		cache:   cache,
		logger:  logger.Named("rbac-cache"),
	}
}

// cachedBool is a wrapper for caching boolean values.
type cachedBool struct {
	Value bool `json:"value"`
}

// cachedStringSlice is a wrapper for caching string slices.
type cachedStringSlice struct {
	Values []string `json:"values"`
}

// Enforce checks if a subject has permission with caching.
func (s *CachedService) Enforce(ctx context.Context, sub, obj, act string) (bool, error) {
	// If no cache, use underlying service
	if s.cache == nil {
		return s.Service.Enforce(ctx, sub, obj, act)
	}

	cacheKey := cache.RBACEnforceKey(sub, obj, act)

	// Try cache first
	var result cachedBool
	if err := s.cache.GetJSON(ctx, cacheKey, &result); err == nil {
		s.logger.Debug("RBAC enforce cache hit",
			zap.String("key", cacheKey),
			zap.Bool("allowed", result.Value))
		return result.Value, nil
	}

	s.logger.Debug("RBAC enforce cache miss", zap.String("key", cacheKey))

	// Cache miss - check policy
	allowed, err := s.Service.Enforce(ctx, sub, obj, act)
	if err != nil {
		return false, err
	}

	// Cache the result
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, cachedBool{Value: allowed}, cache.RBACEnforceTTL); setErr != nil {
			s.logger.Warn("failed to cache RBAC enforce result", zap.Error(setErr))
		}
	}()

	return allowed, nil
}

// EnforceWithContext checks if a user has permission with caching.
func (s *CachedService) EnforceWithContext(ctx context.Context, userID uuid.UUID, resource, action string) (bool, error) {
	return s.Enforce(ctx, userID.String(), resource, action)
}

// GetUserRoles returns all roles for a user with caching.
func (s *CachedService) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]string, error) {
	// If no cache, use underlying service
	if s.cache == nil {
		return s.Service.GetUserRoles(ctx, userID)
	}

	cacheKey := cache.RBACUserRolesKey(userID.String())

	// Try cache first
	var result cachedStringSlice
	if err := s.cache.GetJSON(ctx, cacheKey, &result); err == nil {
		s.logger.Debug("RBAC roles cache hit",
			zap.String("key", cacheKey),
			zap.Strings("roles", result.Values))
		return result.Values, nil
	}

	s.logger.Debug("RBAC roles cache miss", zap.String("key", cacheKey))

	// Cache miss - get from Casbin
	roles, err := s.Service.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Cache the result
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, cachedStringSlice{Values: roles}, cache.RBACPolicyTTL); setErr != nil {
			s.logger.Warn("failed to cache RBAC roles", zap.Error(setErr))
		}
	}()

	return roles, nil
}

// HasRole checks if a user has a specific role with caching.
func (s *CachedService) HasRole(ctx context.Context, userID uuid.UUID, role string) (bool, error) {
	// If no cache, use underlying service
	if s.cache == nil {
		return s.Service.HasRole(ctx, userID, role)
	}

	cacheKey := cache.RBACEnforceKey(userID.String(), "role", role)

	// Try cache first
	var result cachedBool
	if err := s.cache.GetJSON(ctx, cacheKey, &result); err == nil {
		s.logger.Debug("RBAC has role cache hit",
			zap.String("key", cacheKey),
			zap.Bool("hasRole", result.Value))
		return result.Value, nil
	}

	s.logger.Debug("RBAC has role cache miss", zap.String("key", cacheKey))

	// Cache miss - check from Casbin
	hasRole, err := s.Service.HasRole(ctx, userID, role)
	if err != nil {
		return false, err
	}

	// Cache the result
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, cachedBool{Value: hasRole}, cache.RBACPolicyTTL); setErr != nil {
			s.logger.Warn("failed to cache RBAC has role result", zap.Error(setErr))
		}
	}()

	return hasRole, nil
}

// AssignRole assigns a role to a user and invalidates cache.
func (s *CachedService) AssignRole(ctx context.Context, userID uuid.UUID, role string) error {
	if err := s.Service.AssignRole(ctx, userID, role); err != nil {
		return err
	}

	// Invalidate cache for this user
	s.invalidateUserCache(ctx, userID)
	return nil
}

// RemoveRole removes a role from a user and invalidates cache.
func (s *CachedService) RemoveRole(ctx context.Context, userID uuid.UUID, role string) error {
	if err := s.Service.RemoveRole(ctx, userID, role); err != nil {
		return err
	}

	// Invalidate cache for this user
	s.invalidateUserCache(ctx, userID)
	return nil
}

// AddPolicy adds a policy rule and invalidates related caches.
func (s *CachedService) AddPolicy(ctx context.Context, sub, obj, act string) error {
	if err := s.Service.AddPolicy(ctx, sub, obj, act); err != nil {
		return err
	}

	// Invalidate all RBAC caches since policies affect everyone
	s.invalidateAllRBAC(ctx)
	return nil
}

// RemovePolicy removes a policy rule and invalidates related caches.
func (s *CachedService) RemovePolicy(ctx context.Context, sub, obj, act string) error {
	if err := s.Service.RemovePolicy(ctx, sub, obj, act); err != nil {
		return err
	}

	// Invalidate all RBAC caches since policies affect everyone
	s.invalidateAllRBAC(ctx)
	return nil
}

// LoadPolicy reloads policies from the database and invalidates all caches.
func (s *CachedService) LoadPolicy(ctx context.Context) error {
	if err := s.Service.LoadPolicy(ctx); err != nil {
		return err
	}

	// Invalidate all RBAC caches
	s.invalidateAllRBAC(ctx)
	return nil
}

// invalidateUserCache invalidates RBAC cache entries for a specific user.
func (s *CachedService) invalidateUserCache(ctx context.Context, userID uuid.UUID) {
	if s.cache == nil {
		return
	}

	if err := s.cache.InvalidateRBACForUser(ctx, userID.String()); err != nil {
		s.logger.Warn("failed to invalidate RBAC cache for user",
			zap.String("user_id", userID.String()),
			zap.Error(err))
	}
}

// invalidateAllRBAC invalidates all RBAC cache entries.
func (s *CachedService) invalidateAllRBAC(ctx context.Context) {
	if s.cache == nil {
		return
	}

	if err := s.cache.InvalidateAllRBAC(ctx); err != nil {
		s.logger.Warn("failed to invalidate all RBAC cache", zap.Error(err))
	}
}
