package auth

import (
	"encoding/hex"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// JWT Manager Tests
// ============================================================================

func TestJWTManager_GenerateAccessToken(t *testing.T) {
	t.Parallel()

	manager := NewTokenManager("test-secret-key", 15*time.Minute)
	userID := uuid.Must(uuid.NewV7())
	username := "testuser"

	token, err := manager.GenerateAccessToken(userID, username)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// JWT should have 3 parts separated by dots
	parts := strings.Split(token, ".")
	assert.Len(t, parts, 3)

	// Should be able to validate the token
	claims, err := manager.ValidateAccessToken(token)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)
	assert.True(t, claims.ExpiresAt > time.Now().Unix())
}

func TestJWTManager_ValidateAccessToken(t *testing.T) {
	t.Parallel()

	manager := NewTokenManager("test-secret-key", 15*time.Minute)
	userID := uuid.Must(uuid.NewV7())

	t.Run("valid token", func(t *testing.T) {
		token, err := manager.GenerateAccessToken(userID, "user1")
		require.NoError(t, err)

		claims, err := manager.ValidateAccessToken(token)
		require.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, "user1", claims.Username)
	})

	t.Run("expired token", func(t *testing.T) {
		expiredManager := NewTokenManager("test-secret-key", -1*time.Hour) // Already expired
		token, err := expiredManager.GenerateAccessToken(userID, "user2")
		require.NoError(t, err)

		_, err = expiredManager.ValidateAccessToken(token)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "token expired")
	})

	t.Run("invalid signature", func(t *testing.T) {
		token, err := manager.GenerateAccessToken(userID, "user3")
		require.NoError(t, err)

		// Tamper with signature
		parts := strings.Split(token, ".")
		parts[2] = "invalid-signature"
		tamperedToken := strings.Join(parts, ".")

		_, err = manager.ValidateAccessToken(tamperedToken)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid token signature")
	})

	t.Run("malformed token", func(t *testing.T) {
		_, err := manager.ValidateAccessToken("not.a.valid.token")
		require.Error(t, err)
	})

	t.Run("empty token", func(t *testing.T) {
		_, err := manager.ValidateAccessToken("")
		require.Error(t, err)
	})

	t.Run("wrong secret", func(t *testing.T) {
		manager1 := NewTokenManager("secret1", 15*time.Minute)
		manager2 := NewTokenManager("secret2", 15*time.Minute)

		token, err := manager1.GenerateAccessToken(userID, "user4")
		require.NoError(t, err)

		// Try to validate with different secret
		_, err = manager2.ValidateAccessToken(token)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid token signature")
	})
}

func TestJWTManager_GenerateRefreshToken(t *testing.T) {
	t.Parallel()

	manager := NewTokenManager("test-secret-key", 15*time.Minute)

	t.Run("generates unique tokens", func(t *testing.T) {
		token1, err := manager.GenerateRefreshToken()
		require.NoError(t, err)
		assert.NotEmpty(t, token1)
		assert.Len(t, token1, 64) // 32 bytes = 64 hex chars

		token2, err := manager.GenerateRefreshToken()
		require.NoError(t, err)
		assert.NotEmpty(t, token2)

		// Tokens should be different
		assert.NotEqual(t, token1, token2)
	})

	t.Run("generates cryptographically random tokens", func(t *testing.T) {
		tokens := make(map[string]bool)
		for range 100 {
			token, err := manager.GenerateRefreshToken()
			require.NoError(t, err)

			// Should not have duplicates
			assert.False(t, tokens[token], "duplicate token generated")
			tokens[token] = true
		}
	})
}

