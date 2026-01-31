# Source of Truth

<!-- SOURCES: air, casbin, casbin-pgx-adapter, dragonfly, embedded-postgres, embedded-postgres-docs, ffmpeg, ffmpeg-codecs, ffmpeg-formats, fx, genqlient, genqlient-docs, go-astiav, go-astiav-docs, go-blurhash, go-faster-errors, go-faster-jx, go-faster-yaml, go-fcm, go-fcm-docs, go-io, go-mail, go-mail-docs, gobreaker, gofeed, gofeed-docs, gohlslib, goimagehash, goimagehash-docs, golang-migrate, golang-x-crypto, google-uuid, gorilla-feeds, gorilla-feeds-docs, govips, hashicorp-raft, hashicorp-raft-docs, koanf, lastfm-api, m3u8, mockery, ogen, opa, otter, pgx, postgresql-arrays, postgresql-json, prometheus, prometheus-metrics, raft-boltdb, resty, river, rueidis, rueidis-docs, shadcn-svelte, sqlc, sqlc-config, sturdyc, sturdyc-docs, svelte-runes, svelte5, sveltekit, testcontainers, testify, tint, typesense, typesense-go, validator, xmltv, zap -->

> Single-page reference for all modules, packages, integrations, and versions

**Last Updated**: 2026-01-31
**Go Version**: 1.25.6
**Node.js**: 20.x (LTS)
**Python**: 3.12
**PostgreSQL**: 18.1 (ONLY - no SQLite support)
**Build Command**: `GOEXPERIMENT=greenteagc,jsonv2 go build ./...`

---

## Documentation Map

> Navigate to any part of the documentation from here

| Category | Index | Description |
|----------|-------|-------------|
| **Architecture** | [INDEX](architecture/INDEX.md) | System design, principles, metadata system |
| **Features** | [INDEX](features/INDEX.md) | Content modules, playback, shared features |
| **Integrations** | [INDEX](integrations/INDEX.md) | Metadata providers, Arr stack, auth providers |
| **Services** | [INDEX](services/INDEX.md) | Backend services (auth, user, session, etc.) |
| **Operations** | [INDEX](operations/INDEX.md) | Setup, deployment, best practices |
| **Technical** | [INDEX](technical/INDEX.md) | API, frontend, configuration |
| **Research** | [INDEX](research/INDEX.md) | User pain points, UX/UI resources |

### Quick Links

| Topic | Document |
|-------|----------|
| Full Design Index | [DESIGN_INDEX.md](DESIGN_INDEX.md) |
| Navigation Map | [NAVIGATION.md](NAVIGATION.md) |
| External Sources | [SOURCES_INDEX.md](../sources/SOURCES_INDEX.md) |
| Design ↔ Sources | [DESIGN_CROSSREF.md](../sources/DESIGN_CROSSREF.md) |

### Deep Directory Shortcuts

> Direct access to nested documentation (depth 3+)

| Path | Key Documents |
|------|---------------|
| **Metadata Providers** | |
| └ [integrations/metadata/video/](integrations/metadata/video/INDEX.md) | [TMDb](integrations/metadata/video/TMDB.md), [TheTVDB](integrations/metadata/video/THETVDB.md) |
| └ [integrations/metadata/music/](integrations/metadata/music/INDEX.md) | [MusicBrainz](integrations/metadata/music/MUSICBRAINZ.md), [Last.fm](integrations/metadata/music/LASTFM.md) |
| └ [integrations/metadata/books/](integrations/metadata/books/INDEX.md) | [OpenLibrary](integrations/metadata/books/OPENLIBRARY.md), [Audible](integrations/metadata/books/AUDIBLE.md) |
| └ [integrations/metadata/adult/](integrations/metadata/adult/INDEX.md) | [StashDB](integrations/metadata/adult/STASHDB.md), [ThePornDB](integrations/metadata/adult/THEPORNDB.md) |
| **Wiki Sources** | |
| └ [integrations/wiki/adult/](integrations/wiki/adult/INDEX.md) | [IAFD](integrations/wiki/adult/IAFD.md), [Babepedia](integrations/wiki/adult/BABEPEDIA.md) |

---

## Core Design Principles

### Database Strategy
**PostgreSQL ONLY** - No SQLite support. Simplifies codebase, pgx is the best driver.

### Package Update Policy
**1 Minor Behind** - Use newest STABLE version, never alphas/RCs. Monitor via Dependabot.

