# Module Implementation TODO

> Detailed implementation checklist for the modular Revenge architecture.
> Each module is fully self-contained and can be implemented independently.

## Implementation Order

1. **Phase 1: Core Infrastructure** (shared services, River, cache, search)
2. **Phase 2: Movie Module** (reference implementation)
3. **Phase 3: TV Show Module** (shares patterns with movie)
4. **Phase 4: Music Module** (complex relations)
5. **Phase 5: Remaining Modules** (audiobook, book, podcast, photo, livetv)
6. **Phase 6: Adult Modules** (isolated `c` schema)

---

## Phase 1: Core Infrastructure

### 1.1 New Dependencies

- [ ] Add `github.com/redis/go-redis/v9` (Dragonfly/cache)
- [ ] Add `github.com/typesense/typesense-go/v4` (search)
- [ ] Add `github.com/riverqueue/river` (job queue)
- [ ] Add `github.com/riverqueue/river/riverdriver/riverpgxv5`
- [ ] Add `github.com/ogen-go/ogen` (OpenAPI codegen)

### 1.2 Infrastructure Services

- [ ] `internal/infra/cache/cache.go` - Dragonfly client (go-redis)
- [ ] `internal/infra/search/search.go` - Typesense client
- [ ] `internal/infra/jobs/river.go` - River job queue setup
- [ ] `internal/infra/jobs/workers.go` - Base worker registration

### 1.3 Shared Tables Migration

- [ ] `migrations/shared/000001_extensions.sql` - Already exists
- [ ] `migrations/shared/000002_users.sql` - Already exists, verify
- [ ] `migrations/shared/000003_sessions.sql` - Already exists, verify
- [ ] `migrations/shared/000004_oidc.sql` - Already exists, verify
- [ ] `migrations/shared/000005_libraries.sql` - Add module_type column
- [ ] `migrations/shared/000006_api_keys.sql` - NEW
- [ ] `migrations/shared/000007_server_settings.sql` - NEW
- [ ] `migrations/shared/000008_activity_log.sql` - Add module column
- [ ] River migrations - Handled by River CLI

### 1.4 Video Playlists (shared movie + tvshow)

- [ ] `migrations/shared/000010_video_playlists.sql`
  - `video_playlists` (id, user_id, name, description, is_public, created_at, updated_at)
  - `video_playlist_items` (id, playlist_id, item_type, item_id, position, added_at)

### 1.5 Audio Playlists (shared music + audiobook + podcast)

- [ ] `migrations/shared/000011_audio_playlists.sql`
  - `audio_playlists`
  - `audio_playlist_items`

### 1.6 Video Collections (shared movie + tvshow)

- [ ] `migrations/shared/000012_video_collections.sql`
  - `video_collections`
  - `video_collection_movies`
  - `video_collection_episodes`

### 1.7 Audio Collections (shared music + audiobook)

- [ ] `migrations/shared/000013_audio_collections.sql`
  - `audio_collections`
  - `audio_collection_tracks`
  - `audio_collection_audiobooks`

### 1.8 OpenAPI Setup (ogen)

- [ ] `api/openapi/revenge.yaml` - Main OpenAPI spec
- [ ] `api/openapi/movies.yaml` - Movie endpoints
- [ ] `api/openapi/shows.yaml` - TV show endpoints
- [ ] `api/generate.go` - `//go:generate ogen ...`

---

## Phase 2: Movie Module

### 2.1 Database Migration `migrations/movie/`

#### Core Tables
- [ ] `000001_movies.up.sql`
  ```sql
  movies (
    id UUID PRIMARY KEY,
    library_id UUID REFERENCES libraries(id),
    title VARCHAR(500) NOT NULL,
    sort_title VARCHAR(500),
    original_title VARCHAR(500),
    tagline TEXT,
    overview TEXT,
    path TEXT NOT NULL,
    runtime_ticks BIGINT,
    release_date DATE,
    year INT,
    budget BIGINT,
    revenue BIGINT,
    status VARCHAR(50),
    tmdb_id INT,
    imdb_id VARCHAR(20),
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
  )
  ```

