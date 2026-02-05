# TODO A8: Cluster Readiness

**Phase**: A8
**Priority**: P0 (Blocker for Production)
**Effort**: 24-40 hours
**Status**: Pending
**Dependencies**: A7 (Security Fixes)
**Created**: 2026-02-05

---

## Overview

Make Revenge fully ready for cluster deployment on Kubernetes, k3s, and Docker Swarm:
- **Hybrid Storage**: NFS for media files, S3/MinIO for user content (avatars)
- **Leader Election**: River + Raft for periodic jobs
- **Request Correlation**: X-Request-ID tracking
- **Observability**: Enhanced monitoring

**Current Status**: 84% cluster-ready
**Target**: 100% production-ready

**Source**: [REPORT_3_CLUSTER_READINESS.md](REPORT_3_CLUSTER_READINESS.md)

---

## Decision Log

| Topic | Decision | Reason |
|-------|----------|--------|
| Media Storage | NFS/CephFS (ReadOnlyMany) | Best performance, no code changes needed |
| Avatar Storage | S3/MinIO | Cloud-native, scalable, interface ready |
| Leader Election | River + Raft | River for jobs, Raft for periodic cleanup |
| Deployment Target | K8s primary, k3s/Swarm supported | Enterprise + lightweight options |

---

## Tasks

### A8.1: Hybrid Storage Implementation 🔴 CRITICAL

**Priority**: P0
**Effort**: 16-24h

#### A8.1.1: NFS Volume Support for Media

**Location**:
- `charts/revenge/values.yaml`
- `charts/revenge/templates/deployment.yaml`
- `docker-compose.yml`

**Goal**: Mount shared NFS volume for media files across all pods.

**Helm Chart Changes**:

```yaml
# values.yaml
media:
  persistence:
    enabled: true
    storageClass: "nfs-client"  # or manual PV
    accessMode: ReadOnlyMany    # Multiple pods can read
    size: 1Ti
    nfs:
      server: ""                # User configurable
      path: ""                  # User configurable
      readOnly: true
    # Or use existingClaim
    existingClaim: ""

  # Movie library paths (inside container)
  moviePaths:
    - /media/movies

  # TV library paths (for future)
  tvPaths:
    - /media/tv

  # Music library paths (for future)
  musicPaths:
    - /media/music
```

```yaml
# templates/deployment.yaml
spec:
  template:
    spec:
      volumes:
        {{- if .Values.media.persistence.enabled }}
        - name: media
          {{- if .Values.media.persistence.existingClaim }}
          persistentVolumeClaim:
            claimName: {{ .Values.media.persistence.existingClaim }}
          {{- else if .Values.media.persistence.nfs.server }}
          nfs:
            server: {{ .Values.media.persistence.nfs.server }}
            path: {{ .Values.media.persistence.nfs.path }}
            readOnly: {{ .Values.media.persistence.nfs.readOnly }}
          {{- end }}
        {{- end }}

      containers:
        - name: revenge
          volumeMounts:
            {{- if .Values.media.persistence.enabled }}
            - name: media
              mountPath: /media
              readOnly: true
            {{- end }}
```

**PersistentVolume Template** (for manual setup):

```yaml
# templates/media-pv.yaml (optional, user can create manually)
{{- if and .Values.media.persistence.enabled (not .Values.media.persistence.existingClaim) (not .Values.media.persistence.nfs.server) }}
apiVersion: v1
kind: PersistentVolume
metadata:
  name: {{ include "revenge.fullname" . }}-media
spec:
  capacity:
    storage: {{ .Values.media.persistence.size }}
  accessModes:
    - ReadOnlyMany
  storageClassName: {{ .Values.media.persistence.storageClass }}
  # User must configure their storage backend
  # Example: NFS
  # nfs:
  #   server: nas.example.com
  #   path: /volume1/media
{{- end }}
```

**Subtasks**:
- [ ] Add media volume configuration to Helm values
- [ ] Update deployment template with volume mounts
- [ ] Create example PersistentVolume YAML
- [ ] Update documentation with NFS setup guide
- [ ] Test with multiple pods (ensure all can read media)

---

#### A8.1.2: S3 Storage for Avatars

**Location**:
- `internal/service/storage/s3.go` (new)
- `internal/service/storage/module.go`
- `internal/config/config.go`

