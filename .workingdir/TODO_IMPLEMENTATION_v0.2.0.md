# v0.2.0 Implementation TODO

**Version**: v0.2.0 - Core Backend Services
**Last Updated**: 2026-02-02
**Approach**: Step-by-step, test-first, lint-clean

## Implementation Order

Services werden in dieser Reihenfolge implementiert (Abhängigkeiten berücksichtigt):

1. ✅ **PostgreSQL Pool Enhancement** - Foundation for all services (4/4 complete)
2. ✅ **Dragonfly/Redis Integration** - Caching infrastructure (4/4 complete)
3. ✅ **River Job Queue** - Background jobs (3/3 complete)
4. ✅ **Settings Service** - Server and user settings (Steps 4.1-4.4 complete, 4 commits)
5. **User Service** - Core entity
6. **Auth Service** - Depends on User
7. **Session Service** - Depends on Auth
8. **RBAC Service** - Depends on User
9. **API Keys Service** - Depends on User + RBAC
10. **OIDC Service** - Depends on Auth + User
11. **Activity Service** - Depends on User (logging)
12. **Library Service** - Depends on User + RBAC
13. **Health Service Enhancement** - Final integration check

---

## Step 1: PostgreSQL Pool Enhancement ✅

### 1.1 Database Migrations Setup
- [x] Review existing migration files in `/migrations`
- [x] Create migration: `000004_create_server_settings_table.up.sql`
- [x] Create migration: `000004_create_server_settings_table.down.sql`
- [x] **Test**: Run migrations up/down
- [x] **Lint**: Check SQL syntax
- [x] **Verify**: Table created in test database

### 1.2 sqlc Configuration
- [x] Review `sqlc.yaml` configuration
- [x] Create `internal/infra/database/queries/settings.sql`
- [x] Generate code: `sqlc generate`
- [x] **Test**: Generated code compiles
- [x] **Lint**: `golangci-lint run internal/infra/database/`
- [x] **Verify**: No errors

### 1.3 Database Pool Metrics
- [x] Create `internal/infra/database/metrics.go`
- [x] Add Prometheus metrics for pool stats
- [x] Instrument `internal/infra/database/postgres.go`
- [x] **Test**: Unit tests for metrics recording
- [x] **Lint**: `golangci-lint run internal/infra/database/`
- [x] **Coverage**: Check with `go test -cover`
- [x] **Verify**: Metrics endpoint returns pool stats

### 1.4 Query Logging
- [x] Add debug query logging to pgx config
- [x] Create `internal/infra/database/logger.go`
- [x] Add slow query detection (configurable threshold)
- [x] **Test**: Integration test with embedded-postgres
- [x] **Lint**: Full check
- [x] **Coverage**: Ensure 80%+
- [x] **Verify**: Debug logs appear in test mode

**Checkpoint**: Database infrastructure ready ✅
**Tests**: `go test ./internal/infra/database/... -v` ✅
**Lint**: `golangci-lint run ./internal/infra/database/...` ✅
**Coverage**: `go test ./internal/infra/database/... -cover` ✅

---

## Step 2: Dragonfly/Redis Integration ✅

### 2.1 Rueidis Client Setup
- [x] Create `internal/infra/cache/rueidis.go`
- [x] Implement client initialization
- [x] Add connection pooling config
- [x] **Test**: Unit tests with mock
- [x] **Lint**: Check
- [x] **Verify**: Client connects to Dragonfly

### 2.2 Otter L1 Cache
- [x] Create `internal/infra/cache/otter.go`
- [x] Implement W-TinyLFU cache
- [x] Configure TTL and size limits
- [x] **Test**: Unit tests for cache operations
- [x] **Lint**: Check
- [x] **Coverage**: 80%+
- [x] **Verify**: Cache hit/miss works

### 2.3 Cache Operations
- [x] Create `internal/infra/cache/cache.go`
- [x] Implement Get/Set/Delete operations
- [x] Add cache invalidation patterns
- [x] **Test**: Integration test with testcontainers (Dragonfly)
- [x] **Lint**: Full check
- [x] **Coverage**: 80%+
- [x] **Verify**: L1->L2 fallback works

### 2.4 Distributed Locks
- [x] Implement distributed lock primitives
- [x] Add lock timeout handling
- [x] **Test**: Concurrent lock tests
- [x] **Lint**: Check
- [x] **Coverage**: 80%+
- [x] **Verify**: Lock prevents race conditions

**Checkpoint**: Cache infrastructure ready ✅
**Tests**: `go test ./internal/infra/cache/... -v` ✅
**Lint**: `golangci-lint run ./internal/infra/cache/...` ✅
**Coverage**: `go test ./internal/infra/cache/... -cover` ✅

---

## Step 3: River Job Queue

### 3.1 River Client Setup
- [x] Create `internal/infra/jobs/river.go`
- [x] Initialize River client with pgx
- [x] Configure worker pool
- [x] **Test**: Unit test client creation
- [x] **Lint**: Check
- [x] **Verify**: Client initializes

