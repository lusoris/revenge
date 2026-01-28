# Authentication Providers

> OIDC/SSO integration for user authentication

---

## Overview

Revenge supports external authentication via OpenID Connect (OIDC):
- Single Sign-On (SSO)
- Centralized user management
- Multi-factor authentication (via provider)
- Session management

---

## Providers

| Provider | Type | Status |
|----------|------|--------|
| [Authelia](AUTHELIA.md) | Self-hosted SSO | 🟢 Recommended |
| [Authentik](AUTHENTIK.md) | Self-hosted IdP | 🟢 Supported |
| [Keycloak](KEYCLOAK.md) | Enterprise IdP | 🟢 Supported |
| [Generic OIDC](GENERIC_OIDC.md) | Any OIDC provider | 🟢 Supported |

---

## Provider Details

### Authelia
**Lightweight self-hosted SSO**

- ✅ Simple setup
- ✅ 2FA support
- ✅ Reverse proxy integration
- ✅ Low resource usage
- 🎯 **Recommended for homelab**

### Authentik
**Modern identity provider**

- ✅ Full IdP features
- ✅ Beautiful UI
- ✅ LDAP/SCIM support
- ✅ Application management
- ⚠️ Higher resource usage

### Keycloak
**Enterprise-grade IdP**

- ✅ Full enterprise features
- ✅ Federation support
- ✅ Fine-grained permissions
- ⚠️ Complex setup
- ⚠️ High resource usage

### Generic OIDC
**Any OIDC-compliant provider**

- ✅ Google, GitHub, etc.
- ✅ Azure AD
- ✅ Any OIDC provider
- ⚠️ Manual configuration

---

## Authentication Flow

```
User → Revenge Login
    ↓
Redirect to OIDC Provider
    ↓
User authenticates (+ 2FA if enabled)
    ↓
Provider redirects back with code
    ↓
Revenge exchanges code for tokens
    ↓
Validate ID token, create session
    ↓
User logged in
```

---

## Configuration

```yaml
auth:
  # Built-in auth (default)
  local:
    enabled: true

  # OIDC providers
  oidc:
    enabled: true

    # Default provider
    default_provider: authelia

    providers:
      authelia:
        enabled: true
        issuer: "https://auth.example.com"
        client_id: "${OIDC_CLIENT_ID}"
        client_secret: "${OIDC_CLIENT_SECRET}"
        scopes: ["openid", "profile", "email"]

      authentik:
        enabled: false
        issuer: "https://authentik.example.com/application/o/revenge/"
        client_id: "${AUTHENTIK_CLIENT_ID}"
        client_secret: "${AUTHENTIK_CLIENT_SECRET}"
```

---

## User Mapping

Map OIDC claims to Revenge user attributes:

```yaml
auth:
  oidc:
    claim_mapping:
      username: "preferred_username"
      email: "email"
      name: "name"
      groups: "groups"

    # Auto-create users from OIDC
    auto_provision: true

    # Default role for new users
    default_role: "user"
```

---

## Related Documentation

- [OIDC Implementation](../../features/OIDC_IMPLEMENTATION.md)
- [User Management](../../features/USER_MANAGEMENT.md)
