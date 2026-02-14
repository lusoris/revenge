# Pre-Frontend Readiness Analysis

## TL;DR: Nothing is blocking. Start building.

All **161 OpenAPI-defined routes** have real handler implementations. Zero operations fall through to `UnimplementedHandler`. The backend is production-ready for frontend integration.

---

## What's Complete

### Auth & Sessions
- JWT login/register/logout/refresh
- Cookie-based auth + CSRF double-submit pattern
- API key auth (create, list, get, revoke)
- OIDC SSO (full flow: authorize → callback, admin provider CRUD, user linking)
- MFA: TOTP setup/verify/disable, WebAuthn register/login, backup codes
- Password reset, email verification, change password
- Session management: list, revoke specific, revoke all, refresh

### User Management
- Profile CRUD, preferences, avatar upload
- Admin: list users with search/filters, soft-delete

### Content Browsing — Movies (21 handlers)
- List, get, search, recently added, top rated
- Continue watching, watch history, stats
- Files, cast, crew, genres, collections, similar movies
- Watch progress CRUD, mark as watched, metadata refresh

### Content Browsing — TV Shows (26 handlers)
- List, get, search, recently added, continue watching, stats
- Seasons, episodes, cast, crew, genres, networks
- Episode progress CRUD, mark watched (single + bulk)
- Next episode, recent/upcoming episodes, metadata refresh

### Search (Typesense) — 10 handlers
- Full-text search with faceted filtering
- Autocomplete (movies + TV shows)
- Multi-collection unified search
- Search facets for filter dropdowns
- Admin reindex endpoints

### Metadata (TMDb proxy) — ~30 handlers
- Movie/TV/season/episode/person metadata, credits, images
- Image proxy (`/images/{type}/{size}/{path}`) — serves posters/backdrops without exposing API keys
- Similar/recommendations, external IDs, content ratings

### Library Management — 11 handlers
- Full CRUD, scan trigger, scan history
- Per-user permissions (grant/revoke/list)
- Cross-content genre listing

### Admin / RBAC — 10 handlers
- Full Casbin RBAC: policies, roles, permissions CRUD
- User role assignment/removal
- Activity logs with search, user/resource filters, stats

### Integrations — 10 handlers
- Radarr: status, sync trigger, quality profiles, root folders, webhook
- Sonarr: same

### Playback — 4 handlers + HLS stream
- Start/stop/get session, heartbeat
- HLS stream handler at `/api/v1/playback/stream/`
- astiav in-process transcoding

### Settings
- Server settings: list, get, update (admin)
- User settings: list, get, update, delete

### Real-Time Events (SSE)
- `GET /api/v1/events` — per-client category filtering
- Auth via Bearer or `?token=` query param
- Bridges `notification.Agent` → SSE fanout

### Security
- CORS fully configured (default `["*"]`, configurable)
- CSP + security headers middleware
- API docs at `GET /api/docs` (Scalar), spec at `GET /api/openapi.yaml`

---

## What's Missing (Not Blocking)

### Frontend Static File Serving
- No embedded SPA serving — the server only serves API routes
- **Decision needed**: Embed in Go binary vs. deploy frontend separately (nginx/CDN/SvelteKit standalone)
- **Not blocking**: Frontend dev can run locally with a dev server proxying to the API

### Nice-to-Have Features (from TODO.md)

| Feature | Status | Priority |
|---|---|---|
| Collections & playlists | Not started | Medium — TMDb collections work, custom playlists missing |
| Skip Intro / Credits detection | Not started | Low — player enhancement |
| SyncPlay (watch together) | Not started | Low — social feature |
| Trickplay (timeline thumbnails) | Not started | Low — player enhancement |
| Release calendar | Not started | Low — discovery feature |
| Content request system | Not started | Low — user → admin flow |
| Music/Audiobook/Book modules | Not started | Low — new content types |
| Circuit breaker for external APIs | Missing | Medium — resilience |
| Cache warming on startup | Missing | Low — performance |

---

## Architecture Decision: SPA Serving

Options:
1. **Embedded in Go binary** (`embed.FS` + `http.FileServer` with SPA fallback) — single binary deployment
2. **Separate nginx/caddy** — traditional reverse proxy setup, already used in docker-compose
3. **SvelteKit standalone** — SSR-capable, runs its own Node server
4. **CDN** — static hosting (Cloudflare Pages, Vercel, etc.)

Recommendation: Start with option 2 (nginx) for dev, decide on production architecture later. The API is CORS-ready for any approach.
