-- Revert fine-grained permissions back to coarse-grained

-- Delete fine-grained permissions
DELETE FROM shared.casbin_rule WHERE ptype = 'p' AND v0 IN ('user', 'guest', 'moderator');

-- Restore coarse-grained moderator permissions
INSERT INTO shared.casbin_rule (ptype, v0, v1, v2) VALUES
    ('p', 'moderator', 'users', 'read'),
    ('p', 'moderator', 'users', 'write'),
    ('p', 'moderator', 'requests', 'read'),
    ('p', 'moderator', 'requests', 'write'),
    ('p', 'moderator', 'requests', 'delete'),
    ('p', 'moderator', 'movies', 'read'),
    ('p', 'moderator', 'movies', 'write'),
    ('p', 'moderator', 'audit', 'read'),
    ('p', 'moderator', 'library', 'read'),
    ('p', 'moderator', 'profile', 'read'),
    ('p', 'moderator', 'profile', 'write')
ON CONFLICT DO NOTHING;

-- Restore coarse-grained user permissions
INSERT INTO shared.casbin_rule (ptype, v0, v1, v2) VALUES
    ('p', 'user', 'profile', 'read'),
    ('p', 'user', 'profile', 'write'),
    ('p', 'user', 'library', 'read'),
    ('p', 'user', 'playback', 'read'),
    ('p', 'user', 'playback', 'write'),
    ('p', 'user', 'requests', 'read'),
    ('p', 'user', 'requests', 'write'),
    ('p', 'user', 'movies', 'read')
ON CONFLICT DO NOTHING;

-- Restore coarse-grained guest permissions
INSERT INTO shared.casbin_rule (ptype, v0, v1, v2) VALUES
    ('p', 'guest', 'library', 'read')
ON CONFLICT DO NOTHING;

COMMENT ON TABLE shared.casbin_rule IS 'RBAC policies reverted to coarse-grained';
