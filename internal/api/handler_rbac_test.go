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
	"github.com/lusoris/revenge/internal/service/activity"
	"github.com/lusoris/revenge/internal/service/rbac"
	"github.com/lusoris/revenge/internal/testutil"
)

func setupRBACTestHandler(t *testing.T) (*Handler, *testutil.TestDB, uuid.UUID, uuid.UUID) {
	t.Helper()
	testDB := testutil.NewTestDB(t)

	// Clear any existing policies
	_, err := testDB.Pool().Exec(context.Background(), "DELETE FROM shared.casbin_rule")
	require.NoError(t, err)

	// Insert default policies for built-in roles
	_, err = testDB.Pool().Exec(context.Background(), `
		INSERT INTO shared.casbin_rule (ptype, v0, v1, v2) VALUES
			('p', 'admin', '*', '*'),
			('p', 'user', 'profile', 'read'),
			('p', 'user', 'profile', 'write'),
			('p', 'user', 'library', 'read'),
			('p', 'user', 'playback', 'read'),
			('p', 'user', 'playback', 'write'),
			('p', 'guest', 'library', 'read')
	`)
	require.NoError(t, err)

	// Set up RBAC service with Casbin
	adapter := rbac.NewAdapter(testDB.Pool())
	modelPath := "../../config/casbin_model.conf"
	enforcer, err := casbin.NewEnforcer(modelPath, adapter)
	require.NoError(t, err)
	rbacService := rbac.NewService(enforcer, zap.NewNop(), activity.NewNoopLogger())

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

// ListRoles tests

func TestHandler_ListRoles_NoAuth(t *testing.T) {
	t.Parallel()
	handler, _, _, _ := setupRBACTestHandler(t)

	ctx := context.Background()

	result, err := handler.ListRoles(ctx)
	require.NoError(t, err)

	_, ok := result.(*ogen.ListRolesUnauthorized)
	require.True(t, ok)
}

func TestHandler_ListRoles_NotAdmin(t *testing.T) {
	t.Parallel()
	handler, _, _, regularUserID := setupRBACTestHandler(t)

	ctx := contextWithUserID(context.Background(), regularUserID)

	result, err := handler.ListRoles(ctx)
	require.NoError(t, err)

	forbidden, ok := result.(*ogen.ListRolesForbidden)
	require.True(t, ok)
	assert.Equal(t, 403, forbidden.Code)
}

func TestHandler_ListRoles_Success(t *testing.T) {
	t.Parallel()
	handler, _, adminID, _ := setupRBACTestHandler(t)

	ctx := contextWithUserID(context.Background(), adminID)

	result, err := handler.ListRoles(ctx)
	require.NoError(t, err)

	response, ok := result.(*ogen.RolesResponse)
	require.True(t, ok)
	assert.NotEmpty(t, response.Roles)

	// Check that admin role exists and is built-in
	var foundAdmin bool
	for _, role := range response.Roles {
		if role.Name == "admin" {
			foundAdmin = true
			assert.True(t, role.IsBuiltIn)
		}
	}
	assert.True(t, foundAdmin, "admin role should be present")
}

// GetRole tests

func TestHandler_GetRole_NoAuth(t *testing.T) {
	t.Parallel()
	handler, _, _, _ := setupRBACTestHandler(t)

	ctx := context.Background()
	params := ogen.GetRoleParams{RoleName: "admin"}

	result, err := handler.GetRole(ctx, params)
	require.NoError(t, err)

	_, ok := result.(*ogen.GetRoleUnauthorized)
	require.True(t, ok)
}

func TestHandler_GetRole_NotAdmin(t *testing.T) {
	t.Parallel()
	handler, _, _, regularUserID := setupRBACTestHandler(t)

	ctx := contextWithUserID(context.Background(), regularUserID)
	params := ogen.GetRoleParams{RoleName: "admin"}

	result, err := handler.GetRole(ctx, params)
	require.NoError(t, err)

	forbidden, ok := result.(*ogen.GetRoleForbidden)
	require.True(t, ok)
	assert.Equal(t, 403, forbidden.Code)
}

func TestHandler_GetRole_Success(t *testing.T) {
	t.Parallel()
	handler, _, adminID, _ := setupRBACTestHandler(t)

	ctx := contextWithUserID(context.Background(), adminID)
	params := ogen.GetRoleParams{RoleName: "admin"}

	result, err := handler.GetRole(ctx, params)
	require.NoError(t, err)

	role, ok := result.(*ogen.RoleDetail)
	require.True(t, ok)
	assert.Equal(t, "admin", role.Name)
	assert.True(t, role.IsBuiltIn)
}

func TestHandler_GetRole_NotFound(t *testing.T) {
	t.Parallel()
	handler, _, adminID, _ := setupRBACTestHandler(t)

	ctx := contextWithUserID(context.Background(), adminID)
	params := ogen.GetRoleParams{RoleName: "nonexistent"}

	result, err := handler.GetRole(ctx, params)
	require.NoError(t, err)

	_, ok := result.(*ogen.GetRoleNotFound)
	require.True(t, ok)
}

// CreateRole tests

func TestHandler_CreateRole_NoAuth(t *testing.T) {
	t.Parallel()
	handler, _, _, _ := setupRBACTestHandler(t)

	ctx := context.Background()
	req := &ogen.CreateRoleRequest{
		Name: "testrole",
	}

	result, err := handler.CreateRole(ctx, req)
	require.NoError(t, err)

	_, ok := result.(*ogen.CreateRoleUnauthorized)
	require.True(t, ok)
}

func TestHandler_CreateRole_NotAdmin(t *testing.T) {
	t.Parallel()
	handler, _, _, regularUserID := setupRBACTestHandler(t)

	ctx := contextWithUserID(context.Background(), regularUserID)
	req := &ogen.CreateRoleRequest{
		Name: "testrole",
	}

	result, err := handler.CreateRole(ctx, req)
	require.NoError(t, err)

	forbidden, ok := result.(*ogen.CreateRoleForbidden)
	require.True(t, ok)
	assert.Equal(t, 403, forbidden.Code)
}

func TestHandler_CreateRole_Success(t *testing.T) {
	t.Parallel()
	handler, _, adminID, _ := setupRBACTestHandler(t)

	ctx := contextWithUserID(context.Background(), adminID)
	req := &ogen.CreateRoleRequest{
		Name:        "newrole",
		Description: ogen.NewOptString("A new custom role"),
		Permissions: []ogen.Permission{
			{Resource: "library", Action: "read"},
			{Resource: "movies", Action: "read"},
		},
	}

	result, err := handler.CreateRole(ctx, req)
	require.NoError(t, err)

	role, ok := result.(*ogen.RoleDetail)
	require.True(t, ok)
	assert.Equal(t, "newrole", role.Name)
	assert.Equal(t, "A new custom role", role.Description.Value)
	assert.False(t, role.IsBuiltIn)
	assert.Len(t, role.Permissions, 2)
}

func TestHandler_CreateRole_AlreadyExists(t *testing.T) {
	t.Parallel()
	handler, _, adminID, _ := setupRBACTestHandler(t)

	ctx := contextWithUserID(context.Background(), adminID)

	// Create role first with at least one permission
	req := &ogen.CreateRoleRequest{
		Name: "existingrole",
		Permissions: []ogen.Permission{
			{Resource: "library", Action: "read"},
		},
	}
	_, err := handler.CreateRole(ctx, req)
	require.NoError(t, err)

	// Try to create again
	result, err := handler.CreateRole(ctx, req)
	require.NoError(t, err)

	conflict, ok := result.(*ogen.CreateRoleConflict)
	require.True(t, ok)
	assert.Equal(t, 409, conflict.Code)
}

// DeleteRole tests

func TestHandler_DeleteRole_NoAuth(t *testing.T) {
	t.Parallel()
	handler, _, _, _ := setupRBACTestHandler(t)

	ctx := context.Background()
	params := ogen.DeleteRoleParams{RoleName: "testrole"}

	result, err := handler.DeleteRole(ctx, params)
	require.NoError(t, err)

	_, ok := result.(*ogen.DeleteRoleUnauthorized)
	require.True(t, ok)
}

func TestHandler_DeleteRole_NotAdmin(t *testing.T) {
	t.Parallel()
	handler, _, _, regularUserID := setupRBACTestHandler(t)

	ctx := contextWithUserID(context.Background(), regularUserID)
	params := ogen.DeleteRoleParams{RoleName: "testrole"}

	result, err := handler.DeleteRole(ctx, params)
	require.NoError(t, err)

	forbidden, ok := result.(*ogen.DeleteRoleForbidden)
	require.True(t, ok)
	assert.Equal(t, 403, forbidden.Code)
}

func TestHandler_DeleteRole_Success(t *testing.T) {
	t.Parallel()
	handler, _, adminID, _ := setupRBACTestHandler(t)

	ctx := contextWithUserID(context.Background(), adminID)

	// Create role first with at least one permission
	createReq := &ogen.CreateRoleRequest{
		Name: "deleterole",
		Permissions: []ogen.Permission{
			{Resource: "library", Action: "read"},
		},
	}
	_, err := handler.CreateRole(ctx, createReq)
	require.NoError(t, err)

	// Delete role
	params := ogen.DeleteRoleParams{RoleName: "deleterole"}
	result, err := handler.DeleteRole(ctx, params)
	require.NoError(t, err)

	_, ok := result.(*ogen.DeleteRoleNoContent)
	require.True(t, ok)

	// Verify role no longer exists
	getParams := ogen.GetRoleParams{RoleName: "deleterole"}
	getResult, err := handler.GetRole(ctx, getParams)
	require.NoError(t, err)
	_, ok = getResult.(*ogen.GetRoleNotFound)
	assert.True(t, ok)
}

func TestHandler_DeleteRole_BuiltIn(t *testing.T) {
	t.Parallel()
	handler, _, adminID, _ := setupRBACTestHandler(t)

	ctx := contextWithUserID(context.Background(), adminID)
	params := ogen.DeleteRoleParams{RoleName: "admin"}

	result, err := handler.DeleteRole(ctx, params)
	require.NoError(t, err)

	badReq, ok := result.(*ogen.DeleteRoleBadRequest)
	require.True(t, ok)
	assert.Equal(t, 400, badReq.Code)
	assert.Contains(t, badReq.Message, "built-in")
}

func TestHandler_DeleteRole_NotFound(t *testing.T) {
	t.Parallel()
	handler, _, adminID, _ := setupRBACTestHandler(t)

	ctx := contextWithUserID(context.Background(), adminID)
	params := ogen.DeleteRoleParams{RoleName: "nonexistent"}

	result, err := handler.DeleteRole(ctx, params)
	require.NoError(t, err)

	_, ok := result.(*ogen.DeleteRoleNotFound)
	require.True(t, ok)
}

// UpdateRolePermissions tests

func TestHandler_UpdateRolePermissions_NoAuth(t *testing.T) {
	t.Parallel()
	handler, _, _, _ := setupRBACTestHandler(t)

	ctx := context.Background()
	req := &ogen.UpdatePermissionsRequest{
		Permissions: []ogen.Permission{
			{Resource: "library", Action: "read"},
		},
	}
	params := ogen.UpdateRolePermissionsParams{RoleName: "testrole"}

	result, err := handler.UpdateRolePermissions(ctx, req, params)
	require.NoError(t, err)

	_, ok := result.(*ogen.UpdateRolePermissionsUnauthorized)
	require.True(t, ok)
}

func TestHandler_UpdateRolePermissions_NotAdmin(t *testing.T) {
	t.Parallel()
	handler, _, _, regularUserID := setupRBACTestHandler(t)

	ctx := contextWithUserID(context.Background(), regularUserID)
	req := &ogen.UpdatePermissionsRequest{
		Permissions: []ogen.Permission{
			{Resource: "library", Action: "read"},
		},
	}
	params := ogen.UpdateRolePermissionsParams{RoleName: "testrole"}

	result, err := handler.UpdateRolePermissions(ctx, req, params)
	require.NoError(t, err)

	forbidden, ok := result.(*ogen.UpdateRolePermissionsForbidden)
	require.True(t, ok)
	assert.Equal(t, 403, forbidden.Code)
}

func TestHandler_UpdateRolePermissions_Success(t *testing.T) {
	t.Parallel()
	handler, _, adminID, _ := setupRBACTestHandler(t)

	ctx := contextWithUserID(context.Background(), adminID)

	// Create role first
	createReq := &ogen.CreateRoleRequest{
		Name: "permrole",
		Permissions: []ogen.Permission{
			{Resource: "library", Action: "read"},
		},
	}
	_, err := handler.CreateRole(ctx, createReq)
	require.NoError(t, err)

	// Update permissions
	updateReq := &ogen.UpdatePermissionsRequest{
		Permissions: []ogen.Permission{
			{Resource: "library", Action: "read"},
			{Resource: "library", Action: "write"},
			{Resource: "movies", Action: "read"},
		},
	}
	params := ogen.UpdateRolePermissionsParams{RoleName: "permrole"}

	result, err := handler.UpdateRolePermissions(ctx, updateReq, params)
	require.NoError(t, err)

	role, ok := result.(*ogen.RoleDetail)
	require.True(t, ok)
	assert.Equal(t, "permrole", role.Name)
	assert.Len(t, role.Permissions, 3)
}

func TestHandler_UpdateRolePermissions_NotFound(t *testing.T) {
	t.Parallel()
	handler, _, adminID, _ := setupRBACTestHandler(t)

	ctx := contextWithUserID(context.Background(), adminID)
	req := &ogen.UpdatePermissionsRequest{
		Permissions: []ogen.Permission{
			{Resource: "library", Action: "read"},
		},
	}
	params := ogen.UpdateRolePermissionsParams{RoleName: "nonexistent"}

	result, err := handler.UpdateRolePermissions(ctx, req, params)
	require.NoError(t, err)

	_, ok := result.(*ogen.UpdateRolePermissionsNotFound)
	require.True(t, ok)
}

// ListPermissions tests

func TestHandler_ListPermissions_NoAuth(t *testing.T) {
	t.Parallel()
	handler, _, _, _ := setupRBACTestHandler(t)

	ctx := context.Background()

	result, err := handler.ListPermissions(ctx)
	require.NoError(t, err)

	_, ok := result.(*ogen.ListPermissionsUnauthorized)
	require.True(t, ok)
}

func TestHandler_ListPermissions_NotAdmin(t *testing.T) {
	t.Parallel()
	handler, _, _, regularUserID := setupRBACTestHandler(t)

	ctx := contextWithUserID(context.Background(), regularUserID)

	result, err := handler.ListPermissions(ctx)
	require.NoError(t, err)

	forbidden, ok := result.(*ogen.ListPermissionsForbidden)
	require.True(t, ok)
	assert.Equal(t, 403, forbidden.Code)
}

func TestHandler_ListPermissions_Success(t *testing.T) {
	t.Parallel()
	handler, _, adminID, _ := setupRBACTestHandler(t)

	ctx := contextWithUserID(context.Background(), adminID)

	result, err := handler.ListPermissions(ctx)
	require.NoError(t, err)

	response, ok := result.(*ogen.PermissionsResponse)
	require.True(t, ok)
	assert.NotEmpty(t, response.Permissions)

	// Verify that known resources are present
	var foundLibrary, foundUsers bool
	for _, p := range response.Permissions {
		if p.Resource == "library" {
			foundLibrary = true
		}
		if p.Resource == "users" {
			foundUsers = true
		}
	}
	assert.True(t, foundLibrary, "library resource should be present")
	assert.True(t, foundUsers, "users resource should be present")
}
