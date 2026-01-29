# Revenge - Development Roadmap

> Modular media server with complete content isolation

**Last Updated**: 2026-01-29
**Current Phase**: Design Audit Fixes
**Build**: `GOEXPERIMENT=greenteagc,jsonv2 go build ./...`

---

## Quick Status

```
Foundation (Week 1-2)     ████████████████████████ 100%
Design Audit              ████████████████████████ 100%
Critical Fixes            ████████░░░░░░░░░░░░░░░░  30%  <- CURRENT
Library Refactor          ░░░░░░░░░░░░░░░░░░░░░░░░   0%  <- BLOCKING
Movie Module              ████████████████████░░░░  85%
TV Shows Module           ████████████████████░░░░  80%
Adult Module (QAR)        ████████░░░░░░░░░░░░░░░░  35%
Music Module              ░░░░░░░░░░░░░░░░░░░░░░░░   0%
Books Module              ░░░░░░░░░░░░░░░░░░░░░░░░   0%
Comics Module             ░░░░░░░░░░░░░░░░░░░░░░░░   0%
Frontend                  ░░░░░░░░░░░░░░░░░░░░░░░░   0%
```

---

## Phase 1: Critical Fixes (Blocking)

> These must be fixed before any other work

### 1.1 Wiring & Registration Issues
- [x] **Register TVShow module** in `cmd/revenge/main.go`
  - Added `tvshow.ModuleWithRiver` to fx.New()
- [x] **Add TVShow handler deps** to `internal/api/module.go`
  - Added TvshowService to handler dependencies
- [ ] **Register Adult modules** in `cmd/revenge/main.go`
  - ⚠️ Blocked: Requires full QAR obfuscation first (Phase 3)

### 1.2 Service Signature Fixes
- [ ] **Fix Session.UpdateActivity** in `internal/service/session/service.go`
  - Design: `UpdateActivity(ctx, sessionID, ipAddress *netip.Addr)`
  - Code: `UpdateActivity(ctx, sessionID, profileID *uuid.UUID)`
  - Update signature + all callers

### 1.3 Configuration Location
- [ ] **Move config** from `pkg/config/` to `internal/config/`
  - Update all imports
  - Design specifies `internal/config/config.go`

### 1.4 Error Handling
- [x] **Remove os.Exit()** from `cmd/revenge/main.go`
  - Now triggers graceful shutdown via shutdowner.Trigger()
  - Design principle: never panic/exit for errors

### 1.5 Metadata Service Core
- [ ] **Implement Radarr provider** in `internal/service/metadata/radarr/`
  - Currently stub returning ErrUnavailable
  - Must work per Servarr-first principle (Priority 1)
- [ ] **Implement central MetadataService**
  - Create aggregation service with provider priority/fallback
  - Design: `internal/service/metadata/service.go`

---

## Phase 2: TV Shows Module Completion

### 2.1 API Handlers
- [ ] **Implement TV Shows handlers** in `internal/api/tvshows.go`
  - Wire to OpenAPI spec `api/openapi/tvshows.yaml`
  - All CRUD operations for series, seasons, episodes
  - User data endpoints (favorites, ratings, watch history)

### 2.2 Watch Next Feature
- [ ] **Implement GetNextEpisode handler**
  - Service exists but no API handler
  - Endpoint defined in OpenAPI spec
- [ ] **Implement Continue Watching handler** for TV shows
  - Movies done, TV shows missing

### 2.3 Missing Query Features
- [ ] **Add 30-day filter** to continue watching queries
  - Design: `last_played_at > NOW() - INTERVAL '30 days'`
  - Currently no time window filtering

---

## Phase 2.5: Library Architecture Refactor

> Per LIBRARY_TYPES.md - migrate from shared library table to per-module tables

### 2.5.1 Create Per-Module Library Tables
- [ ] **movie_libraries** table in `movie/000005_movie_libraries.up.sql`
  - Movie-specific settings: tmdb_enabled, imdb_enabled, download_trailers, etc.
- [ ] **tv_libraries** table in `tvshow/000005_tv_libraries.up.sql`
  - TV-specific settings: sonarr_sync, tvdb_enabled, anime_mode, etc.
- [ ] **music_libraries** table in `music/000001_*.up.sql`
- [ ] **audiobook_libraries** table in `audiobook/000001_*.up.sql`
- [ ] **book_libraries** table in `book/000001_*.up.sql`
- [ ] **podcast_libraries** table in `podcast/000001_*.up.sql`
- [ ] **photo_libraries** table in `photo/000001_*.up.sql`
- [ ] **livetv_sources** table in `livetv/000001_*.up.sql`
- [ ] **comic_libraries** table in `comics/000001_*.up.sql`
- [ ] **qar.fleets** table in `qar/000001_fleets.up.sql` (adult)

### 2.5.2 Update Foreign Keys
- [ ] `movies.library_id` → `REFERENCES movie_libraries(id)`
- [ ] `series.library_id` → `REFERENCES tv_libraries(id)`
- [ ] Remove `library_type` enum (no shared enum)

