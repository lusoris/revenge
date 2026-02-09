package metadata

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMovieProviderBase_AllMethodsReturnErrNotFound(t *testing.T) {
	t.Parallel()
	var base MovieProviderBase
	ctx := context.Background()
	opts := SearchOptions{}

	results, err := base.SearchMovie(ctx, "test", opts)
	assert.Nil(t, results)
	assert.ErrorIs(t, err, ErrNotFound)

	movie, err := base.GetMovie(ctx, "1", "en")
	assert.Nil(t, movie)
	assert.ErrorIs(t, err, ErrNotFound)

	credits, err := base.GetMovieCredits(ctx, "1")
	assert.Nil(t, credits)
	assert.ErrorIs(t, err, ErrNotFound)

	images, err := base.GetMovieImages(ctx, "1")
	assert.Nil(t, images)
	assert.ErrorIs(t, err, ErrNotFound)

	releaseDates, err := base.GetMovieReleaseDates(ctx, "1")
	assert.Nil(t, releaseDates)
	assert.ErrorIs(t, err, ErrNotFound)

	translations, err := base.GetMovieTranslations(ctx, "1")
	assert.Nil(t, translations)
	assert.ErrorIs(t, err, ErrNotFound)

	extIDs, err := base.GetMovieExternalIDs(ctx, "1")
	assert.Nil(t, extIDs)
	assert.ErrorIs(t, err, ErrNotFound)

	similar, count, err := base.GetSimilarMovies(ctx, "1", opts)
	assert.Nil(t, similar)
	assert.Equal(t, 0, count)
	assert.ErrorIs(t, err, ErrNotFound)

	recs, count, err := base.GetMovieRecommendations(ctx, "1", opts)
	assert.Nil(t, recs)
	assert.Equal(t, 0, count)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestTVShowProviderBase_AllMethodsReturnErrNotFound(t *testing.T) {
	t.Parallel()
	var base TVShowProviderBase
	ctx := context.Background()
	opts := SearchOptions{}

	results, err := base.SearchTVShow(ctx, "test", opts)
	assert.Nil(t, results)
	assert.ErrorIs(t, err, ErrNotFound)

	show, err := base.GetTVShow(ctx, "1", "en")
	assert.Nil(t, show)
	assert.ErrorIs(t, err, ErrNotFound)

	credits, err := base.GetTVShowCredits(ctx, "1")
	assert.Nil(t, credits)
	assert.ErrorIs(t, err, ErrNotFound)

	images, err := base.GetTVShowImages(ctx, "1")
	assert.Nil(t, images)
	assert.ErrorIs(t, err, ErrNotFound)

	ratings, err := base.GetTVShowContentRatings(ctx, "1")
	assert.Nil(t, ratings)
	assert.ErrorIs(t, err, ErrNotFound)

	translations, err := base.GetTVShowTranslations(ctx, "1")
	assert.Nil(t, translations)
	assert.ErrorIs(t, err, ErrNotFound)

	extIDs, err := base.GetTVShowExternalIDs(ctx, "1")
	assert.Nil(t, extIDs)
	assert.ErrorIs(t, err, ErrNotFound)

	season, err := base.GetSeason(ctx, "1", 1, "en")
	assert.Nil(t, season)
	assert.ErrorIs(t, err, ErrNotFound)

	seasonCredits, err := base.GetSeasonCredits(ctx, "1", 1)
	assert.Nil(t, seasonCredits)
	assert.ErrorIs(t, err, ErrNotFound)

	seasonImages, err := base.GetSeasonImages(ctx, "1", 1)
	assert.Nil(t, seasonImages)
	assert.ErrorIs(t, err, ErrNotFound)

	episode, err := base.GetEpisode(ctx, "1", 1, 1, "en")
	assert.Nil(t, episode)
	assert.ErrorIs(t, err, ErrNotFound)

	epCredits, err := base.GetEpisodeCredits(ctx, "1", 1, 1)
	assert.Nil(t, epCredits)
	assert.ErrorIs(t, err, ErrNotFound)

	epImages, err := base.GetEpisodeImages(ctx, "1", 1, 1)
	assert.Nil(t, epImages)
	assert.ErrorIs(t, err, ErrNotFound)
}

// TestProviderBase_Embedding verifies that a concrete provider type can embed
// the base structs and override only the methods it needs.
func TestProviderBase_Embedding(t *testing.T) {
	t.Parallel()

	type testProvider struct {
		MovieProviderBase
		TVShowProviderBase
	}

	p := &testProvider{}
	ctx := context.Background()

	// Inherited from MovieProviderBase
	_, err := p.GetMovieCredits(ctx, "1")
	assert.ErrorIs(t, err, ErrNotFound)

	// Inherited from TVShowProviderBase
	_, err = p.GetTVShowCredits(ctx, "1")
	assert.ErrorIs(t, err, ErrNotFound)
}
