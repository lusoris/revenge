package rbac

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/testutil"
)

func setupCachedTestService(t *testing.T) (*CachedService, *testutil.TestDB) {
	t.Helper()
	svc, testDB := setupTestService(t)

	// Create cache instance (L1-only for tests)
	testCache, err := cache.NewCache(nil, 1000, 15*time.Minute)
	require.NoError(t, err)

	cachedSvc := NewCachedService(svc, testCache, zap.NewNop())
	return cachedSvc, testDB
}

func TestNewCachedService(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)

	testCache, err := cache.NewCache(nil, 1000, 15*time.Minute)
	require.NoError(t, err)

	cached := NewCachedService(svc, testCache, zap.NewNop())
	require.NotNil(t, cached)
	assert.NotNil(t, cached.Service)
	assert.NotNil(t, cached.cache)
}

func TestCachedService_Enforce_CacheHit(t *testing.T) {
	t.Parallel()
	svc, _ := setupCachedTestService(t)
	ctx := context.Background()

	// Add a policy
	err := svc.AddPolicy(ctx, "alice", "data1", "read")
	require.NoError(t, err)

	// First call - cache miss
	allowed, err := svc.Enforce(ctx, "alice", "data1", "read")
	require.NoError(t, err)
	assert.True(t, allowed)

	// Second call - should be cache hit
	allowed, err = svc.Enforce(ctx, "alice", "data1", "read")
	require.NoError(t, err)
	assert.True(t, allowed)
}

func TestCachedService_Enforce_NoCache(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Create cached service without cache (nil cache)
	cachedSvc := NewCachedService(svc, nil, zap.NewNop())

	// Add a policy
	err := svc.AddPolicy(ctx, "bob", "data2", "write")
	require.NoError(t, err)

	// Enforce should work without cache
	allowed, err := cachedSvc.Enforce(ctx, "bob", "data2", "write")
	require.NoError(t, err)
	assert.True(t, allowed)
}

func TestCachedService_EnforceWithContext(t *testing.T) {
	t.Parallel()
	svc, _ := setupCachedTestService(t)
	ctx := context.Background()

	userID := uuid.New()

	// Create role and assign to user
	_, err := svc.CreateRole(ctx, "reader", "", []Permission{{Resource: "books", Action: "read"}})
	require.NoError(t, err)

	err = svc.AssignRole(ctx, userID, "reader")
	require.NoError(t, err)

	// Check permission
	allowed, err := svc.EnforceWithContext(ctx, userID, "books", "read")
	require.NoError(t, err)
	assert.True(t, allowed)
}

func TestCachedService_GetUserRoles(t *testing.T) {
	t.Parallel()
	svc, _ := setupCachedTestService(t)
	ctx := context.Background()

	userID := uuid.New()

	// Create roles and assign
	_, err := svc.CreateRole(ctx, "admin", "", []Permission{{Resource: "users", Action: "manage"}})
	require.NoError(t, err)
	_, err = svc.CreateRole(ctx, "editor", "", []Permission{{Resource: "content", Action: "write"}})
	require.NoError(t, err)

	err = svc.AssignRole(ctx, userID, "admin")
	require.NoError(t, err)
	err = svc.AssignRole(ctx, userID, "editor")
	require.NoError(t, err)

	// First call - cache miss
	roles, err := svc.GetUserRoles(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, 2, len(roles))

	// Second call - cache hit
	roles, err = svc.GetUserRoles(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, 2, len(roles))
}

func TestCachedService_GetUserRoles_NoCache(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	cachedSvc := NewCachedService(svc, nil, zap.NewNop())
	userID := uuid.New()

	// Should work without cache
	roles, err := cachedSvc.GetUserRoles(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, 0, len(roles))
}

func TestCachedService_HasRole(t *testing.T) {
	t.Parallel()
	svc, _ := setupCachedTestService(t)
	ctx := context.Background()

	userID := uuid.New()

	// Create role and assign
	_, err := svc.CreateRole(ctx, "moderator", "", []Permission{{Resource: "posts", Action: "moderate"}})
	require.NoError(t, err)

	err = svc.AssignRole(ctx, userID, "moderator")
	require.NoError(t, err)

	// First call - cache miss
	hasRole, err := svc.HasRole(ctx, userID, "moderator")
	require.NoError(t, err)
	assert.True(t, hasRole)

	// Second call - cache hit
	hasRole, err = svc.HasRole(ctx, userID, "moderator")
	require.NoError(t, err)
	assert.True(t, hasRole)

	// Check non-existent role
	hasRole, err = svc.HasRole(ctx, userID, "admin")
	require.NoError(t, err)
	assert.False(t, hasRole)
}

