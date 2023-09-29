/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 *
 * Establishes a Watcher for the Kubernetes Secrets that contain the various certificates and keys used to generate a tls.Config object;
 * exposes the certificates and keys.
 */

package certification

type Certificates struct {
	CACertificate     string
	ClientCertificate string
	ClientKey         string
}

func NewCertificates() (*Certificates, error) {
	return &Certificates{}, nil
}

// GetCACertificate returns the Certificate Authority certificate.
func (c *Certificates) GetCACertificate() []byte {
	return []byte(c.CACertificate)
}

// GetClientCertificate returns the Client certificate and key.
func (c *Certificates) GetClientCertificate() ([]byte, []byte) {
	clientKey := []byte(c.ClientKey)
	clientCertificate := []byte(c.ClientCertificate)

	return clientKey, clientCertificate
}
