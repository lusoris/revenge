-- Enable required PostgreSQL extensions
CREATE EXTENSION IF NOT EXISTS pgcrypto;      -- UUID generation, encryption
CREATE EXTENSION IF NOT EXISTS pg_trgm;       -- Trigram similarity for fuzzy search
CREATE EXTENSION IF NOT EXISTS unaccent;      -- Remove accents for search normalization
