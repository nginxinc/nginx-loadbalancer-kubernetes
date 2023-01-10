// Copyright 2023 f5 Inc. All rights reserved.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE file.

package main

import (
	config2 "github.com/nginxinc/kubernetes-nginx-ingress/internal/config"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/eventing"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/synchronization"
	"github.com/sirupsen/logrus"
)

func main() {
	config, err := config2.NewSettings()
	if err != nil {
		logrus.Error(err)
	}

	synchronizer, err := synchronization.NewNginxPlusSynchronizer(config)
	if err != nil {
		logrus.Error(err)
	}

	watcher, err := eventing.NewWatcher(synchronizer)
	if err != nil {
		logrus.Error(err)
	}

	watcher.Watch()
}
