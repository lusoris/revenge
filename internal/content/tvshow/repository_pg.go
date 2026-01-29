package tvshow

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	db "github.com/lusoris/revenge/internal/content/tvshow/db"
)

// pgRepository implements Repository using PostgreSQL.
type pgRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

// NewRepository creates a new PostgreSQL-backed TV show repository.
func NewRepository(pool *pgxpool.Pool) Repository {
	return &pgRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

// Series CRUD

// GetSeriesByID retrieves a series by its ID.
func (r *pgRepository) GetSeriesByID(ctx context.Context, id uuid.UUID) (*Series, error) {
	s, err := r.queries.GetSeriesByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSeriesNotFound
		}
		return nil, err
	}
	return FromDBSeries(&s), nil
}

// GetSeriesByTmdbID retrieves a series by its TMDb ID.
func (r *pgRepository) GetSeriesByTmdbID(ctx context.Context, tmdbID int) (*Series, error) {
	id := int32(tmdbID)
	s, err := r.queries.GetSeriesByTmdbID(ctx, &id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSeriesNotFound
		}
		return nil, err
	}
	return FromDBSeries(&s), nil
}

// GetSeriesByImdbID retrieves a series by its IMDb ID.
func (r *pgRepository) GetSeriesByImdbID(ctx context.Context, imdbID string) (*Series, error) {
	s, err := r.queries.GetSeriesByImdbID(ctx, &imdbID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSeriesNotFound
		}
		return nil, err
	}
	return FromDBSeries(&s), nil
}

// GetSeriesByTvdbID retrieves a series by its TVDB ID.
func (r *pgRepository) GetSeriesByTvdbID(ctx context.Context, tvdbID int) (*Series, error) {
	id := int32(tvdbID)
	s, err := r.queries.GetSeriesByTvdbID(ctx, &id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSeriesNotFound
		}
		return nil, err
	}
	return FromDBSeries(&s), nil
}

// ListSeries retrieves series with pagination.
func (r *pgRepository) ListSeries(ctx context.Context, params ListParams) ([]*Series, error) {
	rows, err := r.queries.ListSeries(ctx, db.ListSeriesParams{
		Limit:     int32(params.Limit),
		Offset:    int32(params.Offset),
		SortBy:    params.SortBy,
		SortOrder: params.SortOrder,
	})
	if err != nil {
		return nil, err
	}

	series := make([]*Series, len(rows))
	for i, row := range rows {
		series[i] = FromDBSeries(&row)
	}
	return series, nil
}

// ListSeriesByLibrary retrieves series from a specific library.
func (r *pgRepository) ListSeriesByLibrary(ctx context.Context, libraryID uuid.UUID, params ListParams) ([]*Series, error) {
	rows, err := r.queries.ListSeriesByLibrary(ctx, db.ListSeriesByLibraryParams{
		TvLibraryID: libraryID,
		Limit:       int32(params.Limit),
		Offset:      int32(params.Offset),
		SortBy:      params.SortBy,
		SortOrder:   params.SortOrder,
	})
	if err != nil {
		return nil, err
	}

	series := make([]*Series, len(rows))
	for i, row := range rows {
		series[i] = FromDBSeries(&row)
	}
	return series, nil
}

// ListRecentlyAddedSeries retrieves recently added series.
func (r *pgRepository) ListRecentlyAddedSeries(ctx context.Context, libraryIDs []uuid.UUID, limit int) ([]*Series, error) {
	rows, err := r.queries.ListRecentlyAddedSeries(ctx, db.ListRecentlyAddedSeriesParams{
		LibraryIds: libraryIDs,
		Limit:      int32(limit),
	})
	if err != nil {
		return nil, err
	}

	series := make([]*Series, len(rows))
	for i, row := range rows {
		series[i] = FromDBSeries(&row)
	}
	return series, nil
}

