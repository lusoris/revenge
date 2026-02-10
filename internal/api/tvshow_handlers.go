package api

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/content/tvshow"
	"github.com/lusoris/revenge/internal/util"
)

// TV Show errors
var (
	ErrSeriesNotFound  = errors.New("series not found")
	ErrSeasonNotFound  = errors.New("season not found")
	ErrEpisodeNotFound = errors.New("episode not found")
)

// ListTVShows returns a paginated list of TV shows.
func (h *Handler) ListTVShows(ctx context.Context, params ogen.ListTVShowsParams) (ogen.ListTVShowsRes, error) {
	filters := tvshow.SeriesListFilters{
		OrderBy: string(params.OrderBy.Or("created_at")),
		Limit:   util.SafeIntToInt32(params.Limit.Or(20)),
		Offset:  util.SafeIntToInt32(params.Offset.Or(0)),
	}

	series, err := h.tvshowService.ListSeries(ctx, filters)
	if err != nil {
		return nil, err
	}

	// Get total count for pagination
	total, err := h.tvshowService.CountSeries(ctx)
	if err != nil {
		total = int64(len(series))
	}

	lang := h.GetMetadataLanguage(ctx)
	localized := LocalizeSeriesList(series, lang)

	result := make([]ogen.TVSeries, len(localized))
	for i, s := range localized {
		result[i] = *seriesToOgen(&s)
	}

	return &ogen.TVShowListResponse{
		Items: result,
		Total: total,
	}, nil
}

// SearchTVShows searches TV shows by title.
func (h *Handler) SearchTVShows(ctx context.Context, params ogen.SearchTVShowsParams) (ogen.SearchTVShowsRes, error) {
	series, err := h.tvshowService.SearchSeries(ctx, params.Query, util.SafeIntToInt32(params.Limit.Or(20)), util.SafeIntToInt32(params.Offset.Or(0)))
	if err != nil {
		return nil, err
	}

	lang := h.GetMetadataLanguage(ctx)
	localized := LocalizeSeriesList(series, lang)

	result := make([]ogen.TVSeries, len(localized))
	for i, s := range localized {
		result[i] = *seriesToOgen(&s)
	}

	return (*ogen.SearchTVShowsOKApplicationJSON)(&result), nil
}

