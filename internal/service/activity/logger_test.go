package activity

import (
	"context"
	"net"
	"testing"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/infra/logging"
	"github.com/lusoris/revenge/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoopLogger(t *testing.T) {
	logger := NewNoopLogger()

	t.Run("LogAction does nothing", func(t *testing.T) {
		err := logger.LogAction(context.Background(), LogActionRequest{
			UserID:       uuid.Must(uuid.NewV7()),
			Username:     "testuser",
			Action:       ActionUserLogin,
			ResourceType: ResourceTypeUser,
			ResourceID:   uuid.Must(uuid.NewV7()),
		})
		assert.NoError(t, err)
	})

	t.Run("LogFailure does nothing", func(t *testing.T) {
		err := logger.LogFailure(context.Background(), LogFailureRequest{
			Action:       ActionUserLogin,
			ErrorMessage: "test error",
		})
		assert.NoError(t, err)
	})
}

func setupLoggerTestService(t *testing.T) (*Service, testutil.DB) {
	t.Helper()
	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())
	repo := NewRepositoryPg(queries)
	logger := logging.NewTestLogger()
	svc := NewService(repo, logger)
	return svc, testDB
}

func TestServiceLogger(t *testing.T) {
	t.Parallel()
	service, _ := setupLoggerTestService(t)
	logger := NewLogger(service)

	t.Run("LogAction creates entry", func(t *testing.T) {
		ctx := context.Background()
		// Use uuid.Nil since we don't have actual users in test DB
		resourceID := uuid.Must(uuid.NewV7())
		ip := net.ParseIP("192.168.1.100")

		err := logger.LogAction(ctx, LogActionRequest{
			UserID:       uuid.Nil, // Use Nil to avoid FK constraint
			Username:     "testuser",
			Action:       ActionUserLogin,
			ResourceType: ResourceTypeUser,
			ResourceID:   resourceID,
			IPAddress:    ip,
			UserAgent:    "TestAgent/1.0",
			Metadata: map[string]interface{}{
				"device": "desktop",
			},
		})
		require.NoError(t, err)

		// Verify entry was created by action search (since userID is nil)
		action := ActionUserLogin
		entries, _, err := service.Search(ctx, SearchFilters{
			Action: &action,
			Limit:  10,
		})
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(entries), 1)

		// Find our entry
		var found bool
		for _, entry := range entries {
			if entry.Username != nil && *entry.Username == "testuser" {
				assert.Equal(t, ActionUserLogin, entry.Action)
				assert.Equal(t, true, entry.Success)
				found = true
				break
			}
		}
		assert.True(t, found, "entry should be found")
	})

	t.Run("LogFailure creates failed entry", func(t *testing.T) {
		ctx := context.Background()
		username := "faileduser"
		ip := net.ParseIP("10.0.0.1")
		userAgent := "FailedAgent/1.0"

		err := logger.LogFailure(ctx, LogFailureRequest{
			UserID:       nil, // No user for failed login
			Username:     &username,
			Action:       ActionUserLogin,
			ErrorMessage: "invalid credentials",
			IPAddress:    &ip,
			UserAgent:    &userAgent,
		})
		require.NoError(t, err)

		// Verify entry was created
		success := false
		entries, _, err := service.Search(ctx, SearchFilters{
			Success: &success,
			Limit:   10,
		})
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(entries), 1)

		// Find our entry
		var found bool
		for _, entry := range entries {
			if entry.Username != nil && *entry.Username == "faileduser" {
				assert.Equal(t, ActionUserLogin, entry.Action)
				assert.Equal(t, false, entry.Success)
				assert.Equal(t, "invalid credentials", *entry.ErrorMessage)
				found = true
				break
			}
		}
		assert.True(t, found, "entry should be found")
	})
}
