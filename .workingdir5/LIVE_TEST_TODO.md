# Live Test Coverage TODO

**Created**: 2026-02-07
**Updated**: 2026-02-08
**Stack**: docker-compose.dev.yml (all services healthy)
**Test file**: tests/live/smoke_test.go

## Final Results

**319 subtests PASSING** | 0 failures | ~125s runtime

## Bugs Found & Fixed

| Bug | Severity | Fix |
|-----|----------|-----|
| Image proxy nil pointer panic (crash) | CRITICAL | Nil check in `GetProxiedImage` handler |
| Image service never wired in fx app | CRITICAL | Added `image.Module` to app, made non-optional |
| Webhook job kinds not registered in River | HIGH | Added worker constructors + `registerWorkers` in radarr/sonarr modules |
| Metadata returns 500 when no providers configured | HIGH | Handle `ErrNoProviders` in all 7 metadata handlers |
| GetUserById returns 500 instead of 404 | MEDIUM | Don't return both typed response AND error with ogen |
| Webhook worker panic on nil Movie pointer | MEDIUM | Guard `args.Payload.Movie` before accessing `.ID` |
| Webhook handler panic on nil syncService | MEDIUM | Guard `h.syncService` in `HandleWebhook` |

## Coverage Summary

| Category | Tested | Total | Coverage |
|----------|--------|-------|----------|
| Health/Infra (9) | 9 | 9 | 100% |
| Auth (11) | 11 | 11 | 100% |
| Auth Lifecycle (7) | 7 | 7 | 100% |
| Sessions (7) | 7 | 7 | 100% |
| API Keys (7) | 7 | 7 | 100% |
| User Settings (5) | 5 | 5 | 100% |
| Preferences (3) | 3 | 3 | 100% |
| MFA (14) | 14 | 14 | 100% |
| Movies list/search (8) | 8 | 8 | 100% |
| Movies detail/subresources (16) | 16 | 16 | 100% |
| TV Shows list/search (7) | 7 | 7 | 100% |
| TV Shows detail/subresources (19) | 19 | 19 | 100% |
| Collections (2) | 2 | 2 | 100% |
| Playback (4) | 4 | 4 | 100% |
| Metadata (7) | 7 | 7 | 100% |
| Images (1) | 1 | 1 | 100% |
| Search (6) | 6 | 6 | 100% |
| Settings server (4) | 4 | 4 | 100% |
| Users/Misc (4) | 4 | 4 | 100% |
| RBAC (7) | 7 | 7 | 100% |
| Libraries (8) | 8 | 8 | 100% |
| Admin Activity (4) | 4 | 4 | 100% |
| OIDC Admin (8) | 8 | 8 | 100% |
| OIDC User (2) | 2 | 2 | 100% |
| Integrations (8) | 8 | 8 | 100% |
| Auth Edge Cases (2) | 2 | 2 | 100% |
| Unauthed Access (41) | 41 | 41 | 100% |
| **TOTAL** | **319** | **319** | **100%** |

## Test Functions (34 total)

1. TestLive_Infrastructure — health checks, DB, Typesense, Dragonfly, Casbin
2. TestLive_AuthEndpoints — register, login, refresh, forgot/reset, verify
3. TestLive_Chain_AuthLifecycle — profile, password change, token refresh, logout
4. TestLive_Chain_SessionManagement — list, current, logout all
5. TestLive_Chain_APIKeyLifecycle — create, list, auth-with, revoke
6. TestLive_Chain_UserSettings — set, get, update, delete
7. TestLive_Chain_MFA — status, setup TOTP, backup codes, webauthn
8. TestLive_Chain_Preferences — get, update, verify
9. TestLive_ContentEndpoints_EmptyLibrary — all list/search endpoints
10. TestLive_SearchInfrastructure — Typesense search, autocomplete, facets
11. TestLive_TVShowSearchInfrastructure — TV search, autocomplete, facets
12. TestLive_AdminEndpoints — users, libraries, activity, server settings, RBAC
13. TestLive_UnauthenticatedAccess — all public vs protected endpoints
14. TestLive_AdminOnlyEndpoints — RBAC enforcement on admin endpoints
15. TestLive_SecurityHeaders — CORS, rate limiting, content-type
16. TestLive_ConnectionSafety — long requests, concurrent access, DB pooling
17. TestLive_MovieDetailSubresources — files, cast, crew, genres, collection, similar, progress
18. TestLive_TVShowDetailSubresources — seasons, episodes, cast, crew, genres, networks, watch-stats
19. TestLive_TVShowEpisodeSubresources — season episodes, episode files, episode progress
20. TestLive_Collections — collection detail, collection movies
21. TestLive_PlaybackSessions — create, get, heartbeat, end session
22. TestLive_MFACompleteFlows — enable/disable, TOTP verify/delete, backup codes, WebAuthn
23. TestLive_SessionsGranular — refresh, delete current, delete specific
24. TestLive_AdminActivityGranular — per-user, per-resource, action types
25. TestLive_LibraryGranular — create, update, scans, permissions grant/revoke
26. TestLive_RBACGranular — role CRUD, permission update, policy CRUD
27. TestLive_SettingsGranular — server setting set/get, admin-only enforcement
28. TestLive_AuthEdgeCases — resend verification (authed + unauthed)
29. TestLive_MetadataEndpoints — search/lookup for movie, TV, collection, season, episode
30. TestLive_OIDCAdminCRUD — create, get, update, disable, enable, set default, delete
31. TestLive_OIDCUserEndpoints — list linked, unlink nonexistent
32. TestLive_IntegrationEndpoints — radarr/sonarr sync, profiles, folders, webhooks
33. TestLive_MiscEndpoints — get user by ID, image proxy, avatar upload
34. TestLive_UnauthenticatedAccess_Extended — 41 endpoints verified require auth
