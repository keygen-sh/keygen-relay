package server

import (
	"crypto/hmac"
	"crypto/sha256"
)

type Signer struct {
	secret []byte
}

func NewSigner(cfg *Config) *Signer {
	if cfg.SigningSecret != nil {
		return &Signer{secret: []byte(*cfg.SigningSecret)}
	}

	return &Signer{secret: nil}
}

// Sign generates an HMAC-SHA256 signature for the given message
func (s *Signer) Sign(message []byte) []byte {
	mac := hmac.New(sha256.New, s.secret)
	mac.Write(message)

	return mac.Sum(nil)
}

// Enabled returns true if the signer has a secret configured
func (s *Signer) Enabled() bool {
	return len(s.secret) > 0
}
