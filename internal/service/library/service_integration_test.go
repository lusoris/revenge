package library

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/service/activity"
	"github.com/lusoris/revenge/internal/testutil"
	"go.uber.org/zap"
)

// setupTestService creates a Service backed by a real PostgreSQL database.
func setupTestService(t *testing.T) (*Service, testutil.DB) {
	t.Helper()
	repo, testDB := setupTestRepository(t)
	logger := zap.NewNop()
	activityLogger := activity.NewNoopLogger()
	svc := NewService(repo, logger, activityLogger)
	return svc, testDB
}

// ============================================================================
// Library CRUD Integration Tests
// ============================================================================

func TestServiceIntegration_Create(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	t.Run("success with defaults", func(t *testing.T) {
		lib, err := svc.Create(ctx, CreateLibraryRequest{
			Name:    "Movies",
			Type:    LibraryTypeMovie,
			Paths:   []string{"/media/movies"},
			Enabled: true,
		})
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, lib.ID)
		assert.Equal(t, "Movies", lib.Name)
		assert.Equal(t, LibraryTypeMovie, lib.Type)
		assert.Equal(t, "en", lib.PreferredLanguage) // default
		assert.True(t, lib.Enabled)
		assert.False(t, lib.CreatedAt.IsZero())
		assert.False(t, lib.UpdatedAt.IsZero())
	})

	t.Run("success with custom preferred language", func(t *testing.T) {
		lib, err := svc.Create(ctx, CreateLibraryRequest{
			Name:              "Filme",
			Type:              LibraryTypeMovie,
			Paths:             []string{"/media/filme"},
			Enabled:           true,
			PreferredLanguage: "de",
		})
		require.NoError(t, err)
		assert.Equal(t, "de", lib.PreferredLanguage)
	})

	t.Run("success with scanner config", func(t *testing.T) {
		lib, err := svc.Create(ctx, CreateLibraryRequest{
			Name:    "Music Collection",
			Type:    LibraryTypeMusic,
			Paths:   []string{"/media/music"},
			Enabled: true,
			ScannerConfig: map[string]interface{}{
				"skip_hidden": true,
				"depth":       float64(5),
			},
		})
		require.NoError(t, err)
		assert.NotNil(t, lib.ScannerConfig)
		assert.Equal(t, true, lib.ScannerConfig["skip_hidden"])
	})

	t.Run("invalid library type", func(t *testing.T) {
		lib, err := svc.Create(ctx, CreateLibraryRequest{
			Name:    "Invalid",
			Type:    "invalid_type",
			Paths:   []string{"/media/invalid"},
			Enabled: true,
		})
		assert.Nil(t, lib)
		assert.ErrorIs(t, err, ErrInvalidLibraryType)
	})

	t.Run("duplicate name", func(t *testing.T) {
		_, err := svc.Create(ctx, CreateLibraryRequest{
			Name:    "Duplicate",
			Type:    LibraryTypeTVShow,
			Paths:   []string{"/media/tv1"},
			Enabled: true,
		})
		require.NoError(t, err)

		lib, err := svc.Create(ctx, CreateLibraryRequest{
			Name:    "Duplicate",
			Type:    LibraryTypeTVShow,
			Paths:   []string{"/media/tv2"},
			Enabled: true,
		})
		assert.Nil(t, lib)
		assert.ErrorIs(t, err, ErrLibraryExists)
	})
}

