/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package application

import "testing"

func TestBorderClient_CreatesHttpBorderClient(t *testing.T) {
	client, err := NewBorderClient("http", nil)
	if err != nil {
		t.Errorf(`error creating border client: %v`, err)
	}

	if _, ok := client.(*HttpBorderClient); !ok {
		t.Errorf(`expected client to be of type HttpBorderClient`)
	}
}

func TestBorderClient_CreatesTcpBorderClient(t *testing.T) {
	client, err := NewBorderClient("tcp", nil)
	if err != nil {
		t.Errorf(`error creating border client: %v`, err)
	}

	if _, ok := client.(*TcpBorderClient); !ok {
		t.Errorf(`expected client to be of type TcpBorderClient`)
	}
}

func TestBorderClient_UnknownClientType(t *testing.T) {
	_, err := NewBorderClient("unknown", nil)
	if err == nil {
		t.Errorf(`expected error creating border client`)
	}

	if err.Error() != `unknown border client type: unknown` {
		t.Errorf(`expected error to be 'unknown border client type: unknown', got: %v`, err)
	}
}
