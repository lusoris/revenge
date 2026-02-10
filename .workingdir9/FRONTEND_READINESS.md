# SvelteKit + Tailwind Frontend Readiness

**Datum:** 2026-02-10
**Geplanter Stack:** SvelteKit 2, Svelte 5, Tailwind CSS 4, shadcn-svelte

---

## Gesamtbewertung: 🟡 Fast bereit (95%)

Das Backend ist funktional komplett für ein MVP-Frontend. Es gibt keine fehlenden Endpoints — alle 188 Operationen sind implementiert. Die Blocker sind API-Konsistenz-Issues die das Frontend-Development erschweren würden.

---

## Bereit ✅

### Auth & Session Management
- [x] JWT + Refresh Token Flow (access + refresh)
- [x] Cookie-based Auth (HttpOnly access + refresh cookies)
- [x] CSRF Protection (double-submit cookie + `X-CSRF-Token` header)
- [x] MFA (TOTP + WebAuthn + Backup Codes)
- [x] OIDC/SSO (Generic, Authentik, Keycloak)
- [x] Session Management (list, revoke, current)
- [x] API Key Auth als Alternative
- [x] Account Registration + Email Verification
- [x] Password Reset Flow

**SvelteKit Integration:**
- SSR `load()` → Cookie wird automatisch mitgesendet
- Client-Side → CSRF Token aus JS-readable Cookie lesen
- OIDC → Server-Side Auth Hook ideal für `/oidc/callback`

### Content APIs
- [x] Movies: Full CRUD, Files, Cast/Crew, Genres, Collections, Similar, Watch Progress
- [x] TV Shows: Series/Seasons/Episodes, Files, Cast/Crew, Genres, Networks
- [x] Metadata Proxy (TMDb, TVDb + 10 weitere Provider)
- [x] Image Proxy (eliminiert CORS + Key Exposure)
- [x] Search: Full-Text + Autocomplete + Facets (Typesense)
- [x] Library Management (CRUD, Scan, Permissions)

### Streaming
- [x] HLS Playback Sessions (Create, Get, Delete)
- [x] HLS Stream Handler (`/api/v1/playback/stream/`)
- [x] Quality Profiles (multiple resolutions)
- [x] Audio Track Renditions
- [x] Subtitle Tracks (WebVTT)

### Admin
- [x] User Management (list, delete)
- [x] RBAC (roles, permissions, policies)
- [x] Activity Logging (audit trail)
- [x] Integration Management (Radarr, Sonarr)
- [x] Settings (server-wide + per-user)
- [x] Library Scans + Permissions

### Real-Time
- [x] SSE Events (`GET /api/v1/events`)
- [x] Category-based Filtering (`?categories=library,content,system`)
- [x] Auth via Bearer oder `?token=` Query Param

### Developer Experience
- [x] OpenAPI Spec (`GET /api/openapi.yaml`)
- [x] Scalar API Docs (`GET /api/docs`)
- [x] Consistent Error Schema (`{ code, message, details }`)
- [x] Health Probes for Dev Containers

---

## Blocker 🔴 (vor Frontend-Start fixen)

### 1. Property Naming: snake_case vs camelCase

**Problem:** Gemischte Konventionen in der API. TypeScript-Types können nicht einheitlich generiert werden.

**Betroffene Fläche:**
```
snake_case: auth, mfa, movies, tvshows, users, settings, sessions, search, apikeys, metadata
camelCase:  oidc, activity, libraries, radarr, sonarr, playback
```

**Fix:** Eine Konvention wählen, OpenAPI Spec anpassen, ogen regenerieren.

**Impact auf Frontend:** Jeder API-Call, jedes TypeScript-Interface, jede Komponente.

### 2. Fehlende `total` in List Responses

**Problem:** 17+ Endpoints geben nackte Arrays zurück. Kein `total` Count = kein Pagination-UI.

**Betroffene Routen (SvelteKit):**
```
/movies (search, continue-watching, watch-history)
/movies/[id] (files, genres)
/tvshows (search, continue-watching)
/tvshows/episodes (recent, upcoming)
/tvshows/[id] (seasons, episodes, genres, networks)
/genres
/collections/[id]/movies
```

**Fix:** Alle List-Endpoints in `{ items, total, limit, offset }` Envelope wrappen.

---

## Warnings 🟡 (sollte gefixt werden)

### 3. Pagination-Inkonsistenz

