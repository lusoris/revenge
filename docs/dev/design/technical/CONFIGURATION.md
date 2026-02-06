# Configuration

<!-- DESIGN: technical -->

**Package**: `internal/config`
**Loader**: koanf (YAML file + environment variables)
**Default path**: `config/revenge.yaml`

> Application configuration with defaults, environment overrides, and validation

---

## Config Struct

Top-level `Config` with koanf namespaces:

```go
type Config struct {
    Server       ServerConfig       `koanf:"server"`
    Database     DatabaseConfig     `koanf:"database"`
    Cache        CacheConfig        `koanf:"cache"`
    Search       SearchConfig       `koanf:"search"`
    Jobs         JobsConfig         `koanf:"jobs"`
    Logging      LoggingConfig      `koanf:"logging"`
    Auth         AuthConfig         `koanf:"auth"`
    Session      SessionConfig      `koanf:"session"`
    RBAC         RBACConfig         `koanf:"rbac"`
    Movie        MovieConfig        `koanf:"movie"`
    Integrations IntegrationsConfig `koanf:"integrations"`
    Legacy       LegacyConfig       `koanf:"legacy"`
    Email        EmailConfig        `koanf:"email"`
    Avatar       AvatarConfig       `koanf:"avatar"`
    Storage      StorageConfig      `koanf:"storage"`
    Activity     ActivityConfig     `koanf:"activity"`
    Raft         RaftConfig         `koanf:"raft"`
}
```

## Sections and Defaults

### server.*

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| host | string | `0.0.0.0` | Bind address |
| port | int | `8080` | Server port |
| read_timeout | duration | `30s` | Max request read time |
| write_timeout | duration | `30s` | Max response write time |
| idle_timeout | duration | `120s` | Keep-alive idle timeout |
| shutdown_timeout | duration | `10s` | Graceful shutdown timeout |
| rate_limit.enabled | bool | `true` | Enable rate limiting |
| rate_limit.backend | string | `memory` | `memory` or `redis` |
| rate_limit.global.requests_per_second | float64 | `10.0` | Global RPS per IP |
| rate_limit.global.burst | int | `20` | Global burst |
| rate_limit.auth.requests_per_second | float64 | `1.0` | Auth RPS per IP |
| rate_limit.auth.burst | int | `5` | Auth burst |

### database.*

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| url | string | `postgres://revenge:changeme@localhost:5432/revenge?sslmode=disable` | PostgreSQL URL |
| max_conns | int | `0` | Max pool size (0 = CPU*2+1) |
| min_conns | int | `2` | Min pool size |
| max_conn_lifetime | duration | `30m` | Max connection age |
| max_conn_idle_time | duration | `5m` | Max idle time |
| health_check_period | duration | `30s` | Health check interval |

### cache.*

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| url | string | `""` | Dragonfly/Redis URL |
| enabled | bool | `false` | Enable distributed cache |

### search.*

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| url | string | `""` | Typesense URL |
| api_key | string | `""` | Typesense API key |
| enabled | bool | `false` | Enable search |

### jobs.*

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| max_workers | int | `100` | Max concurrent River workers |
| fetch_cooldown | duration | `200ms` | Cooldown between fetches |
| fetch_poll_interval | duration | `2s` | Poll interval for new jobs |
| rescue_stuck_jobs_after | duration | `30m` | Rescue stuck jobs after |

### logging.*

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| level | string | `info` | `debug`, `info`, `warn`, `error` |
| format | string | `text` | `text` (dev) or `json` (prod) |
| development | bool | `false` | Enable dev mode (pretty printing) |

### auth.*

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| jwt_secret | string | `""` | JWT signing secret (min 32 chars) |
| jwt_expiry | duration | `24h` | JWT token lifetime |
| refresh_expiry | duration | `168h` | Refresh token lifetime (7 days) |
| lockout_threshold | int | `5` | Failed attempts before lockout |
| lockout_window | duration | `15m` | Lockout counting window |
| lockout_enabled | bool | `true` | Enable account lockout |

### session.*

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| cache_enabled | bool | `true` | Cache sessions in Dragonfly |
| cache_ttl | duration | `5m` | Session cache TTL |
| max_per_user | int | `10` | Max active sessions per user |
| token_length | int | `32` | Session token length (bytes) |

### rbac.*

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| model_path | string | `config/casbin_model.conf` | Casbin model file |
| policy_reload_interval | duration | `5m` | Policy reload interval |

### integrations.radarr.* / integrations.sonarr.*

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| enabled | bool | `false` | Enable integration |
| base_url | string | `http://localhost:7878` (radarr) | Server URL |
| api_key | string | `""` | API key |
| auto_sync | bool | `false` | Automatic library sync |
| sync_interval | int | `300` | Sync interval (seconds) |

### email.*

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| enabled | bool | `false` | Enable email sending |
| provider | string | `smtp` | `smtp` or `sendgrid` |
| from_address | string | `""` | Sender email |
| from_name | string | `Revenge Media Server` | Sender name |
| base_url | string | `http://localhost:8080` | App URL for email links |
| smtp.host | string | `""` | SMTP host |
| smtp.port | int | `587` | SMTP port |
| smtp.use_starttls | bool | `true` | STARTTLS |
| sendgrid.api_key | string | `""` | SendGrid API key |

### storage.*

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| backend | string | `local` | `local` or `s3` |
| local.path | string | `/data/storage` | Local storage path |
| s3.endpoint | string | `""` | S3/MinIO endpoint |
| s3.region | string | `us-east-1` | S3 region |
| s3.bucket | string | `""` | S3 bucket |
| s3.use_path_style | bool | `false` | Path-style URLs (MinIO) |

### Other Sections

| Namespace | Key Defaults |
|-----------|-------------|
| movie.tmdb.* | api_key: "", rate_limit: 40, cache_ttl: 5m |
| movie.library.* | paths: [], scan_interval: 0s |
| legacy.* | enabled: false, require_pin: true, audit_all_access: true |
| avatar.* | storage_path: /data/avatars, max_size: 2MB, types: jpeg/png/webp |
| activity.* | retention_days: 90 |
| raft.* | enabled: false, bind_addr: 0.0.0.0:7000 |

## Loading

```go
var Module = fx.Module("config",
    fx.Provide(ProvideConfig),
)

func ProvideConfig() (*Config, error) {
    cfg, err := Load(DefaultConfigPath)
    // Environment variables override file values
    return cfg, err
}
```

## Related Documentation

- [../operations/DEVELOPMENT.md](../operations/DEVELOPMENT.md) - Development setup
- [../operations/SETUP.md](../operations/SETUP.md) - Production deployment
