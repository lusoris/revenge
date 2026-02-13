package api

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/govalues/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/content"
	"github.com/lusoris/revenge/internal/content/movie"
)

// ============================================================================
// externalRatingsToOgen
// ============================================================================

func TestExternalRatingsToOgen(t *testing.T) {
	t.Run("converts ratings", func(t *testing.T) {
		ratings := []content.ExternalRating{
			{Source: "IMDb", Value: "9.0/10", Score: 90.0},
			{Source: "Rotten Tomatoes", Value: "94%", Score: 94.0},
		}
		result := externalRatingsToOgen(ratings)
		require.Len(t, result, 2)
		assert.Equal(t, "IMDb", result[0].Source)
		assert.Equal(t, "9.0/10", result[0].Value)
		assert.InDelta(t, float32(90.0), result[0].Score, 0.1)
	})

	t.Run("nil returns nil", func(t *testing.T) {
		assert.Nil(t, externalRatingsToOgen(nil))
	})

	t.Run("empty returns nil", func(t *testing.T) {
		assert.Nil(t, externalRatingsToOgen([]content.ExternalRating{}))
	})
}

// ============================================================================
// movieToOgen
// ============================================================================

func TestMovieToOgen(t *testing.T) {
	t.Run("full movie", func(t *testing.T) {
		id := uuid.Must(uuid.NewV7())
		va, _ := decimal.NewFromFloat64(8.5)
		pop, _ := decimal.NewFromFloat64(100.5)
		now := time.Now()
		budget := int64(200000000)
		revenue := int64(1000000000)

		m := &movie.Movie{
			ID:                id,
			Title:             "Inception",
			OriginalTitle:     stringPtr("Inception"),
			OriginalLanguage:  stringPtr("en"),
			Year:              int32Ptr(2010),
			ReleaseDate:       &now,
			Runtime:           int32Ptr(148),
			Overview:          stringPtr("A thief enters dreams."),
			Tagline:           stringPtr("Your mind is the scene of the crime"),
			Status:            stringPtr("Released"),
			PosterPath:        stringPtr("/poster.jpg"),
			BackdropPath:      stringPtr("/backdrop.jpg"),
			TrailerURL:        stringPtr("https://youtube.com/watch?v=abc"),
			VoteAverage:       &va,
			VoteCount:         int32Ptr(30000),
			Popularity:        &pop,
			Budget:            &budget,
			Revenue:           &revenue,
			TMDbID:            int32Ptr(27205),
			IMDbID:            stringPtr("tt1375666"),
			MetadataUpdatedAt: &now,
			RadarrID:          int32Ptr(1),
			CreatedAt:         now,
			UpdatedAt:         now,
			LibraryAddedAt:    now,
			ExternalRatings: []movie.ExternalRating{
				{Source: "IMDb", Value: "8.8/10", Score: 88.0},
			},
		}

		o := movieToOgen(m)
		assert.Equal(t, id, o.ID.Value)
		assert.Equal(t, "Inception", o.Title.Value)
		assert.True(t, o.TmdbID.Set)
		assert.Equal(t, 27205, o.TmdbID.Value)
		assert.True(t, o.ImdbID.Set)
		assert.True(t, o.Year.Set)
		assert.Equal(t, 2010, o.Year.Value)
		assert.True(t, o.Runtime.Set)
		assert.Equal(t, 148, o.Runtime.Value)
		assert.True(t, o.Overview.Set)
		assert.True(t, o.Tagline.Set)
		assert.True(t, o.Status.Set)
		assert.True(t, o.PosterPath.Set)
		assert.True(t, o.BackdropPath.Set)
		assert.True(t, o.TrailerURL.Set)
		assert.True(t, o.VoteAverage.Set)
		assert.True(t, o.VoteCount.Set)
		assert.True(t, o.Popularity.Set)
		assert.True(t, o.Budget.Set)
		assert.Equal(t, int64(200000000), o.Budget.Value)
		assert.True(t, o.Revenue.Set)
		assert.True(t, o.MetadataUpdatedAt.Set)
		assert.True(t, o.RadarrID.Set)
		require.Len(t, o.ExternalRatings, 1)
	})

	t.Run("minimal movie", func(t *testing.T) {
		m := &movie.Movie{
			ID:    uuid.Must(uuid.NewV7()),
			Title: "Test",
		}
		o := movieToOgen(m)
		assert.Equal(t, "Test", o.Title.Value)
		assert.False(t, o.TmdbID.Set)
		assert.False(t, o.ImdbID.Set)
		assert.False(t, o.Year.Set)
		assert.False(t, o.Runtime.Set)
		assert.False(t, o.Overview.Set)
		assert.Nil(t, o.ExternalRatings)
	})
}

