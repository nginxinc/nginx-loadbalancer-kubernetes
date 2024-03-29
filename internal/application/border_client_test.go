/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package application

import (
	"testing"

	"github.com/nginxinc/kubernetes-nginx-ingress/test/mocks"
)

func TestBorderClient_CreatesHttpBorderClient(t *testing.T) {
	t.Parallel()
	borderClient := mocks.MockNginxClient{}
	client, err := NewBorderClient("http", borderClient)
	if err != nil {
		t.Errorf(`error creating border client: %v`, err)
	}

	if _, ok := client.(*NginxHTTPBorderClient); !ok {
		t.Errorf(`expected client to be of type NginxHTTPBorderClient`)
	}
}

func TestBorderClient_CreatesTcpBorderClient(t *testing.T) {
	t.Parallel()
	borderClient := mocks.MockNginxClient{}
	client, err := NewBorderClient("stream", borderClient)
	if err != nil {
		t.Errorf(`error creating border client: %v`, err)
	}

	if _, ok := client.(*NginxStreamBorderClient); !ok {
		t.Errorf(`expected client to be of type NginxStreamBorderClient`)
	}
}

func TestBorderClient_UnknownClientType(t *testing.T) {
	t.Parallel()
	unknownClientType := "unknown"
	borderClient := mocks.MockNginxClient{}
	client, err := NewBorderClient(unknownClientType, borderClient)
	if err == nil {
		t.Errorf(`expected error creating border client`)
	}

	if err.Error() != `unknown border client type: unknown` {
		t.Errorf(`expected error to be 'unknown border client type: unknown', got: %v`, err)
	}

	if _, ok := client.(*NullBorderClient); !ok {
		t.Errorf(`expected client to be of type NullBorderClient`)
	}
}
