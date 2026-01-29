-- Collections - curated groups of related content
-- Different from playlists: collections are typically system/admin created
-- and represent logical groupings (franchises, series, themes)

CREATE TYPE collection_type AS ENUM (
    'video',        -- Movies, shows
    'audio',        -- Albums, artists, playlists
    'mixed',        -- Cross-type collections
    'franchise',    -- Movie/TV franchises (MCU, Star Wars)
    'box_set',      -- Box sets / special editions
    'smart'         -- Auto-generated based on rules
);

-- Main collections table
CREATE TABLE collections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Basic info
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE,
    description TEXT,
    collection_type collection_type NOT NULL,

    -- Display
    poster_url TEXT,
    backdrop_url TEXT,
    logo_url TEXT,
    theme_color VARCHAR(7), -- Hex color for UI theming

    -- Metadata
    item_count INT NOT NULL DEFAULT 0,
    total_duration_ms BIGINT NOT NULL DEFAULT 0,

    -- External IDs (for matching with metadata providers)
    tmdb_collection_id INT,
    tvdb_id INT,
    imdb_id VARCHAR(20),

    -- Smart collection rules (JSON for flexibility)
    -- Example: {"genres": ["action"], "year_min": 2000, "rating_min": 7.0}
    smart_rules JSONB,

    -- Visibility
    is_visible BOOLEAN NOT NULL DEFAULT TRUE,
    sort_order INT NOT NULL DEFAULT 0,

    -- Ownership (NULL = system collection)
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Collection items - content in collections
CREATE TABLE collection_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    collection_id UUID NOT NULL REFERENCES collections(id) ON DELETE CASCADE,

    -- Content reference (polymorphic)
    content_type VARCHAR(50) NOT NULL, -- movie, series, album, artist, etc.
    content_id UUID NOT NULL,

    -- Ordering
    position INT NOT NULL,

    -- Item-specific metadata
    release_order INT,          -- Original release order
    chronological_order INT,    -- In-universe chronological order
    custom_title VARCHAR(500),  -- Override title for this collection context
    notes TEXT,                 -- Notes about this item in collection context

    -- Timestamps
    added_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (collection_id, content_type, content_id)
);

-- Collection tags for categorization
CREATE TABLE collection_tags (
    collection_id UUID NOT NULL REFERENCES collections(id) ON DELETE CASCADE,
    tag VARCHAR(100) NOT NULL,
    PRIMARY KEY (collection_id, tag)
);

-- User collection subscriptions (for notifications, etc.)
CREATE TABLE collection_subscriptions (
    collection_id UUID NOT NULL REFERENCES collections(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    notify_new_items BOOLEAN NOT NULL DEFAULT TRUE,
    subscribed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (collection_id, user_id)
);

-- Indexes
CREATE INDEX idx_collections_type ON collections(collection_type);
CREATE INDEX idx_collections_slug ON collections(slug);
CREATE INDEX idx_collections_tmdb ON collections(tmdb_collection_id) WHERE tmdb_collection_id IS NOT NULL;
CREATE INDEX idx_collections_visible ON collections(is_visible, sort_order);
CREATE INDEX idx_collection_items_collection ON collection_items(collection_id, position);
CREATE INDEX idx_collection_items_content ON collection_items(content_type, content_id);
CREATE INDEX idx_collection_tags_tag ON collection_tags(tag);
CREATE INDEX idx_collection_subscriptions_user ON collection_subscriptions(user_id);

-- Trigger to update collection metadata on item changes
CREATE OR REPLACE FUNCTION update_collection_metadata()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE collections
    SET
        item_count = (SELECT COUNT(*) FROM collection_items WHERE collection_id = COALESCE(NEW.collection_id, OLD.collection_id)),
        updated_at = NOW()
    WHERE id = COALESCE(NEW.collection_id, OLD.collection_id);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_collection_items_update_metadata
AFTER INSERT OR UPDATE OR DELETE ON collection_items
FOR EACH ROW EXECUTE FUNCTION update_collection_metadata();