// ============================================================================
// movieFileToOgen
// ============================================================================

func TestMovieFileToOgen(t *testing.T) {
	t.Run("full file", func(t *testing.T) {
		now := time.Now()
		fr, _ := decimal.NewFromFloat64(23.976)
		f := &movie.MovieFile{
			ID:                uuid.Must(uuid.NewV7()),
			MovieID:           uuid.Must(uuid.NewV7()),
			FilePath:          "/movies/inception.mkv",
			FileSize:          42000000000,
			FileName:          "inception.mkv",
			Resolution:        stringPtr("2160p"),
			QualityProfile:    stringPtr("UHD Remux"),
			VideoCodec:        stringPtr("HEVC"),
			AudioCodec:        stringPtr("TrueHD Atmos"),
			Container:         stringPtr("mkv"),
			DurationSeconds:   int32Ptr(8880),
			BitrateKbps:       int32Ptr(35000),
			Framerate:         &fr,
			DynamicRange:      stringPtr("Dolby Vision"),
			ColorSpace:        stringPtr("BT.2020"),
			AudioChannels:     stringPtr("7.1"),
			AudioLanguages:    []string{"eng", "fra"},
			SubtitleLanguages: []string{"eng", "deu", "fra"},
			RadarrFileID:      int32Ptr(42),
			LastScannedAt:     &now,
			IsMonitored:       boolPtr(true),
			CreatedAt:         now,
			UpdatedAt:         now,
		}

		o := movieFileToOgen(f)
		assert.Equal(t, f.ID, o.ID.Value)
		assert.Equal(t, "inception.mkv", o.FileName.Value)
		assert.True(t, o.Resolution.Set)
		assert.Equal(t, "2160p", o.Resolution.Value)
		assert.True(t, o.VideoCodec.Set)
		assert.True(t, o.AudioCodec.Set)
		assert.True(t, o.DurationSeconds.Set)
		assert.True(t, o.BitrateKbps.Set)
		assert.True(t, o.Framerate.Set)
		assert.True(t, o.DynamicRange.Set)
		assert.True(t, o.ColorSpace.Set)
		assert.True(t, o.AudioChannels.Set)
		assert.Len(t, o.AudioLanguages, 2)
		assert.Len(t, o.SubtitleLanguages, 3)
		assert.True(t, o.RadarrFileID.Set)
		assert.True(t, o.LastScannedAt.Set)
		assert.True(t, o.IsMonitored.Set)
	})

	t.Run("minimal file", func(t *testing.T) {
		f := &movie.MovieFile{
			ID:       uuid.Must(uuid.NewV7()),
			MovieID:  uuid.Must(uuid.NewV7()),
			FilePath: "/test.mkv",
			FileName: "test.mkv",
		}
		o := movieFileToOgen(f)
		assert.False(t, o.Resolution.Set)
		assert.False(t, o.VideoCodec.Set)
		assert.Nil(t, o.AudioLanguages)
		assert.Nil(t, o.SubtitleLanguages)
	})
}

// ============================================================================
// movieCreditToOgen
// ============================================================================

func TestMovieCreditToOgen(t *testing.T) {
	t.Run("cast credit", func(t *testing.T) {
		c := &movie.MovieCredit{
			ID:           uuid.Must(uuid.NewV7()),
			MovieID:      uuid.Must(uuid.NewV7()),
			TMDbPersonID: 6193,
			Name:         "Leonardo DiCaprio",
			Character:    stringPtr("Cobb"),
			CreditType:   "cast",
			CastOrder:    int32Ptr(0),
			ProfilePath:  stringPtr("/leo.jpg"),
		}
		o := movieCreditToOgen(c)
		assert.Equal(t, "Leonardo DiCaprio", o.Name.Value)
		assert.True(t, o.Character.Set)
		assert.Equal(t, "Cobb", o.Character.Value)
		assert.True(t, o.CastOrder.Set)
		assert.False(t, o.Job.Set)
		assert.False(t, o.Department.Set)
	})

	t.Run("crew credit", func(t *testing.T) {
		c := &movie.MovieCredit{
			ID:           uuid.Must(uuid.NewV7()),
			MovieID:      uuid.Must(uuid.NewV7()),
			TMDbPersonID: 525,
			Name:         "Christopher Nolan",
			Job:          stringPtr("Director"),
			Department:   stringPtr("Directing"),
			CreditType:   "crew",
		}
		o := movieCreditToOgen(c)
		assert.True(t, o.Job.Set)
		assert.Equal(t, "Director", o.Job.Value)
		assert.True(t, o.Department.Set)
		assert.False(t, o.Character.Set)
		assert.False(t, o.CastOrder.Set)
	})
}

