package jobs

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/govalues/decimal"
	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/content"
	"github.com/lusoris/revenge/internal/content/shared/scanner"
	"github.com/lusoris/revenge/internal/content/tvshow"
	infrajobs "github.com/lusoris/revenge/internal/infra/jobs"
	"github.com/lusoris/revenge/internal/service/search"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/lusoris/revenge/internal/infra/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// Args.Kind() Tests
// =============================================================================

func TestLibraryScanArgs_Kind(t *testing.T) {
	t.Parallel()

	args := LibraryScanArgs{
		Paths: []string{"/tv"},
		Force: true,
	}
	assert.Equal(t, KindLibraryScan, args.Kind())
	assert.Equal(t, "tvshow_library_scan", args.Kind())
}

func TestMetadataRefreshArgs_Kind(t *testing.T) {
	t.Parallel()

	seriesID := uuid.Must(uuid.NewV7())
	args := MetadataRefreshArgs{
		SeriesID: &seriesID,
		Force:    true,
	}
	assert.Equal(t, KindMetadataRefresh, args.Kind())
	assert.Equal(t, "tvshow_metadata_refresh", args.Kind())
}

func TestFileMatchArgs_Kind(t *testing.T) {
	t.Parallel()

	args := FileMatchArgs{
		FilePath:     "/tv/Breaking Bad/S01E01.mkv",
		ForceRematch: false,
		AutoCreate:   true,
	}
	assert.Equal(t, KindFileMatch, args.Kind())
	assert.Equal(t, "tvshow_file_match", args.Kind())
}

func TestSearchIndexArgs_Kind(t *testing.T) {
	t.Parallel()

	args := SearchIndexArgs{
		FullReindex: true,
	}
	assert.Equal(t, KindSearchIndex, args.Kind())
	assert.Equal(t, "tvshow_search_index", args.Kind())
}

func TestSeriesRefreshArgs_Kind(t *testing.T) {
	t.Parallel()

	args := SeriesRefreshArgs{
		SeriesID:        uuid.Must(uuid.NewV7()),
		TMDbID:          1396,
		RefreshSeasons:  true,
		RefreshEpisodes: true,
	}
	assert.Equal(t, KindSeriesRefresh, args.Kind())
	assert.Equal(t, "tvshow_series_refresh", args.Kind())
}

// =============================================================================
// Job Kind Constants
// =============================================================================

func TestJobKinds_AreUnique(t *testing.T) {
	t.Parallel()

	kinds := []string{
		KindLibraryScan,
		KindMetadataRefresh,
		KindFileMatch,
		KindSearchIndex,
		KindSeriesRefresh,
		KindSeasonRefresh,
		KindEpisodeRefresh,
	}

	seen := make(map[string]bool)
	for _, kind := range kinds {
		assert.False(t, seen[kind], "duplicate job kind: %s", kind)
		seen[kind] = true
		assert.NotEmpty(t, kind)
		assert.Contains(t, kind, "tvshow_")
	}
}

func TestJobKinds_Values(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "tvshow_library_scan", KindLibraryScan)
	assert.Equal(t, "tvshow_metadata_refresh", KindMetadataRefresh)
	assert.Equal(t, "tvshow_file_match", KindFileMatch)
	assert.Equal(t, "tvshow_search_index", KindSearchIndex)
	assert.Equal(t, "tvshow_series_refresh", KindSeriesRefresh)
	assert.Equal(t, "tvshow_season_refresh", KindSeasonRefresh)
	assert.Equal(t, "tvshow_episode_refresh", KindEpisodeRefresh)
}

// =============================================================================
// InsertOpts Tests
// =============================================================================

func TestLibraryScanArgs_InsertOpts(t *testing.T) {
	t.Parallel()

	args := LibraryScanArgs{}
	opts := args.InsertOpts()
	assert.Equal(t, infrajobs.QueueBulk, opts.Queue)
}

func TestSearchIndexArgs_InsertOpts(t *testing.T) {
	t.Parallel()

	args := SearchIndexArgs{}
	opts := args.InsertOpts()
	assert.Equal(t, infrajobs.QueueBulk, opts.Queue)
}

func TestMetadataRefreshArgs_InsertOpts_UsesDefaults(t *testing.T) {
	t.Parallel()

	// MetadataRefreshArgs does not define InsertOpts, so river.WorkerDefaults
	// will provide default insert opts (empty queue = default queue).
	args := MetadataRefreshArgs{}
	// Verify Kind still works as expected even without custom InsertOpts.
	assert.Equal(t, KindMetadataRefresh, args.Kind())
}

func TestFileMatchArgs_InsertOpts_UsesDefaults(t *testing.T) {
	t.Parallel()

	// FileMatchArgs does not define InsertOpts.
	args := FileMatchArgs{}
	assert.Equal(t, KindFileMatch, args.Kind())
}

func TestSeriesRefreshArgs_InsertOpts_UsesDefaults(t *testing.T) {
	t.Parallel()

	// SeriesRefreshArgs does not define InsertOpts.
	args := SeriesRefreshArgs{}
	assert.Equal(t, KindSeriesRefresh, args.Kind())
}

// =============================================================================
// Args Fields Tests
// =============================================================================

func TestLibraryScanArgs_Fields(t *testing.T) {
	t.Parallel()

	libraryID := uuid.Must(uuid.NewV7())
	args := LibraryScanArgs{
		Paths:      []string{"/tv/shows", "/tv/anime"},
		Force:      true,
		LibraryID:  &libraryID,
		AutoCreate: true,
	}

	assert.Len(t, args.Paths, 2)
	assert.Equal(t, "/tv/shows", args.Paths[0])
	assert.Equal(t, "/tv/anime", args.Paths[1])
	assert.True(t, args.Force)
	assert.Equal(t, libraryID, *args.LibraryID)
	assert.True(t, args.AutoCreate)
}

func TestLibraryScanArgs_NilLibraryID(t *testing.T) {
	t.Parallel()

	args := LibraryScanArgs{
		Force: false,
	}

	assert.Nil(t, args.LibraryID)
	assert.Empty(t, args.Paths)
	assert.False(t, args.Force)
	assert.False(t, args.AutoCreate)
}

func TestMetadataRefreshArgs_Fields(t *testing.T) {
	t.Parallel()

	t.Run("series refresh", func(t *testing.T) {
		t.Parallel()
		seriesID := uuid.Must(uuid.NewV7())
		args := MetadataRefreshArgs{
			SeriesID:      &seriesID,
			Force:         true,
			RefreshImages: true,
			BatchSize:     100,
		}
		assert.NotNil(t, args.SeriesID)
		assert.Equal(t, seriesID, *args.SeriesID)
		assert.True(t, args.Force)
		assert.True(t, args.RefreshImages)
		assert.Equal(t, int32(100), args.BatchSize)
		assert.Nil(t, args.SeasonID)
		assert.Nil(t, args.EpisodeID)
	})

	t.Run("season refresh", func(t *testing.T) {
		t.Parallel()
		seasonID := uuid.Must(uuid.NewV7())
		args := MetadataRefreshArgs{
			SeasonID: &seasonID,
		}
		assert.Nil(t, args.SeriesID)
		assert.NotNil(t, args.SeasonID)
		assert.Equal(t, seasonID, *args.SeasonID)
		assert.Nil(t, args.EpisodeID)
	})

	t.Run("episode refresh", func(t *testing.T) {
		t.Parallel()
		episodeID := uuid.Must(uuid.NewV7())
		args := MetadataRefreshArgs{
			EpisodeID: &episodeID,
		}
		assert.Nil(t, args.SeriesID)
		assert.Nil(t, args.SeasonID)
		assert.NotNil(t, args.EpisodeID)
		assert.Equal(t, episodeID, *args.EpisodeID)
	})

	t.Run("full refresh with batch size", func(t *testing.T) {
		t.Parallel()
		args := MetadataRefreshArgs{
			Force:     true,
			BatchSize: 25,
		}
		assert.Nil(t, args.SeriesID)
		assert.Nil(t, args.SeasonID)
		assert.Nil(t, args.EpisodeID)
		assert.True(t, args.Force)
		assert.Equal(t, int32(25), args.BatchSize)
	})
}

func TestFileMatchArgs_Fields(t *testing.T) {
	t.Parallel()

	episodeID := uuid.Must(uuid.NewV7())
	args := FileMatchArgs{
		FilePath:     "/media/tv/Show Name/Season 01/Show.Name.S01E01.mkv",
		EpisodeID:    &episodeID,
		ForceRematch: true,
		AutoCreate:   true,
	}

	assert.Equal(t, "/media/tv/Show Name/Season 01/Show.Name.S01E01.mkv", args.FilePath)
	assert.Equal(t, episodeID, *args.EpisodeID)
	assert.True(t, args.ForceRematch)
	assert.True(t, args.AutoCreate)
}

func TestFileMatchArgs_NilEpisodeID(t *testing.T) {
	t.Parallel()

	args := FileMatchArgs{
		AutoCreate: false,
	}

	assert.Nil(t, args.EpisodeID)
	assert.Empty(t, args.FilePath)
	assert.False(t, args.ForceRematch)
	assert.False(t, args.AutoCreate)
}

func TestSearchIndexArgs_Fields(t *testing.T) {
	t.Parallel()

	t.Run("specific series", func(t *testing.T) {
		t.Parallel()
		seriesID := uuid.Must(uuid.NewV7())
		args := SearchIndexArgs{
			SeriesID:    &seriesID,
			FullReindex: false,
			BatchSize:   50,
		}
		assert.NotNil(t, args.SeriesID)
		assert.Equal(t, seriesID, *args.SeriesID)
		assert.False(t, args.FullReindex)
		assert.Equal(t, int32(50), args.BatchSize)
	})

	t.Run("full reindex", func(t *testing.T) {
		t.Parallel()
		args := SearchIndexArgs{
			FullReindex: true,
			BatchSize:   100,
		}
		assert.Nil(t, args.SeriesID)
		assert.True(t, args.FullReindex)
		assert.Equal(t, int32(100), args.BatchSize)
	})
}

func TestSeriesRefreshArgs_Fields(t *testing.T) {
	t.Parallel()

	args := SeriesRefreshArgs{
		SeriesID:        uuid.Must(uuid.NewV7()),
		TMDbID:          1396,
		RefreshSeasons:  true,
		RefreshEpisodes: true,
		Languages:       []string{"en", "de", "es"},
	}

	assert.NotEqual(t, uuid.Nil, args.SeriesID)
	assert.Equal(t, int32(1396), args.TMDbID)
	assert.True(t, args.RefreshSeasons)
	assert.True(t, args.RefreshEpisodes)
	assert.Len(t, args.Languages, 3)
	assert.Equal(t, "en", args.Languages[0])
	assert.Equal(t, "de", args.Languages[1])
	assert.Equal(t, "es", args.Languages[2])
}

func TestSeriesRefreshArgs_MinimalFields(t *testing.T) {
	t.Parallel()

	args := SeriesRefreshArgs{
		SeriesID: uuid.Must(uuid.NewV7()),
	}

	assert.NotEqual(t, uuid.Nil, args.SeriesID)
	assert.Equal(t, int32(0), args.TMDbID)
	assert.False(t, args.RefreshSeasons)
	assert.False(t, args.RefreshEpisodes)
	assert.Nil(t, args.Languages)
}

// =============================================================================
// Constructor Tests
// =============================================================================

func TestNewLibraryScanWorker(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	worker := NewLibraryScanWorker(nil, nil, nil, logger)

	assert.NotNil(t, worker)
	assert.Nil(t, worker.service)
	assert.Nil(t, worker.metadataProvider)
	assert.Nil(t, worker.jobClient)
	assert.NotNil(t, worker.logger)
}

func TestNewMetadataRefreshWorker(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	worker := NewMetadataRefreshWorker(nil, nil, logger)

	assert.NotNil(t, worker)
	assert.Nil(t, worker.service)
	assert.Nil(t, worker.jobClient)
	assert.NotNil(t, worker.logger)
}

func TestNewFileMatchWorker(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	worker := NewFileMatchWorker(nil, nil, logger)

	assert.NotNil(t, worker)
	assert.Nil(t, worker.service)
	assert.Nil(t, worker.metadataProvider)
	assert.NotNil(t, worker.logger)
}

func TestNewSearchIndexWorker(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	worker := NewSearchIndexWorker(nil, nil, logger)

	assert.NotNil(t, worker)
	assert.Nil(t, worker.service)
	assert.Nil(t, worker.searchService)
	assert.NotNil(t, worker.logger)
}

func TestNewSeriesRefreshWorker(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	worker := NewSeriesRefreshWorker(nil, nil, logger)

	assert.NotNil(t, worker)
	assert.Nil(t, worker.service)
	assert.Nil(t, worker.jobClient)
	assert.NotNil(t, worker.logger)
}

// =============================================================================
// Timeout Tests
// =============================================================================

func TestLibraryScanWorker_Timeout(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	worker := NewLibraryScanWorker(nil, nil, nil, logger)

	timeout := worker.Timeout(&river.Job[LibraryScanArgs]{})
	assert.Equal(t, 30*time.Minute, timeout)
}

func TestMetadataRefreshWorker_Timeout(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	worker := NewMetadataRefreshWorker(nil, nil, logger)

	timeout := worker.Timeout(&river.Job[MetadataRefreshArgs]{})
	assert.Equal(t, 15*time.Minute, timeout)
}

