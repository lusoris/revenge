# Revenge - Development Roadmap

> Modular media server with complete content isolation

**Last Updated**: 2026-01-30
**Current Phase**: Pre-Testing Implementation
**Build**: `GOEXPERIMENT=greenteagc,jsonv2 go build ./...`

---

## Quick Status

```
Foundation (Week 1-2)     ████████████████████████ 100%
Design Audit              ████████████████████████ 100%
Critical Fixes            ████████████████████████ 100% ✓
Library Refactor          ████████████████████████ 100% ✓
Movie Module              ████████████████████████ 100% ✓
TV Shows Module           ████████████████████████ 100% ✓
Adult Module (QAR)        ████████████████████████ 100% ✓
Pre-Test Implementation   ██████████░░░░░░░░░░░░░░  40%  <- CURRENT
Unit Tests                ░░░░░░░░░░░░░░░░░░░░░░░░   0%
Integration Tests         ░░░░░░░░░░░░░░░░░░░░░░░░   0%
Music Module              ░░░░░░░░░░░░░░░░░░░░░░░░   0%
Books Module              ░░░░░░░░░░░░░░░░░░░░░░░░   0%
Comics Module             ░░░░░░░░░░░░░░░░░░░░░░░░   0%
Frontend                  ░░░░░░░░░░░░░░░░░░░░░░░░   0%
```

---

## 🎯 Current Goal: Complete Everything Before Tests

> Implement all features that don't require tests to validate, then write comprehensive test suites.

---

## Phase A: Pre-Unit-Test Implementation

> Everything that can be built without needing tests to verify correctness.

### A.1 QAR Relationship Handlers (Complete Stubs)
- [x] **ListAdultPerformerMovies** - `expedition.ListByPerformer()` method
  - [x] Add repository method: `ListByCrewID(ctx, crewID, limit, offset)`
  - [x] Add service method: `ListByPerformer(ctx, performerID, limit, offset)`
  - [x] Wire to handler
- [x] **ListAdultStudioMovies** - `expedition.ListByPort()` method
  - [x] Add repository method: `ListByPortID(ctx, portID, limit, offset)`
  - [x] Add service method: `ListByStudio(ctx, studioID, limit, offset)`
  - [x] Wire to handler
- [x] **ListAdultTagMovies** - `expedition.ListByFlag()` method
  - [x] Add repository method: `ListByFlagID(ctx, flagID, limit, offset)`
  - [x] Add service method: `ListByTag(ctx, tagID, limit, offset)`
  - [x] Wire to handler
- [ ] **ListAdultSimilarMovies** - needs similar recommendation logic
  - Based on shared flags/crew/port
- [ ] **ListAdultMovieMarkers** - needs marker/chapter entity
  - Create `qar/marker/` module or add to voyage

### A.2 QAR Request System
- [ ] **SearchAdultRequests** - implement full search
- [ ] **ListAdultRequests** - implement user request listing
- [ ] **CreateAdultRequest** - create download/metadata requests
- [ ] **GetAdultRequest** - get single request
- [ ] **VoteAdultRequest** - user voting on requests
- [ ] **CommentAdultRequest** - comments on requests
- [ ] **ListAdultAdminRequests** - admin view of all requests
- [ ] **ApproveAdultRequest** - admin approval
- [ ] **DeclineAdultRequest** - admin decline
- [ ] **UpdateAdultRequestQuota** - admin quota management
- [ ] **ListAdultRequestRules** - auto-approval rules
- [ ] **CreateAdultRequestRule** - create auto-rules
- [ ] **UpdateAdultRequestRule** - modify rules
- [ ] **DeleteAdultRequestRule** - remove rules

### A.3 External Metadata Integrations
- [ ] **StashAppClient** → `internal/service/metadata/stash_app/`
  - [ ] types.go - Stash-App GraphQL types
  - [ ] client.go - GraphQL client with circuit breaker
  - [ ] provider.go - Implements metadata provider interface
  - [ ] module.go - fx wiring
  - Features:
    - Import scenes from local Stash instance
    - Sync scene markers as chapters
    - Import user ratings
    - One-way sync (Stash → Revenge)
- [ ] **StashDB Search handlers**
  - [ ] SearchAdultStashDBScenes - search StashDB
  - [ ] GetAdultStashDBScene - get scene details
  - [ ] SearchAdultStashDBPerformers - search performers
  - [ ] GetAdultStashDBPerformer - get performer details
  - [ ] IdentifyAdultStashDBScene - fingerprint lookup
- [ ] **TPDB handlers**
  - [ ] SearchAdultTPDBScenes - search TPDB
  - [ ] GetAdultTPDBScene - get scene
  - [ ] GetAdultTPDBPerformer - get performer
