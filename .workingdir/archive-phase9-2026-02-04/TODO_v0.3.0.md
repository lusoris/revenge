# TODO v0.3.0 - Consolidated

**Last Updated**: 2026-02-04
**Status**: Pre-Frontend Phase
**Next**: Critical Fixes → Implementation Tasks → Frontend → MVP Release

---

## Quick Reference - Decisions Made (2026-02-04)

| Topic | Decision |
|-------|----------|
| Permissions | Fine-grained (`movie:list`, `movie:get`, etc.) |
| Cache Tests | testcontainers for Redis L2 paths |
| Session Storage | Hybrid (Dragonfly L1 + PostgreSQL L2) |
| Queue Priorities | 5 levels: critical, high, default, low, bulk |
| Log Retention | 90 days default, configurable |
| Pagination | Both cursor + offset, cursor default |
| Device Fingerprinting | Deferred to v0.6.0 (Transcoding) |

---

## Phase A0: Critical Fixes (Blockers from Stubs Analysis)

> **Source**: [STUBS_AND_UNIMPLEMENTED_REPORT.md](STUBS_AND_UNIMPLEMENTED_REPORT.md)

### A0.1: Auth Context User ID Extraction [P0-BLOCKER]
**Priority**: CRITICAL | **Effort**: 2-3h

All authenticated API handlers use placeholder `uuid.New()` or hardcoded UUIDs instead of actual user.

**Affected Files** (`internal/api/handler.go`):
- Lines 190-263: Multiple `// TODO: Get user ID from auth context`
- Lines 191, 211, 231, 247, 263: `userID := uuid.New() // Placeholder`
- Lines 413, 474, 499, 560: `userID := uuid.MustParse("550e8400-...")` hardcoded

**Tasks**:
- [ ] Create auth middleware to extract user ID from JWT claims
- [ ] Store user ID in request context (`ctx.Value("userID")`)
- [ ] Create helper: `func GetUserIDFromContext(ctx context.Context) (uuid.UUID, error)`
- [ ] Update all 9+ handler locations to use context extraction
- [ ] Add tests for middleware

---

### A0.2: Email Service Implementation [P0-BLOCKER]
**Priority**: CRITICAL | **Effort**: 4-6h

User registration, password reset, and email verification are non-functional.

**Affected Files**:
| File | Line | Issue |
|------|------|-------|
| `internal/service/auth/service.go` | 96 | `// TODO: Send verification email` |
| `internal/service/auth/service.go` | 156 | `// TODO: Send verification email` |
| `internal/service/auth/service.go` | 409 | `// TODO: Send reset email` |

**Tasks**:
- [ ] Create `internal/service/email/service.go`
- [ ] Implement SMTP transport (configurable)
- [ ] Implement SendGrid transport (optional, configurable)
- [ ] Email templates: verification, password reset, welcome
- [ ] Config: `email.provider: smtp|sendgrid`
- [ ] Config: `email.smtp.host`, `smtp.port`, `smtp.user`, `smtp.password`
- [ ] Config: `email.from_address`, `email.from_name`
- [ ] River job for async email sending
- [ ] Tests with mock SMTP server
- [ ] Update auth service to call email service

---

### A0.3: Session Count Implementation [P0-BLOCKER]
**Priority**: HIGH | **Effort**: 1h

`CountUserSessions` returns hardcoded 0.

**Affected File**: `internal/service/session/service.go:251`
```go
return 0, nil // TODO: Return actual count
```

**Tasks**:
- [ ] Add sqlc query: `CountSessionsByUserID`
- [ ] Update service to call repository
- [ ] Test

---

### A0.4: Avatar Upload Implementation [P1]
**Priority**: HIGH | **Effort**: 3-4h

Avatar upload returns `BadRequest` with "not yet implemented".

**Affected Files**:
- `internal/api/handler.go:562-566` - Handler stub
- `internal/service/user/service.go:326-327` - Service returns placeholder path

**Tasks**:
- [ ] Parse multipart form in handler
- [ ] Validate file type (JPEG, PNG, WebP)
- [ ] Validate file size (max 2MB configurable)
- [ ] Resize image to standard sizes (64x64, 128x128, 256x256)
- [ ] Storage interface: local filesystem initially
- [ ] Config: `avatar.storage: local|s3`
- [ ] Config: `avatar.max_size: 2MB`
- [ ] Config: `avatar.local_path: /data/avatars`
- [ ] Update user record with avatar URL
- [ ] Tests

