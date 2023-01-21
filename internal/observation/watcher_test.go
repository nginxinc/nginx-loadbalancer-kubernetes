// Copyright 2023 f5 Inc. All rights reserved.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE file.

package observation

import (
	"context"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/synchronization"
	"os"
	"testing"
)

func TestNewWatcher(t *testing.T) {
	const EnvVarName = "NGINX_PLUS_HOST"
	const ExpectedValue = "https://demo.nginx.com/api"

	synchronizer, err := synchronization.NewSynchronizer()
	if err != nil {
		t.Fatalf(`should have been no error, %v`, err)
	}

	defer os.Unsetenv(EnvVarName)
	os.Setenv(EnvVarName, ExpectedValue)

	ctx := context.Background()

	handler := NewHandler(synchronizer)
	handler.Initialize()

	watcher, err := NewWatcher(ctx, handler)
	if err != nil {
		t.Fatalf(`failed to create watcher: %v`, err)
	}

	if watcher == nil {
		t.Fatal("watcher should not be nil")
	}
}
