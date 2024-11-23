package licenses_test

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/keygen-sh/keygen-relay/internal/licenses"
	"github.com/keygen-sh/keygen-relay/internal/testutils"
	"github.com/stretchr/testify/assert"
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
	assert.Contains(t, err.Error(), "license with the provided key already exists")
}

func TestAddLicense_FileNotFound(t *testing.T) {
	store, dbConn := testutils.NewMemoryStore(t)
	defer testutils.CloseMemoryStore(dbConn)

	manager := licenses.NewManager(
		&licenses.Config{
			EnabledAudit: true,
		},
		func(filename string) ([]byte, error) {
			return nil, os.ErrNotExist
		},
		func(cert []byte) licenses.LicenseVerifier {
			return &testutils.FakeLicenseVerifier{}
		},
	)

	manager.AttachStore(store)

	err := manager.AddLicense(context.Background(), "non_existent.lic", "test_key", "test_public_key")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "license file not found at 'non_existent.lic'")
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

	// add a license that to be deleted
	err := manager.AddLicense(ctx, "test_license.lic", "test_key", "test_public_key")
	assert.NoError(t, err, "Failed to add license")

	// check that the license was created
	_, err = manager.GetLicenseByID(ctx, "license_test_key")
	assert.NoError(t, err, "License should exist")

	// remove the license
	err = manager.RemoveLicense(ctx, "license_test_key")
	assert.NoError(t, err, "Failed to remove license")

	// check that the license is removed
	_, err = manager.GetLicenseByID(context.Background(), "license_test_key")
	assert.Error(t, err, "License should not exist after deletion")
	assert.Contains(t, err.Error(), "license with ID license_test_key: license not found")
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
	assert.Contains(t, err.Error(), "license with ID invalid_id: license not found")
}

func TestClaimLicense_Success(t *testing.T) {
	ctx := context.Background()
	store, dbConn := testutils.NewMemoryStore(t)
	defer testutils.CloseMemoryStore(dbConn)

	manager := licenses.NewManager(
		&licenses.Config{
			Strategy:          "fifo",
			EnabledAudit:      true,
			ExtendOnHeartbeat: true,
		},
		func(filename string) ([]byte, error) {
			return []byte("mock_certificate"), nil
		},
		func(cert []byte) licenses.LicenseVerifier {
			return &testutils.FakeLicenseVerifier{}
		},
	)
	manager.AttachStore(store)

	err := manager.AddLicense(ctx, "test_license.lic", "test_key", "test_public_key")
	assert.NoError(t, err)

	result, err := manager.ClaimLicense(ctx, "test_fingerprint")
	assert.NoError(t, err)
	assert.Equal(t, result.Status, licenses.OperationStatusCreated)
	assert.NotNil(t, result.License)
	assert.NotNil(t, result.License.LastClaimedAt)
	assert.NotNil(t, result.License.CreatedAt)
	assert.Equal(t, "test_key", result.License.Key)
}

func TestClaimLicense_NoLicensesAvailable(t *testing.T) {
	ctx := context.Background()
	store, dbConn := testutils.NewMemoryStore(t)
	defer testutils.CloseMemoryStore(dbConn)

	manager := licenses.NewManager(
		&licenses.Config{
			Strategy:          "fifo",
			EnabledAudit:      true,
			ExtendOnHeartbeat: true,
		},
		func(filename string) ([]byte, error) {
			return []byte("mock_certificate"), nil
		},
		func(cert []byte) licenses.LicenseVerifier {
			return &testutils.FakeLicenseVerifier{}
		},
	)
	manager.AttachStore(store)

	result, err := manager.ClaimLicense(ctx, "test_fingerprint")
	assert.NoError(t, err)
	assert.Equal(t, result.Status, licenses.OperationStatusNoLicensesAvailable)
	assert.Nil(t, result.License)
}

