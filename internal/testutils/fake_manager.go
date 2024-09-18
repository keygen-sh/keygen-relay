package testutils

import (
	"context"
	"github.com/keygen-sh/keygen-relay/internal/licenses"
)

type FakeManager struct {
	AddLicenseFn     func(ctx context.Context, filePath, key, publicKey string) error
	RemoveLicenseFn  func(ctx context.Context, id string) error
	ListLicensesFn   func(ctx context.Context) ([]licenses.License, error)
	GetLicenseByIDFn func(ctx context.Context, id string) (licenses.License, error)
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
