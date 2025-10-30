package security

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

// LoadTLS loads a TLS certificate and key from files.
func LoadTLS(caFile, certFile, keyFile string) (*tls.Config, error) {
	// Load CA certificate
	caCert, err := os.ReadFile(caFile)
	if err != nil {
		return nil, err
	}

	certPool, err := LoadCABundle(string(caCert))
	if err != nil {
		return nil, err
	}

	// Load server certificate
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{cert},
		ClientCAs:    certPool,
		RootCAs:      certPool,
	}

	return tlsConfig, nil
}

// LoadSystemCABundle loads system root CA bundle with certs from the given PEMs.
func LoadSystemCABundle(cas ...string) (*x509.CertPool, error) {
	certPool, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("can't load system cert pool: %w", err)
	}

	for _, ca := range cas {
		if ok := certPool.AppendCertsFromPEM([]byte(ca)); !ok {
			return nil, fmt.Errorf("can't parse ca pem '%+v'", ca)
		}
	}

	return certPool, nil
}

// LoadCABundle makes new CA bundle with certs from the given PEMs.
func LoadCABundle(cas ...string) (*x509.CertPool, error) {
	certPool := x509.NewCertPool()

	for _, ca := range cas {
		if ok := certPool.AppendCertsFromPEM([]byte(ca)); !ok {
			return nil, fmt.Errorf("can't parse ca pem")
		}
	}

	return certPool, nil
}
