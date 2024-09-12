-- name: InsertAuditLog :exec
INSERT INTO audit_logs (action, entity_type, entity_id)
VALUES (?, ?, ?);

-- name: GetAuditLogs :many
SELECT id, action, entity_type, entity_id, timestamp
FROM audit_logs
ORDER BY timestamp DESC
LIMIT ?;

-- name: GetAuditLogsByEntity :many
SELECT id, action, entity_type, entity_id, timestamp
FROM audit_logs
WHERE entity_id = ? AND entity_type = ?
ORDER BY timestamp DESC;
