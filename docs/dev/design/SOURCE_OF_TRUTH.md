# Source of Truth

> Single-page reference for all modules, packages, integrations, and versions

**Last Updated**: 2026-01-30
**Go Version**: 1.25.6
**PostgreSQL**: 18.1
**SQLite**: 3.45+ (via modernc.org/sqlite)

---

## Content Modules

| Module | Schema | Status | Primary Metadata | Arr Integration | Design Doc |
|--------|--------|--------|------------------|-----------------|------------|
| Movie | `public` | ✅ Complete | TMDb | Radarr | [movie.md](features/movies/) |
| TV Show | `public` | ✅ Complete | TMDb, TheTVDB | Sonarr | [tv.md](features/tv/) |
| Music | `public` | 🟡 Scaffold | MusicBrainz, Last.fm | Lidarr | [music.md](features/music/) |
| Audiobook | `public` | 🟡 Scaffold | Audnexus, OpenLibrary | Chaptarr | [audiobook.md](features/audiobook/) |
| Book | `public` | 🟡 Scaffold | OpenLibrary, Goodreads | Chaptarr | [book.md](features/book/) |
| Podcast | `public` | 🟡 Scaffold | RSS/iTunes | Native | [podcast.md](features/podcasts/) |
| Photo | `public` | 🔴 Planned | EXIF, Immich | - | [photo.md](features/photos/) |
| Comics | `public` | 🔴 Planned | ComicVine | - | [comics.md](features/comics/) |
| LiveTV | `public` | 🔴 Planned | XMLTV | TVHeadend | [livetv.md](features/livetv/) |
| QAR Voyages | `qar` | 🟡 Scaffold | StashDB, ThePornDB | Whisparr | [qar.md](features/adult/) |
| QAR Expeditions | `qar` | 🟡 Scaffold | StashDB, ThePornDB | Whisparr | [qar.md](features/adult/) |
| QAR Treasures | `qar` | 🟡 Scaffold | StashDB | Prowlarr | [GALLERY_MODULE.md](features/adult/) |

---

## Backend Services

| Service | Package | fx Module | Status | Design Doc |
|---------|---------|-----------|--------|------------|
| Auth | `internal/service/auth` | `auth.Module` | ✅ Complete | [AUTH.md](services/AUTH.md) |
| User | `internal/service/user` | `user.Module` | ✅ Complete | [USER.md](services/USER.md) |
| Session | `internal/service/session` | `session.Module` | ✅ Complete | [SESSION.md](services/SESSION.md) |
| RBAC | `internal/service/rbac` | `rbac.Module` | ✅ Complete | [RBAC.md](services/RBAC.md) |
| Activity | `internal/service/activity` | `activity.Module` | ✅ Complete | [ACTIVITY.md](services/ACTIVITY.md) |
| Settings | `internal/service/settings` | `settings.Module` | ✅ Complete | [SETTINGS.md](services/SETTINGS.md) |
| API Keys | `internal/service/apikeys` | `apikeys.Module` | ✅ Complete | [APIKEYS.md](services/APIKEYS.md) |
| OIDC | `internal/service/oidc` | `oidc.Module` | ✅ Complete | [OIDC.md](services/OIDC.md) |
| Grants | `internal/service/grants` | `grants.Module` | ✅ Complete | - |
| Fingerprint | `internal/service/fingerprint` | `fingerprint.Module` | ✅ Complete | - |
| Library | `internal/service/library` | `library.Module` | ✅ Complete | [LIBRARY.md](services/LIBRARY.md) |
| Playback | `internal/service/playback` | `playback.Module` | 🟡 Partial | [PLAYBACK.md](technical/) |
| Metadata | `internal/service/metadata` | `metadata.Module` | 🟡 Partial | [METADATA.md](services/METADATA.md) |
| Search | `internal/service/search` | `search.Module` | 🟡 Partial | - |
| Health | `internal/infra/health` | `health.Module` | ✅ Complete | - |
| Scrobbling | `internal/service/scrobbling` | `scrobbling.Module` | 🔴 Planned | [SCROBBLING.md](features/shared/) |
| Analytics | `internal/service/analytics` | `analytics.Module` | 🔴 Planned | - |
| Notification | `internal/service/notification` | `notification.Module` | 🔴 Planned | - |

---

## Infrastructure Components

