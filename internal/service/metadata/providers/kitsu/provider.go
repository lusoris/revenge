package kitsu

import (
	"context"
	"fmt"

	"github.com/lusoris/revenge/internal/service/metadata"
)

// Ensure Provider implements required interfaces.
// Kitsu is an anime discovery platform with rich anime metadata, episode data,
// categories, and external ID mappings. It does NOT support movies (non-anime),
// people search, or collections.
var (
	_ metadata.Provider       = (*Provider)(nil)
	_ metadata.TVShowProvider = (*Provider)(nil)
)

// Provider implements the metadata provider interface for Kitsu.
type Provider struct {
	metadata.TVShowProviderBase
	client   *Client
	priority int
}

// NewProvider creates a new Kitsu provider.
func NewProvider(config Config) (*Provider, error) {
	client, err := NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("create kitsu provider: %w", err)
	}
	return &Provider{
		client:   client,
		priority: 42, // Below AniList (45), supplementary anime data
	}, nil
}

func (p *Provider) ID() metadata.ProviderID       { return metadata.ProviderKitsu }
func (p *Provider) Name() string                   { return "Kitsu" }
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
	for _, r := range result.Data {
		searchResults = append(searchResults, mapAnimeToTVShowSearchResult(r))
	}
	return searchResults, nil
}

func (p *Provider) GetTVShow(ctx context.Context, id string, _ string) (*metadata.TVShowMetadata, error) {
	result, err := p.client.GetAnime(ctx, id)
	if err != nil {
		return nil, err
	}

	m := mapAnimeToTVShowMetadata(result)
	if m == nil {
		return nil, metadata.ErrNotFound
	}
	return m, nil
}


func (p *Provider) GetTVShowImages(ctx context.Context, id string) (*metadata.Images, error) {
	result, err := p.client.GetAnime(ctx, id)
	if err != nil {
		return nil, err
	}

	images := mapImages(result)
	if images == nil {
		return nil, metadata.ErrNotFound
	}
	return images, nil
}

func (p *Provider) GetTVShowContentRatings(ctx context.Context, id string) ([]metadata.ContentRating, error) {
	result, err := p.client.GetAnime(ctx, id)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, metadata.ErrNotFound
	}

	a := result.Data.Attributes
	if a.AgeRating == nil {
		return nil, metadata.ErrNotFound
	}

	rating := metadata.ContentRating{
		CountryCode: "JP",
		Rating:      *a.AgeRating,
	}
	if a.AgeRatingGuide != nil {
		rating.Descriptors = []string{*a.AgeRatingGuide}
	}

	return []metadata.ContentRating{rating}, nil
}


func (p *Provider) GetTVShowExternalIDs(ctx context.Context, id string) (*metadata.ExternalIDs, error) {
	mappings, err := p.client.GetMappings(ctx, id)
	if err != nil {
		return nil, err
	}

	ids := mapMappingsToExternalIDs(mappings)
	if ids == nil {
		return nil, metadata.ErrNotFound
	}
	return ids, nil
}

func (p *Provider) GetSeason(ctx context.Context, showID string, seasonNum int, _ string) (*metadata.SeasonMetadata, error) {
	// Kitsu anime are typically single-season. Fetch episodes for the show.
	episodes, err := p.client.GetEpisodes(ctx, showID, 20, 0)
	if err != nil {
		return nil, err
	}
	if episodes == nil || len(episodes.Data) == 0 {
		return nil, metadata.ErrNotFound
	}

	sm := &metadata.SeasonMetadata{
		ProviderID:   showID,
		Provider:     metadata.ProviderKitsu,
		ShowID:       showID,
		SeasonNumber: seasonNum,
		Name:         fmt.Sprintf("Season %d", seasonNum),
		Episodes:     mapEpisodesToSummary(episodes, seasonNum),
	}

	return sm, nil
}



func (p *Provider) GetEpisode(ctx context.Context, showID string, seasonNum, episodeNum int, _ string) (*metadata.EpisodeMetadata, error) {
	// Fetch episodes and find the matching one
	// Kitsu limits to 20 per page, so we may need to offset
	offset := 0
	if episodeNum > 20 {
		offset = episodeNum - 20
	}

	episodes, err := p.client.GetEpisodes(ctx, showID, 20, offset)
	if err != nil {
		return nil, err
	}
	if episodes == nil {
		return nil, metadata.ErrNotFound
	}

	for _, ep := range episodes.Data {
		a := ep.Attributes
		sn := 1
		if a.SeasonNumber != nil {
			sn = *a.SeasonNumber
		}
		if a.Number != nil && *a.Number == episodeNum && sn == seasonNum {
			return mapEpisodeToMetadata(ep, showID), nil
		}
	}

	return nil, metadata.ErrNotFound
}


