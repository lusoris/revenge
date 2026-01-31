# Claude Code Instructions for Revenge Project

**Purpose**: General guidance for Claude Code when assisting with the Revenge media server project across all tools

**Last Updated**: 2026-01-31

---

## Project Overview

**Revenge** is a modern media server. For exact versions of all technologies, see:

**📋 SOURCE OF TRUTH**: [../../docs/dev/design/00_SOURCE_OF_TRUTH.md](../../docs/dev/design/00_SOURCE_OF_TRUTH.md)

**CRITICAL**: Always reference SOURCE_OF_TRUTH for:
- Go version and required GOEXPERIMENT flags
- Node.js/npm versions
- Python version
- Database versions (PostgreSQL, Dragonfly, Typesense)
- Frontend framework versions (Svelte, SvelteKit)
- All other tool and package versions

**Never hardcode versions** - they change and must be kept in sync with SOURCE_OF_TRUTH.

---

## Tool-Specific Instructions

When the user's environment or question involves specific tools, refer to tool-specific CLAUDE.md files:

- **VS Code**: [../../.vscode/docs/CLAUDE.md](../../.vscode/docs/CLAUDE.md) (if exists)
- **Zed**: [../../.zed/docs/CLAUDE.md](../../.zed/docs/CLAUDE.md) (if exists)
- **JetBrains**: [../../.jetbrains/docs/CLAUDE.md](../../.jetbrains/docs/CLAUDE.md)
- **Coder**: [../../.coder/docs/CLAUDE.md](../../.coder/docs/CLAUDE.md) (if exists)

---

## Development Environment Detection

### Determine User's Environment

**Local development**:
- User runs commands from project root
- Services run via Docker Compose
- Uses local IDE (VS Code, Zed, GoLand, etc.)

**Remote development (Coder)**:
- User mentions Coder, Gateway, or workspace
- Commands run inside workspace via SSH
- Services run within workspace
- Uses remote IDE (VS Code Remote, Zed SSH, JetBrains Gateway)

**IDE detection**:
- **VS Code**: User mentions extensions, settings.json, launch.json
- **Zed**: User mentions Zed, zed command, .zed/settings.json
- **JetBrains**: User mentions GoLand, IntelliJ IDEA, Gateway, run configurations
- **Other**: User mentions specific tool name

---

## File Structure

### Documentation Organization

```
revenge/
├── .vscode/          # VS Code configuration
│   └── docs/         # VS Code documentation
├── .zed/             # Zed configuration
│   └── docs/         # Zed documentation
├── .jetbrains/       # JetBrains configuration
│   └── docs/         # JetBrains documentation
├── .coder/           # Coder workspace template
│   └── docs/         # Coder documentation
├── .shared/          # Cross-tool integration
│   └── docs/         # Shared documentation
│       ├── INDEX.md              # Hub
│       ├── ONBOARDING.md         # New developer guide
│       ├── INTEGRATION.md        # Tool integration
│       ├── WORKFLOWS.md          # Development workflows
│       ├── SETTINGS_GUIDE.md     # Settings management
│       ├── PROFILES.md           # Developer profiles
│       ├── TOOL_COMPARISON.md    # When to use which tool
│       ├── TROUBLESHOOTING.md    # Cross-tool issues
│       └── CLAUDE.md             # This file
├── .claude/          # Claude Code configuration
│   └── docs/         # Claude Code documentation
├── .githooks/        # Git hooks
│   └── docs/         # Git hooks documentation
├── .github/          # GitHub Actions
│   └── docs/         # GitHub Actions documentation
└── docs/
    └── dev/
        └── design/   # Design documentation
            ├── 00_SOURCE_OF_TRUTH.md      # Tech stack versions
            ├── DESIGN_INDEX.md            # Design docs hub
            ├── architecture/              # Architecture docs
            ├── features/                  # Feature docs
            ├── integrations/              # Integration docs
            ├── services/                  # Service docs
            ├── technical/                 # Technical docs
            └── operations/                # Operations docs
```

### When to Reference Documentation

**Setup issues** → Tool-specific setup docs:
- VS Code: .vscode/docs/SETUP.md (if exists)
- Zed: .zed/docs/SETUP.md
- JetBrains: .jetbrains/docs/SETUP.md
- Coder: .coder/docs/INDEX.md

**Workflow questions** → .shared/docs/WORKFLOWS.md

**Tool comparison** → .shared/docs/TOOL_COMPARISON.md

**Design/architecture** → docs/dev/design/

---

## Code Style and Formatting

### Go

**Style**:
- Hard tabs (tab size 4)
- Follow standard Go conventions
- Use goimports for formatting (via gopls or IDE)

**Imports order**:
1. Standard library
2. Third-party packages
3. Revenge internal packages

