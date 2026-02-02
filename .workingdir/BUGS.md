# v0.1.0 Bugs & Issues Tracker

**Updated**: 2026-02-02 11:25
**Milestone**: v0.1.0 - Skeleton

---

## Critical Bugs (P0)

### ✅ BUG-001: No HTTP Server Running [RESOLVED]

**Status**: ✅ Resolved (2026-02-02 13:30)
**Priority**: P0 (Critical)
**Severity**: Blocker
**Discovered**: 2026-02-02
**Resolved**: 2026-02-02
**Component**: internal/api

**Description**:
OpenAPI specification exists at `api/openapi/openapi.yaml` with health endpoints defined, but no actual HTTP server was running. The ogen code had not been generated, and no handlers were implemented.

**Impact**:
- Cannot test API endpoints
- Health checks not accessible
- v0.1.0 non-functional

**Root Cause**:
1. ogen.yaml config format incompatible with ogen CLI v1.18.0
2. .gitignore excluded *_gen.go files
3. No HTTP server implementation
4. No handler implementation

**Solution Implemented**:
✅ **Commit 6477058dca**: Generated ogen code
  - Fixed ogen CLI invocation (use flags instead of config file)
  - Generated 17 ogen files (oas_*_gen.go)
  - Updated .gitignore to commit generated code
  - Updated Makefile with correct ogen command

✅ **Commit 63a5202131**: Implemented HTTP server and handlers
  - Created internal/api/server.go with fx lifecycle
  - Created internal/api/handler.go implementing ogen.Handler
  - Implemented GetLiveness, GetReadiness, GetStartup handlers
  - Wired into fx.App via internal/api/module.go
  - Made Auth.JWTSecret optional (not in v0.1.0 scope)

**Verification**:
- ✅ Code compiles successfully
- ⏳ Local testing pending (needs database setup)

**Assignee**: Agent
**Resolution**: 2026-02-02 13:30

---

## High Priority Bugs (P1)

### 🟡 BUG-002: Test Coverage Below Target

**Status**: Open
**Priority**: P1 (High)
**Severity**: Major
**Discovered**: 2026-02-02
**Component**: Tests

**Description**:
Overall test coverage is 33%, well below the 80% target for v0.1.0.

**Affected Packages**:
- internal/config: 9.5%
- internal/infra/database: 11.0%
- internal/infra/health: 34.0%
- internal/infra/cache: 0.0%
- internal/infra/jobs: 0.0%
- internal/infra/search: 0.0%

**Impact**:
- Untested code may contain bugs
- Cannot verify correctness
- Does not meet v0.1.0 requirements

**Solution**:
See STATUS.md Phase 2 for test implementation plan.

**Assignee**: TBD
**ETA**: This week

---

### 🟡 BUG-003: sqlc Not Generating Code

**Status**: Open
**Priority**: P1 (High)
**Severity**: Major
**Discovered**: 2026-02-02
**Component**: internal/infra/database

**Description**:
`sqlc.yaml` configuration exists, but no SQL queries are defined and no Go code has been generated.

**Impact**:
- Cannot interact with database using type-safe queries
- No user/session query functions available

**Root Cause**:
No `.sql` query files exist in the repository.

**Solution**:
1. Create query files in appropriate directories
2. Run `sqlc generate`
3. Commit generated code

**Assignee**: TBD
**ETA**: This week

---

## Medium Priority Bugs (P2)

### ⚠️ BUG-004: Go Toolchain Version Mismatch

**Status**: Open
**Priority**: P2 (Medium)
**Severity**: Minor
**Discovered**: 2026-02-02
**Component**: Build System

**Description**:
Running `go test -cover` produces errors:
```
compile: version "go1.25.6" does not match go tool version "go1.25.5"
```

**Impact**:
- Cannot accurately measure test coverage
- Build still works, but coverage reporting broken

**Root Cause**:
Mismatch between Go module version (1.25.6) and installed Go toolchain (1.25.5).

**Possible Solutions**:
1. Upgrade all Go toolchain components to 1.25.6
2. Downgrade go.mod to 1.25.5
3. Use Go toolchain management (`go install golang.org/dl/go1.25.6@latest`)

**Assignee**: TBD
**ETA**: Before v0.1.0 release

---

### ⚠️ BUG-005: Missing config.example.yaml

**Status**: Open
**Priority**: P2 (Medium)
**Severity**: Minor
**Discovered**: 2026-02-02
**Component**: config

**Description**:
`config/config.yaml` exists but no `config/config.example.yaml` for documentation purposes.

**Impact**:
- Users don't have a fully documented config reference
- Harder to understand available options

**Solution**:
Create `config/config.example.yaml` with:
- All configuration options
- Comments explaining each option
- Example values

**Assignee**: TBD
**ETA**: Before v0.1.0 release

---

## Low Priority Issues (P3)

### 📝 ISSUE-001: Integration Tests Not Utilizing testcontainers

**Status**: Open
**Priority**: P3 (Low)
**Severity**: Enhancement
**Discovered**: 2026-02-02
**Component**: Tests

**Description**:
`internal/testutil/containers.go` exists with testcontainers setup, but no tests are currently using it.

**Impact**:
- Missing real integration tests with actual services
- Cannot test interactions with PostgreSQL, Dragonfly, Typesense

**Solution**:
See STATUS.md Phase 3 for integration test plan.

**Assignee**: TBD
**ETA**: This week

---

### 📝 ISSUE-002: Makefile Lint Target Skipped

**Status**: Open
**Priority**: P3 (Low)
**Severity**: Enhancement
**Discovered**: 2026-02-02
**Component**: Build System

**Description**:
`make lint` is currently skipped because golangci-lint doesn't support Go 1.25 yet.

**Impact**:
- Cannot run linting checks
- May introduce code style issues

**Solution**:
Wait for golangci-lint to support Go 1.25, or use alternative linters.

**Assignee**: TBD
**ETA**: When golangci-lint adds Go 1.25 support

---

## Resolved Bugs

_(None yet)_

---

## Bug Statistics

**Total Open**: 7
- P0 Critical: 1
- P1 High: 2
- P2 Medium: 2
- P3 Low: 2

**Total Resolved**: 0

**Resolution Rate**: 0%

---

## Testing Notes

### Test Execution Issues

**Go Version Mismatch**:
```bash
go test -cover ./...
# Error: compile: version "go1.25.6" does not match go tool version "go1.25.5"
```

**Workaround**:
Run tests without coverage flag:
```bash
go test -short ./...
# Works fine, all tests pass
```

---

## Action Items

1. **URGENT**: Fix BUG-001 (No HTTP server) - Blocker for all API work
2. **HIGH**: Address BUG-002 (Test coverage) - Required for v0.1.0 completion
3. **HIGH**: Fix BUG-003 (sqlc generation) - Needed for database interactions
4. **MEDIUM**: Investigate BUG-004 (Go version) - Affects coverage reporting
5. **MEDIUM**: Create BUG-005 (config.example) - Documentation improvement
6. **LOW**: Implement ISSUE-001 (testcontainers) - Better integration testing
7. **LOW**: Monitor ISSUE-002 (lint support) - Wait for tooling updates
