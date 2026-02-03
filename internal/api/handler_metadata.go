package api

import (
	"bytes"
	"context"
	"time"

	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/content/movie"
	"go.uber.org/zap"
)

// SearchMoviesMetadata searches TMDb for movies.
func (h *Handler) SearchMoviesMetadata(ctx context.Context, params ogen.SearchMoviesMetadataParams) (ogen.SearchMoviesMetadataRes, error) {
	// Get year filter if provided
	var year *int
	if params.Year.Set {
		y := int(params.Year.Value)
		year = &y
	}

	// Get limit
	limit := 20
	if params.Limit.Set {
		limit = int(params.Limit.Value)
	}

	// Search TMDb
	results, err := h.metadataService.SearchMovies(ctx, params.Q, year)
	if err != nil {
		h.logger.Error("TMDb search failed", zap.Error(err))
		return nil, err
	}

	// Convert to API response
	response := &ogen.MetadataSearchResults{
		Page:         ogen.NewOptInt(1),
		TotalResults: ogen.NewOptInt(len(results)),
		TotalPages:   ogen.NewOptInt(1),
		Results:      make([]ogen.MetadataSearchResult, 0, min(len(results), limit)),
	}

	for i, m := range results {
		if i >= limit {
			break
		}

		result := ogen.MetadataSearchResult{
			TmdbID: ogen.NewOptInt(getInt(m.TMDbID)),
			Title:  ogen.NewOptString(m.Title),
		}

		if m.OriginalTitle != nil {
			result.OriginalTitle = ogen.NewOptString(*m.OriginalTitle)
		}
		if m.Overview != nil {
			result.Overview = ogen.NewOptNilString(*m.Overview)
		}
		if m.ReleaseDate != nil {
			result.ReleaseDate = ogen.NewOptNilDate(*m.ReleaseDate)
		}
		if m.PosterPath != nil {
			result.PosterPath = ogen.NewOptNilString(*m.PosterPath)
		}
		if m.BackdropPath != nil {
			result.BackdropPath = ogen.NewOptNilString(*m.BackdropPath)
		}
		if m.VoteAverage != nil {
			if f, ok := m.VoteAverage.Float64(); ok {
				result.VoteAverage = ogen.NewOptFloat32(float32(f))
			}
		}
		if m.VoteCount != nil {
			result.VoteCount = ogen.NewOptInt(int(*m.VoteCount))
		}
		if m.Popularity != nil {
			if f, ok := m.Popularity.Float64(); ok {
				result.Popularity = ogen.NewOptFloat32(float32(f))
			}
		}

		response.Results = append(response.Results, result)
	}

	return response, nil
}

// GetMovieMetadata gets detailed movie info from TMDb.
func (h *Handler) GetMovieMetadata(ctx context.Context, params ogen.GetMovieMetadataParams) (ogen.GetMovieMetadataRes, error) {
	// Get movie details from TMDb
	tmdbMovie, err := h.metadataService.GetMovieByTMDbID(ctx, params.TmdbId)
	if err != nil {
		if err == movie.ErrMovieNotFound {
			return &ogen.GetMovieMetadataNotFound{}, nil
		}
		h.logger.Error("TMDb get movie failed", zap.Error(err))
		return nil, err
	}

	// Convert to API response
	response := &ogen.MetadataMovie{
		TmdbID: ogen.NewOptInt(getInt(tmdbMovie.TMDbID)),
		Title:  ogen.NewOptString(tmdbMovie.Title),
	}

	if tmdbMovie.IMDbID != nil {
		response.ImdbID = ogen.NewOptNilString(*tmdbMovie.IMDbID)
	}
	if tmdbMovie.OriginalTitle != nil {
		response.OriginalTitle = ogen.NewOptString(*tmdbMovie.OriginalTitle)
	}
	if tmdbMovie.Tagline != nil {
		response.Tagline = ogen.NewOptNilString(*tmdbMovie.Tagline)
	}
	if tmdbMovie.Overview != nil {
		response.Overview = ogen.NewOptNilString(*tmdbMovie.Overview)
	}
	if tmdbMovie.ReleaseDate != nil {
		response.ReleaseDate = ogen.NewOptNilDate(*tmdbMovie.ReleaseDate)
	}
	if tmdbMovie.Runtime != nil {
		response.Runtime = ogen.NewOptNilInt(int(*tmdbMovie.Runtime))
	}
	if tmdbMovie.Status != nil {
		response.Status = ogen.NewOptString(*tmdbMovie.Status)
	}
	if tmdbMovie.PosterPath != nil {
		response.PosterPath = ogen.NewOptNilString(*tmdbMovie.PosterPath)
	}
	if tmdbMovie.BackdropPath != nil {
		response.BackdropPath = ogen.NewOptNilString(*tmdbMovie.BackdropPath)
	}
	if tmdbMovie.VoteAverage != nil {
		if f, ok := tmdbMovie.VoteAverage.Float64(); ok {
			response.VoteAverage = ogen.NewOptFloat32(float32(f))
		}
	}
	if tmdbMovie.VoteCount != nil {
		response.VoteCount = ogen.NewOptInt(int(*tmdbMovie.VoteCount))
	}
	if tmdbMovie.Popularity != nil {
		if f, ok := tmdbMovie.Popularity.Float64(); ok {
			response.Popularity = ogen.NewOptFloat32(float32(f))
		}
	}

	return response, nil
}

