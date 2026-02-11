package library_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lusoris/revenge/internal/infra/logging"
	"github.com/lusoris/revenge/internal/service/library"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	infrajobs "github.com/lusoris/revenge/internal/infra/jobs"
)

func TestLibraryScanCleanupArgs_Kind(t *testing.T) {
	args := library.LibraryScanCleanupArgs{}
	assert.Equal(t, "library_scan_cleanup", args.Kind())
}

func TestLibraryScanCleanupArgs_InsertOpts(t *testing.T) {
	args := library.LibraryScanCleanupArgs{}
	opts := args.InsertOpts()
	assert.Equal(t, infrajobs.QueueLow, opts.Queue)
}

func TestNewLibraryScanCleanupWorker(t *testing.T) {
	mockRepo := NewMockLibraryRepository(t)
	logger := logging.NewTestLogger()

	worker := library.NewLibraryScanCleanupWorker(mockRepo, logger)
	assert.NotNil(t, worker)
}

func TestLibraryScanCleanupWorker_Timeout(t *testing.T) {
	mockRepo := NewMockLibraryRepository(t)
	logger := logging.NewTestLogger()
	worker := library.NewLibraryScanCleanupWorker(mockRepo, logger)

	job := &river.Job[library.LibraryScanCleanupArgs]{
		JobRow: &rivertype.JobRow{ID: 1},
		Args:   library.LibraryScanCleanupArgs{RetentionDays: 30},
	}

	timeout := worker.Timeout(job)
	assert.Equal(t, 2*time.Minute, timeout)
}

func TestLibraryScanCleanupWorker_Work(t *testing.T) {
	t.Run("success with cleanup", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		logger := logging.NewTestLogger()
		worker := library.NewLibraryScanCleanupWorker(mockRepo, logger)

		mockRepo.On("DeleteOldScans", mock.Anything, mock.AnythingOfType("time.Time")).
			Return(int64(5), nil)

		job := &river.Job[library.LibraryScanCleanupArgs]{
			JobRow: &rivertype.JobRow{ID: 1},
			Args:   library.LibraryScanCleanupArgs{RetentionDays: 30},
		}

		err := worker.Work(context.Background(), job)
		require.NoError(t, err)
		mockRepo.AssertCalled(t, "DeleteOldScans", mock.Anything, mock.AnythingOfType("time.Time"))
	})

	t.Run("success with dry run", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		logger := logging.NewTestLogger()
		worker := library.NewLibraryScanCleanupWorker(mockRepo, logger)

		job := &river.Job[library.LibraryScanCleanupArgs]{
			JobRow: &rivertype.JobRow{ID: 2},
			Args: library.LibraryScanCleanupArgs{
				RetentionDays: 7,
				DryRun:        true,
			},
		}

		err := worker.Work(context.Background(), job)
		require.NoError(t, err)
		// DeleteOldScans should NOT be called in dry run
		mockRepo.AssertNotCalled(t, "DeleteOldScans")
	})

	t.Run("default retention days when zero", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		logger := logging.NewTestLogger()
		worker := library.NewLibraryScanCleanupWorker(mockRepo, logger)

		mockRepo.On("DeleteOldScans", mock.Anything, mock.AnythingOfType("time.Time")).
			Return(int64(0), nil)

		job := &river.Job[library.LibraryScanCleanupArgs]{
			JobRow: &rivertype.JobRow{ID: 3},
			Args:   library.LibraryScanCleanupArgs{RetentionDays: 0},
		}

		err := worker.Work(context.Background(), job)
		require.NoError(t, err)
	})

	t.Run("default retention days when negative", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		logger := logging.NewTestLogger()
		worker := library.NewLibraryScanCleanupWorker(mockRepo, logger)

		mockRepo.On("DeleteOldScans", mock.Anything, mock.AnythingOfType("time.Time")).
			Return(int64(0), nil)

		job := &river.Job[library.LibraryScanCleanupArgs]{
			JobRow: &rivertype.JobRow{ID: 4},
			Args:   library.LibraryScanCleanupArgs{RetentionDays: -5},
		}

		err := worker.Work(context.Background(), job)
		require.NoError(t, err)
	})

	t.Run("error from repo", func(t *testing.T) {
		mockRepo := NewMockLibraryRepository(t)
		logger := logging.NewTestLogger()
		worker := library.NewLibraryScanCleanupWorker(mockRepo, logger)

		mockRepo.On("DeleteOldScans", mock.Anything, mock.AnythingOfType("time.Time")).
			Return(int64(0), errors.New("db error"))

		job := &river.Job[library.LibraryScanCleanupArgs]{
			JobRow: &rivertype.JobRow{ID: 5},
			Args:   library.LibraryScanCleanupArgs{RetentionDays: 30},
		}

		err := worker.Work(context.Background(), job)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to cleanup scans")
	})
}
