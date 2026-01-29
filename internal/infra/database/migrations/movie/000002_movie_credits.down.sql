-- Movie Credits: Rollback
-- Note: video_people and video_credit_role are managed by shared/000017_video_people
BEGIN;

DROP TABLE IF EXISTS movie_credits;

COMMIT;
