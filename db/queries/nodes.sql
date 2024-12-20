-- name: InsertNode :one
INSERT INTO nodes (fingerprint, created_at)
VALUES (?, unixepoch())
RETURNING *;

-- name: GetNodeByFingerprint :one
SELECT *
FROM nodes
WHERE fingerprint = ? AND deactivated_at IS NULL;

-- name: PingNodeHeartbeatByFingerprint :exec
UPDATE nodes
SET last_heartbeat_at = unixepoch()
WHERE fingerprint = ? AND deactivated_at IS NULL;

-- name: DeactivateNodeByFingerprint :exec
UPDATE nodes
SET deactivated_at = unixepoch()
WHERE fingerprint = ? AND deactivated_at IS NULL;

-- name: DeactivateInactiveNodes :many
UPDATE nodes
SET deactivated_at = unixepoch()
WHERE last_heartbeat_at <= strftime('%s', 'now', ?) AND deactivated_at IS NULL
RETURNING *;
