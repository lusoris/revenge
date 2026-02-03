# Comprehensive TODO - v0.3.0 MVP

**Last Updated**: 2026-02-03 23:55
**Current Focus**: Movie Module - Tests
**Status**: Backend Complete ✅ → TMDb Complete ✅ → Library Provider Complete ✅ → River Jobs Complete ✅ → Tests 🟡 (46.7%)

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
- [ ] Rate Limiting (per-endpoint)
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
- [ ] OpenAPI spec integration

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

#### TMDb Service
- [ ] Search movie
- [ ] Get movie details
- [ ] Get movie credits (cast/crew)
- [ ] Get movie images
- [ ] Get similar movies
- [ ] Get collection details

#### Image Handler
- [ ] Poster download/cache
- [ ] Backdrop download/cache
- [ ] Profile image download/cache
- [ ] Image proxy endpoint

#### API Handlers
- [ ] `GET /api/v1/metadata/search/movie?q=`
- [ ] `GET /api/v1/metadata/movie/:tmdbId`
- [ ] `GET /api/v1/images/:type/:path` (proxy)

#### Tests
- [ ] Unit tests with mock API
- [ ] Integration tests (optional, needs API key)

### Search Service (Typesense)

#### Typesense Setup
- [ ] Client configuration
- [ ] Collection schemas
- [ ] Index management

#### Movie Collection Schema
- [ ] Define schema (title, original_title, overview, year, genres, cast, director, rating, added_at)

#### Search Service
- [ ] Index movie
- [ ] Remove from index
- [ ] Search movies (full-text)
- [ ] Faceted search (genre, year)
- [ ] Autocomplete

#### API Handlers
- [ ] `GET /api/v1/search?q=&type=movie`
- [ ] `GET /api/v1/search/autocomplete?q=`

#### River Jobs
- [ ] SearchIndexJob - Index single item
- [ ] SearchReindexJob - Full reindex

#### Tests
- [ ] Unit tests
- [ ] Integration tests with Typesense container

### Radarr Integration

#### Radarr Client
- [ ] API v3 implementation
- [ ] Authentication (API key)
- [ ] Error handling

#### Radarr Service
- [ ] Get all movies
- [ ] Get movie by ID
- [ ] Get movie files
- [ ] Sync library (Radarr → Revenge)
- [ ] Trigger refresh in Radarr
- [ ] Get quality profiles
- [ ] Get root folders

#### Sync Logic
- [ ] Full sync (initial)
- [ ] Incremental sync (changes only)
- [ ] File path mapping
- [ ] Conflict resolution

#### Webhook Handler
- [ ] `POST /api/v1/webhooks/radarr`
- [ ] Handle: Grab, Download, Rename, Delete events

#### API Handlers
- [ ] `GET /api/v1/admin/integrations/radarr/status`
- [ ] `POST /api/v1/admin/integrations/radarr/sync`
- [ ] `GET /api/v1/admin/integrations/radarr/quality-profiles`

#### River Jobs
- [ ] RadarrSyncJob - Full library sync
- [ ] RadarrWebhookJob - Process webhook events

#### Tests
- [ ] Unit tests with mock API
- [ ] Integration tests (optional)

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

### 🔄 In Progress
- **MFA Implementation** (10-20 hours) - Current Sprint
  - Phase 1: Foundation (migrations + encryption)
  - Phase 2: TOTP
  - Phase 3: WebAuthn
  - Phase 4: Integration
  - Phase 5: Hardening

### 🔴 Not Started (v0.3.0 MVP)
- Movie Module (backend)
- TMDb metadata service
- Radarr integration
- Typesense search
- Frontend (SvelteKit)
- Full Docker Compose stack

---

## Execution Strategy

1. **MFA First** (Current Sprint, 10-20 hours)
   - Provides essential security feature
   - Tests auth/session integration
   - Proves encryption patterns

2. **Movie Module Backend** (~20-30 hours)
   - Core business logic
   - Database schema
   - Library scanning
   - TMDb integration
   - Radarr integration

3. **Search Integration** (~8-12 hours)
   - Typesense setup
   - Movie indexing
   - Search API

4. **Frontend Development** (~40-60 hours)
   - SvelteKit setup
   - Authentication UI
   - Movie browser
   - Movie detail pages
   - Basic player
   - Admin panel

5. **Infrastructure & Testing** (~8-16 hours)
   - Docker Compose full stack
   - End-to-end tests
   - Documentation
   - Deployment guides

**Total Estimated Effort**: 86-138 hours (~2-3 weeks full-time)

---

## Notes

- All design documents are complete and ready
- Test-first approach with 80%+ coverage target
- Commit after each major milestone
- Keep TODO updated with progress
- Run tests and linter before each commit
