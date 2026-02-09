# Planning Documentation Analysis

> Comprehensive audit of `docs/dev/design/planning/` against actual codebase and active design docs.
> Generated: 2026-02-06

---

## Table of Contents

- [Executive Summary](#executive-summary)
- [1. Items Marked TODO/Planned That Are Actually DONE](#1-items-marked-todoplanned-that-are-actually-done)
  - [1.1 ROADMAP.md](#11-roadmapmd)
  - [1.2 TODO_v0.1.0.md](#12-todo_v010md)
  - [1.3 TODO_v0.2.0.md (Entire Milestone)](#13-todo_v020md-entire-milestone)
  - [1.4 TODO_v0.3.0.md (Partial)](#14-todo_v030md-partial)
  - [1.5 TODO_v0.4.0.md (Partial)](#15-todo_v040md-partial)
  - [1.6 TODO_v0.8.0.md (Partial)](#16-todo_v080md-partial)
  - [1.7 TODO_v0.9.0.md (Partial - QAR Schema)](#17-todo_v090md-partial---qar-schema)
- [2. Broken Design Doc References](#2-broken-design-doc-references)
  - [2.1 TODO_v0.2.0.md](#21-todo_v020md)
  - [2.2 TODO_v0.3.0.md](#22-todo_v030md)
  - [2.3 TODO_v0.4.0.md](#23-todo_v040md)
  - [2.4 TODO_v0.5.0.md](#24-todo_v050md)
  - [2.5 TODO_v0.6.0.md](#25-todo_v060md)
  - [2.6 TODO_v0.7.0.md](#26-todo_v070md)
  - [2.7 TODO_v0.8.0.md](#27-todo_v080md)
  - [2.8 TODO_v0.9.0.md](#28-todo_v090md)
  - [2.9 TODO_v1.0.0.md](#29-todo_v100md)
  - [2.10 ROADMAP.md and INDEX.md](#210-roadmapmd-and-indexmd)
- [3. Gap Analysis - Features Not Covered by Any Active Design Doc](#3-gap-analysis---features-not-covered-by-any-active-design-doc)
- [4. Stale Status Markers](#4-stale-status-markers)
  - [4.1 ROADMAP.md](#41-roadmapmd)
  - [4.2 INDEX.md](#42-indexmd)
  - [4.3 TODO_v0.2.0.md](#43-todo_v020md)
  - [4.4 TODO_v0.3.0.md](#44-todo_v030md)
  - [4.5 TODO_v0.4.0.md](#45-todo_v040md)
- [5. Missing From Planning - Implemented But Unplanned](#5-missing-from-planning---implemented-but-unplanned)
- [6. Summary of Required Changes Per File](#6-summary-of-required-changes-per-file)

---

## Executive Summary

The planning docs are **severely out of date**. The project has advanced well beyond v0.1.x (where the roadmap claims "Current Version: v0.1.3") and has implemented the majority of v0.2.0 (Core Services), significant portions of v0.3.0 (MVP/Movies), substantial parts of v0.4.0 (TV Shows), and some items from later milestones. Here is the high-level reality:

| Area | Planning Says | Code Reality |
|------|--------------|-------------|
| Current version | v0.1.3 (Skeleton) | Well into v0.3.0+ territory |
| v0.2.0 (Core) status | "Not Started" | **~95% complete** - all 10+ services implemented |
| v0.3.0 (Movies) status | "Not Started" | **~70% backend complete** (module, metadata, Radarr, search) |
| v0.4.0 (Shows) status | "Not Started" | **~60% backend complete** (full module, Sonarr, TVDb) |
| Services implemented | 0 of 10 | 15 services: auth, user, session, rbac, apikeys, oidc, mfa, settings, activity, library, metadata, search, email, notification, storage |
| Content modules | 0 | 3: movie (full), tvshow (full), qar (schema) |
| Migrations | 3 | 32 (64 files) |
| River workers | 0 | 17 workers in 5 queues |
| OpenAPI endpoints | Health only | 153 handler methods |
| Integrations | 0 | Radarr, Sonarr, TMDb, TVDb (all implemented) |
| CI/CD workflows | 7-8 | 9 workflows |
| Design doc 00_SOURCE_OF_TRUTH.md | Referenced everywhere | **Deleted** during doc rewrite |

---

## 1. Items Marked TODO/Planned That Are Actually DONE

### 1.1 ROADMAP.md

**File**: `docs/dev/design/planning/ROADMAP.md`

| Line/Section | Planning Says | Actual State |
|---|---|---|
| Line 69 | "Current Version: v0.1.3" | Should be updated - code is well beyond v0.2.0 |
| Line 108 (v0.2.0) | "Not Started" | **~95% complete** - all core services implemented |
| Line 109 (v0.3.0) | "Not Started" | **~70% backend complete** |
| Line 110 (v0.4.0) | "Not Started" | **~60% backend complete** |
| Lines 235-248 (v0.2.0 deliverables) | All unchecked `[ ]` | Almost all are implemented (see v0.2.0 section below) |
| Lines 269-270 (v0.3.0 Movie Module) | Unchecked | Movie module fully implemented |
| Lines 271-274 (v0.3.0 Metadata/Search/Radarr) | Unchecked | TMDb, Typesense, Radarr all implemented |
| Lines 306-309 (v0.4.0 TV Shows) | Unchecked | TV show module, TVDb, Sonarr all implemented |

### 1.2 TODO_v0.1.0.md

**File**: `docs/dev/design/planning/TODO_v0.1.0.md`

Items marked as "deferred to v0.2.0" that are now **DONE**:

| Line | Item | Status |
|---|---|---|
| Line 101 | `internal/infra/cache/module.go` - rueidis provider (deferred to v0.2.0) | **DONE** - `internal/infra/cache/module.go` exists with otter L1 + rueidis L2 |
| Line 102 | `internal/infra/search/module.go` - typesense provider (deferred to v0.2.0) | **DONE** - `internal/infra/search/module.go` exists |
| Line 103 | `internal/infra/jobs/module.go` - river provider (deferred to v0.2.0) | **DONE** - full River infrastructure with 5 queues |
| Line 216 | Dragonfly ping check (deferred to v0.2.0) | **DONE** - `internal/infra/health/checks.go` exists |
| Line 217 | Typesense health check (deferred to v0.2.0) | **DONE** - health checks implemented |
| Line 218 | River worker check (deferred to v0.2.0) | **DONE** - River client operational |
| Lines 233-235 | `cmd/revenge/migrate.go` subcommands (deferred) | **DONE** - `cmd/revenge/migrate.go` implements up/down/version/create |
| Line 246 | testcontainers-go Dragonfly (deferred to v0.2.0) | Needs verification |
| Line 247 | testcontainers-go Typesense (deferred to v0.2.0) | Needs verification |

### 1.3 TODO_v0.2.0.md (Entire Milestone)

**File**: `docs/dev/design/planning/TODO_v0.2.0.md`

**Status marker on line 34 says "Not Started" -- this is the most egregiously wrong status in the planning docs.** Nearly everything in this milestone is implemented.

#### Auth Service (Lines 50-89) -- DONE

| Item | Code Evidence |
|---|---|
| Database schema (auth_tokens, password_reset_tokens, email_verification_tokens) | Migrations 000008, 000009, 000010 |
| Repository interface + PostgreSQL implementation | `internal/service/auth/repository.go`, `repository_pg.go` |
| Service (login, logout, register, password reset, JWT, token refresh) | `internal/service/auth/service.go`, `jwt.go`, `mfa_integration.go` |
| Handler endpoints | `internal/api/handler.go` (ogen-generated covers auth endpoints) |
| Middleware (JWT validation, token extraction) | `internal/api/middleware/` |
| Tests | `service_exhaustive_test.go`, `service_integration_test.go`, `jwt_test.go` |

#### User Service (Lines 92-134) -- DONE

| Item | Code Evidence |
|---|---|
| Database schema (users, user_preferences, user_avatars) | Migrations 000002, 000006, 000007 |
| Repository + PostgreSQL implementation | `internal/service/user/repository.go`, `repository_pg.go` |
| Service (profile, change password, preferences) | `internal/service/user/service.go`, `cached_service.go` |
| Tests | `service_test.go`, `service_unit_test.go`, `cached_service_test.go` |

#### Session Service (Lines 136-166) -- DONE

| Item | Code Evidence |
|---|---|
| Database schema (sessions) | Migration 000003 |
| Repository (PostgreSQL + cache) | `internal/service/session/repository.go`, `repository_pg.go`, `cached_service.go` |
| Service (create, validate, extend, revoke, list) | `internal/service/session/service.go` |
| Handler | `internal/api/handler_session.go` |
| Tests | `service_test.go`, `service_exhaustive_test.go`, `cached_service_test.go` |

#### RBAC Service (Lines 168-213) -- DONE

| Item | Code Evidence |
|---|---|
| Database schema (casbin_rule, fine-grained permissions) | Migrations 000011, 000027, 000029 |
| casbin-pgx-adapter | `internal/service/rbac/adapter.go` |
| Service (check permission, add/remove policy, roles) | `internal/service/rbac/service.go`, `cached_service.go`, `roles.go`, `permissions.go` |
| Handler | `internal/api/handler_rbac.go` |
| Default roles (admin, user, guest, moderator, legacy:read) | `internal/service/rbac/roles.go` |
| Tests | `service_test.go`, `service_unit_test.go`, `roles_test.go`, `cached_service_test.go` |

#### API Keys Service (Lines 215-246) -- DONE

| Item | Code Evidence |
|---|---|
| Database schema (api_keys) | Migration 000012 |
| Repository + service | `internal/service/apikeys/repository.go`, `service.go` |
| Handler | `internal/api/handler_apikeys.go` |
| Tests | `service_test.go`, `service_unit_test.go` |

#### OIDC Service (Lines 248-281) -- DONE

| Item | Code Evidence |
|---|---|
| Database schema (oidc_providers, oidc_user_links) | Migration 000013 |
| Repository + service | `internal/service/oidc/repository.go`, `service.go` |
| Handler | `internal/api/handler_oidc.go` |
| Admin endpoints (CRUD providers, enable/disable, set default) | Visible in ogen server interface |
| Tests | `service_test.go`, `service_unit_test.go` |

#### Settings Service (Lines 283-308) -- DONE

| Item | Code Evidence |
|---|---|
| Database schema (server_settings, user_settings) | Migrations 000004, 000005 |
| Repository + service | `internal/service/settings/repository.go`, `service.go`, `cached_service.go` |
| Tests | `service_test.go`, `service_unit_test.go`, `cached_service_test.go` |

#### Activity Service (Lines 310-337) -- DONE

| Item | Code Evidence |
|---|---|
| Database schema (activity_logs) | Migration 000014 |
| Repository + service | `internal/service/activity/repository.go`, `service.go`, `logger.go` |
| Handler | `internal/api/handler_activity.go` |
| River Job (ActivityCleanupWorker) | `internal/service/activity/cleanup.go` |
| Tests | `service_test.go`, `service_unit_test.go`, `service_additional_test.go` |

#### Library Service (Lines 339-372) -- DONE

| Item | Code Evidence |
|---|---|
| Database schema (libraries, library_paths, library_access) | Migration 000015 |
| Repository + service | `internal/service/library/repository.go`, `service.go`, `cached_service.go` |
| Handler | `internal/api/handler_library.go` |
| Cleanup worker | `internal/service/library/cleanup.go` |
| Tests | `service_unit_test.go`, `repository_pg_test.go` |

#### Health Service Enhancement (Lines 374-386) -- PARTIALLY DONE

| Item | Status |
|---|---|
| Service-level health checks | DONE - `internal/infra/health/checks.go` |
| River job queue health | DONE |
| Prometheus Metrics | DONE - `internal/infra/observability/metrics.go`, `middleware.go` |

#### PostgreSQL Integration Enhancement (Lines 388-401) -- DONE

| Item | Code Evidence |
|---|---|
| Connection pool metrics | `internal/infra/database/metrics.go` |
| Query logging (debug mode) | `internal/infra/database/logger.go` |
| sqlc queries (user, session, library, settings, activity, auth_tokens, apikeys, oidc, mfa) | `internal/infra/database/db/*.sql.go` |

#### Dragonfly/Redis Integration (Lines 403-419) -- DONE

| Item | Code Evidence |
|---|---|
| rueidis client | `internal/infra/cache/cache.go` |
| otter L1 cache | `internal/infra/cache/otter.go` |
| Cache module | `internal/infra/cache/module.go` |

#### River Job Queue Setup (Lines 421-436) -- DONE

| Item | Code Evidence |
|---|---|
| Client setup, worker pool, graceful shutdown | `internal/infra/jobs/river.go`, `module.go` |
| Queue configuration (5 queues: critical, high, default, low, bulk) | `internal/infra/jobs/queues.go` |
| CleanupJob | `internal/infra/jobs/cleanup_job.go` |
| EmailSendJob | Email service exists at `internal/service/email/service.go` |
| NotificationJob | `internal/infra/jobs/notification_job.go` |

### 1.4 TODO_v0.3.0.md (Partial)

**File**: `docs/dev/design/planning/TODO_v0.3.0.md`

**Status says "Not Started" on line 32 -- should be "In Progress" or "Mostly Complete (Backend)".**

#### Movie Module Backend (Lines 50-125) -- DONE

| Item | Code Evidence |
|---|---|
| Database schema (movies, movie_files, movie_credits, movie_collections, movie_genres, movie_watched) | Migrations 000021-000026, 000031 |
| Entity types | `internal/content/movie/types.go` |
| Repository + PostgreSQL | `internal/content/movie/repository.go`, `repository_postgres.go` |
| Service (get, list, search, watch progress, metadata refresh) | `internal/content/movie/service.go`, `cached_service.go` |
| Library provider (scan, match, handle changes) | `internal/content/movie/library_service.go`, `library_scanner.go`, `library_matcher.go` |
| Handler (all endpoints) | `internal/content/movie/handler.go`, `internal/api/movie_handlers.go`, `movie_converters.go` |
| River jobs (MetadataRefresh, LibraryScan, FileMatch, SearchIndex) | `internal/content/movie/moviejobs/*.go` (4 workers) |
| fx module | `internal/content/movie/module.go` |
| Tests | `service_test.go`, `handler_test.go`, `library_scanner_test.go`, `library_matcher_test.go`, `types_test.go` |
| Mediainfo extraction | `internal/content/movie/mediainfo.go` |

#### Metadata Service / TMDb (Lines 143-173) -- DONE

| Item | Code Evidence |
|---|---|
| TMDb client (rate limiting, retry, response caching) | `internal/service/metadata/providers/tmdb/client.go` |
| TMDb provider (search, details, credits, images, similar, recommendations) | `internal/service/metadata/providers/tmdb/provider.go` |
| Metadata service (aggregation, provider chain) | `internal/service/metadata/service.go` |
| Image service | `internal/infra/image/service.go` |
| Handler | `internal/api/handler_metadata.go` |

#### Search Service / Typesense (Lines 174-218) -- DONE

| Item | Code Evidence |
|---|---|
| Typesense setup, client config | `internal/infra/search/module.go` |
| Movie search schema | `internal/service/search/movie_schema.go` |
| Movie search service | `internal/service/search/movie_service.go` |
| Cached search service | `internal/service/search/cached_service.go` |
| Handler | `internal/api/handler_search.go` |
| Tests | `movie_service_test.go` |

#### Radarr Integration (Lines 219-257) -- DONE

| Item | Code Evidence |
|---|---|
| Radarr client (API v3, auth, error handling) | `internal/integration/radarr/client.go` |
| Radarr service (get movies, sync, quality profiles, root folders) | `internal/integration/radarr/service.go` |
| Mapper (Radarr types to Revenge types) | `internal/integration/radarr/mapper.go` |
| Webhook handler | `internal/integration/radarr/webhook_handler.go` |
| Handler (admin status, sync, quality profiles) | `internal/api/handler_radarr.go` |
| River jobs (RadarrSyncWorker, RadarrWebhookWorker) | `internal/integration/radarr/jobs.go` |
| fx module | `internal/integration/radarr/module.go` |
| Tests | `client_test.go`, `jobs_test.go`, `mapper_test.go` |

#### Frontend (Lines 258-337) -- NOT DONE

No frontend code exists. SvelteKit project has not been started.

#### Collections (Lines 127-142) -- PARTIALLY DONE

Migration 000024 creates collection tables. No dedicated collection service, but collection data flows through TMDb metadata.

### 1.5 TODO_v0.4.0.md (Partial)

**File**: `docs/dev/design/planning/TODO_v0.4.0.md`

**Status says "Not Started" on line 29 -- should be "In Progress" or "Mostly Complete (Backend)".**

#### TV Show Module Backend (Lines 46-143) -- DONE

| Item | Code Evidence |
|---|---|
| Database schema (full tvshow schema) | Migration 000032 (comprehensive: series, seasons, episodes, genres, credits, episode_files, watch_progress, networks) |
| Types | `internal/content/tvshow/types.go` |
| Repository + PostgreSQL | `internal/content/tvshow/repository.go`, `repository_postgres.go` |
| Service (get show, seasons, episodes, watch progress, continue watching) | `internal/content/tvshow/service.go` |
| Metadata provider interface | `internal/content/tvshow/metadata_provider.go` |
| Adapters (metadata, scanner) | `internal/content/tvshow/adapters/metadata_adapter.go`, `scanner_adapter.go` |
| Handler (all endpoints) | `internal/api/tvshow_handlers.go`, `tvshow_converters.go` |
| River jobs (LibraryScan, MetadataRefresh, FileMatch, SearchIndex, SeriesRefresh) | `internal/content/tvshow/jobs/jobs.go` (5 workers) |
| fx module | `internal/content/tvshow/module.go` |
| sqlc generated queries | `internal/content/tvshow/db/*.sql.go` (series, seasons, episodes, genres, credits, episode_files, watch_progress, networks) |
| Tests | `service_test.go`, `types_test.go`, `jobs_test.go`, `adapters/*_test.go` |

#### TheTVDB Integration (Lines 144-167) -- DONE

| Item | Code Evidence |
|---|---|
| TVDb client | `internal/service/metadata/providers/tvdb/client.go` |
| TVDb provider (search, details, seasons, episodes, artwork) | `internal/service/metadata/providers/tvdb/provider.go` |
| Type mapping | `internal/service/metadata/providers/tvdb/mapping.go` |

#### TMDb TV Support (Lines 169-177) -- DONE

TMDb provider implements `TVShowProvider` interface (search, details, seasons, episodes, credits, images).

#### Sonarr Integration (Lines 179-216) -- DONE

| Item | Code Evidence |
|---|---|
| Sonarr client (API v3, auth) | `internal/integration/sonarr/client.go` |
| Sonarr service (get series, sync, quality profiles, root folders) | `internal/integration/sonarr/service.go` |
| Mapper | `internal/integration/sonarr/mapper.go` |
| Webhook handler | `internal/integration/sonarr/webhook_handler.go` |
| Handler | `internal/api/handler_sonarr.go` |
| River jobs (SonarrSyncWorker, SonarrWebhookWorker) | `internal/integration/sonarr/jobs.go` |
| Tests | `client_test.go` |

#### Episode Watch Progress (Lines 218-237) -- DONE

Watch progress tracking implemented via `internal/content/tvshow/db/watch_progress.sql.go` and tvshow service.

#### Search Integration (Lines 240-280) -- PARTIALLY DONE

TV show search schema is not yet visible as a separate schema file like `movie_schema.go`, but the SearchIndexWorker exists in tvshow jobs suggesting it is at least partially wired.

### 1.6 TODO_v0.8.0.md (Partial)

**File**: `docs/dev/design/planning/TODO_v0.8.0.md`

#### Notification Service (Lines 138-189) -- DONE

| Item | Code Evidence |
|---|---|
| Notification dispatcher | `internal/service/notification/dispatcher.go` |
| Notification types | `internal/service/notification/notification.go` |
| Email agent | `internal/service/notification/agents/email.go` |
| Discord agent | `internal/service/notification/agents/discord.go` |
| Gotify agent | `internal/service/notification/agents/gotify.go` |
| Webhook agent | `internal/service/notification/agents/webhook.go` |
| River notification job | `internal/infra/jobs/notification_job.go` |
| Tests | `dispatcher_test.go`, `notification_test.go`, `agents_test.go` |

Note: The planning docs placed notification service in v0.8.0 but it was implemented much earlier. The design doc `docs/dev/design/services/NOTIFICATION.md` is active.

### 1.7 TODO_v0.9.0.md (Partial - QAR Schema)

**File**: `docs/dev/design/planning/TODO_v0.9.0.md`

The QAR module has database schema (migration 000001 creates `qar` schema) and sqlc-generated placeholder code:
- `internal/content/qar/db/db.go`, `models.go`, `placeholder.sql.go`, `querier.go`

No service, handler, or business logic yet -- schema only as expected.

---

## 2. Broken Design Doc References

The doc rewrite moved many docs to `planned/` and reorganized the directory structure. The planning docs still reference the OLD paths.

### 2.1 TODO_v0.2.0.md

| Line | Old Path Referenced | Current Path |
|---|---|---|
| 471 | `../services/AUTH.md` | EXISTS: `services/AUTH.md` |
| 472 | `../services/USER.md` | EXISTS: `services/USER.md` |
| 473 | `../services/SESSION.md` | EXISTS: `services/SESSION.md` |
| 474 | `../services/RBAC.md` | EXISTS: `services/RBAC.md` |
| 475 | `../services/APIKEYS.md` | EXISTS: `services/APIKEYS.md` |
| 476 | `../services/OIDC.md` | EXISTS: `services/OIDC.md` |
| 477 | `../services/SETTINGS.md` | EXISTS: `services/SETTINGS.md` |
| 478 | `../services/ACTIVITY.md` | EXISTS: `services/ACTIVITY.md` |
| 479 | `../services/LIBRARY.md` | EXISTS: `services/LIBRARY.md` |
| 482 | `../integrations/infrastructure/DRAGONFLY.md` | EXISTS: `integrations/infrastructure/DRAGONFLY.md` |
| 483 | `../integrations/infrastructure/RIVER.md` | EXISTS: `integrations/infrastructure/RIVER.md` |
| 484 | `../integrations/infrastructure/POSTGRESQL.md` | EXISTS: `integrations/infrastructure/POSTGRESQL.md` |
| 487 | `../integrations/auth/AUTHELIA.md` | **MOVED** to `planned/integrations/AUTHELIA.md` |
| 488 | `../integrations/auth/AUTHENTIK.md` | **MOVED** to `planned/integrations/AUTHENTIK.md` |
| 489 | `../integrations/auth/KEYCLOAK.md` | **MOVED** to `planned/integrations/KEYCLOAK.md` |
| 490 | `../integrations/auth/GENERIC_OIDC.md` | EXISTS: `integrations/auth/GENERIC_OIDC.md` |
| 493 | `../technical/EMAIL.md` | EXISTS: `technical/EMAIL.md` |
| 494 | `../features/shared/RBAC.md` | EXISTS: `features/shared/RBAC.md` |

### 2.2 TODO_v0.3.0.md

| Line | Old Path Referenced | Current Path |
|---|---|---|
| 429 | `../features/video/MOVIE_MODULE.md` | EXISTS: `features/video/MOVIE_MODULE.md` |
| 430 | `../features/shared/COLLECTIONS.md` | **MOVED** to `planned/features/COLLECTIONS.md` |
| 431 | `../features/shared/LIBRARIES.md` | EXISTS: `features/shared/LIBRARIES.md` |
| 434 | `../integrations/metadata/video/TMDB.md` | EXISTS: `integrations/metadata/video/TMDB.md` |
| 435 | `../integrations/servarr/RADARR.md` | EXISTS: `integrations/servarr/RADARR.md` |
| 436 | `../integrations/infrastructure/TYPESENSE.md` | EXISTS: `integrations/infrastructure/TYPESENSE.md` |
| 441 | `../services/METADATA.md` | EXISTS: `services/METADATA.md` |
| 442 | `../services/SEARCH.md` | EXISTS: `services/SEARCH.md` |
| 443 | `../services/LIBRARY.md` | EXISTS: `services/LIBRARY.md` |
| 447 | `../features/playback/WATCH_NEXT_CONTINUE_WATCHING.md` | **MOVED** to `planned/features/playback/WATCH_NEXT_CONTINUE_WATCHING.md` |
| 448 | `../features/playback/TRICKPLAY.md` | **MOVED** to `planned/features/playback/TRICKPLAY.md` |
| 449 | `../features/playback/SKIP_INTRO.md` | **MOVED** to `planned/features/playback/SKIP_INTRO.md` |
| 451 | `../technical/FRONTEND.md` | EXISTS: `technical/FRONTEND.md` |
| 452 | `../technical/API.md` | EXISTS: `technical/API.md` |
| 453 | `../patterns/HTTP_CLIENT.md` | EXISTS: `patterns/HTTP_CLIENT.md` |

### 2.3 TODO_v0.4.0.md

| Line | Old Path Referenced | Current Path |
|---|---|---|
| 353 | `../features/video/TVSHOW_MODULE.md` | EXISTS: `features/video/TVSHOW_MODULE.md` |
| 354 | `../features/playback/WATCH_NEXT_CONTINUE_WATCHING.md` | **MOVED** to `planned/features/playback/WATCH_NEXT_CONTINUE_WATCHING.md` |
| 357 | `../integrations/metadata/video/THETVDB.md` | EXISTS: `integrations/metadata/video/THETVDB.md` |
| 358 | `../integrations/metadata/video/TMDB.md` | EXISTS |
| 359 | `../integrations/servarr/SONARR.md` | EXISTS: `integrations/servarr/SONARR.md` |
| 362 | `../features/playback/TRICKPLAY.md` | **MOVED** to `planned/features/playback/TRICKPLAY.md` |
| 363 | `../features/playback/SKIP_INTRO.md` | **MOVED** to `planned/features/playback/SKIP_INTRO.md` |

### 2.4 TODO_v0.5.0.md

| Line | Old Path Referenced | Current Path |
|---|---|---|
| 403 | `../features/music/MUSIC_MODULE.md` | **MOVED** to `planned/features/music/MUSIC_MODULE.md` |
| 404 | `../technical/AUDIO_STREAMING.md` | **MOVED** to `planned/technical/AUDIO_STREAMING.md` |
| 407 | `../integrations/metadata/music/MUSICBRAINZ.md` | **MOVED** to `planned/integrations/music/MUSICBRAINZ.md` |
| 408 | `../integrations/metadata/music/LASTFM.md` | **MOVED** to `planned/integrations/music/LASTFM.md` |
| 409 | `../integrations/scrobbling/LASTFM_SCROBBLE.md` | **MOVED** to `planned/integrations/scrobbling/LASTFM_SCROBBLE.md` |
| 410 | `../integrations/scrobbling/LISTENBRAINZ.md` | **MOVED** to `planned/integrations/scrobbling/LISTENBRAINZ.md` |
| 411 | `../integrations/servarr/LIDARR.md` | **MOVED** to `planned/integrations/LIDARR.md` |
| 412 | `../integrations/metadata/music/SPOTIFY.md` | **MOVED** to `planned/integrations/music/SPOTIFY.md` |
| 413 | `../integrations/metadata/music/DISCOGS.md` | **MOVED** to `planned/integrations/music/DISCOGS.md` |
| 416 | `../features/shared/SCROBBLING.md` | **MOVED** to `planned/features/SCROBBLING.md` |

### 2.5 TODO_v0.6.0.md

| Line | Old Path Referenced | Current Path |
|---|---|---|
| 342 | `../features/playback/TRICKPLAY.md` | **MOVED** to `planned/features/playback/TRICKPLAY.md` |
| 343 | `../features/playback/SKIP_INTRO.md` | **MOVED** to `planned/features/playback/SKIP_INTRO.md` |
| 344 | `../features/playback/WATCH_NEXT_CONTINUE_WATCHING.md` | **MOVED** to `planned/features/playback/WATCH_NEXT_CONTINUE_WATCHING.md` |
| 345 | `../features/playback/SYNCPLAY.md` | **MOVED** to `planned/features/playback/SYNCPLAY.md` |
| 346 | `../integrations/casting/CHROMECAST.md` | **MOVED** to `planned/integrations/casting/CHROMECAST.md` |
| 347 | `../integrations/casting/DLNA.md` | **MOVED** to `planned/integrations/casting/DLNA.md` |

### 2.6 TODO_v0.7.0.md

| Line | Old Path Referenced | Current Path |
|---|---|---|
| 370 | `../features/audiobook/AUDIOBOOK_MODULE.md` | **MOVED** to `planned/features/audiobook/AUDIOBOOK_MODULE.md` |
| 371 | `../features/book/BOOK_MODULE.md` | **MOVED** to `planned/features/book/BOOK_MODULE.md` |
| 372 | `../features/podcasts/PODCASTS.md` | **MOVED** to `planned/features/podcasts/PODCASTS.md` |
| 373 | `../integrations/metadata/books/AUDIBLE.md` | **MOVED** to `planned/integrations/books/AUDIBLE.md` |
| 374 | `../integrations/metadata/books/OPENLIBRARY.md` | **MOVED** to `planned/integrations/books/OPENLIBRARY.md` |

### 2.7 TODO_v0.8.0.md

| Line | Old Path Referenced | Current Path |
|---|---|---|
| 330 | `../features/shared/SCROBBLING.md` | **MOVED** to `planned/features/SCROBBLING.md` |
| 331 | `../features/shared/ANALYTICS_SERVICE.md` | **MOVED** to `planned/features/ANALYTICS_SERVICE.md` |
| 332 | `../services/NOTIFICATION.md` | EXISTS: `services/NOTIFICATION.md` |
| 333 | `../features/shared/REQUEST_SYSTEM.md` | **MOVED** to `planned/features/REQUEST_SYSTEM.md` |
| 334 | `../services/GRANTS.md` | **MOVED** to `planned/services/GRANTS.md` |
| 335 | `../features/shared/I18N.md` | **MOVED** to `planned/features/I18N.md` |
| 336 | `../integrations/scrobbling/TRAKT.md` | **MOVED** to `planned/integrations/scrobbling/TRAKT.md` |
| 337 | `../integrations/scrobbling/LASTFM_SCROBBLE.md` | **MOVED** to `planned/integrations/scrobbling/LASTFM_SCROBBLE.md` |
| 338 | `../integrations/scrobbling/LISTENBRAINZ.md` | **MOVED** to `planned/integrations/scrobbling/LISTENBRAINZ.md` |

### 2.8 TODO_v0.9.0.md

| Line | Old Path Referenced | Current Path |
|---|---|---|
| 384 | `../features/adult/ADULT_CONTENT_SYSTEM.md` | EXISTS: `features/adult/ADULT_CONTENT_SYSTEM.md` |
| 385 | `../integrations/metadata/adult/STASHDB.md` | **MOVED** to `planned/integrations/adult/STASHDB.md` |
| 386 | `../integrations/servarr/WHISPARR.md` | **MOVED** to `planned/integrations/WHISPARR.md` |
| 387 | `../features/livetv/LIVE_TV_DVR.md` | **MOVED** to `planned/features/livetv/LIVE_TV_DVR.md` |
| 388 | `../features/photos/PHOTOS_LIBRARY.md` | **MOVED** to `planned/features/photos/PHOTOS_LIBRARY.md` |
| 389 | `../features/comics/COMICS_MODULE.md` | **MOVED** to `planned/features/comics/COMICS_MODULE.md` |

### 2.9 TODO_v1.0.0.md

| Line | Old Path Referenced | Current Path |
|---|---|---|
| 328 | `../technical/TECH_STACK.md` | EXISTS: `technical/TECH_STACK.md` |
| 329 | `../architecture/ARCHITECTURE.md` | EXISTS: `architecture/ARCHITECTURE.md` |
| 330 | `../operations/VERSIONING.md` | EXISTS: `operations/VERSIONING.md` |

### 2.10 ROADMAP.md and INDEX.md

| File | Line | Old Path | Status |
|---|---|---|---|
| ROADMAP.md | 528 | `../00_SOURCE_OF_TRUTH.md` | **DELETED** - no longer exists |
| ROADMAP.md | 531 | `../../../.workingdir/PLANNING_ANALYSIS.md` | May be stale reference |
| TODO_v0.0.0.md | 204 | `../00_SOURCE_OF_TRUTH.md` | **DELETED** |
| TODO_v0.1.0.md | 338 | `../00_SOURCE_OF_TRUTH.md` | **DELETED** |
| TODO_v0.2.0.md | 501 | `../00_SOURCE_OF_TRUTH.md` | **DELETED** |
| TODO_v0.3.0.md | 460 | `../00_SOURCE_OF_TRUTH.md` | **DELETED** |
| TODO_v0.4.0.md | 370 | `../00_SOURCE_OF_TRUTH.md` | **DELETED** |
| TODO_v0.5.0.md | 423 | `../00_SOURCE_OF_TRUTH.md` | **DELETED** |
| TODO_v0.6.0.md | 341 | `../00_SOURCE_OF_TRUTH.md` | **DELETED** |
| TODO_v0.7.0.md | 369 | `../00_SOURCE_OF_TRUTH.md` | **DELETED** |
| TODO_v0.8.0.md | 329 | `../00_SOURCE_OF_TRUTH.md` | **DELETED** |
| TODO_v0.9.0.md | 383 | `../00_SOURCE_OF_TRUTH.md` | **DELETED** |
| TODO_v1.0.0.md | 327 | `../00_SOURCE_OF_TRUTH.md` | **DELETED** |
| INDEX.md | 28 | `../../sources/SOURCES.md` | Likely stale |
| TODO_v0.9.0.md | 50 | `../00_SOURCE_OF_TRUTH.md#qar-obfuscation-terminology` | **DELETED** |

---

## 3. Gap Analysis - Features Not Covered by Any Active Design Doc

Items in the planning docs that reference design docs which have been moved to `planned/` and have **no active replacement**:

| Planning Doc | Feature | Old Design Doc | Active Replacement? |
|---|---|---|---|
| TODO_v0.3.0 | Collections | `features/shared/COLLECTIONS.md` | **NO** - only in `planned/`. Collections exist in code (migration 000024) but no active design doc. |
| TODO_v0.3.0/v0.4.0/v0.6.0 | Watch Next / Continue Watching | `features/playback/WATCH_NEXT_CONTINUE_WATCHING.md` | **NO** - in `planned/`. Basic implementation exists in tvshow service. |
| TODO_v0.6.0 | Trickplay | `features/playback/TRICKPLAY.md` | **NO** - only in `planned/` |
| TODO_v0.6.0 | Skip Intro | `features/playback/SKIP_INTRO.md` | **NO** - only in `planned/` |
| TODO_v0.6.0 | SyncPlay | `features/playback/SYNCPLAY.md` | **NO** - only in `planned/` |
| TODO_v0.6.0 | Chromecast | `integrations/casting/CHROMECAST.md` | **NO** - only in `planned/` |
| TODO_v0.6.0 | DLNA | `integrations/casting/DLNA.md` | **NO** - only in `planned/` |
| TODO_v0.5.0 | Music Module | `features/music/MUSIC_MODULE.md` | **NO** - only in `planned/` |
| TODO_v0.5.0 | Audio Streaming | `technical/AUDIO_STREAMING.md` | **NO** - only in `planned/` |
| TODO_v0.5.0 | All music integrations | Various `integrations/metadata/music/*` | **NO** - all in `planned/` |
| TODO_v0.7.0 | Audiobook Module | `features/audiobook/AUDIOBOOK_MODULE.md` | **NO** - only in `planned/` |
| TODO_v0.7.0 | Book Module | `features/book/BOOK_MODULE.md` | **NO** - only in `planned/` |
| TODO_v0.7.0 | Podcast Module | `features/podcasts/PODCASTS.md` | **NO** - only in `planned/` |
| TODO_v0.8.0 | Scrobbling | `features/shared/SCROBBLING.md` | **NO** - only in `planned/` |
| TODO_v0.8.0 | Analytics | `features/shared/ANALYTICS_SERVICE.md` | **NO** - only in `planned/` |
| TODO_v0.8.0 | Request System | `features/shared/REQUEST_SYSTEM.md` | **NO** - only in `planned/` |
| TODO_v0.8.0 | Grants | `services/GRANTS.md` | **NO** - only in `planned/` |
| TODO_v0.8.0 | i18n | `features/shared/I18N.md` | **NO** - only in `planned/` |
| TODO_v0.8.0 | All scrobbling integrations | `integrations/scrobbling/*` | **NO** - all in `planned/` |
| TODO_v0.9.0 | StashDB, Whisparr | Various integration docs | **NO** - in `planned/` |
| TODO_v0.9.0 | Live TV, Photos, Comics | Various feature docs | **NO** - all in `planned/` |

**Note**: This is expected and correct -- the doc rewrite properly moved unimplemented features to `planned/`. The issue is that the planning docs still link to the OLD non-planned paths, so the links are broken.

---

## 4. Stale Status Markers

### 4.1 ROADMAP.md

| Line | Current Marker | Should Be |
|---|---|---|
| 69 | `Current Version: v0.1.3 (Skeleton Complete + CI Fixes)` | Should reflect actual progress (at least "v0.2.0+ in progress" or similar) |
| 91 | `Phase 2: Implementation Phase IN PROGRESS` | Still correct but "Current Focus" text is wrong |
| 93 | `Current Focus: Building the foundation and MVP (v0.0.0 to v0.3.0)` | Should note v0.2.0 is essentially done, v0.3.0/v0.4.0 backends in progress |
| 108 | v0.2.0: `Designed` | Should be `Complete` or `Nearly Complete` |
| 109 | v0.3.0: `Designed` | Should be `In Progress (Backend ~70%)` |
| 110 | v0.4.0: `Designed` | Should be `In Progress (Backend ~60%)` |
| 230 | v0.2.0 status: `Not Started` | **WRONG** - Should be `Complete` or `Nearly Complete` |
| 262 | v0.3.0 status: `Not Started` | **WRONG** - Should be `In Progress` |
| 300 | v0.4.0 status: `Not Started` | **WRONG** - Should be `In Progress` |
| All v0.2.0 deliverables (lines 235-248) | All `[ ]` unchecked | Most should be `[x]` |
| All v0.3.0 backend deliverables (lines 269-274) | All `[ ]` unchecked | Backend items should be `[x]` |
| All v0.4.0 deliverables (lines 306-309) | All `[ ]` unchecked | Backend items should be `[x]` |

### 4.2 INDEX.md

| Line | Current Status | Should Be |
|---|---|---|
| 13 | ROADMAP.md: `Planned` | Should be `Partial` (some milestones complete, some in progress) |
| 14 | TODO_v0.0.0: `Planned` | Should be `Complete` |
| 15 | TODO_v0.1.0: `Planned` | Should be `Complete` |
| 16 | TODO_v0.2.0: `Planned` | Should be `Complete` or `Nearly Complete` |
| 17 | TODO_v0.3.0: `Planned` | Should be `Partial` (backend mostly done) |
| 18 | TODO_v0.4.0: `Planned` | Should be `Partial` (backend mostly done) |

### 4.3 TODO_v0.2.0.md

| Line | Current | Should Be |
|---|---|---|
| 34 | `Status: Not Started` | `Status: Complete` or `Nearly Complete` |
| Every `[ ]` in Auth Service (lines 53-89) | Unchecked | Should all be `[x]` |
| Every `[ ]` in User Service (lines 93-134) | Unchecked | Should all be `[x]` |
| Every `[ ]` in Session Service (lines 138-166) | Unchecked | Should all be `[x]` |
| Every `[ ]` in RBAC Service (lines 169-213) | Unchecked | Should all be `[x]` |
| Every `[ ]` in API Keys Service (lines 215-246) | Unchecked | Should all be `[x]` |
| Every `[ ]` in OIDC Service (lines 248-281) | Unchecked | Should all be `[x]` |
| Every `[ ]` in Settings Service (lines 283-308) | Unchecked | Should all be `[x]` |
| Every `[ ]` in Activity Service (lines 310-337) | Unchecked | Should all be `[x]` |
| Every `[ ]` in Library Service (lines 339-372) | Unchecked | Should all be `[x]` |
| Every `[ ]` in Health Enhancement (lines 374-386) | Unchecked | Most should be `[x]` |
| Every `[ ]` in PostgreSQL Integration (lines 388-401) | Unchecked | Should all be `[x]` |
| Every `[ ]` in Dragonfly Integration (lines 403-419) | Unchecked | Should all be `[x]` |
| Every `[ ]` in River Setup (lines 421-436) | Unchecked | Should all be `[x]` |

### 4.4 TODO_v0.3.0.md

| Line | Current | Should Be |
|---|---|---|
| 32 | `Status: Not Started` | `Status: In Progress` |
| Movie Module Backend (lines 50-125) | All `[ ]` | All should be `[x]` |
| Collection Support (lines 127-142) | All `[ ]` | Schema done `[x]`, service `[ ]` |
| Metadata Service TMDb (lines 143-173) | All `[ ]` | All should be `[x]` |
| Search Service (lines 174-218) | All `[ ]` | All should be `[x]` |
| Radarr Integration (lines 219-257) | All `[ ]` | All should be `[x]` |
| Frontend (lines 258-337) | All `[ ]` | Remain `[ ]` (not started) |

### 4.5 TODO_v0.4.0.md

| Line | Current | Should Be |
|---|---|---|
| 29 | `Status: Not Started` | `Status: In Progress` |
| TV Show Module Backend (lines 46-143) | All `[ ]` | All should be `[x]` |
| TheTVDB Integration (lines 144-167) | All `[ ]` | All should be `[x]` |
| TMDb TV Support (lines 169-177) | All `[ ]` | All should be `[x]` |
| Sonarr Integration (lines 179-216) | All `[ ]` | All should be `[x]` |
| Episode Watch Progress (lines 218-237) | All `[ ]` | Schema and basic tracking `[x]` |
| Search Integration (lines 240-280) | All `[ ]` | Partially done |
| Frontend Updates (lines 282-323) | All `[ ]` | Remain `[ ]` (not started) |

---

## 5. Missing From Planning - Implemented But Unplanned

The following exist in code but are NOT mentioned in any planning doc:

| Feature | Code Location | Notes |
|---|---|---|
| **MFA Service** (TOTP, WebAuthn, backup codes) | `internal/service/mfa/` + migrations 000016-000020, 000028 | Full MFA with TOTP, WebAuthn, backup codes. 5 dedicated migrations. Not in any TODO. Active design doc: `services/MFA.md` |
| **MFA API Handler** | `internal/api/handler_mfa.go` | Full MFA endpoints in OpenAPI spec |
| **Email Service** | `internal/service/email/service.go`, `module.go` | Listed as "EmailSendJob" in v0.2.0 River section but no dedicated section. Active design doc: `services/EMAIL.md` |
| **Storage Service** | `internal/service/storage/storage.go`, `s3.go`, `mock_storage.go` | S3-compatible storage service with local + S3 backends. Active design doc: `services/STORAGE.md`. Not mentioned in any planning doc. |
| **Image Service** | `internal/infra/image/service.go`, `module.go` | Image proxy and caching. Active design doc: `infrastructure/IMAGE.md`. Not in planning. |
| **Observability Infrastructure** | `internal/infra/observability/metrics.go`, `middleware.go`, `pprof.go`, `server.go` | Full observability with Prometheus metrics, pprof, dedicated metrics server. Active design doc: `infrastructure/OBSERVABILITY.md`. Only vaguely referenced as "Prometheus Metrics" in v0.2.0 Health Enhancement. |
| **Raft / Leader Election** | `internal/infra/raft/election.go`, `module.go` | Distributed leader election. Not in any planning doc. |
| **Notification Service (fully)** | `internal/service/notification/` with 4 agent types | Planned for v0.8.0 but already implemented. Active design doc: `services/NOTIFICATION.md` |
| **Rate Limiting Middleware** | `internal/api/middleware/ratelimit.go`, `ratelimit_redis.go` | Both in-memory and Redis-backed rate limiting. Not specifically planned. |
| **Request ID Middleware** | `internal/api/middleware/request_id.go` | Not specifically planned. |
| **Request Metadata Middleware** | `internal/api/middleware/request_metadata.go` | Not specifically planned. |
| **Localization in API** | `internal/api/localization.go` | API localization support. Not specifically planned. |
| **Crypto Package** | `internal/crypto/encryption.go`, `password.go` | Encryption and password hashing. Not specifically planned as separate package. |
| **Validate Package** | `internal/validate/convert.go` | Type conversion with validation. Not specifically planned. |
| **Content Shared Infrastructure** | `internal/content/shared/` (jobs, library, matcher, metadata, scanner) | Shared content infrastructure for all content modules. Not explicitly planned as reusable layer. |
| **Moderator Role** | Migration 000027 | Additional RBAC role not in original planning (which listed admin, user, guest, legacy:read). |
| **Fine-Grained Permissions** | Migration 000029 | Granular permission system beyond original Casbin plan. |
| **Failed Login Attempts** | Migration 000030 | Security hardening not in planning. |
| **Movie Multilanguage** | Migration 000031 | i18n support for movie metadata, much earlier than the v0.8.0 i18n plan. |
| **9th GitHub Actions Workflow** | `.github/workflows/wiki-sync.yml` | ROADMAP mentions 7-8 workflows. There are now 9. |

---

## 6. Summary of Required Changes Per File

### ROADMAP.md

1. **Update "Current Version"** (line 69): Change from "v0.1.3" to reflect actual state
2. **Update "Current Phase" description** (lines 91-95): Note v0.2.0 complete, v0.3.0/v0.4.0 in progress
3. **Update Milestone Overview table** (lines 101-116): Change v0.2.0 from "Designed" to "Complete", v0.3.0/v0.4.0 to "In Progress"
4. **Check v0.2.0 deliverables** (lines 235-248): Check all items as complete
5. **Check v0.3.0 backend deliverables** (lines 269-274): Check backend items, leave frontend unchecked
6. **Check v0.4.0 backend deliverables** (lines 306-309): Check backend items, leave frontend unchecked
7. **Remove reference to `00_SOURCE_OF_TRUTH.md`** (line 528): File was deleted
8. **Remove/update `.workingdir/PLANNING_ANALYSIS.md` reference** (line 531)
9. **Add MFA service** to v0.2.0 deliverables (missing entirely)
10. **Add email service** to v0.2.0 deliverables (only mentioned as River job)
11. **Add storage service** somewhere in planning
12. **Add notification service** to v0.2.0 or note it was pulled forward from v0.8.0
13. **Update workflow count**: ROADMAP references Design Phase docs count; the CI/CD section should note 9 workflows

### INDEX.md

1. **Update all status indicators** (lines 13-24): v0.0.0 and v0.1.0 are "Complete", v0.2.0 is "Complete", v0.3.0 and v0.4.0 are "Partial"
2. **Remove Sources reference** (line 28): `../../sources/SOURCES.md` likely stale

### TODO_v0.0.0.md

1. **Remove reference to `00_SOURCE_OF_TRUTH.md`** (line 204)
2. No other major changes needed (correctly marked Complete)

### TODO_v0.1.0.md

1. **Check items that were "deferred to v0.2.0" and are now done** (lines 101-103, 216-218, 233-235): Mark as complete
2. **Remove reference to `00_SOURCE_OF_TRUTH.md`** (line 338)
3. **Note**: Version table on line 304 has old dependency versions (these are historical and may stay)

### TODO_v0.2.0.md

1. **Change status from "Not Started" to "Complete"** (line 34)
2. **Check ALL deliverable items as done** (lines 50-436): Nearly all checkboxes should be `[x]`
3. **Add MFA Service section** (missing entirely from v0.2.0 -- 5 migrations, full service, handler)
4. **Add Email Service section** (only mentioned as a River job)
5. **Add Storage Service section** (exists in code, no planning)
6. **Fix broken doc references**: `AUTHELIA.md`, `AUTHENTIK.md`, `KEYCLOAK.md` moved to `planned/`
7. **Remove reference to `00_SOURCE_OF_TRUTH.md`** (line 501)
8. **Add verification checklist completions** (lines 440-448)
9. **Note the notification service was pulled forward** from v0.8.0

### TODO_v0.3.0.md

1. **Change status from "Not Started" to "In Progress"** (line 32)
2. **Check all backend movie deliverables as done** (lines 50-125)
3. **Check Metadata/TMDb deliverables as done** (lines 143-173)
4. **Check Search/Typesense deliverables as done** (lines 174-218)
5. **Check Radarr deliverables as done** (lines 219-257)
6. **Leave Frontend items unchecked** (lines 258-337)
7. **Fix broken doc references**: `COLLECTIONS.md`, `WATCH_NEXT*.md`, `TRICKPLAY.md`, `SKIP_INTRO.md` moved to `planned/`
8. **Remove reference to `00_SOURCE_OF_TRUTH.md`** (line 460)
9. **Update collection status** (migration done, no dedicated service)

### TODO_v0.4.0.md

1. **Change status from "Not Started" to "In Progress"** (line 29)
2. **Check all TV Show Backend deliverables as done** (lines 46-143)
3. **Check TheTVDB Integration as done** (lines 144-167)
4. **Check TMDb TV Support as done** (lines 169-177)
5. **Check Sonarr Integration as done** (lines 179-216)
6. **Check Episode Watch Progress as done** (lines 218-237)
7. **Leave Frontend items unchecked** (lines 282-323)
8. **Fix broken doc references**: `WATCH_NEXT*.md`, `TRICKPLAY.md`, `SKIP_INTRO.md` moved to `planned/`
9. **Remove reference to `00_SOURCE_OF_TRUTH.md`** (line 370)

### TODO_v0.5.0.md

1. **Fix ALL design doc references** (lines 403-416): Every single one moved to `planned/`
2. **Remove reference to `00_SOURCE_OF_TRUTH.md`** (line 423)

### TODO_v0.6.0.md

1. **Fix ALL design doc references** (lines 342-347): All moved to `planned/`
2. **Remove reference to `00_SOURCE_OF_TRUTH.md`** (line 341)

### TODO_v0.7.0.md

1. **Fix ALL design doc references** (lines 370-374): All moved to `planned/`
2. **Remove reference to `00_SOURCE_OF_TRUTH.md`** (line 369)

### TODO_v0.8.0.md

1. **Note notification service is already implemented** (lines 138-189)
2. **Fix design doc references** (lines 330-338): Most moved to `planned/`
3. **Remove reference to `00_SOURCE_OF_TRUTH.md`** (line 329)
4. **Consider removing or noting i18n partial work**: API localization already exists in code

### TODO_v0.9.0.md

1. **Note QAR schema already exists** (migration 000001 creates `qar` schema, sqlc placeholder code exists)
2. **Fix design doc references** (lines 385-389): STASHDB, WHISPARR, LIVE_TV_DVR, PHOTOS, COMICS all moved to `planned/`
3. **Remove reference to `00_SOURCE_OF_TRUTH.md`** (line 383)

### TODO_v1.0.0.md

1. **Remove reference to `00_SOURCE_OF_TRUTH.md`** (line 327)
2. References to ARCHITECTURE.md, TECH_STACK.md, VERSIONING.md are still valid

---

## Appendix: Quick Reference - Actual Codebase State

### Services Implemented (15)

| Service | Path | Migrations |
|---|---|---|
| auth | `internal/service/auth/` | 000008, 000009, 000010, 000030 |
| user | `internal/service/user/` | 000002, 000006, 000007 |
| session | `internal/service/session/` | 000003, 000020 |
| rbac | `internal/service/rbac/` | 000011, 000027, 000029 |
| apikeys | `internal/service/apikeys/` | 000012 |
| oidc | `internal/service/oidc/` | 000013 |
| mfa | `internal/service/mfa/` | 000016, 000017, 000018, 000019, 000028 |
| settings | `internal/service/settings/` | 000004, 000005 |
| activity | `internal/service/activity/` | 000014 |
| library | `internal/service/library/` | 000015 |
| metadata | `internal/service/metadata/` | (uses content tables) |
| search | `internal/service/search/` | (uses Typesense) |
| email | `internal/service/email/` | (no migration) |
| notification | `internal/service/notification/` | (no migration) |
| storage | `internal/service/storage/` | (no migration) |

### Content Modules (3)

| Module | Path | Migrations | Workers |
|---|---|---|---|
| movie | `internal/content/movie/` | 000021-000026, 000031 | 4 (MetadataRefresh, LibraryScan, FileMatch, SearchIndex) |
| tvshow | `internal/content/tvshow/` | 000032 | 5 (LibraryScan, MetadataRefresh, FileMatch, SearchIndex, SeriesRefresh) |
| qar | `internal/content/qar/` | 000001 (schema) | 0 (schema only) |

### Infrastructure (10 modules)

| Module | Path |
|---|---|
| database | `internal/infra/database/` |
| cache | `internal/infra/cache/` (otter L1 + rueidis L2) |
| jobs | `internal/infra/jobs/` (River, 5 queues, 17 total workers) |
| health | `internal/infra/health/` |
| image | `internal/infra/image/` |
| logging | `internal/infra/logging/` |
| observability | `internal/infra/observability/` |
| search | `internal/infra/search/` (Typesense) |
| raft | `internal/infra/raft/` |

### Integrations (4)

| Integration | Path |
|---|---|
| Radarr | `internal/integration/radarr/` |
| Sonarr | `internal/integration/sonarr/` |
| TMDb | `internal/service/metadata/providers/tmdb/` |
| TVDb | `internal/service/metadata/providers/tvdb/` |

### River Workers (17 total)

| Worker | Location | Queue |
|---|---|---|
| CleanupWorker | `infra/jobs/cleanup_job.go` | low |
| NotificationWorker | `infra/jobs/notification_job.go` | high |
| ActivityCleanupWorker | `service/activity/cleanup.go` | low |
| LibraryScanCleanupWorker | `service/library/cleanup.go` | low |
| MovieMetadataRefreshWorker | `content/movie/moviejobs/metadata_refresh.go` | default |
| MovieLibraryScanWorker | `content/movie/moviejobs/library_scan.go` | bulk |
| MovieFileMatchWorker | `content/movie/moviejobs/file_match.go` | default |
| MovieSearchIndexWorker | `content/movie/moviejobs/search_index.go` | default |
| RadarrSyncWorker | `integration/radarr/jobs.go` | default |
| RadarrWebhookWorker | `integration/radarr/jobs.go` | high |
| TVShowLibraryScanWorker | `content/tvshow/jobs/jobs.go` | bulk |
| TVShowMetadataRefreshWorker | `content/tvshow/jobs/jobs.go` | default |
| TVShowFileMatchWorker | `content/tvshow/jobs/jobs.go` | default |
| TVShowSearchIndexWorker | `content/tvshow/jobs/jobs.go` | default |
| TVShowSeriesRefreshWorker | `content/tvshow/jobs/jobs.go` | default |
| SonarrSyncWorker | `integration/sonarr/jobs.go` | default |
| SonarrWebhookWorker | `integration/sonarr/jobs.go` | high |

### Migrations (32 pairs = 64 files)

000001 through 000032 in `internal/infra/database/migrations/shared/`.

### OpenAPI Endpoints

153 handler methods in the ogen-generated server interface.

### GitHub Actions Workflows (9)

ci.yml, coverage.yml, develop.yml, labels.yml, pr-checks.yml, release-please.yml, security.yml, stale.yml, wiki-sync.yml.
