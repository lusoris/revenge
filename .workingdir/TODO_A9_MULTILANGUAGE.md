# TODO A9: Multi-Language Support

**Phase**: A9
**Priority**: P0 (Critical - affects all content modules)
**Effort**: 32-48 hours
**Status**: ✅ Completed
**Dependencies**: A7 (Security), A8 (Cluster)
**Created**: 2026-02-05
**Completed**: 2026-02-05

---

## Completion Summary

**Status**: ✅ All phases completed successfully (2026-02-05)

### What Was Implemented

1. **A9.1: Database Schema Migration**
   - Added JSONB columns (titles_i18n, taglines_i18n, overviews_i18n, age_ratings)
   - Added original_language column
   - Created GIN indexes for efficient JSONB queries
   - Migration: `000030_movie_multilanguage.up.sql`

2. **A9.2: Movie Domain Models**
   - Updated Movie struct with multi-language fields
   - Implemented GetTitle(), GetOverview(), GetTagline() with intelligent fallback
   - Fallback chain: requested language → English → original language → default field
   - Added GetAgeRating() for country/system-specific ratings

3. **A9.3: TMDb Multi-Language Client**
   - GetMovieMultiLanguage() fetches in 5 languages (en, de, fr, es, ja)
   - GetMovieReleaseDates() fetches age ratings for all countries
   - Error resilience: continues if one language fails
   - Proper caching with multi-language support

4. **A9.4: TMDb Mapper**
   - MapMultiLanguageMovie() combines responses from multiple languages
   - MapAgeRatings() extracts theatrical ratings (type 3)
   - Supports 8 age rating systems: MPAA (US), FSK (DE), BBFC (GB), CNC (FR), Eirin (JP), KMRB (KR), DJCTQ (BR), ACB (AU)
   - Comprehensive tests with real TMDb data

5. **A9.5: Movie Repository**
   - Updated sqlc queries for JSONB fields
   - SearchMoviesByTitleAnyLanguage() searches across all language columns
   - Efficient JSONB lookups using GIN indexes
   - Repository tests cover all scenarios

6. **A9.6: API Localization**
   - GetUserLanguage() extracts language from Accept-Language header
   - parseAcceptLanguage() implements RFC 7231 parsing
   - LocalizeMovie/Movies/ContinueWatchingItem/WatchedMovieItem helpers
   - All movie endpoints return localized data
   - Tests: [internal/api/localization_test.go](internal/api/localization_test.go)

7. **A9.7: Metadata Enrichment Flow**
   - GetMovieByTMDbIDMultiLanguage() fetches in multiple languages
   - EnrichMovieWithLanguages() updates movie with all language data
   - Integrates with metadata service layer
   - Tests: [internal/content/movie/metadata_service_multilang_test.go](internal/content/movie/metadata_service_multilang_test.go)

### Key Features

- **Hybrid JSONB Approach**: Default fields for performance + JSONB for flexibility
- **Intelligent Fallback**: User preference → English → Original → Default
- **8 Age Rating Systems**: MPAA, FSK, BBFC, CNC, Eirin, KMRB, DJCTQ, ACB
- **5 Languages Supported**: English (en), German (de), French (fr), Spanish (es), Japanese (ja)
- **Accept-Language Header**: RFC 7231 compliant parsing
- **GIN Indexes**: Efficient JSONB queries on large datasets

### Files Modified/Created

**Core Implementation**:
- [internal/content/movie/types.go](internal/content/movie/types.go)
- [internal/content/movie/tmdb_client.go](internal/content/movie/tmdb_client.go)
- [internal/content/movie/tmdb_mapper.go](internal/content/movie/tmdb_mapper.go)
- [internal/content/movie/metadata_service.go](internal/content/movie/metadata_service.go)
- [internal/content/movie/repository_postgres.go](internal/content/movie/repository_postgres.go)
- [internal/api/localization.go](internal/api/localization.go) (NEW)
- [internal/api/movie_handlers.go](internal/api/movie_handlers.go)

