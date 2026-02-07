package mfa

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// SetRememberDevice / GetRememberDeviceSettings Integration Tests
// ============================================================================

func TestMFAManager_SetRememberDevice(t *testing.T) {
	t.Parallel()

	t.Run("creates settings when none exist and enables remember device", func(t *testing.T) {
		t.Parallel()
		manager, _, ctx, userID := setupMFAManager(t)

		err := manager.SetRememberDevice(ctx, userID, true, 14)
		require.NoError(t, err)

		enabled, days, err := manager.GetRememberDeviceSettings(ctx, userID)
		require.NoError(t, err)
		assert.True(t, enabled)
		assert.Equal(t, int32(14), days)
	})

	t.Run("creates settings when none exist and disables remember device", func(t *testing.T) {
		t.Parallel()
		manager, _, ctx, userID := setupMFAManager(t)

		err := manager.SetRememberDevice(ctx, userID, false, 7)
		require.NoError(t, err)

		enabled, days, err := manager.GetRememberDeviceSettings(ctx, userID)
		require.NoError(t, err)
		assert.False(t, enabled)
		assert.Equal(t, int32(7), days)
	})

	t.Run("updates existing settings", func(t *testing.T) {
		t.Parallel()
		manager, _, ctx, userID := setupMFAManager(t)

		// Create settings first by enabling MFA
		_, err := manager.totp.GenerateSecret(ctx, userID, "test@example.com")
		require.NoError(t, err)
		err = manager.totp.EnableTOTP(ctx, userID)
		require.NoError(t, err)
		err = manager.EnableMFA(ctx, userID)
		require.NoError(t, err)

		// Now update remember device setting
		err = manager.SetRememberDevice(ctx, userID, true, 60)
		require.NoError(t, err)

		enabled, days, err := manager.GetRememberDeviceSettings(ctx, userID)
		require.NoError(t, err)
		assert.True(t, enabled)
		assert.Equal(t, int32(60), days)
	})

	t.Run("toggle remember device on then off", func(t *testing.T) {
		t.Parallel()
		manager, _, ctx, userID := setupMFAManager(t)

		// Enable
		err := manager.SetRememberDevice(ctx, userID, true, 30)
		require.NoError(t, err)

		enabled, _, err := manager.GetRememberDeviceSettings(ctx, userID)
		require.NoError(t, err)
		assert.True(t, enabled)

		// Disable
		err = manager.SetRememberDevice(ctx, userID, false, 30)
		require.NoError(t, err)

		enabled, _, err = manager.GetRememberDeviceSettings(ctx, userID)
		require.NoError(t, err)
		assert.False(t, enabled)
	})

	t.Run("updates duration days", func(t *testing.T) {
		t.Parallel()
		manager, _, ctx, userID := setupMFAManager(t)

		// Set initial
		err := manager.SetRememberDevice(ctx, userID, true, 7)
		require.NoError(t, err)

		// Change duration
		err = manager.SetRememberDevice(ctx, userID, true, 90)
		require.NoError(t, err)

		enabled, days, err := manager.GetRememberDeviceSettings(ctx, userID)
		require.NoError(t, err)
		assert.True(t, enabled)
		assert.Equal(t, int32(90), days)
	})
}

func TestMFAManager_GetRememberDeviceSettings(t *testing.T) {
	t.Parallel()

	t.Run("returns defaults when no settings exist", func(t *testing.T) {
		t.Parallel()
		manager, _, ctx, userID := setupMFAManager(t)

		enabled, days, err := manager.GetRememberDeviceSettings(ctx, userID)
		require.NoError(t, err)
		assert.False(t, enabled, "should default to disabled")
		assert.Equal(t, int32(30), days, "should default to 30 days")
	})

	t.Run("returns stored settings", func(t *testing.T) {
		t.Parallel()
		manager, _, ctx, userID := setupMFAManager(t)

		err := manager.SetRememberDevice(ctx, userID, true, 45)
		require.NoError(t, err)

		enabled, days, err := manager.GetRememberDeviceSettings(ctx, userID)
		require.NoError(t, err)
		assert.True(t, enabled)
		assert.Equal(t, int32(45), days)
	})
}

// ============================================================================
// GetStatus with RememberDevice Integration Tests
// ============================================================================

