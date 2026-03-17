package certgen

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"time"
)

// Input contains all user-provided values used to generate the certificates.
type Input struct {
	Domain string
	IP     string
	Org    string
}

// Output contains the generated files, stored in memory.
type Output struct {
	// Files is a map of relative file paths to bytes, for example:
	// "certs/ca.crt" -> PEM bytes
	Files map[string][]byte
}

// Generate creates an internal CA and a server certificate, returning all files in memory.
// This is WASM-friendly (no filesystem usage).
func Generate(in Input) (Output, error) {
	// ------------------------------------------------
	// CA
	// ------------------------------------------------
	caKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return Output{}, err
	}

	caTemplate := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   "Internal CA",
			Organization: []string{in.Org},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}

	caDER, err := x509.CreateCertificate(rand.Reader, &caTemplate, &caTemplate, &caKey.PublicKey, caKey)
	if err != nil {
		return Output{}, err
	}

	caCert, err := x509.ParseCertificate(caDER)
	if err != nil {
		return Output{}, err
	}

	// ------------------------------------------------
	// SERVER CERT
	// ------------------------------------------------
	serverKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return Output{}, err
	}

	var ips []net.IP
	parsedIP := net.ParseIP(in.IP)
	if parsedIP != nil {
		ips = append(ips, parsedIP)
	}

	serverTemplate := x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			CommonName:   in.Domain,
			Organization: []string{in.Org},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(5, 0, 0),

		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},

		DNSNames:    []string{in.Domain},
		IPAddresses: ips,
	}

	serverDER, err := x509.CreateCertificate(rand.Reader, &serverTemplate, caCert, &serverKey.PublicKey, caKey)
	if err != nil {
		return Output{}, err
	}

	files := map[string][]byte{
		"certs/ca.crt":         pemCert(caDER),
		"certs/ca.key":         pemRSAKey(caKey),
		"certs/server.crt":     pemCert(serverDER),
		"certs/server.key":     pemRSAKey(serverKey),
		"certs/fullchain.crt": pemFullchain(serverDER, caDER),
		// Helper installer script lives next to ca.crt inside certs/.
		"certs/install_ca_cert.sh": installCACertScript,
	}

	return Output{Files: files}, nil
}

// pemCert encodes a DER certificate as PEM.
func pemCert(der []byte) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: der,
	})
}

// pemRSAKey encodes an RSA private key as PEM (PKCS#1).
func pemRSAKey(key *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
}

// pemFullchain builds a "fullchain" PEM (server cert + CA cert).
func pemFullchain(serverDER, caDER []byte) []byte {
	serverPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: serverDER})
	caPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caDER})
	out := make([]byte, 0, len(serverPEM)+len(caPEM))
	out = append(out, serverPEM...)
	out = append(out, caPEM...)
	return out
}

