-- Drop movie collections tables
DROP TRIGGER IF EXISTS update_movie_collections_updated_at ON public.movie_collections;

DROP INDEX IF EXISTS idx_collection_members_order;
DROP INDEX IF EXISTS idx_collection_members_movie;
DROP INDEX IF EXISTS idx_collection_members_collection;
DROP INDEX IF EXISTS idx_collections_name_trgm;
DROP INDEX IF EXISTS idx_collections_tmdb_id;

DROP TABLE IF EXISTS public.movie_collection_members;
DROP TABLE IF EXISTS public.movie_collections;
