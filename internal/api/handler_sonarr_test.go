package api

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/integration/sonarr"
)

// mockSonarrService is a mock implementation of sonarrService for testing.
type mockSonarrService struct {
	healthy         bool
	status          sonarr.SyncStatus
	systemStatus    *sonarr.SystemStatus
	qualityProfiles []sonarr.QualityProfile
	rootFolders     []sonarr.RootFolder
	syncError       error
}

func (m *mockSonarrService) GetStatus() sonarr.SyncStatus {
	return m.status
}

func (m *mockSonarrService) IsHealthy(ctx context.Context) bool {
	return m.healthy
}

func (m *mockSonarrService) GetSystemStatus(ctx context.Context) (*sonarr.SystemStatus, error) {
	return m.systemStatus, nil
}

func (m *mockSonarrService) GetQualityProfiles(ctx context.Context) ([]sonarr.QualityProfile, error) {
	return m.qualityProfiles, nil
}

func (m *mockSonarrService) GetRootFolders(ctx context.Context) ([]sonarr.RootFolder, error) {
	return m.rootFolders, nil
}

func (m *mockSonarrService) SyncLibrary(ctx context.Context) (*sonarr.SyncResult, error) {
	if m.syncError != nil {
		return nil, m.syncError
	}
	return &sonarr.SyncResult{
		SeriesAdded:   5,
		SeriesUpdated: 10,
	}, nil
}

func (m *mockSonarrService) LookupSeries(_ context.Context, _ string) ([]sonarr.Series, error) {
	return nil, nil
}

// ============================================================================
// AdminGetSonarrStatus Tests
// ============================================================================

func TestHandler_AdminGetSonarrStatus_Unauthenticated(t *testing.T) {
	t.Parallel()
	handler, _, _ := setupRadarrTestHandler(t)

	ctx := context.Background()

	result, err := handler.AdminGetSonarrStatus(ctx)
	require.NoError(t, err)

	unauthorized, ok := result.(*ogen.AdminGetSonarrStatusUnauthorized)
	require.True(t, ok, "expected *ogen.AdminGetSonarrStatusUnauthorized, got %T", result)
	assert.Equal(t, 401, unauthorized.Code)
	assert.Contains(t, unauthorized.Message, "Authentication required")
}

func TestHandler_AdminGetSonarrStatus_NotConfigured(t *testing.T) {
	t.Parallel()
	handler, _, adminID := setupRadarrTestHandler(t)

	ctx := WithUserID(context.Background(), adminID)

	result, err := handler.AdminGetSonarrStatus(ctx)
	require.NoError(t, err)

	unavailable, ok := result.(*ogen.AdminGetSonarrStatusServiceUnavailable)
	require.True(t, ok)
	assert.Equal(t, 503, unavailable.Code)
	assert.Contains(t, unavailable.Message, "not configured")
}

func TestHandler_AdminGetSonarrStatus_Connected(t *testing.T) {
	t.Parallel()
	handler, _, adminID := setupRadarrTestHandler(t)

	mockService := &mockSonarrService{
		healthy: true,
		status: sonarr.SyncStatus{
			IsRunning:       false,
			LastSync:        time.Now().Add(-1 * time.Hour),
			SeriesAdded:     5,
			SeriesUpdated:   10,
			TotalSeries:     100,
			EpisodesAdded:   50,
			EpisodesUpdated: 20,
		},
		systemStatus: &sonarr.SystemStatus{
			Version:      "4.0.0",
			InstanceName: "Test Sonarr",
			StartTime:    time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
		},
	}
	handler.sonarrService = mockService

	ctx := WithUserID(context.Background(), adminID)

	result, err := handler.AdminGetSonarrStatus(ctx)
	require.NoError(t, err)

	status, ok := result.(*ogen.SonarrStatus)
	require.True(t, ok)
	assert.True(t, status.Connected)
	assert.Equal(t, "4.0.0", status.Version.Value)
	assert.Equal(t, "Test Sonarr", status.InstanceName.Value)
	assert.False(t, status.SyncStatus.IsRunning)
	assert.Equal(t, 5, status.SyncStatus.SeriesAdded)
	assert.Equal(t, 10, status.SyncStatus.SeriesUpdated)
	assert.Equal(t, 100, status.SyncStatus.TotalSeries)
	assert.Equal(t, 50, status.SyncStatus.EpisodesAdded.Value)
	assert.Equal(t, 20, status.SyncStatus.EpisodesUpdated.Value)
}

