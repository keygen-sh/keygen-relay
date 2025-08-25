-- name: InsertLicense :one
INSERT INTO licenses (pool_id, guid, file, key)
VALUES (?, ?, ?, ?)
RETURNING *;

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

-- name: GetLicenseWithoutPoolByNodeID :one
SELECT *
FROM licenses
WHERE node_id = ? AND pool_id IS NULL;

-- name: GetLicenseWithPoolByNodeID :one
SELECT *
FROM licenses
WHERE node_id = ? AND pool_id = ?;

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

-- name: ReleaseLicenseWithoutPoolByNodeID :exec
UPDATE licenses
SET node_id = NULL, last_released_at = unixepoch()
WHERE node_id = ? AND pool_id IS NULL;

-- name: ReleaseLicenseWithPoolByNodeID :exec
UPDATE licenses
SET node_id = NULL, last_released_at = unixepoch()
WHERE node_id = ? AND pool_id = ?;

-- name: ClaimLicenseWithoutPoolFIFO :one
UPDATE licenses
SET node_id = ?, last_claimed_at = unixepoch(), claims = claims + 1
WHERE id = (
    SELECT l.id
    FROM licenses l
    WHERE l.node_id IS NULL AND l.pool_id IS NULL
    ORDER BY l.created_at ASC
    LIMIT 1
)
RETURNING *;

-- name: ClaimLicenseWithPoolFIFO :one
UPDATE licenses
SET node_id = ?, last_claimed_at = unixepoch(), claims = claims + 1
WHERE id = (
    SELECT l.id
    FROM licenses l
    WHERE l.node_id IS NULL AND l.pool_id = ?
    ORDER BY l.created_at ASC
    LIMIT 1
)
RETURNING *;

-- name: ClaimLicenseWithoutPoolLIFO :one
UPDATE licenses
SET node_id = ?, last_claimed_at = unixepoch(), claims = claims + 1
WHERE id = (
    SELECT l.id
    FROM licenses l
    WHERE l.node_id IS NULL AND l.pool_id IS NULL
    ORDER BY l.created_at DESC
    LIMIT 1
)
RETURNING *;

-- name: ClaimLicenseWithPoolLIFO :one
UPDATE licenses
SET node_id = ?, last_claimed_at = unixepoch(), claims = claims + 1
WHERE id = (
     SELECT l.id
    FROM licenses l
    WHERE l.node_id IS NULL AND l.pool_id = ?
    ORDER BY l.created_at DESC
    LIMIT 1
)
RETURNING *;

-- name: ClaimLicenseWithoutPoolRandom :one
UPDATE licenses
SET node_id = ?, last_claimed_at = unixepoch(), claims = claims + 1
WHERE id = (
    SELECT l.id
    FROM licenses l
    WHERE l.node_id IS NULL AND l.pool_id IS NULL
    ORDER BY RANDOM()
    LIMIT 1
)
RETURNING *;

-- name: ClaimLicenseWithPoolRandom :one
UPDATE licenses
SET node_id = ?, last_claimed_at = unixepoch(), claims = claims + 1
WHERE id = (
    SELECT l.id
    FROM licenses l
    WHERE l.node_id IS NULL AND l.pool_id = ?
    ORDER BY RANDOM()
    LIMIT 1
)
RETURNING *;

-- name: ReleaseLicensesFromDeadNodes :many
UPDATE licenses
SET node_id = NULL, last_released_at = unixepoch()
WHERE node_id IN (
    SELECT id FROM nodes
    WHERE last_heartbeat_at <= strftime('%s', 'now', ?) AND deactivated_at IS NULL
)
RETURNING *;

