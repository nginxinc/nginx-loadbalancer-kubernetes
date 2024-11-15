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
	"github.com/nginxinc/kubernetes-nginx-ingress/pkg/pointer"
	v1 "k8s.io/api/core/v1"
	discovery "k8s.io/api/discovery/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	corelisters "k8s.io/client-go/listers/core/v1"
	discoverylisters "k8s.io/client-go/listers/discovery/v1"
)

const (
	AssertionFailureFormat = "expected %v events, got %v"
	ManyNodes              = 7
	NoNodes                = 0
	OneNode                = 1
	ManyEndpointSlices     = 7
	NoEndpointSlices       = 0
	OneEndpointSlice       = 1
	TranslateErrorFormat   = "Translate() error = %v"
)

/*
 * Created Event Tests
 */

func TestCreatedTranslateNoPorts(t *testing.T) {
	t.Parallel()
	testcases := map[string]struct{ serviceType v1.ServiceType }{
		"nodePort":  {v1.ServiceTypeNodePort},
		"clusterIP": {v1.ServiceTypeClusterIP},
	}

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			const expectedEventCount = 0

			service := defaultService(tc.serviceType)
			event := buildCreatedEvent(service)

			translator := NewTranslator(
				NewFakeEndpointSliceLister([]*discovery.EndpointSlice{}, nil),
				NewFakeNodeLister([]*v1.Node{}, nil),
			)

			translatedEvents, err := translator.Translate(&event)
			if err != nil {
				t.Fatalf(TranslateErrorFormat, err)
			}

			actualEventCount := len(translatedEvents)
			if actualEventCount != expectedEventCount {
				t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
			}
		})
	}
}

func TestCreatedTranslateNoInterestingPorts(t *testing.T) {
	t.Parallel()
	testcases := map[string]struct{ serviceType v1.ServiceType }{
		"nodePort":  {v1.ServiceTypeNodePort},
		"clusterIP": {v1.ServiceTypeClusterIP},
	}

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			const expectedEventCount = 0
			const portCount = 1

			ports := generateUpdatablePorts(portCount, 0)
			service := serviceWithPorts(tc.serviceType, ports)
			event := buildCreatedEvent(service)

			translator := NewTranslator(
				NewFakeEndpointSliceLister([]*discovery.EndpointSlice{}, nil),
				NewFakeNodeLister([]*v1.Node{}, nil),
			)

			translatedEvents, err := translator.Translate(&event)
			if err != nil {
				t.Fatalf(TranslateErrorFormat, err)
			}

			actualEventCount := len(translatedEvents)
			if actualEventCount != expectedEventCount {
				t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
			}
		})
	}
}

//nolint:dupl
func TestCreatedTranslateOneInterestingPort(t *testing.T) {
	t.Parallel()
	testcases := map[string]struct {
		serviceType         v1.ServiceType
		nodes               []*v1.Node
		endpoints           []*discovery.EndpointSlice
		expectedServerCount int
	}{
		"nodePort": {
			serviceType:         v1.ServiceTypeNodePort,
			nodes:               generateNodes(OneNode),
			expectedServerCount: OneNode,
		},
		"clusterIP": {
			serviceType:         v1.ServiceTypeClusterIP,
			endpoints:           generateEndpointSlices(OneEndpointSlice, 1, 1),
			expectedServerCount: OneEndpointSlice,
		},
	}

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			const expectedEventCount = 1
			const portCount = 1

			ports := generatePorts(portCount)
			service := serviceWithPorts(tc.serviceType, ports)
			event := buildCreatedEvent(service)

			translator := NewTranslator(NewFakeEndpointSliceLister(tc.endpoints, nil), NewFakeNodeLister(tc.nodes, nil))
			translatedEvents, err := translator.Translate(&event)
			if err != nil {
				t.Fatalf(TranslateErrorFormat, err)
			}

			actualEventCount := len(translatedEvents)
			if actualEventCount != expectedEventCount {
				t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
			}

			assertExpectedServerCount(t, tc.expectedServerCount, translatedEvents)
		})
	}
}

