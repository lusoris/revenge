package oidc

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

// Errors
var (
	ErrProviderNotFound    = errors.New("OIDC provider not found")
	ErrProviderDisabled    = errors.New("OIDC provider is disabled")
	ErrUserLinkNotFound    = errors.New("OIDC user link not found")
	ErrUserLinkExists      = errors.New("user already linked to this provider")
	ErrStateNotFound       = errors.New("OAuth state not found")
	ErrStateExpired        = errors.New("OAuth state expired")
	ErrInvalidState        = errors.New("invalid OAuth state")
	ErrInvalidCode         = errors.New("invalid authorization code")
	ErrTokenExchange       = errors.New("failed to exchange token")
	ErrUserInfoFetch       = errors.New("failed to fetch user info")
	ErrAutoCreateDisabled  = errors.New("auto-create users is disabled for this provider")
	ErrLinkingDisabled     = errors.New("account linking is disabled for this provider")
	ErrProviderNameExists  = errors.New("provider with this name already exists")
	ErrInvalidProviderType = errors.New("invalid provider type")
	ErrInvalidIssuerURL    = errors.New("invalid issuer URL")
	ErrDiscoveryFailed     = errors.New("OIDC discovery failed")
	ErrInvalidCallbackURL  = errors.New("invalid callback URL")
	ErrEncryptionFailed    = errors.New("failed to encrypt sensitive data")
	ErrDecryptionFailed    = errors.New("failed to decrypt sensitive data")
)

// Constants
const (
	StateExpiry           = 10 * time.Minute
	CodeVerifierLen       = 64
	StateLen              = 32
	ProviderTypeGeneric   = "generic"
	ProviderTypeAuthentik = "authentik"
	ProviderTypeKeycloak  = "keycloak"
)

// Service implements OIDC business logic
type Service struct {
	repo        Repository
	logger      *zap.Logger
	callbackURL string
	encryptKey  []byte // For encrypting client secrets and tokens
}

// NewService creates a new OIDC service
func NewService(repo Repository, logger *zap.Logger, callbackURL string, encryptKey []byte) *Service {
	return &Service{
		repo:        repo,
		logger:      logger,
		callbackURL: callbackURL,
		encryptKey:  encryptKey,
	}
}

// ============================================================================
// Provider Management
// ============================================================================

// AddProvider creates a new OIDC provider
func (s *Service) AddProvider(ctx context.Context, req CreateProviderRequest) (*Provider, error) {
	// Validate provider type
	if !isValidProviderType(req.ProviderType) {
		return nil, ErrInvalidProviderType
	}

	// Validate issuer URL
	if _, err := url.Parse(req.IssuerURL); err != nil {
		return nil, ErrInvalidIssuerURL
	}

	// Check if provider name already exists
	existing, err := s.repo.GetProviderByName(ctx, req.Name)
	if err == nil && existing != nil {
		return nil, ErrProviderNameExists
	}
	if err != nil && !errors.Is(err, ErrProviderNotFound) {
		return nil, err
	}

	// Encrypt client secret
	encrypted, err := s.encryptSecret(req.ClientSecretEncrypted)
	if err != nil {
		return nil, ErrEncryptionFailed
	}
	req.ClientSecretEncrypted = encrypted

	// Set default scopes if empty
	if len(req.Scopes) == 0 {
		req.Scopes = []string{"openid", "profile", "email"}
	}

	// Set default claim mappings
	if req.ClaimMappings.Username == "" {
		req.ClaimMappings.Username = "preferred_username"
	}
	if req.ClaimMappings.Email == "" {
		req.ClaimMappings.Email = "email"
	}
	if req.ClaimMappings.Name == "" {
		req.ClaimMappings.Name = "name"
	}
	if req.ClaimMappings.Picture == "" {
		req.ClaimMappings.Picture = "picture"
	}

	return s.repo.CreateProvider(ctx, req)
}

// GetProvider gets a provider by ID
func (s *Service) GetProvider(ctx context.Context, id uuid.UUID) (*Provider, error) {
	return s.repo.GetProvider(ctx, id)
}

