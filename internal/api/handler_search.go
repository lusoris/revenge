package api

import (
	"context"
	"fmt"
	"sync"
	"time"

	"log/slog"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/content/movie/moviejobs"
	tvshowjobs "github.com/lusoris/revenge/internal/content/tvshow/jobs"
	"github.com/lusoris/revenge/internal/service/search"
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

// ReindexTVShowSearch triggers a full reindex of all TV shows via River job queue.
func (h *Handler) ReindexTVShowSearch(ctx context.Context) (ogen.ReindexTVShowSearchRes, error) {
	if h.riverClient == nil {
		return nil, fmt.Errorf("job queue not available")
	}

	result, err := h.riverClient.Insert(ctx, tvshowjobs.SearchIndexArgs{
		FullReindex: true,
	}, nil)
	if err != nil {
		h.logger.Error("failed to enqueue TV show reindex job", slog.Any("error", err))
		return nil, fmt.Errorf("failed to enqueue TV show reindex job: %w", err)
	}

	h.logger.Info("TV show reindex job enqueued", slog.Int64("job_id", result.Job.ID))

	return &ogen.ReindexTVShowSearchAccepted{
		Message: ogen.NewOptString("TV show reindex job enqueued"),
		JobID:   ogen.NewOptUUID(uuid.Must(uuid.NewV7())),
	}, nil
}

// SearchLibraryTVShows searches TV shows in the library using Typesense.
func (h *Handler) SearchLibraryTVShows(ctx context.Context, params ogen.SearchLibraryTVShowsParams) (ogen.SearchLibraryTVShowsRes, error) {
	if h.tvshowSearchService == nil {
		return &ogen.TVShowSearchResults{
			TotalHits:    ogen.NewOptInt(0),
			TotalPages:   ogen.NewOptInt(0),
			CurrentPage:  ogen.NewOptInt(1),
			SearchTimeMs: ogen.NewOptInt(0),
			Hits:         []ogen.TVShowSearchHit{},
		}, nil
	}

	searchParams := search.TVShowSearchParams{
		Query:             params.Q,
		Page:              1,
		PerPage:           20,
		SortBy:            "popularity:desc",
		IncludeHighlights: true,
		FacetBy:           []string{"genres", "year", "status", "type", "networks", "has_file"},
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

	result, err := h.tvshowSearchService.SearchSeries(ctx, searchParams)
	if err != nil {
		h.logger.Warn("tvshow search unavailable, returning empty results", slog.Any("error", err))
		return &ogen.TVShowSearchResults{
			TotalHits:    ogen.NewOptInt(0),
			TotalPages:   ogen.NewOptInt(0),
			CurrentPage:  ogen.NewOptInt(searchParams.Page),
			SearchTimeMs: ogen.NewOptInt(0),
			Hits:         []ogen.TVShowSearchHit{},
		}, nil
	}

	response := &ogen.TVShowSearchResults{
		TotalHits:    ogen.NewOptInt(result.TotalHits),
		TotalPages:   ogen.NewOptInt(result.TotalPages),
		CurrentPage:  ogen.NewOptInt(result.CurrentPage),
		SearchTimeMs: ogen.NewOptInt(int(result.SearchTime.Milliseconds())),
		Hits:         make([]ogen.TVShowSearchHit, 0, len(result.Hits)),
	}

	for _, hit := range result.Hits {
		apiHit := ogen.TVShowSearchHit{
			Score: ogen.NewOptFloat32(float32(hit.Score)),
		}

		doc := hit.Document
		apiDoc := ogen.TVShowSearchDocument{
			ID:            ogen.NewOptUUID(parseUUID(doc.ID)),
			TmdbID:        ogen.NewOptInt(int(doc.TMDbID)),
			TvdbID:        ogen.NewOptInt(int(doc.TVDbID)),
			ImdbID:        ogen.NewOptString(doc.IMDbID),
			Title:         ogen.NewOptString(doc.Title),
			OriginalTitle: ogen.NewOptString(doc.OriginalTitle),
			Year:          ogen.NewOptInt(int(doc.Year)),
			Overview:      ogen.NewOptString(doc.Overview),
			Status:        ogen.NewOptString(doc.Status),
			Type:          ogen.NewOptString(doc.Type),
			PosterPath:    ogen.NewOptString(doc.PosterPath),
			BackdropPath:  ogen.NewOptString(doc.BackdropPath),
			VoteAverage:   ogen.NewOptFloat32(float32(doc.VoteAverage)),
			Popularity:    ogen.NewOptFloat32(float32(doc.Popularity)),
			HasFile:       ogen.NewOptBool(doc.HasFile),
			TotalSeasons:  ogen.NewOptInt(int(doc.TotalSeasons)),
			TotalEpisodes: ogen.NewOptInt(int(doc.TotalEpisodes)),
		}

		if doc.FirstAirDate > 0 {
			t := time.Unix(doc.FirstAirDate, 0)
			apiDoc.FirstAirDate = ogen.NewOptDate(t)
		}

		if len(doc.Genres) > 0 {
			apiDoc.Genres = make([]string, len(doc.Genres))
			copy(apiDoc.Genres, doc.Genres)
		}
		if len(doc.Cast) > 0 {
			apiDoc.Cast = make([]string, len(doc.Cast))
			copy(apiDoc.Cast, doc.Cast)
		}
		if len(doc.Networks) > 0 {
			apiDoc.Networks = make([]string, len(doc.Networks))
			copy(apiDoc.Networks, doc.Networks)
		}

		apiHit.Document = ogen.NewOptTVShowSearchDocument(apiDoc)

		if len(hit.Highlights) > 0 {
			highlightMap := make(ogen.TVShowSearchHitHighlights)
			for field, snippets := range hit.Highlights {
				highlightMap[field] = snippets
			}
			apiHit.Highlights = ogen.OptTVShowSearchHitHighlights{
				Value: highlightMap,
				Set:   true,
			}
		}

		response.Hits = append(response.Hits, apiHit)
	}

	if len(result.Facets) > 0 {
		facets := ogen.OptTVShowSearchResultsFacets{
			Set:   true,
			Value: make(ogen.TVShowSearchResultsFacets),
		}
		for facetName, values := range result.Facets {
			facets.Value[facetName] = convertFacetValues(values)
		}
		response.Facets = facets
	}

	return response, nil
}

// AutocompleteTVShows provides autocomplete suggestions for TV show titles.
func (h *Handler) AutocompleteTVShows(ctx context.Context, params ogen.AutocompleteTVShowsParams) (ogen.AutocompleteTVShowsRes, error) {
	if h.tvshowSearchService == nil {
		return &ogen.AutocompleteResults{Suggestions: []string{}}, nil
	}

	limit := 5
	if params.Limit.Set {
		limit = params.Limit.Value
	}

	suggestions, err := h.tvshowSearchService.AutocompleteSeries(ctx, params.Q, limit)
	if err != nil {
		h.logger.Warn("tvshow autocomplete unavailable, returning empty", slog.Any("error", err))
		return &ogen.AutocompleteResults{Suggestions: []string{}}, nil
	}

	return &ogen.AutocompleteResults{Suggestions: suggestions}, nil
}

// GetTVShowSearchFacets returns available facet values for TV show filtering.
func (h *Handler) GetTVShowSearchFacets(ctx context.Context) (ogen.GetTVShowSearchFacetsRes, error) {
	if h.tvshowSearchService == nil {
		return &ogen.TVShowSearchFacets{}, nil
	}

	facetNames := []string{"genres", "year", "status", "type", "networks", "has_file"}

	facets, err := h.tvshowSearchService.GetFacets(ctx, facetNames)
	if err != nil {
		h.logger.Warn("tvshow facets unavailable, returning empty", slog.Any("error", err))
		return &ogen.TVShowSearchFacets{}, nil
	}

	response := &ogen.TVShowSearchFacets{}

	if values, ok := facets["genres"]; ok {
		response.Genres = convertFacetValues(values)
	}
	if values, ok := facets["year"]; ok {
		response.Years = convertFacetValues(values)
	}
	if values, ok := facets["status"]; ok {
		response.Status = convertFacetValues(values)
	}
	if values, ok := facets["type"]; ok {
		response.Type = convertFacetValues(values)
	}
	if values, ok := facets["networks"]; ok {
		response.Networks = convertFacetValues(values)
	}
	if values, ok := facets["has_file"]; ok {
		response.HasFile = convertFacetValues(values)
	}

	return response, nil
}

// Helper functions

// SearchMulti searches across all collections in parallel and returns merged results.
func (h *Handler) SearchMulti(ctx context.Context, params ogen.SearchMultiParams) (ogen.SearchMultiRes, error) {
	limit := 5
	if params.Limit.Set {
		limit = params.Limit.Value
	}

	response := &ogen.MultiSearchResults{}
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Search movies
	if h.searchService != nil && h.searchService.IsEnabled() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result, err := h.searchService.Search(ctx, search.SearchParams{
				Query:   params.Q,
				Page:    1,
				PerPage: limit,
				SortBy:  "_text_match:desc",
			})
			if err != nil {
				h.logger.Warn("multi-search: movies unavailable", slog.Any("error", err))
				return
			}
			mu.Lock()
			response.Movies = ogen.NewOptSearchResults(h.convertMovieResults(result))
			mu.Unlock()
		}()
	}

	// Search TV shows
	if h.tvshowSearchService != nil && h.tvshowSearchService.IsEnabled() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result, err := h.tvshowSearchService.SearchSeries(ctx, search.TVShowSearchParams{
				Query:   params.Q,
				Page:    1,
				PerPage: limit,
				SortBy:  "_text_match:desc",
			})
			if err != nil {
				h.logger.Warn("multi-search: tvshows unavailable", slog.Any("error", err))
				return
			}
			mu.Lock()
			response.Tvshows = ogen.NewOptTVShowSearchResults(h.convertTVShowResults(result))
			mu.Unlock()
		}()
	}

	// Search episodes
	if h.episodeSearchService != nil && h.episodeSearchService.IsEnabled() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result, err := h.episodeSearchService.SearchEpisodes(ctx, search.EpisodeSearchParams{
				Query:   params.Q,
				Page:    1,
				PerPage: limit,
				SortBy:  "_text_match:desc",
			})
			if err != nil {
				h.logger.Warn("multi-search: episodes unavailable", slog.Any("error", err))
				return
			}
			mu.Lock()
			response.Episodes = ogen.NewOptEpisodeSearchResults(h.convertEpisodeResults(result))
			mu.Unlock()
		}()
	}

	// Search seasons
	if h.seasonSearchService != nil && h.seasonSearchService.IsEnabled() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result, err := h.seasonSearchService.SearchSeasons(ctx, search.SeasonSearchParams{
				Query:   params.Q,
				Page:    1,
				PerPage: limit,
				SortBy:  "_text_match:desc",
			})
			if err != nil {
				h.logger.Warn("multi-search: seasons unavailable", slog.Any("error", err))
				return
			}
			mu.Lock()
			response.Seasons = ogen.NewOptSeasonSearchResults(h.convertSeasonResults(result))
			mu.Unlock()
		}()
	}

	// Search people
	if h.personSearchService != nil && h.personSearchService.IsEnabled() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result, err := h.personSearchService.SearchPersons(ctx, search.PersonSearchParams{
				Query:   params.Q,
				Page:    1,
				PerPage: limit,
				SortBy:  "_text_match:desc",
			})
			if err != nil {
				h.logger.Warn("multi-search: people unavailable", slog.Any("error", err))
				return
			}
			mu.Lock()
			response.People = ogen.NewOptPersonSearchResults(h.convertPersonResults(result))
			mu.Unlock()
		}()
	}

	wg.Wait()

	return response, nil
}

