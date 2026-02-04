package session

import (
	"context"
	"net/netip"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
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

func setupTestRepository(t *testing.T) (*RepositoryPG, *testutil.TestDB) {
	t.Helper()
	testDB := testutil.NewTestDB(t)
	queries := db.New(testDB.Pool())
	repo := &RepositoryPG{queries: queries}
	return repo, testDB
}

func createTestUser(t *testing.T, testDB *testutil.TestDB) uuid.UUID {
	t.Helper()
	queries := db.New(testDB.Pool())
	user, err := queries.CreateUser(context.Background(), db.CreateUserParams{
		Username:     "testuser_" + uuid.New().String()[:8],
		Email:        "test_" + uuid.New().String()[:8] + "@example.com",
		PasswordHash: "hash",
	})
	require.NoError(t, err)
	return user.ID
}

func stringPtr(s string) *string {
	return &s
}

func TestRepositoryPG_CreateSession(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)
	ipAddr := netip.MustParseAddr("192.168.1.1")

	session, err := repo.CreateSession(ctx, CreateSessionParams{
		UserID:           userID,
		TokenHash:        "token_hash_123",
		RefreshTokenHash: stringPtr("refresh_hash_123"),
		IPAddress:        &ipAddr,
		UserAgent:        stringPtr("Mozilla/5.0"),
		DeviceName:       stringPtr("Chrome Browser"),
		Scopes:           []string{"read", "write"},
		ExpiresAt:        time.Now().Add(24 * time.Hour),
	})
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, session.ID)
	assert.Equal(t, userID, session.UserID)
}

func TestRepositoryPG_GetSessionByTokenHash(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	created, err := repo.CreateSession(ctx, CreateSessionParams{
		UserID:           userID,
		TokenHash:        "get_token_123",
		RefreshTokenHash: stringPtr("refresh_123"),
		ExpiresAt:        time.Now().Add(24 * time.Hour),
	})
	require.NoError(t, err)

	retrieved, err := repo.GetSessionByTokenHash(ctx, "get_token_123")
	require.NoError(t, err)
	require.NotNil(t, retrieved)
	assert.Equal(t, created.ID, retrieved.ID)
}

func TestRepositoryPG_GetSessionByID(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	created, err := repo.CreateSession(ctx, CreateSessionParams{
		UserID:           userID,
		TokenHash:        "id_test",
		RefreshTokenHash: stringPtr("refresh_id"),
		ExpiresAt:        time.Now().Add(24 * time.Hour),
	})
	require.NoError(t, err)

	retrieved, err := repo.GetSessionByID(ctx, created.ID)
	require.NoError(t, err)
	require.NotNil(t, retrieved)
	assert.Equal(t, created.ID, retrieved.ID)
	assert.Equal(t, userID, retrieved.UserID)
}

func TestRepositoryPG_GetSessionByRefreshTokenHash(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	refreshHash := "refresh_test_hash"
	created, err := repo.CreateSession(ctx, CreateSessionParams{
		UserID:           userID,
		TokenHash:        "token_refresh",
		RefreshTokenHash: &refreshHash,
		ExpiresAt:        time.Now().Add(24 * time.Hour),
	})
	require.NoError(t, err)

	retrieved, err := repo.GetSessionByRefreshTokenHash(ctx, refreshHash)
	require.NoError(t, err)
	require.NotNil(t, retrieved)
	assert.Equal(t, created.ID, retrieved.ID)
}

func TestRepositoryPG_ListUserSessions(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	for i := 0; i < 3; i++ {
		_, err := repo.CreateSession(ctx, CreateSessionParams{
			UserID:           userID,
			TokenHash:        "list_" + uuid.New().String()[:8],
			RefreshTokenHash: stringPtr("ref_" + uuid.New().String()[:8]),
			ExpiresAt:        time.Now().Add(24 * time.Hour),
		})
		require.NoError(t, err)
	}

	sessions, err := repo.ListUserSessions(ctx, userID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(sessions), 3)
}

func TestRepositoryPG_ListAllUserSessions(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	_, err := repo.CreateSession(ctx, CreateSessionParams{
		UserID:           userID,
		TokenHash:        "active",
		RefreshTokenHash: stringPtr("ref_active"),
		ExpiresAt:        time.Now().Add(24 * time.Hour),
	})
	require.NoError(t, err)

	_, err = repo.CreateSession(ctx, CreateSessionParams{
		UserID:           userID,
		TokenHash:        "expired",
		RefreshTokenHash: stringPtr("ref_expired"),
		ExpiresAt:        time.Now().Add(-24 * time.Hour),
	})
	require.NoError(t, err)

	allSessions, err := repo.ListAllUserSessions(ctx, userID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(allSessions), 2)
}

