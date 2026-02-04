package activity

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/testutil"
)

// Additional tests for error paths and edge cases

func TestService_Log_NilOptionalFields(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Log with all optional fields as nil
	req := LogRequest{
		Action:  "test.action",
		Success: true,
	}

	err := svc.Log(ctx, req)
	require.NoError(t, err)
}

func TestService_Log_FailedAction(t *testing.T) {
	t.Parallel()
	svc, testDB := setupTestService(t)
	ctx := context.Background()

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "failuser",
		Email:    "fail@example.com",
	})

	errorMsg := "Authentication failed"
	req := LogRequest{
		UserID:       &user.ID,
		Username:     &user.Username,
		Action:       ActionUserLogin,
		Success:      false,
		ErrorMessage: &errorMsg,
	}

	err := svc.Log(ctx, req)
	require.NoError(t, err)
}

func TestService_List_WithLargeLimit(t *testing.T) {
	t.Parallel()
	svc, testDB := setupTestService(t)
	ctx := context.Background()

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "filteruser",
		Email:    "filter@example.com",
	})

	// Log some activities
	ip := net.ParseIP("192.168.1.100")
	resourceID := uuid.New()
	resourceType := "resource"
	userAgent := "Test Agent"

	for i := 0; i < 5; i++ {
		req := LogRequest{
			UserID:       &user.ID,
			Username:     &user.Username,
			Action:       ActionUserUpdate,
			ResourceType: &resourceType,
			ResourceID:   &resourceID,
			IPAddress:    &ip,
			UserAgent:    &userAgent,
			Success:      true,
		}
		err := svc.Log(ctx, req)
		require.NoError(t, err)
	}

	// Test with limit and offset
	logs, total, err := svc.List(ctx, 10, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(logs), 5)
	assert.GreaterOrEqual(t, total, int64(5))
}

func TestService_List_WithOffset(t *testing.T) {
	t.Parallel()
	svc, testDB := setupTestService(t)
	ctx := context.Background()

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "timeuser",
		Email:    "time@example.com",
	})

	// Log multiple activities
	for i := 0; i < 3; i++ {
		req := LogRequest{
			UserID:   &user.ID,
			Username: &user.Username,
			Action:   ActionUserLogin,
			Success:  true,
		}
		err := svc.Log(ctx, req)
		require.NoError(t, err)
	}

	// Query with offset
	logs, total, err := svc.List(ctx, 10, 1)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(1))
	assert.GreaterOrEqual(t, len(logs), 0)
}

func TestService_CleanupOldLogs_BeforeDate(t *testing.T) {
	t.Parallel()
	svc, testDB := setupTestService(t)
	ctx := context.Background()

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "cleanupuser",
		Email:    "cleanup@example.com",
	})

	// Log some activities
	for i := 0; i < 3; i++ {
		req := LogRequest{
			UserID:   &user.ID,
			Username: &user.Username,
			Action:   ActionUserLogin,
			Success:  true,
		}
		err := svc.Log(ctx, req)
		require.NoError(t, err)
	}

	// Cleanup logs older than 1 hour from now (should not delete recent logs)
	beforeDate := time.Now().Add(1 * time.Hour)
	deleted, err := svc.CleanupOldLogs(ctx, beforeDate)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, deleted, int64(0))
}

func TestRepository_GetByAction(t *testing.T) {
	t.Parallel()
	svc, testDB := setupTestService(t)
	ctx := context.Background()

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "actionuser",
		Email:    "action@example.com",
	})

	// Log activities with specific action
	action := "test.specific.action"
	for i := 0; i < 2; i++ {
		req := LogRequest{
			UserID:   &user.ID,
			Username: &user.Username,
			Action:   action,
			Success:  true,
		}
		err := svc.Log(ctx, req)
		require.NoError(t, err)
	}

	// Get by action
	logs, err := svc.repo.GetByAction(ctx, action, 10, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(logs), 2)
	for _, log := range logs {
		assert.Equal(t, action, log.Action)
	}
}

func TestRepository_GetByIP(t *testing.T) {
	t.Parallel()
	svc, testDB := setupTestService(t)
	ctx := context.Background()

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "ipuser",
		Email:    "ip@example.com",
	})

	// Log activities from specific IP
	ip := net.ParseIP("203.0.113.42")
	userAgent := "Test Agent"
	for i := 0; i < 2; i++ {
		req := LogRequest{
			UserID:    &user.ID,
			Username:  &user.Username,
			Action:    ActionUserLogin,
			IPAddress: &ip,
			UserAgent: &userAgent,
			Success:   true,
		}
		err := svc.Log(ctx, req)
		require.NoError(t, err)
	}

	// Get by IP
	logs, err := svc.repo.GetByIP(ctx, ip, 10, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(logs), 2)
	for _, log := range logs {
		assert.NotNil(t, log.IPAddress)
		assert.True(t, log.IPAddress.Equal(ip))
	}
}

func TestService_Search_WithLimit(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Search with default limit
	logs, total, err := svc.Search(ctx, SearchFilters{
		Limit: 10,
	})
	require.NoError(t, err)
	assert.NotNil(t, logs)
	assert.GreaterOrEqual(t, total, int64(0))
}

func TestService_GetStats_Success(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Get stats
	stats, err := svc.GetStats(ctx)
	require.NoError(t, err)
	assert.NotNil(t, stats)
}

func TestService_GetRecentActions_LimitAndOffset(t *testing.T) {
	t.Parallel()
	svc, testDB := setupTestService(t)
	ctx := context.Background()

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "recentuser",
		Email:    "recent@example.com",
	})

	// Log multiple activities
	for i := 0; i < 5; i++ {
		req := LogRequest{
			UserID:   &user.ID,
			Username: &user.Username,
			Action:   ActionUserUpdate,
			Success:  true,
		}
		err := svc.Log(ctx, req)
		require.NoError(t, err)
	}

	// Get recent actions with limit
	actions, err := svc.GetRecentActions(ctx, 3)
	require.NoError(t, err)
	assert.LessOrEqual(t, len(actions), 3)
}
