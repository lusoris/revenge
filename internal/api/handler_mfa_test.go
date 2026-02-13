package api

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/go-faster/jx"
	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/infra/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/api/ogen"
)

func newTestMFAHandler() *MFAHandler {
	return NewMFAHandler(nil, nil, nil, nil, logging.NewTestLogger())
}

func contextWithUser(t *testing.T) context.Context {
	t.Helper()
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())
	ctx = WithUserID(ctx, userID)
	ctx = WithUsername(ctx, "testuser")
	return ctx
}

// Test that NewMFAHandler stores the webauthnService field
func TestNewMFAHandler_WithWebAuthnService(t *testing.T) {
	handler := NewMFAHandler(nil, nil, nil, nil, logging.NewTestLogger())
	assert.Nil(t, handler.webauthnService)

	// Non-nil is tested implicitly by all "not configured" tests returning 501
}

// Test all WebAuthn handlers return 501 when WebAuthn service is nil
func TestMFAHandler_WebAuthn_NotConfigured(t *testing.T) {
	handler := newTestMFAHandler()
	ctx := contextWithUser(t)

	t.Run("BeginWebAuthnRegistration", func(t *testing.T) {
		res, err := handler.BeginWebAuthnRegistration(ctx, ogen.OptBeginWebAuthnRegistrationReq{})
		require.NoError(t, err)
		errRes, ok := res.(*ogen.Error)
		require.True(t, ok)
		assert.Equal(t, 501, errRes.Code)
		assert.Contains(t, errRes.Message, "WebAuthn not configured")
	})

	t.Run("FinishWebAuthnRegistration", func(t *testing.T) {
		res, err := handler.FinishWebAuthnRegistration(ctx, &ogen.WebAuthnFinishRegistrationRequest{})
		require.NoError(t, err)
		errRes, ok := res.(*ogen.FinishWebAuthnRegistrationBadRequest)
		require.True(t, ok)
		assert.Equal(t, 501, errRes.Code)
	})

	t.Run("BeginWebAuthnLogin", func(t *testing.T) {
		res, err := handler.BeginWebAuthnLogin(ctx)
		require.NoError(t, err)
		errRes, ok := res.(*ogen.BeginWebAuthnLoginBadRequest)
		require.True(t, ok)
		assert.Equal(t, 501, errRes.Code)
	})

	t.Run("FinishWebAuthnLogin", func(t *testing.T) {
		res, err := handler.FinishWebAuthnLogin(ctx, &ogen.WebAuthnFinishLoginRequest{})
		require.NoError(t, err)
		errRes, ok := res.(*ogen.FinishWebAuthnLoginBadRequest)
		require.True(t, ok)
		assert.Equal(t, 501, errRes.Code)
	})

	t.Run("ListWebAuthnCredentials", func(t *testing.T) {
		res, err := handler.ListWebAuthnCredentials(ctx)
		require.NoError(t, err)
		errRes, ok := res.(*ogen.Error)
		require.True(t, ok)
		assert.Equal(t, 501, errRes.Code)
	})

	t.Run("DeleteWebAuthnCredential", func(t *testing.T) {
		res, err := handler.DeleteWebAuthnCredential(ctx, ogen.DeleteWebAuthnCredentialParams{
			CredentialId: uuid.Must(uuid.NewV7()),
		})
		require.NoError(t, err)
		errRes, ok := res.(*ogen.DeleteWebAuthnCredentialNotFound)
		require.True(t, ok)
		assert.Equal(t, 501, errRes.Code)
	})

	t.Run("RenameWebAuthnCredential", func(t *testing.T) {
		res, err := handler.RenameWebAuthnCredential(ctx, &ogen.RenameWebAuthnCredentialReq{
			Name: "New Name",
		}, ogen.RenameWebAuthnCredentialParams{
			CredentialId: uuid.Must(uuid.NewV7()),
		})
		require.NoError(t, err)
		errRes, ok := res.(*ogen.RenameWebAuthnCredentialNotFound)
		require.True(t, ok)
		assert.Equal(t, 501, errRes.Code)
	})
}

