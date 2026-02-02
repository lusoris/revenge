-- Create user_settings table in shared schema
-- User-specific configuration settings

CREATE TABLE IF NOT EXISTS shared.user_settings (
    -- Composite primary key
    user_id UUID NOT NULL REFERENCES shared.users(id) ON DELETE CASCADE,
    key VARCHAR(255) NOT NULL,
    
    -- Setting value (JSON for flexibility)
    value JSONB NOT NULL,

    -- Metadata
    description TEXT,
    category VARCHAR(100), -- e.g., 'ui', 'notifications', 'privacy', 'playback'

    -- Validation
    data_type VARCHAR(50) NOT NULL, -- e.g., 'string', 'number', 'boolean', 'json'

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (user_id, key)
);

-- Indexes
CREATE INDEX idx_user_settings_user_id ON shared.user_settings(user_id);
CREATE INDEX idx_user_settings_category ON shared.user_settings(user_id, category);
CREATE INDEX idx_user_settings_updated_at ON shared.user_settings(updated_at DESC);

-- Comments
COMMENT ON TABLE shared.user_settings IS 'User-specific configuration settings';
COMMENT ON COLUMN shared.user_settings.value IS 'Setting value stored as JSONB for type flexibility';
COMMENT ON COLUMN shared.user_settings.data_type IS 'Expected data type for validation';
COMMENT ON COLUMN shared.user_settings.category IS 'Setting category for organization';