//nolint:dupl
func TestCreatedTranslateManyInterestingPorts(t *testing.T) {
	t.Parallel()
	testcases := map[string]struct {
		serviceType         v1.ServiceType
		nodes               []*v1.Node
		endpoints           []*discovery.EndpointSlice
		expectedServerCount int
	}{
		"nodePort": {
			serviceType:         v1.ServiceTypeNodePort,
			nodes:               generateNodes(OneNode),
			expectedServerCount: OneNode,
		},
		"clusterIP": {
			serviceType:         v1.ServiceTypeClusterIP,
			endpoints:           generateEndpointSlices(OneEndpointSlice, 4, 4),
			expectedServerCount: OneEndpointSlice,
		},
	}

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			const expectedEventCount = 4
			const portCount = 4

			ports := generatePorts(portCount)
			service := serviceWithPorts(tc.serviceType, ports)
			event := buildCreatedEvent(service)

			translator := NewTranslator(NewFakeEndpointSliceLister(tc.endpoints, nil), NewFakeNodeLister(tc.nodes, nil))
			translatedEvents, err := translator.Translate(&event)
			if err != nil {
				t.Fatalf(TranslateErrorFormat, err)
			}

			actualEventCount := len(translatedEvents)
			if actualEventCount != expectedEventCount {
				t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
			}

			assertExpectedServerCount(t, tc.expectedServerCount, translatedEvents)
		})
	}
}

//nolint:dupl
func TestCreatedTranslateManyMixedPorts(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		serviceType         v1.ServiceType
		nodes               []*v1.Node
		endpoints           []*discovery.EndpointSlice
		expectedServerCount int
	}{
		"nodePort": {
			serviceType:         v1.ServiceTypeNodePort,
			nodes:               generateNodes(OneNode),
			expectedServerCount: OneNode,
		},
		"clusterIP": {
			serviceType:         v1.ServiceTypeClusterIP,
			endpoints:           generateEndpointSlices(OneEndpointSlice, 6, 2),
			expectedServerCount: OneEndpointSlice,
		},
	}

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			const expectedEventCount = 2
			const portCount = 6
			const updatablePortCount = 2

			ports := generateUpdatablePorts(portCount, updatablePortCount)
			service := serviceWithPorts(tc.serviceType, ports)
			event := buildCreatedEvent(service)

			translator := NewTranslator(NewFakeEndpointSliceLister(tc.endpoints, nil), NewFakeNodeLister(tc.nodes, nil))
			translatedEvents, err := translator.Translate(&event)
			if err != nil {
				t.Fatalf(TranslateErrorFormat, err)
			}

			actualEventCount := len(translatedEvents)
			if actualEventCount != expectedEventCount {
				t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
			}

			assertExpectedServerCount(t, tc.expectedServerCount, translatedEvents)
		})
	}
}

func TestCreatedTranslateManyMixedPortsAndManyNodes(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		serviceType         v1.ServiceType
		nodes               []*v1.Node
		endpoints           []*discovery.EndpointSlice
		expectedServerCount int
	}{
		"nodePort": {
			serviceType:         v1.ServiceTypeNodePort,
			nodes:               generateNodes(ManyNodes),
			expectedServerCount: ManyNodes,
		},
		"clusterIP": {
			serviceType:         v1.ServiceTypeClusterIP,
			endpoints:           generateEndpointSlices(ManyEndpointSlices, 6, 2),
			expectedServerCount: ManyEndpointSlices,
		},
	}

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			const expectedEventCount = 2
			const portCount = 6
			const updatablePortCount = 2

			ports := generateUpdatablePorts(portCount, updatablePortCount)
			service := serviceWithPorts(tc.serviceType, ports)
			event := buildCreatedEvent(service)

			translator := NewTranslator(NewFakeEndpointSliceLister(tc.endpoints, nil), NewFakeNodeLister(tc.nodes, nil))
			translatedEvents, err := translator.Translate(&event)
			if err != nil {
				t.Fatalf(TranslateErrorFormat, err)
			}

			actualEventCount := len(translatedEvents)
			if actualEventCount != expectedEventCount {
				t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
			}

			assertExpectedServerCount(t, ManyNodes, translatedEvents)
		})
	}
}

/*
 * Updated Event Tests
 */

func TestUpdatedTranslateNoPorts(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		serviceType         v1.ServiceType
		nodes               []*v1.Node
		endpoints           []*discovery.EndpointSlice
		expectedServerCount int
	}{
		"nodePort": {
			serviceType:         v1.ServiceTypeNodePort,
			nodes:               generateNodes(OneNode),
			expectedServerCount: OneNode,
		},
		"clusterIP": {
			serviceType:         v1.ServiceTypeClusterIP,
			endpoints:           generateEndpointSlices(OneEndpointSlice, 0, 0),
			expectedServerCount: OneEndpointSlice,
		},
	}

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			const expectedEventCount = 0

			service := defaultService(tc.serviceType)
			event := buildUpdatedEvent(service)

			translator := NewTranslator(NewFakeEndpointSliceLister(tc.endpoints, nil), NewFakeNodeLister(tc.nodes, nil))
			translatedEvents, err := translator.Translate(&event)
			if err != nil {
				t.Fatalf(TranslateErrorFormat, err)
			}

			actualEventCount := len(translatedEvents)
			if actualEventCount != expectedEventCount {
				t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
			}
		})
	}
}

