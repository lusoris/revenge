package tmdb

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lusoris/revenge/internal/service/metadata"
	"github.com/lusoris/revenge/internal/util"
)

// Ensure Provider implements all interfaces.
var (
	_ metadata.Provider           = (*Provider)(nil)
	_ metadata.MovieProvider      = (*Provider)(nil)
	_ metadata.TVShowProvider     = (*Provider)(nil)
	_ metadata.PersonProvider     = (*Provider)(nil)
	_ metadata.ImageProvider      = (*Provider)(nil)
	_ metadata.CollectionProvider = (*Provider)(nil)
)

// Provider implements the metadata provider interface for TMDb.
type Provider struct {
	client   *Client
	priority int
}

// NewProvider creates a new TMDb provider.
func NewProvider(config Config) (*Provider, error) {
	client, err := NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("create tmdb provider: %w", err)
	}
	return &Provider{
		client:   client,
		priority: 100, // TMDb is the primary provider
	}, nil
}

// NewProviderWithClient creates a provider with an existing client.
func NewProviderWithClient(client *Client) *Provider {
	return &Provider{
		client:   client,
		priority: 100,
	}
}

// ID returns the provider identifier.
func (p *Provider) ID() metadata.ProviderID {
	return metadata.ProviderTMDb
}

// Name returns the human-readable provider name.
func (p *Provider) Name() string {
	return "The Movie Database"
}

// Priority returns the provider priority.
func (p *Provider) Priority() int {
	return p.priority
}

// SupportsMovies returns true.
func (p *Provider) SupportsMovies() bool {
	return true
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
	// TMDb supports virtually all languages
	return lang != ""
}

// SearchMovie searches for movies.
func (p *Provider) SearchMovie(ctx context.Context, query string, opts metadata.SearchOptions) ([]metadata.MovieSearchResult, error) {
	lang := normalizeLang(opts.Language)
	resp, err := p.client.SearchMovie(ctx, query, opts.Year, lang)
	if err != nil {
		return nil, err
	}

	results := make([]metadata.MovieSearchResult, len(resp.Results))
	for i, r := range resp.Results {
		results[i] = mapMovieSearchResult(&r)
	}
	return results, nil
}

// GetMovie retrieves movie details.
func (p *Provider) GetMovie(ctx context.Context, id string, lang string) (*metadata.MovieMetadata, error) {
	tmdbID, err := strconv.Atoi(id)
	if err != nil {
		return nil, metadata.NewProviderError(metadata.ProviderTMDb, 400, "invalid id", metadata.ErrInvalidID)
	}

	resp, err := p.client.GetMovie(ctx, tmdbID, normalizeLang(lang), "")
	if err != nil {
		return nil, err
	}

	return mapMovieMetadata(resp), nil
}

// GetMovieCredits retrieves movie credits.
func (p *Provider) GetMovieCredits(ctx context.Context, id string) (*metadata.Credits, error) {
	tmdbID, err := strconv.Atoi(id)
	if err != nil {
		return nil, metadata.NewProviderError(metadata.ProviderTMDb, 400, "invalid id", metadata.ErrInvalidID)
	}

	resp, err := p.client.GetMovieCredits(ctx, tmdbID)
	if err != nil {
		return nil, err
	}

	return mapCredits(resp), nil
}

// GetMovieImages retrieves movie images.
func (p *Provider) GetMovieImages(ctx context.Context, id string) (*metadata.Images, error) {
	tmdbID, err := strconv.Atoi(id)
	if err != nil {
		return nil, metadata.NewProviderError(metadata.ProviderTMDb, 400, "invalid id", metadata.ErrInvalidID)
	}

	resp, err := p.client.GetMovieImages(ctx, tmdbID)
	if err != nil {
		return nil, err
	}

	return mapImages(resp), nil
}