// GetRecentlyAddedTVShows returns recently added TV shows.
func (h *Handler) GetRecentlyAddedTVShows(ctx context.Context, params ogen.GetRecentlyAddedTVShowsParams) (ogen.GetRecentlyAddedTVShowsRes, error) {
	limit := util.SafeIntToInt32(params.Limit.Or(20))
	offset := util.SafeIntToInt32(params.Offset.Or(0))

	series, total, err := h.tvshowService.ListRecentlyAdded(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	lang := h.GetMetadataLanguage(ctx)
	localized := LocalizeSeriesList(series, lang)

	items := make([]ogen.TVSeries, len(localized))
	for i, s := range localized {
		items[i] = *seriesToOgen(&s)
	}

	page := (offset / limit) + 1
	return &ogen.TVSeriesListResponse{
		Items:    items,
		Total:    total,
		Page:     ogen.NewOptInt(int(page)),
		PageSize: ogen.NewOptInt(int(limit)),
	}, nil
}

// GetTVContinueWatching returns the continue watching list for TV shows.
func (h *Handler) GetTVContinueWatching(ctx context.Context, params ogen.GetTVContinueWatchingParams) (ogen.GetTVContinueWatchingRes, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	items, err := h.tvshowService.GetContinueWatching(ctx, userID, util.SafeIntToInt32(params.Limit.Or(20)))
	if err != nil {
		return nil, err
	}

	lang := h.GetMetadataLanguage(ctx)
	for i := range items {
		if items[i].Series != nil {
			items[i].Series = LocalizeSeries(items[i].Series, lang)
		}
	}

	result := make([]ogen.TVContinueWatchingItem, len(items))
	for i, item := range items {
		result[i] = *tvContinueWatchingItemToOgen(&item)
	}

	return (*ogen.GetTVContinueWatchingOKApplicationJSON)(&result), nil
}

// GetUserTVStats returns user's TV watching statistics.
func (h *Handler) GetUserTVStats(ctx context.Context) (ogen.GetUserTVStatsRes, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	stats, err := h.tvshowService.GetUserStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	return userTVStatsToOgen(stats, userID.String()), nil
}

// GetRecentEpisodes returns recently aired episodes.
func (h *Handler) GetRecentEpisodes(ctx context.Context, params ogen.GetRecentEpisodesParams) (ogen.GetRecentEpisodesRes, error) {
	episodes, err := h.tvshowService.ListRecentEpisodes(ctx, util.SafeIntToInt32(params.Limit.Or(20)), util.SafeIntToInt32(params.Offset.Or(0)))
	if err != nil {
		return nil, err
	}

	result := make([]ogen.EpisodeWithSeriesInfo, len(episodes))
	for i, e := range episodes {
		result[i] = *episodeWithSeriesInfoToOgen(&e)
	}

	return (*ogen.GetRecentEpisodesOKApplicationJSON)(&result), nil
}

// GetUpcomingEpisodes returns upcoming episodes.
func (h *Handler) GetUpcomingEpisodes(ctx context.Context, params ogen.GetUpcomingEpisodesParams) (ogen.GetUpcomingEpisodesRes, error) {
	episodes, err := h.tvshowService.ListUpcomingEpisodes(ctx, util.SafeIntToInt32(params.Limit.Or(20)), util.SafeIntToInt32(params.Offset.Or(0)))
	if err != nil {
		return nil, err
	}

	result := make([]ogen.EpisodeWithSeriesInfo, len(episodes))
	for i, e := range episodes {
		result[i] = *episodeWithSeriesInfoToOgen(&e)
	}

	return (*ogen.GetUpcomingEpisodesOKApplicationJSON)(&result), nil
}

// GetTVShow returns details about a specific TV show.
func (h *Handler) GetTVShow(ctx context.Context, params ogen.GetTVShowParams) (ogen.GetTVShowRes, error) {
	seriesID, err := uuid.Parse(params.ID.String())
	if err != nil {
		return &ogen.GetTVShowNotFound{}, nil
	}

	series, err := h.tvshowService.GetSeries(ctx, seriesID)
	if err != nil {
		return &ogen.GetTVShowNotFound{}, nil
	}

	lang := h.GetMetadataLanguage(ctx)
	localized := LocalizeSeries(series, lang)

	return seriesToOgen(localized), nil
}

// GetTVShowSeasons returns all seasons for a TV show.
func (h *Handler) GetTVShowSeasons(ctx context.Context, params ogen.GetTVShowSeasonsParams) (ogen.GetTVShowSeasonsRes, error) {
	seriesID, err := uuid.Parse(params.ID.String())
	if err != nil {
		return &ogen.GetTVShowSeasonsNotFound{}, nil
	}

	seasons, err := h.tvshowService.ListSeasons(ctx, seriesID)
	if err != nil {
		return &ogen.GetTVShowSeasonsNotFound{}, nil
	}

	result := make([]ogen.TVSeason, len(seasons))
	for i, s := range seasons {
		result[i] = *seasonToOgen(&s)
	}

	return (*ogen.GetTVShowSeasonsOKApplicationJSON)(&result), nil
}

// GetTVShowEpisodes returns all episodes for a TV show.
func (h *Handler) GetTVShowEpisodes(ctx context.Context, params ogen.GetTVShowEpisodesParams) (ogen.GetTVShowEpisodesRes, error) {
	seriesID, err := uuid.Parse(params.ID.String())
	if err != nil {
		return &ogen.GetTVShowEpisodesNotFound{}, nil
	}

	episodes, err := h.tvshowService.ListEpisodesBySeries(ctx, seriesID)
	if err != nil {
		return &ogen.GetTVShowEpisodesNotFound{}, nil
	}

	result := make([]ogen.TVEpisode, len(episodes))
	for i, e := range episodes {
		result[i] = *episodeToOgen(&e)
	}

	return (*ogen.GetTVShowEpisodesOKApplicationJSON)(&result), nil
}

// GetTVShowCast returns the cast for a TV show.
func (h *Handler) GetTVShowCast(ctx context.Context, params ogen.GetTVShowCastParams) (ogen.GetTVShowCastRes, error) {
	seriesID := params.ID
	limit := util.SafeIntToInt32(params.Limit.Or(50))
	offset := util.SafeIntToInt32(params.Offset.Or(0))

	cast, total, err := h.tvshowService.GetSeriesCast(ctx, seriesID, limit, offset)
	if err != nil {
		return &ogen.GetTVShowCastNotFound{}, nil
	}

	items := make([]ogen.TVSeriesCredit, len(cast))
	for i, c := range cast {
		items[i] = *seriesCreditToOgen(&c)
	}

	page := (offset / limit) + 1
	return &ogen.TVSeriesCreditListResponse{
		Items:    items,
		Total:    total,
		Page:     ogen.NewOptInt(int(page)),
		PageSize: ogen.NewOptInt(int(limit)),
	}, nil
}

// GetTVShowCrew returns the crew for a TV show.
func (h *Handler) GetTVShowCrew(ctx context.Context, params ogen.GetTVShowCrewParams) (ogen.GetTVShowCrewRes, error) {
	seriesID := params.ID
	limit := util.SafeIntToInt32(params.Limit.Or(50))
	offset := util.SafeIntToInt32(params.Offset.Or(0))

	crew, total, err := h.tvshowService.GetSeriesCrew(ctx, seriesID, limit, offset)
	if err != nil {
		return &ogen.GetTVShowCrewNotFound{}, nil
	}

	items := make([]ogen.TVSeriesCredit, len(crew))
	for i, c := range crew {
		items[i] = *seriesCreditToOgen(&c)
	}

	page := (offset / limit) + 1
	return &ogen.TVSeriesCreditListResponse{
		Items:    items,
		Total:    total,
		Page:     ogen.NewOptInt(int(page)),
		PageSize: ogen.NewOptInt(int(limit)),
	}, nil
}

// GetTVShowGenres returns the genres for a TV show.
func (h *Handler) GetTVShowGenres(ctx context.Context, params ogen.GetTVShowGenresParams) (ogen.GetTVShowGenresRes, error) {
	seriesID, err := uuid.Parse(params.ID.String())
	if err != nil {
		return &ogen.GetTVShowGenresNotFound{}, nil
	}

	genres, err := h.tvshowService.GetSeriesGenres(ctx, seriesID)
	if err != nil {
		return &ogen.GetTVShowGenresNotFound{}, nil
	}

	result := make([]ogen.TVGenre, len(genres))
	for i, g := range genres {
		result[i] = *seriesGenreToOgen(&g)
	}

	return (*ogen.GetTVShowGenresOKApplicationJSON)(&result), nil
}

// GetTVShowNetworks returns the networks for a TV show.
func (h *Handler) GetTVShowNetworks(ctx context.Context, params ogen.GetTVShowNetworksParams) (ogen.GetTVShowNetworksRes, error) {
	seriesID, err := uuid.Parse(params.ID.String())
	if err != nil {
		return &ogen.GetTVShowNetworksNotFound{}, nil
	}

	networks, err := h.tvshowService.GetSeriesNetworks(ctx, seriesID)
	if err != nil {
		return &ogen.GetTVShowNetworksNotFound{}, nil
	}

	result := make([]ogen.TVNetwork, len(networks))
	for i, n := range networks {
		result[i] = *networkToOgen(&n)
	}

	return (*ogen.GetTVShowNetworksOKApplicationJSON)(&result), nil
}

// GetTVShowWatchStats returns watch statistics for a TV show.
func (h *Handler) GetTVShowWatchStats(ctx context.Context, params ogen.GetTVShowWatchStatsParams) (ogen.GetTVShowWatchStatsRes, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	seriesID, err := uuid.Parse(params.ID.String())
	if err != nil {
		return &ogen.GetTVShowWatchStatsNotFound{}, nil
	}

	stats, err := h.tvshowService.GetSeriesWatchStats(ctx, userID, seriesID)
	if err != nil {
		return &ogen.GetTVShowWatchStatsNotFound{}, nil
	}

	return seriesWatchStatsToOgen(stats, seriesID.String()), nil
}

// GetTVShowNextEpisode returns the next episode to watch for a TV show.
func (h *Handler) GetTVShowNextEpisode(ctx context.Context, params ogen.GetTVShowNextEpisodeParams) (ogen.GetTVShowNextEpisodeRes, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	seriesID, err := uuid.Parse(params.ID.String())
	if err != nil {
		return &ogen.GetTVShowNextEpisodeNotFound{}, nil
	}

	episode, err := h.tvshowService.GetNextEpisode(ctx, userID, seriesID)
	if err != nil || episode == nil {
		return &ogen.GetTVShowNextEpisodeNotFound{}, nil
	}

	return episodeToOgen(episode), nil
}

// RefreshTVShowMetadata triggers a metadata refresh for a TV show.
func (h *Handler) RefreshTVShowMetadata(ctx context.Context, params ogen.RefreshTVShowMetadataParams) (ogen.RefreshTVShowMetadataRes, error) {
	seriesID, err := uuid.Parse(params.ID.String())
	if err != nil {
		return &ogen.RefreshTVShowMetadataNotFound{}, nil
	}

	err = h.tvshowService.RefreshSeriesMetadata(ctx, seriesID)
	if err != nil {
		return &ogen.RefreshTVShowMetadataNotFound{}, nil
	}

	return &ogen.RefreshTVShowMetadataAccepted{}, nil
}

// GetTVSeason returns details about a specific season.
func (h *Handler) GetTVSeason(ctx context.Context, params ogen.GetTVSeasonParams) (ogen.GetTVSeasonRes, error) {
	seasonID, err := uuid.Parse(params.ID.String())
	if err != nil {
		return &ogen.GetTVSeasonNotFound{}, nil
	}

	season, err := h.tvshowService.GetSeason(ctx, seasonID)
	if err != nil {
		return &ogen.GetTVSeasonNotFound{}, nil
	}

	return seasonToOgen(season), nil
}

// GetTVSeasonEpisodes returns all episodes for a season.
func (h *Handler) GetTVSeasonEpisodes(ctx context.Context, params ogen.GetTVSeasonEpisodesParams) (ogen.GetTVSeasonEpisodesRes, error) {
	seasonID, err := uuid.Parse(params.ID.String())
	if err != nil {
		return &ogen.GetTVSeasonEpisodesNotFound{}, nil
	}

	episodes, err := h.tvshowService.ListEpisodesBySeason(ctx, seasonID)
	if err != nil {
		return &ogen.GetTVSeasonEpisodesNotFound{}, nil
	}

	result := make([]ogen.TVEpisode, len(episodes))
	for i, e := range episodes {
		result[i] = *episodeToOgen(&e)
	}

	return (*ogen.GetTVSeasonEpisodesOKApplicationJSON)(&result), nil
}

// GetTVEpisode returns details about a specific episode.
func (h *Handler) GetTVEpisode(ctx context.Context, params ogen.GetTVEpisodeParams) (ogen.GetTVEpisodeRes, error) {
	episodeID, err := uuid.Parse(params.ID.String())
	if err != nil {
		return &ogen.GetTVEpisodeNotFound{}, nil
	}

	episode, err := h.tvshowService.GetEpisode(ctx, episodeID)
	if err != nil {
		return &ogen.GetTVEpisodeNotFound{}, nil
	}

	return episodeToOgen(episode), nil
}

// GetTVEpisodeFiles returns files for an episode.
func (h *Handler) GetTVEpisodeFiles(ctx context.Context, params ogen.GetTVEpisodeFilesParams) (ogen.GetTVEpisodeFilesRes, error) {
	episodeID, err := uuid.Parse(params.ID.String())
	if err != nil {
		return &ogen.GetTVEpisodeFilesNotFound{}, nil
	}

	files, err := h.tvshowService.ListEpisodeFiles(ctx, episodeID)
	if err != nil {
		return &ogen.GetTVEpisodeFilesNotFound{}, nil
	}

	result := make([]ogen.TVEpisodeFile, len(files))
	for i, f := range files {
		result[i] = *episodeFileToOgen(&f)
	}

	return (*ogen.GetTVEpisodeFilesOKApplicationJSON)(&result), nil
}

// GetTVEpisodeProgress returns watch progress for an episode.
func (h *Handler) GetTVEpisodeProgress(ctx context.Context, params ogen.GetTVEpisodeProgressParams) (ogen.GetTVEpisodeProgressRes, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	episodeID, err := uuid.Parse(params.ID.String())
	if err != nil {
		return &ogen.GetTVEpisodeProgressNotFound{}, nil
	}

	progress, err := h.tvshowService.GetEpisodeProgress(ctx, userID, episodeID)
	if err != nil {
		return &ogen.GetTVEpisodeProgressNotFound{}, nil
	}

	return episodeWatchProgressToOgen(progress), nil
}

// UpdateTVEpisodeProgress updates watch progress for an episode.
func (h *Handler) UpdateTVEpisodeProgress(ctx context.Context, req *ogen.UpdateEpisodeProgressRequest, params ogen.UpdateTVEpisodeProgressParams) (ogen.UpdateTVEpisodeProgressRes, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	episodeID, err := uuid.Parse(params.ID.String())
	if err != nil {
		return &ogen.UpdateTVEpisodeProgressNotFound{}, nil
	}

	progress, err := h.tvshowService.UpdateEpisodeProgress(ctx, userID, episodeID, util.SafeIntToInt32(req.ProgressSeconds), util.SafeIntToInt32(req.DurationSeconds))
	if err != nil {
		return &ogen.UpdateTVEpisodeProgressNotFound{}, nil
	}

	return episodeWatchProgressToOgen(progress), nil
}

// DeleteTVEpisodeProgress deletes watch progress for an episode.
func (h *Handler) DeleteTVEpisodeProgress(ctx context.Context, params ogen.DeleteTVEpisodeProgressParams) (ogen.DeleteTVEpisodeProgressRes, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	episodeID, err := uuid.Parse(params.ID.String())
	if err != nil {
		return &ogen.DeleteTVEpisodeProgressNotFound{}, nil
	}

	err = h.tvshowService.RemoveEpisodeProgress(ctx, userID, episodeID)
	if err != nil {
		return &ogen.DeleteTVEpisodeProgressNotFound{}, nil
	}

	return &ogen.DeleteTVEpisodeProgressNoContent{}, nil
}

// MarkTVEpisodeWatched marks an episode as watched.
func (h *Handler) MarkTVEpisodeWatched(ctx context.Context, req ogen.OptMarkTVEpisodeWatchedReq, params ogen.MarkTVEpisodeWatchedParams) (ogen.MarkTVEpisodeWatchedRes, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	episodeID, err := uuid.Parse(params.ID.String())
	if err != nil {
		return &ogen.MarkTVEpisodeWatchedNotFound{}, nil
	}

	err = h.tvshowService.MarkEpisodeWatched(ctx, userID, episodeID)
	if err != nil {
		return &ogen.MarkTVEpisodeWatchedNotFound{}, nil
	}

	return &ogen.MarkTVEpisodeWatchedNoContent{}, nil
}

// MarkTVEpisodesBulkWatched marks multiple episodes as watched in a single request.
func (h *Handler) MarkTVEpisodesBulkWatched(ctx context.Context, req *ogen.BulkEpisodesWatchedRequest) (ogen.MarkTVEpisodesBulkWatchedRes, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	affected, err := h.tvshowService.MarkEpisodesWatchedBulk(ctx, userID, req.EpisodeIds)
	if err != nil {
		return &ogen.Error{
			Code:    500,
			Message: "failed to mark episodes as watched",
		}, nil
	}

	return &ogen.BulkEpisodesWatchedResponse{
		MarkedCount: affected,
	}, nil
}
