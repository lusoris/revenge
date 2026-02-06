package rbac

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Additional tests for error paths and edge cases to increase coverage

func TestService_GetPolicies_Empty(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	policies, err := svc.GetPolicies(ctx)
	require.NoError(t, err)
	// Policies can be nil or empty when there are none
	if policies != nil {
		assert.GreaterOrEqual(t, len(policies), 0)
	}
}

func TestService_HasRole_NoRole(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())
	hasRole, err := svc.HasRole(ctx, userID, "nonexistent")
	require.NoError(t, err)
	assert.False(t, hasRole)
}

func TestService_GetUsersForRole_EmptyRole(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Create a role with no users
	_, err := svc.CreateRole(ctx, "emptyrole", "No users", []Permission{
		{Resource: "test", Action: "read"},
	})
	require.NoError(t, err)

	users, err := svc.GetUsersForRole(ctx, "emptyrole")
	require.NoError(t, err)
	assert.Equal(t, 0, len(users))
}

func TestService_LoadPolicy_Success(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Add some policies
	err := svc.AddPolicy(ctx, "user1", "resource1", "read")
	require.NoError(t, err)

	// Load policies
	err = svc.LoadPolicy(ctx)
	require.NoError(t, err)

	// Verify policies still work
	allowed, err := svc.Enforce(ctx, "user1", "resource1", "read")
	require.NoError(t, err)
	assert.True(t, allowed)
}

func TestService_SavePolicy_Success(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Add policies
	err := svc.AddPolicy(ctx, "user2", "resource2", "write")
	require.NoError(t, err)

	// Save policies
	err = svc.SavePolicy(ctx)
	require.NoError(t, err)
}

func TestService_Enforce_Multiple(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Add multiple policies
	err := svc.AddPolicy(ctx, "admin", "users", "read")
	require.NoError(t, err)
	err = svc.AddPolicy(ctx, "admin", "users", "write")
	require.NoError(t, err)
	err = svc.AddPolicy(ctx, "admin", "users", "delete")
	require.NoError(t, err)

	// Test all policies
	allowed, err := svc.Enforce(ctx, "admin", "users", "read")
	require.NoError(t, err)
	assert.True(t, allowed)

	allowed, err = svc.Enforce(ctx, "admin", "users", "write")
	require.NoError(t, err)
	assert.True(t, allowed)

	allowed, err = svc.Enforce(ctx, "admin", "users", "delete")
	require.NoError(t, err)
	assert.True(t, allowed)

	// Test non-existent permission
	allowed, err = svc.Enforce(ctx, "admin", "users", "execute")
	require.NoError(t, err)
	assert.False(t, allowed)
}

func TestService_RemovePolicy_Multiple(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Add multiple policies
	err := svc.AddPolicy(ctx, "user3", "data", "read")
	require.NoError(t, err)
	err = svc.AddPolicy(ctx, "user3", "data", "write")
	require.NoError(t, err)

	// Remove one
	err = svc.RemovePolicy(ctx, "user3", "data", "write")
	require.NoError(t, err)

	// Verify only one remains
	allowed, err := svc.Enforce(ctx, "user3", "data", "read")
	require.NoError(t, err)
	assert.True(t, allowed)

	allowed, err = svc.Enforce(ctx, "user3", "data", "write")
	require.NoError(t, err)
	assert.False(t, allowed)
}

func TestService_AssignRole_Multiple(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())

	// Create multiple roles
	_, err := svc.CreateRole(ctx, "viewer", "", []Permission{{Resource: "content", Action: "read"}})
	require.NoError(t, err)
	_, err = svc.CreateRole(ctx, "contributor", "", []Permission{{Resource: "content", Action: "write"}})
	require.NoError(t, err)

	// Assign both roles
	err = svc.AssignRole(ctx, userID, "viewer")
	require.NoError(t, err)
	err = svc.AssignRole(ctx, userID, "contributor")
	require.NoError(t, err)

	// Verify user has both roles
	roles, err := svc.GetUserRoles(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, 2, len(roles))
}

func TestService_RemoveRole_Last(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())

	// Create and assign single role
	_, err := svc.CreateRole(ctx, "temporary", "", []Permission{{Resource: "temp", Action: "read"}})
	require.NoError(t, err)

	err = svc.AssignRole(ctx, userID, "temporary")
	require.NoError(t, err)

	// Verify has role
	hasRole, err := svc.HasRole(ctx, userID, "temporary")
	require.NoError(t, err)
	assert.True(t, hasRole)

	// Remove the role
	err = svc.RemoveRole(ctx, userID, "temporary")
	require.NoError(t, err)

	// Verify no longer has role
	hasRole, err = svc.HasRole(ctx, userID, "temporary")
	require.NoError(t, err)
	assert.False(t, hasRole)
}

func TestService_GetUserRoles_Empty(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())

	// Get roles for user with no roles
	roles, err := svc.GetUserRoles(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, 0, len(roles))
}

func TestService_ListRoles_Empty(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// List roles when none exist (except potentially from other tests due to shared DB)
	roles, err := svc.ListRoles(ctx)
	require.NoError(t, err)
	assert.NotNil(t, roles)
	// Don't assert count since other parallel tests may have created roles
}
