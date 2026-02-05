# Revenge Codebase: Cluster Deployment Readiness Assessment

**Generated**: 2026-02-05
**Target Platforms**: Kubernetes, k3s, Docker Swarm
**Overall Readiness**: 84% (Near Production-Ready)

---

## Executive Summary

The Revenge codebase is **84% ready for cluster deployment** with **excellent state management**, **12-factor compliance**, and **production-grade observability**. The primary blockers are:
1. **Media file storage** (local filesystem not cluster-compatible)
2. **Avatar storage** (needs S3 or shared volume)
3. **Leader election** for periodic jobs (not implemented)

With these three issues resolved, the application will be fully ready for horizontal scaling across Kubernetes, k3s, and Docker Swarm.

---

## Cluster Readiness Scorecard

| Category | Score | Status | Notes |
|----------|-------|--------|-------|
| **State Management** | 90% | 🟡 | Sessions/cache ✅, media storage ⚠️ |
| **Configuration** | 100% | ✅ | 12-factor compliant |
| **Service Discovery** | 100% | ✅ | Health probes excellent |
| **Distributed Concerns** | 70% | 🟡 | Jobs ✅, leader election 🔴 |
| **Observability** | 85% | 🟡 | Metrics ✅, tracing partial |
| **Storage** | 60% | ⚠️ | DB ✅, media/avatars need work |
| **Overall** | **84%** | 🟡 | **Near production-ready** |

---

## 1. State Management ✅ EXCELLENT

### Externalized State (Production-Ready)

**PostgreSQL 18+**: All persistent state stored in PostgreSQL ✅
- User accounts, sessions, movies, libraries, metadata
- Connection pooling configured (`pgxpool`)
- Self-healing pool with health checks
- Location: `internal\infra\database\`

**Dragonfly (Redis-compatible)**: Distributed cache ✅
- Session caching with L2 distributed layer
- Rate limiting state (Redis-backed with fallback)
- Cache invalidation handled correctly across instances
- Location: `internal\infra\cache\`

**Configuration**: `internal\config\config.go:113-132`
```yaml
database:
  url: postgres://...
  max_conns: 0  # Auto: (CPU * 2) + 1
  min_conns: 2
  max_conn_lifetime: 30m
  max_conn_idle_time: 5m
  health_check_period: 30s

cache:
  url: dragonfly:6379
  password: ${CACHE_PASSWORD}
```

---

### No Problematic In-Memory State ✅

**L1 Cache (Otter)**: Local process cache - Cluster-safe ✅
- Read-through pattern
- Cache misses fall back to L2 (Dragonfly) then PostgreSQL
- Invalidation via L2 ensures consistency
- No cross-pod state dependency

**Session Handling**: Distributed ✅
- Sessions stored in PostgreSQL (`shared.sessions` table)
- Optional L2 cache (Dragonfly) for performance
- `CachedService` wraps session service with distributed caching
- Session validation: `TokenHash` lookup → Cache or DB
- No local-only session storage

**Location**: `internal\service\session\cached_service.go:119-163`
```go
// Session revocation invalidates L2 cache
s.cache.Delete(ctx, cacheKey)
s.cache.InvalidateUserSessions(ctx, userID)
```

---

### File Upload Handling ⚠️ NEEDS WORK

**Current State**: Local filesystem storage
- **Location**: `internal\service\storage\storage.go:34-55`
- Implementation: `LocalStorage` uses filesystem at `avatar.storage_path`
- Mutex-protected (`sync.RWMutex`) for local consistency
- File paths: `/data/avatars/{userId}/{uuid}.ext`

**Problem**:
- Local filesystem storage won't work in multi-pod deployments
- Each pod would have isolated storage
- No shared volume configuration in Helm charts

**Solution Required** (Choose one):

#### Option 1: S3-Compatible Storage (Recommended)
```go
// Interface already exists!
type Storage interface {
    Store(ctx context.Context, key string, reader io.Reader, contentType string) (string, error)
    Get(ctx context.Context, key string) (io.ReadCloser, error)
    Delete(ctx context.Context, key string) error
}

