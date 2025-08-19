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
	ErrBadPool         = errors.New("pool not found")
)

type LicenseOperationResult struct {
	License *db.License
	Status  OperationStatus
}

type FileReaderFunc func(filename string) ([]byte, error)

type Manager interface {
	AddLicense(ctx context.Context, pool *string, licenseFilePath string, licenseKey string, publicKeyPath string) (*db.License, error)
	RemoveLicense(ctx context.Context, pool *string, id string) error
	ListLicenses(ctx context.Context, pool *string) ([]db.License, error)
	GetLicenseByGUID(ctx context.Context, pool *string, id string) (*db.License, error)
	AttachStore(store db.Store)
	ClaimLicense(ctx context.Context, pool *string, fingerprint string) (*LicenseOperationResult, error)
	ReleaseLicense(ctx context.Context, pool *string, fingerprint string) (*LicenseOperationResult, error)
	Config() *Config
	CullDeadNodes(ctx context.Context, ttl time.Duration) ([]db.Node, error)
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

func (m *manager) AddLicense(ctx context.Context, poolName *string, licenseFilePath string, licenseKey string, publicKey string) (*db.License, error) {
	slog.Debug("starting to add a new license", "pool", deref(poolName), "filePath", licenseFilePath)

	cert, err := m.dataReader(licenseFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			slog.Warn("license file not found", "filePath", licenseFilePath)

			return nil, fmt.Errorf("license file not found at '%s'", licenseFilePath)
		}

		slog.Error("failed to read license file", "filePath", licenseFilePath, "error", err)

		return nil, fmt.Errorf("failed to read license file: %w", err)
	}

	slog.Debug("successfully read the license file", "filePath", licenseFilePath)

	lic := m.verifier(cert)
	keygen.PublicKey = publicKey

	if err := lic.Verify(); err != nil {
		return nil, fmt.Errorf("license verification failed: %w", err)
	}

	dec, err := lic.Decrypt(licenseKey)
	if err != nil {
		return nil, fmt.Errorf("license decryption failed: %w", err)
	}

	guid := dec.License.ID
	key := dec.License.Key
	var pool *db.Pool

	tx, err := m.store.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	qtx := m.store.WithTx(tx)
	defer tx.Rollback()

	if poolName != nil {
		pool, err = m.findOrCreatePool(ctx, *qtx, *poolName)
		if err != nil {
			return nil, fmt.Errorf("failed to find or create pool: %w", err)
		}
	}

	license, err := qtx.InsertLicense(ctx, pool, guid, cert, key)
	if err != nil {
		slog.Debug("failed to insert license", "licenseGuid", guid, "error", err)

		if isUniqueConstraintError(err) {
			return nil, fmt.Errorf("license with the provided key already exists")
		}

		return nil, fmt.Errorf("failed to insert license: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Log audit, but do not fail the operation if it fails
	if m.config.EnabledAudit {
		if err := m.store.InsertAuditLog(ctx, pool, db.EventTypeLicenseAdded, db.EntityTypeLicense, license.ID); err != nil {
			slog.Debug("failed to insert audit log", "licenseGuid", guid, "error", err)
		}
	}

	slog.Debug("added license successfully", "licenseGuid", guid)

	return license, nil
}

func (m *manager) RemoveLicense(ctx context.Context, poolName *string, guid string) error {
	slog.Debug("starting to remove license", "pool", deref(poolName), "licenseGuid", guid)

	var pool *db.Pool
	var err error

	if poolName != nil {
		pool, err = m.store.GetPoolByName(ctx, *poolName)
		if err != nil {
			slog.Debug("failed to fetch pool", "poolName", *poolName, "error", err)

			return ErrBadPool
		}
	}

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
		if err := m.store.InsertAuditLog(ctx, pool, db.EventTypeLicenseRemoved, db.EntityTypeLicense, license.ID); err != nil {
			slog.Debug("failed to insert audit log", "licenseGuid", guid, "error", err)
		}
	}

	slog.Debug("removed license successfully", "licenseGuid", guid)

	return nil
}

func (m *manager) ListLicenses(ctx context.Context, poolName *string) ([]db.License, error) {
	slog.Debug("fetching licenses", "pool", deref(poolName))

	var licenses []db.License
	var pool *db.Pool
	var err error

	if poolName != nil {
		pool, err = m.store.GetPoolByName(ctx, *poolName)
		if err != nil {
			slog.Debug("failed to fetch pool", "poolName", *poolName, "error", err)

			return nil, ErrBadPool
		}

		licenses, err = m.store.GetPooledLicenses(ctx, pool)
	} else {
		licenses, err = m.store.GetAllLicenses(ctx)
	}
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

func (m *manager) GetLicenseByGUID(ctx context.Context, poolName *string, guid string) (*db.License, error) {
	slog.Debug("fetching license", "licenseGuid", guid)

	var license *db.License
	var pool *db.Pool
	var err error

	if poolName != nil {
		pool, err = m.store.GetPoolByName(ctx, *poolName)
		if err != nil {
			slog.Debug("failed to fetch pool", "poolName", *poolName, "error", err)

			return nil, ErrBadPool
		}

		license, err = m.store.GetPooledLicenseByGUID(ctx, pool, guid)
	} else {
		license, err = m.store.GetLicenseByGUID(ctx, guid)
	}

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

func (m *manager) ClaimLicense(ctx context.Context, poolName *string, fingerprint string) (*LicenseOperationResult, error) {
	tx, err := m.store.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	qtx := m.store.WithTx(tx)
	defer tx.Rollback()

	var pool *db.Pool
	if poolName != nil {
		pool, err = m.store.GetPoolByName(ctx, *poolName)
		if err != nil {
			slog.Debug("failed to fetch pool", "poolName", *poolName, "error", err)

			return nil, ErrBadPool
		}
	}

	node, err := m.findOrActivateNode(ctx, *qtx, pool, fingerprint)
	if err != nil {
		return nil, fmt.Errorf("failed to find or activate node: %w", err)
	}

	var license *db.License
	if pool != nil {
		license, err = qtx.GetPooledLicenseByNodeID(ctx, pool, &node.ID)
	} else {
		license, err = qtx.GetLicenseByNodeID(ctx, &node.ID)
	}

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
				{EventTypeID: db.EventTypeLicenseLeaseExtended, EntityTypeID: db.EntityTypeLicense, EntityID: license.ID, Pool: pool},
				{EventTypeID: db.EventTypeNodeHeartbeatPing, EntityTypeID: db.EntityTypeNode, EntityID: node.ID, Pool: pool},
			}); err != nil {
				slog.Warn("failed to insert audit logs", "error", err)
			}
		}

		slog.Info("lease extended successfully", "licenseGuid", license.Guid)

		return &LicenseOperationResult{
			License: license,
			Status:  OperationStatusExtended,
		}, nil
	}

	// claim a new lease on a license if node doesn't have a lease
	if pool != nil {
		license, err = qtx.ClaimPooledLicenseByStrategy(ctx, pool, m.config.Strategy, &node.ID)
	} else {
		license, err = qtx.ClaimLicenseByStrategy(ctx, m.config.Strategy, &node.ID)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("no licenses available in pool", "nodeId", node.ID, "nodeFingerprint", node.Fingerprint)

			return &LicenseOperationResult{Status: OperationStatusNoLicensesAvailable}, nil
		}

		if isUniqueConstraintError(err) {
			return &LicenseOperationResult{Status: OperationStatusConflict}, nil
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
			{Pool: pool, EventTypeID: db.EventTypeLicenseLeased, EntityTypeID: db.EntityTypeLicense, EntityID: license.ID},
			{Pool: pool, EventTypeID: db.EventTypeNodeHeartbeatPing, EntityTypeID: db.EntityTypeNode, EntityID: node.ID},
		}); err != nil {
			slog.Warn("failed to insert audit logs", "error", err)
		}
	}

	slog.Info("new lease claimed successfully", "licenseGuid", license.Guid)

	return &LicenseOperationResult{
		License: license,
		Status:  OperationStatusCreated,
	}, nil
}

