package core

import (
	v1 "k8s.io/api/core/v1"
	"testing"
)

func TestNewEvent(t *testing.T) {
	expectedType := Created
	expectedService := &v1.Service{}
	expectedPreviousService := &v1.Service{}
	expectedNodeIps := []string{"127.0.0.1"}

	event := NewEvent(expectedType, expectedService, expectedPreviousService, expectedNodeIps)

	if event.Type != expectedType {
		t.Errorf("expected Created, got %v", event.Type)
	}

	if event.Service != expectedService {
		t.Errorf("expected service, got %#v", event.Service)
	}

	if event.PreviousService != expectedPreviousService {
		t.Errorf("expected previous service, got %#v", event.PreviousService)
	}

	if event.NodeIps[0] != expectedNodeIps[0] {
		t.Errorf("expected node ips, got %#v", event.NodeIps)
	}
}