// ListRecentlyPlayedSeries retrieves recently played series.
func (r *pgRepository) ListRecentlyPlayedSeries(ctx context.Context, libraryIDs []uuid.UUID, limit int) ([]*Series, error) {
	rows, err := r.queries.ListRecentlyPlayedSeries(ctx, db.ListRecentlyPlayedSeriesParams{
		LibraryIds: libraryIDs,
		Limit:      int32(limit),
	})
	if err != nil {
		return nil, err
	}

	series := make([]*Series, len(rows))
	for i, row := range rows {
		series[i] = FromDBSeries(&row)
	}
	return series, nil
}

// ListCurrentlyAiringSeries retrieves currently airing series.
func (r *pgRepository) ListCurrentlyAiringSeries(ctx context.Context, libraryIDs []uuid.UUID, limit int) ([]*Series, error) {
	rows, err := r.queries.ListCurrentlyAiringSeries(ctx, db.ListCurrentlyAiringSeriesParams{
		LibraryIds: libraryIDs,
		Limit:      int32(limit),
	})
	if err != nil {
		return nil, err
	}

	series := make([]*Series, len(rows))
	for i, row := range rows {
		series[i] = FromDBSeries(&row)
	}
	return series, nil
}

// SearchSeries searches series by title and overview.
func (r *pgRepository) SearchSeries(ctx context.Context, query string, params ListParams) ([]*Series, error) {
	rows, err := r.queries.SearchSeries(ctx, db.SearchSeriesParams{
		PlaintoTsquery: query,
		Limit:          int32(params.Limit),
		Offset:         int32(params.Offset),
	})
	if err != nil {
		return nil, err
	}

	series := make([]*Series, len(rows))
	for i, row := range rows {
		series[i] = FromDBSeries(&row)
	}
	return series, nil
}

// CountSeries returns the total number of series.
func (r *pgRepository) CountSeries(ctx context.Context) (int64, error) {
	return r.queries.CountSeries(ctx)
}

// CountSeriesByLibrary returns the number of series in a library.
func (r *pgRepository) CountSeriesByLibrary(ctx context.Context, libraryID uuid.UUID) (int64, error) {
	return r.queries.CountSeriesByLibrary(ctx, libraryID)
}

// CreateSeries creates a new series.
func (r *pgRepository) CreateSeries(ctx context.Context, series *Series) error {
	s, err := r.queries.CreateSeries(ctx, series.ToDBCreateParams())
	if err != nil {
		return err
	}
	series.ID = s.ID
	series.CreatedAt = s.CreatedAt
	series.UpdatedAt = s.UpdatedAt
	series.DateAdded = s.DateAdded
	return nil
}

// UpdateSeries updates an existing series.
func (r *pgRepository) UpdateSeries(ctx context.Context, series *Series) error {
	params := db.UpdateSeriesParams{
		ID: series.ID,
	}

	if series.Title != "" {
		params.Title = &series.Title
	}
	if series.SortTitle != "" {
		params.SortTitle = &series.SortTitle
	}
	if series.OriginalTitle != "" {
		params.OriginalTitle = &series.OriginalTitle
	}
	if series.Tagline != "" {
		params.Tagline = &series.Tagline
	}
	if series.Overview != "" {
		params.Overview = &series.Overview
	}
	if series.FirstAirDate != nil {
		params.FirstAirDate = pgtype.Date{Time: *series.FirstAirDate, Valid: true}
	}
	if series.LastAirDate != nil {
		params.LastAirDate = pgtype.Date{Time: *series.LastAirDate, Valid: true}
	}
	if series.Year > 0 {
		y := int32(series.Year)
		params.Year = &y
	}
	if series.Status != "" {
		params.Status = &series.Status
	}
	if series.Type != "" {
		params.Type = &series.Type
	}
	if series.ContentRating != "" {
		params.ContentRating = &series.ContentRating
	}
	if series.RatingLevel >= 0 {
		rl := int32(series.RatingLevel)
		params.RatingLevel = &rl
	}
	if series.CommunityRating > 0 {
		params.CommunityRating = numericFromFloat(series.CommunityRating)
	}
	if series.VoteCount > 0 {
		vc := int32(series.VoteCount)
		params.VoteCount = &vc
	}
	if series.PosterPath != "" {
		params.PosterPath = &series.PosterPath
	}
	if series.PosterBlurhash != "" {
		params.PosterBlurhash = &series.PosterBlurhash
	}
	if series.BackdropPath != "" {
		params.BackdropPath = &series.BackdropPath
	}
	if series.BackdropBlurhash != "" {
		params.BackdropBlurhash = &series.BackdropBlurhash
	}
	if series.LogoPath != "" {
		params.LogoPath = &series.LogoPath
	}
	if series.TmdbID > 0 {
		id := int32(series.TmdbID)
		params.TmdbID = &id
	}
	if series.ImdbID != "" {
		params.ImdbID = &series.ImdbID
	}
	if series.TvdbID > 0 {
		id := int32(series.TvdbID)
		params.TvdbID = &id
	}
	if series.NetworkName != "" {
		params.NetworkName = &series.NetworkName
	}
	if series.NetworkLogoPath != "" {
		params.NetworkLogoPath = &series.NetworkLogoPath
	}
	params.IsLocked = &series.IsLocked

	s, err := r.queries.UpdateSeries(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrSeriesNotFound
		}
		return err
	}
	series.UpdatedAt = s.UpdatedAt
	return nil
}

