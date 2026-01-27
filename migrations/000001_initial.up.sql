-- 001_initial.up.sql
-- Initial schema for Jellyfin Go
-- PostgreSQL 18+

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- =============================================================================
-- USERS
-- =============================================================================

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) UNIQUE,
    password_hash VARCHAR(255),  -- NULL for OIDC-only users
    display_name VARCHAR(255),
    is_admin BOOLEAN NOT NULL DEFAULT false,
    is_disabled BOOLEAN NOT NULL DEFAULT false,
    last_login_at TIMESTAMPTZ,
    last_activity_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email) WHERE email IS NOT NULL;
CREATE INDEX idx_users_is_admin ON users(is_admin) WHERE is_admin = true;

-- =============================================================================
-- SESSIONS
-- =============================================================================

CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(64) NOT NULL UNIQUE,  -- SHA-256 of access token
    refresh_token_hash VARCHAR(64) UNIQUE,   -- SHA-256 of refresh token
    device_id VARCHAR(255),
    device_name VARCHAR(255),
    client_name VARCHAR(255),
    client_version VARCHAR(50),
    ip_address INET,
    expires_at TIMESTAMPTZ NOT NULL,
    refresh_expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_token_hash ON sessions(token_hash);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);

-- =============================================================================
-- OIDC PROVIDERS
-- =============================================================================

CREATE TABLE oidc_providers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,       -- Internal name (keycloak, authentik, etc.)
    display_name VARCHAR(255) NOT NULL,       -- UI display name
    issuer_url VARCHAR(512) NOT NULL,         -- OIDC issuer URL
    client_id VARCHAR(255) NOT NULL,
    client_secret_encrypted BYTEA NOT NULL,   -- Encrypted with server key
    scopes TEXT[] NOT NULL DEFAULT ARRAY['openid', 'profile', 'email'],
    enabled BOOLEAN NOT NULL DEFAULT true,
    auto_create_users BOOLEAN NOT NULL DEFAULT true,
    default_admin BOOLEAN NOT NULL DEFAULT false,  -- New users are admins
    claim_mappings JSONB NOT NULL DEFAULT '{}',    -- Custom claim mappings
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_oidc_providers_enabled ON oidc_providers(enabled) WHERE enabled = true;

-- =============================================================================
-- OIDC USER LINKS
-- =============================================================================

CREATE TABLE oidc_user_links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider_id UUID NOT NULL REFERENCES oidc_providers(id) ON DELETE CASCADE,
    subject VARCHAR(255) NOT NULL,            -- OIDC 'sub' claim
    email VARCHAR(255),                       -- Email from OIDC (may differ from user email)
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_login_at TIMESTAMPTZ,
    UNIQUE(provider_id, subject)
);

CREATE INDEX idx_oidc_user_links_user_id ON oidc_user_links(user_id);
CREATE INDEX idx_oidc_user_links_provider_subject ON oidc_user_links(provider_id, subject);

-- =============================================================================
-- LIBRARIES
-- =============================================================================

CREATE TYPE library_type AS ENUM ('movies', 'tvshows', 'music', 'photos', 'mixed');

CREATE TABLE libraries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    type library_type NOT NULL,
    paths TEXT[] NOT NULL DEFAULT '{}',       -- Array of filesystem paths
    settings JSONB NOT NULL DEFAULT '{}',     -- Library-specific settings
    is_visible BOOLEAN NOT NULL DEFAULT true,
    scan_interval_hours INT DEFAULT 24,
    last_scan_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_libraries_type ON libraries(type);

-- =============================================================================
-- MEDIA ITEMS
-- =============================================================================

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

CREATE INDEX idx_media_items_library_id ON media_items(library_id);
CREATE INDEX idx_media_items_parent_id ON media_items(parent_id);
CREATE INDEX idx_media_items_type ON media_items(type);
CREATE INDEX idx_media_items_name ON media_items(name);
CREATE INDEX idx_media_items_path ON media_items(path) WHERE path IS NOT NULL;
CREATE INDEX idx_media_items_year ON media_items(year) WHERE year IS NOT NULL;

-- Full-text search index
ALTER TABLE media_items ADD COLUMN search_vector tsvector
    GENERATED ALWAYS AS (
        setweight(to_tsvector('english', coalesce(name, '')), 'A') ||
        setweight(to_tsvector('english', coalesce(overview, '')), 'B') ||
        setweight(to_tsvector('english', coalesce(tagline, '')), 'C')
    ) STORED;

