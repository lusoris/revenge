# Revenge Codebase — Vollständige Analyse

**Datum:** 2026-02-10
**Branch:** develop
**Scope:** Gesamter Backend-Code, API, Infra, Services, Content Module

---

## Zusammenfassung

**Gesamtbewertung: Solider, produktionsnaher Backend-Code mit wenigen echten Problemen.**

Die Architektur ist sauber geschichtet, der Code ist gut getestet, und die API ist zu ~95% bereit für ein SvelteKit-Frontend. Es gibt aber konkrete Issues die behoben werden sollten.

### Zahlen auf einen Blick

| Metrik | Wert |
|---|---|
| Go-Dependencies (direkt) | ~50 |
| API Endpoints | 187 |
| Handler-Implementierung | 188/188 (100%) |
| DB Migrations | 37 (keine Lücken) |
| sqlc Queries | ~200+ (shared) + ~50 (movie) + ~100 (tvshow) |
| Service-Packages | 16 + 3 Content Module |
| Middleware-Komponenten | 11 |
| River Job Workers | 9 |
| Metadata Provider | 12 |
| Test-Dateien | 100+ |
| Hand-written Handler Lines | ~9,520 |
| Generated Code (ogen) | ~152,000 Lines |

---

## 1. FEHLER & KRITISCHE ISSUES

### CRITICAL: `RefreshPersonWorker` ist ein verkleideter No-Op

**Location:** `internal/service/metadata/jobs/workers.go:188-217`

Der Worker akzeptiert Jobs, loggt eine Warning und gibt `nil` zurück. Jobs werden **stillschweigend verschluckt** ohne Fehler. Sollte entweder `ErrNotImplemented` returnen oder gar nicht erst als Worker registriert sein.

```go
func (w *RefreshPersonWorker) Work(ctx context.Context, job *river.Job[RefreshPersonArgs]) error {
    w.logger.Warn("person metadata refresh not yet implemented — job accepted but no-op", ...)
    return nil  // silently drops the job
}
```

### WARNING: `panic()` in Production Code

| Location | Trigger |
|---|---|
| `internal/infra/jobs/river.go:187` | `Subscribe()` panics wenn River Client nil ist |
| `internal/validate/convert.go:19-53` | `Must*` Funktionen paniken bei Overflow (Go-Konvention, aber risky) |
| `internal/config/loader.go:123-126` | `MustLoad` panics bei Config-Fehler (standard für main init) |

### WARNING: Ungenutzter DB-Pool in `migrate.go`

**Location:** `cmd/revenge/migrate.go:41-46`

Erstellt einen `database.NewPool`, pingt ihn, defers `pool.Close()`, nutzt den Pool aber nie. Die `Migrate*` Funktionen öffnen eigene Connections via `sql.Open`. **Dead Code.**

---

## 2. STUBS & UNIMPLEMENTIERTES

### Echte Stubs im Code

| Was | Location | Status |
|---|---|---|
| **QAR-Modul** | `internal/content/qar/` | Nur `SELECT 1 AS placeholder`. Kein Service, kein Handler, keine Tests |
| **RefreshPersonWorker** | `internal/service/metadata/jobs/workers.go` | Akzeptiert Jobs, tut aber nichts |
| **Migration-Dirs** `qar/`, `tvshow/` | `internal/infra/database/migrations/` | Enthalten nur `.gitkeep` — alle Migrations laufen über `shared/` |

### Nicht-implementierte Features (aus TODO.md)

**Content Module (0/7):**
- Music, Audiobooks, Books, Podcasts, Comics, Photos, Live TV

**Services:**
- Transcoding-Engine (Blackbeard), Analytics-Dashboard, Grants (fine-grained sharing), Fingerprint (media identification)

**Integrations:**
- Lidarr, Whisparr, Authelia/Authentik/Keycloak SSO, Last.fm, ListenBrainz, MusicBrainz/Spotify

**Features:**
- Collections/Playlists, Watch Next/Continue Watching, Skip Intro/Credits Detection, SyncPlay, Trickplay (timeline thumbnails), Release Calendar, Content Request System

**Infra:**
- Circuit Breaker (gobreaker Dep existiert, ungenutzt), Request Coalescing, Cache Warming, Service Watchdog

---

## 3. UNGEREIMTHEITEN & INKONSISTENZEN

### 3.1 API Naming: snake_case vs camelCase (HOCH)