#### People & Credits
- [ ] `000002_movie_people.up.sql`
  ```sql
  movie_people (id, name, tmdb_id, imdb_id, image_path, biography, birth_date, death_date)
  movie_cast (movie_id, person_id, character, "order")
  movie_crew (movie_id, person_id, department, job)
  ```

#### Studios
- [ ] `000003_movie_studios.up.sql`
  ```sql
  movie_studios_list (id, name, tmdb_id, logo_path)
  movie_studios (movie_id, studio_id)
  ```

#### Media Files
- [ ] `000004_movie_streams.up.sql`
  ```sql
  movie_streams (
    id, movie_id, stream_index, stream_type,
    codec, language, title, is_default, is_forced,
    width, height, bitrate, channels, sample_rate
  )
  movie_subtitles (
    id, movie_id, language, path, format,
    is_external, is_forced, is_default
  )
  movie_chapters (id, movie_id, start_ticks, title)
  ```

#### Genres & Tags
- [ ] `000005_movie_genres.up.sql`
  ```sql
  movie_genres (movie_id, genre_id)  -- FK to public.genres where domain='movie'
  movie_tags (id, movie_id, user_id, tag)
  ```

#### Images
- [ ] `000006_movie_images.up.sql`
  ```sql
  movie_images (
    id, movie_id, type, path,
    width, height, blurhash, provider
  )
  ```

#### User Data
- [ ] `000007_movie_user_data.up.sql`
  ```sql
  movie_user_ratings (id, user_id, movie_id, score, created_at, updated_at)
  movie_external_ratings (id, movie_id, source, score, vote_count, url, fetched_at)
  movie_favorites (user_id, movie_id, added_at)
  movie_watchlist (user_id, movie_id, added_at)
  movie_history (id, user_id, movie_id, position_ticks, completed, watched_at)
  ```

### 2.2 SQLC Queries `queries/movie/`

- [ ] `movies.sql` - CRUD for movies
- [ ] `movie_people.sql` - Cast and crew queries
- [ ] `movie_streams.sql` - Stream/subtitle/chapter queries
- [ ] `movie_ratings.sql` - User and external ratings
- [ ] `movie_user_data.sql` - Favorites, watchlist, history

### 2.3 Domain Layer `internal/content/movie/`

- [ ] `entity.go` - Movie, MoviePerson, MovieStudio, etc.
- [ ] `repository.go` - Interface definitions
- [ ] `repository_pg.go` - PostgreSQL implementation

### 2.4 Service Layer

- [ ] `service.go` - Business logic
- [ ] `scanner.go` - File system scanner
- [ ] `provider_tmdb.go` - TMDb metadata provider
- [ ] `provider_omdb.go` - OMDb metadata provider

### 2.5 API Layer (ogen-generated)

- [ ] `api/openapi/movies.yaml` - OpenAPI spec
- [ ] `api/generated/movies/` - ogen-generated handlers
- [ ] `handler.go` - Implement ogen interfaces
- [ ] Routes:
  - `GET /api/v1/movies` - List movies
  - `GET /api/v1/movies/{id}` - Get movie
  - `GET /api/v1/movies/{id}/similar` - Similar movies
  - `GET /api/v1/movies/{id}/cast` - Movie cast
  - `GET /api/v1/movies/{id}/streams` - Stream info
  - `POST /api/v1/movies/{id}/rate` - Rate movie
  - `POST /api/v1/movies/{id}/favorite` - Add to favorites
  - `DELETE /api/v1/movies/{id}/favorite` - Remove from favorites
  - `POST /api/v1/movies/{id}/watchlist` - Add to watchlist
  - `GET /api/v1/movies/{id}/play` - Get playback URL (returns Blackbeard URL)
  - `POST /api/v1/movies/{id}/progress` - Update progress

### 2.6 River Jobs

- [ ] `jobs.go` - River job definitions:
  ```go
  type ScanMovieLibraryArgs struct { ... }
  type FetchMovieMetadataArgs struct { ... }
  type IndexMovieArgs struct { ... }  // Typesense
  ```

### 2.7 Module Registration

- [ ] `module.go` - fx.Module with all providers

---

## Phase 3: TV Show Module

