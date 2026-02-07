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

// ============================================================================
// Full Lifecycle: Create -> Validate -> Revoke -> Validate (fails)
// ============================================================================

func TestServiceIntegration_CreateValidateRevokeValidate(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	logger := zaptest.NewLogger(t)
	svc := NewService(repo, logger, 10, 0)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Create key
	resp, err := svc.CreateKey(ctx, userID, CreateKeyRequest{
		Name:   "Lifecycle Key",
		Scopes: []string{"read", "write"},
	})
	require.NoError(t, err)
	rawKey := resp.RawKey
	keyID := resp.Key.ID

	// Validate key
	key, err := svc.ValidateKey(ctx, rawKey)
	require.NoError(t, err)
	assert.Equal(t, keyID, key.ID)
	assert.True(t, key.IsActive)
	assert.ElementsMatch(t, []string{"read", "write"}, key.Scopes)

	// Revoke key
	err = svc.RevokeKey(ctx, keyID)
	require.NoError(t, err)

	// Validate again should fail
	_, err = svc.ValidateKey(ctx, rawKey)
	assert.ErrorIs(t, err, ErrKeyInactive)
}

// ============================================================================
// ValidateKey with expired key
// ============================================================================

func TestServiceIntegration_ValidateKey_ExpiredKey(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	logger := zaptest.NewLogger(t)
	svc := NewService(repo, logger, 10, 0)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Create key with very short expiry
	resp, err := svc.CreateKey(ctx, userID, CreateKeyRequest{
		Name:      "Expiring Key",
		Scopes:    []string{"read"},
		ExpiresAt: timePtr(time.Now().Add(500 * time.Millisecond)),
	})
	require.NoError(t, err)

	// Should validate before expiry
	key, err := svc.ValidateKey(ctx, resp.RawKey)
	require.NoError(t, err)
	assert.Equal(t, resp.Key.ID, key.ID)

	// Wait for expiry
	time.Sleep(600 * time.Millisecond)

	// Should fail after expiry
	_, err = svc.ValidateKey(ctx, resp.RawKey)
	assert.ErrorIs(t, err, ErrKeyExpired)
}

// ============================================================================
// UpdateScopes -> ValidateKey returns new scopes
// ============================================================================

func TestServiceIntegration_UpdateScopes_ValidateReturnsNewScopes(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	logger := zaptest.NewLogger(t)
	svc := NewService(repo, logger, 10, 0)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Create key with read scope
	resp, err := svc.CreateKey(ctx, userID, CreateKeyRequest{
		Name:   "Scope Update Key",
		Scopes: []string{"read"},
	})
	require.NoError(t, err)

	// Update scopes to read+write
	err = svc.UpdateScopes(ctx, resp.Key.ID, []string{"read", "write"})
	require.NoError(t, err)

	// Validate should return new scopes
	key, err := svc.ValidateKey(ctx, resp.RawKey)
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"read", "write"}, key.Scopes)

	// CheckScope should reflect new scopes
	hasWrite, err := svc.CheckScope(ctx, resp.Key.ID, "write")
	require.NoError(t, err)
	assert.True(t, hasWrite)

	// Wait for async last_used_at update to settle
	time.Sleep(50 * time.Millisecond)
}

// ============================================================================
// DeleteExpiredAPIKeys cleanup
// ============================================================================

