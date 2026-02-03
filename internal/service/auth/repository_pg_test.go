package auth

import (
	"context"
	"net/netip"
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

func setupTestRepository(t *testing.T) (*RepositoryPG, *testutil.TestDB) {
	t.Helper()
	testDB := testutil.NewTestDB(t)
	queries := db.New(testDB.Pool())
	repo := NewRepositoryPG(queries)
	return repo, testDB
}

func createTestUser(t *testing.T, testDB *testutil.TestDB, username, email string) db.SharedUser {
	t.Helper()
	queries := db.New(testDB.Pool())
	user, err := queries.CreateUser(context.Background(), db.CreateUserParams{
		Username:     username,
		Email:        email,
		PasswordHash: "hash123",
	})
	require.NoError(t, err)
	return user
}

// ============================================================================
// User Operations Tests
// ============================================================================

func TestRepositoryPG_CreateUser(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	t.Run("create user successfully", func(t *testing.T) {
		user, err := repo.CreateUser(ctx, db.CreateUserParams{
			Username:     "testuser",
			Email:        "test@example.com",
			PasswordHash: "hash123",
		})
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, user.ID)
		assert.Equal(t, "testuser", user.Username)
		assert.Equal(t, "test@example.com", user.Email)
		assert.Equal(t, "hash123", user.PasswordHash)
	})

	t.Run("duplicate username fails", func(t *testing.T) {
		_, err := repo.CreateUser(ctx, db.CreateUserParams{
			Username:     "duplicate",
			Email:        "user1@example.com",
			PasswordHash: "hash123",
		})
		require.NoError(t, err)

		_, err = repo.CreateUser(ctx, db.CreateUserParams{
			Username:     "duplicate",
			Email:        "user2@example.com",
			PasswordHash: "hash123",
		})
		require.Error(t, err)
	})

	t.Run("duplicate email fails", func(t *testing.T) {
		_, err := repo.CreateUser(ctx, db.CreateUserParams{
			Username:     "user1",
			Email:        "same@example.com",
			PasswordHash: "hash123",
		})
		require.NoError(t, err)

		_, err = repo.CreateUser(ctx, db.CreateUserParams{
			Username:     "user2",
			Email:        "same@example.com",
			PasswordHash: "hash123",
		})
		require.Error(t, err)
	})
}

func TestRepositoryPG_GetUserByID(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	user := createTestUser(t, testDB, "getbyid", "getbyid@example.com")

	t.Run("existing user", func(t *testing.T) {
		retrieved, err := repo.GetUserByID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, user.ID, retrieved.ID)
		assert.Equal(t, user.Username, retrieved.Username)
	})

	t.Run("non-existent user", func(t *testing.T) {
		_, err := repo.GetUserByID(ctx, uuid.New())
		require.Error(t, err)
	})
}

func TestRepositoryPG_GetUserByUsername(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	user := createTestUser(t, testDB, "getbyusername", "getbyusername@example.com")

	t.Run("existing username", func(t *testing.T) {
		retrieved, err := repo.GetUserByUsername(ctx, user.Username)
		require.NoError(t, err)
		assert.Equal(t, user.ID, retrieved.ID)
		assert.Equal(t, user.Username, retrieved.Username)
	})

	t.Run("non-existent username", func(t *testing.T) {
		_, err := repo.GetUserByUsername(ctx, "nonexistent")
		require.Error(t, err)
	})
}

func TestRepositoryPG_GetUserByEmail(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	user := createTestUser(t, testDB, "getbyemail", "getbyemail@example.com")

	t.Run("existing email", func(t *testing.T) {
		retrieved, err := repo.GetUserByEmail(ctx, user.Email)
		require.NoError(t, err)
		assert.Equal(t, user.ID, retrieved.ID)
		assert.Equal(t, user.Email, retrieved.Email)
	})

	t.Run("non-existent email", func(t *testing.T) {
		_, err := repo.GetUserByEmail(ctx, "nonexistent@example.com")
		require.Error(t, err)
	})
}

func TestRepositoryPG_UpdateUserPassword(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	user := createTestUser(t, testDB, "updatepass", "updatepass@example.com")

	err := repo.UpdateUserPassword(ctx, user.ID, "newhash456")
	require.NoError(t, err)

	// Verify password was updated
	updated, err := repo.GetUserByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, "newhash456", updated.PasswordHash)
}

