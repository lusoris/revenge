package moviejobs

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestMovieSearchIndexArgsKind(t *testing.T) {
	args := MovieSearchIndexArgs{}
	assert.Equal(t, "movie_search_index", args.Kind())
}

func TestSearchIndexOperationConstants(t *testing.T) {
	assert.Equal(t, SearchIndexOperation("index"), SearchIndexOperationIndex)
	assert.Equal(t, SearchIndexOperation("remove"), SearchIndexOperationRemove)
	assert.Equal(t, SearchIndexOperation("reindex"), SearchIndexOperationReindex)
}

func TestMovieSearchIndexArgs(t *testing.T) {
	movieID := uuid.New()

	tests := []struct {
		name      string
		operation SearchIndexOperation
		movieID   uuid.UUID
	}{
		{
			name:      "index operation",
			operation: SearchIndexOperationIndex,
			movieID:   movieID,
		},
		{
			name:      "remove operation",
			operation: SearchIndexOperationRemove,
			movieID:   movieID,
		},
		{
			name:      "reindex operation",
			operation: SearchIndexOperationReindex,
			movieID:   uuid.Nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := MovieSearchIndexArgs{
				Operation: tt.operation,
				MovieID:   tt.movieID,
			}

			assert.Equal(t, tt.operation, args.Operation)
			assert.Equal(t, tt.movieID, args.MovieID)
			assert.Equal(t, "movie_search_index", args.Kind())
		})
	}
}

func TestNewMovieSearchIndexWorker(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()
	worker := NewMovieSearchIndexWorker(nil, nil, logger)

	assert.NotNil(t, worker)
	assert.Nil(t, worker.movieRepo)
	assert.Nil(t, worker.searchService)
	assert.NotNil(t, worker.logger)
}
