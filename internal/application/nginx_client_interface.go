/*
 * Copyright (c) 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package application

import nginxClient "github.com/nginxinc/nginx-plus-go-client/client"

type NginxClientInterface interface {
	DeleteStreamServer(upstream string, server string) error
	UpdateStreamServers(upstream string, servers []nginxClient.StreamUpstreamServer) ([]nginxClient.StreamUpstreamServer, []nginxClient.StreamUpstreamServer, []nginxClient.StreamUpstreamServer, error)

	DeleteHTTPServer(upstream string, server string) error
	UpdateHTTPServers(upstream string, servers []nginxClient.UpstreamServer) ([]nginxClient.UpstreamServer, []nginxClient.UpstreamServer, []nginxClient.UpstreamServer, error)
}
