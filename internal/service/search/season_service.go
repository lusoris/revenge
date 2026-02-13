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
	"github.com/typesense/typesense-go/v2/typesense/api"
)

// SeasonSearchService provides season search operations using Typesense.
type SeasonSearchService struct {
	client *search.Client
	logger *slog.Logger
}

// NewSeasonSearchService creates a new season search service.
func NewSeasonSearchService(client *search.Client, logger *slog.Logger) *SeasonSearchService {
	return &SeasonSearchService{
		client: client,
		logger: logger.With("service", "season_search"),
	}
}

// IsEnabled returns true if search is enabled.
func (s *SeasonSearchService) IsEnabled() bool {
	return s.client != nil && s.client.IsEnabled()
}

// InitializeCollection creates the seasons collection if it doesn't exist.
func (s *SeasonSearchService) InitializeCollection(ctx context.Context) error {
	if !s.IsEnabled() {
		s.logger.Debug("search disabled, skipping collection initialization")
		return nil
	}

	_, err := s.client.GetCollection(ctx, SeasonCollectionName)
	if err == nil {
		s.logger.Debug("seasons collection already exists")
		return nil
	}

	schema := SeasonCollectionSchema()
	if err := s.client.CreateCollection(ctx, schema); err != nil {
		return fmt.Errorf("failed to create seasons collection: %w", err)
	}

	s.logger.Info("created seasons collection")
	return nil
}

// SeasonWithContext bundles a season with its parent series info for indexing.
type SeasonWithContext struct {
	Season           *tvshow.Season
	SeriesTitle      string
	SeriesPosterPath string
}

// IndexSeason indexes a single season in Typesense.
func (s *SeasonSearchService) IndexSeason(ctx context.Context, season *tvshow.Season, seriesTitle, seriesPosterPath string) error {
	if !s.IsEnabled() {
		return nil
	}

	doc := s.seasonToDocument(season, seriesTitle, seriesPosterPath)
	_, err := s.client.IndexDocument(ctx, SeasonCollectionName, doc)
	if err != nil {
		return fmt.Errorf("failed to index season %s: %w", season.ID, err)
	}

	s.logger.Debug("indexed season", "id", season.ID, "name", season.Name)
	return nil
}

// UpdateSeason updates a season in the search index.
func (s *SeasonSearchService) UpdateSeason(ctx context.Context, season *tvshow.Season, seriesTitle, seriesPosterPath string) error {
	if !s.IsEnabled() {
		return nil
	}

	doc := s.seasonToDocument(season, seriesTitle, seriesPosterPath)
	_, err := s.client.UpdateDocument(ctx, SeasonCollectionName, season.ID.String(), doc)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return s.IndexSeason(ctx, season, seriesTitle, seriesPosterPath)
		}
		return fmt.Errorf("failed to update season %s: %w", season.ID, err)
	}

	s.logger.Debug("updated season", "id", season.ID, "name", season.Name)
	return nil
}

// RemoveSeason removes a season from the search index.
func (s *SeasonSearchService) RemoveSeason(ctx context.Context, seasonID uuid.UUID) error {
	if !s.IsEnabled() {
		return nil
	}

	_, err := s.client.DeleteDocument(ctx, SeasonCollectionName, seasonID.String())
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return fmt.Errorf("failed to remove season %s: %w", seasonID, err)
	}

	s.logger.Debug("removed season", "id", seasonID)
	return nil
}

// RemoveSeasonsBySeries removes all seasons for a series from the search index.
func (s *SeasonSearchService) RemoveSeasonsBySeries(ctx context.Context, seriesID uuid.UUID) error {
	if !s.IsEnabled() {
		return nil
	}

	filterBy := fmt.Sprintf("series_id:=%s", seriesID.String())
	searchParams := &api.SearchCollectionParams{
		Q:        new("*"),
		QueryBy:  new("name"),
		FilterBy: &filterBy,
		PerPage:  new(250),
	}

	result, err := s.client.Search(ctx, SeasonCollectionName, searchParams)
	if err != nil {
		return fmt.Errorf("failed to search seasons for series %s: %w", seriesID, err)
	}

	if result.Hits == nil {
		return nil
	}

	for _, hit := range *result.Hits {
		if hit.Document != nil {
			doc := parseSeasonDocument(*hit.Document)
			if doc.ID != "" {
				_, _ = s.client.DeleteDocument(ctx, SeasonCollectionName, doc.ID)
			}
		}
	}

	s.logger.Debug("removed seasons for series", "series_id", seriesID)
	return nil
}

