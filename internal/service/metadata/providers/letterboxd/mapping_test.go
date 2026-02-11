package letterboxd

import (
	"testing"
	"time"

	"github.com/lusoris/revenge/internal/service/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenreNameToID(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{"Action", 28},
		{"Drama", 18},
		{"Science Fiction", 878},
		{"Unknown", 0},
		// case-insensitive fallback
		{"action", 28},
		{"DRAMA", 18},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, genreNameToID(tt.name))
		})
	}
}

func TestMapDepartment(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"ProductionDesign", "Art"},
		{"Costumes", "Costume & Make-Up"},
		{"VisualEffects", "Visual Effects"},
		{"Lighting", "Lighting"},
		{"Sound", "Sound"},
		{"TitleDesign", "Art"},
		{"Stunts", "Crew"},
		{"Casting", "Production"},
		{"Studio", "Production"},
		{"SomethingElse", "Crew"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, mapDepartment(tt.input))
		})
	}
}

func TestLargestImageURL(t *testing.T) {
	t.Run("nil image", func(t *testing.T) {
		assert.Nil(t, largestImageURL(nil))
	})

	t.Run("empty sizes", func(t *testing.T) {
		img := &Image{Sizes: []ImageSize{}}
		assert.Nil(t, largestImageURL(img))
	})

	t.Run("single size", func(t *testing.T) {
		img := &Image{Sizes: []ImageSize{
			{Width: 200, Height: 300, URL: "https://example.com/small.jpg"},
		}}
		result := largestImageURL(img)
		require.NotNil(t, result)
		assert.Equal(t, "https://example.com/small.jpg", *result)
	})

	t.Run("picks largest area", func(t *testing.T) {
		img := &Image{Sizes: []ImageSize{
			{Width: 200, Height: 300, URL: "https://example.com/small.jpg"},
			{Width: 1000, Height: 1500, URL: "https://example.com/large.jpg"},
			{Width: 500, Height: 750, URL: "https://example.com/medium.jpg"},
		}}
		result := largestImageURL(img)
		require.NotNil(t, result)
		assert.Equal(t, "https://example.com/large.jpg", *result)
	})

	t.Run("empty URL in largest", func(t *testing.T) {
		img := &Image{Sizes: []ImageSize{
			{Width: 1000, Height: 1500, URL: ""},
		}}
		assert.Nil(t, largestImageURL(img))
	})
}

func TestExtractExternalIDs(t *testing.T) {
	links := []Link{
		{Type: "tmdb", ID: "123"},
		{Type: "imdb", ID: "tt1234567"},
		{Type: "letterboxd", URL: "https://letterboxd.com/film/test/"},
	}
	tmdbID, imdbID := extractExternalIDs(links)
	assert.Equal(t, 123, tmdbID)
	assert.Equal(t, "tt1234567", imdbID)

	// Invalid TMDb ID
	links2 := []Link{{Type: "tmdb", ID: "bad"}}
	tmdbID2, imdbID2 := extractExternalIDs(links2)
	assert.Equal(t, 0, tmdbID2)
	assert.Equal(t, "", imdbID2)

	// Empty links
	tmdbID3, imdbID3 := extractExternalIDs(nil)
	assert.Equal(t, 0, tmdbID3)
	assert.Equal(t, "", imdbID3)
}

func TestMapFilmSummaryToSearchResult(t *testing.T) {
	t.Run("full summary", func(t *testing.T) {
		summary := &FilmSummary{
			ID:           "abc123",
			Name:         "Inception",
			OriginalName: "Inception Original",
			ReleaseYear:  2010,
			Rating:       4.5,
			Adult:        false,
			Poster: &Image{Sizes: []ImageSize{
				{Width: 500, Height: 750, URL: "https://example.com/poster.jpg"},
			}},
			Genres: []Genre{{Name: "Action"}, {Name: "Science Fiction"}},
		}
		result := mapFilmSummaryToSearchResult(summary)
		assert.Equal(t, "abc123", result.ProviderID)
		assert.Equal(t, metadata.ProviderLetterboxd, result.Provider)
		assert.Equal(t, "Inception", result.Title)
		assert.Equal(t, "Inception Original", result.OriginalTitle)
		require.NotNil(t, result.Year)
		assert.Equal(t, 2010, *result.Year)
		require.NotNil(t, result.ReleaseDate)
		assert.Equal(t, 2010, result.ReleaseDate.Year())
		assert.InDelta(t, 9.0, result.VoteAverage, 0.01) // 4.5 * 2.0
		require.NotNil(t, result.PosterPath)
		assert.Equal(t, "https://example.com/poster.jpg", *result.PosterPath)
		assert.Equal(t, []int{28, 878}, result.GenreIDs)
	})

	t.Run("no year", func(t *testing.T) {
		summary := &FilmSummary{ID: "x", Name: "Test", ReleaseYear: 0}
		result := mapFilmSummaryToSearchResult(summary)
		assert.Nil(t, result.Year)
		assert.Nil(t, result.ReleaseDate)
	})

	t.Run("no rating", func(t *testing.T) {
		summary := &FilmSummary{ID: "x", Name: "Test"}
		result := mapFilmSummaryToSearchResult(summary)
		assert.InDelta(t, 0.0, result.VoteAverage, 0.01)
	})
}

