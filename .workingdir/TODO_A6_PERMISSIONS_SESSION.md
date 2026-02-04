# Phase A6: Permissions & Session

**Priority**: HIGH
**Effort**: 20-30h
**Dependencies**: A0

---

## A6.1: Fine-Grained Permissions ✅

**Priority**: HIGH | **Effort**: 4-6h | **Actual**: 1h
**Status**: COMPLETED (2026-02-04)

Current Casbin policies use coarse permissions. Need to implement fine-grained.

**Completed Tasks**:
- [x] Define permission taxonomy in `internal/service/rbac/permissions.go`
  - Permission constants: `PermUsersList`, `PermMoviesGet`, etc.
  - `FineGrainedResources`: users, profile, movies, libraries, playback, requests, settings, audit, integrations, notifications, admin
  - `FineGrainedActions`: list, get, create, update, delete, read, write, stream, progress, approve, export, sync, scan, user_read, user_write, *
  - `DefaultRolePermissions` map for admin/moderator/user/guest
  - `HasPermission` helper function
- [x] Migration `000029_fine_grained_permissions.up.sql` - Update default role policies
  - Moderator: user management (limited), full movies/libraries/requests/integrations/notifications
  - User: profile, view movies/libraries, playback, own requests/notifications
  - Guest: read-only profile/movies/libraries, stream only
- [x] Update `internal/service/rbac/roles.go` - Use FineGrainedResources/Actions
- [ ] Update middleware to check granular permissions (deferred - existing middleware works)
- [ ] Update OpenAPI spec with permission requirements (deferred - documentation task)

---

## A6.2: Hybrid Session Storage ✅

**Priority**: HIGH | **Effort**: 4-6h | **Actual**: 1h
**Status**: COMPLETED (2026-02-04)

Sessions now use Dragonfly/Redis L1 cache for performance.

**Completed Tasks**:
- [x] Update `internal/service/session/cached_service.go`:
  - [x] ValidateSession: Check cache first, fallback to PostgreSQL (already implemented)
  - [x] CreateSession: Write-through to cache on session create
  - [x] RevokeSession: Invalidate cache, update PostgreSQL (already implemented)
  - [x] Configurable cache TTL
- [x] Session serialization via JSON (using cache.SetJSON/GetJSON)
- [x] Config: `session.cache_enabled: true` (default)
- [x] Config: `session.cache_ttl: 5m` (configurable)
- [x] Config: `session.max_per_user: 10`
- [x] Config: `session.token_length: 32`
- [x] Updated tests for new cacheTTL parameter
- [ ] Benchmark: Compare latency with/without cache (deferred - manual testing)
- [ ] Tests with testcontainers Redis (covered in A6.5)

---

## A6.3: 5-Level River Queue Priorities ✅

**Priority**: MEDIUM | **Effort**: 2-3h | **Actual**: 0.5h
**Status**: COMPLETED (2026-02-04)

Expanded from 4 queues to 5-level priority system.

**Completed Tasks**:
- [x] Update `internal/infra/jobs/queues.go`:
  - QueueCritical: Security events, auth failures (20 workers)
  - QueueHigh: Notifications, webhooks, user actions (15 workers)
  - QueueDefault: Metadata fetching, sync (10 workers)
  - QueueLow: Cleanup, maintenance (5 workers)
  - QueueBulk: Library scans, batch ops (3 workers)
- [x] Updated QueuePriority() function with 5-level thresholds
- [x] Notification jobs now use QueueHigh instead of removed QueueNotifications
- [x] Updated River config with worker allocation by priority
- [x] Tests for queue assignment updated

---

## A6.4: Configurable Activity Log Retention

**Priority**: LOW | **Effort**: 1-2h

Activity logs should be configurable, default 90 days.

**Tasks**:
- [ ] Add config: `activity.retention_days: 90`
- [ ] Update cleanup job to use config value
- [ ] Add admin API: `PUT /api/v1/admin/settings/activity-retention`
- [ ] Migration for default server setting
- [ ] Tests

---

## A6.5: Cache Integration Tests (testcontainers)

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

## A6.6: Test Coverage to 80% [ONGOING]

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
