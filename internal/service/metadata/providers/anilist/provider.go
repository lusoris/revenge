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
	metadata.TVShowProviderBase
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

func (p *Provider) ID() metadata.ProviderID { return metadata.ProviderAniList }
func (p *Provider) Name() string            { return "AniList" }
func (p *Provider) Priority() int           { return p.priority }
func (p *Provider) SupportsMovies() bool    { return false }
func (p *Provider) SupportsTVShows() bool   { return true }
func (p *Provider) SupportsPeople() bool    { return false }
func (p *Provider) ClearCache()             { p.client.clearCache() }

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
	page := max(opts.Page, 1)

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
