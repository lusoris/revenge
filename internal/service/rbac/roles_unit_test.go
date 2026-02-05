package rbac_test

import (
	"context"
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/service/activity"
	"github.com/lusoris/revenge/internal/service/rbac"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func setupRolesTestService(t *testing.T) *rbac.Service {
	t.Helper()

	m, err := model.NewModelFromString(casbinModelConf)
	require.NoError(t, err)

	enforcer, err := casbin.NewEnforcer(m)
	require.NoError(t, err)

	logger := zap.NewNop()
	return rbac.NewService(enforcer, logger, activity.NewNoopLogger())
}

func TestPermission_String(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	perm := rbac.Permission{Resource: "movies", Action: "read"}
	assert.Equal(t, "movies:read", perm.String())
}

func TestParsePermission(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("valid permission", func(t *testing.T) {
		perm, err := rbac.ParsePermission("movies:read")
		require.NoError(t, err)
		assert.Equal(t, "movies", perm.Resource)
		assert.Equal(t, "read", perm.Action)
	})

	t.Run("invalid format", func(t *testing.T) {
		_, err := rbac.ParsePermission("invalid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid permission format")
	})

	t.Run("complex resource:action", func(t *testing.T) {
		perm, err := rbac.ParsePermission("admin:settings:write")
		require.NoError(t, err)
		assert.Equal(t, "admin", perm.Resource)
		assert.Equal(t, "settings:write", perm.Action)
	})
}

func TestService_ListRoles_Unit(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("empty when no roles", func(t *testing.T) {
		svc := setupRolesTestService(t)
		ctx := context.Background()

		roles, err := svc.ListRoles(ctx)
		require.NoError(t, err)
		assert.Empty(t, roles)
	})

	t.Run("returns all roles", func(t *testing.T) {
		svc := setupRolesTestService(t)
		ctx := context.Background()

		// Create roles
		_, err := svc.CreateRole(ctx, "editor", "Edit content", []rbac.Permission{
			{Resource: "movies", Action: "write"},
		})
		require.NoError(t, err)

		_, err = svc.CreateRole(ctx, "viewer", "View content", []rbac.Permission{
			{Resource: "movies", Action: "read"},
		})
		require.NoError(t, err)

		roles, err := svc.ListRoles(ctx)
		require.NoError(t, err)
		assert.Len(t, roles, 2)

		// Should be sorted by name
		assert.Equal(t, "editor", roles[0].Name)
		assert.Equal(t, "viewer", roles[1].Name)
	})
}

func TestService_GetRole_Unit(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		svc := setupRolesTestService(t)
		ctx := context.Background()

		_, err := svc.CreateRole(ctx, "testrole", "Test description", []rbac.Permission{
			{Resource: "movies", Action: "read"},
			{Resource: "movies", Action: "write"},
		})
		require.NoError(t, err)

		role, err := svc.GetRole(ctx, "testrole")
		require.NoError(t, err)
		assert.Equal(t, "testrole", role.Name)
		assert.Len(t, role.Permissions, 2)
	})

	t.Run("not found", func(t *testing.T) {
		svc := setupRolesTestService(t)
		ctx := context.Background()

		_, err := svc.GetRole(ctx, "nonexistent")
		assert.ErrorIs(t, err, rbac.ErrRoleNotFound)
	})
}

