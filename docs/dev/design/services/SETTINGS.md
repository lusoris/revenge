# Settings Service

> Server settings persistence and retrieval


<!-- TOC-START -->

## Table of Contents

- [Developer Resources](#developer-resources)
- [Status](#status)
- [Overview](#overview)
- [Categories](#categories)
- [Common Setting Keys](#common-setting-keys)
- [Operations](#operations)
  - [Get Settings](#get-settings)
  - [Set Settings](#set-settings)
  - [Delete Settings](#delete-settings)
  - [List Settings](#list-settings)
  - [Convenience Methods](#convenience-methods)
- [Storage Format](#storage-format)
- [Errors](#errors)
- [Implementation Checklist](#implementation-checklist)
  - [Phase 1: Core Infrastructure](#phase-1-core-infrastructure)
  - [Phase 2: Database](#phase-2-database)
  - [Phase 3: Service Layer](#phase-3-service-layer)
  - [Phase 4: API Integration](#phase-4-api-integration)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related Documents](#related-documents)

<!-- TOC-END -->

**Module**: `internal/service/settings`

## Developer Resources

> Package versions: [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#go-dependencies-core)

| Package | Purpose |
|---------|---------|
| encoding/json | JSON value marshaling |
| pgx | PostgreSQL JSONB storage |
| koanf | Configuration management |
| otter | Settings caching |

## Status

| Dimension | Status |
|-----------|--------|
| Design | ✅ |
| Sources | ✅ |
| Instructions | ✅ |
| Code | 🔴 |
| Linting | 🔴 |
| Unit Testing | 🔴 |
| Integration Testing | 🔴 |
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

## Implementation Checklist

### Phase 1: Core Infrastructure
- [ ] Create `internal/service/settings/` package structure
- [ ] Define setting keys as constants
- [ ] Create repository interface
- [ ] Add fx module wiring

### Phase 2: Database
- [ ] Create migration for `server_settings` table
- [ ] Add indexes (category, is_public)
- [ ] Write sqlc queries

### Phase 3: Service Layer
- [ ] Implement CRUD operations with caching
- [ ] Implement type-safe getters (string, bool, int)
- [ ] Add convenience methods for common settings
- [ ] Implement cache invalidation

### Phase 4: API Integration
- [ ] Define OpenAPI endpoints
- [ ] Generate ogen handlers
- [ ] Wire handlers to service
- [ ] Add admin authorization for private settings

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
| [PostgreSQL Arrays](https://www.postgresql.org/docs/current/arrays.html) | [Local](../../sources/database/postgresql-arrays.md) |
| [PostgreSQL JSON Functions](https://www.postgresql.org/docs/current/functions-json.html) | [Local](../../sources/database/postgresql-json.md) |
| [Uber fx](https://pkg.go.dev/go.uber.org/fx) | [Local](../../sources/tooling/fx.md) |
| [koanf](https://pkg.go.dev/github.com/knadh/koanf/v2) | [Local](../../sources/tooling/koanf.md) |
| [ogen OpenAPI Generator](https://pkg.go.dev/github.com/ogen-go/ogen) | [Local](../../sources/tooling/ogen.md) |
| [pgx PostgreSQL Driver](https://pkg.go.dev/github.com/jackc/pgx/v5) | [Local](../../sources/database/pgx.md) |
| [sqlc](https://docs.sqlc.dev/en/stable/) | [Local](../../sources/database/sqlc.md) |
| [sqlc Configuration](https://docs.sqlc.dev/en/stable/reference/config.html) | [Local](../../sources/database/sqlc-config.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Services](INDEX.md)

### In This Section

- [Activity Service](ACTIVITY.md)
- [Analytics Service](ANALYTICS.md)
- [API Keys Service](APIKEYS.md)
- [Auth Service](AUTH.md)
- [Fingerprint Service](FINGERPRINT.md)
- [Grants Service](GRANTS.md)
- [Library Service](LIBRARY.md)
- [Metadata Service](METADATA.md)

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

## Related Documents

- [Configuration](../technical/CONFIGURATION.md) - File-based config
- [Auth Service](AUTH.md) - Registration settings
- [User Service](USER.md) - User preferences
- [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md) - Configuration keys reference
