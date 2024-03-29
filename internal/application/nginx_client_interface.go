/*
 * Copyright (c) 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package application

import nginxClient "github.com/nginxinc/nginx-plus-go-client/client"

// NginxClientInterface defines the functions used on the NGINX Plus client,
// abstracting away the full details of that client.
type NginxClientInterface interface {
	// DeleteStreamServer is used by the NginxStreamBorderClient.
	DeleteStreamServer(upstream string, server string) error

	// UpdateStreamServers is used by the NginxStreamBorderClient.
	UpdateStreamServers(
		upstream string,
		servers []nginxClient.StreamUpstreamServer,
	) ([]nginxClient.StreamUpstreamServer, []nginxClient.StreamUpstreamServer, []nginxClient.StreamUpstreamServer, error)

	// DeleteHTTPServer is used by the NginxHTTPBorderClient.
	DeleteHTTPServer(upstream string, server string) error

	// UpdateHTTPServers is used by the NginxHTTPBorderClient.
	UpdateHTTPServers(
		upstream string,
		servers []nginxClient.UpstreamServer,
	) ([]nginxClient.UpstreamServer, []nginxClient.UpstreamServer, []nginxClient.UpstreamServer, error)
}
