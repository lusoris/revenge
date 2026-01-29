-- name: GetAdultMovieByID :one
SELECT * FROM qar.movies WHERE id = $1;

-- name: ListAdultMovies :many
SELECT * FROM qar.movies
ORDER BY title
LIMIT $1 OFFSET $2;

-- name: ListAdultMoviesByLibrary :many
SELECT * FROM qar.movies
WHERE library_id = $1
ORDER BY title
LIMIT $2 OFFSET $3;

-- name: CreateAdultMovie :one
INSERT INTO qar.movies (
    library_id,
    title,
    sort_title,
    original_title,
    overview,
    release_date,
    runtime_ticks,
    studio_id,
    director,
    series,
    path,
    size_bytes,
    container,
    video_codec,
    audio_codec,
    resolution,
    phash,
    oshash,
    whisparr_id,
    stashdb_id,
    tpdb_id,
    has_file,
    is_hdr,
    is_3d
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8,
    $9, $10, $11, $12, $13, $14, $15, $16,
    $17, $18, $19, $20, $21, $22, $23, $24
)
RETURNING *;

-- name: UpdateAdultMovie :one
UPDATE qar.movies
SET
    library_id = $2,
    title = $3,
    sort_title = $4,
    original_title = $5,
    overview = $6,
    release_date = $7,
    runtime_ticks = $8,
    studio_id = $9,
    director = $10,
    series = $11,
    path = $12,
    size_bytes = $13,
    container = $14,
    video_codec = $15,
    audio_codec = $16,
    resolution = $17,
    phash = $18,
    oshash = $19,
    whisparr_id = $20,
    stashdb_id = $21,
    tpdb_id = $22,
    has_file = $23,
    is_hdr = $24,
    is_3d = $25
WHERE id = $1
RETURNING *;

-- name: DeleteAdultMovie :exec
DELETE FROM qar.movies WHERE id = $1;