// GetMovieReleaseDates retrieves movie release dates.
func (p *Provider) GetMovieReleaseDates(ctx context.Context, id string) ([]metadata.ReleaseDate, error) {
	tmdbID, err := strconv.Atoi(id)
	if err != nil {
		return nil, metadata.NewProviderError(metadata.ProviderTMDb, 400, "invalid id", metadata.ErrInvalidID)
	}

	resp, err := p.client.GetMovieReleaseDates(ctx, tmdbID)
	if err != nil {
		return nil, err
	}

	return mapReleaseDates(resp), nil
}

// GetMovieTranslations retrieves movie translations.
func (p *Provider) GetMovieTranslations(ctx context.Context, id string) ([]metadata.Translation, error) {
	tmdbID, err := strconv.Atoi(id)
	if err != nil {
		return nil, metadata.NewProviderError(metadata.ProviderTMDb, 400, "invalid id", metadata.ErrInvalidID)
	}

	resp, err := p.client.GetMovieTranslations(ctx, tmdbID)
	if err != nil {
		return nil, err
	}

	return mapTranslations(resp), nil
}

// GetMovieExternalIDs retrieves movie external IDs.
func (p *Provider) GetMovieExternalIDs(ctx context.Context, id string) (*metadata.ExternalIDs, error) {
	tmdbID, err := strconv.Atoi(id)
	if err != nil {
		return nil, metadata.NewProviderError(metadata.ProviderTMDb, 400, "invalid id", metadata.ErrInvalidID)
	}

	resp, err := p.client.GetMovieExternalIDs(ctx, tmdbID)
	if err != nil {
		return nil, err
	}

	return mapExternalIDs(resp, util.SafeIntToInt32(tmdbID)), nil
}

// GetSimilarMovies retrieves movies similar to the given movie.
func (p *Provider) GetSimilarMovies(ctx context.Context, id string, opts metadata.SearchOptions) ([]metadata.MovieSearchResult, int, error) {
	tmdbID, err := strconv.Atoi(id)
	if err != nil {
		return nil, 0, metadata.NewProviderError(metadata.ProviderTMDb, 400, "invalid id", metadata.ErrInvalidID)
	}

	page := 1
	if opts.Page > 0 {
		page = opts.Page
	}

	resp, err := p.client.GetSimilarMovies(ctx, tmdbID, normalizeLang(opts.Language), page)
	if err != nil {
		return nil, 0, err
	}

	return mapMovieSearchResults(resp), resp.TotalResults, nil
}

// GetMovieRecommendations retrieves recommended movies based on the given movie.
func (p *Provider) GetMovieRecommendations(ctx context.Context, id string, opts metadata.SearchOptions) ([]metadata.MovieSearchResult, int, error) {
	tmdbID, err := strconv.Atoi(id)
	if err != nil {
		return nil, 0, metadata.NewProviderError(metadata.ProviderTMDb, 400, "invalid id", metadata.ErrInvalidID)
	}

	page := 1
	if opts.Page > 0 {
		page = opts.Page
	}

	resp, err := p.client.GetMovieRecommendations(ctx, tmdbID, normalizeLang(opts.Language), page)
	if err != nil {
		return nil, 0, err
	}

	return mapMovieSearchResults(resp), resp.TotalResults, nil
}

// SearchTVShow searches for TV shows.
func (p *Provider) SearchTVShow(ctx context.Context, query string, opts metadata.SearchOptions) ([]metadata.TVShowSearchResult, error) {
	lang := normalizeLang(opts.Language)
	resp, err := p.client.SearchTV(ctx, query, opts.Year, lang)
	if err != nil {
		return nil, err
	}

	results := make([]metadata.TVShowSearchResult, len(resp.Results))
	for i, r := range resp.Results {
		results[i] = mapTVSearchResult(&r)
	}
	return results, nil
}

// GetTVShow retrieves TV show details.
func (p *Provider) GetTVShow(ctx context.Context, id string, lang string) (*metadata.TVShowMetadata, error) {
	tmdbID, err := strconv.Atoi(id)
	if err != nil {
		return nil, metadata.NewProviderError(metadata.ProviderTMDb, 400, "invalid id", metadata.ErrInvalidID)
	}

	resp, err := p.client.GetTV(ctx, tmdbID, normalizeLang(lang), "")
	if err != nil {
		return nil, err
	}

	return mapTVShowMetadata(resp), nil
}

