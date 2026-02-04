package mfa

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
	"go.uber.org/zap"

	db "github.com/lusoris/revenge/internal/infra/database/db"
)

var (
	ErrCredentialNotFound    = errors.New("webauthn credential not found")
	ErrCloneDetected         = errors.New("webauthn authenticator clone detected")
	ErrInvalidCounter        = errors.New("sign counter did not increment")
	ErrNoCredentials         = errors.New("user has no webauthn credentials")
	ErrCredentialAlreadyUsed = errors.New("credential ID already registered")
)

// Safe conversion helpers to prevent integer overflow

// safeUint32ToInt32 safely converts uint32 to int32, capping at max int32
func safeUint32ToInt32(val uint32) int32 {
	const maxInt32 = 2147483647
	if val > maxInt32 {
		return maxInt32
	}
	return int32(val) // #nosec G115 -- validated above
}

// safeInt32ToUint32 safely converts int32 to uint32, treating negative as 0
func safeInt32ToUint32(val int32) uint32 {
	if val < 0 {
		return 0
	}
	return uint32(val) // #nosec G115 -- validated above
}

// WebAuthnService handles WebAuthn credential management and authentication.
type WebAuthnService struct {
	queries  *db.Queries
	logger   *zap.Logger
	webAuthn *webauthn.WebAuthn
}

// NewWebAuthnService creates a new WebAuthn service.
func NewWebAuthnService(
	queries *db.Queries,
	logger *zap.Logger,
	rpDisplayName string,
	rpID string,
	rpOrigins []string,
) (*WebAuthnService, error) {
	wconfig := &webauthn.Config{
		RPDisplayName: rpDisplayName, // e.g., "Revenge"
		RPID:          rpID,          // e.g., "revenge.example.com"
		RPOrigins:     rpOrigins,     // e.g., ["https://revenge.example.com"]
		Timeouts: webauthn.TimeoutsConfig{
			Login: webauthn.TimeoutConfig{
				Enforce:    true,
				Timeout:    60 * time.Second,
				TimeoutUVD: 60 * time.Second,
			},
			Registration: webauthn.TimeoutConfig{
				Enforce:    true,
				Timeout:    60 * time.Second,
				TimeoutUVD: 60 * time.Second,
			},
		},
	}

	wa, err := webauthn.New(wconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create webauthn instance: %w", err)
	}

	return &WebAuthnService{
		queries:  queries,
		logger:   logger,
		webAuthn: wa,
	}, nil
}

// WebAuthnUser implements the webauthn.User interface for our user model.
type WebAuthnUser struct {
	ID          []byte
	Name        string
	DisplayName string
	Credentials []webauthn.Credential
}

// WebAuthnID returns the user's ID in bytes (required by webauthn.User interface).
func (u *WebAuthnUser) WebAuthnID() []byte {
	return u.ID
}

// WebAuthnName returns the user's username (required by webauthn.User interface).
func (u *WebAuthnUser) WebAuthnName() string {
	return u.Name
}

// WebAuthnDisplayName returns the user's display name (required by webauthn.User interface).
func (u *WebAuthnUser) WebAuthnDisplayName() string {
	return u.DisplayName
}

// WebAuthnCredentials returns the user's credentials (required by webauthn.User interface).
func (u *WebAuthnUser) WebAuthnCredentials() []webauthn.Credential {
	return u.Credentials
}

// WebAuthnIcon returns the user's icon URL (deprecated but required by webauthn.User interface).
func (u *WebAuthnUser) WebAuthnIcon() string {
	return ""
}

