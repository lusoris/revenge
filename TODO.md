# Project TODO

## Current Priorities

### Documentation Rewrite
**Status**: Planning complete, execution starting
**Plan**: [.workingdir3/PLAN.md](.workingdir3/PLAN.md)

218 design docs are out of sync with codebase. Incremental 11-step plan to:
1. Fix root MDs (this file, README, CONTRIBUTING)
2. Fix status tables in design docs
3. Move ~150 unimplemented feature docs to `planned/`
4. Clean up auto-generated templates and wiki duplicates
5. Rewrite architecture, service, infra, module docs from actual code

### Frontend Development
**Status**: Not started
**Impact**: No UI for end users

Design docs exist (SvelteKit 2, Svelte 5, Tailwind CSS 4, shadcn-svelte) but zero frontend code.

---

## Recently Completed

### Metadata System - Force + Languages Plumbing (2026-02-06)
- Added `MetadataRefreshOptions{Force, Languages}` to movie and tvshow modules
- Plumbed through entire stack: service interfaces -> adapters -> workers
- Added `ClearCache()` to `metadata.Provider` and `metadata.Service` interfaces
- Adapters clear cache on `Force=true`, use per-request languages when provided
- Workers construct opts from job args and pass to service

### River Job Workers (2026-02-06)
- 9 workers fully implemented (was previously a stub):
  - Movie: MetadataRefresh, LibraryScan, FileMatch, SearchIndex
  - TV Show: MetadataRefresh, LibraryScan, FileMatch, SeriesRefresh, SearchIndex
  - Shared: LibraryScanCleanup, ActivityCleanup
- Progress reporting via `river.JobProgress`
- Worker timeouts and retry configuration

### CI Fixes (2026-02-06)
- Fixed govulncheck job missing CGO dependencies (libvips-dev, libav*-dev)
- Fixed test migration paths (was `migrations/`, now `internal/infra/database/migrations/shared/`)
- Fixed service test signatures (`NewService(repo)` -> `NewService(repo, nil)`)

### Password Hash Migration (2026-02-03)
- Hybrid password verifier with bcrypt backward compatibility
- `NeedsMigration()` helper, 4 comprehensive tests

### Dependency Upgrades (2026-02-05)
- Replaced shopspring/decimal with govalues/decimal (performance)
- Replaced x/image with govips for image processing

---

## Known Issues

### Security: G602 Slice Bounds in Generated Code
**Location**: `internal/api/ogen/oas_router_gen.go` (generated code)
**Count**: 10 issues
**Note**: In ogen-generated code, not hand-written. Will be fixed upstream or via ogen config.

### Security: G101 False Positives in SQLC Code
**Location**: `internal/infra/database/db/*.sql.go` (generated code)
**Count**: 43 issues
**Note**: SQL query names flagged as potential credentials. All false positives.

---

## Future Work (by area)

### Content Modules
- [ ] Music module (highest priority new module)
- [ ] Audiobook module
- [ ] Book module
- [ ] Podcast module
- [ ] Comics module
- [ ] Photos module
- [ ] Live TV/DVR module

### Services
- [ ] Playback service (progress tracking, watch history)
- [ ] Transcoding service (Blackbeard integration)
- [ ] Analytics service
- [ ] Grants service (fine-grained sharing)
- [ ] Fingerprint service (media identification)

### Integrations
- [ ] Lidarr (music)
- [ ] Whisparr (adult content)
- [ ] Authelia, Authentik, Keycloak (SSO providers)
- [ ] Trakt, Last.fm, ListenBrainz (scrobbling)
- [ ] Additional metadata providers (MusicBrainz, Spotify, etc.)

### Features
- [ ] Collections and playlists
- [ ] Watch Next / Continue Watching
- [ ] Skip Intro / Credits detection
- [ ] SyncPlay (watch together)
- [ ] Trickplay (timeline thumbnails)
- [ ] Release calendar
- [ ] Content request system

### Infrastructure
- [ ] Circuit breaker integration (gobreaker dep exists, unused)
- [ ] Request coalescing (sturdyc dep exists, unused)
- [ ] Cache warming on startup
- [ ] Service self-healing / watchdog
