// Copyright 2023 f5 Inc. All rights reserved.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE file.

package translation

import (
	"fmt"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/core"
	nginxClient "github.com/nginxinc/nginx-plus-go-client/client"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

func Translate(event *core.Event) (*core.Event, error) {
	logrus.Debug("Translate::Translate")

	addresses, err := extractAddresses(event.Service)
	if err != nil {
		return event, fmt.Errorf(`error translating Servuce: %#v`, err)
	}

	buildAndAppendUpstreams(event, addresses)

	return event, nil
}

func buildAndAppendUpstreams(event *core.Event, addresses []string) {
	for _, address := range addresses {
		event.NginxUpstreams = append(event.NginxUpstreams, nginxClient.UpstreamServer{
			Server: address,
		})
	}
}

func extractAddresses(ingress *v1.Service) ([]string, error) {
	logrus.Infof("extractAddresses::ingress: %#v", ingress)
	var addresses []string

	//ingresses := ingress.Status.LoadBalancer.Ingress
	//
	//for _, ingress := range ingresses {
	//	if ingress.IP != "" {
	//		addresses = append(addresses, ingress.IP)
	//	} else if ingress.Hostname != "" {
	//		addresses = append(addresses, ingress.Hostname)
	//	} else {
	//		return nil, errors.New("ingress status does not contain IP or Hostname")
	//	}
	//}

	return addresses, nil
}
