// Copyright 2023 f5 Inc. All rights reserved.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE file.

package observation

import (
	"context"
	"fmt"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/core"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/networking/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

const ResyncPeriod = 0

type Watcher struct {
	ctx                      context.Context
	client                   *kubernetes.Clientset
	eventHandlerRegistration interface{}
	handler                  *Handler
	informer                 cache.SharedIndexInformer
}

func NewWatcher(ctx context.Context, handler *Handler) (*Watcher, error) {
	return &Watcher{
		ctx:     ctx,
		handler: handler,
	}, nil
}

func (w *Watcher) Initialize() error {
	logrus.Debug("Watcher::Initialize")
	var err error

	w.client, err = w.buildKubernetesClient()
	if err != nil {
		return fmt.Errorf(`initalization error: %w`, err)
	}

	w.informer, err = w.buildInformer()
	if err != nil {
		return fmt.Errorf(`initialization error: %w`, err)
	}

	err = w.initializeEventListeners()
	if err != nil {
		return fmt.Errorf(`initialization error: %w`, err)
	}

	return nil
}

func (w *Watcher) Watch() error {
	logrus.Debug("Watcher::Watch")
	defer utilruntime.HandleCrash()
	defer w.handler.ShutDown()

	go w.informer.Run(w.ctx.Done())

	if !cache.WaitForNamedCacheSync(WatcherQueueName, w.ctx.Done(), w.informer.HasSynced) {
		return fmt.Errorf(`error occurred waiting for the cache to sync`)
	}

	<-w.ctx.Done()
	return nil
}

func (w *Watcher) buildEventHandlerForAdd() func(interface{}) {
	logrus.Debug("Watcher::buildEventHandlerForAdd")
	return func(obj interface{}) {
		ingress := obj.(*v1.Ingress)
		var previousIngress *v1.Ingress
		e := core.NewEvent(core.Created, ingress, previousIngress)
		w.handler.AddRateLimitedEvent(&e)
	}
}

func (w *Watcher) buildEventHandlerForDelete() func(interface{}) {
	logrus.Debug("Watcher::buildEventHandlerForDelete")
	return func(obj interface{}) {
		ingress := obj.(*v1.Ingress)
		var previousIngress *v1.Ingress
		e := core.NewEvent(core.Deleted, ingress, previousIngress)
		w.handler.AddRateLimitedEvent(&e)
	}
}

func (w *Watcher) buildEventHandlerForUpdate() func(interface{}, interface{}) {
	logrus.Debug("Watcher::buildEventHandlerForUpdate")
	return func(previous, updated interface{}) {
		ingress := updated.(*v1.Ingress)
		previousIngress := previous.(*v1.Ingress)
		e := core.NewEvent(core.Updated, ingress, previousIngress)
		w.handler.AddRateLimitedEvent(&e)
	}
}

func (w *Watcher) buildInformer() (cache.SharedIndexInformer, error) {
	logrus.Debug("Watcher::buildInformer")

	factory := informers.NewSharedInformerFactory(w.client, ResyncPeriod)
	informer := factory.Networking().V1().Ingresses().Informer()

	return informer, nil
}

func (w *Watcher) buildKubernetesClient() (*kubernetes.Clientset, error) {
	logrus.Debug("Watcher::buildKubernetesClient")
	k8sConfig, err := rest.InClusterConfig()
	if err == rest.ErrNotInCluster {
		return nil, fmt.Errorf(`not running in a Cluster: %w`, err)
	} else if err != nil {
		return nil, fmt.Errorf(`error occurred getting the Cluster config: %w`, err)
	}

	client, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		return nil, fmt.Errorf(`error occurred creating a client: %w`, err)
	}

	return client, nil
}

func (w *Watcher) initializeEventListeners() error {
	logrus.Debug("Watcher::initializeEventListeners")
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
