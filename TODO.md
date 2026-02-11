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

### Cache Lifecycle, Retry Filtering & Observability Hardening (2026-02-08)
- otter L1 cache: fixed TTL bug (ttl<=0 → ttl<0), added L1Option system (WithExpiryAccessing, WithOnDeletion)
- Added ExpiryAccessing to all 15 provider/integration caches (read-on-access refreshes TTL)
- Added Close() to all 14 provider/integration clients (prevents otter goroutine leaks)
- Added OnDeletion callback in transcode pipeline (kills evicted FFmpeg processes)
- Added SetCommonRetryCondition to all 17 HTTP clients (skip retries on 4xx, retry only 5xx/network)
- Replaced rueidis.NewClient with rueidisotel.NewClient for OTel instrumentation
- S3: use s3manager.Uploader for automatic multipart on large files
- Raft: merged raft-log.db + raft-stable.db into single raft.db
- Typesense: replaced manual URL parsing with net/url.Parse, added 5s connection timeout
- vips: added StartupConfig with ConcurrencyLevel and MaxCacheSize
- AdminListUsers: fixed error codes (400→403 with descriptive messages)
- Deprecated JobsQueueSize metric (dead code, replaced by river_queue_size)
- Cleanup job: added MaxAttempts:5 (was global default 25)
- Fixed argon2id test params (p=4→p=2 to match production)
- Fixed misleading cache module comment

### Cookie Auth + CSRF Protection (2026-02-07)
- Cookie-based auth middleware: extracts HttpOnly access token cookie → injects as Bearer header
- CSRF double-submit cookie pattern with `X-CSRF-Token` header validation
- ResponseWriter context injection for ogen handlers to set cookies
- Login/refresh set HttpOnly cookies (access, refresh) + JS-readable CSRF cookie
- Logout clears all auth cookies
- Fully opt-in via `server.cookie_auth.enabled` config

### SSE Real-Time Events (2026-02-07)
- `GET /api/v1/events` Server-Sent Events endpoint
- Broker with per-client category filtering and non-blocking broadcast
- Auth via Bearer header or `?token=` query param
- Bridges `notification.Agent` interface → SSE fanout
- 30s keepalive, ResponseController for write deadline management

### Tier 3 Metadata Providers (2026-02-07)
- **Trakt** (priority 38): Movie + TV show metadata, extended info, rate-limited
- **Simkl** (priority 36): Movie + TV + anime metadata, cross-reference IDs
- **Letterboxd** (priority 34): Movie-only metadata, OAuth2 client credentials, rating conversion

### Tier 2 Metadata Providers (2026-02-06)
- **AniList** (priority 45): GraphQL-based anime/manga metadata
- **AniDB** (priority 44): UDP + HTTP anime metadata with rate limiting
- **MAL** (priority 43): MyAnimeList metadata via Jikan API
- **Kitsu** (priority 42): JSON:API anime/manga metadata

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
- [x] Trakt (metadata provider, completed 2026-02-07)
- [ ] Last.fm, ListenBrainz (scrobbling)
- [ ] Additional metadata providers (MusicBrainz, Spotify, etc.)

### Features
- [ ] Collections and playlists
- [ ] Watch Next / Continue Watching
- [ ] Skip Intro / Credits detection
- [ ] SyncPlay (watch together)
- [ ] Trickplay (timeline thumbnails)
- [x] SSE real-time events (completed 2026-02-07)
- [x] Cookie-based authentication + CSRF (completed 2026-02-07)
- [ ] Release calendar
- [ ] Content request system

### Infrastructure
- [ ] Circuit breaker integration (gobreaker dep exists, unused)
- [ ] Request coalescing (sturdyc dep exists, unused)
- [ ] Cache warming on startup
- [ ] Service self-healing / watchdog