**Example**:
```go
import (
    "context"
    "fmt"

    "github.com/jmoiron/sqlx"
    "go.uber.org/fx"

    "github.com/lusoris/revenge/internal/config"
    "github.com/lusoris/revenge/internal/database"
)
```

### Python (Scripts)

**Style**:
- 4 spaces (no tabs)
- Follow PEP 8
- Use ruff for linting and formatting

**Format command**:
```bash
ruff format scripts/
ruff check scripts/
```

### TypeScript/Svelte

**Style**:
- 2 spaces (no tabs)
- Use Prettier for formatting
- Follow Svelte conventions

**Format command**:
```bash
cd web
npm run format
```

**EditorConfig**:
All styles are enforced via `.editorconfig` - IDEs respect this automatically.

---

## Testing

### Go Tests

**Running tests**:
```bash
# All tests
go test ./...

# Specific package
go test ./internal/api/handlers

# With coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# With race detector
go test -race ./...
```

**Test conventions**:
- Test files end in `_test.go`
- Test functions start with `Test`
- Table-driven tests preferred
- Aim for 80%+ coverage

### Frontend Tests

```bash
cd web

# Run tests
npm run test

# Watch mode
npm run test:watch

# Coverage
npm run test:coverage
```

---

## Git Workflow

### Branch Strategy

**Main branches**:
- `main` - Production releases (protected)
- `develop` - Development (main working branch, protected)

**Feature branches**:
- `feature/feature-name` - New features
- `fix/bug-name` - Bug fixes
- `hotfix/issue-name` - Emergency fixes (from main)

### Commit Messages

**Format**: Conventional Commits

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types**:
- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation
- `refactor` - Code refactoring
- `test` - Tests
- `chore` - Maintenance

**Example**:
```
feat(api): add movie search endpoint

- Implement /api/v1/movies/search handler
- Add query parameter validation
- Include unit tests with 85% coverage
- Update OpenAPI spec

Closes #123
```

### Workflow

See [WORKFLOWS.md](WORKFLOWS.md) for complete workflows:
- Local development workflow
- Remote development workflow (Coder)
- Code review workflow
- Release workflow

---

## Architecture Patterns

### Dependency Injection (uber-go/fx)

Revenge uses fx for dependency injection:

```go
fx.New(
    fx.Provide(
        config.New,
        database.New,
        handlers.New,
    ),
    fx.Invoke(server.Run),
)
```

When adding new services, use fx.Provide pattern.

### Service Layer

**Pattern**:
```
handlers/ → services/ → repositories/ → database
```

- **Handlers**: HTTP request handling (thin layer)
- **Services**: Business logic
- **Repositories**: Data access
- **Database**: SQLC-generated queries

### Configuration (koanf)

Configuration via koanf with multiple sources:
1. Default values
2. Config file (YAML)
3. Environment variables
4. Command-line flags

See: [../../docs/dev/design/technical/CONFIGURATION.md](../../docs/dev/design/technical/CONFIGURATION.md)

---

## Common Commands

### Development

```bash
# Start services (local)
docker-compose -f docker-compose.dev.yml up -d

# Run backend (with hot reload)
air

# Run frontend
cd web
npm run dev

# Run tests
go test ./...

# Lint
golangci-lint run

# Format
goimports -w .
```

### Remote (Coder)

```bash
# Login
coder login https://coder.ancilla.lol

# List workspaces
coder list

# Start workspace
coder start revenge-dev

# SSH into workspace
coder ssh revenge-dev

# Port forward
coder port-forward revenge-dev --tcp 8096:8096
```

---

## Environment Variables

### Required

- `GOEXPERIMENT=greenteagc,jsonv2` - Go experimental features
- Database connection (auto-configured in dev)
- Cache connection (auto-configured in dev)

### Optional

- `LOG_LEVEL` - Logging level (debug, info, warn, error)
- `PORT` - API port (default: 8096)

See: [../../docs/dev/design/technical/CONFIGURATION.md](../../docs/dev/design/technical/CONFIGURATION.md)

---

## Database

### Schema Management

**Tool**: golang-migrate

**Migrations**:
```bash
# Create migration
migrate create -ext sql -dir migrations -seq add_users_table

# Apply migrations
migrate -path migrations -database "postgres://revenge:revenge@localhost:5432/revenge?sslmode=disable" up

# Rollback
migrate -path migrations -database "..." down 1
```

### Queries (SQLC)

**Generate queries**:
```bash
sqlc generate
```

Queries defined in `internal/database/queries/` → Generated code in `internal/database/sqlc/`

See: [../../docs/dev/design/operations/DEVELOPMENT.md](../../docs/dev/design/operations/DEVELOPMENT.md)

