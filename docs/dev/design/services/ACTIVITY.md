# Activity Service

> Audit logging and event tracking

**Location**: `internal/service/activity/`

---

## Overview

The Activity service provides comprehensive audit logging:

- User actions (login, logout, settings changes)
- Library events (created, updated, scanned)
- Content events (played, rated)
- Security events and API errors
- Queryable activity history

---

## Activity Types

```go
const (
    TypeUserLogin       = "user_login"
    TypeUserLogout      = "user_logout"
    TypeUserCreated     = "user_created"
    TypeUserUpdated     = "user_updated"
    TypeUserDeleted     = "user_deleted"
    TypePasswordChanged = "password_changed"
    TypeSessionCreated  = "session_created"
    TypeSessionExpired  = "session_expired"
    TypeLibraryCreated  = "library_created"
    TypeLibraryUpdated  = "library_updated"
    TypeLibraryDeleted  = "library_deleted"
    TypeLibraryScanned  = "library_scanned"
    TypeContentPlayed   = "content_played"
    TypeContentRated    = "content_rated"
    TypeSettingsChanged = "settings_changed"
    TypeAPIError        = "api_error"
    TypeSecurityEvent   = "security_event"
)
```

## Severity Levels

```go
const (
    SeverityInfo     = "info"
    SeverityWarning  = "warning"
    SeverityError    = "error"
    SeverityCritical = "critical"
)
```

---

## Operations

### Log Activity

```go
type LogParams struct {
    UserID    *uuid.UUID
    Type      string
    Severity  string
    Message   string
    Metadata  map[string]any
    IPAddress netip.Addr
    UserAgent *string
}

func (s *Service) Log(ctx context.Context, params LogParams) (*db.ActivityLog, error)
```

### Convenience Methods

```go
func (s *Service) LogUserLogin(ctx context.Context, userID uuid.UUID, ip netip.Addr, userAgent *string) error
func (s *Service) LogUserLogout(ctx context.Context, userID uuid.UUID) error
func (s *Service) LogSecurityEvent(ctx context.Context, userID *uuid.UUID, message string, metadata map[string]any, ip netip.Addr) error
func (s *Service) LogAPIError(ctx context.Context, userID *uuid.UUID, message string, metadata map[string]any) error
```

### Query Activities

```go
func (s *Service) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]db.ActivityLog, error)
func (s *Service) ListByType(ctx context.Context, activityType string, limit, offset int32) ([]db.ActivityLog, error)
func (s *Service) ListBySeverity(ctx context.Context, severity string, limit, offset int32) ([]db.ActivityLog, error)
func (s *Service) ListRecent(ctx context.Context, limit, offset int32) ([]db.ActivityLog, error)
```

### Cleanup

```go
func (s *Service) DeleteOlderThan(ctx context.Context, before time.Time) error
```

---

## Metadata

Activities support arbitrary JSON metadata:

```go
s.Log(ctx, LogParams{
    UserID:   &userID,
    Type:     TypeLibraryScanned,
    Severity: SeverityInfo,
    Message:  "Library scan completed",
    Metadata: map[string]any{
        "library_id":  libraryID,
        "files_added": 150,
        "files_updated": 23,
        "duration_ms": 4500,
    },
})
```

---

## Related

- [Auth Service](AUTH.md) - Login/logout events
- [Library Service](LIBRARY.md) - Library events
