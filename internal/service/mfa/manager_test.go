package mfa

import (
	"context"
	"crypto/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/lusoris/revenge/internal/crypto"
	db "github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/testutil"
)

func TestVerifyMethod_Constants(t *testing.T) {
	// Verify the method constants are correctly defined
	assert.Equal(t, VerifyMethod("totp"), VerifyMethodTOTP)
	assert.Equal(t, VerifyMethod("webauthn"), VerifyMethodWebAuthn)
	assert.Equal(t, VerifyMethod("backup_code"), VerifyMethodBackupCode)
}

func TestMFAStatus_Structure(t *testing.T) {
	userID := uuid.New()
	status := MFAStatus{
		UserID:                userID,
		HasTOTP:               true,
		WebAuthnCount:         2,
		UnusedBackupCodes:     8,
		RequireMFA:            true,
		RememberDeviceEnabled: false,
	}

	assert.Equal(t, userID, status.UserID)
	assert.True(t, status.HasTOTP)
	assert.Equal(t, int64(2), status.WebAuthnCount)
	assert.Equal(t, int64(8), status.UnusedBackupCodes)
	assert.True(t, status.RequireMFA)
	assert.False(t, status.RememberDeviceEnabled)
}

func TestVerificationResult_Structure(t *testing.T) {
	userID := uuid.New()
	result := VerificationResult{
		Success: true,
		Method:  VerifyMethodTOTP,
		UserID:  userID,
	}

	assert.True(t, result.Success)
	assert.Equal(t, VerifyMethodTOTP, result.Method)
	assert.Equal(t, userID, result.UserID)
}

// setupMFAManager creates a test MFA manager with all services
func setupMFAManager(t *testing.T) (*MFAManager, *db.Queries, context.Context, uuid.UUID) {
	t.Helper()

	testDB := testutil.NewTestDB(t)
	queries := db.New(testDB.Pool())
	ctx := context.Background()
	logger := zaptest.NewLogger(t)

	// Create an encryption key for TOTP secrets
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	encryptor, err := crypto.NewEncryptor(key)
	require.NoError(t, err)

	// Create TOTP service
	totpService := NewTOTPService(queries, encryptor, logger, "TestApp")

	// Create BackupCodes service
	backupCodesService := NewBackupCodesService(queries, logger)

	// Create WebAuthn service (can be nil for some tests)
	webauthnService, _ := NewWebAuthnService(queries, logger, nil, "Test App", "localhost", []string{"http://localhost:3000"})

	// Create the manager
	manager := NewMFAManager(queries, totpService, webauthnService, backupCodesService, logger)

	// Create a test user
	userID := createTestUserForManager(t, queries, ctx)

	return manager, queries, ctx, userID
}

// createTestUserForManager creates a test user for MFA manager tests
func createTestUserForManager(t *testing.T, queries *db.Queries, ctx context.Context) uuid.UUID {
	t.Helper()

	isActive := true
	user, err := queries.CreateUser(ctx, db.CreateUserParams{
		Email:        "mfatest_" + uuid.New().String()[:8] + "@example.com",
		Username:     "mfatest_" + uuid.New().String()[:8],
		PasswordHash: "$argon2id$v=19$m=65536,t=3,p=4$test$test",
		IsActive:     &isActive,
	})
	require.NoError(t, err)

	return user.ID
}

// generateTOTPCodeForTest generates a valid TOTP code for testing
func generateTOTPCodeForTest(secret string) (string, error) {
	return totp.GenerateCode(secret, time.Now())
}

func TestNewMFAManager(t *testing.T) {
	t.Parallel()

	testDB := testutil.NewTestDB(t)
	queries := db.New(testDB.Pool())
	logger := zaptest.NewLogger(t)

	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	encryptor, err := crypto.NewEncryptor(key)
	require.NoError(t, err)

	totpService := NewTOTPService(queries, encryptor, logger, "Test")
	backupCodesService := NewBackupCodesService(queries, logger)
	webauthnService, _ := NewWebAuthnService(queries, logger, nil, "Test", "localhost", []string{"http://localhost"})

	manager := NewMFAManager(queries, totpService, webauthnService, backupCodesService, logger)

	assert.NotNil(t, manager)
	assert.NotNil(t, manager.queries)
	assert.NotNil(t, manager.totp)
	assert.NotNil(t, manager.webauthn)
	assert.NotNil(t, manager.backupCodes)
	assert.NotNil(t, manager.logger)
}

