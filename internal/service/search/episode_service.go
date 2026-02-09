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
	"github.com/lusoris/revenge/internal/util/ptr"
	"github.com/typesense/typesense-go/v2/typesense/api"
)

// EpisodeSearchService provides episode search operations using Typesense.
type EpisodeSearchService struct {
	client *search.Client
	logger *slog.Logger
}

// NewEpisodeSearchService creates a new episode search service.
func NewEpisodeSearchService(client *search.Client, logger *slog.Logger) *EpisodeSearchService {
	return &EpisodeSearchService{
		client: client,
		logger: logger.With("service", "episode_search"),
	}
}

// IsEnabled returns true if search is enabled.
func (s *EpisodeSearchService) IsEnabled() bool {
	return s.client != nil && s.client.IsEnabled()
}

// InitializeCollection creates the episodes collection if it doesn't exist.
func (s *EpisodeSearchService) InitializeCollection(ctx context.Context) error {
	if !s.IsEnabled() {
		s.logger.Debug("search disabled, skipping collection initialization")
		return nil
	}

	_, err := s.client.GetCollection(ctx, EpisodeCollectionName)
	if err == nil {
		s.logger.Debug("episodes collection already exists")
		return nil
	}

	schema := EpisodeCollectionSchema()
	if err := s.client.CreateCollection(ctx, schema); err != nil {
		return fmt.Errorf("failed to create episodes collection: %w", err)
	}

	s.logger.Info("created episodes collection")
	return nil
}

// EpisodeWithContext bundles an episode with its parent series info for indexing.
type EpisodeWithContext struct {
	Episode         *tvshow.Episode
	SeriesTitle     string
	SeriesPosterPath string
	HasFile         bool
}

// IndexEpisode indexes a single episode in Typesense.
func (s *EpisodeSearchService) IndexEpisode(ctx context.Context, ep *tvshow.Episode, seriesTitle, seriesPosterPath string, hasFile bool) error {
	if !s.IsEnabled() {
		return nil
	}

	doc := s.episodeToDocument(ep, seriesTitle, seriesPosterPath, hasFile)
	_, err := s.client.IndexDocument(ctx, EpisodeCollectionName, doc)
	if err != nil {
		return fmt.Errorf("failed to index episode %s: %w", ep.ID, err)
	}

	s.logger.Debug("indexed episode", "id", ep.ID, "title", ep.Title)
	return nil
}

// UpdateEpisode updates an episode in the search index.
func (s *EpisodeSearchService) UpdateEpisode(ctx context.Context, ep *tvshow.Episode, seriesTitle, seriesPosterPath string, hasFile bool) error {
	if !s.IsEnabled() {
		return nil
	}

	doc := s.episodeToDocument(ep, seriesTitle, seriesPosterPath, hasFile)
	_, err := s.client.UpdateDocument(ctx, EpisodeCollectionName, ep.ID.String(), doc)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return s.IndexEpisode(ctx, ep, seriesTitle, seriesPosterPath, hasFile)
		}
		return fmt.Errorf("failed to update episode %s: %w", ep.ID, err)
	}

	s.logger.Debug("updated episode", "id", ep.ID, "title", ep.Title)
	return nil
}

// RemoveEpisode removes an episode from the search index.
func (s *EpisodeSearchService) RemoveEpisode(ctx context.Context, episodeID uuid.UUID) error {
	if !s.IsEnabled() {
		return nil
	}

	_, err := s.client.DeleteDocument(ctx, EpisodeCollectionName, episodeID.String())
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return fmt.Errorf("failed to remove episode %s: %w", episodeID, err)
	}

	s.logger.Debug("removed episode", "id", episodeID)
	return nil
}

// RemoveEpisodesBySeries removes all episodes for a series from the search index.
func (s *EpisodeSearchService) RemoveEpisodesBySeries(ctx context.Context, seriesID uuid.UUID) error {
	if !s.IsEnabled() {
		return nil
	}

	filterBy := fmt.Sprintf("series_id:=%s", seriesID.String())
	searchParams := &api.SearchCollectionParams{
		Q:        ptr.To("*"),
		QueryBy:  ptr.To("title"),
		FilterBy: &filterBy,
		PerPage:  ptr.To(250),
	}

	result, err := s.client.Search(ctx, EpisodeCollectionName, searchParams)
	if err != nil {
		return fmt.Errorf("failed to search episodes for series %s: %w", seriesID, err)
	}

	if result.Hits == nil {
		return nil
	}

	for _, hit := range *result.Hits {
		if hit.Document != nil {
			doc := parseEpisodeDocument(*hit.Document)
			if doc.ID != "" {
				_, _ = s.client.DeleteDocument(ctx, EpisodeCollectionName, doc.ID)
			}
		}
	}

	s.logger.Debug("removed episodes for series", "series_id", seriesID)
	return nil
}

