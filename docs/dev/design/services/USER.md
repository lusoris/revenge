# User Service

> User account management and authentication


<!-- TOC-START -->

## Table of Contents

- [Developer Resources](#developer-resources)
- [Status](#status)
- [Overview](#overview)
- [Roles](#roles)
- [Operations](#operations)
  - [Create User](#create-user)
  - [Get User](#get-user)
  - [Authentication](#authentication)
  - [Check Users](#check-users)
- [Password Security](#password-security)
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

**Module**: `internal/service/user`

## Developer Resources

> Package versions: [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#go-dependencies-core)

| Package | Purpose |
|---------|---------|
| golang.org/x/crypto/bcrypt | Password hashing |
| pgx | PostgreSQL driver |
| otter | User caching |
| fx | Dependency injection |

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
---

## Overview

The User service handles all user account operations:

- User CRUD operations
- Password hashing (bcrypt)
- Role management
- Authentication validation

---

## Roles

```go
const (
    RoleAdmin     = "admin"
    RoleModerator = "moderator"
    RoleUser      = "user"
    RoleGuest     = "guest"
)
```

---

## Operations

### Create User

```go
type CreateParams struct {
    Username          string
    Email             *string
    Password          string   // Plain text (hashed internally)
    Role              string   // admin, moderator, user, guest
    IsAdmin           bool     // Deprecated: use Role
    MaxRatingLevel    int32
    AdultEnabled      bool
    PreferredLanguage *string
}

func (s *Service) Create(ctx context.Context, params CreateParams) (*db.User, error)
```

### Get User

```go
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*db.User, error)
func (s *Service) GetByUsername(ctx context.Context, username string) (*db.User, error)
```

### Authentication

```go
// Authenticate validates username/password
func (s *Service) Authenticate(ctx context.Context, username, password string) (*db.User, error)

// ValidatePassword checks password against stored hash
func (s *Service) ValidatePassword(ctx context.Context, user *db.User, password string) error

// UpdatePassword changes user's password
func (s *Service) UpdatePassword(ctx context.Context, userID uuid.UUID, newPassword string) error
```

### Check Users

```go
func (s *Service) HasAnyUsers(ctx context.Context) (bool, error)
```

---

## Password Security

- **Algorithm**: bcrypt
- **Default Cost**: 12
- **Storage**: Only hash stored, never plain text

---

## Errors

| Error | Description |
|-------|-------------|
| `ErrUserNotFound` | User does not exist |
| `ErrUserExists` | Username/email already taken |
| `ErrInvalidCredentials` | Invalid username or password |
| `ErrUserDisabled` | Account is disabled |
| `ErrInvalidRole` | Invalid role specified |

---

## Implementation Checklist

### Phase 1: Core Infrastructure
- [ ] Create `internal/service/user/` package structure
- [ ] Define entity types in `entity.go`
- [ ] Create repository interface
- [ ] Add fx module wiring

### Phase 2: Database
- [ ] Create migration for `users` table
- [ ] Add unique constraints (username, email)
- [ ] Add indexes
- [ ] Write sqlc queries

### Phase 3: Service Layer
- [ ] Implement CRUD operations with caching
- [ ] Implement bcrypt password hashing (cost 12)
- [ ] Implement authentication validation
- [ ] Add role management

### Phase 4: API Integration
- [ ] Define OpenAPI endpoints
- [ ] Generate ogen handlers
- [ ] Wire handlers to service
- [ ] Add admin authorization for user management

---


<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

### Referenced Sources

| Source | Documentation |
|--------|---------------|
| [PostgreSQL Arrays](https://www.postgresql.org/docs/current/arrays.html) | [Local](../../sources/database/postgresql-arrays.md) |
| [PostgreSQL JSON Functions](https://www.postgresql.org/docs/current/functions-json.html) | [Local](../../sources/database/postgresql-json.md) |
| [River Job Queue](https://pkg.go.dev/github.com/riverqueue/river) | [Local](../../sources/tooling/river.md) |
| [Uber fx](https://pkg.go.dev/go.uber.org/fx) | [Local](../../sources/tooling/fx.md) |
| [golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto) | [Local](../../sources/go/x/crypto.md) |
| [ogen OpenAPI Generator](https://pkg.go.dev/github.com/ogen-go/ogen) | [Local](../../sources/tooling/ogen.md) |
| [pgx PostgreSQL Driver](https://pkg.go.dev/github.com/jackc/pgx/v5) | [Local](../../sources/database/pgx.md) |
| [sqlc](https://docs.sqlc.dev/en/stable/) | [Local](../../sources/database/sqlc.md) |
| [sqlc Configuration](https://docs.sqlc.dev/en/stable/reference/config.html) | [Local](../../sources/database/sqlc-config.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Services](INDEX.md)

### In This Section

- [Activity Service](ACTIVITY.md)
- [Analytics Service](ANALYTICS.md)
- [API Keys Service](APIKEYS.md)
- [Auth Service](AUTH.md)
- [Fingerprint Service](FINGERPRINT.md)
- [Grants Service](GRANTS.md)
- [Library Service](LIBRARY.md)
- [Metadata Service](METADATA.md)

### Related Topics

- [Revenge - Architecture v2](../architecture/01_ARCHITECTURE.md) _Architecture_
- [Revenge - Design Principles](../architecture/02_DESIGN_PRINCIPLES.md) _Architecture_
- [Revenge - Metadata System](../architecture/03_METADATA_SYSTEM.md) _Architecture_
- [Revenge - Player Architecture](../architecture/04_PLAYER_ARCHITECTURE.md) _Architecture_
- [Plugin Architecture Decision](../architecture/05_PLUGIN_ARCHITECTURE_DECISION.md) _Architecture_

### Indexes

- [Design Index](../DESIGN_INDEX.md) - All design docs by category/topic
- [Source of Truth](../00_SOURCE_OF_TRUTH.md) - Package versions and status

<!-- DESIGN-BREADCRUMBS-END -->

## Related Documents

- [Auth Service](AUTH.md) - Authentication flows
- [Session Service](SESSION.md) - Session management
- [RBAC Service](RBAC.md) - Role permissions
- [Activity Service](ACTIVITY.md) - User event logging
- [OIDC Service](OIDC.md) - SSO user linking
- [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md) - Service inventory
