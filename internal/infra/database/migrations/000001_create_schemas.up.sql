-- Create database schemas
-- public: Main content (movies, TV, music, etc.)
-- shared: Shared services (users, sessions, settings, etc.)
-- qar: Adult content (isolated for privacy)

-- public schema already exists by default

-- Create shared schema for shared services
CREATE SCHEMA IF NOT EXISTS shared;

-- Create qar schema for adult content (isolated)
CREATE SCHEMA IF NOT EXISTS qar;

-- Grant usage permissions
GRANT USAGE ON SCHEMA public TO PUBLIC;
GRANT USAGE ON SCHEMA shared TO PUBLIC;
GRANT USAGE ON SCHEMA qar TO PUBLIC;

-- Set search path to include all schemas
-- This allows unqualified table names to be resolved
ALTER DATABASE CURRENT_DATABASE() SET search_path TO public, shared, qar;

COMMENT ON SCHEMA public IS 'Main content: movies, TV shows, music, audiobooks, books, podcasts, photos, comics, live TV';
COMMENT ON SCHEMA shared IS 'Shared services: users, sessions, settings, API keys, RBAC, activity logs';
COMMENT ON SCHEMA qar IS 'Adult content (QAR): voyages, expeditions, treasures, crew, ports, flags (isolated for privacy)';
