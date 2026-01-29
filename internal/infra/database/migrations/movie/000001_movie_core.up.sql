-- Movies: Core movie table
BEGIN;

CREATE TABLE movies (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    library_id          UUID NOT NULL REFERENCES libraries(id) ON DELETE CASCADE,

    -- File info
    path                TEXT NOT NULL UNIQUE,
    container           VARCHAR(20),
    size_bytes          BIGINT,
    runtime_ticks       BIGINT,          -- Duration in ticks (100ns units)

    -- Core metadata
    title               VARCHAR(500) NOT NULL,
    sort_title          VARCHAR(500),
    original_title      VARCHAR(500),
    tagline             TEXT,
    overview            TEXT,

    -- Release info
    release_date        DATE,
    year                INT,
    content_rating      VARCHAR(20),     -- PG-13, R, etc.
    rating_level        INT DEFAULT 0,   -- 0-100 normalized content level

    -- Financial (optional, from TMDb)
    budget              BIGINT,
    revenue             BIGINT,

    -- Community ratings
    community_rating    DECIMAL(3,1),    -- 0.0 - 10.0
    vote_count          INT,
    critic_rating       DECIMAL(3,1),    -- From external sources
    critic_count        INT,

    -- Primary images (cached locally)
    poster_path         TEXT,
    poster_blurhash     VARCHAR(50),
    backdrop_path       TEXT,
    backdrop_blurhash   VARCHAR(50),
    logo_path           TEXT,

    -- External IDs
    tmdb_id             INT UNIQUE,
    imdb_id             VARCHAR(20) UNIQUE,
    tvdb_id             INT,

    -- Status
    date_added          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_played_at      TIMESTAMPTZ,
    play_count          INT NOT NULL DEFAULT 0,
    is_locked           BOOLEAN NOT NULL DEFAULT false,  -- Prevent metadata refresh

    -- Timestamps
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for common queries
CREATE INDEX idx_movies_library ON movies(library_id);
CREATE INDEX idx_movies_title ON movies(title);
CREATE INDEX idx_movies_sort_title ON movies(sort_title);
CREATE INDEX idx_movies_year ON movies(year);
CREATE INDEX idx_movies_release_date ON movies(release_date);
CREATE INDEX idx_movies_rating ON movies(community_rating DESC NULLS LAST);
CREATE INDEX idx_movies_date_added ON movies(date_added DESC);
CREATE INDEX idx_movies_tmdb ON movies(tmdb_id) WHERE tmdb_id IS NOT NULL;
CREATE INDEX idx_movies_imdb ON movies(imdb_id) WHERE imdb_id IS NOT NULL;

-- Full-text search index
CREATE INDEX idx_movies_search ON movies USING GIN (
    to_tsvector('english', COALESCE(title, '') || ' ' || COALESCE(original_title, '') || ' ' || COALESCE(overview, ''))
);

-- Trigger for updated_at
CREATE TRIGGER movies_updated_at
    BEFORE UPDATE ON movies
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- Movie Collections: Groups of movies (e.g., "The Dark Knight Trilogy")
CREATE TABLE movie_collections (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name                VARCHAR(255) NOT NULL,
    sort_name           VARCHAR(255),
    overview            TEXT,
    poster_path         TEXT,
    poster_blurhash     VARCHAR(50),
    backdrop_path       TEXT,
    backdrop_blurhash   VARCHAR(50),

    -- External IDs
    tmdb_id             INT UNIQUE,

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_movie_collections_name ON movie_collections(name);

CREATE TRIGGER movie_collections_updated_at
    BEFORE UPDATE ON movie_collections
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- Movie to Collection relationship
ALTER TABLE movies ADD COLUMN collection_id UUID REFERENCES movie_collections(id) ON DELETE SET NULL;
ALTER TABLE movies ADD COLUMN collection_order INT;

CREATE INDEX idx_movies_collection ON movies(collection_id) WHERE collection_id IS NOT NULL;

-- Movie Studios
CREATE TABLE movie_studios (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name                VARCHAR(255) NOT NULL,
    logo_path           TEXT,

    -- External IDs
    tmdb_id             INT UNIQUE,

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_movie_studios_name ON movie_studios(name);

-- Junction: Movie <-> Studio
CREATE TABLE movie_studio_link (
    movie_id            UUID NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    studio_id           UUID NOT NULL REFERENCES movie_studios(id) ON DELETE CASCADE,
    display_order       INT NOT NULL DEFAULT 0,

    PRIMARY KEY (movie_id, studio_id)
);

CREATE INDEX idx_movie_studio_link_studio ON movie_studio_link(studio_id);

COMMIT;
