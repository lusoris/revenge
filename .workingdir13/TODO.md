# Workingdir 13 — Master TODO

## Phase 1: Critical Bugs & River Jobs (Priority: IMMEDIATE)

- [x] 1.1 Create .workingdir13 directory
- [x] 1.2 Add `InsertOpts()` with `MaxAttempts` + `UniqueOpts` to ALL job args types that lack them
  - [x] `metadata/jobs/refresh.go`: RefreshMovieArgs, RefreshTVShowArgs, RefreshSeasonArgs, RefreshEpisodeArgs, RefreshPersonArgs, EnrichContentArgs, DownloadImageArgs
  - [x] `moviejobs/file_match.go`: MovieFileMatchArgs
  - [x] `moviejobs/search_index.go`: MovieSearchIndexArgs
  - [x] `tvshow/jobs/jobs.go`: MetadataRefreshArgs, FileMatchArgs, SeriesRefreshArgs
  - [x] `radarr/jobs.go`: RadarrSyncJobArgs, RadarrWebhookJobArgs (add MaxAttempts + UniqueOpts)
  - [x] `sonarr/jobs.go`: SonarrSyncJobArgs, SonarrWebhookJobArgs (add MaxAttempts + UniqueOpts)
  - [x] `activity/cleanup.go`: ActivityCleanupArgs (add MaxAttempts)
  - [x] `analytics/stats_worker.go`: StatsAggregationArgs (add MaxAttempts)
  - [x] `library/cleanup.go`: LibraryScanCleanupArgs (add MaxAttempts)
  - [x] `playback/jobs/cleanup.go`: CleanupArgs (use QueueLow constant instead of "low")
- [x] 1.3 Fix `movie_file_match` returning error on unmatched files (causes 25 retries)
- [x] 1.4 Fix `stats_aggregation` panic when queries is nil
- [x] 1.5 Add `CompletedJobRetentionPeriod` + `DiscardedJobRetentionPeriod` to River config
- [x] 1.6 Reduce global default MaxAttempts from 25 → 5
- [x] 1.7 Add `JobTimeout: -1` to River config (per-worker Timeout() handles this)

## Phase 2: Flow Integration Fixes (Priority: HIGH)

- [x] 2.1 Fix TV show library scan API trigger (currently only enqueues movie scan)
- [x] 2.2 Wire notification dispatch at lifecycle points
- [x] 2.3 Auto-enqueue search index jobs after library scan completion
- [x] 2.4 Fix RefreshPersonWorker to return `river.JobCancel` instead of silently succeeding

## Phase 3: Cache Coverage (Priority: HIGH)

- [x] 3.1 Fix OIDC `sync.Map` → `L1Cache` (violates project rules)
- [x] 3.2 Fix session CachedService: override `RefreshSession` + `RevokeAllUserSessionsExcept`
- [x] 3.3 Cache `ValidateKey` in API Keys service (interface extraction + CachedService wrapper + SHA-256 cache keys)
- [x] 3.4 Cache `GetMovieFiles` (called on every playback)
- [x] 3.5 Cache `ListDistinctGenres` (movie + tvshow, near-static data)

## Phase 4: Code Deduplication (Priority: MEDIUM)

- [x] 4.1 Create `internal/integration/arrbase` shared package (15 shared types, type aliases for backward compat)
- [x] 4.2 Create generic `CacheAside[T]` helper function
- [x] 4.3 Create `setOpt*` generic helpers for ogen optional field mapping (reduced boilerplate in converters)

## Phase 5: Additional Features (Priority: LOW)

- [x] 5.1 Add periodic library scan job (configurable schedule)
- [x] 5.2 Add notification agent configuration from config/DB (NotificationsConfig structs + registerNotificationAgents)
- [x] 5.3 Add playback heartbeat endpoint (POST /api/v1/playback/sessions/{sessionId}/heartbeat)
- [x] 5.4 Update stale DownloadImageWorker comment

## Phase 6: Testing & Validation (Priority: HIGH)