func TestService_CreateRole_Unit(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		svc := setupRolesTestService(t)
		ctx := context.Background()

		role, err := svc.CreateRole(ctx, "newrole", "A new role", []rbac.Permission{
			{Resource: "movies", Action: "read"},
			{Resource: "shows", Action: "read"},
		})
		require.NoError(t, err)
		assert.Equal(t, "newrole", role.Name)
		assert.Equal(t, "A new role", role.Description)
		assert.Len(t, role.Permissions, 2)
		assert.False(t, role.IsBuiltIn)
		assert.Equal(t, 0, role.UserCount)
	})

	t.Run("already exists", func(t *testing.T) {
		svc := setupRolesTestService(t)
		ctx := context.Background()

		_, err := svc.CreateRole(ctx, "existing", "", []rbac.Permission{
			{Resource: "movies", Action: "read"},
		})
		require.NoError(t, err)

		_, err = svc.CreateRole(ctx, "existing", "", []rbac.Permission{
			{Resource: "shows", Action: "read"},
		})
		assert.ErrorIs(t, err, rbac.ErrRoleAlreadyExists)
	})

	t.Run("empty name", func(t *testing.T) {
		svc := setupRolesTestService(t)
		ctx := context.Background()

		_, err := svc.CreateRole(ctx, "", "", []rbac.Permission{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be empty")
	})

	t.Run("name with whitespace", func(t *testing.T) {
		svc := setupRolesTestService(t)
		ctx := context.Background()

		_, err := svc.CreateRole(ctx, "invalid name", "", []rbac.Permission{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "whitespace")
	})
}

func TestService_DeleteRole_Unit(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		svc := setupRolesTestService(t)
		ctx := context.Background()

		_, err := svc.CreateRole(ctx, "todelete", "", []rbac.Permission{
			{Resource: "movies", Action: "read"},
		})
		require.NoError(t, err)

		err = svc.DeleteRole(ctx, "todelete")
		require.NoError(t, err)

		_, err = svc.GetRole(ctx, "todelete")
		assert.ErrorIs(t, err, rbac.ErrRoleNotFound)
	})

	t.Run("built-in role", func(t *testing.T) {
		svc := setupRolesTestService(t)
		ctx := context.Background()

		err := svc.DeleteRole(ctx, "admin")
		assert.ErrorIs(t, err, rbac.ErrBuiltInRole)
	})

	t.Run("not found", func(t *testing.T) {
		svc := setupRolesTestService(t)
		ctx := context.Background()

		err := svc.DeleteRole(ctx, "nonexistent")
		assert.ErrorIs(t, err, rbac.ErrRoleNotFound)
	})

	t.Run("role in use", func(t *testing.T) {
		svc := setupRolesTestService(t)
		ctx := context.Background()

		_, err := svc.CreateRole(ctx, "inuse", "", []rbac.Permission{
			{Resource: "movies", Action: "read"},
		})
		require.NoError(t, err)

		userID := uuid.New()
		require.NoError(t, svc.AssignRole(ctx, userID, "inuse"))

		err = svc.DeleteRole(ctx, "inuse")
		assert.ErrorIs(t, err, rbac.ErrRoleInUse)
	})
}

func TestService_UpdateRolePermissions_Unit(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		svc := setupRolesTestService(t)
		ctx := context.Background()

		_, err := svc.CreateRole(ctx, "updatable", "", []rbac.Permission{
			{Resource: "movies", Action: "read"},
		})
		require.NoError(t, err)

		newPerms := []rbac.Permission{
			{Resource: "movies", Action: "read"},
			{Resource: "movies", Action: "write"},
			{Resource: "shows", Action: "read"},
		}

		role, err := svc.UpdateRolePermissions(ctx, "updatable", newPerms)
		require.NoError(t, err)
		assert.Len(t, role.Permissions, 3)
	})

	t.Run("not found", func(t *testing.T) {
		svc := setupRolesTestService(t)
		ctx := context.Background()

		_, err := svc.UpdateRolePermissions(ctx, "nonexistent", []rbac.Permission{})
		assert.ErrorIs(t, err, rbac.ErrRoleNotFound)
	})
}

func TestService_ListPermissions_Unit(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	svc := setupRolesTestService(t)
	ctx := context.Background()

	perms := svc.ListPermissions(ctx)
	assert.NotEmpty(t, perms)

	// Should contain combinations of resources and actions
	hasMoviesRead := false
	for _, p := range perms {
		if p.Resource == "movies" && p.Action == "read" {
			hasMoviesRead = true
			break
		}
	}
	assert.True(t, hasMoviesRead, "should contain movies:read permission")
}

func TestService_GetRolePermissions_Unit(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("returns permissions", func(t *testing.T) {
		svc := setupRolesTestService(t)
		ctx := context.Background()

		_, err := svc.CreateRole(ctx, "witherms", "", []rbac.Permission{
			{Resource: "movies", Action: "read"},
			{Resource: "movies", Action: "write"},
		})
		require.NoError(t, err)

		perms, err := svc.GetRolePermissions(ctx, "witherms")
		require.NoError(t, err)
		assert.Len(t, perms, 2)
	})

	t.Run("empty for nonexistent", func(t *testing.T) {
		svc := setupRolesTestService(t)
		ctx := context.Background()

		perms, err := svc.GetRolePermissions(ctx, "nonexistent")
		require.NoError(t, err)
		assert.Empty(t, perms)
	})
}

