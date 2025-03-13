/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package synchronization

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	nginxClient "github.com/nginx/nginx-plus-go-client/v2/client"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/application"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/communication"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/configuration"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/core"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/util/workqueue"
)

// Interface defines the interface needed to implement a synchronizer.
type Interface interface {
	// AddEvent adds an event to the queue.
	AddEvent(event core.Event)

	// Run starts the synchronizer.
	Run(ctx context.Context) error

	// ShutDown shuts down the synchronizer.
	ShutDown()
}

type Translator interface {
	Translate(*core.Event) (core.ServerUpdateEvents, error)
}

type ServiceKey struct {
	Name      string
	Namespace string
}

// Synchronizer is responsible for synchronizing the state of the Border Servers.
// Operating against the "nlk-synchronizer", it handles events by creating
// a Border Client as specified in the Service annotation for the Upstream.
// See application/border_client.go and application/application_constants.go for details.
type Synchronizer struct {
	eventQueue    workqueue.TypedRateLimitingInterface[ServiceKey]
	settings      configuration.Settings
	translator    Translator
	cache         *cache
	serviceLister corelisters.ServiceLister
}

// NewSynchronizer creates a new Synchronizer.
func NewSynchronizer(
	settings configuration.Settings,
	eventQueue workqueue.TypedRateLimitingInterface[ServiceKey],
	translator Translator,
	serviceLister corelisters.ServiceLister,
) (*Synchronizer, error) {
	synchronizer := Synchronizer{
		eventQueue:    eventQueue,
		settings:      settings,
		cache:         newCache(),
		translator:    translator,
		serviceLister: serviceLister,
	}

	return &synchronizer, nil
}

// AddEvent adds an event to the rate-limited queue. If no hosts are specified this is a null operation.
func (s *Synchronizer) AddEvent(event core.Event) {
	slog.Debug(`Synchronizer::AddEvent`)

	if len(s.settings.NginxPlusHosts) == 0 {
		slog.Warn(`No Nginx Plus hosts were specified. Skipping synchronization.`)
		return
	}

	key := ServiceKey{Name: event.Service.Name, Namespace: event.Service.Namespace}
	var deletedAt time.Time
	if event.Type == core.Deleted {
		deletedAt = time.Now()
	}

	s.cache.add(key, service{event.Service, deletedAt})
	s.eventQueue.AddRateLimited(key)
}

// Run starts the Synchronizer, spins up Goroutines to process events, and waits for a stop signal.
func (s *Synchronizer) Run(ctx context.Context) error {
	slog.Debug(`Synchronizer::Run`)

	// worker is the main message loop
	worker := func() {
		slog.Debug(`Synchronizer::worker`)
		for s.handleNextServiceEvent(ctx) {
		}
	}

	for i := 0; i < s.settings.Synchronizer.Threads; i++ {
		go wait.Until(worker, 0, ctx.Done())
	}

	<-ctx.Done()
	return nil
}

// ShutDown stops the Synchronizer and shuts down the event queue
func (s *Synchronizer) ShutDown() {
	slog.Debug(`Synchronizer::ShutDown`)
	s.eventQueue.ShutDownWithDrain()
}

// buildBorderClient creates a Border Client for the specified event.
// NOTE: There is an open issue (https://github.com/nginxinc/nginx-loadbalancer-kubernetes/issues/36) to move creation
// of the underlying Border Server client to the NewBorderClient function.
func (s *Synchronizer) buildBorderClient(event *core.ServerUpdateEvent) (application.Interface, error) {
	slog.Debug(`Synchronizer::buildBorderClient`)

	var err error

	httpClient, err := communication.NewHTTPClient(s.settings)
	if err != nil {
		return nil, fmt.Errorf(`error creating HTTP client: %v`, err)
	}

	ngxClient, err := nginxClient.NewNginxClient(event.NginxHost, nginxClient.WithHTTPClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf(`error creating Nginx Plus client: %v`, err)
	}

	return application.NewBorderClient(event.ClientType, ngxClient)
}

// fanOutEventToHosts takes a list of events and returns a list of events, one for each Border Server.
func (s *Synchronizer) fanOutEventToHosts(event core.ServerUpdateEvents) core.ServerUpdateEvents {
	slog.Debug(`Synchronizer::fanOutEventToHosts`)

	var events core.ServerUpdateEvents

	for _, host := range s.settings.NginxPlusHosts {
		for _, event := range event {
			updatedEvent := core.ServerUpdateEventWithHost(event, host)

			events = append(events, updatedEvent)
		}
	}

	return events
}

