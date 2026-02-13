package mfa

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/google/uuid"

	db "github.com/lusoris/revenge/internal/infra/database/db"
)

var (
	ErrMFANotEnabled   = errors.New("mfa not enabled for user")
	ErrMFAAlreadySetup = errors.New("mfa already setup for user")
	ErrNoMFAMethod     = errors.New("user has no mfa methods configured")
)

// MFAManager coordinates all MFA methods (TOTP, WebAuthn, backup codes).
type MFAManager struct {
	queries     *db.Queries
	totp        *TOTPService
	webauthn    *WebAuthnService
	backupCodes *BackupCodesService
	logger      *slog.Logger
}

// NewMFAManager creates a new MFA manager.
func NewMFAManager(
	queries *db.Queries,
	totp *TOTPService,
	webauthn *WebAuthnService,
	backupCodes *BackupCodesService,
	logger *slog.Logger,
) *MFAManager {
	return &MFAManager{
		queries:     queries,
		totp:        totp,
		webauthn:    webauthn,
		backupCodes: backupCodes,
		logger:      logger,
	}
}

// MFAStatus represents the current MFA configuration for a user.
type MFAStatus struct {
	UserID                uuid.UUID `json:"user_id"`
	HasTOTP               bool      `json:"has_totp"`
	WebAuthnCount         int64     `json:"webauthn_count"`
	UnusedBackupCodes     int64     `json:"unused_backup_codes"`
	RequireMFA            bool      `json:"require_mfa"`
	RememberDeviceEnabled bool      `json:"remember_device_enabled"`
}

// GetStatus returns the current MFA status for a user.
func (m *MFAManager) GetStatus(ctx context.Context, userID uuid.UUID) (*MFAStatus, error) {
	status, err := m.queries.GetUserMFAStatus(ctx, userID)
	if err != nil {
		// If no status record exists, return default status
		return &MFAStatus{
			UserID:                userID,
			HasTOTP:               false,
			WebAuthnCount:         0,
			UnusedBackupCodes:     0,
			RequireMFA:            false,
			RememberDeviceEnabled: false,
		}, nil
	}

	// Get remember device setting from MFA settings
	rememberDeviceEnabled := false
	settings, err := m.queries.GetUserMFASettings(ctx, userID)
	if err == nil {
		rememberDeviceEnabled = settings.RememberDeviceEnabled
	}
	// If settings don't exist, rememberDeviceEnabled stays false (default)

	return &MFAStatus{
		UserID:                userID,
		HasTOTP:               status.HasTotp,
		WebAuthnCount:         status.WebauthnCount,
		UnusedBackupCodes:     status.UnusedBackupCodes,
		RequireMFA:            status.RequireMfa,
		RememberDeviceEnabled: rememberDeviceEnabled,
	}, nil
}

// HasAnyMethod checks if a user has any MFA method configured.
func (m *MFAManager) HasAnyMethod(ctx context.Context, userID uuid.UUID) (bool, error) {
	hasAny, err := m.queries.HasAnyMFAMethod(ctx, userID)
	if err != nil {
		return false, err
	}
	if hasAny == nil {
		return false, nil
	}
	return *hasAny, nil
}

// RequiresMFA checks if MFA is required for a user.
func (m *MFAManager) RequiresMFA(ctx context.Context, userID uuid.UUID) (bool, error) {
	status, err := m.GetStatus(ctx, userID)
	if err != nil {
		return false, err
	}
	return status.RequireMFA, nil
}