// GetProviderByName gets a provider by name
func (s *Service) GetProviderByName(ctx context.Context, name string) (*Provider, error) {
	return s.repo.GetProviderByName(ctx, name)
}

// GetDefaultProvider gets the default provider
func (s *Service) GetDefaultProvider(ctx context.Context) (*Provider, error) {
	return s.repo.GetDefaultProvider(ctx)
}

// ListProviders lists all providers
func (s *Service) ListProviders(ctx context.Context) ([]Provider, error) {
	return s.repo.ListProviders(ctx)
}

// ListEnabledProviders lists all enabled providers
func (s *Service) ListEnabledProviders(ctx context.Context) ([]Provider, error) {
	return s.repo.ListEnabledProviders(ctx)
}

// UpdateProvider updates a provider
func (s *Service) UpdateProvider(ctx context.Context, id uuid.UUID, req UpdateProviderRequest) (*Provider, error) {
	// Validate provider type if provided
	if req.ProviderType != nil && !isValidProviderType(*req.ProviderType) {
		return nil, ErrInvalidProviderType
	}

	// Validate issuer URL if provided
	if req.IssuerURL != nil {
		if _, err := url.Parse(*req.IssuerURL); err != nil {
			return nil, ErrInvalidIssuerURL
		}
	}

	// Encrypt client secret if provided
	if req.ClientSecretEncrypted != nil {
		encrypted, err := s.encryptSecret(req.ClientSecretEncrypted)
		if err != nil {
			return nil, ErrEncryptionFailed
		}
		req.ClientSecretEncrypted = encrypted
	}

	return s.repo.UpdateProvider(ctx, id, req)
}

// DeleteProvider deletes a provider
func (s *Service) DeleteProvider(ctx context.Context, id uuid.UUID) error {
	// Delete any pending states for this provider
	if err := s.repo.DeleteStatesByProvider(ctx, id); err != nil {
		s.logger.Warn("failed to delete states for provider", zap.String("provider_id", id.String()), zap.Error(err))
	}
	return s.repo.DeleteProvider(ctx, id)
}

// EnableProvider enables a provider
func (s *Service) EnableProvider(ctx context.Context, id uuid.UUID) error {
	return s.repo.EnableProvider(ctx, id)
}

// DisableProvider disables a provider
func (s *Service) DisableProvider(ctx context.Context, id uuid.UUID) error {
	// Delete any pending states for this provider
	if err := s.repo.DeleteStatesByProvider(ctx, id); err != nil {
		s.logger.Warn("failed to delete states for disabled provider", zap.String("provider_id", id.String()), zap.Error(err))
	}
	return s.repo.DisableProvider(ctx, id)
}

// SetDefaultProvider sets a provider as default
func (s *Service) SetDefaultProvider(ctx context.Context, id uuid.UUID) error {
	return s.repo.SetDefaultProvider(ctx, id)
}

// ============================================================================
// OAuth2 Flow
// ============================================================================

// AuthURLResult contains the authorization URL and state
type AuthURLResult struct {
	URL   string
	State string
}

// GetAuthURL generates an authorization URL for a provider
func (s *Service) GetAuthURL(ctx context.Context, providerName string, redirectURL string, userID *uuid.UUID) (*AuthURLResult, error) {
	// Get provider
	provider, err := s.repo.GetProviderByName(ctx, providerName)
	if err != nil {
		return nil, err
	}
	if !provider.IsEnabled {
		return nil, ErrProviderDisabled
	}

	// Generate state and PKCE verifier
	state, err := generateRandomString(StateLen)
	if err != nil {
		return nil, err
	}
	codeVerifier, err := generateRandomString(CodeVerifierLen)
	if err != nil {
		return nil, err
	}

	// Store state
	_, err = s.repo.CreateState(ctx, CreateStateRequest{
		State:        state,
		CodeVerifier: &codeVerifier,
		ProviderID:   provider.ID,
		UserID:       userID,
		RedirectURL:  &redirectURL,
		ExpiresAt:    time.Now().Add(StateExpiry),
	})
	if err != nil {
		return nil, err
	}

	// Build OAuth2 config
	oauth2Config := s.buildOAuth2Config(provider)

	// Generate code challenge (S256)
	codeChallenge := generateCodeChallenge(codeVerifier)

	// Generate auth URL with PKCE
	authURL := oauth2Config.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)

	return &AuthURLResult{
		URL:   authURL,
		State: state,
	}, nil
}

