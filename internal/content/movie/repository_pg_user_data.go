package movie

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/lusoris/revenge/internal/content/movie/db"
)

// User Ratings

// GetUserRating retrieves a user's rating for a movie.
func (r *pgRepository) GetUserRating(ctx context.Context, userID, movieID uuid.UUID) (*UserRating, error) {
	row, err := r.queries.GetMovieUserRating(ctx, db.GetMovieUserRatingParams{
		UserID:  userID,
		MovieID: movieID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	rating := &UserRating{
		UserID:    row.UserID,
		MovieID:   row.MovieID,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
	if row.Rating.Valid {
		f, _ := row.Rating.Float64Value()
		rating.Rating = f.Float64
	}
	if row.Review != nil {
		rating.Review = *row.Review
	}
	return rating, nil
}

// SetUserRating sets a user's rating for a movie.
func (r *pgRepository) SetUserRating(ctx context.Context, userID, movieID uuid.UUID, rating float64, review string) error {
	_, err := r.queries.SetMovieUserRating(ctx, db.SetMovieUserRatingParams{
		UserID:  userID,
		MovieID: movieID,
		Rating:  numericFromFloat(rating),
		Review:  ptrString(review),
	})
	return err
}

// DeleteUserRating deletes a user's rating for a movie.
func (r *pgRepository) DeleteUserRating(ctx context.Context, userID, movieID uuid.UUID) error {
	return r.queries.DeleteMovieUserRating(ctx, db.DeleteMovieUserRatingParams{
		UserID:  userID,
		MovieID: movieID,
	})
}

// Favorites

// IsFavorite checks if a movie is in a user's favorites.
func (r *pgRepository) IsFavorite(ctx context.Context, userID, movieID uuid.UUID) (bool, error) {
	return r.queries.IsMovieFavorite(ctx, db.IsMovieFavoriteParams{
		UserID:  userID,
		MovieID: movieID,
	})
}

// AddFavorite adds a movie to a user's favorites.
func (r *pgRepository) AddFavorite(ctx context.Context, userID, movieID uuid.UUID) error {
	return r.queries.AddMovieFavorite(ctx, db.AddMovieFavoriteParams{
		UserID:  userID,
		MovieID: movieID,
	})
}

// RemoveFavorite removes a movie from a user's favorites.
func (r *pgRepository) RemoveFavorite(ctx context.Context, userID, movieID uuid.UUID) error {
	return r.queries.RemoveMovieFavorite(ctx, db.RemoveMovieFavoriteParams{
		UserID:  userID,
		MovieID: movieID,
	})
}

// ListFavorites retrieves a user's favorite movies.
func (r *pgRepository) ListFavorites(ctx context.Context, userID uuid.UUID, params ListParams) ([]*Movie, error) {
	rows, err := r.queries.ListUserFavoriteMovies(ctx, db.ListUserFavoriteMoviesParams{
		UserID: userID,
		Limit:  int32(params.Limit),
		Offset: int32(params.Offset),
	})
	if err != nil {
		return nil, err
	}

	movies := make([]*Movie, len(rows))
	for i, row := range rows {
		movies[i] = FromDBMovie(&row)
	}
	return movies, nil
}

// CountFavorites returns the number of favorite movies for a user.
func (r *pgRepository) CountFavorites(ctx context.Context, userID uuid.UUID) (int64, error) {
	return r.queries.CountUserFavoriteMovies(ctx, userID)
}

// Watch History

// GetWatchHistory retrieves the current watch history for a movie.
func (r *pgRepository) GetWatchHistory(ctx context.Context, userID, movieID uuid.UUID) (*WatchHistory, error) {
	row, err := r.queries.GetMovieWatchHistory(ctx, db.GetMovieWatchHistoryParams{
		UserID:  userID,
		MovieID: movieID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return fromDBWatchHistory(&row), nil
}

// CreateWatchHistory creates a new watch history entry.
func (r *pgRepository) CreateWatchHistory(ctx context.Context, history *WatchHistory) error {
	var profileID pgtype.UUID
	if history.ProfileID != nil {
		profileID = pgtype.UUID{Bytes: *history.ProfileID, Valid: true}
	}

	row, err := r.queries.CreateWatchHistory(ctx, db.CreateWatchHistoryParams{
		UserID:        history.UserID,
		ProfileID:     profileID,
		MovieID:       history.MovieID,
		PositionTicks: history.PositionTicks,
		DurationTicks: &history.DurationTicks,
		DeviceName:    ptrString(history.DeviceName),
		DeviceType:    ptrString(history.DeviceType),
		ClientName:    ptrString(history.ClientName),
		PlayMethod:    ptrString(history.PlayMethod),
	})
	if err != nil {
		return err
	}
	history.ID = row.ID
	history.StartedAt = row.StartedAt
	history.LastUpdatedAt = row.LastUpdatedAt
	return nil
}

// UpdateWatchHistory updates a watch history entry.
func (r *pgRepository) UpdateWatchHistory(ctx context.Context, id uuid.UUID, positionTicks int64, durationTicks *int64) error {
	_, err := r.queries.UpdateWatchHistory(ctx, db.UpdateWatchHistoryParams{
		ID:            id,
		PositionTicks: positionTicks,
		DurationTicks: durationTicks,
	})
	return err
}

// MarkWatchHistoryCompleted marks a watch history entry as completed.
func (r *pgRepository) MarkWatchHistoryCompleted(ctx context.Context, id uuid.UUID) error {
	return r.queries.MarkWatchHistoryCompleted(ctx, id)
}

// ListResumeableMovies retrieves movies that the user can resume.
func (r *pgRepository) ListResumeableMovies(ctx context.Context, userID uuid.UUID, limit int) ([]WatchHistory, error) {
	rows, err := r.queries.ListResumeableMovies(ctx, db.ListResumeableMoviesParams{
		UserID: userID,
		Limit:  int32(limit),
	})
	if err != nil {
		return nil, err
	}

	history := make([]WatchHistory, len(rows))
	for i, row := range rows {
		history[i] = WatchHistory{
			ID:            row.ID,
			UserID:        row.UserID,
			MovieID:       row.MovieID,
			PositionTicks: row.PositionTicks,
			Completed:     row.Completed,
			StartedAt:     row.StartedAt,
			LastUpdatedAt: row.LastUpdatedAt,
		}
		if row.ProfileID.Valid {
			id := uuid.UUID(row.ProfileID.Bytes)
			history[i].ProfileID = &id
		}
		if row.DurationTicks != nil {
			history[i].DurationTicks = *row.DurationTicks
		}
		if row.PlayedPercentage.Valid {
			f, _ := row.PlayedPercentage.Float64Value()
			history[i].PlayedPercentage = f.Float64
		}
	}
	return history, nil
}

// IsWatched checks if a movie has been watched by a user.
func (r *pgRepository) IsWatched(ctx context.Context, userID, movieID uuid.UUID) (bool, error) {
	return r.queries.IsMovieWatched(ctx, db.IsMovieWatchedParams{
		UserID:  userID,
		MovieID: movieID,
	})
}

// Watchlist

// IsInWatchlist checks if a movie is in a user's watchlist.
func (r *pgRepository) IsInWatchlist(ctx context.Context, userID, movieID uuid.UUID) (bool, error) {
	return r.queries.IsMovieInWatchlist(ctx, db.IsMovieInWatchlistParams{
		UserID:  userID,
		MovieID: movieID,
	})
}

// AddToWatchlist adds a movie to a user's watchlist.
func (r *pgRepository) AddToWatchlist(ctx context.Context, userID, movieID uuid.UUID) error {
	return r.queries.AddMovieToWatchlist(ctx, db.AddMovieToWatchlistParams{
		UserID:  userID,
		MovieID: movieID,
	})
}

// RemoveFromWatchlist removes a movie from a user's watchlist.
func (r *pgRepository) RemoveFromWatchlist(ctx context.Context, userID, movieID uuid.UUID) error {
	return r.queries.RemoveMovieFromWatchlist(ctx, db.RemoveMovieFromWatchlistParams{
		UserID:  userID,
		MovieID: movieID,
	})
}

// ListWatchlist retrieves a user's watchlist.
func (r *pgRepository) ListWatchlist(ctx context.Context, userID uuid.UUID, params ListParams) ([]*Movie, error) {
	rows, err := r.queries.ListUserWatchlist(ctx, db.ListUserWatchlistParams{
		UserID: userID,
		Limit:  int32(params.Limit),
		Offset: int32(params.Offset),
	})
	if err != nil {
		return nil, err
	}

	movies := make([]*Movie, len(rows))
	for i, row := range rows {
		movies[i] = FromDBMovie(&row)
	}
	return movies, nil
}

// CountWatchlist returns the number of movies in a user's watchlist.
func (r *pgRepository) CountWatchlist(ctx context.Context, userID uuid.UUID) (int64, error) {
	return r.queries.CountUserWatchlist(ctx, userID)
}

// DeleteWatchHistory removes a watch history entry.
func (r *pgRepository) DeleteWatchHistory(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteWatchHistory(ctx, id)
}

// Helper functions

func fromDBWatchHistory(row *db.MovieWatchHistory) *WatchHistory {
	if row == nil {
		return nil
	}
	history := &WatchHistory{
		ID:            row.ID,
		UserID:        row.UserID,
		MovieID:       row.MovieID,
		PositionTicks: row.PositionTicks,
		Completed:     row.Completed,
		StartedAt:     row.StartedAt,
		LastUpdatedAt: row.LastUpdatedAt,
	}
	if row.ProfileID.Valid {
		id := uuid.UUID(row.ProfileID.Bytes)
		history.ProfileID = &id
	}
	if row.DurationTicks != nil {
		history.DurationTicks = *row.DurationTicks
	}
	if row.PlayedPercentage.Valid {
		f, _ := row.PlayedPercentage.Float64Value()
		history.PlayedPercentage = f.Float64
	}
	if row.CompletedAt.Valid {
		t := row.CompletedAt.Time
		history.CompletedAt = &t
	}
	if row.DeviceName != nil {
		history.DeviceName = *row.DeviceName
	}
	if row.DeviceType != nil {
		history.DeviceType = *row.DeviceType
	}
	if row.ClientName != nil {
		history.ClientName = *row.ClientName
	}
	if row.PlayMethod != nil {
		history.PlayMethod = *row.PlayMethod
	}
	return history
}
