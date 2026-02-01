

---
sources:
  - name: Last.fm API
    url: ../../../sources/apis/lastfm.md
    note: Auto-resolved from lastfm-api
  - name: River Job Queue
    url: ../../../sources/tooling/river.md
    note: Auto-resolved from river
design_refs:
  - title: 01_ARCHITECTURE
    path: ../../architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: ../../architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: ../../architecture/03_METADATA_SYSTEM.md
---

## Table of Contents

- [ListenBrainz](#listenbrainz)
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
- [No global config needed - users provide their own tokens](#no-global-config-needed-users-provide-their-own-tokens)
    - [Config Keys](#config-keys)
  - [Testing Strategy](#testing-strategy)
    - [Unit Tests](#unit-tests)
    - [Integration Tests](#integration-tests)
    - [Test Coverage](#test-coverage)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)


# ListenBrainz


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: integration


> Integration with ListenBrainz

> Open-source music listening history tracking (MusicBrainz project)
**API Base URL**: `https://api.listenbrainz.org/1`
**Authentication**: token

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
internal/integration/listenbrainz/
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
// ListenBrainz integration service
type ListenBrainzService interface {
  // Connection
  ConnectUser(ctx context.Context, userID uuid.UUID, userToken string) (*LBConnection, error)
  ValidateToken(ctx context.Context, userToken string) (string, error)  // Returns username

  // Submission
  SubmitListen(ctx context.Context, userID uuid.UUID, trackID uuid.UUID, listenedAt time.Time) error
  SubmitPlayingNow(ctx context.Context, userID uuid.UUID, trackID uuid.UUID) error

  // Import
  ImportHistory(ctx context.Context, connectionID uuid.UUID, minTimestamp *time.Time) error
  GetUserStats(ctx context.Context, connectionID uuid.UUID) (*LBStats, error)
}

// Listen payload for submission
type ListenPayload struct {
  ListenType string   `json:"listen_type"`  // "single", "playing_now", "import"
  Payload    []Listen `json:"payload"`
}

type Listen struct {
  ListenedAt    int64         `json:"listened_at,omitempty"`
  TrackMetadata TrackMetadata `json:"track_metadata"`
}

type TrackMetadata struct {
  ArtistName      string                 `json:"artist_name"`
  TrackName       string                 `json:"track_name"`
  ReleaseName     string                 `json:"release_name,omitempty"`
  AdditionalInfo  map[string]interface{} `json:"additional_info,omitempty"`
}
```


### Dependencies

**Go Packages**:
- `net/http` - HTTP client
- `github.com/google/uuid` - UUID support (for MBIDs)
- `github.com/jackc/pgx/v5` - PostgreSQL driver
- `github.com/riverqueue/river` - Background jobs
- `go.uber.org/fx` - Dependency injection

**External Services**:
- ListenBrainz account (free, open-source)






## Configuration
### Environment Variables

```bash
# No global config needed - users provide their own tokens
LISTENBRAINZ_AUTO_SUBMIT=true
LISTENBRAINZ_IMPORT_ON_CONNECT=true
```


### Config Keys

```yaml
integrations:
  listenbrainz:
    auto_submit: true              # Auto-submit listens as user plays tracks
    import_on_connect: true        # Import existing history when user connects
    max_import_age_days: 365       # Only import listens from last year
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
- [Last.fm API](../../../sources/apis/lastfm.md) - Auto-resolved from lastfm-api
- [River Job Queue](../../../sources/tooling/river.md) - Auto-resolved from river

