package library_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/service/activity"
	"github.com/lusoris/revenge/internal/service/library"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func makeTestLibrary(id uuid.UUID, name, libType string) *library.Library {
	now := time.Now()
	return &library.Library{
		ID:                id,
		Name:              name,
		Type:              libType,
		Paths:             []string{"/media/" + name},
		Enabled:           true,
		ScanOnStartup:     true,
		PreferredLanguage: "en",
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

func makeTestScan(id, libraryID uuid.UUID, scanType, status string) *library.LibraryScan {
	now := time.Now()
	return &library.LibraryScan{
		ID:        id,
		LibraryID: libraryID,
		ScanType:  scanType,
		Status:    status,
		CreatedAt: now,
	}
}

func setupLibraryService(repo library.Repository) *library.Service {
	logger := zap.NewNop()
	activityLogger := activity.NewNoopLogger()
	return library.NewService(repo, logger, activityLogger)
}

func TestLibraryService_Create_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		mockRepo.On("GetByName", mock.Anything, "Movies").
			Return(nil, library.ErrNotFound)
		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*library.Library")).
			Return(nil)

		req := library.CreateLibraryRequest{
			Name:    "Movies",
			Type:    library.LibraryTypeMovie,
			Paths:   []string{"/media/movies"},
			Enabled: true,
		}

		lib, err := svc.Create(context.Background(), req)

		require.NoError(t, err)
		assert.Equal(t, "Movies", lib.Name)
		assert.Equal(t, library.LibraryTypeMovie, lib.Type)
	})

	t.Run("invalid library type", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		req := library.CreateLibraryRequest{
			Name:    "Invalid",
			Type:    "invalid_type",
			Paths:   []string{"/media/invalid"},
			Enabled: true,
		}

		lib, err := svc.Create(context.Background(), req)

		assert.Nil(t, lib)
		assert.ErrorIs(t, err, library.ErrInvalidLibraryType)
	})

	t.Run("library exists", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		existingLib := makeTestLibrary(uuid.Must(uuid.NewV7()), "Movies", library.LibraryTypeMovie)
		mockRepo.On("GetByName", mock.Anything, "Movies").Return(existingLib, nil)

		req := library.CreateLibraryRequest{
			Name:    "Movies",
			Type:    library.LibraryTypeMovie,
			Paths:   []string{"/media/movies"},
			Enabled: true,
		}

		lib, err := svc.Create(context.Background(), req)

		assert.Nil(t, lib)
		assert.ErrorIs(t, err, library.ErrLibraryExists)
	})
}

func TestLibraryService_Get_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		expected := makeTestLibrary(libID, "Movies", library.LibraryTypeMovie)
		mockRepo.On("Get", mock.Anything, libID).Return(expected, nil)

		lib, err := svc.Get(context.Background(), libID)

		require.NoError(t, err)
		assert.Equal(t, expected.ID, lib.ID)
		assert.Equal(t, "Movies", lib.Name)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		mockRepo.On("Get", mock.Anything, libID).Return(nil, library.ErrNotFound)

		lib, err := svc.Get(context.Background(), libID)

		assert.Nil(t, lib)
		assert.ErrorIs(t, err, library.ErrNotFound)
	})
}

func TestLibraryService_GetByName_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		expected := makeTestLibrary(uuid.Must(uuid.NewV7()), "Movies", library.LibraryTypeMovie)
		mockRepo.On("GetByName", mock.Anything, "Movies").Return(expected, nil)

		lib, err := svc.GetByName(context.Background(), "Movies")

		require.NoError(t, err)
		assert.Equal(t, "Movies", lib.Name)
	})
}

func TestLibraryService_List_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libs := []library.Library{
			*makeTestLibrary(uuid.Must(uuid.NewV7()), "Movies", library.LibraryTypeMovie),
			*makeTestLibrary(uuid.Must(uuid.NewV7()), "TV Shows", library.LibraryTypeTVShow),
		}
		mockRepo.On("List", mock.Anything).Return(libs, nil)

		result, err := svc.List(context.Background())

		require.NoError(t, err)
		assert.Len(t, result, 2)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		mockRepo.On("List", mock.Anything).Return(nil, errors.New("db error"))

		result, err := svc.List(context.Background())

		assert.Nil(t, result)
		assert.Error(t, err)
	})
}

func TestLibraryService_ListEnabled_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libs := []library.Library{
			*makeTestLibrary(uuid.Must(uuid.NewV7()), "Movies", library.LibraryTypeMovie),
		}
		mockRepo.On("ListEnabled", mock.Anything).Return(libs, nil)

		result, err := svc.ListEnabled(context.Background())

		require.NoError(t, err)
		assert.Len(t, result, 1)
	})
}

func TestLibraryService_ListByType_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libs := []library.Library{
			*makeTestLibrary(uuid.Must(uuid.NewV7()), "Movies", library.LibraryTypeMovie),
		}
		mockRepo.On("ListByType", mock.Anything, library.LibraryTypeMovie).Return(libs, nil)

		result, err := svc.ListByType(context.Background(), library.LibraryTypeMovie)

		require.NoError(t, err)
		assert.Len(t, result, 1)
	})

	t.Run("invalid type", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		result, err := svc.ListByType(context.Background(), "invalid")

		assert.Nil(t, result)
		assert.ErrorIs(t, err, library.ErrInvalidLibraryType)
	})
}

func TestLibraryService_ListAccessible_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		userID := uuid.Must(uuid.NewV7())
		libs := []library.Library{
			*makeTestLibrary(uuid.Must(uuid.NewV7()), "Movies", library.LibraryTypeMovie),
		}
		mockRepo.On("GetUserAccessibleLibraries", mock.Anything, userID).Return(libs, nil)

		result, err := svc.ListAccessible(context.Background(), userID)

		require.NoError(t, err)
		assert.Len(t, result, 1)
	})
}

