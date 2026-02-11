package playback

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/infra/cache"
)

// SessionManager manages active playback sessions using L1Cache for O(1) lookups.
type SessionManager struct {
	cache       *cache.L1Cache[uuid.UUID, *Session]
	mu          sync.Mutex
	activeCount int
	maxSessions int
	timeout     time.Duration
	logger      *slog.Logger
}

// NewSessionManager creates a new session manager backed by L1Cache.
func NewSessionManager(maxSessions int, timeout time.Duration, logger *slog.Logger) (*SessionManager, error) {
	c, err := cache.NewL1Cache[uuid.UUID, *Session](maxSessions*2, timeout)
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
func (m *SessionManager) Create(session *Session) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.activeCount >= m.maxSessions {
		return fmt.Errorf("maximum concurrent sessions (%d) reached", m.maxSessions)
	}

	now := time.Now()
	session.CreatedAt = now
	session.LastAccessedAt = now
	session.ExpiresAt = now.Add(m.timeout)

	m.cache.Set(session.ID, session)
	m.activeCount++

	m.logger.Info("playback session created",
		slog.String("session_id", session.ID.String()),
		slog.String("user_id", session.UserID.String()),
		slog.String("media_type", string(session.MediaType)),
		slog.String("media_id", session.MediaID.String()),
		slog.Int("active_sessions", m.activeCount),
	)

	return nil
}

// Get retrieves a session by ID. Returns nil, false if not found.
func (m *SessionManager) Get(id uuid.UUID) (*Session, bool) {
	return m.cache.Get(id)
}

// Touch updates the last-accessed timestamp, keeping the session alive.
func (m *SessionManager) Touch(id uuid.UUID) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

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

// Delete removes a session and decrements the active count.
// Returns the removed session, or nil if not found.
func (m *SessionManager) Delete(id uuid.UUID) *Session {
	session, ok := m.cache.Get(id)
	if !ok {
		return nil
	}

	m.cache.Delete(id)

	m.mu.Lock()
	if m.activeCount > 0 {
		m.activeCount--
	}
	m.mu.Unlock()

	m.logger.Info("playback session deleted",
		slog.String("session_id", id.String()),
		slog.Int("active_sessions", m.activeCount),
	)

	return session
}

// ActiveCount returns the current number of active sessions.
func (m *SessionManager) ActiveCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.activeCount
}

// Close shuts down the session manager and its cache.
func (m *SessionManager) Close() {
	m.cache.Close()
}
