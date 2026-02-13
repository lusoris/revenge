package search

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/content/movie"
	"github.com/lusoris/revenge/internal/infra/observability"
	"github.com/lusoris/revenge/internal/infra/search"
	"github.com/typesense/typesense-go/v2/typesense/api"
)

// MovieSearchService provides movie search operations using Typesense.
type MovieSearchService struct {
	client *search.Client
	logger *slog.Logger
}

// NewMovieSearchService creates a new movie search service.
func NewMovieSearchService(client *search.Client, logger *slog.Logger) *MovieSearchService {
	return &MovieSearchService{
		client: client,
		logger: logger.With("service", "movie_search"),
	}
}

// IsEnabled returns true if search is enabled.
func (s *MovieSearchService) IsEnabled() bool {
	return s.client != nil && s.client.IsEnabled()
}

// InitializeCollection creates the movies collection if it doesn't exist.
func (s *MovieSearchService) InitializeCollection(ctx context.Context) error {
	if !s.IsEnabled() {
		s.logger.Debug("search disabled, skipping collection initialization")
		return nil
	}

	// Check if collection already exists
	_, err := s.client.GetCollection(ctx, MovieCollectionName)
	if err == nil {
		s.logger.Debug("movies collection already exists")
		return nil
	}

	// Create collection
	schema := MovieCollectionSchema()
	if err := s.client.CreateCollection(ctx, schema); err != nil {
		return fmt.Errorf("failed to create movies collection: %w", err)
	}

	s.logger.Info("created movies collection")
	return nil
}

// IndexMovie indexes a single movie in Typesense.
func (s *MovieSearchService) IndexMovie(ctx context.Context, m *movie.Movie, genres []movie.MovieGenre, credits []movie.MovieCredit, file *movie.MovieFile) error {
	if !s.IsEnabled() {
		return nil
	}

	doc := s.movieToDocument(m, genres, credits, file)
	_, err := s.client.IndexDocument(ctx, MovieCollectionName, doc)
	if err != nil {
		return fmt.Errorf("failed to index movie %s: %w", m.ID, err)
	}

	s.logger.Debug("indexed movie", "id", m.ID, "title", m.Title)
	return nil
}

// UpdateMovie updates a movie in the search index.
func (s *MovieSearchService) UpdateMovie(ctx context.Context, m *movie.Movie, genres []movie.MovieGenre, credits []movie.MovieCredit, file *movie.MovieFile) error {
	if !s.IsEnabled() {
		return nil
	}

	doc := s.movieToDocument(m, genres, credits, file)
	_, err := s.client.UpdateDocument(ctx, MovieCollectionName, m.ID.String(), doc)
	if err != nil {
		// If document doesn't exist, create it
		if strings.Contains(err.Error(), "not found") {
			return s.IndexMovie(ctx, m, genres, credits, file)
		}
		return fmt.Errorf("failed to update movie %s: %w", m.ID, err)
	}

	s.logger.Debug("updated movie", "id", m.ID, "title", m.Title)
	return nil
}

// RemoveMovie removes a movie from the search index.
func (s *MovieSearchService) RemoveMovie(ctx context.Context, movieID uuid.UUID) error {
	if !s.IsEnabled() {
		return nil
	}

	_, err := s.client.DeleteDocument(ctx, MovieCollectionName, movieID.String())
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return fmt.Errorf("failed to remove movie %s: %w", movieID, err)
	}

	s.logger.Debug("removed movie", "id", movieID)
	return nil
}

// BulkIndexMovies indexes multiple movies at once.
func (s *MovieSearchService) BulkIndexMovies(ctx context.Context, movies []MovieWithRelations) error {
	if !s.IsEnabled() {
		return nil
	}

	if len(movies) == 0 {
		return nil
	}

	docs := make([]any, len(movies))
	for i, m := range movies {
		docs[i] = s.movieToDocument(m.Movie, m.Genres, m.Credits, m.File)
	}

	results, err := s.client.ImportDocuments(ctx, MovieCollectionName, docs, "upsert")
	if err != nil {
		return fmt.Errorf("failed to bulk index movies: %w", err)
	}

	// Count successes and failures
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

	s.logger.Info("bulk indexed movies", "success", successCount, "errors", errorCount)

	if errorCount > 0 {
		return fmt.Errorf("bulk index completed with %d errors out of %d documents", errorCount, successCount+errorCount)
	}
	return nil
}

