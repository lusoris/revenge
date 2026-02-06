## Table of Contents

- [RBAC Service](#rbac-service)
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
- [Policy management (admin only)](#policy-management-admin-only)
- [Role management](#role-management)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# RBAC Service

<!-- DESIGN: services, README, test_output_claude, test_output_wiki -->


**Created**: 2026-01-31
**Status**: 🟡 In Progress
**Category**: service


> > Role-based access control with Casbin

**Package**: `internal/service/rbac`
**fx Module**: `rbac.Module`

---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | - |
| Sources | ✅ | - |
| Instructions | ✅ | - |
| Code | 🟡 Partial | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: 🟡 In Progress


---


## Architecture

```mermaid
flowchart LR
    subgraph Layer1["Layer 1"]
        node1(["Client<br/>(Web/App)"])
        node2["Middleware<br/>(Auth)"]
        node3["Casbin<br/>Enforcer"]
    end

    subgraph Layer2["Layer 2"]
        node4[("PostgreSQL<br/>casbin_rule")]
    end

    %% Connections
    node3 --> node4

    %% Styling
    style Layer1 fill:#1976D2,stroke:#1976D2,color:#fff
    style Layer2 fill:#388E3C,stroke:#388E3C,color:#fff
```

### Service Structure

```
internal/service/rbac/
├── module.go              # fx module (NewService, custom pgx Casbin adapter)
├── service.go             # Service struct + business logic (23 methods)
├── permissions.go         # 40+ permission constants (Resource + Action pairs)
├── roles.go               # Role definitions and defaults
├── adapter.go             # Custom Casbin pgx adapter (shared.casbin_rule table)
├── cached_service.go      # CachedService wrapping Service with cache layer
└── service_test.go        # Tests (🔴 not yet)
```

### Dependencies
**Go Packages**:
- `github.com/google/uuid`
- `github.com/casbin/casbin/v2` - RBAC policy engine
- `github.com/jackc/pgx/v5` - Custom pgx adapter (not casbin/pgx-adapter)
- `go.uber.org/fx`, `go.uber.org/zap`

**Internal Dependencies**:
- `internal/service/activity` - `activity.Logger` for audit logging
- `internal/infra/cache` - `cache.Cache` for CachedService

### Provides

`rbac.Module` provides: `NewService`, custom pgx adapter

## Implementation

### Key Interfaces (from code) ✅

```go
// Service is a concrete struct (not interface).
// Source: internal/service/rbac/service.go
type Service struct {
  enforcer       *casbin.Enforcer
  logger         *zap.Logger
  activityLogger activity.Logger
}

// Enforcement (2 methods)
func (s *Service) Enforce(sub, obj, act string) (bool, error)
func (s *Service) EnforceWithContext(ctx context.Context, userID uuid.UUID, resource, action string) (bool, error)

// Policy management (3 methods)
func (s *Service) AddPolicy(ctx context.Context, sub, obj, act string) error
func (s *Service) RemovePolicy(ctx context.Context, sub, obj, act string) error
func (s *Service) GetPolicies(ctx context.Context) ([][]string, error)

// Role management (15 methods)
func (s *Service) AssignRole(ctx context.Context, userID uuid.UUID, role string) error
func (s *Service) RemoveRole(ctx context.Context, userID uuid.UUID, role string) error
func (s *Service) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]string, error)
func (s *Service) GetUsersForRole(ctx context.Context, role string) ([]string, error)
func (s *Service) HasRole(ctx context.Context, userID uuid.UUID, role string) (bool, error)
func (s *Service) ListRoles(ctx context.Context) ([]Role, error)
func (s *Service) GetRole(ctx context.Context, name string) (*Role, error)
func (s *Service) CreateRole(ctx context.Context, role Role) error
func (s *Service) DeleteRole(ctx context.Context, name string) error
func (s *Service) UpdateRolePermissions(ctx context.Context, role string, permissions []Permission) error
func (s *Service) GetRolePermissions(ctx context.Context, role string) ([]Permission, error)
func (s *Service) AddPermissionToRole(ctx context.Context, role string, perm Permission) error
func (s *Service) RemovePermissionFromRole(ctx context.Context, role string, perm Permission) error
func (s *Service) GetAllRoleNames(ctx context.Context) ([]string, error)
func (s *Service) CheckUserPermission(ctx context.Context, userID uuid.UUID, resource, action string) (bool, error)

// Lifecycle (3 methods)
func (s *Service) LoadPolicy() error
func (s *Service) SavePolicy() error
func (s *Service) ListPermissions() []Permission
```

**Key Types**:
- `Role` - Role name + description + permissions
- `Permission` - `Resource` + `Action` pair
- `CachedService` - Cache wrapper using `cache.Cache`
- 40+ permission constants in `permissions.go` (users, profile, movies, libraries, playback, requests, settings, audit, integrations, notifications, admin)

## Configuration

### Current Config (from code) ✅

From `config.go` `RBACConfig` (koanf namespace `rbac.*`):
```yaml
rbac:
  model_path: /config/casbin_model.conf    # Casbin model file path
  policy_reload_interval: 5m               # Auto-reload interval
```

## API Endpoints
```
# Policy management (admin only)
GET    /api/v1/rbac/policies              # List policies
POST   /api/v1/rbac/policies              # Add policy
DELETE /api/v1/rbac/policies              # Remove policy

# Role management
POST   /api/v1/rbac/users/:id/roles       # Assign role
DELETE /api/v1/rbac/users/:id/roles/:role # Remove role
GET    /api/v1/rbac/users/:id/roles       # Get user roles
```

## Related Documentation
### Design Documents
- [services](INDEX.md)
- [01_ARCHITECTURE](../architecture/ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../architecture/DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../architecture/METADATA_SYSTEM.md)

### External Sources
- [Casbin](../../sources/security/casbin.md) - Auto-resolved from casbin
- [Uber fx](../../sources/tooling/fx.md) - Auto-resolved from fx
- [pgx PostgreSQL Driver](../../sources/database/pgx.md) - Auto-resolved from pgx
- [PostgreSQL Arrays](../../sources/database/postgresql-arrays.md) - Auto-resolved from postgresql-arrays
- [PostgreSQL JSON Functions](../../sources/database/postgresql-json.md) - Auto-resolved from postgresql-json
- [River Job Queue](../../sources/tooling/river.md) - Auto-resolved from river