// UpdateSeriesPlaybackStats updates the playback statistics for a series.
func (r *pgRepository) UpdateSeriesPlaybackStats(ctx context.Context, id uuid.UUID) error {
	return r.queries.UpdateSeriesPlaybackStats(ctx, id)
}

// DeleteSeries deletes a series.
func (r *pgRepository) DeleteSeries(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteSeries(ctx, id)
}

// DeleteSeriesByLibrary deletes all series in a library.
func (r *pgRepository) DeleteSeriesByLibrary(ctx context.Context, libraryID uuid.UUID) error {
	return r.queries.DeleteSeriesByLibrary(ctx, libraryID)
}

// SeriesExistsByTmdbID checks if a series exists by TMDb ID.
func (r *pgRepository) SeriesExistsByTmdbID(ctx context.Context, tmdbID int) (bool, error) {
	id := int32(tmdbID)
	return r.queries.SeriesExistsByTmdbID(ctx, &id)
}

// SeriesExistsByTvdbID checks if a series exists by TVDB ID.
func (r *pgRepository) SeriesExistsByTvdbID(ctx context.Context, tvdbID int) (bool, error) {
	id := int32(tvdbID)
	return r.queries.SeriesExistsByTvdbID(ctx, &id)
}

// Seasons

// GetSeasonByID retrieves a season by its ID.
func (r *pgRepository) GetSeasonByID(ctx context.Context, id uuid.UUID) (*Season, error) {
	s, err := r.queries.GetSeasonByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSeasonNotFound
		}
		return nil, err
	}
	return FromDBSeason(&s), nil
}

// GetSeasonByNumber retrieves a season by series ID and season number.
func (r *pgRepository) GetSeasonByNumber(ctx context.Context, seriesID uuid.UUID, seasonNumber int) (*Season, error) {
	s, err := r.queries.GetSeasonByNumber(ctx, db.GetSeasonByNumberParams{
		SeriesID:     seriesID,
		SeasonNumber: int32(seasonNumber),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSeasonNotFound
		}
		return nil, err
	}
	return FromDBSeason(&s), nil
}

// ListSeasons retrieves all seasons for a series (excluding specials).
func (r *pgRepository) ListSeasons(ctx context.Context, seriesID uuid.UUID) ([]*Season, error) {
	rows, err := r.queries.ListSeasons(ctx, seriesID)
	if err != nil {
		return nil, err
	}

	seasons := make([]*Season, len(rows))
	for i, row := range rows {
		seasons[i] = FromDBSeason(&row)
	}
	return seasons, nil
}

