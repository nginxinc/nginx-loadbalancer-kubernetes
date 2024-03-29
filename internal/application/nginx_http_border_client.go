/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package application

import (
	"fmt"

	"github.com/nginxinc/kubernetes-nginx-ingress/internal/core"
	nginxClient "github.com/nginxinc/nginx-plus-go-client/client"
)

// NginxHttpBorderClient implements the BorderClient interface for HTTP upstreams.
type NginxHttpBorderClient struct {
	BorderClient
	nginxClient NginxClientInterface
}

// NewNginxHttpBorderClient is the Factory function for creating an NginxHttpBorderClient.
func NewNginxHttpBorderClient(client interface{}) (Interface, error) {
	ngxClient, ok := client.(NginxClientInterface)
	if !ok {
		return nil, fmt.Errorf(`expected a NginxClientInterface, got a %v`, client)
	}

	return &NginxHttpBorderClient{
		nginxClient: ngxClient,
	}, nil
}

// Update manages the Upstream servers for the Upstream Name given in the ServerUpdateEvent.
func (hbc *NginxHttpBorderClient) Update(event *core.ServerUpdateEvent) error {
	httpUpstreamServers := asNginxHttpUpstreamServers(event.UpstreamServers)
	_, _, _, err := hbc.nginxClient.UpdateHTTPServers(event.UpstreamName, httpUpstreamServers)
	if err != nil {
		return fmt.Errorf(`error occurred updating the nginx+ upstream server: %w`, err)
	}

	return nil
}

// Delete deletes the Upstream server for the Upstream Name given in the ServerUpdateEvent.
func (hbc *NginxHttpBorderClient) Delete(event *core.ServerUpdateEvent) error {
	err := hbc.nginxClient.DeleteHTTPServer(event.UpstreamName, event.UpstreamServers[0].Host)
	if err != nil {
		return fmt.Errorf(`error occurred deleting the nginx+ upstream server: %w`, err)
	}

	return nil
}

// asNginxHttpUpstreamServer converts a core.UpstreamServer to a nginxClient.UpstreamServer.
func asNginxHttpUpstreamServer(server *core.UpstreamServer) nginxClient.UpstreamServer {
	return nginxClient.UpstreamServer{
		Server: server.Host,
	}
}

// asNginxHttpUpstreamServers converts a core.UpstreamServers to a []nginxClient.UpstreamServer.
func asNginxHttpUpstreamServers(servers core.UpstreamServers) []nginxClient.UpstreamServer {
	var upstreamServers []nginxClient.UpstreamServer

	for _, server := range servers {
		upstreamServers = append(upstreamServers, asNginxHttpUpstreamServer(server))
	}

	return upstreamServers
}
