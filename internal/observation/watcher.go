/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package observation

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/nginxinc/kubernetes-nginx-ingress/internal/configuration"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/core"
	v1 "k8s.io/api/core/v1"
	discovery "k8s.io/api/discovery/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	coreinformers "k8s.io/client-go/informers/core/v1"
	discoveryinformers "k8s.io/client-go/informers/discovery/v1"
	"k8s.io/client-go/tools/cache"
)

// Watcher is responsible for watching for changes to Kubernetes resources.
// Particularly, Services in the namespace defined in the WatcherSettings::NginxIngressNamespace setting.
// When a change is detected, an Event is generated and added to the Handler's queue.
type Watcher struct {
	// handler is the event handler
	handler HandlerInterface

	// settings is the configuration settings
	settings configuration.Settings

	// servicesInformer is the informer used to watch for changes to services
	servicesInformer cache.SharedIndexInformer

	// endpointSliceInformer is the informer used to watch for changes to endpoint slices
	endpointSliceInformer cache.SharedIndexInformer

	// nodesInformer is the informer used to watch for changes to nodes
	nodesInformer cache.SharedIndexInformer

	register *register
}

// NewWatcher creates a new Watcher
func NewWatcher(
	settings configuration.Settings,
	handler HandlerInterface,
	serviceInformer coreinformers.ServiceInformer,
	endpointSliceInformer discoveryinformers.EndpointSliceInformer,
	nodeInformer coreinformers.NodeInformer,
) (*Watcher, error) {
	if serviceInformer == nil {
		return nil, fmt.Errorf("service informer cannot be nil")
	}

	if endpointSliceInformer == nil {
		return nil, fmt.Errorf("endpoint slice informer cannot be nil")
	}

	if nodeInformer == nil {
		return nil, fmt.Errorf("node informer cannot be nil")
	}

	servicesInformer := serviceInformer.Informer()
	endpointSlicesInformer := endpointSliceInformer.Informer()
	nodesInformer := nodeInformer.Informer()

	w := &Watcher{
		handler:               handler,
		settings:              settings,
		servicesInformer:      servicesInformer,
		endpointSliceInformer: endpointSlicesInformer,
		nodesInformer:         nodesInformer,
		register:              newRegister(),
	}

	if err := w.initializeEventListeners(servicesInformer); err != nil {
		return nil, err
	}

	return w, nil
}