func TestUpdatedTranslateNoInterestingPorts(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		serviceType         v1.ServiceType
		nodes               []*v1.Node
		endpoints           []*discovery.EndpointSlice
		expectedServerCount int
	}{
		"nodePort": {
			serviceType:         v1.ServiceTypeNodePort,
			nodes:               generateNodes(OneNode),
			expectedServerCount: OneNode,
		},
		"clusterIP": {
			serviceType:         v1.ServiceTypeClusterIP,
			endpoints:           generateEndpointSlices(OneEndpointSlice, 1, 0),
			expectedServerCount: OneEndpointSlice,
		},
	}

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			const expectedEventCount = 0
			const portCount = 1

			ports := generateUpdatablePorts(portCount, 0)
			service := serviceWithPorts(tc.serviceType, ports)
			event := buildUpdatedEvent(service)

			translator := NewTranslator(NewFakeEndpointSliceLister(tc.endpoints, nil), NewFakeNodeLister(tc.nodes, nil))
			translatedEvents, err := translator.Translate(&event)
			if err != nil {
				t.Fatalf(TranslateErrorFormat, err)
			}

			actualEventCount := len(translatedEvents)
			if actualEventCount != expectedEventCount {
				t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
			}
		})
	}
}

func TestUpdatedTranslateOneInterestingPort(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		serviceType         v1.ServiceType
		nodes               []*v1.Node
		endpoints           []*discovery.EndpointSlice
		expectedServerCount int
	}{
		"nodePort": {
			serviceType:         v1.ServiceTypeNodePort,
			nodes:               generateNodes(OneNode),
			expectedServerCount: OneNode,
		},
		"clusterIP": {
			serviceType:         v1.ServiceTypeClusterIP,
			endpoints:           generateEndpointSlices(OneEndpointSlice, 1, 1),
			expectedServerCount: OneEndpointSlice,
		},
	}

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			const expectedEventCount = 1
			const portCount = 1

			ports := generatePorts(portCount)
			service := serviceWithPorts(tc.serviceType, ports)
			event := buildUpdatedEvent(service)

			translator := NewTranslator(NewFakeEndpointSliceLister(tc.endpoints, nil), NewFakeNodeLister(tc.nodes, nil))
			translatedEvents, err := translator.Translate(&event)
			if err != nil {
				t.Fatalf(TranslateErrorFormat, err)
			}

			actualEventCount := len(translatedEvents)
			if actualEventCount != expectedEventCount {
				t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
			}

			assertExpectedServerCount(t, OneNode, translatedEvents)
		})
	}
}

//nolint:dupl
func TestUpdatedTranslateManyInterestingPorts(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		serviceType         v1.ServiceType
		nodes               []*v1.Node
		endpoints           []*discovery.EndpointSlice
		expectedServerCount int
	}{
		"nodePort": {
			serviceType:         v1.ServiceTypeNodePort,
			nodes:               generateNodes(OneNode),
			expectedServerCount: OneNode,
		},
		"clusterIP": {
			serviceType:         v1.ServiceTypeClusterIP,
			endpoints:           generateEndpointSlices(OneEndpointSlice, 4, 4),
			expectedServerCount: OneEndpointSlice,
		},
	}

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			const expectedEventCount = 4
			const portCount = 4

			ports := generatePorts(portCount)
			service := serviceWithPorts(tc.serviceType, ports)
			event := buildUpdatedEvent(service)

			translator := NewTranslator(NewFakeEndpointSliceLister(tc.endpoints, nil), NewFakeNodeLister(tc.nodes, nil))
			translatedEvents, err := translator.Translate(&event)
			if err != nil {
				t.Fatalf(TranslateErrorFormat, err)
			}

			actualEventCount := len(translatedEvents)
			if actualEventCount != expectedEventCount {
				t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
			}

			assertExpectedServerCount(t, tc.expectedServerCount, translatedEvents)
		})
	}
}

//nolint:dupl
func TestUpdatedTranslateManyMixedPorts(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		serviceType         v1.ServiceType
		nodes               []*v1.Node
		endpoints           []*discovery.EndpointSlice
		expectedServerCount int
	}{
		"nodePort": {
			serviceType:         v1.ServiceTypeNodePort,
			nodes:               generateNodes(OneNode),
			expectedServerCount: OneNode,
		},
		"clusterIP": {
			serviceType:         v1.ServiceTypeClusterIP,
			endpoints:           generateEndpointSlices(OneEndpointSlice, 6, 2),
			expectedServerCount: OneEndpointSlice,
		},
	}

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			const expectedEventCount = 2
			const portCount = 6
			const updatablePortCount = 2

			ports := generateUpdatablePorts(portCount, updatablePortCount)
			service := serviceWithPorts(tc.serviceType, ports)
			event := buildUpdatedEvent(service)

			translator := NewTranslator(NewFakeEndpointSliceLister(tc.endpoints, nil), NewFakeNodeLister(tc.nodes, nil))
			translatedEvents, err := translator.Translate(&event)
			if err != nil {
				t.Fatalf(TranslateErrorFormat, err)
			}

			actualEventCount := len(translatedEvents)
			if actualEventCount != expectedEventCount {
				t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
			}

			assertExpectedServerCount(t, tc.expectedServerCount, translatedEvents)
		})
	}
}