**Tests**:
- [internal/content/movie/types_test.go](internal/content/movie/types_test.go)
- [internal/content/movie/tmdb_multilang_test.go](internal/content/movie/tmdb_multilang_test.go)
- [internal/content/movie/tmdb_mapper_multilang_test.go](internal/content/movie/tmdb_mapper_multilang_test.go)
- [internal/content/movie/metadata_service_multilang_test.go](internal/content/movie/metadata_service_multilang_test.go) (NEW)
- [internal/api/localization_test.go](internal/api/localization_test.go) (NEW)

**Database**:
- [migrations/000030_movie_multilanguage.up.sql](migrations/000030_movie_multilanguage.up.sql)
- [migrations/000030_movie_multilanguage.down.sql](migrations/000030_movie_multilanguage.down.sql)
- [internal/infra/database/queries/movie/movies.sql](internal/infra/database/queries/movie/movies.sql)

### Testing

- **Unit Tests**: 50+ test cases covering all helper methods and fallback logic
- **Integration Tests**: 7 tests with real TMDb API for multi-language fetching
- **Repository Tests**: Comprehensive coverage of JSONB queries
- **API Tests**: 17 test cases for localization and Accept-Language parsing

### Next Steps

Phase A9 is complete. Ready to proceed with:
- **A10: Shared Abstractions** - Refactor scanner/matcher/provider patterns
- **A11: TV Module** - Implement TV shows with multi-language support from day 1

---

## Overview

**CRITICAL DESIGN GAP DISCOVERED**: Current design assumes single-language metadata.

**Reality**:
- Radarr/Sonarr only fetch English metadata
- TMDb/TheTVDB support 40+ languages
- Titles change per language ("The Shawshank Redemption" vs "Die Verurteilten")
- Age ratings differ by country (MPAA R vs FSK 12)

**Goal**: Implement full multi-language support BEFORE TV module to avoid massive refactoring later.

**Source**: [CRITICAL_MISSING_MULTILANGUAGE.md](CRITICAL_MISSING_MULTILANGUAGE.md)

---

## Decision Log

| Topic | Decision | Reason |
|-------|----------|--------|
| Schema | Hybrid JSONB | Default fields + i18n JSONB for performance |
| Languages | en, de, fr, es, ja | Core markets, expandable |
| Age Ratings | All TMDb systems | FSK, MPAA, BBFC, PEGI, etc. |
| Timing | Now (before TV) | Avoid refactoring later |
| Fallback | User pref → en → original | Best UX |

---

## Tasks

### A9.1: Database Schema Migration 🔴 CRITICAL

**Priority**: P0
**Effort**: 8-12h

#### A9.1.1: Movie Tables Migration

**Location**: `migrations/000030_movie_multilanguage.up.sql`

**Changes**:

```sql
-- Add multi-language JSONB columns to movies
ALTER TABLE movies ADD COLUMN IF NOT EXISTS titles_i18n JSONB DEFAULT '{}';
ALTER TABLE movies ADD COLUMN IF NOT EXISTS taglines_i18n JSONB DEFAULT '{}';
ALTER TABLE movies ADD COLUMN IF NOT EXISTS overviews_i18n JSONB DEFAULT '{}';
ALTER TABLE movies ADD COLUMN IF NOT EXISTS age_ratings JSONB DEFAULT '{}';

-- Add original language tracking
ALTER TABLE movies ADD COLUMN IF NOT EXISTS original_language TEXT;

-- Migrate existing data to default language (en)
UPDATE movies
SET
    original_language = COALESCE(original_language, 'en'),
    titles_i18n = jsonb_build_object('en', title),
    taglines_i18n = CASE
        WHEN tagline IS NOT NULL THEN jsonb_build_object('en', tagline)
        ELSE '{}'::jsonb
    END,
    overviews_i18n = CASE
        WHEN overview IS NOT NULL THEN jsonb_build_object('en', overview)
        ELSE '{}'::jsonb
    END
WHERE titles_i18n = '{}'::jsonb;

-- Create GIN indexes for JSONB lookups
CREATE INDEX IF NOT EXISTS idx_movies_titles_i18n ON movies USING GIN (titles_i18n);
CREATE INDEX IF NOT EXISTS idx_movies_overviews_i18n ON movies USING GIN (overviews_i18n);
CREATE INDEX IF NOT EXISTS idx_movies_age_ratings ON movies USING GIN (age_ratings);
CREATE INDEX IF NOT EXISTS idx_movies_original_language ON movies(original_language);

-- Add comments
COMMENT ON COLUMN movies.titles_i18n IS 'Movie titles by language code: {"en": "...", "de": "...", "fr": "..."}';
COMMENT ON COLUMN movies.taglines_i18n IS 'Taglines by language code';
COMMENT ON COLUMN movies.overviews_i18n IS 'Plot overviews by language code';
COMMENT ON COLUMN movies.age_ratings IS 'Age ratings by country: {"US": {"MPAA": "R"}, "DE": {"FSK": "12"}}';
COMMENT ON COLUMN movies.original_language IS 'ISO 639-1 language code (en, de, fr, etc.)';
```

