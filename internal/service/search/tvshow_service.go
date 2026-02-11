package search

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/content/tvshow"
	"github.com/lusoris/revenge/internal/infra/search"
	"github.com/lusoris/revenge/internal/util"
	"github.com/lusoris/revenge/internal/util/ptr"
	"github.com/typesense/typesense-go/v2/typesense/api"
)

// TVShowSearchService provides TV show search operations using Typesense.
type TVShowSearchService struct {
	client *search.Client
	logger *slog.Logger
}

// NewTVShowSearchService creates a new TV show search service.
func NewTVShowSearchService(client *search.Client, logger *slog.Logger) *TVShowSearchService {
	return &TVShowSearchService{
		client: client,
		logger: logger.With("service", "tvshow_search"),
	}
}

// IsEnabled returns true if search is enabled.
func (s *TVShowSearchService) IsEnabled() bool {
	return s.client != nil && s.client.IsEnabled()
}

// InitializeCollection creates the TV shows collection if it doesn't exist.
func (s *TVShowSearchService) InitializeCollection(ctx context.Context) error {
	if !s.IsEnabled() {
		s.logger.Debug("search disabled, skipping collection initialization")
		return nil
	}

	_, err := s.client.GetCollection(ctx, TVShowCollectionName)
	if err == nil {
		s.logger.Debug("tvshows collection already exists")
		return nil
	}

	schema := TVShowCollectionSchema()
	if err := s.client.CreateCollection(ctx, schema); err != nil {
		return fmt.Errorf("failed to create tvshows collection: %w", err)
	}

	s.logger.Info("created tvshows collection")
	return nil
}

// TVShowWithRelations bundles a series with its related data for indexing.
type TVShowWithRelations struct {
	Series   *tvshow.Series
	Genres   []tvshow.SeriesGenre
	Credits  []tvshow.SeriesCredit
	Networks []tvshow.Network
	HasFile  bool
}

// IndexSeries indexes a single series in Typesense.
func (s *TVShowSearchService) IndexSeries(ctx context.Context, series *tvshow.Series, genres []tvshow.SeriesGenre, credits []tvshow.SeriesCredit, networks []tvshow.Network, hasFile bool) error {
	if !s.IsEnabled() {
		return nil
	}

	doc := s.seriesToDocument(series, genres, credits, networks, hasFile)
	_, err := s.client.IndexDocument(ctx, TVShowCollectionName, doc)
	if err != nil {
		return fmt.Errorf("failed to index series %s: %w", series.ID, err)
	}

	s.logger.Debug("indexed series", "id", series.ID, "title", series.Title)
	return nil
}

// UpdateSeries updates a series in the search index.
func (s *TVShowSearchService) UpdateSeries(ctx context.Context, series *tvshow.Series, genres []tvshow.SeriesGenre, credits []tvshow.SeriesCredit, networks []tvshow.Network, hasFile bool) error {
	if !s.IsEnabled() {
		return nil
	}

	doc := s.seriesToDocument(series, genres, credits, networks, hasFile)
	_, err := s.client.UpdateDocument(ctx, TVShowCollectionName, series.ID.String(), doc)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return s.IndexSeries(ctx, series, genres, credits, networks, hasFile)
		}
		return fmt.Errorf("failed to update series %s: %w", series.ID, err)
	}

	s.logger.Debug("updated series", "id", series.ID, "title", series.Title)
	return nil
}

// RemoveSeries removes a series from the search index.
func (s *TVShowSearchService) RemoveSeries(ctx context.Context, seriesID uuid.UUID) error {
	if !s.IsEnabled() {
		return nil
	}

	_, err := s.client.DeleteDocument(ctx, TVShowCollectionName, seriesID.String())
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return fmt.Errorf("failed to remove series %s: %w", seriesID, err)
	}

	s.logger.Debug("removed series", "id", seriesID)
	return nil
}

// BulkIndexSeries indexes multiple series at once.
func (s *TVShowSearchService) BulkIndexSeries(ctx context.Context, shows []TVShowWithRelations) error {
	if !s.IsEnabled() {
		return nil
	}

	if len(shows) == 0 {
		return nil
	}

	docs := make([]interface{}, len(shows))
	for i, show := range shows {
		docs[i] = s.seriesToDocument(show.Series, show.Genres, show.Credits, show.Networks, show.HasFile)
	}

	results, err := s.client.ImportDocuments(ctx, TVShowCollectionName, docs, "upsert")
	if err != nil {
		return fmt.Errorf("failed to bulk index series: %w", err)
	}

	var successCount, errorCount int
	for _, result := range results {
		if result.Success {
			successCount++
		} else {
			errorCount++
			if result.Error != "" {
				s.logger.Warn("failed to index document", "error", result.Error)
			}
		}
	}

	s.logger.Info("bulk indexed series", "success", successCount, "errors", errorCount)
	return nil
}

