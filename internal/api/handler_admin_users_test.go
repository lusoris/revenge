package api

import (
	"context"
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/infra/logging"
	"github.com/lusoris/revenge/internal/service/activity"
	"github.com/lusoris/revenge/internal/service/rbac"
	"github.com/lusoris/revenge/internal/service/storage"
	"github.com/lusoris/revenge/internal/service/user"
	"github.com/lusoris/revenge/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupAdminUserTestHandler(t *testing.T) (*Handler, testutil.DB, uuid.UUID) {
	t.Helper()
	testDB := testutil.NewFastTestDB(t)
	pool := testDB.Pool()
	queries := db.New(pool)

	// Clear casbin policies for test isolation
	_, err := pool.Exec(context.Background(), "DELETE FROM shared.casbin_rule")
	require.NoError(t, err)

	// Set up services
	userRepo := user.NewPostgresRepository(queries)
	userService := user.NewCachedService(
		user.NewService(pool, userRepo, activity.NewNoopLogger(), &storage.MockStorage{}, config.AvatarConfig{}),
		nil, // nil cache = pass-through mode
		logging.NewTestLogger(),
	)

	adapter := rbac.NewAdapter(pool)
	enforcer, err := casbin.NewSyncedEnforcer("../../config/casbin_model.conf", adapter)
	require.NoError(t, err)
	rbacService := rbac.NewService(enforcer, logging.NewTestLogger(), activity.NewNoopLogger())

	// Create admin user
	admin := testutil.CreateUser(t, pool, testutil.User{
		Username: "admin",
		Email:    "admin@example.com",
		IsActive: true,
		IsAdmin:  true,
	})
	err = rbacService.AssignRole(context.Background(), admin.ID, "admin")
	require.NoError(t, err)

	handler := &Handler{
		logger:      logging.NewTestLogger(),
		userService: userService,
		rbacService: rbacService,
		cfg:         &config.Config{},
	}

	return handler, testDB, admin.ID
}

// ============================================================================
// AdminListUsers Tests
// ============================================================================

func TestHandler_AdminListUsers_Unauthorized(t *testing.T) {
	t.Parallel()
	handler, _, _ := setupAdminUserTestHandler(t)

	result, err := handler.AdminListUsers(context.Background(), ogen.AdminListUsersParams{})
	require.NoError(t, err)

	resp, ok := result.(*ogen.AdminListUsersUnauthorized)
	require.True(t, ok)
	assert.Equal(t, 401, resp.Code)
}

func TestHandler_AdminListUsers_Success(t *testing.T) {
	t.Parallel()
	handler, testDB, adminID := setupAdminUserTestHandler(t)

	// Create additional users
	testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "alice", Email: "alice@example.com", IsActive: true,
	})
	testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "bob", Email: "bob@example.com", IsActive: true,
	})

	ctx := contextWithUserID(context.Background(), adminID)
	result, err := handler.AdminListUsers(ctx, ogen.AdminListUsersParams{})
	require.NoError(t, err)

	resp, ok := result.(*ogen.AdminUserListResponse)
	require.True(t, ok)
	assert.Equal(t, int64(3), resp.Total) // admin + alice + bob
	assert.Len(t, resp.Users, 3)
}

func TestHandler_AdminListUsers_SearchByUsername(t *testing.T) {
	t.Parallel()
	handler, testDB, adminID := setupAdminUserTestHandler(t)

	testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "alice", Email: "alice@example.com", IsActive: true,
	})
	testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "bob", Email: "bob@example.com", IsActive: true,
	})

	ctx := contextWithUserID(context.Background(), adminID)
	result, err := handler.AdminListUsers(ctx, ogen.AdminListUsersParams{
		Query: ogen.NewOptString("alice"),
	})
	require.NoError(t, err)

	resp, ok := result.(*ogen.AdminUserListResponse)
	require.True(t, ok)
	assert.Equal(t, int64(1), resp.Total)
	assert.Len(t, resp.Users, 1)
	assert.Equal(t, "alice", resp.Users[0].Username)
}

