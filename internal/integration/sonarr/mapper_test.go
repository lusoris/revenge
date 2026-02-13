package sonarr

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// --- NewMapper ---

func TestNewMapper(t *testing.T) {
	m := NewMapper()
	assert.NotNil(t, m)
}

// --- ToSeries ---

func TestMapper_ToSeries(t *testing.T) {
	mapper := NewMapper()

	firstAired := time.Date(2008, 1, 20, 0, 0, 0, 0, time.UTC)
	lastAired := time.Date(2013, 9, 29, 0, 0, 0, 0, time.UTC)

	ss := &Series{
		ID:    456,
		Title: "Breaking Bad",
		OriginalLanguage: Language{
			ID:   1,
			Name: "English",
		},
		Overview:   "A chemistry teacher diagnosed with terminal lung cancer...",
		Status:     "ended",
		FirstAired: &firstAired,
		LastAired:  &lastAired,
		TVDbID:     81189,
		IMDbID:     "tt0903747",
		SeriesType: SeriesTypeStandard,
		Seasons: []SeasonInfo{
			{SeasonNumber: 0},
			{SeasonNumber: 1},
			{SeasonNumber: 2},
			{SeasonNumber: 3},
			{SeasonNumber: 4},
			{SeasonNumber: 5},
		},
		Images: []Image{
			{CoverType: "poster", RemoteURL: "https://artworks.thetvdb.com/poster.jpg"},
			{CoverType: "fanart", RemoteURL: "https://artworks.thetvdb.com/fanart.jpg"},
		},
		Ratings: Ratings{
			Value: 9.5,
			Votes: 50000,
		},
		Statistics: &Statistics{
			TotalEpisodeCount: 62,
		},
	}

	series := mapper.ToSeries(ss)

	assert.NotEqual(t, uuid.Nil, series.ID)
	assert.Equal(t, "Breaking Bad", series.Title)
	assert.Equal(t, "English", series.OriginalLanguage)
	assert.Equal(t, int32(456), *series.SonarrID)
	assert.Equal(t, int32(81189), *series.TVDbID)
	assert.Equal(t, "tt0903747", *series.IMDbID)
	assert.Equal(t, "ended", *series.Status)
	assert.Equal(t, "A chemistry teacher diagnosed with terminal lung cancer...", *series.Overview)
	assert.Equal(t, &firstAired, series.FirstAirDate)
	assert.Equal(t, &lastAired, series.LastAirDate)
	assert.Equal(t, int32(6), series.TotalSeasons)
	assert.Equal(t, int32(62), series.TotalEpisodes)

	// Ratings
	assert.NotNil(t, series.VoteAverage)
	f, ok := series.VoteAverage.Float64()
	assert.True(t, ok)
	assert.Equal(t, 9.5, f)
	assert.Equal(t, int32(50000), *series.VoteCount)

	// Images
	assert.Equal(t, "https://artworks.thetvdb.com/poster.jpg", *series.PosterPath)
	assert.Equal(t, "https://artworks.thetvdb.com/fanart.jpg", *series.BackdropPath)

	// Series type
	assert.Equal(t, "Scripted", *series.Type)

	// Timestamps
	assert.False(t, series.CreatedAt.IsZero())
	assert.False(t, series.UpdatedAt.IsZero())
}

func TestMapper_ToSeries_ImageFallbackToLocalURL(t *testing.T) {
	mapper := NewMapper()

	ss := &Series{
		ID:    1,
		Title: "Test Series",
		Images: []Image{
			{CoverType: "poster", URL: "/MediaCover/1/poster.jpg", RemoteURL: ""},
			{CoverType: "fanart", URL: "/MediaCover/1/fanart.jpg", RemoteURL: ""},
		},
	}

	series := mapper.ToSeries(ss)

	assert.Equal(t, "/MediaCover/1/poster.jpg", *series.PosterPath)
	assert.Equal(t, "/MediaCover/1/fanart.jpg", *series.BackdropPath)
}

