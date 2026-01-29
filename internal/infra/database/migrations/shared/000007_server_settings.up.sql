-- 000007_server_settings.up.sql
-- Persisted server configuration and system settings

BEGIN;

-- Server settings table - key-value store for dynamic configuration
CREATE TABLE server_settings (
    key VARCHAR(255) PRIMARY KEY,
    value JSONB NOT NULL,
    category VARCHAR(50) NOT NULL DEFAULT 'general', -- general, security, media, etc.
    description TEXT,
    is_public BOOLEAN NOT NULL DEFAULT false,  -- If true, can be read by non-admin users
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_server_settings_category ON server_settings(category);
CREATE INDEX idx_server_settings_public ON server_settings(is_public) WHERE is_public = true;

-- Trigger for updated_at
CREATE TRIGGER update_server_settings_updated_at
    BEFORE UPDATE ON server_settings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Default settings
INSERT INTO server_settings (key, value, category, description, is_public) VALUES
    ('server.name', '"Revenge Media Server"', 'general', 'Server display name', true),
    ('server.version', '"1.0.0"', 'general', 'Server version', true),
    ('server.timezone', '"UTC"', 'general', 'Server timezone', false),
    ('security.require_authentication', 'true', 'security', 'Require authentication for all requests', false),
    ('security.allow_registration', 'false', 'security', 'Allow new user registration', false),
    ('media.default_transcoding_profile', '"720p"', 'media', 'Default transcoding profile', false),
    ('media.enable_hardware_acceleration', 'true', 'media', 'Enable hardware transcoding', false)
ON CONFLICT (key) DO NOTHING;

COMMIT;
