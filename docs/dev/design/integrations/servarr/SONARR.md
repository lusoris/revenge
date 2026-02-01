## Table of Contents

- [Sonarr](#sonarr)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Integration Structure](#integration-structure)
    - [Data Flow](#data-flow)
    - [Provides](#provides)
  - [Implementation](#implementation)
    - [Key Interfaces](#key-interfaces)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# Sonarr


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: integration


> Integration with Sonarr

> TV show management automation
**API Base URL**: `http://localhost:8989/api/v3`
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

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│  Revenge    │────▶│   Sonarr     │────▶│   Sonarr    │
│  Request    │◀────│ Integration  │◀────│   Server    │
│  System     │     │              │     └─────────────┘
└─────────────┘     └──────┬───────┘
                           │
                    ┌──────▼────────┐
                    │   Webhook     │
                    │   Handler     │
                    └───────────────┘
```


### Integration Structure

```
internal/integration/sonarr/
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

### Key Interfaces

```go
type SonarrService interface {
  // Series management
  AddSeries(ctx context.Context, tvdbID int, qualityProfileID int, rootFolder string) (*SonarrSeries, error)
  DeleteSeries(ctx context.Context, sonarrID int, deleteFiles bool) error
  SearchSeason(ctx context.Context, seriesID int, seasonNumber int) error

  // Episode management
  GetEpisodes(ctx context.Context, seriesID int) ([]SonarrEpisode, error)
  GetCalendar(ctx context.Context, start, end time.Time) ([]CalendarEpisode, error)

  // Sync
  SyncLibrary(ctx context.Context, instanceID uuid.UUID) error
}

type SonarrSeries struct {
  ID              int      `json:"id"`
  Title           string   `json:"title"`
  Year            int      `json:"year"`
  TVDbID          int      `json:"tvdbId"`
  TVRageID        int      `json:"tvRageId"`
  IMDbID          string   `json:"imdbId"`
  Monitored       bool     `json:"monitored"`
  QualityProfile  int      `json:"qualityProfileId"`
  SeasonFolder    bool     `json:"seasonFolder"`
  Path            string   `json:"path"`
  Seasons         []Season `json:"seasons"`
}

type SonarrEpisode struct {
  ID              int      `json:"id"`
  SeriesID        int      `json:"seriesId"`
  SeasonNumber    int      `json:"seasonNumber"`
  EpisodeNumber   int      `json:"episodeNumber"`
  Title           string   `json:"title"`
  AirDate         string   `json:"airDate"`
  HasFile         bool     `json:"hasFile"`
  Monitored       bool     `json:"monitored"`
}
```
















## Related Documentation
### Design Documents
- [01_ARCHITECTURE](../../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../../architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../../architecture/03_METADATA_SYSTEM.md)

### External Sources
- [Uber fx](../../../sources/tooling/fx.md) - Auto-resolved from fx
- [pgx PostgreSQL Driver](../../../sources/database/pgx.md) - Auto-resolved from pgx
- [PostgreSQL Arrays](../../../sources/database/postgresql-arrays.md) - Auto-resolved from postgresql-arrays
- [PostgreSQL JSON Functions](../../../sources/database/postgresql-json.md) - Auto-resolved from postgresql-json
- [River Job Queue](../../../sources/tooling/river.md) - Auto-resolved from river
- [Servarr Wiki](../../../sources/apis/servarr-wiki.md) - Auto-resolved from servarr-wiki
- [Sonarr API Docs](../../../sources/apis/sonarr-docs.md) - Auto-resolved from sonarr-docs
- [Typesense API](../../../sources/infrastructure/typesense.md) - Auto-resolved from typesense
- [Typesense Go Client](../../../sources/infrastructure/typesense-go.md) - Auto-resolved from typesense-go

