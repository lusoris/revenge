package api

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/service/metadata"
	"go.uber.org/zap"
)

// SearchMoviesMetadata searches TMDb for movies.
func (h *Handler) SearchMoviesMetadata(ctx context.Context, params ogen.SearchMoviesMetadataParams) (ogen.SearchMoviesMetadataRes, error) {
	// Build search options
	opts := metadata.SearchOptions{}
	if params.Year.Set {
		y := int(params.Year.Value)
		opts.Year = &y
	}

	// Get limit
	limit := 20
	if params.Limit.Set {
		limit = int(params.Limit.Value)
	}

	// Search via shared metadata service
	results, err := h.metadataService.SearchMovie(ctx, params.Q, opts)
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

	for i, r := range results {
		if i >= limit {
			break
		}

		result := ogen.MetadataSearchResult{
			Title: ogen.NewOptString(r.Title),
		}

		// Parse TMDb ID from provider ID
		if r.ProviderID != "" {
			var tmdbID int
			fmt.Sscanf(r.ProviderID, "%d", &tmdbID)
			if tmdbID > 0 {
				result.TmdbID = ogen.NewOptInt(tmdbID)
			}
		}

		if r.OriginalTitle != "" {
			result.OriginalTitle = ogen.NewOptString(r.OriginalTitle)
		}
		if r.Overview != "" {
			result.Overview = ogen.NewOptNilString(r.Overview)
		}
		if r.ReleaseDate != nil {
			result.ReleaseDate = ogen.NewOptNilDate(*r.ReleaseDate)
		}
		if r.PosterPath != nil {
			result.PosterPath = ogen.NewOptNilString(*r.PosterPath)
		}
		if r.BackdropPath != nil {
			result.BackdropPath = ogen.NewOptNilString(*r.BackdropPath)
		}
		if r.VoteAverage > 0 {
			result.VoteAverage = ogen.NewOptFloat32(float32(r.VoteAverage))
		}
		if r.VoteCount > 0 {
			result.VoteCount = ogen.NewOptInt(r.VoteCount)
		}
		if r.Popularity > 0 {
			result.Popularity = ogen.NewOptFloat32(float32(r.Popularity))
		}

		response.Results = append(response.Results, result)
	}

	return response, nil
}

// GetMovieMetadata gets detailed movie info from TMDb.
func (h *Handler) GetMovieMetadata(ctx context.Context, params ogen.GetMovieMetadataParams) (ogen.GetMovieMetadataRes, error) {
	// Get movie details from shared metadata service
	movieMeta, err := h.metadataService.GetMovieMetadata(ctx, int32(params.TmdbId), nil)
	if err != nil {
		h.logger.Error("TMDb get movie failed", zap.Error(err))
		return nil, err
	}

	if movieMeta == nil {
		return &ogen.GetMovieMetadataNotFound{}, nil
	}

	// Convert to API response
	response := &ogen.MetadataMovie{
		Title: ogen.NewOptString(movieMeta.Title),
	}

	if movieMeta.TMDbID != nil {
		response.TmdbID = ogen.NewOptInt(int(*movieMeta.TMDbID))
	}
	if movieMeta.IMDbID != nil {
		response.ImdbID = ogen.NewOptNilString(*movieMeta.IMDbID)
	}
	if movieMeta.OriginalTitle != "" {
		response.OriginalTitle = ogen.NewOptString(movieMeta.OriginalTitle)
	}
	if movieMeta.Tagline != nil && *movieMeta.Tagline != "" {
		response.Tagline = ogen.NewOptNilString(*movieMeta.Tagline)
	}
	if movieMeta.Overview != nil && *movieMeta.Overview != "" {
		response.Overview = ogen.NewOptNilString(*movieMeta.Overview)
	}
	if movieMeta.ReleaseDate != nil {
		response.ReleaseDate = ogen.NewOptNilDate(*movieMeta.ReleaseDate)
	}
	if movieMeta.Runtime != nil {
		response.Runtime = ogen.NewOptNilInt(int(*movieMeta.Runtime))
	}
	if movieMeta.Status != "" {
		response.Status = ogen.NewOptString(movieMeta.Status)
	}
	if movieMeta.PosterPath != nil {
		response.PosterPath = ogen.NewOptNilString(*movieMeta.PosterPath)
	}
	if movieMeta.BackdropPath != nil {
		response.BackdropPath = ogen.NewOptNilString(*movieMeta.BackdropPath)
	}
	if movieMeta.VoteAverage > 0 {
		response.VoteAverage = ogen.NewOptFloat32(float32(movieMeta.VoteAverage))
	}
	if movieMeta.VoteCount > 0 {
		response.VoteCount = ogen.NewOptInt(movieMeta.VoteCount)
	}
	if movieMeta.Popularity > 0 {
		response.Popularity = ogen.NewOptFloat32(float32(movieMeta.Popularity))
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
	// Get collection details from shared metadata service
	collection, err := h.metadataService.GetCollectionMetadata(ctx, int32(params.TmdbId), nil)
	if err != nil {
		h.logger.Error("TMDb get collection failed", zap.Error(err))
		return nil, err
	}

	if collection == nil {
		return &ogen.GetCollectionMetadataNotFound{}, nil
	}

	// Convert to API response
	response := &ogen.MetadataCollection{
		Name: ogen.NewOptString(collection.Name),
	}

	// Parse collection ID from provider ID
	if collection.ProviderID != "" {
		var collectionID int
		fmt.Sscanf(collection.ProviderID, "%d", &collectionID)
		if collectionID > 0 {
			response.ID = ogen.NewOptInt(collectionID)
		}
	}
	if collection.Overview != nil && *collection.Overview != "" {
		response.Overview = ogen.NewOptNilString(*collection.Overview)
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
			Title: ogen.NewOptString(part.Title),
		}

		// Parse movie ID from provider ID
		if part.ProviderID != "" {
			var partID int
			fmt.Sscanf(part.ProviderID, "%d", &partID)
			if partID > 0 {
				p.ID = ogen.NewOptInt(partID)
			}
		}
		if part.OriginalTitle != "" {
			p.OriginalTitle = ogen.NewOptString(part.OriginalTitle)
		}
		if part.Overview != "" {
			p.Overview = ogen.NewOptNilString(part.Overview)
		}
		if part.ReleaseDate != nil {
			p.ReleaseDate = ogen.NewOptNilDate(*part.ReleaseDate)
		}
		if part.PosterPath != nil {
			p.PosterPath = ogen.NewOptNilString(*part.PosterPath)
		}
		if part.BackdropPath != nil {
			p.BackdropPath = ogen.NewOptNilString(*part.BackdropPath)
		}
		if part.VoteAverage > 0 {
			p.VoteAverage = ogen.NewOptFloat32(float32(part.VoteAverage))
		}
		if part.VoteCount > 0 {
			p.VoteCount = ogen.NewOptInt(part.VoteCount)
		}
		if part.Popularity > 0 {
			p.Popularity = ogen.NewOptFloat32(float32(part.Popularity))
		}

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
