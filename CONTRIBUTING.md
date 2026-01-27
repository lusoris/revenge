# Contributing to Jellyfin Go

First off, thank you for considering contributing to Jellyfin Go! 🎉

## Code of Conduct

Be respectful, inclusive, and professional. We're all here to build something great.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check existing issues. When creating a bug report, include:

- **Clear title and description**
- **Steps to reproduce**
- **Expected vs actual behavior**
- **Environment** (OS, Go version, deployment method)
- **Logs** (if applicable)

### Suggesting Features

Feature suggestions are welcome! Please:

- Check if the feature already exists in Jellyfin (C#)
- Explain the use case
- Consider if it fits the single-server focus (Phase 1-4)

### Pull Requests

1. Fork the repo
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Write/update tests
5. Ensure tests pass (`go test ./...`)
6. Run linters (`golangci-lint run`)
7. Commit with clear messages
8. Push and open a PR

## Development Setup

```bash
# Prerequisites
- Go 1.22+
- Docker & Docker Compose
- Git

# Clone your fork
git clone https://github.com/YOUR_USERNAME/jellyfin-go.git
cd jellyfin-go

# Install dependencies
go mod download

# Start development environment
docker-compose up -d postgres

# Run tests
go test ./...

# Run application
go run ./cmd/jellyfin
```

## Project Structure

```
jellyfin-go/
├── cmd/
│   └── jellyfin/          # Main application entry point
├── internal/              # Private application code
│   ├── domain/           # Domain entities and business rules
│   ├── service/          # Business logic
│   ├── infra/            # Infrastructure (DB, cache, etc.)
│   └── api/              # HTTP handlers
├── pkg/                   # Public libraries
│   ├── middleware/       # HTTP middleware
│   └── util/             # Utilities
├── docs/                  # Documentation
├── migrations/            # Database migrations
├── configs/               # Configuration files
└── tests/                 # Integration tests
```

## Coding Standards

### Go Style

Follow standard Go conventions:
- `gofmt` for formatting
- `golangci-lint` for linting
- Effective Go practices

### Naming

```go
// Good
func GetUserByID(id string) (*User, error)
type MediaRepository interface{}

// Bad
func get_user(id string) (*User, error)
type mediaRepo interface{}
```

### Error Handling

```go
// Always wrap errors with context
func (s *Service) GetMedia(id string) (*Media, error) {
    media, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get media %s: %w", id, err)
    }
    return media, nil
}
```

### Testing

```go
// Write table-driven tests
func TestUserService_Create(t *testing.T) {
    tests := []struct {
        name    string
        input   CreateUserRequest
        want    *User
        wantErr bool
    }{
        {
            name: "valid user",
            input: CreateUserRequest{Username: "test"},
            want: &User{Username: "test"},
            wantErr: false,
        },
        // More test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## Commit Messages

Use conventional commits:

```
feat: add user authentication
fix: resolve transcoding memory leak
docs: update setup guide
test: add integration tests for library scanner
refactor: simplify FFmpeg command builder
perf: optimize database queries
```

## Testing Guidelines

### Unit Tests

- Test business logic in isolation
- Mock external dependencies
- Aim for 80%+ coverage

```bash
go test ./internal/service/...
```

### Integration Tests

- Test API endpoints
- Use testcontainers for DB
- Test real FFmpeg integration

```bash
go test -tags=integration ./tests/...
```

### Performance Tests

```bash
go test -bench=. -benchmem ./internal/...
```

## Documentation

- Update README.md for user-facing changes
- Update docs/ for architecture changes
- Add inline comments for complex logic
- Update API docs (Swagger)

## Review Process

1. CI/CD checks must pass (tests, linting)
2. At least one maintainer approval required
3. No unresolved review comments
4. Up-to-date with main branch

## Release Process

We follow semantic versioning (SemVer):

- **v0.x.x**: Pre-release, breaking changes OK
- **v1.x.x**: Stable, backwards compatible
- **v2.x.x**: Major version, breaking changes

## Phase Priorities

**Current Focus: Phase 1 (Single-Server MVP)**

Priority areas for contributions:
1. Core API endpoints (Jellyfin compatibility)
2. Library scanning and metadata
3. SQLite/PostgreSQL support
4. Authentication and authorization
5. Bug fixes and performance

**Later Phases:**
- Phase 2: Transcoding and streaming
- Phase 3: Advanced features
- Phase 4: Security and polish
- Phase 5: Optional clustering

## Questions?

- Open a GitHub Discussion
- Join our Discord
- Check existing issues

Thank you for contributing! 🚀
