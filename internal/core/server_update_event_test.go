/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package core

import (
	"testing"
)

const clientType = "clientType"

var emptyUpstreamServers UpstreamServers

func TestServerUpdateEventWithIdAndHost(t *testing.T) {
	t.Parallel()
	event := NewServerUpdateEvent(Created, "upstream", clientType, emptyUpstreamServers)

	if event.ID != "" {
		t.Errorf("expected empty ID, got %s", event.ID)
	}

	if event.NginxHost != "" {
		t.Errorf("expected empty NginxHost, got %s", event.NginxHost)
	}

	eventWithIDAndHost := ServerUpdateEventWithIDAndHost(event, "id", "host")

	if eventWithIDAndHost.ID != "id" {
		t.Errorf("expected Id to be 'id', got %s", eventWithIDAndHost.ID)
	}

	if eventWithIDAndHost.NginxHost != "host" {
		t.Errorf("expected NginxHost to be 'host', got %s", eventWithIDAndHost.NginxHost)
	}

	if eventWithIDAndHost.ClientType != clientType {
		t.Errorf("expected ClientType to be '%s', got %s", clientType, eventWithIDAndHost.ClientType)
	}
}

func TestTypeNameCreated(t *testing.T) {
	t.Parallel()
	event := NewServerUpdateEvent(Created, "upstream", clientType, emptyUpstreamServers)

	if event.TypeName() != "Created" {
		t.Errorf("expected 'Created', got %s", event.TypeName())
	}
}

func TestTypeNameUpdated(t *testing.T) {
	t.Parallel()
	event := NewServerUpdateEvent(Updated, "upstream", clientType, emptyUpstreamServers)

	if event.TypeName() != "Updated" {
		t.Errorf("expected 'Updated', got %s", event.TypeName())
	}
}

func TestTypeNameDeleted(t *testing.T) {
	t.Parallel()
	event := NewServerUpdateEvent(Deleted, "upstream", clientType, emptyUpstreamServers)

	if event.TypeName() != "Deleted" {
		t.Errorf("expected 'Deleted', got %s", event.TypeName())
	}
}

func TestTypeNameUnknown(t *testing.T) {
	t.Parallel()
	event := NewServerUpdateEvent(EventType(100), "upstream", clientType, emptyUpstreamServers)

	if event.TypeName() != "Unknown" {
		t.Errorf("expected 'Unknown', got %s", event.TypeName())
	}
}
