# RBAC Service

> Role-based access control with Casbin

**Location**: `internal/service/rbac/`

---

## Overview

The RBAC service provides permission management using Casbin:

- Policy enforcement
- Role assignment
- Permission checking
- Dynamic policy updates

---

## Model

Uses RBAC with resource/action model:

```ini
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
```

---

## Operations

### Check Permission

```go
func (s *Service) Enforce(ctx context.Context, sub, obj, act string) (bool, error)
```

Example:
```go
allowed, err := rbac.Enforce(ctx, "user:123", "library:456", "read")
```

### Role Management

```go
// Add role to user
func (s *Service) AddRoleForUser(ctx context.Context, user, role string) error

// Remove role from user
func (s *Service) DeleteRoleForUser(ctx context.Context, user, role string) error

// Get user's roles
func (s *Service) GetRolesForUser(ctx context.Context, user string) ([]string, error)
```

### Policy Management

```go
// Add policy
func (s *Service) AddPolicy(ctx context.Context, sub, obj, act string) error

// Remove policy
func (s *Service) RemovePolicy(ctx context.Context, sub, obj, act string) error

// Get policies for subject
func (s *Service) GetPoliciesForSubject(ctx context.Context, sub string) ([][]string, error)
```

---

## Default Roles

| Role | Permissions |
|------|-------------|
| `admin` | Full access to all resources |
| `moderator` | Manage content, limited user management |
| `user` | Access own libraries, standard features |
| `guest` | Read-only access to public content |

---

## Resource Types

| Resource | Actions |
|----------|---------|
| `library:*` | read, write, delete, manage |
| `user:*` | read, write, delete |
| `settings` | read, write |
| `admin` | access |

---

## Related

- [RBAC_CASBIN.md](../features/shared/RBAC_CASBIN.md) - Full design doc
- [User Service](USER.md) - User management
- [Library Service](LIBRARY.md) - Library access
