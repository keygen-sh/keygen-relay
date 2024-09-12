-- name: InsertNode :exec
INSERT INTO nodes (fingerprint, claimed_at, last_heartbeat_at)
VALUES (?, NULL, NULL);

-- name: GetNodeByID :one
SELECT id, fingerprint, claimed_at, last_heartbeat_at
FROM nodes
WHERE id = ?;

-- name: GetAllNodes :many
SELECT id, fingerprint, claimed_at, last_heartbeat_at
FROM nodes
ORDER BY id;

-- name: UpdateNodeHeartbeat :exec
UPDATE nodes
SET last_heartbeat_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: ClaimNode :exec
UPDATE nodes
SET claimed_at = CURRENT_TIMESTAMP
WHERE id = ?;
