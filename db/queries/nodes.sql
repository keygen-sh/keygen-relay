-- name: InsertNode :one
INSERT INTO nodes (fingerprint, claimed_at, last_heartbeat_at, created_at)
VALUES (?, NULL, NULL, unixepoch())
RETURNING id, fingerprint, claimed_at, last_heartbeat_at, created_at;

-- name: GetNodeByFingerprint :one
SELECT id, fingerprint, claimed_at, last_heartbeat_at, created_at
FROM nodes
WHERE fingerprint = ?;

-- name: UpdateNodeHeartbeatByFingerprint :exec
UPDATE nodes
SET last_heartbeat_at = unixepoch()
WHERE fingerprint = ?;

-- name: UpdateNodeHeartbeatAndClaimedAtByFingerprint :exec
UPDATE nodes
SET last_heartbeat_at = unixepoch(), claimed_at = unixepoch()
WHERE fingerprint = ?;

-- name: DeleteNodeByFingerprint :exec
DELETE FROM nodes WHERE fingerprint = ?;

-- name: DeleteInactiveNodes :exec
DELETE FROM nodes
WHERE last_heartbeat_at <= strftime('%s', 'now', ?);
