// Package oidc provides OIDC/SSO authentication services for Revenge Go.
package oidc

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/lusoris/revenge/internal/domain"
)

// Service implements OIDC authentication.
type Service struct {
	providers  domain.OIDCProviderRepository
	userLinks  domain.OIDCUserLinkRepository
	users      domain.UserRepository
	sessions   domain.SessionRepository
	tokens     domain.TokenService
	passwords  domain.PasswordService
	httpClient *http.Client

	// In-memory state store (could be moved to Redis for multi-instance)
	stateMu    sync.RWMutex
	stateStore map[string]*authState

	// Discovery cache
	discoveryCacheMu sync.RWMutex
	discoveryCache   map[string]*cachedDiscovery

	// Configuration
	stateExpiry     time.Duration
	accessDuration  time.Duration
	refreshDuration time.Duration
}

type authState struct {
	ProviderID   uuid.UUID
	Nonce        string
	RedirectURI  string
	CodeVerifier string // PKCE
	CreatedAt    time.Time
}

type cachedDiscovery struct {
	Document  *DiscoveryDocument
	FetchedAt time.Time
}

// DiscoveryDocument represents the OIDC discovery document.
type DiscoveryDocument struct {
	Issuer                        string   `json:"issuer"`
	AuthorizationEndpoint         string   `json:"authorization_endpoint"`
	TokenEndpoint                 string   `json:"token_endpoint"`
	UserInfoEndpoint              string   `json:"userinfo_endpoint,omitempty"`
	JwksURI                       string   `json:"jwks_uri"`
	EndSessionEndpoint            string   `json:"end_session_endpoint,omitempty"`
	ScopesSupported               []string `json:"scopes_supported,omitempty"`
	ResponseTypesSupported        []string `json:"response_types_supported,omitempty"`
	CodeChallengeMethodsSupported []string `json:"code_challenge_methods_supported,omitempty"`
}

// ServiceParams holds dependencies for the OIDC service.
type ServiceParams struct {
	Providers       domain.OIDCProviderRepository
	UserLinks       domain.OIDCUserLinkRepository
	Users           domain.UserRepository
	Sessions        domain.SessionRepository
	Tokens          domain.TokenService
	Passwords       domain.PasswordService
	AccessDuration  time.Duration
	RefreshDuration time.Duration
}

// GetProviders returns all enabled OIDC providers.
func (s *Service) GetProviders(ctx context.Context) ([]*domain.OIDCProvider, error) {
	return s.providers.ListEnabled(ctx)
}

// GetAuthorizationURL generates the authorization URL for a provider.
func (s *Service) GetAuthorizationURL(ctx context.Context, providerID uuid.UUID, redirectURI string) (string, error) {
	provider, err := s.providers.GetByID(ctx, providerID)
	if err != nil {
		if errors.Is(err, domain.ErrOIDCProviderNotFound) {
			return "", domain.ErrOIDCProviderNotFound
		}
		return "", fmt.Errorf("failed to get provider: %w", err)
	}

	if !provider.Enabled {
		return "", fmt.Errorf("OIDC provider is disabled")
	}

	// Fetch discovery document
	discovery, err := s.fetchDiscovery(ctx, provider.IssuerURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch discovery document: %w", err)
	}

	// Generate state, nonce, and PKCE code verifier
	state, err := generateRandomString(32)
	if err != nil {
		return "", fmt.Errorf("failed to generate state: %w", err)
	}

	nonce, err := generateRandomString(32)
	if err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	codeVerifier, err := generateRandomString(64)
	if err != nil {
		return "", fmt.Errorf("failed to generate code verifier: %w", err)
	}

	// Store state for callback validation
	s.stateMu.Lock()
	s.stateStore[state] = &authState{
		ProviderID:   providerID,
		Nonce:        nonce,
		RedirectURI:  redirectURI,
		CodeVerifier: codeVerifier,
		CreatedAt:    time.Now(),
	}
	s.stateMu.Unlock()

	// Build authorization URL
	authURL, err := url.Parse(discovery.AuthorizationEndpoint)
	if err != nil {
		return "", fmt.Errorf("invalid authorization endpoint: %w", err)
	}

	// Generate code challenge (PKCE S256)
	codeChallenge := generateCodeChallenge(codeVerifier)

	query := authURL.Query()
	query.Set("client_id", provider.ClientID)
	query.Set("redirect_uri", redirectURI)
	query.Set("response_type", "code")
	query.Set("scope", strings.Join(provider.Scopes, " "))
	query.Set("state", state)
	query.Set("nonce", nonce)
	query.Set("code_challenge", codeChallenge)
	query.Set("code_challenge_method", "S256")
	authURL.RawQuery = query.Encode()

	slog.Info("generated OIDC authorization URL",
		slog.String("provider", provider.Name),
		slog.String("provider_id", providerID.String()))

	return authURL.String(), nil
}

