package movie

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	db "github.com/lusoris/revenge/internal/content/movie/db"
)

// pgRepository implements Repository using PostgreSQL.
type pgRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

// NewRepository creates a new PostgreSQL-backed movie repository.
func NewRepository(pool *pgxpool.Pool) Repository {
	return &pgRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

// GetByID retrieves a movie by its ID.
func (r *pgRepository) GetByID(ctx context.Context, id uuid.UUID) (*Movie, error) {
	m, err := r.queries.GetMovieByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMovieNotFound
		}
		return nil, err
	}
	return FromDBMovie(&m), nil
}

// GetByPath retrieves a movie by its file path.
func (r *pgRepository) GetByPath(ctx context.Context, path string) (*Movie, error) {
	m, err := r.queries.GetMovieByPath(ctx, path)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMovieNotFound
		}
		return nil, err
	}
	return FromDBMovie(&m), nil
}

// GetByTmdbID retrieves a movie by its TMDb ID.
func (r *pgRepository) GetByTmdbID(ctx context.Context, tmdbID int) (*Movie, error) {
	id := int32(tmdbID)
	m, err := r.queries.GetMovieByTmdbID(ctx, &id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMovieNotFound
		}
		return nil, err
	}
	return FromDBMovie(&m), nil
}

// GetByImdbID retrieves a movie by its IMDb ID.
func (r *pgRepository) GetByImdbID(ctx context.Context, imdbID string) (*Movie, error) {
	m, err := r.queries.GetMovieByImdbID(ctx, &imdbID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMovieNotFound
		}
		return nil, err
	}
	return FromDBMovie(&m), nil
}

