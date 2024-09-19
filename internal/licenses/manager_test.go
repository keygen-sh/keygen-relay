package licenses_test

import (
	"context"
	"github.com/keygen-sh/keygen-relay/internal/licenses"
	"github.com/keygen-sh/keygen-relay/internal/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddLicense_Success(t *testing.T) {
	store, dbConn := testutils.NewMemoryStore(t)
	defer testutils.CloseMemoryStore(dbConn)

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

	manager.AttachStore(store)

	err := manager.AddLicense(context.Background(), "test_license.lic", "test_key", "test_public_key")
	assert.NoError(t, err)

	license, err := manager.GetLicenseByID(context.Background(), "license_test_key")
	assert.NoError(t, err)
	assert.Equal(t, "test_key", license.Key)
}

func TestAddLicense_Failure(t *testing.T) {
	store, dbConn := testutils.NewMemoryStore(t)
	defer testutils.CloseMemoryStore(dbConn)

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

	manager.AttachStore(store)

	err := manager.AddLicense(context.Background(), "test_license.lic", "test_key", "test_public_key")
	assert.NoError(t, err)

	err = manager.AddLicense(context.Background(), "test_license.lic", "test_key", "test_public_key")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to insert license")
}

func TestRemoveLicense_Success(t *testing.T) {
	ctx := context.Background()
	store, dbConn := testutils.NewMemoryStore(t)
	defer testutils.CloseMemoryStore(dbConn)

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

	manager.AttachStore(store)

	// Add a license that to be deleted
	err := manager.AddLicense(ctx, "test_license.lic", "test_key", "test_public_key")
	assert.NoError(t, err, "Failed to add license")

	// Check that the license was created
	_, err = manager.GetLicenseByID(ctx, "license_test_key")
	assert.NoError(t, err, "License should exist")

	// Remove the license
	err = manager.RemoveLicense(ctx, "license_test_key")
	assert.NoError(t, err, "Failed to remove license")

	// Ensure the license is removed
	_, err = manager.GetLicenseByID(context.Background(), "license_test_key")
	assert.Error(t, err, "License should not exist after deletion")
	assert.Contains(t, err.Error(), "license with ID license_test_key not found")
}

func TestRemoveLicense_Failure(t *testing.T) {
	store, dbConn := testutils.NewMemoryStore(t)
	defer testutils.CloseMemoryStore(dbConn)

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

	manager.AttachStore(store)

	err := manager.RemoveLicense(context.Background(), "invalid_id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "license with ID invalid_id not found")
}

func TestListLicenses_Success(t *testing.T) {
	store, dbConn := testutils.NewMemoryStore(t)
	defer testutils.CloseMemoryStore(dbConn)

	manager := licenses.NewManager(
		&licenses.Config{
			EnabledAudit: true,
		},
		func(filename string) ([]byte, error) {
			if filename == "test_license_1.lic" {
				return []byte("mock_certificate_1"), nil
			}
			return []byte("mock_certificate_2"), nil
		},
		func(cert []byte) licenses.LicenseVerifier {
			return &testutils.FakeLicenseVerifier{}
		},
	)

	manager.AttachStore(store)

	err := manager.AddLicense(context.Background(), "test_license_1.lic", "test_key_1", "test_public_key_1")
	assert.NoError(t, err)
	err = manager.AddLicense(context.Background(), "test_license_2.lic", "test_key_2", "test_public_key_2")
	assert.NoError(t, err)

	licenseList, err := manager.ListLicenses(context.Background())
	assert.NoError(t, err)
	assert.Len(t, licenseList, 2)
	assert.Equal(t, "test_key_1", licenseList[0].Key)
	assert.Equal(t, "test_key_2", licenseList[1].Key)
}

func TestGetLicenseByID_Success(t *testing.T) {
	store, dbConn := testutils.NewMemoryStore(t)
	defer testutils.CloseMemoryStore(dbConn)

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

	manager.AttachStore(store)

	err := manager.AddLicense(context.Background(), "test_license.lic", "test_key", "test_public_key")
	assert.NoError(t, err)

	license, err := manager.GetLicenseByID(context.Background(), "license_test_key")
	assert.NoError(t, err)
	assert.Equal(t, "test_key", license.Key)
}

func TestGetLicenseByID_Failure(t *testing.T) {
	store, dbConn := testutils.NewMemoryStore(t)
	defer testutils.CloseMemoryStore(dbConn)

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

	manager.AttachStore(store)

	_, err := manager.GetLicenseByID(context.Background(), "invalid_id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "license with ID invalid_id not found")
}