### Test Coverage
**80% minimum** - Required for all packages.

### Metadata Priority Chain

> **CORE PRINCIPLE**: If data exists locally, ALWAYS use local first, then fallback to external.

```
Priority Order (ALWAYS):
1. LOCAL CACHE     → First, instant UI display
2. ARR SERVICES    → Radarr, Sonarr, Whisparr (cached metadata)
3. INTERNAL        → Stash-App (if connected)
4. EXTERNAL        → TMDb, StashDB.org, MusicBrainz, etc.
5. ENRICHMENT      → Background jobs, lower priority, seamless
```

This applies to ALL data types across all modules.

### Design Patterns

| Pattern | Decision | Notes |
|---------|----------|-------|
| Error Handling | Sentinels (internal) + Custom APIError (external) | Type-safe errors + API responses |
| Testing | Table-driven + testify + mockery | mockery for auto-generated mocks |
| Logging | Text (Dev, tint) + JSON (Prod, zap) | slog/tint dev, zap prod |
| Metrics | Prometheus + OpenTelemetry | Both - Prometheus for K8s, OTel for traces |
| Validation | ogen (API) + go-playground/validator (Business) | Separate layers |
| Pagination | Cursor (default) + Offset (option) | Cursor for performance, Offset for compatibility |
| Integration Tests | testcontainers-go | Coder-compatible, real containers |
| Unit Tests | embedded-postgres-go | Fast, no containers needed |

---

## Content Modules

| Module | Schema | Status | Primary Metadata | Arr Integration | Design Doc |
|--------|--------|--------|------------------|-----------------|------------|
| Movie | `public` | ✅ Complete | TMDb | Radarr | [MOVIE_MODULE.md](features/video/MOVIE_MODULE.md) |
| TV Show | `public` | ✅ Complete | TMDb, TheTVDB | Sonarr | [TVSHOW_MODULE.md](features/video/TVSHOW_MODULE.md) |
| Music | `public` | 🟡 Scaffold | MusicBrainz, Last.fm | Lidarr | [MUSIC_MODULE.md](features/music/MUSIC_MODULE.md) |
| Audiobook | `public` | 🟡 Scaffold | Audnexus, OpenLibrary | Chaptarr | [AUDIOBOOK_MODULE.md](features/audiobook/AUDIOBOOK_MODULE.md) |
| Book | `public` | 🟡 Scaffold | OpenLibrary, Goodreads | Chaptarr | [BOOK_MODULE.md](features/book/BOOK_MODULE.md) |
| Podcast | `public` | 🟡 Scaffold | RSS/iTunes | Native | [PODCASTS.md](features/podcasts/PODCASTS.md) |
| Photo | `public` | 🔴 Planned | EXIF, Immich | - | [PHOTOS_LIBRARY.md](features/photos/PHOTOS_LIBRARY.md) |
| Comics | `public` | 🔴 Planned | ComicVine | - | [COMICS_MODULE.md](features/comics/COMICS_MODULE.md) |
| LiveTV | `public` | 🔴 Planned | XMLTV | TVHeadend | [LIVE_TV_DVR.md](features/livetv/LIVE_TV_DVR.md) |
| QAR Voyages | `qar` | 🟡 Scaffold | StashDB, ThePornDB | Whisparr | [ADULT_CONTENT_SYSTEM.md](features/adult/ADULT_CONTENT_SYSTEM.md) |
| QAR Expeditions | `qar` | 🟡 Scaffold | StashDB, ThePornDB | Whisparr | [ADULT_CONTENT_SYSTEM.md](features/adult/ADULT_CONTENT_SYSTEM.md) |
| QAR Treasures | `qar` | 🟡 Scaffold | StashDB | Prowlarr | [GALLERY_MODULE.md](features/adult/GALLERY_MODULE.md) |

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
| Grants | `internal/service/grants` | `grants.Module` | 🔵 Planned | [GRANTS.md](services/GRANTS.md) |
| Fingerprint | `internal/service/fingerprint` | `fingerprint.Module` | 🔵 Planned | [FINGERPRINT.md](services/FINGERPRINT.md) |
| Library | `internal/service/library` | `library.Module` | ✅ Complete | [LIBRARY.md](services/LIBRARY.md) |
| Playback | `internal/service/playback` | `playback.Module` | 🟡 Partial | [PLAYBACK.md](technical/) |
| Metadata | `internal/service/metadata` | `metadata.Module` | 🟡 Partial | [METADATA.md](services/METADATA.md) |
| Search | `internal/service/search` | `search.Module` | 🟡 Partial | [SEARCH.md](services/SEARCH.md) |
| Health | `internal/infra/health` | `health.Module` | ✅ Complete | - |
| Scrobbling | `internal/service/scrobbling` | `scrobbling.Module` | 🔴 Planned | [SCROBBLING.md](features/shared/) |
| Analytics | `internal/service/analytics` | `analytics.Module` | 🔴 Planned | [ANALYTICS.md](services/ANALYTICS.md) |
| Notification | `internal/service/notification` | `notification.Module` | 🔴 Planned | [NOTIFICATION.md](services/NOTIFICATION.md) |

