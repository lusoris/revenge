package simkl

import (
	"context"
	"fmt"

	"github.com/lusoris/revenge/internal/service/metadata"
)

// Ensure Provider implements required interfaces.
// Simkl provides movie, TV show, and anime metadata with cross-referenced IDs
// (IMDb, TMDb, TVDb, MAL, AniDB, AniList, Kitsu), community ratings,
// and external ratings (IMDb, MAL). It also provides images (posters, fanart).
// Note: Simkl API terms require linking back to simkl.com and restrict use
// with competing services unless Simkl sync is integrated.
var (
	_ metadata.Provider       = (*Provider)(nil)
	_ metadata.MovieProvider  = (*Provider)(nil)
	_ metadata.TVShowProvider = (*Provider)(nil)
)

// Provider implements the metadata provider interface for Simkl.
type Provider struct {
	metadata.MovieProviderBase
	metadata.TVShowProviderBase
	client   *Client
	priority int
}

// NewProvider creates a new Simkl provider.
func NewProvider(config Config) (*Provider, error) {
	client, err := NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("create simkl provider: %w", err)
	}
	return &Provider{
		client:   client,
		priority: 36, // Below Trakt (38), above Letterboxd (34)
	}, nil
}

func (p *Provider) ID() metadata.ProviderID       { return metadata.ProviderSimkl }
func (p *Provider) Name() string                   { return "Simkl" }
func (p *Provider) Priority() int                  { return p.priority }
func (p *Provider) SupportsMovies() bool           { return true }
func (p *Provider) SupportsTVShows() bool          { return true }
func (p *Provider) SupportsPeople() bool           { return false }
func (p *Provider) SupportsLanguage(_ string) bool { return true } // English-primary
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
		searchResults = append(searchResults, mapSearchResultToMovieSearchResult(r))
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

func (p *Provider) GetMovieCredits(_ context.Context, _ string) (*metadata.Credits, error) {
	// Simkl does not provide credits
	return nil, metadata.ErrNotFound
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
	// Search both TV and anime
	tvResults, tvErr := p.client.SearchShows(ctx, query)
	animeResults, animeErr := p.client.SearchAnime(ctx, query)

	if tvErr != nil && animeErr != nil {
		return nil, tvErr
	}

	var allResults []SearchResult
	allResults = append(allResults, tvResults...)
	allResults = append(allResults, animeResults...)

	if len(allResults) == 0 {
		return nil, metadata.ErrNotFound
	}

	searchResults := make([]metadata.TVShowSearchResult, 0, len(allResults))
	for _, r := range allResults {
		searchResults = append(searchResults, mapSearchResultToTVShowSearchResult(r))
	}
	return searchResults, nil
}

func (p *Provider) GetTVShow(ctx context.Context, id string, _ string) (*metadata.TVShowMetadata, error) {
	// Try TV first, then anime
	show, err := p.client.GetShow(ctx, id)
	if err != nil {
		show, err = p.client.GetAnime(ctx, id)
		if err != nil {
			return nil, err
		}
	}

	md := mapShowToMetadata(show)
	if md == nil {
		return nil, metadata.ErrNotFound
	}

	// Fetch episodes to count seasons
	episodes, err := p.client.GetShowEpisodes(ctx, id)
	if err == nil && len(episodes) > 0 {
		seasonMap := make(map[int]bool)
		for _, ep := range episodes {
			seasonMap[ep.Season] = true
		}
		md.NumberOfSeasons = len(seasonMap)
		if md.NumberOfEpisodes == 0 {
			md.NumberOfEpisodes = len(episodes)
		}
	}

	return md, nil
}





func (p *Provider) GetTVShowExternalIDs(ctx context.Context, id string) (*metadata.ExternalIDs, error) {
	show, err := p.client.GetShow(ctx, id)
	if err != nil {
		show, err = p.client.GetAnime(ctx, id)
		if err != nil {
			return nil, err
		}
	}
	return mapExternalIDs(show.IDs), nil
}

func (p *Provider) GetSeason(ctx context.Context, showID string, seasonNum int, _ string) (*metadata.SeasonMetadata, error) {
	// Get all episodes and filter by season
	episodes, err := p.client.GetShowEpisodes(ctx, showID)
	if err != nil {
		episodes, err = p.client.GetAnimeEpisodes(ctx, showID)
		if err != nil {
			return nil, err
		}
	}

	seasonEpisodes := mapEpisodesToSummaries(episodes, seasonNum)
	if len(seasonEpisodes) == 0 {
		return nil, metadata.ErrNotFound
	}

	sm := &metadata.SeasonMetadata{
		ProviderID:   fmt.Sprintf("%s-s%d", showID, seasonNum),
		Provider:     metadata.ProviderSimkl,
		ShowID:       showID,
		SeasonNumber: seasonNum,
		Episodes:     seasonEpisodes,
	}

	// Set air date from first episode
	if len(seasonEpisodes) > 0 && seasonEpisodes[0].AirDate != nil {
		sm.AirDate = seasonEpisodes[0].AirDate
	}

	return sm, nil
}





