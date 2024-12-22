package licenses

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/keygen-sh/keygen-go/v3"
	"github.com/keygen-sh/keygen-relay/internal/db"
	"github.com/mattn/go-sqlite3"
)

type OperationStatus int

const (
	OperationStatusSuccess OperationStatus = iota
	OperationStatusCreated
	OperationStatusExtended
	OperationStatusConflict
	OperationStatusNotFound
	OperationStatusNoLicensesAvailable
)

var (
	ErrNoLicenses      = errors.New("license pool is empty")
	ErrLicenseNotFound = errors.New("license not found")
)

type LicenseOperationResult struct {
	License *db.License
	Status  OperationStatus
}

type FileReaderFunc func(filename string) ([]byte, error)

type Manager interface {
	AddLicense(ctx context.Context, licenseFilePath string, licenseKey string, publicKeyPath string) error
	RemoveLicense(ctx context.Context, id string) error
	ListLicenses(ctx context.Context) ([]db.License, error)
	GetLicenseByGUID(ctx context.Context, id string) (*db.License, error)
	AttachStore(store db.Store)
	ClaimLicense(ctx context.Context, fingerprint string) (*LicenseOperationResult, error)
	ReleaseLicense(ctx context.Context, fingerprint string) (*LicenseOperationResult, error)
	Config() *Config
	CullDeadNodes(ctx context.Context, ttl time.Duration) error
}

type manager struct {
	store      db.Store
	config     *Config
	dataReader FileReaderFunc
	verifier   func(cert []byte) LicenseVerifier
}

func NewManager(config *Config, dataReader FileReaderFunc, verifier func(cert []byte) LicenseVerifier) Manager {
	return &manager{
		config:     config,
		dataReader: dataReader,
		verifier:   verifier,
	}
}

func (m *manager) AttachStore(store db.Store) {
	m.store = store
}

func (m *manager) AddLicense(ctx context.Context, licenseFilePath string, licenseKey string, publicKey string) error {
	slog.Debug("starting to add a new license", "filePath", licenseFilePath)

	cert, err := m.dataReader(licenseFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			slog.Warn("license file not found", "filePath", licenseFilePath)

			return fmt.Errorf("license file not found at '%s'", licenseFilePath)
		}

		slog.Error("failed to read license file", "filePath", licenseFilePath, "error", err)

		return fmt.Errorf("failed to read license file: %w", err)
	}

	slog.Debug("successfully read the license file", "filePath", licenseFilePath)

	lic := m.verifier(cert)
	keygen.PublicKey = publicKey

	if err := lic.Verify(); err != nil {
		return fmt.Errorf("license verification failed: %w", err)
	}

	dec, err := lic.Decrypt(licenseKey)
	if err != nil {
		return fmt.Errorf("license decryption failed: %w", err)
	}

	guid := dec.License.ID
	key := dec.License.Key

	license, err := m.store.InsertLicense(ctx, guid, cert, key)
	if err != nil {
		slog.Debug("failed to insert license", "licenseGuid", guid, "error", err)

		if isUniqueConstraintError(err) {
			return fmt.Errorf("license with the provided key already exists")
		}

		return fmt.Errorf("failed to insert license: %w", err)
	}

	// Log audit, but do not fail the operation if it fails
	if m.config.EnabledAudit {
		if err := m.store.InsertAuditLog(ctx, db.EventTypeLicenseAdded, db.EntityTypeLicense, license.ID); err != nil {
			slog.Debug("failed to insert audit log", "licenseGuid", guid, "error", err)
		}
	}

	slog.Debug("added license successfully", "licenseGuid", guid)

	return nil
}

func (m *manager) RemoveLicense(ctx context.Context, guid string) error {
	slog.Debug("starting to remove license", "licenseGuid", guid)

	license, err := m.store.DeleteLicenseByGUID(ctx, guid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("license %s not found", guid)
		}

		slog.Debug("failed to delete license", "licenseGuid", guid, "error", err)

		return fmt.Errorf("failed to delete license: %w", err)
	}

	// Log audit, but do not fail the operation if it fails
	if m.config.EnabledAudit {
		if err := m.store.InsertAuditLog(ctx, db.EventTypeLicenseRemoved, db.EntityTypeLicense, license.ID); err != nil {
			slog.Debug("failed to insert audit log", "licenseGuid", guid, "error", err)
		}
	}

	slog.Debug("removed license successfully", "licenseGuid", guid)

	return nil
}

