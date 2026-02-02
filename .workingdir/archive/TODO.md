# v0.1.0 Implementation TODO

**Updated**: 2026-02-02 12:06
**Milestone**: v0.1.0 - Skeleton
**Priority Order**: P0 → P1 → P2

---

## ✅ COMPLETED: Integration Testing Infrastructure (2026-02-02)

**Goal**: Establish comprehensive integration testing with testcontainers-go
**Status**: COMPLETE (4 phases, 3.7h actual vs 6-9h estimated)
**Commits**:
- 25cedcafb9: Phase 1 - testcontainers PostgreSQL (testutil/containers.go, 5 tests)
- 073c4c93c9: Phase 2 - Health endpoint tests (10 tests)
- ca780a6648: Phase 2 - Database integration tests (12 tests)
- 20c381b778: Phase 3 - ogen client tests (10 tests)
- 3787369a2b: Phase 4 - CI workflow configuration (.github/workflows/dev.yml)
- 990ab8e25d: Fix - Remove ALTER DATABASE from migrations (CI failure fix)

**Deliverables**:
- ✅ 40 integration tests across 5 files (tests/integration/*.go)
- ✅ testcontainers-go infrastructure (internal/testutil/containers.go)
- ✅ CI pipeline configured with testcontainers support
- ✅ Real PostgreSQL 18.1-alpine in containers (no mocks)
- ✅ Full server lifecycle testing (startup, graceful shutdown, signals)
- ✅ Health endpoint E2E tests (liveness, readiness, startup)
- ✅ Database integration tests (migrations, transactions, constraints)
- ✅ ogen client type-safe API tests (contract validation)

**Migration Fix (990ab8e25d)**:
- Fixed CI failure: "database revenge does not exist"
- Removed `ALTER DATABASE revenge SET search_path` from 000001_create_schemas.up.sql
- Added `?search_path=public,shared` to testcontainers connection URL
- Schema isolation now via connection parameter (testable + production-ready)
- CI run 21588234982 verifying fix (in progress)

**Next Steps**:
- Monitor CI run 21588234982 to verify all 40 tests pass
- Update .workingdir/TODO.md with Phase 1 HTTP Server progress
- Document K8s/Helm/Docker Swarm deployment testing scope for v0.1.0
- Ensure documentation updates use template-based pipelines (scripts/automation/)

---

## 🔴 Phase 1: HTTP Server Implementation (P0 - CRITICAL)

**Goal**: Get HTTP server running with health endpoints
**Estimated**: 4-6 hours
**Start Date**: TBD
**Target**: Today

### Step 1.1: Generate ogen Code ⏳
```bash
# Commands to run:
make generate
# or
go generate ./...
```

**Tasks**:
- [ ] Verify `ogen.yaml` configuration is correct
- [ ] Run code generation
- [ ] Verify generated files in `internal/api/oas_*.go`
- [ ] Review generated interfaces
- [ ] Commit generated code with message: "chore: generate ogen server code"

**Expected Output**:
- `internal/api/oas_server_gen.go`
- `internal/api/oas_types_gen.go`
- `internal/api/oas_client_gen.go` (if enabled)
- Other ogen-generated files

**Verification**:
```bash
ls -la internal/api/oas_*.go
```

---

### Step 1.2: Implement HTTP Server ⏳

**File**: `internal/api/server.go`

**Tasks**:
- [ ] Create file with Server struct
- [ ] Implement ogen.Handler interface
- [ ] Add HTTP server configuration (port, timeouts)
- [ ] Add lifecycle hooks (Start, Stop)
- [ ] Wire into fx.App
- [ ] Add graceful shutdown handling

**Template**:
```go
package api

import (
    "context"
    "net/http"
    "time"

    "go.uber.org/fx"
    "go.uber.org/zap"
)

type Server struct {
    logger *zap.Logger
    health HealthService // from infra/health
    server *http.Server
}

func NewServer(
    lc fx.Lifecycle,
    logger *zap.Logger,
    health HealthService,
    config *config.Config,
) (*Server, error) {
    // Implementation
}

func (s *Server) Start(ctx context.Context) error {
    // Start HTTP server
}

func (s *Server) Stop(ctx context.Context) error {
    // Graceful shutdown
}
```

---

### Step 1.3: Implement Health Handlers ⏳

**File**: `internal/api/health_handler.go`

**Tasks**:
- [ ] Create handler struct
- [ ] Implement GetLiveness (ogen method)
- [ ] Implement GetReadiness (ogen method)
- [ ] Implement GetStartup (ogen method)
- [ ] Wire health service from infra/health
- [ ] Add error handling
- [ ] Add response formatting

**Template**:
```go
package api

import (
    "context"

    "github.com/lusoris/revenge/internal/infra/health"
)

type HealthHandler struct {
    healthService *health.Service
}

func NewHealthHandler(hs *health.Service) *HealthHandler {
    return &HealthHandler{healthService: hs}
}

// GetLiveness implements ogen GetLiveness method
func (h *HealthHandler) GetLiveness(ctx context.Context) (*HealthCheck, error) {
    // Implementation
}

// GetReadiness implements ogen GetReadiness method
func (h *HealthHandler) GetReadiness(ctx context.Context) (*HealthCheck, error) {
    // Implementation
}

// GetStartup implements ogen GetStartup method
func (h *HealthHandler) GetStartup(ctx context.Context) (*HealthCheck, error) {
    // Implementation
}
```

---

### Step 1.4: Wire into fx.App ⏳

**File**: `internal/api/module.go`

**Tasks**:
- [ ] Create fx.Module for API
- [ ] Provide Server
- [ ] Provide HealthHandler
- [ ] Register lifecycle hooks
- [ ] Add to main.go imports

**Template**:
```go
package api

import "go.uber.org/fx"

var Module = fx.Module("api",
    fx.Provide(
        NewServer,
        NewHealthHandler,
    ),
)
```

**Update**: `cmd/revenge/main.go`
```go
fx.New(
    // ... existing modules
    api.Module,
)
```

---

### Step 1.5: Test Locally ⏳

**Tasks**:
- [ ] Start server: `make dev` or `go run cmd/revenge/main.go`
- [ ] Verify server starts without errors
- [ ] Test liveness: `curl http://localhost:8080/health/live`
- [ ] Test readiness: `curl http://localhost:8080/health/ready`
- [ ] Test startup: `curl http://localhost:8080/health/startup`
- [ ] Verify JSON responses
- [ ] Check logs for errors
- [ ] Test graceful shutdown (Ctrl+C)

**Expected Responses**:
```json
// GET /health/live
{
  "status": "ok",
  "timestamp": "2026-02-02T11:35:00Z"
}

// GET /health/ready
{
  "status": "ok",
  "checks": {
    "database": "ok",
    "cache": "ok",
    "search": "ok"
  },
  "timestamp": "2026-02-02T11:35:00Z"
}
```

---

### Step 1.6: Add Tests ⏳

**File**: `internal/api/server_test.go`

**Tasks**:
- [ ] Test server startup
- [ ] Test graceful shutdown
- [ ] Test configuration loading
- [ ] Test lifecycle hooks

**File**: `internal/api/health_handler_test.go`

**Tasks**:
- [ ] Test GetLiveness handler
- [ ] Test GetReadiness handler
- [ ] Test GetStartup handler
- [ ] Test error responses
- [ ] Test JSON serialization

**Run Tests**:
```bash
go test ./internal/api/... -v
```

---

### Step 1.7: Commit & Push ⏳

**Tasks**:
- [ ] Stage all changes: `git add internal/api/`
- [ ] Commit with message:
```
feat: implement HTTP server with health endpoints

- Generate ogen server code from OpenAPI spec
- Implement HTTP server with lifecycle management
- Implement health check endpoints (live, ready, startup)
- Wire server into fx.App
- Add tests for server and handlers

Closes BUG-001
Refs: TODO_v0.1.0.md Phase 1
```
- [ ] Push to develop: `git push origin develop`
- [ ] Verify CI passes

**Verification**:
```bash
gh run list --branch develop --limit 1
```

---

## 🟡 Phase 2: Improve Test Coverage (P1 - HIGH)

**Goal**: Increase test coverage from 33% to 80%+
**Estimated**: 8-12 hours
**Dependencies**: Phase 1 complete
**Target**: This week

### 2.1: Config Package Tests (9.5% → 80%) ⏳

**File**: `internal/config/config_test.go` (enhance)

**Tasks**:
- [ ] Test default config loading
- [ ] Test YAML file loading from different paths
- [ ] Test environment variable overrides (REVENGE_*)
- [ ] Test config validation failures
- [ ] Test config merging (defaults + file + env)
- [ ] Test invalid YAML handling
- [ ] Test missing file handling

**Run**:
```bash
go test ./internal/config/... -v -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

### 2.2: Database Package Tests (11% → 80%) ⏳

**File**: `internal/infra/database/pool_test.go` (enhance)

**Tasks**:
- [ ] Test pgxpool connection with valid config
- [ ] Test pgxpool connection with invalid config
- [ ] Test health check (ping)
- [ ] Test connection pool configuration
- [ ] Test graceful shutdown
- [ ] Test connection retry logic
- [ ] Use embedded-postgres for real database tests

**File**: `internal/infra/database/migrate_test.go` (create)

**Tasks**:
- [ ] Test migration up
- [ ] Test migration down
- [ ] Test migration version
- [ ] Test migration with invalid SQL
- [ ] Test migration rollback

---

### 2.3: Health Package Tests (34% → 80%) ⏳

**File**: `internal/infra/health/service_test.go` (enhance)

**Tasks**:
- [ ] Test liveness check (always healthy)
- [ ] Test readiness check (all dependencies healthy)
- [ ] Test readiness check (one dependency down)
- [ ] Test startup check
- [ ] Test check timeout handling
- [ ] Test check error handling
- [ ] Mock dependency checks

**File**: `internal/infra/health/checks_test.go` (create)

**Tasks**:
- [ ] Test PostgreSQL ping check
- [ ] Test Dragonfly ping check
- [ ] Test Typesense health check
- [ ] Test River worker check
- [ ] Test check failures

---

### 2.4: App Package Tests (27.8% → 80%) ⏳

**File**: `internal/app/module_test.go` (enhance)

**Tasks**:
- [ ] Test fx.App startup with all modules
- [ ] Test fx.App graceful shutdown
- [ ] Test module registration
- [ ] Test lifecycle hooks execution order
- [ ] Test dependency injection
- [ ] Test missing dependency errors

---

### 2.5: Cache Package Tests (0% → 80%) ⏳

**File**: `internal/infra/cache/module_test.go` (enhance)

**Tasks**:
- [ ] Test rueidis client creation
- [ ] Test connection to Dragonfly/Redis
- [ ] Test connection errors
- [ ] Test health check
- [ ] Test graceful shutdown
- [ ] Mock rueidis for unit tests

---

### 2.6: Jobs Package Tests (0% → 80%) ⏳

**File**: `internal/infra/jobs/module_test.go` (enhance)

**Tasks**:
- [ ] Test river client creation
- [ ] Test worker registration
- [ ] Test job queue operations
- [ ] Test connection errors
- [ ] Test graceful shutdown
- [ ] Mock river for unit tests

---

### 2.7: Search Package Tests (0% → 80%) ⏳

**File**: `internal/infra/search/module_test.go` (enhance)

**Tasks**:
- [ ] Test typesense client creation
- [ ] Test connection to Typesense
- [ ] Test health check
- [ ] Test connection errors
- [ ] Test graceful shutdown
- [ ] Mock typesense for unit tests

---

### 2.8: Errors Package Tests (44.4% → 80%) ⏳

**File**: `internal/errors/wrap_test.go` (create)

**Tasks**:
- [ ] Test error wrapping with context
- [ ] Test error unwrapping
- [ ] Test stack trace preservation
- [ ] Test error formatting
- [ ] Test multiple levels of wrapping

---

## 🟡 Phase 3: Integration Testing (P1 - HIGH)

**Goal**: Add full-stack integration tests with real services
**Estimated**: 4-6 hours
**Dependencies**: Phase 2 mostly complete
**Target**: This week

### 3.1: testcontainers Setup ⏳

**File**: `internal/testutil/integration_test.go` (create)

**Tasks**:
- [ ] Set up testcontainers PostgreSQL
- [ ] Set up testcontainers Dragonfly
- [ ] Set up testcontainers Typesense
- [ ] Test service startup
- [ ] Test service connectivity

---

### 3.2: Database Integration Tests ⏳

**Tasks**:
- [ ] Test migrations with real PostgreSQL
- [ ] Test connection pooling
- [ ] Test concurrent queries
- [ ] Test transaction handling

---

### 3.3: HTTP Server Integration Tests ⏳

**Tasks**:
- [ ] Start full fx.App with all services
- [ ] Test HTTP server startup
- [ ] Test health endpoints with real dependencies
- [ ] Test graceful shutdown
- [ ] Test error handling

---

### 3.4: End-to-End Test ⏳

**File**: `internal/api/e2e_test.go` (create)

**Tasks**:
- [ ] Start all services (PostgreSQL, Dragonfly, Typesense)
- [ ] Start HTTP server
- [ ] Run migrations
- [ ] Test full health check flow
- [ ] Test dependency failures
- [ ] Clean up services

---

## ⚪ Phase 4: Polish & Documentation (P2 - MEDIUM)

**Goal**: Finish documentation and prepare for release
**Estimated**: 2-3 hours
**Dependencies**: All phases mostly complete
**Target**: Before v0.1.0 release

### 4.1: Configuration Documentation ⏳

**File**: `config/config.example.yaml` (create)

**Tasks**:
- [ ] Copy config.yaml structure
- [ ] Add comments for every option
- [ ] Add example values
- [ ] Document environment variable overrides
- [ ] Add recommended production values

---

### 4.2: Go Toolchain Fix ⏳

**Tasks**:
- [ ] Investigate version mismatch (1.25.6 vs 1.25.5)
- [ ] Update all Go toolchain components to 1.25.6
- [ ] Verify `go test -cover` works
- [ ] Update CI if needed

---

### 4.3: README Updates ⏳

**File**: `README.md`

**Tasks**:
- [ ] Add v0.1.0 status badge
- [ ] Update feature list
- [ ] Add "Getting Started" section
- [ ] Document prerequisites
- [ ] Add build instructions
- [ ] Add test instructions
- [ ] Document available endpoints

---

### 4.4: Final Verification ⏳

**Tasks**:
- [ ] Run full test suite: `make test`
- [ ] Verify coverage ≥80%: `make test-coverage`
- [ ] Run linters: `make lint`
- [ ] Build binary: `make build`
- [ ] Test binary: `./bin/revenge`
- [ ] Test migrations: `make migrate-up && make migrate-down`
- [ ] Start Docker Compose: `make docker-up`
- [ ] Verify all services healthy
- [ ] Stop Docker Compose: `make docker-down`
- [ ] Verify CI passes

---

### 4.5: Release Preparation ⏳

**Tasks**:
- [ ] Update CHANGELOG.md with v0.1.0 changes
- [ ] Create Git tag: `git tag v0.1.0`
- [ ] Push tag: `git push origin v0.1.0`
- [ ] Create GitHub release
- [ ] Update project board
- [ ] Close v0.1.0 milestone

---

## Progress Tracking

### Daily Checklist

**Today (2026-02-02)**:
- [x] Complete gap analysis
- [x] Create tracking files
- [ ] Start Phase 1: Generate ogen code
- [ ] Implement HTTP server
- [ ] Test health endpoints
- [ ] Commit Phase 1

**This Week**:
- [ ] Complete Phase 1
- [ ] Complete Phase 2 (test coverage)
- [ ] Complete Phase 3 (integration tests)
- [ ] Start Phase 4 (polish)

**Before v0.1.0**:
- [ ] All phases complete
- [ ] All tests passing
- [ ] Documentation updated
- [ ] CI passing
- [ ] Ready for feature testing

---

## Notes

- Work in order: P0 → P1 → P2
- Commit after each major task
- Update STATUS.md daily
- Update BUGS.md when issues found
- Ask questions in QUESTIONS.md
- Target: v0.1.0 complete this week
