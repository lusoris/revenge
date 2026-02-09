// Package movie provides an adapter that bridges the shared metadata service
// to the movie content module.
package movie

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/govalues/decimal"

	contentmovie "github.com/lusoris/revenge/internal/content/movie"
	"github.com/lusoris/revenge/internal/service/metadata"
)

// Adapter wraps the shared metadata service for movie-specific operations.
// This adapter implements the movie.MetadataProvider interface using the shared service.
type Adapter struct {
	service   metadata.Service
	languages []string
}

// NewAdapter creates a new adapter that uses the shared metadata service.
func NewAdapter(service metadata.Service, languages []string) *Adapter {
	if len(languages) == 0 {
		languages = []string{"en", "de", "fr", "es", "ja"}
	}
	return &Adapter{
		service:   service,
		languages: languages,
	}
}

// Ensure Adapter implements MetadataProvider.
var _ contentmovie.MetadataProvider = (*Adapter)(nil)

// SearchMovies searches for movies using the shared metadata service.
func (a *Adapter) SearchMovies(ctx context.Context, query string, year *int) ([]*contentmovie.Movie, error) {
	opts := metadata.SearchOptions{
		Year:     year,
		Language: a.languages[0],
	}

	results, err := a.service.SearchMovie(ctx, query, opts)
	if err != nil {
		return nil, fmt.Errorf("search movies: %w", err)
	}

	movies := make([]*contentmovie.Movie, len(results))
	for i, r := range results {
		movies[i] = mapSearchResultToMovie(&r)
	}

	return movies, nil
}

// EnrichMovie enriches a movie with metadata from the shared service.
func (a *Adapter) EnrichMovie(ctx context.Context, mov *contentmovie.Movie, opts ...contentmovie.MetadataRefreshOptions) error {
	if mov.TMDbID == nil {
		return fmt.Errorf("movie has no TMDb ID")
	}

	// Determine languages and force from options
	languages := a.languages
	if len(opts) > 0 {
		if len(opts[0].Languages) > 0 {
			languages = opts[0].Languages
		}
		if opts[0].Force {
			a.service.ClearCache()
		}
	}

	meta, err := a.service.GetMovieMetadata(ctx, *mov.TMDbID, languages)
	if err != nil {
		return fmt.Errorf("get movie metadata: %w", err)
	}

	// Enrich with external ratings from secondary providers (OMDb, etc.)
	// Safe to call here — this runs inside a River background worker, not in API request path.
	a.service.EnrichMovieRatings(ctx, meta)

	// Get release dates for age ratings
	releaseDates, err := a.service.GetMovieReleaseDates(ctx, *mov.TMDbID)
	if err != nil {
		// Continue without release dates
		releaseDates = nil
	}

	// Map to movie domain type
	mapMetadataToMovie(mov, meta, releaseDates)

	return nil
}

// GetMovieCredits retrieves movie credits using the shared service.
func (a *Adapter) GetMovieCredits(ctx context.Context, movieID uuid.UUID, tmdbID int) ([]contentmovie.MovieCredit, error) {
	credits, err := a.service.GetMovieCredits(ctx, int32(tmdbID))
	if err != nil {
		return nil, fmt.Errorf("get movie credits: %w", err)
	}

	return mapCreditsToMovieCredits(movieID, credits), nil
}

// GetMovieGenres retrieves movie genres using the shared service.
func (a *Adapter) GetMovieGenres(ctx context.Context, movieID uuid.UUID, tmdbID int) ([]contentmovie.MovieGenre, error) {
	meta, err := a.service.GetMovieMetadata(ctx, int32(tmdbID), []string{a.languages[0]})
	if err != nil {
		return nil, fmt.Errorf("get movie metadata for genres: %w", err)
	}

	genres := make([]contentmovie.MovieGenre, len(meta.Genres))
	for i, g := range meta.Genres {
		genres[i] = contentmovie.MovieGenre{
			ID:          uuid.Must(uuid.NewV7()),
			MovieID:     movieID,
			TMDbGenreID: int32(g.ID),
			Name:        g.Name,
		}
	}

	return genres, nil
}

