-- Video People: Shared between movie and tvshow modules
-- Data overlaps 100% after background worker enrichment (TMDB, TVDB, IMDB)
BEGIN;

CREATE TABLE video_people (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name                VARCHAR(255) NOT NULL,
    sort_name           VARCHAR(255),
    original_name       VARCHAR(255),

    -- Bio
    biography           TEXT,
    birthdate           DATE,
    deathdate           DATE,
    birthplace          VARCHAR(255),
    gender              VARCHAR(20),

    -- Images
    primary_image_url   TEXT,
    primary_image_blurhash VARCHAR(50),

    -- External IDs (shared across video providers)
    tmdb_id             INT,
    imdb_id             VARCHAR(20),
    tvdb_id             INT,

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_video_people_name ON video_people(name);
CREATE INDEX idx_video_people_sort_name ON video_people(sort_name) WHERE sort_name IS NOT NULL;
CREATE INDEX idx_video_people_tmdb ON video_people(tmdb_id) WHERE tmdb_id IS NOT NULL;
CREATE INDEX idx_video_people_imdb ON video_people(imdb_id) WHERE imdb_id IS NOT NULL;
CREATE INDEX idx_video_people_tvdb ON video_people(tvdb_id) WHERE tvdb_id IS NOT NULL;

CREATE TRIGGER video_people_updated_at
    BEFORE UPDATE ON video_people
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- Credit role enum (shared between movie and tvshow)
CREATE TYPE video_credit_role AS ENUM (
    'actor',
    'director',
    'writer',
    'creator',           -- TV: series creator
    'showrunner',        -- TV: showrunner
    'producer',
    'executive_producer',
    'composer',
    'cinematographer',
    'editor',
    'production_designer',
    'costume_designer',
    'makeup_artist',
    'visual_effects',
    'stunt_coordinator',
    'sound_designer',
    'guest_star'         -- TV: guest appearances
);

COMMIT;
