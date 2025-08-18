package testutils

import (
	"context"
	"time"

	"github.com/keygen-sh/keygen-relay/internal/db"
	"github.com/keygen-sh/keygen-relay/internal/licenses"
)

type FakeManager struct {
	store              db.Store
	AddLicenseFn       func(ctx context.Context, pool *string, filePath, key, publicKey string) (*db.License, error)
	RemoveLicenseFn    func(ctx context.Context, pool *string, id string) error
	ListLicensesFn     func(ctx context.Context, pool *string) ([]db.License, error)
	GetLicenseByGUIDFn func(ctx context.Context, pool *string, id string) (*db.License, error)
	ClaimLicenseFn     func(ctx context.Context, pool *string, fingerprint string) (*licenses.LicenseOperationResult, error)
	ReleaseLicenseFn   func(ctx context.Context, pool *string, fingerprint string) (*licenses.LicenseOperationResult, error)
	CullDeadNodesFn    func(ctx context.Context, ttl time.Duration) ([]db.Node, error)
	ConfigFn           func() *licenses.Config
}

func (f *FakeManager) AddLicense(ctx context.Context, pool *string, filePath, key, publicKey string) (*db.License, error) {
	if f.AddLicenseFn != nil {
		return f.AddLicenseFn(ctx, pool, filePath, key, publicKey)
	}
	return &db.License{}, nil
}

func (f *FakeManager) RemoveLicense(ctx context.Context, pool *string, id string) error {
	if f.RemoveLicenseFn != nil {
		return f.RemoveLicenseFn(ctx, pool, id)
	}
	return nil
}

func (f *FakeManager) ListLicenses(ctx context.Context, pool *string) ([]db.License, error) {
	if f.ListLicensesFn != nil {
		return f.ListLicensesFn(ctx, pool)
	}
	return []db.License{}, nil
}

func (f *FakeManager) GetLicenseByGUID(ctx context.Context, pool *string, id string) (*db.License, error) {
	if f.GetLicenseByGUIDFn != nil {
		return f.GetLicenseByGUIDFn(ctx, pool, id)
	}
	return &db.License{}, nil
}

func (f *FakeManager) AttachStore(store db.Store) {
	f.store = store
}

func (f *FakeManager) ClaimLicense(ctx context.Context, pool *string, fingerprint string) (*licenses.LicenseOperationResult, error) {
	if f.ClaimLicenseFn != nil {
		return f.ClaimLicenseFn(ctx, pool, fingerprint)
	}

	return nil, nil
}

func (f *FakeManager) ReleaseLicense(ctx context.Context, pool *string, fingerprint string) (*licenses.LicenseOperationResult, error) {
	if f.ReleaseLicenseFn != nil {
		return f.ReleaseLicenseFn(ctx, pool, fingerprint)
	}

	return nil, nil
}

func (f *FakeManager) CullDeadNodes(ctx context.Context, ttl time.Duration) ([]db.Node, error) {
	if f.CullDeadNodesFn != nil {
		return f.CullDeadNodesFn(ctx, ttl)
	}

	return nil, nil
}

func (f *FakeManager) Config() *licenses.Config {
	if f.ConfigFn != nil {
		return f.ConfigFn()
	}

	return &licenses.Config{}
}
