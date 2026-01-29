-- TV Libraries: Per-module library table
-- Replaces dependency on shared libraries table
BEGIN;

-- TV-specific library table
CREATE TABLE tv_libraries (
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
    tvdb_enabled            BOOLEAN NOT NULL DEFAULT true,
    download_backdrops      BOOLEAN NOT NULL DEFAULT true,
    download_nfo            BOOLEAN NOT NULL DEFAULT false,
    generate_chapters       BOOLEAN NOT NULL DEFAULT true,

    -- TV-specific settings
    season_folder_format    VARCHAR(100) DEFAULT 'Season {season}',
    episode_naming_format   VARCHAR(255) DEFAULT '{series} - S{season:00}E{episode:00} - {title}',
    auto_add_missing        BOOLEAN NOT NULL DEFAULT false,

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
CREATE INDEX idx_tv_libraries_owner ON tv_libraries(owner_user_id) WHERE owner_user_id IS NOT NULL;

-- Trigger for updated_at
CREATE TRIGGER tv_libraries_updated_at
    BEFORE UPDATE ON tv_libraries
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- TV library user access: Which users can access which libraries
CREATE TABLE tv_library_access (
    library_id      UUID NOT NULL REFERENCES tv_libraries(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    can_manage      BOOLEAN NOT NULL DEFAULT false,

    PRIMARY KEY (library_id, user_id)
);

CREATE INDEX idx_tv_library_access_user ON tv_library_access(user_id);

-- Add new FK column referencing tv_libraries (nullable initially for migration)
ALTER TABLE series ADD COLUMN tv_library_id UUID REFERENCES tv_libraries(id) ON DELETE CASCADE;

-- Create index for the new FK
CREATE INDEX idx_series_tv_library ON series(tv_library_id) WHERE tv_library_id IS NOT NULL;

COMMIT;