// List retrieves movies with pagination.
func (r *pgRepository) List(ctx context.Context, params ListParams) ([]*Movie, error) {
	rows, err := r.queries.ListMovies(ctx, db.ListMoviesParams{
		Limit:     int32(params.Limit),
		Offset:    int32(params.Offset),
		SortBy:    params.SortBy,
		SortOrder: params.SortOrder,
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

// ListByLibrary retrieves movies from a specific library.
func (r *pgRepository) ListByLibrary(ctx context.Context, libraryID uuid.UUID, params ListParams) ([]*Movie, error) {
	rows, err := r.queries.ListMoviesByLibrary(ctx, db.ListMoviesByLibraryParams{
		LibraryID: libraryID,
		Limit:     int32(params.Limit),
		Offset:    int32(params.Offset),
		SortBy:    params.SortBy,
		SortOrder: params.SortOrder,
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

// ListByCollection retrieves movies in a collection.
func (r *pgRepository) ListByCollection(ctx context.Context, collectionID uuid.UUID) ([]*Movie, error) {
	rows, err := r.queries.ListMoviesByCollection(ctx, pgtype.UUID{Bytes: collectionID, Valid: true})
	if err != nil {
		return nil, err
	}

	movies := make([]*Movie, len(rows))
	for i, row := range rows {
		movies[i] = FromDBMovie(&row)
	}
	return movies, nil
}

// ListRecentlyAdded retrieves recently added movies.
func (r *pgRepository) ListRecentlyAdded(ctx context.Context, libraryIDs []uuid.UUID, limit int) ([]*Movie, error) {
	rows, err := r.queries.ListRecentlyAddedMovies(ctx, db.ListRecentlyAddedMoviesParams{
		LibraryIds: libraryIDs,
		Limit:      int32(limit),
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

// ListRecentlyPlayed retrieves recently played movies.
func (r *pgRepository) ListRecentlyPlayed(ctx context.Context, libraryIDs []uuid.UUID, limit int) ([]*Movie, error) {
	rows, err := r.queries.ListRecentlyPlayedMovies(ctx, db.ListRecentlyPlayedMoviesParams{
		LibraryIds: libraryIDs,
		Limit:      int32(limit),
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

// Search searches movies by title and overview.
func (r *pgRepository) Search(ctx context.Context, query string, params ListParams) ([]*Movie, error) {
	rows, err := r.queries.SearchMovies(ctx, db.SearchMoviesParams{
		PlaintoTsquery: query,
		Limit:          int32(params.Limit),
		Offset:         int32(params.Offset),
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

// Count returns the total number of movies.
func (r *pgRepository) Count(ctx context.Context) (int64, error) {
	return r.queries.CountMovies(ctx)
}

// CountByLibrary returns the number of movies in a library.
func (r *pgRepository) CountByLibrary(ctx context.Context, libraryID uuid.UUID) (int64, error) {
	return r.queries.CountMoviesByLibrary(ctx, libraryID)
}

// Create creates a new movie.
func (r *pgRepository) Create(ctx context.Context, movie *Movie) error {
	m, err := r.queries.CreateMovie(ctx, movie.ToDBCreateParams())
	if err != nil {
		return err
	}
	movie.ID = m.ID
	movie.CreatedAt = m.CreatedAt
	movie.UpdatedAt = m.UpdatedAt
	movie.DateAdded = m.DateAdded
	return nil
}

// Update updates an existing movie.
func (r *pgRepository) Update(ctx context.Context, movie *Movie) error {
	params := db.UpdateMovieParams{
		ID: movie.ID,
	}

	if movie.Title != "" {
		params.Title = &movie.Title
	}
	if movie.SortTitle != "" {
		params.SortTitle = &movie.SortTitle
	}
	if movie.OriginalTitle != "" {
		params.OriginalTitle = &movie.OriginalTitle
	}
	if movie.Tagline != "" {
		params.Tagline = &movie.Tagline
	}
	if movie.Overview != "" {
		params.Overview = &movie.Overview
	}
	if movie.ReleaseDate != nil {
		params.ReleaseDate = pgtype.Date{Time: *movie.ReleaseDate, Valid: true}
	}
	if movie.Year > 0 {
		y := int32(movie.Year)
		params.Year = &y
	}
	if movie.ContentRating != "" {
		params.ContentRating = &movie.ContentRating
	}
	if movie.RatingLevel >= 0 {
		rl := int32(movie.RatingLevel)
		params.RatingLevel = &rl
	}
	if movie.Budget > 0 {
		params.Budget = &movie.Budget
	}
	if movie.Revenue > 0 {
		params.Revenue = &movie.Revenue
	}
	if movie.CommunityRating > 0 {
		params.CommunityRating = numericFromFloat(movie.CommunityRating)
	}
	if movie.VoteCount > 0 {
		vc := int32(movie.VoteCount)
		params.VoteCount = &vc
	}
	if movie.PosterPath != "" {
		params.PosterPath = &movie.PosterPath
	}
	if movie.PosterBlurhash != "" {
		params.PosterBlurhash = &movie.PosterBlurhash
	}
	if movie.BackdropPath != "" {
		params.BackdropPath = &movie.BackdropPath
	}
	if movie.BackdropBlurhash != "" {
		params.BackdropBlurhash = &movie.BackdropBlurhash
	}
	if movie.LogoPath != "" {
		params.LogoPath = &movie.LogoPath
	}
	if movie.TmdbID > 0 {
		id := int32(movie.TmdbID)
		params.TmdbID = &id
	}
	if movie.ImdbID != "" {
		params.ImdbID = &movie.ImdbID
	}
	if movie.TvdbID > 0 {
		id := int32(movie.TvdbID)
		params.TvdbID = &id
	}
	if movie.CollectionID != nil {
		params.CollectionID = pgtype.UUID{Bytes: *movie.CollectionID, Valid: true}
	}
	if movie.CollectionOrder > 0 {
		co := int32(movie.CollectionOrder)
		params.CollectionOrder = &co
	}
	params.IsLocked = &movie.IsLocked

	m, err := r.queries.UpdateMovie(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrMovieNotFound
		}
		return err
	}
	movie.UpdatedAt = m.UpdatedAt
	return nil
}

// UpdatePlaybackStats updates the playback statistics for a movie.
func (r *pgRepository) UpdatePlaybackStats(ctx context.Context, id uuid.UUID) error {
	return r.queries.UpdateMoviePlaybackStats(ctx, id)
}

// Delete deletes a movie.
func (r *pgRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteMovie(ctx, id)
}

// DeleteByLibrary deletes all movies in a library.
func (r *pgRepository) DeleteByLibrary(ctx context.Context, libraryID uuid.UUID) error {
	return r.queries.DeleteMoviesByLibrary(ctx, libraryID)
}

// ExistsByPath checks if a movie exists by path.
func (r *pgRepository) ExistsByPath(ctx context.Context, path string) (bool, error) {
	return r.queries.MovieExistsByPath(ctx, path)
}

// ExistsByTmdbID checks if a movie exists by TMDb ID.
func (r *pgRepository) ExistsByTmdbID(ctx context.Context, tmdbID int) (bool, error) {
	id := int32(tmdbID)
	return r.queries.MovieExistsByTmdbID(ctx, &id)
}

// ListPaths returns all movie paths in a library.
func (r *pgRepository) ListPaths(ctx context.Context, libraryID uuid.UUID) (map[uuid.UUID]string, error) {
	rows, err := r.queries.ListMoviePaths(ctx, libraryID)
	if err != nil {
		return nil, err
	}

	paths := make(map[uuid.UUID]string, len(rows))
	for _, row := range rows {
		paths[row.ID] = row.Path
	}
	return paths, nil
}

// Collections

// GetCollectionByID retrieves a collection by ID.
func (r *pgRepository) GetCollectionByID(ctx context.Context, id uuid.UUID) (*Collection, error) {
	c, err := r.queries.GetCollectionByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCollectionNotFound
		}
		return nil, err
	}
	return fromDBCollection(&c), nil
}

// GetCollectionByTmdbID retrieves a collection by TMDb ID.
func (r *pgRepository) GetCollectionByTmdbID(ctx context.Context, tmdbID int) (*Collection, error) {
	id := int32(tmdbID)
	c, err := r.queries.GetCollectionByTmdbID(ctx, &id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCollectionNotFound
		}
		return nil, err
	}
	return fromDBCollection(&c), nil
}

// ListCollections lists all collections.
func (r *pgRepository) ListCollections(ctx context.Context, params ListParams) ([]*Collection, error) {
	rows, err := r.queries.ListCollections(ctx, db.ListCollectionsParams{
		Limit:  int32(params.Limit),
		Offset: int32(params.Offset),
	})
	if err != nil {
		return nil, err
	}

	collections := make([]*Collection, len(rows))
	for i, row := range rows {
		collections[i] = fromDBCollection(&row)
	}
	return collections, nil
}

// CountCollections returns the total number of collections.
func (r *pgRepository) CountCollections(ctx context.Context) (int64, error) {
	return r.queries.CountCollections(ctx)
}

// CreateCollection creates a new collection.
func (r *pgRepository) CreateCollection(ctx context.Context, collection *Collection) error {
	var tmdbID *int32
	if collection.TmdbID > 0 {
		id := int32(collection.TmdbID)
		tmdbID = &id
	}

	c, err := r.queries.CreateCollection(ctx, db.CreateCollectionParams{
		Name:             collection.Name,
		SortName:         &collection.SortName,
		Overview:         &collection.Overview,
		PosterPath:       &collection.PosterPath,
		PosterBlurhash:   &collection.PosterBlurhash,
		BackdropPath:     &collection.BackdropPath,
		BackdropBlurhash: &collection.BackdropBlurhash,
		TmdbID:           tmdbID,
	})
	if err != nil {
		return err
	}
	collection.ID = c.ID
	collection.CreatedAt = c.CreatedAt
	collection.UpdatedAt = c.UpdatedAt
	return nil
}

// UpdateCollection updates an existing collection.
func (r *pgRepository) UpdateCollection(ctx context.Context, collection *Collection) error {
	_, err := r.queries.UpdateCollection(ctx, db.UpdateCollectionParams{
		ID:               collection.ID,
		Name:             &collection.Name,
		SortName:         &collection.SortName,
		Overview:         &collection.Overview,
		PosterPath:       &collection.PosterPath,
		PosterBlurhash:   &collection.PosterBlurhash,
		BackdropPath:     &collection.BackdropPath,
		BackdropBlurhash: &collection.BackdropBlurhash,
	})
	return err
}

// DeleteCollection deletes a collection.
func (r *pgRepository) DeleteCollection(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteCollection(ctx, id)
}

// Studios

// GetStudioByID retrieves a studio by ID.
func (r *pgRepository) GetStudioByID(ctx context.Context, id uuid.UUID) (*Studio, error) {
	s, err := r.queries.GetStudioByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrStudioNotFound
		}
		return nil, err
	}
	return fromDBStudio(&s), nil
}

// GetStudioByTmdbID retrieves a studio by TMDb ID.
func (r *pgRepository) GetStudioByTmdbID(ctx context.Context, tmdbID int) (*Studio, error) {
	id := int32(tmdbID)
	s, err := r.queries.GetStudioByTmdbID(ctx, &id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrStudioNotFound
		}
		return nil, err
	}
	return fromDBStudio(&s), nil
}

// ListStudios lists all studios.
func (r *pgRepository) ListStudios(ctx context.Context, params ListParams) ([]*Studio, error) {
	rows, err := r.queries.ListStudios(ctx, db.ListStudiosParams{
		Limit:  int32(params.Limit),
		Offset: int32(params.Offset),
	})
	if err != nil {
		return nil, err
	}

	studios := make([]*Studio, len(rows))
	for i, row := range rows {
		studios[i] = fromDBStudio(&row)
	}
	return studios, nil
}

// CreateStudio creates a new studio.
func (r *pgRepository) CreateStudio(ctx context.Context, studio *Studio) error {
	var tmdbID *int32
	if studio.TmdbID > 0 {
		id := int32(studio.TmdbID)
		tmdbID = &id
	}

	s, err := r.queries.CreateStudio(ctx, db.CreateStudioParams{
		Name:     studio.Name,
		LogoPath: &studio.LogoPath,
		TmdbID:   tmdbID,
	})
	if err != nil {
		return err
	}
	studio.ID = s.ID
	studio.CreatedAt = s.CreatedAt
	return nil
}

// LinkMovieStudio links a studio to a movie.
func (r *pgRepository) LinkMovieStudio(ctx context.Context, movieID, studioID uuid.UUID, order int) error {
	return r.queries.LinkMovieStudio(ctx, db.LinkMovieStudioParams{
		MovieID:      movieID,
		StudioID:     studioID,
		DisplayOrder: int32(order),
	})
}

// UnlinkMovieStudios removes all studio links from a movie.
func (r *pgRepository) UnlinkMovieStudios(ctx context.Context, movieID uuid.UUID) error {
	return r.queries.UnlinkMovieStudios(ctx, movieID)
}

// GetMovieStudios retrieves studios for a movie.
func (r *pgRepository) GetMovieStudios(ctx context.Context, movieID uuid.UUID) ([]Studio, error) {
	rows, err := r.queries.ListMovieStudios(ctx, movieID)
	if err != nil {
		return nil, err
	}

	studios := make([]Studio, len(rows))
	for i, row := range rows {
		studios[i] = *fromDBStudio(&row)
	}
	return studios, nil
}

// Helper functions

func fromDBCollection(c *db.MovieCollection) *Collection {
	if c == nil {
		return nil
	}
	collection := &Collection{
		ID:        c.ID,
		Name:      c.Name,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
	if c.SortName != nil {
		collection.SortName = *c.SortName
	}
	if c.Overview != nil {
		collection.Overview = *c.Overview
	}
	if c.PosterPath != nil {
		collection.PosterPath = *c.PosterPath
	}
	if c.PosterBlurhash != nil {
		collection.PosterBlurhash = *c.PosterBlurhash
	}
	if c.BackdropPath != nil {
		collection.BackdropPath = *c.BackdropPath
	}
	if c.BackdropBlurhash != nil {
		collection.BackdropBlurhash = *c.BackdropBlurhash
	}
	if c.TmdbID != nil {
		collection.TmdbID = int(*c.TmdbID)
	}
	return collection
}

func fromDBStudio(s *db.MovieStudio) *Studio {
	if s == nil {
		return nil
	}
	studio := &Studio{
		ID:        s.ID,
		Name:      s.Name,
		CreatedAt: s.CreatedAt,
	}
	if s.LogoPath != nil {
		studio.LogoPath = *s.LogoPath
	}
	if s.TmdbID != nil {
		studio.TmdbID = int(*s.TmdbID)
	}
	return studio
}
