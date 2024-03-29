/*
 * Copyright (c) 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package mocks

import nginxClient "github.com/nginxinc/nginx-plus-go-client/client"

type MockNginxClient struct {
	CalledFunctions map[string]bool
	Error           error
}

func NewMockNginxClient() *MockNginxClient {
	return &MockNginxClient{
		CalledFunctions: make(map[string]bool),
		Error:           nil,
	}
}

func NewErroringMockClient(err error) *MockNginxClient {
	return &MockNginxClient{
		CalledFunctions: make(map[string]bool),
		Error:           err,
	}
}

func (m MockNginxClient) DeleteStreamServer(_ string, _ string) error {
	m.CalledFunctions["DeleteStreamServer"] = true

	if m.Error != nil {
		return m.Error
	}

	return nil
}

func (m MockNginxClient) UpdateStreamServers(
	_ string,
	_ []nginxClient.StreamUpstreamServer,
) ([]nginxClient.StreamUpstreamServer, []nginxClient.StreamUpstreamServer, []nginxClient.StreamUpstreamServer, error) {
	m.CalledFunctions["UpdateStreamServers"] = true

	if m.Error != nil {
		return nil, nil, nil, m.Error
	}

	return nil, nil, nil, nil
}

func (m MockNginxClient) DeleteHTTPServer(_ string, _ string) error {
	m.CalledFunctions["DeleteHTTPServer"] = true

	if m.Error != nil {
		return m.Error
	}

	return nil
}

func (m MockNginxClient) UpdateHTTPServers(
	_ string,
	_ []nginxClient.UpstreamServer,
) ([]nginxClient.UpstreamServer, []nginxClient.UpstreamServer, []nginxClient.UpstreamServer, error) {
	m.CalledFunctions["UpdateHTTPServers"] = true

	if m.Error != nil {
		return nil, nil, nil, m.Error
	}

	return nil, nil, nil, nil
}
