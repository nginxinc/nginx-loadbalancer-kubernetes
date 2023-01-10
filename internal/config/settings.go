// Copyright 2023 f5 Inc. All rights reserved.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE file.

package config

import (
	"errors"
	"os"
)

type Settings struct {
	NginxPlusHost string
}

func NewSettings() (*Settings, error) {
	config := new(Settings)

	config.NginxPlusHost = os.Getenv("NGINX_PLUS_HOST")
	if config.NginxPlusHost == "" {
		return nil, errors.New("the NGINX_PLUS_HOST variable is not defined. This is required")
	}

	return config, nil
}
