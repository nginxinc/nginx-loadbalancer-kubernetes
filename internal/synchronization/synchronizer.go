/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package synchronization

import (
	"fmt"

	"github.com/nginxinc/kubernetes-nginx-ingress/internal/application"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/communication"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/configuration"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/core"
	nginxClient "github.com/nginxinc/nginx-plus-go-client/client"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/workqueue"
)

// Interface defines the interface needed to implement a synchronizer.
type Interface interface {

	// AddEvents adds a list of events to the queue.
	AddEvents(events core.ServerUpdateEvents)

	// AddEvent adds an event to the queue.
	AddEvent(event *core.ServerUpdateEvent)

	// Run starts the synchronizer.
	Run(stopCh <-chan struct{})

	// ShutDown shuts down the synchronizer.
	ShutDown()
}

// Synchronizer is responsible for synchronizing the state of the Border Servers.
// Operating against the "nlk-synchronizer", it handles events by creating a Border Client as specified in the
// Service annotation for the Upstream. see application/border_client.go and application/application_constants.go for details.
type Synchronizer struct {
	eventQueue workqueue.RateLimitingInterface
	settings   *configuration.Settings
}

// NewSynchronizer creates a new Synchronizer.
func NewSynchronizer(settings *configuration.Settings, eventQueue workqueue.RateLimitingInterface) (*Synchronizer, error) {
	synchronizer := Synchronizer{
		eventQueue: eventQueue,
		settings:   settings,
	}

	return &synchronizer, nil
}

// AddEvents adds a list of events to the queue. If no hosts are specified this is a null operation.
// Events will fan out to the number of hosts specified before being added to the queue.
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

// AddEvent adds an event to the queue. If no hosts are specified this is a null operation.
// Events will be added to the queue after a random delay between MinMillisecondsJitter and MaxMillisecondsJitter.
func (s *Synchronizer) AddEvent(event *core.ServerUpdateEvent) {
	logrus.Debugf(`Synchronizer::AddEvent: %#v`, event)

	if event.NginxHost == `` {
		logrus.Warnf(`Nginx host was not specified. Skipping synchronization.`)
		return
	}

	after := RandomMilliseconds(s.settings.Synchronizer.MinMillisecondsJitter, s.settings.Synchronizer.MaxMillisecondsJitter)
	s.eventQueue.AddAfter(event, after)
}

// Run starts the Synchronizer, spins up Goroutines to process events, and waits for a stop signal.
func (s *Synchronizer) Run(stopCh <-chan struct{}) {
	logrus.Debug(`Synchronizer::Run`)

	for i := 0; i < s.settings.Synchronizer.Threads; i++ {
		go wait.Until(s.worker, 0, stopCh)
	}

	<-stopCh
}

// ShutDown stops the Synchronizer and shuts down the event queue
func (s *Synchronizer) ShutDown() {
	logrus.Debugf(`Synchronizer::ShutDown`)
	s.eventQueue.ShutDownWithDrain()
}

// buildBorderClient creates a Border Client for the specified event.
// NOTE: There is an open issue (https://github.com/nginxinc/nginx-loadbalancer-kubernetes/issues/36) to move creation
// of the underlying Border Server client to the NewBorderClient function.
func (s *Synchronizer) buildBorderClient(event *core.ServerUpdateEvent) (application.Interface, error) {
	logrus.Debugf(`Synchronizer::buildBorderClient`)

	var err error

	httpClient, err := communication.NewHttpClient(s.settings)
	if err != nil {
		return nil, fmt.Errorf(`error creating HTTP client: %v`, err)
	}

	ngxClient, err := nginxClient.NewNginxClient(httpClient, event.NginxHost)
	if err != nil {
		return nil, fmt.Errorf(`error creating Nginx Plus client: %v`, err)
	}

	return application.NewBorderClient(event.ClientType, ngxClient)
}

// fanOutEventToHosts takes a list of events and returns a list of events, one for each Border Server.
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

// handleEvent dispatches an event to the proper handler function.
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

// handleCreatedUpdatedEvent handles events of type Created or Updated.
func (s *Synchronizer) handleCreatedUpdatedEvent(serverUpdateEvent *core.ServerUpdateEvent) error {
	logrus.Debugf(`Synchronizer::handleCreatedUpdatedEvent: Id: %s`, serverUpdateEvent.Id)

	var err error

	borderClient, err := s.buildBorderClient(serverUpdateEvent)
	if err != nil {
		return fmt.Errorf(`error occurred creating the border client: %w`, err)
	}

	if err = borderClient.Update(serverUpdateEvent); err != nil {
		return fmt.Errorf(`error occurred updating the %s upstream servers: %w`, serverUpdateEvent.ClientType, err)
	}

	return nil
}

// handleDeletedEvent handles events of type Deleted.
func (s *Synchronizer) handleDeletedEvent(serverUpdateEvent *core.ServerUpdateEvent) error {
	logrus.Debugf(`Synchronizer::handleDeletedEvent: Id: %s`, serverUpdateEvent.Id)

	var err error

	borderClient, err := s.buildBorderClient(serverUpdateEvent)
	if err != nil {
		return fmt.Errorf(`error occurred creating the border client: %w`, err)
	}

	if err = borderClient.Delete(serverUpdateEvent); err != nil {
		return fmt.Errorf(`error occurred deleting the %s upstream servers: %w`, serverUpdateEvent.ClientType, err)
	}

	return nil
}

// handleNextEvent pulls an event from the event queue and feeds it to the event handler with retry logic
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

// worker is the main message loop
func (s *Synchronizer) worker() {
	logrus.Debug(`Synchronizer::worker`)
	for s.handleNextEvent() {
	}
}

// withRetry handles errors from the event handler and requeues events that fail
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
