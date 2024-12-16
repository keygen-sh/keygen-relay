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
)

type EntityTypeId int

const (
	EntityTypeUnknown EntityTypeId = iota
	EntityTypeLicense
	EntityTypeNode
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

func (s *Store) InsertLicense(ctx context.Context, id string, file []byte, key string) error {
	params := InsertLicenseParams{
		ID:   id,
		File: file,
		Key:  key,
	}
	return s.queries.InsertLicense(ctx, params)
}

func (s *Store) DeleteLicenseByIDTx(ctx context.Context, id string) error {
	tx, err := s.connection.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	qtx := s.queries.WithTx(tx)
	defer tx.Rollback()

	_, err = qtx.GetLicenseByID(ctx, id)
	if err != nil {
		return err
	}

	err = qtx.DeleteLicenseByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete license: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *Store) GetAllLicenses(ctx context.Context) ([]License, error) {
	licenses, err := s.queries.GetAllLicenses(ctx)
	if err != nil {
		return nil, err
	}

	return licenses, nil
}

func (s *Store) GetLicenseByID(ctx context.Context, id string) (*License, error) {
	license, err := s.queries.GetLicenseByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &license, nil
}

func (s *Store) ReleaseLicenseByNodeID(ctx context.Context, nodeID *int64) error {
	return s.queries.ReleaseLicenseByNodeID(ctx, nodeID)
}

func (s *Store) InsertNode(ctx context.Context, fingerprint string) (*Node, error) {
	node, err := s.queries.InsertNode(ctx, fingerprint)
	if err != nil {
		return nil, err
	}

	return &node, nil
}

func (s *Store) DeleteNodeByFingerprint(ctx context.Context, fingerprint string) error {
	return s.queries.DeleteNodeByFingerprint(ctx, fingerprint)
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

// TODO(ezekg) allow event data? e.g. license.lease_extended {from:x,to:y} or license.leased {node:n} or node.heartbeat_ping {count:n}
//
//	but doing so would pose problems for future aggregation...
func (s *Store) InsertAuditLog(ctx context.Context, eventTypeId EventTypeId, entityTypeId EntityTypeId, entityID string) error {
	params := InsertAuditLogParams{
		EventTypeID:  int64(eventTypeId),
		EntityTypeID: int64(entityTypeId),
		EntityID:     entityID,
	}

	return s.queries.InsertAuditLog(ctx, params)
}

type BulkInsertAuditLogParams struct {
	EventTypeID  EventTypeId
	EntityTypeID EntityTypeId
	EntityID     string
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

		if err := qtx.queries.InsertAuditLog(ctx, params); err != nil {
			return fmt.Errorf("failed to insert audit log: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *Store) ClaimLicenseFIFO(ctx context.Context, nodeID *int64) (*License, error) {
	license, err := s.queries.ClaimLicenseFIFO(ctx, nodeID)
	if err != nil {
		return nil, err
	}

	return &license, nil
}

func (s *Store) ClaimLicenseLIFO(ctx context.Context, nodeID *int64) (*License, error) {
	license, err := s.queries.ClaimLicenseLIFO(ctx, nodeID)
	if err != nil {
		return nil, err
	}

	return &license, nil
}

func (s *Store) ClaimLicenseRandom(ctx context.Context, nodeID *int64) (*License, error) {
	license, err := s.queries.ClaimLicenseRandom(ctx, nodeID)
	if err != nil {
		return nil, err
	}

	return &license, nil
}

func (s *Store) ClaimLicenseByStrategy(ctx context.Context, strategy string, nodeID *int64) (*License, error) {
	switch strategy {
	case "fifo":
		return s.ClaimLicenseFIFO(ctx, nodeID)
	case "lifo":
		return s.ClaimLicenseLIFO(ctx, nodeID)
	case "rand":
		return s.ClaimLicenseRandom(ctx, nodeID)
	default:
		return s.ClaimLicenseFIFO(ctx, nodeID)
	}
}

func (s *Store) GetLicenseByNodeID(ctx context.Context, nodeID *int64) (*License, error) {
	license, err := s.queries.GetLicenseByNodeID(ctx, nodeID)
	if err != nil {
		return nil, err
	}

	return &license, nil
}

func (s *Store) ReleaseLicensesFromInactiveNodes(ctx context.Context, ttl time.Duration) ([]License, error) {
	t := fmt.Sprintf("-%d seconds", int(ttl.Seconds()))

	licenses, err := s.queries.ReleaseLicensesFromInactiveNodes(ctx, t)
	if err != nil {
		slog.Error("failed to release licenses from inactive nodes", "error", err)

		return nil, err
	}

	return licenses, nil
}

func (s *Store) DeleteInactiveNodes(ctx context.Context, ttl time.Duration) ([]Node, error) {
	t := fmt.Sprintf("-%d seconds", int(ttl.Seconds()))

	nodes, err := s.queries.DeleteInactiveNodes(ctx, t)
	if err != nil {
		slog.Error("failed to delete inactive nodes", "error", err)

		return nil, err
	}

	return nodes, nil
}
