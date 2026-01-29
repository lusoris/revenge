-- Movie Credits Queries

-- name: GetMovieCast :many
SELECT mc.*, p.name, p.primary_image_url, p.primary_image_blurhash
FROM movie_credits mc
JOIN movie_people p ON mc.person_id = p.id
WHERE mc.movie_id = $1 AND mc.role = 'actor'
ORDER BY mc.billing_order ASC;

-- name: GetMovieCrew :many
SELECT mc.*, p.name, p.primary_image_url, p.primary_image_blurhash
FROM movie_credits mc
JOIN movie_people p ON mc.person_id = p.id
WHERE mc.movie_id = $1 AND mc.role != 'actor'
ORDER BY
    CASE mc.role
        WHEN 'director' THEN 1
        WHEN 'writer' THEN 2
        WHEN 'producer' THEN 3
        ELSE 10
    END,
    mc.billing_order ASC;

-- name: GetMovieDirectors :many
SELECT mc.*, p.name, p.primary_image_url, p.primary_image_blurhash
FROM movie_credits mc
JOIN movie_people p ON mc.person_id = p.id
WHERE mc.movie_id = $1 AND mc.role = 'director'
ORDER BY mc.billing_order ASC;

-- name: GetMovieWriters :many
SELECT mc.*, p.name, p.primary_image_url, p.primary_image_blurhash
FROM movie_credits mc
JOIN movie_people p ON mc.person_id = p.id
WHERE mc.movie_id = $1 AND mc.role = 'writer'
ORDER BY mc.billing_order ASC;

-- name: GetPersonMovieCredits :many
SELECT mc.*, m.title, m.year, m.poster_path, m.poster_blurhash
FROM movie_credits mc
JOIN movies m ON mc.movie_id = m.id
WHERE mc.person_id = $1
ORDER BY m.release_date DESC NULLS LAST;

-- name: CreateMovieCredit :one
INSERT INTO movie_credits (
    movie_id, person_id, role, character_name, department, job, billing_order, is_guest, tmdb_credit_id
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: DeleteMovieCredits :exec
DELETE FROM movie_credits WHERE movie_id = $1;

-- name: DeleteMovieCreditsByRole :exec
DELETE FROM movie_credits WHERE movie_id = $1 AND role = $2;

-- name: CountMovieCast :one
SELECT COUNT(*) FROM movie_credits WHERE movie_id = $1 AND role = 'actor';

-- name: CountMovieCrew :one
SELECT COUNT(*) FROM movie_credits WHERE movie_id = $1 AND role != 'actor';
