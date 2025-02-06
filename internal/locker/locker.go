package locker

import (
	"fmt"
	"os"
	"runtime"

	"github.com/keygen-sh/keygen-go/v3"
	"github.com/keygen-sh/keygen-relay/internal/licenses"
	"github.com/keygen-sh/machineid"
)

// Machine attributes for node-locking Relay, embedded at compile time. When set, this
// locks Relay to a specific machine, depending on provided attributes. Relay will
// error on mismatch, e.g. underlying IP address is different than expected IP.
var (
	PublicKey   string // required
	Fingerprint string // required
	Platform    string // optional
	Hostname    string // optional
	IP          string // optional
	Addr        string // optional
	Port        string // optional
)

func init() {
	// FIXME(ezekg) add support for non-global config to SDK
	keygen.PublicKey = PublicKey
}

// Locked returns a boolean whether or not Relay is node-locked
func Locked() bool {
	return Fingerprint != "" // decent proxy for node-locked
}

// LockedAddr returns a boolean whether or not Relay's bind address is locked
func LockedAddr() bool {
	return Addr != ""
}

// LockedPort returns a boolean whether or not Relay's port is locked
func LockedPort() bool {
	return Port != ""
}

// Unlock attempts to unlock Relay via a machine file and license key using the
// current machine's fingerprint
func Unlock(config Config) (*keygen.MachineFileDataset, error) {
	path, key := config.MachineFilePath, config.LicenseKey
	if path == "" {
		return nil, fmt.Errorf("machine file path is required")
	}

	if key == "" {
		return nil, fmt.Errorf("license key is required")
	}

	fingerprint, err := machineid.ProtectedID("keygen-relay")
	if err != nil {
		return nil, fmt.Errorf("machine could not determine fingerprint: %w", err)
	}

	cert, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("machine file could not be read: %w", err)
	}

	verifier := licenses.NewKeygenMachineVerifier(cert)
	err = verifier.Verify()
	switch {
	case err == keygen.ErrLicenseFileNotGenuine:
		return nil, fmt.Errorf("machine file is not genuine: %w", err)
	case err != nil:
		return nil, fmt.Errorf("machine file could not be verified: %w", err)
	}

	dataset, err := verifier.Decrypt(key, fingerprint)
	switch {
	case err == keygen.ErrSystemClockUnsynced:
		return nil, fmt.Errorf("machine file is desynced with system clock: %w", err)
	case err == keygen.ErrLicenseFileExpired:
		return nil, fmt.Errorf("machine file is expired: %w", err)
	case err != nil:
		return nil, fmt.Errorf("machine file could not be decrypted: %w", err)
	}

	if expected, actual := Fingerprint, fingerprint; dataset.Machine.Fingerprint != expected || actual != expected {
		return nil, fmt.Errorf("machine file fingerprint mismatch")
	}

	if Platform != "" {
		platform := runtime.GOOS + "/" + runtime.GOARCH

		if expected, actual := Platform, platform; dataset.Machine.Platform != expected || actual != expected {
			return nil, fmt.Errorf("machine file platform mismatch")
		}
	}

	if Hostname != "" {
		hostname, err := os.Hostname()
		if err != nil {
			return nil, fmt.Errorf("machine could not determine hostname: %w", err)
		}

		if expected, actual := Hostname, hostname; dataset.Machine.Hostname != expected || actual != expected {
			return nil, fmt.Errorf("machine file hostname mismatch")
		}
	}

	if IP != "" {
		ip, err := getPrivateIP()
		if err != nil {
			return nil, fmt.Errorf("machine could not determine ip: %w", err)
		}

		if expected, actual := IP, ip; dataset.Machine.IP != expected || actual != expected {
			return nil, fmt.Errorf("machine file ip mismatch")
		}
	}

	return dataset, nil
}
