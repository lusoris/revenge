-- name: GetAdultSceneByID :one
SELECT * FROM qar.scenes WHERE id = $1;

-- name: ListAdultScenes :many
SELECT * FROM qar.scenes
ORDER BY title
LIMIT $1 OFFSET $2;

-- name: ListAdultScenesByLibrary :many
SELECT * FROM qar.scenes
WHERE library_id = $1
ORDER BY title
LIMIT $2 OFFSET $3;

-- name: CreateAdultScene :one
INSERT INTO qar.scenes (
    library_id,
    title,
    sort_title,
    overview,
    release_date,
    runtime_minutes,
    studio_id,
    whisparr_id,
    stash_id,
    stashdb_id,
    tpdb_id,
    path,
    size_bytes,
    video_codec,
    audio_codec,
    resolution,
    oshash,
    phash,
    md5,
    cover_path
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
    $11, $12, $13, $14, $15, $16, $17, $18, $19, $20
)
RETURNING *;

-- name: UpdateAdultScene :one
UPDATE qar.scenes
SET
    library_id = $2,
    title = $3,
    sort_title = $4,
    overview = $5,
    release_date = $6,
    runtime_minutes = $7,
    studio_id = $8,
    whisparr_id = $9,
    stash_id = $10,
    stashdb_id = $11,
    tpdb_id = $12,
    path = $13,
    size_bytes = $14,
    video_codec = $15,
    audio_codec = $16,
    resolution = $17,
    oshash = $18,
    phash = $19,
    md5 = $20,
    cover_path = $21
WHERE id = $1
RETURNING *;

-- name: DeleteAdultScene :exec
DELETE FROM qar.scenes WHERE id = $1;