func TestLibraryService_Update_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		newName := "Updated Movies"
		updated := makeTestLibrary(libID, newName, library.LibraryTypeMovie)

		mockRepo.On("GetByName", mock.Anything, newName).Return(nil, library.ErrNotFound)
		mockRepo.On("Update", mock.Anything, libID, mock.AnythingOfType("*library.LibraryUpdate")).
			Return(updated, nil)

		update := &library.LibraryUpdate{Name: &newName}
		lib, err := svc.Update(context.Background(), libID, update)

		require.NoError(t, err)
		assert.Equal(t, newName, lib.Name)
	})

	t.Run("invalid type in update", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		invalidType := "invalid"
		update := &library.LibraryUpdate{Type: &invalidType}

		lib, err := svc.Update(context.Background(), libID, update)

		assert.Nil(t, lib)
		assert.ErrorIs(t, err, library.ErrInvalidLibraryType)
	})

	t.Run("name already exists", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		existingID := uuid.Must(uuid.NewV7())
		newName := "Existing Library"
		existing := makeTestLibrary(existingID, newName, library.LibraryTypeMovie)

		mockRepo.On("GetByName", mock.Anything, newName).Return(existing, nil)

		update := &library.LibraryUpdate{Name: &newName}
		lib, err := svc.Update(context.Background(), libID, update)

		assert.Nil(t, lib)
		assert.ErrorIs(t, err, library.ErrLibraryExists)
	})
}

func TestLibraryService_Delete_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		lib := makeTestLibrary(libID, "Movies", library.LibraryTypeMovie)

		mockRepo.On("Get", mock.Anything, libID).Return(lib, nil)
		mockRepo.On("RevokeAllPermissions", mock.Anything, libID).Return(nil)
		mockRepo.On("Delete", mock.Anything, libID).Return(nil)

		err := svc.Delete(context.Background(), libID)

		assert.NoError(t, err)
	})

	t.Run("revoke permissions error", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())

		mockRepo.On("Get", mock.Anything, libID).Return(nil, library.ErrNotFound)
		mockRepo.On("RevokeAllPermissions", mock.Anything, libID).Return(errors.New("db error"))

		err := svc.Delete(context.Background(), libID)

		assert.Error(t, err)
	})
}

func TestLibraryService_Count_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		mockRepo.On("Count", mock.Anything).Return(int64(5), nil)

		count, err := svc.Count(context.Background())

		require.NoError(t, err)
		assert.Equal(t, int64(5), count)
	})
}

func TestLibraryService_TriggerScan_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		lib := makeTestLibrary(libID, "Movies", library.LibraryTypeMovie)

		mockRepo.On("Get", mock.Anything, libID).Return(lib, nil)
		mockRepo.On("GetRunningScans", mock.Anything).Return([]library.LibraryScan{}, nil)
		mockRepo.On("CreateScan", mock.Anything, mock.AnythingOfType("*library.LibraryScan")).Return(nil)

		scan, err := svc.TriggerScan(context.Background(), libID, library.ScanTypeFull)

		require.NoError(t, err)
		assert.Equal(t, libID, scan.LibraryID)
		assert.Equal(t, library.ScanTypeFull, scan.ScanType)
		assert.Equal(t, library.ScanStatusPending, scan.Status)
	})

	t.Run("invalid scan type", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())

		scan, err := svc.TriggerScan(context.Background(), libID, "invalid")

		assert.Nil(t, scan)
		assert.ErrorIs(t, err, library.ErrInvalidScanType)
	})

	t.Run("scan already in progress", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		lib := makeTestLibrary(libID, "Movies", library.LibraryTypeMovie)

		runningScans := []library.LibraryScan{
			*makeTestScan(uuid.Must(uuid.NewV7()), libID, library.ScanTypeFull, library.ScanStatusRunning),
		}

		mockRepo.On("Get", mock.Anything, libID).Return(lib, nil)
		mockRepo.On("GetRunningScans", mock.Anything).Return(runningScans, nil)

		scan, err := svc.TriggerScan(context.Background(), libID, library.ScanTypeFull)

		assert.Nil(t, scan)
		assert.ErrorIs(t, err, library.ErrScanInProgress)
	})
}

func TestLibraryService_GetScan_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		scanID := uuid.Must(uuid.NewV7())
		libID := uuid.Must(uuid.NewV7())
		expected := makeTestScan(scanID, libID, library.ScanTypeFull, library.ScanStatusCompleted)

		mockRepo.On("GetScan", mock.Anything, scanID).Return(expected, nil)

		scan, err := svc.GetScan(context.Background(), scanID)

		require.NoError(t, err)
		assert.Equal(t, scanID, scan.ID)
	})
}

func TestLibraryService_ListScans_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success with default limit", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		scans := []library.LibraryScan{
			*makeTestScan(uuid.Must(uuid.NewV7()), libID, library.ScanTypeFull, library.ScanStatusCompleted),
		}

		mockRepo.On("ListScans", mock.Anything, libID, int32(20), int32(0)).Return(scans, nil)
		mockRepo.On("CountScans", mock.Anything, libID).Return(int64(1), nil)

		result, count, err := svc.ListScans(context.Background(), libID, 0, 0)

		require.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, int64(1), count)
	})

	t.Run("limit capped at 100", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		scans := []library.LibraryScan{}

		mockRepo.On("ListScans", mock.Anything, libID, int32(100), int32(0)).Return(scans, nil)
		mockRepo.On("CountScans", mock.Anything, libID).Return(int64(0), nil)

		result, count, err := svc.ListScans(context.Background(), libID, 200, 0)

		require.NoError(t, err)
		assert.Len(t, result, 0)
		assert.Equal(t, int64(0), count)
	})
}

