package moviejobs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestMovieFileMatchArgs_Kind(t *testing.T) {
	t.Parallel()

	args := MovieFileMatchArgs{}
	assert.Equal(t, MovieFileMatchJobKind, args.Kind())
	assert.Equal(t, "movie_file_match", args.Kind())
}

func TestMovieFileMatchArgs_Fields(t *testing.T) {
	t.Parallel()

	args := MovieFileMatchArgs{
		FilePath:     "/media/movies/Inception (2010)/Inception.mkv",
		ForceRematch: true,
	}

	assert.Equal(t, "/media/movies/Inception (2010)/Inception.mkv", args.FilePath)
	assert.True(t, args.ForceRematch)
}

func TestNewMovieFileMatchWorker(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()
	worker := NewMovieFileMatchWorker(nil, logger)

	assert.NotNil(t, worker)
	assert.Nil(t, worker.libraryService)
	assert.NotNil(t, worker.logger)
}

func TestMovieFileMatchWorker_Kind(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()
	worker := NewMovieFileMatchWorker(nil, logger)

	assert.Equal(t, MovieFileMatchJobKind, worker.Kind())
	assert.Equal(t, "movie_file_match", worker.Kind())
}