func TestMapper_ToSeries_ImagePrefersRemoteURL(t *testing.T) {
	mapper := NewMapper()

	ss := &Series{
		ID:    1,
		Title: "Test Series",
		Images: []Image{
			{CoverType: "poster", URL: "/local/poster.jpg", RemoteURL: "https://remote/poster.jpg"},
			{CoverType: "fanart", URL: "/local/fanart.jpg", RemoteURL: "https://remote/fanart.jpg"},
		},
	}

	series := mapper.ToSeries(ss)

	assert.Equal(t, "https://remote/poster.jpg", *series.PosterPath)
	assert.Equal(t, "https://remote/fanart.jpg", *series.BackdropPath)
}

func TestMapper_ToSeries_NoImages(t *testing.T) {
	mapper := NewMapper()

	ss := &Series{
		ID:    1,
		Title: "Test Series",
	}

	series := mapper.ToSeries(ss)

	assert.Nil(t, series.PosterPath)
	assert.Nil(t, series.BackdropPath)
}

func TestMapper_ToSeries_BannerImageIgnored(t *testing.T) {
	mapper := NewMapper()

	ss := &Series{
		ID:    1,
		Title: "Test Series",
		Images: []Image{
			{CoverType: "banner", RemoteURL: "https://example.com/banner.jpg"},
		},
	}

	series := mapper.ToSeries(ss)

	// Banner is not mapped to poster or backdrop
	assert.Nil(t, series.PosterPath)
	assert.Nil(t, series.BackdropPath)
}

func TestMapper_ToSeries_NoRatings(t *testing.T) {
	mapper := NewMapper()

	ss := &Series{
		ID:    1,
		Title: "Test Series",
		Ratings: Ratings{
			Value: 0,
			Votes: 0,
		},
	}

	series := mapper.ToSeries(ss)

	assert.Nil(t, series.VoteAverage)
	assert.Nil(t, series.VoteCount)
}

func TestMapper_ToSeries_NoStatistics(t *testing.T) {
	mapper := NewMapper()

	ss := &Series{
		ID:    1,
		Title: "Test Series",
	}

	series := mapper.ToSeries(ss)

	assert.Equal(t, int32(0), series.TotalEpisodes)
}

func TestMapper_ToSeries_EmptyIMDbID(t *testing.T) {
	mapper := NewMapper()

	ss := &Series{
		ID:     1,
		Title:  "Test Series",
		IMDbID: "",
	}

	series := mapper.ToSeries(ss)

	// ptrString returns nil for empty string
	assert.Nil(t, series.IMDbID)
}

func TestMapper_ToSeries_EmptyOverview(t *testing.T) {
	mapper := NewMapper()

	ss := &Series{
		ID:    1,
		Title: "Test Series",
	}

	series := mapper.ToSeries(ss)

	assert.Nil(t, series.Overview)
}

func TestMapper_ToSeries_EmptySeriesType(t *testing.T) {
	mapper := NewMapper()

	ss := &Series{
		ID:         1,
		Title:      "Test Series",
		SeriesType: "",
	}

	series := mapper.ToSeries(ss)

	assert.Nil(t, series.Type)
}

func TestMapper_ToSeries_SeriesTypes(t *testing.T) {
	mapper := NewMapper()

	tests := []struct {
		name       string
		seriesType string
		expected   string
	}{
		{"standard maps to Scripted", SeriesTypeStandard, "Scripted"},
		{"daily maps to Talk Show", SeriesTypeDaily, "Talk Show"},
		{"anime maps to Animation", SeriesTypeAnime, "Animation"},
		{"unknown maps to Scripted", "unknowntype", "Scripted"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := &Series{
				ID:         1,
				Title:      "Test",
				SeriesType: tt.seriesType,
			}
			series := mapper.ToSeries(ss)
			assert.Equal(t, tt.expected, *series.Type)
		})
	}
}

// --- ToSeason ---