---

## Infrastructure Components

| Component | Package/Service | Version | Purpose | Design Doc |
|-----------|-----------------|---------|---------|------------|
| PostgreSQL | `pgx/v5` | v5.7.5 | Primary database | [POSTGRESQL.md](integrations/infrastructure/) |
| Dragonfly | `rueidis` | v1.0.49 | Cache/sessions | [DRAGONFLY.md](integrations/infrastructure/) |
| Typesense | `typesense-go/v4` | v4.x | Full-text search | [TYPESENSE.md](integrations/infrastructure/) |
| River | `riverqueue/river` | v0.26.0 | Job queue | [RIVER.md](integrations/infrastructure/) |

---

## Go Dependencies (Core)

| Package | Version | Purpose | Notes |
|---------|---------|---------|-------|
| `go.uber.org/fx` | v1.23.0 | Dependency injection | Lifecycle hooks, graceful shutdown |
| `github.com/jackc/pgx/v5` | v5.7.5 | PostgreSQL driver | Self-healing pool via pgxpool |
| `github.com/riverqueue/river` | v0.26.0 | Job queue | PostgreSQL-backed, transactional |
| `github.com/redis/rueidis` | v1.0.49 | Redis/Dragonfly client | Pipelining, client-side cache |
| `github.com/maypok86/otter/v2` | v2.x | In-memory cache | W-TinyLFU, faster than Ristretto |
| `github.com/knadh/koanf/v2` | v2.3.0 | Configuration | Hot reload via Watch() |
| `github.com/ogen-go/ogen` | v1.18.0 | OpenAPI codegen | Type-safe handlers |
| `github.com/golang-migrate/migrate/v4` | v4.19.1 | DB migrations | Embedded SQL support |
| `github.com/gobwas/ws` | latest | WebSocket | Zero-alloc, maximum performance |
| `github.com/go-resty/resty/v2` | v2.17.1 | HTTP client | Retry, backoff, middleware |
| `github.com/google/uuid` | v1.6.0 | UUIDs | V7 time-ordered support |
| `github.com/stretchr/testify` | v1.11.1 | Testing | Assertions, mocks |

## Go Dependencies (Security & RBAC)

| Package | Version | Purpose | Notes |
|---------|---------|---------|-------|
| `github.com/casbin/casbin/v2` | v2.135.0 | RBAC framework | Role-based policies |
| `github.com/pckhoi/casbin-pgx-adapter/v3` | v3.2.0 | Casbin PostgreSQL | Async policy sync |
| `github.com/open-policy-agent/opa` | v1.5.0 | Policy engine | Complex ABAC/data-driven policies |
| `golang.org/x/crypto` | v0.47.0 | Cryptography | Argon2, bcrypt, etc. |

## Go Dependencies (Observability)

| Package | Version | Purpose | Notes |
|---------|---------|---------|-------|
| `github.com/lmittmann/tint` | v1.1.2 | Structured logging (dev) | Colorized slog handler |
| `go.uber.org/zap` | v1.28.0 | Structured logging (prod) | High-performance JSON logs |
| `go.opentelemetry.io/otel` | v1.39.0 | Telemetry | Traces, metrics |
| `github.com/heptiolabs/healthcheck` | latest | Health probes | K8s liveness/readiness |
| `github.com/go-faster/errors` | v0.7.1 | Error handling | Stack traces, wrapping |

## Go Dependencies (Resilience)

| Package | Version | Purpose | Notes |
|---------|---------|---------|-------|
| `github.com/viccon/sturdyc` | v1.1.5 | Circuit breaker | Bulkhead pattern |
| `github.com/sony/gobreaker` | v1.0.0 | Circuit breaker | Simple, fast |

