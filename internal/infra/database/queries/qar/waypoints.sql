-- Waypoints (Scene Markers) - QAR obfuscated queries

-- name: GetWaypointByID :one
SELECT * FROM qar.voyage_waypoints WHERE id = $1;

-- name: ListWaypointsByVoyage :many
SELECT * FROM qar.voyage_waypoints
WHERE voyage_id = $1
ORDER BY start_seconds;

-- name: CreateWaypoint :one
INSERT INTO qar.voyage_waypoints (
    voyage_id,
    title,
    start_seconds,
    end_seconds,
    flag_id,
    stash_marker_id
) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateWaypoint :one
UPDATE qar.voyage_waypoints
SET
    title = $2,
    start_seconds = $3,
    end_seconds = $4,
    flag_id = $5
WHERE id = $1
RETURNING *;

-- name: DeleteWaypoint :exec
DELETE FROM qar.voyage_waypoints WHERE id = $1;

-- name: DeleteWaypointsByVoyage :exec
DELETE FROM qar.voyage_waypoints WHERE voyage_id = $1;

-- name: GetWaypointByStashMarkerID :one
SELECT * FROM qar.voyage_waypoints WHERE stash_marker_id = $1;
