package fanarttv

import (
	"context"
	"fmt"

	"github.com/lusoris/revenge/internal/service/metadata"
)

// Ensure Provider implements required interfaces.
// Fanart.tv is an image-focused provider: it provides high-quality artwork
// (logos, clearart, disc art, banners) but does NOT support search, credits,
// or detailed metadata. Non-image methods return ErrNotFound.
var (
	_ metadata.Provider      = (*Provider)(nil)
	_ metadata.MovieProvider = (*Provider)(nil)
	_ metadata.TVShowProvider = (*Provider)(nil)
	_ metadata.ImageProvider = (*Provider)(nil)
)

// Provider implements the metadata provider interface for Fanart.tv.
type Provider struct {
	client   *Client
	priority int
}

// NewProvider creates a new Fanart.tv provider.
func NewProvider(config Config) (*Provider, error) {
	client, err := NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("create fanarttv provider: %w", err)
	}
	return &Provider{
		client:   client,
		priority: 60, // Lower than TMDb (100) and TVDb (80)
	}, nil
}

func (p *Provider) ID() metadata.ProviderID     { return metadata.ProviderFanartTV }
func (p *Provider) Name() string                 { return "Fanart.tv" }
func (p *Provider) Priority() int                { return p.priority }
func (p *Provider) SupportsMovies() bool         { return true }
func (p *Provider) SupportsTVShows() bool        { return true }
func (p *Provider) SupportsPeople() bool         { return false }
func (p *Provider) SupportsLanguage(_ string) bool { return true }
func (p *Provider) ClearCache()                  { p.client.clearCache() }

// --- ImageProvider ---

func (p *Provider) GetImageURL(path string, _ metadata.ImageSize) string {
	// Fanart.tv returns full URLs, no size transformation needed.
	return path
}

func (p *Provider) GetImageBaseURL() string {
	return "https://assets.fanart.tv/fanart"
}

func (p *Provider) DownloadImage(_ context.Context, _ string, _ metadata.ImageSize) ([]byte, error) {
	return nil, metadata.ErrNotFound
}

// --- MovieProvider (image methods only) ---

func (p *Provider) SearchMovie(_ context.Context, _ string, _ metadata.SearchOptions) ([]metadata.MovieSearchResult, error) {
	return nil, metadata.ErrNotFound
}

func (p *Provider) GetMovie(_ context.Context, _ string, _ string) (*metadata.MovieMetadata, error) {
	return nil, metadata.ErrNotFound
}

func (p *Provider) GetMovieCredits(_ context.Context, _ string) (*metadata.Credits, error) {
	return nil, metadata.ErrNotFound
}

func (p *Provider) GetMovieImages(ctx context.Context, id string) (*metadata.Images, error) {
	resp, err := p.client.GetMovieImages(ctx, id)
	if err != nil {
		return nil, err
	}
	images := mapMovieImages(resp)
	if images == nil {
		return nil, metadata.ErrNotFound
	}
	return images, nil
}

func (p *Provider) GetMovieReleaseDates(_ context.Context, _ string) ([]metadata.ReleaseDate, error) {
	return nil, metadata.ErrNotFound
}

func (p *Provider) GetMovieTranslations(_ context.Context, _ string) ([]metadata.Translation, error) {
	return nil, metadata.ErrNotFound
}

func (p *Provider) GetMovieExternalIDs(_ context.Context, _ string) (*metadata.ExternalIDs, error) {
	return nil, metadata.ErrNotFound
}

func (p *Provider) GetSimilarMovies(_ context.Context, _ string, _ metadata.SearchOptions) ([]metadata.MovieSearchResult, int, error) {
	return nil, 0, metadata.ErrNotFound
}

func (p *Provider) GetMovieRecommendations(_ context.Context, _ string, _ metadata.SearchOptions) ([]metadata.MovieSearchResult, int, error) {
	return nil, 0, metadata.ErrNotFound
}

// --- TVShowProvider (image methods only) ---
// NOTE: Fanart.tv uses TVDb IDs for TV shows. When the service passes TMDb IDs,
// this may return 404. A TMDb→TVDb ID mapping at the service layer would fix this.

func (p *Provider) SearchTVShow(_ context.Context, _ string, _ metadata.SearchOptions) ([]metadata.TVShowSearchResult, error) {
	return nil, metadata.ErrNotFound
}

func (p *Provider) GetTVShow(_ context.Context, _ string, _ string) (*metadata.TVShowMetadata, error) {
	return nil, metadata.ErrNotFound
}

func (p *Provider) GetTVShowCredits(_ context.Context, _ string) (*metadata.Credits, error) {
	return nil, metadata.ErrNotFound
}

func (p *Provider) GetTVShowImages(ctx context.Context, id string) (*metadata.Images, error) {
	resp, err := p.client.GetTVShowImages(ctx, id)
	if err != nil {
		return nil, err
	}
	images := mapTVShowImages(resp)
	if images == nil {
		return nil, metadata.ErrNotFound
	}
	return images, nil
}

func (p *Provider) GetTVShowContentRatings(_ context.Context, _ string) ([]metadata.ContentRating, error) {
	return nil, metadata.ErrNotFound
}

func (p *Provider) GetTVShowTranslations(_ context.Context, _ string) ([]metadata.Translation, error) {
	return nil, metadata.ErrNotFound
}

func (p *Provider) GetTVShowExternalIDs(_ context.Context, _ string) (*metadata.ExternalIDs, error) {
	return nil, metadata.ErrNotFound
}

func (p *Provider) GetSeason(_ context.Context, _ string, _ int, _ string) (*metadata.SeasonMetadata, error) {
	return nil, metadata.ErrNotFound
}

func (p *Provider) GetSeasonCredits(_ context.Context, _ string, _ int) (*metadata.Credits, error) {
	return nil, metadata.ErrNotFound
}

func (p *Provider) GetSeasonImages(ctx context.Context, showID string, seasonNum int) (*metadata.Images, error) {
	resp, err := p.client.GetTVShowImages(ctx, showID)
	if err != nil {
		return nil, err
	}
	images := mapSeasonImages(resp, seasonNum)
	if images == nil {
		return nil, metadata.ErrNotFound
	}
	return images, nil
}

func (p *Provider) GetEpisode(_ context.Context, _ string, _, _ int, _ string) (*metadata.EpisodeMetadata, error) {
	return nil, metadata.ErrNotFound
}

func (p *Provider) GetEpisodeCredits(_ context.Context, _ string, _, _ int) (*metadata.Credits, error) {
	return nil, metadata.ErrNotFound
}

func (p *Provider) GetEpisodeImages(_ context.Context, _ string, _, _ int) (*metadata.Images, error) {
	// Fanart.tv doesn't have per-episode images
	return nil, metadata.ErrNotFound
}