// BeginRegistration starts the WebAuthn registration ceremony.
// Returns credential creation options to be sent to the client.
func (s *WebAuthnService) BeginRegistration(
	ctx context.Context,
	userID uuid.UUID,
	username string,
	displayName string,
) (*protocol.CredentialCreation, error) {
	// Get existing credentials to exclude them from re-registration
	existingCreds, err := s.queries.ListWebAuthnCredentials(ctx, userID)
	if err != nil {
		s.logger.Error("failed to list existing credentials",
			zap.String("user_id", userID.String()),
			zap.Error(err))
		// Continue without exclusions
		existingCreds = []db.WebauthnCredential{}
	}

	// Convert existing credentials to webauthn format
	credentials := make([]webauthn.Credential, 0, len(existingCreds))
	for _, cred := range existingCreds {
		credentials = append(credentials, webauthn.Credential{
			ID:              cred.CredentialID,
			PublicKey:       cred.PublicKey,
			AttestationType: cred.AttestationType,
			Transport:       convertTransportsFromDB(cred.Transports),
			Flags: webauthn.CredentialFlags{
				UserPresent:    cred.UserPresent,
				UserVerified:   cred.UserVerified,
				BackupEligible: cred.BackupEligible,
				BackupState:    cred.BackupState,
			},
			Authenticator: webauthn.Authenticator{
				AAGUID:    cred.Aaguid,
				SignCount: safeInt32ToUint32(cred.SignCount),
			},
		})
	}

	user := &WebAuthnUser{
		ID:          userID[:],
		Name:        username,
		DisplayName: displayName,
		Credentials: credentials,
	}

	options, session, err := s.webAuthn.BeginRegistration(user)
	if err != nil {
		return nil, fmt.Errorf("failed to begin registration: %w", err)
	}

	// TODO: Store session in cache (Redis/Dragonfly) with 5min TTL
	// For now, we'll rely on the client to store the session data
	// In production, you should:
	// - Marshal session to JSON
	// - Store in cache with key: webauthn:registration:{userID}
	// - Set TTL to 5 minutes
	_ = session

	s.logger.Info("webauthn registration started",
		zap.String("user_id", userID.String()),
		zap.String("challenge", fmt.Sprintf("%x", options.Response.Challenge)))

	return options, nil
}

// FinishRegistration completes the WebAuthn registration ceremony.
// Validates the credential and stores it in the database.
func (s *WebAuthnService) FinishRegistration(
	ctx context.Context,
	userID uuid.UUID,
	username string,
	displayName string,
	response *protocol.ParsedCredentialCreationData,
	sessionData webauthn.SessionData,
	credentialName string,
) error {
	// Get existing credentials
	existingCreds, _ := s.queries.ListWebAuthnCredentials(ctx, userID)
	credentials := make([]webauthn.Credential, 0, len(existingCreds))
	for _, cred := range existingCreds {
		credentials = append(credentials, webauthn.Credential{
			ID:              cred.CredentialID,
			PublicKey:       cred.PublicKey,
			AttestationType: cred.AttestationType,
			Transport:       convertTransportsFromDB(cred.Transports),
			Flags: webauthn.CredentialFlags{
				UserPresent:    cred.UserPresent,
				UserVerified:   cred.UserVerified,
				BackupEligible: cred.BackupEligible,
				BackupState:    cred.BackupState,
			},
			Authenticator: webauthn.Authenticator{
				AAGUID:    cred.Aaguid,
				SignCount: safeInt32ToUint32(cred.SignCount),
			},
		})
	}

	user := &WebAuthnUser{
		ID:          userID[:],
		Name:        username,
		DisplayName: displayName,
		Credentials: credentials,
	}

	credential, err := s.webAuthn.CreateCredential(user, sessionData, response)
	if err != nil {
		return fmt.Errorf("failed to create credential: %w", err)
	}

	// Check if credential ID already exists
	existing, err := s.queries.GetWebAuthnCredentialByCredentialID(ctx, credential.ID)
	if err == nil && existing.ID != uuid.Nil {
		return ErrCredentialAlreadyUsed
	}

	// Store credential in database
	now := time.Now()
	namePtr := func() *string {
		var name string
		if credentialName != "" {
			name = credentialName
		} else {
			name = fmt.Sprintf("Authenticator registered on %s", now.Format("2006-01-02"))
		}
		return &name
	}()

	_, err = s.queries.CreateWebAuthnCredential(ctx, db.CreateWebAuthnCredentialParams{
		UserID:          userID,
		CredentialID:    credential.ID,
		PublicKey:       credential.PublicKey,
		AttestationType: credential.AttestationType,
		Transports:      convertTransportsToDB(credential.Transport),
		BackupEligible:  credential.Flags.BackupEligible,
		BackupState:     credential.Flags.BackupState,
		UserPresent:     credential.Flags.UserPresent,
		UserVerified:    credential.Flags.UserVerified,
		Aaguid:          credential.Authenticator.AAGUID,
		Name:            namePtr,
	})
	if err != nil {
		return fmt.Errorf("failed to store credential: %w", err)
	}

	s.logger.Info("webauthn credential registered",
		zap.String("user_id", userID.String()),
		zap.String("credential_id", fmt.Sprintf("%x", credential.ID)),
		zap.String("attestation_type", credential.AttestationType))

	return nil
}

