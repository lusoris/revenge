# Pipeline "Reduced Mode" Analysis

**Date**: 2026-02-01
**Question**: Are pipelines in "reduced mode" because they were created before code existed?

---

## Findings

### ✅ FULLY ACTIVE: Regular Unit Tests

**Status**: Running normally, NO reduced mode

**Evidence:**
- Go test files exist: 4 test files in `internal/`
- Test job runs: `go test -v -race -coverprofile=coverage.out -covermode=atomic ./...`
- No skip conditions

**Test Files:**
- `internal/version/version_test.go`
- `internal/config/config_test.go`
- `internal/config/config_bugfix_test.go`
- `internal/infra/database/pool_bugfix_test.go`

**Workflow**: `.github/workflows/dev.yml` lines 66-113 (test job)

---

### ⚠️ REDUCED MODE: Integration Tests

**Status**: SKIPPED if no integration tests exist

**Skip Condition** (.github/workflows/dev.yml:285-290):
```bash
# Skip if no Go integration tests exist yet
if go list -tags=integration ./tests/integration/... 2>/dev/null | grep -q .; then
  go test -v -tags=integration ./tests/integration/...
else
  echo "No Go integration tests found, skipping"
fi
```

**Current State:**
- `tests/integration/` exists
- Contains only Python tests (test_doc_pipeline.py)
- NO Go files with `//go:build integration` tag
- Integration tests ARE currently being skipped

**When to Remove Skip:**
- When we create actual Go integration tests
- Files should be in `tests/integration/` with `//go:build integration` tag
- Should use testcontainers for PostgreSQL, Redis, Typesense

---

### ✅ FULLY ACTIVE: Linting

**Status**: Running normally, NO reduced mode

**Evidence:**
- golangci-lint runs on all Go code
- No skip conditions
- Recently fixed for v2.8.0 (ISSUE-006)

---

### ✅ FULLY ACTIVE: Build Matrix

**Status**: Running normally, NO reduced mode

**Evidence:**
- Builds for 5 platforms: linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64
- All builds running
- No skip conditions

---

## Recommendation

**KEEP the integration test skip condition for now**

**Reasoning:**
1. It's a sensible guard - tests can't run if they don't exist
2. Regular unit tests ARE running
3. When we add integration tests, the skip will automatically stop

**Future Action Item:**
Create integration tests in `tests/integration/` when we have actual services to integration test (database, cache, search, jobs).

**Priority**: LOW - The important tests (unit tests, linting) are all running

---

## Other Workflows Checked

**pr-checks.yml**: Has skip for "exempt" label - intentional feature, not reduced mode

---

## Conclusion

Only integration tests are in "reduced mode" (intentionally skipped). All other pipeline stages are fully active and testing real code.