// ListSeasonsWithSpecials retrieves all seasons for a series including specials.
func (r *pgRepository) ListSeasonsWithSpecials(ctx context.Context, seriesID uuid.UUID) ([]*Season, error) {
	rows, err := r.queries.ListSeasonsWithSpecials(ctx, seriesID)
	if err != nil {
		return nil, err
	}

	seasons := make([]*Season, len(rows))
	for i, row := range rows {
		seasons[i] = FromDBSeason(&row)
	}
	return seasons, nil
}

// CountSeasons returns the number of seasons for a series.
func (r *pgRepository) CountSeasons(ctx context.Context, seriesID uuid.UUID) (int64, error) {
	return r.queries.CountSeasons(ctx, seriesID)
}

// CreateSeason creates a new season.
func (r *pgRepository) CreateSeason(ctx context.Context, season *Season) error {
	params := db.CreateSeasonParams{
		SeriesID:     season.SeriesID,
		SeasonNumber: int32(season.SeasonNumber),
	}

	if season.Name != "" {
		params.Name = &season.Name
	}
	if season.Overview != "" {
		params.Overview = &season.Overview
	}
	if season.AirDate != nil {
		params.AirDate = pgtype.Date{Time: *season.AirDate, Valid: true}
	}
	if season.Year > 0 {
		y := int32(season.Year)
		params.Year = &y
	}
	if season.PosterPath != "" {
		params.PosterPath = &season.PosterPath
	}
	if season.PosterBlurhash != "" {
		params.PosterBlurhash = &season.PosterBlurhash
	}
	if season.TmdbID > 0 {
		id := int32(season.TmdbID)
		params.TmdbID = &id
	}
	if season.TvdbID > 0 {
		id := int32(season.TvdbID)
		params.TvdbID = &id
	}

	s, err := r.queries.CreateSeason(ctx, params)
	if err != nil {
		return err
	}
	season.ID = s.ID
	season.CreatedAt = s.CreatedAt
	season.UpdatedAt = s.UpdatedAt
	return nil
}

// UpdateSeason updates an existing season.
func (r *pgRepository) UpdateSeason(ctx context.Context, season *Season) error {
	params := db.UpdateSeasonParams{
		ID: season.ID,
	}

	if season.Name != "" {
		params.Name = &season.Name
	}
	if season.Overview != "" {
		params.Overview = &season.Overview
	}
	if season.AirDate != nil {
		params.AirDate = pgtype.Date{Time: *season.AirDate, Valid: true}
	}
	if season.Year > 0 {
		y := int32(season.Year)
		params.Year = &y
	}
	if season.PosterPath != "" {
		params.PosterPath = &season.PosterPath
	}
	if season.PosterBlurhash != "" {
		params.PosterBlurhash = &season.PosterBlurhash
	}
	if season.TmdbID > 0 {
		id := int32(season.TmdbID)
		params.TmdbID = &id
	}
	if season.TvdbID > 0 {
		id := int32(season.TvdbID)
		params.TvdbID = &id
	}

	s, err := r.queries.UpdateSeason(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrSeasonNotFound
		}
		return err
	}
	season.UpdatedAt = s.UpdatedAt
	return nil
}

// DeleteSeason deletes a season.
func (r *pgRepository) DeleteSeason(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteSeason(ctx, id)
}

// DeleteSeasonsBySeries deletes all seasons for a series.
func (r *pgRepository) DeleteSeasonsBySeries(ctx context.Context, seriesID uuid.UUID) error {
	return r.queries.DeleteSeasonsBySeries(ctx, seriesID)
}

// GetOrCreateSeason gets or creates a season.
func (r *pgRepository) GetOrCreateSeason(ctx context.Context, seriesID uuid.UUID, seasonNumber int, name string) (*Season, error) {
	var n interface{} = nil
	if name != "" {
		n = name
	}
	s, err := r.queries.GetOrCreateSeason(ctx, db.GetOrCreateSeasonParams{
		SeriesID:     seriesID,
		SeasonNumber: int32(seasonNumber),
		Column3:      n,
	})
	if err != nil {
		return nil, err
	}
	return FromDBSeason(&s), nil
}

// Episodes