//nolint:dupl
func TestUpdatedTranslateManyMixedPortsAndManyNodes(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		serviceType         v1.ServiceType
		nodes               []*v1.Node
		endpoints           []*discovery.EndpointSlice
		expectedServerCount int
	}{
		"nodePort": {
			serviceType:         v1.ServiceTypeNodePort,
			nodes:               generateNodes(ManyNodes),
			expectedServerCount: ManyNodes,
		},
		"clusterIP": {
			serviceType:         v1.ServiceTypeClusterIP,
			endpoints:           generateEndpointSlices(ManyEndpointSlices, 6, 2),
			expectedServerCount: ManyEndpointSlices,
		},
	}

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			const expectedEventCount = 2
			const portCount = 6
			const updatablePortCount = 2

			ports := generateUpdatablePorts(portCount, updatablePortCount)
			service := serviceWithPorts(tc.serviceType, ports)
			event := buildUpdatedEvent(service)

			translator := NewTranslator(NewFakeEndpointSliceLister(tc.endpoints, nil), NewFakeNodeLister(tc.nodes, nil))
			translatedEvents, err := translator.Translate(&event)
			if err != nil {
				t.Fatalf(TranslateErrorFormat, err)
			}

			actualEventCount := len(translatedEvents)
			if actualEventCount != expectedEventCount {
				t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
			}

			assertExpectedServerCount(t, tc.expectedServerCount, translatedEvents)
		})
	}
}

/*
 * Deleted Event Tests
 */

func TestDeletedTranslateNoPortsAndNoNodes(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		serviceType v1.ServiceType
		nodes       []*v1.Node
		endpoints   []*discovery.EndpointSlice
	}{
		"nodePort": {
			serviceType: v1.ServiceTypeNodePort,
		},
		"clusterIP": {
			serviceType: v1.ServiceTypeClusterIP,
		},
	}

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			const expectedEventCount = 0

			service := defaultService(tc.serviceType)
			event := buildDeletedEvent(service)

			translator := NewTranslator(NewFakeEndpointSliceLister(tc.endpoints, nil), NewFakeNodeLister(tc.nodes, nil))
			translatedEvents, err := translator.Translate(&event)
			if err != nil {
				t.Fatalf(TranslateErrorFormat, err)
			}

			actualEventCount := len(translatedEvents)
			if actualEventCount != expectedEventCount {
				t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
			}

			assertExpectedServerCount(t, 0, translatedEvents)
		})
	}
}

func TestDeletedTranslateNoInterestingPortsAndNoNodes(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		serviceType v1.ServiceType
		nodes       []*v1.Node
		endpoints   []*discovery.EndpointSlice
	}{
		"nodePort": {
			serviceType: v1.ServiceTypeNodePort,
		},
		"clusterIP": {
			serviceType: v1.ServiceTypeClusterIP,
		},
	}

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			const expectedEventCount = 0
			const portCount = 1

			ports := generateUpdatablePorts(portCount, 0)
			service := serviceWithPorts(tc.serviceType, ports)
			event := buildDeletedEvent(service)

			translator := NewTranslator(NewFakeEndpointSliceLister(tc.endpoints, nil), NewFakeNodeLister(tc.nodes, nil))
			translatedEvents, err := translator.Translate(&event)
			if err != nil {
				t.Fatalf(TranslateErrorFormat, err)
			}

			actualEventCount := len(translatedEvents)
			if actualEventCount != expectedEventCount {
				t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
			}

			assertExpectedServerCount(t, 0, translatedEvents)
		})
	}
}

//nolint:dupl
func TestDeletedTranslateOneInterestingPortAndNoNodes(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		serviceType v1.ServiceType
		nodes       []*v1.Node
		endpoints   []*discovery.EndpointSlice
	}{
		"nodePort": {
			serviceType: v1.ServiceTypeNodePort,
		},
		"clusterIP": {
			serviceType: v1.ServiceTypeClusterIP,
			endpoints:   generateEndpointSlices(0, 1, 1),
		},
	}

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			const expectedEventCount = 1
			const portCount = 1

			ports := generatePorts(portCount)
			service := serviceWithPorts(tc.serviceType, ports)
			event := buildDeletedEvent(service)

			translator := NewTranslator(NewFakeEndpointSliceLister(tc.endpoints, nil), NewFakeNodeLister(tc.nodes, nil))
			translatedEvents, err := translator.Translate(&event)
			if err != nil {
				t.Fatalf(TranslateErrorFormat, err)
			}

			actualEventCount := len(translatedEvents)
			if actualEventCount != expectedEventCount {
				t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
			}

			assertExpectedServerCount(t, 0, translatedEvents)
		})
	}
}

