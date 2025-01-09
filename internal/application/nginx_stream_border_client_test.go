/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */
// dupl complains about duplicates with nginx_http_border_client_test.go
//nolint:dupl
package application

import (
	"context"
	"testing"
)

func TestTcpBorderClient_Delete(t *testing.T) {
	t.Parallel()
	event := buildServerUpdateEvent(deletedEventType, ClientTypeNginxStream)
	borderClient, nginxClient, err := buildBorderClient(ClientTypeNginxStream)
	if err != nil {
		t.Fatalf(`error occurred creating a new border client: %v`, err)
	}

	err = borderClient.Delete(context.Background(), event)
	if err != nil {
		t.Fatalf(`error occurred deleting the nginx+ upstream server: %v`, err)
	}

	if !nginxClient.CalledFunctions["DeleteStreamServer"] {
		t.Fatalf(`expected DeleteStreamServer to be called`)
	}
}

func TestTcpBorderClient_Update(t *testing.T) {
	t.Parallel()
	event := buildServerUpdateEvent(createEventType, ClientTypeNginxStream)
	borderClient, nginxClient, err := buildBorderClient(ClientTypeNginxStream)
	if err != nil {
		t.Fatalf(`error occurred creating a new border client: %v`, err)
	}

	err = borderClient.Update(context.Background(), event)
	if err != nil {
		t.Fatalf(`error occurred deleting the nginx+ upstream server: %v`, err)
	}

	if !nginxClient.CalledFunctions["UpdateStreamServers"] {
		t.Fatalf(`expected UpdateStreamServers to be called`)
	}
}

func TestTcpBorderClient_BadNginxClient(t *testing.T) {
	t.Parallel()
	var emptyInterface interface{}
	_, err := NewBorderClient(ClientTypeNginxStream, emptyInterface)
	if err == nil {
		t.Fatalf(`expected an error to occur when creating a new border client`)
	}
}

func TestTcpBorderClient_DeleteReturnsError(t *testing.T) {
	t.Parallel()
	event := buildServerUpdateEvent(deletedEventType, ClientTypeNginxStream)
	borderClient, err := buildTerrorizingBorderClient(ClientTypeNginxStream)
	if err != nil {
		t.Fatalf(`error occurred creating a new border client: %v`, err)
	}

	err = borderClient.Delete(context.Background(), event)

	if err == nil {
		t.Fatalf(`expected an error to occur when deleting the nginx+ upstream server`)
	}
}

func TestTcpBorderClient_UpdateReturnsError(t *testing.T) {
	t.Parallel()
	event := buildServerUpdateEvent(createEventType, ClientTypeNginxStream)
	borderClient, err := buildTerrorizingBorderClient(ClientTypeNginxStream)
	if err != nil {
		t.Fatalf(`error occurred creating a new border client: %v`, err)
	}

	err = borderClient.Update(context.Background(), event)

	if err == nil {
		t.Fatalf(`expected an error to occur when deleting the nginx+ upstream server`)
	}
}