// CallbackResult contains the result of handling an OAuth callback
type CallbackResult struct {
	UserLink    *UserLink
	UserID      uuid.UUID
	IsNewUser   bool
	AccessToken string
	IDToken     string
	UserInfo    *UserInfo
}

// UserInfo contains user information from the OIDC provider
type UserInfo struct {
	Subject  string
	Email    string
	Name     string
	Username string
	Picture  string
	Roles    []string
	Claims   map[string]any
}

// HandleCallback handles the OAuth callback
func (s *Service) HandleCallback(ctx context.Context, stateParam, code string) (*CallbackResult, error) {
	// Validate state
	state, err := s.repo.GetState(ctx, stateParam)
	if err != nil {
		return nil, ErrInvalidState
	}

	// Delete state immediately (one-time use)
	defer func() {
		_ = s.repo.DeleteState(ctx, stateParam)
	}()

	// Check expiration
	if time.Now().After(state.ExpiresAt) {
		return nil, ErrStateExpired
	}

	// Get provider
	provider, err := s.repo.GetProvider(ctx, state.ProviderID)
	if err != nil {
		return nil, err
	}
	if !provider.IsEnabled {
		return nil, ErrProviderDisabled
	}

	// Create OIDC provider and verifier
	oidcProvider, err := oidc.NewProvider(ctx, provider.IssuerURL)
	if err != nil {
		s.logger.Error("failed to create OIDC provider", zap.String("issuer", provider.IssuerURL), zap.Error(err))
		return nil, ErrDiscoveryFailed
	}

	// Build OAuth2 config
	oauth2Config := s.buildOAuth2Config(provider)

	// Exchange code for tokens
	var tokenOpts []oauth2.AuthCodeOption
	if state.CodeVerifier != nil {
		tokenOpts = append(tokenOpts, oauth2.SetAuthURLParam("code_verifier", *state.CodeVerifier))
	}

	token, err := oauth2Config.Exchange(ctx, code, tokenOpts...)
	if err != nil {
		s.logger.Error("token exchange failed", zap.Error(err))
		return nil, ErrTokenExchange
	}

	// Extract ID token
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("no id_token in token response")
	}

	// Verify ID token
	verifier := oidcProvider.Verifier(&oidc.Config{ClientID: provider.ClientID})
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		s.logger.Error("ID token verification failed", zap.Error(err))
		return nil, fmt.Errorf("invalid id_token: %w", err)
	}

	// Extract claims
	var claims map[string]any
	if err := idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to extract claims: %w", err)
	}

	// Build user info
	userInfo := s.extractUserInfo(provider, claims)
	userInfo.Subject = idToken.Subject

	// Check if user is already linked
	link, err := s.repo.GetUserLinkBySubject(ctx, provider.ID, userInfo.Subject)
	if err != nil && !errors.Is(err, ErrUserLinkNotFound) {
		return nil, err
	}

	result := &CallbackResult{
		AccessToken: token.AccessToken,
		IDToken:     rawIDToken,
		UserInfo:    userInfo,
	}

	if link != nil {
		// Existing user - update link and return
		updateReq := UpdateUserLinkRequest{
			Email:      &userInfo.Email,
			Name:       &userInfo.Name,
			PictureURL: &userInfo.Picture,
		}
		now := time.Now()
		updateReq.LastLoginAt = &now

		if token.AccessToken != "" {
			encrypted, err := s.encryptSecret([]byte(token.AccessToken))
			if err == nil {
				updateReq.AccessTokenEncrypted = encrypted
			}
		}
		if token.RefreshToken != "" {
			encrypted, err := s.encryptSecret([]byte(token.RefreshToken))
			if err == nil {
				updateReq.RefreshTokenEncrypted = encrypted
			}
		}
		if !token.Expiry.IsZero() {
			updateReq.TokenExpiresAt = &token.Expiry
		}

		updatedLink, err := s.repo.UpdateUserLink(ctx, link.ID, updateReq)
		if err != nil {
			s.logger.Warn("failed to update user link", zap.Error(err))
		} else {
			result.UserLink = updatedLink
		}

		result.UserID = link.UserID
		result.IsNewUser = false
		return result, nil
	}

	// No existing link
	// If this is a linking flow (userID provided in state), create the link
	if state.UserID != nil {
		if !provider.AllowLinking {
			return nil, ErrLinkingDisabled
		}
		link, err := s.createUserLink(ctx, provider, *state.UserID, userInfo, token)
		if err != nil {
			return nil, err
		}
		result.UserLink = link
		result.UserID = *state.UserID
		result.IsNewUser = false
		return result, nil
	}

	// New user - auto-create must be enabled
	if !provider.AutoCreateUsers {
		return nil, ErrAutoCreateDisabled
	}

	// Return info for user creation (caller must create user and link)
	result.IsNewUser = true
	return result, nil
}

