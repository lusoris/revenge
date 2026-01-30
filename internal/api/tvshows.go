package api

import (
	"context"
	"errors"

	"github.com/google/uuid"

	gen "github.com/lusoris/revenge/api/generated"
	"github.com/lusoris/revenge/internal/content/tvshow"
)

// tvshowModuleDisabledError returns a generic error for endpoints that accept *gen.Error.
func tvshowModuleDisabledError() *gen.Error {
	return &gen.Error{
		Code:    "module_disabled",
		Message: "TV show module is not enabled",
	}
}

func (h *Handler) requireTVShowService() (*tvshow.Service, error) {
	if h.tvshowService == nil {
		return nil, ErrModuleDisabled
	}
	return h.tvshowService, nil
}

// =============================================================================
// Series Handlers
// =============================================================================

// ListSeries implements the listSeries operation.
func (h *Handler) ListSeries(ctx context.Context, params gen.ListSeriesParams) (gen.ListSeriesRes, error) {
	_, err := requireUser(ctx)
	if err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireTVShowService()
	if err != nil {
		return tvshowModuleDisabledError(), nil
	}

	listParams := tvshow.ListParams{
		Limit:     params.Limit.Or(20),
		Offset:    params.Offset.Or(0),
		SortBy:    string(params.SortBy.Or(gen.ListSeriesSortByTitle)),
		SortOrder: string(params.SortOrder.Or(gen.ListSeriesSortOrderAsc)),
	}

	var series []*tvshow.Series
	var total int64

	if params.LibraryId.IsSet() {
		series, total, err = svc.ListSeries(ctx, params.LibraryId.Value, listParams)
	} else if params.Query.IsSet() {
		series, err = svc.SearchSeries(ctx, params.Query.Value, listParams)
		if err == nil {
			total = int64(len(series))
		}
	} else {
		series, total, err = svc.ListAllSeries(ctx, listParams)
	}

	if err != nil {
		h.logger.Error("List series failed", "error", err)
		return &gen.Error{
			Code:    "list_failed",
			Message: "Failed to list series",
		}, nil
	}

	result := gen.SeriesListResponse{
		Series:     make([]gen.Series, len(series)),
		Pagination: paginationMeta(total, listParams.Limit, listParams.Offset),
	}

	for i, s := range series {
		result.Series[i] = seriesToAPI(s, nil)
	}

	return &result, nil
}

// ListRecentlyAddedSeries implements the listRecentlyAddedSeries operation.
func (h *Handler) ListRecentlyAddedSeries(ctx context.Context, params gen.ListRecentlyAddedSeriesParams) (gen.ListRecentlyAddedSeriesRes, error) {
	_, err := requireUser(ctx)
	if err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireTVShowService()
	if err != nil {
		return tvshowModuleDisabledError(), nil
	}

	listParams := tvshow.ListParams{
		Limit:     params.Limit.Or(20),
		SortBy:    "date_added",
		SortOrder: "desc",
	}

	series, _, err := svc.ListAllSeries(ctx, listParams)
	if err != nil {
		h.logger.Error("List recent series failed", "error", err)
		return &gen.Error{
			Code:    "list_failed",
			Message: "Failed to list series",
		}, nil
	}

	result := make(gen.ListRecentlyAddedSeriesOKApplicationJSON, len(series))
	for i, s := range series {
		result[i] = seriesToAPI(s, nil)
	}

	return &result, nil
}

// ListContinueWatchingSeries implements the listContinueWatchingSeries operation.
func (h *Handler) ListContinueWatchingSeries(ctx context.Context, params gen.ListContinueWatchingSeriesParams) (gen.ListContinueWatchingSeriesRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireTVShowService()
	if err != nil {
		return tvshowModuleDisabledError(), nil
	}

	progresses, err := svc.ListContinueWatchingSeries(ctx, usr.ID, params.Limit.Or(10))
	if err != nil {
		h.logger.Error("List continue watching series failed", "error", err)
		return &gen.Error{
			Code:    "list_failed",
			Message: "Failed to list series",
		}, nil
	}

	result := make(gen.ListContinueWatchingSeriesOKApplicationJSON, 0, len(progresses))
	for _, progress := range progresses {
		s, err := svc.GetSeries(ctx, progress.SeriesID)
		if err != nil {
			continue
		}
		result = append(result, seriesWithProgressToAPI(s, progress))
	}

	return &result, nil
}