// convertMovieResults converts movie search results to API response type.
func (h *Handler) convertMovieResults(result *search.SearchResult) ogen.SearchResults {
	if result == nil {
		return ogen.SearchResults{Hits: []ogen.SearchHit{}}
	}

	resp := ogen.SearchResults{
		TotalHits:    ogen.NewOptInt(result.TotalHits),
		SearchTimeMs: ogen.NewOptInt(int(result.SearchTime.Milliseconds())),
		Hits:         make([]ogen.SearchHit, 0, len(result.Hits)),
	}

	for _, hit := range result.Hits {
		doc := hit.Document
		apiDoc := ogen.SearchDocument{
			ID:            ogen.NewOptUUID(parseUUID(doc.ID)),
			TmdbID:        ogen.NewOptInt(int(doc.TMDbID)),
			Title:         ogen.NewOptString(doc.Title),
			Overview:      ogen.NewOptString(doc.Overview),
			PosterPath:    ogen.NewOptString(doc.PosterPath),
			VoteAverage:   ogen.NewOptFloat32(float32(doc.VoteAverage)),
			HasFile:       ogen.NewOptBool(doc.HasFile),
			Year:          ogen.NewOptInt(int(doc.Year)),
		}
		if doc.ReleaseDate > 0 {
			apiDoc.ReleaseDate = ogen.NewOptDate(time.Unix(doc.ReleaseDate, 0))
		}
		resp.Hits = append(resp.Hits, ogen.SearchHit{
			Document: ogen.NewOptSearchDocument(apiDoc),
			Score:    ogen.NewOptFloat32(float32(hit.Score)),
		})
	}

	return resp
}

