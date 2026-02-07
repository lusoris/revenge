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
)

// =============================================================================
// MovieFileMatchArgs Tests
// =============================================================================

func TestMovieFileMatchArgs_Kind(t *testing.T) {
	t.Parallel()

	args := MovieFileMatchArgs{}
	assert.Equal(t, MovieFileMatchJobKind, args.Kind())
	assert.Equal(t, "movie_file_match", args.Kind())
}

func TestMovieFileMatchJobKind_Constant(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "movie_file_match", MovieFileMatchJobKind)
	assert.NotEmpty(t, MovieFileMatchJobKind)
}

func TestMovieFileMatchArgs_Fields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		args               MovieFileMatchArgs
		expectedPath       string
		expectedForceMatch bool
	}{
		{
			name: "with force rematch",
			args: MovieFileMatchArgs{
				FilePath:     "/media/movies/Inception (2010)/Inception.mkv",
				ForceRematch: true,
			},
			expectedPath:       "/media/movies/Inception (2010)/Inception.mkv",
			expectedForceMatch: true,
		},
		{
			name: "without force rematch",
			args: MovieFileMatchArgs{
				FilePath:     "/media/movies/The Matrix (1999)/The.Matrix.1999.mkv",
				ForceRematch: false,
			},
			expectedPath:       "/media/movies/The Matrix (1999)/The.Matrix.1999.mkv",
			expectedForceMatch: false,
		},
		{
			name:               "zero value",
			args:               MovieFileMatchArgs{},
			expectedPath:       "",
			expectedForceMatch: false,
		},
		{
			name: "path with special characters",
			args: MovieFileMatchArgs{
				FilePath:     "/media/movies/Amélie (2001)/Amélie.mkv",
				ForceRematch: false,
			},
			expectedPath:       "/media/movies/Amélie (2001)/Amélie.mkv",
			expectedForceMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expectedPath, tt.args.FilePath)
			assert.Equal(t, tt.expectedForceMatch, tt.args.ForceRematch)
		})
	}
}

// =============================================================================
// MovieFileMatchWorker Constructor Tests
// =============================================================================

func TestNewMovieFileMatchWorker(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()
	worker := NewMovieFileMatchWorker(nil, logger)

	assert.NotNil(t, worker)
	assert.Nil(t, worker.libraryService)
	assert.NotNil(t, worker.logger)
}

func TestNewMovieFileMatchWorker_NilLogger(t *testing.T) {
	t.Parallel()

	worker := NewMovieFileMatchWorker(nil, nil)
	assert.NotNil(t, worker)
	assert.Nil(t, worker.libraryService)
	assert.Nil(t, worker.logger)
}

// =============================================================================
// MovieFileMatchWorker.Kind() Tests
// =============================================================================

func TestMovieFileMatchWorker_Kind(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()
	worker := NewMovieFileMatchWorker(nil, logger)

	assert.Equal(t, MovieFileMatchJobKind, worker.Kind())
	assert.Equal(t, "movie_file_match", worker.Kind())
}

func TestMovieFileMatchWorker_Kind_MatchesArgs(t *testing.T) {
	t.Parallel()

	worker := NewMovieFileMatchWorker(nil, zap.NewNop())
	args := MovieFileMatchArgs{}

	// Worker kind and args kind must match for River to route jobs correctly.
	assert.Equal(t, args.Kind(), worker.Kind())
}

// =============================================================================
// MovieFileMatchWorker.Timeout() Tests
// =============================================================================

func TestMovieFileMatchWorker_Timeout(t *testing.T) {
	t.Parallel()

	worker := NewMovieFileMatchWorker(nil, zap.NewNop())

	job := &river.Job[MovieFileMatchArgs]{
		JobRow: &rivertype.JobRow{ID: 1, Kind: MovieFileMatchJobKind},
		Args: MovieFileMatchArgs{
			FilePath:     "/media/movies/test.mkv",
			ForceRematch: false,
		},
	}

	assert.Equal(t, 5*time.Minute, worker.Timeout(job))
}

// =============================================================================
// MovieFileMatchWorker.Work() Tests
// =============================================================================

func TestMovieFileMatchWorker_Work_NilLibraryService_NonexistentFile(t *testing.T) {
	t.Parallel()

	worker := NewMovieFileMatchWorker(nil, zap.NewNop())

	job := &river.Job[MovieFileMatchArgs]{
		JobRow: &rivertype.JobRow{ID: 1, Kind: MovieFileMatchJobKind},
		Args: MovieFileMatchArgs{
			FilePath:     "/nonexistent/path/movie.mkv",
			ForceRematch: false,
		},
	}

	// MatchFile on a nil *LibraryService first calls os.Stat which returns an error
	// for a nonexistent file path, returning "file not found" before accessing any
	// struct fields. This confirms the error propagation path.
	err := worker.Work(context.Background(), job)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "file not found")
}
