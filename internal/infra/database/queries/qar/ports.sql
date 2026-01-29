-- Ports (Studios) - QAR obfuscated queries

-- name: GetPortByID :one
SELECT * FROM qar.ports WHERE id = $1;

-- name: ListPorts :many
SELECT * FROM qar.ports
ORDER BY name
LIMIT $1 OFFSET $2;

-- name: CreatePort :one
INSERT INTO qar.ports (
    name,
    parent_id,
    stashdb_id,
    tpdb_id,
    url,
    logo_path
) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdatePort :one
UPDATE qar.ports
SET
    name = $2,
    parent_id = $3,
    stashdb_id = $4,
    tpdb_id = $5,
    url = $6,
    logo_path = $7
WHERE id = $1
RETURNING *;

-- name: DeletePort :exec
DELETE FROM qar.ports WHERE id = $1;

-- name: GetPortByStashDBID :one
SELECT * FROM qar.ports WHERE stashdb_id = $1;

-- name: GetPortByTPDBID :one
SELECT * FROM qar.ports WHERE tpdb_id = $1;

-- name: SearchPorts :many
SELECT * FROM qar.ports
WHERE name ILIKE '%' || $1 || '%'
ORDER BY name
LIMIT $2 OFFSET $3;

-- name: ListPortChildren :many
SELECT * FROM qar.ports
WHERE parent_id = $1
ORDER BY name;

-- name: ListRootPorts :many
SELECT * FROM qar.ports
WHERE parent_id IS NULL
ORDER BY name
LIMIT $1 OFFSET $2;

-- name: CountExpeditionsByPort :one
SELECT COUNT(*) FROM qar.expeditions WHERE port_id = $1;

-- name: CountVoyagesByPort :one
SELECT COUNT(*) FROM qar.voyages WHERE port_id = $1;