// BulkIndexEpisodes indexes multiple episodes at once.
func (s *EpisodeSearchService) BulkIndexEpisodes(ctx context.Context, episodes []EpisodeWithContext) error {
	if !s.IsEnabled() {
		return nil
	}

	if len(episodes) == 0 {
		return nil
	}

	docs := make([]interface{}, len(episodes))
	for i, ep := range episodes {
		docs[i] = s.episodeToDocument(ep.Episode, ep.SeriesTitle, ep.SeriesPosterPath, ep.HasFile)
	}

	results, err := s.client.ImportDocuments(ctx, EpisodeCollectionName, docs, "upsert")
	if err != nil {
		return fmt.Errorf("failed to bulk index episodes: %w", err)
	}

	var successCount, errorCount int
	for _, result := range results {
		if result.Success {
			successCount++
		} else {
			errorCount++
			if result.Error != "" {
				s.logger.Warn("failed to index episode document", "error", result.Error)
			}
		}
	}

	s.logger.Info("bulk indexed episodes", "success", successCount, "errors", errorCount)
	return nil
}

// EpisodeSearchResult represents an episode search result with hits.
type EpisodeSearchResult struct {
	Hits        []EpisodeHit
	TotalHits   int
	TotalPages  int
	CurrentPage int
	Facets      map[string][]FacetValue
	SearchTime  time.Duration
}

// EpisodeHit represents a single episode hit in search results.
type EpisodeHit struct {
	Document   EpisodeDocument
	Score      float64
	Highlights map[string][]string
}

// EpisodeSearchParams contains parameters for episode search.
type EpisodeSearchParams struct {
	Query             string
	Page              int
	PerPage           int
	SortBy            string
	FilterBy          string
	FacetBy           []string
	IncludeHighlights bool
}

// DefaultEpisodeSearchParams returns default search parameters for episodes.
func DefaultEpisodeSearchParams() EpisodeSearchParams {
	return EpisodeSearchParams{
		Page:              1,
		PerPage:           20,
		SortBy:            "air_date:desc",
		FacetBy:           []string{"season_number", "has_file"},
		IncludeHighlights: true,
	}
}