// BeginLogin starts the WebAuthn authentication ceremony.
// Returns credential request options to be sent to the client.
func (s *WebAuthnService) BeginLogin(
	ctx context.Context,
	userID uuid.UUID,
	username string,
	displayName string,
) (*protocol.CredentialAssertion, error) {
	// Get user's credentials
	existingCreds, err := s.queries.ListWebAuthnCredentials(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list credentials: %w", err)
	}

	if len(existingCreds) == 0 {
		return nil, ErrNoCredentials
	}

	// Filter out credentials with clone detection flags
	validCreds := make([]db.WebauthnCredential, 0, len(existingCreds))
	for _, cred := range existingCreds {
		if !cred.CloneDetected {
			validCreds = append(validCreds, cred)
		}
	}

	if len(validCreds) == 0 {
		return nil, errors.New("all credentials marked as cloned")
	}

	// Convert credentials to webauthn format
	credentials := make([]webauthn.Credential, 0, len(validCreds))
	for _, cred := range validCreds {
		credentials = append(credentials, webauthn.Credential{
			ID:              cred.CredentialID,
			PublicKey:       cred.PublicKey,
			AttestationType: cred.AttestationType,
			Transport:       convertTransportsFromDB(cred.Transports),
			Flags: webauthn.CredentialFlags{
				UserPresent:    cred.UserPresent,
				UserVerified:   cred.UserVerified,
				BackupEligible: cred.BackupEligible,
				BackupState:    cred.BackupState,
			},
			Authenticator: webauthn.Authenticator{
				AAGUID:    cred.Aaguid,
				SignCount: safeInt32ToUint32(cred.SignCount),
			},
		})
	}

	user := &WebAuthnUser{
		ID:          userID[:],
		Name:        username,
		DisplayName: displayName,
		Credentials: credentials,
	}

	options, session, err := s.webAuthn.BeginLogin(user)
	if err != nil {
		return nil, fmt.Errorf("failed to begin login: %w", err)
	}

	// TODO: Store session in cache (Redis/Dragonfly) with 5min TTL
	_ = session

	s.logger.Info("webauthn login started",
		zap.String("user_id", userID.String()),
		zap.Int("credentials_count", len(credentials)))

	return options, nil
}