**Goal**: Implement S3-compatible storage backend for user-generated content.

**Interface** (already exists):
```go
// internal/service/storage/storage.go
type Storage interface {
    Store(ctx context.Context, key string, reader io.Reader, contentType string) (string, error)
    Get(ctx context.Context, key string) (io.ReadCloser, error)
    Delete(ctx context.Context, key string) error
}
```

**Implementation**:

```go
// internal/service/storage/s3.go
package storage

import (
    "context"
    "fmt"
    "io"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/credentials"
    "github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
    client *s3.Client
    bucket string
}

func NewS3Storage(cfg S3Config) (*S3Storage, error) {
    // Load AWS config
    awsCfg, err := config.LoadDefaultConfig(context.Background(),
        config.WithRegion(cfg.Region),
        config.WithCredentialsProvider(
            credentials.NewStaticCredentialsProvider(
                cfg.AccessKeyID,
                cfg.SecretAccessKey,
                "",
            ),
        ),
    )
    if err != nil {
        return nil, fmt.Errorf("failed to load AWS config: %w", err)
    }

    // Create S3 client with custom endpoint (for MinIO)
    client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
        if cfg.Endpoint != "" {
            o.BaseEndpoint = aws.String(cfg.Endpoint)
            o.UsePathStyle = true  // Required for MinIO
        }
    })

    return &S3Storage{
        client: client,
        bucket: cfg.Bucket,
    }, nil
}

func (s *S3Storage) Store(ctx context.Context, key string, reader io.Reader, contentType string) (string, error) {
    _, err := s.client.PutObject(ctx, &s3.PutObjectInput{
        Bucket:      aws.String(s.bucket),
        Key:         aws.String(key),
        Body:        reader,
        ContentType: aws.String(contentType),
    })
    if err != nil {
        return "", fmt.Errorf("failed to upload to S3: %w", err)
    }

    // Return URL (public or signed depending on bucket config)
    url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucket, key)
    return url, nil
}

func (s *S3Storage) Get(ctx context.Context, key string) (io.ReadCloser, error) {
    output, err := s.client.GetObject(ctx, &s3.GetObjectInput{
        Bucket: aws.String(s.bucket),
        Key:    aws.String(key),
    })
    if err != nil {
        return nil, fmt.Errorf("failed to get from S3: %w", err)
    }

    return output.Body, nil
}

func (s *S3Storage) Delete(ctx context.Context, key string) error {
    _, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
        Bucket: aws.String(s.bucket),
        Key:    aws.String(key),
    })
    if err != nil {
        return fmt.Errorf("failed to delete from S3: %w", err)
    }

    return nil
}
```

**Configuration**:

```go
// internal/config/config.go
type Config struct {
    // ... existing fields ...

    Storage StorageConfig `koanf:"storage"`
}

type StorageConfig struct {
    Backend string    `koanf:"backend" validate:"oneof=local s3"`  // "local" or "s3"
    Local   LocalStorageConfig `koanf:"local"`
    S3      S3Config           `koanf:"s3"`
}

type LocalStorageConfig struct {
    Path string `koanf:"path" validate:"required_if=Backend local"`
}

type S3Config struct {
    Endpoint        string `koanf:"endpoint"`         // For MinIO: "http://minio:9000"
    Region          string `koanf:"region"`           // "us-east-1"
    Bucket          string `koanf:"bucket" validate:"required_if=Backend s3"`
    AccessKeyID     string `koanf:"access_key_id" validate:"required_if=Backend s3"`
    SecretAccessKey string `koanf:"secret_access_key" validate:"required_if=Backend s3"`
    UsePathStyle    bool   `koanf:"use_path_style"`   // true for MinIO
}
```

**Dependency Injection**:

```go
// internal/service/storage/module.go
func provideStorage(cfg *config.Config) (Storage, error) {
    switch cfg.Storage.Backend {
    case "s3":
        return NewS3Storage(cfg.Storage.S3)
    case "local":
        return NewLocalStorage(cfg.Storage.Local)
    default:
        return nil, fmt.Errorf("unknown storage backend: %s", cfg.Storage.Backend)
    }
}
```

**Docker Compose with MinIO**:

