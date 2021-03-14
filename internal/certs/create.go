package certs

import (
	"crypto"
	"crypto/x509"
	"os"
)

const (
	fileFlags = os.O_CREATE | os.O_WRONLY
	fileMode  = 0600
	keyBits   = 4096
)

func Create(fnCAKey string, fnCACert string, fnServerKey string, fnServerCert string) error {
	caKey, caCert, err := createCA(fnCAKey, fnCACert)
	if err != nil {
		return err
	}

	return createServer(caKey, caCert, fnServerKey, fnServerCert)
}

// createCA will create a CA key and certificate
// and save them to fnKey and fnCert.
func createCA(fnKey string, fnCert string) (crypto.Signer, *x509.Certificate, error) {
	key, err := createKey(keyBits)
	if err != nil {
		return nil, nil, err
	}

	if err := writeKey(key, fnKey); err != nil {
		return nil, nil, err
	}

	raw, cert, err := createCACert(key)
	if err != nil {
		return nil, nil, err
	}

	if err := writeCert(raw, fnCert); err != nil {
		return nil, nil, err
	}

	return key, cert, nil
}

// createServer will create a server key and certificate and sign
// it with caKey and caCert and save them to fnKey and fnCert
func createServer(caKey crypto.Signer, caCert *x509.Certificate, fnKey string, fnCert string) error {
	key, err := createKey(keyBits)
	if err != nil {
		return err
	}

	if err := writeKey(key, fnKey); err != nil {
		return err
	}

	raw, _, err := createServerCert(key, caKey, caCert)
	if err != nil {
		return err
	}

	if err := writeCert(raw, fnCert); err != nil {
		return err
	}

	return nil
}
