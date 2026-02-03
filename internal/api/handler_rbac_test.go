package api

import (
	"context"
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/service/rbac"
	"github.com/lusoris/revenge/internal/testutil"
)

func setupRBACTestHandler(t *testing.T) (*Handler, *testutil.TestDB, uuid.UUID, uuid.UUID) {
	t.Helper()
	testDB := testutil.NewTestDB(t)

	// Clear any existing policies
	_, err := testDB.Pool().Exec(context.Background(), "DELETE FROM shared.casbin_rule")
	require.NoError(t, err)

	// Set up RBAC service with Casbin
	adapter := rbac.NewAdapter(testDB.Pool())
	modelPath := "../../config/casbin_model.conf"
	enforcer, err := casbin.NewEnforcer(modelPath, adapter)
	require.NoError(t, err)
	rbacService := rbac.NewService(enforcer, zap.NewNop())

	// Create admin user
	adminUser := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "admin",
		Email:    "admin@example.com",
	})

	// Create regular user
	regularUser := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "regular",
		Email:    "regular@example.com",
	})

	// Grant admin role
	err = rbacService.AssignRole(context.Background(), adminUser.ID, "admin")
	require.NoError(t, err)

	handler := &Handler{
		logger:      zap.NewNop(),
		rbacService: rbacService,
	}

	return handler, testDB, adminUser.ID, regularUser.ID
}

// ListPolicies tests

func TestHandler_ListPolicies_NoAuth(t *testing.T) {
	t.Parallel()
	handler, _, _, _ := setupRBACTestHandler(t)

	ctx := context.Background()

	result, err := handler.ListPolicies(ctx)
	require.NoError(t, err)

	_, ok := result.(*ogen.ListPoliciesUnauthorized)
	require.True(t, ok)
}

func TestHandler_ListPolicies_NotAdmin(t *testing.T) {
	t.Parallel()
	handler, _, _, regularUserID := setupRBACTestHandler(t)

	ctx := contextWithUserID(context.Background(), regularUserID)

	result, err := handler.ListPolicies(ctx)
	require.NoError(t, err)

	forbidden, ok := result.(*ogen.ListPoliciesForbidden)
	require.True(t, ok)
	assert.Equal(t, 403, forbidden.Code)
}

func TestHandler_ListPolicies_Success(t *testing.T) {
	t.Parallel()
	handler, _, adminID, _ := setupRBACTestHandler(t)

	// Add some policies
	err := handler.rbacService.AddPolicy(context.Background(), "alice", "data1", "read")
	require.NoError(t, err)
	err = handler.rbacService.AddPolicy(context.Background(), "bob", "data2", "write")
	require.NoError(t, err)

	ctx := contextWithUserID(context.Background(), adminID)

	result, err := handler.ListPolicies(ctx)
	require.NoError(t, err)

	response, ok := result.(*ogen.PolicyListResponse)
	require.True(t, ok)
	assert.GreaterOrEqual(t, len(response.Policies), 2)
}

// AddPolicy tests

func TestHandler_AddPolicy_NoAuth(t *testing.T) {
	t.Parallel()
	handler, _, _, _ := setupRBACTestHandler(t)

	ctx := context.Background()
	req := &ogen.PolicyRequest{
		Subject: "alice",
		Object:  "data1",
		Action:  "read",
	}

	result, err := handler.AddPolicy(ctx, req)
	require.NoError(t, err)

	_, ok := result.(*ogen.AddPolicyUnauthorized)
	require.True(t, ok)
}

func TestHandler_AddPolicy_NotAdmin(t *testing.T) {
	t.Parallel()
	handler, _, _, regularUserID := setupRBACTestHandler(t)

	ctx := contextWithUserID(context.Background(), regularUserID)
	req := &ogen.PolicyRequest{
		Subject: "alice",
		Object:  "data1",
		Action:  "read",
	}

	result, err := handler.AddPolicy(ctx, req)
	require.NoError(t, err)

	forbidden, ok := result.(*ogen.AddPolicyForbidden)
	require.True(t, ok)
	assert.Equal(t, 403, forbidden.Code)
}

