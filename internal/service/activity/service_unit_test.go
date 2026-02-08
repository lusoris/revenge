package activity_test

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/infra/logging"
	"github.com/lusoris/revenge/internal/service/activity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func makeTestEntry(id uuid.UUID, action string, success bool) *activity.Entry {
	now := time.Now()
	return &activity.Entry{
		ID:        id,
		Action:    action,
		Success:   success,
		CreatedAt: now,
	}
}

func setupActivityService(repo activity.Repository) *activity.Service {
	logger := logging.NewTestLogger()
	return activity.NewService(repo, logger)
}

func TestActivityService_Log_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockActivityRepository(t)
		svc := setupActivityService(mockRepo)

		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*activity.Entry")).Return(nil)

		req := activity.LogRequest{
			Action:  activity.ActionUserLogin,
			Success: true,
		}

		err := svc.Log(context.Background(), req)

		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo := NewMockActivityRepository(t)
		svc := setupActivityService(mockRepo)

		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*activity.Entry")).Return(errors.New("db error"))

		req := activity.LogRequest{
			Action:  activity.ActionUserLogin,
			Success: true,
		}

		err := svc.Log(context.Background(), req)

		assert.Error(t, err)
	})
}

func TestActivityService_LogWithContext_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockActivityRepository(t)
		svc := setupActivityService(mockRepo)

		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*activity.Entry")).Return(nil)

		userID := uuid.Must(uuid.NewV7())
		username := "testuser"
		resourceType := "user"
		resourceID := uuid.Must(uuid.NewV7())
		ipAddress := net.ParseIP("192.168.1.1")
		userAgent := "Mozilla/5.0"

		err := svc.LogWithContext(
			context.Background(),
			userID,
			username,
			activity.ActionUserLogin,
			resourceType,
			resourceID,
			nil,
			ipAddress,
			userAgent,
		)

		assert.NoError(t, err)
	})
}

func TestActivityService_LogFailure_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockActivityRepository(t)
		svc := setupActivityService(mockRepo)

		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*activity.Entry")).Return(nil)

		userID := uuid.Must(uuid.NewV7())
		username := "testuser"
		ipAddress := net.ParseIP("192.168.1.1")
		userAgent := "Mozilla/5.0"

		err := svc.LogFailure(
			context.Background(),
			&userID,
			&username,
			activity.ActionUserLogin,
			"invalid credentials",
			&ipAddress,
			&userAgent,
		)

		assert.NoError(t, err)
	})
}

func TestActivityService_Get_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockActivityRepository(t)
		svc := setupActivityService(mockRepo)

		entryID := uuid.Must(uuid.NewV7())
		expected := makeTestEntry(entryID, activity.ActionUserLogin, true)

		mockRepo.On("Get", mock.Anything, entryID).Return(expected, nil)

		entry, err := svc.Get(context.Background(), entryID)

		require.NoError(t, err)
		assert.Equal(t, entryID, entry.ID)
		assert.Equal(t, activity.ActionUserLogin, entry.Action)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := NewMockActivityRepository(t)
		svc := setupActivityService(mockRepo)

		entryID := uuid.Must(uuid.NewV7())

		mockRepo.On("Get", mock.Anything, entryID).Return(nil, activity.ErrNotFound)

		entry, err := svc.Get(context.Background(), entryID)

		assert.Nil(t, entry)
		assert.ErrorIs(t, err, activity.ErrNotFound)
	})
}

