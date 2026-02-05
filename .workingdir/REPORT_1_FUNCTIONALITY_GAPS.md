# Revenge Codebase Analysis: Functionality Gaps and Incomplete Implementations

**Generated**: 2026-02-05
**Analysis Type**: Comprehensive Codebase Review
**Overall Progress**: ~30% complete

---

## Executive Summary

The Revenge media server codebase is in **pre-MVP phase (v0.3.0)** with a **complete backend foundation** but **significant gaps in content modules, frontend, and integrations**. The project has strong infrastructure (auth, database, caching, job queue) but is missing ~80% of planned content types and the entire frontend application.

**Overall Progress**: ~30% complete
- **Infrastructure & Core Services**: 85% complete
- **Content Modules**: 15% complete (1 of 11 modules)
- **Integrations**: 10% complete (1 of 9 Arr services)
- **Frontend**: 0% (does not exist)

---

## 1. Content Modules Analysis

### 1.1 Fully Implemented Modules

| Module | Status | Implementation | Files | Tests | Coverage |
|--------|--------|----------------|-------|-------|----------|
| **Movie** | ✅ Complete | Full CRUD, TMDb metadata, library scanner, search indexing, Radarr sync | 50 files | Yes | 41.0% |

**Movie Module Details**:
- **Location**: `internal\content\movie\`
- **Database**: 6 tables (movies, movie_files, movie_credits, movie_collections, movie_genres, movie_watched)
- **Repository**: Full PostgreSQL implementation with all CRUD operations
- **Service**: Business logic complete with caching
- **Jobs**: Library scan, metadata refresh, search indexing (River jobs)
- **API Handlers**: Complete CRUD endpoints
- **Integrations**: TMDb client + Radarr sync service
- **Tests**: Unit tests, integration tests, mock repositories
- **Gaps**: None identified - this is the reference implementation

### 1.2 Stub/Placeholder Modules

| Module | Status | Evidence | Missing Components |
|--------|--------|----------|-------------------|
| **TV Show** | 🔴 Stub Only | Only placeholder.sql query | Repository, service, handlers, migrations, API endpoints |
| **QAR (Adult)** | 🔴 Stub Only | Only placeholder.sql query | Repository, service, handlers, migrations, API endpoints |
| **Music** | 🔴 Not Started | No files found | Everything |
| **Audiobook** | 🔴 Not Started | No files found | Everything |
| **Book** | 🔴 Not Started | No files found | Everything |
| **Podcast** | 🔴 Not Started | No files found | Everything |
| **Photo** | 🔴 Not Started | No files found | Everything |
| **Comics** | 🔴 Not Started | No files found | Everything |
| **LiveTV** | 🔴 Not Started | No files found | Everything |

**Stub Implementation Evidence**:

`internal\content\tvshow\db\placeholder.sql.go`:
```go
// Placeholder query for TV show content module (v0.3.0+)
// This minimal query is here to satisfy sqlc's requirement for non-empty query directories
// Real TV show queries will be implemented in v0.3.0
func (q *Queries) TVShowPlaceholder(ctx context.Context) (int32, error) {
    row := q.db.QueryRow(ctx, tVShowPlaceholder)
    var placeholder int32
    err := row.Scan(&placeholder)
    return placeholder, err
}
```

### 1.3 Missing CRUD Operations

**All stub modules** are missing:
1. Entity definitions (types.go)
2. Repository interface + PostgreSQL implementation
3. Service layer with business logic
4. Cached service wrappers
5. API handlers (ogen-generated interfaces are defined but not implemented)
6. Database migrations (tables don't exist)
7. River jobs for background processing
8. Library scanners/matchers
9. Metadata clients
10. Search schema definitions

---

## 2. Core Services Analysis

### 2.1 Fully Complete Services (✅)

| Service | Location | Status | Coverage | Notes |
|---------|----------|--------|----------|-------|
| **User** | `internal\service\user\` | ✅ Complete | 45.9% | Full user management, avatar handling |
| **Session** | `internal\service\session\` | ✅ Complete | 63.2% | Hybrid Dragonfly+PostgreSQL storage |
| **API Keys** | `internal\service\apikeys\` | ✅ Complete | 73.9% | Full API key lifecycle |
| **RBAC** | `internal\service\rbac\` | ✅ Complete | 39.6% | Casbin integration, 4 default roles |
| **Activity** | `internal\service\activity\` | ✅ Complete | 31.6% | Activity logging with cleanup jobs |
| **Settings** | `internal\service\settings\` | ✅ Complete | 41.0% | Server + user settings |
| **Library** | `internal\service\library\` | ✅ Complete | 43.6% | Library management |
| **OIDC** | `internal\service\oidc\` | ✅ Complete | 27.6% | OAuth2/OIDC provider integration |
| **Fingerprint** | `internal\service\fingerprint\` | ✅ Complete | N/A | Device fingerprinting |
| **Grants** | `internal\service\grants\` | ✅ Complete | N/A | Permission grants |

### 2.2 Partial Implementations (🟡)

| Service | Status | What's Complete | What's Missing | Coverage |
|---------|--------|-----------------|----------------|----------|
| **Auth** | 🟡 Partial | Login, JWT, password hashing, basic MFA flow | IP/User-Agent extraction, device fingerprinting integration | 38.4% |
| **MFA** | 🟡 Partial | TOTP + backup codes + WebAuthn infrastructure | Higher test coverage, more WebAuthn tests | 10.4% |
| **Email** | 🟡 Partial | SMTP sending complete | SendGrid API not implemented | 59.0% |
| **Search** | 🟡 Partial | Movie search complete, Typesense client | TV, music, multi-content search | 37.0% |
| **Notification** | 🟡 Partial | Dispatcher + 4 agents (Email, Discord, Gotify, Webhook) | More agents (Slack, Telegram, etc.), user preferences UI | 97.6% |

**Auth Service TODOs**:
- `internal\api\handler.go:709`: "TODO: extract IP, user agent, fingerprint from request"
- `internal\api\handler.go:823`: "TODO: Extract IP address and user agent from request"

**Email Service TODO**:
- `internal\service\email\service.go:205`: "TODO: Implement SendGrid API call"

### 2.3 Not Implemented Services (🔴)

| Service | Status | Design Doc | Notes |
|---------|--------|------------|-------|
| **Scrobbling** | 🔴 Planned | Yes | No files exist (Trakt, Last.fm, ListenBrainz, Letterboxd, Simkl) |
| **Analytics** | 🔴 Planned | Yes | No files exist |
| **Playback** | 🔴 Planned | Partial | No files exist (streaming, transcoding, progress tracking) |

---

## 3. Infrastructure Services

### 3.1 Database (PostgreSQL) - ✅ Complete

- **Status**: Fully implemented
- **Coverage**: 34.4%
- **Files**: 29 migrations (58 files: up/down)
- **Features**:
  - pgxpool with self-healing
  - Connection health checks
  - sqlc code generation
  - Transactional support
- **Migrations**: Up to `000029_fine_grained_permissions.sql`

**Schemas**:
- `public`: Main content (movies only, TV/music/etc. planned)
- `shared`: Users, sessions, settings, RBAC, activity, libraries
- `qar`: Adult content (planned, placeholder only)

### 3.2 Cache (Dragonfly/Redis) - ✅ Complete

- **Status**: Fully implemented with testcontainers tests
- **Coverage**: 65.7%
- **Location**: `internal\infra\cache\`
- **Features**:
  - L1 cache (Otter in-memory)
  - L2 cache (Rueidis client for Dragonfly/Redis)
  - Key namespacing
  - TTL support
  - Pattern-based invalidation
  - JSON serialization
  - Integration tests with real Dragonfly container

**No gaps identified**.

### 3.3 Search (Typesense) - 🟡 Partial

- **Status**: Client complete, only movie search implemented
- **Coverage**: 37.0% (search service), 56.1% (infra client)
- **Location**: `internal\infra\search\` + `internal\service\search\`
- **Complete**:
  - Typesense client with health checks
  - Movie collection schema
  - Movie indexing service
  - Search API handlers
- **Missing**:
  - TV show search schema
  - Music search schema
  - Multi-content search
  - Faceted search implementation
  - Real-time index updates via webhooks

**TODO Found**:
- `internal\api\handler_search.go:194`: "TODO: This should be an async job via River" (ReindexSearch endpoint)

### 3.4 Job Queue (River) - ✅ Complete

- **Status**: Fully implemented
- **Coverage**: 54.4%
- **Location**: `internal\infra\jobs\`
- **Features**:
  - 5-level priority queues (critical, high, default, low, bulk)
  - PostgreSQL-backed
  - Transactional job enqueueing
  - Graceful shutdown
  - Worker pool configuration
- **Job Types Implemented**:
  - Movie library scan
  - Movie metadata refresh
  - Movie search indexing
  - File matching
  - Radarr sync
  - Activity log cleanup
  - Notification dispatch

**Gaps**:
- Jobs for TV, music, and other content types not yet defined
- Max attempts hardcoded (`TODO: Make configurable` at `internal\infra\jobs\module.go:40`)

### 3.5 Webhooks - ✅ Complete (for Radarr)

- **Status**: Implemented for Radarr only
- **Location**: `internal\integration\radarr\webhook_handler.go`
- **Complete**:
  - Radarr webhook receiver
  - Event-driven movie sync
  - Webhook validation
- **Missing**:
  - Sonarr webhooks
  - Lidarr webhooks
  - Whisparr webhooks
  - Generic webhook system for external integrations

---

## 4. Integration Services

### 4.1 Arr Ecosystem

| Service | Status | Location | Implementation | Tests |
|---------|--------|----------|----------------|-------|
| **Radarr** | ✅ Complete | `internal\integration\radarr\` | Client, sync service, webhooks, jobs, mapper | Yes (43.3%) |
| **Sonarr** | 🔴 Not Started | - | Nothing | No |
| **Lidarr** | 🔴 Not Started | - | Nothing | No |
| **Whisparr** | 🔴 Not Started | - | Nothing | No |
| **Readarr** | 🔴 Not Started | - | Nothing | No |
| **Prowlarr** | 🔴 Not Started | - | Nothing | No |

**Radarr Integration Details** (Reference Implementation):
- **Files**: 12 Go files
- **Features**:
  - REST API client with authentication
  - Full movie sync service
  - Webhook handler for real-time updates
  - River jobs for background sync
  - Mapper for Radarr→Revenge entities
  - Quality profile + root folder queries
  - Sync status tracking
- **API Endpoints Defined But Not Implemented**:
  - `GET /api/v1/admin/integrations/radarr/status` (stub)
  - `GET /api/v1/admin/integrations/radarr/quality-profiles` (stub)
  - `GET /api/v1/admin/integrations/radarr/root-folders` (stub)

### 4.2 Metadata Providers

| Provider | Module(s) | Status | Location | Notes |
|----------|-----------|--------|----------|-------|
| **TMDb** | Movie | ✅ Complete | `internal\content\movie\tmdb_*.go` | Full client, mapper, types, tests |
| **TheTVDB** | TV | 🔴 Not Started | - | Design doc exists |
| **MusicBrainz** | Music | 🔴 Not Started | - | Design doc exists |
| **Last.fm** | Music | 🔴 Not Started | - | Design doc exists |
| **Audnexus** | Audiobook | 🔴 Not Started | - | Design doc exists |
| **OpenLibrary** | Book | 🔴 Not Started | - | Design doc exists |
| **ComicVine** | Comics | 🔴 Not Started | - | Design doc exists |
| **StashDB** | QAR | 🔴 Not Started | - | Design doc exists |
| **ThePornDB** | QAR | 🔴 Not Started | - | Design doc exists |

### 4.3 Scrobbling Services

| Service | Content | Status | Location |
|---------|---------|--------|----------|
| Trakt | Movies, TV | 🔴 Not Started | - |
| Last.fm | Music | 🔴 Not Started | - |
| ListenBrainz | Music | 🔴 Not Started | - |
| Letterboxd | Movies | 🔴 Not Started | - |
| Simkl | Movies, TV, Anime | 🔴 Not Started | - |

**No files exist** for any scrobbling service.

---

## 5. Frontend (SvelteKit)

### 5.1 Frontend Status: **🔴 DOES NOT EXIST**

**Evidence**: No frontend directory found in the codebase.

**Glob searches returned no results** for:
- `frontend/**/*`
- `web/**/*`
- `ui/**/*`
- `**/*.svelte`

### 5.2 Frontend Roadmap (TODO_B_FRONTEND.md)

**Total Effort**: 40-60 hours

**Planned Components** (all pending):
1. SvelteKit 2 + Svelte 5 setup
2. Authentication pages (login, register, MFA)
3. Layout & navigation
4. Library browser (movies grid)
5. Movie detail pages
6. Search interface
7. Video player (HLS.js)
8. Settings pages
9. Admin dashboard

**API Coverage**: 0% - Backend has 92 API endpoints defined, frontend cannot consume any of them

### 5.3 API Client Generation

- OpenAPI spec exists: `api\openapi\openapi.yaml` (6,594 lines, 92 endpoints)
- No TypeScript client generated
- No fetch/axios wrapper
- No Svelte stores for API state

---

## 6. Missing Test Coverage (Below 80% Target)

### 6.1 Services Below Coverage Target

| Service | Current | Target | Gap | Priority |
|---------|---------|--------|-----|----------|
| MFA | 10.4% | 80% | 69.6% | HIGH |
| RBAC | 39.6% | 80% | 40.4% | HIGH |
| OIDC | 27.6% | 80% | 52.4% | MEDIUM |
| Activity | 31.6% | 80% | 48.4% | MEDIUM |
| Auth | 38.4% | 80% | 41.6% | HIGH |
| Search | 37.0% | 80% | 43.0% | MEDIUM |
| Settings | 41.0% | 80% | 39.0% | MEDIUM |
| User | 45.9% | 80% | 34.1% | MEDIUM |
| Movie | 41.0% | 80% | 39.0% | MEDIUM |
| Library | 43.6% | 80% | 36.4% | MEDIUM |

### 6.2 Content Modules with 0% Coverage

All generated database query packages have **0% coverage**:
- `internal\content\movie\db` - 0%
- `internal\content\tvshow\db` - 0%
- `internal\content\qar\db` - 0%
- `internal\infra\database\db` - 0%

These are sqlc-generated files and typically don't need direct testing, but integration tests should exercise them.

---

## 7. TODO Comments & Stub Implementations

### 7.1 Critical TODOs in Code

| File | Line | TODO | Impact |
|------|------|------|--------|
| `internal\api\handler.go` | 524 | Handle notification settings (JSONB fields) | User preferences incomplete |
| `internal\api\handler.go` | 709 | Extract IP, user agent, fingerprint from request | Auth security incomplete |
| `internal\api\handler.go` | 823 | Extract IP address and user agent from request | Session tracking incomplete |
| `internal\api\handler_oidc.go` | 77 | Implement custom redirect middleware or JSON response | OIDC UX issue |
| `internal\api\handler_search.go` | 194 | ReindexSearch should be async River job | Performance issue |
| `internal\service\email\service.go` | 205 | Implement SendGrid API call | Feature incomplete |
| `internal\content\movie\service.go` | 275 | Implement metadata refresh via River job | Feature stub |
| `internal\infra\jobs\module.go` | 40 | Make MaxAttempts configurable | Hardcoded config |

### 7.2 Unimplemented API Endpoints

**OpenAPI spec defines 92 endpoints**, but many are stubs returning `ErrNotImplemented`:

**RBAC Endpoints** (stubs):
- `POST /api/v1/rbac/policies` - AddPolicy
- All policy management endpoints

**OIDC Admin Endpoints** (stubs):
- `POST /api/v1/admin/oidc/providers` - AdminCreateOIDCProvider
- `GET /api/v1/admin/oidc/providers/{providerId}` - AdminGetOIDCProvider
- `DELETE /api/v1/admin/oidc/providers/{providerId}` - AdminDeleteOIDCProvider
- `POST /api/v1/admin/oidc/providers/{providerId}/enable` - AdminEnableOIDCProvider
- `POST /api/v1/admin/oidc/providers/{providerId}/disable` - AdminDisableOIDCProvider

**Radarr Admin Endpoints** (stubs):
- `GET /api/v1/admin/integrations/radarr/status` - AdminGetRadarrStatus
- `GET /api/v1/admin/integrations/radarr/quality-profiles` - AdminGetRadarrQualityProfiles
- `GET /api/v1/admin/integrations/radarr/root-folders` - AdminGetRadarrRootFolders

---

## 8. Feature Completeness Matrix

### 8.1 Core Features

| Feature | Backend | Frontend | Integration | Status |
|---------|---------|----------|-------------|--------|
| User registration | ✅ | 🔴 | N/A | Backend only |
| Login/Logout | ✅ | 🔴 | N/A | Backend only |
| MFA (TOTP) | ✅ | 🔴 | N/A | Backend only |
| MFA (WebAuthn) | ✅ | 🔴 | N/A | Backend only |
| Session management | ✅ | 🔴 | N/A | Backend only |
| RBAC | ✅ | 🔴 | N/A | Backend only |
| OIDC SSO | ✅ | 🔴 | N/A | Backend only |
| API Keys | ✅ | 🔴 | N/A | Backend only |
| Activity logs | ✅ | 🔴 | N/A | Backend only |
| Settings | ✅ | 🔴 | N/A | Backend only |

### 8.2 Content Features

| Feature | Movies | TV | Music | Books | QAR | Other |
|---------|--------|----|----|-------|-----|-------|
| Library management | ✅ | 🔴 | 🔴 | 🔴 | 🔴 | 🔴 |
| File scanning | ✅ | 🔴 | 🔴 | 🔴 | 🔴 | 🔴 |
| Metadata fetching | ✅ | 🔴 | 🔴 | 🔴 | 🔴 | 🔴 |
| Search indexing | ✅ | 🔴 | 🔴 | 🔴 | 🔴 | 🔴 |
| Collections | ✅ | 🔴 | 🔴 | 🔴 | 🔴 | 🔴 |
| Watch progress | ✅ | 🔴 | N/A | N/A | 🔴 | N/A |
| Playback | 🔴 | 🔴 | 🔴 | 🔴 | 🔴 | 🔴 |
| Transcoding | 🔴 | 🔴 | 🔴 | 🔴 | 🔴 | 🔴 |

### 8.3 Integration Features

| Integration | Status | Completeness |
|-------------|--------|--------------|
| Radarr | ✅ | 90% (admin endpoints stub) |
| Sonarr | 🔴 | 0% |
| Lidarr | 🔴 | 0% |
| Whisparr | 🔴 | 0% |
| Readarr | 🔴 | 0% |
| Prowlarr | 🔴 | 0% |
| TMDb | ✅ | 100% |
| TheTVDB | 🔴 | 0% |
| MusicBrainz | 🔴 | 0% |
| Trakt | 🔴 | 0% |
| Notifications | ✅ | 80% (4 agents working) |

---

## 9. Key Findings Summary

### 9.1 Strengths

1. **Solid Infrastructure Foundation**
   - PostgreSQL with 29 migrations and self-healing pool
   - L1+L2 caching with testcontainers tests
   - River job queue with 5-level priorities
   - Comprehensive auth system (JWT, MFA, RBAC, OIDC, API keys)
   - Activity logging and health checks

2. **Movie Module as Reference Implementation**
   - Complete CRUD operations
   - TMDb metadata integration
   - Library scanning with MediaInfo
   - Radarr sync with webhooks
   - Search indexing
   - River jobs for background tasks
   - 41% test coverage (approaching target)

3. **Good Code Organization**
   - Clear separation: `content/`, `service/`, `infra/`, `integration/`
   - Repository pattern with interfaces
   - Dependency injection via fx
   - ogen-generated type-safe API handlers
   - sqlc-generated database queries

### 9.2 Critical Gaps

1. **Frontend Does Not Exist** (🚨 Blocker)
   - No SvelteKit app
   - No UI components
   - No API client
   - 40-60 hours of work required

2. **Only 1 of 11 Content Modules Complete**
   - TV Show: 0% (placeholder only)
   - Music: 0%
   - Audiobooks: 0%
   - Books: 0%
   - Podcasts: 0%
   - Photos: 0%
   - Comics: 0%
   - LiveTV: 0%
   - QAR (Adult): 0% (placeholder only)

3. **Only 1 of 9 Arr Integrations Complete**
   - Radarr: ✅ 90%
   - Sonarr: 🔴 0%
   - Lidarr: 🔴 0%
   - Whisparr: 🔴 0%
   - Others: 🔴 0%

4. **No Playback/Streaming**
   - No video player
   - No transcoding
   - No HLS/DASH streaming
   - No subtitle support

5. **Test Coverage Below 80% Target**
   - 11 services below target
   - MFA at 10.4% (critical gap)
   - RBAC at 39.6% (critical gap)

---

## 10. Recommendations

### 10.1 Phase Prioritization

**Critical Path to MVP**:
1. **Complete A6.6 test coverage** (8-12h) - Ensure stability before frontend work
2. **Build frontend Phase B** (40-60h) - Unblock end-to-end functionality
3. **Complete Phase C infrastructure** (8-16h) - Enable deployment
4. **Post-MVP**: Expand to TV, music, other content types

### 10.2 Test Coverage Strategy

**Focus on high-risk, low-coverage areas**:
1. MFA (10.4% → 80%): Critical for security
2. RBAC (39.6% → 80%): Critical for authorization
3. Auth (38.4% → 80%): Critical for authentication

### 10.3 Content Module Strategy

**After MVP**:
1. Clone movie module structure for TV shows
2. Reuse patterns: repository, service, jobs, scanner, search
3. Add TV-specific logic: seasons, episodes, series tracking
4. Integrate Sonarr (similar to Radarr)
5. Repeat for music, books, etc.

---

## 11. Timeline Estimates

From `TODO_INDEX.md`:
- **Phase A remaining**: 8-12h (test coverage)
- **Phase B**: 40-60h (frontend)
- **Phase C**: 8-16h (infrastructure)
- **Total remaining**: 56-88 hours (1.5-2 weeks full-time)

---

## Conclusion

The Revenge codebase has **excellent infrastructure foundations** and a **complete movie module** that serves as a reference implementation. However, the project is **30% complete overall** with the **frontend entirely missing** and **90% of content modules not yet started**.

**Path to MVP**: Focus on completing test coverage (A6.6), building the frontend (Phase B), and deploying infrastructure (Phase C). Post-MVP, replicate the movie module pattern for TV, music, and other content types.

**Estimated completion**: 56-88 hours (1.5-2 weeks full-time) for MVP, then expand content modules iteratively.
