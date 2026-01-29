-- TV Show Credits Queries
-- Uses shared video_people table

-- Series Cast & Crew

-- name: GetSeriesCast :many
SELECT sc.*, p.name, p.primary_image_url, p.primary_image_blurhash
FROM series_credits sc
JOIN video_people p ON sc.person_id = p.id
WHERE sc.series_id = $1 AND sc.role = 'actor'
ORDER BY sc.billing_order ASC;

-- name: GetSeriesCrew :many
SELECT sc.*, p.name, p.primary_image_url, p.primary_image_blurhash
FROM series_credits sc
JOIN video_people p ON sc.person_id = p.id
WHERE sc.series_id = $1 AND sc.role != 'actor'
ORDER BY
    CASE sc.role
        WHEN 'creator' THEN 1
        WHEN 'showrunner' THEN 2
        WHEN 'executive_producer' THEN 3
        WHEN 'director' THEN 4
        WHEN 'writer' THEN 5
        WHEN 'producer' THEN 6
        ELSE 10
    END,
    sc.billing_order ASC;

-- name: GetSeriesCreators :many
SELECT sc.*, p.name, p.primary_image_url, p.primary_image_blurhash
FROM series_credits sc
JOIN video_people p ON sc.person_id = p.id
WHERE sc.series_id = $1 AND sc.role IN ('creator', 'showrunner')
ORDER BY sc.billing_order ASC;

-- name: GetPersonSeriesCredits :many
SELECT sc.*, s.title, s.year, s.poster_path, s.poster_blurhash
FROM series_credits sc
JOIN series s ON sc.series_id = s.id
WHERE sc.person_id = $1
ORDER BY s.first_air_date DESC NULLS LAST;

-- name: CreateSeriesCredit :one
INSERT INTO series_credits (
    series_id, person_id, role, character_name, department, job, billing_order, tmdb_credit_id
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: DeleteSeriesCredits :exec
DELETE FROM series_credits WHERE series_id = $1;

-- name: DeleteSeriesCreditsByRole :exec
DELETE FROM series_credits WHERE series_id = $1 AND role = $2;

-- name: CountSeriesCast :one
SELECT COUNT(*) FROM series_credits WHERE series_id = $1 AND role = 'actor';

-- name: CountSeriesCrew :one
SELECT COUNT(*) FROM series_credits WHERE series_id = $1 AND role != 'actor';

-- Episode Cast & Crew (Guest Stars, Episode Directors, Writers)

-- name: GetEpisodeCast :many
SELECT ec.*, p.name, p.primary_image_url, p.primary_image_blurhash
FROM episode_credits ec
JOIN video_people p ON ec.person_id = p.id
WHERE ec.episode_id = $1 AND ec.role = 'actor'
ORDER BY ec.is_guest ASC, ec.billing_order ASC;

-- name: GetEpisodeGuestStars :many
SELECT ec.*, p.name, p.primary_image_url, p.primary_image_blurhash
FROM episode_credits ec
JOIN video_people p ON ec.person_id = p.id
WHERE ec.episode_id = $1 AND ec.role = 'actor' AND ec.is_guest = true
ORDER BY ec.billing_order ASC;

-- name: GetEpisodeCrew :many
SELECT ec.*, p.name, p.primary_image_url, p.primary_image_blurhash
FROM episode_credits ec
JOIN video_people p ON ec.person_id = p.id
WHERE ec.episode_id = $1 AND ec.role != 'actor'
ORDER BY
    CASE ec.role
        WHEN 'director' THEN 1
        WHEN 'writer' THEN 2
        ELSE 10
    END,
    ec.billing_order ASC;

-- name: GetEpisodeDirectors :many
SELECT ec.*, p.name, p.primary_image_url, p.primary_image_blurhash
FROM episode_credits ec
JOIN video_people p ON ec.person_id = p.id
WHERE ec.episode_id = $1 AND ec.role = 'director'
ORDER BY ec.billing_order ASC;

-- name: GetEpisodeWriters :many
SELECT ec.*, p.name, p.primary_image_url, p.primary_image_blurhash
FROM episode_credits ec
JOIN video_people p ON ec.person_id = p.id
WHERE ec.episode_id = $1 AND ec.role = 'writer'
ORDER BY ec.billing_order ASC;

-- name: GetPersonEpisodeCredits :many
SELECT ec.*, e.title as episode_title, e.season_number, e.episode_number,
       s.id as series_id, s.title as series_title, s.poster_path as series_poster
FROM episode_credits ec
JOIN episodes e ON ec.episode_id = e.id
JOIN series s ON e.series_id = s.id
WHERE ec.person_id = $1
ORDER BY e.air_date DESC NULLS LAST;

-- name: CreateEpisodeCredit :one
INSERT INTO episode_credits (
    episode_id, person_id, role, character_name, department, job, billing_order, is_guest, tmdb_credit_id
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: DeleteEpisodeCredits :exec
DELETE FROM episode_credits WHERE episode_id = $1;

-- name: DeleteEpisodeCreditsBySeries :exec
DELETE FROM episode_credits WHERE episode_id IN (SELECT id FROM episodes WHERE series_id = $1);

-- name: CountEpisodeGuestStars :one
SELECT COUNT(*) FROM episode_credits WHERE episode_id = $1 AND role = 'actor' AND is_guest = true;
