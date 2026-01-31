# Revenge - Architecture v2

> Complete modular architecture for a next-generation media server.
> Ground-up design with full module isolation, no shared content tables.

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Complete architecture specification |
| Sources | ⚪ | N/A - internal architecture doc |
| Instructions | 🔴 | |
| Code | 🔴 | Reset to template |
| Linting | 🔴 | |
| Unit Testing | 🔴 | |
| Integration Testing | 🔴 | |

**Priority**: 🔴 HIGH
**Module**: Core architecture
**Dependencies**: [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md)

---

## Implementation Checklist

### Phase 1: Project Foundation
- [ ] Initialize Go module structure
- [ ] Set up fx dependency injection framework
- [ ] Configure structured logging (slog)
- [ ] Set up configuration management (koanf)

### Phase 2: Database Layer
- [ ] Set up PostgreSQL with pgx driver
- [ ] Initialize sqlc code generation
- [ ] Create base migrations (extensions, shared schema)
- [ ] Set up River job queue

### Phase 3: API Layer
- [ ] Set up ogen code generation from OpenAPI specs
- [ ] Configure HTTP server with Go 1.22 routing
- [ ] Add middleware (authentication, logging, CORS)
- [ ] Set up WebSocket support

### Phase 4: Infrastructure
- [ ] Set up Typesense search client
- [ ] Configure Dragonfly caching (Redis-compatible)
- [ ] Add health check endpoints
- [ ] Set up metrics/observability

---

## Design Principles

