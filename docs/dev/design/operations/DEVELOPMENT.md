## Development Quick Start

### Prerequisites

- Go 1.25+
- Docker & Docker Compose (optional)
- Git

### Setup (Windows PowerShell)

```powershell
# Clone repository

<!-- SOURCES: conventional-commits, dragonfly, koanf, ogen, pgx, postgresql-arrays, postgresql-json, river, sqlc, sqlc-config, typesense, typesense-go -->

<!-- DESIGN: operations, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES, 03_METADATA_SYSTEM -->


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

| Dimension | Status |
|-----------|--------|
| Design | 🔴 |
| Sources | 🔴 |
| Instructions | 🔴 |
| Code | 🔴 |
| Linting | 🔴 |
| Unit Testing | 🔴 |
| Integration Testing | 🔴 |
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


---

