# Content Rating System

> Universal age restriction and content rating system for revenge.

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Full design with DB schema, Go filtering logic, per-module systems |
| Sources | 🟡 | International rating systems documented |
| Instructions | ✅ | Implementation checklist complete |
| Code | 🔴 | |
| Linting | 🔴 | |
| Unit Testing | 🔴 | |
| Integration Testing | 🔴 | |

**Location**: `internal/service/rating/`

---

## Overview

revenge implements a comprehensive content rating system that:
- Supports **international rating systems** (MPAA, FSK, BBFC, PEGI, etc.)
- Filters content based on **user age/permissions**
- Handles **person visibility** based on their content appearances
- Manages **image ratings** separately (SFW/NSFW)

## Design Principles

### Private Server Model

Revenge is a **private, invite-only** server. There is no public registration, so:
- ❌ No legal FSK18/age verification required
- ✅ Admin creates accounts, manages child profiles
- ✅ Child accounts get age-based filtering
- ✅ Parental controls via PIN protection

### Content-Based Filtering (NOT Person-Based)

```
CONTENT has Rating (Movie, Episode, Song, etc.)
PERSON has NO inherent Rating

Person Visibility = Does person have ANY visible content for this user?

Examples:
• Adult-Only Actor (only in adult content) → hidden for child users
• Action Star in FSK18 Movie → VISIBLE (FSK18 ≠ Adult/XXX)
• Actor with Mixed Content → visible, only adult roles hidden
```

### FSK18 vs Adult Distinction

| Category | Examples | Visibility |
|----------|----------|------------|
| `age_18` | Violence, Horror, War | Normal actors, visible to 18+ |
| `adult` | Explicit sexual content | Adult industry, requires opt-in |

## International Rating Systems

### Supported Systems

| System | Region | Ratings |
|--------|--------|---------|
| MPAA | USA | G, PG, PG-13, R, NC-17 |
| FSK | Germany | 0, 6, 12, 16, 18 |
| BBFC | UK | U, PG, 12A, 15, 18, R18 |
| PEGI | Europe (Games) | 3, 7, 12, 16, 18 |
| ACB | Australia | G, PG, M, MA15+, R18+, X18+ |
| CERO | Japan | A, B, C, D, Z |
| Kijkwijzer | Netherlands | AL, 6, 9, 12, 16, 18 |
| CNC | France | U, 10, 12, 16, 18, X |
| EIRIN | Japan (Film) | G, PG12, R15+, R18+ |
| CBFC | India | U, UA, A, S |
| ... | 30+ more | ... |

### Database Schema

```sql
-- Rating systems table (seeded with known systems)
CREATE TABLE rating_systems (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(20) NOT NULL UNIQUE,      -- 'mpaa', 'fsk', 'bbfc'
    name VARCHAR(100) NOT NULL,             -- 'Motion Picture Association'
    country_codes TEXT[] NOT NULL,          -- ['US', 'CA']
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Individual ratings within each system
CREATE TABLE ratings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    system_id UUID NOT NULL REFERENCES rating_systems(id),
    code VARCHAR(20) NOT NULL,              -- 'PG-13', 'FSK 16'
    name VARCHAR(100) NOT NULL,             -- 'Parental Guidance 13'
    description TEXT,
    min_age INT,                            -- Minimum age (0, 6, 12, 16, 18)
    normalized_level INT NOT NULL,          -- 0-100 scale for cross-system comparison
    sort_order INT NOT NULL,
    is_adult BOOLEAN DEFAULT false,         -- Explicit adult content flag
    icon_url VARCHAR(512),
    UNIQUE(system_id, code)
);

-- Cross-reference equivalents (for display)
CREATE TABLE rating_equivalents (
    rating_id UUID NOT NULL REFERENCES ratings(id),
    equivalent_rating_id UUID NOT NULL REFERENCES ratings(id),
    PRIMARY KEY (rating_id, equivalent_rating_id)
);

-- Content can have multiple ratings (from different systems)
CREATE TABLE content_ratings (
    content_id UUID NOT NULL,               -- FK to media_items, images, etc.
    content_type VARCHAR(50) NOT NULL,      -- 'media_item', 'image'
    rating_id UUID NOT NULL REFERENCES ratings(id),
    source VARCHAR(100),                    -- 'tmdb', 'manual', 'imdb'
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (content_id, rating_id)
);
```

