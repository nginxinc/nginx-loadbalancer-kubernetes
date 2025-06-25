/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package mocks

import "github.com/nginxinc/kubernetes-nginx-ingress/internal/core"

type MockSynchronizer struct {
	Events []core.ServerUpdateEvent
}

func (s *MockSynchronizer) AddEvents(events core.ServerUpdateEvents) {
	for _, event := range events {
		s.Events = append(s.Events, *event)
	}
}

func (s *MockSynchronizer) AddEvent(event *core.ServerUpdateEvent) {
	s.Events = append(s.Events, *event)
}

func (s *MockSynchronizer) Initialize() error {
	return nil
}

func (s *MockSynchronizer) Run(stopCh <-chan struct{}) {
	<-stopCh
}

func (s *MockSynchronizer) ShutDown() {

}
