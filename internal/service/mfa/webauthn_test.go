package mfa

import (
	"context"
	"crypto/rand"
	"log/slog"
	"math"
	"testing"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/crypto"
	"github.com/lusoris/revenge/internal/infra/cache"
	db "github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/testutil"
	"github.com/lusoris/revenge/internal/util"
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
			service, err := NewWebAuthnService(nil, logger, nil, tt.rpDisplayName, tt.rpID, tt.rpOrigins)

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
	userID := uuid.Must(uuid.NewV7())
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
	userID := uuid.Must(uuid.NewV7())

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

// Unit tests for safe integer conversions

func TestSafeUint32ToInt32(t *testing.T) {
	tests := []struct {
		name     string
		input    uint32
		expected int32
	}{
		{
			name:     "zero",
			input:    0,
			expected: 0,
		},
		{
			name:     "small value",
			input:    12345,
			expected: 12345,
		},
		{
			name:     "max int32 value",
			input:    math.MaxInt32,
			expected: math.MaxInt32,
		},
		{
			name:     "just over max int32",
			input:    math.MaxInt32 + 1,
			expected: math.MaxInt32, // Should cap at max
		},
		{
			name:     "max uint32 value",
			input:    math.MaxUint32,
			expected: math.MaxInt32, // Should cap at max
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.SafeUint32ToInt32(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSafeInt32ToUint32(t *testing.T) {
	tests := []struct {
		name     string
		input    int32
		expected uint32
	}{
		{
			name:     "zero",
			input:    0,
			expected: 0,
		},
		{
			name:     "positive value",
			input:    12345,
			expected: 12345,
		},
		{
			name:     "max int32 value",
			input:    math.MaxInt32,
			expected: math.MaxInt32,
		},
		{
			name:     "negative value",
			input:    -1,
			expected: 0, // Should treat negative as 0
		},
		{
			name:     "min int32 value",
			input:    math.MinInt32,
			expected: 0, // Should treat negative as 0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.SafeInt32ToUint32(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Integration tests

func setupWebAuthnService(t *testing.T) (*WebAuthnService, *db.Queries, context.Context, uuid.UUID) {
	t.Helper()

	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())
	ctx := context.Background()
	logger := zaptest.NewLogger(t)

	service, err := NewWebAuthnService(queries, logger, nil, "Test App", "localhost", []string{"http://localhost:3000"})
	require.NoError(t, err)

	// Create a test user
	userID := createTestUserForWebAuthn(t, queries, ctx)

	return service, queries, ctx, userID
}

func createTestUserForWebAuthn(t *testing.T, queries *db.Queries, ctx context.Context) uuid.UUID {
	t.Helper()

	isActive := true
	user, err := queries.CreateUser(ctx, db.CreateUserParams{
		Email:        "webauthn_" + uuid.Must(uuid.NewV7()).String()[:8] + "@example.com",
		Username:     "webauthn_" + uuid.Must(uuid.NewV7()).String()[:8],
		PasswordHash: "$argon2id$v=19$m=65536,t=3,p=4$test$test",
		IsActive:     &isActive,
	})
	require.NoError(t, err)

	return user.ID
}

func TestWebAuthnService_HasWebAuthn(t *testing.T) {
	t.Parallel()

	t.Run("returns false for new user", func(t *testing.T) {
		t.Parallel()
		service, _, ctx, userID := setupWebAuthnService(t)

		hasWebAuthn, err := service.HasWebAuthn(ctx, userID)
		require.NoError(t, err)
		assert.False(t, hasWebAuthn)
	})

	t.Run("returns true after credential added", func(t *testing.T) {
		t.Parallel()
		service, queries, ctx, userID := setupWebAuthnService(t)

		// Manually insert a credential for testing
		name := "Test Security Key"
		_, err := queries.CreateWebAuthnCredential(ctx, db.CreateWebAuthnCredentialParams{
			UserID:          userID,
			CredentialID:    []byte("test-credential-id-12345"),
			PublicKey:       []byte("test-public-key-data"),
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

		hasWebAuthn, err := service.HasWebAuthn(ctx, userID)
		require.NoError(t, err)
		assert.True(t, hasWebAuthn)
	})
}

func TestWebAuthnService_ListCredentials(t *testing.T) {
	t.Parallel()

	t.Run("returns empty list for new user", func(t *testing.T) {
		t.Parallel()
		service, _, ctx, userID := setupWebAuthnService(t)

		creds, err := service.ListCredentials(ctx, userID)
		require.NoError(t, err)
		assert.Empty(t, creds)
	})

	t.Run("returns credentials after adding", func(t *testing.T) {
		t.Parallel()
		service, queries, ctx, userID := setupWebAuthnService(t)

		// Add two credentials
		name1 := "YubiKey 5"
		name2 := "TouchID"
		_, err := queries.CreateWebAuthnCredential(ctx, db.CreateWebAuthnCredentialParams{
			UserID:          userID,
			CredentialID:    []byte("cred-1-id"),
			PublicKey:       []byte("pub-key-1"),
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
			CredentialID:    []byte("cred-2-id"),
			PublicKey:       []byte("pub-key-2"),
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

		creds, err := service.ListCredentials(ctx, userID)
		require.NoError(t, err)
		assert.Len(t, creds, 2)
	})
}

func TestWebAuthnService_DeleteCredential(t *testing.T) {
	t.Parallel()

	service, queries, ctx, userID := setupWebAuthnService(t)

	// Add a credential
	name := "Test Key"
	cred, err := queries.CreateWebAuthnCredential(ctx, db.CreateWebAuthnCredentialParams{
		UserID:          userID,
		CredentialID:    []byte("delete-test-cred-id"),
		PublicKey:       []byte("pub-key"),
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

	// Verify credential exists
	hasWebAuthn, err := service.HasWebAuthn(ctx, userID)
	require.NoError(t, err)
	assert.True(t, hasWebAuthn)

	// Delete credential
	err = service.DeleteCredential(ctx, userID, cred.ID)
	require.NoError(t, err)

	// Verify credential is gone
	hasWebAuthn, err = service.HasWebAuthn(ctx, userID)
	require.NoError(t, err)
	assert.False(t, hasWebAuthn)
}

func TestWebAuthnService_RenameCredential(t *testing.T) {
	t.Parallel()

	service, queries, ctx, userID := setupWebAuthnService(t)

	// Add a credential
	originalName := "Original Name"
	cred, err := queries.CreateWebAuthnCredential(ctx, db.CreateWebAuthnCredentialParams{
		UserID:          userID,
		CredentialID:    []byte("rename-test-cred-id"),
		PublicKey:       []byte("pub-key"),
		AttestationType: "none",
		Transports:      []string{"usb"},
		BackupEligible:  false,
		BackupState:     false,
		UserPresent:     true,
		UserVerified:    true,
		Aaguid:          []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		Name:            &originalName,
	})
	require.NoError(t, err)

	// Rename credential
	newName := "My YubiKey 5C"
	err = service.RenameCredential(ctx, cred.ID, newName)
	require.NoError(t, err)

	// Verify name changed
	creds, err := service.ListCredentials(ctx, userID)
	require.NoError(t, err)
	require.Len(t, creds, 1)
	assert.Equal(t, newName, *creds[0].Name)
}

func TestWebAuthnService_BeginRegistration(t *testing.T) {
	t.Parallel()

	t.Run("returns credential creation options", func(t *testing.T) {
		t.Parallel()
		service, _, ctx, userID := setupWebAuthnService(t)

		options, err := service.BeginRegistration(ctx, userID, "testuser", "Test User")
		require.NoError(t, err)
		require.NotNil(t, options)

		// Verify options structure
		assert.NotEmpty(t, options.Response.Challenge)
		assert.Equal(t, "Test App", options.Response.RelyingParty.Name)
		assert.Equal(t, "localhost", options.Response.RelyingParty.ID)
		assert.Equal(t, "testuser", options.Response.User.Name)
		assert.Equal(t, "Test User", options.Response.User.DisplayName)
	})

	t.Run("succeeds with existing credentials", func(t *testing.T) {
		t.Parallel()
		service, queries, ctx, userID := setupWebAuthnService(t)

		// Add an existing credential
		name := "Existing Key"
		_, err := queries.CreateWebAuthnCredential(ctx, db.CreateWebAuthnCredentialParams{
			UserID:          userID,
			CredentialID:    []byte("existing-cred-id"),
			PublicKey:       []byte("pub-key"),
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

		// Begin registration should work even when user has existing credentials
		// The webauthn library handles credential exclusion internally
		options, err := service.BeginRegistration(ctx, userID, "testuser", "Test User")
		require.NoError(t, err)
		require.NotNil(t, options)

		// Verify basic options are still correct
		assert.NotEmpty(t, options.Response.Challenge)
	})
}

func TestWebAuthnService_BeginLogin(t *testing.T) {
	t.Parallel()

	t.Run("fails if no credentials exist", func(t *testing.T) {
		t.Parallel()
		service, _, ctx, userID := setupWebAuthnService(t)

		_, err := service.BeginLogin(ctx, userID, "testuser", "Test User")
		assert.ErrorIs(t, err, ErrNoCredentials)
	})

	t.Run("returns assertion options with credentials", func(t *testing.T) {
		t.Parallel()
		service, queries, ctx, userID := setupWebAuthnService(t)

		// Add a credential
		name := "Test Key"
		_, err := queries.CreateWebAuthnCredential(ctx, db.CreateWebAuthnCredentialParams{
			UserID:          userID,
			CredentialID:    []byte("login-test-cred-id"),
			PublicKey:       []byte("pub-key"),
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

		options, err := service.BeginLogin(ctx, userID, "testuser", "Test User")
		require.NoError(t, err)
		require.NotNil(t, options)

		// Verify options
		assert.NotEmpty(t, options.Response.Challenge)
		assert.Equal(t, "localhost", options.Response.RelyingPartyID)
		assert.Len(t, options.Response.AllowedCredentials, 1)
	})

	t.Run("excludes cloned credentials", func(t *testing.T) {
		t.Parallel()
		service, queries, ctx, userID := setupWebAuthnService(t)

		// Add a credential then mark it as cloned
		name := "Cloned Key"
		cred, err := queries.CreateWebAuthnCredential(ctx, db.CreateWebAuthnCredentialParams{
			UserID:          userID,
			CredentialID:    []byte("cloned-cred-id"),
			PublicKey:       []byte("pub-key"),
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

		// Mark as cloned
		err = queries.MarkWebAuthnCloneDetected(ctx, cred.CredentialID)
		require.NoError(t, err)

		// Should fail because all credentials are marked as cloned
		_, err = service.BeginLogin(ctx, userID, "testuser", "Test User")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "all credentials marked as cloned")
	})
}

func TestWebAuthnService_FinishRegistration_InvalidResponse(t *testing.T) {
	t.Parallel()

	t.Run("fails with invalid attestation data", func(t *testing.T) {
		t.Parallel()
		service, _, ctx, userID := setupWebAuthnService(t)

		// Create a minimal but invalid parsed response
		response := &protocol.ParsedCredentialCreationData{
			ParsedPublicKeyCredential: protocol.ParsedPublicKeyCredential{
				ParsedCredential: protocol.ParsedCredential{
					ID:   "fake-credential-id",
					Type: "public-key",
				},
				RawID: []byte("fake-credential-id"),
			},
			Response: protocol.ParsedAttestationResponse{
				CollectedClientData: protocol.CollectedClientData{
					Type:      protocol.CreateCeremony,
					Challenge: "wrong-challenge",
					Origin:    "http://localhost:3000",
				},
			},
		}

		sessionData := webauthn.SessionData{
			Challenge: "test-challenge-that-doesnt-match",
			UserID:    userID[:],
		}

		err := service.FinishRegistration(ctx, userID, "testuser", "Test User", response, sessionData, "Test Key")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create credential")
	})

	t.Run("fails with invalid attestation data and existing credentials", func(t *testing.T) {
		t.Parallel()
		service, queries, ctx, userID := setupWebAuthnService(t)

		// Add an existing credential so the loop body executes
		name := "Existing Key"
		_, err := queries.CreateWebAuthnCredential(ctx, db.CreateWebAuthnCredentialParams{
			UserID:          userID,
			CredentialID:    []byte("existing-cred-for-finish-reg"),
			PublicKey:       []byte("pub-key"),
			AttestationType: "none",
			Transports:      []string{"usb", "nfc"},
			BackupEligible:  true,
			BackupState:     false,
			UserPresent:     true,
			UserVerified:    true,
			Aaguid:          []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			Name:            &name,
		})
		require.NoError(t, err)

		response := &protocol.ParsedCredentialCreationData{
			ParsedPublicKeyCredential: protocol.ParsedPublicKeyCredential{
				ParsedCredential: protocol.ParsedCredential{
					ID:   "fake-credential-id",
					Type: "public-key",
				},
				RawID: []byte("fake-credential-id"),
			},
			Response: protocol.ParsedAttestationResponse{
				CollectedClientData: protocol.CollectedClientData{
					Type:      protocol.CreateCeremony,
					Challenge: "wrong-challenge",
					Origin:    "http://localhost:3000",
				},
			},
		}

		sessionData := webauthn.SessionData{
			Challenge: "test-challenge-that-doesnt-match",
			UserID:    userID[:],
		}

		err = service.FinishRegistration(ctx, userID, "testuser", "Test User", response, sessionData, "My Key")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create credential")
	})
}

func TestWebAuthnService_FinishLogin_InvalidResponse(t *testing.T) {
	t.Parallel()

	t.Run("fails with no credentials", func(t *testing.T) {
		t.Parallel()
		service, _, ctx, userID := setupWebAuthnService(t)

		response := &protocol.ParsedCredentialAssertionData{
			ParsedPublicKeyCredential: protocol.ParsedPublicKeyCredential{
				ParsedCredential: protocol.ParsedCredential{
					ID:   "fake-credential-id",
					Type: "public-key",
				},
				RawID: []byte("fake-credential-id"),
			},
			Response: protocol.ParsedAssertionResponse{
				CollectedClientData: protocol.CollectedClientData{
					Type:      protocol.AssertCeremony,
					Challenge: "wrong-challenge",
					Origin:    "http://localhost:3000",
				},
			},
		}

		sessionData := webauthn.SessionData{
			Challenge: "test-challenge",
			UserID:    userID[:],
		}

		err := service.FinishLogin(ctx, userID, "testuser", "Test User", response, sessionData)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to validate login")
	})

	t.Run("fails with invalid assertion data", func(t *testing.T) {
		t.Parallel()
		service, queries, ctx, userID := setupWebAuthnService(t)

		// Add a credential
		name := "Test Key"
		_, err := queries.CreateWebAuthnCredential(ctx, db.CreateWebAuthnCredentialParams{
			UserID:          userID,
			CredentialID:    []byte("login-finish-test-cred"),
			PublicKey:       []byte("pub-key-data"),
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

		response := &protocol.ParsedCredentialAssertionData{
			ParsedPublicKeyCredential: protocol.ParsedPublicKeyCredential{
				ParsedCredential: protocol.ParsedCredential{
					ID:   "fake-credential-id",
					Type: "public-key",
				},
				RawID: []byte("fake-credential-id"),
			},
			Response: protocol.ParsedAssertionResponse{
				CollectedClientData: protocol.CollectedClientData{
					Type:      protocol.AssertCeremony,
					Challenge: "wrong-challenge",
					Origin:    "http://localhost:3000",
				},
			},
		}

		sessionData := webauthn.SessionData{
			Challenge: "test-challenge-that-doesnt-match",
			UserID:    userID[:],
		}

		err = service.FinishLogin(ctx, userID, "testuser", "Test User", response, sessionData)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to validate login")
	})
}

// Module tests

func TestNewTOTPServiceFromConfig(t *testing.T) {
	t.Parallel()

	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())
	logger := zaptest.NewLogger(t)

	// Create an encryption key
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	encryptor, err := crypto.NewEncryptor(key)
	require.NoError(t, err)

	cfg := &config.Config{}

	service := NewTOTPServiceFromConfig(queries, encryptor, logger, cfg)
	assert.NotNil(t, service)
	assert.Equal(t, "Revenge", service.issuer)
}

// Session cache tests

func TestWebAuthnService_HasCache(t *testing.T) {
	t.Parallel()

	logger := zaptest.NewLogger(t)

	t.Run("returns false when no cache configured", func(t *testing.T) {
		service, err := NewWebAuthnService(nil, logger, nil, "Test", "localhost", []string{"http://localhost"})
		require.NoError(t, err)
		assert.False(t, service.HasCache())
	})

	t.Run("returns true when cache configured", func(t *testing.T) {
		sessionCache, err := cache.NewNamedCache(nil, 100, 5*time.Minute, "webauthn-test")
		require.NoError(t, err)
		defer sessionCache.Close()

		service, err := NewWebAuthnService(nil, logger, sessionCache, "Test", "localhost", []string{"http://localhost"})
		require.NoError(t, err)
		assert.True(t, service.HasCache())
	})
}

func TestWebAuthnService_SessionCache(t *testing.T) {
	t.Parallel()

	logger := zaptest.NewLogger(t)
	ctx := context.Background()

	t.Run("stores and retrieves registration session", func(t *testing.T) {
		sessionCache, err := cache.NewNamedCache(nil, 100, 5*time.Minute, "webauthn-test")
		require.NoError(t, err)
		defer sessionCache.Close()

		testDB := testutil.NewFastTestDB(t)
		queries := db.New(testDB.Pool())
		userID := createTestUserForWebAuthn(t, queries, ctx)

		service, err := NewWebAuthnService(queries, logger, sessionCache, "Test", "localhost", []string{"http://localhost:3000"})
		require.NoError(t, err)

		// Begin registration (stores session in cache)
		_, err = service.BeginRegistration(ctx, userID, "testuser", "Test User")
		require.NoError(t, err)

		// Retrieve session from cache
		session, err := service.GetRegistrationSession(ctx, userID)
		require.NoError(t, err)
		assert.NotNil(t, session)
		assert.NotEmpty(t, session.Challenge)

		// Delete session
		service.DeleteRegistrationSession(ctx, userID)

		// Session should be gone
		_, err = service.GetRegistrationSession(ctx, userID)
		assert.Error(t, err)
	})

	t.Run("stores and retrieves login session", func(t *testing.T) {
		sessionCache, err := cache.NewNamedCache(nil, 100, 5*time.Minute, "webauthn-test")
		require.NoError(t, err)
		defer sessionCache.Close()

		testDB := testutil.NewFastTestDB(t)
		queries := db.New(testDB.Pool())
		userID := createTestUserForWebAuthn(t, queries, ctx)

		// Add a credential first
		name := "Test Key"
		_, err = queries.CreateWebAuthnCredential(ctx, db.CreateWebAuthnCredentialParams{
			UserID:          userID,
			CredentialID:    []byte("session-test-cred-id"),
			PublicKey:       []byte("pub-key"),
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

		service, err := NewWebAuthnService(queries, logger, sessionCache, "Test", "localhost", []string{"http://localhost:3000"})
		require.NoError(t, err)

		// Begin login (stores session in cache)
		_, err = service.BeginLogin(ctx, userID, "testuser", "Test User")
		require.NoError(t, err)

		// Retrieve session from cache
		session, err := service.GetLoginSession(ctx, userID)
		require.NoError(t, err)
		assert.NotNil(t, session)
		assert.NotEmpty(t, session.Challenge)

		// Delete session
		service.DeleteLoginSession(ctx, userID)

		// Session should be gone
		_, err = service.GetLoginSession(ctx, userID)
		assert.Error(t, err)
	})

	t.Run("GetRegistrationSession fails without cache", func(t *testing.T) {
		service, err := NewWebAuthnService(nil, logger, nil, "Test", "localhost", []string{"http://localhost"})
		require.NoError(t, err)

		_, err = service.GetRegistrationSession(ctx, uuid.Must(uuid.NewV7()))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cache not configured")
	})

	t.Run("GetLoginSession fails without cache", func(t *testing.T) {
		service, err := NewWebAuthnService(nil, logger, nil, "Test", "localhost", []string{"http://localhost"})
		require.NoError(t, err)

		_, err = service.GetLoginSession(ctx, uuid.Must(uuid.NewV7()))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cache not configured")
	})
}

func TestNewWebAuthnServiceFromConfig(t *testing.T) {
	t.Parallel()

	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())
	logger := zaptest.NewLogger(t)

	t.Run("uses localhost when host is empty", func(t *testing.T) {
		cfg := &config.Config{
			Server: config.ServerConfig{
				Host: "",
				Port: 3000,
			},
		}

		service, err := NewWebAuthnServiceFromConfig(queries, logger, cfg, nil)
		require.NoError(t, err)
		assert.NotNil(t, service)
	})

	t.Run("uses localhost when host is 0.0.0.0", func(t *testing.T) {
		cfg := &config.Config{
			Server: config.ServerConfig{
				Host: "0.0.0.0",
				Port: 3000,
			},
		}

		service, err := NewWebAuthnServiceFromConfig(queries, logger, cfg, nil)
		require.NoError(t, err)
		assert.NotNil(t, service)
	})

	t.Run("uses configured host", func(t *testing.T) {
		cfg := &config.Config{
			Server: config.ServerConfig{
				Host: "revenge.example.com",
				Port: 443,
			},
		}

		service, err := NewWebAuthnServiceFromConfig(queries, logger, cfg, nil)
		require.NoError(t, err)
		assert.NotNil(t, service)
	})

	t.Run("with cache client creates session cache", func(t *testing.T) {
		cfg := &config.Config{
			Server: config.ServerConfig{
				Host: "localhost",
				Port: 3000,
			},
		}

		// Create a disabled cache client (no Redis needed)
		cacheClient, err := cache.NewClient(&config.Config{
			Cache: config.CacheConfig{Enabled: false},
		}, slog.Default())
		require.NoError(t, err)

		service, err := NewWebAuthnServiceFromConfig(queries, logger, cfg, cacheClient)
		require.NoError(t, err)
		assert.NotNil(t, service)
		assert.True(t, service.HasCache(), "should have session cache when cacheClient is provided")
	})
}

