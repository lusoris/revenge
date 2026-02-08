package rbac

import (
	"context"
	"testing"

	"github.com/casbin/casbin/v2"
	casbinmodel "github.com/casbin/casbin/v2/model"
	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/service/activity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// Tests to fill coverage gaps in service.go, roles.go, and permissions.go.
// Uses in-memory enforcer (no database required).

func TestNewService_FieldsSet(t *testing.T) {
	t.Parallel()

	m, err := casbinmodel.NewModelFromString(casbinModelConfInternal)
	require.NoError(t, err)

	enforcer, err := casbin.NewSyncedEnforcer(m)
	require.NoError(t, err)

	logger := zap.NewNop()
	actLogger := activity.NewNoopLogger()

	svc := NewService(enforcer, logger, actLogger)

	assert.NotNil(t, svc)
	assert.NotNil(t, svc.enforcer)
	assert.NotNil(t, svc.logger)
	assert.NotNil(t, svc.activityLogger)
}

func TestService_Enforce_DeniedNoPolicy(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	// No policies - should deny
	allowed, err := svc.Enforce(ctx, "nobody", "resource", "read")
	require.NoError(t, err)
	assert.False(t, allowed)
}

func TestService_Enforce_AllowedWithRoleAndPolicy(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())

	// Add policy for role
	require.NoError(t, svc.AddPolicy(ctx, "viewer", "movies", "list"))

	// Before role assignment - deny
	allowed, err := svc.Enforce(ctx, userID.String(), "movies", "list")
	require.NoError(t, err)
	assert.False(t, allowed)

	// After role assignment - allow
	require.NoError(t, svc.AssignRole(ctx, userID, "viewer"))

	allowed, err = svc.Enforce(ctx, userID.String(), "movies", "list")
	require.NoError(t, err)
	assert.True(t, allowed)
}

func TestService_EnforceWithContext_DelegatesCorrectly(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())
	require.NoError(t, svc.AddPolicy(ctx, userID.String(), "settings", "read"))

	allowed, err := svc.EnforceWithContext(ctx, userID, "settings", "read")
	require.NoError(t, err)
	assert.True(t, allowed)

	allowed, err = svc.EnforceWithContext(ctx, userID, "settings", "write")
	require.NoError(t, err)
	assert.False(t, allowed)
}

func TestService_AddPolicy_DuplicateNoError(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	// First add
	err := svc.AddPolicy(ctx, "alice", "data", "read")
	require.NoError(t, err)

	// Duplicate add - should not error (just log warning)
	err = svc.AddPolicy(ctx, "alice", "data", "read")
	require.NoError(t, err)

	// Verify only one policy exists
	policies, err := svc.GetPolicies(ctx)
	require.NoError(t, err)
	assert.Len(t, policies, 1)
}

func TestService_RemovePolicy_NotFoundError(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	err := svc.RemovePolicy(ctx, "nonexistent", "data", "read")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestService_GetPolicies_EmptyAndPopulated(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	// Empty
	policies, err := svc.GetPolicies(ctx)
	require.NoError(t, err)
	assert.Empty(t, policies)

	// Add policies
	require.NoError(t, svc.AddPolicy(ctx, "sub1", "obj1", "act1"))
	require.NoError(t, svc.AddPolicy(ctx, "sub2", "obj2", "act2"))

	policies, err = svc.GetPolicies(ctx)
	require.NoError(t, err)
	assert.Len(t, policies, 2)
}

func TestService_AssignRole_ActivityLogged(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())
	err := svc.AssignRole(ctx, userID, "admin")
	require.NoError(t, err)

	// Verify role was assigned
	hasRole, err := svc.HasRole(ctx, userID, "admin")
	require.NoError(t, err)
	assert.True(t, hasRole)
}

func TestService_AssignRole_DuplicateNoError(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())

	err := svc.AssignRole(ctx, userID, "admin")
	require.NoError(t, err)

	// Duplicate should not error
	err = svc.AssignRole(ctx, userID, "admin")
	require.NoError(t, err)

	// Only one role
	roles, err := svc.GetUserRoles(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, roles, 1)
}

