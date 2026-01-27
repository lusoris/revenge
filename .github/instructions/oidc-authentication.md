# OIDC/SSO Integration Guide

> OpenID Connect authentication for Jellyfin Go

## Overview

Jellyfin Go supports OIDC (OpenID Connect) for Single Sign-On with identity providers like:
- Keycloak
- Authentik
- Auth0
- Okta
- Azure AD / Entra ID
- Google Workspace
- Any OIDC-compliant provider

## Dependencies

```go
import (
    "github.com/coreos/go-oidc/v3/oidc"
    "golang.org/x/oauth2"
)
```

## Provider Configuration

### Database Schema

```sql
CREATE TABLE oidc_providers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    issuer_url VARCHAR(512) NOT NULL,
    client_id VARCHAR(255) NOT NULL,
    client_secret_encrypted BYTEA NOT NULL,
    scopes TEXT[] DEFAULT ARRAY['openid', 'profile', 'email'],
    enabled BOOLEAN DEFAULT true,
    auto_create_users BOOLEAN DEFAULT true,
    default_role VARCHAR(50) DEFAULT 'user',
    claim_mappings JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE oidc_user_links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider_id UUID NOT NULL REFERENCES oidc_providers(id) ON DELETE CASCADE,
    subject VARCHAR(255) NOT NULL,  -- OIDC 'sub' claim
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(provider_id, subject)
);
```

## Authorization Code Flow with PKCE

```go
package auth

import (
    "context"
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"

    "github.com/coreos/go-oidc/v3/oidc"
    "golang.org/x/oauth2"
)

type OIDCProvider struct {
    Name     string
    Provider *oidc.Provider
    Config   oauth2.Config
    Verifier *oidc.IDTokenVerifier
}

// NewOIDCProvider creates a new OIDC provider from configuration
func NewOIDCProvider(ctx context.Context, cfg ProviderConfig) (*OIDCProvider, error) {
    // Discover OIDC configuration from issuer
    provider, err := oidc.NewProvider(ctx, cfg.IssuerURL)
    if err != nil {
        return nil, fmt.Errorf("failed to discover OIDC provider: %w", err)
    }

    // Configure OAuth2
    oauth2Config := oauth2.Config{
        ClientID:     cfg.ClientID,
        ClientSecret: cfg.ClientSecret,
        Endpoint:     provider.Endpoint(),
        RedirectURL:  cfg.RedirectURL,
        Scopes:       append([]string{oidc.ScopeOpenID}, cfg.Scopes...),
    }

    // Create ID token verifier
    verifier := provider.Verifier(&oidc.Config{
        ClientID: cfg.ClientID,
    })

    return &OIDCProvider{
        Name:     cfg.Name,
        Provider: provider,
        Config:   oauth2Config,
        Verifier: verifier,
    }, nil
}

// GeneratePKCE creates PKCE code verifier and challenge
func GeneratePKCE() (verifier, challenge string, err error) {
    // Generate random verifier (43-128 chars)
    b := make([]byte, 32)
    if _, err := rand.Read(b); err != nil {
        return "", "", err
    }
    verifier = base64.RawURLEncoding.EncodeToString(b)

    // Create S256 challenge
    h := sha256.Sum256([]byte(verifier))
    challenge = base64.RawURLEncoding.EncodeToString(h[:])

    return verifier, challenge, nil
}

// AuthCodeURL generates the authorization URL
func (p *OIDCProvider) AuthCodeURL(state, codeChallenge string) string {
    return p.Config.AuthCodeURL(state,
        oauth2.SetAuthURLParam("code_challenge", codeChallenge),
        oauth2.SetAuthURLParam("code_challenge_method", "S256"),
    )
}

// Exchange exchanges auth code for tokens
func (p *OIDCProvider) Exchange(ctx context.Context, code, codeVerifier string) (*oauth2.Token, error) {
    return p.Config.Exchange(ctx, code,
        oauth2.SetAuthURLParam("code_verifier", codeVerifier),
    )
}

// VerifyIDToken validates the ID token
func (p *OIDCProvider) VerifyIDToken(ctx context.Context, rawIDToken string) (*oidc.IDToken, error) {
    return p.Verifier.Verify(ctx, rawIDToken)
}
```