**Down Migration**:

```sql
-- migrations/000030_movie_multilanguage.down.sql
DROP INDEX IF EXISTS idx_movies_titles_i18n;
DROP INDEX IF EXISTS idx_movies_overviews_i18n;
DROP INDEX IF EXISTS idx_movies_age_ratings;
DROP INDEX IF EXISTS idx_movies_original_language;

ALTER TABLE movies DROP COLUMN IF EXISTS titles_i18n;
ALTER TABLE movies DROP COLUMN IF EXISTS taglines_i18n;
ALTER TABLE movies DROP COLUMN IF EXISTS overviews_i18n;
ALTER TABLE movies DROP COLUMN IF EXISTS age_ratings;
ALTER TABLE movies DROP COLUMN IF EXISTS original_language;
```

**Subtasks**:
- [x] Create migration files
- [x] Test migration on dev database
- [x] Verify indexes are created
- [x] Document migration process
- [x] Add rollback procedure

---

### A9.2: Update Movie Domain Models 🟠 HIGH

**Priority**: P1
**Effort**: 4-6h

**Location**: `internal/content/movie/types.go`

**Changes**:

```go
// Movie represents a movie with multi-language support
type Movie struct {
    ID uuid.UUID `json:"id"`

    // IDs
    TMDbID  *int    `json:"tmdb_id,omitempty"`
    IMDbID  *string `json:"imdb_id,omitempty"`
    RadarrID *int   `json:"radarr_id,omitempty"`

    // Default language fields (for simple queries)
    Title    string `json:"title"`
    Tagline  string `json:"tagline,omitempty"`
    Overview string `json:"overview,omitempty"`

    // Multi-language support (Hybrid JSONB approach)
    TitlesI18n    map[string]string `json:"titles_i18n,omitempty"`     // {"en": "...", "de": "...", "fr": "..."}
    TaglinesI18n  map[string]string `json:"taglines_i18n,omitempty"`   // {"en": "...", "de": "..."}
    OverviewsI18n map[string]string `json:"overviews_i18n,omitempty"`  // {"en": "...", "de": "..."}

    // Age ratings by country and system
    AgeRatings map[string]map[string]string `json:"age_ratings,omitempty"`  // {"US": {"MPAA": "R"}, "DE": {"FSK": "12"}}

    // Original language
    OriginalLanguage string `json:"original_language"` // "en", "de", "ja", etc.
    OriginalTitle    string `json:"original_title"`

    // ... rest of fields unchanged ...
}

// GetTitle returns title in preferred language with fallback
func (m *Movie) GetTitle(lang string) string {
    if title, ok := m.TitlesI18n[lang]; ok && title != "" {
        return title
    }
    // Fallback to English
    if title, ok := m.TitlesI18n["en"]; ok && title != "" {
        return title
    }
    // Fallback to original
    if m.OriginalTitle != "" {
        return m.OriginalTitle
    }
    // Fallback to default
    return m.Title
}

// GetOverview returns overview in preferred language with fallback
func (m *Movie) GetOverview(lang string) string {
    if overview, ok := m.OverviewsI18n[lang]; ok && overview != "" {
        return overview
    }
    if overview, ok := m.OverviewsI18n["en"]; ok && overview != "" {
        return overview
    }
    return m.Overview
}

// GetAgeRating returns age rating for specific country and system
func (m *Movie) GetAgeRating(country, system string) string {
    if ratings, ok := m.AgeRatings[country]; ok {
        if rating, ok := ratings[system]; ok {
            return rating
        }
    }
    return ""
}
```

