/*
 * Copyright (c) 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package mocks

import nginxClient "github.com/nginxinc/nginx-plus-go-client/client"

type MockNginxClient struct {
	CalledFunctions map[string]bool
}

func NewMockNginxClient() *MockNginxClient {
	return &MockNginxClient{
		CalledFunctions: make(map[string]bool),
	}
}

func (m MockNginxClient) DeleteStreamServer(upstream string, server string) error {
	m.CalledFunctions["DeleteStreamServer"] = true
	return nil
}

func (m MockNginxClient) UpdateStreamServers(upstream string, servers []nginxClient.StreamUpstreamServer) ([]nginxClient.StreamUpstreamServer, []nginxClient.StreamUpstreamServer, []nginxClient.StreamUpstreamServer, error) {
	m.CalledFunctions["UpdateStreamServers"] = true
	return nil, nil, nil, nil
}

func (m MockNginxClient) DeleteHTTPServer(upstream string, server string) error {
	m.CalledFunctions["DeleteHTTPServer"] = true
	return nil
}

func (m MockNginxClient) UpdateHTTPServers(upstream string, servers []nginxClient.UpstreamServer) ([]nginxClient.UpstreamServer, []nginxClient.UpstreamServer, []nginxClient.UpstreamServer, error) {
	m.CalledFunctions["UpdateHTTPServers"] = true
	return nil, nil, nil, nil
}