SvelteKit bräuchte 3 verschiedene Pagination-Utilities:
- `limit`/`offset` für Content-APIs
- `page`/`per_page` für Search-APIs
- Keine Pagination für bare-array Endpoints

**Empfehlung:** Alles auf `limit`/`offset` standardisieren.

### 4. HTTP Method Inkonsistenz (Progress)

```
POST /movies/{id}/progress        ← SvelteKit form action?
PUT  /tvshows/episodes/{id}/progress  ← fetch() PUT?
```

Unterschiedliche Verben für semantisch identische Operationen.

### 5. Sort-Parameter Chaos

```
Movies:  order_by=title|year|added|rating
TV:      order_by=created_at|title|first_air_date|vote_average|popularity
Search:  sort_by=...
```

Frontend braucht separate Sort-Logik pro Ressource.

### 6. Typesense DSL Exposure

`filter_by` erwartet Typesense-Syntax: `genres:=Action && year:>=2020`. Das Frontend müsste Typesense-Queries bauen — leakt Backend-Implementierung.

---

## SvelteKit Route-Mapping (Vorschlag)

```
src/routes/
├── +layout.svelte              ← Auth check, theme, nav
├── +page.svelte                ← Dashboard / Home
├── login/+page.svelte          ← /api/v1/auth/login
├── register/+page.svelte       ← /api/v1/auth/register
├── movies/
│   ├── +page.svelte            ← GET /api/v1/movies
│   ├── +page.server.ts         ← SSR load() mit Cookie-Auth
│   └── [id]/
│       ├── +page.svelte        ← GET /api/v1/movies/{id}
│       ├── +page.server.ts
│       └── play/+page.svelte   ← POST /api/v1/playback/sessions
├── tv/
│   ├── +page.svelte            ← GET /api/v1/tvshows
│   └── [id]/
│       ├── +page.svelte        ← GET /api/v1/tvshows/{id}
│       └── season/[sn]/
│           └── episode/[en]/
│               └── +page.svelte
├── search/+page.svelte         ← GET /api/v1/search/multi
├── libraries/
│   ├── +page.svelte            ← GET /api/v1/libraries
│   └── [id]/+page.svelte       ← GET /api/v1/libraries/{id}
├── settings/
│   ├── +page.svelte            ← GET /api/v1/settings/user
│   ├── security/+page.svelte   ← MFA, Sessions, API Keys
│   └── oidc/+page.svelte       ← OIDC Link/Unlink
├── admin/
│   ├── users/+page.svelte      ← GET /api/v1/admin/users
│   ├── activity/+page.svelte   ← GET /api/v1/admin/activity
│   ├── oidc/+page.svelte       ← Admin OIDC providers
│   ├── integrations/+page.svelte ← Radarr/Sonarr status
│   └── libraries/+page.svelte  ← Library management
└── api/                        ← SvelteKit API routes (proxy/BFF)
    └── auth/
        └── callback/[provider]/+server.ts  ← OIDC callback
```

---

## Empfohlene Frontend-TypeScript-Client-Strategie

### Option A: openapi-typescript + openapi-fetch (Empfohlen)
```bash
npx openapi-typescript http://localhost:8096/api/openapi.yaml -o src/lib/api/schema.d.ts
```
- Type-safe, zero runtime overhead
- Nutzt die bereits vorhandene OpenAPI-Spec
- Automatische Typen-Generierung bei API-Änderungen

### Option B: ogen generiert bereits einen Go-Client
- Nicht relevant für SvelteKit

### Option C: Manuell
- Nicht empfohlen bei 187 Endpoints

---

## Checkliste vor Frontend-Start

- [ ] **P0:** Property Naming vereinheitlichen (snake_case ODER camelCase)
- [ ] **P0:** Alle List-Endpoints mit `{ items, total }` Envelope
- [ ] **P1:** Pagination auf `limit`/`offset` standardisieren
- [ ] **P1:** HTTP Methods angleichen (Progress: POST vs PUT)
- [ ] **P2:** OpenAPI Tags aufräumen
- [ ] **P2:** Dupliziertes TVShowListResponse entfernen
- [ ] **P2:** Sort-Parameter vereinheitlichen
- [ ] Frontend-Repo initialisieren (SvelteKit 2 + Svelte 5 + Tailwind CSS 4 + shadcn-svelte)
- [ ] openapi-typescript für Type-Generierung einrichten
- [ ] Auth-Flow implementieren (Cookie-Auth + CSRF)
- [ ] SSE-Client für Real-Time Events
- [ ] HLS.js Integration für Video-Playback