func TestMFAManager_GetStatus(t *testing.T) {
	t.Parallel()

	t.Run("returns default status for new user", func(t *testing.T) {
		t.Parallel()
		manager, _, ctx, userID := setupMFAManager(t)

		status, err := manager.GetStatus(ctx, userID)

		require.NoError(t, err)
		assert.Equal(t, userID, status.UserID)
		assert.False(t, status.HasTOTP)
		assert.Equal(t, int64(0), status.WebAuthnCount)
		assert.Equal(t, int64(0), status.UnusedBackupCodes)
		assert.False(t, status.RequireMFA)
		assert.False(t, status.RememberDeviceEnabled)
	})

	t.Run("returns status after TOTP setup", func(t *testing.T) {
		t.Parallel()
		manager, _, ctx, userID := setupMFAManager(t)

		// Generate and enable TOTP
		_, err := manager.totp.GenerateSecret(ctx, userID, "test@example.com")
		require.NoError(t, err)

		err = manager.totp.EnableTOTP(ctx, userID)
		require.NoError(t, err)

		status, err := manager.GetStatus(ctx, userID)
		require.NoError(t, err)
		assert.True(t, status.HasTOTP)
	})

	t.Run("returns status after backup codes generation", func(t *testing.T) {
		t.Parallel()
		manager, _, ctx, userID := setupMFAManager(t)

		// Generate backup codes
		_, err := manager.backupCodes.GenerateCodes(ctx, userID)
		require.NoError(t, err)

		status, err := manager.GetStatus(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, int64(10), status.UnusedBackupCodes)
	})
}

func TestMFAManager_HasAnyMethod(t *testing.T) {
	t.Parallel()

	t.Run("returns false for new user", func(t *testing.T) {
		manager, _, ctx, userID := setupMFAManager(t)

		hasAny, err := manager.HasAnyMethod(ctx, userID)
		require.NoError(t, err)
		assert.False(t, hasAny)
	})

	t.Run("returns true after TOTP setup", func(t *testing.T) {
		manager, _, ctx, userID := setupMFAManager(t)

		// Setup TOTP
		_, err := manager.totp.GenerateSecret(ctx, userID, "test@example.com")
		require.NoError(t, err)
		err = manager.totp.EnableTOTP(ctx, userID)
		require.NoError(t, err)

		hasAny, err := manager.HasAnyMethod(ctx, userID)
		require.NoError(t, err)
		assert.True(t, hasAny)
	})

	t.Run("returns false with only backup codes", func(t *testing.T) {
		// Backup codes alone don't count as a primary MFA method
		manager, _, ctx, userID := setupMFAManager(t)

		// Generate backup codes only
		_, err := manager.backupCodes.GenerateCodes(ctx, userID)
		require.NoError(t, err)

		// HasAnyMethod checks TOTP and WebAuthn, not backup codes
		hasAny, err := manager.HasAnyMethod(ctx, userID)
		require.NoError(t, err)
		assert.False(t, hasAny) // Backup codes are fallback, not primary MFA
	})
}

func TestMFAManager_RequiresMFA(t *testing.T) {
	t.Parallel()

	t.Run("returns false for new user", func(t *testing.T) {
		manager, _, ctx, userID := setupMFAManager(t)

		requires, err := manager.RequiresMFA(ctx, userID)
		require.NoError(t, err)
		assert.False(t, requires)
	})

	t.Run("returns true after EnableMFA called", func(t *testing.T) {
		manager, _, ctx, userID := setupMFAManager(t)

		// Setup TOTP first (required for EnableMFA)
		_, err := manager.totp.GenerateSecret(ctx, userID, "test@example.com")
		require.NoError(t, err)
		err = manager.totp.EnableTOTP(ctx, userID)
		require.NoError(t, err)

		// Enable MFA
		err = manager.EnableMFA(ctx, userID)
		require.NoError(t, err)

		requires, err := manager.RequiresMFA(ctx, userID)
		require.NoError(t, err)
		assert.True(t, requires)
	})
}