// CallbackParams contains parameters for handling the OIDC callback.
type CallbackParams struct {
	Code          string
	State         string
	RedirectURI   string
	DeviceID      *string
	DeviceName    *string
	ClientName    *string
	ClientVersion *string
}

// HandleCallback processes the OIDC callback and returns auth result.
func (s *Service) HandleCallback(ctx context.Context, params CallbackParams) (*domain.AuthResult, error) {
	// Validate and retrieve state
	s.stateMu.Lock()
	state, exists := s.stateStore[params.State]
	if exists {
		delete(s.stateStore, params.State)
	}
	s.stateMu.Unlock()

	if !exists {
		return nil, fmt.Errorf("invalid or expired state")
	}

	if time.Since(state.CreatedAt) > s.stateExpiry {
		return nil, fmt.Errorf("authentication request expired")
	}

	// Get provider
	provider, err := s.providers.GetByID(ctx, state.ProviderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}

	// Fetch discovery document
	discovery, err := s.fetchDiscovery(ctx, provider.IssuerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch discovery document: %w", err)
	}

	// Exchange code for tokens
	tokenResp, err := s.exchangeCode(ctx, provider, discovery, params.Code, state.RedirectURI, state.CodeVerifier)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange authorization code: %w", err)
	}

	// Parse and validate ID token
	claims, err := s.parseIDToken(tokenResp.IDToken, provider, state.Nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to validate ID token: %w", err)
	}

	// Find or create user
	user, err := s.findOrCreateUser(ctx, provider, claims)
	if err != nil {
		return nil, err
	}

	// Update last login on link
	link, err := s.userLinks.Get(ctx, provider.ID, claims.Subject)
	if err == nil {
		if err := s.userLinks.UpdateLastLogin(ctx, link.ID); err != nil {
			slog.Warn("failed to update OIDC link last login",
				slog.String("link_id", link.ID.String()),
				slog.Any("error", err))
		}
	}

	// Create session
	result, err := s.createSession(ctx, user, params.DeviceID, params.DeviceName, params.ClientName, params.ClientVersion)
	if err != nil {
		return nil, err
	}

	// Update user's last login
	if err := s.users.UpdateLastLogin(ctx, user.ID); err != nil {
		slog.Warn("failed to update user last login",
			slog.String("user_id", user.ID.String()),
			slog.Any("error", err))
	}

	slog.Info("OIDC login successful",
		slog.String("user_id", user.ID.String()),
		slog.String("username", user.Username),
		slog.String("provider", provider.Name))

	return result, nil
}

// IDTokenClaims represents claims from the OIDC ID token.
type IDTokenClaims struct {
	Subject           string
	Email             string
	EmailVerified     bool
	Name              string
	PreferredUsername string
	Groups            []string
}

// TokenResponse represents the OIDC token endpoint response.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IDToken      string `json:"id_token"`
}

func (s *Service) exchangeCode(ctx context.Context, provider *domain.OIDCProvider, discovery *DiscoveryDocument, code, redirectURI, codeVerifier string) (*TokenResponse, error) {
	// Decrypt client secret
	clientSecret := s.decryptClientSecret(provider.ClientSecretEncrypted)

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)
	data.Set("client_id", provider.ClientID)
	data.Set("client_secret", clientSecret)
	data.Set("code_verifier", codeVerifier)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, discovery.TokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token endpoint returned %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	return &tokenResp, nil
}

