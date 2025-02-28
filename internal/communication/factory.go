/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package communication

import (
	"crypto/tls"
	"fmt"
	"log/slog"
	netHttp "net/http"
	"time"

	"github.com/nginxinc/kubernetes-nginx-ingress/internal/authentication"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/configuration"
	"github.com/nginxinc/kubernetes-nginx-ingress/pkg/buildinfo"
)

// NewHTTPClient is a factory method to create a new Http Client with a default configuration.
// RoundTripper is a wrapper around the default net/communication Transport to add additional headers, in this case,
// the Headers are configured for JSON.
func NewHTTPClient(settings configuration.Settings) (*netHttp.Client, error) {
	headers := NewHeaders(settings.APIKey)
	tlsConfig := NewTLSConfig(settings)
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
func NewHeaders(apiKey string) []string {
	headers := []string{
		"Content-Type: application/json",
		"Accept: application/json",
		fmt.Sprintf("X-NLK-Version: %s", buildinfo.SemVer()),
	}

	if apiKey != "" {
		headers = append(headers, fmt.Sprintf("Authorization: ApiKey %s", apiKey))
	}

	return headers
}

// NewTLSConfig is a factory method to create a new basic Tls Config.
// More attention should be given to the use of `InsecureSkipVerify: true`, as it is not recommended for production use.
func NewTLSConfig(settings configuration.Settings) *tls.Config {
	tlsConfig, err := authentication.NewTLSConfig(settings)
	if err != nil {
		slog.Warn("Failed to create TLS config", "error", err)
		return &tls.Config{InsecureSkipVerify: true} //nolint:gosec
	}

	return tlsConfig
}

// NewTransport is a factory method to create a new basic Http Transport.
func NewTransport(config *tls.Config) *netHttp.Transport {
	transport := netHttp.DefaultTransport.(*netHttp.Transport).Clone()
	transport.TLSClientConfig = config

	return transport
}
