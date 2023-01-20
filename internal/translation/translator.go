// Copyright 2023 f5 Inc. All rights reserved.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE file.

package translation

import (
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/core"
)

type Translator interface {
	Translate(event core.Event) (interface{}, error)
}
