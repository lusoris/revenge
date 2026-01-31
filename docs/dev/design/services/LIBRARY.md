# Library Service

> Library management and access control

**Module**: `internal/service/library`

## Developer Resources

> Package versions: [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#go-dependencies-core)

| Package | Purpose |
|---------|---------|
| pgx | PostgreSQL driver |
| otter | Library caching |
| sqlc | Type-safe SQL queries |
| fx | Dependency injection |

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | |
| Sources | ✅ | |
| Instructions | ✅ | |
| Code | 🔴 | |
| Linting | 🔴 | |
| Unit Testing | 🔴 | |
| Integration Testing | 🔴 | |

---

## Overview

The Library service manages media libraries including:

- Library CRUD operations
- User access control (grants/revokes)
- Library type validation
- Access-aware queries

---

## Library Types

```go
const (
    LibraryTypeMovie      = "movie"
    LibraryTypeTvshow     = "tvshow"
    LibraryTypeMusic      = "music"
    LibraryTypeAudiobook  = "audiobook"
    LibraryTypeBook       = "book"
    LibraryTypePodcast    = "podcast"
    LibraryTypePhoto      = "photo"
    LibraryTypeLivetv     = "livetv"
    LibraryTypeComics     = "comics"
    LibraryTypeAdultMovie = "adult_movie"
    LibraryTypeAdultScene = "adult_scene"
)
```

---

## Operations

### Create Library

```go
type CreateParams struct {
    Name              string
    LibraryType       string
    Paths             []string
    ScanEnabled       bool
    ScanIntervalHours int32
    PreferredLanguage *string
    DownloadImages    bool
    DownloadNfo       bool
    GenerateChapters  bool
    IsPrivate         bool
    OwnerUserID       pgtype.UUID
    SortOrder         int32
    Icon              *string
}

func (s *Service) Create(ctx context.Context, params CreateParams) (*db.Library, error)
```

### Get Library

```go
// Get by ID (no access check)
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*db.Library, error)

// Get by ID with user access check
func (s *Service) GetByIDWithAccess(ctx context.Context, id, userID uuid.UUID, isAdmin bool) (*db.Library, error)
```

### List Libraries

```go
// List all (admin only)
func (s *Service) List(ctx context.Context) ([]db.Library, error)

// List accessible to user
func (s *Service) ListAccessible(ctx context.Context, userID uuid.UUID) ([]db.Library, error)

// List by type
func (s *Service) ListByType(ctx context.Context, libraryType db.LibraryType) ([]db.Library, error)

// List with access info (admin)
func (s *Service) ListAll(ctx context.Context) ([]WithAccess, error)

// List with access info (user)
func (s *Service) ListForUser(ctx context.Context, userID uuid.UUID) ([]WithAccess, error)
```

### Update Library

```go
type UpdateParams struct {
    ID                uuid.UUID
    Name              *string
    Paths             []string
    ScanEnabled       *bool
    ScanIntervalHours *int32
    PreferredLanguage *string
    DownloadImages    *bool
    DownloadNfo       *bool
    GenerateChapters  *bool
    IsPrivate         *bool
    SortOrder         *int32
    Icon              *string
}

func (s *Service) Update(ctx context.Context, params UpdateParams) (*db.Library, error)
```

### Delete Library

```go
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error
```

---

## Access Control

### Grant Access

```go
func (s *Service) GrantAccess(ctx context.Context, libraryID, userID uuid.UUID, canManage bool) error
```

### Revoke Access

```go
func (s *Service) RevokeAccess(ctx context.Context, libraryID, userID uuid.UUID) error
```

### Check Access

```go
func (s *Service) UserCanAccess(ctx context.Context, libraryID, userID uuid.UUID) (bool, error)
```

### List Users

```go
func (s *Service) ListUsers(ctx context.Context, libraryID uuid.UUID) ([]db.ListLibraryUsersRow, error)
```

---

## Access Model

```go
type WithAccess struct {
    Library   db.Library
    CanManage bool
}
```

Access rules:
- **Admins**: Access all libraries, can manage all
- **Owners**: Access and manage own private libraries
- **Granted users**: Access based on `library_access` table

---

## Errors

| Error | Description |
|-------|-------------|
| `ErrLibraryNotFound` | Library does not exist |
| `ErrAccessDenied` | User lacks access to library |
| `ErrInvalidLibraryType` | Invalid library type specified |

---

## Implementation Checklist

### Phase 1: Core Infrastructure
- [ ] Create `internal/service/library/` package structure
- [ ] Define entity types in `entity.go`
- [ ] Create repository interface
- [ ] Add fx module wiring

### Phase 2: Database
- [ ] Create migration for `libraries` table
- [ ] Create `library_access` table for permissions
- [ ] Add indexes and constraints
- [ ] Write sqlc queries

### Phase 3: Service Layer
- [ ] Implement CRUD operations with caching
- [ ] Implement access control logic
- [ ] Add library type validation
- [ ] Implement cache invalidation

### Phase 4: API Integration
- [ ] Define OpenAPI endpoints
- [ ] Generate ogen handlers
- [ ] Wire handlers to service
- [ ] Add admin/user authorization checks

---


<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

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

- [Library Types](../features/shared/LIBRARY_TYPES.md) - Per-module library design
- [RBAC](../features/shared/RBAC_CASBIN.md) - Permission model
- [Grants Service](GRANTS.md) - Fine-grained sharing
- [Metadata Service](METADATA.md) - Library metadata enrichment
- [Search Service](SEARCH.md) - Library content search
- [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md) - Service inventory
