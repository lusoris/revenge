# Revenge - Technology Stack

## 🚀 Modern Go Setup (2025)

### Core Language
- **Go 1.25** (Latest stable)
  - Built-in container/cgroup support (no automaxprocs needed)
  - `sync.WaitGroup.Go()` for cleaner goroutine management
  - `testing.B.Loop()` for benchmarks
  - Enhanced iterator protocol
  - Improved generics performance

### Standard Library First
- **net/http.ServeMux** (Go 1.22+) - Enhanced routing with method & path patterns
- **log/slog** (Go 1.21+) - Structured logging built-in
- **context** - First-class context support
- **errors** - Modern error handling with `errors.Is` and `errors.As`

### Dependencies (Carefully Selected)

#### Configuration
- **koanf v2** - Modern, type-safe configuration
  - Replaces deprecated Viper
  - Environment variables with `REVENGE_` prefix
  - YAML file support

#### Logging
- **tint** - Beautiful console logging for slog
  - Colored output in development
  - Human-readable format
  - Minimal overhead

#### Dependency Injection
- **uber-go/fx v1.24** - Latest stable
  - Improved performance
  - Better error messages
  - Lifecycle management

#### Database
- **pgx v5** - PostgreSQL 18+ driver
  - Better performance than lib/pq
  - Native prepared statements
  - Connection pooling
- **sqlc** - Type-safe SQL query generation
  - Compile-time query validation
  - No ORM overhead
- **golang-migrate** - Database migrations
  - Per-module migration folders

#### Cache & Search

