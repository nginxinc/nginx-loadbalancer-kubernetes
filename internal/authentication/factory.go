/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 *
 * Factory for creating tls.Config objects based on the provided `tls-mode`.
 */

package authentication

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/nginxinc/kubernetes-nginx-ingress/internal/certification"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/configuration"
	"github.com/sirupsen/logrus"
)

func NewTlsConfig(settings *configuration.Settings) (*tls.Config, error) {
	logrus.Debugf("authentication::NewTlsConfig Creating TLS config for mode: '%s'", settings.TlsMode)
	switch settings.TlsMode {

	case configuration.NoTLS:
		return buildBasicTlsConfig(true), nil

	case configuration.SelfSignedTLS: // needs ca cert
		return buildSelfSignedTlsConfig(settings.Certificates)

	case configuration.SelfSignedMutualTLS: // needs ca cert and client cert
		return buildSelfSignedMtlsConfig(settings.Certificates)

	case configuration.CertificateAuthorityTLS: // needs nothing
		return buildBasicTlsConfig(false), nil

	case configuration.CertificateAuthorityMutualTLS: // needs client cert
		return buildCaTlsConfig(settings.Certificates)

	default:
		return nil, fmt.Errorf("unknown TLS mode: %s", settings.TlsMode)
	}
}

func buildSelfSignedTlsConfig(certificates *certification.Certificates) (*tls.Config, error) {
	logrus.Debug("authentication::buildSelfSignedTlsConfig Building self-signed TLS config")
	certPool, err := buildCaCertificatePool(certificates.GetCACertificate())
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		InsecureSkipVerify: false,
		RootCAs:            certPool,
	}, nil
}

func buildSelfSignedMtlsConfig(certificates *certification.Certificates) (*tls.Config, error) {
	logrus.Debug("authentication::buildSelfSignedMtlsConfig Building self-signed mTLS config")
	certPool, err := buildCaCertificatePool(certificates.GetCACertificate())
	if err != nil {
		return nil, err
	}

	certificate, err := buildCertificates(certificates.GetClientCertificate())
	if err != nil {
		return nil, err
	}
	logrus.Debugf("buildSelfSignedMtlsConfig Certificate: %v", certificate)

	return &tls.Config{
		InsecureSkipVerify: false,
		RootCAs:            certPool,
		ClientAuth:         tls.RequireAndVerifyClientCert,
		Certificates:       []tls.Certificate{certificate},
	}, nil
}

func buildBasicTlsConfig(skipVerify bool) *tls.Config {
	logrus.Debugf("authentication::buildBasicTlsConfig skipVerify(%v)", skipVerify)
	return &tls.Config{
		InsecureSkipVerify: skipVerify,
	}
}

func buildCaTlsConfig(certificates *certification.Certificates) (*tls.Config, error) {
	logrus.Debug("authentication::buildCaTlsConfig")
	certificate, err := buildCertificates(certificates.GetClientCertificate())
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		InsecureSkipVerify: false,
		Certificates:       []tls.Certificate{certificate},
	}, nil
}

func buildCertificates(privateKeyPEM []byte, certificatePEM []byte) (tls.Certificate, error) {
	logrus.Debug("authentication::buildCertificates")
	return tls.X509KeyPair(certificatePEM, privateKeyPEM)
}

func buildCaCertificatePool(caCert []byte) (*x509.CertPool, error) {
	logrus.Debug("authentication::buildCaCertificatePool")
	block, _ := pem.Decode(caCert)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing CA certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing certificate: %w", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AddCert(cert)

	return caCertPool, nil
}