// SearchEpisodes searches for episodes with the given parameters.
func (s *EpisodeSearchService) SearchEpisodes(ctx context.Context, params EpisodeSearchParams) (*EpisodeSearchResult, error) {
	if !s.IsEnabled() {
		return &EpisodeSearchResult{}, nil
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

	queryBy := "title,overview,series_title"
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
	result, err := s.client.Search(ctx, EpisodeCollectionName, searchParams)
	searchTime := time.Since(start)

	if err != nil {
		return nil, fmt.Errorf("episode search failed: %w", err)
	}

	searchResult := &EpisodeSearchResult{
		TotalHits:   deref(result.Found),
		CurrentPage: params.Page,
		TotalPages:  (deref(result.Found) + params.PerPage - 1) / params.PerPage,
		SearchTime:  searchTime,
		Hits:        make([]EpisodeHit, 0, len(derefHits(result.Hits))),
		Facets:      make(map[string][]FacetValue),
	}

	if result.Hits != nil {
		for _, hit := range *result.Hits {
			epHit := EpisodeHit{
				Highlights: make(map[string][]string),
			}

			if hit.Document != nil {
				epHit.Document = parseEpisodeDocument(*hit.Document)
			}

			if hit.TextMatch != nil {
				epHit.Score = float64(*hit.TextMatch)
			}

			if params.IncludeHighlights && hit.Highlights != nil {
				for _, hl := range *hit.Highlights {
					if hl.Field != nil && hl.Snippets != nil {
						epHit.Highlights[*hl.Field] = *hl.Snippets
					}
				}
			}

			searchResult.Hits = append(searchResult.Hits, epHit)
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

// AutocompleteEpisodes provides search suggestions for episode titles.
func (s *EpisodeSearchService) AutocompleteEpisodes(ctx context.Context, query string, limit int) ([]string, error) {
	if !s.IsEnabled() {
		return nil, nil
	}

	if limit < 1 {
		limit = 5
	}
	if limit > 20 {
		limit = 20
	}

	queryBy := "title,series_title"
	perPage := limit

	searchParams := &api.SearchCollectionParams{
		Q:                   &query,
		QueryBy:             &queryBy,
		PerPage:             &perPage,
		Prefix:              ptr.To("true"),
		DropTokensThreshold: ptr.To(0),
	}

	result, err := s.client.Search(ctx, EpisodeCollectionName, searchParams)
	if err != nil {
		return nil, fmt.Errorf("episode autocomplete failed: %w", err)
	}

	suggestions := make([]string, 0, limit)
	if result.Hits != nil {
		seen := make(map[string]bool)
		for _, hit := range *result.Hits {
			if hit.Document != nil {
				doc := parseEpisodeDocument(*hit.Document)
				label := fmt.Sprintf("%s - S%02dE%02d - %s", doc.SeriesTitle, doc.SeasonNumber, doc.EpisodeNumber, doc.Title)
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

// ReindexAll reindexes all episodes in the database.
func (s *EpisodeSearchService) ReindexAll(ctx context.Context, service tvshow.Service) error {
	if !s.IsEnabled() {
		return nil
	}

	s.logger.Info("starting full episode reindex")

	_ = s.client.DeleteCollection(ctx, EpisodeCollectionName)
	if err := s.InitializeCollection(ctx); err != nil {
		return fmt.Errorf("failed to reinitialize episodes collection: %w", err)
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
			episodes, err := service.ListEpisodesBySeries(ctx, series.ID)
			if err != nil {
				s.logger.Warn("failed to list episodes for series",
					slog.String("series_id", series.ID.String()),
					slog.Any("error", err),
				)
				continue
			}

			if len(episodes) == 0 {
				continue
			}

			batch := make([]EpisodeWithContext, 0, len(episodes))
			for i := range episodes {
				hasFile := false
				files, err := service.ListEpisodeFiles(ctx, episodes[i].ID)
				if err == nil && len(files) > 0 {
					hasFile = true
				}

				posterPath := ""
				if series.PosterPath != nil {
					posterPath = *series.PosterPath
				}

				batch = append(batch, EpisodeWithContext{
					Episode:          &episodes[i],
					SeriesTitle:      series.Title,
					SeriesPosterPath: posterPath,
					HasFile:          hasFile,
				})
			}

			if err := s.BulkIndexEpisodes(ctx, batch); err != nil {
				s.logger.Error("failed to index episode batch",
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

	s.logger.Info("completed full episode reindex", "total", totalIndexed)
	return nil
}

// episodeToDocument converts an episode with series context to a search document.
func (s *EpisodeSearchService) episodeToDocument(ep *tvshow.Episode, seriesTitle, seriesPosterPath string, hasFile bool) EpisodeDocument {
	now := time.Now().Unix()

	doc := EpisodeDocument{
		ID:               ep.ID.String(),
		SeriesID:         ep.SeriesID.String(),
		SeasonID:         ep.SeasonID.String(),
		SeasonNumber:     ep.SeasonNumber,
		EpisodeNumber:    ep.EpisodeNumber,
		Title:            ep.Title,
		HasFile:          hasFile,
		SeriesTitle:      seriesTitle,
		SeriesPosterPath: seriesPosterPath,
		CreatedAt:        ep.CreatedAt.Unix(),
		UpdatedAt:        now,
	}

	if ep.TMDbID != nil {
		doc.TMDbID = *ep.TMDbID
	}
	if ep.TVDbID != nil {
		doc.TVDbID = *ep.TVDbID
	}
	if ep.IMDbID != nil {
		doc.IMDbID = *ep.IMDbID
	}
	if ep.Overview != nil {
		doc.Overview = *ep.Overview
	}
	if ep.AirDate != nil {
		doc.AirDate = ep.AirDate.Unix()
	}
	if ep.Runtime != nil {
		doc.Runtime = *ep.Runtime
	}
	if ep.VoteAverage != nil {
		f, _ := ep.VoteAverage.Float64()
		doc.VoteAverage = f
	}
	if ep.VoteCount != nil {
		doc.VoteCount = *ep.VoteCount
	}
	if ep.StillPath != nil {
		doc.StillPath = *ep.StillPath
	}

	return doc
}

// parseEpisodeDocument parses a map into an EpisodeDocument.
func parseEpisodeDocument(data map[string]interface{}) EpisodeDocument {
	doc := EpisodeDocument{}

	if v, ok := data["id"].(string); ok {
		doc.ID = v
	}
	if v, ok := data["series_id"].(string); ok {
		doc.SeriesID = v
	}
	if v, ok := data["season_id"].(string); ok {
		doc.SeasonID = v
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
	if v, ok := data["season_number"].(float64); ok {
		doc.SeasonNumber = int32(v)
	}
	if v, ok := data["episode_number"].(float64); ok {
		doc.EpisodeNumber = int32(v)
	}
	if v, ok := data["title"].(string); ok {
		doc.Title = v
	}
	if v, ok := data["overview"].(string); ok {
		doc.Overview = v
	}
	if v, ok := data["air_date"].(float64); ok {
		doc.AirDate = int64(v)
	}
	if v, ok := data["runtime"].(float64); ok {
		doc.Runtime = int32(v)
	}
	if v, ok := data["vote_average"].(float64); ok {
		doc.VoteAverage = v
	}
	if v, ok := data["vote_count"].(float64); ok {
		doc.VoteCount = int32(v)
	}
	if v, ok := data["still_path"].(string); ok {
		doc.StillPath = v
	}
	if v, ok := data["has_file"].(bool); ok {
		doc.HasFile = v
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

// derefHits safely dereferences a hits pointer, returning an empty slice if nil.
func derefHits(hits *[]api.SearchResultHit) []api.SearchResultHit {
	if hits == nil {
		return nil
	}
	return *hits
}
