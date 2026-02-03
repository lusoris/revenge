package crypto

import (
	"fmt"
	"testing"

	"github.com/alexedwards/argon2id"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPasswordHasher_HashPassword(t *testing.T) {
	hasher := NewPasswordHasher()

	password := "MySecurePassword123!"
	hash, err := hasher.HashPassword(password)

	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)
	assert.Contains(t, hash, "$argon2id$", "hash should use Argon2id format")
}

func TestPasswordHasher_HashPasswordEmpty(t *testing.T) {
	hasher := NewPasswordHasher()

	hash, err := hasher.HashPassword("")

	assert.Error(t, err)
	assert.Empty(t, hash)
	assert.Contains(t, err.Error(), "password cannot be empty")
}

func TestPasswordHasher_VerifyPassword(t *testing.T) {
	hasher := NewPasswordHasher()

	password := "MySecurePassword123!"
	hash, err := hasher.HashPassword(password)
	require.NoError(t, err)

	// Correct password
	match, err := hasher.VerifyPassword(password, hash)
	require.NoError(t, err)
	assert.True(t, match)

	// Wrong password
	match, err = hasher.VerifyPassword("WrongPassword", hash)
	require.NoError(t, err)
	assert.False(t, match)
}

func TestPasswordHasher_VerifyPasswordEmpty(t *testing.T) {
	hasher := NewPasswordHasher()

	// Empty password
	match, err := hasher.VerifyPassword("", "somehash")
	assert.Error(t, err)
	assert.False(t, match)

	// Empty hash
	match, err = hasher.VerifyPassword("password", "")
	assert.Error(t, err)
	assert.False(t, match)
}

func TestPasswordHasher_CustomParams(t *testing.T) {
	params := &argon2id.Params{
		Memory:      32 * 1024, // 32 MiB
		Iterations:  2,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}

	hasher := NewPasswordHasherWithParams(params)

	password := "TestPassword"
	hash, err := hasher.HashPassword(password)
	require.NoError(t, err)

	match, err := hasher.VerifyPassword(password, hash)
	require.NoError(t, err)
	assert.True(t, match)
}

func TestGenerateSecureToken(t *testing.T) {
	token1, err := GenerateSecureToken(32)
	require.NoError(t, err)
	assert.Len(t, token1, 64, "32 bytes should produce 64 hex characters")

	token2, err := GenerateSecureToken(32)
	require.NoError(t, err)
	assert.NotEqual(t, token1, token2, "tokens should be unique")
}

func TestGenerateSecureTokenInvalidLength(t *testing.T) {
	token, err := GenerateSecureToken(0)
	assert.Error(t, err)
	assert.Empty(t, token)

	token, err = GenerateSecureToken(-1)
	assert.Error(t, err)
	assert.Empty(t, token)
}

func TestGenerateSecureTokenVariousLengths(t *testing.T) {
	tests := []struct {
		byteLength int
		hexLength  int
	}{
		{16, 32},
		{32, 64},
		{64, 128},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d_bytes", tt.byteLength), func(t *testing.T) {
			token, err := GenerateSecureToken(tt.byteLength)
			require.NoError(t, err)
			assert.Len(t, token, tt.hexLength)
		})
	}
}

func TestPasswordHasher_MultipleHashes(t *testing.T) {
	hasher := NewPasswordHasher()
	password := "TestPass123!"

	// Hash the same password twice
	hash1, err := hasher.HashPassword(password)
	require.NoError(t, err)

	hash2, err := hasher.HashPassword(password)
	require.NoError(t, err)

	// Hashes should be different (due to salt)
	assert.NotEqual(t, hash1, hash2, "same password should produce different hashes")

	// Both hashes should verify
	match1, err := hasher.VerifyPassword(password, hash1)
	require.NoError(t, err)
	assert.True(t, match1)

	match2, err := hasher.VerifyPassword(password, hash2)
	require.NoError(t, err)
	assert.True(t, match2)
}