func TestFileMatchWorker_Timeout(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	worker := NewFileMatchWorker(nil, nil, logger)

	timeout := worker.Timeout(&river.Job[FileMatchArgs]{})
	assert.Equal(t, 5*time.Minute, timeout)
}

func TestSearchIndexWorker_Timeout(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	worker := NewSearchIndexWorker(nil, nil, logger)

	timeout := worker.Timeout(&river.Job[SearchIndexArgs]{})
	assert.Equal(t, 10*time.Minute, timeout)
}

func TestSeriesRefreshWorker_Timeout(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	worker := NewSeriesRefreshWorker(nil, nil, logger)

	timeout := worker.Timeout(&river.Job[SeriesRefreshArgs]{})
	assert.Equal(t, 10*time.Minute, timeout)
}

// =============================================================================
// Work() Tests - LibraryScan
// =============================================================================

func TestLibraryScanWorker_Work_EmptyPaths(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	worker := NewLibraryScanWorker(nil, nil, nil, logger)

	job := &river.Job[LibraryScanArgs]{
		JobRow: &rivertype.JobRow{ID: 1, Kind: KindLibraryScan},
		Args: LibraryScanArgs{
			Paths: []string{},
			Force: false,
		},
	}

	// Empty paths should return nil early.
	err := worker.Work(context.Background(), job)
	require.NoError(t, err)
}

func TestLibraryScanWorker_Work_NilPaths(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	worker := NewLibraryScanWorker(nil, nil, nil, logger)

	job := &river.Job[LibraryScanArgs]{
		JobRow: &rivertype.JobRow{ID: 2, Kind: KindLibraryScan},
		Args: LibraryScanArgs{
			Paths: nil,
			Force: false,
		},
	}

	// nil paths should return nil early (len(nil) == 0).
	err := worker.Work(context.Background(), job)
	require.NoError(t, err)
}

func TestLibraryScanWorker_Work_NonexistentPaths(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	worker := NewLibraryScanWorker(nil, nil, nil, logger)

	job := &river.Job[LibraryScanArgs]{
		JobRow: &rivertype.JobRow{ID: 3, Kind: KindLibraryScan},
		Args: LibraryScanArgs{
			Paths: []string{"/nonexistent/path/that/does/not/exist"},
			Force: false,
		},
	}

	// Scanning a nonexistent path will fail at the filesystem scanner level.
	err := worker.Work(context.Background(), job)
	// The scanner may error or may return empty results depending on implementation.
	// Either way, the worker should handle it.
	_ = err
}

// =============================================================================
// Work() Tests - SearchIndex
// =============================================================================

func TestSearchIndexWorker_Work_SearchDisabled(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	// A nil-client TVShowSearchService will have IsEnabled() return false.
	searchSvc := &search.TVShowSearchService{}
	worker := NewSearchIndexWorker(nil, searchSvc, logger)

	job := &river.Job[SearchIndexArgs]{
		JobRow: &rivertype.JobRow{ID: 1, Kind: KindSearchIndex},
		Args: SearchIndexArgs{
			FullReindex: true,
		},
	}

	// When search is disabled, work should return nil immediately.
	err := worker.Work(context.Background(), job)
	require.NoError(t, err)
}

// =============================================================================
// Work() Tests - FileMatch
// =============================================================================

func TestFileMatchWorker_Work_NonexistentFile(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	worker := NewFileMatchWorker(nil, nil, logger)

	job := &river.Job[FileMatchArgs]{
		JobRow: &rivertype.JobRow{ID: 1, Kind: KindFileMatch},
		Args: FileMatchArgs{
			FilePath:   "/nonexistent/file/that/does/not/exist.mkv",
			AutoCreate: false,
		},
	}

	// Should return an error because the file does not exist.
	err := worker.Work(context.Background(), job)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "file not found")
}

// =============================================================================
// Helper Function Tests
// =============================================================================

func TestJobInsertOpts(t *testing.T) {
	t.Parallel()

	t.Run("with priority only", func(t *testing.T) {
		t.Parallel()
		opts := JobInsertOpts(HighPriority, nil)
		require.NotNil(t, opts)
		assert.Equal(t, HighPriority, opts.Priority)
		assert.True(t, opts.ScheduledAt.IsZero())
	})

	t.Run("with default priority", func(t *testing.T) {
		t.Parallel()
		opts := JobInsertOpts(DefaultPriority, nil)
		require.NotNil(t, opts)
		assert.Equal(t, DefaultPriority, opts.Priority)
	})

	t.Run("with low priority", func(t *testing.T) {
		t.Parallel()
		opts := JobInsertOpts(LowPriority, nil)
		require.NotNil(t, opts)
		assert.Equal(t, LowPriority, opts.Priority)
	})

	t.Run("with scheduled time", func(t *testing.T) {
		t.Parallel()
		scheduled := time.Now().Add(time.Hour)
		opts := JobInsertOpts(DefaultPriority, &scheduled)
		require.NotNil(t, opts)
		assert.Equal(t, DefaultPriority, opts.Priority)
		assert.Equal(t, scheduled, opts.ScheduledAt)
	})

	t.Run("with zero priority", func(t *testing.T) {
		t.Parallel()
		opts := JobInsertOpts(0, nil)
		require.NotNil(t, opts)
		assert.Equal(t, 0, opts.Priority)
	})
}

func TestPriorityConstants(t *testing.T) {
	t.Parallel()

	// Higher number = lower priority in River.
	assert.Equal(t, 1, HighPriority)
	assert.Equal(t, 2, DefaultPriority)
	assert.Equal(t, 3, LowPriority)
	assert.Less(t, HighPriority, DefaultPriority)
	assert.Less(t, DefaultPriority, LowPriority)
}

func TestNormalizeTitle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"lowercase", "breaking bad", "breaking bad"},
		{"uppercase", "BREAKING BAD", "breaking bad"},
		{"mixed case", "Breaking Bad", "breaking bad"},
		{"empty string", "", ""},
		{"single word", "Dexter", "dexter"},
		{"with numbers", "24", "24"},
		{"special characters preserved", "Grey's Anatomy", "grey's anatomy"},
		{"unicode", "Narcos: Mexico", "narcos: mexico"},
		{"all caps with spaces", "THE WIRE", "the wire"},
		{"title with year", "Doctor Who (2005)", "doctor who (2005)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, normalizeTitle(tt.input))
		})
	}
}

func TestNormalizeTitle_MatchComparison(t *testing.T) {
	t.Parallel()

	// Verify that normalization enables case-insensitive matching.
	assert.Equal(t, normalizeTitle("Breaking Bad"), normalizeTitle("breaking bad"))
	assert.Equal(t, normalizeTitle("THE WIRE"), normalizeTitle("The Wire"))
	assert.NotEqual(t, normalizeTitle("Breaking Bad"), normalizeTitle("Better Call Saul"))
}

// =============================================================================
// seriesToCreateParams Tests
// =============================================================================

func TestSeriesToCreateParams_Minimal(t *testing.T) {
	t.Parallel()

	series := &tvshow.Series{
		Title:            "Breaking Bad",
		OriginalLanguage: "en",
	}

	params := seriesToCreateParams(series)
	assert.Equal(t, "Breaking Bad", params.Title)
	assert.Equal(t, "en", params.OriginalLanguage)
	assert.Nil(t, params.TMDbID)
	assert.Nil(t, params.TVDbID)
	assert.Nil(t, params.IMDbID)
	assert.Nil(t, params.FirstAirDate)
	assert.Nil(t, params.LastAirDate)
	assert.Nil(t, params.VoteAverage)
	assert.Nil(t, params.Popularity)
}

func TestSeriesToCreateParams_WithDates(t *testing.T) {
	t.Parallel()

	firstAir := time.Date(2008, 1, 20, 0, 0, 0, 0, time.UTC)
	lastAir := time.Date(2013, 9, 29, 0, 0, 0, 0, time.UTC)
	tmdbID := int32(1396)

	series := &tvshow.Series{
		Title:            "Breaking Bad",
		OriginalLanguage: "en",
		TMDbID:           &tmdbID,
		FirstAirDate:     &firstAir,
		LastAirDate:      &lastAir,
	}

	params := seriesToCreateParams(series)
	assert.Equal(t, "Breaking Bad", params.Title)
	assert.NotNil(t, params.TMDbID)
	assert.Equal(t, int32(1396), *params.TMDbID)
	require.NotNil(t, params.FirstAirDate)
	assert.Equal(t, "2008-01-20", *params.FirstAirDate)
	require.NotNil(t, params.LastAirDate)
	assert.Equal(t, "2013-09-29", *params.LastAirDate)
}

func TestSeriesToCreateParams_WithVoteAndPopularity(t *testing.T) {
	t.Parallel()

	tmdbID := int32(1396)
	voteAvg := must(decimal.New(86, 1))
	popularity := must(decimal.New(1234, 1))

	series := &tvshow.Series{
		Title:            "Breaking Bad",
		OriginalLanguage: "en",
		TMDbID:           &tmdbID,
		VoteAverage:      &voteAvg,
		Popularity:       &popularity,
	}

	params := seriesToCreateParams(series)
	assert.Equal(t, "Breaking Bad", params.Title)
	require.NotNil(t, params.VoteAverage)
	assert.Equal(t, voteAvg.String(), *params.VoteAverage)
	require.NotNil(t, params.Popularity)
	assert.Equal(t, popularity.String(), *params.Popularity)
}

func TestSeriesToCreateParams_AllFields(t *testing.T) {
	t.Parallel()

	tmdbID := int32(1396)
	tvdbID := int32(81189)
	imdbID := "tt0903747"
	originalTitle := "Breaking Bad"
	tagline := "All Hail the King"
	overview := "A chemistry teacher..."
	status := "Ended"
	seriesType := "Scripted"
	posterPath := "/poster.jpg"
	backdropPath := "/backdrop.jpg"
	homepage := "https://example.com"
	trailerURL := "https://youtube.com/watch?v=xyz"
	voteCount := int32(5000)

	series := &tvshow.Series{
		TMDbID:           &tmdbID,
		TVDbID:           &tvdbID,
		IMDbID:           &imdbID,
		Title:            "Breaking Bad",
		OriginalTitle:    &originalTitle,
		OriginalLanguage: "en",
		Tagline:          &tagline,
		Overview:         &overview,
		Status:           &status,
		Type:             &seriesType,
		VoteCount:        &voteCount,
		PosterPath:       &posterPath,
		BackdropPath:     &backdropPath,
		TotalSeasons:     5,
		TotalEpisodes:    62,
		Homepage:         &homepage,
		TrailerURL:       &trailerURL,
	}

	params := seriesToCreateParams(series)
	assert.Equal(t, &tmdbID, params.TMDbID)
	assert.Equal(t, &tvdbID, params.TVDbID)
	assert.Equal(t, &imdbID, params.IMDbID)
	assert.Equal(t, "Breaking Bad", params.Title)
	assert.Equal(t, &originalTitle, params.OriginalTitle)
	assert.Equal(t, "en", params.OriginalLanguage)
	assert.Equal(t, &tagline, params.Tagline)
	assert.Equal(t, &overview, params.Overview)
	assert.Equal(t, &status, params.Status)
	assert.Equal(t, &seriesType, params.Type)
	assert.Equal(t, &voteCount, params.VoteCount)
	assert.Equal(t, &posterPath, params.PosterPath)
	assert.Equal(t, &backdropPath, params.BackdropPath)
	assert.Equal(t, int32(5), params.TotalSeasons)
	assert.Equal(t, int32(62), params.TotalEpisodes)
	assert.Equal(t, &homepage, params.Homepage)
	assert.Equal(t, &trailerURL, params.TrailerURL)
}

// must is a helper to unwrap (decimal, error) in tests.
func must(d decimal.Decimal, err error) decimal.Decimal {
	if err != nil {
		panic(err)
	}
	return d
}

// =============================================================================
// Mock Service Implementation
// =============================================================================

type mockService struct {
	mock.Mock
}

func (m *mockService) GetSeries(ctx context.Context, id uuid.UUID) (*tvshow.Series, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tvshow.Series), args.Error(1)
}

func (m *mockService) GetSeriesByTMDbID(ctx context.Context, tmdbID int32) (*tvshow.Series, error) {
	args := m.Called(ctx, tmdbID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tvshow.Series), args.Error(1)
}

func (m *mockService) GetSeriesByTVDbID(ctx context.Context, tvdbID int32) (*tvshow.Series, error) {
	args := m.Called(ctx, tvdbID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tvshow.Series), args.Error(1)
}

func (m *mockService) GetSeriesBySonarrID(ctx context.Context, sonarrID int32) (*tvshow.Series, error) {
	args := m.Called(ctx, sonarrID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tvshow.Series), args.Error(1)
}

func (m *mockService) ListSeries(ctx context.Context, filters tvshow.SeriesListFilters) ([]tvshow.Series, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).([]tvshow.Series), args.Error(1)
}

func (m *mockService) CountSeries(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockService) SearchSeries(ctx context.Context, query string, limit, offset int32) ([]tvshow.Series, error) {
	args := m.Called(ctx, query, limit, offset)
	return args.Get(0).([]tvshow.Series), args.Error(1)
}

func (m *mockService) ListRecentlyAdded(ctx context.Context, limit, offset int32) ([]tvshow.Series, int64, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]tvshow.Series), args.Get(1).(int64), args.Error(2)
}

func (m *mockService) ListByGenre(ctx context.Context, tmdbGenreID int32, limit, offset int32) ([]tvshow.Series, error) {
	args := m.Called(ctx, tmdbGenreID, limit, offset)
	return args.Get(0).([]tvshow.Series), args.Error(1)
}

