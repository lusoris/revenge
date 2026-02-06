package sonarr

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

func TestNewClient(t *testing.T) {
	client := NewClient(Config{
		BaseURL: "http://localhost:8989",
		APIKey:  "test-api-key",
	})

	assert.NotNil(t, client)
	assert.Equal(t, "http://localhost:8989", client.baseURL)
	assert.Equal(t, "test-api-key", client.apiKey)
}

func TestNewClient_Defaults(t *testing.T) {
	client := NewClient(Config{
		BaseURL: "http://localhost:8989",
		APIKey:  "test-api-key",
	})

	assert.NotNil(t, client.rateLimiter)
	assert.Equal(t, 5*time.Minute, client.cacheTTL)
}

func TestNewClient_CustomConfig(t *testing.T) {
	client := NewClient(Config{
		BaseURL:   "http://sonarr.local:8989",
		APIKey:    "custom-key",
		RateLimit: rate.Limit(5.0),
		CacheTTL:  10 * time.Minute,
		Timeout:   60 * time.Second,
	})

	assert.NotNil(t, client)
	assert.Equal(t, "http://sonarr.local:8989", client.baseURL)
	assert.Equal(t, "custom-key", client.apiKey)
	assert.Equal(t, 10*time.Minute, client.cacheTTL)
}

func TestClient_Cache(t *testing.T) {
	client := NewClient(Config{
		BaseURL:  "http://localhost:8989",
		APIKey:   "test-api-key",
		CacheTTL: 1 * time.Second,
	})

	// Set a cache entry
	client.setCache("test-key", "test-value")

	// Get it back
	result := client.getFromCache("test-key")
	assert.Equal(t, "test-value", result)

	// Wait for expiration
	time.Sleep(1100 * time.Millisecond)

	// Should be nil after expiration
	result = client.getFromCache("test-key")
	assert.Nil(t, result)
}

func TestClient_ClearCache(t *testing.T) {
	client := NewClient(Config{
		BaseURL:  "http://localhost:8989",
		APIKey:   "test-api-key",
		CacheTTL: 1 * time.Hour,
	})

	// Set cache entries
	client.setCache("key1", "value1")
	client.setCache("key2", "value2")

	// Verify they exist
	assert.NotNil(t, client.getFromCache("key1"))
	assert.NotNil(t, client.getFromCache("key2"))

	// Clear cache
	client.ClearCache()

	// Should be nil after clear
	assert.Nil(t, client.getFromCache("key1"))
	assert.Nil(t, client.getFromCache("key2"))
}

func TestErrors(t *testing.T) {
	assert.Error(t, ErrSeriesNotFound)
	assert.Error(t, ErrEpisodeNotFound)
	assert.Error(t, ErrEpisodeFileNotFound)
	assert.Error(t, ErrUnauthorized)
	assert.Error(t, ErrConnectionFailed)
	assert.Error(t, ErrRateLimited)

	assert.Contains(t, ErrSeriesNotFound.Error(), "series not found")
	assert.Contains(t, ErrEpisodeNotFound.Error(), "episode not found")
}

func TestSeriesTypes(t *testing.T) {
	series := Series{
		ID:       1,
		Title:    "Breaking Bad",
		TVDbID:   81189,
		Status:   StatusEnded,
		Ended:    true,
		Year:     2008,
		Seasons: []SeasonInfo{
			{SeasonNumber: 1, Monitored: true},
			{SeasonNumber: 2, Monitored: true},
		},
	}

	assert.Equal(t, 1, series.ID)
	assert.Equal(t, "Breaking Bad", series.Title)
	assert.Equal(t, 81189, series.TVDbID)
	assert.Equal(t, StatusEnded, series.Status)
	assert.True(t, series.Ended)
	assert.Equal(t, 2008, series.Year)
	assert.Len(t, series.Seasons, 2)
	assert.Equal(t, 1, series.Seasons[0].SeasonNumber)
	assert.True(t, series.Seasons[0].Monitored)
}

