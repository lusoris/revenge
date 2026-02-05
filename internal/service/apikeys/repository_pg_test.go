package apikeys

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/testutil"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testutil.StopSharedPostgres()
	os.Exit(code)
}

func setupTestRepository(t *testing.T) (Repository, testutil.DB) {
	t.Helper()
	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())
	repo := NewRepositoryPg(queries)
	return repo, testDB
}

func createTestUser(t *testing.T, testDB testutil.DB) uuid.UUID {
	t.Helper()
	queries := db.New(testDB.Pool())
	user, err := queries.CreateUser(context.Background(), db.CreateUserParams{
		Username:     "apikey_user_" + uuid.New().String()[:8],
		Email:        "apikey_" + uuid.New().String()[:8] + "@example.com",
		PasswordHash: "hash",
	})
	require.NoError(t, err)
	return user.ID
}

func TestRepositoryPg_CreateAPIKey(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	key, err := repo.CreateAPIKey(ctx, db.CreateAPIKeyParams{
		UserID:      userID,
		Name:        "Test Key",
		Description: stringPtr("Test Description"),
		KeyHash:     "test_hash_123",
		KeyPrefix:   "rv_test",
		Scopes:      []string{"read", "write"},
	})
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, key.ID)
	assert.Equal(t, userID, key.UserID)
	assert.Equal(t, "Test Key", key.Name)
	assert.Equal(t, "test_hash_123", key.KeyHash)
	assert.True(t, key.IsActive)
}

func TestRepositoryPg_GetAPIKey(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	created, err := repo.CreateAPIKey(ctx, db.CreateAPIKeyParams{
		UserID:    userID,
		Name:      "Get Test",
		KeyHash:   "get_hash",
		KeyPrefix: "rv_get",
		Scopes:    []string{"read"},
	})
	require.NoError(t, err)

	retrieved, err := repo.GetAPIKey(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, retrieved.ID)
	assert.Equal(t, created.Name, retrieved.Name)
}

func TestRepositoryPg_GetAPIKeyByHash(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	created, err := repo.CreateAPIKey(ctx, db.CreateAPIKeyParams{
		UserID:    userID,
		Name:      "Hash Test",
		KeyHash:   "unique_hash_123",
		KeyPrefix: "rv_hash",
		Scopes:    []string{"read"},
	})
	require.NoError(t, err)

	retrieved, err := repo.GetAPIKeyByHash(ctx, "unique_hash_123")
	require.NoError(t, err)
	assert.Equal(t, created.ID, retrieved.ID)
}

func TestRepositoryPg_GetAPIKeyByPrefix(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	created, err := repo.CreateAPIKey(ctx, db.CreateAPIKeyParams{
		UserID:    userID,
		Name:      "Prefix Test",
		KeyHash:   "prefix_hash",
		KeyPrefix: "rv_unique_prefix",
		Scopes:    []string{"read"},
	})
	require.NoError(t, err)

	retrieved, err := repo.GetAPIKeyByPrefix(ctx, "rv_unique_prefix")
	require.NoError(t, err)
	assert.Equal(t, created.ID, retrieved.ID)
}

func TestRepositoryPg_ListUserAPIKeys(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Create multiple keys
	for i := 0; i < 3; i++ {
		_, err := repo.CreateAPIKey(ctx, db.CreateAPIKeyParams{
			UserID:    userID,
			Name:      "Key " + uuid.New().String()[:8],
			KeyHash:   "hash_" + uuid.New().String()[:8],
			KeyPrefix: "rv_" + uuid.New().String()[:8],
			Scopes:    []string{"read"},
		})
		require.NoError(t, err)
	}

	keys, err := repo.ListUserAPIKeys(ctx, userID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(keys), 3)
}

func TestRepositoryPg_ListActiveUserAPIKeys(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Create active key
	activeKey, err := repo.CreateAPIKey(ctx, db.CreateAPIKeyParams{
		UserID:    userID,
		Name:      "Active Key",
		KeyHash:   "active_hash",
		KeyPrefix: "rv_active",
		Scopes:    []string{"read"},
	})
	require.NoError(t, err)

	// Create and revoke a key
	revokedKey, err := repo.CreateAPIKey(ctx, db.CreateAPIKeyParams{
		UserID:    userID,
		Name:      "Revoked Key",
		KeyHash:   "revoked_hash",
		KeyPrefix: "rv_revoked",
		Scopes:    []string{"read"},
	})
	require.NoError(t, err)

	err = repo.RevokeAPIKey(ctx, revokedKey.ID)
	require.NoError(t, err)

	// List only active keys
	activeKeys, err := repo.ListActiveUserAPIKeys(ctx, userID)
	require.NoError(t, err)

	// Should contain active key
	found := false
	for _, key := range activeKeys {
		if key.ID == activeKey.ID {
			found = true
			break
		}
	}
	assert.True(t, found, "Active key should be in the list")

	// Should not contain revoked key
	for _, key := range activeKeys {
		assert.NotEqual(t, revokedKey.ID, key.ID, "Revoked key should not be in active list")
	}
}

