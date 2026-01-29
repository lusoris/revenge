# Service Layer Patterns

> Instructions for implementing service layer components in Revenge.

## Service Structure

Every service MUST follow this structure:

```go
package servicename

import (
    "context"
    "log/slog"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/lusoris/revenge/internal/infra/database/db"
)

type Service struct {
    pool    *pgxpool.Pool
    queries *db.Queries
    logger  *slog.Logger
    // Additional dependencies
}

func NewService(pool *pgxpool.Pool, logger *slog.Logger) *Service {
    return &Service{
        pool:    pool,
        queries: db.New(pool),
        logger:  logger.With(slog.String("service", "servicename")),
    }
}
```

---

## Core Service Patterns

### Auth Service

Location: `internal/service/auth/`

```go
type Service struct {
    pool     *pgxpool.Pool
    queries  *db.Queries
    users    *user.Service
    sessions *session.Service
    logger   *slog.Logger
}

// Login authenticates user and creates session
type LoginParams struct {
    Username   string
    Password   string
    DeviceName *string
    DeviceType *string
    ClientName *string
    RememberMe bool
}

type LoginResult struct {
    User         *User
    AccessToken  string
    RefreshToken string
    SessionID    uuid.UUID
    ExpiresAt    time.Time
}

func (s *Service) Login(ctx context.Context, params LoginParams) (*LoginResult, error) {
    // 1. Find user by username
    // 2. Verify password with bcrypt
    // 3. Create session
    // 4. Return tokens
}

func (s *Service) Logout(ctx context.Context, sessionID uuid.UUID) error
func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (*LoginResult, error)
func (s *Service) ValidateToken(ctx context.Context, token string) (*Session, error)
```

### User Service

Location: `internal/service/user/`

```go
type Service struct {
    pool    *pgxpool.Pool
    queries *db.Queries
    rbac    *rbac.Service
    logger  *slog.Logger
}

type CreateParams struct {
    Username string
    Email    string
    Password string
    Role     Role
}

func (s *Service) Create(ctx context.Context, params CreateParams) (*User, error) {
    // 1. Validate input
    // 2. Check username/email uniqueness
    // 3. Hash password (bcrypt, cost 12)
    // 4. Insert user
    // 5. Assign RBAC role
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*User, error)
func (s *Service) GetByUsername(ctx context.Context, username string) (*User, error)
func (s *Service) Update(ctx context.Context, id uuid.UUID, params UpdateParams) (*User, error)
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error
func (s *Service) VerifyPassword(ctx context.Context, user *User, password string) error
func (s *Service) ChangePassword(ctx context.Context, id uuid.UUID, oldPass, newPass string) error
```

### Session Service

Location: `internal/service/session/`

```go
type Service struct {
    pool        *pgxpool.Pool
    queries     *db.Queries
    config      *SessionConfig
    logger      *slog.Logger
}

type SessionConfig struct {
    TokenLength     int           // default: 32 bytes
    SessionDuration time.Duration // default: 24h
    MaxSessions     int           // default: 10 per user
    RefreshWindow   time.Duration // default: 7 days
}

func (s *Service) Create(ctx context.Context, params CreateParams) (*Session, error) {
    // 1. Generate random token (32 bytes)
    // 2. Hash token with SHA-256 (store hash, return plain)
    // 3. Enforce max sessions per user
    // 4. Insert session
}

func (s *Service) Validate(ctx context.Context, token string) (*Session, error) {
    // 1. Hash incoming token
    // 2. Lookup by hash
    // 3. Check expiration
    // 4. Update last_used_at
}

func (s *Service) Revoke(ctx context.Context, sessionID uuid.UUID) error
func (s *Service) RevokeAll(ctx context.Context, userID uuid.UUID) error
func (s *Service) ListByUser(ctx context.Context, userID uuid.UUID) ([]*Session, error)
```

### Library Service

Location: `internal/service/library/`

```go
type Service struct {
    pool       *pgxpool.Pool
    queries    *db.Queries
    riverQueue *river.Client[pgx.Tx]
    logger     *slog.Logger
}

func (s *Service) Create(ctx context.Context, params CreateParams) (*Library, error)
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Library, error)
func (s *Service) List(ctx context.Context, userID uuid.UUID) ([]*Library, error)
func (s *Service) Update(ctx context.Context, id uuid.UUID, params UpdateParams) (*Library, error)
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error

// Scanning
func (s *Service) Scan(ctx context.Context, id uuid.UUID, full bool) error {
    // Queue River job for scanning
    _, err := s.riverQueue.Insert(ctx, &ScanLibraryArgs{
        LibraryID: id,
        FullScan:  full,
    }, nil)
    return err
}

func (s *Service) GetScanStatus(ctx context.Context, id uuid.UUID) (*ScanStatus, error)
```