//nolint:dupl
func TestDeletedTranslateManyInterestingPortsAndNoNodes(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		serviceType v1.ServiceType
		nodes       []*v1.Node
		endpoints   []*discovery.EndpointSlice
	}{
		"nodePort": {
			serviceType: v1.ServiceTypeNodePort,
		},
		"clusterIP": {
			serviceType: v1.ServiceTypeClusterIP,
			endpoints:   generateEndpointSlices(0, 4, 4),
		},
	}

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			const portCount = 4
			const expectedEventCount = 4

			ports := generatePorts(portCount)
			service := serviceWithPorts(tc.serviceType, ports)
			event := buildDeletedEvent(service)

			translator := NewTranslator(NewFakeEndpointSliceLister(tc.endpoints, nil), NewFakeNodeLister(tc.nodes, nil))
			translatedEvents, err := translator.Translate(&event)
			if err != nil {
				t.Fatalf(TranslateErrorFormat, err)
			}

			actualEventCount := len(translatedEvents)
			if actualEventCount != expectedEventCount {
				t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
			}

			assertExpectedServerCount(t, 0, translatedEvents)
		})
	}
}

func TestDeletedTranslateManyMixedPortsAndNoNodes(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		serviceType v1.ServiceType
		nodes       []*v1.Node
		endpoints   []*discovery.EndpointSlice
	}{
		"nodePort": {
			serviceType: v1.ServiceTypeNodePort,
		},
		"clusterIP": {
			serviceType: v1.ServiceTypeClusterIP,
			endpoints:   generateEndpointSlices(0, 6, 2),
		},
	}

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			const portCount = 6
			const updatablePortCount = 2
			const expectedEventCount = 2

			ports := generateUpdatablePorts(portCount, updatablePortCount)
			service := serviceWithPorts(tc.serviceType, ports)
			event := buildDeletedEvent(service)

			translator := NewTranslator(NewFakeEndpointSliceLister(tc.endpoints, nil), NewFakeNodeLister(tc.nodes, nil))
			translatedEvents, err := translator.Translate(&event)
			if err != nil {
				t.Fatalf(TranslateErrorFormat, err)
			}

			actualEventCount := len(translatedEvents)
			if actualEventCount != expectedEventCount {
				t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
			}

			assertExpectedServerCount(t, 0, translatedEvents)
		})
	}
}

//nolint:dupl
func TestDeletedTranslateNoPortsAndOneNode(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		serviceType v1.ServiceType
		nodes       []*v1.Node
		endpoints   []*discovery.EndpointSlice
	}{
		"nodePort": {
			serviceType: v1.ServiceTypeNodePort,
			nodes:       generateNodes(OneNode),
		},
		"clusterIP": {
			serviceType: v1.ServiceTypeClusterIP,
			endpoints:   generateEndpointSlices(OneEndpointSlice, 0, 0),
		},
	}

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			const expectedEventCount = 0

			service := defaultService(tc.serviceType)
			event := buildDeletedEvent(service)

			translator := NewTranslator(NewFakeEndpointSliceLister(tc.endpoints, nil), NewFakeNodeLister(tc.nodes, nil))
			translatedEvents, err := translator.Translate(&event)
			if err != nil {
				t.Fatalf(TranslateErrorFormat, err)
			}

			actualEventCount := len(translatedEvents)
			if actualEventCount != expectedEventCount {
				t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
			}

			assertExpectedServerCount(t, 0, translatedEvents)
		})
	}
}

func TestDeletedTranslateNoInterestingPortsAndOneNode(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		serviceType v1.ServiceType
		nodes       []*v1.Node
		endpoints   []*discovery.EndpointSlice
	}{
		"nodePort": {
			serviceType: v1.ServiceTypeNodePort,
			nodes:       generateNodes(OneNode),
		},
		"clusterIP": {
			serviceType: v1.ServiceTypeClusterIP,
			endpoints:   generateEndpointSlices(OneEndpointSlice, 1, 0),
		},
	}

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			const portCount = 1
			const expectedEventCount = 0

			ports := generateUpdatablePorts(portCount, 0)
			service := serviceWithPorts(tc.serviceType, ports)
			event := buildDeletedEvent(service)

			translator := NewTranslator(NewFakeEndpointSliceLister(tc.endpoints, nil), NewFakeNodeLister(tc.nodes, nil))
			translatedEvents, err := translator.Translate(&event)
			if err != nil {
				t.Fatalf(TranslateErrorFormat, err)
			}

			actualEventCount := len(translatedEvents)
			if actualEventCount != expectedEventCount {
				t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
			}

			assertExpectedServerCount(t, 0, translatedEvents)
		})
	}
}

