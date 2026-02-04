# Comprehensive TODO - v0.3.0 MVP

**Last Updated**: 2026-02-04 02:45
**Current Focus**: Phase 8 (Test Coverage 80%) 🔄
**Status**: Backend Complete ✅ → TMDb Complete ✅ → Library Provider Complete ✅ → River Jobs Complete ✅ → Typesense Complete ✅ → Radarr Complete ✅ → Rate Limiting Complete ✅ → **Feature Gaps Identified ✅** → Pre-Frontend 🟡 (Phase 7/8 ✅) → Tests 🟡 (5.8% → target 80%)

**Reports erstellt**:
- [FEATURE_GAP_ANALYSIS.md](./FEATURE_GAP_ANALYSIS.md) - Umfassende Feature-Analyse

---

## Pre-MFA: Quick Fixes

### Standardize Health Endpoints (30 minutes) ✅ COMPLETE
**Previous**: `/health/live`, `/health/ready`, `/health/startup`
**Standard**: `/healthz`, `/readyz`, `/startupz` (Kubernetes convention)

- [x] Update OpenAPI spec: Rename endpoints to `/healthz`, `/readyz`, `/startupz` ✅
- [x] Regenerate ogen code ✅
- [x] Update integration tests ✅
- [x] Update API tests ✅

**Commit**: 39fd6653c0 - refactor(api): standardize health endpoints to Kubernetes conventions

