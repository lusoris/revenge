package api

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/content/movie/moviejobs"
	"github.com/lusoris/revenge/internal/service/search"
	"log/slog"
)

// SearchLibraryMovies searches movies in the library using Typesense.
func (h *Handler) SearchLibraryMovies(ctx context.Context, params ogen.SearchLibraryMoviesParams) (ogen.SearchLibraryMoviesRes, error) {
	// Build search params
	searchParams := search.SearchParams{
		Query:             params.Q,
		Page:              1,
		PerPage:           20,
		SortBy:            "popularity:desc",
		IncludeHighlights: true,
		FacetBy:           []string{"genres", "year", "status", "directors", "resolution", "has_file"},
	}

	if params.Page.Set {
		searchParams.Page = params.Page.Value
	}
	if params.PerPage.Set {
		searchParams.PerPage = params.PerPage.Value
	}
	if params.SortBy.Set {
		searchParams.SortBy = string(params.SortBy.Value)
	}
	if params.FilterBy.Set {
		searchParams.FilterBy = params.FilterBy.Value
	}

	// Execute search — return empty results on error (e.g. collection not found)
	result, err := h.searchService.Search(ctx, searchParams)
	if err != nil {
		h.logger.Warn("search unavailable, returning empty results", slog.Any("error",err))
		return &ogen.SearchResults{
			TotalHits:    ogen.NewOptInt(0),
			TotalPages:   ogen.NewOptInt(0),
			CurrentPage:  ogen.NewOptInt(searchParams.Page),
			SearchTimeMs: ogen.NewOptInt(0),
			Hits:         []ogen.SearchHit{},
		}, nil
	}

	// Convert to API response
	response := &ogen.SearchResults{
		TotalHits:    ogen.NewOptInt(result.TotalHits),
		TotalPages:   ogen.NewOptInt(result.TotalPages),
		CurrentPage:  ogen.NewOptInt(result.CurrentPage),
		SearchTimeMs: ogen.NewOptInt(int(result.SearchTime.Milliseconds())),
		Hits:         make([]ogen.SearchHit, 0, len(result.Hits)),
	}

	// Convert hits
	for _, hit := range result.Hits {
		apiHit := ogen.SearchHit{
			Score: ogen.NewOptFloat32(float32(hit.Score)),
		}

		// Convert document
		doc := hit.Document
		apiDoc := ogen.SearchDocument{
			ID:            ogen.NewOptUUID(parseUUID(doc.ID)),
			TmdbID:        ogen.NewOptInt(int(doc.TMDbID)),
			ImdbID:        ogen.NewOptString(doc.IMDbID),
			Title:         ogen.NewOptString(doc.Title),
			OriginalTitle: ogen.NewOptString(doc.OriginalTitle),
			Year:          ogen.NewOptInt(int(doc.Year)),
			Runtime:       ogen.NewOptInt(int(doc.Runtime)),
			Overview:      ogen.NewOptString(doc.Overview),
			Status:        ogen.NewOptString(doc.Status),
			PosterPath:    ogen.NewOptString(doc.PosterPath),
			BackdropPath:  ogen.NewOptString(doc.BackdropPath),
			VoteAverage:   ogen.NewOptFloat32(float32(doc.VoteAverage)),
			Popularity:    ogen.NewOptFloat32(float32(doc.Popularity)),
			HasFile:       ogen.NewOptBool(doc.HasFile),
			Resolution:    ogen.NewOptString(doc.Resolution),
			QualityProfile: ogen.NewOptString(doc.QualityProfile),
		}

		// Convert release date from unix timestamp
		if doc.ReleaseDate > 0 {
			t := time.Unix(doc.ReleaseDate, 0)
			apiDoc.ReleaseDate = ogen.NewOptDate(t)
		}

		// Convert arrays
		if len(doc.Genres) > 0 {
			apiDoc.Genres = make([]string, len(doc.Genres))
			copy(apiDoc.Genres, doc.Genres)
		}
		if len(doc.Cast) > 0 {
			apiDoc.Cast = make([]string, len(doc.Cast))
			copy(apiDoc.Cast, doc.Cast)
		}
		if len(doc.Directors) > 0 {
			apiDoc.Directors = make([]string, len(doc.Directors))
			copy(apiDoc.Directors, doc.Directors)
		}

		apiHit.Document = ogen.NewOptSearchDocument(apiDoc)

		// Convert highlights
		if len(hit.Highlights) > 0 {
			highlightMap := make(ogen.SearchHitHighlights)
			for field, snippets := range hit.Highlights {
				highlightMap[field] = snippets
			}
			apiHit.Highlights = ogen.OptSearchHitHighlights{
				Value: highlightMap,
				Set:   true,
			}
		}

		response.Hits = append(response.Hits, apiHit)
	}

	// Convert facets
	if len(result.Facets) > 0 {
		facets := ogen.OptSearchResultsFacets{
			Set:   true,
			Value: make(ogen.SearchResultsFacets),
		}
		for facetName, values := range result.Facets {
			apiFacets := make([]ogen.FacetValue, len(values))
			for i, v := range values {
				apiFacets[i] = ogen.FacetValue{
					Value: ogen.NewOptString(v.Value),
					Count: ogen.NewOptInt(v.Count),
				}
			}
			facets.Value[facetName] = apiFacets
		}
		response.Facets = facets
	}

	return response, nil
}