## Go Dependencies (Distributed/Clustering)

| Package | Version | Purpose | Notes |
|---------|---------|---------|-------|
| `github.com/hashicorp/raft` | v1.7.x | Consensus | Leader election, state sync |
| `github.com/hashicorp/raft-boltdb/v2` | v2.3.x | Raft storage | BoltDB log store |

## Go Dependencies (Kubernetes)

| Package | Version | Purpose | Notes |
|---------|---------|---------|-------|
| `sigs.k8s.io/controller-runtime` | v0.20.0 | K8s Operator SDK | Operator pattern, reconcilers |
| `k8s.io/client-go` | v0.32.0 | K8s API client | Used by controller-runtime |

## Go Dependencies (Serialization)

| Package | Version | Purpose | Notes |
|---------|---------|---------|-------|
| `github.com/go-faster/jx` | v1.2.0 | JSON parsing | Zero-allocation, ogen-compatible |
| `github.com/go-faster/yaml` | v0.4.6 | YAML parsing | Fast, streaming |
| `github.com/Khan/genqlient` | latest | GraphQL client | Codegen, type-safe queries |

## Go Dependencies (Media Processing) — Planned

| Package | Version | Purpose | Notes |
|---------|---------|---------|-------|
| `github.com/corona10/goimagehash` | latest | Perceptual hashing | pHash, aHash, dHash, wHash |
| `github.com/davidbyttow/govips/v2` | latest | Image processing | libvips bindings (faster than bimg) |
| `github.com/bbrks/go-blurhash` | latest | Blurhash generation | Placeholder images |
| `github.com/asticode/go-astiav` | latest | FFmpeg bindings | Video processing |
| `github.com/wtolson/go-taglib` | latest | Audio metadata | CGo, Read+Write ALL formats |
| `github.com/mmcdole/gofeed` | latest | RSS/Atom parsing | Podcast feeds |

## Go Dependencies (Notifications)

| Package | Version | Purpose | Notes |
|---------|---------|---------|-------|
| `github.com/wneessen/go-mail` | v0.6.2 | SMTP email | Modern, secure, replaces gomail |
| `github.com/appleboy/go-fcm` | v0.2.1 | Firebase Cloud Messaging | Push notifications |
| `github.com/gorilla/feeds` | v1.2.0 | RSS/Atom generation | News feed output |

## Go Dependencies (Development)

| Package | Version | Purpose | Notes |
|---------|---------|---------|-------|
| `github.com/air-verse/air` | latest | Hot reload | File watcher, auto-rebuild |
| `github.com/go-playground/validator/v10` | v10.28.0 | Validation | Struct tags, thread-safe |
| `github.com/sqlc-dev/sqlc` | v1.30.0 | SQL codegen | Type-safe queries |

## Go Dependencies (Testing)

| Package | Version | Purpose | Notes |
|---------|---------|---------|-------|
| `github.com/stretchr/testify` | v1.11.1 | Assertions | Assertions, require, suite |
| `github.com/vektra/mockery/v3` | v3.3.0 | Mock generation | Interface mocking |
| `github.com/testcontainers/testcontainers-go` | v0.37.0 | Integration tests | Real containers |
| `github.com/fergusstrange/embedded-postgres` | v1.30.0 | Unit tests | Fast embedded PostgreSQL |

---

## External Integrations

> **Detailed Docs**: See [integrations/](integrations/) for implementation details, code examples, and API specifics.

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

### Migration File Naming Convention

```
{version}_{description}.{direction}.sql

Examples:
000001_create_users_table.up.sql
000001_create_users_table.down.sql
000002_add_user_email_index.up.sql
000002_add_user_email_index.down.sql
```

| Component | Format | Example |
|-----------|--------|---------|
| Version | 6-digit zero-padded | `000001`, `000042` |
| Description | snake_case | `create_users_table` |
| Direction | `up` or `down` | `up.sql`, `down.sql` |

**Rules:**
- Sequential versioning (no gaps)
- One migration per logical change
- Always provide both `up` and `down`
- Use schema prefix for non-public: `CREATE TABLE qar.crew`

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
| `database` | `driver` | string | postgres | `postgres` (only supported driver) |
| `database` | `url` | string | - | PostgreSQL URL |
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
| `REVENGE_DATABASE_URL` | `database.url` | string | Yes | `postgres://...` |
| `REVENGE_CACHE_URL` | `cache.url` | string | No | `redis://dragonfly:6379` |
| `REVENGE_SEARCH_URL` | `search.url` | string | No | `http://typesense:8108` |
| `REVENGE_SEARCH_API_KEY` | `search.api_key` | string | No | `xyz...` |

