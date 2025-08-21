package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"
)

type EventTypeId int

const (
	EventTypeUnknown EventTypeId = iota
	EventTypeLicenseAdded
	EventTypeLicenseRemoved
	EventTypeLicenseLeased
	EventTypeLicenseLeaseExtended
	EventTypeLicenseReleased
	EventTypeLicenseLeaseExpired
	EventTypeNodeActivated
	EventTypeNodeHeartbeatPing
	EventTypeNodeDeactivated
	EventTypeNodeCulled
	EventTypePoolAdded
)

type EntityTypeId int

const (
	EntityTypeUnknown EntityTypeId = iota
	EntityTypeLicense
	EntityTypeNode
	EntityTypePool
)

type Store struct {
	queries    *Queries
	connection *sql.DB
}

func NewStore(queries *Queries, connection *sql.DB) *Store {
	return &Store{
		queries:    queries,
		connection: connection,
	}
}

func (s *Store) BeginTx(ctx context.Context) (*sql.Tx, error) {
	tx, err := s.connection.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return tx, nil
}

// WithTx returns a new Store that uses the transaction queries
func (s *Store) WithTx(tx *sql.Tx) *Store {
	return &Store{
		queries:    s.queries.WithTx(tx),
		connection: s.connection,
	}
}

func (s *Store) InsertLicense(ctx context.Context, pool *Pool, guid string, file []byte, key string) (*License, error) {
	params := InsertLicenseParams{
		Guid: guid,
		File: file,
		Key:  key,
	}

	if pool != nil {
		params.PoolID = &pool.ID
	}

	license, err := s.queries.InsertLicense(ctx, params)
	if err != nil {
		return nil, err
	}

	return &license, nil
}

func (s *Store) DeleteLicenseByGUID(ctx context.Context, id string) (*License, error) {
	license, err := s.queries.DeleteLicenseByGUID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete license: %w", err)
	}

	return &license, nil
}

func (s *Store) GetLicenses(ctx context.Context) ([]License, error) {
	licenses, err := s.queries.GetLicenses(ctx)
	if err != nil {
		return nil, err
	}

	return licenses, nil
}

func (s *Store) GetLicensesWithPool(ctx context.Context, pool *Pool) ([]License, error) {
	licenses, err := s.queries.GetLicensesWithPool(ctx, &pool.ID)
	if err != nil {
		return nil, err
	}

	return licenses, nil
}

func (s *Store) GetLicenseByGUID(ctx context.Context, id string) (*License, error) {
	license, err := s.queries.GetLicenseByGUID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &license, nil
}

func (s *Store) GetLicenseWithoutPoolByGUID(ctx context.Context, id string) (*License, error) {
	license, err := s.queries.GetLicenseWithoutPoolByGUID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &license, nil
}

func (s *Store) GetLicenseWithPoolByGUID(ctx context.Context, pool *Pool, id string) (*License, error) {
	license, err := s.queries.GetLicenseWithPoolByGUID(ctx, GetLicenseWithPoolByGUIDParams{id, &pool.ID})
	if err != nil {
		return nil, err
	}

	return &license, nil
}

func (s *Store) ReleaseLicenseWithPoolByNodeID(ctx context.Context, pool *Pool, nodeID *int64) error {
	return s.queries.ReleaseLicenseWithPoolByNodeID(ctx, ReleaseLicenseWithPoolByNodeIDParams{nodeID, &pool.ID})
}

func (s *Store) ReleaseLicenseWithoutPoolByNodeID(ctx context.Context, nodeID *int64) error {
	return s.queries.ReleaseLicenseWithoutPoolByNodeID(ctx, nodeID)
}

func (s *Store) ActivateNode(ctx context.Context, fingerprint string) (*Node, error) {
	node, err := s.queries.ActivateNode(ctx, fingerprint)
	if err != nil {
		return nil, err
	}

	return &node, nil
}

func (s *Store) DeactivateNodeByFingerprint(ctx context.Context, fingerprint string) error {
	return s.queries.DeactivateNodeByFingerprint(ctx, fingerprint)
}

func (s *Store) GetNodeByFingerprint(ctx context.Context, fingerprint string) (*Node, error) {
	node, err := s.queries.GetNodeByFingerprint(ctx, fingerprint)
	if err != nil {
		return nil, err
	}

	return &node, nil
}

func (s *Store) PingNodeHeartbeatByFingerprint(ctx context.Context, fingerprint string) error {
	return s.queries.PingNodeHeartbeatByFingerprint(ctx, fingerprint)
}

func (s *Store) CreatePool(ctx context.Context, name string) (*Pool, error) {
	pool, err := s.queries.CreatePool(ctx, name)
	if err != nil {
		return nil, err
	}

	return &pool, nil
}

func (s *Store) GetPoolByID(ctx context.Context, id int64) (*Pool, error) {
	pool, err := s.queries.GetPoolByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &pool, nil
}

func (s *Store) GetPoolByName(ctx context.Context, name string) (*Pool, error) {
	pool, err := s.queries.GetPoolByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return &pool, nil
}

func (s *Store) DeletePoolByID(ctx context.Context, id int64) (*Pool, error) {
	pool, err := s.queries.DeletePoolByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &pool, nil
}

// TODO(ezekg) allow event data? e.g. license.lease_extended {from:x,to:y} or license.leased {node:n} or node.heartbeat_ping {count:n}
//
//	but doing so would pose problems for future aggregation...
func (s *Store) InsertAuditLog(ctx context.Context, pool *Pool, eventTypeId EventTypeId, entityTypeId EntityTypeId, entityID int64) error {
	params := InsertAuditLogParams{
		EventTypeID:  int64(eventTypeId),
		EntityTypeID: int64(entityTypeId),
		EntityID:     entityID,
	}

	if pool != nil {
		params.PoolID = &pool.ID
	}

	return s.queries.InsertAuditLog(ctx, params)
}