func TestMapper_ToSeason(t *testing.T) {
	mapper := NewMapper()
	seriesID := uuid.Must(uuid.NewV7())

	si := &SeasonInfo{
		SeasonNumber: 3,
		Monitored:    true,
		Statistics: &Statistics{
			EpisodeCount: 13,
		},
	}

	season := mapper.ToSeason(si, seriesID)

	assert.NotEqual(t, uuid.Nil, season.ID)
	assert.Equal(t, seriesID, season.SeriesID)
	assert.Equal(t, int32(3), season.SeasonNumber)
	assert.Equal(t, "Season 3", season.Name)
	assert.Equal(t, int32(13), season.EpisodeCount)
	assert.False(t, season.CreatedAt.IsZero())
	assert.False(t, season.UpdatedAt.IsZero())
}

func TestMapper_ToSeason_Specials(t *testing.T) {
	mapper := NewMapper()
	seriesID := uuid.Must(uuid.NewV7())

	si := &SeasonInfo{
		SeasonNumber: 0,
		Statistics: &Statistics{
			EpisodeCount: 5,
		},
	}

	season := mapper.ToSeason(si, seriesID)

	assert.Equal(t, int32(0), season.SeasonNumber)
	assert.Equal(t, "Specials", season.Name)
	assert.Equal(t, int32(5), season.EpisodeCount)
}

func TestMapper_ToSeason_NoStatistics(t *testing.T) {
	mapper := NewMapper()
	seriesID := uuid.Must(uuid.NewV7())

	si := &SeasonInfo{
		SeasonNumber: 1,
	}

	season := mapper.ToSeason(si, seriesID)

	assert.Equal(t, int32(0), season.EpisodeCount)
}

func TestMapper_ToSeason_VariousNumbers(t *testing.T) {
	mapper := NewMapper()
	seriesID := uuid.Must(uuid.NewV7())

	tests := []struct {
		seasonNumber int
		expectedName string
	}{
		{0, "Specials"},
		{1, "Season 1"},
		{10, "Season 10"},
		{25, "Season 25"},
	}

	for _, tt := range tests {
		t.Run(tt.expectedName, func(t *testing.T) {
			si := &SeasonInfo{SeasonNumber: tt.seasonNumber}
			season := mapper.ToSeason(si, seriesID)
			assert.Equal(t, tt.expectedName, season.Name)
			assert.Equal(t, int32(tt.seasonNumber), season.SeasonNumber)
		})
	}
}

// --- ToEpisode ---

func TestMapper_ToEpisode(t *testing.T) {
	mapper := NewMapper()
	seriesID := uuid.Must(uuid.NewV7())
	seasonID := uuid.Must(uuid.NewV7())

	airDate := time.Date(2013, 9, 29, 21, 0, 0, 0, time.UTC)

	se := &Episode{
		ID:            501,
		TVDbID:        4752975,
		SeasonNumber:  5,
		EpisodeNumber: 16,
		Title:         "Felina",
		Overview:      "All bad things must come to an end.",
		AirDateUtc:    &airDate,
		Runtime:       55,
		Images: []Image{
			{CoverType: "screenshot", RemoteURL: "https://artworks.thetvdb.com/screenshot.jpg"},
		},
	}

	ep := mapper.ToEpisode(se, seriesID, seasonID)

	assert.NotEqual(t, uuid.Nil, ep.ID)
	assert.Equal(t, seriesID, ep.SeriesID)
	assert.Equal(t, seasonID, ep.SeasonID)
	assert.Equal(t, int32(4752975), *ep.TVDbID)
	assert.Equal(t, int32(5), ep.SeasonNumber)
	assert.Equal(t, int32(16), ep.EpisodeNumber)
	assert.Equal(t, "Felina", ep.Title)
	assert.Equal(t, "All bad things must come to an end.", *ep.Overview)
	assert.Equal(t, &airDate, ep.AirDate)
	assert.Equal(t, int32(55), *ep.Runtime)
	assert.Equal(t, "https://artworks.thetvdb.com/screenshot.jpg", *ep.StillPath)
	assert.False(t, ep.CreatedAt.IsZero())
	assert.False(t, ep.UpdatedAt.IsZero())
}