---

### A0.5: Request Metadata Extraction [P1]
**Priority**: HIGH | **Effort**: 2h

IP address, user agent, fingerprint not extracted from requests.

**Affected Files**:
- `internal/api/handler.go:621` - Login handler
- `internal/api/handler.go:735` - Another handler

**Tasks**:
- [ ] Create middleware to extract:
  - IP address (with X-Forwarded-For support)
  - User-Agent header
  - Accept-Language header
- [ ] Store in request context
- [ ] Helper: `GetRequestMetadata(ctx) RequestMeta`
- [ ] Use in session creation, activity logging
- [ ] Tests

---

### A0.6: WebAuthn Session Cache [P1]
**Priority**: MEDIUM | **Effort**: 2h

WebAuthn challenge sessions not stored in cache.

**Affected Files**:
- `internal/service/mfa/webauthn.go:174` - `// TODO: Store session in cache`
- `internal/service/mfa/webauthn.go:340` - `// TODO: Store session in cache`
- `internal/service/auth/mfa_integration.go:112` - Returns error "not yet implemented"

**Tasks**:
- [ ] Store WebAuthn challenge in Dragonfly with 5min TTL
- [ ] Key format: `webauthn:session:{userID}:{sessionID}`
- [ ] Retrieve challenge during verification
- [ ] Delete after successful verification
- [ ] Tests

---

### A0.7: OIDC New User Creation [P1]
**Priority**: MEDIUM | **Effort**: 2h

First OIDC login doesn't auto-create user account.

**Affected File**: `internal/api/handler_oidc.go:121`
```go
// TODO: If IsNewUser, create the user account via user service
```

**Tasks**:
- [ ] Check if user exists by OIDC subject
- [ ] If not, create user with:
  - Email from OIDC claims
  - Display name from claims (or email prefix)
  - Default role (user)
  - Linked OIDC identity
- [ ] Create session for new user
- [ ] Tests

---

### A0.8: MFA Remember Device Setting [P2]
**Priority**: LOW | **Effort**: 1h

`RememberDeviceEnabled` hardcoded to false.

**Affected File**: `internal/service/mfa/manager.go:77`

**Tasks**:
- [ ] Query `user_mfa_settings` table for setting
- [ ] Return actual value in MFA status
- [ ] Test

---

## Phase A1: Movie Repository Completion [P2]

> 25+ repository methods returning "not implemented"

**Affected File**: `internal/content/movie/repository_postgres.go:136-390`

### A1.1: Core Movie Queries
- [ ] `ListMoviesByIDs(ctx, ids []uuid.UUID) ([]*Movie, error)`
- [ ] `CountMovies(ctx) (int64, error)`
- [ ] `SearchMovies(ctx, query string, limit int) ([]*Movie, error)`
- [ ] `ListMoviesByGenre(ctx, genreID, limit, offset) ([]*Movie, error)`
- [ ] `ListRecentMovies(ctx, limit) ([]*Movie, error)`
- [ ] `ListPopularMovies(ctx, limit) ([]*Movie, error)`
- [ ] `GetMovieWithDetails(ctx, id) (*MovieWithDetails, error)`
- [ ] `GetMovieByExternalID(ctx, provider, externalID) (*Movie, error)`

### A1.2: Movie Metadata Operations
- [ ] `UpdateMovieMetadata(ctx, id, metadata) error`
- [ ] `DeleteMovie(ctx, id) error`

### A1.3: Movie Credits
- [ ] `ListMovieCredits(ctx, movieID) ([]*Credit, error)`
- [ ] `GetMovieCredit(ctx, creditID) (*Credit, error)`
- [ ] `CreateMovieCredit(ctx, credit) (*Credit, error)`
- [ ] `DeleteMovieCredit(ctx, creditID) error`
- [ ] `DeleteMovieCredits(ctx, movieID) error`

### A1.4: Movie Collections
- [ ] `GetMovieCollection(ctx, movieID) (*Collection, error)`
- [ ] `ListCollectionMovies(ctx, collectionID) ([]*Movie, error)`

