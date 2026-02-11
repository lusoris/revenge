-- Migration 000041: Make tmdb_id nullable on networks table
-- This allows creating networks from non-TMDb providers that don't have a TMDb ID.

ALTER TABLE tvshow.networks ALTER COLUMN tmdb_id DROP NOT NULL;
