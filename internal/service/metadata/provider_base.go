package metadata

import "context"

// MovieProviderBase provides default ErrNotFound implementations for all
// MovieProvider methods. Embed this struct in provider types that only
// implement a subset of the MovieProvider interface, then override the
// methods that the provider actually supports.
//
// Example:
//
//	type Provider struct {
//	    metadata.MovieProviderBase
//	    client *Client
//	}
//
//	// Only override what this provider supports:
//	func (p *Provider) SearchMovie(ctx context.Context, ...) { ... }
//	func (p *Provider) GetMovie(ctx context.Context, ...) { ... }
type MovieProviderBase struct{}

func (MovieProviderBase) SearchMovie(_ context.Context, _ string, _ SearchOptions) ([]MovieSearchResult, error) {
	return nil, ErrNotFound
}

func (MovieProviderBase) GetMovie(_ context.Context, _ string, _ string) (*MovieMetadata, error) {
	return nil, ErrNotFound
}

func (MovieProviderBase) GetMovieCredits(_ context.Context, _ string) (*Credits, error) {
	return nil, ErrNotFound
}

func (MovieProviderBase) GetMovieImages(_ context.Context, _ string) (*Images, error) {
	return nil, ErrNotFound
}

func (MovieProviderBase) GetMovieReleaseDates(_ context.Context, _ string) ([]ReleaseDate, error) {
	return nil, ErrNotFound
}

func (MovieProviderBase) GetMovieTranslations(_ context.Context, _ string) ([]Translation, error) {
	return nil, ErrNotFound
}

func (MovieProviderBase) GetMovieExternalIDs(_ context.Context, _ string) (*ExternalIDs, error) {
	return nil, ErrNotFound
}

func (MovieProviderBase) GetSimilarMovies(_ context.Context, _ string, _ SearchOptions) ([]MovieSearchResult, int, error) {
	return nil, 0, ErrNotFound
}

func (MovieProviderBase) GetMovieRecommendations(_ context.Context, _ string, _ SearchOptions) ([]MovieSearchResult, int, error) {
	return nil, 0, ErrNotFound
}

// TVShowProviderBase provides default ErrNotFound implementations for all
// TVShowProvider methods. Embed this struct in provider types that only
// implement a subset of the TVShowProvider interface, then override the
// methods that the provider actually supports.
type TVShowProviderBase struct{}

func (TVShowProviderBase) SearchTVShow(_ context.Context, _ string, _ SearchOptions) ([]TVShowSearchResult, error) {
	return nil, ErrNotFound
}

func (TVShowProviderBase) GetTVShow(_ context.Context, _ string, _ string) (*TVShowMetadata, error) {
	return nil, ErrNotFound
}

func (TVShowProviderBase) GetTVShowCredits(_ context.Context, _ string) (*Credits, error) {
	return nil, ErrNotFound
}

func (TVShowProviderBase) GetTVShowImages(_ context.Context, _ string) (*Images, error) {
	return nil, ErrNotFound
}

func (TVShowProviderBase) GetTVShowContentRatings(_ context.Context, _ string) ([]ContentRating, error) {
	return nil, ErrNotFound
}

func (TVShowProviderBase) GetTVShowTranslations(_ context.Context, _ string) ([]Translation, error) {
	return nil, ErrNotFound
}

func (TVShowProviderBase) GetTVShowExternalIDs(_ context.Context, _ string) (*ExternalIDs, error) {
	return nil, ErrNotFound
}

func (TVShowProviderBase) GetSeason(_ context.Context, _ string, _ int, _ string) (*SeasonMetadata, error) {
	return nil, ErrNotFound
}

func (TVShowProviderBase) GetSeasonCredits(_ context.Context, _ string, _ int) (*Credits, error) {
	return nil, ErrNotFound
}

func (TVShowProviderBase) GetSeasonImages(_ context.Context, _ string, _ int) (*Images, error) {
	return nil, ErrNotFound
}

func (TVShowProviderBase) GetEpisode(_ context.Context, _ string, _, _ int, _ string) (*EpisodeMetadata, error) {
	return nil, ErrNotFound
}

func (TVShowProviderBase) GetEpisodeCredits(_ context.Context, _ string, _, _ int) (*Credits, error) {
	return nil, ErrNotFound
}

func (TVShowProviderBase) GetEpisodeImages(_ context.Context, _ string, _, _ int) (*Images, error) {
	return nil, ErrNotFound
}
