-- Dynamic RBAC: Replace static role ENUM with dynamic roles using Casbin
-- This migration adds support for admin-defined custom roles with granular permissions
BEGIN;

-- Create roles table for custom role management
-- Note: Casbin handles the actual permission policies in casbin_rule table
CREATE TABLE roles (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(100) NOT NULL UNIQUE,
    display_name    VARCHAR(255) NOT NULL,
    description     TEXT,
    color           VARCHAR(7),                          -- Hex color for UI (e.g., #FF5733)
    icon            VARCHAR(50),                         -- Icon name for UI
    is_system       BOOLEAN NOT NULL DEFAULT false,      -- System roles cannot be deleted
    is_default      BOOLEAN NOT NULL DEFAULT false,      -- Default role for new users
    priority        INT NOT NULL DEFAULT 0,              -- Higher = more important (for UI sorting)
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      UUID REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX idx_roles_name ON roles(name);
CREATE INDEX idx_roles_is_default ON roles(is_default) WHERE is_default = true;
CREATE INDEX idx_roles_priority ON roles(priority DESC);

CREATE TRIGGER roles_updated_at
    BEFORE UPDATE ON roles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- Add role_id column to users (FK to roles table)
ALTER TABLE users ADD COLUMN role_id UUID REFERENCES roles(id) ON DELETE SET NULL;

-- Insert default system roles (matching previous ENUM values)
INSERT INTO roles (name, display_name, description, is_system, priority, color) VALUES
    ('admin', 'Administrator', 'Full access to all features and settings', true, 1000, '#EF4444'),
    ('moderator', 'Moderator', 'Can manage libraries, metadata, and moderate content', true, 500, '#F59E0B'),
    ('user', 'User', 'Standard user with browse, play, and social features', true, 100, '#3B82F6'),
    ('guest', 'Guest', 'Limited access for browsing only', true, 0, '#6B7280');

-- Set the 'user' role as default for new users
UPDATE roles SET is_default = true WHERE name = 'user';

-- Migrate existing users from ENUM to role_id
UPDATE users u SET role_id = r.id
FROM roles r
WHERE u.role::text = r.name;

-- Create permission definitions table (for UI to show available permissions)
-- Note: This is for UI/API reference - actual enforcement is via Casbin policies
CREATE TABLE permission_definitions (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(100) NOT NULL UNIQUE,
    display_name    VARCHAR(255) NOT NULL,
    description     TEXT NOT NULL,
    category        VARCHAR(50) NOT NULL,
    is_dangerous    BOOLEAN NOT NULL DEFAULT false,      -- Marks permissions that need extra confirmation
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_permission_definitions_category ON permission_definitions(category);

-- Seed permission definitions (same as before, but with display names)
INSERT INTO permission_definitions (name, display_name, description, category, is_dangerous) VALUES
    -- System permissions
    ('system.settings.read', 'View Server Settings', 'View server configuration', 'System', false),
    ('system.settings.write', 'Modify Server Settings', 'Change server configuration', 'System', true),
    ('system.logs.read', 'View Activity Logs', 'View system activity logs', 'System', false),
    ('system.jobs.read', 'View Background Jobs', 'View background job status', 'System', false),
    ('system.jobs.manage', 'Manage Background Jobs', 'Cancel or retry background jobs', 'System', true),
    ('system.apikeys.manage', 'Manage API Keys', 'Create and revoke API keys', 'System', true),
    ('system.roles.read', 'View Roles', 'View role definitions', 'System', false),
    ('system.roles.manage', 'Manage Roles', 'Create, edit, and delete roles', 'System', true),

    -- User permissions
    ('users.read', 'View Users', 'View user list and profiles', 'Users', false),
    ('users.create', 'Create Users', 'Create new user accounts', 'Users', true),
    ('users.update', 'Update Users', 'Modify user accounts', 'Users', true),
    ('users.delete', 'Delete Users', 'Delete user accounts', 'Users', true),
    ('users.sessions.manage', 'Manage Sessions', 'Force logout users', 'Users', true),

    -- Library permissions
    ('libraries.read', 'View Libraries', 'View library list', 'Libraries', false),
    ('libraries.create', 'Create Libraries', 'Create new libraries', 'Libraries', true),
    ('libraries.update', 'Update Libraries', 'Modify library settings', 'Libraries', true),
    ('libraries.delete', 'Delete Libraries', 'Delete libraries', 'Libraries', true),
    ('libraries.scan', 'Scan Libraries', 'Trigger library scans', 'Libraries', false),

    -- Content permissions
    ('content.browse', 'Browse Content', 'Browse movies, shows, music, etc.', 'Content', false),
    ('content.metadata.read', 'View Metadata', 'View content metadata', 'Content', false),
    ('content.metadata.write', 'Edit Metadata', 'Edit content metadata', 'Content', false),
    ('content.images.manage', 'Manage Images', 'Upload and manage content images', 'Content', false),
    ('content.delete', 'Delete Content', 'Delete content items', 'Content', true),

    -- Playback permissions
    ('playback.stream', 'Stream Content', 'Stream media content', 'Playback', false),
    ('playback.download', 'Download Content', 'Download media files', 'Playback', false),
    ('playback.transcode', 'Request Transcoding', 'Request transcoded streams', 'Playback', false),

    -- Social permissions
    ('social.rate', 'Rate Content', 'Rate movies, shows, albums, etc.', 'Social', false),
    ('social.playlists.create', 'Create Playlists', 'Create new playlists', 'Social', false),
    ('social.playlists.manage', 'Manage Playlists', 'Edit and delete own playlists', 'Social', false),
    ('social.collections.create', 'Create Collections', 'Create new collections', 'Social', false),
    ('social.collections.manage', 'Manage Collections', 'Edit and delete own collections', 'Social', false),
    ('social.history.read', 'View History', 'View own watch/play history', 'Social', false),
    ('social.favorites.manage', 'Manage Favorites', 'Add and remove favorites', 'Social', false),

    -- Adult content permissions
    ('adult.browse', 'Browse Adult Content', 'Access adult content library', 'Adult', false),
    ('adult.stream', 'Stream Adult Content', 'Stream adult media', 'Adult', false),
    ('adult.metadata.write', 'Edit Adult Metadata', 'Edit adult content metadata', 'Adult', false);

-- Casbin will create its own casbin_rule table via the adapter
-- We'll seed the default policies in the application code after Casbin initializes

-- Drop old static tables (keep for rollback reference, but not used)
-- Note: We keep role_permissions and permissions tables for migration rollback
-- They will be dropped in a future cleanup migration once Casbin is stable

COMMIT;
