package api

import (
	"context"
	"testing"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/google/uuid"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/integration/radarr"
	"github.com/lusoris/revenge/internal/service/activity"
	"github.com/lusoris/revenge/internal/service/rbac"
	"github.com/lusoris/revenge/internal/testutil"
)

// mockRadarrService is a mock implementation of radarrService for testing.
type mockRadarrService struct {
	healthy         bool
	status          radarr.SyncStatus
	systemStatus    *radarr.SystemStatus
	qualityProfiles []radarr.QualityProfile
	rootFolders     []radarr.RootFolder
	syncError       error
}

func (m *mockRadarrService) GetStatus() radarr.SyncStatus {
	return m.status
}

func (m *mockRadarrService) IsHealthy(ctx context.Context) bool {
	return m.healthy
}

func (m *mockRadarrService) GetSystemStatus(ctx context.Context) (*radarr.SystemStatus, error) {
	return m.systemStatus, nil
}

func (m *mockRadarrService) GetQualityProfiles(ctx context.Context) ([]radarr.QualityProfile, error) {
	return m.qualityProfiles, nil
}

func (m *mockRadarrService) GetRootFolders(ctx context.Context) ([]radarr.RootFolder, error) {
	return m.rootFolders, nil
}

func (m *mockRadarrService) SyncLibrary(ctx context.Context) (*radarr.SyncResult, error) {
	if m.syncError != nil {
		return nil, m.syncError
	}
	return &radarr.SyncResult{
		MoviesAdded:   5,
		MoviesUpdated: 10,
	}, nil
}

// mockRiverClient is a mock implementation of riverClient for testing.
type mockRiverClient struct {
	insertedArgs []river.JobArgs
	insertError  error
}

func (m *mockRiverClient) Insert(ctx context.Context, args river.JobArgs, opts *river.InsertOpts) (*rivertype.JobInsertResult, error) {
	if m.insertError != nil {
		return nil, m.insertError
	}
	m.insertedArgs = append(m.insertedArgs, args)
	return &rivertype.JobInsertResult{}, nil
}

func setupRadarrTestHandler(t *testing.T) (*Handler, *testutil.TestDB, uuid.UUID) {
	t.Helper()
	testDB := testutil.NewTestDB(t)

	// Clear any existing policies from the table to ensure test isolation
	_, err := testDB.Pool().Exec(context.Background(), "DELETE FROM shared.casbin_rule")
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

	// Grant admin role
	err = rbacService.AssignRole(context.Background(), adminUser.ID, "admin")
	require.NoError(t, err)

	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTExpiry: 15 * time.Minute,
		},
	}

	handler := &Handler{
		logger:      zap.NewNop(),
		rbacService: rbacService,
		cfg:         cfg,
	}

	return handler, testDB, adminUser.ID
}

// ============================================================================
// AdminGetRadarrStatus Tests
// ============================================================================

func TestHandler_AdminGetRadarrStatus_NotAdmin(t *testing.T) {
	t.Parallel()
	handler, _, _ := setupRadarrTestHandler(t)

	ctx := context.Background()

	result, err := handler.AdminGetRadarrStatus(ctx)
	require.NoError(t, err)

	forbidden, ok := result.(*ogen.AdminGetRadarrStatusForbidden)
	require.True(t, ok)
	assert.Equal(t, 403, forbidden.Code)
	assert.Contains(t, forbidden.Message, "Admin access required")
}

func TestHandler_AdminGetRadarrStatus_NotConfigured(t *testing.T) {
	t.Parallel()
	handler, _, adminID := setupRadarrTestHandler(t)

	ctx := WithUserID(context.Background(), adminID)

	result, err := handler.AdminGetRadarrStatus(ctx)
	require.NoError(t, err)

	unavailable, ok := result.(*ogen.AdminGetRadarrStatusServiceUnavailable)
	require.True(t, ok)
	assert.Equal(t, 503, unavailable.Code)
	assert.Contains(t, unavailable.Message, "not configured")
}

func TestHandler_AdminGetRadarrStatus_Connected(t *testing.T) {
	t.Parallel()
	handler, _, adminID := setupRadarrTestHandler(t)

	mockService := &mockRadarrService{
		healthy: true,
		status: radarr.SyncStatus{
			IsRunning:     false,
			LastSync:      time.Now().Add(-1 * time.Hour),
			MoviesAdded:   5,
			MoviesUpdated: 10,
			TotalMovies:   100,
		},
		systemStatus: &radarr.SystemStatus{
			Version:      "4.7.0",
			InstanceName: "Test Radarr",
			StartTime:    time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
		},
	}
	handler.radarrService = mockService

	ctx := WithUserID(context.Background(), adminID)

	result, err := handler.AdminGetRadarrStatus(ctx)
	require.NoError(t, err)

	status, ok := result.(*ogen.RadarrStatus)
	require.True(t, ok)
	assert.True(t, status.Connected)
	assert.Equal(t, "4.7.0", status.Version.Value)
	assert.Equal(t, "Test Radarr", status.InstanceName.Value)
	assert.False(t, status.SyncStatus.IsRunning)
	assert.Equal(t, 5, status.SyncStatus.MoviesAdded)
	assert.Equal(t, 10, status.SyncStatus.MoviesUpdated)
	assert.Equal(t, 100, status.SyncStatus.TotalMovies)
}