func TestEpisodeTypes(t *testing.T) {
	episode := Episode{
		ID:            1,
		SeriesID:      1,
		SeasonNumber:  1,
		EpisodeNumber: 1,
		Title:         "Pilot",
		HasFile:       true,
		Monitored:     true,
	}

	assert.Equal(t, 1, episode.ID)
	assert.Equal(t, 1, episode.SeriesID)
	assert.Equal(t, 1, episode.SeasonNumber)
	assert.Equal(t, 1, episode.EpisodeNumber)
	assert.Equal(t, "Pilot", episode.Title)
	assert.True(t, episode.HasFile)
	assert.True(t, episode.Monitored)
}

func TestWebhookPayload(t *testing.T) {
	payload := WebhookPayload{
		EventType: EventDownload,
		Series: &WebhookSeries{
			ID:     1,
			Title:  "Breaking Bad",
			TVDbID: 81189,
		},
		Episodes: []WebhookEpisode{
			{
				ID:            1,
				SeasonNumber:  1,
				EpisodeNumber: 1,
				Title:         "Pilot",
			},
		},
		IsUpgrade: false,
	}

	assert.Equal(t, EventDownload, payload.EventType)
	assert.NotNil(t, payload.Series)
	assert.Equal(t, 1, payload.Series.ID)
	assert.Equal(t, "Breaking Bad", payload.Series.Title)
	assert.Equal(t, 81189, payload.Series.TVDbID)
	assert.Len(t, payload.Episodes, 1)
	assert.Equal(t, 1, payload.Episodes[0].ID)
	assert.Equal(t, 1, payload.Episodes[0].SeasonNumber)
	assert.Equal(t, 1, payload.Episodes[0].EpisodeNumber)
	assert.Equal(t, "Pilot", payload.Episodes[0].Title)
	assert.False(t, payload.IsUpgrade)
}

func TestEventConstants(t *testing.T) {
	events := []string{
		EventGrab,
		EventDownload,
		EventRename,
		EventSeriesAdd,
		EventSeriesDelete,
		EventEpisodeFileDelete,
		EventHealth,
		EventHealthRestored,
		EventApplicationUpdate,
		EventManualInteractionRequired,
		EventTest,
	}

	// Verify all events are unique
	seen := make(map[string]bool)
	for _, event := range events {
		assert.False(t, seen[event], "duplicate event: %s", event)
		seen[event] = true
		assert.NotEmpty(t, event)
	}
}

func TestStatusConstants(t *testing.T) {
	assert.Equal(t, "continuing", StatusContinuing)
	assert.Equal(t, "ended", StatusEnded)
	assert.Equal(t, "upcoming", StatusUpcoming)
	assert.Equal(t, "deleted", StatusDeleted)
}

func TestSeriesTypeConstants(t *testing.T) {
	assert.Equal(t, "standard", SeriesTypeStandard)
	assert.Equal(t, "daily", SeriesTypeDaily)
	assert.Equal(t, "anime", SeriesTypeAnime)
}

func TestMonitorConstants(t *testing.T) {
	monitors := []string{
		MonitorAll,
		MonitorFuture,
		MonitorMissing,
		MonitorExisting,
		MonitorPilot,
		MonitorFirstSeason,
		MonitorLastSeason,
		MonitorMonitorSpecials,
		MonitorUnmonitorSpecials,
		MonitorNone,
	}

	for _, monitor := range monitors {
		assert.NotEmpty(t, monitor)
	}
}

func TestAddSeriesRequest(t *testing.T) {
	req := AddSeriesRequest{
		Title:            "Breaking Bad",
		QualityProfileID: 1,
		TVDbID:           81189,
		RootFolderPath:   "/tv",
		Monitored:        true,
		SeasonFolder:     true,
		SeriesType:       SeriesTypeStandard,
		AddOptions: AddSeriesOptions{
			Monitor:                  MonitorAll,
			SearchForMissingEpisodes: true,
		},
	}

	assert.Equal(t, "Breaking Bad", req.Title)
	assert.Equal(t, 1, req.QualityProfileID)
	assert.Equal(t, 81189, req.TVDbID)
	assert.Equal(t, "/tv", req.RootFolderPath)
	assert.True(t, req.Monitored)
	assert.True(t, req.SeasonFolder)
	assert.Equal(t, SeriesTypeStandard, req.SeriesType)
	assert.Equal(t, MonitorAll, req.AddOptions.Monitor)
	assert.True(t, req.AddOptions.SearchForMissingEpisodes)
}
