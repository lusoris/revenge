package api

import (
	"context"
	"errors"

	"github.com/google/uuid"

	gen "github.com/lusoris/revenge/api/generated"
	"github.com/lusoris/revenge/internal/content/movie"
)

func movieModuleDisabled() *gen.Error {
	return &gen.Error{
		Code:    "module_disabled",
		Message: "Movie module disabled",
	}
}

// ListMovies implements the listMovies operation.
func (h *Handler) ListMovies(ctx context.Context, params gen.ListMoviesParams) (gen.ListMoviesRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireMovieService()
	if err != nil {
		return movieModuleDisabled(), nil
	}

	listParams := movie.ListParams{
		Limit:     params.Limit.Or(20),
		Offset:    params.Offset.Or(0),
		SortBy:    string(params.SortBy.Or(gen.ListMoviesSortByTitle)),
		SortOrder: string(params.SortOrder.Or(gen.ListMoviesSortOrderAsc)),
	}

	var movies []*movie.Movie
	var total int64

	if params.LibraryId.IsSet() {
		movies, total, err = svc.ListMovies(ctx, params.LibraryId.Value, listParams)
	} else if params.Query.IsSet() {
		movies, err = svc.SearchMovies(ctx, params.Query.Value, listParams)
		if err == nil {
			total = int64(len(movies))
		}
	} else {
		movies, total, err = svc.ListAllMovies(ctx, listParams)
	}

	if err != nil {
		h.logger.Error("List movies failed", "error", err, "user_id", usr.ID)
		return &gen.Error{
			Code:    "list_failed",
			Message: "Failed to list movies",
		}, nil
	}

	result := gen.MovieListResponse{
		Movies:     make([]gen.Movie, len(movies)),
		Pagination: paginationMeta(total, listParams.Limit, listParams.Offset),
	}

	for i, m := range movies {
		result.Movies[i] = movieToAPI(m, nil)
	}

	return &result, nil
}

// ListRecentMovies implements the listRecentMovies operation.
func (h *Handler) ListRecentMovies(ctx context.Context, params gen.ListRecentMoviesParams) (gen.ListRecentMoviesRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireMovieService()
	if err != nil {
		return movieModuleDisabled(), nil
	}

	// Get accessible library IDs for user
	libs, err := h.libraryService.ListForUser(ctx, usr.ID)
	if err != nil {
		h.logger.Error("Get libraries failed", "error", err)
		return &gen.Error{
			Code:    "list_failed",
			Message: "Failed to get libraries",
		}, nil
	}

	libraryIDs := make([]uuid.UUID, len(libs))
	for i, l := range libs {
		libraryIDs[i] = l.Library.ID
	}

	movies, err := svc.ListRecentlyAdded(ctx, libraryIDs, params.Limit.Or(20))
	if err != nil {
		h.logger.Error("List recent movies failed", "error", err)
		return &gen.Error{
			Code:    "list_failed",
			Message: "Failed to list movies",
		}, nil
	}

	result := make(gen.ListRecentMoviesOKApplicationJSON, len(movies))
	for i, m := range movies {
		result[i] = movieToAPI(m, nil)
	}

	return &result, nil
}

// ListContinueWatching implements the listContinueWatching operation.
func (h *Handler) ListContinueWatching(ctx context.Context, params gen.ListContinueWatchingParams) (gen.ListContinueWatchingRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireMovieService()
	if err != nil {
		return movieModuleDisabled(), nil
	}

	history, err := svc.ListResumeableMovies(ctx, usr.ID, params.Limit.Or(10))
	if err != nil {
		h.logger.Error("List continue watching failed", "error", err)
		return &gen.Error{
			Code:    "list_failed",
			Message: "Failed to list movies",
		}, nil
	}

	result := make(gen.ListContinueWatchingOKApplicationJSON, 0, len(history))
	for _, hist := range history {
		m, err := svc.GetMovie(ctx, hist.MovieID)
		if err != nil {
			continue
		}
		result = append(result, movieWithProgressToAPI(m, &hist))
	}

	return &result, nil
}

