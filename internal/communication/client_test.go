// Copyright 2023 f5 Inc. All rights reserved.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE file.

package communication

import (
	"testing"
)

func TestNewHttpClient(t *testing.T) {
	client, err := NewHttpClient()

	if err != nil {
		t.Fatalf(`Unexpected error: %v`, err)
	}

	if client == nil {
		t.Fatalf(`client should not be nil`)
	}
}

func TestNewHeaders(t *testing.T) {
	headers := NewHeaders()

	if headers == nil {
		t.Fatalf(`headers should not be nil`)
	}

	if len(headers) != 2 {
		t.Fatalf(`headers should have 2 elements`)
	}

	if headers[0] != "Content-Type: application/json" {
		t.Fatalf(`headers[0] should be "Content-Type: application/json"`)
	}

	if headers[1] != "Accept: application/json" {
		t.Fatalf(`headers[1] should be "Accept: application/json"`)
	}
}

func TestNewTlsConfig(t *testing.T) {
	config := NewTlsConfig()

	if config == nil {
		t.Fatalf(`config should not be nil`)
	}

	if !config.InsecureSkipVerify {
		t.Fatalf(`config.InsecureSkipVerify should be true`)
	}
}

func TestNewTransport(t *testing.T) {
	config := NewTlsConfig()
	transport := NewTransport(config)

	if transport == nil {
		t.Fatalf(`transport should not be nil`)
	}

	if transport.TLSClientConfig == nil {
		t.Fatalf(`transport.TLSClientConfig should not be nil`)
	}

	if transport.TLSClientConfig != config {
		t.Fatalf(`transport.TLSClientConfig should be the same as config`)
	}

	if !transport.TLSClientConfig.InsecureSkipVerify {
		t.Fatalf(`transport.TLSClientConfig.InsecureSkipVerify should be true`)
	}
}