func TestActivityService_List_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockActivityRepository(t)
		svc := setupActivityService(mockRepo)

		entries := []activity.Entry{
			*makeTestEntry(uuid.Must(uuid.NewV7()), activity.ActionUserLogin, true),
			*makeTestEntry(uuid.Must(uuid.NewV7()), activity.ActionUserLogout, true),
		}

		mockRepo.On("List", mock.Anything, int32(50), int32(0)).Return(entries, nil)
		mockRepo.On("Count", mock.Anything).Return(int64(2), nil)

		result, count, err := svc.List(context.Background(), 50, 0)

		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, int64(2), count)
	})

	t.Run("list error", func(t *testing.T) {
		mockRepo := NewMockActivityRepository(t)
		svc := setupActivityService(mockRepo)

		mockRepo.On("List", mock.Anything, int32(50), int32(0)).Return(nil, errors.New("db error"))

		result, count, err := svc.List(context.Background(), 50, 0)

		assert.Nil(t, result)
		assert.Equal(t, int64(0), count)
		assert.Error(t, err)
	})

	t.Run("count error", func(t *testing.T) {
		mockRepo := NewMockActivityRepository(t)
		svc := setupActivityService(mockRepo)

		entries := []activity.Entry{}
		mockRepo.On("List", mock.Anything, int32(50), int32(0)).Return(entries, nil)
		mockRepo.On("Count", mock.Anything).Return(int64(0), errors.New("count error"))

		result, count, err := svc.List(context.Background(), 50, 0)

		assert.Nil(t, result)
		assert.Equal(t, int64(0), count)
		assert.Error(t, err)
	})
}

func TestActivityService_Search_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success with default limit", func(t *testing.T) {
		mockRepo := NewMockActivityRepository(t)
		svc := setupActivityService(mockRepo)

		entries := []activity.Entry{
			*makeTestEntry(uuid.Must(uuid.NewV7()), activity.ActionUserLogin, true),
		}

		filters := activity.SearchFilters{}
		mockRepo.On("Search", mock.Anything, mock.MatchedBy(func(f activity.SearchFilters) bool {
			return f.Limit == 50
		})).Return(entries, int64(1), nil)

		result, count, err := svc.Search(context.Background(), filters)

		require.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, int64(1), count)
	})

	t.Run("limit capped at 100", func(t *testing.T) {
		mockRepo := NewMockActivityRepository(t)
		svc := setupActivityService(mockRepo)

		entries := []activity.Entry{}

		filters := activity.SearchFilters{Limit: 200}
		mockRepo.On("Search", mock.Anything, mock.MatchedBy(func(f activity.SearchFilters) bool {
			return f.Limit == 100
		})).Return(entries, int64(0), nil)

		result, count, err := svc.Search(context.Background(), filters)

		require.NoError(t, err)
		assert.Len(t, result, 0)
		assert.Equal(t, int64(0), count)
	})
}

func TestActivityService_GetUserActivity_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success with default limit", func(t *testing.T) {
		mockRepo := NewMockActivityRepository(t)
		svc := setupActivityService(mockRepo)

		userID := uuid.Must(uuid.NewV7())
		entries := []activity.Entry{
			*makeTestEntry(uuid.Must(uuid.NewV7()), activity.ActionUserLogin, true),
		}

		mockRepo.On("GetByUser", mock.Anything, userID, int32(50), int32(0)).Return(entries, int64(1), nil)

		result, count, err := svc.GetUserActivity(context.Background(), userID, 0, 0)

		require.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, int64(1), count)
	})

	t.Run("limit capped at 100", func(t *testing.T) {
		mockRepo := NewMockActivityRepository(t)
		svc := setupActivityService(mockRepo)

		userID := uuid.Must(uuid.NewV7())
		entries := []activity.Entry{}

		mockRepo.On("GetByUser", mock.Anything, userID, int32(100), int32(0)).Return(entries, int64(0), nil)

		result, count, err := svc.GetUserActivity(context.Background(), userID, 200, 0)

		require.NoError(t, err)
		assert.Len(t, result, 0)
		assert.Equal(t, int64(0), count)
	})
}

