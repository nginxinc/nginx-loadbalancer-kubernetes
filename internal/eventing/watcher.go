// Copyright 2023 f5 Inc. All rights reserved.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE file.

package eventing

import (
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/synchronization"
)

type Watcher struct {
	synchronizer *synchronization.NginxPlusSynchronizer
}

func NewWatcher(synchronizer *synchronization.NginxPlusSynchronizer) (*Watcher, error) {
	return &Watcher{
		synchronizer: synchronizer,
	}, nil
}

func (*Watcher) Watch() {

}
