// Copyright 2023 f5 Inc. All rights reserved.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE file.

package http

import (
	"net/http"
	netHttp "net/http"
	"strings"
)

type RoundTripper struct {
	Headers      []string
	RoundTripper http.RoundTripper
}

func NewRoundTripper(headers []string, transport *netHttp.Transport) *RoundTripper {
	return &RoundTripper{
		Headers:      headers,
		RoundTripper: transport,
	}
}

// RoundTrip Merge Headers
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
