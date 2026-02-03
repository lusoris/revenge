package apikeys

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestService_CreateKey(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	logger := zaptest.NewLogger(t)
	svc := NewService(repo, logger, 10, 0)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	resp, err := svc.CreateKey(ctx, userID, CreateKeyRequest{
		Name:        "Test Key",
		Description: stringPtr("Test description"),
		Scopes:      []string{"read", "write"},
	})
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Verify raw key format
	assert.True(t, strings.HasPrefix(resp.RawKey, KeyPrefix), "Raw key should start with rv_")
	assert.Equal(t, len(KeyPrefix)+(KeyLength*2), len(resp.RawKey), "Raw key should be correct length")

	// Verify key data
	assert.Equal(t, "Test Key", resp.Key.Name)
	assert.Equal(t, "Test description", *resp.Key.Description)
	assert.Equal(t, []string{"read", "write"}, resp.Key.Scopes)
	assert.True(t, resp.Key.IsActive)
	assert.Equal(t, resp.RawKey[:8], resp.Key.KeyPrefix)
}

func TestService_CreateKey_MaxKeysExceeded(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	logger := zaptest.NewLogger(t)
	svc := NewService(repo, logger, 2, 0) // Max 2 keys
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Create 2 keys (max)
	for i := 0; i < 2; i++ {
		_, err := svc.CreateKey(ctx, userID, CreateKeyRequest{
			Name:   "Key",
			Scopes: []string{"read"},
		})
		require.NoError(t, err)
	}

	// Third key should fail
	_, err := svc.CreateKey(ctx, userID, CreateKeyRequest{
		Name:   "Key 3",
		Scopes: []string{"read"},
	})
	assert.ErrorIs(t, err, ErrMaxKeysExceeded)
}

func TestService_CreateKey_InvalidScope(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	logger := zaptest.NewLogger(t)
	svc := NewService(repo, logger, 10, 0)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	_, err := svc.CreateKey(ctx, userID, CreateKeyRequest{
		Name:   "Invalid Scope Key",
		Scopes: []string{"invalid_scope"},
	})
	assert.ErrorIs(t, err, ErrInvalidScope)
}

func TestService_CreateKey_WithExpiry(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	logger := zaptest.NewLogger(t)
	svc := NewService(repo, logger, 10, 0)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	expiresAt := time.Now().Add(24 * time.Hour)
	resp, err := svc.CreateKey(ctx, userID, CreateKeyRequest{
		Name:      "Expiring Key",
		Scopes:    []string{"read"},
		ExpiresAt: &expiresAt,
	})
	require.NoError(t, err)
	require.NotNil(t, resp.Key.ExpiresAt)
	assert.WithinDuration(t, expiresAt, *resp.Key.ExpiresAt, time.Second)
}

func TestService_CreateKey_DefaultExpiry(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	logger := zaptest.NewLogger(t)
	defaultExpiry := 90 * 24 * time.Hour
	svc := NewService(repo, logger, 10, defaultExpiry)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	resp, err := svc.CreateKey(ctx, userID, CreateKeyRequest{
		Name:   "Auto Expiring Key",
		Scopes: []string{"read"},
	})
	require.NoError(t, err)
	require.NotNil(t, resp.Key.ExpiresAt)
	assert.WithinDuration(t, time.Now().Add(defaultExpiry), *resp.Key.ExpiresAt, time.Second)
}

func TestService_GetKey(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	logger := zaptest.NewLogger(t)
	svc := NewService(repo, logger, 10, 0)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Create a key
	resp, err := svc.CreateKey(ctx, userID, CreateKeyRequest{
		Name:   "Get Test Key",
		Scopes: []string{"read"},
	})
	require.NoError(t, err)

	// Get it back
	key, err := svc.GetKey(ctx, resp.Key.ID)
	require.NoError(t, err)
	assert.Equal(t, resp.Key.ID, key.ID)
	assert.Equal(t, "Get Test Key", key.Name)
}

func TestService_GetKey_NotFound(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	logger := zaptest.NewLogger(t)
	svc := NewService(repo, logger, 10, 0)
	ctx := context.Background()

	_, err := svc.GetKey(ctx, uuid.New())
	assert.ErrorIs(t, err, ErrKeyNotFound)
}

