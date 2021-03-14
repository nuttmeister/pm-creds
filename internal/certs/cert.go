package certs

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"math"
	"math/big"
	"os"
	"time"
)

// createCACert will create a CA certificate from key.
// Both the bytes and parsed cert is returned.
func createCACert(key *rsa.PrivateKey) ([]byte, *x509.Certificate, error) {
	sn, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't create serial number. %w", err)
	}

	skid, err := createSKID(key.Public())
	if err != nil {
		return nil, nil, err
	}

	template := &x509.Certificate{
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLenZero:        true,

		Subject:      pkix.Name{CommonName: "pm-creds ca"},
		SerialNumber: sn,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(100, 0, 0),

		SubjectKeyId:   skid,
		AuthorityKeyId: skid,
		KeyUsage:       x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
	}

	raw, err := x509.CreateCertificate(rand.Reader, template, template, key.Public(), key)
	if err != nil {
		return nil, nil, err
	}

	cert, err := x509.ParseCertificate(raw)
	if err != nil {
		return nil, nil, err
	}

	return raw, cert, nil
}

// createServerCert creates a new server cert and signs it with caCert and caKey.
// Both bytes and parsed cert is returned.
func createServerCert(key crypto.Signer, caKey crypto.Signer, caCert *x509.Certificate) ([]byte, *x509.Certificate, error) {
	sn, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return nil, nil, err
	}

	template := &x509.Certificate{
		BasicConstraintsValid: true,
		IsCA:                  false,

		Subject:      pkix.Name{CommonName: "localhost"},
		SerialNumber: sn,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),

		KeyUsage: x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
	}

	raw, err := x509.CreateCertificate(rand.Reader, template, caCert, key.Public(), caKey)
	if err != nil {
		return nil, nil, err
	}

	cert, err := x509.ParseCertificate(raw)
	if err != nil {
		return nil, nil, err
	}

	return raw, cert, nil
}

// createSKID will create a sha1 subject public key info from public key.
func createSKID(key crypto.PublicKey) ([]byte, error) {
	raw, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return nil, fmt.Errorf("couldn't create pkix from public key. %w", err)
	}

	data := &struct {
		Algorithm        pkix.AlgorithmIdentifier
		SubjectPublicKey asn1.BitString
	}{}

	if _, err := asn1.Unmarshal(raw, data); err != nil {
		return nil, fmt.Errorf("couldn't asn1 unmarshal pkix public key. %w", err)
	}

	res := sha1.Sum(data.SubjectPublicKey.Bytes)
	return res[:], nil
}

// writeCACert will write ca cert as a pem file fn.
func writeCert(raw []byte, fn string) error {
	writer, err := os.OpenFile(fn, fileFlags, fileMode)
	if err != nil {
		return fmt.Errorf("couldn't create file %q. %w", fn, err)
	}
	defer writer.Close()

	if err := pem.Encode(
		writer,
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: raw,
		},
	); err != nil {
		return fmt.Errorf("couldn't write cert file %q. %w", fn, err)
	}

	return nil
}
