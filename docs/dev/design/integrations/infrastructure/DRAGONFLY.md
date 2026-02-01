## Table of Contents

- [Dragonfly](#dragonfly)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Integration Structure](#integration-structure)
    - [Data Flow](#data-flow)
    - [Provides](#provides)
  - [Implementation](#implementation)
    - [Key Interfaces](#key-interfaces)
    - [Dependencies](#dependencies)
  - [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
- [Dragonfly connection](#dragonfly-connection)
- [L1 cache (otter)](#l1-cache-otter)
- [L2 cache (Dragonfly)](#l2-cache-dragonfly)
    - [Config Keys](#config-keys)
  - [API Endpoints](#api-endpoints)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# Dragonfly


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: integration


> Integration with Dragonfly

> High-performance Redis-compatible cache
**Authentication**: password

---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | - |
| Sources | ✅ | - |
| Instructions | ✅ | - |
| Code | 🔴 | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: ✅ Complete



---


## Architecture

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│   Server    │────▶│   rueidis    │────▶│  Dragonfly  │
│ (Services)  │◀────│ (Redis Client│◀────│   Server    │
└─────────────┘     │  with Auto-  │     └─────────────┘
                    │   Pipelining)│
                    └──────────────┘
                           │
                    ┌──────┴───────┐
                    │   sturdyc    │
                    │ (Coalescing) │
                    └──────────────┘
```


### Integration Structure

```
internal/integration/dragonfly/
├── client.go              # API client
├── types.go               # Response types
├── mapper.go              # Map external → internal types
├── cache.go               # Response caching
└── client_test.go         # Tests
```

### Data Flow

<!-- Data flow diagram -->

### Provides
<!-- Data provided by integration -->


## Implementation

### Key Interfaces

```go
// Cache interface (unified L1 + L2)
type Cache interface {
  Get(ctx context.Context, key string) ([]byte, error)
  Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
  Delete(ctx context.Context, key string) error
  Exists(ctx context.Context, key string) (bool, error)
  Invalidate(ctx context.Context, pattern string) error
}

// Configuration
type DragonflyConfig struct {
  Addresses  []string      `yaml:"addresses"`
  Password   string        `yaml:"password"`
  DB         int           `yaml:"db"`
  PoolSize   int           `yaml:"pool_size"`
  MaxRetries int           `yaml:"max_retries"`
  TLS        bool          `yaml:"tls"`
}

// L1 Cache Config (otter)
type L1Config struct {
  MaxSize  int           `yaml:"max_size"`
  TTL      time.Duration `yaml:"ttl"`
}
```


### Dependencies
**Go Packages**:
- `github.com/redis/rueidis` - High-performance Redis client
- `github.com/maypok86/otter` - L1 in-memory cache
- `github.com/viccon/sturdyc` - Request coalescing
- `go.uber.org/fx` - Dependency injection

**External Services**:
- Dragonfly server (recommended) or Redis 6.0+ (fallback)







## Configuration

### Environment Variables

```bash
# Dragonfly connection
DRAGONFLY_ADDRESSES=localhost:6379
DRAGONFLY_PASSWORD=secret
DRAGONFLY_DB=0
DRAGONFLY_POOL_SIZE=10
DRAGONFLY_TLS=false

# L1 cache (otter)
CACHE_L1_MAX_SIZE=10000
CACHE_L1_TTL=5m

# L2 cache (Dragonfly)
CACHE_L2_TTL=1h
```


### Config Keys
```yaml
cache:
  dragonfly:
    addresses:
      - localhost:6379
    password: ${DRAGONFLY_PASSWORD}
    db: 0
    pool_size: 10
    max_retries: 3
    tls: false

  l1:
    max_size: 10000
    ttl: 5m

  l2:
    ttl: 1h
```



## API Endpoints
**Health Check**:
```
GET /api/v1/health/cache
```

**Response**:
```json
{
  "status": "healthy",
  "dragonfly": {
    "connected": true,
    "version": "1.15.0"
  },
  "l1_stats": {
    "size": 4523,
    "hit_rate": 0.87,
    "evictions": 123
  },
  "l2_stats": {
    "hit_rate": 0.72
  }
}
```

**Cache Stats**:
```
GET /api/v1/admin/cache/stats
```








## Related Documentation
### Design Documents
- [01_ARCHITECTURE](../../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../../architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../../architecture/03_METADATA_SYSTEM.md)

### External Sources
- [Dragonfly Documentation](../../../sources/infrastructure/dragonfly.md) - Auto-resolved from dragonfly
- [pgx PostgreSQL Driver](../../../sources/database/pgx.md) - Auto-resolved from pgx
- [PostgreSQL Arrays](../../../sources/database/postgresql-arrays.md) - Auto-resolved from postgresql-arrays
- [PostgreSQL JSON Functions](../../../sources/database/postgresql-json.md) - Auto-resolved from postgresql-json
- [Prometheus Go Client](../../../sources/observability/prometheus.md) - Auto-resolved from prometheus
- [Prometheus Metric Types](../../../sources/observability/prometheus-metrics.md) - Auto-resolved from prometheus-metrics
- [River Job Queue](../../../sources/tooling/river.md) - Auto-resolved from river
- [rueidis](../../../sources/tooling/rueidis.md) - Auto-resolved from rueidis
- [rueidis GitHub README](../../../sources/tooling/rueidis-guide.md) - Auto-resolved from rueidis-docs
- [Typesense API](../../../sources/infrastructure/typesense.md) - Auto-resolved from typesense
- [Typesense Go Client](../../../sources/infrastructure/typesense-go.md) - Auto-resolved from typesense-go