// TVShowSearchResult represents a search result with hits and facets.
type TVShowSearchResult struct {
	Hits        []TVShowHit
	TotalHits   int
	TotalPages  int
	CurrentPage int
	Facets      map[string][]FacetValue
	SearchTime  time.Duration
}

// TVShowHit represents a single series hit in search results.
type TVShowHit struct {
	Document   TVShowDocument
	Score      float64
	Highlights map[string][]string
}

// TVShowSearchParams contains parameters for TV show search.
type TVShowSearchParams struct {
	Query             string
	Page              int
	PerPage           int
	SortBy            string
	FilterBy          string
	FacetBy           []string
	IncludeHighlights bool
}

// DefaultTVShowSearchParams returns default search parameters for TV shows.
func DefaultTVShowSearchParams() TVShowSearchParams {
	return TVShowSearchParams{
		Page:              1,
		PerPage:           20,
		SortBy:            "popularity:desc",
		FacetBy:           []string{"genres", "year", "status", "type", "networks", "has_file"},
		IncludeHighlights: true,
	}
}

// SearchSeries searches for TV shows with the given parameters.
func (s *TVShowSearchService) SearchSeries(ctx context.Context, params TVShowSearchParams) (*TVShowSearchResult, error) {
	if !s.IsEnabled() {
		return &TVShowSearchResult{}, nil
	}

	if params.Page < 1 {
		params.Page = 1
	}
	if params.PerPage < 1 {
		params.PerPage = 20
	}
	if params.PerPage > 100 {
		params.PerPage = 100
	}

	queryBy := "title,original_title,overview,cast,networks"
	page := params.Page
	perPage := params.PerPage

	searchParams := &api.SearchCollectionParams{
		Q:       &params.Query,
		QueryBy: &queryBy,
		Page:    &page,
		PerPage: &perPage,
	}

	if params.SortBy != "" {
		searchParams.SortBy = &params.SortBy
	}

	if params.FilterBy != "" {
		searchParams.FilterBy = &params.FilterBy
	}

	if len(params.FacetBy) > 0 {
		facetBy := strings.Join(params.FacetBy, ",")
		searchParams.FacetBy = &facetBy
	}

	start := time.Now()
	result, err := s.client.Search(ctx, TVShowCollectionName, searchParams)
	searchTime := time.Since(start)

	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	searchResult := &TVShowSearchResult{
		TotalHits:   deref(result.Found),
		CurrentPage: params.Page,
		TotalPages:  (deref(result.Found) + params.PerPage - 1) / params.PerPage,
		SearchTime:  searchTime,
		Hits:        make([]TVShowHit, 0, len(*result.Hits)),
		Facets:      make(map[string][]FacetValue),
	}

	if result.Hits != nil {
		for _, hit := range *result.Hits {
			tvHit := TVShowHit{
				Highlights: make(map[string][]string),
			}

			if hit.Document != nil {
				tvHit.Document = parseTVShowDocument(*hit.Document)
			}

			if hit.TextMatch != nil {
				tvHit.Score = float64(*hit.TextMatch)
			}

			if params.IncludeHighlights && hit.Highlights != nil {
				for _, hl := range *hit.Highlights {
					if hl.Field != nil && hl.Snippets != nil {
						tvHit.Highlights[*hl.Field] = *hl.Snippets
					}
				}
			}

			searchResult.Hits = append(searchResult.Hits, tvHit)
		}
	}

	if result.FacetCounts != nil {
		for _, fc := range *result.FacetCounts {
			if fc.FieldName == nil || fc.Counts == nil {
				continue
			}
			values := make([]FacetValue, 0, len(*fc.Counts))
			for _, count := range *fc.Counts {
				if count.Value == nil || count.Count == nil {
					continue
				}
				values = append(values, FacetValue{
					Value: *count.Value,
					Count: *count.Count,
				})
			}
			searchResult.Facets[*fc.FieldName] = values
		}
	}

	return searchResult, nil
}

