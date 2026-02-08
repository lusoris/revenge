package user

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/infra/database/db"
)

// CachedService wraps the user service with caching.
type CachedService struct {
	*Service
	cache  *cache.Cache
	logger *slog.Logger
}

// NewCachedService creates a new cached user service.
func NewCachedService(svc *Service, c *cache.Cache, logger *slog.Logger) *CachedService {
	return &CachedService{
		Service: svc,
		cache:   c,
		logger:  logger.With("component", "user-cache"),
	}
}

// GetUser retrieves a user by ID with caching (1 min TTL).
func (s *CachedService) GetUser(ctx context.Context, userID uuid.UUID) (*db.SharedUser, error) {
	if s.cache == nil {
		return s.Service.GetUser(ctx, userID)
	}

	cacheKey := cache.UserKey(userID.String())

	var user db.SharedUser
	if err := s.cache.GetJSON(ctx, cacheKey, &user); err == nil {
		s.logger.Debug("user cache hit", slog.String("id", userID.String()))
		return &user, nil
	}

	s.logger.Debug("user cache miss", slog.String("id", userID.String()))

	result, err := s.Service.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Cache async
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), cache.SessionTTL)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, result, cache.UserTTL); setErr != nil {
			s.logger.Warn("failed to cache user", slog.Any("error",setErr))
		}
	}()

	return result, nil
}

// GetUserByUsername retrieves a user by username with caching.
func (s *CachedService) GetUserByUsername(ctx context.Context, username string) (*db.SharedUser, error) {
	if s.cache == nil {
		return s.Service.GetUserByUsername(ctx, username)
	}

	cacheKey := cache.UserByNameKey(username)

	var user db.SharedUser
	if err := s.cache.GetJSON(ctx, cacheKey, &user); err == nil {
		s.logger.Debug("user by name cache hit", slog.String("username", username))
		return &user, nil
	}

	result, err := s.Service.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	// Cache async
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), cache.SessionTTL)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, result, cache.UserTTL); setErr != nil {
			s.logger.Warn("failed to cache user by name", slog.Any("error",setErr))
		}
	}()

	return result, nil
}

// Write operations - invalidate cache

// UpdateUser updates a user and invalidates cache.
func (s *CachedService) UpdateUser(ctx context.Context, userID uuid.UUID, params UpdateUserParams) (*db.SharedUser, error) {
	result, err := s.Service.UpdateUser(ctx, userID, params)
	if err != nil {
		return nil, err
	}

	// Invalidate cache
	if s.cache != nil {
		go func() {
			cacheCtx, cancel := context.WithTimeout(context.Background(), cache.SessionTTL)
			defer cancel()
			if err := s.cache.InvalidateUser(cacheCtx, userID.String()); err != nil {
				s.logger.Warn("failed to invalidate user cache", slog.Any("error",err))
			}
		}()
	}

	return result, nil
}

// DeleteUser deletes a user and invalidates cache.
func (s *CachedService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	// Get user first to invalidate by username too
	user, _ := s.Service.GetUser(ctx, userID)

	err := s.Service.DeleteUser(ctx, userID)
	if err != nil {
		return err
	}

	// Invalidate cache
	if s.cache != nil {
		go func() {
			cacheCtx, cancel := context.WithTimeout(context.Background(), cache.SessionTTL)
			defer cancel()
			if err := s.cache.InvalidateUser(cacheCtx, userID.String()); err != nil {
				s.logger.Warn("failed to invalidate user cache", slog.Any("error",err))
			}
			if user != nil {
				if err := s.cache.Delete(cacheCtx, cache.UserByNameKey(user.Username)); err != nil {
					s.logger.Warn("failed to invalidate user by name cache", slog.Any("error",err))
				}
			}
		}()
	}

	return nil
}

// InvalidateUserCache invalidates all cache entries for a user.
func (s *CachedService) InvalidateUserCache(ctx context.Context, userID uuid.UUID) error {
	if s.cache == nil {
		return nil
	}
	return s.cache.InvalidateUser(ctx, userID.String())
}
