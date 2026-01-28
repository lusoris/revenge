---
applyTo: "**/pkg/health/**/*.go,**/internal/api/handlers/health*.go"
---

# Health Check Patterns

> Implement proper health checking for service reliability.

## Health Check Categories

| Category   | Requirement       | Example Services |
| ---------- | ----------------- | ---------------- |
| `critical` | Must be healthy   | Database, Auth   |
| `warm`     | Should be healthy | Cache, Search    |
| `cold`     | Can be unhealthy  | Email, OIDC      |

## Registration

```go
// Good: Register checks at startup
func RegisterHealthChecks(checker *health.Checker, db *pgxpool.Pool, cache *redis.Client) {
    // Critical - failure = unhealthy
    checker.Register(health.Check{
        Name:     "database",
        Category: health.CategoryCritical,
        Timeout:  5 * time.Second,
        Check: func(ctx context.Context) error {
            return db.Ping(ctx)
        },
    })

    // Warm - failure = degraded
    checker.Register(health.Check{
        Name:     "cache",
        Category: health.CategoryWarm,
        Timeout:  2 * time.Second,
        Check: func(ctx context.Context) error {
            return cache.Ping(ctx).Err()
        },
    })

    // Cold - failure = still okay
    checker.RegisterFunc("email", health.CategoryCold, func(ctx context.Context) error {
        return emailClient.Ping(ctx)
    })
}
```

## HTTP Endpoints

```go
// Good: Separate liveness and readiness
func SetupHealthRoutes(mux *http.ServeMux, checker *health.Checker) {
    // Liveness - is the process alive?
    mux.HandleFunc("GET /health/live", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("ok"))
    })

    // Readiness - can we serve traffic?
    mux.HandleFunc("GET /health/ready", func(w http.ResponseWriter, r *http.Request) {
        if checker.IsReady(r.Context()) {
            w.WriteHeader(http.StatusOK)
        } else {
            w.WriteHeader(http.StatusServiceUnavailable)
        }

        status := checker.Check(r.Context())
        json.NewEncoder(w).Encode(status)
    })

    // Detailed health - for monitoring
    mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
        status := checker.Check(r.Context())

        code := http.StatusOK
        if status.Status == health.StatusUnhealthy {
            code = http.StatusServiceUnavailable
        }

        w.WriteHeader(code)
        json.NewEncoder(w).Encode(status)
    })
}
```

## Response Format

```json
{
  "status": "healthy",
  "services": {
    "database": {
      "name": "database",
      "healthy": true,
      "category": "critical",
      "latency_ms": 2
    },
    "cache": {
      "name": "cache",
      "healthy": true,
      "category": "warm",
      "latency_ms": 1
    },
    "email": {
      "name": "email",
      "healthy": false,
      "category": "cold",
      "latency_ms": 5000,
      "error": "connection refused"
    }
  },
  "checked_at": "2026-01-28T12:00:00Z"
}
```

## Kubernetes Integration

```yaml
# Good: Proper probe configuration
livenessProbe:
  httpGet:
    path: /health/live
    port: 8096
  initialDelaySeconds: 5
  periodSeconds: 10
  failureThreshold: 3

readinessProbe:
  httpGet:
    path: /health/ready
    port: 8096
  initialDelaySeconds: 10
  periodSeconds: 5
  failureThreshold: 2

startupProbe:
  httpGet:
    path: /health/ready
    port: 8096
  initialDelaySeconds: 5
  periodSeconds: 5
  failureThreshold: 30 # 2.5 minutes to start
```

## Check Implementation

```go
// Good: Meaningful health check with timeout
func (s *SearchService) HealthCheck(ctx context.Context) error {
    ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
    defer cancel()

    // Actually test functionality, not just connection
    _, err := s.client.Health(ctx, 2*time.Second)
    return err
}

// Good: Check lazy services
func (s *TranscoderService) HealthCheck(ctx context.Context) error {
    if !s.lazy.IsInitialized() {
        return nil // Not initialized = not unhealthy
    }

    client, err := s.lazy.Get()
    if err != nil {
        return err
    }

    return client.Ping(ctx)
}
```

## Caching

```go
// Good: Cache results to avoid overload
type Checker struct {
    cacheDuration time.Duration // 5 seconds
    lastStatus    HealthStatus
    lastCheckTime time.Time
}

func (c *Checker) Check(ctx context.Context) HealthStatus {
    if time.Since(c.lastCheckTime) < c.cacheDuration {
        return c.lastStatus
    }

    // Perform actual checks...
}
```

## DO's and DON'Ts

### DO

```go
// ✅ Set appropriate timeouts
check.Timeout = 5 * time.Second

// ✅ Use correct categories
checker.Register(health.Check{
    Name:     "optional-service",
    Category: health.CategoryCold, // Won't fail overall health
})

// ✅ Check actual functionality
func dbHealthCheck(ctx context.Context) error {
    return db.QueryRowContext(ctx, "SELECT 1").Err()
}

// ✅ Handle gracefully in handlers
if !checker.IsReady(ctx) {
    // Return 503, not panic
}
```

### DON'T

```go
// ❌ Make all checks critical
checker.Register(health.Check{
    Name:     "email",
    Category: health.CategoryCritical, // Email down = all down?
})

// ❌ No timeout
func dbHealthCheck(ctx context.Context) error {
    return db.Ping(ctx) // May hang forever
}

// ❌ Check too frequently
cacheDuration: 100 * time.Millisecond // Overload on /health spam

// ❌ Expose sensitive info
return fmt.Errorf("connection failed: %s", password)
```
---

## Related

- [INDEX.instructions.md](INDEX.instructions.md) - Main instruction index with all cross-references
- [ARCHITECTURE_V2.md](../../docs/dev/design/architecture/ARCHITECTURE_V2.md) - System architecture
- [BEST_PRACTICES.md](../../docs/dev/design/operations/BEST_PRACTICES.md) - Best practices
