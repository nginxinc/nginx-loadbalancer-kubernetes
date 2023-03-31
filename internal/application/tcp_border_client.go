/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package application

import (
	"fmt"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/core"
)

type TcpBorderClient struct {
	BorderClient
	nginxClient NginxClientInterface
}

func NewTcpBorderClient(client interface{}) (Interface, error) {
	nginxClient, ok := client.(NginxClientInterface)
	if !ok {
		return nil, fmt.Errorf(`expected a NginxClientInterface, got a %v`, client)
	}

	return &TcpBorderClient{
		nginxClient: nginxClient,
	}, nil
}

func (tbc *TcpBorderClient) Update(event *core.ServerUpdateEvent) error {
	_, _, _, err := tbc.nginxClient.UpdateStreamServers(event.NginxHost, nil)
	if err != nil {
		return fmt.Errorf(`error occurred updating the nginx+ upstream server: %w`, err)
	}

	return nil
}

func (tbc *TcpBorderClient) Delete(event *core.ServerUpdateEvent) error {
	err := tbc.nginxClient.DeleteStreamServer(event.NginxHost, event.TcpServers[0].Server)
	if err != nil {
		return fmt.Errorf(`error occurred deleting the nginx+ upstream server: %w`, err)
	}

	return nil
}
