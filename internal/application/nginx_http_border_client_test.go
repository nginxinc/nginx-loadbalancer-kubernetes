/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

// dupl complains about duplicates with nginx_stream_border_client_test.go
//
//nolint:dupl
package application

import (
	"context"
	"testing"
)

func TestHttpBorderClient_Delete(t *testing.T) {
	t.Parallel()
	event := buildServerUpdateEvent(deletedEventType, ClientTypeNginxHTTP)
	borderClient, nginxClient, err := buildBorderClient(ClientTypeNginxHTTP)
	if err != nil {
		t.Fatalf(`error occurred creating a new border client: %v`, err)
	}

	err = borderClient.Delete(context.Background(), event)
	if err != nil {
		t.Fatalf(`error occurred deleting the nginx+ upstream server: %v`, err)
	}

	if !nginxClient.CalledFunctions["DeleteHTTPServer"] {
		t.Fatalf(`expected DeleteHTTPServer to be called`)
	}
}

func TestHttpBorderClient_Update(t *testing.T) {
	t.Parallel()
	event := buildServerUpdateEvent(createEventType, ClientTypeNginxHTTP)
	borderClient, nginxClient, err := buildBorderClient(ClientTypeNginxHTTP)
	if err != nil {
		t.Fatalf(`error occurred creating a new border client: %v`, err)
	}

	err = borderClient.Update(context.Background(), event)
	if err != nil {
		t.Fatalf(`error occurred deleting the nginx+ upstream server: %v`, err)
	}

	if !nginxClient.CalledFunctions["UpdateHTTPServers"] {
		t.Fatalf(`expected UpdateHTTPServers to be called`)
	}
}

func TestHttpBorderClient_BadNginxClient(t *testing.T) {
	t.Parallel()
	var emptyInterface interface{}
	_, err := NewBorderClient(ClientTypeNginxHTTP, emptyInterface)
	if err == nil {
		t.Fatalf(`expected an error to occur when creating a new border client`)
	}
}

func TestHttpBorderClient_DeleteReturnsError(t *testing.T) {
	t.Parallel()
	event := buildServerUpdateEvent(deletedEventType, ClientTypeNginxHTTP)
	borderClient, err := buildTerrorizingBorderClient(ClientTypeNginxHTTP)
	if err != nil {
		t.Fatalf(`error occurred creating a new border client: %v`, err)
	}

	err = borderClient.Delete(context.Background(), event)

	if err == nil {
		t.Fatalf(`expected an error to occur when deleting the nginx+ upstream server`)
	}
}

func TestHttpBorderClient_UpdateReturnsError(t *testing.T) {
	t.Parallel()
	event := buildServerUpdateEvent(createEventType, ClientTypeNginxHTTP)
	borderClient, err := buildTerrorizingBorderClient(ClientTypeNginxHTTP)
	if err != nil {
		t.Fatalf(`error occurred creating a new border client: %v`, err)
	}

	err = borderClient.Update(context.Background(), event)

	if err == nil {
		t.Fatalf(`expected an error to occur when deleting the nginx+ upstream server`)
	}
}
