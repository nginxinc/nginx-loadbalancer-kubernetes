package core

import (
	nginxClient "github.com/nginxinc/nginx-plus-go-client/client"
	v1 "k8s.io/api/networking/v1"
)

type EventType int

const (
	Created EventType = iota
	Updated
	Deleted
)

type Event struct {
	Type            EventType
	Ingress         *v1.Ingress
	PreviousIngress *v1.Ingress
	NginxUpstream   *nginxClient.UpstreamServer
}

func NewEvent(eventType EventType, ingress *v1.Ingress, previousIngress *v1.Ingress) Event {
	return Event{
		Type:            eventType,
		Ingress:         ingress,
		PreviousIngress: previousIngress,
		NginxUpstream:   nil,
	}
}
