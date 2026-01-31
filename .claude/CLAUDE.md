# Claude Code Instructions - Claude Code Agent

**Tool**: Claude Code
**Purpose**: AI-powered development assistant for Revenge project
**Documentation**: [docs/INDEX.md](docs/INDEX.md)

---

## Entry Point for Claude Code

When working on the Revenge project, **ALWAYS START HERE**:

### 1. Source of Truth (REQUIRED READING)

**[/docs/dev/design/00_SOURCE_OF_TRUTH.md](../docs/dev/design/00_SOURCE_OF_TRUTH.md)**

This document contains:
- Technology stack versions (Go, PostgreSQL, etc.)
- All Go dependencies with exact versions
- Frontend stack (SvelteKit, Svelte, Tailwind)
- Infrastructure components (Dragonfly, Typesense, River)
- API namespaces and structure
- Database schemas
- Configuration keys
- QAR (adult content) obfuscation terminology
- Cross-reference to external sources

**⚠️ CRITICAL**: Always reference this document for package versions, API structure, and design decisions. DO NOT use outdated versions or deprecated packages.

---

### 2. Tech Stack Details

**[/docs/dev/design/technical/TECH_STACK.md](../docs/dev/design/technical/TECH_STACK.md)**

Provides detailed rationale for technology choices:
- Modern Go (check SOURCE_OF_TRUTH) features used
- Dependency philosophy
- Frontend stack details
- Deployment platforms
- Development tools

---

### 3. Design Documentation

**[/docs/dev/design/DESIGN_INDEX.md](../docs/dev/design/DESIGN_INDEX.md)**

Navigate to specific design docs:
- Architecture
- Features (movies, TV, music, QAR, etc.)
- Integrations (metadata providers, Arr stack, auth)
- Services (backend services)
- Operations (setup, deployment, best practices)
- Technical (API, frontend, configuration)

---

## Project Overview

**Revenge** is a modern, self-hosted media server with Go backend, SvelteKit frontend, PostgreSQL database, and distributed caching.

**For complete tech stack and architecture details, see**:
- [00_SOURCE_OF_TRUTH.md](../docs/dev/design/00_SOURCE_OF_TRUTH.md) - All versions and dependencies
- [TECH_STACK.md](../docs/dev/design/technical/TECH_STACK.md) - Technology choices and rationale
- [ARCHITECTURE.md](../docs/dev/design/architecture/01_ARCHITECTURE.md) - System architecture

---

## Claude Code Configuration

### Skills Available

Located in [skills/](skills/):

1. **add-design-doc** - Create new design documents following conventions
2. **check-sources** - Check external documentation sources and cross-references
3. **run-pipeline** - Run doc/source pipelines for automation
4. **update-status** - Update design document status tables
5. **coder-template** - Manage Coder workspace templates
6. **coder-workspace** - Manage Coder workspace operations

### Permissions

Configured in `settings.local.json`:

- Git operations (status, diff, add, commit, push)
- Build commands (go build, test)
- Source fetching (scripts/fetch-sources.py)
- Web access to documentation sites

---

## Common Tasks

### Working on Design Documentation

1. **Always start** with [00_SOURCE_OF_TRUTH.md](../docs/dev/design/00_SOURCE_OF_TRUTH.md)
2. Use `add-design-doc` skill to create new design docs
3. Use `update-status` skill to update status tables
4. Use `check-sources` skill to verify external references
5. Follow existing document structure and conventions

### Fetching External Sources

1. Use `check-sources` skill
2. Or manually run: `python scripts/fetch-sources.py`
3. Sources are defined in `docs/dev/sources/SOURCES.yaml`
4. Fetched sources go to `docs/dev/sources/`

### Running Automation Pipelines

1. Use `run-pipeline` skill for:
   - Documentation generation
   - Index updates
   - Cross-reference generation
   - Source fetching

### Managing Coder Workspaces

1. Use `coder-template` skill to modify `.coder/template.tf`
2. Use `coder-workspace` skill to manage workspaces

---

## Development Workflow

### 1. Understanding a Feature

Before implementing, read:
1. SOURCE_OF_TRUTH for versions and architecture
2. Relevant design doc (see DESIGN_INDEX.md)
3. Integration docs (if working with external services)
4. Existing code patterns

### 2. Writing Code

