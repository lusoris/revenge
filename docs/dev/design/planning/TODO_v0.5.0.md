# TODO v0.5.0 - Audio

<!-- DESIGN: planning, README, test_output_claude, test_output_wiki -->


<!-- TOC-START -->

## Table of Contents

- [Overview](#overview)
- [Deliverables](#deliverables)
  - [Music Module (Backend)](#music-module-backend)
  - [MusicBrainz Integration](#musicbrainz-integration)
  - [Last.fm Integration](#lastfm-integration)
  - [Lidarr Integration](#lidarr-integration)
  - [Audio Player](#audio-player)
  - [Search Integration](#search-integration)
  - [Frontend Updates](#frontend-updates)
- [Verification Checklist](#verification-checklist)
- [Dependencies](#dependencies)
- [Related Documentation](#related-documentation)

<!-- TOC-END -->


> Music Module

**Status**: 🔴 Not Started
**Tag**: `v0.5.0`
**Focus**: Music Module (Artists, Albums, Tracks)

**Depends On**: [v0.4.0](TODO_v0.4.0.md) (Pattern established for content modules)

---

## Overview

This milestone adds music library support with artist/album/track hierarchy, MusicBrainz and Last.fm metadata, Lidarr integration, lyrics support, and a dedicated audio player.

---

## Deliverables

### Music Module (Backend)

- [ ] **Database Schema** (`migrations/`)
  - [ ] `public.music_artists` table
    - [ ] id, name, sort_name
    - [ ] overview, biography
    - [ ] type (person, group, orchestra, etc.)
    - [ ] begin_date, end_date
    - [ ] country, area
    - [ ] image_path
    - [ ] musicbrainz_id, lastfm_url, spotify_id
  - [ ] `public.music_albums` table
    - [ ] id, title
    - [ ] artist_id (primary artist)
    - [ ] release_date, year
    - [ ] type (album, EP, single, compilation)
    - [ ] total_tracks, total_discs
    - [ ] duration_ms
    - [ ] cover_path
    - [ ] musicbrainz_id, spotify_id
  - [ ] `public.music_album_artists` table (many-to-many)
  - [ ] `public.music_tracks` table
    - [ ] id, title
    - [ ] album_id, artist_id
    - [ ] track_number, disc_number
    - [ ] duration_ms
    - [ ] explicit
    - [ ] musicbrainz_id, isrc
  - [ ] `public.music_track_artists` table (featuring artists)
  - [ ] `public.music_files` table
    - [ ] id, track_id
    - [ ] path, format, bitrate, sample_rate
    - [ ] codec, channels
    - [ ] size_bytes
  - [ ] `public.music_genres` table
  - [ ] `public.music_album_genres` table
  - [ ] `public.music_playback_history` table
  - [ ] `public.music_lyrics` table
    - [ ] track_id
    - [ ] plain_lyrics, synced_lyrics (LRC format)
    - [ ] source
  - [ ] Indexes on musicbrainz_id, spotify_id, name

- [ ] **Entity** (`internal/content/music/entity.go`)
  - [ ] Artist struct
  - [ ] Album struct
  - [ ] Track struct
  - [ ] TrackFile struct
  - [ ] Lyrics struct
  - [ ] PlaybackHistory struct

- [ ] **Repository** (`internal/content/music/`)
  - [ ] `repository.go` - Interface
  - [ ] `repository_pg.go` - Implementation
  - [ ] Artist CRUD
  - [ ] Album CRUD
  - [ ] Track CRUD
  - [ ] Playback history operations
  - [ ] List/filter operations

- [ ] **Service** (`internal/content/music/service.go`)
  - [ ] Get artist by ID
  - [ ] Get artist albums
  - [ ] Get artist top tracks
  - [ ] Get album by ID
  - [ ] Get album tracks
  - [ ] Get track by ID
  - [ ] List artists (paginated)
  - [ ] List albums (paginated)
  - [ ] Search music
  - [ ] Record playback
  - [ ] Get recently played
  - [ ] Get most played
  - [ ] Get lyrics for track

- [ ] **Library Provider** (`internal/content/music/library_service.go`)
  - [ ] Implement LibraryProvider interface
  - [ ] Scan library path
  - [ ] Parse folder structure (Artist/Album/Track)
  - [ ] Read audio file metadata (go-taglib)
  - [ ] Match files to tracks
  - [ ] Handle file changes

- [ ] **Handler** (`internal/api/music_handler.go`)
  - [ ] `GET /api/v1/music/artists` (list)
  - [ ] `GET /api/v1/music/artists/:id`
  - [ ] `GET /api/v1/music/artists/:id/albums`
  - [ ] `GET /api/v1/music/artists/:id/top-tracks`
  - [ ] `GET /api/v1/music/artists/:id/similar`
  - [ ] `GET /api/v1/music/albums` (list)
  - [ ] `GET /api/v1/music/albums/:id`
  - [ ] `GET /api/v1/music/albums/:id/tracks`
  - [ ] `GET /api/v1/music/tracks/:id`
  - [ ] `GET /api/v1/music/tracks/:id/lyrics`
  - [ ] `GET /api/v1/music/tracks/:id/stream` (audio stream)
  - [ ] `POST /api/v1/music/tracks/:id/play` (record playback)
  - [ ] `GET /api/v1/music/recently-played`
  - [ ] `GET /api/v1/music/most-played`
  - [ ] `POST /api/v1/music/artists/:id/refresh`
  - [ ] `POST /api/v1/music/albums/:id/refresh`

- [ ] **River Jobs** (`internal/content/music/jobs.go`)
  - [ ] MusicLibraryScanJob
  - [ ] ArtistMetadataRefreshJob
  - [ ] AlbumMetadataRefreshJob
  - [ ] LyricsFetchJob

- [ ] **fx Module** (`internal/content/music/module.go`)

- [ ] **Tests**
  - [ ] Unit tests (80%+ coverage)
  - [ ] Integration tests

### MusicBrainz Integration

- [ ] **MusicBrainz Client** (`internal/service/metadata/musicbrainz/client.go`)
  - [ ] API implementation
  - [ ] Rate limiting (1 req/s strict)
  - [ ] User-Agent header (required)
  - [ ] Response caching

- [ ] **MusicBrainz Service** (`internal/service/metadata/musicbrainz/service.go`)
  - [ ] Search artist
  - [ ] Get artist details
  - [ ] Search release group (album)
  - [ ] Get release group details
  - [ ] Get release details
  - [ ] Search recording (track)
  - [ ] Get cover art (via Cover Art Archive)
  - [ ] Get artist relations

- [ ] **Tests**
  - [ ] Unit tests with mock API

### Last.fm Integration

- [ ] **Last.fm Client** (`internal/service/metadata/lastfm/client.go`)
  - [ ] API implementation
  - [ ] API key authentication
  - [ ] Rate limiting (5 req/s)

- [ ] **Last.fm Service** (`internal/service/metadata/lastfm/service.go`)
  - [ ] Get artist info
  - [ ] Get artist top albums
  - [ ] Get artist top tracks
  - [ ] Get artist similar
  - [ ] Get album info
  - [ ] Get track info
  - [ ] Get artist tags
  - [ ] Search artist
  - [ ] Search album
  - [ ] Search track

- [ ] **Metadata Enrichment**
  - [ ] Artist biography
  - [ ] Album description
  - [ ] Genre/tag extraction
  - [ ] Listener/play counts

- [ ] **Tests**
  - [ ] Unit tests with mock API

### Lidarr Integration

- [ ] **Lidarr Client** (`internal/service/metadata/lidarr/client.go`)
  - [ ] API v1 implementation
  - [ ] Authentication (API key)
  - [ ] Error handling

- [ ] **Lidarr Service** (`internal/service/metadata/lidarr/service.go`)
  - [ ] Get all artists
  - [ ] Get artist by ID
  - [ ] Get albums for artist
  - [ ] Get tracks for album
  - [ ] Sync library (Lidarr → Revenge)
  - [ ] Trigger artist refresh
  - [ ] Get quality profiles
  - [ ] Get metadata profiles
  - [ ] Get root folders

- [ ] **Sync Logic** (`internal/service/metadata/lidarr/sync.go`)
  - [ ] Full sync (initial)
  - [ ] Incremental sync
  - [ ] File path mapping
  - [ ] Album/track matching

- [ ] **Webhook Handler**
  - [ ] `POST /api/v1/webhooks/lidarr`
  - [ ] Handle: Grab, Download, Rename, Delete, Retag events

- [ ] **Handler** (`internal/api/lidarr_handler.go`)
  - [ ] `GET /api/v1/admin/integrations/lidarr/status`
  - [ ] `POST /api/v1/admin/integrations/lidarr/sync`

- [ ] **River Jobs**
  - [ ] LidarrSyncJob
  - [ ] LidarrWebhookJob

- [ ] **Tests**
  - [ ] Unit tests with mock API

### Audio Player

- [ ] **Audio Streaming Endpoint** (`internal/api/stream_handler.go`)
  - [ ] `GET /api/v1/music/tracks/:id/stream`
  - [ ] Range request support
  - [ ] Format detection
  - [ ] Optional transcoding endpoint

- [ ] **Lyrics Support** (`internal/service/lyrics/`)
  - [ ] LRC parser (synced lyrics)
  - [ ] Plain lyrics fallback
  - [ ] Lyrics search providers (LRCLIB, etc.)
  - [ ] Cache fetched lyrics

- [ ] **Playback Recording**
  - [ ] Track play start
  - [ ] Track play complete (scrobble threshold: 50% or 4 min)
  - [ ] Listening history

### Search Integration

- [ ] **Artist Collection Schema**
  ```json
  {
    "name": "music_artists",
    "fields": [
      {"name": "id", "type": "string"},
      {"name": "name", "type": "string"},
      {"name": "overview", "type": "string"},
      {"name": "type", "type": "string"},
      {"name": "genres", "type": "string[]"}
    ]
  }
  ```

- [ ] **Album Collection Schema**
  ```json
  {
    "name": "music_albums",
    "fields": [
      {"name": "id", "type": "string"},
      {"name": "title", "type": "string"},
      {"name": "artist_name", "type": "string"},
      {"name": "year", "type": "int32"},
      {"name": "type", "type": "string"},
      {"name": "genres", "type": "string[]"}
    ]
  }
  ```

- [ ] **Track Collection Schema**
  ```json
  {
    "name": "music_tracks",
    "fields": [
      {"name": "id", "type": "string"},
      {"name": "title", "type": "string"},
      {"name": "artist_name", "type": "string"},
      {"name": "album_title", "type": "string"}
    ]
  }
  ```

- [ ] **Search Service Updates**
  - [ ] Index artist, album, track
  - [ ] Multi-type music search

### Frontend Updates

- [ ] **Music Landing** (`/music`)
  - [ ] Recently played
  - [ ] Recently added albums
  - [ ] Featured artists

- [ ] **Artists View** (`/music/artists`)
  - [ ] Artist cards/list
  - [ ] Filtering (genre)
  - [ ] Sorting (name, recently added)

- [ ] **Artist Detail** (`/music/artists/[id]`)
  - [ ] Artist image
  - [ ] Biography
  - [ ] Discography (albums)
  - [ ] Top tracks
  - [ ] Similar artists

- [ ] **Albums View** (`/music/albums`)
  - [ ] Album grid with covers
  - [ ] Filtering (genre, year, type)
  - [ ] Sorting options

- [ ] **Album Detail** (`/music/albums/[id]`)
  - [ ] Album cover (large)
  - [ ] Album info
  - [ ] Track list with durations
  - [ ] Play all / Shuffle buttons

- [ ] **Audio Player Component**
  - [ ] Fixed bottom player bar
  - [ ] Now playing info
  - [ ] Play/Pause, Next/Previous
  - [ ] Progress bar (seekable)
  - [ ] Volume control
  - [ ] Queue management
  - [ ] Shuffle / Repeat modes
  - [ ] Full-screen player view

- [ ] **Lyrics Display**
  - [ ] Synced lyrics (karaoke mode)
  - [ ] Plain lyrics fallback
  - [ ] Toggle lyrics view

- [ ] **Queue Management**
  - [ ] View queue
  - [ ] Reorder queue
  - [ ] Clear queue
  - [ ] Add to queue

- [ ] **Search Updates**
  - [ ] Include artists, albums, tracks
  - [ ] Type filter

- [ ] **Admin: Lidarr Integration**
  - [ ] Settings page
  - [ ] Connection test
  - [ ] Manual sync

---

## Verification Checklist

- [ ] Music library scans and imports
- [ ] Artists, albums, tracks display correctly
- [ ] Lidarr sync imports music
- [ ] Audio player plays music
- [ ] Lyrics display (synced when available)
- [ ] Queue management works
- [ ] Playback history records
- [ ] Search includes music
- [ ] All tests pass (80%+ coverage)
- [ ] CI pipeline passes

---

## Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| github.com/wtolson/go-taglib | latest | Audio metadata (CGo) |
| github.com/go-resty/resty/v2 | v2.17.1 | HTTP client |

---

## Design Documentation

> **Note**: Design work for v0.5.0 scope is **COMPLETE**. The following design documents exist and should be referenced during implementation:

### Core Module Designs
- [MUSIC_MODULE.md](../features/music/MUSIC_MODULE.md) - Complete music module design (artists, albums, tracks, playlists)
- [AUDIO_STREAMING.md](../technical/AUDIO_STREAMING.md) - Audio streaming and progress tracking architecture

### Integration Designs
- [MUSICBRAINZ.md](../integrations/metadata/music/MUSICBRAINZ.md) - MusicBrainz metadata provider
- [LASTFM.md](../integrations/metadata/music/LASTFM.md) - Last.fm integration (scrobbling, metadata)
- [LASTFM_SCROBBLE.md](../integrations/scrobbling/LASTFM_SCROBBLE.md) - Last.fm scrobbling service
- [LISTENBRAINZ.md](../integrations/scrobbling/LISTENBRAINZ.md) - ListenBrainz scrobbling (designed, bonus feature)
- [LIDARR.md](../integrations/servarr/LIDARR.md) - Lidarr integration design
- [SPOTIFY.md](../integrations/metadata/music/SPOTIFY.md) - Spotify metadata enrichment (designed, bonus)
- [DISCOGS.md](../integrations/metadata/music/DISCOGS.md) - Discogs metadata (designed, bonus)

### Related Features (Designed, Future Enhancement)
- [SCROBBLING.md](../features/shared/SCROBBLING.md) - External scrobbling architecture (Trakt, Last.fm, etc.)

---

## Related Documentation

- [ROADMAP.md](ROADMAP.md) - Full roadmap overview
- [DESIGN_INDEX.md](../DESIGN_INDEX.md) - Full design documentation index
