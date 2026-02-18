package testutils

import "github.com/keygen-sh/keygen-go/v3"

type FakeLicenseVerifier struct {
	LicenseID       string
	LicenseKey      string
	LicenseMetadata map[string]interface{}
}

func (f *FakeLicenseVerifier) Verify() error {
	return nil
}

func (f *FakeLicenseVerifier) Decrypt(key string) (*keygen.LicenseFileDataset, error) {
	licenseID := "license_" + key
	f.LicenseID = licenseID
	f.LicenseKey = key

	return &keygen.LicenseFileDataset{
		License: keygen.License{
			ID:       licenseID,
			Key:      key,
			Metadata: f.LicenseMetadata,
		},
	}, nil
}
