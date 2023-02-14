// Copyright 2023 f5 Inc. All rights reserved.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE file.

package configuration

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"strings"
)

const (
	ConfigMapsNamespace = "nkl"
	ResyncPeriod        = 0
)

type Settings struct {
	ctx                      context.Context
	NginxPlusHosts           []string
	k8sClient                *kubernetes.Clientset
	informer                 cache.SharedInformer
	eventHandlerRegistration cache.ResourceEventHandlerRegistration
}

func NewSettings(ctx context.Context, k8sClient *kubernetes.Clientset) (*Settings, error) {
	config := new(Settings)

	config.k8sClient = k8sClient
	config.ctx = ctx

	return config, nil
}

func (s *Settings) Initialize() error {
	logrus.Info("Settings::Initialize")

	var err error

	informer, err := s.buildInformer()
	if err != nil {
		return fmt.Errorf(`error occurred building ConfigMap informer: %w`, err)
	}

	s.informer = informer

	err = s.initializeEventListeners()
	if err != nil {
		return fmt.Errorf(`error occurred initializing event listeners: %w`, err)
	}

	return nil
}

func (s *Settings) Run() {
	logrus.Debug("Settings::Run")

	defer utilruntime.HandleCrash()

	go s.informer.Run(s.ctx.Done())

	<-s.ctx.Done()
}

func (s *Settings) buildInformer() (cache.SharedInformer, error) {
	options := informers.WithNamespace(ConfigMapsNamespace)
	factory := informers.NewSharedInformerFactoryWithOptions(s.k8sClient, ResyncPeriod, options)
	informer := factory.Core().V1().ConfigMaps().Informer()

	return informer, nil
}

func (s *Settings) initializeEventListeners() error {
	logrus.Debug("Settings::initializeEventListeners")

	var err error

	handlers := cache.ResourceEventHandlerFuncs{
		AddFunc:    s.handleAddEvent,
		UpdateFunc: s.handleUpdateEvent,
		DeleteFunc: s.handleDeleteEvent,
	}

	s.eventHandlerRegistration, err = s.informer.AddEventHandler(handlers)
	if err != nil {
		return fmt.Errorf(`error occurred registering event handlers: %w`, err)
	}

	return nil
}

func (s *Settings) handleAddEvent(obj interface{}) {
	logrus.Debug("Settings::handleAddEvent")

	s.handleUpdateEvent(obj, nil)
}

func (s *Settings) handleDeleteEvent(_ interface{}) {
	logrus.Debug("Settings::handleDeleteEvent")

	s.updateHosts([]string{})
}

func (s *Settings) handleUpdateEvent(obj interface{}, _ interface{}) {
	logrus.Debug("Settings::handleUpdateEvent")

	configMap, ok := obj.(*corev1.ConfigMap)
	if !ok {
		logrus.Errorf("Settings::handleUpdateEvent: could not convert obj to ConfigMap")
		return
	}

	hosts, found := configMap.Data["nginx-hosts"]
	if !found {
		logrus.Errorf("Settings::handleUpdateEvent: nginx-hosts key not found in ConfigMap")
		return
	}

	newHosts := s.parseHosts(hosts)
	s.updateHosts(newHosts)
}

func (s *Settings) parseHosts(hosts string) []string {
	return strings.Split(hosts, ",")
}

func (s *Settings) updateHosts(hosts []string) {
	s.NginxPlusHosts = hosts
}
