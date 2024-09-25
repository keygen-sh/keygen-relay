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

func setup(env *testscript.Env) error {
	binPath := filepath.Join(env.WorkDir, "relay")
	cmd := exec.Command("go", "build", "-o", binPath, "./cmd/relay")
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	env.Setenv("PATH", env.Getenv("PATH")+string(os.PathListSeparator)+env.WorkDir)

	licenseSrc := filepath.Join(env.Getenv("TESTSCRIPT_ROOT"), "testdata", "license.lic")
	licenseDst := filepath.Join(env.WorkDir, "license.lic")
	input, err := os.ReadFile(licenseSrc)
	if err != nil {
		return fmt.Errorf("failed to read license file: %v", err)
	}
	if err := os.WriteFile(licenseDst, input, 0644); err != nil {
		return fmt.Errorf("failed to write license file: %v", err)
	}

	return nil
}
