-- TV Network Queries

-- name: GetNetworkByID :one
SELECT * FROM tv_networks WHERE id = $1;

-- name: GetNetworkByTmdbID :one
SELECT * FROM tv_networks WHERE tmdb_id = $1;

-- name: GetNetworkByName :one
SELECT * FROM tv_networks WHERE name = $1;

-- name: ListNetworks :many
SELECT * FROM tv_networks ORDER BY name ASC LIMIT $1 OFFSET $2;

-- name: ListNetworksWithCounts :many
SELECT n.*, COUNT(snl.series_id) as series_count
FROM tv_networks n
LEFT JOIN series_network_link snl ON n.id = snl.network_id
GROUP BY n.id
ORDER BY series_count DESC, n.name ASC
LIMIT $1 OFFSET $2;

-- name: CreateNetwork :one
INSERT INTO tv_networks (name, logo_path, origin_country, tmdb_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetOrCreateNetwork :one
INSERT INTO tv_networks (name, logo_path, origin_country, tmdb_id)
VALUES ($1, $2, $3, $4)
ON CONFLICT (tmdb_id) DO UPDATE SET
    name = COALESCE(EXCLUDED.name, tv_networks.name),
    logo_path = COALESCE(EXCLUDED.logo_path, tv_networks.logo_path)
RETURNING *;

-- name: UpdateNetwork :one
UPDATE tv_networks SET
    name = COALESCE(sqlc.narg('name'), name),
    logo_path = COALESCE(sqlc.narg('logo_path'), logo_path),
    origin_country = COALESCE(sqlc.narg('origin_country'), origin_country)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteNetwork :exec
DELETE FROM tv_networks WHERE id = $1;

-- Series <-> Network Links

-- name: GetSeriesNetworks :many
SELECT n.* FROM tv_networks n
JOIN series_network_link snl ON n.id = snl.network_id
WHERE snl.series_id = $1
ORDER BY snl.display_order ASC;

-- name: LinkSeriesNetwork :exec
INSERT INTO series_network_link (series_id, network_id, display_order)
VALUES ($1, $2, $3)
ON CONFLICT (series_id, network_id) DO UPDATE SET display_order = $3;

-- name: UnlinkSeriesNetworks :exec
DELETE FROM series_network_link WHERE series_id = $1;

-- name: ListSeriesByNetwork :many
SELECT s.* FROM series s
JOIN series_network_link snl ON s.id = snl.series_id
WHERE snl.network_id = $1
ORDER BY s.sort_title ASC
LIMIT $2 OFFSET $3;

-- name: CountSeriesByNetwork :one
SELECT COUNT(*) FROM series_network_link WHERE network_id = $1;