### Normalized Level Scale

For cross-system filtering, all ratings map to a 0-100 normalized scale:

| Level | Description | Examples |
|-------|-------------|----------|
| 0 | All Ages | G, FSK 0, U |
| 25 | Young Children (6+) | PG, FSK 6 |
| 50 | Older Children (12+) | PG-13, FSK 12, 12A |
| 75 | Teens (16+) | R, FSK 16, 15 |
| 90 | Adults (18+) | NC-17, FSK 18, 18 |
| 100 | Adult Only (XXX) | R18, X18+ |

### Filtering Logic

```go
// User has max_normalized_level (from birthdate or parental setting)
// Content has ratings with normalized_levels

func (f *Filter) IsContentAllowed(ctx context.Context, contentID uuid.UUID, userLevel int) bool {
    // Get the MINIMUM normalized_level from all content ratings
    // (most restrictive rating wins)
    minLevel := f.getMinContentLevel(ctx, contentID)
    return minLevel <= userLevel
}
```

### Display Logic

```go
// User has preferred_rating_system (from locale or settings)
// Show rating in user's preferred system, fallback to any available

func (f *Filter) GetDisplayRating(ctx context.Context, contentID uuid.UUID, preferredSystem string) *Rating {
    ratings := f.getContentRatings(ctx, contentID)

    // Try preferred system first
    for _, r := range ratings {
        if r.SystemCode == preferredSystem {
            return r
        }
    }

    // Fallback to most common (MPAA, FSK, BBFC)
    // Or return nil if no rating
}
```

## User Settings

### User Table Extensions

```sql
ALTER TABLE users ADD COLUMN birthdate DATE;
ALTER TABLE users ADD COLUMN max_rating_level INT DEFAULT 100;  -- Parental override
ALTER TABLE users ADD COLUMN adult_content_enabled BOOLEAN DEFAULT false;
ALTER TABLE users ADD COLUMN preferred_rating_system VARCHAR(20);  -- 'fsk', 'mpaa'
ALTER TABLE users ADD COLUMN parental_pin_hash VARCHAR(255);
ALTER TABLE users ADD COLUMN hide_restricted BOOLEAN DEFAULT true;  -- Hide vs show locked
```

### Age Calculation

```go
func (u *User) EffectiveMaxLevel() int {
    // Calculate age from birthdate
    if u.Birthdate != nil {
        age := calculateAge(*u.Birthdate)
        ageLevel := ageToNormalizedLevel(age)

        // Parental override can only REDUCE, not increase
        if u.MaxRatingLevel < ageLevel {
            return u.MaxRatingLevel
        }
        return ageLevel
    }

    // No birthdate = use parental setting
    return u.MaxRatingLevel
}

func ageToNormalizedLevel(age int) int {
    switch {
    case age >= 18: return 100
    case age >= 16: return 75
    case age >= 12: return 50
    case age >= 6:  return 25
    default:        return 0
    }
}
```

## Person Visibility

### Logic

```go
func (f *Filter) IsPersonVisible(ctx context.Context, personID, userID uuid.UUID) bool {
    user := f.getUser(ctx, userID)
    maxLevel := user.EffectiveMaxLevel()

    // Person is visible if they appear in ANY content the user can see
    return f.personHasVisibleContent(ctx, personID, maxLevel)
}

// SQL Query
const personVisibilityQuery = `
SELECT EXISTS (
    SELECT 1 FROM media_people mp
    JOIN media_items mi ON mp.media_item_id = mi.id
    JOIN content_ratings cr ON cr.content_id = mi.id AND cr.content_type = 'media_item'
    JOIN ratings r ON cr.rating_id = r.id
    WHERE mp.person_id = $1
    AND r.normalized_level <= $2
)
`
```