1. **Maximum Isolation** - Each content module is fully self-contained
2. **No Shared Content Tables** - Every module has its own optimized tables
3. **Per-Module Everything** - Ratings, history, favorites, metadata all per module
4. **Adult Content Isolation** - Separate PostgreSQL schema (`qar`) for complete separation (see [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#qar-obfuscation-terminology))
5. **Enable/Disable Modules** - Each module can be independently enabled/disabled
6. **Type Safety** - No polymorphic references, proper FK constraints everywhere
7. **External Transcoding** - Delegate to "Blackbeard" service for scalability
8. **Clustering** - Raft consensus for HA deployments (see [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#distributed-consensus-raft))
9. **Network QoS** - Priority hierarchy (P0-P4) ensures user experience over background jobs (see [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#network-qos-design-principle))
10. **Orchestration-Ready** - Docker Compose, K8s, K3s, Swarm support (see [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#container-orchestration))

---

## Technology Stack

| Component | Technology | Purpose |
|-----------|------------|---------|
| Language | Go 1.25+ | Core runtime |
| Database | PostgreSQL 18+ | Primary data store |
| Cache | Dragonfly | Redis-compatible, session/cache |
| Search | Typesense | Full-text search per module |
| Job Queue | River | PostgreSQL-native background jobs |
| API Docs | ogen | OpenAPI spec-first code generation |
| DI | uber-go/fx | Dependency injection |
| Config | koanf | Configuration management |
| Migrations | golang-migrate | Database schema versioning |
| SQL | sqlc | Type-safe query generation |
| HTTP | net/http (stdlib) | Go 1.22+ routing patterns |
| Logging | log/slog (stdlib) | Structured logging |

### Dependencies

> See [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#go-dependencies-core) for complete package list with versions.

Key packages by category:
- **Core**: pgx, koanf, fx
- **Infrastructure**: rueidis, otter, sturdyc, typesense-go, river
- **API**: ogen, go-faster/jx
- **HTTP & WebSocket**: resty, gobwas/ws
- **Media**: fsnotify, go-astiav, govips
- **Database**: golang-migrate, sqlc

---

## Module Overview

### Content Modules (12)

| Module | Schema | Description |
|--------|--------|-------------|
| `movie` | public | Movies, trailers |
| `tvshow` | public | Series, seasons, episodes |
| `music` | public | Artists, albums, tracks, music videos |
| `audiobook` | public | Audiobooks, chapters |
| `book` | public | E-books |
| `podcast` | public | Podcasts, episodes |
| `photo` | public | Photos, albums |
| `livetv` | public | Channels, programs, DVR recordings |
| `comics` | public | Comics, manga, graphic novels |
| `collection` | public | Cross-module collections (video/audio pools) |
| `qar` | qar | Adult content (expeditions, voyages, crew) |

> **Note:** Adult content uses schema `qar` with "Queen Anne's Revenge" obfuscation terminology. See [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#qar-obfuscation-terminology) for mapping.

### Shared Infrastructure

| Component | Description |
|-----------|-------------|
| `users` | User accounts, authentication |
| `profiles` | User profiles under accounts (Netflix-style) |
| `sessions` | Active sessions, devices |
| `resource_grants` | Polymorphic resource access grants |
| `api_keys` | External API authentication |
| `server_settings` | Persisted configuration |
| `activity_log` | Partitioned audit log (async writes, auto-cleanup) |
| `oidc` | SSO/OIDC provider configuration |
| `river_*` | Job queue tables (managed by River) |

> **Note:** Libraries are per-module (e.g., `movie_libraries`, `tv_libraries`, `qar.fleets`).
> See [LIBRARY_TYPES.md](../features/shared/LIBRARY_TYPES.md) for details.

---

## Project Structure

```
internal/
  service/                    # Shared services (cross-cutting)
    auth/                     # Authentication
    user/                     # User management
    session/                  # Session handling
    oidc/                     # SSO/OIDC
    grants/                   # Polymorphic resource grants
    playback/                 # Playback session management
      client.go               # Client detection & capabilities
      bandwidth.go            # Bandwidth monitoring (external)
      transcoder.go           # Blackbeard integration
      session.go              # Playback session state
      buffer.go               # HLS/DASH segment buffering
      fileserver.go           # Raw file HTTP streaming to Blackbeard
      stream_handler.go       # Unified stream handler with buffering
      transcode_cache.go      # Memory-aware transcode cache
      disk_cache.go           # Persistent disk cache with quotas
      profile.go              # Transcode profiles & device groups
      module.go               # fx module registration

  content/                    # Content modules (isolated)
    movie/
      entity.go               # Domain entities
      repository.go           # Repository interface
      service.go              # Business logic
      handler.go              # HTTP handlers
      scanner.go              # File scanner
      provider_tmdb.go        # Metadata provider
      jobs.go                 # River job definitions
      module.go               # fx.Module registration
    tvshow/
    music/
    audiobook/
    book/
    podcast/
    photo/
    livetv/
    collection/
    qar/                      # Adult modules (Queen Anne's Revenge)
      expedition/             # Full-length movies
      voyage/                 # Individual scenes
      crew/                   # Performers
      port/                   # Studios
      flag/                   # Tags
      fleet/                  # Libraries

  infra/
    database/                 # Shared DB infrastructure
      migrations/
        shared/               # Users, sessions, grants, RBAC
        movie/                # Movie module (incl. movie_libraries)
        tvshow/               # TV module (incl. tv_libraries)
        music/                # Music module (incl. music_libraries)
        qar/                  # Adult schema (incl. qar.fleets)
      queries/
        movie/
        tvshow/
        ...
    cache/                    # Dragonfly client
    search/                   # Typesense client
    jobs/                     # River job queue setup

api/
  openapi/
    revenge.yaml              # Main OpenAPI spec
    movies.yaml               # Movie endpoints
    shows.yaml                # TV show endpoints
    ...
  generated/                  # ogen-generated code
```

---

## Job Queue (River)

All background work runs through River:

| Job Type | Description |
|----------|-------------|
| Library Scanning | Scan folders for new media |
| Metadata Fetching | TMDb, TheTVDB, MusicBrainz, etc. |
| Image Processing | Download posters, generate blurhash |
| Search Indexing | Update Typesense on changes |
| Cleanup Tasks | Remove orphaned files |
| Scheduled Refresh | Re-fetch metadata periodically |
| Notification Jobs | Webhook calls, email alerts |

```go
// Example job definition
type ScanLibraryArgs struct {
    LibraryID uuid.UUID `json:"library_id"`
    FullScan  bool      `json:"full_scan"`
}

func (ScanLibraryArgs) Kind() string { return "scan_library" }

type ScanLibraryWorker struct {
    river.WorkerDefaults[ScanLibraryArgs]
    scanner *LibraryScanner
}

func (w *ScanLibraryWorker) Work(ctx context.Context, job *river.Job[ScanLibraryArgs]) error {
    return w.scanner.Scan(ctx, job.Args.LibraryID, job.Args.FullScan)
}
```

---

## External Transcoding (Blackbeard)

Revenge does **not** transcode internally. Instead, it delegates to an external service called "Blackbeard". The stream always flows through Revenge to maintain access control and session management.

### Architecture

```
┌──────────┐    1. Play     ┌──────────┐   2. Transcode   ┌────────────┐
│  Client  │ ─────────────→ │ Revenge  │ ───────────────→ │ Blackbeard │
└──────────┘                └──────────┘                  └────────────┘
     ▲                           │ ▲                            │
     │                           │ │                            │
     │    5. Stream (proxied)    │ │    4. Transcoded stream    │
     └───────────────────────────┘ └────────────────────────────┘
                                   │
                                   │ 3. Raw file
                                   ▼
                             ┌──────────┐
                             │  Storage │
                             └──────────┘
```

### Playback Flow

1. **Client requests playback** from Revenge with session info
2. **Revenge detects client capabilities** (resolution, codecs, bandwidth)
3. **Revenge determines transcoding needs** based on source file + client
4. **If transcoding needed**: Revenge requests stream from Blackbeard with profile
5. **Blackbeard fetches raw file** from Revenge (internal API)
6. **Blackbeard transcodes on-the-fly** and returns stream to Revenge
7. **Revenge proxies stream** to client (maintains auth, tracks progress)

### Client Detection

Revenge tracks client capabilities for optimal streaming:

| Data | Source | Purpose |
|------|--------|---------|
| Device type | User-Agent, app registration | Base profile selection |
| Max resolution | Client capability report | Resolution cap |
| Supported codecs | Client capability report | Direct play vs transcode |
| HDR support | Client capability report | Tone mapping decision |
| Audio channels | Client capability report | Audio downmix |

### Bandwidth Monitoring (External Clients)

For clients outside the local network, Revenge monitors bandwidth:

```go
type ClientSession struct {
    ID              uuid.UUID
    UserID          uuid.UUID
    DeviceID        string
    IsExternal      bool              // Outside local network

    // Bandwidth tracking (external only)
    BandwidthKbps   int               // Current measured bandwidth
    BandwidthJitter int               // Variance in kbps
    BandwidthSamples []BandwidthSample // Rolling window

    // Client capabilities
    MaxResolution   string            // "4k", "1080p", "720p", etc.
    SupportedCodecs []string          // ["h264", "hevc", "av1"]
    SupportsHDR     bool
    MaxAudioChannels int
}

type BandwidthSample struct {
    Timestamp   time.Time
    Kbps        int
    Latency     time.Duration
}
```

**Bandwidth adaptation:**
- Measure during initial buffer and periodically
- Track jitter (variance) to detect unstable connections
- Adapt quality profile dynamically (ABR-like behavior)
- Conservative bitrate selection: `targetBitrate = bandwidth * 0.8 - jitter`

### Transcode Request

```go
type TranscodeRequest struct {
    // Source
    MediaID     uuid.UUID
    StreamIndex int               // Video/audio stream index

    // Client constraints
    MaxWidth       int            // From client resolution
    MaxHeight      int
    MaxBitrate     int            // From bandwidth measurement
    TargetCodec    string         // Preferred output codec

    // Bandwidth info (external clients)
    BandwidthKbps  int
    JitterKbps     int
    IsExternal     bool

    // Playback
    StartPosition  time.Duration
    AudioStreamIdx int
    SubtitleIdx    *int
}
```

### Why Proxy Through Revenge?

1. **Access control** - Validate session on every chunk
2. **Progress tracking** - Know exactly what client watched
3. **Bandwidth monitoring** - Measure actual throughput
4. **Session management** - Handle pause, seek, stop
5. **Analytics** - Track quality switches, buffering events
6. **Single endpoint** - Client only talks to Revenge

### Blackbeard API (Internal)

Blackbeard exposes internal API for Revenge only:

```
POST /transcode/start
  → Returns stream ID + HLS/DASH manifest URL

GET /transcode/{id}/master.m3u8
  → HLS master playlist (multi-quality)

GET /transcode/{id}/{quality}/segment_{n}.ts
  → Video segment

DELETE /transcode/{id}
  → Stop transcoding, cleanup
```

### Internal Raw File Serving

Revenge exposes internal endpoints for Blackbeard to fetch raw media files:

```
GET /internal/stream/{mediaID}
  → Raw file with HTTP Range support
  → Requires internal token authentication
  → Supports byte-range requests for seeking

GET /internal/probe/{mediaID}
  → File metadata for Blackbeard analysis
```

**HTTP Range Support:**
- Full `bytes=start-end` range requests
- Partial content (206) responses
- Chunked streaming (64KB chunks)
- Context cancellation on client disconnect
- MIME type detection for video/audio/subtitle formats

### Stream Buffering

Revenge buffers transcoded segments between Blackbeard and clients for stability:

```
┌────────────┐       ┌───────────────────┐       ┌──────────┐
│ Blackbeard │ ────→ │ Segment Buffer    │ ────→ │  Client  │
└────────────┘       │ (5-8 segments)    │       └──────────┘
                     │ • Prefetch ahead  │
                     │ • Retry on error  │
                     │ • LRU eviction    │
                     └───────────────────┘
```

**Buffer Configuration:**
| Setting | Default | Description |
|---------|---------|-------------|
| `max_segments` | 8 | Maximum buffered segments |
| `min_buffer_duration` | 15s | Minimum buffer before playback |
| `max_buffer_duration` | 60s | Maximum buffer size |
| `prefetch_count` | 3 | Segments to prefetch ahead |
| `max_retries` | 3 | Retries per segment fetch |
| `idle_timeout` | 5m | Cleanup idle buffers |

**Buffer Benefits:**
1. **Error recovery** - Time to retry failed fetches
2. **Smooth playback** - Absorbs network jitter
3. **Quality switches** - Buffer during profile changes
4. **Blackbeard restart** - Keep playing during restarts

### Transcode Cache

Memory-aware caching of transcoded segments for reuse:

```
┌────────────────────────────────────────────────────────────┐
│                    Transcode Cache                         │
│                                                            │
│  Session A (active, priority=3)    Session B (idle, p=1)  │
│  ├─ segment_0.ts  [1.2MB]          ├─ segment_0.ts        │
│  ├─ segment_1.ts  [1.1MB]          └─ segment_1.ts        │
│  └─ segment_2.ts  [1.3MB]                                 │
│                                                            │
│  Memory: 450MB / 2GB (22%)                                │
│  Eviction: LRU with priority awareness                    │
└────────────────────────────────────────────────────────────┘
```

**Cache Configuration:**
| Setting | Default | Description |
|---------|---------|-------------|
| `max_memory` | 25% RAM | Maximum cache size |
| `max_segments_per_session` | 50 | ~5 min per transcode |
| `min_retention` | 30s | Minimum before eviction |
| `high_pressure` | 80% | Start evicting old segments |
| `critical_pressure` | 95% | Aggressive eviction |

**Cache Benefits:**
1. **Seek back** - Previously watched segments available
2. **Resume playback** - Segments survive pause/resume
3. **Multiple viewers** - Same transcode shared if parameters match
4. **Memory efficient** - Evict only when pressure requires

### Benefits

- **Revenge stays simple** - No FFmpeg, no heavy processing
- **Scalable transcoding** - Multiple Blackbeard instances
- **Regional deployment** - Blackbeard near storage, Revenge near users
- **Replaceable** - Swap transcoder without touching Revenge
- **Centralized control** - All access through Revenge API

---

## API Structure

### OpenAPI Spec-First with ogen

```yaml
# api/openapi/revenge.yaml
openapi: 3.1.0
info:
  title: Revenge API
  version: 1.0.0

paths:
  /api/v1/movies:
    $ref: './movies.yaml#/paths/~1movies'
  /api/v1/movies/{id}:
    $ref: './movies.yaml#/paths/~1movies~1{id}'
```

Generate handlers:
```bash
go generate ./api/...
# or: ogen --target api/generated --package api --clean api/openapi/revenge.yaml
```

### URL Structure

```
/api/v1/
  /movies/...
  /shows/...
  /music/...
  /audiobooks/...
  /books/...
  /podcasts/...
  /photos/...
  /livetv/...
  /collections/...

  /legacy/expeditions/...  # Adult movies (QAR namespace)
  /legacy/voyages/...      # Adult scenes (QAR namespace)
  /legacy/crew/...         # Adult performers (QAR namespace)

  /users/...            # Shared
  /libraries/...        # Shared
  /system/...           # Shared
  /jobs/...             # Job status (admin)
```

> **Security:** `/legacy/` endpoints require `legacy:read` auth scope, are not listed in public API docs, have separate rate limiting, and all access is audit-logged.
>
> See [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#api-namespaces) for complete API namespace documentation.

---

## Database Migrations

### Per-Module Structure

```
migrations/
  shared/
    000001_extensions.up.sql
    000002_users.up.sql
    000003_sessions.up.sql
    000004_libraries.up.sql
    000005_oidc.up.sql
  movie/
    000001_movies.up.sql
    000002_movie_people.up.sql
    000003_movie_streams.up.sql
    000004_movie_user_data.up.sql
  tvshow/
    000001_series.up.sql
    ...
  qar/
    000001_qar_schema.up.sql    # CREATE SCHEMA qar;
    000002_qar_expeditions.up.sql
    000003_qar_crew.up.sql
    ...
```

### Migration Order

1. Run `shared/` first (always)
2. Run module migrations based on enabled modules
3. River manages its own migrations

---

## Caching Strategy (Dragonfly)

| Data Type | TTL | Key Pattern |
|-----------|-----|-------------|
| Sessions | 24h | `session:{token}` |
| User profiles | 5m | `user:{id}` |
| Library lists | 1m | `libraries:{user_id}` |
| Search results | 30s | `search:{module}:{hash}` |
| Metadata cache | 1h | `meta:{provider}:{id}` |

```go
// Example usage
type CacheService struct {
    rdb *redis.Client
}

func (c *CacheService) GetSession(ctx context.Context, token string) (*Session, error) {
    key := fmt.Sprintf("session:%s", token)
    data, err := c.rdb.Get(ctx, key).Bytes()
    // ...
}
```

---

## Search Architecture (Typesense)

One collection per module:

| Collection | Fields |
|------------|--------|
| `movies` | title, overview, genres, year, cast |
| `series` | title, overview, genres, year, cast |
| `tracks` | title, artist, album, genre |
| `audiobooks` | title, author, narrator |
| `qar_expeditions` | title, crew, port, flags, year |
| `qar_voyages` | title, crew, expedition_id, flags |
| `qar_crew` | name, aliases, ports, flags |
| `qar_ports` | name, parent_port |
| `qar_treasures` | title, crew, port, flags |

> **Note**: QAR collections use obfuscated terminology. See [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#qar-obfuscation-terminology) for full mapping (expeditions=movies, voyages=scenes, crew=performers, ports=studios, treasures=galleries, flags=tags).

```go
// Example: Index a movie
func (s *SearchService) IndexMovie(ctx context.Context, movie *Movie) error {
    doc := map[string]interface{}{
        "id":       movie.ID.String(),
        "title":    movie.Title,
        "overview": movie.Overview,
        "year":     movie.Year,
        "genres":   movie.Genres,
    }
    _, err := s.client.Collection("movies").Documents().Upsert(ctx, doc)
    return err
}
```

---

## Shared Resources

### Playlist Pools (3)

| Pool | Modules | Tables |
|------|---------|--------|
| Video | movie, tvshow | `video_playlists`, `video_playlist_items` |
| Audio | music, audiobook, podcast | `audio_playlists`, `audio_playlist_items` |
| Adult | qar (expeditions, voyages) | Planned (`qar.playlists`, `qar.playlist_items`) |

### Collection Pools (3)

| Pool | Modules | Tables |
|------|---------|--------|
| Video | movie, tvshow | `video_collections`, `video_collection_movies`, `video_collection_episodes` |
| Audio | music, audiobook | `audio_collections`, `audio_collection_tracks`, `audio_collection_audiobooks` |
| Adult | qar (expeditions, voyages) | Planned (`qar.collections`, `qar.collection_items`) |

---

## Rating Systems

### Content Ratings (Age Restriction)
- Only for `movie` and `tvshow` modules
- Shared `content_ratings` system (MPAA, FSK, PEGI, etc.)
- Adult content has no content ratings (implicit 18+)

### User Ratings (Per Module)

| Module | Ratable Entities | Sync Services |
|--------|------------------|---------------|
| movie | movies | Trakt, Letterboxd, Simkl |
| tvshow | series, episodes | Trakt, Simkl |
| music | artists, albums, tracks | Last.fm, ListenBrainz |
| audiobook | audiobooks | Goodreads, Audible |
| book | books | Goodreads, OpenLibrary |
| adult_movie | movies, scenes, performers | — |
| adult_scene | scenes, performers | — |

### External Ratings (Per Module)

| Module | Sources |
|--------|---------|
| movie | IMDb, Rotten Tomatoes, Metacritic, TMDb |
| tvshow | IMDb, Rotten Tomatoes, TheTVDB |
| music | Last.fm, Spotify, RateYourMusic |
| audiobook | Audible, Goodreads |
| book | Goodreads, OpenLibrary |

---

## QAR Schema Isolation (Adult Content)

> See [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#qar-obfuscation-terminology) for complete QAR terminology mapping.

All adult content in separate PostgreSQL schema `qar` (Queen Anne's Revenge):

```sql
CREATE SCHEMA qar;

-- Tables (QAR terminology)
qar.expeditions      -- Movies (full-length releases)
qar.voyages          -- Scenes (individual scenes)
qar.crew             -- Performers
qar.ports            -- Studios
qar.flags            -- Tags
qar.fleets           -- Libraries
qar.treasures        -- Gallery images
qar.expedition_crew  -- Movie-performer relations
qar.voyage_crew      -- Scene-performer relations
qar.user_ratings
qar.user_favorites
qar.watch_history
```

**Benefits:**
- Complete data isolation
- Separate backup/restore (`pg_dump -n qar`)
- PostgreSQL-level access control (`REVOKE ALL ON SCHEMA qar FROM public_role`)
- Easy to purge entire schema
- Legal/compliance separation
- Obfuscated terminology for discretion

---

## Per-Module Tables Template

Each content module typically has:

### Core Tables
- `{items}` - Main content table (movies, tracks, etc.)
- `{item}_images` - Artwork/posters
- `{item}_streams` - Audio/video/subtitle streams in files
- `{item}_genres` - Genre assignments (domain-scoped)
- `{item}_tags` - User-defined tags

### People/Credits (domain-dependent)
**Video modules (movie, tvshow):** Share `video_people` table in `shared/` migrations
- Background workers enrich people data (TMDB, TVDB) and data overlaps 100% after enrichment
- `movie_credits` and `series_credits` both reference `video_people`

**Other modules:** Isolated people tables per module
- `music_artists` - Music artists (different metadata: discography, genres, etc.)
- `book_authors` - Book authors (different metadata: bibliography, awards, etc.)
- `comic_creators` - Comic creators (writers, artists, colorists, etc.)

**QAR module:** Completely isolated in schema `qar`
- `qar.crew` - Performer data with NSFW images, StashDB/TPDB metadata sources

### Studios (per module)
- `{module}_studios` - Production studios

### Video-specific (movie, tvshow, adult)
- `{item}_subtitles` - External subtitle files
- `{item}_chapters` - Chapter markers

### Music-specific
- `track_lyrics` - Lyrics with optional sync

### Photo-specific
- `photo_exif` - EXIF metadata, GPS, camera info

### User Data (per module)
- `{item}_user_ratings` - User scores
- `{item}_external_ratings` - Third-party scores
- `{item}_favorites` - User favorites
- `{item}_watchlist` - Watch later (video only)
- `{item}_history` - Watch/play history with progress

---

## Summary

| Aspect | Decision |
|--------|----------|
| Content isolation | Per-module tables, no shared content |
| Adult namespace | `/legacy/` API with `qar` schema |
| Transcoding | External "Blackbeard" service |
| Job queue | River (PostgreSQL-native) |
| API docs | ogen (OpenAPI spec-first) |
| Migrations | Per-module folders |
| Cache | Dragonfly (rueidis) + otter (local) + sturdyc (API) |
| Search | Typesense (per-module collections) |
| Frontend | SvelteKit 2 + Tailwind CSS 4 |
| Deployment | Docker Compose (primary), K8s/K3s/Swarm supported |
| Clustering | Raft consensus (hashicorp/raft) |

---

## Frontend Architecture

### Technology Stack

| Component | Technology | Purpose |
|-----------|------------|---------|
| Framework | SvelteKit 2 | SSR, routing, API |
| UI | Tailwind CSS 4 + shadcn-svelte | Styling, components |
| State | Svelte Stores + TanStack Query | Client/server state |
| Auth | JWT + OIDC | Authentication |
| i18n | Built-in | Internationalization |
| PWA | Service Worker | Offline support |

### Frontend Structure

```
web/
  src/
    lib/
      components/
        ui/              # shadcn-svelte components
        media/           # Media cards, players
        admin/           # Admin-specific components
      stores/            # Svelte stores (auth, theme, playback)
      api/               # Generated API client from OpenAPI
      utils/             # Helpers, formatters
    routes/
      (app)/
        (admin)/         # Admin panel (/admin/...)
          users/
          libraries/
          settings/
          logs/
        (media)/         # Media browsing
          movies/
          shows/
          music/
        (player)/        # Video/audio player (gapless, crossfade)
      (auth)/            # Login, register, OIDC
      api/               # API routes (BFF pattern)
    app.css
    app.html
  static/
  tailwind.config.ts
  svelte.config.js
```

### Player Features (WebUI)

| Feature | Implementation |
|---------|----------------|
| Gapless audio | Web Audio API, 30s prefetch + instant switch |
| Crossfade | Dual gain nodes, 5s overlapping crossfade |
| Synced lyrics | LRC format, binary search lookup |
| Visualizations | Canvas frequency bars OR fanart pulse effects |
| Quality switching | Seamless via WebSocket to Blackbeard |
| Subtitles | WebVTT external + container extraction |
| WebSocket sync | Watch Party, position tracking, quality changes |

**Technology:**
- Video: Shaka Player (DASH) + hls.js (HLS for non-Safari)
- Audio: Web Audio API + Howler.js wrapper
- Streaming: HLS primary, DASH fallback, Progressive last resort

### Role-Based Access Control (RBAC)

| Role | Permissions |
|------|-------------|
| `admin` | Full access, server settings, user management |
| `moderator` | Manage libraries, metadata, moderate content |
| `user` | Browse, play, rate, create playlists |
| `guest` | Browse only (configurable) |

### Theme System

- Light/Dark mode with system preference detection
- CSS variables for full customization
- Per-user theme preference stored in settings
- Admin-configurable default theme

---

## Client Profiles (Blackbeard)

### Device Group Mapping

| Group | User-Agent Patterns | Default Profile |
|-------|---------------------|-----------------|
| `tv_4k` | tizen, webos, android tv (4K capable) | `hevc_4k_hdr` |
| `tv_hd` | roku, fire tv stick, older TVs | `h264_1080p` |
| `mobile_ios` | iphone, ipad | `h264_1080p_hls` |
| `mobile_android` | android mobile | `h264_1080p_hls` |
| `desktop_app` | electron, revenge-desktop | `hevc_4k` |
| `browser_modern` | chrome, firefox, edge | `vp9_1080p_dash` |
| `browser_legacy` | safari, older browsers | `h264_720p_hls` |
| `low_bandwidth` | any (bandwidth < 3 Mbps) | `h264_480p_hls` |

### Profile Selection Flow

```
┌──────────────────────────────────────────────────────────┐
│                    Client Request                        │
└──────────────────────────────────────────────────────────┘
                           │
                           ▼
┌──────────────────────────────────────────────────────────┐
│ 1. Parse User-Agent → Detect Device Group                │
└──────────────────────────────────────────────────────────┘
                           │
                           ▼
┌──────────────────────────────────────────────────────────┐
│ 2. Check Client Capabilities (if reported)               │
│    - Max resolution                                      │
│    - Supported codecs                                    │
│    - HDR support                                         │
└──────────────────────────────────────────────────────────┘
                           │
                           ▼
┌──────────────────────────────────────────────────────────┐
│ 3. Get Base Profile for Device Group                     │
└──────────────────────────────────────────────────────────┘
                           │
                           ▼
┌──────────────────────────────────────────────────────────┐
│ 4. External Client? Check Bandwidth                      │
│    - Measured bandwidth (rolling average)                │
│    - Jitter (variance)                                   │
│    - Override profile if bandwidth too low               │
└──────────────────────────────────────────────────────────┘
                           │
                           ▼
┌──────────────────────────────────────────────────────────┐
│ 5. Apply User Preferences (quality setting)              │
└──────────────────────────────────────────────────────────┘
                           │
                           ▼
┌──────────────────────────────────────────────────────────┐
│ 6. Send Profile ID to Blackbeard                         │
└──────────────────────────────────────────────────────────┘
```

### Bandwidth-Based Override

| Measured Bandwidth | Profile Override |
|--------------------|------------------|
| > 25 Mbps | Use device default |
| 15-25 Mbps | `hevc_1080p` |
| 8-15 Mbps | `h264_1080p` |
| 3-8 Mbps | `h264_720p` |
| 1.5-3 Mbps | `h264_480p` |
| < 1.5 Mbps | `h264_360p` |

---

## Deployment

### Docker Compose (Recommended)

```yaml
version: "3.9"
services:
  revenge:
    image: ghcr.io/lusoris/revenge:latest
    ports: ["8096:8096"]
    environment:
      REVENGE_DATABASE_URL: postgres://revenge:pass@postgres/revenge
      REVENGE_CACHE_URL: redis://dragonfly:6379
      REVENGE_SEARCH_URL: http://typesense:8108
      REVENGE_BLACKBEARD_URL: http://blackbeard:9000
    volumes:
      - ./config:/config
      - /media:/media:ro
    depends_on: [postgres, dragonfly, typesense, blackbeard]

  postgres:
    image: postgres:18-alpine
    environment:
      POSTGRES_USER: revenge
      POSTGRES_PASSWORD: pass
      POSTGRES_DB: revenge
    volumes:
      - postgres_data:/var/lib/postgresql/data

  dragonfly:
    image: docker.dragonflydb.io/dragonflydb/dragonfly

  typesense:
    image: typesense/typesense:0.25.2
    environment:
      TYPESENSE_API_KEY: xyz
      TYPESENSE_DATA_DIR: /data
    volumes:
      - typesense_data:/data

  blackbeard:
    image: ghcr.io/lusoris/blackbeard:latest
    environment:
      BLACKBEARD_REVENGE_URL: http://revenge:8096
    deploy:
      resources:
        reservations:
          devices:
            - capabilities: [gpu]

volumes:
  postgres_data:
  typesense_data:
```

### Platform Support

| Platform | Support |
|----------|---------|
| Linux amd64 | ✅ Full |
| Linux arm64 | ✅ Full |
| macOS arm64 | ✅ Full |
| macOS amd64 | ✅ Full |
| Windows amd64 | ✅ Full |
| FreeBSD | 🔶 Community |

---

## Internal Packages

### pkg/resilience - Fault Tolerance

| Component | Purpose |
|-----------|---------|
| `CircuitBreaker` | Prevent cascade failures to external services |
| `Bulkhead` | Isolate resources with concurrency limits |
| `TokenBucketLimiter` | Rate limiting with burst support |
| `Retry` | Exponential backoff with jitter |

### pkg/supervisor - Self-Healing

| Component | Purpose |
|-----------|---------|
| `Supervisor` | Service supervision with auto-restart |
| `ServiceFunc` | Wrap functions as supervised services |
| Strategies | OneForOne, OneForAll, RestForOne |

### pkg/graceful - Shutdown

| Component | Purpose |
|-----------|---------|
| `Shutdowner` | Ordered shutdown hooks |
| `DrainConnections` | HTTP server connection draining |
| `WaitGroupContext` | Context-aware goroutine tracking |

### pkg/hotreload - Runtime Config

| Component | Purpose |
|-----------|---------|
| `ConfigWatcher` | File change detection with reload |
| `AtomicValue` | Lock-free config access |
| `FeatureFlags` | Runtime feature toggles with percentage rollout |

### pkg/metrics - Observability

| Component | Purpose |
|-----------|---------|
| `Counter`, `Gauge` | Basic metrics |
| `Timer`, `Histogram` | Latency and distribution tracking |
| `HTTPMetrics` | HTTP middleware metrics |

### pkg/lazy - Lazy Initialization

| Component | Purpose |
|-----------|---------|
| `Service[T]` | On-demand initialization |
| `ServiceWithCleanup[T]` | With cleanup on shutdown |
| `Pool[T]` | Round-robin lazy service pool |

### pkg/health - Health Checks

| Component | Purpose |
|-----------|---------|
| `Checker` | Tiered health check management |
| Categories | Critical, Warm, Cold |
| Caching | 5s result cache to prevent overload |

---

## Documentation Index

### Core Architecture

| Document | Description |
|----------|-------------|
| [01_ARCHITECTURE.md](01_ARCHITECTURE.md) | This document - system architecture |
| [TECH_STACK.md](TECH_STACK.md) | Technology choices and rationale |
| [Source of Truth - Project Structure](../00_SOURCE_OF_TRUTH.md#project-structure) | Directory layout |

### Frontend & UI

| Document | Description |
|----------|-------------|
| [FRONTEND.md](FRONTEND.md) | SvelteKit frontend architecture, RBAC, themes |
| [I18N.md](I18N.md) | Internationalization (UI, metadata, audio/subtitle) |

### Content & Metadata

| Document | Description |
|----------|-------------|
| [METADATA_SYSTEM.md](METADATA_SYSTEM.md) | Servarr integration, fallback providers, images |
| [LIBRARY_TYPES.md](LIBRARY_TYPES.md) | Content module definitions |
| [CONTENT_RATING.md](CONTENT_RATING.md) | Content rating systems (MPAA, FSK, etc.) |

### Streaming & Playback

| Document | Description |
|----------|-------------|
| [AUDIO_STREAMING.md](AUDIO_STREAMING.md) | Audio streaming, progress tracking, bandwidth adaptation |
| [CLIENT_SUPPORT.md](CLIENT_SUPPORT.md) | Client capabilities, Chromecast, DLNA |
| [MEDIA_ENHANCEMENTS.md](MEDIA_ENHANCEMENTS.md) | Trailers, themes, intros, trickplay, chapters, Live TV |
| [OFFLOADING.md](OFFLOADING.md) | Blackbeard transcoding integration |
| [SCROBBLING.md](SCROBBLING.md) | External sync (Trakt, Last.fm, ListenBrainz) |

### Operations & Deployment

| Document | Description |
|----------|-------------|
| [REVERSE_PROXY.md](REVERSE_PROXY.md) | Nginx, Caddy, Traefik configuration |
| [BEST_PRACTICES.md](BEST_PRACTICES.md) | Resilience, self-healing, observability |
| [DEVELOPMENT.md](DEVELOPMENT.md) | Development environment setup |
| [SETUP.md](SETUP.md) | Production deployment |
| [Source of Truth - Orchestration](../00_SOURCE_OF_TRUTH.md#container-orchestration) | K8s, K3s, Swarm deployment patterns |
| [Source of Truth - Clustering](../00_SOURCE_OF_TRUTH.md#distributed-consensus-raft) | Raft consensus for HA deployments |

### API & Integration

| Document | Description |
|----------|-------------|
| [API.md](API.md) | API design guidelines |


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

**Category**: [Architecture](INDEX.md)

### In This Section

- [Revenge - Design Principles](02_DESIGN_PRINCIPLES.md)
- [Revenge - Metadata System](03_METADATA_SYSTEM.md)
- [Revenge - Player Architecture](04_PLAYER_ARCHITECTURE.md)
- [Plugin Architecture Decision](05_PLUGIN_ARCHITECTURE_DECISION.md)

### Related Topics

- [Revenge - Adult Content System](../features/adult/ADULT_CONTENT_SYSTEM.md) _Adult_
- [Revenge - Adult Content Metadata System](../features/adult/ADULT_METADATA.md) _Adult_
- [Adult Data Reconciliation](../features/adult/DATA_RECONCILIATION.md) _Adult_
- [Adult Gallery Module (QAR: Treasures)](../features/adult/GALLERY_MODULE.md) _Adult_
- [Whisparr v3 & StashDB Schema Integration](../features/adult/WHISPARR_STASHDB_SCHEMA.md) _Adult_

### Indexes

- [Design Index](../DESIGN_INDEX.md) - All design docs by category/topic
- [Source of Truth](../00_SOURCE_OF_TRUTH.md) - Package versions and status

<!-- DESIGN-BREADCRUMBS-END -->

---

## Cross-References

| Related Document | Relationship |
|------------------|--------------|
| [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md) | Package versions, module lists, config keys |
| [02_DESIGN_PRINCIPLES.md](02_DESIGN_PRINCIPLES.md) | Core design philosophy |
| [03_METADATA_SYSTEM.md](03_METADATA_SYSTEM.md) | Metadata flow and providers |
| [04_PLAYER_ARCHITECTURE.md](04_PLAYER_ARCHITECTURE.md) | WebUI player implementation |
| [05_PLUGIN_ARCHITECTURE_DECISION.md](05_PLUGIN_ARCHITECTURE_DECISION.md) | Plugin vs native decision |
