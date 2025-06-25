/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package communication

import (
	"bytes"
	"context"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/configuration"
	"k8s.io/client-go/kubernetes/fake"
	netHttp "net/http"
	"testing"
)

func TestNewRoundTripper(t *testing.T) {
	k8sClient := fake.NewSimpleClientset()
	settings, _ := configuration.NewSettings(context.Background(), k8sClient)
	headers := NewHeaders()
	transport := NewTransport(NewTlsConfig(settings))
	roundTripper := NewRoundTripper(headers, transport)

	if roundTripper == nil {
		t.Fatalf(`roundTripper should not be nil`)
	}

	if roundTripper.Headers == nil {
		t.Fatalf(`roundTripper.Headers should not be nil`)
	}

	if len(roundTripper.Headers) != 2 {
		t.Fatalf(`roundTripper.Headers should have 2 elements`)
	}

	if roundTripper.Headers[0] != "Content-Type: application/json" {
		t.Fatalf(`roundTripper.Headers[0] should be "Content-Type: application/json"`)
	}

	if roundTripper.Headers[1] != "Accept: application/json" {
		t.Fatalf(`roundTripper.Headers[1] should be "Accept: application/json"`)
	}

	if roundTripper.RoundTripper == nil {
		t.Fatalf(`roundTripper.RoundTripper should not be nil`)
	}
}

func TestRoundTripperRoundTrip(t *testing.T) {
	k8sClient := fake.NewSimpleClientset()
	settings, err := configuration.NewSettings(context.Background(), k8sClient)
	headers := NewHeaders()
	transport := NewTransport(NewTlsConfig(settings))
	roundTripper := NewRoundTripper(headers, transport)

	request, err := NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf(`Unexpected error: %v`, err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-nginx-loadbalancer-kubernetes", "nlk")

	response, err := roundTripper.RoundTrip(request)
	if err != nil {
		t.Fatalf(`Unexpected error: %v`, err)
	}

	if response == nil {
		t.Fatalf(`response should not be nil`)
	}

	headerLen := len(response.Header)
	if headerLen <= 2 {
		t.Fatalf(`response.Header should have at least 2 elements, found %d`, headerLen)
	}
}

func NewRequest(method string, url string, body []byte) (*netHttp.Request, error) {
	request, err := netHttp.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	return request, nil
}