type BulkInsertAuditLogParams struct {
	EventTypeID  EventTypeId
	EntityTypeID EntityTypeId
	EntityID     int64
	Pool         *Pool
}

func (s *Store) BulkInsertAuditLogs(ctx context.Context, logs []BulkInsertAuditLogParams) error {
	tx, err := s.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	qtx := s.WithTx(tx)
	defer tx.Rollback()

	for _, log := range logs {
		params := InsertAuditLogParams{
			EventTypeID:  int64(log.EventTypeID),
			EntityTypeID: int64(log.EntityTypeID),
			EntityID:     log.EntityID,
		}

		if log.Pool != nil {
			params.PoolID = &log.Pool.ID
		}

		if err := qtx.queries.InsertAuditLog(ctx, params); err != nil {
			return fmt.Errorf("failed to insert audit log: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *Store) ClaimLicenseWithoutPoolFIFO(ctx context.Context, nodeID *int64) (*License, error) {
	license, err := s.queries.ClaimLicenseWithoutPoolFIFO(ctx, nodeID)
	if err != nil {
		return nil, err
	}

	return &license, nil
}

func (s *Store) ClaimLicenseWithoutPoolLIFO(ctx context.Context, nodeID *int64) (*License, error) {
	license, err := s.queries.ClaimLicenseWithoutPoolLIFO(ctx, nodeID)
	if err != nil {
		return nil, err
	}

	return &license, nil
}

func (s *Store) ClaimLicenseWithoutPoolRandom(ctx context.Context, nodeID *int64) (*License, error) {
	license, err := s.queries.ClaimLicenseWithoutPoolRandom(ctx, nodeID)
	if err != nil {
		return nil, err
	}

	return &license, nil
}

func (s *Store) ClaimLicenseWithoutPoolByStrategy(ctx context.Context, strategy string, nodeID *int64) (*License, error) {
	switch strategy {
	case "fifo":
		return s.ClaimLicenseWithoutPoolFIFO(ctx, nodeID)
	case "lifo":
		return s.ClaimLicenseWithoutPoolLIFO(ctx, nodeID)
	case "rand":
		return s.ClaimLicenseWithoutPoolRandom(ctx, nodeID)
	default:
		return s.ClaimLicenseWithoutPoolFIFO(ctx, nodeID)
	}
}

func (s *Store) ClaimLicenseWithPoolFIFO(ctx context.Context, pool *Pool, nodeID *int64) (*License, error) {
	license, err := s.queries.ClaimLicenseWithPoolFIFO(ctx, ClaimLicenseWithPoolFIFOParams{nodeID, &pool.ID})
	if err != nil {
		return nil, err
	}

	return &license, nil
}

func (s *Store) ClaimLicenseWithPoolLIFO(ctx context.Context, pool *Pool, nodeID *int64) (*License, error) {
	license, err := s.queries.ClaimLicenseWithPoolLIFO(ctx, ClaimLicenseWithPoolLIFOParams{nodeID, &pool.ID})
	if err != nil {
		return nil, err
	}

	return &license, nil
}

func (s *Store) ClaimLicenseWithPoolRandom(ctx context.Context, pool *Pool, nodeID *int64) (*License, error) {
	license, err := s.queries.ClaimLicenseWithPoolRandom(ctx, ClaimLicenseWithPoolRandomParams{nodeID, &pool.ID})
	if err != nil {
		return nil, err
	}

	return &license, nil
}

func (s *Store) ClaimLicenseWithPoolByStrategy(ctx context.Context, pool *Pool, strategy string, nodeID *int64) (*License, error) {
	switch strategy {
	case "fifo":
		return s.ClaimLicenseWithPoolFIFO(ctx, pool, nodeID)
	case "lifo":
		return s.ClaimLicenseWithPoolLIFO(ctx, pool, nodeID)
	case "rand":
		return s.ClaimLicenseWithPoolRandom(ctx, pool, nodeID)
	default:
		return s.ClaimLicenseWithPoolFIFO(ctx, pool, nodeID)
	}
}

func (s *Store) GetLicenseWithoutPoolByNodeID(ctx context.Context, nodeID *int64) (*License, error) {
	license, err := s.queries.GetLicenseWithoutPoolByNodeID(ctx, nodeID)
	if err != nil {
		return nil, err
	}

	return &license, nil
}

func (s *Store) GetLicenseWithPoolByNodeID(ctx context.Context, pool *Pool, nodeID *int64) (*License, error) {
	license, err := s.queries.GetLicenseWithPoolByNodeID(ctx, GetLicenseWithPoolByNodeIDParams{nodeID, &pool.ID})
	if err != nil {
		return nil, err
	}

	return &license, nil
}

func (s *Store) ReleaseLicensesFromDeadNodes(ctx context.Context, ttl time.Duration) ([]License, error) {
	t := fmt.Sprintf("-%d seconds", int(ttl.Seconds()))

	licenses, err := s.queries.ReleaseLicensesFromDeadNodes(ctx, t)
	if err != nil {
		slog.Error("failed to release licenses from dead nodes", "error", err)

		return nil, err
	}

	return licenses, nil
}

func (s *Store) DeactivateDeadNodes(ctx context.Context, ttl time.Duration) ([]Node, error) {
	t := fmt.Sprintf("-%d seconds", int(ttl.Seconds()))

	nodes, err := s.queries.DeactivateDeadNodes(ctx, t)
	if err != nil {
		slog.Error("failed to deactivate dead nodes", "error", err)

		return nil, err
	}

	return nodes, nil
}
