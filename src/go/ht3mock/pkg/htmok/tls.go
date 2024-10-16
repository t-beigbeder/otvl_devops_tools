package htmok

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

func getTlsClientConfig(config *Config) (*tls.Config, error) {
	var certs []tls.Certificate
	var err error
	if config.CertFile != "" && config.KeyFile != "" {
		scert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("in getTlsClientConfig: %v", err)
		}
		certs = append(certs, scert)
	}
	if config.CCertFile == "" {
		return &tls.Config{InsecureSkipVerify: config.UnsafeTls, Certificates: certs}, nil
	}
	cb, err := os.ReadFile(config.CCertFile)
	if err != nil {
		return nil, fmt.Errorf("in getTlsClientConfig: %v", err)
	}
	caCertPool, _ := x509.SystemCertPool()
	ok := caCertPool.AppendCertsFromPEM(cb)
	if !ok {
		return nil, fmt.Errorf("in AppendCertsFromPEM: nok")
	}
	return &tls.Config{RootCAs: caCertPool, InsecureSkipVerify: config.UnsafeTls, Certificates: certs}, nil
}
