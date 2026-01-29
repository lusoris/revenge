-- Playlists - user-created ordered lists of content
-- Supports both video and audio content types

CREATE TYPE playlist_type AS ENUM (
    'video',      -- Movies, episodes, clips
    'audio',      -- Songs, podcast episodes, audiobook chapters
    'mixed'       -- Allow both types (future use)
);

CREATE TYPE playlist_visibility AS ENUM (
    'private',    -- Only owner can see
    'shared',     -- Owner can share link
    'public'      -- Visible to all users
);

-- Main playlists table
CREATE TABLE playlists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Basic info
    name VARCHAR(255) NOT NULL,
    description TEXT,
    playlist_type playlist_type NOT NULL,
    visibility playlist_visibility NOT NULL DEFAULT 'private',

    -- Display
    thumbnail_url TEXT,
    sort_order INT NOT NULL DEFAULT 0,

    -- Playback settings
    shuffle_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    repeat_mode VARCHAR(20) NOT NULL DEFAULT 'none', -- none, one, all

    -- Metadata
    item_count INT NOT NULL DEFAULT 0,
    total_duration_ms BIGINT NOT NULL DEFAULT 0,

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_played_at TIMESTAMPTZ
);

-- Playlist items - the actual content in playlists
CREATE TABLE playlist_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    playlist_id UUID NOT NULL REFERENCES playlists(id) ON DELETE CASCADE,

    -- Content reference (polymorphic)
    content_type VARCHAR(50) NOT NULL, -- movie, episode, track, podcast_episode, etc.
    content_id UUID NOT NULL,

    -- Ordering
    position INT NOT NULL,

    -- Item-specific settings
    start_time_ms BIGINT, -- Start playback at specific time
    end_time_ms BIGINT,   -- End playback at specific time

    -- Metadata snapshot (for display even if content is deleted)
    title VARCHAR(500),
    duration_ms BIGINT,
    thumbnail_url TEXT,

    -- Timestamps
    added_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_played_at TIMESTAMPTZ,

    UNIQUE (playlist_id, position)
);

-- Playlist collaborators - users who can edit shared playlists
CREATE TABLE playlist_collaborators (
    playlist_id UUID NOT NULL REFERENCES playlists(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    can_edit BOOLEAN NOT NULL DEFAULT FALSE,
    added_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (playlist_id, user_id)
);

-- Indexes
CREATE INDEX idx_playlists_user_id ON playlists(user_id);
CREATE INDEX idx_playlists_type ON playlists(playlist_type);
CREATE INDEX idx_playlists_visibility ON playlists(visibility);
CREATE INDEX idx_playlist_items_playlist_id ON playlist_items(playlist_id);
CREATE INDEX idx_playlist_items_content ON playlist_items(content_type, content_id);
CREATE INDEX idx_playlist_collaborators_user ON playlist_collaborators(user_id);

-- Trigger to update playlist metadata on item changes
CREATE OR REPLACE FUNCTION update_playlist_metadata()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE playlists
    SET
        item_count = (SELECT COUNT(*) FROM playlist_items WHERE playlist_id = COALESCE(NEW.playlist_id, OLD.playlist_id)),
        total_duration_ms = (SELECT COALESCE(SUM(duration_ms), 0) FROM playlist_items WHERE playlist_id = COALESCE(NEW.playlist_id, OLD.playlist_id)),
        updated_at = NOW()
    WHERE id = COALESCE(NEW.playlist_id, OLD.playlist_id);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_playlist_items_update_metadata
AFTER INSERT OR UPDATE OR DELETE ON playlist_items
FOR EACH ROW EXECUTE FUNCTION update_playlist_metadata();
