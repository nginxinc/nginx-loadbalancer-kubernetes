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
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	coreinformers "k8s.io/client-go/informers/core/v1"
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
}

// NewWatcher creates a new Watcher
func NewWatcher(
	settings configuration.Settings,
	handler HandlerInterface,
	serviceInformer coreinformers.ServiceInformer,
) (*Watcher, error) {
	if serviceInformer == nil {
		return nil, fmt.Errorf("service informer cannot be nil")
	}

	servicesInformer := serviceInformer.Informer()

	w := &Watcher{
		handler:          handler,
		settings:         settings,
		servicesInformer: servicesInformer,
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

// buildServiceEventHandlerForAdd creates a function that is used as an event handler
// for the informer when Add events are raised.
func (w *Watcher) buildServiceEventHandlerForAdd() func(interface{}) {
	slog.Info("Watcher::buildServiceEventHandlerForAdd")
	return func(obj interface{}) {
		service := obj.(*v1.Service)
		if !w.isDesiredService(service) {
			return
		}

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
		service := updated.(*v1.Service)
		if !w.isDesiredService(service) {
			return
		}

		previousService := previous.(*v1.Service)
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

	handlers := cache.ResourceEventHandlerFuncs{
		AddFunc:    w.buildServiceEventHandlerForAdd(),
		DeleteFunc: w.buildServiceEventHandlerForDelete(),
		UpdateFunc: w.buildServiceEventHandlerForUpdate(),
	}

	_, err = servicesInformer.AddEventHandler(handlers)
	if err != nil {
		return fmt.Errorf(`error occurred adding event handlers: %w`, err)
	}

	return nil
}