func TestService_RemoveRole_NotFoundError(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())
	err := svc.RemoveRole(ctx, userID, "nonexistent")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestService_RemoveRole_ActivityLogged(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())
	require.NoError(t, svc.AssignRole(ctx, userID, "editor"))

	err := svc.RemoveRole(ctx, userID, "editor")
	require.NoError(t, err)

	hasRole, err := svc.HasRole(ctx, userID, "editor")
	require.NoError(t, err)
	assert.False(t, hasRole)
}

func TestService_GetUserRoles_EmptyAndPopulated(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())

	// Empty
	roles, err := svc.GetUserRoles(ctx, userID)
	require.NoError(t, err)
	assert.Empty(t, roles)

	// Add roles
	require.NoError(t, svc.AssignRole(ctx, userID, "admin"))
	require.NoError(t, svc.AssignRole(ctx, userID, "editor"))

	roles, err = svc.GetUserRoles(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, roles, 2)
	assert.Contains(t, roles, "admin")
	assert.Contains(t, roles, "editor")
}

func TestService_GetUsersForRole_WithMultipleUsersAndRoles(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	u1 := uuid.Must(uuid.NewV7())
	u2 := uuid.Must(uuid.NewV7())
	u3 := uuid.Must(uuid.NewV7())

	require.NoError(t, svc.AssignRole(ctx, u1, "admin"))
	require.NoError(t, svc.AssignRole(ctx, u2, "admin"))
	require.NoError(t, svc.AssignRole(ctx, u3, "viewer"))

	admins, err := svc.GetUsersForRole(ctx, "admin")
	require.NoError(t, err)
	assert.Len(t, admins, 2)
	assert.Contains(t, admins, u1)
	assert.Contains(t, admins, u2)

	viewers, err := svc.GetUsersForRole(ctx, "viewer")
	require.NoError(t, err)
	assert.Len(t, viewers, 1)
	assert.Contains(t, viewers, u3)
}

func TestService_GetUsersForRole_NonexistentRoleReturnsEmpty(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	users, err := svc.GetUsersForRole(ctx, "nonexistent_role")
	require.NoError(t, err)
	assert.Empty(t, users)
}

func TestService_HasRole_TrueAndFalse(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())

	// False before assignment
	hasRole, err := svc.HasRole(ctx, userID, "admin")
	require.NoError(t, err)
	assert.False(t, hasRole)

	// True after assignment
	require.NoError(t, svc.AssignRole(ctx, userID, "admin"))
	hasRole, err = svc.HasRole(ctx, userID, "admin")
	require.NoError(t, err)
	assert.True(t, hasRole)

	// False for different role
	hasRole, err = svc.HasRole(ctx, userID, "editor")
	require.NoError(t, err)
	assert.False(t, hasRole)
}

func TestService_CheckUserPermission_Coverage(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())

	// No role -> denied
	allowed, err := svc.CheckUserPermission(ctx, userID, "movies", "list")
	require.NoError(t, err)
	assert.False(t, allowed)

	// Create role with permission and assign
	_, err = svc.CreateRole(ctx, "movieviewer", "", []Permission{
		{Resource: "movies", Action: "list"},
		{Resource: "movies", Action: "get"},
	})
	require.NoError(t, err)
	require.NoError(t, svc.AssignRole(ctx, userID, "movieviewer"))

	// Now allowed
	allowed, err = svc.CheckUserPermission(ctx, userID, "movies", "list")
	require.NoError(t, err)
	assert.True(t, allowed)

	// Different action - still denied
	allowed, err = svc.CheckUserPermission(ctx, userID, "movies", "delete")
	require.NoError(t, err)
	assert.False(t, allowed)
}

