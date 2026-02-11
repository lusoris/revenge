package api

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"log/slog"

	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/service/metadata"
	"github.com/lusoris/revenge/internal/util"
)

// SearchMoviesMetadata searches for movies via metadata provider.
func (h *Handler) SearchMoviesMetadata(ctx context.Context, params ogen.SearchMoviesMetadataParams) (ogen.SearchMoviesMetadataRes, error) {
	// Handle Radarr provider directly via the *arr integration
	if params.Provider.Set && params.Provider.Value == ogen.SearchMoviesMetadataProviderRadarr {
		return h.searchMoviesViaRadarr(ctx, params)
	}

	// Build search options
	opts := metadata.SearchOptions{}
	if params.Year.Set {
		y := int(params.Year.Value)
		opts.Year = &y
	}
	if params.Provider.Set {
		opts.ProviderID = metadata.ProviderID(params.Provider.Value)
	}
	if params.Language.Set {
		opts.Language = params.Language.Value
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
		h.logger.Error("metadata movie search failed", slog.Any("error", err))
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
	movieMeta, err := h.metadataService.GetMovieMetadata(ctx, util.SafeIntToInt32(params.TmdbId), nil)
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
	collection, err := h.metadataService.GetCollectionMetadata(ctx, util.SafeIntToInt32(params.TmdbId), nil)
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

// SearchTVShowsMetadata searches for TV shows via metadata provider.
func (h *Handler) SearchTVShowsMetadata(ctx context.Context, params ogen.SearchTVShowsMetadataParams) (ogen.SearchTVShowsMetadataRes, error) {
	// Handle Sonarr provider directly via the *arr integration
	if params.Provider.Set && params.Provider.Value == ogen.SearchTVShowsMetadataProviderSonarr {
		return h.searchTVShowsViaSonarr(ctx, params)
	}

	// Build search options
	opts := metadata.SearchOptions{}
	if params.Year.Set {
		y := int(params.Year.Value)
		opts.Year = &y
	}
	if params.Provider.Set {
		opts.ProviderID = metadata.ProviderID(params.Provider.Value)
	}
	if params.Language.Set {
		opts.Language = params.Language.Value
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
		h.logger.Error("metadata TV search failed", slog.Any("error", err))
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
	tvMeta, err := h.metadataService.GetTVShowMetadata(ctx, util.SafeIntToInt32(params.TmdbId), nil)
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
	seasonMeta, err := h.metadataService.GetSeasonMetadata(ctx, util.SafeIntToInt32(params.TmdbId), int(params.SeasonNumber), nil)
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
	episodeMeta, err := h.metadataService.GetEpisodeMetadata(ctx, util.SafeIntToInt32(params.TmdbId), int(params.SeasonNumber), int(params.EpisodeNumber), nil)
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

// GetMovieMetadataCredits gets movie credits from TMDb.
func (h *Handler) GetMovieMetadataCredits(ctx context.Context, params ogen.GetMovieMetadataCreditsParams) (ogen.GetMovieMetadataCreditsRes, error) {
	// Get movie credits from shared metadata service
	credits, err := h.metadataService.GetMovieCredits(ctx, util.SafeIntToInt32(params.TmdbId))
	if err != nil {
		if errors.Is(err, metadata.ErrNoProviders) || errors.Is(err, metadata.ErrNotFound) {
			return &ogen.GetMovieMetadataCreditsNotFound{}, nil
		}
		h.logger.Error("TMDb get movie credits failed", slog.Any("error", err))
		return nil, err
	}

	if credits == nil {
		return &ogen.GetMovieMetadataCreditsNotFound{}, nil
	}

	// Convert to API response
	response := &ogen.MetadataCredits{}

	// Map cast
	if len(credits.Cast) > 0 {
		cast := make([]ogen.MetadataCastMember, 0, len(credits.Cast))
		for _, c := range credits.Cast {
			member := ogen.MetadataCastMember{
				Name:      ogen.NewOptString(c.Name),
				Character: ogen.NewOptString(c.Character),
				Order:     ogen.NewOptInt(c.Order),
			}
			// Parse cast member ID from ProviderID
			if c.ProviderID != "" {
				var castID int
				_, _ = fmt.Sscanf(c.ProviderID, "%d", &castID)
				if castID > 0 {
					member.ID = ogen.NewOptInt(castID)
				}
			}
			if c.ProfilePath != nil && *c.ProfilePath != "" {
				member.ProfilePath = ogen.NewOptNilString(*c.ProfilePath)
			}
			cast = append(cast, member)
		}
		response.Cast = cast
	}

	// Map crew
	if len(credits.Crew) > 0 {
		crew := make([]ogen.MetadataCrewMember, 0, len(credits.Crew))
		for _, c := range credits.Crew {
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

	return response, nil
}

// GetMovieMetadataImages gets movie images from metadata provider.
func (h *Handler) GetMovieMetadataImages(ctx context.Context, params ogen.GetMovieMetadataImagesParams) (ogen.GetMovieMetadataImagesRes, error) {
	images, err := h.metadataService.GetMovieImages(ctx, util.SafeIntToInt32(params.TmdbId))
	if err != nil {
		if errors.Is(err, metadata.ErrNoProviders) || errors.Is(err, metadata.ErrNotFound) {
			return &ogen.GetMovieMetadataImagesNotFound{}, nil
		}
		h.logger.Error("get movie images failed", slog.Any("error", err))
		return nil, err
	}

	if images == nil {
		return &ogen.GetMovieMetadataImagesNotFound{}, nil
	}

	return convertImages(images), nil
}

// GetSimilarMoviesMetadata gets similar movies from metadata provider.
func (h *Handler) GetSimilarMoviesMetadata(ctx context.Context, params ogen.GetSimilarMoviesMetadataParams) (ogen.GetSimilarMoviesMetadataRes, error) {
	limit := 20
	if params.Limit.Set {
		limit = int(params.Limit.Value)
	}

	results, total, err := h.metadataService.GetSimilarMovies(ctx, util.SafeIntToInt32(params.TmdbId), metadata.SearchOptions{})
	if err != nil {
		if errors.Is(err, metadata.ErrNoProviders) || errors.Is(err, metadata.ErrNotFound) {
			return &ogen.GetSimilarMoviesMetadataNotFound{}, nil
		}
		h.logger.Error("get similar movies failed", slog.Any("error", err))
		return nil, err
	}

	return convertMovieSearchResults(results, limit, total), nil
}

// GetMovieRecommendationsMetadata gets movie recommendations from metadata provider.
func (h *Handler) GetMovieRecommendationsMetadata(ctx context.Context, params ogen.GetMovieRecommendationsMetadataParams) (ogen.GetMovieRecommendationsMetadataRes, error) {
	limit := 20
	if params.Limit.Set {
		limit = int(params.Limit.Value)
	}

	results, total, err := h.metadataService.GetMovieRecommendations(ctx, util.SafeIntToInt32(params.TmdbId), metadata.SearchOptions{})
	if err != nil {
		if errors.Is(err, metadata.ErrNoProviders) || errors.Is(err, metadata.ErrNotFound) {
			return &ogen.GetMovieRecommendationsMetadataNotFound{}, nil
		}
		h.logger.Error("get movie recommendations failed", slog.Any("error", err))
		return nil, err
	}

	return convertMovieSearchResults(results, limit, total), nil
}

// GetMovieExternalIDs gets movie external IDs.
func (h *Handler) GetMovieExternalIDs(ctx context.Context, params ogen.GetMovieExternalIDsParams) (ogen.GetMovieExternalIDsRes, error) {
	ids, err := h.metadataService.GetMovieExternalIDs(ctx, util.SafeIntToInt32(params.TmdbId))
	if err != nil {
		if errors.Is(err, metadata.ErrNoProviders) || errors.Is(err, metadata.ErrNotFound) {
			return &ogen.GetMovieExternalIDsNotFound{}, nil
		}
		h.logger.Error("get movie external IDs failed", slog.Any("error", err))
		return nil, err
	}

	if ids == nil {
		return &ogen.GetMovieExternalIDsNotFound{}, nil
	}

	return convertExternalIDs(ids), nil
}

// GetTVShowMetadataCredits gets TV show credits from metadata provider.
func (h *Handler) GetTVShowMetadataCredits(ctx context.Context, params ogen.GetTVShowMetadataCreditsParams) (ogen.GetTVShowMetadataCreditsRes, error) {
	credits, err := h.metadataService.GetTVShowCredits(ctx, util.SafeIntToInt32(params.TmdbId))
	if err != nil {
		if errors.Is(err, metadata.ErrNoProviders) || errors.Is(err, metadata.ErrNotFound) {
			return &ogen.GetTVShowMetadataCreditsNotFound{}, nil
		}
		h.logger.Error("get TV show credits failed", slog.Any("error", err))
		return nil, err
	}

	if credits == nil {
		return &ogen.GetTVShowMetadataCreditsNotFound{}, nil
	}

	return convertCredits(credits), nil
}

// GetTVShowMetadataImages gets TV show images from metadata provider.
func (h *Handler) GetTVShowMetadataImages(ctx context.Context, params ogen.GetTVShowMetadataImagesParams) (ogen.GetTVShowMetadataImagesRes, error) {
	images, err := h.metadataService.GetTVShowImages(ctx, util.SafeIntToInt32(params.TmdbId))
	if err != nil {
		if errors.Is(err, metadata.ErrNoProviders) || errors.Is(err, metadata.ErrNotFound) {
			return &ogen.GetTVShowMetadataImagesNotFound{}, nil
		}
		h.logger.Error("get TV show images failed", slog.Any("error", err))
		return nil, err
	}

	if images == nil {
		return &ogen.GetTVShowMetadataImagesNotFound{}, nil
	}

	return convertImages(images), nil
}

// GetTVShowContentRatings gets TV show content ratings.
func (h *Handler) GetTVShowContentRatings(ctx context.Context, params ogen.GetTVShowContentRatingsParams) (ogen.GetTVShowContentRatingsRes, error) {
	ratings, err := h.metadataService.GetTVShowContentRatings(ctx, util.SafeIntToInt32(params.TmdbId))
	if err != nil {
		if errors.Is(err, metadata.ErrNoProviders) || errors.Is(err, metadata.ErrNotFound) {
			return &ogen.GetTVShowContentRatingsNotFound{}, nil
		}
		h.logger.Error("get TV show content ratings failed", slog.Any("error", err))
		return nil, err
	}

	results := make([]ogen.MetadataContentRating, 0, len(ratings))
	for _, r := range ratings {
		rating := ogen.MetadataContentRating{
			CountryCode: ogen.NewOptString(r.CountryCode),
			Rating:      ogen.NewOptString(r.Rating),
			Descriptors: r.Descriptors,
		}
		results = append(results, rating)
	}

	return &ogen.MetadataContentRatings{Results: results}, nil
}

// GetTVShowExternalIDs gets TV show external IDs.
func (h *Handler) GetTVShowExternalIDs(ctx context.Context, params ogen.GetTVShowExternalIDsParams) (ogen.GetTVShowExternalIDsRes, error) {
	ids, err := h.metadataService.GetTVShowExternalIDs(ctx, util.SafeIntToInt32(params.TmdbId))
	if err != nil {
		if errors.Is(err, metadata.ErrNoProviders) || errors.Is(err, metadata.ErrNotFound) {
			return &ogen.GetTVShowExternalIDsNotFound{}, nil
		}
		h.logger.Error("get TV show external IDs failed", slog.Any("error", err))
		return nil, err
	}

	if ids == nil {
		return &ogen.GetTVShowExternalIDsNotFound{}, nil
	}

	return convertExternalIDs(ids), nil
}

// GetSeasonMetadataImages gets season images from metadata provider.
func (h *Handler) GetSeasonMetadataImages(ctx context.Context, params ogen.GetSeasonMetadataImagesParams) (ogen.GetSeasonMetadataImagesRes, error) {
	images, err := h.metadataService.GetSeasonImages(ctx, util.SafeIntToInt32(params.TmdbId), int(params.SeasonNumber))
	if err != nil {
		if errors.Is(err, metadata.ErrNoProviders) || errors.Is(err, metadata.ErrNotFound) {
			return &ogen.GetSeasonMetadataImagesNotFound{}, nil
		}
		h.logger.Error("get season images failed", slog.Any("error", err))
		return nil, err
	}

	if images == nil {
		return &ogen.GetSeasonMetadataImagesNotFound{}, nil
	}

	return convertImages(images), nil
}

// GetEpisodeMetadataImages gets episode images from metadata provider.
func (h *Handler) GetEpisodeMetadataImages(ctx context.Context, params ogen.GetEpisodeMetadataImagesParams) (ogen.GetEpisodeMetadataImagesRes, error) {
	images, err := h.metadataService.GetEpisodeImages(ctx, util.SafeIntToInt32(params.TmdbId), int(params.SeasonNumber), int(params.EpisodeNumber))
	if err != nil {
		if errors.Is(err, metadata.ErrNoProviders) || errors.Is(err, metadata.ErrNotFound) {
			return &ogen.GetEpisodeMetadataImagesNotFound{}, nil
		}
		h.logger.Error("get episode images failed", slog.Any("error", err))
		return nil, err
	}

	if images == nil {
		return &ogen.GetEpisodeMetadataImagesNotFound{}, nil
	}

	return convertImages(images), nil
}

// SearchPersonMetadata searches for people via metadata provider.
func (h *Handler) SearchPersonMetadata(ctx context.Context, params ogen.SearchPersonMetadataParams) (ogen.SearchPersonMetadataRes, error) {
	opts := metadata.SearchOptions{}
	if params.Provider.Set {
		opts.ProviderID = metadata.ProviderID(params.Provider.Value)
	}
	if params.Language.Set {
		opts.Language = params.Language.Value
	}

	limit := 20
	if params.Limit.Set {
		limit = int(params.Limit.Value)
	}

	results, err := h.metadataService.SearchPerson(ctx, params.Q, opts)
	if err != nil {
		if errors.Is(err, metadata.ErrNoProviders) {
			return &ogen.MetadataPersonSearchResults{Results: []ogen.MetadataPersonSearchResult{}}, nil
		}
		h.logger.Error("person search failed", slog.Any("error", err))
		return nil, err
	}

	response := &ogen.MetadataPersonSearchResults{
		Page:         ogen.NewOptInt(1),
		TotalResults: ogen.NewOptInt(len(results)),
		TotalPages:   ogen.NewOptInt(1),
		Results:      make([]ogen.MetadataPersonSearchResult, 0, min(len(results), limit)),
	}

	for i, r := range results {
		if i >= limit {
			break
		}

		result := ogen.MetadataPersonSearchResult{
			Name: ogen.NewOptString(r.Name),
		}

		if r.ProviderID != "" {
			var tmdbID int
			_, _ = fmt.Sscanf(r.ProviderID, "%d", &tmdbID)
			if tmdbID > 0 {
				result.TmdbID = ogen.NewOptInt(tmdbID)
			}
		}
		if r.ProfilePath != nil {
			result.ProfilePath = ogen.NewOptNilString(*r.ProfilePath)
		}
		if r.Popularity > 0 {
			result.Popularity = ogen.NewOptFloat32(float32(r.Popularity))
		}
		if r.KnownFor != nil {
			knownFor := make([]ogen.MetadataMediaReference, 0, len(r.KnownFor))
			for _, kf := range r.KnownFor {
				ref := ogen.MetadataMediaReference{
					Title: ogen.NewOptString(kf.Title),
				}
				switch kf.MediaType {
				case "movie":
					ref.MediaType = ogen.NewOptMetadataMediaReferenceMediaType(ogen.MetadataMediaReferenceMediaTypeMovie)
				case "tv":
					ref.MediaType = ogen.NewOptMetadataMediaReferenceMediaType(ogen.MetadataMediaReferenceMediaTypeTv)
				}
				if kf.ID != "" {
					var id int
					_, _ = fmt.Sscanf(kf.ID, "%d", &id)
					if id > 0 {
						ref.TmdbID = ogen.NewOptInt(id)
					}
				}
				if kf.PosterPath != nil {
					ref.PosterPath = ogen.NewOptNilString(*kf.PosterPath)
				}
				knownFor = append(knownFor, ref)
			}
			result.KnownFor = knownFor
		}

		response.Results = append(response.Results, result)
	}

	return response, nil
}

// GetPersonMetadata gets person details from metadata provider.
func (h *Handler) GetPersonMetadata(ctx context.Context, params ogen.GetPersonMetadataParams) (ogen.GetPersonMetadataRes, error) {
	person, err := h.metadataService.GetPersonMetadata(ctx, util.SafeIntToInt32(params.TmdbId), nil)
	if err != nil {
		if errors.Is(err, metadata.ErrNoProviders) || errors.Is(err, metadata.ErrNotFound) {
			return &ogen.GetPersonMetadataNotFound{}, nil
		}
		h.logger.Error("get person metadata failed", slog.Any("error", err))
		return nil, err
	}

	if person == nil {
		return &ogen.GetPersonMetadataNotFound{}, nil
	}

	response := &ogen.MetadataPerson{
		Name: ogen.NewOptString(person.Name),
	}

	if person.TMDbID != nil {
		response.TmdbID = ogen.NewOptInt(int(*person.TMDbID))
	}
	if person.IMDbID != nil {
		response.ImdbID = ogen.NewOptNilString(*person.IMDbID)
	}
	if len(person.AlsoKnownAs) > 0 {
		response.AlsoKnownAs = person.AlsoKnownAs
	}
	if person.Biography != nil {
		response.Biography = ogen.NewOptNilString(*person.Biography)
	}
	if person.Birthday != nil {
		response.Birthday = ogen.NewOptNilDate(*person.Birthday)
	}
	if person.Deathday != nil {
		response.Deathday = ogen.NewOptNilDate(*person.Deathday)
	}
	response.Gender = ogen.NewOptInt(person.Gender)
	if person.PlaceOfBirth != nil {
		response.PlaceOfBirth = ogen.NewOptNilString(*person.PlaceOfBirth)
	}
	if person.ProfilePath != nil {
		response.ProfilePath = ogen.NewOptNilString(*person.ProfilePath)
	}
	if person.Homepage != nil {
		response.Homepage = ogen.NewOptNilString(*person.Homepage)
	}
	if person.Popularity > 0 {
		response.Popularity = ogen.NewOptFloat32(float32(person.Popularity))
	}
	if person.KnownForDept != "" {
		response.KnownForDepartment = ogen.NewOptString(person.KnownForDept)
	}

	return response, nil
}

// GetPersonMetadataCredits gets person credits from metadata provider.
func (h *Handler) GetPersonMetadataCredits(ctx context.Context, params ogen.GetPersonMetadataCreditsParams) (ogen.GetPersonMetadataCreditsRes, error) {
	credits, err := h.metadataService.GetPersonCredits(ctx, util.SafeIntToInt32(params.TmdbId))
	if err != nil {
		if errors.Is(err, metadata.ErrNoProviders) || errors.Is(err, metadata.ErrNotFound) {
			return &ogen.GetPersonMetadataCreditsNotFound{}, nil
		}
		h.logger.Error("get person credits failed", slog.Any("error", err))
		return nil, err
	}

	if credits == nil {
		return &ogen.GetPersonMetadataCreditsNotFound{}, nil
	}

	response := &ogen.MetadataPersonCredits{
		TmdbID: ogen.NewOptInt(int(params.TmdbId)),
	}

	if len(credits.CastCredits) > 0 {
		cast := make([]ogen.MetadataMediaCredit, 0, len(credits.CastCredits))
		for _, c := range credits.CastCredits {
			credit := convertMediaCredit(c)
			cast = append(cast, credit)
		}
		response.Cast = cast
	}

	if len(credits.CrewCredits) > 0 {
		crew := make([]ogen.MetadataMediaCredit, 0, len(credits.CrewCredits))
		for _, c := range credits.CrewCredits {
			credit := convertMediaCredit(c)
			crew = append(crew, credit)
		}
		response.Crew = crew
	}

	return response, nil
}

// GetPersonMetadataImages gets person images from metadata provider.
func (h *Handler) GetPersonMetadataImages(ctx context.Context, params ogen.GetPersonMetadataImagesParams) (ogen.GetPersonMetadataImagesRes, error) {
	images, err := h.metadataService.GetPersonImages(ctx, util.SafeIntToInt32(params.TmdbId))
	if err != nil {
		if errors.Is(err, metadata.ErrNoProviders) || errors.Is(err, metadata.ErrNotFound) {
			return &ogen.GetPersonMetadataImagesNotFound{}, nil
		}
		h.logger.Error("get person images failed", slog.Any("error", err))
		return nil, err
	}

	if images == nil {
		return &ogen.GetPersonMetadataImagesNotFound{}, nil
	}

	return convertImages(images), nil
}

// ListMetadataProviders lists available metadata providers.
func (h *Handler) ListMetadataProviders(ctx context.Context) (ogen.ListMetadataProvidersRes, error) {
	providers := h.metadataService.GetProviders()

	result := make([]ogen.MetadataProvider, 0, len(providers))
	for _, p := range providers {
		capabilities := make([]string, 0, 5)
		if p.SupportsMovies() {
			capabilities = append(capabilities, "movie")
		}
		if p.SupportsTVShows() {
			capabilities = append(capabilities, "tvshow")
		}
		if p.SupportsPeople() {
			capabilities = append(capabilities, "person")
		}

		result = append(result, ogen.MetadataProvider{
			ID:           ogen.NewOptString(string(p.ID())),
			Name:         ogen.NewOptString(p.Name()),
			Capabilities: capabilities,
		})
	}

	return &ogen.MetadataProviderList{Providers: result}, nil
}

// --- Metadata conversion helpers ---

func convertImages(images *metadata.Images) *ogen.MetadataImages {
	response := &ogen.MetadataImages{}

	if len(images.Posters) > 0 {
		response.Posters = convertImageSlice(images.Posters)
	}
	if len(images.Backdrops) > 0 {
		response.Backdrops = convertImageSlice(images.Backdrops)
	}
	if len(images.Logos) > 0 {
		response.Logos = convertImageSlice(images.Logos)
	}
	if len(images.Stills) > 0 {
		response.Stills = convertImageSlice(images.Stills)
	}
	if len(images.Profiles) > 0 {
		response.Profiles = convertImageSlice(images.Profiles)
	}

	return response
}

func convertImageSlice(images []metadata.Image) []ogen.MetadataImage {
	result := make([]ogen.MetadataImage, 0, len(images))
	for _, img := range images {
		mi := ogen.MetadataImage{
			FilePath:    ogen.NewOptString(img.FilePath),
			AspectRatio: ogen.NewOptFloat64(img.AspectRatio),
			Width:       ogen.NewOptInt(img.Width),
			Height:      ogen.NewOptInt(img.Height),
		}
		if img.VoteAverage > 0 {
			mi.VoteAverage = ogen.NewOptFloat32(float32(img.VoteAverage))
		}
		if img.VoteCount > 0 {
			mi.VoteCount = ogen.NewOptInt(img.VoteCount)
		}
		if img.Language != nil {
			mi.Language = ogen.NewOptNilString(*img.Language)
		}
		result = append(result, mi)
	}
	return result
}

func convertExternalIDs(ids *metadata.ExternalIDs) *ogen.MetadataExternalIDs {
	response := &ogen.MetadataExternalIDs{}

	if ids.IMDbID != nil {
		response.ImdbID = ogen.NewOptNilString(*ids.IMDbID)
	}
	if ids.TVDbID != nil {
		response.TvdbID = ogen.NewOptNilInt(int(*ids.TVDbID))
	}
	if ids.TMDbID != nil {
		response.TmdbID = ogen.NewOptNilInt(int(*ids.TMDbID))
	}
	if ids.WikidataID != nil {
		response.WikidataID = ogen.NewOptNilString(*ids.WikidataID)
	}
	if ids.FacebookID != nil {
		response.FacebookID = ogen.NewOptNilString(*ids.FacebookID)
	}
	if ids.InstagramID != nil {
		response.InstagramID = ogen.NewOptNilString(*ids.InstagramID)
	}
	if ids.TwitterID != nil {
		response.TwitterID = ogen.NewOptNilString(*ids.TwitterID)
	}
	if ids.TikTokID != nil {
		response.TiktokID = ogen.NewOptNilString(*ids.TikTokID)
	}
	if ids.YouTubeID != nil {
		response.YoutubeID = ogen.NewOptNilString(*ids.YouTubeID)
	}

	return response
}

func convertCredits(credits *metadata.Credits) *ogen.MetadataCredits {
	response := &ogen.MetadataCredits{}

	if len(credits.Cast) > 0 {
		cast := make([]ogen.MetadataCastMember, 0, len(credits.Cast))
		for _, c := range credits.Cast {
			member := ogen.MetadataCastMember{
				Name:      ogen.NewOptString(c.Name),
				Character: ogen.NewOptString(c.Character),
				Order:     ogen.NewOptInt(c.Order),
			}
			if c.ProviderID != "" {
				var castID int
				_, _ = fmt.Sscanf(c.ProviderID, "%d", &castID)
				if castID > 0 {
					member.ID = ogen.NewOptInt(castID)
				}
			}
			if c.ProfilePath != nil && *c.ProfilePath != "" {
				member.ProfilePath = ogen.NewOptNilString(*c.ProfilePath)
			}
			cast = append(cast, member)
		}
		response.Cast = cast
	}

	if len(credits.Crew) > 0 {
		crew := make([]ogen.MetadataCrewMember, 0, len(credits.Crew))
		for _, c := range credits.Crew {
			member := ogen.MetadataCrewMember{
				Name:       ogen.NewOptString(c.Name),
				Job:        ogen.NewOptString(c.Job),
				Department: ogen.NewOptString(c.Department),
			}
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

	return response
}

func convertMovieSearchResults(results []metadata.MovieSearchResult, limit, total int) *ogen.MetadataSearchResults {
	response := &ogen.MetadataSearchResults{
		Page:         ogen.NewOptInt(1),
		TotalResults: ogen.NewOptInt(total),
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

	return response
}

func convertMediaCredit(c metadata.MediaCredit) ogen.MetadataMediaCredit {
	credit := ogen.MetadataMediaCredit{
		Title: ogen.NewOptString(c.Title),
	}

	switch c.MediaType {
	case "movie":
		credit.MediaType = ogen.NewOptMetadataMediaCreditMediaType(ogen.MetadataMediaCreditMediaTypeMovie)
	case "tv":
		credit.MediaType = ogen.NewOptMetadataMediaCreditMediaType(ogen.MetadataMediaCreditMediaTypeTv)
	}

	if c.MediaID != "" {
		var id int
		_, _ = fmt.Sscanf(c.MediaID, "%d", &id)
		if id > 0 {
			credit.TmdbID = ogen.NewOptInt(id)
		}
	}

	if c.Character != nil {
		credit.Character = ogen.NewOptNilString(*c.Character)
	}
	if c.Job != nil {
		credit.Job = ogen.NewOptNilString(*c.Job)
	}
	if c.Department != nil {
		credit.Department = ogen.NewOptNilString(*c.Department)
	}
	if c.ReleaseDate != nil {
		credit.ReleaseDate = ogen.NewOptNilDate(*c.ReleaseDate)
	}
	if c.PosterPath != nil {
		credit.PosterPath = ogen.NewOptNilString(*c.PosterPath)
	}
	if c.VoteAverage > 0 {
		credit.VoteAverage = ogen.NewOptFloat32(float32(c.VoteAverage))
	}
	if c.EpisodeCount != nil {
		credit.EpisodeCount = ogen.NewOptNilInt(*c.EpisodeCount)
	}

	return credit
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// searchMoviesViaRadarr searches for movies via the configured Radarr instance's lookup API.
func (h *Handler) searchMoviesViaRadarr(ctx context.Context, params ogen.SearchMoviesMetadataParams) (ogen.SearchMoviesMetadataRes, error) {
	if h.radarrService == nil {
		badReq := ogen.SearchMoviesMetadataBadRequest(ogen.Error{
			Code:    400,
			Message: "Radarr integration not configured",
		})
		return &badReq, nil
	}

	movies, err := h.radarrService.LookupMovie(ctx, params.Q)
	if err != nil {
		h.logger.Error("radarr movie lookup failed", slog.Any("error", err))
		return nil, err
	}

	limit := 20
	if params.Limit.Set {
		limit = int(params.Limit.Value)
	}

	response := &ogen.MetadataSearchResults{
		Page:         ogen.NewOptInt(1),
		TotalResults: ogen.NewOptInt(len(movies)),
		TotalPages:   ogen.NewOptInt(1),
		Results:      make([]ogen.MetadataSearchResult, 0, min(len(movies), limit)),
	}

	for i, m := range movies {
		if i >= limit {
			break
		}

		result := ogen.MetadataSearchResult{
			Title: ogen.NewOptString(m.Title),
		}
		if m.TMDbID > 0 {
			result.TmdbID = ogen.NewOptInt(m.TMDbID)
		}
		if m.OriginalTitle != "" {
			result.OriginalTitle = ogen.NewOptString(m.OriginalTitle)
		}
		if m.Overview != "" {
			result.Overview = ogen.NewOptNilString(m.Overview)
		}
		// Map best available rating
		if m.Ratings.TMDb != nil && m.Ratings.TMDb.Value > 0 {
			result.VoteAverage = ogen.NewOptFloat32(float32(m.Ratings.TMDb.Value))
			result.VoteCount = ogen.NewOptInt(m.Ratings.TMDb.Votes)
		} else if m.Ratings.IMDb != nil && m.Ratings.IMDb.Value > 0 {
			result.VoteAverage = ogen.NewOptFloat32(float32(m.Ratings.IMDb.Value))
			result.VoteCount = ogen.NewOptInt(m.Ratings.IMDb.Votes)
		}
		// Map poster from Radarr images
		for _, img := range m.Images {
			if img.CoverType == "poster" && img.RemoteURL != "" {
				result.PosterPath = ogen.NewOptNilString(img.RemoteURL)
			}
			if img.CoverType == "fanart" && img.RemoteURL != "" {
				result.BackdropPath = ogen.NewOptNilString(img.RemoteURL)
			}
		}

		response.Results = append(response.Results, result)
	}

	return response, nil
}

// searchTVShowsViaSonarr searches for TV shows via the configured Sonarr instance's lookup API.
func (h *Handler) searchTVShowsViaSonarr(ctx context.Context, params ogen.SearchTVShowsMetadataParams) (ogen.SearchTVShowsMetadataRes, error) {
	if h.sonarrService == nil {
		return &ogen.MetadataTVSearchResults{
			Page:         ogen.NewOptInt(1),
			TotalResults: ogen.NewOptInt(0),
			TotalPages:   ogen.NewOptInt(1),
			Results:      []ogen.MetadataTVSearchResult{},
		}, nil
	}

	series, err := h.sonarrService.LookupSeries(ctx, params.Q)
	if err != nil {
		h.logger.Error("sonarr series lookup failed", slog.Any("error", err))
		return nil, err
	}

	limit := 20
	if params.Limit.Set {
		limit = int(params.Limit.Value)
	}

	response := &ogen.MetadataTVSearchResults{
		Page:         ogen.NewOptInt(1),
		TotalResults: ogen.NewOptInt(len(series)),
		TotalPages:   ogen.NewOptInt(1),
		Results:      make([]ogen.MetadataTVSearchResult, 0, min(len(series), limit)),
	}

	for i, s := range series {
		if i >= limit {
			break
		}

		result := ogen.MetadataTVSearchResult{
			Name: ogen.NewOptString(s.Title),
		}
		if s.Overview != "" {
			result.Overview = ogen.NewOptNilString(s.Overview)
		}
		if s.Ratings.Value > 0 {
			result.VoteAverage = ogen.NewOptFloat32(float32(s.Ratings.Value))
			result.VoteCount = ogen.NewOptInt(s.Ratings.Votes)
		}
		// Map poster from Sonarr images
		for _, img := range s.Images {
			if img.CoverType == "poster" && img.RemoteURL != "" {
				result.PosterPath = ogen.NewOptNilString(img.RemoteURL)
			}
			if img.CoverType == "fanart" && img.RemoteURL != "" {
				result.BackdropPath = ogen.NewOptNilString(img.RemoteURL)
			}
		}

		response.Results = append(response.Results, result)
	}

	return response, nil
}
