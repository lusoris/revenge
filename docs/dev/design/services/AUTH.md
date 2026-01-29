# Auth Service

> Authentication, registration, and password management

**Location**: `internal/service/auth/`

---

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

## Related Services

- [User Service](USER.md) - User account management
- [Session Service](SESSION.md) - Session token handling
- [OIDC Service](OIDC.md) - SSO authentication