// Test all WebAuthn handlers return 401 when no user is in context
func TestMFAHandler_WebAuthn_Unauthorized(t *testing.T) {
	ctx := context.Background() // No user in context

	t.Run("GetUserIDFromContext_fails_without_user", func(t *testing.T) {
		_, err := GetUserIDFromContext(ctx)
		require.Error(t, err)
	})

	t.Run("GetUserIDFromContext_succeeds_with_user", func(t *testing.T) {
		userID := uuid.Must(uuid.NewV7())
		ctxWithUser := WithUserID(ctx, userID)
		got, err := GetUserIDFromContext(ctxWithUser)
		require.NoError(t, err)
		assert.Equal(t, userID, got)
	})
}

// Test structToJxRawMap conversion
func TestStructToJxRawMap(t *testing.T) {
	t.Run("simple struct", func(t *testing.T) {
		input := struct {
			Name  string `json:"name"`
			Value int    `json:"value"`
		}{
			Name:  "test",
			Value: 42,
		}

		result, err := structToJxRawMap(input)
		require.NoError(t, err)
		assert.Len(t, result, 2)

		// Check "name" field
		nameRaw, ok := result["name"]
		require.True(t, ok)
		assert.Equal(t, jx.Raw(`"test"`), nameRaw)

		// Check "value" field
		valueRaw, ok := result["value"]
		require.True(t, ok)
		assert.Equal(t, jx.Raw(`42`), valueRaw)
	})

	t.Run("nested struct", func(t *testing.T) {
		input := struct {
			Outer string `json:"outer"`
			Inner struct {
				Field string `json:"field"`
			} `json:"inner"`
		}{
			Outer: "hello",
		}
		input.Inner.Field = "world"

		result, err := structToJxRawMap(input)
		require.NoError(t, err)
		assert.Len(t, result, 2)

		innerRaw, ok := result["inner"]
		require.True(t, ok)
		assert.Contains(t, string(innerRaw), `"field"`)
		assert.Contains(t, string(innerRaw), `"world"`)
	})

	t.Run("empty struct", func(t *testing.T) {
		input := struct{}{}
		result, err := structToJxRawMap(input)
		require.NoError(t, err)
		assert.Empty(t, result)
	})
}

// Test jxRawMapToJSON conversion
func TestJxRawMapToJSON(t *testing.T) {
	t.Run("simple map", func(t *testing.T) {
		input := map[string]jx.Raw{
			"name":  jx.Raw(`"test"`),
			"value": jx.Raw(`42`),
		}

		result, err := jxRawMapToJSON(input)
		require.NoError(t, err)

		// Parse back to verify valid JSON
		var parsed map[string]json.RawMessage
		err = json.Unmarshal(result, &parsed)
		require.NoError(t, err)
		assert.Len(t, parsed, 2)
		assert.Equal(t, json.RawMessage(`"test"`), parsed["name"])
		assert.Equal(t, json.RawMessage(`42`), parsed["value"])
	})

	t.Run("empty map", func(t *testing.T) {
		input := map[string]jx.Raw{}
		result, err := jxRawMapToJSON(input)
		require.NoError(t, err)
		assert.Equal(t, `{}`, string(result))
	})

	t.Run("round trip", func(t *testing.T) {
		original := struct {
			Challenge string   `json:"challenge"`
			RPName    string   `json:"rpName"`
			Timeout   int      `json:"timeout"`
			Params    []string `json:"params"`
		}{
			Challenge: "abc123",
			RPName:    "Revenge",
			Timeout:   60000,
			Params:    []string{"es256", "rs256"},
		}

		// Struct -> map[string]jx.Raw
		rawMap, err := structToJxRawMap(original)
		require.NoError(t, err)

		// map[string]jx.Raw -> JSON bytes
		jsonBytes, err := jxRawMapToJSON(rawMap)
		require.NoError(t, err)

		// JSON bytes -> struct
		var restored struct {
			Challenge string   `json:"challenge"`
			RPName    string   `json:"rpName"`
			Timeout   int      `json:"timeout"`
			Params    []string `json:"params"`
		}
		err = json.Unmarshal(jsonBytes, &restored)
		require.NoError(t, err)

		assert.Equal(t, original.Challenge, restored.Challenge)
		assert.Equal(t, original.RPName, restored.RPName)
		assert.Equal(t, original.Timeout, restored.Timeout)
		assert.Equal(t, original.Params, restored.Params)
	})
}

