package mfa

import (
	"testing"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestNewWebAuthnService(t *testing.T) {
	logger := zaptest.NewLogger(t)

	tests := []struct {
		name          string
		rpDisplayName string
		rpID          string
		rpOrigins     []string
		wantErr       bool
	}{
		{
			name:          "valid configuration",
			rpDisplayName: "Revenge Test",
			rpID:          "localhost",
			rpOrigins:     []string{"http://localhost:3000"},
			wantErr:       false,
		},
		{
			name:          "multiple origins",
			rpDisplayName: "Revenge Production",
			rpID:          "revenge.example.com",
			rpOrigins:     []string{"https://revenge.example.com", "https://app.revenge.example.com"},
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := NewWebAuthnService(nil, logger, tt.rpDisplayName, tt.rpID, tt.rpOrigins)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, service)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, service)
				assert.NotNil(t, service.webAuthn)
			}
		})
	}
}

func TestWebAuthnUser_Interface(t *testing.T) {
	userID := uuid.New()
	user := &WebAuthnUser{
		ID:          userID[:],
		Name:        "testuser",
		DisplayName: "Test User",
		Credentials: []webauthn.Credential{},
	}

	// Test interface compliance
	assert.Equal(t, userID[:], user.WebAuthnID())
	assert.Equal(t, "testuser", user.WebAuthnName())
	assert.Equal(t, "Test User", user.WebAuthnDisplayName())
	assert.Empty(t, user.WebAuthnCredentials())
	assert.Empty(t, user.WebAuthnIcon())
}

func TestConvertTransports(t *testing.T) {
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
			name: "single transport",
			transports: []protocol.AuthenticatorTransport{
				protocol.USB,
			},
			expected: []string{"usb"},
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert to DB format
			dbTransports := convertTransportsToDB(tt.transports)
			assert.Equal(t, tt.expected, dbTransports)

			// Convert back
			originalTransports := convertTransportsFromDB(dbTransports)
			assert.Equal(t, tt.transports, originalTransports)
		})
	}
}

func TestSessionDataSerialization(t *testing.T) {
	// Create test session data
	userID := uuid.New()

	session := webauthn.SessionData{
		Challenge:        "test-challenge-data-32-bytes-long",
		UserID:           userID[:],
		UserVerification: protocol.VerificationRequired,
	}

	// Serialize
	data, err := SessionDataToJSON(session)
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	// Deserialize
	restored, err := SessionDataFromJSON(data)
	require.NoError(t, err)
	assert.Equal(t, session.Challenge, restored.Challenge)
	assert.Equal(t, session.UserID, restored.UserID)
	assert.Equal(t, session.UserVerification, restored.UserVerification)
}

func TestSessionDataSerialization_InvalidJSON(t *testing.T) {
	invalidJSON := []byte(`{"challenge": "invalid}`)

	_, err := SessionDataFromJSON(invalidJSON)
	assert.Error(t, err)
}

// Integration tests would go here
// These would require actual database connection and mock WebAuthn responses

func TestWebAuthnService_CredentialLifecycle(t *testing.T) {
	t.Skip("Integration test - requires database and WebAuthn mock")

	// This test would cover:
	// 1. BeginRegistration -> creates challenge
	// 2. FinishRegistration -> validates and stores credential
	// 3. ListCredentials -> retrieves stored credentials
	// 4. BeginLogin -> creates authentication challenge
	// 5. FinishLogin -> validates assertion and updates counter
	// 6. RenameCredential -> updates credential name
	// 7. DeleteCredential -> removes credential
}

func TestWebAuthnService_CloneDetection(t *testing.T) {
	t.Skip("Integration test - requires database and WebAuthn mock")

	// This test would cover:
	// 1. Register a credential (counter = 0)
	// 2. Authenticate successfully (counter = 1)
	// 3. Authenticate successfully (counter = 2)
	// 4. Attempt authentication with old counter (counter = 1)
	// 5. Should fail with ErrCloneDetected
	// 6. Credential should be marked as clone_detected = true
}

func TestWebAuthnService_MultipleCredentials(t *testing.T) {
	t.Skip("Integration test - requires database")

	// This test would cover:
	// 1. Register multiple credentials for same user (e.g., YubiKey + TouchID)
	// 2. BeginLogin should return allowCredentials with both
	// 3. Authenticate with first credential
	// 4. Authenticate with second credential
	// 5. Delete first credential
	// 6. BeginLogin should only return second credential
}