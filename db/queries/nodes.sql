-- name: ActivateNode :one
INSERT INTO nodes (fingerprint)
VALUES (?)
ON CONFLICT (fingerprint) DO UPDATE SET deactivated_at = NULL
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

-- name: DeactivateDeadNodes :many
UPDATE nodes
SET deactivated_at = unixepoch()
WHERE last_heartbeat_at <= strftime('%s', 'now', ?) AND deactivated_at IS NULL
RETURNING *;
