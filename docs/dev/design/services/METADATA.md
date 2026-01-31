# Metadata Service

<!-- SOURCES: fx, lastfm-api, river -->

<!-- DESIGN: services, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES, 03_METADATA_SYSTEM -->


> External metadata providers for media enrichment


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Developer Resources](#developer-resources)
- [Overview](#overview)
- [Provider Inventory](#provider-inventory)
  - [Video Providers](#video-providers)
  - [Music Providers](#music-providers)
  - [Book Providers](#book-providers)
  - [Comics Providers](#comics-providers)
  - [QAR Providers (Adult)](#qar-providers-adult)
  - [Arr Services (Primary)](#arr-services-primary)
- [Provider Architecture](#provider-architecture)
- [Provider Interface](#provider-interface)
- [TMDb Provider](#tmdb-provider)
  - [Rate Limits](#rate-limits)
- [Radarr Provider](#radarr-provider)
  - [Priority](#priority)
- [StashDB Provider (QAR)](#stashdb-provider-qar)
  - [GraphQL Client](#graphql-client)
- [Configuration](#configuration)
- [Implementation Checklist](#implementation-checklist)
  - [Phase 1: Core Infrastructure](#phase-1-core-infrastructure)
  - [Phase 2: Video Providers](#phase-2-video-providers)
  - [Phase 3: Audio Providers](#phase-3-audio-providers)
  - [Phase 4: QAR Providers](#phase-4-qar-providers)
  - [Phase 5: Book Providers](#phase-5-book-providers)
  - [Phase 6: Service Layer](#phase-6-service-layer)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related](#related)

<!-- TOC-END -->

**Module**: `internal/service/metadata/`
**Dependencies**: [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#metadata-providers)

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


## Related

- [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#metadata-providers) - Provider inventory and status
- [03_METADATA_SYSTEM.md](../architecture/03_METADATA_SYSTEM.md) - Architecture
- [integrations/metadata/](../integrations/metadata/) - Per-provider API details
- [integrations/servarr/](../integrations/servarr/) - Arr service integrations
