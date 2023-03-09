// Copyright 2023 f5 Inc. All rights reserved.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE file.

package observation

import (
	"context"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/configuration"
	"github.com/nginxinc/kubernetes-nginx-ingress/test/mocks"
	"k8s.io/client-go/kubernetes"
	"testing"
)

func TestWatcher_MustInitialize(t *testing.T) {
	watcher, _ := buildWatcher()
	if err := watcher.Watch(); err == nil {
		t.Errorf("Expected error, got %s", err)
	}
}

func buildWatcher() (*Watcher, error) {
	k8sClient := &kubernetes.Clientset{}
	settings, _ := configuration.NewSettings(context.Background(), k8sClient)
	handler := &mocks.MockHandler{}

	return NewWatcher(settings, handler)
}
