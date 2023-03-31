/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package application

import (
	"github.com/nginxinc/kubernetes-nginx-ingress/test/mocks"
	"testing"
)

func TestBorderClient_CreatesHttpBorderClient(t *testing.T) {
	borderClient := mocks.MockNginxClient{}
	client, err := NewBorderClient("http", borderClient)
	if err != nil {
		t.Errorf(`error creating border client: %v`, err)
	}

	if _, ok := client.(*HttpBorderClient); !ok {
		t.Errorf(`expected client to be of type HttpBorderClient`)
	}
}

func TestBorderClient_CreatesTcpBorderClient(t *testing.T) {
	borderClient := mocks.MockNginxClient{}
	client, err := NewBorderClient("tcp", borderClient)
	if err != nil {
		t.Errorf(`error creating border client: %v`, err)
	}

	if _, ok := client.(*TcpBorderClient); !ok {
		t.Errorf(`expected client to be of type TcpBorderClient`)
	}
}

func TestBorderClient_UnknownClientType(t *testing.T) {
	unknownClientType := "unknown"
	borderClient := mocks.MockNginxClient{}
	_, err := NewBorderClient(unknownClientType, borderClient)
	if err == nil {
		t.Errorf(`expected error creating border client`)
	}

	if err.Error() != `unknown border client type: unknown` {
		t.Errorf(`expected error to be 'unknown border client type: unknown', got: %v`, err)
	}
}