func TestMapper_ToEpisode_ScreenshotFallbackToLocalURL(t *testing.T) {
	mapper := NewMapper()
	seriesID := uuid.Must(uuid.NewV7())
	seasonID := uuid.Must(uuid.NewV7())

	se := &Episode{
		ID:            1,
		SeasonNumber:  1,
		EpisodeNumber: 1,
		Title:         "Pilot",
		Images: []Image{
			{CoverType: "screenshot", URL: "/MediaCover/screenshot.jpg", RemoteURL: ""},
		},
	}

	ep := mapper.ToEpisode(se, seriesID, seasonID)

	assert.Equal(t, "/MediaCover/screenshot.jpg", *ep.StillPath)
}

func TestMapper_ToEpisode_ScreenshotPrefersRemoteURL(t *testing.T) {
	mapper := NewMapper()
	seriesID := uuid.Must(uuid.NewV7())
	seasonID := uuid.Must(uuid.NewV7())

	se := &Episode{
		ID:            1,
		SeasonNumber:  1,
		EpisodeNumber: 1,
		Title:         "Pilot",
		Images: []Image{
			{CoverType: "screenshot", URL: "/local/screenshot.jpg", RemoteURL: "https://remote/screenshot.jpg"},
		},
	}

	ep := mapper.ToEpisode(se, seriesID, seasonID)

	assert.Equal(t, "https://remote/screenshot.jpg", *ep.StillPath)
}

func TestMapper_ToEpisode_NoScreenshotImage(t *testing.T) {
	mapper := NewMapper()
	seriesID := uuid.Must(uuid.NewV7())
	seasonID := uuid.Must(uuid.NewV7())

	se := &Episode{
		ID:            1,
		SeasonNumber:  1,
		EpisodeNumber: 1,
		Title:         "Pilot",
		Images: []Image{
			{CoverType: "poster", RemoteURL: "https://example.com/poster.jpg"},
		},
	}

	ep := mapper.ToEpisode(se, seriesID, seasonID)

	assert.Nil(t, ep.StillPath)
}

func TestMapper_ToEpisode_NoImages(t *testing.T) {
	mapper := NewMapper()
	seriesID := uuid.Must(uuid.NewV7())
	seasonID := uuid.Must(uuid.NewV7())

	se := &Episode{
		ID:            1,
		SeasonNumber:  1,
		EpisodeNumber: 1,
		Title:         "Pilot",
	}

	ep := mapper.ToEpisode(se, seriesID, seasonID)

	assert.Nil(t, ep.StillPath)
}

func TestMapper_ToEpisode_EmptyOverview(t *testing.T) {
	mapper := NewMapper()
	seriesID := uuid.Must(uuid.NewV7())
	seasonID := uuid.Must(uuid.NewV7())

	se := &Episode{
		ID:            1,
		SeasonNumber:  1,
		EpisodeNumber: 1,
		Title:         "Pilot",
		Overview:      "",
	}

	ep := mapper.ToEpisode(se, seriesID, seasonID)

	assert.Nil(t, ep.Overview)
}

func TestMapper_ToEpisode_OnlyFirstScreenshotUsed(t *testing.T) {
	mapper := NewMapper()
	seriesID := uuid.Must(uuid.NewV7())
	seasonID := uuid.Must(uuid.NewV7())

	se := &Episode{
		ID:            1,
		SeasonNumber:  1,
		EpisodeNumber: 1,
		Title:         "Pilot",
		Images: []Image{
			{CoverType: "screenshot", RemoteURL: "https://first-screenshot.jpg"},
			{CoverType: "screenshot", RemoteURL: "https://second-screenshot.jpg"},
		},
	}

	ep := mapper.ToEpisode(se, seriesID, seasonID)

	// Only first screenshot is used (break after first match)
	assert.Equal(t, "https://first-screenshot.jpg", *ep.StillPath)
}

// --- ToEpisodeFile ---

