package api

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/lusoris/revenge/internal/infra/logging"

	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/content/movie/moviejobs"
)

// ============================================================================
// ReindexSearch Tests
// ============================================================================

func TestHandler_ReindexSearch_NilRiverClient(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger:      logging.NewTestLogger(),
		riverClient: nil,
	}

	result, err := handler.ReindexSearch(context.Background())
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "job queue not available")
}

func TestHandler_ReindexSearch_Success(t *testing.T) {
	t.Parallel()

	mockRiver := &mockRiverClient{}
	handler := &Handler{
		logger:      logging.NewTestLogger(),
		riverClient: mockRiver,
	}

	result, err := handler.ReindexSearch(context.Background())
	require.NoError(t, err)
	require.NotNil(t, result)

	accepted, ok := result.(*ogen.ReindexSearchAccepted)
	require.True(t, ok)
	assert.Equal(t, "Reindex job enqueued", accepted.Message.Value)
	assert.True(t, accepted.JobID.Set, "JobID should be set")

	// Verify job was inserted with correct args
	require.Len(t, mockRiver.insertedArgs, 1)
	indexArgs, ok := mockRiver.insertedArgs[0].(moviejobs.MovieSearchIndexArgs)
	require.True(t, ok)
	assert.Equal(t, moviejobs.SearchIndexOperationReindex, indexArgs.Operation)
}

func TestHandler_ReindexSearch_InsertError(t *testing.T) {
	t.Parallel()

	mockRiver := &mockRiverClient{
		insertError: errors.New("queue connection failed"),
	}
	handler := &Handler{
		logger:      logging.NewTestLogger(),
		riverClient: mockRiver,
	}

	result, err := handler.ReindexSearch(context.Background())
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to enqueue reindex job")
	assert.Contains(t, err.Error(), "queue connection failed")
}