func TestActivityService_GetResourceActivity_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success with default limit", func(t *testing.T) {
		mockRepo := NewMockActivityRepository(t)
		svc := setupActivityService(mockRepo)

		resourceID := uuid.Must(uuid.NewV7())
		resourceType := "user"
		entries := []activity.Entry{
			*makeTestEntry(uuid.Must(uuid.NewV7()), activity.ActionUserUpdate, true),
		}

		mockRepo.On("GetByResource", mock.Anything, resourceType, resourceID, int32(50), int32(0)).
			Return(entries, int64(1), nil)

		result, count, err := svc.GetResourceActivity(context.Background(), resourceType, resourceID, 0, 0)

		require.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, int64(1), count)
	})

	t.Run("limit capped at 100", func(t *testing.T) {
		mockRepo := NewMockActivityRepository(t)
		svc := setupActivityService(mockRepo)

		resourceID := uuid.Must(uuid.NewV7())
		resourceType := "user"
		entries := []activity.Entry{}

		mockRepo.On("GetByResource", mock.Anything, resourceType, resourceID, int32(100), int32(0)).
			Return(entries, int64(0), nil)

		result, count, err := svc.GetResourceActivity(context.Background(), resourceType, resourceID, 200, 0)

		require.NoError(t, err)
		assert.Len(t, result, 0)
		assert.Equal(t, int64(0), count)
	})
}

func TestActivityService_GetFailedActivity_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success with default limit", func(t *testing.T) {
		mockRepo := NewMockActivityRepository(t)
		svc := setupActivityService(mockRepo)

		entries := []activity.Entry{
			*makeTestEntry(uuid.Must(uuid.NewV7()), activity.ActionUserLogin, false),
		}

		mockRepo.On("GetFailed", mock.Anything, int32(50), int32(0)).Return(entries, nil)

		result, err := svc.GetFailedActivity(context.Background(), 0, 0)

		require.NoError(t, err)
		assert.Len(t, result, 1)
	})

	t.Run("limit capped at 100", func(t *testing.T) {
		mockRepo := NewMockActivityRepository(t)
		svc := setupActivityService(mockRepo)

		entries := []activity.Entry{}

		mockRepo.On("GetFailed", mock.Anything, int32(100), int32(0)).Return(entries, nil)

		result, err := svc.GetFailedActivity(context.Background(), 200, 0)

		require.NoError(t, err)
		assert.Len(t, result, 0)
	})
}

func TestActivityService_GetStats_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockActivityRepository(t)
		svc := setupActivityService(mockRepo)

		stats := &activity.Stats{
			TotalCount:   100,
			SuccessCount: 95,
			FailedCount:  5,
		}

		mockRepo.On("GetStats", mock.Anything).Return(stats, nil)

		result, err := svc.GetStats(context.Background())

		require.NoError(t, err)
		assert.Equal(t, int64(100), result.TotalCount)
		assert.Equal(t, int64(95), result.SuccessCount)
		assert.Equal(t, int64(5), result.FailedCount)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo := NewMockActivityRepository(t)
		svc := setupActivityService(mockRepo)

		mockRepo.On("GetStats", mock.Anything).Return(nil, errors.New("db error"))

		result, err := svc.GetStats(context.Background())

		assert.Nil(t, result)
		assert.Error(t, err)
	})
}

func TestActivityService_GetRecentActions_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success with default limit", func(t *testing.T) {
		mockRepo := NewMockActivityRepository(t)
		svc := setupActivityService(mockRepo)

		actions := []activity.ActionCount{
			{Action: activity.ActionUserLogin, Count: 50},
			{Action: activity.ActionUserLogout, Count: 30},
		}

		mockRepo.On("GetRecentActions", mock.Anything, int32(20)).Return(actions, nil)

		result, err := svc.GetRecentActions(context.Background(), 0)

		require.NoError(t, err)
		assert.Len(t, result, 2)
	})

	t.Run("limit capped at 50", func(t *testing.T) {
		mockRepo := NewMockActivityRepository(t)
		svc := setupActivityService(mockRepo)

		actions := []activity.ActionCount{}

		mockRepo.On("GetRecentActions", mock.Anything, int32(50)).Return(actions, nil)

		result, err := svc.GetRecentActions(context.Background(), 100)

		require.NoError(t, err)
		assert.Len(t, result, 0)
	})
}