### 3.1 Database Migration `migrations/tvshow/`

#### Core Tables
- [ ] `000001_series.up.sql`
  ```sql
  series (id, library_id, title, sort_title, overview, status,
          network, first_air_date, last_air_date, tvdb_id, tmdb_id, imdb_id, ...)
  seasons (id, series_id, season_number, name, overview, air_date, ...)
  episodes (id, season_id, series_id, episode_number, title, overview,
            runtime_ticks, air_date, path, ...)
  ```

#### People & Credits
- [ ] `000002_series_people.up.sql`
  ```sql
  series_people (id, name, tvdb_id, tmdb_id, imdb_id, ...)
  series_cast (series_id, person_id, character, "order")
  series_crew (series_id, person_id, department, job)
  episode_cast (episode_id, person_id, character)  -- guest stars
  ```

#### Studios
- [ ] `000003_series_studios.up.sql`

#### Media Files
- [ ] `000004_episode_streams.up.sql`
  ```sql
  episode_streams (id, episode_id, ...)
  episode_subtitles (id, episode_id, ...)
  episode_chapters (id, episode_id, ...)
  ```

#### User Data
- [ ] `000005_tvshow_user_data.up.sql`
  ```sql
  series_user_ratings, episode_user_ratings
  series_external_ratings, episode_external_ratings
  series_favorites, series_watchlist
  episode_history
  ```

### 3.2 - 3.6 (Same pattern as movie)

---

## Phase 4: Music Module

### 4.1 Database Migration `migrations/music/`

#### Core Tables
- [ ] `000001_music_core.up.sql`
  ```sql
  artists (id, library_id, name, sort_name, overview,
           musicbrainz_id, spotify_id, ...)
  albums (id, library_id, artist_id, title, sort_title,
          release_date, album_type, musicbrainz_id, ...)
  tracks (id, album_id, artist_id, title, disc_number, track_number,
          duration_ticks, path, isrc, musicbrainz_id, ...)
  artist_albums (artist_id, album_id, role)  -- for multi-artist albums
  track_artists (track_id, artist_id, role)  -- for featuring artists
  ```

#### Music-specific
- [ ] `000002_track_streams.up.sql`
  ```sql
  track_streams (id, track_id, codec, bitrate, channels, sample_rate, ...)
  track_lyrics (id, track_id, content, is_synced, language, source)
  ```

#### User Data
- [ ] `000003_music_user_data.up.sql`
  ```sql
  artist_user_ratings, album_user_ratings, track_user_ratings
  artist_external_ratings, album_external_ratings, track_external_ratings
  artist_favorites, album_favorites, track_favorites
  track_history (id, user_id, track_id, played_at)  -- no position, tracks are short
  ```

---

## Phase 5: Remaining Modules

### 5.1 Audiobook Module
- [ ] `audiobooks`, `audiobook_chapters`
- [ ] `audiobook_authors`, `audiobook_narrators`
- [ ] `audiobook_user_ratings`, `audiobook_favorites`
- [ ] `chapter_progress`

### 5.2 Book Module
- [ ] `books`, `book_authors`
- [ ] `book_user_ratings`, `book_favorites`
- [ ] `book_progress`

### 5.3 Podcast Module
- [ ] `podcasts`, `podcast_episodes`
- [ ] `podcast_favorites`, `episode_progress`

### 5.4 Photo Module
- [ ] `photos`, `photo_albums`
- [ ] `photo_exif` (camera, GPS, etc.)

### 5.5 Live TV Module
- [ ] `channels`, `programs`, `recordings`
- [ ] `recording_streams`, `channel_favorites`

---

## Phase 6: Adult Modules (Isolated `c` Schema)

> **Note:** Adult modules use obscured schema `c` and API namespace `/c/`.
> No content ratings (implicit 18+). Separate access scope required.

### 6.1 Schema Setup
- [ ] `migrations/c/000001_schema.up.sql`
  ```sql
  CREATE SCHEMA IF NOT EXISTS c;
  ```

### 6.2 Adult Movie Module
- [ ] `000002_c_movies.up.sql`
  ```sql
  c.movies (id, library_id, title, path, release_date, ...)
  c.scenes (id, movie_id, title, start_ticks, end_ticks, ...)
  ```

