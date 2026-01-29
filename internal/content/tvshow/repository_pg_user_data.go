package tvshow

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/lusoris/revenge/internal/content/tvshow/db"
)

// Series User Data

// GetSeriesUserRating retrieves a user's rating for a series.
func (r *pgRepository) GetSeriesUserRating(ctx context.Context, userID, seriesID uuid.UUID) (*SeriesUserRating, error) {
	rating, err := r.queries.GetSeriesUserRating(ctx, db.GetSeriesUserRatingParams{
		UserID:   userID,
		SeriesID: seriesID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	f, _ := rating.Rating.Float64Value()
	result := &SeriesUserRating{
		UserID:    rating.UserID,
		SeriesID:  rating.SeriesID,
		Rating:    f.Float64,
		CreatedAt: rating.CreatedAt,
		UpdatedAt: rating.UpdatedAt,
	}
	if rating.Review != nil {
		result.Review = *rating.Review
	}
	return result, nil
}

// SetSeriesUserRating sets a user's rating for a series.
func (r *pgRepository) SetSeriesUserRating(ctx context.Context, userID, seriesID uuid.UUID, rating float64, review string) error {
	var rev *string
	if review != "" {
		rev = &review
	}
	_, err := r.queries.SetSeriesUserRating(ctx, db.SetSeriesUserRatingParams{
		UserID:   userID,
		SeriesID: seriesID,
		Rating:   numericFromFloat(rating),
		Review:   rev,
	})
	return err
}

// DeleteSeriesUserRating deletes a user's rating for a series.
func (r *pgRepository) DeleteSeriesUserRating(ctx context.Context, userID, seriesID uuid.UUID) error {
	return r.queries.DeleteSeriesUserRating(ctx, db.DeleteSeriesUserRatingParams{
		UserID:   userID,
		SeriesID: seriesID,
	})
}

// IsSeriesFavorite checks if a series is favorited by a user.
func (r *pgRepository) IsSeriesFavorite(ctx context.Context, userID, seriesID uuid.UUID) (bool, error) {
	return r.queries.IsSeriesFavorite(ctx, db.IsSeriesFavoriteParams{
		UserID:   userID,
		SeriesID: seriesID,
	})
}

// AddSeriesFavorite adds a series to a user's favorites.
func (r *pgRepository) AddSeriesFavorite(ctx context.Context, userID, seriesID uuid.UUID) error {
	return r.queries.AddSeriesFavorite(ctx, db.AddSeriesFavoriteParams{
		UserID:   userID,
		SeriesID: seriesID,
	})
}

// RemoveSeriesFavorite removes a series from a user's favorites.
func (r *pgRepository) RemoveSeriesFavorite(ctx context.Context, userID, seriesID uuid.UUID) error {
	return r.queries.RemoveSeriesFavorite(ctx, db.RemoveSeriesFavoriteParams{
		UserID:   userID,
		SeriesID: seriesID,
	})
}

// ListFavoriteSeries lists a user's favorite series.
func (r *pgRepository) ListFavoriteSeries(ctx context.Context, userID uuid.UUID, params ListParams) ([]*Series, error) {
	rows, err := r.queries.ListUserFavoriteSeries(ctx, db.ListUserFavoriteSeriesParams{
		UserID: userID,
		Limit:  int32(params.Limit),
		Offset: int32(params.Offset),
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

// CountFavoriteSeries counts a user's favorite series.
func (r *pgRepository) CountFavoriteSeries(ctx context.Context, userID uuid.UUID) (int64, error) {
	return r.queries.CountUserFavoriteSeries(ctx, userID)
}

// IsSeriesInWatchlist checks if a series is in a user's watchlist.
func (r *pgRepository) IsSeriesInWatchlist(ctx context.Context, userID, seriesID uuid.UUID) (bool, error) {
	return r.queries.IsSeriesInWatchlist(ctx, db.IsSeriesInWatchlistParams{
		UserID:   userID,
		SeriesID: seriesID,
	})
}

// AddSeriesToWatchlist adds a series to a user's watchlist.
func (r *pgRepository) AddSeriesToWatchlist(ctx context.Context, userID, seriesID uuid.UUID) error {
	return r.queries.AddSeriesToWatchlist(ctx, db.AddSeriesToWatchlistParams{
		UserID:   userID,
		SeriesID: seriesID,
	})
}

// RemoveSeriesFromWatchlist removes a series from a user's watchlist.
func (r *pgRepository) RemoveSeriesFromWatchlist(ctx context.Context, userID, seriesID uuid.UUID) error {
	return r.queries.RemoveSeriesFromWatchlist(ctx, db.RemoveSeriesFromWatchlistParams{
		UserID:   userID,
		SeriesID: seriesID,
	})
}

// ListSeriesWatchlist lists a user's series watchlist.
func (r *pgRepository) ListSeriesWatchlist(ctx context.Context, userID uuid.UUID, params ListParams) ([]*Series, error) {
	rows, err := r.queries.ListUserSeriesWatchlist(ctx, db.ListUserSeriesWatchlistParams{
		UserID: userID,
		Limit:  int32(params.Limit),
		Offset: int32(params.Offset),
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

// CountSeriesWatchlist counts a user's series watchlist.
func (r *pgRepository) CountSeriesWatchlist(ctx context.Context, userID uuid.UUID) (int64, error) {
	return r.queries.CountUserSeriesWatchlist(ctx, userID)
}

// Episode User Data

// GetEpisodeUserRating retrieves a user's rating for an episode.
func (r *pgRepository) GetEpisodeUserRating(ctx context.Context, userID, episodeID uuid.UUID) (*EpisodeUserRating, error) {
	rating, err := r.queries.GetEpisodeUserRating(ctx, db.GetEpisodeUserRatingParams{
		UserID:    userID,
		EpisodeID: episodeID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	f, _ := rating.Rating.Float64Value()
	return &EpisodeUserRating{
		UserID:    rating.UserID,
		EpisodeID: rating.EpisodeID,
		Rating:    f.Float64,
		CreatedAt: rating.CreatedAt,
		UpdatedAt: rating.UpdatedAt,
	}, nil
}

// SetEpisodeUserRating sets a user's rating for an episode.
func (r *pgRepository) SetEpisodeUserRating(ctx context.Context, userID, episodeID uuid.UUID, rating float64) error {
	_, err := r.queries.SetEpisodeUserRating(ctx, db.SetEpisodeUserRatingParams{
		UserID:    userID,
		EpisodeID: episodeID,
		Rating:    numericFromFloat(rating),
	})
	return err
}

// DeleteEpisodeUserRating deletes a user's rating for an episode.
func (r *pgRepository) DeleteEpisodeUserRating(ctx context.Context, userID, episodeID uuid.UUID) error {
	return r.queries.DeleteEpisodeUserRating(ctx, db.DeleteEpisodeUserRatingParams{
		UserID:    userID,
		EpisodeID: episodeID,
	})
}

// Episode Watch History

// GetEpisodeWatchHistory retrieves a user's watch history for an episode.
func (r *pgRepository) GetEpisodeWatchHistory(ctx context.Context, userID, episodeID uuid.UUID) (*EpisodeWatchHistory, error) {
	history, err := r.queries.GetEpisodeWatchHistory(ctx, db.GetEpisodeWatchHistoryParams{
		UserID:    userID,
		EpisodeID: episodeID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return fromDBEpisodeWatchHistory(&history), nil
}

// CreateEpisodeWatchHistory creates a watch history entry for an episode.
func (r *pgRepository) CreateEpisodeWatchHistory(ctx context.Context, history *EpisodeWatchHistory) error {
	params := db.CreateEpisodeWatchHistoryParams{
		UserID:        history.UserID,
		EpisodeID:     history.EpisodeID,
		PositionTicks: history.PositionTicks,
	}

	if history.ProfileID != nil {
		params.ProfileID = pgtype.UUID{Bytes: *history.ProfileID, Valid: true}
	}
	if history.DurationTicks > 0 {
		params.DurationTicks = &history.DurationTicks
	}
	if history.DeviceName != "" {
		params.DeviceName = &history.DeviceName
	}
	if history.DeviceType != "" {
		params.DeviceType = &history.DeviceType
	}
	if history.ClientName != "" {
		params.ClientName = &history.ClientName
	}
	if history.PlayMethod != "" {
		params.PlayMethod = &history.PlayMethod
	}

	h, err := r.queries.CreateEpisodeWatchHistory(ctx, params)
	if err != nil {
		return err
	}
	history.ID = h.ID
	history.StartedAt = h.StartedAt
	history.LastUpdatedAt = h.LastUpdatedAt
	return nil
}

// UpdateEpisodeWatchHistory updates a watch history entry.
func (r *pgRepository) UpdateEpisodeWatchHistory(ctx context.Context, id uuid.UUID, positionTicks int64, durationTicks *int64) error {
	_, err := r.queries.UpdateEpisodeWatchHistory(ctx, db.UpdateEpisodeWatchHistoryParams{
		ID:            id,
		PositionTicks: positionTicks,
		DurationTicks: durationTicks,
	})
	return err
}

// MarkEpisodeWatchHistoryCompleted marks a watch history entry as completed.
func (r *pgRepository) MarkEpisodeWatchHistoryCompleted(ctx context.Context, id uuid.UUID) error {
	return r.queries.MarkEpisodeWatchHistoryCompleted(ctx, id)
}

// DeleteEpisodeWatchHistory deletes a watch history entry.
func (r *pgRepository) DeleteEpisodeWatchHistory(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteEpisodeWatchHistory(ctx, id)
}

// ListResumeableEpisodes lists episodes a user can resume.
func (r *pgRepository) ListResumeableEpisodes(ctx context.Context, userID uuid.UUID, limit int) ([]EpisodeWatchHistory, error) {
	rows, err := r.queries.ListResumeableEpisodes(ctx, db.ListResumeableEpisodesParams{
		UserID: userID,
		Limit:  int32(limit),
	})
	if err != nil {
		return nil, err
	}

	history := make([]EpisodeWatchHistory, len(rows))
	for i, row := range rows {
		history[i] = EpisodeWatchHistory{
			ID:            row.ID,
			UserID:        row.UserID,
			EpisodeID:     row.EpisodeID,
			PositionTicks: row.PositionTicks,
			Completed:     row.Completed,
			StartedAt:     row.StartedAt,
			LastUpdatedAt: row.LastUpdatedAt,
		}
		if row.ProfileID.Valid {
			pid := uuid.UUID(row.ProfileID.Bytes)
			history[i].ProfileID = &pid
		}
		if row.DurationTicks != nil {
			history[i].DurationTicks = *row.DurationTicks
		}
		if row.PlayedPercentage.Valid {
			f, _ := row.PlayedPercentage.Float64Value()
			history[i].PlayedPercentage = f.Float64
		}
		if row.DeviceName != nil {
			history[i].DeviceName = *row.DeviceName
		}
		if row.DeviceType != nil {
			history[i].DeviceType = *row.DeviceType
		}
		if row.ClientName != nil {
			history[i].ClientName = *row.ClientName
		}
		if row.PlayMethod != nil {
			history[i].PlayMethod = *row.PlayMethod
		}
	}
	return history, nil
}

// IsEpisodeWatched checks if an episode has been watched by a user.
func (r *pgRepository) IsEpisodeWatched(ctx context.Context, userID, episodeID uuid.UUID) (bool, error) {
	return r.queries.IsEpisodeWatched(ctx, db.IsEpisodeWatchedParams{
		UserID:    userID,
		EpisodeID: episodeID,
	})
}

// CountWatchedEpisodes counts episodes watched by a user.
func (r *pgRepository) CountWatchedEpisodes(ctx context.Context, userID uuid.UUID) (int64, error) {
	return r.queries.CountUserWatchedEpisodes(ctx, userID)
}

// CountWatchedEpisodesBySeries counts episodes watched by a user for a specific series.
func (r *pgRepository) CountWatchedEpisodesBySeries(ctx context.Context, userID, seriesID uuid.UUID) (int64, error) {
	return r.queries.CountUserWatchedEpisodesBySeries(ctx, db.CountUserWatchedEpisodesBySeriesParams{
		UserID:   userID,
		SeriesID: seriesID,
	})
}

// Series Watch Progress

// GetSeriesWatchProgress retrieves a user's watch progress for a series.
func (r *pgRepository) GetSeriesWatchProgress(ctx context.Context, userID, seriesID uuid.UUID) (*SeriesWatchProgress, error) {
	progress, err := r.queries.GetSeriesWatchProgress(ctx, db.GetSeriesWatchProgressParams{
		UserID:   userID,
		SeriesID: seriesID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return fromDBSeriesWatchProgress(&progress), nil
}

// ListContinueWatchingSeries lists series a user is currently watching.
func (r *pgRepository) ListContinueWatchingSeries(ctx context.Context, userID uuid.UUID, limit int) ([]*SeriesWatchProgress, error) {
	rows, err := r.queries.ListContinueWatchingSeries(ctx, db.ListContinueWatchingSeriesParams{
		UserID: userID,
		Limit:  int32(limit),
	})
	if err != nil {
		return nil, err
	}

	progress := make([]*SeriesWatchProgress, len(rows))
	for i, row := range rows {
		progress[i] = &SeriesWatchProgress{
			UserID:          row.UserID,
			SeriesID:        row.SeriesID,
			TotalEpisodes:   int(row.TotalEpisodes),
			WatchedEpisodes: int(row.WatchedEpisodes),
			IsWatching:      row.IsWatching,
		}
		if row.LastEpisodeID.Valid {
			id := uuid.UUID(row.LastEpisodeID.Bytes)
			progress[i].LastEpisodeID = &id
		}
		if row.LastSeasonNumber != nil {
			progress[i].LastSeasonNumber = int(*row.LastSeasonNumber)
		}
		if row.LastEpisodeNumber != nil {
			progress[i].LastEpisodeNumber = int(*row.LastEpisodeNumber)
		}
		if row.ProgressPercentage.Valid {
			f, _ := row.ProgressPercentage.Float64Value()
			progress[i].ProgressPercent = f.Float64
		}
		if row.StartedAt.Valid {
			t := row.StartedAt.Time
			progress[i].StartedAt = &t
		}
		if row.LastWatchedAt.Valid {
			t := row.LastWatchedAt.Time
			progress[i].LastWatchedAt = &t
		}
		if row.CompletedAt.Valid {
			t := row.CompletedAt.Time
			progress[i].CompletedAt = &t
		}
	}
	return progress, nil
}

// ListCompletedSeries lists series a user has completed.
func (r *pgRepository) ListCompletedSeries(ctx context.Context, userID uuid.UUID, params ListParams) ([]*SeriesWatchProgress, error) {
	rows, err := r.queries.ListCompletedSeries(ctx, db.ListCompletedSeriesParams{
		UserID: userID,
		Limit:  int32(params.Limit),
		Offset: int32(params.Offset),
	})
	if err != nil {
		return nil, err
	}

	progress := make([]*SeriesWatchProgress, len(rows))
	for i, row := range rows {
		progress[i] = &SeriesWatchProgress{
			UserID:          row.UserID,
			SeriesID:        row.SeriesID,
			TotalEpisodes:   int(row.TotalEpisodes),
			WatchedEpisodes: int(row.WatchedEpisodes),
			IsWatching:      row.IsWatching,
		}
		if row.LastEpisodeID.Valid {
			id := uuid.UUID(row.LastEpisodeID.Bytes)
			progress[i].LastEpisodeID = &id
		}
		if row.LastSeasonNumber != nil {
			progress[i].LastSeasonNumber = int(*row.LastSeasonNumber)
		}
		if row.LastEpisodeNumber != nil {
			progress[i].LastEpisodeNumber = int(*row.LastEpisodeNumber)
		}
		if row.ProgressPercentage.Valid {
			f, _ := row.ProgressPercentage.Float64Value()
			progress[i].ProgressPercent = f.Float64
		}
		if row.StartedAt.Valid {
			t := row.StartedAt.Time
			progress[i].StartedAt = &t
		}
		if row.LastWatchedAt.Valid {
			t := row.LastWatchedAt.Time
			progress[i].LastWatchedAt = &t
		}
		if row.CompletedAt.Valid {
			t := row.CompletedAt.Time
			progress[i].CompletedAt = &t
		}
	}
	return progress, nil
}

// DeleteSeriesWatchProgress deletes a user's watch progress for a series.
func (r *pgRepository) DeleteSeriesWatchProgress(ctx context.Context, userID, seriesID uuid.UUID) error {
	return r.queries.DeleteSeriesWatchProgress(ctx, db.DeleteSeriesWatchProgressParams{
		UserID:   userID,
		SeriesID: seriesID,
	})
}

// External Ratings

// GetSeriesExternalRatings retrieves external ratings for a series.
func (r *pgRepository) GetSeriesExternalRatings(ctx context.Context, seriesID uuid.UUID) (map[string]float64, error) {
	rows, err := r.queries.GetSeriesExternalRatings(ctx, seriesID)
	if err != nil {
		return nil, err
	}

	ratings := make(map[string]float64, len(rows))
	for _, row := range rows {
		if row.Rating.Valid {
			f, _ := row.Rating.Float64Value()
			ratings[row.Source] = f.Float64
		}
	}
	return ratings, nil
}

// UpsertSeriesExternalRating upserts an external rating for a series.
func (r *pgRepository) UpsertSeriesExternalRating(ctx context.Context, seriesID uuid.UUID, source string, rating float64, voteCount int, certified bool) error {
	return r.queries.UpsertSeriesExternalRating(ctx, db.UpsertSeriesExternalRatingParams{
		SeriesID:  seriesID,
		Source:    source,
		Rating:    numericFromFloat(rating),
		VoteCount: intPtr(int32(voteCount)),
		Certified: &certified,
	})
}

// Helper functions

func fromDBEpisodeWatchHistory(h *db.EpisodeWatchHistory) *EpisodeWatchHistory {
	if h == nil {
		return nil
	}
	result := &EpisodeWatchHistory{
		ID:            h.ID,
		UserID:        h.UserID,
		EpisodeID:     h.EpisodeID,
		PositionTicks: h.PositionTicks,
		Completed:     h.Completed,
		StartedAt:     h.StartedAt,
		LastUpdatedAt: h.LastUpdatedAt,
	}
	if h.ProfileID.Valid {
		pid := uuid.UUID(h.ProfileID.Bytes)
		result.ProfileID = &pid
	}
	if h.DurationTicks != nil {
		result.DurationTicks = *h.DurationTicks
	}
	if h.PlayedPercentage.Valid {
		f, _ := h.PlayedPercentage.Float64Value()
		result.PlayedPercentage = f.Float64
	}
	if h.CompletedAt.Valid {
		t := h.CompletedAt.Time
		result.CompletedAt = &t
	}
	if h.DeviceName != nil {
		result.DeviceName = *h.DeviceName
	}
	if h.DeviceType != nil {
		result.DeviceType = *h.DeviceType
	}
	if h.ClientName != nil {
		result.ClientName = *h.ClientName
	}
	if h.PlayMethod != nil {
		result.PlayMethod = *h.PlayMethod
	}
	return result
}

func fromDBSeriesWatchProgress(p *db.SeriesWatchProgress) *SeriesWatchProgress {
	if p == nil {
		return nil
	}
	result := &SeriesWatchProgress{
		UserID:          p.UserID,
		SeriesID:        p.SeriesID,
		TotalEpisodes:   int(p.TotalEpisodes),
		WatchedEpisodes: int(p.WatchedEpisodes),
		IsWatching:      p.IsWatching,
	}
	if p.LastEpisodeID.Valid {
		id := uuid.UUID(p.LastEpisodeID.Bytes)
		result.LastEpisodeID = &id
	}
	if p.LastSeasonNumber != nil {
		result.LastSeasonNumber = int(*p.LastSeasonNumber)
	}
	if p.LastEpisodeNumber != nil {
		result.LastEpisodeNumber = int(*p.LastEpisodeNumber)
	}
	if p.ProgressPercentage.Valid {
		f, _ := p.ProgressPercentage.Float64Value()
		result.ProgressPercent = f.Float64
	}
	if p.StartedAt.Valid {
		t := p.StartedAt.Time
		result.StartedAt = &t
	}
	if p.LastWatchedAt.Valid {
		t := p.LastWatchedAt.Time
		result.LastWatchedAt = &t
	}
	if p.CompletedAt.Valid {
		t := p.CompletedAt.Time
		result.CompletedAt = &t
	}
	return result
}

func intPtr(i int32) *int32 {
	return &i
}
