package mfa

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/infra/logging"
)

// ============================================================================
// Error Variables Tests
// ============================================================================

func TestNoDB_ErrorVariables(t *testing.T) {
	t.Parallel()

	t.Run("MFA manager errors", func(t *testing.T) {
		t.Parallel()
		assert.EqualError(t, ErrMFANotEnabled, "mfa not enabled for user")
		assert.EqualError(t, ErrMFAAlreadySetup, "mfa already setup for user")
		assert.EqualError(t, ErrNoMFAMethod, "user has no mfa methods configured")
	})

	t.Run("WebAuthn errors", func(t *testing.T) {
		t.Parallel()
		assert.EqualError(t, ErrCredentialNotFound, "webauthn credential not found")
		assert.EqualError(t, ErrCloneDetected, "webauthn authenticator clone detected")
		assert.EqualError(t, ErrInvalidCounter, "sign counter did not increment")
		assert.EqualError(t, ErrNoCredentials, "user has no webauthn credentials")
		assert.EqualError(t, ErrCredentialAlreadyUsed, "credential ID already registered")
	})

	t.Run("backup code errors", func(t *testing.T) {
		t.Parallel()
		assert.EqualError(t, ErrNoBackupCodes, "no backup codes available")
		assert.EqualError(t, ErrInvalidBackupCode, "invalid backup code")
		assert.EqualError(t, ErrBackupCodeUsed, "backup code already used")
	})
}

// ============================================================================
// Constants Tests
// ============================================================================

func TestNoDB_Constants(t *testing.T) {
	t.Parallel()

	t.Run("backup code constants", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, 10, BackupCodeCount)
		assert.Equal(t, 8, BackupCodeLength)
	})

	t.Run("verify method constants", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, VerifyMethod("totp"), VerifyMethodTOTP)
		assert.Equal(t, VerifyMethod("webauthn"), VerifyMethodWebAuthn)
		assert.Equal(t, VerifyMethod("backup_code"), VerifyMethodBackupCode)
	})

	t.Run("webauthn session constants", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, 5*time.Minute, webAuthnSessionTTL)
		assert.Equal(t, "webauthn:registration:", webAuthnRegistrationKeyPrefix)
		assert.Equal(t, "webauthn:login:", webAuthnLoginKeyPrefix)
	})
}

// ============================================================================
// MFAStatus Struct Tests
// ============================================================================

func TestNoDB_MFAStatus(t *testing.T) {
	t.Parallel()

	t.Run("zero value", func(t *testing.T) {
		t.Parallel()
		status := MFAStatus{}
		assert.Equal(t, uuid.Nil, status.UserID)
		assert.False(t, status.HasTOTP)
		assert.Equal(t, int64(0), status.WebAuthnCount)
		assert.Equal(t, int64(0), status.UnusedBackupCodes)
		assert.False(t, status.RequireMFA)
		assert.False(t, status.RememberDeviceEnabled)
	})

	t.Run("full values", func(t *testing.T) {
		t.Parallel()
		userID := uuid.Must(uuid.NewV7())
		status := MFAStatus{
			UserID:                userID,
			HasTOTP:               true,
			WebAuthnCount:         3,
			UnusedBackupCodes:     7,
			RequireMFA:            true,
			RememberDeviceEnabled: true,
		}
		assert.Equal(t, userID, status.UserID)
		assert.True(t, status.HasTOTP)
		assert.Equal(t, int64(3), status.WebAuthnCount)
		assert.Equal(t, int64(7), status.UnusedBackupCodes)
		assert.True(t, status.RequireMFA)
		assert.True(t, status.RememberDeviceEnabled)
	})

	t.Run("json marshalling", func(t *testing.T) {
		t.Parallel()
		userID := uuid.Must(uuid.NewV7())
		status := MFAStatus{
			UserID:                userID,
			HasTOTP:               true,
			WebAuthnCount:         2,
			UnusedBackupCodes:     8,
			RequireMFA:            true,
			RememberDeviceEnabled: false,
		}
		data, err := json.Marshal(status)
		require.NoError(t, err)

		var decoded MFAStatus
		err = json.Unmarshal(data, &decoded)
		require.NoError(t, err)
		assert.Equal(t, status, decoded)
	})
}

