/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package mocks

import "github.com/nginxinc/kubernetes-nginx-ingress/internal/core"

type MockHandler struct {
}

func (h *MockHandler) AddRateLimitedEvent(_ *core.Event) {

}

func (h *MockHandler) Initialize() {

}

func (h *MockHandler) Run(_ <-chan struct{}) {

}

func (h *MockHandler) ShutDown() {

}