func TestHandler_AdminGetSonarrStatus_Disconnected(t *testing.T) {
	t.Parallel()
	handler, _, adminID := setupRadarrTestHandler(t)

	mockService := &mockSonarrService{
		healthy: false,
		status: sonarr.SyncStatus{
			IsRunning:     false,
			LastSyncError: "connection refused",
		},
	}
	handler.sonarrService = mockService

	ctx := WithUserID(context.Background(), adminID)

	result, err := handler.AdminGetSonarrStatus(ctx)
	require.NoError(t, err)

	status, ok := result.(*ogen.SonarrStatus)
	require.True(t, ok)
	assert.False(t, status.Connected)
	assert.Equal(t, "connection refused", status.SyncStatus.LastSyncError.Value)
}

// ============================================================================
// AdminTriggerSonarrSync Tests
// ============================================================================

func TestHandler_AdminTriggerSonarrSync_Unauthenticated(t *testing.T) {
	t.Parallel()
	handler, _, _ := setupRadarrTestHandler(t)

	ctx := context.Background()

	result, err := handler.AdminTriggerSonarrSync(ctx)
	require.NoError(t, err)

	unauthorized, ok := result.(*ogen.AdminTriggerSonarrSyncUnauthorized)
	require.True(t, ok, "expected *ogen.AdminTriggerSonarrSyncUnauthorized, got %T", result)
	assert.Equal(t, 401, unauthorized.Code)
}

func TestHandler_AdminTriggerSonarrSync_NotConfigured(t *testing.T) {
	t.Parallel()
	handler, _, adminID := setupRadarrTestHandler(t)

	ctx := WithUserID(context.Background(), adminID)

	result, err := handler.AdminTriggerSonarrSync(ctx)
	require.NoError(t, err)

	unavailable, ok := result.(*ogen.AdminTriggerSonarrSyncServiceUnavailable)
	require.True(t, ok)
	assert.Equal(t, 503, unavailable.Code)
}

func TestHandler_AdminTriggerSonarrSync_AlreadyRunning(t *testing.T) {
	t.Parallel()
	handler, _, adminID := setupRadarrTestHandler(t)

	mockService := &mockSonarrService{
		status: sonarr.SyncStatus{
			IsRunning: true,
		},
	}
	handler.sonarrService = mockService

	ctx := WithUserID(context.Background(), adminID)

	result, err := handler.AdminTriggerSonarrSync(ctx)
	require.NoError(t, err)

	conflict, ok := result.(*ogen.AdminTriggerSonarrSyncConflict)
	require.True(t, ok)
	assert.Equal(t, 409, conflict.Code)
	assert.Contains(t, conflict.Message, "already in progress")
}

func TestHandler_AdminTriggerSonarrSync_WithRiver(t *testing.T) {
	t.Parallel()
	handler, _, adminID := setupRadarrTestHandler(t)

	mockService := &mockSonarrService{
		status: sonarr.SyncStatus{IsRunning: false},
	}
	mockRiver := &mockRiverClient{}
	handler.sonarrService = mockService
	handler.riverClient = mockRiver

	ctx := WithUserID(context.Background(), adminID)

	result, err := handler.AdminTriggerSonarrSync(ctx)
	require.NoError(t, err)

	response, ok := result.(*ogen.SonarrSyncResponse)
	require.True(t, ok)
	assert.Equal(t, ogen.SonarrSyncResponseStatusQueued, response.Status)
	assert.Contains(t, response.Message, "queued")

	// Verify job was inserted
	require.Len(t, mockRiver.insertedArgs, 1)
	syncArgs, ok := mockRiver.insertedArgs[0].(*sonarr.SonarrSyncJobArgs)
	require.True(t, ok)
	assert.Equal(t, sonarr.SonarrSyncOperationFull, syncArgs.Operation)
}

func TestHandler_AdminTriggerSonarrSync_DirectSync(t *testing.T) {
	t.Parallel()
	handler, _, adminID := setupRadarrTestHandler(t)

	mockService := &mockSonarrService{
		status: sonarr.SyncStatus{IsRunning: false},
	}
	handler.sonarrService = mockService
	// No river client - returns 503 (job queue required)

	ctx := WithUserID(context.Background(), adminID)

	result, err := handler.AdminTriggerSonarrSync(ctx)
	require.NoError(t, err)

	response, ok := result.(*ogen.AdminTriggerSonarrSyncServiceUnavailable)
	require.True(t, ok)
	assert.Equal(t, 503, response.Code)
	assert.Contains(t, response.Message, "Job queue not available")
}

// ============================================================================
// AdminGetSonarrQualityProfiles Tests
// ============================================================================

func TestHandler_AdminGetSonarrQualityProfiles_Unauthenticated(t *testing.T) {
	t.Parallel()
	handler, _, _ := setupRadarrTestHandler(t)

	ctx := context.Background()

	result, err := handler.AdminGetSonarrQualityProfiles(ctx)
	require.NoError(t, err)

	unauthorized, ok := result.(*ogen.AdminGetSonarrQualityProfilesUnauthorized)
	require.True(t, ok, "expected *ogen.AdminGetSonarrQualityProfilesUnauthorized, got %T", result)
	assert.Equal(t, 401, unauthorized.Code)
}

