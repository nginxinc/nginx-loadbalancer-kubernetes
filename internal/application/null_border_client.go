/*
 * Copyright (c) 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package application

import "github.com/nginxinc/kubernetes-nginx-ingress/internal/core"

type NullBorderClient struct {
}

func NewNullBorderClient() (Interface, error) {
	return &NullBorderClient{}, nil
}

func (nbc *NullBorderClient) Update(_ *core.ServerUpdateEvent) error {
	return nil
}

func (nbc *NullBorderClient) Delete(_ *core.ServerUpdateEvent) error {
	return nil
}