// GetSeries implements the getSeries operation.
func (h *Handler) GetSeries(ctx context.Context, params gen.GetSeriesParams) (gen.GetSeriesRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.GetSeriesUnauthorized{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireTVShowService()
	if err != nil {
		return &gen.GetSeriesNotFound{Code: "module_disabled", Message: "TV show module is not enabled"}, nil
	}

	s, err := svc.GetSeriesWithRelations(ctx, params.SeriesId)
	if err != nil {
		if errors.Is(err, tvshow.ErrSeriesNotFoundInService) {
			return &gen.GetSeriesNotFound{
				Code:    "not_found",
				Message: "Series not found",
			}, nil
		}
		h.logger.Error("Get series failed", "error", err, "series_id", params.SeriesId)
		return &gen.GetSeriesNotFound{
			Code:    "get_failed",
			Message: "Failed to get series",
		}, nil
	}

	userData := h.getUserSeriesData(ctx, usr.ID, s.ID)
	result := seriesFullToAPI(s, userData)
	return &result, nil
}

// UpdateSeries implements the updateSeries operation.
func (h *Handler) UpdateSeries(ctx context.Context, req *gen.SeriesUpdate, params gen.UpdateSeriesParams) (gen.UpdateSeriesRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.UpdateSeriesUnauthorized{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireTVShowService()
	if err != nil {
		return &gen.UpdateSeriesNotFound{Code: "module_disabled", Message: "TV show module is not enabled"}, nil
	}

	if !usr.IsAdmin {
		return &gen.UpdateSeriesForbidden{
			Code:    "forbidden",
			Message: "Admin access required",
		}, nil
	}

	s, err := svc.GetSeries(ctx, params.SeriesId)
	if err != nil {
		if errors.Is(err, tvshow.ErrSeriesNotFoundInService) {
			return &gen.UpdateSeriesNotFound{
				Code:    "not_found",
				Message: "Series not found",
			}, nil
		}
		return &gen.UpdateSeriesNotFound{
			Code:    "get_failed",
			Message: "Failed to get series",
		}, nil
	}

	if req.Title.IsSet() {
		s.Title = req.Title.Value
	}
	if req.OriginalTitle.IsSet() {
		s.OriginalTitle = req.OriginalTitle.Value
	}
	if req.SortTitle.IsSet() {
		s.SortTitle = req.SortTitle.Value
	}
	if req.Overview.IsSet() {
		s.Overview = req.Overview.Value
	}
	if req.Tagline.IsSet() {
		s.Tagline = req.Tagline.Value
	}
	if req.ContentRating.IsSet() {
		s.ContentRating = req.ContentRating.Value
	}

	if err := svc.UpdateSeries(ctx, s); err != nil {
		h.logger.Error("Update series failed", "error", err, "series_id", params.SeriesId)
		return &gen.UpdateSeriesNotFound{
			Code:    "update_failed",
			Message: "Failed to update series",
		}, nil
	}

	result := seriesToAPI(s, nil)
	return &result, nil
}

// DeleteSeries implements the deleteSeries operation.
func (h *Handler) DeleteSeries(ctx context.Context, params gen.DeleteSeriesParams) (gen.DeleteSeriesRes, error) {
	svc, err := h.requireTVShowService()
	if err != nil {
		return &gen.DeleteSeriesNotFound{Code: "module_disabled", Message: "TV show module is not enabled"}, nil
	}

	_, err = requireAdmin(ctx)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) {
			return &gen.DeleteSeriesUnauthorized{
				Code:    "unauthorized",
				Message: "Not authenticated",
			}, nil
		}
		return &gen.DeleteSeriesForbidden{
			Code:    "forbidden",
			Message: "Admin access required",
		}, nil
	}

	if err := svc.DeleteSeries(ctx, params.SeriesId); err != nil {
		if errors.Is(err, tvshow.ErrSeriesNotFoundInService) {
			return &gen.DeleteSeriesNotFound{
				Code:    "not_found",
				Message: "Series not found",
			}, nil
		}
		h.logger.Error("Delete series failed", "error", err, "series_id", params.SeriesId)
		return &gen.DeleteSeriesNotFound{
			Code:    "delete_failed",
			Message: "Failed to delete series",
		}, nil
	}

	return &gen.DeleteSeriesNoContent{}, nil
}

