-- TV Shows: Core tables (series, seasons, episodes)
BEGIN;

-- TV Series: Main series table
CREATE TABLE series (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    library_id          UUID NOT NULL REFERENCES libraries(id) ON DELETE CASCADE,

    -- Core metadata
    title               VARCHAR(500) NOT NULL,
    sort_title          VARCHAR(500),
    original_title      VARCHAR(500),
    tagline             TEXT,
    overview            TEXT,

    -- Airing info
    first_air_date      DATE,
    last_air_date       DATE,
    year                INT,                     -- First air year
    status              VARCHAR(50),             -- Returning Series, Ended, Canceled, etc.
    type                VARCHAR(50),             -- Scripted, Documentary, Reality, etc.

    -- Ratings
    content_rating      VARCHAR(20),             -- TV-14, TV-MA, etc.
    rating_level        INT DEFAULT 0,           -- 0-100 normalized content level
    community_rating    DECIMAL(3,1),            -- 0.0 - 10.0
    vote_count          INT,

    -- Counts (cached for performance)
    season_count        INT NOT NULL DEFAULT 0,
    episode_count       INT NOT NULL DEFAULT 0,
    special_count       INT NOT NULL DEFAULT 0,

    -- Primary images
    poster_path         TEXT,
    poster_blurhash     VARCHAR(50),
    backdrop_path       TEXT,
    backdrop_blurhash   VARCHAR(50),
    logo_path           TEXT,

    -- External IDs
    tmdb_id             INT UNIQUE,
    imdb_id             VARCHAR(20) UNIQUE,
    tvdb_id             INT UNIQUE,

    -- Network info (primary network)
    network_name        VARCHAR(255),
    network_logo_path   TEXT,

    -- Status
    date_added          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_played_at      TIMESTAMPTZ,
    is_locked           BOOLEAN NOT NULL DEFAULT false,

    -- Timestamps
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for series
CREATE INDEX idx_series_library ON series(library_id);
CREATE INDEX idx_series_title ON series(title);
CREATE INDEX idx_series_sort_title ON series(sort_title);
CREATE INDEX idx_series_year ON series(year);
CREATE INDEX idx_series_first_air ON series(first_air_date);
CREATE INDEX idx_series_rating ON series(community_rating DESC NULLS LAST);
CREATE INDEX idx_series_date_added ON series(date_added DESC);
CREATE INDEX idx_series_status ON series(status);
CREATE INDEX idx_series_tmdb ON series(tmdb_id) WHERE tmdb_id IS NOT NULL;
CREATE INDEX idx_series_imdb ON series(imdb_id) WHERE imdb_id IS NOT NULL;
CREATE INDEX idx_series_tvdb ON series(tvdb_id) WHERE tvdb_id IS NOT NULL;

-- Full-text search
CREATE INDEX idx_series_search ON series USING GIN (
    to_tsvector('english', COALESCE(title, '') || ' ' || COALESCE(original_title, '') || ' ' || COALESCE(overview, ''))
);

CREATE TRIGGER series_updated_at
    BEFORE UPDATE ON series
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- Seasons
CREATE TABLE seasons (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    series_id           UUID NOT NULL REFERENCES series(id) ON DELETE CASCADE,

    -- Season info
    season_number       INT NOT NULL,
    name                VARCHAR(255),            -- "Season 1" or custom name
    overview            TEXT,

    -- Air dates
    air_date            DATE,
    year                INT,

    -- Episode count (cached)
    episode_count       INT NOT NULL DEFAULT 0,

    -- Images
    poster_path         TEXT,
    poster_blurhash     VARCHAR(50),

    -- External IDs
    tmdb_id             INT,
    tvdb_id             INT,

    -- Timestamps
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (series_id, season_number)
);

CREATE INDEX idx_seasons_series ON seasons(series_id);
CREATE INDEX idx_seasons_number ON seasons(series_id, season_number);
CREATE INDEX idx_seasons_air_date ON seasons(air_date);

CREATE TRIGGER seasons_updated_at
    BEFORE UPDATE ON seasons
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- Episodes
CREATE TABLE episodes (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    series_id           UUID NOT NULL REFERENCES series(id) ON DELETE CASCADE,
    season_id           UUID NOT NULL REFERENCES seasons(id) ON DELETE CASCADE,

    -- File info
    path                TEXT NOT NULL UNIQUE,
    container           VARCHAR(20),
    size_bytes          BIGINT,
    runtime_ticks       BIGINT,

    -- Episode info
    season_number       INT NOT NULL,
    episode_number      INT NOT NULL,
    absolute_number     INT,                     -- For anime ordering

    -- Metadata
    title               VARCHAR(500) NOT NULL,
    overview            TEXT,
    production_code     VARCHAR(50),

    -- Air info
    air_date            DATE,
    air_date_utc        TIMESTAMPTZ,

    -- Ratings
    community_rating    DECIMAL(3,1),
    vote_count          INT,

    -- Images
    still_path          TEXT,
    still_blurhash      VARCHAR(50),

    -- External IDs
    tmdb_id             INT,
    imdb_id             VARCHAR(20),
    tvdb_id             INT,

    -- Status
    date_added          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_played_at      TIMESTAMPTZ,
    play_count          INT NOT NULL DEFAULT 0,
    is_locked           BOOLEAN NOT NULL DEFAULT false,

    -- Timestamps
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (series_id, season_number, episode_number)
);

-- Indexes for episodes
CREATE INDEX idx_episodes_series ON episodes(series_id);
CREATE INDEX idx_episodes_season ON episodes(season_id);
CREATE INDEX idx_episodes_number ON episodes(series_id, season_number, episode_number);
CREATE INDEX idx_episodes_absolute ON episodes(series_id, absolute_number) WHERE absolute_number IS NOT NULL;
CREATE INDEX idx_episodes_air_date ON episodes(air_date);
CREATE INDEX idx_episodes_date_added ON episodes(date_added DESC);
CREATE INDEX idx_episodes_tmdb ON episodes(tmdb_id) WHERE tmdb_id IS NOT NULL;
CREATE INDEX idx_episodes_tvdb ON episodes(tvdb_id) WHERE tvdb_id IS NOT NULL;

CREATE TRIGGER episodes_updated_at
    BEFORE UPDATE ON episodes
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- TV Networks
CREATE TABLE tv_networks (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name                VARCHAR(255) NOT NULL,
    logo_path           TEXT,
    origin_country      VARCHAR(10),

    -- External IDs
    tmdb_id             INT UNIQUE,

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tv_networks_name ON tv_networks(name);

-- Series <-> Network junction
CREATE TABLE series_network_link (
    series_id           UUID NOT NULL REFERENCES series(id) ON DELETE CASCADE,
    network_id          UUID NOT NULL REFERENCES tv_networks(id) ON DELETE CASCADE,
    display_order       INT NOT NULL DEFAULT 0,

    PRIMARY KEY (series_id, network_id)
);

CREATE INDEX idx_series_network_link_network ON series_network_link(network_id);

-- Function to update series episode/season counts
CREATE OR REPLACE FUNCTION update_series_counts()
RETURNS TRIGGER AS $$
BEGIN
    -- Update season episode count
    IF TG_TABLE_NAME = 'episodes' THEN
        UPDATE seasons SET
            episode_count = (SELECT COUNT(*) FROM episodes WHERE season_id = COALESCE(NEW.season_id, OLD.season_id))
        WHERE id = COALESCE(NEW.season_id, OLD.season_id);
    END IF;

    -- Update series counts
    UPDATE series SET
        season_count = (SELECT COUNT(*) FROM seasons WHERE series_id = COALESCE(NEW.series_id, OLD.series_id) AND season_number > 0),
        episode_count = (SELECT COUNT(*) FROM episodes WHERE series_id = COALESCE(NEW.series_id, OLD.series_id) AND season_number > 0),
        special_count = (SELECT COUNT(*) FROM episodes WHERE series_id = COALESCE(NEW.series_id, OLD.series_id) AND season_number = 0)
    WHERE id = COALESCE(NEW.series_id, OLD.series_id);

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Triggers to maintain counts
CREATE TRIGGER episodes_update_counts
    AFTER INSERT OR UPDATE OR DELETE ON episodes
    FOR EACH ROW EXECUTE FUNCTION update_series_counts();

CREATE TRIGGER seasons_update_counts
    AFTER INSERT OR UPDATE OR DELETE ON seasons
    FOR EACH ROW EXECUTE FUNCTION update_series_counts();

COMMIT;
