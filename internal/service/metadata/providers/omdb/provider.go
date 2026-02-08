package omdb

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lusoris/revenge/internal/service/metadata"
)

// Ensure Provider implements required interfaces.
// OMDb is a ratings-focused provider: it provides IMDb, Rotten Tomatoes,
// and Metacritic ratings plus basic metadata. It does NOT support images,
// credits, translations, or detailed season/episode info.
var (
	_ metadata.Provider      = (*Provider)(nil)
	_ metadata.MovieProvider = (*Provider)(nil)
	_ metadata.TVShowProvider = (*Provider)(nil)
)

// Provider implements the metadata provider interface for OMDb.
type Provider struct {
	client   *Client
	priority int
}

// NewProvider creates a new OMDb provider.
func NewProvider(config Config) (*Provider, error) {
	client, err := NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("create omdb provider: %w", err)
	}
	return &Provider{
		client:   client,
		priority: 40, // Lower than TMDb (100), TVDb (80), Fanart.tv (60)
	}, nil
}

func (p *Provider) ID() metadata.ProviderID       { return metadata.ProviderOMDb }
func (p *Provider) Name() string                   { return "OMDb" }
func (p *Provider) Priority() int                  { return p.priority }
func (p *Provider) SupportsMovies() bool           { return true }
func (p *Provider) SupportsTVShows() bool          { return true }
func (p *Provider) SupportsPeople() bool           { return false }
func (p *Provider) SupportsLanguage(_ string) bool { return true } // English-only API
func (p *Provider) ClearCache()                    { p.client.clearCache() }

// --- MovieProvider ---

func (p *Provider) SearchMovie(ctx context.Context, query string, opts metadata.SearchOptions) ([]metadata.MovieSearchResult, error) {
	year := ""
	if opts.Year != nil {
		year = strconv.Itoa(*opts.Year)
	}
	resp, err := p.client.Search(ctx, query, year, "movie", opts.Page)
	if err != nil {
		return nil, err
	}
	results := mapMovieSearchResults(resp)
	if results == nil {
		return nil, metadata.ErrNotFound
	}
	return results, nil
}

func (p *Provider) GetMovie(ctx context.Context, id string, _ string) (*metadata.MovieMetadata, error) {
	// OMDb uses IMDb IDs (tt-prefixed)
	resp, err := p.client.GetByIMDbID(ctx, id)
	if err != nil {
		return nil, err
	}
	m := mapMovieMetadata(resp)
	if m == nil {
		return nil, metadata.ErrNotFound
	}
	return m, nil
}

func (p *Provider) GetMovieCredits(_ context.Context, _ string) (*metadata.Credits, error) {
	return nil, metadata.ErrNotFound
}

func (p *Provider) GetMovieImages(_ context.Context, _ string) (*metadata.Images, error) {
	return nil, metadata.ErrNotFound
}

func (p *Provider) GetMovieReleaseDates(_ context.Context, _ string) ([]metadata.ReleaseDate, error) {
	return nil, metadata.ErrNotFound
}

func (p *Provider) GetMovieTranslations(_ context.Context, _ string) ([]metadata.Translation, error) {
	return nil, metadata.ErrNotFound
}

func (p *Provider) GetMovieExternalIDs(ctx context.Context, id string) (*metadata.ExternalIDs, error) {
	resp, err := p.client.GetByIMDbID(ctx, id)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, metadata.ErrNotFound
	}
	return &metadata.ExternalIDs{
		IMDbID: &resp.IMDbID,
	}, nil
}

func (p *Provider) GetSimilarMovies(_ context.Context, _ string, _ metadata.SearchOptions) ([]metadata.MovieSearchResult, int, error) {
	return nil, 0, metadata.ErrNotFound
}

func (p *Provider) GetMovieRecommendations(_ context.Context, _ string, _ metadata.SearchOptions) ([]metadata.MovieSearchResult, int, error) {
	return nil, 0, metadata.ErrNotFound
}

// --- TVShowProvider ---

func (p *Provider) SearchTVShow(ctx context.Context, query string, opts metadata.SearchOptions) ([]metadata.TVShowSearchResult, error) {
	year := ""
	if opts.Year != nil {
		year = strconv.Itoa(*opts.Year)
	}
	resp, err := p.client.Search(ctx, query, year, "series", opts.Page)
	if err != nil {
		return nil, err
	}
	results := mapTVShowSearchResults(resp)
	if results == nil {
		return nil, metadata.ErrNotFound
	}
	return results, nil
}

func (p *Provider) GetTVShow(ctx context.Context, id string, _ string) (*metadata.TVShowMetadata, error) {
	resp, err := p.client.GetByIMDbID(ctx, id)
	if err != nil {
		return nil, err
	}
	m := mapTVShowMetadata(resp)
	if m == nil {
		return nil, metadata.ErrNotFound
	}
	return m, nil
}

func (p *Provider) GetTVShowCredits(_ context.Context, _ string) (*metadata.Credits, error) {
	return nil, metadata.ErrNotFound
}

func (p *Provider) GetTVShowImages(_ context.Context, _ string) (*metadata.Images, error) {
	return nil, metadata.ErrNotFound
}

func (p *Provider) GetTVShowContentRatings(_ context.Context, _ string) ([]metadata.ContentRating, error) {
	return nil, metadata.ErrNotFound
}

func (p *Provider) GetTVShowTranslations(_ context.Context, _ string) ([]metadata.Translation, error) {
	return nil, metadata.ErrNotFound
}

func (p *Provider) GetTVShowExternalIDs(ctx context.Context, id string) (*metadata.ExternalIDs, error) {
	resp, err := p.client.GetByIMDbID(ctx, id)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, metadata.ErrNotFound
	}
	return &metadata.ExternalIDs{
		IMDbID: &resp.IMDbID,
	}, nil
}

func (p *Provider) GetSeason(_ context.Context, _ string, _ int, _ string) (*metadata.SeasonMetadata, error) {
	return nil, metadata.ErrNotFound
}

func (p *Provider) GetSeasonCredits(_ context.Context, _ string, _ int) (*metadata.Credits, error) {
	return nil, metadata.ErrNotFound
}

func (p *Provider) GetSeasonImages(_ context.Context, _ string, _ int) (*metadata.Images, error) {
	return nil, metadata.ErrNotFound
}

func (p *Provider) GetEpisode(_ context.Context, _ string, _, _ int, _ string) (*metadata.EpisodeMetadata, error) {
	return nil, metadata.ErrNotFound
}

func (p *Provider) GetEpisodeCredits(_ context.Context, _ string, _, _ int) (*metadata.Credits, error) {
	return nil, metadata.ErrNotFound
}

func (p *Provider) GetEpisodeImages(_ context.Context, _ string, _, _ int) (*metadata.Images, error) {
	return nil, metadata.ErrNotFound
}
