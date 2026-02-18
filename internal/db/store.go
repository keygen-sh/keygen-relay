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

func (s *Store) InsertLicense(ctx context.Context, pool *Pool, guid string, file []byte, key string, seats int64) (*License, error) {
	params := InsertLicenseParams{
		Guid:  guid,
		File:  file,
		Key:   key,
		Seats: seats,
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

func (s *Store) GetLicenseByNodeID(ctx context.Context, nodeID *int64, predicates ...LicensePredicateFunc) (*License, error) {
	predicate := applyLicensePredicates(predicates...)

	var license License
	var err error

	switch {
	case predicate.pool == AnyPool:
		return nil, ErrAnyPoolNotSupported
	case predicate.pool != nil:
		license, err = s.queries.GetLicenseByNodeIDWithPool(ctx, GetLicenseByNodeIDWithPoolParams{NodeID: *nodeID, PoolID: &predicate.pool.ID})
	default:
		license, err = s.queries.GetLicenseByNodeIDWithoutPool(ctx, *nodeID)
	}

	if err != nil {
		return nil, err
	}

	return &license, nil
}

func (s *Store) ReleaseLicenseByNodeID(ctx context.Context, nodeID *int64, predicates ...LicensePredicateFunc) error {
	predicate := applyLicensePredicates(predicates...)

	license, err := s.GetLicenseByNodeID(ctx, nodeID, predicates...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}

	switch {
	case predicate.pool == AnyPool:
		return ErrAnyPoolNotSupported
	case predicate.pool != nil:
		if err := s.queries.ReleaseLicenseNodeWithPoolByNodeID(ctx, ReleaseLicenseNodeWithPoolByNodeIDParams{NodeID: *nodeID, PoolID: &predicate.pool.ID}); err != nil {
			return err
		}
	default:
		if err := s.queries.ReleaseLicenseNodeWithoutPoolByNodeID(ctx, *nodeID); err != nil {
			return err
		}
	}

	return s.queries.TouchLicenseLastReleasedAt(ctx, license.ID)
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

	var licenseID int64
	var err error

	switch {
	case predicate.pool == AnyPool:
		return nil, ErrAnyPoolNotSupported
	case predicate.pool != nil:
		switch strategy {
		case "fifo":
			licenseID, err = s.queries.ClaimLicenseNodeWithPoolFIFO(ctx, ClaimLicenseNodeWithPoolFIFOParams{NodeID: *nodeID, PoolID: &predicate.pool.ID})
		case "lifo":
			licenseID, err = s.queries.ClaimLicenseNodeWithPoolLIFO(ctx, ClaimLicenseNodeWithPoolLIFOParams{NodeID: *nodeID, PoolID: &predicate.pool.ID})
		case "rand":
			licenseID, err = s.queries.ClaimLicenseNodeWithPoolRandom(ctx, ClaimLicenseNodeWithPoolRandomParams{NodeID: *nodeID, PoolID: &predicate.pool.ID})
		default:
			return nil, ErrBadStrategy
		}
	default:
		switch strategy {
		case "fifo":
			licenseID, err = s.queries.ClaimLicenseNodeWithoutPoolFIFO(ctx, *nodeID)
		case "lifo":
			licenseID, err = s.queries.ClaimLicenseNodeWithoutPoolLIFO(ctx, *nodeID)
		case "rand":
			licenseID, err = s.queries.ClaimLicenseNodeWithoutPoolRandom(ctx, *nodeID)
		default:
			return nil, ErrBadStrategy
		}
	}

	if err != nil {
		return nil, err
	}

	if err := s.queries.IncrementLicenseClaims(ctx, licenseID); err != nil {
		return nil, fmt.Errorf("failed to increment license claims: %w", err)
	}

	license, err := s.queries.GetLicenseByID(ctx, licenseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get license after claim: %w", err)
	}

	return &license, nil
}

func (s *Store) GetLicenseNodeCount(ctx context.Context, licenseID int64) (int64, error) {
	return s.queries.GetLicenseNodeCount(ctx, licenseID)
}

func (s *Store) ReleaseLicensesFromDeadNodes(ctx context.Context, ttl time.Duration) ([]License, error) {
	t := fmt.Sprintf("-%d seconds", int(ttl.Seconds()))

	licenseIDs, err := s.queries.ReleaseLicenseNodesFromDeadNodes(ctx, t)
	if err != nil {
		logger.Error("failed to release license nodes from dead nodes", "error", err)

		return nil, err
	}

	// deduplicate license IDs
	seen := make(map[int64]struct{}, len(licenseIDs))
	var licenses []License

	for _, id := range licenseIDs {
		if _, ok := seen[id]; ok {
			continue
		}

		seen[id] = struct{}{}

		if err := s.queries.TouchLicenseLastReleasedAt(ctx, id); err != nil {
			logger.Error("failed to touch license last_released_at", "licenseID", id, "error", err)

			return nil, err
		}

		license, err := s.queries.GetLicenseByID(ctx, id)
		if err != nil {
			logger.Error("failed to get license by ID after release", "licenseID", id, "error", err)

			return nil, err
		}

		licenses = append(licenses, license)
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