func TestActivityService_CleanupOldLogs_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockActivityRepository(t)
		svc := setupActivityService(mockRepo)

		olderThan := time.Now().Add(-30 * 24 * time.Hour)

		mockRepo.On("DeleteOld", mock.Anything, mock.AnythingOfType("time.Time")).Return(int64(100), nil)

		count, err := svc.CleanupOldLogs(context.Background(), olderThan)

		require.NoError(t, err)
		assert.Equal(t, int64(100), count)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo := NewMockActivityRepository(t)
		svc := setupActivityService(mockRepo)

		olderThan := time.Now().Add(-30 * 24 * time.Hour)

		mockRepo.On("DeleteOld", mock.Anything, mock.AnythingOfType("time.Time")).Return(int64(0), errors.New("db error"))

		count, err := svc.CleanupOldLogs(context.Background(), olderThan)

		assert.Equal(t, int64(0), count)
		assert.Error(t, err)
	})
}

func TestActivityService_CountOldLogs_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockActivityRepository(t)
		svc := setupActivityService(mockRepo)

		olderThan := time.Now().Add(-30 * 24 * time.Hour)

		mockRepo.On("CountOld", mock.Anything, mock.AnythingOfType("time.Time")).Return(int64(50), nil)

		count, err := svc.CountOldLogs(context.Background(), olderThan)

		require.NoError(t, err)
		assert.Equal(t, int64(50), count)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo := NewMockActivityRepository(t)
		svc := setupActivityService(mockRepo)

		olderThan := time.Now().Add(-30 * 24 * time.Hour)

		mockRepo.On("CountOld", mock.Anything, mock.AnythingOfType("time.Time")).Return(int64(0), errors.New("db error"))

		count, err := svc.CountOldLogs(context.Background(), olderThan)

		assert.Equal(t, int64(0), count)
		assert.Error(t, err)
	})
}

// ========== ServiceLogger Unit Tests ==========

func TestServiceLogger_LogAction_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success with all fields", func(t *testing.T) {
		mockRepo := NewMockActivityRepository(t)
		svc := setupActivityService(mockRepo)
		logger := activity.NewLogger(svc)

		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*activity.Entry")).Return(nil)

		userID := uuid.Must(uuid.NewV7())
		resourceID := uuid.Must(uuid.NewV7())
		ipAddr := net.ParseIP("192.168.1.1")

		err := logger.LogAction(context.Background(), activity.LogActionRequest{
			UserID:       userID,
			Username:     "testuser",
			Action:       activity.ActionUserLogin,
			ResourceType: activity.ResourceTypeUser,
			ResourceID:   resourceID,
			Changes:      map[string]interface{}{"field": "value"},
			Metadata:     map[string]interface{}{"key": "value"},
			IPAddress:    ipAddr,
			UserAgent:    "Mozilla/5.0",
		})

		assert.NoError(t, err)
	})

	t.Run("success with nil user ID", func(t *testing.T) {
		mockRepo := NewMockActivityRepository(t)
		svc := setupActivityService(mockRepo)
		logger := activity.NewLogger(svc)

		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(e *activity.Entry) bool {
			// UserID should be nil when uuid.Nil is passed
			return e.UserID == nil
		})).Return(nil)

		err := logger.LogAction(context.Background(), activity.LogActionRequest{
			UserID:       uuid.Nil,
			Username:     "testuser",
			Action:       activity.ActionUserLogin,
			ResourceType: activity.ResourceTypeUser,
			ResourceID:   uuid.Nil,
		})

		assert.NoError(t, err)
	})

	t.Run("success with empty username", func(t *testing.T) {
		mockRepo := NewMockActivityRepository(t)
		svc := setupActivityService(mockRepo)
		logger := activity.NewLogger(svc)

		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(e *activity.Entry) bool {
			// Username should be nil when empty string is passed
			return e.Username == nil
		})).Return(nil)

		err := logger.LogAction(context.Background(), activity.LogActionRequest{
			UserID:       uuid.Nil,
			Username:     "",
			Action:       activity.ActionUserLogin,
			ResourceType: "",
			ResourceID:   uuid.Nil,
		})

		assert.NoError(t, err)
	})

	t.Run("success with nil IP address", func(t *testing.T) {
		mockRepo := NewMockActivityRepository(t)
		svc := setupActivityService(mockRepo)
		logger := activity.NewLogger(svc)

		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(e *activity.Entry) bool {
			return e.IPAddress == nil
		})).Return(nil)

		err := logger.LogAction(context.Background(), activity.LogActionRequest{
			UserID:    uuid.Must(uuid.NewV7()),
			Username:  "testuser",
			Action:    activity.ActionUserLogin,
			IPAddress: nil,
			UserAgent: "",
		})

		assert.NoError(t, err)
	})

	t.Run("error from repository", func(t *testing.T) {
		mockRepo := NewMockActivityRepository(t)
		svc := setupActivityService(mockRepo)
		logger := activity.NewLogger(svc)

		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*activity.Entry")).Return(errors.New("db error"))

		err := logger.LogAction(context.Background(), activity.LogActionRequest{
			UserID:   uuid.Must(uuid.NewV7()),
			Username: "testuser",
			Action:   activity.ActionUserLogin,
		})

		assert.Error(t, err)
	})
}

