-- Create server_settings table in shared schema
-- Key-value store for server-wide configuration

CREATE TABLE IF NOT EXISTS shared.server_settings (
    -- Setting key (primary key)
    key VARCHAR(255) PRIMARY KEY,

    -- Setting value (JSON for flexibility)
    value JSONB NOT NULL,

    -- Metadata
    description TEXT,
    category VARCHAR(100), -- e.g., 'auth', 'security', 'features', 'ui'

    -- Validation
    data_type VARCHAR(50) NOT NULL, -- e.g., 'string', 'number', 'boolean', 'json'
    is_secret BOOLEAN DEFAULT FALSE, -- Mark sensitive values (encrypted)
    is_public BOOLEAN DEFAULT FALSE, -- Exposed to public API

    -- Constraints
    allowed_values JSONB, -- Array of allowed values (for enums)
    min_value NUMERIC,
    max_value NUMERIC,
    pattern TEXT, -- Regex pattern for validation

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID REFERENCES shared.users(id) ON DELETE SET NULL
);

-- Indexes
CREATE INDEX idx_server_settings_category ON shared.server_settings(category);
CREATE INDEX idx_server_settings_is_public ON shared.server_settings(is_public) WHERE is_public = TRUE;
CREATE INDEX idx_server_settings_updated_at ON shared.server_settings(updated_at DESC);

-- Comments
COMMENT ON TABLE shared.server_settings IS 'Server-wide configuration settings with validation';
COMMENT ON COLUMN shared.server_settings.value IS 'Setting value stored as JSONB for type flexibility';
COMMENT ON COLUMN shared.server_settings.is_secret IS 'Indicates sensitive values that should be encrypted';
COMMENT ON COLUMN shared.server_settings.is_public IS 'Settings exposed via public API (e.g., server name)';
COMMENT ON COLUMN shared.server_settings.data_type IS 'Expected data type for validation';

-- Insert default settings
INSERT INTO shared.server_settings (key, value, description, category, data_type, is_public) VALUES
    ('server.name', '"Revenge Media Server"', 'Server display name', 'general', 'string', TRUE),
    ('auth.jwt.access_token_expiry', '3600', 'JWT access token expiry in seconds (1 hour)', 'auth', 'number', FALSE),
    ('auth.jwt.refresh_token_expiry', '2592000', 'JWT refresh token expiry in seconds (30 days)', 'auth', 'number', FALSE),
    ('auth.password.min_length', '8', 'Minimum password length', 'auth', 'number', FALSE),
    ('auth.session.max_per_user', '10', 'Maximum active sessions per user', 'auth', 'number', FALSE),
    ('features.registration_enabled', 'true', 'Allow new user registration', 'features', 'boolean', TRUE),
    ('features.oidc_enabled', 'false', 'Enable OIDC authentication', 'features', 'boolean', TRUE)
ON CONFLICT (key) DO NOTHING;
