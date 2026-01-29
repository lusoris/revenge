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

### Server

| Variable | Config Key | Default | Description |
|----------|------------|---------|-------------|
| `REVENGE_HOST` | `server.host` | `0.0.0.0` | Listen address |
| `REVENGE_PORT` | `server.port` | `8096` | Listen port |
| `REVENGE_BASE_URL` | `server.base_url` | `/` | Base URL path |

### Database (PostgreSQL)

| Variable | Config Key | Default | Description |
|----------|------------|---------|-------------|
| `REVENGE_DB_HOST` | `database.host` | `localhost` | PostgreSQL host |
| `REVENGE_DB_PORT` | `database.port` | `5432` | PostgreSQL port |
| `REVENGE_DB_USER` | `database.user` | `revenge` | Database user |
| `REVENGE_DB_PASSWORD` | `database.password` | | Database password |
| `REVENGE_DB_NAME` | `database.name` | `revenge` | Database name |
| `REVENGE_DB_SSLMODE` | `database.sslmode` | `disable` | SSL mode |

### Cache (Dragonfly/Redis)

| Variable | Config Key | Default | Description |
|----------|------------|---------|-------------|
| `REVENGE_CACHE_URL` | `cache.addr` | `localhost:6379` | Cache address |
| `REVENGE_CACHE_PASSWORD` | `cache.password` | | Cache password |

### Search (Typesense)

| Variable | Config Key | Default | Description |
|----------|------------|---------|-------------|
| `REVENGE_TYPESENSE_URL` | `search.host`/`search.port` | `localhost:8108` | Typesense URL |
| `REVENGE_TYPESENSE_API_KEY` | `search.api_key` | | API key |

### Logging

| Variable | Config Key | Default | Description |
|----------|------------|---------|-------------|
| `REVENGE_LOG_LEVEL` | `logging.level` | `info` | Log level (debug, info, warn, error) |
| `REVENGE_LOG_FORMAT` | `logging.format` | `json` | Log format (json, text) |

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

## Related

- [Setup Guide](../operations/SETUP.md) - Production setup
- [Development Guide](../operations/DEVELOPMENT.md) - Development environment
- [koanf-configuration.instructions.md](../../../.github/instructions/koanf-configuration.instructions.md) - Configuration patterns
