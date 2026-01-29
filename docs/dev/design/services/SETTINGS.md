# Settings Service

> Server settings persistence and retrieval

**Location**: `internal/service/settings/`

---

## Overview

The Settings service manages server-wide configuration:

- Key-value settings storage
- Category-based organization
- Type-safe getters
- Public/private visibility

---

## Categories

```go
const (
    CategoryGeneral  = "general"
    CategorySecurity = "security"
    CategoryMedia    = "media"
    CategoryCache    = "cache"
    CategorySearch   = "search"
    CategoryAdult    = "adult"
)
```

## Common Setting Keys

```go
const (
    KeyServerName                = "server.name"
    KeyServerVersion             = "server.version"
    KeyServerTimezone            = "server.timezone"
    KeySecurityRequireAuth       = "security.require_authentication"
    KeySecurityAllowRegistration = "security.allow_registration"
    KeyMediaDefaultProfile       = "media.default_transcoding_profile"
    KeyMediaHWAccel              = "media.enable_hardware_acceleration"
    KeyAdultGloballyEnabled      = "adult.globally_enabled"
)
```

---

## Operations

### Get Settings

```go
// Get raw setting
func (s *Service) Get(ctx context.Context, key string) (*db.ServerSetting, error)

// Get and unmarshal to destination
func (s *Service) GetValue(ctx context.Context, key string, dest any) error

// Type-safe getters
func (s *Service) GetString(ctx context.Context, key string) (string, error)
func (s *Service) GetBool(ctx context.Context, key string) (bool, error)
func (s *Service) GetInt(ctx context.Context, key string) (int, error)
```

### Set Settings

```go
type SetParams struct {
    Key         string
    Value       any      // JSON-serializable
    Category    string
    Description *string
    IsPublic    *bool
}

func (s *Service) Set(ctx context.Context, params SetParams) (*db.ServerSetting, error)
```

### Delete Settings

```go
func (s *Service) Delete(ctx context.Context, key string) error
```

### List Settings

```go
func (s *Service) ListAll(ctx context.Context) ([]db.ServerSetting, error)
func (s *Service) ListByCategory(ctx context.Context, category string) ([]db.ServerSetting, error)
func (s *Service) ListPublic(ctx context.Context) ([]db.ServerSetting, error)
```

### Convenience Methods

```go
func (s *Service) GetServerName(ctx context.Context) (string, error)
func (s *Service) IsRegistrationAllowed(ctx context.Context) (bool, error)
func (s *Service) IsAuthRequired(ctx context.Context) (bool, error)
func (s *Service) IsAdultContentGloballyEnabled(ctx context.Context) (bool, error)
```

---

## Storage Format

Settings are stored as JSON values in PostgreSQL:

```sql
CREATE TABLE server_settings (
    key         TEXT PRIMARY KEY,
    value       JSONB NOT NULL,
    category    TEXT NOT NULL,
    description TEXT,
    is_public   BOOLEAN DEFAULT FALSE,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

## Errors

| Error | Description |
|-------|-------------|
| `ErrSettingNotFound` | Setting key does not exist |

---

## Related

- [Configuration](../technical/CONFIGURATION.md) - File-based config
- [Auth Service](AUTH.md) - Registration settings
