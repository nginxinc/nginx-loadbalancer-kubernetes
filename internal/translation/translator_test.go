/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package translation

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/nginxinc/kubernetes-nginx-ingress/internal/core"
	v1 "k8s.io/api/core/v1"
)

const (
	AssertionFailureFormat = "expected %v events, got %v"
	ManyNodes              = 7
	NoNodes                = 0
	OneNode                = 1
	TranslateErrorFormat   = "Translate() error = %v"
)

/*
 * Created Event Tests
 */

func TestCreatedTranslateNoPorts(t *testing.T) {
	t.Parallel()
	const expectedEventCount = 0

	service := defaultService()
	event := buildCreatedEvent(service, OneNode)

	translatedEvents, err := Translate(&event)
	if err != nil {
		t.Fatalf(TranslateErrorFormat, err)
	}

	actualEventCount := len(translatedEvents)
	if actualEventCount != expectedEventCount {
		t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
	}
}

func TestCreatedTranslateNoInterestingPorts(t *testing.T) {
	t.Parallel()
	const expectedEventCount = 0
	const portCount = 1

	ports := generateUpdatablePorts(portCount, 0)
	service := serviceWithPorts(ports)
	event := buildCreatedEvent(service, OneNode)

	translatedEvents, err := Translate(&event)
	if err != nil {
		t.Fatalf(TranslateErrorFormat, err)
	}

	actualEventCount := len(translatedEvents)
	if actualEventCount != expectedEventCount {
		t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
	}
}

func TestCreatedTranslateOneInterestingPort(t *testing.T) {
	t.Parallel()
	const expectedEventCount = 1
	const portCount = 1

	ports := generatePorts(portCount)
	service := serviceWithPorts(ports)
	event := buildCreatedEvent(service, OneNode)

	translatedEvents, err := Translate(&event)
	if err != nil {
		t.Fatalf(TranslateErrorFormat, err)
	}

	actualEventCount := len(translatedEvents)
	if actualEventCount != expectedEventCount {
		t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
	}

	assertExpectedServerCount(t, OneNode, translatedEvents)
}

func TestCreatedTranslateManyInterestingPorts(t *testing.T) {
	t.Parallel()
	const expectedEventCount = 4
	const portCount = 4

	ports := generatePorts(portCount)
	service := serviceWithPorts(ports)
	event := buildCreatedEvent(service, OneNode)

	translatedEvents, err := Translate(&event)
	if err != nil {
		t.Fatalf(TranslateErrorFormat, err)
	}

	actualEventCount := len(translatedEvents)
	if actualEventCount != expectedEventCount {
		t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
	}

	assertExpectedServerCount(t, OneNode, translatedEvents)
}

func TestCreatedTranslateManyMixedPorts(t *testing.T) {
	t.Parallel()
	const expectedEventCount = 2
	const portCount = 6
	const updatablePortCount = 2

	ports := generateUpdatablePorts(portCount, updatablePortCount)
	service := serviceWithPorts(ports)
	event := buildCreatedEvent(service, OneNode)

	translatedEvents, err := Translate(&event)
	if err != nil {
		t.Fatalf(TranslateErrorFormat, err)
	}

	actualEventCount := len(translatedEvents)
	if actualEventCount != expectedEventCount {
		t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
	}

	assertExpectedServerCount(t, OneNode, translatedEvents)
}

func TestCreatedTranslateManyMixedPortsAndManyNodes(t *testing.T) {
	t.Parallel()
	const expectedEventCount = 2
	const portCount = 6
	const updatablePortCount = 2

	ports := generateUpdatablePorts(portCount, updatablePortCount)
	service := serviceWithPorts(ports)
	event := buildCreatedEvent(service, ManyNodes)

	translatedEvents, err := Translate(&event)
	if err != nil {
		t.Fatalf(TranslateErrorFormat, err)
	}

	actualEventCount := len(translatedEvents)
	if actualEventCount != expectedEventCount {
		t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
	}

	assertExpectedServerCount(t, ManyNodes, translatedEvents)
}