| Component | Package/Service | Version | Purpose | Design Doc |
|-----------|-----------------|---------|---------|------------|
| PostgreSQL | `pgx/v5` | v5.8.0 | Primary database | [POSTGRESQL.md](integrations/infrastructure/) |
| SQLite | `modernc.org/sqlite` | latest | Embedded database (alt) | [SQLITE.md](integrations/infrastructure/) |
| Dragonfly | `rueidis` | v1.0.71 | Cache/sessions | [DRAGONFLY.md](integrations/infrastructure/) |
| Typesense | `typesense-go/v3` | v3.2.0 | Full-text search | [TYPESENSE.md](integrations/infrastructure/) |
| River | `riverqueue/river` | v0.30.2 | Job queue | [RIVER.md](integrations/infrastructure/) |

### Dual Database Support

| Feature | PostgreSQL | SQLite |
|---------|------------|--------|
| **Driver** | `pgx/v5` | `modernc.org/sqlite` (CGo-free) |
| **Query Gen** | sqlc + pgx | sqlc + database/sql |
| **Migrations** | golang-migrate | golang-migrate |
| **Connection Pool** | pgxpool (self-healing) | Single connection + WAL |
| **Use Case** | Production, multi-user | Single-user, embedded |
| **Admin Switch** | ✅ Runtime switchable | ✅ Runtime switchable |

**Notes:**
- `modernc.org/sqlite` is pure Go (no CGo) - easier cross-compilation
- Performance ~75% of CGo sqlite3, but sufficient for single-user
- WAL mode + NORMAL sync required for acceptable performance
- Database switch available in admin panel + initial setup wizard

---

## Go Dependencies (Core)

| Package | Version | Purpose | Notes |
|---------|---------|---------|-------|
| `go.uber.org/fx` | v1.24.0 | Dependency injection | Lifecycle hooks, graceful shutdown |
| `github.com/jackc/pgx/v5` | v5.8.0 | PostgreSQL driver | Self-healing pool via pgxpool |
| `modernc.org/sqlite` | latest | SQLite driver | CGo-free, ~75% native perf |
| `github.com/riverqueue/river` | v0.30.2 | Job queue | PostgreSQL-backed, transactional |
| `github.com/redis/rueidis` | v1.0.71 | Redis/Dragonfly client | Pipelining, client-side cache |
| `github.com/maypok86/otter` | v1.2.4 | In-memory cache | Faster than Ristretto, S3-FIFO |
| `github.com/knadh/koanf/v2` | v2.3.2 | Configuration | Hot reload via Watch() |
| `github.com/ogen-go/ogen` | v1.18.0 | OpenAPI codegen | Type-safe handlers |
| `github.com/golang-migrate/migrate/v4` | v4.19.1 | DB migrations | Embedded SQL support |
| `github.com/coder/websocket` | v1.8.14 | WebSocket | Low-allocation, concurrent |
| `github.com/go-resty/resty/v2` | v2.17.1 | HTTP client | Retry, backoff, middleware |
| `github.com/google/uuid` | v1.6.0 | UUIDs | V7 time-ordered support |
| `github.com/stretchr/testify` | v1.11.1 | Testing | Assertions, mocks |

## Go Dependencies (Security & RBAC)

| Package | Version | Purpose | Notes |
|---------|---------|---------|-------|
| `github.com/casbin/casbin/v2` | v2.135.0 | RBAC framework | Policy engine |
| `github.com/pckhoi/casbin-pgx-adapter/v3` | v3.2.0 | Casbin PostgreSQL | Async policy sync |
| `golang.org/x/crypto` | v0.47.0 | Cryptography | Argon2, bcrypt, etc. |

## Go Dependencies (Observability)

| Package | Version | Purpose | Notes |
|---------|---------|---------|-------|
| `github.com/lmittmann/tint` | v1.1.2 | Structured logging | Colorized slog handler |
| `go.opentelemetry.io/otel` | v1.39.0 | Telemetry | Traces, metrics |
| `github.com/heptiolabs/healthcheck` | latest | Health probes | K8s liveness/readiness |
| `github.com/go-faster/errors` | v0.7.1 | Error handling | Stack traces, wrapping |

## Go Dependencies (Resilience)

| Package | Version | Purpose | Notes |
|---------|---------|---------|-------|
| `github.com/viccon/sturdyc` | v1.1.5 | Circuit breaker | Bulkhead pattern |
| `github.com/sony/gobreaker` | v1.0.0 | Circuit breaker | Simple, fast |

