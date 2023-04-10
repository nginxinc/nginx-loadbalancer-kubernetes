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

type TcpBorderClient struct {
	BorderClient
	nginxClient NginxClientInterface
}

func NewTcpBorderClient(client interface{}) (Interface, error) {
	ngxClient, ok := client.(NginxClientInterface)
	if !ok {
		return nil, fmt.Errorf(`expected a NginxClientInterface, got a %v`, client)
	}

	return &TcpBorderClient{
		nginxClient: ngxClient,
	}, nil
}

func (tbc *TcpBorderClient) Update(event *core.ServerUpdateEvent) error {
	streamUpstreamServers := asNginxStreamUpstreamServers(event.UpstreamServers)
	_, _, _, err := tbc.nginxClient.UpdateStreamServers(event.UpstreamName, streamUpstreamServers)
	if err != nil {
		return fmt.Errorf(`error occurred updating the nginx+ upstream server: %w`, err)
	}

	return nil
}

func (tbc *TcpBorderClient) Delete(event *core.ServerUpdateEvent) error {
	err := tbc.nginxClient.DeleteStreamServer(event.UpstreamName, event.UpstreamServers[0].Host)
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
	var upstreamServers []nginxClient.StreamUpstreamServer

	for _, server := range servers {
		upstreamServers = append(upstreamServers, asNginxStreamUpstreamServer(server))
	}

	return upstreamServers
}