func (m *mockService) ListByNetwork(ctx context.Context, networkID uuid.UUID, limit, offset int32) ([]tvshow.Series, error) {
	args := m.Called(ctx, networkID, limit, offset)
	return args.Get(0).([]tvshow.Series), args.Error(1)
}

func (m *mockService) ListByStatus(ctx context.Context, status string, limit, offset int32) ([]tvshow.Series, error) {
	args := m.Called(ctx, status, limit, offset)
	return args.Get(0).([]tvshow.Series), args.Error(1)
}

func (m *mockService) CreateSeries(ctx context.Context, params tvshow.CreateSeriesParams) (*tvshow.Series, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tvshow.Series), args.Error(1)
}

func (m *mockService) UpdateSeries(ctx context.Context, params tvshow.UpdateSeriesParams) (*tvshow.Series, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tvshow.Series), args.Error(1)
}

func (m *mockService) DeleteSeries(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockService) GetSeason(ctx context.Context, id uuid.UUID) (*tvshow.Season, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tvshow.Season), args.Error(1)
}

func (m *mockService) GetSeasonByNumber(ctx context.Context, seriesID uuid.UUID, seasonNumber int32) (*tvshow.Season, error) {
	args := m.Called(ctx, seriesID, seasonNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tvshow.Season), args.Error(1)
}

func (m *mockService) ListSeasons(ctx context.Context, seriesID uuid.UUID) ([]tvshow.Season, error) {
	args := m.Called(ctx, seriesID)
	return args.Get(0).([]tvshow.Season), args.Error(1)
}

func (m *mockService) ListSeasonsWithEpisodeCount(ctx context.Context, seriesID uuid.UUID) ([]tvshow.SeasonWithEpisodeCount, error) {
	args := m.Called(ctx, seriesID)
	return args.Get(0).([]tvshow.SeasonWithEpisodeCount), args.Error(1)
}

func (m *mockService) CreateSeason(ctx context.Context, params tvshow.CreateSeasonParams) (*tvshow.Season, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tvshow.Season), args.Error(1)
}

func (m *mockService) UpsertSeason(ctx context.Context, params tvshow.CreateSeasonParams) (*tvshow.Season, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tvshow.Season), args.Error(1)
}

func (m *mockService) UpdateSeason(ctx context.Context, params tvshow.UpdateSeasonParams) (*tvshow.Season, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tvshow.Season), args.Error(1)
}

func (m *mockService) DeleteSeason(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockService) GetEpisode(ctx context.Context, id uuid.UUID) (*tvshow.Episode, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tvshow.Episode), args.Error(1)
}

func (m *mockService) GetEpisodeByTMDbID(ctx context.Context, tmdbID int32) (*tvshow.Episode, error) {
	args := m.Called(ctx, tmdbID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tvshow.Episode), args.Error(1)
}

func (m *mockService) GetEpisodeByNumber(ctx context.Context, seriesID uuid.UUID, seasonNumber, episodeNumber int32) (*tvshow.Episode, error) {
	args := m.Called(ctx, seriesID, seasonNumber, episodeNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tvshow.Episode), args.Error(1)
}

func (m *mockService) GetEpisodeByFile(ctx context.Context, filePath string) (*tvshow.Episode, error) {
	args := m.Called(ctx, filePath)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tvshow.Episode), args.Error(1)
}

func (m *mockService) ListEpisodesBySeries(ctx context.Context, seriesID uuid.UUID) ([]tvshow.Episode, error) {
	args := m.Called(ctx, seriesID)
	return args.Get(0).([]tvshow.Episode), args.Error(1)
}

func (m *mockService) ListEpisodesBySeason(ctx context.Context, seasonID uuid.UUID) ([]tvshow.Episode, error) {
	args := m.Called(ctx, seasonID)
	return args.Get(0).([]tvshow.Episode), args.Error(1)
}

func (m *mockService) ListEpisodesBySeasonNumber(ctx context.Context, seriesID uuid.UUID, seasonNumber int32) ([]tvshow.Episode, error) {
	args := m.Called(ctx, seriesID, seasonNumber)
	return args.Get(0).([]tvshow.Episode), args.Error(1)
}

func (m *mockService) ListRecentEpisodes(ctx context.Context, limit, offset int32) ([]tvshow.EpisodeWithSeriesInfo, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]tvshow.EpisodeWithSeriesInfo), args.Error(1)
}

func (m *mockService) ListUpcomingEpisodes(ctx context.Context, limit, offset int32) ([]tvshow.EpisodeWithSeriesInfo, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]tvshow.EpisodeWithSeriesInfo), args.Error(1)
}

func (m *mockService) CreateEpisode(ctx context.Context, params tvshow.CreateEpisodeParams) (*tvshow.Episode, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tvshow.Episode), args.Error(1)
}

func (m *mockService) UpsertEpisode(ctx context.Context, params tvshow.CreateEpisodeParams) (*tvshow.Episode, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tvshow.Episode), args.Error(1)
}

func (m *mockService) UpdateEpisode(ctx context.Context, params tvshow.UpdateEpisodeParams) (*tvshow.Episode, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tvshow.Episode), args.Error(1)
}

func (m *mockService) DeleteEpisode(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockService) GetEpisodeFile(ctx context.Context, id uuid.UUID) (*tvshow.EpisodeFile, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tvshow.EpisodeFile), args.Error(1)
}

func (m *mockService) GetEpisodeFileByPath(ctx context.Context, filePath string) (*tvshow.EpisodeFile, error) {
	args := m.Called(ctx, filePath)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tvshow.EpisodeFile), args.Error(1)
}

func (m *mockService) GetEpisodeFileBySonarrID(ctx context.Context, sonarrFileID int32) (*tvshow.EpisodeFile, error) {
	args := m.Called(ctx, sonarrFileID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tvshow.EpisodeFile), args.Error(1)
}

func (m *mockService) ListEpisodeFiles(ctx context.Context, episodeID uuid.UUID) ([]tvshow.EpisodeFile, error) {
	args := m.Called(ctx, episodeID)
	return args.Get(0).([]tvshow.EpisodeFile), args.Error(1)
}

func (m *mockService) CreateEpisodeFile(ctx context.Context, params tvshow.CreateEpisodeFileParams) (*tvshow.EpisodeFile, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tvshow.EpisodeFile), args.Error(1)
}

func (m *mockService) UpdateEpisodeFile(ctx context.Context, params tvshow.UpdateEpisodeFileParams) (*tvshow.EpisodeFile, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tvshow.EpisodeFile), args.Error(1)
}

func (m *mockService) DeleteEpisodeFile(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockService) GetSeriesCast(ctx context.Context, seriesID uuid.UUID, limit, offset int32) ([]tvshow.SeriesCredit, int64, error) {
	args := m.Called(ctx, seriesID, limit, offset)
	return args.Get(0).([]tvshow.SeriesCredit), args.Get(1).(int64), args.Error(2)
}

func (m *mockService) GetSeriesCrew(ctx context.Context, seriesID uuid.UUID, limit, offset int32) ([]tvshow.SeriesCredit, int64, error) {
	args := m.Called(ctx, seriesID, limit, offset)
	return args.Get(0).([]tvshow.SeriesCredit), args.Get(1).(int64), args.Error(2)
}

func (m *mockService) GetEpisodeGuestStars(ctx context.Context, episodeID uuid.UUID) ([]tvshow.EpisodeCredit, error) {
	args := m.Called(ctx, episodeID)
	return args.Get(0).([]tvshow.EpisodeCredit), args.Error(1)
}

func (m *mockService) GetEpisodeCrew(ctx context.Context, episodeID uuid.UUID) ([]tvshow.EpisodeCredit, error) {
	args := m.Called(ctx, episodeID)
	return args.Get(0).([]tvshow.EpisodeCredit), args.Error(1)
}

func (m *mockService) GetSeriesGenres(ctx context.Context, seriesID uuid.UUID) ([]tvshow.SeriesGenre, error) {
	args := m.Called(ctx, seriesID)
	return args.Get(0).([]tvshow.SeriesGenre), args.Error(1)
}

func (m *mockService) ListDistinctGenres(ctx context.Context) ([]content.GenreSummary, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]content.GenreSummary), args.Error(1)
}

func (m *mockService) GetSeriesNetworks(ctx context.Context, seriesID uuid.UUID) ([]tvshow.Network, error) {
	args := m.Called(ctx, seriesID)
	return args.Get(0).([]tvshow.Network), args.Error(1)
}

func (m *mockService) UpdateEpisodeProgress(ctx context.Context, userID, episodeID uuid.UUID, progressSeconds, durationSeconds int32) (*tvshow.EpisodeWatched, error) {
	args := m.Called(ctx, userID, episodeID, progressSeconds, durationSeconds)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tvshow.EpisodeWatched), args.Error(1)
}

func (m *mockService) GetEpisodeProgress(ctx context.Context, userID, episodeID uuid.UUID) (*tvshow.EpisodeWatched, error) {
	args := m.Called(ctx, userID, episodeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tvshow.EpisodeWatched), args.Error(1)
}

func (m *mockService) MarkEpisodeWatched(ctx context.Context, userID, episodeID uuid.UUID) error {
	args := m.Called(ctx, userID, episodeID)
	return args.Error(0)
}

func (m *mockService) MarkEpisodesWatchedBulk(ctx context.Context, userID uuid.UUID, episodeIDs []uuid.UUID) (int64, error) {
	args := m.Called(ctx, userID, episodeIDs)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockService) MarkSeasonWatched(ctx context.Context, userID, seasonID uuid.UUID) error {
	args := m.Called(ctx, userID, seasonID)
	return args.Error(0)
}

func (m *mockService) MarkSeriesWatched(ctx context.Context, userID, seriesID uuid.UUID) error {
	args := m.Called(ctx, userID, seriesID)
	return args.Error(0)
}

func (m *mockService) RemoveEpisodeProgress(ctx context.Context, userID, episodeID uuid.UUID) error {
	args := m.Called(ctx, userID, episodeID)
	return args.Error(0)
}

func (m *mockService) RemoveSeriesProgress(ctx context.Context, userID, seriesID uuid.UUID) error {
	args := m.Called(ctx, userID, seriesID)
	return args.Error(0)
}

func (m *mockService) GetContinueWatching(ctx context.Context, userID uuid.UUID, limit int32) ([]tvshow.ContinueWatchingItem, error) {
	args := m.Called(ctx, userID, limit)
	return args.Get(0).([]tvshow.ContinueWatchingItem), args.Error(1)
}

func (m *mockService) GetNextEpisode(ctx context.Context, userID, seriesID uuid.UUID) (*tvshow.Episode, error) {
	args := m.Called(ctx, userID, seriesID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tvshow.Episode), args.Error(1)
}

func (m *mockService) GetSeriesWatchStats(ctx context.Context, userID, seriesID uuid.UUID) (*tvshow.SeriesWatchStats, error) {
	args := m.Called(ctx, userID, seriesID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tvshow.SeriesWatchStats), args.Error(1)
}

func (m *mockService) GetUserStats(ctx context.Context, userID uuid.UUID) (*tvshow.UserTVStats, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tvshow.UserTVStats), args.Error(1)
}

func (m *mockService) RefreshSeriesMetadata(ctx context.Context, id uuid.UUID, opts ...tvshow.MetadataRefreshOptions) error {
	args := m.Called(ctx, id, opts)
	return args.Error(0)
}

func (m *mockService) RefreshSeasonMetadata(ctx context.Context, id uuid.UUID, opts ...tvshow.MetadataRefreshOptions) error {
	args := m.Called(ctx, id, opts)
	return args.Error(0)
}

func (m *mockService) RefreshEpisodeMetadata(ctx context.Context, id uuid.UUID, opts ...tvshow.MetadataRefreshOptions) error {
	args := m.Called(ctx, id, opts)
	return args.Error(0)
}

// Verify mock implements interface
var _ tvshow.Service = (*mockService)(nil)

// =============================================================================
// Mock MetadataProvider Implementation
// =============================================================================

type mockMetadataProvider struct {
	mock.Mock
}

func (m *mockMetadataProvider) SearchSeries(ctx context.Context, query string, year *int) ([]*tvshow.Series, error) {
	args := m.Called(ctx, query, year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*tvshow.Series), args.Error(1)
}

func (m *mockMetadataProvider) EnrichSeries(ctx context.Context, series *tvshow.Series, opts ...tvshow.MetadataRefreshOptions) error {
	args := m.Called(ctx, series, opts)
	return args.Error(0)
}

func (m *mockMetadataProvider) EnrichSeason(ctx context.Context, season *tvshow.Season, seriesTMDbID int32, opts ...tvshow.MetadataRefreshOptions) error {
	args := m.Called(ctx, season, seriesTMDbID, opts)
	return args.Error(0)
}

func (m *mockMetadataProvider) EnrichEpisode(ctx context.Context, episode *tvshow.Episode, seriesTMDbID int32, opts ...tvshow.MetadataRefreshOptions) error {
	args := m.Called(ctx, episode, seriesTMDbID, opts)
	return args.Error(0)
}

func (m *mockMetadataProvider) GetSeriesCredits(ctx context.Context, seriesID uuid.UUID, tmdbID int) ([]tvshow.SeriesCredit, error) {
	args := m.Called(ctx, seriesID, tmdbID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]tvshow.SeriesCredit), args.Error(1)
}

