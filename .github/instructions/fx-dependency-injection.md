# Dependency Injection with uber-go/fx

> Quick reference for fx v1.24+ patterns

## Basic App Structure

```go
func main() {
    fx.New(
        // Module grouping
        fx.Module("core",
            fx.Provide(NewConfig, NewLogger),
        ),
        fx.Module("database",
            fx.Provide(NewDatabase),
        ),
        fx.Module("api",
            fx.Provide(NewHTTPServer),
            fx.Invoke(RegisterRoutes),
        ),
    ).Run()
}
```

## Provide vs Supply vs Invoke

```go
// Provide: register constructor (lazy, only called if needed)
fx.Provide(NewDatabase)

// Supply: provide existing value (for testing, config)
fx.Supply(existingLogger)

// Invoke: execute function at startup (eager)
fx.Invoke(StartMetricsServer)
```

## Parameter Structs (fx.In)

```go
// Many dependencies? Use struct
type ServerParams struct {
    fx.In

    Config   *Config
    Logger   *slog.Logger
    DB       *pgxpool.Pool
    Cache    *redis.Client `optional:"true"` // Optional!
    Handlers []Handler     `group:"handlers"` // Value groups
}

func NewServer(p ServerParams) *Server {
    return &Server{
        config: p.Config,
        logger: p.Logger,
        // ...
    }
}
```

## Result Structs (fx.Out)

```go
// Return multiple types from constructor
type DatabaseResult struct {
    fx.Out

    Pool    *pgxpool.Pool
    Queries *db.Queries
}

func NewDatabase(cfg *Config) (DatabaseResult, error) {
    pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
    if err != nil {
        return DatabaseResult{}, err
    }
    return DatabaseResult{
        Pool:    pool,
        Queries: db.New(pool),
    }, nil
}
```

## Named Dependencies

```go
// Provide named instances
type Connections struct {
    fx.Out

    ReadConn  *sql.DB `name:"read"`
    WriteConn *sql.DB `name:"write"`
}

// Consume named instances
type ServiceParams struct {
    fx.In

    ReadDB  *sql.DB `name:"read"`
    WriteDB *sql.DB `name:"write"`
}
```

## Value Groups

```go
// Provide to group
type HandlerResult struct {
    fx.Out
    Handler http.Handler `group:"handlers"`
}

func NewUserHandler(...) HandlerResult {
    return HandlerResult{Handler: userHandler}
}

func NewMediaHandler(...) HandlerResult {
    return HandlerResult{Handler: mediaHandler}
}

// Consume group
type RouterParams struct {
    fx.In
    Handlers []http.Handler `group:"handlers"`
}

func NewRouter(p RouterParams) *http.ServeMux {
    mux := http.NewServeMux()
    for _, h := range p.Handlers {
        // register handlers
    }
    return mux
}
```

## Lifecycle Hooks

```go
func NewServer(lc fx.Lifecycle, cfg *Config) *http.Server {
    srv := &http.Server{Addr: cfg.Addr}

    lc.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            ln, err := net.Listen("tcp", srv.Addr)
            if err != nil {
                return err
            }
            go srv.Serve(ln)
            return nil
        },
        OnStop: func(ctx context.Context) error {
            return srv.Shutdown(ctx)
        },
    })

    return srv
}

// Alternative: StartHook/StopHook helpers
lc.Append(fx.StartStopHook(
    func() { slog.Info("starting") },
    func() { slog.Info("stopping") },
))
```

## Annotate (without structs)

```go
fx.Provide(
    fx.Annotate(
        NewHandler,
        fx.ParamTags(`name:"read"`, ``),      // param annotations
        fx.ResultTags(`group:"handlers"`),     // result annotations
        fx.As(new(Handler)),                   // provide as interface
        fx.OnStart(func(h *MyHandler) error {  // lifecycle
            return h.Init()
        }),
        fx.OnStop(func(h *MyHandler) error {
            return h.Shutdown()
        }),
    ),
)
```

## Private Dependencies

```go
fx.Module("internal",
    fx.Provide(
        fx.Private,  // Only visible within this module
        NewInternalService,
    ),
)
```

## Decorators

```go
// Wrap/modify dependencies
fx.Decorate(func(log *slog.Logger) *slog.Logger {
    return log.With("module", "api")
})
```

## Testing with fxtest

```go
import "go.uber.org/fx/fxtest"

func TestServer(t *testing.T) {
    var server *Server

    app := fxtest.New(t,
        fx.Provide(NewConfig, NewServer),
        fx.Populate(&server),  // Extract for testing
    )
    app.RequireStart()
    defer app.RequireStop()

    // Test server...
}
```

## Error Handling

```go
// Constructors can return errors
func NewDatabase(cfg *Config) (*Database, error) {
    if cfg.URL == "" {
        return nil, errors.New("database URL required")
    }
    // ...
}

// Shutdown can provide exit code
fx.Invoke(func(shutdowner fx.Shutdowner) {
    // Later...
    shutdowner.Shutdown(fx.ExitCode(1))
})
```

## Module Pattern

```go
// pkg/database/module.go
package database

var Module = fx.Module("database",
    fx.Provide(
        NewPool,
        NewQueries,
    ),
)

// cmd/app/main.go
fx.New(
    database.Module,
    api.Module,
    // ...
)
```

## Recommended Structure

```
internal/
├── api/
│   ├── module.go      # fx.Module("api", ...)
│   └── handlers/
├── service/
│   └── module.go
├── infra/
│   ├── database/
│   │   └── module.go
│   └── cache/
│       └── module.go
└── app/
    └── app.go         # fx.New() combining all modules
```