func TestLibraryService_StartScan_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		scanID := uuid.Must(uuid.NewV7())
		libID := uuid.Must(uuid.NewV7())
		expected := makeTestScan(scanID, libID, library.ScanTypeFull, library.ScanStatusRunning)

		mockRepo.On("UpdateScanStatus", mock.Anything, scanID, mock.AnythingOfType("*library.ScanStatusUpdate")).
			Return(expected, nil)

		scan, err := svc.StartScan(context.Background(), scanID)

		require.NoError(t, err)
		assert.Equal(t, scanID, scan.ID)
	})
}

func TestLibraryService_CompleteScan_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success with progress", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		scanID := uuid.Must(uuid.NewV7())
		libID := uuid.Must(uuid.NewV7())
		startedAt := time.Now().Add(-10 * time.Minute)

		runningScan := makeTestScan(scanID, libID, library.ScanTypeFull, library.ScanStatusRunning)
		runningScan.StartedAt = &startedAt

		completed := makeTestScan(scanID, libID, library.ScanTypeFull, library.ScanStatusCompleted)

		progress := &library.ScanProgress{
			ItemsScanned: 100,
			ItemsAdded:   50,
		}

		mockRepo.On("UpdateScanProgress", mock.Anything, scanID, progress).Return(runningScan, nil)
		mockRepo.On("GetScan", mock.Anything, scanID).Return(runningScan, nil)
		mockRepo.On("UpdateScanStatus", mock.Anything, scanID, mock.AnythingOfType("*library.ScanStatusUpdate")).
			Return(completed, nil)

		scan, err := svc.CompleteScan(context.Background(), scanID, progress)

		require.NoError(t, err)
		assert.Equal(t, library.ScanStatusCompleted, scan.Status)
	})

	t.Run("success without progress", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		scanID := uuid.Must(uuid.NewV7())
		libID := uuid.Must(uuid.NewV7())

		runningScan := makeTestScan(scanID, libID, library.ScanTypeFull, library.ScanStatusRunning)
		completed := makeTestScan(scanID, libID, library.ScanTypeFull, library.ScanStatusCompleted)

		mockRepo.On("GetScan", mock.Anything, scanID).Return(runningScan, nil)
		mockRepo.On("UpdateScanStatus", mock.Anything, scanID, mock.AnythingOfType("*library.ScanStatusUpdate")).
			Return(completed, nil)

		scan, err := svc.CompleteScan(context.Background(), scanID, nil)

		require.NoError(t, err)
		assert.Equal(t, library.ScanStatusCompleted, scan.Status)
	})
}

func TestLibraryService_FailScan_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		scanID := uuid.Must(uuid.NewV7())
		libID := uuid.Must(uuid.NewV7())

		runningScan := makeTestScan(scanID, libID, library.ScanTypeFull, library.ScanStatusRunning)
		failed := makeTestScan(scanID, libID, library.ScanTypeFull, library.ScanStatusFailed)

		mockRepo.On("GetScan", mock.Anything, scanID).Return(runningScan, nil)
		mockRepo.On("UpdateScanStatus", mock.Anything, scanID, mock.AnythingOfType("*library.ScanStatusUpdate")).
			Return(failed, nil)

		scan, err := svc.FailScan(context.Background(), scanID, "test error")

		require.NoError(t, err)
		assert.Equal(t, library.ScanStatusFailed, scan.Status)
	})
}

func TestLibraryService_CancelScan_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		scanID := uuid.Must(uuid.NewV7())
		libID := uuid.Must(uuid.NewV7())

		cancelled := makeTestScan(scanID, libID, library.ScanTypeFull, library.ScanStatusCancelled)

		mockRepo.On("UpdateScanStatus", mock.Anything, scanID, mock.AnythingOfType("*library.ScanStatusUpdate")).
			Return(cancelled, nil)

		scan, err := svc.CancelScan(context.Background(), scanID)

		require.NoError(t, err)
		assert.Equal(t, library.ScanStatusCancelled, scan.Status)
	})
}

func TestLibraryService_GrantPermission_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		userID := uuid.Must(uuid.NewV7())
		lib := makeTestLibrary(libID, "Movies", library.LibraryTypeMovie)

		mockRepo.On("Get", mock.Anything, libID).Return(lib, nil)
		mockRepo.On("GrantPermission", mock.Anything, mock.AnythingOfType("*library.Permission")).Return(nil)

		err := svc.GrantPermission(context.Background(), libID, userID, library.PermissionView)

		assert.NoError(t, err)
	})

	t.Run("invalid permission", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		userID := uuid.Must(uuid.NewV7())

		err := svc.GrantPermission(context.Background(), libID, userID, "invalid")

		assert.ErrorIs(t, err, library.ErrInvalidPermission)
	})

	t.Run("library not found", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		userID := uuid.Must(uuid.NewV7())

		mockRepo.On("Get", mock.Anything, libID).Return(nil, library.ErrNotFound)

		err := svc.GrantPermission(context.Background(), libID, userID, library.PermissionView)

		assert.ErrorIs(t, err, library.ErrNotFound)
	})
}