func TestClaimLicense_AlreadyClaimed_WithHeartbeatEnabled(t *testing.T) {
	ctx := context.Background()
	store, dbConn := testutils.NewMemoryStore(t)
	defer testutils.CloseMemoryStore(dbConn)

	manager := licenses.NewManager(
		&licenses.Config{
			Strategy:          "fifo",
			EnabledAudit:      true,
			ExtendOnHeartbeat: true,
		},
		func(filename string) ([]byte, error) {
			return []byte("mock_certificate"), nil
		},
		func(cert []byte) licenses.LicenseVerifier {
			return &testutils.FakeLicenseVerifier{}
		},
	)
	manager.AttachStore(store)

	err := manager.AddLicense(ctx, "test_license.lic", "test_key", "test_public_key")
	assert.NoError(t, err)

	// First getting the license
	result1, err1 := manager.ClaimLicense(ctx, "test_fingerprint")
	assert.NoError(t, err1)
	assert.Equal(t, result1.Status, licenses.OperationStatusCreated)
	assert.NotNil(t, result1.License)

	// The second getting the license with the same node.fingerprint
	result2, err2 := manager.ClaimLicense(ctx, "test_fingerprint")
	assert.NoError(t, err2)
	assert.Equal(t, result2.Status, licenses.OperationStatusExtended)
	assert.NotNil(t, result2.License)
}

func TestClaimLicense_AlreadyClaimed_WithHeartbeatDisabled(t *testing.T) {
	ctx := context.Background()
	store, dbConn := testutils.NewMemoryStore(t)
	defer testutils.CloseMemoryStore(dbConn)

	manager := licenses.NewManager(
		&licenses.Config{
			Strategy:          "fifo",
			EnabledAudit:      true,
			ExtendOnHeartbeat: false,
		},
		func(filename string) ([]byte, error) {
			return []byte("mock_certificate"), nil
		},
		func(cert []byte) licenses.LicenseVerifier {
			return &testutils.FakeLicenseVerifier{}
		},
	)
	manager.AttachStore(store)

	// add the license
	err := manager.AddLicense(ctx, "test_license.lic", "test_key", "test_public_key")
	assert.NoError(t, err)

	// first getting the license
	result1, err1 := manager.ClaimLicense(ctx, "test_fingerprint")
	assert.NoError(t, err1)
	assert.Equal(t, result1.Status, licenses.OperationStatusCreated)
	assert.NotNil(t, result1.License)

	// the second getting the license with the same node.fingerprint
	result2, err2 := manager.ClaimLicense(ctx, "test_fingerprint")
	assert.NoError(t, err2)
	assert.Equal(t, result2.Status, licenses.OperationStatusConflict)
	assert.Nil(t, result2.License)
}

func TestClaimLicense_FIFO_Strategy(t *testing.T) {
	ctx := context.Background()
	store, dbConn := testutils.NewMemoryStore(t)
	defer testutils.CloseMemoryStore(dbConn)

	fakeLicenseVerifier := func(cert []byte) licenses.LicenseVerifier {
		return &testutils.FakeLicenseVerifier{}
	}

	manager := licenses.NewManager(
		&licenses.Config{
			Strategy:          "fifo",
			EnabledAudit:      true,
			ExtendOnHeartbeat: true,
		},
		func(filename string) ([]byte, error) {
			if filename == "license1.lic" {
				return []byte("mock_certificate_1"), nil
			}
			if filename == "license2.lic" {
				return []byte("mock_certificate_2"), nil
			}
			return []byte("mock_certificate_3"), nil
		},
		fakeLicenseVerifier,
	)
	manager.AttachStore(store)

	err := manager.AddLicense(ctx, "license1.lic", "key1", "public_key")
	assert.NoError(t, err)
	err = manager.AddLicense(ctx, "license2.lic", "key2", "public_key")
	assert.NoError(t, err)
	err = manager.AddLicense(ctx, "license3.lic", "key3", "public_key")
	assert.NoError(t, err)

	result, err := manager.ClaimLicense(ctx, "test_fingerprint")
	assert.NoError(t, err)
	assert.Equal(t, licenses.OperationStatusCreated, result.Status)
	assert.NotNil(t, result.License)
	assert.Equal(t, "license_key1", result.License.ID)
	assert.Equal(t, "key1", result.License.Key)
}

