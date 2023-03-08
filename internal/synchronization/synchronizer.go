// Copyright 2023 f5 Inc. All rights reserved.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE file.

package synchronization

import (
	"fmt"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/communication"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/configuration"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/core"
	nginxClient "github.com/nginxinc/nginx-plus-go-client/client"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/workqueue"
)

type Interface interface {
	AddEvents(events core.ServerUpdateEvents)
	AddEvent(event *core.ServerUpdateEvent)
	Run(stopCh <-chan struct{})
	ShutDown()
}

type Synchronizer struct {
	eventQueue workqueue.RateLimitingInterface
	settings   *configuration.Settings
}

func NewSynchronizer(settings *configuration.Settings, eventQueue workqueue.RateLimitingInterface) (*Synchronizer, error) {
	synchronizer := Synchronizer{
		eventQueue: eventQueue,
		settings:   settings,
	}

	return &synchronizer, nil
}

func (s *Synchronizer) AddEvents(events core.ServerUpdateEvents) {
	logrus.Debugf(`Synchronizer::AddEvents adding %d events`, len(events))

	if len(s.settings.NginxPlusHosts) == 0 {
		logrus.Warnf(`No Nginx Plus hosts were specified. Skipping synchronization.`)
		return
	}

	updatedEvents := s.fanOutEventToHosts(events)

	for _, event := range updatedEvents {
		s.AddEvent(event)
	}
}

func (s *Synchronizer) AddEvent(event *core.ServerUpdateEvent) {
	logrus.Debugf(`Synchronizer::AddEvent: %#v`, event)

	if event.NginxHost == `` {
		logrus.Warnf(`Nginx host was not specified. Skipping synchronization.`)
		return
	}

	after := RandomMilliseconds(s.settings.Synchronizer.MinMillisecondsJitter, s.settings.Synchronizer.MaxMillisecondsJitter)
	s.eventQueue.AddAfter(event, after)
}

func (s *Synchronizer) Run(stopCh <-chan struct{}) {
	logrus.Debug(`Synchronizer::Run`)

	for i := 0; i < s.settings.Synchronizer.Threads; i++ {
		go wait.Until(s.worker, 0, stopCh)
	}

	<-stopCh
}

func (s *Synchronizer) ShutDown() {
	logrus.Debugf(`Synchronizer::ShutDown`)
	s.eventQueue.ShutDownWithDrain()
}

func (s *Synchronizer) buildNginxPlusClient(nginxHost string) (*nginxClient.NginxClient, error) {
	logrus.Debugf(`Synchronizer::buildNginxPlusClient for host: %s`, nginxHost)

	var err error

	httpClient, err := communication.NewHttpClient()
	if err != nil {
		return nil, fmt.Errorf(`error creating HTTP client: %v`, err)
	}

	client, err := nginxClient.NewNginxClient(httpClient, nginxHost)
	if err != nil {
		return nil, fmt.Errorf(`error creating Nginx Plus client: %v`, err)
	}

	return client, nil
}

func (s *Synchronizer) fanOutEventToHosts(event core.ServerUpdateEvents) core.ServerUpdateEvents {
	logrus.Debugf(`Synchronizer::fanOutEventToHosts: %#v`, event)

	var events core.ServerUpdateEvents

	for hidx, host := range s.settings.NginxPlusHosts {
		for eidx, event := range event {
			id := fmt.Sprintf(`[%d:%d]-[%s]-[%s]-[%s]`, hidx, eidx, RandomString(12), event.UpstreamName, host)
			updatedEvent := core.ServerUpdateEventWithIdAndHost(event, id, host)

			events = append(events, updatedEvent)
		}
	}

	return events
}

func (s *Synchronizer) handleEvent(event *core.ServerUpdateEvent) error {
	logrus.Debugf(`Synchronizer::handleEvent: Id: %s`, event.Id)

	var err error

	switch event.Type {
	case core.Created:
		fallthrough

	case core.Updated:
		err = s.handleCreatedUpdatedEvent(event)

	case core.Deleted:
		err = s.handleDeletedEvent(event)

	default:
		logrus.Warnf(`Synchronizer::handleEvent: unknown event type: %d`, event.Type)
	}

	if err == nil {
		logrus.Infof(`Synchronizer::handleEvent: successfully %s the nginx+ host(s) for Upstream: %s: Id(%s)`, event.TypeName(), event.UpstreamName, event.Id)
	}

	return err
}

func (s *Synchronizer) handleCreatedUpdatedEvent(serverUpdateEvent *core.ServerUpdateEvent) error {
	logrus.Debugf(`Synchronizer::handleCreatedUpdatedEvent: Id: %s`, serverUpdateEvent.Id)

	var err error

	client, err := s.buildNginxPlusClient(serverUpdateEvent.NginxHost)
	if err != nil {
		return fmt.Errorf(`error occurred building the nginx+ client: %w`, err)
	}

	_, _, _, err = client.UpdateStreamServers(serverUpdateEvent.UpstreamName, serverUpdateEvent.Servers)
	if err != nil {
		return fmt.Errorf(`error occurred updating the nginx+ upstream servers: %w`, err)
	}

	return nil
}

func (s *Synchronizer) handleDeletedEvent(serverUpdateEvent *core.ServerUpdateEvent) error {
	logrus.Debugf(`Synchronizer::handleDeletedEvent: Id: %s`, serverUpdateEvent.Id)

	var err error

	client, err := s.buildNginxPlusClient(serverUpdateEvent.NginxHost)
	if err != nil {
		return fmt.Errorf(`error occurred building the nginx+ client: %w`, err)
	}

	err = client.DeleteStreamServer(serverUpdateEvent.UpstreamName, serverUpdateEvent.Servers[0].Server)
	if err != nil {
		return fmt.Errorf(`error occurred deleting the nginx+ upstream server: %w`, err)
	}

	return nil
}

func (s *Synchronizer) handleNextEvent() bool {
	logrus.Debug(`Synchronizer::handleNextEvent`)

	evt, quit := s.eventQueue.Get()
	if quit {
		return false
	}

	defer s.eventQueue.Done(evt)

	event := evt.(*core.ServerUpdateEvent)
	s.withRetry(s.handleEvent(event), event)

	return true
}

func (s *Synchronizer) worker() {
	logrus.Debug(`Synchronizer::worker`)
	for s.handleNextEvent() {
	}
}

func (s *Synchronizer) withRetry(err error, event *core.ServerUpdateEvent) {
	logrus.Debug("Synchronizer::withRetry")
	if err != nil {
		// TODO: Add Telemetry
		if s.eventQueue.NumRequeues(event) < s.settings.Synchronizer.RetryCount { // TODO: Make this configurable
			s.eventQueue.AddRateLimited(event)
			logrus.Infof(`Synchronizer::withRetry: requeued event: %s; error: %v`, event.Id, err)
		} else {
			s.eventQueue.Forget(event)
			logrus.Warnf(`Synchronizer::withRetry: event %#v has been dropped due to too many retries`, event)
		}
	} else {
		s.eventQueue.Forget(event)
	} // TODO: Add error logging
}
