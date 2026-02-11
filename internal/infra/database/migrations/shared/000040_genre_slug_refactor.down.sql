-- Rollback: Restore tmdb_genre_id columns in genre tables

-- ============================================================
-- Movie Genres
-- ============================================================
ALTER TABLE movie.movie_genres
DROP CONSTRAINT IF EXISTS uq_movie_genre;

DROP INDEX IF EXISTS movie.idx_movie_genres_slug;

ALTER TABLE movie.movie_genres
ADD COLUMN tmdb_genre_id INTEGER NOT NULL DEFAULT 0;

ALTER TABLE movie.movie_genres DROP COLUMN slug;

ALTER TABLE movie.movie_genres
ADD CONSTRAINT uq_movie_genre UNIQUE (movie_id, tmdb_genre_id);

CREATE INDEX idx_movie_genres_genre_id ON movie.movie_genres (tmdb_genre_id);

-- ============================================================
-- Series Genres
-- ============================================================
ALTER TABLE tvshow.series_genres
DROP CONSTRAINT IF EXISTS uq_series_genre;

DROP INDEX IF EXISTS tvshow.idx_series_genres_slug;

ALTER TABLE tvshow.series_genres
ADD COLUMN tmdb_genre_id INTEGER NOT NULL DEFAULT 0;

ALTER TABLE tvshow.series_genres DROP COLUMN slug;

ALTER TABLE tvshow.series_genres
ADD CONSTRAINT series_genres_series_id_tmdb_genre_id_key UNIQUE (series_id, tmdb_genre_id);

CREATE INDEX idx_series_genres_genre ON tvshow.series_genres (tmdb_genre_id);
