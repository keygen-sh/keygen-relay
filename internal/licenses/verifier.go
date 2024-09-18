package licenses

import "github.com/keygen-sh/keygen-go/v3"

type LicenseVerifier interface {
	Verify() error
	Decrypt(key string) (*keygen.LicenseFileDataset, error)
}

type KeygenLicenseVerifier struct {
	lic *keygen.LicenseFile
}

func NewKeygenLicenseVerifier(cert []byte) LicenseVerifier {
	return &KeygenLicenseVerifier{
		lic: &keygen.LicenseFile{
			Certificate: string(cert),
		},
	}
}

func (k *KeygenLicenseVerifier) Verify() error {
	return k.lic.Verify()
}

func (k *KeygenLicenseVerifier) Decrypt(key string) (*keygen.LicenseFileDataset, error) {
	return k.lic.Decrypt(key)
}