func TestRepositoryPG_UpdateUserEmailVerified(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	user := createTestUser(t, testDB, "verifyemail", "verifyemail@example.com")

	err := repo.UpdateUserEmailVerified(ctx, user.ID, true)
	require.NoError(t, err)

	// Verify email was verified
	updated, err := repo.GetUserByID(ctx, user.ID)
	require.NoError(t, err)
	require.NotNil(t, updated.EmailVerified)
	assert.True(t, *updated.EmailVerified)
}

func TestRepositoryPG_UpdateUserLastLogin(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	user := createTestUser(t, testDB, "lastlogin", "lastlogin@example.com")

	time.Sleep(10 * time.Millisecond) // Ensure time difference

	err := repo.UpdateUserLastLogin(ctx, user.ID)
	require.NoError(t, err)

	// Verify last login was updated
	updated, err := repo.GetUserByID(ctx, user.ID)
	require.NoError(t, err)
	assert.True(t, updated.LastLoginAt.Valid)
	assert.True(t, updated.LastLoginAt.Time.After(user.CreatedAt))
}

// ============================================================================
// Auth Token Tests
// ============================================================================

func TestRepositoryPG_CreateAuthToken(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	user := createTestUser(t, testDB, "authtoken", "authtoken@example.com")
	ipAddr := netip.MustParseAddr("192.168.1.1")

	token, err := repo.CreateAuthToken(ctx, CreateAuthTokenParams{
		UserID:            user.ID,
		TokenHash:         "hash_abc123",
		TokenType:         "refresh",
		DeviceName:        stringPtr("Chrome Browser"),
		DeviceFingerprint: stringPtr("fingerprint123"),
		IPAddress:         &ipAddr,
		UserAgent:         stringPtr("Mozilla/5.0"),
		ExpiresAt:         time.Now().Add(24 * time.Hour),
	})
	require.NoError(t, err)

	assert.NotEqual(t, uuid.Nil, token.ID)
	assert.Equal(t, user.ID, token.UserID)
	assert.Equal(t, "hash_abc123", token.TokenHash)
	assert.Equal(t, "refresh", token.TokenType)
	assert.Equal(t, "Chrome Browser", *token.DeviceName)
	assert.Nil(t, token.RevokedAt)
}

func TestRepositoryPG_GetAuthTokenByHash(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	user := createTestUser(t, testDB, "gettoken", "gettoken@example.com")

	created, err := repo.CreateAuthToken(ctx, CreateAuthTokenParams{
		UserID:    user.ID,
		TokenHash: "hash_get123",
		TokenType: "refresh",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})
	require.NoError(t, err)

	t.Run("existing token", func(t *testing.T) {
		retrieved, err := repo.GetAuthTokenByHash(ctx, "hash_get123")
		require.NoError(t, err)
		assert.Equal(t, created.ID, retrieved.ID)
		assert.Equal(t, created.TokenHash, retrieved.TokenHash)
	})

	t.Run("non-existent token", func(t *testing.T) {
		_, err := repo.GetAuthTokenByHash(ctx, "hash_nonexistent")
		require.Error(t, err)
	})
}

func TestRepositoryPG_GetAuthTokensByUserID(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	user := createTestUser(t, testDB, "usertokens", "usertokens@example.com")

	// Create multiple tokens
	for i := 0; i < 3; i++ {
		_, err := repo.CreateAuthToken(ctx, CreateAuthTokenParams{
			UserID:    user.ID,
			TokenHash: "hash_user_" + string(rune('a'+i)),
			TokenType: "refresh",
			ExpiresAt: time.Now().Add(24 * time.Hour),
		})
		require.NoError(t, err)
	}

	tokens, err := repo.GetAuthTokensByUserID(ctx, user.ID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(tokens), 3)
}

func TestRepositoryPG_UpdateAuthTokenLastUsed(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	user := createTestUser(t, testDB, "lastused", "lastused@example.com")

	token, err := repo.CreateAuthToken(ctx, CreateAuthTokenParams{
		UserID:    user.ID,
		TokenHash: "hash_lastused",
		TokenType: "refresh",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})
	require.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	err = repo.UpdateAuthTokenLastUsed(ctx, token.ID)
	require.NoError(t, err)

	// Verify last used was updated
	updated, err := repo.GetAuthTokenByHash(ctx, token.TokenHash)
	require.NoError(t, err)
	require.NotNil(t, updated.LastUsedAt)
	assert.True(t, updated.LastUsedAt.After(token.CreatedAt))
}