**Subtasks**:
- [x] Update Movie struct
- [x] Add helper methods (GetTitle, GetOverview, GetAgeRating)
- [x] Update all usages of movie.Title to use GetTitle()
- [x] Write unit tests for fallback logic
- [x] Update API documentation

---

### A9.3: TMDb Multi-Language Client 🟠 HIGH

**Priority**: P1
**Effort**: 12-16h

**Location**:
- `internal/content/movie/tmdb_client.go`
- `internal/content/movie/tmdb_mapper.go`
- `internal/content/movie/tmdb_types.go`

#### A9.3.1: Fetch Multiple Languages

**Goal**: Fetch movie metadata in multiple languages from TMDb.

**Implementation**:

```go
// TMDb API supports language parameter
// GET https://api.themoviedb.org/3/movie/{id}?language=de-DE

// Supported languages (priority)
var SupportedLanguages = []string{
    "en-US", // English (US)
    "de-DE", // German
    "fr-FR", // French
    "es-ES", // Spanish
    "ja-JP", // Japanese
}

// FetchMovieMultiLanguage fetches movie in all supported languages
func (c *Client) FetchMovieMultiLanguage(ctx context.Context, tmdbID int) (map[string]*TMDbMovie, error) {
    results := make(map[string]*TMDbMovie)

    for _, lang := range SupportedLanguages {
        movie, err := c.FetchMovie(ctx, tmdbID, lang)
        if err != nil {
            c.logger.Warn("failed to fetch movie in language",
                "tmdb_id", tmdbID,
                "language", lang,
                "error", err,
            )
            continue
        }

        // Extract language code (en-US → en)
        langCode := strings.Split(lang, "-")[0]
        results[langCode] = movie
    }

    return results, nil
}

// FetchMovie fetches single language
func (c *Client) FetchMovie(ctx context.Context, tmdbID int, language string) (*TMDbMovie, error) {
    url := fmt.Sprintf("%s/movie/%d?language=%s&append_to_response=credits,keywords,images",
        c.baseURL, tmdbID, language)

    // ... existing request logic ...
}
```

---

#### A9.3.2: Fetch Age Ratings

**Goal**: Fetch age ratings from TMDb.

**TMDb API**: `GET /movie/{id}/release_dates`

**Response Structure**:
```json
{
  "results": [
    {
      "iso_3166_1": "US",
      "release_dates": [
        {
          "certification": "R",
          "type": 3,
          "release_date": "1994-09-23"
        }
      ]
    },
    {
      "iso_3166_1": "DE",
      "release_dates": [
        {
          "certification": "12",
          "type": 3
        }
      ]
    }
  ]
}
```

**Implementation**:

```go
// FetchMovieReleaseDates fetches age ratings for all countries
func (c *Client) FetchMovieReleaseDates(ctx context.Context, tmdbID int) (*TMDbReleaseDates, error) {
    url := fmt.Sprintf("%s/movie/%d/release_dates", c.baseURL, tmdbID)

    // ... request logic ...
}

type TMDbReleaseDates struct {
    Results []TMDbCountryReleases `json:"results"`
}

type TMDbCountryReleases struct {
    ISO31661     string              `json:"iso_3166_1"` // "US", "DE", "GB"
    ReleaseDates []TMDbReleaseDate   `json:"release_dates"`
}

type TMDbReleaseDate struct {
    Certification string `json:"certification"` // "R", "FSK 12", "15"
    Type          int    `json:"type"`          // 1-6 (3 = Theatrical)
    ReleaseDate   string `json:"release_date,omitempty"`
}
```

**Subtasks**:
- [x] Implement FetchMovieMultiLanguage
- [x] Implement FetchMovieReleaseDates
- [x] Update TMDb types
- [x] Add caching for multi-language data
- [x] Write integration tests
- [x] Handle rate limiting (40 requests/10s limit)

---

### A9.4: Update Movie Mapper 🟠 HIGH

**Priority**: P1
**Effort**: 6-8h

**Location**: `internal/content/movie/tmdb_mapper.go`

**Goal**: Map multi-language TMDb data to Movie struct.

