/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package observation

import (
	"errors"
	"fmt"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/configuration"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/core"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"time"
)

type Watcher struct {
	eventHandlerRegistration interface{}
	handler                  HandlerInterface
	informer                 cache.SharedIndexInformer
	settings                 *configuration.Settings
}

func NewWatcher(settings *configuration.Settings, handler HandlerInterface) (*Watcher, error) {
	return &Watcher{
		handler:  handler,
		settings: settings,
	}, nil
}

func (w *Watcher) Initialize() error {
	logrus.Debug("Watcher::Initialize")
	var err error

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

	if w.informer == nil {
		return errors.New("error: Initialize must be called before Watch")
	}

	defer utilruntime.HandleCrash()
	defer w.handler.ShutDown()

	go w.informer.Run(w.settings.Context.Done())

	if !cache.WaitForNamedCacheSync(w.settings.Handler.WorkQueueSettings.Name, w.settings.Context.Done(), w.informer.HasSynced) {
		return fmt.Errorf(`error occurred waiting for the cache to sync`)
	}

	<-w.settings.Context.Done()
	return nil
}

func (w *Watcher) buildEventHandlerForAdd() func(interface{}) {
	logrus.Info("Watcher::buildEventHandlerForAdd")
	return func(obj interface{}) {
		nodeIps, err := w.retrieveNodeIps()
		if err != nil {
			logrus.Errorf(`error occurred retrieving node ips: %v`, err)
			return
		}
		service := obj.(*v1.Service)
		var previousService *v1.Service
		e := core.NewEvent(core.Created, service, previousService, nodeIps)
		w.handler.AddRateLimitedEvent(&e)
	}
}

func (w *Watcher) buildEventHandlerForDelete() func(interface{}) {
	logrus.Info("Watcher::buildEventHandlerForDelete")
	return func(obj interface{}) {
		nodeIps, err := w.retrieveNodeIps()
		if err != nil {
			logrus.Errorf(`error occurred retrieving node ips: %v`, err)
			return
		}
		service := obj.(*v1.Service)
		var previousService *v1.Service
		e := core.NewEvent(core.Deleted, service, previousService, nodeIps)
		w.handler.AddRateLimitedEvent(&e)
	}
}

func (w *Watcher) buildEventHandlerForUpdate() func(interface{}, interface{}) {
	logrus.Info("Watcher::buildEventHandlerForUpdate")
	return func(previous, updated interface{}) {
		nodeIps, err := w.retrieveNodeIps()
		if err != nil {
			logrus.Errorf(`error occurred retrieving node ips: %v`, err)
			return
		}
		service := updated.(*v1.Service)
		previousService := previous.(*v1.Service)
		e := core.NewEvent(core.Updated, service, previousService, nodeIps)
		w.handler.AddRateLimitedEvent(&e)
	}
}

func (w *Watcher) buildInformer() (cache.SharedIndexInformer, error) {
	logrus.Debug("Watcher::buildInformer")

	options := informers.WithNamespace(w.settings.Watcher.NginxIngressNamespace)
	factory := informers.NewSharedInformerFactoryWithOptions(w.settings.K8sClient, w.settings.Watcher.ResyncPeriod, options)
	informer := factory.Core().V1().Services().Informer()

	return informer, nil
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
	started := time.Now()
	logrus.Debug("Watcher::retrieveNodeIps")

	var nodeIps []string

	nodes, err := w.settings.K8sClient.CoreV1().Nodes().List(w.settings.Context, metav1.ListOptions{})
	if err != nil {
		logrus.Errorf(`error occurred retrieving the list of nodes: %v`, err)
		return nil, err
	}

	for _, node := range nodes.Items {
		if w.notMasterNode(node) {
			for _, address := range node.Status.Addresses {
				if address.Type == v1.NodeInternalIP {
					nodeIps = append(nodeIps, address.Address)
				}
			}
		}
	}

	logrus.Debugf("Watcher::retrieveNodeIps duration: %d", time.Since(started).Nanoseconds())

	return nodeIps, nil
}

func (w *Watcher) notMasterNode(node v1.Node) bool {
	logrus.Debug("Watcher::notMasterNode")

	_, found := node.Labels["node-role.kubernetes.io/master"]

	return !found
}