// LinkUser links an existing user to an OIDC provider
func (s *Service) LinkUser(ctx context.Context, userID uuid.UUID, providerID uuid.UUID, subject string, userInfo *UserInfo, token *oauth2.Token) (*UserLink, error) {
	provider, err := s.repo.GetProvider(ctx, providerID)
	if err != nil {
		return nil, err
	}
	if !provider.AllowLinking {
		return nil, ErrLinkingDisabled
	}

	// Check if already linked
	existing, err := s.repo.GetUserLinkByUserAndProvider(ctx, userID, providerID)
	if err == nil && existing != nil {
		return nil, ErrUserLinkExists
	}
	if err != nil && !errors.Is(err, ErrUserLinkNotFound) {
		return nil, err
	}

	return s.createUserLink(ctx, provider, userID, userInfo, token)
}

// UnlinkUser unlinks a user from an OIDC provider
func (s *Service) UnlinkUser(ctx context.Context, userID uuid.UUID, providerID uuid.UUID) error {
	return s.repo.DeleteUserLinkByUserAndProvider(ctx, userID, providerID)
}

// ListUserLinks lists all OIDC links for a user
func (s *Service) ListUserLinks(ctx context.Context, userID uuid.UUID) ([]UserLinkWithProvider, error) {
	return s.repo.ListUserLinks(ctx, userID)
}

// ============================================================================
// Helpers
// ============================================================================

func (s *Service) buildOAuth2Config(provider *Provider) *oauth2.Config {
	callbackURL := s.callbackURL
	if !strings.Contains(callbackURL, "/callback/") {
		// Append provider name to callback URL
		callbackURL = strings.TrimSuffix(callbackURL, "/") + "/" + provider.Name
	}

	config := &oauth2.Config{
		ClientID:     provider.ClientID,
		ClientSecret: string(s.decryptSecret(provider.ClientSecretEncrypted)),
		RedirectURL:  callbackURL,
		Scopes:       provider.Scopes,
	}

	// Use discovered endpoints or custom endpoints
	if provider.AuthorizationEndpoint != nil && provider.TokenEndpoint != nil {
		config.Endpoint = oauth2.Endpoint{
			AuthURL:  *provider.AuthorizationEndpoint,
			TokenURL: *provider.TokenEndpoint,
		}
	} else {
		// Use standard OIDC endpoints
		config.Endpoint = oauth2.Endpoint{
			AuthURL:  provider.IssuerURL + "/authorize",
			TokenURL: provider.IssuerURL + "/token",
		}
	}

	return config
}

