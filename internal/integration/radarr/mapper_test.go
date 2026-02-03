package radarr

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMapper_ToMovie(t *testing.T) {
	mapper := NewMapper()

	radarrMovie := &Movie{
		ID:            1,
		Title:         "Inception",
		OriginalTitle: "Inception",
		Year:          2010,
		TMDbID:        27205,
		IMDbID:        "tt1375666",
		Overview:      "A thief who steals corporate secrets...",
		Runtime:       148,
		Status:        "released",
		Genres:        []string{"Action", "Science Fiction", "Thriller"},
		Popularity:    91.123,
		Images: []Image{
			{CoverType: "poster", RemoteURL: "https://image.tmdb.org/t/p/poster.jpg"},
			{CoverType: "fanart", RemoteURL: "https://image.tmdb.org/t/p/backdrop.jpg"},
		},
		YouTubeTrailerID: "YoHD9XEInc0",
		Ratings: Ratings{
			TMDb: &Rating{
				Value: 8.4,
				Votes: 35000,
			},
		},
		Added: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	movie := mapper.ToMovie(radarrMovie)

	assert.NotEqual(t, uuid.Nil, movie.ID)
	assert.Equal(t, "Inception", movie.Title)
	assert.Equal(t, int32(27205), *movie.TMDbID)
	assert.Equal(t, "tt1375666", *movie.IMDbID)
	assert.Equal(t, int32(2010), *movie.Year)
	assert.Equal(t, int32(148), *movie.Runtime)
	assert.Equal(t, "A thief who steals corporate secrets...", *movie.Overview)
	assert.Equal(t, "released", *movie.Status)
	assert.Equal(t, int32(1), *movie.RadarrID)

	// Check ratings
	assert.NotNil(t, movie.VoteAverage)
	assert.Equal(t, 8.4, movie.VoteAverage.InexactFloat64())
	assert.Equal(t, int32(35000), *movie.VoteCount)

	// Check images
	assert.Equal(t, "https://image.tmdb.org/t/p/poster.jpg", *movie.PosterPath)
	assert.Equal(t, "https://image.tmdb.org/t/p/backdrop.jpg", *movie.BackdropPath)

	// Check trailer
	assert.Equal(t, "https://www.youtube.com/watch?v=YoHD9XEInc0", *movie.TrailerURL)
}

func TestMapper_ToMovie_WithIMDbRating(t *testing.T) {
	mapper := NewMapper()

	radarrMovie := &Movie{
		ID:     1,
		Title:  "Test Movie",
		Year:   2024,
		TMDbID: 12345,
		Ratings: Ratings{
			IMDb: &Rating{
				Value: 7.5,
				Votes: 10000,
			},
		},
	}

	movie := mapper.ToMovie(radarrMovie)

	assert.NotNil(t, movie.VoteAverage)
	assert.Equal(t, 7.5, movie.VoteAverage.InexactFloat64())
	assert.Equal(t, int32(10000), *movie.VoteCount)
}

func TestMapper_ToMovie_WithDates(t *testing.T) {
	mapper := NewMapper()

	digitalRelease := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	physicalRelease := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)
	inCinemas := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		movie    *Movie
		expected time.Time
	}{
		{
			name: "prefers digital release",
			movie: &Movie{
				ID:              1,
				Title:           "Test",
				Year:            2024,
				TMDbID:          12345,
				DigitalRelease:  &digitalRelease,
				PhysicalRelease: &physicalRelease,
				InCinemas:       &inCinemas,
			},
			expected: digitalRelease,
		},
		{
			name: "falls back to physical release",
			movie: &Movie{
				ID:              1,
				Title:           "Test",
				Year:            2024,
				TMDbID:          12345,
				PhysicalRelease: &physicalRelease,
				InCinemas:       &inCinemas,
			},
			expected: physicalRelease,
		},
		{
			name: "falls back to in cinemas",
			movie: &Movie{
				ID:        1,
				Title:     "Test",
				Year:      2024,
				TMDbID:    12345,
				InCinemas: &inCinemas,
			},
			expected: inCinemas,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			movie := mapper.ToMovie(tt.movie)
			assert.NotNil(t, movie.ReleaseDate)
			assert.Equal(t, tt.expected, *movie.ReleaseDate)
		})
	}
}

