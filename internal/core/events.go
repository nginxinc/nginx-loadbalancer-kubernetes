/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

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
	ClientType   string
	Id           string
	NginxHost    string
	Type         EventType
	UpstreamName string
	TcpServers   []nginxClient.StreamUpstreamServer
	HttpServers  []nginxClient.UpstreamServer
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

func NewServerUpdateEvent(eventType EventType, upstreamName string, clientType string, tcpServers []nginxClient.StreamUpstreamServer, httpServers []nginxClient.UpstreamServer) *ServerUpdateEvent {
	return &ServerUpdateEvent{
		ClientType:   clientType,
		Type:         eventType,
		UpstreamName: upstreamName,
		TcpServers:   tcpServers,
		HttpServers:  httpServers,
	}
}

func ServerUpdateEventWithIdAndHost(event *ServerUpdateEvent, id string, nginxHost string) *ServerUpdateEvent {
	return &ServerUpdateEvent{
		ClientType:   event.ClientType,
		Id:           id,
		NginxHost:    nginxHost,
		Type:         event.Type,
		UpstreamName: event.UpstreamName,
		TcpServers:   event.TcpServers,
		HttpServers:  event.HttpServers,
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
