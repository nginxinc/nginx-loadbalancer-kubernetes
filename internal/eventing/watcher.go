// Copyright 2023 f5 Inc. All rights reserved.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE file.

package eventing

import (
	"context"
	"fmt"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/synchronization"
	networking "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const ResyncPeriod = 0

type Watcher struct {
	client                   *kubernetes.Clientset
	eventHandlerRegistration interface{}
	eventQueue               workqueue.RateLimitingInterface
	informer                 cache.SharedInformer
	synchronizer             *synchronization.NginxPlusSynchronizer
}

func NewWatcher(synchronizer *synchronization.NginxPlusSynchronizer) (*Watcher, error) {
	return &Watcher{
		synchronizer: synchronizer,
	}, nil
}

func (w *Watcher) Initialize() error {
	var err error

	w.client, err = w.buildKubernetesClient()
	if err != nil {
		return fmt.Errorf(`initalization error: %w`, err)
	}

	w.informer, err = w.buildInformer()
	if err != nil {
		return fmt.Errorf(`initialization error: %w`, err)
	}

	w.eventQueue, err = w.buildEventQueue()

	err = w.initializeEventListeners()
	if err != nil {
		return fmt.Errorf(`initialization error: %w`, err)
	}

	return nil
}

func (w *Watcher) Watch() error {

	return nil
}

func (w *Watcher) buildEventHandlerForAdd() func(interface{}) {
	return func(obj interface{}) {
		e := NewEvent(Created, obj, nil)
		w.eventQueue.AddRateLimited(e)
	}
}

func (w *Watcher) buildEventHandlerForDelete() func(interface{}) {
	return func(obj interface{}) {
		e := NewEvent(Deleted, obj, nil)
		w.eventQueue.AddRateLimited(e)
	}
}

func (w *Watcher) buildEventHandlerForUpdate() func(interface{}, interface{}) {
	return func(previous interface{}, updated interface{}) {
		e := NewEvent(Updated, updated, previous)
		w.eventQueue.AddRateLimited(e)
	}
}

func (w *Watcher) buildEventQueue() (workqueue.RateLimitingInterface, error) {
	eventQueue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	return eventQueue, nil
}

func (w *Watcher) buildInformer() (cache.SharedInformer, error) {
	listerWatcher, err := w.buildListWatcher()
	if err != nil {
		return nil, fmt.Errorf(`error occurred creating a ListerWatcher: %w`, err)
	}

	informer := cache.NewSharedInformer(listerWatcher, &networking.Ingress{}, ResyncPeriod)

	return informer, nil
}

func (w *Watcher) buildKubernetesClient() (*kubernetes.Clientset, error) {
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

func (w *Watcher) buildListWatcher() (*cache.ListWatch, error) {
	lw := cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			return w.client.NetworkingV1beta1().Ingresses(metav1.NamespaceAll).List(context.TODO(), options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return w.client.NetworkingV1beta1().Ingresses(metav1.NamespaceAll).Watch(context.TODO(), options)
		},
	}

	return &lw, nil
}

func (w *Watcher) initializeEventListeners() error {
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
