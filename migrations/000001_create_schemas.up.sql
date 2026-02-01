-- Create schemas for Revenge
-- public: Main content (movies, TV shows, music, etc.)
-- shared: Shared services (users, sessions, settings, etc.)
-- qar: Adult content (isolated with access control)

-- public schema already exists by default in PostgreSQL
-- Ensure it's configured correctly
COMMENT ON SCHEMA public IS 'Main content: movies, TV shows, music, audiobooks, books, podcasts';

-- Create shared schema for shared services
CREATE SCHEMA IF NOT EXISTS shared;
COMMENT ON SCHEMA shared IS 'Shared services: users, sessions, settings, RBAC, activity';

-- Create qar schema for adult content (requires legacy:read scope)
CREATE SCHEMA IF NOT EXISTS qar;
COMMENT ON SCHEMA qar IS 'QAR (Adult content): voyages, expeditions, treasures - requires legacy:read scope';

-- Set search path to include all schemas
-- Application will set this per-connection based on user permissions
-- Default search path: public, shared (qar requires explicit scope)
ALTER DATABASE revenge SET search_path TO public, shared;
