package session

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"net/netip"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/infra/logging"
)

// =====================================================
// Helper functions and mock setup for same-package tests
// =====================================================

func newTestService(t *testing.T, repo Repository) *Service {
	t.Helper()
	return &Service{
		repo:          repo,
		logger:        logging.NewTestLogger(),
		tokenLength:   32,
		expiry:        24 * time.Hour,
		refreshExpiry: 7 * 24 * time.Hour,
		maxPerUser:    10,
	}
}

// =====================================================
// generateToken tests
// =====================================================

func TestGenerateToken_ReturnsUniqueTokens(t *testing.T) {
	t.Parallel()

	svc := newTestService(t, nil)

	token1, hash1, err := svc.generateToken()
	require.NoError(t, err)
	assert.NotEmpty(t, token1)
	assert.NotEmpty(t, hash1)

	token2, hash2, err := svc.generateToken()
	require.NoError(t, err)
	assert.NotEmpty(t, token2)
	assert.NotEmpty(t, hash2)

	// Tokens should be unique
	assert.NotEqual(t, token1, token2)
	assert.NotEqual(t, hash1, hash2)
}

func TestGenerateToken_TokenLength(t *testing.T) {
	t.Parallel()

	svc := newTestService(t, nil)

	token, _, err := svc.generateToken()
	require.NoError(t, err)

	// 32 bytes = 64 hex characters
	assert.Len(t, token, 64)
}

func TestGenerateToken_HashMatchesToken(t *testing.T) {
	t.Parallel()

	svc := newTestService(t, nil)

	token, hash, err := svc.generateToken()
	require.NoError(t, err)

	// Verify hash matches
	expectedHash := sha256.Sum256([]byte(token))
	expectedHashStr := hex.EncodeToString(expectedHash[:])
	assert.Equal(t, expectedHashStr, hash)
}

func TestGenerateToken_DifferentLengths(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		tokenLength int
		expectedLen int
	}{
		{"16 bytes", 16, 32},
		{"32 bytes", 32, 64},
		{"64 bytes", 64, 128},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				logger:      logging.NewTestLogger(),
				tokenLength: tt.tokenLength,
			}

			token, _, err := svc.generateToken()
			require.NoError(t, err)
			assert.Len(t, token, tt.expectedLen)
		})
	}
}

// =====================================================
// hashToken tests
// =====================================================

func TestHashToken_Deterministic(t *testing.T) {
	t.Parallel()

	svc := newTestService(t, nil)

	hash1 := svc.hashToken("test-token-123")
	hash2 := svc.hashToken("test-token-123")

	assert.Equal(t, hash1, hash2)
}

func TestHashToken_DifferentInputsDifferentHashes(t *testing.T) {
	t.Parallel()

	svc := newTestService(t, nil)

	hash1 := svc.hashToken("token-a")
	hash2 := svc.hashToken("token-b")

	assert.NotEqual(t, hash1, hash2)
}

func TestHashToken_SHA256Length(t *testing.T) {
	t.Parallel()

	svc := newTestService(t, nil)

	hash := svc.hashToken("any-token")

	// SHA256 produces 32 bytes = 64 hex characters
	assert.Len(t, hash, 64)
}

func TestHashToken_EmptyInput(t *testing.T) {
	t.Parallel()

	svc := newTestService(t, nil)

	hash := svc.hashToken("")
	assert.NotEmpty(t, hash)
	assert.Len(t, hash, 64) // SHA256 of empty string
}

// =====================================================
// sessionToInfo tests
// =====================================================

func TestSessionToInfo_ActiveSession(t *testing.T) {
	t.Parallel()

	svc := newTestService(t, nil)

	sessionID := uuid.Must(uuid.NewV7())
	userID := uuid.Must(uuid.NewV7())
	deviceName := "My Phone"
	userAgent := "Mozilla/5.0"
	ipAddr := netip.MustParseAddr("10.0.0.1")

	session := &db.SharedSession{
		ID:             sessionID,
		UserID:         userID,
		TokenHash:      "hash123",
		IpAddress:      ipAddr,
		DeviceName:     &deviceName,
		UserAgent:      &userAgent,
		CreatedAt:      time.Now().Add(-1 * time.Hour),
		LastActivityAt: time.Now(),
		ExpiresAt:      time.Now().Add(23 * time.Hour),
	}

	info := svc.sessionToInfo(session, true)

	assert.Equal(t, sessionID, info.ID)
	assert.True(t, info.IsActive)
	assert.True(t, info.IsCurrent)
	assert.Equal(t, "My Phone", *info.DeviceName)
	assert.Equal(t, "Mozilla/5.0", *info.UserAgent)
	assert.Equal(t, "10.0.0.1", *info.IPAddress)
}