// AutocompleteSeries provides search suggestions for series titles.
func (s *TVShowSearchService) AutocompleteSeries(ctx context.Context, query string, limit int) ([]string, error) {
	if !s.IsEnabled() {
		return nil, nil
	}

	if limit < 1 {
		limit = 5
	}
	if limit > 20 {
		limit = 20
	}

	queryBy := "title,original_title"
	perPage := limit

	searchParams := &api.SearchCollectionParams{
		Q:                   &query,
		QueryBy:             &queryBy,
		PerPage:             &perPage,
		Prefix:              ptr.To("true"),
		DropTokensThreshold: ptr.To(0),
	}

	result, err := s.client.Search(ctx, TVShowCollectionName, searchParams)
	if err != nil {
		return nil, fmt.Errorf("autocomplete failed: %w", err)
	}

	suggestions := make([]string, 0, limit)
	if result.Hits != nil {
		seen := make(map[string]bool)
		for _, hit := range *result.Hits {
			if hit.Document != nil {
				doc := parseTVShowDocument(*hit.Document)
				if doc.Title != "" && !seen[doc.Title] {
					suggestions = append(suggestions, doc.Title)
					seen[doc.Title] = true
				}
			}
			if len(suggestions) >= limit {
				break
			}
		}
	}

	return suggestions, nil
}

// GetFacets returns available facet values for filtering.
func (s *TVShowSearchService) GetFacets(ctx context.Context, facetNames []string) (map[string][]FacetValue, error) {
	if !s.IsEnabled() {
		return nil, nil
	}

	params := TVShowSearchParams{
		Query:   "*",
		Page:    1,
		PerPage: 0,
		FacetBy: facetNames,
	}

	result, err := s.SearchSeries(ctx, params)
	if err != nil {
		return nil, err
	}

	return result.Facets, nil
}

// ReindexAll reindexes all series in the database.
func (s *TVShowSearchService) ReindexAll(ctx context.Context, service tvshow.Service) error {
	if !s.IsEnabled() {
		return nil
	}

	s.logger.Info("starting full tvshow reindex")

	_ = s.client.DeleteCollection(ctx, TVShowCollectionName)
	if err := s.InitializeCollection(ctx); err != nil {
		return fmt.Errorf("failed to reinitialize collection: %w", err)
	}

	batchSize := int32(100)
	var offset int32
	totalIndexed := 0

	for {
		seriesList, err := service.ListSeries(ctx, tvshow.SeriesListFilters{
			Limit:  batchSize,
			Offset: offset,
		})
		if err != nil {
			return fmt.Errorf("failed to list series: %w", err)
		}

		if len(seriesList) == 0 {
			break
		}

		showsWithRelations := make([]TVShowWithRelations, 0, len(seriesList))
		for _, series := range seriesList {
			genres, _ := service.GetSeriesGenres(ctx, series.ID)
			cast, _, _ := service.GetSeriesCast(ctx, series.ID, 1000, 0)
			crew, _, _ := service.GetSeriesCrew(ctx, series.ID, 1000, 0)
			networks, _ := service.GetSeriesNetworks(ctx, series.ID)

			credits := append(cast, crew...)

			// Check if series has any episode files
			hasFile := false
			episodes, err := service.ListEpisodesBySeries(ctx, series.ID)
			if err == nil {
				for _, ep := range episodes {
					files, err := service.ListEpisodeFiles(ctx, ep.ID)
					if err == nil && len(files) > 0 {
						hasFile = true
						break
					}
				}
			}

			sCopy := series
			showsWithRelations = append(showsWithRelations, TVShowWithRelations{
				Series:   &sCopy,
				Genres:   genres,
				Credits:  credits,
				Networks: networks,
				HasFile:  hasFile,
			})
		}

		if err := s.BulkIndexSeries(ctx, showsWithRelations); err != nil {
			s.logger.Error("failed to index batch", "offset", offset, "error", err)
		}

		totalIndexed += len(seriesList)
		offset += batchSize

		if len(seriesList) < int(batchSize) {
			break
		}
	}

	s.logger.Info("completed full tvshow reindex", "total", totalIndexed)
	return nil
}

