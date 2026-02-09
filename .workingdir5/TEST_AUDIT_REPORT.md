# Comprehensive Test & Code Audit Report

**Created**: 2026-02-07
**Purpose**: Pre-frontend preparation — find all bugs, gaps, and missing test coverage
**Method**: Full coverage analysis + 6 targeted deep-scan agents + manual verification
**Principle**: "No shortcuts, no half measures — check code when tests fail, don't fix tests to pass bugs"

---

## 1. Test Coverage Summary

All 48 packages pass. Zero failures. `go vet` and `golangci-lint` both report 0 issues.

### Coverage by Package (sorted by coverage)

#### Critical (0% — no tests at all)
| Package | Coverage | Notes |
|---------|----------|-------|
| `infra/raft` | 0.0% | Leader election — no tests |
| `playback/jobs` | 0.0% | Cleanup job — no tests |
| `playback/subtitle` | 0.0% | Subtitle extraction — no tests |
| `playback/playbackfx` | 0.0% | fx module — no tests |
| `service/metadata` | 0.0% | Service interface — no tests |
| `service/metadata/adapters/*` | 0.0% | Movie/TV adapters — no tests |
| `service/metadata/providers/*` | 0.0% | TMDb/TVDb providers — no tests |
| `service/metadata/jobs` | 0.0% | Metadata jobs — no tests |
| `service/metadata/metadatafx` | 0.0% | fx module — no tests |
| `content/movie/db` | 0.0% | Generated sqlc (expected) |
| `content/tvshow/db` | 0.0% | Generated sqlc (expected) |
| `content/qar/db` | 0.0% | Generated sqlc placeholder (expected) |
| `infra/database/db` | 0.0% | Generated sqlc (expected) |
| `api/ogen` | 0.0% | Generated ogen (expected) |

#### Low (<30% — needs significant work)
| Package | Coverage | Notes |
|---------|----------|-------|
| `integration/sonarr` | 2.3% | Only client tests, no sync/webhook tests |
| `content/tvshow/jobs` | 2.4% | Only module test, no worker tests |
| `content/movie/moviejobs` | 8.7% | Only module test, no worker tests |
| `api` | 20.1% | 6 handler files completely untested |
| `playback` | 24.7% | Service layer partially tested |

#### Medium (30-60% — needs targeted additions)
| Package | Coverage | Notes |
|---------|----------|-------|
| `service/search` | 37.9% | Movie search tested, TV show search missing |
| `integration/radarr` | 42.9% | Client tested, sync service partially |
| `infra/observability` | 42.1% | Metrics registration tested, collectors not |
| `testutil` | 48.2% | Test helpers themselves partially tested |
| `service/storage` | 50.9% | Local storage tested, S3 not |
| `infra/jobs` | 51.5% | River client tested, queues partially |
| `playback/hls` | 53.4% | Manifest tested, handler partially |
| `infra/search` | 53.5% | Client tested, indexing partially |

#### Good (60-80% — minor gaps)
| Package | Coverage | Notes |
|---------|----------|-------|
| `service/auth` | 62.9% | Core flows tested, edge cases missing |
| `service/email` | 62.3% | SMTP + SendGrid tested |
| `service/library` | 65.9% | CRUD tested, permissions partially |
| `infra/cache` | 65.7% | L1 + L2 tested |
| `service/mfa` | 66.8% | TOTP + backup codes tested |
| `playback/transcode` | 67.5% | Profiles + decisions tested |
| `service/oidc` | 68.4% | Provider management tested |
| `service/session` | 77.5% | All methods tested |
| `service/rbac` | 77.7% | Casbin integration tested |
| `infra/logging` | 78.7% | Logger creation tested |
| `infra/image` | 78.9% | Resize/crop tested |
| `infra/database` | 79.5% | Migrations + pool tested |
| `api/middleware` | 80.4% | Rate limiting + request ID tested |

