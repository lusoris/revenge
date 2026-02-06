-- Drop movie_genres table
DROP INDEX IF EXISTS idx_movie_genres_name;
DROP INDEX IF EXISTS idx_movie_genres_genre_id;
DROP INDEX IF EXISTS idx_movie_genres_movie_id;
DROP TABLE IF EXISTS public.movie_genres;
