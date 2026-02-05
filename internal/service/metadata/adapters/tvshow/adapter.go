// Package tvshow provides an adapter that bridges the shared metadata service
// to the TV show content module.
package tvshow

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	contenttvshow "github.com/lusoris/revenge/internal/content/tvshow"
	"github.com/lusoris/revenge/internal/service/metadata"
)

// Adapter wraps the shared metadata service for TV show-specific operations.
// This adapter implements the tvshow.MetadataProvider interface using the shared service.
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
var _ contenttvshow.MetadataProvider = (*Adapter)(nil)

// SearchSeries searches for TV series using the shared metadata service.
func (a *Adapter) SearchSeries(ctx context.Context, query string, year *int) ([]*contenttvshow.Series, error) {
	opts := metadata.SearchOptions{
		Year:     year,
		Language: a.languages[0],
	}

	results, err := a.service.SearchTVShow(ctx, query, opts)
	if err != nil {
		return nil, fmt.Errorf("search series: %w", err)
	}

	series := make([]*contenttvshow.Series, len(results))
	for i, r := range results {
		series[i] = mapSearchResultToSeries(&r)
	}

	return series, nil
}

// EnrichSeries enriches a series with metadata from the shared service.
func (a *Adapter) EnrichSeries(ctx context.Context, series *contenttvshow.Series) error {
	if series.TMDbID == nil {
		return fmt.Errorf("series has no TMDb ID")
	}

	meta, err := a.service.GetTVShowMetadata(ctx, *series.TMDbID, a.languages)
	if err != nil {
		return fmt.Errorf("get series metadata: %w", err)
	}

	// Get content ratings for age ratings
	contentRatings, err := a.service.GetTVShowContentRatings(ctx, *series.TMDbID)
	if err != nil {
		// Continue without content ratings
		contentRatings = nil
	}

	// Map to series domain type
	mapMetadataToSeries(series, meta, contentRatings)

	return nil
}

// EnrichSeason enriches a season with metadata from the shared service.
func (a *Adapter) EnrichSeason(ctx context.Context, season *contenttvshow.Season, seriesTMDbID int32) error {
	meta, err := a.service.GetSeasonMetadata(ctx, seriesTMDbID, int(season.SeasonNumber), a.languages)
	if err != nil {
		return fmt.Errorf("get season metadata: %w", err)
	}

	mapSeasonMetadataToSeason(season, meta)
	return nil
}

// EnrichEpisode enriches an episode with metadata from the shared service.
func (a *Adapter) EnrichEpisode(ctx context.Context, episode *contenttvshow.Episode, seriesTMDbID int32) error {
	meta, err := a.service.GetEpisodeMetadata(ctx, seriesTMDbID, int(episode.SeasonNumber), int(episode.EpisodeNumber), a.languages)
	if err != nil {
		return fmt.Errorf("get episode metadata: %w", err)
	}

	mapEpisodeMetadataToEpisode(episode, meta)
	return nil
}

// GetSeriesCredits retrieves series credits using the shared service.
func (a *Adapter) GetSeriesCredits(ctx context.Context, seriesID uuid.UUID, tmdbID int) ([]contenttvshow.SeriesCredit, error) {
	credits, err := a.service.GetTVShowCredits(ctx, int32(tmdbID))
	if err != nil {
		return nil, fmt.Errorf("get series credits: %w", err)
	}

	return mapCreditsToSeriesCredits(seriesID, credits), nil
}

// GetSeriesGenres retrieves series genres using the shared service.
func (a *Adapter) GetSeriesGenres(ctx context.Context, seriesID uuid.UUID, tmdbID int) ([]contenttvshow.SeriesGenre, error) {
	meta, err := a.service.GetTVShowMetadata(ctx, int32(tmdbID), []string{a.languages[0]})
	if err != nil {
		return nil, fmt.Errorf("get series metadata for genres: %w", err)
	}

	genres := make([]contenttvshow.SeriesGenre, len(meta.Genres))
	for i, g := range meta.Genres {
		genres[i] = contenttvshow.SeriesGenre{
			ID:          uuid.New(),
			SeriesID:    seriesID,
			TMDbGenreID: int32(g.ID),
			Name:        g.Name,
		}
	}

	return genres, nil
}

