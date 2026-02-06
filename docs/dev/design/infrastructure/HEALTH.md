# Health Checks

<!-- DESIGN: infrastructure -->

**Package**: `internal/infra/health`
**fx Module**: `health.Module`

> Kubernetes liveness, readiness, and startup probes with dependency checks

---

## Service Structure

```
internal/infra/health/
├── checks.go              # Independent health check functions
├── service.go             # Service with K8s probe logic
├── handler.go             # HTTP endpoint handlers
└── module.go              # fx module
```

## Health Service

```go
type Service struct { /* RWMutex for startup state */ }

func (s *Service) Liveness(ctx context.Context) *CheckResult
func (s *Service) Readiness(ctx context.Context) *CheckResult
func (s *Service) Startup(ctx context.Context) *CheckResult
func (s *Service) FullCheck(ctx context.Context) []CheckResult
func (s *Service) MarkStartupComplete()
```

**Probe behavior**:
- **Liveness**: Always returns healthy (process alive check)
- **Readiness**: Database required, cache/jobs optional (degraded OK)
- **Startup**: Tracks initialization phase, returns unhealthy until `MarkStartupComplete()` called

## Dependency Checks

```go
func CheckDatabase(ctx context.Context, pool *pgxpool.Pool) CheckResult
func CheckCache(ctx context.Context, client *cache.Client) CheckResult
func CheckJobs(ctx context.Context, client *jobs.Client) CheckResult
func CheckAll(ctx context.Context, pool, cache, jobs) []CheckResult
```

**Status hierarchy**: `unhealthy` > `degraded` > `healthy`

Optional services (cache, jobs) return `degraded` when unavailable but don't fail readiness.

## HTTP Endpoints

```go
type Handler struct { service *Service }

func (h *Handler) HandleLiveness(w, r)   // GET /health/live
func (h *Handler) HandleReadiness(w, r)  // GET /health/ready
func (h *Handler) HandleStartup(w, r)    // GET /health/startup
func (h *Handler) HandleFull(w, r)       // GET /health/full
func (h *Handler) RegisterRoutes(mux)
```

**Response**: JSON with `status`, `message`, and optional `details` map.

## Key Types

```go
type Status string  // "healthy", "unhealthy", "degraded"

type CheckResult struct {
    Name    string
    Status  Status
    Message string
    Details map[string]interface{}  // Pool stats, connection info, etc.
}
```

## Related Documentation

- [OBSERVABILITY.md](OBSERVABILITY.md) - Metrics and monitoring
- [DATABASE.md](DATABASE.md) - Pool health check details