- [ ] `000003_c_performers.up.sql`
  ```sql
  c.performers (id, name, aliases, gender, birth_date,
                ethnicity, measurements, image_path, ...)
  c.movie_performers (movie_id, performer_id)
  c.scene_performers (scene_id, performer_id)
  ```

- [ ] `000004_c_studios.up.sql`
  ```sql
  c.studios (id, name, logo_path, website)
  c.movie_studios (movie_id, studio_id)
  ```

- [ ] `000005_c_tags.up.sql`
  ```sql
  c.tags (id, name, category)
  c.movie_tags (movie_id, tag_id)
  c.scene_tags (scene_id, tag_id)
  ```

- [ ] `000006_c_movie_media.up.sql`
  ```sql
  c.movie_streams, c.movie_subtitles, c.movie_chapters
  c.movie_images, c.scene_images, c.performer_images
  ```

- [ ] `000007_c_movie_user_data.up.sql`
  ```sql
  c.movie_user_ratings, c.scene_user_ratings
  c.performer_user_ratings, c.studio_user_ratings
  c.movie_favorites, c.scene_favorites
  c.performer_favorites, c.studio_favorites
  c.movie_history
  ```

### 6.3 Adult Show Module
- [ ] `000010_c_series.up.sql`
  ```sql
  c.series, c.seasons, c.episodes
  c.episode_performers
  c.series_studios
  c.episode_streams, c.episode_subtitles
  c.series_user_ratings, c.episode_user_ratings
  c.series_favorites, c.episode_history
  ```

### 6.4 Adult Playlists & Collections
- [ ] `000015_c_playlists.up.sql`
  ```sql
  c.playlists (id, user_id, name, ...)
  c.playlist_items (id, playlist_id, item_type, item_id, position)
  c.collections (id, name, ...)
  c.collection_items (id, collection_id, item_type, item_id)
  ```

### 6.5 Adult API Endpoints

- [ ] `api/openapi/c_movies.yaml` - NOT included in public docs
- [ ] Routes under `/api/v1/c/`:
  - `GET /c/movies` - List adult movies
  - `GET /c/movies/{id}` - Get adult movie
  - `GET /c/shows` - List adult shows
  - Requires special auth scope `adult:read`

---

## Cross-Cutting Concerns

### Search Integration (Typesense)
- [ ] One collection per module (including `c_movies`, `c_series`)
- [ ] Unified search aggregation endpoint
- [ ] Module-specific facets
- [ ] Adult content NOT in unified search (separate endpoint)

### Job Queue (River)
- [ ] Queue definitions:
  - `scanning` - Library scans (MaxWorkers: 2)
  - `metadata` - Metadata fetching (MaxWorkers: 10)
  - `indexing` - Typesense updates (MaxWorkers: 5)
- [ ] Per-module job types
- [ ] Job history in River tables

### External Transcoding (Blackbeard)
- [ ] `/play` endpoints return Blackbeard URLs
- [ ] Profile negotiation (h264_720p, h265_1080p, etc.)
- [ ] Direct stream option for compatible clients

### Caching (Dragonfly)
- [ ] Session caching
- [ ] Metadata cache
- [ ] Search result cache (30s TTL)

---

## Migration Strategy

### From Current Schema

1. Keep existing shared tables (users, sessions, libraries, oidc)
2. **DELETE** old migrations 000006-000015 (media_items, images, etc.)
3. Create new per-module migration folders
4. Run River migrations
5. Implement modules one at a time

### Folder Restructuring

```
Current:
  internal/service/movie/    ❌ Move
  internal/service/tvshow/   ❌ Move

Target:
  internal/content/movie/    ✅
  internal/content/tvshow/   ✅
  internal/content/c/movie/  ✅
  internal/content/c/show/   ✅

  internal/service/          ✅ Keep shared services only
    auth/
    user/
    oidc/
    library/
```

### Rollback Plan

- Each module has independent down migrations
- Can disable modules without affecting others
- Shared infrastructure remains stable
- River jobs can be paused per queue
