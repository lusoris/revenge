package tvdb

import (
	"context"
	"strconv"

	"github.com/lusoris/revenge/internal/service/metadata"
)

// Ensure Provider implements interfaces.
var (
	_ metadata.Provider       = (*Provider)(nil)
	_ metadata.TVShowProvider = (*Provider)(nil)
	_ metadata.PersonProvider = (*Provider)(nil)
)

// Provider implements the metadata provider interface for TVDb.
type Provider struct {
	client   *Client
	priority int
}

// NewProvider creates a new TVDb provider.
func NewProvider(config Config) *Provider {
	return &Provider{
		client:   NewClient(config),
		priority: 80, // TVDb is secondary to TMDb
	}
}

// NewProviderWithClient creates a provider with an existing client.
func NewProviderWithClient(client *Client) *Provider {
	return &Provider{
		client:   client,
		priority: 80,
	}
}

// ID returns the provider identifier.
func (p *Provider) ID() metadata.ProviderID {
	return metadata.ProviderTVDb
}

// Name returns the human-readable provider name.
func (p *Provider) Name() string {
	return "TheTVDB"
}

// Priority returns the provider priority.
func (p *Provider) Priority() int {
	return p.priority
}

// SupportsMovies returns false (TVDb movie support is limited).
func (p *Provider) SupportsMovies() bool {
	return false // We use TMDb for movies
}

// SupportsTVShows returns true.
func (p *Provider) SupportsTVShows() bool {
	return true
}

// SupportsPeople returns true.
func (p *Provider) SupportsPeople() bool {
	return true
}

// SupportsLanguage returns true for supported languages.
func (p *Provider) SupportsLanguage(lang string) bool {
	return lang != ""
}

// SearchTVShow searches for TV shows.
func (p *Provider) SearchTVShow(ctx context.Context, query string, opts metadata.SearchOptions) ([]metadata.TVShowSearchResult, error) {
	resp, err := p.client.Search(ctx, query, "series")
	if err != nil {
		return nil, err
	}

	results := make([]metadata.TVShowSearchResult, 0, len(resp.Data))
	for _, r := range resp.Data {
		if r.Type == "series" || r.PrimaryType == "series" {
			results = append(results, mapTVSearchResult(&r))
		}
	}
	return results, nil
}

// GetTVShow retrieves TV show details.
func (p *Provider) GetTVShow(ctx context.Context, id string, lang string) (*metadata.TVShowMetadata, error) {
	tvdbID, err := strconv.Atoi(id)
	if err != nil {
		return nil, metadata.NewProviderError(metadata.ProviderTVDb, 400, "invalid id", metadata.ErrInvalidID)
	}

	resp, err := p.client.GetSeriesExtended(ctx, tvdbID, "")
	if err != nil {
		return nil, err
	}

	result := mapTVShowMetadata(&resp.SeriesResponse)

	// Map localized data if available
	if lang != "" && lang != "en" {
		if overview, ok := resp.Overviews[lang]; ok {
			result.Translations = make(map[string]*metadata.LocalizedTVShowData)
			result.Translations[lang] = &metadata.LocalizedTVShowData{
				Language: lang,
				Overview: overview,
			}
		}
	}

	return result, nil
}

// GetTVShowCredits retrieves TV show credits.
func (p *Provider) GetTVShowCredits(ctx context.Context, id string) (*metadata.Credits, error) {
	tvdbID, err := strconv.Atoi(id)
	if err != nil {
		return nil, metadata.NewProviderError(metadata.ProviderTVDb, 400, "invalid id", metadata.ErrInvalidID)
	}

	resp, err := p.client.GetSeriesExtended(ctx, tvdbID, "")
	if err != nil {
		return nil, err
	}

	return mapCharactersToCredits(resp.Characters), nil
}

// GetTVShowImages retrieves TV show images.
func (p *Provider) GetTVShowImages(ctx context.Context, id string) (*metadata.Images, error) {
	tvdbID, err := strconv.Atoi(id)
	if err != nil {
		return nil, metadata.NewProviderError(metadata.ProviderTVDb, 400, "invalid id", metadata.ErrInvalidID)
	}

	artworks, err := p.client.GetSeriesArtworks(ctx, tvdbID, nil, "")
	if err != nil {
		return nil, err
	}

	return mapArtworksToImages(artworks), nil
}

// GetTVShowContentRatings retrieves TV show content ratings.
func (p *Provider) GetTVShowContentRatings(ctx context.Context, id string) ([]metadata.ContentRating, error) {
	tvdbID, err := strconv.Atoi(id)
	if err != nil {
		return nil, metadata.NewProviderError(metadata.ProviderTVDb, 400, "invalid id", metadata.ErrInvalidID)
	}

	resp, err := p.client.GetSeriesExtended(ctx, tvdbID, "")
	if err != nil {
		return nil, err
	}

	return mapContentRatings(resp.ContentRatings), nil
}