// BulkIndexSeasons indexes multiple seasons at once.
func (s *SeasonSearchService) BulkIndexSeasons(ctx context.Context, seasons []SeasonWithContext) error {
	if !s.IsEnabled() {
		return nil
	}

	if len(seasons) == 0 {
		return nil
	}

	docs := make([]any, len(seasons))
	for i, season := range seasons {
		docs[i] = s.seasonToDocument(season.Season, season.SeriesTitle, season.SeriesPosterPath)
	}

	results, err := s.client.ImportDocuments(ctx, SeasonCollectionName, docs, "upsert")
	if err != nil {
		return fmt.Errorf("failed to bulk index seasons: %w", err)
	}

	var successCount, errorCount int
	for _, result := range results {
		if result.Success {
			successCount++
		} else {
			errorCount++
			if result.Error != "" {
				s.logger.Warn("failed to index season document", "error", result.Error)
			}
		}
	}

	s.logger.Info("bulk indexed seasons", "success", successCount, "errors", errorCount)
	return nil
}

// SeasonSearchResult represents a season search result with hits.
type SeasonSearchResult struct {
	Hits        []SeasonHit
	TotalHits   int
	TotalPages  int
	CurrentPage int
	Facets      map[string][]FacetValue
	SearchTime  time.Duration
}

// SeasonHit represents a single season hit in search results.
type SeasonHit struct {
	Document   SeasonDocument
	Score      float64
	Highlights map[string][]string
}

// SeasonSearchParams contains parameters for season search.
type SeasonSearchParams struct {
	Query             string
	Page              int
	PerPage           int
	SortBy            string
	FilterBy          string
	FacetBy           []string
	IncludeHighlights bool
}

// DefaultSeasonSearchParams returns default search parameters for seasons.
func DefaultSeasonSearchParams() SeasonSearchParams {
	return SeasonSearchParams{
		Page:              1,
		PerPage:           20,
		SortBy:            "air_date:desc",
		FacetBy:           []string{"season_number"},
		IncludeHighlights: true,
	}
}

