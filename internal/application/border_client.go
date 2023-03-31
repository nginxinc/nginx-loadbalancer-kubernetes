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
	Update(core.ServerUpdateEvent) error
	Delete(core.ServerUpdateEvent) error
}

type BorderClient struct {
}

func NewBorderClient(whichType string, borderClient interface{}) (Interface, error) {
	logrus.Debugf(`NewBorderClient for type: %s`, whichType)

	switch whichType {
	case "tcp":
		return &TcpBorderClient{
			BorderClient: BorderClient{},
		}, nil
	case "http":
		return NewHttpBorderClient(borderClient)

	default:
		return nil, fmt.Errorf(`unknown border client type: %s`, whichType)
	}
}