- [ ] **Stash-App sync handlers**
  - [ ] SyncAdultStash - sync with Stash-App
  - [ ] ImportAdultStash - import from Stash-App
  - [ ] GetAdultStashStatus - connection status

### A.4 Playback Service
- [x] **Create service** → `internal/service/playback/`
  - [x] service.go - Playback orchestration
  - [x] types.go - PlaybackSession, StreamInfo, etc.
  - [x] module.go - fx wiring
- [x] **Implement core methods**
  - [x] StartPlayback(ctx, userID, mediaID, mediaType) → PlaybackSession
  - [x] UpdateProgress(ctx, sessionID, positionTicks)
  - [x] StopPlayback(ctx, sessionID)
  - [x] GetActiveSession(ctx, userID, mediaID)
- [ ] **Up Next / Auto-Play Queue** (stub implemented)
  - [ ] BuildUpNextQueue(ctx, userID, currentMediaID) → []MediaItem
  - [ ] TV: next episode in series
  - [ ] Movie: similar movies or collection next
  - [ ] QAR: similar expeditions
- [ ] **API endpoints**
  - [ ] POST /api/playback/start
  - [ ] PUT /api/playback/{sessionId}/progress
  - [ ] POST /api/playback/{sessionId}/stop
  - [ ] GET /api/playback/up-next

### A.5 Cross-Device Sync (Basic)
- [ ] **Polling endpoint** `/api/sync/playback?since={ts}`
  - Returns playback state changes since timestamp
  - Lightweight alternative to WebSocket
- [ ] **BroadcastToUser()** - notify all user sessions of changes

### A.6 RBAC Completion
- [ ] **Missing Casbin methods** in `internal/service/rbac/casbin.go`
  - [ ] Enforce(subject, object, action) bool
  - [ ] AddRoleForUser(userID, role)
  - [ ] RemoveRoleForUser(userID, role)
  - [ ] GetRolesForUser(userID) []string
  - [ ] GetUsersForRole(role) []string
  - [ ] GetPermissionsForUser(userID) [][]string
- [ ] **Resource grants table** (polymorphic permissions)
  - [ ] Migration: `shared/000021_resource_grants.up.sql`
  - [ ] HasGrant(userID, resourceType, resourceID, permission) bool
  - [ ] CreateGrant(userID, resourceType, resourceID, permission)
  - [ ] DeleteGrant(grantID)
  - [ ] DeleteByResource(resourceType, resourceID)
- [ ] **Missing permissions** to seed
  - [ ] access.rules.view, access.rules.manage, access.bypass
  - [ ] request.* permissions (15 total)
  - [ ] adult.request.* permissions (7 total)

### A.7 Health Checks
- [x] **Enable cache health check** via infra/health module
- [x] **Enable search health check** via infra/health module
- [x] **Enable jobs health check** via infra/health module
- [ ] **Add QAR-specific health** - check adult module enabled

### A.8 User Preferences
- [ ] **Add preference fields** to user_profiles table
  - [ ] auto_play_enabled boolean DEFAULT true
  - [ ] auto_play_delay_seconds int DEFAULT 10
  - [ ] continue_watching_days int DEFAULT 30
  - [ ] mark_watched_percent int DEFAULT 90
  - [ ] adult_pin_hash text (for PIN protection)
- [ ] **API endpoints**
  - [ ] GET /api/users/me/preferences
  - [ ] PUT /api/users/me/preferences
- [ ] **Implement PIN protection** (optional adult content lock)

### A.9 Audit Logging
- [ ] **Redesign activity_log schema**
  - [ ] Add: module, entity_id, entity_type, changes (JSONB)
  - [ ] Partition by month for performance
- [ ] **River async worker** for audit writes
  - [ ] AuditLogWorker - fire-and-forget audit entries
- [ ] **Adult access audit** - log all QAR access
  - [ ] user_id, resource_type, resource_id, action, timestamp, ip