// convertTVShowResults converts TV show search results to API response type.
func (h *Handler) convertTVShowResults(result *search.TVShowSearchResult) ogen.TVShowSearchResults {
	if result == nil {
		return ogen.TVShowSearchResults{Hits: []ogen.TVShowSearchHit{}}
	}

	resp := ogen.TVShowSearchResults{
		TotalHits:    ogen.NewOptInt(result.TotalHits),
		SearchTimeMs: ogen.NewOptInt(int(result.SearchTime.Milliseconds())),
		Hits:         make([]ogen.TVShowSearchHit, 0, len(result.Hits)),
	}

	for _, hit := range result.Hits {
		doc := hit.Document
		apiDoc := ogen.TVShowSearchDocument{
			ID:           ogen.NewOptUUID(parseUUID(doc.ID)),
			TmdbID:       ogen.NewOptInt(int(doc.TMDbID)),
			Title:        ogen.NewOptString(doc.Title),
			Overview:     ogen.NewOptString(doc.Overview),
			PosterPath:   ogen.NewOptString(doc.PosterPath),
			VoteAverage:  ogen.NewOptFloat32(float32(doc.VoteAverage)),
			HasFile:      ogen.NewOptBool(doc.HasFile),
			Year:         ogen.NewOptInt(int(doc.Year)),
		}
		if doc.FirstAirDate > 0 {
			apiDoc.FirstAirDate = ogen.NewOptDate(time.Unix(doc.FirstAirDate, 0))
		}
		resp.Hits = append(resp.Hits, ogen.TVShowSearchHit{
			Document: ogen.NewOptTVShowSearchDocument(apiDoc),
			Score:    ogen.NewOptFloat32(float32(hit.Score)),
		})
	}

	return resp
}

