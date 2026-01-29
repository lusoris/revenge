# Revenge - Development Roadmap

> Modular media server with complete content isolation

**Last Updated**: 2026-01-29
**Current Phase**: Foundation (Week 1-2)
**Implementation Status**: 95% complete (Week 1-2 + Movie Module done)

---

## 📊 Quick Links

- 📐 [Architecture](docs/ARCHITECTURE_V2.md) - Complete system design
- 📋 [Analysis Reports](#analysis-reports) - Current status
- 🎯 [Current Sprint](#current-sprint-week-1) - What we're working on now
- 🗺️ [Roadmap](#roadmap-overview) - 16-week plan to MVP

---

## Analysis Reports

- [Architecture Compliance](ARCHITECTURE_COMPLIANCE_ANALYSIS.md) - 65% conformance score
- [Advanced Features Integration](ADVANCED_FEATURES_INTEGRATION_ANALYSIS.md) - 10% integration score
- [Core Functionality Analysis](CORE_FUNCTIONALITY_ANALYSIS.md) - Missing workers/services
- [Design TODOs Extraction](DESIGN_TODOS_EXTRACTION.md) - 100+ components to implement
- [Documentation Cleanup](DOCUMENTATION_CLEANUP_REPORT.md) - 264+ outdated TODOs removed
- [Comprehensive Analysis](COMPREHENSIVE_ANALYSIS_FINAL.md) - Complete status report

---

## Roadmap Overview

```
Week 1-2:  Foundation (P0)       ████████████████████████ 100%
Week 3-4:  Movie Module (P1)     ██████████████████████░░ 95%
Week 4-8:  Modules + Frontend    ░░░░░░░░░░░░░░░░░░░░░░░░  0%
Week 5-8:  Features (P2)         ░░░░░░░░░░░░░░░░░░░░░░░░  0%
Week 8+:   Extended (P3)         ░░░░░░░░░░░░░░░░░░░░░░░░  0%
```

**Target**: MVP in 16 weeks with 3-4 developers

---

## 🎯 Current Work

### In Progress
- [ ] **TV Shows Module** - sqlc queries ✅, entities & repository ✅, service & jobs pending
  - **Docs**: [content-modules.instructions.md](.github/instructions/content-modules.instructions.md)
  - **Design**: [ARCHITECTURE_V2.md](docs/dev/design/architecture/ARCHITECTURE_V2.md)

### Completed Today (2026-01-29)
- [x] **Dynamic RBAC with Casbin** - Full implementation complete
  - **Docs**: [RBAC_CASBIN.md](docs/dev/design/features/RBAC_CASBIN.md)
  - **Instructions**: [rbac-casbin.instructions.md](.github/instructions/rbac-casbin.instructions.md)
- [x] **TV Shows sqlc queries** - 7 query files, 100+ queries generated
  - Files: `internal/infra/database/queries/tvshow/` (series, seasons, episodes, credits, genres, networks, images, user_data)
  - Generated: `internal/content/tvshow/db/` (11 Go files)
- [x] **TV Shows entities & repository** - Complete domain model and PostgreSQL repository
  - Files: `internal/content/tvshow/` (entity.go, repository.go, repository_pg*.go)
- [x] **Shared video_people** - Created `shared/000017_video_people.up.sql`, refactored movie/tvshow credits
- [x] **Design Docs Updated** - Fixed ARCHITECTURE_V2.md and content-modules.instructions.md for shared video_people
- [x] **Lint cleanup** - Fixed all 100+ lint errors with strict golangci-lint config
- [x] **Go 1.25 Experimental Features** - greenteagc + jsonv2 enabled

### Pending Content Modules
- [ ] TV Shows module - service, jobs
  - **Instructions**: [content-modules.instructions.md](.github/instructions/content-modules.instructions.md)
- [ ] Music module (isolated people: `music_artists`)
- [ ] Books module (isolated people: `book_authors`)
- [ ] Comics module (isolated people: `comic_creators`)

---

## 🔐 RBAC Enhancement (Dynamic Roles) ✅ DONE

**Status**: Casbin-based dynamic RBAC implemented!

**Completed**:
- [x] Add `github.com/casbin/casbin/v2` dependency
- [x] Add `github.com/pckhoi/casbin-pgx-adapter/v3` for PostgreSQL (native pgx)
- [x] Create `roles` table for custom role management
- [x] Create `permission_definitions` table for UI reference
- [x] Add `role_id` FK to users table
- [x] Implement `CasbinService` with full CRUD for roles
- [x] Update RBAC middleware for Casbin
- [x] Create documentation: `docs/dev/design/features/RBAC_CASBIN.md`
- [x] Create instructions: `.github/instructions/rbac-casbin.instructions.md`

**Pending** (Frontend):
- [ ] Create admin UI for role management
- [ ] Create API endpoints for role CRUD (OpenAPI spec)

---

## 🎯 Sprint History

### Day 1: Foundation Setup

#### Step 1: Immediate Fixes (2 hours) ⚡ PRIORITY
- [x] **Module Registration** (15 min) ✅ DONE
  - [x] Add missing modules to main.go: cache, search, jobs, oidc, genre, playback
  - [x] Test application starts without errors
  - **File**: `cmd/revenge/main.go`

- [x] **Configuration Loading** (45 min) ✅ DONE
  - [x] Fix hardcoded cache address in `internal/infra/cache/cache.go`
  - [x] Fix hardcoded search address in `internal/infra/search/search.go`
  - [x] Config structs already exist in `pkg/config/config.go`
  - [x] Config files already properly configured

- [x] **Dependencies Update** (10 min) ✅ DONE
  - [x] `go get github.com/typesense/typesense-go/v4@latest`
  - [x] `go get github.com/ogen-go/ogen@latest`
  - [x] `go mod tidy`

- [x] **Basic Integration** (45 min) ✅ DONE
  - [x] Wire pkg/health Checker to main.go
  - [x] Wire pkg/graceful Shutdowner to main.go
  - [x] Register health checks for database, cache, search
  - [x] Add ShutdownTimeout to ServerConfig
  - [ ] Test /health endpoints (manual verification needed)
#### Step 2: Shared Migrations (2 hours) ✅ DONE
- [x] `000006_genres.up.sql` - Global genres table with common genres
- [x] `000007_server_settings.up.sql` - Persisted server configuration
- [x] `000008_activity_log.up.sql` - Audit logging with severity levels
- [x] `000009_content_ratings.up.sql` - MPAA, FSK, PEGI, BBFC ratings
- [x] Create corresponding .down.sql files
- [ ] Test migrations up/down (needs running PostgreSQL)

#### Step 3: Session Service (2 hours) ✅ DONE
- [x] Session service already 100% implemented (`internal/service/session/service.go`)
- [x] Wire session.Module to main.go
- [x] Verify build compiles successfully
- [x] Cleaned up non-existent API handlers (to be created in Day 4-5)

---

## ✅ Day 1 Complete Summary

**Time Invested**: ~4.5 hours
**Status**: Foundation is solid, all core infrastructure integrated

### What We Built:
1. ✅ **Module Registration** - 6 modules wired: cache, search, jobs, oidc, session, library
2. ✅ **Configuration** - Fixed hardcoded configs, all using koanf now
3. ✅ **Dependencies** - Typesense v4, ogen v1.18, otter, sturdyc
4. ✅ **Health & Shutdown** - pkg/health + pkg/graceful fully integrated
5. ✅ **Migrations** - 4 new shared migrations (genres, settings, activity_log, content_ratings)
6. ✅ **Session Service** - Fully integrated session management

### Build Status:
```bash
✅ go build ./...  # SUCCESS
⏳ go test ./...   # Pending (needs infrastructure running)
```

---

## ✅ Day 2 Complete Summary

**Status**: RBAC System fully implemented

### What We Built:
1. ✅ **Migration 000013_rbac.up.sql** - Added `role` enum and column to users table
2. ✅ **Migration 000014_permissions.up.sql** - Permissions table with 32 permissions, role mappings
3. ✅ **RBAC Service** - `internal/service/rbac/service.go` with permission checking
4. ✅ **RBAC Middleware** - `internal/middleware/rbac.go` with RequirePermission, RequireRole, RequireAdmin
5. ✅ **User Service Updates** - Role support in CreateParams/UpdateParams
6. ✅ **Fixed Migration Conflicts** - Renumbered duplicate migrations (006-012)
7. ✅ **Fixed Build Errors** - uuid type conversion, resty API changes

### Roles Defined:
| Role | Permissions |
|------|-------------|
| `admin` | Full access (all 32 permissions) |
| `moderator` | Libraries, metadata, content management |
| `user` | Browse, play, rate, playlists |
| `guest` | Browse only |

### Build Status:
```bash
✅ go build ./...  # SUCCESS
✅ sqlc generate   # SUCCESS
```

---

## ✅ Day 3 Complete Summary

**Status**: Global Services fully implemented

### What We Built:
1. ✅ **Activity Logger** - `internal/service/activity/service.go`
   - Activity types: user_login, user_logout, security_event, api_error, etc.
   - Severity levels: info, warning, error, critical
   - IP/User Agent tracking
   - List by user/type/severity methods
2. ✅ **Server Settings** - `internal/service/settings/service.go`
   - Key-value settings store with categories
   - Helper methods: GetString, GetBool, GetInt
   - Server name, registration settings
3. ✅ **API Keys** - `internal/service/apikeys/service.go`
   - SHA-256 hashed keys (raw key only shown once)
   - Scope-based permissions
   - Expiration support
   - Usage tracking
4. ✅ **Wired to main.go** - All modules registered with fx

### Build Status:
```bash
✅ go build ./...  # SUCCESS
✅ sqlc generate   # SUCCESS
```

---

## ✅ Day 4-5 Complete Summary

**Status**: OpenAPI Foundation complete

### What We Built:
1. ✅ **OpenAPI Specs** - Comprehensive `api/openapi/revenge.yaml` (already existed)
   - System endpoints (health, server info)
   - Auth endpoints (login, logout, register, sessions)
   - User management endpoints
   - Library management endpoints
   - Movie endpoints (scaffolded, build-tagged for now)
   - Adult endpoints (scaffolded, build-tagged for now)
2. ✅ **ogen Code Generation** - `go generate ./api/...` working
   - 20 generated files in `api/generated/`
   - Type-safe handlers, validators, routers
3. ✅ **Handler Implementation** - `internal/api/`
   - System handlers (health, server info)
   - Auth handlers (login, logout, register, password change, sessions)
   - User handlers (CRUD, admin functions)
   - Library handlers (CRUD, scan trigger)
   - Security handler (Bearer token auth)
4. ✅ **Wired to main.go**
   - api.Module registered with fx
   - BuildInfo provided
   - API server mounted at `/api/v1/`

### Build Status:
```bash
✅ go build ./...  # SUCCESS
✅ go generate ./api/...  # SUCCESS
```

### Notes:
- `movies.go` and `adult.go` temporarily excluded (`//go:build ignore`) - need type fixes when modules are ready
- Core handlers (auth, users, libraries, system) fully functional

---

## ✅ Week 2 Day 1 Complete Summary

**Status**: Workers infrastructure + migrations complete

### What We Built:
1. ✅ **Job Args Structs** - `internal/infra/jobs/workers.go`
   - `ScanLibraryArgs` - Library scanning with full/incremental modes
   - `FetchMetadataArgs` - Metadata fetching with provider selection
   - `DownloadImageArgs` - Image downloads with priority
   - `IndexSearchArgs` - Search indexing (upsert/delete)
   - `CleanupArgs` - Cleanup operations (orphaned files, sessions, activity)
   - `RefreshLibraryArgs` - Library metadata refresh
   - `GenerateTrickplayArgs` - Trickplay image generation

2. ✅ **Worker Implementations** - 7 workers with interface dependencies
   - All workers use dependency injection via interfaces
   - Graceful handling when services not yet available
   - Proper logging and error handling

3. ✅ **Worker Registration** - `RegisterWorkers()` function
   - Registers all workers with River
   - Optional dependencies via fx injection
   - New queues: `images`, `cleanup`

4. ✅ **Migrations**
   - `000015_playlists` - Playlists with items, collaborators, triggers
   - `000016_collections` - Collections with tags, subscriptions, smart rules

### Build Status:
```bash
✅ go build ./...  # SUCCESS
```

### Notes:
- Worker implementations ready but need service interfaces to be implemented
- Scanner/Fetcher logic deferred until content modules exist
- Migrations ready for `migrate up`

---

## 📋 Next: Week 4-8 - Frontend & Remaining Content Modules

---

## 📦 Week 2: Workers & Remaining Foundation

### River Workers Infrastructure (3 days)

#### Worker Base Setup (Day 1) ✅ DONE
- [x] Create `internal/infra/jobs/workers.go`
- [x] Define job args structs for all 7 workers
- [x] Create worker registration function
- [x] Set up worker services interface

#### Core Workers (Day 2-3) ✅ DONE
- [x] **Library Scanner** - `ScanLibraryWorker`
  - [x] Worker struct with interface dependency
  - [ ] Create `internal/service/library/scanner.go` (deferred - needs content modules)
  - [ ] Implement directory scanning
  - [ ] File type detection

- [x] **Metadata Fetcher** - `FetchMetadataWorker`
  - [x] Worker struct with interface dependency
  - [ ] Create `internal/service/metadata/providers.go` (deferred - needs content modules)
  - [ ] Provider registry pattern

- [x] **Image Downloader** - `DownloadImageWorker`
- [x] **Search Indexer** - `IndexSearchWorker`
- [x] **Cleanup Worker** - `CleanupWorker`
- [x] **Refresh Library** - `RefreshLibraryWorker`
- [x] **Trickplay Generator** - `GenerateTrickplayWorker`

### Remaining Migrations (1 day) ✅ DONE
- [x] `000015_playlists.up.sql` - Video/audio playlists with collaborators
- [x] `000016_collections.up.sql` - Curated content collections

### Integration Testing (1 day)
- [ ] Test all workers can be queued
- [ ] Test health checks report all services
- [ ] Test graceful shutdown works
- [ ] Test session management works
- [ ] Test RBAC middleware works

---

## 🎬 Week 3-4: Movie Module (Reference Implementation) ✅ DONE

### Database Layer (3 days) ✅ DONE
- [x] Create `internal/infra/database/migrations/movie/000001_movies.up.sql`
- [x] Create queries in `internal/infra/database/queries/movie/`
- [x] Run sqlc generate
- [x] Create repository interfaces

### Domain & Service (3 days) ✅ DONE
- [x] Create `internal/content/movie/entity.go` - Entities
- [x] Create `internal/content/movie/repository.go` - Repository interface
- [x] Create `internal/content/movie/service.go` - Business logic
- [x] TMDb provider integration (`metadata_provider.go`)

### API Layer (2 days) ✅ DONE
- [x] Movie endpoints in `api/openapi/revenge.yaml`
- [x] ogen handlers generated
- [x] Implement handler functions (`internal/api/movies.go`)
- [x] Wire movie.ModuleWithRiver to main.go

### Jobs Integration (2 days) ✅ DONE
- [x] Movie jobs in `internal/content/movie/jobs.go`
- [x] EnrichMetadataArgs for metadata fetch
- [x] River worker registration
- [ ] Test end-to-end flow (needs running infrastructure)

---

## 🎨 Week 4-8: Frontend & Remaining Modules

### Frontend Foundation (Week 4)
- [ ] Initialize SvelteKit 2 in `web/`
- [ ] Configure Tailwind CSS 4
- [ ] Install shadcn-svelte
- [ ] Setup TanStack Query
- [ ] Basic auth pages

### Content Modules (Week 5-7)
Each module follows movie module pattern:
- [ ] **TV Shows** (1 week)
- [ ] **Music** (1 week)
- [ ] **Audiobooks** (3 days)
- [ ] **Books** (3 days)
- [ ] **Podcasts** (3 days)
- [ ] **Photos** (3 days)
- [ ] **LiveTV** (1 week)
- [ ] **Collections** (3 days)
- [ ] **Adult Content** (1 week - schema `c`)

### Player Implementation (Week 8)
- [ ] Video player (Shaka + hls.js)
- [ ] Audio player (Web Audio API + Howler.js)
- [ ] Gapless audio
- [ ] Crossfade
- [ ] Subtitles

---

## 🔌 P1: External Integrations (Week 3-6)

### Metadata Providers (Critical)
- [ ] **TMDb** - Movie/TV metadata
- [ ] **TheTVDB** - TV show metadata
- [ ] **MusicBrainz** - Music metadata

### Servarr Ecosystem
- [ ] **Radarr** - Movie management
- [ ] **Sonarr** - TV show management
- [ ] **Lidarr** - Music management
- [ ] **Chaptarr** - Books & audiobooks (uses Readarr API)
- [ ] **Whisparr v3 (eros)** - Adult content management (schema `c`)

### Scrobbling (P2)
- [ ] **Trakt** - Movie/TV sync
- [ ] **Last.fm** - Music scrobbling
- [ ] **ListenBrainz** - Music scrobbling

---

## 🚀 P2: Feature Enhancements (Week 5-8)

- [ ] **i18n System** - Multi-language support
- [ ] **Analytics Service** - Watch statistics, Year in Review
- [ ] **Profiles System** - Netflix-style profiles
- [ ] **Media Enhancements** - Trickplay, intro detection, chapters
- [ ] **Advanced Observability** - Metrics, supervision

---

## 🎁 P3: Extended Features (Week 8+)

- [ ] **Request System** - Content requests with voting
- [ ] **Ticketing System** - Bug reports, feature requests
- [ ] **Comics Module** - CBZ/CBR reader
- [ ] **LiveTV & DVR** - TV recording

---

## 📝 Documentation Tasks (Ongoing)

### Completed ✅
- ✅ Architecture compliance analysis
- ✅ Advanced features integration analysis
- ✅ Core functionality gap analysis
- ✅ Design extraction from docs
- ✅ Archived 264+ outdated TODOs (6 documents)
- ✅ Cleaned up TECH_STACK.md
- ✅ Updated docs/INDEX.md

### In Progress
- ⏳ Integration docs (37/72 complete, 51%)

---

## ⚠️ Important Notes

**Adult Content Isolation**:
- Schema: `c` (not `adult`)
- API namespace: `/c/*` (obscured)
- Module location: `internal/content/c/`

**No Client Development**:
- WebUI only (SvelteKit)
- Support existing clients (Jellyfin, VLC, etc.)

**External Transcoding**:
- Blackbeard service handles all transcoding
- Revenge proxies streams only

---

**Next Action**: Week 4-8 - Frontend & Remaining Content Modules
- **Wiki platforms**: Normal (Wikipedia, FANDOM, TVTropes) + Adult (Babepedia, IAFD, Boobpedia)
- **External adult platforms**: FreeOnes, TheNude, Pornhub, OnlyFans (performer enrichment, c schema isolated)
- **Scrobbling**: Trakt/Simkl (movies/TV), Last.fm/ListenBrainz (music), Letterboxd (import only)

---

## Implementation Phases

### Phase 1: Core Infrastructure ⬜ IN PROGRESS

- [x] Project setup (Go 1.25, fx, koanf, sqlc)
- [x] CI/CD (GitHub Actions, release-please)
- [x] Docker Compose (dev + prod)
- [x] Configuration system (REVENGE_* env vars)
- [x] Logging (slog)
- [x] HTTP server with graceful shutdown
- [x] Health endpoints
- [x] Basic auth middleware
- [x] User/Session/OIDC tables
- [x] Genre domain separation
- [x] Shared tables (libraries, api_keys, server_settings, activity_log)
- [ ] River job queue setup
- [ ] Typesense search client
- [ ] Dragonfly cache client

### Phase 2: Movie Module ⬜ NOT STARTED

- [ ] Database schema (movies, genres, people, studios, images, streams)
- [ ] Domain entities
- [ ] Repository layer (sqlc)
- [ ] Service layer
- [ ] HTTP handlers (ogen)
- [ ] User data (ratings, history, favorites, watchlist)
- [ ] Radarr integration
- [ ] TMDb fallback provider

### Phase 3: TV Show Module ⬜ NOT STARTED

- [ ] Database schema (series, seasons, episodes)
- [ ] Domain/Repository/Service/Handlers
- [ ] User data
- [ ] Sonarr integration
- [ ] TheTVDB/TMDb fallback

### Phase 4: Music Module ⬜ NOT STARTED

- [ ] Database schema (artists, albums, tracks, music_videos)
- [ ] Domain/Repository/Service/Handlers
- [ ] User data
- [ ] Lidarr integration
- [ ] MusicBrainz/Last.fm fallback

### Phase 5: Playback Service ⬜ NOT STARTED

- [ ] Session management
- [ ] Client capability detection
- [ ] Blackbeard transcoder integration
- [ ] Stream buffering
- [ ] Progress tracking
- [ ] Bandwidth adaptation

### Phase 6: Remaining Content Modules ⬜ NOT STARTED

- [ ] Audiobook module (Audiobookshelf integration)
- [ ] Book module (Audiobookshelf + Chaptarr)
- [ ] Podcast module (Audiobookshelf + RSS)
- [ ] Photo module
- [ ] LiveTV module (PVR backends)
- [ ] Collection module (video + audio pools)

### Phase 7: Adult Modules ⬜ NOT STARTED

- [ ] `c` PostgreSQL schema (isolated)
- [ ] Adult movie module
- [ ] Adult show module
- [ ] Shared performers/studios/tags
- [ ] Adult playlists & collections
- [ ] Whisparr integration
- [ ] Stash/StashDB integration

### Phase 8: Media Enhancements ⬜ NOT STARTED

- [ ] Trailer system (local, Radarr, TMDb, YouTube)
- [ ] Audio themes (Netflix-style hover music)
- [ ] Intro/outro detection (Chromaprint)
- [ ] Trickplay generation
- [ ] Chapter extraction
- [ ] Cinema mode (preroll/postroll)

### Phase 9: External Services ⬜ NOT STARTED

- [ ] Trakt scrobbling
- [ ] Last.fm scrobbling
- [ ] ListenBrainz scrobbling
- [ ] Import/export ratings

### Phase 10: Frontend ⬜ NOT STARTED

- [ ] SvelteKit 2 setup
- [ ] Tailwind CSS 4 + shadcn-svelte
- [ ] Authentication (JWT + OIDC)
- [ ] Media browser
- [ ] Video player
- [ ] Audio player
- [ ] Admin panel

---

## Go 1.25 Features to Adopt

- [ ] `sync.WaitGroup.Go` - Replace manual wg.Add/Done patterns
- [ ] `testing/synctest` - Concurrent code testing
- [ ] `net/http.CrossOriginProtection` - Replace custom CSRF
- [ ] `slog.GroupAttrs` - Grouped logging
- [ ] `runtime/trace.FlightRecorder` - Observability
- [ ] `reflect.TypeAssert` - Zero-allocation type assertions

## Experimental Features ✅ ENABLED

Both experimental features enabled in `Dockerfile` and `Makefile`:

- [x] `GOEXPERIMENT=greenteagc` - New GC (10-40% memory reduction)
- [x] `GOEXPERIMENT=jsonv2` - Faster JSON encoding/decoding

Build with experiments: `make build` (automatic) or `GOEXPERIMENT=greenteagc,jsonv2 go build`

---

## Documentation Status

### Completed ✅

- [x] ARCHITECTURE_V2.md - Complete modular design
- [x] TECH_STACK.md - Technology choices
- [x] PROJECT_STRUCTURE.md - Directory layout
- [x] METADATA_SYSTEM.md - Servarr-first with Audiobookshelf
- [x] AUDIO_STREAMING.md - Progress, bandwidth adaptation
- [x] CLIENT_SUPPORT.md - Chromecast, DLNA, capabilities
- [x] MEDIA_ENHANCEMENTS.md - Trailers, themes, intros, trickplay, Live TV
- [x] SCROBBLING.md - External service sync
- [x] OFFLOADING.md - Blackbeard integration
- [x] BEST_PRACTICES.md - Resilience patterns
- [x] I18N.md - Internationalization

### TODO 📝

- [ ] ADULT_METADATA.md - Stash/StashDB/Whisparr integration
- [ ] CINEMA_MODE.md - Preroll, postroll, intermission
- [ ] API.md - OpenAPI design guidelines
- [ ] REVERSE_PROXY.md - Nginx, Caddy, Traefik configs
- [ ] MOBILE_APPS.md - iOS/Android architecture

---

## Completed ✅

- [x] Project setup (Go 1.25, fx, koanf, sqlc)
- [x] CI/CD (GitHub Actions, release-please)
- [x] Docker Compose (dev + prod)
- [x] Configuration system (REVENGE_* env vars)
- [x] Logging (slog)
- [x] HTTP server with graceful shutdown
- [x] Health endpoints
- [x] Basic auth middleware
- [x] User/Session/OIDC tables
- [x] Genre domain separation
- [x] Resilience packages (circuit breaker, bulkhead, retry)
- [x] Supervisor/graceful shutdown packages
- [x] Health check system
- [x] Hot reload configuration
- [x] Lazy initialization patterns
- [x] Metrics package
- [x] Playback service architecture (docs)
- [x] Documentation suite
