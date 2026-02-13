package library

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/testutil"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testutil.StopSharedPostgres()
	os.Exit(code)
}

func setupTestRepository(t *testing.T) (*RepositoryPg, testutil.DB) {
	t.Helper()
	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())
	repo := NewRepositoryPg(queries)
	return repo, testDB
}

func createTestLibrary(name, libType string) *Library {
	return &Library{
		Name:              name,
		Type:              libType,
		Paths:             []string{"/media/" + name},
		Enabled:           true,
		ScanOnStartup:     false,
		PreferredLanguage: "en",
	}
}

// ============================================================================
// Library CRUD Tests
// ============================================================================

func TestRepositoryPg_CreateLibrary(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	lib := createTestLibrary("test_library", "movies")
	err := repo.Create(ctx, lib)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, lib.ID)
	assert.False(t, lib.CreatedAt.IsZero())
	assert.False(t, lib.UpdatedAt.IsZero())
}

func TestRepositoryPg_CreateLibrary_Minimal(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	lib := &Library{
		Name:              "minimal",
		Type:              "tv",
		Paths:             []string{"/tv"},
		Enabled:           true,
		PreferredLanguage: "en",
	}

	err := repo.Create(ctx, lib)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, lib.ID)
}

func TestRepositoryPg_Get(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	lib := createTestLibrary("get_test", "movies")
	require.NoError(t, repo.Create(ctx, lib))

	retrieved, err := repo.Get(ctx, lib.ID)
	require.NoError(t, err)
	assert.Equal(t, lib.ID, retrieved.ID)
	assert.Equal(t, lib.Name, retrieved.Name)
	assert.Equal(t, lib.Type, retrieved.Type)
	assert.Equal(t, lib.Paths, retrieved.Paths)
}

func TestRepositoryPg_Get_NotFound(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	_, err := repo.Get(ctx, uuid.Must(uuid.NewV7()))
	assert.Error(t, err)
}

func TestRepositoryPg_GetByName(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	lib := createTestLibrary("unique_name", "movies")
	require.NoError(t, repo.Create(ctx, lib))

	retrieved, err := repo.GetByName(ctx, "unique_name")
	require.NoError(t, err)
	assert.Equal(t, lib.ID, retrieved.ID)
	assert.Equal(t, "unique_name", retrieved.Name)
}

func TestRepositoryPg_GetByName_NotFound(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	_, err := repo.GetByName(ctx, "nonexistent")
	assert.Error(t, err)
}

func TestRepositoryPg_List(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	// Create multiple libraries
	for i := range 3 {
		lib := createTestLibrary("list_test_"+string(rune('a'+i)), "movies")
		require.NoError(t, repo.Create(ctx, lib))
	}

	libraries, err := repo.List(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(libraries), 3)
}

func TestRepositoryPg_ListEnabled(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	// Create enabled library
	enabled := createTestLibrary("enabled", "movies")
	enabled.Enabled = true
	require.NoError(t, repo.Create(ctx, enabled))

	// Create disabled library
	disabled := createTestLibrary("disabled", "tv")
	disabled.Enabled = false
	require.NoError(t, repo.Create(ctx, disabled))

	libraries, err := repo.ListEnabled(ctx)
	require.NoError(t, err)

	// Check all returned libraries are enabled
	for _, lib := range libraries {
		assert.True(t, lib.Enabled)
	}
}

func TestRepositoryPg_ListByType(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	// Create movies libraries
	for i := range 2 {
		lib := createTestLibrary("movie_"+string(rune('a'+i)), "movies")
		require.NoError(t, repo.Create(ctx, lib))
	}

	// Create tv library
	tv := createTestLibrary("tv_show", "tv")
	require.NoError(t, repo.Create(ctx, tv))

	movies, err := repo.ListByType(ctx, "movies")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(movies), 2)
	for _, lib := range movies {
		assert.Equal(t, "movies", lib.Type)
	}
}

func TestRepositoryPg_Update(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	lib := createTestLibrary("update_test", "movies")
	require.NoError(t, repo.Create(ctx, lib))

	update := &LibraryUpdate{
		Name:    new("updated_name"),
		Enabled: new(false),
	}

	updated, err := repo.Update(ctx, lib.ID, update)
	require.NoError(t, err)
	assert.Equal(t, "updated_name", updated.Name)
	assert.False(t, updated.Enabled)
	assert.True(t, updated.UpdatedAt.After(lib.UpdatedAt))
}

func TestRepositoryPg_Delete(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	lib := createTestLibrary("delete_test", "movies")
	require.NoError(t, repo.Create(ctx, lib))

	err := repo.Delete(ctx, lib.ID)
	require.NoError(t, err)

	_, err = repo.Get(ctx, lib.ID)
	assert.Error(t, err)
}

