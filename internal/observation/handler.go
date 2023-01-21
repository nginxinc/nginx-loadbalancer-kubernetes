// Copyright 2023 f5 Inc. All rights reserved.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE file.

package observation

import (
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/core"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/synchronization"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/translation"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/util/workqueue"
)

const WatcherQueueName = `nec-handler`

type Handler struct {
	eventQueue   workqueue.RateLimitingInterface
	synchronizer *synchronization.Synchronizer
}

func NewHandler(synchronizer *synchronization.Synchronizer) *Handler {
	return &Handler{
		synchronizer: synchronizer,
	}
}

func (h *Handler) AddRateLimitedEvent(event *core.Event) {
	logrus.Infof(`Handler::AddRateLimitedEvent: %#v`, event)
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

func (h *Handler) handleEvent(e *core.Event) {
	logrus.Info("Handler::handleEvent")
	// TODO: Add Telemetry

	event, err := translation.Translate(e)
	if err != nil {
		logrus.Errorf(`Handler::handleEvent error translating: %v`, err)
	} else {
		h.synchronizer.AddRateLimitedEvent(event)
	}
}

func (h *Handler) handleNextEvent() bool {
	logrus.Info("Handler::handleNextEvent")
	event, quit := h.eventQueue.Get()
	if quit {
		return false
	}

	defer h.eventQueue.Done(event)

	// TODO: use withRetry
	h.handleEvent(event.(*core.Event))

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
