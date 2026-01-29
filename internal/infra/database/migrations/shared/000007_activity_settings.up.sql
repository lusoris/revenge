-- Activity Log: Audit trail for important actions
CREATE TABLE activity_log (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID REFERENCES users(id) ON DELETE SET NULL,
    profile_id      UUID REFERENCES profiles(id) ON DELETE SET NULL,

    -- Action details
    action          VARCHAR(100) NOT NULL,           -- e.g., 'login', 'playback.start', 'library.scan'
    module          VARCHAR(50),                     -- Content module (movie, tvshow, etc.)
    item_id         UUID,                            -- Related item ID (if applicable)
    item_type       VARCHAR(50),                     -- Item type for display

    -- Context
    details         JSONB,                           -- Additional action-specific data
    ip_address      INET,
    user_agent      TEXT,

    -- Severity
    severity        VARCHAR(20) NOT NULL DEFAULT 'info',  -- debug, info, warn, error

    -- Timestamp
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for common queries
CREATE INDEX idx_activity_log_user ON activity_log(user_id) WHERE user_id IS NOT NULL;
CREATE INDEX idx_activity_log_created ON activity_log(created_at DESC);
CREATE INDEX idx_activity_log_action ON activity_log(action);
CREATE INDEX idx_activity_log_module ON activity_log(module) WHERE module IS NOT NULL;

-- Partitioning hint: Consider partitioning by month for large deployments
-- CREATE TABLE activity_log_YYYY_MM PARTITION OF activity_log FOR VALUES FROM (...) TO (...);

-- Server Settings: Persisted configuration
CREATE TABLE server_settings (
    key             VARCHAR(255) PRIMARY KEY,
    value           JSONB NOT NULL,
    description     TEXT,
    updated_by      UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Insert default settings
INSERT INTO server_settings (key, value, description) VALUES
    ('server.name', '"Revenge Media Server"', 'Server display name'),
    ('server.public', 'false', 'Allow public access without login'),
    ('server.registration_enabled', 'false', 'Allow new user registration'),
    ('server.default_language', '"en"', 'Default language for new users'),
    ('transcoding.enabled', 'true', 'Enable transcoding via Blackbeard'),
    ('transcoding.blackbeard_url', '""', 'Blackbeard service URL'),
    ('library.scan_on_startup', 'true', 'Scan libraries on server start'),
    ('library.watch_filesystem', 'true', 'Watch for file changes'),
    ('cache.enabled', 'true', 'Enable caching layer'),
    ('search.enabled', 'true', 'Enable Typesense search'),
    ('adult.globally_enabled', 'false', 'Allow adult content access server-wide');

-- Genres: Domain-scoped genre definitions
CREATE TYPE genre_domain AS ENUM (
    'movie',
    'tv',
    'music',
    'book',
    'podcast',
    'game',
    'adult'
);

CREATE TABLE genres (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    domain          genre_domain NOT NULL,
    name            VARCHAR(100) NOT NULL,
    slug            VARCHAR(100) NOT NULL,
    description     TEXT,
    parent_id       UUID REFERENCES genres(id) ON DELETE SET NULL,
    external_ids    JSONB NOT NULL DEFAULT '{}'::jsonb,  -- {tmdb: "28", musicbrainz: "..."}

    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (domain, slug)
);

CREATE INDEX idx_genres_domain ON genres(domain);
CREATE INDEX idx_genres_parent ON genres(parent_id) WHERE parent_id IS NOT NULL;

CREATE TRIGGER genres_updated_at
    BEFORE UPDATE ON genres
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- People: Shared people table (actors, directors, artists, authors)
CREATE TABLE people (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) NOT NULL,
    sort_name       VARCHAR(255),                    -- For sorting (e.g., "Spielberg, Steven")
    original_name   VARCHAR(255),                    -- Name in original language

    -- Bio
    biography       TEXT,
    birthdate       DATE,
    deathdate       DATE,
    birthplace      VARCHAR(255),
    gender          VARCHAR(20),

    -- Images
    primary_image_url   TEXT,
    primary_image_blurhash VARCHAR(50),

    -- External IDs
    tmdb_id         INT,
    imdb_id         VARCHAR(20),
    tvdb_id         INT,
    musicbrainz_id  UUID,

    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_people_name ON people(name);
CREATE INDEX idx_people_tmdb ON people(tmdb_id) WHERE tmdb_id IS NOT NULL;
CREATE INDEX idx_people_imdb ON people(imdb_id) WHERE imdb_id IS NOT NULL;

CREATE TRIGGER people_updated_at
    BEFORE UPDATE ON people
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
