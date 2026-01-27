-- 000011_rating_systems.up.sql
-- Content rating systems and age restriction support

-- Rating systems table (seeded with known systems like MPAA, FSK, BBFC, etc.)
CREATE TABLE rating_systems (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(20) NOT NULL UNIQUE,           -- 'mpaa', 'fsk', 'bbfc'
    name VARCHAR(100) NOT NULL,                  -- 'Motion Picture Association'
    country_codes TEXT[] NOT NULL DEFAULT '{}', -- ['US', 'CA']
    is_active BOOLEAN NOT NULL DEFAULT true,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_rating_systems_code ON rating_systems(code);
CREATE INDEX idx_rating_systems_active ON rating_systems(is_active) WHERE is_active = true;

-- Individual ratings within each system
CREATE TABLE ratings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    system_id UUID NOT NULL REFERENCES rating_systems(id) ON DELETE CASCADE,
    code VARCHAR(20) NOT NULL,                   -- 'PG-13', 'FSK 16'
    name VARCHAR(100) NOT NULL,                  -- 'Parental Guidance 13'
    description TEXT,
    min_age INT,                                 -- Minimum age (0, 6, 12, 16, 18)
    normalized_level INT NOT NULL,              -- 0-100 scale for cross-system comparison
    sort_order INT NOT NULL DEFAULT 0,
    is_adult BOOLEAN NOT NULL DEFAULT false,    -- Explicit adult content flag
    icon_url VARCHAR(512),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(system_id, code),
    CONSTRAINT ratings_normalized_level_check CHECK (normalized_level >= 0 AND normalized_level <= 100)
);

CREATE INDEX idx_ratings_system_id ON ratings(system_id);
CREATE INDEX idx_ratings_normalized_level ON ratings(normalized_level);
CREATE INDEX idx_ratings_is_adult ON ratings(is_adult) WHERE is_adult = true;

-- Cross-reference equivalents between rating systems (for display)
CREATE TABLE rating_equivalents (
    rating_id UUID NOT NULL REFERENCES ratings(id) ON DELETE CASCADE,
    equivalent_rating_id UUID NOT NULL REFERENCES ratings(id) ON DELETE CASCADE,
    PRIMARY KEY (rating_id, equivalent_rating_id),
    CONSTRAINT rating_equivalents_not_self CHECK (rating_id != equivalent_rating_id)
);

-- Content ratings - many-to-many between content and ratings
-- Supports multiple ratings per content (from different systems)
CREATE TABLE content_ratings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content_id UUID NOT NULL,                    -- FK to media_items, images, etc.
    content_type VARCHAR(50) NOT NULL,           -- 'media_item', 'image', 'person_image'
    rating_id UUID NOT NULL REFERENCES ratings(id) ON DELETE CASCADE,
    source VARCHAR(100),                         -- 'tmdb', 'manual', 'imdb', 'stash-box'
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(content_id, rating_id)
);

CREATE INDEX idx_content_ratings_content ON content_ratings(content_id, content_type);
CREATE INDEX idx_content_ratings_rating_id ON content_ratings(rating_id);

-- Materialized view for efficient content level lookups
-- Returns the MINIMUM (most restrictive) normalized level for each content
CREATE MATERIALIZED VIEW content_min_rating_levels AS
SELECT
    cr.content_id,
    cr.content_type,
    MIN(r.normalized_level) AS min_level,
    BOOL_OR(r.is_adult) AS is_adult
FROM content_ratings cr
JOIN ratings r ON cr.rating_id = r.id
GROUP BY cr.content_id, cr.content_type;

CREATE UNIQUE INDEX idx_content_min_rating_levels_pk ON content_min_rating_levels(content_id, content_type);
CREATE INDEX idx_content_min_rating_levels_level ON content_min_rating_levels(min_level);
CREATE INDEX idx_content_min_rating_levels_adult ON content_min_rating_levels(is_adult) WHERE is_adult = true;

-- Function to refresh the materialized view
CREATE OR REPLACE FUNCTION refresh_content_min_rating_levels()
RETURNS TRIGGER AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY content_min_rating_levels;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Trigger to refresh materialized view on content_ratings changes
CREATE TRIGGER trigger_refresh_content_min_rating_levels
AFTER INSERT OR UPDATE OR DELETE ON content_ratings
FOR EACH STATEMENT
EXECUTE FUNCTION refresh_content_min_rating_levels();
