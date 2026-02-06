# fx Module Patterns

> Three wiring patterns for dependency injection with uber/fx. Written from code as of 2026-02-06.

---

## Overview

Every package in `internal/` exposes an `fx.Module` that provides its dependencies. The app module (`internal/app/module.go`) composes all modules:

```go
var Module = fx.Module("app",
    // Infrastructure (order matters for dependencies)
    config.Module,
    logging.Module,
    database.Module,
    cache.Module,
    jobs.Module,
    health.Module,
    observability.Module,
    image.Module,
    storage.Module,

    // Services
    settings.Module,
    user.Module,
    auth.Module,
    session.Module,
    mfa.Module,
    oidc.Module,
    rbac.Module,
    apikeys.Module,
    activity.Module,
    email.Module,
    notification.Module,
    library.Module,
    search.Module,
    metadata.Module,

    // Content modules
    movie.Module,
    tvshow.Module,

    // Integrations
    radarr.Module,
    sonarr.Module,

    // Workers
    moviejobs.Module,
    tvshowjobs.Module,

    // API (always last — depends on everything)
    api.Module,
)
```

---

## Pattern 1: Simple (Direct Constructors)

Used when constructors match fx's expectations exactly (dependencies as parameters, result as return).

```go
// internal/service/settings/module.go
var Module = fx.Module("settings",
    fx.Provide(
        NewPostgresRepository,  // func(*pgxpool.Pool) Repository
        NewService,             // func(Repository) Service
    ),
)
```

fx resolves the dependency graph automatically: `*pgxpool.Pool` is provided by `database.Module`, passed to `NewPostgresRepository`, which returns `Repository`, which is passed to `NewService`.

**Used by:** settings, search, email, storage, health, logging

---

## Pattern 2: Config-Extracting (Inline Functions)

Used when the constructor needs config values extracted from `*config.Config`, or when the provided type needs to be an interface but the constructor returns a concrete type.

```go
// internal/service/apikeys/module.go
var Module = fx.Module("apikeys",
    fx.Provide(
        func(queries *db.Queries) Repository {
            return NewRepositoryPg(queries)
        },
        func(repo Repository, logger *zap.Logger, cfg *config.Config) *Service {
            return NewService(
                repo,
                logger,
                cfg.APIKeys.MaxKeysPerUser,
                cfg.APIKeys.DefaultExpiry,
            )
        },
    ),
)
```

The inline functions act as adapters: they extract the specific config values the constructor needs, keeping the constructor's signature clean (no dependency on `*config.Config`).

**Used by:** apikeys, user, session, notification, library

---

## Pattern 3: Complex (Multiple Providers + Inline)

Used when a module provides multiple types with interdependencies, or needs special construction logic.

```go
// internal/service/auth/module.go
var Module = fx.Module("auth",
    fx.Provide(
        // TokenManager from config
        func(cfg *config.Config) TokenManager {
            return NewTokenManager(cfg.Auth.JWTSecret, cfg.Auth.JWTExpiry)
        },
        // Repository from queries
        func(queries *db.Queries) Repository {
            return NewRepositoryPG(queries)
        },
        // Service with many dependencies
        func(pool *pgxpool.Pool, repo Repository, tm TokenManager,
            activityLogger activity.Logger, emailService *email.Service,
            cfg *config.Config) *Service {
            return NewService(
                pool, repo, tm, activityLogger, emailService,
                cfg.Auth.JWTExpiry, cfg.Auth.RefreshExpiry,
                cfg.Auth.LockoutThreshold, cfg.Auth.LockoutWindow,
                cfg.Auth.LockoutEnabled,
            )
        },
    ),
)
```

**Used by:** auth, mfa, oidc, rbac, metadata

---

## Pattern 4: With Lifecycle Hooks

Used for infrastructure that needs startup/shutdown behavior.

```go
// internal/infra/jobs/module.go
var Module = fx.Module("jobs",
    fx.Provide(
        NewRiverWorkers,   // func() *river.Workers
        NewRiverClient,    // func(*pgxpool.Pool, *river.Workers, *config.Config, *slog.Logger) (*Client, error)
    ),
    fx.Invoke(registerHooks),  // Lifecycle hooks
)

func registerHooks(lc fx.Lifecycle, client *Client, logger *slog.Logger) {
    lc.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            return client.Start(ctx)
        },
        OnStop: func(ctx context.Context) error {
            return client.Stop(ctx)
        },
    })
}
```

`fx.Invoke` runs after all providers are resolved. Lifecycle hooks ensure proper startup/shutdown order.

**Used by:** jobs, cache, api (HTTP server)

---

## Pattern 5: With CachedService

Used when a service has both a base service and a cached wrapper.

```go
// internal/service/rbac/module.go
var Module = fx.Module("rbac",
    fx.Provide(
        func(pool *pgxpool.Pool, cfg *config.Config) (*Service, error) {
            return NewService(pool, cfg.RBAC.ModelPath)
        },
        func(svc *Service, c *cache.Cache, logger *zap.Logger) *CachedService {
            return NewCachedService(svc, c, logger)
        },
    ),
)
```

Other modules that depend on RBAC receive `*CachedService` (which embeds `*Service`), getting caching transparently.

---

## fx.In / fx.Out Structs

For modules with many dependencies, use parameter/result structs:

```go
// internal/api/server.go
type ServerParams struct {
    fx.In

    Config          *config.Config
    Logger          *zap.Logger
    HealthService   *health.Service
    UserService     *user.Service
    AuthService     *auth.Service
    // ... 20+ dependencies ...
    MetadataService metadata.Service     `optional:"true"`
    RadarrService   *radarr.SyncService  `optional:"true"`
    RiverClient     *jobs.Client         `optional:"true"`
}

func NewServer(p ServerParams) (*Server, error) {
    // Use p.Config, p.Logger, etc.
}
```

The `optional:"true"` tag means fx won't fail if that dependency isn't provided. Used for optional integrations (Radarr, Sonarr, metadata) that may not be configured.

---

## Common Mistakes

1. **Circular dependencies**: fx detects these at startup. Fix by introducing an interface or restructuring.
2. **Missing `optional:"true"`**: If an integration module isn't loaded, fx fails on startup. Always mark optional deps.
3. **Wrong provide order**: Order in `fx.Module` doesn't matter — fx resolves the graph. But order in `app.Module` affects startup order for lifecycle hooks.
4. **Providing interface vs concrete**: If other modules depend on `Repository` (interface), your provider must return `Repository`, not `*pgRepository`. Inline functions handle this.
5. **Forgetting to add to app.Module**: New modules must be added to `internal/app/module.go` or they won't be loaded.
