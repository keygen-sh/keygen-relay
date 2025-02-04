package licenses

import (
	"github.com/keygen-sh/keygen-go/v3"
)

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

type MachineVerifier interface {
	Verify() error
	Decrypt(fingerprint string, key string) (*keygen.MachineFileDataset, error)
}

type KeygenMachineVerifier struct {
	lic *keygen.MachineFile
}

func NewKeygenMachineVerifier(cert []byte) MachineVerifier {
	return &KeygenMachineVerifier{
		lic: &keygen.MachineFile{
			Certificate: string(cert),
		},
	}
}

func (k *KeygenMachineVerifier) Verify() error {
	return k.lic.Verify()
}

func (k *KeygenMachineVerifier) Decrypt(key string, fingerprint string) (*keygen.MachineFileDataset, error) {
	return k.lic.Decrypt(key + fingerprint)
}
