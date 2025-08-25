package cmd_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/keygen-sh/keygen-relay/internal/cmd"
	"github.com/keygen-sh/keygen-relay/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestDelCmd_Success(t *testing.T) {
	manager := &testutils.FakeManager{
		RemoveLicenseFn: func(ctx context.Context, pool *string, id string) error {
			return nil
		},
	}

	delCmd := cmd.DelCmd(manager)

	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	delCmd.SetOut(outBuf)
	delCmd.SetErr(errBuf)

	delCmd.SetArgs([]string{"--license=test-id"})

	err := delCmd.Execute()
	assert.NoError(t, err)
	assert.Empty(t, errBuf.String())
	assert.Contains(t, outBuf.String(), "license deleted successfully: test-id")
}

func TestDelCmd_MultiSuccess(t *testing.T) {
	manager := &testutils.FakeManager{
		RemoveLicenseFn: func(ctx context.Context, pool *string, id string) error {
			return nil
		},
	}

	delCmd := cmd.DelCmd(manager)

	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	delCmd.SetOut(outBuf)
	delCmd.SetErr(errBuf)

	delCmd.SetArgs([]string{"--license=1", "--license=2"})

	err := delCmd.Execute()
	assert.NoError(t, err)
	assert.Empty(t, errBuf.String())
	assert.Contains(t, outBuf.String(), "license deleted successfully: 1")
	assert.Contains(t, outBuf.String(), "license deleted successfully: 2")
}

func TestDelCmd_Error(t *testing.T) {
	manager := &testutils.FakeManager{
		RemoveLicenseFn: func(ctx context.Context, pool *string, id string) error {
			return errors.New("failed to remove license")
		},
	}

	delCmd := cmd.DelCmd(manager)

	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	delCmd.SetOut(outBuf)
	delCmd.SetErr(errBuf)

	delCmd.SetArgs([]string{"--license=test-id"})

	_ = delCmd.Execute()

	assert.Contains(t, errBuf.String(), "error: failed to remove license")
}

func TestDelCmd_MissingID(t *testing.T) {
	manager := &testutils.FakeManager{}

	delCmd := cmd.DelCmd(manager)

	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	delCmd.SetOut(outBuf)
	delCmd.SetErr(errBuf)

	delCmd.SetArgs([]string{})

	err := delCmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, errBuf.String(), `required flag(s) "license" not set`)
}
