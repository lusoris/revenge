# Phase A5: Library Matcher

**Priority**: P2
**Effort**: 4-6h
**Dependencies**: A0, A1
**Status**: ✅ Complete (2026-02-04)

---

## A5.1: Library File Matching ✅

**Affected File**: `internal/content/movie/library_matcher.go`

**Completed Tasks**:
- [x] Implement `findExistingMovie` to search DB before TMDb lookup
- [x] Use Levenshtein distance for title similarity scoring
- [x] Score existing movies considering title and year
- [x] Return high-confidence matches (≥0.8) from existing DB

**Implementation Details**:
- `findExistingMovie` searches via `SearchMoviesByTitle`
- `scoreExistingMovie` calculates match score with Levenshtein distance
- Only returns matches with score ≥ 0.8 to avoid false positives

---

## A5.2: Credits & Genres Saving ✅

**Affected File**: `internal/content/movie/library_matcher.go`

**Completed Tasks**:
- [x] Save credits when creating movie from TMDb
- [x] Save genres when creating movie from TMDb
- [x] Use `CreateMovieCredit` and `AddMovieGenre` repository methods

---

## A5.3: Improved Confidence Scoring ✅

**Completed Tasks**:
- [x] Use Levenshtein distance for title similarity
- [x] Check original title as bonus
- [x] Refined year matching weights
- [x] Popularity-based confidence boost

**Scoring Breakdown**:
- Title match (Levenshtein): up to 0.5
- Original title bonus: +0.1
- Year exact match: +0.3 (±1 year: +0.15)
- Popularity boost: +0.05 to +0.1

---

## Matching Algorithm

1. **Check existing DB** - `findExistingMovie` searches local database first
2. **TMDb lookup** - If no DB match, search TMDb API
3. **Score candidates** - Calculate confidence using Levenshtein distance
4. **Create movie** - Save movie with credits and genres
5. **Return result** - Match type, confidence, created flag
