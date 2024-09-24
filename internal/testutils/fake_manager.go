package testutils

import (
	"context"
	"github.com/keygen-sh/keygen-relay/internal/licenses"
	"time"
)

type FakeManager struct {
	store                  licenses.Store
	AddLicenseFn           func(ctx context.Context, filePath, key, publicKey string) error
	RemoveLicenseFn        func(ctx context.Context, id string) error
	ListLicensesFn         func(ctx context.Context) ([]licenses.License, error)
	GetLicenseByIDFn       func(ctx context.Context, id string) (licenses.License, error)
	ClaimLicenseFn         func(ctx context.Context, fingerprint string) (*licenses.LicenseOperationResult, error)
	ReleaseLicenseFn       func(ctx context.Context, fingerprint string) (*licenses.LicenseOperationResult, error)
	CleanupInactiveNodesFn func(ctx context.Context, ttl time.Duration) error
	ConfigFn               func() *licenses.Config
}

func (f *FakeManager) AddLicense(ctx context.Context, filePath, key, publicKey string) error {
	if f.AddLicenseFn != nil {
		return f.AddLicenseFn(ctx, filePath, key, publicKey)
	}
	return nil
}

func (f *FakeManager) RemoveLicense(ctx context.Context, id string) error {
	if f.RemoveLicenseFn != nil {
		return f.RemoveLicenseFn(ctx, id)
	}
	return nil
}

func (f *FakeManager) ListLicenses(ctx context.Context) ([]licenses.License, error) {
	if f.ListLicensesFn != nil {
		return f.ListLicensesFn(ctx)
	}
	return []licenses.License{}, nil
}

func (f *FakeManager) GetLicenseByID(ctx context.Context, id string) (licenses.License, error) {
	if f.GetLicenseByIDFn != nil {
		return f.GetLicenseByIDFn(ctx, id)
	}
	return licenses.License{}, nil
}

func (f *FakeManager) AttachStore(store licenses.Store) {
	f.store = store
}

func (f *FakeManager) ClaimLicense(ctx context.Context, fingerprint string) (*licenses.LicenseOperationResult, error) {
	if f.ClaimLicenseFn != nil {
		return f.ClaimLicenseFn(ctx, fingerprint)
	}

	return nil, nil
}

func (f *FakeManager) ReleaseLicense(ctx context.Context, fingerprint string) (*licenses.LicenseOperationResult, error) {
	if f.ReleaseLicenseFn != nil {
		return f.ReleaseLicenseFn(ctx, fingerprint)
	}

	return nil, nil
}

func (f *FakeManager) CleanupInactiveNodes(ctx context.Context, ttl time.Duration) error {
	if f.CleanupInactiveNodesFn != nil {
		return f.CleanupInactiveNodesFn(ctx, ttl)
	}

	return nil
}

func (f *FakeManager) Config() *licenses.Config {
	if f.ConfigFn != nil {
		return f.ConfigFn()
	}

	return &licenses.Config{}
}