## HTTP Handlers

```go
// GET /Auth/OIDC/Providers - List available providers
func (h *AuthHandler) ListOIDCProviders(w http.ResponseWriter, r *http.Request) {
    providers, err := h.oidcService.ListEnabledProviders(r.Context())
    if err != nil {
        writeError(w, http.StatusInternalServerError, "Failed to list providers")
        return
    }

    // Return only public info (no secrets)
    response := make([]ProviderInfo, len(providers))
    for i, p := range providers {
        response[i] = ProviderInfo{
            ID:          p.ID,
            Name:        p.Name,
            DisplayName: p.DisplayName,
        }
    }
    writeJSON(w, http.StatusOK, response)
}

// GET /Auth/OIDC/Authorize/{providerId} - Start OIDC flow
func (h *AuthHandler) OIDCAuthorize(w http.ResponseWriter, r *http.Request) {
    providerID, err := uuid.Parse(r.PathValue("providerId"))
    if err != nil {
        writeError(w, http.StatusBadRequest, "Invalid provider ID")
        return
    }

    provider, err := h.oidcService.GetProvider(r.Context(), providerID)
    if err != nil {
        writeError(w, http.StatusNotFound, "Provider not found")
        return
    }

    // Generate state and PKCE
    state := generateSecureToken(32)
    verifier, challenge, _ := GeneratePKCE()

    // Store in session/cache (expires in 10 minutes)
    h.cache.Set(r.Context(), "oidc:state:"+state, OIDCState{
        ProviderID:   providerID,
        CodeVerifier: verifier,
        RedirectURL:  r.URL.Query().Get("redirect"),
    }, 10*time.Minute)

    // Redirect to IdP
    authURL := provider.AuthCodeURL(state, challenge)
    http.Redirect(w, r, authURL, http.StatusFound)
}

// GET /Auth/OIDC/Callback - Handle OIDC callback
func (h *AuthHandler) OIDCCallback(w http.ResponseWriter, r *http.Request) {
    code := r.URL.Query().Get("code")
    state := r.URL.Query().Get("state")

    if code == "" || state == "" {
        writeError(w, http.StatusBadRequest, "Missing code or state")
        return
    }

    // Retrieve and validate state
    var oidcState OIDCState
    if err := h.cache.Get(r.Context(), "oidc:state:"+state, &oidcState); err != nil {
        writeError(w, http.StatusBadRequest, "Invalid or expired state")
        return
    }
    h.cache.Delete(r.Context(), "oidc:state:"+state)

    // Get provider
    provider, err := h.oidcService.GetProvider(r.Context(), oidcState.ProviderID)
    if err != nil {
        writeError(w, http.StatusInternalServerError, "Provider error")
        return
    }

    // Exchange code for tokens
    token, err := provider.Exchange(r.Context(), code, oidcState.CodeVerifier)
    if err != nil {
        writeError(w, http.StatusUnauthorized, "Token exchange failed")
        return
    }

    // Extract and verify ID token
    rawIDToken, ok := token.Extra("id_token").(string)
    if !ok {
        writeError(w, http.StatusUnauthorized, "Missing ID token")
        return
    }

    idToken, err := provider.VerifyIDToken(r.Context(), rawIDToken)
    if err != nil {
        writeError(w, http.StatusUnauthorized, "Invalid ID token")
        return
    }

    // Extract claims
    var claims struct {
        Subject string `json:"sub"`
        Email   string `json:"email"`
        Name    string `json:"name"`
        Groups  []string `json:"groups"`
    }
    if err := idToken.Claims(&claims); err != nil {
        writeError(w, http.StatusInternalServerError, "Failed to parse claims")
        return
    }

    // Find or create user
    user, err := h.oidcService.FindOrCreateUser(r.Context(), FindOrCreateParams{
        ProviderID: oidcState.ProviderID,
        Subject:    claims.Subject,
        Email:      claims.Email,
        Name:       claims.Name,
        Groups:     claims.Groups,
    })
    if err != nil {
        writeError(w, http.StatusInternalServerError, "User provisioning failed")
        return
    }

    // Create session and return JWT
    session, err := h.sessionService.CreateSession(r.Context(), user.ID)
    if err != nil {
        writeError(w, http.StatusInternalServerError, "Session creation failed")
        return
    }

    // Redirect with token or return JSON
    if oidcState.RedirectURL != "" {
        redirectURL := fmt.Sprintf("%s?token=%s", oidcState.RedirectURL, session.AccessToken)
        http.Redirect(w, r, redirectURL, http.StatusFound)
        return
    }

    writeJSON(w, http.StatusOK, AuthResponse{
        AccessToken:  session.AccessToken,
        RefreshToken: session.RefreshToken,
        ExpiresIn:    int(session.ExpiresIn.Seconds()),
        User:         userToDTO(user),
    })
}
```

