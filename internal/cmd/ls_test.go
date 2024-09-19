package cmd_test

import (
	"bytes"
	"context"
	"errors"
	"github.com/keygen-sh/keygen-relay/internal/licenses"
	"testing"

	"github.com/keygen-sh/keygen-relay/internal/cmd"
	"github.com/keygen-sh/keygen-relay/internal/testutils"
	"github.com/keygen-sh/keygen-relay/internal/ui"
	"github.com/stretchr/testify/assert"
)

func TestLsCmd_Success(t *testing.T) {
	manager := &testutils.FakeManager{
		ListLicensesFn: func(ctx context.Context) ([]licenses.License, error) {
			return []licenses.License{
				{ID: "License_1", Key: "License_Key_1", Claims: 5},
				{ID: "License_2", Key: "License_Key_2", Claims: 10},
			}, nil
		},
	}

	outBuf := new(bytes.Buffer)
	renderer := ui.NewSimpleTableRenderer(outBuf)

	lsCmd := cmd.LsCmd(manager, renderer)
	lsCmd.SetOut(outBuf)

	err := lsCmd.Execute()
	assert.NoError(t, err)

	assert.Contains(t, outBuf.String(), "License_Key_1")
	assert.Contains(t, outBuf.String(), "License_Key_2")
}

func TestLsCmd_NoLicenses(t *testing.T) {
	manager := &testutils.FakeManager{
		ListLicensesFn: func(ctx context.Context) ([]licenses.License, error) {
			return []licenses.License{}, nil
		},
	}

	outBuf := new(bytes.Buffer)
	renderer := ui.NewSimpleTableRenderer(outBuf)

	lsCmd := cmd.LsCmd(manager, renderer)
	lsCmd.SetOut(outBuf)

	err := lsCmd.Execute()
	assert.NoError(t, err)

	assert.Contains(t, outBuf.String(), "No licenses found.")
}

func TestLsCmd_Error(t *testing.T) {
	manager := &testutils.FakeManager{
		ListLicensesFn: func(ctx context.Context) ([]licenses.License, error) {
			return nil, errors.New("failed to list licenses")
		},
	}

	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)

	renderer := ui.NewSimpleTableRenderer(outBuf)

	lsCmd := cmd.LsCmd(manager, renderer)
	lsCmd.SetOut(outBuf)
	lsCmd.SetErr(errBuf)

	err := lsCmd.Execute()
	assert.Error(t, err)

	assert.Contains(t, errBuf.String(), "failed to list licenses")
}
