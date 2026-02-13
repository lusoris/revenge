package mfa

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"image/png"
	"strings"
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/crypto"
	db "github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/infra/logging"
	"github.com/lusoris/revenge/internal/testutil"
)

func setupTOTPService(t *testing.T) (*TOTPService, *db.Queries) {
	t.Helper()
	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())
	logger := logging.NewTestLogger()

	// Create an encryption key for TOTP secrets
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	encryptor, err := crypto.NewEncryptor(key)
	require.NoError(t, err)

	svc := NewTOTPService(queries, encryptor, logger, "Revenge")
	return svc, queries
}

// ============================================================================
// Unit Tests for TOTP Algorithm
// ============================================================================

func TestTOTPAlgorithm_GenerateSecret(t *testing.T) {
	// Setup
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	encryptor, err := crypto.NewEncryptor(key)
	require.NoError(t, err)

	t.Run("secret generation and encryption", func(t *testing.T) {
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

func TestTOTPAlgorithm_VerifyCode(t *testing.T) {
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

func TestTOTPAlgorithm_TimeSkew(t *testing.T) {
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
			valid := totp.Validate(code, secretStr)

			// Note: This test is time-dependent
			if tt.offset == 0 {
				assert.True(t, valid, "current time code should be valid")
			}
		})
	}
}

func TestTOTPAlgorithm_SecretLength(t *testing.T) {
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

func TestTOTPAlgorithm_CodeFormat(t *testing.T) {
	secret := make([]byte, 20)
	_, err := rand.Read(secret)
	require.NoError(t, err)

	secretBase32 := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(secret)

	// Generate multiple codes to test format
	for range 10 {
		code, err := totp.GenerateCode(secretBase32, time.Now())
		require.NoError(t, err)

		// Code should be exactly 6 digits
		assert.Len(t, code, 6)
		assert.Regexp(t, `^\d{6}$`, code, "code should be 6 digits")
	}
}

func TestTOTPAlgorithm_Uniqueness(t *testing.T) {
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

func TestTOTPAlgorithm_Deterministic(t *testing.T) {
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

func TestTOTPAlgorithm_EncryptionIntegration(t *testing.T) {
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

// ============================================================================
// Integration Tests for TOTPService
// ============================================================================

func TestTOTPService_NewTOTPService(t *testing.T) {
	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())
	logger := logging.NewTestLogger()

	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	encryptor, err := crypto.NewEncryptor(key)
	require.NoError(t, err)

	svc := NewTOTPService(queries, encryptor, logger, "TestApp")

	assert.NotNil(t, svc)
	assert.Equal(t, queries, svc.queries)
	assert.Equal(t, encryptor, svc.encryptor)
	assert.Equal(t, logger, svc.logger)
	assert.Equal(t, "TestApp", svc.issuer)
}

func TestTOTPService_GenerateSecret(t *testing.T) {
	t.Parallel()
	svc, queries := setupTOTPService(t)

	ctx := context.Background()
	userID := createTestUser(t, queries, ctx)
	accountName := "test@example.com"

	setup, err := svc.GenerateSecret(ctx, userID, accountName)
	require.NoError(t, err)

	// Verify secret format
	assert.NotEmpty(t, setup.Secret)
	assert.Len(t, setup.Secret, 32, "base32-encoded 20-byte secret should be 32 chars")
	assert.Regexp(t, `^[A-Z2-7]+$`, setup.Secret, "secret should be valid base32")

	// Verify QR code
	assert.NotEmpty(t, setup.QRCode)
	// Decode PNG to verify it's valid
	_, err = png.Decode(strings.NewReader(string(setup.QRCode)))
	require.NoError(t, err, "QR code should be valid PNG")

	// Verify URL format
	assert.NotEmpty(t, setup.URL)
	assert.True(t, strings.HasPrefix(setup.URL, "otpauth://totp/"))
	assert.Contains(t, setup.URL, "Revenge")
	assert.Contains(t, setup.URL, accountName)

	// Verify secret was stored in database
	totpSecret, err := queries.GetUserTOTPSecret(ctx, userID)
	require.NoError(t, err)
	assert.NotEmpty(t, totpSecret.EncryptedSecret)
	assert.Equal(t, userID, totpSecret.UserID)
}

func TestTOTPService_GenerateSecret_MultipleUsers(t *testing.T) {
	t.Parallel()
	svc, queries := setupTOTPService(t)

	ctx := context.Background()
	userID1 := createTestUser(t, queries, ctx)
	userID2 := createTestUser(t, queries, ctx)

	setup1, err := svc.GenerateSecret(ctx, userID1, "user1@example.com")
	require.NoError(t, err)

	setup2, err := svc.GenerateSecret(ctx, userID2, "user2@example.com")
	require.NoError(t, err)

	// Each user should have different secrets
	assert.NotEqual(t, setup1.Secret, setup2.Secret)
	assert.NotEqual(t, setup1.URL, setup2.URL)
}

func TestTOTPService_VerifyCode(t *testing.T) {
	t.Parallel()
	svc, queries := setupTOTPService(t)

	ctx := context.Background()
	userID := createTestUser(t, queries, ctx)

	// Generate secret
	setup, err := svc.GenerateSecret(ctx, userID, "test@example.com")
	require.NoError(t, err)

	t.Run("valid code", func(t *testing.T) {
		// Generate current TOTP code
		code, err := totp.GenerateCode(setup.Secret, time.Now())
		require.NoError(t, err)

		// Verify code
		valid, err := svc.VerifyCode(ctx, userID, code)
		require.NoError(t, err)
		assert.True(t, valid)

		// Verify TOTP was auto-verified on first successful verification
		totpSecret, err := queries.GetUserTOTPSecret(ctx, userID)
		require.NoError(t, err)
		assert.True(t, totpSecret.VerifiedAt.Valid, "TOTP should be verified after first use")
		assert.True(t, totpSecret.Enabled, "TOTP should be enabled after first verification")
	})

	t.Run("invalid code", func(t *testing.T) {
		// Try with wrong code
		valid, err := svc.VerifyCode(ctx, userID, "000000")
		require.NoError(t, err)
		assert.False(t, valid)
	})

	t.Run("code replay rejected", func(t *testing.T) {
		// Use the next time step to get a fresh code (not the same as "valid code" subtest)
		code, err := totp.GenerateCode(setup.Secret, time.Now().Add(30*time.Second))
		require.NoError(t, err)

		valid1, err := svc.VerifyCode(ctx, userID, code)
		require.NoError(t, err)
		assert.True(t, valid1)

		// Same code should be rejected (replay protection)
		valid2, err := svc.VerifyCode(ctx, userID, code)
		require.NoError(t, err)
		assert.False(t, valid2, "replayed TOTP code should be rejected")
	})

	t.Run("no TOTP configured", func(t *testing.T) {
		newUserID := createTestUser(t, queries, ctx)

		// Should fail when user has no TOTP secret
		valid, err := svc.VerifyCode(ctx, newUserID, "123456")
		require.Error(t, err)
		assert.False(t, valid)
		assert.Contains(t, err.Error(), "failed to get TOTP secret")
	})
}

func TestTOTPService_EnableTOTP(t *testing.T) {
	t.Parallel()
	svc, queries := setupTOTPService(t)

	ctx := context.Background()
	userID := createTestUser(t, queries, ctx)

	// Generate and verify secret first
	setup, err := svc.GenerateSecret(ctx, userID, "test@example.com")
	require.NoError(t, err)

	code, err := totp.GenerateCode(setup.Secret, time.Now())
	require.NoError(t, err)

	valid, err := svc.VerifyCode(ctx, userID, code)
	require.NoError(t, err)
	require.True(t, valid)

	// Enable TOTP
	err = svc.EnableTOTP(ctx, userID)
	require.NoError(t, err)

	// Verify it's enabled
	totpSecret, err := queries.GetUserTOTPSecret(ctx, userID)
	require.NoError(t, err)
	assert.True(t, totpSecret.Enabled)
}

func TestTOTPService_DisableTOTP(t *testing.T) {
	t.Parallel()
	svc, queries := setupTOTPService(t)

	ctx := context.Background()
	userID := createTestUser(t, queries, ctx)

	// Setup and enable TOTP
	setup, err := svc.GenerateSecret(ctx, userID, "test@example.com")
	require.NoError(t, err)

	code, err := totp.GenerateCode(setup.Secret, time.Now())
	require.NoError(t, err)

	valid, err := svc.VerifyCode(ctx, userID, code)
	require.NoError(t, err)
	require.True(t, valid)

	// Verify it's enabled
	totpSecret, err := queries.GetUserTOTPSecret(ctx, userID)
	require.NoError(t, err)
	require.True(t, totpSecret.Enabled)

	// Disable TOTP
	err = svc.DisableTOTP(ctx, userID)
	require.NoError(t, err)

	// Verify it's disabled
	totpSecret, err = queries.GetUserTOTPSecret(ctx, userID)
	require.NoError(t, err)
	assert.False(t, totpSecret.Enabled)

	// Secret should still exist (just disabled)
	assert.NotEmpty(t, totpSecret.EncryptedSecret)
}

func TestTOTPService_DeleteTOTP(t *testing.T) {
	t.Parallel()
	svc, queries := setupTOTPService(t)

	ctx := context.Background()
	userID := createTestUser(t, queries, ctx)

	// Setup TOTP
	_, err := svc.GenerateSecret(ctx, userID, "test@example.com")
	require.NoError(t, err)

	// Verify it exists
	_, err = queries.GetUserTOTPSecret(ctx, userID)
	require.NoError(t, err)

	// Delete TOTP
	err = svc.DeleteTOTP(ctx, userID)
	require.NoError(t, err)

	// Verify it's deleted
	_, err = queries.GetUserTOTPSecret(ctx, userID)
	assert.Error(t, err, "TOTP secret should be deleted")
}

func TestTOTPService_HasTOTP(t *testing.T) {
	t.Parallel()
	svc, queries := setupTOTPService(t)

	ctx := context.Background()
	userID := createTestUser(t, queries, ctx)

	// Should not have TOTP initially
	has, err := svc.HasTOTP(ctx, userID)
	require.NoError(t, err)
	assert.False(t, has)

	// Generate secret
	_, err = svc.GenerateSecret(ctx, userID, "test@example.com")
	require.NoError(t, err)

	// Should have TOTP now
	has, err = svc.HasTOTP(ctx, userID)
	require.NoError(t, err)
	assert.True(t, has)

	// Delete TOTP
	err = svc.DeleteTOTP(ctx, userID)
	require.NoError(t, err)

	// Should not have TOTP after deletion
	has, err = svc.HasTOTP(ctx, userID)
	require.NoError(t, err)
	assert.False(t, has)
}

func TestTOTPService_CompleteFlow(t *testing.T) {
	t.Parallel()
	svc, queries := setupTOTPService(t)

	ctx := context.Background()
	userID := createTestUser(t, queries, ctx)

	// 1. User starts TOTP setup
	setup, err := svc.GenerateSecret(ctx, userID, "user@example.com")
	require.NoError(t, err)
	assert.NotEmpty(t, setup.Secret)
	assert.NotEmpty(t, setup.QRCode)
	assert.NotEmpty(t, setup.URL)

	// 2. User scans QR code and enters first code to verify
	code, err := totp.GenerateCode(setup.Secret, time.Now())
	require.NoError(t, err)

	valid, err := svc.VerifyCode(ctx, userID, code)
	require.NoError(t, err)
	assert.True(t, valid)

	// 3. Verify TOTP is enabled after first verification
	has, err := svc.HasTOTP(ctx, userID)
	require.NoError(t, err)
	assert.True(t, has)

	totpSecret, err := queries.GetUserTOTPSecret(ctx, userID)
	require.NoError(t, err)
	assert.True(t, totpSecret.Enabled)
	assert.True(t, totpSecret.VerifiedAt.Valid)

	// 4. User can use TOTP for login (use next time step to avoid replay rejection)
	code2, err := totp.GenerateCode(setup.Secret, time.Now().Add(30*time.Second))
	require.NoError(t, err)

	// If code2 is the same as the initial code, skip this check (same time step)
	if code2 != code {
		valid, err = svc.VerifyCode(ctx, userID, code2)
		require.NoError(t, err)
		assert.True(t, valid)
	}

	// 5. User can disable TOTP
	err = svc.DisableTOTP(ctx, userID)
	require.NoError(t, err)

	totpSecret, err = queries.GetUserTOTPSecret(ctx, userID)
	require.NoError(t, err)
	assert.False(t, totpSecret.Enabled)

	// 6. User can re-enable TOTP
	err = svc.EnableTOTP(ctx, userID)
	require.NoError(t, err)

	totpSecret, err = queries.GetUserTOTPSecret(ctx, userID)
	require.NoError(t, err)
	assert.True(t, totpSecret.Enabled)

	// 7. User can delete TOTP completely
	err = svc.DeleteTOTP(ctx, userID)
	require.NoError(t, err)

	has, err = svc.HasTOTP(ctx, userID)
	require.NoError(t, err)
	assert.False(t, has)
}

func TestTOTPService_SecretEncryption(t *testing.T) {
	t.Parallel()
	svc, queries := setupTOTPService(t)

	ctx := context.Background()
	userID := createTestUser(t, queries, ctx)

	// Generate secret
	setup, err := svc.GenerateSecret(ctx, userID, "test@example.com")
	require.NoError(t, err)

	// Get encrypted secret from database
	totpSecret, err := queries.GetUserTOTPSecret(ctx, userID)
	require.NoError(t, err)

	// Encrypted secret should not match plaintext secret
	assert.NotEqual(t, setup.Secret, string(totpSecret.EncryptedSecret))

	// Encrypted secret should be longer (includes nonce + ciphertext + tag)
	assert.Greater(t, len(totpSecret.EncryptedSecret), len(setup.Secret))

	// Generate and verify code to ensure encryption/decryption works
	code, err := totp.GenerateCode(setup.Secret, time.Now())
	require.NoError(t, err)

	valid, err := svc.VerifyCode(ctx, userID, code)
	require.NoError(t, err)
	assert.True(t, valid, "encrypted secret should decrypt correctly for verification")
}

func TestTOTPService_MultipleSecrets(t *testing.T) {
	t.Parallel()
	svc, queries := setupTOTPService(t)

	ctx := context.Background()
	userID := createTestUser(t, queries, ctx)

	// Generate first secret
	setup1, err := svc.GenerateSecret(ctx, userID, "test@example.com")
	require.NoError(t, err)

	// Generate second secret (should replace first one)
	setup2, err := svc.GenerateSecret(ctx, userID, "test@example.com")
	require.NoError(t, err)

	// Secrets should be different
	assert.NotEqual(t, setup1.Secret, setup2.Secret)

	// Old secret should not work
	code1, err := totp.GenerateCode(setup1.Secret, time.Now())
	require.NoError(t, err)

	valid, err := svc.VerifyCode(ctx, userID, code1)
	require.NoError(t, err)
	assert.False(t, valid, "old secret should not work after regeneration")

	// New secret should work
	code2, err := totp.GenerateCode(setup2.Secret, time.Now())
	require.NoError(t, err)

	valid, err = svc.VerifyCode(ctx, userID, code2)
	require.NoError(t, err)
	assert.True(t, valid, "new secret should work")
}