// MovieWithRelations bundles a movie with its related data for indexing.
type MovieWithRelations struct {
	Movie   *movie.Movie
	Genres  []movie.MovieGenre
	Credits []movie.MovieCredit
	File    *movie.MovieFile
}

// SearchResult represents a search result with hits and facets.
type SearchResult struct {
	Hits        []MovieHit
	TotalHits   int
	TotalPages  int
	CurrentPage int
	Facets      map[string][]FacetValue
	SearchTime  time.Duration
}

// MovieHit represents a single movie hit in search results.
type MovieHit struct {
	Document   MovieDocument
	Score      float64
	Highlights map[string][]string
}

// FacetValue represents a facet value with count.
type FacetValue struct {
	Value string
	Count int
}

// SearchParams contains parameters for movie search.
type SearchParams struct {
	Query             string
	Page              int
	PerPage           int
	SortBy            string
	FilterBy          string
	FacetBy           []string
	IncludeHighlights bool
}

// DefaultSearchParams returns default search parameters.
func DefaultSearchParams() SearchParams {
	return SearchParams{
		Page:              1,
		PerPage:           20,
		SortBy:            "popularity:desc",
		FacetBy:           []string{"genres", "year", "status", "directors", "resolution", "has_file"},
		IncludeHighlights: true,
	}
}

