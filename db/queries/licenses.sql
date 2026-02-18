-- name: InsertLicense :one
INSERT INTO licenses (pool_id, guid, file, key, seats)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: GetLicenseByID :one
SELECT *
FROM licenses
WHERE id = ?;

-- name: GetLicenseByGUID :one
SELECT *
FROM licenses
WHERE guid = ?;

-- name: GetLicenseWithoutPoolByGUID :one
SELECT *
FROM licenses
WHERE guid = ? AND pool_id IS NULL;

-- name: GetLicenseWithPoolByGUID :one
SELECT *
FROM licenses
WHERE guid = ? AND pool_id = ?;

-- name: GetLicenseByNodeIDWithoutPool :one
SELECT l.*
FROM licenses l
JOIN license_nodes ln ON ln.license_id = l.id
WHERE ln.node_id = ? AND l.pool_id IS NULL
LIMIT 1;

-- name: GetLicenseByNodeIDWithPool :one
SELECT l.*
FROM licenses l
JOIN license_nodes ln ON ln.license_id = l.id
WHERE ln.node_id = ? AND l.pool_id = ?
LIMIT 1;

-- name: GetLicenses :many
SELECT *
FROM licenses
ORDER BY id;

-- name: GetLicensesWithoutPool :many
SELECT *
FROM licenses
WHERE pool_id IS NULL
ORDER BY id;

-- name: GetLicensesWithPool :many
SELECT *
FROM licenses
WHERE pool_id = ?
ORDER BY id;

-- name: DeleteLicenseByGUID :one
DELETE FROM licenses
WHERE guid = ?
RETURNING *;

-- name: IncrementLicenseClaims :exec
UPDATE licenses
SET claims = claims + 1, last_claimed_at = unixepoch()
WHERE id = ?;

-- name: TouchLicenseLastReleasedAt :exec
UPDATE licenses
SET last_released_at = unixepoch()
WHERE id = ?;

-- name: ClaimLicenseNodeWithoutPoolFIFO :one
INSERT INTO license_nodes (license_id, node_id)
SELECT l.id, ?
FROM licenses l
WHERE l.pool_id IS NULL
  AND (SELECT COUNT(*) FROM license_nodes ln WHERE ln.license_id = l.id) < l.seats
ORDER BY l.created_at ASC
LIMIT 1
RETURNING license_id;

-- name: ClaimLicenseNodeWithPoolFIFO :one
INSERT INTO license_nodes (license_id, node_id)
SELECT l.id, ?
FROM licenses l
WHERE l.pool_id = ?
  AND (SELECT COUNT(*) FROM license_nodes ln WHERE ln.license_id = l.id) < l.seats
ORDER BY l.created_at ASC
LIMIT 1
RETURNING license_id;

-- name: ClaimLicenseNodeWithoutPoolLIFO :one
INSERT INTO license_nodes (license_id, node_id)
SELECT l.id, ?
FROM licenses l
WHERE l.pool_id IS NULL
  AND (SELECT COUNT(*) FROM license_nodes ln WHERE ln.license_id = l.id) < l.seats
ORDER BY l.created_at DESC
LIMIT 1
RETURNING license_id;

-- name: ClaimLicenseNodeWithPoolLIFO :one
INSERT INTO license_nodes (license_id, node_id)
SELECT l.id, ?
FROM licenses l
WHERE l.pool_id = ?
  AND (SELECT COUNT(*) FROM license_nodes ln WHERE ln.license_id = l.id) < l.seats
ORDER BY l.created_at DESC
LIMIT 1
RETURNING license_id;

-- name: ClaimLicenseNodeWithoutPoolRandom :one
INSERT INTO license_nodes (license_id, node_id)
SELECT l.id, ?
FROM licenses l
WHERE l.pool_id IS NULL
  AND (SELECT COUNT(*) FROM license_nodes ln WHERE ln.license_id = l.id) < l.seats
ORDER BY RANDOM()
LIMIT 1
RETURNING license_id;

-- name: ClaimLicenseNodeWithPoolRandom :one
INSERT INTO license_nodes (license_id, node_id)
SELECT l.id, ?
FROM licenses l
WHERE l.pool_id = ?
  AND (SELECT COUNT(*) FROM license_nodes ln WHERE ln.license_id = l.id) < l.seats
ORDER BY RANDOM()
LIMIT 1
RETURNING license_id;

-- name: ReleaseLicenseNodeWithoutPoolByNodeID :exec
DELETE FROM license_nodes
WHERE node_id = ?
  AND license_id IN (SELECT id FROM licenses WHERE pool_id IS NULL);

-- name: ReleaseLicenseNodeWithPoolByNodeID :exec
DELETE FROM license_nodes
WHERE node_id = ?
  AND license_id IN (SELECT id FROM licenses WHERE pool_id = ?);

-- name: ReleaseLicenseNodesFromDeadNodes :many
DELETE FROM license_nodes
WHERE node_id IN (
    SELECT id FROM nodes
    WHERE last_heartbeat_at <= strftime('%s', 'now', ?) AND deactivated_at IS NULL
)
RETURNING license_id;
