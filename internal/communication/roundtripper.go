/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package communication

import (
	"net/http"
)

// Header represents a structured header key-value pair.
type Header struct {
	Key   string
	Value string
}

// RoundTripper wraps a RoundTripper to add custom headers to requests.
type RoundTripper struct {
	Headers      []Header
	RoundTripper http.RoundTripper
}

// NewRoundTripper creates a new RoundTripper with the given headers and transport.
func NewRoundTripper(headers []Header, baseTransport http.RoundTripper) *RoundTripper {
	if baseTransport == nil {
		baseTransport = http.DefaultTransport
	}
	return &RoundTripper{
		Headers:      headers,
		RoundTripper: baseTransport,
	}
}

// RoundTrip adds custom headers to the request before passing it to the base RoundTripper.
func (rt *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	newRequest := req.Clone(req.Context()) // Clone the request
	for _, header := range rt.Headers {
		if _, exists := newRequest.Header[header.Key]; !exists {
			newRequest.Header.Add(header.Key, header.Value)
		}
	}
	return rt.RoundTripper.RoundTrip(newRequest)
}
