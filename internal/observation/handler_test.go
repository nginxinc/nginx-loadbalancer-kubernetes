/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package observation

import (
	"testing"

	"github.com/nginxinc/kubernetes-nginx-ingress/internal/configuration"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/core"
	"github.com/nginxinc/kubernetes-nginx-ingress/test/mocks"
	v1 "k8s.io/api/core/v1"
)

func TestHandler_AddsEventToSynchronizer(t *testing.T) {
	t.Parallel()
	synchronizer, handler := buildHandler()

	event := &core.Event{
		Type: core.Created,
		Service: &v1.Service{
			Spec: v1.ServiceSpec{
				Ports: []v1.ServicePort{
					{
						Name: "http-back",
					},
				},
			},
		},
	}

	handler.AddRateLimitedEvent(event)

	handler.handleNextEvent()

	if len(synchronizer.Events) != 1 {
		t.Errorf(`handler.AddRateLimitedEvent did not add the event to the queue`)
	}
}

func buildHandler() (
	*mocks.MockSynchronizer, *Handler,
) {
	eventQueue := &mocks.MockRateLimiter{}
	synchronizer := &mocks.MockSynchronizer{}

	handler := NewHandler(configuration.Settings{}, synchronizer, eventQueue)

	return synchronizer, handler
}
