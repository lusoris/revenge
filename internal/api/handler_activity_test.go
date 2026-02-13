package api

import (
	"context"
	"testing"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/infra/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/service/activity"
	"github.com/lusoris/revenge/internal/service/rbac"
	"github.com/lusoris/revenge/internal/testutil"
)

func setupActivityTestHandler(t *testing.T) (*Handler, testutil.DB, uuid.UUID) {
	t.Helper()
	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())

	// Clear any existing policies from the table to ensure test isolation
	_, err := testDB.Pool().Exec(context.Background(), "DELETE FROM shared.casbin_rule")
	require.NoError(t, err)

	// Set up activity service
	activityRepo := activity.NewRepositoryPg(queries)
	activityService := activity.NewService(activityRepo, logging.NewTestLogger())

	// Set up RBAC service with Casbin
	adapter := rbac.NewAdapter(testDB.Pool())
	modelPath := "../../config/casbin_model.conf"
	enforcer, err := casbin.NewSyncedEnforcer(modelPath, adapter)
	require.NoError(t, err)
	rbacService := rbac.NewService(enforcer, logging.NewTestLogger(), activity.NewNoopLogger())

	// Create admin user
	adminUser := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "admin",
		Email:    "admin@example.com",
	})

	// Grant admin role
	err = rbacService.AssignRole(context.Background(), adminUser.ID, "admin")
	require.NoError(t, err)

	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTExpiry: 15 * time.Minute,
		},
	}

	handler := &Handler{
		logger:          logging.NewTestLogger(),
		activityService: activityService,
		rbacService:     rbacService,
		cfg:             cfg,
	}

	return handler, testDB, adminUser.ID
}

func TestHandler_SearchActivityLogs_Unauthenticated(t *testing.T) {
	t.Parallel()
	handler, _, _ := setupActivityTestHandler(t)

	ctx := context.Background()
	params := ogen.SearchActivityLogsParams{}

	result, err := handler.SearchActivityLogs(ctx, params)
	require.NoError(t, err)

	unauthorized, ok := result.(*ogen.SearchActivityLogsUnauthorized)
	require.True(t, ok, "expected *ogen.SearchActivityLogsUnauthorized, got %T", result)
	assert.Equal(t, 401, unauthorized.Code)
	assert.Contains(t, unauthorized.Message, "Authentication required")
}

func TestHandler_SearchActivityLogs_Success(t *testing.T) {
	t.Parallel()
	handler, testDB, adminID := setupActivityTestHandler(t)

	// Create a user for activity log
	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "testuser",
		Email:    "test@example.com",
	})

	// Create activity log entry
	err := handler.activityService.Log(context.Background(), activity.LogRequest{
		UserID:       &user.ID,
		Action:       "test_action",
		ResourceType: stringPtr("test_resource"),
		ResourceID:   uuidPtr(uuid.MustParse("00000000-0000-0000-0000-000000000123")),
		Success:      true,
	})
	require.NoError(t, err)

	// Create admin context
	ctx := contextWithUserID(context.Background(), adminID)

	params := ogen.SearchActivityLogsParams{
		Action: ogen.NewOptString("test_action"),
	}

	result, err := handler.SearchActivityLogs(ctx, params)
	require.NoError(t, err)

	response, ok := result.(*ogen.ActivityLogListResponse)
	require.True(t, ok)
	assert.GreaterOrEqual(t, response.Total, int64(1))
	assert.GreaterOrEqual(t, len(response.Entries), 1)
}

