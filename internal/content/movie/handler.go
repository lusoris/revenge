package movie

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/content"
)

// MetadataQueue allows enqueuing metadata refresh jobs.
type MetadataQueue interface {
	EnqueueRefreshMovie(ctx context.Context, movieID uuid.UUID, force bool, languages []string) error
}

// Handler handles HTTP requests for movies
type Handler struct {
	service       Service
	metadataQueue MetadataQueue
}

// NewHandler creates a new movie handler
func NewHandler(service Service, metadataQueue MetadataQueue) *Handler {
	return &Handler{
		service:       service,
		metadataQueue: metadataQueue,
	}
}

// GetMovie handles GET /api/v1/movies/:id
func (h *Handler) GetMovie(ctx context.Context, id string) (*Movie, error) {
	movieID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid movie ID: %w", err)
	}

	return h.service.GetMovie(ctx, movieID)
}

// ListMovies handles GET /api/v1/movies
func (h *Handler) ListMovies(ctx context.Context, params ListMoviesParams) ([]Movie, error) {
	filters := ListFilters(params)

	return h.service.ListMovies(ctx, filters)
}

// CountMovies returns the total number of movies
func (h *Handler) CountMovies(ctx context.Context) (int64, error) {
	return h.service.CountMovies(ctx)
}

// SearchMovies handles GET /api/v1/movies/search
func (h *Handler) SearchMovies(ctx context.Context, params SearchMoviesParams) ([]Movie, error) {
	filters := SearchFilters{
		Limit:  params.Limit,
		Offset: params.Offset,
	}

	return h.service.SearchMovies(ctx, params.Query, filters)
}

// GetRecentlyAdded handles GET /api/v1/movies/recently-added
func (h *Handler) GetRecentlyAdded(ctx context.Context, params PaginationParams) ([]Movie, int64, error) {
	return h.service.ListRecentlyAdded(ctx, params.Limit, params.Offset)
}

// GetTopRated handles GET /api/v1/movies/top-rated
func (h *Handler) GetTopRated(ctx context.Context, params TopRatedParams) ([]Movie, int64, error) {
	minVotes := int32(100) // Default minimum votes
	if params.MinVotes != nil {
		minVotes = *params.MinVotes
	}

	return h.service.ListTopRated(ctx, minVotes, params.Limit, params.Offset)
}

// GetMovieFiles handles GET /api/v1/movies/:id/files
func (h *Handler) GetMovieFiles(ctx context.Context, id string) ([]MovieFile, error) {
	movieID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid movie ID: %w", err)
	}

	return h.service.GetMovieFiles(ctx, movieID)
}

// CreditPaginationParams contains pagination params for credit queries
type CreditPaginationParams struct {
	Limit  int32
	Offset int32
}

// GetMovieCast handles GET /api/v1/movies/:id/cast
func (h *Handler) GetMovieCast(ctx context.Context, id string, params CreditPaginationParams) ([]MovieCredit, int64, error) {
	movieID, err := uuid.Parse(id)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid movie ID: %w", err)
	}

	return h.service.GetMovieCast(ctx, movieID, params.Limit, params.Offset)
}

// GetMovieCrew handles GET /api/v1/movies/:id/crew
func (h *Handler) GetMovieCrew(ctx context.Context, id string, params CreditPaginationParams) ([]MovieCredit, int64, error) {
	movieID, err := uuid.Parse(id)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid movie ID: %w", err)
	}

	return h.service.GetMovieCrew(ctx, movieID, params.Limit, params.Offset)
}

// GetMovieGenres handles GET /api/v1/movies/:id/genres
func (h *Handler) GetMovieGenres(ctx context.Context, id string) ([]MovieGenre, error) {
	movieID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid movie ID: %w", err)
	}

	return h.service.GetMovieGenres(ctx, movieID)
}

// GetMoviesByGenre handles GET /api/v1/movies/genre/:genreId
func (h *Handler) GetMoviesByGenre(ctx context.Context, genreID int32, params PaginationParams) ([]Movie, error) {
	return h.service.GetMoviesByGenre(ctx, genreID, params.Limit, params.Offset)
}

// ListDistinctGenres returns all distinct movie genres with item counts.
func (h *Handler) ListDistinctGenres(ctx context.Context) ([]content.GenreSummary, error) {
	return h.service.ListDistinctGenres(ctx)
}

// GetMovieCollection handles GET /api/v1/movies/:id/collection
func (h *Handler) GetMovieCollection(ctx context.Context, id string) (*MovieCollection, error) {
	movieID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid movie ID: %w", err)
	}

	return h.service.GetCollectionForMovie(ctx, movieID)
}