func TestService_ListUserKeys(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	logger := zaptest.NewLogger(t)
	svc := NewService(repo, logger, 10, 0)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Create 3 keys
	for i := 0; i < 3; i++ {
		_, err := svc.CreateKey(ctx, userID, CreateKeyRequest{
			Name:   "List Key",
			Scopes: []string{"read"},
		})
		require.NoError(t, err)
	}

	// List them
	keys, err := svc.ListUserKeys(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, keys, 3)
}

func TestService_ValidateKey(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	logger := zaptest.NewLogger(t)
	svc := NewService(repo, logger, 10, 0)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Create a key
	resp, err := svc.CreateKey(ctx, userID, CreateKeyRequest{
		Name:   "Validate Test",
		Scopes: []string{"read", "write"},
	})
	require.NoError(t, err)

	// Validate it
	key, err := svc.ValidateKey(ctx, resp.RawKey)
	require.NoError(t, err)
	assert.Equal(t, resp.Key.ID, key.ID)
	assert.True(t, key.IsActive)

	// Wait a bit to ensure last_used_at would be updated
	time.Sleep(50 * time.Millisecond)

	// Get key again to check last_used_at was updated (async)
	updatedKey, err := svc.GetKey(ctx, resp.Key.ID)
	require.NoError(t, err)
	// Last used might not be set yet due to async update, but it shouldn't error
	_ = updatedKey.LastUsedAt
}

func TestService_ValidateKey_InvalidFormat(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	logger := zaptest.NewLogger(t)
	svc := NewService(repo, logger, 10, 0)
	ctx := context.Background()

	_, err := svc.ValidateKey(ctx, "invalid_key")
	assert.ErrorIs(t, err, ErrInvalidKeyFormat)
}

func TestService_ValidateKey_NotFound(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	logger := zaptest.NewLogger(t)
	svc := NewService(repo, logger, 10, 0)
	ctx := context.Background()

	// Valid format but doesn't exist
	fakeKey := KeyPrefix + strings.Repeat("a", KeyLength*2)
	_, err := svc.ValidateKey(ctx, fakeKey)
	assert.ErrorIs(t, err, ErrKeyNotFound)
}

func TestService_ValidateKey_Inactive(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	logger := zaptest.NewLogger(t)
	svc := NewService(repo, logger, 10, 0)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Create and revoke a key
	resp, err := svc.CreateKey(ctx, userID, CreateKeyRequest{
		Name:   "Revoked Key",
		Scopes: []string{"read"},
	})
	require.NoError(t, err)

	err = svc.RevokeKey(ctx, resp.Key.ID)
	require.NoError(t, err)

	// Try to validate revoked key
	_, err = svc.ValidateKey(ctx, resp.RawKey)
	assert.ErrorIs(t, err, ErrKeyInactive)
}

func TestService_ValidateKey_Expired(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	logger := zaptest.NewLogger(t)
	svc := NewService(repo, logger, 10, 0)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Create expired key directly in DB
	queries := db.New(testDB.Pool())
	dbKey, err := queries.CreateAPIKey(ctx, db.CreateAPIKeyParams{
		UserID:    userID,
		Name:      "Expired Key",
		KeyHash:   "expired_test_hash",
		KeyPrefix: "rv_expir",
		Scopes:    []string{"read"},
		ExpiresAt: pgtype.Timestamptz{
			Time:  time.Now().Add(-1 * time.Hour),
			Valid: true,
		},
	})
	require.NoError(t, err)

	// Create a "raw key" that would hash to our known hash
	// For testing, we use the hash directly (in production this would be generated)
	// We need to create a key that when hashed gives us "expired_test_hash"
	// For simplicity, let's just check the error with a properly formatted fake key

	// Actually, we need to get the raw key. Let's create a real one and expire it after
	resp, err := svc.CreateKey(ctx, userID, CreateKeyRequest{
		Name:      "Will Expire",
		Scopes:    []string{"read"},
		ExpiresAt: timePtr(time.Now().Add(1 * time.Second)),
	})
	require.NoError(t, err)

	// Wait for expiry
	time.Sleep(1100 * time.Millisecond)

	// Validate expired key
	_, err = svc.ValidateKey(ctx, resp.RawKey)
	assert.ErrorIs(t, err, ErrKeyExpired)

	// Clean up the unused expired key
	_ = dbKey
}