func TestMFAManager_GetStatus_WithRememberDevice(t *testing.T) {
	t.Parallel()

	t.Run("status reflects remember device enabled", func(t *testing.T) {
		t.Parallel()
		manager, _, ctx, userID := setupMFAManager(t)

		// Set up TOTP so GetUserMFAStatus has a record
		_, err := manager.totp.GenerateSecret(ctx, userID, "test@example.com")
		require.NoError(t, err)
		err = manager.totp.EnableTOTP(ctx, userID)
		require.NoError(t, err)
		err = manager.EnableMFA(ctx, userID)
		require.NoError(t, err)

		// Enable remember device
		err = manager.SetRememberDevice(ctx, userID, true, 30)
		require.NoError(t, err)

		status, err := manager.GetStatus(ctx, userID)
		require.NoError(t, err)
		assert.True(t, status.RememberDeviceEnabled)
		assert.True(t, status.HasTOTP)
		assert.True(t, status.RequireMFA)
	})

	t.Run("status reflects remember device disabled", func(t *testing.T) {
		t.Parallel()
		manager, _, ctx, userID := setupMFAManager(t)

		// Set up TOTP and enable MFA
		_, err := manager.totp.GenerateSecret(ctx, userID, "test@example.com")
		require.NoError(t, err)
		err = manager.totp.EnableTOTP(ctx, userID)
		require.NoError(t, err)
		err = manager.EnableMFA(ctx, userID)
		require.NoError(t, err)

		// Set remember device to disabled
		err = manager.SetRememberDevice(ctx, userID, false, 30)
		require.NoError(t, err)

		status, err := manager.GetStatus(ctx, userID)
		require.NoError(t, err)
		assert.False(t, status.RememberDeviceEnabled)
	})
}

// ============================================================================
// DisableMFA Edge Cases
// ============================================================================

func TestMFAManager_DisableMFA_NoSettingsExist(t *testing.T) {
	t.Parallel()

	manager, _, ctx, userID := setupMFAManager(t)

	// DisableMFA on a user that has never had MFA settings created.
	// This calls UpdateMFASettingsRequireMFA which may fail if no row exists.
	err := manager.DisableMFA(ctx, userID)
	// The behavior depends on the DB query - it may error or be a no-op.
	// We just verify it doesn't panic and document the behavior.
	if err != nil {
		assert.Contains(t, err.Error(), "failed to disable mfa")
	}
}

// ============================================================================
// EnableMFA → DisableMFA → Re-enable Cycle with Remember Device
// ============================================================================

