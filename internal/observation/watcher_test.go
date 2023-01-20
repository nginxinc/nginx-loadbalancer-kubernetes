// Copyright 2023 f5 Inc. All rights reserved.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE file.

package observation

import (
	"context"
	"os"
	"testing"
)

func TestNewWatcher(t *testing.T) {
	const EnvVarName = "NGINX_PLUS_HOST"
	const ExpectedValue = "https://demo.nginx.com/api"

	defer os.Unsetenv(EnvVarName)
	os.Setenv(EnvVarName, ExpectedValue)

	ctx := context.Background()

	handler := NewHandler()
	handler.Initialize()

	watcher, err := NewWatcher(ctx, handler)
	if err != nil {
		t.Fatalf(`failed to create watcher: %v`, err)
	}

	if watcher == nil {
		t.Fatal("watcher should not be nil")
	}
}