// convertEpisodeResults converts episode search results to API response type.
func (h *Handler) convertEpisodeResults(result *search.EpisodeSearchResult) ogen.EpisodeSearchResults {
	if result == nil {
		return ogen.EpisodeSearchResults{Hits: []ogen.EpisodeSearchHit{}}
	}

	resp := ogen.EpisodeSearchResults{
		TotalHits: ogen.NewOptInt(result.TotalHits),
		Hits:      make([]ogen.EpisodeSearchHit, 0, len(result.Hits)),
	}

	for _, hit := range result.Hits {
		doc := hit.Document
		apiDoc := ogen.EpisodeSearchDocument{
			ID:               ogen.NewOptUUID(parseUUID(doc.ID)),
			SeriesID:         ogen.NewOptUUID(parseUUID(doc.SeriesID)),
			SeasonNumber:     ogen.NewOptInt(int(doc.SeasonNumber)),
			EpisodeNumber:    ogen.NewOptInt(int(doc.EpisodeNumber)),
			Title:            ogen.NewOptString(doc.Title),
			Overview:         ogen.NewOptString(doc.Overview),
			AirDate:          ogen.NewOptInt(int(doc.AirDate)),
			Runtime:          ogen.NewOptInt(int(doc.Runtime)),
			VoteAverage:      ogen.NewOptFloat32(float32(doc.VoteAverage)),
			StillPath:        ogen.NewOptString(doc.StillPath),
			HasFile:          ogen.NewOptBool(doc.HasFile),
			SeriesTitle:      ogen.NewOptString(doc.SeriesTitle),
			SeriesPosterPath: ogen.NewOptString(doc.SeriesPosterPath),
		}
		resp.Hits = append(resp.Hits, ogen.EpisodeSearchHit{
			Document: ogen.NewOptEpisodeSearchDocument(apiDoc),
			Score:    ogen.NewOptFloat32(float32(hit.Score)),
		})
	}

	return resp
}

