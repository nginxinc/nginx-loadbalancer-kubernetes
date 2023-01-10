// Copyright 2023 f5 Inc. All rights reserved.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE file.

package synchronization

import (
	"errors"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/config"
)

type NginxPlusSynchronizer struct{}

func NewNginxPlusSynchronizer(settings *config.Settings) (*Synchronizer, error) {
	return nil, errors.New("not implemented")
}

func (*NginxPlusSynchronizer) Synchronize() (interface{}, error) {
	return "", nil
}