func TestHandler_SearchActivityLogs_WithFilters(t *testing.T) {
	t.Parallel()
	handler, testDB, adminID := setupActivityTestHandler(t)

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "filteruser",
		Email:    "filter@example.com",
	})

	// Create multiple activity entries
	for i := 0; i < 3; i++ {
		err := handler.activityService.Log(context.Background(), activity.LogRequest{
			UserID:       &user.ID,
			Action:       "filtered_action",
			ResourceType: stringPtr("filtered_resource"),
			ResourceID:   uuidPtr(uuid.MustParse("00000000-0000-0000-0000-000000000456")),
			Success:      true,
		})
		require.NoError(t, err)
	}

	ctx := contextWithUserID(context.Background(), adminID)

	// Test with user filter
	params := ogen.SearchActivityLogsParams{
		UserID: ogen.NewOptUUID(user.ID),
		Limit:  ogen.NewOptInt(10),
		Offset: ogen.NewOptInt(0),
	}

	result, err := handler.SearchActivityLogs(ctx, params)
	require.NoError(t, err)

	response, ok := result.(*ogen.ActivityLogListResponse)
	require.True(t, ok)
	assert.GreaterOrEqual(t, response.Total, int64(3))
}

func TestHandler_GetUserActivityLogs_Unauthenticated(t *testing.T) {
	t.Parallel()
	handler, _, _ := setupActivityTestHandler(t)

	ctx := context.Background()
	params := ogen.GetUserActivityLogsParams{
		UserId: uuid.Must(uuid.NewV7()),
	}

	result, err := handler.GetUserActivityLogs(ctx, params)
	require.NoError(t, err)

	unauthorized, ok := result.(*ogen.GetUserActivityLogsUnauthorized)
	require.True(t, ok, "expected *ogen.GetUserActivityLogsUnauthorized, got %T", result)
	assert.Equal(t, 401, unauthorized.Code)
}

func TestHandler_GetUserActivityLogs_Success(t *testing.T) {
	t.Parallel()
	handler, testDB, adminID := setupActivityTestHandler(t)

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "userlog",
		Email:    "userlog@example.com",
	})

	// Create activity for this user
	err := handler.activityService.Log(context.Background(), activity.LogRequest{
		UserID:       &user.ID,
		Action:       "user_action",
		ResourceType: stringPtr("user_resource"),
		Success:      true,
	})
	require.NoError(t, err)

	ctx := contextWithUserID(context.Background(), adminID)
	params := ogen.GetUserActivityLogsParams{
		UserId: user.ID,
		Limit:  ogen.NewOptInt(50),
		Offset: ogen.NewOptInt(0),
	}

	result, err := handler.GetUserActivityLogs(ctx, params)
	require.NoError(t, err)

	response, ok := result.(*ogen.ActivityLogListResponse)
	require.True(t, ok)
	assert.GreaterOrEqual(t, response.Total, int64(1))
}

func TestHandler_GetResourceActivityLogs_Unauthenticated(t *testing.T) {
	t.Parallel()
	handler, _, _ := setupActivityTestHandler(t)

	ctx := context.Background()
	params := ogen.GetResourceActivityLogsParams{
		ResourceType: "test_type",
		ResourceId:   uuid.Must(uuid.NewV7()),
	}

	result, err := handler.GetResourceActivityLogs(ctx, params)
	require.NoError(t, err)

	unauthorized, ok := result.(*ogen.GetResourceActivityLogsUnauthorized)
	require.True(t, ok, "expected *ogen.GetResourceActivityLogsUnauthorized, got %T", result)
	assert.Equal(t, 401, unauthorized.Code)
}

func TestHandler_GetResourceActivityLogs_Success(t *testing.T) {
	t.Parallel()
	handler, testDB, adminID := setupActivityTestHandler(t)

	resourceID := uuid.MustParse("00000000-0000-0000-0000-000000000999")
	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "resourceuser",
		Email:    "resource@example.com",
	})

	// Create activity for a specific resource
	err := handler.activityService.Log(context.Background(), activity.LogRequest{
		UserID:       &user.ID,
		Action:       "resource_action",
		ResourceType: stringPtr("test_resource"),
		ResourceID:   &resourceID,
		Success:      true,
	})
	require.NoError(t, err)

	ctx := contextWithUserID(context.Background(), adminID)
	params := ogen.GetResourceActivityLogsParams{
		ResourceType: "test_resource",
		ResourceId:   resourceID,
		Limit:        ogen.NewOptInt(50),
		Offset:       ogen.NewOptInt(0),
	}

	result, err := handler.GetResourceActivityLogs(ctx, params)
	require.NoError(t, err)

	response, ok := result.(*ogen.ActivityLogListResponse)
	require.True(t, ok)
	assert.GreaterOrEqual(t, response.Total, int64(1))
}