func TestServiceIntegration_CRUD_Lifecycle(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Create
	lib, err := svc.Create(ctx, CreateLibraryRequest{
		Name:    "CRUD Test Library",
		Type:    LibraryTypeMovie,
		Paths:   []string{"/media/crud"},
		Enabled: true,
	})
	require.NoError(t, err)
	libID := lib.ID

	// Get by ID
	retrieved, err := svc.Get(ctx, libID)
	require.NoError(t, err)
	assert.Equal(t, "CRUD Test Library", retrieved.Name)
	assert.Equal(t, LibraryTypeMovie, retrieved.Type)

	// Get by Name
	byName, err := svc.GetByName(ctx, "CRUD Test Library")
	require.NoError(t, err)
	assert.Equal(t, libID, byName.ID)

	// Update name
	newName := "Updated CRUD Library"
	updated, err := svc.Update(ctx, libID, &LibraryUpdate{
		Name: &newName,
	})
	require.NoError(t, err)
	assert.Equal(t, newName, updated.Name)
	assert.True(t, updated.UpdatedAt.After(lib.UpdatedAt) || updated.UpdatedAt.Equal(lib.UpdatedAt))

	// Update enabled
	disabled := false
	updated2, err := svc.Update(ctx, libID, &LibraryUpdate{
		Enabled: &disabled,
	})
	require.NoError(t, err)
	assert.False(t, updated2.Enabled)

	// List all
	allLibs, err := svc.List(ctx)
	require.NoError(t, err)
	found := false
	for _, l := range allLibs {
		if l.ID == libID {
			found = true
			assert.Equal(t, newName, l.Name)
		}
	}
	assert.True(t, found, "Library should appear in List")

	// Count
	count, err := svc.Count(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(1))

	// Delete
	err = svc.Delete(ctx, libID)
	require.NoError(t, err)

	// Verify deletion
	_, err = svc.Get(ctx, libID)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestServiceIntegration_Update_Errors(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Create two libraries
	lib1, err := svc.Create(ctx, CreateLibraryRequest{
		Name:    "Update Error Lib 1",
		Type:    LibraryTypeMovie,
		Paths:   []string{"/media/err1"},
		Enabled: true,
	})
	require.NoError(t, err)

	lib2, err := svc.Create(ctx, CreateLibraryRequest{
		Name:    "Update Error Lib 2",
		Type:    LibraryTypeTVShow,
		Paths:   []string{"/media/err2"},
		Enabled: true,
	})
	require.NoError(t, err)

	// Update lib1 to name of lib2
	existingName := lib2.Name
	_, err = svc.Update(ctx, lib1.ID, &LibraryUpdate{
		Name: &existingName,
	})
	assert.ErrorIs(t, err, ErrLibraryExists)

	// Invalid type in update
	invalidType := "invalid"
	_, err = svc.Update(ctx, lib1.ID, &LibraryUpdate{
		Type: &invalidType,
	})
	assert.ErrorIs(t, err, ErrInvalidLibraryType)

	// Update same name to itself should succeed
	sameName := lib1.Name
	updated, err := svc.Update(ctx, lib1.ID, &LibraryUpdate{
		Name: &sameName,
	})
	require.NoError(t, err)
	assert.Equal(t, sameName, updated.Name)

	// Update non-existent library
	_, err = svc.Update(ctx, uuid.Must(uuid.NewV7()), &LibraryUpdate{
		Name: stringPtr("ghost"),
	})
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestServiceIntegration_ListByType(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Create several libraries of different types
	_, err := svc.Create(ctx, CreateLibraryRequest{
		Name: "Type Movies 1", Type: LibraryTypeMovie,
		Paths: []string{"/m1"}, Enabled: true,
	})
	require.NoError(t, err)
	_, err = svc.Create(ctx, CreateLibraryRequest{
		Name: "Type Movies 2", Type: LibraryTypeMovie,
		Paths: []string{"/m2"}, Enabled: true,
	})
	require.NoError(t, err)
	_, err = svc.Create(ctx, CreateLibraryRequest{
		Name: "Type TV", Type: LibraryTypeTVShow,
		Paths: []string{"/tv"}, Enabled: true,
	})
	require.NoError(t, err)

	movies, err := svc.ListByType(ctx, LibraryTypeMovie)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(movies), 2)
	for _, m := range movies {
		assert.Equal(t, LibraryTypeMovie, m.Type)
	}

	tvShows, err := svc.ListByType(ctx, LibraryTypeTVShow)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(tvShows), 1)

	// Invalid type
	_, err = svc.ListByType(ctx, "invalid")
	assert.ErrorIs(t, err, ErrInvalidLibraryType)
}

