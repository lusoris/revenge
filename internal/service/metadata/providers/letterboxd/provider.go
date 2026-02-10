package letterboxd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lusoris/revenge/internal/service/metadata"
)

// Ensure Provider implements required interfaces.
// Letterboxd provides movie metadata with community ratings, credits, and
// links to TMDb/IMDb. It is movies-only (no TV show support).
var (
	_ metadata.Provider      = (*Provider)(nil)
	_ metadata.MovieProvider = (*Provider)(nil)
)

// Provider implements the metadata provider interface for Letterboxd.
type Provider struct {
	metadata.MovieProviderBase
	client   *Client
	priority int
}

// NewProvider creates a new Letterboxd provider.
func NewProvider(config Config) (*Provider, error) {
	client, err := NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("create letterboxd provider: %w", err)
	}
	return &Provider{
		client:   client,
		priority: 34, // Below Simkl (36)
	}, nil
}

func (p *Provider) ID() metadata.ProviderID       { return metadata.ProviderLetterboxd }
func (p *Provider) Name() string                   { return "Letterboxd" }
func (p *Provider) Priority() int                  { return p.priority }
func (p *Provider) SupportsMovies() bool           { return true }
func (p *Provider) SupportsTVShows() bool          { return false }
func (p *Provider) SupportsPeople() bool           { return false }
func (p *Provider) SupportsLanguage(_ string) bool { return true } // English-primary
func (p *Provider) ClearCache()                    { p.client.clearCache() }

// --- MovieProvider ---

func (p *Provider) SearchMovie(ctx context.Context, query string, _ metadata.SearchOptions) ([]metadata.MovieSearchResult, error) {
	searchResp, err := p.client.SearchFilms(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("letterboxd search movie: %w", err)
	}

	var results []metadata.MovieSearchResult
	for _, item := range searchResp.Items {
		if item.Type != "FilmSearchItem" || item.Film == nil {
			continue
		}
		results = append(results, mapFilmSummaryToSearchResult(item.Film))
	}

	return results, nil
}

func (p *Provider) GetMovie(ctx context.Context, id string, _ string) (*metadata.MovieMetadata, error) {
	film, err := p.client.GetFilm(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("letterboxd get movie %s: %w", id, err)
	}

	m := mapFilmToMetadata(film)

	// Enrich with statistics (vote count)
	stats, err := p.client.GetFilmStatistics(ctx, id)
	if err == nil && stats != nil {
		m.VoteCount = stats.Counts.Ratings
	}

	return &m, nil
}

func (p *Provider) GetMovieCredits(ctx context.Context, id string) (*metadata.Credits, error) {
	film, err := p.client.GetFilm(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("letterboxd get movie credits %s: %w", id, err)
	}

	credits := mapCredits(film.Contributions)
	return &credits, nil
}




func (p *Provider) GetMovieExternalIDs(ctx context.Context, id string) (*metadata.ExternalIDs, error) {
	film, err := p.client.GetFilm(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("letterboxd get external IDs %s: %w", id, err)
	}

	ids := &metadata.ExternalIDs{}
	for _, link := range film.Links {
		switch link.Type {
		case "imdb":
			ids.IMDbID = &link.ID
		case "tmdb":
			if tmdbID, err := strconv.Atoi(link.ID); err == nil {
				id32 := int32(tmdbID)
				ids.TMDbID = &id32
			}
		}
	}

	return ids, nil
}


