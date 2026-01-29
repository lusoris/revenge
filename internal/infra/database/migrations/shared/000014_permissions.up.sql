-- Permissions table and role-permission mappings
-- Based on ARCHITECTURE_V2.md RBAC definitions

-- Permission categories
CREATE TYPE permission_category AS ENUM (
    'system',      -- Server settings, system administration
    'users',       -- User management
    'libraries',   -- Library management
    'content',     -- Content browsing and metadata
    'playback',    -- Media playback
    'social',      -- Ratings, playlists, collections
    'adult'        -- Adult content access (schema c)
);

-- Permissions table
CREATE TABLE permissions (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100) UNIQUE NOT NULL,
    description TEXT NOT NULL,
    category    permission_category NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Role-permission mapping table
CREATE TABLE role_permissions (
    role        user_role NOT NULL,
    permission_id INT NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role, permission_id)
);

-- Seed permissions
INSERT INTO permissions (name, description, category) VALUES
    -- System permissions
    ('system.settings.read', 'View server settings', 'system'),
    ('system.settings.write', 'Modify server settings', 'system'),
    ('system.logs.read', 'View activity logs', 'system'),
    ('system.jobs.read', 'View background jobs', 'system'),
    ('system.jobs.manage', 'Manage background jobs (cancel, retry)', 'system'),
    ('system.apikeys.manage', 'Manage API keys', 'system'),

    -- User permissions
    ('users.read', 'View user list and profiles', 'users'),
    ('users.create', 'Create new users', 'users'),
    ('users.update', 'Update user accounts', 'users'),
    ('users.delete', 'Delete user accounts', 'users'),
    ('users.sessions.manage', 'Manage user sessions (force logout)', 'users'),

    -- Library permissions
    ('libraries.read', 'View library list', 'libraries'),
    ('libraries.create', 'Create new libraries', 'libraries'),
    ('libraries.update', 'Update library settings', 'libraries'),
    ('libraries.delete', 'Delete libraries', 'libraries'),
    ('libraries.scan', 'Trigger library scans', 'libraries'),

    -- Content permissions
    ('content.browse', 'Browse content (movies, shows, music, etc.)', 'content'),
    ('content.metadata.read', 'View content metadata', 'content'),
    ('content.metadata.write', 'Edit content metadata', 'content'),
    ('content.images.manage', 'Manage content images', 'content'),
    ('content.delete', 'Delete content items', 'content'),

    -- Playback permissions
    ('playback.stream', 'Stream media content', 'playback'),
    ('playback.download', 'Download media files', 'playback'),
    ('playback.transcode', 'Request transcoded streams', 'playback'),

    -- Social permissions
    ('social.rate', 'Rate content', 'social'),
    ('social.playlists.create', 'Create playlists', 'social'),
    ('social.playlists.manage', 'Manage own playlists', 'social'),
    ('social.collections.create', 'Create collections', 'social'),
    ('social.collections.manage', 'Manage own collections', 'social'),
    ('social.history.read', 'View own watch/play history', 'social'),
    ('social.favorites.manage', 'Manage favorites', 'social'),

    -- Adult content permissions
    ('adult.browse', 'Browse adult content (schema c)', 'adult'),
    ('adult.stream', 'Stream adult content', 'adult'),
    ('adult.metadata.write', 'Edit adult content metadata', 'adult');

-- Assign permissions to roles
-- Admin: Full access (all permissions)
INSERT INTO role_permissions (role, permission_id)
SELECT 'admin', id FROM permissions;

-- Moderator: Manage libraries, metadata, moderate content
INSERT INTO role_permissions (role, permission_id)
SELECT 'moderator', id FROM permissions
WHERE name IN (
    'system.logs.read',
    'system.jobs.read',
    'users.read',
    'libraries.read',
    'libraries.create',
    'libraries.update',
    'libraries.scan',
    'content.browse',
    'content.metadata.read',
    'content.metadata.write',
    'content.images.manage',
    'playback.stream',
    'playback.download',
    'playback.transcode',
    'social.rate',
    'social.playlists.create',
    'social.playlists.manage',
    'social.collections.create',
    'social.collections.manage',
    'social.history.read',
    'social.favorites.manage'
);

-- User: Browse, play, rate, create playlists
INSERT INTO role_permissions (role, permission_id)
SELECT 'user', id FROM permissions
WHERE name IN (
    'libraries.read',
    'content.browse',
    'content.metadata.read',
    'playback.stream',
    'playback.download',
    'playback.transcode',
    'social.rate',
    'social.playlists.create',
    'social.playlists.manage',
    'social.collections.create',
    'social.collections.manage',
    'social.history.read',
    'social.favorites.manage'
);

-- Guest: Browse only (no playback by default - configurable)
INSERT INTO role_permissions (role, permission_id)
SELECT 'guest', id FROM permissions
WHERE name IN (
    'libraries.read',
    'content.browse',
    'content.metadata.read'
);

-- Indexes
CREATE INDEX idx_role_permissions_role ON role_permissions(role);
CREATE INDEX idx_permissions_category ON permissions(category);
