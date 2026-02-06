-- Add moderator role and permissions
-- This migration extends the RBAC system with moderator permissions

-- Add moderator role permissions
INSERT INTO shared.casbin_rule (ptype, v0, v1, v2) VALUES
    -- Moderator can read users
    ('p', 'moderator', 'users', 'read'),
    -- Moderator can update users (but not delete)
    ('p', 'moderator', 'users', 'write'),
    -- Moderator has full access to requests
    ('p', 'moderator', 'requests', 'read'),
    ('p', 'moderator', 'requests', 'write'),
    ('p', 'moderator', 'requests', 'delete'),
    -- Moderator can moderate content
    ('p', 'moderator', 'movies', 'read'),
    ('p', 'moderator', 'movies', 'write'),
    -- Moderator can view audit logs
    ('p', 'moderator', 'audit', 'read'),
    -- Moderator can read libraries
    ('p', 'moderator', 'library', 'read'),
    -- Moderator profile access
    ('p', 'moderator', 'profile', 'read'),
    ('p', 'moderator', 'profile', 'write')
ON CONFLICT DO NOTHING;

-- Extend user role with requests permission
INSERT INTO shared.casbin_rule (ptype, v0, v1, v2) VALUES
    -- User can create and manage their own requests
    ('p', 'user', 'requests', 'read'),
    ('p', 'user', 'requests', 'write'),
    -- User can read movies
    ('p', 'user', 'movies', 'read')
ON CONFLICT DO NOTHING;

COMMENT ON TABLE shared.casbin_rule IS 'RBAC policies with moderator role added';
