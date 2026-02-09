# MASTER_QUESTIONS Implementation Progress

**Date**: 2026-02-07
**Branch**: develop

## Completed (12/12) - ALL DONE

| Q# | Task | Commit | Notes |
|----|------|--------|-------|
| Q19 | IP/UA/fingerprint extraction | `7e8b54bd` | Wired existing `middleware.GetRequestMetadata(ctx)` to Login and ForgotPassword handlers |
| Q22 | Library stats from repo | `07fcd623` | Replaced stub with `repo.CountMovies()`. TV show stats exist in tvshow repo but not API-wired yet |
| Q23 | Cleanup job logic | `244d55bf` | Full rewrite: `AuthCleanupRepository` interface, 7 cleanup methods, target routing, `ScheduleCleanup()` helper |
| Q18 | SendGrid email provider | `5e9a9373` | Direct HTTP POST to SendGrid v3 API (no SDK). Tests use httptest server |
| Q16 | OIDC redirect fix | `f8f3bfd6` | Changed from broken 302 (no Location header) to 200 JSON with `auth_url` (SPA pattern). OpenAPI spec updated + ogen regenerated |
| Q24 | MaxAttempts configurable | `6c351564` | Added `jobs.max_attempts` to `JobsConfig`, default 25, removed hardcoded value from module.go |
| Q25 | Stale comments cleanup | `5a4a84df` | Removed misleading "Placeholder implementations" and "TODO: Implement all" from movie repo (methods are implemented) |
| Q20 | Notification settings handler | `bc08fb9f` | Marshal email/push/digest notification ogen types to `json.RawMessage` in UpdateUserPreferences. Added `prefsToOgenResponse` helper for both GET/PUT |
| Q26 | Profile visibility enforcement | `f32d2a79` | GetUserById now checks `profile_visibility` preference. Private/friends profiles return 404 for other users |
| Q21 | Search reindex River job | `b8f0d9a1` | Replaced stub with actual `riverClient.Insert()` using existing `MovieSearchIndexWorker` with reindex operation |
| Q17 | TV show search indexing | `6cf8f834` | Full Typesense integration: schema, `TVShowSearchService`, `SearchIndexWorker` with real implementation, fx wiring, 19 tests |
| Q15 | WebAuthn MFA verification | `c9ce54e4` | 7 OpenAPI endpoints, ogen regen, handler implementation with jx.Raw protocol bridge, delegation from main Handler, 22 tests |

## Bugs Found During Implementation

1. **OIDC 302 redirect was non-functional** (Q16): ogen generated empty struct for 302 response, no way to set Location header. Fixed by switching to JSON response.

2. **Pre-existing API test failure**: `gin_trgm_ops` operator class missing in test database - the `pg_trgm` extension is not installed. This causes migration tests to fail. Not related to any Q-item changes.

3. **CleanupWorker was not registered**: The auth `CleanupWorker` in `internal/infra/jobs/` is defined and tested but never registered with `river.AddWorker()` in any module. It needs to be wired into the fx dependency graph to actually run.

4. **GetUserPreferences/UpdateUserPreferences were missing notification settings in response**: Both handlers only returned non-notification fields (profile_visibility, theme, etc.) but omitted email/push/digest notification settings from the response JSON.

## Files Modified This Session

- `internal/api/handler.go` - Q19, Q20, Q26, Q15 (WebAuthn delegation)
- `internal/api/handler_mfa.go` - Q15 (WebAuthn handlers + helpers)
- `internal/api/handler_mfa_test.go` - Q15 (22 tests: not-configured, unauthorized, helpers, delegation)
- `internal/api/handler_oidc.go` - Q16
- `internal/api/handler_search.go` - Q21
- `internal/api/server.go` - Q15 (WebAuthnService in ServerParams)
- `internal/api/ogen/` - Q16, Q15 (regenerated)
- `api/openapi/openapi.yaml` - Q16, Q15 (7 WebAuthn endpoints + 6 schemas)
- `internal/infra/jobs/cleanup_job.go` - Q23
- `internal/infra/jobs/cleanup_job_test.go` - Q23
- `internal/infra/jobs/module.go` - Q24
- `internal/config/config.go` - Q24
- `internal/service/email/service.go` - Q18
- `internal/service/email/service_test.go` - Q18
- `internal/service/search/tvshow_schema.go` - Q17 (new)
- `internal/service/search/tvshow_service.go` - Q17 (new)
- `internal/service/search/tvshow_service_test.go` - Q17 (new, 19 tests)
- `internal/service/search/module.go` - Q17 (added TVShowSearchService)
- `internal/content/tvshow/jobs/jobs.go` - Q17 (real SearchIndexWorker)
- `internal/content/tvshow/jobs/module.go` - Q17 (search dependency)
- `internal/content/movie/library_service.go` - Q22
- `internal/content/movie/library_service_test.go` - Q22
- `internal/content/movie/repository_postgres.go` - Q25