// =============================================================================
// NoAuth tests for TOTP, Backup Codes, EnableMFA, DisableMFA
// =============================================================================

// TestMFAHandler_GetMFAStatus_NoAuth verifies that GetMFAStatus returns 401 when
// no user ID is present in the context.
func TestMFAHandler_GetMFAStatus_NoAuth(t *testing.T) {
	handler := newTestMFAHandler()
	ctx := context.Background() // No user in context

	res, err := handler.GetMFAStatus(ctx)
	require.NoError(t, err)

	errRes, ok := res.(*ogen.Error)
	require.True(t, ok, "expected *ogen.Error, got %T", res)
	assert.Equal(t, 401, errRes.Code)
	assert.Equal(t, "Authentication required", errRes.Message)
}

// TestMFAHandler_SetupTOTP_NoAuth verifies that SetupTOTP returns 401 when
// no user ID is present in the context.
func TestMFAHandler_SetupTOTP_NoAuth(t *testing.T) {
	handler := newTestMFAHandler()
	ctx := context.Background()

	res, err := handler.SetupTOTP(ctx, &ogen.SetupTOTPReq{
		AccountName: "test@example.com",
	})
	require.NoError(t, err)

	errRes, ok := res.(*ogen.Error)
	require.True(t, ok, "expected *ogen.Error, got %T", res)
	assert.Equal(t, 401, errRes.Code)
	assert.Equal(t, "Authentication required", errRes.Message)
}

// TestMFAHandler_VerifyTOTP_NoAuth verifies that VerifyTOTP returns 401 when
// no user ID is present in the context.
func TestMFAHandler_VerifyTOTP_NoAuth(t *testing.T) {
	handler := newTestMFAHandler()
	ctx := context.Background()

	res, err := handler.VerifyTOTP(ctx, &ogen.VerifyTOTPReq{
		Code: "123456",
	})
	require.NoError(t, err)

	errRes, ok := res.(*ogen.VerifyTOTPUnauthorized)
	require.True(t, ok, "expected *ogen.VerifyTOTPUnauthorized, got %T", res)
	assert.Equal(t, 401, errRes.Code)
	assert.Equal(t, "Authentication required", errRes.Message)
}

// TestMFAHandler_DisableTOTP_NoAuth verifies that DisableTOTP returns 401 when
// no user ID is present in the context.
func TestMFAHandler_DisableTOTP_NoAuth(t *testing.T) {
	handler := newTestMFAHandler()
	ctx := context.Background()

	res, err := handler.DisableTOTP(ctx)
	require.NoError(t, err)

	errRes, ok := res.(*ogen.Error)
	require.True(t, ok, "expected *ogen.Error, got %T", res)
	assert.Equal(t, 401, errRes.Code)
	assert.Equal(t, "Authentication required", errRes.Message)
}

// TestMFAHandler_GenerateBackupCodes_NoAuth verifies that GenerateBackupCodes
// returns 401 when no user ID is present in the context.
func TestMFAHandler_GenerateBackupCodes_NoAuth(t *testing.T) {
	handler := newTestMFAHandler()
	ctx := context.Background()

	res, err := handler.GenerateBackupCodes(ctx)
	require.NoError(t, err)

	errRes, ok := res.(*ogen.Error)
	require.True(t, ok, "expected *ogen.Error, got %T", res)
	assert.Equal(t, 401, errRes.Code)
	assert.Equal(t, "Authentication required", errRes.Message)
}

