-- name: InsertLicense :one
INSERT INTO licenses (pool_id, guid, file, key)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: GetLicenseByGUID :one
SELECT *
FROM licenses
WHERE guid = ?;

-- name: GetUnpooledLicenseByGUID :one
SELECT *
FROM licenses
WHERE guid = ? and pool_id IS NULL;

-- name: GetPooledLicenseByGUID :one
SELECT *
FROM licenses
WHERE guid = ? AND pool_id = ?;

-- name: GetLicenseByNodeID :one
SELECT *
FROM licenses
WHERE node_id = ? and pool_id IS NULL;

-- name: GetPooledLicenseByNodeID :one
SELECT *
FROM licenses
WHERE node_id = ? AND pool_id = ?;

-- name: GetAllLicenses :many
SELECT *
FROM licenses
ORDER BY id;

-- name: GetPooledLicenses :many
SELECT *
FROM licenses
WHERE pool_id = ?
ORDER BY id;

-- name: DeleteLicenseByGUID :one
DELETE FROM licenses
WHERE guid = ?
RETURNING *;

-- name: ReleaseLicenseByNodeID :exec
UPDATE licenses
SET node_id = NULL, last_released_at = unixepoch()
WHERE node_id = ? AND pool_id IS NULL;

-- name: ReleasePooledLicenseByNodeID :exec
UPDATE licenses
SET node_id = NULL, last_released_at = unixepoch()
WHERE node_id = ? AND pool_id = ?;

-- name: ClaimLicenseFIFO :one
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

-- name: ClaimPooledLicenseFIFO :one
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

-- name: ClaimLicenseLIFO :one
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

-- name: ClaimPooledLicenseLIFO :one
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

-- name: ClaimLicenseRandom :one
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

-- name: ClaimPooledLicenseRandom :one
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

