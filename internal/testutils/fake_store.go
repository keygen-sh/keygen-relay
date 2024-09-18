package testutils

import (
	"context"
	"github.com/keygen-sh/keygen-relay/internal/licenses"
)

type FakeStore struct {
	InsertLicenseFn           func(ctx context.Context, id string, file []byte, key string) error
	DeleteLicenseByIDFn       func(ctx context.Context, id string) error
	GetAllLicensesFn          func(ctx context.Context) ([]licenses.License, error)
	GetLicenseByIDFn          func(ctx context.Context, id string) (licenses.License, error)
	InsertAuditLogFn          func(ctx context.Context, action, entityType, entityID string) error
	InsertNodeFn              func(ctx context.Context, fingerprint string) error
	GetNodeByFingerprintFn    func(ctx context.Context, fingerprint string) (licenses.Node, error)
	UpdateNodeHeartbeatFn     func(ctx context.Context, fingerprint string) error
	DeleteNodeByFingerprintFn func(ctx context.Context, fingerprint string) error
}

func (m *FakeStore) InsertLicense(ctx context.Context, id string, file []byte, key string) error {
	return m.InsertLicenseFn(ctx, id, file, key)
}

func (m *FakeStore) DeleteLicenseByID(ctx context.Context, id string) error {
	return m.DeleteLicenseByIDFn(ctx, id)
}

func (m *FakeStore) GetAllLicenses(ctx context.Context) ([]licenses.License, error) {
	return m.GetAllLicensesFn(ctx)
}

func (m *FakeStore) GetLicenseByID(ctx context.Context, id string) (licenses.License, error) {
	return m.GetLicenseByIDFn(ctx, id)
}

func (m *FakeStore) InsertAuditLog(ctx context.Context, action, entityType, entityID string) error {
	return m.InsertAuditLogFn(ctx, action, entityType, entityID)
}

func (m *FakeStore) InsertNode(ctx context.Context, fingerprint string) error {
	return m.InsertNodeFn(ctx, fingerprint)
}

func (m *FakeStore) GetNodeByFingerprint(ctx context.Context, fingerprint string) (licenses.Node, error) {
	return m.GetNodeByFingerprintFn(ctx, fingerprint)
}

func (m *FakeStore) UpdateNodeHeartbeat(ctx context.Context, fingerprint string) error {
	return m.UpdateNodeHeartbeatFn(ctx, fingerprint)
}

func (m *FakeStore) DeleteNodeByFingerprint(ctx context.Context, fingerprint string) error {
	return m.DeleteNodeByFingerprintFn(ctx, fingerprint)
}
