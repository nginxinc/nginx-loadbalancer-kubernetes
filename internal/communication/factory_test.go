/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package communication

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewHTTPClient(t *testing.T) {
	t.Parallel()

	client, err := NewHTTPClient("fakeKey", true)
	if err != nil {
		t.Fatalf(`Unexpected error: %v`, err)
	}

	if client == nil {
		t.Fatalf(`client should not be nil`)
	}
}

//nolint:goconst
func TestNewHeaders(t *testing.T) {
	t.Parallel()
	headers := NewHeaders("fakeKey")

	if headers == nil {
		t.Fatalf(`headers should not be nil`)
	}

	if len(headers) != 4 {
		t.Fatalf(`headers should have 3 elements`)
	}

	if headers[0] != "Content-Type: application/json" {
		t.Fatalf(`headers[0] should be "Content-Type: application/json"`)
	}

	if headers[1] != "Accept: application/json" {
		t.Fatalf(`headers[1] should be "Accept: application/json"`)
	}

	if headers[2] != "X-NLK-Version: " {
		t.Fatalf(`headers[2] should be "X-NLK-Version: "`)
	}

	if headers[3] != "Authorization: ApiKey fakeKey" {
		t.Fatalf(`headers[3] should be "Accept: Authorization: ApiKey fakeKey"`)
	}
}

func TestNewHeadersWithNoAPIKey(t *testing.T) {
	t.Parallel()
	headers := NewHeaders("")

	if headers == nil {
		t.Fatalf(`headers should not be nil`)
	}

	if len(headers) != 3 {
		t.Fatalf(`headers should have 2 elements`)
	}

	if headers[0] != "Content-Type: application/json" {
		t.Fatalf(`headers[0] should be "Content-Type: application/json"`)
	}

	if headers[1] != "Accept: application/json" {
		t.Fatalf(`headers[1] should be "Accept: application/json"`)
	}

	if headers[2] != "X-NLK-Version: " {
		t.Fatalf(`headers[2] should be "X-NLK-Version: "`)
	}
}

func TestNewTransport(t *testing.T) {
	t.Parallel()

	transport := NewTransport(false)

	if transport == nil {
		t.Fatalf(`transport should not be nil`)
	}

	if transport.TLSClientConfig == nil {
		t.Fatalf(`transport.TLSClientConfig should not be nil`)
	}

	require.False(t, transport.TLSClientConfig.InsecureSkipVerify)
}
