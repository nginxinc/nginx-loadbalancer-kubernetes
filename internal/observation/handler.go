// Copyright 2023 f5 Inc. All rights reserved.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE file.

package observation

import (
	"fmt"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/core"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/synchronization"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/translation"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/workqueue"
	"time"
)

const RateLimiterBase = time.Second * 2
const RateLimiterMax = time.Second * 60
const RetryCount = 5
const Threads = 1
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
	logrus.Debugf(`Handler::AddRateLimitedEvent: %#v`, event)
	h.eventQueue.AddRateLimited(event)
}

func (h *Handler) Initialize() {
	rateLimiter := workqueue.NewItemExponentialFailureRateLimiter(RateLimiterBase, RateLimiterMax)
	h.eventQueue = workqueue.NewNamedRateLimitingQueue(rateLimiter, WatcherQueueName)
}

func (h *Handler) Run(stopCh <-chan struct{}) {
	logrus.Debug("Handler::Run")

	for i := 0; i < Threads; i++ {
		go wait.Until(h.worker, 0, stopCh)
	}

	<-stopCh
}

func (h *Handler) ShutDown() {
	logrus.Debug("Handler::ShutDown")
	h.eventQueue.ShutDown()
}

func (h *Handler) handleEvent(e *core.Event) error {
	logrus.Debugf(`Handler::handleEvent: %#v`, e)
	// TODO: Add Telemetry

	event, err := translation.Translate(e)
	if err != nil {
		return fmt.Errorf(`Handler::handleEvent error translating: %v`, err)
	}

	h.synchronizer.AddRateLimitedEvent(event)

	return nil
}

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

func (h *Handler) worker() {
	for h.handleNextEvent() {
		// TODO: Add Telemetry
	}
}

func (h *Handler) withRetry(err error, event *core.Event) {
	logrus.Debug("Handler::withRetry")
	if err != nil {
		// TODO: Add Telemetry
		if h.eventQueue.NumRequeues(event) < RetryCount { // TODO: Make this configurable
			h.eventQueue.AddRateLimited(event)
			logrus.Infof(`Handler::withRetry: requeued event: %#v; error: %v`, event, err)
		} else {
			h.eventQueue.Forget(event)
			logrus.Warnf(`Handler::withRetry: event %#v has been dropped due to too many retries`, event)
		}
	} // TODO: Add error logging
}