*`DATABASE_URL` is required for PostgreSQL connection.

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

## Network QoS Design Principle

> User experience ALWAYS takes priority over background operations

### Priority Hierarchy

| Priority | Traffic Type | Examples | Bandwidth |
|----------|--------------|----------|-----------|
| **P0 (Critical)** | User interface | API responses, WebSocket | Unlimited |
| **P1 (High)** | Active streaming | Video/audio playback | Reserved 80% |
| **P2 (Normal)** | User-initiated | Manual metadata refresh | Fair share |
| **P3 (Low)** | Background jobs | Library scan, thumbnails | Throttled |
| **P4 (Idle)** | Maintenance | Cache cleanup, analytics | Opportunistic |

### Bandwidth Management

```
Available Bandwidth: 100 Mbps
    ↓
P0 + P1 Reserved: 80 Mbps (streaming + UI)
    ↓
P2 + P3 Shared: 20 Mbps (background tasks)
    ↓
When P1 active:
  - P3 throttled to 5 Mbps
  - P4 paused entirely
```

### Implementation Rules

| Rule | Description |
|------|-------------|
| **Stream Priority** | Active playback sessions get priority bandwidth |
| **Job Throttling** | Background jobs pause/slow when user active |
| **Adaptive Quality** | Transcoding adjusts based on available bandwidth |
| **Connection Limits** | External API calls limited per-second |
| **Fair Queuing** | Multiple users share bandwidth fairly |

### River Job QoS

| Queue | Priority | Max Workers (Idle) | Max Workers (Active User) |
|-------|----------|-------------------|---------------------------|
| `playback` | P1 | 50 | 50 (unchanged) |
| `metadata` | P2 | 20 | 10 (reduced) |
| `library` | P3 | 10 | 2 (heavily reduced) |
| `maintenance` | P4 | 5 | 0 (paused) |

### Monitoring Metrics

| Metric | Purpose |
|--------|---------|
| `active_streams_count` | Number of active playback sessions |
| `background_job_throttle_ratio` | How much jobs are throttled (0-1) |
| `bandwidth_utilization_percent` | Current bandwidth usage |
| `job_queue_depth` | Pending jobs per queue |

---

## Distributed Consensus (Raft)

> For multi-instance deployments requiring leader election and state synchronization

### Use Cases

| Feature | Requires Consensus | Reason |
|---------|-------------------|--------|
| Library scan coordination | Yes | Avoid duplicate work |
| Job queue leadership | Yes | Single scheduler |
| Cache invalidation | Yes | Consistent state |
| Session management | No | Dragonfly handles this |
| Playback state | No | Per-user, no conflict |

### Go Packages (Raft)

| Package | Version | Use Case | Notes |
|---------|---------|----------|-------|
| `github.com/hashicorp/raft` | v1.7.x | Batteries-included | Used by Consul, Nomad |
| `go.etcd.io/raft/v3` | v3.6.x | Minimal core | Used by K8s, CockroachDB |

**Recommendation**: Use `hashicorp/raft` for simpler integration (includes storage, transport).

### Cluster Configuration

```yaml
cluster:
  enabled: false              # Single-instance by default
  mode: "raft"                # Consensus mode
  node_id: "node-1"           # Unique node identifier

  raft:
    bind_addr: "0.0.0.0:7000"
    advertise_addr: "192.168.1.10:7000"
    data_dir: "/data/raft"

    # Bootstrap (first node only)
    bootstrap: true

    # Join existing cluster
    join_addrs:
      - "192.168.1.11:7000"
      - "192.168.1.12:7000"

    # Tuning
    heartbeat_timeout: "1s"
    election_timeout: "1s"
    snapshot_interval: "120s"
    snapshot_threshold: 8192
```

### Leader Responsibilities

| Task | Leader Only | Follower Behavior |
|------|-------------|-------------------|
| Schedule library scans | Yes | Forward to leader |
| Process River jobs | All nodes | Distributed workers |
| Serve API requests | All nodes | Proxy writes to leader |
| EPG refresh | Yes | Read from cache |
| Metadata enrichment | All nodes | Idempotent operations |

