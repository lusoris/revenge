-- Create user_avatars table in shared schema
-- Avatar storage with versioning and metadata

CREATE TABLE IF NOT EXISTS shared.user_avatars (
    -- Primary key
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Owner
    user_id UUID NOT NULL REFERENCES shared.users(id) ON DELETE CASCADE,

    -- Storage
    file_path TEXT NOT NULL, -- Path to avatar file in storage
    file_size_bytes BIGINT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,

    -- Image metadata
    width INTEGER NOT NULL,
    height INTEGER NOT NULL,
    is_animated BOOLEAN DEFAULT FALSE,

    -- Versioning
    version INTEGER NOT NULL DEFAULT 1,
    is_current BOOLEAN DEFAULT TRUE,

    -- Upload metadata
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    uploaded_from_ip INET,
    uploaded_from_user_agent TEXT,

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Partial unique index to ensure only one current avatar per user
CREATE UNIQUE INDEX idx_user_avatars_unique_current 
    ON shared.user_avatars(user_id) 
    WHERE is_current = TRUE AND deleted_at IS NULL;

-- Indexes
CREATE INDEX idx_user_avatars_user_id ON shared.user_avatars(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_user_avatars_current ON shared.user_avatars(user_id, is_current) WHERE deleted_at IS NULL;
CREATE INDEX idx_user_avatars_version ON shared.user_avatars(user_id, version DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_user_avatars_uploaded_at ON shared.user_avatars(uploaded_at DESC);

-- Comments
COMMENT ON TABLE shared.user_avatars IS 'User avatar storage with versioning and metadata';
COMMENT ON COLUMN shared.user_avatars.file_path IS 'Storage path relative to upload directory';
COMMENT ON COLUMN shared.user_avatars.version IS 'Avatar version number (incremented on each upload)';
COMMENT ON COLUMN shared.user_avatars.is_current IS 'Whether this is the currently active avatar';
COMMENT ON INDEX shared.idx_user_avatars_unique_current IS 'Ensures only one current avatar per user';
