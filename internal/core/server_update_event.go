/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package core

// ServerUpdateEvent is an internal representation of an event. The Translator produces these events
// from Events received from the Handler. These are then consumed by the Synchronizer and passed along to
// the appropriate BorderClient.
type ServerUpdateEvent struct {

	// ClientType is the type of BorderClient that should handle this event. This is configured via Service Annotations.
	// See application_constants.go for the list of supported types.
	ClientType string

	// Id is the unique identifier for this event.
	ID string

	// NginxHost is the host name of the NGINX Plus instance that should handle this event.
	NginxHost string

	// Type is the type of event. See EventType for the list of supported types.
	Type EventType

	// UpstreamName is the name of the upstream in the Border Server.
	UpstreamName string

	// UpstreamServers is the list of servers in the Upstream.
	UpstreamServers UpstreamServers
}

// ServerUpdateEvents is a list of ServerUpdateEvent.
type ServerUpdateEvents = []*ServerUpdateEvent

// NewServerUpdateEvent creates a new ServerUpdateEvent.
func NewServerUpdateEvent(
	eventType EventType,
	upstreamName string,
	clientType string,
	upstreamServers UpstreamServers,
) *ServerUpdateEvent {
	return &ServerUpdateEvent{
		ClientType:      clientType,
		Type:            eventType,
		UpstreamName:    upstreamName,
		UpstreamServers: upstreamServers,
	}
}

// ServerUpdateEventWithIDAndHost creates a new ServerUpdateEvent with the specified Id and Host.
func ServerUpdateEventWithIDAndHost(event *ServerUpdateEvent, id string, nginxHost string) *ServerUpdateEvent {
	return &ServerUpdateEvent{
		ClientType:      event.ClientType,
		ID:              id,
		NginxHost:       nginxHost,
		Type:            event.Type,
		UpstreamName:    event.UpstreamName,
		UpstreamServers: event.UpstreamServers,
	}
}

// TypeName returns the string representation of the EventType.
func (e *ServerUpdateEvent) TypeName() string {
	switch e.Type {
	case Created:
		return "Created"
	case Updated:
		return "Updated"
	case Deleted:
		return "Deleted"
	default:
		return "Unknown"
	}
}
