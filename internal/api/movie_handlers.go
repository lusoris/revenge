package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/content/movie"
	"github.com/lusoris/revenge/internal/service/metadata"
	"github.com/lusoris/revenge/internal/util"
)

// GetMovie delegates to the movie handler.
func (h *Handler) GetMovie(ctx context.Context, params ogen.GetMovieParams) (ogen.GetMovieRes, error) {
	m, err := h.movieHandler.GetMovie(ctx, params.ID.String())
	if err != nil {
		if errors.Is(err, movie.ErrMovieNotFound) {
			return &ogen.GetMovieNotFound{}, nil
		}
		return nil, err
	}

	// Get user's preferred language
	lang := h.GetMetadataLanguage(ctx)

	// Localize movie fields
	localizedMovie := LocalizeMovie(m, lang)

	return movieToOgen(localizedMovie), nil
}

// ListMovies delegates to the movie handler.
func (h *Handler) ListMovies(ctx context.Context, params ogen.ListMoviesParams) (ogen.ListMoviesRes, error) {
	handlerParams := movie.ListMoviesParams{
		OrderBy: string(params.OrderBy.Or("created_at")),
		Limit:   util.SafeIntToInt32(params.Limit.Or(20)),
		Offset:  util.SafeIntToInt32(params.Offset.Or(0)),
	}

	movies, err := h.movieHandler.ListMovies(ctx, handlerParams)
	if err != nil {
		return nil, err
	}

	// Get user's preferred language and localize results
	lang := h.GetMetadataLanguage(ctx)
	localizedMovies := LocalizeMovies(movies, lang)

	result := make([]ogen.Movie, len(localizedMovies))
	for i, m := range localizedMovies {
		result[i] = *movieToOgen(&m)
	}

	return (*ogen.ListMoviesOKApplicationJSON)(&result), nil
}

// SearchMovies delegates to the movie handler.
func (h *Handler) SearchMovies(ctx context.Context, params ogen.SearchMoviesParams) (ogen.SearchMoviesRes, error) {
	handlerParams := movie.SearchMoviesParams{
		Query:  params.Query,
		Limit:  util.SafeIntToInt32(params.Limit.Or(20)),
		Offset: util.SafeIntToInt32(params.Offset.Or(0)),
	}

	movies, err := h.movieHandler.SearchMovies(ctx, handlerParams)
	if err != nil {
		return nil, err
	}

	// Get user's preferred language and localize results
	lang := h.GetMetadataLanguage(ctx)
	localizedMovies := LocalizeMovies(movies, lang)

	result := make([]ogen.Movie, len(localizedMovies))
	for i, m := range localizedMovies {
		result[i] = *movieToOgen(&m)
	}

	return (*ogen.SearchMoviesOKApplicationJSON)(&result), nil
}

// GetRecentlyAdded delegates to the movie handler.
func (h *Handler) GetRecentlyAdded(ctx context.Context, params ogen.GetRecentlyAddedParams) (ogen.GetRecentlyAddedRes, error) {
	handlerParams := movie.PaginationParams{
		Limit:  util.SafeIntToInt32(params.Limit.Or(20)),
		Offset: 0,
	}

	movies, err := h.movieHandler.GetRecentlyAdded(ctx, handlerParams)
	if err != nil {
		return nil, err
	}

	// Get user's preferred language and localize results
	lang := h.GetMetadataLanguage(ctx)
	localizedMovies := LocalizeMovies(movies, lang)

	result := make([]ogen.Movie, len(localizedMovies))
	for i, m := range localizedMovies {
		result[i] = *movieToOgen(&m)
	}

	return (*ogen.GetRecentlyAddedOKApplicationJSON)(&result), nil
}

// GetTopRated delegates to the movie handler.
func (h *Handler) GetTopRated(ctx context.Context, params ogen.GetTopRatedParams) (ogen.GetTopRatedRes, error) {
	handlerParams := movie.TopRatedParams{
		Limit:    util.SafeIntToInt32(params.Limit.Or(20)),
		Offset:   0,
		MinVotes: func() *int32 { v := util.SafeIntToInt32(params.MinVotes.Or(100)); return &v }(),
	}

	movies, err := h.movieHandler.GetTopRated(ctx, handlerParams)
	if err != nil {
		return nil, err
	}

	// Get user's preferred language and localize results
	lang := h.GetMetadataLanguage(ctx)
	localizedMovies := LocalizeMovies(movies, lang)

	result := make([]ogen.Movie, len(localizedMovies))
	for i, m := range localizedMovies {
		result[i] = *movieToOgen(&m)
	}

	return (*ogen.GetTopRatedOKApplicationJSON)(&result), nil
}

