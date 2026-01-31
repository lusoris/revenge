# Metadata Service

> External metadata providers for media enrichment

**Module**: `internal/service/metadata/`
**Dependencies**: [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#metadata-providers)

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

## Developer Resources

> Package versions: [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#go-dependencies-core)
> External APIs: [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#metadata-providers)

| Package | Purpose |
|---------|---------|
| genqlient | Type-safe GraphQL client |
| resty | HTTP client |

---

## Overview

The Metadata service provides unified access to external APIs for fetching media metadata. It follows the **Servarr-First** principle: when connected to an *arr service, use its cached metadata as primary source.

> See [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#metadata-priority-chain) for the priority chain.

---

## Provider Inventory

> See [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#metadata-providers) for rate limits and implementation status.

### Video Providers

| Provider | Content | Location | Status |
|----------|---------|----------|--------|
| TMDb | Movies, TV | `internal/service/metadata/tmdb` | ✅ |
| TheTVDB | TV | `internal/service/metadata/thetvdb` | 🔴 |

### Music Providers

| Provider | Content | Location | Status |
|----------|---------|----------|--------|
| MusicBrainz | Music | `internal/service/metadata/musicbrainz` | 🔴 |
| Last.fm | Music | `internal/service/metadata/lastfm` | 🔴 |

### Book Providers

| Provider | Content | Location | Status |
|----------|---------|----------|--------|
| Audnexus | Audiobooks | `internal/service/metadata/audnexus` | 🔴 |
| OpenLibrary | Books | `internal/service/metadata/openlibrary` | 🔴 |

### Comics Providers

| Provider | Content | Location | Status |
|----------|---------|----------|--------|
| ComicVine | Comics | `internal/service/metadata/comicvine` | 🔴 |

### QAR Providers (Adult)

| Provider | Content | Location | Status |
|----------|---------|----------|--------|
| StashDB | Voyages, Expeditions | `internal/service/metadata/stashdb` | 🟡 |
| ThePornDB | Voyages, Expeditions | `internal/service/metadata/tpdb` | 🔴 |

### Arr Services (Primary)

| Provider | Content | Location | Status |
|----------|---------|----------|--------|
| Radarr | Movies | `internal/service/metadata/radarr` | ✅ |
| Sonarr | TV | `internal/service/metadata/sonarr` | 🔴 |
| Lidarr | Music | `internal/service/metadata/lidarr` | 🔴 |
| Whisparr | QAR | `internal/service/metadata/whisparr` | 🔴 |

---

## Provider Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Metadata Service                          │
└──────────────────────────┬──────────────────────────────────┘
                           │
         ┌─────────────────┼─────────────────┐
         │                 │                 │
         ▼                 ▼                 ▼
    ┌─────────┐      ┌──────────┐     ┌──────────┐
    │   Arr   │      │ External │     │ Enrichment│
    │Services │      │   APIs   │     │   APIs   │
    │Priority │      │ Priority │     │ Priority │
    │   1-2   │      │   3-4    │     │    5     │
    └─────────┘      └──────────┘     └──────────┘
         │                 │                 │
    ┌────┴────┐     ┌──────┴──────┐    ┌────┴────┐
    │ Radarr  │     │    TMDb     │    │  Fanart │
    │ Sonarr  │     │   TheTVDB   │    │  OMDB   │
    │ Lidarr  │     │ MusicBrainz │    │         │
    │ Whisparr│     │   StashDB   │    │         │
    └─────────┘     └─────────────┘    └─────────┘
```

---

## Provider Interface

All providers implement a common interface:

```go
type Provider interface {
    Name() string
    Priority() int
    ContentTypes() []string
    IsAvailable(ctx context.Context) bool
    Search(ctx context.Context, query SearchQuery) ([]SearchResult, error)
    GetMetadata(ctx context.Context, id ExternalID) (*Metadata, error)
}

type SearchQuery struct {
    Title    string
    Year     int
    IMDbID   string
    TVDbID   int
    // Content-specific fields...
}

type ExternalID struct {
    Provider string // "tmdb", "tvdb", "musicbrainz", etc.
    ID       string
}
```

---

## TMDb Provider

**Location**: `internal/service/metadata/tmdb/`

Primary external source for movies and TV shows.

```go
func (p *Provider) Name() string        { return "tmdb" }
func (p *Provider) Priority() int       { return 4 }
func (p *Provider) ContentTypes() []string { return []string{"movie", "series"} }
```

### Rate Limits

- 50 requests/second (with API key)
- Uses sturdyc for request coalescing

---

## Radarr Provider

**Location**: `internal/service/metadata/radarr/`

Primary metadata source when Radarr is connected. Uses cached metadata from Radarr which itself aggregates from TMDb and other sources.

```go
func (p *Provider) Name() string        { return "radarr" }
func (p *Provider) Priority() int       { return 2 }
func (p *Provider) ContentTypes() []string { return []string{"movie"} }
```

### Priority

Radarr takes priority over TMDb when connected (Servarr-first principle).

---

## StashDB Provider (QAR)

**Location**: `internal/service/metadata/stashdb/`

GraphQL-based metadata for QAR content (adult).

```go
func (p *Provider) Name() string        { return "stashdb" }
func (p *Provider) Priority() int       { return 4 }
func (p *Provider) ContentTypes() []string { return []string{"expedition", "voyage", "crew"} }
```

### GraphQL Client

Uses `Khan/genqlient` for type-safe GraphQL queries.

---

## Configuration

```yaml
metadata:
  # Video
  tmdb:
    api_key: "${REVENGE_TMDB_API_KEY}"
    enabled: true
  thetvdb:
    api_key: "${REVENGE_THETVDB_API_KEY}"
    enabled: true

  # Music
  musicbrainz:
    enabled: true  # No API key required
  lastfm:
    api_key: "${REVENGE_LASTFM_API_KEY}"
    enabled: true

  # Books
  audnexus:
    enabled: true  # No API key required
  openlibrary:
    enabled: true  # No API key required

  # Comics
  comicvine:
    api_key: "${REVENGE_COMICVINE_API_KEY}"
    enabled: true

  # QAR (Adult)
  stashdb:
    api_key: "${REVENGE_STASHDB_API_KEY}"
    enabled: true
  tpdb:
    api_key: "${REVENGE_TPDB_API_KEY}"
    enabled: true

# Arr Services
arr:
  radarr:
    url: "${REVENGE_RADARR_URL}"
    api_key: "${REVENGE_RADARR_API_KEY}"
    enabled: true
  sonarr:
    url: "${REVENGE_SONARR_URL}"
    api_key: "${REVENGE_SONARR_API_KEY}"
    enabled: true
```

---

## Implementation Checklist

### Phase 1: Core Infrastructure

- [ ] Create `internal/service/metadata/` package structure
- [ ] Define `provider.go` interface
- [ ] Implement provider registry with priority ordering
- [ ] Add fx module wiring in `module.go`

### Phase 2: Video Providers

- [ ] Implement TMDb provider
- [ ] Implement TheTVDB provider
- [ ] Add rate limiting with sturdyc
- [ ] Implement Radarr/Sonarr providers (arr-first)

### Phase 3: Audio Providers

- [ ] Implement MusicBrainz provider
- [ ] Implement Last.fm provider
- [ ] Implement Lidarr provider

### Phase 4: QAR Providers

- [ ] Implement StashDB GraphQL provider
- [ ] Implement ThePornDB provider
- [ ] Implement Whisparr provider

### Phase 5: Book Providers

- [ ] Implement Audnexus provider
- [ ] Implement OpenLibrary provider
- [ ] Implement Readarr provider

### Phase 6: Service Layer

- [ ] Implement unified metadata service
- [ ] Add priority-based provider selection
- [ ] Implement caching layer (otter)
- [ ] Add River jobs for background fetching

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
- [Library Service](LIBRARY.md)
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

## Related

- [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#metadata-providers) - Provider inventory and status
- [03_METADATA_SYSTEM.md](../architecture/03_METADATA_SYSTEM.md) - Architecture
- [integrations/metadata/](../integrations/metadata/) - Per-provider API details
- [integrations/servarr/](../integrations/servarr/) - Arr service integrations
