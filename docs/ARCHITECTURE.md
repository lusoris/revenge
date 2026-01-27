# Jellyfin Go - Architecture Design

## Table of Contents

- [Deployment Modes](#deployment-modes)
- [System Overview](#system-overview)
- [Architecture Principles](#architecture-principles)
- [Component Architecture](#component-architecture)
- [Data Flow](#data-flow)
- [Technology Stack](#technology-stack)
- [Deployment Architecture](#deployment-architecture)
- [Scalability Design](#scalability-design)

---

## Deployment Modes

### Single-Server Mode (Default)

**Target Users:** Home users, small teams, self-hosters
**Requirements:** One server/PC with Docker or native binary

Jellyfin Go works out-of-the-box on a single server with minimal dependencies:

```yaml
Minimum Setup:
- SQLite (embedded, no separate database)
- Local file storage
- In-memory cache (optional Redis)
- No clustering required

Recommended Setup:
- PostgreSQL (local instance)
- Redis or Dragonfly (local, optional)
- Local file storage or NFS mount
```

**Features in Single-Server Mode:**
- ✅ All core functionality
- ✅ Hardware transcoding
- ✅ Multiple users
- ✅ Library management
- ✅ Search (PostgreSQL full-text or Typesense optional)
- ✅ API compatibility

### Enterprise/Multi-Instance Mode (Optional)

**Target Users:** Large deployments, high availability requirements
**Requirements:** Kubernetes cluster, load balancer, shared storage

**Additional Features:**
- Horizontal scaling (10+ instances)
- Distributed coordination
- CDN integration
- Advanced monitoring
- High availability (99.9%+)

---

## System Overview

Jellyfin Go is a flexible media server built with Go, designed to run anywhere from a Raspberry Pi to a Kubernetes cluster.

### High-Level Architecture

```
┌──────────────────────────────────────────────────────────────┐
│                     Load Balancer / CDN                       │
│              (NGINX / HAProxy / Cloudflare)                   │
└────────┬─────────────┬─────────────┬──────────────────────────┘
         │             │             │
    ┌────▼────┐   ┌────▼────┐   ┌────▼────┐
    │Jellyfin │   │Jellyfin │   │Jellyfin │  API Instances
    │  Go #1  │   │  Go #2  │   │  Go #3  │  (Stateless)
    │         │   │         │   │         │
    │Ristretto│   │Ristretto│   │Ristretto│  L1 Cache (Local)
    └────┬────┘   └────┬────┘   └────┬────┘
         │             │             │
         └─────────────┴─────────────┘
                       │
         ┌─────────────┴─────────────┐
         │                           │
    ┌────▼─────────────────────┐ ┌──▼──────────────┐
    │  Dragonfly Cluster       │ │   Typesense     │
    │  (3 nodes)               │ │   Cluster       │
    │  - Sessions              │ │   (3 nodes)     │
    │  - Job queues            │ │   - Search      │
    │  - Distributed locks     │ │   - Facets      │
    │  - Pub/sub events        │ │   - Aggregation │
    └────┬─────────────────────┘ └─────────────────┘
         │                           │
         │                           │
    ┌────▼────────────────┬──────────▼─────────────┐
    │  PostgreSQL Primary │  PostgreSQL Replica 1  │
    │  (Write + Read)     │  (Read-only)          │
    └─────────────────────┴────────────────────────┘
              │                    │
              │  Streaming         │
              │  Replication       │
              └────────────────────┘
```

---

## Architecture Principles

### 1. Clean Architecture (Hexagonal)

```
┌─────────────────────────────────────────────────────────┐
│                    API Layer                             │
│  (HTTP Handlers, WebSocket, gRPC)                       │
│  - gorilla/mux routing                                   │
│  - OpenAPI/Swagger docs                                  │
│  - Rate limiting middleware                              │
└───────────────────┬─────────────────────────────────────┘
                    │
┌───────────────────▼─────────────────────────────────────┐
│                Service Layer                             │
│  (Business Logic)                                        │
│  - Media management                                      │
│  - Transcoding orchestration                            │
│  - User management                                       │
│  - Session handling                                      │
└───────────────────┬─────────────────────────────────────┘
                    │
┌───────────────────▼─────────────────────────────────────┐
│                Domain Layer                              │
│  (Core Business Entities)                                │
│  - User, Media, Library, Session                        │
│  - Domain events                                         │
│  - Business rules                                        │
└───────────────────┬─────────────────────────────────────┘
                    │
┌───────────────────▼─────────────────────────────────────┐
│              Infrastructure Layer                        │
│  (External Dependencies)                                 │
│  - PostgreSQL (sqlc)                                    │
│  - Dragonfly (go-redis)                                 │
│  - Typesense (typesense-go)                             │
│  - FFmpeg (os/exec)                                     │
│  - File system                                           │
└─────────────────────────────────────────────────────────┘
```

### 2. Dependency Injection (uber-go/fx)

All components registered via DI for testability and modularity:

```go
fx.New(
    // Infrastructure
    fx.Provide(NewPostgreSQLConnection),
    fx.Provide(NewDragonflyClient),
    fx.Provide(NewTypesenseClient),
    fx.Provide(NewFFmpegExecutor),
    
    // Repositories
    fx.Provide(NewUserRepository),
    fx.Provide(NewMediaRepository),
    fx.Provide(NewSessionRepository),
    
    // Services
    fx.Provide(NewAuthService),
    fx.Provide(NewMediaService),
    fx.Provide(NewTranscodingService),
    
    // HTTP Handlers
    fx.Provide(NewUserHandler),
    fx.Provide(NewMediaHandler),
    fx.Provide(NewStreamHandler),
    
    // Lifecycle
    fx.Invoke(StartHTTPServer),
)
```

### 3. Interface-Based Design

All external dependencies behind interfaces for testing and swapping:

```go
type MediaRepository interface {
    GetByID(ctx context.Context, id string) (*Media, error)
    Search(ctx context.Context, query SearchQuery) ([]*Media, error)
    Create(ctx context.Context, media *Media) error
    Update(ctx context.Context, media *Media) error
}

type CacheProvider interface {
    Get(ctx context.Context, key string) (interface{}, error)
    Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
}

type SearchEngine interface {
    Index(ctx context.Context, doc Document) error
    Search(ctx context.Context, query string, opts SearchOptions) (*SearchResult, error)
}
```

---

## Component Architecture

### 1. API Layer

**Responsibilities:**
- HTTP request handling
- Input validation
- Authentication/authorization
- Rate limiting
- Response formatting

**Components:**
- `handlers/` - HTTP handlers for each API endpoint
- `middleware/` - Authentication, logging, metrics, rate limiting
- `models/` - Request/response DTOs

**Tech Stack:**
- `gorilla/mux` for routing
- `go-playground/validator` for validation
- `swaggo/swag` for OpenAPI docs

### 2. Service Layer

**Responsibilities:**
- Business logic orchestration
- Transaction management
- Event publishing
- Cross-cutting concerns

**Components:**
- `services/auth/` - Authentication and authorization
- `services/media/` - Media library management
- `services/transcoding/` - FFmpeg job orchestration
- `services/search/` - Search indexing coordination

### 3. Domain Layer

**Responsibilities:**
- Core business entities
- Domain rules and validation
- Domain events

**Components:**
- `domain/user/` - User entity and rules
- `domain/media/` - Media entity and rules
- `domain/session/` - Session entity and rules
- `domain/events/` - Domain event definitions

### 4. Infrastructure Layer

**Responsibilities:**
- External system integration
- Data persistence
- Caching
- Message queuing

**Components:**
- `infra/postgres/` - PostgreSQL repository implementations
- `infra/dragonfly/` - Dragonfly cache client
- `infra/typesense/` - Typesense search client
- `infra/ffmpeg/` - FFmpeg process management
- `infra/filesystem/` - File system operations

---

## Data Flow

### 1. Media Upload Flow

```
Client → API Handler → Auth Middleware → Upload Service
                                             │
                                             ├─→ Virus Scan (optional)
                                             ├─→ Save to Storage
                                             ├─→ Extract Metadata (FFmpeg)
                                             ├─→ Generate Thumbnails
                                             ├─→ PostgreSQL (metadata)
                                             ├─→ Typesense (index)
                                             └─→ CDN Upload (thumbnails)
```

### 2. Streaming Flow (HLS)

```
Client Request → LB → API Instance
                        │
                        ├─→ Auth Check (JWT validation)
                        ├─→ Policy Check (can user stream?)
                        ├─→ Get Media Info (PostgreSQL/Cache)
                        │
                        ├─→ Direct Play? → Serve file
                        │
                        └─→ Transcode Needed?
                              │
                              ├─→ Check existing job (Dragonfly)
                              ├─→ Start FFmpeg (Worker Pool)
                              ├─→ Generate HLS segments
                              ├─→ Cache segments (Dragonfly)
                              └─→ Serve playlist + segments
```

### 3. Search Flow

```
Client Search → API → Search Service
                        │
                        ├─→ Check Cache (Ristretto)
                        │   └─→ Cache Hit → Return
                        │
                        ├─→ Query Typesense
                        │   ├─→ Faceted filters
                        │   ├─→ Typo tolerance
                        │   └─→ Ranking
                        │
                        ├─→ Enrich with PostgreSQL data
                        ├─→ Cache results (Ristretto)
                        └─→ Return to client
```

### 4. Multi-Instance Coordination

```
Instance 1: User updates media
    │
    ├─→ Write to PostgreSQL (primary)
    ├─→ Update Typesense index
    ├─→ Invalidate local cache (Ristretto)
    └─→ Publish event (Dragonfly pub/sub)
            │
            ├─→ Instance 2: Receive event → Invalidate cache
            ├─→ Instance 3: Receive event → Invalidate cache
            └─→ Instance N: Receive event → Invalidate cache
```

---

## Technology Stack

### Data Storage

#### PostgreSQL 16+ OR SQLite (Database)

**PostgreSQL (Recommended for production):**
- Use Cases: User accounts, media metadata, library structure, sessions, activity logs
- Configuration: Single instance (default) or with read replicas (enterprise)
- Connection pooling: pgxpool (25-64 connections per instance)
- Performance: Better for >10k media items

**SQLite (Default for single-server):**
- Use Cases: Same as PostgreSQL, embedded database
- Configuration: Single file, no separate service needed
- Performance: Excellent for <10k media items, <10 concurrent users
- Advantages: Zero-configuration, portable, perfect for home users

**Schema Design:**
- Normalized for consistency
- JSONB for flexible metadata
- Partitioned tables for logs (monthly)
- Indexes on frequent queries

#### Dragonfly/Redis (Optional Cache)

**Use Cases:**
- Session cache (hot data, sub-ms access)
- Transcoding job queue (multi-instance only)
- Distributed locks (multi-instance only)
- Pub/sub for cache invalidation (multi-instance only)
- Rate limiting counters

**Configuration:**
- Single-Server: Local Redis/Dragonfly instance OR in-memory cache
- Multi-Instance: 3-node cluster for HA
- Redis protocol compatible
- Memory: 1-2GB (single-server), 8-16GB per node (cluster)
- **Not required:** Falls back to in-memory cache if not available

**Namespaces:**
- `sessions:` - User sessions
- `jobs:` - Transcoding jobs
- `locks:` - Distributed locks
- `ratelimit:` - Rate limit counters
- `cache:` - General cache

#### Ristretto (Local Cache)

**Use Cases:**
- Frequently accessed metadata
- Computed values (policies, thumbnails)
- Database query results

**Configuration:**
- Cost-based eviction (1GB max per instance)
- TinyLFU admission policy
- TTL: 5-15 minutes

#### Typesense (Optional Search Engine)

**Use Cases:**
- Full-text media search
- Faceted filtering (genre, year, rating)
- Typo-tolerant queries
- Real-time indexing

**Configuration:**
- Single-Server: Local Typesense instance OR PostgreSQL full-text search
- Multi-Instance: 3-node Raft cluster
- Collections: media_active, media_archive, users, collections
- Memory: 2-4GB (single-server, <50k items), 10-15GB per node (cluster, 100k+ items)
- **Fallback:** PostgreSQL `tsvector` full-text search if Typesense not available

### Media Processing

#### FFmpeg (jellyfin-ffmpeg)

**Capabilities:**
- Video transcoding (H.264, HEVC, AV1, VP9)
- Audio transcoding (AAC, MP3, OPUS, FLAC)
- HLS/DASH segmentation
- Hardware acceleration (VAAPI, NVENC, QuickSync, AMF, VideoToolbox)
- Thumbnail generation

**Integration:**
- Process pool (max 10 concurrent by default)
- Command builder pattern
- Progress tracking via stderr parsing
- Graceful shutdown (SIGTERM → wait → SIGKILL)

### Observability

#### Prometheus + Grafana

**Metrics:**
- HTTP request duration (histogram)
- Active connections (gauge)
- Cache hit ratio (counter)
- Transcoding jobs active (gauge)
- Database pool size (gauge)

#### OpenTelemetry

**Tracing:**
- Distributed traces across services
- Context propagation
- Span attributes (user_id, media_id, etc.)

#### slog (Structured Logging)

**Log Levels:**
- DEBUG: Development only
- INFO: Normal operations
- WARN: Degraded but functional
- ERROR: Failures requiring attention

**Format:** JSON with context (request_id, user_id, trace_id)

#### Pyroscope (Continuous Profiling)

**Profiles:**
- CPU usage
- Memory allocation
- Goroutines
- Mutex contention

---

## Deployment Architecture

### Single-Server (Home User)

**Option 1: SQLite (Zero Dependencies)**
```bash
# Download binary
wget https://github.com/your-org/jellyfin-go/releases/latest/jellyfin-go
chmod +x jellyfin-go

# Run with defaults (SQLite, in-memory cache)
./jellyfin-go

# Opens web UI at http://localhost:8096
# Data stored in ~/.jellyfin-go/
```

**Option 2: Docker (Recommended)**
```bash
# Single container, all-in-one
docker run -d \
  -p 8096:8096 \
  -v /path/to/media:/media \
  -v jellyfin-data:/data \
  --name jellyfin-go \
  jellyfin/jellyfin-go:latest
```

**Option 3: Docker Compose with PostgreSQL**
```yaml
# docker-compose.yml (recommended for better performance)
services:
  jellyfin-go:
    image: jellyfin/jellyfin-go:latest
    ports: ["8096:8096"]
    environment:
      - DATABASE_URL=postgres://jellyfin:password@postgres:5432/jellyfin
    volumes:
      - /path/to/media:/media
      - jellyfin-config:/config
    depends_on: [postgres]
  
  postgres:
    image: postgres:16-alpine
    environment:
      - POSTGRES_DB=jellyfin
      - POSTGRES_USER=jellyfin
      - POSTGRES_PASSWORD=password
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  jellyfin-config:
  postgres_data:
```

### Development

```yaml
# docker-compose.dev.yml (full stack for development)
services:
  jellyfin-go:
    build: .
    ports: ["8096:8096"]
    environment:
      - DATABASE_URL=postgres://postgres:password@postgres:5432/jellyfin
      - REDIS_URL=redis:6379
      - TYPESENSE_URL=http://typesense:8108
    depends_on: [postgres, redis, typesense]
  
  postgres:
    image: postgres:16-alpine
    volumes: [postgres_data:/var/lib/postgresql/data]
  
  redis:
    image: redis:7-alpine
  
  typesense:
    image: typesense/typesense:26.0
```

### Production (Kubernetes)

```yaml
# Deployment structure
├── StatefulSet: jellyfin-api (10-50 replicas, HPA)
├── StatefulSet: postgres-primary (1 replica)
├── StatefulSet: postgres-replica (2 replicas)
├── StatefulSet: dragonfly (3 replicas)
├── StatefulSet: typesense (3 replicas)
├── Service: jellyfin-api (LoadBalancer)
├── Service: postgres-primary (ClusterIP)
├── Service: postgres-replica (ClusterIP)
├── Service: dragonfly (ClusterIP)
├── Service: typesense (ClusterIP)
├── Ingress: NGINX/HAProxy with TLS
├── PVC: media-storage (ReadWriteMany, NFS/EFS)
└── PVC: per StatefulSet for data persistence
```

**Scaling Strategy:**
- HPA based on CPU (>70%) and memory (>80%)
- Min replicas: 3
- Max replicas: 50
- Scale-up: +5 pods every 30s
- Scale-down: -1 pod every 5min (gradual)

**Resource Requests/Limits:**
```yaml
api:
  requests: {cpu: 2, memory: 4Gi}
  limits: {cpu: 4, memory: 8Gi}
postgres:
  requests: {cpu: 4, memory: 16Gi}
  limits: {cpu: 8, memory: 32Gi}
dragonfly:
  requests: {cpu: 2, memory: 8Gi}
  limits: {cpu: 4, memory: 16Gi}
typesense:
  requests: {cpu: 2, memory: 12Gi}
  limits: {cpu: 4, memory: 24Gi}
```

---

## Scalability Design

### Horizontal Scaling

**Stateless API Instances:**
- All state in PostgreSQL, Dragonfly, or Typesense
- No local file dependencies (media on NFS/S3)
- Can scale to 100+ instances

**Coordination Mechanisms:**
- Redlock for distributed locking (transcoding job assignment)
- Pub/sub for cache invalidation across instances
- Sticky sessions for WebSocket (cookie-based)

### Vertical Scaling

**Database:**
- PostgreSQL: Scale up to 64 cores, 256GB RAM
- Read replicas: Add more for read-heavy workloads

**Cache:**
- Dragonfly: Multi-threaded, scales to 32+ cores
- Add more nodes to cluster for memory

**Search:**
- Typesense: Scale to 32 cores per node
- Horizontal: Add more nodes to cluster

### Data Partitioning

**Time-Based:**
- Activity logs: Partitioned by month
- Old data archived to separate table

**Functional:**
- Separate databases per tenant (optional, enterprise)
- Separate Typesense collections (active vs archive)

### CDN Offloading

**Static Assets:**
- Thumbnails, posters, banners → CDN
- 95%+ traffic served from edge
- Origin only for dynamic content

**Video Segments:**
- HLS segments cached at edge
- Reduces origin load by 80%+

---

## Security Architecture

### Defense in Depth

**Layer 1: Network**
- TLS 1.3 for all external traffic
- Private VPC for internal services
- Security groups (least privilege)

**Layer 2: Application**
- Input validation (all endpoints)
- SQL injection prevention (parameterized queries)
- XSS prevention (CSP headers)
- CSRF protection (tokens)

**Layer 3: Authentication**
- JWT with short expiry (15min)
- Refresh tokens (7 days, rotated)
- API keys with rotation
- MFA support (TOTP)

**Layer 4: Authorization**
- Policy-based (20+ policies)
- Resource ownership checks (BOLA prevention)
- Rate limiting (per-user and per-IP)

**Layer 5: Data**
- Encryption at rest (PostgreSQL TDE optional)
- Encryption in transit (TLS)
- Secrets in Vault/K8s secrets
- PII data minimization

### Threat Model

**Threats:**
- DDoS attacks → CDN protection + rate limiting
- Credential stuffing → Rate limiting + MFA
- SQL injection → Parameterized queries (sqlc)
- XSS → CSP headers + output encoding
- SSRF → URL validation + private IP blocking
- Path traversal → Input sanitization + chroot

**Monitoring:**
- Failed login attempts (alert on >10/min)
- Unusual traffic patterns (alert on spikes)
- Unauthorized access attempts (403/401 rate)
- Slow query attacks (timeout enforcement)

---

## Disaster Recovery

### Backup Strategy

**PostgreSQL:**
- Full backup: Weekly (pgBackRest)
- Differential: Daily
- WAL archiving: Continuous (5min)
- Retention: 30 days
- Storage: S3-compatible

**Configuration:**
- Version controlled (Git)
- Encrypted secrets (Vault)
- Automated provisioning (Terraform/Helm)

**Media Files:**
- Primary: NFS/EFS/S3
- Backup: S3 Glacier (30-day delay)
- Snapshots: Daily

### Recovery Procedures

**RTO (Recovery Time Objective):** 4 hours
**RPO (Recovery Point Objective):** 5 minutes

**Scenarios:**

1. **Single Instance Failure:** Automatic (K8s restarts pod)
2. **Database Failure:** Promote replica (5-10 min)
3. **Complete Cluster Failure:** Restore from backup (2-4 hours)
4. **Data Corruption:** Point-in-time recovery (1-2 hours)

---

## Performance Optimization

### Database Optimization

- Prepared statements (sqlc generates)
- Connection pooling (pgxpool)
- Indexes on frequent queries
- Query timeout (10s default)
- EXPLAIN ANALYZE for slow queries

### Caching Strategy

- Cache-aside pattern
- TTL-based invalidation
- Event-based invalidation (pub/sub)
- Negative caching (404s)
- Probabilistic early expiration

### API Optimization

- Response compression (gzip/brotli)
- Pagination (max 100 items)
- Partial responses (field selection)
- Conditional requests (ETag, If-None-Match)
- HTTP/2 server push (optional)

### Media Streaming Optimization

- Direct play when possible (no transcoding)
- HLS adaptive bitrate
- Segment pre-generation (popular content)
- CDN caching (edge serving)
- Byte-range request support

---

This architecture is designed for production workloads with 10,000+ concurrent users and 100,000+ media items. All components are horizontally scalable and highly available.