func TestRepositoryPG_RevokeAuthToken(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	user := createTestUser(t, testDB, "revoke", "revoke@example.com")

	token, err := repo.CreateAuthToken(ctx, CreateAuthTokenParams{
		UserID:    user.ID,
		TokenHash: "hash_revoke",
		TokenType: "refresh",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})
	require.NoError(t, err)

	err = repo.RevokeAuthToken(ctx, token.ID)
	require.NoError(t, err)

	// Verify token was revoked - GetAuthTokenByHash filters revoked tokens
	_, err = repo.GetAuthTokenByHash(ctx, token.TokenHash)
	assert.Error(t, err) // Should not find revoked token

	// Also verify via direct DB query that revoked_at is set
	var revokedAt pgtype.Timestamp
	err = testDB.Pool().QueryRow(ctx,
		"SELECT revoked_at FROM shared.auth_tokens WHERE id = $1",
		token.ID,
	).Scan(&revokedAt)
	require.NoError(t, err)
	assert.True(t, revokedAt.Valid)
}

func TestRepositoryPG_RevokeAuthTokenByHash(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	user := createTestUser(t, testDB, "revokehash", "revokehash@example.com")

	_, err := repo.CreateAuthToken(ctx, CreateAuthTokenParams{
		UserID:    user.ID,
		TokenHash: "hash_revokehash",
		TokenType: "refresh",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})
	require.NoError(t, err)

	err = repo.RevokeAuthTokenByHash(ctx, "hash_revokehash")
	require.NoError(t, err)

	// Verify token was revoked - GetAuthTokenByHash filters revoked tokens
	_, err = repo.GetAuthTokenByHash(ctx, "hash_revokehash")
	assert.Error(t, err) // Should not find revoked token
}

func TestRepositoryPG_RevokeAllUserAuthTokens(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	user := createTestUser(t, testDB, "revokeall", "revokeall@example.com")

	// Create multiple tokens
	for i := 0; i < 3; i++ {
		_, err := repo.CreateAuthToken(ctx, CreateAuthTokenParams{
			UserID:    user.ID,
			TokenHash: "hash_revokeall_" + string(rune('a'+i)),
			TokenType: "refresh",
			ExpiresAt: time.Now().Add(24 * time.Hour),
		})
		require.NoError(t, err)
	}

	err := repo.RevokeAllUserAuthTokens(ctx, user.ID)
	require.NoError(t, err)

	// Verify all tokens are revoked - GetAuthTokensByUserID filters revoked tokens
	tokens, err := repo.GetAuthTokensByUserID(ctx, user.ID)
	require.NoError(t, err)
	assert.Len(t, tokens, 0) // All tokens should be filtered out
}

func TestRepositoryPG_RevokeAllUserAuthTokensExcept(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	user := createTestUser(t, testDB, "revokeexcept", "revokeexcept@example.com")

	// Create multiple tokens
	var keepToken AuthToken
	for i := 0; i < 3; i++ {
		token, err := repo.CreateAuthToken(ctx, CreateAuthTokenParams{
			UserID:    user.ID,
			TokenHash: "hash_except_" + string(rune('a'+i)),
			TokenType: "refresh",
			ExpiresAt: time.Now().Add(24 * time.Hour),
		})
		require.NoError(t, err)
		if i == 0 {
			keepToken = token
		}
	}

	err := repo.RevokeAllUserAuthTokensExcept(ctx, user.ID, keepToken.ID)
	require.NoError(t, err)

	// Verify only one token remains active (GetAuthTokensByUserID filters revoked)
	tokens, err := repo.GetAuthTokensByUserID(ctx, user.ID)
	require.NoError(t, err)
	assert.Len(t, tokens, 1)
	assert.Equal(t, keepToken.ID, tokens[0].ID)
}

