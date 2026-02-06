package rbac

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPermission_String(t *testing.T) {
	t.Parallel()
	perm := Permission{Resource: "users", Action: "read"}
	assert.Equal(t, "users:read", perm.String())
}

func TestParsePermission(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    Permission
		wantErr bool
	}{
		{
			name:  "valid permission",
			input: "users:read",
			want:  Permission{Resource: "users", Action: "read"},
		},
		{
			name:  "valid permission with write",
			input: "movies:write",
			want:  Permission{Resource: "movies", Action: "write"},
		},
		{
			name:    "invalid format - no colon",
			input:   "usersread",
			wantErr: true,
		},
		{
			name:    "invalid format - empty",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePermission(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_CreateRole(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	perms := []Permission{
		{Resource: "users", Action: "read"},
		{Resource: "users", Action: "write"},
	}

	role, err := svc.CreateRole(ctx, "editor", "Can edit content", perms)
	require.NoError(t, err)
	assert.Equal(t, "editor", role.Name)
	assert.Equal(t, "Can edit content", role.Description)
	assert.Equal(t, 2, len(role.Permissions))
	assert.False(t, role.IsBuiltIn)
	assert.Equal(t, 0, role.UserCount)
}

func TestService_CreateRole_EmptyName(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	_, err := svc.CreateRole(ctx, "", "description", []Permission{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")
}

func TestService_CreateRole_WithWhitespace(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	_, err := svc.CreateRole(ctx, "role name", "description", []Permission{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot contain whitespace")
}

func TestService_CreateRole_AlreadyExists(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	perms := []Permission{{Resource: "users", Action: "read"}}
	_, err := svc.CreateRole(ctx, "editor", "description", perms)
	require.NoError(t, err)

	_, err = svc.CreateRole(ctx, "editor", "description", perms)
	assert.ErrorIs(t, err, ErrRoleAlreadyExists)
}

func TestService_GetRole(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Create a role first
	perms := []Permission{
		{Resource: "movies", Action: "read"},
		{Resource: "movies", Action: "write"},
	}
	created, err := svc.CreateRole(ctx, "movieeditor", "Can edit movies", perms)
	require.NoError(t, err)

	// Get the role
	role, err := svc.GetRole(ctx, "movieeditor")
	require.NoError(t, err)
	assert.Equal(t, created.Name, role.Name)
	assert.Equal(t, 2, len(role.Permissions))
}

func TestService_GetRole_NotFound(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	_, err := svc.GetRole(ctx, "nonexistent")
	assert.ErrorIs(t, err, ErrRoleNotFound)
}

func TestService_ListRoles(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Create several roles
	_, err := svc.CreateRole(ctx, "role1", "First role", []Permission{
		{Resource: "users", Action: "read"},
	})
	require.NoError(t, err)

	_, err = svc.CreateRole(ctx, "role2", "Second role", []Permission{
		{Resource: "movies", Action: "read"},
	})
	require.NoError(t, err)

	// List all roles
	roles, err := svc.ListRoles(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(roles), 2)

	// Verify roles are sorted
	for i := 1; i < len(roles); i++ {
		assert.LessOrEqual(t, roles[i-1].Name, roles[i].Name)
	}
}

func TestService_DeleteRole(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Create a role
	_, err := svc.CreateRole(ctx, "deleteme", "Will be deleted", []Permission{
		{Resource: "users", Action: "read"},
	})
	require.NoError(t, err)

	// Delete the role
	err = svc.DeleteRole(ctx, "deleteme")
	require.NoError(t, err)

	// Verify it's gone
	_, err = svc.GetRole(ctx, "deleteme")
	assert.ErrorIs(t, err, ErrRoleNotFound)
}

func TestService_DeleteRole_BuiltIn(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Try to delete a built-in role
	err := svc.DeleteRole(ctx, "admin")
	assert.ErrorIs(t, err, ErrBuiltInRole)
}

func TestService_DeleteRole_NotFound(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	err := svc.DeleteRole(ctx, "nonexistent")
	assert.ErrorIs(t, err, ErrRoleNotFound)
}

func TestService_DeleteRole_InUse(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Create a role and assign it to a user
	_, err := svc.CreateRole(ctx, "inuse", "Role in use", []Permission{
		{Resource: "users", Action: "read"},
	})
	require.NoError(t, err)

	userID := uuid.Must(uuid.NewV7())
	err = svc.AssignRole(ctx, userID, "inuse")
	require.NoError(t, err)

	// Try to delete the role
	err = svc.DeleteRole(ctx, "inuse")
	assert.ErrorIs(t, err, ErrRoleInUse)
}

func TestService_UpdateRolePermissions(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Create a role
	_, err := svc.CreateRole(ctx, "updatable", "Will be updated", []Permission{
		{Resource: "users", Action: "read"},
	})
	require.NoError(t, err)

	// Update permissions
	newPerms := []Permission{
		{Resource: "movies", Action: "read"},
		{Resource: "movies", Action: "write"},
	}
	updated, err := svc.UpdateRolePermissions(ctx, "updatable", newPerms)
	require.NoError(t, err)
	assert.Equal(t, 2, len(updated.Permissions))

	// Verify permissions were updated
	role, err := svc.GetRole(ctx, "updatable")
	require.NoError(t, err)
	assert.Equal(t, 2, len(role.Permissions))
}

func TestService_UpdateRolePermissions_NotFound(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	_, err := svc.UpdateRolePermissions(ctx, "nonexistent", []Permission{})
	assert.ErrorIs(t, err, ErrRoleNotFound)
}

func TestService_ListPermissions(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	perms := svc.ListPermissions(ctx)
	assert.Greater(t, len(perms), 0)

	// Verify all combinations are present
	expectedCount := len(AvailableResources) * len(AvailableActions)
	assert.Equal(t, expectedCount, len(perms))
}

func TestService_GetRolePermissions(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Create a role with specific permissions
	perms := []Permission{
		{Resource: "users", Action: "read"},
		{Resource: "users", Action: "write"},
	}
	_, err := svc.CreateRole(ctx, "testrole", "Test role", perms)
	require.NoError(t, err)

	// Get permissions
	rolePerms, err := svc.GetRolePermissions(ctx, "testrole")
	require.NoError(t, err)
	assert.Equal(t, 2, len(rolePerms))
}

func TestService_GetRolePermissions_Empty(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Get permissions for non-existent role
	rolePerms, err := svc.GetRolePermissions(ctx, "nonexistent")
	require.NoError(t, err)
	assert.Equal(t, 0, len(rolePerms))
}

func TestService_AddPermissionToRole(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Create a role
	_, err := svc.CreateRole(ctx, "addperm", "Add permission test", []Permission{
		{Resource: "users", Action: "read"},
	})
	require.NoError(t, err)

	// Add a permission
	err = svc.AddPermissionToRole(ctx, "addperm", Permission{Resource: "users", Action: "write"})
	require.NoError(t, err)

	// Verify permission was added
	perms, err := svc.GetRolePermissions(ctx, "addperm")
	require.NoError(t, err)
	assert.Equal(t, 2, len(perms))
}

func TestService_AddPermissionToRole_Duplicate(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Create a role with a permission
	_, err := svc.CreateRole(ctx, "dupeperm", "Duplicate permission test", []Permission{
		{Resource: "users", Action: "read"},
	})
	require.NoError(t, err)

	// Try to add the same permission again
	err = svc.AddPermissionToRole(ctx, "dupeperm", Permission{Resource: "users", Action: "read"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestService_RemovePermissionFromRole(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Create a role with multiple permissions
	_, err := svc.CreateRole(ctx, "remperm", "Remove permission test", []Permission{
		{Resource: "users", Action: "read"},
		{Resource: "users", Action: "write"},
	})
	require.NoError(t, err)

	// Remove a permission
	err = svc.RemovePermissionFromRole(ctx, "remperm", Permission{Resource: "users", Action: "write"})
	require.NoError(t, err)

	// Verify permission was removed
	perms, err := svc.GetRolePermissions(ctx, "remperm")
	require.NoError(t, err)
	assert.Equal(t, 1, len(perms))
}

func TestService_RemovePermissionFromRole_NotFound(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Create a role
	_, err := svc.CreateRole(ctx, "noperm", "No permission test", []Permission{
		{Resource: "users", Action: "read"},
	})
	require.NoError(t, err)

	// Try to remove non-existent permission
	err = svc.RemovePermissionFromRole(ctx, "noperm", Permission{Resource: "users", Action: "delete"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestService_GetAllRoleNames(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Create several roles
	_, err := svc.CreateRole(ctx, "role1", "", []Permission{{Resource: "users", Action: "read"}})
	require.NoError(t, err)
	_, err = svc.CreateRole(ctx, "role2", "", []Permission{{Resource: "movies", Action: "read"}})
	require.NoError(t, err)

	// Get all role names
	names, err := svc.GetAllRoleNames(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(names), 2)

	// Verify sorted
	for i := 1; i < len(names); i++ {
		assert.LessOrEqual(t, names[i-1], names[i])
	}
}

func TestService_CheckUserPermission(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Create a role and assign to user
	userID := uuid.Must(uuid.NewV7())
	_, err := svc.CreateRole(ctx, "checkperm", "", []Permission{
		{Resource: "movies", Action: "read"},
	})
	require.NoError(t, err)

	err = svc.AssignRole(ctx, userID, "checkperm")
	require.NoError(t, err)

	// Check user has permission
	allowed, err := svc.CheckUserPermission(ctx, userID, "movies", "read")
	require.NoError(t, err)
	assert.True(t, allowed)

	// Check user doesn't have other permission
	allowed, err = svc.CheckUserPermission(ctx, userID, "movies", "write")
	require.NoError(t, err)
	assert.False(t, allowed)
}
