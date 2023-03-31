/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package application

import (
	"errors"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/core"
	"github.com/nginxinc/kubernetes-nginx-ingress/test/mocks"
	nginxClient2 "github.com/nginxinc/nginx-plus-go-client/client"
	"testing"
)

const (
	clientType       = "http"
	deletedEventType = core.Deleted
	createEventType  = core.Created
	upstreamName     = "upstreamName"
	server           = "server"
)

var emptyStreamServers = []nginxClient2.StreamUpstreamServer{}

func TestHttpBorderClient_Delete(t *testing.T) {
	event := buildServerUpdateEvent(deletedEventType)
	borderClient, nginxClient, err := buildBorderClient()
	if err != nil {
		t.Fatalf(`error occurred creating a new border client: %v`, err)
	}

	err = borderClient.Delete(event)
	if err != nil {
		t.Fatalf(`error occurred deleting the nginx+ upstream server: %v`, err)
	}

	if !nginxClient.CalledFunctions["DeleteHTTPServer"] {
		t.Fatalf(`expected DeleteHTTPServer to be called`)
	}
}

func TestHttpBorderClient_Update(t *testing.T) {
	event := buildServerUpdateEvent(createEventType)
	borderClient, nginxClient, err := buildBorderClient()
	if err != nil {
		t.Fatalf(`error occurred creating a new border client: %v`, err)
	}

	err = borderClient.Update(event)
	if err != nil {
		t.Fatalf(`error occurred deleting the nginx+ upstream server: %v`, err)
	}

	if !nginxClient.CalledFunctions["UpdateHTTPServers"] {
		t.Fatalf(`expected UpdateHTTPServers to be called`)
	}
}

func TestHttpBorderClient_BadNginxClient(t *testing.T) {
	var emptyInterface interface{}
	_, err := NewBorderClient(clientType, emptyInterface)
	if err == nil {
		t.Fatalf(`expected an error to occur when creating a new border client`)
	}
}

func TestHttpBorderClient_DeleteReturnsError(t *testing.T) {
	event := buildServerUpdateEvent(deletedEventType)
	borderClient, _, err := buildTerrorizingBorderClient()
	if err != nil {
		t.Fatalf(`error occurred creating a new border client: %v`, err)
	}

	err = borderClient.Delete(event)

	if err == nil {
		t.Fatalf(`expected an error to occur when deleting the nginx+ upstream server`)
	}
}

func TestHttpBorderClient_UpdateReturnsError(t *testing.T) {
	event := buildServerUpdateEvent(createEventType)
	borderClient, _, err := buildTerrorizingBorderClient()
	if err != nil {
		t.Fatalf(`error occurred creating a new border client: %v`, err)
	}

	err = borderClient.Update(event)

	if err == nil {
		t.Fatalf(`expected an error to occur when deleting the nginx+ upstream server`)
	}
}

func buildTerrorizingBorderClient() (Interface, *mocks.MockNginxClient, error) {
	nginxClient := mocks.NewErroringMockClient(errors.New(`something went horribly horribly wrong`))
	bc, err := NewBorderClient(clientType, nginxClient)

	return bc, nginxClient, err
}

func buildBorderClient() (Interface, *mocks.MockNginxClient, error) {
	nginxClient := mocks.NewMockNginxClient()
	bc, err := NewBorderClient(clientType, nginxClient)

	return bc, nginxClient, err
}

func buildServerUpdateEvent(eventType core.EventType) core.ServerUpdateEvent {
	servers := []nginxClient2.UpstreamServer{
		{
			Server: server,
		},
	}
	return *core.NewServerUpdateEvent(eventType, upstreamName, emptyStreamServers, servers)
}
