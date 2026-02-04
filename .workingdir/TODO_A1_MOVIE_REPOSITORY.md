# Phase A1: Movie Repository Completion

**Priority**: P2
**Effort**: 6-8h
**Dependencies**: A0
**Affected File**: `internal/content/movie/repository_postgres.go`
**Status**: ✅ COMPLETED (2026-02-04)

> All 25+ repository methods implemented

---

## A1.1: Core Movie Queries - ✅ COMPLETED

- [x] `ListMovies(ctx, filters) ([]Movie, error)`
- [x] `CountMovies(ctx) (int64, error)`
- [x] `SearchMoviesByTitle(ctx, query, limit, offset) ([]Movie, error)`
- [x] `ListMoviesByYear(ctx, year, limit, offset) ([]Movie, error)`
- [x] `ListRecentlyAdded(ctx, limit, offset) ([]Movie, error)`
- [x] `ListTopRated(ctx, minVotes, limit, offset) ([]Movie, error)`

---

## A1.2: Movie Metadata Operations - ✅ COMPLETED

- [x] `CreateMovie(ctx, params) (*Movie, error)`
- [x] `UpdateMovie(ctx, params) (*Movie, error)`
- [x] `DeleteMovie(ctx, id) error`

---

## A1.3: Movie Credits - ✅ COMPLETED

- [x] `CreateMovieCredit(ctx, params) (*MovieCredit, error)`
- [x] `ListMovieCast(ctx, movieID) ([]MovieCredit, error)`
- [x] `ListMovieCrew(ctx, movieID) ([]MovieCredit, error)`
- [x] `DeleteMovieCredits(ctx, movieID) error`

---

## A1.4: Movie Genres - ✅ COMPLETED

- [x] `AddMovieGenre(ctx, movieID, tmdbGenreID, name) error`
- [x] `ListMovieGenres(ctx, movieID) ([]MovieGenre, error)`
- [x] `DeleteMovieGenres(ctx, movieID) error`
- [x] `ListMoviesByGenre(ctx, tmdbGenreID, limit, offset) ([]Movie, error)`

---

## A1.5: Movie Files - ✅ COMPLETED

- [x] `CreateMovieFile(ctx, params) (*MovieFile, error)`
- [x] `GetMovieFile(ctx, id) (*MovieFile, error)`
- [x] `GetMovieFileByPath(ctx, path) (*MovieFile, error)`
- [x] `ListMovieFilesByMovieID(ctx, movieID) ([]MovieFile, error)`
- [x] `UpdateMovieFile(ctx, params) (*MovieFile, error)`
- [x] `DeleteMovieFile(ctx, id) error`

---

## A1.6: Watch Progress - ✅ COMPLETED

- [x] `CreateOrUpdateWatchProgress(ctx, params) (*MovieWatched, error)`
- [x] `GetWatchProgress(ctx, userID, movieID) (*MovieWatched, error)`
- [x] `DeleteWatchProgress(ctx, userID, movieID) error`
- [x] `ListContinueWatching(ctx, userID, limit) ([]ContinueWatchingItem, error)`
- [x] `ListWatchedMovies(ctx, userID, limit, offset) ([]WatchedMovieItem, error)`
- [x] `GetUserMovieStats(ctx, userID) (*UserMovieStats, error)`

---

## Implementation Notes

### Helper Functions Added
- `stringToPgDate` - Convert string date to pgtype.Date
- `stringToPgNumeric` - Convert string to pgtype.Numeric
- `stringToPgTimestamptz` - Convert string to pgtype.Timestamptz
- `dbWatchedToWatched` - Convert DB MovieWatched to domain
- `dbContinueWatchingRowToMovie` - Convert ListContinueWatchingRow to Movie
- `dbWatchedMovieRowToMovie` - Convert ListWatchedMoviesRow to Movie
- `dbCreditToCredit` - Convert DB MovieCredit to domain
- `dbGenreToGenre` - Convert DB MovieGenre to domain
- `pgNumericToInt32Ptr` - Convert pgtype.Numeric to *int32
- `derefInt32` - Safely dereference *int32
- `derefBool` - Safely dereference *bool

### Already Existed (not touched)
- Collection methods (CreateMovieCollection, GetMovieCollection, etc.)
- GetMovie, GetMovieByTMDbID, GetMovieByIMDbID, GetMovieByRadarrID
- GetMovieFileByRadarrID
- All sqlc queries in `internal/content/movie/db/movies.sql.go`

### Build & Tests
- `go build ./internal/content/movie/...` ✅
- `go test ./internal/content/movie/...` ✅
