package sonarr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/content/tvshow"
)

// mockTVShowRepo is a minimal mock of tvshow.Repository for webhook handler tests.
// Only the methods called during webhook handling are implemented; all others panic
// to surface unexpected calls immediately.
type mockTVShowRepo struct {
	tvshow.Repository // embed interface to satisfy compiler; unimplemented methods panic

	getSeriesBySonarrIDFn    func(ctx context.Context, sonarrID int32) (*tvshow.Series, error)
	listSeriesFn             func(ctx context.Context, filters tvshow.SeriesListFilters) ([]tvshow.Series, error)
	createSeriesFn           func(ctx context.Context, params tvshow.CreateSeriesParams) (*tvshow.Series, error)
	updateSeriesFn           func(ctx context.Context, params tvshow.UpdateSeriesParams) (*tvshow.Series, error)
	listSeasonsBySeriesFn    func(ctx context.Context, seriesID uuid.UUID) ([]tvshow.Season, error)
	createSeasonFn           func(ctx context.Context, params tvshow.CreateSeasonParams) (*tvshow.Season, error)
	listEpisodesBySeasonFn   func(ctx context.Context, seasonID uuid.UUID) ([]tvshow.Episode, error)
	createEpisodeFn          func(ctx context.Context, params tvshow.CreateEpisodeParams) (*tvshow.Episode, error)
	updateEpisodeFn          func(ctx context.Context, params tvshow.UpdateEpisodeParams) (*tvshow.Episode, error)
	getEpisodeFileBySonarrFn func(ctx context.Context, sonarrFileID int32) (*tvshow.EpisodeFile, error)
	createEpisodeFileFn      func(ctx context.Context, params tvshow.CreateEpisodeFileParams) (*tvshow.EpisodeFile, error)
	updateEpisodeFileFn      func(ctx context.Context, params tvshow.UpdateEpisodeFileParams) (*tvshow.EpisodeFile, error)
}

func (m *mockTVShowRepo) GetSeriesBySonarrID(ctx context.Context, sonarrID int32) (*tvshow.Series, error) {
	if m.getSeriesBySonarrIDFn != nil {
		return m.getSeriesBySonarrIDFn(ctx, sonarrID)
	}
	return nil, fmt.Errorf("series not found")
}

func (m *mockTVShowRepo) ListSeries(ctx context.Context, filters tvshow.SeriesListFilters) ([]tvshow.Series, error) {
	if m.listSeriesFn != nil {
		return m.listSeriesFn(ctx, filters)
	}
	return nil, nil
}

func (m *mockTVShowRepo) CreateSeries(ctx context.Context, params tvshow.CreateSeriesParams) (*tvshow.Series, error) {
	if m.createSeriesFn != nil {
		return m.createSeriesFn(ctx, params)
	}
	return &tvshow.Series{ID: uuid.Must(uuid.NewV7()), Title: "test"}, nil
}

func (m *mockTVShowRepo) UpdateSeries(ctx context.Context, params tvshow.UpdateSeriesParams) (*tvshow.Series, error) {
	if m.updateSeriesFn != nil {
		return m.updateSeriesFn(ctx, params)
	}
	return &tvshow.Series{ID: params.ID, Title: "test"}, nil
}

func (m *mockTVShowRepo) ListSeasonsBySeries(ctx context.Context, seriesID uuid.UUID) ([]tvshow.Season, error) {
	if m.listSeasonsBySeriesFn != nil {
		return m.listSeasonsBySeriesFn(ctx, seriesID)
	}
	return nil, nil
}

func (m *mockTVShowRepo) CreateSeason(ctx context.Context, params tvshow.CreateSeasonParams) (*tvshow.Season, error) {
	if m.createSeasonFn != nil {
		return m.createSeasonFn(ctx, params)
	}
	return &tvshow.Season{ID: uuid.Must(uuid.NewV7()), SeriesID: params.SeriesID, SeasonNumber: params.SeasonNumber}, nil
}

