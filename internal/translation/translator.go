/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package translation

import (
	"fmt"
	"strings"

	"github.com/nginxinc/kubernetes-nginx-ingress/internal/core"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

// Translate transforms event data into an intermediate format that can be consumed by the BorderClient implementations
// and used to update the Border Servers.
func Translate(event *core.Event) (core.ServerUpdateEvents, error) {
	logrus.Debug("Translate::Translate")

	return buildServerUpdateEvents(event.Service.Spec.Ports, event)
}

// buildServerUpdateEvents builds a list of ServerUpdateEvents based on the event type
// The NGINX+ Client uses a list of servers for Created and Updated events.
// The client performs reconciliation between the list of servers in the NGINX+ Client call
// and the list of servers in NGINX+.
// The NGINX+ Client uses a single server for Deleted events;
// so the list of servers is broken up into individual events.
func buildServerUpdateEvents(ports []v1.ServicePort, event *core.Event) (core.ServerUpdateEvents, error) {
	logrus.Debugf("Translate::buildServerUpdateEvents(ports=%#v)", ports)

	events := core.ServerUpdateEvents{}
	for _, port := range ports {
		context, upstreamName, err := getContextAndUpstreamName(port)
		if err != nil {
			logrus.Info(err)
			continue
		}

		upstreamServers := buildUpstreamServers(event.NodeIps, port)

		switch event.Type {
		case core.Created:
			fallthrough

		case core.Updated:
			events = append(events, core.NewServerUpdateEvent(event.Type, upstreamName, context, upstreamServers))

		case core.Deleted:
			for _, server := range upstreamServers {
				events = append(events, core.NewServerUpdateEvent(
					event.Type, upstreamName, context, core.UpstreamServers{server},
				))
			}

		default:
			logrus.Warnf(`Translator::buildServerUpdateEvents: unknown event type: %d`, event.Type)
		}

	}

	return events, nil
}

func buildUpstreamServers(nodeIPs []string, port v1.ServicePort) core.UpstreamServers {
	var servers core.UpstreamServers

	for _, nodeIP := range nodeIPs {
		host := fmt.Sprintf("%s:%d", nodeIP, port.NodePort)
		server := core.NewUpstreamServer(host)
		servers = append(servers, server)
	}

	return servers
}

// getContextAndUpstreamName returns the nginx context being supplied by the port (either "http" or "stream")
// and the upstream name.
func getContextAndUpstreamName(port v1.ServicePort) (clientType string, appName string, err error) {
	context, upstreamName, found := strings.Cut(port.Name, "-")
	switch {
	case !found:
		return clientType, appName,
			fmt.Errorf("ignoring port %s because it is not in the format [http|stream]-{upstreamName}", port.Name)
	case context != "http" && context != "stream":
		return clientType, appName, fmt.Errorf("port name %s does not include \"http\" or \"stream\" context", port.Name)
	default:
		return context, upstreamName, nil
	}
}
