-- Movie Libraries: Per-module library table
-- Replaces dependency on shared libraries table
BEGIN;

-- Movie-specific library table
CREATE TABLE movie_libraries (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name                    VARCHAR(255) NOT NULL,
    paths                   TEXT[] NOT NULL,

    -- Scanning settings
    scan_enabled            BOOLEAN NOT NULL DEFAULT true,
    scan_interval_hours     INT NOT NULL DEFAULT 24,
    last_scan_at            TIMESTAMPTZ,
    last_scan_duration      INTERVAL,

    -- Metadata settings
    preferred_language      VARCHAR(10) DEFAULT 'en',
    tmdb_enabled            BOOLEAN NOT NULL DEFAULT true,
    imdb_enabled            BOOLEAN NOT NULL DEFAULT true,
    download_trailers       BOOLEAN NOT NULL DEFAULT false,
    download_backdrops      BOOLEAN NOT NULL DEFAULT true,
    download_nfo            BOOLEAN NOT NULL DEFAULT false,
    generate_chapters       BOOLEAN NOT NULL DEFAULT true,

    -- Access control
    is_private              BOOLEAN NOT NULL DEFAULT false,
    owner_user_id           UUID REFERENCES users(id) ON DELETE SET NULL,

    -- Display
    sort_order              INT NOT NULL DEFAULT 0,
    icon                    VARCHAR(50),

    -- Timestamps
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_movie_libraries_owner ON movie_libraries(owner_user_id) WHERE owner_user_id IS NOT NULL;

-- Trigger for updated_at
CREATE TRIGGER movie_libraries_updated_at
    BEFORE UPDATE ON movie_libraries
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- Movie library user access: Which users can access which libraries
CREATE TABLE movie_library_access (
    library_id      UUID NOT NULL REFERENCES movie_libraries(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    can_manage      BOOLEAN NOT NULL DEFAULT false,

    PRIMARY KEY (library_id, user_id)
);

CREATE INDEX idx_movie_library_access_user ON movie_library_access(user_id);

-- Add new FK column referencing movie_libraries (nullable initially for migration)
ALTER TABLE movies ADD COLUMN movie_library_id UUID REFERENCES movie_libraries(id) ON DELETE CASCADE;

-- Create index for the new FK
CREATE INDEX idx_movies_movie_library ON movies(movie_library_id) WHERE movie_library_id IS NOT NULL;

COMMIT;