// GetSeriesNetworks retrieves series networks using the shared service.
func (a *Adapter) GetSeriesNetworks(ctx context.Context, tmdbID int) ([]contenttvshow.Network, error) {
	meta, err := a.service.GetTVShowMetadata(ctx, int32(tmdbID), []string{a.languages[0]})
	if err != nil {
		return nil, fmt.Errorf("get series metadata for networks: %w", err)
	}

	networks := make([]contenttvshow.Network, len(meta.Networks))
	for i, n := range meta.Networks {
		networks[i] = contenttvshow.Network{
			ID:            uuid.New(),
			TMDbID:        int32(n.ID),
			Name:          n.Name,
			LogoPath:      n.LogoPath,
			OriginCountry: ptrString(n.OriginCountry),
		}
	}

	return networks, nil
}

// ClearCache is a no-op for the adapter (cache is managed by the service).
func (a *Adapter) ClearCache() {
	// Cache is managed by the shared metadata service
}

// mapSearchResultToSeries converts a search result to a series domain type.
func mapSearchResultToSeries(r *metadata.TVShowSearchResult) *contenttvshow.Series {
	series := &contenttvshow.Series{
		ID:               uuid.New(),
		Title:            r.Name,
		OriginalTitle:    ptrString(r.OriginalName),
		OriginalLanguage: r.OriginalLanguage,
		PosterPath:       r.PosterPath,
		BackdropPath:     r.BackdropPath,
		FirstAirDate:     r.FirstAirDate,
	}

	if r.VoteAverage > 0 {
		va := decimal.NewFromFloat(r.VoteAverage)
		series.VoteAverage = &va
	}
	if r.VoteCount > 0 {
		vc := int32(r.VoteCount)
		series.VoteCount = &vc
	}
	if r.Popularity > 0 {
		pop := decimal.NewFromFloat(r.Popularity)
		series.Popularity = &pop
	}

	// Set TMDb ID from provider ID
	if r.ProviderID != "" {
		var tmdbID int32
		fmt.Sscanf(r.ProviderID, "%d", &tmdbID)
		series.TMDbID = &tmdbID
	}

	return series
}

// mapMetadataToSeries maps shared metadata to series domain type.
func mapMetadataToSeries(series *contenttvshow.Series, meta *metadata.TVShowMetadata, contentRatings []metadata.ContentRating) {
	series.Title = meta.Name
	series.OriginalTitle = ptrString(meta.OriginalName)
	series.OriginalLanguage = meta.OriginalLanguage
	series.Overview = meta.Overview
	series.Tagline = meta.Tagline
	series.Status = ptrString(meta.Status)
	series.Type = ptrString(meta.Type)
	series.FirstAirDate = meta.FirstAirDate
	series.LastAirDate = meta.LastAirDate
	series.TotalSeasons = int32(meta.NumberOfSeasons)
	series.TotalEpisodes = int32(meta.NumberOfEpisodes)
	series.PosterPath = meta.PosterPath
	series.BackdropPath = meta.BackdropPath
	series.Homepage = meta.Homepage
	series.TrailerURL = meta.TrailerURL
	series.IMDbID = meta.IMDbID
	series.TMDbID = meta.TMDbID
	series.TVDbID = meta.TVDbID

	// Map ratings
	if meta.VoteAverage > 0 {
		va := decimal.NewFromFloat(meta.VoteAverage)
		series.VoteAverage = &va
	}
	if meta.VoteCount > 0 {
		vc := int32(meta.VoteCount)
		series.VoteCount = &vc
	}
	if meta.Popularity > 0 {
		pop := decimal.NewFromFloat(meta.Popularity)
		series.Popularity = &pop
	}

	// Map multi-language data
	if len(meta.Translations) > 0 {
		series.TitlesI18n = make(map[string]string)
		series.TaglinesI18n = make(map[string]string)
		series.OverviewsI18n = make(map[string]string)

		for lang, trans := range meta.Translations {
			if trans.Name != "" {
				series.TitlesI18n[lang] = trans.Name
			}
			if trans.Tagline != "" {
				series.TaglinesI18n[lang] = trans.Tagline
			}
			if trans.Overview != "" {
				series.OverviewsI18n[lang] = trans.Overview
			}
		}
	}

	// Map age ratings from content ratings
	if len(contentRatings) > 0 {
		series.AgeRatings = make(map[string]map[string]string)
		for _, cr := range contentRatings {
			if cr.Rating != "" {
				country := cr.CountryCode
				system := getTVAgeRatingSystem(country)
				if series.AgeRatings[country] == nil {
					series.AgeRatings[country] = make(map[string]string)
				}
				series.AgeRatings[country][system] = cr.Rating
			}
		}
	}
}