func (m *mockTVShowRepo) ListEpisodesBySeason(ctx context.Context, seasonID uuid.UUID) ([]tvshow.Episode, error) {
	if m.listEpisodesBySeasonFn != nil {
		return m.listEpisodesBySeasonFn(ctx, seasonID)
	}
	return nil, nil
}

func (m *mockTVShowRepo) CreateEpisode(ctx context.Context, params tvshow.CreateEpisodeParams) (*tvshow.Episode, error) {
	if m.createEpisodeFn != nil {
		return m.createEpisodeFn(ctx, params)
	}
	return &tvshow.Episode{ID: uuid.Must(uuid.NewV7()), SeriesID: params.SeriesID, SeasonID: params.SeasonID}, nil
}

func (m *mockTVShowRepo) UpdateEpisode(ctx context.Context, params tvshow.UpdateEpisodeParams) (*tvshow.Episode, error) {
	if m.updateEpisodeFn != nil {
		return m.updateEpisodeFn(ctx, params)
	}
	return &tvshow.Episode{ID: params.ID}, nil
}

func (m *mockTVShowRepo) GetEpisodeFileBySonarrID(ctx context.Context, sonarrFileID int32) (*tvshow.EpisodeFile, error) {
	if m.getEpisodeFileBySonarrFn != nil {
		return m.getEpisodeFileBySonarrFn(ctx, sonarrFileID)
	}
	return nil, fmt.Errorf("episode file not found")
}

func (m *mockTVShowRepo) CreateEpisodeFile(ctx context.Context, params tvshow.CreateEpisodeFileParams) (*tvshow.EpisodeFile, error) {
	if m.createEpisodeFileFn != nil {
		return m.createEpisodeFileFn(ctx, params)
	}
	return &tvshow.EpisodeFile{ID: uuid.Must(uuid.NewV7()), EpisodeID: params.EpisodeID}, nil
}

func (m *mockTVShowRepo) UpdateEpisodeFile(ctx context.Context, params tvshow.UpdateEpisodeFileParams) (*tvshow.EpisodeFile, error) {
	if m.updateEpisodeFileFn != nil {
		return m.updateEpisodeFileFn(ctx, params)
	}
	return &tvshow.EpisodeFile{ID: params.ID}, nil
}

// newTestSonarrServer creates an httptest.Server that responds with valid Sonarr API data.
// It handles GET /api/v3/series/{id} and GET /api/v3/episode?seriesId={id}.
func newTestSonarrServer(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v3/series/", func(w http.ResponseWriter, r *http.Request) {
		series := Series{
			ID:     42,
			Title:  "Test Series",
			TVDbID: 12345,
			Status: StatusContinuing,
			Year:   2020,
			Statistics: &Statistics{
				EpisodeFileCount: 1,
				EpisodeCount:     10,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(series)
	})

	mux.HandleFunc("/api/v3/episode", func(w http.ResponseWriter, r *http.Request) {
		episodes := []Episode{
			{
				ID:            1,
				SeriesID:      42,
				SeasonNumber:  1,
				EpisodeNumber: 1,
				Title:         "Pilot",
				HasFile:       true,
				EpisodeFile: &EpisodeFile{
					ID:           100,
					SeriesID:     42,
					SeasonNumber: 1,
					RelativePath: "Season 01/pilot.mkv",
					Path:         "/tv/Test Series/Season 01/pilot.mkv",
					Size:         1073741824,
					DateAdded:    time.Now(),
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(episodes)
	})

	return httptest.NewServer(mux)
}

// newTestSonarrServerError creates an httptest.Server that returns errors.
func newTestSonarrServerError(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}))
}

// newTestWebhookHandler creates a WebhookHandler backed by a mock Sonarr server and mock repo.
func newTestWebhookHandler(t *testing.T, server *httptest.Server, repo *mockTVShowRepo) *WebhookHandler {
	t.Helper()
	client := NewClient(Config{
		BaseURL:  server.URL,
		APIKey:   "test-key",
		CacheTTL: 1 * time.Second,
		Timeout:  5 * time.Second,
	})
	mapper := NewMapper()
	logger := slog.Default()

	syncService := NewSyncService(client, mapper, repo, logger)
	return NewWebhookHandler(syncService, logger)
}

