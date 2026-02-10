-- Create movie schema and move movie tables from public
-- Aligns with tvshow module pattern: each content module owns its schema

CREATE SCHEMA IF NOT EXISTS movie;

COMMENT ON SCHEMA movie IS 'Movie content: movies, files, credits, collections, genres, watch progress';

-- Move movie tables to movie schema
-- Note: ALTER TABLE SET SCHEMA also moves indexes, constraints, and triggers
ALTER TABLE public.movies SET SCHEMA movie;

ALTER TABLE public.movie_files SET SCHEMA movie;

ALTER TABLE public.movie_credits SET SCHEMA movie;

ALTER TABLE public.movie_collections SET SCHEMA movie;

ALTER TABLE public.movie_collection_members SET SCHEMA movie;

ALTER TABLE public.movie_genres SET SCHEMA movie;

ALTER TABLE public.movie_watched SET SCHEMA movie;
