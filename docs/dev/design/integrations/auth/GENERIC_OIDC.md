# Generic OIDC Integration

> Support for any OpenID Connect compatible provider

**Status**: 🟡 PLANNED
**Priority**: 🟡 MEDIUM (Phase 1 - Core Infrastructure)
**Type**: OIDC Identity Provider

---

## Overview

Generic OIDC support allows Revenge to integrate with any OpenID Connect compliant identity provider not specifically documented. This includes:
- Cloud providers (Auth0, Okta, Azure AD, Google)
- Self-hosted solutions (Dex, Ory Hydra, Zitadel)
- Custom OIDC implementations

**Integration Points**:
- **OIDC Discovery**: Auto-configure from `.well-known/openid-configuration`
- **Standard flows**: Authorization Code, PKCE
- **Configurable claims**: Map any claims to Revenge user fields
- **Flexible role mapping**: Support various group/role claim formats

---

## Developer Resources

- 📚 **OIDC Spec**: https://openid.net/specs/openid-connect-core-1_0.html
- 🔗 **Discovery Spec**: https://openid.net/specs/openid-connect-discovery-1_0.html
- 🔗 **JWT.io**: https://jwt.io/ (token debugging)
- 🔗 **OIDC Playground**: https://openidconnect.net/

---

## OIDC Requirements

Any provider must support:

| Requirement | Description |
|-------------|-------------|
| Discovery | `.well-known/openid-configuration` endpoint |
| Authorization Code Flow | Standard OAuth2 auth code grant |
| PKCE (optional) | Proof Key for Code Exchange (S256) |
| Token Endpoint | Exchange code for tokens |
| UserInfo Endpoint | Fetch user claims |
| JWKS | JSON Web Key Set for token validation |

### Required Claims

| Claim | Required | Description |
|-------|----------|-------------|
| `sub` | ✅ Yes | Unique user identifier |
| `email` | 🟡 Recommended | User email address |
| `preferred_username` | 🟡 Recommended | Display name |
| `name` | 🟢 Optional | Full name |
| `groups` | 🟢 Optional | Group membership |

---

## Configuration

### Basic Configuration

```yaml
# configs/config.yaml
integrations:
  oidc:
    enabled: true
    provider: "generic"
    generic:
      # Required: OIDC Discovery URL
      issuer_url: "https://idp.example.com"

      # Required: Client credentials
      client_id: "${REVENGE_OIDC_CLIENT_ID}"
      client_secret: "${REVENGE_OIDC_CLIENT_SECRET}"

      # Required: Callback URL
      redirect_uri: "https://revenge.example.com/api/v1/auth/oidc/callback"

      # Scopes to request
      scopes:
        - openid
        - profile
        - email
        # - groups  # If supported by provider

      # Auto-create users on first login
      auto_provision: true
```

### Advanced Configuration

```yaml
integrations:
  oidc:
    enabled: true
    provider: "generic"
    generic:
      issuer_url: "https://idp.example.com"
      client_id: "${REVENGE_OIDC_CLIENT_ID}"
      client_secret: "${REVENGE_OIDC_CLIENT_SECRET}"
      redirect_uri: "https://revenge.example.com/api/v1/auth/oidc/callback"

      scopes:
        - openid
        - profile
        - email
        - groups

      # Override discovery endpoints (if non-standard)
      endpoints:
        authorization: ""  # Empty = use discovery
        token: ""
        userinfo: ""
        jwks: ""
        logout: ""  # Optional: for single logout

      # Claim mapping
      claims:
        subject: "sub"                    # User identifier
        username: "preferred_username"    # Falls back to email
        email: "email"
        name: "name"
        groups: "groups"                  # Or "roles", "realm_access.roles", etc.

      # Role/group mapping
      role_mappings:
        admin: admin              # group "admin" → role "admin"
        users: user               # group "users" → role "user"
        limited: restricted       # group "limited" → role "restricted"

      # Default role if no mapping matches
      default_role: "user"

      # Security settings
      pkce_enabled: true          # Use PKCE (recommended)
      pkce_method: "S256"         # S256 or plain
      token_validation:
        verify_signature: true
        verify_issuer: true
        verify_audience: true
        clock_skew_seconds: 60    # Allow for clock drift

      # User provisioning
      auto_provision: true
      allow_registration: true
      update_on_login: true       # Update user info on each login
```

---

## Implementation Checklist

- [ ] **OIDC Client** (`internal/service/oidc/provider_generic.go`)
  - [ ] Discovery document fetching and caching
  - [ ] Endpoint extraction from discovery
  - [ ] Authorization URL generation
  - [ ] Code exchange
  - [ ] Token validation (signature, issuer, audience, expiry)
  - [ ] UserInfo fetching
  - [ ] Token refresh

