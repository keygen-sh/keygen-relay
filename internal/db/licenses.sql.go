// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: licenses.sql

package db

import (
	"context"
)

const claimLicense = `-- name: ClaimLicense :exec
UPDATE licenses
SET node_id = ?, last_claimed_at = unixepoch(), claims = claims + 1
WHERE id = ? AND node_id IS NULL
`

type ClaimLicenseParams struct {
	NodeID *int64
	ID     string
}

func (q *Queries) ClaimLicense(ctx context.Context, arg ClaimLicenseParams) error {
	_, err := q.db.ExecContext(ctx, claimLicense, arg.NodeID, arg.ID)
	return err
}

const claimLicenseFIFO = `-- name: ClaimLicenseFIFO :one
UPDATE licenses
SET node_id = ?, last_claimed_at = unixepoch(), claims = claims + 1
WHERE id = (
    SELECT id
    FROM licenses
    WHERE node_id IS NULL
    ORDER BY created_at ASC
    LIMIT 1
)
RETURNING id, file, "key", claims, last_claimed_at, last_released_at, node_id, created_at
`

func (q *Queries) ClaimLicenseFIFO(ctx context.Context, nodeID *int64) (License, error) {
	row := q.db.QueryRowContext(ctx, claimLicenseFIFO, nodeID)
	var i License
	err := row.Scan(
		&i.ID,
		&i.File,
		&i.Key,
		&i.Claims,
		&i.LastClaimedAt,
		&i.LastReleasedAt,
		&i.NodeID,
		&i.CreatedAt,
	)
	return i, err
}

const claimLicenseLIFO = `-- name: ClaimLicenseLIFO :one
UPDATE licenses
SET node_id = ?, last_claimed_at = unixepoch(), claims = claims + 1
WHERE id = (
    SELECT id
    FROM licenses
    WHERE node_id IS NULL
    ORDER BY created_at DESC
    LIMIT 1
)
RETURNING id, file, "key", claims, last_claimed_at, last_released_at, node_id, created_at
`

func (q *Queries) ClaimLicenseLIFO(ctx context.Context, nodeID *int64) (License, error) {
	row := q.db.QueryRowContext(ctx, claimLicenseLIFO, nodeID)
	var i License
	err := row.Scan(
		&i.ID,
		&i.File,
		&i.Key,
		&i.Claims,
		&i.LastClaimedAt,
		&i.LastReleasedAt,
		&i.NodeID,
		&i.CreatedAt,
	)
	return i, err
}

const claimLicenseRandom = `-- name: ClaimLicenseRandom :one
UPDATE licenses
SET node_id = ?, last_claimed_at = unixepoch(), claims = claims + 1
WHERE id = (
    SELECT id
    FROM licenses
    WHERE node_id IS NULL
    ORDER BY RANDOM()
    LIMIT 1
)
RETURNING id, file, "key", claims, last_claimed_at, last_released_at, node_id, created_at
`

func (q *Queries) ClaimLicenseRandom(ctx context.Context, nodeID *int64) (License, error) {
	row := q.db.QueryRowContext(ctx, claimLicenseRandom, nodeID)
	var i License
	err := row.Scan(
		&i.ID,
		&i.File,
		&i.Key,
		&i.Claims,
		&i.LastClaimedAt,
		&i.LastReleasedAt,
		&i.NodeID,
		&i.CreatedAt,
	)
	return i, err
}

const deleteLicenseByID = `-- name: DeleteLicenseByID :exec
DELETE FROM licenses
WHERE id = ?
`

func (q *Queries) DeleteLicenseByID(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, deleteLicenseByID, id)
	return err
}

const getAllLicenses = `-- name: GetAllLicenses :many
SELECT id, file, "key", claims, last_claimed_at, last_released_at, node_id, created_at
FROM licenses
ORDER BY id
`

func (q *Queries) GetAllLicenses(ctx context.Context) ([]License, error) {
	rows, err := q.db.QueryContext(ctx, getAllLicenses)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []License
	for rows.Next() {
		var i License
		if err := rows.Scan(
			&i.ID,
			&i.File,
			&i.Key,
			&i.Claims,
			&i.LastClaimedAt,
			&i.LastReleasedAt,
			&i.NodeID,
			&i.CreatedAt,
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

const getLicenseByID = `-- name: GetLicenseByID :one
SELECT id, file, "key", claims, last_claimed_at, last_released_at, node_id, created_at
FROM licenses
WHERE id = ?
`

func (q *Queries) GetLicenseByID(ctx context.Context, id string) (License, error) {
	row := q.db.QueryRowContext(ctx, getLicenseByID, id)
	var i License
	err := row.Scan(
		&i.ID,
		&i.File,
		&i.Key,
		&i.Claims,
		&i.LastClaimedAt,
		&i.LastReleasedAt,
		&i.NodeID,
		&i.CreatedAt,
	)
	return i, err
}

const getLicenseByNodeID = `-- name: GetLicenseByNodeID :one
SELECT id, file, "key", claims, last_claimed_at, last_released_at, node_id, created_at
FROM licenses
WHERE node_id = ?
`

func (q *Queries) GetLicenseByNodeID(ctx context.Context, nodeID *int64) (License, error) {
	row := q.db.QueryRowContext(ctx, getLicenseByNodeID, nodeID)
	var i License
	err := row.Scan(
		&i.ID,
		&i.File,
		&i.Key,
		&i.Claims,
		&i.LastClaimedAt,
		&i.LastReleasedAt,
		&i.NodeID,
		&i.CreatedAt,
	)
	return i, err
}

const insertLicense = `-- name: InsertLicense :exec
INSERT INTO licenses (id, file, key, claims, node_id)
VALUES (?, ?, ?, ?, NULL)
`

type InsertLicenseParams struct {
	ID     string
	File   []byte
	Key    string
	Claims int64
}

func (q *Queries) InsertLicense(ctx context.Context, arg InsertLicenseParams) error {
	_, err := q.db.ExecContext(ctx, insertLicense,
		arg.ID,
		arg.File,
		arg.Key,
		arg.Claims,
	)
	return err
}

const releaseLicenseByNodeID = `-- name: ReleaseLicenseByNodeID :exec
UPDATE licenses
SET node_id = NULL, last_released_at = unixepoch()
WHERE node_id = ?
`

func (q *Queries) ReleaseLicenseByNodeID(ctx context.Context, nodeID *int64) error {
	_, err := q.db.ExecContext(ctx, releaseLicenseByNodeID, nodeID)
	return err
}

const releaseLicensesFromDeadNodes = `-- name: ReleaseLicensesFromDeadNodes :many
UPDATE licenses
SET node_id = NULL, last_released_at = unixepoch()
WHERE node_id IN (
    SELECT id FROM nodes
    WHERE last_heartbeat_at <= strftime('%s', 'now', ?) AND deactivated_at IS NULL
)
RETURNING id, file, "key", claims, last_claimed_at, last_released_at, node_id, created_at
`

func (q *Queries) ReleaseLicensesFromDeadNodes(ctx context.Context, strftime interface{}) ([]License, error) {
	rows, err := q.db.QueryContext(ctx, releaseLicensesFromDeadNodes, strftime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []License
	for rows.Next() {
		var i License
		if err := rows.Scan(
			&i.ID,
			&i.File,
			&i.Key,
			&i.Claims,
			&i.LastClaimedAt,
			&i.LastReleasedAt,
			&i.NodeID,
			&i.CreatedAt,
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