// Implement S3Storage
type S3Storage struct {
    client *s3.Client
    bucket string
}
```

**Configuration**:
```yaml
avatar:
  storage_backend: s3  # or "local"
  s3:
    endpoint: minio:9000
    bucket: revenge-avatars
    access_key_id: ${S3_ACCESS_KEY}
    secret_access_key: ${S3_SECRET_KEY}
```

#### Option 2: Shared Volume (Less optimal)
```yaml
# Kubernetes PVC
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: revenge-avatars
spec:
  accessModes:
    - ReadWriteMany  # Requires NFS/CephFS
  resources:
    requests:
      storage: 10Gi
```

---

### Media Files 🔴 CRITICAL ISSUE

**Library paths are local filesystem**:
```yaml
movie:
  library:
    paths: ["/media/movies"]
```

**Problem for Clustering**:
- Media scanning requires all pods to access the same files
- Current: Docker volume mount `./media:/media:ro` (read-only)
- Multi-pod: Needs shared storage

**Solutions**:

#### Option 1: NFS/SMB Mount (Recommended)
```yaml
# Kubernetes Volume
apiVersion: v1
kind: PersistentVolume
metadata:
  name: media-storage
spec:
  capacity:
    storage: 1Ti
  accessModes:
    - ReadOnlyMany  # Multiple pods can read
  nfs:
    server: nas.local
    path: /volume1/media
```

**Helm Chart**:
```yaml
volumes:
  - name: media
    nfs:
      server: {{ .Values.media.nfs.server }}
      path: {{ .Values.media.nfs.path }}
      readOnly: true

volumeMounts:
  - name: media
    mountPath: /media
    readOnly: true
```

#### Option 2: Object Storage Adapter
Adapt library scanner to use S3-compatible storage:
- Requires significant refactoring of `library_scanner.go`
- Less performant than filesystem
- Not recommended for large media libraries

#### Option 3: Single-Pod Scanner with Affinity
Pin scan jobs to specific node with media access:
```yaml
# Kubernetes Job Affinity
affinity:
  nodeAffinity:
    requiredDuringSchedulingIgnoredDuringExecution:
      nodeSelectorTerms:
        - matchExpressions:
          - key: media-server
            operator: In
            values:
              - "true"
```

---

## 2. Configuration ✅ 12-FACTOR COMPLIANT

### Environment Variables (Production-Ready)

**Koanf-based config** with clear precedence:
1. Environment variables (`REVENGE_*` prefix)
2. Config file (`config/config.yaml`)
3. Defaults

**Location**: `internal\config\config.go`

**All critical config externalized**:
```bash
# Server
REVENGE_SERVER_PORT=8080
REVENGE_SERVER_HOST=0.0.0.0

# Database
REVENGE_DATABASE_URL=postgres://user:pass@postgres:5432/revenge
REVENGE_DATABASE_MAX_CONNS=0  # Auto

# Cache
REVENGE_CACHE_URL=dragonfly:6379
REVENGE_CACHE_PASSWORD=${CACHE_PASSWORD}

# Search
REVENGE_SEARCH_URL=http://typesense:8108
REVENGE_SEARCH_API_KEY=${TYPESENSE_API_KEY}

# JWT
REVENGE_JWT_SECRET=${JWT_SECRET}
REVENGE_JWT_EXPIRY=24h

# Logging
REVENGE_LOGGING_LEVEL=info
REVENGE_LOGGING_FORMAT=json
```

### Secrets Management ✅

- No hardcoded secrets
- Environment variable injection
- Ready for Kubernetes Secrets/ConfigMaps
- JWT secret, API keys, DB passwords all configurable

**Kubernetes Secret Example**:
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: revenge-secrets
type: Opaque
stringData:
  database-url: postgres://...
  jwt-secret: ...
  cache-password: ...
  typesense-api-key: ...
```

### Configuration Files ✅