func TestHandler_AdminGetSonarrQualityProfiles_Success(t *testing.T) {
	t.Parallel()
	handler, _, adminID := setupRadarrTestHandler(t)

	mockService := &mockSonarrService{
		qualityProfiles: []sonarr.QualityProfile{
			{ID: 1, Name: "HD-1080p", UpgradeAllowed: true, Cutoff: 7},
			{ID: 2, Name: "Any", UpgradeAllowed: false, Cutoff: 0},
		},
	}
	handler.sonarrService = mockService

	ctx := WithUserID(context.Background(), adminID)

	result, err := handler.AdminGetSonarrQualityProfiles(ctx)
	require.NoError(t, err)

	profiles, ok := result.(*ogen.SonarrQualityProfileList)
	require.True(t, ok)
	assert.Len(t, profiles.Profiles, 2)
	assert.Equal(t, 1, profiles.Profiles[0].ID)
	assert.Equal(t, "HD-1080p", profiles.Profiles[0].Name)
	assert.True(t, profiles.Profiles[0].UpgradeAllowed.Value)
}

// ============================================================================
// AdminGetSonarrRootFolders Tests
// ============================================================================

func TestHandler_AdminGetSonarrRootFolders_Unauthenticated(t *testing.T) {
	t.Parallel()
	handler, _, _ := setupRadarrTestHandler(t)

	ctx := context.Background()

	result, err := handler.AdminGetSonarrRootFolders(ctx)
	require.NoError(t, err)

	unauthorized, ok := result.(*ogen.AdminGetSonarrRootFoldersUnauthorized)
	require.True(t, ok, "expected *ogen.AdminGetSonarrRootFoldersUnauthorized, got %T", result)
	assert.Equal(t, 401, unauthorized.Code)
}

func TestHandler_AdminGetSonarrRootFolders_Success(t *testing.T) {
	t.Parallel()
	handler, _, adminID := setupRadarrTestHandler(t)

	mockService := &mockSonarrService{
		rootFolders: []sonarr.RootFolder{
			{ID: 1, Path: "/tv", Accessible: true, FreeSpace: 1000000000000},
			{ID: 2, Path: "/archive", Accessible: false, FreeSpace: 0},
		},
	}
	handler.sonarrService = mockService

	ctx := WithUserID(context.Background(), adminID)

	result, err := handler.AdminGetSonarrRootFolders(ctx)
	require.NoError(t, err)

	folders, ok := result.(*ogen.SonarrRootFolderList)
	require.True(t, ok)
	assert.Len(t, folders.Folders, 2)
	assert.Equal(t, 1, folders.Folders[0].ID)
	assert.Equal(t, "/tv", folders.Folders[0].Path)
	assert.True(t, folders.Folders[0].Accessible)
	assert.Equal(t, int64(1000000000000), folders.Folders[0].FreeSpace.Value)
	assert.False(t, folders.Folders[1].Accessible)
}

// ============================================================================
// HandleSonarrWebhook Tests
// ============================================================================

func TestHandler_HandleSonarrWebhook_Success(t *testing.T) {
	t.Parallel()
	handler, _, _ := setupRadarrTestHandler(t)

	mockRiver := &mockRiverClient{}
	handler.riverClient = mockRiver

	ctx := context.Background()
	payload := &ogen.SonarrWebhookPayload{
		EventType:    ogen.SonarrWebhookPayloadEventTypeDownload,
		InstanceName: ogen.NewOptString("Test Sonarr"),
	}
	payload.Series = ogen.NewOptSonarrWebhookSeries(ogen.SonarrWebhookSeries{
		ID:    ogen.NewOptInt(123),
		Title: ogen.NewOptString("Test Series"),
	})

	result, err := handler.HandleSonarrWebhook(ctx, payload)
	require.NoError(t, err)

	_, ok := result.(*ogen.HandleSonarrWebhookAccepted)
	require.True(t, ok)

	// Verify job was queued
	require.Len(t, mockRiver.insertedArgs, 1)
}

func TestHandler_HandleSonarrWebhook_NoRiver(t *testing.T) {
	t.Parallel()
	handler, _, _ := setupRadarrTestHandler(t)

	// No river client - should still accept webhook
	ctx := context.Background()
	payload := &ogen.SonarrWebhookPayload{
		EventType: ogen.SonarrWebhookPayloadEventTypeTest,
	}

	result, err := handler.HandleSonarrWebhook(ctx, payload)
	require.NoError(t, err)

	_, ok := result.(*ogen.HandleSonarrWebhookAccepted)
	require.True(t, ok)
}
