/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package observation

import (
	"context"
	"fmt"
	"testing"

	"github.com/nginxinc/kubernetes-nginx-ingress/internal/configuration"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/core"
	"github.com/nginxinc/kubernetes-nginx-ingress/test/mocks"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/util/workqueue"
)

func TestHandler_AddsEventToSynchronizer(t *testing.T) {
	t.Parallel()
	_, _, synchronizer, handler, err := buildHandler()
	if err != nil {
		t.Errorf(`should have been no error, %v`, err)
	}

	event := &core.Event{
		Type: core.Created,
		Service: &v1.Service{
			Spec: v1.ServiceSpec{
				Ports: []v1.ServicePort{
					{
						Name: "nlk-back",
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
	*configuration.Settings,
	workqueue.RateLimitingInterface,
	*mocks.MockSynchronizer, *Handler, error,
) {
	settings, err := configuration.NewSettings(context.Background(), nil)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf(`should have been no error, %v`, err)
	}

	eventQueue := &mocks.MockRateLimiter{}
	synchronizer := &mocks.MockSynchronizer{}

	handler := NewHandler(settings, synchronizer, eventQueue)

	return settings, eventQueue, synchronizer, handler, nil
}
