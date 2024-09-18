package licenses

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/keygen-sh/keygen-go/v3"
	"log/slog"
)

type FileReaderFunc func(filename string) ([]byte, error)

type Store interface {
	InsertLicense(ctx context.Context, id string, file []byte, key string) error
	DeleteLicenseByID(ctx context.Context, id string) error
	GetAllLicenses(ctx context.Context) ([]License, error)
	GetLicenseByID(ctx context.Context, id string) (License, error)
	InsertNode(ctx context.Context, fingerprint string) error
	GetNodeByFingerprint(ctx context.Context, fingerprint string) (Node, error)
	UpdateNodeHeartbeat(ctx context.Context, fingerprint string) error
	DeleteNodeByFingerprint(ctx context.Context, fingerprint string) error
	InsertAuditLog(ctx context.Context, action, entityType, entityID string) error
}

type Manager interface {
	AddLicense(ctx context.Context, licenseFilePath string, licenseKey string, publicKeyPath string) error
	RemoveLicense(ctx context.Context, id string) error
	ListLicenses(ctx context.Context) ([]License, error)
	GetLicenseByID(ctx context.Context, id string) (License, error)
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
	LastClaimedAt  sql.NullString
	LastReleasedAt sql.NullString
	NodeID         sql.NullInt64
}

type Node struct {
	ID              int64
	Fingerprint     string
	ClaimedAt       sql.NullString
	LastHeartbeatAt sql.NullString
}

type AuditLog struct {
	ID         int64
	Action     string
	EntityType string
	EntityID   string
	Timestamp  sql.NullString
}

func NewManager(store Store, config *Config, dataReader FileReaderFunc, verifier func(cert []byte) LicenseVerifier) Manager {
	return &manager{
		store:      store,
		config:     config,
		dataReader: dataReader,
		verifier:   verifier,
	}
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

	if m.config.EnabledAudit {
		if err := m.store.InsertAuditLog(ctx, "add", "license", id); err != nil {
			slog.Error("failed to insert audit log", "licenseID", id, "error", err)
			return fmt.Errorf("failed to insert audit log: %w", err)
		}
	}

	slog.Info("added license", "licenseID", id)
	return nil
}

func (m *manager) RemoveLicense(ctx context.Context, id string) error {
	slog.Debug("Starting to remove license", "id", id)

	if err := m.store.DeleteLicenseByID(ctx, id); err != nil {
		slog.Error("failed to remove license", "licenseID", id, "error", err)
		return fmt.Errorf("failed to remove license: %w", err)
	}

	if m.config.EnabledAudit {
		if err := m.store.InsertAuditLog(ctx, "remove", "license", id); err != nil {
			slog.Error("failed to insert audit log", "licenseID", id, "error", err)
			return fmt.Errorf("failed to insert audit log: %w", err)
		}
	}

	slog.Info("removed license", "licenseID", id)
	return nil
}

func (m *manager) ListLicenses(ctx context.Context) ([]License, error) {
	slog.Debug("Fetching licenses")

	licenses, err := m.store.GetAllLicenses(ctx)
	if err != nil {
		slog.Error("failed to fetch licenses", "error", err)
		return nil, err
	}

	slog.Info("fetched licenses", "count", len(licenses))
	return licenses, nil
}

func (m *manager) GetLicenseByID(ctx context.Context, id string) (License, error) {
	slog.Debug("Fetching license", "id", id)

	license, err := m.store.GetLicenseByID(ctx, id)
	if err != nil {
		slog.Error("failed to fetch license by ID", "licenseID", id, "error", err)
		return License{}, err
	}

	slog.Info("fetched license", "licenseID", id)
	return license, nil
}
