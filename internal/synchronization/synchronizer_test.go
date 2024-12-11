/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package synchronization

import (
	"fmt"
	"testing"
	"time"

	"github.com/nginxinc/kubernetes-nginx-ingress/internal/configuration"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/core"
	"github.com/nginxinc/kubernetes-nginx-ingress/test/mocks"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	corelisters "k8s.io/client-go/listers/core/v1"
)

func TestSynchronizer_NewSynchronizer(t *testing.T) {
	t.Parallel()

	rateLimiter := &mocks.MockRateLimiter{}

	synchronizer, err := NewSynchronizer(
		configuration.Settings{},
		rateLimiter,
		&fakeTranslator{},
		newFakeServicesLister(defaultService()),
	)
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

	rateLimiter := &mocks.MockRateLimiter{}

	synchronizer, err := NewSynchronizer(
		defaultSettings(),
		rateLimiter,
		&fakeTranslator{},
		newFakeServicesLister(defaultService()),
	)
	if err != nil {
		t.Fatalf(`should have been no error, %v`, err)
	}

	if synchronizer == nil {
		t.Fatal("should have an Synchronizer instance")
	}

	// NOTE: Ideally we have a custom logger that can be mocked to capture the log message
	// and assert a warning was logged that the NGINX Plus host was not specified.
	synchronizer.AddEvent(core.Event{})
	actualEventCount := rateLimiter.Len()
	if actualEventCount != expectedEventCount {
		t.Fatalf(`expected %v events, got %v`, expectedEventCount, actualEventCount)
	}
}

func TestSynchronizer_AddEventOneHost(t *testing.T) {
	t.Parallel()
	const expectedEventCount = 1
	events := buildServerUpdateEvents(1)

	rateLimiter := &mocks.MockRateLimiter{}

	synchronizer, err := NewSynchronizer(
		defaultSettings("https://localhost:8080"),
		rateLimiter,
		&fakeTranslator{events, nil},
		newFakeServicesLister(defaultService()),
	)
	if err != nil {
		t.Fatalf(`should have been no error, %v`, err)
	}

	if synchronizer == nil {
		t.Fatal("should have an Synchronizer instance")
	}

	synchronizer.AddEvent(buildServiceUpdateEvent(1))
	actualEventCount := rateLimiter.Len()
	if actualEventCount != expectedEventCount {
		t.Fatalf(`expected %v events, got %v`, expectedEventCount, actualEventCount)
	}
}

func TestSynchronizer_AddEventManyHosts(t *testing.T) {
	t.Parallel()
	const expectedEventCount = 1
	events := buildServerUpdateEvents(1)
	hosts := []string{
		"https://localhost:8080",
		"https://localhost:8081",
		"https://localhost:8082",
	}

	rateLimiter := &mocks.MockRateLimiter{}

	synchronizer, err := NewSynchronizer(
		defaultSettings(hosts...),
		rateLimiter,
		&fakeTranslator{events, nil},
		newFakeServicesLister(defaultService()),
	)
	if err != nil {
		t.Fatalf(`should have been no error, %v`, err)
	}

	if synchronizer == nil {
		t.Fatal("should have an Synchronizer instance")
	}

	synchronizer.AddEvent(buildServiceUpdateEvent(1))
	actualEventCount := rateLimiter.Len()
	if actualEventCount != expectedEventCount {
		t.Fatalf(`expected %v events, got %v`, expectedEventCount, actualEventCount)
	}
}

func TestSynchronizer_AddEventsNoHosts(t *testing.T) {
	t.Parallel()
	const expectedEventCount = 0
	events := buildServerUpdateEvents(4)
	rateLimiter := &mocks.MockRateLimiter{}

	synchronizer, err := NewSynchronizer(
		defaultSettings(),
		rateLimiter,
		&fakeTranslator{events, nil},
		newFakeServicesLister(defaultService()),
	)
	if err != nil {
		t.Fatalf(`should have been no error, %v`, err)
	}

	if synchronizer == nil {
		t.Fatal("should have an Synchronizer instance")
	}

	// NOTE: Ideally we have a custom logger that can be mocked to capture the log message
	// and assert a warning was logged that the NGINX Plus host was not specified.
	for i := 0; i < 4; i++ {
		synchronizer.AddEvent(buildServiceUpdateEvent(i))
	}

	actualEventCount := rateLimiter.Len()
	if actualEventCount != expectedEventCount {
		t.Fatalf(`expected %v events, got %v`, expectedEventCount, actualEventCount)
	}
}

