package playback

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/maypok86/otter/v2"
)

// SessionManager manages active playback sessions using L1Cache for O(1) lookups.
type SessionManager struct {
	cache       *cache.L1Cache[uuid.UUID, *Session]
	maxSessions int
	timeout     time.Duration
	logger      *slog.Logger
}

// SessionCleanupFunc is called when a session is evicted or expired from cache.
// It receives the session ID and should kill associated FFmpeg processes.
type SessionCleanupFunc func(sessionID uuid.UUID)

// NewSessionManager creates a new session manager backed by L1Cache.
// The optional cleanupFn is called when sessions are evicted/expired by the cache,
// allowing the caller to kill orphaned FFmpeg processes and clean up resources.
func NewSessionManager(maxSessions int, timeout time.Duration, logger *slog.Logger, cleanupFn ...SessionCleanupFunc) (*SessionManager, error) {
	var opts []cache.L1Option[uuid.UUID, *Session]

	if len(cleanupFn) > 0 && cleanupFn[0] != nil {
		fn := cleanupFn[0]
		opts = append(opts, cache.WithOnDeletion[uuid.UUID, *Session](func(e otter.DeletionEvent[uuid.UUID, *Session]) {
			// Only run cleanup for TTL expiry and size evictions, not explicit deletes
			// (explicit deletes already handle cleanup in StopSession).
			if e.WasEvicted() {
				logger.Warn("session expired/evicted, cleaning up resources",
					slog.String("session_id", e.Key.String()),
					slog.String("reason", e.Cause.String()),
				)
				fn(e.Key)
				// Clean up segment directory
				if e.Value != nil && e.Value.SegmentDir != "" {
					go func() {
						if err := os.RemoveAll(e.Value.SegmentDir); err != nil {
							logger.Warn("failed to clean up segment dir after eviction",
								slog.String("session_id", e.Key.String()),
								slog.String("dir", e.Value.SegmentDir),
								slog.String("error", err.Error()),
							)
						}
					}()
				}
			}
		}))
	}

	c, err := cache.NewL1Cache[uuid.UUID, *Session](maxSessions*2, timeout, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create session cache: %w", err)
	}

	return &SessionManager{
		cache:       c,
		maxSessions: maxSessions,
		timeout:     timeout,
		logger:      logger,
	}, nil
}



// Create stores a new session. Returns error if max concurrent sessions exceeded.
// The active count is derived from the otter cache size, which correctly reflects
// TTL-based evictions without requiring manual bookkeeping.
func (m *SessionManager) Create(session *Session) error {
	if m.cache.Size() >= m.maxSessions {
		return fmt.Errorf("maximum concurrent sessions (%d) reached", m.maxSessions)
	}

	now := time.Now()
	session.CreatedAt = now
	session.LastAccessedAt = now
	session.ExpiresAt = now.Add(m.timeout)

	m.cache.Set(session.ID, session)

	m.logger.Info("playback session created",
		slog.String("session_id", session.ID.String()),
		slog.String("user_id", session.UserID.String()),
		slog.String("media_type", string(session.MediaType)),
		slog.String("media_id", session.MediaID.String()),
		slog.Int("active_sessions", m.cache.Size()),
	)

	return nil
}

// Get retrieves a session by ID. Returns nil, false if not found.
func (m *SessionManager) Get(id uuid.UUID) (*Session, bool) {
	return m.cache.Get(id)
}

// Touch updates the last-accessed timestamp, keeping the session alive.
func (m *SessionManager) Touch(id uuid.UUID) bool {
	session, ok := m.cache.Get(id)
	if !ok {
		return false
	}

	session.LastAccessedAt = time.Now()
	session.ExpiresAt = time.Now().Add(m.timeout)
	// Re-set refreshes the TTL in otter
	m.cache.Set(id, session)
	return true
}

// Delete removes a session.
// Returns the removed session, or nil if not found.
func (m *SessionManager) Delete(id uuid.UUID) *Session {
	session, ok := m.cache.Get(id)
	if !ok {
		return nil
	}

	m.cache.Delete(id)

	m.logger.Info("playback session deleted",
		slog.String("session_id", id.String()),
		slog.Int("active_sessions", m.cache.Size()),
	)

	return session
}

// ActiveCount returns the current number of active sessions.
func (m *SessionManager) ActiveCount() int {
	return m.cache.Size()
}

// Close shuts down the session manager and its cache.
func (m *SessionManager) Close() {
	m.cache.Close()
}