func TestHandler_AdminGetRadarrStatus_Disconnected(t *testing.T) {
	t.Parallel()
	handler, _, adminID := setupRadarrTestHandler(t)

	mockService := &mockRadarrService{
		healthy: false,
		status: radarr.SyncStatus{
			IsRunning:     false,
			LastSyncError: "connection refused",
		},
	}
	handler.radarrService = mockService

	ctx := WithUserID(context.Background(), adminID)

	result, err := handler.AdminGetRadarrStatus(ctx)
	require.NoError(t, err)

	status, ok := result.(*ogen.RadarrStatus)
	require.True(t, ok)
	assert.False(t, status.Connected)
	assert.Equal(t, "connection refused", status.SyncStatus.LastSyncError.Value)
}

// ============================================================================
// AdminTriggerRadarrSync Tests
// ============================================================================

func TestHandler_AdminTriggerRadarrSync_NotAdmin(t *testing.T) {
	t.Parallel()
	handler, _, _ := setupRadarrTestHandler(t)

	ctx := context.Background()

	result, err := handler.AdminTriggerRadarrSync(ctx)
	require.NoError(t, err)

	forbidden, ok := result.(*ogen.AdminTriggerRadarrSyncForbidden)
	require.True(t, ok)
	assert.Equal(t, 403, forbidden.Code)
}

func TestHandler_AdminTriggerRadarrSync_NotConfigured(t *testing.T) {
	t.Parallel()
	handler, _, adminID := setupRadarrTestHandler(t)

	ctx := WithUserID(context.Background(), adminID)

	result, err := handler.AdminTriggerRadarrSync(ctx)
	require.NoError(t, err)

	unavailable, ok := result.(*ogen.AdminTriggerRadarrSyncServiceUnavailable)
	require.True(t, ok)
	assert.Equal(t, 503, unavailable.Code)
}

func TestHandler_AdminTriggerRadarrSync_AlreadyRunning(t *testing.T) {
	t.Parallel()
	handler, _, adminID := setupRadarrTestHandler(t)

	mockService := &mockRadarrService{
		status: radarr.SyncStatus{
			IsRunning: true,
		},
	}
	handler.radarrService = mockService

	ctx := WithUserID(context.Background(), adminID)

	result, err := handler.AdminTriggerRadarrSync(ctx)
	require.NoError(t, err)

	conflict, ok := result.(*ogen.AdminTriggerRadarrSyncConflict)
	require.True(t, ok)
	assert.Equal(t, 409, conflict.Code)
	assert.Contains(t, conflict.Message, "already in progress")
}

func TestHandler_AdminTriggerRadarrSync_WithRiver(t *testing.T) {
	t.Parallel()
	handler, _, adminID := setupRadarrTestHandler(t)

	mockService := &mockRadarrService{
		status: radarr.SyncStatus{IsRunning: false},
	}
	mockRiver := &mockRiverClient{}
	handler.radarrService = mockService
	handler.riverClient = mockRiver

	ctx := WithUserID(context.Background(), adminID)

	result, err := handler.AdminTriggerRadarrSync(ctx)
	require.NoError(t, err)

	response, ok := result.(*ogen.RadarrSyncResponse)
	require.True(t, ok)
	assert.Equal(t, ogen.RadarrSyncResponseStatusQueued, response.Status)
	assert.Contains(t, response.Message, "queued")

	// Verify job was inserted
	require.Len(t, mockRiver.insertedArgs, 1)
	syncArgs, ok := mockRiver.insertedArgs[0].(*radarr.RadarrSyncJobArgs)
	require.True(t, ok)
	assert.Equal(t, radarr.RadarrSyncOperationFull, syncArgs.Operation)
}

func TestHandler_AdminTriggerRadarrSync_DirectSync(t *testing.T) {
	t.Parallel()
	handler, _, adminID := setupRadarrTestHandler(t)

	mockService := &mockRadarrService{
		status: radarr.SyncStatus{IsRunning: false},
	}
	handler.radarrService = mockService
	// No river client - should start sync directly

	ctx := WithUserID(context.Background(), adminID)

	result, err := handler.AdminTriggerRadarrSync(ctx)
	require.NoError(t, err)

	response, ok := result.(*ogen.RadarrSyncResponse)
	require.True(t, ok)
	assert.Equal(t, ogen.RadarrSyncResponseStatusStarted, response.Status)
	assert.Contains(t, response.Message, "started")
}

// ============================================================================
// AdminGetRadarrQualityProfiles Tests
// ============================================================================