### 2.5.3 Deprecate Shared Library Table
- [ ] Add data migration: `shared/000020_deprecate_libraries.up.sql`
- [ ] Eventually remove `shared/000005_libraries.up.sql`

### 2.5.4 Polymorphic Permissions
- [ ] Create `permissions` table with polymorphic `resource_type` + `resource_id`
- [ ] Resource types: `movie_library`, `tv_library`, `qar.fleet`, etc.
- [ ] Implement `LibraryProvider` interface in each module

---

## Phase 3: Adult Module Completion

### 3.1 Schema Obfuscation (Queen Anne's Revenge)
- [x] **Rename directories** from `c/` to `qar/` (done)
- [ ] **Full entity obfuscation** per ADULT_CONTENT_SYSTEM.md:
  - `qar/movie/` → `qar/expedition/` (Movie → Expedition)
  - `qar/scene/` → `qar/voyage/` (Scene → Voyage)
  - Create `qar/crew/` (Performer → Crew)
  - Create `qar/port/` (Studio → Port)
  - Create `qar/flag/` (Tag → Flag)
  - Create `qar/fleet/` (Library → Fleet)
- [ ] **Update SQL tables** to obfuscated names:
  - `qar.movies` → `qar.expeditions`
  - `qar.scenes` → `qar.voyages`
  - `qar.performers` → `qar.crew`
  - `qar.studios` → `qar.ports`
  - `qar.tags` → `qar.flags`
- [ ] **Field obfuscation** (13+ fields):
  - measurements → cargo, aliases → names, tattoos → markings
  - career_start → maiden_voyage, birth_date → christening
  - penis_size → cutlass, has_breasts → figurehead, etc.
- [ ] **Update API namespace** to `/api/v1/qar/`

### 3.2 Access Control Framework
- [ ] **Add scoped permissions** (adult:read, adult:write)
- [ ] **Implement AdultAuthMiddleware**
- [ ] **Add audit logging** for all adult content access
- [ ] **Implement PIN protection** (optional)

### 3.3 External Integrations
- [ ] **Implement WhisparrClient** - acquisition proxy
- [ ] **Implement StashDBClient** - GraphQL enrichment
- [ ] **Implement StashAppClient** - private instance sync
- [ ] **Implement FingerprintService** - hash generation/matching

### 3.4 Missing Modules
- [ ] **Implement c/performer** module (service, repository)
- [ ] **Implement c/studio** module (service, repository)
- [ ] **Implement c/tag** module (service, repository)
- [ ] **Implement c.show** module (directory exists but empty)

### 3.5 Async Processing
- [ ] **Add River jobs** for fingerprinting
- [ ] **Add River jobs** for StashDB enrichment
- [ ] **Add River jobs** for Stash-App sync

### 3.6 Search Isolation
- [ ] **Create separate Typesense collections** for adult content (qar_movies, qar_voyages)
- [ ] **Separate search endpoint** `/api/v1/qar/search` requiring adult:read scope

### 3.7 Update Instructions File
- [ ] **Update adult-modules.instructions.md** to use `qar` schema (currently says `c`)
  - Instructions file is outdated, design doc is truth

---

## Phase 4: RBAC & Security Completion

### 4.1 Missing RBAC Methods
- [ ] **Add missing methods** to `internal/service/rbac/casbin.go`
  - Enforce(), AddRoleForUser(), RemoveRoleForUser()
  - GetRolesForUser(), GetUsersForRole(), etc.

### 4.2 Resource Grants
- [ ] **Create resource_grants table** (polymorphic permissions)
  - Migration: `000019_resource_grants.up.sql`
- [ ] **Implement HasGrant(), CreateGrant(), DeleteByResource()**

### 4.3 Metadata Audit Logging
- [ ] **Redesign activity_log schema**
  - Add: module, entity_id, entity_type, changes (JSONB)
  - Implement partitioning by month
- [ ] **Add River async worker** for audit writes
- [ ] **Implement edit history** with rollback capability

### 4.4 Missing Permissions
- [ ] **Add access.* permissions** to casbin.go
  - access.rules.view, access.rules.manage, access.bypass
- [ ] **Seed request permissions** to database (15 missing)
- [ ] **Seed adult request permissions** (7 missing)

### 4.5 Activity Encryption
- [ ] **Implement AES-256-GCM** for activity data at rest
  - Design principle: Privacy by Default

---

## Phase 5: Playback Features

### 5.1 Up Next / Auto-Play Queue
- [ ] **Create UpNextQueue struct**
- [ ] **Implement BuildUpNextQueue()** service
- [ ] **Add /api/playback/up-next endpoint**

### 5.2 Cross-Device Sync
- [ ] **Implement WebSocket** for playback sync
- [ ] **Implement polling fallback** `/api/sync/playback?since={ts}`
- [ ] **Add BroadcastToUser()** on position updates

### 5.3 User Preferences
- [ ] **Add user preference fields** to database
  - auto_play_enabled, auto_play_delay_seconds
  - continue_watching_days, mark_watched_percent

