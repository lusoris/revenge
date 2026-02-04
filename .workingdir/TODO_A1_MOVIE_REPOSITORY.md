# Phase A1: Movie Repository Completion

**Priority**: P2
**Effort**: 6-8h
**Dependencies**: A0
**Affected File**: `internal/content/movie/repository_postgres.go:136-390`

> 25+ repository methods returning "not implemented"

---

## A1.1: Core Movie Queries

- [ ] `ListMoviesByIDs(ctx, ids []uuid.UUID) ([]*Movie, error)`
- [ ] `CountMovies(ctx) (int64, error)`
- [ ] `SearchMovies(ctx, query string, limit int) ([]*Movie, error)`
- [ ] `ListMoviesByGenre(ctx, genreID, limit, offset) ([]*Movie, error)`
- [ ] `ListRecentMovies(ctx, limit) ([]*Movie, error)`
- [ ] `ListPopularMovies(ctx, limit) ([]*Movie, error)`
- [ ] `GetMovieWithDetails(ctx, id) (*MovieWithDetails, error)`
- [ ] `GetMovieByExternalID(ctx, provider, externalID) (*Movie, error)`

---

## A1.2: Movie Metadata Operations

- [ ] `UpdateMovieMetadata(ctx, id, metadata) error`
- [ ] `DeleteMovie(ctx, id) error`

---

## A1.3: Movie Credits

- [ ] `ListMovieCredits(ctx, movieID) ([]*Credit, error)`
- [ ] `GetMovieCredit(ctx, creditID) (*Credit, error)`
- [ ] `CreateMovieCredit(ctx, credit) (*Credit, error)`
- [ ] `DeleteMovieCredit(ctx, creditID) error`
- [ ] `DeleteMovieCredits(ctx, movieID) error`

---

## A1.4: Movie Collections

- [ ] `GetMovieCollection(ctx, movieID) (*Collection, error)`
- [ ] `ListCollectionMovies(ctx, collectionID) ([]*Movie, error)`

---

## A1.5: Movie Files

- [ ] `ListMovieFiles(ctx, movieID) ([]*File, error)`
- [ ] `GetMovieFile(ctx, fileID) (*File, error)`
- [ ] `CreateMovieFile(ctx, file) (*File, error)`
- [ ] `DeleteMovieFile(ctx, fileID) error`

---

## Notes

**Important**: Many may already have sqlc queries - verify before implementing.

Check `internal/content/movie/db/` for existing generated queries before writing new ones.