func TestServiceIntegration_CleanupExpiredKeys_OnlyDeletesRevokedExpired(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	logger := zaptest.NewLogger(t)
	svc := NewService(repo, logger, 10, 0)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Create an active key with no expiry (should NOT be deleted)
	activeResp, err := svc.CreateKey(ctx, userID, CreateKeyRequest{
		Name:   "Active Key",
		Scopes: []string{"read"},
	})
	require.NoError(t, err)

	// Create an expired but still active key (should NOT be deleted - cleanup requires is_active=false)
	expiredActiveResp, err := svc.CreateKey(ctx, userID, CreateKeyRequest{
		Name:      "Expired Active Key",
		Scopes:    []string{"read"},
		ExpiresAt: timePtr(time.Now().Add(-1 * time.Hour)),
	})
	require.NoError(t, err)

	// Create a revoked key WITHOUT expiry (should NOT be deleted - cleanup requires expires_at)
	revokedNoExpiryResp, err := svc.CreateKey(ctx, userID, CreateKeyRequest{
		Name:   "Revoked No Expiry Key",
		Scopes: []string{"read"},
	})
	require.NoError(t, err)
	err = svc.RevokeKey(ctx, revokedNoExpiryResp.Key.ID)
	require.NoError(t, err)

	// Create a revoked AND expired key (SHOULD be deleted - both conditions met)
	revokedExpiredResp, err := svc.CreateKey(ctx, userID, CreateKeyRequest{
		Name:      "Revoked Expired Key",
		Scopes:    []string{"read"},
		ExpiresAt: timePtr(time.Now().Add(-1 * time.Hour)),
	})
	require.NoError(t, err)
	err = svc.RevokeKey(ctx, revokedExpiredResp.Key.ID)
	require.NoError(t, err)

	// Run cleanup
	err = svc.CleanupExpiredKeys(ctx)
	require.NoError(t, err)

	// Active key should still exist
	key, err := svc.GetKey(ctx, activeResp.Key.ID)
	require.NoError(t, err)
	assert.True(t, key.IsActive)

	// Expired but active key should still exist
	key, err = svc.GetKey(ctx, expiredActiveResp.Key.ID)
	require.NoError(t, err)
	assert.Equal(t, "Expired Active Key", key.Name)

	// Revoked but no expiry key should still exist
	key, err = svc.GetKey(ctx, revokedNoExpiryResp.Key.ID)
	require.NoError(t, err)
	assert.Equal(t, "Revoked No Expiry Key", key.Name)

	// Revoked AND expired key should be gone
	_, err = svc.GetKey(ctx, revokedExpiredResp.Key.ID)
	assert.ErrorIs(t, err, ErrKeyNotFound)
}

// ============================================================================
// CountUserAPIKeys accuracy after create/delete cycles
// ============================================================================

func TestServiceIntegration_CountUserAPIKeys_AfterCreateDeleteCycles(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	logger := zaptest.NewLogger(t)
	svc := NewService(repo, logger, 10, 0)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Count should be 0 initially
	count, err := repo.CountUserAPIKeys(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)

	// Create 3 keys
	var keyIDs []uuid.UUID
	for i := 0; i < 3; i++ {
		resp, err := svc.CreateKey(ctx, userID, CreateKeyRequest{
			Name:   "Count Key",
			Scopes: []string{"read"},
		})
		require.NoError(t, err)
		keyIDs = append(keyIDs, resp.Key.ID)
	}

	// Count should be 3
	count, err = repo.CountUserAPIKeys(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)

	// Delete one key directly through repo
	err = repo.DeleteAPIKey(ctx, keyIDs[0])
	require.NoError(t, err)

	// Count should be 2
	count, err = repo.CountUserAPIKeys(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, int64(2), count)

	// Revoking DOES change count (CountUserAPIKeys only counts active keys)
	err = svc.RevokeKey(ctx, keyIDs[1])
	require.NoError(t, err)

	count, err = repo.CountUserAPIKeys(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, int64(1), count) // Only 1 active key remains
}

// ============================================================================
// Error cases
// ============================================================================

func TestServiceIntegration_ValidateKey_InvalidKeyHash(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	logger := zaptest.NewLogger(t)
	svc := NewService(repo, logger, 10, 0)
	ctx := context.Background()

	// Valid format but key hash not in database
	fakeKey := KeyPrefix + strings.Repeat("ab", KeyLength)
	_, err := svc.ValidateKey(ctx, fakeKey)
	assert.ErrorIs(t, err, ErrKeyNotFound)
}

func TestServiceIntegration_GetKey_NonExistentKey(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	logger := zaptest.NewLogger(t)
	svc := NewService(repo, logger, 10, 0)
	ctx := context.Background()

	_, err := svc.GetKey(ctx, uuid.Must(uuid.NewV7()))
	assert.ErrorIs(t, err, ErrKeyNotFound)
}

func TestServiceIntegration_CheckScope_NonExistentKey(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	logger := zaptest.NewLogger(t)
	svc := NewService(repo, logger, 10, 0)
	ctx := context.Background()

	has, err := svc.CheckScope(ctx, uuid.Must(uuid.NewV7()), "read")
	assert.False(t, has)
	assert.ErrorIs(t, err, ErrKeyNotFound)
}

func TestServiceIntegration_ValidateKey_InvalidFormat(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	logger := zaptest.NewLogger(t)
	svc := NewService(repo, logger, 10, 0)
	ctx := context.Background()

	tests := []struct {
		name string
		key  string
	}{
		{"empty", ""},
		{"too short", "rv_"},
		{"wrong prefix", "xx_" + strings.Repeat("a", KeyLength*2)},
		{"wrong length", KeyPrefix + "short"},
		{"only prefix", KeyPrefix},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.ValidateKey(ctx, tt.key)
			assert.ErrorIs(t, err, ErrInvalidKeyFormat)
		})
	}
}