//nolint:dupl
func TestDeletedTranslateOneInterestingPortAndOneNode(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		serviceType v1.ServiceType
		nodes       []*v1.Node
		endpoints   []*discovery.EndpointSlice
	}{
		"nodePort": {
			serviceType: v1.ServiceTypeNodePort,
			nodes:       generateNodes(OneNode),
		},
		"clusterIP": {
			serviceType: v1.ServiceTypeClusterIP,
			endpoints:   generateEndpointSlices(OneEndpointSlice, 1, 1),
		},
	}

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			const portCount = 1
			const expectedEventCount = 1

			ports := generatePorts(portCount)
			service := serviceWithPorts(tc.serviceType, ports)
			event := buildDeletedEvent(service)

			translator := NewTranslator(NewFakeEndpointSliceLister(tc.endpoints, nil), NewFakeNodeLister(tc.nodes, nil))
			translatedEvents, err := translator.Translate(&event)
			if err != nil {
				t.Fatalf(TranslateErrorFormat, err)
			}

			actualEventCount := len(translatedEvents)
			if actualEventCount != expectedEventCount {
				t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
			}

			assertExpectedServerCount(t, 0, translatedEvents)
		})
	}
}

//nolint:dupl
func TestDeletedTranslateManyInterestingPortsAndOneNode(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		serviceType v1.ServiceType
		nodes       []*v1.Node
		endpoints   []*discovery.EndpointSlice
	}{
		"nodePort": {
			serviceType: v1.ServiceTypeNodePort,
			nodes:       generateNodes(OneNode),
		},
		"clusterIP": {
			serviceType: v1.ServiceTypeClusterIP,
			endpoints:   generateEndpointSlices(OneEndpointSlice, 4, 4),
		},
	}

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			const portCount = 4
			const expectedEventCount = 4

			ports := generatePorts(portCount)
			service := serviceWithPorts(tc.serviceType, ports)
			event := buildDeletedEvent(service)

			translator := NewTranslator(NewFakeEndpointSliceLister(tc.endpoints, nil), NewFakeNodeLister(tc.nodes, nil))
			translatedEvents, err := translator.Translate(&event)
			if err != nil {
				t.Fatalf(TranslateErrorFormat, err)
			}

			actualEventCount := len(translatedEvents)
			if actualEventCount != expectedEventCount {
				t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
			}

			assertExpectedServerCount(t, 0, translatedEvents)
		})
	}
}

func TestDeletedTranslateManyMixedPortsAndOneNode(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		serviceType v1.ServiceType
		nodes       []*v1.Node
		endpoints   []*discovery.EndpointSlice
	}{
		"nodePort": {
			serviceType: v1.ServiceTypeNodePort,
			nodes:       generateNodes(OneNode),
		},
		"clusterIP": {
			serviceType: v1.ServiceTypeClusterIP,
			endpoints:   generateEndpointSlices(OneEndpointSlice, 6, 2),
		},
	}

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			const portCount = 6
			const updatablePortCount = 2
			const expectedEventCount = 2

			ports := generateUpdatablePorts(portCount, updatablePortCount)
			service := serviceWithPorts(tc.serviceType, ports)
			event := buildDeletedEvent(service)

			translator := NewTranslator(NewFakeEndpointSliceLister(tc.endpoints, nil), NewFakeNodeLister(tc.nodes, nil))
			translatedEvents, err := translator.Translate(&event)
			if err != nil {
				t.Fatalf(TranslateErrorFormat, err)
			}

			actualEventCount := len(translatedEvents)
			if actualEventCount != expectedEventCount {
				t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
			}

			assertExpectedServerCount(t, 0, translatedEvents)
		})
	}
}

//nolint:dupl
func TestDeletedTranslateNoPortsAndManyNodes(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		serviceType v1.ServiceType
		nodes       []*v1.Node
		endpoints   []*discovery.EndpointSlice
	}{
		"nodePort": {
			serviceType: v1.ServiceTypeNodePort,
			nodes:       generateNodes(ManyNodes),
		},
		"clusterIP": {
			serviceType: v1.ServiceTypeClusterIP,
			endpoints:   generateEndpointSlices(ManyEndpointSlices, 0, 0),
		},
	}

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			const expectedEventCount = 0

			service := defaultService(tc.serviceType)
			event := buildDeletedEvent(service)

			translator := NewTranslator(NewFakeEndpointSliceLister(tc.endpoints, nil), NewFakeNodeLister(tc.nodes, nil))
			translatedEvents, err := translator.Translate(&event)
			if err != nil {
				t.Fatalf(TranslateErrorFormat, err)
			}

			actualEventCount := len(translatedEvents)
			if actualEventCount != expectedEventCount {
				t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
			}

			assertExpectedServerCount(t, 0, translatedEvents)
		})
	}
}