func TestLibraryService_RevokePermission_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		userID := uuid.Must(uuid.NewV7())

		mockRepo.On("RevokePermission", mock.Anything, libID, userID, library.PermissionView).Return(nil)

		err := svc.RevokePermission(context.Background(), libID, userID, library.PermissionView)

		assert.NoError(t, err)
	})

	t.Run("invalid permission", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		userID := uuid.Must(uuid.NewV7())

		err := svc.RevokePermission(context.Background(), libID, userID, "invalid")

		assert.ErrorIs(t, err, library.ErrInvalidPermission)
	})
}

func TestLibraryService_CheckPermission_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("has permission", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		userID := uuid.Must(uuid.NewV7())

		mockRepo.On("CheckPermission", mock.Anything, libID, userID, library.PermissionView).Return(true, nil)

		has, err := svc.CheckPermission(context.Background(), libID, userID, library.PermissionView)

		require.NoError(t, err)
		assert.True(t, has)
	})

	t.Run("invalid permission", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		userID := uuid.Must(uuid.NewV7())

		has, err := svc.CheckPermission(context.Background(), libID, userID, "invalid")

		assert.False(t, has)
		assert.ErrorIs(t, err, library.ErrInvalidPermission)
	})
}

func TestLibraryService_CanAccess_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("admin always has access", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		userID := uuid.Must(uuid.NewV7())

		can, err := svc.CanAccess(context.Background(), libID, userID, true)

		require.NoError(t, err)
		assert.True(t, can)
	})

	t.Run("non-admin checks permission", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		userID := uuid.Must(uuid.NewV7())

		mockRepo.On("CheckPermission", mock.Anything, libID, userID, library.PermissionView).Return(true, nil)

		can, err := svc.CanAccess(context.Background(), libID, userID, false)

		require.NoError(t, err)
		assert.True(t, can)
	})
}

func TestLibraryService_CanDownload_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("admin always can", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		userID := uuid.Must(uuid.NewV7())

		can, err := svc.CanDownload(context.Background(), libID, userID, true)

		require.NoError(t, err)
		assert.True(t, can)
	})

	t.Run("non-admin checks permission", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		userID := uuid.Must(uuid.NewV7())

		mockRepo.On("CheckPermission", mock.Anything, libID, userID, library.PermissionDownload).Return(false, nil)

		can, err := svc.CanDownload(context.Background(), libID, userID, false)

		require.NoError(t, err)
		assert.False(t, can)
	})
}

func TestLibraryService_CanManage_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("admin always can", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		userID := uuid.Must(uuid.NewV7())

		can, err := svc.CanManage(context.Background(), libID, userID, true)

		require.NoError(t, err)
		assert.True(t, can)
	})
}

func TestLibraryService_ListPermissions_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		perms := []library.Permission{
			{ID: uuid.Must(uuid.NewV7()), LibraryID: libID, UserID: uuid.Must(uuid.NewV7()), Permission: library.PermissionView},
		}

		mockRepo.On("ListPermissions", mock.Anything, libID).Return(perms, nil)

		result, err := svc.ListPermissions(context.Background(), libID)

		require.NoError(t, err)
		assert.Len(t, result, 1)
	})
}

func TestLibraryService_ListUserPermissions_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		userID := uuid.Must(uuid.NewV7())
		perms := []library.Permission{
			{ID: uuid.Must(uuid.NewV7()), LibraryID: uuid.Must(uuid.NewV7()), UserID: userID, Permission: library.PermissionView},
		}

		mockRepo.On("ListUserPermissions", mock.Anything, userID).Return(perms, nil)

		result, err := svc.ListUserPermissions(context.Background(), userID)

		require.NoError(t, err)
		assert.Len(t, result, 1)
	})
}

func TestLibraryService_GetPermission_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		userID := uuid.Must(uuid.NewV7())
		perm := &library.Permission{
			ID:         uuid.Must(uuid.NewV7()),
			LibraryID:  libID,
			UserID:     userID,
			Permission: library.PermissionView,
		}

		mockRepo.On("GetPermission", mock.Anything, libID, userID, library.PermissionView).Return(perm, nil)

		result, err := svc.GetPermission(context.Background(), libID, userID, library.PermissionView)

		require.NoError(t, err)
		assert.Equal(t, library.PermissionView, result.Permission)
	})
}

// ============================================================================
// Additional Tests for Uncovered Code Paths
// ============================================================================

func TestLibraryService_GetLatestScan_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		scanID := uuid.Must(uuid.NewV7())
		expected := makeTestScan(scanID, libID, library.ScanTypeFull, library.ScanStatusCompleted)

		mockRepo.On("GetLatestScan", mock.Anything, libID).Return(expected, nil)

		scan, err := svc.GetLatestScan(context.Background(), libID)

		require.NoError(t, err)
		assert.Equal(t, scanID, scan.ID)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())

		mockRepo.On("GetLatestScan", mock.Anything, libID).Return(nil, library.ErrScanNotFound)

		scan, err := svc.GetLatestScan(context.Background(), libID)

		assert.Nil(t, scan)
		assert.ErrorIs(t, err, library.ErrScanNotFound)
	})
}

func TestLibraryService_GetRunningScans_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success with running scans", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		scans := []library.LibraryScan{
			*makeTestScan(uuid.Must(uuid.NewV7()), libID, library.ScanTypeFull, library.ScanStatusRunning),
			*makeTestScan(uuid.Must(uuid.NewV7()), uuid.Must(uuid.NewV7()), library.ScanTypeIncremental, library.ScanStatusRunning),
		}

		mockRepo.On("GetRunningScans", mock.Anything).Return(scans, nil)

		result, err := svc.GetRunningScans(context.Background())

		require.NoError(t, err)
		assert.Len(t, result, 2)
	})

	t.Run("success empty", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		mockRepo.On("GetRunningScans", mock.Anything).Return([]library.LibraryScan{}, nil)

		result, err := svc.GetRunningScans(context.Background())

		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		mockRepo.On("GetRunningScans", mock.Anything).Return(nil, errors.New("db error"))

		result, err := svc.GetRunningScans(context.Background())

		assert.Nil(t, result)
		assert.Error(t, err)
	})
}

