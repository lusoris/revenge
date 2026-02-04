package session

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/infra/database/db"
	"go.uber.org/zap"
)

// CachedService wraps the session Service with caching support.
type CachedService struct {
	*Service
	cache  *cache.Cache
	logger *zap.Logger
}

// NewCachedService creates a new cached session service.
// If cache is nil, it falls back to the underlying service without caching.
func NewCachedService(svc *Service, cache *cache.Cache, logger *zap.Logger) *CachedService {
	return &CachedService{
		Service: svc,
		cache:   cache,
		logger:  logger.Named("session-cache"),
	}
}

// ValidateSession validates a session token with caching.
// Session lookups are cached for a short period to reduce database load.
func (s *CachedService) ValidateSession(ctx context.Context, token string) (*db.SharedSession, error) {
	// If no cache, use underlying service
	if s.cache == nil {
		return s.Service.ValidateSession(ctx, token)
	}

	tokenHash := s.hashToken(token)
	cacheKey := cache.SessionKey(tokenHash)

	// Try cache first
	var session db.SharedSession
	if err := s.cache.GetJSON(ctx, cacheKey, &session); err == nil {
		s.logger.Debug("session cache hit", zap.String("key", cacheKey))

		// Still update activity in background (fire and forget)
		go func() {
			actCtx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			if err := s.repo.UpdateSessionActivity(actCtx, session.ID); err != nil {
				s.logger.Warn("failed to update session activity",
					zap.String("session_id", session.ID.String()),
					zap.Error(err))
			}
		}()

		return &session, nil
	}

	s.logger.Debug("session cache miss", zap.String("key", cacheKey))

	// Cache miss - validate from database
	result, err := s.Service.ValidateSession(ctx, token)
	if err != nil {
		return nil, err
	}

	// Cache the session
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, result, cache.SessionTTL); setErr != nil {
			s.logger.Warn("failed to cache session", zap.Error(setErr))
		}
	}()

	return result, nil
}

// RevokeSession revokes a session and invalidates cache.
func (s *CachedService) RevokeSession(ctx context.Context, sessionID uuid.UUID) error {
	// First get the session to find the token hash for cache invalidation
	session, err := s.repo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return err
	}

	// Revoke in database
	if err := s.Service.RevokeSession(ctx, sessionID); err != nil {
		return err
	}

	// Invalidate cache
	if s.cache != nil && session != nil {
		cacheKey := cache.SessionKey(session.TokenHash)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			s.logger.Warn("failed to invalidate session cache",
				zap.String("session_id", sessionID.String()),
				zap.Error(err))
		}
	}

	return nil
}

// RevokeAllUserSessions revokes all sessions for a user and invalidates cache.
func (s *CachedService) RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	// Revoke in database
	if err := s.Service.RevokeAllUserSessions(ctx, userID); err != nil {
		return err
	}

	// Invalidate all user sessions from cache
	if s.cache != nil {
		if err := s.cache.InvalidateUserSessions(ctx, userID.String()); err != nil {
			s.logger.Warn("failed to invalidate user sessions cache",
				zap.String("user_id", userID.String()),
				zap.Error(err))
		}
	}

	return nil
}

// InvalidateSessionCache invalidates a specific session from cache.
// Useful when session state changes (e.g., MFA verification).
func (s *CachedService) InvalidateSessionCache(ctx context.Context, tokenHash string) error {
	if s.cache == nil {
		return nil
	}
	return s.cache.InvalidateSession(ctx, tokenHash)
}
