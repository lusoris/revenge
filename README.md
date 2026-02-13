# Revenge

[![CI](https://github.com/lusoris/revenge/actions/workflows/ci.yml/badge.svg?branch=develop)](https://github.com/lusoris/revenge/actions/workflows/ci.yml)
[![Develop Build](https://github.com/lusoris/revenge/actions/workflows/develop.yml/badge.svg?branch=develop)](https://github.com/lusoris/revenge/actions/workflows/develop.yml)
[![Security](https://github.com/lusoris/revenge/actions/workflows/security.yml/badge.svg?branch=develop)](https://github.com/lusoris/revenge/actions/workflows/security.yml)
[![Go Version](https://img.shields.io/badge/go-1.25.7-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-AGPL--3.0-blue)](LICENSE)

> Modular media server with complete content isolation

A ground-up media server built in Go with a fully modular architecture. Each content type is completely isolated with its own tables, services, and handlers.

**WARNING: EARLY DEVELOPMENT** - Core infrastructure and backend services implemented, frontend not yet started. See the [current status](TODO.md) for details.

---

## Design Principles

- **Performance First** - UX never blocked by backend tasks
- **Client Agnostic** - REST API designed for native mobile/TV apps
- **Privacy by Default** - Self-hosted, no telemetry
- **Bleeding Edge Stable** - Latest stable Go/PostgreSQL, no alpha deps

See [DESIGN_PRINCIPLES.md](docs/dev/design/architecture/DESIGN_PRINCIPLES.md) for full details.

---

## Implemented

### Backend Services (15)
Authentication (JWT + refresh tokens), MFA (TOTP + WebAuthn), Sessions, RBAC (Casbin), OIDC/SSO, Users, API Keys, Settings, Activity logging, Library management, Metadata aggregation, Search (Typesense), Email (SMTP), Notifications, Storage (local + S3)

### Content Modules
- **Movies** - Full CRUD, metadata enrichment (TMDb + TVDb), library scanning, file matching, search indexing
- **TV Shows** - Series/season/episode hierarchy, metadata enrichment, library scanning
- **QAR** - Adult content isolation (pirate-themed obfuscation, separate `qar` schema, `/api/v1/legacy/*` namespace)

### Background Jobs (River)
9 workers: metadata refresh, library scan, file match, search index, series refresh, library cleanup, activity cleanup

### Integrations
- **Metadata**: TMDb, TheTVDB (with caching, language support, force refresh)
- **Servarr**: Radarr, Sonarr (library sync)
- **Auth**: Generic OIDC provider

### Infrastructure
- PostgreSQL 18+ with pgxpool, 30+ migrations, sqlc codegen
- Dragonfly (Redis-compatible) distributed cache via rueidis
- otter v2 in-memory cache (W-TinyLFU)
- Typesense full-text search
- River PostgreSQL-native job queue
- K8s health checks (liveness/readiness/startup)
- OpenTelemetry observability
- govips image processing

### CI/CD
8 GitHub Actions workflows: CI, develop auto-build, release-please, security scanning, coverage, PR checks, stale issue cleanup, label sync

## Planned (Not Yet Implemented)

- **Frontend**: SvelteKit 2, Svelte 5, Tailwind CSS 4, shadcn-svelte
- **Content Modules**: Music, Audiobooks, Books, Podcasts, Comics, Photos, Live TV
- **Playback**: Watch Next, Skip Intro, SyncPlay, Trickplay, Transcoding (Blackbeard)
- **Integrations**: Lidarr, Whisparr, Chaptarr, Authelia, Authentik, Keycloak, scrobbling (Trakt, Last.fm)
- **Features**: Collections, Content Ratings, Request System, Release Calendar

See [docs/dev/design/](docs/dev/design/) for comprehensive design documentation covering all planned features.

---

## Architecture

| Component   | Technology          | Notes                            |
|-------------|---------------------|----------------------------------|
| Language    | Go 1.26+            | `GOEXPERIMENT=jsonv2,goroutineleakprofile,simd,runtimesecret` |
| Database    | PostgreSQL 18+      | sqlc for type-safe queries       |
| Cache       | Dragonfly + rueidis | Redis-compatible, auto-pipelining|
| Local Cache | otter v2.x          | W-TinyLFU eviction               |
| Search      | Typesense           | typesense-go/v2 client           |
| Job Queue   | River               | PostgreSQL-native                |
| API         | ogen                | OpenAPI spec-first               |
| DI          | uber-go/fx          | Dependency injection             |
| Config      | koanf v2            | Multi-source configuration       |
| Logging     | slog + zap          | Structured logging               |
| Resilience  | gobreaker + backoff | Circuit breakers, retries        |

See [ARCHITECTURE.md](docs/dev/design/architecture/ARCHITECTURE.md) for the complete design.

---

## Quick Start

### Docker Compose (Recommended)

```bash
git clone https://github.com/lusoris/revenge.git
cd revenge
docker compose up -d
```

### Development

```bash
# Prerequisites: Go 1.25+, Docker (for dependencies)
git clone https://github.com/lusoris/revenge.git
cd revenge

# Start dependencies (PostgreSQL, Dragonfly, Typesense)
docker compose -f docker-compose.dev.yml up -d

# Build and run
make build && ./bin/revenge

# Or run directly
make run
```

---

## Development

All commands are in the [Makefile](Makefile). Run `make help` for the full list.

```bash
# Build
make build                  # Build binary
make build-linux            # Cross-compile for Linux

# Test
make test                   # Unit tests with race detection
make test-integration       # Integration tests (requires Docker)

# Code quality
make lint                   # golangci-lint
make vet                    # go vet
make fmt                    # Format code

# Code generation
make generate               # ogen + sqlc + go generate

# Docker
make docker-build           # Build Docker image
make docker-scan            # Build + Trivy scan
```

Integration tests use [testcontainers-go](https://testcontainers.com/) to spin up real PostgreSQL, Dragonfly, and Typesense instances in Docker.

---

## Documentation

- [Architecture](docs/dev/design/architecture/ARCHITECTURE.md) - System architecture
- [Tech Stack](docs/dev/design/technical/TECH_STACK.md) - Technology choices with rationale
- [Design Index](docs/dev/design/DESIGN_INDEX.md) - Full design documentation index

---

## Project Status

**Current Phase**: Backend services and content modules

See the [project roadmap](TODO.md) for detailed implementation status.

---

## Contributing

See the [contributing guide](CONTRIBUTING.md) for guidelines.

---

## License

GNU Affero General Public License v3.0 - see [LICENSE](LICENSE).

---

## Contact

- **Issues**: [GitHub Issues](https://github.com/lusoris/revenge/issues)
- **Security**: See the [security policy](SECURITY.md)
