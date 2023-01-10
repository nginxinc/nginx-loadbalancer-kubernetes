// Copyright 2023 f5 Inc. All rights reserved.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE file.

package eventing

import (
	"errors"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/synchronization"
)

type Watcher struct{}

func NewWatcher(*synchronization.Synchronizer) (*Watcher, error) {
	return nil, errors.New("not implemented")
}

func (*Watcher) Watch() {

}
