# Configuration Reference

> Complete configuration options for Revenge

**Location**: `internal/config/config.go`

---

## Overview

Configuration is loaded via [koanf v2](https://github.com/knadh/koanf) from:

1. **YAML file** (optional) - `config.yaml`, `configs/config.yaml`, or `/etc/revenge/config.yaml`
2. **Environment variables** - `REVENGE_` prefix, overrides file settings

---

## Configuration File

```yaml
# config.yaml
server:
  host: 0.0.0.0
  port: 8096
  base_url: /
  read_timeout: 30s
  write_timeout: 30s

database:
  host: localhost
  port: 5432
  user: revenge
  password: changeme
  name: revenge
  sslmode: disable
  max_conns: 25
  min_conns: 5

cache:
  addr: localhost:6379
  password: ""
  db: 0
  local_capacity: 10000
  local_ttl: 300
  api_capacity: 5000
  api_ttl: 3600

search:
  host: localhost
  port: 8108
  api_key: xyz

auth:
  jwt_secret: your-secret-key
  session_duration: 24

metadata:
  radarr:
    base_url: http://localhost:7878
    api_key: your-radarr-key
  tmdb:
    api_key: your-tmdb-key

modules:
  movie: true
  tvshow: true
  music: true
  audiobook: false
  book: false
  podcast: false
  photo: false
  livetv: false
  comics: false
  adult: false

logging:
  level: info
  format: json
```

---

## Environment Variables

All settings can be set via environment variables with `REVENGE_` prefix.

> **📋 Complete Environment Variable Reference**: See [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#environment-variable-mapping) for the authoritative list of all environment variables with platform-specific mappings (Docker Compose, K8s, K3s, Swarm).

The sections below document the **config structure and Go types**. For the complete env var list, always refer to SOURCE_OF_TRUTH.

---

## Configuration Sections

### Server Configuration

```go
type ServerConfig struct {
    Host              string        // Listen address
    Port              int           // Listen port
    BaseURL           string        // Base URL path
    ReadTimeout       time.Duration // HTTP read timeout
    WriteTimeout      time.Duration // HTTP write timeout
    IdleTimeout       time.Duration // HTTP idle timeout
    ReadHeaderTimeout time.Duration // Header read timeout
    MaxHeaderBytes    int           // Max header size
}
```

**Defaults**:
- `read_timeout`: 30s
- `write_timeout`: 30s
- `idle_timeout`: 60s
- `read_header_timeout`: 5s
- `max_header_bytes`: 1MB

### Database Configuration

```go
type DatabaseConfig struct {
    Host     string // PostgreSQL host
    Port     int    // PostgreSQL port
    User     string // Database user
    Password string // Database password
    Name     string // Database name
    SSLMode  string // SSL mode
    MaxConns int32  // Max connections
    MinConns int32  // Min connections
}
```

**Defaults**:
- `max_conns`: 25
- `min_conns`: 5
- `sslmode`: disable

### Cache Configuration

```go
type CacheConfig struct {
    Addr     string // Redis/Dragonfly address
    Password string // Cache password
    DB       int    // Redis database number

    // Local cache (otter)
    LocalCapacity int // Max entries (default: 10000)
    LocalTTL      int // TTL in seconds (default: 300)

    // API cache (sturdyc)
    APICapacity  int // Max entries (default: 5000)
    APINumShards int // Number of shards (default: 10)
    APITTL       int // TTL in seconds (default: 3600)
}
```

### Search Configuration

```go
type SearchConfig struct {
    Host   string // Typesense host
    Port   int    // Typesense port
    APIKey string // Typesense API key
}
```

### Auth Configuration

```go
type AuthConfig struct {
    JWTSecret       string // JWT signing secret
    SessionDuration int    // Session duration in hours (default: 24)
}
```

### Metadata Configuration

```go
type MetadataConfig struct {
    Radarr RadarrConfig
    TMDb   TMDbConfig
}

type RadarrConfig struct {
    BaseURL string // Radarr base URL
    APIKey  string // Radarr API key
}

type TMDbConfig struct {
    APIKey     string // TMDb API key
    BaseURL    string // API base URL
    ImageURL   string // Image base URL
    Timeout    int    // Request timeout (seconds)
    CacheTTL   int    // Cache TTL (seconds)
    CacheSize  int    // Max cache entries
    RetryCount int    // Max retries
}
```

### Modules Configuration

```go
type ModulesConfig struct {
    Movie     bool // Movies (default: true)
    TVShow    bool // TV Shows (default: true)
    Music     bool // Music (default: true)
    Audiobook bool // Audiobooks (default: false)
    Book      bool // Books (default: false)
    Podcast   bool // Podcasts (default: false)
    Photo     bool // Photos (default: false)
    LiveTV    bool // Live TV (default: false)
    Comics    bool // Comics (default: false)
    Adult     bool // Adult content (default: false, explicit opt-in)
}
```

### Logging Configuration

```go
type LoggingConfig struct {
    Level  string // Log level: debug, info, warn, error
    Format string // Output format: json, text
}
```

---

## Loading Order

1. Load YAML file (if exists)
2. Apply environment variables (override file values)
3. Apply environment aliases (`REVENGE_DB_*` → `database.*`)
4. Set defaults for missing values

---

## fx Module

Configuration is provided via fx dependency injection:

```go
import "github.com/lusoris/revenge/internal/config"

fx.New(
    config.Module,
    // ... other modules
)
```

---


<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

### Referenced Sources

| Source | Documentation |
|--------|---------------|
| [koanf](https://pkg.go.dev/github.com/knadh/koanf/v2) | [Local](../../sources/tooling/koanf.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Technical](INDEX.md)

### In This Section

- [API Reference](API.md)
- [Revenge - Audio Streaming & Progress Tracking](AUDIO_STREAMING.md)
- [Revenge - Frontend Architecture](FRONTEND.md)
- [Revenge - Advanced Offloading Architecture](OFFLOADING.md)
- [Revenge - Technology Stack](TECH_STACK.md)

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

## Related

- [Setup Guide](../operations/SETUP.md) - Production setup
- [Development Guide](../operations/DEVELOPMENT.md) - Development environment
- [koanf-configuration.instructions.md](../../../.github/instructions/koanf-configuration.instructions.md) - Configuration patterns
