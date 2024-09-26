//go:build integrity
// +build integrity

package main_test

import (
	"fmt"
	"github.com/rogpeppe/go-internal/testscript"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestRelayIntegration(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir:   "testdata",
		Setup: setup,
	})
}

// setup prepares the test environment by building the relay binary and copying required license files
func setup(env *testscript.Env) error {
	binPath := filepath.Join(env.WorkDir, "relay")
	buildCmd := exec.Command("go", "build", "-o", binPath, "./cmd/relay")
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("failed to build relay binary: %w", err)
	}

	env.Setenv("PATH", env.Getenv("PATH")+string(os.PathListSeparator)+env.WorkDir)

	// copy license files into the working directory tests
	testscriptRoot := env.Getenv("TESTSCRIPT_ROOT")
	if err := copyFile(filepath.Join(testscriptRoot, "testdata", "license.lic"), filepath.Join(env.WorkDir, "license.lic")); err != nil {
		return fmt.Errorf("failed to copy license.lic: %w", err)
	}
	if err := copyFile(filepath.Join(testscriptRoot, "testdata", "license_2.lic"), filepath.Join(env.WorkDir, "license_2.lic")); err != nil {
		return fmt.Errorf("failed to copy license_2.lic: %w", err)
	}

	return nil
}

func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("unable to read file %s: %w", src, err)
	}

	if err := os.WriteFile(dst, input, 0644); err != nil {
		return fmt.Errorf("unable to write file to %s: %w", dst, err)
	}

	return nil
}
