---
sources:
  - name: Authelia Documentation
    url: ../../../sources/security/authelia.md
    note: Auto-resolved from authelia
  - name: Authelia OIDC Guide
    url: https://www.authelia.com/integration/openid-connect/introduction/
    note: OIDC implementation details
  - name: Authelia Configuration
    url: https://www.authelia.com/configuration/prologue/introduction/
    note: Configuration reference
design_refs:
  - title: 01_ARCHITECTURE
    path: ../../architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: ../../architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: ../../architecture/03_METADATA_SYSTEM.md
---

## Table of Contents

- [Authelia](#authelia)
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
- [Authelia OIDC configuration](#authelia-oidc-configuration)
    - [Config Keys](#config-keys)
  - [API Endpoints](#api-endpoints)
  - [Testing Strategy](#testing-strategy)
    - [Unit Tests](#unit-tests)
    - [Integration Tests](#integration-tests)
    - [Test Coverage](#test-coverage)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)


# Authelia


**Created**: 2026-02-01
**Status**: ✅ Complete
**Category**: integration


> Integration with Authelia

> Lightweight authentication and authorization server for homelab environments
**API Base URL**: `https://auth.homelab.local`
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
internal/integration/authelia/
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
// AutheliaProvider implements OIDC for Authelia
type AutheliaProvider struct {
    config   *AutheliaConfig
    oauth    *oauth2.Config
    verifier *oidc.IDTokenVerifier
}

type AutheliaConfig struct {
    Issuer           string            // https://auth.homelab.local
    ClientID         string            // revenge-client
    ClientSecret     string            // secret
    RedirectURL      string            // https://revenge.local/api/v1/auth/oidc/callback
    Scopes           []string          // openid, profile, email, groups
    GroupMappings    map[string]string // Authelia group → Revenge role
    AutoCreateUsers  bool              // Create users on first login
    UpdateUserInfo   bool              // Update email/name on each login
    UsePKCE          bool              // Enable PKCE (recommended)
    RequireMFA       bool              // Require two-factor authentication
}

// OIDCProvider interface (generic)
type OIDCProvider interface {
    // Get authorization URL
    GetAuthURL(state string, pkceChallenge string) string

    // Exchange code for tokens
    ExchangeCode(ctx context.Context, code string, pkceVerifier string) (*TokenResponse, error)

    // Verify ID token
    VerifyIDToken(ctx context.Context, rawIDToken string) (*IDToken, error)

    // Get user info from UserInfo endpoint
    GetUserInfo(ctx context.Context, accessToken string) (*UserInfo, error)

    // Map provider groups to Revenge roles
    MapRoles(ctx context.Context, user *UserInfo) ([]string, error)
}

type TokenResponse struct {
    IDToken      string
    AccessToken  string
    RefreshToken string  // Optional in Authelia
    ExpiresIn    int
}

type IDToken struct {
    Subject           string   // User ID in Authelia
    Email             string
    EmailVerified     bool
    Name              string
    PreferredUsername string
    Groups            []string // Authelia groups
    IssuedAt          time.Time
    ExpiresAt         time.Time
    AuthenticationMethodsReferences []string // amr claim (MFA info)
}

type UserInfo struct {
    Sub               string   // Subject (user ID)
    Email             string
    EmailVerified     bool
    Name              string
    PreferredUsername string
    Groups            []string // From Authelia LDAP/file backend
    Picture           string   // Avatar URL (if configured)
    AMR               []string // Authentication methods (pwd, otp, etc.)
}
```


### Dependencies
**Go Packages**:
- `github.com/coreos/go-oidc/v3/oidc` - OIDC client
- `golang.org/x/oauth2` - OAuth2 flow
- `github.com/golang-jwt/jwt/v5` - JWT parsing
- `github.com/google/uuid` - PKCE code verifier generation

**External Services**:
- Authelia server (https://www.authelia.com/)
- Redis (for Authelia session storage, optional)
- LDAP server (optional backend)






## Configuration
### Environment Variables

```bash
# Authelia OIDC configuration
REVENGE_OIDC_PROVIDER=authelia
REVENGE_OIDC_AUTHELIA_ISSUER=https://auth.homelab.local
REVENGE_OIDC_AUTHELIA_CLIENT_ID=revenge-client
REVENGE_OIDC_AUTHELIA_CLIENT_SECRET=very-secret-key
REVENGE_OIDC_AUTHELIA_REDIRECT_URL=https://revenge.local/api/v1/auth/oidc/callback
```


### Config Keys
```yaml
auth:
  oidc:
    enabled: true
    provider: authelia     # authentik, authelia, keycloak, generic
    providers:
      authelia:
        issuer: ${REVENGE_OIDC_AUTHELIA_ISSUER}
        client_id: ${REVENGE_OIDC_AUTHELIA_CLIENT_ID}
        client_secret: ${REVENGE_OIDC_AUTHELIA_CLIENT_SECRET}
        redirect_url: https://revenge.local/api/v1/auth/oidc/callback
        scopes:
          - openid
          - profile
          - email
          - groups          # Authelia groups claim
        group_mappings:
          # Authelia group → Revenge role
          "admins": "admin"
          "users": "user"
          "readonly": "readonly"
          "family": "family"
        auto_create_users: true
        update_user_info: true
        user_claim: "preferred_username"  # Field to use as username
        use_pkce: true                    # Enable PKCE (recommended)
        require_mfa: false                # Require two-factor auth
```



## API Endpoints
**OIDC Endpoints** (Revenge):
```
GET  /api/v1/auth/oidc/login
GET  /api/v1/auth/oidc/callback
POST /api/v1/auth/oidc/refresh
POST /api/v1/auth/oidc/logout
```

**Example - Initiate Login**:
```
GET /api/v1/auth/oidc/login?provider=authelia

→ Generates PKCE verifier and challenge
→ Redirects to:
https://auth.homelab.local/api/oidc/authorization?
  client_id=revenge-client&
  redirect_uri=https://revenge.local/api/v1/auth/oidc/callback&
  response_type=code&
  scope=openid+profile+email+groups&
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
1. Verifies state matches stored value
2. Exchanges code for tokens (includes PKCE verifier)
3. Verifies ID token signature
4. Extracts user info and groups
5. Creates or updates user
6. Maps groups to roles
7. Creates session
8. Sets session cookie
9. Redirects to /
```

**Example - Refresh Token**:
```
POST /api/v1/auth/oidc/refresh
Content-Type: application/json

{
  "refresh_token": "refresh-token-value"
}

→ Response:
{
  "access_token": "new-access-token",
  "id_token": "new-id-token",
  "expires_in": 3600
}
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
- [Authelia OIDC Guide](https://www.authelia.com/integration/openid-connect/introduction/) - OIDC implementation details
- [Authelia Configuration](https://www.authelia.com/configuration/prologue/introduction/) - Configuration reference

