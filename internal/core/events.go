package core

import (
	nginxClient "github.com/nginxinc/nginx-plus-go-client/client"
	v1 "k8s.io/api/core/v1"
)

type EventType int

const (
	Created EventType = iota
	Updated
	Deleted
)

type Event struct {
	Type            EventType
	Service         *v1.Service
	PreviousService *v1.Service
	NodeIps         []string
}

type ServerUpdateEvent struct {
	Id           string
	NginxHost    string
	Type         EventType
	UpstreamName string
	Servers      []nginxClient.StreamUpstreamServer
}

type ServerUpdateEvents = []*ServerUpdateEvent

func NewEvent(eventType EventType, service *v1.Service, previousService *v1.Service, nodeIps []string) Event {
	return Event{
		Type:            eventType,
		Service:         service,
		PreviousService: previousService,
		NodeIps:         nodeIps,
	}
}

func NewServerUpdateEvent(eventType EventType, upstreamName string, servers []nginxClient.StreamUpstreamServer) *ServerUpdateEvent {
	return &ServerUpdateEvent{
		Type:         eventType,
		UpstreamName: upstreamName,
		Servers:      servers,
	}
}

func ServerUpdateEventWithIdAndHost(event *ServerUpdateEvent, id string, nginxHost string) *ServerUpdateEvent {
	return &ServerUpdateEvent{
		Id:           id,
		NginxHost:    nginxHost,
		Type:         event.Type,
		UpstreamName: event.UpstreamName,
		Servers:      event.Servers,
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
