package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/url"

	"github.com/go-faster/jx"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/google/uuid"
	"log/slog"

	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/service/mfa"
)

// MFAHandler handles MFA-related operations
type MFAHandler struct {
	totpService        *mfa.TOTPService
	backupCodesService *mfa.BackupCodesService
	mfaManager         *mfa.MFAManager
	webauthnService    *mfa.WebAuthnService
	logger             *slog.Logger
}

// NewMFAHandler creates a new MFA handler
func NewMFAHandler(
	totpService *mfa.TOTPService,
	backupCodesService *mfa.BackupCodesService,
	mfaManager *mfa.MFAManager,
	webauthnService *mfa.WebAuthnService,
	logger *slog.Logger,
) *MFAHandler {
	return &MFAHandler{
		totpService:        totpService,
		backupCodesService: backupCodesService,
		mfaManager:         mfaManager,
		webauthnService:    webauthnService,
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
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
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
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
		return &ogen.Error{
			Code:    500,
			Message: "Failed to generate TOTP secret",
		}, nil
	}

	// Parse the URL
	parsedURL, err := url.Parse(setup.URL)
	if err != nil {
		h.logger.Error("failed to parse TOTP URL",
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
		return &ogen.Error{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}

	return &ogen.TOTPSetup{
		Secret: ogen.NewOptString(setup.Secret),
		QrCode: setup.QRCode,
		URL:    ogen.NewOptURI(*parsedURL),
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
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
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
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
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
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
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
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
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
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
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
			slog.String("user_id", userID.String()),
			slog.Any("error", err))

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
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
		return &ogen.Error{
			Code:    500,
			Message: "Failed to disable MFA requirement",
		}, nil
	}

	return &ogen.DisableMFAOK{
		Success: ogen.NewOptBool(true),
	}, nil
}

// BeginWebAuthnRegistration starts the WebAuthn registration ceremony.
func (h *MFAHandler) BeginWebAuthnRegistration(ctx context.Context, req ogen.OptBeginWebAuthnRegistrationReq) (ogen.BeginWebAuthnRegistrationRes, error) {
	if h.webauthnService == nil {
		return &ogen.Error{
			Code:    501,
			Message: "WebAuthn not configured",
		}, nil
	}

	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return &ogen.Error{
			Code:    401,
			Message: "Authentication required",
		}, nil
	}

	username, err := GetUsername(ctx)
	if err != nil {
		username = userID.String()
	}

	options, err := h.webauthnService.BeginRegistration(ctx, userID, username, username)
	if err != nil {
		h.logger.Error("failed to begin WebAuthn registration",
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
		return &ogen.Error{
			Code:    500,
			Message: "Failed to begin WebAuthn registration",
		}, nil
	}

	optionsMap, err := structToJxRawMap(options)
	if err != nil {
		h.logger.Error("failed to serialize WebAuthn options",
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
		return &ogen.Error{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}

	resp := &ogen.WebAuthnBeginRegistrationResponse{
		Options: ogen.OptWebAuthnBeginRegistrationResponseOptions{
			Value: ogen.WebAuthnBeginRegistrationResponseOptions(optionsMap),
			Set:   true,
		},
	}

	if req.Set && req.Value.CredentialName.Set {
		resp.CredentialName = req.Value.CredentialName
	}

	return resp, nil
}

// FinishWebAuthnRegistration completes the WebAuthn registration ceremony.
func (h *MFAHandler) FinishWebAuthnRegistration(ctx context.Context, req *ogen.WebAuthnFinishRegistrationRequest) (ogen.FinishWebAuthnRegistrationRes, error) {
	if h.webauthnService == nil {
		return (*ogen.FinishWebAuthnRegistrationBadRequest)(&ogen.Error{
			Code:    501,
			Message: "WebAuthn not configured",
		}), nil
	}

	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return (*ogen.FinishWebAuthnRegistrationUnauthorized)(&ogen.Error{
			Code:    401,
			Message: "Authentication required",
		}), nil
	}

	username, err := GetUsername(ctx)
	if err != nil {
		username = userID.String()
	}

	// Retrieve session data from cache
	sessionData, err := h.webauthnService.GetRegistrationSession(ctx, userID)
	if err != nil {
		h.logger.Error("failed to retrieve registration session",
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
		return (*ogen.FinishWebAuthnRegistrationBadRequest)(&ogen.Error{
			Code:    400,
			Message: "Registration session expired or not found",
		}), nil
	}
	defer h.webauthnService.DeleteRegistrationSession(ctx, userID)

	// Convert map[string]jx.Raw credential back to JSON for protocol parsing
	credJSON, err := jxRawMapToJSON(req.Credential)
	if err != nil {
		h.logger.Error("failed to serialize credential response",
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
		return (*ogen.FinishWebAuthnRegistrationBadRequest)(&ogen.Error{
			Code:    400,
			Message: "Invalid credential response",
		}), nil
	}

	// Parse the credential creation response
	parsedResponse, err := protocol.ParseCredentialCreationResponseBody(bytes.NewReader(credJSON))
	if err != nil {
		h.logger.Error("failed to parse credential creation response",
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
		return (*ogen.FinishWebAuthnRegistrationBadRequest)(&ogen.Error{
			Code:    400,
			Message: "Invalid credential response format",
		}), nil
	}

	// Extract credential name
	credentialName := ""
	if req.CredentialName.Set {
		credentialName = req.CredentialName.Value
	}

	// Complete registration
	if err := h.webauthnService.FinishRegistration(ctx, userID, username, username, parsedResponse, *sessionData, credentialName); err != nil {
		h.logger.Error("failed to finish WebAuthn registration",
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
		return (*ogen.FinishWebAuthnRegistrationBadRequest)(&ogen.Error{
			Code:    400,
			Message: "Failed to complete registration",
		}), nil
	}

	return &ogen.FinishWebAuthnRegistrationOK{
		Success: ogen.NewOptBool(true),
		Message: ogen.NewOptString("WebAuthn credential registered successfully"),
	}, nil
}

// BeginWebAuthnLogin starts the WebAuthn authentication ceremony.
func (h *MFAHandler) BeginWebAuthnLogin(ctx context.Context) (ogen.BeginWebAuthnLoginRes, error) {
	if h.webauthnService == nil {
		return (*ogen.BeginWebAuthnLoginBadRequest)(&ogen.Error{
			Code:    501,
			Message: "WebAuthn not configured",
		}), nil
	}

	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return (*ogen.BeginWebAuthnLoginUnauthorized)(&ogen.Error{
			Code:    401,
			Message: "Authentication required",
		}), nil
	}

	username, err := GetUsername(ctx)
	if err != nil {
		username = userID.String()
	}

	options, err := h.webauthnService.BeginLogin(ctx, userID, username, username)
	if err != nil {
		h.logger.Error("failed to begin WebAuthn login",
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
		return (*ogen.BeginWebAuthnLoginBadRequest)(&ogen.Error{
			Code:    400,
			Message: "Failed to begin WebAuthn login",
		}), nil
	}

	optionsMap, err := structToJxRawMap(options)
	if err != nil {
		h.logger.Error("failed to serialize WebAuthn options",
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
		return (*ogen.BeginWebAuthnLoginBadRequest)(&ogen.Error{
			Code:    500,
			Message: "Internal server error",
		}), nil
	}

	return &ogen.WebAuthnBeginLoginResponse{
		Options: ogen.OptWebAuthnBeginLoginResponseOptions{
			Value: ogen.WebAuthnBeginLoginResponseOptions(optionsMap),
			Set:   true,
		},
	}, nil
}

// FinishWebAuthnLogin completes the WebAuthn authentication ceremony.
func (h *MFAHandler) FinishWebAuthnLogin(ctx context.Context, req *ogen.WebAuthnFinishLoginRequest) (ogen.FinishWebAuthnLoginRes, error) {
	if h.webauthnService == nil {
		return (*ogen.FinishWebAuthnLoginBadRequest)(&ogen.Error{
			Code:    501,
			Message: "WebAuthn not configured",
		}), nil
	}

	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return (*ogen.FinishWebAuthnLoginUnauthorized)(&ogen.Error{
			Code:    401,
			Message: "Authentication required",
		}), nil
	}

	username, err := GetUsername(ctx)
	if err != nil {
		username = userID.String()
	}

	// Retrieve session data from cache
	sessionData, err := h.webauthnService.GetLoginSession(ctx, userID)
	if err != nil {
		h.logger.Error("failed to retrieve login session",
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
		return (*ogen.FinishWebAuthnLoginBadRequest)(&ogen.Error{
			Code:    400,
			Message: "Login session expired or not found",
		}), nil
	}
	defer h.webauthnService.DeleteLoginSession(ctx, userID)

	// Convert map[string]jx.Raw credential back to JSON for protocol parsing
	credJSON, err := jxRawMapToJSON(req.Credential)
	if err != nil {
		h.logger.Error("failed to serialize credential response",
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
		return (*ogen.FinishWebAuthnLoginBadRequest)(&ogen.Error{
			Code:    400,
			Message: "Invalid credential response",
		}), nil
	}

	// Parse the credential request response
	parsedResponse, err := protocol.ParseCredentialRequestResponseBody(bytes.NewReader(credJSON))
	if err != nil {
		h.logger.Error("failed to parse credential request response",
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
		return (*ogen.FinishWebAuthnLoginBadRequest)(&ogen.Error{
			Code:    400,
			Message: "Invalid credential response format",
		}), nil
	}

	// Complete login
	if err := h.webauthnService.FinishLogin(ctx, userID, username, username, parsedResponse, *sessionData); err != nil {
		h.logger.Error("failed to finish WebAuthn login",
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
		return (*ogen.FinishWebAuthnLoginBadRequest)(&ogen.Error{
			Code:    400,
			Message: "WebAuthn authentication failed",
		}), nil
	}

	return &ogen.FinishWebAuthnLoginOK{
		Success: ogen.NewOptBool(true),
		Message: ogen.NewOptString("WebAuthn authentication successful"),
	}, nil
}

// ListWebAuthnCredentials returns all WebAuthn credentials for the user.
func (h *MFAHandler) ListWebAuthnCredentials(ctx context.Context) (ogen.ListWebAuthnCredentialsRes, error) {
	if h.webauthnService == nil {
		return &ogen.Error{
			Code:    501,
			Message: "WebAuthn not configured",
		}, nil
	}

	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return &ogen.Error{
			Code:    401,
			Message: "Authentication required",
		}, nil
	}

	creds, err := h.webauthnService.ListCredentials(ctx, userID)
	if err != nil {
		h.logger.Error("failed to list WebAuthn credentials",
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
		return &ogen.Error{
			Code:    500,
			Message: "Failed to list credentials",
		}, nil
	}

	infos := make([]ogen.WebAuthnCredentialInfo, len(creds))
	for i, cred := range creds {
		info := ogen.WebAuthnCredentialInfo{
			ID:             ogen.NewOptUUID(cred.ID),
			BackupEligible: ogen.NewOptBool(cred.BackupEligible),
			BackupState:    ogen.NewOptBool(cred.BackupState),
			CloneDetected:  ogen.NewOptBool(cred.CloneDetected),
			CreatedAt:      ogen.NewOptDateTime(cred.CreatedAt),
		}
		if cred.Name != nil {
			info.Name = ogen.NewOptString(*cred.Name)
		}
		if cred.LastUsedAt.Valid {
			info.LastUsedAt = ogen.NewOptDateTime(cred.LastUsedAt.Time)
		}
		infos[i] = info
	}

	return &ogen.WebAuthnCredentialsList{
		Credentials: infos,
		Count:       ogen.NewOptInt(len(infos)),
	}, nil
}

// DeleteWebAuthnCredential removes a WebAuthn credential.
func (h *MFAHandler) DeleteWebAuthnCredential(ctx context.Context, params ogen.DeleteWebAuthnCredentialParams) (ogen.DeleteWebAuthnCredentialRes, error) {
	if h.webauthnService == nil {
		return (*ogen.DeleteWebAuthnCredentialNotFound)(&ogen.Error{
			Code:    501,
			Message: "WebAuthn not configured",
		}), nil
	}

	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return (*ogen.DeleteWebAuthnCredentialUnauthorized)(&ogen.Error{
			Code:    401,
			Message: "Authentication required",
		}), nil
	}

	if err := h.webauthnService.DeleteCredential(ctx, userID, params.CredentialId); err != nil {
		h.logger.Error("failed to delete WebAuthn credential",
			slog.String("user_id", userID.String()),
			slog.String("credential_id", params.CredentialId.String()),
			slog.Any("error", err))
		return (*ogen.DeleteWebAuthnCredentialNotFound)(&ogen.Error{
			Code:    404,
			Message: "Credential not found",
		}), nil
	}

	return &ogen.DeleteWebAuthnCredentialNoContent{}, nil
}

// RenameWebAuthnCredential updates the name of a WebAuthn credential.
func (h *MFAHandler) RenameWebAuthnCredential(ctx context.Context, req *ogen.RenameWebAuthnCredentialReq, params ogen.RenameWebAuthnCredentialParams) (ogen.RenameWebAuthnCredentialRes, error) {
	if h.webauthnService == nil {
		return (*ogen.RenameWebAuthnCredentialNotFound)(&ogen.Error{
			Code:    501,
			Message: "WebAuthn not configured",
		}), nil
	}

	_, err := GetUserIDFromContext(ctx)
	if err != nil {
		return (*ogen.RenameWebAuthnCredentialUnauthorized)(&ogen.Error{
			Code:    401,
			Message: "Authentication required",
		}), nil
	}

	if err := h.webauthnService.RenameCredential(ctx, params.CredentialId, req.Name); err != nil {
		h.logger.Error("failed to rename WebAuthn credential",
			slog.String("credential_id", params.CredentialId.String()),
			slog.Any("error", err))
		return (*ogen.RenameWebAuthnCredentialNotFound)(&ogen.Error{
			Code:    404,
			Message: "Credential not found",
		}), nil
	}

	return &ogen.RenameWebAuthnCredentialOK{
		Success: ogen.NewOptBool(true),
	}, nil
}

// structToJxRawMap converts a Go struct to a map[string]jx.Raw via JSON serialization.
func structToJxRawMap(v any) (map[string]jx.Raw, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var rawMap map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return nil, err
	}

	result := make(map[string]jx.Raw, len(rawMap))
	for k, v := range rawMap {
		result[k] = jx.Raw(v)
	}
	return result, nil
}

// jxRawMapToJSON converts a map[string]jx.Raw back to JSON bytes.
func jxRawMapToJSON(m map[string]jx.Raw) ([]byte, error) {
	regular := make(map[string]json.RawMessage, len(m))
	for k, v := range m {
		regular[k] = json.RawMessage(v)
	}
	return json.Marshal(regular)
}

// GetUserIDFromContext extracts user ID from request context
func GetUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	return GetUserID(ctx)
}