func TestLibraryService_UpdateScanProgress_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		scanID := uuid.Must(uuid.NewV7())
		libID := uuid.Must(uuid.NewV7())
		expected := makeTestScan(scanID, libID, library.ScanTypeFull, library.ScanStatusRunning)
		expected.ItemsScanned = 50
		expected.ItemsAdded = 25

		progress := &library.ScanProgress{
			ItemsScanned: 50,
			ItemsAdded:   25,
		}

		mockRepo.On("UpdateScanProgress", mock.Anything, scanID, progress).Return(expected, nil)

		scan, err := svc.UpdateScanProgress(context.Background(), scanID, progress)

		require.NoError(t, err)
		assert.Equal(t, int32(50), scan.ItemsScanned)
		assert.Equal(t, int32(25), scan.ItemsAdded)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		scanID := uuid.Must(uuid.NewV7())
		progress := &library.ScanProgress{
			ItemsScanned: 50,
		}

		mockRepo.On("UpdateScanProgress", mock.Anything, scanID, progress).Return(nil, errors.New("db error"))

		scan, err := svc.UpdateScanProgress(context.Background(), scanID, progress)

		assert.Nil(t, scan)
		assert.Error(t, err)
	})
}

func TestLibraryService_CanManage_NonAdmin_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("non-admin with permission", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		userID := uuid.Must(uuid.NewV7())

		mockRepo.On("CheckPermission", mock.Anything, libID, userID, library.PermissionManage).Return(true, nil)

		can, err := svc.CanManage(context.Background(), libID, userID, false)

		require.NoError(t, err)
		assert.True(t, can)
	})

	t.Run("non-admin without permission", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		userID := uuid.Must(uuid.NewV7())

		mockRepo.On("CheckPermission", mock.Anything, libID, userID, library.PermissionManage).Return(false, nil)

		can, err := svc.CanManage(context.Background(), libID, userID, false)

		require.NoError(t, err)
		assert.False(t, can)
	})
}

func TestLibraryService_Create_ErrorPaths_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("GetByName returns unexpected error", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		mockRepo.On("GetByName", mock.Anything, "Movies").Return(nil, errors.New("db error"))

		req := library.CreateLibraryRequest{
			Name:    "Movies",
			Type:    library.LibraryTypeMovie,
			Paths:   []string{"/media/movies"},
			Enabled: true,
		}

		lib, err := svc.Create(context.Background(), req)

		assert.Nil(t, lib)
		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
	})

	t.Run("repo Create returns error", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		mockRepo.On("GetByName", mock.Anything, "Movies").Return(nil, library.ErrNotFound)
		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*library.Library")).Return(errors.New("create error"))

		req := library.CreateLibraryRequest{
			Name:    "Movies",
			Type:    library.LibraryTypeMovie,
			Paths:   []string{"/media/movies"},
			Enabled: true,
		}

		lib, err := svc.Create(context.Background(), req)

		assert.Nil(t, lib)
		assert.Error(t, err)
		assert.Equal(t, "create error", err.Error())
	})
}

func TestLibraryService_Update_ErrorPaths_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("GetByName returns unexpected error", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		newName := "Updated"

		mockRepo.On("GetByName", mock.Anything, newName).Return(nil, errors.New("db error"))

		update := &library.LibraryUpdate{Name: &newName}
		lib, err := svc.Update(context.Background(), libID, update)

		assert.Nil(t, lib)
		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
	})

	t.Run("repo Update returns error", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		newName := "Updated"

		mockRepo.On("GetByName", mock.Anything, newName).Return(nil, library.ErrNotFound)
		mockRepo.On("Update", mock.Anything, libID, mock.AnythingOfType("*library.LibraryUpdate")).Return(nil, errors.New("update error"))

		update := &library.LibraryUpdate{Name: &newName}
		lib, err := svc.Update(context.Background(), libID, update)

		assert.Nil(t, lib)
		assert.Error(t, err)
		assert.Equal(t, "update error", err.Error())
	})

	t.Run("update without name change", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		enabled := true
		updated := makeTestLibrary(libID, "Movies", library.LibraryTypeMovie)

		mockRepo.On("Update", mock.Anything, libID, mock.AnythingOfType("*library.LibraryUpdate")).Return(updated, nil)

		update := &library.LibraryUpdate{Enabled: &enabled}
		lib, err := svc.Update(context.Background(), libID, update)

		require.NoError(t, err)
		assert.Equal(t, libID, lib.ID)
	})

	t.Run("update same library name allowed", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		newName := "Movies"
		existing := makeTestLibrary(libID, newName, library.LibraryTypeMovie)
		updated := makeTestLibrary(libID, newName, library.LibraryTypeMovie)

		mockRepo.On("GetByName", mock.Anything, newName).Return(existing, nil)
		mockRepo.On("Update", mock.Anything, libID, mock.AnythingOfType("*library.LibraryUpdate")).Return(updated, nil)

		update := &library.LibraryUpdate{Name: &newName}
		lib, err := svc.Update(context.Background(), libID, update)

		require.NoError(t, err)
		assert.Equal(t, newName, lib.Name)
	})
}