func TestDeletedTranslateNoInterestingPortsAndManyNodes(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		serviceType v1.ServiceType
		nodes       []*v1.Node
		endpoints   []*discovery.EndpointSlice
	}{
		"nodePort": {
			serviceType: v1.ServiceTypeNodePort,
			nodes:       generateNodes(ManyNodes),
		},
		"clusterIP": {
			serviceType: v1.ServiceTypeClusterIP,
			endpoints:   generateEndpointSlices(ManyEndpointSlices, 1, 0),
		},
	}

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			const portCount = 1
			const updatablePortCount = 0
			const expectedEventCount = updatablePortCount * ManyNodes

			ports := generateUpdatablePorts(portCount, updatablePortCount)
			service := serviceWithPorts(tc.serviceType, ports)
			event := buildDeletedEvent(service)

			translator := NewTranslator(NewFakeEndpointSliceLister(tc.endpoints, nil), NewFakeNodeLister(tc.nodes, nil))
			translatedEvents, err := translator.Translate(&event)
			if err != nil {
				t.Fatalf(TranslateErrorFormat, err)
			}

			actualEventCount := len(translatedEvents)
			if actualEventCount != expectedEventCount {
				t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
			}

			assertExpectedServerCount(t, 0, translatedEvents)
		})
	}
}

//nolint:dupl
func TestDeletedTranslateOneInterestingPortAndManyNodes(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		serviceType v1.ServiceType
		nodes       []*v1.Node
		endpoints   []*discovery.EndpointSlice
	}{
		"nodePort": {
			serviceType: v1.ServiceTypeNodePort,
			nodes:       generateNodes(ManyNodes),
		},
		"clusterIP": {
			serviceType: v1.ServiceTypeClusterIP,
			endpoints:   generateEndpointSlices(ManyEndpointSlices, 1, 1),
		},
	}

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			const portCount = 1
			const expectedEventCount = 1

			ports := generatePorts(portCount)
			service := serviceWithPorts(tc.serviceType, ports)
			event := buildDeletedEvent(service)

			translator := NewTranslator(NewFakeEndpointSliceLister(tc.endpoints, nil), NewFakeNodeLister(tc.nodes, nil))
			translatedEvents, err := translator.Translate(&event)
			if err != nil {
				t.Fatalf(TranslateErrorFormat, err)
			}

			actualEventCount := len(translatedEvents)
			if actualEventCount != expectedEventCount {
				t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
			}

			assertExpectedServerCount(t, 0, translatedEvents)
		})
	}
}

//nolint:dupl
func TestDeletedTranslateManyInterestingPortsAndManyNodes(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		serviceType v1.ServiceType
		nodes       []*v1.Node
		endpoints   []*discovery.EndpointSlice
	}{
		"nodePort": {
			serviceType: v1.ServiceTypeNodePort,
			nodes:       generateNodes(ManyNodes),
		},
		"clusterIP": {
			serviceType: v1.ServiceTypeClusterIP,
			endpoints:   generateEndpointSlices(ManyEndpointSlices, 4, 4),
		},
	}

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			const portCount = 4
			const expectedEventCount = 4

			ports := generatePorts(portCount)
			service := serviceWithPorts(tc.serviceType, ports)
			event := buildDeletedEvent(service)

			translator := NewTranslator(NewFakeEndpointSliceLister(tc.endpoints, nil), NewFakeNodeLister(tc.nodes, nil))
			translatedEvents, err := translator.Translate(&event)
			if err != nil {
				t.Fatalf(TranslateErrorFormat, err)
			}

			actualEventCount := len(translatedEvents)
			if actualEventCount != expectedEventCount {
				t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
			}

			assertExpectedServerCount(t, 0, translatedEvents)
		})
	}
}

