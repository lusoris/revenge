package search

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/lusoris/revenge/internal/infra/search"
	"github.com/lusoris/revenge/internal/util/ptr"
	"github.com/typesense/typesense-go/v2/typesense/api"
)

// PersonAggregate represents aggregated person data from all credit sources.
// Callers build this from movie and TV show credits before indexing.
type PersonAggregate struct {
	TMDbPersonID int32
	Name         string
	ProfilePath  string
	KnownFor     []string // Movie/show titles
	Characters   []string // Character names
	Departments  []string // Unique departments (Acting, Directing, etc.)
	MovieCount   int32
	TVShowCount  int32
}

// PersonSearchResult contains the search results for a person query.
type PersonSearchResult struct {
	TotalHits   int
	CurrentPage int
	TotalPages  int
	SearchTime  time.Duration
	Hits        []PersonHit
	Facets      map[string][]FacetValue
}

// PersonHit represents a single person search result.
type PersonHit struct {
	Document   PersonDocument
	Score      float64
	Highlights map[string][]string
}

// PersonSearchParams defines parameters for searching people.
type PersonSearchParams struct {
	Query             string
	Page              int
	PerPage           int
	SortBy            string
	FilterBy          string
	FacetBy           []string
	IncludeHighlights bool
}

// DefaultPersonSearchParams returns sensible defaults for person search.
func DefaultPersonSearchParams() PersonSearchParams {
	return PersonSearchParams{
		Page:              1,
		PerPage:           20,
		SortBy:            "total_credits:desc",
		IncludeHighlights: true,
		FacetBy:           []string{"departments"},
	}
}

// PersonSearchService provides person search operations using Typesense.
type PersonSearchService struct {
	client *search.Client
	logger *slog.Logger
}

// NewPersonSearchService creates a new person search service.
func NewPersonSearchService(client *search.Client, logger *slog.Logger) *PersonSearchService {
	return &PersonSearchService{
		client: client,
		logger: logger.With("service", "person_search"),
	}
}

// IsEnabled returns true if search is enabled.
func (s *PersonSearchService) IsEnabled() bool {
	return s.client != nil && s.client.IsEnabled()
}

// InitializeCollection creates the people collection if it doesn't exist.
func (s *PersonSearchService) InitializeCollection(ctx context.Context) error {
	if !s.IsEnabled() {
		s.logger.Debug("search disabled, skipping collection initialization")
		return nil
	}

	_, err := s.client.GetCollection(ctx, PersonCollectionName)
	if err == nil {
		s.logger.Debug("people collection already exists")
		return nil
	}

	schema := PersonCollectionSchema()
	if err := s.client.CreateCollection(ctx, schema); err != nil {
		return fmt.Errorf("failed to create people collection: %w", err)
	}

	s.logger.Info("created people collection")
	return nil
}

// IndexPerson indexes a single person in Typesense.
func (s *PersonSearchService) IndexPerson(ctx context.Context, person PersonAggregate) error {
	if !s.IsEnabled() {
		return nil
	}

	doc := personToDocument(person)
	_, err := s.client.IndexDocument(ctx, PersonCollectionName, doc)
	return err
}

// UpdatePerson updates a person document in Typesense.
func (s *PersonSearchService) UpdatePerson(ctx context.Context, person PersonAggregate) error {
	if !s.IsEnabled() {
		return nil
	}

	doc := personToDocument(person)
	_, err := s.client.UpdateDocument(ctx, PersonCollectionName, doc.ID, doc)
	return err
}

// RemovePerson removes a person from the search index by TMDb person ID.
func (s *PersonSearchService) RemovePerson(ctx context.Context, tmdbPersonID int32) error {
	if !s.IsEnabled() {
		return nil
	}

	docID := fmt.Sprintf("%d", tmdbPersonID)
	_, err := s.client.DeleteDocument(ctx, PersonCollectionName, docID)
	return err
}