func TestLibraryService_Delete_ErrorPaths_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("delete error", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		lib := makeTestLibrary(libID, "Movies", library.LibraryTypeMovie)

		mockRepo.On("Get", mock.Anything, libID).Return(lib, nil)
		mockRepo.On("RevokeAllPermissions", mock.Anything, libID).Return(nil)
		mockRepo.On("Delete", mock.Anything, libID).Return(errors.New("delete error"))

		err := svc.Delete(context.Background(), libID)

		assert.Error(t, err)
		assert.Equal(t, "delete error", err.Error())
	})

	t.Run("delete without library found for logging", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())

		mockRepo.On("Get", mock.Anything, libID).Return(nil, library.ErrNotFound)
		mockRepo.On("RevokeAllPermissions", mock.Anything, libID).Return(nil)
		mockRepo.On("Delete", mock.Anything, libID).Return(nil)

		err := svc.Delete(context.Background(), libID)

		assert.NoError(t, err)
	})
}

func TestLibraryService_TriggerScan_ErrorPaths_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("library not found", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())

		mockRepo.On("Get", mock.Anything, libID).Return(nil, library.ErrNotFound)

		scan, err := svc.TriggerScan(context.Background(), libID, library.ScanTypeFull)

		assert.Nil(t, scan)
		assert.ErrorIs(t, err, library.ErrNotFound)
	})

	t.Run("GetRunningScans error", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		lib := makeTestLibrary(libID, "Movies", library.LibraryTypeMovie)

		mockRepo.On("Get", mock.Anything, libID).Return(lib, nil)
		mockRepo.On("GetRunningScans", mock.Anything).Return(nil, errors.New("db error"))

		scan, err := svc.TriggerScan(context.Background(), libID, library.ScanTypeFull)

		assert.Nil(t, scan)
		assert.Error(t, err)
	})

	t.Run("CreateScan error", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		lib := makeTestLibrary(libID, "Movies", library.LibraryTypeMovie)

		mockRepo.On("Get", mock.Anything, libID).Return(lib, nil)
		mockRepo.On("GetRunningScans", mock.Anything).Return([]library.LibraryScan{}, nil)
		mockRepo.On("CreateScan", mock.Anything, mock.AnythingOfType("*library.LibraryScan")).Return(errors.New("create error"))

		scan, err := svc.TriggerScan(context.Background(), libID, library.ScanTypeFull)

		assert.Nil(t, scan)
		assert.Error(t, err)
	})

	t.Run("scan in progress for different library allows trigger", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		otherLibID := uuid.Must(uuid.NewV7())
		lib := makeTestLibrary(libID, "Movies", library.LibraryTypeMovie)

		runningScans := []library.LibraryScan{
			*makeTestScan(uuid.Must(uuid.NewV7()), otherLibID, library.ScanTypeFull, library.ScanStatusRunning),
		}

		mockRepo.On("Get", mock.Anything, libID).Return(lib, nil)
		mockRepo.On("GetRunningScans", mock.Anything).Return(runningScans, nil)
		mockRepo.On("CreateScan", mock.Anything, mock.AnythingOfType("*library.LibraryScan")).Return(nil)

		scan, err := svc.TriggerScan(context.Background(), libID, library.ScanTypeFull)

		require.NoError(t, err)
		assert.Equal(t, libID, scan.LibraryID)
	})
}

func TestLibraryService_ListScans_ErrorPaths_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("ListScans error", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())

		mockRepo.On("ListScans", mock.Anything, libID, int32(20), int32(0)).Return(nil, errors.New("db error"))

		result, count, err := svc.ListScans(context.Background(), libID, 0, 0)

		assert.Nil(t, result)
		assert.Equal(t, int64(0), count)
		assert.Error(t, err)
	})

	t.Run("CountScans error", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		scans := []library.LibraryScan{
			*makeTestScan(uuid.Must(uuid.NewV7()), libID, library.ScanTypeFull, library.ScanStatusCompleted),
		}

		mockRepo.On("ListScans", mock.Anything, libID, int32(20), int32(0)).Return(scans, nil)
		mockRepo.On("CountScans", mock.Anything, libID).Return(int64(0), errors.New("count error"))

		result, count, err := svc.ListScans(context.Background(), libID, 0, 0)

		assert.Nil(t, result)
		assert.Equal(t, int64(0), count)
		assert.Error(t, err)
	})
}

func TestLibraryService_CompleteScan_ErrorPaths_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("UpdateScanProgress error", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		scanID := uuid.Must(uuid.NewV7())
		progress := &library.ScanProgress{ItemsScanned: 100}

		mockRepo.On("UpdateScanProgress", mock.Anything, scanID, progress).Return(nil, errors.New("progress error"))

		scan, err := svc.CompleteScan(context.Background(), scanID, progress)

		assert.Nil(t, scan)
		assert.Error(t, err)
	})

	t.Run("GetScan error", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		scanID := uuid.Must(uuid.NewV7())

		mockRepo.On("GetScan", mock.Anything, scanID).Return(nil, errors.New("get error"))

		scan, err := svc.CompleteScan(context.Background(), scanID, nil)

		assert.Nil(t, scan)
		assert.Error(t, err)
	})

	t.Run("scan with nil StartedAt", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		scanID := uuid.Must(uuid.NewV7())
		libID := uuid.Must(uuid.NewV7())

		runningScan := makeTestScan(scanID, libID, library.ScanTypeFull, library.ScanStatusRunning)
		runningScan.StartedAt = nil // No start time

		completed := makeTestScan(scanID, libID, library.ScanTypeFull, library.ScanStatusCompleted)

		mockRepo.On("GetScan", mock.Anything, scanID).Return(runningScan, nil)
		mockRepo.On("UpdateScanStatus", mock.Anything, scanID, mock.AnythingOfType("*library.ScanStatusUpdate")).Return(completed, nil)

		scan, err := svc.CompleteScan(context.Background(), scanID, nil)

		require.NoError(t, err)
		assert.Equal(t, library.ScanStatusCompleted, scan.Status)
	})
}