func newPayload(eventType string, series *WebhookSeries) *WebhookPayload {
	return &WebhookPayload{
		EventType:    eventType,
		InstanceName: "test-sonarr",
		Series:       series,
	}
}

func defaultSeries() *WebhookSeries {
	return &WebhookSeries{
		ID:     42,
		Title:  "Test Series",
		TVDbID: 12345,
	}
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestWebhookHandler_NewWebhookHandler(t *testing.T) {
	server := newTestSonarrServer(t)
	defer server.Close()

	repo := &mockTVShowRepo{}
	handler := newTestWebhookHandler(t, server, repo)

	assert.NotNil(t, handler)
	assert.NotNil(t, handler.syncService)
	assert.NotNil(t, handler.logger)
}

func TestWebhookHandler_TestEvent(t *testing.T) {
	server := newTestSonarrServer(t)
	defer server.Close()

	repo := &mockTVShowRepo{}
	handler := newTestWebhookHandler(t, server, repo)

	err := handler.HandleWebhook(context.Background(), newPayload(EventTest, nil))
	assert.NoError(t, err)
}

func TestWebhookHandler_GrabEvent(t *testing.T) {
	server := newTestSonarrServer(t)
	defer server.Close()

	repo := &mockTVShowRepo{}
	handler := newTestWebhookHandler(t, server, repo)

	t.Run("with series info", func(t *testing.T) {
		payload := newPayload(EventGrab, defaultSeries())
		payload.Episodes = []WebhookEpisode{
			{ID: 1, SeasonNumber: 1, EpisodeNumber: 1, Title: "Pilot"},
		}
		payload.Release = &WebhookRelease{
			Quality:      "HDTV-1080p",
			ReleaseTitle: "Test.Series.S01E01.1080p.HDTV",
		}

		err := handler.HandleWebhook(context.Background(), payload)
		assert.NoError(t, err)
	})

	t.Run("without series info returns error", func(t *testing.T) {
		payload := newPayload(EventGrab, nil)

		err := handler.HandleWebhook(context.Background(), payload)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "grab event missing series info")
	})
}

func TestWebhookHandler_DownloadEvent(t *testing.T) {
	server := newTestSonarrServer(t)
	defer server.Close()

	t.Run("successful sync", func(t *testing.T) {
		syncCalled := false
		repo := &mockTVShowRepo{
			getSeriesBySonarrIDFn: func(ctx context.Context, sonarrID int32) (*tvshow.Series, error) {
				syncCalled = true
				return nil, fmt.Errorf("series not found")
			},
		}
		handler := newTestWebhookHandler(t, server, repo)

		payload := newPayload(EventDownload, defaultSeries())
		payload.Episodes = []WebhookEpisode{
			{ID: 1, SeasonNumber: 1, EpisodeNumber: 1, Title: "Pilot"},
		}

		err := handler.HandleWebhook(context.Background(), payload)
		assert.NoError(t, err)
		assert.True(t, syncCalled, "SyncSeries should have been called, which queries GetSeriesBySonarrID")
	})

	t.Run("without series info returns error", func(t *testing.T) {
		repo := &mockTVShowRepo{}
		handler := newTestWebhookHandler(t, server, repo)

		payload := newPayload(EventDownload, nil)

		err := handler.HandleWebhook(context.Background(), payload)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "download event missing series info")
	})

	t.Run("sync failure returns error", func(t *testing.T) {
		errServer := newTestSonarrServerError(t)
		defer errServer.Close()

		repo := &mockTVShowRepo{}
		handler := newTestWebhookHandler(t, errServer, repo)

		payload := newPayload(EventDownload, defaultSeries())
		err := handler.HandleWebhook(context.Background(), payload)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sync series")
	})

	t.Run("is upgrade flag", func(t *testing.T) {
		repo := &mockTVShowRepo{}
		handler := newTestWebhookHandler(t, server, repo)

		payload := newPayload(EventDownload, defaultSeries())
		payload.IsUpgrade = true
		payload.Episodes = []WebhookEpisode{
			{ID: 1, SeasonNumber: 1, EpisodeNumber: 1, Title: "Pilot"},
		}

		err := handler.HandleWebhook(context.Background(), payload)
		assert.NoError(t, err)
	})
}