// EnableMFA enables MFA requirement for a user (requires at least one method to be set up).
func (m *MFAManager) EnableMFA(ctx context.Context, userID uuid.UUID) error {
	hasAny, err := m.HasAnyMethod(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check mfa methods: %w", err)
	}

	if !hasAny {
		return ErrNoMFAMethod
	}

	// Get or create MFA settings
	settings, err := m.queries.GetUserMFASettings(ctx, userID)
	if err != nil {
		// Create settings if they don't exist
		settings, err = m.queries.CreateUserMFASettings(ctx, db.CreateUserMFASettingsParams{
			UserID:                     userID,
			TotpEnabled:                false,
			WebauthnEnabled:            false,
			BackupCodesGenerated:       false,
			RequireMfa:                 false,
			RememberDeviceEnabled:      false,
			RememberDeviceDurationDays: 30,
			TrustedDevices:             []byte("[]"),
		})
		if err != nil {
			return fmt.Errorf("failed to create mfa settings: %w", err)
		}
	}

	// Enable MFA requirement
	err = m.queries.UpdateMFASettingsRequireMFA(ctx, db.UpdateMFASettingsRequireMFAParams{
		UserID:     userID,
		RequireMfa: true,
	})
	if err != nil {
		return fmt.Errorf("failed to enable mfa: %w", err)
	}

	m.logger.Info("mfa enabled",
		slog.String("user_id", userID.String()),
		slog.Bool("has_totp", settings.TotpEnabled),
		slog.Bool("has_webauthn", settings.WebauthnEnabled))

	return nil
}

// DisableMFA disables MFA requirement for a user (does not remove methods).
func (m *MFAManager) DisableMFA(ctx context.Context, userID uuid.UUID) error {
	err := m.queries.UpdateMFASettingsRequireMFA(ctx, db.UpdateMFASettingsRequireMFAParams{
		UserID:     userID,
		RequireMfa: false,
	})
	if err != nil {
		return fmt.Errorf("failed to disable mfa: %w", err)
	}

	m.logger.Info("mfa disabled", slog.String("user_id", userID.String()))
	return nil
}

// SetRememberDevice enables or disables the remember device feature for a user.
// When enabled, users can skip MFA verification on trusted devices for the specified duration.
func (m *MFAManager) SetRememberDevice(ctx context.Context, userID uuid.UUID, enabled bool, durationDays int32) error {
	// Ensure MFA settings exist first
	_, err := m.queries.GetUserMFASettings(ctx, userID)
	if err != nil {
		// Create settings if they don't exist
		_, err = m.queries.CreateUserMFASettings(ctx, db.CreateUserMFASettingsParams{
			UserID:                     userID,
			TotpEnabled:                false,
			WebauthnEnabled:            false,
			BackupCodesGenerated:       false,
			RequireMfa:                 false,
			RememberDeviceEnabled:      enabled,
			RememberDeviceDurationDays: durationDays,
			TrustedDevices:             []byte("[]"),
		})
		if err != nil {
			return fmt.Errorf("failed to create mfa settings: %w", err)
		}
		m.logger.Info("mfa settings created with remember device",
			slog.String("user_id", userID.String()),
			slog.Bool("enabled", enabled),
			slog.Any("duration_days", durationDays))
		return nil
	}

	// Update existing settings
	err = m.queries.UpdateMFASettingsRememberDevice(ctx, db.UpdateMFASettingsRememberDeviceParams{
		UserID:                     userID,
		RememberDeviceEnabled:      enabled,
		RememberDeviceDurationDays: durationDays,
	})
	if err != nil {
		return fmt.Errorf("failed to update remember device setting: %w", err)
	}

	m.logger.Info("remember device setting updated",
		slog.String("user_id", userID.String()),
		slog.Bool("enabled", enabled),
		slog.Any("duration_days", durationDays))
	return nil
}

// GetRememberDeviceSettings returns the remember device settings for a user.
func (m *MFAManager) GetRememberDeviceSettings(ctx context.Context, userID uuid.UUID) (enabled bool, durationDays int32, err error) {
	settings, err := m.queries.GetUserMFASettings(ctx, userID)
	if err != nil {
		// No settings exist, return defaults
		return false, 30, nil
	}
	return settings.RememberDeviceEnabled, settings.RememberDeviceDurationDays, nil
}

// VerifyMethod represents which MFA method was used for verification.
type VerifyMethod string