// GetEpisodeByID retrieves an episode by its ID.
func (r *pgRepository) GetEpisodeByID(ctx context.Context, id uuid.UUID) (*Episode, error) {
	e, err := r.queries.GetEpisodeByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrEpisodeNotFound
		}
		return nil, err
	}
	return FromDBEpisode(&e), nil
}

// GetEpisodeByPath retrieves an episode by its file path.
func (r *pgRepository) GetEpisodeByPath(ctx context.Context, path string) (*Episode, error) {
	e, err := r.queries.GetEpisodeByPath(ctx, path)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrEpisodeNotFound
		}
		return nil, err
	}
	return FromDBEpisode(&e), nil
}

// GetEpisodeByNumber retrieves an episode by series, season, and episode number.
func (r *pgRepository) GetEpisodeByNumber(ctx context.Context, seriesID uuid.UUID, seasonNumber, episodeNumber int) (*Episode, error) {
	e, err := r.queries.GetEpisodeByNumber(ctx, db.GetEpisodeByNumberParams{
		SeriesID:      seriesID,
		SeasonNumber:  int32(seasonNumber),
		EpisodeNumber: int32(episodeNumber),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrEpisodeNotFound
		}
		return nil, err
	}
	return FromDBEpisode(&e), nil
}

// GetEpisodeByAbsoluteNumber retrieves an episode by absolute number (for anime).
func (r *pgRepository) GetEpisodeByAbsoluteNumber(ctx context.Context, seriesID uuid.UUID, absoluteNumber int) (*Episode, error) {
	num := int32(absoluteNumber)
	e, err := r.queries.GetEpisodeByAbsoluteNumber(ctx, db.GetEpisodeByAbsoluteNumberParams{
		SeriesID:       seriesID,
		AbsoluteNumber: &num,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrEpisodeNotFound
		}
		return nil, err
	}
	return FromDBEpisode(&e), nil
}

// ListEpisodes retrieves all episodes for a series.
func (r *pgRepository) ListEpisodes(ctx context.Context, seriesID uuid.UUID) ([]*Episode, error) {
	rows, err := r.queries.ListEpisodes(ctx, seriesID)
	if err != nil {
		return nil, err
	}

	episodes := make([]*Episode, len(rows))
	for i, row := range rows {
		episodes[i] = FromDBEpisode(&row)
	}
	return episodes, nil
}

// ListEpisodesBySeason retrieves all episodes for a season.
func (r *pgRepository) ListEpisodesBySeason(ctx context.Context, seasonID uuid.UUID) ([]*Episode, error) {
	rows, err := r.queries.ListEpisodesBySeason(ctx, seasonID)
	if err != nil {
		return nil, err
	}

	episodes := make([]*Episode, len(rows))
	for i, row := range rows {
		episodes[i] = FromDBEpisode(&row)
	}
	return episodes, nil
}

// ListEpisodesBySeasonNumber retrieves episodes by series ID and season number.
func (r *pgRepository) ListEpisodesBySeasonNumber(ctx context.Context, seriesID uuid.UUID, seasonNumber int) ([]*Episode, error) {
	rows, err := r.queries.ListEpisodesBySeasonNumber(ctx, db.ListEpisodesBySeasonNumberParams{
		SeriesID:     seriesID,
		SeasonNumber: int32(seasonNumber),
	})
	if err != nil {
		return nil, err
	}

	episodes := make([]*Episode, len(rows))
	for i, row := range rows {
		episodes[i] = FromDBEpisode(&row)
	}
	return episodes, nil
}

