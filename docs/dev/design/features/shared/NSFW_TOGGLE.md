## Table of Contents

- [Revenge - NSFW Toggle](#revenge-nsfw-toggle)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Database Schema](#database-schema)
    - [Module Structure](#module-structure)
    - [Component Interaction](#component-interaction)
  - [Implementation](#implementation)
    - [Key Interfaces](#key-interfaces)
    - [Dependencies](#dependencies)
  - [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
    - [Config Keys](#config-keys)
  - [API Endpoints](#api-endpoints)
    - [Content Management](#content-management)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# Revenge - NSFW Toggle


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: feature


> Content module for 

> User preference component for adult content visibility. > Referenced by [WHISPARR_STASHDB_SCHEMA.md](../adult/WHISPARR_STASHDB_SCHEMA.md) and [ADULT_CONTENT_SYSTEM.md](../adult/ADULT_CONTENT_SYSTEM.md).

---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | - |
| Sources | 🟡 | - |
| Instructions | ✅ | - |
| Code | 🔴 | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: ✅ Complete



---


## Architecture

```mermaid
flowchart TD
    node1([Client<br/>(Web/App)])
    node2[[API Handler<br/>(ogen)]]
    node3[[Service<br/>(Logic)]]
    node4["Repository<br/>(sqlc)"]
    node5[[Metadata<br/>Service]]
    node6[(Cache<br/>(otter))]
    node7[(PostgreSQL<br/>(pgx))]
    node8([External<br/>APIs])
    node1 --> node2
    node2 --> node3
    node4 --> node5
    node5 --> node6
    node7 --> node8
    node3 --> node4
    node6 --> node7
```

### Database Schema

**Schema**: `public`

<!-- Schema diagram -->

### Module Structure

```
internal/content/revenge___nsfw_toggle/
├── module.go              # fx module definition
├── repository.go          # Database operations
├── service.go             # Business logic
├── handler.go             # HTTP handlers (ogen)
├── types.go               # Domain types
└── revenge___nsfw_toggle_test.go
```

### Component Interaction

<!-- Component interaction diagram -->
## Implementation

### Key Interfaces

```go
type NSFWService interface {
  EnableNSFW(ctx context.Context, userID uuid.UUID, pin string) error
  DisableNSFW(ctx context.Context, userID uuid.UUID) error
  IsNSFWEnabled(ctx context.Context, userID uuid.UUID) (bool, error)
  SetPIN(ctx context.Context, userID uuid.UUID, pin string) error
  VerifyPIN(ctx context.Context, userID uuid.UUID, pin string) (bool, error)
}
```


### Dependencies
**Go Packages**:
- `github.com/google/uuid`
- `github.com/jackc/pgx/v5`
- `golang.org/x/crypto/bcrypt` - PIN hashing
- `go.uber.org/fx`

## Configuration

### Environment Variables

```bash
NSFW_DEFAULT_SESSION_TIMEOUT=60
```


### Config Keys
```yaml
nsfw:
  default_session_timeout_minutes: 60
  require_pin_by_default: true
```

## API Endpoints

### Content Management
```
POST /api/v1/nsfw/enable          # Enable NSFW mode
POST /api/v1/nsfw/disable         # Disable NSFW mode
GET  /api/v1/nsfw/status          # Get NSFW status
PUT  /api/v1/nsfw/pin             # Set/update PIN
```

## Related Documentation
### Design Documents
- [01_ARCHITECTURE](../../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../../architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../../architecture/03_METADATA_SYSTEM.md)

### External Sources
- [Casbin](../../../sources/security/casbin.md) - Auto-resolved from casbin
- [Dragonfly Documentation](../../../sources/infrastructure/dragonfly.md) - Auto-resolved from dragonfly
- [Uber fx](../../../sources/tooling/fx.md) - Auto-resolved from fx
- [ogen OpenAPI Generator](../../../sources/tooling/ogen.md) - Auto-resolved from ogen
- [pgx PostgreSQL Driver](../../../sources/database/pgx.md) - Auto-resolved from pgx
- [PostgreSQL Arrays](../../../sources/database/postgresql-arrays.md) - Auto-resolved from postgresql-arrays
- [PostgreSQL JSON Functions](../../../sources/database/postgresql-json.md) - Auto-resolved from postgresql-json
- [River Job Queue](../../../sources/tooling/river.md) - Auto-resolved from river
- [rueidis](../../../sources/tooling/rueidis.md) - Auto-resolved from rueidis
- [rueidis GitHub README](../../../sources/tooling/rueidis-guide.md) - Auto-resolved from rueidis-docs
- [sqlc](../../../sources/database/sqlc.md) - Auto-resolved from sqlc
- [sqlc Configuration](../../../sources/database/sqlc-config.md) - Auto-resolved from sqlc-config
- [Svelte 5 Runes](../../../sources/frontend/svelte-runes.md) - Auto-resolved from svelte-runes
- [Svelte 5 Documentation](../../../sources/frontend/svelte5.md) - Auto-resolved from svelte5
- [SvelteKit Documentation](../../../sources/frontend/sveltekit.md) - Auto-resolved from sveltekit