**References**:
- [Kubernetes Liveness/Readiness Probes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/)
- [GKE Health Check Standards](https://cloud.google.com/kubernetes-engine/docs/concepts/ingress#health_checks)

---

## Current Sprint: MFA Implementation (10-20 hours)

> **Design Complete**: [MFA.md](../docs/dev/design/services/MFA.md)
> **Status**: Ready to implement

### Phase 1: Foundation (2-3 hours) ✅ COMPLETE
- [x] **Database Migrations** ✅
  - [x] `000016_create_user_totp_secrets.up.sql`
  - [x] `000017_create_webauthn_credentials.up.sql`
  - [x] `000018_create_mfa_backup_codes.up.sql`
  - [x] `000019_create_user_mfa_settings.up.sql`
  - [x] Migrations tested and applied successfully

- [x] **Encryption Service** (`internal/crypto/encryption.go`) ✅
  - [x] Implement AES-256-GCM encryption
  - [x] Encrypt/Decrypt helpers with nonce handling
  - [x] Tests with 85.1% coverage
  - [x] Extracted from OIDC service for reuse

- [x] **SQLC Queries** (`internal/infra/database/queries/shared/mfa.sql`) ✅
  - [x] TOTP queries (Create, Get, Verify, Enable/Disable, UpdateLastUsed)
  - [x] WebAuthn queries (Create, List, Get, UpdateCounter, CloneDetection)
  - [x] Backup codes queries (Create copyfrom, GetUnused, Use, Count, DeleteAll)
  - [x] Settings queries (CRUD, Enable/Disable methods, TrustedDevices)
  - [x] Combined status queries (GetUserMFAStatus, HasAnyMFAMethod)
  - [x] Generated SQLC code compiles successfully (30+ operations)

**Commits**:
- 782a470b0d: feat(crypto): add shared AES-256-GCM encryption service
- 5e1913a5b3: feat(mfa): add database migrations for MFA tables
- aa3c2b6b7d: feat(mfa): add SQLC queries for MFA operations

### Phase 2: TOTP Implementation (2-3 hours) ✅ COMPLETE
- [x] **TOTP Service** (`internal/service/mfa/totp.go`) ✅
  - [x] Generate secret (20 bytes/160 bits, base32-encoded)
  - [x] Generate QR code (PNG, 256x256, otpauth://totp/...)
  - [x] Verify TOTP code (RFC 6238, 30s window, ±1 step skew)
  - [x] Store encrypted secret (AES-256-GCM)
  - [x] Enable/disable/delete TOTP
  - [x] Auto-enable on first successful verification

- [x] **Tests** (`internal/service/mfa/totp_test.go`) ✅
  - [x] Unit tests for TOTP generation/verification
  - [x] Test time skew tolerance (±30s)
  - [x] Test secret encryption/decryption
  - [x] Test code format (6 digits)
  - [x] Test deterministic generation
  - [x] Test uniqueness across secrets
  - [x] Integration test stubs (database required)

**Features**:
- SHA1 algorithm (most compatible with authenticator apps)
- 6-digit codes (standard)
- 30-second time window
- Encrypted secret storage with AES-256-GCM

**Commit**: 3a7464f322 - feat(mfa): implement TOTP service with encryption

### Phase 3: WebAuthn (3-4 hours) ✅ COMPLETE
- [x] **WebAuthn Service** (`internal/service/mfa/webauthn.go`) ✅
  - [x] Use `github.com/go-webauthn/webauthn` v0.11.2
  - [x] Registration flow (BeginRegistration, FinishRegistration)
  - [x] Authentication flow (BeginLogin, FinishLogin)
  - [x] Credential storage (credential_id, public_key, AAGUID, transports)
  - [x] Clone detection (sign counter verification with rollback detection)
  - [x] Multiple credentials per user support
  - [x] Credential management (list, rename, delete)

- [x] **WebAuthnUser Interface** ✅
  - [x] Implements `webauthn.User` interface
  - [x] Dynamic credential loading from database
  - [x] UUID-based user identification

- [x] **Tests** (`internal/service/mfa/webauthn_test.go`) ✅
  - [x] Unit tests for service initialization
  - [x] Test WebAuthnUser interface compliance
  - [x] Test transport conversion
  - [x] Test session data serialization (JSON)
  - [x] Integration test stubs (database + mock WebAuthn responses required)
  - [x] Test scenarios: lifecycle, clone detection, multiple credentials

**Features**:
- W3C WebAuthn Level 3 compliance
- Discoverable credentials support
- User verification required
- Clone detection with counter rollback prevention
- Multi-device support (USB, NFC, BLE, Internal)

**Commit**: f0c3da69cf - feat(mfa): implement WebAuthn service with clone detection

### Phase 4: Backup Codes + Manager (2-3 hours) ✅ COMPLETE
- [x] **Backup Codes Service** (`internal/service/mfa/backup_codes.go`) ✅
  - [x] Generate 10 backup codes (8 bytes random → 16 hex chars)
  - [x] Format codes (XXXX-XXXX-XXXX-XXXX for UX)
  - [x] Hash codes (bcrypt cost 12)
  - [x] Verify backup code with constant-time comparison
  - [x] Mark as used (single-use with IP tracking)
  - [x] Regenerate codes (delete old, generate new)
  - [x] Get remaining unused count

- [x] **MFA Manager Service** (`internal/service/mfa/manager.go`) ✅
  - [x] Unified MFA coordinator (TOTP, WebAuthn, Backup Codes)
  - [x] GetStatus - aggregated MFA status for user
  - [x] HasAnyMethod - check if any MFA configured
  - [x] RequiresMFA - check if MFA enforcement enabled
  - [x] EnableMFA/DisableMFA - toggle MFA requirement
  - [x] VerifyTOTP/VerifyBackupCode - unified verification
  - [x] RemoveAllMethods - cleanup all MFA data

- [x] **Tests** ✅
  - [x] `backup_codes_test.go` - code generation, formatting, normalization, constant-time comparison
  - [x] `manager_test.go` - data structures, integration test stubs
  - [x] All unit tests passing (21/21)
  - [x] Integration test stubs for database-dependent tests

- [x] **Auth Service Integration** (`internal/service/auth/mfa_integration.go`) ✅
  - [x] Add MFA verification to login flow (MFAAuthenticator)
  - [x] Return MFA challenge if enabled (CheckMFARequired)
  - [x] Verify MFA response (VerifyMFA - TOTP/backup codes)
  - [x] Session enhancement (mark MFA-verified with CompleteMFALogin)
  - [x] LoginWithMFA flow implementation
  - [x] GetSessionMFAInfo for verification status

- [x] **Session Updates** (`migrations/000020_add_mfa_to_sessions.up.sql`) ✅
  - [x] Add `mfa_verified` boolean to session
  - [x] Add `mfa_verified_at` timestamp
  - [x] SQLC queries for MFA session tracking
  - [x] Index on (user_id, mfa_verified)

**Security Features**:
- Constant-time comparison (prevents timing attacks)
- Argon2id hashing (replaced bcrypt for consistency)
- One-time use enforcement with database constraints
- IP tracking for audit trail
- Case-insensitive, dash-tolerant code normalization

**Commits**:
- e72d7f7ff9 - feat(mfa): implement backup codes and MFA manager services
- 3a1ae626ac - feat(mfa): integrate MFA with auth service and unify password hashing

### Phase 5: Production Hardening & API Integration (2-3 hours) ✅ COMPLETE
- [x] **MFA Module** (`internal/service/mfa/module.go`) ✅
  - [x] fx.Module with providers for TOTP, BackupCodes, WebAuthn, MFAManager
  - [x] Configuration from config (issuer, RP ID, origins)
  - [x] Integrated into app module

- [x] **API Handlers** (`internal/api/handler_mfa.go`) ✅
  - [x] `GET /api/v1/mfa/status` - Get MFA configuration status
  - [x] `POST /api/v1/mfa/totp/setup` - Generate secret + QR code
  - [x] `POST /api/v1/mfa/totp/verify` - Verify and enable TOTP
  - [x] `DELETE /api/v1/mfa/totp` - Disable TOTP
  - [x] `POST /api/v1/mfa/backup-codes/generate` - Generate 10 backup codes
  - [x] `POST /api/v1/mfa/backup-codes/regenerate` - Delete old, generate new
  - [x] `POST /api/v1/mfa/enable` - Turn on MFA requirement
  - [x] `POST /api/v1/mfa/disable` - Turn off MFA requirement
  - [x] Update `api/openapi/openapi.yaml` with all endpoints
  - [x] Regenerated ogen code
  - [x] Fixed ogen type handling (url.URL, []byte, operation-specific errors)
  - [x] GetUserIDFromContext implementation

- [x] **Server Integration** (`internal/api/server.go`) ✅
  - [x] Wire MFA services into API server via dependency injection
  - [x] Delegate MFA methods from main Handler to MFAHandler
  - [x] All 8 endpoints operational

- [x] **Testing** ✅
  - [x] Unit tests for all MFA services (TOTP, WebAuthn, BackupCodes, Manager)
  - [x] All tests passing (21/21 unit tests)
  - [x] Integration test stubs for database-dependent flows
  - [x] Application compiles and builds successfully

**Commits**:
- 5cee136167 - feat(mfa): add MFA API handlers and integrate with server
- e8928e6f9d - fix(mfa): implement GetUserIDFromContext using existing context helper

**Pending (Deferred)**:
- [x] Rate Limiting (per-endpoint) ✅ - Added with configurable auth/global tiers
- [ ] Comprehensive audit logging
- [ ] WebAuthn API endpoints (registration/login flows)
- [ ] Full end-to-end integration tests with database
- [ ] API documentation updates
- [ ] User guide for MFA enrollment

**Notes**:
- WebAuthn service implemented but API endpoints deferred (needs more complex flow)
- Basic structure complete, production features can be added incrementally
- All core MFA functionality working (TOTP + backup codes + session tracking)

---

## v0.3.0 MVP Scope (After MFA)

> **Design Complete**: All required design docs exist (see TODO_v0.3.0.md)
> **Focus**: Movie Module + TMDb + Radarr + Typesense + Frontend

### Movie Module (Backend)

#### Database Schema ✅ COMPLETE
- [x] `public.movies` table (UUID v7, title, year, runtime, overview, tmdb_id, imdb_id, poster/backdrop paths) ✅
- [x] `public.movie_genres` table ✅
- [x] `public.movie_credits` table (cast/crew combined with type field) ✅
- [x] `public.movie_files` table ✅
- [x] `public.movie_watched` table (watch progress with generated percent column) ✅
- [x] `public.movie_collections` + `public.movie_collection_members` tables ✅
- [x] Indexes on tmdb_id, imdb_id, title (including trigram for fuzzy search) ✅

**Commit**: a72c8c877c - feat(db): add movie module database schema with 6 migrations

#### Entity Layer ✅ COMPLETE
- [x] Movie, MovieFile, MovieCredit, MovieCollection, MovieGenre, MovieWatched structs ✅
- [x] Parameter types (CreateMovieParams, UpdateMovieParams, etc.) ✅
- [x] Filter types (ListFilters, SearchFilters) ✅

**Commit**: 4f6441a4fb - feat(movie): add repository layer with PostgreSQL implementation

#### Repository Layer ✅ COMPLETE
- [x] Interface definition (Repository with all CRUD operations) ✅
- [x] PostgreSQL implementation (NewPostgresRepository) ✅
- [x] CRUD operations (Get, Create, Update, Delete) ✅
- [x] List with filters (genre, year, recently added, top rated) ✅
- [x] Search by title (trigram similarity) ✅
- [x] Watch progress operations (Create/Update, Get, Delete, Stats) ✅
- [x] SQLC queries (movies.sql with 40+ operations) ✅

**Commits**:
- 31d35ed992 - feat(movie): add SQLC queries for movie module
- 4f6441a4fb - feat(movie): add repository layer with PostgreSQL implementation

#### Service Layer ✅ COMPLETE
- [x] Get movie by ID ✅
- [x] List movies (paginated with filters) ✅
- [x] Search movies (title fuzzy search) ✅
- [x] Update watch progress (auto-complete at 90%) ✅
- [x] Get continue watching ✅
- [x] Get recently added ✅
- [x] Trigger metadata refresh (placeholder) ✅
- [x] Input validation (title required, file uniqueness) ✅
- [x] Business logic (existence checks, completion calculation) ✅

**Commit**: 5ac9fe3131 - feat(movie): add service layer and fx module

#### Library Provider ✅ COMPLETE
- [x] Scanner (walk filesystem, parse filenames, extract title/year) ✅
- [x] Matcher (TMDb search, confidence scoring, create movies) ✅
- [x] Service (ScanLibrary, RefreshMovie workflows) ✅
- [x] Filename parsing patterns (Title (YEAR), Title.YEAR) ✅
- [x] Quality marker removal (1080p, BluRay, x264, etc.) ✅
- [x] Video extensions support (13 formats) ✅
- [x] Confidence algorithm (title similarity + year match + popularity) ✅

**Commit**: d8789fc4d3 - feat(movie): add Library Provider for file scanning and matching

#### API Handlers ✅ COMPLETE
- [x] `GET /api/v1/movies` (list, paginated) ✅
- [x] `GET /api/v1/movies/:id` ✅
- [x] `GET /api/v1/movies/:id/files` ✅
- [x] `GET /api/v1/movies/:id/cast` ✅
- [x] `GET /api/v1/movies/:id/crew` ✅
- [x] `GET /api/v1/movies/:id/genres` ✅
- [x] `POST /api/v1/movies/:id/progress` ✅
- [x] `GET /api/v1/movies/:id/progress` ✅
- [x] `DELETE /api/v1/movies/:id/progress` ✅
- [x] `POST /api/v1/movies/:id/watched` ✅
- [x] `POST /api/v1/movies/:id/refresh` ✅
- [x] `GET /api/v1/movies/recently-added` ✅
- [x] `GET /api/v1/movies/top-rated` ✅
- [x] `GET /api/v1/movies/continue-watching` ✅
- [x] `GET /api/v1/movies/watch-history` ✅
- [x] `GET /api/v1/movies/stats` ✅
- [x] OpenAPI spec integration (ogen) ✅
- [x] Wire handlers into API server ✅
- [x] Type converters (domain ↔ ogen) ✅

**Commits**:
- f18891b880 - feat(movie): add HTTP handlers and integrate into app
- 59fb5d1350 - feat: Add Movie Module backend foundation

#### River Jobs ✅ COMPLETE
- [x] MovieMetadataRefreshJob (refresh TMDb metadata for movie by ID) ✅
- [x] MovieLibraryScanJob (scan library paths for new/changed files) ✅
- [x] MovieFileMatchJob (stub - match single file to movie) ✅
- [x] Worker registration and FX module integration ✅
- [x] Config added (movie.tmdb and movie.library settings) ✅

**Commit**: 033accd17b - feat(movie): add River Jobs for background processing

#### Tests � IN PROGRESS (46.7% Coverage)
- [x] Unit tests for service (mock repository)
- [x] Unit tests for handler
- [x] Unit tests for TMDb mapper
- [x] Lint fixes (handler.go, library_scanner.go, tmdb_client_test.go)
- [ ] Integration tests with database (target: 80%+ coverage)

### Collection Support ✅ COMPLETE (Database + Logic)
- [x] `public.movie_collections` table ✅
- [x] `public.movie_collection_members` junction table ✅
- [x] Collection repository methods ✅
- [x] Collection service (Get, GetMovies) ✅
- [x] API handlers (`GET /collections/:id`, `/collections/:id/movies`) ✅
- [x] OpenAPI spec integration ✅

### Metadata Service (TMDb) ✅ COMPLETE

#### TMDb Client ✅ COMPLETE
- [x] API key configuration ✅
- [x] Rate limiting (40 req/10s) ✅
- [x] Retry with backoff ✅
- [x] Response caching (sync.Map with TTL) ✅
- [x] TMDb types (Movie, Credits, Images, Collections) ✅
- [x] TMDb client (SearchMovies, GetMovie, GetCredits, GetImages, GetCollection) ✅
- [x] TMDb mapper (TMDb → domain types) ✅
- [x] Metadata service (unified interface) ✅
- [x] Image URL construction and downloading ✅
- [x] Proxy/VPN support ✅

**Commit**: a70c7b57e2 - feat(movie): add TMDb metadata service

#### TMDb Service ✅ COMPLETE (included in TMDb Client above)
- [x] Search movie ✅
- [x] Get movie details ✅
- [x] Get movie credits (cast/crew) ✅
- [x] Get movie images ✅
- [x] Get similar movies ✅
- [x] Get collection details ✅

#### Image Handler ✅ COMPLETE (internal/infra/image)
- [x] Poster download/cache ✅
- [x] Backdrop download/cache ✅
- [x] Profile image download/cache ✅
- [x] Image proxy endpoint (`GET /api/v1/images/{type}/{size}/{path}`) ✅

#### API Handlers ✅ COMPLETE (internal/api/handler_metadata.go)
- [x] `GET /api/v1/metadata/search/movie?q=` ✅
- [x] `GET /api/v1/metadata/movie/:tmdbId` ✅
- [x] `GET /api/v1/metadata/collection/:tmdbId` ✅
- [x] `GET /api/v1/images/:type/:size/:path` (proxy) ✅
- [x] `GET /api/v1/movies/:id/similar` ✅

#### Tests
- [ ] Unit tests with mock API
- [ ] Integration tests (optional, needs API key)

### Search Service (Typesense) ✅ COMPLETE

#### Typesense Setup ✅ COMPLETE
- [x] Client configuration (`internal/infra/search/module.go`) ✅
- [x] Collection schemas (`internal/service/search/movie_schema.go`) ✅
- [x] Index management (create collection, bulk indexing, reindex) ✅

#### Movie Collection Schema ✅ COMPLETE
- [x] Define schema (title, original_title, overview, year, genres, cast, director, rating, added_at) ✅
- [x] Facets: genres, year, status, directors, resolution, has_file ✅
- [x] Sortable: popularity, vote_average, release_date, library_added_at ✅
- [x] Infix search on title fields for partial matching ✅

#### Search Service ✅ COMPLETE (`internal/service/search/movie_service.go`)
- [x] Index movie ✅
- [x] Remove from index ✅
- [x] Search movies (full-text) ✅
- [x] Faceted search (genre, year, status, resolution) ✅
- [x] Autocomplete ✅
- [x] Bulk indexing ✅
- [x] Full reindex ✅

#### API Handlers ✅ COMPLETE (`internal/api/handler_search.go`)
- [x] `GET /api/v1/search/movies` (full-text search with facets) ✅
- [x] `GET /api/v1/search/movies/autocomplete` ✅
- [x] `GET /api/v1/search/movies/facets` ✅
- [x] `POST /api/v1/search/reindex` (admin-only) ✅

#### River Jobs ✅ COMPLETE (`internal/content/movie/moviejobs/search_index.go`)
- [x] MovieSearchIndexWorker - Index/remove single movie ✅
- [x] SearchReindexJob - Full reindex ✅
- [x] FX module integration ✅

#### Tests ✅ COMPLETE
- [x] Unit tests (`internal/service/search/movie_service_test.go` - 18+ tests) ✅
- [x] Worker tests (`internal/content/movie/moviejobs/search_index_test.go`) ✅
- [ ] Integration tests with Typesense container (deferred)

**Commits**:
- 6a8701c12f - feat(search): add Typesense movie search service and API endpoints
- 19b3f209e9 - feat(search): add River job for search index operations

### Radarr Integration ✅ CORE COMPLETE

#### Radarr Client ✅ COMPLETE
- [x] API v3 implementation (`internal/integration/radarr/client.go`) ✅
- [x] Authentication (API key) ✅
- [x] Error handling (`internal/integration/radarr/errors.go`) ✅
- [x] Type definitions (`internal/integration/radarr/types.go`) ✅
- [x] Rate limiting and caching ✅

#### Radarr Service ✅ COMPLETE
- [x] Get all movies ✅
- [x] Get movie by ID ✅
- [x] Get movie files ✅
- [x] Sync library (Radarr → Revenge) ✅
- [x] Trigger refresh in Radarr ✅
- [x] Get quality profiles ✅
- [x] Get root folders ✅

#### Sync Logic ✅ COMPLETE
- [x] Full sync (initial) (`SyncLibrary()`) ✅
- [x] Single movie sync (`SyncMovie()`) ✅
- [x] File path mapping ✅
- [x] Collection sync ✅

#### Webhook Handler ✅ COMPLETE
- [x] Webhook handler (`internal/integration/radarr/webhook_handler.go`) ✅
- [x] Handle: Grab, Download, Rename, Delete events ✅
- [x] `POST /api/v1/webhooks/radarr` (API endpoint) ✅

#### API Handlers ✅ COMPLETE
- [x] `GET /api/v1/admin/integrations/radarr/status` ✅
- [x] `POST /api/v1/admin/integrations/radarr/sync` ✅
- [x] `GET /api/v1/admin/integrations/radarr/quality-profiles` ✅
- [x] `GET /api/v1/admin/integrations/radarr/root-folders` ✅

#### River Jobs ✅ COMPLETE
- [x] RadarrSyncJob - Full library sync ✅
- [x] RadarrWebhookJob - Process webhook events ✅

#### Tests ✅ COMPLETE
- [x] Unit tests with mock API (`client_test.go`, `mapper_test.go`) ✅
- [x] Unit tests for API handlers (`handler_radarr_test.go` - 15 tests) ✅
- [ ] Integration tests (optional)

**Commits**:
- 6ad5379d83 - feat(radarr): implement Radarr integration client and sync service

### Frontend (Basic SvelteKit)

#### Project Setup
- [ ] SvelteKit 2 initialization
- [ ] Svelte 5 configuration
- [ ] TypeScript setup
- [ ] Tailwind CSS 4 setup
- [ ] shadcn-svelte components

#### Authentication Flow
- [ ] Login page (`/login`)
- [ ] Registration page (`/register`)
- [ ] Password reset flow
- [ ] JWT storage (httpOnly cookie)
- [ ] Auth store (Svelte store)
- [ ] Protected routes

#### Layout
- [ ] Navigation sidebar
- [ ] Header with user menu
- [ ] Responsive design
- [ ] Dark mode (default)

#### Library Browser
- [ ] Movies grid view (`/movies`)
- [ ] Movie card component
- [ ] Sorting (title, year, added)
- [ ] Filtering (genre, year)
- [ ] Pagination/infinite scroll
- [ ] Search integration

#### Movie Detail Page
- [ ] Hero backdrop
- [ ] Poster image
- [ ] Title, year, runtime
- [ ] Overview
- [ ] Cast carousel
- [ ] Crew list
- [ ] Similar movies
- [ ] Play button
- [ ] Watch progress

#### Search
- [ ] Global search bar
- [ ] Search results page
- [ ] Autocomplete dropdown

#### Basic Player
- [ ] Player page (`/play/[id]`)
- [ ] HLS.js integration
- [ ] Basic controls (play, pause, seek)
- [ ] Progress tracking
- [ ] Quality selection
- [ ] Subtitle selection

#### Settings
- [ ] Profile settings
- [ ] Playback preferences
- [ ] Language preference

#### Admin Pages
- [ ] Dashboard overview
- [ ] Library management
- [ ] User management
- [ ] Integration settings (Radarr)

#### Components (shadcn-svelte)
- [ ] Button, Input, Card
- [ ] Dialog, Sheet
- [ ] Select, Dropdown
- [ ] Avatar, Badge
- [ ] Skeleton loaders
- [ ] Toast notifications

#### API Client
- [ ] Type-safe API client
- [ ] Error handling
- [ ] Token refresh logic
- [ ] TanStack Query integration

### Infrastructure

#### Typesense Deployment
- [ ] Docker Compose service
- [ ] Helm chart subchart
- [ ] Environment variables

#### Full Docker Compose Stack
- [ ] revenge (backend)
- [ ] revenge-frontend
- [ ] postgresql
- [ ] dragonfly
- [ ] typesense
- [ ] traefik (reverse proxy)

#### Docker Images
- [ ] Backend multi-stage Dockerfile
- [ ] Frontend multi-stage Dockerfile
- [ ] Combined nginx config

### Documentation
- [ ] Getting started guide
- [ ] Installation guide (Docker)
- [ ] Configuration reference
- [ ] Radarr setup guide
- [ ] Complete OpenAPI spec
- [ ] Swagger UI endpoint
- [ ] API authentication guide

---

## MVP Verification Checklist

- [ ] Movies display in frontend
- [ ] Search works end-to-end
- [ ] Radarr sync imports movies
- [ ] Watch progress saves and restores
- [ ] Player plays video files
- [ ] Authentication works (login/logout)
- [ ] MFA works (TOTP + WebAuthn + backup codes)
- [ ] RBAC enforced on admin pages
- [ ] All tests pass (80%+ coverage)
- [ ] CI pipeline passes
- [ ] Docker Compose stack works

---

## Design Documentation References

All design work is **COMPLETE**. Reference these during implementation:

### MFA
- [MFA.md](../docs/dev/design/services/MFA.md) - Complete MFA implementation plan

### Movie Module
- [MOVIE_MODULE.md](../docs/dev/design/features/video/MOVIE_MODULE.md)
- [COLLECTIONS.md](../docs/dev/design/features/shared/COLLECTIONS.md)
- [LIBRARY_TYPES.md](../docs/dev/design/features/shared/LIBRARY_TYPES.md)

### Integrations
- [TMDB.md](../docs/dev/design/integrations/metadata/video/TMDB.md)
- [RADARR.md](../docs/dev/design/integrations/servarr/RADARR.md)
- [TYPESENSE.md](../docs/dev/design/integrations/infrastructure/TYPESENSE.md)

### Services
- [METADATA.md](../docs/dev/design/services/METADATA.md)
- [SEARCH.md](../docs/dev/design/services/SEARCH.md)
- [LIBRARY.md](../docs/dev/design/services/LIBRARY.md)
- [USER_SETTINGS.md](../docs/dev/design/services/USER_SETTINGS.md)

### Technical
- [FRONTEND.md](../docs/dev/design/technical/FRONTEND.md)
- [API.md](../docs/dev/design/technical/API.md)
- [HTTP_CLIENT.md](../docs/dev/design/patterns/HTTP_CLIENT.md)

---

## Current Progress Summary

### ✅ Completed (v0.2.0)
- PostgreSQL pool with metrics
- Dragonfly/Redis L2 cache
- Otter L1 cache
- River job queue
- Settings service (server-level)
- User service foundation
- Auth service foundation
- Session service foundation
- Password hashing (argon2id with bcrypt backward compat)

### ✅ Completed (v0.3.0 Sprint)
- **MFA Implementation** (Phases 1-5) - TOTP, WebAuthn, Backup Codes
- **Movie Module Backend** - Entity, Repository, Service, Library Provider, API, River Jobs
- **TMDb Metadata Service** - Client, Mapper, Caching
- **Typesense Search Integration** - Schema, Service, API, River Jobs
- **Radarr Integration** - Client, Sync Service, Webhook Handler, API Handlers, River Jobs

### 🔄 In Progress
- **Movie Module Tests** - Currently 46.7%, target 80%+

### 🔴 Not Started (v0.3.0 MVP)
- Frontend (SvelteKit)
- Full Docker Compose stack

---

## 🚨 Feature Gap Analysis (2026-02-04)

> **Analysis Complete**: Compared against Jellyfin, Plex, Overseerr, Tautulli, Navidrome, Audiobookshelf, Kavita, Immich

### Critical Gaps for v0.3.0 MVP

#### 1. Library Scanner - MediaInfo Integration ✅ COMPLETE
**Previous State**: `ExtractFileInfo()` in `library_scanner.go` was a stub - only got file size and container extension
**Solution**: Implemented go-astiav (FFmpeg bindings) instead of FFprobe CLI for better performance and type safety
**Completed**: 2026-02-04

**Implemented**:
- [x] go-astiav wrapper (`mediainfo.go`) with native FFmpeg bindings
- [x] Video: resolution, codec, profile, framerate, HDR detection, color info
- [x] Audio: all tracks with codec, channels, layout, language
- [x] Subtitles: all tracks with codec, language, forced flag
- [x] Updated `ExtractFileInfo()` to use go-astiav
- [x] Unit tests (8 passing)
- [x] Dockerfile updated for CGO/FFmpeg libs

**Files created/modified**:
- `internal/content/movie/mediainfo.go` (NEW - 540 lines)
- `internal/content/movie/mediainfo_test.go` (NEW - 165 lines)
- `internal/content/movie/library_scanner.go` (MODIFIED)
- `Dockerfile` (MODIFIED - FFmpeg dev libs)

#### 2. Real-time File Watching - Not Implemented ❌
**Current State**: fsnotify is in go.mod but NOT used anywhere in code
**Problem**: Library changes require manual scan trigger - no auto-detection of new/changed/deleted files
**Priority**: Nice-to-have for v0.3.0 (Radarr webhooks cover this), required for standalone mode

**Deferred to v0.4.0**:
- [ ] fsnotify watcher service
- [ ] Debounced file change events
- [ ] Trigger library scan on changes
- [ ] Config: `library.watch_enabled: true`

#### 3. Notification Service - Schema Only ❌
**Current State**: DB schema exists (`user_preferences.email_notifications`, etc.) but NO implementation
**Problem**: No way to notify users about anything (requests, new content, etc.)

**Required for v0.3.0** (minimal):
- [ ] Notification service interface
- [ ] Webhook notification agent (generic)
- [ ] Discord notification agent
- [ ] Config: `notifications.agents[]`

**Deferred**:
- [ ] Email (SMTP) agent
- [ ] Telegram agent
- [ ] Apprise integration
- [ ] Push notifications

### Important Gaps (v0.4.0+)

#### 4. Request System - Design Only ❌
**Current State**: Wiki page exists, NO implementation
**Roadmap**: Scheduled for v0.8.0 (Intelligence milestone)
**Recommendation**: Consider moving to v0.4.0 - high user value

**Features**:
- [ ] Request table and SQLC queries
- [ ] Request service (create, approve, deny, auto-approve rules)
- [ ] Integration with Radarr/Sonarr (add movie on approval)
- [ ] Request API endpoints
- [ ] Request notifications

#### 5. Transcoding Service - Not Started ❌
**Current State**: Design exists (`04_PLAYER_ARCHITECTURE.md`), config keys defined, NO implementation
**Roadmap**: v0.6.0 (Playback milestone)
**Required for**: Playback of unsupported formats, quality selection

**Files referenced in design**:
- `internal/playback/transcoder/` - Not created
- `internal/playback/hls/` - Not created
- FFmpeg/go-astiav integration - Not started

#### 6. Hardware Acceleration - Not Started ❌
**Current State**: Config key `playback.transcode.hw_accel` defined but no implementation
**Support needed**: VAAPI, NVENC, QSV, VideoToolbox
**Priority**: Required for production transcoding performance

### Nice-to-Have Gaps (v1.0)

#### 7. Analytics/Statistics - Minimal ❌
**Current State**: `movie_watched` table tracks progress, but NO analytics dashboard
**Missing**: Play count per movie, user watch time, concurrent streams, bandwidth

#### 8. Skip Intro Detection - Not Started ❌
**Roadmap**: v0.6.0
**Design**: Exists in player architecture

#### 9. Trickplay Thumbnails - Not Started ❌
**Roadmap**: v0.6.0
**Design**: Exists in player architecture

#### 10. SyncPlay (Watch Together) - Not Started ❌
**Roadmap**: v0.6.0
**Design**: Exists in player architecture

---

### Action Items for v0.3.0 (ALLES VOR FRONTEND)

> **Entscheidung**: Alle Items müssen vor Frontend-Start abgeschlossen sein.
> **Geschätzter Gesamtaufwand**: ~45-55 Stunden

#### Phase 1: MediaInfo mit go-astiav (4-6h) ✅ COMPLETE

**Warum go-astiav statt FFprobe CLI?**
- Bereits für Transcoding (v0.6.0) geplant → keine zusätzliche Dependency
- Native Go Bindings → keine Exec-Calls
- Typed API → sicherer als JSON-Parsing
- CGO erforderlich, aber das brauchen wir eh für HW-Acceleration

**Tasks**:
1. [x] `internal/content/movie/mediainfo.go` - go-astiav Wrapper ✅
   - [x] `ProbeFile(path string) (*MediaInfo, error)` - Hauptfunktion
   - [x] Duration, Bitrate, Container Format
   - [x] Video Stream: Codec, Resolution, Framerate, HDR Info
   - [x] Audio Streams: Codec, Channels, Language, Title
   - [x] Subtitle Streams: Codec, Language, Forced
   - [x] Color Space, Color Range, Color Primaries
2. [x] `internal/content/movie/mediainfo_test.go` - Unit Tests ✅ (8 Tests)
3. [x] `ExtractFileInfo()` in `library_scanner.go` updaten ✅
4. [x] go-astiav v0.40.0 zu go.mod hinzugefügt ✅
5. [x] Dockerfile updaten (FFmpeg libs für CGO) ✅

**Completed**: 2026-02-04
**Files Created/Modified**:
- `internal/content/movie/mediainfo.go` (NEW - 540 lines)
- `internal/content/movie/mediainfo_test.go` (NEW - 165 lines)
- `internal/content/movie/library_scanner.go` (MODIFIED - ExtractFileInfo)
- `Dockerfile` (MODIFIED - FFmpeg dev libs for build + runtime)

#### Phase 2: Notification Service (6-8h) ✅ COMPLETE

**Agents für v0.3.0**:
- [x] Webhook (generisch) - Kann alles anbinden ✅
- [x] Discord - Sehr populär bei Self-Hostern ✅
- [x] Gotify/ntfy - Self-Hosted Push Notifications ✅
- [x] Email (SMTP) - Klassisch ✅

**Tasks**:
1. [x] `internal/service/notification/` - Service Package ✅
   - [x] `notification.go` - Interface + Event Types (230 lines)
   - [x] `dispatcher.go` - Event Router + User Preferences (230 lines)
   - [x] `agents/webhook.go` - Generic Webhook Agent (235 lines)
   - [x] `agents/discord.go` - Discord Webhook Agent (310 lines)
   - [x] `agents/email.go` - SMTP Email Agent (340 lines)
   - [x] `agents/gotify.go` - Gotify/ntfy Push Agents (420 lines)
2. [x] River Job für async Notification Dispatch (`notification_job.go`) ✅
3. [ ] User Notification Preferences API (deferred to API phase)
4. [ ] Admin Notification Settings API (deferred to API phase)
5. [x] Tests für alle Agents ✅ (55+ tests)

**Completed**: 2026-02-04
**Files Created**:
- `internal/service/notification/notification.go` - Event types, Agent interface
- `internal/service/notification/dispatcher.go` - Event dispatcher with async support
- `internal/service/notification/notification_test.go` - 25+ unit tests
- `internal/service/notification/dispatcher_test.go` - 15+ unit tests
- `internal/service/notification/agents/webhook.go` - Generic webhook with retry
- `internal/service/notification/agents/discord.go` - Discord embeds + colors
- `internal/service/notification/agents/email.go` - SMTP with TLS/STARTTLS
- `internal/service/notification/agents/gotify.go` - Gotify + ntfy support
- `internal/service/notification/agents/agents_test.go` - Agent unit tests
- `internal/infra/jobs/notification_job.go` - River job worker
- `internal/infra/jobs/queues.go` (MODIFIED - added notifications queue)

**Event Types**:
```go
const (
    EventMovieAdded       EventType = "movie.added"
    EventMovieAvailable   EventType = "movie.available"
    EventRequestCreated   EventType = "request.created"
    EventRequestApproved  EventType = "request.approved"
    EventUserCreated      EventType = "user.created"
    EventPlaybackStarted  EventType = "playback.started"
    EventLibraryScanDone  EventType = "library.scan_done"
    EventLoginSuccess     EventType = "auth.login_success"
    EventLoginFailed      EventType = "auth.login_failed"
    EventMFAEnabled       EventType = "auth.mfa_enabled"
    EventPasswordChanged  EventType = "auth.password_changed"
)
```

#### Phase 3: Audit Logging (3-4h) ✅ COMPLETE

**Schema existiert** (`activity_log`), Activity Service bereits implementiert.

**Implementiert**:
- ✅ `internal/service/activity/logger.go` - Logger Interface
- ✅ Integration in Auth Service (Login/Logout/Password Change)
- ✅ Integration in User Service (Create/Update/Delete)
- ✅ Integration in Library Service (Create/Update/Delete)
- ✅ Integration in RBAC Service (AssignRole/RemoveRole)
- ✅ Alle Tests angepasst und passing

**Dateien erstellt/geändert**:
- `internal/service/activity/logger.go` (NEW)
- `internal/service/activity/logger_test.go` (NEW)
- `internal/service/activity/module.go` (MODIFIED - adds NewLogger)
- `internal/service/auth/service.go` (MODIFIED - activity logging)
- `internal/service/auth/module.go` (MODIFIED - injects activity.Logger)
- `internal/service/user/service.go` (MODIFIED - activity logging)
- `internal/service/user/module.go` (MODIFIED - injects activity.Logger)
- `internal/service/library/service.go` (MODIFIED - activity logging)
- `internal/service/library/module.go` (MODIFIED - injects activity.Logger)
- `internal/service/rbac/service.go` (MODIFIED - activity logging)
- `internal/service/rbac/module.go` (MODIFIED - injects activity.Logger)

**Completion Date**: 2026-02-04

#### Phase 4: RBAC Erweiterungen (6-8h) ✅ COMPLETE

**Alle 4 Rollen implementieren**:

| Rolle | Permissions |
|-------|-------------|
| `admin` | `*:*` - Alles |
| `moderator` | `users:read`, `users:update`, `requests:*`, `content:moderate`, `audit:read` |
| `user` | `movies:read`, `libraries:read`, `requests:create`, `self:*` |
| `guest` | `movies:read`, `libraries:read` (kein Write, keine History) |

**Implemented**:
- [x] Migration für Moderator + Guest Rollen (000027_add_moderator_role)
- [x] `CreateRole(name, permissions)` - Admin API
- [x] `DeleteRole(name)` - Admin API
- [x] `UpdateRolePermissions()` - Admin API
- [x] `ListRoles()` - Admin API
- [x] `ListPermissions()` - Verfügbare Permissions auflisten
- [x] Tests für alle neuen Endpoints (23 tests)
- [x] OpenAPI spec updated with new endpoints

**New Files**:
- `internal/service/rbac/roles.go` (~350 lines) - Role management service
- `migrations/000027_add_moderator_role.up.sql` - Moderator permissions
- `migrations/000027_add_moderator_role.down.sql` - Rollback

**New Endpoints**:
- `GET /api/v1/rbac/roles` - List all roles
- `POST /api/v1/rbac/roles` - Create custom role
- `GET /api/v1/rbac/roles/{roleName}` - Get role details
- `DELETE /api/v1/rbac/roles/{roleName}` - Delete custom role
- `PUT /api/v1/rbac/roles/{roleName}/permissions` - Update permissions
- `GET /api/v1/rbac/permissions` - List available permissions

**Completion Date**: 2026-02-04

#### Phase 5: Rate Limiter Migration (2-3h) ✅ COMPLETE

**Von sync.Map zu Dragonfly** für Multi-Instance Support.

**Implemented**:
1. [x] `internal/api/middleware/ratelimit_redis.go` - Redis Rate Limiter (~250 lines)
   - [x] Sliding Window mit Rueidis
   - [x] Atomare Operationen via Lua Script (slidingWindowScript)
   - [x] Per-IP Rate Limiting mit Endpoint-aware Keys
   - [x] Rate-Limit Headers (X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset)
2. [x] Config: `rate_limit.backend: "memory" | "redis"` in config.go
3. [x] Automatic Fallback zu Memory wenn Redis nicht verfügbar
4. [x] Tests: `ratelimit_redis_test.go` (18 Unit Tests)
5. [x] Server Integration: `server.go` wählt Backend basierend auf Config

**Files Created**:
- `internal/api/middleware/ratelimit_redis.go` (~250 lines)
- `internal/api/middleware/ratelimit_redis_test.go` (18 tests)

**Completion Date**: 2026-02-04

#### Phase 6: Cache für Hot Paths (4-6h) ✅ COMPLETE

**Gecachte Hot Paths**:

| Path | TTL | Invalidierung |
|------|-----|---------------|
| Session Validation | 30s | Bei Revoke/Logout |
| RBAC Enforce | 30s | Bei Policy Change |
| RBAC User Roles | 5min | Bei Role Change |
| Server Settings | 5min | Bei Setting Update |
| User Settings | 2min | Bei User Update |

**Implemented**:
1. [x] Cache Key Infrastructure (`internal/infra/cache/keys.go` ~200 lines)
   - [x] Key Prefixes: `session:`, `rbac:enforce:`, `rbac:roles:`, `settings:server:`, `settings:user:`
   - [x] TTL Constants: SessionTTL=30s, RBACPolicyTTL=5min, RBACEnforceTTL=30s
   - [x] Helper Functions: SessionKey(), RBACEnforceKey(), ServerSettingKey(), UserSettingKey()
   - [x] Invalidation Helpers: InvalidateSession(), InvalidateRBACForUser(), InvalidateAllRBAC()
2. [x] Cached Session Service (`internal/service/session/cached_service.go`)
   - [x] ValidateSession() mit Cache-Aside Pattern
   - [x] RevokeSession() mit Cache Invalidation
   - [x] RevokeAllUserSessions() mit Pattern-based Invalidation
3. [x] Cached RBAC Service (`internal/service/rbac/cached_service.go`)
   - [x] Enforce() - Cached Permission Checks
   - [x] GetUserRoles() - Cached Role Lookups
   - [x] Write-Through Invalidation bei AssignRole/RemoveRole/AddPolicy/RemovePolicy
4. [x] Cached Settings Service (`internal/service/settings/cached_service.go`)
   - [x] GetServerSetting()/ListServerSettings() - Cached
   - [x] GetUserSetting() - Per-User Caching
   - [x] Write-Through Invalidation bei Updates
5. [x] **Cached Movie Service** (`internal/content/movie/cached_service.go`)
   - [x] GetMovie(), GetCast(), GetCrew(), GetGenres(), GetFiles() - 5min TTL
   - [x] ListMovies(), ListRecentlyAdded(), ListTopRated() - 2-5min TTL
   - [x] GetContinueWatching() - Per-User 1min TTL
   - [x] Write-Through Invalidation bei Create/Update/Delete
6. [x] **Cached Search Service** (`internal/service/search/cached_service.go`)
   - [x] Search() - 30sec TTL
   - [x] Autocomplete() - 30sec TTL
   - [x] GetFacets() - 2min TTL
7. [x] **Cached Library Service** (`internal/service/library/cached_service.go`)
   - [x] Get(), List(), Count() - 10min TTL
   - [x] Write-Through Invalidation bei Create/Update/Delete/CompleteScan
8. [x] **Cached User Service** (`internal/service/user/cached_service.go`)
   - [x] GetUser(), GetUserByUsername() - 1min TTL
   - [x] Write-Through Invalidation bei Update/Delete
9. [x] **Image Cache Enhancement** (`internal/infra/image/service.go`)
   - [x] ETags für Conditional Requests (If-None-Match → 304)
   - [x] Immutable Cache-Control für CDN Caching
   - [x] CORS Headers für Frontend Clients
10. [x] Tests: keys_test.go (13 tests), cached_service_test.go (5 tests)

**Files Created/Modified**:
- `internal/infra/cache/keys.go` (~400 lines) - Extended with all key prefixes and TTLs
- `internal/infra/cache/keys_test.go` (13 tests)
- `internal/service/session/cached_service.go`
- `internal/service/session/cached_service_test.go` (5 tests)
- `internal/service/rbac/cached_service.go`
- `internal/service/settings/cached_service.go`
- `internal/content/movie/cached_service.go` (~450 lines)
- `internal/service/search/cached_service.go` (~160 lines)
- `internal/service/library/cached_service.go` (~210 lines) NEW
- `internal/service/user/cached_service.go` (~140 lines) NEW
- `internal/infra/image/service.go` (MODIFIED - ETags + cache headers)

**Cache TTL Summary**:
| Data Type | TTL | Invalidation |
|-----------|-----|--------------|
| Session | 30s | On revoke/logout |
| RBAC Enforce | 30s | On policy change |
| RBAC Roles | 5min | On role change |
| Server Settings | 5min | On setting update |
| User Settings | 2min | On user update |
| Movie | 5min | On update/delete |
| Movie Cast/Crew | 5min | On movie update |
| Recently Added | 2min | On movie create |
| Top Rated | 5min | On rating change |
| Continue Watching | 1min | On progress update |
| Search Results | 30sec | On index update |
| Autocomplete | 30sec | On index update |
| Facets | 2min | On index update |
| Library | 10min | On scan complete |
| User | 1min | On user update |
| Images | 7 days (disk) | Never (immutable) |

**Completion Date**: 2026-02-04

#### Phase 7: Observability (4-5h) ✅ DONE

**pprof (nur Dev Mode)** ✅:
```go
if config.Logging.Development {
    RegisterPprofHandlers(mux)
}
```

**Prometheus Metrics** ✅:
- HTTP Request Latency (Histogram) - `revenge_http_request_duration_seconds`
- HTTP Request Count (Counter) - `revenge_http_requests_total`
- HTTP Requests In Flight (Gauge) - `revenge_http_requests_in_flight`
- Active Sessions (Gauge) - `revenge_sessions_active_total`
- Cache Hit/Miss (Counter) - `revenge_cache_hits_total`, `revenge_cache_misses_total`
- Cache Operation Duration (Histogram) - `revenge_cache_operation_duration_seconds`
- Database Query Latency (Histogram) - `revenge_db_query_duration_seconds`
- River Job Queue Size (Gauge) - `revenge_jobs_queue_size`
- Job Duration (Histogram) - `revenge_jobs_duration_seconds`
- Auth Attempts (Counter) - `revenge_auth_attempts_total`
- Rate Limit Hits (Counter) - `revenge_ratelimit_hits_total`
- Library Scan Metrics - `revenge_library_*`
- Search Query Metrics - `revenge_search_*`

**Tasks**:
1. [x] pprof Endpoint (nur wenn `config.Logging.Development: true`)
2. [x] `/metrics` Endpoint mit Prometheus Registry (Port: main+1000)
3. [x] HTTP Middleware für Request Metrics (ogen-Middleware)
4. [x] Cache Metrics Integration (L1/L2 Hit/Miss, Operation Duration)
5. [x] Database Query Metrics (ready for integration)
6. [x] River/Job Metrics (ready for integration)
7. [x] Grafana Dashboard Template

**Files Created**:
- `internal/infra/observability/metrics.go` (~230 lines) - All Prometheus metrics
- `internal/infra/observability/middleware.go` (~190 lines) - HTTP metrics middleware
- `internal/infra/observability/pprof.go` (~25 lines) - pprof handler registration
- `internal/infra/observability/server.go` (~90 lines) - Observability HTTP server
- `internal/infra/observability/module.go` - fx module
- `internal/infra/observability/metrics_test.go` - Unit tests
- `deploy/grafana-dashboard.json` - Grafana dashboard template

**Modified Files**:
- `internal/app/module.go` - Added observability.Module
- `internal/api/server.go` - Added HTTPMetricsMiddleware to ogen
- `internal/infra/cache/cache.go` - Added cache metrics instrumentation

**Observability Endpoints** (Port: main+1000):
- `/metrics` - Prometheus scrape endpoint
- `/health/live` - Kubernetes liveness probe
- `/health/ready` - Kubernetes readiness probe
- `/debug/pprof/*` - pprof endpoints (development mode only)

**Completion Date**: 2026-02-04

#### Phase 8: Test Coverage 80%+ (8-12h) 🔄 IN PROGRESS

**Current Coverage**: 5.8% → Target: 80%

**Fixes Applied**:
- [x] Fix TestDefaultQueueConfig (jobs package - 4 queues)
- [x] Fix TestMigrationsUpDown (migration 000027 moderator role)

**Tests Added**:
- [x] `internal/infra/logging/logging_test.go` (NEW - 78.7% coverage)
  - TestNewLogger_DefaultConfig, _DevelopmentMode, _ProductionMode, _LogLevels
  - TestNewZapLogger_* (multiple configurations)
  - TestParseLevel (all log levels)
  - TestNewTestLogger, TestConfig_Defaults

- [x] `internal/infra/observability/metrics_test.go` (EXTENDED - 42.1% coverage)
  - TestStatusResponseWriter (3 subtests)
  - TestIsAlphanumericWithHyphens
  - TestHasDigits
  - TestStandardHTTPMetricsMiddleware (3 subtests)
  - TestExtractStatusFromResponse

- [x] `internal/infra/search/module_test.go` (EXTENDED - 56.1% coverage)
  - TestNewClient_Disabled, _EnabledEmptyURL, _EnabledWithURL
  - TestClient_DisabledOperations (10 subtests)
  - TestClient_IsEnabled (2 subtests)

**Coverage Progress**:
| Package | Before | After |
|---------|--------|-------|
| logging | 0% | 78.7% ✅ |
| observability | 28% | 42.1% |
| search | 0% | 56.1% |
| database | FAIL | 76.3% ✅ |
| health | FAIL | 45.2% ✅ |

**Remaining Tasks**:
1. [ ] Service layer tests (rbac, session, auth, settings)
2. [ ] Content module tests (movie, moviejobs)
3. [ ] Integration tests (radarr)
4. [ ] Verify 80% total coverage

**Focus Areas**:
- Movie Service + Library Scanner
- Notification Service
- Audit Service
- RBAC Service Extensions
- Cache Layer
- Rate Limiter

**Tasks**:
1. [ ] Movie Service Tests
2. [ ] Library Scanner Tests (mit go-astiav Mocks)
3. [ ] Notification Agent Tests
4. [ ] Audit Service Tests
5. [ ] RBAC Extension Tests
6. [ ] Cache Integration Tests
7. [ ] E2E Tests für neue APIs

---

### Zusammenfassung Pre-Frontend

| Phase | Aufwand | Status |
|-------|---------|--------|
| 1. MediaInfo (go-astiav) | 4-6h | ✅ DONE |
| 2. Notification Service | 6-8h | ✅ DONE |
| 3. Audit Logging | 3-4h | ✅ DONE |
| 4. RBAC Erweiterungen | 6-8h | ✅ DONE |
| 5. Rate Limiter Migration | 2-3h | ✅ DONE |
| 6. Cache Hot Paths | 4-6h | ✅ DONE |
| 7. Observability | 4-5h | ✅ DONE |
| 8. Test Coverage 80% | 8-12h | 🔄 IN PROGRESS |
| **Total** | **37-52h** | |

**Fortschritt**: Phase 1-7 ✅, Phase 8 🔄 (2026-02-04)

---

### Zusätzliche Erkenntnisse aus Architektur-Review

#### RBAC/Permissions Gaps ✅ RESOLVED (Phase 4)

**Aktueller Stand** (nach Phase 4):
- Casbin mit RBAC Model (sub, obj, act) ✅
- 4 Rollen implementiert: `admin`, `moderator`, `user`, `guest` ✅
- Custom Role CRUD via API ✅
- Per-Library Permissions → v0.4.0

**Service-Methoden vorhanden** ✅:
- `AssignRole()`, `RemoveRole()`, `GetUserRoles()`, `HasRole()`, `GetUsersForRole()`
- `CreateRole()`, `DeleteRole()`, `UpdateRolePermissions()`, `ListRoles()`, `ListPermissions()`

**TODO für RBAC**:
- [x] Moderator-Rolle in Migrations hinzufügen (Phase 4) ✅
- [x] API für Custom Role CRUD (Phase 4) ✅
- [ ] Per-Library Access Control (v0.4.0)

#### Cluster/Multi-Instance Gaps ✅ RESOLVED (Phase 5 + 6)

**Rate Limiter** ✅ FIXED:
- Redis/Dragonfly-basierter Rate Limiter mit Sliding Window
- Lua Script für atomare Operationen
- Fallback zu In-Memory wenn Redis nicht verfügbar

**Caching** ✅ FIXED:
- L1 (Otter) + L2 (Rueidis) Unified Cache
- Cache-Aside Pattern für Session/RBAC/Settings
- Pattern-based Invalidation

**Aktuell Cluster-Ready** ✅:
- PostgreSQL für State
- Dragonfly für Cache + Rate Limiting
- River für Job Queue
- Casbin mit PostgreSQL Adapter
- Rate Limiter mit Redis Backend
- Session/RBAC/Settings Caching mit L1+L2

**Benötigt Anpassung**:
- ~~Rate Limiter → Dragonfly~~ ✅ DONE (Phase 5)
- WebSocket für SyncPlay später → Redis PubSub (v0.6.0)

#### Profiling/Monitoring Gaps

**Benchmarks vorhanden** ✅:
- `internal/crypto/password_bench_test.go` - Password hashing benchmarks

**Wird implementiert** (Phase 7):
- `/debug/pprof/*` Endpoint (nur in Dev Mode)
- `/metrics` Endpoint (Prometheus)
- OpenTelemetry Traces (Dependency da, nicht konfiguriert)

#### Cache Best Practices Review

**Aktuell implementiert** ✅:
- L1: Otter (In-Memory, configurable)
- L2: Rueidis → Dragonfly
- TTL Handling mit L1/L2 Koordination
- TMDb Client hat eigenen Cache

**Wird gecached** (Phase 6):
- Movie Details API
- User Permissions/Roles
- Library Stats
- Search Results

**Cache Invalidation**:
- Event-basierte Invalidierung (Phase 6)
- Cache-Tags für Gruppeninvalidierung

---

## Execution Strategy

1. **MFA First** ✅ COMPLETE
   - Provides essential security feature
   - Tests auth/session integration
   - Proves encryption patterns

2. **Movie Module Backend** ✅ COMPLETE
   - Core business logic
   - Database schema
   - Library scanning
   - TMDb integration

3. **Search Integration** ✅ COMPLETE
   - Typesense setup
   - Movie indexing
   - Search API

4. **Radarr Integration** ✅ COMPLETE
   - Radarr client with rate limiting
   - Sync service
   - Webhook handler
   - Admin API endpoints
   - River jobs for async processing

5. **Frontend Development** (~40-60 hours)
   - SvelteKit setup
   - Authentication UI
   - Movie browser
   - Movie detail pages
   - Basic player
   - Admin panel

6. **Infrastructure & Testing** (~8-16 hours)
   - Docker Compose full stack
   - End-to-end tests
   - Documentation
   - Deployment guides

**Remaining Estimated Effort**: 56-88 hours (~1.5-2 weeks full-time)

---

## Notes

- All design documents are complete and ready
- Test-first approach with 80%+ coverage target
- Commit after each major milestone
- Keep TODO updated with progress
- Run tests and linter before each commit