const (
	VerifyMethodTOTP       VerifyMethod = "totp"
	VerifyMethodWebAuthn   VerifyMethod = "webauthn"
	VerifyMethodBackupCode VerifyMethod = "backup_code"
)

// VerificationResult contains the result of an MFA verification attempt.
type VerificationResult struct {
	Success bool         `json:"success"`
	Method  VerifyMethod `json:"method"`
	UserID  uuid.UUID    `json:"user_id"`
}

// VerifyTOTP verifies a TOTP code for a user.
func (m *MFAManager) VerifyTOTP(ctx context.Context, userID uuid.UUID, code string) (*VerificationResult, error) {
	valid, err := m.totp.VerifyCode(ctx, userID, code)
	if err != nil {
		return nil, err
	}

	return &VerificationResult{
		Success: valid,
		Method:  VerifyMethodTOTP,
		UserID:  userID,
	}, nil
}

// VerifyBackupCode verifies a backup code for a user.
func (m *MFAManager) VerifyBackupCode(ctx context.Context, userID uuid.UUID, code string, clientIP string) (*VerificationResult, error) {
	valid, err := m.backupCodes.VerifyCode(ctx, userID, code, clientIP)
	if err != nil {
		return nil, err
	}

	return &VerificationResult{
		Success: valid,
		Method:  VerifyMethodBackupCode,
		UserID:  userID,
	}, nil
}

// VerifyWebAuthn verifies a WebAuthn assertion for a user.
// This retrieves the cached login session (from BeginLogin) and validates the
// client's assertion response. The session is cleaned up regardless of outcome.
func (m *MFAManager) VerifyWebAuthn(
	ctx context.Context,
	userID uuid.UUID,
	username string,
	assertion *protocol.ParsedCredentialAssertionData,
) (*VerificationResult, error) {
	if m.webauthn == nil {
		return nil, errors.New("webauthn service not configured")
	}
	if assertion == nil {
		return nil, errors.New("webauthn assertion data is required")
	}

	// Retrieve the login session stored by BeginLogin
	sessionData, err := m.webauthn.GetLoginSession(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("login session expired or not found: %w", err)
	}
	defer m.webauthn.DeleteLoginSession(ctx, userID)

	// Complete the WebAuthn ceremony
	if err := m.webauthn.FinishLogin(ctx, userID, username, username, assertion, *sessionData); err != nil {
		return &VerificationResult{
			Success: false,
			Method:  VerifyMethodWebAuthn,
			UserID:  userID,
		}, nil
	}

	return &VerificationResult{
		Success: true,
		Method:  VerifyMethodWebAuthn,
		UserID:  userID,
	}, nil
}

// RemoveAllMethods removes all MFA methods for a user (use with caution!).
func (m *MFAManager) RemoveAllMethods(ctx context.Context, userID uuid.UUID) error {
	// Delete TOTP
	err := m.totp.DeleteTOTP(ctx, userID)
	if err != nil {
		m.logger.Warn("failed to delete totp", slog.Any("error", err))
	}

	// Delete all WebAuthn credentials
	creds, err := m.webauthn.ListCredentials(ctx, userID)
	if err == nil {
		for _, cred := range creds {
			err = m.webauthn.DeleteCredential(ctx, userID, cred.ID)
			if err != nil {
				m.logger.Warn("failed to delete webauthn credential",
					slog.String("credential_id", cred.ID.String()),
					slog.Any("error", err))
			}
		}
	}

	// Delete all backup codes
	err = m.backupCodes.DeleteAllCodes(ctx, userID)
	if err != nil {
		m.logger.Warn("failed to delete backup codes", slog.Any("error", err))
	}

	// Delete MFA settings
	err = m.queries.DeleteUserMFASettings(ctx, userID)
	if err != nil {
		m.logger.Warn("failed to delete mfa settings", slog.Any("error", err))
	}

	m.logger.Info("removed all mfa methods", slog.String("user_id", userID.String()))
	return nil
}