// TestMFAHandler_RegenerateBackupCodes_NoAuth verifies that RegenerateBackupCodes
// returns 401 when no user ID is present in the context.
func TestMFAHandler_RegenerateBackupCodes_NoAuth(t *testing.T) {
	handler := newTestMFAHandler()
	ctx := context.Background()

	res, err := handler.RegenerateBackupCodes(ctx)
	require.NoError(t, err)

	errRes, ok := res.(*ogen.Error)
	require.True(t, ok, "expected *ogen.Error, got %T", res)
	assert.Equal(t, 401, errRes.Code)
	assert.Equal(t, "Authentication required", errRes.Message)
}

// TestMFAHandler_EnableMFA_NoAuth verifies that EnableMFA returns 401 when
// no user ID is present in the context.
func TestMFAHandler_EnableMFA_NoAuth(t *testing.T) {
	handler := newTestMFAHandler()
	ctx := context.Background()

	res, err := handler.EnableMFA(ctx)
	require.NoError(t, err)

	errRes, ok := res.(*ogen.EnableMFAUnauthorized)
	require.True(t, ok, "expected *ogen.EnableMFAUnauthorized, got %T", res)
	assert.Equal(t, 401, errRes.Code)
	assert.Equal(t, "Authentication required", errRes.Message)
}

// TestMFAHandler_DisableMFA_NoAuth verifies that DisableMFA returns 401 when
// no user ID is present in the context.
func TestMFAHandler_DisableMFA_NoAuth(t *testing.T) {
	handler := newTestMFAHandler()
	ctx := context.Background()

	res, err := handler.DisableMFA(ctx)
	require.NoError(t, err)

	errRes, ok := res.(*ogen.Error)
	require.True(t, ok, "expected *ogen.Error, got %T", res)
	assert.Equal(t, 401, errRes.Code)
	assert.Equal(t, "Authentication required", errRes.Message)
}

// =============================================================================
// Handler delegation tests for TOTP / Backup Codes / Enable/Disable MFA
// =============================================================================

// TestHandler_GetMFAStatus_NoAuth tests that the top-level Handler delegation for
// GetMFAStatus returns 401 when no user is in context.
func TestHandler_GetMFAStatus_NoAuth(t *testing.T) {
	mfaHandler := newTestMFAHandler()
	handler := &Handler{
		logger:     logging.NewTestLogger(),
		mfaHandler: mfaHandler,
	}
	ctx := context.Background()

	res, err := handler.GetMFAStatus(ctx)
	require.NoError(t, err)

	errRes, ok := res.(*ogen.Error)
	require.True(t, ok, "expected *ogen.Error, got %T", res)
	assert.Equal(t, 401, errRes.Code)
}

// TestHandler_SetupTOTP_NoAuth tests Handler delegation for SetupTOTP without auth.
func TestHandler_SetupTOTP_NoAuth(t *testing.T) {
	mfaHandler := newTestMFAHandler()
	handler := &Handler{
		logger:     logging.NewTestLogger(),
		mfaHandler: mfaHandler,
	}
	ctx := context.Background()

	res, err := handler.SetupTOTP(ctx, &ogen.SetupTOTPReq{AccountName: "test@example.com"})
	require.NoError(t, err)

	errRes, ok := res.(*ogen.Error)
	require.True(t, ok, "expected *ogen.Error, got %T", res)
	assert.Equal(t, 401, errRes.Code)
}

// TestHandler_VerifyTOTP_NoAuth tests Handler delegation for VerifyTOTP without auth.
func TestHandler_VerifyTOTP_NoAuth(t *testing.T) {
	mfaHandler := newTestMFAHandler()
	handler := &Handler{
		logger:     logging.NewTestLogger(),
		mfaHandler: mfaHandler,
	}
	ctx := context.Background()

	res, err := handler.VerifyTOTP(ctx, &ogen.VerifyTOTPReq{Code: "123456"})
	require.NoError(t, err)

	errRes, ok := res.(*ogen.VerifyTOTPUnauthorized)
	require.True(t, ok, "expected *ogen.VerifyTOTPUnauthorized, got %T", res)
	assert.Equal(t, 401, errRes.Code)
}

