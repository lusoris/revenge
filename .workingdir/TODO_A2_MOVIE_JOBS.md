# Phase A2: Movie Jobs Completion

**Priority**: P2
**Effort**: 4-6h
**Dependencies**: A0, A1
**Status**: ✅ Complete (2026-02-04)

---

## A2.1: File Match Job ✅

**Affected Files**:
- `internal/content/movie/library_service.go` - Added `MatchFile` method
- `internal/content/movie/library_matcher.go` - Made `MatchFile` public
- `internal/content/movie/moviejobs/file_match.go` - Updated worker to use `MatchFile`

**Completed Tasks**:
- [x] Implement `library.Service.MatchFile` method
- [x] Update file match worker to use it
- [x] Match logic: filename parsing → TMDb lookup → confidence scoring
- [x] File record creation after successful match

**Implementation Details**:
- `LibraryService.MatchFile` checks if file exists, validates it's a video file
- Checks for existing match (unless force rematch)
- Parses filename to extract title and year
- Uses `Matcher.MatchFile` for TMDb lookup and confidence scoring
- Creates movie file record on successful match
- Worker logs match results with movie ID, title, match type, confidence

---

## A2.2: Metadata Refresh Job ✅

**Affected Files**:
- `internal/content/movie/library_service.go` - Updated `RefreshMovie` method
- `internal/content/movie/moviejobs/metadata_refresh.go` - Full implementation

**Completed Tasks**:
- [x] Implement full `RefreshMovie` in library service
- [x] Implement metadata refresh worker with credits/genres
- [x] Clear cache on force refresh
- [x] Update movie record in database
- [x] Refresh credits (delete old + fetch from TMDb + create new)
- [x] Refresh genres (delete old + fetch from TMDb + add new)

**Implementation Details**:
- Worker clears metadata cache when `Force=true`
- Enriches movie via `MetadataService.EnrichMovie`
- Updates movie record via `Repository.UpdateMovie`
- Deletes and recreates credits via `DeleteMovieCredits`/`CreateMovieCredit`
- Deletes and recreates genres via `DeleteMovieGenres`/`AddMovieGenre`
- Continues processing even if credits/genres fail (non-fatal)
