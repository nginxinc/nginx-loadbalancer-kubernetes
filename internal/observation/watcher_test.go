/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package observation

import (
	"context"
	"testing"

	"github.com/nginxinc/kubernetes-nginx-ingress/internal/configuration"
	"github.com/nginxinc/kubernetes-nginx-ingress/test/mocks"
	"k8s.io/client-go/kubernetes"
)

func TestWatcher_MustInitialize(t *testing.T) {
	t.Parallel()
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
