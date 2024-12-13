-- name: InsertAuditLog :exec
INSERT INTO audit_logs (event_type_id, entity_type, entity_id)
VALUES (?, ?, ?);

-- name: GetAuditLogs :many
SELECT id, event_type_id, entity_type, entity_id, created_at
FROM audit_logs
ORDER BY created_at DESC
LIMIT ?;

-- name: GetAuditLogsByEntity :many
SELECT id, event_type_id, entity_type, entity_id, created_at
FROM audit_logs
WHERE entity_id = ? AND entity_type = ?
ORDER BY created_at DESC;
