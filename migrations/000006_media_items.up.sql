-- 000006_media_items.up.sql
-- Media items - movies, series, episodes, music, photos

CREATE TYPE media_type AS ENUM (
    'movie', 'series', 'season', 'episode',
    'artist', 'album', 'audio',
    'photo', 'photo_album',
    'folder'
);

CREATE TABLE media_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    library_id UUID NOT NULL REFERENCES libraries(id) ON DELETE CASCADE,
    parent_id UUID REFERENCES media_items(id) ON DELETE CASCADE,
    type media_type NOT NULL,
    name VARCHAR(500) NOT NULL,
    sort_name VARCHAR(500),
    path TEXT,                                -- Filesystem path (NULL for virtual items)

    -- Common metadata
    overview TEXT,
    tagline VARCHAR(500),
    year INT,
    premiere_date DATE,
    end_date DATE,
    runtime_ticks BIGINT,                     -- Duration in ticks (100ns units)

    -- Series/Episode specific
    season_number INT,
    episode_number INT,
    absolute_episode_number INT,

    -- Music specific
    album_artist VARCHAR(255),
    track_number INT,
    disc_number INT,

    -- Ratings
    community_rating DECIMAL(3,1),            -- e.g., 8.5
    critic_rating DECIMAL(3,1),

    -- External IDs
    provider_ids JSONB NOT NULL DEFAULT '{}', -- {"imdb": "tt123", "tmdb": "456", ...}

    -- Metadata
    genres TEXT[] NOT NULL DEFAULT '{}',
    tags TEXT[] NOT NULL DEFAULT '{}',
    studios TEXT[] NOT NULL DEFAULT '{}',

    -- File info (for actual media files)
    container VARCHAR(50),                    -- mkv, mp4, etc.
    video_codec VARCHAR(50),
    audio_codec VARCHAR(50),
    width INT,
    height INT,
    bitrate INT,

    -- Timestamps
    date_created TIMESTAMPTZ,                 -- File creation date
    date_modified TIMESTAMPTZ,                -- File modification date
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Core indexes
CREATE INDEX idx_media_items_library_id ON media_items(library_id);
CREATE INDEX idx_media_items_parent_id ON media_items(parent_id);
CREATE INDEX idx_media_items_type ON media_items(type);
CREATE INDEX idx_media_items_name ON media_items(name);
CREATE INDEX idx_media_items_path ON media_items(path) WHERE path IS NOT NULL;
CREATE INDEX idx_media_items_year ON media_items(year) WHERE year IS NOT NULL;

-- Composite indexes for common queries
CREATE INDEX idx_media_items_library_type ON media_items(library_id, type);
CREATE INDEX idx_media_items_parent_type ON media_items(parent_id, type);

-- Full-text search index (PostgreSQL generated column)
ALTER TABLE media_items ADD COLUMN search_vector tsvector
    GENERATED ALWAYS AS (
        setweight(to_tsvector('english', coalesce(name, '')), 'A') ||
        setweight(to_tsvector('english', coalesce(overview, '')), 'B') ||
        setweight(to_tsvector('english', coalesce(tagline, '')), 'C')
    ) STORED;

CREATE INDEX idx_media_items_search ON media_items USING GIN(search_vector);

CREATE TRIGGER update_media_items_updated_at
    BEFORE UPDATE ON media_items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
