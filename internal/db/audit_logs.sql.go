// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: audit_logs.sql

package db

import (
	"context"
)

const getAuditLogs = `-- name: GetAuditLogs :many
SELECT id, action, entity_type, entity_id, timestamp
FROM audit_logs
ORDER BY timestamp DESC
LIMIT ?
`

func (q *Queries) GetAuditLogs(ctx context.Context, limit int64) ([]AuditLog, error) {
	rows, err := q.db.QueryContext(ctx, getAuditLogs, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AuditLog
	for rows.Next() {
		var i AuditLog
		if err := rows.Scan(
			&i.ID,
			&i.Action,
			&i.EntityType,
			&i.EntityID,
			&i.Timestamp,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAuditLogsByEntity = `-- name: GetAuditLogsByEntity :many
SELECT id, action, entity_type, entity_id, timestamp
FROM audit_logs
WHERE entity_id = ? AND entity_type = ?
ORDER BY timestamp DESC
`

type GetAuditLogsByEntityParams struct {
	EntityID   string `json:"entity_id"`
	EntityType string `json:"entity_type"`
}

func (q *Queries) GetAuditLogsByEntity(ctx context.Context, arg GetAuditLogsByEntityParams) ([]AuditLog, error) {
	rows, err := q.db.QueryContext(ctx, getAuditLogsByEntity, arg.EntityID, arg.EntityType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AuditLog
	for rows.Next() {
		var i AuditLog
		if err := rows.Scan(
			&i.ID,
			&i.Action,
			&i.EntityType,
			&i.EntityID,
			&i.Timestamp,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertAuditLog = `-- name: InsertAuditLog :exec
INSERT INTO audit_logs (action, entity_type, entity_id)
VALUES (?, ?, ?)
`

type InsertAuditLogParams struct {
	Action     string `json:"action"`
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
}

func (q *Queries) InsertAuditLog(ctx context.Context, arg InsertAuditLogParams) error {
	_, err := q.db.ExecContext(ctx, insertAuditLog, arg.Action, arg.EntityType, arg.EntityID)
	return err
}