---

## API

### OpenAPI Spec

**Location**: `api/openapi.yaml`

**Generation**: Using ogen

```bash
go generate ./...
```

Generates:
- Server stubs
- Client code
- Type definitions

See: [../../docs/dev/design/technical/API.md](../../docs/dev/design/technical/API.md)

---

## Helpful References

### When User Asks About...

**"How do I set up locally?"**
→ [ONBOARDING.md](ONBOARDING.md) or tool-specific setup docs

**"How do I use Coder?"**
→ [../../.coder/docs/INDEX.md](../../.coder/docs/INDEX.md)

**"Which IDE should I use?"**
→ [TOOL_COMPARISON.md](TOOL_COMPARISON.md)

**"How do I do X workflow?"**
→ [WORKFLOWS.md](WORKFLOWS.md)

**"How do settings work?"**
→ [SETTINGS_GUIDE.md](SETTINGS_GUIDE.md)

**"What's the tech stack?"**
→ [../../docs/dev/design/00_SOURCE_OF_TRUTH.md](../../docs/dev/design/00_SOURCE_OF_TRUTH.md)

**"How does feature X work?"**
→ [../../docs/dev/design/DESIGN_INDEX.md](../../docs/dev/design/DESIGN_INDEX.md) → Find relevant doc

**"I'm getting error Y"**
→ [TROUBLESHOOTING.md](TROUBLESHOOTING.md) or tool-specific troubleshooting

---

## Best Practices

### When Assisting Users

1. **Detect environment first**:
   - Local or remote (Coder)?
   - Which IDE?
   - What OS?

2. **Provide context-specific commands**:
   - Local: Docker Compose commands
   - Remote: Coder workspace commands
   - IDE-specific: Run configurations vs terminal

3. **Reference documentation**:
   - Link to relevant docs
   - Don't duplicate entire docs in responses
   - Point to specific sections

4. **Use correct file paths**:
   - Absolute paths from project root
   - Use forward slashes (/)
   - Example: `internal/api/handlers/movie.go`

5. **Include verification steps**:
   - After suggesting changes, tell user how to verify
   - Example: "Run `go test ./...` to verify"

6. **Respect tool preferences**:
   - Don't suggest switching tools unless asked
   - Provide tool-appropriate solutions
   - Example: Don't suggest VS Code tasks.json to JetBrains user

---

## Code Examples

### When Providing Code

**Include context**:
```go
// In internal/api/handlers/movie.go

func (h *MovieHandler) GetMovie(ctx context.Context, req *api.GetMovieRequest) (*api.MovieResponse, error) {
    movie, err := h.movieService.GetByID(ctx, req.ID)
    if err != nil {
        return nil, fmt.Errorf("get movie: %w", err)
    }

    return &api.MovieResponse{
        ID: movie.ID,
        Title: movie.Title,
    }, nil
}
```

**Explain where it goes**:
- File path
- Function name
- Any dependencies needed

---

## Troubleshooting Approach

### General Debugging Steps

1. **Gather information**:
   - What's the exact error message?
   - Which tool/environment?
   - What were they trying to do?

2. **Check common issues first**:
   - Services running?
   - Dependencies installed?
   - Correct versions?

3. **Suggest verification**:
   - `go version` - Check Go version
   - `docker ps` - Check services
   - `coder list` - Check workspace (remote)

4. **Reference docs**:
   - Tool-specific troubleshooting docs
   - Shared troubleshooting doc
   - Design docs if architecture question

5. **Escalate if needed**:
   - Point to GitHub issues
   - Suggest checking logs
   - Reference official documentation

---

## Week 1 Priorities

Current focus: **Remote development on Coder ready THIS WEEK**

**Support all 4 IDE options**:
1. VS Code (browser) - In Coder workspace
2. VS Code (desktop + SSH) - Connect to Coder
3. Zed (SSH) - Connect to Coder workspace
4. JetBrains Gateway - Connect to Coder workspace

When user asks about remote development, prioritize Coder-related documentation and workflows.

---

## Version Information

**CRITICAL**: Always verify versions from SOURCE_OF_TRUTH:

📋 **Read**: [../../docs/dev/design/00_SOURCE_OF_TRUTH.md](../../docs/dev/design/00_SOURCE_OF_TRUTH.md)

**Before suggesting any installation or upgrade**:
1. Read SOURCE_OF_TRUTH to get current required versions
2. Use those exact versions in commands and suggestions
3. Never use hardcoded version numbers from memory

**Never suggest outdated versions** or incompatible packages.
**Never hardcode versions** - always reference SOURCE_OF_TRUTH.

---

**Last Updated**: 2026-01-31
**Maintained By**: Revenge Development Team
