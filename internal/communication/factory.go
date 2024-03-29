/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package communication

import (
	"crypto/tls"
	netHttp "net/http"
	"time"

	"github.com/nginxinc/kubernetes-nginx-ingress/internal/authentication"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/configuration"
	"github.com/sirupsen/logrus"
)

// NewHttpClient is a factory method to create a new Http Client with a default configuration.
// RoundTripper is a wrapper around the default net/communication Transport to add additional headers, in this case,
// the Headers are configured for JSON.
func NewHttpClient(settings *configuration.Settings) (*netHttp.Client, error) {
	headers := NewHeaders()
	tlsConfig := NewTlsConfig(settings)
	transport := NewTransport(tlsConfig)
	roundTripper := NewRoundTripper(headers, transport)

	return &netHttp.Client{
		Transport:     roundTripper,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       time.Second * 10,
	}, nil
}

// NewHeaders is a factory method to create a new basic Http Headers slice.
func NewHeaders() []string {
	return []string{
		"Content-Type: application/json",
		"Accept: application/json",
	}
}

// NewTlsConfig is a factory method to create a new basic Tls Config.
// More attention should be given to the use of `InsecureSkipVerify: true`, as it is not recommended for production use.
func NewTlsConfig(settings *configuration.Settings) *tls.Config {
	tlsConfig, err := authentication.NewTlsConfig(settings)
	if err != nil {
		logrus.Warnf("Failed to create TLS config: %v", err)
		return &tls.Config{InsecureSkipVerify: true}
	}

	return tlsConfig
}

// NewTransport is a factory method to create a new basic Http Transport.
func NewTransport(config *tls.Config) *netHttp.Transport {
	transport := netHttp.DefaultTransport.(*netHttp.Transport)
	transport.TLSClientConfig = config

	return transport
}
