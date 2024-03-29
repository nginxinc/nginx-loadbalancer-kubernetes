/*
 * Copyright (c) 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package core

// UpstreamServer represents a single upstream server.
// This is an internal representation used to abstract the definition
// of an upstream server from any specific client.
type UpstreamServer struct {

	// Host is the host name or IP address of the upstream server.
	Host string
}

// UpstreamServers is a slice of UpstreamServer.
type UpstreamServers = []*UpstreamServer

// NewUpstreamServer creates a new UpstreamServer.
func NewUpstreamServer(host string) *UpstreamServer {
	return &UpstreamServer{
		Host: host,
	}
}
