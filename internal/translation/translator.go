/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package translation

import (
	"fmt"
	"strings"

	"github.com/nginxinc/kubernetes-nginx-ingress/internal/application"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/configuration"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/core"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

// Translate transforms event data into an intermediate format that can be consumed by the BorderClient implementations
// and used to update the Border Servers.
func Translate(event *core.Event) (core.ServerUpdateEvents, error) {
	logrus.Debug("Translate::Translate")

	portsOfInterest := filterPorts(event.Service.Spec.Ports)

	return buildServerUpdateEvents(portsOfInterest, event)
}

// filterPorts returns a list of ports that have the NlkPrefix in the port name.
func filterPorts(ports []v1.ServicePort) []v1.ServicePort {
	var portsOfInterest []v1.ServicePort

	for _, port := range ports {
		if strings.HasPrefix(port.Name, configuration.NlkPrefix) {
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
		upstreamServers, _ := buildUpstreamServers(event.NodeIps, port)
		clientType := getClientType(port.Name, event.Service.Annotations)

		switch event.Type {
		case core.Created:
			fallthrough

		case core.Updated:
			events = append(events, core.NewServerUpdateEvent(event.Type, ingressName, clientType, upstreamServers))

		case core.Deleted:
			for _, server := range upstreamServers {
				events = append(events, core.NewServerUpdateEvent(event.Type, ingressName, clientType, core.UpstreamServers{server}))
			}

		default:
			logrus.Warnf(`Translator::buildServerUpdateEvents: unknown event type: %d`, event.Type)
		}

	}

	return events, nil
}

func buildUpstreamServers(nodeIps []string, port v1.ServicePort) (core.UpstreamServers, error) {
	var servers core.UpstreamServers

	for _, nodeIp := range nodeIps {
		host := fmt.Sprintf("%s:%d", nodeIp, port.NodePort)
		server := core.NewUpstreamServer(host)
		servers = append(servers, server)
	}

	return servers, nil
}

// fixIngressName removes the NlkPrefix from the port name
func fixIngressName(name string) string {
	return name[4:]
}

// getClientType returns the client type for the port, defaults to ClientTypeNginxHttp if no Annotation is found.
func getClientType(portName string, annotations map[string]string) string {
	key := fmt.Sprintf("%s/%s", configuration.PortAnnotationPrefix, portName)
	logrus.Infof("getClientType: key=%s", key)
	if annotations != nil {
		if clientType, ok := annotations[key]; ok {
			return clientType
		}
	}

	return application.ClientTypeNginxHttp
}
