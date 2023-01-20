// Copyright 2023 f5 Inc. All rights reserved.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE file.

package communication

import (
	"crypto/tls"
	netHttp "net/http"
	"time"
)

func NewHttpClient() (*netHttp.Client, error) {
	headers := NewHeaders()
	tlsConfig := NewTlsConfig()
	transport := NewTransport(tlsConfig)
	roundTripper := NewRoundTripper(headers, transport)

	return &netHttp.Client{
		Transport:     roundTripper,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       time.Second * 10,
	}, nil
}

func NewHeaders() []string {
	return []string{
		"Content-Type: application/json",
		"Accept: application/json",
	}
}

func NewTlsConfig() *tls.Config {
	return &tls.Config{InsecureSkipVerify: true}
}

func NewTransport(config *tls.Config) *netHttp.Transport {
	transport := netHttp.DefaultTransport.(*netHttp.Transport)
	transport.TLSClientConfig = config

	return transport
}