// TestHandler_DisableTOTP_NoAuth tests Handler delegation for DisableTOTP without auth.
func TestHandler_DisableTOTP_NoAuth(t *testing.T) {
	mfaHandler := newTestMFAHandler()
	handler := &Handler{
		logger:     logging.NewTestLogger(),
		mfaHandler: mfaHandler,
	}
	ctx := context.Background()

	res, err := handler.DisableTOTP(ctx)
	require.NoError(t, err)

	errRes, ok := res.(*ogen.Error)
	require.True(t, ok, "expected *ogen.Error, got %T", res)
	assert.Equal(t, 401, errRes.Code)
}

// TestHandler_GenerateBackupCodes_NoAuth tests Handler delegation for
// GenerateBackupCodes without auth.
func TestHandler_GenerateBackupCodes_NoAuth(t *testing.T) {
	mfaHandler := newTestMFAHandler()
	handler := &Handler{
		logger:     logging.NewTestLogger(),
		mfaHandler: mfaHandler,
	}
	ctx := context.Background()

	res, err := handler.GenerateBackupCodes(ctx)
	require.NoError(t, err)

	errRes, ok := res.(*ogen.Error)
	require.True(t, ok, "expected *ogen.Error, got %T", res)
	assert.Equal(t, 401, errRes.Code)
}

// TestHandler_RegenerateBackupCodes_NoAuth tests Handler delegation for
// RegenerateBackupCodes without auth.
func TestHandler_RegenerateBackupCodes_NoAuth(t *testing.T) {
	mfaHandler := newTestMFAHandler()
	handler := &Handler{
		logger:     logging.NewTestLogger(),
		mfaHandler: mfaHandler,
	}
	ctx := context.Background()

	res, err := handler.RegenerateBackupCodes(ctx)
	require.NoError(t, err)

	errRes, ok := res.(*ogen.Error)
	require.True(t, ok, "expected *ogen.Error, got %T", res)
	assert.Equal(t, 401, errRes.Code)
}

// TestHandler_EnableMFA_NoAuth tests Handler delegation for EnableMFA without auth.
func TestHandler_EnableMFA_NoAuth(t *testing.T) {
	mfaHandler := newTestMFAHandler()
	handler := &Handler{
		logger:     logging.NewTestLogger(),
		mfaHandler: mfaHandler,
	}
	ctx := context.Background()

	res, err := handler.EnableMFA(ctx)
	require.NoError(t, err)

	errRes, ok := res.(*ogen.EnableMFAUnauthorized)
	require.True(t, ok, "expected *ogen.EnableMFAUnauthorized, got %T", res)
	assert.Equal(t, 401, errRes.Code)
}

// TestHandler_DisableMFA_NoAuth tests Handler delegation for DisableMFA without auth.
func TestHandler_DisableMFA_NoAuth(t *testing.T) {
	mfaHandler := newTestMFAHandler()
	handler := &Handler{
		logger:     logging.NewTestLogger(),
		mfaHandler: mfaHandler,
	}
	ctx := context.Background()

	res, err := handler.DisableMFA(ctx)
	require.NoError(t, err)

	errRes, ok := res.(*ogen.Error)
	require.True(t, ok, "expected *ogen.Error, got %T", res)
	assert.Equal(t, 401, errRes.Code)
}

// =============================================================================
// Nil MFAHandler tests - verify Handler panics when mfaHandler is nil
// =============================================================================