func TestMapper_ToEpisodeFile(t *testing.T) {
	mapper := NewMapper()
	episodeID := uuid.Must(uuid.NewV7())

	dateAdded := time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC)

	sef := &EpisodeFile{
		ID:           42,
		RelativePath: "Season 05/Breaking.Bad.S05E16.Felina.mkv",
		Path:         "/tv/Breaking Bad/Season 05/Breaking.Bad.S05E16.Felina.mkv",
		Size:         2500000000,
		DateAdded:    dateAdded,
		Quality: Quality{
			Quality: QualityInfo{
				Name:       "Bluray-1080p",
				Resolution: 1080,
			},
		},
		MediaInfo: &MediaInfo{
			VideoCodec:     "x265",
			AudioCodec:     "DTS-HD MA",
			VideoBitrate:   15000000,
			RunTime:        "55m30s",
			AudioLanguages: "English",
			Subtitles:      "English / German",
		},
	}

	file := mapper.ToEpisodeFile(sef, episodeID)

	assert.NotEqual(t, uuid.Nil, file.ID)
	assert.Equal(t, episodeID, file.EpisodeID)
	assert.Equal(t, "/tv/Breaking Bad/Season 05/Breaking.Bad.S05E16.Felina.mkv", file.FilePath)
	assert.Equal(t, "Season 05/Breaking.Bad.S05E16.Felina.mkv", file.FileName)
	assert.Equal(t, int64(2500000000), file.FileSize)
	assert.Equal(t, int32(42), *file.SonarrFileID)
	assert.Equal(t, dateAdded, file.CreatedAt)
	assert.False(t, file.UpdatedAt.IsZero())

	// Quality
	assert.Equal(t, "Bluray-1080p", *file.QualityProfile)
	assert.Equal(t, "1080p", *file.Resolution)

	// Container
	assert.Equal(t, "mkv", *file.Container)

	// Media info
	assert.Equal(t, "x265", *file.VideoCodec)
	assert.Equal(t, "DTS-HD MA", *file.AudioCodec)
	assert.Equal(t, int32(15000), *file.BitrateKbps)

	// Duration
	assert.NotNil(t, file.DurationSeconds)
	d, ok := file.DurationSeconds.Float64()
	assert.True(t, ok)
	assert.InDelta(t, 3330.0, d, 0.01) // 55m30s = 3330 seconds

	// Audio languages
	assert.Equal(t, []string{"English"}, file.AudioLanguages)

	// Subtitles
	assert.Equal(t, []string{"English / German"}, file.SubtitleLanguages)
}

func TestMapper_ToEpisodeFile_NilMediaInfo(t *testing.T) {
	mapper := NewMapper()
	episodeID := uuid.Must(uuid.NewV7())

	sef := &EpisodeFile{
		ID:        1,
		Path:      "/tv/show/episode.mp4",
		Size:      1000000,
		DateAdded: time.Now(),
		Quality: Quality{
			Quality: QualityInfo{
				Name:       "WEBDL-720p",
				Resolution: 720,
			},
		},
	}

	file := mapper.ToEpisodeFile(sef, episodeID)

	assert.Nil(t, file.VideoCodec)
	assert.Nil(t, file.AudioCodec)
	assert.Nil(t, file.BitrateKbps)
	assert.Nil(t, file.DurationSeconds)
	assert.Nil(t, file.AudioLanguages)
	assert.Nil(t, file.SubtitleLanguages)

	// Quality should still be set
	assert.Equal(t, "WEBDL-720p", *file.QualityProfile)
	assert.Equal(t, "720p", *file.Resolution)
}

func TestMapper_ToEpisodeFile_EmptyQualityName(t *testing.T) {
	mapper := NewMapper()
	episodeID := uuid.Must(uuid.NewV7())

	sef := &EpisodeFile{
		ID:        1,
		Path:      "/tv/show/episode.mkv",
		Size:      1000000,
		DateAdded: time.Now(),
		Quality: Quality{
			Quality: QualityInfo{
				Name:       "",
				Resolution: 0,
			},
		},
	}

	file := mapper.ToEpisodeFile(sef, episodeID)

	assert.Nil(t, file.QualityProfile)
	assert.Nil(t, file.Resolution)
}

