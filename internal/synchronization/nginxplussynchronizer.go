// Copyright 2023 f5 Inc. All rights reserved.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE file.

package synchronization

import (
	nginxClient "github.com/nginxinc/nginx-plus-go-client/client"
)

type NginxPlusSynchronizer struct {
	NginxPlusClient *nginxClient.NginxClient
}

func NewNginxPlusSynchronizer(nginxClient *nginxClient.NginxClient) (*NginxPlusSynchronizer, error) {
	synchronizer := NginxPlusSynchronizer{
		NginxPlusClient: nginxClient,
	}

	return &synchronizer, nil
}

func (*NginxPlusSynchronizer) Synchronize() (interface{}, error) {
	return "", nil
}
