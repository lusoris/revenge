# Authentik Integration

<!-- SOURCES: authelia, authentik, keycloak -->

<!-- DESIGN: integrations/auth, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES, 03_METADATA_SYSTEM -->


> Enterprise-grade identity provider for self-hosted environments


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Overview](#overview)
- [Developer Resources](#developer-resources)
- [OIDC Details](#oidc-details)
  - [Supported Scopes](#supported-scopes)
  - [Authentik Configuration](#authentik-configuration)
- [Data Mapping](#data-mapping)
  - [Authentik Claims → Revenge User](#authentik-claims-revenge-user)
  - [Group → Role Mapping](#group-role-mapping)
- [Implementation Checklist](#implementation-checklist)
- [Configuration](#configuration)
- [Advanced Features](#advanced-features)
  - [Property Mappings](#property-mappings)
  - [Outpost (Proxy)](#outpost-proxy)
  - [API Authentication](#api-authentication)
- [Database Schema](#database-schema)
- [Authentik vs Authelia](#authentik-vs-authelia)
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

Authentik is a powerful, open-source identity provider with advanced features like user flows, policies, and application proxies. Revenge integrates with Authentik as an OIDC provider for:
- User authentication via OIDC
- Single Sign-On across media stack
- Multi-factor authentication
- Advanced user flows and policies
- Application-level access control

**Integration Points**:
- **OIDC Provider**: Standard OpenID Connect flow
- **Application Proxy**: Optional reverse proxy authentication
- **User provisioning**: Auto-create Revenge users from OIDC claims
- **Groups/Roles**: Map Authentik groups to Revenge roles

---

## Developer Resources

- 📚 **Docs**: https://goauthentik.io/docs/
- 🔗 **OIDC Provider Docs**: https://goauthentik.io/docs/providers/oauth2/
- 🔗 **API Docs**: https://goauthentik.io/developer-docs/api/
- 🔗 **GitHub**: https://github.com/goauthentik/authentik

---

## OIDC Details

**Discovery URL**: `https://authentik.example.com/application/o/<app-slug>/.well-known/openid-configuration`
**Authorization Endpoint**: `https://authentik.example.com/application/o/authorize/`
**Token Endpoint**: `https://authentik.example.com/application/o/token/`
**UserInfo Endpoint**: `https://authentik.example.com/application/o/userinfo/`
**JWKS URI**: `https://authentik.example.com/application/o/<app-slug>/jwks/`

### Supported Scopes

| Scope | Claims |
|-------|--------|
| `openid` | `sub`, `iss`, `aud` |
| `profile` | `name`, `preferred_username`, `given_name`, `family_name`, `nickname` |
| `email` | `email`, `email_verified` |
| `groups` | `groups` (array) |
| `ak_proxy` | Proxy-specific claims |

### Authentik Configuration

1. **Create OAuth2/OIDC Provider**:
   - Name: `Revenge OIDC`
   - Authorization flow: `default-provider-authorization-explicit-consent`
   - Client ID: (auto-generated or custom)
   - Client Secret: (generated)
   - Redirect URIs: `https://revenge.example.com/api/v1/auth/oidc/callback`
   - Scopes: `openid`, `profile`, `email`, `groups`
   - Subject mode: Based on User's username

2. **Create Application**:
   - Name: `Revenge Media Server`
   - Slug: `revenge`
   - Provider: `Revenge OIDC`
   - Launch URL: `https://revenge.example.com`

3. **Assign Groups** (optional):
   - Create/use groups: `revenge-admins`, `revenge-users`
   - Assign users to groups

---

## Data Mapping

### Authentik Claims → Revenge User

| Authentik Claim | Revenge Field | Notes |
|-----------------|---------------|-------|
| `sub` | `oidc_subject` | UUID or username based on config |
| `preferred_username` | `username` | Display name |
| `email` | `email` | User email |
| `email_verified` | `email_verified` | Verification status |
| `name` | `display_name` | Full name |
| `groups` | `roles[]` | Role mapping |
| `ak_groups` | `roles[]` | Alternative groups claim |

### Group → Role Mapping

```yaml
integrations:
  oidc:
    authentik:
      group_mappings:
        - group: "revenge-admins"
          role: "admin"
        - group: "revenge-users"
          role: "user"
        - group: "revenge-kids"
          role: "restricted"
```

---

## Implementation Checklist

- [ ] **OIDC Client** (`internal/service/oidc/provider_authentik.go`)
  - [ ] Discovery document parsing
  - [ ] Authorization code flow with PKCE
  - [ ] Token validation (JWT)
  - [ ] UserInfo fetching
  - [ ] Token refresh handling

- [ ] **User Provisioning** (`internal/service/user/oidc_provisioning.go`)
  - [ ] Auto-create user on first login
  - [ ] Map OIDC claims to user fields
  - [ ] Map groups to roles
  - [ ] Handle user updates

- [ ] **Session Management**
  - [ ] Create Revenge session from OIDC token
  - [ ] Session expiry based on token expiry
  - [ ] Support Authentik logout URL

---

## Configuration

```yaml
# configs/config.yaml
integrations:
  oidc:
    enabled: true
    provider: "authentik"
    authentik:
      issuer_url: "https://authentik.example.com/application/o/revenge/"
      client_id: "${REVENGE_OIDC_CLIENT_ID}"
      client_secret: "${REVENGE_OIDC_CLIENT_SECRET}"
      redirect_uri: "https://revenge.example.com/api/v1/auth/oidc/callback"
      scopes:
        - openid
        - profile
        - email
        - groups
      group_mappings:
        revenge-admins: admin
        revenge-users: user
      auto_provision: true
      allow_registration: true
```

---

## Advanced Features

### Property Mappings

Authentik allows custom property mappings to include additional claims:

```python
# Authentik property mapping for Revenge
return {
    "revenge_role": user.group_attributes().get("revenge_role", "user"),
    "revenge_allowed_libraries": user.group_attributes().get("libraries", []),
}
```

### Outpost (Proxy)

Authentik Outpost can protect Revenge at the reverse proxy level:

```yaml
# docker-compose.yml
services:
  authentik-proxy:
    image: ghcr.io/goauthentik/proxy
    environment:
      AUTHENTIK_HOST: https://authentik.example.com
      AUTHENTIK_TOKEN: ${AUTHENTIK_PROXY_TOKEN}
    ports:
      - "9000:9000"
```

### API Authentication

Use Authentik tokens for API access:

```bash
# Get API token from Authentik
curl -X POST https://authentik.example.com/application/o/token/ \
  -d "grant_type=client_credentials" \
  -d "client_id=${CLIENT_ID}" \
  -d "client_secret=${CLIENT_SECRET}"
```

---

## Database Schema

Uses shared OIDC tables from [Authelia Integration](AUTHELIA.md#database-schema).

---

## Authentik vs Authelia

| Feature | Authentik | Authelia |
|---------|-----------|----------|
| Complexity | Higher | Lower |
| UI | Full admin UI | Minimal |
| User management | Built-in | Requires external |
| Policies | Advanced | Basic |
| MFA | Extensive | TOTP, WebAuthn |
| Resource usage | Higher | Lower |
| Best for | Enterprise | Simple SSO |

**Recommendation**: Use Authelia for simple setups, Authentik for complex enterprise needs.

---

## Troubleshooting

| Issue | Cause | Solution |
|-------|-------|----------|
| 404 on discovery URL | Wrong app slug | Check application slug in Authentik |
| Invalid client | Wrong client_id | Verify in Authentik provider settings |
| No groups in token | Missing scope | Add `groups` scope in Authentik |
| PKCE error | S256 required | Ensure PKCE with S256 challenge |

---


## Related Documentation

- [Authelia Integration](AUTHELIA.md)
- [Keycloak Integration](KEYCLOAK.md)
- [Generic OIDC](GENERIC_OIDC.md)