func TestClaimLicense_LIFO_Strategy(t *testing.T) {
	ctx := context.Background()
	store, dbConn := testutils.NewMemoryStore(t)
	defer testutils.CloseMemoryStore(dbConn)

	fakeLicenseVerifier := func(cert []byte) licenses.LicenseVerifier {
		return &testutils.FakeLicenseVerifier{}
	}

	manager := licenses.NewManager(
		&licenses.Config{
			Strategy:          "lifo",
			EnabledAudit:      true,
			ExtendOnHeartbeat: true,
		},
		func(filename string) ([]byte, error) {
			if filename == "license1.lic" {
				return []byte("mock_certificate_1"), nil
			}
			if filename == "license2.lic" {
				return []byte("mock_certificate_2"), nil
			}
			return []byte("mock_certificate_3"), nil
		},
		fakeLicenseVerifier,
	)
	manager.AttachStore(store)

	err := manager.AddLicense(ctx, "license1.lic", "key1", "public_key")
	assert.NoError(t, err)

	err = manager.AddLicense(ctx, "license2.lic", "key2", "public_key")
	assert.NoError(t, err)

	err = manager.AddLicense(ctx, "license3.lic", "key3", "public_key")
	assert.NoError(t, err)

	// we need to update created_at manually, because in tests the records are created very quickly in the same time
	// update created_at for license_key1
	_, err = dbConn.ExecContext(ctx, `UPDATE licenses SET created_at = datetime('now', '-3 seconds') WHERE id = 'license_key1'`)
	assert.NoError(t, err)

	// update created_at for license_key2
	_, err = dbConn.ExecContext(ctx, `UPDATE licenses SET created_at = datetime('now', '-2 seconds') WHERE id = 'license_key2'`)
	assert.NoError(t, err)

	// update created_at for license_key3
	_, err = dbConn.ExecContext(ctx, `UPDATE licenses SET created_at = datetime('now', '-1 seconds') WHERE id = 'license_key3'`)
	assert.NoError(t, err)

	// claim license with LIFO strategy
	result, err := manager.ClaimLicense(ctx, "test_fingerprint")
	assert.NoError(t, err)
	assert.Equal(t, licenses.OperationStatusCreated, result.Status)
	assert.NotNil(t, result.License)
	assert.Equal(t, "license_key3", result.License.ID)
	assert.Equal(t, "key3", result.License.Key)
}

func TestReleaseLicense_Success(t *testing.T) {
	ctx := context.Background()
	store, dbConn := testutils.NewMemoryStore(t)
	defer testutils.CloseMemoryStore(dbConn)

	manager := licenses.NewManager(
		&licenses.Config{
			Strategy:     "fifo",
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

	// adding the license and then getting it
	err := manager.AddLicense(ctx, "test_license.lic", "test_key", "test_public_key")
	assert.NoError(t, err)

	result, err := manager.ClaimLicense(ctx, "test_fingerprint")
	assert.NoError(t, err)
	assert.Equal(t, result.Status, licenses.OperationStatusCreated)
	assert.NotNil(t, result.License)

	// Release the license
	releaseResult, err := manager.ReleaseLicense(ctx, "test_fingerprint")
	assert.NoError(t, err)
	assert.Equal(t, releaseResult.Status, licenses.OperationStatusSuccess)

	// checking the empty node
	_, err = store.GetNodeByFingerprint(ctx, "test_fingerprint")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, sql.ErrNoRows))

	// checking the free license
	claimedLicense, err := store.GetLicenseByID(ctx, result.License.ID)
	assert.NoError(t, err)
	assert.Nil(t, claimedLicense.NodeID)
	assert.NotNil(t, claimedLicense.LastReleasedAt)
}

