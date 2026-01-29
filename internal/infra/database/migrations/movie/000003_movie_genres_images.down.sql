BEGIN;

DROP TABLE IF EXISTS movie_videos;
DROP TYPE IF EXISTS movie_video_site;
DROP TYPE IF EXISTS movie_video_type;
DROP TABLE IF EXISTS movie_images;
DROP TYPE IF EXISTS movie_image_type;
DROP TABLE IF EXISTS movie_genre_link;
DROP TRIGGER IF EXISTS movie_genres_updated_at ON movie_genres;
DROP TABLE IF EXISTS movie_genres;

COMMIT;
