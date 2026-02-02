## Table of Contents

- [Grants Service](#grants-service)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Service Structure](#service-structure)
    - [Dependencies](#dependencies)
    - [Provides](#provides)
    - [Component Diagram](#component-diagram)
  - [Implementation](#implementation)
    - [Key Interfaces](#key-interfaces)
    - [Dependencies](#dependencies)
  - [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
    - [Config Keys](#config-keys)
  - [API Endpoints](#api-endpoints)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# Grants Service


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: service


> > Polymorphic resource access grants for fine-grained sharing

**Package**: `internal/service/grants`
**fx Module**: `grants.Module`

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

```mermaid
flowchart TD
    node1([Client<br/>[Web/App]])
    node2[[API Handler<br/>[ogen]]]
    node3[[Service<br/>[Logic]]]
    node4["Repository<br/>[sqlc]"]
    node5[[RBAC<br/>Service]]
    node6[(Cache<br/>[otter])]
    node7[(PostgreSQL<br/>[pgx])]
    node1 --> node2
    node2 --> node3
    node4 --> node5
    node5 --> node6
    node3 --> node4
    node6 --> node7
```

### Service Structure

```
internal/service/grants/
├── module.go              # fx module definition
├── service.go             # Service implementation
├── repository.go          # Data access (if needed)
├── handler.go             # HTTP handlers (if exposed)
├── middleware.go          # Middleware (if needed)
├── types.go               # Domain types
└── service_test.go        # Tests
```

### Dependencies
**Go Packages**:
- `github.com/google/uuid`
- `github.com/jackc/pgx/v5`
- `github.com/riverqueue/river` - Cleanup jobs
- `github.com/maypok86/otter` - Grant cache
- `go.uber.org/fx`


### Provides
<!-- Service provides -->

### Component Diagram

<!-- Component diagram -->
## Implementation

### Key Interfaces

```go
type GrantsService interface {
  // Grant management
  GrantAccess(ctx context.Context, req GrantRequest) (*AccessGrant, error)
  RevokeGrant(ctx context.Context, grantID uuid.UUID) error
  CheckAccess(ctx context.Context, userID uuid.UUID, resourceType string, resourceID uuid.UUID, permission string) (bool, error)

  // List grants
  GetUserGrants(ctx context.Context, userID uuid.UUID) ([]AccessGrant, error)
  GetResourceGrants(ctx context.Context, resourceType string, resourceID uuid.UUID) ([]AccessGrant, error)

  // Cleanup
  CleanupExpiredGrants(ctx context.Context) (int, error)
}

type GrantRequest struct {
  GrantedToUserID uuid.UUID  `json:"granted_to_user_id"`
  ResourceType    string     `json:"resource_type"`
  ResourceID      uuid.UUID  `json:"resource_id"`
  Permission      string     `json:"permission"`
  ExpiresAt       *time.Time `json:"expires_at,omitempty"`
}
```


### Dependencies
**Go Packages**:
- `github.com/google/uuid`
- `github.com/jackc/pgx/v5`
- `github.com/riverqueue/river` - Cleanup jobs
- `github.com/maypok86/otter` - Grant cache
- `go.uber.org/fx`

## Configuration

### Environment Variables

```bash
GRANTS_CLEANUP_INTERVAL=1h
GRANTS_DEFAULT_EXPIRY=168h  # 7 days
```


### Config Keys
```yaml
grants:
  cleanup_interval: 1h
  default_expiry: 168h
```

## API Endpoints
```
POST   /api/v1/grants                       # Grant access
DELETE /api/v1/grants/:id                   # Revoke grant
GET    /api/v1/grants/me                    # Get my grants
GET    /api/v1/grants/resource/:type/:id    # Get resource grants
```

## Related Documentation
### Design Documents
- [services](INDEX.md)
- [01_ARCHITECTURE](../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../architecture/03_METADATA_SYSTEM.md)

### External Sources
- [Casbin](../../sources/security/casbin.md) - Auto-resolved from casbin
- [Uber fx](../../sources/tooling/fx.md) - Auto-resolved from fx
- [pgx PostgreSQL Driver](../../sources/database/pgx.md) - Auto-resolved from pgx
- [PostgreSQL Arrays](../../sources/database/postgresql-arrays.md) - Auto-resolved from postgresql-arrays
- [PostgreSQL JSON Functions](../../sources/database/postgresql-json.md) - Auto-resolved from postgresql-json
- [River Job Queue](../../sources/tooling/river.md) - Auto-resolved from river

