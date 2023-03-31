/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package application

import (
	"testing"
)

func TestTcpBorderClient_Delete(t *testing.T) {
	event := buildServerUpdateEvent(deletedEventType)
	borderClient, nginxClient, err := buildBorderClient(clientTypeTcp)
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
	event := buildServerUpdateEvent(createEventType)
	borderClient, nginxClient, err := buildBorderClient(clientTypeTcp)
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

func TestTcpBorderClient_BadNginxClient(t *testing.T) {
	var emptyInterface interface{}
	_, err := NewBorderClient(clientTypeTcp, emptyInterface)
	if err == nil {
		t.Fatalf(`expected an error to occur when creating a new border client`)
	}
}

func TestTcpBorderClient_DeleteReturnsError(t *testing.T) {
	event := buildServerUpdateEvent(deletedEventType)
	borderClient, _, err := buildTerrorizingBorderClient(clientTypeTcp)
	if err != nil {
		t.Fatalf(`error occurred creating a new border client: %v`, err)
	}

	err = borderClient.Delete(event)

	if err == nil {
		t.Fatalf(`expected an error to occur when deleting the nginx+ upstream server`)
	}
}

func TestTcpBorderClient_UpdateReturnsError(t *testing.T) {
	event := buildServerUpdateEvent(createEventType)
	borderClient, _, err := buildTerrorizingBorderClient(clientTypeTcp)
	if err != nil {
		t.Fatalf(`error occurred creating a new border client: %v`, err)
	}

	err = borderClient.Update(event)

	if err == nil {
		t.Fatalf(`expected an error to occur when deleting the nginx+ upstream server`)
	}
}