func TestDeletedTranslateManyMixedPortsAndManyNodes(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		serviceType v1.ServiceType
		nodes       []*v1.Node
		endpoints   []*discovery.EndpointSlice
	}{
		"nodePort": {
			serviceType: v1.ServiceTypeNodePort,
			nodes:       generateNodes(ManyNodes),
		},
		"clusterIP": {
			serviceType: v1.ServiceTypeClusterIP,
			endpoints:   generateEndpointSlices(ManyEndpointSlices, 6, 2),
		},
	}

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			const portCount = 6
			const updatablePortCount = 2
			const expectedEventCount = 2

			ports := generateUpdatablePorts(portCount, updatablePortCount)
			service := serviceWithPorts(tc.serviceType, ports)
			event := buildDeletedEvent(service)

			translator := NewTranslator(NewFakeEndpointSliceLister(tc.endpoints, nil), NewFakeNodeLister(tc.nodes, nil))
			translatedEvents, err := translator.Translate(&event)
			if err != nil {
				t.Fatalf(TranslateErrorFormat, err)
			}

			actualEventCount := len(translatedEvents)
			if actualEventCount != expectedEventCount {
				t.Fatalf(AssertionFailureFormat, expectedEventCount, actualEventCount)
			}

			assertExpectedServerCount(t, 0, translatedEvents)
		})
	}
}

func assertExpectedServerCount(t *testing.T, expectedCount int, events core.ServerUpdateEvents) {
	for _, translatedEvent := range events {
		serverCount := len(translatedEvent.UpstreamServers)
		if serverCount != expectedCount {
			t.Fatalf("expected %d servers, got %d", expectedCount, serverCount)
		}
	}
}

func defaultService(serviceType v1.ServiceType) *v1.Service {
	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "default-service",
			Labels: map[string]string{"kubernetes.io/service-name": "default-service"},
		},
		Spec: v1.ServiceSpec{
			Type: serviceType,
		},
	}
}

func serviceWithPorts(serviceType v1.ServiceType, ports []v1.ServicePort) *v1.Service {
	return &v1.Service{
		Spec: v1.ServiceSpec{
			Type:  serviceType,
			Ports: ports,
		},
	}
}

func buildCreatedEvent(service *v1.Service) core.Event {
	return buildEvent(core.Created, service)
}

func buildDeletedEvent(service *v1.Service) core.Event {
	return buildEvent(core.Deleted, service)
}

func buildUpdatedEvent(service *v1.Service) core.Event {
	return buildEvent(core.Updated, service)
}

func buildEvent(eventType core.EventType, service *v1.Service) core.Event {
	previousService := defaultService(service.Spec.Type)

	event := core.NewEvent(eventType, service, previousService)
	event.Service.Name = "default-service"
	return event
}

func generateNodes(count int) (nodes []*v1.Node) {
	for i := 0; i < count; i++ {
		nodes = append(nodes, &v1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("node%d", i),
			},
			Status: v1.NodeStatus{
				Addresses: []v1.NodeAddress{
					{
						Type:    v1.NodeInternalIP,
						Address: fmt.Sprintf("10.0.0.%v", i),
					},
				},
			},
		})
	}

	return nodes
}

func generateEndpointSlices(endpointCount, portCount, updatablePortCount int,
) (endpointSlices []*discovery.EndpointSlice) {
	servicePorts := generateUpdatablePorts(portCount, updatablePortCount)

	ports := make([]discovery.EndpointPort, 0, len(servicePorts))
	for _, servicePort := range servicePorts {
		ports = append(ports, discovery.EndpointPort{
			Name: pointer.To(servicePort.Name),
			Port: pointer.To(int32(8080)),
		})
	}

	var endpoints []discovery.Endpoint
	for i := 0; i < endpointCount; i++ {
		endpoints = append(endpoints, discovery.Endpoint{
			Addresses: []string{
				fmt.Sprintf("10.0.0.%v", i),
			},
		})
	}

	endpointSlices = append(endpointSlices, &discovery.EndpointSlice{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "endpointSlice",
			Labels: map[string]string{"kubernetes.io/service-name": "default-service"},
		},
		Endpoints: endpoints,
		Ports:     ports,
	})

	return endpointSlices
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

func NewFakeEndpointSliceLister(list []*discovery.EndpointSlice, err error) discoverylisters.EndpointSliceLister {
	return &endpointSliceLister{
		list: list,
		err:  err,
	}
}

func NewFakeNodeLister(list []*v1.Node, err error) corelisters.NodeLister {
	return &nodeLister{
		list: list,
		err:  err,
	}
}

type nodeLister struct {
	list []*v1.Node
	err  error
}

func (l *nodeLister) List(selector labels.Selector) (ret []*v1.Node, err error) {
	return l.list, l.err
}

// currently unused
func (l *nodeLister) Get(name string) (*v1.Node, error) {
	return nil, nil
}

type endpointSliceLister struct {
	list []*discovery.EndpointSlice
	err  error
}

func (l *endpointSliceLister) List(selector labels.Selector) (ret []*discovery.EndpointSlice, err error) {
	return l.list, l.err
}

func (l *endpointSliceLister) Get(name string) (*discovery.EndpointSlice, error) {
	return nil, nil
}

func (l *endpointSliceLister) EndpointSlices(name string) discoverylisters.EndpointSliceNamespaceLister {
	return l
}