func TestServiceIntegration_ListEnabled(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Create enabled and disabled libraries
	_, err := svc.Create(ctx, CreateLibraryRequest{
		Name: "Enabled Lib", Type: LibraryTypeMovie,
		Paths: []string{"/en"}, Enabled: true,
	})
	require.NoError(t, err)
	_, err = svc.Create(ctx, CreateLibraryRequest{
		Name: "Disabled Lib", Type: LibraryTypeMovie,
		Paths: []string{"/dis"}, Enabled: false,
	})
	require.NoError(t, err)

	enabled, err := svc.ListEnabled(ctx)
	require.NoError(t, err)
	for _, l := range enabled {
		assert.True(t, l.Enabled, "ListEnabled should only return enabled libraries")
	}
}

func TestServiceIntegration_Get_NotFound(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	_, err := svc.Get(ctx, uuid.Must(uuid.NewV7()))
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestServiceIntegration_GetByName_NotFound(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	_, err := svc.GetByName(ctx, "nonexistent_library_name")
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestServiceIntegration_Delete_CleansUpPermissions(t *testing.T) {
	t.Parallel()
	svc, testDB := setupTestService(t)
	ctx := context.Background()

	lib, err := svc.Create(ctx, CreateLibraryRequest{
		Name: "Delete Perms Lib", Type: LibraryTypeMovie,
		Paths: []string{"/del"}, Enabled: true,
	})
	require.NoError(t, err)

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "del_perm_user",
		Email:    "del_perm@example.com",
	})

	// Grant permission
	err = svc.GrantPermission(ctx, lib.ID, user.ID, PermissionView)
	require.NoError(t, err)

	// Verify permission exists
	has, err := svc.CheckPermission(ctx, lib.ID, user.ID, PermissionView)
	require.NoError(t, err)
	assert.True(t, has)

	// Delete library (should also revoke permissions)
	err = svc.Delete(ctx, lib.ID)
	require.NoError(t, err)

	// Permissions should be revoked (library no longer exists, so permission check would fail)
	// Verify by listing user permissions -- the deleted library's permissions should be gone
	perms, err := svc.ListUserPermissions(ctx, user.ID)
	require.NoError(t, err)
	for _, p := range perms {
		assert.NotEqual(t, lib.ID, p.LibraryID, "Permissions for deleted library should be gone")
	}
}

// ============================================================================
// Scan Lifecycle Integration Tests
// ============================================================================

func TestServiceIntegration_ScanLifecycle_FullSuccess(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Create library
	lib, err := svc.Create(ctx, CreateLibraryRequest{
		Name: "Scan Lifecycle Lib", Type: LibraryTypeMovie,
		Paths: []string{"/scan"}, Enabled: true,
	})
	require.NoError(t, err)

	// TriggerScan
	scan, err := svc.TriggerScan(ctx, lib.ID, ScanTypeFull)
	require.NoError(t, err)
	assert.Equal(t, ScanStatusPending, scan.Status)
	assert.Equal(t, lib.ID, scan.LibraryID)
	assert.Equal(t, ScanTypeFull, scan.ScanType)

	// StartScan
	started, err := svc.StartScan(ctx, scan.ID)
	require.NoError(t, err)
	assert.Equal(t, ScanStatusRunning, started.Status)
	assert.NotNil(t, started.StartedAt)

	// UpdateScanProgress
	progress := &ScanProgress{
		ItemsScanned: 50,
		ItemsAdded:   30,
		ItemsUpdated: 10,
		ItemsRemoved: 5,
		ErrorsCount:  2,
	}
	progressed, err := svc.UpdateScanProgress(ctx, scan.ID, progress)
	require.NoError(t, err)
	assert.Equal(t, int32(50), progressed.ItemsScanned)
	assert.Equal(t, int32(30), progressed.ItemsAdded)
	assert.Equal(t, int32(10), progressed.ItemsUpdated)
	assert.Equal(t, int32(5), progressed.ItemsRemoved)
	assert.Equal(t, int32(2), progressed.ErrorsCount)

	// CompleteScan with final progress
	finalProgress := &ScanProgress{
		ItemsScanned: 100,
		ItemsAdded:   60,
		ItemsUpdated: 20,
		ItemsRemoved: 10,
		ErrorsCount:  3,
	}
	completed, err := svc.CompleteScan(ctx, scan.ID, finalProgress)
	require.NoError(t, err)
	assert.Equal(t, ScanStatusCompleted, completed.Status)
	assert.NotNil(t, completed.CompletedAt)
	assert.NotNil(t, completed.DurationSeconds)

	// GetScan to verify final state
	final, err := svc.GetScan(ctx, scan.ID)
	require.NoError(t, err)
	assert.Equal(t, ScanStatusCompleted, final.Status)
	assert.Equal(t, int32(100), final.ItemsScanned)
	assert.Equal(t, int32(60), final.ItemsAdded)

	// GetLatestScan should return this scan
	latest, err := svc.GetLatestScan(ctx, lib.ID)
	require.NoError(t, err)
	assert.Equal(t, scan.ID, latest.ID)
}