// ListRecentlyAddedEpisodes retrieves recently added episodes.
func (r *pgRepository) ListRecentlyAddedEpisodes(ctx context.Context, libraryIDs []uuid.UUID, limit int) ([]*Episode, error) {
	rows, err := r.queries.ListRecentlyAddedEpisodes(ctx, db.ListRecentlyAddedEpisodesParams{
		LibraryIds: libraryIDs,
		Limit:      int32(limit),
	})
	if err != nil {
		return nil, err
	}

	episodes := make([]*Episode, len(rows))
	for i, row := range rows {
		ep := FromDBEpisode(&db.Episode{
			ID:            row.ID,
			SeriesID:      row.SeriesID,
			SeasonID:      row.SeasonID,
			Path:          row.Path,
			Container:     row.Container,
			SizeBytes:     row.SizeBytes,
			RuntimeTicks:  row.RuntimeTicks,
			SeasonNumber:  row.SeasonNumber,
			EpisodeNumber: row.EpisodeNumber,
			Title:         row.Title,
			DateAdded:     row.DateAdded,
			CreatedAt:     row.CreatedAt,
			UpdatedAt:     row.UpdatedAt,
		})
		episodes[i] = ep
	}
	return episodes, nil
}

// ListRecentlyAiredEpisodes retrieves recently aired episodes.
func (r *pgRepository) ListRecentlyAiredEpisodes(ctx context.Context, libraryIDs []uuid.UUID, limit int) ([]*Episode, error) {
	rows, err := r.queries.ListRecentlyAiredEpisodes(ctx, db.ListRecentlyAiredEpisodesParams{
		LibraryIds: libraryIDs,
		Limit:      int32(limit),
	})
	if err != nil {
		return nil, err
	}

	episodes := make([]*Episode, len(rows))
	for i, row := range rows {
		ep := FromDBEpisode(&db.Episode{
			ID:            row.ID,
			SeriesID:      row.SeriesID,
			SeasonID:      row.SeasonID,
			Path:          row.Path,
			Container:     row.Container,
			SizeBytes:     row.SizeBytes,
			RuntimeTicks:  row.RuntimeTicks,
			SeasonNumber:  row.SeasonNumber,
			EpisodeNumber: row.EpisodeNumber,
			Title:         row.Title,
			AirDate:       row.AirDate,
			DateAdded:     row.DateAdded,
			CreatedAt:     row.CreatedAt,
			UpdatedAt:     row.UpdatedAt,
		})
		episodes[i] = ep
	}
	return episodes, nil
}

// ListUpcomingEpisodes retrieves upcoming episodes.
func (r *pgRepository) ListUpcomingEpisodes(ctx context.Context, libraryIDs []uuid.UUID, limit int) ([]*Episode, error) {
	rows, err := r.queries.ListUpcomingEpisodes(ctx, db.ListUpcomingEpisodesParams{
		LibraryIds: libraryIDs,
		Limit:      int32(limit),
	})
	if err != nil {
		return nil, err
	}

	episodes := make([]*Episode, len(rows))
	for i, row := range rows {
		ep := FromDBEpisode(&db.Episode{
			ID:            row.ID,
			SeriesID:      row.SeriesID,
			SeasonID:      row.SeasonID,
			Path:          row.Path,
			Container:     row.Container,
			SizeBytes:     row.SizeBytes,
			RuntimeTicks:  row.RuntimeTicks,
			SeasonNumber:  row.SeasonNumber,
			EpisodeNumber: row.EpisodeNumber,
			Title:         row.Title,
			AirDate:       row.AirDate,
			DateAdded:     row.DateAdded,
			CreatedAt:     row.CreatedAt,
			UpdatedAt:     row.UpdatedAt,
		})
		episodes[i] = ep
	}
	return episodes, nil
}

// CountEpisodes returns the number of episodes for a series.
func (r *pgRepository) CountEpisodes(ctx context.Context, seriesID uuid.UUID) (int64, error) {
	return r.queries.CountEpisodes(ctx, seriesID)
}

// CountEpisodesBySeason returns the number of episodes for a season.
func (r *pgRepository) CountEpisodesBySeason(ctx context.Context, seasonID uuid.UUID) (int64, error) {
	return r.queries.CountEpisodesBySeason(ctx, seasonID)
}