func TestMapper_ToEpisodeFile_MediaInfoEmptyStrings(t *testing.T) {
	mapper := NewMapper()
	episodeID := uuid.Must(uuid.NewV7())

	sef := &EpisodeFile{
		ID:        1,
		Path:      "/tv/show/episode.avi",
		Size:      500000,
		DateAdded: time.Now(),
		MediaInfo: &MediaInfo{
			VideoCodec:     "",
			AudioCodec:     "",
			VideoBitrate:   0,
			RunTime:        "",
			AudioLanguages: "",
			Subtitles:      "",
		},
	}

	file := mapper.ToEpisodeFile(sef, episodeID)

	// ptrString returns nil for empty strings
	assert.Nil(t, file.VideoCodec)
	assert.Nil(t, file.AudioCodec)
	assert.Equal(t, int32(0), *file.BitrateKbps) // 0/1000 = 0, still gets ptr
	assert.Nil(t, file.DurationSeconds)          // empty RunTime is not parsed
	assert.Nil(t, file.AudioLanguages)
	assert.Nil(t, file.SubtitleLanguages)
}

func TestMapper_ToEpisodeFile_InvalidRunTime(t *testing.T) {
	mapper := NewMapper()
	episodeID := uuid.Must(uuid.NewV7())

	sef := &EpisodeFile{
		ID:        1,
		Path:      "/tv/show/episode.mkv",
		Size:      1000000,
		DateAdded: time.Now(),
		MediaInfo: &MediaInfo{
			RunTime: "not-a-duration",
		},
	}

	file := mapper.ToEpisodeFile(sef, episodeID)

	// Invalid duration should result in nil DurationSeconds
	assert.Nil(t, file.DurationSeconds)
}

func TestMapper_ToEpisodeFile_ContainerExtraction(t *testing.T) {
	mapper := NewMapper()
	episodeID := uuid.Must(uuid.NewV7())

	tests := []struct {
		name              string
		path              string
		expectedContainer *string
	}{
		{"mkv file", "/tv/show/episode.mkv", new("mkv")},
		{"mp4 file", "/tv/show/episode.mp4", new("mp4")},
		{"avi file", "/tv/show/episode.avi", new("avi")},
		{"ts file", "/tv/show/episode.ts", new("ts")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sef := &EpisodeFile{
				ID:        1,
				Path:      tt.path,
				DateAdded: time.Now(),
			}
			file := mapper.ToEpisodeFile(sef, episodeID)
			if tt.expectedContainer != nil {
				assert.Equal(t, *tt.expectedContainer, *file.Container)
			} else {
				assert.Nil(t, file.Container)
			}
		})
	}
}

// --- ToGenres ---

func TestMapper_ToGenres(t *testing.T) {
	mapper := NewMapper()
	seriesID := uuid.Must(uuid.NewV7())

	ss := &Series{
		ID:     1,
		Title:  "Test Series",
		Genres: []string{"Drama", "Thriller", "Crime"},
	}

	genres := mapper.ToGenres(ss, seriesID)

	assert.Len(t, genres, 3)
	assert.Equal(t, "Drama", genres[0].Name)
	assert.Equal(t, "Thriller", genres[1].Name)
	assert.Equal(t, "Crime", genres[2].Name)

	for _, g := range genres {
		assert.Equal(t, seriesID, g.SeriesID)
		assert.NotEqual(t, uuid.Nil, g.ID)
		assert.False(t, g.CreatedAt.IsZero())
	}
}

func TestMapper_ToGenres_Empty(t *testing.T) {
	mapper := NewMapper()
	seriesID := uuid.Must(uuid.NewV7())

	ss := &Series{
		ID:     1,
		Title:  "Test Series",
		Genres: []string{},
	}

	genres := mapper.ToGenres(ss, seriesID)

	assert.Empty(t, genres)
	assert.Len(t, genres, 0)
}

func TestMapper_ToGenres_Single(t *testing.T) {
	mapper := NewMapper()
	seriesID := uuid.Must(uuid.NewV7())

	ss := &Series{
		ID:     1,
		Title:  "Test Series",
		Genres: []string{"Animation"},
	}

	genres := mapper.ToGenres(ss, seriesID)

	assert.Len(t, genres, 1)
	assert.Equal(t, "Animation", genres[0].Name)
}

