# Dynamic RBAC with Casbin

> Role-Based Access Control using Casbin for dynamic permission management.

## Overview

Revenge uses [Casbin](https://casbin.org/) for dynamic Role-Based Access Control (RBAC). This allows administrators to:

- Create custom roles via the admin interface
- Assign granular permissions to roles at runtime
- Modify role-permission mappings without code changes or restarts

## Architecture

### Components

| Component | Description |
|-----------|-------------|
| `CasbinService` | Core RBAC service in `internal/service/rbac/casbin.go` |
| `casbin_rules` | PostgreSQL table storing Casbin policies |
| `roles` | Custom roles table for role metadata (name, description, color, etc.) |
| `permission_definitions` | Reference table of available permissions for UI |
| `RBACMiddleware` | HTTP middleware for permission checks |

### Data Flow

```
Request → Auth Middleware → RBAC Middleware → Handler
                              ↓
                         CasbinService
                              ↓
                    Casbin Enforcer (cached)
                              ↓
                    casbin_rules (PostgreSQL)
```

## Database Schema

### roles Table

Stores role definitions for admin management:

```sql
CREATE TABLE roles (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(100) NOT NULL UNIQUE,  -- e.g., "editor", "viewer"
    display_name    VARCHAR(255) NOT NULL,          -- e.g., "Content Editor"
    description     TEXT,
    color           VARCHAR(7),                      -- Hex color for UI
    icon            VARCHAR(50),                     -- Icon name for UI
    is_system       BOOLEAN NOT NULL DEFAULT false,  -- Cannot delete system roles
    is_default      BOOLEAN NOT NULL DEFAULT false,  -- Default for new users
    priority        INT NOT NULL DEFAULT 0,          -- UI sorting
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      UUID REFERENCES users(id)
);
```

### System Roles (Cannot Be Deleted)

| Role | Description | Priority |
|------|-------------|----------|
| `admin` | Full access to all features | 1000 |
| `moderator` | Manage libraries, metadata, content | 500 |
| `user` | Browse, play, rate, playlists | 100 |
| `guest` | Browse only | 0 |

### casbin_rules Table (Managed by Casbin)

Stores permission policies:

```sql
-- Automatically created by Casbin adapter
CREATE TABLE casbin_rules (
    id SERIAL PRIMARY KEY,
    ptype VARCHAR(10),  -- "p" for policy, "g" for grouping
    v0 VARCHAR(256),    -- role name
    v1 VARCHAR(256),    -- permission name
    v2 VARCHAR(256),    -- action (always "allow")
    v3 VARCHAR(256),
    v4 VARCHAR(256),
    v5 VARCHAR(256)
);
```

### permission_definitions Table

Reference table for UI to display available permissions:

```sql
CREATE TABLE permission_definitions (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(100) NOT NULL UNIQUE,    -- e.g., "content.browse"
    display_name    VARCHAR(255) NOT NULL,           -- e.g., "Browse Content"
    description     TEXT NOT NULL,
    category        VARCHAR(50) NOT NULL,            -- e.g., "Content"
    is_dangerous    BOOLEAN NOT NULL DEFAULT false,  -- Needs extra confirmation
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

## Permission Categories

| Category | Description |
|----------|-------------|
| `System` | Server settings, logs, jobs, API keys, roles |
| `Users` | User management |
| `Libraries` | Library management |
| `Content` | Content browsing and metadata |
| `Playback` | Streaming, downloading, transcoding |
| `Social` | Ratings, playlists, collections, history |
| `Adult` | Adult content access (schema `c`) |

## Available Permissions

### System Permissions

| Permission | Description | Dangerous |
|------------|-------------|-----------|
| `system.settings.read` | View server settings | No |
| `system.settings.write` | Modify server settings | Yes |
| `system.logs.read` | View activity logs | No |
| `system.jobs.read` | View background jobs | No |
| `system.jobs.manage` | Manage background jobs | Yes |
| `system.apikeys.manage` | Manage API keys | Yes |
| `system.roles.read` | View roles | No |
| `system.roles.manage` | Manage roles | Yes |

### User Permissions

| Permission | Description | Dangerous |
|------------|-------------|-----------|
| `users.read` | View user list | No |
| `users.create` | Create users | Yes |
| `users.update` | Update users | Yes |
| `users.delete` | Delete users | Yes |
| `users.sessions.manage` | Force logout | Yes |

### Library Permissions

| Permission | Description | Dangerous |
|------------|-------------|-----------|
| `libraries.read` | View libraries | No |
| `libraries.create` | Create libraries | Yes |
| `libraries.update` | Update libraries | Yes |
| `libraries.delete` | Delete libraries | Yes |
| `libraries.scan` | Trigger scans | No |

### Content Permissions

| Permission | Description | Dangerous |
|------------|-------------|-----------|
| `content.browse` | Browse content | No |
| `content.metadata.read` | View metadata | No |
| `content.metadata.write` | Edit metadata (audited) | No |
| `content.metadata.lock` | Lock metadata from auto-updates | No |
| `content.metadata.audit` | View metadata edit history | No |
| `content.images.manage` | Manage images | No |
| `content.delete` | Delete content | Yes |

### Playback Permissions

| Permission | Description | Dangerous |
|------------|-------------|-----------|
| `playback.stream` | Stream content | No |
| `playback.download` | Download files | No |
| `playback.transcode` | Request transcoding | No |

### Social Permissions

| Permission | Description | Dangerous |
|------------|-------------|-----------|
| `social.rate` | Rate content | No |
| `social.playlists.create` | Create playlists | No |
| `social.playlists.manage` | Manage playlists | No |
| `social.collections.create` | Create collections | No |
| `social.collections.manage` | Manage collections | No |
| `social.history.read` | View history | No |
| `social.favorites.manage` | Manage favorites | No |

### Request Permissions

| Permission | Description | Dangerous |
|------------|-------------|-----------|
| `requests.submit` | Submit content requests | No |
| `requests.view.own` | View own requests | No |
| `requests.vote` | Vote on requests | No |
| `requests.comment` | Comment on requests | No |
| `requests.cancel.own` | Cancel own pending requests | No |
| `requests.view.all` | View all users' requests | No |
| `requests.approve` | Approve requests | No |
| `requests.decline` | Decline requests | No |
| `requests.priority` | Set request priority | No |
| `requests.rules.read` | View request rules | No |
| `requests.rules.manage` | Create/edit request rules | Yes |
| `requests.quotas.read` | View user quotas | No |
| `requests.quotas.manage` | Manage user quotas | Yes |
| `requests.polls.vote` | Vote in polls | No |
| `requests.polls.create` | Create polls | No |
| `requests.polls.manage` | Manage polls | Yes |

### Adult Request Permissions (Schema `c`)

| Permission | Description | Dangerous |
|------------|-------------|-----------|
| `adult.requests.submit` | Submit adult content requests | No |
| `adult.requests.view.own` | View own adult requests | No |
| `adult.requests.vote` | Vote on adult requests | No |
| `adult.requests.approve` | Approve adult requests | No |
| `adult.requests.decline` | Decline adult requests | No |
| `adult.requests.rules.manage` | Manage adult request rules | Yes |
| `adult.requests.polls.create` | Create adult content polls | No |
| `adult.requests.polls.manage` | Manage adult polls | Yes |

### Adult Permissions

| Permission | Description | Dangerous |
|------------|-------------|-----------|
| `adult.browse` | Browse adult content | No |
| `adult.stream` | Stream adult content | No |
| `adult.metadata.write` | Edit adult metadata | No |

## Usage

### Checking Permissions in Handlers

```go
func (h *Handler) DeleteMovie(ctx context.Context, params api.DeleteMovieParams) error {
    user, ok := middleware.UserFromContext(ctx)
    if !ok {
        return ErrUnauthorized
    }

    // Check permission
    if err := h.rbac.RequirePermission(ctx, user.ID, rbac.PermContentDelete); err != nil {
        return ErrForbidden
    }

    // Proceed with deletion
    return h.service.DeleteMovie(ctx, params.ID)
}
```

### Using Middleware

```go
// Single permission
router.Handle("/api/v1/admin/settings",
    rbacMiddleware.RequirePermission(rbac.PermSystemSettingsWrite)(settingsHandler))

// Any of multiple permissions
router.Handle("/api/v1/libraries/{id}/scan",
    rbacMiddleware.RequireAnyPermission(rbac.PermLibrariesScan, rbac.PermSystemJobsManage)(scanHandler))

// Role-based (for backwards compatibility)
router.Handle("/api/v1/admin/users",
    rbacMiddleware.RequireRole("admin", "moderator")(usersHandler))
```

### Managing Roles via API

```go
// Create custom role
role, err := rbacService.CreateRole(ctx, rbac.CreateRoleParams{
    Name:        "editor",
    DisplayName: "Content Editor",
    Description: "Can edit content metadata and images",
    Color:       "#10B981",
    Permissions: []string{
        rbac.PermContentBrowse,
        rbac.PermContentMetadataRead,
        rbac.PermContentMetadataWrite,
        rbac.PermContentImagesManage,
    },
    CreatedBy: adminUserID,
})

// Update role permissions
err = rbacService.SetRolePermissions(ctx, "editor", []string{
    rbac.PermContentBrowse,
    rbac.PermContentMetadataRead,
    rbac.PermContentMetadataWrite,
    rbac.PermContentImagesManage,
    rbac.PermLibrariesScan, // Added new permission
})

// Delete custom role (only non-system roles)
err = rbacService.DeleteRole(ctx, "editor")
```

## Casbin Model

Located at `internal/service/rbac/model.conf`:

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

## Performance

- Policies are cached in memory by Casbin
- Database is only hit on policy changes or explicit reload
- Policy changes auto-save to database
- Thread-safe with read-write mutex

## Migration from Static RBAC

The migration `000018_dynamic_rbac.up.sql`:

1. Creates `roles` table with system roles (admin, moderator, user, guest)
2. Creates `permission_definitions` table
3. Adds `role_id` column to `users` table
4. Migrates existing users from ENUM to role_id
5. Casbin adapter creates `casbin_rules` table on first run
6. Default policies are seeded by `CasbinService` on startup

## Dependencies

```go
github.com/casbin/casbin/v2           // Casbin core
github.com/pckhoi/casbin-pgx-adapter/v3  // PostgreSQL adapter (pgx native)
```

## Metadata Auditing

All metadata edits by users with `content.metadata.write` permission are logged to the `activity_log` table for accountability and rollback capability.

### Audit Log Schema

```sql
-- activity_log table stores all audited actions
CREATE TABLE activity_log (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id),
    action          VARCHAR(50) NOT NULL,      -- 'metadata.edit', 'metadata.lock', 'image.upload'
    module          VARCHAR(50) NOT NULL,      -- 'movie', 'tvshow', 'music', etc.
    entity_id       UUID NOT NULL,             -- ID of the edited item
    entity_type     VARCHAR(50) NOT NULL,      -- 'movie', 'series', 'episode', 'album', etc.
    changes         JSONB NOT NULL,            -- {"field": {"old": "...", "new": "..."}}
    ip_address      INET,
    user_agent      TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_activity_log_entity ON activity_log(module, entity_type, entity_id);
CREATE INDEX idx_activity_log_user ON activity_log(user_id, created_at DESC);
CREATE INDEX idx_activity_log_action ON activity_log(action, created_at DESC);
```

### Audited Actions

| Action | Description |
|--------|-------------|
| `metadata.edit` | Field values changed (title, overview, year, etc.) |
| `metadata.lock` | Metadata locked from automatic provider updates |
| `metadata.unlock` | Metadata unlocked, allows auto-updates again |
| `metadata.refresh` | Manual metadata refresh triggered |
| `image.upload` | Custom image uploaded |
| `image.select` | Different image selected from provider |
| `image.delete` | Image removed |

### Metadata Locking

Editors can lock specific items from automatic metadata updates:

```go
// Lock metadata - prevents automatic updates
err := metadataService.LockMetadata(ctx, rbac.LockMetadataParams{
    EntityType: "movie",
    EntityID:   movieID,
    LockedBy:   userID,
    Reason:     "Manual corrections applied",
})

// Locked metadata is skipped during scheduled refreshes
func (s *MetadataService) RefreshMetadata(ctx context.Context, entityType string, entityID uuid.UUID) error {
    if s.IsLocked(ctx, entityType, entityID) {
        slog.Debug("skipping locked metadata", "type", entityType, "id", entityID)
        return nil
    }
    // ... fetch from provider
}
```

### Viewing Edit History

Users with `content.metadata.audit` permission can view edit history:

```go
// Get edit history for an item
history, err := activityService.GetEntityHistory(ctx, rbac.EntityHistoryParams{
    Module:     "movie",
    EntityType: "movie",
    EntityID:   movieID,
    Limit:      50,
})

// Returns list of changes with diffs
for _, entry := range history {
    fmt.Printf("%s by %s at %s\n", entry.Action, entry.UserName, entry.CreatedAt)
    for field, change := range entry.Changes {
        fmt.Printf("  %s: %q → %q\n", field, change.Old, change.New)
    }
}
```

### Rollback Support

Admins can revert to previous values using the audit log:

```go
// Revert a specific change
err := metadataService.RevertChange(ctx, rbac.RevertChangeParams{
    ActivityLogID: logEntryID,
    RevertedBy:    adminUserID,
    Reason:        "Incorrect edit by user",
})
```

## Related Documents

- [ARCHITECTURE_V2.md](../architecture/ARCHITECTURE_V2.md) - Overall architecture
- [DESIGN_PRINCIPLES.md](../architecture/DESIGN_PRINCIPLES.md) - Design principles