## Go Dependencies (Serialization)

| Package | Version | Purpose | Notes |
|---------|---------|---------|-------|
| `github.com/go-faster/jx` | v1.2.0 | JSON parsing | Zero-allocation |
| `github.com/go-faster/yaml` | v0.4.6 | YAML parsing | Fast, streaming |

## Go Dependencies (Media Processing) — Planned

| Package | Version | Purpose | Notes |
|---------|---------|---------|-------|
| `github.com/corona10/goimagehash` | latest | Perceptual hashing | pHash, aHash, dHash, wHash |
| `github.com/davidbyttow/govips/v2` | latest | Image processing | libvips bindings (faster than bimg) |
| `github.com/bbrks/go-blurhash` | latest | Blurhash generation | Placeholder images |
| `github.com/asticode/go-astiav` | latest | FFmpeg bindings | Video processing |
| `github.com/dhowden/tag` | latest | Audio metadata | ID3, FLAC, etc. |
| `github.com/mmcdole/gofeed` | latest | RSS/Atom parsing | Podcast feeds |

## Go Dependencies (Development)

| Package | Version | Purpose | Notes |
|---------|---------|---------|-------|
| `github.com/air-verse/air` | latest | Hot reload | File watcher, auto-rebuild |
| `github.com/go-playground/validator/v10` | v10.28.0 | Validation | Struct tags, thread-safe |
| `github.com/sqlc-dev/sqlc` | v1.30.0 | SQL codegen | Type-safe queries |

---

## External Integrations

### Metadata Providers

| Provider | Module(s) | Auth | Rate Limit | Package | Status |
|----------|-----------|------|------------|---------|--------|
| TMDb | Movie, TV | API Key | 50 req/s | `internal/service/metadata/tmdb` | ✅ |
| TheTVDB | TV | API Key | 20 req/10s | `internal/service/metadata/thetvdb` | 🔴 |
| MusicBrainz | Music | None | 1 req/s | `internal/service/metadata/musicbrainz` | 🔴 |
| Last.fm | Music | API Key | 5 req/s | `internal/service/metadata/lastfm` | 🔴 |
| Audnexus | Audiobook | None | Fair use | `internal/service/metadata/audnexus` | 🔴 |
| OpenLibrary | Book | None | Fair use | `internal/service/metadata/openlibrary` | 🔴 |
| ComicVine | Comics | API Key | 200/day | `internal/service/metadata/comicvine` | 🔴 |
| StashDB | QAR | API Key | Fair use | `internal/service/metadata/stashdb` | 🟡 |
| ThePornDB | QAR | API Key | Fair use | `internal/service/metadata/tpdb` | 🔴 |

### Arr Ecosystem

| Service | Content | API Version | Package | Status |
|---------|---------|-------------|---------|--------|
| Radarr | Movies | v3 | `internal/service/metadata/radarr` | ✅ |
| Sonarr | TV | v3 | `internal/service/metadata/sonarr` | 🔴 |
| Lidarr | Music | v1 | `internal/service/metadata/lidarr` | 🔴 |
| Whisparr | QAR | v3 | `internal/service/metadata/whisparr` | 🟡 |
| Chaptarr | Books | v1 (Readarr) | `internal/service/metadata/chaptarr` | 🔴 |
| Prowlarr | Indexers | v1 | `internal/service/metadata/prowlarr` | 🔴 |

### Scrobbling

| Service | Content | Auth | Package | Status |
|---------|---------|------|---------|--------|
| Trakt | Movies, TV | OAuth 2.0 | `internal/service/scrobbling/trakt` | 🔴 |
| Last.fm | Music | API Key + Session | `internal/service/scrobbling/lastfm` | 🔴 |
| ListenBrainz | Music | Token | `internal/service/scrobbling/listenbrainz` | 🔴 |
| Letterboxd | Movies | OAuth 2.0 | `internal/service/scrobbling/letterboxd` | 🔴 |
| Simkl | Movies, TV, Anime | OAuth 2.0 | `internal/service/scrobbling/simkl` | 🔴 |

---

## Database Schemas

| Schema | Purpose | Tables | Access |
|--------|---------|--------|--------|
| `public` | Main content | movies, tv_shows, music_*, etc. | All authenticated |
| `shared` | Shared services | users, sessions, settings, etc. | All authenticated |
| `qar` | Adult content | crew, voyages, expeditions, etc. | `legacy:read` scope |