### A1.5: Movie Files
- [ ] `ListMovieFiles(ctx, movieID) ([]*File, error)`
- [ ] `GetMovieFile(ctx, fileID) (*File, error)`
- [ ] `CreateMovieFile(ctx, file) (*File, error)`
- [ ] `DeleteMovieFile(ctx, fileID) error`

**Note**: Many may already have sqlc queries - verify before implementing.

---

## Phase A2: Movie Jobs Completion [P2]

### A2.1: File Match Job
**Affected File**: `internal/content/movie/moviejobs/file_match.go:58-64`

**Tasks**:
- [ ] Implement `library.Service.MatchFile` method
- [ ] Update file match worker to use it
- [ ] Match logic: filename parsing → TMDb lookup → confidence scoring
- [ ] Tests

### A2.2: Metadata Refresh Job
**Affected Files**:
- `internal/content/movie/service.go:275-277`
- `internal/content/movie/moviejobs/metadata_refresh.go:90-92`

**Tasks**:
- [ ] Implement `RefreshMetadata(ctx, movieID)` in service
- [ ] Queue River job for async processing
- [ ] Refresh credits during metadata refresh
- [ ] Tests

---

## Phase A3: Test Infrastructure [P2]

### A3.1: Dragonfly Testcontainer
**Affected File**: `internal/testutil/containers.go:168-171`

**Tasks**:
- [ ] Implement `NewDragonflyContainer(t) (*DragonflyContainer, error)`
- [ ] Use `docker.io/docker/dragonfly:latest` image
- [ ] Return connection string
- [ ] Cleanup on test completion
- [ ] Test the container helper itself

### A3.2: Typesense Testcontainer
**Affected File**: `internal/testutil/containers.go:182-185`

**Tasks**:
- [ ] Implement `NewTypesenseContainer(t) (*TypesenseContainer, error)`
- [ ] Use `typesense/typesense:latest` image
- [ ] Configure API key
- [ ] Return connection details
- [ ] Cleanup on test completion
- [ ] Test the container helper itself

---

## Phase A4: Webhook & Notification Enhancements [P2]

### A4.1: Custom Webhook Templates
**Affected File**: `internal/service/notification/agents/webhook.go:190`

**Tasks**:
- [ ] Parse `PayloadTemplate` as Go template
- [ ] Provide template data: event type, payload, timestamp, etc.
- [ ] Validate template on agent creation
- [ ] Tests

---

## Phase A5: Library Matcher [P2]

### A5.1: Library File Matching
**Affected Files**:
- `internal/content/movie/library_matcher.go:120`
- `internal/content/movie/library_service.go:190`

**Tasks**:
- [ ] Implement file → movie matching algorithm
- [ ] Parse filename for title, year, quality
- [ ] Query TMDb for matches
- [ ] Score and rank matches
- [ ] Return best match or candidates
- [ ] Tests

---

## Phase A6: Permissions & Session [Existing Tasks]

### A6.1: Fine-Grained Permissions
**Priority**: HIGH | **Effort**: 4-6h

Current Casbin policies use coarse permissions. Need to implement fine-grained.

**Tasks**:
- [ ] Define permission taxonomy in `docs/dev/design/services/RBAC_PERMISSIONS.md`
  ```
  movies:list, movies:get, movies:create, movies:update, movies:delete
  libraries:list, libraries:get, libraries:create, libraries:scan
  users:list, users:get, users:create, users:update, users:delete
  settings:read, settings:write
  admin:* (wildcard for admin)
  ```
- [ ] Migration `000029_fine_grained_permissions.up.sql` - Update default role policies
- [ ] Update `internal/service/rbac/permissions.go` - Permission constants
- [ ] Update middleware to check granular permissions
- [ ] Update OpenAPI spec with permission requirements per endpoint
- [ ] Tests for permission checks

---

### A6.2: Hybrid Session Storage
**Priority**: HIGH | **Effort**: 4-6h

Sessions currently PostgreSQL-only. Implement Dragonfly L1 for performance.

**Tasks**:
- [ ] Update `internal/service/session/cached_service.go`:
  - [ ] ValidateSession: Check Dragonfly first, fallback to PostgreSQL
  - [ ] CreateSession: Write to both (write-through)
  - [ ] RevokeSession: Invalidate Dragonfly, update PostgreSQL