## Claim Mapping

```go
// MapClaims maps OIDC claims to user attributes
func MapClaims(claims map[string]any, mappings ClaimMappings) UserAttributes {
    attrs := UserAttributes{}

    // Email mapping (default: "email")
    if emailKey := mappings.Email; emailKey != "" {
        if email, ok := claims[emailKey].(string); ok {
            attrs.Email = email
        }
    }

    // Name mapping (default: "name" or "preferred_username")
    if nameKey := mappings.Name; nameKey != "" {
        if name, ok := claims[nameKey].(string); ok {
            attrs.Name = name
        }
    }

    // Groups/Roles mapping
    if groupsKey := mappings.Groups; groupsKey != "" {
        if groups, ok := claims[groupsKey].([]any); ok {
            for _, g := range groups {
                if group, ok := g.(string); ok {
                    attrs.Groups = append(attrs.Groups, group)
                }
            }
        }
    }

    // Admin check (configurable group name)
    if mappings.AdminGroup != "" {
        for _, g := range attrs.Groups {
            if g == mappings.AdminGroup {
                attrs.IsAdmin = true
                break
            }
        }
    }

    return attrs
}
```

## Configuration Example

### Keycloak

```yaml
# config.yaml
oidc:
  providers:
    - name: keycloak
      display_name: "Login with Keycloak"
      issuer_url: "https://keycloak.example.com/realms/jellyfin"
      client_id: "jellyfin-go"
      client_secret: "${OIDC_KEYCLOAK_SECRET}"
      scopes: ["openid", "profile", "email", "groups"]
      auto_create_users: true
      default_role: "user"
      claim_mappings:
        email: "email"
        name: "preferred_username"
        groups: "groups"
        admin_group: "jellyfin-admins"
```

### Authentik

```yaml
oidc:
  providers:
    - name: authentik
      display_name: "Login with Authentik"
      issuer_url: "https://auth.example.com/application/o/jellyfin/"
      client_id: "jellyfin-go"
      client_secret: "${OIDC_AUTHENTIK_SECRET}"
      scopes: ["openid", "profile", "email", "groups"]
```

### Auth0

```yaml
oidc:
  providers:
    - name: auth0
      display_name: "Login with Auth0"
      issuer_url: "https://your-tenant.auth0.com/"
      client_id: "your-client-id"
      client_secret: "${OIDC_AUTH0_SECRET}"
      scopes: ["openid", "profile", "email"]
      claim_mappings:
        groups: "https://your-namespace/roles"
```

## Security Considerations

1. **Always use PKCE** - Prevents authorization code interception
2. **Validate state parameter** - Prevents CSRF attacks
3. **Short state expiration** - 10 minutes max
4. **Encrypt client secrets** - Never store in plaintext
5. **Validate issuer** - Ensure tokens come from expected IdP
6. **Check audience** - Verify `aud` claim matches client_id
7. **Use HTTPS** - Never transmit tokens over HTTP

## Testing

```go
func TestOIDCCallback(t *testing.T) {
    // Use mock OIDC server for testing
    mockServer := oidctest.NewServer(t)
    defer mockServer.Close()

    // Configure provider with mock issuer
    cfg := ProviderConfig{
        IssuerURL:    mockServer.IssuerURL(),
        ClientID:     "test-client",
        ClientSecret: "test-secret",
    }

    // ... test authorization flow
}
```
