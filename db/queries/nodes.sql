-- name: InsertNode :exec
INSERT INTO nodes (fingerprint, claimed_at, last_heartbeat_at)
VALUES (?, NULL, NULL);

-- name: GetNodeByID :one
SELECT id, fingerprint, claimed_at, last_heartbeat_at
FROM nodes
WHERE id = ?;

-- name: GetNodeByFingerprint :one
SELECT id, fingerprint, claimed_at, last_heartbeat_at
FROM nodes
WHERE fingerprint = ?;

-- name: GetAllNodes :many
SELECT id, fingerprint, claimed_at, last_heartbeat_at
FROM nodes
ORDER BY id;

-- name: UpdateNodeHeartbeatByFingerprint :exec
UPDATE nodes
SET last_heartbeat_at = CURRENT_TIMESTAMP
WHERE fingerprint = ?;

-- name: ClaimNode :exec
UPDATE nodes
SET claimed_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeleteNodeByFingerprint :exec
DELETE FROM nodes WHERE fingerprint = ?;
