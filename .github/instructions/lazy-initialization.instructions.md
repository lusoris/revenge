---
applyTo: "**/pkg/lazy/**/*.go"
---

# Lazy Initialization Patterns

> Initialize services only when first needed, not at startup.

## When to Use Lazy Initialization

### Good Candidates (Cold Start)

| Service            | Why                          |
| ------------------ | ---------------------------- |
| Transcoder Client  | Only needed for playback     |
| Metadata Providers | Only needed for library scan |
| Email Service      | Only for notifications       |
| OIDC Providers     | Only for SSO login           |
| Search Client      | Only for search queries      |

### Poor Candidates (Always Hot)

| Service         | Why                  |
| --------------- | -------------------- |
| HTTP Server     | Handles all requests |
| Auth Middleware | Every request        |
| Session Cache   | Session validation   |
| Database Pool   | Critical path        |

## Basic Usage

```go
// Good: Create lazy service
var transcoder = lazy.New(func() (*TranscoderClient, error) {
    return NewTranscoderClient(config), nil
})

// Good: Use when needed
func HandlePlayback(w http.ResponseWriter, r *http.Request) {
    client, err := transcoder.Get() // Init on first call
    if err != nil {
        http.Error(w, "unavailable", http.StatusServiceUnavailable)
        return
    }
    // Use client...
}
```

## With Cleanup

```go
// Good: Register cleanup for graceful shutdown
var searchClient = lazy.NewWithCleanup(
    func() (*SearchClient, error) {
        return NewSearchClient(config)
    },
    func(c *SearchClient) error {
        return c.Close()
    },
)

// On shutdown
func Shutdown() {
    searchClient.Close() // Only closes if initialized
}
```

## fx Integration

```go
// Good: Provide lazy wrapper
func ProvideLazyTranscoder(config Config, logger *slog.Logger) *lazy.Service[*TranscoderClient] {
    return lazy.New(func() (*TranscoderClient, error) {
        logger.Info("initializing transcoder (lazy)")
        return NewTranscoderClient(config), nil
    })
}

// In module
var Module = fx.Module("playback",
    fx.Provide(ProvideLazyTranscoder),
)

// In handler
type Handler struct {
    transcoder *lazy.Service[*TranscoderClient]
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    client, err := h.transcoder.Get()
    // ...
}
```

## Monitoring

```go
// Good: Track initialization
func handleRequest() {
    wasInit := transcoder.IsInitialized()

    client, err := transcoder.Get()

    if !wasInit && transcoder.IsInitialized() {
        logger.Info("lazy init triggered",
            "service", "transcoder",
            "duration", transcoder.InitTime())
    }
}
```

## Error Handling

```go
// Good: Handle init errors gracefully
client, err := lazyService.Get()
if err != nil {
    if errors.Is(err, ErrConfigMissing) {
        // Config error - won't recover
        http.Error(w, "service not configured", http.StatusServiceUnavailable)
        return
    }
    // Transient error - may recover on retry
    http.Error(w, "service temporarily unavailable", http.StatusServiceUnavailable)
    return
}
```

## DO's and DON'Ts

### DO

```go
// ✅ Use for non-critical services
var emailClient = lazy.New(NewEmailClient)

// ✅ Check initialization state
if transcoder.IsInitialized() {
    // Safe to assume fast response
}

// ✅ Handle errors gracefully
client, err := lazy.Get()
if err != nil {
    return fallbackBehavior()
}

// ✅ Log initialization
logger.Info("service initialized", "duration", lazy.InitTime())
```

### DON'T

```go
// ❌ Use for critical path services
var authMiddleware = lazy.New(NewAuth) // Auth is always needed!

// ❌ Panic on errors
client := lazyService.MustGet() // Only if truly can't fail

// ❌ Assume instant availability
start := time.Now()
client, _ := transcoder.Get() // May take 100ms+
// Don't measure this as request latency

// ❌ Forget cleanup
var conn = lazy.New(func() (*sql.DB, error) {
    return sql.Open("postgres", dsn)
})
// Connection never closed!
```
---

## Related

- [INDEX.instructions.md](INDEX.instructions.md) - Main instruction index with all cross-references
- [ARCHITECTURE_V2.md](../../docs/dev/design/architecture/ARCHITECTURE_V2.md) - System architecture
- [BEST_PRACTICES.md](../../docs/dev/design/operations/BEST_PRACTICES.md) - Best practices
