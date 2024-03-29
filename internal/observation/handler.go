/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package observation

import (
	"fmt"

	"github.com/nginxinc/kubernetes-nginx-ingress/internal/configuration"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/core"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/synchronization"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/translation"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/workqueue"
)

// HandlerInterface is the interface for the event handler
type HandlerInterface interface {

	// AddRateLimitedEvent defines the interface for adding an event to the event queue
	AddRateLimitedEvent(event *core.Event)

	// Run defines the interface used to start the event handler
	Run(stopCh <-chan struct{})

	// ShutDown defines the interface used to stop the event handler
	ShutDown()
}

// Handler is responsible for processing events in the "nlk-handler" queue.
// When processing a message the Translation module is used to translate the event into an internal representation.
// The translation process may result in multiple events being generated. This fan-out mainly supports the differences
// in NGINX Plus API calls for creating/updating Upstreams and deleting Upstreams.
type Handler struct {

	// eventQueue is the queue used to store events
	eventQueue workqueue.RateLimitingInterface

	// settings is the configuration settings
	settings *configuration.Settings

	// synchronizer is the synchronizer used to synchronize the internal representation with a Border Server
	synchronizer synchronization.Interface
}

// NewHandler creates a new event handler
func NewHandler(settings *configuration.Settings, synchronizer synchronization.Interface, eventQueue workqueue.RateLimitingInterface) *Handler {
	return &Handler{
		eventQueue:   eventQueue,
		settings:     settings,
		synchronizer: synchronizer,
	}
}

// AddRateLimitedEvent adds an event to the event queue
func (h *Handler) AddRateLimitedEvent(event *core.Event) {
	logrus.Debugf(`Handler::AddRateLimitedEvent: %#v`, event)
	h.eventQueue.AddRateLimited(event)
}

// Run starts the event handler, spins up Goroutines to process events, and waits for a stop signal
func (h *Handler) Run(stopCh <-chan struct{}) {
	logrus.Debug("Handler::Run")

	for i := 0; i < h.settings.Handler.Threads; i++ {
		go wait.Until(h.worker, 0, stopCh)
	}

	<-stopCh
}

// ShutDown stops the event handler and shuts down the event queue
func (h *Handler) ShutDown() {
	logrus.Debug("Handler::ShutDown")
	h.eventQueue.ShutDown()
}

// handleEvent feeds translated events to the synchronizer
func (h *Handler) handleEvent(e *core.Event) error {
	logrus.Debugf(`Handler::handleEvent: %#v`, e)
	// TODO: Add Telemetry

	events, err := translation.Translate(e)
	if err != nil {
		return fmt.Errorf(`Handler::handleEvent error translating: %v`, err)
	}

	h.synchronizer.AddEvents(events)

	return nil
}

// handleNextEvent pulls an event from the event queue and feeds it to the event handler with retry logic
func (h *Handler) handleNextEvent() bool {
	logrus.Debug("Handler::handleNextEvent")
	evt, quit := h.eventQueue.Get()
	logrus.Debugf(`Handler::handleNextEvent: %#v, quit: %v`, evt, quit)
	if quit {
		return false
	}

	defer h.eventQueue.Done(evt)

	event := evt.(*core.Event)
	h.withRetry(h.handleEvent(event), event)

	return true
}

// worker is the main message loop
func (h *Handler) worker() {
	for h.handleNextEvent() {
		// TODO: Add Telemetry
	}
}

// withRetry handles errors from the event handler and requeues events that fail
func (h *Handler) withRetry(err error, event *core.Event) {
	logrus.Debug("Handler::withRetry")
	if err != nil {
		// TODO: Add Telemetry
		if h.eventQueue.NumRequeues(event) < h.settings.Handler.RetryCount {
			h.eventQueue.AddRateLimited(event)
			logrus.Infof(`Handler::withRetry: requeued event: %#v; error: %v`, event, err)
		} else {
			h.eventQueue.Forget(event)
			logrus.Warnf(`Handler::withRetry: event %#v has been dropped due to too many retries`, event)
		}
	} // TODO: Add error logging
}
