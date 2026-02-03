package crypto

import (
	"fmt"
	"testing"

	"github.com/alexedwards/argon2id"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
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

func TestPasswordHasher_VerifyBcryptPassword(t *testing.T) {
	hasher := NewPasswordHasher()

	// Test with a real bcrypt hash (password: "TestPass123!")
	// Generated with: bcrypt.GenerateFromPassword([]byte("TestPass123!"), bcrypt.DefaultCost)
	password := "TestPass123!"
	bcryptHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)

	// Should successfully verify bcrypt hash
	match, err := hasher.VerifyPassword(password, string(bcryptHash))
	require.NoError(t, err)
	assert.True(t, match, "bcrypt password should match")

	// Wrong password should not match
	match, err = hasher.VerifyPassword("WrongPassword", string(bcryptHash))
	require.NoError(t, err)
	assert.False(t, match, "wrong password should not match")
}

func TestPasswordHasher_VerifyArgon2idPassword(t *testing.T) {
	hasher := NewPasswordHasher()

	password := "TestPass123!"
	hash, err := hasher.HashPassword(password)
	require.NoError(t, err)

	// Should successfully verify argon2id hash
	match, err := hasher.VerifyPassword(password, hash)
	require.NoError(t, err)
	assert.True(t, match, "argon2id password should match")

	// Wrong password should not match
	match, err = hasher.VerifyPassword("WrongPassword", hash)
	require.NoError(t, err)
	assert.False(t, match, "wrong password should not match")
}

func TestPasswordHasher_VerifyBothFormats(t *testing.T) {
	hasher := NewPasswordHasher()
	password := "SecurePass456!"

	// Generate both formats
	bcryptHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)

	argon2idHash, err := hasher.HashPassword(password)
	require.NoError(t, err)

	tests := []struct {
		name   string
		hash   string
		format string
	}{
		{"bcrypt_2a", string(bcryptHash), "bcrypt"},
		{"argon2id", argon2idHash, "argon2id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := hasher.VerifyPassword(password, tt.hash)
			require.NoError(t, err)
			assert.True(t, match, "%s format should verify correctly", tt.format)
		})
	}
}

func TestPasswordHasher_NeedsMigration(t *testing.T) {
	hasher := NewPasswordHasher()

	tests := []struct {
		name          string
		hash          string
		needsMigration bool
	}{
		{"bcrypt_2a", "$2a$12$abc123...", true},
		{"bcrypt_2b", "$2b$12$abc123...", true},
		{"bcrypt_2y", "$2y$12$abc123...", true},
		{"argon2id", "$argon2id$v=19$m=65536,t=3,p=2$...", false},
		{"unknown", "$unknown$...", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasher.NeedsMigration(tt.hash)
			assert.Equal(t, tt.needsMigration, result)
		})
	}
}