// GetProxiedImage proxies images from TMDb.
func (h *Handler) GetProxiedImage(ctx context.Context, params ogen.GetProxiedImageParams) (ogen.GetProxiedImageRes, error) {
	// Map ogen type to image service type
	imageType := string(params.Type)
	size := string(params.Size)
	path := "/" + params.Path

	// Fetch image
	data, contentType, err := h.imageService.FetchImage(ctx, imageType, path, size)
	if err != nil {
		h.logger.Error("Image fetch failed", zap.Error(err))
		return &ogen.GetProxiedImageNotFound{}, nil
	}

	// Return based on content type
	switch contentType {
	case "image/png":
		return &ogen.GetProxiedImageOKImagePNG{
			Data: bytes.NewReader(data),
		}, nil
	default:
		// Default to JPEG
		return &ogen.GetProxiedImageOKImageJpeg{
			Data: bytes.NewReader(data),
		}, nil
	}
}

// GetCollectionMetadata gets detailed collection info from TMDb.
func (h *Handler) GetCollectionMetadata(ctx context.Context, params ogen.GetCollectionMetadataParams) (ogen.GetCollectionMetadataRes, error) {
	// Get collection details from TMDb
	collection, err := h.metadataService.GetCollectionDetails(ctx, params.TmdbId)
	if err != nil {
		if err == movie.ErrCollectionNotFound {
			return &ogen.GetCollectionMetadataNotFound{}, nil
		}
		h.logger.Error("TMDb get collection failed", zap.Error(err))
		return nil, err
	}

	// Convert to API response
	response := &ogen.MetadataCollection{
		ID:   ogen.NewOptInt(collection.ID),
		Name: ogen.NewOptString(collection.Name),
	}

	if collection.Overview != "" {
		response.Overview = ogen.NewOptNilString(collection.Overview)
	}
	if collection.PosterPath != nil {
		response.PosterPath = ogen.NewOptNilString(*collection.PosterPath)
	}
	if collection.BackdropPath != nil {
		response.BackdropPath = ogen.NewOptNilString(*collection.BackdropPath)
	}

	// Map parts (movies in collection)
	parts := make([]ogen.MetadataCollectionPart, 0, len(collection.Parts))
	for _, part := range collection.Parts {
		p := ogen.MetadataCollectionPart{
			ID:    ogen.NewOptInt(part.ID),
			Title: ogen.NewOptString(part.Title),
		}

		if part.OriginalTitle != "" {
			p.OriginalTitle = ogen.NewOptString(part.OriginalTitle)
		}
		if part.Overview != "" {
			p.Overview = ogen.NewOptNilString(part.Overview)
		}
		if part.ReleaseDate != "" {
			// Parse date
			date, dateErr := parseDate(part.ReleaseDate)
			if dateErr == nil {
				p.ReleaseDate = ogen.NewOptNilDate(date)
			}
		}
		if part.PosterPath != nil {
			p.PosterPath = ogen.NewOptNilString(*part.PosterPath)
		}
		if part.BackdropPath != nil {
			p.BackdropPath = ogen.NewOptNilString(*part.BackdropPath)
		}
		p.VoteAverage = ogen.NewOptFloat32(float32(part.VoteAverage))
		p.VoteCount = ogen.NewOptInt(part.VoteCount)
		p.Popularity = ogen.NewOptFloat32(float32(part.Popularity))

		parts = append(parts, p)
	}
	response.Parts = parts

	return response, nil
}

// parseDate parses a date string in YYYY-MM-DD format.
func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

// Helper functions
func getInt(i *int32) int {
	if i == nil {
		return 0
	}
	return int(*i)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