func TestMFAManager_EnableMFA(t *testing.T) {
	t.Parallel()

	t.Run("fails if no methods configured", func(t *testing.T) {
		manager, _, ctx, userID := setupMFAManager(t)

		err := manager.EnableMFA(ctx, userID)
		assert.ErrorIs(t, err, ErrNoMFAMethod)
	})

	t.Run("succeeds if TOTP configured", func(t *testing.T) {
		manager, _, ctx, userID := setupMFAManager(t)

		// Setup TOTP
		_, err := manager.totp.GenerateSecret(ctx, userID, "test@example.com")
		require.NoError(t, err)
		err = manager.totp.EnableTOTP(ctx, userID)
		require.NoError(t, err)

		// Enable MFA should succeed
		err = manager.EnableMFA(ctx, userID)
		assert.NoError(t, err)

		// Verify MFA is now required
		requires, err := manager.RequiresMFA(ctx, userID)
		require.NoError(t, err)
		assert.True(t, requires)
	})

	t.Run("fails with only backup codes", func(t *testing.T) {
		// Backup codes alone are not enough to enable MFA - need TOTP or WebAuthn
		manager, _, ctx, userID := setupMFAManager(t)

		// Generate backup codes only
		_, err := manager.backupCodes.GenerateCodes(ctx, userID)
		require.NoError(t, err)

		// Enable MFA should fail - backup codes are fallback, not primary
		err = manager.EnableMFA(ctx, userID)
		assert.ErrorIs(t, err, ErrNoMFAMethod)
	})

	t.Run("creates settings if they don't exist", func(t *testing.T) {
		manager, queries, ctx, userID := setupMFAManager(t)

		// Setup TOTP
		_, err := manager.totp.GenerateSecret(ctx, userID, "test@example.com")
		require.NoError(t, err)
		err = manager.totp.EnableTOTP(ctx, userID)
		require.NoError(t, err)

		// Enable MFA (this should create settings)
		err = manager.EnableMFA(ctx, userID)
		require.NoError(t, err)

		// Verify settings were created
		settings, err := queries.GetUserMFASettings(ctx, userID)
		require.NoError(t, err)
		assert.True(t, settings.RequireMfa)
	})
}

func TestMFAManager_DisableMFA(t *testing.T) {
	t.Parallel()

	t.Run("sets RequireMFA to false", func(t *testing.T) {
		manager, _, ctx, userID := setupMFAManager(t)

		// Setup and enable MFA first
		_, err := manager.totp.GenerateSecret(ctx, userID, "test@example.com")
		require.NoError(t, err)
		err = manager.totp.EnableTOTP(ctx, userID)
		require.NoError(t, err)
		err = manager.EnableMFA(ctx, userID)
		require.NoError(t, err)

		// Verify MFA is required
		requires, err := manager.RequiresMFA(ctx, userID)
		require.NoError(t, err)
		assert.True(t, requires)

		// Disable MFA
		err = manager.DisableMFA(ctx, userID)
		require.NoError(t, err)

		// Verify MFA is no longer required
		requires, err = manager.RequiresMFA(ctx, userID)
		require.NoError(t, err)
		assert.False(t, requires)
	})

	t.Run("does not remove MFA methods", func(t *testing.T) {
		manager, _, ctx, userID := setupMFAManager(t)

		// Setup and enable MFA
		_, err := manager.totp.GenerateSecret(ctx, userID, "test@example.com")
		require.NoError(t, err)
		err = manager.totp.EnableTOTP(ctx, userID)
		require.NoError(t, err)
		err = manager.EnableMFA(ctx, userID)
		require.NoError(t, err)

		// Disable MFA
		err = manager.DisableMFA(ctx, userID)
		require.NoError(t, err)

		// TOTP should still exist
		hasTOTP, err := manager.totp.HasTOTP(ctx, userID)
		require.NoError(t, err)
		assert.True(t, hasTOTP)
	})

	t.Run("user can re-enable MFA", func(t *testing.T) {
		manager, _, ctx, userID := setupMFAManager(t)

		// Setup and enable MFA
		_, err := manager.totp.GenerateSecret(ctx, userID, "test@example.com")
		require.NoError(t, err)
		err = manager.totp.EnableTOTP(ctx, userID)
		require.NoError(t, err)
		err = manager.EnableMFA(ctx, userID)
		require.NoError(t, err)

		// Disable MFA
		err = manager.DisableMFA(ctx, userID)
		require.NoError(t, err)

		// Re-enable MFA
		err = manager.EnableMFA(ctx, userID)
		require.NoError(t, err)

		requires, err := manager.RequiresMFA(ctx, userID)
		require.NoError(t, err)
		assert.True(t, requires)
	})
}

