# Revenge - Development Roadmap

> Modular media server with complete content isolation

**Last Updated**: 2026-01-30
**Current Phase**: Adult Module (QAR) Completion
**Build**: `GOEXPERIMENT=greenteagc,jsonv2 go build ./...`

---

## Quick Status

```
Foundation (Week 1-2)     ████████████████████████ 100%
Design Audit              ████████████████████████ 100%
Critical Fixes            ████████████████████████ 100% ✓
Library Refactor          ████████████████████████ 100% ✓
Movie Module              ████████████████████████  95%
TV Shows Module           ██████████████████████░░  95%
Adult Module (QAR)        ██████████████████████░░  92%  <- CURRENT
Music Module              ░░░░░░░░░░░░░░░░░░░░░░░░   0%
Books Module              ░░░░░░░░░░░░░░░░░░░░░░░░   0%
Comics Module             ░░░░░░░░░░░░░░░░░░░░░░░░   0%
Frontend                  ░░░░░░░░░░░░░░░░░░░░░░░░   0%
```

---

## 🔴 Critical Path (Optimal Execution Order)

> Based on design doc analysis - dependencies mapped, blockers identified

### P0: Access Control + API (BLOCKING)
- [x] **RBAC adult permissions** → `internal/service/rbac/casbin.go`
  - Added adult.* permissions to moderator role in seedDefaultPolicies
- [x] **AdultAuthMiddleware** → `internal/api/handler.go`
  - Added requireAdultAccess, requireAdultBrowse, requireAdultStream helpers
  - Wired RBAC service to Handler for permission checks
- [x] **QAR OpenAPI spec** → `api/openapi/qar.yaml` (1935 lines)
  - Full spec with Fleet, Expedition, Voyage, Crew, Port, Flag, Search endpoints
  - Integration with revenge.yaml + ogen codegen pending
- [ ] **QAR API handlers** → `internal/api/qar.go` (~50 endpoints)

### P1: Enable QAR Workflow
- [ ] **FingerprintService** → `internal/service/fingerprint/`
- [ ] **WhisparrClient** → `internal/service/metadata/whisparr/`

### P2: Quality & Polish
- [ ] **QAR Search isolation** → Typesense collections
- [ ] **StashAppClient** → `internal/service/metadata/stash_app/` (optional)

### P3: Cross-Module (parallelizable)
- [ ] **Playback service** → `internal/service/playback/` (affects movie, tv, qar)
- [ ] **30-day continue watching filter** (movie + tvshow queries)

---

## Phase 1: Critical Fixes (Blocking)

> These must be fixed before any other work

### 1.1 Wiring & Registration Issues
- [x] **Register TVShow module** in `cmd/revenge/main.go`
  - Added `tvshow.ModuleWithRiver` to fx.New()
- [x] **Add TVShow handler deps** to `internal/api/module.go`
  - Added TvshowService to handler dependencies
- [x] **Register Adult modules** in `cmd/revenge/main.go`
  - qar.Module registered (obfuscation complete)

### 1.2 Service Signature Fixes
- [x] **Fix Session.UpdateActivity** in `internal/service/session/service.go`
  - Changed from `(ctx, sessionID, profileID *uuid.UUID)` to `(ctx, sessionID, ipAddress *netip.Addr)`
  - Updated SQL query to track IP address

### 1.3 Configuration Location
- [x] **Config location** - Kept in `pkg/config/` (intentional)
  - Design doc CONFIGURATION.md needs updating to reflect `pkg/config/`

### 1.4 Error Handling
- [x] **Remove os.Exit()** from `cmd/revenge/main.go`
  - Now triggers graceful shutdown via shutdowner.Trigger()
  - Design principle: never panic/exit for errors

### 1.5 Metadata Service Core ✅
- [x] **Implement Radarr provider** in `internal/service/metadata/radarr/`
  - Full API v3 client with circuit breaker
  - types.go, client.go, provider.go, module.go
- [x] **Implement central MetadataService**
  - Orchestration with Servarr-first fallback
  - `internal/service/metadata/service.go` + `module.go`

---

## Phase 2: TV Shows Module Completion

### 2.1 API Handlers
- [x] **Implement TV Shows handlers** in `internal/api/tvshows.go`
  - Integrated tvshows.yaml into revenge.yaml OpenAPI spec
  - All CRUD operations for series, seasons, episodes
  - User data endpoints (favorites, ratings, watch history)
  - Converters added to `internal/api/converters.go`

### 2.2 Watch Next Feature
- [x] **Implement GetNextEpisode handler**
  - Handler in tvshows.go calls TvshowService.GetNextEpisode