// GetCollection handles GET /api/v1/collections/:id
func (h *Handler) GetCollection(ctx context.Context, id string) (*MovieCollection, error) {
	collectionID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid collection ID: %w", err)
	}

	return h.service.GetMovieCollection(ctx, collectionID)
}

// GetCollectionMovies handles GET /api/v1/collections/:id/movies
func (h *Handler) GetCollectionMovies(ctx context.Context, id string) ([]Movie, error) {
	collectionID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid collection ID: %w", err)
	}

	return h.service.GetMoviesByCollection(ctx, collectionID)
}

// UpdateWatchProgress handles POST /api/v1/movies/:id/progress
func (h *Handler) UpdateWatchProgress(ctx context.Context, userID uuid.UUID, id string, params UpdateWatchProgressParams) (*MovieWatched, error) {
	movieID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid movie ID: %w", err)
	}

	return h.service.UpdateWatchProgress(ctx, userID, movieID, params.ProgressSeconds, params.DurationSeconds)
}

// GetWatchProgress handles GET /api/v1/movies/:id/progress
func (h *Handler) GetWatchProgress(ctx context.Context, userID uuid.UUID, id string) (*MovieWatched, error) {
	movieID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid movie ID: %w", err)
	}

	return h.service.GetWatchProgress(ctx, userID, movieID)
}

// DeleteWatchProgress handles DELETE /api/v1/movies/:id/progress
func (h *Handler) DeleteWatchProgress(ctx context.Context, userID uuid.UUID, id string) error {
	movieID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid movie ID: %w", err)
	}

	return h.service.RemoveWatchProgress(ctx, userID, movieID)
}

// MarkAsWatched handles POST /api/v1/movies/:id/watched
func (h *Handler) MarkAsWatched(ctx context.Context, userID uuid.UUID, id string) error {
	movieID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid movie ID: %w", err)
	}

	return h.service.MarkAsWatched(ctx, userID, movieID)
}

// GetContinueWatching handles GET /api/v1/movies/continue-watching
func (h *Handler) GetContinueWatching(ctx context.Context, userID uuid.UUID, limit int32) ([]ContinueWatchingItem, error) {
	return h.service.GetContinueWatching(ctx, userID, limit)
}

// GetWatchHistory handles GET /api/v1/movies/watch-history
func (h *Handler) GetWatchHistory(ctx context.Context, userID uuid.UUID, params PaginationParams) ([]WatchedMovieItem, error) {
	return h.service.GetWatchHistory(ctx, userID, params.Limit, params.Offset)
}

// GetUserStats handles GET /api/v1/movies/stats
func (h *Handler) GetUserStats(ctx context.Context, userID uuid.UUID) (*UserMovieStats, error) {
	return h.service.GetUserStats(ctx, userID)
}

// RefreshMetadata handles POST /api/v1/movies/:id/refresh.
// It validates that the movie exists, then enqueues an async metadata refresh job.
func (h *Handler) RefreshMetadata(ctx context.Context, id string) error {
	movieID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid movie ID: %w", err)
	}

	// Verify movie exists before enqueuing
	if _, err := h.service.GetMovie(ctx, movieID); err != nil {
		return err
	}

	// Enqueue async refresh — returns immediately with 202 Accepted
	return h.metadataQueue.EnqueueRefreshMovie(ctx, movieID, true, nil)
}

// Request/Response parameter types

// ListMoviesParams contains parameters for listing movies
type ListMoviesParams struct {
	OrderBy string
	Limit   int32
	Offset  int32
}

// SearchMoviesParams contains parameters for searching movies
type SearchMoviesParams struct {
	Query  string
	Limit  int32
	Offset int32
}

// PaginationParams contains pagination parameters
type PaginationParams struct {
	Limit  int32
	Offset int32
}

// TopRatedParams contains parameters for top-rated movies
type TopRatedParams struct {
	MinVotes *int32
	Limit    int32
	Offset   int32
}

// UpdateWatchProgressParams contains parameters for updating watch progress
type UpdateWatchProgressParams struct {
	ProgressSeconds int32 `json:"progressSeconds"`
	DurationSeconds int32 `json:"durationSeconds"`
}

// HTTPError represents an HTTP error response
type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error implements the error interface
func (e *HTTPError) Error() string {
	return e.Message
}

// NewHTTPError creates a new HTTP error
func NewHTTPError(code int, message string) *HTTPError {
	return &HTTPError{
		Code:    code,
		Message: message,
	}
}

// NotFound creates a 404 error
func NotFound(message string) *HTTPError {
	return NewHTTPError(http.StatusNotFound, message)
}

// BadRequest creates a 400 error
func BadRequest(message string) *HTTPError {
	return NewHTTPError(http.StatusBadRequest, message)
}

// InternalError creates a 500 error
func InternalError(message string) *HTTPError {
	return NewHTTPError(http.StatusInternalServerError, message)
}