---

## QAR Obfuscation Terminology

> Queen Anne's Revenge — Pirate-themed naming for adult content isolation

### Module Mapping

| Real Term | QAR Term | Description |
|-----------|----------|-------------|
| Performer | Crew | Adult performers |
| Scene | Voyage | Individual adult scenes |
| Movie | Expedition | Full-length adult films |
| Gallery | Treasure | Image galleries |
| Studio | Port | Production studios |
| Tag | Flag | Content tags/categories |
| Library | Fleet | Content library |

### Fleet Types

| Fleet Type | Content | Arr Integration |
|------------|---------|-----------------|
| `voyage` | Scenes (short clips) | Whisparr |
| `expedition` | Movies (full-length) | Whisparr |
| `treasure` | Galleries (images) | Prowlarr |

### Field Obfuscation (Crew Entity)

| Real Field | QAR Field | Type |
|------------|-----------|------|
| fake_tits | rigged | boolean |
| weight | ballast | int |
| circumcised | trimmed | boolean |
| penis_length | cutlass | string |
| death_date | scuttled | date |
| tattoos | ink | text |
| piercings | hooks | text |
| hair_color | mast_color | string |
| ethnicity | origin | string |
| measurements | cargo | string |

### URL Obfuscation

| Layer | Pattern | Example |
|-------|---------|---------|
| External API | `/api/v1/legacy/*` | `/api/v1/legacy/crew/123` |
| Internal Schema | `qar.*` | `qar.crew WHERE id = '123'` |
| Config Key | `legacy.*` | `legacy.enabled = true` |
| Storage Path | `/data/qar/` | `/data/qar/crew/123/` |

---

## API Namespaces

| Namespace | Purpose | Auth Required | Design Doc |
|-----------|---------|---------------|------------|
| `/api/v1/auth/*` | Authentication | No (public) | [API.md](technical/API.md) |
| `/api/v1/users/*` | User management | Yes | - |
| `/api/v1/movies/*` | Movie content | Yes | - |
| `/api/v1/tv/*` | TV content | Yes | - |
| `/api/v1/music/*` | Music content | Yes | - |
| `/api/v1/audiobooks/*` | Audiobook content | Yes | - |
| `/api/v1/books/*` | Book content | Yes | - |
| `/api/v1/podcasts/*` | Podcast content | Yes | - |
| `/api/v1/photos/*` | Photo content | Yes | - |
| `/api/v1/comics/*` | Comics content | Yes | - |
| `/api/v1/livetv/*` | LiveTV content | Yes | - |
| `/api/v1/legacy/*` | QAR content | Yes + `legacy:read` | [ADULT_CONTENT_SYSTEM.md](features/adult/) |
| `/api/v1/admin/*` | Admin operations | Yes + `admin:*` | - |

---

## Configuration Keys

| Section | Key | Type | Default | Description |
|---------|-----|------|---------|-------------|
| `server` | `port` | int | 8080 | HTTP port |
| `server` | `host` | string | 0.0.0.0 | Bind address |
| `database` | `driver` | string | postgres | `postgres` or `sqlite` |
| `database` | `url` | string | - | PostgreSQL URL |
| `database` | `sqlite_path` | string | - | SQLite file path |
| `cache` | `url` | string | - | Dragonfly URL |
| `search` | `url` | string | - | Typesense URL |
| `search` | `api_key` | string | - | Typesense API key |
| `legacy` | `enabled` | bool | false | Enable QAR module |
| `legacy.privacy` | `require_pin` | bool | true | PIN for QAR access |
| `legacy.privacy` | `audit_all_access` | bool | true | Log all QAR access |

---

## Environment Variable Mapping

> All config keys map to env vars via `REVENGE_` prefix + `_` separator

### Conversion Pattern

```
YAML Key:        server.port
Environment:     REVENGE_SERVER_PORT
Docker Compose:  REVENGE_SERVER_PORT: "8080"
K8s ConfigMap:   REVENGE_SERVER_PORT: "8080"
```

### Core Environment Variables