func TestRepositoryPG_CountActiveUserSessions(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	for i := 0; i < 3; i++ {
		_, err := repo.CreateSession(ctx, CreateSessionParams{
			UserID:           userID,
			TokenHash:        "count_" + uuid.New().String()[:8],
			RefreshTokenHash: stringPtr("cref_" + uuid.New().String()[:8]),
			ExpiresAt:        time.Now().Add(24 * time.Hour),
		})
		require.NoError(t, err)
	}

	count, err := repo.CountActiveUserSessions(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)
}

func TestRepositoryPG_UpdateSessionActivity(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	session, err := repo.CreateSession(ctx, CreateSessionParams{
		UserID:           userID,
		TokenHash:        "activity_test",
		RefreshTokenHash: stringPtr("ref_activity"),
		ExpiresAt:        time.Now().Add(24 * time.Hour),
	})
	require.NoError(t, err)

	originalActivity := session.LastActivityAt

	time.Sleep(10 * time.Millisecond)
	err = repo.UpdateSessionActivity(ctx, session.ID)
	require.NoError(t, err)

	updated, err := repo.GetSessionByID(ctx, session.ID)
	require.NoError(t, err)
	assert.True(t, updated.LastActivityAt.After(originalActivity))
}

func TestRepositoryPG_UpdateSessionActivityByTokenHash(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	session, err := repo.CreateSession(ctx, CreateSessionParams{
		UserID:           userID,
		TokenHash:        "activity_token_test",
		RefreshTokenHash: stringPtr("ref_activity_token"),
		ExpiresAt:        time.Now().Add(24 * time.Hour),
	})
	require.NoError(t, err)

	originalActivity := session.LastActivityAt

	time.Sleep(10 * time.Millisecond)
	err = repo.UpdateSessionActivityByTokenHash(ctx, "activity_token_test")
	require.NoError(t, err)

	updated, err := repo.GetSessionByTokenHash(ctx, "activity_token_test")
	require.NoError(t, err)
	assert.True(t, updated.LastActivityAt.After(originalActivity))
}

func TestRepositoryPG_RevokeSession(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	session, err := repo.CreateSession(ctx, CreateSessionParams{
		UserID:           userID,
		TokenHash:        "revoke_test",
		RefreshTokenHash: stringPtr("refresh_revoke"),
		ExpiresAt:        time.Now().Add(24 * time.Hour),
	})
	require.NoError(t, err)

	reason := "test"
	err = repo.RevokeSession(ctx, session.ID, &reason)
	require.NoError(t, err)

	// GetSessionByTokenHash filters out revoked sessions, so check DB directly
	var revokedAt *time.Time
	err = testDB.Pool().QueryRow(ctx, "SELECT revoked_at FROM shared.sessions WHERE id = $1", session.ID).Scan(&revokedAt)
	require.NoError(t, err)
	assert.NotNil(t, revokedAt, "Session should be revoked")
}

func TestRepositoryPG_RevokeSessionByTokenHash(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	_, err := repo.CreateSession(ctx, CreateSessionParams{
		UserID:           userID,
		TokenHash:        "revoke_by_token",
		RefreshTokenHash: stringPtr("refresh_revoke_token"),
		ExpiresAt:        time.Now().Add(24 * time.Hour),
	})
	require.NoError(t, err)

	reason := "security"
	err = repo.RevokeSessionByTokenHash(ctx, "revoke_by_token", &reason)
	require.NoError(t, err)

	// Revoked sessions are filtered out, check DB directly
	var revokedAt *time.Time
	err = testDB.Pool().QueryRow(ctx, "SELECT revoked_at FROM shared.sessions WHERE token_hash = 'revoke_by_token'").Scan(&revokedAt)
	require.NoError(t, err)
	assert.NotNil(t, revokedAt)
}

