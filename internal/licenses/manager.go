package licenses

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/keygen-sh/keygen-go/v3"
	"log/slog"
	"strconv"
	"time"
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
	ErrNoLicenses      = errors.New("no licenses found")
	ErrLicenseNotFound = errors.New("license not found")
)

type LicenseOperationResult struct {
	License *License
	Status  OperationStatus
}

type FileReaderFunc func(filename string) ([]byte, error)

type Store interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)
	WithTx(tx *sql.Tx) Store
	InsertLicense(ctx context.Context, id string, file []byte, key string) error
	DeleteLicenseByIDTx(ctx context.Context, id string) error
	GetAllLicenses(ctx context.Context) ([]License, error)
	GetLicenseByID(ctx context.Context, id string) (License, error)
	GetLicenseByNodeID(ctx context.Context, nodeID *int64) (License, error)
	InsertNode(ctx context.Context, fingerprint string) (Node, error)
	GetNodeByFingerprint(ctx context.Context, fingerprint string) (Node, error)
	UpdateNodeHeartbeat(ctx context.Context, fingerprint string) error
	UpdateNodeHeartbeatAndClaimedAtByFingerprint(ctx context.Context, fingerprint string) error
	DeleteNodeByFingerprint(ctx context.Context, fingerprint string) error
	ReleaseLicenseByNodeID(ctx context.Context, nodeID *int64) error
	InsertAuditLog(ctx context.Context, action, entityType, entityID string) error
	ClaimUnclaimedLicenseRandom(ctx context.Context, nodeID *int64) (License, error)
	ClaimUnclaimedLicenseLIFO(ctx context.Context, nodeID *int64) (License, error)
	ClaimUnclaimedLicenseFIFO(ctx context.Context, nodeID *int64) (License, error)
	DeleteInactiveNodes(ctx context.Context, ttl time.Duration) error
}

type Manager interface {
	AddLicense(ctx context.Context, licenseFilePath string, licenseKey string, publicKeyPath string) error
	RemoveLicense(ctx context.Context, id string) error
	ListLicenses(ctx context.Context) ([]License, error)
	GetLicenseByID(ctx context.Context, id string) (License, error)
	AttachStore(store Store)
	ClaimLicense(ctx context.Context, fingerprint string) (*LicenseOperationResult, error)
	ReleaseLicense(ctx context.Context, fingerprint string) (*LicenseOperationResult, error)
	Config() *Config
	CleanupInactiveNodes(ctx context.Context, ttl time.Duration) error
}

type manager struct {
	store      Store
	config     *Config
	dataReader FileReaderFunc
	verifier   func(cert []byte) LicenseVerifier
}

type License struct {
	ID             string
	File           []byte
	Key            string
	Claims         int64
	LastClaimedAt  *string
	LastReleasedAt *string
	NodeID         *int64
	CreatedAt      *string
}

type Node struct {
	ID              int64
	Fingerprint     string
	ClaimedAt       *string
	LastHeartbeatAt *string
	CreatedAt       *string
}

type AuditLog struct {
	ID         int64
	Action     string
	EntityType string
	EntityID   string
	Timestamp  *string
	CreatedAt  *string
}

func NewManager(config *Config, dataReader FileReaderFunc, verifier func(cert []byte) LicenseVerifier) Manager {
	return &manager{
		config:     config,
		dataReader: dataReader,
		verifier:   verifier,
	}
}

func (m *manager) AttachStore(store Store) {
	m.store = store
}