func (m *manager) ListLicenses(ctx context.Context) ([]db.License, error) {
	slog.Debug("fetching licenses")

	licenses, err := m.store.GetAllLicenses(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Debug("license pool is empty")

			return nil, ErrNoLicenses
		}

		slog.Debug("failed to fetch licenses", "error", err)

		return nil, err
	}

	slog.Debug("fetched licenses successfully", "count", len(licenses))

	return licenses, nil
}

func (m *manager) GetLicenseByGUID(ctx context.Context, guid string) (*db.License, error) {
	slog.Debug("fetching license", "licenseGuid", guid)

	license, err := m.store.GetLicenseByGUID(ctx, guid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Debug("license not found", "licenseGuid", guid)

			return nil, fmt.Errorf("license %s: %w", guid, ErrLicenseNotFound)
		}

		slog.Debug("failed to fetch license by ID", "licenseGuid", guid, "error", err)

		return nil, err
	}

	slog.Debug("fetched license successfully", "licenseGuid", guid)

	return license, nil
}

func (m *manager) ClaimLicense(ctx context.Context, fingerprint string) (*LicenseOperationResult, error) {
	tx, err := m.store.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	qtx := m.store.WithTx(tx)
	defer tx.Rollback()

	node, err := m.findOrActivateNode(ctx, *qtx, fingerprint)
	if err != nil {
		return nil, fmt.Errorf("failed to find or activate node: %w", err)
	}

	claimedLicense, err := qtx.GetLicenseByNodeID(ctx, &node.ID)

	// extend the lease if the node already has a lease on a license
	if err == nil {
		if !m.config.ExtendOnHeartbeat { // if heartbeat is disabled, we can't extend the claimed license
			slog.Warn("failed to claim license due to conflict due to heartbeat disabled", "nodeID", node.ID, "nodeFingerprint", node.Fingerprint)

			return &LicenseOperationResult{Status: OperationStatusConflict}, nil
		}

		if err := qtx.PingNodeHeartbeatByFingerprint(ctx, fingerprint); err != nil {
			return nil, fmt.Errorf("failed to update node heartbeat: %w", err)
		}

		if err := tx.Commit(); err != nil {
			return nil, fmt.Errorf("failed to commit transaction: %w", err)
		}

		if m.config.EnabledAudit {
			if err := m.store.BulkInsertAuditLogs(ctx, []db.BulkInsertAuditLogParams{
				{EventTypeID: db.EventTypeLicenseLeaseExtended, EntityTypeID: db.EntityTypeLicense, EntityID: claimedLicense.ID},
				{EventTypeID: db.EventTypeNodeHeartbeatPing, EntityTypeID: db.EntityTypeNode, EntityID: node.ID},
			}); err != nil {
				slog.Warn("failed to insert audit logs", "error", err)
			}
		}

		slog.Info("license extended successfully", "licenseGuid", claimedLicense.Guid)

		return &LicenseOperationResult{
			License: claimedLicense,
			Status:  OperationStatusExtended,
		}, nil
	}

	// claim a new lease on a license if node doesn't have a lease
	newLicense, err := qtx.ClaimLicenseByStrategy(ctx, m.config.Strategy, &node.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("no licenses available in pool", "nodeId", node.ID, "nodeFingerprint", node.Fingerprint)

			return &LicenseOperationResult{Status: OperationStatusNoLicensesAvailable}, nil
		}

		return nil, fmt.Errorf("failed to claim license: %w", err)
	}

	if err := qtx.PingNodeHeartbeatByFingerprint(ctx, fingerprint); err != nil {
		return nil, fmt.Errorf("failed to update node claim: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	if m.config.EnabledAudit {
		if err := m.store.BulkInsertAuditLogs(ctx, []db.BulkInsertAuditLogParams{
			{EventTypeID: db.EventTypeLicenseLeased, EntityTypeID: db.EntityTypeLicense, EntityID: newLicense.ID},
			{EventTypeID: db.EventTypeNodeHeartbeatPing, EntityTypeID: db.EntityTypeNode, EntityID: node.ID},
		}); err != nil {
			slog.Warn("failed to insert audit logs", "error", err)
		}
	}

	slog.Info("new license claimed successfully", "licenseGuid", newLicense.Guid)

	return &LicenseOperationResult{
		License: newLicense,
		Status:  OperationStatusCreated,
	}, nil
}