func TestServiceLogger_LogFailure_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success with all fields", func(t *testing.T) {
		mockRepo := NewMockActivityRepository(t)
		svc := setupActivityService(mockRepo)
		logger := activity.NewLogger(svc)

		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(e *activity.Entry) bool {
			return e.Success == false && e.ErrorMessage != nil
		})).Return(nil)

		userID := uuid.Must(uuid.NewV7())
		username := "testuser"
		ipAddr := net.ParseIP("192.168.1.1")
		userAgent := "Mozilla/5.0"

		err := logger.LogFailure(context.Background(), activity.LogFailureRequest{
			UserID:       &userID,
			Username:     &username,
			Action:       activity.ActionUserLogin,
			ErrorMessage: "invalid credentials",
			IPAddress:    &ipAddr,
			UserAgent:    &userAgent,
		})

		assert.NoError(t, err)
	})

	t.Run("success with nil fields", func(t *testing.T) {
		mockRepo := NewMockActivityRepository(t)
		svc := setupActivityService(mockRepo)
		logger := activity.NewLogger(svc)

		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(e *activity.Entry) bool {
			return e.Success == false && e.UserID == nil && e.Username == nil
		})).Return(nil)

		err := logger.LogFailure(context.Background(), activity.LogFailureRequest{
			UserID:       nil,
			Username:     nil,
			Action:       activity.ActionUserLogin,
			ErrorMessage: "invalid credentials",
			IPAddress:    nil,
			UserAgent:    nil,
		})

		assert.NoError(t, err)
	})

	t.Run("error from repository", func(t *testing.T) {
		mockRepo := NewMockActivityRepository(t)
		svc := setupActivityService(mockRepo)
		logger := activity.NewLogger(svc)

		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*activity.Entry")).Return(errors.New("db error"))

		err := logger.LogFailure(context.Background(), activity.LogFailureRequest{
			Action:       activity.ActionUserLogin,
			ErrorMessage: "error",
		})

		assert.Error(t, err)
	})
}

func TestNewLogger_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	mockRepo := NewMockActivityRepository(t)
	svc := setupActivityService(mockRepo)
	logger := activity.NewLogger(svc)

	assert.NotNil(t, logger)
	// Verify it implements the Logger interface
	var _ activity.Logger = logger
}