/*
 * Updated Event Tests
 */

func TestUpdatedTranslateNoPorts(t *testing.T) {
	t.Parallel()
	const expectedEventCount = 0

	service := defaultService()
	event := buildUpdatedEvent(service, OneNode)

	translatedEvents, err := Translate(&event)
	if err != nil {
		t.Fatalf(TranslateErrorFormat, err)
	}

	actualEventCount := len(translatedEvents)
	if actualEventCount != expectedEventCount {
		t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
	}
}

func TestUpdatedTranslateNoInterestingPorts(t *testing.T) {
	t.Parallel()
	const expectedEventCount = 0
	const portCount = 1

	ports := generateUpdatablePorts(portCount, 0)
	service := serviceWithPorts(ports)
	event := buildUpdatedEvent(service, OneNode)

	translatedEvents, err := Translate(&event)
	if err != nil {
		t.Fatalf(TranslateErrorFormat, err)
	}

	actualEventCount := len(translatedEvents)
	if actualEventCount != expectedEventCount {
		t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
	}
}

func TestUpdatedTranslateOneInterestingPort(t *testing.T) {
	t.Parallel()
	const expectedEventCount = 1
	const portCount = 1

	ports := generatePorts(portCount)
	service := serviceWithPorts(ports)
	event := buildUpdatedEvent(service, OneNode)

	translatedEvents, err := Translate(&event)
	if err != nil {
		t.Fatalf(TranslateErrorFormat, err)
	}

	actualEventCount := len(translatedEvents)
	if actualEventCount != expectedEventCount {
		t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
	}

	assertExpectedServerCount(t, OneNode, translatedEvents)
}

func TestUpdatedTranslateManyInterestingPorts(t *testing.T) {
	t.Parallel()
	const expectedEventCount = 4
	const portCount = 4

	ports := generatePorts(portCount)
	service := serviceWithPorts(ports)
	event := buildUpdatedEvent(service, OneNode)

	translatedEvents, err := Translate(&event)
	if err != nil {
		t.Fatalf(TranslateErrorFormat, err)
	}

	actualEventCount := len(translatedEvents)
	if actualEventCount != expectedEventCount {
		t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
	}

	assertExpectedServerCount(t, OneNode, translatedEvents)
}

func TestUpdatedTranslateManyMixedPorts(t *testing.T) {
	t.Parallel()
	const expectedEventCount = 2
	const portCount = 6
	const updatablePortCount = 2

	ports := generateUpdatablePorts(portCount, updatablePortCount)
	service := serviceWithPorts(ports)
	event := buildUpdatedEvent(service, OneNode)

	translatedEvents, err := Translate(&event)
	if err != nil {
		t.Fatalf(TranslateErrorFormat, err)
	}

	actualEventCount := len(translatedEvents)
	if actualEventCount != expectedEventCount {
		t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
	}

	assertExpectedServerCount(t, OneNode, translatedEvents)
}

func TestUpdatedTranslateManyMixedPortsAndManyNodes(t *testing.T) {
	t.Parallel()
	const expectedEventCount = 2
	const portCount = 6
	const updatablePortCount = 2

	ports := generateUpdatablePorts(portCount, updatablePortCount)
	service := serviceWithPorts(ports)
	event := buildUpdatedEvent(service, ManyNodes)

	translatedEvents, err := Translate(&event)
	if err != nil {
		t.Fatalf(TranslateErrorFormat, err)
	}

	actualEventCount := len(translatedEvents)
	if actualEventCount != expectedEventCount {
		t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
	}

	assertExpectedServerCount(t, ManyNodes, translatedEvents)
}

/*
 * Deleted Event Tests
 */