func TestWebhookHandler_RenameEvent(t *testing.T) {
	server := newTestSonarrServer(t)
	defer server.Close()

	t.Run("successful sync", func(t *testing.T) {
		syncCalled := false
		repo := &mockTVShowRepo{
			getSeriesBySonarrIDFn: func(ctx context.Context, sonarrID int32) (*tvshow.Series, error) {
				syncCalled = true
				return nil, fmt.Errorf("series not found")
			},
		}
		handler := newTestWebhookHandler(t, server, repo)

		payload := newPayload(EventRename, defaultSeries())
		err := handler.HandleWebhook(context.Background(), payload)
		assert.NoError(t, err)
		assert.True(t, syncCalled)
	})

	t.Run("without series info returns error", func(t *testing.T) {
		repo := &mockTVShowRepo{}
		handler := newTestWebhookHandler(t, server, repo)

		payload := newPayload(EventRename, nil)
		err := handler.HandleWebhook(context.Background(), payload)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rename event missing series info")
	})
}

func TestWebhookHandler_SeriesAddEvent(t *testing.T) {
	server := newTestSonarrServer(t)
	defer server.Close()

	t.Run("successful sync", func(t *testing.T) {
		syncCalled := false
		repo := &mockTVShowRepo{
			getSeriesBySonarrIDFn: func(ctx context.Context, sonarrID int32) (*tvshow.Series, error) {
				syncCalled = true
				return nil, fmt.Errorf("series not found")
			},
		}
		handler := newTestWebhookHandler(t, server, repo)

		payload := newPayload(EventSeriesAdd, defaultSeries())
		err := handler.HandleWebhook(context.Background(), payload)
		assert.NoError(t, err)
		assert.True(t, syncCalled)
	})

	t.Run("without series info returns error", func(t *testing.T) {
		repo := &mockTVShowRepo{}
		handler := newTestWebhookHandler(t, server, repo)

		payload := newPayload(EventSeriesAdd, nil)
		err := handler.HandleWebhook(context.Background(), payload)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "series add event missing series info")
	})
}

func TestWebhookHandler_EpisodeFileDeleteEvent(t *testing.T) {
	server := newTestSonarrServer(t)
	defer server.Close()

	t.Run("successful sync", func(t *testing.T) {
		syncCalled := false
		repo := &mockTVShowRepo{
			getSeriesBySonarrIDFn: func(ctx context.Context, sonarrID int32) (*tvshow.Series, error) {
				syncCalled = true
				return nil, fmt.Errorf("series not found")
			},
		}
		handler := newTestWebhookHandler(t, server, repo)

		payload := newPayload(EventEpisodeFileDelete, defaultSeries())
		payload.DeletedFiles = []WebhookEpisodeFile{
			{ID: 100, RelativePath: "Season 01/episode.mkv"},
		}
		err := handler.HandleWebhook(context.Background(), payload)
		assert.NoError(t, err)
		assert.True(t, syncCalled)
	})

	t.Run("without series info returns error", func(t *testing.T) {
		repo := &mockTVShowRepo{}
		handler := newTestWebhookHandler(t, server, repo)

		payload := newPayload(EventEpisodeFileDelete, nil)
		err := handler.HandleWebhook(context.Background(), payload)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "episode file delete event missing series info")
	})
}