// --- Helper functions ---

func TestPtr(t *testing.T) {
	intVal := int32(42)
	result := new(intVal)
	assert.Equal(t, intVal, *result)

	strVal := "hello"
	strResult := new(strVal)
	assert.Equal(t, strVal, *strResult)
}

func TestPtrString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *string
	}{
		{"non-empty string", "hello", new("hello")},
		{"empty string", "", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ptrString(tt.input)
			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.Equal(t, *tt.expected, *result)
			}
		})
	}
}

func TestResolutionToString(t *testing.T) {
	tests := []struct {
		name       string
		resolution int
		expected   string
	}{
		{"4K (2160)", 2160, "4K"},
		{"above 4K (4320)", 4320, "4K"},
		{"1080p", 1080, "1080p"},
		{"above 1080 (1440)", 1440, "1080p"},
		{"720p", 720, "720p"},
		{"above 720 (900)", 900, "720p"},
		{"480p", 480, "480p"},
		{"above 480 (576)", 576, "480p"},
		{"SD (360)", 360, "SD"},
		{"SD (240)", 240, "SD"},
		{"SD (0)", 0, "SD"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, resolutionToString(tt.resolution))
		})
	}
}

func TestMapSeriesType(t *testing.T) {
	tests := []struct {
		name       string
		sonarrType string
		expected   string
	}{
		{"standard", SeriesTypeStandard, "Scripted"},
		{"daily", SeriesTypeDaily, "Talk Show"},
		{"anime", SeriesTypeAnime, "Animation"},
		{"unknown defaults to Scripted", "documentary", "Scripted"},
		{"empty defaults to Scripted", "", "Scripted"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, mapSeriesType(tt.sonarrType))
		})
	}
}

func TestExtractContainer(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected *string
	}{
		{"mkv extension", "/path/to/video.mkv", new("mkv")},
		{"mp4 extension", "/path/to/video.mp4", new("mp4")},
		{"avi extension", "/path/to/video.avi", new("avi")},
		{"ts extension", "/path/to/video.ts", new("ts")},
		{"dots in path", "/path/to/my.show.S01E01.mkv", new("mkv")},
		{"short path (exactly 4 chars)", "a.mk", new("mk")},
		{"too short path (3 chars)", "a.m", nil},
		{"too short path (2 chars)", "ab", nil},
		{"empty path", "", nil},
		{"no extension", "/path/to/noextension", nil},
		{"path with trailing dot", "/path/to/file.", new("")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractContainer(tt.path)
			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, *tt.expected, *result)
			}
		})
	}
}

// --- Table-driven integration test covering all branches ---

