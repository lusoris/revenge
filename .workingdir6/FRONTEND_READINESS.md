# Frontend Readiness Audit & Fix Plan

## Audit Summary (2026-02-07)

### What's Done (Good News)
- **166/166 API endpoints implemented** — zero stubs, zero unimplemented methods
- **14/14 services fully implemented** — auth, session, user, library, settings, rbac, apikeys, activity, oidc, mfa, search (movie+tv), playback, notification
- **All fx modules wired** — build clean, vet clean
- **HLS streaming works** — separate handler with CORS headers for HLS.js
- **Image proxy works** — CORS headers set for cross-origin image loading
- **Error responses are JSON** — consistent `{error, message, code}` format
- **Rate limiting works** — dual-tier (auth: 5/s, global: 50/s), memory or Redis backend
- **Request tracing** — `X-Request-ID` in all responses

### What's Broken / Missing

---

## CRITICAL — Frontend cannot make API calls

### C1: No CORS Middleware on API Endpoints

**Problem**: The ogen API server does NOT set `Access-Control-Allow-Origin` on responses. The generated OPTIONS handler in `oas_cfg_gen.go` sets `Allow-Methods` and `Allow-Headers: Content-Type` but never sets `Allow-Origin`. Any browser request from a different origin (e.g. SvelteKit dev server on :5173) will be blocked.

**Scope**: ALL API endpoints except HLS streams and image proxy (those have `*` CORS already).

**Fix**:
1. Add CORS config to `config.go`:
   ```go
   type CORSConfig struct {
       AllowedOrigins   []string      // e.g. ["http://localhost:5173", "https://app.example.com"]
       AllowCredentials bool          // default: true
       MaxAge           time.Duration // default: 12h
   }
   ```
2. Create `internal/api/middleware/cors.go` — proper CORS middleware that:
   - Sets `Access-Control-Allow-Origin` (from config, or reflect request origin if in allowed list)
   - Sets `Access-Control-Allow-Methods: GET, POST, PUT, PATCH, DELETE, OPTIONS`
   - Sets `Access-Control-Allow-Headers: Authorization, Content-Type, X-API-Key, X-Request-ID`
   - Sets `Access-Control-Expose-Headers: X-Request-ID, X-RateLimit-Remaining, Retry-After`
   - Sets `Access-Control-Allow-Credentials: true`
   - Sets `Access-Control-Max-Age` from config
   - Handles OPTIONS preflight with 204
3. Wire into middleware chain in `server.go` (first in chain, before rate limiter)

**Files**: `config.go`, `middleware/cors.go` (new), `server.go`

---

## HIGH — Frontend works but UX is broken

### H1: List Endpoints Missing Pagination

**Problem**: Most list endpoints return bare arrays with no pagination metadata. Frontend cannot implement pagination UI, infinite scroll, or know total counts.

**Affected endpoints** (no pagination params, no total in response):
- `ListSessions` → `{sessions: [...]}`
- `ListAPIKeys` → `{keys: [...]}`
- `ListPolicies` → `{policies: [...]}`
- `ListRoles` → `{roles: [...]}`
- `ListOIDCProviders` → `{providers: [...]}`
- `AdminListOIDCProviders` → `{providers: [...]}`
- `ListUserOIDCLinks` → `{links: [...]}`
- `ListLibraryPermissions` → `{permissions: [...]}`
- `ListMovies` → bare `[...]` array
- `ListTVShows` → bare `[...]` array

**Exception** (already correct): `ListLibraries` → `{libraries: [...], total: N}` ✓

**Fix**: For admin/small-collection endpoints (sessions, API keys, policies, roles, OIDC providers, permissions) — these are bounded lists (max ~100 items). Add `total` field to response schemas but don't add server-side pagination (overkill).

For movies/tvshows — these can be thousands. Add proper `limit`/`offset` params and wrap response in `{items: [...], total: N, limit: N, offset: N}`.

**Files**: `openapi.yaml` (schema changes) → regenerate ogen → update handlers

### H2: Mixed Pagination Scheme (limit/offset vs page/per_page)

**Problem**: Search endpoints use `page`/`per_page` while list endpoints use `limit`/`offset`. Frontend must implement two pagination patterns.

