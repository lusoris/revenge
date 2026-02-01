---
sources:
  - name: Uber fx
    url: ../../sources/tooling/fx.md
    note: Auto-resolved from fx
  - name: ogen OpenAPI Generator
    url: ../../sources/tooling/ogen.md
    note: Auto-resolved from ogen
  - name: sqlc
    url: ../../sources/database/sqlc.md
    note: Auto-resolved from sqlc
  - name: sqlc Configuration
    url: ../../sources/database/sqlc-config.md
    note: Auto-resolved from sqlc-config
design_refs:
  - title: services
    path: INDEX.md
  - title: 01_ARCHITECTURE
    path: ../architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: ../architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: ../architecture/03_METADATA_SYSTEM.md
---

## Table of Contents

- [Activity Service](#activity-service)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Service Structure](#service-structure)
    - [Dependencies](#dependencies)
    - [Provides](#provides)
    - [Component Diagram](#component-diagram)
  - [Implementation](#implementation)
    - [File Structure](#file-structure)
    - [Key Interfaces](#key-interfaces)
    - [Dependencies](#dependencies)
  - [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
    - [Config Keys](#config-keys)
  - [API Endpoints](#api-endpoints)
  - [Testing Strategy](#testing-strategy)
    - [Unit Tests](#unit-tests)
    - [Integration Tests](#integration-tests)
    - [Test Coverage](#test-coverage)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)


# Activity Service


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: service


> > Audit logging and event tracking

**Package**: `internal/service/activity`
**fx Module**: `activity.Module`

---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | - |
| Sources | ✅ | - |
| Instructions | ✅ | - |
| Code | 🔴 | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: ✅ Complete



---


## Architecture

### Service Structure

```
internal/service/activity/
├── module.go              # fx module definition
├── service.go             # Service implementation
├── repository.go          # Data access (if needed)
├── handler.go             # HTTP handlers (if exposed)
├── middleware.go          # Middleware (if needed)
├── types.go               # Domain types
└── service_test.go        # Tests
```

### Dependencies
**Go Packages**:
- `github.com/google/uuid`
- `github.com/jackc/pgx/v5`
- `github.com/riverqueue/river` - Cleanup jobs
- `go.uber.org/fx`


### Provides
<!-- Service provides -->

### Component Diagram

<!-- Component diagram -->


## Implementation

### File Structure

<!-- File structure -->

### Key Interfaces

```go
type ActivityService interface {
  // Logging
  Log(ctx context.Context, entry ActivityEntry) error
  LogWithContext(ctx context.Context, userID uuid.UUID, action, resourceType string, resourceID uuid.UUID, changes map[string]interface{}) error

  // Querying
  GetUserActivity(ctx context.Context, userID uuid.UUID, filters ActivityFilters) ([]ActivityEntry, error)
  GetResourceActivity(ctx context.Context, resourceType string, resourceID uuid.UUID) ([]ActivityEntry, error)
  Search(ctx context.Context, filters ActivityFilters) ([]ActivityEntry, error)

  // Cleanup
  CleanupOldLogs(ctx context.Context, olderThan time.Time) (int, error)
}

type ActivityEntry struct {
  ID           uuid.UUID              `db:"id" json:"id"`
  UserID       *uuid.UUID             `db:"user_id" json:"user_id,omitempty"`
  Username     *string                `db:"username" json:"username,omitempty"`
  Action       string                 `db:"action" json:"action"`
  ResourceType *string                `db:"resource_type" json:"resource_type,omitempty"`
  ResourceID   *uuid.UUID             `db:"resource_id" json:"resource_id,omitempty"`
  Changes      map[string]interface{} `db:"changes" json:"changes,omitempty"`
  IPAddress    *net.IP                `db:"ip_address" json:"ip_address,omitempty"`
  Success      bool                   `db:"success" json:"success"`
  CreatedAt    time.Time              `db:"created_at" json:"created_at"`
}
```


### Dependencies

**Go Packages**:
- `github.com/google/uuid`
- `github.com/jackc/pgx/v5`
- `github.com/riverqueue/river` - Cleanup jobs
- `go.uber.org/fx`






## Configuration
### Environment Variables

```bash
ACTIVITY_RETENTION_DAYS=90
ACTIVITY_CLEANUP_INTERVAL=24h
```


### Config Keys

```yaml
activity:
  retention_days: 90
  cleanup_interval: 24h
  log_failed_attempts: true
```



## API Endpoints
```
GET    /api/v1/activity                # Search activity logs
GET    /api/v1/activity/users/:id      # Get user activity
GET    /api/v1/activity/resources/:type/:id # Get resource activity
```

**Example Response**:
```json
{
  "entries": [
    {
      "id": "uuid-123",
      "user_id": "uuid-456",
      "username": "admin",
      "action": "settings.update",
      "resource_type": "setting",
      "resource_id": "uuid-789",
      "changes": {
        "server.name": {
          "old": "Revenge",
          "new": "My Server"
        }
      },
      "ip_address": "192.168.1.100",
      "success": true,
      "created_at": "2026-02-01T10:00:00Z"
    }
  ],
  "total": 1,
  "page": 1
}
```



## Testing Strategy

### Unit Tests

<!-- Unit test strategy -->

### Integration Tests

<!-- Integration test strategy -->

### Test Coverage

Target: **80% minimum**







## Related Documentation
### Design Documents
- [services](INDEX.md)
- [01_ARCHITECTURE](../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../architecture/03_METADATA_SYSTEM.md)

### External Sources
- [Uber fx](../../sources/tooling/fx.md) - Auto-resolved from fx
- [ogen OpenAPI Generator](../../sources/tooling/ogen.md) - Auto-resolved from ogen
- [sqlc](../../sources/database/sqlc.md) - Auto-resolved from sqlc
- [sqlc Configuration](../../sources/database/sqlc-config.md) - Auto-resolved from sqlc-config