- Optional YAML file (`config.example.yaml`)
- Hot-reload supported via `koanf.Watch()` (already implemented)
- Config validation using `go-playground/validator`

---

## 3. Service Discovery ✅ PRODUCTION READY

### Health Checks (Kubernetes-Ready)

**Location**: `internal\infra\health\service.go`

**Three probe types implemented**:

#### Liveness Probe (Always Healthy)
```go
// Line 64-72: Always returns healthy
func (s *Service) Liveness(ctx context.Context) CheckResult {
    return CheckResult{Healthy: true}
}
```

**Endpoint**: `GET /health/live`

#### Readiness Probe (Checks Dependencies)
```go
// Line 76-110: Checks startup + dependencies
- Startup complete check
- Database health check
- Returns 503 if not ready
```

**Endpoint**: `GET /health/ready`

#### Startup Probe (Initialization Status)
```go
// Line 114-132: Initialization status
- MarkStartupComplete() called on boot (line 55-62)
- Prevents traffic before ready
```

**Endpoint**: `GET /health/startup`

**Helm Chart Configuration**:
```yaml
livenessProbe:
  httpGet:
    path: /health/live
    port: http
  initialDelaySeconds: 10
  periodSeconds: 30
  timeoutSeconds: 5
  failureThreshold: 3

readinessProbe:
  httpGet:
    path: /health/ready
    port: http
  initialDelaySeconds: 5
  periodSeconds: 10
  timeoutSeconds: 3
  failureThreshold: 3

startupProbe:
  httpGet:
    path: /health/startup
    port: http
  initialDelaySeconds: 0
  periodSeconds: 5
  failureThreshold: 30  # 150s total
```

---

### Graceful Shutdown ✅

**Signal handling**:
```go
// cmd/revenge/main.go:36-37
ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

// Graceful stop with timeout (line 53-59)
stopCtx, stopCancel := context.WithTimeout(context.Background(), application.StopTimeout())
application.Stop(stopCtx)
```

**fx lifecycle integration**:
- All services register OnStop hooks
- River job queue graceful shutdown
- Database pool draining
- HTTP server graceful shutdown configured (`ShutdownTimeout: 10s`)

**Kubernetes terminationGracePeriodSeconds**:
```yaml
spec:
  terminationGracePeriodSeconds: 30  # Match app timeout
```

---

## 4. Distributed Concerns ⚠️ PARTIAL

### Leader Election 🔴 NOT IMPLEMENTED

**Current State**:
- No leader election mechanism
- Multiple pods would execute same periodic jobs

**Evidence**:
- River job queue is shared (PostgreSQL-backed) ✅
- But no job deduplication for periodic tasks
- No `PeriodicJob` scheduling found in codebase

**Periodic Tasks at Risk**:
```go
// internal/infra/jobs/cleanup_job.go - CleanupWorker
// internal/service/activity/cleanup.go - Activity log cleanup
```

**Solutions Needed** (Choose one):

#### Option 1: River Periodic Jobs (Recommended)
```go
// River supports periodic jobs with unique keys
periodicJobs := []*river.PeriodicJob{
    river.NewPeriodicJob(
        river.PeriodicInterval(1 * time.Hour),
        func() (river.JobArgs, *river.InsertOpts) {
            return ActivityCleanupArgs{}, &river.InsertOpts{
                UniqueOpts: river.UniqueOpts{
                    ByArgs: true,  // Prevents duplicate scheduling
                },
            }
        },
    ),
}
```

#### Option 2: HashiCorp Raft (Already in dependencies!)
```go
// go.mod line 202: github.com/hashicorp/raft v1.8.0

type LeaderElection struct {
    raft *raft.Raft
}

func (le *LeaderElection) IsLeader() bool {
    return le.raft.State() == raft.Leader
}

// Only leader runs periodic jobs
if le.IsLeader() {
    runPeriodicJob()
}
```

