package api

import (
	"context"
	"net/url"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/service/mfa"
)

// MFAHandler handles MFA-related operations
type MFAHandler struct {
	totpService        *mfa.TOTPService
	backupCodesService *mfa.BackupCodesService
	mfaManager         *mfa.MFAManager
	logger             *zap.Logger
}

// NewMFAHandler creates a new MFA handler
func NewMFAHandler(
	totpService *mfa.TOTPService,
	backupCodesService *mfa.BackupCodesService,
	mfaManager *mfa.MFAManager,
	logger *zap.Logger,
) *MFAHandler {
	return &MFAHandler{
		totpService:        totpService,
		backupCodesService: backupCodesService,
		mfaManager:         mfaManager,
		logger:             logger,
	}
}

// GetMFAStatus returns the current MFA configuration status for the authenticated user
func (h *MFAHandler) GetMFAStatus(ctx context.Context) (ogen.GetMFAStatusRes, error) {
	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return &ogen.Error{
			Code:    401,
			Message: "Authentication required",
		}, nil
	}

	status, err := h.mfaManager.GetStatus(ctx, userID)
	if err != nil {
		h.logger.Error("failed to get MFA status",
			zap.String("user_id", userID.String()),
			zap.Error(err))
		return &ogen.Error{
			Code:    500,
			Message: "Failed to retrieve MFA status",
		}, nil
	}

	return &ogen.MFAStatus{
		UserID:            ogen.NewOptUUID(userID),
		HasTotp:           ogen.NewOptBool(status.HasTOTP),
		WebauthnCount:     ogen.NewOptInt(int(status.WebAuthnCount)),
		UnusedBackupCodes: ogen.NewOptInt(int(status.UnusedBackupCodes)),
		RequireMfa:        ogen.NewOptBool(status.RequireMFA),
	}, nil
}

// SetupTOTP generates TOTP secret and QR code for enrollment
func (h *MFAHandler) SetupTOTP(ctx context.Context, req *ogen.SetupTOTPReq) (ogen.SetupTOTPRes, error) {
	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return &ogen.Error{
			Code:    401,
			Message: "Authentication required",
		}, nil
	}

	accountName := req.AccountName
	setup, err := h.totpService.GenerateSecret(ctx, userID, accountName)
	if err != nil {
		h.logger.Error("failed to generate TOTP secret",
			zap.String("user_id", userID.String()),
			zap.Error(err))
		return &ogen.Error{
			Code:    500,
			Message: "Failed to generate TOTP secret",
		}, nil
	}

	// Parse the URL
	parsedURL, err := url.Parse(setup.URL)
	if err != nil {
		h.logger.Error("failed to parse TOTP URL",
			zap.String("user_id", userID.String()),
			zap.Error(err))
		return &ogen.Error{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}

	return &ogen.TOTPSetup{
		Secret:  ogen.NewOptString(setup.Secret),
		QrCode:  setup.QRCode,
		URL:     ogen.NewOptURI(*parsedURL),
	}, nil
}

// VerifyTOTP verifies TOTP code and enables TOTP
func (h *MFAHandler) VerifyTOTP(ctx context.Context, req *ogen.VerifyTOTPReq) (ogen.VerifyTOTPRes, error) {
	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return &ogen.VerifyTOTPUnauthorized{
			Code:    401,
			Message: "Authentication required",
		}, nil
	}

	valid, err := h.totpService.VerifyCode(ctx, userID, req.Code)
	if err != nil {
		h.logger.Error("failed to verify TOTP code",
			zap.String("user_id", userID.String()),
			zap.Error(err))
		return (*ogen.VerifyTOTPBadRequest)(&ogen.Error{
			Code:    500,
			Message: "Failed to verify TOTP code",
		}), nil
	}

	if !valid {
		return &ogen.VerifyTOTPBadRequest{
			Code:    400,
			Message: "Invalid TOTP code",
		}, nil
	}

	// Enable TOTP
	if err := h.totpService.EnableTOTP(ctx, userID); err != nil {
		h.logger.Error("failed to enable TOTP",
			zap.String("user_id", userID.String()),
			zap.Error(err))
		return (*ogen.VerifyTOTPBadRequest)(&ogen.Error{
			Code:    500,
			Message: "Failed to enable TOTP",
		}), nil
	}

	return &ogen.VerifyTOTPOK{
		Success: ogen.NewOptBool(true),
		Message: ogen.NewOptString("TOTP enabled successfully"),
	}, nil
}