func TestService_AddPermissionToRole_Unit(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		svc := setupRolesTestService(t)
		ctx := context.Background()

		_, err := svc.CreateRole(ctx, "addperm", "", []rbac.Permission{
			{Resource: "movies", Action: "read"},
		})
		require.NoError(t, err)

		err = svc.AddPermissionToRole(ctx, "addperm", rbac.Permission{
			Resource: "movies",
			Action:   "write",
		})
		require.NoError(t, err)

		perms, err := svc.GetRolePermissions(ctx, "addperm")
		require.NoError(t, err)
		assert.Len(t, perms, 2)
	})

	t.Run("duplicate", func(t *testing.T) {
		svc := setupRolesTestService(t)
		ctx := context.Background()

		_, err := svc.CreateRole(ctx, "dupperm", "", []rbac.Permission{
			{Resource: "movies", Action: "read"},
		})
		require.NoError(t, err)

		err = svc.AddPermissionToRole(ctx, "dupperm", rbac.Permission{
			Resource: "movies",
			Action:   "read",
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestService_RemovePermissionFromRole_Unit(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		svc := setupRolesTestService(t)
		ctx := context.Background()

		_, err := svc.CreateRole(ctx, "remperm", "", []rbac.Permission{
			{Resource: "movies", Action: "read"},
			{Resource: "movies", Action: "write"},
		})
		require.NoError(t, err)

		err = svc.RemovePermissionFromRole(ctx, "remperm", rbac.Permission{
			Resource: "movies",
			Action:   "write",
		})
		require.NoError(t, err)

		perms, err := svc.GetRolePermissions(ctx, "remperm")
		require.NoError(t, err)
		assert.Len(t, perms, 1)
	})

	t.Run("not found", func(t *testing.T) {
		svc := setupRolesTestService(t)
		ctx := context.Background()

		_, err := svc.CreateRole(ctx, "noperm", "", []rbac.Permission{
			{Resource: "movies", Action: "read"},
		})
		require.NoError(t, err)

		err = svc.RemovePermissionFromRole(ctx, "noperm", rbac.Permission{
			Resource: "movies",
			Action:   "write",
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestService_GetAllRoleNames_Unit(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("empty", func(t *testing.T) {
		svc := setupRolesTestService(t)
		ctx := context.Background()

		names, err := svc.GetAllRoleNames(ctx)
		require.NoError(t, err)
		assert.Empty(t, names)
	})

	t.Run("returns names sorted", func(t *testing.T) {
		svc := setupRolesTestService(t)
		ctx := context.Background()

		_, err := svc.CreateRole(ctx, "zeta", "", []rbac.Permission{{Resource: "a", Action: "a"}})
		require.NoError(t, err)
		_, err = svc.CreateRole(ctx, "alpha", "", []rbac.Permission{{Resource: "b", Action: "b"}})
		require.NoError(t, err)
		_, err = svc.CreateRole(ctx, "beta", "", []rbac.Permission{{Resource: "c", Action: "c"}})
		require.NoError(t, err)

		names, err := svc.GetAllRoleNames(ctx)
		require.NoError(t, err)
		assert.Equal(t, []string{"alpha", "beta", "zeta"}, names)
	})
}

func TestService_CheckUserPermission_Unit(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	svc := setupRolesTestService(t)
	ctx := context.Background()

	userID := uuid.New()

	// Create role with permission
	_, err := svc.CreateRole(ctx, "checker", "", []rbac.Permission{
		{Resource: "movies", Action: "read"},
	})
	require.NoError(t, err)

	// User without role
	allowed, err := svc.CheckUserPermission(ctx, userID, "movies", "read")
	require.NoError(t, err)
	assert.False(t, allowed)

	// Assign role
	require.NoError(t, svc.AssignRole(ctx, userID, "checker"))

	// User with role
	allowed, err = svc.CheckUserPermission(ctx, userID, "movies", "read")
	require.NoError(t, err)
	assert.True(t, allowed)

	// User without specific permission
	allowed, err = svc.CheckUserPermission(ctx, userID, "movies", "write")
	require.NoError(t, err)
	assert.False(t, allowed)
}

func TestBuiltInRoles(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	// Check built-in roles are defined
	assert.Contains(t, rbac.BuiltInRoles, "admin")
	assert.Contains(t, rbac.BuiltInRoles, "moderator")
	assert.Contains(t, rbac.BuiltInRoles, "user")
	assert.Contains(t, rbac.BuiltInRoles, "guest")
}

func TestAvailableResourcesAndActions(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	// Check resources are defined
	assert.NotEmpty(t, rbac.AvailableResources)
	assert.Contains(t, rbac.AvailableResources, "movies")
	assert.Contains(t, rbac.AvailableResources, "users")

	// Check actions are defined
	assert.NotEmpty(t, rbac.AvailableActions)
	assert.Contains(t, rbac.AvailableActions, "read")
	assert.Contains(t, rbac.AvailableActions, "write")
}