func TestRepositoryPg_Count(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	initialCount, err := repo.Count(ctx)
	require.NoError(t, err)

	// Create libraries
	for i := range 3 {
		lib := createTestLibrary("count_"+string(rune('a'+i)), "movies")
		require.NoError(t, repo.Create(ctx, lib))
	}

	count, err := repo.Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, initialCount+3, count)
}

func TestRepositoryPg_CountByType(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	// Create music libraries
	for i := range 2 {
		lib := createTestLibrary("music_"+string(rune('a'+i)), "music")
		require.NoError(t, repo.Create(ctx, lib))
	}

	count, err := repo.CountByType(ctx, "music")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(2))
}

// ============================================================================
// Library Scan Tests
// ============================================================================

func TestRepositoryPg_CreateScan(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	lib := createTestLibrary("scan_lib", "movies")
	require.NoError(t, repo.Create(ctx, lib))

	scan := &LibraryScan{
		LibraryID: lib.ID,
		ScanType:  ScanTypeFull,
		Status:    "running",
	}

	err := repo.CreateScan(ctx, scan)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, scan.ID)
	assert.False(t, scan.CreatedAt.IsZero())
}

func TestRepositoryPg_GetScan(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	lib := createTestLibrary("get_scan_lib", "movies")
	require.NoError(t, repo.Create(ctx, lib))

	scan := &LibraryScan{
		LibraryID: lib.ID,
		ScanType:  ScanTypeFull,
		Status:    "running",
	}
	require.NoError(t, repo.CreateScan(ctx, scan))

	retrieved, err := repo.GetScan(ctx, scan.ID)
	require.NoError(t, err)
	assert.Equal(t, scan.ID, retrieved.ID)
	assert.Equal(t, lib.ID, retrieved.LibraryID)
}

func TestRepositoryPg_ListScans(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	lib := createTestLibrary("list_scans_lib", "movies")
	require.NoError(t, repo.Create(ctx, lib))

	// Create multiple scans
	for range 3 {
		scan := &LibraryScan{
			LibraryID: lib.ID,
			ScanType:  ScanTypeFull,
			Status:    "completed",
		}
		require.NoError(t, repo.CreateScan(ctx, scan))
	}

	scans, err := repo.ListScans(ctx, lib.ID, 10, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(scans), 3)
}

func TestRepositoryPg_CountScans(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	lib := createTestLibrary("count_scans_lib", "movies")
	require.NoError(t, repo.Create(ctx, lib))

	// Create scans
	for range 2 {
		scan := &LibraryScan{
			LibraryID: lib.ID,
			ScanType:  ScanTypeFull,
			Status:    "completed",
		}
		require.NoError(t, repo.CreateScan(ctx, scan))
	}

	count, err := repo.CountScans(ctx, lib.ID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(2))
}

func TestRepositoryPg_GetLatestScan(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	lib := createTestLibrary("latest_scan_lib", "movies")
	require.NoError(t, repo.Create(ctx, lib))

	// Create scans
	var lastScan *LibraryScan
	for range 3 {
		scan := &LibraryScan{
			LibraryID: lib.ID,
			ScanType:  ScanTypeFull,
			Status:    "completed",
		}
		require.NoError(t, repo.CreateScan(ctx, scan))
		lastScan = scan
		time.Sleep(10 * time.Millisecond) // Ensure time difference
	}

	latest, err := repo.GetLatestScan(ctx, lib.ID)
	require.NoError(t, err)
	assert.Equal(t, lastScan.ID, latest.ID)
}

func TestRepositoryPg_GetRunningScans(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	lib := createTestLibrary("running_scans_lib", "movies")
	require.NoError(t, repo.Create(ctx, lib))

	// Create running scan
	running := &LibraryScan{
		LibraryID: lib.ID,
		ScanType:  ScanTypeFull,
		Status:    "running",
	}
	require.NoError(t, repo.CreateScan(ctx, running))

	// Create completed scan
	completed := &LibraryScan{
		LibraryID: lib.ID,
		ScanType:  ScanTypeFull,
		Status:    "completed",
	}
	require.NoError(t, repo.CreateScan(ctx, completed))

	runningScans, err := repo.GetRunningScans(ctx)
	require.NoError(t, err)

	// Check all returned scans are running
	foundRunning := false
	for _, scan := range runningScans {
		assert.Equal(t, "running", scan.Status)
		if scan.ID == running.ID {
			foundRunning = true
		}
	}
	assert.True(t, foundRunning)
}

func TestRepositoryPg_UpdateScanStatus(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	lib := createTestLibrary("update_status_lib", "movies")
	require.NoError(t, repo.Create(ctx, lib))

	scan := &LibraryScan{
		LibraryID: lib.ID,
		ScanType:  ScanTypeFull,
		Status:    "running",
	}
	require.NoError(t, repo.CreateScan(ctx, scan))

	now := time.Now()
	update := &ScanStatusUpdate{
		Status:       "completed",
		CompletedAt:  &now,
		ErrorMessage: new(""),
	}

	updated, err := repo.UpdateScanStatus(ctx, scan.ID, update)
	require.NoError(t, err)
	assert.Equal(t, "completed", updated.Status)
	assert.NotNil(t, updated.CompletedAt)
}

