package certs

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// createKey will create a new rsa key with length.
func createKey(length int) (*rsa.PrivateKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, length)
	if err != nil {
		return nil, fmt.Errorf("couldn't create new key. %w", err)
	}

	return key, nil
}

// writeKey will write key as a pem file fn.
func writeKey(key *rsa.PrivateKey, fn string) error {
	writer, err := os.OpenFile(fn, fileFlags, fileMode)
	if err != nil {
		return fmt.Errorf("couldn't create file %q. %w", fn, err)
	}
	defer writer.Close()

	if err := pem.Encode(
		writer,
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	); err != nil {
		return fmt.Errorf("couldn't write key file %q. %w", fn, err)
	}

	return nil
}
