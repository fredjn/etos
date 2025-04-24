package events

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/eiffel-community/eiffelevents-sdk-go"
	"github.com/eiffel-community/eiffelevents-sdk-go/signature"
)

// Signer handles signing of Eiffel events
type Signer struct {
	privateKey *rsa.PrivateKey
}

// NewSigner creates a new signer from a private key file
func NewSigner(privateKeyPath string) (*Signer, error) {
	// Read the private key file
	keyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	// Decode the PEM block
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	// Parse the private key
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return &Signer{
		privateKey: privateKey,
	}, nil
}

// SignEvent signs an Eiffel event
func (s *Signer) SignEvent(event *eiffelevents.Any) error {
	// Create a signer with the private key
	signer := signature.NewSigner(s.privateKey, crypto.SHA256)

	// Sign the event
	if err := signer.Sign(event); err != nil {
		return fmt.Errorf("failed to sign event: %w", err)
	}

	return nil
}

// VerifyEvent verifies the signature of an Eiffel event
func VerifyEvent(event *eiffelevents.Any, publicKey *rsa.PublicKey) error {
	// Create a verifier with the public key
	verifier := signature.NewVerifier(publicKey, crypto.SHA256)

	// Verify the event signature
	if err := verifier.Verify(event); err != nil {
		return fmt.Errorf("failed to verify event signature: %w", err)
	}

	return nil
} 