func TestRepositoryPG_CountActiveAuthTokensByUser(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	user := createTestUser(t, testDB, "countactive", "countactive@example.com")

	// Create 3 active and 2 revoked tokens
	for i := 0; i < 5; i++ {
		token, err := repo.CreateAuthToken(ctx, CreateAuthTokenParams{
			UserID:    user.ID,
			TokenHash: "hash_count_" + string(rune('a'+i)),
			TokenType: "refresh",
			ExpiresAt: time.Now().Add(24 * time.Hour),
		})
		require.NoError(t, err)

		if i < 2 {
			err = repo.RevokeAuthToken(ctx, token.ID)
			require.NoError(t, err)
		}
	}

	count, err := repo.CountActiveAuthTokensByUser(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)
}

func TestRepositoryPG_DeleteExpiredAuthTokens(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	user := createTestUser(t, testDB, "delexpired", "delexpired@example.com")

	// Create expired token
	queries := db.New(testDB.Pool())
	_, err := queries.CreateAuthToken(ctx, db.CreateAuthTokenParams{
		UserID:    user.ID,
		TokenHash: "hash_expired",
		TokenType: "refresh",
		ExpiresAt: time.Now().Add(-1 * time.Hour), // Expired
	})
	require.NoError(t, err)

	err = repo.DeleteExpiredAuthTokens(ctx)
	require.NoError(t, err)

	// Verify expired token was deleted (not even in DB)
	var count int
	err = testDB.Pool().QueryRow(ctx,
		"SELECT COUNT(*) FROM shared.auth_tokens WHERE token_hash = $1",
		"hash_expired",
	).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestRepositoryPG_DeleteRevokedAuthTokens(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	user := createTestUser(t, testDB, "delrevoked", "delrevoked@example.com")

	token, err := repo.CreateAuthToken(ctx, CreateAuthTokenParams{
		UserID:    user.ID,
		TokenHash: "hash_delrevoked",
		TokenType: "refresh",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})
	require.NoError(t, err)

	// Revoke token with old timestamp (30+ days ago)
	oldTime := time.Now().Add(-31 * 24 * time.Hour)
	_, err = testDB.Pool().Exec(ctx,
		"UPDATE shared.auth_tokens SET revoked_at = $1 WHERE id = $2",
		oldTime, token.ID,
	)
	require.NoError(t, err)

	err = repo.DeleteRevokedAuthTokens(ctx)
	require.NoError(t, err)

	// Verify token no longer exists (not even in DB)
	var count int
	err = testDB.Pool().QueryRow(ctx,
		"SELECT COUNT(*) FROM shared.auth_tokens WHERE id = $1",
		token.ID,
	).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

// ============================================================================
// Password Reset Token Tests
// ============================================================================

func TestRepositoryPG_CreatePasswordResetToken(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	user := createTestUser(t, testDB, "resettoken", "resettoken@example.com")
	ipAddr := netip.MustParseAddr("10.0.0.1")

	token, err := repo.CreatePasswordResetToken(ctx, CreatePasswordResetTokenParams{
		UserID:    user.ID,
		TokenHash: "reset_hash_123",
		IPAddress: &ipAddr,
		UserAgent: stringPtr("Mozilla/5.0"),
		ExpiresAt: time.Now().Add(1 * time.Hour),
	})
	require.NoError(t, err)

	assert.NotEqual(t, uuid.Nil, token.ID)
	assert.Equal(t, user.ID, token.UserID)
	assert.Equal(t, "reset_hash_123", token.TokenHash)
	assert.Nil(t, token.UsedAt)
}

func TestRepositoryPG_GetPasswordResetToken(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	user := createTestUser(t, testDB, "getreset", "getreset@example.com")

	created, err := repo.CreatePasswordResetToken(ctx, CreatePasswordResetTokenParams{
		UserID:    user.ID,
		TokenHash: "reset_get_123",
		ExpiresAt: time.Now().Add(1 * time.Hour),
	})
	require.NoError(t, err)

	t.Run("existing token", func(t *testing.T) {
		retrieved, err := repo.GetPasswordResetToken(ctx, "reset_get_123")
		require.NoError(t, err)
		assert.Equal(t, created.ID, retrieved.ID)
		assert.Equal(t, created.TokenHash, retrieved.TokenHash)
	})

	t.Run("non-existent token", func(t *testing.T) {
		_, err := repo.GetPasswordResetToken(ctx, "reset_nonexistent")
		require.Error(t, err)
	})
}

func TestRepositoryPG_MarkPasswordResetTokenUsed(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	user := createTestUser(t, testDB, "markused", "markused@example.com")

	token, err := repo.CreatePasswordResetToken(ctx, CreatePasswordResetTokenParams{
		UserID:    user.ID,
		TokenHash: "reset_markused",
		ExpiresAt: time.Now().Add(1 * time.Hour),
	})
	require.NoError(t, err)

	err = repo.MarkPasswordResetTokenUsed(ctx, token.ID)
	require.NoError(t, err)

	// Verify token was marked as used via direct DB query
	var usedAt pgtype.Timestamp
	err = testDB.Pool().QueryRow(ctx,
		"SELECT used_at FROM shared.password_reset_tokens WHERE id = $1",
		token.ID,
	).Scan(&usedAt)
	require.NoError(t, err)
	assert.True(t, usedAt.Valid)
}

func TestRepositoryPG_InvalidateUserPasswordResetTokens(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	user := createTestUser(t, testDB, "invalidate", "invalidate@example.com")

	// Create multiple reset tokens
	for i := 0; i < 3; i++ {
		_, err := repo.CreatePasswordResetToken(ctx, CreatePasswordResetTokenParams{
			UserID:    user.ID,
			TokenHash: "reset_inv_" + string(rune('a'+i)),
			ExpiresAt: time.Now().Add(1 * time.Hour),
		})
		require.NoError(t, err)
	}

	err := repo.InvalidateUserPasswordResetTokens(ctx, user.ID)
	require.NoError(t, err)

	// Verify all tokens are marked as used via direct DB query
	var count int
	err = testDB.Pool().QueryRow(ctx,
		"SELECT COUNT(*) FROM shared.password_reset_tokens WHERE user_id = $1 AND used_at IS NOT NULL",
		user.ID,
	).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 3, count)
}

func TestRepositoryPG_DeleteExpiredPasswordResetTokens(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	user := createTestUser(t, testDB, "delresetexp", "delresetexp@example.com")

	// Create expired reset token
	queries := db.New(testDB.Pool())
	_, err := queries.CreatePasswordResetToken(ctx, db.CreatePasswordResetTokenParams{
		UserID:    user.ID,
		TokenHash: "reset_expired",
		ExpiresAt: time.Now().Add(-1 * time.Hour), // Expired
	})
	require.NoError(t, err)

	err = repo.DeleteExpiredPasswordResetTokens(ctx)
	require.NoError(t, err)

	// Verify expired token was deleted (not in DB)
	var count int
	err = testDB.Pool().QueryRow(ctx,
		"SELECT COUNT(*) FROM shared.password_reset_tokens WHERE token_hash = $1",
		"reset_expired",
	).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestRepositoryPG_DeleteUsedPasswordResetTokens(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	user := createTestUser(t, testDB, "delresetused", "delresetused@example.com")

	token, err := repo.CreatePasswordResetToken(ctx, CreatePasswordResetTokenParams{
		UserID:    user.ID,
		TokenHash: "reset_delused",
		ExpiresAt: time.Now().Add(1 * time.Hour),
	})
	require.NoError(t, err)

	// Mark as used with old timestamp (7+ days ago)
	oldTime := time.Now().Add(-8 * 24 * time.Hour)
	_, err = testDB.Pool().Exec(ctx,
		"UPDATE shared.password_reset_tokens SET used_at = $1 WHERE id = $2",
		oldTime, token.ID,
	)
	require.NoError(t, err)

	err = repo.DeleteUsedPasswordResetTokens(ctx)
	require.NoError(t, err)

	// Verify token was deleted (not in DB)
	var count int
	err = testDB.Pool().QueryRow(ctx,
		"SELECT COUNT(*) FROM shared.password_reset_tokens WHERE id = $1",
		token.ID,
	).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

// ============================================================================
// Email Verification Token Tests
// ============================================================================

func TestRepositoryPG_CreateEmailVerificationToken(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	user := createTestUser(t, testDB, "emailverify", "emailverify@example.com")
	ipAddr := netip.MustParseAddr("172.16.0.1")

	token, err := repo.CreateEmailVerificationToken(ctx, CreateEmailVerificationTokenParams{
		UserID:    user.ID,
		TokenHash: "verify_hash_123",
		Email:     user.Email,
		IPAddress: &ipAddr,
		UserAgent: stringPtr("Mozilla/5.0"),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})
	require.NoError(t, err)

	assert.NotEqual(t, uuid.Nil, token.ID)
	assert.Equal(t, user.ID, token.UserID)
	assert.Equal(t, "verify_hash_123", token.TokenHash)
	assert.Nil(t, token.VerifiedAt)
}

func TestRepositoryPG_GetEmailVerificationToken(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	user := createTestUser(t, testDB, "getverify", "getverify@example.com")

	created, err := repo.CreateEmailVerificationToken(ctx, CreateEmailVerificationTokenParams{
		UserID:    user.ID,
		TokenHash: "verify_get_123",
		Email:     user.Email,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})
	require.NoError(t, err)

	t.Run("existing token", func(t *testing.T) {
		retrieved, err := repo.GetEmailVerificationToken(ctx, "verify_get_123")
		require.NoError(t, err)
		assert.Equal(t, created.ID, retrieved.ID)
		assert.Equal(t, created.TokenHash, retrieved.TokenHash)
	})

	t.Run("non-existent token", func(t *testing.T) {
		_, err := repo.GetEmailVerificationToken(ctx, "verify_nonexistent")
		require.Error(t, err)
	})
}

func TestRepositoryPG_MarkEmailVerificationTokenUsed(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	user := createTestUser(t, testDB, "markverifyused", "markverifyused@example.com")

	token, err := repo.CreateEmailVerificationToken(ctx, CreateEmailVerificationTokenParams{
		UserID:    user.ID,
		TokenHash: "verify_markused",
		Email:     user.Email,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})
	require.NoError(t, err)

	err = repo.MarkEmailVerificationTokenUsed(ctx, token.ID)
	require.NoError(t, err)

	// Verify token was marked as verified via direct DB query
	var verifiedAt pgtype.Timestamp
	err = testDB.Pool().QueryRow(ctx,
		"SELECT verified_at FROM shared.email_verification_tokens WHERE id = $1",
		token.ID,
	).Scan(&verifiedAt)
	require.NoError(t, err)
	assert.True(t, verifiedAt.Valid)
}

func TestRepositoryPG_InvalidateUserEmailVerificationTokens(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	user := createTestUser(t, testDB, "invverify", "invverify@example.com")

	// Create multiple verification tokens
	for i := 0; i < 3; i++ {
		_, err := repo.CreateEmailVerificationToken(ctx, CreateEmailVerificationTokenParams{
			UserID:    user.ID,
			TokenHash: "verify_inv_" + string(rune('a'+i)),
			Email:     user.Email,
			ExpiresAt: time.Now().Add(24 * time.Hour),
		})
		require.NoError(t, err)
	}

	err := repo.InvalidateUserEmailVerificationTokens(ctx, user.ID)
	require.NoError(t, err)

	// Verify all tokens are marked as verified via direct DB query
	var count int
	err = testDB.Pool().QueryRow(ctx,
		"SELECT COUNT(*) FROM shared.email_verification_tokens WHERE user_id = $1 AND verified_at IS NOT NULL",
		user.ID,
	).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 3, count)
}

func TestRepositoryPG_DeleteExpiredEmailVerificationTokens(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	user := createTestUser(t, testDB, "delverifyexp", "delverifyexp@example.com")

	// Create expired verification token
	queries := db.New(testDB.Pool())
	_, err := queries.CreateEmailVerificationToken(ctx, db.CreateEmailVerificationTokenParams{
		UserID:    user.ID,
		TokenHash: "verify_expired",
		ExpiresAt: time.Now().Add(-1 * time.Hour), // Expired
	})
	require.NoError(t, err)

	err = repo.DeleteExpiredEmailVerificationTokens(ctx)
	require.NoError(t, err)

	// Verify expired token was deleted (not in DB)
	var count int
	err = testDB.Pool().QueryRow(ctx,
		"SELECT COUNT(*) FROM shared.email_verification_tokens WHERE token_hash = $1",
		"verify_expired",
	).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

// Helper function
func stringPtr(s string) *string {
	return &s
}
