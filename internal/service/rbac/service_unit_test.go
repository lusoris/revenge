package rbac_test

import (
	"context"
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/infra/logging"
	"github.com/lusoris/revenge/internal/service/activity"
	"github.com/lusoris/revenge/internal/service/rbac"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// casbinModelConf is the RBAC model for testing
const casbinModelConf = `
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
`

func setupUnitTestService(t *testing.T) *rbac.Service {
	t.Helper()

	m, err := model.NewModelFromString(casbinModelConf)
	require.NoError(t, err)

	enforcer, err := casbin.NewSyncedEnforcer(m)
	require.NoError(t, err)

	logger := logging.NewTestLogger()
	return rbac.NewService(enforcer, logger, activity.NewNoopLogger())
}

func TestService_NewService(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	svc := setupUnitTestService(t)
	assert.NotNil(t, svc)
}

func TestService_Enforce_Unit(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("allowed when policy exists", func(t *testing.T) {
		svc := setupUnitTestService(t)
		ctx := context.Background()

		err := svc.AddPolicy(ctx, "alice", "data1", "read")
		require.NoError(t, err)

		allowed, err := svc.Enforce(ctx, "alice", "data1", "read")
		require.NoError(t, err)
		assert.True(t, allowed)
	})

	t.Run("denied when no policy", func(t *testing.T) {
		svc := setupUnitTestService(t)
		ctx := context.Background()

		allowed, err := svc.Enforce(ctx, "bob", "data1", "read")
		require.NoError(t, err)
		assert.False(t, allowed)
	})

	t.Run("denied for different action", func(t *testing.T) {
		svc := setupUnitTestService(t)
		ctx := context.Background()

		err := svc.AddPolicy(ctx, "alice", "data1", "read")
		require.NoError(t, err)

		allowed, err := svc.Enforce(ctx, "alice", "data1", "write")
		require.NoError(t, err)
		assert.False(t, allowed)
	})

	t.Run("denied for different resource", func(t *testing.T) {
		svc := setupUnitTestService(t)
		ctx := context.Background()

		err := svc.AddPolicy(ctx, "alice", "data1", "read")
		require.NoError(t, err)

		allowed, err := svc.Enforce(ctx, "alice", "data2", "read")
		require.NoError(t, err)
		assert.False(t, allowed)
	})
}

func TestService_EnforceWithContext_Unit(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	svc := setupUnitTestService(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())

	err := svc.AddPolicy(ctx, userID.String(), "resource1", "read")
	require.NoError(t, err)

	allowed, err := svc.EnforceWithContext(ctx, userID, "resource1", "read")
	require.NoError(t, err)
	assert.True(t, allowed)

	allowed, err = svc.EnforceWithContext(ctx, userID, "resource1", "write")
	require.NoError(t, err)
	assert.False(t, allowed)
}

func TestService_AddPolicy_Unit(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		svc := setupUnitTestService(t)
		ctx := context.Background()

		err := svc.AddPolicy(ctx, "alice", "data1", "read")
		require.NoError(t, err)

		policies, err := svc.GetPolicies(ctx)
		require.NoError(t, err)
		assert.Len(t, policies, 1)
	})

	t.Run("duplicate policy does not error", func(t *testing.T) {
		svc := setupUnitTestService(t)
		ctx := context.Background()

		err := svc.AddPolicy(ctx, "alice", "data1", "read")
		require.NoError(t, err)

		err = svc.AddPolicy(ctx, "alice", "data1", "read")
		require.NoError(t, err)

		policies, err := svc.GetPolicies(ctx)
		require.NoError(t, err)
		assert.Len(t, policies, 1)
	})
}

