# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- N/A

### Changed
- N/A

### Fixed
- N/A

## [0.1.2] - 2026-02-02

### Added
- **Errors Package Tests**: Complete test coverage for error handling utilities
  - `internal/errors`: 44% → 100% coverage

### Test Infrastructure
- **errors/wrap_test.go**: Comprehensive tests for wrap.go functions
  - `TestWrapf`: 6 subtests covering format args, nil handling, nesting
  - `TestWithStack`: 5 subtests covering stack trace verification
  - `TestWrapSentinel`: 8 subtests covering all sentinel errors
  - `TestFormatError`: 7 subtests covering formatting scenarios
  - `TestWrapChaining`: chaining Wrap, Wrapf, WithStack
  - `TestConcurrentErrorCreation`: concurrent safety (100 goroutines)
  - `TestErrorMessageFormat`: message format verification

## [0.1.1] - 2026-02-02

### Added
- **Comprehensive Unit Tests**: Major test coverage improvements across core packages
  - `internal/api`: 41% → 97.1% coverage with fx.Lifecycle integration tests
  - `internal/config`: 10% → 76.2% coverage with loader and validation tests
  - `internal/infra/health`: 55% → 68.1% coverage with health check tests
  - `internal/infra/database`: 20% → 22% coverage with PoolConfig tests

### Test Infrastructure
- **api/handler_test.go**: Tests for GetLiveness, GetStartup, GetReadiness, NewError
  - Embedded-postgres integration for realistic database testing
  - Concurrent request handling tests
  - Error response validation
- **api/server_test.go**: Full fx.Lifecycle tests using fxtest.New
  - Server startup/shutdown lifecycle
  - Graceful shutdown verification
  - Concurrent request handling (50 parallel requests)
  - Configuration application tests
  - Multiple port sequence tests
- **config/loader_test.go**: Comprehensive configuration tests
  - Default value loading
  - YAML file loading
  - Environment variable overrides (REVENGE_* prefix)
  - Validation failure scenarios
  - MustLoad panic handling
- **database/pool_test.go**: Connection pool configuration tests
  - MaxConns calculation (CPU * 2 + 1)
  - URL parsing and validation
  - Connection timeout settings
- **health/checks_test.go**: Stub health check tests
  - CheckCache, CheckSearch, CheckJobs
  - Status constants validation
  - Concurrent check execution

### Notes
- Tests use embedded-postgres for integration testing
- Run with `-p 1` flag to avoid port conflicts in parallel mode
- Full test suite: `go test ./internal/... -cover -count=1 -p 1`

## [0.1.0] - 2026-02-02

### Added
- **HTTP Server**: Fully functional HTTP server with ogen-generated code from OpenAPI spec
- **Health Endpoints**:
  - `GET /health/live` - Liveness probe (Kubernetes)
  - `GET /health/ready` - Readiness probe with dependency checks
  - `GET /health/startup` - Startup probe
- **Configuration System**: YAML-based configuration with environment variable support
- **Structured Logging**: Dual logging with Zap (JSON) and Slog (structured)
- **Database Support**: PostgreSQL connection pooling with pgxpool and migrations
- **Dependency Injection**: Uber fx for lifecycle management
- **Docker Support**: Multi-stage Dockerfile and docker-compose for development
- **CI/CD Pipeline**: GitHub Actions with build, test, lint, and security scanning
- **Integration Tests**: Comprehensive test suite with testcontainers
- **OpenAPI Spec**: Full API specification at `api/openapi/openapi.yaml`

### Infrastructure
- PostgreSQL 18.1 support
- Dragonfly (Redis-compatible) cache client stub
- Typesense search client stub
- River background jobs client stub

### Developer Experience
- Makefile with common targets (`make build`, `make test`, `make lint`)
- GoReleaser configuration for releases
- Renovate for dependency updates
- CodeQL security scanning
