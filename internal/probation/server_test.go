/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package probation

import (
	"github.com/nginxinc/kubernetes-nginx-ingress/test/mocks"
	"testing"
)

func TestHealthServer_HandleLive(t *testing.T) {
	server := NewHealthServer()
	writer := mocks.NewMockResponseWriter()
	server.HandleLive(writer, nil)

	if string(writer.Body()) != Ok {
		t.Errorf("HandleLive should return %s", Ok)
	}
}

func TestHealthServer_HandleReady(t *testing.T) {
	server := NewHealthServer()
	writer := mocks.NewMockResponseWriter()
	server.HandleReady(writer, nil)

	if string(writer.Body()) != Ok {
		t.Errorf("HandleReady should return %s", Ok)
	}
}

func TestHealthServer_HandleStartup(t *testing.T) {
	server := NewHealthServer()
	writer := mocks.NewMockResponseWriter()
	server.HandleStartup(writer, nil)

	if string(writer.Body()) != Ok {
		t.Errorf("HandleStartup should return %s", Ok)
	}
}

func TestHealthServer_HandleFailCheck(t *testing.T) {
	failCheck := mocks.NewMockCheck(false)
	server := NewHealthServer()
	writer := mocks.NewMockResponseWriter()
	server.handleProbe(writer, nil, failCheck)

	body := string(writer.Body())
	if body != "Service Not Available" {
		t.Errorf("Expected 'Service Not Available', got %v", body)
	}
}