func TestHandler_AddPolicy_Success(t *testing.T) {
	t.Parallel()
	handler, _, adminID, _ := setupRBACTestHandler(t)

	ctx := contextWithUserID(context.Background(), adminID)
	req := &ogen.PolicyRequest{
		Subject: "alice",
		Object:  "data1",
		Action:  "read",
	}

	result, err := handler.AddPolicy(ctx, req)
	require.NoError(t, err)

	_, ok := result.(*ogen.AddPolicyCreated)
	require.True(t, ok)

	// Verify policy was added
	allowed, err := handler.rbacService.Enforce(context.Background(), "alice", "data1", "read")
	require.NoError(t, err)
	assert.True(t, allowed)
}

// RemovePolicy tests

func TestHandler_RemovePolicy_NoAuth(t *testing.T) {
	t.Parallel()
	handler, _, _, _ := setupRBACTestHandler(t)

	ctx := context.Background()
	req := &ogen.PolicyRequest{
		Subject: "alice",
		Object:  "data1",
		Action:  "read",
	}

	result, err := handler.RemovePolicy(ctx, req)
	require.NoError(t, err)

	_, ok := result.(*ogen.RemovePolicyUnauthorized)
	require.True(t, ok)
}

func TestHandler_RemovePolicy_NotAdmin(t *testing.T) {
	t.Parallel()
	handler, _, _, regularUserID := setupRBACTestHandler(t)

	ctx := contextWithUserID(context.Background(), regularUserID)
	req := &ogen.PolicyRequest{
		Subject: "alice",
		Object:  "data1",
		Action:  "read",
	}

	result, err := handler.RemovePolicy(ctx, req)
	require.NoError(t, err)

	forbidden, ok := result.(*ogen.RemovePolicyForbidden)
	require.True(t, ok)
	assert.Equal(t, 403, forbidden.Code)
}

func TestHandler_RemovePolicy_Success(t *testing.T) {
	t.Parallel()
	handler, _, adminID, _ := setupRBACTestHandler(t)

	// Add policy first
	err := handler.rbacService.AddPolicy(context.Background(), "alice", "data1", "read")
	require.NoError(t, err)

	ctx := contextWithUserID(context.Background(), adminID)
	req := &ogen.PolicyRequest{
		Subject: "alice",
		Object:  "data1",
		Action:  "read",
	}

	result, err := handler.RemovePolicy(ctx, req)
	require.NoError(t, err)

	_, ok := result.(*ogen.RemovePolicyNoContent)
	require.True(t, ok)

	// Verify policy was removed
	allowed, err := handler.rbacService.Enforce(context.Background(), "alice", "data1", "read")
	require.NoError(t, err)
	assert.False(t, allowed)
}

func TestHandler_RemovePolicy_NotFound(t *testing.T) {
	t.Parallel()
	handler, _, adminID, _ := setupRBACTestHandler(t)

	ctx := contextWithUserID(context.Background(), adminID)
	req := &ogen.PolicyRequest{
		Subject: "alice",
		Object:  "nonexistent",
		Action:  "read",
	}

	result, err := handler.RemovePolicy(ctx, req)
	require.NoError(t, err)

	_, ok := result.(*ogen.RemovePolicyNotFound)
	require.True(t, ok)
}

// GetUserRoles tests

func TestHandler_GetUserRoles_NoAuth(t *testing.T) {
	t.Parallel()
	handler, _, _, regularUserID := setupRBACTestHandler(t)

	ctx := context.Background()
	params := ogen.GetUserRolesParams{
		UserId: regularUserID,
	}

	result, err := handler.GetUserRoles(ctx, params)
	require.NoError(t, err)

	errResponse, ok := result.(*ogen.Error)
	require.True(t, ok)
	assert.Equal(t, 401, errResponse.Code)
}

func TestHandler_GetUserRoles_Success(t *testing.T) {
	t.Parallel()
	handler, _, adminID, regularUserID := setupRBACTestHandler(t)

	// Assign a role to regular user
	err := handler.rbacService.AssignRole(context.Background(), regularUserID, "editor")
	require.NoError(t, err)

	ctx := contextWithUserID(context.Background(), adminID)
	params := ogen.GetUserRolesParams{
		UserId: regularUserID,
	}

	result, err := handler.GetUserRoles(ctx, params)
	require.NoError(t, err)

	response, ok := result.(*ogen.RoleListResponse)
	require.True(t, ok)
	assert.Contains(t, response.Roles, "editor")
}

// AssignRole tests