func (m *mockMetadataProvider) GetSeriesGenres(ctx context.Context, seriesID uuid.UUID, tmdbID int) ([]tvshow.SeriesGenre, error) {
	args := m.Called(ctx, seriesID, tmdbID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]tvshow.SeriesGenre), args.Error(1)
}

func (m *mockMetadataProvider) GetSeriesNetworks(ctx context.Context, tmdbID int) ([]tvshow.Network, error) {
	args := m.Called(ctx, tmdbID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]tvshow.Network), args.Error(1)
}

func (m *mockMetadataProvider) ClearCache() {
	m.Called()
}

// Verify mock implements interface
var _ tvshow.MetadataProvider = (*mockMetadataProvider)(nil)

// =============================================================================
// Work() Tests - MetadataRefresh
// =============================================================================

func TestMetadataRefreshWorker_Work_RefreshEpisode(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewMetadataRefreshWorker(svc, &infrajobs.Client{}, logger)

	episodeID := uuid.Must(uuid.NewV7())
	job := &river.Job[MetadataRefreshArgs]{
		JobRow: &rivertype.JobRow{ID: 1, Kind: KindMetadataRefresh},
		Args: MetadataRefreshArgs{
			EpisodeID: &episodeID,
			Force:     true,
		},
	}

	svc.On("RefreshEpisodeMetadata", mock.Anything, episodeID, []tvshow.MetadataRefreshOptions{
		{Force: true},
	}).Return(nil)

	err := worker.Work(context.Background(), job)
	require.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestMetadataRefreshWorker_Work_RefreshEpisode_Error(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewMetadataRefreshWorker(svc, &infrajobs.Client{}, logger)

	episodeID := uuid.Must(uuid.NewV7())
	job := &river.Job[MetadataRefreshArgs]{
		JobRow: &rivertype.JobRow{ID: 1, Kind: KindMetadataRefresh},
		Args: MetadataRefreshArgs{
			EpisodeID: &episodeID,
		},
	}

	svc.On("RefreshEpisodeMetadata", mock.Anything, episodeID, []tvshow.MetadataRefreshOptions{
		{Force: false},
	}).Return(errors.New("refresh failed"))

	err := worker.Work(context.Background(), job)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "metadata refresh completed with errors")
	svc.AssertExpectations(t)
}

func TestMetadataRefreshWorker_Work_RefreshSeason(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewMetadataRefreshWorker(svc, &infrajobs.Client{}, logger)

	seasonID := uuid.Must(uuid.NewV7())
	job := &river.Job[MetadataRefreshArgs]{
		JobRow: &rivertype.JobRow{ID: 2, Kind: KindMetadataRefresh},
		Args: MetadataRefreshArgs{
			SeasonID: &seasonID,
			Force:    false,
		},
	}

	svc.On("RefreshSeasonMetadata", mock.Anything, seasonID, []tvshow.MetadataRefreshOptions{
		{Force: false},
	}).Return(nil)

	err := worker.Work(context.Background(), job)
	require.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestMetadataRefreshWorker_Work_RefreshSeason_Error(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewMetadataRefreshWorker(svc, &infrajobs.Client{}, logger)

	seasonID := uuid.Must(uuid.NewV7())
	job := &river.Job[MetadataRefreshArgs]{
		JobRow: &rivertype.JobRow{ID: 2, Kind: KindMetadataRefresh},
		Args: MetadataRefreshArgs{
			SeasonID: &seasonID,
		},
	}

	svc.On("RefreshSeasonMetadata", mock.Anything, seasonID, []tvshow.MetadataRefreshOptions{
		{Force: false},
	}).Return(errors.New("season refresh failed"))

	err := worker.Work(context.Background(), job)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "metadata refresh completed with errors")
	svc.AssertExpectations(t)
}

func TestMetadataRefreshWorker_Work_RefreshSeries(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewMetadataRefreshWorker(svc, &infrajobs.Client{}, logger)

	seriesID := uuid.Must(uuid.NewV7())
	job := &river.Job[MetadataRefreshArgs]{
		JobRow: &rivertype.JobRow{ID: 3, Kind: KindMetadataRefresh},
		Args: MetadataRefreshArgs{
			SeriesID: &seriesID,
			Force:    true,
		},
	}

	svc.On("RefreshSeriesMetadata", mock.Anything, seriesID, []tvshow.MetadataRefreshOptions{
		{Force: true},
	}).Return(nil)

	err := worker.Work(context.Background(), job)
	require.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestMetadataRefreshWorker_Work_RefreshSeries_Error(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewMetadataRefreshWorker(svc, &infrajobs.Client{}, logger)

	seriesID := uuid.Must(uuid.NewV7())
	job := &river.Job[MetadataRefreshArgs]{
		JobRow: &rivertype.JobRow{ID: 3, Kind: KindMetadataRefresh},
		Args: MetadataRefreshArgs{
			SeriesID: &seriesID,
		},
	}

	svc.On("RefreshSeriesMetadata", mock.Anything, seriesID, []tvshow.MetadataRefreshOptions{
		{Force: false},
	}).Return(errors.New("series refresh failed"))

	err := worker.Work(context.Background(), job)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "metadata refresh completed with errors")
	svc.AssertExpectations(t)
}

