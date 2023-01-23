// Copyright 2023 f5 Inc. All rights reserved.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE file.

package synchronization

import (
	"fmt"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/communication"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/config"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/core"
	nginxClient "github.com/nginxinc/nginx-plus-go-client/client"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/workqueue"
)

const Threads = 1
const SynchronizerQueueName = `nec-synchronizer`

type Synchronizer struct {
	NginxPlusClient *nginxClient.NginxClient
	eventQueue      workqueue.RateLimitingInterface
}

func NewSynchronizer() (*Synchronizer, error) {
	synchronizer := Synchronizer{}

	return &synchronizer, nil
}

func (s *Synchronizer) AddRateLimitedEvent(event *core.Event) {
	logrus.Infof(`Synchronizer::AddRateLimitedEvent: %#v`, event)
	s.eventQueue.AddRateLimited(event)
}

func (s *Synchronizer) Initialize() error {
	var err error
	settings, err := config.NewSettings()
	if err != nil {
		return fmt.Errorf(`error loading configuration: %v`, err)
	}

	httpClient, err := communication.NewHttpClient()
	if err != nil {
		return fmt.Errorf(`error creating HTTP client: %v`, err)
	}

	s.NginxPlusClient, err = nginxClient.NewNginxClient(httpClient, settings.NginxPlusHost)
	if err != nil {
		return fmt.Errorf(`error creating Nginx Plus client: %v`, err)
	}

	s.eventQueue = workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), SynchronizerQueueName)

	return nil
}

func (s *Synchronizer) Run(stopCh <-chan struct{}) {
	logrus.Info(`Synchronizer::Run`)
	for i := 0; i < Threads; i++ {
		go wait.Until(s.worker, 0, stopCh)
	}

	<-stopCh
}

func (s *Synchronizer) ShutDown() {
	logrus.Infof(`Synchronizer::ShutDown`)
	s.eventQueue.ShutDown()
}

func (s *Synchronizer) handleEvent(event *core.Event) {
	logrus.Info(`Synchronizer::handleEvent`)
	logrus.Infof(`Synchronizer::handleEvent: %#v`, event)
}

func (s *Synchronizer) handleNextEvent() bool {
	logrus.Info(`Synchronizer::handleNextEvent`)
	event, quit := s.eventQueue.Get()
	if quit {
		return false
	}

	defer s.eventQueue.Done(event)

	s.handleEvent(event.(*core.Event))

	return true
}

func (s *Synchronizer) worker() {
	logrus.Info(`Synchronizer::sync`)
	for s.handleNextEvent() {
		// TODO: Add Telemetry
	}
}
