# Contributing to Revenge

Thank you for considering contributing to Revenge!

## Code of Conduct

Be respectful, inclusive, and professional.

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

- Check if the feature is already in our [design docs](docs/dev/design/DESIGN_INDEX.md)
- Explain the use case
- Open an issue for discussion

### Pull Requests

1. Fork the repo
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Write/update tests
5. Ensure tests pass (`make test`)
6. Run linters (`make lint`)
7. Commit with conventional commit messages
8. Push and open a PR

## Development Setup

```bash
# Prerequisites: Go 1.25+, Docker & Docker Compose

# Clone your fork
git clone https://github.com/YOUR_USERNAME/revenge.git
cd revenge

# Start dependencies (PostgreSQL, Dragonfly, Typesense)
docker compose -f docker-compose.dev.yml up -d

# Build
make build

# Run tests
make test

# Run application
make run
```

See the [Makefile](Makefile) for all available commands (`make help`).

## Project Structure

```
revenge/
├── cmd/revenge/              # Application entrypoint
├── internal/
│   ├── api/                  # HTTP handlers + ogen-generated API
│   │   ├── middleware/       # Rate limiting, request context
│   │   └── ogen/             # Generated OpenAPI code
│   ├── app/                  # fx module wiring
│   ├── config/               # koanf configuration
│   ├── content/              # Content modules
│   │   ├── movie/            # Movie module (service, repo, handler, jobs)
│   │   ├── tvshow/           # TV show module
│   │   ├── qar/              # Adult content (isolated schema)
│   │   └── shared/           # Shared library/scanner/matcher
│   ├── service/              # Backend services (auth, user, rbac, etc.)
│   ├── infra/                # Infrastructure (database, cache, jobs, health)
│   ├── integration/          # External integrations (radarr, sonarr)
│   ├── crypto/               # Encryption, hashing
│   ├── errors/               # Error types
│   ├── validate/             # Input validation
│   └── testutil/             # Test helpers (testcontainers)
├── tests/integration/         # Integration tests
├── api/openapi/               # OpenAPI specification
├── charts/revenge/            # Helm chart
├── scripts/                   # Docker entrypoint
└── docs/dev/design/           # Design documentation
```

## Coding Standards

### Go Style

Follow standard Go conventions:
- `gofmt` for formatting
- `golangci-lint` for linting
- Effective Go practices

### Patterns

- **Context-first APIs**: All service methods take `context.Context` as first parameter
- **Error wrapping**: Always wrap errors with `fmt.Errorf("context: %w", err)`
- **Repository pattern**: Interfaces in domain, implementations in adapters
- **fx modules**: Each package exports `var Module = fx.Module(...)`
- **Table-driven tests**: Use testify for assertions

### Example

```go
// Always wrap errors with context
func (s *Service) GetMedia(ctx context.Context, id uuid.UUID) (*Media, error) {
    media, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("get media %s: %w", id, err)
    }
    return media, nil
}
```

### Testing

```go
// Table-driven tests with testify
func TestService_Create(t *testing.T) {
    tests := []struct {
        name    string
        input   CreateRequest
        wantErr bool
    }{
        {name: "valid", input: CreateRequest{Name: "test"}, wantErr: false},
        {name: "empty name", input: CreateRequest{}, wantErr: true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // ...
        })
    }
}
```

## Commit Messages

Use [conventional commits](https://www.conventionalcommits.org/):

```
feat: add user authentication
fix: resolve database connection issue
docs: update API documentation
test: add integration tests for library scanner
refactor: simplify metadata adapter
perf: optimize database queries
```

## Language Policy

All code comments, documentation, commit messages, and PR descriptions must be in **English**.

## Testing

```bash
make test                # Unit tests with race detection
make test-integration    # Integration tests (requires Docker)
```

- Aim for 80%+ coverage
- Use table-driven tests
- Use testcontainers for database integration tests
- Mock external dependencies with mockery

## Review Process

1. CI/CD checks must pass
2. At least one maintainer approval required
3. No unresolved review comments
4. Up-to-date with develop branch

## Questions?

- Open a [GitHub Issue](https://github.com/lusoris/revenge/issues)
- Check existing [issues](https://github.com/lusoris/revenge/issues) and [discussions](https://github.com/lusoris/revenge/discussions)
