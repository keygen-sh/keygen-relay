package licenses

import (
	"context"
	"fmt"
	"github.com/keygen-sh/keygen-go/v3"
	"github.com/keygen-sh/keygen-relay/internal/db"
	"os"
)

type Store interface {
	InsertLicense(ctx context.Context, id string, file []byte, key string) error
	DeleteLicenseByID(ctx context.Context, id string) error
	GetAllLicenses(ctx context.Context) ([]db.License, error)
	GetLicenseByID(ctx context.Context, id string) (db.License, error)
	InsertNode(ctx context.Context, fingerprint string) error
	GetNodeByFingerprint(ctx context.Context, fingerprint string) (db.Node, error)
	UpdateNodeHeartbeat(ctx context.Context, fingerprint string) error
	DeleteNodeByFingerprint(ctx context.Context, fingerprint string) error
	InsertAuditLog(ctx context.Context, action, entityType, entityID string) error
}

type Manager interface {
	AddLicense(ctx context.Context, licenseFilePath string, licenseKey string, publicKeyPath string) error
	RemoveLicense(ctx context.Context, id string) error
	ListLicenses(ctx context.Context) ([]db.License, error)
	GetLicenseByID(ctx context.Context, id string) (db.License, error)
}

type manager struct {
	store  Store
	config *Config
}

func NewManager(store Store, config *Config) Manager {
	return &manager{
		store:  store,
		config: config,
	}
}

func (m *manager) AddLicense(ctx context.Context, licenseFilePath string, licenseKey string, publicKey string) error {
	cert, err := os.ReadFile(licenseFilePath)
	if err != nil {
		return fmt.Errorf("failed to read license file: %w", err)
	}

	keygen.PublicKey = publicKey
	lic := &keygen.LicenseFile{
		Certificate: string(cert),
	}

	if err := lic.Verify(); err != nil {
		return fmt.Errorf("license verification failed: %w", err)
	}

	dec, err := lic.Decrypt(licenseKey)
	if err != nil {
		return fmt.Errorf("license decryption failed: %w", err)
	}

	id := dec.License.ID

	if err := m.store.InsertLicense(ctx, id, cert, licenseKey); err != nil {
		return fmt.Errorf("failed to insert license: %w", err)
	}

	if m.config.EnabledAudit {
		if err := m.store.InsertAuditLog(ctx, "add", "license", id); err != nil {
			return fmt.Errorf("failed to insert audit log: %w", err)
		}
	}

	fmt.Printf("Added license ID: %s\n", id)
	return nil
}

func (m *manager) RemoveLicense(ctx context.Context, id string) error {
	if err := m.store.DeleteLicenseByID(ctx, id); err != nil {
		return fmt.Errorf("failed to remove license: %w", err)
	}

	if m.config.EnabledAudit {
		if err := m.store.InsertAuditLog(ctx, "remove", "license", id); err != nil {
			return fmt.Errorf("failed to insert audit log: %w", err)
		}
	}

	fmt.Printf("Removed license ID: %s\n", id)
	return nil
}

func (m *manager) ListLicenses(ctx context.Context) ([]db.License, error) {
	return m.store.GetAllLicenses(ctx)
}

func (m *manager) GetLicenseByID(ctx context.Context, id string) (db.License, error) {
	return m.store.GetLicenseByID(ctx, id)
}
