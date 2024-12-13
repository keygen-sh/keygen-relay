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
	EventTypeLicenseClaimed
	EventTypeLicenseReleased
	EventTypeNodeActivated
	EventTypeNodePing
	EventTypeNodeCulled
	EventTypeNodeDeactivated
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
	defer tx.Rollback()

	qtx := s.queries.WithTx(tx)

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

func (s *Store) UpdateNodeHeartbeat(ctx context.Context, fingerprint string) error {
	return s.queries.UpdateNodeHeartbeatByFingerprint(ctx, fingerprint)
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

func (s *Store) InsertAuditLog(ctx context.Context, eventTypeId EventTypeId, entityType string, entityID string) error {
	params := InsertAuditLogParams{
		EventTypeID: int64(eventTypeId),
		EntityType:  entityType,
		EntityID:    entityID,
	}
	return s.queries.InsertAuditLog(ctx, params)
}

func (s *Store) ClaimUnclaimedLicenseFIFO(ctx context.Context, nodeID *int64) (*License, error) {
	license, err := s.queries.ClaimUnclaimedLicenseFIFO(ctx, nodeID)
	if err != nil {
		return nil, err
	}

	return &license, nil
}

func (s *Store) ClaimUnclaimedLicenseLIFO(ctx context.Context, nodeID *int64) (*License, error) {
	license, err := s.queries.ClaimUnclaimedLicenseLIFO(ctx, nodeID)
	if err != nil {
		return nil, err
	}

	return &license, nil
}

func (s *Store) ClaimUnclaimedLicenseRandom(ctx context.Context, nodeID *int64) (*License, error) {
	license, err := s.queries.ClaimUnclaimedLicenseRandom(ctx, nodeID)
	if err != nil {
		return nil, err
	}

	return &license, nil
}

func (s *Store) GetLicenseByNodeID(ctx context.Context, nodeID *int64) (*License, error) {
	license, err := s.queries.GetLicenseByNodeID(ctx, nodeID)
	if err != nil {
		return nil, err
	}

	return &license, nil
}

func (s *Store) UpdateNodeHeartbeatAndClaimedAtByFingerprint(ctx context.Context, fingerprint string) error {
	return s.queries.UpdateNodeHeartbeatAndClaimedAtByFingerprint(ctx, fingerprint)
}

func (s *Store) ReleaseLicensesFromInactiveNodes(ctx context.Context, ttl time.Duration) ([]License, error) {
	ttlDuration := fmt.Sprintf("-%d seconds", int(ttl.Seconds()))
	licenses, err := s.queries.ReleaseLicensesFromInactiveNodes(ctx, ttlDuration)
	if err != nil {
		slog.Error("failed to release licenses from inactive nodes", "error", err)

		return nil, err
	}

	return licenses, nil
}

func (s *Store) DeleteInactiveNodes(ctx context.Context, ttl time.Duration) error {
	ttlDuration := fmt.Sprintf("-%d seconds", int(ttl.Seconds()))

	if err := s.queries.DeleteInactiveNodes(ctx, ttlDuration); err != nil {
		slog.Error("failed to delete inactive nodes", "error", err)
		return err
	}

	return nil
}