func TestDeletedTranslateNoPortsAndNoNodes(t *testing.T) {
	t.Parallel()
	const expectedEventCount = 0

	service := defaultService()
	event := buildDeletedEvent(service, NoNodes)

	translatedEvents, err := Translate(&event)
	if err != nil {
		t.Fatalf(TranslateErrorFormat, err)
	}

	actualEventCount := len(translatedEvents)
	if actualEventCount != expectedEventCount {
		t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
	}

	assertExpectedServerCount(t, ManyNodes, translatedEvents)
}

func TestDeletedTranslateNoInterestingPortsAndNoNodes(t *testing.T) {
	t.Parallel()
	const expectedEventCount = 0
	const portCount = 1

	ports := generateUpdatablePorts(portCount, 0)
	service := serviceWithPorts(ports)
	event := buildDeletedEvent(service, NoNodes)

	translatedEvents, err := Translate(&event)
	if err != nil {
		t.Fatalf(TranslateErrorFormat, err)
	}

	actualEventCount := len(translatedEvents)
	if actualEventCount != expectedEventCount {
		t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
	}

	assertExpectedServerCount(t, ManyNodes, translatedEvents)
}

func TestDeletedTranslateOneInterestingPortAndNoNodes(t *testing.T) {
	t.Parallel()

	const expectedEventCount = 0
	const portCount = 1

	ports := generatePorts(portCount)
	service := serviceWithPorts(ports)
	event := buildDeletedEvent(service, NoNodes)

	translatedEvents, err := Translate(&event)
	if err != nil {
		t.Fatalf(TranslateErrorFormat, err)
	}

	actualEventCount := len(translatedEvents)
	if actualEventCount != expectedEventCount {
		t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
	}

	assertExpectedServerCount(t, ManyNodes, translatedEvents)
}

func TestDeletedTranslateManyInterestingPortsAndNoNodes(t *testing.T) {
	t.Parallel()
	const expectedEventCount = 0
	const portCount = 4

	ports := generatePorts(portCount)
	service := serviceWithPorts(ports)
	event := buildDeletedEvent(service, NoNodes)

	translatedEvents, err := Translate(&event)
	if err != nil {
		t.Fatalf(TranslateErrorFormat, err)
	}

	actualEventCount := len(translatedEvents)
	if actualEventCount != expectedEventCount {
		t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
	}

	assertExpectedServerCount(t, ManyNodes, translatedEvents)
}

func TestDeletedTranslateManyMixedPortsAndNoNodes(t *testing.T) {
	t.Parallel()
	const expectedEventCount = 0
	const portCount = 6
	const updatablePortCount = 2

	ports := generateUpdatablePorts(portCount, updatablePortCount)
	service := serviceWithPorts(ports)
	event := buildDeletedEvent(service, NoNodes)

	translatedEvents, err := Translate(&event)
	if err != nil {
		t.Fatalf(TranslateErrorFormat, err)
	}

	actualEventCount := len(translatedEvents)
	if actualEventCount != expectedEventCount {
		t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
	}

	assertExpectedServerCount(t, ManyNodes, translatedEvents)
}

func TestDeletedTranslateNoPortsAndOneNode(t *testing.T) {
	t.Parallel()
	const expectedEventCount = 0

	service := defaultService()
	event := buildDeletedEvent(service, OneNode)

	translatedEvents, err := Translate(&event)
	if err != nil {
		t.Fatalf(TranslateErrorFormat, err)
	}

	actualEventCount := len(translatedEvents)
	if actualEventCount != expectedEventCount {
		t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
	}

	assertExpectedServerCount(t, ManyNodes, translatedEvents)
}

func TestDeletedTranslateNoInterestingPortsAndOneNode(t *testing.T) {
	t.Parallel()
	const expectedEventCount = 0
	const portCount = 1

	ports := generateUpdatablePorts(portCount, 0)
	service := serviceWithPorts(ports)
	event := buildDeletedEvent(service, OneNode)

	translatedEvents, err := Translate(&event)
	if err != nil {
		t.Fatalf(TranslateErrorFormat, err)
	}

	actualEventCount := len(translatedEvents)
	if actualEventCount != expectedEventCount {
		t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
	}

	assertExpectedServerCount(t, ManyNodes, translatedEvents)
}