func TestMFAManager_VerifyTOTP(t *testing.T) {
	t.Parallel()

	manager, _, ctx, userID := setupMFAManager(t)

	// Setup TOTP
	secret, err := manager.totp.GenerateSecret(ctx, userID, "test@example.com")
	require.NoError(t, err)
	err = manager.totp.EnableTOTP(ctx, userID)
	require.NoError(t, err)

	t.Run("returns success for valid TOTP code", func(t *testing.T) {
		// Generate a valid code
		code, err := generateTOTPCodeForTest(secret.Secret)
		require.NoError(t, err)

		result, err := manager.VerifyTOTP(ctx, userID, code)
		require.NoError(t, err)
		assert.True(t, result.Success)
		assert.Equal(t, VerifyMethodTOTP, result.Method)
		assert.Equal(t, userID, result.UserID)
	})

	t.Run("returns failure for invalid TOTP code", func(t *testing.T) {
		result, err := manager.VerifyTOTP(ctx, userID, "000000")
		require.NoError(t, err)
		assert.False(t, result.Success)
		assert.Equal(t, VerifyMethodTOTP, result.Method)
	})
}

func TestMFAManager_VerifyBackupCode(t *testing.T) {
	t.Parallel()

	manager, _, ctx, userID := setupMFAManager(t)

	// Generate backup codes
	codes, err := manager.backupCodes.GenerateCodes(ctx, userID)
	require.NoError(t, err)
	require.Len(t, codes, 10)

	t.Run("returns success for valid backup code", func(t *testing.T) {
		result, err := manager.VerifyBackupCode(ctx, userID, codes[0], "127.0.0.1")
		require.NoError(t, err)
		assert.True(t, result.Success)
		assert.Equal(t, VerifyMethodBackupCode, result.Method)
		assert.Equal(t, userID, result.UserID)
	})

	t.Run("returns failure for invalid backup code", func(t *testing.T) {
		result, err := manager.VerifyBackupCode(ctx, userID, "INVALID-CODE-HERE", "127.0.0.1")
		require.NoError(t, err)
		assert.False(t, result.Success)
	})

	t.Run("returns failure for already used code", func(t *testing.T) {
		// Use the same code again
		result, err := manager.VerifyBackupCode(ctx, userID, codes[0], "127.0.0.1")
		require.NoError(t, err)
		assert.False(t, result.Success)
	})
}

