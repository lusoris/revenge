-- TV Show Genres and Images: Rollback
BEGIN;

DROP TABLE IF EXISTS series_videos;
DROP TABLE IF EXISTS episode_images;
DROP TABLE IF EXISTS season_images;
DROP TABLE IF EXISTS series_images;
DROP TABLE IF EXISTS series_genre_link;
DROP TABLE IF EXISTS tvshow_genres;

COMMIT;