func (m *manager) ReleaseLicense(ctx context.Context, poolName *string, fingerprint string) (*LicenseOperationResult, error) {
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

	var license *db.License
	var pool *db.Pool

	if poolName != nil {
		pool, err = m.store.GetPoolByName(ctx, *poolName)
		if err != nil {
			slog.Debug("failed to fetch pool", "poolName", *poolName, "error", err)

			return nil, ErrBadPool
		}

		license, err = qtx.GetPooledLicenseByNodeID(ctx, pool, &node.ID)
	} else {
		license, err = qtx.GetLicenseByNodeID(ctx, &node.ID)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("license release failed - claimed license not found", "nodeFingerprint", node.Fingerprint)

			return &LicenseOperationResult{Status: OperationStatusNotFound}, nil
		}

		return nil, fmt.Errorf("failed to fetch claimed license: %w", err)
	}

	if pool != nil {
		if err := qtx.ReleasePooledLicenseByNodeID(ctx, pool, &node.ID); err != nil {
			return nil, fmt.Errorf("failed to release license: %w", err)
		}
	} else {
		if err := qtx.ReleaseLicenseByNodeID(ctx, &node.ID); err != nil {
			return nil, fmt.Errorf("failed to release license: %w", err)
		}
	}

	if err := qtx.DeactivateNodeByFingerprint(ctx, fingerprint); err != nil {
		return nil, fmt.Errorf("failed to deactivate node: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	if m.config.EnabledAudit {
		if err := m.store.BulkInsertAuditLogs(ctx, []db.BulkInsertAuditLogParams{
			{Pool: pool, EventTypeID: db.EventTypeLicenseReleased, EntityTypeID: db.EntityTypeLicense, EntityID: license.ID},
			{Pool: pool, EventTypeID: db.EventTypeNodeDeactivated, EntityTypeID: db.EntityTypeNode, EntityID: node.ID},
		}); err != nil {
			slog.Warn("failed to insert audit logs", "error", err)
		}
	}

	slog.Info("license released successfully", "licenseGuid", license.Guid)

	return &LicenseOperationResult{Status: OperationStatusSuccess}, nil
}

