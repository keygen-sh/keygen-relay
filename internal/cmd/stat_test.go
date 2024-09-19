package cmd_test

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"github.com/keygen-sh/keygen-relay/internal/licenses"
	"testing"

	"github.com/keygen-sh/keygen-relay/internal/cmd"
	"github.com/keygen-sh/keygen-relay/internal/testutils"
	"github.com/keygen-sh/keygen-relay/internal/ui"
	"github.com/stretchr/testify/assert"
)

func TestStatCmd_Success(t *testing.T) {
	manager := &testutils.FakeManager{
		GetLicenseByIDFn: func(ctx context.Context, id string) (licenses.License, error) {
			return licenses.License{
				ID:             "License_1",
				Key:            "License_Key1",
				Claims:         5,
				NodeID:         sql.NullInt64{Int64: 123, Valid: true},
				LastClaimedAt:  sql.NullString{String: "2024-01-01T00:00:00Z", Valid: true},
				LastReleasedAt: sql.NullString{String: "2024-01-05T10:00:00Z", Valid: true},
			}, nil
		},
	}

	outBuf := new(bytes.Buffer)
	renderer := ui.NewSimpleTableRenderer(outBuf)

	statCmd := cmd.StatCmd(manager, renderer)
	statCmd.SetOut(outBuf)
	statCmd.SetArgs([]string{"--id=License_1"})

	err := statCmd.Execute()
	assert.NoError(t, err)

	assert.Contains(t, outBuf.String(), "License_1")
	assert.Contains(t, outBuf.String(), "123")
	assert.Contains(t, outBuf.String(), "2024-01-01 00:00:00")
	assert.Contains(t, outBuf.String(), "2024-01-05 10:00:00")
}

func TestStatCmd_MissingFlag(t *testing.T) {
	manager := &testutils.FakeManager{}

	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)

	renderer := ui.NewSimpleTableRenderer(outBuf)

	statCmd := cmd.StatCmd(manager, renderer)
	statCmd.SetOut(outBuf)
	statCmd.SetErr(errBuf)

	err := statCmd.Execute()
	assert.Error(t, err)

	assert.Contains(t, errBuf.String(), `required flag(s) "id" not set`)
}

func TestStatCmd_Error(t *testing.T) {
	manager := &testutils.FakeManager{
		GetLicenseByIDFn: func(ctx context.Context, id string) (licenses.License, error) {
			return licenses.License{}, errors.New("license not found")
		},
	}

	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)

	renderer := ui.NewSimpleTableRenderer(outBuf)

	statCmd := cmd.StatCmd(manager, renderer)
	statCmd.SetOut(outBuf)
	statCmd.SetErr(errBuf)
	statCmd.SetArgs([]string{"--id=invalid"})

	err := statCmd.Execute()
	assert.Error(t, err)

	assert.Contains(t, errBuf.String(), "license not found")
}