func TestDeletedTranslateOneInterestingPortAndOneNode(t *testing.T) {
	t.Parallel()
	const expectedEventCount = 1
	const portCount = 1

	ports := generatePorts(portCount)
	service := serviceWithPorts(ports)
	event := buildDeletedEvent(service, OneNode)

	translatedEvents, err := Translate(&event)
	if err != nil {
		t.Fatalf(TranslateErrorFormat, err)
	}

	actualEventCount := len(translatedEvents)
	if actualEventCount != expectedEventCount {
		t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
	}

	assertExpectedServerCount(t, OneNode, translatedEvents)
}

func TestDeletedTranslateManyInterestingPortsAndOneNode(t *testing.T) {
	t.Parallel()
	const expectedEventCount = 4
	const portCount = 4

	ports := generatePorts(portCount)
	service := serviceWithPorts(ports)
	event := buildDeletedEvent(service, OneNode)

	translatedEvents, err := Translate(&event)
	if err != nil {
		t.Fatalf(TranslateErrorFormat, err)
	}

	actualEventCount := len(translatedEvents)
	if actualEventCount != expectedEventCount {
		t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
	}

	assertExpectedServerCount(t, OneNode, translatedEvents)
}

func TestDeletedTranslateManyMixedPortsAndOneNode(t *testing.T) {
	t.Parallel()
	const expectedEventCount = 2
	const portCount = 6
	const updatablePortCount = 2

	ports := generateUpdatablePorts(portCount, updatablePortCount)
	service := serviceWithPorts(ports)
	event := buildDeletedEvent(service, OneNode)

	translatedEvents, err := Translate(&event)
	if err != nil {
		t.Fatalf(TranslateErrorFormat, err)
	}

	actualEventCount := len(translatedEvents)
	if actualEventCount != expectedEventCount {
		t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
	}

	assertExpectedServerCount(t, OneNode, translatedEvents)
}

func TestDeletedTranslateNoPortsAndManyNodes(t *testing.T) {
	t.Parallel()
	const expectedEventCount = 0

	service := defaultService()
	event := buildDeletedEvent(service, ManyNodes)

	translatedEvents, err := Translate(&event)
	if err != nil {
		t.Fatalf(TranslateErrorFormat, err)
	}

	actualEventCount := len(translatedEvents)
	if actualEventCount != expectedEventCount {
		t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
	}

	assertExpectedServerCount(t, ManyNodes, translatedEvents)
}

func TestDeletedTranslateNoInterestingPortsAndManyNodes(t *testing.T) {
	t.Parallel()
	const portCount = 1
	const updatablePortCount = 0
	const expectedEventCount = updatablePortCount * ManyNodes

	ports := generateUpdatablePorts(portCount, updatablePortCount)
	service := serviceWithPorts(ports)
	event := buildDeletedEvent(service, ManyNodes)

	translatedEvents, err := Translate(&event)
	if err != nil {
		t.Fatalf(TranslateErrorFormat, err)
	}

	actualEventCount := len(translatedEvents)
	if actualEventCount != expectedEventCount {
		t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
	}

	assertExpectedServerCount(t, ManyNodes, translatedEvents)
}

func TestDeletedTranslateOneInterestingPortAndManyNodes(t *testing.T) {
	t.Parallel()
	const portCount = 1
	const expectedEventCount = portCount * ManyNodes

	ports := generatePorts(portCount)
	service := serviceWithPorts(ports)
	event := buildDeletedEvent(service, ManyNodes)

	translatedEvents, err := Translate(&event)
	if err != nil {
		t.Fatalf(TranslateErrorFormat, err)
	}

	actualEventCount := len(translatedEvents)
	if actualEventCount != expectedEventCount {
		t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
	}

	assertExpectedServerCount(t, OneNode, translatedEvents)
}

