package cmd_test

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/keygen-sh/keygen-relay/internal/cmd"
	"github.com/keygen-sh/keygen-relay/internal/db"
	"github.com/keygen-sh/keygen-relay/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestStatCmd_Success(t *testing.T) {
	nodeID := int64(123)
	lastClaimedAt, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z")
	lastReleasedAt, _ := time.Parse(time.RFC3339, "2024-01-05T10:00:00Z")

	manager := &testutils.FakeManager{
		GetLicenseByIDFn: func(ctx context.Context, id string) (*db.License, error) {
			return &db.License{
				ID:             "License_1",
				Key:            "License_Key1",
				Claims:         5,
				NodeID:         &nodeID,
				LastClaimedAt:  &lastClaimedAt,
				LastReleasedAt: &lastReleasedAt,
			}, nil
		},
	}

	outBuf := new(bytes.Buffer)

	statCmd := cmd.StatCmd(manager)
	statCmd.SetOut(outBuf)
	statCmd.SetArgs([]string{"--license=License_1", "--plain"})

	err := statCmd.Execute()
	assert.NoError(t, err)

	assert.Contains(t, outBuf.String(), "License_1")
	assert.Contains(t, outBuf.String(), "123")
	assert.Contains(t, outBuf.String(), "2024-01-01T00:00:00Z")
	assert.Contains(t, outBuf.String(), "2024-01-05T10:00:00Z")
}

func TestStatCmd_MissingFlag(t *testing.T) {
	manager := &testutils.FakeManager{}

	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)

	statCmd := cmd.StatCmd(manager)
	statCmd.SetOut(outBuf)
	statCmd.SetErr(errBuf)

	err := statCmd.Execute()
	assert.Error(t, err)

	assert.Contains(t, errBuf.String(), `required flag(s) "license" not set`)
}

func TestStatCmd_Error(t *testing.T) {
	manager := &testutils.FakeManager{
		GetLicenseByIDFn: func(ctx context.Context, id string) (*db.License, error) {
			return nil, errors.New("license not found")
		},
	}

	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)

	statCmd := cmd.StatCmd(manager)
	statCmd.SetOut(outBuf)
	statCmd.SetErr(errBuf)
	statCmd.SetArgs([]string{"--license=invalid", "--plain"})

	_ = statCmd.Execute()

	assert.Contains(t, errBuf.String(), "license not found")
}