func TestRepositoryPg_CountUserAPIKeys(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	initialCount, err := repo.CountUserAPIKeys(ctx, userID)
	require.NoError(t, err)

	// Create 3 keys (key_prefix max 16 chars)
	for i := 0; i < 3; i++ {
		_, err := repo.CreateAPIKey(ctx, db.CreateAPIKeyParams{
			UserID:    userID,
			Name:      "Count Key",
			KeyHash:   "count_hash_" + uuid.New().String()[:8],
			KeyPrefix: "rv_" + uuid.New().String()[:8],
			Scopes:    []string{"read"},
		})
		require.NoError(t, err)
	}

	count, err := repo.CountUserAPIKeys(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, initialCount+3, count)
}

func TestRepositoryPg_RevokeAPIKey(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	key, err := repo.CreateAPIKey(ctx, db.CreateAPIKeyParams{
		UserID:    userID,
		Name:      "Revoke Test",
		KeyHash:   "revoke_hash",
		KeyPrefix: "rv_revoke",
		Scopes:    []string{"read"},
	})
	require.NoError(t, err)

	err = repo.RevokeAPIKey(ctx, key.ID)
	require.NoError(t, err)

	// Check it's revoked by looking at DB directly
	var isActive bool
	err = testDB.Pool().QueryRow(ctx, "SELECT is_active FROM shared.api_keys WHERE id = $1", key.ID).Scan(&isActive)
	require.NoError(t, err)
	assert.False(t, isActive, "Key should be revoked (is_active=false)")
}

func TestRepositoryPg_UpdateAPIKeyLastUsed(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	key, err := repo.CreateAPIKey(ctx, db.CreateAPIKeyParams{
		UserID:    userID,
		Name:      "Last Used Test",
		KeyHash:   "lastused_hash",
		KeyPrefix: "rv_lastused",
		Scopes:    []string{"read"},
	})
	require.NoError(t, err)

	originalLastUsed := key.LastUsedAt

	time.Sleep(10 * time.Millisecond)
	err = repo.UpdateAPIKeyLastUsed(ctx, key.ID)
	require.NoError(t, err)

	updated, err := repo.GetAPIKey(ctx, key.ID)
	require.NoError(t, err)

	if originalLastUsed.Valid {
		assert.True(t, updated.LastUsedAt.Time.After(originalLastUsed.Time))
	} else {
		assert.True(t, updated.LastUsedAt.Valid)
	}
}

func TestRepositoryPg_UpdateAPIKeyScopes(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	key, err := repo.CreateAPIKey(ctx, db.CreateAPIKeyParams{
		UserID:    userID,
		Name:      "Scopes Test",
		KeyHash:   "scopes_hash",
		KeyPrefix: "rv_scopes",
		Scopes:    []string{"read"},
	})
	require.NoError(t, err)

	newScopes := []string{"read", "write", "delete"}
	err = repo.UpdateAPIKeyScopes(ctx, key.ID, newScopes)
	require.NoError(t, err)

	updated, err := repo.GetAPIKey(ctx, key.ID)
	require.NoError(t, err)
	assert.ElementsMatch(t, newScopes, updated.Scopes)
}

func TestRepositoryPg_DeleteAPIKey(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	key, err := repo.CreateAPIKey(ctx, db.CreateAPIKeyParams{
		UserID:    userID,
		Name:      "Delete Test",
		KeyHash:   "delete_hash",
		KeyPrefix: "rv_delete",
		Scopes:    []string{"read"},
	})
	require.NoError(t, err)

	err = repo.DeleteAPIKey(ctx, key.ID)
	require.NoError(t, err)

	// Verify deletion
	var count int
	err = testDB.Pool().QueryRow(ctx, "SELECT COUNT(*) FROM shared.api_keys WHERE id = $1", key.ID).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestRepositoryPg_DeleteExpiredAPIKeys(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Create expired AND revoked key (DeleteExpiredAPIKeys only deletes revoked keys)
	queries := db.New(testDB.Pool())
	key, err := queries.CreateAPIKey(ctx, db.CreateAPIKeyParams{
		UserID:    userID,
		Name:      "Expired Key",
		KeyHash:   "expired_hash",
		KeyPrefix: "rv_expired",
		Scopes:    []string{"read"},
		ExpiresAt: pgtype.Timestamptz{
			Time:  time.Now().Add(-91 * 24 * time.Hour),
			Valid: true,
		},
	})
	require.NoError(t, err)

	// Revoke it (SQL requires is_active = false)
	err = repo.RevokeAPIKey(ctx, key.ID)
	require.NoError(t, err)

	err = repo.DeleteExpiredAPIKeys(ctx)
	require.NoError(t, err)

	// Verify expired key was deleted
	var count int
	err = testDB.Pool().QueryRow(ctx, "SELECT COUNT(*) FROM shared.api_keys WHERE key_hash = 'expired_hash'").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func stringPtr(s string) *string {
	return &s
}
