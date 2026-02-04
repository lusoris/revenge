package session

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// testHashToken computes the same hash as Service.hashToken for testing
func testHashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func TestNewCachedService(t *testing.T) {
	logger := zap.NewNop()
	svc := &Service{}

	cached := NewCachedService(svc, nil, logger, 5*time.Minute)

	require.NotNil(t, cached)
	assert.Equal(t, svc, cached.Service)
	assert.Nil(t, cached.cache)
	assert.Equal(t, 5*time.Minute, cached.cacheTTL)
}

func TestNewCachedService_DefaultTTL(t *testing.T) {
	logger := zap.NewNop()
	svc := &Service{}

	// When TTL is 0, should use default SessionTTL
	cached := NewCachedService(svc, nil, logger, 0)

	require.NotNil(t, cached)
	assert.Equal(t, cache.SessionTTL, cached.cacheTTL)
}

func TestCachedService_ValidateSession_NoCache(t *testing.T) {
	// Create a mock repository that returns a session
	mockRepo := &mockRepository{
		session: db.SharedSession{
			ID:        uuid.New(),
			UserID:    uuid.New(),
			TokenHash: "test-token-hash",
		},
		callCount: make(map[string]int),
	}

	svc := &Service{
		repo:   mockRepo,
		logger: zap.NewNop(),
	}

	cached := NewCachedService(svc, nil, zap.NewNop(), 5*time.Minute)

	// Without cache, should still work
	session, err := cached.ValidateSession(context.Background(), "test-token")
	require.NoError(t, err)
	assert.NotNil(t, session)
}

func TestCachedService_ValidateSession_WithCache(t *testing.T) {
	// Create L1 cache for testing - use short TTL so SessionTTL (30s) is cached in L1
	l1Cache, err := cache.NewCache(nil, 1000, 15*time.Second)
	require.NoError(t, err)
	defer l1Cache.Close()

	mockRepo := &mockRepository{
		session: db.SharedSession{
			ID:        uuid.New(),
			UserID:    uuid.New(),
			TokenHash: "cached-token-hash",
		},
		callCount: make(map[string]int),
	}

	svc := &Service{
		repo:   mockRepo,
		logger: zap.NewNop(),
	}

	cached := NewCachedService(svc, l1Cache, zap.NewNop(), 5*time.Minute)

	// First call - cache miss
	session1, err := cached.ValidateSession(context.Background(), "cached-token")
	require.NoError(t, err)
	assert.NotNil(t, session1)

	// Give time for async cache set
	time.Sleep(50 * time.Millisecond)

	// Second call should hit cache
	session2, err := cached.ValidateSession(context.Background(), "cached-token")
	require.NoError(t, err)
	assert.NotNil(t, session2)

	// Repository should only be called once for GetSessionByTokenHash
	assert.Equal(t, 1, mockRepo.callCount["GetSessionByTokenHash"])
}

func TestCachedService_RevokeSession_InvalidatesCache(t *testing.T) {
	l1Cache, err := cache.NewCache(nil, 1000, time.Minute)
	require.NoError(t, err)
	defer l1Cache.Close()

	sessionID := uuid.New()
	// Use the actual hash of the token so cache invalidation works correctly
	tokenHash := testHashToken("revoke-token")
	mockRepo := &mockRepository{
		session: db.SharedSession{
			ID:        sessionID,
			UserID:    uuid.New(),
			TokenHash: tokenHash,
		},
		callCount: make(map[string]int),
	}

	svc := &Service{
		repo:   mockRepo,
		logger: zap.NewNop(),
	}

	cached := NewCachedService(svc, l1Cache, zap.NewNop(), 5*time.Minute)

	// Populate cache
	_, err = cached.ValidateSession(context.Background(), "revoke-token")
	require.NoError(t, err)
	time.Sleep(50 * time.Millisecond)

	// Revoke session
	err = cached.RevokeSession(context.Background(), sessionID)
	require.NoError(t, err)

	// Next call should miss cache (invalidated)
	mockRepo.callCount["GetSessionByTokenHash"] = 0
	_, err = cached.ValidateSession(context.Background(), "revoke-token")
	require.NoError(t, err)

	// Should have called repository again
	assert.Equal(t, 1, mockRepo.callCount["GetSessionByTokenHash"])
}

