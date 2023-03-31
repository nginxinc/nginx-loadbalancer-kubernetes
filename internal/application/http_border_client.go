/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package application

import (
	"fmt"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/core"
)

type HttpBorderClient struct {
	BorderClient
	nginxClient NginxClientInterface
}

func NewHttpBorderClient(client interface{}) (Interface, error) {
	nginxClient, ok := client.(NginxClientInterface)
	if !ok {
		return nil, fmt.Errorf(`expected a NginxClientInterface, got a %v`, client)
	}

	return &HttpBorderClient{
		nginxClient: nginxClient,
	}, nil
}

func (hbc *HttpBorderClient) Update(event core.ServerUpdateEvent) error {
	_, _, _, err := hbc.nginxClient.UpdateHTTPServers(event.NginxHost, nil)
	if err != nil {
		return fmt.Errorf(`error occurred updating the nginx+ upstream server: %w`, err)
	}

	return nil
}

func (hbc *HttpBorderClient) Delete(event core.ServerUpdateEvent) error {
	err := hbc.nginxClient.DeleteHTTPServer(event.NginxHost, event.HttpServers[0].Server)
	if err != nil {
		return fmt.Errorf(`error occurred deleting the nginx+ upstream server: %w`, err)
	}

	return nil
}