// ============================================================================
// VerificationResult Tests
// ============================================================================

func TestNoDB_VerificationResult(t *testing.T) {
	t.Parallel()

	t.Run("successful verification", func(t *testing.T) {
		t.Parallel()
		userID := uuid.Must(uuid.NewV7())
		result := VerificationResult{
			Success: true,
			Method:  VerifyMethodTOTP,
			UserID:  userID,
		}
		assert.True(t, result.Success)
		assert.Equal(t, VerifyMethodTOTP, result.Method)
		assert.Equal(t, userID, result.UserID)
	})

	t.Run("failed verification", func(t *testing.T) {
		t.Parallel()
		userID := uuid.Must(uuid.NewV7())
		result := VerificationResult{
			Success: false,
			Method:  VerifyMethodBackupCode,
			UserID:  userID,
		}
		assert.False(t, result.Success)
		assert.Equal(t, VerifyMethodBackupCode, result.Method)
		assert.Equal(t, userID, result.UserID)
	})

	t.Run("json marshalling", func(t *testing.T) {
		t.Parallel()
		userID := uuid.Must(uuid.NewV7())
		result := VerificationResult{
			Success: true,
			Method:  VerifyMethodWebAuthn,
			UserID:  userID,
		}
		data, err := json.Marshal(result)
		require.NoError(t, err)

		var decoded VerificationResult
		err = json.Unmarshal(data, &decoded)
		require.NoError(t, err)
		assert.Equal(t, result, decoded)
	})
}

// ============================================================================
// TOTPSetup Struct Tests
// ============================================================================

func TestNoDB_TOTPSetup(t *testing.T) {
	t.Parallel()

	setup := TOTPSetup{
		Secret: "JBSWY3DPEHPK3PXP",
		QRCode: []byte{0x89, 0x50, 0x4E, 0x47}, // PNG magic bytes
		URL:    "otpauth://totp/Revenge:user@example.com?secret=JBSWY3DPEHPK3PXP&issuer=Revenge",
	}

	assert.NotEmpty(t, setup.Secret)
	assert.NotEmpty(t, setup.QRCode)
	assert.True(t, strings.HasPrefix(setup.URL, "otpauth://totp/"))
	assert.Contains(t, setup.URL, "Revenge")
}

// ============================================================================
// BackupCode Struct Tests
// ============================================================================

func TestNoDB_BackupCode(t *testing.T) {
	t.Parallel()

	now := time.Now()
	code := BackupCode{
		Code:      "1234-5678-90ab-cdef",
		Hash:      "$argon2id$v=19$m=65536,t=3,p=2$test",
		CreatedAt: now,
	}

	assert.Equal(t, "1234-5678-90ab-cdef", code.Code)
	assert.True(t, strings.HasPrefix(code.Hash, "$argon2id$"))
	assert.Equal(t, now, code.CreatedAt)
}

// ============================================================================
// Backup Code Helper Functions Tests
// ============================================================================

func TestNoDB_generateRandomCode(t *testing.T) {
	t.Parallel()

	t.Run("generates valid hex codes", func(t *testing.T) {
		t.Parallel()
		codes := make(map[string]bool)
		for i := 0; i < 100; i++ {
			code, err := generateRandomCode()
			require.NoError(t, err)

			// 8 bytes = 16 hex characters
			assert.Len(t, code, 16)
			assert.Regexp(t, `^[0-9a-f]{16}$`, code)

			// All codes should be unique
			assert.False(t, codes[code], "code should be unique")
			codes[code] = true
		}
	})
}

