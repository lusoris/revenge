-- Update RBAC policies to fine-grained permissions
-- Changes coarse actions (read/write/delete) to fine-grained actions (list/get/create/update/delete/etc.)

-- Delete old coarse-grained permissions for roles (keep admin wildcard, user role assignments)
DELETE FROM shared.casbin_rule WHERE ptype = 'p' AND v0 IN ('user', 'guest', 'moderator');

-- Moderator role - comprehensive access without full admin
INSERT INTO shared.casbin_rule (ptype, v0, v1, v2) VALUES
    -- User management (limited - no create/delete)
    ('p', 'moderator', 'users', 'list'),
    ('p', 'moderator', 'users', 'get'),
    -- Own profile
    ('p', 'moderator', 'profile', 'read'),
    ('p', 'moderator', 'profile', 'update'),
    -- Movies - full access
    ('p', 'moderator', 'movies', 'list'),
    ('p', 'moderator', 'movies', 'get'),
    ('p', 'moderator', 'movies', 'create'),
    ('p', 'moderator', 'movies', 'update'),
    ('p', 'moderator', 'movies', 'delete'),
    -- Libraries - full access
    ('p', 'moderator', 'libraries', 'list'),
    ('p', 'moderator', 'libraries', 'get'),
    ('p', 'moderator', 'libraries', 'create'),
    ('p', 'moderator', 'libraries', 'update'),
    ('p', 'moderator', 'libraries', 'delete'),
    ('p', 'moderator', 'libraries', 'scan'),
    -- Playback
    ('p', 'moderator', 'playback', 'stream'),
    ('p', 'moderator', 'playback', 'progress'),
    -- Requests - full access
    ('p', 'moderator', 'requests', 'list'),
    ('p', 'moderator', 'requests', 'get'),
    ('p', 'moderator', 'requests', 'create'),
    ('p', 'moderator', 'requests', 'approve'),
    ('p', 'moderator', 'requests', 'delete'),
    -- Settings - read server, full user
    ('p', 'moderator', 'settings', 'read'),
    ('p', 'moderator', 'settings', 'user_read'),
    ('p', 'moderator', 'settings', 'user_write'),
    -- Audit - read only
    ('p', 'moderator', 'audit', 'read'),
    -- Integrations - full access
    ('p', 'moderator', 'integrations', 'list'),
    ('p', 'moderator', 'integrations', 'get'),
    ('p', 'moderator', 'integrations', 'create'),
    ('p', 'moderator', 'integrations', 'update'),
    ('p', 'moderator', 'integrations', 'delete'),
    ('p', 'moderator', 'integrations', 'sync'),
    -- Notifications - full access
    ('p', 'moderator', 'notifications', 'list'),
    ('p', 'moderator', 'notifications', 'get'),
    ('p', 'moderator', 'notifications', 'create'),
    ('p', 'moderator', 'notifications', 'update'),
    ('p', 'moderator', 'notifications', 'delete')
ON CONFLICT DO NOTHING;

-- User role - standard access
INSERT INTO shared.casbin_rule (ptype, v0, v1, v2) VALUES
    -- Own profile
    ('p', 'user', 'profile', 'read'),
    ('p', 'user', 'profile', 'update'),
    -- Movies - list and view only
    ('p', 'user', 'movies', 'list'),
    ('p', 'user', 'movies', 'get'),
    -- Libraries - list and view only
    ('p', 'user', 'libraries', 'list'),
    ('p', 'user', 'libraries', 'get'),
    -- Playback
    ('p', 'user', 'playback', 'stream'),
    ('p', 'user', 'playback', 'progress'),
    -- Requests - can create and view own
    ('p', 'user', 'requests', 'list'),
    ('p', 'user', 'requests', 'get'),
    ('p', 'user', 'requests', 'create'),
    -- Settings - own settings only
    ('p', 'user', 'settings', 'user_read'),
    ('p', 'user', 'settings', 'user_write'),
    -- Notifications - own only
    ('p', 'user', 'notifications', 'list'),
    ('p', 'user', 'notifications', 'get')
ON CONFLICT DO NOTHING;

-- Guest role - minimal read-only access
INSERT INTO shared.casbin_rule (ptype, v0, v1, v2) VALUES
    -- Profile read only
    ('p', 'guest', 'profile', 'read'),
    -- Movies - list and view only
    ('p', 'guest', 'movies', 'list'),
    ('p', 'guest', 'movies', 'get'),
    -- Libraries - list and view only
    ('p', 'guest', 'libraries', 'list'),
    ('p', 'guest', 'libraries', 'get'),
    -- Playback - stream only (no progress tracking)
    ('p', 'guest', 'playback', 'stream')
ON CONFLICT DO NOTHING;

COMMENT ON TABLE shared.casbin_rule IS 'RBAC policies with fine-grained permissions (v0.3.0)';
