# Development

<!-- DESIGN: operations -->

**Go**: 1.26+ with `GOEXPERIMENT=jsonv2,goroutineleakprofile,simd,runtimesecret`
**Build**: Makefile (`make help` for all targets)

> Development environment setup, build commands, and tooling

---

## Prerequisites

- Go 1.26+
- Docker (for dev services + integration tests)
- Make

## Quick Start

```bash
# Install dev tools (golangci-lint, air, migrate, sqlc, mockery, dlv, govulncheck)
make install-tools

# Start dev services (PostgreSQL, Dragonfly, Typesense)
make docker-up

# Run with hot reload
make dev

# Or build and run
make build && ./bin/revenge
```

## Makefile Targets

### Build

| Target | Purpose |
|--------|---------|
| `build` | Build binary to `bin/revenge` |
| `build-linux` | Cross-compile for Linux (amd64 + arm64) |
| `run` | Build and run |
| `dev` | Run with hot reload (requires air) |

### Test

| Target | Purpose |
|--------|---------|
| `test` | Unit tests with race detection + coverage |
| `test-short` | Fast unit tests (skip slow) |
| `test-integration` | Integration tests (requires Docker) |
| `test-all` | Unit + integration |
| `test-coverage` | Run tests and open coverage report |

### Docker

| Target | Purpose |
|--------|---------|
| `docker-build` | Build Docker image |
| `docker-scan` | Build + Trivy security scan |
| `docker-test` | Full stack smoke test |
| `docker-up` | Start dev services |
| `docker-down` | Stop dev services |

### Code Quality

| Target | Purpose |
|--------|---------|
| `lint` | golangci-lint v2.8.0 (5m timeout) |
| `fmt` | Format code |
| `vet` | go vet |
| `vuln` | govulncheck |

### Database

| Target | Purpose |
|--------|---------|
| `migrate-up` | Run migrations |
| `migrate-down` | Rollback one migration |
| `migrate-down-all` | Rollback all migrations |
| `migrate-force VERSION=N` | Force migration version |
| `migrate-version` | Show current version |
| `migrate-create NAME=x` | Create new migration pair |

### Code Generation

| Target | Purpose |
|--------|---------|
| `generate` | Run all generation (ogen + sqlc + go generate) |
| `ogen` | Generate API code from OpenAPI spec |
| `sqlc` | Generate database query code |

### Other

| Target | Purpose |
|--------|---------|
| `ci` | Full CI pipeline (lint + test + build + scan) |
| `all` | All checks and build (clean + deps + lint + test + build) |
| `clean` | Remove build artifacts |
| `deps` | Download and verify dependencies |
| `tidy` | Tidy go.mod |
| `install-tools` | Install all dev tools |

## Environment Variables

```bash
GOEXPERIMENT=jsonv2,goroutineleakprofile,simd,runtimesecret    # Required (set in Makefile)
DATABASE_URL=postgres://revenge:revenge_dev_pass@localhost:5432/revenge?sslmode=disable
```

## Development Environments

### VS Code Devcontainer

`.devcontainer/devcontainer.json` — Full development environment:

- Base: `mcr.microsoft.com/devcontainers/go:1-1.25-bookworm`
- Includes: ffmpeg, postgresql-client, all Go tools
- Services: PostgreSQL 18, Dragonfly, Typesense 0.25.2
- Forwarded ports: 8096 (API), 5173 (SvelteKit), 5432 (PostgreSQL), 6379 (Dragonfly), 8108 (Typesense)
- Post-create: `go mod download && make install-tools`

### Coder Template

`.coder/template.tf` — Terraform-based workspace:

- Backends: Docker (default), Kubernetes, K3s, Docker Swarm
- IDEs: Zed (SSH, default), VS Code (browser/desktop), JetBrains Gateway, Terminal
- Resources: 2-16 CPU cores (4 default), 4-16 GB memory (8 default), 30Gi disk
- Images: Go 1.25.7, PostgreSQL 18-alpine, Dragonfly latest, Typesense 0.25.2

## Linting Configuration

`.golangci.yml`:

- **Linters**: errcheck, govet, ineffassign, staticcheck, unused
- **govet**: All checks except shadow, fieldalignment
- **staticcheck**: All checks except QF1001, QF1008, ST1000 (package comments), ST1003 (naming)
- **Test exclusions**: gosec, errcheck, gocritic, staticcheck excluded in `_test.go`
- **Timeout**: 5m

## Related Documentation

- [SETUP.md](SETUP.md) - Production deployment
- [CI_CD.md](CI_CD.md) - CI/CD pipelines
- [../technical/TESTING.md](../technical/TESTING.md) - Testing infrastructure
