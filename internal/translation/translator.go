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
	"strings"
)

const NklPrefix = "nkl-"

func Translate(event *core.Event) (core.ServerUpdateEvents, error) {
	logrus.Debug("Translate::Translate")

	portsOfInterest := filterPorts(event.Service.Spec.Ports)

	return buildServerUpdateEvents(portsOfInterest, event.NodeIps)
}

func filterPorts(ports []v1.ServicePort) []v1.ServicePort {
	var portsOfInterest []v1.ServicePort

	for _, port := range ports {
		if strings.HasPrefix(port.Name, NklPrefix) {
			portsOfInterest = append(portsOfInterest, port)
		}
	}

	return portsOfInterest
}

// TODO: Get the list of Node IPs from the Kubernetes API and fan out over the port
func buildServerUpdateEvents(ports []v1.ServicePort, nodeIps []string) (core.ServerUpdateEvents, error) {
	logrus.Debugf("Translate::buildServerUpdateEvents(ports=%#v)", ports)

	upstreams := core.ServerUpdateEvents{}
	for _, port := range ports {
		ingressName := fixIngressName(port.Name)
		servers, _ := buildServers(nodeIps, port)

		upstreams = append(upstreams, core.NewServerUpdateEvent(ingressName, servers))
	}

	return upstreams, nil
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
