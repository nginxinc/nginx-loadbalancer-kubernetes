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
	NginxUpstreams  []nginxClient.UpstreamServer
}

func NewEvent(eventType EventType, service *v1.Service, previousService *v1.Service) Event {
	return Event{
		Type:            eventType,
		Service:         service,
		PreviousService: previousService,
		NginxUpstreams:  []nginxClient.UpstreamServer{},
	}
}
