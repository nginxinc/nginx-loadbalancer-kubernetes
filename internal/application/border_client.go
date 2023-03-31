/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package application

import (
	"fmt"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/core"
	"github.com/sirupsen/logrus"
)

type Interface interface {
	Update(*core.ServerUpdateEvent) error
	Delete(*core.ServerUpdateEvent) error
}

type BorderClient struct {
}

// NewBorderClient Returns a NullBorderClient if the type is unknown, this avoids panics due to nil pointer dereferences.
func NewBorderClient(clientType string, borderClient interface{}) (Interface, error) {
	logrus.Debugf(`NewBorderClient for type: %s`, clientType)

	switch clientType {
	case "tcp":
		return NewTcpBorderClient(borderClient)

	case "http":
		return NewHttpBorderClient(borderClient)

	default:
		borderClient, _ := NewNullBorderClient()
		return borderClient, fmt.Errorf(`unknown border client type: %s`, clientType)
	}
}
