package session

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/infra/database/db"
)

// CachedService wraps the session Service with caching support.
type CachedService struct {
	*Service
	cache    *cache.Cache
	logger   *slog.Logger
	cacheTTL time.Duration
}

// NewCachedService creates a new cached session service.
// If cache is nil, it falls back to the underlying service without caching.
func NewCachedService(svc *Service, c *cache.Cache, logger *slog.Logger, cacheTTL time.Duration) *CachedService {
	if cacheTTL == 0 {
		cacheTTL = cache.SessionTTL // Default to 30s if not specified
	}
	return &CachedService{
		Service:  svc,
		cache:    c,
		logger:   logger.With("component", "session-cache"),
		cacheTTL: cacheTTL,
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
		s.logger.Debug("session cache hit", slog.String("key", cacheKey))

		// Still update activity in background (fire and forget)
		go func() {
			actCtx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			if err := s.repo.UpdateSessionActivity(actCtx, session.ID); err != nil {
				s.logger.Warn("failed to update session activity",
					slog.String("session_id", session.ID.String()),
					slog.Any("error", err))
			}
		}()

		return &session, nil
	}

	s.logger.Debug("session cache miss", slog.String("key", cacheKey))

	// Cache miss - validate from database
	result, err := s.Service.ValidateSession(ctx, token)
	if err != nil {
		return nil, err
	}

	// Cache the session
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, result, s.cacheTTL); setErr != nil {
			s.logger.Warn("failed to cache session", slog.Any("error", setErr))
		}
	}()

	return result, nil
}

// CreateSession creates a new session with write-through caching.
// The session is written to both the database and cache simultaneously.
func (s *CachedService) CreateSession(ctx context.Context, userID uuid.UUID, deviceInfo DeviceInfo, scopes []string) (uuid.UUID, string, string, error) {
	// Create session in database
	sessionID, token, refreshToken, err := s.Service.CreateSession(ctx, userID, deviceInfo, scopes)
	if err != nil {
		return uuid.Nil, "", "", err
	}

	// Write-through: cache the new session immediately
	if s.cache != nil {
		tokenHash := s.hashToken(token)
		cacheKey := cache.SessionKey(tokenHash)

		// Get the session we just created to cache it
		session, getErr := s.repo.GetSessionByTokenHash(ctx, tokenHash)
		if getErr == nil && session != nil {
			go func() {
				cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
				defer cancel()
				if setErr := s.cache.SetJSON(cacheCtx, cacheKey, session, s.cacheTTL); setErr != nil {
					s.logger.Warn("failed to cache new session", slog.Any("error", setErr))
				} else {
					s.logger.Debug("session cached on create",
						slog.String("user_id", userID.String()),
						slog.String("key", cacheKey))
				}
			}()
		}
	}

	return sessionID, token, refreshToken, nil
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
				slog.String("session_id", sessionID.String()),
				slog.Any("error", err))
		}
	}

	return nil
}

// RevokeAllUserSessions revokes all sessions for a user and invalidates cache.
func (s *CachedService) RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	// Collect token hashes BEFORE revoking so we can invalidate the correct cache keys.
	// Sessions are cached under "session:<tokenHash>" — pattern-based invalidation
	// by user ID cannot match those keys.
	var tokenHashes []string
	if s.cache != nil {
		sessions, err := s.repo.ListUserSessions(ctx, userID)
		if err == nil {
			for _, sess := range sessions {
				tokenHashes = append(tokenHashes, sess.TokenHash)
			}
		}
	}

	// Revoke in database
	if err := s.Service.RevokeAllUserSessions(ctx, userID); err != nil {
		return err
	}

	// Invalidate each session's cache entry by its actual key
	if s.cache != nil {
		for _, hash := range tokenHashes {
			cacheKey := cache.SessionKey(hash)
			if err := s.cache.Delete(ctx, cacheKey); err != nil {
				s.logger.Warn("failed to invalidate session cache entry",
					slog.String("user_id", userID.String()),
					slog.Any("error", err))
			}
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

// RefreshSession refreshes session tokens and invalidates the old session's cache entry.
// Without this override, the old token's cached entry would remain valid until TTL expiry,
// allowing use of a revoked session token.
func (s *CachedService) RefreshSession(ctx context.Context, refreshToken string) (string, string, error) {
	// Get old session's token hash BEFORE refresh so we can invalidate its cache entry.
	refreshTokenHash := s.hashToken(refreshToken)
	oldSession, _ := s.repo.GetSessionByRefreshTokenHash(ctx, refreshTokenHash)

	// Perform the actual refresh (creates new session, revokes old).
	newToken, newRefreshToken, err := s.Service.RefreshSession(ctx, refreshToken)
	if err != nil {
		return "", "", err
	}

	// Invalidate old session cache entry.
	if s.cache != nil && oldSession != nil {
		cacheKey := cache.SessionKey(oldSession.TokenHash)
		if delErr := s.cache.Delete(ctx, cacheKey); delErr != nil {
			s.logger.Warn("failed to invalidate old session cache on refresh",
				slog.String("session_id", oldSession.ID.String()),
				slog.Any("error", delErr))
		}
	}

	// Write-through: cache the new session.
	if s.cache != nil {
		newTokenHash := s.hashToken(newToken)
		cacheKey := cache.SessionKey(newTokenHash)
		newSession, getErr := s.repo.GetSessionByTokenHash(ctx, newTokenHash)
		if getErr == nil && newSession != nil {
			go func() {
				cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
				defer cancel()
				if setErr := s.cache.SetJSON(cacheCtx, cacheKey, newSession, s.cacheTTL); setErr != nil {
					s.logger.Warn("failed to cache refreshed session", slog.Any("error", setErr))
				}
			}()
		}
	}

	return newToken, newRefreshToken, nil
}

// RevokeAllUserSessionsExcept revokes all sessions except the current one and invalidates cache.
// Without this override, revoked sessions would remain in cache until TTL expiry.
func (s *CachedService) RevokeAllUserSessionsExcept(ctx context.Context, userID uuid.UUID, currentSessionID uuid.UUID) error {
	// Collect token hashes BEFORE revoking so we can invalidate the correct cache keys.
	var tokenHashes []string
	if s.cache != nil {
		sessions, err := s.repo.ListUserSessions(ctx, userID)
		if err == nil {
			for _, sess := range sessions {
				if sess.ID != currentSessionID {
					tokenHashes = append(tokenHashes, sess.TokenHash)
				}
			}
		}
	}

	// Revoke in database.
	if err := s.Service.RevokeAllUserSessionsExcept(ctx, userID, currentSessionID); err != nil {
		return err
	}

	// Invalidate each revoked session's cache entry.
	if s.cache != nil {
		for _, hash := range tokenHashes {
			cacheKey := cache.SessionKey(hash)
			if err := s.cache.Delete(ctx, cacheKey); err != nil {
				s.logger.Warn("failed to invalidate session cache entry",
					slog.String("user_id", userID.String()),
					slog.Any("error", err))
			}
		}
	}

	return nil
}
