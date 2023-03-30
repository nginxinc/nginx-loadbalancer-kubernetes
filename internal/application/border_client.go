/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package application

import (
	"fmt"
	nginxClient "github.com/nginxinc/nginx-plus-go-client/client"
	"github.com/sirupsen/logrus"
)

type Interface interface {
	Update()
	Delete()
}

type BorderClient struct {
	NginxPlusClient *nginxClient.NginxClient
}

func NewBorderClient(whichType string, nginxClient *nginxClient.NginxClient) (Interface, error) {
	logrus.Debugf(`NewBorderClient for type: %s`, whichType)

	switch whichType {
	case "tcp":
		return &TcpBorderClient{
			BorderClient: BorderClient{
				NginxPlusClient: nginxClient,
			},
		}, nil
	case "http":
		return &HttpBorderClient{
			BorderClient: BorderClient{
				NginxPlusClient: nginxClient,
			},
		}, nil
	default:
		return nil, fmt.Errorf(`unknown border client type: %s`, whichType)
	}
}