func TestHandler_AssignRole_NoAuth(t *testing.T) {
	t.Parallel()
	handler, _, _, regularUserID := setupRBACTestHandler(t)

	ctx := context.Background()
	req := &ogen.AssignRoleRequest{
		Role: "editor",
	}
	params := ogen.AssignRoleParams{
		UserId: regularUserID,
	}

	result, err := handler.AssignRole(ctx, req, params)
	require.NoError(t, err)

	_, ok := result.(*ogen.AssignRoleUnauthorized)
	require.True(t, ok)
}

func TestHandler_AssignRole_NotAdmin(t *testing.T) {
	t.Parallel()
	handler, _, _, regularUserID := setupRBACTestHandler(t)

	ctx := contextWithUserID(context.Background(), regularUserID)
	req := &ogen.AssignRoleRequest{
		Role: "editor",
	}
	params := ogen.AssignRoleParams{
		UserId: regularUserID,
	}

	result, err := handler.AssignRole(ctx, req, params)
	require.NoError(t, err)

	forbidden, ok := result.(*ogen.AssignRoleForbidden)
	require.True(t, ok)
	assert.Equal(t, 403, forbidden.Code)
}

func TestHandler_AssignRole_Success(t *testing.T) {
	t.Parallel()
	handler, _, adminID, regularUserID := setupRBACTestHandler(t)

	ctx := contextWithUserID(context.Background(), adminID)
	req := &ogen.AssignRoleRequest{
		Role: "editor",
	}
	params := ogen.AssignRoleParams{
		UserId: regularUserID,
	}

	result, err := handler.AssignRole(ctx, req, params)
	require.NoError(t, err)

	_, ok := result.(*ogen.AssignRoleCreated)
	require.True(t, ok)

	// Verify role was assigned
	hasRole, err := handler.rbacService.HasRole(context.Background(), regularUserID, "editor")
	require.NoError(t, err)
	assert.True(t, hasRole)
}

// RemoveRole tests

func TestHandler_RemoveRole_NoAuth(t *testing.T) {
	t.Parallel()
	handler, _, _, regularUserID := setupRBACTestHandler(t)

	ctx := context.Background()
	params := ogen.RemoveRoleParams{
		UserId: regularUserID,
		Role:   "editor",
	}

	result, err := handler.RemoveRole(ctx, params)
	require.NoError(t, err)

	_, ok := result.(*ogen.RemoveRoleUnauthorized)
	require.True(t, ok)
}

func TestHandler_RemoveRole_NotAdmin(t *testing.T) {
	t.Parallel()
	handler, _, _, regularUserID := setupRBACTestHandler(t)

	ctx := contextWithUserID(context.Background(), regularUserID)
	params := ogen.RemoveRoleParams{
		UserId: regularUserID,
		Role:   "editor",
	}

	result, err := handler.RemoveRole(ctx, params)
	require.NoError(t, err)

	forbidden, ok := result.(*ogen.RemoveRoleForbidden)
	require.True(t, ok)
	assert.Equal(t, 403, forbidden.Code)
}

func TestHandler_RemoveRole_Success(t *testing.T) {
	t.Parallel()
	handler, _, adminID, regularUserID := setupRBACTestHandler(t)

	// Assign role first
	err := handler.rbacService.AssignRole(context.Background(), regularUserID, "editor")
	require.NoError(t, err)

	ctx := contextWithUserID(context.Background(), adminID)
	params := ogen.RemoveRoleParams{
		UserId: regularUserID,
		Role:   "editor",
	}

	result, err := handler.RemoveRole(ctx, params)
	require.NoError(t, err)

	_, ok := result.(*ogen.RemoveRoleNoContent)
	require.True(t, ok)

	// Verify role was removed
	hasRole, err := handler.rbacService.HasRole(context.Background(), regularUserID, "editor")
	require.NoError(t, err)
	assert.False(t, hasRole)
}

func TestHandler_RemoveRole_NotFound(t *testing.T) {
	t.Parallel()
	handler, _, adminID, regularUserID := setupRBACTestHandler(t)

	ctx := contextWithUserID(context.Background(), adminID)
	params := ogen.RemoveRoleParams{
		UserId: regularUserID,
		Role:   "nonexistent",
	}

	result, err := handler.RemoveRole(ctx, params)
	require.NoError(t, err)

	_, ok := result.(*ogen.RemoveRoleNotFound)
	require.True(t, ok)
}
