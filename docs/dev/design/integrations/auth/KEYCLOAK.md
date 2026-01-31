# Keycloak Integration

> Enterprise identity and access management solution


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Overview](#overview)
- [Developer Resources](#developer-resources)
- [OIDC Details](#oidc-details)
  - [Supported Scopes](#supported-scopes)
  - [Keycloak Configuration](#keycloak-configuration)
- [Data Mapping](#data-mapping)
  - [Keycloak Claims → Revenge User](#keycloak-claims-revenge-user)
  - [Role Mapping Strategies](#role-mapping-strategies)
- [Implementation Checklist](#implementation-checklist)
- [Configuration](#configuration)
- [Advanced Features](#advanced-features)
  - [Fine-Grained Authorization](#fine-grained-authorization)
  - [User Federation (LDAP)](#user-federation-ldap)
  - [Social Identity Providers](#social-identity-providers)
- [Database Schema](#database-schema)
- [Keycloak vs Others](#keycloak-vs-others)
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
| Integration Testing | 🔴 |
**Priority**: 🟢 LOW (Phase 1 - Core Infrastructure)
**Type**: OIDC Identity Provider

---

## Overview

Keycloak is a mature, enterprise-grade identity and access management solution by Red Hat. Revenge integrates with Keycloak as an OIDC provider for:
- User authentication via OIDC
- Single Sign-On with enterprise features
- Federation with LDAP/Active Directory
- Fine-grained authorization policies
- Multi-realm support

**Integration Points**:
- **OIDC Provider**: Standard OpenID Connect flow
- **User Federation**: LDAP, Kerberos, social logins
- **Role mapping**: Realm/client roles to Revenge roles
- **Authorization**: Fine-grained resource permissions

---

## Developer Resources

- 📚 **Docs**: https://www.keycloak.org/documentation
- 🔗 **OIDC Docs**: https://www.keycloak.org/docs/latest/securing_apps/#_oidc
- 🔗 **REST API**: https://www.keycloak.org/docs-api/latest/rest-api/
- 🔗 **GitHub**: https://github.com/keycloak/keycloak

---

## OIDC Details

**Discovery URL**: `https://keycloak.example.com/realms/{realm}/.well-known/openid-configuration`
**Authorization Endpoint**: `https://keycloak.example.com/realms/{realm}/protocol/openid-connect/auth`
**Token Endpoint**: `https://keycloak.example.com/realms/{realm}/protocol/openid-connect/token`
**UserInfo Endpoint**: `https://keycloak.example.com/realms/{realm}/protocol/openid-connect/userinfo`
**JWKS URI**: `https://keycloak.example.com/realms/{realm}/protocol/openid-connect/certs`
**Logout Endpoint**: `https://keycloak.example.com/realms/{realm}/protocol/openid-connect/logout`

### Supported Scopes

| Scope | Claims |
|-------|--------|
| `openid` | `sub`, `iss`, `aud`, `exp`, `iat` |
| `profile` | `name`, `preferred_username`, `given_name`, `family_name` |
| `email` | `email`, `email_verified` |
| `roles` | `realm_access.roles`, `resource_access.{client}.roles` |
| `groups` | `groups` (requires mapper) |

### Keycloak Configuration

1. **Create Realm** (or use existing):
   - Name: `media` (or your realm name)

2. **Create Client**:
   - Client ID: `revenge`
   - Client Protocol: `openid-connect`
   - Access Type: `confidential`
   - Valid Redirect URIs: `https://revenge.example.com/api/v1/auth/oidc/callback`
   - Web Origins: `https://revenge.example.com`

3. **Configure Client**:
   - Standard Flow Enabled: `ON`
   - Direct Access Grants: `OFF` (more secure)
   - Service Accounts: `OFF` (unless needed for API)

4. **Create Protocol Mappers** (for groups):
   - Name: `groups`
   - Mapper Type: `Group Membership`
   - Token Claim Name: `groups`
   - Full group path: `OFF`
   - Add to ID token: `ON`
   - Add to access token: `ON`
   - Add to userinfo: `ON`

5. **Create Roles**:
   - Realm Roles: `revenge-admin`, `revenge-user`, `revenge-restricted`
   - Or Client Roles under `revenge` client

---

## Data Mapping

### Keycloak Claims → Revenge User

| Keycloak Claim | Revenge Field | Notes |
|----------------|---------------|-------|
| `sub` | `oidc_subject` | UUID |
| `preferred_username` | `username` | Display name |
| `email` | `email` | User email |
| `email_verified` | `email_verified` | Verification status |
| `name` | `display_name` | Full name |
| `realm_access.roles` | `roles[]` | Realm-level roles |
| `resource_access.revenge.roles` | `roles[]` | Client-specific roles |
| `groups` | `groups[]` | Group membership (if mapper added) |

### Role Mapping Strategies

**Option 1: Realm Roles**
```yaml
group_mappings:
  revenge-admin: admin
  revenge-user: user
  revenge-restricted: restricted
```

**Option 2: Client Roles**
```yaml
group_mappings:
  # Uses resource_access.revenge.roles
  admin: admin
  user: user
  restricted: restricted
use_client_roles: true
```

**Option 3: Groups**
```yaml
# Requires group mapper in Keycloak
group_mappings:
  /media/admins: admin
  /media/users: user
```

---

## Implementation Checklist

- [ ] **OIDC Client** (`internal/service/oidc/provider_keycloak.go`)
  - [ ] Discovery document parsing
  - [ ] Authorization code flow with PKCE
  - [ ] Token validation (JWT)
  - [ ] Role extraction (realm + client roles)
  - [ ] Token refresh handling
  - [ ] Logout URL support

- [ ] **User Provisioning**
  - [ ] Auto-create user on first login
  - [ ] Map OIDC claims to user fields
  - [ ] Support realm roles, client roles, and groups
  - [ ] Handle user updates

- [ ] **Session Management**
  - [ ] Create Revenge session from OIDC token
  - [ ] Support Keycloak single logout
  - [ ] Back-channel logout (optional)

---

## Configuration

```yaml
# configs/config.yaml
integrations:
  oidc:
    enabled: true
    provider: "keycloak"
    keycloak:
      issuer_url: "https://keycloak.example.com/realms/media"
      client_id: "${REVENGE_OIDC_CLIENT_ID}"
      client_secret: "${REVENGE_OIDC_CLIENT_SECRET}"
      redirect_uri: "https://revenge.example.com/api/v1/auth/oidc/callback"
      scopes:
        - openid
        - profile
        - email
        - roles
      role_source: "realm"  # realm, client, or groups
      role_mappings:
        revenge-admin: admin
        revenge-user: user
      auto_provision: true
      allow_registration: true
      logout_redirect: "https://revenge.example.com"
```

---

## Advanced Features

### Fine-Grained Authorization

Keycloak supports resource-based authorization:

```json
// Keycloak authorization settings
{
  "resources": [
    {
      "name": "library",
      "type": "revenge:library",
      "scopes": ["read", "write", "delete"]
    }
  ],
  "policies": [
    {
      "name": "admin-policy",
      "type": "role",
      "roles": ["revenge-admin"]
    }
  ],
  "permissions": [
    {
      "name": "library-admin-permission",
      "type": "resource",
      "resources": ["library"],
      "policies": ["admin-policy"],
      "scopes": ["read", "write", "delete"]
    }
  ]
}
```

### User Federation (LDAP)

Keycloak can federate users from LDAP/AD:

1. Add User Federation > LDAP
2. Configure connection to your LDAP server
3. Map LDAP groups to Keycloak roles
4. Users automatically available for Revenge SSO

### Social Identity Providers

Configure social logins in Keycloak:
- Google, GitHub, Facebook, etc.
- Users login via social → Keycloak → Revenge

---

## Database Schema

Uses shared OIDC tables from [Authelia Integration](AUTHELIA.md#database-schema).

---

## Keycloak vs Others

| Feature | Keycloak | Authentik | Authelia |
|---------|----------|-----------|----------|
| Complexity | Highest | High | Low |
| Enterprise features | Full | Partial | Basic |
| LDAP federation | Excellent | Good | External |
| Resource usage | High | Medium | Low |
| Learning curve | Steep | Moderate | Easy |
| Best for | Enterprise | Self-hosted enterprise | Simple SSO |

**Recommendation**: Use Keycloak when:
- LDAP/AD integration required
- Complex authorization policies needed
- Already using Keycloak in organization
- Enterprise audit/compliance requirements

---

## Troubleshooting

| Issue | Cause | Solution |
|-------|-------|----------|
| Invalid realm | Wrong realm in URL | Check realm name (case-sensitive) |
| Invalid client | Wrong client_id | Verify in Keycloak client settings |
| No roles in token | Missing mapper | Add roles mapper or use correct scope |
| CORS errors | Missing web origin | Add origin in client settings |
| Token expired | Short token lifespan | Implement token refresh |

---


<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

### Referenced Sources

| Source | Documentation |
|--------|---------------|
| [Authelia Documentation](https://www.authelia.com/overview/) | [Local](../../../sources/security/authelia.md) |
| [Authentik Documentation](https://goauthentik.io/docs/) | [Local](../../../sources/security/authentik.md) |
| [Keycloak Documentation](https://www.keycloak.org/documentation) | [Local](../../../sources/security/keycloak.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Auth](INDEX.md)

### In This Section

- [Authelia Integration](AUTHELIA.md)
- [Authentik Integration](AUTHENTIK.md)
- [Generic OIDC Integration](GENERIC_OIDC.md)

### Related Topics

- [Revenge - Architecture v2](../../architecture/01_ARCHITECTURE.md) _Architecture_
- [Revenge - Design Principles](../../architecture/02_DESIGN_PRINCIPLES.md) _Architecture_
- [Revenge - Metadata System](../../architecture/03_METADATA_SYSTEM.md) _Architecture_
- [Revenge - Player Architecture](../../architecture/04_PLAYER_ARCHITECTURE.md) _Architecture_
- [Plugin Architecture Decision](../../architecture/05_PLUGIN_ARCHITECTURE_DECISION.md) _Architecture_

### Indexes

- [Design Index](../../DESIGN_INDEX.md) - All design docs by category/topic
- [Source of Truth](../../00_SOURCE_OF_TRUTH.md) - Package versions and status

<!-- DESIGN-BREADCRUMBS-END -->

## Related Documentation

- [Authelia Integration](AUTHELIA.md)
- [Authentik Integration](AUTHENTIK.md)
- [Generic OIDC](GENERIC_OIDC.md)
