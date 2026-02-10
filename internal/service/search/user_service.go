package search

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/infra/search"
	"github.com/lusoris/revenge/internal/util/ptr"
	"github.com/typesense/typesense-go/v2/typesense/api"
)

// UserSearchResult contains the search results for a user query.
type UserSearchResult struct {
	TotalHits   int
	CurrentPage int
	TotalPages  int
	SearchTime  time.Duration
	Hits        []UserHit
	Facets      map[string][]FacetValue
}

// UserHit represents a single user search result.
type UserHit struct {
	Document   UserDocument
	Score      float64
	Highlights map[string][]string
}

// UserSearchParams defines parameters for searching users.
type UserSearchParams struct {
	Query             string
	Page              int
	PerPage           int
	SortBy            string
	FilterBy          string
	FacetBy           []string
	IncludeHighlights bool
}

// DefaultUserSearchParams returns sensible defaults for user search.
func DefaultUserSearchParams() UserSearchParams {
	return UserSearchParams{
		Page:              1,
		PerPage:           20,
		SortBy:            "created_at:desc",
		IncludeHighlights: true,
		FacetBy:           []string{"is_active", "is_admin"},
	}
}

// UserSearchService provides user search operations using Typesense.
type UserSearchService struct {
	client *search.Client
	logger *slog.Logger
}

// NewUserSearchService creates a new user search service.
func NewUserSearchService(client *search.Client, logger *slog.Logger) *UserSearchService {
	return &UserSearchService{
		client: client,
		logger: logger.With("service", "user_search"),
	}
}

// IsEnabled returns true if search is enabled.
func (s *UserSearchService) IsEnabled() bool {
	return s.client != nil && s.client.IsEnabled()
}

// InitializeCollection creates the users collection if it doesn't exist.
func (s *UserSearchService) InitializeCollection(ctx context.Context) error {
	if !s.IsEnabled() {
		s.logger.Debug("search disabled, skipping collection initialization")
		return nil
	}

	_, err := s.client.GetCollection(ctx, UserCollectionName)
	if err == nil {
		s.logger.Debug("users collection already exists")
		return nil
	}

	schema := UserCollectionSchema()
	if err := s.client.CreateCollection(ctx, schema); err != nil {
		return fmt.Errorf("failed to create users collection: %w", err)
	}

	s.logger.Info("created users collection")
	return nil
}

// IndexUser indexes a single user in Typesense.
func (s *UserSearchService) IndexUser(ctx context.Context, user UserDocument) error {
	if !s.IsEnabled() {
		return nil
	}

	_, err := s.client.IndexDocument(ctx, UserCollectionName, user)
	return err
}

// UpdateUser updates a user document in Typesense.
func (s *UserSearchService) UpdateUser(ctx context.Context, user UserDocument) error {
	if !s.IsEnabled() {
		return nil
	}

	_, err := s.client.UpdateDocument(ctx, UserCollectionName, user.ID, user)
	return err
}

// RemoveUser removes a user from the search index by UUID.
func (s *UserSearchService) RemoveUser(ctx context.Context, userID uuid.UUID) error {
	if !s.IsEnabled() {
		return nil
	}

	_, err := s.client.DeleteDocument(ctx, UserCollectionName, userID.String())
	return err
}

// BulkIndexUsers indexes multiple users at once using batch import.
func (s *UserSearchService) BulkIndexUsers(ctx context.Context, users []UserDocument) error {
	if !s.IsEnabled() {
		return nil
	}

	if len(users) == 0 {
		return nil
	}

	documents := make([]interface{}, 0, len(users))
	for _, u := range users {
		documents = append(documents, u)
	}

	results, err := s.client.ImportDocuments(ctx, UserCollectionName, documents, "upsert")
	if err != nil {
		return fmt.Errorf("failed to bulk index users: %w", err)
	}

	var successCount, errorCount int
	for _, result := range results {
		if result.Success {
			successCount++
		} else {
			errorCount++
			if result.Error != "" {
				s.logger.Warn("failed to index user document", "error", result.Error)
			}
		}
	}

	s.logger.Info("bulk indexed users", "success", successCount, "errors", errorCount)
	return nil
}