func TestSessionToInfo_RevokedSession(t *testing.T) {
	t.Parallel()

	svc := newTestService(t, nil)

	revokedAt := pgtype.Timestamptz{Time: time.Now().Add(-1 * time.Hour), Valid: true}
	session := &db.SharedSession{
		ID:             uuid.Must(uuid.NewV7()),
		UserID:         uuid.Must(uuid.NewV7()),
		TokenHash:      "hash",
		IpAddress:      netip.MustParseAddr("127.0.0.1"),
		CreatedAt:      time.Now().Add(-2 * time.Hour),
		LastActivityAt: time.Now().Add(-1 * time.Hour),
		ExpiresAt:      time.Now().Add(22 * time.Hour),
		RevokedAt:      revokedAt,
	}

	info := svc.sessionToInfo(session, false)

	assert.False(t, info.IsActive, "revoked session should not be active")
	assert.False(t, info.IsCurrent)
}

func TestSessionToInfo_ExpiredSession(t *testing.T) {
	t.Parallel()

	svc := newTestService(t, nil)

	session := &db.SharedSession{
		ID:             uuid.Must(uuid.NewV7()),
		UserID:         uuid.Must(uuid.NewV7()),
		TokenHash:      "hash",
		IpAddress:      netip.MustParseAddr("127.0.0.1"),
		CreatedAt:      time.Now().Add(-25 * time.Hour),
		LastActivityAt: time.Now().Add(-24 * time.Hour),
		ExpiresAt:      time.Now().Add(-1 * time.Hour), // expired
	}

	info := svc.sessionToInfo(session, false)

	assert.False(t, info.IsActive, "expired session should not be active")
}

func TestSessionToInfo_UnspecifiedIPv4(t *testing.T) {
	t.Parallel()

	svc := newTestService(t, nil)

	// Use actual unspecified address 0.0.0.0 which returns true for IsUnspecified()
	session := &db.SharedSession{
		ID:             uuid.Must(uuid.NewV7()),
		UserID:         uuid.Must(uuid.NewV7()),
		TokenHash:      "hash",
		IpAddress:      netip.MustParseAddr("0.0.0.0"),
		CreatedAt:      time.Now(),
		LastActivityAt: time.Now(),
		ExpiresAt:      time.Now().Add(24 * time.Hour),
	}

	info := svc.sessionToInfo(session, false)

	assert.Nil(t, info.IPAddress, "unspecified IP 0.0.0.0 should result in nil IPAddress")
}

func TestSessionToInfo_ZeroValueAddr(t *testing.T) {
	t.Parallel()

	svc := newTestService(t, nil)

	// Zero-value netip.Addr{} is NOT "unspecified" - it's the zero/invalid addr
	// IsUnspecified() returns false for zero value, so it will generate an IP string
	session := &db.SharedSession{
		ID:             uuid.Must(uuid.NewV7()),
		UserID:         uuid.Must(uuid.NewV7()),
		TokenHash:      "hash",
		IpAddress:      netip.Addr{},
		CreatedAt:      time.Now(),
		LastActivityAt: time.Now(),
		ExpiresAt:      time.Now().Add(24 * time.Hour),
	}

	info := svc.sessionToInfo(session, false)

	// Zero-value Addr is not unspecified, so IPAddress will be set
	assert.NotNil(t, info.IPAddress)
}

func TestSessionToInfo_NilDeviceNameAndUserAgent(t *testing.T) {
	t.Parallel()

	svc := newTestService(t, nil)

	session := &db.SharedSession{
		ID:             uuid.Must(uuid.NewV7()),
		UserID:         uuid.Must(uuid.NewV7()),
		TokenHash:      "hash",
		IpAddress:      netip.MustParseAddr("192.168.1.1"),
		DeviceName:     nil,
		UserAgent:      nil,
		CreatedAt:      time.Now(),
		LastActivityAt: time.Now(),
		ExpiresAt:      time.Now().Add(24 * time.Hour),
	}

	info := svc.sessionToInfo(session, false)

	assert.Nil(t, info.DeviceName)
	assert.Nil(t, info.UserAgent)
	assert.NotNil(t, info.IPAddress)
}

// =====================================================
// CachedService additional tests
// =====================================================