func TestRepositoryPg_UpdateScanProgress(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	lib := createTestLibrary("progress_lib", "movies")
	require.NoError(t, repo.Create(ctx, lib))

	scan := &LibraryScan{
		LibraryID: lib.ID,
		ScanType:  ScanTypeFull,
		Status:    "running",
	}
	require.NoError(t, repo.CreateScan(ctx, scan))

	progress := &ScanProgress{
		ItemsScanned: 100,
		ItemsAdded:   80,
		ItemsUpdated: 10,
		ItemsRemoved: 5,
		ErrorsCount:  3,
	}

	updated, err := repo.UpdateScanProgress(ctx, scan.ID, progress)
	require.NoError(t, err)
	assert.Equal(t, int32(100), updated.ItemsScanned)
	assert.Equal(t, int32(80), updated.ItemsAdded)
	assert.Equal(t, int32(5), updated.ItemsRemoved)
}

// ============================================================================
// Permission Tests
// ============================================================================

func TestRepositoryPg_GrantPermission(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	lib := createTestLibrary("perm_lib", "movies")
	require.NoError(t, repo.Create(ctx, lib))

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "permuser",
		Email:    "perm@example.com",
	})

	perm := &Permission{
		LibraryID:  lib.ID,
		UserID:     user.ID,
		Permission: "read",
	}

	err := repo.GrantPermission(ctx, perm)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, perm.ID)
}

func TestRepositoryPg_CheckPermission(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	lib := createTestLibrary("check_perm_lib", "movies")
	require.NoError(t, repo.Create(ctx, lib))

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "checkuser",
		Email:    "check@example.com",
	})

	// Grant permission
	perm := &Permission{
		LibraryID:  lib.ID,
		UserID:     user.ID,
		Permission: "read",
	}
	require.NoError(t, repo.GrantPermission(ctx, perm))

	// Check permission exists
	hasPermission, err := repo.CheckPermission(ctx, lib.ID, user.ID, "read")
	require.NoError(t, err)
	assert.True(t, hasPermission)

	// Check permission doesn't exist
	hasWrite, err := repo.CheckPermission(ctx, lib.ID, user.ID, "write")
	require.NoError(t, err)
	assert.False(t, hasWrite)
}

func TestRepositoryPg_ListPermissions(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	lib := createTestLibrary("list_perm_lib", "movies")
	require.NoError(t, repo.Create(ctx, lib))

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "listpermuser",
		Email:    "listperm@example.com",
	})

	// Grant multiple permissions
	permissions := []string{"read", "write"}
	for _, p := range permissions {
		perm := &Permission{
			LibraryID:  lib.ID,
			UserID:     user.ID,
			Permission: p,
		}
		require.NoError(t, repo.GrantPermission(ctx, perm))
	}

	perms, err := repo.ListPermissions(ctx, lib.ID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(perms), 2)
}

func TestRepositoryPg_GetUserAccessibleLibraries(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "accessuser",
		Email:    "access@example.com",
	})

	// Create libraries and grant permissions
	for i := range 2 {
		lib := createTestLibrary("access_lib_"+string(rune('a'+i)), "movies")
		require.NoError(t, repo.Create(ctx, lib))

		perm := &Permission{
			LibraryID:  lib.ID,
			UserID:     user.ID,
			Permission: "view",
		}
		require.NoError(t, repo.GrantPermission(ctx, perm))
	}

	libraries, err := repo.GetUserAccessibleLibraries(ctx, user.ID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(libraries), 2)
}

func TestRepositoryPg_RevokePermission(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	lib := createTestLibrary("revoke_lib", "movies")
	require.NoError(t, repo.Create(ctx, lib))

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "revokeuser",
		Email:    "revoke@example.com",
	})

	// Grant permission
	perm := &Permission{
		LibraryID:  lib.ID,
		UserID:     user.ID,
		Permission: "read",
	}
	require.NoError(t, repo.GrantPermission(ctx, perm))

	// Revoke permission
	err := repo.RevokePermission(ctx, lib.ID, user.ID, "read")
	require.NoError(t, err)

	// Check permission is gone
	hasPermission, err := repo.CheckPermission(ctx, lib.ID, user.ID, "read")
	require.NoError(t, err)
	assert.False(t, hasPermission)
}

// ============================================================================
// Additional Repository Tests for Coverage
// ============================================================================

