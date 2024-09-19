package licenses_test

import (
	"context"
	"errors"
	"github.com/keygen-sh/keygen-relay/internal/licenses"
	"github.com/keygen-sh/keygen-relay/internal/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddLicense_Success(t *testing.T) {
	fakeStore := &testutils.FakeStore{
		InsertLicenseFn: func(ctx context.Context, id string, file []byte, key string) error {
			return nil
		},
		InsertAuditLogFn: func(ctx context.Context, action, entityType, entityID string) error {
			return nil
		},
	}

	manager := licenses.NewManager(
		&licenses.Config{
			EnabledAudit: true,
		},
		func(filename string) ([]byte, error) {
			return []byte("mock_certificate"), nil
		},
		func(cert []byte) licenses.LicenseVerifier {
			return &testutils.FakeLicenseVerifier{}
		},
	)

	manager.AttachStore(fakeStore)

	err := manager.AddLicense(context.Background(), "test_license.lic", "test_key", "test_public_key")
	assert.NoError(t, err)
}

func TestAddLicense_Failure(t *testing.T) {
	fakeStore := &testutils.FakeStore{
		InsertLicenseFn: func(ctx context.Context, id string, file []byte, key string) error {
			return errors.New("failed to insert license")
		},
	}

	manager := licenses.NewManager(
		&licenses.Config{
			EnabledAudit: true,
		},
		func(filename string) ([]byte, error) {
			return []byte("mock_certificate"), nil
		},
		func(cert []byte) licenses.LicenseVerifier {
			return &testutils.FakeLicenseVerifier{}
		},
	)

	manager.AttachStore(fakeStore)

	err := manager.AddLicense(context.Background(), "test_license.lic", "test_key", "test_public_key")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to insert license")
}

func TestRemoveLicense_Success(t *testing.T) {
	fakeStore := &testutils.FakeStore{
		DeleteLicenseByIDFn: func(ctx context.Context, id string) error {
			return nil
		},
		InsertAuditLogFn: func(ctx context.Context, action, entityType, entityID string) error {
			return nil
		},
	}

	manager := licenses.NewManager(
		&licenses.Config{
			EnabledAudit: true,
		},
		func(filename string) ([]byte, error) {
			return []byte("mock_certificate"), nil
		},
		func(cert []byte) licenses.LicenseVerifier {
			return &testutils.FakeLicenseVerifier{}
		},
	)

	manager.AttachStore(fakeStore)

	err := manager.RemoveLicense(context.Background(), "test_id")
	assert.NoError(t, err)
}

func TestRemoveLicense_Failure(t *testing.T) {
	fakeStore := &testutils.FakeStore{
		DeleteLicenseByIDFn: func(ctx context.Context, id string) error {
			return errors.New("failed to remove license")
		},
	}

	manager := licenses.NewManager(
		&licenses.Config{
			EnabledAudit: true,
		},
		func(filename string) ([]byte, error) {
			return []byte("mock_certificate"), nil
		},
		func(cert []byte) licenses.LicenseVerifier {
			return &testutils.FakeLicenseVerifier{}
		},
	)

	manager.AttachStore(fakeStore)

	err := manager.RemoveLicense(context.Background(), "test_id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to remove license")
}

func TestListLicenses_Success(t *testing.T) {
	fakeStore := &testutils.FakeStore{
		GetAllLicensesFn: func(ctx context.Context) ([]licenses.License, error) {
			return []licenses.License{
				{
					ID:     "test_id_1",
					Key:    "test_key_1",
					Claims: 1,
				},
				{
					ID:     "test_id_2",
					Key:    "test_key_2",
					Claims: 2,
				},
			}, nil
		},
	}

	manager := licenses.NewManager(
		&licenses.Config{
			EnabledAudit: true,
		},
		func(filename string) ([]byte, error) {
			return []byte("mock_certificate"), nil
		},
		func(cert []byte) licenses.LicenseVerifier {
			return &testutils.FakeLicenseVerifier{}
		},
	)

	manager.AttachStore(fakeStore)

	licenseList, err := manager.ListLicenses(context.Background())
	assert.NoError(t, err)
	assert.Len(t, licenseList, 2)
	assert.Equal(t, "test_id_1", licenseList[0].ID)
	assert.Equal(t, "test_key_1", licenseList[0].Key)
}

func TestGetLicenseByID_Success(t *testing.T) {
	fakeStore := &testutils.FakeStore{
		GetLicenseByIDFn: func(ctx context.Context, id string) (licenses.License, error) {
			return licenses.License{
				ID:     "test_id",
				Key:    "test_key",
				Claims: 1,
			}, nil
		},
	}

	manager := licenses.NewManager(
		&licenses.Config{
			EnabledAudit: true,
		},
		func(filename string) ([]byte, error) {
			return []byte("mock_certificate"), nil
		},
		func(cert []byte) licenses.LicenseVerifier {
			return &testutils.FakeLicenseVerifier{}
		},
	)

	manager.AttachStore(fakeStore)

	license, err := manager.GetLicenseByID(context.Background(), "test_id")
	assert.NoError(t, err)
	assert.Equal(t, "test_id", license.ID)
	assert.Equal(t, "test_key", license.Key)
}

func TestGetLicenseByID_Failure(t *testing.T) {
	fakeStore := &testutils.FakeStore{
		GetLicenseByIDFn: func(ctx context.Context, id string) (licenses.License, error) {
			return licenses.License{}, errors.New("license not found")
		},
	}

	manager := licenses.NewManager(
		&licenses.Config{
			EnabledAudit: true,
		},
		func(filename string) ([]byte, error) {
			return []byte("mock_certificate"), nil
		},
		func(cert []byte) licenses.LicenseVerifier {
			return &testutils.FakeLicenseVerifier{}
		},
	)

	manager.AttachStore(fakeStore)

	_, err := manager.GetLicenseByID(context.Background(), "invalid_id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "license not found")
}