// handleServiceEvent gets the latest state for the service from the shared
// informer cache, translates the service event into server update events and
// dispatches these events to the proper handler function.
func (s *Synchronizer) handleServiceEvent(ctx context.Context, key ServiceKey) (err error) {
	logger := slog.With("service", key)
	logger.Debug(`Synchronizer::handleServiceEvent`)

	// if a service exists in the shared informer cache, we can assume that we need to update it
	event := core.Event{Type: core.Updated}

	cachedService, exists := s.cache.get(key)

	namespaceLister := s.serviceLister.Services(key.Namespace)
	k8sService, err := namespaceLister.Get(key.Name)
	switch {
	// the service has been deleted. We need to rely on the local cache to
	// gather the last known state of the service so we can delete its
	// upstream servers
	case err != nil && apierrors.IsNotFound(err):
		if !exists {
			logger.Warn(`Synchronizer::handleServiceEvent: no information could be gained about service`)
			return nil
		}
		// no matter what type the cached event has, the service no longer exists, so the type is Deleted
		event.Type = core.Deleted
		event.Service = cachedService.service
	case err != nil:
		return err
	case exists && !cachedService.removedAt.IsZero():
		event.Type = core.Deleted
		event.Service = cachedService.service
	default:
		event.Service = k8sService
	}

	events, err := s.translator.Translate(&event)
	if err != nil {
		return err
	}

	if len(events) == 0 {
		slog.Warn("Synchronizer::handleServiceEvent: no events to process")
		return nil
	}

	events = s.fanOutEventToHosts(events)

	for _, evt := range events {
		switch event.Type {
		case core.Created, core.Updated:
			if handleErr := s.handleCreatedUpdatedEvent(ctx, evt); handleErr != nil {
				err = errors.Join(err, handleErr)
			}
		case core.Deleted:
			if handleErr := s.handleDeletedEvent(ctx, evt); handleErr != nil {
				err = errors.Join(err, handleErr)
			}
		default:
			slog.Warn(`Synchronizer::handleServiceEvent: unknown event type`, "type", event.Type)
		}
	}

	if err != nil {
		return err
	}

	if event.Type == core.Deleted {
		s.cache.delete(ServiceKey{Name: event.Service.Name, Namespace: event.Service.Namespace})
	}

	slog.Debug(
		"Synchronizer::handleServiceEvent: successfully handled the service change", "service", key,
	)

	return nil
}

// handleCreatedUpdatedEvent handles events of type Created or Updated.
func (s *Synchronizer) handleCreatedUpdatedEvent(ctx context.Context, serverUpdateEvent *core.ServerUpdateEvent) error {
	slog.Debug(`Synchronizer::handleCreatedUpdatedEvent`)

	var err error

	borderClient, err := s.buildBorderClient(serverUpdateEvent)
	if err != nil {
		return fmt.Errorf(`error occurred creating the border client: %w`, err)
	}

	if err = borderClient.Update(ctx, serverUpdateEvent); err != nil {
		return fmt.Errorf(`error occurred updating the %s upstream servers: %w`, serverUpdateEvent.ClientType, err)
	}

	return nil
}

// handleDeletedEvent handles events of type Deleted.
func (s *Synchronizer) handleDeletedEvent(ctx context.Context, serverUpdateEvent *core.ServerUpdateEvent) error {
	slog.Debug(`Synchronizer::handleDeletedEvent`)

	var err error

	borderClient, err := s.buildBorderClient(serverUpdateEvent)
	if err != nil {
		return fmt.Errorf(`error occurred creating the border client: %w`, err)
	}

	err = borderClient.Update(ctx, serverUpdateEvent)

	switch {
	case err == nil:
		return nil
	// checking the string is not ideal, but the plus client gives us no option
	case strings.Contains(err.Error(), "status=404"):
		return nil
	default:
		return fmt.Errorf(`error occurred deleting the %s upstream servers: %w`, serverUpdateEvent.ClientType, err)
	}
}

// handleNextServiceEvent pulls a service from the event queue and feeds it to
// the service event handler with retry logic
func (s *Synchronizer) handleNextServiceEvent(ctx context.Context) bool {
	slog.Debug(`Synchronizer::handleNextServiceEvent`)

	svc, quit := s.eventQueue.Get()
	if quit {
		return false
	}

	defer s.eventQueue.Done(svc)

	s.withRetry(s.handleServiceEvent(ctx, svc), svc)

	return true
}

// withRetry handles errors from the event handler and requeues events that fail
func (s *Synchronizer) withRetry(err error, key ServiceKey) {
	slog.Debug("Synchronizer::withRetry")
	if err != nil {
		// TODO: Add Telemetry
		s.eventQueue.AddRateLimited(key)
		slog.Info(`Synchronizer::withRetry: requeued service update`, "service", key, "error", err)
	} else {
		s.eventQueue.Forget(key)
	} // TODO: Add error logging
}