### Image Handling for Persons

```sql
-- Person images have their own ratings
CREATE TABLE person_images (
    id UUID PRIMARY KEY,
    person_id UUID NOT NULL REFERENCES people(id),
    image_type VARCHAR(50) NOT NULL,  -- 'profile', 'backdrop'
    url VARCHAR(512) NOT NULL,
    rating_id UUID REFERENCES ratings(id),  -- NULL = unrated/SFW
    is_primary BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

```go
func (f *Filter) GetPersonImage(ctx context.Context, personID, userID uuid.UUID) *Image {
    user := f.getUser(ctx, userID)
    maxLevel := user.EffectiveMaxLevel()

    // Get primary image that user can see
    images := f.getPersonImages(ctx, personID)

    for _, img := range images {
        if img.IsPrimary && f.isImageAllowed(img, maxLevel) {
            return img
        }
    }

    // Fallback to any allowed image
    for _, img := range images {
        if f.isImageAllowed(img, maxLevel) {
            return img
        }
    }

    // No allowed image = return placeholder
    return f.getPlaceholderImage()
}
```

## API Behavior

### Content Listing

```go
// GET /Items always filters by user permissions
func (h *ItemsHandler) List(w http.ResponseWriter, r *http.Request) {
    user := middleware.GetUser(r.Context())
    maxLevel := user.EffectiveMaxLevel()

    items := h.service.ListItems(r.Context(), ListParams{
        MaxRatingLevel: maxLevel,
        IncludeAdult:   user.AdultContentEnabled,
    })

    // Each item in response has appropriate rating displayed
}
```

### Restricted Content Response

When content is restricted but `hide_restricted = false`:

```json
{
    "Id": "abc-123",
    "Name": "Restricted Content",
    "IsRestricted": true,
    "RequiresPin": true,
    "MinimumRatingLevel": 75,
    "ImageUrl": "/placeholder/restricted.png"
}
```

## Metadata Provider Integration

| Provider | Rating Source |
|----------|---------------|
| TMDB | Certifications by country |
| IMDB | Parents Guide (parsed) |
| MusicBrainz | - (music has separate system) |
| Stash-Box | Adult (always 100) |
| TPDB | Adult (always 100) |
| OpenLibrary | - (books have separate system) |

---

## Module-Specific Rating Systems

**Important**: Audio, Books, and Comics modules have **separate** age restriction systems from video content (movies/tvshows). These content types use different industry rating systems.

### Video Content (Movies & TV Shows)
Uses the international film/TV rating systems described above (MPAA, FSK, BBFC, etc.)

### Music (Audio Content)

Music uses the **Parental Advisory** system and regional equivalents:

| System | Region | Labels |
|--------|--------|--------|
| RIAA PAL | USA | Explicit, Clean |
| BVMI | Germany | USK-controlled |
| BPI | UK | Parental Advisory |

```sql
-- Music-specific ratings (separate table)
CREATE TABLE music_ratings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(20) NOT NULL UNIQUE,  -- 'explicit', 'clean', 'none'
    name VARCHAR(100) NOT NULL,
    description TEXT,
    min_age INT,  -- 18 for explicit
    sort_order INT NOT NULL
);