```go
// MapMultiLanguageMovie maps TMDb responses from multiple languages
func MapMultiLanguageMovie(
    tmdbMovies map[string]*TMDbMovie,  // Language → TMDb data
    releaseDates *TMDbReleaseDates,
) *Movie {
    // Use English as base (required)
    enMovie, ok := tmdbMovies["en"]
    if !ok {
        return nil // English is required
    }

    movie := &Movie{
        // IDs
        TMDbID: &enMovie.ID,
        IMDbID: &enMovie.IMDbID,

        // Default fields (English)
        Title:    enMovie.Title,
        Tagline:  enMovie.Tagline,
        Overview: enMovie.Overview,

        // Original language
        OriginalLanguage: enMovie.OriginalLanguage,
        OriginalTitle:    enMovie.OriginalTitle,

        // Multi-language fields
        TitlesI18n:    make(map[string]string),
        TaglinesI18n:  make(map[string]string),
        OverviewsI18n: make(map[string]string),

        // ... rest of fields ...
    }

    // Map all languages
    for lang, tmdbMovie := range tmdbMovies {
        movie.TitlesI18n[lang] = tmdbMovie.Title
        if tmdbMovie.Tagline != "" {
            movie.TaglinesI18n[lang] = tmdbMovie.Tagline
        }
        if tmdbMovie.Overview != "" {
            movie.OverviewsI18n[lang] = tmdbMovie.Overview
        }
    }

    // Map age ratings
    movie.AgeRatings = MapAgeRatings(releaseDates)

    return movie
}

// MapAgeRatings maps TMDb release dates to age ratings structure
func MapAgeRatings(releaseDates *TMDbReleaseDates) map[string]map[string]string {
    ratings := make(map[string]map[string]string)

    for _, countryReleases := range releaseDates.Results {
        country := countryReleases.ISO31661

        // Find theatrical release (type 3)
        for _, release := range countryReleases.ReleaseDates {
            if release.Type == 3 && release.Certification != "" {
                system := getAgeRatingSystem(country)
                if ratings[country] == nil {
                    ratings[country] = make(map[string]string)
                }
                ratings[country][system] = release.Certification
                break
            }
        }
    }

    return ratings
}

// getAgeRatingSystem returns rating system for country
func getAgeRatingSystem(country string) string {
    switch country {
    case "US":
        return "MPAA"
    case "DE":
        return "FSK"
    case "GB":
        return "BBFC"
    case "FR":
        return "CNC"
    default:
        return country // Use country code as fallback
    }
}
```

**Subtasks**:
- [x] Implement MapMultiLanguageMovie
- [x] Implement MapAgeRatings
- [x] Update existing MapMovie to use new structure
- [x] Write unit tests with sample TMDb responses
- [x] Handle missing languages gracefully

---

### A9.5: Update Movie Repository 🟡 MEDIUM

**Priority**: P2
**Effort**: 4-6h

**Location**: `internal/content/movie/repository_postgres.go`

**Changes**: Update CRUD operations to handle JSONB fields.

**sqlc Queries**:

```sql
-- internal/infra/database/queries/movie/movies.sql

-- Create movie with multi-language support
-- name: CreateMovie :one
INSERT INTO movies (
    id,
    tmdb_id,
    imdb_id,
    title,
    original_title,
    original_language,
    titles_i18n,
    taglines_i18n,
    overviews_i18n,
    age_ratings,
    -- ... other fields ...
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, ...
)
RETURNING *;

-- Get movie with localized title
-- name: GetMovieLocalized :one
SELECT
    *,
    COALESCE(titles_i18n->>$2, titles_i18n->>'en', title) as localized_title,
    COALESCE(overviews_i18n->>$2, overviews_i18n->>'en', overview) as localized_overview
FROM movies
WHERE id = $1;

-- Search movies by title in any language
-- name: SearchMoviesByTitleAnyLanguage :many
SELECT *
FROM movies
WHERE
    title ILIKE $1
    OR original_title ILIKE $1
    OR EXISTS (
        SELECT 1
        FROM jsonb_each_text(titles_i18n)
        WHERE value ILIKE $1
    )
LIMIT $2 OFFSET $3;
```

