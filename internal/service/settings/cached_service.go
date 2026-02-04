package settings

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/infra/cache"
	"go.uber.org/zap"
)

// CachedService wraps the settings Service with caching support.
type CachedService struct {
	Service
	cache  *cache.Cache
	logger *zap.Logger
}

// NewCachedService creates a new cached settings service.
// If cache is nil, it falls back to the underlying service without caching.
func NewCachedService(svc Service, c *cache.Cache, logger *zap.Logger) *CachedService {
	return &CachedService{
		Service: svc,
		cache:   c,
		logger:  logger.Named("settings-cache"),
	}
}

// GetServerSetting gets a server setting with caching.
func (s *CachedService) GetServerSetting(ctx context.Context, key string) (*ServerSetting, error) {
	// If no cache, use underlying service
	if s.cache == nil {
		return s.Service.GetServerSetting(ctx, key)
	}

	cacheKey := cache.ServerSettingKey(key)

	// Try cache first
	var setting ServerSetting
	if err := s.cache.GetJSON(ctx, cacheKey, &setting); err == nil {
		s.logger.Debug("server setting cache hit", zap.String("key", cacheKey))
		return &setting, nil
	}

	s.logger.Debug("server setting cache miss", zap.String("key", cacheKey))

	// Cache miss - get from database
	result, err := s.Service.GetServerSetting(ctx, key)
	if err != nil {
		return nil, err
	}

	// Cache the result
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, result, cache.ServerSettingsTTL); setErr != nil {
			s.logger.Warn("failed to cache server setting", zap.Error(setErr))
		}
	}()

	return result, nil
}

// ListServerSettings lists all server settings with caching.
func (s *CachedService) ListServerSettings(ctx context.Context) ([]ServerSetting, error) {
	// If no cache, use underlying service
	if s.cache == nil {
		return s.Service.ListServerSettings(ctx)
	}

	cacheKey := cache.KeyPrefixServerSetting + "_all"

	// Try cache first
	var settings []ServerSetting
	if err := s.cache.GetJSON(ctx, cacheKey, &settings); err == nil {
		s.logger.Debug("server settings list cache hit", zap.Int("count", len(settings)))
		return settings, nil
	}

	s.logger.Debug("server settings list cache miss")

	// Cache miss - get from database
	result, err := s.Service.ListServerSettings(ctx)
	if err != nil {
		return nil, err
	}

	// Cache the result
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, result, cache.ServerSettingsTTL); setErr != nil {
			s.logger.Warn("failed to cache server settings list", zap.Error(setErr))
		}
	}()

	return result, nil
}

// ListPublicServerSettings lists public server settings with caching.
func (s *CachedService) ListPublicServerSettings(ctx context.Context) ([]ServerSetting, error) {
	// If no cache, use underlying service
	if s.cache == nil {
		return s.Service.ListPublicServerSettings(ctx)
	}

	cacheKey := cache.KeyPrefixServerSetting + "_public"

	// Try cache first
	var settings []ServerSetting
	if err := s.cache.GetJSON(ctx, cacheKey, &settings); err == nil {
		s.logger.Debug("public server settings cache hit", zap.Int("count", len(settings)))
		return settings, nil
	}

	s.logger.Debug("public server settings cache miss")

	// Cache miss - get from database
	result, err := s.Service.ListPublicServerSettings(ctx)
	if err != nil {
		return nil, err
	}

	// Cache the result
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, result, cache.ServerSettingsTTL); setErr != nil {
			s.logger.Warn("failed to cache public server settings", zap.Error(setErr))
		}
	}()

	return result, nil
}

// SetServerSetting sets a server setting and invalidates cache.
func (s *CachedService) SetServerSetting(ctx context.Context, key string, value interface{}, updatedBy uuid.UUID) (*ServerSetting, error) {
	result, err := s.Service.SetServerSetting(ctx, key, value, updatedBy)
	if err != nil {
		return nil, err
	}

	// Invalidate cache
	if s.cache != nil {
		if err := s.cache.InvalidateServerSettings(ctx); err != nil {
			s.logger.Warn("failed to invalidate server settings cache", zap.Error(err))
		}
	}

	return result, nil
}

// DeleteServerSetting deletes a server setting and invalidates cache.
func (s *CachedService) DeleteServerSetting(ctx context.Context, key string) error {
	if err := s.Service.DeleteServerSetting(ctx, key); err != nil {
		return err
	}

	// Invalidate cache
	if s.cache != nil {
		if err := s.cache.InvalidateServerSettings(ctx); err != nil {
			s.logger.Warn("failed to invalidate server settings cache", zap.Error(err))
		}
	}

	return nil
}

// GetUserSetting gets a user setting with caching.
func (s *CachedService) GetUserSetting(ctx context.Context, userID uuid.UUID, key string) (*UserSetting, error) {
	// If no cache, use underlying service
	if s.cache == nil {
		return s.Service.GetUserSetting(ctx, userID, key)
	}

	cacheKey := cache.UserSettingKey(userID.String(), key)

	// Try cache first
	var setting UserSetting
	if err := s.cache.GetJSON(ctx, cacheKey, &setting); err == nil {
		s.logger.Debug("user setting cache hit", zap.String("key", cacheKey))
		return &setting, nil
	}

	s.logger.Debug("user setting cache miss", zap.String("key", cacheKey))

	// Cache miss - get from database
	result, err := s.Service.GetUserSetting(ctx, userID, key)
	if err != nil {
		return nil, err
	}

	// Cache the result
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, result, cache.UserSettingsTTL); setErr != nil {
			s.logger.Warn("failed to cache user setting", zap.Error(setErr))
		}
	}()

	return result, nil
}

// SetUserSetting sets a user setting and invalidates cache.
func (s *CachedService) SetUserSetting(ctx context.Context, userID uuid.UUID, key string, value interface{}) (*UserSetting, error) {
	result, err := s.Service.SetUserSetting(ctx, userID, key, value)
	if err != nil {
		return nil, err
	}

	// Invalidate cache for this specific setting
	if s.cache != nil {
		cacheKey := cache.UserSettingKey(userID.String(), key)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			s.logger.Warn("failed to invalidate user setting cache", zap.Error(err))
		}
	}

	return result, nil
}

// SetUserSettingsBulk sets multiple user settings and invalidates cache.
func (s *CachedService) SetUserSettingsBulk(ctx context.Context, userID uuid.UUID, settings map[string]interface{}) error {
	if err := s.Service.SetUserSettingsBulk(ctx, userID, settings); err != nil {
		return err
	}

	// Invalidate cache for all user settings
	if s.cache != nil {
		if err := s.cache.InvalidateUserSettings(ctx, userID.String()); err != nil {
			s.logger.Warn("failed to invalidate user settings cache", zap.Error(err))
		}
	}

	return nil
}

// DeleteUserSetting deletes a user setting and invalidates cache.
func (s *CachedService) DeleteUserSetting(ctx context.Context, userID uuid.UUID, key string) error {
	if err := s.Service.DeleteUserSetting(ctx, userID, key); err != nil {
		return err
	}

	// Invalidate cache
	if s.cache != nil {
		cacheKey := cache.UserSettingKey(userID.String(), key)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			s.logger.Warn("failed to invalidate user setting cache", zap.Error(err))
		}
	}

	return nil
}
