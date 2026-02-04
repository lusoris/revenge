-- Remove moderator role permissions
DELETE FROM shared.casbin_rule WHERE ptype = 'p' AND v0 = 'moderator';

-- Remove extended user permissions
DELETE FROM shared.casbin_rule WHERE ptype = 'p' AND v0 = 'user' AND v1 = 'requests';
DELETE FROM shared.casbin_rule WHERE ptype = 'p' AND v0 = 'user' AND v1 = 'movies';
