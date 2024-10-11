-- name: InsertNode :one
INSERT INTO nodes (fingerprint, claimed_at, last_heartbeat_at, created_at)
VALUES (?, NULL, NULL, CURRENT_TIMESTAMP)
RETURNING id, fingerprint, claimed_at, last_heartbeat_at, created_at;

-- name: GetNodeByFingerprint :one
SELECT id, fingerprint, claimed_at, last_heartbeat_at, created_at
FROM nodes
WHERE fingerprint = ?;

-- name: UpdateNodeHeartbeatByFingerprint :exec
UPDATE nodes
SET last_heartbeat_at = CURRENT_TIMESTAMP
WHERE fingerprint = ?;

-- name: UpdateNodeHeartbeatAndClaimedAtByFingerprint :exec
UPDATE nodes
SET last_heartbeat_at = CURRENT_TIMESTAMP, claimed_at = CURRENT_TIMESTAMP
WHERE fingerprint = ?;

-- name: DeleteNodeByFingerprint :exec
DELETE FROM nodes WHERE fingerprint = ?;

-- name: DeleteInactiveNodes :exec
DELETE FROM nodes
WHERE last_heartbeat_at <= datetime('now', ?);