// GetContinueWatching delegates to the movie handler.
func (h *Handler) GetContinueWatching(ctx context.Context, params ogen.GetContinueWatchingParams) (ogen.GetContinueWatchingRes, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, err
	}
	limit := util.SafeIntToInt32(params.Limit.Or(20))

	items, err := h.movieHandler.GetContinueWatching(ctx, userID, limit)
	if err != nil {
		return nil, err
	}

	// Get user's preferred language and localize results
	lang := h.GetMetadataLanguage(ctx)
	localizedItems := LocalizeContinueWatchingItems(items, lang)

	result := make([]ogen.ContinueWatchingItem, len(localizedItems))
	for i, item := range localizedItems {
		result[i] = continueWatchingItemToOgen(&item)
	}

	return (*ogen.GetContinueWatchingOKApplicationJSON)(&result), nil
}

// GetWatchHistory delegates to the movie handler.
func (h *Handler) GetWatchHistory(ctx context.Context, params ogen.GetWatchHistoryParams) (ogen.GetWatchHistoryRes, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, err
	}
	handlerParams := movie.PaginationParams{
		Limit:  util.SafeIntToInt32(params.Limit.Or(20)),
		Offset: util.SafeIntToInt32(params.Offset.Or(0)),
	}

	items, err := h.movieHandler.GetWatchHistory(ctx, userID, handlerParams)
	if err != nil {
		return nil, err
	}

	// Get user's preferred language and localize results
	lang := h.GetMetadataLanguage(ctx)
	localizedItems := LocalizeWatchedMovieItems(items, lang)

	result := make([]ogen.WatchedMovieItem, len(localizedItems))
	for i, item := range localizedItems {
		result[i] = watchedMovieItemToOgen(&item)
	}

	return (*ogen.GetWatchHistoryOKApplicationJSON)(&result), nil
}

// GetUserMovieStats delegates to the movie handler.
func (h *Handler) GetUserMovieStats(ctx context.Context) (ogen.GetUserMovieStatsRes, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	stats, err := h.movieHandler.GetUserStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	return userMovieStatsToOgen(stats), nil
}

// GetMovieFiles delegates to the movie handler.
func (h *Handler) GetMovieFiles(ctx context.Context, params ogen.GetMovieFilesParams) (ogen.GetMovieFilesRes, error) {
	files, err := h.movieHandler.GetMovieFiles(ctx, params.ID.String())
	if err != nil {
		if errors.Is(err, movie.ErrMovieNotFound) {
			return &ogen.GetMovieFilesNotFound{}, nil
		}
		return nil, err
	}

	result := make([]ogen.MovieFile, len(files))
	for i, f := range files {
		result[i] = *movieFileToOgen(&f)
	}

	return (*ogen.GetMovieFilesOKApplicationJSON)(&result), nil
}

// GetMovieCast delegates to the movie handler.
func (h *Handler) GetMovieCast(ctx context.Context, params ogen.GetMovieCastParams) (ogen.GetMovieCastRes, error) {
	cast, err := h.movieHandler.GetMovieCast(ctx, params.ID.String())
	if err != nil {
		if errors.Is(err, movie.ErrMovieNotFound) {
			return &ogen.GetMovieCastNotFound{}, nil
		}
		return nil, err
	}

	result := make([]ogen.MovieCredit, len(cast))
	for i, c := range cast {
		result[i] = *movieCreditToOgen(&c)
	}

	return (*ogen.GetMovieCastOKApplicationJSON)(&result), nil
}

// GetMovieCrew delegates to the movie handler.
func (h *Handler) GetMovieCrew(ctx context.Context, params ogen.GetMovieCrewParams) (ogen.GetMovieCrewRes, error) {
	crew, err := h.movieHandler.GetMovieCrew(ctx, params.ID.String())
	if err != nil {
		if errors.Is(err, movie.ErrMovieNotFound) {
			return &ogen.GetMovieCrewNotFound{}, nil
		}
		return nil, err
	}

	result := make([]ogen.MovieCredit, len(crew))
	for i, c := range crew {
		result[i] = *movieCreditToOgen(&c)
	}

	return (*ogen.GetMovieCrewOKApplicationJSON)(&result), nil
}