func TestServiceIntegration_ScanLifecycle_Failure(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	lib, err := svc.Create(ctx, CreateLibraryRequest{
		Name: "Scan Fail Lib", Type: LibraryTypeMovie,
		Paths: []string{"/scanfail"}, Enabled: true,
	})
	require.NoError(t, err)

	// TriggerScan
	scan, err := svc.TriggerScan(ctx, lib.ID, ScanTypeIncremental)
	require.NoError(t, err)

	// StartScan
	_, err = svc.StartScan(ctx, scan.ID)
	require.NoError(t, err)

	// FailScan
	failed, err := svc.FailScan(ctx, scan.ID, "disk I/O error")
	require.NoError(t, err)
	assert.Equal(t, ScanStatusFailed, failed.Status)
	assert.NotNil(t, failed.CompletedAt)
	assert.NotNil(t, failed.DurationSeconds)

	// Verify error message is stored
	failedScan, err := svc.GetScan(ctx, scan.ID)
	require.NoError(t, err)
	require.NotNil(t, failedScan.ErrorMessage)
	assert.Equal(t, "disk I/O error", *failedScan.ErrorMessage)
}

func TestServiceIntegration_ScanLifecycle_Cancellation(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	lib, err := svc.Create(ctx, CreateLibraryRequest{
		Name: "Scan Cancel Lib", Type: LibraryTypeMovie,
		Paths: []string{"/scancancel"}, Enabled: true,
	})
	require.NoError(t, err)

	// TriggerScan
	scan, err := svc.TriggerScan(ctx, lib.ID, ScanTypeMetadata)
	require.NoError(t, err)

	// CancelScan (before starting)
	cancelled, err := svc.CancelScan(ctx, scan.ID)
	require.NoError(t, err)
	assert.Equal(t, ScanStatusCancelled, cancelled.Status)
	assert.NotNil(t, cancelled.CompletedAt)
}

func TestServiceIntegration_CompleteScan_WithoutProgress(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	lib, err := svc.Create(ctx, CreateLibraryRequest{
		Name: "Scan No Progress Lib", Type: LibraryTypeMovie,
		Paths: []string{"/scannoprog"}, Enabled: true,
	})
	require.NoError(t, err)

	scan, err := svc.TriggerScan(ctx, lib.ID, ScanTypeFull)
	require.NoError(t, err)

	_, err = svc.StartScan(ctx, scan.ID)
	require.NoError(t, err)

	// CompleteScan without progress (nil)
	completed, err := svc.CompleteScan(ctx, scan.ID, nil)
	require.NoError(t, err)
	assert.Equal(t, ScanStatusCompleted, completed.Status)
}