// BulkIndexPersons indexes multiple people at once using batch import.
func (s *PersonSearchService) BulkIndexPersons(ctx context.Context, people []PersonAggregate) error {
	if !s.IsEnabled() {
		return nil
	}

	if len(people) == 0 {
		return nil
	}

	documents := make([]interface{}, 0, len(people))
	for _, p := range people {
		documents = append(documents, personToDocument(p))
	}

	results, err := s.client.ImportDocuments(ctx, PersonCollectionName, documents, "upsert")
	if err != nil {
		return fmt.Errorf("failed to bulk index people: %w", err)
	}

	var successCount, errorCount int
	for _, result := range results {
		if result.Success {
			successCount++
		} else {
			errorCount++
			if result.Error != "" {
				s.logger.Warn("failed to index person document", "error", result.Error)
			}
		}
	}

	s.logger.Info("bulk indexed people", "success", successCount, "errors", errorCount)
	return nil
}

// SearchPersons performs a full-text search across people.
func (s *PersonSearchService) SearchPersons(ctx context.Context, params PersonSearchParams) (*PersonSearchResult, error) {
	if !s.IsEnabled() {
		return nil, nil
	}

	queryBy := "name,known_for,characters"
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
	result, err := s.client.Search(ctx, PersonCollectionName, searchParams)
	searchTime := time.Since(start)

	if err != nil {
		return nil, fmt.Errorf("person search failed: %w", err)
	}

	searchResult := &PersonSearchResult{
		TotalHits:   deref(result.Found),
		CurrentPage: params.Page,
		TotalPages:  (deref(result.Found) + params.PerPage - 1) / params.PerPage,
		SearchTime:  searchTime,
		Hits:        make([]PersonHit, 0, len(derefHits(result.Hits))),
		Facets:      make(map[string][]FacetValue),
	}

	if result.Hits != nil {
		for _, hit := range *result.Hits {
			pHit := PersonHit{
				Highlights: make(map[string][]string),
			}

			if hit.Document != nil {
				pHit.Document = parsePersonDocument(*hit.Document)
			}

			if hit.TextMatch != nil {
				pHit.Score = float64(*hit.TextMatch)
			}

			if params.IncludeHighlights && hit.Highlights != nil {
				for _, hl := range *hit.Highlights {
					if hl.Field != nil && hl.Snippets != nil {
						pHit.Highlights[*hl.Field] = *hl.Snippets
					}
				}
			}

			searchResult.Hits = append(searchResult.Hits, pHit)
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

// AutocompletePersons provides search suggestions for person names.
func (s *PersonSearchService) AutocompletePersons(ctx context.Context, query string, limit int) ([]string, error) {
	if !s.IsEnabled() {
		return nil, nil
	}

	if limit < 1 {
		limit = 5
	}
	if limit > 20 {
		limit = 20
	}

	queryBy := "name"
	perPage := limit

	searchParams := &api.SearchCollectionParams{
		Q:                   &query,
		QueryBy:             &queryBy,
		PerPage:             &perPage,
		Prefix:              ptr.To("true"),
		DropTokensThreshold: ptr.To(0),
	}

	result, err := s.client.Search(ctx, PersonCollectionName, searchParams)
	if err != nil {
		return nil, fmt.Errorf("person autocomplete failed: %w", err)
	}

	suggestions := make([]string, 0, limit)
	if result.Hits != nil {
		seen := make(map[string]bool)
		for _, hit := range *result.Hits {
			if hit.Document != nil {
				doc := parsePersonDocument(*hit.Document)
				if !seen[doc.Name] {
					suggestions = append(suggestions, doc.Name)
					seen[doc.Name] = true
				}
			}
			if len(suggestions) >= limit {
				break
			}
		}
	}

	return suggestions, nil
}

// ReindexAll reindexes all people from pre-aggregated data.
// The caller is responsible for aggregating person data from movie and TV show credits.
func (s *PersonSearchService) ReindexAll(ctx context.Context, people []PersonAggregate) error {
	if !s.IsEnabled() {
		return nil
	}

	s.logger.Info("starting full person reindex")

	_ = s.client.DeleteCollection(ctx, PersonCollectionName)
	if err := s.InitializeCollection(ctx); err != nil {
		return fmt.Errorf("failed to reinitialize people collection: %w", err)
	}

	batchSize := 100
	totalIndexed := 0

	for i := 0; i < len(people); i += batchSize {
		end := i + batchSize
		if end > len(people) {
			end = len(people)
		}

		batch := people[i:end]
		if err := s.BulkIndexPersons(ctx, batch); err != nil {
			s.logger.Error("failed to index person batch",
				slog.Int("offset", i),
				slog.Any("error", err),
			)
		}

		totalIndexed += len(batch)
	}

	s.logger.Info("completed full person reindex", slog.Int("total", totalIndexed))
	return nil
}

// personToDocument converts a PersonAggregate to a PersonDocument.
func personToDocument(p PersonAggregate) PersonDocument {
	doc := PersonDocument{
		ID:           fmt.Sprintf("%d", p.TMDbPersonID),
		TMDbID:       p.TMDbPersonID,
		Name:         p.Name,
		ProfilePath:  p.ProfilePath,
		MovieCount:   p.MovieCount,
		TVShowCount:  p.TVShowCount,
		TotalCredits: p.MovieCount + p.TVShowCount,
	}

	if len(p.KnownFor) > 0 {
		// Deduplicate and cap at 20 titles
		seen := make(map[string]bool)
		knownFor := make([]string, 0, len(p.KnownFor))
		for _, title := range p.KnownFor {
			if !seen[title] {
				knownFor = append(knownFor, title)
				seen[title] = true
			}
		}
		if len(knownFor) > 20 {
			knownFor = knownFor[:20]
		}
		doc.KnownFor = knownFor
	}

	if len(p.Characters) > 0 {
		seen := make(map[string]bool)
		chars := make([]string, 0, len(p.Characters))
		for _, c := range p.Characters {
			if !seen[c] {
				chars = append(chars, c)
				seen[c] = true
			}
		}
		if len(chars) > 20 {
			chars = chars[:20]
		}
		doc.Characters = chars
	}

	if len(p.Departments) > 0 {
		seen := make(map[string]bool)
		depts := make([]string, 0, len(p.Departments))
		for _, d := range p.Departments {
			if !seen[d] {
				depts = append(depts, d)
				seen[d] = true
			}
		}
		doc.Departments = depts
	}

	return doc
}

// parsePersonDocument converts a raw Typesense document map to a PersonDocument.
func parsePersonDocument(data map[string]interface{}) PersonDocument {
	doc := PersonDocument{}

	if v, ok := data["id"].(string); ok {
		doc.ID = v
	}
	if v, ok := data["tmdb_id"].(float64); ok {
		doc.TMDbID = int32(v)
	}
	if v, ok := data["name"].(string); ok {
		doc.Name = v
	}
	if v, ok := data["profile_path"].(string); ok {
		doc.ProfilePath = v
	}
	if v, ok := data["known_for"].([]interface{}); ok {
		doc.KnownFor = toStringSlice(v)
	}
	if v, ok := data["characters"].([]interface{}); ok {
		doc.Characters = toStringSlice(v)
	}
	if v, ok := data["departments"].([]interface{}); ok {
		doc.Departments = toStringSlice(v)
	}
	if v, ok := data["movie_count"].(float64); ok {
		doc.MovieCount = int32(v)
	}
	if v, ok := data["tvshow_count"].(float64); ok {
		doc.TVShowCount = int32(v)
	}
	if v, ok := data["total_credits"].(float64); ok {
		doc.TotalCredits = int32(v)
	}

	return doc
}