-- Per-track and per-album ratings
CREATE TABLE music_content_ratings (
    content_id UUID NOT NULL,
    content_type VARCHAR(50) NOT NULL,  -- 'album', 'track'
    rating_id UUID NOT NULL REFERENCES music_ratings(id),
    source VARCHAR(100),  -- 'musicbrainz', 'spotify', 'manual'
    PRIMARY KEY (content_id, content_type)
);
```

### Books

Books use **age range** recommendations and content warnings:

| System | Region | Ratings |
|--------|--------|---------|
| Publisher Age Range | Universal | Children, YA, Adult, Mature |
| Common Sense Media | USA | Age recommendations (0-18+) |
| Book Trust | UK | Age bands |

```sql
-- Book-specific ratings (separate table)
CREATE TABLE book_ratings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(20) NOT NULL UNIQUE,  -- 'children', 'ya', 'adult', 'mature'
    name VARCHAR(100) NOT NULL,
    age_range_min INT,  -- 0, 12, 16, 18
    age_range_max INT,  -- 12, 16, 100, 100
    sort_order INT NOT NULL
);

-- Book content warnings (separate from rating)
CREATE TABLE book_content_warnings (
    book_id UUID NOT NULL,
    warning_type VARCHAR(50) NOT NULL,  -- 'violence', 'sexual', 'language', 'drugs'
    severity VARCHAR(20) NOT NULL,      -- 'mild', 'moderate', 'graphic'
    PRIMARY KEY (book_id, warning_type)
);

CREATE TABLE book_content_ratings (
    book_id UUID NOT NULL,
    rating_id UUID NOT NULL REFERENCES book_ratings(id),
    source VARCHAR(100),  -- 'goodreads', 'manual'
    PRIMARY KEY (book_id)
);
```

### Comics

Comics use **publisher rating systems** and regional equivalents:

| System | Publisher/Region | Ratings |
|--------|------------------|---------|
| Marvel Comics | Marvel | All Ages, T, T+, Parental Advisory, MAX |
| DC Comics | DC | E, E10+, T, T+, M |
| CCA | Historical | Approved, Not Approved |
| Manga | Japan | All Ages, Shōnen, Seinen, Josei |

```sql
-- Comic-specific ratings (separate table)
CREATE TABLE comic_ratings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(20) NOT NULL UNIQUE,  -- 'all_ages', 'teen', 'teen_plus', 'mature', 'max'
    name VARCHAR(100) NOT NULL,
    publisher VARCHAR(50),  -- 'marvel', 'dc', 'manga', 'universal'
    min_age INT,
    sort_order INT NOT NULL
);

CREATE TABLE comic_content_ratings (
    comic_id UUID NOT NULL,
    rating_id UUID NOT NULL REFERENCES comic_ratings(id),
    source VARCHAR(100),  -- 'comicvine', 'manual'
    PRIMARY KEY (comic_id)
);
```

### Per-Module Filtering

Each module maintains its own filtering logic:

```go
// Video content uses shared normalized_level system
type VideoRatingFilter struct {
    maxNormalizedLevel int
}

// Music uses explicit/clean system
type MusicRatingFilter struct {
    allowExplicit bool
}

// Books use age-range system
type BookRatingFilter struct {
    userAge       int
    allowMature   bool
}