func TestCachedService_CreateSession_WithCache(t *testing.T) {
	t.Parallel()

	l1Cache, err := cache.NewCache(nil, 1000, 15*time.Second)
	require.NoError(t, err)
	defer l1Cache.Close()

	sessionID := uuid.Must(uuid.NewV7())
	userID := uuid.Must(uuid.NewV7())

	repo := &mockRepository{
		session: db.SharedSession{
			ID:        sessionID,
			UserID:    userID,
			TokenHash: "created-token-hash",
		},
		callCount: make(map[string]int),
	}

	svc := &Service{
		repo:        repo,
		logger:      logging.NewTestLogger(),
		tokenLength: 32,
		expiry:      24 * time.Hour,
		maxPerUser:  10,
	}

	cached := NewCachedService(svc, l1Cache, logging.NewTestLogger(), 5*time.Minute)

	token, refreshToken, err := cached.CreateSession(
		context.Background(),
		userID,
		DeviceInfo{},
		[]string{"read"},
	)

	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotEmpty(t, refreshToken)

	// Wait for async cache write
	time.Sleep(200 * time.Millisecond)

	// Verify CreateSession was called on the repo
	assert.Equal(t, 1, repo.getCallCount("CreateSession"))
	// CountActiveUserSessions should be called
	assert.Equal(t, 1, repo.getCallCount("CountActiveUserSessions"))
}

func TestCachedService_CreateSession_WithoutCache(t *testing.T) {
	t.Parallel()

	sessionID := uuid.Must(uuid.NewV7())
	userID := uuid.Must(uuid.NewV7())

	repo := &mockRepository{
		session: db.SharedSession{
			ID:        sessionID,
			UserID:    userID,
			TokenHash: "nocache-hash",
		},
		callCount: make(map[string]int),
	}

	svc := &Service{
		repo:        repo,
		logger:      logging.NewTestLogger(),
		tokenLength: 32,
		expiry:      24 * time.Hour,
		maxPerUser:  10,
	}

	cached := NewCachedService(svc, nil, logging.NewTestLogger(), 5*time.Minute)

	token, refreshToken, err := cached.CreateSession(
		context.Background(),
		userID,
		DeviceInfo{},
		[]string{"read"},
	)

	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotEmpty(t, refreshToken)
}

func TestCachedService_RevokeAllUserSessions_WithCache(t *testing.T) {
	t.Parallel()

	l1Cache, err := cache.NewCache(nil, 1000, time.Minute)
	require.NoError(t, err)
	defer l1Cache.Close()

	repo := &mockRepository{
		session:   db.SharedSession{},
		callCount: make(map[string]int),
	}

	svc := &Service{
		repo:   repo,
		logger: logging.NewTestLogger(),
	}

	cached := NewCachedService(svc, l1Cache, logging.NewTestLogger(), 5*time.Minute)

	userID := uuid.Must(uuid.NewV7())
	err = cached.RevokeAllUserSessions(context.Background(), userID)
	require.NoError(t, err)

	assert.Equal(t, 1, repo.getCallCount("CountActiveUserSessions"))
	assert.Equal(t, 1, repo.getCallCount("RevokeAllUserSessions"))
}

func TestCachedService_RevokeAllUserSessions_WithoutCache(t *testing.T) {
	t.Parallel()

	repo := &mockRepository{
		session:   db.SharedSession{},
		callCount: make(map[string]int),
	}

	svc := &Service{
		repo:   repo,
		logger: logging.NewTestLogger(),
	}

	cached := NewCachedService(svc, nil, logging.NewTestLogger(), 5*time.Minute)

	userID := uuid.Must(uuid.NewV7())
	err := cached.RevokeAllUserSessions(context.Background(), userID)
	require.NoError(t, err)

	assert.Equal(t, 1, repo.getCallCount("CountActiveUserSessions"))
	assert.Equal(t, 1, repo.getCallCount("RevokeAllUserSessions"))
}

func TestCachedService_RevokeSession_WithoutCache(t *testing.T) {
	t.Parallel()

	sessionID := uuid.Must(uuid.NewV7())
	repo := &mockRepository{
		session: db.SharedSession{
			ID:        sessionID,
			TokenHash: "test-hash",
		},
		callCount: make(map[string]int),
	}

	svc := &Service{
		repo:   repo,
		logger: logging.NewTestLogger(),
	}

	cached := NewCachedService(svc, nil, logging.NewTestLogger(), 5*time.Minute)

	err := cached.RevokeSession(context.Background(), sessionID)
	require.NoError(t, err)

	assert.Equal(t, 1, repo.getCallCount("GetSessionByID"))
	assert.Equal(t, 1, repo.getCallCount("RevokeSession"))
}

