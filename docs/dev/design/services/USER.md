# User Service

> User account management and authentication

**Location**: `internal/service/user/`

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

## Related

- [Auth Service](AUTH.md) - Authentication flows
- [Session Service](SESSION.md) - Session management
- [RBAC](RBAC.md) - Role permissions