// GetTVShowCredits retrieves TV show credits.
func (p *Provider) GetTVShowCredits(ctx context.Context, id string) (*metadata.Credits, error) {
	tmdbID, err := strconv.Atoi(id)
	if err != nil {
		return nil, metadata.NewProviderError(metadata.ProviderTMDb, 400, "invalid id", metadata.ErrInvalidID)
	}

	resp, err := p.client.GetTVCredits(ctx, tmdbID)
	if err != nil {
		return nil, err
	}

	return mapCredits(resp), nil
}

// GetTVShowImages retrieves TV show images.
func (p *Provider) GetTVShowImages(ctx context.Context, id string) (*metadata.Images, error) {
	tmdbID, err := strconv.Atoi(id)
	if err != nil {
		return nil, metadata.NewProviderError(metadata.ProviderTMDb, 400, "invalid id", metadata.ErrInvalidID)
	}

	resp, err := p.client.GetTVImages(ctx, tmdbID)
	if err != nil {
		return nil, err
	}

	return mapImages(resp), nil
}

// GetTVShowContentRatings retrieves TV show content ratings.
func (p *Provider) GetTVShowContentRatings(ctx context.Context, id string) ([]metadata.ContentRating, error) {
	tmdbID, err := strconv.Atoi(id)
	if err != nil {
		return nil, metadata.NewProviderError(metadata.ProviderTMDb, 400, "invalid id", metadata.ErrInvalidID)
	}

	resp, err := p.client.GetTVContentRatings(ctx, tmdbID)
	if err != nil {
		return nil, err
	}

	return mapContentRatings(resp), nil
}

// GetTVShowTranslations retrieves TV show translations.
func (p *Provider) GetTVShowTranslations(ctx context.Context, id string) ([]metadata.Translation, error) {
	tmdbID, err := strconv.Atoi(id)
	if err != nil {
		return nil, metadata.NewProviderError(metadata.ProviderTMDb, 400, "invalid id", metadata.ErrInvalidID)
	}

	resp, err := p.client.GetTVTranslations(ctx, tmdbID)
	if err != nil {
		return nil, err
	}

	return mapTranslations(resp), nil
}

// GetTVShowExternalIDs retrieves TV show external IDs.
func (p *Provider) GetTVShowExternalIDs(ctx context.Context, id string) (*metadata.ExternalIDs, error) {
	tmdbID, err := strconv.Atoi(id)
	if err != nil {
		return nil, metadata.NewProviderError(metadata.ProviderTMDb, 400, "invalid id", metadata.ErrInvalidID)
	}

	resp, err := p.client.GetTVExternalIDs(ctx, tmdbID)
	if err != nil {
		return nil, err
	}

	return mapExternalIDs(resp, util.SafeIntToInt32(tmdbID)), nil
}

// GetSeason retrieves season details.
func (p *Provider) GetSeason(ctx context.Context, showID string, seasonNum int, lang string) (*metadata.SeasonMetadata, error) {
	tmdbID, err := strconv.Atoi(showID)
	if err != nil {
		return nil, metadata.NewProviderError(metadata.ProviderTMDb, 400, "invalid id", metadata.ErrInvalidID)
	}

	resp, err := p.client.GetSeason(ctx, tmdbID, seasonNum, normalizeLang(lang), "")
	if err != nil {
		return nil, err
	}

	return mapSeasonMetadata(resp, showID), nil
}

// GetSeasonCredits retrieves season credits.
func (p *Provider) GetSeasonCredits(ctx context.Context, showID string, seasonNum int) (*metadata.Credits, error) {
	tmdbID, err := strconv.Atoi(showID)
	if err != nil {
		return nil, metadata.NewProviderError(metadata.ProviderTMDb, 400, "invalid id", metadata.ErrInvalidID)
	}

	resp, err := p.client.GetSeasonCredits(ctx, tmdbID, seasonNum)
	if err != nil {
		return nil, err
	}

	return mapCredits(resp), nil
}

