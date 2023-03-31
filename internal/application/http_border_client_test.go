/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package application

import (
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

func TestHttpBorderClient_Delete(t *testing.T) {
	servers := []nginxClient2.StreamUpstreamServer{
		{
			Server: server,
		},
	}
	event := core.NewServerUpdateEvent(deletedEventType, upstreamName, servers)
	nginxClient := mocks.NewMockNginxClient()
	borderClient, err := NewBorderClient(clientType, nginxClient)
	if err != nil {
		t.Fatalf(`error occurred creating a new border client: %v`, err)
	}

	err = borderClient.Delete(*event)
	if err != nil {
		t.Fatalf(`error occurred deleting the nginx+ upstream server: %v`, err)
	}

	if !nginxClient.CalledFunctions["DeleteHTTPServer"] {
		t.Fatalf(`expected DeleteHTTPServer to be called`)
	}
}

func TestHttpBorderClient_Update(t *testing.T) {
	servers := []nginxClient2.StreamUpstreamServer{
		{
			Server: server,
		},
	}
	event := core.NewServerUpdateEvent(deletedEventType, upstreamName, servers)
	nginxClient := mocks.NewMockNginxClient()
	borderClient, err := NewBorderClient(clientType, nginxClient)
	if err != nil {
		t.Fatalf(`error occurred creating a new border client: %v`, err)
	}

	err = borderClient.Update(*event)
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