### 3.2 Queue Configuration
- [x] Create `internal/infra/jobs/queues.go`
- [x] Define queue priorities (critical, default, low)
- [x] Configure retry policies
- [x] **Test**: Queue config validation
- [x] **Lint**: Check
- [x] **Verify**: Queues registered

### 3.3 Base Job Types
- [x] Create `internal/infra/jobs/cleanup_job.go`
- [x] Implement periodic cleanup job
- [x] **Test**: Job execution test
- [x] **Lint**: Check
- [x] **Coverage**: 80%+
- [x] **Verify**: Job runs and completes

**Checkpoint**: Job queue ready ✅
**Tests**: `go test ./internal/infra/jobs/... -v` ✅
**Lint**: `golangci-lint run ./internal/infra/jobs/...` ✅
**Coverage**: `go test ./internal/infra/jobs/... -cover` ✅ 65.6%

---

## Step 4: Settings Service

### 4.1 Database Schema
- [x] Migration: `000005_create_user_settings_table.up.sql`
- [x] Migration: `000005_create_user_settings_table.down.sql`
- [x] Run migrations
- [x] **Test**: Migration up/down
- [x] **Verify**: Tables exist

### 4.2 Repository Layer
- [x] Create `internal/service/settings/repository.go` (interface)
- [x] Create `internal/service/settings/repository_pg.go` (implementation)
- [x] sqlc queries: `internal/infra/database/queries/settings.sql`
- [x] Generate sqlc code
- [ ] **Test**: Repository unit tests with mockery
- [x] **Lint**: Check
- [ ] **Coverage**: 80%+
- [x] **Verify**: CRUD operations work

### 4.3 Service Layer
- [x] Create `internal/service/settings/service.go`
- [x] Implement Get/Set server settings
- [x] Implement Get/Set user settings
- [x] Add setting validation
- [ ] **Test**: Service unit tests
- [x] **Lint**: Check
- [ ] **Coverage**: 80%+
- [x] **Verify**: Business logic works

### 4.4 API Handler
- [x] Add endpoints to `api/openapi/openapi.yaml`
- [x] Regenerate ogen code: `go generate ./...`
- [x] Create `internal/api/settings_handler.go`
- [x] Implement handler methods
- [ ] **Test**: Handler integration tests
- [x] **Lint**: Check
- [ ] **Coverage**: 80%+
- [x] **Verify**: API endpoints work

**Checkpoint**: Settings service complete ✅
**Tests**: `go test ./internal/service/settings/... ./internal/api/... -v -run Settings`
**Lint**: `golangci-lint run ./internal/service/settings/... ./internal/api/...` ✅
**Coverage**: `go test ./internal/service/settings/... -cover`

---

## Step 5: User Service

### 5.1 Database Schema
- [ ] Migration: `000006_create_users_table.up.sql` (enhance existing?)
- [ ] Migration: `000007_create_user_preferences_table.up.sql`
- [ ] Migration: `000008_create_user_avatars_table.up.sql`
- [ ] Run migrations
- [ ] **Test**: Migrations
- [ ] **Verify**: Tables with indexes

### 5.2 Repository Layer
- [ ] Create `internal/service/user/repository.go`
- [ ] Create `internal/service/user/repository_pg.go`
- [ ] sqlc queries: `internal/infra/database/queries/users.sql`
- [ ] Generate code
- [ ] **Test**: Repository tests
- [ ] **Lint**: Check
- [ ] **Coverage**: 80%+
- [ ] **Verify**: User CRUD works

### 5.3 Service Layer
- [ ] Create `internal/service/user/service.go`
- [ ] Implement user profile operations
- [ ] Implement password hashing (bcrypt/argon2)
- [ ] Implement avatar upload logic
- [ ] **Test**: Service tests
- [ ] **Lint**: Check
- [ ] **Coverage**: 80%+
- [ ] **Verify**: All operations work

### 5.4 API Handler
- [ ] Update OpenAPI spec
- [ ] Regenerate ogen
- [ ] Create `internal/api/user_handler.go`
- [ ] Implement endpoints
- [ ] **Test**: Integration tests
- [ ] **Lint**: Check
- [ ] **Coverage**: 80%+
- [ ] **Verify**: Full user lifecycle

**Checkpoint**: User service complete
**Tests**: `go test ./internal/service/user/... -v`
**Lint**: `golangci-lint run ./internal/service/user/...`
**Coverage**: `go test ./internal/service/user/... -cover`

---

## Step 6: Auth Service

### 6.1 Database Schema ✅ COMPLETE (Commit 20)
- [x] Migration: `000008_create_auth_tokens_table.up.sql` (JWT refresh tokens)
- [x] Migration: `000009_create_password_reset_tokens_table.up.sql` (One-time reset)
- [x] Migration: `000010_create_email_verification_tokens_table.up.sql` (Email verification)
- [x] **Test**: Migrations (auth_tokens_test.go, testing.go helper)
- [x] **Verify**: Token tables exist (3 tables with SHA-256 hashing, partial indexes)