func (m *manager) Config() *Config {
	return m.config
}

func (m *manager) findOrCreatePool(ctx context.Context, store db.Store, name string) (*db.Pool, error) {
	pool, err := store.GetPoolByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			pool, err = store.CreatePool(ctx, name)
			if err != nil {
				slog.Error("failed to insert pool", "poolName", name, "error", err)

				return nil, fmt.Errorf("failed to insert pool: %w", err)
			}

			if m.config.EnabledAudit {
				if err := store.InsertAuditLog(ctx, nil, db.EventTypePoolAdded, db.EntityTypePool, pool.ID); err != nil {
					slog.Warn("failed to insert audit log", "poolID", pool.ID, "poolName", pool.Name, "error", err)
				}
			}
		} else {
			slog.Error("failed to find pool", "poolName", name, "error", err)

			return nil, fmt.Errorf("failed to fetch pool: %w", err)
		}
	}

	return pool, nil
}

func (m *manager) findOrActivateNode(ctx context.Context, store db.Store, pool *db.Pool, fingerprint string) (*db.Node, error) {
	node, err := store.GetNodeByFingerprint(ctx, fingerprint)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			node, err = store.ActivateNode(ctx, fingerprint)
			if err != nil {
				slog.Error("failed to insert node", "nodeFingerprint", fingerprint, "error", err)

				return nil, fmt.Errorf("failed to insert node: %w", err)
			}

			if m.config.EnabledAudit {
				if err := store.InsertAuditLog(ctx, pool, db.EventTypeNodeActivated, db.EntityTypeNode, node.ID); err != nil {
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

func (m *manager) CullDeadNodes(ctx context.Context, ttl time.Duration) ([]db.Node, error) {
	tx, err := m.store.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	qtx := m.store.WithTx(tx)
	defer tx.Rollback()

	licenses, err := qtx.ReleaseLicensesFromDeadNodes(ctx, ttl)
	if err != nil {
		slog.Error("failed to release licenses from dead nodes", "error", err)

		return nil, err
	}

	nodes, err := qtx.DeactivateDeadNodes(ctx, ttl)
	if err != nil {
		slog.Error("failed to deactivate dead nodes", "error", err)

		return nil, err
	}

	if err := tx.Commit(); err != nil {
		slog.Error("failed to commit transaction", "error", err)

		return nil, err
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

	return nodes, nil
}

func isUniqueConstraintError(err error) bool {
	var sqliteErr sqlite3.Error

	if ok := errors.As(err, &sqliteErr); ok {
		return errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique)
	}

	return false
}

func deref[T any](t *T) T {
	var zero T

	if t == nil {
		return zero
	}

	return *t
}