func TestNoDB_formatCode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "16 character code",
			input:    "1234567890abcdef",
			expected: "1234-5678-90ab-cdef",
		},
		{
			name:     "all zeros",
			input:    "0000000000000000",
			expected: "0000-0000-0000-0000",
		},
		{
			name:     "all f's",
			input:    "ffffffffffffffff",
			expected: "ffff-ffff-ffff-ffff",
		},
		{
			name:     "short code returned as-is",
			input:    "12345",
			expected: "12345",
		},
		{
			name:     "empty code",
			input:    "",
			expected: "",
		},
		{
			name:     "longer code returned as-is",
			input:    "12345678901234567890",
			expected: "12345678901234567890",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := formatCode(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNoDB_normalizeCode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "formatted code with dashes",
			input:    "1234-5678-90AB-CDEF",
			expected: "1234567890abcdef",
		},
		{
			name:     "code with spaces",
			input:    "1234 5678 90ab cdef",
			expected: "1234567890abcdef",
		},
		{
			name:     "mixed case",
			input:    "AbCdEfGh12345678",
			expected: "abcdefgh12345678",
		},
		{
			name:     "already normalized",
			input:    "1234567890abcdef",
			expected: "1234567890abcdef",
		},
		{
			name:     "multiple dashes and spaces",
			input:    "12-34 56-78 90-AB CD-EF",
			expected: "1234567890abcdef",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only dashes",
			input:    "----",
			expected: "",
		},
		{
			name:     "only spaces",
			input:    "    ",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := normalizeCode(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNoDB_ConstantTimeCompare(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		a        string
		b        string
		expected bool
	}{
		{
			name:     "identical strings",
			a:        "1234567890abcdef",
			b:        "1234567890abcdef",
			expected: true,
		},
		{
			name:     "different strings same length",
			a:        "1234567890abcdef",
			b:        "1234567890abcdeg",
			expected: false,
		},
		{
			name:     "different lengths",
			a:        "12345",
			b:        "1234567890",
			expected: false,
		},
		{
			name:     "both empty",
			a:        "",
			b:        "",
			expected: true,
		},
		{
			name:     "one empty",
			a:        "abc",
			b:        "",
			expected: false,
		},
		{
			name:     "first byte different",
			a:        "x234567890abcdef",
			b:        "1234567890abcdef",
			expected: false,
		},
		{
			name:     "last byte different",
			a:        "1234567890abcdex",
			b:        "1234567890abcdef",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := ConstantTimeCompare(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNoDB_formatAndNormalizeRoundTrip(t *testing.T) {
	t.Parallel()

	// Generate a code, format it, then normalize it back
	code, err := generateRandomCode()
	require.NoError(t, err)

	formatted := formatCode(code)
	normalized := normalizeCode(formatted)
	assert.Equal(t, code, normalized, "format+normalize should be identity")
}

// ============================================================================
// WebAuthnUser Tests
// ============================================================================

func TestNoDB_WebAuthnUser(t *testing.T) {
	t.Parallel()

	t.Run("interface methods", func(t *testing.T) {
		t.Parallel()
		userID := uuid.Must(uuid.NewV7())
		user := &WebAuthnUser{
			ID:          userID[:],
			Name:        "testuser",
			DisplayName: "Test User",
			Credentials: []webauthn.Credential{},
		}

		assert.Equal(t, userID[:], user.WebAuthnID())
		assert.Equal(t, "testuser", user.WebAuthnName())
		assert.Equal(t, "Test User", user.WebAuthnDisplayName())
		assert.Empty(t, user.WebAuthnCredentials())
		assert.Empty(t, user.WebAuthnIcon())
	})

	t.Run("with credentials", func(t *testing.T) {
		t.Parallel()
		cred := webauthn.Credential{
			ID:              []byte("credential-id"),
			PublicKey:       []byte("public-key"),
			AttestationType: "none",
		}
		user := &WebAuthnUser{
			ID:          uuid.Must(uuid.NewV7()).NodeID(),
			Name:        "userWithCreds",
			DisplayName: "User With Creds",
			Credentials: []webauthn.Credential{cred},
		}

		assert.Len(t, user.WebAuthnCredentials(), 1)
		assert.Equal(t, []byte("credential-id"), user.WebAuthnCredentials()[0].ID)
	})

	t.Run("implements webauthn.User interface", func(t *testing.T) {
		t.Parallel()
		var _ webauthn.User = (*WebAuthnUser)(nil)
	})
}

// ============================================================================
// Transport Conversion Tests
// ============================================================================

func TestNoDB_convertTransportsToDB(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		transports []protocol.AuthenticatorTransport
		expected   []string
	}{
		{
			name:       "empty",
			transports: []protocol.AuthenticatorTransport{},
			expected:   []string{},
		},
		{
			name:       "single transport",
			transports: []protocol.AuthenticatorTransport{protocol.USB},
			expected:   []string{"usb"},
		},
		{
			name: "multiple transports",
			transports: []protocol.AuthenticatorTransport{
				protocol.USB,
				protocol.NFC,
				protocol.BLE,
				protocol.Internal,
			},
			expected: []string{"usb", "nfc", "ble", "internal"},
		},
		{
			name:       "hybrid transport",
			transports: []protocol.AuthenticatorTransport{protocol.Hybrid},
			expected:   []string{"hybrid"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := convertTransportsToDB(tt.transports)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNoDB_convertTransportsFromDB(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		transports []string
		expected   []protocol.AuthenticatorTransport
	}{
		{
			name:       "empty",
			transports: []string{},
			expected:   []protocol.AuthenticatorTransport{},
		},
		{
			name:       "single transport",
			transports: []string{"usb"},
			expected:   []protocol.AuthenticatorTransport{protocol.USB},
		},
		{
			name:       "multiple transports",
			transports: []string{"usb", "nfc", "ble", "internal"},
			expected: []protocol.AuthenticatorTransport{
				protocol.USB,
				protocol.NFC,
				protocol.BLE,
				protocol.Internal,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := convertTransportsFromDB(tt.transports)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNoDB_convertTransportsRoundTrip(t *testing.T) {
	t.Parallel()

	transports := []protocol.AuthenticatorTransport{
		protocol.USB,
		protocol.NFC,
		protocol.BLE,
		protocol.Internal,
	}

	dbFormat := convertTransportsToDB(transports)
	restored := convertTransportsFromDB(dbFormat)
	assert.Equal(t, transports, restored)
}

// ============================================================================
// Session Data Serialization Tests
// ============================================================================

func TestNoDB_SessionDataToJSON(t *testing.T) {
	t.Parallel()

	t.Run("serializes session data", func(t *testing.T) {
		t.Parallel()
		userID := uuid.Must(uuid.NewV7())
		session := webauthn.SessionData{
			Challenge:        "test-challenge-data-32-bytes-long",
			UserID:           userID[:],
			UserVerification: protocol.VerificationRequired,
		}

		data, err := SessionDataToJSON(session)
		require.NoError(t, err)
		assert.NotEmpty(t, data)

		// Verify it's valid JSON
		var raw map[string]interface{}
		err = json.Unmarshal(data, &raw)
		require.NoError(t, err)
		assert.Contains(t, raw, "challenge")
	})

	t.Run("empty session", func(t *testing.T) {
		t.Parallel()
		session := webauthn.SessionData{}
		data, err := SessionDataToJSON(session)
		require.NoError(t, err)
		assert.NotEmpty(t, data)
	})
}

func TestNoDB_SessionDataFromJSON(t *testing.T) {
	t.Parallel()

	t.Run("deserializes session data", func(t *testing.T) {
		t.Parallel()
		userID := uuid.Must(uuid.NewV7())
		original := webauthn.SessionData{
			Challenge:        "test-challenge-data",
			UserID:           userID[:],
			UserVerification: protocol.VerificationRequired,
		}

		data, err := SessionDataToJSON(original)
		require.NoError(t, err)

		restored, err := SessionDataFromJSON(data)
		require.NoError(t, err)
		assert.Equal(t, original.Challenge, restored.Challenge)
		assert.Equal(t, original.UserID, restored.UserID)
		assert.Equal(t, original.UserVerification, restored.UserVerification)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		t.Parallel()
		_, err := SessionDataFromJSON([]byte(`{"challenge": "invalid}`))
		require.Error(t, err)
	})

	t.Run("empty JSON object", func(t *testing.T) {
		t.Parallel()
		result, err := SessionDataFromJSON([]byte(`{}`))
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result.Challenge)
	})

	t.Run("nil bytes", func(t *testing.T) {
		t.Parallel()
		_, err := SessionDataFromJSON(nil)
		require.Error(t, err)
	})
}

func TestNoDB_SessionDataRoundTrip(t *testing.T) {
	t.Parallel()

	userID := uuid.Must(uuid.NewV7())
	session := webauthn.SessionData{
		Challenge:            "test-challenge-32-bytes-minimum-length",
		UserID:               userID[:],
		AllowedCredentialIDs: [][]byte{[]byte("cred-1"), []byte("cred-2")},
		UserVerification:     protocol.VerificationPreferred,
	}

	data, err := SessionDataToJSON(session)
	require.NoError(t, err)

	restored, err := SessionDataFromJSON(data)
	require.NoError(t, err)

	assert.Equal(t, session.Challenge, restored.Challenge)
	assert.Equal(t, session.UserID, restored.UserID)
	assert.Equal(t, session.UserVerification, restored.UserVerification)
	assert.Len(t, restored.AllowedCredentialIDs, 2)
}

// ============================================================================
// NewWebAuthnService Tests
// ============================================================================

func TestNoDB_NewWebAuthnService(t *testing.T) {
	t.Parallel()
	logger := logging.NewTestLogger()

	tests := []struct {
		name          string
		rpDisplayName string
		rpID          string
		rpOrigins     []string
		wantErr       bool
	}{
		{
			name:          "valid localhost config",
			rpDisplayName: "Test App",
			rpID:          "localhost",
			rpOrigins:     []string{"http://localhost:3000"},
			wantErr:       false,
		},
		{
			name:          "valid production config",
			rpDisplayName: "Revenge",
			rpID:          "revenge.example.com",
			rpOrigins:     []string{"https://revenge.example.com"},
			wantErr:       false,
		},
		{
			name:          "multiple origins",
			rpDisplayName: "Revenge",
			rpID:          "revenge.example.com",
			rpOrigins:     []string{"https://revenge.example.com", "https://app.revenge.example.com"},
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc, err := NewWebAuthnService(nil, logger, nil, tt.rpDisplayName, tt.rpID, tt.rpOrigins)
			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, svc)
			} else {
				require.NoError(t, err)
				require.NotNil(t, svc)
				assert.NotNil(t, svc.webAuthn)
				assert.Nil(t, svc.cache, "cache should be nil when not provided")
			}
		})
	}
}

// ============================================================================
// WebAuthnService HasCache Tests
// ============================================================================

func TestNoDB_WebAuthnService_HasCache(t *testing.T) {
	t.Parallel()
	logger := logging.NewTestLogger()

	t.Run("no cache", func(t *testing.T) {
		t.Parallel()
		svc, err := NewWebAuthnService(nil, logger, nil, "Test", "localhost", []string{"http://localhost"})
		require.NoError(t, err)
		assert.False(t, svc.HasCache())
	})

	t.Run("with cache", func(t *testing.T) {
		t.Parallel()
		sessionCache, err := cache.NewNamedCache(nil, 100, 5*time.Minute, "webauthn-test")
		require.NoError(t, err)
		defer sessionCache.Close()

		svc, err := NewWebAuthnService(nil, logger, sessionCache, "Test", "localhost", []string{"http://localhost"})
		require.NoError(t, err)
		assert.True(t, svc.HasCache())
	})
}

// ============================================================================
// WebAuthnService Session Tests (without DB)
// ============================================================================

func TestNoDB_WebAuthnService_storeSession(t *testing.T) {
	t.Parallel()
	logger := logging.NewTestLogger()
	ctx := context.Background()

	t.Run("nil cache does not error", func(t *testing.T) {
		t.Parallel()
		svc, err := NewWebAuthnService(nil, logger, nil, "Test", "localhost", []string{"http://localhost"})
		require.NoError(t, err)

		session := &webauthn.SessionData{Challenge: "test-challenge"}
		err = svc.storeSession(ctx, webAuthnRegistrationKeyPrefix, uuid.Must(uuid.NewV7()), session)
		require.NoError(t, err)
	})

	t.Run("with cache stores and retrieves session", func(t *testing.T) {
		t.Parallel()
		sessionCache, err := cache.NewNamedCache(nil, 100, 5*time.Minute, "webauthn-test")
		require.NoError(t, err)
		defer sessionCache.Close()

		svc, err := NewWebAuthnService(nil, logger, sessionCache, "Test", "localhost", []string{"http://localhost"})
		require.NoError(t, err)

		userID := uuid.Must(uuid.NewV7())
		session := &webauthn.SessionData{Challenge: "test-challenge-data"}

		err = svc.storeSession(ctx, webAuthnRegistrationKeyPrefix, userID, session)
		require.NoError(t, err)

		retrieved, err := svc.getSession(ctx, webAuthnRegistrationKeyPrefix, userID)
		require.NoError(t, err)
		assert.Equal(t, "test-challenge-data", retrieved.Challenge)
	})
}

func TestNoDB_WebAuthnService_getSession(t *testing.T) {
	t.Parallel()
	logger := logging.NewTestLogger()
	ctx := context.Background()

	t.Run("nil cache returns error", func(t *testing.T) {
		t.Parallel()
		svc, err := NewWebAuthnService(nil, logger, nil, "Test", "localhost", []string{"http://localhost"})
		require.NoError(t, err)

		_, err = svc.getSession(ctx, webAuthnRegistrationKeyPrefix, uuid.Must(uuid.NewV7()))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cache not configured")
	})

	t.Run("cache miss returns error", func(t *testing.T) {
		t.Parallel()
		sessionCache, err := cache.NewNamedCache(nil, 100, 5*time.Minute, "webauthn-test")
		require.NoError(t, err)
		defer sessionCache.Close()

		svc, err := NewWebAuthnService(nil, logger, sessionCache, "Test", "localhost", []string{"http://localhost"})
		require.NoError(t, err)

		_, err = svc.getSession(ctx, webAuthnRegistrationKeyPrefix, uuid.Must(uuid.NewV7()))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "session not found")
	})
}

func TestNoDB_WebAuthnService_deleteSession(t *testing.T) {
	t.Parallel()
	logger := logging.NewTestLogger()
	ctx := context.Background()

	t.Run("nil cache does not panic", func(t *testing.T) {
		t.Parallel()
		svc, err := NewWebAuthnService(nil, logger, nil, "Test", "localhost", []string{"http://localhost"})
		require.NoError(t, err)

		// Should not panic
		svc.deleteSession(ctx, webAuthnRegistrationKeyPrefix, uuid.Must(uuid.NewV7()))
	})

	t.Run("with cache deletes session", func(t *testing.T) {
		t.Parallel()
		sessionCache, err := cache.NewNamedCache(nil, 100, 5*time.Minute, "webauthn-test")
		require.NoError(t, err)
		defer sessionCache.Close()

		svc, err := NewWebAuthnService(nil, logger, sessionCache, "Test", "localhost", []string{"http://localhost"})
		require.NoError(t, err)

		userID := uuid.Must(uuid.NewV7())
		session := &webauthn.SessionData{Challenge: "test-challenge"}

		err = svc.storeSession(ctx, webAuthnRegistrationKeyPrefix, userID, session)
		require.NoError(t, err)

		svc.deleteSession(ctx, webAuthnRegistrationKeyPrefix, userID)

		_, err = svc.getSession(ctx, webAuthnRegistrationKeyPrefix, userID)
		require.Error(t, err, "session should be deleted")
	})
}

func TestNoDB_WebAuthnService_GetRegistrationSession(t *testing.T) {
	t.Parallel()
	logger := logging.NewTestLogger()
	ctx := context.Background()

	t.Run("nil cache returns error", func(t *testing.T) {
		t.Parallel()
		svc, err := NewWebAuthnService(nil, logger, nil, "Test", "localhost", []string{"http://localhost"})
		require.NoError(t, err)

		_, err = svc.GetRegistrationSession(ctx, uuid.Must(uuid.NewV7()))
		require.Error(t, err)
	})
}

func TestNoDB_WebAuthnService_GetLoginSession(t *testing.T) {
	t.Parallel()
	logger := logging.NewTestLogger()
	ctx := context.Background()

	t.Run("nil cache returns error", func(t *testing.T) {
		t.Parallel()
		svc, err := NewWebAuthnService(nil, logger, nil, "Test", "localhost", []string{"http://localhost"})
		require.NoError(t, err)

		_, err = svc.GetLoginSession(ctx, uuid.Must(uuid.NewV7()))
		require.Error(t, err)
	})
}

func TestNoDB_WebAuthnService_DeleteRegistrationSession(t *testing.T) {
	t.Parallel()
	logger := logging.NewTestLogger()
	ctx := context.Background()

	t.Run("nil cache does not panic", func(t *testing.T) {
		t.Parallel()
		svc, err := NewWebAuthnService(nil, logger, nil, "Test", "localhost", []string{"http://localhost"})
		require.NoError(t, err)
		svc.DeleteRegistrationSession(ctx, uuid.Must(uuid.NewV7()))
	})
}

func TestNoDB_WebAuthnService_DeleteLoginSession(t *testing.T) {
	t.Parallel()
	logger := logging.NewTestLogger()
	ctx := context.Background()

	t.Run("nil cache does not panic", func(t *testing.T) {
		t.Parallel()
		svc, err := NewWebAuthnService(nil, logger, nil, "Test", "localhost", []string{"http://localhost"})
		require.NoError(t, err)
		svc.DeleteLoginSession(ctx, uuid.Must(uuid.NewV7()))
	})
}

// ============================================================================
// Session Store/Get/Delete Full Lifecycle Tests
// ============================================================================

func TestNoDB_WebAuthnService_SessionLifecycle(t *testing.T) {
	t.Parallel()
	logger := logging.NewTestLogger()
	ctx := context.Background()

	sessionCache, err := cache.NewNamedCache(nil, 100, 5*time.Minute, "webauthn-lifecycle")
	require.NoError(t, err)
	defer sessionCache.Close()

	svc, err := NewWebAuthnService(nil, logger, sessionCache, "Test", "localhost", []string{"http://localhost"})
	require.NoError(t, err)

	userID := uuid.Must(uuid.NewV7())

	t.Run("registration session lifecycle", func(t *testing.T) {
		session := &webauthn.SessionData{
			Challenge: "reg-challenge",
			UserID:    userID[:],
		}

		// Store
		err := svc.storeSession(ctx, webAuthnRegistrationKeyPrefix, userID, session)
		require.NoError(t, err)

		// Get
		retrieved, err := svc.GetRegistrationSession(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, "reg-challenge", retrieved.Challenge)

		// Delete
		svc.DeleteRegistrationSession(ctx, userID)

		// Verify deleted
		_, err = svc.GetRegistrationSession(ctx, userID)
		require.Error(t, err)
	})

	t.Run("login session lifecycle", func(t *testing.T) {
		session := &webauthn.SessionData{
			Challenge: "login-challenge",
			UserID:    userID[:],
		}

		// Store
		err := svc.storeSession(ctx, webAuthnLoginKeyPrefix, userID, session)
		require.NoError(t, err)

		// Get
		retrieved, err := svc.GetLoginSession(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, "login-challenge", retrieved.Challenge)

		// Delete
		svc.DeleteLoginSession(ctx, userID)

		// Verify deleted
		_, err = svc.GetLoginSession(ctx, userID)
		require.Error(t, err)
	})
}

// ============================================================================
// NewMFAManager Test (constructor only, no DB)
// ============================================================================

func TestNoDB_NewMFAManager(t *testing.T) {
	t.Parallel()
	logger := logging.NewTestLogger()

	manager := NewMFAManager(nil, nil, nil, nil, logger)
	require.NotNil(t, manager)
	assert.Nil(t, manager.queries)
	assert.Nil(t, manager.totp)
	assert.Nil(t, manager.webauthn)
	assert.Nil(t, manager.backupCodes)
	assert.NotNil(t, manager.logger)
}

// ============================================================================
// NewTOTPService Test (constructor only, no DB)
// ============================================================================

func TestNoDB_NewTOTPService(t *testing.T) {
	t.Parallel()
	logger := logging.NewTestLogger()

	svc := NewTOTPService(nil, nil, logger, "TestIssuer")
	require.NotNil(t, svc)
	assert.Nil(t, svc.queries)
	assert.Nil(t, svc.encryptor)
	assert.NotNil(t, svc.logger)
	assert.Equal(t, "TestIssuer", svc.issuer)
}

// ============================================================================
// NewBackupCodesService Test (constructor only, no DB)
// ============================================================================

func TestNoDB_NewBackupCodesService(t *testing.T) {
	t.Parallel()
	logger := logging.NewTestLogger()

	svc := NewBackupCodesService(nil, logger)
	require.NotNil(t, svc)
	assert.Nil(t, svc.queries)
	assert.NotNil(t, svc.hasher)
	assert.NotNil(t, svc.logger)
}

// ============================================================================
// Module Tests (NewTOTPServiceFromConfig, NewWebAuthnServiceFromConfig)
// ============================================================================

func TestNoDB_NewTOTPServiceFromConfig(t *testing.T) {
	t.Parallel()
	logger := logging.NewTestLogger()

	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "localhost",
			Port: 8080,
		},
	}

	svc := NewTOTPServiceFromConfig(nil, nil, logger, cfg)
	require.NotNil(t, svc)
	assert.Equal(t, "Revenge", svc.issuer)
	assert.Nil(t, svc.queries)
	assert.Nil(t, svc.encryptor)
	assert.NotNil(t, svc.logger)
}

func TestNoDB_NewWebAuthnServiceFromConfig(t *testing.T) {
	t.Parallel()
	logger := logging.NewTestLogger()

	t.Run("localhost config", func(t *testing.T) {
		t.Parallel()
		cfg := &config.Config{
			Server: config.ServerConfig{
				Host: "localhost",
				Port: 8080,
			},
		}

		svc, err := NewWebAuthnServiceFromConfig(nil, logger, cfg, nil)
		require.NoError(t, err)
		require.NotNil(t, svc)
		assert.Nil(t, svc.cache, "cache should be nil when cacheClient is nil")
		assert.NotNil(t, svc.webAuthn)
	})

	t.Run("0.0.0.0 defaults to localhost", func(t *testing.T) {
		t.Parallel()
		cfg := &config.Config{
			Server: config.ServerConfig{
				Host: "0.0.0.0",
				Port: 3000,
			},
		}

		svc, err := NewWebAuthnServiceFromConfig(nil, logger, cfg, nil)
		require.NoError(t, err)
		require.NotNil(t, svc)
	})

	t.Run("empty host defaults to localhost", func(t *testing.T) {
		t.Parallel()
		cfg := &config.Config{
			Server: config.ServerConfig{
				Host: "",
				Port: 3000,
			},
		}

		svc, err := NewWebAuthnServiceFromConfig(nil, logger, cfg, nil)
		require.NoError(t, err)
		require.NotNil(t, svc)
	})

	t.Run("custom host", func(t *testing.T) {
		t.Parallel()
		cfg := &config.Config{
			Server: config.ServerConfig{
				Host: "revenge.example.com",
				Port: 443,
			},
		}

		svc, err := NewWebAuthnServiceFromConfig(nil, logger, cfg, nil)
		require.NoError(t, err)
		require.NotNil(t, svc)
	})

}