func (s *Service) extractUserInfo(provider *Provider, claims map[string]any) *UserInfo {
	info := &UserInfo{
		Claims: claims,
	}

	// Extract mapped claims
	if v, ok := getClaim(claims, provider.ClaimMappings.Email).(string); ok {
		info.Email = v
	}
	if v, ok := getClaim(claims, provider.ClaimMappings.Name).(string); ok {
		info.Name = v
	}
	if v, ok := getClaim(claims, provider.ClaimMappings.Username).(string); ok {
		info.Username = v
	}
	if v, ok := getClaim(claims, provider.ClaimMappings.Picture).(string); ok {
		info.Picture = v
	}

	// Extract roles
	if provider.ClaimMappings.Roles != "" {
		if roles := getClaim(claims, provider.ClaimMappings.Roles); roles != nil {
			switch r := roles.(type) {
			case []string:
				info.Roles = r
			case []any:
				for _, v := range r {
					if s, ok := v.(string); ok {
						info.Roles = append(info.Roles, s)
					}
				}
			}
		}
	}

	// Map roles if mappings are defined
	if len(provider.RoleMappings) > 0 {
		var mappedRoles []string
		for _, role := range info.Roles {
			if mapped, ok := provider.RoleMappings[role]; ok {
				mappedRoles = append(mappedRoles, mapped)
			}
		}
		info.Roles = mappedRoles
	}

	return info
}

func (s *Service) createUserLink(ctx context.Context, provider *Provider, userID uuid.UUID, userInfo *UserInfo, token *oauth2.Token) (*UserLink, error) {
	req := CreateUserLinkRequest{
		UserID:     userID,
		ProviderID: provider.ID,
		Subject:    userInfo.Subject,
		Email:      &userInfo.Email,
		Name:       &userInfo.Name,
		PictureURL: &userInfo.Picture,
	}

	if token != nil {
		if token.AccessToken != "" {
			encrypted, err := s.encryptSecret([]byte(token.AccessToken))
			if err == nil {
				req.AccessTokenEncrypted = encrypted
			}
		}
		if token.RefreshToken != "" {
			encrypted, err := s.encryptSecret([]byte(token.RefreshToken))
			if err == nil {
				req.RefreshTokenEncrypted = encrypted
			}
		}
		if !token.Expiry.IsZero() {
			req.TokenExpiresAt = &token.Expiry
		}
	}

	return s.repo.CreateUserLink(ctx, req)
}

// CleanupExpiredStates removes expired OAuth states
func (s *Service) CleanupExpiredStates(ctx context.Context) (int64, error) {
	return s.repo.DeleteExpiredStates(ctx)
}

// ============================================================================
// Crypto Helpers
// ============================================================================

// Simple encryption - in production, use a proper KMS or secure encryption
func (s *Service) encryptSecret(plaintext []byte) ([]byte, error) {
	if len(s.encryptKey) == 0 {
		// No encryption configured - return as-is (for dev only)
		return plaintext, nil
	}

	// Create AES cipher
	block, err := aes.NewCipher(s.encryptKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt and prepend nonce
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func (s *Service) decryptSecret(ciphertext []byte) []byte {
	if len(s.encryptKey) == 0 {
		// No encryption configured - return as-is
		return ciphertext
	}

	// Create AES cipher
	block, err := aes.NewCipher(s.encryptKey)
	if err != nil {
		s.logger.Error("failed to create cipher for decryption", zap.Error(err))
		return ciphertext // Fallback to returning as-is
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		s.logger.Error("failed to create GCM for decryption", zap.Error(err))
		return ciphertext
	}

	// Extract nonce
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		s.logger.Error("ciphertext too short")
		return ciphertext
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		s.logger.Error("failed to decrypt secret", zap.Error(err))
		return ciphertext // Fallback
	}

	return plaintext
}

// ============================================================================
// Utility Functions
// ============================================================================

func isValidProviderType(t string) bool {
	switch t {
	case "oidc", ProviderTypeGeneric, ProviderTypeAuthentik, ProviderTypeKeycloak:
		return true
	default:
		return false
	}
}

func generateRandomString(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func generateCodeChallenge(verifier string) string {
	h := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(h[:])
}

func getClaim(claims map[string]any, path string) any {
	parts := strings.Split(path, ".")
	var current any = claims
	for _, part := range parts {
		if m, ok := current.(map[string]any); ok {
			current = m[part]
		} else {
			return nil
		}
	}
	return current
}