```yaml
# docker-compose.yml
services:
  minio:
    image: minio/minio:latest
    command: server /data --console-address ":9001"
    ports:
      - "9000:9000"
      - "9001:9001"
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin
    volumes:
      - minio_data:/data
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:9000/minio/health/live"]
      interval: 30s
      timeout: 20s
      retries: 3

  revenge:
    # ... existing config ...
    environment:
      REVENGE_STORAGE_BACKEND: s3
      REVENGE_STORAGE_S3_ENDPOINT: http://minio:9000
      REVENGE_STORAGE_S3_REGION: us-east-1
      REVENGE_STORAGE_S3_BUCKET: revenge-avatars
      REVENGE_STORAGE_S3_ACCESS_KEY_ID: minioadmin
      REVENGE_STORAGE_S3_SECRET_ACCESS_KEY: minioadmin
      REVENGE_STORAGE_S3_USE_PATH_STYLE: "true"

volumes:
  minio_data:
```

**Subtasks**:
- [ ] Add AWS SDK dependency (`go get github.com/aws/aws-sdk-go-v2/...`)
- [ ] Implement S3Storage
- [ ] Add storage config to config.go
- [ ] Update storage module with backend selection
- [ ] Add MinIO to docker-compose.yml
- [ ] Write integration tests with testcontainers MinIO
- [ ] Document S3/MinIO setup
- [ ] Create bucket initialization script

---

### A8.2: Leader Election with River + Raft 🟠 HIGH

**Priority**: P1
**Effort**: 12-16h

**Goal**: Prevent duplicate execution of periodic cleanup jobs in multi-pod setup.

**Decision**:
- **River** for all job distribution (already working)
- **Raft** for leader election of periodic cleanup jobs

#### A8.2.1: Raft Leader Election

**Location**:
- `internal/infra/raft/` (new package)
- `internal/infra/raft/module.go`
- `internal/infra/raft/election.go`

**Dependencies**: Already in go.mod!
```
github.com/hashicorp/raft v1.8.0
```

**Implementation**:

```go
// internal/infra/raft/election.go
package raft

import (
    "fmt"
    "net"
    "os"
    "path/filepath"
    "time"

    "github.com/hashicorp/raft"
    raftboltdb "github.com/hashicorp/raft-boltdb"
)

type LeaderElection struct {
    raft *raft.Raft
}

func NewLeaderElection(cfg Config) (*LeaderElection, error) {
    // Setup Raft configuration
    config := raft.DefaultConfig()
    config.LocalID = raft.ServerID(cfg.NodeID)

    // Setup Raft communication
    addr, err := net.ResolveTCPAddr("tcp", cfg.BindAddr)
    if err != nil {
        return nil, err
    }

    transport, err := raft.NewTCPTransport(cfg.BindAddr, addr, 3, 10*time.Second, os.Stderr)
    if err != nil {
        return nil, err
    }

    // Create the snapshot store
    snapshots, err := raft.NewFileSnapshotStore(cfg.DataDir, 2, os.Stderr)
    if err != nil {
        return nil, err
    }

    // Create the log store and stable store
    logStore, err := raftboltdb.NewBoltStore(filepath.Join(cfg.DataDir, "raft-log.db"))
    if err != nil {
        return nil, err
    }

    stableStore, err := raftboltdb.NewBoltStore(filepath.Join(cfg.DataDir, "raft-stable.db"))
    if err != nil {
        return nil, err
    }

    // Create FSM (simple no-op for leader election only)
    fsm := &simpleFSM{}

    // Instantiate the Raft system
    ra, err := raft.NewRaft(config, fsm, logStore, stableStore, snapshots, transport)
    if err != nil {
        return nil, err
    }

    // Bootstrap cluster if first node
    if cfg.Bootstrap {
        configuration := raft.Configuration{
            Servers: []raft.Server{
                {
                    ID:      config.LocalID,
                    Address: transport.LocalAddr(),
                },
            },
        }
        ra.BootstrapCluster(configuration)
    }

    return &LeaderElection{raft: ra}, nil
}

func (le *LeaderElection) IsLeader() bool {
    return le.raft.State() == raft.Leader
}

func (le *LeaderElection) LeaderAddr() string {
    _, id := le.raft.LeaderWithID()
    return string(id)
}

func (le *LeaderElection) Close() error {
    return le.raft.Shutdown().Error()
}

// Simple FSM for leader election (no state needed)
type simpleFSM struct{}

func (f *simpleFSM) Apply(*raft.Log) interface{} { return nil }
func (f *simpleFSM) Snapshot() (raft.FSMSnapshot, error) { return &simpleSnapshot{}, nil }
func (f *simpleFSM) Restore(io.ReadCloser) error { return nil }

type simpleSnapshot struct{}
func (s *simpleSnapshot) Persist(sink raft.SnapshotSink) error { return nil }
func (s *simpleSnapshot) Release() {}
```