// Comics use publisher rating system
type ComicRatingFilter struct {
    maxRatingCode string  // 'teen', 'teen_plus', 'mature'
    allowAdult    bool
}
```

### User Preferences

```sql
-- User settings per module (separate from video rating preferences)
ALTER TABLE users ADD COLUMN music_allow_explicit BOOLEAN DEFAULT false;
ALTER TABLE users ADD COLUMN book_max_age_rating VARCHAR(20) DEFAULT 'adult';
ALTER TABLE users ADD COLUMN comic_max_rating VARCHAR(20) DEFAULT 'teen_plus';
```

---

## Implementation Checklist

### Phase 1: Core Infrastructure
- [ ] Create package structure `internal/service/rating/`
- [ ] Define `RatingSystem` entity (MPAA, FSK, BBFC, PEGI, etc.)
- [ ] Define `Rating` entity with normalized_level field
- [ ] Define `ContentRating` entity (polymorphic content reference)
- [ ] Define `RatingEquivalent` for cross-system mapping
- [ ] Define module-specific rating entities (music, book, comic)
- [ ] Create repository interfaces for all rating types
- [ ] Register fx module `internal/service/rating/module.go`

### Phase 2: Database
- [ ] Create migration `shared/000XXX_rating_systems.up.sql`
- [ ] Create `rating_systems` table with country codes array
- [ ] Create `ratings` table with normalized_level (0-100)
- [ ] Create `rating_equivalents` cross-reference table
- [ ] Create `content_ratings` polymorphic table
- [ ] Seed international rating systems (MPAA, FSK, BBFC, PEGI, ACB, CERO, etc.)
- [ ] Seed individual ratings with correct normalized levels
- [ ] Create `music_ratings` table (explicit/clean system)
- [ ] Create `book_ratings` table (age range system)
- [ ] Create `comic_ratings` table (publisher rating system)
- [ ] Add indexes on system_id, normalized_level, content lookups
- [ ] Generate sqlc queries for rating lookups
- [ ] Generate queries for content filtering by user level

### Phase 3: Service Layer
- [ ] Implement `RatingFilter` with `IsContentAllowed()` logic
- [ ] Implement `GetDisplayRating()` with user preference fallback
- [ ] Implement `EffectiveMaxLevel()` age calculation from birthdate
- [ ] Implement `IsPersonVisible()` for actor/person filtering
- [ ] Implement per-module filters (VideoRatingFilter, MusicRatingFilter, BookRatingFilter, ComicRatingFilter)
- [ ] Add caching for rating systems (static, long TTL)
- [ ] Add caching for user effective levels (Redis with invalidation)
- [ ] Implement parental PIN verification

### Phase 4: Background Jobs
- [ ] Create River job for rating sync from metadata providers (TMDB certifications)
- [ ] Create River job for missing rating detection
- [ ] Create job for person visibility recalculation on content changes

### Phase 5: API Integration
- [ ] Add OpenAPI schema for rating endpoints
- [ ] Implement `GET /api/v1/ratings/systems` (list rating systems)
- [ ] Implement `GET /api/v1/ratings/systems/:code` (get system with ratings)
- [ ] Add rating filter middleware for all content endpoints
- [ ] Add rating data to content responses
- [ ] Implement user rating preference endpoints
- [ ] Add parental PIN verification endpoint
- [ ] Implement restricted content response format
- [ ] Add RBAC permissions for rating management

---


<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Shared](INDEX.md)

### In This Section

- [Time-Based Access Controls](ACCESS_CONTROLS.md)
- [Tracearr Analytics Service](ANALYTICS_SERVICE.md)
- [Revenge - Client Support & Device Capabilities](CLIENT_SUPPORT.md)
- [Revenge - Internationalization (i18n)](I18N.md)
- [Library Types](LIBRARY_TYPES.md)
- [News System](NEWS_SYSTEM.md)
- [Revenge - NSFW Toggle](NSFW_TOGGLE.md)
- [Dynamic RBAC with Casbin](RBAC_CASBIN.md)

### Related Topics

- [Revenge - Architecture v2](../../architecture/01_ARCHITECTURE.md) _Architecture_
- [Revenge - Design Principles](../../architecture/02_DESIGN_PRINCIPLES.md) _Architecture_
- [Revenge - Metadata System](../../architecture/03_METADATA_SYSTEM.md) _Architecture_
- [Revenge - Player Architecture](../../architecture/04_PLAYER_ARCHITECTURE.md) _Architecture_
- [Plugin Architecture Decision](../../architecture/05_PLUGIN_ARCHITECTURE_DECISION.md) _Architecture_

### Indexes

- [Design Index](../../DESIGN_INDEX.md) - All design docs by category/topic
- [Source of Truth](../../00_SOURCE_OF_TRUTH.md) - Package versions and status

<!-- DESIGN-BREADCRUMBS-END -->

## Related

- [LIBRARY_TYPES.md](LIBRARY_TYPES.md) - Extended library types including adult content
- [PHASE1_CHECKLIST.md](PHASE1_CHECKLIST.md) - Implementation status