// GetTVShowTranslations retrieves TV show translations.
func (p *Provider) GetTVShowTranslations(ctx context.Context, id string) ([]metadata.Translation, error) {
	tvdbID, err := strconv.Atoi(id)
	if err != nil {
		return nil, metadata.NewProviderError(metadata.ProviderTVDb, 400, "invalid id", metadata.ErrInvalidID)
	}

	resp, err := p.client.GetSeriesExtended(ctx, tvdbID, "")
	if err != nil {
		return nil, err
	}

	return mapOverviewsToTranslations(resp.Overviews, resp.NameTranslations), nil
}

// GetTVShowExternalIDs retrieves TV show external IDs.
func (p *Provider) GetTVShowExternalIDs(ctx context.Context, id string) (*metadata.ExternalIDs, error) {
	tvdbID, err := strconv.Atoi(id)
	if err != nil {
		return nil, metadata.NewProviderError(metadata.ProviderTVDb, 400, "invalid id", metadata.ErrInvalidID)
	}

	resp, err := p.client.GetSeriesExtended(ctx, tvdbID, "")
	if err != nil {
		return nil, err
	}

	return mapRemoteIDsToExternalIDs(resp.RemoteIDs, int32(tvdbID)), nil
}

// GetSeason retrieves season details.
func (p *Provider) GetSeason(ctx context.Context, showID string, seasonNum int, lang string) (*metadata.SeasonMetadata, error) {
	tvdbID, err := strconv.Atoi(showID)
	if err != nil {
		return nil, metadata.NewProviderError(metadata.ProviderTVDb, 400, "invalid id", metadata.ErrInvalidID)
	}

	// TVDb requires finding the season ID first
	series, err := p.client.GetSeriesExtended(ctx, tvdbID, "")
	if err != nil {
		return nil, err
	}

	var seasonID int
	for _, s := range series.Seasons {
		if s.Number == seasonNum && (s.Type == nil || s.Type.Type == "default" || s.Type.Type == "official") {
			seasonID = s.ID
			break
		}
	}

	if seasonID == 0 {
		return nil, metadata.NewNotFoundError(metadata.ProviderTVDb, "season", strconv.Itoa(seasonNum))
	}

	resp, err := p.client.GetSeasonExtended(ctx, seasonID)
	if err != nil {
		return nil, err
	}

	return mapSeasonMetadata(resp, showID), nil
}

// GetSeasonCredits retrieves season credits.
func (p *Provider) GetSeasonCredits(ctx context.Context, showID string, seasonNum int) (*metadata.Credits, error) {
	// TVDb doesn't have dedicated season credits - use show credits
	return p.GetTVShowCredits(ctx, showID)
}

// GetSeasonImages retrieves season images.
func (p *Provider) GetSeasonImages(ctx context.Context, showID string, seasonNum int) (*metadata.Images, error) {
	// TVDb doesn't have dedicated season images endpoint
	// Return empty for now - could filter show artworks by season
	return &metadata.Images{}, nil
}

// GetEpisode retrieves episode details.
func (p *Provider) GetEpisode(ctx context.Context, showID string, seasonNum, episodeNum int, lang string) (*metadata.EpisodeMetadata, error) {
	tvdbID, err := strconv.Atoi(showID)
	if err != nil {
		return nil, metadata.NewProviderError(metadata.ProviderTVDb, 400, "invalid id", metadata.ErrInvalidID)
	}

	// Get episodes for the season
	episodes, err := p.client.GetSeriesEpisodes(ctx, tvdbID, "default", &seasonNum, 0)
	if err != nil {
		return nil, err
	}

	// Find the episode
	for _, ep := range episodes {
		if ep.SeasonNumber == seasonNum && ep.Number == episodeNum {
			return mapEpisodeMetadata(&ep, showID), nil
		}
	}

	return nil, metadata.NewNotFoundError(metadata.ProviderTVDb, "episode", strconv.Itoa(episodeNum))
}

// GetEpisodeCredits retrieves episode credits.
func (p *Provider) GetEpisodeCredits(ctx context.Context, showID string, seasonNum, episodeNum int) (*metadata.Credits, error) {
	// TVDb stores credits per episode via characters array
	tvdbID, err := strconv.Atoi(showID)
	if err != nil {
		return nil, metadata.NewProviderError(metadata.ProviderTVDb, 400, "invalid id", metadata.ErrInvalidID)
	}

	episodes, err := p.client.GetSeriesEpisodes(ctx, tvdbID, "default", &seasonNum, 0)
	if err != nil {
		return nil, err
	}

	for _, ep := range episodes {
		if ep.SeasonNumber == seasonNum && ep.Number == episodeNum {
			return mapCharactersToCredits(ep.Characters), nil
		}
	}

	return &metadata.Credits{}, nil
}