- [x] **Implement Continue Watching handler** for TV shows
  - Handler in tvshows.go calls TvshowService.GetContinueWatching

### 2.3 Missing Query Features
- [ ] **Add 30-day filter** to continue watching queries
  - Design: `last_played_at > NOW() - INTERVAL '30 days'`
  - Currently no time window filtering

---

## Phase 2.5: Library Architecture Refactor

> Per LIBRARY_TYPES.md - migrate from shared library table to per-module tables

### 2.5.1 Create Per-Module Library Tables
- [x] **movie_libraries** table in `movie/000005_movie_libraries.up.sql`
  - Movie-specific settings: tmdb_enabled, imdb_enabled, download_trailers, etc.
- [x] **tv_libraries** table in `tvshow/000005_tv_libraries.up.sql`
  - TV-specific settings: sonarr_sync, tvdb_enabled, anime_mode, etc.
- [ ] **music_libraries** table in `music/000001_*.up.sql`
- [ ] **audiobook_libraries** table in `audiobook/000001_*.up.sql`
- [ ] **book_libraries** table in `book/000001_*.up.sql`
- [ ] **podcast_libraries** table in `podcast/000001_*.up.sql`
- [ ] **photo_libraries** table in `photo/000001_*.up.sql`
- [ ] **livetv_sources** table in `livetv/000001_*.up.sql`
- [ ] **comic_libraries** table in `comics/000001_*.up.sql`
- [x] **qar.fleets** table in `qar/000003_qar_obfuscation.up.sql` (adult)

### 2.5.2 Update Foreign Keys
- [x] `movies.library_id` → `movies.movie_library_id REFERENCES movie_libraries(id)`
- [x] `series.library_id` → `series.tv_library_id REFERENCES tv_libraries(id)`
- [x] Update sqlc queries to use new library column names
- [ ] Remove `library_type` enum (no shared enum)

### 2.5.3 Deprecate Shared Library Table
- [x] Add deprecation migration: `shared/000020_deprecate_libraries.up.sql`
- [x] Remove `shared/000005_libraries.up.sql` (fully migrated to per-module tables)

### 2.5.4 Polymorphic Permissions
- [ ] Create `permissions` table with polymorphic `resource_type` + `resource_id`
- [ ] Resource types: `movie_library`, `tv_library`, `qar.fleet`, etc.
- [x] Implement `LibraryProvider` interface in `internal/content/shared/interfaces.go`
- [x] Implement `LibraryService` for movie module

---

## Phase 3: Adult Module Completion

### 3.1 Schema Obfuscation (Queen Anne's Revenge)
- [x] **Rename directories** from `c/` to `qar/` (done)
- [x] **Full entity obfuscation** per ADULT_CONTENT_SYSTEM.md:
  - [x] `qar/expedition/` (Movie → Expedition)
  - [x] `qar/voyage/` (Scene → Voyage)
  - [x] `qar/crew/` (Performer → Crew)
  - [x] `qar/port/` (Studio → Port)
  - [x] `qar/flag/` (Tag → Flag)
  - [x] `qar/fleet/` (Library → Fleet)
- [x] **Update SQL tables** in `qar/000003_qar_obfuscation.up.sql`:
  - [x] `qar.movies` → `qar.expeditions`
  - [x] `qar.scenes` → `qar.voyages`
  - [x] `qar.performers` → `qar.crew`
  - [x] `qar.studios` → `qar.ports`
  - [x] `qar.tags` → `qar.flags`
- [x] **Field obfuscation** (13+ fields) in migration:
  - measurements → cargo, aliases → names, tattoos → markings
  - career_start → maiden_voyage, birth_date → christening
  - penis_size → cutlass, has_breasts → figurehead, etc.
- [ ] **Update API namespace** to `/api/v1/qar/`

### 3.2 Access Control Framework (CRITICAL PATH - BLOCKING)
- [ ] **Add adult permissions** to `internal/service/rbac/casbin.go`:
  - `adult.browse` - view adult content listings
  - `adult.stream` - stream adult content
  - `adult.metadata.write` - edit adult metadata
  - Assign to roles: admin=all, user=opt-in, moderator=all, guest=denied
- [ ] **Implement AdultAuthMiddleware** in `internal/service/auth/middleware_adult.go`:
  - Extract user from context
  - Check RBAC for adult.* permissions
  - Optional PIN verification (stored hashed in user profile)
  - Fire audit event via River (async, fire-and-forget)
