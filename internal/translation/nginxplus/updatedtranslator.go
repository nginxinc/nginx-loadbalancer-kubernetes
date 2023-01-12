// Copyright 2023 f5 Inc. All rights reserved.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE file.

package nginxplus

import (
	"errors"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/translation"
)

type UpdatedTranslator struct{}

func NewUpdatedTranslator() (*translation.Translator, error) {
	return nil, errors.New("not implemented")
}

func (ut UpdatedTranslator) Translate() (interface{}, error) {
	return "Updated", errors.New("not implemented")
}