// GetMovie implements the getMovie operation.
func (h *Handler) GetMovie(ctx context.Context, params gen.GetMovieParams) (gen.GetMovieRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.GetMovieUnauthorized{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireMovieService()
	if err != nil {
		return movieModuleDisabled(), nil
	}

	m, err := svc.GetMovieWithRelations(ctx, params.MovieId)
	if err != nil {
		if errors.Is(err, movie.ErrMovieNotFoundInService) {
			return &gen.GetMovieNotFound{
				Code:    "not_found",
				Message: "Movie not found",
			}, nil
		}
		h.logger.Error("Get movie failed", "error", err, "movie_id", params.MovieId)
		return &gen.GetMovieNotFound{
			Code:    "get_failed",
			Message: "Failed to get movie",
		}, nil
	}

	// Get user data
	userData := h.getUserMovieData(ctx, usr.ID, m.ID)

	result := movieFullToAPI(m, userData)
	return &result, nil
}

// UpdateMovie implements the updateMovie operation.
func (h *Handler) UpdateMovie(ctx context.Context, req *gen.MovieUpdate, params gen.UpdateMovieParams) (gen.UpdateMovieRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.UpdateMovieUnauthorized{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireMovieService()
	if err != nil {
		return movieModuleDisabled(), nil
	}

	// Only admins can update movie metadata
	if !usr.IsAdmin {
		return &gen.UpdateMovieForbidden{
			Code:    "forbidden",
			Message: "Admin access required",
		}, nil
	}

	m, err := svc.GetMovie(ctx, params.MovieId)
	if err != nil {
		if errors.Is(err, movie.ErrMovieNotFoundInService) {
			return &gen.UpdateMovieNotFound{
				Code:    "not_found",
				Message: "Movie not found",
			}, nil
		}
		return &gen.UpdateMovieNotFound{
			Code:    "get_failed",
			Message: "Failed to get movie",
		}, nil
	}

	// Apply updates
	if req.Title.IsSet() {
		m.Title = req.Title.Value
	}
	if req.OriginalTitle.IsSet() {
		m.OriginalTitle = req.OriginalTitle.Value
	}
	if req.SortTitle.IsSet() {
		m.SortTitle = req.SortTitle.Value
	}
	if req.Overview.IsSet() {
		m.Overview = req.Overview.Value
	}
	if req.Tagline.IsSet() {
		m.Tagline = req.Tagline.Value
	}
	if req.ContentRating.IsSet() {
		m.ContentRating = req.ContentRating.Value
	}

	if err := svc.UpdateMovie(ctx, m); err != nil {
		h.logger.Error("Update movie failed", "error", err, "movie_id", params.MovieId)
		return &gen.UpdateMovieNotFound{
			Code:    "update_failed",
			Message: "Failed to update movie",
		}, nil
	}

	result := movieToAPI(m, nil)
	return &result, nil
}

// DeleteMovie implements the deleteMovie operation.
func (h *Handler) DeleteMovie(ctx context.Context, params gen.DeleteMovieParams) (gen.DeleteMovieRes, error) {
	svc, err := h.requireMovieService()
	if err != nil {
		return movieModuleDisabled(), nil
	}

	_, err := requireAdmin(ctx)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) {
			return &gen.DeleteMovieUnauthorized{
				Code:    "unauthorized",
				Message: "Not authenticated",
			}, nil
		}
		return &gen.DeleteMovieForbidden{
			Code:    "forbidden",
			Message: "Admin access required",
		}, nil
	}

	if err := svc.DeleteMovie(ctx, params.MovieId); err != nil {
		if errors.Is(err, movie.ErrMovieNotFoundInService) {
			return &gen.DeleteMovieNotFound{
				Code:    "not_found",
				Message: "Movie not found",
			}, nil
		}
		h.logger.Error("Delete movie failed", "error", err, "movie_id", params.MovieId)
		return &gen.DeleteMovieNotFound{
			Code:    "delete_failed",
			Message: "Failed to delete movie",
		}, nil
	}

	return &gen.DeleteMovieNoContent{}, nil
}

