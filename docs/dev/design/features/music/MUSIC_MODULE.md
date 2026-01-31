# Music Module

<!-- SOURCES: fx, lastfm-api, ogen, pgx, postgresql-arrays, postgresql-json, river, sqlc, sqlc-config -->

<!-- DESIGN: features/music, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES, 03_METADATA_SYSTEM -->


> Music content management with metadata enrichment from MusicBrainz and Last.fm


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Developer Resources](#developer-resources)
- [Overview](#overview)
- [Architecture](#architecture)
- [Files (Planned)](#files-planned)
- [Entities (Planned)](#entities-planned)
  - [Artist](#artist)
  - [Album](#album)
  - [Track](#track)
- [Metadata Priority Chain](#metadata-priority-chain)
- [Arr Integration](#arr-integration)
- [Scrobbling](#scrobbling)
- [Database Schema (Planned)](#database-schema-planned)
- [API Endpoints (Planned)](#api-endpoints-planned)
- [Implementation Checklist](#implementation-checklist)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Documents](#related-documents)

<!-- TOC-END -->

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | 🟡 | Scaffold - needs detailed spec |
| Sources | 🔴 | MusicBrainz, Last.fm, Discogs API docs needed |
| Instructions | 🔴 |  |
| Code | 🔴 |  |
| Linting | 🔴 |  |
| Unit Testing | 🔴 |  |
| Integration Testing | 🔴 |  |**Location**: `internal/content/music/`

---

## Developer Resources

| Source | URL | Purpose |
|--------|-----|---------|
| MusicBrainz API | [musicbrainz.org/doc/Development](https://musicbrainz.org/doc/Development/XML_Web_Service/Version_2) | Primary music metadata |
| Last.fm API | [last.fm/api](https://www.last.fm/api) | Scrobbling, listening stats |
| Lidarr API | See [integrations/servarr/LIDARR.md](../../integrations/servarr/LIDARR.md) | Servarr integration |

---

## Overview

The Music module provides complete music library management:

- Entity definitions (Artist, Album, Track, etc.)
- Repository pattern with PostgreSQL implementation
- Service layer with otter caching
- Background jobs for metadata enrichment via River
- User data (ratings, play counts, favorites)
- Scrobbling integration (Last.fm, ListenBrainz)

---

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                       API Layer                              │
│                    (ogen handlers)                           │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                    Music Service                             │
│   - Local cache (otter)                                      │
│   - Business logic                                           │
│   - Scrobbling integration                                   │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                    Repository Layer                          │
│   - PostgreSQL queries (sqlc)                                │
│   - User data (ratings, play history)                        │
│   - Relations (artists, albums, genres)                      │
└─────────────────────────────────────────────────────────────┘
```

---

## Files (Planned)

| File | Description |
|------|-------------|
| `entity.go` | Domain entities (Artist, Album, Track, etc.) |
| `repository.go` | Repository interface definition |
| `repository_pg.go` | PostgreSQL implementation |
| `service.go` | Business logic with caching |
| `jobs.go` | River background jobs |
| `metadata_provider.go` | MusicBrainz/Last.fm interface |
| `module.go` | fx dependency injection |

---

## Entities (Planned)

### Artist

```go
type Artist struct {
    shared.ContentEntity

    Name          string
    SortName      string
    MusicBrainzID *uuid.UUID
    DiscogsID     *int
    Biography     string
    Country       string
    StartDate     *time.Time
    EndDate       *time.Time
    Type          string // person, group, orchestra, choir
}
```

### Album

```go
type Album struct {
    shared.ContentEntity

    Title         string
    ArtistID      uuid.UUID
    MusicBrainzID *uuid.UUID
    ReleaseDate   *time.Time
    ReleaseType   string // album, single, EP, compilation
    TotalTracks   int
    TotalDiscs    int
    Label         string
}
```

### Track

```go
type Track struct {
    shared.ContentEntity

    Title         string
    AlbumID       uuid.UUID
    ArtistID      uuid.UUID
    MusicBrainzID *uuid.UUID
    DiscNumber    int
    TrackNumber   int
    DurationMs    int64
    FilePath      string
    Container     string
    Codec         string
    Bitrate       int
}
```

---

## Metadata Priority Chain

See [00_SOURCE_OF_TRUTH.md](../../00_SOURCE_OF_TRUTH.md) for the core metadata priority principle.

```
1. LOCAL CACHE     → First, instant UI display
2. LIDARR          → Arr-first metadata
3. MUSICBRAINZ     → Primary external source
4. LAST.FM         → Listening stats, similar artists
5. DISCOGS         → Fallback for obscure releases
```

---

## Arr Integration

**Primary**: Lidarr

See [integrations/servarr/LIDARR.md](../../integrations/servarr/LIDARR.md) for:
- Webhook handling
- Import notifications
- Library sync patterns

---

## Scrobbling

Music module integrates with scrobbling services:

- **Last.fm** - Primary scrobbling target
- **ListenBrainz** - Open alternative

See [features/shared/SCROBBLING.md](../shared/SCROBBLING.md) for implementation details.

---

## Database Schema (Planned)

Tables in `public` schema:

- `artists` - Artist entities
- `albums` - Album entities
- `tracks` - Track entities
- `artist_album` - Artist-album relationships
- `track_artist` - Track-artist relationships (feat. artists)
- `music_genres` - Genre mappings
- `user_track_history` - Play history for scrobbling
- `user_artist_favorites` - Favorite artists
- `user_album_favorites` - Favorite albums

---

## API Endpoints (Planned)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/music/artists` | List artists |
| GET | `/api/v1/music/artists/{id}` | Get artist details |
| GET | `/api/v1/music/albums` | List albums |
| GET | `/api/v1/music/albums/{id}` | Get album details |
| GET | `/api/v1/music/tracks` | List tracks |
| GET | `/api/v1/music/tracks/{id}` | Get track details |
| POST | `/api/v1/music/tracks/{id}/scrobble` | Record play |
| GET | `/api/v1/music/now-playing` | Current playback |

---

## Implementation Checklist

- [ ] Define entity structs in `entity.go`
- [ ] Create repository interface
- [ ] Implement PostgreSQL repository
- [ ] Create database migrations
- [ ] Implement service layer with caching
- [ ] Add River jobs for metadata enrichment
- [ ] Integrate MusicBrainz provider
- [ ] Integrate Last.fm provider
- [ ] Add Lidarr webhook handlers
- [ ] Implement scrobbling
- [ ] Write unit tests
- [ ] Write integration tests

---


## Related Documents

- [MusicBrainz Integration](../../integrations/metadata/music/MUSICBRAINZ.md)
- [Last.fm Integration](../../integrations/metadata/music/LASTFM.md)
- [Discogs Integration](../../integrations/metadata/music/DISCOGS.md)
- [Lidarr Integration](../../integrations/servarr/LIDARR.md)
- [Scrobbling](../shared/SCROBBLING.md)