// RefreshSeriesMetadata implements the refreshSeriesMetadata operation.
func (h *Handler) RefreshSeriesMetadata(ctx context.Context, params gen.RefreshSeriesMetadataParams) (gen.RefreshSeriesMetadataRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.RefreshSeriesMetadataUnauthorized{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireTVShowService()
	if err != nil {
		return &gen.RefreshSeriesMetadataNotFound{Code: "module_disabled", Message: "TV show module is not enabled"}, nil
	}

	s, err := svc.GetSeries(ctx, params.SeriesId)
	if err != nil {
		if errors.Is(err, tvshow.ErrSeriesNotFoundInService) {
			return &gen.RefreshSeriesMetadataNotFound{
				Code:    "not_found",
				Message: "Series not found",
			}, nil
		}
		return &gen.RefreshSeriesMetadataNotFound{
			Code:    "get_failed",
			Message: "Failed to get series",
		}, nil
	}

	if h.riverClient != nil {
		_, err = h.riverClient.Insert(ctx, tvshow.EnrichSeriesMetadataArgs{
			SeriesID: s.ID,
			TmdbID:   s.TmdbID,
			TvdbID:   s.TvdbID,
			ImdbID:   s.ImdbID,
			Title:    s.Title,
			Year:     s.Year,
		}, nil)
		if err != nil {
			h.logger.Error("Queue metadata refresh failed", "error", err, "series_id", s.ID, "user_id", usr.ID)
		}
	}

	return &gen.JobResponse{
		JobId:   uuid.New().String(),
		Status:  gen.JobResponseStatusQueued,
		Message: gen.NewOptString("Metadata refresh queued"),
	}, nil
}

// =============================================================================
// Season Handlers
// =============================================================================

// ListSeasons implements the listSeasons operation.
func (h *Handler) ListSeasons(ctx context.Context, params gen.ListSeasonsParams) (gen.ListSeasonsRes, error) {
	_, err := requireUser(ctx)
	if err != nil {
		return &gen.ListSeasonsUnauthorized{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireTVShowService()
	if err != nil {
		return &gen.ListSeasonsNotFound{Code: "module_disabled", Message: "TV show module is not enabled"}, nil
	}

	seasons, err := svc.ListSeasons(ctx, params.SeriesId)
	if err != nil {
		h.logger.Error("List seasons failed", "error", err, "series_id", params.SeriesId)
		return &gen.ListSeasonsNotFound{
			Code:    "list_failed",
			Message: "Failed to list seasons",
		}, nil
	}

	result := make(gen.ListSeasonsOKApplicationJSON, len(seasons))
	for i, s := range seasons {
		result[i] = seasonToAPI(s)
	}

	return &result, nil
}

// GetSeason implements the getSeason operation.
func (h *Handler) GetSeason(ctx context.Context, params gen.GetSeasonParams) (gen.GetSeasonRes, error) {
	_, err := requireUser(ctx)
	if err != nil {
		return &gen.GetSeasonUnauthorized{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireTVShowService()
	if err != nil {
		return &gen.GetSeasonNotFound{Code: "module_disabled", Message: "TV show module is not enabled"}, nil
	}

	season, err := svc.GetSeasonWithEpisodes(ctx, params.SeasonId)
	if err != nil {
		if errors.Is(err, tvshow.ErrSeasonNotFoundInService) {
			return &gen.GetSeasonNotFound{
				Code:    "not_found",
				Message: "Season not found",
			}, nil
		}
		h.logger.Error("Get season failed", "error", err, "season_id", params.SeasonId)
		return &gen.GetSeasonNotFound{
			Code:    "get_failed",
			Message: "Failed to get season",
		}, nil
	}

	result := seasonFullToAPI(season)
	return &result, nil
}

// ListEpisodes implements the listEpisodes operation.
func (h *Handler) ListEpisodes(ctx context.Context, params gen.ListEpisodesParams) (gen.ListEpisodesRes, error) {
	_, err := requireUser(ctx)
	if err != nil {
		return &gen.ListEpisodesUnauthorized{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireTVShowService()
	if err != nil {
		return &gen.ListEpisodesNotFound{Code: "module_disabled", Message: "TV show module is not enabled"}, nil
	}

	episodes, err := svc.ListEpisodesBySeason(ctx, params.SeasonId)
	if err != nil {
		h.logger.Error("List episodes failed", "error", err, "season_id", params.SeasonId)
		return &gen.ListEpisodesNotFound{
			Code:    "list_failed",
			Message: "Failed to list episodes",
		}, nil
	}

	result := make(gen.ListEpisodesOKApplicationJSON, len(episodes))
	for i, e := range episodes {
		result[i] = episodeToAPI(e)
	}

	return &result, nil
}

// =============================================================================
// Episode Handlers
// =============================================================================

// GetEpisode implements the getEpisode operation.
func (h *Handler) GetEpisode(ctx context.Context, params gen.GetEpisodeParams) (gen.GetEpisodeRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.GetEpisodeUnauthorized{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireTVShowService()
	if err != nil {
		return &gen.GetEpisodeNotFound{Code: "module_disabled", Message: "TV show module is not enabled"}, nil
	}

	episode, err := svc.GetEpisodeWithRelations(ctx, params.EpisodeId)
	if err != nil {
		if errors.Is(err, tvshow.ErrEpisodeNotFoundInService) {
			return &gen.GetEpisodeNotFound{
				Code:    "not_found",
				Message: "Episode not found",
			}, nil
		}
		h.logger.Error("Get episode failed", "error", err, "episode_id", params.EpisodeId)
		return &gen.GetEpisodeNotFound{
			Code:    "get_failed",
			Message: "Failed to get episode",
		}, nil
	}

	userData := h.getUserEpisodeData(ctx, usr.ID, episode.ID)
	result := episodeFullToAPI(episode, userData)
	return &result, nil
}

// GetNextEpisode implements the getNextEpisode operation.
func (h *Handler) GetNextEpisode(ctx context.Context, params gen.GetNextEpisodeParams) (gen.GetNextEpisodeRes, error) {
	_, err := requireUser(ctx)
	if err != nil {
		return &gen.GetNextEpisodeUnauthorized{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireTVShowService()
	if err != nil {
		return &gen.GetNextEpisodeNotFound{Code: "module_disabled", Message: "TV show module is not enabled"}, nil
	}

	next, err := svc.GetNextEpisode(ctx, params.EpisodeId)
	if err != nil {
		if errors.Is(err, tvshow.ErrEpisodeNotFoundInService) {
			return &gen.GetNextEpisodeNotFound{
				Code:    "not_found",
				Message: "No next episode",
			}, nil
		}
		h.logger.Error("Get next episode failed", "error", err, "episode_id", params.EpisodeId)
		return &gen.GetNextEpisodeNotFound{
			Code:    "get_failed",
			Message: "Failed to get next episode",
		}, nil
	}

	result := episodeToAPI(next)
	return &result, nil
}

// ListUpcomingEpisodes implements the listUpcomingEpisodes operation.
func (h *Handler) ListUpcomingEpisodes(ctx context.Context, params gen.ListUpcomingEpisodesParams) (gen.ListUpcomingEpisodesRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireTVShowService()
	if err != nil {
		return tvshowModuleDisabledError(), nil
	}

	libs, err := h.libraryService.ListForUser(ctx, usr.ID)
	if err != nil {
		h.logger.Error("Get libraries failed", "error", err)
		return &gen.Error{
			Code:    "list_failed",
			Message: "Failed to get libraries",
		}, nil
	}

	libraryIDs := make([]uuid.UUID, 0, len(libs))
	for _, l := range libs {
		if l.Library.Module == "tvshow" {
			libraryIDs = append(libraryIDs, l.Library.ID)
		}
	}

	episodes, err := svc.ListUpcomingEpisodes(ctx, libraryIDs, params.Limit.Or(20))
	if err != nil {
		h.logger.Error("List upcoming episodes failed", "error", err)
		return &gen.Error{
			Code:    "list_failed",
			Message: "Failed to list episodes",
		}, nil
	}

	result := make(gen.ListUpcomingEpisodesOKApplicationJSON, len(episodes))
	for i, e := range episodes {
		result[i] = episodeWithSeriesToAPI(e, svc)
	}

	return &result, nil
}

// ListRecentlyAiredEpisodes implements the listRecentlyAiredEpisodes operation.
func (h *Handler) ListRecentlyAiredEpisodes(ctx context.Context, params gen.ListRecentlyAiredEpisodesParams) (gen.ListRecentlyAiredEpisodesRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireTVShowService()
	if err != nil {
		return tvshowModuleDisabledError(), nil
	}

	libs, err := h.libraryService.ListForUser(ctx, usr.ID)
	if err != nil {
		h.logger.Error("Get libraries failed", "error", err)
		return &gen.Error{
			Code:    "list_failed",
			Message: "Failed to get libraries",
		}, nil
	}

	libraryIDs := make([]uuid.UUID, 0, len(libs))
	for _, l := range libs {
		if l.Library.Module == "tvshow" {
			libraryIDs = append(libraryIDs, l.Library.ID)
		}
	}

	episodes, err := svc.ListRecentlyAiredEpisodes(ctx, libraryIDs, params.Limit.Or(20))
	if err != nil {
		h.logger.Error("List recent episodes failed", "error", err)
		return &gen.Error{
			Code:    "list_failed",
			Message: "Failed to list episodes",
		}, nil
	}

	result := make(gen.ListRecentlyAiredEpisodesOKApplicationJSON, len(episodes))
	for i, e := range episodes {
		result[i] = episodeWithSeriesToAPI(e, svc)
	}

	return &result, nil
}

// =============================================================================
// User Data Handlers - Series
// =============================================================================

// AddSeriesToFavorites implements the addSeriesToFavorites operation.
func (h *Handler) AddSeriesToFavorites(ctx context.Context, params gen.AddSeriesToFavoritesParams) (gen.AddSeriesToFavoritesRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.AddSeriesToFavoritesUnauthorized{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireTVShowService()
	if err != nil {
		return &gen.AddSeriesToFavoritesNotFound{Code: "module_disabled", Message: "TV show module is not enabled"}, nil
	}

	_, err = svc.GetSeries(ctx, params.SeriesId)
	if err != nil {
		if errors.Is(err, tvshow.ErrSeriesNotFoundInService) {
			return &gen.AddSeriesToFavoritesNotFound{
				Code:    "not_found",
				Message: "Series not found",
			}, nil
		}
		return &gen.AddSeriesToFavoritesNotFound{
			Code:    "get_failed",
			Message: "Failed to get series",
		}, nil
	}

	if err := svc.AddSeriesFavorite(ctx, usr.ID, params.SeriesId); err != nil {
		h.logger.Error("Add favorite failed", "error", err, "series_id", params.SeriesId, "user_id", usr.ID)
		return &gen.AddSeriesToFavoritesNotFound{
			Code:    "add_failed",
			Message: "Failed to add to favorites",
		}, nil
	}

	return &gen.AddSeriesToFavoritesNoContent{}, nil
}

// RemoveSeriesFromFavorites implements the removeSeriesFromFavorites operation.
func (h *Handler) RemoveSeriesFromFavorites(ctx context.Context, params gen.RemoveSeriesFromFavoritesParams) (gen.RemoveSeriesFromFavoritesRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireTVShowService()
	if err != nil {
		return tvshowModuleDisabledError(), nil
	}

	if err := svc.RemoveSeriesFavorite(ctx, usr.ID, params.SeriesId); err != nil {
		h.logger.Error("Remove favorite failed", "error", err, "series_id", params.SeriesId, "user_id", usr.ID)
	}

	return &gen.RemoveSeriesFromFavoritesNoContent{}, nil
}

// AddSeriesToWatchlist implements the addSeriesToWatchlist operation.
func (h *Handler) AddSeriesToWatchlist(ctx context.Context, params gen.AddSeriesToWatchlistParams) (gen.AddSeriesToWatchlistRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.AddSeriesToWatchlistUnauthorized{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireTVShowService()
	if err != nil {
		return &gen.AddSeriesToWatchlistNotFound{Code: "module_disabled", Message: "TV show module is not enabled"}, nil
	}

	_, err = svc.GetSeries(ctx, params.SeriesId)
	if err != nil {
		if errors.Is(err, tvshow.ErrSeriesNotFoundInService) {
			return &gen.AddSeriesToWatchlistNotFound{
				Code:    "not_found",
				Message: "Series not found",
			}, nil
		}
		return &gen.AddSeriesToWatchlistNotFound{
			Code:    "get_failed",
			Message: "Failed to get series",
		}, nil
	}

	if err := svc.AddSeriesToWatchlist(ctx, usr.ID, params.SeriesId); err != nil {
		h.logger.Error("Add to watchlist failed", "error", err, "series_id", params.SeriesId, "user_id", usr.ID)
		return &gen.AddSeriesToWatchlistNotFound{
			Code:    "add_failed",
			Message: "Failed to add to watchlist",
		}, nil
	}

	return &gen.AddSeriesToWatchlistNoContent{}, nil
}

// RemoveSeriesFromWatchlist implements the removeSeriesFromWatchlist operation.
func (h *Handler) RemoveSeriesFromWatchlist(ctx context.Context, params gen.RemoveSeriesFromWatchlistParams) (gen.RemoveSeriesFromWatchlistRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireTVShowService()
	if err != nil {
		return tvshowModuleDisabledError(), nil
	}

	if err := svc.RemoveSeriesFromWatchlist(ctx, usr.ID, params.SeriesId); err != nil {
		h.logger.Error("Remove from watchlist failed", "error", err, "series_id", params.SeriesId, "user_id", usr.ID)
	}

	return &gen.RemoveSeriesFromWatchlistNoContent{}, nil
}

// SetSeriesRating implements the setSeriesRating operation.
func (h *Handler) SetSeriesRating(ctx context.Context, req *gen.UserRatingRequest, params gen.SetSeriesRatingParams) (gen.SetSeriesRatingRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.SetSeriesRatingUnauthorized{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireTVShowService()
	if err != nil {
		return &gen.SetSeriesRatingNotFound{Code: "module_disabled", Message: "TV show module is not enabled"}, nil
	}

	_, err = svc.GetSeries(ctx, params.SeriesId)
	if err != nil {
		if errors.Is(err, tvshow.ErrSeriesNotFoundInService) {
			return &gen.SetSeriesRatingNotFound{
				Code:    "not_found",
				Message: "Series not found",
			}, nil
		}
		return &gen.SetSeriesRatingNotFound{
			Code:    "get_failed",
			Message: "Failed to get series",
		}, nil
	}

	review := ""
	if req.Review.IsSet() {
		review = req.Review.Value
	}

	if err := svc.SetSeriesUserRating(ctx, usr.ID, params.SeriesId, req.Rating, review); err != nil {
		h.logger.Error("Set rating failed", "error", err, "series_id", params.SeriesId, "user_id", usr.ID)
		return &gen.SetSeriesRatingNotFound{
			Code:    "set_failed",
			Message: "Failed to set rating",
		}, nil
	}

	rating, err := svc.GetSeriesUserRating(ctx, usr.ID, params.SeriesId)
	if err != nil || rating == nil {
		return &gen.UserRating{
			Rating: req.Rating,
			Review: req.Review,
		}, nil
	}

	return &gen.UserRating{
		Rating:    rating.Rating,
		Review:    gen.NewOptString(rating.Review),
		CreatedAt: gen.NewOptDateTime(rating.CreatedAt),
		UpdatedAt: gen.NewOptDateTime(rating.UpdatedAt),
	}, nil
}

// DeleteSeriesRating implements the deleteSeriesRating operation.
func (h *Handler) DeleteSeriesRating(ctx context.Context, params gen.DeleteSeriesRatingParams) (gen.DeleteSeriesRatingRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireTVShowService()
	if err != nil {
		return tvshowModuleDisabledError(), nil
	}

	if err := svc.DeleteSeriesUserRating(ctx, usr.ID, params.SeriesId); err != nil {
		h.logger.Error("Delete rating failed", "error", err, "series_id", params.SeriesId, "user_id", usr.ID)
	}

	return &gen.DeleteSeriesRatingNoContent{}, nil
}

// =============================================================================
// User Data Handlers - Episodes
// =============================================================================

// MarkEpisodeAsWatched implements the markEpisodeAsWatched operation.
func (h *Handler) MarkEpisodeAsWatched(ctx context.Context, params gen.MarkEpisodeAsWatchedParams) (gen.MarkEpisodeAsWatchedRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.MarkEpisodeAsWatchedUnauthorized{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireTVShowService()
	if err != nil {
		return &gen.MarkEpisodeAsWatchedNotFound{Code: "module_disabled", Message: "TV show module is not enabled"}, nil
	}

	_, err = svc.GetEpisode(ctx, params.EpisodeId)
	if err != nil {
		if errors.Is(err, tvshow.ErrEpisodeNotFoundInService) {
			return &gen.MarkEpisodeAsWatchedNotFound{
				Code:    "not_found",
				Message: "Episode not found",
			}, nil
		}
		return &gen.MarkEpisodeAsWatchedNotFound{
			Code:    "get_failed",
			Message: "Failed to get episode",
		}, nil
	}

	if err := svc.MarkEpisodeAsWatched(ctx, usr.ID, params.EpisodeId); err != nil {
		h.logger.Error("Mark watched failed", "error", err, "episode_id", params.EpisodeId, "user_id", usr.ID)
		return &gen.MarkEpisodeAsWatchedNotFound{
			Code:    "mark_failed",
			Message: "Failed to mark as watched",
		}, nil
	}

	return &gen.MarkEpisodeAsWatchedNoContent{}, nil
}

// MarkEpisodeAsUnwatched implements the markEpisodeAsUnwatched operation.
func (h *Handler) MarkEpisodeAsUnwatched(ctx context.Context, params gen.MarkEpisodeAsUnwatchedParams) (gen.MarkEpisodeAsUnwatchedRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireTVShowService()
	if err != nil {
		return tvshowModuleDisabledError(), nil
	}

	if err := svc.MarkEpisodeAsUnwatched(ctx, usr.ID, params.EpisodeId); err != nil {
		h.logger.Error("Mark unwatched failed", "error", err, "episode_id", params.EpisodeId, "user_id", usr.ID)
	}

	return &gen.MarkEpisodeAsUnwatchedNoContent{}, nil
}

// SetEpisodeRating implements the setEpisodeRating operation.
func (h *Handler) SetEpisodeRating(ctx context.Context, req *gen.UserRatingRequest, params gen.SetEpisodeRatingParams) (gen.SetEpisodeRatingRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.SetEpisodeRatingUnauthorized{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireTVShowService()
	if err != nil {
		return &gen.SetEpisodeRatingNotFound{Code: "module_disabled", Message: "TV show module is not enabled"}, nil
	}

	_, err = svc.GetEpisode(ctx, params.EpisodeId)
	if err != nil {
		if errors.Is(err, tvshow.ErrEpisodeNotFoundInService) {
			return &gen.SetEpisodeRatingNotFound{
				Code:    "not_found",
				Message: "Episode not found",
			}, nil
		}
		return &gen.SetEpisodeRatingNotFound{
			Code:    "get_failed",
			Message: "Failed to get episode",
		}, nil
	}

	if err := svc.SetEpisodeUserRating(ctx, usr.ID, params.EpisodeId, req.Rating); err != nil {
		h.logger.Error("Set rating failed", "error", err, "episode_id", params.EpisodeId, "user_id", usr.ID)
		return &gen.SetEpisodeRatingNotFound{
			Code:    "set_failed",
			Message: "Failed to set rating",
		}, nil
	}

	rating, err := svc.GetEpisodeUserRating(ctx, usr.ID, params.EpisodeId)
	if err != nil || rating == nil {
		return &gen.UserRating{
			Rating: req.Rating,
		}, nil
	}

	return &gen.UserRating{
		Rating:    rating.Rating,
		CreatedAt: gen.NewOptDateTime(rating.CreatedAt),
		UpdatedAt: gen.NewOptDateTime(rating.UpdatedAt),
	}, nil
}

// DeleteEpisodeRating implements the deleteEpisodeRating operation.
func (h *Handler) DeleteEpisodeRating(ctx context.Context, params gen.DeleteEpisodeRatingParams) (gen.DeleteEpisodeRatingRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireTVShowService()
	if err != nil {
		return tvshowModuleDisabledError(), nil
	}

	if err := svc.DeleteEpisodeUserRating(ctx, usr.ID, params.EpisodeId); err != nil {
		h.logger.Error("Delete rating failed", "error", err, "episode_id", params.EpisodeId, "user_id", usr.ID)
	}

	return &gen.DeleteEpisodeRatingNoContent{}, nil
}

// =============================================================================
// User Lists
// =============================================================================

// ListMyFavoriteSeries implements the listMyFavoriteSeries operation.
func (h *Handler) ListMyFavoriteSeries(ctx context.Context, params gen.ListMyFavoriteSeriesParams) (gen.ListMyFavoriteSeriesRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireTVShowService()
	if err != nil {
		return tvshowModuleDisabledError(), nil
	}

	listParams := tvshow.ListParams{
		Limit:  params.Limit.Or(20),
		Offset: params.Offset.Or(0),
	}

	series, total, err := svc.ListSeriesFavorites(ctx, usr.ID, listParams)
	if err != nil {
		h.logger.Error("List favorites failed", "error", err, "user_id", usr.ID)
		return &gen.Error{
			Code:    "list_failed",
			Message: "Failed to list favorites",
		}, nil
	}

	result := gen.SeriesListResponse{
		Series:     make([]gen.Series, len(series)),
		Pagination: paginationMeta(total, listParams.Limit, listParams.Offset),
	}

	for i, s := range series {
		result.Series[i] = seriesToAPI(s, nil)
	}

	return &result, nil
}

// ListMySeriesWatchlist implements the listMySeriesWatchlist operation.
func (h *Handler) ListMySeriesWatchlist(ctx context.Context, params gen.ListMySeriesWatchlistParams) (gen.ListMySeriesWatchlistRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireTVShowService()
	if err != nil {
		return tvshowModuleDisabledError(), nil
	}

	listParams := tvshow.ListParams{
		Limit:  params.Limit.Or(20),
		Offset: params.Offset.Or(0),
	}

	series, total, err := svc.ListSeriesWatchlist(ctx, usr.ID, listParams)
	if err != nil {
		h.logger.Error("List watchlist failed", "error", err, "user_id", usr.ID)
		return &gen.Error{
			Code:    "list_failed",
			Message: "Failed to list watchlist",
		}, nil
	}

	result := gen.SeriesListResponse{
		Series:     make([]gen.Series, len(series)),
		Pagination: paginationMeta(total, listParams.Limit, listParams.Offset),
	}

	for i, s := range series {
		result.Series[i] = seriesToAPI(s, nil)
	}

	return &result, nil
}

// =============================================================================
// Network Handlers
// =============================================================================

// ListNetworks implements the listNetworks operation.
func (h *Handler) ListNetworks(ctx context.Context, params gen.ListNetworksParams) (gen.ListNetworksRes, error) {
	_, err := requireUser(ctx)
	if err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	// Networks are not yet fully implemented
	return &gen.ListNetworksOKApplicationJSON{}, nil
}

// GetNetwork implements the getNetwork operation.
func (h *Handler) GetNetwork(ctx context.Context, params gen.GetNetworkParams) (gen.GetNetworkRes, error) {
	_, err := requireUser(ctx)
	if err != nil {
		return &gen.GetNetworkUnauthorized{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	return &gen.GetNetworkNotFound{
		Code:    "not_implemented",
		Message: "Networks not yet implemented",
	}, nil
}

// ListNetworkSeries implements the listNetworkSeries operation.
func (h *Handler) ListNetworkSeries(ctx context.Context, params gen.ListNetworkSeriesParams) (gen.ListNetworkSeriesRes, error) {
	_, err := requireUser(ctx)
	if err != nil {
		return &gen.ListNetworkSeriesUnauthorized{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	return &gen.ListNetworkSeriesNotFound{
		Code:    "not_implemented",
		Message: "Networks not yet implemented",
	}, nil
}

// =============================================================================
// Helper Functions
// =============================================================================

func (h *Handler) getUserSeriesData(ctx context.Context, userID, seriesID uuid.UUID) *gen.SeriesUserData {
	data := &gen.SeriesUserData{}

	svc, err := h.requireTVShowService()
	if err != nil {
		return data
	}

	if favorite, err := svc.IsSeriesFavorite(ctx, userID, seriesID); err == nil {
		data.IsFavorite = gen.NewOptBool(favorite)
	}

	if watchlist, err := svc.IsSeriesInWatchlist(ctx, userID, seriesID); err == nil {
		data.IsInWatchlist = gen.NewOptBool(watchlist)
	}

	if rating, err := svc.GetSeriesUserRating(ctx, userID, seriesID); err == nil && rating != nil {
		data.UserRating = gen.NewOptFloat32(float32(rating.Rating))
	}

	// Note: SeriesUserData doesn't include watch progress fields
	// Watch progress is returned separately via SeriesWatchProgress

	return data
}

func (h *Handler) getUserEpisodeData(ctx context.Context, userID, episodeID uuid.UUID) *gen.EpisodeUserData {
	data := &gen.EpisodeUserData{}

	svc, err := h.requireTVShowService()
	if err != nil {
		return data
	}

	if watched, err := svc.IsEpisodeWatched(ctx, userID, episodeID); err == nil {
		data.IsWatched = gen.NewOptBool(watched)
	}

	if rating, err := svc.GetEpisodeUserRating(ctx, userID, episodeID); err == nil && rating != nil {
		data.UserRating = gen.NewOptFloat32(float32(rating.Rating))
	}

	if history, err := svc.GetEpisodeWatchHistory(ctx, userID, episodeID); err == nil && history != nil {
		data.PlaybackPosition = gen.NewOptInt64(history.PositionTicks)
		data.PlayedPercentage = gen.NewOptFloat32(float32(history.PlayedPercentage))
		data.LastPlayedAt = gen.NewOptDateTime(history.LastUpdatedAt)
	}

	return data
}