func TestMFAManager_RemoveAllMethods(t *testing.T) {
	t.Parallel()

	t.Run("removes TOTP secret", func(t *testing.T) {
		manager, _, ctx, userID := setupMFAManager(t)

		// Setup TOTP
		_, err := manager.totp.GenerateSecret(ctx, userID, "test@example.com")
		require.NoError(t, err)
		err = manager.totp.EnableTOTP(ctx, userID)
		require.NoError(t, err)

		// Verify TOTP exists
		hasTOTP, err := manager.totp.HasTOTP(ctx, userID)
		require.NoError(t, err)
		assert.True(t, hasTOTP)

		// Remove all methods
		err = manager.RemoveAllMethods(ctx, userID)
		require.NoError(t, err)

		// Verify TOTP is gone
		hasTOTP, err = manager.totp.HasTOTP(ctx, userID)
		require.NoError(t, err)
		assert.False(t, hasTOTP)
	})

	t.Run("removes all backup codes", func(t *testing.T) {
		manager, _, ctx, userID := setupMFAManager(t)

		// Generate backup codes
		_, err := manager.backupCodes.GenerateCodes(ctx, userID)
		require.NoError(t, err)

		// Verify codes exist
		count, err := manager.backupCodes.GetRemainingCount(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, int64(10), count)

		// Remove all methods
		err = manager.RemoveAllMethods(ctx, userID)
		require.NoError(t, err)

		// Verify codes are gone
		count, err = manager.backupCodes.GetRemainingCount(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})

	t.Run("user has no MFA methods after removal", func(t *testing.T) {
		manager, _, ctx, userID := setupMFAManager(t)

		// Setup multiple methods
		_, err := manager.totp.GenerateSecret(ctx, userID, "test@example.com")
		require.NoError(t, err)
		err = manager.totp.EnableTOTP(ctx, userID)
		require.NoError(t, err)
		_, err = manager.backupCodes.GenerateCodes(ctx, userID)
		require.NoError(t, err)

		// Verify has methods
		hasAny, err := manager.HasAnyMethod(ctx, userID)
		require.NoError(t, err)
		assert.True(t, hasAny)

		// Remove all methods
		err = manager.RemoveAllMethods(ctx, userID)
		require.NoError(t, err)

		// Verify no methods
		hasAny, err = manager.HasAnyMethod(ctx, userID)
		require.NoError(t, err)
		assert.False(t, hasAny)
	})
}

func TestMFAManager_FullWorkflow(t *testing.T) {
	t.Parallel()

	manager, _, ctx, userID := setupMFAManager(t)

	// 1. User starts with no MFA
	hasAny, err := manager.HasAnyMethod(ctx, userID)
	require.NoError(t, err)
	assert.False(t, hasAny)

	// 2. User sets up TOTP
	secret, err := manager.totp.GenerateSecret(ctx, userID, "test@example.com")
	require.NoError(t, err)
	err = manager.totp.EnableTOTP(ctx, userID)
	require.NoError(t, err)

	// 3. User generates backup codes
	codes, err := manager.backupCodes.GenerateCodes(ctx, userID)
	require.NoError(t, err)
	require.Len(t, codes, 10)

	// 4. User enables MFA requirement
	err = manager.EnableMFA(ctx, userID)
	require.NoError(t, err)

	// 5. Verify status
	status, err := manager.GetStatus(ctx, userID)
	require.NoError(t, err)
	assert.True(t, status.HasTOTP)
	assert.Equal(t, int64(10), status.UnusedBackupCodes)
	assert.True(t, status.RequireMFA)

	// 6. User authenticates with TOTP
	totpCode, err := generateTOTPCodeForTest(secret.Secret)
	require.NoError(t, err)
	result, err := manager.VerifyTOTP(ctx, userID, totpCode)
	require.NoError(t, err)
	assert.True(t, result.Success)

	// 7. User authenticates with backup code
	result, err = manager.VerifyBackupCode(ctx, userID, codes[0], "127.0.0.1")
	require.NoError(t, err)
	assert.True(t, result.Success)

	// 8. Backup codes count decreased
	status, err = manager.GetStatus(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, int64(9), status.UnusedBackupCodes)

	// 9. User disables MFA (but keeps methods)
	err = manager.DisableMFA(ctx, userID)
	require.NoError(t, err)

	requires, err := manager.RequiresMFA(ctx, userID)
	require.NoError(t, err)
	assert.False(t, requires)

	// 10. Methods still exist
	hasAny, err = manager.HasAnyMethod(ctx, userID)
	require.NoError(t, err)
	assert.True(t, hasAny)

	// 11. User removes all MFA methods
	err = manager.RemoveAllMethods(ctx, userID)
	require.NoError(t, err)

	hasAny, err = manager.HasAnyMethod(ctx, userID)
	require.NoError(t, err)
	assert.False(t, hasAny)
}
