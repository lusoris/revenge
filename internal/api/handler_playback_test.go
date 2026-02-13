package api

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/infra/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/playback"
)

// ============================================================================
// StartPlaybackSession Tests
// ============================================================================

func TestHandler_StartPlaybackSession_Disabled(t *testing.T) {
	t.Parallel()

	handler := &Handler{logger: logging.NewTestLogger()}

	req := &ogen.StartPlaybackRequest{
		MediaType: "movie",
		MediaID:   uuid.New(),
	}

	result, err := handler.StartPlaybackSession(context.Background(), req)
	require.NoError(t, err)

	badReq, ok := result.(*ogen.StartPlaybackSessionBadRequest)
	require.True(t, ok, "expected *ogen.StartPlaybackSessionBadRequest, got %T", result)
	assert.Equal(t, 400, badReq.Code)
	assert.Equal(t, "Playback is not enabled", badReq.Message)
}

func TestHandler_StartPlaybackSession_Unauthorized(t *testing.T) {
	t.Parallel()

	// Use a non-nil zero-value playback.Service so the nil check passes.
	// The handler calls getUserID(ctx) next, which returns false for a bare
	// context.Background(), so no Service methods are ever invoked.
	handler := &Handler{
		logger:          logging.NewTestLogger(),
		playbackService: new(playback.Service),
	}

	req := &ogen.StartPlaybackRequest{
		MediaType: "movie",
		MediaID:   uuid.New(),
	}

	result, err := handler.StartPlaybackSession(context.Background(), req)
	require.NoError(t, err)

	unauth, ok := result.(*ogen.StartPlaybackSessionUnauthorized)
	require.True(t, ok, "expected *ogen.StartPlaybackSessionUnauthorized, got %T", result)
	assert.Equal(t, 401, unauth.Code)
	assert.Equal(t, "Authentication required", unauth.Message)
}

// ============================================================================
// GetPlaybackSession Tests
// ============================================================================

func TestHandler_GetPlaybackSession_Disabled(t *testing.T) {
	t.Parallel()

	handler := &Handler{logger: logging.NewTestLogger()}

	params := ogen.GetPlaybackSessionParams{SessionId: uuid.New()}

	result, err := handler.GetPlaybackSession(context.Background(), params)
	require.NoError(t, err)

	notFound, ok := result.(*ogen.GetPlaybackSessionNotFound)
	require.True(t, ok, "expected *ogen.GetPlaybackSessionNotFound, got %T", result)
	assert.Equal(t, 404, notFound.Code)
	assert.Equal(t, "Playback is not enabled", notFound.Message)
}

func TestHandler_GetPlaybackSession_Unauthorized(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger:          logging.NewTestLogger(),
		playbackService: new(playback.Service),
	}

	params := ogen.GetPlaybackSessionParams{SessionId: uuid.New()}

	result, err := handler.GetPlaybackSession(context.Background(), params)
	require.NoError(t, err)

	unauth, ok := result.(*ogen.GetPlaybackSessionUnauthorized)
	require.True(t, ok, "expected *ogen.GetPlaybackSessionUnauthorized, got %T", result)
	assert.Equal(t, 401, unauth.Code)
	assert.Equal(t, "Authentication required", unauth.Message)
}

// ============================================================================
// StopPlaybackSession Tests
// ============================================================================

func TestHandler_StopPlaybackSession_Disabled(t *testing.T) {
	t.Parallel()

	handler := &Handler{logger: logging.NewTestLogger()}

	params := ogen.StopPlaybackSessionParams{SessionId: uuid.New()}

	result, err := handler.StopPlaybackSession(context.Background(), params)
	require.NoError(t, err)

	notFound, ok := result.(*ogen.StopPlaybackSessionNotFound)
	require.True(t, ok, "expected *ogen.StopPlaybackSessionNotFound, got %T", result)
	assert.Equal(t, 404, notFound.Code)
	assert.Equal(t, "Playback is not enabled", notFound.Message)
}

func TestHandler_StopPlaybackSession_Unauthorized(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger:          logging.NewTestLogger(),
		playbackService: new(playback.Service),
	}

	params := ogen.StopPlaybackSessionParams{SessionId: uuid.New()}

	result, err := handler.StopPlaybackSession(context.Background(), params)
	require.NoError(t, err)

	unauth, ok := result.(*ogen.StopPlaybackSessionUnauthorized)
	require.True(t, ok, "expected *ogen.StopPlaybackSessionUnauthorized, got %T", result)
	assert.Equal(t, 401, unauth.Code)
	assert.Equal(t, "Authentication required", unauth.Message)
}