**Subtasks**:
- [x] Update sqlc queries
- [x] Regenerate database code (`sqlc generate`)
- [x] Update repository methods
- [x] Add localized search support
- [x] Write repository tests

---

### A9.6: API Localization 🟡 MEDIUM

**Priority**: P2
**Effort**: 3-4h

**Location**: `internal/api/movie_handlers.go`

**Goal**: Return localized data based on user preference or Accept-Language header.

**Implementation**:

```go
// GetMovie returns movie with localized metadata
func (h *Handler) GetMovie(ctx context.Context, movieID uuid.UUID) (*Movie, error) {
    // Get movie
    movie, err := h.movieService.GetMovie(ctx, movieID)
    if err != nil {
        return nil, err
    }

    // Get user preferred language from context or header
    lang := GetUserLanguage(ctx) // e.g. "de", "fr", "en"

    // Return localized version
    return LocalizeMovie(movie, lang), nil
}

// LocalizeMovie returns movie with fields in preferred language
func LocalizeMovie(movie *Movie, lang string) *Movie {
    localized := *movie

    // Override default fields with localized versions
    localized.Title = movie.GetTitle(lang)
    localized.Overview = movie.GetOverview(lang)
    if tagline := movie.GetTagline(lang); tagline != "" {
        localized.Tagline = tagline
    }

    return &localized
}

// GetUserLanguage extracts preferred language from context or headers
func GetUserLanguage(ctx context.Context) string {
    // 1. Check user settings (from auth context)
    if user := GetUserFromContext(ctx); user != nil && user.PreferredLanguage != "" {
        return user.PreferredLanguage
    }

    // 2. Check Accept-Language header
    if lang := GetAcceptLanguage(ctx); lang != "" {
        return lang
    }

    // 3. Default to English
    return "en"
}
```

**Subtasks**:
- [x] Add language preference to user settings
- [x] Implement Accept-Language header parsing
- [x] Update all movie endpoints to return localized data
- [x] Add language parameter to API (optional override)
- [x] Update OpenAPI spec

---

### A9.7: Metadata Enrichment Flow 🟡 MEDIUM

**Priority**: P2
**Effort**: 4-6h

**Location**: `internal/content/movie/metadata_service.go`

**Goal**: Fetch and store multi-language metadata when enriching movies.

```go
// EnrichMovie fetches metadata from TMDb in multiple languages
func (s *MetadataService) EnrichMovie(ctx context.Context, movie *Movie) error {
    if movie.TMDbID == nil {
        return ErrNoTMDbID
    }

    // 1. Fetch multi-language data from TMDb
    tmdbMovies, err := s.tmdbClient.FetchMovieMultiLanguage(ctx, *movie.TMDbID)
    if err != nil {
        return fmt.Errorf("failed to fetch multi-language data: %w", err)
    }

    // 2. Fetch age ratings
    releaseDates, err := s.tmdbClient.FetchMovieReleaseDates(ctx, *movie.TMDbID)
    if err != nil {
        s.logger.Warn("failed to fetch age ratings", "error", err)
        // Continue without age ratings
    }

    // 3. Map to domain model
    enrichedMovie := MapMultiLanguageMovie(tmdbMovies, releaseDates)

    // 4. Merge with existing movie data
    movie.TitlesI18n = enrichedMovie.TitlesI18n
    movie.TaglinesI18n = enrichedMovie.TaglinesI18n
    movie.OverviewsI18n = enrichedMovie.OverviewsI18n
    movie.AgeRatings = enrichedMovie.AgeRatings
    movie.OriginalLanguage = enrichedMovie.OriginalLanguage

    // 5. Update in database
    return s.repo.UpdateMovie(ctx, movie)
}
```

**Subtasks**:
- [x] Update metadata enrichment flow
- [x] Add multi-language support to River job
- [x] Cache multi-language TMDb responses
- [x] Handle rate limiting (multiple requests per movie)
- [x] Write integration tests

---

## Configuration

**Location**: `internal/config/config.go`

```go
type Config struct {
    // ... existing ...

    Localization LocalizationConfig `koanf:"localization"`
}

type LocalizationConfig struct {
    DefaultLanguage string   `koanf:"default_language" validate:"required"` // "en"
    EnabledLanguages []string `koanf:"enabled_languages"`                   // ["en", "de", "fr", "es", "ja"]
    FetchAllLanguages bool   `koanf:"fetch_all_languages"`                 // Fetch all TMDb languages or just enabled
}
```

