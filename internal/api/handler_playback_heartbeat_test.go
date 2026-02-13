package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/lusoris/revenge/internal/infra/logging"
	"github.com/lusoris/revenge/internal/playback"
)

// ============================================================================
// HeartbeatPlaybackSession Tests (custom mux handler)
// ============================================================================

func TestHandler_Heartbeat_NoAuth(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger:          logging.NewTestLogger(),
		playbackService: new(playback.Service),
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/playback/sessions/"+uuid.New().String()+"/heartbeat", nil)
	w := httptest.NewRecorder()

	handler.heartbeatHandler().ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Authentication required")
}

func TestHandler_Heartbeat_InvalidSessionID(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger:          logging.NewTestLogger(),
		playbackService: new(playback.Service),
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/playback/sessions/invalid/heartbeat", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	req.SetPathValue("sessionId", "invalid")
	w := httptest.NewRecorder()

	// Falls through to token validation which fails (nil tokenManager)
	handler.heartbeatHandler().ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Authentication not configured")
}

func TestHandler_Heartbeat_WithBody(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger:          logging.NewTestLogger(),
		playbackService: new(playback.Service),
	}

	body := bytes.NewBufferString(`{"position_seconds": 120}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/playback/sessions/"+uuid.New().String()+"/heartbeat", body)
	req.Header.Set("Authorization", "Bearer test-token")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Will fail at token validation (no real tokenManager)
	handler.heartbeatHandler().ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestHandler_Heartbeat_NilPlaybackService(t *testing.T) {
	t.Parallel()

	// When playbackService is nil, heartbeatHandler() is never mounted on the mux.
	// But if someone calls it directly, the auth check happens first.
	handler := &Handler{
		logger: logging.NewTestLogger(),
		// no playbackService
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/playback/sessions/"+uuid.New().String()+"/heartbeat", nil)
	w := httptest.NewRecorder()

	handler.heartbeatHandler().ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
