package moviejobs

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/riverqueue/river"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	"github.com/lusoris/revenge/internal/content/movie"
	metadatajobs "github.com/lusoris/revenge/internal/service/metadata/jobs"
)

// MovieMetadataRefreshWorker refreshes movie metadata using the shared metadata service.
// It handles jobs of type metadatajobs.RefreshMovieArgs (kind: "metadata_refresh_movie").
type MovieMetadataRefreshWorker struct {
	river.WorkerDefaults[metadatajobs.RefreshMovieArgs]
	movieRepo        movie.Repository
	metadataProvider movie.MetadataProvider
	logger           *zap.Logger
}

// NewMovieMetadataRefreshWorker creates a new metadata refresh worker.
func NewMovieMetadataRefreshWorker(
	movieRepo movie.Repository,
	metadataProvider movie.MetadataProvider,
	logger *zap.Logger,
) *MovieMetadataRefreshWorker {
	return &MovieMetadataRefreshWorker{
		movieRepo:        movieRepo,
		metadataProvider: metadataProvider,
		logger:           logger,
	}
}

// Work executes the metadata refresh job.
func (w *MovieMetadataRefreshWorker) Work(ctx context.Context, job *river.Job[metadatajobs.RefreshMovieArgs]) error {
	args := job.Args

	w.logger.Info("starting movie metadata refresh",
		zap.String("job_id", fmt.Sprintf("%d", job.ID)),
		zap.String("movie_id", args.MovieID.String()),
		zap.Bool("force", args.Force),
		zap.Strings("languages", args.Languages),
	)

	// Get existing movie
	existingMovie, err := w.movieRepo.GetMovie(ctx, args.MovieID)
	if err != nil {
		w.logger.Error("failed to get movie",
			zap.String("movie_id", args.MovieID.String()),
			zap.Error(err),
		)
		return fmt.Errorf("failed to get movie: %w", err)
	}

	// Check if movie has TMDb ID
	if existingMovie.TMDbID == nil {
		w.logger.Warn("movie has no TMDb ID, cannot refresh",
			zap.String("movie_id", args.MovieID.String()),
		)
		return fmt.Errorf("movie has no TMDb ID")
	}

	tmdbID := int(*existingMovie.TMDbID)

	// Clear cache if force refresh
	if args.Force {
		w.metadataProvider.ClearCache()
		w.logger.Info("cleared metadata cache for force refresh",
			zap.String("job_id", fmt.Sprintf("%d", job.ID)),
		)
	}

	// Enrich movie with fresh metadata
	if err := w.metadataProvider.EnrichMovie(ctx, existingMovie); err != nil {
		w.logger.Error("failed to enrich movie with metadata",
			zap.Int32("tmdb_id", *existingMovie.TMDbID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to enrich movie: %w", err)
	}

	// Update movie in database with enriched metadata
	if err := w.updateMovieMetadata(ctx, existingMovie); err != nil {
		w.logger.Error("failed to update movie in database",
			zap.String("movie_id", args.MovieID.String()),
			zap.Error(err),
		)
		return fmt.Errorf("failed to update movie: %w", err)
	}

	// Refresh credits
	if err := w.refreshCredits(ctx, args.MovieID, tmdbID); err != nil {
		w.logger.Error("failed to refresh credits",
			zap.String("movie_id", args.MovieID.String()),
			zap.Error(err),
		)
		// Continue despite credit refresh failure
	}

	// Refresh genres
	if err := w.refreshGenres(ctx, args.MovieID, tmdbID); err != nil {
		w.logger.Error("failed to refresh genres",
			zap.String("movie_id", args.MovieID.String()),
			zap.Error(err),
		)
		// Continue despite genre refresh failure
	}

	w.logger.Info("movie metadata refresh completed",
		zap.String("movie_id", args.MovieID.String()),
		zap.String("title", existingMovie.Title),
	)

	return nil
}

// updateMovieMetadata updates the movie record in the database with enriched metadata.
func (w *MovieMetadataRefreshWorker) updateMovieMetadata(ctx context.Context, mov *movie.Movie) error {
	params := movie.UpdateMovieParams{
		ID:               mov.ID,
		TMDbID:           mov.TMDbID,
		IMDbID:           mov.IMDbID,
		Title:            &mov.Title,
		OriginalTitle:    mov.OriginalTitle,
		Year:             mov.Year,
		ReleaseDate:      formatTimePtr(mov.ReleaseDate),
		Runtime:          mov.Runtime,
		Overview:         mov.Overview,
		Tagline:          mov.Tagline,
		Status:           mov.Status,
		OriginalLanguage: mov.OriginalLanguage,
		PosterPath:       mov.PosterPath,
		BackdropPath:     mov.BackdropPath,
		VoteAverage:      formatDecimalPtr(mov.VoteAverage),
		VoteCount:        mov.VoteCount,
		Popularity:       formatDecimalPtr(mov.Popularity),
		Budget:           mov.Budget,
		Revenue:          mov.Revenue,
		TitlesI18n:       mov.TitlesI18n,
		TaglinesI18n:     mov.TaglinesI18n,
		OverviewsI18n:    mov.OverviewsI18n,
		AgeRatings:       mov.AgeRatings,
	}

	_, err := w.movieRepo.UpdateMovie(ctx, params)
	return err
}

// refreshCredits fetches fresh credits and updates the database.
func (w *MovieMetadataRefreshWorker) refreshCredits(ctx context.Context, movieID uuid.UUID, tmdbID int) error {
	// Fetch credits via metadata provider
	credits, err := w.metadataProvider.GetMovieCredits(ctx, movieID, tmdbID)
	if err != nil {
		return fmt.Errorf("failed to fetch credits: %w", err)
	}

	// Delete existing credits
	if err := w.movieRepo.DeleteMovieCredits(ctx, movieID); err != nil {
		return fmt.Errorf("failed to delete existing credits: %w", err)
	}

	// Create new credits
	for _, credit := range credits {
		params := movie.CreateMovieCreditParams{
			MovieID:      movieID,
			TMDbPersonID: credit.TMDbPersonID,
			Name:         credit.Name,
			CreditType:   credit.CreditType,
			Character:    credit.Character,
			Job:          credit.Job,
			Department:   credit.Department,
			CastOrder:    credit.CastOrder,
			ProfilePath:  credit.ProfilePath,
		}

		if _, err := w.movieRepo.CreateMovieCredit(ctx, params); err != nil {
			w.logger.Warn("failed to create credit",
				zap.String("movie_id", movieID.String()),
				zap.String("name", credit.Name),
				zap.Error(err),
			)
			// Continue with other credits
		}
	}

	w.logger.Info("refreshed credits",
		zap.String("movie_id", movieID.String()),
		zap.Int("credit_count", len(credits)),
	)

	return nil
}

// refreshGenres fetches fresh genres and updates the database.
func (w *MovieMetadataRefreshWorker) refreshGenres(ctx context.Context, movieID uuid.UUID, tmdbID int) error {
	// Fetch genres via metadata provider
	genres, err := w.metadataProvider.GetMovieGenres(ctx, movieID, tmdbID)
	if err != nil {
		return fmt.Errorf("failed to fetch genres: %w", err)
	}

	// Delete existing genres
	if err := w.movieRepo.DeleteMovieGenres(ctx, movieID); err != nil {
		return fmt.Errorf("failed to delete existing genres: %w", err)
	}

	// Add new genres
	for _, genre := range genres {
		if err := w.movieRepo.AddMovieGenre(ctx, movieID, genre.TMDbGenreID, genre.Name); err != nil {
			w.logger.Warn("failed to add genre",
				zap.String("movie_id", movieID.String()),
				zap.String("genre", genre.Name),
				zap.Error(err),
			)
			// Continue with other genres
		}
	}

	w.logger.Info("refreshed genres",
		zap.String("movie_id", movieID.String()),
		zap.Int("genre_count", len(genres)),
	)

	return nil
}

// formatTimePtr formats a time.Time pointer to a string pointer.
func formatTimePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format("2006-01-02")
	return &s
}

// formatDecimalPtr formats a decimal.Decimal pointer to a string pointer.
func formatDecimalPtr(d *decimal.Decimal) *string {
	if d == nil || d.IsZero() {
		return nil
	}
	s := d.String()
	return &s
}
