//go:build integration
// +build integration

package cli_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/keygen-sh/keygen-relay/cli"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	code := testscript.RunMain(m, map[string]func() int{
		"relay": func() int {
			return cli.Run()
		},
	})

	os.Exit(code)
}

func TestIntegration(t *testing.T) {
	t.Parallel()

	testscript.Run(t, testscript.Params{
		Dir:                 "testdata",
		RequireExplicitExec: true,
		TestWork:            true,
		Setup:               setup,
	})
}

func setup(env *testscript.Env) error {
	setupFixtures(env)
	setupEnv(env)

	return nil
}

func setupFixtures(env *testscript.Env) error {
	fixtures := []string{
		"license.lic",
		"license_2.lic",
	}

	for _, fixture := range fixtures {
		if err := copyFile(filepath.Join("testdata", fixture), filepath.Join(env.WorkDir, fixture)); err != nil {
			return fmt.Errorf("failed to copy fixture %s: %w", fixture, err)
		}
	}

	return nil
}

func setupEnv(env *testscript.Env) error {
	// TODO(ezekg) make prestine env
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