Die kritischste Inkonsistenz. Die OpenAPI-Spec mischt zwei Konventionen:

| Konvention | Betroffene Ressourcen |
|---|---|
| **snake_case** | Auth, MFA, Movies, TV Shows, Users, Settings, Sessions, Search, API Keys, Metadata |
| **camelCase** | OIDC, Admin OIDC, Activity, Libraries, Library Scans, Radarr/Sonarr, Playback |

Beispiel: `access_token` vs `accessToken`, `created_at` vs `createdAt`.

### 3.2 Pagination: 3 verschiedene Patterns

| Pattern | Verwendet von |
|---|---|
| `limit` / `offset` | Movies, TV Shows, Admin Users, Activity |
| `page` / `per_page` | Search (Typesense) |
| **Keine Pagination** | 17+ Endpoints die Arrays ohne `total` zurückgeben |

### 3.3 List Response Envelopes: Kein Standard

| Response Schema | Wrapper-Feld | `total` vorhanden? |
|---|---|---|
| `MovieListResponse` | `items` | Ja |
| `TVSeriesListResponse` | `items` | Ja |
| `AdminUserListResponse` | `users` | Nein |
| `SessionListResponse` | `sessions` | Nein |
| `PolicyListResponse` | `policies` | Nein |
| `APIKeyListResponse` | `keys` | Nein |
| `OIDCProviderListResponse` | `providers` | Nein |
| `LibraryListResponse` | `libraries` | Nein |
| `ActivityLogListResponse` | `entries` | Ja |
| `RolesResponse` | `roles` | Nein |

### 3.4 HTTP Method Inkonsistenz

- Movie Progress update: **POST** auf `/movies/{id}/progress`
- Episode Progress update: **PUT** auf `/tvshows/episodes/{id}/progress`

### 3.5 Service-Pattern Inkonsistenz

| Pattern | Services |
|---|---|
| Repository Interface + Concrete Service | activity, apikeys, auth, library, oidc, session, user |
| Repository Interface + Service Interface | settings, movie, tvshow, metadata, storage |
| Kein Repository (raw sqlc) | analytics, mfa |
| Kein Repository (external) | rbac (Casbin), search (Typesense), notification (in-memory) |
| Kein Interface, kein Mock | email |

### 3.6 Fehlende Mocks

`email`, `metadata` (Service Interface), `mfa`, `rbac` und `search` haben keine Mocks — das erschwert isoliertes Testen in abhängigen Services.

### 3.7 Sort-Parameter Naming

- Movies/TV: `order_by` (Werte: `title, year, added, rating`)
- Search: `sort_by`
- TV `order_by` Werte: `created_at, title, first_air_date, vote_average, popularity`

### 3.8 Duplizierte Response Schemas

`TVShowListResponse` und `TVSeriesListResponse` — beide wrappen `TVSeries[]`. Eins sollte entfernt werden.

### 3.9 Fehlende Tag-Definitionen

8 Tags werden in Endpoints genutzt, sind aber nicht im Top-Level `tags` Array der OpenAPI-Spec deklariert: `sessions`, `rbac`, `apikeys`, `OIDC`, `Users Admin`, `OIDC Admin`, `Activity`, `Libraries`.

Tag-Casing inkonsistent: lowercase (`health`, `movies`) vs Title Case (`Integrations Admin`) vs ALL CAPS (`OIDC`).

### 3.10 Bare-Array Responses (kein `total` Count)

Folgende Endpoints akzeptieren `limit`/`offset` aber geben nackte Arrays zurück — **Pagination-UI unmöglich**:

- `GET /api/v1/movies/search`
- `GET /api/v1/movies/continue-watching`
- `GET /api/v1/movies/watch-history`
- `GET /api/v1/movies/{id}/files`, `/genres`
- `GET /api/v1/tvshows/search`
- `GET /api/v1/tvshows/continue-watching`
- `GET /api/v1/tvshows/episodes/recent`, `/upcoming`
- `GET /api/v1/tvshows/{id}/seasons`, `/episodes`, `/genres`, `/networks`
- `GET /api/v1/tvshows/seasons/{id}/episodes`
- `GET /api/v1/tvshows/episodes/{id}/files`
- `GET /api/v1/collections/{id}/movies`
- `GET /api/v1/genres`

### 3.11 Overlapping Utility Code