// DisableTOTP removes TOTP from user's MFA methods
func (h *MFAHandler) DisableTOTP(ctx context.Context) (ogen.DisableTOTPRes, error) {
	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return &ogen.Error{
			Code:    401,
			Message: "Authentication required",
		}, nil
	}

	if err := h.totpService.DeleteTOTP(ctx, userID); err != nil {
		h.logger.Error("failed to disable TOTP",
			zap.String("user_id", userID.String()),
			zap.Error(err))
		return &ogen.Error{
			Code:    500,
			Message: "Failed to disable TOTP",
		}, nil
	}

	return &ogen.DisableTOTPNoContent{}, nil
}

// GenerateBackupCodes generates new set of backup codes
func (h *MFAHandler) GenerateBackupCodes(ctx context.Context) (ogen.GenerateBackupCodesRes, error) {
	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return &ogen.Error{
			Code:    401,
			Message: "Authentication required",
		}, nil
	}

	codes, err := h.backupCodesService.GenerateCodes(ctx, userID)
	if err != nil {
		h.logger.Error("failed to generate backup codes",
			zap.String("user_id", userID.String()),
			zap.Error(err))
		return &ogen.Error{
			Code:    500,
			Message: "Failed to generate backup codes",
		}, nil
	}

	return &ogen.GenerateBackupCodesOK{
		Codes: codes,
		Count: ogen.NewOptInt(len(codes)),
	}, nil
}

// RegenerateBackupCodes deletes existing codes and generates new ones
func (h *MFAHandler) RegenerateBackupCodes(ctx context.Context) (ogen.RegenerateBackupCodesRes, error) {
	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return &ogen.Error{
			Code:    401,
			Message: "Authentication required",
		}, nil
	}

	codes, err := h.backupCodesService.RegenerateCodes(ctx, userID)
	if err != nil {
		h.logger.Error("failed to regenerate backup codes",
			zap.String("user_id", userID.String()),
			zap.Error(err))
		return &ogen.Error{
			Code:    500,
			Message: "Failed to regenerate backup codes",
		}, nil
	}

	return &ogen.RegenerateBackupCodesOK{
		Codes: codes,
		Count: ogen.NewOptInt(len(codes)),
	}, nil
}

// EnableMFA requires MFA for login
func (h *MFAHandler) EnableMFA(ctx context.Context) (ogen.EnableMFARes, error) {
	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return &ogen.EnableMFAUnauthorized{
			Code:    401,
			Message: "Authentication required",
		}, nil
	}

	if err := h.mfaManager.EnableMFA(ctx, userID); err != nil {
		h.logger.Error("failed to enable MFA",
			zap.String("user_id", userID.String()),
			zap.Error(err))

		// Check if error is due to no MFA methods configured
		if err == mfa.ErrNoMFAMethod {
			return (*ogen.EnableMFABadRequest)(&ogen.Error{
				Code:    400,
				Message: "At least one MFA method must be configured before enabling MFA",
			}), nil
		}

		return (*ogen.EnableMFABadRequest)(&ogen.Error{
			Code:    500,
			Message: "Failed to enable MFA requirement",
		}), nil
	}

	return &ogen.EnableMFAOK{
		Success: ogen.NewOptBool(true),
	}, nil
}

// DisableMFA turns off MFA requirement
func (h *MFAHandler) DisableMFA(ctx context.Context) (ogen.DisableMFARes, error) {
	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return &ogen.Error{
			Code:    401,
			Message: "Authentication required",
		}, nil
	}

	if err := h.mfaManager.DisableMFA(ctx, userID); err != nil {
		h.logger.Error("failed to disable MFA",
			zap.String("user_id", userID.String()),
			zap.Error(err))
		return &ogen.Error{
			Code:    500,
			Message: "Failed to disable MFA requirement",
		}, nil
	}

	return &ogen.DisableMFAOK{
		Success: ogen.NewOptBool(true),
	}, nil
}

// GetUserIDFromContext extracts user ID from request context
func GetUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	return GetUserID(ctx)
}