func TestHandler_AdminGetRadarrQualityProfiles_NotAdmin(t *testing.T) {
	t.Parallel()
	handler, _, _ := setupRadarrTestHandler(t)

	ctx := context.Background()

	result, err := handler.AdminGetRadarrQualityProfiles(ctx)
	require.NoError(t, err)

	forbidden, ok := result.(*ogen.AdminGetRadarrQualityProfilesForbidden)
	require.True(t, ok)
	assert.Equal(t, 403, forbidden.Code)
}

func TestHandler_AdminGetRadarrQualityProfiles_Success(t *testing.T) {
	t.Parallel()
	handler, _, adminID := setupRadarrTestHandler(t)

	mockService := &mockRadarrService{
		qualityProfiles: []radarr.QualityProfile{
			{ID: 1, Name: "HD-1080p", UpgradeAllowed: true, Cutoff: 7},
			{ID: 2, Name: "Any", UpgradeAllowed: false, Cutoff: 0},
		},
	}
	handler.radarrService = mockService

	ctx := WithUserID(context.Background(), adminID)

	result, err := handler.AdminGetRadarrQualityProfiles(ctx)
	require.NoError(t, err)

	profiles, ok := result.(*ogen.RadarrQualityProfileList)
	require.True(t, ok)
	assert.Len(t, profiles.Profiles, 2)
	assert.Equal(t, 1, profiles.Profiles[0].ID)
	assert.Equal(t, "HD-1080p", profiles.Profiles[0].Name)
	assert.True(t, profiles.Profiles[0].UpgradeAllowed.Value)
}

// ============================================================================
// AdminGetRadarrRootFolders Tests
// ============================================================================

func TestHandler_AdminGetRadarrRootFolders_NotAdmin(t *testing.T) {
	t.Parallel()
	handler, _, _ := setupRadarrTestHandler(t)

	ctx := context.Background()

	result, err := handler.AdminGetRadarrRootFolders(ctx)
	require.NoError(t, err)

	forbidden, ok := result.(*ogen.AdminGetRadarrRootFoldersForbidden)
	require.True(t, ok)
	assert.Equal(t, 403, forbidden.Code)
}

func TestHandler_AdminGetRadarrRootFolders_Success(t *testing.T) {
	t.Parallel()
	handler, _, adminID := setupRadarrTestHandler(t)

	mockService := &mockRadarrService{
		rootFolders: []radarr.RootFolder{
			{ID: 1, Path: "/movies", Accessible: true, FreeSpace: 1000000000000},
			{ID: 2, Path: "/archive", Accessible: false, FreeSpace: 0},
		},
	}
	handler.radarrService = mockService

	ctx := WithUserID(context.Background(), adminID)

	result, err := handler.AdminGetRadarrRootFolders(ctx)
	require.NoError(t, err)

	folders, ok := result.(*ogen.RadarrRootFolderList)
	require.True(t, ok)
	assert.Len(t, folders.Folders, 2)
	assert.Equal(t, 1, folders.Folders[0].ID)
	assert.Equal(t, "/movies", folders.Folders[0].Path)
	assert.True(t, folders.Folders[0].Accessible)
	assert.Equal(t, int64(1000000000000), folders.Folders[0].FreeSpace.Value)
	assert.False(t, folders.Folders[1].Accessible)
}

// ============================================================================
// HandleRadarrWebhook Tests
// ============================================================================

func TestHandler_HandleRadarrWebhook_Success(t *testing.T) {
	t.Parallel()
	handler, _, _ := setupRadarrTestHandler(t)

	mockRiver := &mockRiverClient{}
	handler.riverClient = mockRiver

	ctx := context.Background()
	payload := &ogen.RadarrWebhookPayload{
		EventType:    ogen.RadarrWebhookPayloadEventTypeDownload,
		InstanceName: ogen.NewOptString("Test Radarr"),
	}
	payload.Movie = ogen.NewOptRadarrWebhookMovie(ogen.RadarrWebhookMovie{
		ID:    ogen.NewOptInt(123),
		Title: ogen.NewOptString("Test Movie"),
		Year:  ogen.NewOptInt(2024),
	})

	result, err := handler.HandleRadarrWebhook(ctx, payload)
	require.NoError(t, err)

	_, ok := result.(*ogen.HandleRadarrWebhookAccepted)
	require.True(t, ok)

	// Verify job was queued
	require.Len(t, mockRiver.insertedArgs, 1)
}

func TestHandler_HandleRadarrWebhook_NoRiver(t *testing.T) {
	t.Parallel()
	handler, _, _ := setupRadarrTestHandler(t)

	// No river client - should still accept webhook
	ctx := context.Background()
	payload := &ogen.RadarrWebhookPayload{
		EventType: ogen.RadarrWebhookPayloadEventTypeTest,
	}

	result, err := handler.HandleRadarrWebhook(ctx, payload)
	require.NoError(t, err)

	_, ok := result.(*ogen.HandleRadarrWebhookAccepted)
	require.True(t, ok)
}