### Failure Handling

| Scenario | Behavior |
|----------|----------|
| Leader fails | Election within 1-2s, new leader elected |
| Follower fails | Cluster continues, rejoin on recovery |
| Network partition | Minority partition read-only |
| Split brain | Prevented by quorum requirement |

### Minimum Cluster Size

| Nodes | Fault Tolerance | Quorum |
|-------|-----------------|--------|
| 1 | 0 (single instance) | 1 |
| 3 | 1 node failure | 2 |
| 5 | 2 node failures | 3 |

**Note**: For most home server deployments, single-instance mode is sufficient. Enable clustering only for high-availability requirements.

---

## External Client Support

> Integrations for external apps, devices, and systems

### Client Types

| Client Type | Protocol | Features | Example Apps |
|-------------|----------|----------|--------------|
| **Web App** | HTTP/WS | Full UI | Browser |
| **Mobile App** | REST API | Native UX | iOS, Android apps |
| **Smart TV** | REST + HLS | 10-foot UI | LG webOS, Samsung Tizen |
| **IPTV Client** | M3U/XMLTV | Live TV | VLC, Kodi, TiviMate |
| **Media Player** | DLNA/UPnP | Direct play | VLC, MPV |
| **Home Automation** | REST API | Webhooks | Home Assistant |

### Frontend Applets

Embeddable components for external systems:

| Applet | Purpose | Embed Target |
|--------|---------|--------------|
| `now-playing` | Current playback widget | Home Assistant dashboard |
| `recently-added` | New content carousel | Home screen widgets |
| `continue-watching` | Resume playback list | Smart TV apps |
| `live-tv-guide` | EPG mini-guide | Kodi addon |

### Export Endpoints

| Export | Format | Endpoint | Auth |
|--------|--------|----------|------|
| IPTV Channels | M3U8 | `/api/v1/livetv/export/m3u` | API Token |
| EPG Guide | XMLTV | `/api/v1/livetv/export/epg` | API Token |
| Calendar | iCal | `/api/v1/calendar/export/ical` | API Token |
| Watchlist | JSON | `/api/v1/users/me/watchlist/export` | Bearer |

### Kodi Integration

| Addon | Purpose | Protocol |
|-------|---------|----------|
| `plugin.video.revenge` | Browse & play | REST API + HLS |
| `service.revenge.sync` | Watch state sync | WebSocket |
| `pvr.revenge.livetv` | Live TV | IPTV Simple Client |

### Home Assistant Integration

```yaml
# configuration.yaml
media_player:
  - platform: revenge
    host: revenge.local
    port: 8080
    api_key: !secret revenge_api_key

sensor:
  - platform: revenge
    host: revenge.local
    resources:
      - now_playing
      - recently_added
      - library_stats
```

### Webhook Events

| Event | Payload | Use Case |
|-------|---------|----------|
| `playback.started` | Media info, user | Lighting automation |
| `playback.stopped` | Duration watched | Scrobble trigger |
| `library.new_item` | Item metadata | Notification |
| `transcode.complete` | Job result | Monitoring |

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

## Container Orchestration

> Deployment patterns for Kubernetes, K3s, and Docker Swarm

### Deployment Modes

| Mode | Use Case | Scaling | State |
|------|----------|---------|-------|
| **Single Instance** | Home server, NAS | None | Local volumes |
| **Docker Compose** | Small deployments | Manual | Named volumes |
| **Docker Swarm** | Multi-node home lab | Replicas | Shared storage |
| **K3s** | Lightweight K8s | HPA | PVC |
| **Kubernetes** | Production, HA | HPA + VPA | PVC + CSI |

### Architecture Requirements

| Component | Stateless | Scalable | Notes |
|-----------|-----------|----------|-------|
| Revenge API | Yes | Yes | Horizontal scaling OK |
| River Workers | Yes | Yes | Distributed job processing |
| PostgreSQL | No | Read replicas | Use managed DB or StatefulSet |
| Dragonfly | No | Cluster mode | Use managed Redis or StatefulSet |
| Typesense | No | Cluster mode | StatefulSet with PVC |
| ErsatzTV | No | Single | Per-channel instance possible |

### Kubernetes Resources