func (m *manager) AddLicense(ctx context.Context, licenseFilePath string, licenseKey string, publicKey string) error {
	slog.Debug("Starting to add a new license", "filePath", licenseFilePath)

	cert, err := m.dataReader(licenseFilePath)
	if err != nil {
		slog.Error("failed to read license file", "filePath", licenseFilePath, "error", err)
		return fmt.Errorf("failed to read license file: %w", err)
	}

	slog.Debug("Successfully read the license file", "filePath", licenseFilePath)

	lic := m.verifier(cert)
	keygen.PublicKey = publicKey

	if err := lic.Verify(); err != nil {
		slog.Error("license verification failed", "filePath", licenseFilePath, "error", err)
		return fmt.Errorf("license verification failed: %w", err)
	}

	dec, err := lic.Decrypt(licenseKey)
	if err != nil {
		slog.Error("license decryption failed", "filePath", licenseFilePath, "error", err)
		return fmt.Errorf("license decryption failed: %w", err)
	}

	id := dec.License.ID
	key := dec.License.Key

	if err := m.store.InsertLicense(ctx, id, cert, key); err != nil {
		slog.Error("failed to insert license into store", "licenseID", id, "error", err)
		return fmt.Errorf("failed to insert license: %w", err)
	}

	// Log audit, but do not fail the operation if it fails
	if m.config.EnabledAudit {
		if err := m.store.InsertAuditLog(ctx, "add", "license", id); err != nil {
			slog.Warn("failed to insert audit log", "licenseID", id, "error", err)
		}
	}

	slog.Info("added license successfully", "licenseID", id)
	return nil
}

func (m *manager) RemoveLicense(ctx context.Context, id string) error {
	slog.Debug("Starting to remove license", "id", id)

	err := m.store.DeleteLicenseByIDTx(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("license with ID %s not found", id)
		}
		slog.Error("failed to delete license", "licenseID", id, "error", err)
		return fmt.Errorf("failed to delete license: %w", err)
	}

	// Log audit, but do not fail the operation if it fails
	if m.config.EnabledAudit {
		if err := m.store.InsertAuditLog(ctx, "remove", "license", id); err != nil {
			slog.Warn("failed to insert audit log", "licenseID", id, "error", err)
		}
	}

	slog.Info("removed license successfully", "licenseID", id)
	return nil
}

func (m *manager) ListLicenses(ctx context.Context) ([]License, error) {
	slog.Debug("Fetching licenses")

	licenses, err := m.store.GetAllLicenses(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("no licenses found")
			return nil, ErrNoLicenses
		}

		slog.Error("failed to fetch licenses", "error", err)
		return nil, err
	}

	slog.Info("fetched licenses successfully", "count", len(licenses))
	return licenses, nil
}

func (m *manager) GetLicenseByID(ctx context.Context, id string) (License, error) {
	slog.Debug("Fetching license", "id", id)

	license, err := m.store.GetLicenseByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("license not found", "licenseID", id)
			return License{}, fmt.Errorf("license with ID %s: %w", id, ErrLicenseNotFound)
		}

		slog.Error("failed to fetch license by ID", "licenseID", id, "error", err)
		return License{}, err
	}

	slog.Info("fetched license successfully", "licenseID", id)
	return license, nil
}

func (m *manager) ClaimLicense(ctx context.Context, fingerprint string) (*LicenseOperationResult, error) {
	tx, err := m.store.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	storeWithTx := m.store.WithTx(tx)
	defer tx.Rollback()

	node, err := m.fetchOrCreateNode(ctx, storeWithTx, fingerprint)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch or create node: %w", err)
	}

	claimedLicense, err := storeWithTx.GetLicenseByNodeID(ctx, &node.ID)

	if err == nil {
		if !m.config.ExtendOnHeartbeat { // if heartbeat is disabled, we can't extend the claimed license
			slog.Warn("license claim conflict due to heartbeat disabled", "nodeID", node.ID)
			return &LicenseOperationResult{Status: OperationStatusConflict}, nil
		}

		if err := storeWithTx.UpdateNodeHeartbeat(ctx, fingerprint); err != nil {
			return nil, fmt.Errorf("failed to update node heartbeat: %w", err)
		}

		if err := tx.Commit(); err != nil {
			return nil, fmt.Errorf("failed to commit transaction: %w", err)
		}

		if m.config.EnabledAudit {
			if err := m.store.InsertAuditLog(ctx, "claimed", "license", claimedLicense.ID); err != nil {
				slog.Warn("failed to insert audit log", "licenseID", claimedLicense.ID, "error", err)
			}
		}

		slog.Info("license extended successfully", "licenseID", claimedLicense.ID)
		return &LicenseOperationResult{
			License: &claimedLicense,
			Status:  OperationStatusExtended,
		}, nil
	}

	// claim a new license based on the strategy
	newLicense, err := m.selectLicenseClaimStrategy(ctx, storeWithTx, &node.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("no licenses available for claim", "nodeID", node.ID)
			return &LicenseOperationResult{Status: OperationStatusNoLicensesAvailable}, nil
		}
		return nil, fmt.Errorf("failed to claim license: %w", err)
	}

	// Update node claim timestamp
	if err := storeWithTx.UpdateNodeHeartbeatAndClaimedAtByFingerprint(ctx, fingerprint); err != nil {
		return nil, fmt.Errorf("failed to update node claim: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	if m.config.EnabledAudit {
		if err := m.store.InsertAuditLog(ctx, "claimed", "license", newLicense.ID); err != nil {
			slog.Warn("failed to insert audit log", "licenseID", newLicense.ID, "error", err)
		}
	}

	slog.Info("new license claimed successfully", "licenseID", newLicense.ID)
	return &LicenseOperationResult{
		License: &newLicense,
		Status:  OperationStatusCreated,
	}, nil
}