- [ ] Add session serialization (JSON or msgpack)
- [ ] Config: `session.cache_enabled: true` (default)
- [ ] Config: `session.cache_ttl: 5m`
- [ ] Benchmark: Compare latency with/without cache
- [ ] Tests with testcontainers Redis

---

### A6.3: 5-Level River Queue Priorities
**Priority**: MEDIUM | **Effort**: 2-3h

Currently using 3 priorities. Expand to 5 for better job management.

**Tasks**:
- [ ] Update `internal/infra/jobs/queues.go`:
  ```go
  const (
      QueueCritical = "critical"  // Security events, auth failures
      QueueHigh     = "high"      // User actions, notifications
      QueueDefault  = "default"   // Metadata fetching, sync
      QueueLow      = "low"       // Cleanup, maintenance
      QueueBulk     = "bulk"      // Library scans, batch ops
  )
  ```
- [ ] Assign existing jobs to appropriate queues:
  - `critical`: Auth audit events
  - `high`: Notification dispatch, webhook processing
  - `default`: Metadata refresh, Radarr sync
  - `low`: Session cleanup, expired token cleanup
  - `bulk`: Library scan, search reindex
- [ ] Update River config with priority weights
- [ ] Tests for queue assignment

---

### A6.4: Configurable Activity Log Retention
**Priority**: LOW | **Effort**: 1-2h

Activity logs should be configurable, default 90 days.

**Tasks**:
- [ ] Add config: `activity.retention_days: 90`
- [ ] Update cleanup job to use config value
- [ ] Add admin API: `PUT /api/v1/admin/settings/activity-retention`
- [ ] Migration for default server setting
- [ ] Tests

---

### A6.5: Cache Integration Tests (testcontainers)
**Priority**: MEDIUM | **Effort**: 3-4h

Cache package needs testcontainers-based Redis tests for L2 paths.

**Tasks**:
- [ ] Create `internal/infra/cache/integration_test.go`
- [ ] Use `internal/testutil/containers.go` for Redis/Dragonfly
- [ ] Test L2 cache operations:
  - [ ] Get/Set/Delete with real Redis
  - [ ] TTL expiration
  - [ ] Pattern-based invalidation
  - [ ] L1→L2 fallback
- [ ] Test distributed locks with real Redis
- [ ] Run only with `-tags=integration`

---

### A6.6: Test Coverage to 80% [ONGOING]
**Priority**: HIGH | **Effort**: 8-12h

Current coverage needs improvement. Focus on critical paths.

**Services needing tests**:
| Package | Current | Target |
|---------|---------|--------|
| session | 59.6% | 80% |
| auth | 29.9% | 80% |
| search | 37.0% | 80% |
| mfa | 12.7% | 80% |
| rbac | 1.3% | 80% |
| oidc | 1.7% | 80% |
| activity | 1.2% | 80% |
| user | 0% | 80% |
| settings | 0% | 80% |
| apikeys | 0% | 80% |
| library | 0% | 80% |

**Tasks**:
- [ ] Session service tests (mock repository)
- [ ] Auth service tests (mock dependencies)
- [ ] RBAC service tests (mock Casbin)
- [ ] User service tests
- [ ] Settings service tests
- [ ] Activity service tests
- [ ] API Keys service tests
- [ ] Library service tests
- [ ] MFA service tests (expand existing)
- [ ] OIDC service tests (expand existing)
- [ ] Search service tests (expand existing)

---

## Phase B: Frontend (SvelteKit)

**Effort**: 40-60h | **Start after**: Phase A complete

### B1: Project Setup
- [ ] SvelteKit 2 initialization
- [ ] Svelte 5 configuration
- [ ] TypeScript setup
- [ ] Tailwind CSS 4 setup
- [ ] shadcn-svelte components
- [ ] API client generation from OpenAPI

### B2: Authentication Flow
- [ ] Login page (`/login`)
- [ ] Registration page (`/register`)
- [ ] Password reset flow
- [ ] MFA setup page
- [ ] JWT storage (httpOnly cookie)
- [ ] Auth store (Svelte store)
- [ ] Protected routes

### B3: Layout & Navigation
- [ ] Navigation sidebar
- [ ] Header with user menu
- [ ] Responsive design
- [ ] Dark mode (default)

