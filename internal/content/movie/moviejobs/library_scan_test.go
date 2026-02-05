package moviejobs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestMovieLibraryScanArgs_Kind(t *testing.T) {
	t.Parallel()

	args := MovieLibraryScanArgs{}
	assert.Equal(t, MovieLibraryScanJobKind, args.Kind())
	assert.Equal(t, "movie_library_scan", args.Kind())
}

func TestMovieLibraryScanArgs_Fields(t *testing.T) {
	t.Parallel()

	args := MovieLibraryScanArgs{
		Paths: []string{"/movies", "/media/films"},
		Force: true,
	}

	assert.Len(t, args.Paths, 2)
	assert.Contains(t, args.Paths, "/movies")
	assert.Contains(t, args.Paths, "/media/films")
	assert.True(t, args.Force)
}

func TestNewMovieLibraryScanWorker(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()
	worker := NewMovieLibraryScanWorker(nil, logger)

	assert.NotNil(t, worker)
	assert.Nil(t, worker.libraryService)
	assert.NotNil(t, worker.logger)
}

func TestMovieLibraryScanWorker_Kind(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()
	worker := NewMovieLibraryScanWorker(nil, logger)

	assert.Equal(t, MovieLibraryScanJobKind, worker.Kind())
	assert.Equal(t, "movie_library_scan", worker.Kind())
}