func (m *manager) ReleaseLicense(ctx context.Context, fingerprint string) (*LicenseOperationResult, error) {
	tx, err := m.store.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	storeWithTx := m.store.WithTx(tx)
	defer tx.Rollback()

	node, err := storeWithTx.GetNodeByFingerprint(ctx, fingerprint)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("license release failed - node not found", "fingerprint", fingerprint)
			return &LicenseOperationResult{Status: OperationStatusNotFound}, nil
		}
		return nil, fmt.Errorf("failed to fetch node: %w", err)
	}

	claimedLicense, err := storeWithTx.GetLicenseByNodeID(ctx, &node.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("license release failed - claimed license not found", "nodeID", node.ID)
			return &LicenseOperationResult{Status: OperationStatusNotFound}, nil
		}
		return nil, fmt.Errorf("failed to fetch claimed license: %w", err)
	}

	if err := storeWithTx.ReleaseLicenseByNodeID(ctx, &node.ID); err != nil {
		return nil, fmt.Errorf("failed to release license: %w", err)
	}

	if err := storeWithTx.DeleteNodeByFingerprint(ctx, fingerprint); err != nil {
		return nil, fmt.Errorf("failed to delete node: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	if m.config.EnabledAudit {
		if err := m.store.InsertAuditLog(ctx, "released", "license", claimedLicense.ID); err != nil {
			slog.Warn("failed to insert audit log", "licenseID", claimedLicense.ID, "error", err)
		}
	}

	slog.Info("license released successfully", "licenseID", claimedLicense.ID)
	return &LicenseOperationResult{Status: OperationStatusSuccess}, nil
}

func (m *manager) Config() *Config {
	return m.config
}

func (m *manager) fetchOrCreateNode(ctx context.Context, store Store, fingerprint string) (Node, error) {
	node, err := store.GetNodeByFingerprint(ctx, fingerprint)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			node, err = store.InsertNode(ctx, fingerprint)

			if err != nil {
				slog.Error("failed to insert node", "Fingerprint", fingerprint, "error", err)

				return Node{}, fmt.Errorf("failed to insert node: %w", err)
			}

			if m.config.EnabledAudit {
				err = store.InsertAuditLog(ctx, "inserted", "node", strconv.FormatInt(node.ID, 10))

				if err != nil {
					slog.Warn("failed to insert audit log", "nodeID", node.ID, "error", err)
				}
			}
		} else {
			slog.Error("failed to fetch node", "Fingerprint", fingerprint, "error", err)

			return Node{}, fmt.Errorf("failed to fetch node: %w", err)
		}
	}

	return node, nil
}

func (m *manager) selectLicenseClaimStrategy(ctx context.Context, store Store, nodeID *int64) (License, error) {
	switch m.config.Strategy {
	case "fifo":
		return store.ClaimUnclaimedLicenseFIFO(ctx, nodeID)
	case "lifo":
		return store.ClaimUnclaimedLicenseLIFO(ctx, nodeID)
	case "random":
		return store.ClaimUnclaimedLicenseRandom(ctx, nodeID)
	default:
		return store.ClaimUnclaimedLicenseFIFO(ctx, nodeID)
	}
}

func (m *manager) CleanupInactiveNodes(ctx context.Context, ttl time.Duration) error {
	return m.store.DeleteInactiveNodes(ctx, ttl)
}