#### Option 3: Kubernetes Lease API
```go
import (
    "k8s.io/client-go/kubernetes"
    coordinationv1 "k8s.io/api/coordination/v1"
)

// Use Kubernetes Lease for leader election
leaseClient := clientset.CoordinationV1().Leases("default")
```

#### Option 4: PostgreSQL Advisory Locks
```go
// Use PostgreSQL for leader election
func (s *Service) tryAcquireLeaderLock(ctx context.Context) bool {
    var acquired bool
    err := s.db.QueryRow(ctx, "SELECT pg_try_advisory_lock(12345)").Scan(&acquired)
    return err == nil && acquired
}
```

---

### Distributed Locks ⚠️ PARTIAL

**PostgreSQL Advisory Locks Available**:
- Not currently used
- Could implement for critical sections
- No evidence of race conditions in current code

**River Job Queue** ✅
- Handles job deduplication
- Transactional job enqueueing
- Multi-worker safe

---

### Job Queue in Multi-Instance ✅ READY

**River Implementation**: `internal\infra\jobs\`

**PostgreSQL-backed coordination**:
```go
// Uses PostgreSQL for coordination
// Each pod runs River workers
// Jobs distributed across all workers
// No single point of failure
```

**Configuration**:
```yaml
workers:
  max: 100  # Per pod
  fetch_cooldown: 200ms
  poll_interval: 2s
  rescue_stuck_jobs: 30m
```

**Multi-pod behavior**:
- Jobs dequeued atomically via PostgreSQL
- No duplicate processing
- Scales horizontally
- Load balances automatically

---

### Cache Invalidation ✅ SOLVED

**Distributed cache invalidation**:
```go
// internal/service/session/cached_service.go:119-163

// Session revocation invalidates L2 cache
s.cache.Delete(ctx, cacheKey)
s.cache.InvalidateUserSessions(ctx, userID)
```

**Pattern**:
- Write-through caching
- Invalidation via Dragonfly (shared across pods)
- L1 cache (Otter) has short TTL, acceptable staleness

---

## 5. Observability ✅ PRODUCTION READY

### Structured Logging ✅

**Implementation**:
- Development: `tint` (slog handler, colorized)
- Production: `zap` (high-performance JSON)
- Configuration: `logging.format: json|text`

**Location**: `internal\infra\logging\`

**Extensive coverage**:
- 22 files use `slog.Logger` (51 occurrences)
- Contextual logging throughout
- Log levels: debug, info, warn, error

**Configuration**:
```yaml
logging:
  level: info  # debug, info, warn, error
  format: json  # json, text
```

---

### Metrics ✅ PROMETHEUS-READY

**Location**: `internal\infra\observability\metrics.go`

**Exposed Metrics**:

#### HTTP Metrics
```
revenge_http_requests_total{method, path, status}
revenge_http_request_duration_seconds{method, path}
revenge_http_requests_in_flight
```

#### Cache Metrics
```
revenge_cache_hits_total{cache, layer}
revenge_cache_misses_total{cache, layer}
revenge_cache_operation_duration_seconds
```

#### Job Metrics
```
revenge_jobs_enqueued_total{job_type}
revenge_jobs_completed_total{job_type, status}
revenge_jobs_duration_seconds{job_type}
```

#### Session Metrics
```
revenge_sessions_active_total
```

**Middleware**:
- HTTP metrics middleware for all requests (line 14-44)
- Path normalization to avoid high cardinality (line 47-95)

**Prometheus Scrape Config**:
```yaml
scrape_configs:
  - job_name: 'revenge'
    kubernetes_sd_configs:
      - role: pod
    relabel_configs:
      - source_labels: [__meta_kubernetes_pod_label_app]
        action: keep
        regex: revenge
    metrics_path: /metrics
```

---

### Tracing Support 🟡 PARTIAL

**OpenTelemetry imported**:
```go
// go.mod: go.opentelemetry.io/otel v1.39.0
```

**Not yet instrumented**:
- No trace context propagation found
- No span creation in handlers
- Easy to add later

**Recommendation**:
```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/trace"
)

