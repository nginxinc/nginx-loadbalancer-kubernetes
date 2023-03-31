/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package application

import "github.com/nginxinc/kubernetes-nginx-ingress/internal/core"

type TcpBorderClient struct {
	BorderClient
}

func (client *TcpBorderClient) Update(_ core.ServerUpdateEvent) error {
	return nil
}

func (client *TcpBorderClient) Delete(_ core.ServerUpdateEvent) error {
	return nil
}
