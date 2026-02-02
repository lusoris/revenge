-- Create user_preferences table in shared schema
-- User notification and system preferences

CREATE TABLE IF NOT EXISTS shared.user_preferences (
    -- Primary key
    user_id UUID PRIMARY KEY REFERENCES shared.users(id) ON DELETE CASCADE,

    -- Notification preferences (JSONB for flexibility)
    email_notifications JSONB DEFAULT '{"enabled": true, "frequency": "instant"}'::jsonb,
    push_notifications JSONB DEFAULT '{"enabled": false}'::jsonb,
    digest_notifications JSONB DEFAULT '{"enabled": true, "frequency": "weekly"}'::jsonb,

    -- Privacy preferences
    profile_visibility VARCHAR(20) DEFAULT 'private' CHECK (profile_visibility IN ('public', 'friends', 'private')),
    show_email BOOLEAN DEFAULT FALSE,
    show_activity BOOLEAN DEFAULT TRUE,

    -- Display preferences
    theme VARCHAR(20) DEFAULT 'system' CHECK (theme IN ('light', 'dark', 'system')),
    display_language VARCHAR(10) DEFAULT 'en-US',
    content_language VARCHAR(255), -- Preferred content languages (comma-separated)

    -- Content preferences
    show_adult_content BOOLEAN DEFAULT FALSE,
    show_spoilers BOOLEAN DEFAULT FALSE,
    auto_play_videos BOOLEAN DEFAULT TRUE,

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_user_preferences_updated_at ON shared.user_preferences(updated_at DESC);

-- Comments
COMMENT ON TABLE shared.user_preferences IS 'User preferences for notifications, privacy, and display settings';
COMMENT ON COLUMN shared.user_preferences.email_notifications IS 'Email notification settings (enabled, frequency, types)';
COMMENT ON COLUMN shared.user_preferences.push_notifications IS 'Push notification settings (enabled, device tokens)';
COMMENT ON COLUMN shared.user_preferences.digest_notifications IS 'Digest notification settings (enabled, frequency)';
COMMENT ON COLUMN shared.user_preferences.profile_visibility IS 'Who can view the user profile (public, friends, private)';
