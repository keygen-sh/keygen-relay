package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/keygen-sh/keygen-relay/internal/licenses"
	"log/slog"
	"time"
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
func (s *Store) WithTx(tx *sql.Tx) licenses.Store {
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

func (s *Store) DeleteLicenseByID(ctx context.Context, id string) error {
	return s.queries.DeleteLicenseByID(ctx, id)
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

func (s *Store) GetAllLicenses(ctx context.Context) ([]licenses.License, error) {
	dbLicenses, err := s.queries.GetAllLicenses(ctx)
	if err != nil {
		return nil, err
	}

	licensesList := make([]licenses.License, len(dbLicenses))
	for i, dbLic := range dbLicenses {
		licensesList[i] = convertToLicense(dbLic)
	}

	return licensesList, nil
}

func (s *Store) GetLicenseByID(ctx context.Context, id string) (licenses.License, error) {
	dbLicense, err := s.queries.GetLicenseByID(ctx, id)
	if err != nil {
		return licenses.License{}, err
	}

	return convertToLicense(dbLicense), nil
}

func (s *Store) ClaimLicense(ctx context.Context, params ClaimLicenseParams) error {
	return s.queries.ClaimLicense(ctx, params)
}

func (s *Store) ReleaseLicenseByNodeID(ctx context.Context, nodeID *int64) error {
	return s.queries.ReleaseLicenseByNodeID(ctx, nodeID)
}

func (s *Store) InsertNode(ctx context.Context, fingerprint string) (licenses.Node, error) {
	node, err := s.queries.InsertNode(ctx, fingerprint)

	if err != nil {
		return licenses.Node{}, err
	}

	return licenses.Node{
		ID:              node.ID,
		Fingerprint:     node.Fingerprint,
		ClaimedAt:       node.ClaimedAt,
		LastHeartbeatAt: node.LastHeartbeatAt,
		CreatedAt:       node.CreatedAt,
	}, nil
}

func (s *Store) UpdateNodeHeartbeat(ctx context.Context, fingerprint string) error {
	return s.queries.UpdateNodeHeartbeatByFingerprint(ctx, fingerprint)
}

func (s *Store) DeleteNodeByFingerprint(ctx context.Context, fingerprint string) error {
	return s.queries.DeleteNodeByFingerprint(ctx, fingerprint)
}

func (s *Store) GetNodeByFingerprint(ctx context.Context, fingerprint string) (licenses.Node, error) {
	node, err := s.queries.GetNodeByFingerprint(ctx, fingerprint)
	if err != nil {
		return licenses.Node{}, err
	}

	return licenses.Node{
		ID:              node.ID,
		Fingerprint:     node.Fingerprint,
		ClaimedAt:       node.ClaimedAt,
		LastHeartbeatAt: node.LastHeartbeatAt,
		CreatedAt:       node.CreatedAt,
	}, nil
}

func (s *Store) InsertAuditLog(ctx context.Context, action, entityType, entityID string) error {
	params := InsertAuditLogParams{
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
	}
	return s.queries.InsertAuditLog(ctx, params)
}

func (s *Store) ClaimUnclaimedLicenseFIFO(ctx context.Context, nodeID *int64) (licenses.License, error) {
	dbLicense, err := s.queries.ClaimUnclaimedLicenseFIFO(ctx, nodeID)

	if err != nil {
		return licenses.License{}, err
	}

	return convertToLicense(dbLicense), nil
}

func (s *Store) ClaimUnclaimedLicenseLIFO(ctx context.Context, nodeID *int64) (licenses.License, error) {
	dbLicense, err := s.queries.ClaimUnclaimedLicenseLIFO(ctx, nodeID)

	if err != nil {
		return licenses.License{}, err
	}

	return convertToLicense(dbLicense), nil
}

func (s *Store) ClaimUnclaimedLicenseRandom(ctx context.Context, nodeID *int64) (licenses.License, error) {
	dbLicense, err := s.queries.ClaimUnclaimedLicenseRandom(ctx, nodeID)

	if err != nil {
		return licenses.License{}, err
	}

	return convertToLicense(dbLicense), nil
}

func (s *Store) GetLicenseByNodeID(ctx context.Context, nodeID *int64) (licenses.License, error) {
	dbLicense, err := s.queries.GetLicenseByNodeID(ctx, nodeID)

	if err != nil {
		return licenses.License{}, err
	}

	return convertToLicense(dbLicense), nil
}

func (s *Store) UpdateNodeHeartbeatAndClaimedAtByFingerprint(ctx context.Context, fingerprint string) error {
	return s.queries.UpdateNodeHeartbeatAndClaimedAtByFingerprint(ctx, fingerprint)
}

func (s *Store) Heartbeat(ctx context.Context, fingerprint string) error {
	if err := s.queries.UpdateNodeHeartbeatByFingerprint(ctx, fingerprint); err != nil {
		slog.Error("failed to update node heartbeat", "fingerprint", fingerprint, "error", err)
		return fmt.Errorf("failed to update node heartbeat: %w", err)
	}
	slog.Info("heartbeat updated", "fingerprint", fingerprint)
	return nil
}

func (s *Store) DeleteInactiveNodes(ctx context.Context, ttl time.Duration) error {
	ttlDuration := fmt.Sprintf("-%d seconds", int(ttl.Seconds()))

	tx, err := s.connection.BeginTx(ctx, nil)
	if err != nil {
		slog.Error("failed to begin transaction", "error", err)
		return err
	}
	defer tx.Rollback()

	qtx := s.queries.WithTx(tx)

	if err := qtx.ReleaseLicensesFromInactiveNodes(ctx, ttlDuration); err != nil {
		slog.Error("failed to release licenses from inactive nodes", "error", err)
		return err
	}

	if err := qtx.DeleteInactiveNodes(ctx, ttlDuration); err != nil {
		slog.Error("failed to delete inactive nodes", "error", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		slog.Error("failed to commit transaction", "error", err)
		return err
	}

	slog.Debug("successfully released licenses and deleted inactive nodes")
	return nil
}

func convertToLicense(dbLic License) licenses.License {
	return licenses.License{
		ID:             dbLic.ID,
		File:           dbLic.File,
		Key:            dbLic.Key,
		Claims:         dbLic.Claims,
		NodeID:         dbLic.NodeID,
		CreatedAt:      dbLic.CreatedAt,
		LastClaimedAt:  dbLic.LastClaimedAt,
		LastReleasedAt: dbLic.LastReleasedAt,
	}
}