func TestServiceIntegration_TriggerScan_ScanAlreadyInProgress(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	lib, err := svc.Create(ctx, CreateLibraryRequest{
		Name: "Scan In Progress Lib", Type: LibraryTypeMovie,
		Paths: []string{"/scanprog"}, Enabled: true,
	})
	require.NoError(t, err)

	// Trigger first scan
	scan1, err := svc.TriggerScan(ctx, lib.ID, ScanTypeFull)
	require.NoError(t, err)

	// Start the first scan (now it's running)
	_, err = svc.StartScan(ctx, scan1.ID)
	require.NoError(t, err)

	// Trigger second scan should fail
	_, err = svc.TriggerScan(ctx, lib.ID, ScanTypeFull)
	assert.ErrorIs(t, err, ErrScanInProgress)

	// Complete the first scan
	_, err = svc.CompleteScan(ctx, scan1.ID, nil)
	require.NoError(t, err)

	// Now triggering a new scan should succeed
	scan2, err := svc.TriggerScan(ctx, lib.ID, ScanTypeFull)
	require.NoError(t, err)
	assert.NotEqual(t, scan1.ID, scan2.ID)
}

func TestServiceIntegration_TriggerScan_InvalidType(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	lib, err := svc.Create(ctx, CreateLibraryRequest{
		Name: "Scan Invalid Type Lib", Type: LibraryTypeMovie,
		Paths: []string{"/scaninv"}, Enabled: true,
	})
	require.NoError(t, err)

	_, err = svc.TriggerScan(ctx, lib.ID, "invalid_scan_type")
	assert.ErrorIs(t, err, ErrInvalidScanType)
}

func TestServiceIntegration_TriggerScan_NonExistentLibrary(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	_, err := svc.TriggerScan(ctx, uuid.Must(uuid.NewV7()), ScanTypeFull)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestServiceIntegration_ListScans(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	lib, err := svc.Create(ctx, CreateLibraryRequest{
		Name: "List Scans Lib", Type: LibraryTypeMovie,
		Paths: []string{"/listscans"}, Enabled: true,
	})
	require.NoError(t, err)

	// Create several scans
	for i := 0; i < 5; i++ {
		scan, err := svc.TriggerScan(ctx, lib.ID, ScanTypeFull)
		require.NoError(t, err)
		// Complete each so the next can be triggered
		_, err = svc.CompleteScan(ctx, scan.ID, nil)
		require.NoError(t, err)
	}

	// Test default limit
	scans, count, err := svc.ListScans(ctx, lib.ID, 0, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(5), count)
	assert.Len(t, scans, 5)

	// Test with limit
	scans, count, err = svc.ListScans(ctx, lib.ID, 2, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(5), count)
	assert.Len(t, scans, 2)

	// Test with offset
	scans, count, err = svc.ListScans(ctx, lib.ID, 3, 3)
	require.NoError(t, err)
	assert.Equal(t, int64(5), count)
	assert.Len(t, scans, 2) // Only 2 remain after offset 3

	// Test limit capped at 100
	scans, _, err = svc.ListScans(ctx, lib.ID, 200, 0)
	require.NoError(t, err)
	assert.LessOrEqual(t, len(scans), 100)
}

func TestServiceIntegration_GetRunningScans(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	lib, err := svc.Create(ctx, CreateLibraryRequest{
		Name: "Running Scans Lib", Type: LibraryTypeMovie,
		Paths: []string{"/running"}, Enabled: true,
	})
	require.NoError(t, err)

	// No running scans initially
	running, err := svc.GetRunningScans(ctx)
	require.NoError(t, err)
	initialRunning := len(running)

	// Trigger and start a scan
	scan, err := svc.TriggerScan(ctx, lib.ID, ScanTypeFull)
	require.NoError(t, err)
	_, err = svc.StartScan(ctx, scan.ID)
	require.NoError(t, err)

	// Now there should be a running scan
	running, err = svc.GetRunningScans(ctx)
	require.NoError(t, err)
	assert.Equal(t, initialRunning+1, len(running))

	foundOur := false
	for _, s := range running {
		if s.ID == scan.ID {
			foundOur = true
			assert.Equal(t, ScanStatusRunning, s.Status)
		}
	}
	assert.True(t, foundOur)

	// Complete the scan
	_, err = svc.CompleteScan(ctx, scan.ID, nil)
	require.NoError(t, err)

	// Should be gone from running
	running, err = svc.GetRunningScans(ctx)
	require.NoError(t, err)
	for _, s := range running {
		assert.NotEqual(t, scan.ID, s.ID)
	}
}

// ============================================================================
// Permission Integration Tests
// ============================================================================

func TestServiceIntegration_PermissionLifecycle(t *testing.T) {
	t.Parallel()
	svc, testDB := setupTestService(t)
	ctx := context.Background()

	lib, err := svc.Create(ctx, CreateLibraryRequest{
		Name: "Perm Lifecycle Lib", Type: LibraryTypeMovie,
		Paths: []string{"/perm"}, Enabled: true,
	})
	require.NoError(t, err)

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "perm_lifecycle_user",
		Email:    "perm_lifecycle@example.com",
	})

	// GrantPermission (view)
	err = svc.GrantPermission(ctx, lib.ID, user.ID, PermissionView)
	require.NoError(t, err)

	// CheckPermission - should be true for view
	has, err := svc.CheckPermission(ctx, lib.ID, user.ID, PermissionView)
	require.NoError(t, err)
	assert.True(t, has)

	// CheckPermission - should be false for download (not granted)
	has, err = svc.CheckPermission(ctx, lib.ID, user.ID, PermissionDownload)
	require.NoError(t, err)
	assert.False(t, has)

	// GetPermission
	perm, err := svc.GetPermission(ctx, lib.ID, user.ID, PermissionView)
	require.NoError(t, err)
	assert.Equal(t, PermissionView, perm.Permission)
	assert.Equal(t, lib.ID, perm.LibraryID)
	assert.Equal(t, user.ID, perm.UserID)

	// Grant more permissions
	err = svc.GrantPermission(ctx, lib.ID, user.ID, PermissionDownload)
	require.NoError(t, err)
	err = svc.GrantPermission(ctx, lib.ID, user.ID, PermissionManage)
	require.NoError(t, err)

	// ListPermissions for library
	perms, err := svc.ListPermissions(ctx, lib.ID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(perms), 3)

	// ListUserPermissions
	userPerms, err := svc.ListUserPermissions(ctx, user.ID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(userPerms), 3)

	// RevokePermission
	err = svc.RevokePermission(ctx, lib.ID, user.ID, PermissionView)
	require.NoError(t, err)

	// Verify revoked
	has, err = svc.CheckPermission(ctx, lib.ID, user.ID, PermissionView)
	require.NoError(t, err)
	assert.False(t, has)

	// Other permissions should remain
	has, err = svc.CheckPermission(ctx, lib.ID, user.ID, PermissionDownload)
	require.NoError(t, err)
	assert.True(t, has)
}