// SearchSeasons searches for seasons with the given parameters.
func (s *SeasonSearchService) SearchSeasons(ctx context.Context, params SeasonSearchParams) (*SeasonSearchResult, error) {
	if !s.IsEnabled() {
		return &SeasonSearchResult{}, nil
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

	queryBy := "name,overview,series_title"
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
	result, err := s.client.Search(ctx, SeasonCollectionName, searchParams)
	searchTime := time.Since(start)

	if err != nil {
		return nil, fmt.Errorf("season search failed: %w", err)
	}

	searchResult := &SeasonSearchResult{
		TotalHits:   deref(result.Found),
		CurrentPage: params.Page,
		TotalPages:  (deref(result.Found) + params.PerPage - 1) / params.PerPage,
		SearchTime:  searchTime,
		Hits:        make([]SeasonHit, 0, len(derefHits(result.Hits))),
		Facets:      make(map[string][]FacetValue),
	}

	if result.Hits != nil {
		for _, hit := range *result.Hits {
			sHit := SeasonHit{
				Highlights: make(map[string][]string),
			}

			if hit.Document != nil {
				sHit.Document = parseSeasonDocument(*hit.Document)
			}

			if hit.TextMatch != nil {
				sHit.Score = float64(*hit.TextMatch)
			}

			if params.IncludeHighlights && hit.Highlights != nil {
				for _, hl := range *hit.Highlights {
					if hl.Field != nil && hl.Snippets != nil {
						sHit.Highlights[*hl.Field] = *hl.Snippets
					}
				}
			}

			searchResult.Hits = append(searchResult.Hits, sHit)
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

// AutocompleteSeasons provides search suggestions for season names.
func (s *SeasonSearchService) AutocompleteSeasons(ctx context.Context, query string, limit int) ([]string, error) {
	if !s.IsEnabled() {
		return nil, nil
	}

	if limit < 1 {
		limit = 5
	}
	if limit > 20 {
		limit = 20
	}

	queryBy := "name,series_title"
	perPage := limit

	searchParams := &api.SearchCollectionParams{
		Q:                   &query,
		QueryBy:             &queryBy,
		PerPage:             &perPage,
		Prefix:              new("true"),
		DropTokensThreshold: new(0),
	}

	result, err := s.client.Search(ctx, SeasonCollectionName, searchParams)
	if err != nil {
		return nil, fmt.Errorf("season autocomplete failed: %w", err)
	}

	suggestions := make([]string, 0, limit)
	if result.Hits != nil {
		seen := make(map[string]bool)
		for _, hit := range *result.Hits {
			if hit.Document != nil {
				doc := parseSeasonDocument(*hit.Document)
				label := fmt.Sprintf("%s - %s", doc.SeriesTitle, doc.Name)
				if !seen[label] {
					suggestions = append(suggestions, label)
					seen[label] = true
				}
			}
			if len(suggestions) >= limit {
				break
			}
		}
	}

	return suggestions, nil
}

// ReindexAll reindexes all seasons in the database.
func (s *SeasonSearchService) ReindexAll(ctx context.Context, service tvshow.Service) error {
	if !s.IsEnabled() {
		return nil
	}

	s.logger.Info("starting full season reindex")

	_ = s.client.DeleteCollection(ctx, SeasonCollectionName)
	if err := s.InitializeCollection(ctx); err != nil {
		return fmt.Errorf("failed to reinitialize seasons collection: %w", err)
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

		for _, series := range seriesList {
			seasons, err := service.ListSeasons(ctx, series.ID)
			if err != nil {
				s.logger.Warn("failed to list seasons for series",
					slog.String("series_id", series.ID.String()),
					slog.Any("error", err),
				)
				continue
			}

			if len(seasons) == 0 {
				continue
			}

			posterPath := ""
			if series.PosterPath != nil {
				posterPath = *series.PosterPath
			}

			batch := make([]SeasonWithContext, 0, len(seasons))
			for i := range seasons {
				batch = append(batch, SeasonWithContext{
					Season:           &seasons[i],
					SeriesTitle:      series.Title,
					SeriesPosterPath: posterPath,
				})
			}

			if err := s.BulkIndexSeasons(ctx, batch); err != nil {
				s.logger.Error("failed to index season batch",
					slog.String("series_id", series.ID.String()),
					slog.Any("error", err),
				)
			}

			totalIndexed += len(batch)
		}

		offset += batchSize

		if len(seriesList) < int(batchSize) {
			break
		}
	}

	s.logger.Info("completed full season reindex", "total", totalIndexed)
	return nil
}

// seasonToDocument converts a season with series context to a search document.
func (s *SeasonSearchService) seasonToDocument(season *tvshow.Season, seriesTitle, seriesPosterPath string) SeasonDocument {
	now := time.Now().Unix()

	doc := SeasonDocument{
		ID:               season.ID.String(),
		SeriesID:         season.SeriesID.String(),
		SeasonNumber:     season.SeasonNumber,
		Name:             season.Name,
		EpisodeCount:     season.EpisodeCount,
		SeriesTitle:      seriesTitle,
		SeriesPosterPath: seriesPosterPath,
		CreatedAt:        season.CreatedAt.Unix(),
		UpdatedAt:        now,
	}

	if season.TMDbID != nil {
		doc.TMDbID = *season.TMDbID
	}
	if season.Overview != nil {
		doc.Overview = *season.Overview
	}
	if season.AirDate != nil {
		doc.AirDate = season.AirDate.Unix()
	}
	if season.VoteAverage != nil {
		f, _ := season.VoteAverage.Float64()
		doc.VoteAverage = f
	}
	if season.PosterPath != nil {
		doc.PosterPath = *season.PosterPath
	}

	return doc
}

// parseSeasonDocument parses a map into a SeasonDocument.
func parseSeasonDocument(data map[string]any) SeasonDocument {
	doc := SeasonDocument{}

	if v, ok := data["id"].(string); ok {
		doc.ID = v
	}
	if v, ok := data["series_id"].(string); ok {
		doc.SeriesID = v
	}
	if v, ok := data["tmdb_id"].(float64); ok {
		doc.TMDbID = int32(v)
	}
	if v, ok := data["season_number"].(float64); ok {
		doc.SeasonNumber = int32(v)
	}
	if v, ok := data["name"].(string); ok {
		doc.Name = v
	}
	if v, ok := data["overview"].(string); ok {
		doc.Overview = v
	}
	if v, ok := data["air_date"].(float64); ok {
		doc.AirDate = int64(v)
	}
	if v, ok := data["episode_count"].(float64); ok {
		doc.EpisodeCount = int32(v)
	}
	if v, ok := data["vote_average"].(float64); ok {
		doc.VoteAverage = v
	}
	if v, ok := data["poster_path"].(string); ok {
		doc.PosterPath = v
	}
	if v, ok := data["series_title"].(string); ok {
		doc.SeriesTitle = v
	}
	if v, ok := data["series_poster_path"].(string); ok {
		doc.SeriesPosterPath = v
	}
	if v, ok := data["created_at"].(float64); ok {
		doc.CreatedAt = int64(v)
	}
	if v, ok := data["updated_at"].(float64); ok {
		doc.UpdatedAt = int64(v)
	}

	return doc
}