---

## Transaction Patterns

### Single Transaction

```go
func (s *Service) CreateWithRelations(ctx context.Context, params Params) (*Entity, error) {
    tx, err := s.pool.Begin(ctx)
    if err != nil {
        return nil, fmt.Errorf("begin transaction: %w", err)
    }
    defer tx.Rollback(ctx)

    qtx := s.queries.WithTx(tx)

    // Multiple operations
    entity, err := qtx.CreateEntity(ctx, params.Entity)
    if err != nil {
        return nil, fmt.Errorf("create entity: %w", err)
    }

    for _, rel := range params.Relations {
        if err := qtx.CreateRelation(ctx, entity.ID, rel); err != nil {
            return nil, fmt.Errorf("create relation: %w", err)
        }
    }

    if err := tx.Commit(ctx); err != nil {
        return nil, fmt.Errorf("commit: %w", err)
    }

    return entity, nil
}
```

---

## Error Handling

Define service-specific errors:

```go
package user

import "errors"

var (
    ErrUserNotFound      = errors.New("user not found")
    ErrUsernameTaken     = errors.New("username already taken")
    ErrEmailTaken        = errors.New("email already taken")
    ErrInvalidPassword   = errors.New("invalid password")
    ErrPasswordMismatch  = errors.New("password mismatch")
)
```

Return wrapped errors for context:

```go
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
    user, err := s.queries.GetUserByID(ctx, id)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, ErrUserNotFound
        }
        return nil, fmt.Errorf("get user by id: %w", err)
    }
    return user, nil
}
```

---

## Caching Pattern

Use otter for local caching in services:

```go
import "github.com/maypok86/otter"

type Service struct {
    pool    *pgxpool.Pool
    queries *db.Queries
    cache   otter.Cache[uuid.UUID, *Entity]
    logger  *slog.Logger
}

func NewService(pool *pgxpool.Pool, logger *slog.Logger) (*Service, error) {
    cache, err := otter.MustBuilder[uuid.UUID, *Entity](10_000).
        WithTTL(5 * time.Minute).
        Build()
    if err != nil {
        return nil, err
    }

    return &Service{
        pool:    pool,
        queries: db.New(pool),
        cache:   cache,
        logger:  logger,
    }, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Entity, error) {
    // Check cache first
    if entity, ok := s.cache.Get(id); ok {
        return entity, nil
    }

    // Fetch from database
    entity, err := s.queries.GetEntityByID(ctx, id)
    if err != nil {
        return nil, err
    }

    // Cache result
    s.cache.Set(id, entity)
    return entity, nil
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, params UpdateParams) (*Entity, error) {
    entity, err := s.queries.UpdateEntity(ctx, id, params)
    if err != nil {
        return nil, err
    }

    // Invalidate cache
    s.cache.Delete(id)
    return entity, nil
}
```

---

## fx Module Registration

```go
package auth

import "go.uber.org/fx"

var Module = fx.Module("service/auth",
    fx.Provide(NewService),
)
```

Combined services module:

```go
// internal/service/module.go
package service

import (
    "go.uber.org/fx"

    "github.com/lusoris/revenge/internal/service/auth"
    "github.com/lusoris/revenge/internal/service/user"
    "github.com/lusoris/revenge/internal/service/session"
    "github.com/lusoris/revenge/internal/service/library"
    "github.com/lusoris/revenge/internal/service/settings"
)

var Module = fx.Module("services",
    auth.Module,
    user.Module,
    session.Module,
    library.Module,
    settings.Module,
)
```

---

## Related

- [INDEX.instructions.md](INDEX.instructions.md) - Main instruction index
- [content-modules.instructions.md](content-modules.instructions.md) - Content module patterns
- [sqlc-database.instructions.md](sqlc-database.instructions.md) - Database queries
- [otter-local-cache.instructions.md](otter-local-cache.instructions.md) - Local caching
- [Services Design](../../docs/dev/design/services/INDEX.md) - Service documentation
