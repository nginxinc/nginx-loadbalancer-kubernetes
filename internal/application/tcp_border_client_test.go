/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package application

import (
	"testing"
)

func TestTcpBorderClient_Delete(t *testing.T) {
	event := buildServerUpdateEvent(deletedEventType, ClientTypeTcp)
	borderClient, nginxClient, err := buildBorderClient(ClientTypeTcp)
	if err != nil {
		t.Fatalf(`error occurred creating a new border client: %v`, err)
	}

	err = borderClient.Delete(event)
	if err != nil {
		t.Fatalf(`error occurred deleting the nginx+ upstream server: %v`, err)
	}

	if !nginxClient.CalledFunctions["DeleteStreamServer"] {
		t.Fatalf(`expected DeleteStreamServer to be called`)
	}
}

func TestTcpBorderClient_Update(t *testing.T) {
	event := buildServerUpdateEvent(createEventType, ClientTypeTcp)
	borderClient, nginxClient, err := buildBorderClient(ClientTypeTcp)
	if err != nil {
		t.Fatalf(`error occurred creating a new border client: %v`, err)
	}

	err = borderClient.Update(event)
	if err != nil {
		t.Fatalf(`error occurred deleting the nginx+ upstream server: %v`, err)
	}

	if !nginxClient.CalledFunctions["UpdateStreamServers"] {
		t.Fatalf(`expected UpdateStreamServers to be called`)
	}
}

func TestTcpBorderClient_BadNginxClient(t *testing.T) {
	var emptyInterface interface{}
	_, err := NewBorderClient(ClientTypeTcp, emptyInterface)
	if err == nil {
		t.Fatalf(`expected an error to occur when creating a new border client`)
	}
}

func TestTcpBorderClient_DeleteReturnsError(t *testing.T) {
	event := buildServerUpdateEvent(deletedEventType, ClientTypeTcp)
	borderClient, _, err := buildTerrorizingBorderClient(ClientTypeTcp)
	if err != nil {
		t.Fatalf(`error occurred creating a new border client: %v`, err)
	}

	err = borderClient.Delete(event)

	if err == nil {
		t.Fatalf(`expected an error to occur when deleting the nginx+ upstream server`)
	}
}

func TestTcpBorderClient_UpdateReturnsError(t *testing.T) {
	event := buildServerUpdateEvent(createEventType, ClientTypeTcp)
	borderClient, _, err := buildTerrorizingBorderClient(ClientTypeTcp)
	if err != nil {
		t.Fatalf(`error occurred creating a new border client: %v`, err)
	}

	err = borderClient.Update(event)

	if err == nil {
		t.Fatalf(`expected an error to occur when deleting the nginx+ upstream server`)
	}
}
