package api

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/content/tvshow"
	"github.com/lusoris/revenge/internal/infra/logging"
)

// bulkWatchedMockService is a minimal mock for tvshow.Service used by bulk watched tests.
type bulkWatchedMockService struct {
	tvshow.Service
	affected int64
	err      error
}

func (m *bulkWatchedMockService) MarkEpisodesWatchedBulk(_ context.Context, _ uuid.UUID, _ []uuid.UUID) (int64, error) {
	return m.affected, m.err
}

func TestHandler_MarkTVEpisodesBulkWatched_Success(t *testing.T) {
	t.Parallel()

	episodeIDs := []uuid.UUID{uuid.New(), uuid.New(), uuid.New()}
	svc := &bulkWatchedMockService{affected: 3}
	handler := &Handler{
		logger:        logging.NewTestLogger(),
		tvshowService: svc,
	}

	ctx := WithUserID(context.Background(), uuid.New())
	req := &ogen.BulkEpisodesWatchedRequest{EpisodeIds: episodeIDs}

	result, err := handler.MarkTVEpisodesBulkWatched(ctx, req)
	require.NoError(t, err)

	resp, ok := result.(*ogen.BulkEpisodesWatchedResponse)
	require.True(t, ok, "expected *ogen.BulkEpisodesWatchedResponse, got %T", result)
	assert.Equal(t, int64(3), resp.MarkedCount)
}

func TestHandler_MarkTVEpisodesBulkWatched_PartialMatch(t *testing.T) {
	t.Parallel()

	// 3 IDs sent, only 2 matched real episodes
	episodeIDs := []uuid.UUID{uuid.New(), uuid.New(), uuid.New()}
	svc := &bulkWatchedMockService{affected: 2}
	handler := &Handler{
		logger:        logging.NewTestLogger(),
		tvshowService: svc,
	}

	ctx := WithUserID(context.Background(), uuid.New())
	req := &ogen.BulkEpisodesWatchedRequest{EpisodeIds: episodeIDs}

	result, err := handler.MarkTVEpisodesBulkWatched(ctx, req)
	require.NoError(t, err)

	resp, ok := result.(*ogen.BulkEpisodesWatchedResponse)
	require.True(t, ok, "expected *ogen.BulkEpisodesWatchedResponse, got %T", result)
	assert.Equal(t, int64(2), resp.MarkedCount)
}

func TestHandler_MarkTVEpisodesBulkWatched_NoAuth(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: logging.NewTestLogger(),
	}

	ctx := context.Background()
	req := &ogen.BulkEpisodesWatchedRequest{EpisodeIds: []uuid.UUID{uuid.New()}}

	_, err := handler.MarkTVEpisodesBulkWatched(ctx, req)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrNoUserInContext)
}

func TestHandler_MarkTVEpisodesBulkWatched_ServiceError(t *testing.T) {
	t.Parallel()

	svc := &bulkWatchedMockService{err: errors.New("db error")}
	handler := &Handler{
		logger:        logging.NewTestLogger(),
		tvshowService: svc,
	}

	ctx := WithUserID(context.Background(), uuid.New())
	req := &ogen.BulkEpisodesWatchedRequest{EpisodeIds: []uuid.UUID{uuid.New()}}

	result, err := handler.MarkTVEpisodesBulkWatched(ctx, req)
	require.NoError(t, err)

	errResp, ok := result.(*ogen.Error)
	require.True(t, ok, "expected *ogen.Error, got %T", result)
	assert.Equal(t, 500, errResp.Code)
}
