/*
 * Copyright (c) 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package core

type UpstreamServer struct {
	Host string
}

type UpstreamServers = []*UpstreamServer

func NewUpstreamServer(host string) *UpstreamServer {
	return &UpstreamServer{
		Host: host,
	}
}
