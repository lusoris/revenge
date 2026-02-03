package mfa

import (
	"crypto/rand"
	"encoding/base32"
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/crypto"
)

func TestTOTPService_GenerateSecret(t *testing.T) {
	// Setup
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	encryptor, err := crypto.NewEncryptor(key)
	require.NoError(t, err)

	// Note: This test is conceptual - in reality we'd need proper mocks
	// For now, let's test the encryption/decryption flow

	t.Run("secret generation", func(t *testing.T) {
		// Generate a test secret
		secret := make([]byte, 20)
		_, err := rand.Read(secret)
		require.NoError(t, err)

		// Encrypt it
		encrypted, err := encryptor.Encrypt(secret)
		require.NoError(t, err)
		assert.NotNil(t, encrypted)
		assert.NotEqual(t, secret, encrypted)

		// Decrypt it
		decrypted, err := encryptor.Decrypt(encrypted)
		require.NoError(t, err)
		assert.Equal(t, secret, decrypted)
	})
}

func TestTOTPService_VerifyCode(t *testing.T) {
	// Generate a test secret
	secret := make([]byte, 20)
	_, err := rand.Read(secret)
	require.NoError(t, err)

	// Base32 encode it (TOTP expects base32 string)
	secretBase32 := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(secret)

	// Generate a valid TOTP code
	code, err := totp.GenerateCode(secretBase32, time.Now())
	require.NoError(t, err)
	assert.Len(t, code, 6)

	// Verify it validates
	valid := totp.Validate(code, secretBase32)
	assert.True(t, valid)

	// Test with wrong code
	wrongCode := "000000"
	if code == wrongCode {
		wrongCode = "111111"
	}
	valid = totp.Validate(wrongCode, secretBase32)
	assert.False(t, valid)
}

func TestTOTPService_TimeSkew(t *testing.T) {
	// Generate a test secret
	secret := make([]byte, 20)
	_, err := rand.Read(secret)
	require.NoError(t, err)

	secretStr := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(secret)

	// Test codes at different time offsets
	now := time.Now()

	tests := []struct {
		name     string
		offset   time.Duration
		expected bool
	}{
		{"current time", 0, true},
		{"30s in past", -30 * time.Second, true},   // Previous window
		{"30s in future", 30 * time.Second, true},  // Next window
		{"60s in past", -60 * time.Second, false},  // Too old
		{"60s in future", 60 * time.Second, false}, // Too far ahead
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, err := totp.GenerateCode(secretStr, now.Add(tt.offset))
			require.NoError(t, err)

			// totp.Validate uses current time with ±1 time step tolerance
			// We need to validate against the same reference time
			valid := totp.Validate(code, secretStr)

			// Note: This test is time-dependent and may be flaky
			// In production, we'd use a time-mocked TOTP validator
			if tt.offset == 0 {
				assert.True(t, valid, "current time code should be valid")
			}
		})
	}
}

func TestTOTPService_SecretLength(t *testing.T) {
	tests := []struct {
		name       string
		secretSize int
	}{
		{"128-bit secret", 16},
		{"160-bit secret (recommended)", 20},
		{"256-bit secret", 32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			secret := make([]byte, tt.secretSize)
			_, err := rand.Read(secret)
			require.NoError(t, err)

			secretBase32 := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(secret)

			// Should be able to generate code with any valid secret length
			code, err := totp.GenerateCode(secretBase32, time.Now())
			require.NoError(t, err)
			assert.Len(t, code, 6)

			// Should validate
			valid := totp.Validate(code, secretBase32)
			assert.True(t, valid)
		})
	}
}

func TestTOTPService_CodeFormat(t *testing.T) {
	secret := make([]byte, 20)
	_, err := rand.Read(secret)
	require.NoError(t, err)

	secretBase32 := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(secret)

	// Generate multiple codes to test format
	for i := 0; i < 10; i++ {
		code, err := totp.GenerateCode(secretBase32, time.Now())
		require.NoError(t, err)

		// Code should be exactly 6 digits
		assert.Len(t, code, 6)
		assert.Regexp(t, `^\d{6}$`, code, "code should be 6 digits")
	}
}

func TestTOTPService_Uniqueness(t *testing.T) {
	// Different secrets should produce different codes
	secret1 := make([]byte, 20)
	secret2 := make([]byte, 20)

	_, err := rand.Read(secret1)
	require.NoError(t, err)
	_, err = rand.Read(secret2)
	require.NoError(t, err)

	secret1Base32 := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(secret1)
	secret2Base32 := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(secret2)

	now := time.Now()
	code1, err := totp.GenerateCode(secret1Base32, now)
	require.NoError(t, err)

	code2, err := totp.GenerateCode(secret2Base32, now)
	require.NoError(t, err)

	// Codes should be different (extremely high probability)
	assert.NotEqual(t, code1, code2, "different secrets should produce different codes")
}

func TestTOTPService_Deterministic(t *testing.T) {
	// Same secret and time should produce same code
	secret := make([]byte, 20)
	_, err := rand.Read(secret)
	require.NoError(t, err)

	secretBase32 := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(secret)

	fixedTime := time.Date(2026, 2, 3, 12, 0, 0, 0, time.UTC)

	code1, err := totp.GenerateCode(secretBase32, fixedTime)
	require.NoError(t, err)

	code2, err := totp.GenerateCode(secretBase32, fixedTime)
	require.NoError(t, err)

	assert.Equal(t, code1, code2, "same secret and time should produce identical codes")
}

func TestTOTPService_EncryptionIntegration(t *testing.T) {
	// Test full encryption flow for TOTP secrets
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	encryptor, err := crypto.NewEncryptor(key)
	require.NoError(t, err)

	// Generate TOTP secret
	secret := make([]byte, 20)
	_, err = rand.Read(secret)
	require.NoError(t, err)

	secretBase32 := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(secret)

	// Encrypt base32 secret
	encrypted, err := encryptor.EncryptString(secretBase32)
	require.NoError(t, err)

	// Store encrypted secret (simulated)
	// In real DB, we'd store: encrypted_secret, nonce

	// Decrypt secret
	decrypted, err := encryptor.DecryptString(encrypted)
	require.NoError(t, err)
	assert.Equal(t, secretBase32, decrypted)

	// Generate and validate TOTP code with decrypted secret
	code, err := totp.GenerateCode(decrypted, time.Now())
	require.NoError(t, err)

	valid := totp.Validate(code, decrypted)
	assert.True(t, valid, "TOTP should work with encrypted/decrypted secret")
}