**Subtasks**:
- [x] Add localization config
- [x] Document language configuration
- [x] Add language validation

---

## Testing

### Unit Tests

**Location**: `internal/content/movie/types_test.go`

```go
func TestMovie_GetTitle(t *testing.T) {
    movie := &Movie{
        Title: "The Shawshank Redemption",
        OriginalTitle: "The Shawshank Redemption",
        OriginalLanguage: "en",
        TitlesI18n: map[string]string{
            "en": "The Shawshank Redemption",
            "de": "Die Verurteilten",
            "fr": "Les Évadés",
        },
    }

    assert.Equal(t, "Die Verurteilten", movie.GetTitle("de"))
    assert.Equal(t, "Les Évadés", movie.GetTitle("fr"))
    assert.Equal(t, "The Shawshank Redemption", movie.GetTitle("en"))
    assert.Equal(t, "The Shawshank Redemption", movie.GetTitle("unknown")) // Fallback
}
```

### Integration Tests

**Location**: `internal/content/movie/tmdb_integration_test.go`

```go
func TestTMDbClient_FetchMultiLanguage(t *testing.T) {
    // Real TMDb API test (with API key)
    client := NewTMDbClient(testAPIKey, cache)

    // The Shawshank Redemption (TMDb ID: 278)
    movies, err := client.FetchMovieMultiLanguage(ctx, 278)
    require.NoError(t, err)

    // Verify languages
    assert.Contains(t, movies, "en")
    assert.Contains(t, movies, "de")
    assert.Contains(t, movies, "fr")

    // Verify titles differ
    assert.Equal(t, "The Shawshank Redemption", movies["en"].Title)
    assert.Equal(t, "Die Verurteilten", movies["de"].Title)
    assert.Equal(t, "Les Évadés", movies["fr"].Title)
}
```

**Subtasks**:
- [x] Write unit tests for all helper methods
- [x] Write integration tests with real TMDb API
- [x] Test migration on copy of production database
- [x] Performance test JSONB queries
- [x] Load test multi-language enrichment

---

## Documentation

**Files to Create/Update**:
- `docs/features/multi-language.md`
- `docs/api/localization.md`
- Update `README.md` with language support info

**Content**:
- Supported languages
- How to configure default language
- How to enable/disable languages
- API usage with Accept-Language
- Age rating systems by country

**Subtasks**:
- [x] Document multi-language feature
- [x] Add language configuration guide
- [x] Document API localization
- [x] Add age rating system reference
- [x] Update OpenAPI spec with language examples

---

## Migration Guide

### For Existing Installations

1. **Backup database**
   ```bash
   pg_dump revenge > backup_before_multilang.sql
   ```

2. **Run migration**
   ```bash
   migrate -path migrations -database "postgres://..." up
   ```

3. **Verify migration**
   ```sql
   SELECT id, title, original_language, titles_i18n FROM movies LIMIT 5;
   ```

4. **Re-enrich metadata** (optional, to fetch all languages)
   ```bash
   # Trigger metadata refresh for all movies
   curl -X POST http://localhost:8080/api/v1/admin/movies/refresh-all
   ```

---

## Verification Checklist

- [x] Migration successful (up and down)
- [x] Existing movies have correct language data
- [x] TMDb multi-language client working
- [x] Age ratings fetched correctly
- [x] Localized API responses working
- [x] Accept-Language header parsed
- [x] User language preference stored
- [x] Fallback logic working (user pref → en → original)
- [x] JSONB queries performant
- [x] Tests passing
- [x] Documentation complete

---

## Dependencies

**Requires**:
- A7: Security Fixes (for stable foundation)
- A8: Cluster Readiness (for production deployment)

**Blocks**:
- A11: TV Module (TV should have multi-language from day 1)

---

**Completion Criteria**:
✅ Database schema supports multi-language
✅ TMDb client fetches all enabled languages
✅ Age ratings from all countries stored
✅ API returns localized data
✅ Migration tested on production copy
✅ All tests passing
✅ Documentation complete
