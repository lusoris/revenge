package rbac

import (
	"context"
	"os"
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/lusoris/revenge/internal/service/activity"
	"github.com/lusoris/revenge/internal/testutil"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testutil.StopSharedPostgres()
	os.Exit(code)
}

func setupTestService(t *testing.T) (*Service, testutil.DB) {
	t.Helper()
	testDB := testutil.NewFastTestDB(t)

	// Clear any existing policies from the table to ensure test isolation
	_, err := testDB.Pool().Exec(context.Background(), "DELETE FROM shared.casbin_rule")
	require.NoError(t, err)

	adapter := NewAdapter(testDB.Pool())

	// Use the actual model file from the project
	modelPath := "../../../config/casbin_model.conf"
	enforcer, err := casbin.NewSyncedEnforcer(modelPath, adapter)
	require.NoError(t, err)

	logger := zaptest.NewLogger(t)
	svc := NewService(enforcer, logger, activity.NewNoopLogger())

	return svc, testDB
}

func TestService_AddPolicy(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	err := svc.AddPolicy(ctx, "alice", "data1", "read")
	require.NoError(t, err)

	allowed, err := svc.Enforce(ctx, "alice", "data1", "read")
	require.NoError(t, err)
	assert.True(t, allowed)

	allowed, err = svc.Enforce(ctx, "alice", "data1", "write")
	require.NoError(t, err)
	assert.False(t, allowed)
}

func TestService_AddPolicy_Duplicate(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	err := svc.AddPolicy(ctx, "alice", "data1", "read")
	require.NoError(t, err)

	err = svc.AddPolicy(ctx, "alice", "data1", "read")
	require.NoError(t, err)
}

func TestService_RemovePolicy(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	err := svc.AddPolicy(ctx, "alice", "data1", "read")
	require.NoError(t, err)

	err = svc.RemovePolicy(ctx, "alice", "data1", "read")
	require.NoError(t, err)

	allowed, err := svc.Enforce(ctx, "alice", "data1", "read")
	require.NoError(t, err)
	assert.False(t, allowed)
}

func TestService_RemovePolicy_NotFound(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	err := svc.RemovePolicy(ctx, "alice", "data1", "read")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestService_GetPolicies(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	require.NoError(t, svc.AddPolicy(ctx, "alice", "data1", "read"))
	require.NoError(t, svc.AddPolicy(ctx, "alice", "data2", "write"))
	require.NoError(t, svc.AddPolicy(ctx, "bob", "data1", "read"))

	policies, err := svc.GetPolicies(ctx)
	require.NoError(t, err)
	assert.Len(t, policies, 3)

	containsPolicy := func(policies [][]string, sub, obj, act string) bool {
		for _, p := range policies {
			if len(p) >= 3 && p[0] == sub && p[1] == obj && p[2] == act {
				return true
			}
		}
		return false
	}

	assert.True(t, containsPolicy(policies, "alice", "data1", "read"))
	assert.True(t, containsPolicy(policies, "alice", "data2", "write"))
	assert.True(t, containsPolicy(policies, "bob", "data1", "read"))
}

func TestService_AssignRole(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())

	err := svc.AssignRole(ctx, userID, "admin")
	require.NoError(t, err)

	hasRole, err := svc.HasRole(ctx, userID, "admin")
	require.NoError(t, err)
	assert.True(t, hasRole)
}

func TestService_AssignRole_Duplicate(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())

	err := svc.AssignRole(ctx, userID, "admin")
	require.NoError(t, err)

	err = svc.AssignRole(ctx, userID, "admin")
	require.NoError(t, err)
}

func TestService_RemoveRole(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())

	err := svc.AssignRole(ctx, userID, "admin")
	require.NoError(t, err)

	err = svc.RemoveRole(ctx, userID, "admin")
	require.NoError(t, err)

	hasRole, err := svc.HasRole(ctx, userID, "admin")
	require.NoError(t, err)
	assert.False(t, hasRole)
}

func TestService_RemoveRole_NotFound(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())

	err := svc.RemoveRole(ctx, userID, "admin")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestService_GetUserRoles(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())

	require.NoError(t, svc.AssignRole(ctx, userID, "admin"))
	require.NoError(t, svc.AssignRole(ctx, userID, "editor"))

	roles, err := svc.GetUserRoles(ctx, userID)
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"admin", "editor"}, roles)
}

