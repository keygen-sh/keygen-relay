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

func TestAddCmd_Success(t *testing.T) {
	manager := &testutils.FakeManager{
		AddLicenseFn: func(ctx context.Context, filePath, key, publicKey string) error {
			return nil
		},
	}

	addCmd := cmd.AddCmd(manager)

	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	addCmd.SetOut(outBuf)
	addCmd.SetErr(errBuf)

	addCmd.SetArgs([]string{"--file=file.lic", "--key=key", "--public-key=testpublickey"})

	err := addCmd.Execute()
	assert.NoError(t, err)

	assert.Contains(t, outBuf.String(), "License added successfully.")
}

func TestAddCmd_Error(t *testing.T) {
	manager := &testutils.FakeManager{
		AddLicenseFn: func(ctx context.Context, filePath, key, publicKey string) error {
			return errors.New("failed to add license")
		},
	}

	addCmd := cmd.AddCmd(manager)

	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	addCmd.SetOut(outBuf)
	addCmd.SetErr(errBuf)

	addCmd.SetArgs([]string{"--file=file.lic", "--key=testkey", "--public-key=testpublickey"})

	_ = addCmd.Execute()

	assert.Contains(t, errBuf.String(), "Error: failed to add license")
}

func TestAddCmd_MissingRequiredFlags(t *testing.T) {
	manager := &testutils.FakeManager{}

	addCmd := cmd.AddCmd(manager)

	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	addCmd.SetOut(outBuf)
	addCmd.SetErr(errBuf)

	err := addCmd.Execute()
	assert.Error(t, err)

	assert.Contains(t, errBuf.String(), `required flag(s) "file", "key", "public-key" not set`)
}
