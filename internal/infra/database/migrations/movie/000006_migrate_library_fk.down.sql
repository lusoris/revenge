-- Rollback: Restore old library_id column
BEGIN;

-- Re-add the old library_id column
ALTER TABLE movies ADD COLUMN library_id UUID REFERENCES libraries(id) ON DELETE CASCADE;

-- Migrate data back from movie_library_id to library_id
UPDATE movies
SET library_id = movie_library_id
WHERE movie_library_id IS NOT NULL
AND EXISTS (SELECT 1 FROM libraries l WHERE l.id = movies.movie_library_id);

-- Make library_id NOT NULL
ALTER TABLE movies ALTER COLUMN library_id SET NOT NULL;

-- Re-create old index
CREATE INDEX idx_movies_library_old ON movies(library_id);

-- Make movie_library_id nullable again
ALTER TABLE movies ALTER COLUMN movie_library_id DROP NOT NULL;

COMMIT;
