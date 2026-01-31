## Development Quick Start

### Prerequisites

- Go 1.25+
- Docker & Docker Compose (optional)
- Git

### Setup (Windows PowerShell)

```powershell
# Clone repository


<!-- TOC-START -->

## Table of Contents

- [Development Quick Start](#development-quick-start)
  - [Prerequisites](#prerequisites)
  - [Setup (Windows PowerShell)](#setup-windows-powershell)
  - [Setup (Linux/macOS)](#setup-linuxmacos)
  - [Available Commands](#available-commands)
    - [Using Make (recommended)](#using-make-recommended)
    - [Using Scripts](#using-scripts)
  - [First Run](#first-run)
  - [Development with Docker](#development-with-docker)
  - [Project Structure](#project-structure)
  - [Configuration](#configuration)
  - [Testing](#testing)
  - [Debugging](#debugging)
  - [Code Quality](#code-quality)
  - [Git Workflow](#git-workflow)
  - [Commit Message Format](#commit-message-format)
  - [Useful Resources](#useful-resources)
  - [Troubleshooting](#troubleshooting)
  - [Next Steps](#next-steps)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)

<!-- TOC-END -->

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | 🔴 |  |
| Sources | 🔴 |  |
| Instructions | 🔴 |  |
| Code | 🔴 |  |
| Linting | 🔴 |  |
| Unit Testing | 🔴 |  |
| Integration Testing | 🔴 |  |

---

git clone https://github.com/lusoris/revenge.git
cd revenge

# Run setup script
.\scripts\dev.ps1 setup

# Start development server
.\scripts\dev.ps1 dev
```

### Setup (Linux/macOS)

```bash
# Clone repository
git clone https://github.com/lusoris/revenge.git
cd revenge

# Make scripts executable
chmod +x scripts/*.sh

# Run setup script
./scripts/dev.sh setup

# Start development server
./scripts/dev.sh dev
```

### Available Commands

#### Using Make (recommended)

```bash
make help              # Show all available commands
make deps              # Download dependencies
make build             # Build binary
make run               # Run application
make test              # Run tests
make lint              # Run linters
make docker-build      # Build Docker image
make docker-compose-up # Start with Docker Compose
```

#### Using Scripts

**Windows:**
```powershell
.\scripts\dev.ps1 check         # Check requirements
.\scripts\dev.ps1 install-tools # Install dev tools
.\scripts\dev.ps1 test          # Run tests
.\scripts\dev.ps1 lint          # Run linter
.\scripts\dev.ps1 build         # Build binary
.\scripts\dev.ps1 dev           # Hot reload dev server
```

**Linux/macOS:**
```bash
./scripts/dev.sh check         # Check requirements
./scripts/dev.sh install-tools # Install dev tools
./scripts/dev.sh test          # Run tests
./scripts/dev.sh lint          # Run linter
./scripts/dev.sh build         # Build binary
./scripts/dev.sh dev           # Hot reload dev server
```

### First Run

```bash
# Download dependencies
go mod download

# Build and run
go run ./cmd/revenge
```

The server will start at http://localhost:8096

Test endpoints:
- http://localhost:8096/health/live
- http://localhost:8096/health/ready
- http://localhost:8096/version

### Development with Docker

```bash
# Development environment (PostgreSQL + Dragonfly + Typesense)
docker-compose -f docker-compose.dev.yml up

# Production-like environment
docker-compose up
```

### Project Structure

```
revenge/
├── api/               # OpenAPI specs, generated code (ogen)
├── cmd/               # Application entry points
├── internal/          # Private application code
│   ├── content/       # Content modules (movie, tvshow, music, qar)
│   └── infra/         # Infrastructure (database, cache, search)
│       └── database/
│           ├── migrations/  # Database migrations (golang-migrate)
│           └── queries/     # sqlc queries
├── pkg/               # Public libraries (resilience, etc.)
├── configs/           # Configuration files (koanf)
├── tests/             # Integration tests
├── testdata/          # Test fixtures
├── scripts/           # Helper scripts
└── docs/              # Documentation
```

### Configuration

Configuration can be set via:

1. **Config file** (`configs/config.yaml`)
2. **Environment variables** (prefixed with `REVENGE_`)
3. **Command-line flags** (coming soon)

Example environment variables:
```bash
export REVENGE_LOG_LEVEL=debug
export REVENGE_DB_TYPE=postgres
export REVENGE_DB_HOST=localhost
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific test
go test -v ./internal/domain/...

# Run with race detection
go test -race ./...
```

### Debugging

VS Code launch configurations are included:

1. Press `F5` to start debugging
2. Set breakpoints in your code
3. Use the Debug Console