// ============================================================================
// movieCollectionToOgen
// ============================================================================

func TestMovieCollectionToOgen(t *testing.T) {
	t.Run("full collection", func(t *testing.T) {
		c := &movie.MovieCollection{
			ID:               uuid.Must(uuid.NewV7()),
			Name:             "The Dark Knight Collection",
			TMDbCollectionID: int32Ptr(263),
			Overview:         stringPtr("Batman trilogy"),
			PosterPath:       stringPtr("/collection.jpg"),
			BackdropPath:     stringPtr("/collection_bg.jpg"),
		}
		o := movieCollectionToOgen(c)
		assert.Equal(t, "The Dark Knight Collection", o.Name.Value)
		assert.True(t, o.TmdbCollectionID.Set)
		assert.Equal(t, 263, o.TmdbCollectionID.Value)
		assert.True(t, o.Overview.Set)
		assert.True(t, o.PosterPath.Set)
		assert.True(t, o.BackdropPath.Set)
	})

	t.Run("minimal collection", func(t *testing.T) {
		c := &movie.MovieCollection{
			ID:   uuid.Must(uuid.NewV7()),
			Name: "Test",
		}
		o := movieCollectionToOgen(c)
		assert.False(t, o.TmdbCollectionID.Set)
		assert.False(t, o.Overview.Set)
	})
}

// ============================================================================
// movieGenreToOgen
// ============================================================================

func TestMovieGenreToOgen(t *testing.T) {
	g := &movie.MovieGenre{
		ID:      uuid.Must(uuid.NewV7()),
		MovieID: uuid.Must(uuid.NewV7()),
		Slug:    "action",
		Name:    "Action",
	}
	o := movieGenreToOgen(g)
	assert.Equal(t, "Action", o.Name.Value)
	assert.Equal(t, "action", o.Slug.Value)
}

// ============================================================================
// movieWatchedToOgen
// ============================================================================

func TestMovieWatchedToOgen(t *testing.T) {
	t.Run("with progress percent", func(t *testing.T) {
		w := &movie.MovieWatched{
			ID:              uuid.Must(uuid.NewV7()),
			UserID:          uuid.Must(uuid.NewV7()),
			MovieID:         uuid.Must(uuid.NewV7()),
			ProgressSeconds: 3600,
			DurationSeconds: 7200,
			IsCompleted:     false,
			WatchCount:      1,
			ProgressPercent: int32Ptr(50),
			LastWatchedAt:   time.Now(),
		}
		o := movieWatchedToOgen(w)
		assert.Equal(t, 3600, o.ProgressSeconds.Value)
		assert.Equal(t, 7200, o.DurationSeconds.Value)
		assert.False(t, o.IsCompleted.Value)
		assert.True(t, o.ProgressPercent.Set)
		assert.Equal(t, 50, o.ProgressPercent.Value)
	})

	t.Run("without progress percent", func(t *testing.T) {
		w := &movie.MovieWatched{
			ID:      uuid.Must(uuid.NewV7()),
			UserID:  uuid.Must(uuid.NewV7()),
			MovieID: uuid.Must(uuid.NewV7()),
		}
		o := movieWatchedToOgen(w)
		assert.False(t, o.ProgressPercent.Set)
	})
}

// ============================================================================
// continueWatchingItemToOgen
// ============================================================================