- [ ] **Claim Mapping** (`internal/service/oidc/claims.go`)
  - [ ] Configurable claim paths (dot notation)
  - [ ] Fallback claims (e.g., username → email)
  - [ ] Array claim handling (groups, roles)
  - [ ] Nested claim extraction

- [ ] **Role Mapping** (`internal/service/oidc/roles.go`)
  - [ ] Direct mapping (group name → role)
  - [ ] Prefix/suffix stripping
  - [ ] Default role fallback
  - [ ] Multiple role support

---

## Claim Path Syntax

Support dot notation for nested claims:

```yaml
claims:
  # Simple claims
  subject: "sub"
  email: "email"

  # Nested claims
  groups: "realm_access.roles"        # Keycloak-style
  groups: "resource_access.app.roles" # Client roles
  groups: "cognito:groups"            # AWS Cognito

  # Array indexing
  primary_email: "emails[0].value"    # First email in array
```

---

## Provider-Specific Examples

### Auth0

```yaml
generic:
  issuer_url: "https://your-tenant.auth0.com"
  scopes: [openid, profile, email]
  claims:
    groups: "https://your-app.com/roles"  # Custom claim
```

### Azure AD

```yaml
generic:
  issuer_url: "https://login.microsoftonline.com/{tenant}/v2.0"
  scopes: [openid, profile, email]
  claims:
    groups: "groups"  # Requires group claims in app registration
    username: "preferred_username"
```

### Google

```yaml
generic:
  issuer_url: "https://accounts.google.com"
  scopes: [openid, profile, email]
  claims:
    username: "email"  # Google uses email as identifier
  # Note: Google doesn't support groups
```

### Okta

```yaml
generic:
  issuer_url: "https://your-org.okta.com"
  scopes: [openid, profile, email, groups]
  claims:
    groups: "groups"
```

### AWS Cognito

```yaml
generic:
  issuer_url: "https://cognito-idp.{region}.amazonaws.com/{pool-id}"
  scopes: [openid, profile, email]
  claims:
    groups: "cognito:groups"
    username: "cognito:username"
```

### Dex

```yaml
generic:
  issuer_url: "https://dex.example.com"
  scopes: [openid, profile, email, groups]
  claims:
    groups: "groups"
```

### Zitadel

```yaml
generic:
  issuer_url: "https://your-instance.zitadel.cloud"
  scopes: [openid, profile, email, "urn:zitadel:iam:org:project:roles"]
  claims:
    groups: "urn:zitadel:iam:org:project:roles"
```

---

## Database Schema

Uses shared OIDC tables from [Authelia Integration](AUTHELIA.md#database-schema).

Additional table for provider configuration:

```sql
-- Allow multiple OIDC providers
CREATE TABLE oidc_providers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,           -- Display name
    slug VARCHAR(50) NOT NULL UNIQUE,            -- URL-safe identifier
    provider_type VARCHAR(50) NOT NULL,          -- generic, authelia, authentik, keycloak
    issuer_url TEXT NOT NULL,
    client_id VARCHAR(255) NOT NULL,
    client_secret_encrypted BYTEA NOT NULL,
    scopes TEXT[] NOT NULL DEFAULT ARRAY['openid', 'profile', 'email'],
    config JSONB NOT NULL DEFAULT '{}',          -- Provider-specific config
    enabled BOOLEAN NOT NULL DEFAULT true,
    auto_provision BOOLEAN NOT NULL DEFAULT true,
    icon_url TEXT,                               -- For login button
    display_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

## Error Handling

| Error | Cause | Solution |
|-------|-------|----------|
| Discovery failed | Network error or invalid URL | Check issuer_url, verify connectivity |
| Invalid token | Signature validation failed | Verify JWKS URL, check key rotation |
| Missing claim | Provider doesn't include claim | Check scopes, add claim mapping |
| Role mapping failed | No matching groups | Configure default_role |
| PKCE error | Provider doesn't support PKCE | Set pkce_enabled: false |

---

## Testing Configuration

Validate OIDC configuration before deployment:

```bash
# Fetch discovery document
curl https://idp.example.com/.well-known/openid-configuration | jq

# Verify JWKS
curl https://idp.example.com/.well-known/jwks.json | jq

# Decode JWT token (after login)
# Use https://jwt.io or:
echo $TOKEN | cut -d. -f2 | base64 -d | jq
```

---

## Security Best Practices

1. **Always use PKCE** when supported
2. **Validate all tokens** (signature, issuer, audience, expiry)
3. **Use state parameter** to prevent CSRF
4. **Store secrets securely** (encrypted in database)
5. **Rotate client secrets** periodically
6. **Use HTTPS only** for all OIDC endpoints
7. **Limit scopes** to what's needed

---

## Related Documentation

- [Authelia Integration](AUTHELIA.md) - Recommended self-hosted
- [Authentik Integration](AUTHENTIK.md) - Enterprise self-hosted
- [Keycloak Integration](KEYCLOAK.md) - Enterprise with LDAP
