package fanarttv

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapFanartImages(t *testing.T) {
	t.Run("nil slice", func(t *testing.T) {
		assert.Nil(t, mapFanartImages(nil, "poster"))
	})

	t.Run("empty slice", func(t *testing.T) {
		assert.Nil(t, mapFanartImages([]FanartImage{}, "poster"))
	})

	t.Run("with images", func(t *testing.T) {
		imgs := []FanartImage{
			{ID: "1", URL: "https://example.com/img1.jpg", Lang: "en", Likes: "10"},
			{ID: "2", URL: "https://example.com/img2.jpg", Lang: "00", Likes: "5"},
			{ID: "3", URL: "https://example.com/img3.jpg", Lang: "", Likes: "bad"},
		}
		result := mapFanartImages(imgs, "poster")
		require.Len(t, result, 3)

		// First image: language "en" should be set
		assert.Equal(t, "https://example.com/img1.jpg", result[0].FilePath)
		require.NotNil(t, result[0].Language)
		assert.Equal(t, "en", *result[0].Language)
		assert.Equal(t, 10, result[0].VoteCount)

		// Second image: language "00" should be nil (filtered)
		assert.Nil(t, result[1].Language)
		assert.Equal(t, 5, result[1].VoteCount)

		// Third image: empty language, bad likes
		assert.Nil(t, result[2].Language)
		assert.Equal(t, 0, result[2].VoteCount)
	})
}

func TestMapSeasonImage(t *testing.T) {
	img := SeasonImage{
		ID: "1", URL: "https://example.com/season.jpg", Lang: "de", Likes: "42", Season: "3",
	}
	result := mapSeasonImage(img, "poster")
	assert.Equal(t, "https://example.com/season.jpg", result.FilePath)
	require.NotNil(t, result.Language)
	assert.Equal(t, "de", *result.Language)
	assert.Equal(t, 42, result.VoteCount)

	// "00" language should be nil
	img2 := SeasonImage{URL: "https://example.com/s2.jpg", Lang: "00", Likes: "3"}
	result2 := mapSeasonImage(img2, "poster")
	assert.Nil(t, result2.Language)
}

func TestMapMovieImages(t *testing.T) {
	t.Run("nil response", func(t *testing.T) {
		assert.Nil(t, mapMovieImages(nil))
	})

	t.Run("full movie response", func(t *testing.T) {
		resp := &MovieResponse{
			HDMovieLogos:   []FanartImage{{URL: "https://example.com/hdlogo.jpg", Lang: "en", Likes: "5"}},
			MovieLogos:     []FanartImage{{URL: "https://example.com/logo.jpg", Lang: "en", Likes: "3"}},
			MoviePosters:   []FanartImage{{URL: "https://example.com/poster.jpg", Lang: "en", Likes: "10"}},
			MovieBackdrops: []FanartImage{{URL: "https://example.com/bg.jpg", Lang: "en", Likes: "7"}},
			MovieBanners:   []FanartImage{{URL: "https://example.com/banner.jpg", Lang: "en", Likes: "2"}},
			MovieThumbs:    []FanartImage{{URL: "https://example.com/thumb.jpg", Lang: "en", Likes: "1"}},
		}
		result := mapMovieImages(resp)
		require.NotNil(t, result)
		assert.Len(t, result.Logos, 2)   // HD + standard
		assert.Len(t, result.Posters, 1)
		assert.Len(t, result.Backdrops, 2) // backgrounds + banners
		assert.Len(t, result.Stills, 1)    // thumbs
	})

	t.Run("empty response", func(t *testing.T) {
		resp := &MovieResponse{}
		result := mapMovieImages(resp)
		require.NotNil(t, result)
		assert.Empty(t, result.Logos)
		assert.Empty(t, result.Posters)
	})
}

func TestMapTVShowImages(t *testing.T) {
	t.Run("nil response", func(t *testing.T) {
		assert.Nil(t, mapTVShowImages(nil))
	})

	t.Run("full TV response", func(t *testing.T) {
		resp := &TVShowResponse{
			HDTVLogos:    []FanartImage{{URL: "https://example.com/hdlogo.jpg", Likes: "5"}},
			ClearLogos:   []FanartImage{{URL: "https://example.com/clearlogo.jpg", Likes: "3"}},
			TVPosters:    []FanartImage{{URL: "https://example.com/poster.jpg", Likes: "10"}},
			ShowBackdrops: []FanartImage{{URL: "https://example.com/bg.jpg", Likes: "7"}},
			TVBanners:    []FanartImage{{URL: "https://example.com/banner.jpg", Likes: "2"}},
			TVThumbs:     []FanartImage{{URL: "https://example.com/thumb.jpg", Likes: "1"}},
			HDClearArt:   []FanartImage{{URL: "https://example.com/hdclearart.jpg", Likes: "4"}},
			ClearArt:     []FanartImage{{URL: "https://example.com/clearart.jpg", Likes: "2"}},
			CharacterArt: []FanartImage{{URL: "https://example.com/charart.jpg", Likes: "1"}},
		}
		result := mapTVShowImages(resp)
		require.NotNil(t, result)
		assert.Len(t, result.Logos, 2)     // HD + clear logos
		assert.Len(t, result.Posters, 1)
		assert.Len(t, result.Backdrops, 2) // backgrounds + banners
		assert.Len(t, result.Stills, 1)    // thumbs
		assert.Len(t, result.Profiles, 3)  // HD clearart + clearart + characterart
	})
}

func TestMapSeasonImages(t *testing.T) {
	t.Run("nil response", func(t *testing.T) {
		assert.Nil(t, mapSeasonImages(nil, 1))
	})

	t.Run("matching season", func(t *testing.T) {
		resp := &TVShowResponse{
			SeasonPosters: []SeasonImage{
				{URL: "https://example.com/s1.jpg", Season: "1", Likes: "5"},
				{URL: "https://example.com/s2.jpg", Season: "2", Likes: "3"},
			},
			SeasonThumbs: []SeasonImage{
				{URL: "https://example.com/s1t.jpg", Season: "1", Likes: "2"},
			},
			SeasonBanners: []SeasonImage{
				{URL: "https://example.com/s1b.jpg", Season: "1", Likes: "1"},
				{URL: "https://example.com/s3b.jpg", Season: "3", Likes: "1"},
			},
		}
		result := mapSeasonImages(resp, 1)
		require.NotNil(t, result)
		assert.Len(t, result.Posters, 1)
		assert.Equal(t, "https://example.com/s1.jpg", result.Posters[0].FilePath)
		assert.Len(t, result.Stills, 1)
		assert.Len(t, result.Backdrops, 1)
	})

	t.Run("no matching season", func(t *testing.T) {
		resp := &TVShowResponse{
			SeasonPosters: []SeasonImage{
				{URL: "https://example.com/s1.jpg", Season: "1", Likes: "5"},
			},
		}
		result := mapSeasonImages(resp, 99)
		assert.Nil(t, result)
	})
}