func TestSynchronizer_AddEventsOneHost(t *testing.T) {
	t.Parallel()
	const expectedEventCount = 4
	events := buildServerUpdateEvents(1)
	rateLimiter := &mocks.MockRateLimiter{}

	synchronizer, err := NewSynchronizer(
		defaultSettings("https://localhost:8080"),
		rateLimiter,
		&fakeTranslator{events, nil},
		newFakeServicesLister(defaultService()),
	)
	if err != nil {
		t.Fatalf(`should have been no error, %v`, err)
	}

	if synchronizer == nil {
		t.Fatal("should have an Synchronizer instance")
	}

	for i := 0; i < 4; i++ {
		synchronizer.AddEvent(buildServiceUpdateEvent(i))
	}

	actualEventCount := rateLimiter.Len()
	if actualEventCount != expectedEventCount {
		t.Fatalf(`expected %v events, got %v`, expectedEventCount, actualEventCount)
	}
}

func TestSynchronizer_AddEventsManyHosts(t *testing.T) {
	t.Parallel()
	const eventCount = 4
	events := buildServerUpdateEvents(eventCount)
	rateLimiter := &mocks.MockRateLimiter{}

	hosts := []string{
		"https://localhost:8080",
		"https://localhost:8081",
		"https://localhost:8082",
	}

	expectedEventCount := 4

	synchronizer, err := NewSynchronizer(
		defaultSettings(hosts...),
		rateLimiter,
		&fakeTranslator{events, nil},
		newFakeServicesLister(defaultService()),
	)
	if err != nil {
		t.Fatalf(`should have been no error, %v`, err)
	}

	if synchronizer == nil {
		t.Fatal("should have an Synchronizer instance")
	}

	for i := 0; i < eventCount; i++ {
		synchronizer.AddEvent(buildServiceUpdateEvent(i))
	}

	actualEventCount := rateLimiter.Len()
	if actualEventCount != expectedEventCount {
		t.Fatalf(`expected %v events, got %v`, expectedEventCount, actualEventCount)
	}
}

func buildServiceUpdateEvent(serviceID int) core.Event {
	return core.Event{
		Service: &v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("test-service%d", serviceID),
				Namespace: "test-namespace",
			},
		},
	}
}

func buildServerUpdateEvents(count int) core.ServerUpdateEvents {
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

func defaultSettings(nginxHosts ...string) configuration.Settings {
	return configuration.Settings{
		NginxPlusHosts: nginxHosts,
		Synchronizer: configuration.SynchronizerSettings{
			MaxMillisecondsJitter: 750,
			MinMillisecondsJitter: 250,
			RetryCount:            5,
			Threads:               1,
			WorkQueueSettings: configuration.WorkQueueSettings{
				RateLimiterBase: time.Second * 2,
				RateLimiterMax:  time.Second * 60,
				Name:            "nlk-synchronizer",
			},
		},
	}
}

type fakeTranslator struct {
	events core.ServerUpdateEvents
	err    error
}

func (t *fakeTranslator) Translate(event *core.Event) (core.ServerUpdateEvents, error) {
	return t.events, t.err
}

func newFakeServicesLister(list ...*v1.Service) corelisters.ServiceLister {
	return &servicesLister{
		list: list,
	}
}

type servicesLister struct {
	list []*v1.Service
	err  error
}

func (l *servicesLister) List(selector labels.Selector) (ret []*v1.Service, err error) {
	return l.list, l.err
}

func (l *servicesLister) Get(name string) (*v1.Service, error) {
	for _, service := range l.list {
		if service.Name == name {
			return service, nil
		}
	}

	return nil, nil
}

func (l *servicesLister) Services(name string) corelisters.ServiceNamespaceLister {
	return l
}

func defaultService() *v1.Service {
	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "default-service",
			Labels: map[string]string{"kubernetes.io/service-name": "default-service"},
		},
		Spec: v1.ServiceSpec{
			Type: v1.ServiceTypeNodePort,
		},
	}
}
