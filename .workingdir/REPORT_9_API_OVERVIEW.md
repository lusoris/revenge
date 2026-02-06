# Report 9: Full API Overview

**Generated**: 2026-02-06
**Scope**: Complete API surface of the Revenge Media Server

---

## Table of Contents

1. [Architecture Overview](#1-architecture-overview)
2. [Framework & Server Setup](#2-framework--server-setup)
3. [Middleware Chain](#3-middleware-chain)
4. [Authentication & Security](#4-authentication--security)
5. [Endpoint Reference](#5-endpoint-reference)
6. [Request/Response Types](#6-requestresponse-types)
7. [Error Handling](#7-error-handling)
8. [Dependency Injection Graph](#8-dependency-injection-graph)
9. [Observability](#9-observability)
10. [Configuration](#10-configuration)

---

## 1. Architecture Overview

```
HTTP Client
    |
[http.Server (configurable host:port)]
    |
RequestID HTTP Wrapper
    |
ogen Server (type-safe, code-generated from OpenAPI 3.1.0)
    |
Middleware Chain:
    1. RequestID (extract/generate X-Request-ID)
    2. RequestMetadata (IP, User-Agent, Accept-Language)
    3. HTTPMetrics (Prometheus counters/histograms)
    4. RateLimit (per-IP, memory or Redis backend)
    5. ErrorHandler (custom error formatting)
    |
Handler (implements ogen interfaces)
    |- JWT Bearer Auth validation
    |- Context injection (userID, username, sessionID)
    |- Delegates to service layer
    |
Service Layer (business logic)
    |- Repository pattern with interfaces
    |- Transaction support
    |
Repository Layer (sqlc-generated)
    |- PostgreSQL (pgx)
    |- Rueidis (Redis/Dragonfly cache)
```

**Key Files**:
- OpenAPI Spec: `api/openapi/openapi.yaml` (v0.1.0)
- Server: `internal/api/server.go`
- Main Handler: `internal/api/handler.go`
- Generated Code: `internal/api/ogen/`
- Middleware: `internal/api/middleware/`
- Entry Point: `cmd/revenge/main.go`

---

## 2. Framework & Server Setup

| Property | Value |
|----------|-------|
| Framework | **ogen** (OpenAPI code generation for Go) |
| OpenAPI Version | 3.1.0 |
| API Version | 0.1.0 |
| DI Framework | uber/fx |
| Database | PostgreSQL 18+ (pgx + sqlc) |
| Cache | Rueidis (Redis/Dragonfly) |
| Search | Typesense |
| Job Queue | River |

**Server Lifecycle**: fx hooks for OnStart (background listener) and OnStop (graceful shutdown with configurable timeout).

---

## 3. Middleware Chain

Executed in order per request:

| Order | Middleware | File | Purpose |
|-------|-----------|------|---------|
| 1 | RequestIDHTTPWrapper | `middleware/request_id.go` | Adds X-Request-ID to response headers |
| 2 | RequestIDMiddleware | `middleware/request_id.go` | Extracts/generates request ID, stores in context |
| 3 | RequestMetadataMiddleware | `middleware/request_metadata.go` | Extracts IP (X-Forwarded-For, X-Real-IP), User-Agent, Accept-Language |
| 4 | HTTPMetricsMiddleware | `internal/infra/observability/middleware.go` | Prometheus: requests_total, requests_duration, requests_in_flight |
| 5 | RateLimitMiddleware (auth) | `middleware/ratelimit.go` | Auth endpoints: 1 req/sec, 5 burst |
| 6 | RateLimitMiddleware (global) | `middleware/ratelimit.go` | All endpoints: 10 req/sec, 20 burst |
| 7 | ErrorHandler | `middleware/errors.go` | Custom error responses, 429 Retry-After header |

**Rate Limiting Backends**:
- **In-Memory** (`ratelimit.go`): sync.Map per-IP, cleanup every 5 min (TTL 10 min)
- **Redis** (`ratelimit_redis.go`): Lua sliding window, health monitoring, auto-fallback to in-memory

---

## 4. Authentication & Security

**Method**: HTTP Bearer (JWT)

**Implementation**: `HandleBearerAuth()` in `internal/api/handler.go`
- Validates JWT access token via `tokenManager.ValidateAccessToken()`
- Injects `userID`, `username` into request context
- Invalid/expired tokens logged at WARN level

**Context Keys**:
```go
GetUserID(ctx)    -> uuid.UUID
GetUsername(ctx)   -> string
GetSessionID(ctx) -> uuid.UUID
```

**Authorization**:
- RBAC via Casbin (role-based access control)
- Admin checks on admin-prefixed endpoints
- Library-level permissions (view, download, manage)

**Multi-Factor Authentication**:
- TOTP (Time-based One-Time Password)
- Backup codes
- WebAuthn support (planned)

**External Auth**:
- OIDC federation (generic, Authentik, Keycloak providers)

---

## 5. Endpoint Reference

### 5.1 Health Checks (No Auth)

| Method | Path | Handler | Purpose |
|--------|------|---------|---------|
| GET | `/healthz` | `GetLiveness` | Liveness probe (always 200) |
| GET | `/readyz` | `GetReadiness` | Readiness probe (checks dependencies) |
| GET | `/startupz` | `GetStartup` | Startup probe |

### 5.2 Authentication (`/api/v1/auth/`)

| Method | Path | Handler | Auth |
|--------|------|---------|------|
| POST | `/api/v1/auth/register` | `Register` | No |
| POST | `/api/v1/auth/login` | `Login` | No |
| POST | `/api/v1/auth/logout` | `Logout` | Yes |
| POST | `/api/v1/auth/refresh` | `RefreshToken` | No |
| POST | `/api/v1/auth/verify-email` | `VerifyEmail` | No |
| POST | `/api/v1/auth/resend-verification` | `ResendVerification` | Yes |
| POST | `/api/v1/auth/forgot-password` | `ForgotPassword` | No |
| POST | `/api/v1/auth/reset-password` | `ResetPassword` | No |
| POST | `/api/v1/auth/change-password` | `ChangePassword` | Yes |

### 5.3 OIDC Authentication (`/api/v1/auth/oidc/`)

| Method | Path | Handler | Auth |
|--------|------|---------|------|
| GET | `/api/v1/auth/oidc/providers` | `ListOIDCProviders` | No |
| GET | `/api/v1/auth/oidc/authorize` | `OidcAuthorize` | No |
| GET | `/api/v1/auth/oidc/callback` | `OidcCallback` | No |
| GET | `/api/v1/auth/oidc/links` | `ListUserOIDCLinks` | Yes |
| POST | `/api/v1/auth/oidc/links/{provider}` | `InitOIDCLink` | Yes |
| DELETE | `/api/v1/auth/oidc/links/{provider}` | `UnlinkOIDCProvider` | Yes |

### 5.4 Multi-Factor Authentication (`/api/v1/mfa/`)

File: `internal/api/handler_mfa.go`

| Method | Path | Handler | Purpose |
|--------|------|---------|---------|
| GET | `/api/v1/mfa/status` | `GetMFAStatus` | Get MFA configuration |
| POST | `/api/v1/mfa/totp/setup` | `SetupTOTP` | Generate TOTP secret + QR code |
| POST | `/api/v1/mfa/totp/verify` | `VerifyTOTP` | Verify & enable TOTP |
| DELETE | `/api/v1/mfa/totp` | `DisableTOTP` | Disable TOTP |
| POST | `/api/v1/mfa/backup-codes/generate` | `GenerateBackupCodes` | Generate backup codes |
| POST | `/api/v1/mfa/backup-codes/regenerate` | `RegenerateBackupCodes` | Replace all codes |
| POST | `/api/v1/mfa/enable` | `EnableMFA` | Require MFA for login |
| POST | `/api/v1/mfa/disable` | `DisableMFA` | Turn off MFA requirement |

### 5.5 Sessions (`/api/v1/sessions/`)

File: `internal/api/handler_session.go`

| Method | Path | Handler | Purpose |
|--------|------|---------|---------|
| GET | `/api/v1/sessions` | `ListSessions` | List active sessions |
| GET | `/api/v1/sessions/current` | `GetCurrentSession` | Current session info |
| POST | `/api/v1/sessions/current/logout` | `LogoutCurrent` | Logout current session |
| POST | `/api/v1/sessions/logout-all` | `LogoutAll` | Logout all sessions |
| POST | `/api/v1/sessions/refresh` | `RefreshSession` | Refresh access token |
| DELETE | `/api/v1/sessions/{sessionId}` | `RevokeSession` | Revoke specific session |

### 5.6 API Keys (`/api/v1/api-keys/`)

File: `internal/api/handler_apikeys.go`

| Method | Path | Handler | Purpose |
|--------|------|---------|---------|
| GET | `/api/v1/api-keys` | `ListAPIKeys` | List user API keys |
| POST | `/api/v1/api-keys` | `CreateAPIKey` | Create new key (scopes: read, write, admin) |
| GET | `/api/v1/api-keys/{keyId}` | `GetAPIKey` | Get key details |
| DELETE | `/api/v1/api-keys/{keyId}` | `RevokeAPIKey` | Revoke key |

### 5.7 Users (`/api/v1/users/`)

| Method | Path | Handler | Purpose |
|--------|------|---------|---------|
| GET | `/api/v1/users/me` | `GetCurrentUser` | Get authenticated user profile |
| PUT | `/api/v1/users/me` | `UpdateCurrentUser` | Update profile |
| GET | `/api/v1/users/{userId}` | `GetUserById` | Get public profile |
| GET | `/api/v1/users/me/preferences` | `GetUserPreferences` | Get preferences |
| PUT | `/api/v1/users/me/preferences` | `UpdateUserPreferences` | Update preferences |
| POST | `/api/v1/users/me/avatar` | `UploadAvatar` | Upload avatar image |

### 5.8 Settings (`/api/v1/settings/`)

| Method | Path | Handler | Purpose | Auth |
|--------|------|---------|---------|------|
| GET | `/api/v1/settings/server` | `ListServerSettings` | List all server settings | Admin |
| GET | `/api/v1/settings/server/{key}` | `GetServerSetting` | Get specific setting | Admin |
| PUT | `/api/v1/settings/server/{key}` | `UpdateServerSetting` | Update setting | Admin |
| GET | `/api/v1/settings/user` | `ListUserSettings` | Get user settings | Yes |
| GET | `/api/v1/settings/user/{key}` | `GetUserSetting` | Get specific user setting | Yes |
| PUT | `/api/v1/settings/user/{key}` | `UpdateUserSetting` | Update user setting | Yes |
| DELETE | `/api/v1/settings/user/{key}` | `DeleteUserSetting` | Delete user setting | Yes |

### 5.9 Movies (`/api/v1/movies/`)

File: `internal/api/movie_handlers.go`

| Method | Path | Handler | Parameters |
|--------|------|---------|------------|
| GET | `/api/v1/movies` | `ListMovies` | orderBy, limit, offset |
| GET | `/api/v1/movies/search` | `SearchMovies` | query, limit, offset |
| GET | `/api/v1/movies/recently-added` | `GetRecentlyAdded` | limit, offset |
| GET | `/api/v1/movies/top-rated` | `GetTopRated` | minVotes, limit, offset |
| GET | `/api/v1/movies/continue-watching` | `GetContinueWatching` | limit |
| GET | `/api/v1/movies/watch-history` | `GetWatchHistory` | limit, offset |
| GET | `/api/v1/movies/stats` | `GetUserMovieStats` | - |
| GET | `/api/v1/movies/{id}` | `GetMovie` | id (UUID) |
| GET | `/api/v1/movies/{id}/files` | `GetMovieFiles` | id |
| GET | `/api/v1/movies/{id}/cast` | `GetMovieCast` | id |
| GET | `/api/v1/movies/{id}/crew` | `GetMovieCrew` | id |
| GET | `/api/v1/movies/{id}/genres` | `GetMovieGenres` | id |
| GET | `/api/v1/movies/{id}/collection` | `GetMovieCollection` | id |
| GET | `/api/v1/movies/{id}/similar` | `GetSimilarMovies` | id |
| GET | `/api/v1/movies/{id}/progress` | `GetWatchProgress` | id |
| POST | `/api/v1/movies/{id}/progress` | `UpdateWatchProgress` | id, durationSeconds |
| DELETE | `/api/v1/movies/{id}/progress` | `DeleteWatchProgress` | id |
| POST | `/api/v1/movies/{id}/watched` | `MarkAsWatched` | id |
| POST | `/api/v1/movies/{id}/refresh` | `RefreshMovieMetadata` | id |

### 5.10 Collections (`/api/v1/collections/`)

| Method | Path | Handler | Purpose |
|--------|------|---------|---------|
| GET | `/api/v1/collections/{id}` | `GetCollection` | Get collection details |
| GET | `/api/v1/collections/{id}/movies` | `GetCollectionMovies` | Get collection movies |

### 5.11 TV Shows (`/api/v1/tvshows/`)

File: `internal/api/tvshow_handlers.go`

| Method | Path | Handler | Parameters |
|--------|------|---------|------------|
| GET | `/api/v1/tvshows` | `ListTVShows` | order_by, limit, offset |
| GET | `/api/v1/tvshows/search` | `SearchTVShows` | query, limit, offset |
| GET | `/api/v1/tvshows/recently-added` | `GetRecentlyAddedTVShows` | limit |
| GET | `/api/v1/tvshows/continue-watching` | `GetTVContinueWatching` | limit |
| GET | `/api/v1/tvshows/stats` | `GetUserTVStats` | - |
| GET | `/api/v1/tvshows/episodes/recent` | `GetRecentEpisodes` | limit, offset |
| GET | `/api/v1/tvshows/episodes/upcoming` | `GetUpcomingEpisodes` | limit, offset |
| GET | `/api/v1/tvshows/{id}` | `GetTVShow` | id |
| GET | `/api/v1/tvshows/{id}/seasons` | `GetTVShowSeasons` | id |
| GET | `/api/v1/tvshows/{id}/episodes` | `GetTVShowEpisodes` | id |
| GET | `/api/v1/tvshows/{id}/cast` | `GetTVShowCast` | id |
| GET | `/api/v1/tvshows/{id}/crew` | `GetTVShowCrew` | id |
| GET | `/api/v1/tvshows/{id}/genres` | `GetTVShowGenres` | id |
| GET | `/api/v1/tvshows/{id}/networks` | `GetTVShowNetworks` | id |
| GET | `/api/v1/tvshows/{id}/watch-stats` | `GetTVShowWatchStats` | id |
| GET | `/api/v1/tvshows/{id}/next-episode` | `GetTVShowNextEpisode` | id |
| POST | `/api/v1/tvshows/{id}/refresh` | `RefreshTVShowMetadata` | id |
| GET | `/api/v1/tvshows/seasons/{id}` | `GetTVSeason` | id |
| GET | `/api/v1/tvshows/seasons/{id}/episodes` | `GetTVSeasonEpisodes` | id |
| GET | `/api/v1/tvshows/episodes/{id}` | `GetTVEpisode` | id |
| GET | `/api/v1/tvshows/episodes/{id}/files` | `GetTVEpisodeFiles` | id |
| GET | `/api/v1/tvshows/episodes/{id}/progress` | `GetTVEpisodeProgress` | id |

### 5.12 Search & Discovery (`/api/v1/search/`)

File: `internal/api/handler_search.go`

| Method | Path | Handler | Purpose |
|--------|------|---------|---------|
| GET | `/api/v1/search/library/movies` | `SearchLibraryMovies` | Full-text Typesense search |
| GET | `/api/v1/search/autocomplete/movies` | `AutocompleteMovies` | Title autocomplete |
| GET | `/api/v1/search/facets` | `GetSearchFacets` | Available filter values |
| POST | `/api/v1/search/reindex` | `ReindexSearch` | Trigger search index rebuild (admin) |

### 5.13 Metadata (`/api/v1/metadata/`)

File: `internal/api/handler_metadata.go`

| Method | Path | Handler | Purpose |
|--------|------|---------|---------|
| GET | `/api/v1/metadata/movies/search` | `SearchMoviesMetadata` | Search TMDb for movies |
| GET | `/api/v1/metadata/movies/{tmdbId}` | `GetMovieMetadata` | Get TMDb movie details |
| GET | `/api/v1/metadata/images/{path}` | `GetProxiedImage` | Proxy TMDb images (poster/backdrop/profile/logo) |
| GET | `/api/v1/metadata/collections/{collectionId}` | `GetCollectionMetadata` | Get TMDb collection |
| GET | `/api/v1/metadata/tvshows/search` | `SearchTVShowsMetadata` | Search TMDb for TV shows |
| GET | `/api/v1/metadata/tvshows/{tmdbId}` | `GetTVShowMetadata` | Get TMDb TV show details |
| GET | `/api/v1/metadata/seasons/{seasonId}` | `GetSeasonMetadata` | Get season metadata |
| GET | `/api/v1/metadata/episodes/{episodeId}` | `GetEpisodeMetadata` | Get episode metadata |

### 5.14 Libraries (`/api/v1/libraries/`)

File: `internal/api/handler_library.go`

| Method | Path | Handler | Admin |
|--------|------|---------|-------|
| GET | `/api/v1/libraries` | `ListLibraries` | No |
| POST | `/api/v1/libraries` | `CreateLibrary` | Yes |
| GET | `/api/v1/libraries/{libraryId}` | `GetLibrary` | No |
| PUT | `/api/v1/libraries/{libraryId}` | `UpdateLibrary` | Yes |
| DELETE | `/api/v1/libraries/{libraryId}` | `DeleteLibrary` | Yes |
| POST | `/api/v1/libraries/{libraryId}/scan` | `TriggerLibraryScan` | Yes |
| GET | `/api/v1/libraries/{libraryId}/scans` | `ListLibraryScans` | Yes |
| GET | `/api/v1/libraries/{libraryId}/permissions` | `ListLibraryPermissions` | Yes |
| POST | `/api/v1/libraries/{libraryId}/permissions` | `GrantLibraryPermission` | Yes |
| DELETE | `/api/v1/libraries/{libraryId}/permissions/{permId}` | `RevokeLibraryPermission` | Yes |

**Library Types**: movie, tvshow, music, photo, book, audiobook, comic, podcast, adult

**Scan Types**: full, incremental

**Permission Levels**: view, download, manage

### 5.15 Admin: RBAC (`/api/v1/admin/rbac/`)

File: `internal/api/handler_rbac.go`

| Method | Path | Handler | Purpose |
|--------|------|---------|---------|
| GET | `/api/v1/admin/rbac/policies` | `ListPolicies` | List all policies |
| POST | `/api/v1/admin/rbac/policies` | `AddPolicy` | Add policy (subject, action, resource, effect) |
| DELETE | `/api/v1/admin/rbac/policies` | `RemovePolicy` | Remove policy |
| GET | `/api/v1/admin/rbac/roles` | `ListRoles` | List all roles |
| POST | `/api/v1/admin/rbac/roles` | `CreateRole` | Create role |
| GET | `/api/v1/admin/rbac/roles/{roleId}` | `GetRole` | Get role details |
| DELETE | `/api/v1/admin/rbac/roles/{roleId}` | `DeleteRole` | Delete role |
| POST | `/api/v1/admin/rbac/roles/{roleId}/permissions` | `UpdateRolePermissions` | Update permissions |
| GET | `/api/v1/admin/rbac/permissions` | `ListPermissions` | List all permissions |
| GET | `/api/v1/admin/rbac/users/{userId}/roles` | `GetUserRoles` | Get user roles |
| POST | `/api/v1/admin/rbac/users/{userId}/roles` | `AssignRole` | Assign role to user |
| DELETE | `/api/v1/admin/rbac/users/{userId}/roles/{roleId}` | `RemoveRole` | Remove role from user |

### 5.16 Admin: OIDC Providers (`/api/v1/admin/oidc/`)

File: `internal/api/handler_oidc.go`

| Method | Path | Handler | Purpose |
|--------|------|---------|---------|
| GET | `/api/v1/admin/oidc/providers` | `AdminListOIDCProviders` | List providers |
| POST | `/api/v1/admin/oidc/providers` | `AdminCreateOIDCProvider` | Create provider |
| GET | `/api/v1/admin/oidc/providers/{providerId}` | `AdminGetOIDCProvider` | Get provider |
| PUT | `/api/v1/admin/oidc/providers/{providerId}` | `AdminUpdateOIDCProvider` | Update provider |
| DELETE | `/api/v1/admin/oidc/providers/{providerId}` | `AdminDeleteOIDCProvider` | Delete provider |
| POST | `/api/v1/admin/oidc/providers/{providerId}/enable` | `AdminEnableOIDCProvider` | Enable provider |
| POST | `/api/v1/admin/oidc/providers/{providerId}/disable` | `AdminDisableOIDCProvider` | Disable provider |
| POST | `/api/v1/admin/oidc/default/{providerId}` | `AdminSetDefaultOIDCProvider` | Set default |

**Provider Types**: generic, authentik, keycloak

### 5.17 Admin: Activity Logs (`/api/v1/admin/activity/`)

File: `internal/api/handler_activity.go`

| Method | Path | Handler | Purpose |
|--------|------|---------|---------|
| GET | `/api/v1/admin/activity` | `SearchActivityLogs` | Search all activity logs |
| GET | `/api/v1/admin/activity/users/{userId}` | `GetUserActivityLogs` | Get user activity |
| GET | `/api/v1/admin/activity/resources/{resourceId}` | `GetResourceActivityLogs` | Get resource activity |
| GET | `/api/v1/admin/activity/stats` | `GetActivityStats` | Get activity statistics |
| GET | `/api/v1/admin/activity/recent` | `GetRecentActions` | Get recent actions |

### 5.18 Admin: Integrations (`/api/v1/admin/radarr/`, `/api/v1/admin/sonarr/`)

Files: `internal/api/handler_radarr.go`, `internal/api/handler_sonarr.go`

| Method | Path | Handler | Purpose |
|--------|------|---------|---------|
| GET | `/api/v1/admin/radarr/status` | `AdminGetRadarrStatus` | Radarr connection status |
| POST | `/api/v1/admin/radarr/sync` | `AdminTriggerRadarrSync` | Trigger Radarr sync |
| GET | `/api/v1/admin/radarr/quality-profiles` | `AdminGetRadarrQualityProfiles` | Get quality profiles |
| GET | `/api/v1/admin/radarr/root-folders` | `AdminGetRadarrRootFolders` | Get root folders |
| POST | `/api/v1/webhooks/radarr` | `HandleRadarrWebhook` | Radarr webhook receiver |
| GET | `/api/v1/admin/sonarr/status` | `AdminGetSonarrStatus` | Sonarr connection status |
| POST | `/api/v1/admin/sonarr/sync` | `AdminTriggerSonarrSync` | Trigger Sonarr sync |
| GET | `/api/v1/admin/sonarr/quality-profiles` | `AdminGetSonarrQualityProfiles` | Get quality profiles |
| GET | `/api/v1/admin/sonarr/root-folders` | `AdminGetSonarrRootFolders` | Get root folders |
| POST | `/api/v1/webhooks/sonarr` | `HandleSonarrWebhook` | Sonarr webhook receiver |

---

## 6. Request/Response Types

### 6.1 Core Entity Types

**Movie** (`ogen/oas_schemas_gen.go`):
```
id, tmdb_id, imdb_id, title, original_title, year, release_date, runtime,
overview, tagline, status, original_language, poster_path, backdrop_path,
trailer_url, vote_average (0-10), vote_count, popularity, budget, revenue,
library_added_at, metadata_updated_at, radarr_id, created_at, updated_at
```

**TVSeries**:
```
id, tmdb_id, tvdb_id, imdb_id, name, original_name, original_language,
overview, first_air_date, last_air_date, number_of_seasons, number_of_episodes,
status, poster_path, backdrop_path, vote_average, vote_count, popularity,
in_production, library_added_at, metadata_updated_at, sonarr_id, created_at, updated_at
```

**TVSeason**: `id, tv_show_id, season_number, name, overview, poster_path, air_date`

**TVEpisode**: `id, tv_show_id, season_number, episode_number, name, overview, air_date, runtime, vote_average, vote_count, still_path`

**User**: `id, email, username, display_name, avatar, is_admin, is_active, require_mfa, last_login, password_changed_at, created_at, updated_at`

**Library**: `id, name, description, type, path, scanner, is_active, created_at, updated_at`

### 6.2 Auth Request/Response Types

| Type | Fields |
|------|--------|
| `LoginRequest` | email, password |
| `LoginResponse` | user, accessToken, refreshToken, expiresIn |
| `RegisterRequest` | email, password, confirmPassword, displayName |
| `RefreshSessionRequest` | refreshToken |
| `RefreshSessionResponse` | accessToken, refreshToken, expiresIn |
| `ChangePasswordRequest` | old_password, new_password |
| `ForgotPasswordRequest` | email |
| `ResetPasswordRequest` | token, password |

### 6.3 API Key Types

| Type | Fields |
|------|--------|
| `APIKeyInfo` | id, user_id, name, description, key_prefix (`rv_xxxxx`), scopes, is_active, expires_at, last_used_at, created_at |
| `CreateAPIKeyRequest` | name, description, scopes (read/write/admin), expires_at |
| `APIKeyListResponse` | keys[] |

### 6.4 Search Types

| Type | Fields |
|------|--------|
| `SearchResults` | hits[], facets, total, limit, offset |
| `SearchHit` | id, type, title, overview, poster_path, score, highlights |
| `SearchFacets` | genres[], years[], languages[], ratings[] |
| `FacetValue` | value, count |
| `AutocompleteResults` | results[] |

### 6.5 Activity Log Types

| Type | Fields |
|------|--------|
| `ActivityLogEntry` | id, userId, username, action, resourceType, resourceId, changes, metadata, ipAddress, userAgent, success, errorMessage, createdAt |
| `ActivityStats` | totalCount, successCount, failedCount, oldestEntry, newestEntry |
| `ActionCount` | action, count |

### 6.6 Watch Progress Types

| Type | Fields |
|------|--------|
| `UpdateWatchProgressReq` | duration_seconds |
| `ContinueWatchingItem` | id, type, title, poster_path, last_watched_at, progress |
| `UserMovieStats` | total_movies, watched_movies, total_duration, average_rating |
| `UserTVStats` | total_shows, watched_shows, total_episodes, watched_episodes, total_duration |

### 6.7 Image Proxy Parameters

| Parameter | Values |
|-----------|--------|
| `type` | poster, backdrop, profile, logo |
| `size` | w185, w342, w500, w780, w300, w1280, w45, h632, original |
| `format` | optional format override |

### 6.8 Optional Type System

All API types use ogen's `Opt*` wrapper pattern instead of Go pointers:

```go
OptString, OptInt, OptInt64, OptFloat32, OptBool, OptUUID, OptDateTime, OptDate
OptNilString, OptNilInt, OptNilBool, OptNilDateTime  // Nullable variants
```

Methods: `IsSet()`, `Get()`, `Or(default)`, `SetTo(v)`, `Reset()`

---

## 7. Error Handling

### Error Response Structure

```json
{
  "code": 500,
  "message": "Error description",
  "details": {}
}
```

### Error Constructors (`internal/api/errors.go`)

| Function | HTTP Status |
|----------|-------------|
| `NotFoundError(msg)` | 404 |
| `UnauthorizedError(msg)` | 401 |
| `ForbiddenError(msg)` | 403 |
| `ConflictError(msg)` | 409 |
| `ValidationError(msg)` | 400 |
| `BadRequestError(msg)` | 400 |
| `InternalError(msg, err)` | 500 |
| `UnavailableError(msg)` | 503 |
| `TimeoutError(msg)` | 504 |

### Standard HTTP Status Codes Used

| Code | Meaning | Usage |
|------|---------|-------|
| 200 | OK | Successful GET/PUT |
| 201 | Created | Successful POST (new resource) |
| 202 | Accepted | Async operations (scan, reindex, refresh) |
| 204 | No Content | Successful DELETE, logout, etc. |
| 400 | Bad Request | Validation failures |
| 401 | Unauthorized | Missing/invalid token |
| 403 | Forbidden | Insufficient permissions |
| 404 | Not Found | Resource not found |
| 429 | Too Many Requests | Rate limited (includes Retry-After header) |
| 500 | Internal Server Error | Unexpected errors |
| 503 | Service Unavailable | Dependency not ready |

---

## 8. Dependency Injection Graph

**Framework**: `go.uber.org/fx`

```
app.Module
 |- config.Module          -> *config.Config
 |- logging.Module         -> *zap.Logger
 |- database.Module        -> *pgxpool.Pool, *db.Queries
 |- cache.Module           -> *cache.Client (rueidis)
 |- search.Module          -> search engine (Typesense)
 |- jobs.Module            -> *jobs.Client (River queue)
 |- raft.Module            -> leader election
 |- health.Module          -> *health.Service
 |- settings.Module        -> settings.Service
 |- user.Module            -> *user.Service
 |- auth.Module            -> *auth.Service, auth.TokenManager
 |- session.Module         -> *session.Service
 |- rbac.Module            -> *rbac.Service (Casbin)
 |- apikeys.Module         -> *apikeys.Service
 |- mfa.Module             -> *mfa.TOTPService, *mfa.BackupCodesService, *mfa.MFAManager
 |- oidc.Module            -> *oidc.Service
 |- activity.Module        -> *activity.Service
 |- library.Module         -> *library.Service
 |- movie.Module           -> *movie.Handler
 |- tvshow.Module          -> tvshow.Service
 |- moviejobs.Module       -> movie job workers
 |- tvshowjobs.Module      -> tvshow job workers
 |- radarr.Module          -> *radarr.SyncService
 |- sonarr.Module          -> *sonarr.SyncService
 |- metadatafx.Module      -> metadata.Service
 |- observability.Module   -> Prometheus metrics
 `- api.Module             -> *api.Server (HTTP server)
```

---

## 9. Observability

### Prometheus Metrics

| Metric | Type | Labels |
|--------|------|--------|
| `http_requests_total` | Counter | method, path, status |
| `http_request_duration_seconds` | Histogram | method, path |
| `http_requests_in_flight` | Gauge | - |

**Path Normalization**: UUIDs in paths are replaced with `{id}` to control cardinality.

### Request Tracing

- X-Request-ID header (generated if missing, propagated in response)
- Structured logging with slog

---

## 10. Configuration

### Server Configuration (`config.yaml`)

```yaml
server:
  host: "0.0.0.0"
  port: 8080
  read_timeout: 30s
  write_timeout: 30s
  idle_timeout: 120s
  shutdown_timeout: 30s
  rate_limit:
    enabled: true
    backend: "redis"     # or "memory"
    global:
      requests_per_second: 10
      burst: 20
    auth:
      requests_per_second: 1
      burst: 5
```

---

## Endpoint Count Summary

| Category | Endpoints |
|----------|-----------|
| Health | 3 |
| Auth | 9 |
| OIDC (user) | 6 |
| MFA | 8 |
| Sessions | 6 |
| API Keys | 4 |
| Users | 6 |
| Settings | 7 |
| Movies | 19 |
| Collections | 2 |
| TV Shows | 22 |
| Search | 4 |
| Metadata | 8 |
| Libraries | 10 |
| Admin: RBAC | 12 |
| Admin: OIDC | 8 |
| Admin: Activity | 5 |
| Admin: Radarr | 4 |
| Admin: Sonarr | 4 |
| Webhooks | 2 |
| **Total** | **~149** |

---

## File Organization

```
internal/api/
 |- server.go                    Server initialization & lifecycle
 |- handler.go                   Main handler (auth, users, settings, health)
 |- handler_activity.go          Activity logging
 |- handler_apikeys.go           API key management
 |- handler_mfa.go               Multi-factor auth
 |- handler_oidc.go              OpenID Connect
 |- handler_rbac.go              Role-based access control
 |- handler_session.go           Session management
 |- handler_radarr.go            Radarr integration
 |- handler_sonarr.go            Sonarr integration
 |- handler_library.go           Library management
 |- handler_metadata.go          Metadata lookup (TMDb)
 |- handler_search.go            Full-text search
 |- movie_handlers.go            Movie endpoints
 |- tvshow_handlers.go           TV show endpoints
 |- movie_converters.go          Movie DTO conversion
 |- tvshow_converters.go         TV show DTO conversion
 |- context.go                   Context helpers
 |- errors.go                    Error utilities
 |- localization.go              i18n support
 |- image_utils.go               Image processing
 |- module.go                    fx module definition
 |- middleware/
 |   |- errors.go                Error handling middleware
 |   |- ratelimit.go             In-memory rate limiting
 |   |- ratelimit_redis.go       Redis rate limiting
 |   |- request_id.go            Request ID middleware
 |   `- request_metadata.go     Request metadata extraction
 `- ogen/                        Generated code (~50 files)
     |- oas_schemas_gen.go       258 generated types
     |- oas_parameters_gen.go    Query/path parameter types
     |- oas_router_gen.go        Auto-generated routing
     `- ...
```