func TestService_GetUserRoles_NoRoles(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())

	roles, err := svc.GetUserRoles(ctx, userID)
	require.NoError(t, err)
	assert.Empty(t, roles)
}

func TestService_GetUsersForRole(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	user1 := uuid.Must(uuid.NewV7())
	user2 := uuid.Must(uuid.NewV7())

	require.NoError(t, svc.AssignRole(ctx, user1, "admin"))
	require.NoError(t, svc.AssignRole(ctx, user2, "admin"))

	users, err := svc.GetUsersForRole(ctx, "admin")
	require.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Contains(t, users, user1)
	assert.Contains(t, users, user2)
}

func TestService_GetUsersForRole_NoUsers(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	users, err := svc.GetUsersForRole(ctx, "admin")
	require.NoError(t, err)
	assert.Empty(t, users)
}

func TestService_HasRole(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())

	hasRole, err := svc.HasRole(ctx, userID, "admin")
	require.NoError(t, err)
	assert.False(t, hasRole)

	require.NoError(t, svc.AssignRole(ctx, userID, "admin"))

	hasRole, err = svc.HasRole(ctx, userID, "admin")
	require.NoError(t, err)
	assert.True(t, hasRole)
}

func TestService_Enforce(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	require.NoError(t, svc.AddPolicy(ctx, "alice", "data1", "read"))

	allowed, err := svc.Enforce(ctx, "alice", "data1", "read")
	require.NoError(t, err)
	assert.True(t, allowed)

	allowed, err = svc.Enforce(ctx, "alice", "data1", "write")
	require.NoError(t, err)
	assert.False(t, allowed)

	allowed, err = svc.Enforce(ctx, "bob", "data1", "read")
	require.NoError(t, err)
	assert.False(t, allowed)
}

func TestService_EnforceWithContext(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())

	require.NoError(t, svc.AddPolicy(ctx, userID.String(), "resource1", "read"))

	allowed, err := svc.EnforceWithContext(ctx, userID, "resource1", "read")
	require.NoError(t, err)
	assert.True(t, allowed)
}

func TestService_RoleBasedAccess(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())

	require.NoError(t, svc.AddPolicy(ctx, "admin", "data1", "write"))

	allowed, err := svc.Enforce(ctx, userID.String(), "data1", "write")
	require.NoError(t, err)
	assert.False(t, allowed)

	require.NoError(t, svc.AssignRole(ctx, userID, "admin"))

	allowed, err = svc.Enforce(ctx, userID.String(), "data1", "write")
	require.NoError(t, err)
	assert.True(t, allowed)
}

func TestService_LoadPolicy(t *testing.T) {
	t.Parallel()
	svc, testDB := setupTestService(t)
	ctx := context.Background()

	// Add policies using the first service
	require.NoError(t, svc.AddPolicy(ctx, "alice", "data1", "read"))
	require.NoError(t, svc.AddPolicy(ctx, "bob", "data2", "write"))

	// Create a new service with the same database (simulating a restart)
	adapter := NewAdapter(testDB.Pool())
	modelPath := "../../../config/casbin_model.conf"
	enforcer, err := casbin.NewSyncedEnforcer(modelPath, adapter)
	require.NoError(t, err)

	logger := zaptest.NewLogger(t)
	newSvc := NewService(enforcer, logger, activity.NewNoopLogger())

	// Load policies from the database
	err = newSvc.LoadPolicy(ctx)
	require.NoError(t, err)

	// Verify the policies were loaded
	allowed, err := newSvc.Enforce(ctx, "alice", "data1", "read")
	require.NoError(t, err)
	assert.True(t, allowed)

	allowed, err = newSvc.Enforce(ctx, "bob", "data2", "write")
	require.NoError(t, err)
	assert.True(t, allowed)
}

func TestService_SavePolicy(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	require.NoError(t, svc.AddPolicy(ctx, "alice", "data1", "read"))
	require.NoError(t, svc.AddPolicy(ctx, "bob", "data2", "write"))

	err := svc.SavePolicy(ctx)
	require.NoError(t, err)

	policies, err := svc.GetPolicies(ctx)
	require.NoError(t, err)
	assert.Len(t, policies, 2)
}