#### Excellent (80%+)
| Package | Coverage | Notes |
|---------|----------|-------|
| `service/activity` | 81.2% | All 13 methods tested |
| `service/user` | 81.2% | All 26 methods tested |
| `content/shared/metadata` | 82.7% | Provider chain tested |
| `content/movie/adapters` | 82.6% | TMDb adapter tested |
| `infra/health` | 83.9% | All probes tested |
| `service/settings` | 84.5% | Server + user settings tested |
| `crypto` | 84.7% | Argon2id + HMAC tested |
| `content/shared/scanner` | 86.1% | File scanning tested |
| `content/tvshow/adapters` | 89.7% | TVDb adapter tested |
| `service/apikeys` | 90.3% | All methods tested |
| `service/notification` | 92.6% | Dispatcher tested |
| `content/shared/matcher` | 92.5% | File matching tested |
| `errors` | 100% | Error types tested |
| `util` | 100% | Utility functions tested |
| `validate` | 100% | Validation tested |
| `version` | 100% | Version info tested |
| `content/shared/jobs` | 100% | Job types tested |
| `content/shared/library` | 100% | Library types tested |
| `service/notification/agents` | 80.1% | Discord/Gotify/Webhook tested |

---

## 2. Untested API Handlers (CRITICAL)

6 handler files with zero test coverage (46 handler methods total):

### handler_library.go — 10 methods
1. `ListLibraries` — GET /api/v1/libraries
2. `CreateLibrary` — POST /api/v1/libraries
3. `GetLibrary` — GET /api/v1/libraries/{libraryId}
4. `UpdateLibrary` — PUT /api/v1/libraries/{libraryId}
5. `DeleteLibrary` — DELETE /api/v1/libraries/{libraryId}
6. `TriggerLibraryScan` — POST /api/v1/libraries/{libraryId}/scan
7. `ListLibraryScans` — GET /api/v1/libraries/{libraryId}/scans
8. `ListLibraryPermissions` — GET /api/v1/libraries/{libraryId}/permissions
9. `GrantLibraryPermission` — POST /api/v1/libraries/{libraryId}/permissions
10. `RevokeLibraryPermission` — DELETE /api/v1/libraries/{libraryId}/permissions/{userId}

### handler_oidc.go — 14 methods (security-critical)
1. `ListOIDCProviders` — GET /api/v1/auth/oidc/providers
2. `OidcAuthorize` — GET /api/v1/auth/oidc/authorize
3. `OidcCallback` — GET /api/v1/auth/oidc/callback
4. `ListUserOIDCLinks` — GET /api/v1/auth/oidc/links
5. `InitOIDCLink` — GET /api/v1/auth/oidc/links/{provider}
6. `UnlinkOIDCProvider` — DELETE /api/v1/auth/oidc/links/{provider}
7. `AdminListOIDCProviders` — GET /api/v1/admin/auth/oidc/providers
8. `AdminCreateOIDCProvider` — POST /api/v1/admin/auth/oidc/providers
9. `AdminGetOIDCProvider` — GET /api/v1/admin/auth/oidc/providers/{providerId}
10. `AdminUpdateOIDCProvider` — PUT /api/v1/admin/auth/oidc/providers/{providerId}
11. `AdminDeleteOIDCProvider` — DELETE /api/v1/admin/auth/oidc/providers/{providerId}
12. `AdminEnableOIDCProvider` — POST /api/v1/admin/auth/oidc/providers/{providerId}/enable
13. `AdminDisableOIDCProvider` — POST /api/v1/admin/auth/oidc/providers/{providerId}/disable
14. `AdminSetDefaultOIDCProvider` — POST /api/v1/admin/auth/oidc/providers/{providerId}/default

### handler_metadata.go — 8 methods
1. `SearchMoviesMetadata` — GET /api/v1/metadata/movies/search
2. `GetMovieMetadata` — GET /api/v1/metadata/movies/{tmdbId}
3. `GetProxiedImage` — GET /api/v1/metadata/images/{type}/{size}/{path}
4. `GetCollectionMetadata` — GET /api/v1/metadata/collections/{tmdbId}
5. `SearchTVShowsMetadata` — GET /api/v1/metadata/tvshows/search
6. `GetTVShowMetadata` — GET /api/v1/metadata/tvshows/{tmdbId}
7. `GetSeasonMetadata` — GET /api/v1/metadata/tvshows/{tmdbId}/seasons/{seasonNumber}
8. `GetEpisodeMetadata` — GET /api/v1/metadata/tvshows/{tmdbId}/seasons/{seasonNumber}/episodes/{episodeNumber}

### handler_search.go — 4 methods
1. `SearchLibraryMovies` — GET /api/v1/search/movies
2. `AutocompleteMovies` — GET /api/v1/search/movies/autocomplete
3. `GetSearchFacets` — GET /api/v1/search/facets
4. `ReindexSearch` — POST /api/v1/search/reindex