// GetEpisodeImages retrieves episode images.
func (p *Provider) GetEpisodeImages(ctx context.Context, showID string, seasonNum, episodeNum int) (*metadata.Images, error) {
	// TVDb episodes have a single image field
	tvdbID, err := strconv.Atoi(showID)
	if err != nil {
		return nil, metadata.NewProviderError(metadata.ProviderTVDb, 400, "invalid id", metadata.ErrInvalidID)
	}

	episodes, err := p.client.GetSeriesEpisodes(ctx, tvdbID, "default", &seasonNum, 0)
	if err != nil {
		return nil, err
	}

	for _, ep := range episodes {
		if ep.SeasonNumber == seasonNum && ep.Number == episodeNum {
			images := &metadata.Images{}
			if ep.Image != nil && *ep.Image != "" {
				images.Stills = []metadata.Image{{FilePath: *ep.Image}}
			}
			return images, nil
		}
	}

	return &metadata.Images{}, nil
}

// SearchPerson searches for people.
func (p *Provider) SearchPerson(ctx context.Context, query string, opts metadata.SearchOptions) ([]metadata.PersonSearchResult, error) {
	resp, err := p.client.Search(ctx, query, "person")
	if err != nil {
		return nil, err
	}

	results := make([]metadata.PersonSearchResult, 0, len(resp.Data))
	for _, r := range resp.Data {
		if r.Type == "person" || r.PrimaryType == "person" {
			results = append(results, mapPersonSearchResult(&r))
		}
	}
	return results, nil
}

// GetPerson retrieves person details.
func (p *Provider) GetPerson(ctx context.Context, id string, lang string) (*metadata.PersonMetadata, error) {
	tvdbID, err := strconv.Atoi(id)
	if err != nil {
		return nil, metadata.NewProviderError(metadata.ProviderTVDb, 400, "invalid id", metadata.ErrInvalidID)
	}

	resp, err := p.client.GetPersonExtended(ctx, tvdbID, "")
	if err != nil {
		return nil, err
	}

	return mapPersonMetadata(resp, lang), nil
}

// GetPersonCredits retrieves person credits.
func (p *Provider) GetPersonCredits(ctx context.Context, id string) (*metadata.PersonCredits, error) {
	tvdbID, err := strconv.Atoi(id)
	if err != nil {
		return nil, metadata.NewProviderError(metadata.ProviderTVDb, 400, "invalid id", metadata.ErrInvalidID)
	}

	resp, err := p.client.GetPersonExtended(ctx, tvdbID, "")
	if err != nil {
		return nil, err
	}

	return mapPersonCredits(resp), nil
}

// GetPersonImages retrieves person images.
func (p *Provider) GetPersonImages(ctx context.Context, id string) (*metadata.Images, error) {
	tvdbID, err := strconv.Atoi(id)
	if err != nil {
		return nil, metadata.NewProviderError(metadata.ProviderTVDb, 400, "invalid id", metadata.ErrInvalidID)
	}

	resp, err := p.client.GetPerson(ctx, tvdbID)
	if err != nil {
		return nil, err
	}

	images := &metadata.Images{}
	if resp.Image != nil && *resp.Image != "" {
		images.Profiles = []metadata.Image{{FilePath: *resp.Image}}
	}
	return images, nil
}

// GetPersonExternalIDs retrieves person external IDs.
func (p *Provider) GetPersonExternalIDs(ctx context.Context, id string) (*metadata.ExternalIDs, error) {
	tvdbID, err := strconv.Atoi(id)
	if err != nil {
		return nil, metadata.NewProviderError(metadata.ProviderTVDb, 400, "invalid id", metadata.ErrInvalidID)
	}

	resp, err := p.client.GetPersonExtended(ctx, tvdbID, "")
	if err != nil {
		return nil, err
	}

	tvdbID32 := int32(tvdbID)
	result := &metadata.ExternalIDs{
		TVDbID: &tvdbID32,
	}

	for _, remote := range resp.RemoteIDs {
		switch remote.Type {
		case RemoteIDTypeIMDb:
			result.IMDbID = &remote.ID
		case RemoteIDTypeTMDb:
			if id, err := strconv.Atoi(remote.ID); err == nil {
				tmdbID := int32(id)
				result.TMDbID = &tmdbID
			}
		}
	}

	return result, nil
}

// ClearCache clears the provider cache.
func (p *Provider) ClearCache() {
	p.client.ClearCache()
}
