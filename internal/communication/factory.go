/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package communication

import (
	"fmt"
	netHttp "net/http"
	"time"

	"github.com/nginxinc/kubernetes-nginx-ingress/pkg/buildinfo"
)

// NewHTTPClient is a factory method to create a new Http Client configured for
// working with NGINXaaS or the N+ api. If skipVerify is set to true, the http
// transport will skip TLS certificate verification.
func NewHTTPClient(apiKey string, skipVerify bool) (*netHttp.Client, error) {
	headers := NewHeaders(apiKey)
	transport := NewTransport(skipVerify)
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

// NewTransport is a factory method to create a new basic Http Transport.
func NewTransport(skipVerify bool) *netHttp.Transport {
	transport := netHttp.DefaultTransport.(*netHttp.Transport).Clone()
	transport.TLSClientConfig.InsecureSkipVerify = skipVerify

	return transport
}