### handler_sonarr.go — 6 methods
1. `AdminGetSonarrStatus` — GET /api/v1/admin/integrations/sonarr/status
2. `AdminTriggerSonarrSync` — POST /api/v1/admin/integrations/sonarr/sync
3. `AdminGetSonarrQualityProfiles` — GET /api/v1/admin/integrations/sonarr/quality-profiles
4. `AdminGetSonarrRootFolders` — GET /api/v1/admin/integrations/sonarr/root-folders
5. `HandleSonarrWebhook` — POST /api/v1/webhooks/sonarr
6. Helper: `convertSonarrWebhookPayload`

### handler_playback.go — 3 methods
1. `StartPlaybackSession` — POST /api/v1/playback/sessions
2. `GetPlaybackSession` — GET /api/v1/playback/sessions/{sessionId}
3. `StopPlaybackSession` — DELETE /api/v1/playback/sessions/{sessionId}

### Tested Handlers (for reference)
- handler_activity.go — 5/5 methods (good coverage)
- handler_apikeys.go — 5/5 methods (excellent)
- handler_mfa.go — 17 methods (stub-only: tests only check 501 when service nil)
- handler_radarr.go — 5/5 methods (good)
- handler_rbac.go — 11/11 methods (excellent)
- handler_session.go — 6/6 methods (excellent)

---

## 3. Service Layer Findings

### Missing Error Path Tests (31 total)

| Service | Method | Gap |
|---------|--------|-----|
| apikeys | ValidateKey | Background goroutine silently fails on last_used_at update |
| apikeys | UpdateScopes | Repo call error path untested |
| apikeys | RevokeKey | Error path untested |
| auth | Register | Transaction rollback not tested |
| auth | Login | Account lockout edge cases |
| email | SendVerificationEmail | enabled=false path |
| email | SendPasswordResetEmail | SMTP timeout/TLS errors |
| email | SendWelcomeEmail | enabled=false path |
| email | N/A | SendGrid non-401 status codes |
| library | Update | Name conflict detection edge |
| library | TriggerScan | Repo error at line 286-293 |
| library | GrantPermission | Comprehensive error scenarios |
| mfa | GetStatus | Non-ErrNoRows error handling |
| mfa | SetRememberDevice | Transaction atomicity |
| notification | Dispatch | Async failure handling, goroutine leak on ctx cancel |
| notification | Dispatch | Error channel draining |
| oidc | AddProvider | OIDC discovery failures |
| oidc | ExchangeAuthorizationCode | Token exchange errors |
| oidc | encryptSecret/decryptSecret | Comprehensive encryption failure paths |
| rbac | AddPolicy | Already-exists vs true-failure distinction |
| rbac | RemovePolicy | Same pattern |
| rbac | List* methods | Error path testing |
| search | UpdateMovie | Non-"not found" error types |
| search | BulkIndex | Partial success (3 of 5 docs fail) |
| session | RefreshSession | Max attempts enforcement |
| session | RevokeAllUserSessionsExcept | Current session not found |
| settings | SetServerSetting | Update vs upsert race |
| settings | SetUserSettingsBulk | Partial failure handling |
| storage | Store | Directory creation failure |
| storage | Get | Permission errors |
| storage | Delete | Permission denied case |

### Bugs Found in Service Layer (19 total)

#### HIGH PRIORITY
1. **auth: fmt.Printf instead of logger** — Lines 135, 262, 287, 395, 463, 531 in `service/auth/service.go`. Multiple `fmt.Printf()` calls used for error logging instead of `s.activityLogger.LogFailure()` or zap logger. Error messages go to stdout instead of structured logs.

2. **notification: Missing timeout for agent.Send()** — `service/notification/dispatcher.go:129-138`. Goroutines calling agent.Send() have no per-agent timeout, could block indefinitely.

3. **storage: file.Close() error ignored** — `service/storage/local.go:74-85`. Deferred Close() ignores error if write fails. Should Sync() before returning.

#### MEDIUM PRIORITY
4. **apikeys: Background goroutine silent failure** — `service/apikeys/service.go:185-195`. ValidateKey() fires goroutine to update last_used_at; context timeout causes silent failure.

5. **session: Max sessions not enforced** — `service/session/service.go:56-62`. When maxPerUser exceeded, logs warning but doesn't revoke oldest session (has comment "Optionally revoke oldest session here").

