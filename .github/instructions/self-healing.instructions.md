---
applyTo: "**/pkg/supervisor/**/*.go,**/pkg/graceful/**/*.go,**/cmd/**/*.go"
---

# Self-Healing & Graceful Shutdown Instructions

## Overview

- `pkg/supervisor` - Service supervision and automatic restart
- `pkg/graceful` - Graceful shutdown with ordered hooks

## Service Supervisor

### When to Use

- Background workers (library scanner, metadata fetcher)
- Long-running services that should restart on failure
- Services that need coordinated lifecycle

### Basic Pattern

```go
// Create supervisor
sup := supervisor.NewSupervisor(
    supervisor.DefaultSupervisorConfig("revenge"),
    logger,
)

// Add services
sup.Add(NewLibraryScanner(db, logger))
sup.Add(NewMetadataFetcher(providers, logger))
sup.Add(NewSearchIndexer(typesense, logger))

// Start supervision
if err := sup.Start(); err != nil {
    return err
}

// Later: stop all
sup.Stop()
```

### Service Interface

```go
type Service interface {
    Name() string
    Start(ctx context.Context) error // Blocks until stopped
    Stop(ctx context.Context) error  // Graceful stop
}

// Example implementation
type LibraryScanner struct {
    db     *pgxpool.Pool
    stopCh chan struct{}
    done   chan struct{}
}

func (s *LibraryScanner) Name() string { return "library-scanner" }

func (s *LibraryScanner) Start(ctx context.Context) error {
    s.done = make(chan struct{})
    defer close(s.done)

    ticker := time.NewTicker(time.Hour)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-s.stopCh:
            return nil
        case <-ticker.C:
            if err := s.scan(ctx); err != nil {
                return err // Will trigger restart
            }
        }
    }
}

func (s *LibraryScanner) Stop(ctx context.Context) error {
    close(s.stopCh)
    select {
    case <-s.done:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

### Function-Based Service

```go
sup.Add(supervisor.NewServiceFunc(
    "cleanup-worker",
    func(ctx context.Context) error {
        // Run cleanup loop
        for {
            select {
            case <-ctx.Done():
                return ctx.Err()
            case <-time.After(time.Hour):
                cleanup()
            }
        }
    },
    nil, // No special stop logic needed
))
```

### Supervision Strategies

```go
// OneForOne (default): Only restart the failed service
supervisor.StrategyOneForOne

// OneForAll: Restart ALL services if one fails
// Use when services are tightly coupled
supervisor.StrategyOneForAll

// RestForOne: Restart failed + services started after it
// Use for dependency chains
supervisor.StrategyRestForOne
```

### Restart Configuration

```go
supervisor.SupervisorConfig{
    MaxRestarts:      5,                  // Max restarts in window
    MaxRestartWindow: time.Minute,        // Time window
    RestartDelay:     100 * time.Millisecond, // Initial delay
    MaxRestartDelay:  30 * time.Second,   // Max delay (backoff)
    ShutdownTimeout:  30 * time.Second,   // Graceful stop timeout
}
```

### Health Check Integration

```go
// In health endpoint
if err := sup.HealthCheck(); err != nil {
    // At least one service is not running
    return health.StatusDegraded
}
```

### Manual Restart

```go
// Restart specific service (clears restart counter)
if err := sup.RestartService("metadata-fetcher"); err != nil {
    logger.Error("restart failed", "error", err)
}
```

## Graceful Shutdown

### Basic Pattern

```go
shutdowner := graceful.NewShutdowner(
    graceful.DefaultShutdownConfig(),
    logger,
)

// Register hooks in priority order (lower = earlier)
shutdowner.RegisterFunc("http-server", 0, stopHTTP)
shutdowner.RegisterFunc("supervisor", 10, stopSupervisor)
shutdowner.RegisterFunc("cache-flush", 20, flushCaches)
shutdowner.RegisterFunc("database", 30, closeDB)

// Start listening for SIGINT/SIGTERM
done := shutdowner.Start()

// Block until shutdown complete
<-done
```

### Priority Guidelines

| Priority | Component           | Example                  |
| -------- | ------------------- | ------------------------ |
| 0-9      | Stop accepting work | HTTP server, gRPC        |
| 10-19    | Drain in-flight     | Active streams, requests |
| 20-29    | Flush state         | Caches, indexes          |
| 30-39    | Close connections   | Database, Redis          |
| 40+      | Final cleanup       | Temp files               |

### HTTP Server Drain

```go
shutdowner.Register(graceful.DrainConnections(
    "http-server",
    httpServer,
    0,
))
```

### Programmatic Shutdown

```go
// Trigger from code (e.g., fatal error)
shutdowner.Trigger()

// Wait for completion
shutdowner.Wait()
```

### Context-Aware Wait Group

```go
wg := graceful.NewWaitGroupContext(ctx)

wg.Go(func(ctx context.Context) {
    processItem(ctx, item1)
})
wg.Go(func(ctx context.Context) {
    processItem(ctx, item2)
})

// Wait or timeout
if err := wg.WaitWithTimeout(30 * time.Second); err != nil {
    logger.Warn("workers did not finish in time")
}
```

## Integration Pattern

```go
func main() {
    ctx := context.Background()
    logger := setupLogger()

    // 1. Load config
    cfg := loadConfig()

    // 2. Setup shutdown handler
    shutdowner := graceful.NewShutdowner(
        graceful.DefaultShutdownConfig(),
        logger,
    )

    // 3. Initialize database
    db := setupDatabase(ctx, cfg)
    shutdowner.RegisterFunc("database", 30, func(ctx context.Context) error {
        db.Close()
        return nil
    })

    // 4. Setup supervisor
    sup := supervisor.NewSupervisor(
        supervisor.DefaultSupervisorConfig("revenge"),
        logger,
    )
    sup.Add(NewScanner(db))
    sup.Add(NewIndexer(db))

    shutdowner.RegisterFunc("supervisor", 10, func(ctx context.Context) error {
        return sup.Stop()
    })

    // 5. Start HTTP server
    server := setupHTTPServer(cfg, db, sup)
    shutdowner.Register(graceful.DrainConnections("http", server, 0))

    go server.ListenAndServe()

    // 6. Start supervisor
    sup.Start()

    // 7. Wait for shutdown
    done := shutdowner.Start()
    <-done
}
```

## DO's

- ✅ Use supervisor for background workers
- ✅ Implement proper Stop() with timeout
- ✅ Use priority-ordered shutdown hooks
- ✅ Handle SIGINT/SIGTERM gracefully
- ✅ Log service state changes

## DON'Ts

- ❌ Use os.Exit() directly (skip cleanup)
- ❌ Forget shutdown timeout
- ❌ Block indefinitely in Stop()
- ❌ Ignore context cancellation in workers
- ❌ Use OneForAll without good reason
---

## Related

- [INDEX.instructions.md](INDEX.instructions.md) - Main instruction index with all cross-references
- [ARCHITECTURE_V2.md](../../docs/dev/design/architecture/ARCHITECTURE_V2.md) - System architecture
- [BEST_PRACTICES.md](../../docs/dev/design/operations/BEST_PRACTICES.md) - Best practices