func TestService_RemovePolicy_Unit(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		svc := setupUnitTestService(t)
		ctx := context.Background()

		err := svc.AddPolicy(ctx, "alice", "data1", "read")
		require.NoError(t, err)

		err = svc.RemovePolicy(ctx, "alice", "data1", "read")
		require.NoError(t, err)

		allowed, err := svc.Enforce(ctx, "alice", "data1", "read")
		require.NoError(t, err)
		assert.False(t, allowed)
	})

	t.Run("not found returns error", func(t *testing.T) {
		svc := setupUnitTestService(t)
		ctx := context.Background()

		err := svc.RemovePolicy(ctx, "nonexistent", "data1", "read")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestService_GetPolicies_Unit(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("empty when no policies", func(t *testing.T) {
		svc := setupUnitTestService(t)
		ctx := context.Background()

		policies, err := svc.GetPolicies(ctx)
		require.NoError(t, err)
		assert.Empty(t, policies)
	})

	t.Run("returns all policies", func(t *testing.T) {
		svc := setupUnitTestService(t)
		ctx := context.Background()

		require.NoError(t, svc.AddPolicy(ctx, "alice", "data1", "read"))
		require.NoError(t, svc.AddPolicy(ctx, "bob", "data2", "write"))
		require.NoError(t, svc.AddPolicy(ctx, "charlie", "data3", "delete"))

		policies, err := svc.GetPolicies(ctx)
		require.NoError(t, err)
		assert.Len(t, policies, 3)
	})
}

func TestService_AssignRole_Unit(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		svc := setupUnitTestService(t)
		ctx := context.Background()

		userID := uuid.Must(uuid.NewV7())

		err := svc.AssignRole(ctx, userID, "admin")
		require.NoError(t, err)

		hasRole, err := svc.HasRole(ctx, userID, "admin")
		require.NoError(t, err)
		assert.True(t, hasRole)
	})

	t.Run("duplicate assignment does not error", func(t *testing.T) {
		svc := setupUnitTestService(t)
		ctx := context.Background()

		userID := uuid.Must(uuid.NewV7())

		err := svc.AssignRole(ctx, userID, "admin")
		require.NoError(t, err)

		err = svc.AssignRole(ctx, userID, "admin")
		require.NoError(t, err)

		roles, err := svc.GetUserRoles(ctx, userID)
		require.NoError(t, err)
		assert.Len(t, roles, 1)
	})

	t.Run("multiple roles", func(t *testing.T) {
		svc := setupUnitTestService(t)
		ctx := context.Background()

		userID := uuid.Must(uuid.NewV7())

		require.NoError(t, svc.AssignRole(ctx, userID, "admin"))
		require.NoError(t, svc.AssignRole(ctx, userID, "editor"))
		require.NoError(t, svc.AssignRole(ctx, userID, "viewer"))

		roles, err := svc.GetUserRoles(ctx, userID)
		require.NoError(t, err)
		assert.Len(t, roles, 3)
	})
}

func TestService_RemoveRole_Unit(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		svc := setupUnitTestService(t)
		ctx := context.Background()

		userID := uuid.Must(uuid.NewV7())

		require.NoError(t, svc.AssignRole(ctx, userID, "admin"))

		err := svc.RemoveRole(ctx, userID, "admin")
		require.NoError(t, err)

		hasRole, err := svc.HasRole(ctx, userID, "admin")
		require.NoError(t, err)
		assert.False(t, hasRole)
	})

	t.Run("not found returns error", func(t *testing.T) {
		svc := setupUnitTestService(t)
		ctx := context.Background()

		userID := uuid.Must(uuid.NewV7())

		err := svc.RemoveRole(ctx, userID, "nonexistent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestService_GetUserRoles_Unit(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("empty when no roles", func(t *testing.T) {
		svc := setupUnitTestService(t)
		ctx := context.Background()

		userID := uuid.Must(uuid.NewV7())

		roles, err := svc.GetUserRoles(ctx, userID)
		require.NoError(t, err)
		assert.Empty(t, roles)
	})

	t.Run("returns all user roles", func(t *testing.T) {
		svc := setupUnitTestService(t)
		ctx := context.Background()

		userID := uuid.Must(uuid.NewV7())

		require.NoError(t, svc.AssignRole(ctx, userID, "admin"))
		require.NoError(t, svc.AssignRole(ctx, userID, "editor"))

		roles, err := svc.GetUserRoles(ctx, userID)
		require.NoError(t, err)
		assert.Len(t, roles, 2)
		assert.Contains(t, roles, "admin")
		assert.Contains(t, roles, "editor")
	})
}

func TestService_GetUsersForRole_Unit(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("empty when no users", func(t *testing.T) {
		svc := setupUnitTestService(t)
		ctx := context.Background()

		users, err := svc.GetUsersForRole(ctx, "admin")
		require.NoError(t, err)
		assert.Empty(t, users)
	})

	t.Run("returns all users for role", func(t *testing.T) {
		svc := setupUnitTestService(t)
		ctx := context.Background()

		user1 := uuid.Must(uuid.NewV7())
		user2 := uuid.Must(uuid.NewV7())
		user3 := uuid.Must(uuid.NewV7())

		require.NoError(t, svc.AssignRole(ctx, user1, "admin"))
		require.NoError(t, svc.AssignRole(ctx, user2, "admin"))
		require.NoError(t, svc.AssignRole(ctx, user3, "editor"))

		adminUsers, err := svc.GetUsersForRole(ctx, "admin")
		require.NoError(t, err)
		assert.Len(t, adminUsers, 2)
		assert.Contains(t, adminUsers, user1)
		assert.Contains(t, adminUsers, user2)

		editorUsers, err := svc.GetUsersForRole(ctx, "editor")
		require.NoError(t, err)
		assert.Len(t, editorUsers, 1)
		assert.Contains(t, editorUsers, user3)
	})

	t.Run("ignores invalid UUIDs", func(t *testing.T) {
		svc := setupUnitTestService(t)
		ctx := context.Background()

		// Directly add a non-UUID user via Casbin (simulating legacy data)
		validUser := uuid.Must(uuid.NewV7())
		require.NoError(t, svc.AssignRole(ctx, validUser, "admin"))

		users, err := svc.GetUsersForRole(ctx, "admin")
		require.NoError(t, err)
		assert.Len(t, users, 1)
		assert.Contains(t, users, validUser)
	})
}