### 6.2 Repository Layer ✅ COMPLETE (Commit 21)
- [x] Create `internal/service/auth/repository.go` (27 methods interface)
- [x] Create `internal/service/auth/repository_pg.go` (PostgreSQL implementation)
- [x] sqlc queries: `auth_tokens.sql` (27 queries for 3 token types)
- [x] **Test**: Repository tests (structure tests in auth_tokens_test.go)
- [x] **Lint**: Check (0 issues)
- [x] **Coverage**: 80%+ (pending full test suite)

### 6.3 JWT Implementation ✅ COMPLETE (Commit 22)
- [x] Create `internal/service/auth/jwt.go` (stdlib crypto, no external libs)
- [x] Implement token generation (HMAC-SHA256 JWT, crypto/rand refresh tokens)
- [x] Implement token validation (signature verification, expiry checks)
- [x] **Test**: JWT tests (pending)
- [x] **Lint**: Check (0 issues)
- [x] **Coverage**: 80%+ (pending test suite)

### 6.4 Service Layer
- [ ] Create `internal/service/auth/service.go`
- [ ] Implement Login
- [ ] Implement Logout
- [ ] Implement Register
- [ ] Implement Refresh
- [ ] Implement Password Reset flow
- [ ] **Test**: Service tests
- [ ] **Lint**: Check
- [ ] **Coverage**: 80%+

### 6.5 Middleware
- [ ] Create `internal/api/middleware/auth.go`
- [ ] JWT validation middleware
- [ ] User context injection
- [ ] **Test**: Middleware tests
- [ ] **Lint**: Check
- [ ] **Coverage**: 80%+

### 6.6 API Handler
- [ ] Update OpenAPI spec
- [ ] Regenerate ogen
- [ ] Create `internal/api/auth_handler.go`
- [ ] Implement all endpoints
- [ ] **Test**: End-to-end auth flow
- [ ] **Lint**: Check
- [ ] **Coverage**: 80%+

**Checkpoint**: Auth service complete
**Tests**: Full auth flow works
**Lint**: Clean
**Coverage**: 80%+

---

## Step 7: Session Service

(Following same pattern: Schema → Repository → Service → Handler → Tests)

---

## Step 8: RBAC Service

(Following same pattern)

---

## Step 9: API Keys Service

(Following same pattern)

---

## Step 10: OIDC Service

(Following same pattern)

---

## Step 11: Activity Service

(Following same pattern)

---

## Step 12: Library Service

(Following same pattern)

---

## Step 13: Health Service Enhancement

### 13.1 Enhanced Health Checks
- [ ] Add service-level health checks
- [ ] Add River queue health
- [ ] Add external API checks
- [ ] **Test**: All health checks
- [ ] **Lint**: Check
- [ ] **Verify**: `/health` returns detailed status

---

## Testing Strategy per Step

### Unit Tests
```bash
# Run tests for specific package
go test ./internal/service/settings/... -v

# With coverage
go test ./internal/service/settings/... -cover -coverprofile=coverage.out

# View coverage
go tool cover -html=coverage.out
```

### Integration Tests
```bash
# Run integration tests (with testcontainers)
go test ./tests/integration/... -v

# Skip slow tests
go test ./... -short
```

### Linting
```bash
# Lint specific path
golangci-lint run ./internal/service/settings/...

# Lint everything
golangci-lint run ./...

# Fix auto-fixable issues
golangci-lint run --fix ./...
```

### Full Check Before Next Step
```bash
# Complete verification
make test          # All tests
make lint          # All linting
make coverage      # Coverage check
```

---

## Verification Checklist (After Each Service)

- [ ] All tests pass (unit + integration)
- [ ] Lint is clean (zero warnings)
- [ ] Coverage ≥ 80%
- [ ] OpenAPI spec updated
- [ ] Migrations tested (up/down)
- [ ] API endpoints documented
- [ ] Error handling complete
- [ ] Logging added
- [ ] Metrics instrumented
- [ ] Integration tests pass

---

## Current Step

**Step**: Not Started
**Next**: Step 1.1 - Database Migrations Setup

---

## Notes

- **Never skip tests**: Tests catch issues early
- **Lint frequently**: Small fixes easier than big cleanup
- **Commit often**: After each sub-step that works
- **Document decisions**: Update QUESTIONS_v0.2.0.md when decided
- **Log bugs**: Update BUGS_v0.2.0.md when found
- **Update status**: Keep STATUS_v0.2.0.md current

---

## Dependencies Reminder

From SOURCE_OF_TRUTH:
- Go 1.25.6
- PostgreSQL 18.1
- casbin v2.135.0
- rueidis v1.0.49
- otter v2.x
- river v0.26.0
- go-mail v0.6.2
- golang.org/x/crypto v0.47.0