func TestMapper_ToMovieFile(t *testing.T) {
	mapper := NewMapper()
	movieID := uuid.New()

	radarrFile := &MovieFile{
		ID:           1,
		MovieID:      123,
		RelativePath: "Inception (2010)/Inception.mkv",
		Path:         "/movies/Inception (2010)/Inception.mkv",
		Size:         4500000000,
		DateAdded:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Quality: Quality{
			Quality: QualityInfo{
				ID:         7,
				Name:       "Bluray-1080p",
				Resolution: 1080,
			},
		},
		ReleaseGroup: "SPARKS",
		MediaInfo: &MediaInfo{
			VideoCodec:        "x265",
			AudioCodec:        "DTS-HD MA",
			AudioChannels:     5.1,
			VideoDynamicRange: "HDR",
			VideoFps:          23.976,
			VideoBitrate:      25000000,
			AudioLanguages:    "English",
			Subtitles:         "English",
		},
	}

	file := mapper.ToMovieFile(radarrFile, movieID)

	assert.NotEqual(t, uuid.Nil, file.ID)
	assert.Equal(t, movieID, file.MovieID)
	assert.Equal(t, "/movies/Inception (2010)/Inception.mkv", file.FilePath)
	assert.Equal(t, "Inception (2010)/Inception.mkv", file.FileName)
	assert.Equal(t, int64(4500000000), file.FileSize)
	assert.Equal(t, int32(1), *file.RadarrFileID)

	// Quality
	assert.Contains(t, *file.QualityProfile, "Bluray-1080p")
	assert.Contains(t, *file.QualityProfile, "SPARKS")
	assert.Equal(t, "1080p", *file.Resolution)

	// Media info
	assert.Equal(t, "x265", *file.VideoCodec)
	assert.Equal(t, "DTS-HD MA", *file.AudioCodec)
	assert.Equal(t, "HDR", *file.DynamicRange)
	assert.Equal(t, int32(25000), *file.BitrateKbps)
	assert.Equal(t, "5.1", *file.AudioChannels)
}

func TestMapper_ToMovieCollection(t *testing.T) {
	mapper := NewMapper()

	radarrCollection := &Collection{
		Name:   "The Dark Knight Collection",
		TMDbID: 263,
		Images: []Image{
			{CoverType: "poster", RemoteURL: "https://image.tmdb.org/t/p/poster.jpg"},
			{CoverType: "fanart", RemoteURL: "https://image.tmdb.org/t/p/backdrop.jpg"},
		},
	}

	collection := mapper.ToMovieCollection(radarrCollection)

	assert.NotEqual(t, uuid.Nil, collection.ID)
	assert.Equal(t, "The Dark Knight Collection", collection.Name)
	assert.Equal(t, int32(263), *collection.TMDbCollectionID)
	assert.Equal(t, "https://image.tmdb.org/t/p/poster.jpg", *collection.PosterPath)
	assert.Equal(t, "https://image.tmdb.org/t/p/backdrop.jpg", *collection.BackdropPath)
}

func TestMapper_ToMovieCollection_Nil(t *testing.T) {
	mapper := NewMapper()
	assert.Nil(t, mapper.ToMovieCollection(nil))
}

func TestMapper_ToGenres(t *testing.T) {
	mapper := NewMapper()
	movieID := uuid.New()

	radarrMovie := &Movie{
		ID:     1,
		Title:  "Test Movie",
		Genres: []string{"Action", "Science Fiction", "Thriller"},
	}

	genres := mapper.ToGenres(radarrMovie, movieID)

	assert.Len(t, genres, 3)
	assert.Equal(t, "Action", genres[0].Name)
	assert.Equal(t, "Science Fiction", genres[1].Name)
	assert.Equal(t, "Thriller", genres[2].Name)

	for _, g := range genres {
		assert.Equal(t, movieID, g.MovieID)
		assert.NotEqual(t, uuid.Nil, g.ID)
	}
}

func TestResolutionToString(t *testing.T) {
	tests := []struct {
		resolution int
		expected   string
	}{
		{2160, "4K"},
		{1080, "1080p"},
		{720, "720p"},
		{480, "480p"},
		{360, "SD"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, resolutionToString(tt.resolution))
		})
	}
}

func TestFormatAudioChannels(t *testing.T) {
	tests := []struct {
		channels float64
		expected string
	}{
		{2.0, "2.0"},
		{5.1, "5.1"},
		{7.1, "7.1"},
		{1.0, "1.0"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, formatAudioChannels(tt.channels))
		})
	}
}
