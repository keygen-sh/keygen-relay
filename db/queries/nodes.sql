-- name: InsertNode :one
INSERT INTO nodes (fingerprint, claimed_at, last_heartbeat_at, created_at)
VALUES (?, NULL, NULL, CURRENT_TIMESTAMP)
RETURNING id, fingerprint, claimed_at, last_heartbeat_at, created_at;

-- name: GetNodeByID :one
SELECT id, fingerprint, claimed_at, last_heartbeat_at, created_at
FROM nodes
WHERE id = ?;

-- name: GetNodeByFingerprint :one
SELECT id, fingerprint, claimed_at, last_heartbeat_at, created_at
FROM nodes
WHERE fingerprint = ?;

-- name: GetAllNodes :many
SELECT id, fingerprint, claimed_at, last_heartbeat_at, created_at
FROM nodes
ORDER BY id;

-- name: UpdateNodeHeartbeatByFingerprint :exec
UPDATE nodes
SET last_heartbeat_at = CURRENT_TIMESTAMP
WHERE fingerprint = ?;

-- name: UpdateNodeHeartbeatAndClaimedAtByFingerprint :exec
UPDATE nodes
SET last_heartbeat_at = CURRENT_TIMESTAMP, claimed_at = CURRENT_TIMESTAMP
WHERE fingerprint = ?;

-- name: ClaimNode :exec
UPDATE nodes
SET claimed_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeleteNodeByFingerprint :exec
DELETE FROM nodes WHERE fingerprint = ?;

-- name: GetInactiveNodes :many
SELECT id, fingerprint, claimed_at, last_heartbeat_at, created_at
FROM nodes
WHERE datetime(last_heartbeat_at) <= datetime('now', ?);

-- name: DeleteInactiveNodes :exec
DELETE FROM nodes
WHERE datetime(last_heartbeat_at) <= datetime('now', ?);