---

## Phase 6: Tech Stack Alignment

### 6.1 Missing Libraries
- [ ] **Add failsafe-go** OR document custom `pkg/resilience/`
- [ ] **Add golang-migrate** OR remove from design doc
- [ ] **Add coder/websocket** for WebSocket support
- [ ] **Add govips** for image processing
- [ ] **Add dhowden/tag** for audio metadata

### 6.2 Documentation Updates
- [ ] **Add casbin** to TECH_STACK.md (actively used but undocumented)
- [ ] **Add OpenTelemetry** to TECH_STACK.md (used by ogen)
- [ ] **Document typesense alpha** status or upgrade to stable

### 6.3 Health Checks
- [ ] **Enable cache health check** in main.go (commented out)
- [ ] **Enable search health check** in main.go (commented out)

### 6.4 Advanced Patterns (per instructions)
- [ ] **Circuit breakers** for external services (TMDb, Radarr, Sonarr)
  - Use `pkg/resilience/` circuit breaker for all external API calls
- [ ] **Lazy initialization** for non-critical services
  - Transcoder client, metadata providers, search client
- [ ] **Rate limiting** at API boundaries
  - Per-user and per-IP rate limits using `pkg/resilience/`
- [ ] **Config hot reload** for runtime settings
  - Feature flags, log levels, rate limits via `pkg/hotreload/`
- [ ] **Update adult-modules.instructions.md**
  - Currently says schema `c`, design doc specifies `qar`

---

## Phase 7: Missing Content Modules

| Module | Priority | Notes |
|--------|----------|-------|
| Music | P2 | Lidarr + MusicBrainz |
| Audiobooks | P2 | Chaptarr + Audible |
| Books | P2 | Chaptarr + OpenLibrary |
| Podcasts | P3 | RSS feeds |
| Comics | P3 | ComicVine |
| Photos | P3 | EXIF, GPS metadata |
| LiveTV | P3 | TVHeadend/NextPVR |
| Collection | P3 | Cross-module pools |

---

## Completed

### Foundation (Week 1-2)
- [x] PostgreSQL + sqlc type-safe queries
- [x] Dragonfly cache (rueidis client)
- [x] Typesense search (typesense-go/v4)
- [x] River job queue (PostgreSQL-native)
- [x] uber-go/fx dependency injection
- [x] koanf v2 configuration
- [x] slog structured logging
- [x] otter local cache (W-TinyLFU)
- [x] Health checks + graceful shutdown
- [x] RBAC with Casbin (dynamic roles)
- [x] Session management + OIDC support
- [x] OpenAPI spec + ogen code generation

### Design Audit (2026-01-29)
- [x] Audit all services against design docs
- [x] Audit architecture docs
- [x] Audit technical docs
- [x] Audit feature docs (playback, RBAC, adult, access)
- [x] Create DESIGN_AUDIT_REPORT.md with 67 issues

### Movie Module (85%)
- [x] Full CRUD with relations
- [x] User data (ratings, favorites, watchlist)
- [x] TMDb metadata provider
- [x] River jobs for metadata enrichment
- [ ] Continue watching 30-day filter (missing)

### TV Shows Module (80%)
- [x] Database migrations
- [x] sqlc queries (100+ queries)
- [x] Entity definitions
- [x] Repository (PostgreSQL)
- [x] Service layer
- [x] module.go (fx registration)
- [x] jobs.go (River workers)
- [x] metadata_provider.go
- [x] Module registration in main.go
- [x] Handler deps in api/module.go
- [ ] API handlers (missing)

---

## Tech Stack Reference

| Component | Package | Notes |
|-----------|---------|-------|
| Cache (distributed) | `github.com/redis/rueidis` | NOT go-redis |
| Cache (local) | `github.com/maypok86/otter` v1.2.4 | W-TinyLFU |
| Search | `github.com/typesense/typesense-go/v4` | NOT v3 |
| Config | `github.com/knadh/koanf/v2` | NOT viper |
| Logging | `log/slog` | NOT zap |
| Jobs | `github.com/riverqueue/river` | PostgreSQL-native |
| RBAC | `github.com/casbin/casbin/v2` | Dynamic roles |
| DI | `go.uber.org/fx` | Dependency injection |

---

## Important Notes

**Adult Content** (Queen Anne's Revenge obfuscation):
- Schema: `qar` (isolated PostgreSQL schema)
- API namespace: `/qar/*`
- Module location: `internal/content/qar/`
- See [ADULT_CONTENT_SYSTEM.md](docs/dev/design/features/adult/ADULT_CONTENT_SYSTEM.md)

**Design Docs are Source of Truth**:
- Only `docs/dev/design/` is authoritative
- Other documentation may be outdated
- Code must match design, not vice versa

**Build Commands**:
```bash
# With experiments
GOEXPERIMENT=greenteagc,jsonv2 go build -o bin/revenge ./cmd/revenge

# Generate code
sqlc generate
go generate ./api/...

# Lint
golangci-lint run
```