func TestCachedService_InvalidateSessionCache(t *testing.T) {
	l1Cache, err := cache.NewCache(nil, 1000, time.Minute)
	require.NoError(t, err)
	defer l1Cache.Close()

	cached := &CachedService{
		cache:    l1Cache,
		logger:   zap.NewNop(),
		cacheTTL: 5 * time.Minute,
	}

	// Should not error even with valid token hash
	err = cached.InvalidateSessionCache(context.Background(), "some-token-hash")
	require.NoError(t, err)

	// Should not error with nil cache
	cached.cache = nil
	err = cached.InvalidateSessionCache(context.Background(), "some-token-hash")
	require.NoError(t, err)
}

// mockRepository is a test mock for the session repository
type mockRepository struct {
	session   db.SharedSession
	callCount map[string]int
}

func (m *mockRepository) CreateSession(ctx context.Context, params CreateSessionParams) (db.SharedSession, error) {
	m.callCount["CreateSession"]++
	return m.session, nil
}

func (m *mockRepository) GetSessionByTokenHash(ctx context.Context, tokenHash string) (*db.SharedSession, error) {
	m.callCount["GetSessionByTokenHash"]++
	return &m.session, nil
}

func (m *mockRepository) GetSessionByID(ctx context.Context, sessionID uuid.UUID) (*db.SharedSession, error) {
	m.callCount["GetSessionByID"]++
	return &m.session, nil
}

func (m *mockRepository) GetSessionByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (*db.SharedSession, error) {
	m.callCount["GetSessionByRefreshTokenHash"]++
	return &m.session, nil
}

func (m *mockRepository) ListUserSessions(ctx context.Context, userID uuid.UUID) ([]db.SharedSession, error) {
	m.callCount["ListUserSessions"]++
	return nil, nil
}

func (m *mockRepository) ListAllUserSessions(ctx context.Context, userID uuid.UUID) ([]db.SharedSession, error) {
	m.callCount["ListAllUserSessions"]++
	return nil, nil
}

func (m *mockRepository) CountActiveUserSessions(ctx context.Context, userID uuid.UUID) (int64, error) {
	m.callCount["CountActiveUserSessions"]++
	return 0, nil
}

func (m *mockRepository) UpdateSessionActivity(ctx context.Context, sessionID uuid.UUID) error {
	m.callCount["UpdateSessionActivity"]++
	return nil
}

func (m *mockRepository) UpdateSessionActivityByTokenHash(ctx context.Context, tokenHash string) error {
	m.callCount["UpdateSessionActivityByTokenHash"]++
	return nil
}

func (m *mockRepository) RevokeSession(ctx context.Context, sessionID uuid.UUID, reason *string) error {
	m.callCount["RevokeSession"]++
	return nil
}

func (m *mockRepository) RevokeSessionByTokenHash(ctx context.Context, tokenHash string, reason *string) error {
	m.callCount["RevokeSessionByTokenHash"]++
	return nil
}

func (m *mockRepository) RevokeAllUserSessions(ctx context.Context, userID uuid.UUID, reason *string) error {
	m.callCount["RevokeAllUserSessions"]++
	return nil
}

func (m *mockRepository) RevokeAllUserSessionsExcept(ctx context.Context, userID uuid.UUID, exceptID uuid.UUID, reason *string) error {
	m.callCount["RevokeAllUserSessionsExcept"]++
	return nil
}

func (m *mockRepository) DeleteExpiredSessions(ctx context.Context) (int64, error) {
	m.callCount["DeleteExpiredSessions"]++
	return 0, nil
}

func (m *mockRepository) DeleteRevokedSessions(ctx context.Context) (int64, error) {
	m.callCount["DeleteRevokedSessions"]++
	return 0, nil
}

func (m *mockRepository) GetInactiveSessions(ctx context.Context, inactiveSince time.Time) ([]db.SharedSession, error) {
	m.callCount["GetInactiveSessions"]++
	return nil, nil
}

func (m *mockRepository) RevokeInactiveSessions(ctx context.Context, inactiveSince time.Time) error {
	m.callCount["RevokeInactiveSessions"]++
	return nil
}