func TestMapper_ToSeries_AllBranches(t *testing.T) {
	mapper := NewMapper()

	firstAired := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		series        *Series
		checkPoster   *string
		checkBackdrop *string
		checkType     *string
		checkVoteAvg  bool
		checkEpisodes int32
		checkIMDb     *string
		checkOverview *string
	}{
		{
			name: "full series with all fields",
			series: &Series{
				ID:               100,
				Title:            "Full Series",
				OriginalLanguage: Language{Name: "Japanese"},
				Overview:         "A great show",
				Status:           "continuing",
				FirstAired:       &firstAired,
				TVDbID:           12345,
				IMDbID:           "tt1234567",
				SeriesType:       SeriesTypeAnime,
				Seasons:          []SeasonInfo{{SeasonNumber: 1}},
				Images: []Image{
					{CoverType: "poster", RemoteURL: "https://poster.jpg"},
					{CoverType: "fanart", RemoteURL: "https://fanart.jpg"},
				},
				Ratings:    Ratings{Value: 8.5, Votes: 1000},
				Statistics: &Statistics{TotalEpisodeCount: 24},
			},
			checkPoster:   new("https://poster.jpg"),
			checkBackdrop: new("https://fanart.jpg"),
			checkType:     new("Animation"),
			checkVoteAvg:  true,
			checkEpisodes: 24,
			checkIMDb:     new("tt1234567"),
			checkOverview: new("A great show"),
		},
		{
			name: "minimal series",
			series: &Series{
				ID:    1,
				Title: "Minimal",
			},
			checkPoster:   nil,
			checkBackdrop: nil,
			checkType:     nil,
			checkVoteAvg:  false,
			checkEpisodes: 0,
			checkIMDb:     nil,
			checkOverview: nil,
		},
		{
			name: "series with local image URLs only",
			series: &Series{
				ID:    2,
				Title: "Local Images",
				Images: []Image{
					{CoverType: "poster", URL: "/local/poster.jpg"},
					{CoverType: "fanart", URL: "/local/fanart.jpg"},
				},
			},
			checkPoster:   new("/local/poster.jpg"),
			checkBackdrop: new("/local/fanart.jpg"),
			checkType:     nil,
			checkVoteAvg:  false,
			checkEpisodes: 0,
		},
		{
			name: "series with daily type",
			series: &Series{
				ID:         3,
				Title:      "Daily Show",
				SeriesType: SeriesTypeDaily,
			},
			checkType: new("Talk Show"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.ToSeries(tt.series)

			assert.NotEqual(t, uuid.Nil, result.ID)
			assert.Equal(t, tt.series.Title, result.Title)

			if tt.checkPoster != nil {
				assert.Equal(t, *tt.checkPoster, *result.PosterPath)
			} else if tt.checkPoster == nil && tt.series.Images == nil {
				assert.Nil(t, result.PosterPath)
			}

			if tt.checkBackdrop != nil {
				assert.Equal(t, *tt.checkBackdrop, *result.BackdropPath)
			} else if tt.checkBackdrop == nil && tt.series.Images == nil {
				assert.Nil(t, result.BackdropPath)
			}

			if tt.checkType != nil {
				assert.Equal(t, *tt.checkType, *result.Type)
			} else if tt.series.SeriesType == "" {
				assert.Nil(t, result.Type)
			}

			if tt.checkVoteAvg {
				assert.NotNil(t, result.VoteAverage)
			} else if tt.series.Ratings.Value == 0 {
				assert.Nil(t, result.VoteAverage)
			}

			assert.Equal(t, tt.checkEpisodes, result.TotalEpisodes)

			if tt.checkIMDb != nil {
				assert.Equal(t, *tt.checkIMDb, *result.IMDbID)
			}

			if tt.checkOverview != nil {
				assert.Equal(t, *tt.checkOverview, *result.Overview)
			}
		})
	}
}

func TestMapper_ToEpisodeFile_Resolutions(t *testing.T) {
	mapper := NewMapper()
	episodeID := uuid.Must(uuid.NewV7())

	tests := []struct {
		name       string
		resolution int
		expected   string
	}{
		{"4K", 2160, "4K"},
		{"1080p", 1080, "1080p"},
		{"720p", 720, "720p"},
		{"480p", 480, "480p"},
		{"SD", 360, "SD"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sef := &EpisodeFile{
				ID:        1,
				Path:      "/tv/show/episode.mkv",
				DateAdded: time.Now(),
				Quality: Quality{
					Quality: QualityInfo{
						Name:       "TestQuality",
						Resolution: tt.resolution,
					},
				},
			}
			file := mapper.ToEpisodeFile(sef, episodeID)
			assert.Equal(t, tt.expected, *file.Resolution)
		})
	}
}

func TestMapper_ToEpisodeFile_ZeroResolution(t *testing.T) {
	mapper := NewMapper()
	episodeID := uuid.Must(uuid.NewV7())

	sef := &EpisodeFile{
		ID:        1,
		Path:      "/tv/show/episode.mkv",
		DateAdded: time.Now(),
		Quality: Quality{
			Quality: QualityInfo{
				Name:       "Unknown",
				Resolution: 0,
			},
		},
	}

	file := mapper.ToEpisodeFile(sef, episodeID)

	// Resolution 0 is not > 0, so Resolution should be nil
	assert.Nil(t, file.Resolution)
}

// strPtr is a test helper to create string pointers.
//
//go:fix inline
func strPtr(s string) *string {
	return new(s)
}
