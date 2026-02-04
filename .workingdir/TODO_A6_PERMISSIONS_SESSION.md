# Phase A6: Permissions & Session

**Priority**: HIGH
**Effort**: 20-30h
**Dependencies**: A0

---

## A6.1: Fine-Grained Permissions

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

## A6.2: Hybrid Session Storage

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

## A6.3: 5-Level River Queue Priorities

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
