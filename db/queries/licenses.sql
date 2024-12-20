-- name: InsertLicense :one
INSERT INTO licenses (guid, file, key)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetLicenseByGUID :one
SELECT *
FROM licenses
WHERE guid = ?;

-- name: GetLicenseByNodeID :one
SELECT *
FROM licenses
WHERE node_id = ?;

-- name: GetAllLicenses :many
SELECT *
FROM licenses
ORDER BY id;

-- name: DeleteLicenseByGUID :one
DELETE FROM licenses
WHERE guid = ?
RETURNING *;

-- name: ReleaseLicenseByNodeID :exec
UPDATE licenses
SET node_id = NULL, last_released_at = unixepoch()
WHERE node_id = ?;

-- name: ClaimLicenseFIFO :one
UPDATE licenses
SET node_id = ?, last_claimed_at = unixepoch(), claims = claims + 1
WHERE id = (
    SELECT id
    FROM licenses
    WHERE node_id IS NULL
    ORDER BY created_at ASC
    LIMIT 1
)
RETURNING *;

-- name: ClaimLicenseLIFO :one
UPDATE licenses
SET node_id = ?, last_claimed_at = unixepoch(), claims = claims + 1
WHERE id = (
    SELECT id
    FROM licenses
    WHERE node_id IS NULL
    ORDER BY created_at DESC
    LIMIT 1
)
RETURNING *;

-- name: ClaimLicenseRandom :one
UPDATE licenses
SET node_id = ?, last_claimed_at = unixepoch(), claims = claims + 1
WHERE id = (
    SELECT id
    FROM licenses
    WHERE node_id IS NULL
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

