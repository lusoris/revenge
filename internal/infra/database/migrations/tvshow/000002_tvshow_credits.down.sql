-- TV Show Credits: Rollback
-- Note: video_people and video_credit_role are managed by shared/000017_video_people
BEGIN;

DROP TABLE IF EXISTS episode_credits;
DROP TABLE IF EXISTS series_credits;

COMMIT;
