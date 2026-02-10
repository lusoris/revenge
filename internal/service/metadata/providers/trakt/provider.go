package trakt

import (
	"context"
	"fmt"

	"github.com/lusoris/revenge/internal/service/metadata"
)

// Ensure Provider implements required interfaces.
// Trakt provides rich movie and TV show metadata with cross-referenced IDs
// (IMDb, TMDb, TVDb), credits, translations, and community ratings.
// It does NOT support images, people search, or collections.
var (
	_ metadata.Provider       = (*Provider)(nil)
	_ metadata.MovieProvider  = (*Provider)(nil)
	_ metadata.TVShowProvider = (*Provider)(nil)
)

// Provider implements the metadata provider interface for Trakt.
type Provider struct {
	metadata.MovieProviderBase
	metadata.TVShowProviderBase
	client   *Client
	priority int
}

// NewProvider creates a new Trakt provider.
func NewProvider(config Config) (*Provider, error) {
	client, err := NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("create trakt provider: %w", err)
	}
	return &Provider{
		client:   client,
		priority: 38, // Below anime providers (42-45), above Simkl/Letterboxd
	}, nil
}

func (p *Provider) ID() metadata.ProviderID       { return metadata.ProviderTrakt }
func (p *Provider) Name() string                   { return "Trakt" }
func (p *Provider) Priority() int                  { return p.priority }
func (p *Provider) SupportsMovies() bool           { return true }
func (p *Provider) SupportsTVShows() bool          { return true }
func (p *Provider) SupportsPeople() bool           { return false }
func (p *Provider) SupportsLanguage(_ string) bool { return true } // English-primary, multilingual
func (p *Provider) ClearCache()                    { p.client.clearCache() }

// --- MovieProvider ---

func (p *Provider) SearchMovie(ctx context.Context, query string, _ metadata.SearchOptions) ([]metadata.MovieSearchResult, error) {
	results, err := p.client.SearchMovies(ctx, query)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, metadata.ErrNotFound
	}

	searchResults := make([]metadata.MovieSearchResult, 0, len(results))
	for _, r := range results {
		if r.Movie != nil {
			searchResults = append(searchResults, mapMovieToSearchResult(r.Movie))
		}
	}
	if len(searchResults) == 0 {
		return nil, metadata.ErrNotFound
	}
	return searchResults, nil
}

func (p *Provider) GetMovie(ctx context.Context, id string, _ string) (*metadata.MovieMetadata, error) {
	movie, err := p.client.GetMovie(ctx, id)
	if err != nil {
		return nil, err
	}
	md := mapMovieToMetadata(movie)
	if md == nil {
		return nil, metadata.ErrNotFound
	}
	return md, nil
}

func (p *Provider) GetMovieCredits(ctx context.Context, id string) (*metadata.Credits, error) {
	credits, err := p.client.GetMovieCredits(ctx, id)
	if err != nil {
		return nil, err
	}
	result := mapCredits(credits)
	if result == nil || (len(result.Cast) == 0 && len(result.Crew) == 0) {
		return nil, metadata.ErrNotFound
	}
	return result, nil
}


func (p *Provider) GetMovieReleaseDates(_ context.Context, _ string) ([]metadata.ReleaseDate, error) {
	// Trakt does not provide detailed release dates
	return nil, metadata.ErrNotFound
}

func (p *Provider) GetMovieTranslations(ctx context.Context, id string) ([]metadata.Translation, error) {
	translations, err := p.client.GetMovieTranslations(ctx, id)
	if err != nil {
		return nil, err
	}
	result := mapTranslations(translations)
	if len(result) == 0 {
		return nil, metadata.ErrNotFound
	}
	return result, nil
}

func (p *Provider) GetMovieExternalIDs(ctx context.Context, id string) (*metadata.ExternalIDs, error) {
	movie, err := p.client.GetMovie(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapExternalIDs(movie.IDs), nil
}



// --- TVShowProvider ---

func (p *Provider) SearchTVShow(ctx context.Context, query string, _ metadata.SearchOptions) ([]metadata.TVShowSearchResult, error) {
	results, err := p.client.SearchShows(ctx, query)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, metadata.ErrNotFound
	}

	searchResults := make([]metadata.TVShowSearchResult, 0, len(results))
	for _, r := range results {
		if r.Show != nil {
			searchResults = append(searchResults, mapShowToSearchResult(r.Show))
		}
	}
	if len(searchResults) == 0 {
		return nil, metadata.ErrNotFound
	}
	return searchResults, nil
}

func (p *Provider) GetTVShow(ctx context.Context, id string, _ string) (*metadata.TVShowMetadata, error) {
	show, err := p.client.GetShow(ctx, id)
	if err != nil {
		return nil, err
	}
	md := mapShowToMetadata(show)
	if md == nil {
		return nil, metadata.ErrNotFound
	}

	// Fetch seasons
	seasons, err := p.client.GetShowSeasons(ctx, id)
	if err == nil && len(seasons) > 0 {
		md.Seasons = mapSeasons(seasons)
		md.NumberOfSeasons = len(seasons)

		totalEpisodes := 0
		for _, s := range seasons {
			totalEpisodes += s.EpisodeCount
		}
		md.NumberOfEpisodes = totalEpisodes
	}

	return md, nil
}

func (p *Provider) GetTVShowCredits(ctx context.Context, id string) (*metadata.Credits, error) {
	credits, err := p.client.GetShowCredits(ctx, id)
	if err != nil {
		return nil, err
	}
	result := mapCredits(credits)
	if result == nil || (len(result.Cast) == 0 && len(result.Crew) == 0) {
		return nil, metadata.ErrNotFound
	}
	return result, nil
}



func (p *Provider) GetTVShowTranslations(ctx context.Context, id string) ([]metadata.Translation, error) {
	translations, err := p.client.GetShowTranslations(ctx, id)
	if err != nil {
		return nil, err
	}
	result := mapTranslations(translations)
	if len(result) == 0 {
		return nil, metadata.ErrNotFound
	}
	return result, nil
}

func (p *Provider) GetTVShowExternalIDs(ctx context.Context, id string) (*metadata.ExternalIDs, error) {
	show, err := p.client.GetShow(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapExternalIDs(show.IDs), nil
}

func (p *Provider) GetSeason(ctx context.Context, showID string, seasonNum int, _ string) (*metadata.SeasonMetadata, error) {
	episodes, err := p.client.GetSeasonEpisodes(ctx, showID, seasonNum)
	if err != nil {
		return nil, err
	}

	// Also get season metadata from seasons list
	seasons, err := p.client.GetShowSeasons(ctx, showID)
	if err != nil {
		return nil, err
	}

	var matchedSeason *Season
	for i := range seasons {
		if seasons[i].Number == seasonNum {
			matchedSeason = &seasons[i]
			break
		}
	}
	if matchedSeason == nil {
		return nil, metadata.ErrNotFound
	}

	sm := &metadata.SeasonMetadata{
		ProviderID:   fmt.Sprintf("%d", matchedSeason.IDs.Trakt),
		Provider:     metadata.ProviderTrakt,
		ShowID:       showID,
		SeasonNumber: matchedSeason.Number,
		Name:         matchedSeason.Title,
		AirDate:      matchedSeason.FirstAired,
		VoteAverage:  matchedSeason.Rating,
	}
	if matchedSeason.Overview != "" {
		sm.Overview = &matchedSeason.Overview
	}
	sm.Episodes = mapEpisodesToSummaries(episodes)
	return sm, nil
}