```yaml
# Deployment pattern for stateless Revenge pods
apiVersion: apps/v1
kind: Deployment
metadata:
  name: revenge
spec:
  replicas: 3
  selector:
    matchLabels:
      app: revenge
  template:
    spec:
      containers:
      - name: revenge
        image: revenge:latest
        ports:
        - containerPort: 8080
        envFrom:
        - configMapRef:
            name: revenge-config
        - secretRef:
            name: revenge-secrets
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 5
        resources:
          requests:
            memory: "256Mi"
            cpu: "100m"
          limits:
            memory: "1Gi"
            cpu: "1000m"
        volumeMounts:
        - name: media
          mountPath: /media
          readOnly: true
      volumes:
      - name: media
        persistentVolumeClaim:
          claimName: media-pvc
```

### Helm Chart Structure

```
charts/revenge/
├── Chart.yaml
├── values.yaml
├── templates/
│   ├── deployment.yaml
│   ├── service.yaml
│   ├── ingress.yaml
│   ├── configmap.yaml
│   ├── secret.yaml
│   ├── hpa.yaml
│   ├── pvc.yaml
│   └── _helpers.tpl
└── charts/
    ├── postgresql/          # Bitnami subchart
    ├── dragonfly/           # Optional subchart
    └── typesense/           # Optional subchart
```

### Docker Swarm Stack

```yaml
version: "3.8"
services:
  revenge:
    image: revenge:latest
    deploy:
      replicas: 2
      update_config:
        parallelism: 1
        delay: 10s
      restart_policy:
        condition: on-failure
      resources:
        limits:
          memory: 1G
        reservations:
          memory: 256M
    environment:
      REVENGE_DATABASE_URL: "postgres://revenge:secret@db:5432/revenge"
      REVENGE_CACHE_URL: "redis://dragonfly:6379"
    volumes:
      - media:/media:ro
    networks:
      - revenge-net
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health/ready"]
      interval: 30s
      timeout: 10s
      retries: 3

  db:
    image: postgres:18
    deploy:
      placement:
        constraints:
          - node.role == manager
    volumes:
      - pgdata:/var/lib/postgresql/data
    networks:
      - revenge-net

  dragonfly:
    image: docker.dragonflydb.io/dragonflydb/dragonfly:latest
    deploy:
      placement:
        constraints:
          - node.role == manager
    volumes:
      - dfdata:/data
    networks:
      - revenge-net

volumes:
  media:
    driver: local
    driver_opts:
      type: nfs
      o: addr=nas.local,rw
      device: ":/volume1/media"
  pgdata:
  dfdata:

networks:
  revenge-net:
    driver: overlay
```

### K3s Specifics

| Feature | K3s Approach | Notes |
|---------|--------------|-------|
| **Ingress** | Traefik (default) | Auto-provisioned |
| **Storage** | Local-path (default) | Or Longhorn for HA |
| **Load Balancer** | ServiceLB (default) | Or MetalLB |
| **Certificates** | cert-manager | Let's Encrypt |

### Scaling Strategy

| Component | Min | Max | Trigger |
|-----------|-----|-----|---------|
| Revenge API | 2 | 10 | CPU > 70% |
| River Workers | 1 | 5 | Queue depth > 100 |
| Transcoding | 0 | 3 | Active transcode jobs |

### Persistent Storage

| Data Type | Storage Class | Access Mode | Backup |
|-----------|---------------|-------------|--------|
| Database | SSD/NVMe | RWO | pg_dump daily |
| Cache | SSD | RWO | Not required |
| Search Index | SSD | RWO | Rebuild from DB |
| Media Files | HDD/NAS | ROX | External backup |
| Config | Any | RWO | GitOps |

### Service Discovery

| Platform | Method | Service Name |
|----------|--------|--------------|
| Docker Compose | DNS | `revenge`, `db`, `dragonfly` |
| Docker Swarm | DNS + VIP | `revenge_revenge`, `revenge_db` |
| Kubernetes | CoreDNS | `revenge.default.svc.cluster.local` |

### Secrets Management

| Platform | Method | Notes |
|----------|--------|-------|
| Docker Compose | `.env` file | Not for production |
| Docker Swarm | Docker Secrets | `docker secret create` |
| Kubernetes | Secrets + Sealed Secrets | Encrypt at rest |
| External | Vault, AWS Secrets Manager | Enterprise |

### Init Containers (K8s)