// GetMovieGenres delegates to the movie handler.
func (h *Handler) GetMovieGenres(ctx context.Context, params ogen.GetMovieGenresParams) (ogen.GetMovieGenresRes, error) {
	genres, err := h.movieHandler.GetMovieGenres(ctx, params.ID.String())
	if err != nil {
		if errors.Is(err, movie.ErrMovieNotFound) {
			return &ogen.GetMovieGenresNotFound{}, nil
		}
		return nil, err
	}

	result := make([]ogen.MovieGenre, len(genres))
	for i, g := range genres {
		result[i] = *movieGenreToOgen(&g)
	}

	return (*ogen.GetMovieGenresOKApplicationJSON)(&result), nil
}

// GetMovieCollection delegates to the movie handler.
func (h *Handler) GetMovieCollection(ctx context.Context, params ogen.GetMovieCollectionParams) (ogen.GetMovieCollectionRes, error) {
	collection, err := h.movieHandler.GetMovieCollection(ctx, params.ID.String())
	if err != nil {
		if errors.Is(err, movie.ErrMovieNotFound) || errors.Is(err, movie.ErrNotInCollection) {
			return &ogen.GetMovieCollectionNotFound{}, nil
		}
		return nil, err
	}

	return movieCollectionToOgen(collection), nil
}

// GetWatchProgress delegates to the movie handler.
func (h *Handler) GetWatchProgress(ctx context.Context, params ogen.GetWatchProgressParams) (ogen.GetWatchProgressRes, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	progress, err := h.movieHandler.GetWatchProgress(ctx, userID, params.ID.String())
	if err != nil {
		if errors.Is(err, movie.ErrProgressNotFound) {
			return &ogen.GetWatchProgressNotFound{}, nil
		}
		return nil, err
	}

	return movieWatchedToOgen(progress), nil
}

// UpdateWatchProgress delegates to the movie handler.
func (h *Handler) UpdateWatchProgress(ctx context.Context, req *ogen.UpdateWatchProgressReq, params ogen.UpdateWatchProgressParams) (ogen.UpdateWatchProgressRes, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	updateParams := movie.UpdateWatchProgressParams{
		ProgressSeconds: util.SafeIntToInt32(req.ProgressSeconds),
		DurationSeconds: util.SafeIntToInt32(req.DurationSeconds),
	}

	progress, err := h.movieHandler.UpdateWatchProgress(ctx, userID, params.ID.String(), updateParams)
	if err != nil {
		if errors.Is(err, movie.ErrMovieNotFound) {
			return &ogen.UpdateWatchProgressNotFound{}, nil
		}
		return nil, err
	}

	return movieWatchedToOgen(progress), nil
}

// DeleteWatchProgress delegates to the movie handler.
func (h *Handler) DeleteWatchProgress(ctx context.Context, params ogen.DeleteWatchProgressParams) (ogen.DeleteWatchProgressRes, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	err = h.movieHandler.DeleteWatchProgress(ctx, userID, params.ID.String())
	if err != nil {
		if errors.Is(err, movie.ErrProgressNotFound) {
			return &ogen.DeleteWatchProgressNotFound{}, nil
		}
		return nil, err
	}

	return &ogen.DeleteWatchProgressNoContent{}, nil
}

// MarkAsWatched delegates to the movie handler.
func (h *Handler) MarkAsWatched(ctx context.Context, params ogen.MarkAsWatchedParams) (ogen.MarkAsWatchedRes, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	err = h.movieHandler.MarkAsWatched(ctx, userID, params.ID.String())
	if err != nil {
		if errors.Is(err, movie.ErrMovieNotFound) {
			return &ogen.MarkAsWatchedNotFound{}, nil
		}
		return nil, err
	}

	return &ogen.MarkAsWatchedNoContent{}, nil
}

// RefreshMovieMetadata delegates to the movie handler.
func (h *Handler) RefreshMovieMetadata(ctx context.Context, params ogen.RefreshMovieMetadataParams) (ogen.RefreshMovieMetadataRes, error) {
	err := h.movieHandler.RefreshMetadata(ctx, params.ID.String())
	if err != nil {
		if errors.Is(err, movie.ErrMovieNotFound) {
			return &ogen.RefreshMovieMetadataNotFound{}, nil
		}
		return nil, err
	}

	return &ogen.RefreshMovieMetadataAccepted{}, nil
}