func TestService_CreateRole_WithTabInName(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	// Tab character
	_, err := svc.CreateRole(ctx, "role\twith\ttabs", "", []Permission{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "whitespace")
}

func TestService_CreateRole_WithNewlineInName(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	// Newline character
	_, err := svc.CreateRole(ctx, "role\nname", "", []Permission{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "whitespace")
}

func TestService_ListRoles_SortedOutput(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	// Create roles in non-alphabetical order
	_, err := svc.CreateRole(ctx, "zrole", "", []Permission{{Resource: "a", Action: "a"}})
	require.NoError(t, err)
	_, err = svc.CreateRole(ctx, "arole", "", []Permission{{Resource: "b", Action: "b"}})
	require.NoError(t, err)
	_, err = svc.CreateRole(ctx, "mrole", "", []Permission{{Resource: "c", Action: "c"}})
	require.NoError(t, err)

	roles, err := svc.ListRoles(ctx)
	require.NoError(t, err)
	assert.Len(t, roles, 3)

	// Verify sorted
	assert.Equal(t, "arole", roles[0].Name)
	assert.Equal(t, "mrole", roles[1].Name)
	assert.Equal(t, "zrole", roles[2].Name)
}

func TestService_ListRoles_WithUserCounts(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	_, err := svc.CreateRole(ctx, "counttest", "", []Permission{{Resource: "a", Action: "a"}})
	require.NoError(t, err)

	u1 := uuid.Must(uuid.NewV7())
	u2 := uuid.Must(uuid.NewV7())
	require.NoError(t, svc.AssignRole(ctx, u1, "counttest"))
	require.NoError(t, svc.AssignRole(ctx, u2, "counttest"))

	roles, err := svc.ListRoles(ctx)
	require.NoError(t, err)

	var foundRole *Role
	for i := range roles {
		if roles[i].Name == "counttest" {
			foundRole = &roles[i]
			break
		}
	}
	require.NotNil(t, foundRole)
	assert.Equal(t, 2, foundRole.UserCount)
}

func TestService_GetRole_BuiltInAndCustom(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	// Create a custom role
	_, err := svc.CreateRole(ctx, "customrole", "", []Permission{
		{Resource: "data", Action: "read"},
	})
	require.NoError(t, err)

	role, err := svc.GetRole(ctx, "customrole")
	require.NoError(t, err)
	assert.Equal(t, "customrole", role.Name)
	assert.False(t, role.IsBuiltIn)
	assert.Equal(t, 0, role.UserCount)
}

func TestService_DeleteRole_AllBuiltInRolesProtected(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	for name := range BuiltInRoles {
		err := svc.DeleteRole(ctx, name)
		assert.ErrorIs(t, err, ErrBuiltInRole, "built-in role %q should be protected", name)
	}
}

func TestService_UpdateRolePermissions_BuiltInRole(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	// Create a "admin" policy manually (to simulate built-in role existing)
	require.NoError(t, svc.AddPolicy(ctx, "admin", "all", "*"))

	newPerms := []Permission{
		{Resource: "users", Action: "list"},
	}

	role, err := svc.UpdateRolePermissions(ctx, "admin", newPerms)
	require.NoError(t, err)
	assert.True(t, role.IsBuiltIn)
	assert.Equal(t, "Full system access", role.Description)
}

func TestService_ListPermissions_CorrectCount(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	perms := svc.ListPermissions(ctx)

	expectedCount := len(FineGrainedResources) * len(FineGrainedActions)
	assert.Len(t, perms, expectedCount)

	// Verify structure
	for _, p := range perms {
		assert.NotEmpty(t, p.Resource)
		assert.NotEmpty(t, p.Action)
	}
}

func TestService_GetAllRoleNames_EmptyAndSorted(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	// Empty
	names, err := svc.GetAllRoleNames(ctx)
	require.NoError(t, err)
	assert.Empty(t, names)

	// Add roles
	_, err = svc.CreateRole(ctx, "zeta", "", []Permission{{Resource: "z", Action: "z"}})
	require.NoError(t, err)
	_, err = svc.CreateRole(ctx, "alpha", "", []Permission{{Resource: "a", Action: "a"}})
	require.NoError(t, err)

	names, err = svc.GetAllRoleNames(ctx)
	require.NoError(t, err)
	assert.Equal(t, []string{"alpha", "zeta"}, names)
}

func TestParsePermission_EdgeCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    Permission
		wantErr bool
	}{
		{
			name:  "wildcard action",
			input: "admin:*",
			want:  Permission{Resource: "admin", Action: "*"},
		},
		{
			name:  "colon at start",
			input: ":read",
			want:  Permission{Resource: "", Action: "read"},
		},
		{
			name:  "multiple colons",
			input: "resource:sub:action",
			want:  Permission{Resource: "resource", Action: "sub:action"},
		},
		{
			name:    "only colon",
			input:   ":",
			want:    Permission{Resource: "", Action: ""},
			wantErr: false, // SplitN with ":" produces ["", ""]
		},
		{
			name:    "no colon",
			input:   "nocolon",
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

func TestPermission_String_Various(t *testing.T) {
	t.Parallel()

	tests := []struct {
		perm     Permission
		expected string
	}{
		{Permission{Resource: "users", Action: "list"}, "users:list"},
		{Permission{Resource: "admin", Action: "*"}, "admin:*"},
		{Permission{Resource: "", Action: ""}, ":"},
		{Permission{Resource: "movies", Action: "delete"}, "movies:delete"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.perm.String())
		})
	}
}

func TestBuiltInRoles_Contents(t *testing.T) {
	t.Parallel()

	assert.Len(t, BuiltInRoles, 4)
	assert.Equal(t, "Full system access", BuiltInRoles["admin"])
	assert.Equal(t, "Content moderation and user management", BuiltInRoles["moderator"])
	assert.Equal(t, "Standard user access", BuiltInRoles["user"])
	assert.Equal(t, "Read-only access", BuiltInRoles["guest"])
}

func TestAvailableResources_SameAsFineGrained(t *testing.T) {
	t.Parallel()

	assert.Equal(t, FineGrainedResources, AvailableResources)
}

func TestAvailableActions_SameAsFineGrained(t *testing.T) {
	t.Parallel()

	assert.Equal(t, FineGrainedActions, AvailableActions)
}

func TestDefaultRolePermissions_UserCannotModify(t *testing.T) {
	t.Parallel()

	userPerms := DefaultRolePermissions["user"]

	// User should NOT have admin, delete, or settings:write
	assert.NotContains(t, userPerms, PermAdminAll)
	assert.NotContains(t, userPerms, PermUsersDelete)
	assert.NotContains(t, userPerms, PermSettingsWrite)
	assert.NotContains(t, userPerms, PermAuditExport)
}

func TestDefaultRolePermissions_GuestMinimal(t *testing.T) {
	t.Parallel()

	guestPerms := DefaultRolePermissions["guest"]

	// Guest should have very limited permissions
	assert.Contains(t, guestPerms, PermProfileRead)
	assert.Contains(t, guestPerms, PermMoviesList)
	assert.Contains(t, guestPerms, PermMoviesGet)
	assert.Contains(t, guestPerms, PermLibrariesList)
	assert.Contains(t, guestPerms, PermLibrariesGet)
	assert.Contains(t, guestPerms, PermPlaybackStream)

	// Guest should NOT have
	assert.NotContains(t, guestPerms, PermProfileUpdate)
	assert.NotContains(t, guestPerms, PermPlaybackProgress)
	assert.NotContains(t, guestPerms, PermRequestsCreate)
}

func TestHasPermission_EdgeCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		permissions []string
		perm        string
		expected    bool
	}{
		{
			name:        "single matching permission",
			permissions: []string{"movies:list"},
			perm:        "movies:list",
			expected:    true,
		},
		{
			name:        "single non-matching permission",
			permissions: []string{"movies:list"},
			perm:        "movies:delete",
			expected:    false,
		},
		{
			name:        "admin wildcard only",
			permissions: []string{PermAdminAll},
			perm:        "anything:anything",
			expected:    true,
		},
		{
			name:        "admin wildcard with empty string search",
			permissions: []string{PermAdminAll},
			perm:        "",
			expected:    true,
		},
		{
			name:        "empty search against non-empty list",
			permissions: []string{"movies:list", "movies:get"},
			perm:        "",
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasPermission(tt.permissions, tt.perm)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// =====================================================
// GetUsersForRole with invalid UUID in role mapping
// =====================================================

func TestService_GetUsersForRole_InvalidUUIDSkipped(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	// Directly add a non-UUID user to a role via the enforcer
	_, err := svc.enforcer.AddRoleForUser("not-a-uuid", "testrole")
	require.NoError(t, err)

	// Also add a valid UUID user
	validID := uuid.Must(uuid.NewV7())
	_, err = svc.enforcer.AddRoleForUser(validID.String(), "testrole")
	require.NoError(t, err)

	// GetUsersForRole should skip the invalid UUID and return only the valid one
	users, err := svc.GetUsersForRole(ctx, "testrole")
	require.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, validID, users[0])
}

// =====================================================
// DeleteRole additional branches
// =====================================================

func TestService_DeleteRole_NonexistentReturnsNotFound(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	err := svc.DeleteRole(ctx, "nonexistentrole")
	assert.ErrorIs(t, err, ErrRoleNotFound)
}

func TestService_DeleteRole_InUseReturnsError(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	// Create a custom role
	_, err := svc.CreateRole(ctx, "inuserole", "", []Permission{
		{Resource: "movies", Action: "list"},
	})
	require.NoError(t, err)

	// Assign to a user
	userID := uuid.Must(uuid.NewV7())
	require.NoError(t, svc.AssignRole(ctx, userID, "inuserole"))

	// Try to delete - should fail because in use
	err = svc.DeleteRole(ctx, "inuserole")
	assert.ErrorIs(t, err, ErrRoleInUse)
}

func TestService_DeleteRole_Success(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	// Create a custom role
	_, err := svc.CreateRole(ctx, "deleterole", "", []Permission{
		{Resource: "data", Action: "read"},
	})
	require.NoError(t, err)

	// Delete should succeed (no users assigned)
	err = svc.DeleteRole(ctx, "deleterole")
	require.NoError(t, err)

	// Verify it's gone
	_, err = svc.GetRole(ctx, "deleterole")
	assert.ErrorIs(t, err, ErrRoleNotFound)
}

// =====================================================
// GetRole additional branches
// =====================================================

func TestService_GetRole_NonexistentReturnsNotFound(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	role, err := svc.GetRole(ctx, "nonexistentrole")
	assert.ErrorIs(t, err, ErrRoleNotFound)
	assert.Nil(t, role)
}

func TestService_GetRole_BuiltInWithPolicies(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	// Add policies for the "admin" built-in role
	require.NoError(t, svc.AddPolicy(ctx, "admin", "all", "*"))
	require.NoError(t, svc.AddPolicy(ctx, "admin", "users", "manage"))

	role, err := svc.GetRole(ctx, "admin")
	require.NoError(t, err)
	assert.Equal(t, "admin", role.Name)
	assert.True(t, role.IsBuiltIn)
	assert.Equal(t, "Full system access", role.Description)
	assert.Len(t, role.Permissions, 2)
}

// =====================================================
// CreateRole additional branches
// =====================================================

func TestService_CreateRole_EmptyNameReturnsError(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	_, err := svc.CreateRole(ctx, "", "", []Permission{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "empty")
}

func TestService_CreateRole_DuplicateReturnsAlreadyExists(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	_, err := svc.CreateRole(ctx, "existingrole", "", []Permission{
		{Resource: "a", Action: "b"},
	})
	require.NoError(t, err)

	// Try to create again
	_, err = svc.CreateRole(ctx, "existingrole", "", []Permission{
		{Resource: "c", Action: "d"},
	})
	assert.ErrorIs(t, err, ErrRoleAlreadyExists)
}

func TestService_CreateRole_WithSpaceInName(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	_, err := svc.CreateRole(ctx, "role with space", "", []Permission{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "whitespace")
}

// =====================================================
// GetRolePermissions tests
// =====================================================

func TestService_GetRolePermissions_EmptyAndPopulated(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	// Empty - no policies for this role
	perms, err := svc.GetRolePermissions(ctx, "emptyrole")
	require.NoError(t, err)
	assert.Empty(t, perms)

	// Populated
	_, err = svc.CreateRole(ctx, "permrole", "", []Permission{
		{Resource: "movies", Action: "list"},
		{Resource: "movies", Action: "get"},
	})
	require.NoError(t, err)

	perms, err = svc.GetRolePermissions(ctx, "permrole")
	require.NoError(t, err)
	assert.Len(t, perms, 2)
}

// =====================================================
// AddPermissionToRole / RemovePermissionFromRole tests
// =====================================================

func TestService_AddPermissionToRole_Success(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	// Create role first
	_, err := svc.CreateRole(ctx, "addrole", "", []Permission{
		{Resource: "movies", Action: "list"},
	})
	require.NoError(t, err)

	// Add another permission
	err = svc.AddPermissionToRole(ctx, "addrole", Permission{Resource: "movies", Action: "get"})
	require.NoError(t, err)

	// Verify
	perms, err := svc.GetRolePermissions(ctx, "addrole")
	require.NoError(t, err)
	assert.Len(t, perms, 2)
}

func TestService_AddPermissionToRole_ExistingReturnsDuplicate(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	_, err := svc.CreateRole(ctx, "duprole", "", []Permission{
		{Resource: "movies", Action: "list"},
	})
	require.NoError(t, err)

	// Try to add same permission
	err = svc.AddPermissionToRole(ctx, "duprole", Permission{Resource: "movies", Action: "list"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestService_RemovePermissionFromRole_Success(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	_, err := svc.CreateRole(ctx, "remrole", "", []Permission{
		{Resource: "movies", Action: "list"},
		{Resource: "movies", Action: "get"},
	})
	require.NoError(t, err)

	// Remove one permission
	err = svc.RemovePermissionFromRole(ctx, "remrole", Permission{Resource: "movies", Action: "list"})
	require.NoError(t, err)

	// Verify only one remains
	perms, err := svc.GetRolePermissions(ctx, "remrole")
	require.NoError(t, err)
	assert.Len(t, perms, 1)
	assert.Equal(t, "get", perms[0].Action)
}

func TestService_RemovePermissionFromRole_MissingReturnsNotFound(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	_, err := svc.CreateRole(ctx, "remrole2", "", []Permission{
		{Resource: "movies", Action: "list"},
	})
	require.NoError(t, err)

	// Try to remove non-existent permission
	err = svc.RemovePermissionFromRole(ctx, "remrole2", Permission{Resource: "movies", Action: "delete"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// =====================================================
// UpdateRolePermissions additional branches
// =====================================================

func TestService_UpdateRolePermissions_NonexistentReturnsNotFound(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	_, err := svc.UpdateRolePermissions(ctx, "nonexistentrole", []Permission{
		{Resource: "a", Action: "b"},
	})
	assert.ErrorIs(t, err, ErrRoleNotFound)
}

func TestService_UpdateRolePermissions_CustomRole(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	// Create custom role
	_, err := svc.CreateRole(ctx, "updaterole", "", []Permission{
		{Resource: "movies", Action: "list"},
	})
	require.NoError(t, err)

	// Update permissions
	newPerms := []Permission{
		{Resource: "movies", Action: "get"},
		{Resource: "users", Action: "list"},
	}

	role, err := svc.UpdateRolePermissions(ctx, "updaterole", newPerms)
	require.NoError(t, err)
	assert.Equal(t, "updaterole", role.Name)
	assert.False(t, role.IsBuiltIn)
	assert.Len(t, role.Permissions, 2)
}

// =====================================================
// Error variable tests
// =====================================================

func TestErrorVariables_Distinct(t *testing.T) {
	t.Parallel()

	assert.NotEqual(t, ErrRoleNotFound, ErrRoleAlreadyExists)
	assert.NotEqual(t, ErrRoleNotFound, ErrBuiltInRole)
	assert.NotEqual(t, ErrRoleNotFound, ErrRoleInUse)
	assert.NotEqual(t, ErrRoleAlreadyExists, ErrBuiltInRole)
	assert.NotEqual(t, ErrRoleAlreadyExists, ErrRoleInUse)
	assert.NotEqual(t, ErrBuiltInRole, ErrRoleInUse)

	assert.Equal(t, "role not found", ErrRoleNotFound.Error())
	assert.Equal(t, "role already exists", ErrRoleAlreadyExists.Error())
	assert.Equal(t, "cannot modify built-in role", ErrBuiltInRole.Error())
	assert.Equal(t, "role is assigned to users", ErrRoleInUse.Error())
}

// =====================================================
// ListRoles with built-in and custom roles mix
// =====================================================

func TestService_ListRoles_MixedBuiltInAndCustom(t *testing.T) {
	t.Parallel()

	svc := setupInternalTestService(t)
	ctx := context.Background()

	// Add built-in role policies
	require.NoError(t, svc.AddPolicy(ctx, "admin", "all", "*"))

	// Add custom role
	_, err := svc.CreateRole(ctx, "customrole", "", []Permission{
		{Resource: "movies", Action: "list"},
	})
	require.NoError(t, err)

	roles, err := svc.ListRoles(ctx)
	require.NoError(t, err)
	assert.Len(t, roles, 2)

	// Find admin and custom role
	var adminRole, customRole *Role
	for i := range roles {
		if roles[i].Name == "admin" {
			adminRole = &roles[i]
		} else if roles[i].Name == "customrole" {
			customRole = &roles[i]
		}
	}

	require.NotNil(t, adminRole)
	assert.True(t, adminRole.IsBuiltIn)
	assert.Equal(t, "Full system access", adminRole.Description)

	require.NotNil(t, customRole)
	assert.False(t, customRole.IsBuiltIn)
	assert.Empty(t, customRole.Description)
}