- [x] 6.1 All 66 unit test packages pass
- [x] 6.2 Docker image rebuilt and deployed (`make docker-build` + `docker compose up`)
- [x] 6.3 Live integration tests — 7 new test functions, all pass (`tests/live/new_features_test.go`)
  - Heartbeat endpoint, API key caching, API key edge cases, concurrent API key lifecycle,
    notification config, rate limiting, integration config endpoints
- [x] 6.4 Full live test suite — all pass (0 failures)

## Phase 7: Load Test Expansion (Priority: MEDIUM)

- [x] 7.1 Updated `helpers.js` with `loginFull()`, `apiKeyRequest()`, `randomInt()` utilities
- [x] 7.2 Created `playback_load.js` — playback session lifecycle, heartbeat burst, session abandon/inspect
- [x] 7.3 Created `api_key_load.js` — cached auth perf, key lifecycle, concurrent keys, invalid rejection
- [x] 7.4 Created `write_operations.js` — movie/episode progress, mark watched, mixed R/W, settings
- [x] 7.5 Updated `run.sh` with 3 new test entries
- [x] 7.6 Smoke-tested all 3 new load tests: all pass
  - `api_key_load`: 99.69% check success, 100% auth success (1157/1157), cache hits avg ~1.9ms
  - `write_operations`: 100% check success, 100% write success, 0% HTTP failures
  - `playback_load`: 0% HTTP failures (graceful skip with no content)
- [x] 7.7 Run gentle profile against live image + monitor Prometheus metrics
- [x] 7.8 Identify performance bottlenecks (see findings below)

## Performance Findings from Load Tests (gentle profile, 50→500 VUs)

### Critical: Login OOM Kill (P0)
- **auth_stress @ 500 VUs → container killed (OOM, exit 137)**
- Root cause: Argon2id `DefaultParams` uses 64MB per hash. 500 concurrent logins = 32GB.
- No concurrency semaphore on `PasswordHasher.VerifyPassword`.
- **Fix needed**: Add a bounded semaphore (e.g., `runtime.NumCPU()`) around Argon2id verify
  to cap concurrent memory, or use a worker pool approach.

### High: Global Rate Limiter Blocks 49% of Traffic (P1)
- Rate limit: 10 req/s per IP, burst 20. Under load, **201,343 requests blocked** out of 411,583.
- This is correct per-IP behavior, but in a deployment behind a reverse proxy,
  all requests may appear from the same IP if `X-Forwarded-For` isn't trusted.
- **Action**: Verify `X-Forwarded-For` / `X-Real-IP` header parsing in rate limiter.
  Consider tuning limits or making them configurable per-deployment.

### Medium: Auth Rate Limiter Not Engaging (P2)
- Auth-specific limiter shows 0 allowed, 0 blocked — all auth requests hit the
  global limiter first. Auth limiter operations: LoginUser, VerifyMFA, etc.
- **Action**: Check middleware ordering — auth limiter needs to run on the auth operation
  name AFTER the global limiter passes (or be an independent check).

