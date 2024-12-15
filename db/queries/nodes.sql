-- name: InsertNode :one
INSERT INTO nodes (fingerprint, last_heartbeat_at, created_at)
VALUES (?, NULL, unixepoch())
RETURNING id, fingerprint, last_heartbeat_at, created_at;

-- name: GetNodeByFingerprint :one
SELECT id, fingerprint, last_heartbeat_at, created_at
FROM nodes
WHERE fingerprint = ?;

-- name: PingNodeByFingerprint :exec
UPDATE nodes
SET last_heartbeat_at = unixepoch()
WHERE fingerprint = ?;

-- name: DeleteNodeByFingerprint :exec
DELETE FROM nodes WHERE fingerprint = ?;

-- name: DeleteInactiveNodes :exec
DELETE FROM nodes
WHERE last_heartbeat_at <= strftime('%s', 'now', ?);
