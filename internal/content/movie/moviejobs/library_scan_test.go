package moviejobs

import (
	"context"
	"testing"
	"time"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/content/movie"
	infrajobs "github.com/lusoris/revenge/internal/infra/jobs"
)

// =============================================================================
// MovieLibraryScanArgs Tests
// =============================================================================

func TestMovieLibraryScanArgs_Kind(t *testing.T) {
	t.Parallel()

	args := MovieLibraryScanArgs{}
	assert.Equal(t, MovieLibraryScanJobKind, args.Kind())
	assert.Equal(t, "movie_library_scan", args.Kind())
}

func TestMovieLibraryScanArgs_InsertOpts(t *testing.T) {
	t.Parallel()

	args := MovieLibraryScanArgs{}
	opts := args.InsertOpts()
	assert.Equal(t, infrajobs.QueueBulk, opts.Queue)
}

func TestMovieLibraryScanArgs_Fields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		args          MovieLibraryScanArgs
		expectedPaths []string
		expectedForce bool
	}{
		{
			name: "multiple paths with force",
			args: MovieLibraryScanArgs{
				Paths: []string{"/movies", "/media/films"},
				Force: true,
			},
			expectedPaths: []string{"/movies", "/media/films"},
			expectedForce: true,
		},
		{
			name: "single path without force",
			args: MovieLibraryScanArgs{
				Paths: []string{"/movies"},
				Force: false,
			},
			expectedPaths: []string{"/movies"},
			expectedForce: false,
		},
		{
			name:          "empty paths",
			args:          MovieLibraryScanArgs{},
			expectedPaths: nil,
			expectedForce: false,
		},
		{
			name: "nil paths explicitly",
			args: MovieLibraryScanArgs{
				Paths: nil,
				Force: true,
			},
			expectedPaths: nil,
			expectedForce: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expectedPaths, tt.args.Paths)
			assert.Equal(t, tt.expectedForce, tt.args.Force)
		})
	}
}

func TestMovieLibraryScanJobKind_Constant(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "movie_library_scan", MovieLibraryScanJobKind)
	assert.NotEmpty(t, MovieLibraryScanJobKind)
}

// =============================================================================
// MovieLibraryScanWorker Constructor Tests
// =============================================================================

func TestNewMovieLibraryScanWorker(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()
	worker := NewMovieLibraryScanWorker(nil, logger)

	assert.NotNil(t, worker)
	assert.Nil(t, worker.libraryService)
	assert.NotNil(t, worker.logger)
}

func TestNewMovieLibraryScanWorker_WithNilLogger(t *testing.T) {
	t.Parallel()

	// Passing nil logger should still create worker (won't panic on construction).
	worker := NewMovieLibraryScanWorker(nil, nil)
	assert.NotNil(t, worker)
	assert.Nil(t, worker.libraryService)
	assert.Nil(t, worker.logger)
}

// =============================================================================
// MovieLibraryScanWorker.Kind() Tests
// =============================================================================

func TestMovieLibraryScanWorker_Kind(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()
	worker := NewMovieLibraryScanWorker(nil, logger)

	assert.Equal(t, MovieLibraryScanJobKind, worker.Kind())
	assert.Equal(t, "movie_library_scan", worker.Kind())
}

func TestMovieLibraryScanWorker_Kind_MatchesArgs(t *testing.T) {
	t.Parallel()

	worker := NewMovieLibraryScanWorker(nil, zap.NewNop())
	args := MovieLibraryScanArgs{}

	// Worker kind and args kind must match for River to route jobs correctly.
	assert.Equal(t, args.Kind(), worker.Kind())
}

// =============================================================================
// MovieLibraryScanWorker.Timeout() Tests
// =============================================================================

func TestMovieLibraryScanWorker_Timeout(t *testing.T) {
	t.Parallel()

	worker := NewMovieLibraryScanWorker(nil, zap.NewNop())

	job := &river.Job[MovieLibraryScanArgs]{
		JobRow: &rivertype.JobRow{ID: 1, Kind: MovieLibraryScanJobKind},
		Args: MovieLibraryScanArgs{
			Paths: []string{"/movies"},
			Force: false,
		},
	}

	assert.Equal(t, 30*time.Minute, worker.Timeout(job))
}

func TestMovieLibraryScanWorker_Timeout_ConsistentAcrossCalls(t *testing.T) {
	t.Parallel()

	worker := NewMovieLibraryScanWorker(nil, zap.NewNop())

	job1 := &river.Job[MovieLibraryScanArgs]{
		JobRow: &rivertype.JobRow{ID: 1, Kind: MovieLibraryScanJobKind},
		Args:   MovieLibraryScanArgs{Paths: []string{"/movies"}, Force: false},
	}
	job2 := &river.Job[MovieLibraryScanArgs]{
		JobRow: &rivertype.JobRow{ID: 2, Kind: MovieLibraryScanJobKind},
		Args:   MovieLibraryScanArgs{Paths: []string{"/media"}, Force: true},
	}

	// Timeout should be the same regardless of job args.
	assert.Equal(t, worker.Timeout(job1), worker.Timeout(job2))
}

// =============================================================================
// MovieLibraryScanWorker.Work() Tests
// =============================================================================

func TestMovieLibraryScanWorker_Work_NilLibraryService(t *testing.T) {
	t.Parallel()

	worker := NewMovieLibraryScanWorker(nil, zap.NewNop())

	job := &river.Job[MovieLibraryScanArgs]{
		JobRow: &rivertype.JobRow{ID: 1, Kind: MovieLibraryScanJobKind},
		Args: MovieLibraryScanArgs{
			Paths: []string{"/movies"},
			Force: false,
		},
	}

	// With a nil libraryService, calling Work should panic due to nil pointer dereference
	// on w.libraryService.ScanLibrary(). This confirms the worker requires a real service.
	assert.Panics(t, func() {
		_ = worker.Work(context.Background(), job)
	})
}

func TestMovieLibraryScanWorker_Work_EmptyLibrary(t *testing.T) {
	t.Parallel()

	// Create a real LibraryService with an empty temp directory.
	// The scanner will find no video files, so ScanLibrary returns
	// a summary with 0 files and no errors, exercising the full
	// success path of the Work method.
	tempDir := t.TempDir()

	libConfig := config.LibraryConfig{
		Paths: []string{tempDir},
	}
	libSvc := movie.NewLibraryService(nil, nil, libConfig, nil)
	worker := NewMovieLibraryScanWorker(libSvc, zap.NewNop())

	job := &river.Job[MovieLibraryScanArgs]{
		JobRow: &rivertype.JobRow{ID: 1, Kind: MovieLibraryScanJobKind},
		Args: MovieLibraryScanArgs{
			Paths: []string{tempDir},
			Force: false,
		},
	}

	err := worker.Work(context.Background(), job)
	require.NoError(t, err)
}