// GetMovieByTMDbID retrieves a movie by TMDb ID.
func (a *Adapter) GetMovieByTMDbID(ctx context.Context, tmdbID int) (*contentmovie.Movie, error) {
	meta, err := a.service.GetMovieMetadata(ctx, int32(tmdbID), a.languages)
	if err != nil {
		return nil, fmt.Errorf("get movie metadata: %w", err)
	}

	mov := &contentmovie.Movie{
		ID: uuid.Must(uuid.NewV7()),
	}
	mapMetadataToMovie(mov, meta, nil)
	return mov, nil
}

// GetMovieImages retrieves movie images using the shared service.
func (a *Adapter) GetMovieImages(ctx context.Context, tmdbID int) (*metadata.Images, error) {
	return a.service.GetMovieImages(ctx, int32(tmdbID))
}

// GetImageURL constructs an image URL.
func (a *Adapter) GetImageURL(path string, size metadata.ImageSize) string {
	return a.service.GetImageURL(path, size)
}

// ClearCache clears all cached metadata by delegating to the shared service.
func (a *Adapter) ClearCache() {
	a.service.ClearCache()
}

// mapSearchResultToMovie converts a search result to a movie domain type.
func mapSearchResultToMovie(r *metadata.MovieSearchResult) *contentmovie.Movie {
	mov := &contentmovie.Movie{
		ID:               uuid.Must(uuid.NewV7()),
		Title:            r.Title,
		OriginalTitle:    ptrString(r.OriginalTitle),
		OriginalLanguage: ptrString(r.OriginalLanguage),
		PosterPath:       r.PosterPath,
		BackdropPath:     r.BackdropPath,
		ReleaseDate:      r.ReleaseDate,
		Year:             ptrInt32FromInt(r.Year),
	}

	if r.VoteAverage > 0 {
		va, _ := decimal.NewFromFloat64(r.VoteAverage)
		mov.VoteAverage = &va
	}
	if r.VoteCount > 0 {
		vc := int32(r.VoteCount)
		mov.VoteCount = &vc
	}
	if r.Popularity > 0 {
		pop, _ := decimal.NewFromFloat64(r.Popularity)
		mov.Popularity = &pop
	}

	// Set TMDb ID from provider ID
	if r.ProviderID != "" {
		var tmdbID int32
		_, _ = fmt.Sscanf(r.ProviderID, "%d", &tmdbID)
		mov.TMDbID = &tmdbID
	}

	return mov
}

