-- Library types enum
CREATE TYPE library_type AS ENUM (
    'movie',
    'tvshow',
    'music',
    'audiobook',
    'book',
    'podcast',
    'photo',
    'livetv',
    'comics',
    'adult_movie',
    'adult_scene'
);

-- Libraries table: Content library definitions
CREATE TABLE libraries (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) NOT NULL,
    type            library_type NOT NULL,
    paths           TEXT[] NOT NULL,                 -- Filesystem paths to scan

    -- Scanning settings
    scan_enabled        BOOLEAN NOT NULL DEFAULT true,
    scan_interval_hours INT NOT NULL DEFAULT 24,
    last_scan_at        TIMESTAMPTZ,
    last_scan_duration  INTERVAL,

    -- Metadata settings
    preferred_language      VARCHAR(10) DEFAULT 'en',
    download_images         BOOLEAN NOT NULL DEFAULT true,
    download_nfo            BOOLEAN NOT NULL DEFAULT false,
    generate_chapters       BOOLEAN NOT NULL DEFAULT true,

    -- Access control
    is_private          BOOLEAN NOT NULL DEFAULT false,  -- Only visible to owner
    owner_user_id       UUID REFERENCES users(id) ON DELETE SET NULL,

    -- Display
    sort_order          INT NOT NULL DEFAULT 0,
    icon                VARCHAR(50),

    -- Timestamps
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_libraries_type ON libraries(type);
CREATE INDEX idx_libraries_owner ON libraries(owner_user_id) WHERE owner_user_id IS NOT NULL;

-- Trigger for updated_at
CREATE TRIGGER libraries_updated_at
    BEFORE UPDATE ON libraries
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- Library user access: Which users can access which libraries
CREATE TABLE library_user_access (
    library_id      UUID NOT NULL REFERENCES libraries(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    can_manage      BOOLEAN NOT NULL DEFAULT false,  -- Can edit library settings

    PRIMARY KEY (library_id, user_id)
);

CREATE INDEX idx_library_user_access_user ON library_user_access(user_id);