| Env Variable | Config Key | Type | Required | Example |
|--------------|------------|------|----------|---------|
| `REVENGE_SERVER_PORT` | `server.port` | int | No | `8080` |
| `REVENGE_SERVER_HOST` | `server.host` | string | No | `0.0.0.0` |
| `REVENGE_DATABASE_DRIVER` | `database.driver` | string | Yes | `postgres` |
| `REVENGE_DATABASE_URL` | `database.url` | string | Yes* | `postgres://...` |
| `REVENGE_DATABASE_SQLITE_PATH` | `database.sqlite_path` | string | Yes* | `/data/revenge.db` |
| `REVENGE_CACHE_URL` | `cache.url` | string | No | `redis://dragonfly:6379` |
| `REVENGE_SEARCH_URL` | `search.url` | string | No | `http://typesense:8108` |
| `REVENGE_SEARCH_API_KEY` | `search.api_key` | string | No | `xyz...` |

*One of `DATABASE_URL` or `DATABASE_SQLITE_PATH` required based on driver.

### Secret Environment Variables

| Env Variable | Config Key | Notes |
|--------------|------------|-------|
| `REVENGE_JWT_SECRET` | `auth.jwt_secret` | Min 32 chars |
| `REVENGE_ENCRYPTION_KEY` | `encryption.key` | AES-256 key |
| `REVENGE_LEGACY_ENCRYPTION_KEY` | `legacy.encryption_key` | QAR field encryption |
| `REVENGE_OIDC_CLIENT_SECRET` | `oidc.client_secret` | OIDC provider |

### External Service Variables

| Env Variable | Config Key | Service |
|--------------|------------|---------|
| `REVENGE_TMDB_API_KEY` | `metadata.tmdb.api_key` | TMDb |
| `REVENGE_RADARR_URL` | `arr.radarr.url` | Radarr |
| `REVENGE_RADARR_API_KEY` | `arr.radarr.api_key` | Radarr |
| `REVENGE_SONARR_URL` | `arr.sonarr.url` | Sonarr |
| `REVENGE_SONARR_API_KEY` | `arr.sonarr.api_key` | Sonarr |
| `REVENGE_WHISPARR_URL` | `arr.whisparr.url` | Whisparr |
| `REVENGE_WHISPARR_API_KEY` | `arr.whisparr.api_key` | Whisparr |
| `REVENGE_PROWLARR_URL` | `arr.prowlarr.url` | Prowlarr |
| `REVENGE_PROWLARR_API_KEY` | `arr.prowlarr.api_key` | Prowlarr |
| `REVENGE_STASHDB_API_KEY` | `legacy.stashdb.api_key` | StashDB |

### Docker Compose Example

```yaml
services:
  revenge:
    image: revenge:latest
    environment:
      # Core
      REVENGE_SERVER_PORT: "8080"
      REVENGE_DATABASE_DRIVER: "postgres"
      REVENGE_DATABASE_URL: "postgres://revenge:secret@db:5432/revenge"
      REVENGE_CACHE_URL: "redis://dragonfly:6379"
      REVENGE_SEARCH_URL: "http://typesense:8108"
      REVENGE_SEARCH_API_KEY: "${TYPESENSE_API_KEY}"

      # Secrets (from .env or secrets manager)
      REVENGE_JWT_SECRET: "${JWT_SECRET}"
      REVENGE_ENCRYPTION_KEY: "${ENCRYPTION_KEY}"

      # Metadata
      REVENGE_TMDB_API_KEY: "${TMDB_API_KEY}"

      # Arr Stack
      REVENGE_RADARR_URL: "http://radarr:7878"
      REVENGE_RADARR_API_KEY: "${RADARR_API_KEY}"
```

### Kubernetes ConfigMap/Secret Example

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: revenge-config
data:
  REVENGE_SERVER_PORT: "8080"
  REVENGE_DATABASE_DRIVER: "postgres"
  REVENGE_CACHE_URL: "redis://dragonfly:6379"
  REVENGE_SEARCH_URL: "http://typesense:8108"
---
apiVersion: v1
kind: Secret
metadata:
  name: revenge-secrets
type: Opaque
stringData:
  REVENGE_DATABASE_URL: "postgres://revenge:secret@db:5432/revenge"
  REVENGE_JWT_SECRET: "your-32-char-secret-here"
  REVENGE_TMDB_API_KEY: "your-tmdb-key"