**Fix**: Standardize on `limit`/`offset` everywhere (more flexible, REST-standard). Update search endpoints in OpenAPI spec to use limit/offset. Map to page/per_page internally for Typesense.

**Files**: `openapi.yaml`, `handler_search.go`, search service layer

### H3: Some 201 Responses Have No Body Schema

**Problem**: Several creation endpoints return 201 with no response schema defined:
- `AddPolicy` — 201 but no response body
- `UpdateUserSetting` — 201 for create, no body

**Fix**: Add response schemas returning the created object. Frontend needs the ID/data of what was just created.

**Files**: `openapi.yaml` → regenerate ogen → update handlers

---

## MEDIUM — Polish for good frontend DX

### M1: No Real-Time Notifications (WebSocket/SSE)

**Problem**: Notification service exists (email, discord, gotify, webhook agents) but there's no frontend push channel. Frontend must poll for updates.

**Fix**: Add SSE endpoint `GET /api/v1/events/stream` that pushes events to connected browsers:
- Library scan progress
- New content added
- Job completion
- Admin notifications

**Approach**: SSE is simpler than WebSocket for one-way server→client events. Use `text/event-stream` with proper auth token in query param (SSE doesn't support custom headers).

**Files**: `handler_events.go` (new), `openapi.yaml`, notification service integration

### M2: No OpenAPI Spec Serving at Runtime

**Problem**: The `openapi.yaml` file is not served by the running server. Frontend devs need it for code generation and documentation.

**Fix**: Serve `GET /api/openapi.yaml` from the embedded spec. Optionally add Scalar or Redoc UI at `/api/docs`.

**Files**: `server.go` (add route), embed directive

### M3: Webhook Endpoints Have No Security Definition

**Problem**: Radarr/Sonarr webhook endpoints have no `security` in the OpenAPI spec. They're actually protected by webhook signing but this isn't documented.

**Fix**: Add webhook secret validation to spec (custom header `X-Webhook-Secret` or signature verification). Document in OpenAPI that these endpoints use webhook signing, not bearer auth.

**Files**: `openapi.yaml`

---

## LOW — Future improvements

### L1: Cookie-Based Auth Option

**Current**: Tokens returned in JSON body only. Frontend stores in memory/localStorage (XSS-vulnerable).

**Better**: Offer optional `Set-Cookie` with HttpOnly, Secure, SameSite=Strict for access/refresh tokens. SvelteKit's server-side rendering benefits from cookie auth.

**Defer**: Works fine without cookies for MVP. Add later when SvelteKit SSR is implemented.

### L2: CSRF Protection

**When**: Only needed if cookie auth (L1) is implemented. Token-in-header auth is inherently CSRF-safe.

---

## Execution Plan

### Phase 1: CORS + Pagination Fix (blocks frontend dev)

| Task | Priority | Est. Effort | Agent |
|------|----------|-------------|-------|
| C1: CORS middleware | CRITICAL | 1h | 1 |
| H1: Add `total` to list response schemas | HIGH | 1h | 2 |
| H2: Standardize pagination to limit/offset | HIGH | 1h | 2 |
| H3: Add 201 response schemas | HIGH | 30m | 2 |

### Phase 2: Polish (improves frontend DX)

| Task | Priority | Est. Effort | Agent |
|------|----------|-------------|-------|
| M2: Serve OpenAPI spec at runtime | MEDIUM | 15m | 1 |
| M3: Document webhook security | MEDIUM | 15m | 1 |
| M1: SSE events endpoint | MEDIUM | 2h | 1 |

### Phase 3: Security hardening (before production)

| Task | Priority | Est. Effort | Agent |
|------|----------|-------------|-------|
| L1: Cookie-based auth | LOW | 2h | 1 |
| L2: CSRF protection | LOW | 1h | 1 |

---

## Verification

After Phase 1:
1. `curl -H "Origin: http://localhost:5173" http://localhost:8096/api/v1/health/live -v` → sees `Access-Control-Allow-Origin` header
2. All list endpoints return `total` field
3. `make test` passes
4. `make test-live` passes
5. Frontend dev can start building against the API