**Configuration**:

```go
// internal/config/config.go
type Config struct {
    // ... existing ...
    Raft RaftConfig `koanf:"raft"`
}

type RaftConfig struct {
    Enabled   bool   `koanf:"enabled"`      // Enable Raft leader election
    NodeID    string `koanf:"node_id"`      // Unique node ID (hostname or UUID)
    BindAddr  string `koanf:"bind_addr"`    // "0.0.0.0:7000"
    DataDir   string `koanf:"data_dir"`     // "/data/raft"
    Bootstrap bool   `koanf:"bootstrap"`    // Bootstrap cluster (first node only)
}
```

**Subtasks**:
- [ ] Implement Raft leader election
- [ ] Add configuration
- [ ] Integrate with fx lifecycle
- [ ] Add health check endpoint
- [ ] Document Raft setup for K8s StatefulSet

---

#### A8.2.2: Integrate with Cleanup Jobs

**Location**:
- `internal/infra/jobs/cleanup_job.go`
- `internal/service/activity/cleanup.go`

**Goal**: Only run periodic cleanup jobs on leader node.

```go
// internal/infra/jobs/cleanup_job.go
type CleanupWorker struct {
    activityService *activity.Service
    leaderElection  *raft.LeaderElection
    logger          *slog.Logger
}

func (w *CleanupWorker) Work(ctx context.Context, job *river.Job[CleanupArgs]) error {
    // Only leader executes cleanup
    if !w.leaderElection.IsLeader() {
        w.logger.Info("skipping cleanup (not leader)",
            "leader", w.leaderElection.LeaderAddr(),
        )
        return nil
    }

    w.logger.Info("starting activity log cleanup (leader)")

    // ... existing cleanup logic ...
}
```

**Periodic Job Scheduling** (River):

```go
// internal/infra/jobs/periodic.go (new file)
package jobs

import (
    "time"

    "github.com/riverqueue/river"
)

func SetupPeriodicJobs(client *river.Client[pgx.Tx]) error {
    periodicJobs := []*river.PeriodicJob{
        river.NewPeriodicJob(
            river.PeriodicInterval(1 * time.Hour),
            func() (river.JobArgs, *river.InsertOpts) {
                return &CleanupArgs{}, &river.InsertOpts{
                    Queue: QueueDefault,
                    // Unique by hour to prevent duplicates
                    UniqueOpts: river.UniqueOpts{
                        ByPeriod: 1 * time.Hour,
                    },
                }
            },
            &river.PeriodicJobOpts{RunOnStart: true},
        ),
    }

    return client.PeriodicJobs().AddMany(periodicJobs)
}
```

**Subtasks**:
- [ ] Add LeaderElection to cleanup workers
- [ ] Implement leader check before execution
- [ ] Setup River periodic jobs
- [ ] Test with 3-pod cluster (only leader runs cleanup)
- [ ] Add metrics for leader changes

---

### A8.3: Request Correlation IDs 🟡 MEDIUM

**Priority**: P2
**Effort**: 3-4h

**Location**:
- `internal/api/middleware/request_id.go` (new)
- `internal/api/server.go`

**Goal**: Track requests across services with X-Request-ID header.

**Implementation**:

```go
// internal/api/middleware/request_id.go
package middleware

import (
    "context"
    "net/http"

    "github.com/google/uuid"
)

type contextKey string

const RequestIDKey contextKey = "request_id"

func RequestID() func(next http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Extract or generate request ID
            requestID := r.Header.Get("X-Request-ID")
            if requestID == "" {
                requestID = uuid.New().String()
            }

            // Inject into context
            ctx := context.WithValue(r.Context(), RequestIDKey, requestID)

            // Return in response header
            w.Header().Set("X-Request-ID", requestID)

            // Continue
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func GetRequestID(ctx context.Context) string {
    if id, ok := ctx.Value(RequestIDKey).(string); ok {
        return id
    }
    return ""
}
```