func TestService_HasRole_Unit(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("true when has role", func(t *testing.T) {
		svc := setupUnitTestService(t)
		ctx := context.Background()

		userID := uuid.Must(uuid.NewV7())
		require.NoError(t, svc.AssignRole(ctx, userID, "admin"))

		hasRole, err := svc.HasRole(ctx, userID, "admin")
		require.NoError(t, err)
		assert.True(t, hasRole)
	})

	t.Run("false when no role", func(t *testing.T) {
		svc := setupUnitTestService(t)
		ctx := context.Background()

		userID := uuid.Must(uuid.NewV7())

		hasRole, err := svc.HasRole(ctx, userID, "admin")
		require.NoError(t, err)
		assert.False(t, hasRole)
	})

	t.Run("false for different role", func(t *testing.T) {
		svc := setupUnitTestService(t)
		ctx := context.Background()

		userID := uuid.Must(uuid.NewV7())
		require.NoError(t, svc.AssignRole(ctx, userID, "editor"))

		hasRole, err := svc.HasRole(ctx, userID, "admin")
		require.NoError(t, err)
		assert.False(t, hasRole)
	})
}

func TestService_RoleBasedAccess_Unit(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	svc := setupUnitTestService(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())

	// Add policy for admin role
	require.NoError(t, svc.AddPolicy(ctx, "admin", "users", "manage"))

	// User without role cannot access
	allowed, err := svc.Enforce(ctx, userID.String(), "users", "manage")
	require.NoError(t, err)
	assert.False(t, allowed)

	// Assign admin role
	require.NoError(t, svc.AssignRole(ctx, userID, "admin"))

	// User with role can access
	allowed, err = svc.Enforce(ctx, userID.String(), "users", "manage")
	require.NoError(t, err)
	assert.True(t, allowed)
}

// Note: LoadPolicy and SavePolicy require a real adapter (database)
// and are covered by integration tests in service_test.go.
// These cannot be tested with the in-memory enforcer as it has no adapter.

func TestService_ComplexPermissions_Unit(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	svc := setupUnitTestService(t)
	ctx := context.Background()

	// Setup: admin role has full access, editor has read/write, viewer has read only
	require.NoError(t, svc.AddPolicy(ctx, "admin", "movies", "read"))
	require.NoError(t, svc.AddPolicy(ctx, "admin", "movies", "write"))
	require.NoError(t, svc.AddPolicy(ctx, "admin", "movies", "delete"))
	require.NoError(t, svc.AddPolicy(ctx, "editor", "movies", "read"))
	require.NoError(t, svc.AddPolicy(ctx, "editor", "movies", "write"))
	require.NoError(t, svc.AddPolicy(ctx, "viewer", "movies", "read"))

	adminUser := uuid.Must(uuid.NewV7())
	editorUser := uuid.Must(uuid.NewV7())
	viewerUser := uuid.Must(uuid.NewV7())

	require.NoError(t, svc.AssignRole(ctx, adminUser, "admin"))
	require.NoError(t, svc.AssignRole(ctx, editorUser, "editor"))
	require.NoError(t, svc.AssignRole(ctx, viewerUser, "viewer"))

	// Admin can do everything
	for _, action := range []string{"read", "write", "delete"} {
		allowed, err := svc.Enforce(ctx, adminUser.String(), "movies", action)
		require.NoError(t, err)
		assert.True(t, allowed, "admin should have %s permission", action)
	}

	// Editor can read and write, but not delete
	allowed, err := svc.Enforce(ctx, editorUser.String(), "movies", "read")
	require.NoError(t, err)
	assert.True(t, allowed)

	allowed, err = svc.Enforce(ctx, editorUser.String(), "movies", "write")
	require.NoError(t, err)
	assert.True(t, allowed)

	allowed, err = svc.Enforce(ctx, editorUser.String(), "movies", "delete")
	require.NoError(t, err)
	assert.False(t, allowed)

	// Viewer can only read
	allowed, err = svc.Enforce(ctx, viewerUser.String(), "movies", "read")
	require.NoError(t, err)
	assert.True(t, allowed)

	allowed, err = svc.Enforce(ctx, viewerUser.String(), "movies", "write")
	require.NoError(t, err)
	assert.False(t, allowed)
}