func TestHandler_AdminListUsers_FilterByAdmin(t *testing.T) {
	t.Parallel()
	handler, testDB, adminID := setupAdminUserTestHandler(t)

	testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "regular", Email: "regular@example.com", IsActive: true,
	})

	ctx := contextWithUserID(context.Background(), adminID)
	result, err := handler.AdminListUsers(ctx, ogen.AdminListUsersParams{
		IsAdmin: ogen.NewOptBool(true),
	})
	require.NoError(t, err)

	resp, ok := result.(*ogen.AdminUserListResponse)
	require.True(t, ok)
	assert.Equal(t, int64(1), resp.Total)
	assert.Equal(t, "admin", resp.Users[0].Username)
}

func TestHandler_AdminListUsers_Pagination(t *testing.T) {
	t.Parallel()
	handler, testDB, adminID := setupAdminUserTestHandler(t)

	for i := 0; i < 5; i++ {
		testutil.CreateUser(t, testDB.Pool(), testutil.User{
			Username: "user" + string(rune('a'+i)),
			Email:    "user" + string(rune('a'+i)) + "@example.com",
			IsActive: true,
		})
	}

	ctx := contextWithUserID(context.Background(), adminID)
	result, err := handler.AdminListUsers(ctx, ogen.AdminListUsersParams{
		Limit:  ogen.NewOptInt(2),
		Offset: ogen.NewOptInt(0),
	})
	require.NoError(t, err)

	resp, ok := result.(*ogen.AdminUserListResponse)
	require.True(t, ok)
	assert.Equal(t, int64(6), resp.Total) // admin + 5 users
	assert.Len(t, resp.Users, 2)
}

// ============================================================================
// AdminDeleteUser Tests
// ============================================================================

func TestHandler_AdminDeleteUser_Unauthorized(t *testing.T) {
	t.Parallel()
	handler, _, _ := setupAdminUserTestHandler(t)

	result, err := handler.AdminDeleteUser(context.Background(), ogen.AdminDeleteUserParams{
		UserId: uuid.Must(uuid.NewV7()),
	})
	require.NoError(t, err)

	resp, ok := result.(*ogen.AdminDeleteUserUnauthorized)
	require.True(t, ok)
	assert.Equal(t, 401, resp.Code)
}

func TestHandler_AdminDeleteUser_Success(t *testing.T) {
	t.Parallel()
	handler, testDB, adminID := setupAdminUserTestHandler(t)

	target := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "doomed", Email: "doomed@example.com", IsActive: true,
	})

	ctx := contextWithUserID(context.Background(), adminID)
	result, err := handler.AdminDeleteUser(ctx, ogen.AdminDeleteUserParams{
		UserId: target.ID,
	})
	require.NoError(t, err)

	_, ok := result.(*ogen.AdminDeleteUserNoContent)
	require.True(t, ok)

	// Verify user is soft-deleted (not returned by GetUser)
	_, err = handler.userService.GetUser(ctx, target.ID)
	assert.Error(t, err)
}

func TestHandler_AdminDeleteUser_SelfDelete(t *testing.T) {
	t.Parallel()
	handler, _, adminID := setupAdminUserTestHandler(t)

	ctx := contextWithUserID(context.Background(), adminID)
	result, err := handler.AdminDeleteUser(ctx, ogen.AdminDeleteUserParams{
		UserId: adminID,
	})
	require.NoError(t, err)

	resp, ok := result.(*ogen.AdminDeleteUserForbidden)
	require.True(t, ok)
	assert.Contains(t, resp.Message, "Cannot delete your own account")
}

func TestHandler_AdminDeleteUser_NotFound(t *testing.T) {
	t.Parallel()
	handler, _, adminID := setupAdminUserTestHandler(t)

	ctx := contextWithUserID(context.Background(), adminID)
	result, err := handler.AdminDeleteUser(ctx, ogen.AdminDeleteUserParams{
		UserId: uuid.Must(uuid.NewV7()),
	})
	require.NoError(t, err)

	_, ok := result.(*ogen.AdminDeleteUserNotFound)
	require.True(t, ok)
}
