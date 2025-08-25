package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/keygen-sh/keygen-relay/internal/logger"
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

type StoreError error

var (
	ErrBadStrategy StoreError = errors.New("invalid strategy")
)

type Store struct {
	queries    *Queries
	connection *sql.DB
}

// TxStore represents a Store within a transaction context
type TxStore struct {
	*Store
	tx *sql.Tx
}

func NewStore(queries *Queries, connection *sql.DB) *Store {
	return &Store{
		queries:    queries,
		connection: connection,
	}
}

// BeginTx begins a transaction and returns a TxStore that encapsulates transaction operations
func (s *Store) BeginTx(ctx context.Context) (*TxStore, error) {
	tx, err := s.connection.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &TxStore{
		Store: &Store{
			queries:    s.queries.WithTx(tx),
			connection: s.connection,
		},
		tx: tx,
	}, nil
}

// Commit commits the transaction
func (ts *TxStore) Commit() error {
	return ts.tx.Commit()
}

// Rollback rolls back the transaction
func (ts *TxStore) Rollback() error {
	return ts.tx.Rollback()
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

func (s *Store) GetLicenses(ctx context.Context, predicates ...LicensePredicateFunc) ([]License, error) {
	predicate := applyLicensePredicates(predicates...)

	switch {
	case predicate.pool == AnyPool:
		return s.queries.GetLicenses(ctx)
	case predicate.pool != nil:
		return s.queries.GetLicensesWithPool(ctx, &predicate.pool.ID)
	default:
		return s.queries.GetLicensesWithoutPool(ctx)
	}
}

func (s *Store) GetLicenseByGUID(ctx context.Context, id string, predicates ...LicensePredicateFunc) (*License, error) {
	predicate := applyLicensePredicates(predicates...)

	var license License
	var err error

	switch {
	case predicate.pool == AnyPool:
		license, err = s.queries.GetLicenseByGUID(ctx, id)
	case predicate.pool != nil:
		license, err = s.queries.GetLicenseWithPoolByGUID(ctx, GetLicenseWithPoolByGUIDParams{id, &predicate.pool.ID})
	default:
		license, err = s.queries.GetLicenseWithoutPoolByGUID(ctx, id)
	}

	if err != nil {
		return nil, err
	}

	return &license, nil
}

func (s *Store) ReleaseLicenseByNodeID(ctx context.Context, nodeID *int64, predicates ...LicensePredicateFunc) error {
	predicate := applyLicensePredicates(predicates...)

	switch {
	case predicate.pool == AnyPool:
		return ErrAnyPoolNotSupported
	case predicate.pool != nil:
		return s.queries.ReleaseLicenseWithPoolByNodeID(ctx, ReleaseLicenseWithPoolByNodeIDParams{nodeID, &predicate.pool.ID})
	default:
		return s.queries.ReleaseLicenseWithoutPoolByNodeID(ctx, nodeID)
	}
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

func (s *Store) GetPools(ctx context.Context) ([]Pool, error) {
	return s.queries.GetPools(ctx)
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
		return err
	}
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

		if err := tx.queries.InsertAuditLog(ctx, params); err != nil {
			return fmt.Errorf("failed to insert audit log: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *Store) ClaimLicenseByStrategy(ctx context.Context, strategy string, nodeID *int64, predicates ...LicensePredicateFunc) (*License, error) {
	predicate := applyLicensePredicates(predicates...)

	var license License
	var err error

	switch {
	case predicate.pool == AnyPool:
		return nil, ErrAnyPoolNotSupported
	case predicate.pool != nil:
		// Pool-specific strategy
		switch strategy {
		case "fifo":
			license, err = s.queries.ClaimLicenseWithPoolFIFO(ctx, ClaimLicenseWithPoolFIFOParams{nodeID, &predicate.pool.ID})
		case "lifo":
			license, err = s.queries.ClaimLicenseWithPoolLIFO(ctx, ClaimLicenseWithPoolLIFOParams{nodeID, &predicate.pool.ID})
		case "rand":
			license, err = s.queries.ClaimLicenseWithPoolRandom(ctx, ClaimLicenseWithPoolRandomParams{nodeID, &predicate.pool.ID})
		default:
			return nil, ErrBadStrategy
		}
	default:
		// No pool strategy
		switch strategy {
		case "fifo":
			license, err = s.queries.ClaimLicenseWithoutPoolFIFO(ctx, nodeID)
		case "lifo":
			license, err = s.queries.ClaimLicenseWithoutPoolLIFO(ctx, nodeID)
		case "rand":
			license, err = s.queries.ClaimLicenseWithoutPoolRandom(ctx, nodeID)
		default:
			return nil, ErrBadStrategy
		}
	}

	if err != nil {
		return nil, err
	}

	return &license, nil
}

func (s *Store) GetLicenseByNodeID(ctx context.Context, nodeID *int64, predicates ...LicensePredicateFunc) (*License, error) {
	predicate := applyLicensePredicates(predicates...)

	var license License
	var err error

	switch {
	case predicate.pool == AnyPool:
		return nil, ErrAnyPoolNotSupported
	case predicate.pool != nil:
		license, err = s.queries.GetLicenseWithPoolByNodeID(ctx, GetLicenseWithPoolByNodeIDParams{nodeID, &predicate.pool.ID})
	default:
		license, err = s.queries.GetLicenseWithoutPoolByNodeID(ctx, nodeID)
	}

	if err != nil {
		return nil, err
	}

	return &license, nil
}

func (s *Store) ReleaseLicensesFromDeadNodes(ctx context.Context, ttl time.Duration) ([]License, error) {
	t := fmt.Sprintf("-%d seconds", int(ttl.Seconds()))

	licenses, err := s.queries.ReleaseLicensesFromDeadNodes(ctx, t)
	if err != nil {
		logger.Error("failed to release licenses from dead nodes", "error", err)

		return nil, err
	}

	return licenses, nil
}

func (s *Store) DeactivateDeadNodes(ctx context.Context, ttl time.Duration) ([]Node, error) {
	t := fmt.Sprintf("-%d seconds", int(ttl.Seconds()))

	nodes, err := s.queries.DeactivateDeadNodes(ctx, t)
	if err != nil {
		logger.Error("failed to deactivate dead nodes", "error", err)

		return nil, err
	}

	return nodes, nil
}
