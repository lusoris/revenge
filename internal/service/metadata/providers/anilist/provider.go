package anilist

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lusoris/revenge/internal/service/metadata"
)

// Ensure Provider implements required interfaces.
// AniList is an anime-focused provider with rich anime metadata, character data,
// staff info, and community ratings. It does NOT support movies (live-action),
// people search, or collections.
var (
	_ metadata.Provider       = (*Provider)(nil)
	_ metadata.TVShowProvider = (*Provider)(nil)
)

// Provider implements the metadata provider interface for AniList.
type Provider struct {
	client   *Client
	priority int
}

// NewProvider creates a new AniList provider.
func NewProvider(config Config) (*Provider, error) {
	client, err := NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("create anilist provider: %w", err)
	}
	return &Provider{
		client:   client,
		priority: 45, // Below TVmaze (50), focused on anime content
	}, nil
}

func (p *Provider) ID() metadata.ProviderID       { return metadata.ProviderAniList }
func (p *Provider) Name() string                   { return "AniList" }
func (p *Provider) Priority() int                  { return p.priority }
func (p *Provider) SupportsMovies() bool           { return false }
func (p *Provider) SupportsTVShows() bool          { return true }
func (p *Provider) SupportsPeople() bool           { return false }
func (p *Provider) ClearCache()                    { p.client.clearCache() }

// SupportsLanguage returns true for Japanese and English content.
func (p *Provider) SupportsLanguage(lang string) bool {
	switch lang {
	case "ja", "en", "ko", "zh":
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

	result, err := p.client.SearchAnime(ctx, query, page, 20, opts.IncludeAdult)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result.Media) == 0 {
		return nil, metadata.ErrNotFound
	}

	searchResults := make([]metadata.TVShowSearchResult, 0, len(result.Media))
	for _, m := range result.Media {
		searchResults = append(searchResults, mapMediaToTVShowSearchResult(m))
	}
	return searchResults, nil
}

func (p *Provider) GetTVShow(ctx context.Context, id string, _ string) (*metadata.TVShowMetadata, error) {
	animeID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("anilist: invalid anime ID %q: %w", id, err)
	}

	media, err := p.client.GetAnime(ctx, animeID)
	if err != nil {
		return nil, err
	}

	m := mapMediaToTVShowMetadata(media)
	if m == nil {
		return nil, metadata.ErrNotFound
	}
	return m, nil
}

func (p *Provider) GetTVShowCredits(ctx context.Context, id string) (*metadata.Credits, error) {
	animeID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("anilist: invalid anime ID %q: %w", id, err)
	}

	media, err := p.client.GetAnime(ctx, animeID)
	if err != nil {
		return nil, err
	}

	credits := mapCredits(media)
	if credits == nil || (len(credits.Cast) == 0 && len(credits.Crew) == 0) {
		return nil, metadata.ErrNotFound
	}
	return credits, nil
}

func (p *Provider) GetTVShowImages(ctx context.Context, id string) (*metadata.Images, error) {
	animeID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("anilist: invalid anime ID %q: %w", id, err)
	}

	media, err := p.client.GetAnime(ctx, animeID)
	if err != nil {
		return nil, err
	}

	images := mapImages(media)
	if images == nil {
		return nil, metadata.ErrNotFound
	}
	return images, nil
}

func (p *Provider) GetTVShowContentRatings(_ context.Context, _ string) ([]metadata.ContentRating, error) {
	// AniList doesn't provide regional content ratings
	return nil, metadata.ErrNotFound
}

func (p *Provider) GetTVShowTranslations(_ context.Context, _ string) ([]metadata.Translation, error) {
	// AniList provides limited translation data (title only)
	return nil, metadata.ErrNotFound
}

func (p *Provider) GetTVShowExternalIDs(ctx context.Context, id string) (*metadata.ExternalIDs, error) {
	animeID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("anilist: invalid anime ID %q: %w", id, err)
	}

	media, err := p.client.GetAnime(ctx, animeID)
	if err != nil {
		return nil, err
	}

	ids := findExternalIDs(media)
	if ids == nil {
		return nil, metadata.ErrNotFound
	}
	return ids, nil
}

// AniList doesn't have per-season/episode data in the same way western TV providers do.
// Anime typically has a single "season" representing the entire series run.

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