func TestServiceIntegration_CanAccess_CanDownload_CanManage(t *testing.T) {
	t.Parallel()
	svc, testDB := setupTestService(t)
	ctx := context.Background()

	lib, err := svc.Create(ctx, CreateLibraryRequest{
		Name: "Can Methods Lib", Type: LibraryTypeMovie,
		Paths: []string{"/can"}, Enabled: true,
	})
	require.NoError(t, err)

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "can_user",
		Email:    "can@example.com",
	})

	// Admin always has access
	can, err := svc.CanAccess(ctx, lib.ID, user.ID, true)
	require.NoError(t, err)
	assert.True(t, can)

	can, err = svc.CanDownload(ctx, lib.ID, user.ID, true)
	require.NoError(t, err)
	assert.True(t, can)

	can, err = svc.CanManage(ctx, lib.ID, user.ID, true)
	require.NoError(t, err)
	assert.True(t, can)

	// Non-admin without permissions
	can, err = svc.CanAccess(ctx, lib.ID, user.ID, false)
	require.NoError(t, err)
	assert.False(t, can)

	can, err = svc.CanDownload(ctx, lib.ID, user.ID, false)
	require.NoError(t, err)
	assert.False(t, can)

	can, err = svc.CanManage(ctx, lib.ID, user.ID, false)
	require.NoError(t, err)
	assert.False(t, can)

	// Grant view only
	err = svc.GrantPermission(ctx, lib.ID, user.ID, PermissionView)
	require.NoError(t, err)

	can, err = svc.CanAccess(ctx, lib.ID, user.ID, false)
	require.NoError(t, err)
	assert.True(t, can)

	can, err = svc.CanDownload(ctx, lib.ID, user.ID, false)
	require.NoError(t, err)
	assert.False(t, can) // only view, not download

	// Grant download
	err = svc.GrantPermission(ctx, lib.ID, user.ID, PermissionDownload)
	require.NoError(t, err)

	can, err = svc.CanDownload(ctx, lib.ID, user.ID, false)
	require.NoError(t, err)
	assert.True(t, can)

	// Grant manage
	err = svc.GrantPermission(ctx, lib.ID, user.ID, PermissionManage)
	require.NoError(t, err)

	can, err = svc.CanManage(ctx, lib.ID, user.ID, false)
	require.NoError(t, err)
	assert.True(t, can)
}

