package db

import (
	"context"
	"database/sql"
)

type Store struct {
	queries *Queries
}

func NewStore(queries *Queries) *Store {
	return &Store{
		queries: queries,
	}
}

func (s *Store) InsertLicense(ctx context.Context, params InsertLicenseParams) error {
	return s.queries.InsertLicense(ctx, params)
}

func (s *Store) DeleteLicenseByID(ctx context.Context, id string) error {
	return s.queries.DeleteLicenseByID(ctx, id)
}

func (s *Store) GetAllLicenses(ctx context.Context) ([]License, error) {
	return s.queries.GetAllLicenses(ctx)
}

func (s *Store) GetLicenseByID(ctx context.Context, id string) (License, error) {
	return s.queries.GetLicenseByID(ctx, id)
}

func (s *Store) ClaimLicense(ctx context.Context, params ClaimLicenseParams) error {
	return s.queries.ClaimLicense(ctx, params)
}

func (s *Store) ReleaseLicenseByNodeID(ctx context.Context, nodeID sql.NullInt64) error {
	return s.queries.ReleaseLicenseByNodeID(ctx, nodeID)
}

func (s *Store) InsertNode(ctx context.Context, fingerprint string) error {
	return s.queries.InsertNode(ctx, fingerprint)
}

func (s *Store) UpdateNodeHeartbeat(ctx context.Context, fingerprint string) error {
	return s.queries.UpdateNodeHeartbeatByFingerprint(ctx, fingerprint)
}

func (s *Store) DeleteNodeByFingerprint(ctx context.Context, fingerprint string) error {
	return s.queries.DeleteNodeByFingerprint(ctx, fingerprint)
}

func (s *Store) GetNodeByFingerprint(ctx context.Context, fingerprint string) (Node, error) {
	return s.queries.GetNodeByFingerprint(ctx, fingerprint)
}

func (s *Store) InsertAuditLog(ctx context.Context, params InsertAuditLogParams) error {
	return s.queries.InsertAuditLog(ctx, params)
}
