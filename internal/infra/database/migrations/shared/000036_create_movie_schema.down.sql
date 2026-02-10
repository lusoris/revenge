-- Move movie tables back to public schema
ALTER TABLE movie.movie_watched SET SCHEMA public;

ALTER TABLE movie.movie_genres SET SCHEMA public;

ALTER TABLE movie.movie_collection_members SET SCHEMA public;

ALTER TABLE movie.movie_collections SET SCHEMA public;

ALTER TABLE movie.movie_credits SET SCHEMA public;

ALTER TABLE movie.movie_files SET SCHEMA public;

ALTER TABLE movie.movies SET SCHEMA public;

-- Drop movie schema
DROP SCHEMA IF EXISTS movie;
