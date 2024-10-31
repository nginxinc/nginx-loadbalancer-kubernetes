/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package translation

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/nginxinc/kubernetes-nginx-ingress/internal/core"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Translator struct {
	k8sClient kubernetes.Interface
}

func NewTranslator(k8sClient kubernetes.Interface) *Translator {
	return &Translator{k8sClient}
}

// Translate transforms event data into an intermediate format that can be consumed by the BorderClient implementations
// and used to update the Border Servers.
func (t *Translator) Translate(ctx context.Context, event *core.Event) (core.ServerUpdateEvents, error) {
	slog.Debug("Translate::Translate")

	return t.buildServerUpdateEvents(ctx, event.Service.Spec.Ports, event)
}

// buildServerUpdateEvents builds a list of ServerUpdateEvents based on the event type
// The NGINX+ Client uses a list of servers for Created and Updated events.
// The client performs reconciliation between the list of servers in the NGINX+ Client call
// and the list of servers in NGINX+.
// The NGINX+ Client uses a single server for Deleted events;
// so the list of servers is broken up into individual events.
func (t *Translator) buildServerUpdateEvents(ctx context.Context, ports []v1.ServicePort, event *core.Event,
) (events core.ServerUpdateEvents, err error) {
	slog.Debug("Translate::buildServerUpdateEvents", "ports", ports)

	switch event.Service.Spec.Type {
	case v1.ServiceTypeNodePort:
		return t.buildNodeIPEvents(ctx, ports, event)
	case v1.ServiceTypeClusterIP:
		return t.buildClusterIPEvents(ctx, event)
	default:
		return events, fmt.Errorf("unsupported service type: %s", event.Service.Spec.Type)
	}
}

type upstream struct {
	context string
	name    string
}

func (t *Translator) buildClusterIPEvents(ctx context.Context, event *core.Event,
) (events core.ServerUpdateEvents, err error) {
	namespace := event.Service.GetObjectMeta().GetNamespace()
	serviceName := event.Service.Name

	logger := slog.With("namespace", namespace, "serviceName", serviceName)
	logger.Debug("Translate::buildClusterIPEvents")

	if event.Type == core.Deleted {
		for _, port := range event.Service.Spec.Ports {
			context, upstreamName, pErr := getContextAndUpstreamName(port.Name)
			if pErr != nil {
				logger.Info(pErr.Error())
				continue
			}
			events = append(events, core.NewServerUpdateEvent(core.Updated, upstreamName, context, nil))
		}
		return events, nil
	}

	s := t.k8sClient.DiscoveryV1().EndpointSlices(namespace)
	list, err := s.List(ctx, metav1.ListOptions{LabelSelector: fmt.Sprintf("kubernetes.io/service-name=%s", serviceName)})
	if err != nil {
		logger.Error(`error occurred retrieving the list of endpoint slices`, "error", err)
		return events, err
	}

	upstreams := make(map[upstream][]*core.UpstreamServer)

	for _, endpointSlice := range list.Items {
		for _, port := range endpointSlice.Ports {
			if port.Name == nil || port.Port == nil {
				continue
			}

			context, upstreamName, err := getContextAndUpstreamName(*port.Name)
			if err != nil {
				logger.Info(err.Error())
				continue
			}

			u := upstream{
				context: context,
				name:    upstreamName,
			}
			servers := upstreams[u]

			for _, endpoint := range endpointSlice.Endpoints {
				for _, address := range endpoint.Addresses {
					host := fmt.Sprintf("%s:%d", address, *port.Port)
					servers = append(servers, core.NewUpstreamServer(host))
				}
			}

			upstreams[u] = servers
		}
	}

	for u, servers := range upstreams {
		events = append(events, core.NewServerUpdateEvent(core.Updated, u.name, u.context, servers))
	}

	return events, nil
}

func (t *Translator) buildNodeIPEvents(ctx context.Context, ports []v1.ServicePort, event *core.Event,
) (core.ServerUpdateEvents, error) {
	slog.Debug("Translate::buildNodeIPEvents", "ports", ports)

	events := core.ServerUpdateEvents{}
	for _, port := range ports {
		context, upstreamName, err := getContextAndUpstreamName(port.Name)
		if err != nil {
			slog.Info(err.Error())
			continue
		}

		addresses, err := t.retrieveNodeIps(ctx)
		if err != nil {
			return nil, err
		}

		upstreamServers := buildUpstreamServers(addresses, port)

		switch event.Type {
		case core.Created:
			fallthrough

		case core.Updated:
			events = append(events, core.NewServerUpdateEvent(event.Type, upstreamName, context, upstreamServers))

		case core.Deleted:
			events = append(events, core.NewServerUpdateEvent(
				core.Updated, upstreamName, context, nil,
			))
		default:
			slog.Warn(`Translator::buildNodeIPEvents: unknown event type`, "type", event.Type)
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
func getContextAndUpstreamName(portName string) (clientType string, appName string, err error) {
	context, upstreamName, found := strings.Cut(portName, "-")
	switch {
	case !found:
		return clientType, appName,
			fmt.Errorf("ignoring port %s because it is not in the format [http|stream]-{upstreamName}", portName)
	case context != "http" && context != "stream":
		return clientType, appName, fmt.Errorf("port name %s does not include \"http\" or \"stream\" context", portName)
	default:
		return context, upstreamName, nil
	}
}

// notMasterNode retrieves the IP Addresses of the nodes in the cluster. Currently, the master node is excluded. This is
// because the master node may or may not be a worker node and thus may not be able to route traffic.
func (t *Translator) retrieveNodeIps(ctx context.Context) ([]string, error) {
	started := time.Now()
	slog.Debug("Translator::retrieveNodeIps")

	var nodeIps []string

	nodes, err := t.k8sClient.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		slog.Error("error occurred retrieving the list of nodes", "error", err)
		return nil, err
	}

	for _, node := range nodes.Items {
		// this is kind of a broad assumption, should probably make this a configurable option
		if notMasterNode(node) {
			for _, address := range node.Status.Addresses {
				if address.Type == v1.NodeInternalIP {
					nodeIps = append(nodeIps, address.Address)
				}
			}
		}
	}

	slog.Debug("Translator::retrieveNodeIps duration", "duration", time.Since(started).Nanoseconds())

	return nodeIps, nil
}

// notMasterNode determines if the node is a master node.
func notMasterNode(node v1.Node) bool {
	slog.Debug("Translator::notMasterNode")

	_, found := node.Labels["node-role.kubernetes.io/master"]

	return !found
}
