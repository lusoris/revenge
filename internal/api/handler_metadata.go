package api

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/service/metadata"
	"log/slog"
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
		if errors.Is(err, metadata.ErrNoProviders) {
			return &ogen.MetadataSearchResults{Results: []ogen.MetadataSearchResult{}}, nil
		}
		h.logger.Error("TMDb search failed", slog.Any("error", err))
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
			_, _ = fmt.Sscanf(r.ProviderID, "%d", &tmdbID)
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
		if errors.Is(err, metadata.ErrNoProviders) || errors.Is(err, metadata.ErrNotFound) {
			return &ogen.GetMovieMetadataNotFound{}, nil
		}
		h.logger.Error("TMDb get movie failed", slog.Any("error", err))
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
	if h.imageService == nil {
		return &ogen.GetProxiedImageNotFound{}, nil
	}

	// Map ogen type to image service type
	imageType := string(params.Type)
	size := string(params.Size)
	path := "/" + params.Path

	// Fetch image
	data, contentType, err := h.imageService.FetchImage(ctx, imageType, path, size)
	if err != nil {
		h.logger.Error("Image fetch failed", slog.Any("error",err))
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
		if errors.Is(err, metadata.ErrNoProviders) || errors.Is(err, metadata.ErrNotFound) {
			return &ogen.GetCollectionMetadataNotFound{}, nil
		}
		h.logger.Error("TMDb get collection failed", slog.Any("error", err))
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
		_, _ = fmt.Sscanf(collection.ProviderID, "%d", &collectionID)
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
			_, _ = fmt.Sscanf(part.ProviderID, "%d", &partID)
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

// SearchTVShowsMetadata searches TMDb for TV shows.
func (h *Handler) SearchTVShowsMetadata(ctx context.Context, params ogen.SearchTVShowsMetadataParams) (ogen.SearchTVShowsMetadataRes, error) {
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
	results, err := h.metadataService.SearchTVShow(ctx, params.Q, opts)
	if err != nil {
		if errors.Is(err, metadata.ErrNoProviders) {
			return &ogen.MetadataTVSearchResults{Results: []ogen.MetadataTVSearchResult{}}, nil
		}
		h.logger.Error("TMDb TV search failed", slog.Any("error", err))
		return nil, err
	}

	// Convert to API response
	response := &ogen.MetadataTVSearchResults{
		Page:         ogen.NewOptInt(1),
		TotalResults: ogen.NewOptInt(len(results)),
		TotalPages:   ogen.NewOptInt(1),
		Results:      make([]ogen.MetadataTVSearchResult, 0, min(len(results), limit)),
	}

	for i, r := range results {
		if i >= limit {
			break
		}

		result := ogen.MetadataTVSearchResult{
			Name: ogen.NewOptString(r.Name),
		}

		// Parse TMDb ID from provider ID
		if r.ProviderID != "" {
			var tmdbID int
			_, _ = fmt.Sscanf(r.ProviderID, "%d", &tmdbID)
			if tmdbID > 0 {
				result.TmdbID = ogen.NewOptInt(tmdbID)
			}
		}

		if r.OriginalName != "" {
			result.OriginalName = ogen.NewOptString(r.OriginalName)
		}
		if r.Overview != "" {
			result.Overview = ogen.NewOptNilString(r.Overview)
		}
		if r.FirstAirDate != nil {
			result.FirstAirDate = ogen.NewOptNilDate(*r.FirstAirDate)
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

// GetTVShowMetadata gets detailed TV show info from TMDb.
func (h *Handler) GetTVShowMetadata(ctx context.Context, params ogen.GetTVShowMetadataParams) (ogen.GetTVShowMetadataRes, error) {
	// Get TV show details from shared metadata service
	tvMeta, err := h.metadataService.GetTVShowMetadata(ctx, int32(params.TmdbId), nil)
	if err != nil {
		if errors.Is(err, metadata.ErrNoProviders) || errors.Is(err, metadata.ErrNotFound) {
			return &ogen.GetTVShowMetadataNotFound{}, nil
		}
		h.logger.Error("TMDb get TV show failed", slog.Any("error", err))
		return nil, err
	}

	if tvMeta == nil {
		return &ogen.GetTVShowMetadataNotFound{}, nil
	}

	// Convert to API response
	response := &ogen.MetadataTVShow{
		Name: ogen.NewOptString(tvMeta.Name),
	}

	if tvMeta.TMDbID != nil {
		response.TmdbID = ogen.NewOptInt(int(*tvMeta.TMDbID))
	}
	if tvMeta.IMDbID != nil && *tvMeta.IMDbID != "" {
		response.ImdbID = ogen.NewOptNilString(*tvMeta.IMDbID)
	}
	if tvMeta.TVDbID != nil {
		response.TvdbID = ogen.NewOptNilInt(int(*tvMeta.TVDbID))
	}
	if tvMeta.OriginalName != "" {
		response.OriginalName = ogen.NewOptString(tvMeta.OriginalName)
	}
	if tvMeta.Tagline != nil && *tvMeta.Tagline != "" {
		response.Tagline = ogen.NewOptNilString(*tvMeta.Tagline)
	}
	if tvMeta.Overview != nil && *tvMeta.Overview != "" {
		response.Overview = ogen.NewOptNilString(*tvMeta.Overview)
	}
	if tvMeta.FirstAirDate != nil {
		response.FirstAirDate = ogen.NewOptNilDate(*tvMeta.FirstAirDate)
	}
	if tvMeta.LastAirDate != nil {
		response.LastAirDate = ogen.NewOptNilDate(*tvMeta.LastAirDate)
	}
	if tvMeta.Status != "" {
		response.Status = ogen.NewOptString(tvMeta.Status)
	}
	if tvMeta.Type != "" {
		response.Type = ogen.NewOptString(tvMeta.Type)
	}
	if tvMeta.NumberOfSeasons > 0 {
		response.NumberOfSeasons = ogen.NewOptInt(tvMeta.NumberOfSeasons)
	}
	if tvMeta.NumberOfEpisodes > 0 {
		response.NumberOfEpisodes = ogen.NewOptInt(tvMeta.NumberOfEpisodes)
	}
	if len(tvMeta.EpisodeRuntime) > 0 {
		response.EpisodeRunTime = tvMeta.EpisodeRuntime
	}
	if tvMeta.PosterPath != nil {
		response.PosterPath = ogen.NewOptNilString(*tvMeta.PosterPath)
	}
	if tvMeta.BackdropPath != nil {
		response.BackdropPath = ogen.NewOptNilString(*tvMeta.BackdropPath)
	}
	if tvMeta.VoteAverage > 0 {
		response.VoteAverage = ogen.NewOptFloat32(float32(tvMeta.VoteAverage))
	}
	if tvMeta.VoteCount > 0 {
		response.VoteCount = ogen.NewOptInt(tvMeta.VoteCount)
	}
	if tvMeta.Popularity > 0 {
		response.Popularity = ogen.NewOptFloat32(float32(tvMeta.Popularity))
	}

	// Map networks
	if len(tvMeta.Networks) > 0 {
		networks := make([]ogen.MetadataNetwork, 0, len(tvMeta.Networks))
		for _, n := range tvMeta.Networks {
			network := ogen.MetadataNetwork{
				ID:   ogen.NewOptInt(n.ID),
				Name: ogen.NewOptString(n.Name),
			}
			if n.LogoPath != nil && *n.LogoPath != "" {
				network.LogoPath = ogen.NewOptNilString(*n.LogoPath)
			}
			networks = append(networks, network)
		}
		response.Networks = networks
	}

	// Map genres
	if len(tvMeta.Genres) > 0 {
		genres := make([]ogen.MetadataGenre, 0, len(tvMeta.Genres))
		for _, g := range tvMeta.Genres {
			genres = append(genres, ogen.MetadataGenre{
				ID:   ogen.NewOptInt(g.ID),
				Name: ogen.NewOptString(g.Name),
			})
		}
		response.Genres = genres
	}

	// Map seasons
	if len(tvMeta.Seasons) > 0 {
		seasons := make([]ogen.MetadataSeasonSummary, 0, len(tvMeta.Seasons))
		for _, s := range tvMeta.Seasons {
			season := ogen.MetadataSeasonSummary{
				SeasonNumber: ogen.NewOptInt(s.SeasonNumber),
				Name:         ogen.NewOptString(s.Name),
				EpisodeCount: ogen.NewOptInt(s.EpisodeCount),
			}
			// Parse season ID from ProviderID
			if s.ProviderID != "" {
				var seasonID int
				_, _ = fmt.Sscanf(s.ProviderID, "%d", &seasonID)
				if seasonID > 0 {
					season.ID = ogen.NewOptInt(seasonID)
				}
			}
			if s.Overview != nil && *s.Overview != "" {
				season.Overview = ogen.NewOptNilString(*s.Overview)
			}
			if s.AirDate != nil {
				season.AirDate = ogen.NewOptNilDate(*s.AirDate)
			}
			if s.PosterPath != nil && *s.PosterPath != "" {
				season.PosterPath = ogen.NewOptNilString(*s.PosterPath)
			}
			seasons = append(seasons, season)
		}
		response.Seasons = seasons
	}

	return response, nil
}

// GetSeasonMetadata gets detailed season info from TMDb.
func (h *Handler) GetSeasonMetadata(ctx context.Context, params ogen.GetSeasonMetadataParams) (ogen.GetSeasonMetadataRes, error) {
	// Get season details from shared metadata service
	seasonMeta, err := h.metadataService.GetSeasonMetadata(ctx, int32(params.TmdbId), int(params.SeasonNumber), nil)
	if err != nil {
		if errors.Is(err, metadata.ErrNoProviders) || errors.Is(err, metadata.ErrNotFound) {
			return &ogen.GetSeasonMetadataNotFound{}, nil
		}
		h.logger.Error("TMDb get season failed", slog.Any("error", err))
		return nil, err
	}

	if seasonMeta == nil {
		return &ogen.GetSeasonMetadataNotFound{}, nil
	}

	// Convert to API response
	response := &ogen.MetadataSeason{
		TmdbShowID:   ogen.NewOptInt(int(params.TmdbId)),
		SeasonNumber: ogen.NewOptInt(seasonMeta.SeasonNumber),
		Name:         ogen.NewOptString(seasonMeta.Name),
	}

	// Use TMDbID if available
	if seasonMeta.TMDbID != nil {
		response.ID = ogen.NewOptInt(int(*seasonMeta.TMDbID))
	}
	if seasonMeta.Overview != nil && *seasonMeta.Overview != "" {
		response.Overview = ogen.NewOptNilString(*seasonMeta.Overview)
	}
	if seasonMeta.AirDate != nil {
		response.AirDate = ogen.NewOptNilDate(*seasonMeta.AirDate)
	}
	if seasonMeta.PosterPath != nil && *seasonMeta.PosterPath != "" {
		response.PosterPath = ogen.NewOptNilString(*seasonMeta.PosterPath)
	}

	// Map episodes
	if len(seasonMeta.Episodes) > 0 {
		episodes := make([]ogen.MetadataEpisodeSummary, 0, len(seasonMeta.Episodes))
		for _, e := range seasonMeta.Episodes {
			episode := ogen.MetadataEpisodeSummary{
				EpisodeNumber: ogen.NewOptInt(e.EpisodeNumber),
				Name:          ogen.NewOptString(e.Name),
			}
			// Parse episode ID from ProviderID
			if e.ProviderID != "" {
				var episodeID int
				_, _ = fmt.Sscanf(e.ProviderID, "%d", &episodeID)
				if episodeID > 0 {
					episode.ID = ogen.NewOptInt(episodeID)
				}
			}
			if e.Overview != nil && *e.Overview != "" {
				episode.Overview = ogen.NewOptNilString(*e.Overview)
			}
			if e.AirDate != nil {
				episode.AirDate = ogen.NewOptNilDate(*e.AirDate)
			}
			if e.Runtime != nil && *e.Runtime > 0 {
				episode.Runtime = ogen.NewOptNilInt(int(*e.Runtime))
			}
			if e.StillPath != nil && *e.StillPath != "" {
				episode.StillPath = ogen.NewOptNilString(*e.StillPath)
			}
			if e.VoteAverage > 0 {
				episode.VoteAverage = ogen.NewOptFloat32(float32(e.VoteAverage))
			}
			if e.VoteCount > 0 {
				episode.VoteCount = ogen.NewOptInt(e.VoteCount)
			}
			episodes = append(episodes, episode)
		}
		response.Episodes = episodes
	}

	return response, nil
}

// GetEpisodeMetadata gets detailed episode info from TMDb.
func (h *Handler) GetEpisodeMetadata(ctx context.Context, params ogen.GetEpisodeMetadataParams) (ogen.GetEpisodeMetadataRes, error) {
	// Get episode details from shared metadata service
	episodeMeta, err := h.metadataService.GetEpisodeMetadata(ctx, int32(params.TmdbId), int(params.SeasonNumber), int(params.EpisodeNumber), nil)
	if err != nil {
		if errors.Is(err, metadata.ErrNoProviders) || errors.Is(err, metadata.ErrNotFound) {
			return &ogen.GetEpisodeMetadataNotFound{}, nil
		}
		h.logger.Error("TMDb get episode failed", slog.Any("error", err))
		return nil, err
	}

	if episodeMeta == nil {
		return &ogen.GetEpisodeMetadataNotFound{}, nil
	}

	// Convert to API response
	response := &ogen.MetadataEpisode{
		TmdbShowID:    ogen.NewOptInt(int(params.TmdbId)),
		SeasonNumber:  ogen.NewOptInt(int(params.SeasonNumber)),
		EpisodeNumber: ogen.NewOptInt(episodeMeta.EpisodeNumber),
		Name:          ogen.NewOptString(episodeMeta.Name),
	}

	// Use TMDbID if available
	if episodeMeta.TMDbID != nil {
		response.ID = ogen.NewOptInt(int(*episodeMeta.TMDbID))
	}
	if episodeMeta.Overview != nil && *episodeMeta.Overview != "" {
		response.Overview = ogen.NewOptNilString(*episodeMeta.Overview)
	}
	if episodeMeta.AirDate != nil {
		response.AirDate = ogen.NewOptNilDate(*episodeMeta.AirDate)
	}
	if episodeMeta.Runtime != nil && *episodeMeta.Runtime > 0 {
		response.Runtime = ogen.NewOptNilInt(int(*episodeMeta.Runtime))
	}
	if episodeMeta.StillPath != nil && *episodeMeta.StillPath != "" {
		response.StillPath = ogen.NewOptNilString(*episodeMeta.StillPath)
	}
	if episodeMeta.VoteAverage > 0 {
		response.VoteAverage = ogen.NewOptFloat32(float32(episodeMeta.VoteAverage))
	}
	if episodeMeta.VoteCount > 0 {
		response.VoteCount = ogen.NewOptInt(episodeMeta.VoteCount)
	}

	// Map crew
	if len(episodeMeta.Crew) > 0 {
		crew := make([]ogen.MetadataCrewMember, 0, len(episodeMeta.Crew))
		for _, c := range episodeMeta.Crew {
			member := ogen.MetadataCrewMember{
				Name:       ogen.NewOptString(c.Name),
				Job:        ogen.NewOptString(c.Job),
				Department: ogen.NewOptString(c.Department),
			}
			// Parse crew member ID from ProviderID
			if c.ProviderID != "" {
				var crewID int
				_, _ = fmt.Sscanf(c.ProviderID, "%d", &crewID)
				if crewID > 0 {
					member.ID = ogen.NewOptInt(crewID)
				}
			}
			if c.ProfilePath != nil && *c.ProfilePath != "" {
				member.ProfilePath = ogen.NewOptNilString(*c.ProfilePath)
			}
			crew = append(crew, member)
		}
		response.Crew = crew
	}

	// Map guest stars
	if len(episodeMeta.GuestStars) > 0 {
		guestStars := make([]ogen.MetadataCastMember, 0, len(episodeMeta.GuestStars))
		for _, g := range episodeMeta.GuestStars {
			member := ogen.MetadataCastMember{
				Name:      ogen.NewOptString(g.Name),
				Character: ogen.NewOptString(g.Character),
				Order:     ogen.NewOptInt(g.Order),
			}
			// Parse cast member ID from ProviderID
			if g.ProviderID != "" {
				var castID int
				_, _ = fmt.Sscanf(g.ProviderID, "%d", &castID)
				if castID > 0 {
					member.ID = ogen.NewOptInt(castID)
				}
			}
			if g.ProfilePath != nil && *g.ProfilePath != "" {
				member.ProfilePath = ogen.NewOptNilString(*g.ProfilePath)
			}
			guestStars = append(guestStars, member)
		}
		response.GuestStars = guestStars
	}

	return response, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
