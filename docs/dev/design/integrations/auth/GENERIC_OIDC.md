

---
sources:
  - name: Authelia Documentation
    url: ../../../sources/security/authelia.md
    note: Auto-resolved from authelia
  - name: Authentik Documentation
    url: ../../../sources/security/authentik.md
    note: Auto-resolved from authentik
  - name: Keycloak Documentation
    url: ../../../sources/security/keycloak.md
    note: Auto-resolved from keycloak
  - name: OpenID Connect Core
    url: ../../../sources/security/oidc-core.md
    note: Auto-resolved from oidc
design_refs:
  - title: 01_ARCHITECTURE
    path: ../../architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: ../../architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: ../../architecture/03_METADATA_SYSTEM.md
---

## Table of Contents

- [Generic OIDC](#generic-oidc)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Integration Structure](#integration-structure)
    - [Data Flow](#data-flow)
    - [Provides](#provides)
  - [Implementation](#implementation)
    - [File Structure](#file-structure)
    - [Key Interfaces](#key-interfaces)
    - [Dependencies](#dependencies)
  - [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
- [Generic OIDC configuration](#generic-oidc-configuration)
    - [Config Keys](#config-keys)
  - [API Endpoints](#api-endpoints)
  - [Testing Strategy](#testing-strategy)
    - [Unit Tests](#unit-tests)
    - [Integration Tests](#integration-tests)
    - [Test Coverage](#test-coverage)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)


# Generic OIDC


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: integration


> Integration with Generic OIDC

> Support for any OpenID Connect compatible provider
**API Base URL**: `https://revenge.example.com/api/v1/auth/oidc/callback`
**Authentication**: oidc

---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | - |
| Sources | ✅ | - |
| Instructions | ✅ | - |
| Code | 🔴 | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: ✅ Complete



---


## Architecture

### Integration Structure

```
internal/integration/generic_oidc/
├── client.go              # API client
├── types.go               # Response types
├── mapper.go              # Map external → internal types
├── cache.go               # Response caching
└── client_test.go         # Tests
```

### Data Flow

<!-- Data flow diagram -->

### Provides
<!-- Data provided by integration -->


## Implementation

### File Structure

<!-- File structure -->

### Key Interfaces

```go
// GenericOIDCProvider implements OIDC for any provider
type GenericOIDCProvider struct {
    config    *GenericOIDCConfig
    discovery *DiscoveryDocument
    oauth     *oauth2.Config
    verifier  *oidc.IDTokenVerifier
}

type GenericOIDCConfig struct {
    Name              string            // Provider name (for display)
    DiscoveryURL      string            // OIDC discovery URL
    Issuer            string            // Optional: override issuer from discovery
    ClientID          string            // OAuth2 client ID
    ClientSecret      string            // OAuth2 client secret
    RedirectURL       string            // Callback URL
    Scopes            []string          // Requested scopes (openid required)
    ClaimMappings     ClaimMappings     // Map claims to user fields
    RoleClaim         string            // Claim containing roles/groups
    RoleMappings      map[string]string // Provider role → Revenge role
    AutoCreateUsers   bool              // Create users on first login
    UpdateUserInfo    bool              // Update user info on each login
    UsernameClaim     string            // Claim to use as username
    EmailClaim        string            // Claim to use as email (default: "email")
    NameClaim         string            // Claim to use as name (default: "name")
    PictureClaim      string            // Claim to use as avatar
    ExtraAuthParams   map[string]string // Additional auth params
}

type ClaimMappings struct {
    Username string // Claim path for username (e.g., "preferred_username")
    Email    string // Claim path for email
    Name     string // Claim path for display name
    Picture  string // Claim path for avatar URL
    Roles    string // Claim path for roles (e.g., "groups", "roles")
}

type DiscoveryDocument struct {
    Issuer                            string   `json:"issuer"`
    AuthorizationEndpoint             string   `json:"authorization_endpoint"`
    TokenEndpoint                     string   `json:"token_endpoint"`
    UserInfoEndpoint                  string   `json:"userinfo_endpoint"`
    JWKSUri                           string   `json:"jwks_uri"`
    IntrospectionEndpoint             string   `json:"introspection_endpoint,omitempty"`
    RevocationEndpoint                string   `json:"revocation_endpoint,omitempty"`
    EndSessionEndpoint                string   `json:"end_session_endpoint,omitempty"`
    ScopesSupported                   []string `json:"scopes_supported"`
    ResponseTypesSupported            []string `json:"response_types_supported"`
    GrantTypesSupported               []string `json:"grant_types_supported"`
    SubjectTypesSupported             []string `json:"subject_types_supported"`
    IDTokenSigningAlgValuesSupported  []string `json:"id_token_signing_alg_values_supported"`
    CodeChallengeMethodsSupported     []string `json:"code_challenge_methods_supported,omitempty"`
}

// OIDCProvider interface (generic)
type OIDCProvider interface {
    // Discover provider metadata
    Discover(ctx context.Context) (*DiscoveryDocument, error)

    // Get authorization URL
    GetAuthURL(state string) string

    // Exchange code for tokens
    ExchangeCode(ctx context.Context, code string) (*TokenResponse, error)

    // Verify ID token
    VerifyIDToken(ctx context.Context, rawIDToken string) (*IDToken, error)

    // Get user info from UserInfo endpoint
    GetUserInfo(ctx context.Context, accessToken string) (*UserInfo, error)

    // Map provider claims/roles to Revenge roles
    MapRoles(ctx context.Context, claims map[string]any) ([]string, error)
}

type TokenResponse struct {
    IDToken      string
    AccessToken  string
    RefreshToken string
    ExpiresIn    int
}

type IDToken struct {
    Subject       string         // User ID
    Email         string
    EmailVerified bool
    Name          string
    PreferredUsername string
    Picture       string
    IssuedAt      time.Time
    ExpiresAt     time.Time
    Claims        map[string]any // All claims
}

type UserInfo struct {
    Sub       string         // Subject (user ID)
    Email     string
    Name      string
    Username  string
    Picture   string
    Claims    map[string]any // All userinfo claims
}
```


### Dependencies

**Go Packages**:
- `github.com/coreos/go-oidc/v3/oidc` - OIDC client
- `golang.org/x/oauth2` - OAuth2 flow
- `github.com/golang-jwt/jwt/v5` - JWT parsing (fallback)

**External Services**:
- Any OIDC-compliant identity provider






## Configuration
### Environment Variables

```bash
# Generic OIDC configuration
REVENGE_OIDC_PROVIDER=generic
REVENGE_OIDC_GENERIC_NAME="Okta"
REVENGE_OIDC_GENERIC_DISCOVERY_URL=https://dev-123456.okta.com/.well-known/openid-configuration
REVENGE_OIDC_GENERIC_CLIENT_ID=0oa1b2c3d4e5f6g7h8i9
REVENGE_OIDC_GENERIC_CLIENT_SECRET=very-secret-key
REVENGE_OIDC_GENERIC_REDIRECT_URL=https://revenge.local/api/v1/auth/oidc/callback
```


### Config Keys

```yaml
auth:
  oidc:
    enabled: true
    provider: generic     # Use generic provider
    providers:
      generic:
        name: "Okta"      # Display name
        discovery_url: ${REVENGE_OIDC_GENERIC_DISCOVERY_URL}
        # Optional: override discovery
        # issuer: https://dev-123456.okta.com
        client_id: ${REVENGE_OIDC_GENERIC_CLIENT_ID}
        client_secret: ${REVENGE_OIDC_GENERIC_CLIENT_SECRET}
        redirect_url: https://revenge.local/api/v1/auth/oidc/callback
        scopes:
          - openid
          - profile
          - email
          - groups        # Provider-specific (optional)
        claim_mappings:
          username: "preferred_username"  # Claim path
          email: "email"
          name: "name"
          picture: "picture"
          roles: "groups"                 # Where to find roles
        role_mappings:
          # Provider role/group → Revenge role
          "Revenge Admins": "admin"
          "Revenge Users": "user"
        auto_create_users: true
        update_user_info: true
        username_claim: "preferred_username"
        email_claim: "email"
        name_claim: "name"
        picture_claim: "picture"
        extra_auth_params:
          # Provider-specific params (e.g., Azure AD resource)
          # resource: "https://graph.microsoft.com"
```



## API Endpoints
**OIDC Endpoints** (Revenge):
```
GET  /api/v1/auth/oidc/login?provider=generic
GET  /api/v1/auth/oidc/callback
POST /api/v1/auth/oidc/refresh
POST /api/v1/auth/oidc/logout
```

**Example - Initiate Login**:
```
GET /api/v1/auth/oidc/login?provider=generic

→ Redirects to:
https://provider.example.com/authorize?
  client_id=abc123&
  redirect_uri=https://revenge.local/api/v1/auth/oidc/callback&
  response_type=code&
  scope=openid+profile+email&
  state=random-state&
  code_challenge=...&
  code_challenge_method=S256
```

**Example - Callback**:
```
GET /api/v1/auth/oidc/callback?
  code=authorization-code&
  state=random-state

→ Revenge backend:
1. Verifies state
2. Exchanges code for tokens
3. Verifies ID token
4. Fetches user info
5. Creates user/session
6. Sets session cookie
7. Redirects to /
```



## Testing Strategy

### Unit Tests

<!-- Unit test strategy -->

### Integration Tests

<!-- Integration test strategy -->

### Test Coverage

Target: **80% minimum**







## Related Documentation
### Design Documents
- [01_ARCHITECTURE](../../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../../architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../../architecture/03_METADATA_SYSTEM.md)

### External Sources
- [Authelia Documentation](../../../sources/security/authelia.md) - Auto-resolved from authelia
- [Authentik Documentation](../../../sources/security/authentik.md) - Auto-resolved from authentik
- [Keycloak Documentation](../../../sources/security/keycloak.md) - Auto-resolved from keycloak
- [OpenID Connect Core](../../../sources/security/oidc-core.md) - Auto-resolved from oidc

