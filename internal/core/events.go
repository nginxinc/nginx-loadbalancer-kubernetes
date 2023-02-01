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
}

type ServerUpdateEvent struct {
	id           string
	UpstreamName string
	Servers      []nginxClient.StreamUpstreamServer
}

type ServerUpdateEvents = []*ServerUpdateEvent

func NewEvent(eventType EventType, service *v1.Service, previousService *v1.Service) Event {
	return Event{
		Type:            eventType,
		Service:         service,
		PreviousService: previousService,
	}
}

func NewServerUpdateEvent(upstreamName string, servers []nginxClient.StreamUpstreamServer) *ServerUpdateEvent {
	return &ServerUpdateEvent{
		UpstreamName: upstreamName,
		Servers:      servers,
	}
}
