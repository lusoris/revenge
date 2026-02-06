# API Reference

<!-- DESIGN: technical -->

**Spec**: `api/openapi/openapi.yaml` (8,889 lines, OpenAPI 3.1)
**Generator**: ogen (generates `internal/api/ogen/`)
**Base Path**: `/api/v1`
**Auth**: Bearer JWT (`/api/v1/auth/login`)

> REST API generated from OpenAPI spec via ogen, with hand-written handlers

---

## Endpoint Groups

### Auth (9 endpoints)

| Method | Path | Purpose |
|--------|------|---------|
| POST | /auth/register | Register new user |
| POST | /auth/login | Login, get JWT |
| POST | /auth/logout | Logout |
| POST | /auth/refresh | Refresh JWT token |
| POST | /auth/verify-email | Verify email address |
| POST | /auth/resend-verification | Resend verification email |
| POST | /auth/forgot-password | Request password reset |
| POST | /auth/reset-password | Reset password with token |
| POST | /auth/change-password | Change password (authenticated) |

### MFA (7 endpoints)

| Method | Path | Purpose |
|--------|------|---------|
| GET | /mfa/status | Get MFA status |
| POST | /mfa/totp/setup | Begin TOTP setup |
| POST | /mfa/totp/verify | Verify TOTP code |
| DELETE | /mfa/totp | Remove TOTP |
| POST | /mfa/backup-codes/generate | Generate backup codes |
| POST | /mfa/backup-codes/regenerate | Regenerate backup codes |
| POST | /mfa/enable | Enable MFA |

### Movies (14 endpoints)

`/movies`, `/movies/{id}`, `/movies/{id}/files`, `/movies/{id}/cast`, `/movies/{id}/crew`, `/movies/{id}/genres`, `/movies/{id}/collection`, `/movies/{id}/similar`, `/movies/{id}/progress`, `/movies/{id}/watched`, `/movies/{id}/refresh`, `/movies/search`, `/movies/recently-added`, `/movies/top-rated`, `/movies/continue-watching`, `/movies/watch-history`, `/movies/stats`

Plus `/collections/{id}` and `/collections/{id}/movies`.

### TV Shows (20 endpoints)

`/tvshows`, `/tvshows/{id}`, `/tvshows/{id}/seasons`, `/tvshows/{id}/episodes`, `/tvshows/{id}/cast`, `/tvshows/{id}/crew`, `/tvshows/{id}/genres`, `/tvshows/{id}/networks`, `/tvshows/{id}/watch-stats`, `/tvshows/{id}/next-episode`, `/tvshows/{id}/refresh`, `/tvshows/search`, `/tvshows/recently-added`, `/tvshows/continue-watching`, `/tvshows/stats`, `/tvshows/episodes/recent`, `/tvshows/episodes/upcoming`, `/tvshows/seasons/{id}`, `/tvshows/seasons/{id}/episodes`, `/tvshows/episodes/{id}`, `/tvshows/episodes/{id}/files`, `/tvshows/episodes/{id}/progress`, `/tvshows/episodes/{id}/watched`

### Metadata (6 endpoints)

`/metadata/search/movie`, `/metadata/search/tv`, `/metadata/movie/{tmdbId}`, `/metadata/tv/{tmdbId}`, `/metadata/tv/{tmdbId}/season/{seasonNumber}`, `/metadata/tv/{tmdbId}/season/{seasonNumber}/episode/{episodeNumber}`

### Search (4 endpoints)

`/search/movies`, `/search/movies/autocomplete`, `/search/movies/facets`, `/search/reindex`

### Users (4 endpoints)

`/users/me`, `/users/me/preferences`, `/users/me/avatar`, `/users/{userId}`

### Sessions (4 endpoints)

`/sessions`, `/sessions/current`, `/sessions/refresh`, `/sessions/{sessionId}`

### RBAC (6 endpoints)

`/rbac/policies`, `/rbac/users/{userId}/roles`, `/rbac/users/{userId}/roles/{role}`, `/rbac/roles`, `/rbac/roles/{roleName}`, `/rbac/roles/{roleName}/permissions`, `/rbac/permissions`

### API Keys (2 endpoints)

`/apikeys`, `/apikeys/{keyId}`

### OIDC (9 endpoints)

`/oidc/providers`, `/oidc/auth/{provider}`, `/oidc/callback/{provider}`, `/users/me/oidc`, `/users/me/oidc/{provider}/link`, `/users/me/oidc/{provider}`, `/admin/oidc/providers`, `/admin/oidc/providers/{providerId}`, `/admin/oidc/providers/{providerId}/enable`, `/admin/oidc/providers/{providerId}/disable`, `/admin/oidc/providers/{providerId}/default`

### Activity (4 endpoints)

`/admin/activity`, `/admin/activity/users/{userId}`, `/admin/activity/resources/{resourceType}/{resourceId}`, `/admin/activity/stats`, `/admin/activity/actions`

### Libraries (6 endpoints)

`/libraries`, `/libraries/{libraryId}`, `/libraries/{libraryId}/scan`, `/libraries/{libraryId}/scans`, `/libraries/{libraryId}/permissions`, `/libraries/{libraryId}/permissions/{userId}`

### Settings (4 endpoints)

`/settings/server`, `/settings/server/{key}`, `/settings/user`, `/settings/user/{key}`

### Integrations (10 endpoints)

Radarr: `/admin/integrations/radarr/status`, `../sync`, `../quality-profiles`, `../root-folders`, `/webhooks/radarr`

Sonarr: `/admin/integrations/sonarr/status`, `../sync`, `../quality-profiles`, `../root-folders`, `/webhooks/sonarr`

### Images (1 endpoint)

`/images/{type}/{size}/{path}` — Proxied image serving.

## Handler Architecture

Hand-written handlers in `internal/api/`:

| File | Handles |
|------|---------|
| handler.go | Main handler struct, implements ogen server interface |
| handler_activity.go | Activity log endpoints |
| handler_apikeys.go | API key CRUD |
| handler_library.go | Library management |
| handler_metadata.go | Metadata search/lookup |
| handler_mfa.go | MFA setup/management |
| handler_oidc.go | OIDC provider management |
| handler_radarr.go | Radarr integration + webhooks |
| handler_rbac.go | RBAC policies/roles |
| handler_search.go | Typesense search |
| handler_session.go | Session management |
| handler_sonarr.go | Sonarr integration + webhooks |
| context.go | Request context helpers |
| errors.go | Error response formatting |

## Security

- **Authentication**: Bearer JWT (from `/api/v1/auth/login`)
- **Security scheme**: `bearerAuth` (HTTP bearer, JWT format)
- **Global security**: None (endpoints opt-in via security requirement)

## Code Generation

ogen generates from `api/openapi/openapi.yaml` into `internal/api/ogen/`:
- Server interface
- Request/response types
- Validation
- Router

Regenerate with `make ogen`.

## Related Documentation

- [CONFIGURATION.md](CONFIGURATION.md) - Server and auth config
- [../architecture/ARCHITECTURE.md](../architecture/ARCHITECTURE.md) - System architecture