// CreateEpisode creates a new episode.
func (r *pgRepository) CreateEpisode(ctx context.Context, episode *Episode) error {
	params := db.CreateEpisodeParams{
		SeriesID:      episode.SeriesID,
		SeasonID:      episode.SeasonID,
		Path:          episode.Path,
		SeasonNumber:  int32(episode.SeasonNumber),
		EpisodeNumber: int32(episode.EpisodeNumber),
		Title:         episode.Title,
	}

	if episode.Container != "" {
		params.Container = &episode.Container
	}
	if episode.SizeBytes > 0 {
		params.SizeBytes = &episode.SizeBytes
	}
	if episode.RuntimeTicks > 0 {
		params.RuntimeTicks = &episode.RuntimeTicks
	}
	if episode.AbsoluteNumber != nil {
		n := int32(*episode.AbsoluteNumber)
		params.AbsoluteNumber = &n
	}
	if episode.Overview != "" {
		params.Overview = &episode.Overview
	}
	if episode.ProductionCode != "" {
		params.ProductionCode = &episode.ProductionCode
	}
	if episode.AirDate != nil {
		params.AirDate = pgtype.Date{Time: *episode.AirDate, Valid: true}
	}
	if episode.AirDateUTC != nil {
		params.AirDateUtc = pgtype.Timestamptz{Time: *episode.AirDateUTC, Valid: true}
	}
	if episode.CommunityRating > 0 {
		params.CommunityRating = numericFromFloat(episode.CommunityRating)
	}
	if episode.VoteCount > 0 {
		vc := int32(episode.VoteCount)
		params.VoteCount = &vc
	}
	if episode.StillPath != "" {
		params.StillPath = &episode.StillPath
	}
	if episode.StillBlurhash != "" {
		params.StillBlurhash = &episode.StillBlurhash
	}
	if episode.TmdbID > 0 {
		id := int32(episode.TmdbID)
		params.TmdbID = &id
	}
	if episode.ImdbID != "" {
		params.ImdbID = &episode.ImdbID
	}
	if episode.TvdbID > 0 {
		id := int32(episode.TvdbID)
		params.TvdbID = &id
	}

	e, err := r.queries.CreateEpisode(ctx, params)
	if err != nil {
		return err
	}
	episode.ID = e.ID
	episode.CreatedAt = e.CreatedAt
	episode.UpdatedAt = e.UpdatedAt
	episode.DateAdded = e.DateAdded
	return nil
}

// UpdateEpisode updates an existing episode.
func (r *pgRepository) UpdateEpisode(ctx context.Context, episode *Episode) error {
	params := db.UpdateEpisodeParams{
		ID: episode.ID,
	}

	if episode.Title != "" {
		params.Title = &episode.Title
	}
	if episode.SeasonID != uuid.Nil {
		params.SeasonID = pgtype.UUID{Bytes: episode.SeasonID, Valid: true}
	}
	if episode.Container != "" {
		params.Container = &episode.Container
	}
	if episode.SizeBytes > 0 {
		params.SizeBytes = &episode.SizeBytes
	}
	if episode.RuntimeTicks > 0 {
		params.RuntimeTicks = &episode.RuntimeTicks
	}
	if episode.SeasonNumber > 0 {
		sn := int32(episode.SeasonNumber)
		params.SeasonNumber = &sn
	}
	if episode.EpisodeNumber > 0 {
		en := int32(episode.EpisodeNumber)
		params.EpisodeNumber = &en
	}
	if episode.AbsoluteNumber != nil {
		n := int32(*episode.AbsoluteNumber)
		params.AbsoluteNumber = &n
	}
	if episode.Overview != "" {
		params.Overview = &episode.Overview
	}
	if episode.ProductionCode != "" {
		params.ProductionCode = &episode.ProductionCode
	}
	if episode.AirDate != nil {
		params.AirDate = pgtype.Date{Time: *episode.AirDate, Valid: true}
	}
	if episode.AirDateUTC != nil {
		params.AirDateUtc = pgtype.Timestamptz{Time: *episode.AirDateUTC, Valid: true}
	}
	if episode.CommunityRating > 0 {
		params.CommunityRating = numericFromFloat(episode.CommunityRating)
	}
	if episode.VoteCount > 0 {
		vc := int32(episode.VoteCount)
		params.VoteCount = &vc
	}
	if episode.StillPath != "" {
		params.StillPath = &episode.StillPath
	}
	if episode.StillBlurhash != "" {
		params.StillBlurhash = &episode.StillBlurhash
	}
	if episode.TmdbID > 0 {
		id := int32(episode.TmdbID)
		params.TmdbID = &id
	}
	if episode.ImdbID != "" {
		params.ImdbID = &episode.ImdbID
	}
	if episode.TvdbID > 0 {
		id := int32(episode.TvdbID)
		params.TvdbID = &id
	}
	params.IsLocked = &episode.IsLocked

	e, err := r.queries.UpdateEpisode(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrEpisodeNotFound
		}
		return err
	}
	episode.UpdatedAt = e.UpdatedAt
	return nil
}

