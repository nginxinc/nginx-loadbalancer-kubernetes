/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package communication

import (
	"net/http"
	"strings"
)

// RoundTripper is a simple type that wraps the default net/communication RoundTripper to add additional headers.
type RoundTripper struct {
	Headers      []string
	RoundTripper http.RoundTripper
}

// NewRoundTripper is a factory method to create a new RoundTripper.
func NewRoundTripper(headers []string, transport *http.Transport) *RoundTripper {
	return &RoundTripper{
		Headers:      headers,
		RoundTripper: transport,
	}
}

// RoundTrip This simply adds our default headers to the request before passing it on to the default RoundTripper.
func (roundTripper *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	newRequest := new(http.Request)
	*newRequest = *req
	newRequest.Header = make(http.Header, len(req.Header))
	for k, s := range req.Header {
		newRequest.Header[k] = append([]string(nil), s...)
	}
	for _, s := range roundTripper.Headers {
		split := strings.SplitN(s, ":", 2)
		if len(split) >= 2 {
			newRequest.Header[split[0]] = append([]string(nil), split[1])
		}
	}
	return roundTripper.RoundTripper.RoundTrip(newRequest)
}