func TestMFAManager_FullLifecycleWithRememberDevice(t *testing.T) {
	t.Parallel()

	manager, _, ctx, userID := setupMFAManager(t)

	// 1. Start with no MFA
	status, err := manager.GetStatus(ctx, userID)
	require.NoError(t, err)
	assert.False(t, status.RequireMFA)
	assert.False(t, status.HasTOTP)
	assert.False(t, status.RememberDeviceEnabled)

	// 2. Set up TOTP
	secret, err := manager.totp.GenerateSecret(ctx, userID, "lifecycle@example.com")
	require.NoError(t, err)
	require.NotEmpty(t, secret.Secret)

	err = manager.totp.EnableTOTP(ctx, userID)
	require.NoError(t, err)

	// 3. Enable MFA
	err = manager.EnableMFA(ctx, userID)
	require.NoError(t, err)

	// 4. Enable remember device
	err = manager.SetRememberDevice(ctx, userID, true, 14)
	require.NoError(t, err)

	// 5. Verify full status
	status, err = manager.GetStatus(ctx, userID)
	require.NoError(t, err)
	assert.True(t, status.RequireMFA)
	assert.True(t, status.HasTOTP)
	assert.True(t, status.RememberDeviceEnabled)

	// 6. Generate and verify a TOTP code
	code, err := totp.GenerateCode(secret.Secret, time.Now())
	require.NoError(t, err)
	result, err := manager.VerifyTOTP(ctx, userID, code)
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, VerifyMethodTOTP, result.Method)

	// 7. Generate backup codes
	backupCodes, err := manager.backupCodes.GenerateCodes(ctx, userID)
	require.NoError(t, err)
	require.Len(t, backupCodes, 10)

	// 8. Verify backup code via manager
	backupResult, err := manager.VerifyBackupCode(ctx, userID, backupCodes[0], "192.168.1.100")
	require.NoError(t, err)
	assert.True(t, backupResult.Success)
	assert.Equal(t, VerifyMethodBackupCode, backupResult.Method)

	// 9. Check backup codes decreased
	status, err = manager.GetStatus(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, int64(9), status.UnusedBackupCodes)

	// 10. Disable MFA
	err = manager.DisableMFA(ctx, userID)
	require.NoError(t, err)

	status, err = manager.GetStatus(ctx, userID)
	require.NoError(t, err)
	assert.False(t, status.RequireMFA)
	// TOTP still exists
	assert.True(t, status.HasTOTP)
	// Remember device setting persists
	assert.True(t, status.RememberDeviceEnabled)

	// 11. Re-enable MFA (TOTP still exists, so this should succeed)
	err = manager.EnableMFA(ctx, userID)
	require.NoError(t, err)

	requires, err := manager.RequiresMFA(ctx, userID)
	require.NoError(t, err)
	assert.True(t, requires)

	// 12. Remove all methods
	err = manager.RemoveAllMethods(ctx, userID)
	require.NoError(t, err)

	hasAny, err := manager.HasAnyMethod(ctx, userID)
	require.NoError(t, err)
	assert.False(t, hasAny)

	// Backup codes should be gone too
	count, err := manager.backupCodes.GetRemainingCount(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

// ============================================================================
// Non-existent User Error Cases
// Each subtest gets its own manager to avoid shared DB cleanup issues.
// ============================================================================

func TestMFAManager_NonExistentUser_GetStatus(t *testing.T) {
	t.Parallel()
	manager, _, ctx, _ := setupMFAManager(t)
	fakeUserID := uuid.Must(uuid.NewV7())

	status, err := manager.GetStatus(ctx, fakeUserID)
	require.NoError(t, err)
	assert.Equal(t, fakeUserID, status.UserID)
	assert.False(t, status.HasTOTP)
	assert.False(t, status.RequireMFA)
}

func TestMFAManager_NonExistentUser_HasAnyMethod(t *testing.T) {
	t.Parallel()
	manager, _, ctx, _ := setupMFAManager(t)
	fakeUserID := uuid.Must(uuid.NewV7())

	hasAny, err := manager.HasAnyMethod(ctx, fakeUserID)
	require.NoError(t, err)
	assert.False(t, hasAny)
}

func TestMFAManager_NonExistentUser_RequiresMFA(t *testing.T) {
	t.Parallel()
	manager, _, ctx, _ := setupMFAManager(t)
	fakeUserID := uuid.Must(uuid.NewV7())

	requires, err := manager.RequiresMFA(ctx, fakeUserID)
	require.NoError(t, err)
	assert.False(t, requires)
}

func TestMFAManager_NonExistentUser_EnableMFA(t *testing.T) {
	t.Parallel()
	manager, _, ctx, _ := setupMFAManager(t)
	fakeUserID := uuid.Must(uuid.NewV7())

	err := manager.EnableMFA(ctx, fakeUserID)
	assert.ErrorIs(t, err, ErrNoMFAMethod)
}

func TestMFAManager_NonExistentUser_GetRememberDeviceSettings(t *testing.T) {
	t.Parallel()
	manager, _, ctx, _ := setupMFAManager(t)
	fakeUserID := uuid.Must(uuid.NewV7())

	enabled, days, err := manager.GetRememberDeviceSettings(ctx, fakeUserID)
	require.NoError(t, err)
	assert.False(t, enabled)
	assert.Equal(t, int32(30), days)
}

func TestMFAManager_NonExistentUser_RemoveAllMethods(t *testing.T) {
	t.Parallel()
	manager, _, ctx, _ := setupMFAManager(t)
	fakeUserID := uuid.Must(uuid.NewV7())

	// RemoveAllMethods logs warnings but does not return errors
	err := manager.RemoveAllMethods(ctx, fakeUserID)
	require.NoError(t, err)
}

func TestMFAManager_NonExistentUser_VerifyTOTP(t *testing.T) {
	t.Parallel()
	manager, _, ctx, _ := setupMFAManager(t)
	fakeUserID := uuid.Must(uuid.NewV7())

	result, err := manager.VerifyTOTP(ctx, fakeUserID, "123456")
	// Should error because no TOTP secret exists
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestMFAManager_NonExistentUser_VerifyBackupCode(t *testing.T) {
	t.Parallel()
	manager, _, ctx, _ := setupMFAManager(t)
	fakeUserID := uuid.Must(uuid.NewV7())

	result, err := manager.VerifyBackupCode(ctx, fakeUserID, "0000-0000-0000-0000", "127.0.0.1")
	// When backupCodes.VerifyCode returns an error, the manager returns (nil, err)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrNoBackupCodes)
	assert.Nil(t, result)
}

// ============================================================================
// RemoveAllMethods with MFA Settings
// ============================================================================

func TestMFAManager_RemoveAllMethods_WithSettings(t *testing.T) {
	t.Parallel()

	t.Run("removes MFA settings along with methods", func(t *testing.T) {
		t.Parallel()
		manager, queries, ctx, userID := setupMFAManager(t)

		// Setup TOTP and enable MFA (creates settings)
		_, err := manager.totp.GenerateSecret(ctx, userID, "test@example.com")
		require.NoError(t, err)
		err = manager.totp.EnableTOTP(ctx, userID)
		require.NoError(t, err)
		err = manager.EnableMFA(ctx, userID)
		require.NoError(t, err)

		// Set remember device
		err = manager.SetRememberDevice(ctx, userID, true, 30)
		require.NoError(t, err)

		// Verify settings exist
		settings, err := queries.GetUserMFASettings(ctx, userID)
		require.NoError(t, err)
		assert.True(t, settings.RequireMfa)

		// Remove all methods
		err = manager.RemoveAllMethods(ctx, userID)
		require.NoError(t, err)

		// Settings should be deleted
		_, err = queries.GetUserMFASettings(ctx, userID)
		require.Error(t, err, "settings should be deleted after RemoveAllMethods")
	})

	t.Run("removes TOTP backup codes and settings together", func(t *testing.T) {
		t.Parallel()
		manager, _, ctx, userID := setupMFAManager(t)

		// Setup everything
		_, err := manager.totp.GenerateSecret(ctx, userID, "test@example.com")
		require.NoError(t, err)
		err = manager.totp.EnableTOTP(ctx, userID)
		require.NoError(t, err)
		_, err = manager.backupCodes.GenerateCodes(ctx, userID)
		require.NoError(t, err)
		err = manager.EnableMFA(ctx, userID)
		require.NoError(t, err)

		// Verify everything exists
		hasTOTP, err := manager.totp.HasTOTP(ctx, userID)
		require.NoError(t, err)
		assert.True(t, hasTOTP)
		count, err := manager.backupCodes.GetRemainingCount(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, int64(10), count)

		// Remove all
		err = manager.RemoveAllMethods(ctx, userID)
		require.NoError(t, err)

		// Verify all gone
		hasTOTP, err = manager.totp.HasTOTP(ctx, userID)
		require.NoError(t, err)
		assert.False(t, hasTOTP)
		count, err = manager.backupCodes.GetRemainingCount(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, int64(0), count)
		hasAny, err := manager.HasAnyMethod(ctx, userID)
		require.NoError(t, err)
		assert.False(t, hasAny)
	})
}

// ============================================================================
// TOTP re-enrollment test (GenerateSecret overwrites existing)
// ============================================================================

func TestMFAManager_TOTPReEnrollment(t *testing.T) {
	t.Parallel()
	manager, _, ctx, userID := setupMFAManager(t)

	// Generate first secret
	secret1, err := manager.totp.GenerateSecret(ctx, userID, "test@example.com")
	require.NoError(t, err)
	require.NotEmpty(t, secret1.Secret)
	require.NotEmpty(t, secret1.QRCode)
	require.NotEmpty(t, secret1.URL)

	// Generate second secret (re-enrollment should overwrite)
	secret2, err := manager.totp.GenerateSecret(ctx, userID, "test@example.com")
	require.NoError(t, err)
	require.NotEmpty(t, secret2.Secret)

	// Secrets should be different
	assert.NotEqual(t, secret1.Secret, secret2.Secret, "re-enrollment should generate new secret")

	// Only the second secret should be verifiable
	code2, err := totp.GenerateCode(secret2.Secret, time.Now())
	require.NoError(t, err)
	valid, err := manager.totp.VerifyCode(ctx, userID, code2)
	require.NoError(t, err)
	assert.True(t, valid, "new secret should verify")
}

// ============================================================================
// Backup code with invalid IP (tests the IP parsing fallback path)
// ============================================================================

func TestMFAManager_VerifyBackupCode_InvalidIP(t *testing.T) {
	t.Parallel()
	manager, _, ctx, userID := setupMFAManager(t)

	codes, err := manager.backupCodes.GenerateCodes(ctx, userID)
	require.NoError(t, err)
	require.Len(t, codes, 10)

	// Use an invalid IP address to test the fallback path in VerifyCode
	result, err := manager.VerifyBackupCode(ctx, userID, codes[0], "not-a-valid-ip")
	require.NoError(t, err)
	assert.True(t, result.Success, "verification should succeed even with invalid IP")
}

// ============================================================================
// RemoveAllMethods with WebAuthn Credentials
// ============================================================================

func TestMFAManager_RemoveAllMethods_WithWebAuthn(t *testing.T) {
	t.Parallel()

	t.Run("removes WebAuthn credentials", func(t *testing.T) {
		t.Parallel()
		manager, queries, ctx, userID := setupMFAManager(t)

		// Insert fake WebAuthn credentials directly via SQL
		name1 := "YubiKey 5"
		name2 := "TouchID"
		_, err := queries.CreateWebAuthnCredential(ctx, db.CreateWebAuthnCredentialParams{
			UserID:          userID,
			CredentialID:    []byte("webauthn-cred-1-remove-test"),
			PublicKey:       []byte("fake-public-key-1"),
			AttestationType: "none",
			Transports:      []string{"usb"},
			BackupEligible:  false,
			BackupState:     false,
			UserPresent:     true,
			UserVerified:    true,
			Aaguid:          []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			Name:            &name1,
		})
		require.NoError(t, err)

		_, err = queries.CreateWebAuthnCredential(ctx, db.CreateWebAuthnCredentialParams{
			UserID:          userID,
			CredentialID:    []byte("webauthn-cred-2-remove-test"),
			PublicKey:       []byte("fake-public-key-2"),
			AttestationType: "none",
			Transports:      []string{"internal"},
			BackupEligible:  true,
			BackupState:     false,
			UserPresent:     true,
			UserVerified:    true,
			Aaguid:          []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			Name:            &name2,
		})
		require.NoError(t, err)

		// Verify WebAuthn credentials exist
		creds, err := manager.webauthn.ListCredentials(ctx, userID)
		require.NoError(t, err)
		assert.Len(t, creds, 2)

		// RemoveAllMethods should delete WebAuthn credentials
		err = manager.RemoveAllMethods(ctx, userID)
		require.NoError(t, err)

		// Verify credentials are gone
		creds, err = manager.webauthn.ListCredentials(ctx, userID)
		require.NoError(t, err)
		assert.Empty(t, creds)
	})

	t.Run("removes WebAuthn + TOTP + backup codes together", func(t *testing.T) {
		t.Parallel()
		manager, queries, ctx, userID := setupMFAManager(t)

		// Setup TOTP
		_, err := manager.totp.GenerateSecret(ctx, userID, "all-methods@example.com")
		require.NoError(t, err)
		err = manager.totp.EnableTOTP(ctx, userID)
		require.NoError(t, err)

		// Generate backup codes
		_, err = manager.backupCodes.GenerateCodes(ctx, userID)
		require.NoError(t, err)

		// Insert WebAuthn credential
		name := "Security Key"
		_, err = queries.CreateWebAuthnCredential(ctx, db.CreateWebAuthnCredentialParams{
			UserID:          userID,
			CredentialID:    []byte("webauthn-all-methods-test"),
			PublicKey:       []byte("fake-public-key"),
			AttestationType: "none",
			Transports:      []string{"usb"},
			BackupEligible:  false,
			BackupState:     false,
			UserPresent:     true,
			UserVerified:    true,
			Aaguid:          []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			Name:            &name,
		})
		require.NoError(t, err)

		// Enable MFA
		err = manager.EnableMFA(ctx, userID)
		require.NoError(t, err)

		// Verify everything exists
		status, err := manager.GetStatus(ctx, userID)
		require.NoError(t, err)
		assert.True(t, status.HasTOTP)
		assert.True(t, status.RequireMFA)
		assert.Equal(t, int64(10), status.UnusedBackupCodes)
		assert.Equal(t, int64(1), status.WebAuthnCount)

		// Remove all methods
		err = manager.RemoveAllMethods(ctx, userID)
		require.NoError(t, err)

		// Verify everything is gone
		hasTOTP, err := manager.totp.HasTOTP(ctx, userID)
		require.NoError(t, err)
		assert.False(t, hasTOTP)

		creds, err := manager.webauthn.ListCredentials(ctx, userID)
		require.NoError(t, err)
		assert.Empty(t, creds)

		count, err := manager.backupCodes.GetRemainingCount(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, int64(0), count)

		hasAny, err := manager.HasAnyMethod(ctx, userID)
		require.NoError(t, err)
		assert.False(t, hasAny)
	})
}

// ============================================================================
// HasAnyMethod with WebAuthn Credentials
// ============================================================================

func TestMFAManager_HasAnyMethod_WithWebAuthn(t *testing.T) {
	t.Parallel()
	manager, queries, ctx, userID := setupMFAManager(t)

	// No methods initially
	hasAny, err := manager.HasAnyMethod(ctx, userID)
	require.NoError(t, err)
	assert.False(t, hasAny)

	// Add WebAuthn credential (without TOTP)
	name := "Security Key"
	_, err = queries.CreateWebAuthnCredential(ctx, db.CreateWebAuthnCredentialParams{
		UserID:          userID,
		CredentialID:    []byte("has-any-webauthn-test"),
		PublicKey:       []byte("fake-public-key"),
		AttestationType: "none",
		Transports:      []string{"usb"},
		BackupEligible:  false,
		BackupState:     false,
		UserPresent:     true,
		UserVerified:    true,
		Aaguid:          []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		Name:            &name,
	})
	require.NoError(t, err)

	// HasAnyMethod should return true with WebAuthn only
	hasAny, err = manager.HasAnyMethod(ctx, userID)
	require.NoError(t, err)
	assert.True(t, hasAny)
}

// ============================================================================
// GetStatus with WebAuthn Credentials
// ============================================================================

func TestMFAManager_GetStatus_WithWebAuthn(t *testing.T) {
	t.Parallel()
	manager, queries, ctx, userID := setupMFAManager(t)

	// Add two WebAuthn credentials
	name1 := "YubiKey"
	name2 := "TouchID"
	_, err := queries.CreateWebAuthnCredential(ctx, db.CreateWebAuthnCredentialParams{
		UserID:          userID,
		CredentialID:    []byte("status-webauthn-cred-1"),
		PublicKey:       []byte("fake-public-key-1"),
		AttestationType: "none",
		Transports:      []string{"usb"},
		BackupEligible:  false,
		BackupState:     false,
		UserPresent:     true,
		UserVerified:    true,
		Aaguid:          []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		Name:            &name1,
	})
	require.NoError(t, err)

	_, err = queries.CreateWebAuthnCredential(ctx, db.CreateWebAuthnCredentialParams{
		UserID:          userID,
		CredentialID:    []byte("status-webauthn-cred-2"),
		PublicKey:       []byte("fake-public-key-2"),
		AttestationType: "none",
		Transports:      []string{"internal"},
		BackupEligible:  true,
		BackupState:     false,
		UserPresent:     true,
		UserVerified:    true,
		Aaguid:          []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		Name:            &name2,
	})
	require.NoError(t, err)

	status, err := manager.GetStatus(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, int64(2), status.WebAuthnCount)
	assert.False(t, status.HasTOTP)
	assert.False(t, status.RequireMFA)
}

// ============================================================================
// EnableMFA with WebAuthn only (no TOTP)
// ============================================================================

func TestMFAManager_EnableMFA_WithWebAuthnOnly(t *testing.T) {
	t.Parallel()
	manager, queries, ctx, userID := setupMFAManager(t)

	// Add WebAuthn credential (no TOTP)
	name := "Security Key"
	_, err := queries.CreateWebAuthnCredential(ctx, db.CreateWebAuthnCredentialParams{
		UserID:          userID,
		CredentialID:    []byte("enable-mfa-webauthn-only"),
		PublicKey:       []byte("fake-public-key"),
		AttestationType: "none",
		Transports:      []string{"usb"},
		BackupEligible:  false,
		BackupState:     false,
		UserPresent:     true,
		UserVerified:    true,
		Aaguid:          []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		Name:            &name,
	})
	require.NoError(t, err)

	// EnableMFA should succeed with WebAuthn only
	err = manager.EnableMFA(ctx, userID)
	require.NoError(t, err)

	requires, err := manager.RequiresMFA(ctx, userID)
	require.NoError(t, err)
	assert.True(t, requires)
}
