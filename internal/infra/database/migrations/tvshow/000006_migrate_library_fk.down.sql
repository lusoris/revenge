-- Rollback: Restore old library_id column
BEGIN;

-- Re-add the old library_id column
ALTER TABLE series ADD COLUMN library_id UUID REFERENCES libraries(id) ON DELETE CASCADE;

-- Migrate data back from tv_library_id to library_id
UPDATE series
SET library_id = tv_library_id
WHERE tv_library_id IS NOT NULL
AND EXISTS (SELECT 1 FROM libraries l WHERE l.id = series.tv_library_id);

-- Make library_id NOT NULL
ALTER TABLE series ALTER COLUMN library_id SET NOT NULL;

-- Re-create old index
CREATE INDEX idx_series_library_old ON series(library_id);

-- Make tv_library_id nullable again
ALTER TABLE series ALTER COLUMN tv_library_id DROP NOT NULL;

COMMIT;
