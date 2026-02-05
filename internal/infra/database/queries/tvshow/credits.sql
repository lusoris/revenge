-- Series Credits

-- name: ListSeriesCast :many
SELECT * FROM tvshow.series_credits
WHERE series_id = $1 AND credit_type = 'cast'
ORDER BY cast_order ASC NULLS LAST, name ASC;

-- name: ListSeriesCrew :many
SELECT * FROM tvshow.series_credits
WHERE series_id = $1 AND credit_type = 'crew'
ORDER BY department ASC, name ASC;

-- name: CreateSeriesCredit :one
INSERT INTO tvshow.series_credits (
    series_id, tmdb_person_id, name, credit_type,
    character, cast_order, job, department, profile_path
) VALUES (
    $1, $2, $3, $4,
    $5, $6, $7, $8, $9
)
RETURNING *;

-- name: DeleteSeriesCredits :exec
DELETE FROM tvshow.series_credits WHERE series_id = $1;

-- Episode Credits (Guest Stars)

-- name: ListEpisodeGuestStars :many
SELECT * FROM tvshow.episode_credits
WHERE episode_id = $1 AND credit_type = 'guest_star'
ORDER BY cast_order ASC NULLS LAST, name ASC;

-- name: ListEpisodeCrew :many
SELECT * FROM tvshow.episode_credits
WHERE episode_id = $1 AND credit_type = 'crew'
ORDER BY department ASC, name ASC;

-- name: CreateEpisodeCredit :one
INSERT INTO tvshow.episode_credits (
    episode_id, tmdb_person_id, name, credit_type,
    character, cast_order, job, department, profile_path
) VALUES (
    $1, $2, $3, $4,
    $5, $6, $7, $8, $9
)
RETURNING *;

-- name: DeleteEpisodeCredits :exec
DELETE FROM tvshow.episode_credits WHERE episode_id = $1;
