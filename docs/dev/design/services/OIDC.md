# OIDC Service

<!-- SOURCES: authelia, authentik, fx, keycloak, ogen, sqlc, sqlc-config -->

<!-- DESIGN: services, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES, 03_METADATA_SYSTEM -->


> OpenID Connect / SSO provider management


<!-- TOC-START -->

## Table of Contents

- [Developer Resources](#developer-resources)
- [Status](#status)
- [Overview](#overview)
- [Provider Configuration](#provider-configuration)
  - [Claim Mapping](#claim-mapping)
  - [Role Mapping](#role-mapping)
- [Operations](#operations)
  - [Provider Management](#provider-management)
  - [User Links](#user-links)
  - [Public Info](#public-info)
- [Authentication Flow](#authentication-flow)
- [Errors](#errors)
- [Implementation Checklist](#implementation-checklist)
  - [Phase 1: Core Infrastructure](#phase-1-core-infrastructure)
  - [Phase 2: Database](#phase-2-database)
  - [Phase 3: Service Layer](#phase-3-service-layer)
  - [Phase 4: API Integration](#phase-4-api-integration)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related Documents](#related-documents)

<!-- TOC-END -->

**Module**: `internal/service/oidc`

## Developer Resources

> See [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#backend-services) for service inventory and status.

## Status

| Dimension | Status |
|-----------|--------|
| Design | ✅ |
| Sources | ✅ |
| Instructions | ✅ |
| Code | 🔴 |
| Linting | 🔴 |
| Unit Testing | 🔴 |
| Integration Testing | 🔴 |---

## Overview

The OIDC service manages SSO authentication:

- Provider configuration (Authelia, Authentik, Keycloak, etc.)
- User link management (connecting OIDC identities to local users)
- Claim mapping (mapping OIDC claims to user attributes)
- Role mapping (mapping OIDC groups to roles)
- Auto-provisioning of new users

---

## Provider Configuration

### Claim Mapping

Maps OIDC claims to user attributes:

```go
type ClaimMapping struct {
    Username string `json:"username"` // default: preferred_username
    Email    string `json:"email"`    // default: email
    Name     string `json:"name"`     // default: name
    Groups   string `json:"groups"`   // default: groups
}
```

### Role Mapping

Maps OIDC groups to admin access:

```go
type RoleMapping struct {
    AdminGroups []string `json:"admin_groups"`
}
```

---

## Operations

### Provider Management

```go
type CreateProviderParams struct {
    Name            string
    Slug            string          // URL-safe identifier
    Enabled         bool
    IssuerURL       string
    ClientID        string
    ClientSecretEnc []byte          // Encrypted
    Scopes          []string        // default: openid, profile, email
    ClaimMapping    *ClaimMapping
    RoleMapping     *RoleMapping
    AutoProvision   bool            // Create users on first login
    DefaultRole     string          // user or admin
}

func (s *Service) CreateProvider(ctx context.Context, params CreateProviderParams) (*db.OidcProvider, error)
func (s *Service) GetProviderByID(ctx context.Context, id uuid.UUID) (*db.OidcProvider, error)
func (s *Service) GetProviderBySlug(ctx context.Context, slug string) (*db.OidcProvider, error)
func (s *Service) GetEnabledProvider(ctx context.Context, slug string) (*db.OidcProvider, error)
func (s *Service) ListProviders(ctx context.Context) ([]db.OidcProvider, error)
func (s *Service) ListEnabledProviders(ctx context.Context) ([]db.OidcProvider, error)
func (s *Service) DeleteProvider(ctx context.Context, id uuid.UUID) error
```

### User Links

Links OIDC identities (subject) to local users:

```go
func (s *Service) GetUserLink(ctx context.Context, providerID uuid.UUID, subject string) (*db.OidcUserLink, error)
func (s *Service) ListUserLinks(ctx context.Context, userID uuid.UUID) ([]db.ListOIDCLinksByUserRow, error)
func (s *Service) CreateLink(ctx context.Context, userID, providerID uuid.UUID, subject string, email, name *string, groups []string) (*db.OidcUserLink, error)
func (s *Service) UpdateLinkLogin(ctx context.Context, linkID uuid.UUID, email, name *string, groups []string) error
func (s *Service) DeleteLink(ctx context.Context, linkID uuid.UUID) error
func (s *Service) DeleteUserLinks(ctx context.Context, userID uuid.UUID) error
```

### Public Info

For login page (no sensitive data):

```go
type ProviderInfo struct {
    Name string `json:"name"`
    Slug string `json:"slug"`
}

func (s *Service) GetPublicProviders(ctx context.Context) ([]ProviderInfo, error)
```

---

## Authentication Flow

```
1. User clicks "Login with {Provider}"
2. Redirect to provider authorize URL
3. Provider authenticates user
4. Callback with authorization code
5. Exchange code for tokens
6. Extract claims from ID token
7. Look up or create user link
8. Create local session
9. Return session token to client
```

---

## Errors

| Error | Description |
|-------|-------------|
| `ErrProviderNotFound` | Provider does not exist |
| `ErrProviderDisabled` | Provider is disabled |
| `ErrLinkNotFound` | No link for this user/provider |
| `ErrAutoProvisionDisabled` | Cannot create new user |
| `ErrSlugTaken` | Provider slug already exists |

---

## Implementation Checklist

### Phase 1: Core Infrastructure
- [ ] Create `internal/service/oidc/` package structure
- [ ] Define OIDC provider configuration
- [ ] Add fx module wiring

### Phase 2: Database
- [ ] Create migration for `oidc_providers` table
- [ ] Create `user_oidc_links` table
- [ ] Write sqlc queries

### Phase 3: Service Layer
- [ ] Implement provider registration
- [ ] Implement authorization URL generation
- [ ] Implement callback handling
- [ ] Implement user linking/creation

### Phase 4: API Integration
- [ ] Define OpenAPI endpoints for OIDC flow
- [ ] Generate ogen handlers
- [ ] Wire handlers to service

---


## Related Documents

- [Auth Service](AUTH.md) - Local authentication
- [User Service](USER.md) - User management
- [Session Service](SESSION.md) - Session token handling
- [Integrations: Auth](../integrations/auth/) - Provider-specific docs
- [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md) - Service inventory
