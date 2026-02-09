# Deep Scan TODO

**Created**: 2026-02-07
**Source**: DEEP_SCAN_REPORT.md

## Priority: Quick Fixes

- [x] **BUG-4**: Remove outdated auth config comment (`config.go:194` says "auth not implemented yet")
- [x] **BUG-3**: Register CleanupWorker with River (cleanup job defined but never scheduled)
- [x] **DEAD-1**: Wire JobsConfig fields to River (FetchPollInterval, RescueStuckJobsAfter, MaxWorkers)
- [x] **TODO-1**: Fix outdated comment in `internal/config/config.go:194`

## Priority: Medium

- **UNUSED-1**: `gobreaker` in go.mod — NOT unused, planned for circuit breaker implementation
- **UNUSED-2**: `sturdyc` in go.mod — NOT unused, planned for request coalescing/caching
- [x] **TODO-2**: Localization TODO in `internal/api/localization.go` — implemented `GetMetadataLanguage()` reading user prefs → Accept-Language → default "en"

## Priority: Known Bugs (not quick)

- [x] **BUG-1**: Health check does not detect database failures — implementation was already correct (`checkDatabase()` calls `Ping()` + `SELECT 1`); unskipped integration test
- [x] **BUG-2**: Typesense health check broken with v2 client — fixed timeout (10ns → 10s) + use boolean return value; unskipped integration test

## Priority: Security (tracked in .workingdir/TODO_A7_SECURITY_FIXES.md)

- [x] **A7.1**: Missing transaction boundaries — added tx to `Register()` and `VerifyEmail()`; `ChangePassword`/`ResetPassword`/`UploadAvatar` already had them; `RefreshSession` uses safe ordering
- [x] **A7.2**: Login timing attack — already fixed (dummy hash constant-time comparison)
- [x] **A7.3**: Goroutine leak in notification dispatcher — already fixed (WaitGroup + stopCh + Close)
- [x] **A7.4**: Password reset token info disclosure — already fixed (returns only error, silent success)
- [x] **A7.5**: No service-level rate limiting for Argon2id — already fixed (account lockout implemented)
- [x] **A7.6**: context.Background() in goroutines — already fixed (all have timeouts)

## Completed This Session (2026-02-07)

- [x] **sync.Map → L1Cache migration**: Replaced all sync.Map usage with otter-based L1Cache, updated all docs
- [x] **sqlc.yaml fix**: Fixed broken schema paths (`migrations/` → `internal/infra/database/migrations/shared/`)
- [x] **metadata_language feature**: New DB column, sqlc/ogen codegen, user service params, API handler, OpenAPI spec
- [x] **Localization rewrite**: `GetMetadataLanguage(ctx)` reads user prefs → Accept-Language → "en"; removed broken `getRequestFromContext()`
- [x] **TV show localization**: Added `LocalizeSeries`/`LocalizeSeriesList`, wired into 5 TV show handlers
- [x] **Stale test fixes**: Fixed sonarr `client_test.go` and `ratelimit_test.go` for L1Cache API

## HLS Streaming (2026-02-07)

- [x] Phase 1: Config + Core Types (PlaybackConfig, Session, profiles)
- [x] Phase 2: FFmpeg wrapper + pipeline + subtitle extraction
- [x] Phase 3: HLS manifest generation + HTTP stream handler
- [x] Phase 4: PlaybackService + OpenAPI spec + ogen codegen + API handler
- [x] Phase 5: River cleanup job + fx module (playbackfx) + app wiring
- [x] Build passes, lint clean (0 issues), all playback tests pass
- [x] Fix `gin_trgm_ops` migration issue (added `CREATE EXTENSION pg_trgm` to 000001)
- [x] Add quality policy to CLAUDE.md ("no shortcuts, no half measures")
- [x] Document full HLS architecture in `.workingdir4/HLS_STREAMING.md`

### Architecture: Separate Audio Renditions
Video-only segments per quality profile + audio-only segments per track.
HLS.js downloads only the active track. Switching is instant, zero bandwidth waste.
HLS-compatible codecs (AAC, AC-3, E-AC-3) copied at original quality.

## Test Infrastructure Fixes (2026-02-07)

### Migration Bugs Fixed
- [x] `shared.update_updated_at_column()` missing — referenced by 6 movie migrations but never created; added to 000001
- [x] `auth.users` FK reference — 000032 had `REFERENCES auth.users(id)`, fixed to `shared.users(id)`
- [x] `pg_trgm` extension — added `CREATE EXTENSION IF NOT EXISTS pg_trgm` to 000001

### Test Bug Fixes
- [x] RBAC test: `"library"` → `"libraries"` resource name (FineGrainedResources rename)
- [x] UUID v7 truncation collisions — `uuid.NewV7().String()[:8]` in loops causes duplicates; fixed in session, apikeys, mfa, settings tests
- [x] Avatar upload test — `UploadAvatar` bypasses mock repo (uses `s.pool` directly); rewrote with real DB user

### Test DB Consolidation
- [x] `internal/infra/database/testing.go` — ONE shared embedded postgres (port 15600, sync.Once)
  - `setupTestDB()` creates per-test databases with migrations on shared instance
  - `setupFreshTestDB()` creates per-test databases WITHOUT migrations (for migration tests)
  - `setupTestPool()` returns pgxpool.Pool on shared instance
  - `freshTestDBURL()` returns URL only (for NewPool/MigrateUp tests)
  - `createTestDatabase()` shared helper for DB creation + cleanup
  - `StopSharedPG()` called from TestMain
- [x] `internal/infra/database/migrations_test.go` — removed 5 embedded postgres instances (ports 15432-15436)
  - Table structure tests use `assertColumnsExist()` — only checks required columns, extra columns from new migrations allowed
  - MigrateDown checks version decrease, not hardcoded migration effects
- [x] `internal/infra/database/migrate_test.go` — removed 12 embedded postgres instances (ports 15540-15553)
- [x] `internal/infra/database/auth_tokens_test.go` — uses shared helpers
- [x] `internal/infra/health/service_test.go` — uses testutil.NewTestDB() instead of own embedded postgres

### Port Allocation (no conflicts)
| Package | Port | Purpose |
|---------|------|---------|
| `testutil/database.go` | 15432 | Legacy test helper |
| `testutil/testdb.go` | 15555 | Shared test DB (template pattern) |
| `api/handler_test.go` | 15438-15440 | API handler tests |
| `api/server_test.go` | 15450-15470 | API server tests |
| `infra/database/testing.go` | 15600 | Database package tests (shared) |

### Remaining (API tests)
- [ ] `internal/api/server_test.go` — 12 embedded postgres instances could be consolidated
- [ ] `internal/api/handler_test.go` — 3 embedded postgres instances could be consolidated

## Upcoming Tasks

- [ ] **Client bandwidth detection** — real-time measurement, server-side ABR decisions
- [ ] **User/admin quality settings** — max quality preference, transcoding limits
- [ ] **QoS prioritization** — stream segments > UI requests, priority queuing
- [ ] **Search infrastructure audit** — verify Typesense integration completeness
- [ ] **Hardware acceleration** — VAAPI, NVENC, QSV support
- [ ] **Custom FFmpeg build** — revenge-ffmpeg with all needed codecs
- [ ] **Watch tracking gaps** — `completed_at` column unused, TV watch history endpoint missing

## Notes

- Placeholder SQL files (movie, tvshow, qar) are intentional — sqlc requires non-empty query dirs
- Windows media prober stub is intentional — platform limitation
- Security issues are tracked separately in `.workingdir/TODO_A7_SECURITY_FIXES.md`
