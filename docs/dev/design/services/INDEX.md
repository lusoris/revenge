# Services Documentation

> Core application services implementing business logic

---

## Overview

Services in Revenge implement business logic and coordinate between repositories, external APIs, and background jobs. All services follow a consistent pattern:

- Constructor via `NewService(...)` accepting dependencies
- Structured logging via `slog`
- Context-aware operations
- Error wrapping with context
- fx module for dependency injection

---

## Service Layer Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        API Layer                             │
│                   (ogen handlers)                            │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                    Service Layer                             │
│  ┌─────────┐ ┌─────────┐ ┌──────────┐ ┌─────────────────┐  │
│  │  Auth   │ │  User   │ │ Session  │ │    Library      │  │
│  └────┬────┘ └────┬────┘ └────┬─────┘ └────────┬────────┘  │
│       │           │           │                 │           │
│  ┌────▼────┐ ┌────▼────┐ ┌────▼─────┐ ┌────────▼────────┐  │
│  │ Metadata│ │  RBAC   │ │ Activity │ │    Settings     │  │
│  └─────────┘ └─────────┘ └──────────┘ └─────────────────┘  │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                  Repository Layer                            │
│                   (sqlc queries)                             │
└─────────────────────────────────────────────────────────────┘
```

---

## Services

### Authentication & Users

| Service | Location | Description |
|---------|----------|-------------|
| [Auth](AUTH.md) | `internal/service/auth/` | Login, logout, registration, password management |
| [User](USER.md) | `internal/service/user/` | User CRUD, authentication, profile management |
| [Session](SESSION.md) | `internal/service/session/` | Session tokens, device tracking, activity |
| [OIDC](OIDC.md) | `internal/service/oidc/` | OIDC/SSO provider configuration |
| [API Keys](APIKEYS.md) | `internal/service/apikeys/` | API key generation and validation |

### Content & Libraries

| Service | Location | Description |
|---------|----------|-------------|
| [Library](LIBRARY.md) | `internal/service/library/` | Library CRUD, access control, scanning |
| [Metadata](METADATA.md) | `internal/service/metadata/` | TMDb, Radarr metadata providers |

### Access Control & Audit

| Service | Location | Description |
|---------|----------|-------------|
| [RBAC](RBAC.md) | `internal/service/rbac/` | Casbin-based role permissions |
| [Activity](ACTIVITY.md) | `internal/service/activity/` | Audit logging, event tracking |

### Audit & Configuration

| Service | Location | Description |
|---------|----------|-------------|
| [Activity](ACTIVITY.md) | `internal/service/activity/` | Audit logging, event tracking |
| [Settings](SETTINGS.md) | `internal/service/settings/` | Server settings persistence |

---

## Common Patterns

### Service Structure

```go
// Package myservice provides X operations.
package myservice

type Service struct {
    queries *db.Queries  // Database access
    logger  *slog.Logger // Structured logging
}

func NewService(queries *db.Queries, logger *slog.Logger) *Service {
    return &Service{
        queries: queries,
        logger:  logger.With(slog.String("service", "myservice")),
    }
}
```

### fx Module Pattern

```go
// module.go
package myservice

import "go.uber.org/fx"

var Module = fx.Options(
    fx.Provide(NewService),
)
```

### Error Handling

```go
var (
    ErrNotFound     = errors.New("not found")
    ErrAccessDenied = errors.New("access denied")
)

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Entity, error) {
    entity, err := s.queries.GetByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("get entity: %w", err)
    }
    return &entity, nil
}
```

---

## Related Documentation

- [Architecture](../architecture/ARCHITECTURE_V2.md) - System design
- [Database](../integrations/infrastructure/POSTGRESQL.md) - PostgreSQL patterns
- [Tech Stack](../technical/TECH_STACK.md) - Technology choices
