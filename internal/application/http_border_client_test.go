/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package application

import (
	"testing"
)

func TestHttpBorderClient_Delete(t *testing.T) {
	event := buildServerUpdateEvent(deletedEventType, ClientTypeHttp)
	borderClient, nginxClient, err := buildBorderClient(ClientTypeHttp)
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
	event := buildServerUpdateEvent(createEventType, ClientTypeHttp)
	borderClient, nginxClient, err := buildBorderClient(ClientTypeHttp)
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
	_, err := NewBorderClient(ClientTypeHttp, emptyInterface)
	if err == nil {
		t.Fatalf(`expected an error to occur when creating a new border client`)
	}
}

func TestHttpBorderClient_DeleteReturnsError(t *testing.T) {
	event := buildServerUpdateEvent(deletedEventType, ClientTypeHttp)
	borderClient, _, err := buildTerrorizingBorderClient(ClientTypeHttp)
	if err != nil {
		t.Fatalf(`error occurred creating a new border client: %v`, err)
	}

	err = borderClient.Delete(event)

	if err == nil {
		t.Fatalf(`expected an error to occur when deleting the nginx+ upstream server`)
	}
}

func TestHttpBorderClient_UpdateReturnsError(t *testing.T) {
	event := buildServerUpdateEvent(createEventType, ClientTypeHttp)
	borderClient, _, err := buildTerrorizingBorderClient(ClientTypeHttp)
	if err != nil {
		t.Fatalf(`error occurred creating a new border client: %v`, err)
	}

	err = borderClient.Update(event)

	if err == nil {
		t.Fatalf(`expected an error to occur when deleting the nginx+ upstream server`)
	}
}