func (s *Service) parseIDToken(idToken string, provider *domain.OIDCProvider, expectedNonce string) (*IDTokenClaims, error) {
	// Parse without verification first to get claims
	// In production, you should verify the signature using the JWKS endpoint
	token, _, err := jwt.NewParser().ParseUnverified(idToken, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse ID token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Validate issuer
	iss, ok := claims["iss"].(string)
	if !ok || iss != provider.IssuerURL {
		return nil, fmt.Errorf("invalid issuer: expected %s, got %s", provider.IssuerURL, iss)
	}

	// Validate audience
	aud, ok := claims["aud"].(string)
	if !ok || aud != provider.ClientID {
		// Check if aud is an array
		if audArray, ok := claims["aud"].([]any); ok {
			found := false
			for _, a := range audArray {
				if aStr, ok := a.(string); ok && aStr == provider.ClientID {
					found = true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("invalid audience")
			}
		} else {
			return nil, fmt.Errorf("invalid audience: expected %s, got %v", provider.ClientID, claims["aud"])
		}
	}

	// Validate nonce
	nonce, ok := claims["nonce"].(string)
	if !ok || nonce != expectedNonce {
		return nil, fmt.Errorf("nonce mismatch")
	}

	// Validate expiration
	exp, ok := claims["exp"].(float64)
	if !ok || time.Unix(int64(exp), 0).Before(time.Now()) {
		return nil, fmt.Errorf("token expired")
	}

	// Extract user info
	subject, ok := claims["sub"].(string)
	if !ok || subject == "" {
		return nil, fmt.Errorf("missing subject claim")
	}
	result := &IDTokenClaims{
		Subject: subject,
	}

	if email, ok := claims["email"].(string); ok {
		result.Email = email
	}
	if emailVerified, ok := claims["email_verified"].(bool); ok {
		result.EmailVerified = emailVerified
	}
	if name, ok := claims["name"].(string); ok {
		result.Name = name
	}
	if preferredUsername, ok := claims["preferred_username"].(string); ok {
		result.PreferredUsername = preferredUsername
	}
	if groups, ok := claims["groups"].([]any); ok {
		for _, g := range groups {
			if gStr, ok := g.(string); ok {
				result.Groups = append(result.Groups, gStr)
			}
		}
	}

	return result, nil
}

func (s *Service) findOrCreateUser(ctx context.Context, provider *domain.OIDCProvider, claims *IDTokenClaims) (*domain.User, error) {
	// Check if user is already linked
	existingUser, err := s.userLinks.GetUserByOIDC(ctx, provider.ID, claims.Subject)
	if err == nil {
		// User already linked
		if existingUser.IsDisabled {
			return nil, domain.ErrUserDisabled
		}
		return existingUser, nil
	}

	if !errors.Is(err, domain.ErrOIDCUserLinkNotFound) {
		return nil, fmt.Errorf("failed to check OIDC link: %w", err)
	}

	// User not linked - check if auto-create is enabled
	if !provider.AutoCreateUsers {
		return nil, fmt.Errorf("automatic user creation is disabled for this provider")
	}

	// Determine username
	username := claims.PreferredUsername
	if username == "" {
		username = claims.Subject
	}

	// Check if username already exists
	exists, err := s.users.UsernameExists(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to check username: %w", err)
	}
	if exists {
		// Try with email prefix
		if claims.Email != "" {
			emailParts := strings.Split(claims.Email, "@")
			username = emailParts[0]
			exists, err = s.users.UsernameExists(ctx, username)
			if err != nil {
				return nil, fmt.Errorf("failed to check username: %w", err)
			}
			if exists {
				// Append random suffix
				suffix, err := generateRandomString(4)
				if err != nil {
					return nil, fmt.Errorf("failed to generate suffix: %w", err)
				}
				username = fmt.Sprintf("%s_%s", username, suffix)
			}
		} else {
			suffix, err := generateRandomString(4)
			if err != nil {
				return nil, fmt.Errorf("failed to generate suffix: %w", err)
			}
			username = fmt.Sprintf("%s_%s", username, suffix)
		}
	}

	// Create user
	var email *string
	if claims.Email != "" {
		email = &claims.Email
	}

	var displayName *string
	if claims.Name != "" {
		displayName = &claims.Name
	}

	user, err := s.users.Create(ctx, domain.CreateUserParams{
		Username:    username,
		Email:       email,
		DisplayName: displayName,
		IsAdmin:     provider.DefaultAdmin,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create OIDC link
	_, err = s.userLinks.Create(ctx, domain.CreateOIDCUserLinkParams{
		UserID:     user.ID,
		ProviderID: provider.ID,
		Subject:    claims.Subject,
		Email:      email,
	})
	if err != nil {
		// Rollback user creation on link failure
		if delErr := s.users.Delete(ctx, user.ID); delErr != nil {
			slog.Error("failed to rollback user creation",
				slog.String("user_id", user.ID.String()),
				slog.Any("error", delErr))
		}
		return nil, fmt.Errorf("failed to create OIDC link: %w", err)
	}

	slog.Info("created new user via OIDC",
		slog.String("user_id", user.ID.String()),
		slog.String("username", username),
		slog.String("provider", provider.Name))

	return user, nil
}

func (s *Service) createSession(ctx context.Context, user *domain.User, deviceID, deviceName, clientName, clientVersion *string) (*domain.AuthResult, error) {
	sessionID := uuid.New()
	now := time.Now()
	accessExpiry := now.Add(s.accessDuration)
	refreshExpiry := now.Add(s.refreshDuration)

	// Generate tokens
	accessToken, err := s.tokens.GenerateAccessToken(domain.TokenClaims{
		UserID:    user.ID,
		SessionID: sessionID,
		Username:  user.Username,
		IsAdmin:   user.IsAdmin,
		IssuedAt:  now,
		ExpiresAt: accessExpiry,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.tokens.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Hash tokens for storage
	accessTokenHash := s.tokens.HashToken(accessToken)
	refreshTokenHash := s.tokens.HashToken(refreshToken)

	// Create session
	_, err = s.sessions.Create(ctx, domain.CreateSessionParams{
		UserID:           user.ID,
		TokenHash:        accessTokenHash,
		RefreshTokenHash: &refreshTokenHash,
		ExpiresAt:        accessExpiry,
		RefreshExpiresAt: &refreshExpiry,
		DeviceID:         deviceID,
		DeviceName:       deviceName,
		ClientName:       clientName,
		ClientVersion:    clientVersion,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &domain.AuthResult{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExpiry,
		SessionID:    sessionID,
	}, nil
}

func (s *Service) fetchDiscovery(ctx context.Context, issuerURL string) (*DiscoveryDocument, error) {
	// Check cache
	s.discoveryCacheMu.RLock()
	cached, exists := s.discoveryCache[issuerURL]
	s.discoveryCacheMu.RUnlock()

	if exists && time.Since(cached.FetchedAt) < time.Hour {
		return cached.Document, nil
	}

	// Fetch discovery document
	discoveryURL := strings.TrimSuffix(issuerURL, "/") + "/.well-known/openid-configuration"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, discoveryURL, http.NoBody)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch discovery document: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("discovery endpoint returned %d", resp.StatusCode)
	}

	var doc DiscoveryDocument
	if err := json.NewDecoder(resp.Body).Decode(&doc); err != nil {
		return nil, fmt.Errorf("failed to parse discovery document: %w", err)
	}

	// Cache the result
	s.discoveryCacheMu.Lock()
	s.discoveryCache[issuerURL] = &cachedDiscovery{
		Document:  &doc,
		FetchedAt: time.Now(),
	}
	s.discoveryCacheMu.Unlock()

	return &doc, nil
}

func (s *Service) decryptClientSecret(encrypted []byte) string {
	// TODO: Implement proper encryption/decryption using server key
	// For now, we'll assume the secret is stored in plain text (base64 encoded)
	// In production, use AES-GCM with a server key from config
	return string(encrypted)
}

// LinkUser links an existing user to an OIDC provider.
func (s *Service) LinkUser(ctx context.Context, userID, providerID uuid.UUID, subject string, email *string) error {
	// Check if link already exists
	exists, err := s.userLinks.Exists(ctx, providerID, subject)
	if err != nil {
		return fmt.Errorf("failed to check existing link: %w", err)
	}
	if exists {
		return domain.ErrDuplicateOIDCLink
	}

	_, err = s.userLinks.Create(ctx, domain.CreateOIDCUserLinkParams{
		UserID:     userID,
		ProviderID: providerID,
		Subject:    subject,
		Email:      email,
	})
	if err != nil {
		return fmt.Errorf("failed to create OIDC link: %w", err)
	}

	return nil
}

// UnlinkUser removes an OIDC link from a user.
func (s *Service) UnlinkUser(ctx context.Context, userID, linkID uuid.UUID) error {
	// Verify the link belongs to the user
	links, err := s.userLinks.GetByUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user links: %w", err)
	}

	found := false
	for _, link := range links {
		if link.ID == linkID {
			found = true
			break
		}
	}

	if !found {
		return domain.ErrOIDCUserLinkNotFound
	}

	return s.userLinks.Delete(ctx, linkID)
}

// GetUserLinks returns all OIDC links for a user.
func (s *Service) GetUserLinks(ctx context.Context, userID uuid.UUID) ([]*domain.OIDCUserLinkWithProvider, error) {
	return s.userLinks.GetByUser(ctx, userID)
}

// CleanupExpiredStates removes expired authentication states.
// Should be called periodically.
func (s *Service) CleanupExpiredStates() {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()

	now := time.Now()
	for state, data := range s.stateStore {
		if now.Sub(data.CreatedAt) > s.stateExpiry {
			delete(s.stateStore, state)
		}
	}
}

// Helper functions

func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes)[:length], nil
}

func generateCodeChallenge(verifier string) string {
	h := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(h[:])
}