// GetCollection returns details about a movie collection.
func (h *Handler) GetCollection(ctx context.Context, params ogen.GetCollectionParams) (ogen.GetCollectionRes, error) {
	collection, err := h.movieHandler.GetCollection(ctx, params.ID.String())
	if err != nil {
		if errors.Is(err, movie.ErrCollectionNotFound) {
			return &ogen.GetCollectionNotFound{}, nil
		}
		return nil, err
	}

	return movieCollectionToOgen(collection), nil
}

// GetCollectionMovies returns all movies in a collection.
func (h *Handler) GetCollectionMovies(ctx context.Context, params ogen.GetCollectionMoviesParams) (ogen.GetCollectionMoviesRes, error) {
	movies, err := h.movieHandler.GetCollectionMovies(ctx, params.ID.String())
	if err != nil {
		if errors.Is(err, movie.ErrCollectionNotFound) {
			return &ogen.GetCollectionMoviesNotFound{}, nil
		}
		return nil, err
	}

	result := make([]ogen.Movie, len(movies))
	for i, m := range movies {
		result[i] = *movieToOgen(&m)
	}

	return (*ogen.GetCollectionMoviesOKApplicationJSON)(&result), nil
}

// GetSimilarMovies returns movies similar to the given movie.
func (h *Handler) GetSimilarMovies(ctx context.Context, params ogen.GetSimilarMoviesParams) (ogen.GetSimilarMoviesRes, error) {
	// Get the movie to find its TMDb ID
	m, err := h.movieHandler.GetMovie(ctx, params.ID.String())
	if err != nil {
		if errors.Is(err, movie.ErrMovieNotFound) {
			return &ogen.GetSimilarMoviesNotFound{}, nil
		}
		return nil, err
	}

	if m.TMDbID == nil {
		return &ogen.SimilarMoviesResponse{
			MovieID:       ogen.NewOptUUID(params.ID),
			SimilarMovies: []ogen.SimilarMovie{},
			TotalResults:  ogen.NewOptInt(0),
		}, nil
	}

	// Get similar movies from shared metadata service
	similar, totalResults, err := h.metadataService.GetSimilarMovies(ctx, *m.TMDbID, metadata.SearchOptions{})
	if err != nil {
		return nil, err
	}

	// Get all TMDb IDs in library for in_library flag
	libraryTMDbIDs := make(map[int32]bool)
	libraryMovies, _ := h.movieHandler.ListMovies(ctx, movie.ListMoviesParams{Limit: 10000})
	for _, lm := range libraryMovies {
		if lm.TMDbID != nil {
			libraryTMDbIDs[*lm.TMDbID] = true
		}
	}

	// Convert to API response
	similarMovies := make([]ogen.SimilarMovie, 0, len(similar))
	for _, sm := range similar {
		sim := ogen.SimilarMovie{
			Title: ogen.NewOptString(sm.Title),
		}

		// Parse TMDb ID from provider ID
		if sm.ProviderID != "" {
			var tmdbID int32
			_, _ = fmt.Sscanf(sm.ProviderID, "%d", &tmdbID)
			if tmdbID > 0 {
				sim.TmdbID = ogen.NewOptInt(int(tmdbID))
				sim.InLibrary = ogen.NewOptBool(libraryTMDbIDs[tmdbID])
			}
		}
		if sm.OriginalTitle != "" {
			sim.OriginalTitle = ogen.NewOptString(sm.OriginalTitle)
		}
		if sm.Overview != "" {
			sim.Overview = ogen.NewOptNilString(sm.Overview)
		}
		if sm.ReleaseDate != nil {
			sim.ReleaseDate = ogen.NewOptNilDate(*sm.ReleaseDate)
		}
		if sm.PosterPath != nil {
			sim.PosterPath = ogen.NewOptNilString(*sm.PosterPath)
		}
		if sm.BackdropPath != nil {
			sim.BackdropPath = ogen.NewOptNilString(*sm.BackdropPath)
		}
		if sm.VoteAverage > 0 {
			sim.VoteAverage = ogen.NewOptFloat32(float32(sm.VoteAverage))
		}
		if sm.VoteCount > 0 {
			sim.VoteCount = ogen.NewOptInt(sm.VoteCount)
		}
		if sm.Popularity > 0 {
			sim.Popularity = ogen.NewOptFloat32(float32(sm.Popularity))
		}

		similarMovies = append(similarMovies, sim)
	}

	return &ogen.SimilarMoviesResponse{
		MovieID:       ogen.NewOptUUID(params.ID),
		SimilarMovies: similarMovies,
		TotalResults:  ogen.NewOptInt(totalResults),
	}, nil
}