**Logging Integration**:

```go
// internal/infra/logging/middleware.go
func LoggingMiddleware(logger *slog.Logger) func(next http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            requestID := middleware.GetRequestID(r.Context())

            // Add request ID to all logs
            logger := logger.With("request_id", requestID)
            ctx := context.WithValue(r.Context(), loggerKey, logger)

            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

**Subtasks**:
- [ ] Implement RequestID middleware
- [ ] Add to middleware chain (before logging)
- [ ] Update logging to include request_id
- [ ] Propagate request_id to all service calls
- [ ] Add to metrics/tracing

---

## Helm Chart Updates

### values.yaml Additions

```yaml
# Storage
storage:
  backend: local  # or s3
  local:
    path: /data/storage
  s3:
    endpoint: ""  # For MinIO: http://minio:9000
    region: us-east-1
    bucket: revenge-avatars
    accessKeyId: ""
    secretAccessKey: ""
    usePathStyle: false

# Media (NFS)
media:
  persistence:
    enabled: false
    storageClass: ""
    accessMode: ReadOnlyMany
    size: 1Ti
    existingClaim: ""
    nfs:
      server: ""
      path: ""
      readOnly: true

# Raft
raft:
  enabled: false  # Enable for multi-pod deployments
  statefulset: false  # Use StatefulSet instead of Deployment
  dataDir: /data/raft
```

**Subtasks**:
- [ ] Update values.yaml
- [ ] Update README with cluster setup guide
- [ ] Create example values for K8s/k3s/Swarm
- [ ] Document NFS setup steps
- [ ] Document MinIO setup steps

---

## Testing

### Integration Tests

**Location**: `internal/infra/storage/s3_integration_test.go`

```go
func TestS3Storage_Integration(t *testing.T) {
    // Use testcontainers MinIO
    ctx := context.Background()

    minioC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: testcontainers.ContainerRequest{
            Image: "minio/minio:latest",
            Cmd:   []string{"server", "/data"},
            ExposedPorts: []string{"9000/tcp"},
            Env: map[string]string{
                "MINIO_ROOT_USER":     "minioadmin",
                "MINIO_ROOT_PASSWORD": "minioadmin",
            },
            WaitingFor: wait.ForHTTP("/minio/health/live").WithPort("9000"),
        },
        Started: true,
    })
    require.NoError(t, err)
    defer minioC.Terminate(ctx)

    // Test S3Storage operations
    // ...
}
```

**Subtasks**:
- [ ] Write S3 integration tests
- [ ] Write Raft leader election tests
- [ ] Test multi-pod scenarios (docker-compose with 3 instances)
- [ ] Load test with concurrent requests

---

## Documentation

**Files to Create**:
- `docs/deployment/kubernetes.md`
- `docs/deployment/nfs-setup.md`
- `docs/deployment/minio-setup.md`
- `docs/deployment/raft-clustering.md`

**Subtasks**:
- [ ] Document K8s deployment with NFS
- [ ] Document MinIO/S3 setup
- [ ] Document Raft clustering
- [ ] Add troubleshooting guide
- [ ] Create example manifests

---

## Verification Checklist

- [ ] NFS volume mounts working in K8s
- [ ] S3/MinIO storage functional
- [ ] Raft leader election working
- [ ] Only leader runs periodic jobs
- [ ] Request IDs in all logs
- [ ] Multi-pod deployment tested (3+ pods)
- [ ] Graceful shutdown working
- [ ] Health checks passing
- [ ] Metrics exposed
- [ ] Documentation complete

---

## Dependencies

**Requires**:
- A7: Security Fixes (for production-ready auth)

**Blocks**:
- Production deployment

---

**Completion Criteria**:
✅ 3-pod cluster running successfully
✅ Media accessible from all pods (NFS)
✅ Avatars stored in S3/MinIO
✅ Only leader runs cleanup jobs
✅ Request correlation working
✅ All tests passing