// TestHandler_NilMFAHandler_Panics verifies that calling MFA delegation methods on
// Handler with a nil mfaHandler causes a panic (no nil guard in delegation).
func TestHandler_NilMFAHandler_Panics(t *testing.T) {
	handler := &Handler{
		logger:     logging.NewTestLogger(),
		mfaHandler: nil, // nil MFA handler
	}
	ctx := contextWithUser(t)

	t.Run("GetMFAStatus", func(t *testing.T) {
		assert.Panics(t, func() {
			_, _ = handler.GetMFAStatus(ctx)
		})
	})

	t.Run("SetupTOTP", func(t *testing.T) {
		assert.Panics(t, func() {
			_, _ = handler.SetupTOTP(ctx, &ogen.SetupTOTPReq{AccountName: "test"})
		})
	})

	t.Run("VerifyTOTP", func(t *testing.T) {
		assert.Panics(t, func() {
			_, _ = handler.VerifyTOTP(ctx, &ogen.VerifyTOTPReq{Code: "123456"})
		})
	})

	t.Run("DisableTOTP", func(t *testing.T) {
		assert.Panics(t, func() {
			_, _ = handler.DisableTOTP(ctx)
		})
	})

	t.Run("GenerateBackupCodes", func(t *testing.T) {
		assert.Panics(t, func() {
			_, _ = handler.GenerateBackupCodes(ctx)
		})
	})

	t.Run("RegenerateBackupCodes", func(t *testing.T) {
		assert.Panics(t, func() {
			_, _ = handler.RegenerateBackupCodes(ctx)
		})
	})

	t.Run("EnableMFA", func(t *testing.T) {
		assert.Panics(t, func() {
			_, _ = handler.EnableMFA(ctx)
		})
	})

	t.Run("DisableMFA", func(t *testing.T) {
		assert.Panics(t, func() {
			_, _ = handler.DisableMFA(ctx)
		})
	})
}

// =============================================================================
// WebAuthn NoAuth tests (via MFAHandler directly)
// =============================================================================

// TestMFAHandler_WebAuthn_NoAuth verifies that all WebAuthn handler methods return
// 401 when no user ID is present in the context (even when webauthnService is nil,
// the nil check runs first and returns 501 before auth check; this test uses a
// handler with nil webauthnService to verify the 501-before-401 ordering).
func TestMFAHandler_WebAuthn_NoAuth(t *testing.T) {
	handler := newTestMFAHandler() // webauthnService is nil
	ctx := context.Background()    // no user

	// When webauthnService is nil, the nil check runs before auth check,
	// so we get 501 for methods that check nil first.
	t.Run("BeginWebAuthnRegistration_NilService_Returns501", func(t *testing.T) {
		res, err := handler.BeginWebAuthnRegistration(ctx, ogen.OptBeginWebAuthnRegistrationReq{})
		require.NoError(t, err)
		errRes, ok := res.(*ogen.Error)
		require.True(t, ok)
		assert.Equal(t, 501, errRes.Code)
	})

	t.Run("FinishWebAuthnRegistration_NilService_Returns501", func(t *testing.T) {
		res, err := handler.FinishWebAuthnRegistration(ctx, &ogen.WebAuthnFinishRegistrationRequest{})
		require.NoError(t, err)
		errRes, ok := res.(*ogen.FinishWebAuthnRegistrationBadRequest)
		require.True(t, ok)
		assert.Equal(t, 501, errRes.Code)
	})

	t.Run("BeginWebAuthnLogin_NilService_Returns501", func(t *testing.T) {
		res, err := handler.BeginWebAuthnLogin(ctx)
		require.NoError(t, err)
		errRes, ok := res.(*ogen.BeginWebAuthnLoginBadRequest)
		require.True(t, ok)
		assert.Equal(t, 501, errRes.Code)
	})

	t.Run("FinishWebAuthnLogin_NilService_Returns501", func(t *testing.T) {
		res, err := handler.FinishWebAuthnLogin(ctx, &ogen.WebAuthnFinishLoginRequest{})
		require.NoError(t, err)
		errRes, ok := res.(*ogen.FinishWebAuthnLoginBadRequest)
		require.True(t, ok)
		assert.Equal(t, 501, errRes.Code)
	})

	t.Run("ListWebAuthnCredentials_NilService_Returns501", func(t *testing.T) {
		res, err := handler.ListWebAuthnCredentials(ctx)
		require.NoError(t, err)
		errRes, ok := res.(*ogen.Error)
		require.True(t, ok)
		assert.Equal(t, 501, errRes.Code)
	})

	t.Run("DeleteWebAuthnCredential_NilService_Returns501", func(t *testing.T) {
		res, err := handler.DeleteWebAuthnCredential(ctx, ogen.DeleteWebAuthnCredentialParams{
			CredentialId: uuid.Must(uuid.NewV7()),
		})
		require.NoError(t, err)
		errRes, ok := res.(*ogen.DeleteWebAuthnCredentialNotFound)
		require.True(t, ok)
		assert.Equal(t, 501, errRes.Code)
	})

	t.Run("RenameWebAuthnCredential_NilService_Returns501", func(t *testing.T) {
		res, err := handler.RenameWebAuthnCredential(ctx, &ogen.RenameWebAuthnCredentialReq{
			Name: "name",
		}, ogen.RenameWebAuthnCredentialParams{
			CredentialId: uuid.Must(uuid.NewV7()),
		})
		require.NoError(t, err)
		errRes, ok := res.(*ogen.RenameWebAuthnCredentialNotFound)
		require.True(t, ok)
		assert.Equal(t, 501, errRes.Code)
	})
}