// SearchUsers performs a full-text search across users.
func (s *UserSearchService) SearchUsers(ctx context.Context, params UserSearchParams) (*UserSearchResult, error) {
	if !s.IsEnabled() {
		return nil, nil
	}

	queryBy := "username,email,display_name"
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
	result, err := s.client.Search(ctx, UserCollectionName, searchParams)
	searchTime := time.Since(start)

	if err != nil {
		return nil, fmt.Errorf("user search failed: %w", err)
	}

	searchResult := &UserSearchResult{
		TotalHits:   deref(result.Found),
		CurrentPage: params.Page,
		TotalPages:  (deref(result.Found) + params.PerPage - 1) / params.PerPage,
		SearchTime:  searchTime,
		Hits:        make([]UserHit, 0, len(derefHits(result.Hits))),
		Facets:      make(map[string][]FacetValue),
	}

	if result.Hits != nil {
		for _, hit := range *result.Hits {
			uHit := UserHit{
				Highlights: make(map[string][]string),
			}

			if hit.Document != nil {
				uHit.Document = parseUserDocument(*hit.Document)
			}

			if hit.TextMatch != nil {
				uHit.Score = float64(*hit.TextMatch)
			}

			if params.IncludeHighlights && hit.Highlights != nil {
				for _, hl := range *hit.Highlights {
					if hl.Field != nil && hl.Snippets != nil {
						uHit.Highlights[*hl.Field] = *hl.Snippets
					}
				}
			}

			searchResult.Hits = append(searchResult.Hits, uHit)
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

// AutocompleteUsers provides search suggestions for usernames.
func (s *UserSearchService) AutocompleteUsers(ctx context.Context, query string, limit int) ([]string, error) {
	if !s.IsEnabled() {
		return nil, nil
	}

	if limit < 1 {
		limit = 5
	}
	if limit > 20 {
		limit = 20
	}

	queryBy := "username,display_name"
	perPage := limit

	searchParams := &api.SearchCollectionParams{
		Q:                   &query,
		QueryBy:             &queryBy,
		PerPage:             &perPage,
		Prefix:              ptr.To("true"),
		DropTokensThreshold: ptr.To(0),
	}

	result, err := s.client.Search(ctx, UserCollectionName, searchParams)
	if err != nil {
		return nil, fmt.Errorf("user autocomplete failed: %w", err)
	}

	suggestions := make([]string, 0, limit)
	if result.Hits != nil {
		seen := make(map[string]bool)
		for _, hit := range *result.Hits {
			if hit.Document != nil {
				doc := parseUserDocument(*hit.Document)
				if !seen[doc.Username] {
					suggestions = append(suggestions, doc.Username)
					seen[doc.Username] = true
				}
			}
			if len(suggestions) >= limit {
				break
			}
		}
	}

	return suggestions, nil
}

// ReindexAll reindexes all users from the provided user documents.
// The caller is responsible for fetching and converting users to UserDocument.
func (s *UserSearchService) ReindexAll(ctx context.Context, users []UserDocument) error {
	if !s.IsEnabled() {
		return nil
	}

	s.logger.Info("starting full user reindex")

	_ = s.client.DeleteCollection(ctx, UserCollectionName)
	if err := s.InitializeCollection(ctx); err != nil {
		return fmt.Errorf("failed to reinitialize users collection: %w", err)
	}

	batchSize := 100
	totalIndexed := 0

	for i := 0; i < len(users); i += batchSize {
		end := i + batchSize
		if end > len(users) {
			end = len(users)
		}

		batch := users[i:end]
		if err := s.BulkIndexUsers(ctx, batch); err != nil {
			s.logger.Error("failed to index user batch",
				slog.Int("offset", i),
				slog.Any("error", err),
			)
		}

		totalIndexed += len(batch)
	}

	s.logger.Info("completed full user reindex", slog.Int("total", totalIndexed))
	return nil
}

// parseUserDocument converts a raw Typesense document map to a UserDocument.
func parseUserDocument(data map[string]interface{}) UserDocument {
	doc := UserDocument{}

	if v, ok := data["id"].(string); ok {
		doc.ID = v
	}
	if v, ok := data["username"].(string); ok {
		doc.Username = v
	}
	if v, ok := data["email"].(string); ok {
		doc.Email = v
	}
	if v, ok := data["display_name"].(string); ok {
		doc.DisplayName = v
	}
	if v, ok := data["avatar_url"].(string); ok {
		doc.AvatarURL = v
	}
	if v, ok := data["is_active"].(bool); ok {
		doc.IsActive = v
	}
	if v, ok := data["is_admin"].(bool); ok {
		doc.IsAdmin = v
	}
	if v, ok := data["created_at"].(float64); ok {
		doc.CreatedAt = int64(v)
	}
	if v, ok := data["last_login_at"].(float64); ok {
		doc.LastLoginAt = int64(v)
	}

	return doc
}