6. **settings: Non-atomic upsert** — `service/settings/service.go:157-158`. SetServerSetting tries update first, then insert if not found. Race condition if two requests create same key simultaneously.

7. **mfa: Silent error swallowing** — `service/mfa/manager.go:74-76`. GetStatus() ignores settings query error without logging.

8. **rbac: Inconsistent error on duplicate** — `service/rbac/service.go:71-76`. AddPolicy warns but doesn't distinguish "already exists" from real failure in return value.

9. **storage: No path length validation** — sanitizeKey() doesn't validate against very long paths or excessive nesting.

10. **oidc: Encryption failure returns nil provider** — `service/oidc/service.go:99`. encryptSecret() failure could return nil provider with error that caller might not handle.

---

## 4. Bugs & Dead Code

### Confirmed Bugs

#### ~~BUG-1: River Workers Not Wired to Config~~ — VERIFIED FIXED
**Status**: FIXED in previous session (DEAD-1). `module.go:44-52` reads `cfg.Jobs.MaxWorkers`, `FetchPollInterval`, `RescueStuckJobsAfter` and passes them to River. All 11 workers registered via `river.AddWorker()` in their respective modules, all modules wired in `app/module.go`.

#### BUG-2: Notification Service Not Wired
**File**: `internal/service/notification/` — dispatcher.go exists, 4 agents exist (Discord, Email, Gotify, Webhook), but NO `module.go` and NOT in `app/module.go`.
**Impact**: Notification dispatchers never initialized. Notifications won't send.

#### BUG-3: Auth Service Uses fmt.Printf for Errors
**File**: `internal/service/auth/service.go`
**Lines**: 135, 262, 287, 395, 463, 531
**Impact**: Error messages go to stdout instead of structured logging. Won't appear in log aggregation systems.

### Stubs & Placeholders (Intentional)

| Item | File | Status |
|------|------|--------|
| Windows media prober | `content/movie/mediainfo_windows.go` | Platform stub — intentional |
| QAR module | `content/qar/db/placeholder.sql.go` | v0.3.0+ planned |
| Movie DB queries | `content/movie/db/placeholder.sql.go` | sqlc requirement |
| TV Show DB queries | `content/tvshow/db/placeholder.sql.go` | sqlc requirement |

### Remaining TODOs in Code

Only 3 files have TODO/FIXME comments (excluding tests and generated code):
1. `internal/api/localization.go:15` — User language from settings (partially implemented via GetMetadataLanguage)
2. `internal/service/auth/service_integration_test.go` — Test for timing attack
3. `internal/infra/logging/logging_test.go` — Minor test improvement

### Skipped Tests

| Test | File | Reason |
|------|------|--------|
| Health readiness without DB | `tests/integration/health_test.go:81` | Health service doesn't detect DB failures (documented behavior) |
| Typesense search health | `tests/integration/search/search_test.go:47` | typesense-go v2 client health endpoint issue |

---

## 5. Security Findings

All tracked in `.workingdir/TODO_A7_SECURITY_FIXES.md`. Status as of this audit:

| ID | Issue | Status |
|----|-------|--------|
| A7.1 | Missing transaction boundaries | FIXED — tx added to Register() and VerifyEmail() |
| A7.2 | Login timing attack | FIXED — dummy hash constant-time comparison |
| A7.3 | Goroutine leak in notification | FIXED — WaitGroup + stopCh + Close() |
| A7.4 | Password reset info disclosure | FIXED — returns only error, silent success |
| A7.5 | No rate limiting for Argon2id | FIXED — account lockout implemented |
| A7.6 | context.Background() in goroutines | FIXED — all have timeouts |

---

## 6. Infrastructure Gaps

### Packages With Zero Test Coverage (non-generated)