func TestWebhookHandler_SeriesDeleteEvent(t *testing.T) {
	server := newTestSonarrServer(t)
	defer server.Close()

	t.Run("series exists in local db", func(t *testing.T) {
		existingID := uuid.Must(uuid.NewV7())
		repo := &mockTVShowRepo{
			getSeriesBySonarrIDFn: func(ctx context.Context, sonarrID int32) (*tvshow.Series, error) {
				return &tvshow.Series{ID: existingID, Title: "Test Series"}, nil
			},
		}
		handler := newTestWebhookHandler(t, server, repo)

		payload := newPayload(EventSeriesDelete, defaultSeries())
		err := handler.HandleWebhook(context.Background(), payload)
		assert.NoError(t, err)
	})

	t.Run("series not in local db", func(t *testing.T) {
		repo := &mockTVShowRepo{
			getSeriesBySonarrIDFn: func(ctx context.Context, sonarrID int32) (*tvshow.Series, error) {
				return nil, fmt.Errorf("series not found")
			},
		}
		handler := newTestWebhookHandler(t, server, repo)

		payload := newPayload(EventSeriesDelete, defaultSeries())
		err := handler.HandleWebhook(context.Background(), payload)
		assert.NoError(t, err) // should not error, just log warning
	})

	t.Run("without series info returns error", func(t *testing.T) {
		repo := &mockTVShowRepo{}
		handler := newTestWebhookHandler(t, server, repo)

		payload := newPayload(EventSeriesDelete, nil)
		err := handler.HandleWebhook(context.Background(), payload)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "series delete event missing series info")
	})
}

func TestWebhookHandler_HealthEvent(t *testing.T) {
	server := newTestSonarrServer(t)
	defer server.Close()

	repo := &mockTVShowRepo{}
	handler := newTestWebhookHandler(t, server, repo)

	payload := newPayload(EventHealth, nil)
	payload.Level = "warning"
	payload.Message = "Indexer is down"
	payload.Type = "IndexerLongTermStatusCheck"
	payload.WikiURL = "https://wiki.servarr.com/sonarr/troubleshooting"

	err := handler.HandleWebhook(context.Background(), payload)
	assert.NoError(t, err)
}

func TestWebhookHandler_HealthRestoredEvent(t *testing.T) {
	server := newTestSonarrServer(t)
	defer server.Close()

	repo := &mockTVShowRepo{}
	handler := newTestWebhookHandler(t, server, repo)

	payload := newPayload(EventHealthRestored, nil)
	payload.Message = "Indexer restored"
	payload.Type = "IndexerLongTermStatusCheck"

	err := handler.HandleWebhook(context.Background(), payload)
	assert.NoError(t, err)
}

func TestWebhookHandler_ApplicationUpdateEvent(t *testing.T) {
	server := newTestSonarrServer(t)
	defer server.Close()

	repo := &mockTVShowRepo{}
	handler := newTestWebhookHandler(t, server, repo)

	payload := newPayload(EventApplicationUpdate, nil)
	payload.PreviousVersion = "3.0.9.1549"
	payload.NewVersion = "3.0.10.1567"

	err := handler.HandleWebhook(context.Background(), payload)
	assert.NoError(t, err)
}

func TestWebhookHandler_ManualInteractionRequiredEvent(t *testing.T) {
	server := newTestSonarrServer(t)
	defer server.Close()

	repo := &mockTVShowRepo{}
	handler := newTestWebhookHandler(t, server, repo)

	payload := newPayload(EventManualInteractionRequired, nil)
	payload.Message = "Manual import required"
	payload.DownloadClient = "qBittorrent"
	payload.DownloadID = "abc123"

	err := handler.HandleWebhook(context.Background(), payload)
	assert.NoError(t, err)
}

func TestWebhookHandler_UnknownEventType(t *testing.T) {
	server := newTestSonarrServer(t)
	defer server.Close()

	repo := &mockTVShowRepo{}
	handler := newTestWebhookHandler(t, server, repo)

	payload := newPayload("SomeNewEvent", nil)
	err := handler.HandleWebhook(context.Background(), payload)
	assert.NoError(t, err) // unknown events should not error
}

func TestWebhookHandler_SyncSeriesCalledForSyncEvents(t *testing.T) {
	// Events that should trigger SyncSeries: Download, EpisodeFileDelete, Rename, SeriesAdd
	syncEvents := []string{EventDownload, EventEpisodeFileDelete, EventRename, EventSeriesAdd}

	for _, eventType := range syncEvents {
		t.Run(eventType, func(t *testing.T) {
			server := newTestSonarrServer(t)
			defer server.Close()

			var syncSeriesCalledWith int32
			repo := &mockTVShowRepo{
				getSeriesBySonarrIDFn: func(ctx context.Context, sonarrID int32) (*tvshow.Series, error) {
					syncSeriesCalledWith = sonarrID
					return nil, fmt.Errorf("series not found")
				},
			}
			handler := newTestWebhookHandler(t, server, repo)

			payload := newPayload(eventType, defaultSeries())
			err := handler.HandleWebhook(context.Background(), payload)
			require.NoError(t, err)
			assert.Equal(t, int32(42), syncSeriesCalledWith, "SyncSeries should be called with correct Sonarr series ID")
		})
	}
}

