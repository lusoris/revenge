# Activity Service

<!-- SOURCES: fx, ogen, sqlc, sqlc-config -->

<!-- DESIGN: services, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES, 03_METADATA_SYSTEM -->


> Audit logging and event tracking


<!-- TOC-START -->

## Table of Contents

- [Developer Resources](#developer-resources)
- [Status](#status)
- [Overview](#overview)
- [Activity Types](#activity-types)
- [Severity Levels](#severity-levels)
- [Operations](#operations)
  - [Log Activity](#log-activity)
  - [Convenience Methods](#convenience-methods)
  - [Query Activities](#query-activities)
  - [Cleanup](#cleanup)
- [Metadata](#metadata)
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

**Module**: `internal/service/activity`

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

## Implementation Checklist

### Phase 1: Core Infrastructure
- [ ] Create `internal/service/activity/` package structure
- [ ] Define activity event types in `entity.go`
- [ ] Create repository interface
- [ ] Add fx module wiring

### Phase 2: Database
- [ ] Create migration for `activity_events` table
- [ ] Add partitioning by date for performance
- [ ] Add indexes (user_id, event_type, created_at)
- [ ] Write sqlc queries

### Phase 3: Service Layer
- [ ] Implement event recording
- [ ] Implement event querying with filters
- [ ] Add event aggregation
- [ ] Implement retention cleanup

### Phase 4: API Integration
- [ ] Define OpenAPI endpoints
- [ ] Generate ogen handlers
- [ ] Add admin-only endpoints

---


## Related Documents

- [Auth Service](AUTH.md) - Login/logout events
- [Library Service](LIBRARY.md) - Library events
- [User Service](USER.md) - User CRUD events
- [Session Service](SESSION.md) - Session activity tracking
- [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md) - Service inventory
