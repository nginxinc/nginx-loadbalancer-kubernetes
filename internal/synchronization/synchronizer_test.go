/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package synchronization

import (
	"context"
	"fmt"
	"testing"

	"github.com/nginxinc/kubernetes-nginx-ingress/internal/configuration"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/core"
	"github.com/nginxinc/kubernetes-nginx-ingress/test/mocks"
)

func TestSynchronizer_NewSynchronizer(t *testing.T) {
	t.Parallel()
	settings, err := configuration.NewSettings(context.Background(), nil)
	if err != nil {
		t.Fatalf(`Unexpected error: %v`, err)
	}

	rateLimiter := &mocks.MockRateLimiter{}

	synchronizer, err := NewSynchronizer(settings, rateLimiter)
	if err != nil {
		t.Fatalf(`should have been no error, %v`, err)
	}

	if synchronizer == nil {
		t.Fatal("should have an Synchronizer instance")
	}
}

func TestSynchronizer_AddEventNoHosts(t *testing.T) {
	t.Parallel()
	const expectedEventCount = 0
	event := &core.ServerUpdateEvent{
		ID:              "",
		NginxHost:       "",
		Type:            0,
		UpstreamName:    "",
		UpstreamServers: nil,
	}
	settings, err := configuration.NewSettings(context.Background(), nil)
	if err != nil {
		t.Fatalf(`Unexpected error: %v`, err)
	}
	rateLimiter := &mocks.MockRateLimiter{}

	synchronizer, err := NewSynchronizer(settings, rateLimiter)
	if err != nil {
		t.Fatalf(`should have been no error, %v`, err)
	}

	if synchronizer == nil {
		t.Fatal("should have an Synchronizer instance")
	}

	// NOTE: Ideally we have a custom logger that can be mocked to capture the log message
	// and assert a warning was logged that the NGINX Plus host was not specified.
	synchronizer.AddEvent(event)
	actualEventCount := rateLimiter.Len()
	if actualEventCount != expectedEventCount {
		t.Fatalf(`expected %v events, got %v`, expectedEventCount, actualEventCount)
	}
}

func TestSynchronizer_AddEventOneHost(t *testing.T) {
	t.Parallel()
	const expectedEventCount = 1
	events := buildEvents(1)
	settings, err := configuration.NewSettings(context.Background(), nil)
	if err != nil {
		t.Fatalf(`Unexpected error: %v`, err)
	}
	settings.NginxPlusHosts = []string{"https://localhost:8080"}
	rateLimiter := &mocks.MockRateLimiter{}

	synchronizer, err := NewSynchronizer(settings, rateLimiter)
	if err != nil {
		t.Fatalf(`should have been no error, %v`, err)
	}

	if synchronizer == nil {
		t.Fatal("should have an Synchronizer instance")
	}

	synchronizer.AddEvent(events[0])
	actualEventCount := rateLimiter.Len()
	if actualEventCount != expectedEventCount {
		t.Fatalf(`expected %v events, got %v`, expectedEventCount, actualEventCount)
	}
}

func TestSynchronizer_AddEventManyHosts(t *testing.T) {
	t.Parallel()
	const expectedEventCount = 1
	events := buildEvents(1)
	settings, err := configuration.NewSettings(context.Background(), nil)
	if err != nil {
		t.Fatalf(`Unexpected error: %v`, err)
	}
	settings.NginxPlusHosts = []string{
		"https://localhost:8080",
		"https://localhost:8081",
		"https://localhost:8082",
	}
	rateLimiter := &mocks.MockRateLimiter{}

	synchronizer, err := NewSynchronizer(settings, rateLimiter)
	if err != nil {
		t.Fatalf(`should have been no error, %v`, err)
	}

	if synchronizer == nil {
		t.Fatal("should have an Synchronizer instance")
	}

	synchronizer.AddEvent(events[0])
	actualEventCount := rateLimiter.Len()
	if actualEventCount != expectedEventCount {
		t.Fatalf(`expected %v events, got %v`, expectedEventCount, actualEventCount)
	}
}

func TestSynchronizer_AddEventsNoHosts(t *testing.T) {
	t.Parallel()
	const expectedEventCount = 0
	events := buildEvents(4)
	settings, err := configuration.NewSettings(context.Background(), nil)
	if err != nil {
		t.Fatalf(`Unexpected error: %v`, err)
	}
	rateLimiter := &mocks.MockRateLimiter{}

	synchronizer, err := NewSynchronizer(settings, rateLimiter)
	if err != nil {
		t.Fatalf(`should have been no error, %v`, err)
	}

	if synchronizer == nil {
		t.Fatal("should have an Synchronizer instance")
	}

	// NOTE: Ideally we have a custom logger that can be mocked to capture the log message
	// and assert a warning was logged that the NGINX Plus host was not specified.
	synchronizer.AddEvents(events)
	actualEventCount := rateLimiter.Len()
	if actualEventCount != expectedEventCount {
		t.Fatalf(`expected %v events, got %v`, expectedEventCount, actualEventCount)
	}
}

func TestSynchronizer_AddEventsOneHost(t *testing.T) {
	t.Parallel()
	const expectedEventCount = 4
	events := buildEvents(4)
	settings, err := configuration.NewSettings(context.Background(), nil)
	if err != nil {
		t.Fatalf(`Unexpected error: %v`, err)
	}
	settings.NginxPlusHosts = []string{"https://localhost:8080"}
	rateLimiter := &mocks.MockRateLimiter{}

	synchronizer, err := NewSynchronizer(settings, rateLimiter)
	if err != nil {
		t.Fatalf(`should have been no error, %v`, err)
	}

	if synchronizer == nil {
		t.Fatal("should have an Synchronizer instance")
	}

	synchronizer.AddEvents(events)
	actualEventCount := rateLimiter.Len()
	if actualEventCount != expectedEventCount {
		t.Fatalf(`expected %v events, got %v`, expectedEventCount, actualEventCount)
	}
}

func TestSynchronizer_AddEventsManyHosts(t *testing.T) {
	t.Parallel()
	const eventCount = 4
	events := buildEvents(eventCount)
	rateLimiter := &mocks.MockRateLimiter{}
	settings, err := configuration.NewSettings(context.Background(), nil)
	if err != nil {
		t.Fatalf(`Unexpected error: %v`, err)
	}
	settings.NginxPlusHosts = []string{
		"https://localhost:8080",
		"https://localhost:8081",
		"https://localhost:8082",
	}
	expectedEventCount := eventCount * len(settings.NginxPlusHosts)

	synchronizer, err := NewSynchronizer(settings, rateLimiter)
	if err != nil {
		t.Fatalf(`should have been no error, %v`, err)
	}

	if synchronizer == nil {
		t.Fatal("should have an Synchronizer instance")
	}

	synchronizer.AddEvents(events)
	actualEventCount := rateLimiter.Len()
	if actualEventCount != expectedEventCount {
		t.Fatalf(`expected %v events, got %v`, expectedEventCount, actualEventCount)
	}
}

func buildEvents(count int) core.ServerUpdateEvents {
	events := make(core.ServerUpdateEvents, count)
	for i := 0; i < count; i++ {
		events[i] = &core.ServerUpdateEvent{
			ID:              fmt.Sprintf("id-%v", i),
			NginxHost:       "https://localhost:8080",
			Type:            0,
			UpstreamName:    "",
			UpstreamServers: nil,
		}
	}
	return events
}
