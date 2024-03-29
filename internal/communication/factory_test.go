/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package communication

import (
	"context"
	"testing"

	"github.com/nginxinc/kubernetes-nginx-ingress/internal/configuration"
	"k8s.io/client-go/kubernetes/fake"
)

func TestNewHTTPClient(t *testing.T) {
	t.Parallel()
	k8sClient := fake.NewSimpleClientset()
	settings, err := configuration.NewSettings(context.Background(), k8sClient)
	if err != nil {
		t.Fatalf(`Unexpected error: %v`, err)
	}
	client, err := NewHTTPClient(settings)

	if err != nil {
		t.Fatalf(`Unexpected error: %v`, err)
	}

	if client == nil {
		t.Fatalf(`client should not be nil`)
	}
}

func TestNewHeaders(t *testing.T) {
	t.Parallel()
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

func TestNewTransport(t *testing.T) {
	t.Parallel()
	k8sClient := fake.NewSimpleClientset()
	settings, _ := configuration.NewSettings(context.Background(), k8sClient)
	config := NewTLSConfig(settings)
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