func TestWebhookHandler_NoSyncForLogOnlyEvents(t *testing.T) {
	// Events that should NOT trigger SyncSeries
	logOnlyEvents := []struct {
		eventType string
		payload   *WebhookPayload
	}{
		{EventTest, newPayload(EventTest, nil)},
		{EventHealth, newPayload(EventHealth, nil)},
		{EventHealthRestored, newPayload(EventHealthRestored, nil)},
		{EventApplicationUpdate, newPayload(EventApplicationUpdate, nil)},
		{EventManualInteractionRequired, newPayload(EventManualInteractionRequired, nil)},
	}

	for _, tc := range logOnlyEvents {
		t.Run(tc.eventType, func(t *testing.T) {
			server := newTestSonarrServer(t)
			defer server.Close()

			syncCalled := false
			repo := &mockTVShowRepo{
				getSeriesBySonarrIDFn: func(ctx context.Context, sonarrID int32) (*tvshow.Series, error) {
					syncCalled = true
					return nil, fmt.Errorf("series not found")
				},
			}
			handler := newTestWebhookHandler(t, server, repo)

			err := handler.HandleWebhook(context.Background(), tc.payload)
			assert.NoError(t, err)
			assert.False(t, syncCalled, "SyncSeries should NOT be called for %s events", tc.eventType)
		})
	}
}

func TestWebhookHandler_GrabEventDoesNotSync(t *testing.T) {
	// Grab is a special case - it has a series, logs it, but does not trigger sync
	server := newTestSonarrServer(t)
	defer server.Close()

	syncCalled := false
	repo := &mockTVShowRepo{
		getSeriesBySonarrIDFn: func(ctx context.Context, sonarrID int32) (*tvshow.Series, error) {
			syncCalled = true
			return nil, fmt.Errorf("series not found")
		},
	}
	handler := newTestWebhookHandler(t, server, repo)

	payload := newPayload(EventGrab, defaultSeries())
	payload.Episodes = []WebhookEpisode{
		{ID: 1, SeasonNumber: 1, EpisodeNumber: 1, Title: "Pilot"},
	}

	err := handler.HandleWebhook(context.Background(), payload)
	assert.NoError(t, err)
	assert.False(t, syncCalled, "Grab event should not trigger SyncSeries")
}

func TestWebhookHandler_ContextCancelled(t *testing.T) {
	server := newTestSonarrServer(t)
	defer server.Close()

	repo := &mockTVShowRepo{}
	handler := newTestWebhookHandler(t, server, repo)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	// Log-only events still succeed because they don't use context for I/O
	err := handler.HandleWebhook(ctx, newPayload(EventTest, nil))
	assert.NoError(t, err)

	// Sync events should fail because the HTTP client will see the cancelled context
	payload := newPayload(EventDownload, defaultSeries())
	err = handler.HandleWebhook(ctx, payload)
	assert.Error(t, err)
}

func TestWebhookHandler_SyncSeriesUpdateExisting(t *testing.T) {
	server := newTestSonarrServer(t)
	defer server.Close()

	existingID := uuid.Must(uuid.NewV7())
	sonarrID := int32(42)
	updateCalled := false

	repo := &mockTVShowRepo{
		getSeriesBySonarrIDFn: func(ctx context.Context, id int32) (*tvshow.Series, error) {
			return &tvshow.Series{ID: existingID, SonarrID: &sonarrID, Title: "Existing Series"}, nil
		},
		updateSeriesFn: func(ctx context.Context, params tvshow.UpdateSeriesParams) (*tvshow.Series, error) {
			updateCalled = true
			assert.Equal(t, existingID, params.ID)
			return &tvshow.Series{ID: existingID}, nil
		},
	}
	handler := newTestWebhookHandler(t, server, repo)

	payload := newPayload(EventDownload, defaultSeries())
	err := handler.HandleWebhook(context.Background(), payload)
	assert.NoError(t, err)
	assert.True(t, updateCalled, "should update existing series instead of creating")
}

