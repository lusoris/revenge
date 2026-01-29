# RBAC Development Instructions

> Instructions for implementing role-based access control in Revenge.

## Quick Reference

### Permission Constants

Always use the constants from `internal/service/rbac/casbin.go`:

```go
import "github.com/lusoris/revenge/internal/service/rbac"

// System
rbac.PermSystemSettingsRead
rbac.PermSystemSettingsWrite
rbac.PermSystemLogsRead
rbac.PermSystemJobsRead
rbac.PermSystemJobsManage
rbac.PermSystemAPIKeysManage
rbac.PermSystemRolesRead
rbac.PermSystemRolesManage

// Users
rbac.PermUsersRead
rbac.PermUsersCreate
rbac.PermUsersUpdate
rbac.PermUsersDelete
rbac.PermUsersSessionsManage

// Libraries
rbac.PermLibrariesRead
rbac.PermLibrariesCreate
rbac.PermLibrariesUpdate
rbac.PermLibrariesDelete
rbac.PermLibrariesScan

// Content
rbac.PermContentBrowse
rbac.PermContentMetadataRead
rbac.PermContentMetadataWrite
rbac.PermContentImagesManage
rbac.PermContentDelete

// Playback
rbac.PermPlaybackStream
rbac.PermPlaybackDownload
rbac.PermPlaybackTranscode

// Social
rbac.PermSocialRate
rbac.PermSocialPlaylistsCreate
rbac.PermSocialPlaylistsManage
rbac.PermSocialCollectionsCreate
rbac.PermSocialCollectionsManage
rbac.PermSocialHistoryRead
rbac.PermSocialFavoritesManage

// Adult
rbac.PermAdultBrowse
rbac.PermAdultStream
rbac.PermAdultMetadataWrite
```

## Checking Permissions

### In Handlers (Recommended)

```go
func (h *Handler) SomeAction(ctx context.Context, params api.Params) error {
    user, ok := middleware.UserFromContext(ctx)
    if !ok {
        return api.ErrUnauthorized
    }

    // Single permission check
    if err := h.rbac.RequirePermission(ctx, user.ID, rbac.PermContentDelete); err != nil {
        return api.ErrForbidden
    }

    // Or check any of multiple permissions
    if err := h.rbac.RequireAnyPermission(ctx, user.ID, []string{
        rbac.PermLibrariesScan,
        rbac.PermSystemJobsManage,
    }); err != nil {
        return api.ErrForbidden
    }

    return h.service.DoAction(ctx, params)
}
```

### Boolean Check (When You Need to Branch)

```go
canDelete, err := h.rbac.HasPermission(ctx, user.ID, rbac.PermContentDelete)
if err != nil {
    return err
}

if canDelete {
    // Show delete button in response
}
```

### Via Middleware

```go
// Protect entire route
mux.Handle("/api/v1/admin/settings",
    rbacMiddleware.RequirePermission(rbac.PermSystemSettingsWrite)(handler))

// Multiple permissions (any)
mux.Handle("/api/v1/scan",
    rbacMiddleware.RequireAnyPermission(
        rbac.PermLibrariesScan,
        rbac.PermSystemJobsManage,
    )(handler))
```

## Adding New Permissions

1. **Add constant** in `internal/service/rbac/casbin.go`:

```go
const (
    // ... existing permissions ...

    // New permission
    PermMyNewPermission = "mymodule.action"
)
```

2. **Add to AllPermissions()** function:

```go
func AllPermissions() []string {
    return []string{
        // ... existing ...
        PermMyNewPermission,
    }
}
```

3. **Add to migration** `000018_dynamic_rbac.up.sql`:

```sql
INSERT INTO permission_definitions (name, display_name, description, category, is_dangerous) VALUES
    ('mymodule.action', 'My Action', 'Description of the permission', 'Category', false);
```

4. **Assign to default roles** in `seedDefaultPolicies()`:

```go
// In moderatorPerms, userPerms, or guestPerms as appropriate
moderatorPerms := []string{
    // ... existing ...
    PermMyNewPermission,
}
```

## Role Management API

### List Roles

```go
roles, err := rbacService.ListRoles(ctx)
// Returns []rbac.Role with permissions included
```

### Create Role

```go
role, err := rbacService.CreateRole(ctx, rbac.CreateRoleParams{
    Name:        "editor",           // Unique identifier
    DisplayName: "Content Editor",   // UI display name
    Description: "Can edit content", // Optional description
    Color:       "#10B981",          // Optional hex color
    Icon:        "edit",             // Optional icon name
    Priority:    200,                // For sorting (higher = more important)
    Permissions: []string{           // Initial permissions
        rbac.PermContentMetadataWrite,
    },
    CreatedBy: adminUserID,
})
```

### Update Role Permissions

```go
// Replace all permissions
err := rbacService.SetRolePermissions(ctx, "editor", []string{
    rbac.PermContentBrowse,
    rbac.PermContentMetadataRead,
    rbac.PermContentMetadataWrite,
})

// Add single permission
err := rbacService.AddRolePermission(ctx, "editor", rbac.PermContentImagesManage)

// Remove single permission
err := rbacService.RemoveRolePermission(ctx, "editor", rbac.PermContentDelete)
```

### Delete Role

```go
// Only non-system roles can be deleted
// Returns ErrRoleInUse if users are assigned to it
err := rbacService.DeleteRole(ctx, "editor")
```

## System Roles

These roles cannot be deleted:

| Role | Purpose |
|------|---------|
| `admin` | Full access |
| `moderator` | Content management |
| `user` | Standard user |
| `guest` | Browse only |

You CAN modify their permissions, but not delete them.

## Adult Content Access

Adult content requires BOTH:
1. `adult_enabled = true` on the user record
2. `adult.browse` permission

```go
canAccess, err := rbacService.CanAccessAdultContent(ctx, user)
```

## Testing

```go
func TestHandler_DeleteMovie(t *testing.T) {
    // Setup mock RBAC service
    mockRBAC := &MockRBACService{
        permissions: map[string]bool{
            rbac.PermContentDelete: true,
        },
    }

    handler := NewHandler(service, mockRBAC)

    // Test with permission
    err := handler.DeleteMovie(ctx, params)
    assert.NoError(t, err)

    // Test without permission
    mockRBAC.permissions[rbac.PermContentDelete] = false
    err = handler.DeleteMovie(ctx, params)
    assert.ErrorIs(t, err, ErrForbidden)
}
```

## Common Patterns

### Optional Permission (Show/Hide UI Elements)

```go
response := &MovieResponse{
    Movie: movie,
}

// Check if user can delete (for UI button visibility)
if canDelete, _ := h.rbac.HasPermission(ctx, user.ID, rbac.PermContentDelete); canDelete {
    response.Actions = append(response.Actions, "delete")
}
```

### Resource Owner Check

```go
// User can always edit their own resources
if playlist.OwnerID == user.ID {
    // Allow
} else {
    // Check permission for editing others' playlists
    if err := h.rbac.RequirePermission(ctx, user.ID, rbac.PermSocialPlaylistsManage); err != nil {
        return ErrForbidden
    }
}
```

## DO NOT

- Hard-code role names in permission checks (use permissions instead)
- Skip permission checks for "admin" users (admins get permissions via Casbin)
- Create permissions without adding to `AllPermissions()`
- Delete system roles
- Assume permission check results are cached (they are, but treat as uncached)

## Related

- [RBAC_CASBIN.md](../../docs/dev/design/features/RBAC_CASBIN.md) - Full documentation
- [ARCHITECTURE_V2.md](../../docs/dev/design/architecture/ARCHITECTURE_V2.md) - System architecture
