/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package core

import (
	nginxClient "github.com/nginxinc/nginx-plus-go-client/client"
	"testing"
)

const clientType = "clientType"

var emptyStreamServers []nginxClient.StreamUpstreamServer
var emptyHttpServers []nginxClient.UpstreamServer

func TestServerUpdateEventWithIdAndHost(t *testing.T) {
	event := NewServerUpdateEvent(Created, "upstream", clientType, emptyStreamServers, emptyHttpServers)

	if event.Id != "" {
		t.Errorf("expected empty Id, got %s", event.Id)
	}

	if event.NginxHost != "" {
		t.Errorf("expected empty NginxHost, got %s", event.NginxHost)
	}

	eventWithIdAndHost := ServerUpdateEventWithIdAndHost(event, "id", "host")

	if eventWithIdAndHost.Id != "id" {
		t.Errorf("expected Id to be 'id', got %s", eventWithIdAndHost.Id)
	}

	if eventWithIdAndHost.NginxHost != "host" {
		t.Errorf("expected NginxHost to be 'host', got %s", eventWithIdAndHost.NginxHost)
	}

	if eventWithIdAndHost.ClientType != clientType {
		t.Errorf("expected ClientType to be '%s', got %s", clientType, eventWithIdAndHost.ClientType)
	}
}

func TestTypeNameCreated(t *testing.T) {
	event := NewServerUpdateEvent(Created, "upstream", clientType, emptyStreamServers, emptyHttpServers)

	if event.TypeName() != "Created" {
		t.Errorf("expected 'Created', got %s", event.TypeName())
	}
}

func TestTypeNameUpdated(t *testing.T) {
	event := NewServerUpdateEvent(Updated, "upstream", clientType, emptyStreamServers, emptyHttpServers)

	if event.TypeName() != "Updated" {
		t.Errorf("expected 'Updated', got %s", event.TypeName())
	}
}

func TestTypeNameDeleted(t *testing.T) {
	event := NewServerUpdateEvent(Deleted, "upstream", clientType, emptyStreamServers, emptyHttpServers)

	if event.TypeName() != "Deleted" {
		t.Errorf("expected 'Deleted', got %s", event.TypeName())
	}
}

func TestTypeNameUnknown(t *testing.T) {
	event := NewServerUpdateEvent(EventType(100), "upstream", clientType, emptyStreamServers, emptyHttpServers)

	if event.TypeName() != "Unknown" {
		t.Errorf("expected 'Unknown', got %s", event.TypeName())
	}
}