func TestMapFilmToMetadata(t *testing.T) {
	t.Run("full film", func(t *testing.T) {
		film := &Film{
			ID:              "abc123",
			Name:            "Inception",
			OriginalName:    "Inception Original",
			Tagline:         "Your mind is the scene of the crime.",
			Description:     "A thief who steals corporate secrets.",
			ReleaseYear:     2010,
			RunTime:         148,
			Rating:          4.5,
			FilmCollectionID: "coll-1",
			Poster: &Image{Sizes: []ImageSize{
				{Width: 500, Height: 750, URL: "https://example.com/poster.jpg"},
			}},
			Backdrop: &Image{Sizes: []ImageSize{
				{Width: 1920, Height: 1080, URL: "https://example.com/backdrop.jpg"},
			}},
			Trailer: &FilmTrailer{URL: "https://youtube.com/watch?v=abc"},
			Genres:  []Genre{{Name: "Action"}, {Name: "Drama"}},
			Countries: []Country{{Code: "US", Name: "United States"}},
			Languages: []Language{{Code: "en", Name: "English"}},
			PrimaryLanguage: &Language{Code: "en", Name: "English"},
			Links: []Link{
				{Type: "tmdb", ID: "27205"},
				{Type: "imdb", ID: "tt1375666"},
				{Type: "letterboxd", URL: "https://letterboxd.com/film/inception/"},
			},
		}
		result := mapFilmToMetadata(film)
		assert.Equal(t, "abc123", result.ProviderID)
		assert.Equal(t, metadata.ProviderLetterboxd, result.Provider)
		assert.Equal(t, "Inception", result.Title)
		assert.Equal(t, "Inception Original", result.OriginalTitle)
		require.NotNil(t, result.Tagline)
		assert.Equal(t, "Your mind is the scene of the crime.", *result.Tagline)
		require.NotNil(t, result.Overview)
		require.NotNil(t, result.ReleaseDate)
		assert.Equal(t, time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC), *result.ReleaseDate)
		require.NotNil(t, result.Runtime)
		assert.Equal(t, int32(148), *result.Runtime)
		assert.InDelta(t, 9.0, result.VoteAverage, 0.01)
		require.NotNil(t, result.PosterPath)
		require.NotNil(t, result.BackdropPath)
		require.NotNil(t, result.TrailerURL)
		assert.Equal(t, "https://youtube.com/watch?v=abc", *result.TrailerURL)
		require.Len(t, result.Genres, 2)
		assert.Equal(t, "Action", result.Genres[0].Name)
		require.Len(t, result.ProductionCountries, 1)
		assert.Equal(t, "US", result.ProductionCountries[0].ISOCode)
		require.Len(t, result.SpokenLanguages, 1)
		assert.Equal(t, "en", result.OriginalLanguage)
		require.NotNil(t, result.IMDbID)
		assert.Equal(t, "tt1375666", *result.IMDbID)
		require.NotNil(t, result.TMDbID)
		assert.Equal(t, int32(27205), *result.TMDbID)
		require.NotNil(t, result.Collection)
		require.NotNil(t, result.Homepage)
		assert.Equal(t, "https://letterboxd.com/film/inception/", *result.Homepage)

		// External ratings
		require.Len(t, result.ExternalRatings, 1)
		assert.Equal(t, "Letterboxd", result.ExternalRatings[0].Source)
		assert.InDelta(t, 90.0, result.ExternalRatings[0].Score, 0.01) // 4.5 * 20
	})

	t.Run("minimal film", func(t *testing.T) {
		film := &Film{ID: "x", Name: "Test"}
		result := mapFilmToMetadata(film)
		assert.Equal(t, "x", result.ProviderID)
		assert.Equal(t, "Test", result.Title)
		assert.Nil(t, result.Tagline)
		assert.Nil(t, result.Overview)
		assert.Nil(t, result.ReleaseDate)
		assert.Nil(t, result.Runtime)
		assert.Nil(t, result.PosterPath)
	})
}

func TestMapCredits(t *testing.T) {
	contributions := []Contributions{
		{
			Type: "Director",
			Contributors: []ContributorSummary{
				{Name: "Christopher Nolan"},
			},
		},
		{
			Type: "Actor",
			Contributors: []ContributorSummary{
				{Name: "Leonardo DiCaprio", CharacterName: "Cobb"},
			},
		},
		{
			Type: "Writer",
			Contributors: []ContributorSummary{
				{Name: "Christopher Nolan"},
			},
		},
		{
			Type: "Composer",
			Contributors: []ContributorSummary{
				{Name: "Hans Zimmer"},
			},
		},
		{
			Type: "Cinematography",
			Contributors: []ContributorSummary{
				{Name: "Wally Pfister"},
			},
		},
		{
			Type: "Editor",
			Contributors: []ContributorSummary{
				{Name: "Lee Smith"},
			},
		},
		{
			Type: "Producer",
			Contributors: []ContributorSummary{
				{Name: "Emma Thomas"},
			},
		},
		{
			Type: "VisualEffects",
			Contributors: []ContributorSummary{
				{Name: "Chris Corbould"},
			},
		},
	}

	result := mapCredits(contributions)
	require.Len(t, result.Cast, 1)
	assert.Equal(t, "Leonardo DiCaprio", result.Cast[0].Name)
	assert.Equal(t, "Cobb", result.Cast[0].Character)

	// Count crew by department
	departments := map[string]int{}
	for _, c := range result.Crew {
		departments[c.Department]++
	}
	assert.Equal(t, 1, departments["Directing"])
	assert.Equal(t, 1, departments["Writing"])
	assert.Equal(t, 1, departments["Sound"])
	assert.Equal(t, 1, departments["Camera"])
	assert.Equal(t, 1, departments["Editing"])
	assert.Equal(t, 1, departments["Production"])
	assert.Equal(t, 1, departments["Visual Effects"])
}
