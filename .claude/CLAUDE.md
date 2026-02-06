# Claude Code Instructions

**Project**: Revenge - self-hosted media server
**Stack**: Go 1.25 backend, SvelteKit frontend, PostgreSQL 18, Dragonfly, Typesense

---

## Entry Point

### Source of Truth

**[/docs/dev/design/00_SOURCE_OF_TRUTH.md](../docs/dev/design/00_SOURCE_OF_TRUTH.md)** - package versions, API structure, database schemas, design decisions.

Always reference this before writing code. Do not use outdated versions or deprecated packages.

### Design Documentation

- [DESIGN_INDEX.md](../docs/dev/design/DESIGN_INDEX.md) - navigate to specific design docs
- [TECH_STACK.md](../docs/dev/design/technical/TECH_STACK.md) - technology choices and rationale
- [ARCHITECTURE.md](../docs/dev/design/architecture/01_ARCHITECTURE.md) - system architecture
- [.workingdir/](../.workingdir/) - working reports and analysis (migration reports, API overviews, etc.)

---

## Build & Test

All commands are in the [Makefile](../Makefile). Run `make help` for the full list.

```bash
# Build
make build                  # Build binary
make build-linux            # Cross-compile for Linux (amd64 + arm64)

# Test (local = same as CI)
make test                   # Unit tests with race detection + coverage
make test-short             # Fast unit tests (skip slow)
make test-integration       # Integration tests (requires Docker)
make test-all               # Unit + integration

# Code quality
make lint                   # golangci-lint
make vet                    # go vet
make vuln                   # govulncheck
make fmt                    # Format code

# Docker
make docker-build           # Build Docker image
make docker-scan            # Build + Trivy scan
make docker-test            # Full stack smoke test

# CI pipeline (same as GitHub Actions)
make ci                     # lint + test + docker-build + docker-scan

# Database
make migrate-up             # Run migrations
make migrate-create NAME=x  # Create new migration

# Code generation
make generate               # ogen + sqlc + go generate
```

### Go Environment

```bash
export GOEXPERIMENT=greenteagc,jsonv2  # Required (set in Makefile)
```

---

## Project Structure

```
cmd/revenge/            # Application entrypoint
internal/
  api/ogen/             # Generated API (ogen from OpenAPI spec)
  content/{module}/     # Content modules (movie, tv, music, etc.)
  service/{service}/    # Backend services
  infra/database/
    migrations/shared/  # SQL migrations (embedded via go:embed)
    migrate.go          # Migration runner (iofs + pgx)
tests/integration/      # Integration tests (testcontainers)
api/openapi/            # OpenAPI spec
charts/revenge/         # Helm chart (full stack)
scripts/                # docker-entrypoint.sh
docs/dev/design/        # Design documentation
```

### Migrations

Migrations live in `internal/infra/database/migrations/shared/` and are embedded into the binary via `//go:embed`. The Docker entrypoint runs `revenge migrate up` automatically on startup. No separate init SQL needed.

```bash
make migrate-create NAME=add_users_table  # Creates new migration pair
make migrate-up                           # Apply migrations
```

---

## Skills

Located in [skills/](skills/):

1. **coder-template** - Manage Coder workspace templates (`.coder/template.tf`)
2. **coder-workspace** - Manage Coder workspace operations
3. **setup-workspace** - Set up development environment

---

## CI/CD

### GitHub Actions Workflows

| Workflow | Trigger | Purpose |
|----------|---------|---------|
| `ci.yml` | push/PR to main, develop | Lint, unit tests, Docker build + Trivy, govulncheck, integration tests |
| `develop.yml` | push to develop | Auto-build `:develop` Docker image + dev Helm chart to GHCR |
| `release-please.yml` | push to main | Automated releases, multi-arch Docker + Helm chart to GHCR |
| `security.yml` | schedule + PR | CodeQL, Trivy, govulncheck, dependency review |
| `coverage.yml` | PR | Coverage report as PR comment |
| `pr-checks.yml` | PR | Title format, branch name, merge conflicts |
| `stale.yml` | schedule | Close stale issues/PRs |
| `labels.yml` | push to main | Sync GitHub labels |

### Branch Strategy

- **develop** - working branch, auto-builds on push
- **main** - stable/release, Release Please creates release PRs

### Container Registry

Images pushed to `ghcr.io/lusoris/revenge`. Helm charts to `ghcr.io/lusoris/charts/revenge`.

---

## Deployment

### Docker Compose (simple)

```bash
docker compose -f docker-compose.prod.yml up -d
```

### Helm (Kubernetes)

```bash
helm install revenge oci://ghcr.io/lusoris/charts/revenge
```

Full stack included: Revenge, PostgreSQL, Dragonfly, Typesense. See [charts/revenge/values.yaml](../charts/revenge/values.yaml).

---

## Code Patterns

- **Go 1.25**: Context-first APIs, error wrapping with `%w`, structured logging with slog
- **DI**: fx modules for dependency injection
- **Repos**: Repository pattern with interfaces
- **Testing**: Table-driven tests, testify assertions, mockery mocks, testcontainers for integration
- **Caching**: otter (L1 in-memory), rueidis (L2 distributed), sturdyc (request coalescing)

### Conventional Commits

```
feat: add user authentication
fix: resolve database connection issue
docs: update API documentation
```

---

## QAR (Adult Content) Module

Pirate-themed obfuscation for adult content isolation.

- **URL Pattern**: `/api/v1/legacy/*`
- **Database Schema**: `qar.*`
- **Access Control**: Requires `legacy:read` scope
- See [ADULT_CONTENT_SYSTEM.md](../docs/dev/design/features/adult/ADULT_CONTENT_SYSTEM.md)

---

## Troubleshooting

- **Build fails**: Check `go version` (1.25+), verify GOEXPERIMENT, run `go mod tidy`
- **Tests fail**: Check Docker running (integration), run with `-race`, check logs
- **LSP**: Check `gopls version`, restart LSP server

---

**Last Updated**: 2026-02-06