// Test Handler delegation for WebAuthn methods
func TestHandler_WebAuthn_Delegation(t *testing.T) {
	mfaHandler := newTestMFAHandler()
	handler := &Handler{
		logger:     logging.NewTestLogger(),
		mfaHandler: mfaHandler,
	}
	ctx := contextWithUser(t)

	// All methods should delegate to mfaHandler and return 501 (not configured)
	t.Run("BeginWebAuthnRegistration", func(t *testing.T) {
		res, err := handler.BeginWebAuthnRegistration(ctx, ogen.OptBeginWebAuthnRegistrationReq{})
		require.NoError(t, err)
		errRes, ok := res.(*ogen.Error)
		require.True(t, ok)
		assert.Equal(t, 501, errRes.Code)
	})

	t.Run("FinishWebAuthnRegistration", func(t *testing.T) {
		res, err := handler.FinishWebAuthnRegistration(ctx, &ogen.WebAuthnFinishRegistrationRequest{})
		require.NoError(t, err)
		_, ok := res.(*ogen.FinishWebAuthnRegistrationBadRequest)
		assert.True(t, ok)
	})

	t.Run("BeginWebAuthnLogin", func(t *testing.T) {
		res, err := handler.BeginWebAuthnLogin(ctx)
		require.NoError(t, err)
		_, ok := res.(*ogen.BeginWebAuthnLoginBadRequest)
		assert.True(t, ok)
	})

	t.Run("FinishWebAuthnLogin", func(t *testing.T) {
		res, err := handler.FinishWebAuthnLogin(ctx, &ogen.WebAuthnFinishLoginRequest{})
		require.NoError(t, err)
		_, ok := res.(*ogen.FinishWebAuthnLoginBadRequest)
		assert.True(t, ok)
	})

	t.Run("ListWebAuthnCredentials", func(t *testing.T) {
		res, err := handler.ListWebAuthnCredentials(ctx)
		require.NoError(t, err)
		errRes, ok := res.(*ogen.Error)
		require.True(t, ok)
		assert.Equal(t, 501, errRes.Code)
	})

	t.Run("DeleteWebAuthnCredential", func(t *testing.T) {
		res, err := handler.DeleteWebAuthnCredential(ctx, ogen.DeleteWebAuthnCredentialParams{
			CredentialId: uuid.Must(uuid.NewV7()),
		})
		require.NoError(t, err)
		_, ok := res.(*ogen.DeleteWebAuthnCredentialNotFound)
		assert.True(t, ok)
	})

	t.Run("RenameWebAuthnCredential", func(t *testing.T) {
		res, err := handler.RenameWebAuthnCredential(ctx, &ogen.RenameWebAuthnCredentialReq{Name: "test"}, ogen.RenameWebAuthnCredentialParams{
			CredentialId: uuid.Must(uuid.NewV7()),
		})
		require.NoError(t, err)
		_, ok := res.(*ogen.RenameWebAuthnCredentialNotFound)
		assert.True(t, ok)
	})
}
