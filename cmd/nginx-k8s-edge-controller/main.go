// Copyright 2023 f5 Inc. All rights reserved.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE file.

package main

import (
	config2 "github.com/nginxinc/kubernetes-nginx-ingress/internal/config"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/eventing"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/http"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/synchronization"
	nginx "github.com/nginxinc/nginx-plus-go-client/client"
	"github.com/sirupsen/logrus"
)

func main() {
	settings, err := config2.NewSettings()
	if err != nil {
		logrus.Error(err)
	}

	httpClient, err := http.NewHttpClient()
	if err != nil {
		logrus.Error(err)
	}

	nginxClient, err := nginx.NewNginxClient(httpClient, settings.NginxPlusHost)
	if err != nil {
		logrus.Error(err)
	}

	synchronizer, err := synchronization.NewNginxPlusSynchronizer(nginxClient)
	if err != nil {
		logrus.Error(err)
	}

	watcher, err := eventing.NewWatcher(synchronizer)
	if err != nil {
		logrus.Error(err)
	}

	err = watcher.Initialize()
	if err != nil {
		logrus.Error(err)
	}

	err = watcher.Watch()
	if err != nil {
		logrus.Error(err)
	}
}
