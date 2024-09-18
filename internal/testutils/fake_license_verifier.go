package testutils

import "github.com/keygen-sh/keygen-go/v3"

type FakeLicenseVerifier struct{}

func (m *FakeLicenseVerifier) Verify() error {
	return nil
}

func (m *FakeLicenseVerifier) Decrypt(_ string) (*keygen.LicenseFileDataset, error) {
	return &keygen.LicenseFileDataset{
		License: keygen.License{
			ID:  "fake_id",
			Key: "fake_key",
		},
	}, nil
}