func TestWebhookHandler_SeriesDeleteEventChecksCorrectSonarrID(t *testing.T) {
	server := newTestSonarrServer(t)
	defer server.Close()

	var queriedSonarrID int32
	repo := &mockTVShowRepo{
		getSeriesBySonarrIDFn: func(ctx context.Context, sonarrID int32) (*tvshow.Series, error) {
			queriedSonarrID = sonarrID
			return nil, errors.New("series not found")
		},
	}
	handler := newTestWebhookHandler(t, server, repo)

	series := &WebhookSeries{ID: 99, Title: "Another Series"}
	payload := newPayload(EventSeriesDelete, series)

	err := handler.HandleWebhook(context.Background(), payload)
	assert.NoError(t, err)
	assert.Equal(t, int32(99), queriedSonarrID, "should query local DB with the correct Sonarr series ID")
}

func TestWebhookHandler_AllEventTypes(t *testing.T) {
	// Exhaustive test that every known event type routes to the correct handler
	// and returns no error when given valid input.
	server := newTestSonarrServer(t)
	defer server.Close()

	repo := &mockTVShowRepo{}
	handler := newTestWebhookHandler(t, server, repo)

	tests := []struct {
		name      string
		eventType string
		series    *WebhookSeries
		wantErr   bool
	}{
		{"Test", EventTest, nil, false},
		{"Grab", EventGrab, defaultSeries(), false},
		{"Download", EventDownload, defaultSeries(), false},
		{"Rename", EventRename, defaultSeries(), false},
		{"SeriesAdd", EventSeriesAdd, defaultSeries(), false},
		{"SeriesDelete", EventSeriesDelete, defaultSeries(), false},
		{"EpisodeFileDelete", EventEpisodeFileDelete, defaultSeries(), false},
		{"Health", EventHealth, nil, false},
		{"HealthRestored", EventHealthRestored, nil, false},
		{"ApplicationUpdate", EventApplicationUpdate, nil, false},
		{"ManualInteractionRequired", EventManualInteractionRequired, nil, false},
		{"Unknown", "FutureEvent", nil, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			payload := newPayload(tc.eventType, tc.series)
			err := handler.HandleWebhook(context.Background(), payload)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestWebhookHandler_NilSeriesForAllSeriesRequiredEvents(t *testing.T) {
	// All events that require a series should return a descriptive error when series is nil.
	server := newTestSonarrServer(t)
	defer server.Close()

	repo := &mockTVShowRepo{}
	handler := newTestWebhookHandler(t, server, repo)

	seriesRequiredEvents := []struct {
		eventType   string
		errContains string
	}{
		{EventGrab, "grab event missing series info"},
		{EventDownload, "download event missing series info"},
		{EventRename, "rename event missing series info"},
		{EventSeriesAdd, "series add event missing series info"},
		{EventSeriesDelete, "series delete event missing series info"},
		{EventEpisodeFileDelete, "episode file delete event missing series info"},
	}

	for _, tc := range seriesRequiredEvents {
		t.Run(tc.eventType, func(t *testing.T) {
			payload := newPayload(tc.eventType, nil)
			err := handler.HandleWebhook(context.Background(), payload)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.errContains)
		})
	}
}

func TestWebhookHandler_CreateSeriesFailure(t *testing.T) {
	server := newTestSonarrServer(t)
	defer server.Close()

	repo := &mockTVShowRepo{
		createSeriesFn: func(ctx context.Context, params tvshow.CreateSeriesParams) (*tvshow.Series, error) {
			return nil, fmt.Errorf("db connection failed")
		},
	}
	handler := newTestWebhookHandler(t, server, repo)

	payload := newPayload(EventDownload, defaultSeries())
	err := handler.HandleWebhook(context.Background(), payload)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sync series")
}