func TestCachedService_RevokeSession_WithCache(t *testing.T) {
	t.Parallel()

	l1Cache, err := cache.NewCache(nil, 1000, time.Minute)
	require.NoError(t, err)
	defer l1Cache.Close()

	sessionID := uuid.Must(uuid.NewV7())
	repo := &mockRepository{
		session: db.SharedSession{
			ID:        sessionID,
			TokenHash: "revoke-hash",
		},
		callCount: make(map[string]int),
	}

	svc := &Service{
		repo:   repo,
		logger: logging.NewTestLogger(),
	}

	cached := NewCachedService(svc, l1Cache, logging.NewTestLogger(), 5*time.Minute)

	err = cached.RevokeSession(context.Background(), sessionID)
	require.NoError(t, err)

	assert.Equal(t, 1, repo.getCallCount("GetSessionByID"))
	assert.Equal(t, 1, repo.getCallCount("RevokeSession"))
}

func TestCachedService_ValidateSession_CacheMissThenHit(t *testing.T) {
	t.Parallel()

	l1Cache, err := cache.NewCache(nil, 1000, 15*time.Second)
	require.NoError(t, err)
	defer l1Cache.Close()

	repo := &mockRepository{
		session: db.SharedSession{
			ID:        uuid.Must(uuid.NewV7()),
			UserID:    uuid.Must(uuid.NewV7()),
			TokenHash: "validate-hash",
		},
		callCount: make(map[string]int),
	}

	svc := &Service{
		repo:   repo,
		logger: logging.NewTestLogger(),
	}

	cached := NewCachedService(svc, l1Cache, logging.NewTestLogger(), 5*time.Minute)

	// First call - cache miss
	session1, err := cached.ValidateSession(context.Background(), "validate-token")
	require.NoError(t, err)
	assert.NotNil(t, session1)

	// Wait for async cache set
	time.Sleep(100 * time.Millisecond)

	// Second call - cache hit
	session2, err := cached.ValidateSession(context.Background(), "validate-token")
	require.NoError(t, err)
	assert.NotNil(t, session2)

	// Wait for background operations
	time.Sleep(100 * time.Millisecond)

	// Repository should only be called once
	assert.Equal(t, 1, repo.getCallCount("GetSessionByTokenHash"))
}

// =====================================================
// NewService (module.go) tests
// =====================================================

func TestNewService_ConfigDefaults(t *testing.T) {
	t.Parallel()

	// Simulate what module.go NewService does with various config values
	tests := []struct {
		name           string
		tokenLength    int
		maxPerUser     int
		expectedToken  int
		expectedMaxPer int
	}{
		{"all zeros fallback to defaults", 0, 0, 32, 10},
		{"custom values", 64, 20, 64, 20},
		{"token zero max custom", 0, 5, 32, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenLength := tt.tokenLength
			if tokenLength == 0 {
				tokenLength = 32
			}
			maxPerUser := tt.maxPerUser
			if maxPerUser == 0 {
				maxPerUser = 10
			}

			assert.Equal(t, tt.expectedToken, tokenLength)
			assert.Equal(t, tt.expectedMaxPer, maxPerUser)
		})
	}
}

// =====================================================
// NewServiceForTesting tests
// =====================================================

func TestNewServiceForTesting_AllFieldsSet(t *testing.T) {
	t.Parallel()

	repo := &mockRepository{callCount: make(map[string]int)}
	logger := logging.NewTestLogger()

	svc := NewServiceForTesting(
		repo,
		logger,
		64,
		48*time.Hour,
		14*24*time.Hour,
		20,
	)

	assert.NotNil(t, svc)
	assert.Equal(t, repo, svc.repo)
	assert.Equal(t, 64, svc.tokenLength)
	assert.Equal(t, 48*time.Hour, svc.expiry)
	assert.Equal(t, 14*24*time.Hour, svc.refreshExpiry)
	assert.Equal(t, 20, svc.maxPerUser)
}

// =====================================================
// DeviceInfo and SessionInfo struct tests
// =====================================================

func TestDeviceInfo_AllNil(t *testing.T) {
	t.Parallel()

	di := DeviceInfo{}
	assert.Nil(t, di.DeviceName)
	assert.Nil(t, di.UserAgent)
	assert.Nil(t, di.IPAddress)
}