// Run starts the process of watching for changes to Kubernetes resources.
// Initialize must be called before Watch.
func (w *Watcher) Run(ctx context.Context) error {
	if w.servicesInformer == nil {
		return fmt.Errorf(`servicesInformer is nil`)
	}

	slog.Debug("Watcher::Watch")

	defer utilruntime.HandleCrash()
	defer w.handler.ShutDown()

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

func (w *Watcher) buildNodesEventHandlerForAdd() func(interface{}) {
	slog.Info("Watcher::buildNodesEventHandlerForAdd")
	return func(obj interface{}) {
		slog.Debug("received node add event")
		for _, service := range w.register.listServices() {
			var previousService *v1.Service
			e := core.NewEvent(core.Updated, service, previousService)
			w.handler.AddRateLimitedEvent(&e)
		}
	}
}

func (w *Watcher) buildNodesEventHandlerForUpdate() func(interface{}, interface{}) {
	slog.Info("Watcher::buildNodesEventHandlerForUpdate")
	return func(previous, updated interface{}) {
		slog.Debug("received node update event")
		for _, service := range w.register.listServices() {
			var previousService *v1.Service
			e := core.NewEvent(core.Updated, service, previousService)
			w.handler.AddRateLimitedEvent(&e)
		}
	}
}

func (w *Watcher) buildNodesEventHandlerForDelete() func(interface{}) {
	slog.Info("Watcher::buildNodesEventHandlerForDelete")
	return func(obj interface{}) {
		slog.Debug("received node delete event")
		for _, service := range w.register.listServices() {
			var previousService *v1.Service
			e := core.NewEvent(core.Updated, service, previousService)
			w.handler.AddRateLimitedEvent(&e)
		}
	}
}

func (w *Watcher) buildEndpointSlicesEventHandlerForAdd() func(interface{}) {
	slog.Info("Watcher::buildEndpointSlicesEventHandlerForAdd")
	return func(obj interface{}) {
		slog.Debug("received endpoint slice add event")
		endpointSlice, ok := obj.(*discovery.EndpointSlice)
		if !ok {
			slog.Error("could not convert event object to EndpointSlice", "obj", obj)
			return
		}

		service, ok := w.register.getService(endpointSlice.Namespace, endpointSlice.Labels["kubernetes.io/service-name"])
		if !ok {
			// not interested in any unregistered service
			return
		}

		var previousService *v1.Service
		e := core.NewEvent(core.Updated, service, previousService)
		w.handler.AddRateLimitedEvent(&e)
	}
}

func (w *Watcher) buildEndpointSlicesEventHandlerForUpdate() func(interface{}, interface{}) {
	slog.Info("Watcher::buildEndpointSlicesEventHandlerForUpdate")
	return func(previous, updated interface{}) {
		slog.Debug("received endpoint slice update event")
		endpointSlice, ok := updated.(*discovery.EndpointSlice)
		if !ok {
			slog.Error("could not convert event object to EndpointSlice", "obj", updated)
			return
		}

		service, ok := w.register.getService(endpointSlice.Namespace, endpointSlice.Labels["kubernetes.io/service-name"])
		if !ok {
			// not interested in any unregistered service
			return
		}

		var previousService *v1.Service
		e := core.NewEvent(core.Updated, service, previousService)
		w.handler.AddRateLimitedEvent(&e)
	}
}

func (w *Watcher) buildEndpointSlicesEventHandlerForDelete() func(interface{}) {
	slog.Info("Watcher::buildEndpointSlicesEventHandlerForDelete")
	return func(obj interface{}) {
		slog.Debug("received endpoint slice delete event")
		endpointSlice, ok := obj.(*discovery.EndpointSlice)
		if !ok {
			slog.Error("could not convert event object to EndpointSlice", "obj", obj)
			return
		}

		service, ok := w.register.getService(endpointSlice.Namespace, endpointSlice.Labels["kubernetes.io/service-name"])
		if !ok {
			// not interested in any unregistered service
			return
		}

		var previousService *v1.Service
		e := core.NewEvent(core.Deleted, service, previousService)
		w.handler.AddRateLimitedEvent(&e)
	}
}

// buildServiceEventHandlerForAdd creates a function that is used as an event handler
// for the informer when Add events are raised.
func (w *Watcher) buildServiceEventHandlerForAdd() func(interface{}) {
	slog.Info("Watcher::buildServiceEventHandlerForAdd")
	return func(obj interface{}) {
		service := obj.(*v1.Service)
		if !w.isDesiredService(service) {
			return
		}

		w.register.addOrUpdateService(service)

		var previousService *v1.Service
		e := core.NewEvent(core.Created, service, previousService)
		w.handler.AddRateLimitedEvent(&e)
	}
}

// buildServiceEventHandlerForDelete creates a function that is used as an event handler
// for the informer when Delete events are raised.
func (w *Watcher) buildServiceEventHandlerForDelete() func(interface{}) {
	slog.Info("Watcher::buildServiceEventHandlerForDelete")
	return func(obj interface{}) {
		service := obj.(*v1.Service)
		if !w.isDesiredService(service) {
			return
		}

		w.register.removeService(service)

		var previousService *v1.Service
		e := core.NewEvent(core.Deleted, service, previousService)
		w.handler.AddRateLimitedEvent(&e)
	}
}

// buildServiceEventHandlerForUpdate creates a function that is used as an event handler
// for the informer when Update events are raised.
func (w *Watcher) buildServiceEventHandlerForUpdate() func(interface{}, interface{}) {
	slog.Info("Watcher::buildServiceEventHandlerForUpdate")
	return func(previous, updated interface{}) {
		// TODO NLB-5435 Check for user removing annotation and send delete request to dataplane API
		previousService := previous.(*v1.Service)
		service := updated.(*v1.Service)

		if w.isDesiredService(previousService) && !w.isDesiredService(service) {
			w.register.removeService(previousService)
			return
		}

		if !w.isDesiredService(service) {
			return
		}

		w.register.addOrUpdateService(service)

		e := core.NewEvent(core.Updated, service, previousService)
		w.handler.AddRateLimitedEvent(&e)
	}
}

// initializeEventListeners initializes the event listeners for the informer.
func (w *Watcher) initializeEventListeners(
	servicesInformer cache.SharedIndexInformer,
) error {
	slog.Debug("Watcher::initializeEventListeners")
	var err error

	serviceHandlers := cache.ResourceEventHandlerFuncs{
		AddFunc:    w.buildServiceEventHandlerForAdd(),
		DeleteFunc: w.buildServiceEventHandlerForDelete(),
		UpdateFunc: w.buildServiceEventHandlerForUpdate(),
	}

	endpointSliceHandlers := cache.ResourceEventHandlerFuncs{
		AddFunc:    w.buildEndpointSlicesEventHandlerForAdd(),
		DeleteFunc: w.buildEndpointSlicesEventHandlerForDelete(),
		UpdateFunc: w.buildEndpointSlicesEventHandlerForUpdate(),
	}

	nodeHandlers := cache.ResourceEventHandlerFuncs{
		AddFunc:    w.buildNodesEventHandlerForAdd(),
		DeleteFunc: w.buildNodesEventHandlerForDelete(),
		UpdateFunc: w.buildNodesEventHandlerForUpdate(),
	}

	_, err = servicesInformer.AddEventHandler(serviceHandlers)
	if err != nil {
		return fmt.Errorf(`error occurred adding service event handlers: %w`, err)
	}

	_, err = w.endpointSliceInformer.AddEventHandler(endpointSliceHandlers)
	if err != nil {
		return fmt.Errorf(`error occurred adding endpoint slice event handlers: %w`, err)
	}

	_, err = w.nodesInformer.AddEventHandler(nodeHandlers)
	if err != nil {
		return fmt.Errorf(`error occurred adding node event handlers: %w`, err)
	}

	return nil
}