### Good: Endpoint Performance (OK)
| Endpoint | Avg Latency | Notes |
|---|---|---|
| GET /users/me | 0.28ms | Excellent |
| GET /genres | 0.18ms | Cached, excellent |
| GET /libraries | 0.31ms | Cached, excellent |
| GET /movies | 0.40ms | Good (empty DB) |
| GET /apikeys/{id} | 0.27ms | Good |
| GET /search/movies | 4.3ms | Typesense query, acceptable |
| PUT /settings/user/* | 7.0ms | Write latency, acceptable |
| DELETE /apikeys/{id} | 7.9ms | Key revocation + cache invalidation |
| POST /auth/login | 44.8ms | Argon2id verify, expected |

### Good: Cache Performance (OK)
- **Dragonfly (L2) hit rate: 99.4%** — excellent
- **L1 cache: 335,493 hits / 1,800 misses** — 99.5% hit rate
- API key cache hits avg 0.7ms → validates the Phase 3.3 `CachedService` work

### Good: Resource Utilization (OK)
- **Goroutines: 85** (stable, no leak)
- **Heap: 23MB** (stable after load)
- **PGX pool: 11 conns** (of 65 max) — not saturated
- **Write success rate: 93.8%** under 500 VUs — good

### Pre-existing: realistic_usage.js Config Error
- Uses both `stages` and `scenarios` in k6 options, which is not allowed.
- ~~Needs fix: remove `stages` from top-level since `scenarios` already defines its own stages.~~
- **Fixed**: removed top-level `stages` from options.

## Phase 8: Performance Fixes (Priority: HIGH)

- [x] 8.1 Add concurrency semaphore to `PasswordHasher` (`internal/crypto/password.go`)
  - Semaphore channel caps concurrent Argon2id ops at `2×runtime.NumCPU()`
  - New `HashPasswordContext()` / `VerifyPasswordContext()` methods accept context for cancellation
  - Original `HashPassword()` / `VerifyPassword()` unchanged (backward-compatible, use `context.Background()`)
  - New constructor `NewPasswordHasherWithConcurrency(params, maxConcurrent)` for custom limits
  - Prevents OOM: e.g., on 8-core machine, max 16 concurrent × 64MB = 1GB instead of unbounded
- [x] 8.2 Fix auth rate limiter operation names to match ogen PascalCase
  - `LoginUser` → `Login`, `VerifyMFA` → `VerifyTOTP`, `RequestPasswordReset` → `ForgotPassword`
  - Added `BeginWebAuthnLogin`, `FinishWebAuthnLogin` coverage
  - Fixed in both in-memory (`ratelimit.go`) and Redis (`ratelimit_redis.go`) configs
  - Updated all tests and benchmarks to match
  - **Root cause**: Auth limiter showed 0 allowed / 0 blocked in Prometheus because operation names
    never matched — auth endpoints were only protected by the global limiter
- [x] 8.3 Fix `realistic_usage.js` k6 config (removed conflicting top-level `stages`)
- [x] 8.4 Added semaphore tests: concurrency, context cancellation, context timeout, defaults
- [x] 8.5 All tests pass: crypto (15 tests), middleware (28 tests), auth service, user service, MFA service

## Phase 9: Live Validation of Fixes (Priority: HIGH)

- [x] 9.1 Rebuilt Docker image with all Phase 8 fixes
- [x] 9.2 Verified container healthy, baseline metrics clean (84 goroutines, 11MB heap, all counters 0)
- [x] 9.3 auth_stress @ 500 VUs: **CONTAINER SURVIVED** (previously OOM-killed)
  - Full 5m run, 0 restarts, 80 goroutines, 20MB heap post-load
  - Login success: 2.44% (2,096/85,832) — correctly rate-limited by auth limiter
  - Successful login latency: avg=48.88ms p95=88.6ms — Argon2id verify time, expected
  - Overall p95=5.07ms — excellent
- [x] 9.4 all_endpoints @ 500 VUs: **100% checks pass** (180,729/180,729)
  - p95=5.03ms, same as before — no regression from semaphore
- [x] 9.5 realistic_usage @ 500 VUs: **Now works!** (previously crashed with config error)
  - 84.94% checks pass, 34,039 complete scenario iterations
  - p95=6.73ms — excellent for mixed workload
- [x] 9.6 Prometheus metrics confirm auth rate limiter now fires:
  - **auth_allowed=11,874** / **auth_blocked=468,009** (was 0/0 before fix!)
  - global_allowed=330,535 / global_blocked=143,252
  - Auth limiter correctly protects login/MFA/password-reset at 1 req/s per IP
- [x] 9.7 Post-load resource check: 86 goroutines, 40MB heap, 16 PGX conns — stable, no leaks
- [x] 9.8 DB query latency under load: SELECT avg=2.27ms, INSERT avg=9.83ms, UPDATE avg=10ms
  - Previously SELECT degraded to 600-900ms during OOM — now stable under load
- [x] 9.9 Login avg latency dropped from 44.8ms → 1.70ms (aggregate) because auth rate limiter
  blocks most attempts before reaching Argon2id — exactly the intended behavior