// UpdateEpisodePlaybackStats updates the playback statistics for an episode.
func (r *pgRepository) UpdateEpisodePlaybackStats(ctx context.Context, id uuid.UUID) error {
	return r.queries.UpdateEpisodePlaybackStats(ctx, id)
}

// DeleteEpisode deletes an episode.
func (r *pgRepository) DeleteEpisode(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteEpisode(ctx, id)
}

// DeleteEpisodesBySeries deletes all episodes for a series.
func (r *pgRepository) DeleteEpisodesBySeries(ctx context.Context, seriesID uuid.UUID) error {
	return r.queries.DeleteEpisodesBySeries(ctx, seriesID)
}

// DeleteEpisodesBySeason deletes all episodes for a season.
func (r *pgRepository) DeleteEpisodesBySeason(ctx context.Context, seasonID uuid.UUID) error {
	return r.queries.DeleteEpisodesBySeason(ctx, seasonID)
}

// EpisodeExistsByPath checks if an episode exists by path.
func (r *pgRepository) EpisodeExistsByPath(ctx context.Context, path string) (bool, error) {
	return r.queries.EpisodeExistsByPath(ctx, path)
}

// EpisodeExistsByNumber checks if an episode exists by series, season, and episode number.
func (r *pgRepository) EpisodeExistsByNumber(ctx context.Context, seriesID uuid.UUID, seasonNumber, episodeNumber int) (bool, error) {
	return r.queries.EpisodeExistsByNumber(ctx, db.EpisodeExistsByNumberParams{
		SeriesID:      seriesID,
		SeasonNumber:  int32(seasonNumber),
		EpisodeNumber: int32(episodeNumber),
	})
}

// ListEpisodePaths returns all episode paths in a library.
func (r *pgRepository) ListEpisodePaths(ctx context.Context, libraryID uuid.UUID) (map[uuid.UUID]string, error) {
	rows, err := r.queries.ListEpisodePaths(ctx, libraryID)
	if err != nil {
		return nil, err
	}

	paths := make(map[uuid.UUID]string, len(rows))
	for _, row := range rows {
		paths[row.ID] = row.Path
	}
	return paths, nil
}

// GetNextEpisode returns the next episode in a series.
func (r *pgRepository) GetNextEpisode(ctx context.Context, seriesID uuid.UUID, seasonNumber, episodeNumber int) (*Episode, error) {
	e, err := r.queries.GetNextEpisode(ctx, db.GetNextEpisodeParams{
		SeriesID:      seriesID,
		SeasonNumber:  int32(seasonNumber),
		EpisodeNumber: int32(episodeNumber),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrEpisodeNotFound
		}
		return nil, err
	}
	return FromDBEpisode(&e), nil
}

// GetPreviousEpisode returns the previous episode in a series.
func (r *pgRepository) GetPreviousEpisode(ctx context.Context, seriesID uuid.UUID, seasonNumber, episodeNumber int) (*Episode, error) {
	e, err := r.queries.GetPreviousEpisode(ctx, db.GetPreviousEpisodeParams{
		SeriesID:      seriesID,
		SeasonNumber:  int32(seasonNumber),
		EpisodeNumber: int32(episodeNumber),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrEpisodeNotFound
		}
		return nil, err
	}
	return FromDBEpisode(&e), nil
}