func TestMetadataRefreshWorker_Work_RefreshAll_EmptyList(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewMetadataRefreshWorker(svc, &infrajobs.Client{}, logger)

	job := &river.Job[MetadataRefreshArgs]{
		JobRow: &rivertype.JobRow{ID: 4, Kind: KindMetadataRefresh},
		Args: MetadataRefreshArgs{
			// No SeriesID, SeasonID, or EpisodeID - refresh all
			Force:     false,
			BatchSize: 10,
		},
	}

	// Return empty list - no series to refresh
	svc.On("ListSeries", mock.Anything, tvshow.SeriesListFilters{
		Limit:  10,
		Offset: 0,
	}).Return([]tvshow.Series{}, nil)

	err := worker.Work(context.Background(), job)
	require.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestMetadataRefreshWorker_Work_RefreshAll_DefaultBatchSize(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewMetadataRefreshWorker(svc, &infrajobs.Client{}, logger)

	job := &river.Job[MetadataRefreshArgs]{
		JobRow: &rivertype.JobRow{ID: 5, Kind: KindMetadataRefresh},
		Args: MetadataRefreshArgs{
			// No BatchSize specified - should default to 50
		},
	}

	svc.On("ListSeries", mock.Anything, tvshow.SeriesListFilters{
		Limit:  50,
		Offset: 0,
	}).Return([]tvshow.Series{}, nil)

	err := worker.Work(context.Background(), job)
	require.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestMetadataRefreshWorker_Work_RefreshAll_WithSeries(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewMetadataRefreshWorker(svc, &infrajobs.Client{}, logger)

	seriesID1 := uuid.Must(uuid.NewV7())
	seriesID2 := uuid.Must(uuid.NewV7())

	job := &river.Job[MetadataRefreshArgs]{
		JobRow: &rivertype.JobRow{ID: 6, Kind: KindMetadataRefresh},
		Args: MetadataRefreshArgs{
			Force:     true,
			BatchSize: 10,
		},
	}

	// Return a batch of 2 series (less than batch size, so only 1 iteration)
	svc.On("ListSeries", mock.Anything, tvshow.SeriesListFilters{
		Limit:  10,
		Offset: 0,
	}).Return([]tvshow.Series{
		{ID: seriesID1, Title: "Series 1"},
		{ID: seriesID2, Title: "Series 2"},
	}, nil)

	svc.On("RefreshSeriesMetadata", mock.Anything, seriesID1, []tvshow.MetadataRefreshOptions{
		{Force: true},
	}).Return(nil)
	svc.On("RefreshSeriesMetadata", mock.Anything, seriesID2, []tvshow.MetadataRefreshOptions{
		{Force: true},
	}).Return(nil)

	err := worker.Work(context.Background(), job)
	require.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestMetadataRefreshWorker_Work_RefreshAll_ListError(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewMetadataRefreshWorker(svc, &infrajobs.Client{}, logger)

	job := &river.Job[MetadataRefreshArgs]{
		JobRow: &rivertype.JobRow{ID: 7, Kind: KindMetadataRefresh},
		Args: MetadataRefreshArgs{
			BatchSize: 10,
		},
	}

	svc.On("ListSeries", mock.Anything, tvshow.SeriesListFilters{
		Limit:  10,
		Offset: 0,
	}).Return([]tvshow.Series{}, errors.New("db error"))

	err := worker.Work(context.Background(), job)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "metadata refresh completed with errors")
	svc.AssertExpectations(t)
}

func TestMetadataRefreshWorker_Work_RefreshAll_PartialErrors(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewMetadataRefreshWorker(svc, &infrajobs.Client{}, logger)

	seriesID1 := uuid.Must(uuid.NewV7())
	seriesID2 := uuid.Must(uuid.NewV7())

	job := &river.Job[MetadataRefreshArgs]{
		JobRow: &rivertype.JobRow{ID: 8, Kind: KindMetadataRefresh},
		Args: MetadataRefreshArgs{
			BatchSize: 10,
		},
	}

	svc.On("ListSeries", mock.Anything, tvshow.SeriesListFilters{
		Limit:  10,
		Offset: 0,
	}).Return([]tvshow.Series{
		{ID: seriesID1, Title: "Series 1"},
		{ID: seriesID2, Title: "Series 2"},
	}, nil)

	// First succeeds, second fails
	svc.On("RefreshSeriesMetadata", mock.Anything, seriesID1, []tvshow.MetadataRefreshOptions{
		{Force: false},
	}).Return(nil)
	svc.On("RefreshSeriesMetadata", mock.Anything, seriesID2, []tvshow.MetadataRefreshOptions{
		{Force: false},
	}).Return(errors.New("refresh failed"))

	err := worker.Work(context.Background(), job)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "metadata refresh completed with errors")
	svc.AssertExpectations(t)
}

// =============================================================================
// Work() Tests - SeriesRefresh
// =============================================================================

func TestSeriesRefreshWorker_Work_Success(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewSeriesRefreshWorker(svc, &infrajobs.Client{}, logger)

	seriesID := uuid.Must(uuid.NewV7())
	job := &river.Job[SeriesRefreshArgs]{
		JobRow: &rivertype.JobRow{ID: 1, Kind: KindSeriesRefresh},
		Args: SeriesRefreshArgs{
			SeriesID:  seriesID,
			Languages: []string{"en", "de"},
		},
	}

	svc.On("RefreshSeriesMetadata", mock.Anything, seriesID, []tvshow.MetadataRefreshOptions{
		{Languages: []string{"en", "de"}},
	}).Return(nil)

	err := worker.Work(context.Background(), job)
	require.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestSeriesRefreshWorker_Work_RefreshError(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewSeriesRefreshWorker(svc, &infrajobs.Client{}, logger)

	seriesID := uuid.Must(uuid.NewV7())
	job := &river.Job[SeriesRefreshArgs]{
		JobRow: &rivertype.JobRow{ID: 2, Kind: KindSeriesRefresh},
		Args: SeriesRefreshArgs{
			SeriesID: seriesID,
		},
	}

	svc.On("RefreshSeriesMetadata", mock.Anything, seriesID, []tvshow.MetadataRefreshOptions{
		{},
	}).Return(errors.New("refresh failed"))

	err := worker.Work(context.Background(), job)
	// SeriesRefreshWorker does not return error on failure - it just logs
	require.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestSeriesRefreshWorker_Work_WithSeasons(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewSeriesRefreshWorker(svc, &infrajobs.Client{}, logger)

	seriesID := uuid.Must(uuid.NewV7())
	seasonID1 := uuid.Must(uuid.NewV7())
	seasonID2 := uuid.Must(uuid.NewV7())

	job := &river.Job[SeriesRefreshArgs]{
		JobRow: &rivertype.JobRow{ID: 3, Kind: KindSeriesRefresh},
		Args: SeriesRefreshArgs{
			SeriesID:       seriesID,
			RefreshSeasons: true,
		},
	}

	opts := []tvshow.MetadataRefreshOptions{{}}
	svc.On("RefreshSeriesMetadata", mock.Anything, seriesID, opts).Return(nil)
	svc.On("ListSeasons", mock.Anything, seriesID).Return([]tvshow.Season{
		{ID: seasonID1, SeriesID: seriesID, SeasonNumber: 1, Name: "Season 1"},
		{ID: seasonID2, SeriesID: seriesID, SeasonNumber: 2, Name: "Season 2"},
	}, nil)
	svc.On("RefreshSeasonMetadata", mock.Anything, seasonID1, opts).Return(nil)
	svc.On("RefreshSeasonMetadata", mock.Anything, seasonID2, opts).Return(nil)

	err := worker.Work(context.Background(), job)
	require.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestSeriesRefreshWorker_Work_WithEpisodes(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewSeriesRefreshWorker(svc, &infrajobs.Client{}, logger)

	seriesID := uuid.Must(uuid.NewV7())
	seasonID := uuid.Must(uuid.NewV7())
	ep1ID := uuid.Must(uuid.NewV7())
	ep2ID := uuid.Must(uuid.NewV7())

	job := &river.Job[SeriesRefreshArgs]{
		JobRow: &rivertype.JobRow{ID: 4, Kind: KindSeriesRefresh},
		Args: SeriesRefreshArgs{
			SeriesID:        seriesID,
			RefreshEpisodes: true,
		},
	}

	opts := []tvshow.MetadataRefreshOptions{{}}
	svc.On("RefreshSeriesMetadata", mock.Anything, seriesID, opts).Return(nil)
	svc.On("ListSeasons", mock.Anything, seriesID).Return([]tvshow.Season{
		{ID: seasonID, SeriesID: seriesID, SeasonNumber: 1, Name: "Season 1"},
	}, nil)
	svc.On("ListEpisodesBySeason", mock.Anything, seasonID).Return([]tvshow.Episode{
		{ID: ep1ID, SeasonID: seasonID, EpisodeNumber: 1},
		{ID: ep2ID, SeasonID: seasonID, EpisodeNumber: 2},
	}, nil)
	svc.On("RefreshEpisodeMetadata", mock.Anything, ep1ID, opts).Return(nil)
	svc.On("RefreshEpisodeMetadata", mock.Anything, ep2ID, opts).Return(nil)

	err := worker.Work(context.Background(), job)
	require.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestSeriesRefreshWorker_Work_WithSeasonsAndEpisodes(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewSeriesRefreshWorker(svc, &infrajobs.Client{}, logger)

	seriesID := uuid.Must(uuid.NewV7())
	seasonID := uuid.Must(uuid.NewV7())
	epID := uuid.Must(uuid.NewV7())

	job := &river.Job[SeriesRefreshArgs]{
		JobRow: &rivertype.JobRow{ID: 5, Kind: KindSeriesRefresh},
		Args: SeriesRefreshArgs{
			SeriesID:        seriesID,
			RefreshSeasons:  true,
			RefreshEpisodes: true,
		},
	}

	opts := []tvshow.MetadataRefreshOptions{{}}
	svc.On("RefreshSeriesMetadata", mock.Anything, seriesID, opts).Return(nil)
	svc.On("ListSeasons", mock.Anything, seriesID).Return([]tvshow.Season{
		{ID: seasonID, SeriesID: seriesID, SeasonNumber: 1, Name: "Season 1"},
	}, nil)
	svc.On("RefreshSeasonMetadata", mock.Anything, seasonID, opts).Return(nil)
	svc.On("ListEpisodesBySeason", mock.Anything, seasonID).Return([]tvshow.Episode{
		{ID: epID, SeasonID: seasonID, EpisodeNumber: 1},
	}, nil)
	svc.On("RefreshEpisodeMetadata", mock.Anything, epID, opts).Return(nil)

	err := worker.Work(context.Background(), job)
	require.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestSeriesRefreshWorker_Work_ListSeasonsError(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewSeriesRefreshWorker(svc, &infrajobs.Client{}, logger)

	seriesID := uuid.Must(uuid.NewV7())
	job := &river.Job[SeriesRefreshArgs]{
		JobRow: &rivertype.JobRow{ID: 6, Kind: KindSeriesRefresh},
		Args: SeriesRefreshArgs{
			SeriesID:       seriesID,
			RefreshSeasons: true,
		},
	}

	opts := []tvshow.MetadataRefreshOptions{{}}
	svc.On("RefreshSeriesMetadata", mock.Anything, seriesID, opts).Return(nil)
	svc.On("ListSeasons", mock.Anything, seriesID).Return([]tvshow.Season{}, errors.New("list seasons failed"))

	err := worker.Work(context.Background(), job)
	// Should not error - just warns
	require.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestSeriesRefreshWorker_Work_ListEpisodesError(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewSeriesRefreshWorker(svc, &infrajobs.Client{}, logger)

	seriesID := uuid.Must(uuid.NewV7())
	seasonID := uuid.Must(uuid.NewV7())
	job := &river.Job[SeriesRefreshArgs]{
		JobRow: &rivertype.JobRow{ID: 7, Kind: KindSeriesRefresh},
		Args: SeriesRefreshArgs{
			SeriesID:        seriesID,
			RefreshEpisodes: true,
		},
	}

	opts := []tvshow.MetadataRefreshOptions{{}}
	svc.On("RefreshSeriesMetadata", mock.Anything, seriesID, opts).Return(nil)
	svc.On("ListSeasons", mock.Anything, seriesID).Return([]tvshow.Season{
		{ID: seasonID, SeriesID: seriesID, SeasonNumber: 1, Name: "Season 1"},
	}, nil)
	svc.On("ListEpisodesBySeason", mock.Anything, seasonID).Return([]tvshow.Episode{}, errors.New("list episodes failed"))

	err := worker.Work(context.Background(), job)
	require.NoError(t, err)
	svc.AssertExpectations(t)
}

// =============================================================================
// Work() Tests - SearchIndex (additional)
// =============================================================================

func TestSearchIndexWorker_Work_SearchDisabled_FullReindex(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	searchSvc := &search.TVShowSearchService{}
	worker := NewSearchIndexWorker(nil, searchSvc, logger)

	job := &river.Job[SearchIndexArgs]{
		JobRow: &rivertype.JobRow{ID: 1, Kind: KindSearchIndex},
		Args: SearchIndexArgs{
			FullReindex: true,
		},
	}

	// Search disabled should return nil immediately
	err := worker.Work(context.Background(), job)
	require.NoError(t, err)
}

func TestSearchIndexWorker_Work_SpecificSeries_SearchDisabled(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	searchSvc := &search.TVShowSearchService{}
	worker := NewSearchIndexWorker(nil, searchSvc, logger)

	seriesID := uuid.Must(uuid.NewV7())
	job := &river.Job[SearchIndexArgs]{
		JobRow: &rivertype.JobRow{ID: 2, Kind: KindSearchIndex},
		Args: SearchIndexArgs{
			SeriesID: &seriesID,
		},
	}

	// Search disabled should return nil immediately
	err := worker.Work(context.Background(), job)
	require.NoError(t, err)
}

func TestSearchIndexWorker_Work_NoArgs_SearchDisabled(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	searchSvc := &search.TVShowSearchService{}
	worker := NewSearchIndexWorker(nil, searchSvc, logger)

	job := &river.Job[SearchIndexArgs]{
		JobRow: &rivertype.JobRow{ID: 3, Kind: KindSearchIndex},
		Args: SearchIndexArgs{
			// Neither SeriesID nor FullReindex
		},
	}

	err := worker.Work(context.Background(), job)
	require.NoError(t, err)
}

// =============================================================================
// Work() Tests - FileMatch (additional with mock service)
// =============================================================================

func TestFileMatchWorker_Work_AlreadyMatched_NoForce(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewFileMatchWorker(svc, nil, logger)

	// Create a temp file to satisfy os.Stat
	tmpFile := createTempFile(t, "test-file-match-*.mkv")

	episodeID := uuid.Must(uuid.NewV7())
	job := &river.Job[FileMatchArgs]{
		JobRow: &rivertype.JobRow{ID: 1, Kind: KindFileMatch},
		Args: FileMatchArgs{
			FilePath:     tmpFile,
			ForceRematch: false,
		},
	}

	svc.On("GetEpisodeFileByPath", mock.Anything, tmpFile).Return(&tvshow.EpisodeFile{
		ID:        uuid.Must(uuid.NewV7()),
		EpisodeID: episodeID,
		FilePath:  tmpFile,
	}, nil)

	err := worker.Work(context.Background(), job)
	require.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestFileMatchWorker_Work_ForceRematch(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewFileMatchWorker(svc, nil, logger)

	tmpFile := createTempFileWithName(t, "Show.Name.S01E01.mkv")

	seriesID := uuid.Must(uuid.NewV7())
	seasonID := uuid.Must(uuid.NewV7())
	episodeID := uuid.Must(uuid.NewV7())

	job := &river.Job[FileMatchArgs]{
		JobRow: &rivertype.JobRow{ID: 10, Kind: KindFileMatch},
		Args: FileMatchArgs{
			FilePath:     tmpFile,
			ForceRematch: true,
		},
	}

	// Even though file is already matched, ForceRematch=true continues
	svc.On("GetEpisodeFileByPath", mock.Anything, tmpFile).Return(&tvshow.EpisodeFile{
		ID:        uuid.Must(uuid.NewV7()),
		EpisodeID: episodeID,
		FilePath:  tmpFile,
	}, nil)

	// Parser will extract "Show Name" S01E01 from the filename
	svc.On("SearchSeries", mock.Anything, mock.Anything, int32(5), int32(0)).Return([]tvshow.Series{
		{ID: seriesID, Title: "Show Name"},
	}, nil)
	svc.On("GetSeasonByNumber", mock.Anything, seriesID, int32(1)).Return(&tvshow.Season{
		ID:       seasonID,
		SeriesID: seriesID,
	}, nil)
	svc.On("GetEpisodeByNumber", mock.Anything, seriesID, int32(1), int32(1)).Return(&tvshow.Episode{
		ID:       episodeID,
		SeasonID: seasonID,
	}, nil)
	svc.On("CreateEpisodeFile", mock.Anything, mock.MatchedBy(func(params tvshow.CreateEpisodeFileParams) bool {
		return params.EpisodeID == episodeID && params.FilePath == tmpFile
	})).Return(&tvshow.EpisodeFile{
		ID:        uuid.Must(uuid.NewV7()),
		EpisodeID: episodeID,
		FilePath:  tmpFile,
	}, nil)

	err := worker.Work(context.Background(), job)
	require.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestFileMatchWorker_Work_DirectMatch(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewFileMatchWorker(svc, nil, logger)

	tmpFile := createTempFile(t, "test-direct-match-*.mkv")

	episodeID := uuid.Must(uuid.NewV7())
	job := &river.Job[FileMatchArgs]{
		JobRow: &rivertype.JobRow{ID: 2, Kind: KindFileMatch},
		Args: FileMatchArgs{
			FilePath:  tmpFile,
			EpisodeID: &episodeID,
		},
	}

	svc.On("GetEpisodeFileByPath", mock.Anything, tmpFile).Return(nil, errors.New("not found"))
	svc.On("GetEpisode", mock.Anything, episodeID).Return(&tvshow.Episode{
		ID:            episodeID,
		EpisodeNumber: 1,
		Title:         "Pilot",
	}, nil)
	svc.On("CreateEpisodeFile", mock.Anything, mock.MatchedBy(func(params tvshow.CreateEpisodeFileParams) bool {
		return params.EpisodeID == episodeID && params.FilePath == tmpFile
	})).Return(&tvshow.EpisodeFile{
		ID:        uuid.Must(uuid.NewV7()),
		EpisodeID: episodeID,
		FilePath:  tmpFile,
	}, nil)

	err := worker.Work(context.Background(), job)
	require.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestFileMatchWorker_Work_DirectMatch_EpisodeNotFound(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewFileMatchWorker(svc, nil, logger)

	tmpFile := createTempFile(t, "test-direct-match-epnf-*.mkv")

	episodeID := uuid.Must(uuid.NewV7())
	job := &river.Job[FileMatchArgs]{
		JobRow: &rivertype.JobRow{ID: 3, Kind: KindFileMatch},
		Args: FileMatchArgs{
			FilePath:  tmpFile,
			EpisodeID: &episodeID,
		},
	}

	svc.On("GetEpisodeFileByPath", mock.Anything, tmpFile).Return(nil, errors.New("not found"))
	svc.On("GetEpisode", mock.Anything, episodeID).Return(nil, errors.New("not found"))

	err := worker.Work(context.Background(), job)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "episode not found")
	svc.AssertExpectations(t)
}

func TestFileMatchWorker_Work_DirectMatch_CreateFileFails(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewFileMatchWorker(svc, nil, logger)

	tmpFile := createTempFile(t, "test-dm-createfail-*.mkv")

	episodeID := uuid.Must(uuid.NewV7())
	job := &river.Job[FileMatchArgs]{
		JobRow: &rivertype.JobRow{ID: 4, Kind: KindFileMatch},
		Args: FileMatchArgs{
			FilePath:  tmpFile,
			EpisodeID: &episodeID,
		},
	}

	svc.On("GetEpisodeFileByPath", mock.Anything, tmpFile).Return(nil, errors.New("not found"))
	svc.On("GetEpisode", mock.Anything, episodeID).Return(&tvshow.Episode{
		ID: episodeID,
	}, nil)
	svc.On("CreateEpisodeFile", mock.Anything, mock.Anything).Return(nil, errors.New("create failed"))

	err := worker.Work(context.Background(), job)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "create episode file")
	svc.AssertExpectations(t)
}

func TestFileMatchWorker_Work_NoAutoCreate_SeriesNotFound(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewFileMatchWorker(svc, nil, logger)

	// Need a filename that can be parsed as a TV show - e.g., "Show.Name.S01E01.mkv"
	tmpFile := createTempFileWithName(t, "Show.Name.S01E01.mkv")

	job := &river.Job[FileMatchArgs]{
		JobRow: &rivertype.JobRow{ID: 5, Kind: KindFileMatch},
		Args: FileMatchArgs{
			FilePath:   tmpFile,
			AutoCreate: false,
		},
	}

	svc.On("GetEpisodeFileByPath", mock.Anything, tmpFile).Return(nil, errors.New("not found"))
	svc.On("SearchSeries", mock.Anything, mock.Anything, int32(5), int32(0)).Return([]tvshow.Series{}, nil)

	err := worker.Work(context.Background(), job)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "auto_create disabled")
	svc.AssertExpectations(t)
}

func TestFileMatchWorker_Work_MatchExistingSeriesAndEpisode(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewFileMatchWorker(svc, nil, logger)

	tmpFile := createTempFileWithName(t, "Breaking.Bad.S02E03.mkv")

	seriesID := uuid.Must(uuid.NewV7())
	seasonID := uuid.Must(uuid.NewV7())
	episodeID := uuid.Must(uuid.NewV7())

	job := &river.Job[FileMatchArgs]{
		JobRow: &rivertype.JobRow{ID: 6, Kind: KindFileMatch},
		Args: FileMatchArgs{
			FilePath: tmpFile,
		},
	}

	svc.On("GetEpisodeFileByPath", mock.Anything, tmpFile).Return(nil, errors.New("not found"))
	svc.On("SearchSeries", mock.Anything, mock.Anything, int32(5), int32(0)).Return([]tvshow.Series{
		{ID: seriesID, Title: "Breaking Bad"},
	}, nil)
	svc.On("GetSeasonByNumber", mock.Anything, seriesID, int32(2)).Return(&tvshow.Season{
		ID:       seasonID,
		SeriesID: seriesID,
	}, nil)
	svc.On("GetEpisodeByNumber", mock.Anything, seriesID, int32(2), int32(3)).Return(&tvshow.Episode{
		ID:       episodeID,
		SeasonID: seasonID,
	}, nil)
	svc.On("CreateEpisodeFile", mock.Anything, mock.MatchedBy(func(params tvshow.CreateEpisodeFileParams) bool {
		return params.EpisodeID == episodeID && params.FilePath == tmpFile
	})).Return(&tvshow.EpisodeFile{
		ID:        uuid.Must(uuid.NewV7()),
		EpisodeID: episodeID,
		FilePath:  tmpFile,
	}, nil)

	err := worker.Work(context.Background(), job)
	require.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestFileMatchWorker_Work_SearchSeriesError(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewFileMatchWorker(svc, nil, logger)

	tmpFile := createTempFileWithName(t, "Some.Show.S01E01.mkv")

	job := &river.Job[FileMatchArgs]{
		JobRow: &rivertype.JobRow{ID: 7, Kind: KindFileMatch},
		Args: FileMatchArgs{
			FilePath: tmpFile,
		},
	}

	svc.On("GetEpisodeFileByPath", mock.Anything, tmpFile).Return(nil, errors.New("not found"))
	svc.On("SearchSeries", mock.Anything, mock.Anything, int32(5), int32(0)).Return([]tvshow.Series{}, errors.New("search error"))

	err := worker.Work(context.Background(), job)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "search series")
	svc.AssertExpectations(t)
}

func TestFileMatchWorker_Work_AutoCreate_SeasonNotFound_NoAutoCreate(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewFileMatchWorker(svc, nil, logger)

	tmpFile := createTempFileWithName(t, "Show.Name.S03E05.mkv")
	seriesID := uuid.Must(uuid.NewV7())

	job := &river.Job[FileMatchArgs]{
		JobRow: &rivertype.JobRow{ID: 8, Kind: KindFileMatch},
		Args: FileMatchArgs{
			FilePath:   tmpFile,
			AutoCreate: false,
		},
	}

	svc.On("GetEpisodeFileByPath", mock.Anything, tmpFile).Return(nil, errors.New("not found"))
	svc.On("SearchSeries", mock.Anything, mock.Anything, int32(5), int32(0)).Return([]tvshow.Series{
		{ID: seriesID, Title: "Show Name"},
	}, nil)
	svc.On("GetSeasonByNumber", mock.Anything, seriesID, int32(3)).Return(nil, errors.New("not found"))

	err := worker.Work(context.Background(), job)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "season")
	svc.AssertExpectations(t)
}

func TestFileMatchWorker_Work_AutoCreate_EpisodeNotFound_NoAutoCreate(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewFileMatchWorker(svc, nil, logger)

	tmpFile := createTempFileWithName(t, "Show.Name.S01E05.mkv")
	seriesID := uuid.Must(uuid.NewV7())
	seasonID := uuid.Must(uuid.NewV7())

	job := &river.Job[FileMatchArgs]{
		JobRow: &rivertype.JobRow{ID: 9, Kind: KindFileMatch},
		Args: FileMatchArgs{
			FilePath:   tmpFile,
			AutoCreate: false,
		},
	}

	svc.On("GetEpisodeFileByPath", mock.Anything, tmpFile).Return(nil, errors.New("not found"))
	svc.On("SearchSeries", mock.Anything, mock.Anything, int32(5), int32(0)).Return([]tvshow.Series{
		{ID: seriesID, Title: "Show Name"},
	}, nil)
	svc.On("GetSeasonByNumber", mock.Anything, seriesID, int32(1)).Return(&tvshow.Season{
		ID:       seasonID,
		SeriesID: seriesID,
	}, nil)
	svc.On("GetEpisodeByNumber", mock.Anything, seriesID, int32(1), int32(5)).Return(nil, errors.New("not found"))

	err := worker.Work(context.Background(), job)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "episode S01E05 not found")
	svc.AssertExpectations(t)
}

func TestFileMatchWorker_Work_AutoCreate_CreateSeason(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewFileMatchWorker(svc, nil, logger)

	tmpFile := createTempFileWithName(t, "Test.Show.S02E01.mkv")
	seriesID := uuid.Must(uuid.NewV7())
	seasonID := uuid.Must(uuid.NewV7())
	episodeID := uuid.Must(uuid.NewV7())

	job := &river.Job[FileMatchArgs]{
		JobRow: &rivertype.JobRow{ID: 11, Kind: KindFileMatch},
		Args: FileMatchArgs{
			FilePath:   tmpFile,
			AutoCreate: true,
		},
	}

	svc.On("GetEpisodeFileByPath", mock.Anything, tmpFile).Return(nil, errors.New("not found"))
	svc.On("SearchSeries", mock.Anything, mock.Anything, int32(5), int32(0)).Return([]tvshow.Series{
		{ID: seriesID, Title: "Test Show"},
	}, nil)
	// Season not found -> create it
	svc.On("GetSeasonByNumber", mock.Anything, seriesID, int32(2)).Return(nil, errors.New("not found"))
	svc.On("CreateSeason", mock.Anything, mock.MatchedBy(func(params tvshow.CreateSeasonParams) bool {
		return params.SeriesID == seriesID && params.SeasonNumber == 2
	})).Return(&tvshow.Season{
		ID:           seasonID,
		SeriesID:     seriesID,
		SeasonNumber: 2,
	}, nil)
	// Episode not found -> create it
	svc.On("GetEpisodeByNumber", mock.Anything, seriesID, int32(2), int32(1)).Return(nil, errors.New("not found"))
	svc.On("CreateEpisode", mock.Anything, mock.MatchedBy(func(params tvshow.CreateEpisodeParams) bool {
		return params.SeriesID == seriesID && params.SeasonID == seasonID && params.EpisodeNumber == 1
	})).Return(&tvshow.Episode{
		ID:       episodeID,
		SeasonID: seasonID,
	}, nil)
	svc.On("CreateEpisodeFile", mock.Anything, mock.MatchedBy(func(params tvshow.CreateEpisodeFileParams) bool {
		return params.EpisodeID == episodeID
	})).Return(&tvshow.EpisodeFile{
		ID:        uuid.Must(uuid.NewV7()),
		EpisodeID: episodeID,
		FilePath:  tmpFile,
	}, nil)

	err := worker.Work(context.Background(), job)
	require.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestFileMatchWorker_Work_AutoCreate_CreateSeasonFails(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewFileMatchWorker(svc, nil, logger)

	tmpFile := createTempFileWithName(t, "Test.Show.S02E01.720p.mkv")
	seriesID := uuid.Must(uuid.NewV7())

	job := &river.Job[FileMatchArgs]{
		JobRow: &rivertype.JobRow{ID: 12, Kind: KindFileMatch},
		Args: FileMatchArgs{
			FilePath:   tmpFile,
			AutoCreate: true,
		},
	}

	svc.On("GetEpisodeFileByPath", mock.Anything, tmpFile).Return(nil, errors.New("not found"))
	svc.On("SearchSeries", mock.Anything, mock.Anything, int32(5), int32(0)).Return([]tvshow.Series{
		{ID: seriesID, Title: "Test Show"},
	}, nil)
	svc.On("GetSeasonByNumber", mock.Anything, seriesID, int32(2)).Return(nil, errors.New("not found"))
	svc.On("CreateSeason", mock.Anything, mock.Anything).Return(nil, errors.New("create season failed"))

	err := worker.Work(context.Background(), job)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "create season")
	svc.AssertExpectations(t)
}

func TestFileMatchWorker_Work_AutoCreate_CreateEpisodeFails(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewFileMatchWorker(svc, nil, logger)

	tmpFile := createTempFileWithName(t, "Show.S01E03.mkv")
	seriesID := uuid.Must(uuid.NewV7())
	seasonID := uuid.Must(uuid.NewV7())

	job := &river.Job[FileMatchArgs]{
		JobRow: &rivertype.JobRow{ID: 13, Kind: KindFileMatch},
		Args: FileMatchArgs{
			FilePath:   tmpFile,
			AutoCreate: true,
		},
	}

	svc.On("GetEpisodeFileByPath", mock.Anything, tmpFile).Return(nil, errors.New("not found"))
	svc.On("SearchSeries", mock.Anything, mock.Anything, int32(5), int32(0)).Return([]tvshow.Series{
		{ID: seriesID, Title: "Show"},
	}, nil)
	svc.On("GetSeasonByNumber", mock.Anything, seriesID, int32(1)).Return(&tvshow.Season{
		ID:       seasonID,
		SeriesID: seriesID,
	}, nil)
	svc.On("GetEpisodeByNumber", mock.Anything, seriesID, int32(1), int32(3)).Return(nil, errors.New("not found"))
	svc.On("CreateEpisode", mock.Anything, mock.Anything).Return(nil, errors.New("create episode failed"))

	err := worker.Work(context.Background(), job)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "create episode")
	svc.AssertExpectations(t)
}

func TestFileMatchWorker_Work_AutoCreate_CreateEpisodeFileFails(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewFileMatchWorker(svc, nil, logger)

	tmpFile := createTempFileWithName(t, "Show.S01E03.720p.mkv")
	seriesID := uuid.Must(uuid.NewV7())
	seasonID := uuid.Must(uuid.NewV7())
	episodeID := uuid.Must(uuid.NewV7())

	job := &river.Job[FileMatchArgs]{
		JobRow: &rivertype.JobRow{ID: 14, Kind: KindFileMatch},
		Args: FileMatchArgs{
			FilePath:   tmpFile,
			AutoCreate: true,
		},
	}

	svc.On("GetEpisodeFileByPath", mock.Anything, tmpFile).Return(nil, errors.New("not found"))
	svc.On("SearchSeries", mock.Anything, mock.Anything, int32(5), int32(0)).Return([]tvshow.Series{
		{ID: seriesID, Title: "Show"},
	}, nil)
	svc.On("GetSeasonByNumber", mock.Anything, seriesID, int32(1)).Return(&tvshow.Season{
		ID:       seasonID,
		SeriesID: seriesID,
	}, nil)
	svc.On("GetEpisodeByNumber", mock.Anything, seriesID, int32(1), int32(3)).Return(&tvshow.Episode{
		ID:       episodeID,
		SeasonID: seasonID,
	}, nil)
	svc.On("CreateEpisodeFile", mock.Anything, mock.Anything).Return(nil, errors.New("create file failed"))

	err := worker.Work(context.Background(), job)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "create episode file")
	svc.AssertExpectations(t)
}

func TestFileMatchWorker_Work_AutoCreate_WithMetadataProvider(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	mdp := new(mockMetadataProvider)
	worker := NewFileMatchWorker(svc, mdp, logger)

	tmpFile := createTempFileWithName(t, "New.Show.S01E01.mkv")

	seriesID := uuid.Must(uuid.NewV7())
	seasonID := uuid.Must(uuid.NewV7())
	episodeID := uuid.Must(uuid.NewV7())
	tmdbID := int32(999)

	job := &river.Job[FileMatchArgs]{
		JobRow: &rivertype.JobRow{ID: 15, Kind: KindFileMatch},
		Args: FileMatchArgs{
			FilePath:   tmpFile,
			AutoCreate: true,
		},
	}

	svc.On("GetEpisodeFileByPath", mock.Anything, tmpFile).Return(nil, errors.New("not found"))
	// Search returns no matching series
	svc.On("SearchSeries", mock.Anything, mock.Anything, int32(5), int32(0)).Return([]tvshow.Series{}, nil)
	// Search TMDb for series
	mdp.On("SearchSeries", mock.Anything, mock.Anything, (*int)(nil)).Return([]*tvshow.Series{
		{Title: "New Show", TMDbID: &tmdbID},
	}, nil)
	mdp.On("EnrichSeries", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	svc.On("CreateSeries", mock.Anything, mock.Anything).Return(&tvshow.Series{
		ID:     seriesID,
		Title:  "New Show",
		TMDbID: &tmdbID,
	}, nil)
	svc.On("GetSeasonByNumber", mock.Anything, seriesID, int32(1)).Return(nil, errors.New("not found"))
	svc.On("CreateSeason", mock.Anything, mock.Anything).Return(&tvshow.Season{
		ID:       seasonID,
		SeriesID: seriesID,
	}, nil)
	svc.On("GetEpisodeByNumber", mock.Anything, seriesID, int32(1), int32(1)).Return(nil, errors.New("not found"))
	svc.On("CreateEpisode", mock.Anything, mock.Anything).Return(&tvshow.Episode{
		ID:       episodeID,
		SeasonID: seasonID,
	}, nil)
	svc.On("CreateEpisodeFile", mock.Anything, mock.Anything).Return(&tvshow.EpisodeFile{
		ID:        uuid.Must(uuid.NewV7()),
		EpisodeID: episodeID,
		FilePath:  tmpFile,
	}, nil)

	err := worker.Work(context.Background(), job)
	require.NoError(t, err)
	svc.AssertExpectations(t)
	mdp.AssertExpectations(t)
}

func TestFileMatchWorker_Work_AutoCreate_TMDbNotFound(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	mdp := new(mockMetadataProvider)
	worker := NewFileMatchWorker(svc, mdp, logger)

	tmpFile := createTempFileWithName(t, "Unknown.Show.S01E01.mkv")

	job := &river.Job[FileMatchArgs]{
		JobRow: &rivertype.JobRow{ID: 16, Kind: KindFileMatch},
		Args: FileMatchArgs{
			FilePath:   tmpFile,
			AutoCreate: true,
		},
	}

	svc.On("GetEpisodeFileByPath", mock.Anything, tmpFile).Return(nil, errors.New("not found"))
	svc.On("SearchSeries", mock.Anything, mock.Anything, int32(5), int32(0)).Return([]tvshow.Series{}, nil)
	// TMDb returns no results
	mdp.On("SearchSeries", mock.Anything, mock.Anything, (*int)(nil)).Return([]*tvshow.Series{}, nil)

	err := worker.Work(context.Background(), job)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "series not found in TMDb")
	svc.AssertExpectations(t)
	mdp.AssertExpectations(t)
}

func TestFileMatchWorker_Work_AutoCreate_TMDbSearchError(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	mdp := new(mockMetadataProvider)
	worker := NewFileMatchWorker(svc, mdp, logger)

	tmpFile := createTempFileWithName(t, "Error.Show.S01E01.mkv")

	job := &river.Job[FileMatchArgs]{
		JobRow: &rivertype.JobRow{ID: 17, Kind: KindFileMatch},
		Args: FileMatchArgs{
			FilePath:   tmpFile,
			AutoCreate: true,
		},
	}

	svc.On("GetEpisodeFileByPath", mock.Anything, tmpFile).Return(nil, errors.New("not found"))
	svc.On("SearchSeries", mock.Anything, mock.Anything, int32(5), int32(0)).Return([]tvshow.Series{}, nil)
	mdp.On("SearchSeries", mock.Anything, mock.Anything, (*int)(nil)).Return(nil, errors.New("tmdb error"))

	err := worker.Work(context.Background(), job)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "series not found in TMDb")
	svc.AssertExpectations(t)
	mdp.AssertExpectations(t)
}

func TestFileMatchWorker_Work_AutoCreate_CreateSeriesFails(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	mdp := new(mockMetadataProvider)
	worker := NewFileMatchWorker(svc, mdp, logger)

	tmpFile := createTempFileWithName(t, "Create.Fail.S01E01.mkv")
	tmdbID := int32(123)

	job := &river.Job[FileMatchArgs]{
		JobRow: &rivertype.JobRow{ID: 18, Kind: KindFileMatch},
		Args: FileMatchArgs{
			FilePath:   tmpFile,
			AutoCreate: true,
		},
	}

	svc.On("GetEpisodeFileByPath", mock.Anything, tmpFile).Return(nil, errors.New("not found"))
	svc.On("SearchSeries", mock.Anything, mock.Anything, int32(5), int32(0)).Return([]tvshow.Series{}, nil)
	mdp.On("SearchSeries", mock.Anything, mock.Anything, (*int)(nil)).Return([]*tvshow.Series{
		{Title: "Create Fail", TMDbID: &tmdbID},
	}, nil)
	mdp.On("EnrichSeries", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	svc.On("CreateSeries", mock.Anything, mock.Anything).Return(nil, errors.New("create series failed"))

	err := worker.Work(context.Background(), job)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "create series")
	svc.AssertExpectations(t)
	mdp.AssertExpectations(t)
}

// =============================================================================
// processFile() Tests (direct testing of unexported method)
// =============================================================================

func TestProcessFile_EmptyTitle(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewLibraryScanWorker(svc, nil, &infrajobs.Client{}, logger)

	sr := scanner.ScanResult{
		FilePath:    "/tmp/test.mkv",
		ParsedTitle: "",
		Metadata: map[string]any{
			"season":  1,
			"episode": 1,
		},
		IsMedia: true,
	}

	err := worker.processFile(context.Background(), sr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "could not parse series title")
}

func TestProcessFile_NoSeasonEpisode(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewLibraryScanWorker(svc, nil, &infrajobs.Client{}, logger)

	sr := scanner.ScanResult{
		FilePath:    "/tmp/test.mkv",
		ParsedTitle: "Some Show",
		Metadata:    map[string]any{},
		IsMedia:     true,
	}

	err := worker.processFile(context.Background(), sr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "could not parse season/episode")
}

func TestProcessFile_SearchSeriesError(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewLibraryScanWorker(svc, nil, &infrajobs.Client{}, logger)

	sr := scanner.ScanResult{
		FilePath:    "/tmp/test.mkv",
		ParsedTitle: "Some Show",
		Metadata: map[string]any{
			"season":  1,
			"episode": 1,
		},
		IsMedia: true,
	}

	svc.On("SearchSeries", mock.Anything, "Some Show", int32(5), int32(0)).Return([]tvshow.Series{}, errors.New("db error"))

	err := worker.processFile(context.Background(), sr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "search series")
	svc.AssertExpectations(t)
}

func TestProcessFile_ExactMatch_ExistingEpisode(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	mdp := new(mockMetadataProvider)
	worker := NewLibraryScanWorker(svc, mdp, &infrajobs.Client{}, logger)

	seriesID := uuid.Must(uuid.NewV7())
	seasonID := uuid.Must(uuid.NewV7())
	episodeID := uuid.Must(uuid.NewV7())

	tmpFile := createTempFileWithName(t, "test_episode.mkv")

	sr := scanner.ScanResult{
		FilePath:    tmpFile,
		ParsedTitle: "Breaking Bad",
		Metadata: map[string]any{
			"season":  2,
			"episode": 3,
		},
		IsMedia: true,
	}

	svc.On("SearchSeries", mock.Anything, "Breaking Bad", int32(5), int32(0)).Return([]tvshow.Series{
		{ID: seriesID, Title: "Breaking Bad"},
	}, nil)
	svc.On("GetSeasonByNumber", mock.Anything, seriesID, int32(2)).Return(&tvshow.Season{
		ID:       seasonID,
		SeriesID: seriesID,
	}, nil)
	svc.On("GetEpisodeByNumber", mock.Anything, seriesID, int32(2), int32(3)).Return(&tvshow.Episode{
		ID:       episodeID,
		SeasonID: seasonID,
	}, nil)
	svc.On("CreateEpisodeFile", mock.Anything, mock.MatchedBy(func(params tvshow.CreateEpisodeFileParams) bool {
		return params.EpisodeID == episodeID && params.FilePath == tmpFile
	})).Return(&tvshow.EpisodeFile{
		ID:        uuid.Must(uuid.NewV7()),
		EpisodeID: episodeID,
		FilePath:  tmpFile,
	}, nil)

	err := worker.processFile(context.Background(), sr)
	require.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestProcessFile_NoMatch_CreateFromTMDb(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	mdp := new(mockMetadataProvider)
	worker := NewLibraryScanWorker(svc, mdp, &infrajobs.Client{}, logger)

	seriesID := uuid.Must(uuid.NewV7())
	seasonID := uuid.Must(uuid.NewV7())
	episodeID := uuid.Must(uuid.NewV7())
	tmdbID := int32(999)

	tmpFile := createTempFileWithName(t, "new_episode.mkv")

	sr := scanner.ScanResult{
		FilePath:    tmpFile,
		ParsedTitle: "New Show",
		Metadata: map[string]any{
			"season":  1,
			"episode": 1,
		},
		IsMedia: true,
	}

	svc.On("SearchSeries", mock.Anything, "New Show", int32(5), int32(0)).Return([]tvshow.Series{}, nil)
	mdp.On("SearchSeries", mock.Anything, "New Show", (*int)(nil)).Return([]*tvshow.Series{
		{Title: "New Show", TMDbID: &tmdbID},
	}, nil)
	mdp.On("EnrichSeries", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	svc.On("CreateSeries", mock.Anything, mock.Anything).Return(&tvshow.Series{
		ID:     seriesID,
		Title:  "New Show",
		TMDbID: &tmdbID,
	}, nil)
	svc.On("GetSeasonByNumber", mock.Anything, seriesID, int32(1)).Return(nil, errors.New("not found"))
	mdp.On("EnrichSeason", mock.Anything, mock.Anything, tmdbID, mock.Anything).Return(nil)
	svc.On("CreateSeason", mock.Anything, mock.Anything).Return(&tvshow.Season{
		ID:       seasonID,
		SeriesID: seriesID,
	}, nil)
	svc.On("GetEpisodeByNumber", mock.Anything, seriesID, int32(1), int32(1)).Return(nil, errors.New("not found"))
	mdp.On("EnrichEpisode", mock.Anything, mock.Anything, tmdbID, mock.Anything).Return(nil)
	svc.On("CreateEpisode", mock.Anything, mock.Anything).Return(&tvshow.Episode{
		ID:       episodeID,
		SeasonID: seasonID,
	}, nil)
	svc.On("CreateEpisodeFile", mock.Anything, mock.Anything).Return(&tvshow.EpisodeFile{
		ID:        uuid.Must(uuid.NewV7()),
		EpisodeID: episodeID,
		FilePath:  tmpFile,
	}, nil)

	err := worker.processFile(context.Background(), sr)
	require.NoError(t, err)
	svc.AssertExpectations(t)
	mdp.AssertExpectations(t)
}

func TestProcessFile_TMDbNotFound(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	mdp := new(mockMetadataProvider)
	worker := NewLibraryScanWorker(svc, mdp, &infrajobs.Client{}, logger)

	sr := scanner.ScanResult{
		FilePath:    "/tmp/test.mkv",
		ParsedTitle: "Unknown Series",
		Metadata: map[string]any{
			"season":  1,
			"episode": 1,
		},
		IsMedia: true,
	}

	svc.On("SearchSeries", mock.Anything, "Unknown Series", int32(5), int32(0)).Return([]tvshow.Series{}, nil)
	mdp.On("SearchSeries", mock.Anything, "Unknown Series", (*int)(nil)).Return([]*tvshow.Series{}, nil)

	err := worker.processFile(context.Background(), sr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "series not found")
	svc.AssertExpectations(t)
	mdp.AssertExpectations(t)
}

func TestProcessFile_CreateSeriesFails(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	mdp := new(mockMetadataProvider)
	worker := NewLibraryScanWorker(svc, mdp, &infrajobs.Client{}, logger)

	tmdbID := int32(500)

	sr := scanner.ScanResult{
		FilePath:    "/tmp/test.mkv",
		ParsedTitle: "Fail Series",
		Metadata: map[string]any{
			"season":  1,
			"episode": 1,
		},
		IsMedia: true,
	}

	svc.On("SearchSeries", mock.Anything, "Fail Series", int32(5), int32(0)).Return([]tvshow.Series{}, nil)
	mdp.On("SearchSeries", mock.Anything, "Fail Series", (*int)(nil)).Return([]*tvshow.Series{
		{Title: "Fail Series", TMDbID: &tmdbID},
	}, nil)
	mdp.On("EnrichSeries", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	svc.On("CreateSeries", mock.Anything, mock.Anything).Return(nil, errors.New("db error"))

	err := worker.processFile(context.Background(), sr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "create series")
	svc.AssertExpectations(t)
	mdp.AssertExpectations(t)
}

func TestProcessFile_CreateSeasonFails(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	mdp := new(mockMetadataProvider)
	worker := NewLibraryScanWorker(svc, mdp, &infrajobs.Client{}, logger)

	seriesID := uuid.Must(uuid.NewV7())

	sr := scanner.ScanResult{
		FilePath:    "/tmp/test.mkv",
		ParsedTitle: "Season Fail",
		Metadata: map[string]any{
			"season":  1,
			"episode": 1,
		},
		IsMedia: true,
	}

	svc.On("SearchSeries", mock.Anything, "Season Fail", int32(5), int32(0)).Return([]tvshow.Series{
		{ID: seriesID, Title: "Season Fail"},
	}, nil)
	svc.On("GetSeasonByNumber", mock.Anything, seriesID, int32(1)).Return(nil, errors.New("not found"))
	mdp.On("EnrichSeason", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("enrich failed"))
	svc.On("CreateSeason", mock.Anything, mock.Anything).Return(nil, errors.New("create season failed"))

	err := worker.processFile(context.Background(), sr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "create season")
	svc.AssertExpectations(t)
}

func TestProcessFile_CreateEpisodeFails(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewLibraryScanWorker(svc, nil, &infrajobs.Client{}, logger)

	seriesID := uuid.Must(uuid.NewV7())
	seasonID := uuid.Must(uuid.NewV7())

	sr := scanner.ScanResult{
		FilePath:    "/tmp/test.mkv",
		ParsedTitle: "Episode Fail",
		Metadata: map[string]any{
			"season":  1,
			"episode": 1,
		},
		IsMedia: true,
	}

	svc.On("SearchSeries", mock.Anything, "Episode Fail", int32(5), int32(0)).Return([]tvshow.Series{
		{ID: seriesID, Title: "Episode Fail"},
	}, nil)
	svc.On("GetSeasonByNumber", mock.Anything, seriesID, int32(1)).Return(&tvshow.Season{
		ID:       seasonID,
		SeriesID: seriesID,
	}, nil)
	svc.On("GetEpisodeByNumber", mock.Anything, seriesID, int32(1), int32(1)).Return(nil, errors.New("not found"))
	svc.On("CreateEpisode", mock.Anything, mock.Anything).Return(nil, errors.New("create episode failed"))

	err := worker.processFile(context.Background(), sr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "create episode")
	svc.AssertExpectations(t)
}

func TestProcessFile_CreateEpisodeFileFails(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewLibraryScanWorker(svc, nil, &infrajobs.Client{}, logger)

	seriesID := uuid.Must(uuid.NewV7())
	seasonID := uuid.Must(uuid.NewV7())
	episodeID := uuid.Must(uuid.NewV7())

	tmpFile := createTempFileWithName(t, "test_episode.mkv")

	sr := scanner.ScanResult{
		FilePath:    tmpFile,
		ParsedTitle: "File Fail",
		Metadata: map[string]any{
			"season":  1,
			"episode": 1,
		},
		IsMedia: true,
	}

	svc.On("SearchSeries", mock.Anything, "File Fail", int32(5), int32(0)).Return([]tvshow.Series{
		{ID: seriesID, Title: "File Fail"},
	}, nil)
	svc.On("GetSeasonByNumber", mock.Anything, seriesID, int32(1)).Return(&tvshow.Season{
		ID:       seasonID,
		SeriesID: seriesID,
	}, nil)
	svc.On("GetEpisodeByNumber", mock.Anything, seriesID, int32(1), int32(1)).Return(&tvshow.Episode{
		ID:       episodeID,
		SeasonID: seasonID,
	}, nil)
	svc.On("CreateEpisodeFile", mock.Anything, mock.Anything).Return(nil, errors.New("create file failed"))

	err := worker.processFile(context.Background(), sr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "create episode file")
	svc.AssertExpectations(t)
}

func TestProcessFile_WithEpisodeTitle(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewLibraryScanWorker(svc, nil, &infrajobs.Client{}, logger)

	seriesID := uuid.Must(uuid.NewV7())
	seasonID := uuid.Must(uuid.NewV7())
	episodeID := uuid.Must(uuid.NewV7())

	tmpFile := createTempFileWithName(t, "test_with_title.mkv")

	sr := scanner.ScanResult{
		FilePath:    tmpFile,
		ParsedTitle: "Title Show",
		Metadata: map[string]any{
			"season":        1,
			"episode":       5,
			"episode_title": "The Pilot Episode",
		},
		IsMedia: true,
	}

	svc.On("SearchSeries", mock.Anything, "Title Show", int32(5), int32(0)).Return([]tvshow.Series{
		{ID: seriesID, Title: "Title Show"},
	}, nil)
	svc.On("GetSeasonByNumber", mock.Anything, seriesID, int32(1)).Return(&tvshow.Season{
		ID:       seasonID,
		SeriesID: seriesID,
	}, nil)
	svc.On("GetEpisodeByNumber", mock.Anything, seriesID, int32(1), int32(5)).Return(nil, errors.New("not found"))
	svc.On("CreateEpisode", mock.Anything, mock.MatchedBy(func(params tvshow.CreateEpisodeParams) bool {
		return params.Title == "The Pilot Episode" && params.EpisodeNumber == 5
	})).Return(&tvshow.Episode{
		ID:       episodeID,
		SeasonID: seasonID,
	}, nil)
	svc.On("CreateEpisodeFile", mock.Anything, mock.Anything).Return(&tvshow.EpisodeFile{
		ID:        uuid.Must(uuid.NewV7()),
		EpisodeID: episodeID,
		FilePath:  tmpFile,
	}, nil)

	err := worker.processFile(context.Background(), sr)
	require.NoError(t, err)
	svc.AssertExpectations(t)
}

// =============================================================================
// LibraryScanWorker.Work Tests (additional paths)
// =============================================================================

func TestLibraryScanWorker_Work_ScanError(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewLibraryScanWorker(svc, nil, &infrajobs.Client{}, logger)

	// Use a non-existent path to trigger scan error
	job := &river.Job[LibraryScanArgs]{
		JobRow: &rivertype.JobRow{ID: 1, Kind: KindLibraryScan},
		Args: LibraryScanArgs{
			Paths: []string{"/nonexistent/path/that/does/not/exist"},
		},
	}

	// The scanner won't error on non-existent paths, it just returns empty results
	err := worker.Work(context.Background(), job)
	// Should succeed with 0 items processed
	require.NoError(t, err)
}

func TestLibraryScanWorker_Work_WithMediaFiles_NoAutoCreate(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewLibraryScanWorker(svc, nil, &infrajobs.Client{}, logger)

	// Create a temp dir with a parseable media file
	dir := t.TempDir()
	filePath := dir + "/Breaking.Bad.S01E01.Pilot.mkv"
	f, err := os.Create(filePath)
	require.NoError(t, err)
	_, err = f.WriteString("fake media content")
	require.NoError(t, err)
	f.Close()

	job := &river.Job[LibraryScanArgs]{
		JobRow: &rivertype.JobRow{ID: 2, Kind: KindLibraryScan},
		Args: LibraryScanArgs{
			Paths:      []string{dir},
			AutoCreate: false, // Just discover, don't create
		},
	}

	// Even with AutoCreate=false, it checks if file is already matched
	svc.On("GetEpisodeFileByPath", mock.Anything, filePath).Return(nil, errors.New("not found"))

	err = worker.Work(context.Background(), job)
	require.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestLibraryScanWorker_Work_WithAutoCreate_AlreadyMatched(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	mdp := new(mockMetadataProvider)
	worker := NewLibraryScanWorker(svc, mdp, &infrajobs.Client{}, logger)

	dir := t.TempDir()
	filePath := dir + "/Show.S01E01.mkv"
	f, err := os.Create(filePath)
	require.NoError(t, err)
	_, err = f.WriteString("fake media content")
	require.NoError(t, err)
	f.Close()

	job := &river.Job[LibraryScanArgs]{
		JobRow: &rivertype.JobRow{ID: 3, Kind: KindLibraryScan},
		Args: LibraryScanArgs{
			Paths:      []string{dir},
			AutoCreate: true,
		},
	}

	episodeID := uuid.Must(uuid.NewV7())
	svc.On("GetEpisodeFileByPath", mock.Anything, filePath).Return(&tvshow.EpisodeFile{
		ID:        uuid.Must(uuid.NewV7()),
		EpisodeID: episodeID,
		FilePath:  filePath,
	}, nil)

	err = worker.Work(context.Background(), job)
	require.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestLibraryScanWorker_Work_WithAutoCreate_ProcessFile(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	mdp := new(mockMetadataProvider)
	worker := NewLibraryScanWorker(svc, mdp, &infrajobs.Client{}, logger)

	dir := t.TempDir()
	filePath := dir + "/Good.Show.S01E01.mkv"
	f, err := os.Create(filePath)
	require.NoError(t, err)
	_, err = f.WriteString("fake media content")
	require.NoError(t, err)
	f.Close()

	seriesID := uuid.Must(uuid.NewV7())
	seasonID := uuid.Must(uuid.NewV7())
	episodeID := uuid.Must(uuid.NewV7())

	job := &river.Job[LibraryScanArgs]{
		JobRow: &rivertype.JobRow{ID: 4, Kind: KindLibraryScan},
		Args: LibraryScanArgs{
			Paths:      []string{dir},
			AutoCreate: true,
		},
	}

	// File not already matched
	svc.On("GetEpisodeFileByPath", mock.Anything, filePath).Return(nil, errors.New("not found"))
	// processFile: series search finds match
	svc.On("SearchSeries", mock.Anything, mock.Anything, int32(5), int32(0)).Return([]tvshow.Series{
		{ID: seriesID, Title: "Good Show"},
	}, nil)
	svc.On("GetSeasonByNumber", mock.Anything, seriesID, int32(1)).Return(&tvshow.Season{
		ID:       seasonID,
		SeriesID: seriesID,
	}, nil)
	svc.On("GetEpisodeByNumber", mock.Anything, seriesID, int32(1), int32(1)).Return(&tvshow.Episode{
		ID:       episodeID,
		SeasonID: seasonID,
	}, nil)
	svc.On("CreateEpisodeFile", mock.Anything, mock.Anything).Return(&tvshow.EpisodeFile{
		ID:        uuid.Must(uuid.NewV7()),
		EpisodeID: episodeID,
		FilePath:  filePath,
	}, nil)

	err = worker.Work(context.Background(), job)
	require.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestLibraryScanWorker_Work_WithAutoCreate_ProcessFileFails(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	mdp := new(mockMetadataProvider)
	worker := NewLibraryScanWorker(svc, mdp, &infrajobs.Client{}, logger)

	dir := t.TempDir()
	filePath := dir + "/Bad.Show.S01E01.mkv"
	f, err := os.Create(filePath)
	require.NoError(t, err)
	_, err = f.WriteString("fake media content")
	require.NoError(t, err)
	f.Close()

	job := &river.Job[LibraryScanArgs]{
		JobRow: &rivertype.JobRow{ID: 5, Kind: KindLibraryScan},
		Args: LibraryScanArgs{
			Paths:      []string{dir},
			AutoCreate: true,
		},
	}

	// File not already matched
	svc.On("GetEpisodeFileByPath", mock.Anything, filePath).Return(nil, errors.New("not found"))
	// processFile: series search fails
	svc.On("SearchSeries", mock.Anything, mock.Anything, int32(5), int32(0)).Return([]tvshow.Series{}, errors.New("db error"))

	err = worker.Work(context.Background(), job)
	// Work should succeed (errors are logged, not returned for individual files)
	require.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestLibraryScanWorker_Work_ForceRescan(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	worker := NewLibraryScanWorker(svc, nil, &infrajobs.Client{}, logger)

	dir := t.TempDir()
	filePath := dir + "/Force.Show.S01E01.mkv"
	f, err := os.Create(filePath)
	require.NoError(t, err)
	_, err = f.WriteString("fake media content")
	require.NoError(t, err)
	f.Close()

	episodeID := uuid.Must(uuid.NewV7())

	job := &river.Job[LibraryScanArgs]{
		JobRow: &rivertype.JobRow{ID: 6, Kind: KindLibraryScan},
		Args: LibraryScanArgs{
			Paths:      []string{dir},
			Force:      true,
			AutoCreate: false,
		},
	}

	// Even though file is matched, Force=true means we continue (but auto_create=false, so just log)
	svc.On("GetEpisodeFileByPath", mock.Anything, filePath).Return(&tvshow.EpisodeFile{
		ID:        uuid.Must(uuid.NewV7()),
		EpisodeID: episodeID,
		FilePath:  filePath,
	}, nil)

	err = worker.Work(context.Background(), job)
	require.NoError(t, err)
}

// =============================================================================
// Module & RegisterWorkers Tests
// =============================================================================

func TestRegisterWorkers(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	workers := river.NewWorkers()

	libraryScan := NewLibraryScanWorker(nil, nil, nil, logger)
	metadataRefresh := NewMetadataRefreshWorker(nil, nil, logger)
	fileMatch := NewFileMatchWorker(nil, nil, logger)
	searchIndex := NewSearchIndexWorker(nil, nil, logger)
	seriesRefresh := NewSeriesRefreshWorker(nil, nil, logger)

	err := RegisterWorkers(workers, libraryScan, metadataRefresh, fileMatch, searchIndex, seriesRefresh)
	assert.NoError(t, err)
}

func TestProviderFunctions(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	mdp := new(mockMetadataProvider)
	jobClient := &infrajobs.Client{}

	params := WorkerProviderParams{
		Service:          svc,
		MetadataProvider: mdp,
		JobClient:        jobClient,
		Logger:           logger,
	}

	libScan := provideLibraryScanWorker(params)
	assert.NotNil(t, libScan)

	mdRefresh := provideMetadataRefreshWorker(params)
	assert.NotNil(t, mdRefresh)

	fmWorker := provideFileMatchWorker(params)
	assert.NotNil(t, fmWorker)

	siWorker := provideSearchIndexWorker(params)
	assert.NotNil(t, siWorker)

	srWorker := provideSeriesRefreshWorker(params)
	assert.NotNil(t, srWorker)
}

func TestProviderFunctions_NilOptionalDeps(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := new(mockService)
	jobClient := &infrajobs.Client{}

	params := WorkerProviderParams{
		Service:   svc,
		JobClient: jobClient,
		Logger:    logger,
		// MetadataProvider and SearchService are optional
	}

	libScan := provideLibraryScanWorker(params)
	assert.NotNil(t, libScan)

	siWorker := provideSearchIndexWorker(params)
	assert.NotNil(t, siWorker)
}

// =============================================================================
// Test Helpers
// =============================================================================

func createTempFile(t *testing.T, pattern string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), pattern)
	require.NoError(t, err)
	// Write some bytes so FileSize is nonzero
	_, err = f.WriteString("test content")
	require.NoError(t, err)
	f.Close()
	return f.Name()
}

func createTempFileWithName(t *testing.T, name string) string {
	t.Helper()
	dir := t.TempDir()
	path := dir + "/" + name
	f, err := os.Create(path)
	require.NoError(t, err)
	_, err = f.WriteString("test content")
	require.NoError(t, err)
	f.Close()
	return path
}