1. **infra/raft** — Leader election implementation exists (170+ lines) but zero tests
2. **playback/jobs** — Cleanup worker exists but zero tests
3. **playback/subtitle** — Subtitle extraction exists but zero tests
4. **service/metadata** — Service interface exists but zero tests
5. **service/metadata/adapters/** — Movie/TV adapters exist but zero tests
6. **service/metadata/providers/** — TMDb/TVDb providers exist but zero tests
7. **service/metadata/jobs** — Metadata refresh jobs exist but zero tests

### API Test Infrastructure

The API test suite uses 15 embedded postgres instances across `server_test.go` (12) and `handler_test.go` (3). These could be consolidated to a shared instance similar to the database package pattern.

---

## 7. Weighted Priority Matrix

### P0 — Must Fix (Bugs That Affect Functionality)

| # | Issue | Package | Effort |
|---|-------|---------|--------|
| 1 | fmt.Printf → logger in auth (9 instances) | service/auth | 30min |
| 2 | Wire notification service into app (no module.go) | service/notification | 1h |
| 3 | Notification agent.Send() no timeout | service/notification | 30min |

### P1 — High Priority (Critical Test Gaps)

| # | Issue | Package | Effort |
|---|-------|---------|--------|
| 4 | handler_library.go tests (10 methods) | api | 4-6h |
| 5 | handler_oidc.go tests (14 methods, security) | api | 6-8h |
| 6 | handler_metadata.go tests (8 methods) | api | 3-4h |
| 7 | handler_search.go tests (4 methods) | api | 2-3h |
| 8 | handler_sonarr.go tests (6 methods) | api | 2-3h |
| 9 | handler_playback.go tests (3 methods) | api | 2h |
| 10 | MFA handler tests beyond 501 stubs | api | 4-6h |
| 11 | Metadata service + providers tests | service/metadata | 6-8h |

### P2 — Medium Priority (Coverage Improvements)

| # | Issue | Package | Effort |
|---|-------|---------|--------|
| 12 | Sonarr integration sync/webhook tests | integration/sonarr | 3-4h |
| 13 | Movie/TV show worker tests | content/*/jobs | 4-6h |
| 14 | Raft leader election tests | infra/raft | 2-3h |
| 15 | Playback subtitle + jobs tests | playback/* | 3-4h |
| 16 | Observability metrics collector tests | infra/observability | 2h |
| 17 | Storage S3 tests | service/storage | 2h |

### P3 — Nice To Have (Polish)

| # | Issue | Package | Effort |
|---|-------|---------|--------|
| 18 | API test DB consolidation (15 instances) | api | 4-6h |
| 19 | Service layer error path tests (31 gaps) | service/* | 8-12h |
| 20 | Non-atomic settings upsert | service/settings | 1h |
| 21 | Session max-per-user enforcement | service/session | 1h |

---

## 8. Implementation Plan

### Phase 1: Bug Fixes (1 day)
- [ ] Fix auth service fmt.Printf → logger, 9 instances (P0-1)
- [ ] Wire notification service into app — create module.go (P0-2)
- [ ] Add timeout to notification agent.Send() (P0-3)
- [ ] Fix non-atomic settings upsert (P3-20)

### Phase 2: Critical Handler Tests (3-4 days)
- [ ] handler_library.go tests (P1-4)
- [ ] handler_oidc.go tests (P1-5)
- [ ] handler_metadata.go tests (P1-6)
- [ ] handler_search.go tests (P1-7)
- [ ] handler_sonarr.go tests (P1-8)
- [ ] handler_playback.go tests (P1-9)
- [ ] Upgrade MFA handler tests beyond stubs (P1-10)

### Phase 3: Service & Provider Tests (2-3 days)
- [ ] Metadata service + providers tests (P1-11)
- [ ] Sonarr integration tests (P2-12)
- [ ] Movie/TV show worker tests (P2-13)
- [ ] Playback subtitle + jobs tests (P2-15)

### Phase 4: Infrastructure Tests (1-2 days)
- [ ] Raft leader election tests (P2-14)
- [ ] Observability collector tests (P2-16)
- [ ] Storage S3 tests (P2-17)

### Phase 5: Polish & Error Paths (2-3 days)
- [ ] API test DB consolidation (P3-18)
- [ ] Service layer error path tests (P3-19)
- [ ] Session max enforcement (P3-21)

### Target Coverage After All Phases (minimum 80%)

| Package | Current | Target |
|---------|---------|--------|
| api | 20.1% | 80%+ |
| integration/sonarr | 2.3% | 80%+ |
| content/*/jobs | 2-9% | 80%+ |
| playback/* | 0-25% | 80%+ |
| service/metadata | 0% | 80%+ |
| infra/raft | 0% | 80%+ |
| infra/jobs | 51.5% | 80%+ |
| infra/observability | 42.1% | 80%+ |
| service/search | 37.9% | 80%+ |
| service/storage | 50.9% | 80%+ |
| Overall average | ~55% | ~80%+ |
