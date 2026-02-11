-- Rollback: Make tmdb_id required again on networks table
-- Note: This will fail if any rows have NULL tmdb_id values.

ALTER TABLE tvshow.networks ALTER COLUMN tmdb_id SET NOT NULL;
