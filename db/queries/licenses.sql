-- name: InsertLicense :exec
INSERT INTO licenses (id, file, key, claims, node_id)
VALUES (?, ?, ?, ?, NULL);

-- name: GetLicenseByID :one
SELECT id, file, key, claims, last_claimed_at, last_released_at, node_id
FROM licenses
WHERE id = ?;

-- name: GetAllLicenses :many
SELECT id, file, key, claims, last_claimed_at, last_released_at, node_id
FROM licenses
ORDER BY id;

-- name: DeleteLicenseByID :exec
DELETE FROM licenses
WHERE id = ?;

-- name: ClaimLicense :exec
UPDATE licenses
SET node_id = ?, last_claimed_at = CURRENT_TIMESTAMP, claims = claims + 1
WHERE id = ? AND node_id IS NULL;

-- name: ReleaseLicenseByID :exec
UPDATE licenses
SET node_id = NULL, last_released_at = CURRENT_TIMESTAMP
WHERE id = ? AND node_id IS NOT NULL;

-- name: ReleaseLicenseByNodeID :exec
UPDATE licenses
SET node_id = NULL, last_released_at = CURRENT_TIMESTAMP
WHERE node_id = ?;

-- name: GetUnclaimedLicense :one
SELECT id, file, key, claims, last_claimed_at, last_released_at, node_id
FROM licenses
WHERE node_id IS NULL
LIMIT 1;
