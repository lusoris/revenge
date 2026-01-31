# User Service

<!-- SOURCES: fx, golang-x-crypto, ogen, pgx, postgresql-arrays, postgresql-json, river, sqlc, sqlc-config -->

<!-- DESIGN: services, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES, 03_METADATA_SYSTEM -->


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
| Integration Testing | 🔴 |---

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


## Related Documents

- [Auth Service](AUTH.md) - Authentication flows
- [Session Service](SESSION.md) - Session management
- [RBAC Service](RBAC.md) - Role permissions
- [Activity Service](ACTIVITY.md) - User event logging
- [OIDC Service](OIDC.md) - SSO user linking
- [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md) - Service inventory
