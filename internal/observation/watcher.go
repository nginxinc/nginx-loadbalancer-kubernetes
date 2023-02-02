// Copyright 2023 f5 Inc. All rights reserved.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE file.

package observation

import (
	"context"
	"fmt"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/core"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

const NginxIngressNamespace = "nginx-ingress"
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

	nodeIps, err := w.retrieveNodeIps()
	if err != nil {
		return fmt.Errorf(`error occurred retrieving node ips: %w`, err)
	}

	logrus.Infof("Watcher::Watch::nodeIps: %v", nodeIps)
	
	go w.informer.Run(w.ctx.Done())

	if !cache.WaitForNamedCacheSync(WatcherQueueName, w.ctx.Done(), w.informer.HasSynced) {
		return fmt.Errorf(`error occurred waiting for the cache to sync`)
	}

	<-w.ctx.Done()
	return nil
}

func (w *Watcher) buildEventHandlerForAdd() func(interface{}) {
	logrus.Info("Watcher::buildEventHandlerForAdd")
	return func(obj interface{}) {
		service := obj.(*v1.Service)
		var previousService *v1.Service
		e := core.NewEvent(core.Created, service, previousService)
		w.handler.AddRateLimitedEvent(&e)
	}
}

func (w *Watcher) buildEventHandlerForDelete() func(interface{}) {
	logrus.Info("Watcher::buildEventHandlerForDelete")
	return func(obj interface{}) {
		service := obj.(*v1.Service)
		var previousService *v1.Service
		e := core.NewEvent(core.Deleted, service, previousService)
		w.handler.AddRateLimitedEvent(&e)
	}
}

func (w *Watcher) buildEventHandlerForUpdate() func(interface{}, interface{}) {
	logrus.Info("Watcher::buildEventHandlerForUpdate")
	return func(previous, updated interface{}) {
		service := updated.(*v1.Service)
		previousService := previous.(*v1.Service)
		e := core.NewEvent(core.Updated, service, previousService)
		w.handler.AddRateLimitedEvent(&e)
	}
}

func (w *Watcher) buildInformer() (cache.SharedIndexInformer, error) {
	logrus.Debug("Watcher::buildInformer")

	options := informers.WithNamespace(NginxIngressNamespace)
	factory := informers.NewSharedInformerFactoryWithOptions(w.client, ResyncPeriod, options)
	informer := factory.Core().V1().Services().Informer()

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

func (w *Watcher) retrieveNodeIps() ([]string, error) {
	logrus.Debug("Watcher::retrieveNodeIps")

	var nodeIps []string

	nodes, err := w.client.CoreV1().Nodes().List(w.ctx, metav1.ListOptions{})
	if err != nil {
		logrus.Errorf(`error occurred retrieving the list of nodes: %v`, err)
		return nil, err
	}

	for _, node := range nodes.Items {
		for _, address := range node.Status.Addresses {
			if address.Type == v1.NodeInternalIP {
				nodeIps = append(nodeIps, address.Address)
			}
		}
	}

	return nodeIps, nil
}