### A.10 Documentation Cleanup
- [x] **Remove bogus UPSTREAM_SYNC.md** (was hallucinated)
- [x] **Remove sync scripts** (sync-upstream.sh, sync-upstream.ps1)
- [ ] **Update TECH_STACK.md** - add casbin, otel
- [ ] **Update CONFIGURATION.md** - reflects pkg/config/
- [ ] **Review all docs/dev/design/** for accuracy

---

## Phase B: Unit Test Suite

> After Phase A, create comprehensive unit tests.

### B.1 Repository Tests
- [ ] **Movie repository tests** - all CRUD + queries
- [ ] **TVShow repository tests** - series, seasons, episodes
- [ ] **QAR repository tests** - expedition, voyage, crew, port, flag, fleet

### B.2 Service Tests
- [ ] **Movie service tests** - business logic
- [ ] **TVShow service tests** - business logic
- [ ] **QAR service tests** - all 5 services
- [ ] **Playback service tests** - session management
- [ ] **RBAC service tests** - permission checks

### B.3 Handler Tests
- [ ] **Movie handler tests** - HTTP layer
- [ ] **TVShow handler tests** - HTTP layer
- [ ] **QAR handler tests** - HTTP layer, auth checks

### B.4 Utility Tests
- [ ] **Fingerprint service tests** - oshash, phash, md5
- [ ] **Config tests** - loading, validation
- [ ] **Cache tests** - otter, sturdyc integration

---

## Phase C: Integration/Feature Tests

> After unit tests, create integration tests.

### C.1 Database Integration
- [ ] **Migration tests** - up/down migrations work
- [ ] **Transaction tests** - rollback behavior
- [ ] **Concurrent access tests** - race conditions

### C.2 External Service Integration
- [ ] **TMDb integration tests** - mock server
- [ ] **Radarr integration tests** - mock server
- [ ] **Sonarr integration tests** - mock server
- [ ] **StashDB integration tests** - mock GraphQL
- [ ] **Typesense integration tests** - test container

### C.3 End-to-End Workflows
- [ ] **User registration → login → browse → play**
- [ ] **Library scan → metadata fetch → index**
- [ ] **Watch progress → continue watching → complete**
- [ ] **QAR access control → browse → stream**

### C.4 Performance Tests
- [ ] **API response time benchmarks**
- [ ] **Database query benchmarks**
- [ ] **Cache hit rate tests**
- [ ] **Concurrent user load tests**

---

## Completed (2026-01-30)

### QAR Module (100%)
- [x] Full schema obfuscation (Queen Anne's Revenge)
- [x] All 6 domain packages (expedition, voyage, crew, port, flag, fleet)
- [x] All repositories with sqlc
- [x] All services with business logic
- [x] QAR API handlers (~50 endpoints) in `internal/api/adult.go`
- [x] QAR converters in `internal/api/converters.go`
- [x] Handler wiring in `internal/api/module.go`
- [x] Search support in List handlers (Query parameter)
- [x] Fingerprinting handlers (Identify, Match)
- [x] RBAC adult permissions (adult.browse, adult.stream, adult.metadata.write)
- [x] WhisparrClient with circuit breaker
- [x] StashDB GraphQL client
- [x] FingerprintService (oshash + pHash + MD5)
- [x] Typesense collections (5 isolated collections)

### Continue Watching
- [x] 30-day filter for movies (`ListResumeableMovies`)
- [x] 30-day filter for TV episodes (`ListResumeableEpisodes`)
- [x] 30-day filter for TV series (`ListContinueWatchingSeries`)

### Movie Module (100%)
- [x] Full CRUD with relations
- [x] User data (ratings, favorites, watchlist)
- [x] TMDb metadata provider
- [x] Radarr metadata provider
- [x] River jobs for metadata enrichment
- [x] 30-day continue watching filter

### TV Shows Module (100%)
- [x] Database migrations
- [x] sqlc queries (100+ queries)
- [x] Entity definitions
- [x] Repository (PostgreSQL)
- [x] Service layer
- [x] API handlers
- [x] 30-day continue watching filter

---

## Tech Stack Reference

| Component | Package | Notes |
|-----------|---------|-------|
| Cache (distributed) | `github.com/redis/rueidis` | NOT go-redis |
| Cache (local) | `github.com/maypok86/otter` v1.2.4 | W-TinyLFU |
| Cache (API) | `github.com/viccon/sturdyc` v1.1.5 | Request coalescing |
| Search | `github.com/typesense/typesense-go/v4` | NOT v3 |
| Config | `github.com/knadh/koanf/v2` | NOT viper |
| Logging | `log/slog` | NOT zap |
| Jobs | `github.com/riverqueue/river` | PostgreSQL-native |
| RBAC | `github.com/casbin/casbin/v2` | Dynamic roles |
| DI | `go.uber.org/fx` | Dependency injection |
| HTTP client | `github.com/go-resty/resty/v2` | External APIs |

---

## Build Commands

```bash
# With experiments
GOEXPERIMENT=greenteagc,jsonv2 go build -o bin/revenge ./cmd/revenge

# Generate code
sqlc generate
go generate ./api/...

# Lint
golangci-lint run

# Test (after Phase B)
go test ./...
```

---

## Important Notes

**Adult Content** (Queen Anne's Revenge obfuscation):
- Schema: `qar` (isolated PostgreSQL schema)
- API namespace: `/api/v1/qar/*`
- Module location: `internal/content/qar/`
- See [ADULT_CONTENT_SYSTEM.md](docs/dev/design/features/adult/ADULT_CONTENT_SYSTEM.md)

**Design Docs are Source of Truth**:
- Only `docs/dev/design/` is authoritative
- Other documentation may be outdated
- Code must match design, not vice versa
