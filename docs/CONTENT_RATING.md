# Content Rating System

> Universal age restriction and content rating system for Jellyfin Go.

## Overview

Jellyfin Go implements a comprehensive content rating system that:
- Supports **international rating systems** (MPAA, FSK, BBFC, PEGI, etc.)
- Filters content based on **user age/permissions**
- Handles **person visibility** based on their content appearances
- Manages **image ratings** separately (SFW/NSFW)

## Design Principles

### Private Server Model

Jellyfin is a **private, invite-only** server. There is no public registration, so:
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
    created_at TIMESTAMPTZ DEFAULT NOW()
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
    created_at TIMESTAMPTZ DEFAULT NOW(),
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
    created_at TIMESTAMPTZ DEFAULT NOW()
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
| MusicBrainz | - (music has no ratings) |
| Stash-Box | Adult (always 100) |
| TPDB | Adult (always 100) |
| OpenLibrary | - (books need manual) |

## See Also

- [LIBRARY_TYPES.md](LIBRARY_TYPES.md) - Extended library types including adult content
- [PHASE1_CHECKLIST.md](PHASE1_CHECKLIST.md) - Implementation status