// GetSeasonImages retrieves season images.
func (p *Provider) GetSeasonImages(ctx context.Context, showID string, seasonNum int) (*metadata.Images, error) {
	tmdbID, err := strconv.Atoi(showID)
	if err != nil {
		return nil, metadata.NewProviderError(metadata.ProviderTMDb, 400, "invalid id", metadata.ErrInvalidID)
	}

	resp, err := p.client.GetSeasonImages(ctx, tmdbID, seasonNum)
	if err != nil {
		return nil, err
	}

	return mapImages(resp), nil
}

// GetEpisode retrieves episode details.
func (p *Provider) GetEpisode(ctx context.Context, showID string, seasonNum, episodeNum int, lang string) (*metadata.EpisodeMetadata, error) {
	tmdbID, err := strconv.Atoi(showID)
	if err != nil {
		return nil, metadata.NewProviderError(metadata.ProviderTMDb, 400, "invalid id", metadata.ErrInvalidID)
	}

	resp, err := p.client.GetEpisode(ctx, tmdbID, seasonNum, episodeNum, normalizeLang(lang), "")
	if err != nil {
		return nil, err
	}

	return mapEpisodeMetadata(resp, showID), nil
}

// GetEpisodeCredits retrieves episode credits.
func (p *Provider) GetEpisodeCredits(ctx context.Context, showID string, seasonNum, episodeNum int) (*metadata.Credits, error) {
	tmdbID, err := strconv.Atoi(showID)
	if err != nil {
		return nil, metadata.NewProviderError(metadata.ProviderTMDb, 400, "invalid id", metadata.ErrInvalidID)
	}

	resp, err := p.client.GetEpisodeCredits(ctx, tmdbID, seasonNum, episodeNum)
	if err != nil {
		return nil, err
	}

	return mapCredits(resp), nil
}

// GetEpisodeImages retrieves episode images.
func (p *Provider) GetEpisodeImages(ctx context.Context, showID string, seasonNum, episodeNum int) (*metadata.Images, error) {
	tmdbID, err := strconv.Atoi(showID)
	if err != nil {
		return nil, metadata.NewProviderError(metadata.ProviderTMDb, 400, "invalid id", metadata.ErrInvalidID)
	}

	resp, err := p.client.GetEpisodeImages(ctx, tmdbID, seasonNum, episodeNum)
	if err != nil {
		return nil, err
	}

	return mapImages(resp), nil
}

// SearchPerson searches for people.
func (p *Provider) SearchPerson(ctx context.Context, query string, opts metadata.SearchOptions) ([]metadata.PersonSearchResult, error) {
	lang := normalizeLang(opts.Language)
	resp, err := p.client.SearchPerson(ctx, query, lang)
	if err != nil {
		return nil, err
	}

	results := make([]metadata.PersonSearchResult, len(resp.Results))
	for i, r := range resp.Results {
		results[i] = mapPersonSearchResult(&r)
	}
	return results, nil
}

// GetPerson retrieves person details.
func (p *Provider) GetPerson(ctx context.Context, id string, lang string) (*metadata.PersonMetadata, error) {
	tmdbID, err := strconv.Atoi(id)
	if err != nil {
		return nil, metadata.NewProviderError(metadata.ProviderTMDb, 400, "invalid id", metadata.ErrInvalidID)
	}

	resp, err := p.client.GetPerson(ctx, tmdbID, normalizeLang(lang), "")
	if err != nil {
		return nil, err
	}

	return mapPersonMetadata(resp), nil
}

