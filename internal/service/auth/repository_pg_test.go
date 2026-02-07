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

func setupTestRepository(t *testing.T) (*RepositoryPG, testutil.DB) {
	t.Helper()
	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())
	repo := NewRepositoryPG(queries)
	return repo, testDB
}

func createTestUser(t *testing.T, testDB testutil.DB, username, email string) db.SharedUser {
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
		_, err := repo.GetUserByID(ctx, uuid.Must(uuid.NewV7()))
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

// ============================================================================
// Email Verification Token Tests (additional)
// ============================================================================

func TestRepositoryPG_InvalidateEmailVerificationTokensByEmail(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	// Create two users with different emails
	user1 := createTestUser(t, testDB, "invbyemail1", "invbyemail1@example.com")
	user2 := createTestUser(t, testDB, "invbyemail2", "invbyemail2@example.com")

	// Create tokens for user1's email
	for i := 0; i < 3; i++ {
		_, err := repo.CreateEmailVerificationToken(ctx, CreateEmailVerificationTokenParams{
			UserID:    user1.ID,
			TokenHash: "verify_byemail_u1_" + string(rune('a'+i)),
			Email:     user1.Email,
			ExpiresAt: time.Now().Add(24 * time.Hour),
		})
		require.NoError(t, err)
	}

	// Create tokens for user2's email
	for i := 0; i < 2; i++ {
		_, err := repo.CreateEmailVerificationToken(ctx, CreateEmailVerificationTokenParams{
			UserID:    user2.ID,
			TokenHash: "verify_byemail_u2_" + string(rune('a'+i)),
			Email:     user2.Email,
			ExpiresAt: time.Now().Add(24 * time.Hour),
		})
		require.NoError(t, err)
	}

	// Invalidate only user1's email tokens
	err := repo.InvalidateEmailVerificationTokensByEmail(ctx, user1.Email)
	require.NoError(t, err)

	// Verify all user1's tokens are invalidated
	var count1 int
	err = testDB.Pool().QueryRow(ctx,
		"SELECT COUNT(*) FROM shared.email_verification_tokens WHERE email = $1 AND verified_at IS NOT NULL",
		user1.Email,
	).Scan(&count1)
	require.NoError(t, err)
	assert.Equal(t, 3, count1)

	// Verify user2's tokens are NOT invalidated
	var count2 int
	err = testDB.Pool().QueryRow(ctx,
		"SELECT COUNT(*) FROM shared.email_verification_tokens WHERE email = $1 AND verified_at IS NULL",
		user2.Email,
	).Scan(&count2)
	require.NoError(t, err)
	assert.Equal(t, 2, count2)
}

func TestRepositoryPG_InvalidateEmailVerificationTokensByEmail_NoTokens(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	// Invalidating tokens for an email with no tokens should not error
	err := repo.InvalidateEmailVerificationTokensByEmail(ctx, "nobody@example.com")
	require.NoError(t, err)
}

func TestRepositoryPG_DeleteVerifiedEmailTokens(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	user := createTestUser(t, testDB, "delverified", "delverified@example.com")

	// Create a token and mark it as verified with an old timestamp (7+ days ago)
	token, err := repo.CreateEmailVerificationToken(ctx, CreateEmailVerificationTokenParams{
		UserID:    user.ID,
		TokenHash: "verify_delverified",
		Email:     user.Email,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})
	require.NoError(t, err)

	oldTime := time.Now().Add(-8 * 24 * time.Hour)
	_, err = testDB.Pool().Exec(ctx,
		"UPDATE shared.email_verification_tokens SET verified_at = $1 WHERE id = $2",
		oldTime, token.ID,
	)
	require.NoError(t, err)

	// Create another token that is NOT verified (should survive deletion)
	_, err = repo.CreateEmailVerificationToken(ctx, CreateEmailVerificationTokenParams{
		UserID:    user.ID,
		TokenHash: "verify_delverified_keep",
		Email:     user.Email,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})
	require.NoError(t, err)

	// Delete verified tokens
	err = repo.DeleteVerifiedEmailTokens(ctx)
	require.NoError(t, err)

	// Verify the verified token was deleted
	var countDeleted int
	err = testDB.Pool().QueryRow(ctx,
		"SELECT COUNT(*) FROM shared.email_verification_tokens WHERE token_hash = $1",
		"verify_delverified",
	).Scan(&countDeleted)
	require.NoError(t, err)
	assert.Equal(t, 0, countDeleted)

	// Verify the unverified token still exists
	var countKept int
	err = testDB.Pool().QueryRow(ctx,
		"SELECT COUNT(*) FROM shared.email_verification_tokens WHERE token_hash = $1",
		"verify_delverified_keep",
	).Scan(&countKept)
	require.NoError(t, err)
	assert.Equal(t, 1, countKept)
}

// ============================================================================
// Failed Login Attempts Tests
// ============================================================================

func TestRepositoryPG_RecordFailedLoginAttempt(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	err := repo.RecordFailedLoginAttempt(ctx, "testuser", "192.168.1.100")
	require.NoError(t, err)

	// Verify the record was created
	var count int
	err = testDB.Pool().QueryRow(ctx,
		"SELECT COUNT(*) FROM shared.failed_login_attempts WHERE username = $1 AND ip_address = $2",
		"testuser", "192.168.1.100",
	).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestRepositoryPG_RecordFailedLoginAttempt_Multiple(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	// Record multiple failed attempts
	for i := 0; i < 5; i++ {
		err := repo.RecordFailedLoginAttempt(ctx, "bruteforce", "10.0.0.1")
		require.NoError(t, err)
	}

	// Verify all records were created
	var count int
	err := testDB.Pool().QueryRow(ctx,
		"SELECT COUNT(*) FROM shared.failed_login_attempts WHERE username = $1",
		"bruteforce",
	).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 5, count)
}

func TestRepositoryPG_CountFailedLoginAttemptsByUsername(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	t.Run("no attempts", func(t *testing.T) {
		count, err := repo.CountFailedLoginAttemptsByUsername(ctx, "noattempts", time.Now().Add(-1*time.Hour))
		require.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})

	t.Run("counts only recent attempts", func(t *testing.T) {
		// Record some attempts now
		for i := 0; i < 3; i++ {
			err := repo.RecordFailedLoginAttempt(ctx, "countuser", "192.168.1.50")
			require.NoError(t, err)
		}

		// Count attempts in the last hour (should find all 3)
		count, err := repo.CountFailedLoginAttemptsByUsername(ctx, "countuser", time.Now().Add(-1*time.Hour))
		require.NoError(t, err)
		assert.Equal(t, int64(3), count)

		// Count attempts since the future (should find 0)
		count, err = repo.CountFailedLoginAttemptsByUsername(ctx, "countuser", time.Now().Add(1*time.Hour))
		require.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})

	t.Run("different usernames are isolated", func(t *testing.T) {
		err := repo.RecordFailedLoginAttempt(ctx, "userA", "10.0.0.1")
		require.NoError(t, err)
		err = repo.RecordFailedLoginAttempt(ctx, "userA", "10.0.0.1")
		require.NoError(t, err)
		err = repo.RecordFailedLoginAttempt(ctx, "userB", "10.0.0.1")
		require.NoError(t, err)

		countA, err := repo.CountFailedLoginAttemptsByUsername(ctx, "userA", time.Now().Add(-1*time.Hour))
		require.NoError(t, err)
		assert.Equal(t, int64(2), countA)

		countB, err := repo.CountFailedLoginAttemptsByUsername(ctx, "userB", time.Now().Add(-1*time.Hour))
		require.NoError(t, err)
		assert.Equal(t, int64(1), countB)
	})
}

func TestRepositoryPG_CountFailedLoginAttemptsByIP(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	t.Run("no attempts", func(t *testing.T) {
		count, err := repo.CountFailedLoginAttemptsByIP(ctx, "172.16.0.1", time.Now().Add(-1*time.Hour))
		require.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})

	t.Run("counts by IP across usernames", func(t *testing.T) {
		// Record attempts from same IP for different usernames
		err := repo.RecordFailedLoginAttempt(ctx, "ipuser1", "172.16.0.100")
		require.NoError(t, err)
		err = repo.RecordFailedLoginAttempt(ctx, "ipuser2", "172.16.0.100")
		require.NoError(t, err)
		err = repo.RecordFailedLoginAttempt(ctx, "ipuser3", "172.16.0.100")
		require.NoError(t, err)

		// Record attempt from a different IP
		err = repo.RecordFailedLoginAttempt(ctx, "ipuser1", "172.16.0.200")
		require.NoError(t, err)

		// Count for 172.16.0.100 should be 3
		count, err := repo.CountFailedLoginAttemptsByIP(ctx, "172.16.0.100", time.Now().Add(-1*time.Hour))
		require.NoError(t, err)
		assert.Equal(t, int64(3), count)

		// Count for 172.16.0.200 should be 1
		count, err = repo.CountFailedLoginAttemptsByIP(ctx, "172.16.0.200", time.Now().Add(-1*time.Hour))
		require.NoError(t, err)
		assert.Equal(t, int64(1), count)
	})

	t.Run("time window filtering", func(t *testing.T) {
		err := repo.RecordFailedLoginAttempt(ctx, "timeipuser", "172.16.0.50")
		require.NoError(t, err)

		// Count since the future should find 0
		count, err := repo.CountFailedLoginAttemptsByIP(ctx, "172.16.0.50", time.Now().Add(1*time.Hour))
		require.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})
}

func TestRepositoryPG_ClearFailedLoginAttemptsByUsername(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	// Record multiple failed attempts for the user
	for i := 0; i < 5; i++ {
		err := repo.RecordFailedLoginAttempt(ctx, "clearuser", "10.0.0."+string(rune('1'+i)))
		require.NoError(t, err)
	}

	// Record attempt for a different user
	err := repo.RecordFailedLoginAttempt(ctx, "otheruser", "10.0.0.1")
	require.NoError(t, err)

	// Clear attempts for "clearuser"
	err = repo.ClearFailedLoginAttemptsByUsername(ctx, "clearuser")
	require.NoError(t, err)

	// Verify clearuser's attempts are gone
	var clearCount int
	err = testDB.Pool().QueryRow(ctx,
		"SELECT COUNT(*) FROM shared.failed_login_attempts WHERE username = $1",
		"clearuser",
	).Scan(&clearCount)
	require.NoError(t, err)
	assert.Equal(t, 0, clearCount)

	// Verify otheruser's attempts are preserved
	var otherCount int
	err = testDB.Pool().QueryRow(ctx,
		"SELECT COUNT(*) FROM shared.failed_login_attempts WHERE username = $1",
		"otheruser",
	).Scan(&otherCount)
	require.NoError(t, err)
	assert.Equal(t, 1, otherCount)
}

func TestRepositoryPG_ClearFailedLoginAttemptsByUsername_NoAttempts(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	// Clearing when no attempts exist should not error
	err := repo.ClearFailedLoginAttemptsByUsername(ctx, "nonexistentuser")
	require.NoError(t, err)
}

func TestRepositoryPG_DeleteOldFailedLoginAttempts(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	// Insert an old attempt (> 24 hours ago) directly via SQL
	oldTime := time.Now().Add(-48 * time.Hour)
	_, err := testDB.Pool().Exec(ctx,
		"INSERT INTO shared.failed_login_attempts (username, ip_address, attempted_at) VALUES ($1, $2, $3)",
		"olduser", "10.0.0.1", oldTime,
	)
	require.NoError(t, err)

	// Insert a recent attempt
	err = repo.RecordFailedLoginAttempt(ctx, "recentuser", "10.0.0.2")
	require.NoError(t, err)

	// Delete old attempts
	err = repo.DeleteOldFailedLoginAttempts(ctx)
	require.NoError(t, err)

	// Verify old attempt was deleted
	var oldCount int
	err = testDB.Pool().QueryRow(ctx,
		"SELECT COUNT(*) FROM shared.failed_login_attempts WHERE username = $1",
		"olduser",
	).Scan(&oldCount)
	require.NoError(t, err)
	assert.Equal(t, 0, oldCount)

	// Verify recent attempt still exists
	var recentCount int
	err = testDB.Pool().QueryRow(ctx,
		"SELECT COUNT(*) FROM shared.failed_login_attempts WHERE username = $1",
		"recentuser",
	).Scan(&recentCount)
	require.NoError(t, err)
	assert.Equal(t, 1, recentCount)
}

func TestRepositoryPG_DeleteOldFailedLoginAttempts_NoOldAttempts(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	// Delete old attempts when none exist should not error
	err := repo.DeleteOldFailedLoginAttempts(ctx)
	require.NoError(t, err)
}

// ============================================================================
// Auth Token Tests (additional: GetAuthTokensByDeviceFingerprint)
// ============================================================================

func TestRepositoryPG_GetAuthTokensByDeviceFingerprint(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	user := createTestUser(t, testDB, "devicefp", "devicefp@example.com")

	fingerprint := "fp-abc-123"
	// Create tokens with matching fingerprint
	for i := 0; i < 2; i++ {
		_, err := repo.CreateAuthToken(ctx, CreateAuthTokenParams{
			UserID:            user.ID,
			TokenHash:         "hash_fp_match_" + string(rune('a'+i)),
			TokenType:         "refresh",
			DeviceFingerprint: &fingerprint,
			ExpiresAt:         time.Now().Add(24 * time.Hour),
		})
		require.NoError(t, err)
	}

	// Create token with different fingerprint
	otherFP := "fp-other"
	_, err := repo.CreateAuthToken(ctx, CreateAuthTokenParams{
		UserID:            user.ID,
		TokenHash:         "hash_fp_other",
		TokenType:         "refresh",
		DeviceFingerprint: &otherFP,
		ExpiresAt:         time.Now().Add(24 * time.Hour),
	})
	require.NoError(t, err)

	tokens, err := repo.GetAuthTokensByDeviceFingerprint(ctx, user.ID, fingerprint)
	require.NoError(t, err)
	assert.Len(t, tokens, 2)
	for _, tok := range tokens {
		assert.Equal(t, fingerprint, *tok.DeviceFingerprint)
	}
}

// Helper function
func stringPtr(s string) *string {
	return &s
}