// mapMetadataToMovie maps shared metadata to movie domain type.
func mapMetadataToMovie(mov *contentmovie.Movie, meta *metadata.MovieMetadata, releaseDates []metadata.ReleaseDate) {
	mov.Title = meta.Title
	mov.OriginalTitle = ptrString(meta.OriginalTitle)
	mov.OriginalLanguage = ptrString(meta.OriginalLanguage)
	mov.Overview = meta.Overview
	mov.Tagline = meta.Tagline
	mov.Status = ptrString(meta.Status)
	mov.ReleaseDate = meta.ReleaseDate
	mov.Runtime = meta.Runtime
	mov.Budget = meta.Budget
	mov.Revenue = meta.Revenue
	mov.PosterPath = meta.PosterPath
	mov.BackdropPath = meta.BackdropPath
	mov.IMDbID = meta.IMDbID
	mov.TMDbID = meta.TMDbID

	// Extract year from release date
	if meta.ReleaseDate != nil {
		year := int32(meta.ReleaseDate.Year())
		mov.Year = &year
	}

	// Map ratings
	if meta.VoteAverage > 0 {
		va, _ := decimal.NewFromFloat64(meta.VoteAverage)
		mov.VoteAverage = &va
	}
	if meta.VoteCount > 0 {
		vc := int32(meta.VoteCount)
		mov.VoteCount = &vc
	}
	if meta.Popularity > 0 {
		pop, _ := decimal.NewFromFloat64(meta.Popularity)
		mov.Popularity = &pop
	}

	// Map multi-language data
	if len(meta.Translations) > 0 {
		mov.TitlesI18n = make(map[string]string)
		mov.TaglinesI18n = make(map[string]string)
		mov.OverviewsI18n = make(map[string]string)

		for lang, trans := range meta.Translations {
			if trans.Title != "" {
				mov.TitlesI18n[lang] = trans.Title
			}
			if trans.Tagline != "" {
				mov.TaglinesI18n[lang] = trans.Tagline
			}
			if trans.Overview != "" {
				mov.OverviewsI18n[lang] = trans.Overview
			}
		}
	}

	// Map age ratings from release dates
	if len(releaseDates) > 0 {
		mov.AgeRatings = make(map[string]map[string]string)
		for _, rd := range releaseDates {
			if rd.Certification != "" {
				country := rd.CountryCode
				system := getAgeRatingSystem(country)
				if mov.AgeRatings[country] == nil {
					mov.AgeRatings[country] = make(map[string]string)
				}
				mov.AgeRatings[country][system] = rd.Certification
			}
		}
	}

	// Map external ratings (IMDb, Rotten Tomatoes, Metacritic, etc.)
	if len(meta.ExternalRatings) > 0 {
		mov.ExternalRatings = make([]contentmovie.ExternalRating, len(meta.ExternalRatings))
		for i, er := range meta.ExternalRatings {
			mov.ExternalRatings[i] = contentmovie.ExternalRating{
				Source: er.Source,
				Value:  er.Value,
				Score:  er.Score,
			}
		}
	}
}

// getAgeRatingSystem returns the rating system for a country code.
func getAgeRatingSystem(country string) string {
	switch country {
	case "US":
		return "MPAA"
	case "DE":
		return "FSK"
	case "GB":
		return "BBFC"
	case "FR":
		return "CNC"
	case "JP":
		return "Eirin"
	case "KR":
		return "KMRB"
	case "BR":
		return "DJCTQ"
	case "AU":
		return "ACB"
	default:
		return country // Use country code as fallback
	}
}

// mapCreditsToMovieCredits converts shared credits to movie credits.
func mapCreditsToMovieCredits(movieID uuid.UUID, credits *metadata.Credits) []contentmovie.MovieCredit {
	var result []contentmovie.MovieCredit

	// Map cast
	for _, c := range credits.Cast {
		var personID int32
		_, _ = fmt.Sscanf(c.ProviderID, "%d", &personID)

		credit := contentmovie.MovieCredit{
			ID:           uuid.Must(uuid.NewV7()),
			MovieID:      movieID,
			TMDbPersonID: personID,
			Name:         c.Name,
			Character:    ptrString(c.Character),
			CreditType:   "cast",
			CastOrder:    ptrInt32(&c.Order),
			ProfilePath:  c.ProfilePath,
		}
		result = append(result, credit)
	}

	// Map crew
	for _, c := range credits.Crew {
		var personID int32
		_, _ = fmt.Sscanf(c.ProviderID, "%d", &personID)

		credit := contentmovie.MovieCredit{
			ID:           uuid.Must(uuid.NewV7()),
			MovieID:      movieID,
			TMDbPersonID: personID,
			Name:         c.Name,
			Job:          ptrString(c.Job),
			Department:   ptrString(c.Department),
			CreditType:   "crew",
			ProfilePath:  c.ProfilePath,
		}
		result = append(result, credit)
	}

	return result
}

// Helper functions

func ptrString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func ptrInt32(i *int) *int32 {
	if i == nil {
		return nil
	}
	v := int32(*i)
	return &v
}

func ptrInt32FromInt(i *int) *int32 {
	if i == nil {
		return nil
	}
	v := int32(*i)
	return &v
}