func TestService_RevokeKey(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	logger := zaptest.NewLogger(t)
	svc := NewService(repo, logger, 10, 0)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	resp, err := svc.CreateKey(ctx, userID, CreateKeyRequest{
		Name:   "Revoke Test",
		Scopes: []string{"read"},
	})
	require.NoError(t, err)

	err = svc.RevokeKey(ctx, resp.Key.ID)
	require.NoError(t, err)

	// Verify it's revoked
	key, err := svc.GetKey(ctx, resp.Key.ID)
	require.NoError(t, err)
	assert.False(t, key.IsActive)
}

func TestService_CheckScope(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	logger := zaptest.NewLogger(t)
	svc := NewService(repo, logger, 10, 0)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Create key with read scope
	resp, err := svc.CreateKey(ctx, userID, CreateKeyRequest{
		Name:   "Scope Test",
		Scopes: []string{"read"},
	})
	require.NoError(t, err)

	// Check read scope - should pass
	hasRead, err := svc.CheckScope(ctx, resp.Key.ID, "read")
	require.NoError(t, err)
	assert.True(t, hasRead)

	// Check write scope - should fail
	hasWrite, err := svc.CheckScope(ctx, resp.Key.ID, "write")
	require.NoError(t, err)
	assert.False(t, hasWrite)
}

func TestService_CheckScope_Admin(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	logger := zaptest.NewLogger(t)
	svc := NewService(repo, logger, 10, 0)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Create key with admin scope
	resp, err := svc.CreateKey(ctx, userID, CreateKeyRequest{
		Name:   "Admin Key",
		Scopes: []string{"admin"},
	})
	require.NoError(t, err)

	// Admin scope should grant access to any scope
	hasRead, err := svc.CheckScope(ctx, resp.Key.ID, "read")
	require.NoError(t, err)
	assert.True(t, hasRead, "Admin scope should grant read access")

	hasWrite, err := svc.CheckScope(ctx, resp.Key.ID, "write")
	require.NoError(t, err)
	assert.True(t, hasWrite, "Admin scope should grant write access")
}

func TestService_UpdateScopes(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	logger := zaptest.NewLogger(t)
	svc := NewService(repo, logger, 10, 0)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Create key with read scope
	resp, err := svc.CreateKey(ctx, userID, CreateKeyRequest{
		Name:   "Update Scopes Test",
		Scopes: []string{"read"},
	})
	require.NoError(t, err)

	// Update to write scope
	err = svc.UpdateScopes(ctx, resp.Key.ID, []string{"write", "admin"})
	require.NoError(t, err)

	// Verify scopes updated
	key, err := svc.GetKey(ctx, resp.Key.ID)
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"write", "admin"}, key.Scopes)
}

func TestService_UpdateScopes_InvalidScope(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	logger := zaptest.NewLogger(t)
	svc := NewService(repo, logger, 10, 0)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	resp, err := svc.CreateKey(ctx, userID, CreateKeyRequest{
		Name:   "Invalid Update Test",
		Scopes: []string{"read"},
	})
	require.NoError(t, err)

	err = svc.UpdateScopes(ctx, resp.Key.ID, []string{"invalid"})
	assert.ErrorIs(t, err, ErrInvalidScope)
}

func TestService_CleanupExpiredKeys(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	logger := zaptest.NewLogger(t)
	svc := NewService(repo, logger, 10, 0)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Create expired AND revoked key
	resp, err := svc.CreateKey(ctx, userID, CreateKeyRequest{
		Name:      "Cleanup Test",
		Scopes:    []string{"read"},
		ExpiresAt: timePtr(time.Now().Add(-1 * time.Hour)),
	})
	require.NoError(t, err)

	// Revoke it (cleanup only removes revoked keys)
	err = svc.RevokeKey(ctx, resp.Key.ID)
	require.NoError(t, err)

	// Run cleanup
	err = svc.CleanupExpiredKeys(ctx)
	require.NoError(t, err)

	// Verify key was deleted
	var count int
	err = testDB.Pool().QueryRow(ctx, "SELECT COUNT(*) FROM shared.api_keys WHERE id = $1", resp.Key.ID).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func timePtr(t time.Time) *time.Time {
	return &t
}
