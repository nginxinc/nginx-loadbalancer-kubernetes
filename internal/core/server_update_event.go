/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package core

type ServerUpdateEvent struct {
	ClientType      string
	Id              string
	NginxHost       string
	Type            EventType
	UpstreamName    string
	UpstreamServers UpstreamServers
}

type ServerUpdateEvents = []*ServerUpdateEvent

func NewServerUpdateEvent(eventType EventType, upstreamName string, clientType string, upstreamServers UpstreamServers) *ServerUpdateEvent {
	return &ServerUpdateEvent{
		ClientType:      clientType,
		Type:            eventType,
		UpstreamName:    upstreamName,
		UpstreamServers: upstreamServers,
	}
}

func ServerUpdateEventWithIdAndHost(event *ServerUpdateEvent, id string, nginxHost string) *ServerUpdateEvent {
	return &ServerUpdateEvent{
		ClientType:      event.ClientType,
		Id:              id,
		NginxHost:       nginxHost,
		Type:            event.Type,
		UpstreamName:    event.UpstreamName,
		UpstreamServers: event.UpstreamServers,
	}
}

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
