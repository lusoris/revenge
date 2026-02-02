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
