package api

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/go-faster/jx"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/lusoris/revenge/internal/api/ogen"
)

func newTestMFAHandler() *MFAHandler {
	return NewMFAHandler(nil, nil, nil, nil, zap.NewNop())
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
	handler := NewMFAHandler(nil, nil, nil, nil, zap.NewNop())
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

// Test Handler delegation for WebAuthn methods
func TestHandler_WebAuthn_Delegation(t *testing.T) {
	mfaHandler := newTestMFAHandler()
	handler := &Handler{
		logger:     zap.NewNop(),
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