// AutocompleteMovies provides autocomplete suggestions for movie titles.
func (h *Handler) AutocompleteMovies(ctx context.Context, params ogen.AutocompleteMoviesParams) (ogen.AutocompleteMoviesRes, error) {
	limit := 5
	if params.Limit.Set {
		limit = params.Limit.Value
	}

	suggestions, err := h.searchService.Autocomplete(ctx, params.Q, limit)
	if err != nil {
		h.logger.Warn("autocomplete unavailable, returning empty", slog.Any("error",err))
		return &ogen.AutocompleteResults{
			Suggestions: []string{},
		}, nil
	}

	return &ogen.AutocompleteResults{
		Suggestions: suggestions,
	}, nil
}

// GetSearchFacets returns available facet values for filtering.
func (h *Handler) GetSearchFacets(ctx context.Context) (ogen.GetSearchFacetsRes, error) {
	facetNames := []string{"genres", "year", "status", "directors", "resolution", "has_file"}

	facets, err := h.searchService.GetFacets(ctx, facetNames)
	if err != nil {
		h.logger.Warn("facets unavailable, returning empty", slog.Any("error",err))
		return &ogen.SearchFacets{}, nil
	}

	response := &ogen.SearchFacets{}

	if values, ok := facets["genres"]; ok {
		response.Genres = convertFacetValues(values)
	}
	if values, ok := facets["year"]; ok {
		response.Years = convertFacetValues(values)
	}
	if values, ok := facets["status"]; ok {
		response.Status = convertFacetValues(values)
	}
	if values, ok := facets["directors"]; ok {
		response.Directors = convertFacetValues(values)
	}
	if values, ok := facets["resolution"]; ok {
		response.Resolution = convertFacetValues(values)
	}
	if values, ok := facets["has_file"]; ok {
		response.HasFile = convertFacetValues(values)
	}

	return response, nil
}

// ReindexSearch triggers a full reindex of all movies via River job queue.
func (h *Handler) ReindexSearch(ctx context.Context) (ogen.ReindexSearchRes, error) {
	if h.riverClient == nil {
		return nil, fmt.Errorf("job queue not available")
	}

	result, err := h.riverClient.Insert(ctx, moviejobs.MovieSearchIndexArgs{
		Operation: moviejobs.SearchIndexOperationReindex,
	}, nil)
	if err != nil {
		h.logger.Error("failed to enqueue reindex job", slog.Any("error",err))
		return nil, fmt.Errorf("failed to enqueue reindex job: %w", err)
	}

	h.logger.Info("reindex job enqueued", slog.Int64("job_id", result.Job.ID))

	return &ogen.ReindexSearchAccepted{
		Message: ogen.NewOptString("Reindex job enqueued"),
		JobID:   ogen.NewOptUUID(uuid.Must(uuid.NewV7())),
	}, nil
}

// Helper functions

func parseUUID(s string) uuid.UUID {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil
	}
	return id
}

func convertFacetValues(values []search.FacetValue) []ogen.FacetValue {
	result := make([]ogen.FacetValue, len(values))
	for i, v := range values {
		result[i] = ogen.FacetValue{
			Value: ogen.NewOptString(v.Value),
			Count: ogen.NewOptInt(v.Count),
		}
	}
	return result
}
