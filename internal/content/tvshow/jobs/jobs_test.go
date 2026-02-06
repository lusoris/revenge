package jobs

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestLibraryScanArgs_Kind(t *testing.T) {
	args := LibraryScanArgs{
		Paths: []string{"/tv"},
		Force: true,
	}
	assert.Equal(t, KindLibraryScan, args.Kind())
}

func TestMetadataRefreshArgs_Kind(t *testing.T) {
	seriesID := uuid.Must(uuid.NewV7())
	args := MetadataRefreshArgs{
		SeriesID: &seriesID,
		Force:    true,
	}
	assert.Equal(t, KindMetadataRefresh, args.Kind())
}

func TestFileMatchArgs_Kind(t *testing.T) {
	args := FileMatchArgs{
		FilePath:     "/tv/Breaking Bad/S01E01.mkv",
		ForceRematch: false,
		AutoCreate:   true,
	}
	assert.Equal(t, KindFileMatch, args.Kind())
}

func TestSearchIndexArgs_Kind(t *testing.T) {
	args := SearchIndexArgs{
		FullReindex: true,
	}
	assert.Equal(t, KindSearchIndex, args.Kind())
}

func TestSeriesRefreshArgs_Kind(t *testing.T) {
	args := SeriesRefreshArgs{
		SeriesID:        uuid.Must(uuid.NewV7()),
		TMDbID:          1396,
		RefreshSeasons:  true,
		RefreshEpisodes: true,
	}
	assert.Equal(t, KindSeriesRefresh, args.Kind())
}

func TestJobKinds(t *testing.T) {
	// Verify all job kinds are unique
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

func TestJobInsertOpts(t *testing.T) {
	t.Run("with priority only", func(t *testing.T) {
		opts := JobInsertOpts(HighPriority, nil)
		assert.NotNil(t, opts)
		assert.Equal(t, HighPriority, opts.Priority)
		assert.True(t, opts.ScheduledAt.IsZero())
	})

	t.Run("with scheduled time", func(t *testing.T) {
		scheduled := time.Now().Add(time.Hour)
		opts := JobInsertOpts(DefaultPriority, &scheduled)
		assert.NotNil(t, opts)
		assert.Equal(t, DefaultPriority, opts.Priority)
		assert.Equal(t, scheduled, opts.ScheduledAt)
	})
}

func TestPriorityConstants(t *testing.T) {
	// Higher number = lower priority in River
	assert.Less(t, HighPriority, DefaultPriority)
	assert.Less(t, DefaultPriority, LowPriority)
}

func TestLibraryScanArgs_Fields(t *testing.T) {
	libraryID := uuid.Must(uuid.NewV7())
	args := LibraryScanArgs{
		Paths:     []string{"/tv/shows", "/tv/anime"},
		Force:     true,
		LibraryID: &libraryID,
	}

	assert.Len(t, args.Paths, 2)
	assert.True(t, args.Force)
	assert.Equal(t, libraryID, *args.LibraryID)
}

func TestMetadataRefreshArgs_Fields(t *testing.T) {
	t.Run("series refresh", func(t *testing.T) {
		seriesID := uuid.Must(uuid.NewV7())
		args := MetadataRefreshArgs{
			SeriesID:      &seriesID,
			Force:         true,
			RefreshImages: true,
		}
		assert.NotNil(t, args.SeriesID)
		assert.True(t, args.Force)
		assert.True(t, args.RefreshImages)
		assert.Nil(t, args.SeasonID)
		assert.Nil(t, args.EpisodeID)
	})

	t.Run("season refresh", func(t *testing.T) {
		seasonID := uuid.Must(uuid.NewV7())
		args := MetadataRefreshArgs{
			SeasonID: &seasonID,
		}
		assert.Nil(t, args.SeriesID)
		assert.NotNil(t, args.SeasonID)
		assert.Nil(t, args.EpisodeID)
	})

	t.Run("episode refresh", func(t *testing.T) {
		episodeID := uuid.Must(uuid.NewV7())
		args := MetadataRefreshArgs{
			EpisodeID: &episodeID,
		}
		assert.Nil(t, args.SeriesID)
		assert.Nil(t, args.SeasonID)
		assert.NotNil(t, args.EpisodeID)
	})
}

func TestFileMatchArgs_Fields(t *testing.T) {
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

func TestSearchIndexArgs_Fields(t *testing.T) {
	t.Run("specific series", func(t *testing.T) {
		seriesID := uuid.Must(uuid.NewV7())
		args := SearchIndexArgs{
			SeriesID:    &seriesID,
			FullReindex: false,
		}
		assert.NotNil(t, args.SeriesID)
		assert.False(t, args.FullReindex)
	})

	t.Run("full reindex", func(t *testing.T) {
		args := SearchIndexArgs{
			FullReindex: true,
		}
		assert.Nil(t, args.SeriesID)
		assert.True(t, args.FullReindex)
	})
}

func TestSeriesRefreshArgs_Fields(t *testing.T) {
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
}
