// Copyright 2023 f5 Inc. All rights reserved.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE file.

package translation

import (
	"fmt"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/configuration"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/core"
	nginxClient "github.com/nginxinc/nginx-plus-go-client/client"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"strings"
)

func Translate(event *core.Event) (core.ServerUpdateEvents, error) {
	logrus.Debug("Translate::Translate")

	portsOfInterest := filterPorts(event.Service.Spec.Ports)

	return buildServerUpdateEvents(portsOfInterest, event)
}

func filterPorts(ports []v1.ServicePort) []v1.ServicePort {
	var portsOfInterest []v1.ServicePort

	for _, port := range ports {
		if strings.HasPrefix(port.Name, configuration.NklPrefix) {
			portsOfInterest = append(portsOfInterest, port)
		}
	}

	return portsOfInterest
}

// buildServerUpdateEvents builds a list of ServerUpdateEvents based on the event type
// The NGINX+ Client uses a list of servers for Created and Updated events; the client performs reconciliation between
// the list of servers in the NGINX+ Client call and the list of servers in NGINX+.
// The NGINX+ Client uses a single server for Deleted events; so the list of servers is broken up into individual events.
func buildServerUpdateEvents(ports []v1.ServicePort, event *core.Event) (core.ServerUpdateEvents, error) {
	logrus.Debugf("Translate::buildServerUpdateEvents(ports=%#v)", ports)

	events := core.ServerUpdateEvents{}
	for _, port := range ports {
		ingressName := fixIngressName(port.Name)
		servers, _ := buildServers(event.NodeIps, port)

		switch event.Type {
		case core.Created:
			fallthrough
		case core.Updated:
			events = append(events, core.NewServerUpdateEvent(event.Type, ingressName, servers))
		case core.Deleted:
			for _, server := range servers {
				events = append(events, core.NewServerUpdateEvent(event.Type, ingressName, []nginxClient.StreamUpstreamServer{server}))
			}
		default:
			logrus.Warnf(`Translator::buildServerUpdateEvents: unknown event type: %d`, event.Type)
		}

	}

	return events, nil
}

func buildServers(nodeIps []string, port v1.ServicePort) ([]nginxClient.StreamUpstreamServer, error) {
	var servers []nginxClient.StreamUpstreamServer

	for _, nodeIp := range nodeIps {
		server := nginxClient.StreamUpstreamServer{
			Server: fmt.Sprintf("%s:%d", nodeIp, port.NodePort),
		}
		servers = append(servers, server)
	}

	return servers, nil
}

func fixIngressName(name string) string {
	return name[4:]
}