// ============================================================================
// NewService constructor
// ============================================================================

func TestServiceIntegration_NewService_DefaultMaxKeys(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	logger := zaptest.NewLogger(t)

	// maxKeysPerUser <= 0 should use default
	svc := NewService(repo, logger, 0, 0)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Should be able to create at least 1 key (default is 10)
	resp, err := svc.CreateKey(ctx, userID, CreateKeyRequest{
		Name:   "Default Max Key",
		Scopes: []string{"read"},
	})
	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestServiceIntegration_NewService_NegativeMaxKeys(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	logger := zaptest.NewLogger(t)

	// Negative should use default
	svc := NewService(repo, logger, -1, 0)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	resp, err := svc.CreateKey(ctx, userID, CreateKeyRequest{
		Name:   "Negative Max Key",
		Scopes: []string{"read"},
	})
	require.NoError(t, err)
	assert.NotNil(t, resp)
}

// ============================================================================
// ListUserKeys - empty + with keys
// ============================================================================

func TestServiceIntegration_ListUserKeys_Empty(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	logger := zaptest.NewLogger(t)
	svc := NewService(repo, logger, 10, 0)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	keys, err := svc.ListUserKeys(ctx, userID)
	require.NoError(t, err)
	assert.Empty(t, keys)
}

// ============================================================================
// CreateKey with default expiry from service config
// ============================================================================

func TestServiceIntegration_CreateKey_ServiceDefaultExpiry(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	logger := zaptest.NewLogger(t)
	defaultExpiry := 30 * 24 * time.Hour
	svc := NewService(repo, logger, 10, defaultExpiry)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Create without explicit expiry - should use service default
	resp, err := svc.CreateKey(ctx, userID, CreateKeyRequest{
		Name:   "Default Expiry Key",
		Scopes: []string{"read"},
	})
	require.NoError(t, err)
	require.NotNil(t, resp.Key.ExpiresAt)
	assert.WithinDuration(t, time.Now().Add(defaultExpiry), *resp.Key.ExpiresAt, 2*time.Second)
}

func TestServiceIntegration_CreateKey_ExplicitExpiryOverridesDefault(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	logger := zaptest.NewLogger(t)
	defaultExpiry := 30 * 24 * time.Hour
	svc := NewService(repo, logger, 10, defaultExpiry)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Create with explicit expiry - should use explicit, not default
	explicitExpiry := time.Now().Add(7 * 24 * time.Hour)
	resp, err := svc.CreateKey(ctx, userID, CreateKeyRequest{
		Name:      "Explicit Expiry Key",
		Scopes:    []string{"read"},
		ExpiresAt: &explicitExpiry,
	})
	require.NoError(t, err)
	require.NotNil(t, resp.Key.ExpiresAt)
	assert.WithinDuration(t, explicitExpiry, *resp.Key.ExpiresAt, time.Second)
}

// ============================================================================
// Key format validation
// ============================================================================

func TestServiceIntegration_KeyFormat(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	logger := zaptest.NewLogger(t)
	svc := NewService(repo, logger, 10, 0)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	resp, err := svc.CreateKey(ctx, userID, CreateKeyRequest{
		Name:   "Format Key",
		Scopes: []string{"read"},
	})
	require.NoError(t, err)

	// Verify key format
	assert.True(t, strings.HasPrefix(resp.RawKey, KeyPrefix))
	expectedLen := len(KeyPrefix) + (KeyLength * 2)
	assert.Equal(t, expectedLen, len(resp.RawKey))

	// Prefix stored should be first 8 chars
	assert.Equal(t, resp.RawKey[:8], resp.Key.KeyPrefix)
}

// ============================================================================
// Validate scopes
// ============================================================================

func TestServiceIntegration_CreateKey_MultipleInvalidScopes(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	logger := zaptest.NewLogger(t)
	svc := NewService(repo, logger, 10, 0)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Mix of valid and invalid scopes
	_, err := svc.CreateKey(ctx, userID, CreateKeyRequest{
		Name:   "Mixed Scopes Key",
		Scopes: []string{"read", "invalid_scope"},
	})
	assert.ErrorIs(t, err, ErrInvalidScope)
}

func TestServiceIntegration_UpdateScopes_Invalid(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	logger := zaptest.NewLogger(t)
	svc := NewService(repo, logger, 10, 0)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	resp, err := svc.CreateKey(ctx, userID, CreateKeyRequest{
		Name:   "Update Invalid Scopes Key",
		Scopes: []string{"read"},
	})
	require.NoError(t, err)

	err = svc.UpdateScopes(ctx, resp.Key.ID, []string{"nonexistent"})
	assert.ErrorIs(t, err, ErrInvalidScope)

	// Original scopes should remain
	key, err := svc.GetKey(ctx, resp.Key.ID)
	require.NoError(t, err)
	assert.Equal(t, []string{"read"}, key.Scopes)
}

// ============================================================================
// ValidateKey sets LastUsedAt
// ============================================================================

func TestServiceIntegration_ValidateKey_SetsLastUsedAt(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	logger := zaptest.NewLogger(t)
	svc := NewService(repo, logger, 10, 0)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	resp, err := svc.CreateKey(ctx, userID, CreateKeyRequest{
		Name:   "Last Used Key",
		Scopes: []string{"read"},
	})
	require.NoError(t, err)

	// Initially, last_used_at should be nil
	key, err := svc.GetKey(ctx, resp.Key.ID)
	require.NoError(t, err)
	assert.Nil(t, key.LastUsedAt)

	// Validate the key
	_, err = svc.ValidateKey(ctx, resp.RawKey)
	require.NoError(t, err)

	// Wait for the async last_used_at update
	time.Sleep(200 * time.Millisecond)

	// Now last_used_at should be set
	key, err = svc.GetKey(ctx, resp.Key.ID)
	require.NoError(t, err)
	assert.NotNil(t, key.LastUsedAt, "LastUsedAt should be set after ValidateKey")
}

// ============================================================================
// dbKeyToAPIKey conversion (ExpiresAt, LastUsedAt)
// ============================================================================

func TestServiceIntegration_dbKeyToAPIKey_Fields(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	logger := zaptest.NewLogger(t)
	svc := NewService(repo, logger, 10, 0)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Create key with all fields set
	expiresAt := time.Now().Add(24 * time.Hour)
	resp, err := svc.CreateKey(ctx, userID, CreateKeyRequest{
		Name:        "Full Fields Key",
		Description: stringPtr("A descriptive description"),
		Scopes:      []string{"read", "write", "admin"},
		ExpiresAt:   &expiresAt,
	})
	require.NoError(t, err)

	key, err := svc.GetKey(ctx, resp.Key.ID)
	require.NoError(t, err)

	assert.Equal(t, "Full Fields Key", key.Name)
	assert.Equal(t, "A descriptive description", *key.Description)
	assert.ElementsMatch(t, []string{"read", "write", "admin"}, key.Scopes)
	assert.True(t, key.IsActive)
	assert.NotNil(t, key.ExpiresAt)
	assert.WithinDuration(t, expiresAt, *key.ExpiresAt, time.Second)
	assert.Nil(t, key.LastUsedAt) // Not validated yet
	assert.False(t, key.CreatedAt.IsZero())
	assert.False(t, key.UpdatedAt.IsZero())
}

// ============================================================================
// CreateKey directly through DB with expired + revoked for cleanup
// ============================================================================

func TestServiceIntegration_CleanupExpiredKeys_DirectDB(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	logger := zaptest.NewLogger(t)
	svc := NewService(repo, logger, 10, 0)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Create expired key directly through DB (so we can bypass the service's hash generation)
	queries := db.New(testDB.Pool())
	dbKey, err := queries.CreateAPIKey(ctx, db.CreateAPIKeyParams{
		UserID:    userID,
		Name:      "Direct Expired Key",
		KeyHash:   "cleanup_test_hash_" + uuid.Must(uuid.NewV7()).String()[:8],
		KeyPrefix: "rv_clean",
		Scopes:    []string{"read"},
		ExpiresAt: pgtype.Timestamptz{
			Time:  time.Now().Add(-91 * 24 * time.Hour),
			Valid: true,
		},
	})
	require.NoError(t, err)

	// Revoke it
	err = repo.RevokeAPIKey(ctx, dbKey.ID)
	require.NoError(t, err)

	// Run cleanup through service
	err = svc.CleanupExpiredKeys(ctx)
	require.NoError(t, err)

	// Should be deleted
	_, err = svc.GetKey(ctx, dbKey.ID)
	assert.ErrorIs(t, err, ErrKeyNotFound)
}
