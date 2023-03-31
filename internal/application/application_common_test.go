/*
 * Copyright (c) 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package application

import (
	"errors"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/core"
	"github.com/nginxinc/kubernetes-nginx-ingress/test/mocks"
	nginxClient2 "github.com/nginxinc/nginx-plus-go-client/client"
)

const (
	deletedEventType = core.Deleted
	createEventType  = core.Created
	upstreamName     = "upstreamName"
	server           = "server"
)

func buildTerrorizingBorderClient(clientType string) (Interface, *mocks.MockNginxClient, error) {
	nginxClient := mocks.NewErroringMockClient(errors.New(`something went horribly horribly wrong`))
	bc, err := NewBorderClient(clientType, nginxClient)

	return bc, nginxClient, err
}

func buildBorderClient(clientType string) (Interface, *mocks.MockNginxClient, error) {
	nginxClient := mocks.NewMockNginxClient()
	bc, err := NewBorderClient(clientType, nginxClient)

	return bc, nginxClient, err
}

func buildServerUpdateEvent(eventType core.EventType, clientType string) *core.ServerUpdateEvent {
	httpServers := []nginxClient2.UpstreamServer{
		{
			Server: server,
		},
	}

	tcpServers := []nginxClient2.StreamUpstreamServer{
		{
			Server: server,
		},
	}
	return core.NewServerUpdateEvent(eventType, upstreamName, clientType, tcpServers, httpServers)
}
