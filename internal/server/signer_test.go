package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSigner_WithSecret(t *testing.T) {
	secret := "hunter2"
	cfg := &Config{SigningSecret: &secret}

	signer := NewSigner(cfg)

	assert.NotNil(t, signer)
	assert.True(t, signer.Enabled())
}

func TestNewSigner_WithoutSecret(t *testing.T) {
	cfg := &Config{SigningSecret: nil}

	signer := NewSigner(cfg)

	assert.NotNil(t, signer)
	assert.False(t, signer.Enabled())
}

func TestSigner_EmptySecret(t *testing.T) {
	empty := ""
	cfg := &Config{SigningSecret: &empty}

	signer := NewSigner(cfg)

	assert.NotNil(t, signer)
	assert.False(t, signer.Enabled())
}

func TestSigner_Sign_WithSecret(t *testing.T) {
	secret := "hunter2"
	cfg := &Config{SigningSecret: &secret}
	signer := NewSigner(cfg)

	sig := signer.Sign([]byte("message"))

	assert.Len(t, sig, 32)
}

func TestSigner_Sign_WithoutSecret(t *testing.T) {
	cfg := &Config{SigningSecret: nil}
	signer := NewSigner(cfg)

	sig := signer.Sign([]byte("message"))

	assert.Len(t, sig, 32)
}

func TestSigner_Sign_DifferentMessages(t *testing.T) {
	secret := "hunter2"
	cfg := &Config{SigningSecret: &secret}
	signer := NewSigner(cfg)

	sig1 := signer.Sign([]byte("message one"))
	sig2 := signer.Sign([]byte("message two"))

	assert.NotEqual(t, sig1, sig2)
}

func TestSigner_Sign_DifferentSecrets(t *testing.T) {
	secret1 := "hunter2"
	secret2 := "****"
	cfg1 := &Config{SigningSecret: &secret1}
	cfg2 := &Config{SigningSecret: &secret2}

	signer1 := NewSigner(cfg1)
	signer2 := NewSigner(cfg2)

	message := []byte("message")
	sig1 := signer1.Sign(message)
	sig2 := signer2.Sign(message)

	assert.NotEqual(t, sig1, sig2)
}