func TestContinueWatchingItemToOgen(t *testing.T) {
	va, _ := decimal.NewFromFloat64(8.8)
	pop, _ := decimal.NewFromFloat64(50.0)
	now := time.Now()
	budget := int64(100000000)
	revenue := int64(500000000)

	item := &movie.ContinueWatchingItem{
		Movie: movie.Movie{
			ID:                uuid.Must(uuid.NewV7()),
			Title:             "Inception",
			TMDbID:            int32Ptr(27205),
			IMDbID:            stringPtr("tt1375666"),
			OriginalTitle:     stringPtr("Inception"),
			Year:              int32Ptr(2010),
			ReleaseDate:       &now,
			Runtime:           int32Ptr(148),
			Overview:          stringPtr("Dreams"),
			Tagline:           stringPtr("Your mind..."),
			Status:            stringPtr("Released"),
			OriginalLanguage:  stringPtr("en"),
			PosterPath:        stringPtr("/poster.jpg"),
			BackdropPath:      stringPtr("/backdrop.jpg"),
			TrailerURL:        stringPtr("https://yt.com/watch"),
			VoteAverage:       &va,
			VoteCount:         int32Ptr(30000),
			Popularity:        &pop,
			Budget:            &budget,
			Revenue:           &revenue,
			MetadataUpdatedAt: &now,
			RadarrID:          int32Ptr(1),
			ExternalRatings: []movie.ExternalRating{
				{Source: "IMDb", Value: "8.8/10", Score: 88.0},
			},
		},
		ProgressSeconds: 3600,
		DurationSeconds: 8880,
		ProgressPercent: int32Ptr(40),
		LastWatchedAt:   now,
	}

	o := continueWatchingItemToOgen(item)
	assert.Equal(t, "Inception", o.Title.Value)
	assert.Equal(t, 3600, o.ProgressSeconds.Value)
	assert.Equal(t, 8880, o.DurationSeconds.Value)
	assert.True(t, o.TmdbID.Set)
	assert.True(t, o.Year.Set)
	assert.True(t, o.Runtime.Set)
	assert.True(t, o.Overview.Set)
	assert.True(t, o.VoteAverage.Set)
	assert.True(t, o.Budget.Set)
	assert.True(t, o.Revenue.Set)
	assert.True(t, o.ProgressPercent.Set)
	assert.Equal(t, 40, o.ProgressPercent.Value)
	require.Len(t, o.ExternalRatings, 1)
}

// ============================================================================
// watchedMovieItemToOgen
// ============================================================================

func TestWatchedMovieItemToOgen(t *testing.T) {
	va, _ := decimal.NewFromFloat64(9.0)
	now := time.Now()

	item := &movie.WatchedMovieItem{
		Movie: movie.Movie{
			ID:          uuid.Must(uuid.NewV7()),
			Title:       "The Dark Knight",
			TMDbID:      int32Ptr(155),
			Year:        int32Ptr(2008),
			VoteAverage: &va,
			PosterPath:  stringPtr("/poster.jpg"),
			ExternalRatings: []movie.ExternalRating{
				{Source: "IMDb", Value: "9.0/10", Score: 90.0},
			},
		},
		WatchCount:    3,
		LastWatchedAt: now,
	}

	o := watchedMovieItemToOgen(item)
	assert.Equal(t, "The Dark Knight", o.Title.Value)
	assert.Equal(t, 3, o.WatchCount.Value)
	assert.True(t, o.TmdbID.Set)
	assert.True(t, o.Year.Set)
	assert.True(t, o.VoteAverage.Set)
	require.Len(t, o.ExternalRatings, 1)
}

// ============================================================================
// userMovieStatsToOgen
// ============================================================================

func TestUserMovieStatsToOgen(t *testing.T) {
	t.Run("with total watches", func(t *testing.T) {
		total := int64(42)
		stats := &movie.UserMovieStats{
			WatchedCount:    10,
			InProgressCount: 3,
			TotalWatches:    &total,
		}
		o := userMovieStatsToOgen(stats)
		assert.Equal(t, int64(10), o.WatchedCount.Value)
		assert.Equal(t, int64(3), o.InProgressCount.Value)
		assert.True(t, o.TotalWatches.Set)
		assert.Equal(t, int64(42), o.TotalWatches.Value)
	})

	t.Run("without total watches", func(t *testing.T) {
		stats := &movie.UserMovieStats{
			WatchedCount:    5,
			InProgressCount: 1,
		}
		o := userMovieStatsToOgen(stats)
		assert.False(t, o.TotalWatches.Set)
	})
}

// ============================================================================
// Helpers
// ============================================================================

func int32Ptr(v int32) *int32 { return &v }

func boolPtr(v bool) *bool { return &v }
