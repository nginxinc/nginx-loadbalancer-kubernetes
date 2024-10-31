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

	"github.com/nginxinc/kubernetes-nginx-ingress/internal/configuration"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/core"
	v1 "k8s.io/api/core/v1"
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
func (w *Watcher) Initialize() error {
	slog.Debug("Watcher::Initialize")
	var err error

	w.informer = w.buildInformer()

	err = w.initializeEventListeners()
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
func (w *Watcher) buildEventHandlerForAdd() func(interface{}) {
	slog.Info("Watcher::buildEventHandlerForAdd")
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

// buildEventHandlerForDelete creates a function that is used as an event handler
// for the informer when Delete events are raised.
func (w *Watcher) buildEventHandlerForDelete() func(interface{}) {
	slog.Info("Watcher::buildEventHandlerForDelete")
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

// buildEventHandlerForUpdate creates a function that is used as an event handler
// for the informer when Update events are raised.
func (w *Watcher) buildEventHandlerForUpdate() func(interface{}, interface{}) {
	slog.Info("Watcher::buildEventHandlerForUpdate")
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
func (w *Watcher) initializeEventListeners() error {
	slog.Debug("Watcher::initializeEventListeners")
	var err error

	handlers := cache.ResourceEventHandlerFuncs{
		AddFunc:    w.buildEventHandlerForAdd(),
		DeleteFunc: w.buildEventHandlerForDelete(),
		UpdateFunc: w.buildEventHandlerForUpdate(),
	}

	w.eventHandlerRegistration, err = w.informer.AddEventHandler(handlers)
	if err != nil {
		return fmt.Errorf(`error occurred adding event handlers: %w`, err)
	}

	return nil
}