> **📋 Package Versions**: See [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#go-dependencies-core) for exact versions.

- **Dragonfly** (via rueidis) - Redis-compatible cache
  - **rueidis** - 14x faster than go-redis, auto-pipelining
  - **otter** - Local W-TinyLFU cache, 50% less memory than ristretto
  - **sturdyc** - 90% API call reduction via request coalescing
- **Typesense** (via typesense-go/v4) - Search engine
  - Lightning-fast typo-tolerant search
  - Faceted filtering
  - Lower latency than Elasticsearch

#### Job Queue
- **River** - PostgreSQL-native job queue
  - No separate queue infrastructure
  - Transactional job enqueuing
  - Built-in retries and dead letter queue
  - Worker pool management

#### API Documentation
- **ogen** - OpenAPI spec-first code generation
  - OpenAPI 3.1 support
  - Type-safe generated handlers
  - Built-in validation

### Removed/Replaced

❌ **gorilla/mux** → ✅ **net/http (stdlib)**
- Reason: Go 1.22+ has enhanced routing built-in

❌ **Viper** → ✅ **koanf v2**
- Reason: Viper is in maintenance mode

✅ **slog + tint (dev)** + ✅ **zap (prod)**
- Reason: tint for colorized dev logs, zap for high-perf JSON in production

❌ **automaxprocs** → ✅ **Go 1.25 built-in**
- Reason: Go 1.25 has native container support

## 📦 Dependency Philosophy

1. **Stdlib First** - Use Go standard library when possible
2. **Minimal Dependencies** - Only add when truly needed
3. **Active Maintenance** - No deprecated or abandoned packages
4. **Performance** - Choose performant, well-tested libraries
5. **Type Safety** - Prefer strongly-typed APIs
6. **PostgreSQL-native** - Leverage PostgreSQL features (River, sqlc)

## 🎯 Modern Go Features Used

### Go 1.25 Features (NEW)
- **Built-in container support**: No automaxprocs needed
- **`sync.WaitGroup.Go()`**: Cleaner goroutine management
- **`testing.B.Loop()`**: Better benchmark iteration

### Go 1.22+ Features
- **Enhanced ServeMux Patterns**: `mux.HandleFunc("GET /api/users/{id}", ...)`
- **Ranging over functions**: Iterator support
- **Profile-guided optimization**: Better compiler performance

### Go 1.21+ Features
- **log/slog**: Structured logging
- **min/max built-ins**: No more custom functions
- **clear built-in**: Clear maps and slices

### Modern Patterns
- **Context-first APIs**: All functions accept context
- **Generics**: Type-safe collections where appropriate
- **Functional options**: Clean configuration
- **Error wrapping**: Proper error chains with `%w`

## 🔄 Technology Stack Summary

> **📋 Authoritative Package Versions**: See [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#go-dependencies-core) for current package versions and imports.

This section documents the **rationale** for technology choices. For exact versions, always refer to SOURCE_OF_TRUTH.

---

## 🎨 Frontend Stack

### Core Framework
- **SvelteKit 2** - Modern, fast, SSR-capable
  - Excellent performance (smaller bundles than React/Vue)
  - Built-in routing and SSR
  - TypeScript-first
  - Minimal boilerplate

### UI Framework
- **Tailwind CSS 4** - Utility-first CSS
  - Full dark/light mode support
  - Responsive by default
  - Custom theme system
- **shadcn-svelte** - High-quality components
  - Accessible (WCAG 2.1 AA)
  - Customizable
  - No runtime overhead (copy-paste components)

### State Management
- **Svelte Stores** - Built-in reactivity
- **TanStack Query** - Server state management
  - Caching, background refresh
  - Optimistic updates

### Features
| Feature | Implementation |
|---------|----------------|
| Authentication | JWT + refresh tokens, OIDC |
| Authorization | Full RBAC (admin, mod, user, guest) |
| Themes | CSS variables, user-selectable |
| Responsive | Mobile-first, tablet, desktop |
| i18n | Built-in internationalization |
| PWA | Offline support, installable |
| Real-time | WebSocket for live updates |

### Frontend Structure
```
web/
  src/
    lib/
      components/        # Shared UI components
      stores/            # Svelte stores
      api/               # API client (generated)
      utils/             # Helpers
    routes/
      (app)/             # Main app routes
        (admin)/         # Admin panel
        (media)/         # Media browsing
        (player)/        # Video player
      (auth)/            # Login/register
      api/               # API routes (if needed)
    app.css              # Global styles
    app.html             # HTML template
  static/                # Static assets
  tailwind.config.ts
  svelte.config.js
  vite.config.ts
```

### Admin Features
- User management (CRUD, roles, permissions)
- Library management (scan, refresh, delete)
- Server settings (all configurable)
- Activity logs and analytics
- Module enable/disable
- Theme customization
- Scheduled tasks management

---

## 📱 Client Profiles (Blackbeard)

Pre-configured transcode profiles for device groups:

### Device Groups
| Group | Devices | Default Profile |
|-------|---------|-----------------|
| `tv_4k` | LG WebOS, Samsung Tizen, Android TV | `hevc_4k_hdr` |
| `tv_hd` | Roku, Fire TV Stick, older TVs | `h264_1080p` |
| `mobile_ios` | iPhone, iPad | `h264_1080p_hls` |
| `mobile_android` | Android phones/tablets | `h264_1080p_hls` |
| `desktop_app` | Electron app, native clients | `hevc_4k` |
| `browser_modern` | Chrome, Firefox, Edge | `vp9_1080p_dash` |
| `browser_legacy` | Safari, older browsers | `h264_720p_hls` |
| `low_bandwidth` | Any device, slow connection | `h264_480p_hls` |

### Profile Definition
```go
type TranscodeProfile struct {
    ID              string   `json:"id"`
    Name            string   `json:"name"`

    // Video
    VideoCodec      string   `json:"video_codec"`      // h264, hevc, av1, vp9
    MaxWidth        int      `json:"max_width"`
    MaxHeight       int      `json:"max_height"`
    MaxBitrate      int      `json:"max_bitrate_kbps"`

    // Audio
    AudioCodec      string   `json:"audio_codec"`      // aac, ac3, opus
    AudioChannels   int      `json:"audio_channels"`
    AudioBitrate    int      `json:"audio_bitrate_kbps"`

    // Container
    Container       string   `json:"container"`        // mp4, webm, ts
    StreamFormat    string   `json:"stream_format"`    // hls, dash, progressive

    // Features
    AllowHDR        bool     `json:"allow_hdr"`
    HardwareDecode  bool     `json:"hardware_decode"`
    HardwareEncode  bool     `json:"hardware_encode"`
}
```

### Profile Selection Flow
```
1. Client connects with User-Agent + capability report
2. Revenge maps to device group
3. Get base profile for group
4. Adjust for:
   - Measured bandwidth (external clients)
   - User quality preference
   - Server load
5. Send profile ID to Blackbeard
```

### Bandwidth-Based Override
| Bandwidth | Profile Override |
|-----------|------------------|
| < 1.5 Mbps | `h264_360p` |
| 1.5-3 Mbps | `h264_480p` |
| 3-8 Mbps | `h264_720p` |
| 8-15 Mbps | `h264_1080p` |
| 15-25 Mbps | `hevc_1080p` |
| > 25 Mbps | Use device default |


<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

### Referenced Sources

| Source | Documentation |
|--------|---------------|
| [Go log/slog](https://pkg.go.dev/log/slog) | [Local](../../sources/go/stdlib/slog.md) |
| [River Documentation](https://riverqueue.com/docs) | [Local](../../sources/tooling/river-guide.md) |
| [Typesense Go Client](https://github.com/typesense/typesense-go) | [Local](../../sources/infrastructure/typesense-go.md) |
| [koanf](https://pkg.go.dev/github.com/knadh/koanf/v2) | [Local](../../sources/tooling/koanf.md) |
| [ogen Documentation](https://ogen.dev/docs/intro) | [Local](../../sources/tooling/ogen-guide.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Technical](INDEX.md)

### In This Section

- [API Reference](API.md)
- [Revenge - Audio Streaming & Progress Tracking](AUDIO_STREAMING.md)
- [Configuration Reference](CONFIGURATION.md)
- [Revenge - Frontend Architecture](FRONTEND.md)
- [Revenge - Advanced Offloading Architecture](OFFLOADING.md)

### Related Topics

- [Revenge - Architecture v2](../architecture/01_ARCHITECTURE.md) _Architecture_
- [Revenge - Design Principles](../architecture/02_DESIGN_PRINCIPLES.md) _Architecture_
- [Revenge - Metadata System](../architecture/03_METADATA_SYSTEM.md) _Architecture_
- [Revenge - Player Architecture](../architecture/04_PLAYER_ARCHITECTURE.md) _Architecture_
- [Plugin Architecture Decision](../architecture/05_PLUGIN_ARCHITECTURE_DECISION.md) _Architecture_

### Indexes

- [Design Index](../DESIGN_INDEX.md) - All design docs by category/topic
- [Source of Truth](../00_SOURCE_OF_TRUTH.md) - Package versions and status

<!-- DESIGN-BREADCRUMBS-END -->

---

## 🖥️ Deployment & Platform Support

### Primary Target
- **Docker** (recommended)
  - Single container for Revenge
  - Compose for full stack (PostgreSQL, Dragonfly, Typesense, Blackbeard)
  - Works on Docker Desktop (Windows, macOS, Linux)

### Supported Platforms

| Platform | Support Level | Notes |
|----------|---------------|-------|
| Linux (amd64) | ✅ Full | Primary development target |
| Linux (arm64) | ✅ Full | Raspberry Pi 4+, ARM servers |
| macOS (arm64) | ✅ Full | Apple Silicon |
| macOS (amd64) | ✅ Full | Intel Macs |
| Windows (amd64) | ✅ Full | Native or WSL2 |
| FreeBSD | 🔶 Community | Should work, not tested |

### Deployment Options

#### 1. Docker Compose (Recommended)
```yaml
services:
  revenge:
    image: ghcr.io/lusoris/revenge:latest
    ports: ["8096:8096"]
    volumes:
      - ./config:/config
      - /media:/media:ro
    depends_on: [postgres, dragonfly, typesense]

  postgres:
    image: postgres:18
    volumes: [postgres_data:/var/lib/postgresql/data]

  dragonfly:
    image: docker.dragonflydb.io/dragonflydb/dragonfly

  typesense:
    image: typesense/typesense:0.25.2

  blackbeard:
    image: ghcr.io/lusoris/blackbeard:latest
    deploy:
      resources:
        reservations:
          devices:
            - capabilities: [gpu]  # Optional: Hardware transcoding
```

#### 2. Native Binary
```bash
# Download
curl -LO https://github.com/lusoris/revenge/releases/latest/download/revenge_linux_amd64.tar.gz
tar xzf revenge_linux_amd64.tar.gz

# Configure
cp config.example.yaml config.yaml
# Edit config.yaml

# Run (requires external PostgreSQL, Dragonfly, Typesense)
./revenge serve
```

#### 3. Container Orchestration

> **📋 Full orchestration patterns**: See [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#container-orchestration) for complete Helm, K3s, and Swarm configurations.

##### Kubernetes / Helm

```bash
helm repo add revenge https://lusoris.github.io/revenge-charts
helm install revenge revenge/revenge -f values.yaml
```

##### K3s (Lightweight K8s)

```bash
# K3s comes with Traefik ingress built-in
kubectl apply -f https://lusoris.github.io/revenge/k3s-quickstart.yaml
```

##### Docker Swarm

```bash
docker stack deploy -c docker-stack.yml revenge
```

### Hardware Transcoding (Blackbeard)

| Hardware | Support | Notes |
|----------|---------|-------|
| NVIDIA GPU | ✅ NVENC | Best quality, lowest CPU |
| Intel QSV | ✅ QuickSync | Good for Intel systems |
| AMD AMF | ✅ VCE/VCN | Good for AMD systems |
| Apple VideoToolbox | ✅ | macOS only |
| VAAPI | ✅ | Linux generic |
| Software | ✅ | Fallback, CPU-intensive |

### Minimum Requirements

| Component | Minimum | Recommended |
|-----------|---------|-------------|
| CPU | 2 cores | 4+ cores |
| RAM | 2 GB | 4+ GB |
| Storage | 1 GB (app) | SSD recommended |
| PostgreSQL | 16+ | 18+ |
| Network | 100 Mbps | 1 Gbps |

### Environment Variables
```bash
# Core
REVENGE_DATABASE_URL=postgres://user:pass@localhost/revenge
REVENGE_CACHE_URL=redis://localhost:6379
REVENGE_SEARCH_URL=http://localhost:8108
REVENGE_BLACKBEARD_URL=http://localhost:9000

# Security
REVENGE_JWT_SECRET=your-secret-key
REVENGE_ADMIN_PASSWORD=initial-admin-password

# Features
REVENGE_MODULES_ENABLED=movie,tvshow,music
REVENGE_ADULT_ENABLED=false
```

## ⚡ Performance Improvements

- **Faster routing**: stdlib ServeMux is optimized
- **Less allocations**: slog is allocation-efficient
- **Better GC**: Modern Go runtime improvements
- **Profile-guided optimization**: Compiler optimizations
- **No external queue**: River uses PostgreSQL directly

## 🛡️ Security Updates

- **Latest Go runtime**: Security fixes included
- **No deprecated packages**: Reduced vulnerability surface
- **Modern crypto**: Using latest stdlib crypto
- **Secure defaults**: Better default configurations
- **Adult content isolation**: Separate PostgreSQL schema

## 📝 Code Style

### Modern Go Code
```go
// ✅ Good: Modern Go 1.22+ routing
mux.HandleFunc("GET /users/{id}", handleGetUser)

// ✅ Good: Structured logging with slog
logger.Info("user created", slog.String("id", id), slog.Int("age", age))

// ✅ Good: Error wrapping
return fmt.Errorf("failed to create user: %w", err)

// ✅ Good: Context-first
func GetUser(ctx context.Context, id string) (*User, error)

// ✅ Good: Go 1.25 WaitGroup.Go
var wg sync.WaitGroup
wg.Go(func() { processItem(ctx, item) })
```

### Deprecated Patterns
```go
// ❌ Bad: Old gorilla/mux
r.HandleFunc("/users/{id}", handleGetUser).Methods("GET")

// ❌ Bad: Printf-style logging
log.Printf("user created: id=%s, age=%d", id, age)

// ❌ Bad: Error messages without wrapping
return errors.New("failed to create user")

// ❌ Bad: No context
func GetUser(id string) (*User, error)

// ❌ Bad: Old WaitGroup pattern
wg.Add(1)
go func() { defer wg.Done(); processItem(ctx, item) }()
```

## 🔧 Development Tools

> **📋 Package Versions**: See [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#go-dependencies-development) for exact versions.

- **golangci-lint**: Latest linters
- **gopls**: Latest Go language server
- **govulncheck**: Security scanning
- **air**: Hot reload
- **migrate**: Database migrations
- **sqlc**: SQL code generation
- **ogen**: OpenAPI code generation

## 📚 Resources

- [Go 1.25 Release Notes](https://go.dev/doc/go1.25)
- [Go 1.22 Enhanced Routing](https://go.dev/blog/routing-enhancements)
- [log/slog Package](https://pkg.go.dev/log/slog)
- [koanf Documentation](https://github.com/knadh/koanf)
- [River Job Queue](https://riverqueue.com/docs)
- [ogen Documentation](https://ogen.dev/)
- [Typesense Go Client](https://github.com/typesense/typesense-go)

## ✅ Checklist

- [x] Go 1.25 (latest stable)
- [x] stdlib routing (no gorilla/mux)
- [x] slog/tint for dev, zap for prod
- [x] koanf for config (no viper)
- [x] Modern error handling
- [x] Context-first APIs
- [x] Type-safe patterns
- [x] No deprecated dependencies
- [x] Security-focused defaults
- [x] Performance optimized
- [x] River for job queue (infrastructure implemented)
- [x] ogen for API docs (planned)
- [x] Dragonfly for caching (client implemented)
- [x] Typesense for search (client implemented)