`validate.SafeInt32` (gibt Error zurück) vs `util.SafeIntToInt32` (saturiert) — ähnliche Namen, verschiedene Semantik, keine Cross-Reference-Docs.

### 3.12 Refresh Endpoints ohne Job-ID Response

`POST .../movies/{id}/refresh` und `POST .../tvshows/{id}/refresh` geben 202 ohne Body zurück. Die Search Reindex-Endpoints geben `{ message, job_id }` zurück — das Refresh-Pattern sollte folgen.

### 3.13 MFA Inline-Schemas

Viele MFA-Endpoints definieren ad-hoc `{ success, message }` Inline-Schemas statt `$ref` auf wiederverwendbare Komponenten.

### 3.14 Typesense-Syntax Leakage

`filter_by` in Search-Endpoints exposed Typesense DSL (`genres:=Action && year:>=2020`). Das Frontend müsste Typesense-Syntax kennen.

---

## 4. STÄRKEN

| Bereich | Bewertung | Details |
|---|---|---|
| **Architektur** | Exzellent | Saubere Schichtung, keine zirkulären Abhängigkeiten, fx-DI gut eingesetzt |
| **Migrations** | Exzellent | 37 Migrations, keine Lücken, alle mit up+down, getestet |
| **API Coverage** | Perfekt | 188/188 Endpoints implementiert — Zero Gaps |
| **Middleware** | Exzellent | 11 Komponenten: CORS, CSP, CSRF, Rate Limiting (in-memory + Redis), Request ID, Cache-Control, Cookie Auth |
| **Auth** | Exzellent | JWT + Refresh + MFA (TOTP + WebAuthn + Backup Codes) + OIDC + Cookie Auth + CSRF |
| **Caching** | Exzellent | L1 (otter/W-TinyLFU) + L2 (Redis/Dragonfly), graceful degradation, voll instrumentiert |
| **Playback** | Voll funktional | HLS-Streaming mit FFmpeg, Quality Profiles, Audio/Subtitle Tracks |
| **Search** | Exzellent | Typesense über 6 Entity-Typen mit Autocomplete und Facets |
| **Observability** | Umfassend | Prometheus Metrics für HTTP, DB, Cache, Jobs, Pool Stats, pprof |
| **Error Handling** | Sauber | Sentinel-Pattern mit `go-faster/errors` für Stack Traces |
| **Crypto** | Stark | AES-256-GCM, Argon2id, bcrypt-Migration, UUIDv7 |
| **Integrations** | Komplett | Radarr + Sonarr voll (Client, Sync, Mapper, Webhooks, Jobs, Rate Limiting) |
| **SSE** | Komplett | Server-Sent Events mit Auth, Category-Filtering, 30s Keepalive |
| **Config** | Flexibel | koanf v2, YAML + ENV override, sane defaults |
| **Health Probes** | K8s-ready | liveness/readiness/startup mit Dependency-Checks |
| **Security** | Solid | Security Headers (CSP, X-Frame, Referrer-Policy, Permissions-Policy), CORS, CSRF, Rate Limiting |

---

## 5. API ENDPOINT ÜBERSICHT (187 Operationen)

### Health (3)
`GET /healthz`, `GET /readyz`, `GET /startupz`

### Auth (9)
Register, Login, Logout, Refresh, Verify Email, Resend Verification, Forgot/Reset/Change Password

### MFA + WebAuthn (15)
TOTP Setup/Verify/Delete, Backup Codes Generate/Regenerate, Enable/Disable MFA, WebAuthn Register/Login Begin/Finish, Credential CRUD

### Movies (19)
List, Search, Recently Added, Top Rated, Continue Watching, Watch History, Stats, Get/Files/Cast/Crew/Genres/Collection/Similar/Progress CRUD/Watched/Refresh

### Collections (2)
Get Collection, Get Collection Movies

### TV Shows (26)
List, Search, Recently Added, Continue Watching, Stats, Recent/Upcoming Episodes, Series CRUD, Seasons, Episodes, Cast/Crew/Genres/Networks, Watch Stats, Next Episode, Refresh, Season/Episode Detail, Episode Files/Progress/Watched/Bulk Watched

### Metadata/TMDb Proxy (23)
Movie/TV/Person search, Detail/Credits/Images/Similar/Recommendations/External IDs, Collection Detail, TV Season/Episode Detail, Person Credits/Images, Providers List

