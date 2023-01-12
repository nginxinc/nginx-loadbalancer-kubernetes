// Copyright 2023 f5 Inc. All rights reserved.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE file.

package nginxplus

import (
	"errors"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/translation"
)

type DeletedTranslator struct{}

func NewDeletedTranslator() (*translation.Translator, error) {
	return nil, errors.New("not implemented")
}

func (dt DeletedTranslator) Translate() (interface{}, error) {
	return "Deleted", errors.New("not implemented")
}
