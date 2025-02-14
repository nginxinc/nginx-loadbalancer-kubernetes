/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */
// dupl complains about duplicates with nginx_http_border_client.go
//nolint:dupl
package application

import (
	"context"
	"fmt"

	nginxClient "github.com/nginx/nginx-plus-go-client/v2/client"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/core"
)

// NginxStreamBorderClient implements the BorderClient interface for stream upstreams.
type NginxStreamBorderClient struct {
	BorderClient
	nginxClient NginxClientInterface
	ctx         context.Context
}

// NewNginxStreamBorderClient is the Factory function for creating an NginxStreamBorderClient.
func NewNginxStreamBorderClient(client interface{}) (Interface, error) {
	ngxClient, ok := client.(NginxClientInterface)
	if !ok {
		return nil, fmt.Errorf(`expected a NginxClientInterface, got a %v`, client)
	}

	return &NginxStreamBorderClient{
		nginxClient: ngxClient,
		ctx:         context.Background(),
	}, nil
}

// Update manages the Upstream servers for the Upstream Name given in the ServerUpdateEvent.
func (tbc *NginxStreamBorderClient) Update(ctx context.Context, event *core.ServerUpdateEvent) error {
	streamUpstreamServers := asNginxStreamUpstreamServers(event.UpstreamServers)
	_, _, _, err := tbc.nginxClient.UpdateStreamServers(ctx, event.UpstreamName, streamUpstreamServers)
	if err != nil {
		return fmt.Errorf(`error occurred updating the nginx+ upstream server: %w`, err)
	}

	return nil
}

// Delete deletes the Upstream server for the Upstream Name given in the ServerUpdateEvent.
func (tbc *NginxStreamBorderClient) Delete(ctx context.Context, event *core.ServerUpdateEvent) error {
	err := tbc.nginxClient.DeleteStreamServer(ctx, event.UpstreamName, event.UpstreamServers[0].Host)
	if err != nil {
		return fmt.Errorf(`error occurred deleting the nginx+ upstream server: %w`, err)
	}

	return nil
}

func asNginxStreamUpstreamServer(server *core.UpstreamServer) nginxClient.StreamUpstreamServer {
	return nginxClient.StreamUpstreamServer{
		Server: server.Host,
	}
}

func asNginxStreamUpstreamServers(servers core.UpstreamServers) []nginxClient.StreamUpstreamServer {
	upstreamServers := []nginxClient.StreamUpstreamServer{}

	for _, server := range servers {
		upstreamServers = append(upstreamServers, asNginxStreamUpstreamServer(server))
	}

	return upstreamServers
}