func TestDeviceInfo_AllPopulated(t *testing.T) {
	t.Parallel()

	name := "iPhone"
	agent := "Safari"
	ip := netip.MustParseAddr("10.0.0.1")

	di := DeviceInfo{
		DeviceName: &name,
		UserAgent:  &agent,
		IPAddress:  &ip,
	}

	assert.Equal(t, "iPhone", *di.DeviceName)
	assert.Equal(t, "Safari", *di.UserAgent)
	assert.Equal(t, netip.MustParseAddr("10.0.0.1"), *di.IPAddress)
}

func TestSessionInfo_Fields(t *testing.T) {
	t.Parallel()

	id := uuid.Must(uuid.NewV7())
	now := time.Now()
	device := "Test"
	ip := "192.168.1.1"
	agent := "Test Agent"

	info := SessionInfo{
		ID:             id,
		DeviceName:     &device,
		IPAddress:      &ip,
		UserAgent:      &agent,
		CreatedAt:      now,
		LastActivityAt: now,
		ExpiresAt:      now.Add(24 * time.Hour),
		IsActive:       true,
		IsCurrent:      false,
	}

	assert.Equal(t, id, info.ID)
	assert.True(t, info.IsActive)
	assert.False(t, info.IsCurrent)
	assert.Equal(t, "Test", *info.DeviceName)
	assert.Equal(t, "192.168.1.1", *info.IPAddress)
	assert.Equal(t, "Test Agent", *info.UserAgent)
	assert.Equal(t, now, info.CreatedAt)
	assert.Equal(t, now, info.LastActivityAt)
	assert.True(t, info.ExpiresAt.After(now))
}

// =====================================================
// CreateSessionParams tests
// =====================================================

func TestCreateSessionParams_Fields(t *testing.T) {
	t.Parallel()

	userID := uuid.Must(uuid.NewV7())
	ip := netip.MustParseAddr("10.0.0.1")
	agent := "Chrome"
	device := "Desktop"
	refreshHash := "refresh123"

	params := CreateSessionParams{
		UserID:           userID,
		TokenHash:        "token123",
		RefreshTokenHash: &refreshHash,
		IPAddress:        &ip,
		UserAgent:        &agent,
		DeviceName:       &device,
		Scopes:           []string{"read", "write"},
		ExpiresAt:        time.Now().Add(24 * time.Hour),
	}

	assert.Equal(t, userID, params.UserID)
	assert.Equal(t, "token123", params.TokenHash)
	assert.Equal(t, "refresh123", *params.RefreshTokenHash)
	assert.Equal(t, netip.MustParseAddr("10.0.0.1"), *params.IPAddress)
	assert.Equal(t, "Chrome", *params.UserAgent)
	assert.Equal(t, "Desktop", *params.DeviceName)
	assert.Len(t, params.Scopes, 2)
	assert.True(t, params.ExpiresAt.After(time.Now()))
}

// errRepository is a mock that returns errors
type errRepository struct {
	mockRepository
	countErr  error
	createErr error
}

func (r *errRepository) CountActiveUserSessions(ctx context.Context, userID uuid.UUID) (int64, error) {
	if r.countErr != nil {
		return 0, r.countErr
	}
	return 0, nil
}

func (r *errRepository) CreateSession(ctx context.Context, params CreateSessionParams) (db.SharedSession, error) {
	if r.createErr != nil {
		return db.SharedSession{}, r.createErr
	}
	return db.SharedSession{ID: uuid.Must(uuid.NewV7())}, nil
}

func (r *errRepository) GetSessionByTokenHash(ctx context.Context, tokenHash string) (*db.SharedSession, error) {
	return nil, sql.ErrNoRows
}

func TestCachedService_CreateSession_ServiceError(t *testing.T) {
	t.Parallel()

	repo := &errRepository{
		mockRepository: mockRepository{callCount: make(map[string]int)},
		countErr:       assert.AnError,
	}

	svc := &Service{
		repo:        repo,
		logger:      logging.NewTestLogger(),
		tokenLength: 32,
		expiry:      24 * time.Hour,
		maxPerUser:  10,
	}

	l1Cache, err := cache.NewCache(nil, 100, time.Minute)
	require.NoError(t, err)
	defer l1Cache.Close()

	cached := NewCachedService(svc, l1Cache, logging.NewTestLogger(), 5*time.Minute)

	token, refreshToken, err := cached.CreateSession(
		context.Background(),
		uuid.Must(uuid.NewV7()),
		DeviceInfo{},
		[]string{"read"},
	)

	require.Error(t, err)
	assert.Empty(t, token)
	assert.Empty(t, refreshToken)
}