func TestLibraryService_FailScan_ErrorPaths_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("GetScan error", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		scanID := uuid.Must(uuid.NewV7())

		mockRepo.On("GetScan", mock.Anything, scanID).Return(nil, errors.New("get error"))

		scan, err := svc.FailScan(context.Background(), scanID, "error message")

		assert.Nil(t, scan)
		assert.Error(t, err)
	})

	t.Run("scan with StartedAt calculates duration", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		scanID := uuid.Must(uuid.NewV7())
		libID := uuid.Must(uuid.NewV7())
		startedAt := time.Now().Add(-5 * time.Minute)

		runningScan := makeTestScan(scanID, libID, library.ScanTypeFull, library.ScanStatusRunning)
		runningScan.StartedAt = &startedAt

		failed := makeTestScan(scanID, libID, library.ScanTypeFull, library.ScanStatusFailed)

		mockRepo.On("GetScan", mock.Anything, scanID).Return(runningScan, nil)
		mockRepo.On("UpdateScanStatus", mock.Anything, scanID, mock.MatchedBy(func(update *library.ScanStatusUpdate) bool {
			return update.Status == library.ScanStatusFailed &&
				update.DurationSeconds != nil &&
				*update.DurationSeconds >= 300 // At least 5 minutes
		})).Return(failed, nil)

		scan, err := svc.FailScan(context.Background(), scanID, "error message")

		require.NoError(t, err)
		assert.Equal(t, library.ScanStatusFailed, scan.Status)
	})
}

func TestLibraryService_GrantPermission_ErrorPaths_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("GrantPermission repo error", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		userID := uuid.Must(uuid.NewV7())
		lib := makeTestLibrary(libID, "Movies", library.LibraryTypeMovie)

		mockRepo.On("Get", mock.Anything, libID).Return(lib, nil)
		mockRepo.On("GrantPermission", mock.Anything, mock.AnythingOfType("*library.Permission")).Return(errors.New("grant error"))

		err := svc.GrantPermission(context.Background(), libID, userID, library.PermissionView)

		assert.Error(t, err)
		assert.Equal(t, "grant error", err.Error())
	})
}

func TestLibraryService_RevokePermission_ErrorPaths_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("RevokePermission repo error", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		svc := setupLibraryService(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		userID := uuid.Must(uuid.NewV7())

		mockRepo.On("RevokePermission", mock.Anything, libID, userID, library.PermissionView).Return(errors.New("revoke error"))

		err := svc.RevokePermission(context.Background(), libID, userID, library.PermissionView)

		assert.Error(t, err)
		assert.Equal(t, "revoke error", err.Error())
	})
}

// ============================================================================
// CachedService Tests
// ============================================================================

func setupCachedServiceWithNilCache(repo library.Repository) *library.CachedService {
	logger := zap.NewNop()
	activityLogger := activity.NewNoopLogger()
	svc := library.NewService(repo, logger, activityLogger)
	return library.NewCachedService(svc, nil, logger)
}

func TestCachedService_NewCachedService_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	mockRepo := NewMockLibraryRepository(t)
	logger := zap.NewNop()
	activityLogger := activity.NewNoopLogger()
	svc := library.NewService(mockRepo, logger, activityLogger)

	cachedSvc := library.NewCachedService(svc, nil, logger)

	assert.NotNil(t, cachedSvc)
}

func TestCachedService_Get_NilCache_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success with nil cache falls back to service", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		cachedSvc := setupCachedServiceWithNilCache(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		expected := makeTestLibrary(libID, "Movies", library.LibraryTypeMovie)

		mockRepo.On("Get", mock.Anything, libID).Return(expected, nil)

		lib, err := cachedSvc.Get(context.Background(), libID)

		require.NoError(t, err)
		assert.Equal(t, expected.ID, lib.ID)
		assert.Equal(t, "Movies", lib.Name)
	})

	t.Run("not found with nil cache", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		cachedSvc := setupCachedServiceWithNilCache(mockRepo)

		libID := uuid.Must(uuid.NewV7())

		mockRepo.On("Get", mock.Anything, libID).Return(nil, library.ErrNotFound)

		lib, err := cachedSvc.Get(context.Background(), libID)

		assert.Nil(t, lib)
		assert.ErrorIs(t, err, library.ErrNotFound)
	})
}

func TestCachedService_List_NilCache_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success with nil cache falls back to service", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		cachedSvc := setupCachedServiceWithNilCache(mockRepo)

		libs := []library.Library{
			*makeTestLibrary(uuid.Must(uuid.NewV7()), "Movies", library.LibraryTypeMovie),
			*makeTestLibrary(uuid.Must(uuid.NewV7()), "TV Shows", library.LibraryTypeTVShow),
		}

		mockRepo.On("List", mock.Anything).Return(libs, nil)

		result, err := cachedSvc.List(context.Background())

		require.NoError(t, err)
		assert.Len(t, result, 2)
	})

	t.Run("error with nil cache", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		cachedSvc := setupCachedServiceWithNilCache(mockRepo)

		mockRepo.On("List", mock.Anything).Return(nil, errors.New("db error"))

		result, err := cachedSvc.List(context.Background())

		assert.Nil(t, result)
		assert.Error(t, err)
	})
}