// seriesToDocument converts a series with related data to a search document.
func (s *TVShowSearchService) seriesToDocument(series *tvshow.Series, genres []tvshow.SeriesGenre, credits []tvshow.SeriesCredit, networks []tvshow.Network, hasFile bool) TVShowDocument {
	now := time.Now().Unix()

	doc := TVShowDocument{
		ID:        series.ID.String(),
		Title:     series.Title,
		HasFile:   hasFile,
		CreatedAt: series.CreatedAt.Unix(),
		UpdatedAt: now,
	}

	if series.TMDbID != nil {
		doc.TMDbID = *series.TMDbID
	}
	if series.TVDbID != nil {
		doc.TVDbID = *series.TVDbID
	}
	if series.IMDbID != nil {
		doc.IMDbID = *series.IMDbID
	}
	if series.OriginalTitle != nil {
		doc.OriginalTitle = *series.OriginalTitle
	}
	if series.FirstAirDate != nil {
		doc.FirstAirDate = series.FirstAirDate.Unix()
		doc.Year = util.SafeIntToInt32(series.FirstAirDate.Year())
	}
	if series.Overview != nil {
		doc.Overview = *series.Overview
	}
	if series.Status != nil {
		doc.Status = *series.Status
	}
	if series.Type != nil {
		doc.Type = *series.Type
	}
	doc.OriginalLanguage = series.OriginalLanguage
	if series.PosterPath != nil {
		doc.PosterPath = *series.PosterPath
	}
	if series.BackdropPath != nil {
		doc.BackdropPath = *series.BackdropPath
	}
	if series.VoteAverage != nil {
		f, _ := series.VoteAverage.Float64()
		doc.VoteAverage = f
	}
	if series.VoteCount != nil {
		doc.VoteCount = *series.VoteCount
	}
	if series.Popularity != nil {
		f, _ := series.Popularity.Float64()
		doc.Popularity = f
	}

	doc.TotalSeasons = series.TotalSeasons
	doc.TotalEpisodes = series.TotalEpisodes

	// Genres
	if len(genres) > 0 {
		genreNames := make([]string, len(genres))
		genreSlugs := make([]string, len(genres))
		for i, g := range genres {
			genreNames[i] = g.Name
			genreSlugs[i] = g.Slug
		}
		doc.Genres = genreNames
		doc.GenreSlugs = genreSlugs
	}

	// Credits - extract cast names
	if len(credits) > 0 {
		var cast []string
		for _, c := range credits {
			if c.CreditType == "cast" {
				cast = append(cast, c.Name)
			}
		}
		if len(cast) > 20 {
			cast = cast[:20]
		}
		doc.Cast = cast
	}

	// Networks
	if len(networks) > 0 {
		networkNames := make([]string, len(networks))
		for i, n := range networks {
			networkNames[i] = n.Name
		}
		doc.Networks = networkNames
	}

	return doc
}

// parseTVShowDocument parses a map into a TVShowDocument.
func parseTVShowDocument(data map[string]interface{}) TVShowDocument {
	doc := TVShowDocument{}

	if v, ok := data["id"].(string); ok {
		doc.ID = v
	}
	if v, ok := data["tmdb_id"].(float64); ok {
		doc.TMDbID = int32(v)
	}
	if v, ok := data["tvdb_id"].(float64); ok {
		doc.TVDbID = int32(v)
	}
	if v, ok := data["imdb_id"].(string); ok {
		doc.IMDbID = v
	}
	if v, ok := data["title"].(string); ok {
		doc.Title = v
	}
	if v, ok := data["original_title"].(string); ok {
		doc.OriginalTitle = v
	}
	if v, ok := data["year"].(float64); ok {
		doc.Year = int32(v)
	}
	if v, ok := data["first_air_date"].(float64); ok {
		doc.FirstAirDate = int64(v)
	}
	if v, ok := data["overview"].(string); ok {
		doc.Overview = v
	}
	if v, ok := data["status"].(string); ok {
		doc.Status = v
	}
	if v, ok := data["type"].(string); ok {
		doc.Type = v
	}
	if v, ok := data["original_language"].(string); ok {
		doc.OriginalLanguage = v
	}
	if v, ok := data["poster_path"].(string); ok {
		doc.PosterPath = v
	}
	if v, ok := data["backdrop_path"].(string); ok {
		doc.BackdropPath = v
	}
	if v, ok := data["vote_average"].(float64); ok {
		doc.VoteAverage = v
	}
	if v, ok := data["vote_count"].(float64); ok {
		doc.VoteCount = int32(v)
	}
	if v, ok := data["popularity"].(float64); ok {
		doc.Popularity = v
	}
	if v, ok := data["has_file"].(bool); ok {
		doc.HasFile = v
	}
	if v, ok := data["total_seasons"].(float64); ok {
		doc.TotalSeasons = int32(v)
	}
	if v, ok := data["total_episodes"].(float64); ok {
		doc.TotalEpisodes = int32(v)
	}
	if v, ok := data["created_at"].(float64); ok {
		doc.CreatedAt = int64(v)
	}
	if v, ok := data["updated_at"].(float64); ok {
		doc.UpdatedAt = int64(v)
	}

	// Parse string arrays
	if v, ok := data["genres"].([]interface{}); ok {
		doc.Genres = toStringSlice(v)
	}
	if v, ok := data["cast"].([]interface{}); ok {
		doc.Cast = toStringSlice(v)
	}
	if v, ok := data["networks"].([]interface{}); ok {
		doc.Networks = toStringSlice(v)
	}
	if v, ok := data["genre_slugs"].([]interface{}); ok {
		doc.GenreSlugs = toStringSlice(v)
	}

	return doc
}