// convertSeasonResults converts season search results to API response type.
func (h *Handler) convertSeasonResults(result *search.SeasonSearchResult) ogen.SeasonSearchResults {
	if result == nil {
		return ogen.SeasonSearchResults{Hits: []ogen.SeasonSearchHit{}}
	}

	resp := ogen.SeasonSearchResults{
		TotalHits: ogen.NewOptInt(result.TotalHits),
		Hits:      make([]ogen.SeasonSearchHit, 0, len(result.Hits)),
	}

	for _, hit := range result.Hits {
		doc := hit.Document
		apiDoc := ogen.SeasonSearchDocument{
			ID:               ogen.NewOptUUID(parseUUID(doc.ID)),
			SeriesID:         ogen.NewOptUUID(parseUUID(doc.SeriesID)),
			SeasonNumber:     ogen.NewOptInt(int(doc.SeasonNumber)),
			Name:             ogen.NewOptString(doc.Name),
			Overview:         ogen.NewOptString(doc.Overview),
			AirDate:          ogen.NewOptInt(int(doc.AirDate)),
			EpisodeCount:     ogen.NewOptInt(int(doc.EpisodeCount)),
			VoteAverage:      ogen.NewOptFloat32(float32(doc.VoteAverage)),
			PosterPath:       ogen.NewOptString(doc.PosterPath),
			SeriesTitle:      ogen.NewOptString(doc.SeriesTitle),
			SeriesPosterPath: ogen.NewOptString(doc.SeriesPosterPath),
		}
		resp.Hits = append(resp.Hits, ogen.SeasonSearchHit{
			Document: ogen.NewOptSeasonSearchDocument(apiDoc),
			Score:    ogen.NewOptFloat32(float32(hit.Score)),
		})
	}

	return resp
}

// convertPersonResults converts person search results to API response type.
func (h *Handler) convertPersonResults(result *search.PersonSearchResult) ogen.PersonSearchResults {
	if result == nil {
		return ogen.PersonSearchResults{Hits: []ogen.PersonSearchHit{}}
	}

	resp := ogen.PersonSearchResults{
		TotalHits: ogen.NewOptInt(result.TotalHits),
		Hits:      make([]ogen.PersonSearchHit, 0, len(result.Hits)),
	}

	for _, hit := range result.Hits {
		doc := hit.Document
		apiDoc := ogen.PersonSearchDocument{
			ID:           ogen.NewOptString(doc.ID),
			TmdbID:       ogen.NewOptInt(int(doc.TMDbID)),
			Name:         ogen.NewOptString(doc.Name),
			ProfilePath:  ogen.NewOptString(doc.ProfilePath),
			MovieCount:   ogen.NewOptInt(int(doc.MovieCount)),
			TvshowCount:  ogen.NewOptInt(int(doc.TVShowCount)),
			TotalCredits: ogen.NewOptInt(int(doc.TotalCredits)),
		}
		if len(doc.KnownFor) > 0 {
			apiDoc.KnownFor = make([]string, len(doc.KnownFor))
			copy(apiDoc.KnownFor, doc.KnownFor)
		}
		if len(doc.Characters) > 0 {
			apiDoc.Characters = make([]string, len(doc.Characters))
			copy(apiDoc.Characters, doc.Characters)
		}
		if len(doc.Departments) > 0 {
			apiDoc.Departments = make([]string, len(doc.Departments))
			copy(apiDoc.Departments, doc.Departments)
		}
		resp.Hits = append(resp.Hits, ogen.PersonSearchHit{
			Document: ogen.NewOptPersonSearchDocument(apiDoc),
			Score:    ogen.NewOptFloat32(float32(hit.Score)),
		})
	}

	return resp
}

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
