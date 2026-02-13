package auth

import (
	"context"
	"errors"
	"fmt"
	"net/netip"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/service/mfa"
)

var (
	// ErrMFARequired indicates that MFA verification is needed
	ErrMFARequired = errors.New("mfa verification required")
	// ErrInvalidMFACode indicates that the provided MFA code is invalid
	ErrInvalidMFACode = errors.New("invalid mfa code")
	// ErrMFANotEnabled indicates that MFA is not enabled for the user
	ErrMFANotEnabled = errors.New("mfa not enabled for user")
)

// MFAAuthenticator handles MFA verification during authentication
type MFAAuthenticator struct {
	mfaManager *mfa.MFAManager
}

// NewMFAAuthenticator creates a new MFA authenticator
func NewMFAAuthenticator(mfaManager *mfa.MFAManager) *MFAAuthenticator {
	return &MFAAuthenticator{
		mfaManager: mfaManager,
	}
}

// MFALoginResponse contains the MFA challenge information
type MFALoginResponse struct {
	RequiresMFA      bool      `json:"requires_mfa"`
	UserID           uuid.UUID `json:"user_id"`
	AvailableMethods []string  `json:"available_methods"`
	// WebAuthn options would go here when implemented
	// WebAuthnOptions *WebAuthnLoginOptions `json:"webauthn_options,omitempty"`
}

// MFAVerifyRequest contains MFA verification data
type MFAVerifyRequest struct {
	UserID            uuid.UUID                               `json:"user_id"`
	Method            string                                  `json:"method"`         // "totp", "webauthn", "backup_code"
	Code              string                                  `json:"code,omitempty"` // For TOTP or backup code
	ClientIP          *netip.Addr                             `json:"-"`              // For backup code IP tracking
	Username          string                                  `json:"-"`              // For WebAuthn user lookup
	WebAuthnAssertion *protocol.ParsedCredentialAssertionData `json:"-"`              // For WebAuthn
}

// CheckMFARequired checks if the user requires MFA verification
func (m *MFAAuthenticator) CheckMFARequired(ctx context.Context, userID uuid.UUID) (*MFALoginResponse, error) {
	// Get MFA status
	status, err := m.mfaManager.GetStatus(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get MFA status: %w", err)
	}

	// Check if MFA is required
	if !status.RequireMFA {
		return &MFALoginResponse{
			RequiresMFA: false,
			UserID:      userID,
		}, nil
	}

	// Build list of available methods
	availableMethods := []string{}
	if status.HasTOTP {
		availableMethods = append(availableMethods, "totp")
	}
	if status.WebAuthnCount > 0 {
		availableMethods = append(availableMethods, "webauthn")
	}
	if status.UnusedBackupCodes > 0 {
		availableMethods = append(availableMethods, "backup_code")
	}

	response := &MFALoginResponse{
		RequiresMFA:      true,
		UserID:           userID,
		AvailableMethods: availableMethods,
	}

	return response, nil
}

// VerifyMFA verifies the provided MFA credential
func (m *MFAAuthenticator) VerifyMFA(ctx context.Context, req MFAVerifyRequest) (*mfa.VerificationResult, error) {
	switch req.Method {
	case "totp":
		if req.Code == "" {
			return nil, errors.New("totp code is required")
		}
		return m.mfaManager.VerifyTOTP(ctx, req.UserID, req.Code)

	case "backup_code":
		if req.Code == "" {
			return nil, errors.New("backup code is required")
		}
		clientIP := ""
		if req.ClientIP != nil {
			clientIP = req.ClientIP.String()
		}
		return m.mfaManager.VerifyBackupCode(ctx, req.UserID, req.Code, clientIP)

	case "webauthn":
		if req.WebAuthnAssertion == nil {
			return nil, errors.New("webauthn assertion data is required")
		}
		return m.mfaManager.VerifyWebAuthn(ctx, req.UserID, req.Username, req.WebAuthnAssertion)

	default:
		return nil, fmt.Errorf("unsupported MFA method: %s", req.Method)
	}
}

// LoginWithMFA performs a complete login with MFA check
// This is a convenience method that combines password verification and MFA check
func (s *Service) LoginWithMFA(ctx context.Context, username, password string, ipAddress *netip.Addr, userAgent, deviceName, deviceFingerprint *string, mfaAuthenticator *MFAAuthenticator) (*LoginResponse, *MFALoginResponse, error) {
	// Step 1: Verify credentials (username + password)
	loginResp, err := s.Login(ctx, username, password, ipAddress, userAgent, deviceName, deviceFingerprint)
	if err != nil {
		return nil, nil, err
	}

	// If no MFA authenticator provided, return normal login response
	if mfaAuthenticator == nil {
		return loginResp, nil, nil
	}

	// Step 2: Check if MFA is required
	mfaCheck, err := mfaAuthenticator.CheckMFARequired(ctx, loginResp.User.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to check MFA requirement: %w", err)
	}

	// Step 3: If MFA is required, return MFA challenge instead of tokens
	if mfaCheck.RequiresMFA {
		// Don't return tokens yet - user must complete MFA first
		return nil, mfaCheck, ErrMFARequired
	}

	// No MFA required, return login response
	return loginResp, nil, nil
}

// CompleteMFALogin completes the login after MFA verification
func (s *Service) CompleteMFALogin(ctx context.Context, sessionID uuid.UUID, verificationResult *mfa.VerificationResult) error {
	if !verificationResult.Success {
		return ErrInvalidMFACode
	}

	// Mark session as MFA verified
	err := s.repo.MarkSessionMFAVerified(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to mark session as MFA verified: %w", err)
	}

	return nil
}

// SessionMFAInfo contains MFA information for a session
type SessionMFAInfo struct {
	MFAVerified   bool       `json:"mfa_verified"`
	MFAVerifiedAt *time.Time `json:"mfa_verified_at,omitempty"`
}

// GetSessionMFAInfo retrieves MFA verification status for a session
func (s *Service) GetSessionMFAInfo(ctx context.Context, sessionID uuid.UUID) (*SessionMFAInfo, error) {
	session, err := s.repo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	info := &SessionMFAInfo{
		MFAVerified: session.MfaVerified,
	}

	if session.MfaVerifiedAt.Valid {
		t := session.MfaVerifiedAt.Time
		info.MFAVerifiedAt = &t
	}

	return info, nil
}
