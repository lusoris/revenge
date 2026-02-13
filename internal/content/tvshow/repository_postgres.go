package tvshow

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/govalues/decimal"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/lusoris/revenge/internal/content"
	tvshowdb "github.com/lusoris/revenge/internal/content/tvshow/db"
)

// postgresRepository implements the Repository interface using PostgreSQL
type postgresRepository struct {
	pool    *pgxpool.Pool
	queries *tvshowdb.Queries
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(pool *pgxpool.Pool) Repository {
	return &postgresRepository{
		pool:    pool,
		queries: tvshowdb.New(pool),
	}
}

// =============================================================================
// Series Operations
// =============================================================================

func (r *postgresRepository) GetSeries(ctx context.Context, id uuid.UUID) (*Series, error) {
	series, err := r.queries.GetSeries(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("series not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get series: %w", err)
	}
	return dbSeriesToSeries(series), nil
}

func (r *postgresRepository) GetSeriesByTMDbID(ctx context.Context, tmdbID int32) (*Series, error) {
	series, err := r.queries.GetSeriesByTMDbID(ctx, &tmdbID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("series not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get series by TMDb ID: %w", err)
	}
	return dbSeriesToSeries(series), nil
}

func (r *postgresRepository) GetSeriesByTVDbID(ctx context.Context, tvdbID int32) (*Series, error) {
	series, err := r.queries.GetSeriesByTVDbID(ctx, &tvdbID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("series not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get series by TVDb ID: %w", err)
	}
	return dbSeriesToSeries(series), nil
}

func (r *postgresRepository) GetSeriesBySonarrID(ctx context.Context, sonarrID int32) (*Series, error) {
	series, err := r.queries.GetSeriesBySonarrID(ctx, &sonarrID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("series not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get series by Sonarr ID: %w", err)
	}
	return dbSeriesToSeries(series), nil
}

func (r *postgresRepository) ListSeries(ctx context.Context, filters SeriesListFilters) ([]Series, error) {
	dbSeries, err := r.queries.ListSeries(ctx, tvshowdb.ListSeriesParams{
		Limit:  filters.Limit,
		Offset: filters.Offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list series: %w", err)
	}

	result := make([]Series, len(dbSeries))
	for i, s := range dbSeries {
		result[i] = *dbSeriesToSeries(s)
	}
	return result, nil
}

func (r *postgresRepository) CountSeries(ctx context.Context) (int64, error) {
	return r.queries.CountSeries(ctx)
}

func (r *postgresRepository) SearchSeriesByTitle(ctx context.Context, query string, limit, offset int32) ([]Series, error) {
	searchPattern := "%" + query + "%"
	dbSeries, err := r.queries.SearchSeriesByTitle(ctx, tvshowdb.SearchSeriesByTitleParams{
		Column1: &searchPattern,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search series: %w", err)
	}

	result := make([]Series, len(dbSeries))
	for i, s := range dbSeries {
		result[i] = *dbSeriesToSeries(s)
	}
	return result, nil
}

func (r *postgresRepository) SearchSeriesByTitleAnyLanguage(ctx context.Context, query string, limit, offset int32) ([]Series, error) {
	searchPattern := "%" + query + "%"
	dbSeries, err := r.queries.SearchSeriesByTitleAnyLanguage(ctx, tvshowdb.SearchSeriesByTitleAnyLanguageParams{
		Column1: &searchPattern,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search series: %w", err)
	}

	result := make([]Series, len(dbSeries))
	for i, s := range dbSeries {
		result[i] = *dbSeriesToSeries(s)
	}
	return result, nil
}

func (r *postgresRepository) ListRecentlyAddedSeries(ctx context.Context, limit, offset int32) ([]Series, error) {
	dbSeries, err := r.queries.ListRecentlyAddedSeries(ctx, tvshowdb.ListRecentlyAddedSeriesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list recently added series: %w", err)
	}

	result := make([]Series, len(dbSeries))
	for i, s := range dbSeries {
		result[i] = *dbSeriesToSeries(s)
	}
	return result, nil
}

func (r *postgresRepository) ListSeriesByGenre(ctx context.Context, slug string, limit, offset int32) ([]Series, error) {
	dbSeries, err := r.queries.ListSeriesByGenre(ctx, tvshowdb.ListSeriesByGenreParams{
		Slug:   slug,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list series by genre: %w", err)
	}

	result := make([]Series, len(dbSeries))
	for i, s := range dbSeries {
		result[i] = *dbSeriesToSeries(s)
	}
	return result, nil
}

func (r *postgresRepository) ListSeriesByNetwork(ctx context.Context, networkID uuid.UUID, limit, offset int32) ([]Series, error) {
	dbSeries, err := r.queries.ListSeriesByNetwork(ctx, tvshowdb.ListSeriesByNetworkParams{
		NetworkID: networkID,
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list series by network: %w", err)
	}

	result := make([]Series, len(dbSeries))
	for i, s := range dbSeries {
		result[i] = *dbSeriesToSeries(s)
	}
	return result, nil
}

func (r *postgresRepository) ListSeriesByStatus(ctx context.Context, status string, limit, offset int32) ([]Series, error) {
	dbSeries, err := r.queries.ListSeriesByStatus(ctx, tvshowdb.ListSeriesByStatusParams{
		Status: &status,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list series by status: %w", err)
	}

	result := make([]Series, len(dbSeries))
	for i, s := range dbSeries {
		result[i] = *dbSeriesToSeries(s)
	}
	return result, nil
}

func (r *postgresRepository) CreateSeries(ctx context.Context, params CreateSeriesParams) (*Series, error) {
	dbSeries, err := r.queries.CreateSeries(ctx, tvshowdb.CreateSeriesParams{
		TmdbID:            params.TMDbID,
		TvdbID:            params.TVDbID,
		ImdbID:            params.IMDbID,
		SonarrID:          params.SonarrID,
		Title:             params.Title,
		Tagline:           params.Tagline,
		Overview:          params.Overview,
		TitlesI18n:        marshalStringMap(params.TitlesI18n),
		TaglinesI18n:      marshalStringMap(params.TaglinesI18n),
		OverviewsI18n:     marshalStringMap(params.OverviewsI18n),
		AgeRatings:        marshalNestedStringMap(params.AgeRatings),
		ExternalRatings:   marshalExternalRatings(params.ExternalRatings),
		OriginalLanguage:  params.OriginalLanguage,
		OriginalTitle:     params.OriginalTitle,
		Status:            params.Status,
		Type:              params.Type,
		FirstAirDate:      stringToPgDate(params.FirstAirDate),
		LastAirDate:       stringToPgDate(params.LastAirDate),
		VoteAverage:       stringToPgNumeric(params.VoteAverage),
		VoteCount:         params.VoteCount,
		Popularity:        stringToPgNumeric(params.Popularity),
		PosterPath:        params.PosterPath,
		BackdropPath:      params.BackdropPath,
		TotalSeasons:      params.TotalSeasons,
		TotalEpisodes:     params.TotalEpisodes,
		TrailerUrl:        params.TrailerURL,
		Homepage:          params.Homepage,
		MetadataUpdatedAt: stringToPgTimestamptz(params.MetadataUpdatedAt),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create series: %w", err)
	}
	return dbSeriesToSeries(dbSeries), nil
}

func (r *postgresRepository) UpdateSeries(ctx context.Context, params UpdateSeriesParams) (*Series, error) {
	dbSeries, err := r.queries.UpdateSeries(ctx, tvshowdb.UpdateSeriesParams{
		ID:                params.ID,
		TmdbID:            params.TMDbID,
		TvdbID:            params.TVDbID,
		ImdbID:            params.IMDbID,
		SonarrID:          params.SonarrID,
		Title:             params.Title,
		Tagline:           params.Tagline,
		Overview:          params.Overview,
		TitlesI18n:        marshalStringMapToBytes(params.TitlesI18n),
		TaglinesI18n:      marshalStringMapToBytes(params.TaglinesI18n),
		OverviewsI18n:     marshalStringMapToBytes(params.OverviewsI18n),
		AgeRatings:        marshalNestedStringMapToBytes(params.AgeRatings),
		ExternalRatings:   marshalExternalRatingsToBytes(params.ExternalRatings),
		OriginalLanguage:  params.OriginalLanguage,
		OriginalTitle:     params.OriginalTitle,
		Status:            params.Status,
		Type:              params.Type,
		FirstAirDate:      stringToPgDate(params.FirstAirDate),
		LastAirDate:       stringToPgDate(params.LastAirDate),
		VoteAverage:       stringToPgNumeric(params.VoteAverage),
		VoteCount:         params.VoteCount,
		Popularity:        stringToPgNumeric(params.Popularity),
		PosterPath:        params.PosterPath,
		BackdropPath:      params.BackdropPath,
		TotalSeasons:      params.TotalSeasons,
		TotalEpisodes:     params.TotalEpisodes,
		TrailerUrl:        params.TrailerURL,
		Homepage:          params.Homepage,
		MetadataUpdatedAt: stringToPgTimestamptz(params.MetadataUpdatedAt),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update series: %w", err)
	}
	return dbSeriesToSeries(dbSeries), nil
}

func (r *postgresRepository) UpdateSeriesStats(ctx context.Context, seriesID uuid.UUID) error {
	return r.queries.UpdateSeriesStats(ctx, seriesID)
}

func (r *postgresRepository) DeleteSeries(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteSeries(ctx, id)
}

// =============================================================================
// Season Operations
// =============================================================================

func (r *postgresRepository) GetSeason(ctx context.Context, id uuid.UUID) (*Season, error) {
	season, err := r.queries.GetSeason(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("season not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get season: %w", err)
	}
	return dbSeasonToSeason(season), nil
}

func (r *postgresRepository) GetSeasonByNumber(ctx context.Context, seriesID uuid.UUID, seasonNumber int32) (*Season, error) {
	season, err := r.queries.GetSeasonByNumber(ctx, tvshowdb.GetSeasonByNumberParams{
		SeriesID:     seriesID,
		SeasonNumber: seasonNumber,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("season not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get season by number: %w", err)
	}
	return dbSeasonToSeason(season), nil
}

func (r *postgresRepository) ListSeasonsBySeries(ctx context.Context, seriesID uuid.UUID) ([]Season, error) {
	dbSeasons, err := r.queries.ListSeasonsBySeries(ctx, seriesID)
	if err != nil {
		return nil, fmt.Errorf("failed to list seasons: %w", err)
	}

	result := make([]Season, len(dbSeasons))
	for i, s := range dbSeasons {
		result[i] = *dbSeasonToSeason(s)
	}
	return result, nil
}

func (r *postgresRepository) ListSeasonsBySeriesWithEpisodeCount(ctx context.Context, seriesID uuid.UUID) ([]SeasonWithEpisodeCount, error) {
	dbSeasons, err := r.queries.ListSeasonsBySeriesWithEpisodeCount(ctx, seriesID)
	if err != nil {
		return nil, fmt.Errorf("failed to list seasons with episode count: %w", err)
	}

	result := make([]SeasonWithEpisodeCount, len(dbSeasons))
	for i, s := range dbSeasons {
		result[i] = SeasonWithEpisodeCount{
			Season: Season{
				ID:            s.ID,
				SeriesID:      s.SeriesID,
				TMDbID:        s.TmdbID,
				SeasonNumber:  s.SeasonNumber,
				Name:          s.Name,
				Overview:      s.Overview,
				PosterPath:    s.PosterPath,
				EpisodeCount:  s.EpisodeCount,
				AirDate:       pgDateToTimePtr(s.AirDate),
				VoteAverage:   pgNumericToDecimalPtr(s.VoteAverage),
				NamesI18n:     unmarshalStringMap(s.NamesI18n),
				OverviewsI18n: unmarshalStringMap(s.OverviewsI18n),
				CreatedAt:     s.CreatedAt,
				UpdatedAt:     s.UpdatedAt,
			},
			ActualEpisodeCount: s.ActualEpisodeCount,
		}
	}
	return result, nil
}

func (r *postgresRepository) CreateSeason(ctx context.Context, params CreateSeasonParams) (*Season, error) {
	dbSeason, err := r.queries.CreateSeason(ctx, tvshowdb.CreateSeasonParams{
		SeriesID:      params.SeriesID,
		TmdbID:        params.TMDbID,
		SeasonNumber:  params.SeasonNumber,
		Name:          params.Name,
		Overview:      params.Overview,
		PosterPath:    params.PosterPath,
		EpisodeCount:  params.EpisodeCount,
		AirDate:       stringToPgDate(params.AirDate),
		VoteAverage:   stringToPgNumeric(params.VoteAverage),
		NamesI18n:     marshalStringMap(params.NamesI18n),
		OverviewsI18n: marshalStringMap(params.OverviewsI18n),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create season: %w", err)
	}
	return dbSeasonToSeason(dbSeason), nil
}

func (r *postgresRepository) UpsertSeason(ctx context.Context, params CreateSeasonParams) (*Season, error) {
	dbSeason, err := r.queries.UpsertSeason(ctx, tvshowdb.UpsertSeasonParams{
		SeriesID:      params.SeriesID,
		TmdbID:        params.TMDbID,
		SeasonNumber:  params.SeasonNumber,
		Name:          params.Name,
		Overview:      params.Overview,
		PosterPath:    params.PosterPath,
		EpisodeCount:  params.EpisodeCount,
		AirDate:       stringToPgDate(params.AirDate),
		VoteAverage:   stringToPgNumeric(params.VoteAverage),
		NamesI18n:     marshalStringMap(params.NamesI18n),
		OverviewsI18n: marshalStringMap(params.OverviewsI18n),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upsert season: %w", err)
	}
	return dbSeasonToSeason(dbSeason), nil
}

func (r *postgresRepository) UpdateSeason(ctx context.Context, params UpdateSeasonParams) (*Season, error) {
	dbSeason, err := r.queries.UpdateSeason(ctx, tvshowdb.UpdateSeasonParams{
		ID:            params.ID,
		TmdbID:        params.TMDbID,
		Name:          params.Name,
		Overview:      params.Overview,
		PosterPath:    params.PosterPath,
		EpisodeCount:  params.EpisodeCount,
		AirDate:       stringToPgDate(params.AirDate),
		VoteAverage:   stringToPgNumeric(params.VoteAverage),
		NamesI18n:     marshalStringMapToBytes(params.NamesI18n),
		OverviewsI18n: marshalStringMapToBytes(params.OverviewsI18n),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update season: %w", err)
	}
	return dbSeasonToSeason(dbSeason), nil
}

func (r *postgresRepository) DeleteSeason(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteSeason(ctx, id)
}

func (r *postgresRepository) DeleteSeasonsBySeries(ctx context.Context, seriesID uuid.UUID) error {
	return r.queries.DeleteSeasonsBySeries(ctx, seriesID)
}

// =============================================================================
// Episode Operations
// =============================================================================

func (r *postgresRepository) GetEpisode(ctx context.Context, id uuid.UUID) (*Episode, error) {
	episode, err := r.queries.GetEpisode(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("episode not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get episode: %w", err)
	}
	return dbEpisodeToEpisode(episode), nil
}

func (r *postgresRepository) GetEpisodeByTMDbID(ctx context.Context, tmdbID int32) (*Episode, error) {
	episode, err := r.queries.GetEpisodeByTMDbID(ctx, &tmdbID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("episode not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get episode by TMDb ID: %w", err)
	}
	return dbEpisodeToEpisode(episode), nil
}

func (r *postgresRepository) GetEpisodeByNumber(ctx context.Context, seriesID uuid.UUID, seasonNumber, episodeNumber int32) (*Episode, error) {
	episode, err := r.queries.GetEpisodeByNumber(ctx, tvshowdb.GetEpisodeByNumberParams{
		SeriesID:      seriesID,
		SeasonNumber:  seasonNumber,
		EpisodeNumber: episodeNumber,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("episode not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get episode by number: %w", err)
	}
	return dbEpisodeToEpisode(episode), nil
}

func (r *postgresRepository) ListEpisodesBySeries(ctx context.Context, seriesID uuid.UUID) ([]Episode, error) {
	dbEpisodes, err := r.queries.ListEpisodesBySeries(ctx, seriesID)
	if err != nil {
		return nil, fmt.Errorf("failed to list episodes: %w", err)
	}

	result := make([]Episode, len(dbEpisodes))
	for i, e := range dbEpisodes {
		result[i] = *dbEpisodeToEpisode(e)
	}
	return result, nil
}

func (r *postgresRepository) ListEpisodesBySeason(ctx context.Context, seasonID uuid.UUID) ([]Episode, error) {
	dbEpisodes, err := r.queries.ListEpisodesBySeason(ctx, seasonID)
	if err != nil {
		return nil, fmt.Errorf("failed to list episodes: %w", err)
	}

	result := make([]Episode, len(dbEpisodes))
	for i, e := range dbEpisodes {
		result[i] = *dbEpisodeToEpisode(e)
	}
	return result, nil
}

func (r *postgresRepository) ListEpisodesBySeasonNumber(ctx context.Context, seriesID uuid.UUID, seasonNumber int32) ([]Episode, error) {
	dbEpisodes, err := r.queries.ListEpisodesBySeasonNumber(ctx, tvshowdb.ListEpisodesBySeasonNumberParams{
		SeriesID:     seriesID,
		SeasonNumber: seasonNumber,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list episodes: %w", err)
	}

	result := make([]Episode, len(dbEpisodes))
	for i, e := range dbEpisodes {
		result[i] = *dbEpisodeToEpisode(e)
	}
	return result, nil
}

func (r *postgresRepository) ListRecentEpisodes(ctx context.Context, limit, offset int32) ([]EpisodeWithSeriesInfo, error) {
	dbRows, err := r.queries.ListRecentEpisodes(ctx, tvshowdb.ListRecentEpisodesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list recent episodes: %w", err)
	}

	result := make([]EpisodeWithSeriesInfo, len(dbRows))
	for i, row := range dbRows {
		result[i] = EpisodeWithSeriesInfo{
			Episode: Episode{
				ID:             row.ID,
				SeriesID:       row.SeriesID,
				SeasonID:       row.SeasonID,
				TMDbID:         row.TmdbID,
				TVDbID:         row.TvdbID,
				IMDbID:         row.ImdbID,
				SeasonNumber:   row.SeasonNumber,
				EpisodeNumber:  row.EpisodeNumber,
				Title:          row.Title,
				Overview:       row.Overview,
				AirDate:        pgDateToTimePtr(row.AirDate),
				Runtime:        row.Runtime,
				VoteAverage:    pgNumericToDecimalPtr(row.VoteAverage),
				VoteCount:      row.VoteCount,
				StillPath:      row.StillPath,
				ProductionCode: row.ProductionCode,
				TitlesI18n:     unmarshalStringMap(row.TitlesI18n),
				OverviewsI18n:  unmarshalStringMap(row.OverviewsI18n),
				CreatedAt:      row.CreatedAt,
				UpdatedAt:      row.UpdatedAt,
			},
			SeriesID:         row.SeriesID,
			SeriesTitle:      row.SeriesTitle,
			SeriesPosterPath: row.SeriesPosterPath,
		}
	}
	return result, nil
}

func (r *postgresRepository) ListUpcomingEpisodes(ctx context.Context, limit, offset int32) ([]EpisodeWithSeriesInfo, error) {
	dbRows, err := r.queries.ListUpcomingEpisodes(ctx, tvshowdb.ListUpcomingEpisodesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list upcoming episodes: %w", err)
	}

	result := make([]EpisodeWithSeriesInfo, len(dbRows))
	for i, row := range dbRows {
		result[i] = EpisodeWithSeriesInfo{
			Episode: Episode{
				ID:             row.ID,
				SeriesID:       row.SeriesID,
				SeasonID:       row.SeasonID,
				TMDbID:         row.TmdbID,
				TVDbID:         row.TvdbID,
				IMDbID:         row.ImdbID,
				SeasonNumber:   row.SeasonNumber,
				EpisodeNumber:  row.EpisodeNumber,
				Title:          row.Title,
				Overview:       row.Overview,
				AirDate:        pgDateToTimePtr(row.AirDate),
				Runtime:        row.Runtime,
				VoteAverage:    pgNumericToDecimalPtr(row.VoteAverage),
				VoteCount:      row.VoteCount,
				StillPath:      row.StillPath,
				ProductionCode: row.ProductionCode,
				TitlesI18n:     unmarshalStringMap(row.TitlesI18n),
				OverviewsI18n:  unmarshalStringMap(row.OverviewsI18n),
				CreatedAt:      row.CreatedAt,
				UpdatedAt:      row.UpdatedAt,
			},
			SeriesID:         row.SeriesID,
			SeriesTitle:      row.SeriesTitle,
			SeriesPosterPath: row.SeriesPosterPath,
		}
	}
	return result, nil
}

func (r *postgresRepository) CountEpisodesBySeries(ctx context.Context, seriesID uuid.UUID) (int64, error) {
	return r.queries.CountEpisodesBySeries(ctx, seriesID)
}

func (r *postgresRepository) CountEpisodesBySeason(ctx context.Context, seasonID uuid.UUID) (int64, error) {
	return r.queries.CountEpisodesBySeason(ctx, seasonID)
}

func (r *postgresRepository) CreateEpisode(ctx context.Context, params CreateEpisodeParams) (*Episode, error) {
	dbEpisode, err := r.queries.CreateEpisode(ctx, tvshowdb.CreateEpisodeParams{
		SeriesID:       params.SeriesID,
		SeasonID:       params.SeasonID,
		TmdbID:         params.TMDbID,
		TvdbID:         params.TVDbID,
		ImdbID:         params.IMDbID,
		SeasonNumber:   params.SeasonNumber,
		EpisodeNumber:  params.EpisodeNumber,
		Title:          params.Title,
		Overview:       params.Overview,
		AirDate:        stringToPgDate(params.AirDate),
		Runtime:        params.Runtime,
		VoteAverage:    stringToPgNumeric(params.VoteAverage),
		VoteCount:      params.VoteCount,
		StillPath:      params.StillPath,
		ProductionCode: params.ProductionCode,
		TitlesI18n:     marshalStringMap(params.TitlesI18n),
		OverviewsI18n:  marshalStringMap(params.OverviewsI18n),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create episode: %w", err)
	}
	return dbEpisodeToEpisode(dbEpisode), nil
}

func (r *postgresRepository) UpsertEpisode(ctx context.Context, params CreateEpisodeParams) (*Episode, error) {
	dbEpisode, err := r.queries.UpsertEpisode(ctx, tvshowdb.UpsertEpisodeParams{
		SeriesID:       params.SeriesID,
		SeasonID:       params.SeasonID,
		TmdbID:         params.TMDbID,
		TvdbID:         params.TVDbID,
		ImdbID:         params.IMDbID,
		SeasonNumber:   params.SeasonNumber,
		EpisodeNumber:  params.EpisodeNumber,
		Title:          params.Title,
		Overview:       params.Overview,
		AirDate:        stringToPgDate(params.AirDate),
		Runtime:        params.Runtime,
		VoteAverage:    stringToPgNumeric(params.VoteAverage),
		VoteCount:      params.VoteCount,
		StillPath:      params.StillPath,
		ProductionCode: params.ProductionCode,
		TitlesI18n:     marshalStringMap(params.TitlesI18n),
		OverviewsI18n:  marshalStringMap(params.OverviewsI18n),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upsert episode: %w", err)
	}
	return dbEpisodeToEpisode(dbEpisode), nil
}

func (r *postgresRepository) UpdateEpisode(ctx context.Context, params UpdateEpisodeParams) (*Episode, error) {
	dbEpisode, err := r.queries.UpdateEpisode(ctx, tvshowdb.UpdateEpisodeParams{
		ID:             params.ID,
		TmdbID:         params.TMDbID,
		TvdbID:         params.TVDbID,
		ImdbID:         params.IMDbID,
		Title:          params.Title,
		Overview:       params.Overview,
		AirDate:        stringToPgDate(params.AirDate),
		Runtime:        params.Runtime,
		VoteAverage:    stringToPgNumeric(params.VoteAverage),
		VoteCount:      params.VoteCount,
		StillPath:      params.StillPath,
		ProductionCode: params.ProductionCode,
		TitlesI18n:     marshalStringMapToBytes(params.TitlesI18n),
		OverviewsI18n:  marshalStringMapToBytes(params.OverviewsI18n),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update episode: %w", err)
	}
	return dbEpisodeToEpisode(dbEpisode), nil
}

func (r *postgresRepository) DeleteEpisode(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteEpisode(ctx, id)
}

func (r *postgresRepository) DeleteEpisodesBySeason(ctx context.Context, seasonID uuid.UUID) error {
	return r.queries.DeleteEpisodesBySeason(ctx, seasonID)
}

func (r *postgresRepository) DeleteEpisodesBySeries(ctx context.Context, seriesID uuid.UUID) error {
	return r.queries.DeleteEpisodesBySeries(ctx, seriesID)
}

// =============================================================================
// Episode File Operations
// =============================================================================

func (r *postgresRepository) GetEpisodeFile(ctx context.Context, id uuid.UUID) (*EpisodeFile, error) {
	file, err := r.queries.GetEpisodeFile(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("episode file not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get episode file: %w", err)
	}
	return dbEpisodeFileToEpisodeFile(file), nil
}

func (r *postgresRepository) GetEpisodeFileByPath(ctx context.Context, path string) (*EpisodeFile, error) {
	file, err := r.queries.GetEpisodeFileByPath(ctx, path)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("episode file not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get episode file by path: %w", err)
	}
	return dbEpisodeFileToEpisodeFile(file), nil
}

func (r *postgresRepository) GetEpisodeFileBySonarrID(ctx context.Context, sonarrFileID int32) (*EpisodeFile, error) {
	file, err := r.queries.GetEpisodeFileBySonarrID(ctx, &sonarrFileID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("episode file not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get episode file by Sonarr ID: %w", err)
	}
	return dbEpisodeFileToEpisodeFile(file), nil
}

func (r *postgresRepository) ListEpisodeFilesByEpisode(ctx context.Context, episodeID uuid.UUID) ([]EpisodeFile, error) {
	dbFiles, err := r.queries.ListEpisodeFilesByEpisode(ctx, episodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to list episode files: %w", err)
	}

	result := make([]EpisodeFile, len(dbFiles))
	for i, f := range dbFiles {
		result[i] = *dbEpisodeFileToEpisodeFile(f)
	}
	return result, nil
}

func (r *postgresRepository) CreateEpisodeFile(ctx context.Context, params CreateEpisodeFileParams) (*EpisodeFile, error) {
	audioLangs := params.AudioLanguages
	if audioLangs == nil {
		audioLangs = []string{}
	}
	subtitleLangs := params.SubtitleLanguages
	if subtitleLangs == nil {
		subtitleLangs = []string{}
	}
	dbFile, err := r.queries.CreateEpisodeFile(ctx, tvshowdb.CreateEpisodeFileParams{
		EpisodeID:         params.EpisodeID,
		FilePath:          params.FilePath,
		FileName:          params.FileName,
		FileSize:          params.FileSize,
		Container:         params.Container,
		Resolution:        params.Resolution,
		QualityProfile:    params.QualityProfile,
		VideoCodec:        params.VideoCodec,
		AudioCodec:        params.AudioCodec,
		BitrateKbps:       params.BitrateKbps,
		DurationSeconds:   stringToPgNumeric(params.DurationSeconds),
		AudioLanguages:    audioLangs,
		SubtitleLanguages: subtitleLangs,
		SonarrFileID:      params.SonarrFileID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create episode file: %w", err)
	}
	return dbEpisodeFileToEpisodeFile(dbFile), nil
}

func (r *postgresRepository) UpdateEpisodeFile(ctx context.Context, params UpdateEpisodeFileParams) (*EpisodeFile, error) {
	dbFile, err := r.queries.UpdateEpisodeFile(ctx, tvshowdb.UpdateEpisodeFileParams{
		ID:                params.ID,
		FilePath:          params.FilePath,
		FileName:          params.FileName,
		FileSize:          params.FileSize,
		Container:         params.Container,
		Resolution:        params.Resolution,
		QualityProfile:    params.QualityProfile,
		VideoCodec:        params.VideoCodec,
		AudioCodec:        params.AudioCodec,
		BitrateKbps:       params.BitrateKbps,
		DurationSeconds:   stringToPgNumeric(params.DurationSeconds),
		AudioLanguages:    params.AudioLanguages,
		SubtitleLanguages: params.SubtitleLanguages,
		SonarrFileID:      params.SonarrFileID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update episode file: %w", err)
	}
	return dbEpisodeFileToEpisodeFile(dbFile), nil
}

func (r *postgresRepository) DeleteEpisodeFile(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteEpisodeFile(ctx, id)
}

func (r *postgresRepository) DeleteEpisodeFilesByEpisode(ctx context.Context, episodeID uuid.UUID) error {
	return r.queries.DeleteEpisodeFilesByEpisode(ctx, episodeID)
}

// =============================================================================
// Credits Operations
// =============================================================================

func (r *postgresRepository) CreateSeriesCredit(ctx context.Context, params CreateSeriesCreditParams) (*SeriesCredit, error) {
	dbCredit, err := r.queries.CreateSeriesCredit(ctx, tvshowdb.CreateSeriesCreditParams{
		SeriesID:     params.SeriesID,
		TmdbPersonID: params.TMDbPersonID,
		Name:         params.Name,
		CreditType:   params.CreditType,
		Character:    params.Character,
		CastOrder:    params.CastOrder,
		Job:          params.Job,
		Department:   params.Department,
		ProfilePath:  params.ProfilePath,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create series credit: %w", err)
	}
	return dbSeriesCreditToSeriesCredit(dbCredit), nil
}

func (r *postgresRepository) ListSeriesCast(ctx context.Context, seriesID uuid.UUID, limit, offset int32) ([]SeriesCredit, error) {
	dbCredits, err := r.queries.ListSeriesCast(ctx, tvshowdb.ListSeriesCastParams{
		SeriesID: seriesID,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list series cast: %w", err)
	}

	result := make([]SeriesCredit, len(dbCredits))
	for i, c := range dbCredits {
		result[i] = *dbSeriesCreditToSeriesCredit(c)
	}
	return result, nil
}

func (r *postgresRepository) CountSeriesCast(ctx context.Context, seriesID uuid.UUID) (int64, error) {
	return r.queries.CountSeriesCast(ctx, seriesID)
}

func (r *postgresRepository) ListSeriesCrew(ctx context.Context, seriesID uuid.UUID, limit, offset int32) ([]SeriesCredit, error) {
	dbCredits, err := r.queries.ListSeriesCrew(ctx, tvshowdb.ListSeriesCrewParams{
		SeriesID: seriesID,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list series crew: %w", err)
	}

	result := make([]SeriesCredit, len(dbCredits))
	for i, c := range dbCredits {
		result[i] = *dbSeriesCreditToSeriesCredit(c)
	}
	return result, nil
}

func (r *postgresRepository) CountSeriesCrew(ctx context.Context, seriesID uuid.UUID) (int64, error) {
	return r.queries.CountSeriesCrew(ctx, seriesID)
}

func (r *postgresRepository) DeleteSeriesCredits(ctx context.Context, seriesID uuid.UUID) error {
	return r.queries.DeleteSeriesCredits(ctx, seriesID)
}

func (r *postgresRepository) CreateEpisodeCredit(ctx context.Context, params CreateEpisodeCreditParams) (*EpisodeCredit, error) {
	dbCredit, err := r.queries.CreateEpisodeCredit(ctx, tvshowdb.CreateEpisodeCreditParams{
		EpisodeID:    params.EpisodeID,
		TmdbPersonID: params.TMDbPersonID,
		Name:         params.Name,
		CreditType:   params.CreditType,
		Character:    params.Character,
		CastOrder:    params.CastOrder,
		Job:          params.Job,
		Department:   params.Department,
		ProfilePath:  params.ProfilePath,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create episode credit: %w", err)
	}
	return dbEpisodeCreditToEpisodeCredit(dbCredit), nil
}

func (r *postgresRepository) ListEpisodeGuestStars(ctx context.Context, episodeID uuid.UUID) ([]EpisodeCredit, error) {
	dbCredits, err := r.queries.ListEpisodeGuestStars(ctx, episodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to list episode guest stars: %w", err)
	}

	result := make([]EpisodeCredit, len(dbCredits))
	for i, c := range dbCredits {
		result[i] = *dbEpisodeCreditToEpisodeCredit(c)
	}
	return result, nil
}

func (r *postgresRepository) ListEpisodeCrew(ctx context.Context, episodeID uuid.UUID) ([]EpisodeCredit, error) {
	dbCredits, err := r.queries.ListEpisodeCrew(ctx, episodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to list episode crew: %w", err)
	}

	result := make([]EpisodeCredit, len(dbCredits))
	for i, c := range dbCredits {
		result[i] = *dbEpisodeCreditToEpisodeCredit(c)
	}
	return result, nil
}

func (r *postgresRepository) DeleteEpisodeCredits(ctx context.Context, episodeID uuid.UUID) error {
	return r.queries.DeleteEpisodeCredits(ctx, episodeID)
}

// =============================================================================
// Genres Operations
// =============================================================================

func (r *postgresRepository) AddSeriesGenre(ctx context.Context, seriesID uuid.UUID, slug, name string) error {
	return r.queries.AddSeriesGenre(ctx, tvshowdb.AddSeriesGenreParams{
		SeriesID: seriesID,
		Slug:     slug,
		Name:     name,
	})
}

func (r *postgresRepository) ListSeriesGenres(ctx context.Context, seriesID uuid.UUID) ([]SeriesGenre, error) {
	dbGenres, err := r.queries.ListSeriesGenres(ctx, seriesID)
	if err != nil {
		return nil, fmt.Errorf("failed to list series genres: %w", err)
	}

	result := make([]SeriesGenre, len(dbGenres))
	for i, g := range dbGenres {
		result[i] = SeriesGenre{
			ID:        g.ID,
			SeriesID:  g.SeriesID,
			Slug:      g.Slug,
			Name:      g.Name,
			CreatedAt: g.CreatedAt,
		}
	}
	return result, nil
}

func (r *postgresRepository) ListDistinctSeriesGenres(ctx context.Context) ([]content.GenreSummary, error) {
	rows, err := r.queries.ListDistinctSeriesGenres(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list distinct series genres: %w", err)
	}
	genres := make([]content.GenreSummary, len(rows))
	for i, row := range rows {
		genres[i] = content.GenreSummary{
			Slug:      row.Slug,
			Name:      row.Name,
			ItemCount: row.ItemCount,
		}
	}
	return genres, nil
}

func (r *postgresRepository) DeleteSeriesGenres(ctx context.Context, seriesID uuid.UUID) error {
	return r.queries.DeleteSeriesGenres(ctx, seriesID)
}

// =============================================================================
// Networks Operations
// =============================================================================

func (r *postgresRepository) CreateNetwork(ctx context.Context, params CreateNetworkParams) (*Network, error) {
	dbNetwork, err := r.queries.CreateNetwork(ctx, tvshowdb.CreateNetworkParams{
		TmdbID:        params.TMDbID,
		Name:          params.Name,
		LogoPath:      params.LogoPath,
		OriginCountry: params.OriginCountry,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create network: %w", err)
	}
	return &Network{
		ID:            dbNetwork.ID,
		TMDbID:        dbNetwork.TmdbID,
		Name:          dbNetwork.Name,
		LogoPath:      dbNetwork.LogoPath,
		OriginCountry: dbNetwork.OriginCountry,
		CreatedAt:     dbNetwork.CreatedAt,
	}, nil
}

func (r *postgresRepository) GetNetwork(ctx context.Context, id uuid.UUID) (*Network, error) {
	network, err := r.queries.GetNetwork(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("network not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get network: %w", err)
	}
	return &Network{
		ID:            network.ID,
		TMDbID:        network.TmdbID,
		Name:          network.Name,
		LogoPath:      network.LogoPath,
		OriginCountry: network.OriginCountry,
		CreatedAt:     network.CreatedAt,
	}, nil
}

func (r *postgresRepository) GetNetworkByTMDbID(ctx context.Context, tmdbID int32) (*Network, error) {
	network, err := r.queries.GetNetworkByTMDbID(ctx, &tmdbID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("network not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get network by TMDb ID: %w", err)
	}
	return &Network{
		ID:            network.ID,
		TMDbID:        network.TmdbID,
		Name:          network.Name,
		LogoPath:      network.LogoPath,
		OriginCountry: network.OriginCountry,
		CreatedAt:     network.CreatedAt,
	}, nil
}

func (r *postgresRepository) ListNetworksBySeries(ctx context.Context, seriesID uuid.UUID) ([]Network, error) {
	dbNetworks, err := r.queries.ListNetworksBySeries(ctx, seriesID)
	if err != nil {
		return nil, fmt.Errorf("failed to list networks: %w", err)
	}

	result := make([]Network, len(dbNetworks))
	for i, n := range dbNetworks {
		result[i] = Network{
			ID:            n.ID,
			TMDbID:        n.TmdbID,
			Name:          n.Name,
			LogoPath:      n.LogoPath,
			OriginCountry: n.OriginCountry,
			CreatedAt:     n.CreatedAt,
		}
	}
	return result, nil
}

func (r *postgresRepository) AddSeriesNetwork(ctx context.Context, seriesID, networkID uuid.UUID) error {
	return r.queries.AddSeriesNetwork(ctx, tvshowdb.AddSeriesNetworkParams{
		SeriesID:  seriesID,
		NetworkID: networkID,
	})
}

func (r *postgresRepository) DeleteSeriesNetworks(ctx context.Context, seriesID uuid.UUID) error {
	return r.queries.DeleteSeriesNetworks(ctx, seriesID)
}

// =============================================================================
// Watch Progress Operations
// =============================================================================

func (r *postgresRepository) CreateOrUpdateWatchProgress(ctx context.Context, params CreateWatchProgressParams) (*EpisodeWatched, error) {
	dbWatched, err := r.queries.CreateOrUpdateWatchProgress(ctx, tvshowdb.CreateOrUpdateWatchProgressParams{
		UserID:          params.UserID,
		EpisodeID:       params.EpisodeID,
		ProgressSeconds: params.ProgressSeconds,
		DurationSeconds: params.DurationSeconds,
		IsCompleted:     params.IsCompleted,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create/update watch progress: %w", err)
	}
	return dbEpisodeWatchedToEpisodeWatched(dbWatched), nil
}

func (r *postgresRepository) MarkEpisodeWatched(ctx context.Context, userID, episodeID uuid.UUID, durationSeconds int32) (*EpisodeWatched, error) {
	dbWatched, err := r.queries.MarkEpisodeWatched(ctx, tvshowdb.MarkEpisodeWatchedParams{
		UserID:          userID,
		EpisodeID:       episodeID,
		ProgressSeconds: durationSeconds,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to mark episode watched: %w", err)
	}
	return dbEpisodeWatchedToEpisodeWatched(dbWatched), nil
}

func (r *postgresRepository) MarkEpisodesWatchedBulk(ctx context.Context, userID uuid.UUID, episodeIDs []uuid.UUID) (int64, error) {
	affected, err := r.queries.MarkEpisodesWatchedBulk(ctx, tvshowdb.MarkEpisodesWatchedBulkParams{
		UserID:     userID,
		EpisodeIds: episodeIDs,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to bulk mark episodes watched: %w", err)
	}
	return affected, nil
}

func (r *postgresRepository) GetWatchProgress(ctx context.Context, userID, episodeID uuid.UUID) (*EpisodeWatched, error) {
	watched, err := r.queries.GetEpisodeWatchProgress(ctx, tvshowdb.GetEpisodeWatchProgressParams{
		UserID:    userID,
		EpisodeID: episodeID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("watch progress not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get watch progress: %w", err)
	}
	return dbEpisodeWatchedToEpisodeWatched(watched), nil
}

func (r *postgresRepository) DeleteWatchProgress(ctx context.Context, userID, episodeID uuid.UUID) error {
	return r.queries.DeleteWatchProgress(ctx, tvshowdb.DeleteWatchProgressParams{
		UserID:    userID,
		EpisodeID: episodeID,
	})
}

func (r *postgresRepository) DeleteSeriesWatchProgress(ctx context.Context, userID, seriesID uuid.UUID) error {
	return r.queries.DeleteSeriesWatchProgress(ctx, tvshowdb.DeleteSeriesWatchProgressParams{
		UserID:   userID,
		SeriesID: seriesID,
	})
}

func (r *postgresRepository) ListContinueWatchingSeries(ctx context.Context, userID uuid.UUID, limit int32) ([]ContinueWatchingItem, error) {
	dbRows, err := r.queries.ListContinueWatchingSeries(ctx, tvshowdb.ListContinueWatchingSeriesParams{
		UserID: userID,
		Limit:  limit,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list continue watching: %w", err)
	}

	result := make([]ContinueWatchingItem, len(dbRows))
	for i, row := range dbRows {
		result[i] = ContinueWatchingItem{
			Series: &Series{
				ID:                row.ID,
				TMDbID:            row.TmdbID,
				TVDbID:            row.TvdbID,
				IMDbID:            row.ImdbID,
				SonarrID:          row.SonarrID,
				Title:             row.Title,
				OriginalTitle:     row.OriginalTitle,
				OriginalLanguage:  row.OriginalLanguage,
				Tagline:           row.Tagline,
				Overview:          row.Overview,
				Status:            row.Status,
				Type:              row.Type,
				FirstAirDate:      pgDateToTimePtr(row.FirstAirDate),
				LastAirDate:       pgDateToTimePtr(row.LastAirDate),
				VoteAverage:       pgNumericToDecimalPtr(row.VoteAverage),
				VoteCount:         row.VoteCount,
				Popularity:        pgNumericToDecimalPtr(row.Popularity),
				PosterPath:        row.PosterPath,
				BackdropPath:      row.BackdropPath,
				TotalSeasons:      row.TotalSeasons,
				TotalEpisodes:     row.TotalEpisodes,
				TrailerURL:        row.TrailerUrl,
				Homepage:          row.Homepage,
				TitlesI18n:        unmarshalStringMap(row.TitlesI18n),
				TaglinesI18n:      unmarshalStringMap(row.TaglinesI18n),
				OverviewsI18n:     unmarshalStringMap(row.OverviewsI18n),
				AgeRatings:        unmarshalNestedStringMap(row.AgeRatings),
				MetadataUpdatedAt: pgTimestamptzToTimePtr(row.MetadataUpdatedAt),
				CreatedAt:         row.CreatedAt,
				UpdatedAt:         row.UpdatedAt,
			},
			LastEpisodeID:     row.LastEpisodeID,
			LastSeasonNumber:  row.LastSeasonNumber,
			LastEpisodeNumber: row.LastEpisodeNumber,
			LastEpisodeTitle:  row.LastEpisodeTitle,
			ProgressSeconds:   row.ProgressSeconds,
			DurationSeconds:   row.DurationSeconds,
			LastWatchedAt:     pgTimestamptzToTime(row.LastWatchedAt),
		}
	}
	return result, nil
}

func (r *postgresRepository) ListWatchedEpisodesBySeries(ctx context.Context, userID, seriesID uuid.UUID) ([]WatchedEpisodeItem, error) {
	dbRows, err := r.queries.ListWatchedEpisodesBySeries(ctx, tvshowdb.ListWatchedEpisodesBySeriesParams{
		UserID:   userID,
		SeriesID: seriesID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list watched episodes: %w", err)
	}

	result := make([]WatchedEpisodeItem, len(dbRows))
	for i, row := range dbRows {
		result[i] = WatchedEpisodeItem{
			EpisodeWatched: EpisodeWatched{
				ID:              row.ID,
				UserID:          row.UserID,
				EpisodeID:       row.EpisodeID,
				ProgressSeconds: row.ProgressSeconds,
				DurationSeconds: row.DurationSeconds,
				IsCompleted:     row.IsCompleted,
				WatchCount:      row.WatchCount,
				LastWatchedAt:   pgTimestamptzToTime(row.LastWatchedAt),
				CreatedAt:       row.CreatedAt,
				UpdatedAt:       row.UpdatedAt,
			},
			SeasonNumber:  row.SeasonNumber,
			EpisodeNumber: row.EpisodeNumber,
			EpisodeTitle:  row.EpisodeTitle,
		}
	}
	return result, nil
}

func (r *postgresRepository) ListWatchedEpisodesByUser(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]WatchedEpisodeItem, error) {
	dbRows, err := r.queries.ListWatchedEpisodesByUser(ctx, tvshowdb.ListWatchedEpisodesByUserParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list watched episodes: %w", err)
	}

	result := make([]WatchedEpisodeItem, len(dbRows))
	for i, row := range dbRows {
		result[i] = WatchedEpisodeItem{
			EpisodeWatched: EpisodeWatched{
				ID:              row.ID,
				UserID:          row.UserID,
				EpisodeID:       row.EpisodeID,
				ProgressSeconds: row.ProgressSeconds,
				DurationSeconds: row.DurationSeconds,
				IsCompleted:     row.IsCompleted,
				WatchCount:      row.WatchCount,
				LastWatchedAt:   pgTimestamptzToTime(row.LastWatchedAt),
				CreatedAt:       row.CreatedAt,
				UpdatedAt:       row.UpdatedAt,
			},
			SeasonNumber:     row.SeasonNumber,
			EpisodeNumber:    row.EpisodeNumber,
			EpisodeTitle:     row.EpisodeTitle,
			EpisodeStillPath: row.EpisodeStillPath,
			SeriesID:         row.SeriesID,
			SeriesTitle:      row.SeriesTitle,
			SeriesPosterPath: row.SeriesPosterPath,
		}
	}
	return result, nil
}

func (r *postgresRepository) GetSeriesWatchStats(ctx context.Context, userID, seriesID uuid.UUID) (*SeriesWatchStats, error) {
	stats, err := r.queries.GetSeriesWatchStats(ctx, tvshowdb.GetSeriesWatchStatsParams{
		UserID:   userID,
		SeriesID: seriesID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get series watch stats: %w", err)
	}
	return &SeriesWatchStats{
		WatchedCount:    stats.WatchedCount,
		InProgressCount: stats.InProgressCount,
		TotalWatches:    stats.TotalWatches,
		TotalEpisodes:   stats.TotalEpisodes,
	}, nil
}

func (r *postgresRepository) GetUserTVStats(ctx context.Context, userID uuid.UUID) (*UserTVStats, error) {
	stats, err := r.queries.GetUserTVStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user TV stats: %w", err)
	}
	// TotalWatches is interface{} from sqlc, handle type assertion
	var totalWatches int64
	if tw, ok := stats.TotalWatches.(int64); ok {
		totalWatches = tw
	}
	return &UserTVStats{
		SeriesCount:        stats.SeriesCount,
		EpisodesWatched:    stats.EpisodesWatched,
		EpisodesInProgress: stats.EpisodesInProgress,
		TotalWatches:       totalWatches,
	}, nil
}

func (r *postgresRepository) GetNextUnwatchedEpisode(ctx context.Context, userID, seriesID uuid.UUID) (*Episode, error) {
	episode, err := r.queries.GetNextUnwatchedEpisode(ctx, tvshowdb.GetNextUnwatchedEpisodeParams{
		UserID:   userID,
		SeriesID: seriesID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No unwatched episodes
		}
		return nil, fmt.Errorf("failed to get next unwatched episode: %w", err)
	}
	return dbEpisodeToEpisode(episode), nil
}

// =============================================================================
// Conversion Helpers
// =============================================================================

func dbSeriesToSeries(s tvshowdb.TvshowSeries) *Series {
	return &Series{
		ID:                s.ID,
		TMDbID:            s.TmdbID,
		TVDbID:            s.TvdbID,
		IMDbID:            s.ImdbID,
		SonarrID:          s.SonarrID,
		Title:             s.Title,
		OriginalTitle:     s.OriginalTitle,
		OriginalLanguage:  s.OriginalLanguage,
		Tagline:           s.Tagline,
		Overview:          s.Overview,
		Status:            s.Status,
		Type:              s.Type,
		FirstAirDate:      pgDateToTimePtr(s.FirstAirDate),
		LastAirDate:       pgDateToTimePtr(s.LastAirDate),
		VoteAverage:       pgNumericToDecimalPtr(s.VoteAverage),
		VoteCount:         s.VoteCount,
		Popularity:        pgNumericToDecimalPtr(s.Popularity),
		PosterPath:        s.PosterPath,
		BackdropPath:      s.BackdropPath,
		TotalSeasons:      s.TotalSeasons,
		TotalEpisodes:     s.TotalEpisodes,
		TrailerURL:        s.TrailerUrl,
		Homepage:          s.Homepage,
		TitlesI18n:        unmarshalStringMap(s.TitlesI18n),
		TaglinesI18n:      unmarshalStringMap(s.TaglinesI18n),
		OverviewsI18n:     unmarshalStringMap(s.OverviewsI18n),
		AgeRatings:        unmarshalNestedStringMap(s.AgeRatings),
		ExternalRatings:   unmarshalExternalRatings(s.ExternalRatings),
		MetadataUpdatedAt: pgTimestamptzToTimePtr(s.MetadataUpdatedAt),
		CreatedAt:         s.CreatedAt,
		UpdatedAt:         s.UpdatedAt,
	}
}

func dbSeasonToSeason(s tvshowdb.TvshowSeason) *Season {
	return &Season{
		ID:            s.ID,
		SeriesID:      s.SeriesID,
		TMDbID:        s.TmdbID,
		SeasonNumber:  s.SeasonNumber,
		Name:          s.Name,
		Overview:      s.Overview,
		PosterPath:    s.PosterPath,
		EpisodeCount:  s.EpisodeCount,
		AirDate:       pgDateToTimePtr(s.AirDate),
		VoteAverage:   pgNumericToDecimalPtr(s.VoteAverage),
		NamesI18n:     unmarshalStringMap(s.NamesI18n),
		OverviewsI18n: unmarshalStringMap(s.OverviewsI18n),
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
	}
}

func dbEpisodeToEpisode(e tvshowdb.TvshowEpisode) *Episode {
	return &Episode{
		ID:             e.ID,
		SeriesID:       e.SeriesID,
		SeasonID:       e.SeasonID,
		TMDbID:         e.TmdbID,
		TVDbID:         e.TvdbID,
		IMDbID:         e.ImdbID,
		SeasonNumber:   e.SeasonNumber,
		EpisodeNumber:  e.EpisodeNumber,
		Title:          e.Title,
		Overview:       e.Overview,
		AirDate:        pgDateToTimePtr(e.AirDate),
		Runtime:        e.Runtime,
		VoteAverage:    pgNumericToDecimalPtr(e.VoteAverage),
		VoteCount:      e.VoteCount,
		StillPath:      e.StillPath,
		ProductionCode: e.ProductionCode,
		TitlesI18n:     unmarshalStringMap(e.TitlesI18n),
		OverviewsI18n:  unmarshalStringMap(e.OverviewsI18n),
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}
}

func dbEpisodeFileToEpisodeFile(f tvshowdb.TvshowEpisodeFile) *EpisodeFile {
	return &EpisodeFile{
		ID:                f.ID,
		EpisodeID:         f.EpisodeID,
		FilePath:          f.FilePath,
		FileName:          f.FileName,
		FileSize:          f.FileSize,
		Container:         f.Container,
		Resolution:        f.Resolution,
		QualityProfile:    f.QualityProfile,
		VideoCodec:        f.VideoCodec,
		AudioCodec:        f.AudioCodec,
		BitrateKbps:       f.BitrateKbps,
		DurationSeconds:   pgNumericToDecimalPtr(f.DurationSeconds),
		AudioLanguages:    f.AudioLanguages,
		SubtitleLanguages: f.SubtitleLanguages,
		SonarrFileID:      f.SonarrFileID,
		CreatedAt:         f.CreatedAt,
		UpdatedAt:         f.UpdatedAt,
	}
}

func dbSeriesCreditToSeriesCredit(c tvshowdb.TvshowSeriesCredit) *SeriesCredit {
	return &SeriesCredit{
		ID:           c.ID,
		SeriesID:     c.SeriesID,
		TMDbPersonID: c.TmdbPersonID,
		Name:         c.Name,
		CreditType:   c.CreditType,
		Character:    c.Character,
		Job:          c.Job,
		Department:   c.Department,
		CastOrder:    c.CastOrder,
		ProfilePath:  c.ProfilePath,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}
}

func dbEpisodeCreditToEpisodeCredit(c tvshowdb.TvshowEpisodeCredit) *EpisodeCredit {
	return &EpisodeCredit{
		ID:           c.ID,
		EpisodeID:    c.EpisodeID,
		TMDbPersonID: c.TmdbPersonID,
		Name:         c.Name,
		CreditType:   c.CreditType,
		Character:    c.Character,
		Job:          c.Job,
		Department:   c.Department,
		CastOrder:    c.CastOrder,
		ProfilePath:  c.ProfilePath,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}
}

func dbEpisodeWatchedToEpisodeWatched(w tvshowdb.TvshowEpisodeWatched) *EpisodeWatched {
	return &EpisodeWatched{
		ID:              w.ID,
		UserID:          w.UserID,
		EpisodeID:       w.EpisodeID,
		ProgressSeconds: w.ProgressSeconds,
		DurationSeconds: w.DurationSeconds,
		IsCompleted:     w.IsCompleted,
		WatchCount:      w.WatchCount,
		LastWatchedAt:   pgTimestamptzToTime(w.LastWatchedAt),
		CreatedAt:       w.CreatedAt,
		UpdatedAt:       w.UpdatedAt,
	}
}

// =============================================================================
// Type Conversion Helpers
// =============================================================================

func pgDateToTimePtr(d pgtype.Date) *time.Time {
	if !d.Valid {
		return nil
	}
	return &d.Time
}

func pgTimestamptzToTimePtr(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

func pgTimestamptzToTime(t pgtype.Timestamptz) time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}

func pgNumericToDecimalPtr(n pgtype.Numeric) *decimal.Decimal {
	if !n.Valid {
		return nil
	}
	// Convert to string first
	var buf []byte
	buf, _ = n.MarshalJSON()
	d, err := decimal.Parse(string(buf))
	if err != nil {
		return nil
	}
	return &d
}

func stringToPgDate(s *string) pgtype.Date {
	if s == nil || *s == "" {
		return pgtype.Date{Valid: false}
	}
	t, err := time.Parse("2006-01-02", *s)
	if err != nil {
		return pgtype.Date{Valid: false}
	}
	return pgtype.Date{Time: t, Valid: true}
}

func stringToPgNumeric(s *string) pgtype.Numeric {
	if s == nil || *s == "" {
		return pgtype.Numeric{Valid: false}
	}
	var n pgtype.Numeric
	if err := n.Scan(*s); err != nil {
		return pgtype.Numeric{Valid: false}
	}
	return n
}

func stringToPgTimestamptz(s *string) pgtype.Timestamptz {
	if s == nil || *s == "" {
		return pgtype.Timestamptz{Valid: false}
	}
	t, err := time.Parse(time.RFC3339, *s)
	if err != nil {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: t, Valid: true}
}

func marshalStringMap(m map[string]string) json.RawMessage {
	if m == nil {
		return json.RawMessage("{}")
	}
	b, _ := json.Marshal(m)
	return b
}

func marshalStringMapToBytes(m map[string]string) []byte {
	if m == nil {
		return []byte("{}")
	}
	b, _ := json.Marshal(m)
	return b
}

func marshalNestedStringMap(m map[string]map[string]string) json.RawMessage {
	if m == nil {
		return json.RawMessage("{}")
	}
	b, _ := json.Marshal(m)
	return b
}

func marshalNestedStringMapToBytes(m map[string]map[string]string) []byte {
	if m == nil {
		return []byte("{}")
	}
	b, _ := json.Marshal(m)
	return b
}

func unmarshalStringMap(data json.RawMessage) map[string]string {
	if len(data) == 0 {
		return nil
	}
	var m map[string]string
	_ = json.Unmarshal(data, &m)
	return m
}

func unmarshalNestedStringMap(data json.RawMessage) map[string]map[string]string {
	if len(data) == 0 {
		return nil
	}
	var m map[string]map[string]string
	_ = json.Unmarshal(data, &m)
	return m
}

// marshalExternalRatings marshals []ExternalRating to JSONB json.RawMessage
func marshalExternalRatings(ratings []ExternalRating) json.RawMessage {
	if ratings == nil {
		return json.RawMessage("[]")
	}
	b, _ := json.Marshal(ratings)
	return b
}

// marshalExternalRatingsToBytes marshals []ExternalRating to JSONB []byte
func marshalExternalRatingsToBytes(ratings []ExternalRating) []byte {
	if ratings == nil {
		return []byte("[]")
	}
	b, _ := json.Marshal(ratings)
	return b
}

// unmarshalExternalRatings unmarshals JSONB json.RawMessage to []ExternalRating
func unmarshalExternalRatings(data json.RawMessage) []ExternalRating {
	if len(data) == 0 {
		return nil
	}
	var result []ExternalRating
	_ = json.Unmarshal(data, &result)
	return result
}