### Search/Typesense (9)
Movies/TVShows Search + Autocomplete + Facets, Multi Search, Reindex (Movies + TVShows)

### Settings (7), Users (6), Sessions (6), RBAC (12), API Keys (4), OIDC (6)

### Admin: Users (2), OIDC (8), Activity (5), Radarr (4), Sonarr (4)

### Libraries (10), Webhooks (2), Playback/HLS (3), Images (1), Genres (1), Events/SSE (1)

---

## 6. INFRASTRUKTUR-DETAILS

### Database
- **PostgreSQL 18+** mit pgxpool, Connection Pool Tuning
- **37 Migrations** (000001–000037), continous numbering, up+down
- **sqlc** Codegen: 4 Targets (shared, movie, tvshow, qar-placeholder)
- **River** PostgreSQL-native Job Queue mit eigenem Migrator

### Caching
- **L1:** otter v2 (W-TinyLFU, 10k entries, 5min TTL default)
- **L2:** rueidis (Redis/Dragonfly, server-assisted client caching, 16MB per connection)
- **~30+ Cache Key Prefixes** für alle Domains
- Graceful Degradation: Alles funktioniert wenn Cache disabled

### Jobs (River)
- 5-Priority-Queue System
- 9 Worker Types: Movie (4), TV (5), Shared (2)
- Progress Reporting via `river.JobProgress`
- Leader-only Execution via Raft

### Raft
- HashiCorp Raft für Leader Election (Clustering)
- Potentiell heavyweight — pg advisory locks wären einfacher

### Observability
- Prometheus Metrics: HTTP, DB Queries, Cache Hit/Miss, Job Execution, Pool Stats
- pprof endpoints
- Structured Logging (slog + tint dev / JSON prod)
- Slow Query Detection mit konfigurierbarem Threshold

---

## 7. SERVICE-ARCHITEKTUR MATRIX

| Service | Service Type | Repository | Interface | Cached | Mocks | Tests |
|---|---|---|---|---|---|---|
| activity | Concrete | Yes | Logger iface | No | mockery | 10 files |
| analytics | Worker only | No (raw sqlc) | No | No | No | 1 file |
| apikeys | Concrete | Yes | No | No | mockery | 4 files |
| auth | Concrete | Yes | TokenManager | No | mockery | 5 files |
| email | Concrete | No | No | No | **No** | 1 file |
| library | Concrete | Yes | No | Yes | mockery | 4 files |
| metadata | Private struct | Provider pattern | **Yes** | No | **No** | 1 file |
| mfa | Manager | No (raw sqlc) | No | No | **No** | 6 files |
| notification | Dispatcher | No (in-memory) | Agent iface | No | mockery | 3 files |
| oidc | Concrete | Yes | No | No | mockery | 5 files |
| rbac | Concrete | Casbin adapter | No | Yes | **No** | 6 files |
| search | Concrete(s) | No (Typesense) | No | Yes | **No** | 9 files |
| session | Concrete | Yes | No | Yes | mockery | 5 files |
| settings | Private struct | Yes | **Yes** | Yes | mockery | 4 files |
| storage | Impls | N/A | **Yes** | No | hand-written | 1 file |
| user | Concrete | Yes | No | Yes | mockery | 6 files |
| movie | Private struct | Yes | **Yes** | Yes | mockery | 9+ files |
| tvshow | Private struct | Yes | **Yes** | Yes | partial | 2 files |
| qar | **STUB** | N/A | N/A | N/A | N/A | 0 files |

---

## 8. DEPENDENCY GRAPH (vereinfacht)

```
Layer 1 (leaf):    config, logging, errors, validate, util, crypto
Layer 2 (infra):   database, cache, search, jobs, raft, health, image, observability
Layer 3 (domain):  content/movie, content/tvshow, content/shared
Layer 4 (service): auth, session, user, rbac, apikeys, mfa, oidc, activity, etc.
Layer 5 (integ):   radarr, sonarr (→ content repos)
Layer 6 (stream):  playback (→ content services)
Layer 7 (api):     api handlers + middleware + sse (→ services via fx)
Layer 8 (app):     app/module.go (wires everything)
```

**Keine zirkulären Abhängigkeiten.** Movie↔Metadata-Zyklus wird sauber über Interface-Adapter (`movie.MetadataQueue`) aufgelöst.