Or use Delve directly:
```bash
dlv debug ./cmd/revenge
```

### Code Quality

```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Fix linting issues
golangci-lint run --fix

# Vet code
go vet ./...
```

### Git Workflow

1. Create feature branch: `git checkout -b feature/my-feature`
2. Make changes and commit: `git commit -m "feat: add my feature"`
3. Push branch: `git push origin feature/my-feature`
4. Open Pull Request on GitHub

### Commit Message Format

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add user authentication
fix: resolve database connection issue
docs: update API documentation
test: add unit tests for media service
chore: update dependencies
```

### Useful Resources

- [Architecture Documentation](../architecture/01_ARCHITECTURE.md)
- [Setup Guide](SETUP.md)
- [Contributing Guidelines](../../../../CONTRIBUTING.md)

### Troubleshooting

**Port 8096 already in use:**
```bash
# Windows
netstat -ano | findstr :8096
taskkill /PID <PID> /F

# Linux/macOS
lsof -ti:8096 | xargs kill -9
```

**Database connection fails (PostgreSQL):**
```bash
# Check PostgreSQL is running
docker-compose -f docker-compose.dev.yml ps

# Check connection
PGPASSWORD=password psql -h localhost -U revenge -d revenge -c "SELECT 1"

# Restart PostgreSQL
docker-compose -f docker-compose.dev.yml restart postgres
```

**Module download fails:**
```bash
# Clean module cache
go clean -modcache
go mod download
```

### Next Steps

1. Review the [TODO.md](../../../../TODO.md) for current sprint tasks
3. Pick a task and create a feature branch
4. Implement, test, and submit a Pull Request


<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Operations](INDEX.md)

### In This Section

- [Advanced Patterns & Best Practices](BEST_PRACTICES.md)
- [Branch Protection Rules](BRANCH_PROTECTION.md)
- [Database Auto-Healing & Consistency Restoration](DATABASE_AUTO_HEALING.md)
- [GitFlow Workflow Guide](GITFLOW.md)
- [Revenge - Reverse Proxy & Deployment Best Practices](REVERSE_PROXY.md)
- [revenge - Setup Guide](SETUP.md)

### Related Topics

- [Revenge - Architecture v2](../architecture/01_ARCHITECTURE.md) _Architecture_
- [Revenge - Design Principles](../architecture/02_DESIGN_PRINCIPLES.md) _Architecture_
- [Revenge - Metadata System](../architecture/03_METADATA_SYSTEM.md) _Architecture_
- [Revenge - Player Architecture](../architecture/04_PLAYER_ARCHITECTURE.md) _Architecture_
- [Plugin Architecture Decision](../architecture/05_PLUGIN_ARCHITECTURE_DECISION.md) _Architecture_

### Indexes

- [Design Index](../DESIGN_INDEX.md) - All design docs by category/topic
- [Source of Truth](../00_SOURCE_OF_TRUTH.md) - Package versions and status

<!-- DESIGN-BREADCRUMBS-END -->

---

<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

### Referenced Sources

| Source | Documentation |
|--------|---------------|
| [Conventional Commits](https://www.conventionalcommits.org/) | [Local](../../sources/standards/conventional-commits.md) |
| [Dragonfly Documentation](https://www.dragonflydb.io/docs) | [Local](../../sources/infrastructure/dragonfly.md) |
| [PostgreSQL Arrays](https://www.postgresql.org/docs/current/arrays.html) | [Local](../../sources/database/postgresql-arrays.md) |
| [PostgreSQL JSON Functions](https://www.postgresql.org/docs/current/functions-json.html) | [Local](../../sources/database/postgresql-json.md) |
| [Typesense API](https://typesense.org/docs/latest/api/) | [Local](../../sources/infrastructure/typesense.md) |
| [Typesense Go Client](https://github.com/typesense/typesense-go) | [Local](../../sources/infrastructure/typesense-go.md) |
| [koanf](https://pkg.go.dev/github.com/knadh/koanf/v2) | [Local](../../sources/tooling/koanf.md) |
| [ogen OpenAPI Generator](https://pkg.go.dev/github.com/ogen-go/ogen) | [Local](../../sources/tooling/ogen.md) |
| [pgx PostgreSQL Driver](https://pkg.go.dev/github.com/jackc/pgx/v5) | [Local](../../sources/database/pgx.md) |
| [sqlc](https://docs.sqlc.dev/en/stable/) | [Local](../../sources/database/sqlc.md) |
| [sqlc Configuration](https://docs.sqlc.dev/en/stable/reference/config.html) | [Local](../../sources/database/sqlc-config.md) |

<!-- SOURCE-BREADCRUMBS-END -->