package cmd_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/keygen-sh/keygen-relay/internal/cmd"
	"github.com/keygen-sh/keygen-relay/internal/db"
	"github.com/keygen-sh/keygen-relay/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestAddCmd_Success(t *testing.T) {
	manager := &testutils.FakeManager{
		AddLicenseFn: func(ctx context.Context, filePath, key, publicKey string) (*db.License, error) {
			return &db.License{Guid: "test" + key}, nil
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

	assert.Contains(t, outBuf.String(), "license added successfully: testkey")
}

func TestAddCmd_Error(t *testing.T) {
	manager := &testutils.FakeManager{
		AddLicenseFn: func(ctx context.Context, filePath, key, publicKey string) (*db.License, error) {
			return nil, errors.New("failed to add license")
		},
	}

	addCmd := cmd.AddCmd(manager)

	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	addCmd.SetOut(outBuf)
	addCmd.SetErr(errBuf)

	addCmd.SetArgs([]string{"--file=file.lic", "--key=testkey", "--public-key=testpublickey"})

	err := addCmd.Execute()
	assert.NoError(t, err)

	assert.Contains(t, errBuf.String(), "error: failed to add license")
}

func TestAdd_MultiSuccess(t *testing.T) {
	manager := &testutils.FakeManager{
		AddLicenseFn: func(ctx context.Context, filePath, key, publicKey string) (*db.License, error) {
			return &db.License{Guid: "test" + key}, nil
		},
	}

	addCmd := cmd.AddCmd(manager)

	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	addCmd.SetOut(outBuf)
	addCmd.SetErr(errBuf)

	addCmd.SetArgs([]string{"--file=1.lic", "--key=1", "--file=2.lic", "--key=2", "--public-key=testpublickey"})

	err := addCmd.Execute()
	assert.NoError(t, err)

	assert.Contains(t, outBuf.String(), "license added successfully: test1")
	assert.Contains(t, outBuf.String(), "license added successfully: test2")
}

func TestAdd_MultiError(t *testing.T) {
	manager := &testutils.FakeManager{
		AddLicenseFn: func(ctx context.Context, filePath, key, publicKey string) (*db.License, error) {
			return &db.License{Guid: "test_" + key}, nil
		},
	}

	addCmd := cmd.AddCmd(manager)

	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	addCmd.SetOut(outBuf)
	addCmd.SetErr(errBuf)

	addCmd.SetArgs([]string{"--file=file1.lic", "--key=key1", "--file=file2.lic", "--public-key=testpublickey"})

	err := addCmd.Execute()
	assert.NoError(t, err)

	assert.Contains(t, errBuf.String(), "error: number of key and file flags must match")
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
