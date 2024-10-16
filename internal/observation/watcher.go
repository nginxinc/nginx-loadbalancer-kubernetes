/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package observation

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/nginxinc/kubernetes-nginx-ingress/internal/configuration"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/core"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

// Watcher is responsible for watching for changes to Kubernetes resources.
// Particularly, Services in the namespace defined in the WatcherSettings::NginxIngressNamespace setting.
// When a change is detected, an Event is generated and added to the Handler's queue.
type Watcher struct {
	// eventHandlerRegistration is used to track the event handlers
	eventHandlerRegistration interface{}

	// handler is the event handler
	handler HandlerInterface

	// informer is the informer used to watch for changes to Kubernetes resources
	informer cache.SharedIndexInformer

	k8sClient kubernetes.Interface

	// settings is the configuration settings
	settings configuration.Settings
}

// NewWatcher creates a new Watcher
func NewWatcher(
	settings configuration.Settings, handler HandlerInterface, k8sClient kubernetes.Interface,
) (*Watcher, error) {
	return &Watcher{
		handler:   handler,
		settings:  settings,
		k8sClient: k8sClient,
	}, nil
}

// Initialize initializes the Watcher, must be called before Watch
func (w *Watcher) Initialize(ctx context.Context) error {
	slog.Debug("Watcher::Initialize")
	var err error

	w.informer = w.buildInformer()

	err = w.initializeEventListeners(ctx)
	if err != nil {
		return fmt.Errorf(`initialization error: %w`, err)
	}

	return nil
}

// Watch starts the process of watching for changes to Kubernetes resources.
// Initialize must be called before Watch.
func (w *Watcher) Watch(ctx context.Context) error {
	slog.Debug("Watcher::Watch")

	if w.informer == nil {
		return errors.New("error: Initialize must be called before Watch")
	}

	defer utilruntime.HandleCrash()
	defer w.handler.ShutDown()

	go w.informer.Run(ctx.Done())

	if !cache.WaitForNamedCacheSync(
		w.settings.Handler.WorkQueueSettings.Name,
		ctx.Done(),
		w.informer.HasSynced,
	) {
		return fmt.Errorf(`error occurred waiting for the cache to sync`)
	}

	<-ctx.Done()
	return nil
}

// isDesiredService returns whether the user has configured the given service for watching.
func (w *Watcher) isDesiredService(service *v1.Service) bool {
	annotation, ok := service.Annotations["nginx.com/nginxaas"]
	if !ok {
		return false
	}

	return annotation == w.settings.Watcher.ServiceAnnotation
}

// buildEventHandlerForAdd creates a function that is used as an event handler
// for the informer when Add events are raised.
func (w *Watcher) buildEventHandlerForAdd(ctx context.Context) func(interface{}) {
	slog.Info("Watcher::buildEventHandlerForAdd")
	return func(obj interface{}) {
		service := obj.(*v1.Service)
		if !w.isDesiredService(service) {
			return
		}

		nodeIps, err := w.retrieveNodeIps(ctx)
		if err != nil {
			slog.Error("error occurred retrieving node ips", "error", err)
			return
		}

		var previousService *v1.Service
		e := core.NewEvent(core.Created, service, previousService, nodeIps)
		w.handler.AddRateLimitedEvent(&e)
	}
}

// buildEventHandlerForDelete creates a function that is used as an event handler
// for the informer when Delete events are raised.
func (w *Watcher) buildEventHandlerForDelete(ctx context.Context) func(interface{}) {
	slog.Info("Watcher::buildEventHandlerForDelete")
	return func(obj interface{}) {
		service := obj.(*v1.Service)
		if !w.isDesiredService(service) {
			return
		}

		nodeIps, err := w.retrieveNodeIps(ctx)
		if err != nil {
			slog.Error("error occurred retrieving node ips", "error", err)
			return
		}

		var previousService *v1.Service
		e := core.NewEvent(core.Deleted, service, previousService, nodeIps)
		w.handler.AddRateLimitedEvent(&e)
	}
}

// buildEventHandlerForUpdate creates a function that is used as an event handler
// for the informer when Update events are raised.
func (w *Watcher) buildEventHandlerForUpdate(ctx context.Context) func(interface{}, interface{}) {
	slog.Info("Watcher::buildEventHandlerForUpdate")
	return func(previous, updated interface{}) {
		// TODO NLB-5435 Check for user removing annotation and send delete request to dataplane API
		service := updated.(*v1.Service)
		if !w.isDesiredService(service) {
			return
		}

		nodeIps, err := w.retrieveNodeIps(ctx)
		if err != nil {
			slog.Error("error occurred retrieving node ips", "error", err)
			return
		}

		previousService := previous.(*v1.Service)
		e := core.NewEvent(core.Updated, service, previousService, nodeIps)
		w.handler.AddRateLimitedEvent(&e)
	}
}

// buildInformer creates the informer used to watch for changes to Kubernetes resources.
func (w *Watcher) buildInformer() cache.SharedIndexInformer {
	slog.Debug("Watcher::buildInformer")

	factory := informers.NewSharedInformerFactoryWithOptions(
		w.k8sClient, w.settings.Watcher.ResyncPeriod,
	)
	informer := factory.Core().V1().Services().Informer()

	return informer
}

// initializeEventListeners initializes the event listeners for the informer.
func (w *Watcher) initializeEventListeners(ctx context.Context) error {
	slog.Debug("Watcher::initializeEventListeners")
	var err error

	handlers := cache.ResourceEventHandlerFuncs{
		AddFunc:    w.buildEventHandlerForAdd(ctx),
		DeleteFunc: w.buildEventHandlerForDelete(ctx),
		UpdateFunc: w.buildEventHandlerForUpdate(ctx),
	}

	w.eventHandlerRegistration, err = w.informer.AddEventHandler(handlers)
	if err != nil {
		return fmt.Errorf(`error occurred adding event handlers: %w`, err)
	}

	return nil
}

// notMasterNode retrieves the IP Addresses of the nodes in the cluster. Currently, the master node is excluded. This is
// because the master node may or may not be a worker node and thus may not be able to route traffic.
func (w *Watcher) retrieveNodeIps(ctx context.Context) ([]string, error) {
	started := time.Now()
	slog.Debug("Watcher::retrieveNodeIps")

	var nodeIps []string

	nodes, err := w.k8sClient.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		slog.Error("error occurred retrieving the list of nodes", "error", err)
		return nil, err
	}

	for _, node := range nodes.Items {
		// this is kind of a broad assumption, should probably make this a configurable option
		if w.notMasterNode(node) {
			for _, address := range node.Status.Addresses {
				if address.Type == v1.NodeInternalIP {
					nodeIps = append(nodeIps, address.Address)
				}
			}
		}
	}

	slog.Debug("Watcher::retrieveNodeIps duration", "duration", time.Since(started).Nanoseconds())

	return nodeIps, nil
}

// notMasterNode determines if the node is a master node.
func (w *Watcher) notMasterNode(node v1.Node) bool {
	slog.Debug("Watcher::notMasterNode")

	_, found := node.Labels["node-role.kubernetes.io/master"]

	return !found
}
