package mal

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lusoris/revenge/internal/service/metadata"
)

// Ensure Provider implements required interfaces.
// MAL is a major anime database with community ratings, rankings, and user lists.
// The v2 API provides anime metadata with a simple client ID auth (no OAuth for reads).
// It does NOT support movies (non-anime), people search, or per-episode data.
var (
	_ metadata.Provider       = (*Provider)(nil)
	_ metadata.TVShowProvider = (*Provider)(nil)
)

// Provider implements the metadata provider interface for MyAnimeList.
type Provider struct {
	client   *Client
	priority int
}

// NewProvider creates a new MAL provider.
func NewProvider(config Config) (*Provider, error) {
	client, err := NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("create mal provider: %w", err)
	}
	return &Provider{
		client:   client,
		priority: 43, // Below AniList (45), above Kitsu (42)
	}, nil
}

func (p *Provider) ID() metadata.ProviderID       { return metadata.ProviderMAL }
func (p *Provider) Name() string                   { return "MyAnimeList" }
func (p *Provider) Priority() int                  { return p.priority }
func (p *Provider) SupportsMovies() bool           { return false }
func (p *Provider) SupportsTVShows() bool          { return true }
func (p *Provider) SupportsPeople() bool           { return false }
func (p *Provider) ClearCache()                    { p.client.clearCache() }

// SupportsLanguage returns true for Japanese and English content.
func (p *Provider) SupportsLanguage(lang string) bool {
	switch lang {
	case "ja", "en":
		return true
	default:
		return false
	}
}

// --- TVShowProvider ---

func (p *Provider) SearchTVShow(ctx context.Context, query string, opts metadata.SearchOptions) ([]metadata.TVShowSearchResult, error) {
	page := opts.Page
	if page < 1 {
		page = 1
	}
	limit := 20
	offset := (page - 1) * limit

	result, err := p.client.SearchAnime(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result.Data) == 0 {
		return nil, metadata.ErrNotFound
	}

	searchResults := make([]metadata.TVShowSearchResult, 0, len(result.Data))
	for _, n := range result.Data {
		searchResults = append(searchResults, mapAnimeToTVShowSearchResult(n.Node))
	}
	return searchResults, nil
}

func (p *Provider) GetTVShow(ctx context.Context, id string, _ string) (*metadata.TVShowMetadata, error) {
	animeID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("mal: invalid anime ID %q: %w", id, err)
	}

	anime, err := p.client.GetAnime(ctx, animeID)
	if err != nil {
		return nil, err
	}

	m := mapAnimeToTVShowMetadata(anime)
	if m == nil {
		return nil, metadata.ErrNotFound
	}
	return m, nil
}

func (p *Provider) GetTVShowCredits(_ context.Context, _ string) (*metadata.Credits, error) {
	// MAL API v2 does not expose characters or staff
	return nil, metadata.ErrNotFound
}

func (p *Provider) GetTVShowImages(ctx context.Context, id string) (*metadata.Images, error) {
	animeID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("mal: invalid anime ID %q: %w", id, err)
	}

	anime, err := p.client.GetAnime(ctx, animeID)
	if err != nil {
		return nil, err
	}

	images := mapImages(anime)
	if images == nil {
		return nil, metadata.ErrNotFound
	}
	return images, nil
}

func (p *Provider) GetTVShowContentRatings(ctx context.Context, id string) ([]metadata.ContentRating, error) {
	animeID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("mal: invalid anime ID %q: %w", id, err)
	}

	anime, err := p.client.GetAnime(ctx, animeID)
	if err != nil {
		return nil, err
	}
	if anime == nil || anime.Rating == "" {
		return nil, metadata.ErrNotFound
	}

	return []metadata.ContentRating{
		{CountryCode: "JP", Rating: mapRating(anime.Rating)},
	}, nil
}

func (p *Provider) GetTVShowTranslations(_ context.Context, _ string) ([]metadata.Translation, error) {
	return nil, metadata.ErrNotFound
}

func (p *Provider) GetTVShowExternalIDs(_ context.Context, _ string) (*metadata.ExternalIDs, error) {
	// MAL doesn't provide external ID cross-references in the API
	return nil, metadata.ErrNotFound
}

// MAL API v2 does not have per-season or per-episode data.

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