// GetPersonCredits retrieves person credits.
func (p *Provider) GetPersonCredits(ctx context.Context, id string) (*metadata.PersonCredits, error) {
	tmdbID, err := strconv.Atoi(id)
	if err != nil {
		return nil, metadata.NewProviderError(metadata.ProviderTMDb, 400, "invalid id", metadata.ErrInvalidID)
	}

	resp, err := p.client.GetPersonCredits(ctx, tmdbID, "")
	if err != nil {
		return nil, err
	}

	return mapPersonCredits(resp), nil
}

// GetPersonImages retrieves person images.
func (p *Provider) GetPersonImages(ctx context.Context, id string) (*metadata.Images, error) {
	tmdbID, err := strconv.Atoi(id)
	if err != nil {
		return nil, metadata.NewProviderError(metadata.ProviderTMDb, 400, "invalid id", metadata.ErrInvalidID)
	}

	resp, err := p.client.GetPersonImages(ctx, tmdbID)
	if err != nil {
		return nil, err
	}

	return mapPersonImages(resp), nil
}

// GetPersonExternalIDs retrieves person external IDs.
func (p *Provider) GetPersonExternalIDs(ctx context.Context, id string) (*metadata.ExternalIDs, error) {
	tmdbID, err := strconv.Atoi(id)
	if err != nil {
		return nil, metadata.NewProviderError(metadata.ProviderTMDb, 400, "invalid id", metadata.ErrInvalidID)
	}

	resp, err := p.client.GetPersonExternalIDs(ctx, tmdbID)
	if err != nil {
		return nil, err
	}

	return mapExternalIDs(resp, util.SafeIntToInt32(tmdbID)), nil
}

// GetImageURL constructs a full image URL.
func (p *Provider) GetImageURL(path string, size metadata.ImageSize) string {
	return p.client.GetImageURL(path, string(size))
}

// GetImageBaseURL returns the base URL for images.
func (p *Provider) GetImageBaseURL() string {
	return ImageBaseURL
}

// DownloadImage downloads an image.
func (p *Provider) DownloadImage(ctx context.Context, path string, size metadata.ImageSize) ([]byte, error) {
	return p.client.DownloadImage(ctx, path, string(size))
}

// GetCollection retrieves collection details.
func (p *Provider) GetCollection(ctx context.Context, id string, lang string) (*metadata.CollectionMetadata, error) {
	collID, err := strconv.Atoi(id)
	if err != nil {
		return nil, metadata.NewProviderError(metadata.ProviderTMDb, 400, "invalid id", metadata.ErrInvalidID)
	}

	resp, err := p.client.GetCollection(ctx, collID, normalizeLang(lang))
	if err != nil {
		return nil, err
	}

	return mapCollectionMetadata(resp), nil
}

// GetCollectionImages retrieves collection images.
func (p *Provider) GetCollectionImages(ctx context.Context, id string) (*metadata.Images, error) {
	// TMDb doesn't have a dedicated collection images endpoint
	// We would need to fetch via append_to_response
	collID, err := strconv.Atoi(id)
	if err != nil {
		return nil, metadata.NewProviderError(metadata.ProviderTMDb, 400, "invalid id", metadata.ErrInvalidID)
	}

	// For now, return empty images - could be enhanced later
	_ = collID
	return &metadata.Images{}, nil
}

// ClearCache clears the provider cache.
func (p *Provider) ClearCache() {
	p.client.ClearCache()
}

// normalizeLang converts ISO 639-1 (en) to TMDb format (en-US) if needed.
func normalizeLang(lang string) string {
	if lang == "" {
		return "en-US"
	}
	if len(lang) == 2 {
		switch lang {
		case "en":
			return "en-US"
		case "de":
			return "de-DE"
		case "fr":
			return "fr-FR"
		case "es":
			return "es-ES"
		case "it":
			return "it-IT"
		case "pt":
			return "pt-BR"
		case "ja":
			return "ja-JP"
		case "ko":
			return "ko-KR"
		case "zh":
			return "zh-CN"
		case "ru":
			return "ru-RU"
		default:
			return fmt.Sprintf("%s-%s", lang, lang)
		}
	}
	return lang
}