func TestHandler_GetActivityStats_Unauthenticated(t *testing.T) {
	t.Parallel()
	handler, _, _ := setupActivityTestHandler(t)

	ctx := context.Background()

	result, err := handler.GetActivityStats(ctx)
	require.NoError(t, err)

	unauthorized, ok := result.(*ogen.GetActivityStatsUnauthorized)
	require.True(t, ok, "expected *ogen.GetActivityStatsUnauthorized, got %T", result)
	assert.Equal(t, 401, unauthorized.Code)
}

func TestHandler_GetActivityStats_Success(t *testing.T) {
	t.Parallel()
	handler, testDB, adminID := setupActivityTestHandler(t)

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "statsuser",
		Email:    "stats@example.com",
	})

	// Create some activity
	for i := 0; i < 5; i++ {
		err := handler.activityService.Log(context.Background(), activity.LogRequest{
			UserID:  &user.ID,
			Action:  "stats_action",
			Success: true,
		})
		require.NoError(t, err)
	}

	ctx := contextWithUserID(context.Background(), adminID)

	result, err := handler.GetActivityStats(ctx)
	require.NoError(t, err)

	response, ok := result.(*ogen.ActivityStats)
	require.True(t, ok)
	assert.NotNil(t, response)
}

func TestHandler_GetRecentActions_Unauthenticated(t *testing.T) {
	t.Parallel()
	handler, _, _ := setupActivityTestHandler(t)

	ctx := context.Background()
	params := ogen.GetRecentActionsParams{}

	result, err := handler.GetRecentActions(ctx, params)
	require.NoError(t, err)

	unauthorized, ok := result.(*ogen.GetRecentActionsUnauthorized)
	require.True(t, ok, "expected *ogen.GetRecentActionsUnauthorized, got %T", result)
	assert.Equal(t, 401, unauthorized.Code)
}

func TestHandler_GetRecentActions_Success(t *testing.T) {
	t.Parallel()
	handler, testDB, adminID := setupActivityTestHandler(t)

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "actionsuser",
		Email:    "actions@example.com",
	})

	// Create activity with various actions
	actions := []string{"login", "logout", "update_profile", "create_post"}
	for _, action := range actions {
		err := handler.activityService.Log(context.Background(), activity.LogRequest{
			UserID:  &user.ID,
			Action:  action,
			Success: true,
		})
		require.NoError(t, err)
	}

	ctx := contextWithUserID(context.Background(), adminID)
	params := ogen.GetRecentActionsParams{
		Limit: ogen.NewOptInt(10),
	}

	result, err := handler.GetRecentActions(ctx, params)
	require.NoError(t, err)

	response, ok := result.(*ogen.ActionCountListResponse)
	require.True(t, ok)
	assert.NotNil(t, response)
	assert.Greater(t, len(response.Actions), 0)
}

func TestHandler_GetMyActivity_NoUserID(t *testing.T) {
	t.Parallel()
	handler, _, _ := setupActivityTestHandler(t)

	ctx := context.Background()
	params := ogen.GetRecentActionsParams{}

	result, err := handler.GetRecentActions(ctx, params)
	require.NoError(t, err)

	unauthorized, ok := result.(*ogen.GetRecentActionsUnauthorized)
	require.True(t, ok, "expected *ogen.GetRecentActionsUnauthorized, got %T", result)
	assert.Equal(t, 401, unauthorized.Code)
}

// Helper functions

func contextWithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, contextKeyUserID, userID)
}

func stringPtr(s string) *string {
	return &s
}

func uuidPtr(u uuid.UUID) *uuid.UUID {
	return &u
}