### B4: Library Browser
- [ ] Movies grid view (`/movies`)
- [ ] Movie card component
- [ ] Sorting (title, year, added)
- [ ] Filtering (genre, year)
- [ ] Pagination/infinite scroll
- [ ] Search integration

### B5: Movie Detail Page
- [ ] Hero backdrop
- [ ] Poster image
- [ ] Title, year, runtime, overview
- [ ] Cast carousel
- [ ] Crew list
- [ ] Similar movies
- [ ] Play button
- [ ] Watch progress

### B6: Search
- [ ] Global search bar
- [ ] Search results page
- [ ] Autocomplete dropdown
- [ ] Faceted filtering

### B7: Basic Player
- [ ] Player page (`/play/[id]`)
- [ ] HLS.js integration
- [ ] Basic controls (play, pause, seek)
- [ ] Progress tracking
- [ ] Quality selection (stub)
- [ ] Subtitle selection (stub)

### B8: Settings Pages
- [ ] Profile settings
- [ ] MFA management
- [ ] Playback preferences
- [ ] Language preference

### B9: Admin Pages
- [ ] Dashboard overview
- [ ] Library management
- [ ] User management
- [ ] Integration settings (Radarr)
- [ ] Activity logs viewer

---

## Phase C: Infrastructure & Release

**Effort**: 8-16h | **Start after**: Phase B complete

### C1: Docker Compose Stack
- [ ] `docker-compose.yml` with all services:
  - revenge (backend)
  - revenge-frontend
  - postgresql
  - dragonfly
  - typesense
  - traefik (reverse proxy)
- [ ] `.env.example` with all config
- [ ] Health check integration
- [ ] Volume mounts for persistence

### C2: Docker Images
- [ ] Backend multi-stage Dockerfile (verified)
- [ ] Frontend multi-stage Dockerfile
- [ ] Combined nginx config
- [ ] GitHub Actions for image builds

### C3: Documentation
- [ ] Getting started guide
- [ ] Installation guide (Docker)
- [ ] Configuration reference
- [ ] Radarr setup guide
- [ ] API authentication guide

### C4: MVP Verification
- [ ] Movies display in frontend
- [ ] Search works end-to-end
- [ ] Radarr sync imports movies
- [ ] Watch progress saves and restores
- [ ] Authentication works (login/logout)
- [ ] MFA works (TOTP + backup codes)
- [ ] RBAC enforced on admin pages
- [ ] All tests pass (80%+ coverage)
- [ ] CI pipeline passes
- [ ] Docker Compose stack works

---

## Completed (Reference)

### v0.2.0 Services ✅
- PostgreSQL pool with metrics
- Dragonfly/Redis L2 cache
- Otter L1 cache
- River job queue
- Settings service
- User service
- Auth service (JWT, Argon2id)
- Session service
- RBAC service (Casbin)
- API Keys service
- OIDC service
- Activity service
- Library service
- Health service
- MFA (TOTP, WebAuthn, Backup Codes)

### v0.3.0 Backend ✅
- Movie module (entity, repository, service, API)
- TMDb metadata service
- Library scanner with go-astiav MediaInfo
- Typesense search integration
- Radarr integration (client, sync, webhooks)
- Notification service (Webhook, Discord, Email, Gotify)
- Rate limiting (Redis-based)
- Cache hot paths
- Observability (Prometheus, pprof)
- RBAC extensions (4 roles, custom role API)

---

## Timeline Estimate

| Phase | Effort | Dependencies |
|-------|--------|--------------|
| A0: Critical Fixes | 15-20h | None |
| A1-A5: Stubs Completion | 20-30h | A0 |
| A6: Existing Tasks | 20-30h | A0 |
| B: Frontend | 40-60h | A0-A6 |
| C: Infrastructure | 8-16h | Phase B |
| **Total** | **103-156h** | |

**Estimated completion**: ~3-4 weeks full-time

### Priority Order
1. **A0.1-A0.3** (Critical blockers) - Must fix before any testing
2. **A0.4-A0.7** (P1 issues) - Required for functional MVP
3. **A6.1-A6.6** (Existing tasks) - Can parallel with A1-A5
4. **A1-A5** (P2 stubs) - Can defer some to post-MVP

---

## Notes

- Run `golangci-lint run ./...` before commits
- Run `go test ./... -short` frequently
- Update this TODO as tasks complete
- All design docs in `docs/dev/design/`
