/*
 * Copyright (c) 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package application

import (
	"log/slog"

	"github.com/nginxinc/kubernetes-nginx-ingress/internal/core"
)

// NullBorderClient is a BorderClient that does nothing.
// It serves only to prevent a panic if the BorderClient
// is not set correctly and errors from the factory methods are ignored.
type NullBorderClient struct{}

// NewNullBorderClient is the Factory function for creating a NullBorderClient
func NewNullBorderClient() (Interface, error) {
	return &NullBorderClient{}, nil
}

// Update logs a Warning. It is, after all, a NullObject Pattern implementation.
func (nbc *NullBorderClient) Update(_ *core.ServerUpdateEvent) error {
	slog.Warn("NullBorderClient.Update called")
	return nil
}

// Delete logs a Warning. It is, after all, a NullObject Pattern implementation.
func (nbc *NullBorderClient) Delete(_ *core.ServerUpdateEvent) error {
	slog.Warn("NullBorderClient.Delete called")
	return nil
}
