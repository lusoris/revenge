# Session Service

> Session token management and device tracking

**Location**: `internal/service/session/`

---

## Overview

The Session service manages user sessions:

- Token generation and validation
- Device tracking (name, type, client info)
- Session expiration
- Activity updates
- Session limits per user

---

## Configuration

```go
type Service struct {
    queries            *db.Queries
    logger             *slog.Logger
    sessionDuration    time.Duration  // Default: 24h
    maxSessionsPerUser int            // 0 = unlimited
}

func (s *Service) SetSessionDuration(d time.Duration)
func (s *Service) SetMaxSessionsPerUser(maxSessions int)
```

---

## Operations

### Create Session

```go
type CreateParams struct {
    UserID        uuid.UUID
    ProfileID     *uuid.UUID
    DeviceName    *string
    DeviceType    *string
    ClientName    *string
    ClientVersion *string
    IPAddress     netip.Addr
    UserAgent     *string
}

type CreateResult struct {
    Session *db.Session
    Token   string  // Raw token - only returned on creation
}

func (s *Service) Create(ctx context.Context, params CreateParams) (*CreateResult, error)
```

### Validate Token

```go
func (s *Service) ValidateToken(ctx context.Context, token string) (*db.Session, error)
```

Checks:
1. Token hash exists in database
2. Session is active
3. Session not expired

### Deactivate Sessions

```go
// Single session
func (s *Service) Deactivate(ctx context.Context, sessionID uuid.UUID) error

// All sessions for user
func (s *Service) DeactivateAllForUser(ctx context.Context, userID uuid.UUID) error
```

### Update Activity

```go
func (s *Service) UpdateActivity(ctx context.Context, sessionID uuid.UUID, ipAddress *netip.Addr) error
```

---

## Token Security

- **Generation**: 32 bytes random via `crypto/rand`
- **Encoding**: Base64 URL-safe
- **Storage**: SHA-256 hash only (raw token never stored)
- **Lookup**: Hash-based lookup for O(1) validation

```go
// Token generation
tokenBytes := make([]byte, 32)
rand.Read(tokenBytes)
token := base64.URLEncoding.EncodeToString(tokenBytes)

// Storage: hash only
hash := sha256.Sum256([]byte(token))
tokenHash := base64.URLEncoding.EncodeToString(hash[:])
```

---

## Errors

| Error | Description |
|-------|-------------|
| `ErrSessionNotFound` | Session not found or invalid token |
| `ErrSessionExpired` | Session has expired |
| `ErrSessionInactive` | Session was deactivated |
| `ErrTooManySessions` | Max sessions per user exceeded |

---

## Related

- [Auth Service](AUTH.md) - Login/logout flows
- [User Service](USER.md) - User accounts