- [ ] **Add audit logging** for all adult content access
  - River worker: `AuditAdultAccessWorker`
  - Logs: user_id, resource_type, resource_id, action, timestamp, ip_address
- [ ] **Implement PIN protection** (optional)
  - `user_profiles.adult_pin_hash` column
  - Middleware checks if PIN required + validates

### 3.2.5 QAR API Handlers (CRITICAL PATH - after 3.2)
- [ ] **Create OpenAPI spec** `api/openapi/qar.yaml`:
  - Expedition endpoints (CRUD + user data + crew/flags)
  - Voyage endpoints (CRUD + fingerprint lookup + crew/flags)
  - Crew endpoints (CRUD + names/portraits + expedition/voyage relations)
  - Port endpoints (CRUD + hierarchy)
  - Flag endpoints (CRUD + hierarchy + tagging)
  - Fleet endpoints (CRUD + stats)
- [ ] **Integrate qar.yaml** into `api/openapi/revenge.yaml`
- [ ] **Run ogen codegen** `go generate ./api/...`
- [ ] **Implement handlers** in `internal/api/qar.go`:
  - ~50 endpoints total (copy movie pattern)
  - All require AdultAuthMiddleware
- [ ] **Add converters** in `internal/api/converters.go`:
  - Domain ↔ API types for all QAR entities
- [ ] **Wire handler deps** in `internal/api/module.go`

### 3.3 External Integrations
- [ ] **Implement WhisparrClient** in `internal/service/metadata/whisparr/`:
  - Mirror Radarr client structure (types.go, client.go, provider.go, module.go)
  - API v3: GET /api/v3/movie, /api/v3/person
  - Circuit breaker integration
- [x] **Implement StashDBClient** - GraphQL enrichment (`internal/service/metadata/stashdb/`)
- [ ] **Implement StashAppClient** in `internal/service/metadata/stash_app/`:
  - Sync scene markers as chapters
  - Import user ratings from local Stash instance
  - One-way sync (Stash → Revenge)
- [ ] **Implement FingerprintService** in `internal/service/fingerprint/`:
  - `Fingerprint(path) → VideoFingerprint` (oshash + optional pHash)
  - `MatchScene(fingerprint) → StashDBScene` (query StashDB by hash)
  - Requires ffprobe binary

### 3.4 QAR Modules
- [x] **qar/crew** repository fully implemented (21 methods)
- [x] **qar/port** repository fully implemented (12 methods)
- [x] **qar/flag** repository fully implemented (18 methods)
- [x] **qar/expedition** repository fully implemented (11 methods)
- [x] **qar/voyage** repository fully implemented (13 methods)
- [x] **qar/fleet** repository fully implemented (10 methods)
- [x] **Implement full repository methods** (all complete)

### 3.5 Async Processing
- [x] **Add River jobs** for fingerprinting (expedition/jobs.go, voyage/jobs.go)
- [x] **Add River jobs** for StashDB enrichment (EnrichMetadataWorker)
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
- [x] otter local cache (W-TinyLFU) with config integration
- [x] sturdyc API cache with request coalescing
- [x] Health checks + graceful shutdown
- [x] RBAC with Casbin (dynamic roles)
- [x] Session management + OIDC support
- [x] OpenAPI spec + ogen code generation

### Metadata Infrastructure (2026-01-30)
- [x] Radarr API v3 client (full implementation)
- [x] TMDb provider (existing)
- [x] Central MetadataService with Servarr-first fallback
- [x] Per-module library tables (movie, tv, qar)
- [x] StashDB GraphQL client for QAR metadata enrichment
- [x] LibraryAggregator service with provider interface

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

### TV Shows Module (95%)
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
- [x] API handlers in `internal/api/tvshows.go`
- [ ] 30-day filter for continue watching (missing)

---

## Tech Stack Reference

| Component | Package | Notes |
|-----------|---------|-------|
| Cache (distributed) | `github.com/redis/rueidis` | NOT go-redis |
| Cache (local) | `github.com/maypok86/otter` v1.2.4 | W-TinyLFU |
| Cache (API) | `github.com/viccon/sturdyc` v1.1.5 | Request coalescing |
| Search | `github.com/typesense/typesense-go/v4` | NOT v3 |
| Config | `github.com/knadh/koanf/v2` | NOT viper |
| Logging | `log/slog` | NOT zap |
| Jobs | `github.com/riverqueue/river` | PostgreSQL-native |
| RBAC | `github.com/casbin/casbin/v2` | Dynamic roles |
| DI | `go.uber.org/fx` | Dependency injection |
| HTTP client | `github.com/go-resty/resty/v2` | External APIs |

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
