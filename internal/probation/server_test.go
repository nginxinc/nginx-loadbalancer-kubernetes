/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package probation

import (
	"net/http"
	"testing"

	"github.com/nginxinc/kubernetes-nginx-ingress/test/mocks"
	"github.com/sirupsen/logrus"
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

func TestHealthServer_Start(t *testing.T) {
	server := NewHealthServer()
	server.Start()

	defer server.Stop()

	response, err := http.Get("http://localhost:51031/livez")
	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusAccepted, response.StatusCode)
	}

	logrus.Infof("received a response from the probe server: %v", response)
}