func TestCachedService_Count_NilCache_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success with nil cache falls back to service", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		cachedSvc := setupCachedServiceWithNilCache(mockRepo)

		mockRepo.On("Count", mock.Anything).Return(int64(5), nil)

		count, err := cachedSvc.Count(context.Background())

		require.NoError(t, err)
		assert.Equal(t, int64(5), count)
	})

	t.Run("error with nil cache", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		cachedSvc := setupCachedServiceWithNilCache(mockRepo)

		mockRepo.On("Count", mock.Anything).Return(int64(0), errors.New("db error"))

		count, err := cachedSvc.Count(context.Background())

		assert.Equal(t, int64(0), count)
		assert.Error(t, err)
	})
}

func TestCachedService_Create_NilCache_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success with nil cache", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		cachedSvc := setupCachedServiceWithNilCache(mockRepo)

		mockRepo.On("GetByName", mock.Anything, "Movies").Return(nil, library.ErrNotFound)
		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*library.Library")).Return(nil)

		req := library.CreateLibraryRequest{
			Name:    "Movies",
			Type:    library.LibraryTypeMovie,
			Paths:   []string{"/media/movies"},
			Enabled: true,
		}

		lib, err := cachedSvc.Create(context.Background(), req)

		require.NoError(t, err)
		assert.Equal(t, "Movies", lib.Name)
	})

	t.Run("error propagates with nil cache", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		cachedSvc := setupCachedServiceWithNilCache(mockRepo)

		req := library.CreateLibraryRequest{
			Name:    "Invalid",
			Type:    "invalid_type",
			Paths:   []string{"/media"},
			Enabled: true,
		}

		lib, err := cachedSvc.Create(context.Background(), req)

		assert.Nil(t, lib)
		assert.ErrorIs(t, err, library.ErrInvalidLibraryType)
	})
}

func TestCachedService_Update_NilCache_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success with nil cache", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		cachedSvc := setupCachedServiceWithNilCache(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		newName := "Updated Movies"
		updated := makeTestLibrary(libID, newName, library.LibraryTypeMovie)

		mockRepo.On("GetByName", mock.Anything, newName).Return(nil, library.ErrNotFound)
		mockRepo.On("Update", mock.Anything, libID, mock.AnythingOfType("*library.LibraryUpdate")).Return(updated, nil)

		update := &library.LibraryUpdate{Name: &newName}
		lib, err := cachedSvc.Update(context.Background(), libID, update)

		require.NoError(t, err)
		assert.Equal(t, newName, lib.Name)
	})

	t.Run("error propagates with nil cache", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		cachedSvc := setupCachedServiceWithNilCache(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		invalidType := "invalid"

		update := &library.LibraryUpdate{Type: &invalidType}
		lib, err := cachedSvc.Update(context.Background(), libID, update)

		assert.Nil(t, lib)
		assert.ErrorIs(t, err, library.ErrInvalidLibraryType)
	})
}

func TestCachedService_Delete_NilCache_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success with nil cache", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		cachedSvc := setupCachedServiceWithNilCache(mockRepo)

		libID := uuid.Must(uuid.NewV7())
		lib := makeTestLibrary(libID, "Movies", library.LibraryTypeMovie)

		mockRepo.On("Get", mock.Anything, libID).Return(lib, nil)
		mockRepo.On("RevokeAllPermissions", mock.Anything, libID).Return(nil)
		mockRepo.On("Delete", mock.Anything, libID).Return(nil)

		err := cachedSvc.Delete(context.Background(), libID)

		assert.NoError(t, err)
	})

	t.Run("error propagates with nil cache", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		cachedSvc := setupCachedServiceWithNilCache(mockRepo)

		libID := uuid.Must(uuid.NewV7())

		mockRepo.On("Get", mock.Anything, libID).Return(nil, library.ErrNotFound)
		mockRepo.On("RevokeAllPermissions", mock.Anything, libID).Return(errors.New("db error"))

		err := cachedSvc.Delete(context.Background(), libID)

		assert.Error(t, err)
	})
}

func TestCachedService_CompleteScan_NilCache_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success with nil cache", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		cachedSvc := setupCachedServiceWithNilCache(mockRepo)

		scanID := uuid.Must(uuid.NewV7())
		libID := uuid.Must(uuid.NewV7())

		runningScan := makeTestScan(scanID, libID, library.ScanTypeFull, library.ScanStatusRunning)
		completed := makeTestScan(scanID, libID, library.ScanTypeFull, library.ScanStatusCompleted)

		mockRepo.On("GetScan", mock.Anything, scanID).Return(runningScan, nil)
		mockRepo.On("UpdateScanStatus", mock.Anything, scanID, mock.AnythingOfType("*library.ScanStatusUpdate")).Return(completed, nil)

		scan, err := cachedSvc.CompleteScan(context.Background(), scanID, nil)

		require.NoError(t, err)
		assert.Equal(t, library.ScanStatusCompleted, scan.Status)
	})

	t.Run("error propagates with nil cache", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		cachedSvc := setupCachedServiceWithNilCache(mockRepo)

		scanID := uuid.Must(uuid.NewV7())

		mockRepo.On("GetScan", mock.Anything, scanID).Return(nil, errors.New("get error"))

		scan, err := cachedSvc.CompleteScan(context.Background(), scanID, nil)

		assert.Nil(t, scan)
		assert.Error(t, err)
	})
}

func TestCachedService_InvalidateLibraryCache_NilCache_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("nil cache returns nil error", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		cachedSvc := setupCachedServiceWithNilCache(mockRepo)

		err := cachedSvc.InvalidateLibraryCache(context.Background(), uuid.Must(uuid.NewV7()))

		assert.NoError(t, err)
	})
}