func TestCachedService_HasRole_NoCache(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	cachedSvc := NewCachedService(svc, nil, zap.NewNop())
	userID := uuid.New()

	// Should work without cache
	hasRole, err := cachedSvc.HasRole(ctx, userID, "admin")
	require.NoError(t, err)
	assert.False(t, hasRole)
}

func TestCachedService_AssignRole_InvalidatesCache(t *testing.T) {
	t.Parallel()
	svc, _ := setupCachedTestService(t)
	ctx := context.Background()

	userID := uuid.New()

	// Create role
	_, err := svc.CreateRole(ctx, "writer", "", []Permission{{Resource: "articles", Action: "write"}})
	require.NoError(t, err)

	// Get roles (empty) - this caches the empty result
	roles, err := svc.GetUserRoles(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, 0, len(roles))

	// Assign role - should invalidate cache
	err = svc.AssignRole(ctx, userID, "writer")
	require.NoError(t, err)

	// Get roles again - should get fresh data (not cached empty result)
	roles, err = svc.GetUserRoles(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, 1, len(roles))
	assert.Contains(t, roles, "writer")
}

func TestCachedService_RemoveRole_InvalidatesCache(t *testing.T) {
	t.Parallel()
	svc, _ := setupCachedTestService(t)
	ctx := context.Background()

	userID := uuid.New()

	// Create and assign role
	_, err := svc.CreateRole(ctx, "guest", "", []Permission{{Resource: "public", Action: "read"}})
	require.NoError(t, err)

	err = svc.AssignRole(ctx, userID, "guest")
	require.NoError(t, err)

	// Get roles - caches the result
	roles, err := svc.GetUserRoles(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, 1, len(roles))

	// Remove role - should invalidate cache
	err = svc.RemoveRole(ctx, userID, "guest")
	require.NoError(t, err)

	// Get roles again - should get fresh data (empty)
	roles, err = svc.GetUserRoles(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, 0, len(roles))
}

func TestCachedService_AddPolicy_InvalidatesCache(t *testing.T) {
	t.Parallel()
	svc, _ := setupCachedTestService(t)
	ctx := context.Background()

	// Add and check policy - caches the result
	err := svc.AddPolicy(ctx, "charlie", "resource1", "read")
	require.NoError(t, err)

	allowed, err := svc.Enforce(ctx, "charlie", "resource1", "read")
	require.NoError(t, err)
	assert.True(t, allowed)

	// Remove policy via underlying service (cache should be invalidated on next add)
	err = svc.Service.RemovePolicy(ctx, "charlie", "resource1", "read")
	require.NoError(t, err)

	// Add new policy - should invalidate all RBAC cache
	err = svc.AddPolicy(ctx, "charlie", "resource2", "write")
	require.NoError(t, err)

	// Check old permission - should get fresh data (false, not cached true)
	allowed, err = svc.Enforce(ctx, "charlie", "resource1", "read")
	require.NoError(t, err)
	assert.False(t, allowed)
}

func TestCachedService_RemovePolicy_InvalidatesCache(t *testing.T) {
	t.Parallel()
	svc, _ := setupCachedTestService(t)
	ctx := context.Background()

	// Add policy and cache it
	err := svc.AddPolicy(ctx, "diana", "data1", "delete")
	require.NoError(t, err)

	allowed, err := svc.Enforce(ctx, "diana", "data1", "delete")
	require.NoError(t, err)
	assert.True(t, allowed)

	// Remove policy - should invalidate cache
	err = svc.RemovePolicy(ctx, "diana", "data1", "delete")
	require.NoError(t, err)

	// Check again - should get fresh data (false)
	allowed, err = svc.Enforce(ctx, "diana", "data1", "delete")
	require.NoError(t, err)
	assert.False(t, allowed)
}

func TestCachedService_LoadPolicy_InvalidatesCache(t *testing.T) {
	t.Parallel()
	svc, _ := setupCachedTestService(t)
	ctx := context.Background()

	// Add policy and cache it
	err := svc.AddPolicy(ctx, "eve", "data1", "manage")
	require.NoError(t, err)

	allowed, err := svc.Enforce(ctx, "eve", "data1", "manage")
	require.NoError(t, err)
	assert.True(t, allowed)

	// Load policy - should invalidate all cache
	err = svc.LoadPolicy(ctx)
	require.NoError(t, err)

	// Check permission - should get fresh data from database
	allowed, err = svc.Enforce(ctx, "eve", "data1", "manage")
	require.NoError(t, err)
	assert.True(t, allowed) // Should still be true since we loaded from DB
}