func (h *Handler) GetMovie(ctx context.Context, id uuid.UUID) (*Movie, error) {
    ctx, span := otel.Tracer("movie").Start(ctx, "GetMovie")
    defer span.End()

    movie, err := h.service.GetMovie(ctx, id)
    if err != nil {
        span.RecordError(err)
        return nil, err
    }

    return movie, nil
}
```

---

### Request Correlation IDs 🔴 MISSING

**Not implemented**:
- No `X-Request-ID` generation
- No correlation ID in logs
- No trace context

**Recommendation**:
Add middleware to:
1. Generate/extract request ID from `X-Request-ID` header
2. Inject into context
3. Include in all logs
4. Return in response header

**Implementation**:
```go
// internal/api/middleware/request_id.go
func RequestID() func(next http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            requestID := r.Header.Get("X-Request-ID")
            if requestID == "" {
                requestID = uuid.New().String()
            }

            ctx := context.WithValue(r.Context(), "request_id", requestID)
            w.Header().Set("X-Request-ID", requestID)

            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

---

## 6. Storage ⚠️ NEEDS ATTENTION

### Media Files 🔴 BLOCKING ISSUE

**Current**: Local filesystem
```yaml
movie:
  library:
    paths: ["/path/to/movies"]
```

**Cluster Requirements**:
1. **ReadOnlyMany volume** (NFS/SMB/CephFS)
2. All pods mount same media storage
3. Scanner jobs need consistent view

**Helm Chart Missing**:
```yaml
# Add to values.yaml
media:
  persistence:
    enabled: true
    storageClass: nfs-client
    accessMode: ReadOnlyMany
    size: 1Ti
    nfs:
      server: nas.local
      path: /volume1/media
```

---

### Avatar Storage 🟡 NEEDS UPGRADE

**Current**: LocalStorage (filesystem)
```yaml
avatar:
  storage_path: /data/avatars
```

**Interface exists for S3**:
```go
type Storage interface {
    Store(ctx context.Context, key string, reader io.Reader, contentType string) (string, error)
    Get(ctx context.Context, key string) (io.ReadCloser, error)
    Delete(ctx context.Context, key string) error
}
```

**Solutions**:
1. Implement S3Storage backend (MinIO/S3/GCS)
2. Or use RWX volume (less optimal)

---

### Database Connection Pooling ✅

**Excellent implementation**:
```go
// config/config.go:113-132
max_conns: 0  # Auto: (CPU * 2) + 1
min_conns: 2
max_conn_lifetime: 30m
max_conn_idle_time: 5m
health_check_period: 30s
```

**Multi-pod safe**:
- Each pod has own connection pool
- PostgreSQL handles total connections
- Graceful pool draining on shutdown

---

## 7. Rate Limiting ✅ CLUSTER-READY

### Redis-Backed Rate Limiter

**Location**: `internal\api\middleware\ratelimit_redis.go`

**Implementation**:
```go
// Sliding window algorithm via Lua script (line 162-186)
// Distributed across all pods via Dragonfly
// Automatic fallback to in-memory (line 238-240)
```

**Configuration**:
```yaml
rate_limit:
  enabled: true
  backend: redis  # or "memory"
  global:
    requests_per_second: 10
    burst: 20
  auth:
    requests_per_second: 1
    burst: 5
```

**Cluster behavior**:
- Redis backend: Shared state across pods ✅
- Memory fallback: Per-pod limits (degraded but functional)

---

## 8. Existing Kubernetes Support ✅

### Helm Chart

**Location**: `charts/revenge/`

**Features**:
- Deployment with configurable replicas (default: 2)
- HPA (Horizontal Pod Autoscaler):
  - Min: 2, Max: 10
  - Target CPU: 70%
- Service (ClusterIP)
- Health probes configured
- Resource limits/requests

**Missing**:
- Media volume mounts ⚠️
- Avatar storage PVC ⚠️
- Ingress (disabled by default)

---

### Docker Compose

**Production-ready**:
- PostgreSQL 18
- Dragonfly cache
- Typesense search
- Health checks on all services
- Volume management

---

## Required Changes for Full Cluster Support

### Priority 1: Blocking Issues

#### 1. Media File Storage
**Effort**: 4-8 hours

**Tasks**:
- Add NFS/CephFS PersistentVolume support to Helm chart
- Update deployment with media volume mounts
- Document shared storage setup
- Test multi-pod media scanning

**Helm Chart Changes**:
```yaml
# values.yaml
media:
  persistence:
    enabled: true
    storageClass: nfs-client
    accessMode: ReadOnlyMany
    size: 1Ti

# deployment.yaml
volumeMounts:
  - name: media
    mountPath: /media
    readOnly: true

volumes:
  - name: media
    persistentVolumeClaim:
      claimName: {{ .Values.media.persistence.existingClaim | default "revenge-media" }}
```

---

#### 2. Avatar Storage
**Effort**: 8-16 hours

**Tasks**:
- Implement S3Storage backend
- Add MinIO to docker-compose/helm
- Configuration for storage backend selection
- Migration script from local to S3

**Implementation**:
```go
// internal/service/storage/s3.go
type S3Storage struct {
    client *s3.Client
    bucket string
}

func (s *S3Storage) Store(ctx context.Context, key string, reader io.Reader, contentType string) (string, error) {
    _, err := s.client.PutObject(ctx, &s3.PutObjectInput{
        Bucket:      aws.String(s.bucket),
        Key:         aws.String(key),
        Body:        reader,
        ContentType: aws.String(contentType),
    })
    if err != nil {
        return "", err
    }

    return fmt.Sprintf("s3://%s/%s", s.bucket, key), nil
}
```

---

### Priority 2: Important Enhancements

#### 3. Leader Election for Periodic Jobs
**Effort**: 8-16 hours

**Tasks**:
- Implement River periodic jobs with unique constraints
- Or add HashiCorp Raft leader election
- Update cleanup jobs to only run on leader
- Test multi-pod periodic job execution

**River Periodic Jobs** (Recommended):
```go
// internal/infra/jobs/periodic.go
periodicJobs := []*river.PeriodicJob{
    river.NewPeriodicJob(
        river.PeriodicInterval(1 * time.Hour),
        func() (river.JobArgs, *river.InsertOpts) {
            return &ActivityCleanupArgs{}, &river.InsertOpts{
                Queue: QueueDefault,
                UniqueOpts: river.UniqueOpts{
                    ByPeriod: 1 * time.Hour,
                },
            }
        },
    ),
}
```

---

#### 4. Request Correlation IDs
**Effort**: 2-4 hours

**Tasks**:
- Add middleware for X-Request-ID
- Propagate through context
- Include in logs
- Return in response headers

---

### Priority 3: Nice to Have

#### 5. Distributed Tracing
**Effort**: 16-24 hours

**Tasks**:
- Wire up OpenTelemetry (already imported)
- Add Jaeger/Tempo support
- Instrument critical paths (handlers, services, repositories)
- Configure trace sampling

---

#### 6. Startup Probe in Helm
**Effort**: 1 hour

**Task**: Add to Helm chart (already implemented in code)
```yaml
startupProbe:
  httpGet:
    path: /health/startup
    port: http
  failureThreshold: 30
  periodSeconds: 10
```

---

## Deployment Scenarios

### ✅ What Works Today

1. **Multiple pod deployment**
   - Session sharing via PostgreSQL ✅
   - Distributed caching via Dragonfly ✅
   - Job distribution via River ✅
   - Load balancing across pods ✅

2. **Auto-scaling**
   - HPA configured for CPU ✅
   - Graceful shutdown prevents dropped requests ✅
   - Health probes ensure traffic only to ready pods ✅

3. **Database high availability**
   - Connection pooling per pod ✅
   - Self-healing on connection loss ✅
   - Prepared for PostgreSQL HA (Patroni/Stolon) ✅

---

### ⚠️ What Needs Work

1. **Media scanning with multiple pods**
   - Will fail without shared storage 🔴
   - Could work with job affinity to single node 🟡

2. **Avatar uploads in multi-pod**
   - Files stored on single pod 🔴
   - Other pods won't see uploaded files 🔴

3. **Periodic job deduplication**
   - Multiple cleanup jobs might run 🟡
   - Needs leader election 🟡

---

## Deployment Recommendation

### Current State: Ready for STAGING with Caveats

**Use Cases**:
- ✅ API-only workloads (no media scanning)
- ✅ Read-heavy workloads with shared media (NFS)
- ⚠️ Full media server (needs shared storage)

---

### For Production

**Required** (Blocking):
1. Implement media storage solution (NFS/S3) - 4-8h
2. Configure object storage for avatars - 8-16h

**Recommended** (Important):
3. Add leader election for periodic jobs - 8-16h
4. Add request correlation IDs - 2-4h

**Optional** (Enhancement):
5. Add distributed tracing - 16-24h
6. Load test with 5+ pods - 8h

**Total Effort**: 22-44 hours (3-5 days)

---

## Platform-Specific Notes

### Kubernetes / k3s

**Advantages**:
- Native PersistentVolume support
- Service discovery via DNS
- HPA for auto-scaling
- ConfigMaps and Secrets

**Configuration**:
```yaml
# Use existing Helm chart
helm install revenge ./charts/revenge \
  --set media.persistence.enabled=true \
  --set media.persistence.storageClass=nfs-client \
  --set replicaCount=3
```

---

### Docker Swarm

**Advantages**:
- Simpler than Kubernetes
- Built-in load balancing
- Stack deployment

**Configuration**:
```yaml
# docker-compose.swarm.yml
version: '3.8'
services:
  revenge:
    image: revenge:latest
    deploy:
      replicas: 3
      update_config:
        parallelism: 1
        delay: 10s
    volumes:
      - type: volume
        source: media
        target: /media
        read_only: true
        volume:
          nocopy: true

volumes:
  media:
    driver: local
    driver_opts:
      type: nfs
      o: addr=nas.local,ro
      device: ":/volume1/media"
```

---

## Testing Strategy

### Unit Tests
- [x] Health checks
- [x] Configuration loading
- [x] Cache operations
- [x] Session management

### Integration Tests
- [x] Database connection pooling
- [x] Cache invalidation
- [x] Job queue distribution
- [ ] Multi-pod session sharing
- [ ] Multi-pod cache coherence

### Chaos Tests
- [ ] Pod termination during requests
- [ ] Database connection loss
- [ ] Cache unavailability
- [ ] Network partitions

---

## Monitoring Checklist

### Infrastructure
- [x] CPU usage per pod
- [x] Memory usage per pod
- [x] Network I/O
- [ ] Disk I/O (if using local storage)
- [x] Pod restart count

### Application
- [x] HTTP request rate
- [x] HTTP error rate
- [x] Response times (p50, p95, p99)
- [x] Cache hit/miss ratio
- [x] Database connection pool usage
- [x] Active sessions count
- [x] Job queue length
- [x] Job processing time

### Business
- [ ] User registrations
- [ ] Login success/failure rate
- [ ] Media library size
- [ ] Playback sessions
- [ ] API usage by endpoint

---

## Conclusion

The Revenge codebase has an **excellent foundation for cluster deployment** with proper state externalization, 12-factor compliance, and production-grade observability. The remaining issues are **well-defined and solvable**:

1. **Media storage** - Add NFS volume support (4-8h)
2. **Avatar storage** - Implement S3 backend (8-16h)
3. **Leader election** - Add River periodic jobs (8-16h)

**Total effort**: 20-40 hours (3-5 days) to achieve **100% cluster readiness**.

**Recommendation**: Start with media storage (blocking for core functionality), then avatar storage, then leader election (nice to have but not critical).

The architecture is **well-designed for horizontal scaling** and already handles the hard parts (distributed sessions, cache invalidation, job distribution). Excellent work!