// RefreshMovieMetadata implements the refreshMovieMetadata operation.
func (h *Handler) RefreshMovieMetadata(ctx context.Context, params gen.RefreshMovieMetadataParams) (gen.RefreshMovieMetadataRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.RefreshMovieMetadataUnauthorized{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireMovieService()
	if err != nil {
		return movieModuleDisabled(), nil
	}

	m, err := svc.GetMovie(ctx, params.MovieId)
	if err != nil {
		if errors.Is(err, movie.ErrMovieNotFoundInService) {
			return &gen.RefreshMovieMetadataNotFound{
				Code:    "not_found",
				Message: "Movie not found",
			}, nil
		}
		return &gen.RefreshMovieMetadataNotFound{
			Code:    "get_failed",
			Message: "Failed to get movie",
		}, nil
	}

	// Queue metadata enrichment job
	if h.riverClient != nil {
		_, err = h.riverClient.Insert(ctx, movie.EnrichMetadataArgs{
			MovieID: m.ID,
			TmdbID:  m.TmdbID,
			ImdbID:  m.ImdbID,
			Title:   m.Title,
			Year:    m.Year,
		}, nil)
		if err != nil {
			h.logger.Error("Queue metadata refresh failed", "error", err, "movie_id", m.ID, "user_id", usr.ID)
		}
	}

	return &gen.JobResponse{
		JobId:   uuid.New().String(),
		Status:  gen.JobResponseStatusQueued,
		Message: gen.NewOptString("Metadata refresh queued"),
	}, nil
}

// AddMovieToFavorites implements the addMovieToFavorites operation.
func (h *Handler) AddMovieToFavorites(ctx context.Context, params gen.AddMovieToFavoritesParams) (gen.AddMovieToFavoritesRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.AddMovieToFavoritesUnauthorized{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	// Verify movie exists
	svc, err := h.requireMovieService()
	if err != nil {
		return movieModuleDisabled(), nil
	}

	_, err = svc.GetMovie(ctx, params.MovieId)
	if err != nil {
		if errors.Is(err, movie.ErrMovieNotFoundInService) {
			return &gen.AddMovieToFavoritesNotFound{
				Code:    "not_found",
				Message: "Movie not found",
			}, nil
		}
		return &gen.AddMovieToFavoritesNotFound{
			Code:    "get_failed",
			Message: "Failed to get movie",
		}, nil
	}

	if err := svc.AddFavorite(ctx, usr.ID, params.MovieId); err != nil {
		h.logger.Error("Add favorite failed", "error", err, "movie_id", params.MovieId, "user_id", usr.ID)
		return &gen.AddMovieToFavoritesNotFound{
			Code:    "add_failed",
			Message: "Failed to add to favorites",
		}, nil
	}

	return &gen.AddMovieToFavoritesNoContent{}, nil
}

// RemoveMovieFromFavorites implements the removeMovieFromFavorites operation.
func (h *Handler) RemoveMovieFromFavorites(ctx context.Context, params gen.RemoveMovieFromFavoritesParams) (gen.RemoveMovieFromFavoritesRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireMovieService()
	if err != nil {
		return movieModuleDisabled(), nil
	}

	if err := svc.RemoveFavorite(ctx, usr.ID, params.MovieId); err != nil {
		h.logger.Error("Remove favorite failed", "error", err, "movie_id", params.MovieId, "user_id", usr.ID)
	}

	return &gen.RemoveMovieFromFavoritesNoContent{}, nil
}

// AddMovieToWatchlist implements the addMovieToWatchlist operation.
func (h *Handler) AddMovieToWatchlist(ctx context.Context, params gen.AddMovieToWatchlistParams) (gen.AddMovieToWatchlistRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.AddMovieToWatchlistUnauthorized{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireMovieService()
	if err != nil {
		return movieModuleDisabled(), nil
	}

	// Verify movie exists
	_, err = svc.GetMovie(ctx, params.MovieId)
	if err != nil {
		if errors.Is(err, movie.ErrMovieNotFoundInService) {
			return &gen.AddMovieToWatchlistNotFound{
				Code:    "not_found",
				Message: "Movie not found",
			}, nil
		}
		return &gen.AddMovieToWatchlistNotFound{
			Code:    "get_failed",
			Message: "Failed to get movie",
		}, nil
	}

	if err := svc.AddToWatchlist(ctx, usr.ID, params.MovieId); err != nil {
		h.logger.Error("Add to watchlist failed", "error", err, "movie_id", params.MovieId, "user_id", usr.ID)
		return &gen.AddMovieToWatchlistNotFound{
			Code:    "add_failed",
			Message: "Failed to add to watchlist",
		}, nil
	}

	return &gen.AddMovieToWatchlistNoContent{}, nil
}

// RemoveMovieFromWatchlist implements the removeMovieFromWatchlist operation.
func (h *Handler) RemoveMovieFromWatchlist(ctx context.Context, params gen.RemoveMovieFromWatchlistParams) (gen.RemoveMovieFromWatchlistRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireMovieService()
	if err != nil {
		return movieModuleDisabled(), nil
	}

	if err := svc.RemoveFromWatchlist(ctx, usr.ID, params.MovieId); err != nil {
		h.logger.Error("Remove from watchlist failed", "error", err, "movie_id", params.MovieId, "user_id", usr.ID)
	}

	return &gen.RemoveMovieFromWatchlistNoContent{}, nil
}

// MarkMovieWatched implements the markMovieWatched operation.
func (h *Handler) MarkMovieWatched(ctx context.Context, params gen.MarkMovieWatchedParams) (gen.MarkMovieWatchedRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.MarkMovieWatchedUnauthorized{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireMovieService()
	if err != nil {
		return movieModuleDisabled(), nil
	}

	// Verify movie exists
	_, err = svc.GetMovie(ctx, params.MovieId)
	if err != nil {
		if errors.Is(err, movie.ErrMovieNotFoundInService) {
			return &gen.MarkMovieWatchedNotFound{
				Code:    "not_found",
				Message: "Movie not found",
			}, nil
		}
		return &gen.MarkMovieWatchedNotFound{
			Code:    "get_failed",
			Message: "Failed to get movie",
		}, nil
	}

	if err := svc.MarkAsWatched(ctx, usr.ID, params.MovieId); err != nil {
		h.logger.Error("Mark watched failed", "error", err, "movie_id", params.MovieId, "user_id", usr.ID)
		return &gen.MarkMovieWatchedNotFound{
			Code:    "mark_failed",
			Message: "Failed to mark as watched",
		}, nil
	}

	return &gen.MarkMovieWatchedNoContent{}, nil
}

// MarkMovieUnwatched implements the markMovieUnwatched operation.
func (h *Handler) MarkMovieUnwatched(ctx context.Context, params gen.MarkMovieUnwatchedParams) (gen.MarkMovieUnwatchedRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireMovieService()
	if err != nil {
		return movieModuleDisabled(), nil
	}

	if err := svc.MarkAsUnwatched(ctx, usr.ID, params.MovieId); err != nil {
		h.logger.Error("Mark unwatched failed", "error", err, "movie_id", params.MovieId, "user_id", usr.ID)
	}

	return &gen.MarkMovieUnwatchedNoContent{}, nil
}

// SetMovieRating implements the setMovieRating operation.
func (h *Handler) SetMovieRating(ctx context.Context, req *gen.UserRatingRequest, params gen.SetMovieRatingParams) (gen.SetMovieRatingRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.SetMovieRatingUnauthorized{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireMovieService()
	if err != nil {
		return movieModuleDisabled(), nil
	}

	// Verify movie exists
	_, err = svc.GetMovie(ctx, params.MovieId)
	if err != nil {
		if errors.Is(err, movie.ErrMovieNotFoundInService) {
			return &gen.SetMovieRatingNotFound{
				Code:    "not_found",
				Message: "Movie not found",
			}, nil
		}
		return &gen.SetMovieRatingNotFound{
			Code:    "get_failed",
			Message: "Failed to get movie",
		}, nil
	}

	review := ""
	if req.Review.IsSet() {
		review = req.Review.Value
	}

	if err := svc.SetUserRating(ctx, usr.ID, params.MovieId, req.Rating, review); err != nil {
		h.logger.Error("Set rating failed", "error", err, "movie_id", params.MovieId, "user_id", usr.ID)
		return &gen.SetMovieRatingNotFound{
			Code:    "set_failed",
			Message: "Failed to set rating",
		}, nil
	}

	rating, err := svc.GetUserRating(ctx, usr.ID, params.MovieId)
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

// DeleteMovieRating implements the deleteMovieRating operation.
func (h *Handler) DeleteMovieRating(ctx context.Context, params gen.DeleteMovieRatingParams) (gen.DeleteMovieRatingRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireMovieService()
	if err != nil {
		return movieModuleDisabled(), nil
	}

	if err := svc.DeleteUserRating(ctx, usr.ID, params.MovieId); err != nil {
		h.logger.Error("Delete rating failed", "error", err, "movie_id", params.MovieId, "user_id", usr.ID)
	}

	return &gen.DeleteMovieRatingNoContent{}, nil
}

// ListMyFavoriteMovies implements the listMyFavoriteMovies operation.
func (h *Handler) ListMyFavoriteMovies(ctx context.Context, params gen.ListMyFavoriteMoviesParams) (gen.ListMyFavoriteMoviesRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireMovieService()
	if err != nil {
		return movieModuleDisabled(), nil
	}

	listParams := movie.ListParams{
		Limit:  params.Limit.Or(20),
		Offset: params.Offset.Or(0),
	}

	movies, total, err := svc.ListFavorites(ctx, usr.ID, listParams)
	if err != nil {
		h.logger.Error("List favorites failed", "error", err, "user_id", usr.ID)
		return &gen.Error{
			Code:    "list_failed",
			Message: "Failed to list favorites",
		}, nil
	}

	result := gen.MovieListResponse{
		Movies:     make([]gen.Movie, len(movies)),
		Pagination: paginationMeta(total, listParams.Limit, listParams.Offset),
	}

	for i, m := range movies {
		result.Movies[i] = movieToAPI(m, nil)
	}

	return &result, nil
}

// ListMyWatchlist implements the listMyWatchlist operation.
func (h *Handler) ListMyWatchlist(ctx context.Context, params gen.ListMyWatchlistParams) (gen.ListMyWatchlistRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireMovieService()
	if err != nil {
		return movieModuleDisabled(), nil
	}

	listParams := movie.ListParams{
		Limit:  params.Limit.Or(20),
		Offset: params.Offset.Or(0),
	}

	movies, total, err := svc.ListWatchlist(ctx, usr.ID, listParams)
	if err != nil {
		h.logger.Error("List watchlist failed", "error", err, "user_id", usr.ID)
		return &gen.Error{
			Code:    "list_failed",
			Message: "Failed to list watchlist",
		}, nil
	}

	result := gen.MovieListResponse{
		Movies:     make([]gen.Movie, len(movies)),
		Pagination: paginationMeta(total, listParams.Limit, listParams.Offset),
	}

	for i, m := range movies {
		result.Movies[i] = movieToAPI(m, nil)
	}

	return &result, nil
}

// ListCollections implements the listCollections operation.
func (h *Handler) ListCollections(ctx context.Context, params gen.ListCollectionsParams) (gen.ListCollectionsRes, error) {
	_, err := requireUser(ctx)
	if err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireMovieService()
	if err != nil {
		return movieModuleDisabled(), nil
	}

	listParams := movie.ListParams{
		Limit:  params.Limit.Or(20),
		Offset: params.Offset.Or(0),
	}

	collections, total, err := svc.ListCollections(ctx, listParams)
	if err != nil {
		h.logger.Error("List collections failed", "error", err)
		return &gen.Error{
			Code:    "list_failed",
			Message: "Failed to list collections",
		}, nil
	}

	result := gen.CollectionListResponse{
		Collections: make([]gen.Collection, len(collections)),
		Pagination:  paginationMeta(total, listParams.Limit, listParams.Offset),
	}

	for i, c := range collections {
		result.Collections[i] = collectionToAPI(c)
	}

	return &result, nil
}

// GetCollection implements the getCollection operation.
func (h *Handler) GetCollection(ctx context.Context, params gen.GetCollectionParams) (gen.GetCollectionRes, error) {
	_, err := requireUser(ctx)
	if err != nil {
		return &gen.GetCollectionUnauthorized{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	svc, err := h.requireMovieService()
	if err != nil {
		return movieModuleDisabled(), nil
	}

	collection, err := svc.GetCollection(ctx, params.CollectionId)
	if err != nil {
		if errors.Is(err, movie.ErrCollectionNotFound) {
			return &gen.GetCollectionNotFound{
				Code:    "not_found",
				Message: "Collection not found",
			}, nil
		}
		h.logger.Error("Get collection failed", "error", err, "collection_id", params.CollectionId)
		return &gen.GetCollectionNotFound{
			Code:    "get_failed",
			Message: "Failed to get collection",
		}, nil
	}

	// Get movies in collection
	movies, err := svc.ListMoviesByCollection(ctx, params.CollectionId)
	if err != nil {
		h.logger.Error("List collection movies failed", "error", err, "collection_id", params.CollectionId)
	}

	result := collectionFullToAPI(collection, movies)
	return &result, nil
}

// Helper functions

func (h *Handler) getUserMovieData(ctx context.Context, userID, movieID uuid.UUID) *gen.MovieUserData {
	data := &gen.MovieUserData{}

	svc, err := h.requireMovieService()
	if err != nil {
		return data
	}

	if favorite, err := svc.IsFavorite(ctx, userID, movieID); err == nil {
		data.IsFavorite = gen.NewOptBool(favorite)
	}

	if watchlist, err := svc.IsInWatchlist(ctx, userID, movieID); err == nil {
		data.IsInWatchlist = gen.NewOptBool(watchlist)
	}

	if watched, err := svc.IsWatched(ctx, userID, movieID); err == nil {
		data.IsWatched = gen.NewOptBool(watched)
	}

	if rating, err := svc.GetUserRating(ctx, userID, movieID); err == nil && rating != nil {
		data.UserRating = gen.NewOptFloat64(rating.Rating)
	}

	if history, err := svc.GetWatchHistory(ctx, userID, movieID); err == nil && history != nil {
		data.PlaybackPosition = gen.NewOptInt64(history.PositionTicks)
		data.PlayedPercentage = gen.NewOptFloat64(history.PlayedPercentage)
		data.LastPlayedAt = gen.NewOptDateTime(history.LastUpdatedAt)
	}

	return data
}

func paginationMeta(total int64, limit, offset int) gen.PaginationMeta {
	return gen.PaginationMeta{
		Total:  int(total),
		Limit:  limit,
		Offset: offset,
	}
}