// mapSeasonMetadataToSeason maps shared metadata to season domain type.
func mapSeasonMetadataToSeason(season *contenttvshow.Season, meta *metadata.SeasonMetadata) {
	season.Name = meta.Name
	season.Overview = meta.Overview
	season.PosterPath = meta.PosterPath
	season.AirDate = meta.AirDate
	season.EpisodeCount = int32(len(meta.Episodes))
	season.TMDbID = meta.TMDbID

	if meta.VoteAverage > 0 {
		va := decimal.NewFromFloat(meta.VoteAverage)
		season.VoteAverage = &va
	}

	// Map multi-language data
	if len(meta.Translations) > 0 {
		season.NamesI18n = make(map[string]string)
		season.OverviewsI18n = make(map[string]string)

		for lang, trans := range meta.Translations {
			if trans.Name != "" {
				season.NamesI18n[lang] = trans.Name
			}
			if trans.Overview != "" {
				season.OverviewsI18n[lang] = trans.Overview
			}
		}
	}
}

// mapEpisodeMetadataToEpisode maps shared metadata to episode domain type.
func mapEpisodeMetadataToEpisode(episode *contenttvshow.Episode, meta *metadata.EpisodeMetadata) {
	episode.Title = meta.Name
	episode.Overview = meta.Overview
	episode.AirDate = meta.AirDate
	episode.Runtime = meta.Runtime
	episode.StillPath = meta.StillPath
	episode.ProductionCode = meta.ProductionCode
	episode.TMDbID = meta.TMDbID
	episode.IMDbID = meta.IMDbID

	if meta.VoteAverage > 0 {
		va := decimal.NewFromFloat(meta.VoteAverage)
		episode.VoteAverage = &va
	}
	if meta.VoteCount > 0 {
		vc := int32(meta.VoteCount)
		episode.VoteCount = &vc
	}

	// Map multi-language data
	if len(meta.Translations) > 0 {
		episode.TitlesI18n = make(map[string]string)
		episode.OverviewsI18n = make(map[string]string)

		for lang, trans := range meta.Translations {
			if trans.Name != "" {
				episode.TitlesI18n[lang] = trans.Name
			}
			if trans.Overview != "" {
				episode.OverviewsI18n[lang] = trans.Overview
			}
		}
	}
}

// getTVAgeRatingSystem returns the rating system for a country code.
func getTVAgeRatingSystem(country string) string {
	switch country {
	case "US":
		return "TV Parental Guidelines"
	case "DE":
		return "FSK"
	case "GB":
		return "BBFC"
	case "FR":
		return "CSA"
	case "JP":
		return "EIRIN"
	case "KR":
		return "KMRB"
	case "BR":
		return "DJCTQ"
	case "AU":
		return "ACB"
	case "CA":
		return "CHVRS"
	default:
		return country // Use country code as fallback
	}
}

// mapCreditsToSeriesCredits converts shared credits to series credits.
func mapCreditsToSeriesCredits(seriesID uuid.UUID, credits *metadata.Credits) []contenttvshow.SeriesCredit {
	var result []contenttvshow.SeriesCredit

	// Map cast
	for _, c := range credits.Cast {
		var personID int32
		fmt.Sscanf(c.ProviderID, "%d", &personID)

		credit := contenttvshow.SeriesCredit{
			ID:           uuid.New(),
			SeriesID:     seriesID,
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
		fmt.Sscanf(c.ProviderID, "%d", &personID)

		credit := contenttvshow.SeriesCredit{
			ID:           uuid.New(),
			SeriesID:     seriesID,
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