// FinishLogin completes the WebAuthn authentication ceremony.
// Validates the assertion and updates the credential's sign counter.
func (s *WebAuthnService) FinishLogin(
	ctx context.Context,
	userID uuid.UUID,
	username string,
	displayName string,
	response *protocol.ParsedCredentialAssertionData,
	sessionData webauthn.SessionData,
) error {
	// Get user's credentials
	existingCreds, err := s.queries.ListWebAuthnCredentials(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to list credentials: %w", err)
	}

	credentials := make([]webauthn.Credential, 0, len(existingCreds))
	for _, cred := range existingCreds {
		credentials = append(credentials, webauthn.Credential{
			ID:              cred.CredentialID,
			PublicKey:       cred.PublicKey,
			AttestationType: cred.AttestationType,
			Transport:       convertTransportsFromDB(cred.Transports),
			Flags: webauthn.CredentialFlags{
				UserPresent:    cred.UserPresent,
				UserVerified:   cred.UserVerified,
				BackupEligible: cred.BackupEligible,
				BackupState:    cred.BackupState,
			},
			Authenticator: webauthn.Authenticator{
				AAGUID:    cred.Aaguid,
				SignCount: safeInt32ToUint32(cred.SignCount),
			},
		})
	}

	user := &WebAuthnUser{
		ID:          userID[:],
		Name:        username,
		DisplayName: displayName,
		Credentials: credentials,
	}

	credential, err := s.webAuthn.ValidateLogin(user, sessionData, response)
	if err != nil {
		return fmt.Errorf("failed to validate login: %w", err)
	}

	// Get the credential from DB to check sign counter
	dbCred, err := s.queries.GetWebAuthnCredentialByCredentialID(ctx, credential.ID)
	if err != nil {
		return fmt.Errorf("failed to get credential: %w", err)
	}

	// Clone detection: sign counter must always increment
	newCounter := credential.Authenticator.SignCount
	oldCounter := safeInt32ToUint32(dbCred.SignCount)

	if newCounter <= oldCounter {
		// Sign counter did not increment - possible clone!
		s.logger.Warn("webauthn clone detected",
			zap.String("user_id", userID.String()),
			zap.String("credential_id", fmt.Sprintf("%x", credential.ID)),
			zap.Uint32("old_counter", oldCounter),
			zap.Uint32("new_counter", newCounter))

		// Mark credential as cloned
		err = s.queries.MarkWebAuthnCloneDetected(ctx, credential.ID)
		if err != nil {
			s.logger.Error("failed to mark credential as cloned", zap.Error(err))
		}

		return ErrCloneDetected
	}

	// Update sign counter
	err = s.queries.UpdateWebAuthnCounter(ctx, db.UpdateWebAuthnCounterParams{
		CredentialID: credential.ID,
		SignCount:    safeUint32ToInt32(newCounter),
	})
	if err != nil {
		return fmt.Errorf("failed to update sign counter: %w", err)
	}

	s.logger.Info("webauthn login successful",
		zap.String("user_id", userID.String()),
		zap.String("credential_id", fmt.Sprintf("%x", credential.ID)),
		zap.Uint32("new_counter", newCounter))

	return nil
}

// ListCredentials returns all WebAuthn credentials for a user.
func (s *WebAuthnService) ListCredentials(ctx context.Context, userID uuid.UUID) ([]db.WebauthnCredential, error) {
	return s.queries.ListWebAuthnCredentials(ctx, userID)
}

// DeleteCredential removes a WebAuthn credential.
func (s *WebAuthnService) DeleteCredential(ctx context.Context, userID uuid.UUID, credentialID uuid.UUID) error {
	return s.queries.DeleteWebAuthnCredential(ctx, db.DeleteWebAuthnCredentialParams{
		ID:     credentialID,
		UserID: userID,
	})
}

// RenameCredential updates the user-facing name of a credential.
func (s *WebAuthnService) RenameCredential(ctx context.Context, credentialID uuid.UUID, newName string) error {
	return s.queries.UpdateWebAuthnCredentialName(ctx, db.UpdateWebAuthnCredentialNameParams{
		ID:   credentialID,
		Name: &newName,
	})
}

// HasWebAuthn checks if a user has any WebAuthn credentials.
func (s *WebAuthnService) HasWebAuthn(ctx context.Context, userID uuid.UUID) (bool, error) {
	count, err := s.queries.CountWebAuthnCredentials(ctx, userID)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Helper functions

func convertTransportsToDB(transports []protocol.AuthenticatorTransport) []string {
	result := make([]string, len(transports))
	for i, t := range transports {
		result[i] = string(t)
	}
	return result
}

func convertTransportsFromDB(transports []string) []protocol.AuthenticatorTransport {
	result := make([]protocol.AuthenticatorTransport, len(transports))
	for i, t := range transports {
		result[i] = protocol.AuthenticatorTransport(t)
	}
	return result
}

// SessionDataToJSON serializes session data for caching.
func SessionDataToJSON(session webauthn.SessionData) ([]byte, error) {
	return json.Marshal(session)
}

// SessionDataFromJSON deserializes session data from cache.
func SessionDataFromJSON(data []byte) (*webauthn.SessionData, error) {
	var session webauthn.SessionData
	err := json.Unmarshal(data, &session)
	if err != nil {
		return nil, err
	}
	return &session, nil
}