// Search searches for movies with the given parameters.
func (s *MovieSearchService) Search(ctx context.Context, params SearchParams) (*SearchResult, error) {
	if !s.IsEnabled() {
		return &SearchResult{}, nil
	}

	// Set defaults
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PerPage < 1 {
		params.PerPage = 20
	}
	if params.PerPage > 100 {
		params.PerPage = 100
	}

	// Build search parameters
	queryBy := "title,original_title,overview,cast,directors"
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

	// Execute search
	start := time.Now()
	result, err := s.client.Search(ctx, MovieCollectionName, searchParams)
	searchTime := time.Since(start)
	observability.SearchQueriesTotal.WithLabelValues("search").Inc()
	observability.SearchQueryDuration.WithLabelValues("search").Observe(searchTime.Seconds())

	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	// Convert to our result type
	searchResult := &SearchResult{
		TotalHits:   deref(result.Found),
		CurrentPage: params.Page,
		TotalPages:  (deref(result.Found) + params.PerPage - 1) / params.PerPage,
		SearchTime:  searchTime,
		Hits:        make([]MovieHit, 0, len(*result.Hits)),
		Facets:      make(map[string][]FacetValue),
	}

	// Convert hits
	if result.Hits != nil {
		for _, hit := range *result.Hits {
			movieHit := MovieHit{
				Highlights: make(map[string][]string),
			}

			// Parse document
			if hit.Document != nil {
				movieHit.Document = parseMovieDocument(*hit.Document)
			}

			// Get text match score
			if hit.TextMatch != nil {
				movieHit.Score = float64(*hit.TextMatch)
			}

			// Parse highlights
			if params.IncludeHighlights && hit.Highlights != nil {
				for _, hl := range *hit.Highlights {
					if hl.Field != nil && hl.Snippets != nil {
						movieHit.Highlights[*hl.Field] = *hl.Snippets
					}
				}
			}

			searchResult.Hits = append(searchResult.Hits, movieHit)
		}
	}

	// Convert facets
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

// Autocomplete provides search suggestions for movie titles.
func (s *MovieSearchService) Autocomplete(ctx context.Context, query string, limit int) ([]string, error) {
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
		Prefix:              new("true"),
		DropTokensThreshold: new(0),
	}

	start := time.Now()
	result, err := s.client.Search(ctx, MovieCollectionName, searchParams)
	observability.SearchQueriesTotal.WithLabelValues("autocomplete").Inc()
	observability.SearchQueryDuration.WithLabelValues("autocomplete").Observe(time.Since(start).Seconds())
	if err != nil {
		return nil, fmt.Errorf("autocomplete failed: %w", err)
	}

	suggestions := make([]string, 0, limit)
	if result.Hits != nil {
		seen := make(map[string]bool)
		for _, hit := range *result.Hits {
			if hit.Document != nil {
				doc := parseMovieDocument(*hit.Document)
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
func (s *MovieSearchService) GetFacets(ctx context.Context, facetNames []string) (map[string][]FacetValue, error) {
	if !s.IsEnabled() {
		return nil, nil
	}

	// Use an empty query with facets to get all facet values
	params := SearchParams{
		Query:   "*",
		Page:    1,
		PerPage: 0, // No hits, just facets
		FacetBy: facetNames,
	}

	result, err := s.Search(ctx, params)
	if err != nil {
		return nil, err
	}

	return result.Facets, nil
}

// ReindexAll reindexes all movies in the database.
func (s *MovieSearchService) ReindexAll(ctx context.Context, movieRepo movie.Repository) error {
	if !s.IsEnabled() {
		return nil
	}

	s.logger.Info("starting full reindex")

	// Delete and recreate collection
	_ = s.client.DeleteCollection(ctx, MovieCollectionName)
	if err := s.InitializeCollection(ctx); err != nil {
		return fmt.Errorf("failed to reinitialize collection: %w", err)
	}

	// Get all movies in batches
	batchSize := int32(100)
	offset := int32(0)
	totalIndexed := 0

	for {
		movies, err := movieRepo.ListMovies(ctx, movie.ListFilters{
			Limit:  batchSize,
			Offset: offset,
		})
		if err != nil {
			return fmt.Errorf("failed to list movies: %w", err)
		}

		if len(movies) == 0 {
			break
		}

		// Get related data for each movie
		moviesWithRelations := make([]MovieWithRelations, 0, len(movies))
		for _, m := range movies {
			genres, _ := movieRepo.ListMovieGenres(ctx, m.ID)
			cast, _ := movieRepo.ListMovieCast(ctx, m.ID, 1000, 0)
			crew, _ := movieRepo.ListMovieCrew(ctx, m.ID, 1000, 0)
			files, _ := movieRepo.ListMovieFilesByMovieID(ctx, m.ID)

			// Combine cast and crew into credits
			credits := append(cast, crew...)

			// Get first file as primary
			var file *movie.MovieFile
			if len(files) > 0 {
				file = &files[0]
			}

			mCopy := m // Create a copy to take address of
			moviesWithRelations = append(moviesWithRelations, MovieWithRelations{
				Movie:   &mCopy,
				Genres:  genres,
				Credits: credits,
				File:    file,
			})
		}

		// Bulk index
		if err := s.BulkIndexMovies(ctx, moviesWithRelations); err != nil {
			s.logger.Error("failed to index batch", "offset", offset, "error", err)
		}

		totalIndexed += len(movies)
		offset += batchSize

		if len(movies) < int(batchSize) {
			break
		}
	}

	s.logger.Info("completed full reindex", "total", totalIndexed)
	return nil
}

// movieToDocument converts a movie with related data to a search document.
func (s *MovieSearchService) movieToDocument(m *movie.Movie, genres []movie.MovieGenre, credits []movie.MovieCredit, file *movie.MovieFile) MovieDocument {
	now := time.Now().Unix()

	doc := MovieDocument{
		ID:        m.ID.String(),
		Title:     m.Title,
		CreatedAt: m.CreatedAt.Unix(),
		UpdatedAt: now,
	}

	// TMDb ID
	if m.TMDbID != nil {
		doc.TMDbID = *m.TMDbID
	}

	// IMDb ID
	if m.IMDbID != nil {
		doc.IMDbID = *m.IMDbID
	}

	// Original title
	if m.OriginalTitle != nil {
		doc.OriginalTitle = *m.OriginalTitle
	}

	// Year
	if m.Year != nil {
		doc.Year = *m.Year
	}

	// Release date
	if m.ReleaseDate != nil {
		doc.ReleaseDate = m.ReleaseDate.Unix()
	}

	// Runtime
	if m.Runtime != nil {
		doc.Runtime = *m.Runtime
	}

	// Overview
	if m.Overview != nil {
		doc.Overview = *m.Overview
	}

	// Tagline
	if m.Tagline != nil {
		doc.Tagline = *m.Tagline
	}

	// Status
	if m.Status != nil {
		doc.Status = *m.Status
	}

	// Original language
	if m.OriginalLanguage != nil {
		doc.OriginalLanguage = *m.OriginalLanguage
	}

	// Images
	if m.PosterPath != nil {
		doc.PosterPath = *m.PosterPath
	}
	if m.BackdropPath != nil {
		doc.BackdropPath = *m.BackdropPath
	}

	// Ratings
	if m.VoteAverage != nil {
		f, _ := m.VoteAverage.Float64()
		doc.VoteAverage = f
	}
	if m.VoteCount != nil {
		doc.VoteCount = *m.VoteCount
	}
	if m.Popularity != nil {
		f, _ := m.Popularity.Float64()
		doc.Popularity = f
	}

	// Library added timestamp
	doc.LibraryAddedAt = m.LibraryAddedAt.Unix()

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

	// Credits - extract cast and directors
	if len(credits) > 0 {
		var cast, directors []string
		for _, c := range credits {
			if c.CreditType == "cast" {
				cast = append(cast, c.Name)
			} else if c.CreditType == "crew" && c.Job != nil && *c.Job == "Director" {
				directors = append(directors, c.Name)
			}
		}
		if len(cast) > 20 {
			cast = cast[:20] // Limit cast to top 20
		}
		doc.Cast = cast
		doc.Directors = directors
	}

	// File info
	if file != nil {
		doc.HasFile = true
		if file.Resolution != nil {
			doc.Resolution = *file.Resolution
		}
		if file.QualityProfile != nil {
			doc.QualityProfile = *file.QualityProfile
		}
	}

	return doc
}

// parseMovieDocument parses a map into a MovieDocument.
func parseMovieDocument(data map[string]any) MovieDocument {
	doc := MovieDocument{}

	if v, ok := data["id"].(string); ok {
		doc.ID = v
	}
	if v, ok := data["tmdb_id"].(float64); ok {
		doc.TMDbID = int32(v)
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
	if v, ok := data["release_date"].(float64); ok {
		doc.ReleaseDate = int64(v)
	}
	if v, ok := data["runtime"].(float64); ok {
		doc.Runtime = int32(v)
	}
	if v, ok := data["overview"].(string); ok {
		doc.Overview = v
	}
	if v, ok := data["tagline"].(string); ok {
		doc.Tagline = v
	}
	if v, ok := data["status"].(string); ok {
		doc.Status = v
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
	if v, ok := data["resolution"].(string); ok {
		doc.Resolution = v
	}
	if v, ok := data["quality_profile"].(string); ok {
		doc.QualityProfile = v
	}
	if v, ok := data["library_added_at"].(float64); ok {
		doc.LibraryAddedAt = int64(v)
	}
	if v, ok := data["created_at"].(float64); ok {
		doc.CreatedAt = int64(v)
	}
	if v, ok := data["updated_at"].(float64); ok {
		doc.UpdatedAt = int64(v)
	}

	// Parse string arrays
	if v, ok := data["genres"].([]any); ok {
		doc.Genres = toStringSlice(v)
	}
	if v, ok := data["cast"].([]any); ok {
		doc.Cast = toStringSlice(v)
	}
	if v, ok := data["directors"].([]any); ok {
		doc.Directors = toStringSlice(v)
	}
	if v, ok := data["genre_slugs"].([]any); ok {
		doc.GenreSlugs = toStringSlice(v)
	}

	return doc
}

// toStringSlice converts an interface slice to a string slice.
func toStringSlice(v []any) []string {
	result := make([]string, 0, len(v))
	for _, item := range v {
		if s, ok := item.(string); ok {
			result = append(result, s)
		}
	}
	return result
}

// toInt32Slice converts an interface slice to an int32 slice.
func toInt32Slice(v []any) []int32 {
	result := make([]int32, 0, len(v))
	for _, item := range v {
		if f, ok := item.(float64); ok {
			result = append(result, int32(f))
		}
	}
	return result
}

// deref safely dereferences a pointer, returning 0 if nil.
func deref(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}