```

### Nested Key Translation

| YAML Path | Environment Variable |
|-----------|---------------------|
| `server.port` | `REVENGE_SERVER_PORT` |
| `database.pool.max_conns` | `REVENGE_DATABASE_POOL_MAX_CONNS` |
| `legacy.privacy.require_pin` | `REVENGE_LEGACY_PRIVACY_REQUIRE_PIN` |
| `arr.radarr.url` | `REVENGE_ARR_RADARR_URL` |
| `metadata.tmdb.api_key` | `REVENGE_METADATA_TMDB_API_KEY` |

### koanf Provider Priority

```
1. Environment variables (REVENGE_*)  ← Highest priority
2. Config file (config.yaml)
3. Default values                     ← Lowest priority
```

Environment variables always override config file values, enabling:
- Container orchestration (K8s, Swarm, Compose)
- CI/CD pipelines
- Secrets management integration

---

## Performance Patterns

### Connection Pool (pgxpool)

| Setting | Default | Production | Notes |
|---------|---------|------------|-------|
| `MaxConns` | 4 | `(CPU * 2) + 1` | Based on CPU cores |
| `MinConns` | 0 | 2 | Keep warm connections |
| `MaxConnLifetime` | 1h | 30m | Prevent stale connections |
| `MaxConnIdleTime` | 30m | 5m | Release idle connections |
| `HealthCheckPeriod` | 1m | 30s | Self-healing checks |

### Self-Healing Features

| Feature | Package | Pattern |
|---------|---------|---------|
| Connection recovery | pgxpool | `Reset()` on network errors |
| Circuit breaker | sturdyc/gobreaker | Open after 5 failures |
| Retry with backoff | resty | Exponential backoff |
| Graceful shutdown | fx/river | Context cancellation |

### Caching Strategy

| Layer | Package | TTL | Purpose |
|-------|---------|-----|---------|
| L1 (In-memory) | otter | 5m | Hot data, zero latency |
| L2 (Distributed) | rueidis | 1h | Shared across instances |
| L3 (Database) | pgx | - | Source of truth |

### GC Tuning (Multi-user)

```bash
GOGC=50          # Frequent short pauses (low latency)
GOMEMLIMIT=2GiB  # Prevent OOM
```

---

## Health Check Patterns

### Kubernetes Probes

| Probe | Endpoint | Timeout | Purpose |
|-------|----------|---------|---------|
| Liveness | `/health/live` | 1s | Is process alive? |
| Readiness | `/health/ready` | 5s | Can accept traffic? |
| Startup | `/health/startup` | 30s | Initial boot complete? |

### Dependency Checks

| Dependency | Check Type | Interval |
|------------|------------|----------|
| PostgreSQL | TCP + Query | 30s |
| Dragonfly | PING | 10s |
| Typesense | /health | 30s |
| River (Jobs) | Worker count | 60s |

### Prometheus Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `health_check_status` | Gauge | 0=unhealthy, 1=healthy |
| `health_check_duration_seconds` | Histogram | Check latency |
| `db_pool_connections` | Gauge | Active connections |
| `cache_hit_ratio` | Gauge | L1/L2 hit rates |

---

## Guided Setup Wizard

### Initial Setup Flow

| Step | Screen | Required | Notes |
|------|--------|----------|-------|
| 1 | Database Selection | Yes | PostgreSQL or SQLite |
| 2 | Admin Account | Yes | Email, password |
| 3 | Library Paths | Yes | Media directories |
| 4 | Cache/Search | No | Optional Dragonfly/Typesense |
| 5 | Arr Integration | No | Radarr, Sonarr, etc. |
| 6 | User Groups | No | Family, guests, etc. |

### Feature Activation Wizards

| Feature | Trigger | Steps |
|---------|---------|-------|
| QAR Module | Admin enables | PIN setup → Library path → Privacy settings |
| Scrobbling | User enables | Select service → OAuth → Sync preferences |
| OIDC | Admin enables | Provider URL → Client ID/Secret → Mapping |
| Notifications | User enables | Select channels → Test → Preferences |

### Per-User Setup (First Login)

| Step | Purpose | Skippable |
|------|---------|-----------|
| Profile | Avatar, display name | Yes |
| Preferences | Language, theme, quality | Yes |
| Libraries | Select visible libraries | No |
| PIN | QAR access (if enabled) | Depends |

---

## Status Legend

| Status | Meaning |
|--------|---------|
| ✅ Complete | Fully implemented and tested |
| 🟡 Scaffold | Structure exists, stubs return "not implemented" |
| 🟡 Partial | Some features implemented |
| 🔴 Planned | Designed but not yet implemented |
| ❌ Deprecated | Scheduled for removal |

---

## Cross-Reference: Module → Integration

| Module | Metadata | Arr | Scrobble | Cloud |
|--------|----------|-----|----------|-------|
| Movie | TMDb | Radarr | Trakt, Letterboxd | - |
| TV | TMDb, TheTVDB | Sonarr | Trakt, Simkl | - |
| Music | MusicBrainz, Last.fm | Lidarr | Last.fm, ListenBrainz | - |
| Audiobook | Audnexus | Chaptarr | - | - |
| Book | OpenLibrary | Chaptarr | Goodreads | - |
| Podcast | iTunes/RSS | - | - | - |
| Photo | EXIF | - | - | Immich, Photoprism |
| Comics | ComicVine | - | - | - |
| LiveTV | XMLTV | - | - | TVHeadend |
| QAR | StashDB, TPDB | Whisparr, Prowlarr | - | Stash-App |

---

## River Job Queue Patterns

### Worker Configuration

| Setting | Default | Production | Notes |
|---------|---------|------------|-------|
| `MaxWorkers` | 100 | CPU * 10 | Per job type |
| `FetchCooldown` | 100ms | 200ms | Reduce DB load |
| `FetchPollInterval` | 1s | 2s | Polling frequency |
| `RescueStuckJobsAfter` | 1h | 30m | Self-healing |

### Graceful Shutdown

| Signal | Action | Timeout |
|--------|--------|---------|
| SIGINT (1st) | Soft stop, finish jobs | 10s |
| SIGINT (2nd) | Hard stop, cancel jobs | 10s |
| SIGINT (3rd) | Force exit | 0s |

### Job Types

| Job | Queue | Priority | Retry |
|-----|-------|----------|-------|
| MetadataEnrich | `metadata` | Normal | 3x exponential |
| LibraryScan | `library` | Low | 1x |
| TranscodeStart | `transcode` | High | 2x |
| ScrobbleSync | `scrobble` | Low | 5x exponential |
| NotificationSend | `notification` | Normal | 3x |

### Transactional Enqueueing

```go
// Jobs enqueued with InsertTx are atomic with other DB changes
tx.InsertTx(ctx, tx, MetadataEnrichJob{MovieID: id})
```

---

## Document Deduplication Policy

> This document is the **single source of truth** for inventory data.
> Other design docs should **reference** this document, not duplicate it.

### What belongs HERE (SOURCE_OF_TRUTH)

- Package names and versions
- Module/service lists with status
- API namespace structure
- Configuration key names
- Schema definitions
- Obfuscation terminology tables

### What belongs in SPECIALIZED DOCS

- Design rationale and philosophy
- Implementation patterns and examples
- Business logic and workflows
- Architecture diagrams
- Feature specifications
- Operations procedures

### Reference Pattern

Other docs should use:
```markdown
See [SOURCE_OF_TRUTH.md](SOURCE_OF_TRUTH.md) for current package versions.
```

Instead of duplicating version tables.

---

## External Sources Cross-Reference

> Live documentation sources for verification and updates

### Source Registry

| Category | Registry Location | Purpose |
|----------|-------------------|---------|
| All Sources | [SOURCES.yaml](../sources/SOURCES.yaml) | Auto-fetch registry |
| Fetch Status | [INDEX.yaml](../sources/INDEX.yaml) | Last fetch timestamps |

### Category → Source Mapping

| This SOT Section | Sources Category | Key Sources |
|------------------|------------------|-------------|
| Go Dependencies | `tooling/` | fx, koanf, ogen, river, rueidis, otter |
| Infrastructure | `infrastructure/` | dragonfly, typesense |
| Database | `database/` | pgx, sqlc, migrations |
| External APIs | `apis/` | tmdb, stashdb, trakt, musicbrainz |
| Frontend | `frontend/` | svelte5, sveltekit, shadcn-svelte |
| Protocols | `protocols/` | hls, dash, http-range |
| Security | `security/` | oidc, oauth2, jwt |

### Version Verification Flow

1. **Weekly CI** fetches live docs to `sources/`
2. **Compare** fetched versions against this SOT
3. **Update** SOT if versions changed
4. **Document** breaking changes in changelogs

### Discrepancies

> Track version mismatches between go.mod, SOURCES.yaml, and this document

See [QUESTIONS_TO_DISCUSS.md](QUESTIONS_TO_DISCUSS.md) for current discrepancies requiring resolution.
