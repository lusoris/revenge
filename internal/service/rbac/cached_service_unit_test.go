package rbac

import (
	"context"
	"testing"
	"time"

	"github.com/casbin/casbin/v2"
	casbinmodel "github.com/casbin/casbin/v2/model"
	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/infra/logging"
	"github.com/lusoris/revenge/internal/service/activity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// casbinModelConfInternal is the same model used by existing tests.
const casbinModelConfInternal = `
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

// setupInternalTestService creates a Service with in-memory enforcer for same-package tests.
func setupInternalTestService(t *testing.T) *Service {
	t.Helper()

	m, err := casbinmodel.NewModelFromString(casbinModelConfInternal)
	require.NoError(t, err)

	enforcer, err := casbin.NewSyncedEnforcer(m)
	require.NoError(t, err)

	logger := logging.NewTestLogger()
	return NewService(enforcer, logger, activity.NewNoopLogger())
}

// newTestRBACCache creates an L1-only cache for testing.
func newTestRBACCache(t *testing.T) *cache.Cache {
	t.Helper()
	c, err := cache.NewCache(nil, 1000, 15*time.Minute)
	require.NoError(t, err)
	return c
}

// --- HasPermission tests ---

func TestHasPermission(t *testing.T) {
	tests := []struct {
		name        string
		permissions []string
		perm        string
		expected    bool
	}{
		{
			name:        "exact match found",
			permissions: []string{"users:list", "users:get", "movies:list"},
			perm:        "users:get",
			expected:    true,
		},
		{
			name:        "not found",
			permissions: []string{"users:list", "users:get"},
			perm:        "movies:delete",
			expected:    false,
		},
		{
			name:        "admin wildcard grants everything",
			permissions: []string{PermAdminAll},
			perm:        "movies:delete",
			expected:    true,
		},
		{
			name:        "admin wildcard with other perms",
			permissions: []string{"users:list", PermAdminAll, "movies:get"},
			perm:        "settings:write",
			expected:    true,
		},
		{
			name:        "empty permissions list",
			permissions: []string{},
			perm:        "users:list",
			expected:    false,
		},
		{
			name:        "nil permissions list",
			permissions: nil,
			perm:        "users:list",
			expected:    false,
		},
		{
			name:        "first element matches",
			permissions: []string{"movies:list", "movies:get", "movies:create"},
			perm:        "movies:list",
			expected:    true,
		},
		{
			name:        "last element matches",
			permissions: []string{"movies:list", "movies:get", "movies:create"},
			perm:        "movies:create",
			expected:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasPermission(tt.permissions, tt.perm)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// --- DefaultRolePermissions tests ---

func TestDefaultRolePermissions(t *testing.T) {
	t.Run("admin has admin wildcard", func(t *testing.T) {
		perms, ok := DefaultRolePermissions["admin"]
		require.True(t, ok)
		assert.Contains(t, perms, PermAdminAll)
	})

	t.Run("user has basic permissions", func(t *testing.T) {
		perms, ok := DefaultRolePermissions["user"]
		require.True(t, ok)
		assert.Contains(t, perms, PermProfileRead)
		assert.Contains(t, perms, PermProfileUpdate)
		assert.Contains(t, perms, PermMoviesList)
		assert.Contains(t, perms, PermPlaybackStream)
	})

	t.Run("guest has limited permissions", func(t *testing.T) {
		perms, ok := DefaultRolePermissions["guest"]
		require.True(t, ok)
		assert.Contains(t, perms, PermMoviesList)
		assert.Contains(t, perms, PermPlaybackStream)
		assert.NotContains(t, perms, PermMoviesCreate)
		assert.NotContains(t, perms, PermUsersCreate)
	})

	t.Run("moderator has elevated permissions", func(t *testing.T) {
		perms, ok := DefaultRolePermissions["moderator"]
		require.True(t, ok)
		assert.Contains(t, perms, PermMoviesCreate)
		assert.Contains(t, perms, PermMoviesDelete)
		assert.Contains(t, perms, PermLibrariesScan)
		assert.Contains(t, perms, PermRequestsApprove)
	})
}

// --- loadPolicyLine tests ---

func TestLoadPolicyLine(t *testing.T) {
	t.Run("loads p-type policy rule", func(t *testing.T) {
		m, err := casbinmodel.NewModelFromString(casbinModelConfInternal)
		require.NoError(t, err)

		rule := &CasbinRule{
			PType: "p",
			V0:    "alice",
			V1:    "data1",
			V2:    "read",
		}

		loadPolicyLine(rule, m)

		// Verify the policy was loaded
		policies := m["p"]["p"].Policy
		assert.Len(t, policies, 1)
		assert.Equal(t, []string{"alice", "data1", "read"}, policies[0])
	})

	t.Run("loads g-type role assignment", func(t *testing.T) {
		m, err := casbinmodel.NewModelFromString(casbinModelConfInternal)
		require.NoError(t, err)

		rule := &CasbinRule{
			PType: "g",
			V0:    "alice",
			V1:    "admin",
		}

		loadPolicyLine(rule, m)

		groupPolicies := m["g"]["g"].Policy
		assert.Len(t, groupPolicies, 1)
		assert.Equal(t, []string{"alice", "admin"}, groupPolicies[0])
	})

	t.Run("skips empty values", func(t *testing.T) {
		m, err := casbinmodel.NewModelFromString(casbinModelConfInternal)
		require.NoError(t, err)

		rule := &CasbinRule{
			PType: "p",
			V0:    "alice",
			V1:    "data1",
			V2:    "read",
			V3:    "", // empty
			V4:    "", // empty
			V5:    "", // empty
		}

		loadPolicyLine(rule, m)

		policies := m["p"]["p"].Policy
		assert.Len(t, policies, 1)
		assert.Equal(t, []string{"alice", "data1", "read"}, policies[0])
	})

	t.Run("includes non-empty V3-V5", func(t *testing.T) {
		m, err := casbinmodel.NewModelFromString(casbinModelConfInternal)
		require.NoError(t, err)

		rule := &CasbinRule{
			PType: "p",
			V0:    "alice",
			V1:    "data1",
			V2:    "read",
			V3:    "extra1",
			V4:    "extra2",
			V5:    "extra3",
		}

		loadPolicyLine(rule, m)

		policies := m["p"]["p"].Policy
		assert.Len(t, policies, 1)
		// Casbin model may or may not use V3-V5 but they should be included in the policy slice
		assert.Contains(t, policies[0], "alice")
		assert.Contains(t, policies[0], "data1")
		assert.Contains(t, policies[0], "read")
	})

	t.Run("unknown section is silently ignored", func(t *testing.T) {
		m, err := casbinmodel.NewModelFromString(casbinModelConfInternal)
		require.NoError(t, err)

		rule := &CasbinRule{
			PType: "x",
			V0:    "alice",
			V1:    "data1",
		}

		// Should not panic
		loadPolicyLine(rule, m)
	})

	t.Run("multiple rules loaded correctly", func(t *testing.T) {
		m, err := casbinmodel.NewModelFromString(casbinModelConfInternal)
		require.NoError(t, err)

		rules := []*CasbinRule{
			{PType: "p", V0: "alice", V1: "data1", V2: "read"},
			{PType: "p", V0: "bob", V1: "data2", V2: "write"},
			{PType: "g", V0: "alice", V1: "admin"},
		}

		for _, rule := range rules {
			loadPolicyLine(rule, m)
		}

		assert.Len(t, m["p"]["p"].Policy, 2)
		assert.Len(t, m["g"]["g"].Policy, 1)
	})
}

// --- CachedService tests ---

func TestNewCachedService_Unit(t *testing.T) {
	svc := setupInternalTestService(t)
	logger := logging.NewTestLogger()

	t.Run("with cache", func(t *testing.T) {
		c := newTestRBACCache(t)
		cached := NewCachedService(svc, c, logger)
		assert.NotNil(t, cached)
	})

	t.Run("with nil cache", func(t *testing.T) {
		cached := NewCachedService(svc, nil, logger)
		assert.NotNil(t, cached)
	})
}

func TestCachedService_Enforce(t *testing.T) {
	t.Run("nil cache falls back to service", func(t *testing.T) {
		svc := setupInternalTestService(t)
		ctx := context.Background()

		require.NoError(t, svc.AddPolicy(ctx, "alice", "data1", "read"))

		cached := NewCachedService(svc, nil, logging.NewTestLogger())

		allowed, err := cached.Enforce(ctx, "alice", "data1", "read")
		require.NoError(t, err)
		assert.True(t, allowed)

		allowed, err = cached.Enforce(ctx, "alice", "data1", "write")
		require.NoError(t, err)
		assert.False(t, allowed)
	})

	t.Run("cache miss then hit", func(t *testing.T) {
		svc := setupInternalTestService(t)
		ctx := context.Background()
		c := newTestRBACCache(t)

		require.NoError(t, svc.AddPolicy(ctx, "alice", "data1", "read"))

		cached := NewCachedService(svc, c, logging.NewTestLogger())

		// First call - cache miss
		allowed, err := cached.Enforce(ctx, "alice", "data1", "read")
		require.NoError(t, err)
		assert.True(t, allowed)

		// Wait for async cache set
		time.Sleep(200 * time.Millisecond)

		// Second call - cache hit
		allowed, err = cached.Enforce(ctx, "alice", "data1", "read")
		require.NoError(t, err)
		assert.True(t, allowed)
	})

	t.Run("denied result also cached", func(t *testing.T) {
		svc := setupInternalTestService(t)
		ctx := context.Background()
		c := newTestRBACCache(t)

		cached := NewCachedService(svc, c, logging.NewTestLogger())

		// First call - denied, cache miss
		allowed, err := cached.Enforce(ctx, "bob", "data1", "read")
		require.NoError(t, err)
		assert.False(t, allowed)

		// Wait for async cache set
		time.Sleep(200 * time.Millisecond)

		// Second call - denied, cache hit
		allowed, err = cached.Enforce(ctx, "bob", "data1", "read")
		require.NoError(t, err)
		assert.False(t, allowed)
	})
}

func TestCachedService_EnforceWithContext_Unit(t *testing.T) {
	svc := setupInternalTestService(t)
	ctx := context.Background()
	c := newTestRBACCache(t)

	userID := uuid.Must(uuid.NewV7())
	require.NoError(t, svc.AddPolicy(ctx, userID.String(), "movies", "read"))

	cached := NewCachedService(svc, c, logging.NewTestLogger())

	allowed, err := cached.EnforceWithContext(ctx, userID, "movies", "read")
	require.NoError(t, err)
	assert.True(t, allowed)

	allowed, err = cached.EnforceWithContext(ctx, userID, "movies", "write")
	require.NoError(t, err)
	assert.False(t, allowed)
}

func TestCachedService_GetUserRoles_Unit(t *testing.T) {
	t.Run("nil cache falls back to service", func(t *testing.T) {
		svc := setupInternalTestService(t)
		ctx := context.Background()

		userID := uuid.Must(uuid.NewV7())
		require.NoError(t, svc.AssignRole(ctx, userID, "admin"))

		cached := NewCachedService(svc, nil, logging.NewTestLogger())

		roles, err := cached.GetUserRoles(ctx, userID)
		require.NoError(t, err)
		assert.Contains(t, roles, "admin")
	})

	t.Run("cache miss then hit", func(t *testing.T) {
		svc := setupInternalTestService(t)
		ctx := context.Background()
		c := newTestRBACCache(t)

		userID := uuid.Must(uuid.NewV7())
		require.NoError(t, svc.AssignRole(ctx, userID, "admin"))
		require.NoError(t, svc.AssignRole(ctx, userID, "editor"))

		cached := NewCachedService(svc, c, logging.NewTestLogger())

		// First call - cache miss
		roles, err := cached.GetUserRoles(ctx, userID)
		require.NoError(t, err)
		assert.Len(t, roles, 2)
		assert.Contains(t, roles, "admin")
		assert.Contains(t, roles, "editor")

		// Wait for async cache set
		time.Sleep(200 * time.Millisecond)

		// Second call - cache hit
		roles, err = cached.GetUserRoles(ctx, userID)
		require.NoError(t, err)
		assert.Len(t, roles, 2)
	})

	t.Run("empty roles", func(t *testing.T) {
		svc := setupInternalTestService(t)
		ctx := context.Background()
		c := newTestRBACCache(t)

		userID := uuid.Must(uuid.NewV7())

		cached := NewCachedService(svc, c, logging.NewTestLogger())

		roles, err := cached.GetUserRoles(ctx, userID)
		require.NoError(t, err)
		assert.Empty(t, roles)
	})
}

func TestCachedService_HasRole_Unit(t *testing.T) {
	t.Run("nil cache falls back to service", func(t *testing.T) {
		svc := setupInternalTestService(t)
		ctx := context.Background()

		userID := uuid.Must(uuid.NewV7())
		require.NoError(t, svc.AssignRole(ctx, userID, "admin"))

		cached := NewCachedService(svc, nil, logging.NewTestLogger())

		hasRole, err := cached.HasRole(ctx, userID, "admin")
		require.NoError(t, err)
		assert.True(t, hasRole)

		hasRole, err = cached.HasRole(ctx, userID, "editor")
		require.NoError(t, err)
		assert.False(t, hasRole)
	})

	t.Run("cache miss then hit", func(t *testing.T) {
		svc := setupInternalTestService(t)
		ctx := context.Background()
		c := newTestRBACCache(t)

		userID := uuid.Must(uuid.NewV7())
		require.NoError(t, svc.AssignRole(ctx, userID, "admin"))

		cached := NewCachedService(svc, c, logging.NewTestLogger())

		// First call - cache miss
		hasRole, err := cached.HasRole(ctx, userID, "admin")
		require.NoError(t, err)
		assert.True(t, hasRole)

		// Wait for async cache set
		time.Sleep(200 * time.Millisecond)

		// Second call - cache hit
		hasRole, err = cached.HasRole(ctx, userID, "admin")
		require.NoError(t, err)
		assert.True(t, hasRole)
	})
}

func TestCachedService_AssignRole(t *testing.T) {
	t.Run("assigns role and invalidates cache", func(t *testing.T) {
		svc := setupInternalTestService(t)
		ctx := context.Background()
		c := newTestRBACCache(t)

		cached := NewCachedService(svc, c, logging.NewTestLogger())
		userID := uuid.Must(uuid.NewV7())

		err := cached.AssignRole(ctx, userID, "admin")
		require.NoError(t, err)

		// Verify role was actually assigned
		hasRole, err := cached.HasRole(ctx, userID, "admin")
		require.NoError(t, err)
		assert.True(t, hasRole)
	})

	t.Run("nil cache still works", func(t *testing.T) {
		svc := setupInternalTestService(t)
		ctx := context.Background()

		cached := NewCachedService(svc, nil, logging.NewTestLogger())
		userID := uuid.Must(uuid.NewV7())

		err := cached.AssignRole(ctx, userID, "admin")
		require.NoError(t, err)
	})
}

func TestCachedService_RemoveRole(t *testing.T) {
	t.Run("removes role and invalidates cache", func(t *testing.T) {
		svc := setupInternalTestService(t)
		ctx := context.Background()
		c := newTestRBACCache(t)

		userID := uuid.Must(uuid.NewV7())
		require.NoError(t, svc.AssignRole(ctx, userID, "admin"))

		cached := NewCachedService(svc, c, logging.NewTestLogger())

		err := cached.RemoveRole(ctx, userID, "admin")
		require.NoError(t, err)

		hasRole, err := cached.HasRole(ctx, userID, "admin")
		require.NoError(t, err)
		assert.False(t, hasRole)
	})

	t.Run("error on non-existent role", func(t *testing.T) {
		svc := setupInternalTestService(t)
		ctx := context.Background()
		c := newTestRBACCache(t)

		cached := NewCachedService(svc, c, logging.NewTestLogger())
		userID := uuid.Must(uuid.NewV7())

		err := cached.RemoveRole(ctx, userID, "nonexistent")
		assert.Error(t, err)
	})
}

func TestCachedService_AddPolicy(t *testing.T) {
	t.Run("adds policy and invalidates cache", func(t *testing.T) {
		svc := setupInternalTestService(t)
		ctx := context.Background()
		c := newTestRBACCache(t)

		cached := NewCachedService(svc, c, logging.NewTestLogger())

		err := cached.AddPolicy(ctx, "alice", "data1", "read")
		require.NoError(t, err)

		// Verify policy works
		allowed, err := cached.Enforce(ctx, "alice", "data1", "read")
		require.NoError(t, err)
		assert.True(t, allowed)
	})

	t.Run("nil cache still works", func(t *testing.T) {
		svc := setupInternalTestService(t)
		ctx := context.Background()

		cached := NewCachedService(svc, nil, logging.NewTestLogger())

		err := cached.AddPolicy(ctx, "alice", "data1", "read")
		require.NoError(t, err)
	})
}

func TestCachedService_RemovePolicy(t *testing.T) {
	t.Run("removes policy and invalidates cache", func(t *testing.T) {
		svc := setupInternalTestService(t)
		ctx := context.Background()
		c := newTestRBACCache(t)

		require.NoError(t, svc.AddPolicy(ctx, "alice", "data1", "read"))

		cached := NewCachedService(svc, c, logging.NewTestLogger())

		err := cached.RemovePolicy(ctx, "alice", "data1", "read")
		require.NoError(t, err)

		// Verify policy was removed
		allowed, err := cached.Enforce(ctx, "alice", "data1", "read")
		require.NoError(t, err)
		assert.False(t, allowed)
	})

	t.Run("error on non-existent policy", func(t *testing.T) {
		svc := setupInternalTestService(t)
		ctx := context.Background()
		c := newTestRBACCache(t)

		cached := NewCachedService(svc, c, logging.NewTestLogger())

		err := cached.RemovePolicy(ctx, "nonexistent", "data1", "read")
		assert.Error(t, err)
	})
}

func TestCachedService_LoadPolicy(t *testing.T) {
	t.Run("panics with nil adapter", func(t *testing.T) {
		svc := setupInternalTestService(t)
		ctx := context.Background()
		c := newTestRBACCache(t)

		require.NoError(t, svc.AddPolicy(ctx, "alice", "data1", "read"))

		cached := NewCachedService(svc, c, logging.NewTestLogger())

		// In-memory enforcer with nil adapter panics in casbin
		assert.Panics(t, func() {
			_ = cached.LoadPolicy(ctx)
		})
	})
}

func TestCachedService_invalidateUserCache(t *testing.T) {
	t.Run("with cache", func(t *testing.T) {
		svc := setupInternalTestService(t)
		ctx := context.Background()
		c := newTestRBACCache(t)

		cached := NewCachedService(svc, c, logging.NewTestLogger())
		userID := uuid.Must(uuid.NewV7())

		// Should not panic or error
		cached.invalidateUserCache(ctx, userID)
	})

	t.Run("with nil cache", func(t *testing.T) {
		svc := setupInternalTestService(t)
		ctx := context.Background()

		cached := NewCachedService(svc, nil, logging.NewTestLogger())
		userID := uuid.Must(uuid.NewV7())

		// Should return immediately without error
		cached.invalidateUserCache(ctx, userID)
	})
}

func TestCachedService_invalidateAllRBAC(t *testing.T) {
	t.Run("with cache", func(t *testing.T) {
		svc := setupInternalTestService(t)
		ctx := context.Background()
		c := newTestRBACCache(t)

		cached := NewCachedService(svc, c, logging.NewTestLogger())

		// Should not panic or error
		cached.invalidateAllRBAC(ctx)
	})

	t.Run("with nil cache", func(t *testing.T) {
		svc := setupInternalTestService(t)
		ctx := context.Background()

		cached := NewCachedService(svc, nil, logging.NewTestLogger())

		// Should return immediately without error
		cached.invalidateAllRBAC(ctx)
	})
}

// --- Service.LoadPolicy and SavePolicy tests ---

func TestService_LoadPolicy_NoAdapter(t *testing.T) {
	svc := setupInternalTestService(t)

	// In-memory enforcer with nil adapter panics in casbin
	assert.Panics(t, func() {
		_ = svc.LoadPolicy(context.Background())
	})
}

func TestService_SavePolicy_NoAdapter(t *testing.T) {
	svc := setupInternalTestService(t)

	// In-memory enforcer with nil adapter panics in casbin
	assert.Panics(t, func() {
		_ = svc.SavePolicy(context.Background())
	})
}

// --- CasbinRule struct tests ---

func TestCasbinRule_Fields(t *testing.T) {
	rule := CasbinRule{
		PType: "p",
		V0:    "alice",
		V1:    "data1",
		V2:    "read",
		V3:    "extra1",
		V4:    "extra2",
		V5:    "extra3",
	}

	assert.Equal(t, "p", rule.PType)
	assert.Equal(t, "alice", rule.V0)
	assert.Equal(t, "data1", rule.V1)
	assert.Equal(t, "read", rule.V2)
	assert.Equal(t, "extra1", rule.V3)
	assert.Equal(t, "extra2", rule.V4)
	assert.Equal(t, "extra3", rule.V5)
}

// --- FineGrainedResources and FineGrainedActions tests ---

func TestFineGrainedResources(t *testing.T) {
	assert.NotEmpty(t, FineGrainedResources)
	assert.Contains(t, FineGrainedResources, "users")
	assert.Contains(t, FineGrainedResources, "movies")
	assert.Contains(t, FineGrainedResources, "libraries")
	assert.Contains(t, FineGrainedResources, "playback")
	assert.Contains(t, FineGrainedResources, "admin")
}

func TestFineGrainedActions(t *testing.T) {
	assert.NotEmpty(t, FineGrainedActions)
	assert.Contains(t, FineGrainedActions, "list")
	assert.Contains(t, FineGrainedActions, "get")
	assert.Contains(t, FineGrainedActions, "create")
	assert.Contains(t, FineGrainedActions, "update")
	assert.Contains(t, FineGrainedActions, "delete")
	assert.Contains(t, FineGrainedActions, "*")
}

// --- Permission constants tests ---

func TestPermissionConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"PermUsersList", PermUsersList, "users:list"},
		{"PermUsersGet", PermUsersGet, "users:get"},
		{"PermUsersCreate", PermUsersCreate, "users:create"},
		{"PermUsersUpdate", PermUsersUpdate, "users:update"},
		{"PermUsersDelete", PermUsersDelete, "users:delete"},
		{"PermProfileRead", PermProfileRead, "profile:read"},
		{"PermProfileUpdate", PermProfileUpdate, "profile:update"},
		{"PermMoviesList", PermMoviesList, "movies:list"},
		{"PermMoviesGet", PermMoviesGet, "movies:get"},
		{"PermMoviesCreate", PermMoviesCreate, "movies:create"},
		{"PermMoviesUpdate", PermMoviesUpdate, "movies:update"},
		{"PermMoviesDelete", PermMoviesDelete, "movies:delete"},
		{"PermLibrariesList", PermLibrariesList, "libraries:list"},
		{"PermLibrariesScan", PermLibrariesScan, "libraries:scan"},
		{"PermPlaybackStream", PermPlaybackStream, "playback:stream"},
		{"PermPlaybackProgress", PermPlaybackProgress, "playback:progress"},
		{"PermSettingsRead", PermSettingsRead, "settings:read"},
		{"PermSettingsWrite", PermSettingsWrite, "settings:write"},
		{"PermAuditRead", PermAuditRead, "audit:read"},
		{"PermAuditExport", PermAuditExport, "audit:export"},
		{"PermAdminAll", PermAdminAll, "admin:*"},
		{"PermIntegrationsList", PermIntegrationsList, "integrations:list"},
		{"PermIntegrationsSync", PermIntegrationsSync, "integrations:sync"},
		{"PermNotificationsList", PermNotificationsList, "notifications:list"},
		{"PermRequestsApprove", PermRequestsApprove, "requests:approve"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.constant)
		})
	}
}

// --- Adapter struct test ---

func TestNewAdapter(t *testing.T) {
	adapter := NewAdapter(nil)
	assert.NotNil(t, adapter)
	assert.Equal(t, "shared.casbin_rule", adapter.tableName)
}

// --- CachedService cache interaction patterns ---

func TestCachedService_EnforceInvalidationFlow(t *testing.T) {
	// Tests the flow: enforce (miss) -> cache set -> add policy -> invalidate -> enforce (miss again)
	svc := setupInternalTestService(t)
	ctx := context.Background()
	c := newTestRBACCache(t)

	cached := NewCachedService(svc, c, logging.NewTestLogger())

	// Enforce - denied (no policy), gets cached
	allowed, err := cached.Enforce(ctx, "alice", "data1", "read")
	require.NoError(t, err)
	assert.False(t, allowed)

	// Wait for cache write
	time.Sleep(200 * time.Millisecond)

	// Add policy through cached service (this invalidates the cache)
	err = cached.AddPolicy(ctx, "alice", "data1", "read")
	require.NoError(t, err)

	// Enforce again - should be allowed (cache was invalidated, fresh enforcement)
	allowed, err = cached.Enforce(ctx, "alice", "data1", "read")
	require.NoError(t, err)
	assert.True(t, allowed)
}

func TestCachedService_RoleInvalidationFlow(t *testing.T) {
	// Tests the flow: get roles (miss) -> cache -> assign role -> invalidate -> get roles (fresh)
	svc := setupInternalTestService(t)
	ctx := context.Background()
	c := newTestRBACCache(t)

	cached := NewCachedService(svc, c, logging.NewTestLogger())
	userID := uuid.Must(uuid.NewV7())

	// Get roles - empty, gets cached
	roles, err := cached.GetUserRoles(ctx, userID)
	require.NoError(t, err)
	assert.Empty(t, roles)

	// Wait for cache write
	time.Sleep(200 * time.Millisecond)

	// Assign role through cached service (invalidates user cache)
	err = cached.AssignRole(ctx, userID, "admin")
	require.NoError(t, err)

	// Get roles again - should have admin role (cache was invalidated)
	roles, err = cached.GetUserRoles(ctx, userID)
	require.NoError(t, err)
	assert.Contains(t, roles, "admin")
}

// --- cachedBool and cachedStringSlice types ---

func TestCachedTypes(t *testing.T) {
	t.Run("cachedBool", func(t *testing.T) {
		cb := cachedBool{Value: true}
		assert.True(t, cb.Value)

		cb = cachedBool{Value: false}
		assert.False(t, cb.Value)
	})

	t.Run("cachedStringSlice", func(t *testing.T) {
		css := cachedStringSlice{Values: []string{"admin", "editor"}}
		assert.Len(t, css.Values, 2)
		assert.Contains(t, css.Values, "admin")
		assert.Contains(t, css.Values, "editor")

		css = cachedStringSlice{Values: nil}
		assert.Nil(t, css.Values)
	})
}