func TestDeletedTranslateManyInterestingPortsAndManyNodes(t *testing.T) {
	t.Parallel()
	const portCount = 4
	const expectedEventCount = portCount * ManyNodes

	ports := generatePorts(portCount)
	service := serviceWithPorts(ports)
	event := buildDeletedEvent(service, ManyNodes)

	translatedEvents, err := Translate(&event)
	if err != nil {
		t.Fatalf(TranslateErrorFormat, err)
	}

	actualEventCount := len(translatedEvents)
	if actualEventCount != expectedEventCount {
		t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
	}

	assertExpectedServerCount(t, OneNode, translatedEvents)
}

func TestDeletedTranslateManyMixedPortsAndManyNodes(t *testing.T) {
	t.Parallel()
	const portCount = 6
	const updatablePortCount = 2
	const expectedEventCount = updatablePortCount * ManyNodes

	ports := generateUpdatablePorts(portCount, updatablePortCount)
	service := serviceWithPorts(ports)
	event := buildDeletedEvent(service, ManyNodes)

	translatedEvents, err := Translate(&event)
	if err != nil {
		t.Fatalf(TranslateErrorFormat, err)
	}

	actualEventCount := len(translatedEvents)
	if actualEventCount != expectedEventCount {
		t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
	}

	assertExpectedServerCount(t, OneNode, translatedEvents)
}

func assertExpectedServerCount(t *testing.T, expectedCount int, events core.ServerUpdateEvents) {
	for _, translatedEvent := range events {
		serverCount := len(translatedEvent.UpstreamServers)
		if serverCount != expectedCount {
			t.Fatalf("expected %d servers, got %d", expectedCount, serverCount)
		}
	}
}

func defaultService() *v1.Service {
	return &v1.Service{}
}

func serviceWithPorts(ports []v1.ServicePort) *v1.Service {
	return &v1.Service{
		Spec: v1.ServiceSpec{
			Ports: ports,
		},
	}
}

func buildCreatedEvent(service *v1.Service, nodeCount int) core.Event {
	return buildEvent(core.Created, service, nodeCount)
}

func buildDeletedEvent(service *v1.Service, nodeCount int) core.Event {
	return buildEvent(core.Deleted, service, nodeCount)
}

func buildUpdatedEvent(service *v1.Service, nodeCount int) core.Event {
	return buildEvent(core.Updated, service, nodeCount)
}

func buildEvent(eventType core.EventType, service *v1.Service, nodeCount int) core.Event {
	previousService := defaultService()

	nodeIps := generateNodeIps(nodeCount)

	return core.NewEvent(eventType, service, previousService, nodeIps)
}

func generateNodeIps(count int) []string {
	var nodeIps []string

	for i := 0; i < count; i++ {
		nodeIps = append(nodeIps, fmt.Sprintf("10.0.0.%v", i))
	}

	return nodeIps
}

func generatePorts(portCount int) []v1.ServicePort {
	return generateUpdatablePorts(portCount, portCount)
}

// This is probably A Little Bit of Too Much™, but helps to ensure ordering is not a factor.
func generateUpdatablePorts(portCount int, updatableCount int) []v1.ServicePort {
	ports := []v1.ServicePort{}

	updatable := make([]string, updatableCount)
	nonupdatable := make([]string, portCount-updatableCount)
	contexts := []string{"http-", "stream-"}

	for i := range updatable {
		randomIndex := int(rand.Float32() * 2.0)
		updatable[i] = contexts[randomIndex]
	}

	for j := range nonupdatable {
		nonupdatable[j] = "olm-"
	}

	var prefixes []string
	prefixes = append(prefixes, updatable...)
	prefixes = append(prefixes, nonupdatable...)

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	random.Shuffle(len(prefixes), func(i, j int) { prefixes[i], prefixes[j] = prefixes[j], prefixes[i] })

	for i, prefix := range prefixes {
		ports = append(ports, v1.ServicePort{
			Name: fmt.Sprintf("%supstream%d", prefix, i),
		})
	}

	return ports
}
