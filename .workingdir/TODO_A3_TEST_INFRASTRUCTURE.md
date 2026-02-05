# Phase A3: Test Infrastructure

**Priority**: P2
**Effort**: 3-4h
**Dependencies**: A0
**Status**: ✅ Complete (2026-02-05)

---

## A3.0: PostgreSQL Test Database (NEW) ✅

**Affected Files**:
- `internal/testutil/pgtestdb.go` (NEW)
- `internal/testutil/testdb.go` (updated)
- All `*_test.go` files using database

**Completed Tasks**:
- [x] Implement `NewFastTestDB(t) *FastTestDB` using testcontainers
- [x] Use `postgres:17-alpine` container image
- [x] Implement template database pattern for fast cloning (~50ms per test)
- [x] Add `DB` interface for both TestDB and FastTestDB
- [x] Migrate all service tests from embedded-postgres to testcontainers
- [x] Fix `replaceDBName()` bug that caused parallel test failures
- [x] Automatic container cleanup after tests

**Implementation Details**:
- Shared PostgreSQL container across all tests (started once, ~7s)
- Template database with all migrations applied
- Each test gets isolated database via `CREATE DATABASE ... TEMPLATE`
- Parallel-safe: each test has unique database name (test_1, test_2, ...)
- Automatic cleanup: test databases dropped, container terminated

**Test Results**:
- 32/35 packages passing
- Remaining failures are legacy tests (api, database, health) with embedded-postgres conflicts
- All service tests now use testcontainers infrastructure

---

## A3.1: Dragonfly Testcontainer ✅

**Affected File**: `internal/testutil/containers.go`

**Completed Tasks**:
- [x] Implement `NewDragonflyContainer(t) *DragonflyContainer`
- [x] Use `docker.io/dragonflydb/dragonfly:latest` image
- [x] Expose port 6379 (Redis-compatible)
- [x] Wait for "accepting connections" log message
- [x] Return connection URL, Host, Port
- [x] Cleanup via `Close()` method

**Implementation Details**:
- Returns Redis-compatible URL: `redis://host:port`
- Follows same pattern as PostgreSQL container
- 60 second startup timeout

---

## A3.2: Typesense Testcontainer ✅

**Affected File**: `internal/testutil/containers.go`

**Completed Tasks**:
- [x] Implement `NewTypesenseContainer(t) *TypesenseContainer`
- [x] Use `typesense/typesense:27.1` image
- [x] Configure test API key
- [x] Expose port 8108 (HTTP API)
- [x] Wait for `/health` endpoint
- [x] Return connection URL, Host, Port, APIKey
- [x] Cleanup via `Close()` method

**Implementation Details**:
- Uses test API key: `test-api-key-for-integration-tests`
- Returns HTTP URL: `http://host:port`
- Health check via HTTP endpoint
- 60 second startup timeout
