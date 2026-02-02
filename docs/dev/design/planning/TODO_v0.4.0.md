# TODO v0.4.0 - Shows

<!-- DESIGN: planning, README, test_output_claude, test_output_wiki -->


<!-- TOC-START -->

## Table of Contents

- [Overview](#overview)
- [Deliverables](#deliverables)
  - [TV Show Module (Backend)](#tv-show-module-backend)
  - [TheTVDB Integration](#thetvdb-integration)
  - [TMDb TV Support](#tmdb-tv-support)
  - [Sonarr Integration](#sonarr-integration)
  - [Episode Watch Progress](#episode-watch-progress)
  - [Search Integration](#search-integration)
  - [Frontend Updates](#frontend-updates)
- [Verification Checklist](#verification-checklist)
- [Dependencies from SOURCE_OF_TRUTH](#dependencies-from-source-of-truth)
- [Related Documentation](#related-documentation)

<!-- TOC-END -->


> TV Shows Module

**Status**: 🔴 Not Started
**Tag**: `v0.4.0`
**Focus**: TV Shows Module (Series, Seasons, Episodes)

**Depends On**: [v0.3.0](TODO_v0.3.0.md) (Movie module patterns to follow)

---

## Overview

This milestone adds TV show support with full series/season/episode hierarchy, TheTVDB integration, Sonarr sync, and enhanced watch progress tracking for sequential content.

---

## Deliverables

### TV Show Module (Backend)

- [ ] **Database Schema** (`migrations/`)
  - [ ] `public.tv_shows` table
    - [ ] id, title, original_title
    - [ ] first_air_date, last_air_date
    - [ ] status (continuing, ended, canceled)
    - [ ] overview, tagline
    - [ ] poster_path, backdrop_path
    - [ ] tmdb_id, tvdb_id, imdb_id
    - [ ] episode_run_time
    - [ ] total_seasons, total_episodes
  - [ ] `public.tv_seasons` table
    - [ ] id, tv_show_id
    - [ ] season_number
    - [ ] name, overview
    - [ ] air_date
    - [ ] poster_path
    - [ ] episode_count
  - [ ] `public.tv_episodes` table
    - [ ] id, season_id, tv_show_id
    - [ ] episode_number, season_number
    - [ ] name, overview
    - [ ] air_date, runtime_minutes
    - [ ] still_path
    - [ ] tmdb_id, tvdb_id
  - [ ] `public.tv_show_genres` table
  - [ ] `public.tv_show_cast` table
  - [ ] `public.tv_episode_cast` table (guest stars)
  - [ ] `public.tv_episode_files` table
  - [ ] `public.tv_episode_watch_progress` table
  - [ ] Indexes on tmdb_id, tvdb_id, imdb_id

- [ ] **Entity** (`internal/content/tvshow/entity.go`)
  - [ ] TVShow struct
  - [ ] Season struct
  - [ ] Episode struct
  - [ ] EpisodeFile struct
  - [ ] TVShowCast struct
  - [ ] EpisodeWatchProgress struct

- [ ] **Repository** (`internal/content/tvshow/`)
  - [ ] `repository.go` - Interface
  - [ ] `repository_pg.go` - Implementation
  - [ ] Show CRUD
  - [ ] Season CRUD
  - [ ] Episode CRUD
  - [ ] Watch progress operations
  - [ ] List with filters

- [ ] **Service** (`internal/content/tvshow/service.go`)
  - [ ] Get show by ID
  - [ ] Get seasons for show
  - [ ] Get episodes for season
  - [ ] Get episode by ID
  - [ ] List shows (paginated)
  - [ ] Search shows
  - [ ] Update episode watch progress
  - [ ] Get next episode to watch
  - [ ] Get continue watching (shows)
  - [ ] Get recently added episodes
  - [ ] Mark season as watched
  - [ ] Mark show as watched

- [ ] **Library Provider** (`internal/content/tvshow/library_service.go`)
  - [ ] Implement LibraryProvider interface
  - [ ] Scan library path
  - [ ] Parse folder structure (Show/Season/Episode)
  - [ ] Match files to episodes
  - [ ] Handle file changes

- [ ] **Handler** (`internal/api/tvshow_handler.go`)
  - [ ] `GET /api/v1/tv` (list shows)
  - [ ] `GET /api/v1/tv/:id` (show details)
  - [ ] `GET /api/v1/tv/:id/seasons`
  - [ ] `GET /api/v1/tv/:id/seasons/:seasonNum`
  - [ ] `GET /api/v1/tv/:id/seasons/:seasonNum/episodes`
  - [ ] `GET /api/v1/tv/:id/episodes/:episodeId`
  - [ ] `GET /api/v1/tv/:id/cast`
  - [ ] `GET /api/v1/tv/:id/similar`
  - [ ] `GET /api/v1/tv/:id/next-up` (next episode to watch)
  - [ ] `POST /api/v1/tv/episodes/:id/progress`
  - [ ] `GET /api/v1/tv/episodes/:id/progress`
  - [ ] `POST /api/v1/tv/:id/mark-watched`
  - [ ] `POST /api/v1/tv/:id/seasons/:seasonNum/mark-watched`
  - [ ] `POST /api/v1/tv/:id/refresh`

- [ ] **River Jobs** (`internal/content/tvshow/jobs.go`)
  - [ ] TVShowMetadataRefreshJob
  - [ ] TVShowLibraryScanJob
  - [ ] TVShowFileMatchJob
  - [ ] SeasonMetadataRefreshJob
  - [ ] EpisodeMetadataRefreshJob

- [ ] **fx Module** (`internal/content/tvshow/module.go`)

- [ ] **Tests**
  - [ ] Unit tests (80%+ coverage)
  - [ ] Integration tests

### TheTVDB Integration

- [ ] **TheTVDB Client** (`internal/service/metadata/thetvdb/client.go`)
  - [ ] API v4 implementation
  - [ ] JWT authentication
  - [ ] Token refresh
  - [ ] Rate limiting (20 req/10s)

- [ ] **TheTVDB Service** (`internal/service/metadata/thetvdb/service.go`)
  - [ ] Search series
  - [ ] Get series details
  - [ ] Get season details
  - [ ] Get episode details
  - [ ] Get series episodes (all)
  - [ ] Get artwork (posters, banners, backgrounds)
  - [ ] Get actors

- [ ] **Image Handler** (`internal/service/metadata/thetvdb/images.go`)
  - [ ] Artwork download/cache
  - [ ] Image proxy

- [ ] **Tests**
  - [ ] Unit tests with mock API

### TMDb TV Support

- [ ] **TMDb TV Service Extension** (`internal/service/metadata/tmdb/tv.go`)
  - [ ] Search TV show
  - [ ] Get TV show details
  - [ ] Get season details
  - [ ] Get episode details
  - [ ] Get TV credits
  - [ ] Get TV images

### Sonarr Integration

- [ ] **Sonarr Client** (`internal/service/metadata/sonarr/client.go`)
  - [ ] API v3 implementation
  - [ ] Authentication (API key)
  - [ ] Error handling

- [ ] **Sonarr Service** (`internal/service/metadata/sonarr/service.go`)
  - [ ] Get all series
  - [ ] Get series by ID
  - [ ] Get episodes for series
  - [ ] Get episode files
  - [ ] Sync library (Sonarr → Revenge)
  - [ ] Trigger series refresh
  - [ ] Get quality profiles
  - [ ] Get root folders

- [ ] **Sync Logic** (`internal/service/metadata/sonarr/sync.go`)
  - [ ] Full sync (initial)
  - [ ] Incremental sync
  - [ ] File path mapping
  - [ ] Episode matching

- [ ] **Webhook Handler**
  - [ ] `POST /api/v1/webhooks/sonarr`
  - [ ] Handle: Grab, Download, Rename, Delete events

- [ ] **Handler** (`internal/api/sonarr_handler.go`)
  - [ ] `GET /api/v1/admin/integrations/sonarr/status`
  - [ ] `POST /api/v1/admin/integrations/sonarr/sync`
  - [ ] `GET /api/v1/admin/integrations/sonarr/quality-profiles`

- [ ] **River Jobs**
  - [ ] SonarrSyncJob
  - [ ] SonarrWebhookJob

- [ ] **Tests**
  - [ ] Unit tests with mock API

### Episode Watch Progress

- [ ] **Continue Watching Logic** (`internal/content/tvshow/continue_watching.go`)
  - [ ] Track per-episode progress
  - [ ] Determine "next episode" logic
  - [ ] Handle season transitions
  - [ ] Handle series completion
  - [ ] Resume position calculation

- [ ] **Watch Status Types**
  - [ ] Unwatched
  - [ ] In Progress (with position)
  - [ ] Watched
  - [ ] Mark as watched (without playing)

- [ ] **Aggregation**
  - [ ] Season watch percentage
  - [ ] Show watch percentage
  - [ ] Recently watched shows
  - [ ] Continue watching list

### Search Integration

- [ ] **TV Show Collection Schema**
  ```json
  {
    "name": "tv_shows",
    "fields": [
      {"name": "id", "type": "string"},
      {"name": "title", "type": "string"},
      {"name": "original_title", "type": "string"},
      {"name": "overview", "type": "string"},
      {"name": "year", "type": "int32"},
      {"name": "genres", "type": "string[]"},
      {"name": "cast", "type": "string[]"},
      {"name": "status", "type": "string"},
      {"name": "rating", "type": "float"},
      {"name": "added_at", "type": "int64"}
    ]
  }
  ```

- [ ] **Episode Collection Schema** (optional, for episode search)
  ```json
  {
    "name": "tv_episodes",
    "fields": [
      {"name": "id", "type": "string"},
      {"name": "show_id", "type": "string"},
      {"name": "show_title", "type": "string"},
      {"name": "title", "type": "string"},
      {"name": "season_number", "type": "int32"},
      {"name": "episode_number", "type": "int32"},
      {"name": "overview", "type": "string"}
    ]
  }
  ```

- [ ] **Search Service Updates**
  - [ ] Index TV show
  - [ ] Index episode (optional)
  - [ ] Search TV shows
  - [ ] Multi-type search (movies + shows)

### Frontend Updates

- [ ] **TV Shows Grid** (`/tv`)
  - [ ] Show cards with poster
  - [ ] Watch progress indicator
  - [ ] Filtering by genre, year, status
  - [ ] Sorting options

- [ ] **TV Show Detail** (`/tv/[id]`)
  - [ ] Hero backdrop
  - [ ] Show metadata
  - [ ] Season selector
  - [ ] Episode list per season
  - [ ] Cast carousel
  - [ ] Similar shows
  - [ ] Continue watching button

- [ ] **Season View** (`/tv/[id]/season/[num]`)
  - [ ] Episode list with thumbnails
  - [ ] Watch status per episode
  - [ ] Mark season as watched

- [ ] **Episode Detail** (modal or page)
  - [ ] Episode still image
  - [ ] Title, overview
  - [ ] Air date
  - [ ] Guest stars
  - [ ] Play button

- [ ] **Continue Watching Widget**
  - [ ] Dashboard component
  - [ ] Show next episode card
  - [ ] Resume button

- [ ] **Search Updates**
  - [ ] Include TV shows in results
  - [ ] Type filter (movies/tv)

- [ ] **Admin: Sonarr Integration**
  - [ ] Settings page
  - [ ] Connection test
  - [ ] Manual sync button

---

## Verification Checklist

- [ ] TV shows display in frontend
- [ ] Seasons and episodes hierarchically organized
- [ ] Sonarr sync imports shows and episodes
- [ ] Episode watch progress tracks correctly
- [ ] "Next episode" logic works
- [ ] Search includes TV shows
- [ ] TheTVDB metadata enriches shows
- [ ] All tests pass (80%+ coverage)
- [ ] CI pipeline passes

---

## Dependencies from SOURCE_OF_TRUTH

| Package | Version | Purpose |
|---------|---------|---------|
| github.com/go-resty/resty/v2 | v2.17.1 | HTTP client (Sonarr, TheTVDB) |

---

## Related Documentation

- [ROADMAP.md](ROADMAP.md) - Full roadmap overview
- [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md) - Authoritative versions
- [TVSHOW_MODULE.md](../features/video/TVSHOW_MODULE.md) - TV show module design
- [THETVDB.md](../integrations/metadata/video/THETVDB.md) - TheTVDB integration
- [SONARR.md](../integrations/servarr/SONARR.md) - Sonarr integration
- [WATCH_NEXT_CONTINUE_WATCHING.md](../features/playback/WATCH_NEXT_CONTINUE_WATCHING.md) - Watch progress design
