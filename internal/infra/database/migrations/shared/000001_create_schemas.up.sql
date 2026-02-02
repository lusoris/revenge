-- Create database schemas
-- public: Main content (movies, TV, music, etc.)
-- shared: Shared services (users, sessions, settings, etc.)
-- qar: Adult content (isolated for privacy)

-- public schema already exists by default
COMMENT ON SCHEMA public IS 'Main content: movies, TV shows, music, audiobooks, books, podcasts';

-- Create shared schema for shared services
CREATE SCHEMA IF NOT EXISTS shared;
COMMENT ON SCHEMA shared IS 'Shared services: users, sessions, settings, RBAC, activity';

-- Create qar schema for adult content (isolated)
CREATE SCHEMA IF NOT EXISTS qar;
COMMENT ON SCHEMA qar IS 'QAR (Adult content): requires legacy:read scope';

-- Grant usage permissions
GRANT USAGE ON SCHEMA public TO PUBLIC;
GRANT USAGE ON SCHEMA shared TO PUBLIC;
GRANT USAGE ON SCHEMA qar TO PUBLIC;
