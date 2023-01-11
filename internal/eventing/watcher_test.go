// Copyright 2023 f5 Inc. All rights reserved.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE file.

package eventing

import (
	config2 "github.com/nginxinc/kubernetes-nginx-ingress/internal/config"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/http"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/synchronization"
	nginx "github.com/nginxinc/nginx-plus-go-client/client"
	"github.com/sirupsen/logrus"
	"os"
	"testing"
)

func TestNewWatcher(t *testing.T) {
	const EnvVarName = "NGINX_PLUS_HOST"
	const ExpectedValue = "https://demo.nginx.com/api"

	defer os.Unsetenv(EnvVarName)
	os.Setenv(EnvVarName, ExpectedValue)

	settings, err := config2.NewSettings()
	if err != nil {
		t.Fatalf(`unable to create settings: %v`, err)
	}

	httpClient, err := http.NewHttpClient()
	if err != nil {
		t.Fatalf(`unable to create http client: %v`, err)
	}

	nginxClient, err := nginx.NewNginxClient(httpClient, settings.NginxPlusHost)
	if err != nil {
		logrus.Error(err)
	}

	synchronizer, err := synchronization.NewNginxPlusSynchronizer(nginxClient)
	if err != nil {
		t.Fatalf(`unable to create synchronizer: %v`, err)
	}

	watcher, err := NewWatcher(synchronizer)
	if err != nil {
		t.Fatalf(`failed to create watcher: %v`, err)
	}

	if watcher == nil {
		t.Fatal("watcher should not be nil")
	}
}