func TestReleaseLicense_NodeNotFound(t *testing.T) {
	ctx := context.Background()
	store, dbConn := testutils.NewMemoryStore(t)
	defer testutils.CloseMemoryStore(dbConn)

	manager := licenses.NewManager(
		&licenses.Config{
			Strategy:     "fifo",
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

	// release the license for non existent node
	result, err := manager.ReleaseLicense(ctx, "non_existent_fingerprint")
	assert.Nil(t, err)
	assert.Equal(t, result.Status, licenses.OperationStatusNotFound)
}

func TestReleaseLicense_LicenseNotFound(t *testing.T) {
	ctx := context.Background()
	store, dbConn := testutils.NewMemoryStore(t)
	defer testutils.CloseMemoryStore(dbConn)

	manager := licenses.NewManager(
		&licenses.Config{},
		func(filename string) ([]byte, error) {
			return []byte("mock_certificate"), nil
		},
		func(cert []byte) licenses.LicenseVerifier {
			return &testutils.FakeLicenseVerifier{}
		},
	)
	manager.AttachStore(store)

	_, err := store.InsertNode(ctx, "test_fingerprint")
	assert.NoError(t, err)

	result, err := manager.ReleaseLicense(ctx, "test_fingerprint")
	assert.Nil(t, err)
	assert.Equal(t, result.Status, licenses.OperationStatusNotFound)
}

func TestCleanupInactiveNodes(t *testing.T) {
	ctx := context.Background()
	store, dbConn := testutils.NewMemoryStore(t)
	defer testutils.CloseMemoryStore(dbConn)

	manager := licenses.NewManager(
		&licenses.Config{},
		func(filename string) ([]byte, error) {
			if filename == "license1.lic" {
				return []byte("mock_certificate_1"), nil
			}
			return []byte("mock_certificate_2"), nil
		},
		func(cert []byte) licenses.LicenseVerifier {
			return &testutils.FakeLicenseVerifier{}
		},
	)
	manager.AttachStore(store)

	err := manager.AddLicense(ctx, "license1.lic", "test_key", "test_public_key")
	assert.NoError(t, err)
	err = manager.AddLicense(ctx, "license2.lic", "test_key_2", "test_public_key")
	assert.NoError(t, err)

	// сlaim the first license
	result1, err := manager.ClaimLicense(ctx, "test_fingerprint")
	assert.NoError(t, err)
	assert.Equal(t, licenses.OperationStatusCreated, result1.Status)
	assert.NotNil(t, result1.License)

	// claim the second license
	result2, err := manager.ClaimLicense(ctx, "test_fingerprint_2")
	assert.NoError(t, err)
	assert.Equal(t, licenses.OperationStatusCreated, result2.Status)
	assert.NotNil(t, result2.License)

	// simulate inactive node by updating last_heartbeat_at
	_, err = dbConn.ExecContext(ctx, `
        UPDATE nodes SET last_heartbeat_at = datetime('now', '-120 seconds') WHERE fingerprint = ?;
    `, "test_fingerprint")
	assert.NoError(t, err)

	err = manager.CleanupInactiveNodes(ctx, 30*time.Second)
	assert.NoError(t, err)

	license, err := manager.GetLicenseByID(ctx, result1.License.ID)
	assert.NoError(t, err)
	assert.Nil(t, license.NodeID)
	assert.NotNil(t, license.ID)

	license2, err := manager.GetLicenseByID(ctx, result2.License.ID)
	assert.NoError(t, err)
	assert.NotNil(t, license2.NodeID)

	node, err := store.GetNodeByFingerprint(ctx, "test_fingerprint")
	assert.Error(t, err)
	assert.Equal(t, node.Fingerprint, "")

	node2, err := store.GetNodeByFingerprint(ctx, "test_fingerprint_2")
	assert.NoError(t, err)
	assert.Equal(t, node2.Fingerprint, "test_fingerprint_2")
}