func TestJWTManager_HashRefreshToken(t *testing.T) {
	t.Parallel()

	manager := NewTokenManager("test-secret-key", 15*time.Minute)

	t.Run("hashes consistently", func(t *testing.T) {
		token := "test-refresh-token"
		hash1 := manager.HashRefreshToken(token)
		hash2 := manager.HashRefreshToken(token)

		assert.Equal(t, hash1, hash2)
		assert.NotEqual(t, token, hash1) // Should be hashed
		assert.Len(t, hash1, 64)         // SHA-256 = 64 hex chars
	})

	t.Run("different tokens have different hashes", func(t *testing.T) {
		token1 := "token1"
		token2 := "token2"

		hash1 := manager.HashRefreshToken(token1)
		hash2 := manager.HashRefreshToken(token2)

		assert.NotEqual(t, hash1, hash2)
	})

	t.Run("hash format is hex", func(t *testing.T) {
		token := "test-token"
		hash := manager.HashRefreshToken(token)

		// Should be valid hex
		_, err := hex.DecodeString(hash)
		require.NoError(t, err)
	})
}

func TestJWTManager_ExtractClaims(t *testing.T) {
	t.Parallel()

	manager := NewTokenManager("test-secret-key", 15*time.Minute)
	userID := uuid.Must(uuid.NewV7())

	t.Run("extracts claims from valid token", func(t *testing.T) {
		token, err := manager.GenerateAccessToken(userID, "testuser")
		require.NoError(t, err)

		claims, err := manager.ExtractClaims(token)
		require.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, "testuser", claims.Username)
	})

	t.Run("extracts claims from expired token", func(t *testing.T) {
		expiredManager := NewTokenManager("test-secret-key", -1*time.Hour)
		token, err := expiredManager.GenerateAccessToken(userID, "expireduser")
		require.NoError(t, err)

		// ExtractClaims should work even for expired tokens
		claims, err := manager.ExtractClaims(token)
		require.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
	})

	t.Run("fails on malformed token", func(t *testing.T) {
		_, err := manager.ExtractClaims("not.a.token")
		require.Error(t, err)
	})
}

func TestJWTManager_TokenExpiry(t *testing.T) {
	t.Parallel()

	userID := uuid.Must(uuid.NewV7())

	t.Run("short expiry", func(t *testing.T) {
		manager := NewTokenManager("test-secret-key", 1*time.Second)
		token, err := manager.GenerateAccessToken(userID, "user")
		require.NoError(t, err)

		// Should be valid immediately
		_, err = manager.ValidateAccessToken(token)
		require.NoError(t, err)

		// Wait for expiry
		time.Sleep(2 * time.Second)

		// Should be expired now
		_, err = manager.ValidateAccessToken(token)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "token expired")
	})

	t.Run("long expiry", func(t *testing.T) {
		manager := NewTokenManager("test-secret-key", 24*time.Hour)
		token, err := manager.GenerateAccessToken(userID, "user")
		require.NoError(t, err)

		claims, err := manager.ValidateAccessToken(token)
		require.NoError(t, err)

		// Expiry should be ~24 hours from now (timestamps are in milliseconds)
		expectedExpiry := time.Now().Add(24 * time.Hour).UnixMilli()
		assert.InDelta(t, expectedExpiry, claims.ExpiresAt, 5000) // Within 5 seconds (5000ms)
	})
}

func TestJWTManager_ClaimsFields(t *testing.T) {
	t.Parallel()

	manager := NewTokenManager("test-secret-key", 15*time.Minute)
	userID := uuid.Must(uuid.NewV7())
	username := "testuser"

	token, err := manager.GenerateAccessToken(userID, username)
	require.NoError(t, err)

	claims, err := manager.ValidateAccessToken(token)
	require.NoError(t, err)

	// Check all fields are populated
	assert.NotEqual(t, uuid.Nil, claims.UserID)
	assert.NotEmpty(t, claims.Username)
	assert.Greater(t, claims.IssuedAt, int64(0))
	assert.Greater(t, claims.ExpiresAt, claims.IssuedAt)

	// IssuedAt should be recent (timestamps are in milliseconds)
	now := time.Now().UnixMilli()
	assert.InDelta(t, now, claims.IssuedAt, 5000) // Within 5 seconds (5000ms)
}
