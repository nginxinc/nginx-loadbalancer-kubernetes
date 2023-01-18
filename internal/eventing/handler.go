// Copyright 2023 f5 Inc. All rights reserved.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE file.

package eventing

import (
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/util/workqueue"
)

const WatcherQueueName = `nginx-k8s-edge-controller`

type Handler struct {
	eventQueue workqueue.RateLimitingInterface
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) AddRateLimitedEvent(event *Event) {
	logrus.Infof(`Handler::AddRateLimitedEvent: %v`, event)
	h.eventQueue.AddRateLimited(event)
}

func (h *Handler) Initialize() {
	h.eventQueue = workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), WatcherQueueName)
}

func (h *Handler) Run() {
	logrus.Info("Handler::Run")
	for h.handleNextEvent() {
		// TODO: Add Telemetry
	}
}

func (h *Handler) ShutDown() {
	logrus.Info("Handler::ShutDown")
	h.eventQueue.ShutDown()
}

func (h *Handler) handleEvent(e *Event) {
	logrus.Info("Handler::handleEvent")
	// TODO: Implement translation logic
	// TODO: Add Telemetry
	logrus.Infof(`processing event: %v`, e)
}

func (h *Handler) handleNextEvent() bool {
	logrus.Info("Handler::handleNextEvent")
	event, quit := h.eventQueue.Get()
	if quit {
		return false
	}

	defer h.eventQueue.Done(event)

	h.handleEvent(event.(*Event))

	return true
}

func (h *Handler) withRetry(err error, event interface{}) {
	logrus.Info("Handler::withRetry")
	if err != nil {
		// TODO: Add Telemetry
		if h.eventQueue.NumRequeues(event) < 5 { // TODO: Make this configurable
			h.eventQueue.AddRateLimited(event)
		} else {
			h.eventQueue.Forget(event)
		}
	} // TODO: Add error logging
}