func TestServiceIntegration_ListAccessible_MixedPermissions(t *testing.T) {
	t.Parallel()
	svc, testDB := setupTestService(t)
	ctx := context.Background()

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "accessible_user",
		Email:    "accessible@example.com",
	})

	// Create libraries
	lib1, err := svc.Create(ctx, CreateLibraryRequest{
		Name: "Accessible Lib 1", Type: LibraryTypeMovie,
		Paths: []string{"/a1"}, Enabled: true,
	})
	require.NoError(t, err)

	lib2, err := svc.Create(ctx, CreateLibraryRequest{
		Name: "Accessible Lib 2", Type: LibraryTypeTVShow,
		Paths: []string{"/a2"}, Enabled: true,
	})
	require.NoError(t, err)

	_, err = svc.Create(ctx, CreateLibraryRequest{
		Name: "Inaccessible Lib", Type: LibraryTypeMusic,
		Paths: []string{"/a3"}, Enabled: true,
	})
	require.NoError(t, err)

	// Grant view on lib1 and lib2 only
	err = svc.GrantPermission(ctx, lib1.ID, user.ID, PermissionView)
	require.NoError(t, err)
	err = svc.GrantPermission(ctx, lib2.ID, user.ID, PermissionView)
	require.NoError(t, err)

	// ListAccessible should return only the 2 accessible libraries
	accessible, err := svc.ListAccessible(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, 2, len(accessible))

	accessibleIDs := make(map[uuid.UUID]bool)
	for _, l := range accessible {
		accessibleIDs[l.ID] = true
	}
	assert.True(t, accessibleIDs[lib1.ID])
	assert.True(t, accessibleIDs[lib2.ID])
}

func TestServiceIntegration_GrantPermission_Errors(t *testing.T) {
	t.Parallel()
	svc, testDB := setupTestService(t)
	ctx := context.Background()

	lib, err := svc.Create(ctx, CreateLibraryRequest{
		Name: "Grant Error Lib", Type: LibraryTypeMovie,
		Paths: []string{"/ge"}, Enabled: true,
	})
	require.NoError(t, err)

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "grant_err_user",
		Email:    "grant_err@example.com",
	})

	// Invalid permission
	err = svc.GrantPermission(ctx, lib.ID, user.ID, "invalid_permission")
	assert.ErrorIs(t, err, ErrInvalidPermission)

	// Non-existent library
	err = svc.GrantPermission(ctx, uuid.Must(uuid.NewV7()), user.ID, PermissionView)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestServiceIntegration_CheckPermission_InvalidPermission(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	has, err := svc.CheckPermission(ctx, uuid.Must(uuid.NewV7()), uuid.Must(uuid.NewV7()), "invalid")
	assert.False(t, has)
	assert.ErrorIs(t, err, ErrInvalidPermission)
}

func TestServiceIntegration_RevokePermission_InvalidPermission(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	err := svc.RevokePermission(ctx, uuid.Must(uuid.NewV7()), uuid.Must(uuid.NewV7()), "invalid")
	assert.ErrorIs(t, err, ErrInvalidPermission)
}

// ============================================================================
// NewService constructor test
// ============================================================================

func TestServiceIntegration_NewService(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	logger := zap.NewNop()

	// With noop activity logger
	svc := NewService(repo, logger, activity.NewNoopLogger())
	assert.NotNil(t, svc)

	// Verify it works
	ctx := context.Background()
	count, err := svc.Count(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(0))
}
