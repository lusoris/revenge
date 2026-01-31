# Auth Service

<!-- SOURCES: fx, ogen -->

<!-- DESIGN: services, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES, 03_METADATA_SYSTEM -->


> Authentication, registration, and password management


<!-- TOC-START -->

## Table of Contents

- [Developer Resources](#developer-resources)
- [Status](#status)
- [Overview](#overview)
- [Dependencies](#dependencies)
- [Operations](#operations)
  - [Login](#login)
  - [Logout](#logout)
  - [Logout All](#logout-all)
  - [Register](#register)
  - [Validate Token](#validate-token)
  - [Change Password](#change-password)
  - [Is Setup Required](#is-setup-required)
- [Errors](#errors)
- [Implementation Checklist](#implementation-checklist)
  - [Phase 1: Core Infrastructure](#phase-1-core-infrastructure)
  - [Phase 2: Service Layer](#phase-2-service-layer)
  - [Phase 3: Middleware](#phase-3-middleware)
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

**Module**: `internal/service/auth`

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

The Auth service coordinates authentication flows by combining the User and Session services. It provides:

- User login/logout
- Registration (first user becomes admin)
- Password management
- Token validation
- Setup status checking

---

## Dependencies

```go
type Service struct {
    userService    *user.Service
    sessionService *session.Service
    logger         *slog.Logger
}
```

---

## Operations

### Login

Authenticates user and creates a session.

```go
type LoginParams struct {
    Username      string
    Password      string
    DeviceName    *string
    DeviceType    *string
    ClientName    *string
    ClientVersion *string
    IPAddress     netip.Addr
    UserAgent     *string
}

type LoginResult struct {
    User    *db.User
    Session *db.Session
    Token   string  // Raw token to return to client
}

func (s *Service) Login(ctx context.Context, params LoginParams) (*LoginResult, error)
```

**Flow**:
1. Authenticate user via `userService.Authenticate()`
2. Create session via `sessionService.Create()`
3. Return user, session, and token

### Logout

Deactivates the current session.

```go
func (s *Service) Logout(ctx context.Context, token string) error
```

### Logout All

Deactivates all sessions for a user.

```go
func (s *Service) LogoutAll(ctx context.Context, userID uuid.UUID) error
```

### Register

Creates a new user account. First user registered becomes admin.

```go
type RegisterParams struct {
    Username          string
    Email             *string
    Password          string
    PreferredLanguage *string
}

func (s *Service) Register(ctx context.Context, params RegisterParams) (*db.User, error)
```

**Behavior**:
- Checks if any users exist
- First user: `isAdmin = true`, full access
- Subsequent users: `isAdmin = false`
- Default settings: `MaxRatingLevel = 100`, `AdultEnabled = false`

### Validate Token

Validates a session token and returns user + session.

```go
func (s *Service) ValidateToken(ctx context.Context, token string) (*db.User, *db.Session, error)
```

**Behavior**:
- Validates token via session service
- Fetches user by ID
- Checks if user is disabled (deactivates session if so)
- Updates session activity timestamp

### Change Password

Changes a user's password after validating current password.

```go
func (s *Service) ChangePassword(ctx context.Context, userID uuid.UUID, currentPassword, newPassword string) error
```

### Is Setup Required

Returns `true` if no users exist (initial setup needed).

```go
func (s *Service) IsSetupRequired(ctx context.Context) (bool, error)
```

---

## Errors

| Error | Description |
|-------|-------------|
| `ErrSetupRequired` | Initial setup has not been completed |
| `user.ErrUserDisabled` | User account is disabled |
| `user.ErrInvalidCredentials` | Invalid username or password |

---

## Implementation Checklist

### Phase 1: Core Infrastructure
- [ ] Create `internal/service/auth/` package structure
- [ ] Define auth types and interfaces
- [ ] Add fx module wiring

### Phase 2: Service Layer
- [ ] Implement login flow (user + session)
- [ ] Implement logout flow
- [ ] Implement password change
- [ ] Add rate limiting for auth endpoints

### Phase 3: Middleware
- [ ] Implement auth middleware
- [ ] Add session validation
- [ ] Add API key validation fallback

### Phase 4: API Integration
- [ ] Define OpenAPI endpoints
- [ ] Generate ogen handlers
- [ ] Wire handlers to service

---


## Related Documents

- [User Service](USER.md) - User account management
- [Session Service](SESSION.md) - Session token handling
- [OIDC Service](OIDC.md) - SSO authentication
- [API Keys Service](APIKEYS.md) - Programmatic access
- [Activity Service](ACTIVITY.md) - Login/logout event tracking
- [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md) - Service inventory