Follow patterns in:
- **Go**: [TECH_STACK.md](../docs/dev/design/technical/TECH_STACK.md)
- **Module Structure**: [00_SOURCE_OF_TRUTH.md#project-structure](../docs/dev/design/00_SOURCE_OF_TRUTH.md#project-structure)
- **Testing**: [operations/BEST_PRACTICES.md](../docs/dev/design/operations/BEST_PRACTICES.md)

### 3. Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -coverprofile=coverage.out ./...

# Run specific package tests
go test ./internal/content/movie/...
```

### 4. Committing

Use conventional commits:
```
feat: add user authentication
fix: resolve database connection issue
docs: update API documentation
```

Git hooks will enforce this format.

---

## Technology-Specific Guidance

### Go Development

**Version**: [version from SOURCE_OF_TRUTH]

**Build Command**:
```bash
GOEXPERIMENT=greenteagc,jsonv2 go build ./...
```

**For complete dependency list with versions**, see [00_SOURCE_OF_TRUTH.md](../docs/dev/design/00_SOURCE_OF_TRUTH.md#go-dependencies)

**Patterns**:
- Context-first APIs
- Error wrapping with `%w`
- Structured logging with slog
- Table-driven tests

### Frontend Development

**Framework**: SvelteKit 2 with Svelte 5

**Styling**: Tailwind CSS 4

**Components**: shadcn-svelte

**State Management**:
- Svelte stores for local state
- TanStack Query for server state

### Database

**PostgreSQL 18+ only** (no SQLite support)

**For schema details, migrations, and query patterns, see**:
- [00_SOURCE_OF_TRUTH.md](../docs/dev/design/00_SOURCE_OF_TRUTH.md#database-schemas)
- [ARCHITECTURE.md](../docs/dev/design/architecture/01_ARCHITECTURE.md)

---

## QAR (Adult Content) Module

**Pirate-themed obfuscation** for adult content isolation.

**URL Pattern**: `/api/v1/legacy/*`
**Database Schema**: `qar.*`
**Access Control**: Requires `legacy:read` scope

**For full terminology mapping and details, see**:
- [00_SOURCE_OF_TRUTH.md](../docs/dev/design/00_SOURCE_OF_TRUTH.md#qar-obfuscation-terminology)
- [ADULT_CONTENT_SYSTEM.md](../docs/dev/design/features/adult/ADULT_CONTENT_SYSTEM.md)

---

## Best Practices

### 1. Always Reference SOURCE_OF_TRUTH

- For package versions
- For API structure
- For configuration keys
- For design decisions

### 2. Follow Existing Patterns

- Module structure: `internal/content/{module}/`
- Service structure: `internal/service/{service}/`
- Repository pattern with interfaces
- fx modules for dependency injection

### 3. Test Coverage

- **Minimum 80%** required
- Use table-driven tests
- Use testify for assertions
- Use mockery for mocks
- Integration tests with testcontainers

### 4. Documentation

- Update design docs when architecture changes
- Add inline comments for complex logic
- Update SOURCES.yaml when adding external dependencies
- Run pipelines to update indexes

### 5. Performance

- Use otter for L1 cache (in-memory)
- Use rueidis for L2 cache (distributed)
- Use sturdyc for request coalescing
- Follow [TECH_STACK.md](../docs/dev/design/technical/TECH_STACK.md) performance patterns

---

## Troubleshooting

### Build Fails

1. Check Go version: `go version` (must be 1.25+)
2. Verify GOEXPERIMENT flags: `greenteagc,jsonv2`
3. Check go.mod for version mismatches
4. Run `go mod tidy`

### Tests Fail

1. Check PostgreSQL is running (for integration tests)
2. Verify test database exists
3. Check for race conditions: `go test -race ./...`
4. Review test logs

### LSP Not Working

1. Check gopls is installed: `gopls version`
2. Verify IDE configuration (see `.vscode/CLAUDE.md` or `.zed/CLAUDE.md`)
3. Restart LSP server

---

## Related Documentation

### Tool-Specific
- **VS Code**: [.vscode/CLAUDE.md](../.vscode/CLAUDE.md)
- **Zed**: [.zed/CLAUDE.md](../.zed/CLAUDE.md)
- **Coder**: [.coder/docs/INDEX.md](../.coder/docs/INDEX.md)
- **Git Hooks**: [.githooks/docs/INDEX.md](../.githooks/docs/INDEX.md)
- **GitHub Actions**: [.github/docs/INDEX.md](../.github/docs/INDEX.md)

### Design Documentation
- **Design Index**: [../docs/dev/design/DESIGN_INDEX.md](../docs/dev/design/DESIGN_INDEX.md)
- **Architecture**: [../docs/dev/design/architecture/INDEX.md](../docs/dev/design/architecture/INDEX.md)
- **Features**: [../docs/dev/design/features/INDEX.md](../docs/dev/design/features/INDEX.md)
- **Operations**: [../docs/dev/design/operations/INDEX.md](../docs/dev/design/operations/INDEX.md)

### External Sources
- **Sources Index**: [../docs/dev/sources/SOURCES_INDEX.md](../docs/dev/sources/SOURCES_INDEX.md)
- **Design ↔ Sources**: [../docs/dev/sources/DESIGN_CROSSREF.md](../docs/dev/sources/DESIGN_CROSSREF.md)

---

## Quick Commands

```bash
# Build project
GOEXPERIMENT=greenteagc,jsonv2 go build ./...

# Run tests
go test ./...

# Run with coverage
go test -coverprofile=coverage.out ./...

# Lint code
golangci-lint run

# Format code
go fmt ./...

# Run migrations
make migrate-up

# Start dev server (with hot reload)
air

# Fetch external sources
python scripts/fetch-sources.py

# Generate documentation indexes
python scripts/generate-design-indexes.py
```

---

**⚠️ REMEMBER**: Always start with SOURCE_OF_TRUTH, follow existing patterns, and maintain test coverage above 80%.

---

**Last Updated**: 2026-01-31
**Maintained By**: Development Team
