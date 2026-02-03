package mfa

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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

// Integration tests would go here
// These would require actual database connection and service mocks

func TestMFAManager_GetStatus(t *testing.T) {
	t.Skip("Integration test - requires database")

	// This test would cover:
	// 1. Returns default status for new user
	// 2. Returns actual status after TOTP setup
	// 3. Returns actual status after WebAuthn setup
	// 4. Reflects backup codes count
	// 5. Shows RequireMFA flag correctly
}

func TestMFAManager_HasAnyMethod(t *testing.T) {
	t.Skip("Integration test - requires database")

	// This test would cover:
	// 1. Returns false for new user
	// 2. Returns true after TOTP setup
	// 3. Returns true after WebAuthn setup
	// 4. Returns true after backup codes generation
	// 5. Returns false after all methods removed
}

func TestMFAManager_EnableMFA(t *testing.T) {
	t.Skip("Integration test - requires database")

	// This test would cover:
	// 1. Fails if no methods configured (ErrNoMFAMethod)
	// 2. Succeeds if TOTP configured
	// 3. Succeeds if WebAuthn configured
	// 4. Creates MFA settings if they don't exist
	// 5. Updates existing MFA settings
}

func TestMFAManager_DisableMFA(t *testing.T) {
	t.Skip("Integration test - requires database")

	// This test would cover:
	// 1. Sets RequireMFA to false
	// 2. Does not remove MFA methods
	// 3. User can still use MFA methods
	// 4. User can re-enable MFA
}

func TestMFAManager_VerifyTOTP(t *testing.T) {
	t.Skip("Integration test - requires database and services")

	// This test would cover:
	// 1. Returns success for valid TOTP code
	// 2. Returns failure for invalid TOTP code
	// 3. Sets correct method in result
	// 4. Updates last used timestamp
}

func TestMFAManager_VerifyBackupCode(t *testing.T) {
	t.Skip("Integration test - requires database and services")

	// This test would cover:
	// 1. Returns success for valid backup code
	// 2. Returns failure for invalid backup code
	// 3. Marks code as used
	// 4. Stores client IP
	// 5. Sets correct method in result
}

func TestMFAManager_RemoveAllMethods(t *testing.T) {
	t.Skip("Integration test - requires database and services")

	// This test would cover:
	// 1. Removes TOTP secret
	// 2. Removes all WebAuthn credentials
	// 3. Removes all backup codes
	// 4. Removes MFA settings
	// 5. User has no MFA methods after removal
	// 6. Handles partial failures gracefully
}

func TestMFAManager_RequiresMFA(t *testing.T) {
	t.Skip("Integration test - requires database")

	// This test would cover:
	// 1. Returns false for new user
	// 2. Returns true after EnableMFA called
	// 3. Returns false after DisableMFA called
	// 4. Persists across sessions
}
