-- name: InsertLicense :exec
INSERT INTO licenses (id, file, key, claims, node_id)
VALUES (?, ?, ?, ?, NULL);

-- name: GetLicenseByID :one
SELECT *
FROM licenses
WHERE id = ?;

-- name: GetLicenseByNodeID :one
SELECT *
FROM licenses
WHERE node_id = ?;

-- name: GetAllLicenses :many
SELECT *
FROM licenses
ORDER BY id;

-- name: DeleteLicenseByID :exec
DELETE FROM licenses
WHERE id = ?;

-- name: ClaimLicense :exec
UPDATE licenses
SET node_id = ?, last_claimed_at = unixepoch(), claims = claims + 1
WHERE id = ? AND node_id IS NULL;

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

-- name: ReleaseLicensesFromInactiveNodes :many
UPDATE licenses
SET node_id = NULL, last_released_at = unixepoch()
WHERE node_id IN (
    SELECT id FROM nodes
    WHERE last_heartbeat_at <= strftime('%s', 'now', ?) AND deactivated_at IS NULL
)
RETURNING *;