CREATE INDEX idx_media_items_search ON media_items USING GIN(search_vector);

-- =============================================================================
-- IMAGES
-- =============================================================================

CREATE TYPE image_type AS ENUM (
    'primary', 'backdrop', 'logo', 'thumb', 'banner',
    'art', 'disc', 'box', 'screenshot', 'chapter'
);

CREATE TABLE images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    item_id UUID NOT NULL REFERENCES media_items(id) ON DELETE CASCADE,
    type image_type NOT NULL,
    index INT NOT NULL DEFAULT 0,             -- For multiple images of same type
    path TEXT NOT NULL,                       -- Filesystem path or URL
    width INT,
    height INT,
    blurhash VARCHAR(100),                    -- BlurHash for placeholder
    provider VARCHAR(50),                     -- Source provider (local, tmdb, tvdb, etc.)
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(item_id, type, index)
);

CREATE INDEX idx_images_item_id ON images(item_id);

-- =============================================================================
-- PLAYBACK PROGRESS
-- =============================================================================

CREATE TABLE playback_progress (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    item_id UUID NOT NULL REFERENCES media_items(id) ON DELETE CASCADE,
    position_ticks BIGINT NOT NULL DEFAULT 0, -- Playback position in ticks
    played BOOLEAN NOT NULL DEFAULT false,    -- Has been fully watched
    play_count INT NOT NULL DEFAULT 0,
    last_played_at TIMESTAMPTZ,
    audio_stream_index INT,                   -- Selected audio track
    subtitle_stream_index INT,                -- Selected subtitle track
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, item_id)
);

CREATE INDEX idx_playback_progress_user_id ON playback_progress(user_id);
CREATE INDEX idx_playback_progress_item_id ON playback_progress(item_id);
CREATE INDEX idx_playback_progress_last_played ON playback_progress(last_played_at DESC);

-- =============================================================================
-- PEOPLE (Actors, Directors, etc.)
-- =============================================================================

CREATE TYPE person_type AS ENUM ('actor', 'director', 'writer', 'producer', 'composer', 'guest_star');

CREATE TABLE people (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    sort_name VARCHAR(255),
    overview TEXT,
    birth_date DATE,
    death_date DATE,
    birth_place VARCHAR(255),
    provider_ids JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_people_name ON people(name);

CREATE TABLE media_people (
    item_id UUID NOT NULL REFERENCES media_items(id) ON DELETE CASCADE,
    person_id UUID NOT NULL REFERENCES people(id) ON DELETE CASCADE,
    type person_type NOT NULL,
    role VARCHAR(255),                        -- Character name for actors
    sort_order INT NOT NULL DEFAULT 0,
    PRIMARY KEY (item_id, person_id, type)
);

CREATE INDEX idx_media_people_item_id ON media_people(item_id);
CREATE INDEX idx_media_people_person_id ON media_people(person_id);

-- =============================================================================
-- ACTIVITY LOG
-- =============================================================================

CREATE TYPE activity_type AS ENUM (
    'user_login', 'user_logout', 'user_created', 'user_deleted',
    'playback_start', 'playback_stop', 'playback_progress',
    'library_scan_start', 'library_scan_complete',
    'item_added', 'item_removed', 'item_updated',
    'system_start', 'system_stop', 'system_update'
);

CREATE TABLE activity_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    type activity_type NOT NULL,
    item_id UUID REFERENCES media_items(id) ON DELETE SET NULL,
    severity VARCHAR(20) NOT NULL DEFAULT 'info',  -- info, warning, error
    overview TEXT,
    short_overview VARCHAR(500),
    data JSONB NOT NULL DEFAULT '{}',
    ip_address INET,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_activity_log_user_id ON activity_log(user_id);
CREATE INDEX idx_activity_log_type ON activity_log(type);
CREATE INDEX idx_activity_log_created_at ON activity_log(created_at DESC);

-- Partition by month for better performance (optional, can be enabled later)
-- CREATE TABLE activity_log (...) PARTITION BY RANGE (created_at);

-- =============================================================================
-- TRIGGERS FOR updated_at
-- =============================================================================

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_oidc_providers_updated_at
    BEFORE UPDATE ON oidc_providers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_libraries_updated_at
    BEFORE UPDATE ON libraries
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_media_items_updated_at
    BEFORE UPDATE ON media_items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_playback_progress_updated_at
    BEFORE UPDATE ON playback_progress
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_people_updated_at
    BEFORE UPDATE ON people
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