func TestRepositoryPG_RevokeAllUserSessions(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	for i := 0; i < 3; i++ {
		_, err := repo.CreateSession(ctx, CreateSessionParams{
			UserID:           userID,
			TokenHash:        "revoke_all_" + uuid.New().String()[:8],
			RefreshTokenHash: stringPtr("ref_all_" + uuid.New().String()[:8]),
			ExpiresAt:        time.Now().Add(24 * time.Hour),
		})
		require.NoError(t, err)
	}

	reason := "logout_all"
	err := repo.RevokeAllUserSessions(ctx, userID, &reason)
	require.NoError(t, err)

	sessions, err := repo.ListUserSessions(ctx, userID)
	require.NoError(t, err)
	assert.Empty(t, sessions)
}

func TestRepositoryPG_RevokeAllUserSessionsExcept(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	var keepSession db.SharedSession
	for i := 0; i < 3; i++ {
		s, err := repo.CreateSession(ctx, CreateSessionParams{
			UserID:           userID,
			TokenHash:        "except_" + uuid.New().String()[:8],
			RefreshTokenHash: stringPtr("ref_except_" + uuid.New().String()[:8]),
			ExpiresAt:        time.Now().Add(24 * time.Hour),
		})
		require.NoError(t, err)
		if i == 0 {
			keepSession = s
		}
	}

	reason := "revoke_others"
	err := repo.RevokeAllUserSessionsExcept(ctx, userID, keepSession.ID, &reason)
	require.NoError(t, err)

	sessions, err := repo.ListUserSessions(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, sessions, 1)
	assert.Equal(t, keepSession.ID, sessions[0].ID)
}

func TestRepositoryPG_GetInactiveSessions(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	_, err := repo.CreateSession(ctx, CreateSessionParams{
		UserID:           userID,
		TokenHash:        "inactive_test",
		RefreshTokenHash: stringPtr("ref_inactive"),
		ExpiresAt:        time.Now().Add(24 * time.Hour),
	})
	require.NoError(t, err)

	inactiveSince := time.Now().Add(1 * time.Hour)
	sessions, err := repo.GetInactiveSessions(ctx, inactiveSince)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(sessions), 1)
}

func TestRepositoryPG_RevokeInactiveSessions(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	session, err := repo.CreateSession(ctx, CreateSessionParams{
		UserID:           userID,
		TokenHash:        "revoke_inactive",
		RefreshTokenHash: stringPtr("ref_revoke_inactive"),
		ExpiresAt:        time.Now().Add(24 * time.Hour),
	})
	require.NoError(t, err)

	inactiveSince := time.Now().Add(1 * time.Hour)
	err = repo.RevokeInactiveSessions(ctx, inactiveSince)
	require.NoError(t, err)

	var revokedAt *time.Time
	err = testDB.Pool().QueryRow(ctx, "SELECT revoked_at FROM shared.sessions WHERE id = $1", session.ID).Scan(&revokedAt)
	require.NoError(t, err)
	assert.NotNil(t, revokedAt)
}

func TestRepositoryPG_DeleteExpiredSessions(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	queries := db.New(testDB.Pool())
	_, err := queries.CreateSession(ctx, db.CreateSessionParams{
		UserID:    userID,
		TokenHash: "delete_expired",
		ExpiresAt: time.Now().Add(-91 * 24 * time.Hour),
	})
	require.NoError(t, err)

	deletedCount, err := repo.DeleteExpiredSessions(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, deletedCount, int64(1))

	var count int
	err = testDB.Pool().QueryRow(ctx, "SELECT COUNT(*) FROM shared.sessions WHERE token_hash = 'delete_expired'").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestRepositoryPG_DeleteRevokedSessions(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	session, err := repo.CreateSession(ctx, CreateSessionParams{
		UserID:           userID,
		TokenHash:        "delete_revoked",
		RefreshTokenHash: stringPtr("ref_delete_revoked"),
		ExpiresAt:        time.Now().Add(24 * time.Hour),
	})
	require.NoError(t, err)

	reason := "delete_test"
	err = repo.RevokeSession(ctx, session.ID, &reason)
	require.NoError(t, err)

	_, err = testDB.Pool().Exec(ctx, "UPDATE shared.sessions SET revoked_at = $1 WHERE id = $2",
		time.Now().Add(-31*24*time.Hour), session.ID)
	require.NoError(t, err)

	deletedCount, err := repo.DeleteRevokedSessions(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, deletedCount, int64(1))

	var count int
	err = testDB.Pool().QueryRow(ctx, "SELECT COUNT(*) FROM shared.sessions WHERE id = $1", session.ID).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}
