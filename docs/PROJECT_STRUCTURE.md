# Project Structure

This document describes the organization of the Jellyfin Go codebase.

## Directory Layout

```
jellyfin-go/
├── cmd/                    # Application entry points
│   └── jellyfin/          # Main application
│       └── main.go        # Application entry point
│
├── internal/              # Private application code
│   ├── api/              # HTTP API layer
│   │   ├── handlers/     # HTTP handlers
│   │   └── middleware/   # HTTP middleware
│   │
│   ├── domain/           # Domain entities and business rules
│   │
│   ├── service/          # Business logic
│   │
│   └── infra/            # Infrastructure layer
│       ├── database/     # Database implementations
│       └── cache/        # Cache implementations
│
├── pkg/                   # Public libraries
│   ├── config/           # Configuration management
│   └── logger/           # Logging utilities
│
├── migrations/            # Database migrations
│
├── configs/               # Configuration files
│   ├── config.yaml       # Default configuration
│   └── config.dev.yaml   # Development configuration
│
├── tests/                 # Test files
│   └── integration/      # Integration tests
│
├── scripts/               # Build and development scripts
│   ├── dev.sh            # Development helper (Unix)
│   └── dev.ps1           # Development helper (Windows)
│
├── docs/                  # Documentation
│   ├── ARCHITECTURE.md   # Architecture documentation
│   ├── SETUP.md          # Setup instructions
│   └── ...
│
├── .github/              # GitHub configuration
│   └── workflows/        # CI/CD workflows
│       ├── ci.yml        # Continuous Integration
│       └── release.yml   # Release automation
│
├── .vscode/              # VS Code configuration
│   ├── settings.json     # Editor settings
│   ├── launch.json       # Debug configurations
│   └── extensions.json   # Recommended extensions
│
├── Dockerfile             # Production Docker image
├── docker-compose.yml     # Docker Compose for production
├── docker-compose.dev.yml # Docker Compose for development
├── Makefile              # Build automation
├── go.mod                # Go module definition
├── go.sum                # Go module checksums
├── .gitignore            # Git ignore rules
├── .gitattributes        # Git attributes
├── .editorconfig         # Editor configuration
├── .golangci.yml         # Linter configuration
├── .goreleaser.yml       # Release configuration
├── .air.toml             # Hot reload configuration
├── LICENSE               # License file
├── README.md             # Project README
├── CONTRIBUTING.md       # Contributing guidelines
└── CHANGELOG.md          # Version history
```

## Package Organization

### cmd/
Contains application entry points. Each subdirectory represents a separate binary.

### internal/
Private application code that cannot be imported by other projects.

- **api/**: HTTP layer (handlers, middleware, routing)
- **domain/**: Core business entities and domain logic
- **service/**: Business logic and use cases
- **infra/**: Infrastructure implementations (database, cache, etc.)

### pkg/
Public libraries that can be imported by other projects.

### migrations/
Database migration files managed by golang-migrate.

### configs/
Configuration files for different environments.

### tests/
Test files, organized by type (unit, integration, e2e).

### scripts/
Helper scripts for development, building, and deployment.

## Naming Conventions

- **Packages**: lowercase, single word when possible
- **Files**: lowercase with underscores (snake_case)
- **Types**: PascalCase
- **Functions**: camelCase (unexported) or PascalCase (exported)
- **Constants**: PascalCase or SCREAMING_SNAKE_CASE
- **Interfaces**: PascalCase, often ending with -er suffix

## Import Organization

Imports should be organized in three groups:

1. Standard library packages
2. External dependencies
3. Internal packages

Example:
```go
import (
    "context"
    "fmt"
    
    "github.com/gorilla/mux"
    "go.uber.org/zap"
    
    "github.com/jellyfin/jellyfin-go/internal/domain"
    "github.com/jellyfin/jellyfin-go/pkg/logger"
)
```

## Testing

- Unit tests: `*_test.go` files alongside source code
- Integration tests: `tests/integration/`
- Test helpers: `tests/testutil/`

## Documentation

- Package documentation: godoc comments on package declaration
- Function documentation: comments above function declarations
- Complex logic: inline comments explaining the why, not the what
