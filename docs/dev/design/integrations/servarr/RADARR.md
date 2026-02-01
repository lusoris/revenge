

---
sources:
  - name: Uber fx
    url: ../../../sources/tooling/fx.md
    note: Auto-resolved from fx
  - name: Go context
    url: ../../../sources/go/stdlib/context.md
    note: Auto-resolved from go-context
  - name: pgx PostgreSQL Driver
    url: ../../../sources/database/pgx.md
    note: Auto-resolved from pgx
  - name: PostgreSQL Arrays
    url: ../../../sources/database/postgresql-arrays.md
    note: Auto-resolved from postgresql-arrays
  - name: PostgreSQL JSON Functions
    url: ../../../sources/database/postgresql-json.md
    note: Auto-resolved from postgresql-json
  - name: Radarr API Docs
    url: ../../../sources/apis/radarr-docs.md
    note: Auto-resolved from radarr-docs
  - name: River Job Queue
    url: ../../../sources/tooling/river.md
    note: Auto-resolved from river
  - name: Servarr Wiki
    url: ../../../sources/apis/servarr-wiki.md
    note: Auto-resolved from servarr-wiki
  - name: Typesense API
    url: ../../../sources/infrastructure/typesense.md
    note: Auto-resolved from typesense
  - name: Typesense Go Client
    url: ../../../sources/infrastructure/typesense-go.md
    note: Auto-resolved from typesense-go
design_refs:
  - title: 01_ARCHITECTURE
    path: ../../architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: ../../architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: ../../architecture/03_METADATA_SYSTEM.md
---

## Table of Contents

- [Radarr](#radarr)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Integration Structure](#integration-structure)
    - [Data Flow](#data-flow)
    - [Provides](#provides)
  - [Implementation](#implementation)
    - [File Structure](#file-structure)
    - [Key Interfaces](#key-interfaces)
    - [Dependencies](#dependencies)
  - [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
- [Radarr instance](#radarr-instance)
- [Sync settings](#sync-settings)
    - [Config Keys](#config-keys)
  - [Testing Strategy](#testing-strategy)
    - [Unit Tests](#unit-tests)
    - [Integration Tests](#integration-tests)
    - [Test Coverage](#test-coverage)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)


# Radarr


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: integration


> Integration with Radarr

> Movie management automation and metadata synchronization
**API Base URL**: `http://localhost:7878/api/v3`
**Authentication**: api_key

---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | - |
| Sources | ✅ | - |
| Instructions | 🟡 | - |
| Code | 🔴 | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: ✅ Complete



---


## Architecture

### Integration Structure

```
internal/integration/radarr/
├── client.go              # API client
├── types.go               # Response types
├── mapper.go              # Map external → internal types
├── cache.go               # Response caching
└── client_test.go         # Tests
```

### Data Flow

<!-- Data flow diagram -->

### Provides
<!-- Data provided by integration -->


## Implementation

### File Structure

<!-- File structure -->

### Key Interfaces

```go
// Radarr integration service
type RadarrService interface {
  // Movie management
  AddMovie(ctx context.Context, tmdbID int, qualityProfileID int, rootFolder string) (*RadarrMovie, error)
  DeleteMovie(ctx context.Context, radarrID int, deleteFiles bool) error
  SearchMovie(ctx context.Context, radarrID int) error  // Trigger download

  // Sync
  SyncLibrary(ctx context.Context, instanceID uuid.UUID) error
  GetMovieStatus(ctx context.Context, movieID uuid.UUID) (*MovieStatus, error)

  // Calendar
  GetUpcoming(ctx context.Context, start, end time.Time) ([]CalendarEntry, error)
}

// Radarr movie structure
type RadarrMovie struct {
  ID              int      `json:"id"`
  Title           string   `json:"title"`
  Year            int      `json:"year"`
  TMDbID          int      `json:"tmdbId"`
  IMDbID          string   `json:"imdbId"`
  Monitored       bool     `json:"monitored"`
  QualityProfile  int      `json:"qualityProfileId"`
  RootFolderPath  string   `json:"rootFolderPath"`
  Path            string   `json:"path"`
  HasFile         bool     `json:"hasFile"`
  SizeOnDisk      int64    `json:"sizeOnDisk"`
}
```


### Dependencies

**Go Packages**:
- `net/http` - HTTP client
- `github.com/google/uuid` - UUID support
- `github.com/jackc/pgx/v5` - PostgreSQL driver
- `github.com/riverqueue/river` - Background sync jobs
- `go.uber.org/fx` - Dependency injection

**External Services**:
- Radarr v3+ (self-hosted)






## Configuration
### Environment Variables

```bash
# Radarr instance
RADARR_URL=http://localhost:7878
RADARR_API_KEY=your_api_key_here

# Sync settings
RADARR_AUTO_SYNC=true
RADARR_SYNC_INTERVAL=300  # 5 minutes
```


### Config Keys

```yaml
integrations:
  radarr:
    instances:
      - name: Main Radarr
        base_url: http://localhost:7878
        api_key: ${RADARR_API_KEY}
        enabled: true
        auto_sync: true
        sync_interval: 300
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
- [01_ARCHITECTURE](../../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../../architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../../architecture/03_METADATA_SYSTEM.md)

### External Sources
- [Uber fx](../../../sources/tooling/fx.md) - Auto-resolved from fx
- [Go context](../../../sources/go/stdlib/context.md) - Auto-resolved from go-context
- [pgx PostgreSQL Driver](../../../sources/database/pgx.md) - Auto-resolved from pgx
- [PostgreSQL Arrays](../../../sources/database/postgresql-arrays.md) - Auto-resolved from postgresql-arrays
- [PostgreSQL JSON Functions](../../../sources/database/postgresql-json.md) - Auto-resolved from postgresql-json
- [Radarr API Docs](../../../sources/apis/radarr-docs.md) - Auto-resolved from radarr-docs
- [River Job Queue](../../../sources/tooling/river.md) - Auto-resolved from river
- [Servarr Wiki](../../../sources/apis/servarr-wiki.md) - Auto-resolved from servarr-wiki
- [Typesense API](../../../sources/infrastructure/typesense.md) - Auto-resolved from typesense
- [Typesense Go Client](../../../sources/infrastructure/typesense-go.md) - Auto-resolved from typesense-go