func (m *manager) ReleaseLicense(ctx context.Context, fingerprint string) (*LicenseOperationResult, error) {
	tx, err := m.store.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	qtx := m.store.WithTx(tx)
	defer tx.Rollback()

	node, err := qtx.GetNodeByFingerprint(ctx, fingerprint)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("license release failed - node not found", "fingerprint", fingerprint)

			return &LicenseOperationResult{Status: OperationStatusNotFound}, nil
		}

		return nil, fmt.Errorf("failed to fetch node: %w", err)
	}

	claimedLicense, err := qtx.GetLicenseByNodeID(ctx, &node.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("license release failed - claimed license not found", "nodeFingerprint", node.Fingerprint)

			return &LicenseOperationResult{Status: OperationStatusNotFound}, nil
		}

		return nil, fmt.Errorf("failed to fetch claimed license: %w", err)
	}

	if err := qtx.ReleaseLicenseByNodeID(ctx, &node.ID); err != nil {
		return nil, fmt.Errorf("failed to release license: %w", err)
	}

	if err := qtx.DeactivateNodeByFingerprint(ctx, fingerprint); err != nil {
		return nil, fmt.Errorf("failed to deactivate node: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	if m.config.EnabledAudit {
		if err := m.store.BulkInsertAuditLogs(ctx, []db.BulkInsertAuditLogParams{
			{EventTypeID: db.EventTypeLicenseReleased, EntityTypeID: db.EntityTypeLicense, EntityID: claimedLicense.ID},
			{EventTypeID: db.EventTypeNodeDeactivated, EntityTypeID: db.EntityTypeNode, EntityID: node.ID},
		}); err != nil {
			slog.Warn("failed to insert audit logs", "error", err)
		}
	}

	slog.Info("license released successfully", "licenseGuid", claimedLicense.Guid)

	return &LicenseOperationResult{Status: OperationStatusSuccess}, nil
}

func (m *manager) Config() *Config {
	return m.config
}

func (m *manager) findOrActivateNode(ctx context.Context, store db.Store, fingerprint string) (*db.Node, error) {
	node, err := store.GetNodeByFingerprint(ctx, fingerprint)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			node, err = store.ActivateNode(ctx, fingerprint)
			if err != nil {
				slog.Error("failed to insert node", "nodeFingerprint", fingerprint, "error", err)

				return nil, fmt.Errorf("failed to insert node: %w", err)
			}

			if m.config.EnabledAudit {
				if err := store.InsertAuditLog(ctx, db.EventTypeNodeActivated, db.EntityTypeNode, node.ID); err != nil {
					slog.Warn("failed to insert audit log", "nodeID", node.ID, "nodeFingerprint", node.Fingerprint, "error", err)
				}
			}
		} else {
			slog.Error("failed to find node", "nodeFingerprint", fingerprint, "error", err)

			return nil, fmt.Errorf("failed to fetch node: %w", err)
		}
	}

	return node, nil
}

func (m *manager) CullDeadNodes(ctx context.Context, ttl time.Duration) error {
	tx, err := m.store.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	qtx := m.store.WithTx(tx)
	defer tx.Rollback()

	licenses, err := qtx.ReleaseLicensesFromDeadNodes(ctx, ttl)
	if err != nil {
		slog.Error("failed to release licenses from dead nodes", "error", err)

		return err
	}

	nodes, err := qtx.DeactivateDeadNodes(ctx, ttl)
	if err != nil {
		slog.Error("failed to deactivate dead nodes", "error", err)

		return err
	}

	if err := tx.Commit(); err != nil {
		slog.Error("failed to commit transaction", "error", err)

		return err
	}

	if m.config.EnabledAudit {
		var logs []db.BulkInsertAuditLogParams

		for _, license := range licenses {
			logs = append(logs, db.BulkInsertAuditLogParams{
				EventTypeID:  db.EventTypeLicenseLeaseExpired,
				EntityTypeID: db.EntityTypeLicense,
				EntityID:     license.ID,
			})
		}

		for _, node := range nodes {
			logs = append(logs, db.BulkInsertAuditLogParams{
				EventTypeID:  db.EventTypeNodeCulled,
				EntityTypeID: db.EntityTypeNode,
				EntityID:     node.ID,
			})
		}

		if err := m.store.BulkInsertAuditLogs(ctx, logs); err != nil {
			slog.Warn("failed to insert audit logs", "error", err)
		}
	}

	slog.Debug("successfully released licenses and culled dead nodes", "licenseCount", len(licenses), "nodeCount", len(nodes))

	return nil
}

func isUniqueConstraintError(err error) bool {
	var sqliteErr sqlite3.Error

	ok := errors.As(err, &sqliteErr)

	if ok && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
		return true
	}

	return false
}
