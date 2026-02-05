package jobs

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestJobResult(t *testing.T) {
	t.Run("AddError", func(t *testing.T) {
		result := &JobResult{
			Success:        true,
			ItemsProcessed: 10,
		}

		err1 := errors.New("error 1")
		err2 := errors.New("error 2")

		result.AddError(err1)
		result.AddError(err2)

		assert.Len(t, result.Errors, 2)
		assert.Equal(t, 2, result.ItemsFailed)
		assert.Equal(t, err1, result.Errors[0])
		assert.Equal(t, err2, result.Errors[1])
	})

	t.Run("HasErrors", func(t *testing.T) {
		result := &JobResult{}
		assert.False(t, result.HasErrors())

		result.AddError(errors.New("test"))
		assert.True(t, result.HasErrors())
	})

	t.Run("LogSummary success", func(t *testing.T) {
		logger := zap.NewNop()
		result := &JobResult{
			Success:        true,
			ItemsProcessed: 10,
			Duration:       time.Second,
		}

		// Should not panic
		result.LogSummary(logger, "test_job")
	})

	t.Run("LogSummary with errors", func(t *testing.T) {
		logger := zap.NewNop()
		result := &JobResult{
			Success:        true,
			ItemsProcessed: 10,
			ItemsFailed:    2,
			Duration:       time.Second,
			Errors:         []error{errors.New("err1"), errors.New("err2")},
		}

		// Should not panic
		result.LogSummary(logger, "test_job")
	})

	t.Run("LogSummary failure", func(t *testing.T) {
		logger := zap.NewNop()
		result := &JobResult{
			Success:     false,
			ItemsFailed: 10,
			Duration:    time.Second,
			Errors:      []error{errors.New("fatal")},
		}

		// Should not panic
		result.LogSummary(logger, "test_job")
	})

	t.Run("LogErrors", func(t *testing.T) {
		logger := zap.NewNop()
		result := &JobResult{
			Errors: make([]error, 20),
		}
		for i := range result.Errors {
			result.Errors[i] = errors.New("error")
		}

		// Should truncate at maxErrors
		result.LogErrors(logger, 10)
	})
}

func TestLibraryScanArgs(t *testing.T) {
	args := LibraryScanArgs{
		Paths: []string{"/media/movies", "/media/tv"},
		Force: true,
	}

	assert.Len(t, args.Paths, 2)
	assert.True(t, args.Force)
}

func TestFileMatchArgs(t *testing.T) {
	args := FileMatchArgs{
		FilePath:     "/media/movies/test.mkv",
		ForceRematch: true,
	}

	assert.Equal(t, "/media/movies/test.mkv", args.FilePath)
	assert.True(t, args.ForceRematch)
}

func TestMetadataRefreshArgs(t *testing.T) {
	id := uuid.New()
	args := MetadataRefreshArgs{
		ContentID: id,
		Force:     true,
	}

	assert.Equal(t, id, args.ContentID)
	assert.True(t, args.Force)
}

func TestSearchIndexArgs(t *testing.T) {
	t.Run("full reindex", func(t *testing.T) {
		args := SearchIndexArgs{
			ContentID:   nil,
			FullReindex: true,
		}

		assert.Nil(t, args.ContentID)
		assert.True(t, args.FullReindex)
	})

	t.Run("single item index", func(t *testing.T) {
		id := uuid.New()
		args := SearchIndexArgs{
			ContentID:   &id,
			FullReindex: false,
		}

		assert.NotNil(t, args.ContentID)
		assert.Equal(t, id, *args.ContentID)
		assert.False(t, args.FullReindex)
	})
}

func TestJobContext(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	t.Run("NewJobContext", func(t *testing.T) {
		jc := NewJobContext(ctx, logger, 123, "test_job")

		assert.NotNil(t, jc)
		assert.Equal(t, int64(123), jc.JobID)
		assert.Equal(t, "test_job", jc.JobKind)
		assert.NotNil(t, jc.Logger)
		assert.False(t, jc.StartTime.IsZero())
	})

	t.Run("Elapsed", func(t *testing.T) {
		jc := NewJobContext(ctx, logger, 123, "test_job")

		time.Sleep(10 * time.Millisecond)
		elapsed := jc.Elapsed()

		assert.GreaterOrEqual(t, elapsed, 10*time.Millisecond)
	})

	t.Run("LogStart", func(t *testing.T) {
		jc := NewJobContext(ctx, logger, 123, "test_job")

		// Should not panic
		jc.LogStart(zap.String("test", "value"))
	})

	t.Run("LogComplete", func(t *testing.T) {
		jc := NewJobContext(ctx, logger, 123, "test_job")

		// Should not panic
		jc.LogComplete(zap.Int("items", 10))
	})

	t.Run("LogError", func(t *testing.T) {
		jc := NewJobContext(ctx, logger, 123, "test_job")

		// Should not panic
		jc.LogError("test error", errors.New("something failed"), zap.String("context", "test"))
	})
}

func TestJobKind(t *testing.T) {
	tests := []struct {
		contentType string
		action      string
		expected    string
	}{
		{"movie", "library_scan", "movie_library_scan"},
		{"tvshow", "metadata_refresh", "tvshow_metadata_refresh"},
		{"music", "search_index", "music_search_index"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := JobKind(tt.contentType, tt.action)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestActionConstants(t *testing.T) {
	assert.Equal(t, "library_scan", ActionLibraryScan)
	assert.Equal(t, "file_match", ActionFileMatch)
	assert.Equal(t, "metadata_refresh", ActionMetadataRefresh)
	assert.Equal(t, "search_index", ActionSearchIndex)
	assert.Equal(t, "media_probe", ActionMediaProbe)
}