```yaml
initContainers:
- name: wait-for-db
  image: busybox
  command: ['sh', '-c', 'until nc -z db 5432; do sleep 2; done']
- name: run-migrations
  image: revenge:latest
  command: ['revenge', 'migrate', 'up']
  envFrom:
  - secretRef:
      name: revenge-secrets
```

### Resource Recommendations

| Deployment | CPU Request | Memory Request | CPU Limit | Memory Limit |
|------------|-------------|----------------|-----------|--------------|
| Minimal | 100m | 256Mi | 500m | 512Mi |
| Standard | 250m | 512Mi | 1000m | 1Gi |
| Production | 500m | 1Gi | 2000m | 2Gi |
| + Transcoding | 2000m | 2Gi | 4000m | 4Gi |

---

## Guided Setup Wizard

### Initial Setup Flow

| Step | Screen | Required | Notes |
|------|--------|----------|-------|
| 1 | Database Connection | Yes | PostgreSQL URL |
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

## Status System

### Status Values

| Emoji | Meaning |
|-------|---------|
| ✅ | Complete |
| 🟡 | Partial |
| 🔴 | Not Started |
| ⚪ | N/A |

### Status Dimensions

Each module/service/integration tracks progress across 7 dimensions:

| Dimension | Description |
|-----------|-------------|
| **Design** | Feature spec, architecture, DB schema, API endpoints documented |
| **Sources** | External docs fetched (API specs, GraphQL schemas, best practices) |
| **Instructions** | Claude Code instructions for implementation |
| **Code** | Go implementation |
| **Linting** | golangci-lint rules, formatting |
| **Unit Testing** | Unit tests with embedded-postgres |
| **Integration Testing** | Integration tests with testcontainers |

### Status Format in Design Docs

```markdown
| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | |
| Sources | 🟡 | API docs fetched, GraphQL schema missing |
| Instructions | 🔴 | |
| Code | 🔴 | Reset to template |
| Linting | 🔴 | |
| Unit Testing | 🔴 | |
| Integration Testing | 🔴 | |
```

### Current Reality (2026-01-30)

> **Code was reset to Go template** - All code status is 🔴 NOT STARTED
> Design docs are the focus of M1 milestone

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

## Project Structure

### Backend (Go)

```
revenge/
├── cmd/revenge/          # Main entry point
├── internal/
│   ├── api/              # HTTP handlers (ogen-generated)
│   ├── config/           # Configuration loading
│   ├── content/          # Content modules
│   │   ├── movie/        # Movie module
│   │   ├── tvshow/       # TV Show module
│   │   ├── music/        # Music module (scaffold)
│   │   ├── qar/          # QAR adult module
│   │   │   ├── expedition/
│   │   │   ├── voyage/
│   │   │   ├── crew/
│   │   │   ├── port/
│   │   │   ├── flag/
│   │   │   └── fleet/
│   │   └── ...
│   ├── service/          # Business services
│   │   ├── auth/
│   │   ├── user/
│   │   ├── session/
│   │   ├── rbac/
│   │   ├── metadata/     # Metadata providers
│   │   │   ├── tmdb/
│   │   │   ├── radarr/
│   │   │   ├── stashdb/
│   │   │   └── whisparr/
│   │   └── ...
│   └── infra/            # Infrastructure
│       ├── database/     # pgx, sqlc, migrations
│       ├── cache/        # rueidis, otter, sturdyc
│       ├── search/       # Typesense client
│       ├── jobs/         # River queue
│       └── health/       # Health checks
├── api/openapi/          # OpenAPI specs
├── pkg/                  # Public packages
└── docs/                 # Documentation
```

### Module Pattern

Each content module follows:
```
internal/content/{module}/
├── entity.go           # Domain types
├── repository.go       # Repository interface
├── repository_pg.go    # PostgreSQL implementation
├── service.go          # Business logic
├── jobs.go             # River job definitions
├── library_service.go  # LibraryProvider impl
└── module.go           # fx module
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
See [00_SOURCE_OF_TRUTH.md](00_SOURCE_OF_TRUTH.md) for current package versions.
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
| Sources Index | [SOURCES_INDEX.md](../sources/SOURCES_INDEX.md) | Browsable source list |
| Design ↔ Sources | [DESIGN_CROSSREF.md](../sources/DESIGN_CROSSREF.md) | Doc-to-source mapping |

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

See [02_QUESTIONS_TO_DISCUSS.md](02_QUESTIONS_TO_DISCUSS.md) for current discrepancies requiring resolution.
