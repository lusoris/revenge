# Activity Service

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
| Integration Testing | 🔴 |
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


<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

### Referenced Sources

| Source | Documentation |
|--------|---------------|
| [Uber fx](https://pkg.go.dev/go.uber.org/fx) | [Local](../../sources/tooling/fx.md) |
| [ogen OpenAPI Generator](https://pkg.go.dev/github.com/ogen-go/ogen) | [Local](../../sources/tooling/ogen.md) |
| [sqlc](https://docs.sqlc.dev/en/stable/) | [Local](../../sources/database/sqlc.md) |
| [sqlc Configuration](https://docs.sqlc.dev/en/stable/reference/config.html) | [Local](../../sources/database/sqlc-config.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Services](INDEX.md)

### In This Section

- [Analytics Service](ANALYTICS.md)
- [API Keys Service](APIKEYS.md)
- [Auth Service](AUTH.md)
- [Fingerprint Service](FINGERPRINT.md)
- [Grants Service](GRANTS.md)
- [Library Service](LIBRARY.md)
- [Metadata Service](METADATA.md)
- [Notification Service](NOTIFICATION.md)

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

- [Auth Service](AUTH.md) - Login/logout events
- [Library Service](LIBRARY.md) - Library events
- [User Service](USER.md) - User CRUD events
- [Session Service](SESSION.md) - Session activity tracking
- [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md) - Service inventory