func TestRepositoryPg_DeleteOldScans(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	lib := createTestLibrary("delete_old_scans", "movies")
	require.NoError(t, repo.Create(ctx, lib))

	// Create some scans
	for range 3 {
		scan := &LibraryScan{
			LibraryID: lib.ID,
			ScanType:  ScanTypeFull,
			Status:    "completed",
		}
		require.NoError(t, repo.CreateScan(ctx, scan))
	}

	// Delete scans older than tomorrow (should delete all)
	deleted, err := repo.DeleteOldScans(ctx, time.Now().Add(24*time.Hour))
	require.NoError(t, err)
	assert.GreaterOrEqual(t, deleted, int64(3))

	// Verify scans are gone
	count, err := repo.CountScans(ctx, lib.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func TestRepositoryPg_DeleteOldScans_NoMatch(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	lib := createTestLibrary("delete_old_no_match", "movies")
	require.NoError(t, repo.Create(ctx, lib))

	scan := &LibraryScan{
		LibraryID: lib.ID,
		ScanType:  ScanTypeFull,
		Status:    "completed",
	}
	require.NoError(t, repo.CreateScan(ctx, scan))

	// Delete scans older than yesterday (should delete nothing)
	deleted, err := repo.DeleteOldScans(ctx, time.Now().Add(-24*time.Hour))
	require.NoError(t, err)
	assert.Equal(t, int64(0), deleted)
}

func TestRepositoryPg_RevokeUserPermissions(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "revoke_all_user",
		Email:    "revokeall@example.com",
	})

	// Create libraries and grant permissions
	for i := range 2 {
		lib := createTestLibrary("revoke_all_lib_"+string(rune('a'+i)), "movies")
		require.NoError(t, repo.Create(ctx, lib))

		perm := &Permission{
			LibraryID:  lib.ID,
			UserID:     user.ID,
			Permission: "view",
		}
		require.NoError(t, repo.GrantPermission(ctx, perm))
	}

	// Verify permissions exist
	perms, err := repo.ListUserPermissions(ctx, user.ID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(perms), 2)

	// Revoke all user permissions
	err = repo.RevokeUserPermissions(ctx, user.ID)
	require.NoError(t, err)

	// Verify all gone
	perms, err = repo.ListUserPermissions(ctx, user.ID)
	require.NoError(t, err)
	assert.Empty(t, perms)
}

func TestRepositoryPg_CountPermissions(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	lib := createTestLibrary("count_perm_lib", "movies")
	require.NoError(t, repo.Create(ctx, lib))

	user1 := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "countperm1",
		Email:    "cp1@example.com",
	})
	user2 := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "countperm2",
		Email:    "cp2@example.com",
	})

	// Grant permissions
	require.NoError(t, repo.GrantPermission(ctx, &Permission{LibraryID: lib.ID, UserID: user1.ID, Permission: "view"}))
	require.NoError(t, repo.GrantPermission(ctx, &Permission{LibraryID: lib.ID, UserID: user2.ID, Permission: "view"}))
	require.NoError(t, repo.GrantPermission(ctx, &Permission{LibraryID: lib.ID, UserID: user1.ID, Permission: "download"}))

	count, err := repo.CountPermissions(ctx, lib.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)
}

func TestRepositoryPg_GetPermission(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	lib := createTestLibrary("get_perm_lib", "movies")
	require.NoError(t, repo.Create(ctx, lib))

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "getpermuser",
		Email:    "getperm@example.com",
	})

	require.NoError(t, repo.GrantPermission(ctx, &Permission{LibraryID: lib.ID, UserID: user.ID, Permission: "view"}))

	perm, err := repo.GetPermission(ctx, lib.ID, user.ID, "view")
	require.NoError(t, err)
	assert.Equal(t, lib.ID, perm.LibraryID)
	assert.Equal(t, user.ID, perm.UserID)
	assert.Equal(t, "view", perm.Permission)

	// Not found
	_, err = repo.GetPermission(ctx, lib.ID, user.ID, "manage")
	assert.ErrorIs(t, err, ErrPermissionNotFound)
}

func TestRepositoryPg_Update_WithScannerConfig(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	lib := createTestLibrary("update_scanner", "movies")
	require.NoError(t, repo.Create(ctx, lib))

	update := &LibraryUpdate{
		ScannerConfig: map[string]any{
			"skip_hidden": true,
			"depth":       float64(5),
		},
	}

	updated, err := repo.Update(ctx, lib.ID, update)
	require.NoError(t, err)
	assert.NotNil(t, updated.ScannerConfig)
	assert.Equal(t, true, updated.ScannerConfig["skip_hidden"])
}

func TestRepositoryPg_Update_NotFound(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	update := &LibraryUpdate{Name: new("ghost")}
	_, err := repo.Update(ctx, uuid.Must(uuid.NewV7()), update)
	assert.ErrorIs(t, err, ErrNotFound)
}

// ============================================================================
// Helper Functions
// ============================================================================
