# Authelia Integration

<!-- SOURCES: authelia, authentik, keycloak -->

<!-- DESIGN: integrations/auth, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES, 03_METADATA_SYSTEM -->


> Self-hosted authentication and authorization server


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Overview](#overview)
- [Developer Resources](#developer-resources)
- [OIDC Details](#oidc-details)
  - [Supported Scopes](#supported-scopes)
  - [Authelia Configuration](#authelia-configuration)
- [Data Mapping](#data-mapping)
  - [Authelia Claims → Revenge User](#authelia-claims-revenge-user)
  - [Group → Role Mapping](#group-role-mapping)
- [Implementation Checklist](#implementation-checklist)
- [Configuration](#configuration)
- [Database Schema](#database-schema)
- [Authentication Flow](#authentication-flow)
- [Security Considerations](#security-considerations)
- [Troubleshooting](#troubleshooting)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related Documentation](#related-documentation)

<!-- TOC-END -->

## Status

| Dimension | Status |
|-----------|--------|
| Design | ✅ |
| Sources | ✅ |
| Instructions | ✅ |
| Code | 🔴 |
| Linting | 🔴 |
| Unit Testing | 🔴 |
| Integration Testing | 🔴 |**Priority**: 🟡 MEDIUM (Phase 1 - Core Infrastructure)
**Type**: OIDC Identity Provider

---

## Overview

Authelia is a popular open-source authentication server that provides Single Sign-On (SSO) and two-factor authentication for self-hosted applications. Revenge integrates with Authelia as an OIDC provider for:
- User authentication via OIDC
- Single Sign-On across media stack
- Two-factor authentication (TOTP, WebAuthn, Duo)
- Access control and authorization

**Integration Points**:
- **OIDC Provider**: Standard OpenID Connect flow
- **Forward Auth**: Traefik/Nginx forward authentication
- **User provisioning**: Auto-create Revenge users from OIDC claims
- **Groups/Roles**: Map Authelia groups to Revenge roles

---

## Developer Resources

- 📚 **Docs**: https://www.authelia.com/docs/
- 🔗 **OIDC Docs**: https://www.authelia.com/integration/openid-connect/introduction/
- 🔗 **GitHub**: https://github.com/authelia/authelia
- 🔗 **Configuration Reference**: https://www.authelia.com/configuration/

---

## OIDC Details

**Discovery URL**: `https://auth.example.com/.well-known/openid-configuration`
**Authorization Endpoint**: `https://auth.example.com/api/oidc/authorization`
**Token Endpoint**: `https://auth.example.com/api/oidc/token`
**UserInfo Endpoint**: `https://auth.example.com/api/oidc/userinfo`
**JWKS URI**: `https://auth.example.com/api/oidc/jwks`

### Supported Scopes

| Scope | Claims |
|-------|--------|
| `openid` | `sub` |
| `profile` | `name`, `preferred_username`, `given_name`, `family_name` |
| `email` | `email`, `email_verified` |
| `groups` | `groups` (array) |

### Authelia Configuration

```yaml
# authelia/configuration.yml
identity_providers:
  oidc:
    enabled: true
    cors:
      endpoints:
        - authorization
        - token
        - userinfo
    clients:
      - client_id: revenge
        client_name: Revenge Media Server
        client_secret: '$pbkdf2-sha512$...'  # hashed secret
        public: false
        authorization_policy: two_factor
        redirect_uris:
          - https://revenge.example.com/api/v1/auth/oidc/callback
        scopes:
          - openid
          - profile
          - email
          - groups
        userinfo_signed_response_alg: none
        token_endpoint_auth_method: client_secret_basic
```

---

## Data Mapping

### Authelia Claims → Revenge User

| Authelia Claim | Revenge Field | Notes |
|----------------|---------------|-------|
| `sub` | `oidc_subject` | Unique identifier |
| `preferred_username` | `username` | Display name |
| `email` | `email` | User email |
| `email_verified` | `email_verified` | Verification status |
| `name` | `display_name` | Full name |
| `groups` | `roles[]` | Role mapping |

### Group → Role Mapping

Configure in Revenge:

```yaml
integrations:
  oidc:
    authelia:
      group_mappings:
        - group: "admins"
          role: "admin"
        - group: "media-users"
          role: "user"
        - group: "kids"
          role: "restricted"
```

---

## Implementation Checklist

- [ ] **OIDC Client** (`internal/service/oidc/provider_authelia.go`)
  - [ ] Discovery document parsing
  - [ ] Authorization code flow
  - [ ] Token validation (JWT)
  - [ ] UserInfo fetching
  - [ ] Token refresh handling

- [ ] **User Provisioning** (`internal/service/user/oidc_provisioning.go`)
  - [ ] Auto-create user on first login
  - [ ] Map OIDC claims to user fields
  - [ ] Map groups to roles
  - [ ] Handle user updates (name, email changes)

- [ ] **Session Management**
  - [ ] Create Revenge session from OIDC token
  - [ ] Session expiry based on token expiry
  - [ ] Single logout (if supported)

---

## Configuration

```yaml
# configs/config.yaml
integrations:
  oidc:
    enabled: true
    provider: "authelia"  # or authentik, keycloak, generic
    authelia:
      issuer_url: "https://auth.example.com"
      client_id: "${REVENGE_OIDC_CLIENT_ID}"
      client_secret: "${REVENGE_OIDC_CLIENT_SECRET}"
      redirect_uri: "https://revenge.example.com/api/v1/auth/oidc/callback"
      scopes:
        - openid
        - profile
        - email
        - groups
      group_mappings:
        admins: admin
        media-users: user
        kids: restricted
      auto_provision: true
      allow_registration: true  # Allow new users via OIDC
```

---

## Database Schema

```sql
-- OIDC provider configuration (admin-managed)
CREATE TABLE oidc_providers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    provider_type VARCHAR(50) NOT NULL,  -- authelia, authentik, keycloak, generic
    issuer_url TEXT NOT NULL,
    client_id VARCHAR(255) NOT NULL,
    client_secret_encrypted BYTEA NOT NULL,
    scopes TEXT[] NOT NULL DEFAULT ARRAY['openid', 'profile', 'email'],
    enabled BOOLEAN NOT NULL DEFAULT true,
    auto_provision BOOLEAN NOT NULL DEFAULT true,
    group_mappings JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- User OIDC links
CREATE TABLE user_oidc_links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider_id UUID NOT NULL REFERENCES oidc_providers(id) ON DELETE CASCADE,
    oidc_subject VARCHAR(255) NOT NULL,
    oidc_issuer TEXT NOT NULL,
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(provider_id, oidc_subject)
);

CREATE INDEX idx_user_oidc_links_user ON user_oidc_links(user_id);
CREATE INDEX idx_user_oidc_links_subject ON user_oidc_links(oidc_subject);
```

---

## Authentication Flow

```
┌──────────┐     ┌─────────┐     ┌──────────┐     ┌─────────┐
│  User    │     │ Revenge │     │ Authelia │     │ Revenge │
│ Browser  │     │  Login  │     │   OIDC   │     │   API   │
└────┬─────┘     └────┬────┘     └────┬─────┘     └────┬────┘
     │                │               │                │
     │ Click "Login   │               │                │
     │ with Authelia" │               │                │
     │───────────────>│               │                │
     │                │               │                │
     │                │ Redirect to   │                │
     │                │ /authorize    │                │
     │<───────────────│               │                │
     │                │               │                │
     │ Login + 2FA    │               │                │
     │───────────────────────────────>│                │
     │                │               │                │
     │                │  Auth code    │                │
     │<───────────────────────────────│                │
     │                │               │                │
     │ Callback with  │               │                │
     │ auth code      │               │                │
     │───────────────>│               │                │
     │                │               │                │
     │                │ Exchange code │                │
     │                │ for tokens    │                │
     │                │──────────────>│                │
     │                │               │                │
     │                │ ID + Access   │                │
     │                │ tokens        │                │
     │                │<──────────────│                │
     │                │               │                │
     │                │ Create/update │                │
     │                │ user + session│                │
     │                │──────────────────────────────>│
     │                │               │                │
     │ Set session    │               │                │
     │ cookie         │               │                │
     │<───────────────│               │                │
     │                │               │                │
```

---

## Security Considerations

1. **Token Validation**: Always validate JWT signature using Authelia's JWKS
2. **State Parameter**: Use cryptographic state to prevent CSRF
3. **PKCE**: Use PKCE (S256) for additional security
4. **Secret Storage**: Encrypt client secret at rest
5. **HTTPS Only**: Require HTTPS for all OIDC endpoints

---

## Troubleshooting

| Issue | Cause | Solution |
|-------|-------|----------|
| Invalid redirect URI | Mismatch in Authelia config | Verify redirect_uri in both configs |
| Token validation failed | Clock skew | Sync server clocks, allow 5min leeway |
| Groups not received | Missing scope | Add `groups` scope in Authelia client |
| User not created | auto_provision disabled | Enable `auto_provision: true` |

---


## Related Documentation

- [Authentik Integration](AUTHENTIK.md)
- [Keycloak Integration](KEYCLOAK.md)
- [Generic OIDC](GENERIC_OIDC.md)
- [OIDC Authentication](../../.github/instructions/oidc-authentication.instructions.md)
