/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */
// dupl complains about duplicates with nginx_stream_border_client.go
//nolint:dupl
package application

import (
	"context"
	"fmt"

	nginxClient "github.com/nginx/nginx-plus-go-client/v2/client"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/core"
)

// NginxHttpBorderClient implements the BorderClient interface for HTTP upstreams.
type NginxHTTPBorderClient struct {
	BorderClient
	nginxClient NginxClientInterface
}

// NewNginxHTTPBorderClient is the Factory function for creating an NewNginxHTTPBorderClient.
func NewNginxHTTPBorderClient(client interface{}) (Interface, error) {
	ngxClient, ok := client.(NginxClientInterface)
	if !ok {
		return nil, fmt.Errorf(`expected a NginxClientInterface, got a %v`, client)
	}

	return &NginxHTTPBorderClient{
		nginxClient: ngxClient,
	}, nil
}

// Update manages the Upstream servers for the Upstream Name given in the ServerUpdateEvent.
func (hbc *NginxHTTPBorderClient) Update(ctx context.Context, event *core.ServerUpdateEvent) error {
	httpUpstreamServers := asNginxHTTPUpstreamServers(event.UpstreamServers)
	_, _, _, err := hbc.nginxClient.UpdateHTTPServers(ctx, event.UpstreamName, httpUpstreamServers)
	if err != nil {
		return fmt.Errorf(`error occurred updating the nginx+ upstream server: %w`, err)
	}

	return nil
}

// Delete deletes the Upstream server for the Upstream Name given in the ServerUpdateEvent.
func (hbc *NginxHTTPBorderClient) Delete(ctx context.Context, event *core.ServerUpdateEvent) error {
	err := hbc.nginxClient.DeleteHTTPServer(ctx, event.UpstreamName, event.UpstreamServers[0].Host)
	if err != nil {
		return fmt.Errorf(`error occurred deleting the nginx+ upstream server: %w`, err)
	}

	return nil
}

// asNginxHttpUpstreamServer converts a core.UpstreamServer to a nginxClient.UpstreamServer.
func asNginxHTTPUpstreamServer(server *core.UpstreamServer) nginxClient.UpstreamServer {
	return nginxClient.UpstreamServer{
		Server: server.Host,
	}
}

// asNginxHTTPUpstreamServers converts a core.UpstreamServers to a []nginxClient.UpstreamServer.
func asNginxHTTPUpstreamServers(servers core.UpstreamServers) []nginxClient.UpstreamServer {
	upstreamServers := []nginxClient.UpstreamServer{}

	for _, server := range servers {
		upstreamServers = append(upstreamServers, asNginxHTTPUpstreamServer(server))
	}

	return upstreamServers
}
